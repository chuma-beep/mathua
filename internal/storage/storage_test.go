package storage

import (
	"strings"
	"testing"
	"time"
)

func newTestStore(t *testing.T) *SQLiteStore {
	t.Helper()
	store, err := NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("create store: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	return store
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
	if err := store.RecordAttempt(entry); err != nil {
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
	members, err := store.GetLeagueStandings()
	if err != nil {
		t.Fatalf("standings: %v", err)
	}
	if len(members) != 1 || members[0].Username != "ada" {
		t.Fatalf("expected username on member, got %+v", members)
	}
	if members[0].AvatarURL != "https://photo/x.jpg" || members[0].AvatarDicebear == nil {
		t.Errorf("expected avatar fields on member, got %+v", members[0])
	}
	if c, p, v := AvatarBits("not-json"); c || p != nil || v != 0 {
		t.Errorf("corrupt settings should yield zero bits")
	}
	if c, p, v := AvatarBits(""); c || p != nil || v != 0 {
		t.Errorf("empty settings should yield zero bits")
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
