package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

const notificationPreferenceColumns = `student_id, email_enabled, push_enabled,
	digest, in_app_enabled, reviews_enabled, study_enabled, milestones_enabled,
	inactivity_enabled, preferred_time_min, study_days, study_time_min, timezone,
	quiet_start_min, quiet_end_min, paused_until, created_at, updated_at`

const ensureNotificationPreferencesSQL = `INSERT INTO notification_preferences (
	student_id, email_enabled, push_enabled, digest, in_app_enabled, reviews_enabled,
	study_enabled, milestones_enabled, inactivity_enabled, preferred_time_min,
	study_days, study_time_min, timezone, quiet_start_min, quiet_end_min,
	paused_until, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(student_id) DO NOTHING`

const upsertNotificationPreferencesSQL = `INSERT INTO notification_preferences (
	student_id, email_enabled, push_enabled, digest, in_app_enabled, reviews_enabled,
	study_enabled, milestones_enabled, inactivity_enabled, preferred_time_min,
	study_days, study_time_min, timezone, quiet_start_min, quiet_end_min,
	paused_until, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(student_id) DO UPDATE SET
	email_enabled = excluded.email_enabled,
	push_enabled = excluded.push_enabled,
	digest = excluded.digest,
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
		PushEnabled:       false,
		Digest:            NotificationDigestImmediate,
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
		p.StudentID, boolToInt(p.EmailEnabled), boolToInt(p.PushEnabled),
		defaultDigest(p.Digest), boolToInt(p.InAppEnabled),
		boolToInt(p.ReviewsEnabled), boolToInt(p.StudyEnabled),
		boolToInt(p.MilestonesEnabled), boolToInt(p.InactivityEnabled),
		p.PreferredTimeMin, p.StudyDays, p.StudyTimeMin, p.Timezone,
		p.QuietStartMin, p.QuietEndMin, paused,
		p.CreatedAt.UTC().Unix(), p.UpdatedAt.UTC().Unix(),
	}
}

// defaultDigest keeps a row written before the digest column existed on the
// immediate path rather than batching it silently.
func defaultDigest(digest string) string {
	switch digest {
	case NotificationDigestDaily, NotificationDigestWeekly:
		return digest
	default:
		return NotificationDigestImmediate
	}
}

type notificationRowScanner interface {
	Scan(dest ...interface{}) error
}

func scanNotificationPreferences(row notificationRowScanner) (*NotificationPreferences, error) {
	var p NotificationPreferences
	var email, push, inApp, reviews, study, milestones, inactivity int
	var paused int64
	var createdAt, updatedAt int64
	err := row.Scan(
		&p.StudentID, &email, &push, &p.Digest, &inApp, &reviews, &study,
		&milestones, &inactivity, &p.PreferredTimeMin, &p.StudyDays,
		&p.StudyTimeMin,
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
	p.PushEnabled = push != 0
	p.Digest = defaultDigest(p.Digest)
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

// notificationRecipientColumns is the preference projection plus the student
// fields delivery needs. The preference join is intentional: a learner with no
// notification preferences row has not opted into any channel.
const notificationRecipientColumns = `p.student_id, p.email_enabled, p.push_enabled,
	p.digest, p.in_app_enabled, p.reviews_enabled, p.study_enabled, p.milestones_enabled,
	p.inactivity_enabled, p.preferred_time_min, p.study_days, p.study_time_min,
	p.timezone, p.quiet_start_min, p.quiet_end_min, p.paused_until,
	p.created_at, p.updated_at,
	s.name, s.email, s.email_verified, s.xp_total, s.xp_date, s.created_at,
	(SELECT COUNT(*) FROM concept_progress cp
		WHERE cp.student_id = s.id AND cp.status = 'MASTERED')`

func scanNotificationRecipient(row notificationRowScanner) (*NotificationRecipient, error) {
	var r NotificationRecipient
	var email, push, inApp, reviews, study, milestones, inactivity, emailVerified int
	var paused, createdAt, updatedAt int64
	var accountCreated interface{}
	var xpDate sql.NullString
	err := row.Scan(
		&r.Preferences.StudentID, &email, &push, &r.Preferences.Digest,
		&inApp, &reviews, &study, &milestones, &inactivity,
		&r.Preferences.PreferredTimeMin,
		&r.Preferences.StudyDays, &r.Preferences.StudyTimeMin,
		&r.Preferences.Timezone, &r.Preferences.QuietStartMin,
		&r.Preferences.QuietEndMin, &paused, &createdAt, &updatedAt,
		&r.Name, &r.Email, &emailVerified, &r.XPTotal, &xpDate,
		&accountCreated, &r.MasteredCount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	r.StudentID = r.Preferences.StudentID
	r.XPDate = xpDate.String
	r.Preferences.EmailEnabled = email != 0
	r.Preferences.PushEnabled = push != 0
	r.Preferences.Digest = defaultDigest(r.Preferences.Digest)
	r.Preferences.InAppEnabled = inApp != 0
	r.Preferences.ReviewsEnabled = reviews != 0
	r.Preferences.StudyEnabled = study != 0
	r.Preferences.MilestonesEnabled = milestones != 0
	r.Preferences.InactivityEnabled = inactivity != 0
	if paused > 0 {
		t := time.Unix(paused, 0).UTC()
		r.Preferences.PausedUntil = &t
	}
	r.Preferences.CreatedAt = time.Unix(createdAt, 0).UTC()
	r.Preferences.UpdatedAt = time.Unix(updatedAt, 0).UTC()
	r.EmailVerified = emailVerified != 0
	r.AccountCreated = parseStudentTimestamp(accountCreated)
	return &r, nil
}

// parseStudentTimestamp reads the account creation column, which SQLite writes
// as a datetime string and Postgres returns as a timestamp or an integer,
// depending on the driver and the column age. The zero value means the shape
// was not recognised, and callers treat that as "old enough to nudge".
func parseStudentTimestamp(v interface{}) time.Time {
	switch t := v.(type) {
	case int64:
		if t > 0 {
			return time.Unix(t, 0).UTC()
		}
	case float64:
		if t > 0 {
			return time.Unix(int64(t), 0).UTC()
		}
	case time.Time:
		return t.UTC()
	case []byte:
		return parseStudentTimestampText(string(t))
	case string:
		return parseStudentTimestampText(t)
	}
	return time.Time{}
}

// parseStudentTimestampText accepts the two text shapes the stores produce:
// SQLite's datetime('now') and an RFC 3339 stamp.
func parseStudentTimestampText(raw string) time.Time {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return time.Time{}
	}
	for _, layout := range []string{
		"2006-01-02 15:04:05.999999999-07:00",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02",
	} {
		if t, err := time.Parse(layout, raw); err == nil {
			return t.UTC()
		}
	}
	return time.Time{}
}

func listNotificationRecipients(db *sql.DB, afterStudentID string, limit int, rebind func(string) string) ([]NotificationRecipient, error) {
	if limit <= 0 {
		limit = 100
	}
	if limit > 500 {
		limit = 500
	}
	rows, err := db.Query(rebind(`SELECT `+notificationRecipientColumns+`
FROM notification_preferences p
JOIN students s ON s.id = p.student_id
WHERE p.student_id > ?
ORDER BY p.student_id
LIMIT ?`), afterStudentID, limit)
	if err != nil {
		return nil, fmt.Errorf("list notification recipients: %w", err)
	}
	defer rows.Close()
	out := make([]NotificationRecipient, 0, limit)
	for rows.Next() {
		r, err := scanNotificationRecipient(rows)
		if err != nil {
			return nil, fmt.Errorf("scan notification recipient: %w", err)
		}
		if r != nil {
			out = append(out, *r)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("notification recipient rows: %w", err)
	}
	return out, nil
}

const notificationDeliveryColumns = `notification_id, channel, status, attempts,
	next_attempt_at, claimed_until, provider_id, last_error, suppressed_reason,
	sent_at, opened_at, clicked_at`

const insertNotificationDeliverySQL = `INSERT INTO notification_deliveries (
	notification_id, channel, status, attempts, next_attempt_at, claimed_until,
	provider_id, last_error, suppressed_reason, sent_at, opened_at, clicked_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(notification_id, channel) DO NOTHING`

func createNotificationDelivery(db *sql.DB, query string, d *NotificationDelivery) (bool, error) {
	if d == nil {
		return false, fmt.Errorf("create notification delivery: nil delivery")
	}
	if d.NotificationID == "" || d.Channel == "" {
		return false, fmt.Errorf("create notification delivery: notification_id and channel are required")
	}
	status := d.Status
	if status == "" {
		status = NotificationDeliveryPending
	}
	nextAttempt := int64(0)
	if !d.NextAttemptAt.IsZero() {
		nextAttempt = d.NextAttemptAt.UTC().Unix()
	} else if status == NotificationDeliveryPending || status == NotificationDeliveryRetry {
		nextAttempt = time.Now().UTC().Unix()
	}
	res, err := db.Exec(query, d.NotificationID, d.Channel, status, d.Attempts,
		nextAttempt, notificationUnixOrZero(d.ClaimedUntil), d.ProviderID, d.LastError,
		d.SuppressedReason, notificationUnixOrZero(d.SentAt), notificationUnixOrZero(d.OpenedAt), notificationUnixOrZero(d.ClickedAt))
	if err != nil {
		return false, fmt.Errorf("create notification delivery: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("create notification delivery rows: %w", err)
	}
	return rows > 0, nil
}

func notificationUnixOrZero(t time.Time) int64 {
	if t.IsZero() {
		return 0
	}
	return t.UTC().Unix()
}

func claimNotificationDeliveriesPostgres(db *sql.DB, now time.Time, lease time.Duration, limit int) ([]NotificationDeliveryJob, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	nowUnix := now.UTC().Unix()
	claimUntil := now.Add(lease).UTC().Unix()
	rows, err := db.Query(rebindDollars(`WITH due AS (
	SELECT d.notification_id, d.channel
	FROM notification_deliveries d
	WHERE d.channel IN ('email', 'push')
	  AND d.status IN ('pending', 'retry', 'sending')
	  AND d.next_attempt_at <= $1
	  AND (d.claimed_until = 0 OR d.claimed_until <= $1)
	ORDER BY d.next_attempt_at, d.notification_id
	LIMIT $2
	FOR UPDATE SKIP LOCKED
)
UPDATE notification_deliveries d
SET status = 'sending', attempts = d.attempts + 1, claimed_until = $3
FROM due, notifications n, students s
WHERE d.notification_id = due.notification_id
  AND d.channel = due.channel
  AND n.id = d.notification_id
  AND s.id = n.student_id
RETURNING d.notification_id, n.student_id, d.channel, d.attempts,
	n.kind, n.title, n.body, n.href, s.name, s.email, s.email_verified`), nowUnix, limit, claimUntil)
	if err != nil {
		return nil, fmt.Errorf("claim notification deliveries: %w", err)
	}
	defer rows.Close()
	out := make([]NotificationDeliveryJob, 0, limit)
	for rows.Next() {
		var job NotificationDeliveryJob
		var emailVerified int
		if err := rows.Scan(&job.NotificationID, &job.StudentID, &job.Channel, &job.Attempts,
			&job.Kind, &job.Title, &job.Body, &job.Href, &job.StudentName, &job.Email, &emailVerified); err != nil {
			return nil, fmt.Errorf("scan claimed notification delivery: %w", err)
		}
		job.EmailVerified = emailVerified != 0
		out = append(out, job)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("claimed notification delivery rows: %w", err)
	}
	return out, nil
}

func claimNotificationDeliveriesSQLite(db *sql.DB, now time.Time, lease time.Duration, limit int) ([]NotificationDeliveryJob, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	nowUnix := now.UTC().Unix()
	claimUntil := now.Add(lease).UTC().Unix()
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("claim notification deliveries: begin: %w", err)
	}
	defer tx.Rollback()

	rows, err := tx.Query(`UPDATE notification_deliveries
SET status = 'sending', attempts = attempts + 1, claimed_until = ?
WHERE rowid IN (
	SELECT rowid FROM notification_deliveries
	WHERE channel IN ('email', 'push')
	  AND status IN ('pending', 'retry', 'sending')
	  AND next_attempt_at <= ?
	  AND (claimed_until = 0 OR claimed_until <= ?)
	ORDER BY next_attempt_at, notification_id
	LIMIT ?
)
RETURNING notification_id, channel, attempts`, claimUntil, nowUnix, nowUnix, limit)
	if err != nil {
		return nil, fmt.Errorf("claim notification deliveries: %w", err)
	}
	type claimed struct {
		notificationID string
		channel        string
		attempts       int
	}
	var claimedRows []claimed
	for rows.Next() {
		var c claimed
		if err := rows.Scan(&c.notificationID, &c.channel, &c.attempts); err != nil {
			rows.Close()
			return nil, fmt.Errorf("scan claimed notification delivery: %w", err)
		}
		claimedRows = append(claimedRows, c)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, fmt.Errorf("claimed notification delivery rows: %w", err)
	}
	rows.Close()

	out := make([]NotificationDeliveryJob, 0, len(claimedRows))
	for _, c := range claimedRows {
		var job NotificationDeliveryJob
		var emailVerified int
		err := tx.QueryRow(`SELECT n.student_id, n.kind, n.title, n.body, n.href, s.name, s.email, s.email_verified
FROM notifications n JOIN students s ON s.id = n.student_id
WHERE n.id = ?`, c.notificationID).Scan(
			&job.StudentID, &job.Kind, &job.Title, &job.Body, &job.Href, &job.StudentName, &job.Email, &emailVerified)
		if err == sql.ErrNoRows {
			// The notification was deleted between claim and read. Leave the
			// row to the next claim, which will also see the missing parent.
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("load claimed notification: %w", err)
		}
		job.NotificationID = c.notificationID
		job.Channel = c.channel
		job.Attempts = c.attempts
		job.EmailVerified = emailVerified != 0
		out = append(out, job)
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("claim notification deliveries: commit: %w", err)
	}
	return out, nil
}

func completeNotificationDelivery(db *sql.DB, notificationID, channel, providerID string, now time.Time, rebind func(string) string) error {
	_, err := db.Exec(rebind(`UPDATE notification_deliveries
SET status = 'sent', provider_id = ?, sent_at = ?, claimed_until = 0,
	next_attempt_at = 0, last_error = ''
WHERE notification_id = ? AND channel = ?`),
		providerID, now.UTC().Unix(), notificationID, channel)
	if err != nil {
		return fmt.Errorf("complete notification delivery: %w", err)
	}
	return nil
}

func retryNotificationDelivery(db *sql.DB, notificationID, channel, lastError string, nextAttempt time.Time, rebind func(string) string) error {
	_, err := db.Exec(rebind(`UPDATE notification_deliveries
SET status = 'retry', last_error = ?, next_attempt_at = ?, claimed_until = 0
WHERE notification_id = ? AND channel = ?`),
		lastError, nextAttempt.UTC().Unix(), notificationID, channel)
	if err != nil {
		return fmt.Errorf("retry notification delivery: %w", err)
	}
	return nil
}

func suppressNotificationDelivery(db *sql.DB, notificationID, channel, reason string, now time.Time, rebind func(string) string) error {
	_, err := db.Exec(rebind(`UPDATE notification_deliveries
SET status = 'suppressed', suppressed_reason = ?, claimed_until = 0, next_attempt_at = 0
WHERE notification_id = ? AND channel = ?`),
		reason, notificationID, channel)
	if err != nil {
		return fmt.Errorf("suppress notification delivery: %w", err)
	}
	return nil
}

func failNotificationDelivery(db *sql.DB, notificationID, channel, lastError string, now time.Time, rebind func(string) string) error {
	_, err := db.Exec(rebind(`UPDATE notification_deliveries
SET status = 'failed', last_error = ?, claimed_until = 0, next_attempt_at = 0
WHERE notification_id = ? AND channel = ?`),
		lastError, notificationID, channel)
	if err != nil {
		return fmt.Errorf("fail notification delivery: %w", err)
	}
	return nil
}

const queryGetNotificationScanState = `SELECT watermark FROM notification_scan_state
WHERE student_id = ? AND kind = ?`

func getNotificationScanState(db *sql.DB, query, studentID, kind string) (int64, error) {
	var watermark int64
	err := db.QueryRow(query, studentID, kind).Scan(&watermark)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("get notification scan state: %w", err)
	}
	return watermark, nil
}

const upsertNotificationScanStateSQL = `INSERT INTO notification_scan_state (
	student_id, kind, watermark, updated_at
) VALUES (?, ?, ?, ?)
ON CONFLICT(student_id, kind) DO UPDATE SET
	watermark = excluded.watermark,
	updated_at = excluded.updated_at`

func setNotificationScanState(db *sql.DB, query, studentID, kind string, watermark int64, now time.Time) error {
	if _, err := db.Exec(query, studentID, kind, watermark, now.UTC().Unix()); err != nil {
		return fmt.Errorf("set notification scan state: %w", err)
	}
	return nil
}

func (s *SQLiteStore) ListNotificationRecipients(afterStudentID string, limit int) ([]NotificationRecipient, error) {
	return listNotificationRecipients(s.db, afterStudentID, limit, func(q string) string { return q })
}

func (s *PostgresStore) ListNotificationRecipients(afterStudentID string, limit int) ([]NotificationRecipient, error) {
	return listNotificationRecipients(s.db, afterStudentID, limit, rebindDollars)
}

func (s *SQLiteStore) CreateNotificationDelivery(d *NotificationDelivery) (bool, error) {
	return createNotificationDelivery(s.db, insertNotificationDeliverySQL, d)
}

func (s *PostgresStore) CreateNotificationDelivery(d *NotificationDelivery) (bool, error) {
	return createNotificationDelivery(s.db, rebindDollars(insertNotificationDeliverySQL), d)
}

func (s *SQLiteStore) ClaimNotificationDeliveries(now time.Time, lease time.Duration, limit int) ([]NotificationDeliveryJob, error) {
	return claimNotificationDeliveriesSQLite(s.db, now, lease, limit)
}

func (s *PostgresStore) ClaimNotificationDeliveries(now time.Time, lease time.Duration, limit int) ([]NotificationDeliveryJob, error) {
	return claimNotificationDeliveriesPostgres(s.db, now, lease, limit)
}

func (s *SQLiteStore) CompleteNotificationDelivery(notificationID, channel, providerID string, now time.Time) error {
	return completeNotificationDelivery(s.db, notificationID, channel, providerID, now, func(q string) string { return q })
}

func (s *PostgresStore) CompleteNotificationDelivery(notificationID, channel, providerID string, now time.Time) error {
	return completeNotificationDelivery(s.db, notificationID, channel, providerID, now, rebindDollars)
}

func (s *SQLiteStore) RetryNotificationDelivery(notificationID, channel, lastError string, nextAttempt time.Time) error {
	return retryNotificationDelivery(s.db, notificationID, channel, lastError, nextAttempt, func(q string) string { return q })
}

func (s *PostgresStore) RetryNotificationDelivery(notificationID, channel, lastError string, nextAttempt time.Time) error {
	return retryNotificationDelivery(s.db, notificationID, channel, lastError, nextAttempt, rebindDollars)
}

func (s *SQLiteStore) SuppressNotificationDelivery(notificationID, channel, reason string, now time.Time) error {
	return suppressNotificationDelivery(s.db, notificationID, channel, reason, now, func(q string) string { return q })
}

func (s *PostgresStore) SuppressNotificationDelivery(notificationID, channel, reason string, now time.Time) error {
	return suppressNotificationDelivery(s.db, notificationID, channel, reason, now, rebindDollars)
}

func (s *SQLiteStore) FailNotificationDelivery(notificationID, channel, lastError string, now time.Time) error {
	return failNotificationDelivery(s.db, notificationID, channel, lastError, now, func(q string) string { return q })
}

func (s *PostgresStore) FailNotificationDelivery(notificationID, channel, lastError string, now time.Time) error {
	return failNotificationDelivery(s.db, notificationID, channel, lastError, now, rebindDollars)
}

func (s *SQLiteStore) GetNotificationScanState(studentID, kind string) (int64, error) {
	return getNotificationScanState(s.db, queryGetNotificationScanState, studentID, kind)
}

func (s *PostgresStore) GetNotificationScanState(studentID, kind string) (int64, error) {
	return getNotificationScanState(s.db, rebindDollars(queryGetNotificationScanState), studentID, kind)
}

func (s *SQLiteStore) SetNotificationScanState(studentID, kind string, watermark int64, now time.Time) error {
	return setNotificationScanState(s.db, upsertNotificationScanStateSQL, studentID, kind, watermark, now)
}

func (s *PostgresStore) SetNotificationScanState(studentID, kind string, watermark int64, now time.Time) error {
	return setNotificationScanState(s.db, rebindDollars(upsertNotificationScanStateSQL), studentID, kind, watermark, now)
}

// notificationPushSubscriptionColumns is the projection order the scanner
// expects. The endpoint is the identity: a browser re-subscribing with the same
// endpoint updates the keys instead of adding a second device.
const notificationPushSubscriptionColumns = `endpoint, student_id, p256dh, auth,
	user_agent, created_at, updated_at, last_success_at, failures`

const upsertNotificationPushSubscriptionSQL = `INSERT INTO notification_push_subscriptions (
	endpoint, student_id, p256dh, auth, user_agent, created_at, updated_at,
	last_success_at, failures
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(endpoint) DO UPDATE SET
	student_id = excluded.student_id,
	p256dh = excluded.p256dh,
	auth = excluded.auth,
	user_agent = excluded.user_agent,
	updated_at = excluded.updated_at`

func upsertNotificationPushSubscription(db *sql.DB, query string, sub *NotificationPushSubscription) error {
	if sub == nil {
		return fmt.Errorf("upsert push subscription: nil subscription")
	}
	endpoint := strings.TrimSpace(sub.Endpoint)
	studentID := strings.TrimSpace(sub.StudentID)
	if endpoint == "" || studentID == "" {
		return fmt.Errorf("upsert push subscription: endpoint and student_id are required")
	}
	sub.Endpoint = endpoint
	sub.StudentID = studentID
	now := time.Now().UTC()
	if sub.CreatedAt.IsZero() {
		sub.CreatedAt = now
	}
	sub.UpdatedAt = now
	if _, err := db.Exec(query, endpoint, studentID, sub.P256DH, sub.Auth,
		sub.UserAgent, sub.CreatedAt.UTC().Unix(), sub.UpdatedAt.UTC().Unix(),
		notificationUnixOrZero(sub.LastSuccess), sub.Failures); err != nil {
		return fmt.Errorf("upsert push subscription: %w", err)
	}
	return nil
}

func deleteNotificationPushSubscription(db *sql.DB, query string, args ...interface{}) error {
	if _, err := db.Exec(query, args...); err != nil {
		return fmt.Errorf("delete push subscription: %w", err)
	}
	return nil
}

func listNotificationPushSubscriptions(db *sql.DB, query, studentID string) ([]NotificationPushSubscription, error) {
	rows, err := db.Query(query, studentID)
	if err != nil {
		return nil, fmt.Errorf("list push subscriptions: %w", err)
	}
	defer rows.Close()
	out := make([]NotificationPushSubscription, 0, 4)
	for rows.Next() {
		var sub NotificationPushSubscription
		var createdAt, updatedAt, lastSuccess int64
		if err := rows.Scan(&sub.Endpoint, &sub.StudentID, &sub.P256DH, &sub.Auth,
			&sub.UserAgent, &createdAt, &updatedAt, &lastSuccess, &sub.Failures); err != nil {
			return nil, fmt.Errorf("scan push subscription: %w", err)
		}
		sub.CreatedAt = time.Unix(createdAt, 0).UTC()
		sub.UpdatedAt = time.Unix(updatedAt, 0).UTC()
		if lastSuccess > 0 {
			sub.LastSuccess = time.Unix(lastSuccess, 0).UTC()
		}
		out = append(out, sub)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("push subscription rows: %w", err)
	}
	return out, nil
}

func markPushSubscriptionSuccess(db *sql.DB, query, endpoint string, now time.Time) error {
	at := now.UTC().Unix()
	if _, err := db.Exec(query, at, at, endpoint); err != nil {
		return fmt.Errorf("mark push subscription success: %w", err)
	}
	return nil
}

func markPushSubscriptionFailure(db *sql.DB, query, endpoint string, now time.Time) error {
	if _, err := db.Exec(query, now.UTC().Unix(), endpoint); err != nil {
		return fmt.Errorf("mark push subscription failure: %w", err)
	}
	return nil
}

const queryNotificationDeliveryStats = `SELECT channel, status, COUNT(*)
FROM notification_deliveries
GROUP BY channel, status
ORDER BY channel, status`

func notificationDeliveryStats(db *sql.DB, query string) ([]NotificationDeliveryStat, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("notification delivery stats: %w", err)
	}
	defer rows.Close()
	out := make([]NotificationDeliveryStat, 0, 8)
	for rows.Next() {
		var stat NotificationDeliveryStat
		if err := rows.Scan(&stat.Channel, &stat.Status, &stat.Count); err != nil {
			return nil, fmt.Errorf("scan notification delivery stat: %w", err)
		}
		out = append(out, stat)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("notification delivery stat rows: %w", err)
	}
	return out, nil
}

// queryNotificationDeliveryProblems returns the rows an operator most needs:
// terminal failures and suppressions, newest first. Retries are excluded
// because they are still moving.
const queryNotificationDeliveryProblems = `SELECT d.notification_id, n.student_id,
	s.name, n.title, d.channel, d.status, d.attempts, d.last_error,
	d.suppressed_reason, n.created_at
FROM notification_deliveries d
JOIN notifications n ON n.id = d.notification_id
JOIN students s ON s.id = n.student_id
WHERE d.status IN ('failed', 'suppressed')
ORDER BY n.created_at DESC, d.notification_id
LIMIT ?`

func listNotificationDeliveryProblems(db *sql.DB, query string, limit int) ([]NotificationDeliveryRecord, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("list notification delivery problems: %w", err)
	}
	defer rows.Close()
	out := make([]NotificationDeliveryRecord, 0, limit)
	for rows.Next() {
		var rec NotificationDeliveryRecord
		var createdAt int64
		if err := rows.Scan(&rec.NotificationID, &rec.StudentID, &rec.StudentName,
			&rec.Title, &rec.Channel, &rec.Status, &rec.Attempts, &rec.LastError,
			&rec.SuppressedReason, &createdAt); err != nil {
			return nil, fmt.Errorf("scan notification delivery problem: %w", err)
		}
		rec.CreatedAt = time.Unix(createdAt, 0).UTC()
		out = append(out, rec)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("notification delivery problem rows: %w", err)
	}
	return out, nil
}

func (s *SQLiteStore) UpsertNotificationPushSubscription(sub *NotificationPushSubscription) error {
	return upsertNotificationPushSubscription(s.db, upsertNotificationPushSubscriptionSQL, sub)
}

func (s *PostgresStore) UpsertNotificationPushSubscription(sub *NotificationPushSubscription) error {
	return upsertNotificationPushSubscription(s.db, rebindDollars(upsertNotificationPushSubscriptionSQL), sub)
}

func (s *SQLiteStore) DeleteNotificationPushSubscription(studentID, endpoint string) error {
	return deleteNotificationPushSubscription(s.db,
		`DELETE FROM notification_push_subscriptions WHERE student_id = ? AND endpoint = ?`,
		studentID, endpoint)
}

func (s *PostgresStore) DeleteNotificationPushSubscription(studentID, endpoint string) error {
	return deleteNotificationPushSubscription(s.db, rebindDollars(
		`DELETE FROM notification_push_subscriptions WHERE student_id = ? AND endpoint = ?`),
		studentID, endpoint)
}

func (s *SQLiteStore) DeletePushSubscriptionByEndpoint(endpoint string) error {
	return deleteNotificationPushSubscription(s.db,
		`DELETE FROM notification_push_subscriptions WHERE endpoint = ?`, endpoint)
}

func (s *PostgresStore) DeletePushSubscriptionByEndpoint(endpoint string) error {
	return deleteNotificationPushSubscription(s.db, rebindDollars(
		`DELETE FROM notification_push_subscriptions WHERE endpoint = ?`), endpoint)
}

func (s *SQLiteStore) ListNotificationPushSubscriptions(studentID string) ([]NotificationPushSubscription, error) {
	return listNotificationPushSubscriptions(s.db,
		`SELECT `+notificationPushSubscriptionColumns+` FROM notification_push_subscriptions
		 WHERE student_id = ? ORDER BY created_at, endpoint`, studentID)
}

func (s *PostgresStore) ListNotificationPushSubscriptions(studentID string) ([]NotificationPushSubscription, error) {
	return listNotificationPushSubscriptions(s.db, rebindDollars(
		`SELECT `+notificationPushSubscriptionColumns+` FROM notification_push_subscriptions
		 WHERE student_id = ? ORDER BY created_at, endpoint`), studentID)
}

func (s *SQLiteStore) MarkPushSubscriptionSuccess(endpoint string, now time.Time) error {
	return markPushSubscriptionSuccess(s.db,
		`UPDATE notification_push_subscriptions SET last_success_at = ?, failures = 0,
		 updated_at = ? WHERE endpoint = ?`, endpoint, now)
}

func (s *PostgresStore) MarkPushSubscriptionSuccess(endpoint string, now time.Time) error {
	return markPushSubscriptionSuccess(s.db, rebindDollars(
		`UPDATE notification_push_subscriptions SET last_success_at = ?, failures = 0,
		 updated_at = ? WHERE endpoint = ?`), endpoint, now)
}

func (s *SQLiteStore) MarkPushSubscriptionFailure(endpoint string, now time.Time) error {
	return markPushSubscriptionFailure(s.db,
		`UPDATE notification_push_subscriptions SET failures = failures + 1,
		 updated_at = ? WHERE endpoint = ?`, endpoint, now)
}

func (s *PostgresStore) MarkPushSubscriptionFailure(endpoint string, now time.Time) error {
	return markPushSubscriptionFailure(s.db, rebindDollars(
		`UPDATE notification_push_subscriptions SET failures = failures + 1,
		 updated_at = ? WHERE endpoint = ?`), endpoint, now)
}

func (s *SQLiteStore) NotificationDeliveryStats() ([]NotificationDeliveryStat, error) {
	return notificationDeliveryStats(s.db, queryNotificationDeliveryStats)
}

func (s *PostgresStore) NotificationDeliveryStats() ([]NotificationDeliveryStat, error) {
	return notificationDeliveryStats(s.db, rebindDollars(queryNotificationDeliveryStats))
}

// queryNotificationEngagementStats answers the question the counts above
// cannot: of the notifications that reached someone, how many did they open and
// click? Only sent rows are counted, because a suppressed row never had a
// reader and would drag every rate toward zero for a configuration mistake
// rather than for disinterest.
const queryNotificationEngagementStats = `SELECT n.kind, d.channel, COUNT(*),
	COALESCE(SUM(CASE WHEN d.opened_at > 0 THEN 1 ELSE 0 END), 0),
	COALESCE(SUM(CASE WHEN d.clicked_at > 0 THEN 1 ELSE 0 END), 0)
FROM notification_deliveries d
JOIN notifications n ON n.id = d.notification_id
WHERE d.status = 'sent'
GROUP BY n.kind, d.channel
ORDER BY n.kind, d.channel`

func notificationEngagementStats(db *sql.DB, query string) ([]NotificationEngagementStat, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("notification engagement stats: %w", err)
	}
	defer rows.Close()
	out := make([]NotificationEngagementStat, 0, 8)
	for rows.Next() {
		var stat NotificationEngagementStat
		if err := rows.Scan(&stat.Kind, &stat.Channel, &stat.Sent, &stat.Opened, &stat.Clicked); err != nil {
			return nil, fmt.Errorf("scan notification engagement stat: %w", err)
		}
		out = append(out, stat)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("notification engagement stat rows: %w", err)
	}
	return out, nil
}

func (s *SQLiteStore) NotificationEngagementStats() ([]NotificationEngagementStat, error) {
	return notificationEngagementStats(s.db, queryNotificationEngagementStats)
}

func (s *PostgresStore) NotificationEngagementStats() ([]NotificationEngagementStat, error) {
	return notificationEngagementStats(s.db, rebindDollars(queryNotificationEngagementStats))
}

func (s *SQLiteStore) ListNotificationDeliveryProblems(limit int) ([]NotificationDeliveryRecord, error) {
	return listNotificationDeliveryProblems(s.db, queryNotificationDeliveryProblems, limit)
}

func (s *PostgresStore) ListNotificationDeliveryProblems(limit int) ([]NotificationDeliveryRecord, error) {
	return listNotificationDeliveryProblems(s.db, rebindDollars(queryNotificationDeliveryProblems), limit)
}

func markNotificationOpened(db *sql.DB, notificationID, studentID string, now time.Time, rebind func(string) string) error {
	return markNotificationTimestamp(db, notificationID, studentID, now, rebind,
		`UPDATE notifications SET read_at = ? WHERE id = ? AND student_id = ? AND read_at = 0`,
		`UPDATE notification_deliveries SET opened_at = ? WHERE notification_id = ? AND opened_at = 0`)
}

func markNotificationClicked(db *sql.DB, notificationID, studentID string, now time.Time, rebind func(string) string) error {
	return markNotificationTimestamp(db, notificationID, studentID, now, rebind,
		`UPDATE notifications SET read_at = CASE WHEN read_at = 0 THEN ? ELSE read_at END
		 WHERE id = ? AND student_id = ?`,
		`UPDATE notification_deliveries
		 SET clicked_at = ?, opened_at = CASE WHEN opened_at = 0 THEN ? ELSE opened_at END
		 WHERE notification_id = ?`)
}

// markNotificationTimestamp applies an engagement stamp to the notification and
// its deliveries. The ownership check runs after the update so a repeat open is
// idempotent, while an id the learner does not own is still reported missing.
func markNotificationTimestamp(db *sql.DB, notificationID, studentID string, now time.Time, rebind func(string) string, notificationSQL, deliverySQL string) error {
	if notificationID == "" || studentID == "" {
		return ErrNotificationNotFound
	}
	at := now.UTC().Unix()
	if _, err := db.Exec(rebind(notificationSQL), at, notificationID, studentID); err != nil {
		return fmt.Errorf("mark notification engagement: %w", err)
	}
	var owned int
	if err := db.QueryRow(rebind(`SELECT COUNT(*) FROM notifications WHERE id = ? AND student_id = ?`),
		notificationID, studentID).Scan(&owned); err != nil {
		return fmt.Errorf("mark notification engagement: %w", err)
	}
	if owned == 0 {
		return ErrNotificationNotFound
	}
	var err error
	if strings.Contains(deliverySQL, "clicked_at") {
		_, err = db.Exec(rebind(deliverySQL), at, at, notificationID)
	} else {
		_, err = db.Exec(rebind(deliverySQL), at, notificationID)
	}
	if err != nil {
		return fmt.Errorf("mark notification engagement: %w", err)
	}
	return nil
}

func (s *SQLiteStore) MarkNotificationOpened(studentID, notificationID string, now time.Time) error {
	return markNotificationOpened(s.db, notificationID, studentID, now, func(q string) string { return q })
}

func (s *PostgresStore) MarkNotificationOpened(studentID, notificationID string, now time.Time) error {
	return markNotificationOpened(s.db, notificationID, studentID, now, rebindDollars)
}

func (s *SQLiteStore) MarkNotificationClicked(studentID, notificationID string, now time.Time) error {
	return markNotificationClicked(s.db, notificationID, studentID, now, func(q string) string { return q })
}

func (s *PostgresStore) MarkNotificationClicked(studentID, notificationID string, now time.Time) error {
	return markNotificationClicked(s.db, notificationID, studentID, now, rebindDollars)
}

const replayNotificationDeliverySQL = `UPDATE notification_deliveries
SET status = 'pending', attempts = 0, next_attempt_at = ?, claimed_until = 0,
	last_error = '', suppressed_reason = ''
WHERE notification_id = ? AND channel = ? AND status IN ('failed', 'suppressed')`

func replayNotificationDelivery(db *sql.DB, query, notificationID, channel string, now time.Time) error {
	res, err := db.Exec(query, now.UTC().Unix(), notificationID, channel)
	if err != nil {
		return fmt.Errorf("replay notification delivery: %w", err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return fmt.Errorf("notification delivery is not replayable: %w", ErrNotificationNotFound)
	}
	return nil
}

func (s *SQLiteStore) ReplayNotificationDelivery(notificationID, channel string, now time.Time) error {
	return replayNotificationDelivery(s.db, replayNotificationDeliverySQL, notificationID, channel, now)
}

func (s *PostgresStore) ReplayNotificationDelivery(notificationID, channel string, now time.Time) error {
	return replayNotificationDelivery(s.db, rebindDollars(replayNotificationDeliverySQL), notificationID, channel, now)
}

const queryListNotificationDeliveries = `SELECT ` + notificationDeliveryColumns + `
FROM notification_deliveries WHERE notification_id = ? ORDER BY channel`

func listNotificationDeliveries(db *sql.DB, query, notificationID string) ([]NotificationDelivery, error) {
	rows, err := db.Query(query, notificationID)
	if err != nil {
		return nil, fmt.Errorf("list notification deliveries: %w", err)
	}
	defer rows.Close()
	out := make([]NotificationDelivery, 0, 4)
	for rows.Next() {
		var d NotificationDelivery
		var nextAttempt, claimedUntil, sentAt, openedAt, clickedAt int64
		if err := rows.Scan(&d.NotificationID, &d.Channel, &d.Status, &d.Attempts,
			&nextAttempt, &claimedUntil, &d.ProviderID, &d.LastError,
			&d.SuppressedReason, &sentAt, &openedAt, &clickedAt); err != nil {
			return nil, fmt.Errorf("scan notification delivery: %w", err)
		}
		d.NextAttemptAt = unixOrZeroTime(nextAttempt)
		d.ClaimedUntil = unixOrZeroTime(claimedUntil)
		d.SentAt = unixOrZeroTime(sentAt)
		d.OpenedAt = unixOrZeroTime(openedAt)
		d.ClickedAt = unixOrZeroTime(clickedAt)
		out = append(out, d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("notification delivery rows: %w", err)
	}
	return out, nil
}

func unixOrZeroTime(v int64) time.Time {
	if v <= 0 {
		return time.Time{}
	}
	return time.Unix(v, 0).UTC()
}

func (s *SQLiteStore) ListNotificationDeliveries(notificationID string) ([]NotificationDelivery, error) {
	return listNotificationDeliveries(s.db, queryListNotificationDeliveries, notificationID)
}

func (s *PostgresStore) ListNotificationDeliveries(notificationID string) ([]NotificationDelivery, error) {
	return listNotificationDeliveries(s.db, rebindDollars(queryListNotificationDeliveries), notificationID)
}

const queryListStudentNotificationDeliveries = `SELECT d.notification_id, n.student_id,
	s.name, n.title, d.channel, d.status, d.attempts, d.last_error,
	d.suppressed_reason, n.created_at
FROM notification_deliveries d
JOIN notifications n ON n.id = d.notification_id
JOIN students s ON s.id = n.student_id
WHERE n.student_id = ?
ORDER BY n.created_at DESC, d.notification_id, d.channel
LIMIT ?`

func listStudentNotificationDeliveries(db *sql.DB, query, studentID string, limit int) ([]NotificationDeliveryRecord, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	rows, err := db.Query(query, studentID, limit)
	if err != nil {
		return nil, fmt.Errorf("list student notification deliveries: %w", err)
	}
	defer rows.Close()
	out := make([]NotificationDeliveryRecord, 0, limit)
	for rows.Next() {
		var rec NotificationDeliveryRecord
		var createdAt int64
		if err := rows.Scan(&rec.NotificationID, &rec.StudentID, &rec.StudentName,
			&rec.Title, &rec.Channel, &rec.Status, &rec.Attempts, &rec.LastError,
			&rec.SuppressedReason, &createdAt); err != nil {
			return nil, fmt.Errorf("scan student notification delivery: %w", err)
		}
		rec.CreatedAt = time.Unix(createdAt, 0).UTC()
		out = append(out, rec)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("student notification delivery rows: %w", err)
	}
	return out, nil
}

func (s *SQLiteStore) ListStudentNotificationDeliveries(studentID string, limit int) ([]NotificationDeliveryRecord, error) {
	return listStudentNotificationDeliveries(s.db, queryListStudentNotificationDeliveries, studentID, limit)
}

func (s *PostgresStore) ListStudentNotificationDeliveries(studentID string, limit int) ([]NotificationDeliveryRecord, error) {
	return listStudentNotificationDeliveries(s.db, rebindDollars(queryListStudentNotificationDeliveries), studentID, limit)
}

const queryListDueNotificationDeliveries = `SELECT d.notification_id, n.student_id,
	d.channel, d.attempts, n.kind, n.title, n.body, n.href, s.name, s.email,
	s.email_verified
FROM notification_deliveries d
JOIN notifications n ON n.id = d.notification_id
JOIN students s ON s.id = n.student_id
WHERE n.student_id = ? AND d.channel = ? AND d.status IN ('pending', 'retry', 'sending')
	AND d.next_attempt_at > 0 AND d.next_attempt_at <= ?
ORDER BY d.next_attempt_at, d.notification_id
LIMIT 100`

func listDueNotificationDeliveries(db *sql.DB, query, studentID, channel string, now time.Time) ([]NotificationDeliveryJob, error) {
	rows, err := db.Query(query, studentID, channel, now.UTC().Unix())
	if err != nil {
		return nil, fmt.Errorf("list due notification deliveries: %w", err)
	}
	defer rows.Close()
	out := make([]NotificationDeliveryJob, 0, 8)
	for rows.Next() {
		var job NotificationDeliveryJob
		var verified int
		if err := rows.Scan(&job.NotificationID, &job.StudentID, &job.Channel,
			&job.Attempts, &job.Kind, &job.Title, &job.Body, &job.Href,
			&job.StudentName, &job.Email, &verified); err != nil {
			return nil, fmt.Errorf("scan due notification delivery: %w", err)
		}
		job.EmailVerified = verified != 0
		out = append(out, job)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("due notification delivery rows: %w", err)
	}
	return out, nil
}

func (s *SQLiteStore) ListDueNotificationDeliveries(studentID, channel string, now time.Time) ([]NotificationDeliveryJob, error) {
	return listDueNotificationDeliveries(s.db, queryListDueNotificationDeliveries, studentID, channel, now)
}

func (s *PostgresStore) ListDueNotificationDeliveries(studentID, channel string, now time.Time) ([]NotificationDeliveryJob, error) {
	return listDueNotificationDeliveries(s.db, rebindDollars(queryListDueNotificationDeliveries), studentID, channel, now)
}

const deletePrunedDeliveriesSQL = `DELETE FROM notification_deliveries
WHERE notification_id IN (SELECT id FROM notifications WHERE created_at < ?)`

const deletePrunedNotificationsSQL = `DELETE FROM notifications WHERE created_at < ?`

func pruneNotificationHistory(db *sql.DB, before time.Time, rebind func(string) string) (NotificationPruneResult, error) {
	var result NotificationPruneResult
	cutoff := before.UTC().Unix()
	res, err := db.Exec(rebind(deletePrunedDeliveriesSQL), cutoff)
	if err != nil {
		return result, fmt.Errorf("prune notification deliveries: %w", err)
	}
	if n, _ := res.RowsAffected(); n > 0 {
		result.Deliveries = int(n)
	}
	res, err = db.Exec(rebind(deletePrunedNotificationsSQL), cutoff)
	if err != nil {
		return result, fmt.Errorf("prune notifications: %w", err)
	}
	if n, _ := res.RowsAffected(); n > 0 {
		result.Notifications = int(n)
	}
	return result, nil
}

func (s *SQLiteStore) PruneNotificationHistory(before time.Time) (NotificationPruneResult, error) {
	return pruneNotificationHistory(s.db, before, func(q string) string { return q })
}

func (s *PostgresStore) PruneNotificationHistory(before time.Time) (NotificationPruneResult, error) {
	return pruneNotificationHistory(s.db, before, rebindDollars)
}

const prunePushSubscriptionsSQL = `DELETE FROM notification_push_subscriptions
WHERE failures >= ? AND updated_at < ?`

func prunePushSubscriptions(db *sql.DB, query string, before time.Time, maxFailures int) (int, error) {
	if maxFailures <= 0 {
		maxFailures = 10
	}
	res, err := db.Exec(query, maxFailures, before.UTC().Unix())
	if err != nil {
		return 0, fmt.Errorf("prune push subscriptions: %w", err)
	}
	n, _ := res.RowsAffected()
	return int(n), nil
}

func (s *SQLiteStore) PrunePushSubscriptions(before time.Time, maxFailures int) (int, error) {
	return prunePushSubscriptions(s.db, prunePushSubscriptionsSQL, before, maxFailures)
}

func (s *PostgresStore) PrunePushSubscriptions(before time.Time, maxFailures int) (int, error) {
	return prunePushSubscriptions(s.db, rebindDollars(prunePushSubscriptionsSQL), before, maxFailures)
}
