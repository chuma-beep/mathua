package storage

import (
	"fmt"
	"testing"
	"time"
)

func TestReviewRunFinalizationIsAtomicAndIdempotent(t *testing.T) {
	store := newTestStore(t)
	es := store.(EffortStore)
	rs := store.(ReviewRunStore)
	st := seedStudent(t, store)
	if err := store.AwardXP(st.ID, 10, "seed", "seed:review-run"); err != nil {
		t.Fatalf("seed xp: %v", err)
	}
	const runID = "review-run-1"
	if err := rs.EnsureReviewRun(runID, st.ID, time.Now().UTC()); err != nil {
		t.Fatalf("ensure run: %v", err)
	}
	for i := 0; i < 6; i++ {
		if err := es.RecordEffortSignal(EffortSignal{
			StudentID: st.ID, TaskKind: "review", TaskID: fmt.Sprintf("%s:c%d", runID, i),
			ReviewRunID: runID, AttemptKey: fmt.Sprintf("attempt-%d", i),
			Correct: false, ElapsedSeconds: 10, CreatedAt: time.Now().UTC(),
		}); err != nil {
			t.Fatalf("signal %d: %v", i, err)
		}
	}
	signals, err := rs.ReviewRunSignals(st.ID, runID)
	if err != nil || len(signals) != 6 {
		t.Fatalf("signals = %d, err = %v; want 6", len(signals), err)
	}
	if n, err := rs.CountReviewRunEffortDeductions(st.ID, runID); err != nil || n != 0 {
		t.Fatalf("effort deductions = %d, %v; want 0", n, err)
	}

	fin := ReviewFinalization{
		At: time.Now().UTC(), Amount: -2,
		Reason:   "accountability:extremely_poor_performance",
		EventKey: "accountability:review-run:" + runID,
		NoteForApplied: func(applied int) (string, string, string) {
			return fmt.Sprintf("%d XP: extremely poor review", applied), "extremely_poor_performance", "performance"
		},
	}
	run, now, err := rs.FinalizeReviewRun(st.ID, runID, fin)
	if err != nil || !now {
		t.Fatalf("finalize = %+v, %v, %v; want newly finalized", run, now, err)
	}
	if run.XPAdjustment != -2 || run.XPNote == "" || run.XPReasonCode != "extremely_poor_performance" || run.XPSource != "performance" {
		t.Fatalf("stored run = %+v", run)
	}
	total, _, err := store.GetXP(st.ID)
	if err != nil || total != 8 {
		t.Fatalf("balance = %d, %v; want 8", total, err)
	}

	again, now, err := rs.FinalizeReviewRun(st.ID, runID, fin)
	if err != nil || now {
		t.Fatalf("second finalize = %+v, %v, %v; want stored, false", again, now, err)
	}
	if again.XPAdjustment != -2 || again.XPNote == "" {
		t.Fatalf("second finalize lost stored outcome: %+v", again)
	}
	total, _, err = store.GetXP(st.ID)
	if err != nil || total != 8 {
		t.Fatalf("replay moved balance to %d, %v", total, err)
	}
}

func TestReviewRunFinalizationCanPushBalanceNegative(t *testing.T) {
	store := newTestStore(t)
	rs := store.(ReviewRunStore)
	st := seedStudent(t, store)
	if err := store.AwardXP(st.ID, 1, "seed", "seed:review-negative"); err != nil {
		t.Fatalf("seed: %v", err)
	}
	const runID = "review-run-negative"
	if err := rs.EnsureReviewRun(runID, st.ID, time.Now().UTC()); err != nil {
		t.Fatalf("ensure run: %v", err)
	}
	run, now, err := rs.FinalizeReviewRun(st.ID, runID, ReviewFinalization{
		At: time.Now().UTC(), Amount: -2,
		Reason:   "accountability:extremely_poor_performance",
		EventKey: "accountability:review-run:" + runID,
		NoteForApplied: func(applied int) (string, string, string) {
			return fmt.Sprintf("%d XP: negative review", applied), "extremely_poor_performance", "performance"
		},
	})
	if err != nil || !now || run.XPAdjustment != -2 {
		t.Fatalf("finalize = %+v, %v, %v; want -2", run, now, err)
	}
	total, _, err := store.GetXP(st.ID)
	if err != nil || total != -1 {
		t.Fatalf("balance = %d, %v; want -1", total, err)
	}
}

func TestReviewRunFinalizationWithZeroAdjustmentStoresExplanation(t *testing.T) {
	store := newTestStore(t)
	rs := store.(ReviewRunStore)
	st := seedStudent(t, store)
	const runID = "review-run-zero"
	if err := rs.EnsureReviewRun(runID, st.ID, time.Now().UTC()); err != nil {
		t.Fatalf("ensure run: %v", err)
	}
	run, now, err := rs.FinalizeReviewRun(st.ID, runID, ReviewFinalization{
		At: time.Now().UTC(), Amount: 0,
		NoteForApplied: func(int) (string, string, string) {
			return "0 XP: short review", "poor_performance", "performance"
		},
	})
	if err != nil || !now {
		t.Fatalf("finalize = %+v, %v, %v", run, now, err)
	}
	if run.XPAdjustment != 0 || run.XPNote != "0 XP: short review" || run.XPReasonCode != "poor_performance" {
		t.Fatalf("zero finalization = %+v", run)
	}
}
