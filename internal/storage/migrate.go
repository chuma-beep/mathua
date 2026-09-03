package storage

const schema = `
CREATE TABLE IF NOT EXISTS students (
    id                   TEXT PRIMARY KEY,
    name                 TEXT NOT NULL,
    username             TEXT,
    password_hash        TEXT,
    created_at           TEXT NOT NULL DEFAULT (datetime('now')),
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
    email                TEXT NOT NULL DEFAULT '',
    google_id            TEXT NOT NULL DEFAULT '',
    avatar_url           TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS concept_progress (
    student_id       TEXT    NOT NULL,
    concept_id       TEXT    NOT NULL,
    status           TEXT    NOT NULL DEFAULT 'UNSEEN',
    streak           INTEGER NOT NULL DEFAULT 0,
    best_streak      INTEGER NOT NULL DEFAULT 0,
    avg_response_time REAL   NOT NULL DEFAULT 0,
    attempts         INTEGER NOT NULL DEFAULT 0,
    last_attempted   TEXT,
    last_reviewed    TEXT,
    next_review_due  TEXT,
    sm2_repetitions  INTEGER NOT NULL DEFAULT 0,
    sm2_interval     INTEGER NOT NULL DEFAULT 0,
    sm2_efactor      REAL    NOT NULL DEFAULT 2.5,
    mastered_at      TEXT,
    weakness_score   REAL    NOT NULL DEFAULT 0,
    PRIMARY KEY (student_id, concept_id),
    FOREIGN KEY (student_id) REFERENCES students(id)
);

CREATE TABLE IF NOT EXISTS sessions (
    id         TEXT PRIMARY KEY,
    student_id TEXT NOT NULL,
    started_at TEXT NOT NULL DEFAULT (datetime('now')),
    FOREIGN KEY (student_id) REFERENCES students(id)
);

CREATE TABLE IF NOT EXISTS attempts (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id      TEXT    NOT NULL,
    student_id      TEXT    NOT NULL,
    concept_id      TEXT    NOT NULL,
    answer          TEXT    NOT NULL,
    expected        TEXT    NOT NULL,
    correct         INTEGER NOT NULL DEFAULT 0,
    elapsed_seconds REAL    NOT NULL DEFAULT 0,
    timestamp       TEXT    NOT NULL DEFAULT (datetime('now')),
    FOREIGN KEY (session_id) REFERENCES sessions(id),
    FOREIGN KEY (student_id) REFERENCES students(id)
);

CREATE INDEX IF NOT EXISTS idx_progress_student  ON concept_progress(student_id);
CREATE INDEX IF NOT EXISTS idx_progress_mastered ON concept_progress(student_id, status, mastered_at);
CREATE INDEX IF NOT EXISTS idx_sessions_student  ON sessions(student_id);
CREATE INDEX IF NOT EXISTS idx_attempts_session  ON attempts(session_id);
CREATE INDEX IF NOT EXISTS idx_attempts_student  ON attempts(student_id, timestamp);
CREATE INDEX IF NOT EXISTS idx_attempts_cover    ON attempts(student_id, concept_id, timestamp);
CREATE INDEX IF NOT EXISTS idx_students_share    ON students(share_token);
CREATE INDEX IF NOT EXISTS idx_students_email     ON students(email);
CREATE INDEX IF NOT EXISTS idx_students_google_id ON students(google_id);

CREATE TABLE IF NOT EXISTS questions (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    concept_id TEXT    NOT NULL,
    question   TEXT    NOT NULL,
    answer     TEXT    NOT NULL,
    explanation TEXT   NOT NULL DEFAULT '',
    source     TEXT    NOT NULL DEFAULT '',
    difficulty REAL   NOT NULL DEFAULT 0.5
);

CREATE INDEX IF NOT EXISTS idx_questions_concept ON questions(concept_id);

CREATE TABLE IF NOT EXISTS active_sessions (
    session_id      TEXT PRIMARY KEY REFERENCES sessions(id) ON DELETE CASCADE,
    student_id      TEXT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
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
    updated_at      TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_active_sessions_student ON active_sessions(student_id);
CREATE INDEX IF NOT EXISTS idx_active_sessions_attempt ON active_sessions(attempt_id);

CREATE TABLE IF NOT EXISTS student_topic_speed (
    student_id    TEXT NOT NULL,
    concept_id    TEXT NOT NULL,
    efactor       REAL NOT NULL DEFAULT 2.5,
    interval      INTEGER NOT NULL DEFAULT 0,
    repetitions   INTEGER NOT NULL DEFAULT 0,
    learning_speed REAL NOT NULL DEFAULT 1.0,
    updated_at    TEXT NOT NULL DEFAULT (datetime('now')),
    PRIMARY KEY (student_id, concept_id),
    FOREIGN KEY (student_id) REFERENCES students(id)
);
CREATE INDEX IF NOT EXISTS idx_topic_speed_student ON student_topic_speed(student_id);
`
