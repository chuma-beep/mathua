package storage

import (
	"database/sql"
	"fmt"
	"time"
)

const notificationPreferenceColumns = `student_id, email_enabled, in_app_enabled,
	reviews_enabled, study_enabled, milestones_enabled, inactivity_enabled,
	preferred_time_min, study_days, study_time_min, timezone,
	quiet_start_min, quiet_end_min, paused_until, created_at, updated_at`

const ensureNotificationPreferencesSQL = `INSERT INTO notification_preferences (
	student_id, email_enabled, in_app_enabled, reviews_enabled, study_enabled,
	milestones_enabled, inactivity_enabled, preferred_time_min, study_days,
	study_time_min, timezone, quiet_start_min, quiet_end_min, paused_until,
	created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(student_id) DO NOTHING`

const upsertNotificationPreferencesSQL = `INSERT INTO notification_preferences (
	student_id, email_enabled, in_app_enabled, reviews_enabled, study_enabled,
	milestones_enabled, inactivity_enabled, preferred_time_min, study_days,
	study_time_min, timezone, quiet_start_min, quiet_end_min, paused_until,
	created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(student_id) DO UPDATE SET
	email_enabled = excluded.email_enabled,
	in_app_enabled = excluded.in_app_enabled,
	reviews_enabled = excluded.reviews_enabled,
	study_enabled = excluded.study_enabled,
	milestones_enabled = excluded.milestones_enabled,
	inactivity_enabled = excluded.inactivity_enabled,
	preferred_time_min = excluded.preferred_time_min,
	study_days = excluded.study_days,
	study_time_min = excluded.study_time_min,
	timezone = excluded.timezone,
	quiet_start_min = excluded.quiet_start_min,
	quiet_end_min = excluded.quiet_end_min,
	paused_until = excluded.paused_until,
	updated_at = excluded.updated_at`

// DefaultNotificationPreferences returns the Phase 1 defaults for a new
// learner: in-app reviews and milestones on, email and the schedule-based
// nudges off.
func DefaultNotificationPreferences(studentID string, now time.Time) NotificationPreferences {
	return NotificationPreferences{
		StudentID:         studentID,
		EmailEnabled:      false,
		InAppEnabled:      true,
		ReviewsEnabled:    true,
		StudyEnabled:      false,
		MilestonesEnabled: true,
		InactivityEnabled: false,
		PreferredTimeMin:  NotificationDefaultPreferredTimeMin,
		StudyDays:         "",
		StudyTimeMin:      NotificationDefaultStudyTimeMin,
		Timezone:          NotificationDefaultTimezone,
		QuietStartMin:     NotificationDefaultQuietStartMin,
		QuietEndMin:       NotificationDefaultQuietEndMin,
		CreatedAt:         now.UTC(),
		UpdatedAt:         now.UTC(),
	}
}

func notificationPreferenceArgs(p *NotificationPreferences) []interface{} {
	paused := int64(0)
	if p.PausedUntil != nil {
		paused = p.PausedUntil.UTC().Unix()
	}
	return []interface{}{
		p.StudentID, boolToInt(p.EmailEnabled), boolToInt(p.InAppEnabled),
		boolToInt(p.ReviewsEnabled), boolToInt(p.StudyEnabled),
		boolToInt(p.MilestonesEnabled), boolToInt(p.InactivityEnabled),
		p.PreferredTimeMin, p.StudyDays, p.StudyTimeMin, p.Timezone,
		p.QuietStartMin, p.QuietEndMin, paused,
		p.CreatedAt.UTC().Unix(), p.UpdatedAt.UTC().Unix(),
	}
}

type notificationRowScanner interface {
	Scan(dest ...interface{}) error
}

func scanNotificationPreferences(row notificationRowScanner) (*NotificationPreferences, error) {
	var p NotificationPreferences
	var email, inApp, reviews, study, milestones, inactivity int
	var paused int64
	var createdAt, updatedAt int64
	err := row.Scan(
		&p.StudentID, &email, &inApp, &reviews, &study, &milestones,
		&inactivity, &p.PreferredTimeMin, &p.StudyDays, &p.StudyTimeMin,
		&p.Timezone, &p.QuietStartMin, &p.QuietEndMin, &paused,
		&createdAt, &updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	p.EmailEnabled = email != 0
	p.InAppEnabled = inApp != 0
	p.ReviewsEnabled = reviews != 0
	p.StudyEnabled = study != 0
	p.MilestonesEnabled = milestones != 0
	p.InactivityEnabled = inactivity != 0
	if paused > 0 {
		t := time.Unix(paused, 0).UTC()
		p.PausedUntil = &t
	}
	p.CreatedAt = time.Unix(createdAt, 0).UTC()
	p.UpdatedAt = time.Unix(updatedAt, 0).UTC()
	return &p, nil
}

func getNotificationPreferences(db *sql.DB, query, studentID string) (*NotificationPreferences, error) {
	row := db.QueryRow(query, studentID)
	p, err := scanNotificationPreferences(row)
	if err != nil {
		return nil, fmt.Errorf("get notification preferences: %w", err)
	}
	return p, nil
}

func upsertNotificationPreferences(db *sql.DB, query string, p *NotificationPreferences) error {
	now := time.Now().UTC()
	if p.CreatedAt.IsZero() {
		p.CreatedAt = now
	}
	p.UpdatedAt = now
	if _, err := db.Exec(query, notificationPreferenceArgs(p)...); err != nil {
		return fmt.Errorf("upsert notification preferences: %w", err)
	}
	return nil
}

const queryGetNotificationPreferences = `SELECT ` + notificationPreferenceColumns + `
FROM notification_preferences WHERE student_id = ?`

func (s *SQLiteStore) GetNotificationPreferences(studentID string) (*NotificationPreferences, error) {
	return getNotificationPreferences(s.db, queryGetNotificationPreferences, studentID)
}

func (s *SQLiteStore) EnsureNotificationPreferences(studentID string, now time.Time) (*NotificationPreferences, error) {
	p := DefaultNotificationPreferences(studentID, now)
	if _, err := s.db.Exec(ensureNotificationPreferencesSQL, notificationPreferenceArgs(&p)...); err != nil {
		return nil, fmt.Errorf("ensure notification preferences: %w", err)
	}
	return s.GetNotificationPreferences(studentID)
}

func (s *SQLiteStore) UpsertNotificationPreferences(p *NotificationPreferences) error {
	return upsertNotificationPreferences(s.db, upsertNotificationPreferencesSQL, p)
}

func (s *PostgresStore) GetNotificationPreferences(studentID string) (*NotificationPreferences, error) {
	return getNotificationPreferences(s.db, rebindDollars(queryGetNotificationPreferences), studentID)
}

func (s *PostgresStore) EnsureNotificationPreferences(studentID string, now time.Time) (*NotificationPreferences, error) {
	p := DefaultNotificationPreferences(studentID, now)
	if _, err := s.db.Exec(rebindDollars(ensureNotificationPreferencesSQL), notificationPreferenceArgs(&p)...); err != nil {
		return nil, fmt.Errorf("ensure notification preferences: %w", err)
	}
	return s.GetNotificationPreferences(studentID)
}

func (s *PostgresStore) UpsertNotificationPreferences(p *NotificationPreferences) error {
	return upsertNotificationPreferences(s.db, rebindDollars(upsertNotificationPreferencesSQL), p)
}
