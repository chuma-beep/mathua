package engine

import (
	"time"

	"github.com/chuma-beep/mathua/internal/storage"
	xppolicy "github.com/chuma-beep/mathua/internal/xp"
)

// ReviewPerformanceResult is the durable outcome of finalizing one review run.
// Zero values mean the run completed without an adjustment.
type ReviewPerformanceResult struct {
	XPAdjustment int
	XPNote       string
	XPReasonCode string
	XPSource     string
}

// EnsureReviewRun creates the durable review-run row when a review session
// starts. Stores without review-run persistence keep review effort-only.
func (e *Engine) EnsureReviewRun(studentID, reviewRunID string) error {
	store, ok := e.repo.(storage.ReviewRunStore)
	if !ok {
		return nil
	}
	return store.EnsureReviewRun(reviewRunID, studentID, time.Now().UTC())
}

// FinalizeReviewPerformance is called only when NextReviewQuestion has
// authoritatively reported that no review questions remain. It groups the
// review run's signals, evaluates performance once, and applies at most one
// task-level adjustment. A run that already received an effort deduction is
// finalized without a performance adjustment to avoid double charging.
func (e *Engine) FinalizeReviewPerformance(studentID, reviewRunID string) (ReviewPerformanceResult, error) {
	store, ok := e.repo.(storage.ReviewRunStore)
	if !ok {
		return ReviewPerformanceResult{}, nil
	}
	if err := store.EnsureReviewRun(reviewRunID, studentID, time.Now().UTC()); err != nil {
		return ReviewPerformanceResult{}, err
	}
	signals, err := store.ReviewRunSignals(studentID, reviewRunID)
	if err != nil {
		return ReviewPerformanceResult{}, err
	}

	fin := storage.ReviewFinalization{
		At:       time.Now().UTC(),
		EventKey: "accountability:review-run:" + reviewRunID,
	}
	effortDeductions, err := store.CountReviewRunEffortDeductions(studentID, reviewRunID)
	if err != nil {
		return ReviewPerformanceResult{}, err
	}

	var decision xppolicy.Decision
	if effortDeductions == 0 {
		attempts := make([]xppolicy.AttemptSignal, 0, len(signals))
		for _, sig := range signals {
			attempts = append(attempts, effortAttemptSignal(sig))
		}
		available, _, xerr := e.repo.GetXP(studentID)
		if xerr != nil {
			return ReviewPerformanceResult{}, xerr
		}
		decision = xppolicy.EvaluatePerformance(xppolicy.Input{
			TaskKind:     xppolicy.TaskReview,
			TaskID:       reviewRunID,
			Attempts:     attempts,
			AvailableXP:  available,
			TaskComplete: true,
			Config:       e.effortConfig,
		})
		fin.Amount = decision.Amount
		fin.Reason = decision.LedgerReason()
		fin.NoteForApplied = func(applied int) (string, string, string) {
			_, note, reason := noteForApplied(decision, applied)
			return note, reason, string(decision.Source)
		}
	}

	run, finalizedNow, err := store.FinalizeReviewRun(studentID, reviewRunID, fin)
	if err != nil {
		return ReviewPerformanceResult{}, err
	}
	if finalizedNow && effortDeductions == 0 {
		e.accountabilityStats.recordDecision(decision)
		e.accountabilityStats.recordApplied(decision, run.XPAdjustment)
	}
	return ReviewPerformanceResult{
		XPAdjustment: run.XPAdjustment,
		XPNote:       run.XPNote,
		XPReasonCode: run.XPReasonCode,
		XPSource:     run.XPSource,
	}, nil
}
