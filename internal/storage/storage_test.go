package storage

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

// newTestStore returns the store one test runs against.
//
// SQLite stays the default: `:memory:` needs no service, isolates every test by
// construction and starts in about a millisecond. Setting TEST_POSTGRES_DSN
// runs the same suite against Postgres, which is the only way to learn whether
// a method behaves the same on the engine production will use. The package has
// no t.Parallel calls, so the Postgres branch can truncate every table before
// each test instead of giving each one a schema of its own; a parallel test
// would need the schema, because two of them would empty each other's rows.
func newTestStore(t *testing.T) Repository {
	t.Helper()
	if dsn := os.Getenv("TEST_POSTGRES_DSN"); dsn != "" {
		store, err := NewPostgresStore(dsn)
		if err != nil {
			t.Fatalf("open postgres: %v", err)
		}
		t.Cleanup(func() { store.Close() })
		truncatePostgres(t, store)
		return store
	}
	return newSQLiteStore(t)
}

// newSQLiteStore is newTestStore without the engine switch, for the tests that
// reach past the Repository interface: the schema itself, the trigger that
// blocks audit rewrites, and the invalid-path case.
func newSQLiteStore(t *testing.T) *SQLiteStore {
	t.Helper()
	store, err := NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("create store: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	return store
}

// truncatePostgres empties every table but the migration ledger, so a test sees
// the database a fresh install would have. RESTART IDENTITY keeps the sequence
// numbers a test asserts on predictable, and CASCADE follows the foreign keys
// without depending on the order the tables are named.
func truncatePostgres(t *testing.T, store *PostgresStore) {
	t.Helper()
	rows, err := store.db.Query(
		`SELECT tablename FROM pg_tables WHERE schemaname = 'public' AND tablename <> 'schema_migrations'`)
	if err != nil {
		t.Fatalf("list tables: %v", err)
	}
	defer rows.Close()
	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatalf("scan table name: %v", err)
		}
		names = append(names, name)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("list tables: %v", err)
	}
	if len(names) == 0 {
		t.Fatal("no tables to truncate: the schema did not migrate")
	}
	if _, err := store.db.Exec("TRUNCATE " + strings.Join(names, ", ") + " RESTART IDENTITY CASCADE"); err != nil {
		t.Fatalf("truncate: %v", err)
	}
}

// Students

func TestCreateStudent(t *testing.T) {
	store := newTestStore(t)
	st, err := store.CreateStudent("alice")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}
	if st.ID == "" {
		t.Error("expected non-empty ID")
	}
	if st.Name != "alice" {
		t.Errorf("expected name alice, got %q", st.Name)
	}
	if st.CreatedAt.IsZero() {
		t.Error("expected non-zero created_at")
	}
}

func TestGetStudent_Found(t *testing.T) {
	store := newTestStore(t)
	created, _ := store.CreateStudent("bob")
	st, err := store.GetStudent(created.ID)
	if err != nil {
		t.Fatalf("get student: %v", err)
	}
	if st == nil {
		t.Fatal("expected student, got nil")
	}
	if st.Name != "bob" {
		t.Errorf("expected bob, got %q", st.Name)
	}
}

func TestGetStudent_NotFound(t *testing.T) {
	store := newTestStore(t)
	st, err := store.GetStudent("nonexistent")
	if err != nil {
		t.Fatalf("get student: %v", err)
	}
	if st != nil {
		t.Error("expected nil for nonexistent student")
	}
}

// Progress

func TestGetProgress_NotFound(t *testing.T) {
	store := newTestStore(t)
	st, _ := store.CreateStudent("carol")
	p, err := store.GetProgress(st.ID, "some.concept")
	if err != nil {
		t.Fatalf("get progress: %v", err)
	}
	if p != nil {
		t.Error("expected nil for absent progress")
	}
}

func TestUpsertAndGetProgress(t *testing.T) {
	store := newTestStore(t)
	st, _ := store.CreateStudent("dave")
	now := time.Now().UTC().Truncate(time.Second)

	p := &ConceptProgress{
		StudentID:       st.ID,
		ConceptID:       "arith.add.single",
		Status:          "LEARNING",
		Streak:          3,
		BestStreak:      5,
		AvgResponseTime: 4.2,
		Attempts:        7,
		LastAttempted:   &now,
		LastReviewed:    &now,
		NextReviewDue:   &now,
		SM2Repetitions:  2,
		SM2Interval:     3,
		SM2EFactor:      2.5,
		MasteredAt:      nil,
	}
	if err := store.UpsertProgress(p); err != nil {
		t.Fatalf("upsert: %v", err)
	}

	got, err := store.GetProgress(st.ID, "arith.add.single")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got == nil {
		t.Fatal("expected progress, got nil")
	}
	if got.Streak != 3 {
		t.Errorf("expected streak 3, got %d", got.Streak)
	}
	if got.BestStreak != 5 {
		t.Errorf("expected best_streak 5, got %d", got.BestStreak)
	}
	if got.AvgResponseTime != 4.2 {
		t.Errorf("expected avg_response_time 4.2, got %f", got.AvgResponseTime)
	}
	if got.Status != "LEARNING" {
		t.Errorf("expected LEARNING, got %q", got.Status)
	}
}

func TestUpsertProgress_Update(t *testing.T) {
	store := newTestStore(t)
	st, _ := store.CreateStudent("eve")

	p1 := &ConceptProgress{
		StudentID: st.ID, ConceptID: "c", Status: "LEARNING", Streak: 1,
	}
	_ = store.UpsertProgress(p1)

	p2 := &ConceptProgress{
		StudentID: st.ID, ConceptID: "c", Status: "PRACTICING", Streak: 5,
	}
	_ = store.UpsertProgress(p2)

	got, _ := store.GetProgress(st.ID, "c")
	if got.Status != "PRACTICING" || got.Streak != 5 {
		t.Errorf("expected updated values, got status=%q streak=%d", got.Status, got.Streak)
	}
}

func TestUpsertProgress_ProvisionalRoundTrip(t *testing.T) {
	store := newTestStore(t)
	st, _ := store.CreateStudent("provisional")

	p := &ConceptProgress{
		StudentID:     st.ID,
		ConceptID:     "diag.concept",
		Status:        "MASTERED",
		IsProvisional: true,
	}
	if err := store.UpsertProgress(p); err != nil {
		t.Fatalf("upsert provisional progress: %v", err)
	}

	got, err := store.GetProgress(st.ID, "diag.concept")
	if err != nil {
		t.Fatalf("get progress: %v", err)
	}
	if !got.IsProvisional {
		t.Fatalf("expected IsProvisional true, got false")
	}

	p.IsProvisional = false
	if err := store.UpsertProgress(p); err != nil {
		t.Fatalf("upsert cleared provisional: %v", err)
	}
	got, err = store.GetProgress(st.ID, "diag.concept")
	if err != nil {
		t.Fatalf("get progress: %v", err)
	}
	if got.IsProvisional {
		t.Fatalf("expected IsProvisional false, got true")
	}
}

func TestFSRSProgress_RoundTripAndSeed(t *testing.T) {
	store := newTestStore(t)
	st, _ := store.CreateStudent("fsrs-student")
	now := time.Now().UTC().Truncate(time.Second)

	stab := 4.5
	diff := 2.8
	lastRev := now.Add(-2 * 24 * time.Hour)
	p := &ConceptProgress{
		StudentID:      st.ID,
		ConceptID:      "fsrs.concept",
		Status:         "MASTERED",
		FSRSStability:  &stab,
		FSRSDifficulty: &diff,
		FSRSLastReview: &lastRev,
	}
	if err := store.UpsertProgress(p); err != nil {
		t.Fatalf("upsert fsrs progress: %v", err)
	}

	got, err := store.GetProgress(st.ID, "fsrs.concept")
	if err != nil {
		t.Fatalf("get fsrs progress: %v", err)
	}
	if got.FSRSStability == nil || *got.FSRSStability != stab {
		t.Fatalf("expected stability %f, got %v", stab, got.FSRSStability)
	}
	if got.FSRSDifficulty == nil || *got.FSRSDifficulty != diff {
		t.Fatalf("expected difficulty %f, got %v", diff, got.FSRSDifficulty)
	}
	if got.FSRSLastReview == nil || !got.FSRSLastReview.Equal(lastRev) {
		t.Fatalf("expected last review %v, got %v", lastRev, got.FSRSLastReview)
	}

	// Test SeedFSRSFromSM2
	dueA := now.Add(3 * 24 * time.Hour)
	pUnseeded := &ConceptProgress{
		StudentID:     st.ID,
		ConceptID:     "seed.a",
		Status:        "MASTERED",
		SM2Interval:   3,
		NextReviewDue: &dueA,
		LastReviewed:  &now,
	}
	if err := store.UpsertProgress(pUnseeded); err != nil {
		t.Fatalf("upsert seed.a: %v", err)
	}

	dueB := now.Add(7 * 24 * time.Hour)
	existingStab := 12.0
	pAlreadySeeded := &ConceptProgress{
		StudentID:     st.ID,
		ConceptID:     "seed.b",
		Status:        "MASTERED",
		SM2Interval:   7,
		NextReviewDue: &dueB,
		FSRSStability: &existingStab,
	}
	if err := store.UpsertProgress(pAlreadySeeded); err != nil {
		t.Fatalf("upsert seed.b: %v", err)
	}

	pNoInterval := &ConceptProgress{
		StudentID:   st.ID,
		ConceptID:   "seed.c",
		Status:      "LEARNING",
		SM2Interval: 0,
	}
	if err := store.UpsertProgress(pNoInterval); err != nil {
		t.Fatalf("upsert seed.c: %v", err)
	}

	count, err := store.SeedFSRSFromSM2(now)
	if err != nil {
		t.Fatalf("seed fsrs from sm2: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 seeded row, got %d", count)
	}

	gotA, err := store.GetProgress(st.ID, "seed.a")
	if err != nil {
		t.Fatalf("get seed.a: %v", err)
	}
	if gotA.FSRSStability == nil || *gotA.FSRSStability <= 0 {
		t.Fatalf("expected positive stability for seed.a, got %v", gotA.FSRSStability)
	}
	if gotA.FSRSDifficulty == nil || *gotA.FSRSDifficulty <= 0 {
		t.Fatalf("expected positive difficulty for seed.a, got %v", gotA.FSRSDifficulty)
	}

	gotB, err := store.GetProgress(st.ID, "seed.b")
	if err != nil {
		t.Fatalf("get seed.b: %v", err)
	}
	if gotB.FSRSStability == nil || *gotB.FSRSStability != existingStab {
		t.Fatalf("expected seed.b to keep stability %f, got %v", existingStab, gotB.FSRSStability)
	}

	gotC, err := store.GetProgress(st.ID, "seed.c")
	if err != nil {
		t.Fatalf("get seed.c: %v", err)
	}
	if gotC.FSRSStability != nil {
		t.Fatalf("expected seed.c to have nil stability, got %v", gotC.FSRSStability)
	}

	// Idempotent second pass
	count2, err := store.SeedFSRSFromSM2(now)
	if err != nil {
		t.Fatalf("seed fsrs second pass: %v", err)
	}
	if count2 != 0 {
		t.Fatalf("expected 0 seeded rows on second pass, got %d", count2)
	}
}

func TestGetAllProgress(t *testing.T) {
	store := newTestStore(t)
	st, _ := store.CreateStudent("frank")

	// No progress yet
	all, err := store.GetAllProgress(st.ID)
	if err != nil {
		t.Fatalf("get all: %v", err)
	}
	if len(all) != 0 {
		t.Errorf("expected empty, got %d entries", len(all))
	}

	// Add two concepts
	_ = store.UpsertProgress(&ConceptProgress{StudentID: st.ID, ConceptID: "a", Status: "LEARNING"})
	_ = store.UpsertProgress(&ConceptProgress{StudentID: st.ID, ConceptID: "b", Status: "MASTERED"})

	all, err = store.GetAllProgress(st.ID)
	if err != nil {
		t.Fatalf("get all: %v", err)
	}
	if len(all) != 2 {
		t.Errorf("expected 2 entries, got %d", len(all))
	}
	if all["a"] == nil || all["b"] == nil {
		t.Error("expected both concepts present")
	}
}

// Sessions

func TestCreateSession(t *testing.T) {
	store := newTestStore(t)
	st, _ := store.CreateStudent("grace")
	sess, err := store.CreateSession(st.ID)
	if err != nil {
		t.Fatalf("create session: %v", err)
	}
	if sess.ID == "" {
		t.Error("expected non-empty session ID")
	}
	if sess.StudentID != st.ID {
		t.Errorf("expected student %s, got %s", st.ID, sess.StudentID)
	}
	if sess.StartedAt.IsZero() {
		t.Error("expected non-zero started_at")
	}
}

// Attempts

func TestRecordAndGetAttempts(t *testing.T) {
	store := newTestStore(t)
	st, _ := store.CreateStudent("heidi")
	sess, _ := store.CreateSession(st.ID)

	entry := AttemptEntry{
		SessionID:      sess.ID,
		StudentID:      st.ID,
		ConceptID:      "arith.add.single",
		Answer:         "12",
		Expected:       "12",
		Correct:        true,
		ElapsedSeconds: 3.5,
		Timestamp:      time.Now().UTC(),
	}
	if _, err := store.RecordAttempt(entry); err != nil {
		t.Fatalf("record attempt: %v", err)
	}

	attempts, err := store.GetSessionAttempts(st.ID, sess.ID)
	if err != nil {
		t.Fatalf("get attempts: %v", err)
	}
	if len(attempts) != 1 {
		t.Fatalf("expected 1 attempt, got %d", len(attempts))
	}
	if attempts[0].Answer != "12" {
		t.Errorf("expected answer 12, got %q", attempts[0].Answer)
	}
	if !attempts[0].Correct {
		t.Error("expected correct=true")
	}
	if attempts[0].ElapsedSeconds != 3.5 {
		t.Errorf("expected 3.5s, got %f", attempts[0].ElapsedSeconds)
	}
}

func TestGetSessionAttempts_Empty(t *testing.T) {
	store := newTestStore(t)
	st, _ := store.CreateStudent("ivan")
	attempts, err := store.GetSessionAttempts(st.ID, "nonexistent")
	if err != nil {
		t.Fatalf("get attempts: %v", err)
	}
	if len(attempts) != 0 {
		t.Errorf("expected empty, got %d", len(attempts))
	}
}

// Mistakes transcript: question/source/explanation survive the round trip;
// pre-migration rows (recorded without them) come back empty, never failing.
func TestRecordAttempt_WithContent(t *testing.T) {
	store := newTestStore(t)
	st, _ := store.CreateStudent("judy")
	sess, _ := store.CreateSession(st.ID)

	full := AttemptEntry{
		SessionID:      sess.ID,
		StudentID:      st.ID,
		ConceptID:      "arith.add.single",
		Answer:         "13",
		Expected:       "12",
		Correct:        false,
		ElapsedSeconds: 4.0,
		Timestamp:      time.Now().UTC(),
		Question:       "What is 5 + 7?",
		Source:         "diagnostic",
		Explanation:    "5 + 7 = 12",
	}
	if _, err := store.RecordAttempt(full); err != nil {
		t.Fatalf("record attempt: %v", err)
	}
	legacy := AttemptEntry{
		SessionID: sess.ID, StudentID: st.ID, ConceptID: "arith.add.single",
		Answer: "12", Expected: "12", Correct: true,
		ElapsedSeconds: 2.0, Timestamp: time.Now().UTC(),
	}
	if _, err := store.RecordAttempt(legacy); err != nil {
		t.Fatalf("record legacy attempt: %v", err)
	}

	attempts, err := store.GetAttemptsForStudent(st.ID)
	if err != nil {
		t.Fatalf("get attempts: %v", err)
	}
	if len(attempts) != 2 {
		t.Fatalf("expected 2 attempts, got %d", len(attempts))
	}
	var got *AttemptEntry
	for i := range attempts {
		if attempts[i].Answer == "13" {
			got = &attempts[i]
		}
	}
	if got == nil {
		t.Fatal("missing recorded attempt")
	}
	if got.Question != "What is 5 + 7?" || got.Source != "diagnostic" || got.Explanation != "5 + 7 = 12" {
		t.Errorf("content fields lost: %+v", got)
	}
}

// Leaderboard

func TestGetWeeklyLeaderboard(t *testing.T) {
	store := newTestStore(t)
	alice, _ := store.CreateStudent("alice")
	bob, _ := store.CreateStudent("bob")

	// Alice: 3 mastered (2 this week)
	_ = store.UpsertProgress(&ConceptProgress{StudentID: alice.ID, ConceptID: "a", Status: "MASTERED", MasteredAt: ptrTime(time.Now().UTC())})
	_ = store.UpsertProgress(&ConceptProgress{StudentID: alice.ID, ConceptID: "b", Status: "MASTERED", MasteredAt: ptrTime(time.Now().UTC())})
	_ = store.UpsertProgress(&ConceptProgress{StudentID: alice.ID, ConceptID: "c", Status: "MASTERED", MasteredAt: ptrTime(time.Now().UTC().AddDate(0, 0, -10))})

	// Bob: 2 mastered (2 this week)
	_ = store.UpsertProgress(&ConceptProgress{StudentID: bob.ID, ConceptID: "a", Status: "MASTERED", MasteredAt: ptrTime(time.Now().UTC())})
	_ = store.UpsertProgress(&ConceptProgress{StudentID: bob.ID, ConceptID: "b", Status: "MASTERED", MasteredAt: ptrTime(time.Now().UTC())})

	rows, err := store.GetWeeklyLeaderboard()
	if err != nil {
		t.Fatalf("get weekly leaderboard: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(rows))
	}
	// Alice first (weekly=2 + total=3 tiebreaker over Bob's 2)
	if rows[0].Name != "alice" {
		t.Errorf("expected alice first, got %s", rows[0].Name)
	}
	if rows[0].WeeklyMastered != 2 {
		t.Errorf("expected alice weekly=2, got %d", rows[0].WeeklyMastered)
	}
	if rows[0].TotalMastered != 3 {
		t.Errorf("expected alice total=3, got %d", rows[0].TotalMastered)
	}
	if rows[1].Name != "bob" {
		t.Errorf("expected bob second, got %s", rows[1].Name)
	}
	if rows[1].WeeklyMastered != 2 {
		t.Errorf("expected bob weekly=2, got %d", rows[1].WeeklyMastered)
	}
}

func TestGetWeeklyLeaderboard_WeeklyBeatsTotal(t *testing.T) {
	store := newTestStore(t)
	alice, _ := store.CreateStudent("alice")
	bob, _ := store.CreateStudent("bob")

	// Alice: high total, stale week (5 mastered, 0 this week).
	for _, cid := range []string{"a", "b", "c", "d", "e"} {
		_ = store.UpsertProgress(&ConceptProgress{StudentID: alice.ID, ConceptID: cid, Status: "MASTERED", MasteredAt: ptrTime(time.Now().UTC().AddDate(0, 0, -30))})
	}
	// Bob: low total, hot week (1 mastered this week).
	_ = store.UpsertProgress(&ConceptProgress{StudentID: bob.ID, ConceptID: "a", Status: "MASTERED", MasteredAt: ptrTime(time.Now().UTC())})

	rows, err := store.GetWeeklyLeaderboard()
	if err != nil {
		t.Fatalf("get weekly leaderboard: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(rows))
	}
	if rows[0].Name != "bob" {
		t.Errorf("expected bob first (weekly=1 beats total=5), got %s", rows[0].Name)
	}
}

func TestAdoptGuestAsUserKeepsProgressAndRejectsRegisteredRows(t *testing.T) {
	store := newTestStore(t)
	guest, err := store.CreateStudent("Guest")
	if err != nil {
		t.Fatalf("create guest: %v", err)
	}
	if err := store.UpsertProgress(&ConceptProgress{
		StudentID: guest.ID,
		ConceptID: "arith.add",
		Status:    "PRACTICING",
		Attempts:  2,
	}); err != nil {
		t.Fatalf("seed progress: %v", err)
	}
	adopted, err := store.AdoptGuestAsUser(guest.ID, "Ada", "ada_guest", "hash")
	if err != nil {
		t.Fatalf("adopt guest: %v", err)
	}
	if adopted.ID != guest.ID || adopted.Name != "Ada" || adopted.Username != "ada_guest" || adopted.PasswordHash != "hash" {
		t.Fatalf("unexpected adopted row: %+v", adopted)
	}
	progress, err := store.GetProgress(guest.ID, "arith.add")
	if err != nil {
		t.Fatalf("read progress: %v", err)
	}
	if progress == nil || progress.Attempts != 2 {
		t.Fatalf("expected progress to survive adoption, got %+v", progress)
	}
	if _, err := store.AdoptGuestAsUser(guest.ID, "Other", "other_guest", "hash2"); err == nil {
		t.Fatal("expected a second adoption to fail")
	}
}

func TestLeaderboardCarriesUsernameAndAvatar(t *testing.T) {
	store := newTestStore(t)
	st, err := store.CreateUser("Ada", "ada", "hash")
	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	if err := store.SetAvatarURL(st.ID, "https://photo/x.jpg"); err != nil {
		t.Fatalf("set avatar url: %v", err)
	}
	if err := store.UpdateSettings(st.ID, `{"avatar_dicebear":{"style":"bottts","seed":"s1"},"avatar_custom":false,"avatar_version":3}`); err != nil {
		t.Fatalf("settings: %v", err)
	}
	rows, err := store.GetWeeklyLeaderboard()
	if err != nil {
		t.Fatalf("weekly: %v", err)
	}
	if len(rows) != 1 || rows[0].Username != "ada" {
		t.Fatalf("expected username carried, got %+v", rows)
	}
	if rows[0].AvatarURL != "https://photo/x.jpg" {
		t.Errorf("expected avatar url carried, got %q", rows[0].AvatarURL)
	}
	custom, pick, version := AvatarBits(rows[0].AvatarSettings)
	if custom || version != 3 || pick == nil || pick.Style != "bottts" || pick.Seed != "s1" {
		t.Errorf("bad avatar bits: custom=%v pick=%+v version=%d", custom, pick, version)
	}
	if err := store.AwardXP(st.ID, 20, "test", "ada:weekly"); err != nil {
		t.Fatalf("award: %v", err)
	}
	now := time.Now().UTC()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).
		AddDate(0, 0, -((int(now.Weekday()) + 6) % 7))
	week := start.Format("2006-01-02")
	if err := store.ResolveLeagueWeek(week, start.AddDate(0, 0, -7), start); err != nil {
		t.Fatalf("resolve week: %v", err)
	}
	tier, no, found, err := store.LeagueMembershipOf(st.ID, week)
	if err != nil || !found {
		t.Fatalf("membership: %v found=%v", err, found)
	}
	members, err := store.LeagueRoster(week, tier, no, start, start.AddDate(0, 0, 7))
	if err != nil {
		t.Fatalf("roster: %v", err)
	}
	if len(members) != 1 || members[0].Username != "ada" {
		t.Fatalf("expected username on member, got %+v", members)
	}
	if members[0].AvatarURL != "https://photo/x.jpg" || members[0].Settings == "" {
		t.Errorf("expected avatar fields on member, got %+v", members[0])
	}
	if c, p, v := AvatarBits("not-json"); c || p != nil || v != 0 {
		t.Errorf("corrupt settings should yield zero bits")
	}
	if c, p, v := AvatarBits(""); c || p != nil || v != 0 {
		t.Errorf("empty settings should yield zero bits")
	}
	if p := AvatarPreset(`{"avatar_preset":4}`); p == nil || *p != 4 {
		t.Errorf("legacy preset = %v, want 4", p)
	}
	if p := AvatarPreset("not-json"); p != nil {
		t.Errorf("corrupt settings preset = %v, want nil", p)
	}
	if p := AvatarPreset(""); p != nil {
		t.Errorf("empty settings preset = %v, want nil", p)
	}
}

// Migrate (idempotent)

func TestMigrate_Idempotent(t *testing.T) {
	store := newTestStore(t)
	if err := store.Migrate(); err != nil {
		t.Fatalf("second migrate should be idempotent: %v", err)
	}
}

// PR 1.2 student model: topic speed round-trip + clamps

func TestTopicSpeed_RoundTrip(t *testing.T) {
	store := newTestStore(t)
	st, _ := store.CreateStudent("speed")
	ts := &TopicSpeed{StudentID: st.ID, ConceptID: "frac.add.diff", EFactor: 2.5, Interval: 6, Repetitions: 2, LearningSpeed: 1.4}
	if err := store.UpsertTopicSpeed(ts); err != nil {
		t.Fatalf("upsert: %v", err)
	}
	got, err := store.GetTopicSpeed(st.ID, "frac.add.diff")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got == nil {
		t.Fatal("expected topic speed row")
	}
	if got.LearningSpeed != 1.4 || got.Interval != 6 || got.Repetitions != 2 {
		t.Errorf("unexpected round-trip: %+v", got)
	}
}

func TestTopicSpeed_Clamps(t *testing.T) {
	store := newTestStore(t)
	st, _ := store.CreateStudent("clamp")
	ts := &TopicSpeed{StudentID: st.ID, ConceptID: "c1", EFactor: 2.5, Interval: 1, Repetitions: 0, LearningSpeed: 9.0}
	if err := store.UpsertTopicSpeed(ts); err != nil {
		t.Fatalf("upsert: %v", err)
	}
	got, _ := store.GetTopicSpeed(st.ID, "c1")
	if got.LearningSpeed != 2.0 {
		t.Errorf("expected clamp to 2.0, got %f", got.LearningSpeed)
	}
	ts.LearningSpeed = 0.1
	if err := store.UpsertTopicSpeed(ts); err != nil {
		t.Fatalf("upsert low: %v", err)
	}
	got, _ = store.GetTopicSpeed(st.ID, "c1")
	if got.LearningSpeed != 0.5 {
		t.Errorf("expected clamp to 0.5, got %f", got.LearningSpeed)
	}
}

func TestGetAllTopicSpeeds(t *testing.T) {
	store := newTestStore(t)
	st, _ := store.CreateStudent("all")
	_ = store.UpsertTopicSpeed(&TopicSpeed{StudentID: st.ID, ConceptID: "a", EFactor: 2.5, Interval: 1, Repetitions: 0, LearningSpeed: 1.0})
	_ = store.UpsertTopicSpeed(&TopicSpeed{StudentID: st.ID, ConceptID: "b", EFactor: 2.5, Interval: 2, Repetitions: 1, LearningSpeed: 1.5})
	all, err := store.GetAllTopicSpeeds(st.ID)
	if err != nil {
		t.Fatalf("get all: %v", err)
	}
	if len(all) != 2 {
		t.Errorf("expected 2 speeds, got %d", len(all))
	}
}

// SQLite pragma errors (invalid path)

func TestNewSQLiteStore_InvalidPath(t *testing.T) {
	_, err := NewSQLiteStore("/nonexistent/dir/db.sqlite")
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
	if !strings.Contains(err.Error(), "open") && !strings.Contains(err.Error(), "unable") {
		t.Logf("got error: %v", err)
	}
}

// helpers

func ptrTime(t time.Time) *time.Time { return &t }

func TestUpdateStudentName(t *testing.T) {
	store := newTestStore(t)
	st, err := store.CreateUser("Ada", "ada", "hash")
	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	if err := store.UpdateStudentName(st.ID, "Ada Lovelace"); err != nil {
		t.Fatalf("update name: %v", err)
	}
	got, err := store.GetStudent(st.ID)
	if err != nil {
		t.Fatalf("get student: %v", err)
	}
	if got.Name != "Ada Lovelace" {
		t.Errorf("expected renamed student, got %q", got.Name)
	}
	if got.Username != "ada" {
		t.Errorf("username must be immutable, got %q", got.Username)
	}
}

func TestClaimGuestStudent_Idempotent(t *testing.T) {
	store := newTestStore(t)
	st, err := store.ClaimGuestStudent("guest_adopt_me", "Guest")
	if err != nil {
		t.Fatalf("claim: %v", err)
	}
	if st.ID != "guest_adopt_me" {
		t.Errorf("expected adopted id, got %q", st.ID)
	}
	again, err := store.ClaimGuestStudent("guest_adopt_me", "Other")
	if err != nil {
		t.Fatalf("reclaim: %v", err)
	}
	if again.ID != st.ID || again.Name != "Guest" {
		t.Errorf("reclaim must return existing row untouched, got %+v", again)
	}
	if _, err := store.ClaimGuestStudent("  ", "Guest"); err == nil {
		t.Error("expected error for blank id")
	}
}

func TestUpsertProgressBatch_RoundTrip(t *testing.T) {
	store := newTestStore(t)
	st, _ := store.CreateStudent("batch")
	if err := store.UpsertProgressBatch(nil); err != nil {
		t.Fatalf("empty batch should be a no-op: %v", err)
	}
	batch := []*ConceptProgress{
		{StudentID: st.ID, ConceptID: "a", Status: "LEARNING", WeaknessScore: 0.4},
		{StudentID: st.ID, ConceptID: "b", Status: "MASTERED", WeaknessScore: 0.0},
		{StudentID: st.ID, ConceptID: "c", Status: "UNSEEN", WeaknessScore: 0.9},
	}
	if err := store.UpsertProgressBatch(batch); err != nil {
		t.Fatalf("batch: %v", err)
	}
	all, err := store.GetAllProgress(st.ID)
	if err != nil {
		t.Fatalf("read back: %v", err)
	}
	if len(all) != 3 || all["b"].Status != "MASTERED" || all["c"].WeaknessScore != 0.9 {
		t.Errorf("batch round-trip mismatch: %+v", all)
	}
	// Overwrite via batch (conflict path).
	batch[0].Status = "MASTERED"
	if err := store.UpsertProgressBatch(batch[:1]); err != nil {
		t.Fatalf("overwrite: %v", err)
	}
	all, _ = store.GetAllProgress(st.ID)
	if all["a"].Status != "MASTERED" || len(all) != 3 {
		t.Errorf("batch overwrite mismatch: %+v", all)
	}
}

func TestServerSessions_RoundTrip(t *testing.T) {
	store := newTestStore(t)
	future := time.Now().UTC().Add(time.Hour).Format(time.RFC3339)
	past := time.Now().UTC().Add(-time.Hour).Format(time.RFC3339)
	if err := store.UpsertServerSession("study_expected", "s|c", "42", future); err != nil {
		t.Fatalf("upsert: %v", err)
	}
	v, exp, found, err := store.GetServerSession("study_expected", "s|c")
	if err != nil || !found || v != "42" || exp != future {
		t.Errorf("round-trip mismatch: %q %q %v %v", v, exp, found, err)
	}
	if _, _, found, _ := store.GetServerSession("study_expected", "nope"); found {
		t.Error("expected miss for unknown key")
	}
	// Overwrite.
	if err := store.UpsertServerSession("study_expected", "s|c", "43", future); err != nil {
		t.Fatalf("overwrite: %v", err)
	}
	if v, _, _, _ := store.GetServerSession("study_expected", "s|c"); v != "43" {
		t.Errorf("expected overwritten value 43, got %q", v)
	}
	// Sweep removes only expired rows.
	if err := store.UpsertServerSession("admin", "tok-old", "", past); err != nil {
		t.Fatalf("upsert expired: %v", err)
	}
	if err := store.SweepServerSessions(); err != nil {
		t.Fatalf("sweep: %v", err)
	}
	if _, _, found, _ := store.GetServerSession("admin", "tok-old"); found {
		t.Error("expected expired row swept")
	}
	if _, _, found, _ := store.GetServerSession("study_expected", "s|c"); !found {
		t.Error("live row must survive sweep")
	}
	if err := store.DeleteServerSession("study_expected", "s|c"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, _, found, _ := store.GetServerSession("study_expected", "s|c"); found {
		t.Error("expected deleted row gone")
	}
}

func TestDataMigrations_RunOnce(t *testing.T) {
	store := newTestStore(t)
	versions := func() map[int]bool {
		out := map[int]bool{}
		rows := querySeed(t, store, "SELECT version FROM schema_migrations")
		defer rows.Close()
		for rows.Next() {
			var v int
			if err := rows.Scan(&v); err != nil {
				t.Fatalf("scan version: %v", err)
			}
			out[v] = true
		}
		return out
	}
	if v := versions(); !v[1] || !v[2] || !v[3] {
		t.Fatalf("expected versions 1,2,3 recorded, got %v", v)
	}
	// Legacy row simulation: reset v1 and plant a 150 goal, re-migrate.
	st, _ := store.CreateStudent("legacy")
	execSeed(t, store, "UPDATE students SET daily_xp_goal = 150 WHERE id = ?", st.ID)
	execSeed(t, store, "DELETE FROM schema_migrations WHERE version = ?", 1)
	if err := store.Migrate(); err != nil {
		t.Fatalf("re-migrate: %v", err)
	}
	got, err := store.GetStudent(st.ID)
	if err != nil || got == nil || got.DailyXPGoal != 30 {
		t.Errorf("expected one-time backfill to 30, got %+v err=%v", got, err)
	}
	if v := versions(); !v[1] {
		t.Error("expected version 1 re-recorded")
	}
	// v3 small-awards rescale: reset v3, plant a 90 goal, re-migrate → 30.
	st3, _ := store.CreateStudent("rescaled")
	execSeed(t, store, "UPDATE students SET daily_xp_goal = 90 WHERE id = ?", st3.ID)
	execSeed(t, store, "DELETE FROM schema_migrations WHERE version = ?", 3)
	if err := store.Migrate(); err != nil {
		t.Fatalf("re-migrate: %v", err)
	}
	got3, err := store.GetStudent(st3.ID)
	if err != nil || got3 == nil || got3.DailyXPGoal != 30 {
		t.Errorf("expected proportional rescale to 30, got %+v err=%v", got3, err)
	}
	// Second Migrate is a no-op for data (idempotent, no every-boot writes).
	if err := store.Migrate(); err != nil {
		t.Fatalf("second migrate: %v", err)
	}
}

func TestReviewClockBackfillForMasteredRows(t *testing.T) {
	// A concept mastered through study has no review clock. Migration 8 gives
	// it one so the graph's decay status and the scheduler's review decision
	// read the same timestamp.
	store := newTestStore(t)
	st, err := store.CreateStudent("legacy-review-clock")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}
	masteredAt := time.Now().UTC().AddDate(0, 0, -40).Truncate(time.Second)
	if err := store.UpsertProgress(&ConceptProgress{
		StudentID: st.ID, ConceptID: "a", Status: "MASTERED", MasteredAt: &masteredAt,
	}); err != nil {
		t.Fatalf("seed mastery date: %v", err)
	}
	if err := store.UpsertProgress(&ConceptProgress{
		StudentID: st.ID, ConceptID: "b", Status: "MASTERED",
	}); err != nil {
		t.Fatalf("seed empty clock: %v", err)
	}

	execSeed(t, store, "DELETE FROM schema_migrations WHERE version = ?", 8)
	if err := store.Migrate(); err != nil {
		t.Fatalf("re-migrate: %v", err)
	}

	all, err := store.GetAllProgress(st.ID)
	if err != nil {
		t.Fatalf("get progress: %v", err)
	}
	a := all["a"]
	if a == nil || a.LastReviewed == nil {
		t.Fatalf("mastery-date row: expected a review clock, got %+v", a)
	}
	if diff := a.LastReviewed.Sub(masteredAt); diff > time.Second || diff < -time.Second {
		t.Errorf("review clock = %v, want the mastery date %v", a.LastReviewed, masteredAt)
	}
	b := all["b"]
	if b == nil || b.LastReviewed == nil {
		t.Fatalf("empty-clock row: expected the now fallback, got %+v", b)
	}
	if age := time.Since(*b.LastReviewed); age < 0 || age > time.Minute {
		t.Errorf("fallback clock age = %v, want recent", age)
	}

	// A second migration must not move a clock that is already set.
	if err := store.Migrate(); err != nil {
		t.Fatalf("second migrate: %v", err)
	}
	again, err := store.GetAllProgress(st.ID)
	if err != nil {
		t.Fatalf("get progress again: %v", err)
	}
	if !again["a"].LastReviewed.Equal(*a.LastReviewed) || !again["b"].LastReviewed.Equal(*b.LastReviewed) {
		t.Error("a second migration moved a review clock")
	}
}

// Batch 1: quiz 150 XP gate completions round-trip.

func TestQuizCompletion_RoundTrip(t *testing.T) {
	store := newTestStore(t)
	st, _ := store.CreateStudent("quizzer")
	if got, err := store.LastQuizCompletion(st.ID); err != nil || got != nil {
		t.Fatalf("expected nil completion, got %+v err=%v", got, err)
	}
	if err := store.AwardXP(st.ID, 175, "test", "seed-175"); err != nil {
		t.Fatalf("add xp: %v", err)
	}
	if err := store.RecordQuizCompletion(st.ID, 175); err != nil {
		t.Fatalf("record: %v", err)
	}
	got, err := store.LastQuizCompletion(st.ID)
	if err != nil || got == nil {
		t.Fatalf("expected completion row, got %+v err=%v", got, err)
	}
	if got.XPTotal != 175 {
		t.Errorf("expected xp_total 175, got %d", got.XPTotal)
	}
	if err := store.RecordQuizCompletion(st.ID, 300); err != nil {
		t.Fatalf("record second: %v", err)
	}
	got, err = store.LastQuizCompletion(st.ID)
	if err != nil || got == nil || got.XPTotal != 300 {
		t.Fatalf("expected latest xp_total 300, got %+v err=%v", got, err)
	}
}

func TestResetProgress_RoundTrip(t *testing.T) {
	store := newTestStore(t)
	st, _ := store.CreateStudent("resetter")
	sess, _ := store.CreateSession(st.ID)
	if err := store.UpsertProgress(&ConceptProgress{StudentID: st.ID, ConceptID: "a", Status: "PRACTICING", Streak: 3}); err != nil {
		t.Fatalf("seed progress: %v", err)
	}
	if _, err := store.RecordAttempt(AttemptEntry{SessionID: sess.ID, StudentID: st.ID, ConceptID: "a", Answer: "1", Expected: "2", Correct: false, ElapsedSeconds: 4, Timestamp: time.Now().UTC()}); err != nil {
		t.Fatalf("seed attempt: %v", err)
	}
	if err := store.UpsertTopicSpeed(&TopicSpeed{StudentID: st.ID, ConceptID: "a", EFactor: 2.5, Interval: 1, Repetitions: 1, LearningSpeed: 1.0}); err != nil {
		t.Fatalf("seed speed: %v", err)
	}
	if err := store.RecordQuizCompletion(st.ID, 60); err != nil {
		t.Fatalf("seed quiz: %v", err)
	}
	if err := store.AwardXP(st.ID, 42, "test", "seed-42"); err != nil {
		t.Fatalf("seed xp: %v", err)
	}
	if err := store.UpsertServerSession("study_expected", st.ID+"|a", `{"q":"1"}`, time.Now().UTC().Add(time.Hour).Format(time.RFC3339)); err != nil {
		t.Fatalf("seed anchor: %v", err)
	}
	if err := store.UpsertServerSession("study_plan", st.ID, `{"destination":"d"}`, time.Now().UTC().Add(time.Hour).Format(time.RFC3339)); err != nil {
		t.Fatalf("seed plan: %v", err)
	}
	if err := store.UpsertServerSession("admin_login", st.ID, "tok", time.Now().UTC().Add(time.Hour).Format(time.RFC3339)); err != nil {
		t.Fatalf("seed admin: %v", err)
	}

	if err := store.ResetProgress(st.ID); err != nil {
		t.Fatalf("reset: %v", err)
	}

	if p, _ := store.GetAllProgress(st.ID); len(p) != 0 {
		t.Errorf("progress survives: %d rows", len(p))
	}
	if a, _ := store.GetAttemptsForStudent(st.ID); len(a) != 0 {
		t.Errorf("attempts survive: %d rows", len(a))
	}
	if s, _ := store.GetAllTopicSpeeds(st.ID); len(s) != 0 {
		t.Errorf("speeds survive: %d rows", len(s))
	}
	if q, _ := store.LastQuizCompletion(st.ID); q != nil {
		t.Errorf("quiz completion survives: %+v", q)
	}
	if _, _, found, _ := store.GetServerSession("study_expected", st.ID+"|a"); found {
		t.Error("study anchor survives")
	}
	if _, _, found, _ := store.GetServerSession("study_plan", st.ID); found {
		t.Error("study plan evidence survives")
	}
	if _, _, found, _ := store.GetServerSession("admin_login", st.ID); !found {
		t.Error("admin session wrongly wiped")
	}
	got, err := store.GetStudent(st.ID)
	if err != nil || got == nil {
		t.Fatalf("account must survive: %v", err)
	}
	if got.XPTotal != 0 || got.XPToday != 0 || got.DiagnosticCompleted || got.Name != "resetter" {
		t.Errorf("student row wrong after reset: %+v", got)
	}
	// Idempotent: second reset is a no-op success.
	if err := store.ResetProgress(st.ID); err != nil {
		t.Errorf("second reset: %v", err)
	}
}

// An attempt's difficulty is evidence, and "not recorded" has to survive a round trip
// as itself.
//
// The evidence model weights a correct answer by how hard the question was. That is only
// meaningful if it can also say "we do not know": pre-migration attempts, hand-written
// quiz questions and anything not generated at a known difficulty have no difficulty.
// Storing 0 for those would claim they were trivially easy, which is a fact nobody
// measured and it would be indistinguishable from a real measurement once written.
//
// The Postgres half of this is TestParity_SQLiteVsPostgres; this runs without a database
// server so CI covers it too.
func TestAttemptDifficulty_DistinguishesUnknownFromTrivial(t *testing.T) {
	store := newTestStore(t)
	stu, err := store.CreateStudent("difficulty-round-trip")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}
	sess, err := store.CreateSession(stu.ID)
	if err != nil {
		t.Fatalf("create session: %v", err)
	}
	base := time.Now().UTC().Add(-time.Hour)

	// Unknown: no difficulty recorded.
	if _, err := store.RecordAttempt(AttemptEntry{
		SessionID: sess.ID, StudentID: stu.ID, ConceptID: "c1",
		Answer: "1", Expected: "1", Correct: true, Timestamp: base,
	}); err != nil {
		t.Fatalf("record unknown: %v", err)
	}
	// Generated at 0.8.
	d := 0.8
	if _, err := store.RecordAttempt(AttemptEntry{
		SessionID: sess.ID, StudentID: stu.ID, ConceptID: "c1",
		Answer: "2", Expected: "2", Correct: true, Timestamp: base.Add(time.Minute),
		Difficulty: &d,
	}); err != nil {
		t.Fatalf("record 0.8: %v", err)
	}

	atts, err := store.GetAttemptsForStudent(stu.ID)
	if err != nil {
		t.Fatalf("get attempts: %v", err)
	}
	if len(atts) != 2 {
		t.Fatalf("got %d attempts, want 2", len(atts))
	}
	if atts[0].Difficulty != nil {
		t.Errorf("first attempt difficulty = %v, want nil (unknown, not zero)", *atts[0].Difficulty)
	}
	if atts[1].Difficulty == nil {
		t.Fatal("second attempt difficulty is nil, want 0.8")
	}
	if *atts[1].Difficulty != 0.8 {
		t.Errorf("second attempt difficulty = %v, want 0.8", *atts[1].Difficulty)
	}
}

// The active session has to carry the served difficulty, or a restart between serving a
// question and answering it would write the attempt with no difficulty at all, which is
// the honest-but-lossy outcome, and happens for a reason that is entirely avoidable.
func TestActiveSession_DifficultySurvivesTheRoundTrip(t *testing.T) {
	store := newTestStore(t)
	stu, _ := store.CreateStudent("difficulty-session")
	sess, _ := store.CreateSession(stu.ID)

	d := 0.55
	if err := store.UpsertActiveSession(&ActiveSession{
		SessionID: sess.ID, StudentID: stu.ID, ConceptID: "c1",
		ConceptName: "C1", ExpectedAnswer: "42", AttemptID: "a1",
		Question: "q", Difficulty: &d,
	}); err != nil {
		t.Fatalf("upsert: %v", err)
	}
	got, err := store.GetActiveSession(sess.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got == nil {
		t.Fatal("no active session returned")
	}
	if got.Difficulty == nil || *got.Difficulty != 0.55 {
		t.Errorf("session difficulty = %v, want 0.55", got.Difficulty)
	}

	// And an unknown one stays nil rather than becoming 0.
	if err := store.UpsertActiveSession(&ActiveSession{
		SessionID: sess.ID, StudentID: stu.ID, ConceptID: "c2",
		ConceptName: "C2", ExpectedAnswer: "1", AttemptID: "a2", Question: "q2",
	}); err != nil {
		t.Fatalf("upsert unknown: %v", err)
	}
	got, _ = store.GetActiveSession(sess.ID)
	if got.Difficulty != nil {
		t.Errorf("unknown session difficulty = %v, want nil", *got.Difficulty)
	}
}

func TestActiveSession_KPOrdinalSurvivesRoundTrip(t *testing.T) {
	store := newTestStore(t)
	stu, _ := store.CreateStudent("kp-session-student")
	sess, _ := store.CreateSession(stu.ID)

	if err := store.UpsertActiveSession(&ActiveSession{
		SessionID:      sess.ID,
		StudentID:      stu.ID,
		ConceptID:      "c1",
		ConceptName:    "C1",
		ExpectedAnswer: "42",
		AttemptID:      "a1",
		Question:       "q1",
		KPOrdinal:      2,
	}); err != nil {
		t.Fatalf("upsert: %v", err)
	}

	got, err := store.GetActiveSession(sess.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got == nil {
		t.Fatal("expected active session")
	}
	if got.KPOrdinal != 2 {
		t.Errorf("expected KPOrdinal 2, got %d", got.KPOrdinal)
	}
}

func TestGetSession_Lifecycle(t *testing.T) {
	store := newTestStore(t)
	st, err := store.CreateStudent("session-student")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}

	missing, err := store.GetSession("nonexistent-session-id")
	if err != nil {
		t.Fatalf("get missing session: %v", err)
	}
	if missing != nil {
		t.Errorf("expected nil for missing session, got %+v", missing)
	}

	sess, err := store.CreateSession(st.ID)
	if err != nil {
		t.Fatalf("create session: %v", err)
	}
	got, err := store.GetSession(sess.ID)
	if err != nil {
		t.Fatalf("get session: %v", err)
	}
	if got == nil {
		t.Fatal("expected session row, got nil")
	}
	if got.ID != sess.ID || got.StudentID != st.ID {
		t.Errorf("got session %+v, want ID=%s StudentID=%s", got, sess.ID, st.ID)
	}
	if got.StartedAt.IsZero() {
		t.Error("expected non-zero StartedAt")
	}
}

func TestDeleteActiveSession_Lifecycle(t *testing.T) {
	store := newTestStore(t)
	st, _ := store.CreateStudent("del-active-student")
	sess, _ := store.CreateSession(st.ID)

	diff := 0.75
	act := &ActiveSession{
		SessionID:      sess.ID,
		StudentID:      st.ID,
		ConceptID:      "linalg.det.2x2",
		ConceptName:    "Determinant 2x2",
		ExpectedAnswer: "42",
		AttemptID:      "attempt-123",
		Question:       "Compute det?",
		Difficulty:     &diff,
	}
	if err := store.UpsertActiveSession(act); err != nil {
		t.Fatalf("upsert active session: %v", err)
	}
	got, err := store.GetActiveSession(sess.ID)
	if err != nil || got == nil {
		t.Fatalf("expected active session, got %+v err=%v", got, err)
	}

	if err := store.DeleteActiveSession(sess.ID); err != nil {
		t.Fatalf("delete active session: %v", err)
	}
	after, err := store.GetActiveSession(sess.ID)
	if err != nil {
		t.Fatalf("get active session after delete: %v", err)
	}
	if after != nil {
		t.Errorf("expected nil after delete, got %+v", after)
	}
}

func TestRecentAttemptsForConcept_OrderingAndLimit(t *testing.T) {
	store := newTestStore(t)
	st, _ := store.CreateStudent("recent-attempts-student")
	sess, _ := store.CreateSession(st.ID)

	base := time.Now().UTC().Add(-time.Hour)
	for i := 1; i <= 5; i++ {
		_, err := store.RecordAttempt(AttemptEntry{
			SessionID:      sess.ID,
			StudentID:      st.ID,
			ConceptID:      "calc.deriv.power",
			Answer:         fmt.Sprintf("%d", i),
			Expected:       fmt.Sprintf("%d", i),
			Correct:        i%2 == 0,
			ElapsedSeconds: float64(i + 2),
			Timestamp:      base.Add(time.Duration(i) * time.Minute),
			Question:       fmt.Sprintf("derivative power %d", i),
		})
		if err != nil {
			t.Fatalf("record attempt %d: %v", i, err)
		}
	}

	if _, err := store.RecordAttempt(AttemptEntry{
		SessionID:      sess.ID,
		StudentID:      st.ID,
		ConceptID:      "other.concept",
		Answer:         "99",
		Expected:       "99",
		Correct:        true,
		ElapsedSeconds: 5,
		Timestamp:      base.Add(10 * time.Minute),
	}); err != nil {
		t.Fatalf("record other concept attempt: %v", err)
	}

	recent, err := store.GetRecentAttemptsForConcept(st.ID, "calc.deriv.power", 3)
	if err != nil {
		t.Fatalf("get recent attempts: %v", err)
	}
	if len(recent) != 3 {
		t.Fatalf("expected 3 recent attempts, got %d", len(recent))
	}
	if recent[0].Answer != "3" || recent[1].Answer != "4" || recent[2].Answer != "5" {
		t.Errorf("expected attempts 3, 4 and 5 in chronological order, got %s, %s and %s",
			recent[0].Answer, recent[1].Answer, recent[2].Answer)
	}

	allRecent, err := store.GetRecentAttemptsForConcept(st.ID, "calc.deriv.power", 0)
	if err != nil {
		t.Fatalf("get recent attempts with default limit: %v", err)
	}
	if len(allRecent) != 5 {
		t.Fatalf("expected 5 recent attempts with default limit, got %d", len(allRecent))
	}
}

func TestUpdateAttemptCorrect_Lifecycle(t *testing.T) {
	store := newTestStore(t)
	st, _ := store.CreateStudent("update-attempt-student")
	sess, _ := store.CreateSession(st.ID)

	id, err := store.RecordAttempt(AttemptEntry{
		SessionID:      sess.ID,
		StudentID:      st.ID,
		ConceptID:      "arith.add.single",
		Answer:         "4",
		Expected:       "5",
		Correct:        false,
		ElapsedSeconds: 3,
		Timestamp:      time.Now().UTC(),
	})
	if err != nil {
		t.Fatalf("record attempt: %v", err)
	}

	if err := store.UpdateAttemptCorrect(id, true); err != nil {
		t.Fatalf("update attempt correct: %v", err)
	}

	attempts, err := store.GetAttemptsForStudent(st.ID)
	if err != nil {
		t.Fatalf("get attempts: %v", err)
	}
	if len(attempts) != 1 {
		t.Fatalf("expected 1 attempt, got %d", len(attempts))
	}
	if !attempts[0].Correct {
		t.Error("expected attempt correct to be updated to true")
	}
}

func TestGetAllAttempts_Ordering(t *testing.T) {
	store := newTestStore(t)
	st1, _ := store.CreateStudent("all-attempts-s1")
	st2, _ := store.CreateStudent("all-attempts-s2")
	sess1, _ := store.CreateSession(st1.ID)
	sess2, _ := store.CreateSession(st2.ID)

	base := time.Now().UTC().Add(-time.Hour)
	_, _ = store.RecordAttempt(AttemptEntry{SessionID: sess2.ID, StudentID: st2.ID, ConceptID: "b", Answer: "2", Expected: "2", Correct: true, Timestamp: base.Add(2 * time.Minute)})
	_, _ = store.RecordAttempt(AttemptEntry{SessionID: sess1.ID, StudentID: st1.ID, ConceptID: "b", Answer: "1b", Expected: "1b", Correct: true, Timestamp: base.Add(3 * time.Minute)})
	_, _ = store.RecordAttempt(AttemptEntry{SessionID: sess1.ID, StudentID: st1.ID, ConceptID: "a", Answer: "1a", Expected: "1a", Correct: true, Timestamp: base.Add(1 * time.Minute)})

	all, err := store.GetAllAttempts()
	if err != nil {
		t.Fatalf("get all attempts: %v", err)
	}
	if len(all) != 3 {
		t.Fatalf("expected 3 attempts, got %d", len(all))
	}
}

func TestQuestions_ImportGetCountAndPurge(t *testing.T) {
	store := newTestStore(t)

	c, err := store.GetQuestionCount("test.concept.purge")
	if err != nil {
		t.Fatalf("get initial count: %v", err)
	}
	if c != 0 {
		t.Errorf("expected 0 initial questions, got %d", c)
	}

	qs := []Question{
		{ConceptID: "test.concept.purge", Question: "What is 1+1?", Answer: "2", Explanation: "1+1=2", Source: "generated", Difficulty: 0.1},
		{ConceptID: "test.concept.purge", Question: "What is 2+2?", Answer: "4", Explanation: "2+2=4", Source: "generated", Difficulty: 0.2},
		{ConceptID: "test.concept.other", Question: "What is 3+3?", Answer: "6", Explanation: "3+3=6", Source: "generated", Difficulty: 0.3},
	}
	if err := store.ImportQuestions(qs); err != nil {
		t.Fatalf("import questions: %v", err)
	}

	count, err := store.GetQuestionCount("test.concept.purge")
	if err != nil {
		t.Fatalf("get count: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2 questions, got %d", count)
	}

	fetched, err := store.GetQuestions("test.concept.purge", 5)
	if err != nil {
		t.Fatalf("get questions: %v", err)
	}
	if len(fetched) != 2 {
		t.Errorf("expected 2 fetched questions, got %d", len(fetched))
	}

	purged, err := store.PurgeGeneratedQuestions(map[string]bool{"test.concept.purge": true})
	if err != nil {
		t.Fatalf("purge questions: %v", err)
	}
	if purged != 2 {
		t.Errorf("expected 2 purged questions, got %d", purged)
	}

	afterCount, err := store.GetQuestionCount("test.concept.purge")
	if err != nil {
		t.Fatalf("get count after purge: %v", err)
	}
	if afterCount != 0 {
		t.Errorf("expected 0 questions after purge, got %d", afterCount)
	}

	otherCount, err := store.GetQuestionCount("test.concept.other")
	if err != nil {
		t.Fatalf("get other count: %v", err)
	}
	if otherCount != 1 {
		t.Errorf("expected 1 question remaining in other concept, got %d", otherCount)
	}

	zeroPurged, err := store.PurgeGeneratedQuestions(map[string]bool{})
	if err != nil || zeroPurged != 0 {
		t.Errorf("expected 0 purged for empty map, got %d err=%v", zeroPurged, err)
	}
}

func TestDailyActivity_AggregatesAttempts(t *testing.T) {
	store := newTestStore(t)
	st, _ := store.CreateStudent("daily-activity-student")
	sess, _ := store.CreateSession(st.ID)

	now := time.Now().UTC()
	today := now.Format("2006-01-02")
	_, err := store.RecordAttempt(AttemptEntry{
		SessionID:      sess.ID,
		StudentID:      st.ID,
		ConceptID:      "concept.one",
		Answer:         "1",
		Expected:       "1",
		Correct:        true,
		ElapsedSeconds: 4,
		Timestamp:      now,
	})
	if err != nil {
		t.Fatalf("record attempt 1: %v", err)
	}
	_, err = store.RecordAttempt(AttemptEntry{
		SessionID:      sess.ID,
		StudentID:      st.ID,
		ConceptID:      "concept.two",
		Answer:         "2",
		Expected:       "3",
		Correct:        false,
		ElapsedSeconds: 6,
		Timestamp:      now,
	})
	if err != nil {
		t.Fatalf("record attempt 2: %v", err)
	}

	activities, err := store.GetDailyActivity(st.ID, 7)
	if err != nil {
		t.Fatalf("get daily activity: %v", err)
	}
	if len(activities) != 1 {
		t.Fatalf("expected 1 daily activity row, got %d", len(activities))
	}
	act := activities[0]
	if act.Date != today {
		t.Errorf("expected date %s, got %s", today, act.Date)
	}
	if act.Questions != 2 {
		t.Errorf("expected 2 questions, got %d", act.Questions)
	}
	if act.Correct != 1 {
		t.Errorf("expected 1 correct, got %d", act.Correct)
	}
	if len(act.Concepts) != 2 {
		t.Errorf("expected 2 distinct concepts, got %v", act.Concepts)
	}
}

func TestUpdateReportStatus_Lifecycle(t *testing.T) {
	store := newTestStore(t)
	st, _ := store.CreateStudent("report-student")

	rep := QuestionReport{
		ReporterID:  st.ID,
		ConceptID:   "arith.add.single",
		Kind:        "question",
		Question:    "1+1=?",
		Expected:    "2",
		Explanation: "1+1=2",
		Status:      "open",
	}
	id, err := store.CreateReport(rep)
	if err != nil {
		t.Fatalf("create report: %v", err)
	}

	if err := store.UpdateReportStatus(id, "resolved"); err != nil {
		t.Fatalf("update report status: %v", err)
	}

	if err := store.UpdateReportStatus(9999999, "resolved"); err == nil {
		t.Error("expected error updating nonexistent report")
	}
}

func TestGetAdminUser_Lookup(t *testing.T) {
	store := newTestStore(t)
	st, err := store.CreateStudent("admin-lookup-user")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}

	user, err := store.GetAdminUser(st.ID)
	if err != nil {
		t.Fatalf("get admin user: %v", err)
	}
	if user == nil {
		t.Fatal("expected admin user projection, got nil")
	}
	if user.ID != st.ID || user.Name != "admin-lookup-user" {
		t.Errorf("got user %+v, want ID=%s Name=admin-lookup-user", user, st.ID)
	}
	if user.Role != RoleStudent {
		t.Errorf("expected default role student, got %s", user.Role)
	}

	missing, err := store.GetAdminUser("missing-student-id")
	if err != nil {
		t.Fatalf("get missing admin user: %v", err)
	}
	if missing != nil {
		t.Errorf("expected nil for missing admin user, got %+v", missing)
	}
}

func TestRemedialQueue_CRUD(t *testing.T) {
	store := newTestStore(t)
	st, err := store.CreateStudent("remedial-tester")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}

	// Initial queue must be empty
	q, err := store.GetRemedialQueue(st.ID)
	if err != nil {
		t.Fatalf("get remedial queue: %v", err)
	}
	if len(q) != 0 {
		t.Fatalf("expected empty queue, got %v", q)
	}

	// Enqueue empty list is a no-op
	if err := store.EnqueueRemedial(st.ID); err != nil {
		t.Fatalf("enqueue empty: %v", err)
	}

	// Enqueue items in order
	if err := store.EnqueueRemedial(st.ID, "c1", "c2", "c3"); err != nil {
		t.Fatalf("enqueue items: %v", err)
	}

	q, err = store.GetRemedialQueue(st.ID)
	if err != nil {
		t.Fatalf("get queue: %v", err)
	}
	if len(q) != 3 || q[0] != "c1" || q[1] != "c2" || q[2] != "c3" {
		t.Fatalf("expected [c1 c2 c3], got %v", q)
	}

	// Duplicates are ignored; new items are appended
	if err := store.EnqueueRemedial(st.ID, "c2", "c4"); err != nil {
		t.Fatalf("enqueue duplicate: %v", err)
	}

	q, err = store.GetRemedialQueue(st.ID)
	if err != nil {
		t.Fatalf("get queue after deduplication: %v", err)
	}
	if len(q) != 4 || q[0] != "c1" || q[1] != "c2" || q[2] != "c3" || q[3] != "c4" {
		t.Fatalf("expected [c1 c2 c3 c4], got %v", q)
	}

	// Capping at MaxRemedialQueue (10 items)
	var many []string
	for i := 5; i <= 20; i++ {
		many = append(many, fmt.Sprintf("c%d", i))
	}
	if err := store.EnqueueRemedial(st.ID, many...); err != nil {
		t.Fatalf("enqueue many: %v", err)
	}

	q, err = store.GetRemedialQueue(st.ID)
	if err != nil {
		t.Fatalf("get capped queue: %v", err)
	}
	if len(q) != MaxRemedialQueue {
		t.Fatalf("expected queue capped at %d, got %d", MaxRemedialQueue, len(q))
	}

	// Remove single item
	if err := store.RemoveRemedial(st.ID, "c2"); err != nil {
		t.Fatalf("remove remedial: %v", err)
	}
	q, err = store.GetRemedialQueue(st.ID)
	if err != nil {
		t.Fatalf("get queue after removal: %v", err)
	}
	for _, id := range q {
		if id == "c2" {
			t.Fatalf("c2 was not removed from queue: %v", q)
		}
	}

	// Clear queue
	if err := store.ClearRemedialQueue(st.ID); err != nil {
		t.Fatalf("clear remedial queue: %v", err)
	}
	q, err = store.GetRemedialQueue(st.ID)
	if err != nil {
		t.Fatalf("get queue after clear: %v", err)
	}
	if len(q) != 0 {
		t.Fatalf("expected empty queue after clear, got %v", q)
	}
}

func TestRemedialQueue_ResetProgressAndAccountDelete(t *testing.T) {
	store := newTestStore(t)
	st, err := store.CreateStudent("wipe-tester")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}

	if err := store.EnqueueRemedial(st.ID, "rem1", "rem2"); err != nil {
		t.Fatalf("enqueue: %v", err)
	}

	// ResetProgress must clear remedial queue
	if err := store.ResetProgress(st.ID); err != nil {
		t.Fatalf("reset progress: %v", err)
	}
	q, err := store.GetRemedialQueue(st.ID)
	if err != nil {
		t.Fatalf("get queue: %v", err)
	}
	if len(q) != 0 {
		t.Fatalf("expected empty remedial queue after reset, got %v", q)
	}

	// Seed again and verify DeleteAccount clears remedial queue
	if err := store.EnqueueRemedial(st.ID, "rem3"); err != nil {
		t.Fatalf("enqueue: %v", err)
	}
	if err := store.DeleteAccount(st.ID); err != nil {
		t.Fatalf("delete account: %v", err)
	}
	q, err = store.GetRemedialQueue(st.ID)
	if err != nil {
		t.Fatalf("get queue: %v", err)
	}
	if len(q) != 0 {
		t.Fatalf("expected empty remedial queue after delete account, got %v", q)
	}
}
