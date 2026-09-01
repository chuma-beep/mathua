package storage

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresStore struct {
	db *sql.DB
}

const pgSchema = `
CREATE TABLE IF NOT EXISTS students (
    id                   TEXT PRIMARY KEY,
    name                 TEXT NOT NULL,
    username             TEXT,
    password_hash        TEXT,
    created_at           TEXT NOT NULL DEFAULT (now()::text),
    xp_total             INTEGER NOT NULL DEFAULT 0,
    xp_today             INTEGER NOT NULL DEFAULT 0,
    xp_date              TEXT,
    diagnostic_completed INTEGER NOT NULL DEFAULT 0,
    daily_xp_goal        INTEGER NOT NULL DEFAULT 30,
    settings             TEXT NOT NULL DEFAULT '{}',
    league               TEXT NOT NULL DEFAULT 'bronze',
    league_week          TEXT NOT NULL DEFAULT '',
    league_moved         INTEGER NOT NULL DEFAULT 0,
    share_token          TEXT NOT NULL DEFAULT '',
    course_id            TEXT
);

CREATE TABLE IF NOT EXISTS concept_progress (
    student_id       TEXT    NOT NULL,
    concept_id       TEXT    NOT NULL,
    status           TEXT    NOT NULL DEFAULT 'UNSEEN',
    streak           INTEGER NOT NULL DEFAULT 0,
    best_streak      INTEGER NOT NULL DEFAULT 0,
    avg_response_time DOUBLE PRECISION NOT NULL DEFAULT 0,
    attempts         INTEGER NOT NULL DEFAULT 0,
    last_attempted   TEXT,
    last_reviewed    TEXT,
    next_review_due  TEXT,
    sm2_repetitions  INTEGER NOT NULL DEFAULT 0,
    sm2_interval     INTEGER NOT NULL DEFAULT 0,
    sm2_efactor      DOUBLE PRECISION NOT NULL DEFAULT 2.5,
    mastered_at      TEXT,
    weakness_score   DOUBLE PRECISION NOT NULL DEFAULT 0,
    PRIMARY KEY (student_id, concept_id)
);

CREATE TABLE IF NOT EXISTS sessions (
    id         TEXT PRIMARY KEY,
    student_id TEXT NOT NULL,
    started_at TEXT NOT NULL DEFAULT (now()::text)
);

CREATE TABLE IF NOT EXISTS attempts (
    id              SERIAL PRIMARY KEY,
    session_id      TEXT    NOT NULL,
    student_id      TEXT    NOT NULL,
    concept_id      TEXT    NOT NULL,
    answer          TEXT    NOT NULL,
    expected        TEXT    NOT NULL,
    correct         INTEGER NOT NULL DEFAULT 0,
    elapsed_seconds DOUBLE PRECISION NOT NULL DEFAULT 0,
    timestamp       TEXT    NOT NULL DEFAULT (now()::text)
);

CREATE INDEX IF NOT EXISTS idx_progress_student  ON concept_progress(student_id);
CREATE INDEX IF NOT EXISTS idx_progress_mastered ON concept_progress(student_id, status, mastered_at);
CREATE INDEX IF NOT EXISTS idx_sessions_student  ON sessions(student_id);
CREATE INDEX IF NOT EXISTS idx_attempts_session  ON attempts(session_id);
CREATE INDEX IF NOT EXISTS idx_attempts_student  ON attempts(student_id, timestamp);
CREATE INDEX IF NOT EXISTS idx_attempts_cover    ON attempts(student_id, concept_id, timestamp);
CREATE INDEX IF NOT EXISTS idx_students_share    ON students(share_token);
CREATE INDEX IF NOT EXISTS idx_students_username ON students(username);

CREATE TABLE IF NOT EXISTS questions (
    id         SERIAL PRIMARY KEY,
    concept_id TEXT    NOT NULL,
    question   TEXT    NOT NULL,
    answer     TEXT    NOT NULL,
    explanation TEXT   NOT NULL DEFAULT '',
    source     TEXT    NOT NULL DEFAULT '',
    difficulty DOUBLE PRECISION NOT NULL DEFAULT 0.5
);

CREATE INDEX IF NOT EXISTS idx_questions_concept ON questions(concept_id);

CREATE TABLE IF NOT EXISTS active_sessions (
    session_id      TEXT PRIMARY KEY,
    student_id      TEXT NOT NULL,
    concept_id      TEXT NOT NULL,
    concept_name    TEXT NOT NULL DEFAULT '',
    expected_answer TEXT NOT NULL,
    attempt_id      TEXT NOT NULL UNIQUE,
    question        TEXT NOT NULL,
    explanation     TEXT NOT NULL DEFAULT '',
    diagram         TEXT NOT NULL DEFAULT '',
    is_review       INTEGER NOT NULL DEFAULT 0,
    answered        INTEGER NOT NULL DEFAULT 0,
    last_concept_id TEXT NOT NULL DEFAULT '',
    session_review  INTEGER NOT NULL DEFAULT 0,
    session_new     INTEGER NOT NULL DEFAULT 0,
    updated_at      TEXT NOT NULL DEFAULT (now()::text)
);

CREATE INDEX IF NOT EXISTS idx_active_sessions_student ON active_sessions(student_id);
CREATE INDEX IF NOT EXISTS idx_active_sessions_attempt ON active_sessions(attempt_id);

CREATE TABLE IF NOT EXISTS student_topic_speed (
    student_id    TEXT NOT NULL,
    concept_id    TEXT NOT NULL,
    efactor       DOUBLE PRECISION NOT NULL DEFAULT 2.5,
    interval      INTEGER NOT NULL DEFAULT 0,
    repetitions   INTEGER NOT NULL DEFAULT 0,
    learning_speed DOUBLE PRECISION NOT NULL DEFAULT 1.0,
    updated_at    TEXT NOT NULL DEFAULT (now()::text),
    PRIMARY KEY (student_id, concept_id)
);
CREATE INDEX IF NOT EXISTS idx_topic_speed_student ON student_topic_speed(student_id);
`

func NewPostgresStore(dsn string) (*PostgresStore, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping postgres: %w", err)
	}
	store := &PostgresStore{db: db}
	if err := store.Migrate(); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return store, nil
}

func (s *PostgresStore) Migrate() error {
	if _, err := s.db.Exec(pgSchema); err != nil {
		return err
	}
	return pgAuthMigrate(s.db)
}

func pgAuthMigrate(db *sql.DB) error {
	migrations := []string{
		"ALTER TABLE students ADD COLUMN IF NOT EXISTS username TEXT",
		"ALTER TABLE students ADD COLUMN IF NOT EXISTS password_hash TEXT",
		"ALTER TABLE students ADD COLUMN IF NOT EXISTS course_id TEXT",
		"ALTER TABLE students ADD COLUMN IF NOT EXISTS xp_total INTEGER NOT NULL DEFAULT 0",
		"ALTER TABLE students ADD COLUMN IF NOT EXISTS xp_today INTEGER NOT NULL DEFAULT 0",
		"ALTER TABLE students ADD COLUMN IF NOT EXISTS xp_date TEXT",
		"ALTER TABLE students ADD COLUMN IF NOT EXISTS diagnostic_completed INTEGER NOT NULL DEFAULT 0",
		"ALTER TABLE students ADD COLUMN IF NOT EXISTS daily_xp_goal INTEGER NOT NULL DEFAULT 30",
		"ALTER TABLE students ADD COLUMN IF NOT EXISTS settings TEXT NOT NULL DEFAULT '{}'",
		"ALTER TABLE students ADD COLUMN IF NOT EXISTS league TEXT NOT NULL DEFAULT 'bronze'",
		"ALTER TABLE students ADD COLUMN IF NOT EXISTS league_week TEXT NOT NULL DEFAULT ''",
		"ALTER TABLE students ADD COLUMN IF NOT EXISTS league_moved INTEGER NOT NULL DEFAULT 0",
		"ALTER TABLE students ADD COLUMN IF NOT EXISTS share_token TEXT NOT NULL DEFAULT ''",
		"CREATE INDEX IF NOT EXISTS idx_students_username ON students(username)",
		"ALTER TABLE concept_progress ADD COLUMN IF NOT EXISTS weakness_score DOUBLE PRECISION NOT NULL DEFAULT 0",
		"ALTER TABLE active_sessions ADD COLUMN IF NOT EXISTS last_concept_id TEXT NOT NULL DEFAULT ''",
		"ALTER TABLE active_sessions ADD COLUMN IF NOT EXISTS session_review INTEGER NOT NULL DEFAULT 0",
		"ALTER TABLE active_sessions ADD COLUMN IF NOT EXISTS session_new INTEGER NOT NULL DEFAULT 0",
	}
	for _, m := range migrations {
		if _, err := db.Exec(m); err != nil {
			// IF NOT EXISTS handles duplicates; ignore any other duplicate errors
			if !isPgDuplicateColumn(err) && !isDuplicateColumn(err) {
				return err
			}
		}
	}
	_, _ = db.Exec("UPDATE students SET daily_xp_goal = 30 WHERE daily_xp_goal = 150")
	return nil
}

func isPgDuplicateColumn(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "already exists") || strings.Contains(msg, "duplicate")
}

func (s *PostgresStore) Close() error {
	return s.db.Close()
}

// Students

func (s *PostgresStore) CreateStudent(name string) (*Student, error) {
	id := newUUID()
	now := time.Now().UTC()
	_, err := s.db.Exec(
		"INSERT INTO students (id, name, settings, created_at) VALUES ($1, $2, '{}', $3)",
		id, name, now.Format(time.RFC3339),
	)
	if err != nil {
		return nil, fmt.Errorf("create student: %w", err)
	}
	return &Student{ID: id, Name: name, Settings: "{}", CreatedAt: now}, nil
}

func (s *PostgresStore) CreateUser(name, username, passwordHash string) (*Student, error) {
	id := newUUID()
	now := time.Now().UTC()
	_, err := s.db.Exec(
		"INSERT INTO students (id, name, username, password_hash, settings, created_at) VALUES ($1, $2, $3, $4, '{}', $5)",
		id, name, username, passwordHash, now.Format(time.RFC3339),
	)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return &Student{ID: id, Name: name, Username: username, PasswordHash: passwordHash, Settings: "{}", CreatedAt: now}, nil
}

func (s *PostgresStore) GetStudent(id string) (*Student, error) {
	row := s.db.QueryRow("SELECT id, name, username, password_hash, course_id, xp_total, xp_today, xp_date, diagnostic_completed, daily_xp_goal, settings, created_at, league, league_week, league_moved, share_token FROM students WHERE id = $1", id)
	return scanStudent(row)
}

func (s *PostgresStore) FindByUsername(username string) (*Student, error) {
	row := s.db.QueryRow("SELECT id, name, username, password_hash, course_id, xp_total, xp_today, xp_date, diagnostic_completed, daily_xp_goal, settings, created_at, league, league_week, league_moved, share_token FROM students WHERE username = $1", username)
	return scanStudent(row)
}

func (s *PostgresStore) SetShareToken(studentID, token string) error {
	_, err := s.db.Exec("UPDATE students SET share_token = $1 WHERE id = $2", token, studentID)
	if err != nil {
		return fmt.Errorf("set share token: %w", err)
	}
	return nil
}

func (s *PostgresStore) GetStudentByShareToken(token string) (*Student, error) {
	if token == "" {
		return nil, nil
	}
	row := s.db.QueryRow("SELECT id, name, username, password_hash, course_id, xp_total, xp_today, xp_date, diagnostic_completed, daily_xp_goal, settings, created_at, league, league_week, league_moved, share_token FROM students WHERE share_token = $1", token)
	return scanStudent(row)
}

func (s *PostgresStore) SetCourseID(studentID, courseID string) error {
	_, err := s.db.Exec("UPDATE students SET course_id = $1 WHERE id = $2", courseID, studentID)
	if err != nil {
		return fmt.Errorf("set course_id: %w", err)
	}
	return nil
}

func (s *PostgresStore) AddXP(studentID string, amount int) error {
	today := time.Now().UTC().Format("2006-01-02")
	_, err := s.db.Exec(`
		UPDATE students
		SET xp_total = GREATEST(0, xp_total + $1),
		    xp_today = CASE WHEN xp_date = $2 THEN GREATEST(0, xp_today + $3) ELSE GREATEST(0, $4) END,
		    xp_date  = $5
		WHERE id = $6
	`, amount, today, amount, amount, today, studentID)
	if err != nil {
		return fmt.Errorf("add xp: %w", err)
	}
	return nil
}

func (s *PostgresStore) GetXP(studentID string) (int, int, error) {
	row := s.db.QueryRow("SELECT xp_total, xp_today, xp_date FROM students WHERE id = $1", studentID)
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

func (s *PostgresStore) SetDiagnosticCompleted(studentID string) error {
	_, err := s.db.Exec("UPDATE students SET diagnostic_completed = 1 WHERE id = $1", studentID)
	if err != nil {
		return fmt.Errorf("set diagnostic_completed: %w", err)
	}
	return nil
}

func (s *PostgresStore) SetDailyXPGoal(studentID string, goal int) error {
	_, err := s.db.Exec("UPDATE students SET daily_xp_goal = $1 WHERE id = $2", goal, studentID)
	if err != nil {
		return fmt.Errorf("set daily xp goal: %w", err)
	}
	return nil
}

func (s *PostgresStore) GetSettings(studentID string) (string, error) {
	row := s.db.QueryRow("SELECT settings FROM students WHERE id = $1", studentID)
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

func (s *PostgresStore) UpdateSettings(studentID string, settings string) error {
	_, err := s.db.Exec("UPDATE students SET settings = $1 WHERE id = $2", settings, studentID)
	if err != nil {
		return fmt.Errorf("update settings: %w", err)
	}
	return nil
}

// Progress

func (s *PostgresStore) GetProgress(studentID, conceptID string) (*ConceptProgress, error) {
	row := s.db.QueryRow(`
		SELECT student_id, concept_id, status, streak, best_streak,
		       avg_response_time, attempts, last_attempted, last_reviewed,
		       next_review_due, sm2_repetitions, sm2_interval, sm2_efactor,
		       mastered_at, weakness_score
		FROM concept_progress
		WHERE student_id = $1 AND concept_id = $2
	`, studentID, conceptID)

	p, err := scanProgress(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return p, err
}

func (s *PostgresStore) GetAllProgress(studentID string) (map[string]*ConceptProgress, error) {
	rows, err := s.db.Query(`
		SELECT student_id, concept_id, status, streak, best_streak,
		       avg_response_time, attempts, last_attempted, last_reviewed,
		       next_review_due, sm2_repetitions, sm2_interval, sm2_efactor,
		       mastered_at, weakness_score
		FROM concept_progress
		WHERE student_id = $1
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

func (s *PostgresStore) UpsertProgress(p *ConceptProgress) error {
	_, err := s.db.Exec(`
		INSERT INTO concept_progress
			(student_id, concept_id, status, streak, best_streak,
			 avg_response_time, attempts, last_attempted, last_reviewed,
			 next_review_due, sm2_repetitions, sm2_interval, sm2_efactor,
			 mastered_at, weakness_score)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		ON CONFLICT(student_id, concept_id) DO UPDATE SET
			status            = EXCLUDED.status,
			streak            = EXCLUDED.streak,
			best_streak       = EXCLUDED.best_streak,
			avg_response_time = EXCLUDED.avg_response_time,
			attempts          = EXCLUDED.attempts,
			last_attempted    = EXCLUDED.last_attempted,
			last_reviewed     = EXCLUDED.last_reviewed,
			next_review_due   = EXCLUDED.next_review_due,
			sm2_repetitions   = EXCLUDED.sm2_repetitions,
			sm2_interval      = EXCLUDED.sm2_interval,
			sm2_efactor       = EXCLUDED.sm2_efactor,
			mastered_at       = EXCLUDED.mastered_at,
			weakness_score    = EXCLUDED.weakness_score
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

func (s *PostgresStore) CreateSession(studentID string) (*Session, error) {
	id := newUUID()
	now := time.Now().UTC()
	_, err := s.db.Exec(
		"INSERT INTO sessions (id, student_id, started_at) VALUES ($1, $2, $3)",
		id, studentID, now.Format(time.RFC3339),
	)
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}
	return &Session{ID: id, StudentID: studentID, StartedAt: now}, nil
}

func (s *PostgresStore) GetSession(id string) (*Session, error) {
	row := s.db.QueryRow("SELECT id, student_id, started_at FROM sessions WHERE id = $1", id)
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

func (s *PostgresStore) GetActiveSession(sessionID string) (*ActiveSession, error) {
	row := s.db.QueryRow(`
		SELECT session_id, student_id, concept_id, concept_name, expected_answer,
		       attempt_id, question, explanation, diagram, is_review, answered,
		       last_concept_id, session_review, session_new, updated_at
		FROM active_sessions WHERE session_id = $1
	`, sessionID)
	var a ActiveSession
	var isReview, answered int
	var updatedAt string
	if err := row.Scan(&a.SessionID, &a.StudentID, &a.ConceptID, &a.ConceptName, &a.ExpectedAnswer,
		&a.AttemptID, &a.Question, &a.Explanation, &a.Diagram, &isReview, &answered,
		&a.LastConceptID, &a.SessionReview, &a.SessionNew, &updatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get active session: %w", err)
	}
	a.IsReview = isReview != 0
	a.Answered = answered != 0
	if updatedAt != "" {
		if t, err := time.Parse(time.RFC3339, updatedAt); err == nil {
			a.UpdatedAt = t
		}
	}
	return &a, nil
}

func (s *PostgresStore) UpsertActiveSession(a *ActiveSession) error {
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := s.db.Exec(`
		INSERT INTO active_sessions
			(session_id, student_id, concept_id, concept_name, expected_answer,
			 attempt_id, question, explanation, diagram, is_review, answered,
			 last_concept_id, session_review, session_new, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		ON CONFLICT(session_id) DO UPDATE SET
			student_id      = EXCLUDED.student_id,
			concept_id      = EXCLUDED.concept_id,
			concept_name    = EXCLUDED.concept_name,
			expected_answer = EXCLUDED.expected_answer,
			attempt_id      = EXCLUDED.attempt_id,
			question        = EXCLUDED.question,
			explanation     = EXCLUDED.explanation,
			diagram         = EXCLUDED.diagram,
			is_review       = EXCLUDED.is_review,
			answered        = EXCLUDED.answered,
			last_concept_id = EXCLUDED.last_concept_id,
			session_review  = EXCLUDED.session_review,
			session_new     = EXCLUDED.session_new,
			updated_at      = EXCLUDED.updated_at
	`, a.SessionID, a.StudentID, a.ConceptID, a.ConceptName, a.ExpectedAnswer,
		a.AttemptID, a.Question, a.Explanation, a.Diagram, boolToInt(a.IsReview), boolToInt(a.Answered),
		a.LastConceptID, a.SessionReview, a.SessionNew, now)
	if err != nil {
		return fmt.Errorf("upsert active session: %w", err)
	}
	return nil
}

func (s *PostgresStore) DeleteActiveSession(sessionID string) error {
	if _, err := s.db.Exec("DELETE FROM active_sessions WHERE session_id = $1", sessionID); err != nil {
		return fmt.Errorf("delete active session: %w", err)
	}
	return nil
}

// Attempts

func (s *PostgresStore) RecordAttempt(entry AttemptEntry) error {
	_, err := s.db.Exec(`
		INSERT INTO attempts
			(session_id, student_id, concept_id, answer, expected,
			 correct, elapsed_seconds, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
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

func (s *PostgresStore) GetAllAttempts() ([]AttemptEntry, error) {
	rows, err := s.db.Query(`
		SELECT session_id, student_id, concept_id, answer, expected,
		       correct, elapsed_seconds, timestamp
		FROM attempts
		ORDER BY student_id, concept_id, timestamp ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("get all attempts: %w", err)
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

func (s *PostgresStore) GetAttemptsForStudent(studentID string) ([]AttemptEntry, error) {
	rows, err := s.db.Query(`
		SELECT session_id, student_id, concept_id, answer, expected,
		       correct, elapsed_seconds, timestamp
		FROM attempts
		WHERE student_id = $1
		ORDER BY concept_id, timestamp ASC
	`, studentID)
	if err != nil {
		return nil, fmt.Errorf("get attempts for student: %w", err)
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

func (s *PostgresStore) GetSessionAttempts(studentID, sessionID string) ([]AttemptEntry, error) {
	rows, err := s.db.Query(`
		SELECT session_id, student_id, concept_id, answer, expected,
		       correct, elapsed_seconds, timestamp
		FROM attempts
		WHERE student_id = $1 AND session_id = $2
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

func (s *PostgresStore) GetQuestions(conceptID string, count int) ([]Question, error) {
	rows, err := s.db.Query(`
		SELECT id, concept_id, question, answer, explanation, source, difficulty
		FROM questions
		WHERE concept_id = $1
		ORDER BY RANDOM()
		LIMIT $2
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

func (s *PostgresStore) GetQuestionCount(conceptID string) (int, error) {
	row := s.db.QueryRow("SELECT COUNT(*) FROM questions WHERE concept_id = $1", conceptID)
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("get question count: %w", err)
	}
	return count, nil
}

func (s *PostgresStore) ImportQuestions(qs []Question) error {
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
		INSERT INTO questions (concept_id, question, answer, explanation, source, difficulty)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT DO NOTHING
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

func (s *PostgresStore) PurgeGeneratedQuestions(conceptIDs map[string]bool) (int64, error) {
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
		placeholders := make([]string, len(chunk))
		args := make([]interface{}, len(chunk))
		for i, id := range chunk {
			placeholders[i] = fmt.Sprintf("$%d", i+1)
			args[i] = id
		}
		query := fmt.Sprintf("DELETE FROM questions WHERE concept_id IN (%s)", strings.Join(placeholders, ","))
		res, err := s.db.Exec(query, args...)
		if err != nil {
			return deleted, fmt.Errorf("purge questions: %w", err)
		}
		n, _ := res.RowsAffected()
		deleted += n
	}
	return deleted, nil
}

// Leaderboard

func (s *PostgresStore) GetWeeklyLeaderboard() ([]LeaderboardRow, error) {
	monday := weekStart(time.Now().UTC())
	rows, err := s.db.Query(`
		SELECT s.id, s.name,
			COALESCE((SELECT COUNT(*) FROM concept_progress
			          WHERE student_id = s.id AND status = 'MASTERED'), 0),
			COALESCE((SELECT COUNT(*) FROM concept_progress
			          WHERE student_id = s.id AND status = 'MASTERED'
			          AND mastered_at >= $1), 0)
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

func (s *PostgresStore) GetLeagueStandings() ([]LeagueMember, error) {
	monday := weekStart(time.Now().UTC())
	rows, err := s.db.Query(`
		SELECT s.id, s.name, COALESCE(s.league, 'bronze'), s.league_moved,
			COALESCE((SELECT COUNT(*) FROM concept_progress
			          WHERE student_id = s.id AND status = 'MASTERED'), 0),
			COALESCE((SELECT COUNT(*) FROM concept_progress
			          WHERE student_id = s.id AND status = 'MASTERED'
			          AND mastered_at >= $1), 0)
		FROM students s
	`, monday.Format(time.RFC3339))
	if err != nil {
		return nil, fmt.Errorf("get league standings: %w", err)
	}
	defer rows.Close()

	var out []LeagueMember
	for rows.Next() {
		var m LeagueMember
		var moved int
		if err := rows.Scan(&m.StudentID, &m.Name, &m.Tier, &moved, &m.TotalMastered, &m.WeeklyMastered); err != nil {
			return nil, fmt.Errorf("scan league member: %w", err)
		}
		m.Moved = moved
		if m.Tier == "" {
			m.Tier = "bronze"
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func (s *PostgresStore) SetLeague(studentID, tier, week string, moved int) error {
	_, err := s.db.Exec(
		"UPDATE students SET league = $1, league_week = $2, league_moved = $3 WHERE id = $4",
		tier, week, moved, studentID,
	)
	if err != nil {
		return fmt.Errorf("set league: %w", err)
	}
	return nil
}

// DailyActivity

func (s *PostgresStore) GetDailyActivity(studentID string, days int) ([]DailyActivity, error) {
	// For postgres, use now() - interval
	rows, err := s.db.Query(`
		SELECT date(timestamp)::text as day,
		       COUNT(*) as total,
		       COALESCE(SUM(correct), 0) as correct_count,
		       COALESCE(STRING_AGG(DISTINCT concept_id, ','), '') as concepts
		FROM attempts
		WHERE student_id = $1 AND timestamp >= now() - ($2 || ' days')::interval
		GROUP BY day
		ORDER BY day
	`, studentID, fmt.Sprintf("%d", days))
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

// TopicSpeed

func (s *PostgresStore) GetTopicSpeed(studentID, conceptID string) (*TopicSpeed, error) {
	row := s.db.QueryRow(`SELECT student_id, concept_id, efactor, interval, repetitions, learning_speed FROM student_topic_speed WHERE student_id = $1 AND concept_id = $2`, studentID, conceptID)
	var ts TopicSpeed
	if err := row.Scan(&ts.StudentID, &ts.ConceptID, &ts.EFactor, &ts.Interval, &ts.Repetitions, &ts.LearningSpeed); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get topic speed: %w", err)
	}
	return &ts, nil
}

func (s *PostgresStore) GetAllTopicSpeeds(studentID string) (map[string]*TopicSpeed, error) {
	rows, err := s.db.Query(`SELECT student_id, concept_id, efactor, interval, repetitions, learning_speed FROM student_topic_speed WHERE student_id = $1`, studentID)
	if err != nil {
		return nil, fmt.Errorf("get all topic speeds: %w", err)
	}
	defer rows.Close()
	out := make(map[string]*TopicSpeed)
	for rows.Next() {
		var ts TopicSpeed
		if err := rows.Scan(&ts.StudentID, &ts.ConceptID, &ts.EFactor, &ts.Interval, &ts.Repetitions, &ts.LearningSpeed); err != nil {
			return nil, err
		}
		out[ts.ConceptID] = &ts
	}
	return out, rows.Err()
}

func (s *PostgresStore) UpsertTopicSpeed(ts *TopicSpeed) error {
	if ts.LearningSpeed < 0.5 {
		ts.LearningSpeed = 0.5
	}
	if ts.LearningSpeed > 2.0 {
		ts.LearningSpeed = 2.0
	}
	if ts.EFactor == 0 {
		ts.EFactor = 2.5
	}
	_, err := s.db.Exec(`
		INSERT INTO student_topic_speed (student_id, concept_id, efactor, interval, repetitions, learning_speed, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, now()::text)
		ON CONFLICT(student_id, concept_id) DO UPDATE SET
			efactor = EXCLUDED.efactor,
			interval = EXCLUDED.interval,
			repetitions = EXCLUDED.repetitions,
			learning_speed = EXCLUDED.learning_speed,
			updated_at = now()::text
	`, ts.StudentID, ts.ConceptID, ts.EFactor, ts.Interval, ts.Repetitions, ts.LearningSpeed)
	if err != nil {
		return fmt.Errorf("upsert topic speed: %w", err)
	}
	return nil
}
