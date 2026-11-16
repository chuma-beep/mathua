package storage

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"time"
)

// Weekly leagues.
//
// A league is a fixed cohort for one week: membership is written once, at the
// start of the week, and never changes while the competition runs. Standings are
// read live from the XP ledger. That split is what makes the system safe to
// retry: joining a cohort and closing one are both idempotent writes keyed by
// (week, student), so a crash midway through a transition repeats work instead
// of moving somebody twice.
//
// The previous model kept a single `students.league` string and recomputed
// promotions on every read. Retrying after a partial failure therefore promoted
// a different set of people, because the group it ranked had already shrunk.
// Those columns remain in `students` for compatibility, but nothing reads them.

// LeagueTiers, lowest first. Diamond is the ceiling: its top finishers stay.
var LeagueTiers = []string{"bronze", "silver", "gold", "platinum", "diamond"}

const (
	// LeagueSize is the target cohort size. Leagues fill to this before a new
	// one opens, and the last league of a tier may be short.
	LeagueSize = 30
	// PromoteCount and DemoteCount cap the movers per league per week. Short
	// leagues move proportionally fewer: see moversFor.
	PromoteCount = 5
	DemoteCount  = 5
	// minLeagueToMove keeps tiny cohorts still. With two or three learners,
	// moving half of them would shuffle the same people between tiers every
	// week without measuring anything.
	minLeagueToMove = 6
)

// LeagueRow is one learner inside a cohort, already ranked.
type LeagueRow struct {
	StudentID string
	Name      string
	Username  string
	AvatarURL string
	Settings  string
	XP        int
	Rank      int
}

// LeagueMembership is where a learner sits this week.
type LeagueMembership struct {
	Tier     string
	LeagueNo int
	Rank     int
	XP       int
	Size     int
}

// LeagueCohort is one league's shape: which tier, which number inside it, and
// how many learners hold a seat. The standings are read separately, because a
// week can hold dozens of leagues and an administrator browsing the pyramid
// does not need every roster at once.
type LeagueCohort struct {
	Tier     string
	LeagueNo int
	Size     int
}

// moversFor decides how many learners leave a league upwards and downwards.
//
// The cap is five, so a full league behaves the way the product documents. A
// short league moves a third of itself rather than half, which keeps promotion
// and relegation disjoint without a special case: two thirds of a cohort cannot
// fill both ends. Cohorts below minLeagueToMove do not move at all.
// LeagueMovers is how many learners move between tiers from a league of this
// size. It is exported so the API states the rule it actually applies rather
// than restating the constant in a second place.
func LeagueMovers(size int) int { return moversFor(size) }

func moversFor(size int) int {
	if size < minLeagueToMove {
		return 0
	}
	k := size / 3
	if k > PromoteCount {
		k = PromoteCount
	}
	return k
}

// nextTier returns the tier above or below, clamped at the ends.
func nextTier(tier string, delta int) string {
	i := tierIndex(tier)
	if i < 0 {
		i = 0
	}
	j := i + delta
	if j < 0 {
		j = 0
	}
	if j >= len(LeagueTiers) {
		j = len(LeagueTiers) - 1
	}
	return LeagueTiers[j]
}

func tierIndex(tier string) int {
	for i, t := range LeagueTiers {
		if t == tier {
			return i
		}
	}
	return -1
}

// validTier normalises a stored tier so a corrupt row cannot produce an
// unrenderable league.
func validTier(tier string) string {
	t := strings.ToLower(strings.TrimSpace(tier))
	if tierIndex(t) < 0 {
		return LeagueTiers[0]
	}
	return t
}

// leagueWeekKey is the Monday that opens the competition containing now. It
// mirrors leaderboard.WeekStart without importing that package, because the
// leaderboard already imports storage.
func LeagueWeekKey(now time.Time) string {
	d := now.UTC()
	offset := (int(d.Weekday()) + 6) % 7
	return d.AddDate(0, 0, -offset).Format("2006-01-02")
}

// WeekFromString parses a competition week key (a Monday, YYYY-MM-DD, UTC).
func WeekFromString(week string) (time.Time, error) {
	t, err := time.ParseInLocation("2006-01-02", week, time.UTC)
	if err != nil {
		return time.Time{}, fmt.Errorf("week key %q: %w", week, err)
	}
	if t.Weekday() != time.Monday {
		return time.Time{}, fmt.Errorf("week key %q is not a Monday", week)
	}
	return t, nil
}

type leagueCohort struct {
	Tier     string
	LeagueNo int
	Members  []leagueMember
}

type leagueMember struct {
	StudentID string
	Name      string
	XP        int
	Movement  string
	FromTier  string
}

// planWeek maps last week's cohorts to this week's tiers.
//
// Cohorts arrive with their final XP and a deterministic order (XP desc, name,
// id). The result is one tier per learner. Pure and deterministic: the same
// input always yields the same mapping, so a retried transition cannot quietly
// move somebody else.
func planWeek(cohorts []leagueCohort) (map[string]string, map[string]string) {
	tierOf := make(map[string]string, len(cohorts)*LeagueSize)
	moveOf := make(map[string]string, len(cohorts)*LeagueSize)
	for _, c := range cohorts {
		k := moversFor(len(c.Members))
		for i, m := range c.Members {
			from := validTier(c.Tier)
			tier := from
			switch {
			case k > 0 && i < k:
				tier = nextTier(tier, +1)
			case k > 0 && i >= len(c.Members)-k:
				tier = nextTier(tier, -1)
			}
			tierOf[m.StudentID] = tier
			switch {
			case tierIndex(tier) > tierIndex(from):
				moveOf[m.StudentID] = "up"
			case tierIndex(tier) < tierIndex(from):
				moveOf[m.StudentID] = "down"
			default:
				moveOf[m.StudentID] = "stay"
			}
			c.Members[i].FromTier = from
		}
	}
	return tierOf, moveOf
}

// cutCohorts orders one tier's learners and cuts them into leagues of
// LeagueSize. League numbers restart at 1 inside every tier.
func cutCohorts(tier string, members []leagueMember) []leagueCohort {
	leagueOrder(members)
	var out []leagueCohort
	for i := 0; i < len(members); i += LeagueSize {
		end := i + LeagueSize
		if end > len(members) {
			end = len(members)
		}
		out = append(out, leagueCohort{
			Tier:     tier,
			LeagueNo: len(out) + 1,
			Members:  append([]leagueMember(nil), members[i:end]...),
		})
	}
	return out
}

// leagueOrder sorts members by the documented deterministic order: XP
// descending, then display name, then student id. The last key is what keeps
// two learners with the same name from swapping places between renders.
func leagueOrder(m []leagueMember) {
	sort.SliceStable(m, func(i, j int) bool {
		if m[i].XP != m[j].XP {
			return m[i].XP > m[j].XP
		}
		if m[i].Name != m[j].Name {
			return m[i].Name < m[j].Name
		}
		return m[i].StudentID < m[j].StudentID
	})
}

// leagueCohortRow is the flat shape the stores return from their standings
// query, already ordered by tier, league, XP, name and id.
type leagueCohortRow struct {
	Tier      string
	LeagueNo  int
	StudentID string
	Name      string
	XP        int
}

// groupCohorts rebuilds cohorts from flat rows.
func groupCohorts(rows []leagueCohortRow) []leagueCohort {
	var out []leagueCohort
	for _, r := range rows {
		tier := validTier(r.Tier)
		last := len(out) - 1
		if last < 0 || out[last].Tier != tier || out[last].LeagueNo != r.LeagueNo {
			out = append(out, leagueCohort{Tier: tier, LeagueNo: r.LeagueNo})
			last = len(out) - 1
		}
		out[last].Members = append(out[last].Members, leagueMember{
			StudentID: r.StudentID,
			Name:      r.Name,
			XP:        r.XP,
		})
	}
	return out
}

// rankInCohort fills Rank for display: equal XP shares a rank and the next
// distinct score skips the tied positions (1, 1, 3). The row order itself stays
// the deterministic one, so the policy is about what the number says, not about
// who appears first.
func rankInCohort(rows []LeagueRow) {
	for i := range rows {
		if i > 0 && rows[i].XP == rows[i-1].XP {
			rows[i].Rank = rows[i-1].Rank
			continue
		}
		rows[i].Rank = i + 1
	}
}

// BackfillLeagueMemberships seats every non-guest account with positive
// lifetime XP that has no row in week. It exists because membership used to be
// created lazily by MyLeague, so an account that earned XP but never opened
// My League or its profile had no tier. The call is idempotent and returns the
// number of accounts it seated.
func (s *SQLiteStore) BackfillLeagueMemberships(week string) (int, error) {
	return backfillLeagueMemberships(s.db,
		`SELECT s.id FROM students s
		 WHERE s.xp_total > 0
		   AND `+rankedStudentWhere+`
		   AND NOT EXISTS (SELECT 1 FROM league_members lm
		                   WHERE lm.student_id = s.id AND lm.week = ?)
		 ORDER BY s.id`, week, s.EnsureLeagueMembership)
}

func (s *PostgresStore) BackfillLeagueMemberships(week string) (int, error) {
	return backfillLeagueMemberships(s.db, rebindDollars(
		`SELECT s.id FROM students s
		 WHERE s.xp_total > 0
		   AND `+rankedStudentWhere+`
		   AND NOT EXISTS (SELECT 1 FROM league_members lm
		                   WHERE lm.student_id = s.id AND lm.week = $1)
		 ORDER BY s.id`), week, s.EnsureLeagueMembership)
}

// backfillLeagueMemberships reads the ids before seating anyone: SQLite holds a
// read lock for the life of an open rows cursor, and EnsureLeagueMembership
// writes. Closing the cursor first also keeps one failed insert from aborting
// the rest of the batch.
func backfillLeagueMemberships(db *sql.DB, query, week string, ensure func(string, string) error) (int, error) {
	rows, err := db.Query(query, week)
	if err != nil {
		return 0, fmt.Errorf("list league backfill: %w", err)
	}
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return 0, fmt.Errorf("scan league backfill: %w", err)
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return 0, fmt.Errorf("league backfill rows: %w", err)
	}
	rows.Close()

	seated := 0
	for _, id := range ids {
		if err := ensure(id, week); err != nil {
			return seated, fmt.Errorf("seat %s: %w", id, err)
		}
		seated++
	}
	return seated, nil
}

// BackfillLeagueMovement reconstructs movement for weeks written before the
// column existed. A learner who was absent for a week returns as new, not as a
// promotion, so the comparison only runs between adjacent weeks.
func (s *SQLiteStore) BackfillLeagueMovement() (int, error) {
	return backfillLeagueMovement(s.db, `SELECT week, student_id, tier FROM league_members ORDER BY student_id, week`,
		func(week, id, movement, fromTier string) error {
			_, err := s.db.Exec(`UPDATE league_members SET movement = ?, from_tier = ? WHERE week = ? AND student_id = ?`, movement, fromTier, week, id)
			return err
		})
}

func (s *PostgresStore) BackfillLeagueMovement() (int, error) {
	return backfillLeagueMovement(s.db, `SELECT week, student_id, tier FROM league_members ORDER BY student_id, week`,
		func(week, id, movement, fromTier string) error {
			_, err := s.db.Exec(`UPDATE league_members SET movement = $1, from_tier = $2 WHERE week = $3 AND student_id = $4`, movement, fromTier, week, id)
			return err
		})
}

func backfillLeagueMovement(db *sql.DB, query string, update func(week, id, movement, fromTier string) error) (int, error) {
	rows, err := db.Query(query)
	if err != nil {
		return 0, fmt.Errorf("list league movement: %w", err)
	}
	type seat struct{ week, id, tier string }
	var seats []seat
	for rows.Next() {
		var s seat
		if err := rows.Scan(&s.week, &s.id, &s.tier); err != nil {
			rows.Close()
			return 0, fmt.Errorf("scan league movement: %w", err)
		}
		seats = append(seats, s)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return 0, fmt.Errorf("league movement rows: %w", err)
	}
	rows.Close()

	lastWeek := map[string]string{}
	lastTier := map[string]string{}
	changed := 0
	for _, s := range seats {
		movement := "new"
		fromTier := ""
		if prevWeek, ok := lastWeek[s.id]; ok && prevWeek == shiftWeekKey(s.week, -7) {
			fromTier = lastTier[s.id]
			switch {
			case tierIndex(s.tier) > tierIndex(lastTier[s.id]):
				movement = "up"
			case tierIndex(s.tier) < tierIndex(lastTier[s.id]):
				movement = "down"
			default:
				movement = "stay"
			}
		}
		if err := update(s.week, s.id, movement, fromTier); err != nil {
			return changed, fmt.Errorf("update league movement: %w", err)
		}
		changed++
		lastWeek[s.id] = s.week
		lastTier[s.id] = s.tier
	}
	return changed, nil
}

func shiftWeekKey(week string, days int) string {
	t, err := time.Parse("2006-01-02", week)
	if err != nil {
		return ""
	}
	return t.AddDate(0, 0, days).Format("2006-01-02")
}
