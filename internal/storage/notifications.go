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

const notificationColumns = `id, student_id, kind, dedupe_key, title, body, href,
	created_at, read_at, canceled_at`

const insertNotificationSQL = `INSERT INTO notifications (
	id, student_id, kind, dedupe_key, title, body, href, created_at, read_at, canceled_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(student_id, dedupe_key) DO NOTHING`

type notificationScanner interface {
	Scan(dest ...interface{}) error
}

func scanNotification(row notificationScanner) (*Notification, error) {
	var n Notification
	var createdAt, readAt, canceledAt int64
	if err := row.Scan(
		&n.ID, &n.StudentID, &n.Kind, &n.DedupeKey, &n.Title, &n.Body,
		&n.Href, &createdAt, &readAt, &canceledAt,
	); err != nil {
		return nil, err
	}
	n.CreatedAt = time.Unix(createdAt, 0).UTC()
	if readAt > 0 {
		n.ReadAt = time.Unix(readAt, 0).UTC()
	}
	if canceledAt > 0 {
		n.CanceledAt = time.Unix(canceledAt, 0).UTC()
	}
	return &n, nil
}

func createNotification(db *sql.DB, query string, n *Notification) (bool, error) {
	if n == nil {
		return false, fmt.Errorf("create notification: nil notification")
	}
	if n.ID == "" || n.StudentID == "" || n.Title == "" {
		return false, fmt.Errorf("create notification: id, student_id and title are required")
	}
	created := n.CreatedAt.UTC()
	if created.IsZero() {
		created = time.Now().UTC()
	}
	var readAt, canceledAt int64
	if !n.ReadAt.IsZero() {
		readAt = n.ReadAt.UTC().Unix()
	}
	if !n.CanceledAt.IsZero() {
		canceledAt = n.CanceledAt.UTC().Unix()
	}
	res, err := db.Exec(query, n.ID, n.StudentID, n.Kind, n.DedupeKey, n.Title,
		n.Body, n.Href, created.Unix(), readAt, canceledAt)
	if err != nil {
		return false, fmt.Errorf("create notification: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("create notification rows: %w", err)
	}
	return rows > 0, nil
}

func listNotifications(db *sql.DB, studentID string, opts NotificationListOptions, rebind func(string) string) (NotificationPage, error) {
	limit := opts.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 50 {
		limit = 50
	}
	query := `SELECT ` + notificationColumns + `
FROM notifications
WHERE student_id = ? AND canceled_at = 0`
	args := []interface{}{studentID}
	if opts.UnreadOnly {
		query += ` AND read_at = 0`
	}
	if opts.BeforeCreatedAt > 0 {
		query += ` AND (created_at < ? OR (created_at = ? AND id < ?))`
		args = append(args, opts.BeforeCreatedAt, opts.BeforeCreatedAt, opts.BeforeID)
	}
	query += ` ORDER BY created_at DESC, id DESC LIMIT ?`
	args = append(args, limit+1)

	rows, err := db.Query(rebind(query), args...)
	if err != nil {
		return NotificationPage{}, fmt.Errorf("list notifications: %w", err)
	}
	defer rows.Close()
	out := make([]Notification, 0, limit)
	for rows.Next() {
		n, err := scanNotification(rows)
		if err != nil {
			return NotificationPage{}, fmt.Errorf("scan notification: %w", err)
		}
		out = append(out, *n)
	}
	if err := rows.Err(); err != nil {
		return NotificationPage{}, fmt.Errorf("notification rows: %w", err)
	}
	hasMore := len(out) > limit
	if hasMore {
		out = out[:limit]
	}
	return NotificationPage{Notifications: out, HasMore: hasMore}, nil
}

func notificationUnreadCount(db *sql.DB, query, studentID string) (int, error) {
	var count int
	if err := db.QueryRow(query, studentID).Scan(&count); err != nil {
		return 0, fmt.Errorf("notification unread count: %w", err)
	}
	return count, nil
}

func markNotificationRead(db *sql.DB, studentID, notificationID string, now time.Time, rebind func(string) string) error {
	readAt := now.UTC().Unix()
	res, err := db.Exec(rebind(`UPDATE notifications SET read_at = ?
		WHERE id = ? AND student_id = ? AND canceled_at = 0 AND read_at = 0`),
		readAt, notificationID, studentID)
	if err != nil {
		return fmt.Errorf("mark notification read: %w", err)
	}
	changed, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("mark notification read rows: %w", err)
	}
	if changed > 0 {
		return nil
	}
	// Already read and still owned is success (idempotent). Missing, cancelled
	// or owned by somebody else is deliberately the same not-found error.
	var one int
	err = db.QueryRow(rebind(`SELECT 1 FROM notifications
		WHERE id = ? AND student_id = ? AND canceled_at = 0`), notificationID, studentID).Scan(&one)
	if err == sql.ErrNoRows {
		return ErrNotificationNotFound
	}
	if err != nil {
		return fmt.Errorf("mark notification read ownership: %w", err)
	}
	return nil
}

func markAllNotificationsRead(db *sql.DB, studentID string, now time.Time, rebind func(string) string) error {
	if _, err := db.Exec(rebind(`UPDATE notifications SET read_at = ?
		WHERE student_id = ? AND canceled_at = 0 AND read_at = 0`), now.UTC().Unix(), studentID); err != nil {
		return fmt.Errorf("mark all notifications read: %w", err)
	}
	return nil
}

const queryNotificationUnreadCount = `SELECT COUNT(*) FROM notifications
WHERE student_id = ? AND canceled_at = 0 AND read_at = 0`

func (s *SQLiteStore) CreateNotification(n *Notification) (bool, error) {
	return createNotification(s.db, insertNotificationSQL, n)
}

func (s *SQLiteStore) ListNotifications(studentID string, opts NotificationListOptions) (NotificationPage, error) {
	return listNotifications(s.db, studentID, opts, func(q string) string { return q })
}

func (s *SQLiteStore) NotificationUnreadCount(studentID string) (int, error) {
	return notificationUnreadCount(s.db, queryNotificationUnreadCount, studentID)
}

func (s *SQLiteStore) MarkNotificationRead(studentID, notificationID string, now time.Time) error {
	return markNotificationRead(s.db, studentID, notificationID, now, func(q string) string { return q })
}

func (s *SQLiteStore) MarkAllNotificationsRead(studentID string, now time.Time) error {
	return markAllNotificationsRead(s.db, studentID, now, func(q string) string { return q })
}

func (s *PostgresStore) CreateNotification(n *Notification) (bool, error) {
	return createNotification(s.db, rebindDollars(insertNotificationSQL), n)
}

func (s *PostgresStore) ListNotifications(studentID string, opts NotificationListOptions) (NotificationPage, error) {
	return listNotifications(s.db, studentID, opts, rebindDollars)
}

func (s *PostgresStore) NotificationUnreadCount(studentID string) (int, error) {
	return notificationUnreadCount(s.db, rebindDollars(queryNotificationUnreadCount), studentID)
}

func (s *PostgresStore) MarkNotificationRead(studentID, notificationID string, now time.Time) error {
	return markNotificationRead(s.db, studentID, notificationID, now, rebindDollars)
}

func (s *PostgresStore) MarkAllNotificationsRead(studentID string, now time.Time) error {
	return markAllNotificationsRead(s.db, studentID, now, rebindDollars)
}
