package storage

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
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
		"ALTER TABLE students ADD COLUMN username TEXT",
		"ALTER TABLE students ADD COLUMN password_hash TEXT",
		"ALTER TABLE students ADD COLUMN course_id TEXT",
		"ALTER TABLE students ADD COLUMN xp_total INTEGER NOT NULL DEFAULT 0",
		"ALTER TABLE students ADD COLUMN xp_today INTEGER NOT NULL DEFAULT 0",
		"ALTER TABLE students ADD COLUMN xp_date TEXT",
		"ALTER TABLE students ADD COLUMN diagnostic_completed INTEGER NOT NULL DEFAULT 0",
		"ALTER TABLE students ADD COLUMN daily_xp_goal INTEGER NOT NULL DEFAULT 150",
		"ALTER TABLE students ADD COLUMN settings TEXT NOT NULL DEFAULT '{}'",
		"CREATE INDEX IF NOT EXISTS idx_students_username ON students(username)",
		"ALTER TABLE concept_progress ADD COLUMN weakness_score REAL NOT NULL DEFAULT 0",
	}
	for _, m := range migrations {
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
		"INSERT INTO students (id, name, settings, created_at) VALUES (?, ?, '{}', ?)",
		id, name, now.Format(time.RFC3339),
	)
	if err != nil {
		return nil, fmt.Errorf("create student: %w", err)
	}
	return &Student{ID: id, Name: name, Settings: "{}", CreatedAt: now}, nil
}

func (s *SQLiteStore) CreateUser(name, username, passwordHash string) (*Student, error) {
	id := newUUID()
	now := time.Now().UTC()
	_, err := s.db.Exec(
		"INSERT INTO students (id, name, username, password_hash, settings, created_at) VALUES (?, ?, ?, ?, '{}', ?)",
		id, name, username, passwordHash, now.Format(time.RFC3339),
	)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return &Student{ID: id, Name: name, Username: username, PasswordHash: passwordHash, Settings: "{}", CreatedAt: now}, nil
}

func scanStudent(row interface{ Scan(...interface{}) error }) (*Student, error) {
	var st Student
	var username, passwordHash, courseID, xpDate, settings sql.NullString
	var xpTotal, xpToday, diagCompleted, dailyGoal sql.NullInt64
	var createdAt string
	err := row.Scan(&st.ID, &st.Name, &username, &passwordHash, &courseID, &xpTotal, &xpToday, &xpDate, &diagCompleted, &dailyGoal, &settings, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	st.Username = username.String
	st.PasswordHash = passwordHash.String
	st.CourseID = courseID.String
	if xpTotal.Valid {
		st.XPTotal = int(xpTotal.Int64)
	}
	if xpToday.Valid {
		st.XPToday = int(xpToday.Int64)
	}
	st.XPTodayDate = xpDate.String
	st.DiagnosticCompleted = diagCompleted.Valid && diagCompleted.Int64 == 1
	if dailyGoal.Valid {
		st.DailyXPGoal = int(dailyGoal.Int64)
	}
	st.Settings = settings.String
	if st.Settings == "" {
		st.Settings = "{}"
	}
	if createdAt != "" {
		st.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
		if err != nil {
			st.CreatedAt = time.Now()
		}
	}
	return &st, nil
}

func (s *SQLiteStore) GetStudent(id string) (*Student, error) {
	row := s.db.QueryRow("SELECT id, name, username, password_hash, course_id, xp_total, xp_today, xp_date, diagnostic_completed, daily_xp_goal, settings, created_at FROM students WHERE id = ?", id)
	return scanStudent(row)
}

func (s *SQLiteStore) FindByUsername(username string) (*Student, error) {
	row := s.db.QueryRow("SELECT id, name, username, password_hash, course_id, xp_total, xp_today, xp_date, diagnostic_completed, daily_xp_goal, settings, created_at FROM students WHERE username = ?", username)
	return scanStudent(row)
}

func (s *SQLiteStore) SetCourseID(studentID, courseID string) error {
	_, err := s.db.Exec("UPDATE students SET course_id = ? WHERE id = ?", courseID, studentID)
	if err != nil {
		return fmt.Errorf("set course_id: %w", err)
	}
	return nil
}

func (s *SQLiteStore) AddXP(studentID string, amount int) error {
	today := time.Now().UTC().Format("2006-01-02")
	_, err := s.db.Exec(`
		UPDATE students
		SET xp_total = xp_total + ?,
		    xp_today = CASE WHEN xp_date = ? THEN xp_today + ? ELSE ? END,
		    xp_date  = ?
		WHERE id = ?
	`, amount, today, amount, amount, today, studentID)
	if err != nil {
		return fmt.Errorf("add xp: %w", err)
	}
	return nil
}

func (s *SQLiteStore) GetXP(studentID string) (int, int, error) {
	row := s.db.QueryRow("SELECT xp_total, xp_today, xp_date FROM students WHERE id = ?", studentID)
	var total, today sql.NullInt64
	var xpDate sql.NullString
	if err := row.Scan(&total, &today, &xpDate); err != nil {
		if err == sql.ErrNoRows {
			return 0, 0, nil
		}
		return 0, 0, fmt.Errorf("get xp: %w", err)
	}
	if xpDate.String != time.Now().UTC().Format("2006-01-02") {
		today.Int64 = 0
	}
	return int(total.Int64), int(today.Int64), nil
}

func (s *SQLiteStore) SetDiagnosticCompleted(studentID string) error {
	_, err := s.db.Exec("UPDATE students SET diagnostic_completed = 1 WHERE id = ?", studentID)
	if err != nil {
		return fmt.Errorf("set diagnostic_completed: %w", err)
	}
	return nil
}

func (s *SQLiteStore) SetDailyXPGoal(studentID string, goal int) error {
	_, err := s.db.Exec("UPDATE students SET daily_xp_goal = ? WHERE id = ?", goal, studentID)
	if err != nil {
		return fmt.Errorf("set daily xp goal: %w", err)
	}
	return nil
}

func (s *SQLiteStore) GetSettings(studentID string) (string, error) {
	row := s.db.QueryRow("SELECT settings FROM students WHERE id = ?", studentID)
	var settings string
	if err := row.Scan(&settings); err != nil {
		if err == sql.ErrNoRows {
			return "{}", nil
		}
		return "", fmt.Errorf("get settings: %w", err)
	}
	if settings == "" {
		return "{}", nil
	}
	return settings, nil
}

func (s *SQLiteStore) UpdateSettings(studentID string, settings string) error {
	_, err := s.db.Exec("UPDATE students SET settings = ? WHERE id = ?", settings, studentID)
	if err != nil {
		return fmt.Errorf("update settings: %w", err)
	}
	return nil
}

// Progress

func (s *SQLiteStore) GetProgress(studentID, conceptID string) (*ConceptProgress, error) {
	row := s.db.QueryRow(`
		SELECT student_id, concept_id, status, streak, best_streak,
		       avg_response_time, attempts, last_attempted, last_reviewed,
		       next_review_due, sm2_repetitions, sm2_interval, sm2_efactor,
		       mastered_at, weakness_score
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
		       mastered_at, weakness_score
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
			 mastered_at, weakness_score)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
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
			mastered_at       = excluded.mastered_at,
			weakness_score    = excluded.weakness_score
	`,
		p.StudentID, p.ConceptID, p.Status, p.Streak, p.BestStreak,
		p.AvgResponseTime, p.Attempts,
		nullTime(p.LastAttempted), nullTime(p.LastReviewed),
		nullTime(p.NextReviewDue), p.SM2Repetitions, p.SM2Interval,
		p.SM2EFactor, nullTime(p.MasteredAt), p.WeaknessScore,
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
	if startedAt != "" {
		if t, err := time.Parse(time.RFC3339, startedAt); err == nil {
			ses.StartedAt = t
		}
	}
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
		if ts != "" {
			if t, err := time.Parse(time.RFC3339, ts); err == nil {
				e.Timestamp = t
			}
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

// Questions

func (s *SQLiteStore) GetQuestions(conceptID string, count int) ([]Question, error) {
	rows, err := s.db.Query(`
		SELECT id, concept_id, question, answer, explanation, source, difficulty
		FROM questions
		WHERE concept_id = ?
		ORDER BY RANDOM()
		LIMIT ?
	`, conceptID, count)
	if err != nil {
		return nil, fmt.Errorf("get questions: %w", err)
	}
	defer rows.Close()

	var out []Question
	for rows.Next() {
		var q Question
		if err := rows.Scan(&q.ID, &q.ConceptID, &q.Question, &q.Answer, &q.Explanation, &q.Source, &q.Difficulty); err != nil {
			return nil, fmt.Errorf("scan question: %w", err)
		}
		out = append(out, q)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) GetQuestionCount(conceptID string) (int, error) {
	row := s.db.QueryRow("SELECT COUNT(*) FROM questions WHERE concept_id = ?", conceptID)
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("get question count: %w", err)
	}
	return count, nil
}

func (s *SQLiteStore) ImportQuestions(qs []Question) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			log.Printf("warning: transaction rollback failed: %v", err)
		}
	}()

	stmt, err := tx.Prepare(`
		INSERT OR IGNORE INTO questions (concept_id, question, answer, explanation, source, difficulty)
		VALUES (?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	for _, q := range qs {
		if _, err := stmt.Exec(q.ConceptID, q.Question, q.Answer, q.Explanation, q.Source, q.Difficulty); err != nil {
			return fmt.Errorf("insert question: %w", err)
		}
	}
	return tx.Commit()
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

// DailyActivity

func (s *SQLiteStore) GetDailyActivity(studentID string, days int) ([]DailyActivity, error) {
	limit := fmt.Sprintf("-%d days", days)
	rows, err := s.db.Query(`
		SELECT date(timestamp) as day,
		       COUNT(*) as total,
		       COALESCE(SUM(correct), 0) as correct_count,
		       COALESCE(GROUP_CONCAT(DISTINCT concept_id), '') as concepts
		FROM attempts
		WHERE student_id = ? AND timestamp >= datetime('now', ?)
		GROUP BY day
		ORDER BY day
	`, studentID, limit)
	if err != nil {
		return nil, fmt.Errorf("get daily activity: %w", err)
	}
	defer rows.Close()

	var out []DailyActivity
	for rows.Next() {
		var da DailyActivity
		var conceptsStr string
		if err := rows.Scan(&da.Date, &da.Questions, &da.Correct, &conceptsStr); err != nil {
			return nil, fmt.Errorf("scan daily activity: %w", err)
		}
		if conceptsStr != "" {
			for _, cid := range strings.Split(conceptsStr, ",") {
				cid = strings.TrimSpace(cid)
				if cid != "" {
					da.Concepts = append(da.Concepts, cid)
				}
			}
		} else {
			da.Concepts = []string{}
		}
		out = append(out, da)
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
		&masterAt, &p.WeaknessScore,
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

func (s *SQLiteStore) PurgeGeneratedQuestions(conceptIDs map[string]bool) (int64, error) {
	if len(conceptIDs) == 0 {
		return 0, nil
	}
	ids := make([]string, 0, len(conceptIDs))
	for id := range conceptIDs {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	var deleted int64
	batch := 400
	for start := 0; start < len(ids); start += batch {
		end := start + batch
		if end > len(ids) {
			end = len(ids)
		}
		chunk := ids[start:end]
		placeholders := strings.Repeat("?,", len(chunk))
		placeholders = placeholders[:len(placeholders)-1]
		args := make([]interface{}, len(chunk))
		for i, id := range chunk {
			args[i] = id
		}
		res, err := s.db.Exec("DELETE FROM questions WHERE concept_id IN ("+placeholders+")", args...)
		if err != nil {
			return deleted, fmt.Errorf("purge questions: %w", err)
		}
		n, _ := res.RowsAffected()
		deleted += n
	}
	return deleted, nil
}
