package scoring

import (
	"encoding/json"
	"math"
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/levels"
	"github.com/chuma-beep/mathua/internal/storage"
	"github.com/chuma-beep/mathua/internal/xp"
)

func TestComputeLevel(t *testing.T) {
	cases := []struct {
		mastered int
		expected string
	}{
		{0, "Novice"},
		{31, "Novice"},
		{32, "Apprentice"},
		{64, "Student"},
		{96, "Scholar"},
		{128, "Adept"},
		{160, "Expert"},
		{192, "Master"},
		{224, "Grandmaster"},
		{256, "Math Architect"},
		{284, "Math Architect"},
	}
	for _, tc := range cases {
		got := levels.Compute(tc.mastered)
		if got != tc.expected {
			t.Errorf("mastered=%d: expected %q, got %q", tc.mastered, tc.expected, got)
		}
	}
}

func TestComputeSpeedBonus(t *testing.T) {
	p := &storage.ConceptProgress{AvgResponseTime: 5.0}
	bonus := computeSpeedBonus(p)
	if bonus != 10 {
		t.Errorf("expected 10, got %f", bonus)
	}
	p.AvgResponseTime = 0
	if computeSpeedBonus(p) != 0 {
		t.Error("expected 0 for zero time")
	}
	p.AvgResponseTime = 0.5
	bonus = computeSpeedBonus(p)
	if bonus != 25 {
		t.Errorf("expected capped at 25, got %f", bonus)
	}
}

func TestComputeCurrentStreak_Empty(t *testing.T) {
	s := computeCurrentStreak(map[string]*storage.ConceptProgress{})
	if s != 0 {
		t.Errorf("expected 0, got %d", s)
	}
}

func TestComputeCurrentStreak_Today(t *testing.T) {
	now := time.Now().UTC()
	s := computeCurrentStreak(map[string]*storage.ConceptProgress{
		"a": {LastAttempted: &now},
	})
	if s != 1 {
		t.Errorf("expected 1, got %d", s)
	}
}

func TestComputeCurrentStreak_Yesterday(t *testing.T) {
	yesterday := time.Now().UTC().AddDate(0, 0, -1)
	s := computeCurrentStreak(map[string]*storage.ConceptProgress{
		"a": {LastAttempted: &yesterday},
	})
	// No activity today, so streak is 0.
	if s != 0 {
		t.Errorf("expected 0 (no today), got %d", s)
	}
}

func TestCompute(t *testing.T) {
	d, _ := concepts.Load("../../data/concepts.json")
	store, _ := storage.NewSQLiteStore(":memory:")
	defer store.Close()
	st, _ := store.CreateStudent("tester")
	now := time.Now().UTC()
	_ = store.UpsertProgress(&storage.ConceptProgress{
		StudentID: st.ID, ConceptID: "count.basics.objects",
		Status: "MASTERED", AvgResponseTime: 5.0, MasteredAt: &now, LastAttempted: &now,
	})
	_ = store.UpsertProgress(&storage.ConceptProgress{
		StudentID: st.ID, ConceptID: "count.basics.cardinality",
		Status: "LEARNING", AvgResponseTime: 8.0, LastAttempted: &now,
	})

	u := NewUpdater(d, store)
	scores, err := u.Compute(st.ID)
	if err != nil {
		t.Fatalf("compute: %v", err)
	}
	if scores.ConceptsMastered != 1 {
		t.Errorf("expected 1 mastered, got %d", scores.ConceptsMastered)
	}
	if scores.Level != "Novice" {
		t.Errorf("expected Novice, got %s", scores.Level)
	}
	if scores.WeeklyScore == 0 {
		t.Errorf("expected non-zero weekly score, got %d", scores.WeeklyScore)
	}
}

func TestWeekStart(t *testing.T) {
	// 2026-05-12 is a Tuesday. Monday is May 11.
	ts := time.Date(2026, 5, 12, 15, 30, 0, 0, time.UTC)
	m := weekStart(ts)
	if m.Year() != 2026 || m.Month() != 5 || m.Day() != 11 {
		t.Errorf("expected 2026-05-11, got %s", m.Format("2006-01-02"))
	}
	if m.Hour() != 0 || m.Minute() != 0 {
		t.Error("expected midnight")
	}
}

// The gate was a client-side literal in three places, the copy, the progress bar, and an
// `xp_total >= 50` fallback for a flag the server does send. The server carries it now, and
// this asserts the path the page actually reads. A zero-valued Scores would pass a field
// comparison, because QuizGateXP defaults to 0 rather than the gate, so the test drives Compute.
func TestScoresCarryTheQuizGate(t *testing.T) {
	d, _ := concepts.Load("../../data/concepts.json")
	store, _ := storage.NewSQLiteStore(":memory:")
	defer store.Close()
	st, _ := store.CreateStudent("tester")

	u := NewUpdater(d, store)
	scores, err := u.Compute(st.ID)
	if err != nil {
		t.Fatalf("compute: %v", err)
	}
	if scores.QuizGateXP != xp.QuizGateXP {
		t.Errorf("QuizGateXP = %d, want xp.QuizGateXP = %d", scores.QuizGateXP, xp.QuizGateXP)
	}

	b, err := json.Marshal(scores)
	if err != nil {
		t.Fatal(err)
	}
	var round map[string]any
	if err := json.Unmarshal(b, &round); err != nil {
		t.Fatal(err)
	}
	// Optional in the client schema, so absence has to be visible as absence rather than
	// silently defaulted, that leniency is what let the page keep its own copy for so long.
	v, ok := round["quiz_gate_xp"]
	if !ok {
		t.Fatal("quiz_gate_xp missing from the Scores payload")
	}
	if got := v.(float64); got != float64(xp.QuizGateXP) {
		t.Errorf("quiz_gate_xp = %v, want %d", got, xp.QuizGateXP)
	}
}

// The fallback daily goal for a learner with no stored target was a literal 10 while
// xp.DefaultDailyGoal and the planner said 30 (code-defects.md §1). This asserts that
// an unconfigured learner evaluates against xp.DefaultDailyGoal, and that an explicit
// goal takes precedence.
func TestScoresDailyGoalFallback(t *testing.T) {
	d, _ := concepts.Load("../../data/concepts.json")
	store, _ := storage.NewSQLiteStore(":memory:")
	defer store.Close()
	st, _ := store.CreateStudent("tester")

	u := NewUpdater(d, store)
	scores, err := u.Compute(st.ID)
	if err != nil {
		t.Fatalf("compute: %v", err)
	}
	if scores.DailyXPGoal != xp.DefaultDailyGoal {
		t.Errorf("DailyXPGoal = %d, want xp.DefaultDailyGoal = %d", scores.DailyXPGoal, xp.DefaultDailyGoal)
	}

	if err := store.SetDailyXPGoal(st.ID, 50); err != nil {
		t.Fatalf("set daily goal: %v", err)
	}
	scoresExplicit, err := u.Compute(st.ID)
	if err != nil {
		t.Fatalf("compute explicit: %v", err)
	}
	if scoresExplicit.DailyXPGoal != 50 {
		t.Errorf("DailyXPGoal = %d, want 50", scoresExplicit.DailyXPGoal)
	}
}

func TestWeekStart_Sunday(t *testing.T) {
	// 2026-05-17 is a Sunday. Monday of that week is May 11.
	ts := time.Date(2026, 5, 17, 22, 15, 0, 0, time.UTC)
	m := weekStart(ts)
	if m.Year() != 2026 || m.Month() != 5 || m.Day() != 11 {
		t.Errorf("expected 2026-05-11, got %s", m.Format("2006-01-02"))
	}
	if m.Hour() != 0 || m.Minute() != 0 {
		t.Error("expected midnight")
	}
}

func TestComputeCurrentStreak_ConsecutiveAndGap(t *testing.T) {
	now := time.Now().UTC()
	yesterday := now.AddDate(0, 0, -1)
	twoDaysAgo := now.AddDate(0, 0, -2)
	threeDaysAgo := now.AddDate(0, 0, -3)

	// 3 consecutive days: today, yesterday and two days ago -> streak 3
	progress3 := map[string]*storage.ConceptProgress{
		"a": {LastAttempted: &now},
		"b": {LastAttempted: &yesterday},
		"c": {LastAttempted: &twoDaysAgo},
	}
	if s := computeCurrentStreak(progress3); s != 3 {
		t.Errorf("expected streak 3, got %d", s)
	}

	// Gap between yesterday and three days ago -> streak 2
	progressGap := map[string]*storage.ConceptProgress{
		"a": {LastAttempted: &now},
		"b": {LastAttempted: &yesterday},
		"c": {LastAttempted: &threeDaysAgo},
	}
	if s := computeCurrentStreak(progressGap); s != 2 {
		t.Errorf("expected streak 2, got %d", s)
	}
}

func TestPausedUntilCases(t *testing.T) {
	store, _ := storage.NewSQLiteStore(":memory:")
	defer store.Close()
	st, _ := store.CreateStudent("paused-tester")
	u := NewUpdater(nil, store)

	// Case 1: empty settings
	if got := pausedUntil(u, st.ID); got != "" {
		t.Errorf("expected empty string for empty settings, got %q", got)
	}

	// Case 2: literal "{}"
	_ = store.UpdateSettings(st.ID, "{}")
	if got := pausedUntil(u, st.ID); got != "" {
		t.Errorf("expected empty string for {}, got %q", got)
	}

	// Case 3: invalid json
	_ = store.UpdateSettings(st.ID, "{bad json")
	if got := pausedUntil(u, st.ID); got != "" {
		t.Errorf("expected empty string for malformed json, got %q", got)
	}

	// Case 4: empty pause_until
	_ = store.UpdateSettings(st.ID, `{"pause_until":""}`)
	if got := pausedUntil(u, st.ID); got != "" {
		t.Errorf("expected empty string for empty pause_until, got %q", got)
	}

	// Case 5: unparseable date
	_ = store.UpdateSettings(st.ID, `{"pause_until":"not-a-date"}`)
	if got := pausedUntil(u, st.ID); got != "" {
		t.Errorf("expected empty string for unparseable date, got %q", got)
	}

	// Case 6: past date
	_ = store.UpdateSettings(st.ID, `{"pause_until":"2020-01-01"}`)
	if got := pausedUntil(u, st.ID); got != "" {
		t.Errorf("expected empty string for past date, got %q", got)
	}

	// Case 7: future date
	future := time.Now().UTC().AddDate(1, 0, 0).Format("2006-01-02")
	_ = store.UpdateSettings(st.ID, `{"pause_until":"`+future+`"}`)
	if got := pausedUntil(u, st.ID); got != future {
		t.Errorf("expected future date %q, got %q", future, got)
	}
}

func TestCompute_TopicSpeedsAndQuizGate(t *testing.T) {
	store, _ := storage.NewSQLiteStore(":memory:")
	defer store.Close()
	st, _ := store.CreateStudent("quiz-tester")
	_ = store.SetDailyXPGoal(st.ID, 40)
	_ = store.AwardXP(st.ID, 50, "practice", "seed-xp-1")
	_ = store.RecordQuizCompletion(st.ID, 50)
	_ = store.AwardXP(st.ID, 70, "practice", "seed-xp-2")

	// Topic speed 1: high repetitions/speed causing darkness clamping to 1.0
	_ = store.UpsertTopicSpeed(&storage.TopicSpeed{
		StudentID:     st.ID,
		ConceptID:     "c1",
		LearningSpeed: 2.5,
		Repetitions:   8,
	})
	// Topic speed 2: normal repetitions/speed
	_ = store.UpsertTopicSpeed(&storage.TopicSpeed{
		StudentID:     st.ID,
		ConceptID:     "c2",
		LearningSpeed: 1.2,
		Repetitions:   2,
	})

	// Completion happened at 50 XP, then 70 XP was earned. Current is 120.
	// XPSinceQuiz = 70 >= xp.QuizGateXP (50) -> QuizDue is true.
	u := NewUpdater(nil, store)
	scores, err := u.Compute(st.ID)
	if err != nil {
		t.Fatalf("compute: %v", err)
	}

	if scores.DailyXPGoal != 40 {
		t.Errorf("DailyXPGoal = %d, want 40", scores.DailyXPGoal)
	}
	if scores.XPTotal != 120 {
		t.Errorf("XPTotal = %d, want 120", scores.XPTotal)
	}
	if scores.XPSinceQuiz != 70 {
		t.Errorf("XPSinceQuiz = %d, want 70", scores.XPSinceQuiz)
	}
	if !scores.QuizDue {
		t.Errorf("QuizDue expected true (70 >= %d)", xp.QuizGateXP)
	}
	if scores.SpacedReps["c1"] != 1.0 {
		t.Errorf("SpacedReps[c1] = %f, want clamped 1.0", scores.SpacedReps["c1"])
	}
	expectedC2 := math.Round(1.2*(0.3+float64(2)*0.15)*100) / 100
	if scores.SpacedReps["c2"] != expectedC2 {
		t.Errorf("SpacedReps[c2] = %f, want %f", scores.SpacedReps["c2"], expectedC2)
	}

	// Test negative xpSinceQuiz clamping (e.g. if last quiz recorded higher xp total than current)
	_ = store.RecordQuizCompletion(st.ID, 200)
	scores2, err := u.Compute(st.ID)
	if err != nil {
		t.Fatalf("compute: %v", err)
	}
	if scores2.XPSinceQuiz != 0 {
		t.Errorf("XPSinceQuiz = %d, want clamped 0", scores2.XPSinceQuiz)
	}
	if scores2.QuizDue {
		t.Errorf("QuizDue expected false when XPSinceQuiz is 0")
	}
}

func TestCompute_StorageError(t *testing.T) {
	store, _ := storage.NewSQLiteStore(":memory:")
	st, _ := store.CreateStudent("err-tester")
	u := NewUpdater(nil, store)
	_ = store.Close() // close store to induce error

	if _, err := u.Compute(st.ID); err == nil {
		t.Error("expected error from Compute with closed store")
	}
}
