package storage

import (
	"fmt"
	"testing"
	"time"
)

// TestPostgresReviewRunFinalization runs the review-run transaction against
// the real PostgreSQL backend. It is skipped unless TEST_POSTGRES_DSN is set.
func TestPostgresReviewRunFinalization(t *testing.T) {
	pg := newPostgresStore(t)
	es := EffortStore(pg)
	rrs := ReviewRunStore(pg)
	st := seedStudent(t, pg)
	runID := fmt.Sprintf("review-pg-%d", time.Now().UnixNano())
	if err := rrs.EnsureReviewRun(runID, st.ID, time.Now().UTC()); err != nil {
		t.Fatalf("ensure run: %v", err)
	}
	if err := pg.AwardXP(st.ID, 10, "seed", "seed:"+runID); err != nil {
		t.Fatalf("seed: %v", err)
	}
	for i := 0; i < 6; i++ {
		if err := es.RecordEffortSignal(EffortSignal{
			StudentID: st.ID, TaskKind: "review", TaskID: runID + ":concept",
			ReviewRunID: runID, ConceptID: "concept", AttemptKey: fmt.Sprintf("attempt-%d", i),
			Correct: false, ElapsedSeconds: 10, Difficulty: 0.3,
			AnswerSignature: fmt.Sprintf("wrong-%d", i), CreatedAt: time.Now().UTC(),
		}); err != nil {
			t.Fatalf("record signal %d: %v", i, err)
		}
	}
	fin := ReviewFinalization{
		At: time.Now().UTC(), Amount: -2, Reason: "accountability:poor_performance",
		EventKey: "accountability:review-run:" + runID,
		NoteForApplied: func(applied int) (string, string, string) {
			return fmt.Sprintf("%d XP: review performance", applied), "poor_performance", "performance"
		},
	}
	run, finalized, err := rrs.FinalizeReviewRun(st.ID, runID, fin)
	if err != nil || !finalized {
		t.Fatalf("finalize = %+v, %v, %v; want finalized", run, finalized, err)
	}
	if run.XPAdjustment != -2 || run.XPNote == "" || run.XPSource != "performance" {
		t.Fatalf("stored run = %+v", run)
	}
	total, _, err := pg.GetXP(st.ID)
	if err != nil || total != 8 {
		t.Fatalf("balance = %d, %v; want 8", total, err)
	}
	again, finalized, err := rrs.FinalizeReviewRun(st.ID, runID, ReviewFinalization{
		At: time.Now().UTC(), Amount: -9, Reason: "accountability:changed",
		EventKey: "accountability:review-run:" + runID,
	})
	if err != nil || finalized {
		t.Fatalf("second finalize = %+v, %v, %v; want stored result", again, finalized, err)
	}
	if again.XPAdjustment != -2 || again.XPNote == "" {
		t.Fatalf("second finalize changed stored outcome: %+v", again)
	}
	total, _, err = pg.GetXP(st.ID)
	if err != nil || total != 8 {
		t.Fatalf("balance after replay = %d, %v; want 8", total, err)
	}
}
