package engine

import (
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/storage"
)

// TestTimedQuizPerformanceAdjustmentIsAtomicAndReplayable pins the idempotent
// answer transaction used by timed quizzes. A recovery window there could
// store a response claiming an adjustment the ledger never held. The committed
// response and the ledger must agree, and a replay must not charge twice.
func TestTimedQuizPerformanceAdjustmentIsAtomicAndReplayable(t *testing.T) {
	d, err := concepts.Build([]concepts.Concept{{
		ID: "a", Label: "A", Domain: "d", GradingType: "numeric",
		Prerequisites:    []string{},
		MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60},
	}})
	if err != nil {
		t.Fatalf("build DAG: %v", err)
	}
	store, err := storage.NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	reg := generator.NewRegistry()
	reg.Register("a", &testGen{answer: "42"})
	e := New(store, d, reg, nil, nil)
	st, err := e.CreateStudent("performance")
	if err != nil {
		t.Fatalf("student: %v", err)
	}
	if err := store.AwardXP(st.ID, 10, "seed", "seed:1"); err != nil {
		t.Fatalf("seed xp: %v", err)
	}
	e.SetStudyAnchor(st.ID, "a", "42")

	var last *AnswerResult
	lastSubmission := ""
	for i := 0; i < 6; i++ {
		lastSubmission = "sub-" + string(rune('a'+i))
		last, err = e.SubmitTimedQuizAnswerWithCompletion(st.ID, "a", "wrong", "42", 10, "Q1", nil, lastSubmission, "run-1", i == 5)
		if err != nil {
			t.Fatalf("answer %d: %v", i, err)
		}
	}
	if last == nil || last.XPAdjustment != -2 {
		t.Fatalf("last adjustment = %+v, want -2", last)
	}
	if !strings.Contains(last.ResponseJSON, `"xp_adjustment":-2`) || !strings.Contains(last.ResponseJSON, `"xp_source":"performance"`) {
		t.Fatalf("stored response does not carry the actual adjustment: %s", last.ResponseJSON)
	}

	total, _, err := store.GetXP(st.ID)
	if err != nil {
		t.Fatalf("get xp: %v", err)
	}
	if total != 8 {
		t.Fatalf("xp total = %d, want 8 after a 10 XP seed and a -2 deduction", total)
	}

	replay, err := e.SubmitTimedQuizAnswerWithCompletion(st.ID, "a", "wrong", "42", 10, "Q1", nil, lastSubmission, "run-1", true)
	if err != nil {
		t.Fatalf("replay: %v", err)
	}
	totalAfter, _, err := store.GetXP(st.ID)
	if err != nil {
		t.Fatalf("get xp after replay: %v", err)
	}
	if totalAfter != total {
		t.Fatalf("replay moved XP from %d to %d", total, totalAfter)
	}
	if !strings.Contains(replay.ResponseJSON, `"xp_adjustment":-2`) {
		t.Fatalf("replay lost the stored explanation: %s", replay.ResponseJSON)
	}
}

func TestTimedQuizAdjustmentAppliesFullCalculatedAmount(t *testing.T) {
	e := testEngine(t)
	st, err := e.CreateStudent("capped")
	if err != nil {
		t.Fatalf("student: %v", err)
	}
	store := e.repo.(*storage.SQLiteStore)
	if err := store.AwardXP(st.ID, 1, "seed", "seed:capped"); err != nil {
		t.Fatalf("seed xp: %v", err)
	}
	var last *AnswerResult
	for i := 0; i < 6; i++ {
		last, err = e.SubmitTimedQuizAnswerWithCompletion(st.ID, "a", "wrong", "42", 10, "Q1", nil, "capped-"+string(rune('a'+i)), "run-capped", i == 5)
		if err != nil {
			t.Fatalf("answer %d: %v", i, err)
		}
	}
	if last == nil || last.XPAdjustment != -2 {
		t.Fatalf("last = %+v, want actual -2", last)
	}
	if !strings.Contains(last.XPNote, "-2 XP") || last.XPSource != "performance" {
		t.Fatalf("note/source = %q/%q, want actual -2 from performance", last.XPNote, last.XPSource)
	}
	total, _, err := store.GetXP(st.ID)
	if err != nil {
		t.Fatalf("get xp: %v", err)
	}
	if total != -1 {
		t.Fatalf("balance = %d, want -1", total)
	}
}

func TestQuizPerformanceWaitsForFinalAnswer(t *testing.T) {
	e := testEngine(t)
	st, err := e.CreateStudent("quiz-complete")
	if err != nil {
		t.Fatalf("student: %v", err)
	}
	store := e.repo.(*storage.SQLiteStore)
	if err := store.AwardXP(st.ID, 10, "seed", "seed:quiz-complete"); err != nil {
		t.Fatalf("seed xp: %v", err)
	}
	var last *AnswerResult
	for i := 0; i < 6; i++ {
		last, err = e.SubmitQuizAnswerWithCompletion(st.ID, "a", "wrong", "42", 10, "Q1", nil, "quiz-session-1", i == 5)
		if err != nil {
			t.Fatalf("quiz answer %d: %v", i, err)
		}
		if i < 5 && last.XPAdjustment != 0 {
			t.Fatalf("answer %d produced a premature performance adjustment: %+v", i, last)
		}
	}
	if last == nil || last.XPAdjustment != -2 || last.XPSource != "performance" {
		t.Fatalf("final quiz answer = %+v, want -2 performance", last)
	}
}

func TestStudyPerformanceWaitsForCompletion(t *testing.T) {
	e := testEngine(t)
	st, err := e.CreateStudent("unfinished")
	if err != nil {
		t.Fatalf("student: %v", err)
	}
	store := e.repo.(*storage.SQLiteStore)
	if err := store.AwardXP(st.ID, 10, "seed", "seed:unfinished"); err != nil {
		t.Fatalf("seed xp: %v", err)
	}
	e.SetStudyAnchor(st.ID, "a", "42")

	var last *AnswerResult
	for i := 0; i < 6; i++ {
		last, err = e.SubmitStudyAnswerWithSubmission(st.ID, "a", "wrong", 10, "Q1", "unfinished-"+string(rune('a'+i)))
		if err != nil {
			t.Fatalf("answer %d: %v", i, err)
		}
	}
	if last == nil || last.XPAdjustment != 0 {
		t.Fatalf("unfinished study task = %+v, want no performance adjustment", last)
	}
	total, _, err := store.GetXP(st.ID)
	if err != nil {
		t.Fatalf("get xp: %v", err)
	}
	if total != 10 {
		t.Fatalf("unfinished study moved XP to %d, want 10", total)
	}
}

func TestQuizGateCountsEarnedXPNotAccountabilityDeductions(t *testing.T) {
	e := testEngine(t)
	st, err := e.CreateStudent("gate")
	if err != nil {
		t.Fatalf("student: %v", err)
	}
	store := e.repo.(*storage.SQLiteStore)
	es := storage.EffortStore(store)
	if err := store.AwardXP(st.ID, 60, "seed", "seed:gate"); err != nil {
		t.Fatalf("seed: %v", err)
	}
	if err := e.RecordQuizCompletion(st.ID); err != nil {
		t.Fatalf("quiz completion: %v", err)
	}
	if err := store.AwardXP(st.ID, 10, "study", "award:after-quiz"); err != nil {
		t.Fatalf("award: %v", err)
	}
	if _, err := es.ApplyXPAdjustment(st.ID, -5, "accountability:poor_performance", "accountability:practice:gate"); err != nil {
		t.Fatalf("deduction: %v", err)
	}
	if err := store.AwardXP(st.ID, -5, "study", "rush:gate"); err != nil {
		t.Fatalf("legacy rush deduction: %v", err)
	}
	since, err := e.QuizXPSince(st.ID)
	if err != nil {
		t.Fatalf("quiz xp since: %v", err)
	}
	if since != 10 {
		t.Fatalf("QuizXPSince = %d, want 10 earned XP even after a 5 XP deduction", since)
	}
}
