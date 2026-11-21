package storage

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// CountryDay is one UTC day in the activity series behind the admin map.
type CountryDay struct {
	Day    string
	Active int
}

// StudentLocations is the aggregate the admin map renders.
//
// Every map is keyed by an ISO 3166-1 alpha-2 code. An account that has never
// been seen through the Cloudflare edge has no country: it counts in Total but
// not in Known. The three windows are deliberately separate facts. Known says
// where accounts are, the active windows say where people are practising, and
// the new windows say where the newest accounts came from.
type StudentLocations struct {
	Known       map[string]int
	Active7     map[string]int
	ActivePrev7 map[string]int
	New7        map[string]int
	NewPrev7    map[string]int
	Daily       []CountryDay
	Located     int
	Total       int
	// The country maps count a learner once per country. Someone who appears
	// under two countries in one window (travel, a VPN switch) must not be
	// counted twice in the headline, so the window totals are queried
	// separately rather than summed from the maps.
	Active7Total     int
	ActivePrev7Total int
	// Visits count learner-days, not people: a learner who practises on three
	// days is three visits. They answer "how much did people come back" while
	// Active7Total answers "how many people", and the two are deliberately not
	// the same number.
	Visits7     int
	VisitsPrev7 int
	// Device maps count distinct learners per coarse class in each window. The
	// class is derived from the user agent and the agent itself is never kept.
	Device7     map[string]int
	DevicePrev7 map[string]int
}

// LocationWindow names the day boundaries of one admin locations read.
//
// RecentFrom is the first day of the current seven-day window. PriorFrom is the
// first day of the previous window, which ends the day before RecentFrom.
// DailyFrom and Today bound the thirty-day series. Days are YYYY-MM-DD in UTC.
type LocationWindow struct {
	RecentFrom string
	PriorFrom  string
	DailyFrom  string
	Today      string
}

// StudentActivity is the last coarse activity recorded for one account. It is
// read in a batch for the admin user list so that listing two hundred accounts
// does not run two hundred extra queries.
type StudentActivity struct {
	Country    string
	Device     string
	LastSeenAt string
}

type locationQuerier interface {
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

// RecordStudentActivity stores the coarse country and device class the edge
// reported and notes that the learner was active today.
//
// The students row keeps only the latest country, device and the moment the
// learner was last seen. The daily row is what makes growth measurable: it is
// one row per learner per UTC day, so last week's activity is still answerable
// after this week's capture overwrites the student row. An empty country or
// device leaves the stored value alone rather than erasing it, because a
// missing edge header is not evidence that the learner moved. No IP address or
// user agent is stored in either place.
func (s *SQLiteStore) RecordStudentActivity(studentID, country, device, day string, at time.Time) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("record activity: begin: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.Exec(
		`UPDATE students SET
			country = COALESCE(NULLIF(?, ''), country),
			country_at = CASE WHEN ? = '' THEN country_at ELSE ? END,
			device = COALESCE(NULLIF(?, ''), device),
			last_seen_at = ?
		 WHERE id = ?`,
		country, country, at.UTC().Format(time.RFC3339), device, at.UTC().Format(time.RFC3339), studentID,
	); err != nil {
		return fmt.Errorf("record activity: student: %w", err)
	}
	if _, err := tx.Exec(
		`INSERT INTO learner_country_days (student_id, day, country, device) VALUES (?, ?, ?, ?)
		 ON CONFLICT(student_id, day) DO UPDATE SET country = excluded.country, device = excluded.device`,
		studentID, day, country, device,
	); err != nil {
		return fmt.Errorf("record activity: day: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("record activity: commit: %w", err)
	}
	return nil
}

// RecordStudentActivity is the Postgres twin of the SQLite method above.
func (s *PostgresStore) RecordStudentActivity(studentID, country, device, day string, at time.Time) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("record activity: begin: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.Exec(
		`UPDATE students SET
			country = COALESCE(NULLIF($1, ''), country),
			country_at = CASE WHEN $2 = '' THEN country_at ELSE $3 END,
			device = COALESCE(NULLIF($4, ''), device),
			last_seen_at = $5
		 WHERE id = $6`,
		country, country, at.UTC().Format(time.RFC3339), device, at.UTC().Format(time.RFC3339), studentID,
	); err != nil {
		return fmt.Errorf("record activity: student: %w", err)
	}
	if _, err := tx.Exec(
		`INSERT INTO learner_country_days (student_id, day, country, device) VALUES ($1, $2, $3, $4)
		 ON CONFLICT (student_id, day) DO UPDATE SET country = excluded.country, device = excluded.device`,
		studentID, day, country, device,
	); err != nil {
		return fmt.Errorf("record activity: day: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("record activity: commit: %w", err)
	}
	return nil
}

// StudentActivities returns the stored country, device class and last-seen
// stamp for the given account ids. An id with no row is absent from the map.
func (s *SQLiteStore) StudentActivities(ids []string) (map[string]StudentActivity, error) {
	return queryStudentActivities(s.db, func(q string) string { return q }, ids)
}

// StudentActivities is the Postgres twin of the SQLite method above.
func (s *PostgresStore) StudentActivities(ids []string) (map[string]StudentActivity, error) {
	return queryStudentActivities(s.db, rebindDollars, ids)
}

func queryStudentActivities(q locationQuerier, rebind func(string) string, ids []string) (map[string]StudentActivity, error) {
	out := make(map[string]StudentActivity, len(ids))
	if len(ids) == 0 {
		return out, nil
	}
	placeholders := strings.TrimSuffix(strings.Repeat("?,", len(ids)), ",")
	args := make([]any, len(ids))
	for i, id := range ids {
		args[i] = id
	}
	rows, err := q.Query(rebind(
		`SELECT id, country, device, last_seen_at FROM students WHERE id IN (`+placeholders+`)`),
		args...)
	if err != nil {
		return nil, fmt.Errorf("student activities: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		var a StudentActivity
		if err := rows.Scan(&id, &a.Country, &a.Device, &a.LastSeenAt); err != nil {
			return nil, fmt.Errorf("student activities scan: %w", err)
		}
		out[id] = a
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("student activities rows: %w", err)
	}
	return out, nil
}

// StudentLocations returns the country aggregate for one admin read.
func (s *SQLiteStore) StudentLocations(window LocationWindow) (*StudentLocations, error) {
	return queryStudentLocations(s.db, func(q string) string { return q }, window)
}

// StudentLocations is the Postgres twin of the SQLite method above.
func (s *PostgresStore) StudentLocations(window LocationWindow) (*StudentLocations, error) {
	return queryStudentLocations(s.db, rebindDollars, window)
}

func queryStudentLocations(q locationQuerier, rebind func(string) string, w LocationWindow) (*StudentLocations, error) {
	out := &StudentLocations{
		Known:       map[string]int{},
		Active7:     map[string]int{},
		ActivePrev7: map[string]int{},
		New7:        map[string]int{},
		NewPrev7:    map[string]int{},
		Device7:     map[string]int{},
		DevicePrev7: map[string]int{},
	}

	if err := queryCounts(q, rebind,
		"SELECT country, COUNT(*) FROM students WHERE country <> '' GROUP BY country",
		nil, out.Known); err != nil {
		return nil, err
	}
	for _, n := range out.Known {
		out.Located += n
	}

	if err := queryCounts(q, rebind,
		`SELECT country, COUNT(DISTINCT student_id) FROM learner_country_days
		 WHERE day >= ? GROUP BY country`,
		[]any{w.RecentFrom}, out.Active7); err != nil {
		return nil, err
	}
	if err := queryCounts(q, rebind,
		`SELECT country, COUNT(DISTINCT student_id) FROM learner_country_days
		 WHERE day >= ? AND day < ? GROUP BY country`,
		[]any{w.PriorFrom, w.RecentFrom}, out.ActivePrev7); err != nil {
		return nil, err
	}
	if err := q.QueryRow(rebind(
		`SELECT COUNT(DISTINCT student_id) FROM learner_country_days WHERE day >= ?`),
		w.RecentFrom).Scan(&out.Active7Total); err != nil {
		return nil, fmt.Errorf("locations: active total: %w", err)
	}
	if err := q.QueryRow(rebind(
		`SELECT COUNT(DISTINCT student_id) FROM learner_country_days WHERE day >= ? AND day < ?`),
		w.PriorFrom, w.RecentFrom).Scan(&out.ActivePrev7Total); err != nil {
		return nil, fmt.Errorf("locations: prior active total: %w", err)
	}

	// Visits are learner-days, so one row per learner per day is one visit.
	if err := q.QueryRow(rebind(
		`SELECT COUNT(*) FROM learner_country_days WHERE day >= ?`),
		w.RecentFrom).Scan(&out.Visits7); err != nil {
		return nil, fmt.Errorf("locations: visits: %w", err)
	}
	if err := q.QueryRow(rebind(
		`SELECT COUNT(*) FROM learner_country_days WHERE day >= ? AND day < ?`),
		w.PriorFrom, w.RecentFrom).Scan(&out.VisitsPrev7); err != nil {
		return nil, fmt.Errorf("locations: prior visits: %w", err)
	}
	if err := queryCounts(q, rebind,
		`SELECT device, COUNT(DISTINCT student_id) FROM learner_country_days
		 WHERE day >= ? AND device <> '' GROUP BY device`,
		[]any{w.RecentFrom}, out.Device7); err != nil {
		return nil, err
	}
	if err := queryCounts(q, rebind,
		`SELECT device, COUNT(DISTINCT student_id) FROM learner_country_days
		 WHERE day >= ? AND day < ? AND device <> '' GROUP BY device`,
		[]any{w.PriorFrom, w.RecentFrom}, out.DevicePrev7); err != nil {
		return nil, err
	}

	// created_at is a datetime string in both stores; its first ten characters
	// are the UTC day, so a string comparison is the day comparison.
	if err := queryCounts(q, rebind,
		`SELECT country, COUNT(*) FROM students
		 WHERE country <> '' AND substr(created_at, 1, 10) >= ? GROUP BY country`,
		[]any{w.RecentFrom}, out.New7); err != nil {
		return nil, err
	}
	if err := queryCounts(q, rebind,
		`SELECT country, COUNT(*) FROM students
		 WHERE country <> '' AND substr(created_at, 1, 10) >= ? AND substr(created_at, 1, 10) < ?
		 GROUP BY country`,
		[]any{w.PriorFrom, w.RecentFrom}, out.NewPrev7); err != nil {
		return nil, err
	}

	rows, err := q.Query(rebind(
		`SELECT day, COUNT(DISTINCT student_id) FROM learner_country_days
		 WHERE day >= ? AND day <= ? GROUP BY day ORDER BY day`),
		w.DailyFrom, w.Today)
	if err != nil {
		return nil, fmt.Errorf("locations: daily: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var d CountryDay
		if err := rows.Scan(&d.Day, &d.Active); err != nil {
			return nil, fmt.Errorf("locations: daily scan: %w", err)
		}
		out.Daily = append(out.Daily, d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("locations: daily rows: %w", err)
	}

	if err := q.QueryRow("SELECT COUNT(*) FROM students").Scan(&out.Total); err != nil {
		return nil, fmt.Errorf("locations: total: %w", err)
	}
	return out, nil
}

func queryCounts(q locationQuerier, rebind func(string) string, query string, args []any, into map[string]int) error {
	rows, err := q.Query(rebind(query), args...)
	if err != nil {
		return fmt.Errorf("locations: counts: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var code string
		var n int
		if err := rows.Scan(&code, &n); err != nil {
			return fmt.Errorf("locations: count scan: %w", err)
		}
		into[code] = n
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("locations: count rows: %w", err)
	}
	return nil
}
