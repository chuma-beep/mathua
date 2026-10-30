package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/chuma-beep/mathua/internal/ai"
	"github.com/chuma-beep/mathua/internal/fsrs"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresStore struct {
	db *sql.DB

	// creditProbe runs inside a projection transaction, immediately before the
	// real due date is re-read. It exists only for tests: the conflict guard is
	// otherwise unreachable because the projection holds the concept_progress
	// row FOR UPDATE, so no competing writer can move the due date
	// mid-projection. Per-store rather than package-level, so concurrent tests
	// cannot observe each other's probe. Nil in production, and nothing outside
	// this package can set it.
	creditProbe func(exec func(query string, args ...any) error)
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
    diagnostic_completed_at TEXT NOT NULL DEFAULT '',
    daily_xp_goal        INTEGER NOT NULL DEFAULT 30,
    settings             TEXT NOT NULL DEFAULT '{}',
    league               TEXT NOT NULL DEFAULT 'bronze',
    league_week          TEXT NOT NULL DEFAULT '',
    league_moved         INTEGER NOT NULL DEFAULT 0,
    share_token          TEXT NOT NULL DEFAULT '',
    profile_id           TEXT NOT NULL DEFAULT '',
    friend_code          TEXT NOT NULL DEFAULT '',
    course_id            TEXT,
    email                TEXT NOT NULL DEFAULT '',
    google_id            TEXT NOT NULL DEFAULT '',
    avatar_url           TEXT NOT NULL DEFAULT '',
    email_verified       INTEGER NOT NULL DEFAULT 0,
    role                 TEXT NOT NULL DEFAULT 'student'
);

-- The XP ledger; see the SQLite schema for why UNIQUE(student_id, event_key)
-- is the duplicate-award guard and why amounts may be negative.
CREATE TABLE IF NOT EXISTS xp_events (
    id         BIGSERIAL PRIMARY KEY,
    student_id TEXT NOT NULL,
    amount     INTEGER NOT NULL,
    reason     TEXT NOT NULL DEFAULT '',
    event_key  TEXT NOT NULL,
    created_at BIGINT NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_xp_events_event ON xp_events(student_id, event_key);
CREATE INDEX IF NOT EXISTS idx_xp_events_week ON xp_events(created_at);
CREATE INDEX IF NOT EXISTS idx_xp_events_student_week ON xp_events(student_id, created_at);

-- Weekly league membership. See the SQLite schema for why the week marker and
-- the seats are written in one transaction.
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
    created_at BIGINT NOT NULL
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
    attempt_id    BIGINT  NOT NULL DEFAULT 0,
    created_at    BIGINT  NOT NULL,
    PRIMARY KEY (student_id, operation, submission_id),
    FOREIGN KEY (student_id) REFERENCES students(id)
);

CREATE TABLE IF NOT EXISTS study_answer_effects (
    id            BIGSERIAL PRIMARY KEY,
    student_id    TEXT    NOT NULL,
    submission_id TEXT    NOT NULL,
    effect        TEXT    NOT NULL,
    payload_json  TEXT    NOT NULL DEFAULT '',
    created_at    BIGINT  NOT NULL,
    attempts      INTEGER NOT NULL DEFAULT 0,
    next_attempt_at BIGINT NOT NULL,
    claimed_until BIGINT  NOT NULL DEFAULT 0,
    processed_at  BIGINT  NOT NULL DEFAULT 0,
    last_error    TEXT    NOT NULL DEFAULT '',
    UNIQUE (student_id, submission_id, effect),
    FOREIGN KEY (student_id) REFERENCES students(id)
);

CREATE INDEX IF NOT EXISTS idx_study_answer_effects_pending
    ON study_answer_effects (processed_at, next_attempt_at, claimed_until);

CREATE TABLE IF NOT EXISTS league_weeks (
    week      TEXT PRIMARY KEY,
    closed_at BIGINT NOT NULL
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
    placement_seeded INTEGER NOT NULL DEFAULT 0,
    is_provisional   INTEGER NOT NULL DEFAULT 0,
    fsrs_stability   DOUBLE PRECISION,
    fsrs_difficulty  DOUBLE PRECISION,
    fsrs_last_review TEXT,
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
    timestamp       TEXT    NOT NULL DEFAULT (now()::text),
    question        TEXT    NOT NULL DEFAULT '',
    source          TEXT    NOT NULL DEFAULT '',
    explanation     TEXT    NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_progress_student  ON concept_progress(student_id);
CREATE INDEX IF NOT EXISTS idx_progress_mastered ON concept_progress(student_id, status, mastered_at);
CREATE INDEX IF NOT EXISTS idx_sessions_student  ON sessions(student_id);
CREATE INDEX IF NOT EXISTS idx_attempts_session  ON attempts(session_id);
CREATE INDEX IF NOT EXISTS idx_attempts_student  ON attempts(student_id, timestamp);
CREATE INDEX IF NOT EXISTS idx_attempts_cover    ON attempts(student_id, concept_id, timestamp);

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
    difficulty      DOUBLE PRECISION,
    kp_ordinal      INTEGER NOT NULL DEFAULT 0,
    updated_at      TEXT NOT NULL DEFAULT (now()::text)
);

CREATE INDEX IF NOT EXISTS idx_active_sessions_student ON active_sessions(student_id);
CREATE INDEX IF NOT EXISTS idx_active_sessions_attempt ON active_sessions(attempt_id);

CREATE TABLE IF NOT EXISTS student_topic_speed (
    student_id    TEXT NOT NULL,
    concept_id    TEXT NOT NULL,
    efactor       REAL NOT NULL DEFAULT 2.5,
    interval      INTEGER NOT NULL DEFAULT 0,
    repetitions   INTEGER NOT NULL DEFAULT 0,
    learning_speed DOUBLE PRECISION NOT NULL DEFAULT 1.0,
    updated_at    TEXT NOT NULL DEFAULT (now()::text),
    PRIMARY KEY (student_id, concept_id),
    FOREIGN KEY (student_id) REFERENCES students(id)
);
CREATE INDEX IF NOT EXISTS idx_topic_speed_student ON student_topic_speed(student_id);

CREATE TABLE IF NOT EXISTS quiz_completions (
    id           SERIAL PRIMARY KEY,
    student_id   TEXT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    completed_at TEXT NOT NULL DEFAULT (now()::text),
    xp_total     INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_quiz_completions_student ON quiz_completions(student_id, completed_at);

CREATE TABLE IF NOT EXISTS timed_quiz_baselines (
    student_id   TEXT PRIMARY KEY REFERENCES students(id) ON DELETE CASCADE,
    baseline_xp  INTEGER NOT NULL DEFAULT 0,
    baseline_at  BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS timed_quiz_runs (
    id                TEXT PRIMARY KEY,
    student_id        TEXT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    started_at        BIGINT NOT NULL,
    deadline_at       BIGINT NOT NULL,
    completed_at      BIGINT NOT NULL DEFAULT 0,
    status            TEXT NOT NULL DEFAULT 'active',
    question_count    INTEGER NOT NULL DEFAULT 0,
    correct_count     INTEGER NOT NULL DEFAULT 0,
    unanswered_count  INTEGER NOT NULL DEFAULT 0,
    xp_earned         INTEGER NOT NULL DEFAULT 0,
    created_at        BIGINT NOT NULL
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
    difficulty       DOUBLE PRECISION,
    served_at        BIGINT NOT NULL,
    answered_at      BIGINT NOT NULL DEFAULT 0,
    elapsed_seconds  DOUBLE PRECISION NOT NULL DEFAULT 0,
    answer           TEXT NOT NULL DEFAULT '',
    correct          INTEGER NOT NULL DEFAULT 0,
    status           TEXT NOT NULL DEFAULT 'served',
    submission_id    TEXT NOT NULL DEFAULT '',
    xp_earned        INTEGER NOT NULL DEFAULT 0,
    UNIQUE(run_id, concept_id)
);
CREATE INDEX IF NOT EXISTS idx_timed_quiz_items_run
    ON timed_quiz_items(run_id, position);

CREATE TABLE IF NOT EXISTS review_events (
    id                BIGSERIAL PRIMARY KEY,
    student_id        TEXT NOT NULL,
    concept_id        TEXT NOT NULL,
    attempt_id        TEXT NOT NULL UNIQUE,
    rating            INTEGER NOT NULL,
    correct           INTEGER NOT NULL,
    elapsed_seconds   DOUBLE PRECISION NOT NULL,
    reviewed_at       BIGINT NOT NULL,
    stability_before  DOUBLE PRECISION,
    difficulty_before DOUBLE PRECISION,
    due_before        TEXT,
    stability_after   DOUBLE PRECISION,
    difficulty_after  DOUBLE PRECISION,
    due_after         TEXT,
    created_at        BIGINT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_review_events_student
    ON review_events(student_id, reviewed_at);

CREATE TABLE IF NOT EXISTS question_reports (
    id          SERIAL PRIMARY KEY,
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
    answer_given      TEXT NOT NULL DEFAULT '',
    grading_type      TEXT NOT NULL DEFAULT '',
    answer_hint       TEXT NOT NULL DEFAULT '',
    answer_input_mode TEXT NOT NULL DEFAULT '',
    reason      TEXT NOT NULL DEFAULT 'other',
    detail      TEXT NOT NULL DEFAULT '',
    status      TEXT NOT NULL DEFAULT 'open',
    resolution  TEXT NOT NULL DEFAULT '',
    resolved_by TEXT NOT NULL DEFAULT '',
    resolved_at TEXT NOT NULL DEFAULT '',
    created_at  TEXT NOT NULL DEFAULT (now()::text),
    origin_report_id TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_reports_status  ON question_reports(status, created_at);
CREATE INDEX IF NOT EXISTS idx_reports_concept ON question_reports(concept_id);

-- The last time an account opened its own report history. See the SQLite
-- schema for why one timestamp per account is enough.
CREATE TABLE IF NOT EXISTS report_views (
    student_id TEXT PRIMARY KEY REFERENCES students(id) ON DELETE CASCADE,
    seen_at    TEXT NOT NULL DEFAULT '',
    cleared_at TEXT NOT NULL DEFAULT ''
);

-- One optional screenshot per report, beside the moderation row. ON DELETE
-- CASCADE is what makes account deletion remove the image too.
CREATE TABLE IF NOT EXISTS report_images (
    report_id    BIGINT PRIMARY KEY REFERENCES question_reports(id) ON DELETE CASCADE,
    content_type TEXT NOT NULL,
    bytes        BYTEA NOT NULL,
    created_at   TEXT NOT NULL DEFAULT (now()::text)
);

-- Live diagnostic runs. See the SQLite schema for why the adaptive model is
-- stored as an opaque snapshot rather than derived from attempts.
CREATE TABLE IF NOT EXISTS diagnostic_sessions (
    id         TEXT PRIMARY KEY,
    student_id TEXT NOT NULL DEFAULT '',
    data       BYTEA NOT NULL,
    updated_at BIGINT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_diagnostic_sessions_student ON diagnostic_sessions(student_id);
CREATE INDEX IF NOT EXISTS idx_diagnostic_sessions_updated ON diagnostic_sessions(updated_at);

-- Staff invitations. Only the token hash is stored; a database read cannot yield a usable
-- credential. Rows are never deleted, so the history of who offered access survives.
CREATE TABLE IF NOT EXISTS admin_invitations (
    id          BIGSERIAL PRIMARY KEY,
    token_hash  TEXT NOT NULL,
    email       TEXT NOT NULL DEFAULT '',
    role        TEXT NOT NULL,
    invited_by  TEXT NOT NULL DEFAULT '',
    created_at  TEXT NOT NULL DEFAULT (now()::text),
    expires_at  TEXT NOT NULL,
    accepted_at TEXT NOT NULL DEFAULT '',
    accepted_by TEXT NOT NULL DEFAULT '',
    revoked_at  TEXT NOT NULL DEFAULT '',
    revoked_by  TEXT NOT NULL DEFAULT ''
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_admin_invitations_token ON admin_invitations(token_hash);
CREATE INDEX IF NOT EXISTS idx_admin_invitations_email ON admin_invitations(email);

CREATE TABLE IF NOT EXISTS server_sessions (
    kind       TEXT NOT NULL,
    key        TEXT NOT NULL,
    value      TEXT NOT NULL DEFAULT '',
    expires_at TEXT NOT NULL DEFAULT '',
    updated_at TEXT NOT NULL DEFAULT (now()::text),
    PRIMARY KEY (kind, key)
);

CREATE INDEX IF NOT EXISTS idx_server_sessions_expiry ON server_sessions(expires_at);

-- Per-identity daily AI request counters. See the SQLite schema for why this is
-- a day-scoped row keyed by an opaque usage_key rather than a students column.
CREATE TABLE IF NOT EXISTS ai_usage (
    usage_key TEXT NOT NULL,
    day       TEXT NOT NULL,
    requests  INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (usage_key, day)
);

CREATE INDEX IF NOT EXISTS idx_ai_usage_day ON ai_usage(day);

-- FIRe-inspired implicit credit, shadow only. See the SQLite schema.
CREATE TABLE IF NOT EXISTS implicit_credit_targets (
    student_id        TEXT    NOT NULL,
    target_concept_id TEXT    NOT NULL,
    credit_epoch      BIGINT  NOT NULL DEFAULT 1,
    created_at        BIGINT  NOT NULL,
    PRIMARY KEY (student_id, target_concept_id)
);

CREATE TABLE IF NOT EXISTS implicit_credit_allocations (
    student_id        TEXT   NOT NULL,
    source_attempt_id TEXT   NOT NULL,
    raw_units         BIGINT NOT NULL,
    allocated_units   BIGINT NOT NULL,
    target_count      INTEGER NOT NULL,
    dropped_count     INTEGER NOT NULL,
    created_at        BIGINT NOT NULL,
    PRIMARY KEY (student_id, source_attempt_id)
);

CREATE TABLE IF NOT EXISTS implicit_credit_events (
    student_id        TEXT   NOT NULL,
    source_attempt_id TEXT   NOT NULL,
    target_concept_id TEXT   NOT NULL,
    credit_epoch      BIGINT NOT NULL,
    scaled_units      BIGINT NOT NULL,
    retained_units    BIGINT NOT NULL,
    discarded_units   BIGINT NOT NULL,
    discard_reason    TEXT   NOT NULL DEFAULT '',
    voided_at         BIGINT,
    voided_by         TEXT   NOT NULL DEFAULT '',
    void_reason       TEXT   NOT NULL DEFAULT '',
    created_at        BIGINT NOT NULL,
    PRIMARY KEY (student_id, source_attempt_id, target_concept_id)
);

CREATE INDEX IF NOT EXISTS idx_credit_events_target
    ON implicit_credit_events(student_id, target_concept_id, credit_epoch);

CREATE TABLE IF NOT EXISTS implicit_credit_ordinals (
    student_id         TEXT   NOT NULL,
    target_concept_id  TEXT   NOT NULL,
    credit_epoch       BIGINT NOT NULL,
    ordinal            INTEGER NOT NULL,
    state              TEXT   NOT NULL CHECK (state IN ('banked','projected','stale')),
    banked_at          BIGINT NOT NULL,
    projected_at       BIGINT,
    baseline_due       TEXT,
    counterfactual_due TEXT,
    attempts           INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (student_id, target_concept_id, credit_epoch, ordinal)
);

CREATE TABLE IF NOT EXISTS implicit_credit_applications (
    student_id        TEXT   NOT NULL,
    target_concept_id TEXT   NOT NULL,
    credit_epoch      BIGINT NOT NULL,
    ordinal           INTEGER NOT NULL,
    source_concept_id TEXT   NOT NULL DEFAULT '',
    source_attempt_id TEXT   NOT NULL,
    applied_at        BIGINT NOT NULL,
    actual_due        TEXT   NOT NULL DEFAULT '',
    proposed_due      TEXT   NOT NULL DEFAULT '',
    stability_before  DOUBLE PRECISION,
    stability_after   DOUBLE PRECISION,
    difficulty_before DOUBLE PRECISION,
    difficulty_after  DOUBLE PRECISION,
    reviews_deferred  INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (student_id, target_concept_id, credit_epoch, ordinal),
    UNIQUE (student_id, source_attempt_id, target_concept_id, credit_epoch)
);

CREATE INDEX IF NOT EXISTS idx_credit_applications_attempt
    ON implicit_credit_applications(student_id, source_attempt_id);

CREATE TABLE IF NOT EXISTS fire_activations (
    id              TEXT PRIMARY KEY,
    mode            TEXT NOT NULL,
    cohort          TEXT NOT NULL,
    acknowledged_by TEXT NOT NULL,
    acknowledged_at BIGINT NOT NULL,
    updated_at      BIGINT NOT NULL,
    revision        BIGINT NOT NULL
);

-- Targeted remedial queue: 2-miss halt remedial prerequisites survive restarts.
CREATE TABLE IF NOT EXISTS student_remedial_queue (
    student_id TEXT NOT NULL,
    concept_id TEXT NOT NULL,
    priority   INTEGER NOT NULL,
    created_at BIGINT NOT NULL,
    PRIMARY KEY (student_id, concept_id)
);

CREATE INDEX IF NOT EXISTS idx_remedial_queue_student ON student_remedial_queue(student_id, priority, created_at);

CREATE TABLE IF NOT EXISTS schema_migrations (
    version    INTEGER PRIMARY KEY,
    applied_at TEXT NOT NULL DEFAULT (now()::text)
);

-- Product announcements. They are global rows, so ResetProgress and account
-- deletion never remove the announcement itself. Reads are per learner and are
-- removed with the account.
CREATE TABLE IF NOT EXISTS announcements (
    id           TEXT PRIMARY KEY,
    title        TEXT NOT NULL,
    body         TEXT NOT NULL,
    published    INTEGER NOT NULL DEFAULT 0,
    created_at   BIGINT NOT NULL,
    published_at BIGINT NOT NULL DEFAULT 0,
    created_by   TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_announcements_published ON announcements(published, published_at);

CREATE TABLE IF NOT EXISTS announcement_reads (
    student_id      TEXT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    announcement_id TEXT NOT NULL REFERENCES announcements(id) ON DELETE CASCADE,
    read_at         BIGINT NOT NULL,
    PRIMARY KEY (student_id, announcement_id)
);

CREATE INDEX IF NOT EXISTS idx_announcement_reads_announcement ON announcement_reads(announcement_id);

-- Announcement clears hide the current feed from one learner without deleting
-- the global announcement or another learner's feed.
CREATE TABLE IF NOT EXISTS announcement_views (
    student_id TEXT PRIMARY KEY REFERENCES students(id) ON DELETE CASCADE,
    cleared_at BIGINT NOT NULL
);

-- Administration. actor_id is deliberately not a foreign key to students(id): a trail that is
-- deleted along with the account it records is not a trail. before_json/after_json carry the
-- non-sensitive AdminUser projection only, so a credential cannot reach them.
CREATE TABLE IF NOT EXISTS admin_audit (
    id          BIGSERIAL PRIMARY KEY,
    actor_id    TEXT NOT NULL DEFAULT '',
    action      TEXT NOT NULL,
    entity_type TEXT NOT NULL DEFAULT '',
    entity_id   TEXT NOT NULL DEFAULT '',
    before_json TEXT NOT NULL DEFAULT '',
    after_json  TEXT NOT NULL DEFAULT '',
    created_at  TEXT NOT NULL DEFAULT (now()::text)
);

CREATE INDEX IF NOT EXISTS idx_admin_audit_created ON admin_audit(created_at);
CREATE INDEX IF NOT EXISTS idx_admin_audit_actor   ON admin_audit(actor_id);
CREATE INDEX IF NOT EXISTS idx_admin_audit_entity   ON admin_audit(entity_type, entity_id);

CREATE TABLE IF NOT EXISTS friendships (
    student_a    TEXT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    student_b    TEXT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    status       TEXT NOT NULL CHECK (status IN ('pending', 'accepted')),
    requested_by TEXT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    created_at   BIGINT NOT NULL,
    responded_at BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (student_a, student_b),
    CHECK (student_a < student_b),
    CHECK (requested_by = student_a OR requested_by = student_b)
);

CREATE INDEX IF NOT EXISTS idx_friendships_a_status ON friendships(student_a, status);
CREATE INDEX IF NOT EXISTS idx_friendships_b_status ON friendships(student_b, status);

CREATE TABLE IF NOT EXISTS friend_blocks (
    blocker_id TEXT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    blocked_id TEXT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    created_at BIGINT NOT NULL,
    PRIMARY KEY (blocker_id, blocked_id),
    CHECK (blocker_id <> blocked_id)
);

CREATE INDEX IF NOT EXISTS idx_friend_blocks_blocked ON friend_blocks(blocked_id);
`

func NewPostgresStore(dsn string) (*PostgresStore, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}
	// Fix 7: bounded pool, previously unlimited (database/sql default).
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
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

// rebindDollars renumbers `?` placeholders for Postgres. Statements shared by
// both engines are written once, in the SQLite spelling, and the Postgres
// caller rewrites them here rather than keeping a second copy that can drift.
func rebindDollars(query string) string {
	var b strings.Builder
	n := 0
	for _, r := range query {
		if r == '?' {
			n++
			b.WriteString(fmt.Sprintf("$%d", n))
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

func pgAuthMigrate(db *sql.DB) error {
	migrations := []string{
		"ALTER TABLE students ADD COLUMN IF NOT EXISTS username TEXT",
		"ALTER TABLE students ADD COLUMN IF NOT EXISTS password_hash TEXT",
		// REAL is four bytes in Postgres and eight in SQLite, so 1.4 came back as
		// 1.39999997 and a round trip stopped being equal to the value written.
		"ALTER TABLE student_topic_speed ALTER COLUMN learning_speed TYPE DOUBLE PRECISION",
		"ALTER TABLE students ADD COLUMN IF NOT EXISTS course_id TEXT",
		"ALTER TABLE students ADD COLUMN IF NOT EXISTS xp_total INTEGER NOT NULL DEFAULT 0",
		"ALTER TABLE students ADD COLUMN IF NOT EXISTS xp_today INTEGER NOT NULL DEFAULT 0",
		"ALTER TABLE students ADD COLUMN IF NOT EXISTS xp_date TEXT",
		"ALTER TABLE students ADD COLUMN IF NOT EXISTS diagnostic_completed INTEGER NOT NULL DEFAULT 0",
		"ALTER TABLE students ADD COLUMN IF NOT EXISTS diagnostic_completed_at TEXT NOT NULL DEFAULT ''",
		"ALTER TABLE students ADD COLUMN IF NOT EXISTS daily_xp_goal INTEGER NOT NULL DEFAULT 30",
		"ALTER TABLE students ADD COLUMN IF NOT EXISTS settings TEXT NOT NULL DEFAULT '{}'",
		"ALTER TABLE students ADD COLUMN IF NOT EXISTS league TEXT NOT NULL DEFAULT 'bronze'",
		"ALTER TABLE students ADD COLUMN IF NOT EXISTS league_week TEXT NOT NULL DEFAULT ''",
		"ALTER TABLE students ADD COLUMN IF NOT EXISTS league_moved INTEGER NOT NULL DEFAULT 0",
		"ALTER TABLE students ADD COLUMN IF NOT EXISTS share_token TEXT NOT NULL DEFAULT ''",
		"ALTER TABLE students ADD COLUMN IF NOT EXISTS profile_id TEXT NOT NULL DEFAULT ''",
		"ALTER TABLE students ADD COLUMN IF NOT EXISTS friend_code TEXT NOT NULL DEFAULT ''",
		"ALTER TABLE report_views ADD COLUMN IF NOT EXISTS cleared_at TEXT NOT NULL DEFAULT ''",
		"CREATE INDEX IF NOT EXISTS idx_students_username ON students(username)",
		"CREATE INDEX IF NOT EXISTS idx_students_share ON students(share_token)",
		"ALTER TABLE students ADD COLUMN IF NOT EXISTS email TEXT NOT NULL DEFAULT ''",
		"ALTER TABLE students ADD COLUMN IF NOT EXISTS google_id TEXT NOT NULL DEFAULT ''",
		"ALTER TABLE students ADD COLUMN IF NOT EXISTS avatar_url TEXT NOT NULL DEFAULT ''",
		"ALTER TABLE students ADD COLUMN IF NOT EXISTS email_verified INTEGER NOT NULL DEFAULT 0",
		"CREATE INDEX IF NOT EXISTS idx_students_email ON students(email)",
		"CREATE INDEX IF NOT EXISTS idx_students_google_id ON students(google_id)",
		"ALTER TABLE concept_progress ADD COLUMN IF NOT EXISTS weakness_score DOUBLE PRECISION NOT NULL DEFAULT 0",
		"ALTER TABLE concept_progress ADD COLUMN IF NOT EXISTS placement_seeded INTEGER NOT NULL DEFAULT 0",
		"ALTER TABLE concept_progress ADD COLUMN IF NOT EXISTS is_provisional INTEGER NOT NULL DEFAULT 0",
		"ALTER TABLE concept_progress ADD COLUMN IF NOT EXISTS fsrs_stability DOUBLE PRECISION",
		"ALTER TABLE concept_progress ADD COLUMN IF NOT EXISTS fsrs_difficulty DOUBLE PRECISION",
		"ALTER TABLE concept_progress ADD COLUMN IF NOT EXISTS fsrs_last_review TEXT",
		"ALTER TABLE active_sessions ADD COLUMN IF NOT EXISTS last_concept_id TEXT NOT NULL DEFAULT ''",
		"ALTER TABLE active_sessions ADD COLUMN IF NOT EXISTS session_review INTEGER NOT NULL DEFAULT 0",
		"ALTER TABLE active_sessions ADD COLUMN IF NOT EXISTS session_new INTEGER NOT NULL DEFAULT 0",
		"ALTER TABLE attempts ADD COLUMN IF NOT EXISTS question TEXT NOT NULL DEFAULT ''",
		"ALTER TABLE attempts ADD COLUMN IF NOT EXISTS source TEXT NOT NULL DEFAULT ''",
		"ALTER TABLE attempts ADD COLUMN IF NOT EXISTS explanation TEXT NOT NULL DEFAULT ''",
		// Nullable on purpose: existing rows have no difficulty, and backfilling a
		// plausible-looking number would be fabricating evidence.
		"ALTER TABLE attempts ADD COLUMN IF NOT EXISTS difficulty DOUBLE PRECISION",
		"ALTER TABLE active_sessions ADD COLUMN IF NOT EXISTS difficulty DOUBLE PRECISION",
		"ALTER TABLE active_sessions ADD COLUMN IF NOT EXISTS kp_ordinal INTEGER NOT NULL DEFAULT 0",
		// The administrative role. NOT NULL DEFAULT 'student' migrates every existing account
		// to a learner, which is the only safe direction: nothing in a pre-role database
		// carries evidence about who was administering anything, because it could not.
		"ALTER TABLE students ADD COLUMN IF NOT EXISTS role TEXT NOT NULL DEFAULT 'student'",
		"CREATE INDEX IF NOT EXISTS idx_students_role ON students(role)",
		// Admin V2 moderation state on reports; mirrors the SQLite migration.
		"ALTER TABLE question_reports ADD COLUMN IF NOT EXISTS resolution TEXT NOT NULL DEFAULT ''",
		"ALTER TABLE question_reports ADD COLUMN IF NOT EXISTS resolved_by TEXT NOT NULL DEFAULT ''",
		"ALTER TABLE question_reports ADD COLUMN IF NOT EXISTS resolved_at TEXT NOT NULL DEFAULT ''",
		"ALTER TABLE question_reports ADD COLUMN IF NOT EXISTS answer_given TEXT NOT NULL DEFAULT ''",
		"ALTER TABLE question_reports ADD COLUMN IF NOT EXISTS grading_type TEXT NOT NULL DEFAULT ''",
		"ALTER TABLE question_reports ADD COLUMN IF NOT EXISTS answer_hint TEXT NOT NULL DEFAULT ''",
		"ALTER TABLE question_reports ADD COLUMN IF NOT EXISTS answer_input_mode TEXT NOT NULL DEFAULT ''",
		"ALTER TABLE question_reports ADD COLUMN IF NOT EXISTS origin_report_id TEXT NOT NULL DEFAULT ''",
		"CREATE UNIQUE INDEX IF NOT EXISTS idx_reports_origin ON question_reports(origin_report_id) WHERE origin_report_id != ''",
		`CREATE TABLE IF NOT EXISTS avatar_images (
			student_id TEXT PRIMARY KEY REFERENCES students(id) ON DELETE CASCADE,
			content_type TEXT NOT NULL,
			bytes BYTEA NOT NULL,
			updated_at TEXT NOT NULL DEFAULT (now()::text)
		)`,
		`CREATE TABLE IF NOT EXISTS identities (
			provider TEXT NOT NULL,
			provider_id TEXT NOT NULL,
			student_id TEXT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
			email TEXT NOT NULL DEFAULT '',
			email_verified INTEGER NOT NULL DEFAULT 0,
			linked_at TEXT NOT NULL DEFAULT (now()::text),
			PRIMARY KEY (provider, provider_id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_identities_student ON identities(student_id)`,
		`CREATE TABLE IF NOT EXISTS email_verifications (
			token_hash TEXT PRIMARY KEY,
			student_id TEXT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
			expires_at TEXT NOT NULL,
			used INTEGER NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS link_tokens (
			token_hash TEXT PRIMARY KEY,
			student_id TEXT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
			expires_at TEXT NOT NULL,
			used INTEGER NOT NULL DEFAULT 0
		)`,
		`INSERT INTO identities (provider, provider_id, student_id, email, email_verified)
			SELECT 'google', google_id, id, email, 1 FROM students
			WHERE google_id IS NOT NULL AND google_id != ''
			ON CONFLICT DO NOTHING`,
		`CREATE TABLE IF NOT EXISTS password_resets (
			token_hash TEXT PRIMARY KEY,
			student_id TEXT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
			expires_at TEXT NOT NULL,
			used INTEGER NOT NULL DEFAULT 0
		)`,
	}
	for _, m := range migrations {
		if _, err := db.Exec(m); err != nil {
			// IF NOT EXISTS handles duplicates; ignore any other duplicate errors
			if !isPgDuplicateColumn(err) && !isDuplicateColumn(err) {
				return err
			}
		}
	}
	// One-time data backfills, version-guarded (Fix 7), see sqlite twin.
	dataMigrations := map[int][]string{
		1: {"UPDATE students SET daily_xp_goal = 30 WHERE daily_xp_goal = 150"},
		2: {`UPDATE students SET username = lower(username)
			WHERE lower(username) NOT IN (
				SELECT lower(username) FROM students WHERE username IS NOT NULL AND username != ''
				GROUP BY lower(username) HAVING COUNT(*) > 1
			)`},
		// v3 small-awards economy: proportional rescale preserves intent
		// (30→10, 150→50, 100→33); goals already at/below 10 are untouched.
		3: {"UPDATE students SET daily_xp_goal = GREATEST(1, ROUND(daily_xp_goal / 3.0)) WHERE daily_xp_goal > 10"},
		// v4 daily goal range: bounded 10..100, default 30. Mirrors the SQLite migration
		// exactly, ADR-012 requires the two stores agree on an observable result.
		4: {`UPDATE students SET daily_xp_goal = 30 WHERE daily_xp_goal = 10`,
			`UPDATE students SET daily_xp_goal = 100 WHERE daily_xp_goal > 100`,
			`UPDATE students SET daily_xp_goal = 10 WHERE daily_xp_goal < 10`},
		// v5 Admin V2. Mirrors the SQLite migration exactly: V1 `admin` becomes `owner`, and the
		// report statuses move to the V2 vocabulary. ADR-012 requires the two stores agree.
		5: {`UPDATE students SET role = 'owner' WHERE role = 'admin'`,
			`UPDATE question_reports SET status = 'reviewing' WHERE status = 'confirmed'`,
			`UPDATE question_reports SET status = 'resolved' WHERE status = 'fixed'`},
		6: {`INSERT INTO timed_quiz_baselines (student_id, baseline_xp, baseline_at)
			SELECT id, xp_total, (extract(epoch from now()) * 1000000000)::bigint FROM students
			ON CONFLICT (student_id) DO NOTHING`},
	}
	for v := 1; v <= len(dataMigrations); v++ {
		if err := runOncePostgres(db, v, dataMigrations[v]); err != nil {
			return err
		}
	}
	// Give existing accounts their stable public identity before the unique
	// indexes below make either lookup key ambiguous.
	if err := backfillFriendIdentities(db, rebindDollars); err != nil {
		return err
	}
	// Back the app-level one-email-one-account 409s with a partial unique
	// index (case-insensitive, legacy '' rows excluded). Retried each boot
	// until clean so dirty databases still start; see the SQLite twin.
	if err := ensureEmailUniqueIndexPostgres(db); err != nil {
		return err
	}
	// Usernames are the public handle; see the SQLite twin for the retry
	// policy and the legacy-duplicate skip.
	if err := ensureUsernameUniqueIndexPostgres(db); err != nil {
		return err
	}
	return nil
}

// ensureEmailUniqueIndexPostgres enforces one-email-one-account at the DB
// layer for every writer at once, the app-level 409 checks are
// check-then-insert and race. Expression index on lower(email) matches
// FindByEmail's lookup semantics; the partial predicate keeps legacy ”
// rows legal. While real duplicates exist the index is skipped with a log
// (boot proceeds) and retried next boot.
func ensureEmailUniqueIndexPostgres(db *sql.DB) error {
	var name string
	if err := db.QueryRow("SELECT indexname FROM pg_indexes WHERE schemaname = current_schema() AND indexname = 'uidx_students_email'").Scan(&name); err == nil && name != "" {
		return nil
	} else if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("check email unique index: %w", err)
	}
	var dupes int
	if err := db.QueryRow("SELECT COUNT(*) FROM (SELECT lower(email) FROM students WHERE email IS NOT NULL AND email <> '' GROUP BY lower(email) HAVING COUNT(*) > 1) d").Scan(&dupes); err != nil {
		return fmt.Errorf("audit duplicate emails: %w", err)
	}
	if dupes > 0 {
		log.Printf("storage: %d duplicate email groups; skipping unique email index until resolved", dupes)
		return nil
	}
	if _, err := db.Exec("CREATE UNIQUE INDEX uidx_students_email ON students(lower(email)) WHERE email IS NOT NULL AND email <> ''"); err != nil {
		return fmt.Errorf("create email unique index: %w", err)
	}
	return nil
}

// ensureUsernameUniqueIndexPostgres mirrors the SQLite username index. The
// app-level FindByUsername check is check-then-insert; this index is what makes
// the uniqueness guarantee hold under two simultaneous signups.
func ensureUsernameUniqueIndexPostgres(db *sql.DB) error {
	var name string
	if err := db.QueryRow("SELECT indexname FROM pg_indexes WHERE schemaname = current_schema() AND indexname = 'uidx_students_username'").Scan(&name); err == nil && name != "" {
		return nil
	} else if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("check username unique index: %w", err)
	}
	var dupes int
	if err := db.QueryRow("SELECT COUNT(*) FROM (SELECT lower(username) FROM students WHERE username IS NOT NULL AND username <> '' GROUP BY lower(username) HAVING COUNT(*) > 1) d").Scan(&dupes); err != nil {
		return fmt.Errorf("audit duplicate usernames: %w", err)
	}
	if dupes > 0 {
		log.Printf("storage: %d duplicate username groups; skipping unique username index until resolved", dupes)
		return nil
	}
	if _, err := db.Exec("CREATE UNIQUE INDEX uidx_students_username ON students(lower(username)) WHERE username IS NOT NULL AND username <> ''"); err != nil {
		return fmt.Errorf("create username unique index: %w", err)
	}
	return nil
}

// runOncePostgres applies versioned data backfills exactly once.
// ON CONFLICT DO NOTHING keeps concurrent first boots safe.
func runOncePostgres(db *sql.DB, version int, stmts []string) error {
	var one int
	if err := db.QueryRow("SELECT 1 FROM schema_migrations WHERE version = $1", version).Scan(&one); err == nil {
		return nil
	} else if err != sql.ErrNoRows {
		return fmt.Errorf("check schema migration %d: %w", version, err)
	}
	for _, stmt := range stmts {
		if _, err := db.Exec(stmt); err != nil {
			return fmt.Errorf("apply schema migration %d: %w", version, err)
		}
	}
	if _, err := db.Exec("INSERT INTO schema_migrations (version) VALUES ($1) ON CONFLICT DO NOTHING", version); err != nil {
		return fmt.Errorf("record schema migration %d: %w", version, err)
	}
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
	profileID, err := newProfileID()
	if err != nil {
		return nil, fmt.Errorf("create student: %w", err)
	}
	friendCode, err := newFriendCode()
	if err != nil {
		return nil, fmt.Errorf("create student: %w", err)
	}
	_, err = s.db.Exec(
		"INSERT INTO students (id, profile_id, friend_code, name, settings, created_at) VALUES ($1, $2, $3, $4, '{}', $5)",
		id, profileID, friendCode, name, now.Format(time.RFC3339),
	)
	if err != nil {
		return nil, fmt.Errorf("create student: %w", err)
	}
	if err := s.timedQuiz().ensureBaselineAt(id, now); err != nil {
		return nil, err
	}
	return &Student{ID: id, Name: name, Settings: "{}", CreatedAt: now}, nil
}

func (s *PostgresStore) ClaimGuestStudent(id, name string) (*Student, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("guest id is required")
	}
	name = strings.TrimSpace(name)
	if name == "" {
		name = "Guest"
	}
	now := time.Now().UTC()
	profileID, err := newProfileID()
	if err != nil {
		return nil, fmt.Errorf("claim guest student: %w", err)
	}
	friendCode, err := newFriendCode()
	if err != nil {
		return nil, fmt.Errorf("claim guest student: %w", err)
	}
	_, err = s.db.Exec(
		"INSERT INTO students (id, profile_id, friend_code, name, settings, created_at) VALUES ($1, $2, $3, $4, '{}', $5) ON CONFLICT (id) DO NOTHING",
		id, profileID, friendCode, name, now.Format(time.RFC3339),
	)
	if err != nil {
		return nil, fmt.Errorf("claim guest student: %w", err)
	}
	st, err := s.GetStudent(id)
	if err != nil {
		return nil, err
	}
	if st == nil {
		return nil, fmt.Errorf("claim guest student: row missing after upsert")
	}
	if err := s.timedQuiz().ensureBaselineAt(st.ID, st.CreatedAt); err != nil {
		return nil, err
	}
	return st, nil
}

func (s *PostgresStore) CreateUser(name, username, passwordHash string) (*Student, error) {
	id := newUUID()
	now := time.Now().UTC()
	profileID, err := newProfileID()
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	friendCode, err := newFriendCode()
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	_, err = s.db.Exec(
		"INSERT INTO students (id, profile_id, friend_code, name, username, password_hash, settings, created_at) VALUES ($1, $2, $3, $4, $5, $6, '{}', $7)",
		id, profileID, friendCode, name, username, passwordHash, now.Format(time.RFC3339),
	)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	if err := s.timedQuiz().ensureBaselineAt(id, now); err != nil {
		return nil, err
	}
	return &Student{ID: id, Name: name, Username: username, PasswordHash: passwordHash, Settings: "{}", CreatedAt: now}, nil
}

// AdoptGuestAsUser attaches credentials to an existing guest row. The WHERE
// clause repeats the claimability check so a concurrent guest claim cannot
// turn an already registered row into a second account.
func (s *PostgresStore) AdoptGuestAsUser(id, name, username, passwordHash string) (*Student, error) {
	res, err := s.db.Exec(
		"UPDATE students SET name = $1, username = $2, password_hash = $3 WHERE id = $4 AND COALESCE(password_hash, '') = '' AND COALESCE(email, '') = ''",
		name, username, passwordHash, id,
	)
	if err != nil {
		return nil, fmt.Errorf("adopt guest as user: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("adopt guest as user: %w", err)
	}
	if n != 1 {
		return nil, fmt.Errorf("adopt guest as user: guest row is not claimable")
	}
	st, err := s.GetStudent(id)
	if err != nil {
		return nil, err
	}
	if st != nil {
		if err := s.timedQuiz().ensureBaselineAt(st.ID, st.CreatedAt); err != nil {
			return nil, err
		}
	}
	return st, nil
}

func (s *PostgresStore) GetStudent(id string) (*Student, error) {
	row := s.db.QueryRow("SELECT id, name, username, password_hash, course_id, xp_total, xp_today, xp_date, diagnostic_completed, daily_xp_goal, settings, created_at, league, league_week, league_moved, share_token, email, google_id, avatar_url, email_verified FROM students WHERE id = $1", id)
	return scanStudent(row)
}

func (s *PostgresStore) FindByUsername(username string) (*Student, error) {
	row := s.db.QueryRow("SELECT id, name, username, password_hash, course_id, xp_total, xp_today, xp_date, diagnostic_completed, daily_xp_goal, settings, created_at, league, league_week, league_moved, share_token, email, google_id, avatar_url, email_verified FROM students WHERE username = $1", username)
	return scanStudent(row)
}

func (s *PostgresStore) FindStudentByIdentity(provider, providerID string) (*Student, error) {
	if provider == "" || providerID == "" {
		return nil, nil
	}
	row := s.db.QueryRow(`SELECT id, name, username, password_hash, course_id, xp_total, xp_today, xp_date, diagnostic_completed, daily_xp_goal, settings, created_at, league, league_week, league_moved, share_token, email, google_id, avatar_url, email_verified FROM students WHERE id = (
		SELECT student_id FROM identities WHERE provider = $1 AND provider_id = $2
	)`, provider, providerID)
	return scanStudent(row)
}

func (s *PostgresStore) FindByEmail(email string) (*Student, error) {
	if email == "" {
		return nil, nil
	}
	row := s.db.QueryRow("SELECT id, name, username, password_hash, course_id, xp_total, xp_today, xp_date, diagnostic_completed, daily_xp_goal, settings, created_at, league, league_week, league_moved, share_token, email, google_id, avatar_url, email_verified FROM students WHERE lower(email) = lower($1) AND email <> '' LIMIT 1", email)
	return scanStudent(row)
}

func (s *PostgresStore) CreateGoogleUser(name, email, googleID, avatarURL string) (*Student, error) {
	id := newUUID()
	now := time.Now().UTC()
	profileID, err := newProfileID()
	if err != nil {
		return nil, fmt.Errorf("create google user: %w", err)
	}
	friendCode, err := newFriendCode()
	if err != nil {
		return nil, fmt.Errorf("create google user: %w", err)
	}
	_, err = s.db.Exec(
		"INSERT INTO students (id, profile_id, friend_code, name, email, email_verified, google_id, avatar_url, settings, created_at) VALUES ($1, $2, $3, $4, $5, 1, $6, $7, '{}', $8)",
		id, profileID, friendCode, name, email, googleID, avatarURL, now.Format(time.RFC3339),
	)
	if err != nil {
		return nil, fmt.Errorf("create google user: %w", err)
	}
	if googleID != "" {
		if err := s.CreateIdentity("google", googleID, id, email, true); err != nil {
			return nil, err
		}
	}
	if err := s.timedQuiz().ensureBaselineAt(id, now); err != nil {
		return nil, err
	}
	return &Student{ID: id, Name: name, Email: email, EmailVerified: true, GoogleID: googleID, AvatarURL: avatarURL, Settings: "{}", CreatedAt: now}, nil
}

func (s *PostgresStore) SetEmailVerified(studentID string, verified bool) error {
	v := 0
	if verified {
		v = 1
	}
	_, err := s.db.Exec("UPDATE students SET email_verified = $1 WHERE id = $2", v, studentID)
	if err != nil {
		return fmt.Errorf("set email verified: %w", err)
	}
	return nil
}

func (s *PostgresStore) CreateIdentity(provider, providerID, studentID, email string, emailVerified bool) error {
	v := 0
	if emailVerified {
		v = 1
	}
	_, err := s.db.Exec(`INSERT INTO identities (provider, provider_id, student_id, email, email_verified, linked_at)
		VALUES ($1, $2, $3, $4, $5, now()::text)`, provider, providerID, studentID, email, v)
	if err != nil {
		return fmt.Errorf("create identity: %w", err)
	}
	return nil
}

func (s *PostgresStore) ListIdentities(studentID string) ([]Identity, error) {
	rows, err := s.db.Query("SELECT provider, provider_id, student_id, email, email_verified, linked_at FROM identities WHERE student_id = $1 ORDER BY provider", studentID)
	if err != nil {
		return nil, fmt.Errorf("list identities: %w", err)
	}
	defer rows.Close()
	var out []Identity
	for rows.Next() {
		var id Identity
		var verified int
		var linkedAt string
		if err := rows.Scan(&id.Provider, &id.ProviderID, &id.StudentID, &id.Email, &verified, &linkedAt); err != nil {
			return nil, fmt.Errorf("scan identity: %w", err)
		}
		id.EmailVerified = verified == 1
		if t, err := time.Parse(time.RFC3339, linkedAt); err == nil {
			id.LinkedAt = t
		}
		out = append(out, id)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list identities: %w", err)
	}
	return out, nil
}

func (s *PostgresStore) DeleteIdentity(provider, studentID string) error {
	_, err := s.db.Exec("DELETE FROM identities WHERE provider = $1 AND student_id = $2", provider, studentID)
	if err != nil {
		return fmt.Errorf("delete identity: %w", err)
	}
	return nil
}

func (s *PostgresStore) CreateEmailVerification(tokenHash string, studentID string, expiresAt time.Time) error {
	_, err := s.db.Exec("INSERT INTO email_verifications (token_hash, student_id, expires_at, used) VALUES ($1, $2, $3, 0)",
		tokenHash, studentID, expiresAt.UTC().Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("create email verification: %w", err)
	}
	return nil
}

func (s *PostgresStore) ConsumeEmailVerification(tokenHash string) (string, bool, error) {
	return consumeSingleUseTokenPG(s.db, "email_verifications", tokenHash)
}

func (s *PostgresStore) CreateLinkToken(tokenHash string, studentID string, expiresAt time.Time) error {
	_, err := s.db.Exec("INSERT INTO link_tokens (token_hash, student_id, expires_at, used) VALUES ($1, $2, $3, 0)",
		tokenHash, studentID, expiresAt.UTC().Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("create link token: %w", err)
	}
	return nil
}

func (s *PostgresStore) ConsumeLinkToken(tokenHash string) (string, bool, error) {
	return consumeSingleUseTokenPG(s.db, "link_tokens", tokenHash)
}

// consumeSingleUseTokenPG atomically validates (exists, unused, unexpired) and
// burns a token row, returning the student it belongs to.
func consumeSingleUseTokenPG(db *sql.DB, table, tokenHash string) (string, bool, error) {
	tx, err := db.Begin()
	if err != nil {
		return "", false, fmt.Errorf("token tx: %w", err)
	}
	defer tx.Rollback()
	var studentID, expiresAt string
	var used int
	err = tx.QueryRow("SELECT student_id, expires_at, used FROM "+table+" WHERE token_hash = $1", tokenHash).
		Scan(&studentID, &expiresAt, &used)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, fmt.Errorf("lookup token: %w", err)
	}
	exp, err := time.Parse(time.RFC3339, expiresAt)
	if err != nil || used != 0 || time.Now().UTC().After(exp) {
		return "", false, nil
	}
	if _, err := tx.Exec("UPDATE "+table+" SET used = 1 WHERE token_hash = $1", tokenHash); err != nil {
		return "", false, fmt.Errorf("burn token: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return "", false, fmt.Errorf("commit token: %w", err)
	}
	return studentID, true, nil
}

func (s *PostgresStore) SetShareToken(studentID, token string) error {
	_, err := s.db.Exec("UPDATE students SET share_token = $1 WHERE id = $2", token, studentID)
	if err != nil {
		return fmt.Errorf("set share token: %w", err)
	}
	return nil
}

func (s *PostgresStore) SetAvatarURL(studentID, avatarURL string) error {
	_, err := s.db.Exec("UPDATE students SET avatar_url = $1 WHERE id = $2", avatarURL, studentID)
	if err != nil {
		return fmt.Errorf("set avatar url: %w", err)
	}
	return nil
}

func (s *PostgresStore) CreateOAuthUser(provider, providerID, name, email string, emailVerified bool, avatarURL string) (*Student, error) {
	id := newUUID()
	now := time.Now().UTC()
	ev := 0
	if emailVerified {
		ev = 1
	}
	if name == "" {
		name = email
	}
	profileID, err := newProfileID()
	if err != nil {
		return nil, fmt.Errorf("create oauth user: %w", err)
	}
	friendCode, err := newFriendCode()
	if err != nil {
		return nil, fmt.Errorf("create oauth user: %w", err)
	}
	_, err = s.db.Exec(
		"INSERT INTO students (id, profile_id, friend_code, name, email, email_verified, avatar_url, settings, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, '{}', $8)",
		id, profileID, friendCode, name, email, ev, avatarURL, now.Format(time.RFC3339),
	)
	if err != nil {
		return nil, fmt.Errorf("create oauth user: %w", err)
	}
	if err := s.CreateIdentity(provider, providerID, id, email, emailVerified); err != nil {
		return nil, err
	}
	if err := s.timedQuiz().ensureBaselineAt(id, now); err != nil {
		return nil, err
	}
	return &Student{ID: id, Name: name, Email: email, EmailVerified: emailVerified, AvatarURL: avatarURL, Settings: "{}", CreatedAt: now}, nil
}

func (s *PostgresStore) GetStudentByShareToken(token string) (*Student, error) {
	if token == "" {
		return nil, nil
	}
	row := s.db.QueryRow("SELECT id, name, username, password_hash, course_id, xp_total, xp_today, xp_date, diagnostic_completed, daily_xp_goal, settings, created_at, league, league_week, league_moved, share_token, email, google_id, avatar_url, email_verified FROM students WHERE share_token = $1", token)
	return scanStudent(row)
}

func (s *PostgresStore) SetCourseID(studentID, courseID string) error {
	_, err := s.db.Exec("UPDATE students SET course_id = $1 WHERE id = $2", courseID, studentID)
	if err != nil {
		return fmt.Errorf("set course_id: %w", err)
	}
	return nil
}

func (s *PostgresStore) AwardXP(studentID string, amount int, reason, eventKey string) error {
	if studentID == "" {
		return fmt.Errorf("award xp: empty student")
	}
	if eventKey == "" {
		return fmt.Errorf("award xp: empty event key")
	}
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("award xp: %w", err)
	}
	defer tx.Rollback()

	// ON CONFLICT DO NOTHING plus the unique index is the concurrency guard;
	// see the SQLite implementation for the reasoning.
	res, err := tx.Exec(`
		INSERT INTO xp_events (student_id, amount, reason, event_key, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (student_id, event_key) DO NOTHING
	`, studentID, amount, reason, eventKey, time.Now().UTC().UnixNano())
	if err != nil {
		return fmt.Errorf("award xp: insert event: %w", err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return tx.Commit()
	}

	today := time.Now().UTC().Format("2006-01-02")
	if _, err := tx.Exec(`
		UPDATE students
		SET xp_total = GREATEST(0, xp_total + $1),
		    xp_today = CASE WHEN xp_date = $2 THEN GREATEST(0, xp_today + $3) ELSE GREATEST(0, $4) END,
		    xp_date  = $5
		WHERE id = $6
	`, amount, today, amount, amount, today, studentID); err != nil {
		return fmt.Errorf("award xp: update counters: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("award xp: %w", err)
	}
	return nil
}

func (s *PostgresStore) WeeklyXPByStudent(since time.Time) (map[string]int, error) {
	rows, err := s.db.Query(`
		SELECT student_id, COALESCE(SUM(amount), 0)
		FROM xp_events
		WHERE created_at >= $1
		GROUP BY student_id
	`, since.UTC().UnixNano())
	if err != nil {
		return nil, fmt.Errorf("weekly xp: %w", err)
	}
	defer rows.Close()
	out := map[string]int{}
	for rows.Next() {
		var id string
		var sum int
		if err := rows.Scan(&id, &sum); err != nil {
			return nil, fmt.Errorf("weekly xp: %w", err)
		}
		out[id] = sum
	}
	return out, rows.Err()
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
	_, err := s.db.Exec(
		"UPDATE students SET diagnostic_completed = 1, diagnostic_completed_at = CASE WHEN diagnostic_completed_at = '' THEN $1 ELSE diagnostic_completed_at END WHERE id = $2",
		time.Now().UTC().Format(time.RFC3339), studentID,
	)
	if err != nil {
		return fmt.Errorf("set diagnostic_completed: %w", err)
	}
	return nil
}

// DiagnosticCompletedAt mirrors SQLiteStore.DiagnosticCompletedAt: found is
// false when the flag predates the timestamp column, so no date is invented.
func (s *PostgresStore) DiagnosticCompletedAt(studentID string) (time.Time, bool, error) {
	var raw string
	err := s.db.QueryRow(
		"SELECT diagnostic_completed_at FROM students WHERE id = $1", studentID,
	).Scan(&raw)
	if err == sql.ErrNoRows {
		return time.Time{}, false, nil
	}
	if err != nil {
		return time.Time{}, false, fmt.Errorf("get diagnostic completed_at: %w", err)
	}
	if raw == "" {
		return time.Time{}, false, nil
	}
	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return time.Time{}, false, nil
	}
	return t, true, nil
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

func (s *PostgresStore) SetAvatarImage(studentID string, contentType string, data []byte) error {
	_, err := s.db.Exec(`INSERT INTO avatar_images (student_id, content_type, bytes, updated_at)
		VALUES ($1, $2, $3, now()::text)
		ON CONFLICT(student_id) DO UPDATE SET content_type = excluded.content_type, bytes = excluded.bytes, updated_at = now()::text`,
		studentID, contentType, data)
	if err != nil {
		return fmt.Errorf("set avatar image: %w", err)
	}
	return nil
}

func (s *PostgresStore) GetAvatarImage(studentID string) (string, []byte, bool, error) {
	var contentType string
	var data []byte
	err := s.db.QueryRow("SELECT content_type, bytes FROM avatar_images WHERE student_id = $1", studentID).Scan(&contentType, &data)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil, false, nil
		}
		return "", nil, false, fmt.Errorf("get avatar image: %w", err)
	}
	return contentType, data, true, nil
}

func (s *PostgresStore) ClearAvatarImage(studentID string) error {
	_, err := s.db.Exec("DELETE FROM avatar_images WHERE student_id = $1", studentID)
	if err != nil {
		return fmt.Errorf("clear avatar image: %w", err)
	}
	return nil
}

func (s *PostgresStore) UpdateStudentName(studentID string, name string) error {
	_, err := s.db.Exec("UPDATE students SET name = $1 WHERE id = $2", name, studentID)
	if err != nil {
		return fmt.Errorf("update student name: %w", err)
	}
	return nil
}

func (s *PostgresStore) SetUsername(studentID string, username string) error {
	_, err := s.db.Exec("UPDATE students SET username = $1 WHERE id = $2", username, studentID)
	if err != nil {
		return fmt.Errorf("set username: %w", err)
	}
	return nil
}

func (s *PostgresStore) ListStudentsMissingUsernames() ([]string, error) {
	rows, err := s.db.Query("SELECT id FROM students WHERE username IS NULL OR TRIM(username) = ''")
	if err != nil {
		return nil, fmt.Errorf("list missing usernames: %w", err)
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan missing username: %w", err)
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (s *PostgresStore) SetEmail(studentID string, email string) error {
	_, err := s.db.Exec("UPDATE students SET email = $1 WHERE id = $2", email, studentID)
	if err != nil {
		return fmt.Errorf("set email: %w", err)
	}
	return nil
}

func (s *PostgresStore) SetPasswordHash(studentID string, hash string) error {
	_, err := s.db.Exec("UPDATE students SET password_hash = $1 WHERE id = $2", hash, studentID)
	if err != nil {
		return fmt.Errorf("set password hash: %w", err)
	}
	return nil
}

func (s *PostgresStore) CreatePasswordReset(tokenHash string, studentID string, expiresAt time.Time) error {
	_, err := s.db.Exec("INSERT INTO password_resets (token_hash, student_id, expires_at, used) VALUES ($1, $2, $3, 0)",
		tokenHash, studentID, expiresAt.UTC().Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("create password reset: %w", err)
	}
	return nil
}

// ConsumePasswordReset atomically validates (exists, unused, unexpired) and
// burns a reset token, returning the student it belongs to.
func (s *PostgresStore) ConsumePasswordReset(tokenHash string) (string, bool, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return "", false, fmt.Errorf("reset tx: %w", err)
	}
	defer tx.Rollback()
	var studentID, expiresAt string
	var used int
	err = tx.QueryRow("SELECT student_id, expires_at, used FROM password_resets WHERE token_hash = $1", tokenHash).
		Scan(&studentID, &expiresAt, &used)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, fmt.Errorf("lookup reset: %w", err)
	}
	exp, err := time.Parse(time.RFC3339, expiresAt)
	if err != nil || used != 0 || time.Now().UTC().After(exp) {
		return "", false, nil
	}
	if _, err := tx.Exec("UPDATE password_resets SET used = 1 WHERE token_hash = $1", tokenHash); err != nil {
		return "", false, fmt.Errorf("burn reset: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return "", false, fmt.Errorf("commit reset: %w", err)
	}
	return studentID, true, nil
}

// Progress

func (s *PostgresStore) GetProgress(studentID, conceptID string) (*ConceptProgress, error) {
	row := s.db.QueryRow(`
		SELECT student_id, concept_id, status, streak, best_streak,
		       avg_response_time, attempts, last_attempted, last_reviewed,
		       next_review_due, sm2_repetitions, sm2_interval, sm2_efactor,
		       mastered_at, weakness_score, placement_seeded, is_provisional,
		       fsrs_stability, fsrs_difficulty, fsrs_last_review
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
		       mastered_at, weakness_score, placement_seeded, is_provisional,
		       fsrs_stability, fsrs_difficulty, fsrs_last_review
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
			 mastered_at, weakness_score, placement_seeded, is_provisional,
			 fsrs_stability, fsrs_difficulty, fsrs_last_review)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
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
			weakness_score    = EXCLUDED.weakness_score,
			placement_seeded  = EXCLUDED.placement_seeded,
			is_provisional    = EXCLUDED.is_provisional,
			fsrs_stability    = EXCLUDED.fsrs_stability,
			fsrs_difficulty   = EXCLUDED.fsrs_difficulty,
			fsrs_last_review  = EXCLUDED.fsrs_last_review
	`,
		p.StudentID, p.ConceptID, p.Status, p.Streak, p.BestStreak,
		p.AvgResponseTime, p.Attempts,
		nullTime(p.LastAttempted), nullTime(p.LastReviewed),
		nullTime(p.NextReviewDue), p.SM2Repetitions, p.SM2Interval,
		p.SM2EFactor, nullTime(p.MasteredAt), p.WeaknessScore, boolToInt(p.PlacementSeeded),
		boolToInt(p.IsProvisional),
		nullFloat(p.FSRSStability), nullFloat(p.FSRSDifficulty), nullTime(p.FSRSLastReview),
	)
	if err != nil {
		return fmt.Errorf("upsert progress: %w", err)
	}
	return nil
}

func (s *PostgresStore) UpsertProgressBatch(ps []*ConceptProgress) error {
	if len(ps) == 0 {
		return nil
	}
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("begin progress batch: %w", err)
	}
	stmt, err := tx.Prepare(`
		INSERT INTO concept_progress
			(student_id, concept_id, status, streak, best_streak,
			 avg_response_time, attempts, last_attempted, last_reviewed,
			 next_review_due, sm2_repetitions, sm2_interval, sm2_efactor,
			 mastered_at, weakness_score, placement_seeded, is_provisional,
			 fsrs_stability, fsrs_difficulty, fsrs_last_review)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
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
			weakness_score    = EXCLUDED.weakness_score,
			placement_seeded  = EXCLUDED.placement_seeded,
			is_provisional    = EXCLUDED.is_provisional,
			fsrs_stability    = EXCLUDED.fsrs_stability,
			fsrs_difficulty   = EXCLUDED.fsrs_difficulty,
			fsrs_last_review  = EXCLUDED.fsrs_last_review
	`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("prepare progress batch: %w", err)
	}
	defer stmt.Close()
	for _, p := range ps {
		if _, err := stmt.Exec(
			p.StudentID, p.ConceptID, p.Status, p.Streak, p.BestStreak,
			p.AvgResponseTime, p.Attempts,
			nullTime(p.LastAttempted), nullTime(p.LastReviewed),
			nullTime(p.NextReviewDue), p.SM2Repetitions, p.SM2Interval,
			p.SM2EFactor, nullTime(p.MasteredAt), p.WeaknessScore, boolToInt(p.PlacementSeeded),
			boolToInt(p.IsProvisional),
			nullFloat(p.FSRSStability), nullFloat(p.FSRSDifficulty), nullTime(p.FSRSLastReview),
		); err != nil {
			tx.Rollback()
			return fmt.Errorf("exec progress batch: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit progress batch: %w", err)
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

// EnsureSession inserts a sessions row idempotently for ephemeral
// diagnostic/quiz UUIDs (attempts.session_id references sessions(id)).
func (s *PostgresStore) EnsureSession(id, studentID string) error {
	_, err := s.db.Exec(
		`INSERT INTO sessions (id, student_id, started_at) VALUES ($1, $2, now()::text)
		 ON CONFLICT (id) DO NOTHING`,
		id, studentID,
	)
	if err != nil {
		return fmt.Errorf("ensure session: %w", err)
	}
	return nil
}

// ServerSessions is a durable KV for restart-proof server state.

func (s *PostgresStore) UpsertServerSession(kind, key, value, expiresAt string) error {
	_, err := s.db.Exec(`
		INSERT INTO server_sessions (kind, key, value, expires_at, updated_at)
		VALUES ($1, $2, $3, $4, now()::text)
		ON CONFLICT(kind, key) DO UPDATE SET
			value = EXCLUDED.value,
			expires_at = EXCLUDED.expires_at,
			updated_at = now()::text
	`, kind, key, value, expiresAt)
	if err != nil {
		return fmt.Errorf("upsert server session: %w", err)
	}
	return nil
}

func (s *PostgresStore) GetServerSession(kind, key string) (string, string, bool, error) {
	row := s.db.QueryRow("SELECT value, expires_at FROM server_sessions WHERE kind = $1 AND key = $2", kind, key)
	var value, expiresAt string
	if err := row.Scan(&value, &expiresAt); err != nil {
		if err == sql.ErrNoRows {
			return "", "", false, nil
		}
		return "", "", false, fmt.Errorf("get server session: %w", err)
	}
	return value, expiresAt, true, nil
}

func (s *PostgresStore) DeleteServerSession(kind, key string) error {
	if _, err := s.db.Exec("DELETE FROM server_sessions WHERE kind = $1 AND key = $2", kind, key); err != nil {
		return fmt.Errorf("delete server session: %w", err)
	}
	return nil
}

// ReserveAIUsage mirrors the SQLite twin: one transaction over every limit,
// committed only when each counter is within its ceiling, so a denied request
// leaves no increment behind and concurrent callers serialise on the row locks.
func (s *PostgresStore) ReserveAIUsage(ctx context.Context, day string, limits []ai.AIUsageLimit) (bool, error) {
	if len(limits) == 0 {
		return true, nil
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return false, fmt.Errorf("reserve ai usage: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	for _, l := range limits {
		var n int
		err := tx.QueryRowContext(ctx, `
			INSERT INTO ai_usage (usage_key, day, requests) VALUES ($1, $2, 1)
			ON CONFLICT (usage_key, day) DO UPDATE SET requests = ai_usage.requests + 1
			RETURNING requests
		`, l.Key, day).Scan(&n)
		if err != nil {
			return false, fmt.Errorf("reserve ai usage: %w", err)
		}
		if l.Cap > 0 && n > l.Cap {
			return false, nil
		}
	}
	if err := tx.Commit(); err != nil {
		return false, fmt.Errorf("commit ai usage: %w", err)
	}
	return true, nil
}

// AIUsage reads the counter without incrementing. Diagnostics only.
func (s *PostgresStore) AIUsage(ctx context.Context, usageKey, day string) (int, error) {
	var n int
	err := s.db.QueryRowContext(ctx, "SELECT requests FROM ai_usage WHERE usage_key = $1 AND day = $2", usageKey, day).Scan(&n)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("read ai usage: %w", err)
	}
	return n, nil
}

// PruneAIUsage mirrors the SQLite twin: a range delete over day-scoped counters
// that have all passed. `day` is zero-padded YYYY-MM-DD, so the text comparison
// is a chronological one.
func (s *PostgresStore) PruneAIUsage(beforeDay string) (int64, error) {
	res, err := s.db.Exec("DELETE FROM ai_usage WHERE day < $1", beforeDay)
	if err != nil {
		return 0, fmt.Errorf("prune ai usage: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("prune ai usage: %w", err)
	}
	return n, nil
}

func (s *PostgresStore) SweepServerSessions() error {
	// expires_at is RFC3339 UTC, compare in Go (see sqlite twin).
	rows, err := s.db.Query("SELECT kind, key, expires_at FROM server_sessions WHERE expires_at != ''")
	if err != nil {
		return fmt.Errorf("sweep server sessions: %w", err)
	}
	defer rows.Close()
	type entry struct{ kind, key string }
	var expired []entry
	now := time.Now().UTC()
	for rows.Next() {
		var k, key, exp string
		if err := rows.Scan(&k, &key, &exp); err != nil {
			return fmt.Errorf("scan server sessions: %w", err)
		}
		if t, err := time.Parse(time.RFC3339, exp); err == nil && !t.After(now) {
			expired = append(expired, entry{k, key})
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("sweep server sessions: %w", err)
	}
	rows.Close()
	for _, e := range expired {
		if _, err := s.db.Exec("DELETE FROM server_sessions WHERE kind = $1 AND key = $2", e.kind, e.key); err != nil {
			return fmt.Errorf("sweep server sessions: %w", err)
		}
	}
	return nil
}

func (s *PostgresStore) GetActiveSession(sessionID string) (*ActiveSession, error) {
	row := s.db.QueryRow(`
		SELECT session_id, student_id, concept_id, concept_name, expected_answer,
		       attempt_id, question, explanation, diagram, is_review, answered,
		       last_concept_id, session_review, session_new, difficulty, kp_ordinal, updated_at
		FROM active_sessions WHERE session_id = $1
	`, sessionID)
	var a ActiveSession
	var isReview, answered int
	var updatedAt string
	if err := row.Scan(&a.SessionID, &a.StudentID, &a.ConceptID, &a.ConceptName, &a.ExpectedAnswer,
		&a.AttemptID, &a.Question, &a.Explanation, &a.Diagram, &isReview, &answered,
		&a.LastConceptID, &a.SessionReview, &a.SessionNew, &a.Difficulty, &a.KPOrdinal, &updatedAt); err != nil {
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
			 last_concept_id, session_review, session_new, difficulty, kp_ordinal, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
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
			difficulty      = EXCLUDED.difficulty,
			kp_ordinal      = EXCLUDED.kp_ordinal,
			updated_at      = EXCLUDED.updated_at
	`, a.SessionID, a.StudentID, a.ConceptID, a.ConceptName, a.ExpectedAnswer,
		a.AttemptID, a.Question, a.Explanation, a.Diagram, boolToInt(a.IsReview), boolToInt(a.Answered),
		a.LastConceptID, a.SessionReview, a.SessionNew, a.Difficulty, a.KPOrdinal, now)
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

// ResetProgress mirrors the SQLite twin (see Repository docs for scope).
func (s *PostgresStore) ResetProgress(studentID string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("begin progress reset: %w", err)
	}
	if err := lockStudentCreditTx(context.Background(), tx, studentID); err != nil {
		tx.Rollback()
		return err
	}
	exec := func(stmt string, args ...interface{}) error {
		if _, err := tx.Exec(stmt, args...); err != nil {
			tx.Rollback()
			return err
		}
		return nil
	}
	stmts := []struct {
		q    string
		args []interface{}
	}{
		{"DELETE FROM attempts WHERE student_id = $1", []interface{}{studentID}},
		{"DELETE FROM answer_submissions WHERE student_id = $1", []interface{}{studentID}},
		{"DELETE FROM study_answer_effects WHERE student_id = $1", []interface{}{studentID}},
		{"DELETE FROM active_sessions WHERE student_id = $1", []interface{}{studentID}},
		{"DELETE FROM sessions WHERE student_id = $1", []interface{}{studentID}},
		{"DELETE FROM implicit_credit_applications WHERE student_id = $1", []interface{}{studentID}},
		{"DELETE FROM implicit_credit_ordinals WHERE student_id = $1", []interface{}{studentID}},
		{"DELETE FROM implicit_credit_events WHERE student_id = $1", []interface{}{studentID}},
		{"DELETE FROM implicit_credit_allocations WHERE student_id = $1", []interface{}{studentID}},
		{"DELETE FROM implicit_credit_targets WHERE student_id = $1", []interface{}{studentID}},
		{"DELETE FROM concept_progress WHERE student_id = $1", []interface{}{studentID}},
		{"DELETE FROM student_topic_speed WHERE student_id = $1", []interface{}{studentID}},
		{"DELETE FROM quiz_completions WHERE student_id = $1", []interface{}{studentID}},
		{"DELETE FROM timed_quiz_items WHERE run_id IN (SELECT id FROM timed_quiz_runs WHERE student_id = $1)", []interface{}{studentID}},
		{"DELETE FROM timed_quiz_runs WHERE student_id = $1", []interface{}{studentID}},
		{"DELETE FROM timed_quiz_baselines WHERE student_id = $1", []interface{}{studentID}},
		{"DELETE FROM review_events WHERE student_id = $1", []interface{}{studentID}},
		{"DELETE FROM student_remedial_queue WHERE student_id = $1", []interface{}{studentID}},
		{"DELETE FROM server_sessions WHERE kind IN ('study_expected','study_plan') AND (key = $1 OR key LIKE $1 || '|%')", []interface{}{studentID}},
		{`UPDATE students SET xp_total = 0, xp_today = 0, diagnostic_completed = 0, diagnostic_completed_at = '',
			league = 'bronze', league_week = '', league_moved = 0 WHERE id = $1`, []interface{}{studentID}},
	}
	for _, st := range stmts {
		if err := exec(st.q, st.args...); err != nil {
			return fmt.Errorf("reset progress: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit progress reset: %w", err)
	}
	return nil
}

// DeleteAccount mirrors the SQLite twin (see Repository docs for scope).
func (s *PostgresStore) DeleteAccount(studentID string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("begin account delete: %w", err)
	}
	if err := lockStudentCreditTx(context.Background(), tx, studentID); err != nil {
		tx.Rollback()
		return err
	}
	exec := func(stmt string, args ...interface{}) error {
		if _, err := tx.Exec(stmt, args...); err != nil {
			tx.Rollback()
			return err
		}
		return nil
	}
	stmts := []struct {
		q    string
		args []interface{}
	}{
		{"DELETE FROM email_verifications WHERE student_id = $1", []interface{}{studentID}},
		{"DELETE FROM password_resets WHERE student_id = $1", []interface{}{studentID}},
		{"DELETE FROM link_tokens WHERE student_id = $1", []interface{}{studentID}},
		{"DELETE FROM identities WHERE student_id = $1", []interface{}{studentID}},
		{"DELETE FROM avatar_images WHERE student_id = $1", []interface{}{studentID}},
		{"DELETE FROM attempts WHERE student_id = $1", []interface{}{studentID}},
		{"DELETE FROM answer_submissions WHERE student_id = $1", []interface{}{studentID}},
		{"DELETE FROM study_answer_effects WHERE student_id = $1", []interface{}{studentID}},
		{"DELETE FROM active_sessions WHERE student_id = $1", []interface{}{studentID}},
		{"DELETE FROM sessions WHERE student_id = $1", []interface{}{studentID}},
		{"DELETE FROM implicit_credit_applications WHERE student_id = $1", []interface{}{studentID}},
		{"DELETE FROM implicit_credit_ordinals WHERE student_id = $1", []interface{}{studentID}},
		{"DELETE FROM implicit_credit_events WHERE student_id = $1", []interface{}{studentID}},
		{"DELETE FROM implicit_credit_allocations WHERE student_id = $1", []interface{}{studentID}},
		{"DELETE FROM implicit_credit_targets WHERE student_id = $1", []interface{}{studentID}},
		{"DELETE FROM concept_progress WHERE student_id = $1", []interface{}{studentID}},
		{"DELETE FROM student_topic_speed WHERE student_id = $1", []interface{}{studentID}},
		{"DELETE FROM quiz_completions WHERE student_id = $1", []interface{}{studentID}},
		{"DELETE FROM timed_quiz_items WHERE run_id IN (SELECT id FROM timed_quiz_runs WHERE student_id = $1)", []interface{}{studentID}},
		{"DELETE FROM timed_quiz_runs WHERE student_id = $1", []interface{}{studentID}},
		{"DELETE FROM timed_quiz_baselines WHERE student_id = $1", []interface{}{studentID}},
		{"DELETE FROM review_events WHERE student_id = $1", []interface{}{studentID}},
		{"DELETE FROM student_remedial_queue WHERE student_id = $1", []interface{}{studentID}},
		{"DELETE FROM server_sessions WHERE kind IN ('study_expected','study_plan','study_plan_prefs') AND (key = $1 OR key LIKE $1 || '|%')", []interface{}{studentID}},
		{"DELETE FROM question_reports WHERE reporter_id = $1", []interface{}{studentID}},
		{"DELETE FROM announcement_reads WHERE student_id = $1", []interface{}{studentID}},
		{"DELETE FROM announcement_views WHERE student_id = $1", []interface{}{studentID}},
		{"DELETE FROM ai_usage WHERE usage_key = $1", []interface{}{studentID}},
		{"DELETE FROM friendships WHERE student_a = $1 OR student_b = $1", []interface{}{studentID}},
		{"DELETE FROM friend_blocks WHERE blocker_id = $1 OR blocked_id = $1", []interface{}{studentID}},
		{"DELETE FROM students WHERE id = $1", []interface{}{studentID}},
	}
	for _, st := range stmts {
		if err := exec(st.q, st.args...); err != nil {
			return fmt.Errorf("delete account: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit account delete: %w", err)
	}
	return nil
}

// Attempts

func (s *PostgresStore) RecordAttempt(entry AttemptEntry) (int64, error) {
	return insertAttemptPostgres(s.db, entry)
}

func insertAttemptPostgres(exec interface {
	QueryRow(query string, args ...any) *sql.Row
}, entry AttemptEntry) (int64, error) {
	var id int64
	err := exec.QueryRow(`
		INSERT INTO attempts
			(session_id, student_id, concept_id, answer, expected,
			 correct, elapsed_seconds, timestamp, question, source, explanation, difficulty)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id
	`,
		entry.SessionID, entry.StudentID, entry.ConceptID,
		entry.Answer, entry.Expected,
		boolToInt(entry.Correct), entry.ElapsedSeconds,
		entry.Timestamp.Format(time.RFC3339),
		entry.Question, entry.Source, entry.Explanation, entry.Difficulty,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("record attempt: %w", err)
	}
	return id, nil
}

func insertAttemptPostgresTx(tx *sql.Tx, entry AttemptEntry) (int64, error) {
	return insertAttemptPostgres(tx, entry)
}

func (s *PostgresStore) CommitStudySubmission(c StudySubmissionCommit) (StudySubmissionOutcome, error) {
	return commitStudySubmissionTx(s.db, rebindDollars, insertAttemptPostgresTx, c)
}

func (s *PostgresStore) CommitReviewSubmission(c ReviewSubmissionCommit) (ReviewSubmissionOutcome, error) {
	return commitReviewSubmissionTx(s.db, rebindDollars, insertAttemptPostgresTx, c)
}

func (s *PostgresStore) GetStudySubmission(studentID, operation, submissionID string) (*StudySubmission, error) {
	return getStudySubmission(s.db, rebindDollars, studentID, operation, submissionID)
}

func (s *PostgresStore) ClaimStudyEffects(studentID string, now time.Time, limit int) ([]StudyAnswerEffect, error) {
	return claimStudyEffects(s.db, rebindDollars, true, studentID, now, limit)
}

func (s *PostgresStore) CompleteStudyEffect(id int64) error {
	return completeStudyEffect(s.db, rebindDollars, id, time.Now().UTC())
}

func (s *PostgresStore) FailStudyEffect(id int64, reason string, retryAt time.Time) error {
	return failStudyEffect(s.db, rebindDollars, id, reason, retryAt)
}

func (s *PostgresStore) PruneStudyEffects(olderThan time.Time, limit int) (int64, error) {
	return pruneStudyEffects(s.db, rebindDollars, olderThan, limit)
}

// GetRecentAttemptsForConcept implements the bounded per-concept window the evidence
// model reads on each graded answer. Newest last, matching BuildEvidence's expectation.
//
// Ordered descending so the LIMIT takes the most recent rows, then reversed on the way
// out: the query needs newest-first to bound itself, and the evidence model reads
// oldest-last so its recency weights increase toward the present.
func (s *PostgresStore) GetRecentAttemptsForConcept(studentID, conceptID string, limit int) ([]AttemptEntry, error) {
	if limit <= 0 {
		limit = 10
	}
	rows, err := s.db.Query(`
		SELECT id, session_id, student_id, concept_id, answer, expected,
		       correct, elapsed_seconds, timestamp, question, source, explanation, difficulty
		FROM attempts
		WHERE student_id = $1 AND concept_id = $2
		ORDER BY timestamp DESC, id DESC
		LIMIT $3
	`, studentID, conceptID, limit)

	if err != nil {
		return nil, fmt.Errorf("recent attempts for concept: %w", err)
	}
	defer rows.Close()
	var out []AttemptEntry
	for rows.Next() {
		var e AttemptEntry
		var correct int
		var ts string
		if err := rows.Scan(
			&e.ID, &e.SessionID, &e.StudentID, &e.ConceptID,
			&e.Answer, &e.Expected, &correct, &e.ElapsedSeconds, &ts,
			&e.Question, &e.Source, &e.Explanation, &e.Difficulty,
		); err != nil {
			return nil, fmt.Errorf("scan recent attempt: %w", err)
		}
		e.Correct = correct != 0
		e.Timestamp = parseAttemptTimestamp(ts)
		out = append(out, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	return out, nil
}

// UpdateAttemptCorrect flips an attempt's correctness (grading data-repair).
func (s *PostgresStore) UpdateAttemptCorrect(id int64, correct bool) error {
	if _, err := s.db.Exec(`UPDATE attempts SET correct = $1 WHERE id = $2`, boolToInt(correct), id); err != nil {
		return fmt.Errorf("update attempt correct: %w", err)
	}
	return nil
}

func (s *PostgresStore) GetAllAttempts() ([]AttemptEntry, error) {
	rows, err := s.db.Query(`
		SELECT id, session_id, student_id, concept_id, answer, expected,
		       correct, elapsed_seconds, timestamp, question, source, explanation, difficulty
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
			&e.ID, &e.SessionID, &e.StudentID, &e.ConceptID,
			&e.Answer, &e.Expected, &correct, &e.ElapsedSeconds, &ts,
			&e.Question, &e.Source, &e.Explanation, &e.Difficulty,
		); err != nil {
			return nil, fmt.Errorf("scan attempt: %w", err)
		}
		e.Correct = correct != 0
		e.Timestamp = parseAttemptTimestamp(ts)
		out = append(out, e)
	}
	return out, rows.Err()
}

func (s *PostgresStore) GetAttemptsForStudent(studentID string) ([]AttemptEntry, error) {
	rows, err := s.db.Query(`
		SELECT id, session_id, student_id, concept_id, answer, expected,
		       correct, elapsed_seconds, timestamp, question, source, explanation, difficulty
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
			&e.ID, &e.SessionID, &e.StudentID, &e.ConceptID,
			&e.Answer, &e.Expected, &correct, &e.ElapsedSeconds, &ts,
			&e.Question, &e.Source, &e.Explanation, &e.Difficulty,
		); err != nil {
			return nil, fmt.Errorf("scan attempt: %w", err)
		}
		e.Correct = correct != 0
		e.Timestamp = parseAttemptTimestamp(ts)
		out = append(out, e)
	}
	return out, rows.Err()
}

func (s *PostgresStore) GetSessionAttempts(studentID, sessionID string) ([]AttemptEntry, error) {
	rows, err := s.db.Query(`
		SELECT id, session_id, student_id, concept_id, answer, expected,
		       correct, elapsed_seconds, timestamp, question, source, explanation, difficulty
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
			&e.ID, &e.SessionID, &e.StudentID, &e.ConceptID,
			&e.Answer, &e.Expected, &correct, &e.ElapsedSeconds, &ts,
			&e.Question, &e.Source, &e.Explanation, &e.Difficulty,
		); err != nil {
			return nil, fmt.Errorf("scan attempt: %w", err)
		}
		e.Correct = correct != 0
		e.Timestamp = parseAttemptTimestamp(ts)
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
		SELECT s.id, s.name, s.username, s.avatar_url, s.settings,
			COUNT(CASE WHEN cp.status = 'MASTERED' THEN 1 END),
			COUNT(CASE WHEN cp.status = 'MASTERED' AND cp.mastered_at >= $1 THEN 1 END)
		FROM students s
		LEFT JOIN concept_progress cp ON cp.student_id = s.id
		GROUP BY s.id, s.name, s.username, s.avatar_url, s.settings
		ORDER BY 7 DESC, 6 DESC
	`, monday.Format(time.RFC3339))
	if err != nil {
		return nil, fmt.Errorf("get weekly leaderboard: %w", err)
	}
	defer rows.Close()

	var out []LeaderboardRow
	for rows.Next() {
		var r LeaderboardRow
		var username, avatarURL, settings sql.NullString
		if err := rows.Scan(&r.StudentID, &r.Name, &username, &avatarURL, &settings, &r.TotalMastered, &r.WeeklyMastered); err != nil {
			return nil, fmt.Errorf("scan leaderboard: %w", err)
		}
		r.Username = username.String
		r.AvatarURL = avatarURL.String
		r.AvatarSettings = settings.String
		out = append(out, r)
	}
	return out, rows.Err()
}

// DailyActivity

func (s *PostgresStore) GetDailyActivity(studentID string, days int) ([]DailyActivity, error) {
	// For postgres, use now() - interval
	rows, err := s.db.Query(`
		SELECT date(timestamp::timestamptz)::text as day,
		       COUNT(*) as total,
		       COALESCE(SUM(correct), 0) as correct_count,
		       COALESCE(STRING_AGG(DISTINCT concept_id, ','), '') as concepts
		FROM attempts
		WHERE student_id = $1 AND timestamp::timestamptz >= now() - ($2 || ' days')::interval
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

// Batch 1: Quiz 150 XP gate completions.

func (s *PostgresStore) RecordQuizCompletion(studentID string, xpTotal int) error {
	if _, err := s.db.Exec(`INSERT INTO quiz_completions (student_id, completed_at, xp_total) VALUES ($1, now()::text, $2)`, studentID, xpTotal); err != nil {
		return fmt.Errorf("record quiz completion: %w", err)
	}
	return nil
}

func (s *PostgresStore) LastQuizCompletion(studentID string) (*QuizCompletion, error) {
	row := s.db.QueryRow(`SELECT student_id, completed_at, xp_total FROM quiz_completions WHERE student_id = $1 ORDER BY id DESC LIMIT 1`, studentID)
	var qc QuizCompletion
	var at sql.NullString
	if err := row.Scan(&qc.StudentID, &at, &qc.XPTotal); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("last quiz completion: %w", err)
	}
	qc.CompletedAt = at.String
	return &qc, nil
}

// Question reports (user complaints about questions/explanations/lessons).

func (s *PostgresStore) CreateReport(r QuestionReport) (int64, error) {
	if r.Status == "" {
		r.Status = "open"
	}
	created := r.CreatedAt.UTC().Format(time.RFC3339)
	if r.CreatedAt.IsZero() {
		created = time.Now().UTC().Format(time.RFC3339)
	}
	var id int64
	err := s.db.QueryRow(`
		INSERT INTO question_reports
			(reporter_id, concept_id, kind, question, expected, explanation,
			 lesson_id, source, session_id, attempt_id, answer_given, grading_type,
			 answer_hint, answer_input_mode, reason, detail, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
		RETURNING id
	`,
		r.ReporterID, r.ConceptID, r.Kind, r.Question, r.Expected, r.Explanation,
		r.LessonID, r.Source, r.SessionID, r.AttemptID,
		r.AnswerGiven, r.GradingType, r.AnswerHint, r.AnswerInputMode,
		r.Reason, r.Detail, r.Status,
		created,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create report: %w", err)
	}
	return id, nil
}

// ImportReport mirrors the SQLite implementation: the origin id is what makes
// a repeated import a no-op, and the row is otherwise an ordinary report.
func (s *PostgresStore) ImportReport(originID string, r QuestionReport) (bool, error) {
	originID = strings.TrimSpace(originID)
	if originID == "" {
		return false, fmt.Errorf("import report: origin id is required")
	}
	if r.Status == "" {
		r.Status = "open"
	}
	created := r.CreatedAt.UTC().Format(time.RFC3339)
	if r.CreatedAt.IsZero() {
		created = time.Now().UTC().Format(time.RFC3339)
	}
	res, err := s.db.Exec(`
		INSERT INTO question_reports
			(reporter_id, concept_id, kind, question, expected, explanation,
			 lesson_id, source, session_id, attempt_id, answer_given, grading_type,
			 answer_hint, answer_input_mode, reason, detail, status, created_at,
			 origin_report_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
		ON CONFLICT DO NOTHING
	`,
		r.ReporterID, r.ConceptID, r.Kind, r.Question, r.Expected, r.Explanation,
		r.LessonID, r.Source, r.SessionID, r.AttemptID,
		r.AnswerGiven, r.GradingType, r.AnswerHint, r.AnswerInputMode,
		r.Reason, r.Detail, r.Status,
		created, originID,
	)
	if err != nil {
		return false, fmt.Errorf("import report: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("import report: %w", err)
	}
	return n > 0, nil
}

func (s *PostgresStore) ListReports(status string, limit, offset int) ([]QuestionReport, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	var rows *sql.Rows
	var err error
	if status == "" || status == "all" {
		rows, err = s.db.Query(`
			SELECT id, reporter_id,
			       COALESCE((SELECT username FROM students WHERE students.id = question_reports.reporter_id), ''),
			       COALESCE((SELECT name FROM students WHERE students.id = question_reports.reporter_id), ''),
			       concept_id, kind, question, expected, explanation,
			       lesson_id, source, session_id, attempt_id,
			       answer_given, grading_type, answer_hint, answer_input_mode,
			       reason, detail, status,
			       resolution, resolved_by, resolved_at, created_at, origin_report_id,
			       EXISTS(SELECT 1 FROM report_images WHERE report_images.report_id = question_reports.id)
			FROM question_reports
			ORDER BY created_at DESC
			LIMIT $1 OFFSET $2
		`, limit, offset)
	} else {
		rows, err = s.db.Query(`
			SELECT id, reporter_id,
			       COALESCE((SELECT username FROM students WHERE students.id = question_reports.reporter_id), ''),
			       COALESCE((SELECT name FROM students WHERE students.id = question_reports.reporter_id), ''),
			       concept_id, kind, question, expected, explanation,
			       lesson_id, source, session_id, attempt_id,
			       answer_given, grading_type, answer_hint, answer_input_mode,
			       reason, detail, status,
			       resolution, resolved_by, resolved_at, created_at, origin_report_id,
			       EXISTS(SELECT 1 FROM report_images WHERE report_images.report_id = question_reports.id)
			FROM question_reports
			WHERE status = $1
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3
		`, status, limit, offset)
	}
	if err != nil {
		return nil, fmt.Errorf("list reports: %w", err)
	}
	defer rows.Close()
	var out []QuestionReport
	for rows.Next() {
		var r QuestionReport
		var created string
		if err := rows.Scan(
			&r.ID, &r.ReporterID, &r.ReporterUsername, &r.ReporterName,
			&r.ConceptID, &r.Kind, &r.Question, &r.Expected,
			&r.Explanation, &r.LessonID, &r.Source, &r.SessionID, &r.AttemptID,
			&r.AnswerGiven, &r.GradingType, &r.AnswerHint, &r.AnswerInputMode,
			&r.Reason, &r.Detail, &r.Status, &r.Resolution, &r.ResolvedBy, &r.ResolvedAt, &created,
			&r.OriginReportID, &r.HasImage,
		); err != nil {
			return nil, fmt.Errorf("scan report: %w", err)
		}
		if t, err := time.Parse(time.RFC3339, created); err == nil {
			r.CreatedAt = t
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

func (s *PostgresStore) UpdateReportStatus(id int64, status string) error {
	res, err := s.db.Exec(`UPDATE question_reports SET status = $1 WHERE id = $2`, status, id)
	if err != nil {
		return fmt.Errorf("update report: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("report not found")
	}
	return nil
}

// ── Administration (Postgres) ───────────────────────────────────────────────

const pgAdminUserColumns = `id, name, username, email, role, created_at, (password_hash IS NOT NULL AND password_hash != '')`

func (s *PostgresStore) GetStudentRole(studentID string) (Role, error) {
	var raw string
	err := s.db.QueryRow("SELECT role FROM students WHERE id = $1", studentID).Scan(&raw)
	if err == sql.ErrNoRows {
		// See the SQLite store: a role lookup runs on every authenticated request, so a deleted
		// id must read as "no privilege" rather than fail the middleware.
		return RoleStudent, nil
	}
	if err != nil {
		return RoleStudent, err
	}
	return NormalizeRole(raw), nil
}

func (s *PostgresStore) SetStudentRole(studentID string, role Role) error {
	res, err := s.db.Exec("UPDATE students SET role = $1 WHERE id = $2", string(role), studentID)
	if err != nil {
		return fmt.Errorf("set student role: %w", err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// DemoteAdminSafely mirrors the SQLite implementation statement for statement, so the two
// stores cannot disagree about what the last-admin rule is. See the SQLite comment for why the
// count is inside the UPDATE rather than a read before it.
func (s *PostgresStore) DemoteAdminSafely(studentID string) error {
	res, err := s.db.Exec(
		`UPDATE students SET role = 'student'
		 WHERE id = $1 AND role = 'admin' AND (SELECT COUNT(*) FROM students WHERE role = 'admin') > 1`,
		studentID)
	if err != nil {
		return fmt.Errorf("demote admin: %w", err)
	}
	n, _ := res.RowsAffected()
	if n > 0 {
		return nil
	}
	role, err := s.GetStudentRole(studentID)
	if err != nil {
		return err
	}
	if role != RoleAdmin {
		return sql.ErrNoRows
	}
	return ErrLastAdmin
}

func (s *PostgresStore) CountAdmins() (int, error) {
	var n int
	err := s.db.QueryRow("SELECT COUNT(*) FROM students WHERE role = 'admin'").Scan(&n)
	return n, err
}

func (s *PostgresStore) CountAdminUsers() (int, error) {
	var n int
	err := s.db.QueryRow("SELECT COUNT(*) FROM students").Scan(&n)
	return n, err
}

func (s *PostgresStore) CountStoredQuestions() (int, error) {
	var n int
	err := s.db.QueryRow("SELECT COUNT(*) FROM questions").Scan(&n)
	return n, err
}

// adminQueryerPG is the same indirection the SQLite store has, so its projection read can run
// inside SetRoleAudited's transaction instead of on the pool. The two stores are kept at
// parity deliberately: a rule that only one of them enforces is a rule a Postgres deployment
// does not have.
type adminQueryerPG interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

// SearchAdminUsers escapes LIKE metacharacters the same way the SQLite store does. Postgres
// would also treat a backslash as the default escape character, but it is spelled explicitly
// here so the two stores cannot disagree about which character escapes what.
func (s *PostgresStore) SearchAdminUsers(query string, limit, offset int) ([]AdminUser, error) {
	return searchAdminUsersPG(s.db, query, limit, offset)
}

func searchAdminUsersPG(q adminQueryerPG, query string, limit, offset int) ([]AdminUser, error) {
	term := strings.TrimSpace(query)
	var (
		query_ string
		args   []interface{}
	)
	if term == "" {
		query_ = `SELECT ` + pgAdminUserColumns + ` FROM students ORDER BY created_at DESC, id DESC LIMIT $1 OFFSET $2`
		args = []interface{}{limit, offset}
	} else {
		needle := "%" + escapeLike(term) + "%"
		query_ = `SELECT ` + pgAdminUserColumns + ` FROM students
		 WHERE id LIKE $1 ESCAPE '\' OR email LIKE $2 ESCAPE '\'
		    OR username LIKE $3 ESCAPE '\' OR name LIKE $4 ESCAPE '\'
		 ORDER BY created_at DESC, id DESC LIMIT $5 OFFSET $6`
		args = []interface{}{needle, needle, needle, needle, limit, offset}
	}

	rows, err := q.Query(query_, args...)
	if err != nil {
		return nil, fmt.Errorf("search admin users: %w", err)
	}
	defer rows.Close()
	out := []AdminUser{}
	for rows.Next() {
		var u AdminUser
		var name, username, email, role, createdAt sql.NullString
		var hasPw bool
		if err := rows.Scan(&u.ID, &name, &username, &email, &role, &createdAt, &hasPw); err != nil {
			return nil, err
		}
		u.Name = name.String
		u.Username = username.String
		u.Email = email.String
		u.Role = NormalizeRole(role.String)
		u.CreatedAt = createdAt.String
		u.HasPassword = hasPw
		out = append(out, u)
	}
	return out, rows.Err()
}

func (s *PostgresStore) RecordAdminAudit(rec AdminAuditRecord) (int64, error) {
	// lib/pq does not implement LastInsertId, so the id has to come back from
	// the insert itself. Without this every audit write failed on Postgres.
	var id int64
	err := s.db.QueryRow(
		`INSERT INTO admin_audit (actor_id, action, entity_type, entity_id, before_json, after_json, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		rec.ActorID, rec.Action, rec.EntityType, rec.EntityID,
		rec.BeforeJSON, rec.AfterJSON, formatReportTime(parseReportTime(rec.CreatedAt))).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("record admin audit: %w", err)
	}
	return id, nil
}

func (s *PostgresStore) ListAdminAudit(limit, offset int) ([]AdminAuditRecord, error) {
	rows, err := s.db.Query(
		`SELECT id, actor_id, action, entity_type, entity_id, before_json, after_json, created_at
		 FROM admin_audit ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list admin audit: %w", err)
	}
	defer rows.Close()
	out := []AdminAuditRecord{}
	for rows.Next() {
		var r AdminAuditRecord
		var createdAt sql.NullString
		if err := rows.Scan(&r.ID, &r.ActorID, &r.Action, &r.EntityType, &r.EntityID,
			&r.BeforeJSON, &r.AfterJSON, &createdAt); err != nil {
			return nil, err
		}
		r.CreatedAt = createdAt.String
		out = append(out, r)
	}
	return out, rows.Err()
}

// adminUserByIDPG is the exact-id read, mirroring adminUserByIDSQLite.
func adminUserByIDPG(q adminQueryerPG, id string) (*AdminUser, error) {
	rows, err := searchAdminUsersPG(q, id, 200, 0)
	if err != nil {
		return nil, err
	}
	for i := range rows {
		if rows[i].ID == id {
			return &rows[i], nil
		}
	}
	return nil, nil
}

// SetRoleAudited mirrors the SQLite implementation statement for statement, including the
// transaction boundary, so a deployment on Postgres gets the same guarantee: a role change and
// its audit event are one unit of work.
func (s *PostgresStore) SetRoleAudited(actorID, studentID string, role Role) (RoleChange, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return RoleChange{}, fmt.Errorf("begin role change: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	before, err := adminUserByIDPG(tx, studentID)
	if err != nil {
		return RoleChange{}, fmt.Errorf("read user before: %w", err)
	}
	if before == nil {
		return RoleChange{}, sql.ErrNoRows
	}
	if before.Role == role {
		return RoleChange{Before: *before, After: *before, Changed: false}, ErrRoleUnchanged
	}

	action := ActionAdminDemoteUser
	if role.AtLeast(before.Role) {
		action = ActionAdminPromoteUser
	}
	// Removing an owner is the one change that can empty the top of the hierarchy, so it carries
	// the guard. Mirrors SQLite; see that comment for why the count is inside the UPDATE.
	if before.Role == RoleOwner && role != RoleOwner {
		res, err := tx.Exec(
			`UPDATE students SET role = $1
			 WHERE id = $2 AND role = 'owner' AND (SELECT COUNT(*) FROM students WHERE role = 'owner') > 1`,
			string(role), studentID)
		if err != nil {
			return RoleChange{}, fmt.Errorf("demote owner: %w", err)
		}
		if n, _ := res.RowsAffected(); n == 0 {
			return RoleChange{}, ErrLastOwner
		}
	} else {
		if _, err := tx.Exec("UPDATE students SET role = $1 WHERE id = $2", string(role), studentID); err != nil {
			return RoleChange{}, fmt.Errorf("set role: %w", err)
		}
	}

	after, err := adminUserByIDPG(tx, studentID)
	if err != nil {
		return RoleChange{}, fmt.Errorf("read user after: %w", err)
	}
	if after == nil {
		return RoleChange{}, sql.ErrNoRows
	}

	beforeJSON, err := json.Marshal(before)
	if err != nil {
		return RoleChange{}, fmt.Errorf("encode before: %w", err)
	}
	afterJSON, err := json.Marshal(after)
	if err != nil {
		return RoleChange{}, fmt.Errorf("encode after: %w", err)
	}
	if _, err := tx.Exec(
		`INSERT INTO admin_audit (actor_id, action, entity_type, entity_id, before_json, after_json, created_at)
		 VALUES ($1, $2, 'user', $3, $4, $5, $6)`,
		actorID, action, studentID, string(beforeJSON), string(afterJSON),
		time.Now().UTC().Format(time.RFC3339)); err != nil {
		return RoleChange{}, fmt.Errorf("record audit: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return RoleChange{}, fmt.Errorf("commit role change: %w", err)
	}
	return RoleChange{Before: *before, After: *after, Changed: true, Action: action}, nil
}

// GetAdminUser reads one non-sensitive projection, or nil when the id does not exist.
func (s *PostgresStore) GetAdminUser(studentID string) (*AdminUser, error) {
	return adminUserByIDPG(s.db, studentID)
}

// RecordXPEvent appends to the ledger only; see the SQLite implementation for
// why the counters must not move.
func (s *PostgresStore) RecordXPEvent(studentID string, amount int, reason, eventKey string) (bool, error) {
	if studentID == "" || eventKey == "" {
		return false, fmt.Errorf("xp event needs a student and an event key")
	}
	res, err := s.db.Exec(`INSERT INTO xp_events (student_id, amount, reason, event_key, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (student_id, event_key) DO NOTHING`,
		studentID, amount, reason, eventKey, time.Now().UTC().UnixNano())
	if err != nil {
		return false, fmt.Errorf("record xp event: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("record xp event: %w", err)
	}
	return n > 0, nil
}

func (s *PostgresStore) AttemptsSince(since time.Time) ([]AttemptRow, error) {
	rows, err := s.db.Query(
		`SELECT id, student_id, concept_id, correct, source
		 FROM attempts
		 WHERE timestamp::timestamptz >= $1
		 ORDER BY id ASC`, since.UTC())
	if err != nil {
		return nil, fmt.Errorf("read attempts since: %w", err)
	}
	defer rows.Close()
	var out []AttemptRow
	for rows.Next() {
		var r AttemptRow
		var correct int
		if err := rows.Scan(&r.ID, &r.StudentID, &r.ConceptID, &correct, &r.Source); err != nil {
			return nil, fmt.Errorf("scan attempt: %w", err)
		}
		r.Correct = correct != 0
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read attempts since: %w", err)
	}
	return out, nil
}

// RecordDiagnosticAttempt mirrors the SQLite method; see the schema comment.
func (s *PostgresStore) RecordDiagnosticAttempt(attemptID, response string) error {
	if attemptID == "" || response == "" {
		return fmt.Errorf("diagnostic attempt: empty id or response")
	}
	_, err := s.db.Exec(
		`INSERT INTO diagnostic_answers (attempt_id, response, created_at) VALUES ($1, $2, $3)
		 ON CONFLICT (attempt_id) DO NOTHING`,
		attemptID, response, time.Now().UTC().UnixNano())
	if err != nil {
		return fmt.Errorf("record diagnostic attempt: %w", err)
	}
	return nil
}

func (s *PostgresStore) DiagnosticAttemptResponse(attemptID string) (string, bool, error) {
	var response string
	err := s.db.QueryRow(`SELECT response FROM diagnostic_answers WHERE attempt_id = $1`, attemptID).Scan(&response)
	if err == sql.ErrNoRows {
		return "", false, nil
	}
	if err != nil {
		return "", false, fmt.Errorf("read diagnostic attempt: %w", err)
	}
	return response, true, nil
}

// SeedFSRSFromSM2 backfills FSRS state from existing SM-2 intervals for unseeded rows.
// It is idempotent: only rows where fsrs_stability is NULL are seeded.
func (s *PostgresStore) SeedFSRSFromSM2(now time.Time) (int, error) {
	rows, err := s.db.Query(`
		SELECT student_id, concept_id, sm2_interval, next_review_due, last_reviewed
		FROM concept_progress
		WHERE fsrs_stability IS NULL AND sm2_interval >= 1 AND next_review_due IS NOT NULL AND next_review_due != ''
		ORDER BY student_id, concept_id
	`)
	if err != nil {
		return 0, fmt.Errorf("query unseeded fsrs rows: %w", err)
	}
	defer rows.Close()

	type toSeed struct {
		studentID string
		conceptID string
		interval  int
		nextDue   string
		lastRev   sql.NullString
	}
	var batch []toSeed
	for rows.Next() {
		var item toSeed
		if err := rows.Scan(&item.studentID, &item.conceptID, &item.interval, &item.nextDue, &item.lastRev); err != nil {
			return 0, fmt.Errorf("scan unseeded fsrs row: %w", err)
		}
		batch = append(batch, item)
	}
	if err := rows.Err(); err != nil {
		return 0, fmt.Errorf("iterate unseeded fsrs rows: %w", err)
	}

	params := fsrs.DefaultParameters()
	defaultDiff := params.InitialDifficulty(fsrs.Good)

	stmt, err := s.db.Prepare(`
		UPDATE concept_progress
		SET fsrs_stability = $1, fsrs_difficulty = $2, fsrs_last_review = $3
		WHERE student_id = $4 AND concept_id = $5 AND fsrs_stability IS NULL
	`)
	if err != nil {
		return 0, fmt.Errorf("prepare update fsrs seed: %w", err)
	}
	defer stmt.Close()

	count := 0
	for _, item := range batch {
		due, err := time.Parse(time.RFC3339, item.nextDue)
		if err != nil {
			due, err = time.Parse("2006-01-02 15:04:05", item.nextDue)
		}
		if err != nil {
			continue
		}
		seeded, err := fsrs.SeedFromInterval(item.interval, due, now)
		if err != nil {
			continue
		}
		lastReview := now.Format(time.RFC3339)
		if item.lastRev.Valid && item.lastRev.String != "" {
			lastReview = item.lastRev.String
		}
		res, err := stmt.Exec(seeded.Stability, defaultDiff, lastReview, item.studentID, item.conceptID)
		if err != nil {
			return count, fmt.Errorf("seed fsrs progress: %w", err)
		}
		affected, _ := res.RowsAffected()
		count += int(affected)
	}

	return count, nil
}

// EnqueueRemedial mirrors the SQLite method; see the schema comment.
func (s *PostgresStore) EnqueueRemedial(studentID string, conceptIDs ...string) error {
	if len(conceptIDs) == 0 {
		return nil
	}
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("begin enqueue remedial: %w", err)
	}
	defer tx.Rollback()

	rows, err := tx.Query(`SELECT concept_id, priority FROM student_remedial_queue WHERE student_id = $1 ORDER BY priority ASC, created_at ASC FOR UPDATE`, studentID)
	if err != nil {
		return fmt.Errorf("query remedial queue: %w", err)
	}
	seen := make(map[string]bool)
	maxPriority := -1
	count := 0
	for rows.Next() {
		var cid string
		var pri int
		if err := rows.Scan(&cid, &pri); err != nil {
			rows.Close()
			return fmt.Errorf("scan remedial queue: %w", err)
		}
		seen[cid] = true
		if pri > maxPriority {
			maxPriority = pri
		}
		count++
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate remedial queue: %w", err)
	}

	now := time.Now().Unix()
	stmt, err := tx.Prepare(`INSERT INTO student_remedial_queue (student_id, concept_id, priority, created_at) VALUES ($1, $2, $3, $4) ON CONFLICT (student_id, concept_id) DO NOTHING`)
	if err != nil {
		return fmt.Errorf("prepare insert remedial: %w", err)
	}
	defer stmt.Close()

	for _, cid := range conceptIDs {
		if seen[cid] || count >= MaxRemedialQueue {
			continue
		}
		seen[cid] = true
		maxPriority++
		count++
		if _, err := stmt.Exec(studentID, cid, maxPriority, now); err != nil {
			return fmt.Errorf("insert remedial queue: %w", err)
		}
	}

	return tx.Commit()
}

// GetRemedialQueue mirrors the SQLite method; see the schema comment.
func (s *PostgresStore) GetRemedialQueue(studentID string) ([]string, error) {
	rows, err := s.db.Query(`SELECT concept_id FROM student_remedial_queue WHERE student_id = $1 ORDER BY priority ASC, created_at ASC`, studentID)
	if err != nil {
		return nil, fmt.Errorf("get remedial queue: %w", err)
	}
	defer rows.Close()

	var result []string
	for rows.Next() {
		var cid string
		if err := rows.Scan(&cid); err != nil {
			return nil, fmt.Errorf("scan remedial queue item: %w", err)
		}
		result = append(result, cid)
	}
	return result, rows.Err()
}

// RemoveRemedial mirrors the SQLite method; see the schema comment.
func (s *PostgresStore) RemoveRemedial(studentID, conceptID string) error {
	if _, err := s.db.Exec(`DELETE FROM student_remedial_queue WHERE student_id = $1 AND concept_id = $2`, studentID, conceptID); err != nil {
		return fmt.Errorf("remove remedial: %w", err)
	}
	return nil
}

// ClearRemedialQueue mirrors the SQLite method; see the schema comment.
func (s *PostgresStore) ClearRemedialQueue(studentID string) error {
	if _, err := s.db.Exec(`DELETE FROM student_remedial_queue WHERE student_id = $1`, studentID); err != nil {
		return fmt.Errorf("clear remedial queue: %w", err)
	}
	return nil
}

// AssessmentAttempts returns quiz and diagnostic answers in a window, newest
// first. Only quiz and diagnostic sources are assessments; practice and review
// answers are activity.
func (s *PostgresStore) AssessmentAttempts(studentID string, since, until time.Time, limit int) ([]AttemptEntry, error) {
	rows, err := s.db.Query(`
		SELECT id, session_id, student_id, concept_id, answer, expected, correct,
		       elapsed_seconds, timestamp, question, source, explanation, difficulty
		FROM attempts
		WHERE student_id = $1
		  AND source IN ('quiz', 'diagnostic')
		  AND timestamp::timestamptz >= $2
		  AND timestamp::timestamptz < $3
		ORDER BY timestamp::timestamptz DESC, id DESC
		LIMIT $4`, studentID, since.UTC(), until.UTC(), limit)
	if err != nil {
		return nil, fmt.Errorf("assessment attempts: %w", err)
	}
	defer rows.Close()

	var out []AttemptEntry
	for rows.Next() {
		var e AttemptEntry
		var correct int
		var ts string
		var difficulty sql.NullFloat64
		if err := rows.Scan(&e.ID, &e.SessionID, &e.StudentID, &e.ConceptID, &e.Answer,
			&e.Expected, &correct, &e.ElapsedSeconds, &ts, &e.Question, &e.Source,
			&e.Explanation, &difficulty); err != nil {
			return nil, fmt.Errorf("scan assessment attempt: %w", err)
		}
		e.Correct = correct != 0
		e.Timestamp = parseAttemptTimestamp(ts)
		if difficulty.Valid {
			v := difficulty.Float64
			e.Difficulty = &v
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

// XPInRange sums the ledger in [since, until). xp_events.created_at is UnixNano
// in both stores, so the comparison is numeric and timezone independent.
func (s *PostgresStore) XPInRange(studentID string, since, until time.Time) (int, error) {
	var total int
	err := s.db.QueryRow(`
		SELECT COALESCE(SUM(amount), 0)
		FROM xp_events
		WHERE student_id = $1 AND created_at >= $2 AND created_at < $3`,
		studentID, since.UTC().UnixNano(), until.UTC().UnixNano()).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("xp in range: %w", err)
	}
	return total, nil
}

const postgresLearningEvents = `
	WITH events AS (
		SELECT 'answer' AS kind, CAST(id AS TEXT) AS ref,
		       EXTRACT(EPOCH FROM timestamp::timestamptz)::bigint AS at,
		       concept_id, source, correct AS value, 0 AS xp
		FROM attempts
		WHERE student_id = $1 AND timestamp IS NOT NULL AND timestamp <> ''
		UNION ALL
		SELECT 'mastery', concept_id, EXTRACT(EPOCH FROM mastered_at::timestamptz)::bigint,
		       concept_id, '', 0, 0
		FROM concept_progress
		WHERE student_id = $2 AND status = 'MASTERED'
		  AND mastered_at IS NOT NULL AND mastered_at <> ''
		UNION ALL
		SELECT 'quiz', CAST(id AS TEXT), EXTRACT(EPOCH FROM completed_at::timestamptz)::bigint,
		       '', 'quiz', 0, COALESCE(xp_total, 0)
		FROM quiz_completions
		WHERE student_id = $3
		UNION ALL
		SELECT 'diagnostic', 'diagnostic', EXTRACT(EPOCH FROM diagnostic_completed_at::timestamptz)::bigint,
		       '', 'diagnostic', 0, 0
		FROM students
		WHERE id = $4 AND diagnostic_completed_at IS NOT NULL AND diagnostic_completed_at <> ''
	)
	SELECT kind, ref, at, concept_id, source, value, xp
	FROM events
	WHERE at IS NOT NULL
	ORDER BY at DESC, kind ASC, ref DESC
	LIMIT $5 OFFSET $6`

// LearningEventsPage unions the tables that carry dated events. A concept row
// without mastered_at is a current state, not a dated transition, so it is not
// turned into one here.
func (s *PostgresStore) LearningEventsPage(studentID string, limit, offset int) ([]LearningEvent, int, error) {
	var total int
	err := s.db.QueryRow(`
		WITH events AS (
			SELECT EXTRACT(EPOCH FROM timestamp::timestamptz)::bigint AS at FROM attempts
			WHERE student_id = $1 AND timestamp IS NOT NULL AND timestamp <> ''
			UNION ALL
			SELECT EXTRACT(EPOCH FROM mastered_at::timestamptz)::bigint FROM concept_progress
			WHERE student_id = $2 AND status = 'MASTERED'
			  AND mastered_at IS NOT NULL AND mastered_at <> ''
			UNION ALL
			SELECT EXTRACT(EPOCH FROM completed_at::timestamptz)::bigint FROM quiz_completions
			WHERE student_id = $3
			UNION ALL
			SELECT EXTRACT(EPOCH FROM diagnostic_completed_at::timestamptz)::bigint FROM students
			WHERE id = $4 AND diagnostic_completed_at IS NOT NULL AND diagnostic_completed_at <> ''
		)
		SELECT COUNT(*) FROM events WHERE at IS NOT NULL`,
		studentID, studentID, studentID, studentID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count learning events: %w", err)
	}

	rows, err := s.db.Query(postgresLearningEvents, studentID, studentID, studentID, studentID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("learning events: %w", err)
	}
	defer rows.Close()

	events := []LearningEvent{}
	for rows.Next() {
		var e LearningEvent
		var at int64
		var correct int
		if err := rows.Scan(&e.Kind, &e.Ref, &at, &e.ConceptID, &e.Source, &correct, &e.XP); err != nil {
			return nil, 0, fmt.Errorf("scan learning event: %w", err)
		}
		e.At = time.Unix(at, 0).UTC()
		e.Correct = correct != 0
		events = append(events, e)
	}
	return events, total, rows.Err()
}

// QuizCompletions returns mastery-check completions at or after since, newest
// first. The stored xp_total is the lifetime snapshot at completion, not a
// score, so callers must not derive a percentage from it.
func (s *PostgresStore) QuizCompletions(studentID string, since time.Time, limit int) ([]QuizCompletion, error) {
	rows, err := s.db.Query(`
		SELECT student_id, completed_at, COALESCE(xp_total, 0)
		FROM quiz_completions
		WHERE student_id = $1
		ORDER BY id DESC
		LIMIT $2`, studentID, limit)
	if err != nil {
		return nil, fmt.Errorf("quiz completions: %w", err)
	}
	defer rows.Close()

	var out []QuizCompletion
	for rows.Next() {
		var q QuizCompletion
		if err := rows.Scan(&q.StudentID, &q.CompletedAt, &q.XPTotal); err != nil {
			return nil, fmt.Errorf("scan quiz completion: %w", err)
		}
		if at, ok := parseQuizCompletedAt(q.CompletedAt); ok && at.Before(since) {
			continue
		}
		out = append(out, q)
	}
	return out, rows.Err()
}
