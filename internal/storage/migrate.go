package storage

const schema = `
CREATE TABLE IF NOT EXISTS students (
    id            TEXT PRIMARY KEY,
    name          TEXT NOT NULL,
    username      TEXT,
    password_hash TEXT,
    created_at    TEXT NOT NULL DEFAULT (datetime('now'))
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
CREATE INDEX IF NOT EXISTS idx_sessions_student  ON sessions(student_id);
CREATE INDEX IF NOT EXISTS idx_attempts_session  ON attempts(session_id);
CREATE INDEX IF NOT EXISTS idx_attempts_student  ON attempts(student_id, timestamp);
`
