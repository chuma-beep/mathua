package storage

import (
	"fmt"
	"time"
)

// ActiveLearnerCounts mirrors the SQLite query. `updated_at` and `timestamp`
// are stored as text, so both sides are cast to timestamptz before comparing;
// the argument is a time.Time rather than a formatted string.
func (s *PostgresStore) ActiveLearnerCounts(since, dayStart time.Time) (int, int, int, int, error) {
	var activeNow, activeNowGuests int
	if err := s.db.QueryRow(
		`SELECT
		   COUNT(DISTINCT CASE WHEN substr(student_id, 1, 6) <> 'guest_' THEN student_id END),
		   COUNT(DISTINCT CASE WHEN substr(student_id, 1, 6) = 'guest_' THEN student_id END)
		 FROM active_sessions
		 WHERE updated_at::timestamptz >= $1`,
		since.UTC()).Scan(&activeNow, &activeNowGuests); err != nil {
		return 0, 0, 0, 0, fmt.Errorf("count active learners: %w", err)
	}

	var activeToday, activeTodayGuests int
	if err := s.db.QueryRow(
		`SELECT
		   COUNT(DISTINCT CASE WHEN substr(student_id, 1, 6) <> 'guest_' THEN student_id END),
		   COUNT(DISTINCT CASE WHEN substr(student_id, 1, 6) = 'guest_' THEN student_id END)
		 FROM attempts
		 WHERE timestamp::timestamptz >= $1`,
		dayStart.UTC()).Scan(&activeToday, &activeTodayGuests); err != nil {
		return 0, 0, 0, 0, fmt.Errorf("count learners today: %w", err)
	}
	return activeNow, activeNowGuests, activeToday, activeTodayGuests, nil
}
