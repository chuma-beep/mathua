package storage

import (
	"database/sql"
	"fmt"
	"time"
)

// announcementAccess is the shared SQLite/PostgreSQL announcement
// implementation. Statements are written with ? and rebound for PostgreSQL.
type announcementAccess struct {
	db     *sql.DB
	rebind func(string) string
}

func (s *SQLiteStore) announcement() *announcementAccess {
	return &announcementAccess{db: s.db, rebind: rebindQuestion}
}

func (s *PostgresStore) announcement() *announcementAccess {
	return &announcementAccess{db: s.db, rebind: rebindDollars}
}

func (s *SQLiteStore) ListPublishedAnnouncements(studentID string) ([]Announcement, error) {
	return s.announcement().listPublished(studentID)
}
func (s *PostgresStore) ListPublishedAnnouncements(studentID string) ([]Announcement, error) {
	return s.announcement().listPublished(studentID)
}
func (s *SQLiteStore) MarkAnnouncementsRead(studentID string, now time.Time) error {
	return s.announcement().markRead(studentID, now)
}
func (s *PostgresStore) MarkAnnouncementsRead(studentID string, now time.Time) error {
	return s.announcement().markRead(studentID, now)
}
func (s *SQLiteStore) ClearAnnouncements(studentID string, clearedAt time.Time) error {
	return s.announcement().clear(studentID, clearedAt)
}
func (s *PostgresStore) ClearAnnouncements(studentID string, clearedAt time.Time) error {
	return s.announcement().clear(studentID, clearedAt)
}
func (s *SQLiteStore) ListAnnouncements() ([]Announcement, error) {
	return s.announcement().listAll()
}
func (s *PostgresStore) ListAnnouncements() ([]Announcement, error) {
	return s.announcement().listAll()
}
func (s *SQLiteStore) GetAnnouncement(id string) (*Announcement, error) {
	return s.announcement().get(id)
}
func (s *PostgresStore) GetAnnouncement(id string) (*Announcement, error) {
	return s.announcement().get(id)
}
func (s *SQLiteStore) CreateAnnouncement(a Announcement) (*Announcement, error) {
	return s.announcement().create(a)
}
func (s *PostgresStore) CreateAnnouncement(a Announcement) (*Announcement, error) {
	return s.announcement().create(a)
}
func (s *SQLiteStore) UpdateAnnouncement(id, title, body string, published bool, now time.Time) (*Announcement, error) {
	return s.announcement().update(id, title, body, published, now)
}
func (s *PostgresStore) UpdateAnnouncement(id, title, body string, published bool, now time.Time) (*Announcement, error) {
	return s.announcement().update(id, title, body, published, now)
}

func (a *announcementAccess) listPublished(studentID string) ([]Announcement, error) {
	q := a.rebind(`
		SELECT a.id, a.title, a.body, a.published, a.created_at, a.published_at, a.created_by,
		       CASE WHEN r.announcement_id IS NULL THEN 0 ELSE 1 END
		FROM announcements a
		LEFT JOIN announcement_reads r
		  ON r.announcement_id = a.id AND r.student_id = ?
		LEFT JOIN announcement_views v
		  ON v.student_id = ?
		WHERE a.published = 1
		  AND a.published_at > COALESCE(v.cleared_at, 0)
		ORDER BY a.published_at DESC, a.created_at DESC`)
	rows, err := a.db.Query(q, studentID, studentID)
	if err != nil {
		return nil, fmt.Errorf("list announcements: %w", err)
	}
	defer rows.Close()
	var out []Announcement
	for rows.Next() {
		var item Announcement
		var created, publishedAt int64
		var published, read int
		if err := rows.Scan(&item.ID, &item.Title, &item.Body, &published, &created, &publishedAt, &item.CreatedBy, &read); err != nil {
			return nil, fmt.Errorf("scan announcement: %w", err)
		}
		item.Published = published != 0
		item.CreatedAt = announcementTime(created)
		item.PublishedAt = announcementTime(publishedAt)
		item.Read = read != 0
		out = append(out, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list announcements: %w", err)
	}
	return out, nil
}

func (a *announcementAccess) markRead(studentID string, now time.Time) error {
	q := a.rebind(`
		INSERT INTO announcement_reads (student_id, announcement_id, read_at)
		SELECT ?, id, ? FROM announcements WHERE published = 1
		ON CONFLICT (student_id, announcement_id) DO NOTHING`)
	if _, err := a.db.Exec(q, studentID, now.UTC().UnixNano()); err != nil {
		return fmt.Errorf("mark announcements read: %w", err)
	}
	return nil
}

func (a *announcementAccess) clear(studentID string, clearedAt time.Time) error {
	q := a.rebind(`
		INSERT INTO announcement_views (student_id, cleared_at) VALUES (?, ?)
		ON CONFLICT (student_id) DO UPDATE SET cleared_at = excluded.cleared_at`)
	if _, err := a.db.Exec(q, studentID, clearedAt.UTC().UnixNano()); err != nil {
		return fmt.Errorf("clear announcements: %w", err)
	}
	return nil
}

func (a *announcementAccess) listAll() ([]Announcement, error) {
	q := a.rebind(`
		SELECT id, title, body, published, created_at, published_at, created_by
		FROM announcements
		ORDER BY created_at DESC`)
	rows, err := a.db.Query(q)
	if err != nil {
		return nil, fmt.Errorf("list admin announcements: %w", err)
	}
	defer rows.Close()
	var out []Announcement
	for rows.Next() {
		var item Announcement
		var created, publishedAt int64
		var published int
		if err := rows.Scan(&item.ID, &item.Title, &item.Body, &published, &created, &publishedAt, &item.CreatedBy); err != nil {
			return nil, fmt.Errorf("scan admin announcement: %w", err)
		}
		item.Published = published != 0
		item.CreatedAt = announcementTime(created)
		item.PublishedAt = announcementTime(publishedAt)
		out = append(out, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list admin announcements: %w", err)
	}
	return out, nil
}

func (a *announcementAccess) get(id string) (*Announcement, error) {
	q := a.rebind(`
		SELECT id, title, body, published, created_at, published_at, created_by
		FROM announcements WHERE id = ?`)
	var item Announcement
	var created, publishedAt int64
	var published int
	err := a.db.QueryRow(q, id).Scan(&item.ID, &item.Title, &item.Body, &published, &created, &publishedAt, &item.CreatedBy)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get announcement: %w", err)
	}
	item.Published = published != 0
	item.CreatedAt = announcementTime(created)
	item.PublishedAt = announcementTime(publishedAt)
	return &item, nil
}

func (a *announcementAccess) create(item Announcement) (*Announcement, error) {
	now := time.Now().UTC()
	if item.ID == "" {
		item.ID = newUUID()
	}
	if item.CreatedAt.IsZero() {
		item.CreatedAt = now
	}
	if item.Published && item.PublishedAt.IsZero() {
		item.PublishedAt = now
	}
	q := a.rebind(`
		INSERT INTO announcements (id, title, body, published, created_at, published_at, created_by)
		VALUES (?, ?, ?, ?, ?, ?, ?)`)
	if _, err := a.db.Exec(q, item.ID, item.Title, item.Body, boolInt(item.Published),
		item.CreatedAt.UTC().UnixNano(), announcementUnix(item.PublishedAt), item.CreatedBy); err != nil {
		return nil, fmt.Errorf("create announcement: %w", err)
	}
	return a.get(item.ID)
}

func (a *announcementAccess) update(id, title, body string, published bool, now time.Time) (*Announcement, error) {
	q := a.rebind(`
		UPDATE announcements
		SET title = ?, body = ?, published = ?,
		    published_at = CASE WHEN ? = 1 AND published_at = 0 THEN ? ELSE published_at END
		WHERE id = ?`)
	res, err := a.db.Exec(q, title, body, boolInt(published), boolInt(published), now.UTC().UnixNano(), id)
	if err != nil {
		return nil, fmt.Errorf("update announcement: %w", err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return nil, sql.ErrNoRows
	}
	return a.get(id)
}

func announcementTime(ns int64) time.Time {
	if ns <= 0 {
		return time.Time{}
	}
	return time.Unix(0, ns).UTC()
}

func announcementUnix(t time.Time) int64 {
	if t.IsZero() {
		return 0
	}
	return t.UTC().UnixNano()
}
