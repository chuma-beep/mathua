package storage

import (
	"database/sql"
	"fmt"
	"time"
)

// LeagueRange is the admin read model for more than one competition week.
// Movement is read from the seat rows, so a promotion is counted where it
// happened rather than recomputed from two tier snapshots.
type LeagueRange struct {
	From       string               `json:"from"`
	To         string               `json:"to"`
	Group      string               `json:"group"`
	Weeks      int                  `json:"weeks"`
	WeekKeys   []string             `json:"weeks_available"`
	Seats      int                  `json:"seats"`
	Unique     int                  `json:"unique_learners"`
	Up         int                  `json:"promoted"`
	Down       int                  `json:"demoted"`
	Stay       int                  `json:"stayed"`
	New        int                  `json:"new"`
	UniqueUp   int                  `json:"unique_promoted"`
	UniqueDown int                  `json:"unique_demoted"`
	Left       int                  `json:"left"`
	Matrix     []LeagueTransition   `json:"matrix"`
	Tiers      []LeagueRangeTier    `json:"tiers"`
	Buckets    []LeagueRangeBucket  `json:"buckets"`
	Top        []LeagueRangeLearner `json:"top_range"`
	Latest     []LeagueRangeLearner `json:"top_latest"`
}

type LeagueRangeTier struct {
	Tier    string `json:"tier"`
	Leagues int    `json:"leagues"`
	Seats   int    `json:"seats"`
	Unique  int    `json:"unique_learners"`
}

type LeagueRangeBucket struct {
	Key    string `json:"key"`
	Label  string `json:"label"`
	Weeks  int    `json:"weeks"`
	Seats  int    `json:"seats"`
	Unique int    `json:"unique_learners"`
	Up     int    `json:"promoted"`
	Down   int    `json:"demoted"`
	Stay   int    `json:"stayed"`
	New    int    `json:"new"`
}

type LeagueTransition struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Count int    `json:"count"`
}

type LeagueRangeLearner struct {
	StudentID string `json:"student_id"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	XP        int    `json:"xp"`
	Tier      string `json:"tier"`
}

func (s *SQLiteStore) LeagueRange(from, to, group string) (*LeagueRange, error) {
	return leagueRange(s.db, from, to, group, func(q string) string { return q })
}

func (s *PostgresStore) LeagueRange(from, to, group string) (*LeagueRange, error) {
	return leagueRange(s.db, from, to, group, rebindDollars)
}

func leagueRange(db *sql.DB, from, to, group string, rebind func(string) string) (*LeagueRange, error) {
	out := &LeagueRange{From: from, To: to, Group: group, Tiers: []LeagueRangeTier{}, Buckets: []LeagueRangeBucket{}, Matrix: []LeagueTransition{}, Top: []LeagueRangeLearner{}, Latest: []LeagueRangeLearner{}, WeekKeys: []string{}}
	bucketExpr := "week"
	if group == "month" {
		bucketExpr = "substr(week, 1, 7)"
	}
	rows, err := db.Query(rebind(`SELECT `+bucketExpr+` AS bucket, COUNT(DISTINCT week),
		COUNT(DISTINCT tier || '/' || league_no), COUNT(*), COUNT(DISTINCT student_id),
		SUM(CASE WHEN movement = 'up' THEN 1 ELSE 0 END),
		SUM(CASE WHEN movement = 'down' THEN 1 ELSE 0 END),
		SUM(CASE WHEN movement = 'stay' THEN 1 ELSE 0 END),
		SUM(CASE WHEN movement = 'new' THEN 1 ELSE 0 END)
		FROM league_members WHERE week >= ? AND week <= ?
		GROUP BY bucket ORDER BY bucket`), from, to)
	if err != nil {
		return nil, fmt.Errorf("league range buckets: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var b LeagueRangeBucket
		var leagues int
		if err := rows.Scan(&b.Key, &b.Weeks, &leagues, &b.Seats, &b.Unique, &b.Up, &b.Down, &b.Stay, &b.New); err != nil {
			return nil, fmt.Errorf("scan league range bucket: %w", err)
		}
		b.Label = b.Key
		out.Buckets = append(out.Buckets, b)
		out.Weeks += b.Weeks
		out.Seats += b.Seats
		out.Up += b.Up
		out.Down += b.Down
		out.Stay += b.Stay
		out.New += b.New
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("league range bucket rows: %w", err)
	}
	wrows, err := db.Query(rebind(`SELECT DISTINCT week FROM league_members WHERE week >= ? AND week <= ? ORDER BY week`), from, to)
	if err != nil {
		return nil, fmt.Errorf("league range weeks: %w", err)
	}
	for wrows.Next() {
		var w string
		if err := wrows.Scan(&w); err != nil {
			wrows.Close()
			return nil, fmt.Errorf("scan league range week: %w", err)
		}
		out.WeekKeys = append(out.WeekKeys, w)
	}
	wrows.Close()
	if err := db.QueryRow(rebind(`SELECT COUNT(DISTINCT student_id) FROM league_members WHERE week >= ? AND week <= ?`), from, to).Scan(&out.Unique); err != nil {
		return nil, fmt.Errorf("league range uniques: %w", err)
	}
	trows, err := db.Query(rebind(`SELECT tier, COUNT(*), COUNT(DISTINCT student_id), COUNT(DISTINCT league_no)
		FROM league_members WHERE week >= ? AND week <= ? GROUP BY tier`), from, to)
	if err != nil {
		return nil, fmt.Errorf("league range tiers: %w", err)
	}
	defer trows.Close()
	for trows.Next() {
		var t LeagueRangeTier
		if err := trows.Scan(&t.Tier, &t.Seats, &t.Unique, &t.Leagues); err != nil {
			return nil, fmt.Errorf("scan league range tier: %w", err)
		}
		out.Tiers = append(out.Tiers, t)
	}
	if err := trows.Err(); err != nil {
		return nil, fmt.Errorf("league range tier rows: %w", err)
	}

	if err := db.QueryRow(rebind(`SELECT COUNT(DISTINCT student_id) FROM league_members WHERE week >= ? AND week <= ? AND movement = 'up'`), from, to).Scan(&out.UniqueUp); err != nil {
		return nil, fmt.Errorf("league range unique promoted: %w", err)
	}
	if err := db.QueryRow(rebind(`SELECT COUNT(DISTINCT student_id) FROM league_members WHERE week >= ? AND week <= ? AND movement = 'down'`), from, to).Scan(&out.UniqueDown); err != nil {
		return nil, fmt.Errorf("league range unique demoted: %w", err)
	}
	mrows, err := db.Query(rebind(`SELECT from_tier, tier, COUNT(*) FROM league_members
		WHERE week >= ? AND week <= ? AND from_tier <> '' GROUP BY from_tier, tier`), from, to)
	if err != nil {
		return nil, fmt.Errorf("league range matrix: %w", err)
	}
	for mrows.Next() {
		var cell LeagueTransition
		if err := mrows.Scan(&cell.From, &cell.To, &cell.Count); err != nil {
			mrows.Close()
			return nil, fmt.Errorf("scan league range matrix: %w", err)
		}
		out.Matrix = append(out.Matrix, cell)
	}
	mrows.Close()
	out.Left, err = leagueRangeLeft(db, rebind, from, to)
	if err != nil {
		return nil, err
	}

	start, end := weekBounds(from, to)
	out.Top, err = leagueRangeTop(db, rebind, start, end, to)
	if err != nil {
		return nil, err
	}
	lastStart, lastEnd := weekBounds(to, to)
	out.Latest, err = leagueRangeTop(db, rebind, lastStart, lastEnd, to)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func weekBounds(from, to string) (time.Time, time.Time) {
	a, _ := time.Parse("2006-01-02", from)
	b, _ := time.Parse("2006-01-02", to)
	return a.UTC(), b.AddDate(0, 0, 7).UTC()
}

func leagueRangeTop(db *sql.DB, rebind func(string) string, since, until time.Time, week string) ([]LeagueRangeLearner, error) {
	rows, err := db.Query(rebind(`SELECT s.id, s.name, COALESCE(s.username, ''), SUM(e.amount) AS xp,
		COALESCE((SELECT lm.tier FROM league_members lm
		          WHERE lm.student_id = s.id AND lm.week <= ?
		          ORDER BY lm.week DESC LIMIT 1), '')
		FROM xp_events e JOIN students s ON s.id = e.student_id
		WHERE e.created_at >= ? AND e.created_at < ? AND `+rankedStudentWhere+`
		GROUP BY s.id HAVING SUM(e.amount) > 0
		ORDER BY xp DESC, s.name ASC, s.id ASC LIMIT 10`), week, since.UnixNano(), until.UnixNano())
	if err != nil {
		return nil, fmt.Errorf("league range top: %w", err)
	}
	defer rows.Close()
	out := []LeagueRangeLearner{}
	for rows.Next() {
		var r LeagueRangeLearner
		if err := rows.Scan(&r.StudentID, &r.Name, &r.Username, &r.XP, &r.Tier); err != nil {
			return nil, fmt.Errorf("scan league range top: %w", err)
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("league range top rows: %w", err)
	}
	return out, nil
}

func leagueRangeLeft(db *sql.DB, rebind func(string) string, from, to string) (int, error) {
	rows, err := db.Query(rebind(`SELECT week, student_id FROM league_members WHERE week >= ? AND week <= ? ORDER BY student_id, week`), from, to)
	if err != nil {
		return 0, fmt.Errorf("league range left: %w", err)
	}
	defer rows.Close()
	last := map[string]string{}
	left := 0
	for rows.Next() {
		var week, id string
		if err := rows.Scan(&week, &id); err != nil {
			return 0, fmt.Errorf("scan league range left: %w", err)
		}
		if prev, ok := last[id]; ok && prev < to && week != shiftWeekKey(prev, 7) {
			left++
		}
		last[id] = week
	}
	if err := rows.Err(); err != nil {
		return 0, fmt.Errorf("league range left rows: %w", err)
	}
	for _, week := range last {
		if week < to {
			left++
		}
	}
	return left, nil
}
