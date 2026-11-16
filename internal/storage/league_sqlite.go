package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"time"
)

// Weekly leagues (SQLite). The Postgres twin lives in league_postgres.go and
// mirrors this statement for statement.

// ResolveLeagueWeek writes the membership for `week`, once.
//
// It is called on every read of a league, so the first statement is the marker:
// a week that has been resolved returns immediately, and a second process that
// raced the first loses the unique insert and does no work. Everything after the
// marker happens in the same transaction, so a crash leaves either the whole
// cohort set or none of it.
//
// prevFrom and prevTo bound the XP window of the week being closed.
func (s *SQLiteStore) ResolveLeagueWeek(week string, prevFrom, prevTo time.Time) error {
	start, err := WeekFromString(week)
	if err != nil {
		return err
	}
	prev := start.AddDate(0, 0, -7).Format("2006-01-02")

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("resolve league week: %w", err)
	}
	defer tx.Rollback()

	res, err := tx.Exec(
		`INSERT INTO league_weeks (week, closed_at) VALUES (?, ?)
		 ON CONFLICT(week) DO NOTHING`,
		week, time.Now().UTC().UnixNano())
	if err != nil {
		return fmt.Errorf("claim league week: %w", err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return nil // another request closed it first
	}

	rows, err := tx.Query(
		`SELECT m.tier, m.league_no, m.student_id, s.name,
		        COALESCE((SELECT SUM(e.amount) FROM xp_events e
		                  WHERE e.student_id = m.student_id
		                    AND e.created_at >= ? AND e.created_at < ?), 0) AS xp
		 FROM league_members m
		 JOIN students s ON s.id = m.student_id
		 WHERE m.week = ?
		 ORDER BY m.tier, m.league_no, xp DESC, s.name ASC, m.student_id ASC`,
		prevFrom.UnixNano(), prevTo.UnixNano(), prev)
	if err != nil {
		return fmt.Errorf("read previous leagues: %w", err)
	}
	var flat []leagueCohortRow
	for rows.Next() {
		var r leagueCohortRow
		if err := rows.Scan(&r.Tier, &r.LeagueNo, &r.StudentID, &r.Name, &r.XP); err != nil {
			rows.Close()
			return fmt.Errorf("scan previous league: %w", err)
		}
		flat = append(flat, r)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return fmt.Errorf("read previous leagues: %w", err)
	}

	if len(flat) == 0 {
		// First week of the competition. Seed from lifetime XP so the very first
		// cohorts are not ordered by chance, but store no weekly score: the
		// weekly competition starts at zero for everyone and the UI says so.
		if err := s.bootstrapLeagues(tx, week); err != nil {
			return err
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit league seed: %w", err)
		}
		return nil
	}

	cohorts := groupCohorts(flat)
	tierOf, moveOf := planWeek(cohorts)

	// Record the final XP and position of the week being closed. This is the
	// historical result, so it must not depend on later ledger entries.
	for _, c := range cohorts {
		for i, m := range c.Members {
			if _, err := tx.Exec(
				`UPDATE league_members SET xp = ?, rank = ? WHERE week = ? AND student_id = ?`,
				m.XP, i+1, prev, m.StudentID); err != nil {
				return fmt.Errorf("finalise league member: %w", err)
			}
		}
	}

	byTier := make(map[string][]leagueMember)
	for _, c := range cohorts {
		for _, m := range c.Members {
			tier := tierOf[m.StudentID]
			m.Movement = moveOf[m.StudentID]
			byTier[tier] = append(byTier[tier], m)
		}
	}
	var next []leagueCohort
	for _, tier := range LeagueTiers {
		next = append(next, cutCohorts(tier, byTier[tier])...)
	}
	if err := insertCohorts(tx, week, next, rebindQuestion); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit league week: %w", err)
	}
	return nil
}

// bootstrapLeagues seeds the first competition week from lifetime XP.
func (s *SQLiteStore) bootstrapLeagues(tx *sql.Tx, week string) error {
	rows, err := tx.Query(
		`SELECT id, name, COALESCE(username, ''), xp_total
		 FROM students
		 WHERE xp_total > 0 AND id NOT LIKE 'guest_%'
		 ORDER BY xp_total DESC, name ASC, id ASC`)
	if err != nil {
		return fmt.Errorf("seed leagues: %w", err)
	}
	var members []leagueMember
	for rows.Next() {
		var id, name, username string
		var xp int
		if err := rows.Scan(&id, &name, &username, &xp); err != nil {
			rows.Close()
			return fmt.Errorf("scan seed learner: %w", err)
		}
		members = append(members, leagueMember{StudentID: id, Name: name, XP: xp})
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return fmt.Errorf("seed leagues: %w", err)
	}
	return insertCohorts(tx, week, cutCohorts(LeagueTiers[0], members), rebindQuestion)
}

// insertCohorts writes one week's seats. Both engines call it, so the caller
// supplies the placeholder spelling: SQLite takes `?`, Postgres takes `$n`. The
// statement used to hardcode `?`, which made a Postgres bootstrap fail with a
// syntax error at the first comma, so a league week could never form there.
func insertCohorts(tx *sql.Tx, week string, cohorts []leagueCohort, rebind func(string) string) error {
	const insert = `INSERT INTO league_members (week, student_id, tier, league_no, xp, rank, movement, from_tier)
		 VALUES (?, ?, ?, ?, 0, 0, ?, ?)
		 ON CONFLICT(week, student_id) DO NOTHING`
	query := rebind(insert)
	for _, c := range cohorts {
		for _, m := range c.Members {
			movement := m.Movement
			if movement == "" {
				movement = "new"
			}
			if _, err := tx.Exec(query, week, m.StudentID, c.Tier, c.LeagueNo, movement, m.FromTier); err != nil {
				return fmt.Errorf("insert league member: %w", err)
			}
		}
	}
	return nil
}

// EnsureLeagueMembership places a learner who became eligible after the week
// was resolved. It is idempotent: an existing row returns untouched, and the
// (week, student) primary key settles a race between two tabs.
func (s *SQLiteStore) EnsureLeagueMembership(studentID, week string) error {
	if studentID == "" {
		return nil
	}
	var exists int
	err := s.db.QueryRow(
		`SELECT 1 FROM league_members WHERE week = ? AND student_id = ?`,
		week, studentID).Scan(&exists)
	if err == nil {
		return nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("check league membership: %w", err)
	}

	var xpTotal, ranked int
	if err := s.db.QueryRow(`SELECT s.xp_total, CASE WHEN `+rankedStudentWhere+` THEN 1 ELSE 0 END
		FROM students s WHERE s.id = ?`, studentID).
		Scan(&xpTotal, &ranked); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("check league eligibility: %w", err)
	}
	// A learner enters the league the week they first earn XP. Before that they
	// are not competing, and a zero-XP row would be a participant with nothing
	// to rank. An adopted guest keeps its guest_ id but has credentials, so the
	// ranked predicate, not the id prefix, decides.
	if xpTotal <= 0 || ranked == 0 {
		return nil
	}

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("join league: %w", err)
	}
	defer tx.Rollback()

	if err := tx.QueryRow(
		`SELECT 1 FROM league_members WHERE week = ? AND student_id = ?`,
		week, studentID).Scan(&exists); err == nil {
		return nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("recheck league membership: %w", err)
	}

	no := 0
	err = tx.QueryRow(
		`SELECT league_no FROM league_members
		 WHERE week = ? AND tier = ?
		 GROUP BY league_no HAVING COUNT(*) < ?
		 ORDER BY league_no LIMIT 1`,
		week, LeagueTiers[0], LeagueSize).Scan(&no)
	if errors.Is(err, sql.ErrNoRows) {
		if err := tx.QueryRow(
			`SELECT COALESCE(MAX(league_no), 0) + 1 FROM league_members
			 WHERE week = ? AND tier = ?`,
			week, LeagueTiers[0]).Scan(&no); err != nil {
			return fmt.Errorf("open a league: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("find a league with room: %w", err)
	}

	if _, err := tx.Exec(
		`INSERT INTO league_members (week, student_id, tier, league_no, xp, rank, movement, from_tier)
		 VALUES (?, ?, ?, ?, 0, 0, 'new', '')
		 ON CONFLICT(week, student_id) DO NOTHING`,
		week, studentID, LeagueTiers[0], no); err != nil {
		return fmt.Errorf("insert league member: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit league membership: %w", err)
	}
	return nil
}

// LeagueMembershipOf returns the cohort a learner belongs to this week.
func (s *SQLiteStore) LeagueMembershipOf(studentID, week string) (string, int, bool, error) {
	var tier string
	var no int
	err := s.db.QueryRow(
		`SELECT tier, league_no FROM league_members WHERE week = ? AND student_id = ?`,
		week, studentID).Scan(&tier, &no)
	if errors.Is(err, sql.ErrNoRows) {
		return "", 0, false, nil
	}
	if err != nil {
		return "", 0, false, fmt.Errorf("read league membership: %w", err)
	}
	return validTier(tier), no, true, nil
}

// LeagueRoster returns one cohort's live standings for the week.
func (s *SQLiteStore) LeagueRoster(week, tier string, leagueNo int, from, to time.Time) ([]LeagueRow, error) {
	rows, err := s.db.Query(
		`SELECT m.student_id, s.name, COALESCE(s.username, ''), s.avatar_url, s.settings,
		        COALESCE((SELECT SUM(e.amount) FROM xp_events e
		                  WHERE e.student_id = m.student_id
		                    AND e.created_at >= ? AND e.created_at < ?), 0) AS xp
		 FROM league_members m
		 JOIN students s ON s.id = m.student_id
		 WHERE m.week = ? AND m.tier = ? AND m.league_no = ?
		 ORDER BY xp DESC, s.name ASC, m.student_id ASC`,
		from.UnixNano(), to.UnixNano(), week, validTier(tier), leagueNo)
	if err != nil {
		return nil, fmt.Errorf("read league roster: %w", err)
	}
	defer rows.Close()
	var out []LeagueRow
	for rows.Next() {
		var r LeagueRow
		if err := rows.Scan(&r.StudentID, &r.Name, &r.Username, &r.AvatarURL, &r.Settings, &r.XP); err != nil {
			return nil, fmt.Errorf("scan league roster: %w", err)
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read league roster: %w", err)
	}
	rankInCohort(out)
	return out, nil
}

// LeagueCohorts lists the week's leagues in pyramid order: bronze first, then
// by league number inside each tier. The database groups and counts; the tier
// order comes from LeagueTiers rather than from an alphabetical sort, which
// would put bronze between diamond and gold.
func (s *SQLiteStore) LeagueCohorts(week string) ([]LeagueCohort, error) {
	rows, err := s.db.Query(
		`SELECT tier, league_no, COUNT(*) FROM league_members
		 WHERE week = ?
		 GROUP BY tier, league_no`, week)
	if err != nil {
		return nil, fmt.Errorf("read league cohorts: %w", err)
	}
	defer rows.Close()
	var out []LeagueCohort
	for rows.Next() {
		var c LeagueCohort
		if err := rows.Scan(&c.Tier, &c.LeagueNo, &c.Size); err != nil {
			return nil, fmt.Errorf("scan league cohort: %w", err)
		}
		c.Tier = validTier(c.Tier)
		out = append(out, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read league cohorts: %w", err)
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Tier != out[j].Tier {
			return tierIndex(out[i].Tier) < tierIndex(out[j].Tier)
		}
		return out[i].LeagueNo < out[j].LeagueNo
	})
	return out, nil
}
