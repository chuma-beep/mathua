package scoring

import (
	"testing"

	"github.com/chuma-beep/mathua/internal/storage"
)

func TestQuizEligibilityUsesEarnedXPWhenBalanceIsNegative(t *testing.T) {
	store, err := storage.NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	st, err := store.CreateStudent("negative-gate")
	if err != nil {
		t.Fatalf("student: %v", err)
	}
	if err := store.AwardXP(st.ID, 60, "seed", "seed:gate"); err != nil {
		t.Fatalf("seed: %v", err)
	}
	if err := store.RecordQuizCompletion(st.ID, 60); err != nil {
		t.Fatalf("quiz completion: %v", err)
	}
	if err := store.AwardXP(st.ID, 10, "study", "award:after-quiz"); err != nil {
		t.Fatalf("award after quiz: %v", err)
	}
	es := storage.EffortStore(store)
	if _, err := es.ApplyXPAdjustment(st.ID, -20, "accountability:extremely_poor_performance", "accountability:quiz:gate"); err != nil {
		t.Fatalf("deduction: %v", err)
	}
	total, _, err := store.GetXP(st.ID)
	if err != nil {
		t.Fatalf("get xp: %v", err)
	}
	if total != 50 {
		t.Fatalf("balance = %d, want 50", total)
	}

	earned, err := storage.XPSinceLastQuiz(store, st.ID)
	if err != nil {
		t.Fatalf("storage gate: %v", err)
	}
	if earned != 10 {
		t.Fatalf("earned since quiz = %d, want 10", earned)
	}
	scores, err := NewUpdater(nil, store).Compute(st.ID)
	if err != nil {
		t.Fatalf("compute: %v", err)
	}
	if scores.XPSinceQuiz != earned || scores.QuizDue != (earned >= scores.QuizGateXP) {
		t.Fatalf("dashboard gate = %d due=%v, storage = %d", scores.XPSinceQuiz, scores.QuizDue, earned)
	}
}
