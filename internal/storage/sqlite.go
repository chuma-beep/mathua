package storage

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(path string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		return nil, fmt.Errorf("enable WAL: %w", err)
	}
	if _, err := db.Exec("PRAGMA foreign_keys=ON"); err != nil {
		return nil, fmt.Errorf("enable foreign keys: %w", err)
	}
	if _, err := db.Exec("PRAGMA busy_timeout=5000"); err != nil {
		return nil, fmt.Errorf("set busy timeout: %w", err)
	}
	store := &SQLiteStore{db: db}
	if err := store.Migrate(); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return store, nil
}

func (s *SQLiteStore) Migrate() error {
	if _, err := s.db.Exec(schema); err != nil {
		return err
	}
	return authMigrate(s.db)
}

func authMigrate(db *sql.DB) error {
	migrations := []string{
		`ALTER TABLE students ADD COLUMN username TEXT`,
		`ALTER TABLE students ADD COLUMN password_hash TEXT`,
		`CREATE INDEX IF NOT EXISTS idx_students_username ON students(username)`,
	}
	for _, m := range migrations {
		// Ignore "duplicate column" errors for existing DBs.
		if _, err := db.Exec(m); err != nil {
			if !isDuplicateColumn(err) {
				return err
			}
		}
	}
	return nil
}

func isDuplicateColumn(err error) bool {
	return strings.Contains(err.Error(), "duplicate column name")
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

// Students

func (s *SQLiteStore) CreateStudent(name string) (*Student, error) {
	id := newUUID()
	now := time.Now().UTC()
	_, err := s.db.Exec(
		"INSERT INTO students (id, name, created_at) VALUES (?, ?, ?)",
		id, name, now.Format(time.RFC3339),
	)
	if err != nil {
		return nil, fmt.Errorf("create student: %w", err)
	}
	return &Student{ID: id, Name: name, CreatedAt: now}, nil
}

func (s *SQLiteStore) CreateUser(name, username, passwordHash string) (*Student, error) {
	id := newUUID()
	now := time.Now().UTC()
	_, err := s.db.Exec(
		"INSERT INTO students (id, name, username, password_hash, created_at) VALUES (?, ?, ?, ?, ?)",
		id, name, username, passwordHash, now.Format(time.RFC3339),
	)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return &Student{ID: id, Name: name, Username: username, PasswordHash: passwordHash, CreatedAt: now}, nil
}

func (s *SQLiteStore) GetStudent(id string) (*Student, error) {
	row := s.db.QueryRow("SELECT id, name, username, password_hash, created_at FROM students WHERE id = ?", id)
	var st Student
	var username, passwordHash sql.NullString
	var createdAt string
	if err := row.Scan(&st.ID, &st.Name, &username, &passwordHash, &createdAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get student: %w", err)
	}
	st.Username = username.String
	st.PasswordHash = passwordHash.String
	st.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	return &st, nil
}

func (s *SQLiteStore) FindByUsername(username string) (*Student, error) {
	row := s.db.QueryRow("SELECT id, name, username, password_hash, created_at FROM students WHERE username = ?", username)
	var st Student
	var un, ph sql.NullString
	var createdAt string
	if err := row.Scan(&st.ID, &st.Name, &un, &ph, &createdAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("find by username: %w", err)
	}
	st.Username = un.String
	st.PasswordHash = ph.String
	st.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	return &st, nil
}

// Progress

func (s *SQLiteStore) GetProgress(studentID, conceptID string) (*ConceptProgress, error) {
	row := s.db.QueryRow(`
		SELECT student_id, concept_id, status, streak, best_streak,
		       avg_response_time, attempts, last_attempted, last_reviewed,
		       next_review_due, sm2_repetitions, sm2_interval, sm2_efactor,
		       mastered_at
		FROM concept_progress
		WHERE student_id = ? AND concept_id = ?
	`, studentID, conceptID)

	p, err := scanProgress(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return p, err
}

func (s *SQLiteStore) GetAllProgress(studentID string) (map[string]*ConceptProgress, error) {
	rows, err := s.db.Query(`
		SELECT student_id, concept_id, status, streak, best_streak,
		       avg_response_time, attempts, last_attempted, last_reviewed,
		       next_review_due, sm2_repetitions, sm2_interval, sm2_efactor,
		       mastered_at
		FROM concept_progress
		WHERE student_id = ?
	`, studentID)
	if err != nil {
		return nil, fmt.Errorf("get all progress: %w", err)
	}
	defer rows.Close()

	out := make(map[string]*ConceptProgress)
	for rows.Next() {
		p, err := scanProgress(rows)
		if err != nil {
			return nil, fmt.Errorf("scan progress: %w", err)
		}
		out[p.ConceptID] = p
	}
	return out, rows.Err()
}

func (s *SQLiteStore) UpsertProgress(p *ConceptProgress) error {
	_, err := s.db.Exec(`
		INSERT INTO concept_progress
			(student_id, concept_id, status, streak, best_streak,
			 avg_response_time, attempts, last_attempted, last_reviewed,
			 next_review_due, sm2_repetitions, sm2_interval, sm2_efactor,
			 mastered_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(student_id, concept_id) DO UPDATE SET
			status            = excluded.status,
			streak            = excluded.streak,
			best_streak       = excluded.best_streak,
			avg_response_time = excluded.avg_response_time,
			attempts          = excluded.attempts,
			last_attempted    = excluded.last_attempted,
			last_reviewed     = excluded.last_reviewed,
			next_review_due   = excluded.next_review_due,
			sm2_repetitions   = excluded.sm2_repetitions,
			sm2_interval      = excluded.sm2_interval,
			sm2_efactor       = excluded.sm2_efactor,
			mastered_at       = excluded.mastered_at
	`,
		p.StudentID, p.ConceptID, p.Status, p.Streak, p.BestStreak,
		p.AvgResponseTime, p.Attempts,
		nullTime(p.LastAttempted), nullTime(p.LastReviewed),
		nullTime(p.NextReviewDue), p.SM2Repetitions, p.SM2Interval,
		p.SM2EFactor, nullTime(p.MasteredAt),
	)
	if err != nil {
		return fmt.Errorf("upsert progress: %w", err)
	}
	return nil
}

// Sessions

func (s *SQLiteStore) CreateSession(studentID string) (*Session, error) {
	id := newUUID()
	now := time.Now().UTC()
	_, err := s.db.Exec(
		"INSERT INTO sessions (id, student_id, started_at) VALUES (?, ?, ?)",
		id, studentID, now.Format(time.RFC3339),
	)
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}
	return &Session{ID: id, StudentID: studentID, StartedAt: now}, nil
}

func (s *SQLiteStore) GetSession(id string) (*Session, error) {
	row := s.db.QueryRow("SELECT id, student_id, started_at FROM sessions WHERE id = ?", id)
	var ses Session
	var startedAt string
	if err := row.Scan(&ses.ID, &ses.StudentID, &startedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get session: %w", err)
	}
	ses.StartedAt, _ = time.Parse(time.RFC3339, startedAt)
	return &ses, nil
}

// Attempts

func (s *SQLiteStore) RecordAttempt(entry AttemptEntry) error {
	_, err := s.db.Exec(`
		INSERT INTO attempts
			(session_id, student_id, concept_id, answer, expected,
			 correct, elapsed_seconds, timestamp)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`,
		entry.SessionID, entry.StudentID, entry.ConceptID,
		entry.Answer, entry.Expected,
		boolToInt(entry.Correct), entry.ElapsedSeconds,
		entry.Timestamp.Format(time.RFC3339),
	)
	if err != nil {
		return fmt.Errorf("record attempt: %w", err)
	}
	return nil
}

func (s *SQLiteStore) GetSessionAttempts(studentID, sessionID string) ([]AttemptEntry, error) {
	rows, err := s.db.Query(`
		SELECT session_id, student_id, concept_id, answer, expected,
		       correct, elapsed_seconds, timestamp
		FROM attempts
		WHERE student_id = ? AND session_id = ?
		ORDER BY timestamp ASC
	`, studentID, sessionID)
	if err != nil {
		return nil, fmt.Errorf("get session attempts: %w", err)
	}
	defer rows.Close()

	var out []AttemptEntry
	for rows.Next() {
		var e AttemptEntry
		var correct int
		var ts string
		if err := rows.Scan(
			&e.SessionID, &e.StudentID, &e.ConceptID,
			&e.Answer, &e.Expected, &correct, &e.ElapsedSeconds, &ts,
		); err != nil {
			return nil, fmt.Errorf("scan attempt: %w", err)
		}
		e.Correct = correct != 0
		e.Timestamp, _ = time.Parse(time.RFC3339, ts)
		out = append(out, e)
	}
	return out, rows.Err()
}

// Leaderboard

func (s *SQLiteStore) GetWeeklyLeaderboard() ([]LeaderboardRow, error) {
	monday := weekStart(time.Now().UTC())
	rows, err := s.db.Query(`
		SELECT s.id, s.name,
			COALESCE((SELECT COUNT(*) FROM concept_progress
			          WHERE student_id = s.id AND status = 'MASTERED'), 0),
			COALESCE((SELECT COUNT(*) FROM concept_progress
			          WHERE student_id = s.id AND status = 'MASTERED'
			          AND mastered_at >= ?), 0)
		FROM students s
		ORDER BY 4 DESC, 3 DESC
	`, monday.Format(time.RFC3339))
	if err != nil {
		return nil, fmt.Errorf("get weekly leaderboard: %w", err)
	}
	defer rows.Close()

	var out []LeaderboardRow
	for rows.Next() {
		var r LeaderboardRow
		if err := rows.Scan(&r.StudentID, &r.Name, &r.TotalMastered, &r.WeeklyMastered); err != nil {
			return nil, fmt.Errorf("scan leaderboard: %w", err)
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// Helpers

func scanProgress(scanner interface{ Scan(...interface{}) error }) (*ConceptProgress, error) {
	p := &ConceptProgress{}
	var lastAtt, lastRev, nextRev, masterAt sql.NullString
	err := scanner.Scan(
		&p.StudentID, &p.ConceptID, &p.Status, &p.Streak, &p.BestStreak,
		&p.AvgResponseTime, &p.Attempts, &lastAtt, &lastRev,
		&nextRev, &p.SM2Repetitions, &p.SM2Interval, &p.SM2EFactor,
		&masterAt,
	)
	if err != nil {
		return nil, err
	}
	p.LastAttempted = parseDate(lastAtt)
	p.LastReviewed = parseDate(lastRev)
	p.NextReviewDue = parseDate(nextRev)
	p.MasteredAt = parseDate(masterAt)
	return p, nil
}

func parseDate(ns sql.NullString) *time.Time {
	if !ns.Valid {
		return nil
	}
	t, err := time.Parse(time.RFC3339, ns.String)
	if err != nil {
		return nil
	}
	return &t
}

func nullTime(t *time.Time) interface{} {
	if t == nil {
		return nil
	}
	return t.Format(time.RFC3339)
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func weekStart(t time.Time) time.Time {
	weekday := t.Weekday()
	if weekday == time.Sunday {
		weekday = 7
	}
	offset := int(weekday) - int(time.Monday)
	start := t.AddDate(0, 0, -offset)
	return time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
}
