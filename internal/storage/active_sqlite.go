package storage

import (
	"fmt"
	"time"
)

// ActiveLearnerCounts reads the two live numbers behind the admin overview.
//
// The window is passed in rather than computed here so the definition belongs
// to the caller and can be tested at a fixed instant. `datetime()` on both
// sides normalises the RFC3339 argument against the stored timestamp format,
// which is what the attempts window already does.
func (s *SQLiteStore) ActiveLearnerCounts(since, dayStart time.Time) (int, int, int, int, error) {
	var activeNow, activeNowGuests int
	if err := s.db.QueryRow(
		`SELECT
		   COUNT(DISTINCT CASE WHEN substr(student_id, 1, 6) <> 'guest_' THEN student_id END),
		   COUNT(DISTINCT CASE WHEN substr(student_id, 1, 6) = 'guest_' THEN student_id END)
		 FROM active_sessions
		 WHERE datetime(updated_at) >= datetime(?)`,
		since.UTC().Format(time.RFC3339)).Scan(&activeNow, &activeNowGuests); err != nil {
		return 0, 0, 0, 0, fmt.Errorf("count active learners: %w", err)
	}

	var activeToday, activeTodayGuests int
	if err := s.db.QueryRow(
		`SELECT
		   COUNT(DISTINCT CASE WHEN substr(student_id, 1, 6) <> 'guest_' THEN student_id END),
		   COUNT(DISTINCT CASE WHEN substr(student_id, 1, 6) = 'guest_' THEN student_id END)
		 FROM attempts
		 WHERE datetime(timestamp) >= datetime(?)`,
		dayStart.UTC().Format(time.RFC3339)).Scan(&activeToday, &activeTodayGuests); err != nil {
		return 0, 0, 0, 0, fmt.Errorf("count learners today: %w", err)
	}
	return activeNow, activeNowGuests, activeToday, activeTodayGuests, nil
}
