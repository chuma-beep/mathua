package storage

import (
	"database/sql"
	"encoding/json"
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
	// Fix 7 (revised): small pool instead of a single connection. WAL
	// allows one writer plus concurrent readers; concurrent writers
	// serialize inside SQLite and ride the 5s busy_timeout above. A pool
	// of 8 lets reads interleave with the write queue instead of
	// queueing behind it — loadtest-verified against MaxOpenConns(1).
	//
	// This pool and the ":memory:" DSN do not mix, and the failure is silent until something
	// runs two things at once. Each connection to ":memory:" opens its OWN database, so a second
	// concurrent statement gets an empty one: "no such table: students", and any count taken
	// through it reads as zero. Sequentially the pool reuses its one idle connection and every
	// test passes, which is what makes it hard to find — a concurrent test on a :memory: store
	// reports whatever the empty database contains rather than what the code did.
	//
	// Production uses a file, and so does any test that needs real concurrency; see
	// TestConcurrentDemotionsCannotReachZeroAdmins, which uses t.TempDir() for this reason and
	// says so. Fixing it properly would mean opening ":memory:" as "file::memory:?cache=shared"
	// with _txlock, which changes behaviour for every existing in-memory test to buy nothing.
	db.SetMaxOpenConns(8)
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
		"ALTER TABLE students ADD COLUMN daily_xp_goal INTEGER NOT NULL DEFAULT 30",
		"ALTER TABLE students ADD COLUMN settings TEXT NOT NULL DEFAULT '{}'",
		"ALTER TABLE students ADD COLUMN league TEXT NOT NULL DEFAULT 'bronze'",
		"ALTER TABLE students ADD COLUMN league_week TEXT NOT NULL DEFAULT ''",
		"ALTER TABLE students ADD COLUMN league_moved INTEGER NOT NULL DEFAULT 0",
		"ALTER TABLE students ADD COLUMN share_token TEXT NOT NULL DEFAULT ''",
		"CREATE INDEX IF NOT EXISTS idx_students_username ON students(username)",
		"CREATE INDEX IF NOT EXISTS idx_students_share ON students(share_token)",
		"ALTER TABLE students ADD COLUMN email TEXT NOT NULL DEFAULT ''",
		"ALTER TABLE students ADD COLUMN google_id TEXT NOT NULL DEFAULT ''",
		"ALTER TABLE students ADD COLUMN avatar_url TEXT NOT NULL DEFAULT ''",
		"ALTER TABLE students ADD COLUMN email_verified INTEGER NOT NULL DEFAULT 0",
		`INSERT OR IGNORE INTO identities (provider, provider_id, student_id, email, email_verified)
			SELECT 'google', google_id, id, email, 1 FROM students
			WHERE google_id IS NOT NULL AND google_id != ''`,
		"CREATE INDEX IF NOT EXISTS idx_students_email ON students(email)",
		"CREATE INDEX IF NOT EXISTS idx_students_google_id ON students(google_id)",
		"ALTER TABLE concept_progress ADD COLUMN weakness_score REAL NOT NULL DEFAULT 0",
		"ALTER TABLE active_sessions ADD COLUMN last_concept_id TEXT NOT NULL DEFAULT ''",
		"ALTER TABLE active_sessions ADD COLUMN session_review INTEGER NOT NULL DEFAULT 0",
		"ALTER TABLE active_sessions ADD COLUMN session_new INTEGER NOT NULL DEFAULT 0",
		"ALTER TABLE attempts ADD COLUMN question TEXT NOT NULL DEFAULT ''",
		"ALTER TABLE attempts ADD COLUMN source TEXT NOT NULL DEFAULT ''",
		"ALTER TABLE attempts ADD COLUMN explanation TEXT NOT NULL DEFAULT ''",
		// Nullable on purpose: existing rows have no difficulty, and backfilling a
		// plausible-looking number would be fabricating evidence.
		"ALTER TABLE attempts ADD COLUMN difficulty REAL",
		"ALTER TABLE active_sessions ADD COLUMN difficulty REAL",
		// The administrative role. NOT NULL DEFAULT 'student' means every existing account is
		// migrated to a learner, which is the only direction that can be migrated safely: there
		// is no evidence in a pre-role database about who was administering anything, because
		// nothing in it could.
		"ALTER TABLE students ADD COLUMN role TEXT NOT NULL DEFAULT 'student'",
		"CREATE INDEX IF NOT EXISTS idx_students_role ON students(role)",
	}
	for _, m := range migrations {
		if _, err := db.Exec(m); err != nil {
			if !isDuplicateColumn(err) {
				return err
			}
		}
	}
	// Back the app-level one-email-one-account 409s with a partial unique
	// index (case-insensitive, legacy '' rows excluded). Retried each boot
	// until clean so dirty databases still start.
	if err := ensureEmailUniqueIndexSQLite(db); err != nil {
		return err
	}
	// One-time data backfills (drift fixes), version-guarded:
	// v1 daily goal 150 → 30 (MA 20-40); v2 username lowercasing except
	// colliding groups (same lower form twice — manual rename, never merge).
	// Fix 7: both run exactly once via schema_migrations, not every boot.
	dataMigrations := map[int][]string{
		1: {"UPDATE students SET daily_xp_goal = 30 WHERE daily_xp_goal = 150"},
		2: {`UPDATE students SET username = lower(username)
			WHERE lower(username) NOT IN (
				SELECT lower(username) FROM students WHERE username IS NOT NULL AND username != ''
				GROUP BY lower(username) HAVING COUNT(*) > 1
			)`},
		// v3 small-awards economy: proportional rescale preserves intent
		// (30→10, 150→50, 100→33); goals already at/below 10 are untouched.
		3: {"UPDATE students SET daily_xp_goal = MAX(1, CAST(ROUND(daily_xp_goal / 3.0) AS INTEGER)) WHERE daily_xp_goal > 10"},
		// v4 daily goal range: the goal is now bounded 10..100 (xp.MinDailyGoal..MaxDailyGoal)
		// and defaults to 30. v3 left every learner at or below 10, so a learner who had
		// accepted the old default sat at the *new minimum* rather than the default, and goals
		// above 100 are no longer accepted by the setter. Clamp into range and lift the old
		// default to the new one; anyone who chose their own goal inside 10..100 keeps it.
		4: {`UPDATE students SET daily_xp_goal = 30 WHERE daily_xp_goal = 10`,
			`UPDATE students SET daily_xp_goal = 100 WHERE daily_xp_goal > 100`,
			`UPDATE students SET daily_xp_goal = 10 WHERE daily_xp_goal < 10`},
	}
	for v := 1; v <= len(dataMigrations); v++ {
		if err := runOnceSQLite(db, v, dataMigrations[v]); err != nil {
			return err
		}
	}
	return nil
}

// runOnceSQLite applies versioned data backfills exactly once. INSERT OR
// IGNORE keeps concurrent first boots from failing on the version row.
func runOnceSQLite(db *sql.DB, version int, stmts []string) error {
	var one int
	if err := db.QueryRow("SELECT 1 FROM schema_migrations WHERE version = ?", version).Scan(&one); err == nil {
		return nil
	} else if err != sql.ErrNoRows {
		return fmt.Errorf("check schema migration %d: %w", version, err)
	}
	for _, stmt := range stmts {
		if _, err := db.Exec(stmt); err != nil {
			return fmt.Errorf("apply schema migration %d: %w", version, err)
		}
	}
	if _, err := db.Exec("INSERT OR IGNORE INTO schema_migrations (version) VALUES (?)", version); err != nil {
		return fmt.Errorf("record schema migration %d: %w", version, err)
	}
	return nil
}

func isDuplicateColumn(err error) bool {
	return strings.Contains(err.Error(), "duplicate column name")
}

// ensureEmailUniqueIndexSQLite enforces one-email-one-account at the DB layer
// for every writer at once (signup, OAuth create/merge, profile update,
// verification) — the app-level 409 checks are check-then-insert and race.
// Expression index on lower(email) matches FindByEmail's lookup semantics;
// the partial predicate keeps legacy ” rows legal. While real duplicates
// exist the index is skipped with a log (boot proceeds) and retried next
// boot; the audit names no addresses, only a count.
func ensureEmailUniqueIndexSQLite(db *sql.DB) error {
	var name string
	if err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='index' AND name='uidx_students_email'").Scan(&name); err == nil && name != "" {
		return nil
	} else if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("check email unique index: %w", err)
	}
	var dupes int
	if err := db.QueryRow("SELECT COUNT(*) FROM (SELECT lower(email) FROM students WHERE email IS NOT NULL AND email != '' GROUP BY lower(email) HAVING COUNT(*) > 1)").Scan(&dupes); err != nil {
		return fmt.Errorf("audit duplicate emails: %w", err)
	}
	if dupes > 0 {
		log.Printf("storage: %d duplicate email groups; skipping unique email index until resolved", dupes)
		return nil
	}
	if _, err := db.Exec("CREATE UNIQUE INDEX uidx_students_email ON students(lower(email)) WHERE email IS NOT NULL AND email != ''"); err != nil {
		return fmt.Errorf("create email unique index: %w", err)
	}
	return nil
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

func (s *SQLiteStore) ClaimGuestStudent(id, name string) (*Student, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("guest id is required")
	}
	name = strings.TrimSpace(name)
	if name == "" {
		name = "Guest"
	}
	now := time.Now().UTC()
	_, err := s.db.Exec(
		"INSERT OR IGNORE INTO students (id, name, settings, created_at) VALUES (?, ?, '{}', ?)",
		id, name, now.Format(time.RFC3339),
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
	return st, nil
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
	var xpTotal, xpToday, diagCompleted, dailyGoal, leagueMoved sql.NullInt64
	var league, leagueWeek, shareToken, email, googleID, avatarURL sql.NullString
	var emailVerified sql.NullInt64
	var createdAt string
	err := row.Scan(&st.ID, &st.Name, &username, &passwordHash, &courseID, &xpTotal, &xpToday, &xpDate, &diagCompleted, &dailyGoal, &settings, &createdAt, &league, &leagueWeek, &leagueMoved, &shareToken, &email, &googleID, &avatarURL, &emailVerified)
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
	st.League = league.String
	if st.League == "" {
		st.League = "bronze"
	}
	st.LeagueWeek = leagueWeek.String
	if leagueMoved.Valid {
		st.LeagueMoved = int(leagueMoved.Int64)
	}
	st.ShareToken = shareToken.String
	st.Email = email.String
	st.EmailVerified = emailVerified.Valid && emailVerified.Int64 == 1
	st.GoogleID = googleID.String
	st.AvatarURL = avatarURL.String
	if createdAt != "" {
		st.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
		if err != nil {
			st.CreatedAt = time.Now()
		}
	}
	return &st, nil
}

func (s *SQLiteStore) GetStudent(id string) (*Student, error) {
	row := s.db.QueryRow("SELECT id, name, username, password_hash, course_id, xp_total, xp_today, xp_date, diagnostic_completed, daily_xp_goal, settings, created_at, league, league_week, league_moved, share_token, email, google_id, avatar_url, email_verified FROM students WHERE id = ?", id)
	return scanStudent(row)
}

func (s *SQLiteStore) FindByUsername(username string) (*Student, error) {
	row := s.db.QueryRow("SELECT id, name, username, password_hash, course_id, xp_total, xp_today, xp_date, diagnostic_completed, daily_xp_goal, settings, created_at, league, league_week, league_moved, share_token, email, google_id, avatar_url, email_verified FROM students WHERE username = ?", username)
	return scanStudent(row)
}

func (s *SQLiteStore) FindStudentByIdentity(provider, providerID string) (*Student, error) {
	if provider == "" || providerID == "" {
		return nil, nil
	}
	row := s.db.QueryRow(`SELECT id, name, username, password_hash, course_id, xp_total, xp_today, xp_date, diagnostic_completed, daily_xp_goal, settings, created_at, league, league_week, league_moved, share_token, email, google_id, avatar_url, email_verified FROM students WHERE id = (
		SELECT student_id FROM identities WHERE provider = ? AND provider_id = ?
	)`, provider, providerID)
	return scanStudent(row)
}

func (s *SQLiteStore) FindByEmail(email string) (*Student, error) {
	if email == "" {
		return nil, nil
	}
	row := s.db.QueryRow("SELECT id, name, username, password_hash, course_id, xp_total, xp_today, xp_date, diagnostic_completed, daily_xp_goal, settings, created_at, league, league_week, league_moved, share_token, email, google_id, avatar_url, email_verified FROM students WHERE lower(email) = lower(?) AND email != '' LIMIT 1", email)
	return scanStudent(row)
}

func (s *SQLiteStore) CreateGoogleUser(name, email, googleID, avatarURL string) (*Student, error) {
	id := newUUID()
	now := time.Now().UTC()
	_, err := s.db.Exec(
		"INSERT INTO students (id, name, email, email_verified, google_id, avatar_url, settings, created_at) VALUES (?, ?, ?, 1, ?, ?, '{}', ?)",
		id, name, email, googleID, avatarURL, now.Format(time.RFC3339),
	)
	if err != nil {
		return nil, fmt.Errorf("create google user: %w", err)
	}
	if googleID != "" {
		if err := s.CreateIdentity("google", googleID, id, email, true); err != nil {
			return nil, err
		}
	}
	return &Student{ID: id, Name: name, Email: email, EmailVerified: true, GoogleID: googleID, AvatarURL: avatarURL, Settings: "{}", CreatedAt: now}, nil
}

func (s *SQLiteStore) SetEmailVerified(studentID string, verified bool) error {
	v := 0
	if verified {
		v = 1
	}
	_, err := s.db.Exec("UPDATE students SET email_verified = ? WHERE id = ?", v, studentID)
	if err != nil {
		return fmt.Errorf("set email verified: %w", err)
	}
	return nil
}

func (s *SQLiteStore) CreateIdentity(provider, providerID, studentID, email string, emailVerified bool) error {
	v := 0
	if emailVerified {
		v = 1
	}
	_, err := s.db.Exec(`INSERT INTO identities (provider, provider_id, student_id, email, email_verified, linked_at)
		VALUES (?, ?, ?, ?, ?, datetime('now'))`, provider, providerID, studentID, email, v)
	if err != nil {
		return fmt.Errorf("create identity: %w", err)
	}
	return nil
}

func (s *SQLiteStore) ListIdentities(studentID string) ([]Identity, error) {
	rows, err := s.db.Query("SELECT provider, provider_id, student_id, email, email_verified, linked_at FROM identities WHERE student_id = ? ORDER BY provider", studentID)
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
		} else if t, err := time.Parse("2006-01-02 15:04:05", linkedAt); err == nil {
			id.LinkedAt = t
		}
		out = append(out, id)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list identities: %w", err)
	}
	return out, nil
}

func (s *SQLiteStore) DeleteIdentity(provider, studentID string) error {
	_, err := s.db.Exec("DELETE FROM identities WHERE provider = ? AND student_id = ?", provider, studentID)
	if err != nil {
		return fmt.Errorf("delete identity: %w", err)
	}
	return nil
}

func (s *SQLiteStore) CreateEmailVerification(tokenHash string, studentID string, expiresAt time.Time) error {
	_, err := s.db.Exec("INSERT INTO email_verifications (token_hash, student_id, expires_at, used) VALUES (?, ?, ?, 0)",
		tokenHash, studentID, expiresAt.UTC().Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("create email verification: %w", err)
	}
	return nil
}

func (s *SQLiteStore) ConsumeEmailVerification(tokenHash string) (string, bool, error) {
	return consumeSingleUseToken(s.db, "email_verifications", tokenHash)
}

func (s *SQLiteStore) CreateLinkToken(tokenHash string, studentID string, expiresAt time.Time) error {
	_, err := s.db.Exec("INSERT INTO link_tokens (token_hash, student_id, expires_at, used) VALUES (?, ?, ?, 0)",
		tokenHash, studentID, expiresAt.UTC().Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("create link token: %w", err)
	}
	return nil
}

func (s *SQLiteStore) ConsumeLinkToken(tokenHash string) (string, bool, error) {
	return consumeSingleUseToken(s.db, "link_tokens", tokenHash)
}

// consumeSingleUseToken atomically validates (exists, unused, unexpired) and
// burns a token row, returning the student it belongs to.
func consumeSingleUseToken(db *sql.DB, table, tokenHash string) (string, bool, error) {
	tx, err := db.Begin()
	if err != nil {
		return "", false, fmt.Errorf("token tx: %w", err)
	}
	defer tx.Rollback()
	var studentID, expiresAt string
	var used int
	err = tx.QueryRow("SELECT student_id, expires_at, used FROM "+table+" WHERE token_hash = ?", tokenHash).
		Scan(&studentID, &expiresAt, &used)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, fmt.Errorf("lookup token: %w", err)
	}
	exp, err := time.Parse(time.RFC3339, expiresAt)
	if err != nil {
		exp, err = time.Parse("2006-01-02 15:04:05", expiresAt)
	}
	if err != nil || used != 0 || time.Now().UTC().After(exp) {
		return "", false, nil
	}
	if _, err := tx.Exec("UPDATE "+table+" SET used = 1 WHERE token_hash = ?", tokenHash); err != nil {
		return "", false, fmt.Errorf("burn token: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return "", false, fmt.Errorf("commit token: %w", err)
	}
	return studentID, true, nil
}

func (s *SQLiteStore) SetShareToken(studentID, token string) error {
	_, err := s.db.Exec("UPDATE students SET share_token = ? WHERE id = ?", token, studentID)
	if err != nil {
		return fmt.Errorf("set share token: %w", err)
	}
	return nil
}

func (s *SQLiteStore) SetAvatarURL(studentID, avatarURL string) error {
	_, err := s.db.Exec("UPDATE students SET avatar_url = ? WHERE id = ?", avatarURL, studentID)
	if err != nil {
		return fmt.Errorf("set avatar url: %w", err)
	}
	return nil
}

func (s *SQLiteStore) CreateOAuthUser(provider, providerID, name, email string, emailVerified bool, avatarURL string) (*Student, error) {
	id := newUUID()
	now := time.Now().UTC()
	ev := 0
	if emailVerified {
		ev = 1
	}
	if name == "" {
		name = email
	}
	_, err := s.db.Exec(
		"INSERT INTO students (id, name, email, email_verified, avatar_url, settings, created_at) VALUES (?, ?, ?, ?, ?, '{}', ?)",
		id, name, email, ev, avatarURL, now.Format(time.RFC3339),
	)
	if err != nil {
		return nil, fmt.Errorf("create oauth user: %w", err)
	}
	if err := s.CreateIdentity(provider, providerID, id, email, emailVerified); err != nil {
		return nil, err
	}
	return &Student{ID: id, Name: name, Email: email, EmailVerified: emailVerified, AvatarURL: avatarURL, Settings: "{}", CreatedAt: now}, nil
}

func (s *SQLiteStore) GetStudentByShareToken(token string) (*Student, error) {
	if token == "" {
		return nil, nil
	}
	row := s.db.QueryRow("SELECT id, name, username, password_hash, course_id, xp_total, xp_today, xp_date, diagnostic_completed, daily_xp_goal, settings, created_at, league, league_week, league_moved, share_token, email, google_id, avatar_url, email_verified FROM students WHERE share_token = ?", token)
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
		SET xp_total = MAX(0, xp_total + ?),
		    xp_today = CASE WHEN xp_date = ? THEN MAX(0, xp_today + ?) ELSE MAX(0, ?) END,
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

func (s *SQLiteStore) UpdateStudentName(studentID string, name string) error {
	_, err := s.db.Exec("UPDATE students SET name = ? WHERE id = ?", name, studentID)
	if err != nil {
		return fmt.Errorf("update student name: %w", err)
	}
	return nil
}

func (s *SQLiteStore) SetUsername(studentID string, username string) error {
	_, err := s.db.Exec("UPDATE students SET username = ? WHERE id = ?", username, studentID)
	if err != nil {
		return fmt.Errorf("set username: %w", err)
	}
	return nil
}

func (s *SQLiteStore) ListStudentsMissingUsernames() ([]string, error) {
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

func (s *SQLiteStore) SetEmail(studentID string, email string) error {
	_, err := s.db.Exec("UPDATE students SET email = ? WHERE id = ?", email, studentID)
	if err != nil {
		return fmt.Errorf("set email: %w", err)
	}
	return nil
}

func (s *SQLiteStore) SetPasswordHash(studentID string, hash string) error {
	_, err := s.db.Exec("UPDATE students SET password_hash = ? WHERE id = ?", hash, studentID)
	if err != nil {
		return fmt.Errorf("set password hash: %w", err)
	}
	return nil
}

func (s *SQLiteStore) CreatePasswordReset(tokenHash string, studentID string, expiresAt time.Time) error {
	_, err := s.db.Exec("INSERT INTO password_resets (token_hash, student_id, expires_at, used) VALUES (?, ?, ?, 0)",
		tokenHash, studentID, expiresAt.UTC().Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("create password reset: %w", err)
	}
	return nil
}

// ConsumePasswordReset atomically validates (exists, unused, unexpired) and
// burns a reset token, returning the student it belongs to.
func (s *SQLiteStore) ConsumePasswordReset(tokenHash string) (string, bool, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return "", false, fmt.Errorf("reset tx: %w", err)
	}
	defer tx.Rollback()
	var studentID, expiresAt string
	var used int
	err = tx.QueryRow("SELECT student_id, expires_at, used FROM password_resets WHERE token_hash = ?", tokenHash).
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
	if _, err := tx.Exec("UPDATE password_resets SET used = 1 WHERE token_hash = ?", tokenHash); err != nil {
		return "", false, fmt.Errorf("burn reset: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return "", false, fmt.Errorf("commit reset: %w", err)
	}
	return studentID, true, nil
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

func (s *SQLiteStore) SetAvatarImage(studentID string, contentType string, data []byte) error {
	_, err := s.db.Exec(`INSERT INTO avatar_images (student_id, content_type, bytes, updated_at)
		VALUES (?, ?, ?, datetime('now'))
		ON CONFLICT(student_id) DO UPDATE SET content_type = excluded.content_type, bytes = excluded.bytes, updated_at = datetime('now')`,
		studentID, contentType, data)
	if err != nil {
		return fmt.Errorf("set avatar image: %w", err)
	}
	return nil
}

func (s *SQLiteStore) GetAvatarImage(studentID string) (string, []byte, bool, error) {
	var contentType string
	var data []byte
	err := s.db.QueryRow("SELECT content_type, bytes FROM avatar_images WHERE student_id = ?", studentID).Scan(&contentType, &data)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil, false, nil
		}
		return "", nil, false, fmt.Errorf("get avatar image: %w", err)
	}
	return contentType, data, true, nil
}

func (s *SQLiteStore) ClearAvatarImage(studentID string) error {
	_, err := s.db.Exec("DELETE FROM avatar_images WHERE student_id = ?", studentID)
	if err != nil {
		return fmt.Errorf("clear avatar image: %w", err)
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

func (s *SQLiteStore) UpsertProgressBatch(ps []*ConceptProgress) error {
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
			p.SM2EFactor, nullTime(p.MasteredAt), p.WeaknessScore,
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

// EnsureSession inserts a sessions row idempotently for ephemeral
// diagnostic/quiz UUIDs (attempts.session_id references sessions(id)).
func (s *SQLiteStore) EnsureSession(id, studentID string) error {
	_, err := s.db.Exec(
		"INSERT OR IGNORE INTO sessions (id, student_id, started_at) VALUES (?, ?, ?)",
		id, studentID, time.Now().UTC().Format(time.RFC3339),
	)
	if err != nil {
		return fmt.Errorf("ensure session: %w", err)
	}
	return nil
}

// ServerSessions is a durable KV for restart-proof server state.

func (s *SQLiteStore) UpsertServerSession(kind, key, value, expiresAt string) error {
	_, err := s.db.Exec(`
		INSERT INTO server_sessions (kind, key, value, expires_at, updated_at)
		VALUES (?, ?, ?, ?, datetime('now'))
		ON CONFLICT(kind, key) DO UPDATE SET
			value = excluded.value,
			expires_at = excluded.expires_at,
			updated_at = datetime('now')
	`, kind, key, value, expiresAt)
	if err != nil {
		return fmt.Errorf("upsert server session: %w", err)
	}
	return nil
}

func (s *SQLiteStore) GetServerSession(kind, key string) (string, string, bool, error) {
	row := s.db.QueryRow("SELECT value, expires_at FROM server_sessions WHERE kind = ? AND key = ?", kind, key)
	var value, expiresAt string
	if err := row.Scan(&value, &expiresAt); err != nil {
		if err == sql.ErrNoRows {
			return "", "", false, nil
		}
		return "", "", false, fmt.Errorf("get server session: %w", err)
	}
	return value, expiresAt, true, nil
}

func (s *SQLiteStore) DeleteServerSession(kind, key string) error {
	if _, err := s.db.Exec("DELETE FROM server_sessions WHERE kind = ? AND key = ?", kind, key); err != nil {
		return fmt.Errorf("delete server session: %w", err)
	}
	return nil
}

func (s *SQLiteStore) SweepServerSessions() error {
	// expires_at is RFC3339 UTC written by callers — compare in Go, not SQL,
	// so both dialects share one correct clock. Table stays tiny.
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
		if _, err := s.db.Exec("DELETE FROM server_sessions WHERE kind = ? AND key = ?", e.kind, e.key); err != nil {
			return fmt.Errorf("sweep server sessions: %w", err)
		}
	}
	return nil
}

func (s *SQLiteStore) GetActiveSession(sessionID string) (*ActiveSession, error) {
	row := s.db.QueryRow(`
		SELECT session_id, student_id, concept_id, concept_name, expected_answer,
		       attempt_id, question, explanation, diagram, is_review, answered,
		       last_concept_id, session_review, session_new, difficulty, updated_at
		FROM active_sessions WHERE session_id = ?
	`, sessionID)
	var a ActiveSession
	var isReview, answered int
	var updatedAt string
	if err := row.Scan(&a.SessionID, &a.StudentID, &a.ConceptID, &a.ConceptName, &a.ExpectedAnswer,
		&a.AttemptID, &a.Question, &a.Explanation, &a.Diagram, &isReview, &answered,
		&a.LastConceptID, &a.SessionReview, &a.SessionNew, &a.Difficulty, &updatedAt); err != nil {
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

func (s *SQLiteStore) UpsertActiveSession(a *ActiveSession) error {
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := s.db.Exec(`
		INSERT INTO active_sessions
			(session_id, student_id, concept_id, concept_name, expected_answer,
			 attempt_id, question, explanation, diagram, is_review, answered,
			 last_concept_id, session_review, session_new, difficulty, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(session_id) DO UPDATE SET
			student_id      = excluded.student_id,
			concept_id      = excluded.concept_id,
			concept_name    = excluded.concept_name,
			expected_answer = excluded.expected_answer,
			attempt_id      = excluded.attempt_id,
			question        = excluded.question,
			explanation     = excluded.explanation,
			diagram         = excluded.diagram,
			is_review       = excluded.is_review,
			answered        = excluded.answered,
			last_concept_id = excluded.last_concept_id,
			session_review  = excluded.session_review,
			session_new     = excluded.session_new,
			difficulty      = excluded.difficulty,
			updated_at      = excluded.updated_at
	`, a.SessionID, a.StudentID, a.ConceptID, a.ConceptName, a.ExpectedAnswer,
		a.AttemptID, a.Question, a.Explanation, a.Diagram, boolToInt(a.IsReview), boolToInt(a.Answered),
		a.LastConceptID, a.SessionReview, a.SessionNew, a.Difficulty, now)
	if err != nil {
		return fmt.Errorf("upsert active session: %w", err)
	}
	return nil
}

func (s *SQLiteStore) DeleteActiveSession(sessionID string) error {
	if _, err := s.db.Exec("DELETE FROM active_sessions WHERE session_id = ?", sessionID); err != nil {
		return fmt.Errorf("delete active session: %w", err)
	}
	return nil
}

// ResetProgress wipes one student's learning record (see Repository docs
// for the exact scope). One transaction, child tables first.
func (s *SQLiteStore) ResetProgress(studentID string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("begin progress reset: %w", err)
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
		{"DELETE FROM attempts WHERE student_id = ?", []interface{}{studentID}},
		{"DELETE FROM active_sessions WHERE student_id = ?", []interface{}{studentID}},
		{"DELETE FROM sessions WHERE student_id = ?", []interface{}{studentID}},
		{"DELETE FROM concept_progress WHERE student_id = ?", []interface{}{studentID}},
		{"DELETE FROM student_topic_speed WHERE student_id = ?", []interface{}{studentID}},
		{"DELETE FROM quiz_completions WHERE student_id = ?", []interface{}{studentID}},
		{"DELETE FROM server_sessions WHERE kind IN ('study_expected','study_plan') AND (key = ? OR key LIKE ? || '|%')", []interface{}{studentID, studentID}},
		{`UPDATE students SET xp_total = 0, xp_today = 0, diagnostic_completed = 0,
			league = 'bronze', league_week = '', league_moved = 0 WHERE id = ?`, []interface{}{studentID}},
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

// DeleteAccount removes one student and every owned row (see Repository
// docs for the exact scope). One transaction, child tables first, explicit
// deletes rather than cascade reliance.
func (s *SQLiteStore) DeleteAccount(studentID string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("begin account delete: %w", err)
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
		{"DELETE FROM email_verifications WHERE student_id = ?", []interface{}{studentID}},
		{"DELETE FROM password_resets WHERE student_id = ?", []interface{}{studentID}},
		{"DELETE FROM link_tokens WHERE student_id = ?", []interface{}{studentID}},
		{"DELETE FROM identities WHERE student_id = ?", []interface{}{studentID}},
		{"DELETE FROM avatar_images WHERE student_id = ?", []interface{}{studentID}},
		{"DELETE FROM attempts WHERE student_id = ?", []interface{}{studentID}},
		{"DELETE FROM active_sessions WHERE student_id = ?", []interface{}{studentID}},
		{"DELETE FROM sessions WHERE student_id = ?", []interface{}{studentID}},
		{"DELETE FROM concept_progress WHERE student_id = ?", []interface{}{studentID}},
		{"DELETE FROM student_topic_speed WHERE student_id = ?", []interface{}{studentID}},
		{"DELETE FROM quiz_completions WHERE student_id = ?", []interface{}{studentID}},
		{"DELETE FROM server_sessions WHERE kind IN ('study_expected','study_plan','study_plan_prefs') AND (key = ? OR key LIKE ? || '|%')", []interface{}{studentID, studentID}},
		{"DELETE FROM question_reports WHERE reporter_id = ?", []interface{}{studentID}},
		{"DELETE FROM students WHERE id = ?", []interface{}{studentID}},
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

func (s *SQLiteStore) RecordAttempt(entry AttemptEntry) error {
	_, err := s.db.Exec(`
		INSERT INTO attempts
			(session_id, student_id, concept_id, answer, expected,
			 correct, elapsed_seconds, timestamp, question, source, explanation, difficulty)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		entry.SessionID, entry.StudentID, entry.ConceptID,
		entry.Answer, entry.Expected,
		boolToInt(entry.Correct), entry.ElapsedSeconds,
		entry.Timestamp.Format(time.RFC3339),
		entry.Question, entry.Source, entry.Explanation, entry.Difficulty,
	)
	if err != nil {
		return fmt.Errorf("record attempt: %w", err)
	}
	return nil
}

// GetRecentAttemptsForConcept implements the bounded per-concept window the evidence
// model reads on each graded answer. Newest last, matching BuildEvidence's expectation.
//
// Ordered descending so the LIMIT takes the most recent rows, then reversed on the way
// out: the query needs newest-first to bound itself, and the evidence model reads
// oldest-last so its recency weights increase toward the present.
func (s *SQLiteStore) GetRecentAttemptsForConcept(studentID, conceptID string, limit int) ([]AttemptEntry, error) {
	if limit <= 0 {
		limit = 10
	}
	rows, err := s.db.Query(`
		SELECT id, session_id, student_id, concept_id, answer, expected,
		       correct, elapsed_seconds, timestamp, question, source, explanation, difficulty
		FROM attempts
		WHERE student_id = ? AND concept_id = ?
		ORDER BY timestamp DESC, id DESC
		LIMIT ?
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
func (s *SQLiteStore) UpdateAttemptCorrect(id int64, correct bool) error {
	if _, err := s.db.Exec(`UPDATE attempts SET correct = ? WHERE id = ?`, boolToInt(correct), id); err != nil {
		return fmt.Errorf("update attempt correct: %w", err)
	}
	return nil
}

// GetAllAttempts returns every attempt ordered by student, concept, time
// (aggregate efficacy instrumentation).
func (s *SQLiteStore) GetAllAttempts() ([]AttemptEntry, error) {
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

// GetAttemptsForStudent returns all attempts ordered by concept then time
// (efficacy instrumentation: first-pass / second-pass rates).
func (s *SQLiteStore) GetAttemptsForStudent(studentID string) ([]AttemptEntry, error) {
	rows, err := s.db.Query(`
		SELECT id, session_id, student_id, concept_id, answer, expected,
		       correct, elapsed_seconds, timestamp, question, source, explanation, difficulty
		FROM attempts
		WHERE student_id = ?
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

func (s *SQLiteStore) GetSessionAttempts(studentID, sessionID string) ([]AttemptEntry, error) {
	rows, err := s.db.Query(`
		SELECT id, session_id, student_id, concept_id, answer, expected,
		       correct, elapsed_seconds, timestamp, question, source, explanation, difficulty
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
	// Single scan: aggregate per student in one GROUP BY instead of two
	// correlated COUNT subqueries per row.
	rows, err := s.db.Query(`
		SELECT s.id, s.name, s.username, s.avatar_url, s.settings,
			COUNT(CASE WHEN cp.status = 'MASTERED' THEN 1 END),
			COUNT(CASE WHEN cp.status = 'MASTERED' AND cp.mastered_at >= ? THEN 1 END)
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

// GetLeagueStandings returns every student's league tier + weekly mastered
// (partitioned into tiers by the caller).
func (s *SQLiteStore) GetLeagueStandings() ([]LeagueMember, error) {
	monday := weekStart(time.Now().UTC())
	rows, err := s.db.Query(`
		SELECT s.id, s.name, s.username, s.avatar_url, s.settings, COALESCE(s.league, 'bronze'), s.league_moved,
			COUNT(CASE WHEN cp.status = 'MASTERED' THEN 1 END),
			COUNT(CASE WHEN cp.status = 'MASTERED' AND cp.mastered_at >= ? THEN 1 END)
		FROM students s
		LEFT JOIN concept_progress cp ON cp.student_id = s.id
		GROUP BY s.id, s.name, s.username, s.avatar_url, s.settings, s.league, s.league_moved
	`, monday.Format(time.RFC3339))
	if err != nil {
		return nil, fmt.Errorf("get league standings: %w", err)
	}
	defer rows.Close()

	var out []LeagueMember
	for rows.Next() {
		var m LeagueMember
		var moved int
		var username, avatarURL, settings sql.NullString
		if err := rows.Scan(&m.StudentID, &m.Name, &username, &avatarURL, &settings, &m.Tier, &moved, &m.TotalMastered, &m.WeeklyMastered); err != nil {
			return nil, fmt.Errorf("scan league member: %w", err)
		}
		m.Username = username.String
		m.AvatarURL = avatarURL.String
		m.AvatarSettings = settings.String
		m.AvatarCustom, m.AvatarDicebear, m.AvatarVersion = AvatarBits(settings.String)
		m.Moved = moved
		if m.Tier == "" {
			m.Tier = "bronze"
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) SetLeague(studentID, tier, week string, moved int) error {
	_, err := s.db.Exec(
		"UPDATE students SET league = ?, league_week = ?, league_moved = ? WHERE id = ?",
		tier, week, moved, studentID,
	)
	if err != nil {
		return fmt.Errorf("set league: %w", err)
	}
	return nil
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

func (s *SQLiteStore) GetTopicSpeed(studentID, conceptID string) (*TopicSpeed, error) {
	row := s.db.QueryRow(`SELECT student_id, concept_id, efactor, interval, repetitions, learning_speed FROM student_topic_speed WHERE student_id = ? AND concept_id = ?`, studentID, conceptID)
	var ts TopicSpeed
	if err := row.Scan(&ts.StudentID, &ts.ConceptID, &ts.EFactor, &ts.Interval, &ts.Repetitions, &ts.LearningSpeed); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get topic speed: %w", err)
	}
	return &ts, nil
}

func (s *SQLiteStore) GetAllTopicSpeeds(studentID string) (map[string]*TopicSpeed, error) {
	rows, err := s.db.Query(`SELECT student_id, concept_id, efactor, interval, repetitions, learning_speed FROM student_topic_speed WHERE student_id = ?`, studentID)
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

func (s *SQLiteStore) UpsertTopicSpeed(ts *TopicSpeed) error {
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
		VALUES (?, ?, ?, ?, ?, ?, datetime('now'))
		ON CONFLICT(student_id, concept_id) DO UPDATE SET
			efactor = excluded.efactor,
			interval = excluded.interval,
			repetitions = excluded.repetitions,
			learning_speed = excluded.learning_speed,
			updated_at = datetime('now')
	`, ts.StudentID, ts.ConceptID, ts.EFactor, ts.Interval, ts.Repetitions, ts.LearningSpeed)
	if err != nil {
		return fmt.Errorf("upsert topic speed: %w", err)
	}
	return nil
}

// Batch 1: Quiz 150 XP gate completions.

func (s *SQLiteStore) RecordQuizCompletion(studentID string, xpTotal int) error {
	if _, err := s.db.Exec(`INSERT INTO quiz_completions (student_id, completed_at, xp_total) VALUES (?, datetime('now'), ?)`, studentID, xpTotal); err != nil {
		return fmt.Errorf("record quiz completion: %w", err)
	}
	return nil
}

func (s *SQLiteStore) LastQuizCompletion(studentID string) (*QuizCompletion, error) {
	row := s.db.QueryRow(`SELECT student_id, completed_at, xp_total FROM quiz_completions WHERE student_id = ? ORDER BY id DESC LIMIT 1`, studentID)
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

// Question reports (user complaints about questions/explanations/lessons).

func (s *SQLiteStore) CreateReport(r QuestionReport) (int64, error) {
	if r.Status == "" {
		r.Status = "open"
	}
	res, err := s.db.Exec(`
		INSERT INTO question_reports
			(reporter_id, concept_id, kind, question, expected, explanation,
			 lesson_id, source, session_id, attempt_id, reason, detail, status, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		r.ReporterID, r.ConceptID, r.Kind, r.Question, r.Expected, r.Explanation,
		r.LessonID, r.Source, r.SessionID, r.AttemptID, r.Reason, r.Detail, r.Status,
		formatReportTime(r.CreatedAt),
	)
	if err != nil {
		return 0, fmt.Errorf("create report: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("report id: %w", err)
	}
	return id, nil
}

func (s *SQLiteStore) ListReports(status string, limit, offset int) ([]QuestionReport, error) {
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
			SELECT id, reporter_id, concept_id, kind, question, expected, explanation,
			       lesson_id, source, session_id, attempt_id, reason, detail, status, created_at
			FROM question_reports
			ORDER BY created_at DESC
			LIMIT ? OFFSET ?
		`, limit, offset)
	} else {
		rows, err = s.db.Query(`
			SELECT id, reporter_id, concept_id, kind, question, expected, explanation,
			       lesson_id, source, session_id, attempt_id, reason, detail, status, created_at
			FROM question_reports
			WHERE status = ?
			ORDER BY created_at DESC
			LIMIT ? OFFSET ?
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
			&r.ID, &r.ReporterID, &r.ConceptID, &r.Kind, &r.Question, &r.Expected,
			&r.Explanation, &r.LessonID, &r.Source, &r.SessionID, &r.AttemptID,
			&r.Reason, &r.Detail, &r.Status, &created,
		); err != nil {
			return nil, fmt.Errorf("scan report: %w", err)
		}
		r.CreatedAt = parseReportTime(created)
		out = append(out, r)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) UpdateReportStatus(id int64, status string) error {
	res, err := s.db.Exec(`UPDATE question_reports SET status = ? WHERE id = ?`, status, id)
	if err != nil {
		return fmt.Errorf("update report: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("report not found")
	}
	return nil
}

func formatReportTime(t time.Time) string {
	if t.IsZero() {
		return time.Now().UTC().Format(time.RFC3339)
	}
	return t.UTC().Format(time.RFC3339)
}

func parseReportTime(s string) time.Time {
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t
	}
	// SQLite datetime('now') yields "YYYY-MM-DD HH:MM:SS"
	if t, err := time.Parse("2006-01-02 15:04:05", s); err == nil {
		return t.UTC()
	}
	return time.Time{}
}

// ── Administration (SQLite) ────────────────────────────────────────────────

// adminUserColumns is one constant for the projection, so the SELECT that reads a user for
// administration and the SELECT that counts them cannot drift into disagreeing about which
// columns exist. `password_hash IS NOT NULL AND password_hash != ”` becomes the HasPassword
// boolean: a fact about the account, never the hash itself.
const adminUserColumns = `id, name, username, email, role, created_at, (password_hash IS NOT NULL AND password_hash != '')`

func (s *SQLiteStore) GetStudentRole(studentID string) (Role, error) {
	var raw string
	err := s.db.QueryRow("SELECT role FROM students WHERE id = ?", studentID).Scan(&raw)
	if err == sql.ErrNoRows {
		// A missing account is not an error here. A role lookup happens on every authenticated
		// request, and an id that has been deleted in the last second should read as "no
		// privilege" rather than take down the middleware. The caller's own resource lookup is
		// what reports absence.
		return RoleStudent, nil
	}
	if err != nil {
		return RoleStudent, err
	}
	return NormalizeRole(raw), nil
}

func (s *SQLiteStore) SetStudentRole(studentID string, role Role) error {
	res, err := s.db.Exec("UPDATE students SET role = ? WHERE id = ?", string(role), studentID)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// DemoteAdminSafely refuses to remove the last administrator.
//
// The condition lives inside the UPDATE. A read of COUNT(*) followed by an UPDATE is a race:
// two administrators demoting each other at the same moment both observe two, and the system
// reaches zero — the failure this exists to prevent, reached by the very feature meant to
// prevent it. One statement means the second demotion matches no rows.
//
// The same statement covers the no-op cases for free. An admin who is not the last admin, but
// is already a student, matches nothing and is reported as ErrNoRows, so the caller can tell
// "already a student" from "would have removed the last admin".
func (s *SQLiteStore) DemoteAdminSafely(studentID string) error {
	res, err := s.db.Exec(
		`UPDATE students SET role = 'student'
		 WHERE id = ? AND role = 'admin' AND (SELECT COUNT(*) FROM students WHERE role = 'admin') > 1`,
		studentID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n > 0 {
		return nil
	}
	var role Role
	var err2 error
	role, err2 = s.GetStudentRole(studentID)
	if err2 != nil {
		return err2
	}
	if role != RoleAdmin {
		return sql.ErrNoRows
	}
	return ErrLastAdmin
}

func (s *SQLiteStore) CountAdmins() (int, error) {
	var n int
	err := s.db.QueryRow("SELECT COUNT(*) FROM students WHERE role = 'admin'").Scan(&n)
	return n, err
}

func (s *SQLiteStore) CountAdminUsers() (int, error) {
	var n int
	err := s.db.QueryRow("SELECT COUNT(*) FROM students").Scan(&n)
	return n, err
}

func (s *SQLiteStore) CountStoredQuestions() (int, error) {
	var n int
	err := s.db.QueryRow("SELECT COUNT(*) FROM questions").Scan(&n)
	return n, err
}

// adminQueryer is what a projection read needs, and naming it lets the same code serve both the
// pool and an open transaction. SetRoleAudited has to read the before and after snapshots
// inside its transaction, and a read that went to the pool instead would be looking at
// committed state outside the unit of work it is describing.
type adminQueryer interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

// SearchAdminUsers matches on the four identifiers an administrator actually has for someone.
// An empty query lists newest first, so "who signed up most recently" needs no sorting
// parameter.
//
// The LIKE is over an escaped needle rather than a bare `%` + input, so a search for `x%`
// cannot turn into a full scan-and-everything, and a needle containing `_` does not match any
// character. ESCAPE is what makes the escape character itself matchable.
func (s *SQLiteStore) SearchAdminUsers(query string, limit, offset int) ([]AdminUser, error) {
	return searchAdminUsersSQLite(s.db, query, limit, offset)
}

func searchAdminUsersSQLite(q adminQueryer, query string, limit, offset int) ([]AdminUser, error) {
	term := strings.TrimSpace(query)
	if term == "" {
		return queryAdminUsers(q,
			`SELECT `+adminUserColumns+` FROM students ORDER BY created_at DESC, id DESC LIMIT ? OFFSET ?`,
			limit, offset)
	}
	needle := "%" + escapeLike(term) + "%"
	return queryAdminUsers(q,
		`SELECT `+adminUserColumns+` FROM students
		 WHERE id LIKE ? ESCAPE '\' OR email LIKE ? ESCAPE '\' OR username LIKE ? ESCAPE '\' OR name LIKE ? ESCAPE '\'
		 ORDER BY created_at DESC, id DESC LIMIT ? OFFSET ?`,
		needle, needle, needle, needle, limit, offset)
}

// adminUserByIDSQLite is the exact-id read SetRoleAudited uses. Searching by the id and
// filtering keeps it to the one query in the package that can produce an AdminUser, so there
// is still a single place that decides what an administrator may see.
func adminUserByIDSQLite(q adminQueryer, id string) (*AdminUser, error) {
	rows, err := searchAdminUsersSQLite(q, id, 200, 0)
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

func escapeLike(s string) string {
	r := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	return r.Replace(s)
}

func queryAdminUsers(q adminQueryer, query string, args ...interface{}) ([]AdminUser, error) {
	rows, err := q.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []AdminUser{}
	for rows.Next() {
		var u AdminUser
		var name, username, email, role, createdAt sql.NullString
		var hasPw int
		if err := rows.Scan(&u.ID, &name, &username, &email, &role, &createdAt, &hasPw); err != nil {
			return nil, err
		}
		u.Name = name.String
		u.Username = username.String
		u.Email = email.String
		u.Role = NormalizeRole(role.String)
		u.CreatedAt = createdAt.String
		u.HasPassword = hasPw == 1
		out = append(out, u)
	}
	return out, rows.Err()
}

// RecordAdminAudit is an INSERT and nothing else. There is no update path and no delete path
// anywhere in this package, so the log cannot be rewritten through the repository.
func (s *SQLiteStore) RecordAdminAudit(rec AdminAuditRecord) (int64, error) {
	res, err := s.db.Exec(
		`INSERT INTO admin_audit (actor_id, action, entity_type, entity_id, before_json, after_json, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		rec.ActorID, rec.Action, rec.EntityType, rec.EntityID,
		rec.BeforeJSON, rec.AfterJSON, formatReportTime(parseReportTime(rec.CreatedAt)))
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *SQLiteStore) ListAdminAudit(limit, offset int) ([]AdminAuditRecord, error) {
	rows, err := s.db.Query(
		`SELECT id, actor_id, action, entity_type, entity_id, before_json, after_json, created_at
		 FROM admin_audit ORDER BY id DESC LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, err
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

// SetRoleAudited is the network-facing role change: the mutation and its audit event land
// together or not at all.
//
// The transaction spans the three reads and two writes deliberately. The before-snapshot has to
// be read inside it, or it describes committed state rather than the state being changed; the
// after-snapshot likewise; and the audit row can only honestly name what happened if it is
// written while that is still uncommitted. The last-admin guard stays inside the UPDATE rather
// than becoming a prior SELECT for the reason given on DemoteAdminSafely.
func (s *SQLiteStore) SetRoleAudited(actorID, studentID string, role Role) (RoleChange, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return RoleChange{}, fmt.Errorf("begin role change: %w", err)
	}
	// A no-op after a successful commit, and the safety net for every early return below.
	defer func() { _ = tx.Rollback() }()

	before, err := adminUserByIDSQLite(tx, studentID)
	if err != nil {
		return RoleChange{}, fmt.Errorf("read user before: %w", err)
	}
	if before == nil {
		return RoleChange{}, sql.ErrNoRows
	}
	if before.Role == role {
		// Nothing happened, so nothing is recorded. Returning a change object alongside the
		// sentinel is deliberate: the caller reports `changed: false` from it rather than
		// re-reading the row it was just told was already correct.
		return RoleChange{Before: *before, After: *before, Changed: false}, ErrRoleUnchanged
	}

	action := ActionAdminDemoteUser
	if role == RoleAdmin {
		action = ActionAdminPromoteUser
		if _, err := tx.Exec("UPDATE students SET role = ? WHERE id = ?", string(role), studentID); err != nil {
			return RoleChange{}, fmt.Errorf("promote: %w", err)
		}
	} else {
		res, err := tx.Exec(
			`UPDATE students SET role = 'student'
			 WHERE id = ? AND role = 'admin' AND (SELECT COUNT(*) FROM students WHERE role = 'admin') > 1`,
			studentID)
		if err != nil {
			return RoleChange{}, fmt.Errorf("demote: %w", err)
		}
		// `before.Role` was admin inside this transaction, so zero rows affected can only mean
		// the count condition failed: this was the last administrator.
		if n, _ := res.RowsAffected(); n == 0 {
			return RoleChange{}, ErrLastAdmin
		}
	}

	after, err := adminUserByIDSQLite(tx, studentID)
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
	// The only place an administrative event is written. If this insert fails the deferred
	// rollback discards the role change above, so a role can never be changed without a record.
	if _, err := tx.Exec(
		`INSERT INTO admin_audit (actor_id, action, entity_type, entity_id, before_json, after_json, created_at)
		 VALUES (?, ?, 'user', ?, ?, ?, ?)`,
		actorID, action, studentID, string(beforeJSON), string(afterJSON),
		time.Now().UTC().Format(time.RFC3339)); err != nil {
		return RoleChange{}, fmt.Errorf("record audit: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return RoleChange{}, fmt.Errorf("commit role change: %w", err)
	}
	return RoleChange{Before: *before, After: *after, Changed: true, Action: action}, nil
}
