package engine

import (
	"testing"
	"time"
)

// A notification scan may ask how many reviews are due, but it must not use
// the question-serving path that satisfies implicit reviews. This test pins the
// observable contract: the count is available and every stored scheduling field
// is byte-for-byte unchanged.
func TestReviewsDueReadOnlyDoesNotMutateReviewState(t *testing.T) {
	e := testEngine(t)
	st, err := e.CreateStudent("readonly-review-count")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}
	now := time.Now().UTC()
	p := masteredReviewRow(st.ID, now.Add(-time.Hour))
	if err := e.repo.UpsertProgress(p); err != nil {
		t.Fatalf("seed progress: %v", err)
	}

	before, err := e.repo.GetProgress(st.ID, "a")
	if err != nil || before == nil {
		t.Fatalf("read before: %v, %v", before, err)
	}
	attemptsBefore, err := e.repo.GetAttemptsForStudent(st.ID)
	if err != nil {
		t.Fatalf("attempts before: %v", err)
	}

	due, err := e.ReviewsDueReadOnly(st.ID)
	if err != nil {
		t.Fatalf("read-only due count: %v", err)
	}
	if due != 1 {
		t.Fatalf("due count = %d, want 1", due)
	}

	after, err := e.repo.GetProgress(st.ID, "a")
	if err != nil || after == nil {
		t.Fatalf("read after: %v, %v", after, err)
	}
	if after.NextReviewDue == nil || !after.NextReviewDue.Equal(*before.NextReviewDue) {
		t.Errorf("NextReviewDue changed: %v -> %v", before.NextReviewDue, after.NextReviewDue)
	}
	if after.LastReviewed == nil || !after.LastReviewed.Equal(*before.LastReviewed) {
		t.Errorf("LastReviewed changed: %v -> %v", before.LastReviewed, after.LastReviewed)
	}
	if after.FSRSStability == nil || *after.FSRSStability != *before.FSRSStability {
		t.Errorf("FSRSStability changed: %v -> %v", before.FSRSStability, after.FSRSStability)
	}
	if after.FSRSDifficulty == nil || *after.FSRSDifficulty != *before.FSRSDifficulty {
		t.Errorf("FSRSDifficulty changed: %v -> %v", before.FSRSDifficulty, after.FSRSDifficulty)
	}
	if after.FSRSLastReview == nil || !after.FSRSLastReview.Equal(*before.FSRSLastReview) {
		t.Errorf("FSRSLastReview changed: %v -> %v", before.FSRSLastReview, after.FSRSLastReview)
	}
	if after.SM2Interval != before.SM2Interval || after.SM2Repetitions != before.SM2Repetitions || after.SM2EFactor != before.SM2EFactor {
		t.Errorf("SM-2 state changed: before %+v after %+v", before, after)
	}

	attemptsAfter, err := e.repo.GetAttemptsForStudent(st.ID)
	if err != nil {
		t.Fatalf("attempts after: %v", err)
	}
	if len(attemptsAfter) != len(attemptsBefore) {
		t.Errorf("read-only count wrote attempts: before %d, after %d", len(attemptsBefore), len(attemptsAfter))
	}
}
