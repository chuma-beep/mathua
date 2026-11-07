package engine

import (
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/storage"
	xppolicy "github.com/chuma-beep/mathua/internal/xp"
)

func seedReviewSignals(t *testing.T, e *Engine, studentID, runID string, difficulty float64) {
	t.Helper()
	store, ok := e.repo.(storage.EffortStore)
	if !ok {
		t.Fatal("test store does not implement EffortStore")
	}
	now := time.Now().UTC()
	for i := 0; i < 6; i++ {
		if err := store.RecordEffortSignal(storage.EffortSignal{
			StudentID:       studentID,
			TaskKind:        xppolicy.TaskReview,
			TaskID:          runID + ":concept",
			ReviewRunID:     runID,
			ConceptID:       "concept",
			AttemptKey:      runID + ":" + string(rune('a'+i)),
			Correct:         false,
			ElapsedSeconds:  10,
			Difficulty:      difficulty,
			XPEarned:        0,
			AnswerSignature: string(rune('a' + i)),
			CreatedAt:       now.Add(time.Duration(i) * time.Second),
		}); err != nil {
			t.Fatalf("record signal %d: %v", i, err)
		}
	}
}

func TestFinalizeReviewPerformanceAppliesOnce(t *testing.T) {
	e := testEngine(t)
	st, err := e.CreateStudent("review-finalize")
	if err != nil {
		t.Fatalf("student: %v", err)
	}
	store := e.repo.(*storage.SQLiteStore)
	if err := store.AwardXP(st.ID, 10, "seed", "seed:review"); err != nil {
		t.Fatalf("seed xp: %v", err)
	}
	const runID = "review-run-1"
	if err := e.EnsureReviewRun(st.ID, runID); err != nil {
		t.Fatalf("ensure run: %v", err)
	}
	seedReviewSignals(t, e, st.ID, runID, 0.3)

	first, err := e.FinalizeReviewPerformance(st.ID, runID)
	if err != nil {
		t.Fatalf("finalize: %v", err)
	}
	if first.XPAdjustment != -2 || first.XPSource != "performance" || first.XPNote == "" {
		t.Fatalf("first finalization = %+v, want -2 performance", first)
	}
	total, _, err := store.GetXP(st.ID)
	if err != nil {
		t.Fatalf("get xp: %v", err)
	}
	if total != 8 {
		t.Fatalf("balance = %d, want 8", total)
	}

	second, err := e.FinalizeReviewPerformance(st.ID, runID)
	if err != nil {
		t.Fatalf("second finalize: %v", err)
	}
	if second.XPAdjustment != -2 || second.XPNote != first.XPNote {
		t.Fatalf("second finalization = %+v, want stored outcome", second)
	}
	total, _, err = store.GetXP(st.ID)
	if err != nil {
		t.Fatalf("get xp after replay: %v", err)
	}
	if total != 8 {
		t.Fatalf("replay moved balance to %d", total)
	}
}

func TestFinalizeReviewPerformanceSkipsWhenEffortDeducted(t *testing.T) {
	e := testEngine(t)
	st, err := e.CreateStudent("review-effort")
	if err != nil {
		t.Fatalf("student: %v", err)
	}
	store := e.repo.(*storage.SQLiteStore)
	es := storage.EffortStore(store)
	if err := store.AwardXP(st.ID, 10, "seed", "seed:review-effort"); err != nil {
		t.Fatalf("seed xp: %v", err)
	}
	const runID = "review-run-effort"
	if err := e.EnsureReviewRun(st.ID, runID); err != nil {
		t.Fatalf("ensure run: %v", err)
	}
	seedReviewSignals(t, e, st.ID, runID, 0.3)
	if _, err := es.ApplyXPAdjustment(st.ID, -1, "accountability:insufficient_effort", "accountability:review:"+runID+":concept"); err != nil {
		t.Fatalf("effort deduction: %v", err)
	}

	result, err := e.FinalizeReviewPerformance(st.ID, runID)
	if err != nil {
		t.Fatalf("finalize: %v", err)
	}
	if result.XPAdjustment != 0 || result.XPNote != "" {
		t.Fatalf("performance after effort = %+v, want no adjustment", result)
	}
	total, _, err := store.GetXP(st.ID)
	if err != nil {
		t.Fatalf("get xp: %v", err)
	}
	if total != 9 {
		t.Fatalf("balance = %d, want 9 after only the effort deduction", total)
	}
}

func TestFinalizeReviewPerformanceHardTaskIsExempt(t *testing.T) {
	e := testEngine(t)
	st, err := e.CreateStudent("review-hard")
	if err != nil {
		t.Fatalf("student: %v", err)
	}
	store := e.repo.(*storage.SQLiteStore)
	if err := store.AwardXP(st.ID, 10, "seed", "seed:review-hard"); err != nil {
		t.Fatalf("seed xp: %v", err)
	}
	const runID = "review-run-hard"
	if err := e.EnsureReviewRun(st.ID, runID); err != nil {
		t.Fatalf("ensure run: %v", err)
	}
	seedReviewSignals(t, e, st.ID, runID, 0.8)

	result, err := e.FinalizeReviewPerformance(st.ID, runID)
	if err != nil {
		t.Fatalf("finalize: %v", err)
	}
	if result.XPAdjustment != 0 {
		t.Fatalf("hard task = %+v, want exemption", result)
	}
	total, _, err := store.GetXP(st.ID)
	if err != nil {
		t.Fatalf("get xp: %v", err)
	}
	if total != 10 {
		t.Fatalf("balance = %d, want 10", total)
	}
}
