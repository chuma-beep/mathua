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
    diagnostic_completed_at TEXT NOT NULL DEFAULT '',
    daily_xp_goal        INTEGER NOT NULL DEFAULT 30,
    settings             TEXT NOT NULL DEFAULT '{}',
    league               TEXT NOT NULL DEFAULT 'bronze',
    league_week          TEXT NOT NULL DEFAULT '',
    league_moved         INTEGER NOT NULL DEFAULT 0,
    share_token          TEXT NOT NULL DEFAULT '',
    profile_id           TEXT NOT NULL DEFAULT '',
    friend_code          TEXT NOT NULL DEFAULT '',
    email                TEXT NOT NULL DEFAULT '',
    google_id            TEXT NOT NULL DEFAULT '',
    avatar_url           TEXT NOT NULL DEFAULT '',
    email_verified       INTEGER NOT NULL DEFAULT 0
);

-- The XP ledger. xp_total and xp_today are running counters for the profile;
-- this table is the history behind them, and the only honest source for "XP
-- earned this week". Without it a weekly board can only count mastered
-- concepts, which is not effort. UNIQUE(student_id, event_key) is the
-- duplicate-award guard: the key names the attempt, and the insert is part of
-- the same transaction that moves the counters.
CREATE TABLE IF NOT EXISTS xp_events (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    student_id TEXT NOT NULL,
    amount     INTEGER NOT NULL,
    reason     TEXT NOT NULL DEFAULT '',
    event_key  TEXT NOT NULL,
    created_at INTEGER NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_xp_events_event ON xp_events(student_id, event_key);
CREATE INDEX IF NOT EXISTS idx_xp_events_week ON xp_events(created_at);
CREATE INDEX IF NOT EXISTS idx_xp_events_student_week ON xp_events(student_id, created_at);

-- Weekly league membership. A row is a learner's seat in one league for one
-- week, written by the transition and never rewritten while the week is open:
-- the seat is what makes "my league" stable even as XP moves the standings.
-- league_weeks is the idempotency marker, the transition inserts it first in
-- the same transaction, so a retry after a crash sees it and stops instead of
-- promoting the same cohort twice.
CREATE TABLE IF NOT EXISTS league_members (
    week       TEXT    NOT NULL,
    student_id TEXT    NOT NULL,
    tier       TEXT    NOT NULL,
    league_no  INTEGER NOT NULL,
    xp         INTEGER NOT NULL DEFAULT 0,
    rank       INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (week, student_id)
);

CREATE INDEX IF NOT EXISTS idx_league_members_cohort ON league_members(week, tier, league_no, xp DESC);

-- One row per answered diagnostic question, keyed by session and concept.
-- The diagnostic has no attempt token in its wire format, so the pair is the
-- identity; replaying it returns the same verdict instead of grading twice.
CREATE TABLE IF NOT EXISTS diagnostic_answers (
    attempt_id TEXT PRIMARY KEY,
    response   TEXT NOT NULL,
    created_at INTEGER NOT NULL
);

-- One row per accepted study answer, keyed by the client's submission id.
-- The row is inserted in the same transaction as the attempt, progress, topic
-- speed and XP award, so a retry either replays the stored response or sees no
-- row at all. request_hash rejects the same key carrying a different answer.
CREATE TABLE IF NOT EXISTS answer_submissions (
    student_id    TEXT    NOT NULL,
    operation     TEXT    NOT NULL,
    submission_id TEXT    NOT NULL,
    request_hash  TEXT    NOT NULL,
    response_json TEXT    NOT NULL DEFAULT '',
    attempt_id    INTEGER NOT NULL DEFAULT 0,
    created_at    INTEGER NOT NULL,
    PRIMARY KEY (student_id, operation, submission_id),
    FOREIGN KEY (student_id) REFERENCES students(id)
);

CREATE TABLE IF NOT EXISTS study_answer_effects (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    student_id    TEXT    NOT NULL,
    submission_id TEXT    NOT NULL,
    effect        TEXT    NOT NULL,
    payload_json  TEXT    NOT NULL DEFAULT '',
    created_at    INTEGER NOT NULL,
    attempts      INTEGER NOT NULL DEFAULT 0,
    next_attempt_at INTEGER NOT NULL,
    claimed_until INTEGER NOT NULL DEFAULT 0,
    processed_at  INTEGER NOT NULL DEFAULT 0,
    last_error    TEXT    NOT NULL DEFAULT '',
    UNIQUE (student_id, submission_id, effect),
    FOREIGN KEY (student_id) REFERENCES students(id)
);

CREATE INDEX IF NOT EXISTS idx_study_answer_effects_pending
    ON study_answer_effects (processed_at, next_attempt_at, claimed_until);

CREATE TABLE IF NOT EXISTS league_weeks (
    week      TEXT PRIMARY KEY,
    closed_at INTEGER NOT NULL
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
    placement_seeded INTEGER NOT NULL DEFAULT 0,
    is_provisional   INTEGER NOT NULL DEFAULT 0,
    fsrs_stability   REAL,
    fsrs_difficulty  REAL,
    fsrs_last_review TEXT,
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
    question        TEXT    NOT NULL DEFAULT '',
    source          TEXT    NOT NULL DEFAULT '',
    explanation     TEXT    NOT NULL DEFAULT '',
    -- Difficulty the generator served this question at, 0.3-1.0. NULL means the attempt
    -- predates the column, or the question was not generated at all, so difficulty is
    -- *unknown* rather than zero. Evidence has to be able to tell those apart, so this is
    -- nullable and never backfilled with a plausible-looking number.
    difficulty      REAL,
    FOREIGN KEY (session_id) REFERENCES sessions(id),
    FOREIGN KEY (student_id) REFERENCES students(id)
);

CREATE INDEX IF NOT EXISTS idx_progress_student  ON concept_progress(student_id);
CREATE INDEX IF NOT EXISTS idx_progress_mastered ON concept_progress(student_id, status, mastered_at);
CREATE INDEX IF NOT EXISTS idx_sessions_student  ON sessions(student_id);
CREATE INDEX IF NOT EXISTS idx_attempts_session  ON attempts(session_id);
CREATE INDEX IF NOT EXISTS idx_attempts_student  ON attempts(student_id, timestamp);
CREATE INDEX IF NOT EXISTS idx_attempts_cover    ON attempts(student_id, concept_id, timestamp);

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
    difficulty      REAL,
    kp_ordinal      INTEGER NOT NULL DEFAULT 0,
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

-- Batch 1 (Quiz 150 XP gate): one row per completed mastery-check quiz.
-- xp_total snapshots lifetime XP so xp_since_quiz = current - last.
CREATE TABLE IF NOT EXISTS quiz_completions (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    student_id   TEXT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    completed_at TEXT NOT NULL DEFAULT (datetime('now')),
    xp_total     INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_quiz_completions_student ON quiz_completions(student_id, completed_at);

-- Timed quizzes are separate from the 50 XP mastery check. baseline_at is
-- the rollout boundary: only positive XP ledger events at or after it count
-- toward the 150 XP gate.
CREATE TABLE IF NOT EXISTS timed_quiz_baselines (
    student_id   TEXT PRIMARY KEY REFERENCES students(id) ON DELETE CASCADE,
    baseline_xp  INTEGER NOT NULL DEFAULT 0,
    baseline_at  INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS timed_quiz_runs (
    id                TEXT PRIMARY KEY,
    student_id        TEXT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    started_at        INTEGER NOT NULL,
    deadline_at       INTEGER NOT NULL,
    completed_at      INTEGER NOT NULL DEFAULT 0,
    status            TEXT NOT NULL DEFAULT 'active',
    question_count    INTEGER NOT NULL DEFAULT 0,
    correct_count     INTEGER NOT NULL DEFAULT 0,
    unanswered_count  INTEGER NOT NULL DEFAULT 0,
    xp_earned         INTEGER NOT NULL DEFAULT 0,
    created_at        INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_timed_quiz_runs_student
    ON timed_quiz_runs(student_id, status, deadline_at);
CREATE UNIQUE INDEX IF NOT EXISTS uidx_timed_quiz_active
    ON timed_quiz_runs(student_id) WHERE status = 'active';

CREATE TABLE IF NOT EXISTS timed_quiz_items (
    id               TEXT PRIMARY KEY,
    run_id           TEXT NOT NULL REFERENCES timed_quiz_runs(id) ON DELETE CASCADE,
    concept_id       TEXT NOT NULL,
    position         INTEGER NOT NULL,
    question         TEXT NOT NULL,
    expected_answer  TEXT NOT NULL,
    explanation      TEXT NOT NULL DEFAULT '',
    grading_type     TEXT NOT NULL DEFAULT '',
    difficulty       REAL,
    served_at        INTEGER NOT NULL,
    answered_at      INTEGER NOT NULL DEFAULT 0,
    elapsed_seconds  REAL NOT NULL DEFAULT 0,
    answer           TEXT NOT NULL DEFAULT '',
    correct          INTEGER NOT NULL DEFAULT 0,
    status           TEXT NOT NULL DEFAULT 'served',
    submission_id    TEXT NOT NULL DEFAULT '',
    xp_earned        INTEGER NOT NULL DEFAULT 0,
    UNIQUE(run_id, concept_id)
);
CREATE INDEX IF NOT EXISTS idx_timed_quiz_items_run
    ON timed_quiz_items(run_id, position);

-- One row per spaced-review attempt. attempt_id is the idempotency key.
CREATE TABLE IF NOT EXISTS review_events (
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    student_id        TEXT NOT NULL,
    concept_id        TEXT NOT NULL,
    attempt_id        TEXT NOT NULL UNIQUE,
    rating            INTEGER NOT NULL,
    correct           INTEGER NOT NULL,
    elapsed_seconds   REAL NOT NULL,
    reviewed_at       INTEGER NOT NULL,
    stability_before  REAL,
    difficulty_before REAL,
    due_before        TEXT,
    stability_after   REAL,
    difficulty_after  REAL,
    due_after         TEXT,
    created_at        INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_review_events_student
    ON review_events(student_id, reviewed_at);

CREATE TABLE IF NOT EXISTS avatar_images (
    student_id   TEXT PRIMARY KEY REFERENCES students(id) ON DELETE CASCADE,
    content_type TEXT NOT NULL,
    bytes        BLOB NOT NULL,
    updated_at   TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS password_resets (
    token_hash TEXT PRIMARY KEY,
    student_id TEXT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    expires_at TEXT NOT NULL,
    used       INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS identities (
    provider       TEXT NOT NULL,
    provider_id    TEXT NOT NULL,
    student_id     TEXT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    email          TEXT NOT NULL DEFAULT '',
    email_verified INTEGER NOT NULL DEFAULT 0,
    linked_at      TEXT NOT NULL DEFAULT (datetime('now')),
    PRIMARY KEY (provider, provider_id)
);
CREATE INDEX IF NOT EXISTS idx_identities_student ON identities(student_id);

CREATE TABLE IF NOT EXISTS email_verifications (
    token_hash TEXT PRIMARY KEY,
    student_id TEXT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    expires_at TEXT NOT NULL,
    used       INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS link_tokens (
    token_hash TEXT PRIMARY KEY,
    student_id TEXT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    expires_at TEXT NOT NULL,
    used       INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS question_reports (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    reporter_id TEXT NOT NULL DEFAULT '',
    concept_id  TEXT NOT NULL DEFAULT '',
    kind        TEXT NOT NULL DEFAULT 'question',
    question    TEXT NOT NULL DEFAULT '',
    expected    TEXT NOT NULL DEFAULT '',
    explanation TEXT NOT NULL DEFAULT '',
    lesson_id   TEXT NOT NULL DEFAULT '',
    source      TEXT NOT NULL DEFAULT '',
    session_id  TEXT NOT NULL DEFAULT '',
    attempt_id  TEXT NOT NULL DEFAULT '',
    answer_given     TEXT NOT NULL DEFAULT '',
    grading_type     TEXT NOT NULL DEFAULT '',
    answer_hint      TEXT NOT NULL DEFAULT '',
    answer_input_mode TEXT NOT NULL DEFAULT '',
    reason      TEXT NOT NULL DEFAULT 'other',
    detail      TEXT NOT NULL DEFAULT '',
    status      TEXT NOT NULL DEFAULT 'open',
    resolution  TEXT NOT NULL DEFAULT '',
    resolved_by TEXT NOT NULL DEFAULT '',
    resolved_at TEXT NOT NULL DEFAULT '',
    created_at  TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_reports_status  ON question_reports(status, created_at);
CREATE INDEX IF NOT EXISTS idx_reports_concept ON question_reports(concept_id);

-- The last time an account opened its own report history. A resolved report
-- newer than this is unread; one shared timestamp is enough because the page
-- lists every report at once. Deleting the account removes the marker.
CREATE TABLE IF NOT EXISTS report_views (
    student_id TEXT PRIMARY KEY REFERENCES students(id) ON DELETE CASCADE,
    seen_at    TEXT NOT NULL DEFAULT '',
    cleared_at TEXT NOT NULL DEFAULT ''
);

-- One optional screenshot per report. The bytes live beside the moderation
-- record so deleting an account removes the image with the report that named
-- it, and no upload can outlive the queue entry that explains it.
CREATE TABLE IF NOT EXISTS report_images (
    report_id    INTEGER PRIMARY KEY REFERENCES question_reports(id) ON DELETE CASCADE,
    content_type TEXT NOT NULL,
    bytes        BLOB NOT NULL,
    created_at   TEXT NOT NULL DEFAULT (datetime('now'))
);

-- Live diagnostic runs. The adaptive model is a graph of beliefs over the whole
-- path, so it cannot be rebuilt from the attempts table; it is stored as the
-- diagnostic package's own encoded snapshot. Before this table the run lived
-- only in memory and every deploy answered the next question with
-- "diagnostic session not found". student_id is indexed because ResetProgress
-- and account deletion must evict the rows, not just the map entries.
CREATE TABLE IF NOT EXISTS diagnostic_sessions (
    id         TEXT PRIMARY KEY,
    student_id TEXT NOT NULL DEFAULT '',
    data       BLOB NOT NULL,
    updated_at INTEGER NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_diagnostic_sessions_student ON diagnostic_sessions(student_id);
CREATE INDEX IF NOT EXISTS idx_diagnostic_sessions_updated ON diagnostic_sessions(updated_at);

-- Staff invitations. The raw token is never stored, only its SHA-256 hash, so a database read
-- cannot yield a usable credential. A row is never deleted: a revoked or expired invitation is
-- history, and deleting it would erase the record of who offered production access and to whom.
CREATE TABLE IF NOT EXISTS admin_invitations (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    token_hash  TEXT NOT NULL,
    email       TEXT NOT NULL DEFAULT '',
    role        TEXT NOT NULL,
    invited_by  TEXT NOT NULL DEFAULT '',
    created_at  TEXT NOT NULL DEFAULT (datetime('now')),
    expires_at  TEXT NOT NULL,
    accepted_at TEXT NOT NULL DEFAULT '',
    accepted_by TEXT NOT NULL DEFAULT '',
    revoked_at  TEXT NOT NULL DEFAULT '',
    revoked_by  TEXT NOT NULL DEFAULT ''
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_admin_invitations_token ON admin_invitations(token_hash);
CREATE INDEX IF NOT EXISTS idx_admin_invitations_email ON admin_invitations(email);

-- Product announcements. They are global rows, so ResetProgress and account
-- deletion never remove the announcement itself. Reads are per learner and are
-- removed with the account. Body is Markdown rendered through the lesson
-- sanitizer on the client.
CREATE TABLE IF NOT EXISTS announcements (
    id           TEXT PRIMARY KEY,
    title        TEXT NOT NULL,
    body         TEXT NOT NULL,
    published    INTEGER NOT NULL DEFAULT 0,
    created_at   INTEGER NOT NULL,
    published_at INTEGER NOT NULL DEFAULT 0,
    created_by   TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_announcements_published ON announcements(published, published_at);

CREATE TABLE IF NOT EXISTS announcement_reads (
    student_id      TEXT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    announcement_id TEXT NOT NULL REFERENCES announcements(id) ON DELETE CASCADE,
    read_at         INTEGER NOT NULL,
    PRIMARY KEY (student_id, announcement_id)
);

CREATE INDEX IF NOT EXISTS idx_announcement_reads_announcement ON announcement_reads(announcement_id);

-- Announcement clears hide the current feed from one learner without deleting
-- the global announcement or another learner's feed.
CREATE TABLE IF NOT EXISTS announcement_views (
    student_id TEXT PRIMARY KEY REFERENCES students(id) ON DELETE CASCADE,
    cleared_at INTEGER NOT NULL
);

-- Administration. The role column itself is added by the migration list below rather than
-- here: this DDL runs on every boot, and ALTER TABLE ADD COLUMN is not idempotent. The
-- default is 'student', so every existing row is a learner until an operator promotes it.
--
-- One row per privileged mutation. actor_id is deliberately not a foreign key: a trail that
-- is deleted along with the account it records is not a trail. before_json/after_json hold the
-- non-sensitive AdminUser projection only, a password hash cannot reach them because the
-- projection cannot express one.
CREATE TABLE IF NOT EXISTS admin_audit (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    actor_id    TEXT NOT NULL DEFAULT '',
    action      TEXT NOT NULL,
    entity_type TEXT NOT NULL DEFAULT '',
    entity_id   TEXT NOT NULL DEFAULT '',
    before_json TEXT NOT NULL DEFAULT '',
    after_json  TEXT NOT NULL DEFAULT '',
    created_at  TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_admin_audit_created ON admin_audit(created_at);
CREATE INDEX IF NOT EXISTS idx_admin_audit_actor   ON admin_audit(actor_id);
CREATE INDEX IF NOT EXISTS idx_admin_audit_entity   ON admin_audit(entity_type, entity_id);

-- Fix 6: durable server-side sessions (study anti-cheat expected answers,
-- admin triage logins). Short-lived diag/quiz sessions stay in memory by
-- design (capability UUIDs, 1h janitor; a restart just means a retake).
CREATE TABLE IF NOT EXISTS server_sessions (
    kind       TEXT NOT NULL,
    key        TEXT NOT NULL,
    value      TEXT NOT NULL DEFAULT '',
    expires_at TEXT NOT NULL DEFAULT '',
    updated_at TEXT NOT NULL DEFAULT (datetime('now')),
    PRIMARY KEY (kind, key)
);

CREATE INDEX IF NOT EXISTS idx_server_sessions_expiry ON server_sessions(expires_at);

-- Per-identity daily AI request counters.
--
-- A counter is a day-scoped row keyed by an opaque usage_key, not a column on
-- students, for two reasons. A guest has no students row, and the guest's
-- identity is an IP rather than an id, so the key can not be a foreign key. And
-- a row per day lets the daily budget roll over without a reset job.
--
-- ResetProgress deliberately does not touch this table. A learner who could
-- refill their AI budget by resetting progress would have no budget at all;
-- DeleteAccount does clear it, because the account is gone.
CREATE TABLE IF NOT EXISTS ai_usage (
    usage_key TEXT NOT NULL,
    day       TEXT NOT NULL,
    requests  INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (usage_key, day)
);

CREATE INDEX IF NOT EXISTS idx_ai_usage_day ON ai_usage(day);

-- Fix 7: version tracking for one-time data backfills. DDL stays
-- idempotent (IF NOT EXISTS); data fixes below run exactly once.
-- NOTE: no UNIQUE constraints are added here on purpose, legacy username
-- colliding groups (same lower form twice) still exist and are renamed
-- manually; a UNIQUE would fail migration on real databases.
-- FIRe-inspired implicit credit, shadow only. These tables record what an
-- implicit repetition would have been worth; no query in the scheduler reads
-- them, so production due dates and review order are unchanged.
CREATE TABLE IF NOT EXISTS implicit_credit_targets (
    student_id        TEXT    NOT NULL,
    target_concept_id TEXT    NOT NULL,
    credit_epoch      INTEGER NOT NULL DEFAULT 1,
    created_at        INTEGER NOT NULL,
    PRIMARY KEY (student_id, target_concept_id)
);

CREATE TABLE IF NOT EXISTS implicit_credit_allocations (
    student_id        TEXT    NOT NULL,
    source_attempt_id TEXT    NOT NULL,
    raw_units         INTEGER NOT NULL,
    allocated_units   INTEGER NOT NULL,
    target_count      INTEGER NOT NULL,
    dropped_count     INTEGER NOT NULL,
    created_at        INTEGER NOT NULL,
    PRIMARY KEY (student_id, source_attempt_id)
);

CREATE TABLE IF NOT EXISTS implicit_credit_events (
    student_id        TEXT    NOT NULL,
    source_attempt_id TEXT    NOT NULL,
    target_concept_id TEXT    NOT NULL,
    credit_epoch      INTEGER NOT NULL,
    scaled_units      INTEGER NOT NULL,
    retained_units    INTEGER NOT NULL,
    discarded_units   INTEGER NOT NULL,
    discard_reason    TEXT    NOT NULL DEFAULT '',
    voided_at         INTEGER,
    voided_by         TEXT    NOT NULL DEFAULT '',
    void_reason       TEXT    NOT NULL DEFAULT '',
    created_at        INTEGER NOT NULL,
    PRIMARY KEY (student_id, source_attempt_id, target_concept_id)
);

CREATE INDEX IF NOT EXISTS idx_credit_events_target
    ON implicit_credit_events(student_id, target_concept_id, credit_epoch);

CREATE TABLE IF NOT EXISTS implicit_credit_ordinals (
    student_id          TEXT    NOT NULL,
    target_concept_id   TEXT    NOT NULL,
    credit_epoch        INTEGER NOT NULL,
    ordinal             INTEGER NOT NULL,
    state               TEXT    NOT NULL CHECK (state IN ('banked','projected','stale')),
    banked_at           INTEGER NOT NULL,
    projected_at        INTEGER,
    baseline_due        TEXT,
    counterfactual_due  TEXT,
    attempts            INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (student_id, target_concept_id, credit_epoch, ordinal)
);

-- Active FIRe applications. One row per ordinal actually applied to a real
-- due date. The unique keys make a retry or a second worker replay, never a
-- second scheduling transition.
CREATE TABLE IF NOT EXISTS implicit_credit_applications (
    student_id         TEXT    NOT NULL,
    target_concept_id  TEXT    NOT NULL,
    credit_epoch       INTEGER NOT NULL,
    ordinal            INTEGER NOT NULL,
    source_concept_id  TEXT    NOT NULL DEFAULT '',
    source_attempt_id  TEXT    NOT NULL,
    applied_at         INTEGER NOT NULL,
    actual_due         TEXT    NOT NULL DEFAULT '',
    proposed_due       TEXT    NOT NULL DEFAULT '',
    stability_before   REAL,
    stability_after    REAL,
    difficulty_before  REAL,
    difficulty_after   REAL,
    reviews_deferred   INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (student_id, target_concept_id, credit_epoch, ordinal),
    UNIQUE (student_id, source_attempt_id, target_concept_id, credit_epoch)
);

CREATE INDEX IF NOT EXISTS idx_credit_applications_attempt
    ON implicit_credit_applications(student_id, source_attempt_id);

-- Durable activation state. Environment configuration may request active
-- mode, but this row is the server-side acknowledgement an operator must
-- write before any credit can be applied.
CREATE TABLE IF NOT EXISTS fire_activations (
    id              TEXT PRIMARY KEY,
    mode            TEXT NOT NULL,
    cohort          TEXT NOT NULL,
    acknowledged_by TEXT NOT NULL,
    acknowledged_at INTEGER NOT NULL,
    updated_at      INTEGER NOT NULL,
    revision        INTEGER NOT NULL
);

-- Targeted remedial queue: 2-miss halt remedial prerequisites survive Fly.io restarts.
CREATE TABLE IF NOT EXISTS student_remedial_queue (
    student_id TEXT NOT NULL,
    concept_id TEXT NOT NULL,
    priority   INTEGER NOT NULL,
    created_at INTEGER NOT NULL,
    PRIMARY KEY (student_id, concept_id),
    FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_remedial_queue_student ON student_remedial_queue(student_id, priority, created_at);

-- Mutual friendships. The pair is stored in lexicographic order so a request
-- in either direction collides on the same primary key instead of creating a
-- second row. requested_by preserves who asked; a reciprocal request promotes
-- the pending row to accepted.
CREATE TABLE IF NOT EXISTS friendships (
    student_a    TEXT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    student_b    TEXT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    status       TEXT NOT NULL CHECK (status IN ('pending', 'accepted')),
    requested_by TEXT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    created_at   INTEGER NOT NULL,
    responded_at INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (student_a, student_b),
    CHECK (student_a < student_b),
    CHECK (requested_by = student_a OR requested_by = student_b)
);

CREATE INDEX IF NOT EXISTS idx_friendships_a_status ON friendships(student_a, status);
CREATE INDEX IF NOT EXISTS idx_friendships_b_status ON friendships(student_b, status);

CREATE TABLE IF NOT EXISTS friend_blocks (
    blocker_id TEXT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    blocked_id TEXT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    created_at INTEGER NOT NULL,
    PRIMARY KEY (blocker_id, blocked_id),
    CHECK (blocker_id <> blocked_id)
);

CREATE INDEX IF NOT EXISTS idx_friend_blocks_blocked ON friend_blocks(blocked_id);

CREATE TABLE IF NOT EXISTS schema_migrations (
    version    INTEGER PRIMARY KEY,
    applied_at TEXT NOT NULL DEFAULT (datetime('now'))
);
`
