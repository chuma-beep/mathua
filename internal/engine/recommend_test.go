package engine

import (
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/storage"
)

// ReviewsDue counted a row only when the *stored* status was not MASTERED. That is the one
// population the scheduler's own effectiveState explicitly puts in the review bucket:
//
//	if snap.Status == mastery.StatusMastered && snap.NextReviewDue != nil && !due { isReview = true }
//
// So the count and the candidates disagreed, and the consequence was not a wrong number but a
// missing task: scheduler.selectTasks only takes the review bucket when DueReviews > 0, so a
// learner whose sole due review was a decayed-mastered concept was recommended nothing to
// review, while NextReviewQuestion would happily have served one. The comment on the old guard
// described the behaviour the code did not have.
func TestReviewsDue_IncludesDecayedMasteredConcepts(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("dana")
	past := time.Now().Add(-time.Hour)

	rows := []struct {
		name string
		prog storage.ConceptProgress
	}{
		{"decayed mastered, review due", storage.ConceptProgress{
			StudentID: st.ID, ConceptID: "a", Status: "MASTERED",
			NextReviewDue: &past, LastReviewed: &past, LastAttempted: &past,
		}},
		{"mastered, review not yet due", storage.ConceptProgress{
			StudentID: st.ID, ConceptID: "b", Status: "MASTERED",
			NextReviewDue: ptr(time.Now().Add(48 * time.Hour)),
		}},
		{"mastered, no review scheduled", storage.ConceptProgress{
			StudentID: st.ID, ConceptID: "c", Status: "MASTERED",
		}},
		{"learning, review due", storage.ConceptProgress{
			StudentID: st.ID, ConceptID: "d", Status: "LEARNING",
			NextReviewDue: &past, LastAttempted: &past,
		}},
	}
	for _, r := range rows {
		if err := e.repo.UpsertProgress(&r.prog); err != nil {
			t.Fatalf("%s: upsert: %v", r.name, err)
		}
	}

	n, err := e.ReviewsDue(st.ID)
	if err != nil {
		t.Fatalf("reviews due: %v", err)
	}
	// a and d. `b` is scheduled but not yet due; `c` was never scheduled.
	if n != 2 {
		t.Errorf("ReviewsDue = %d, want 2 (decayed-mastered and learning are both due)", n)
	}
}

func ptr(t time.Time) *time.Time { return &t }

// ConceptsMastered counted every non-UNSEEN row, so a concept the learner had merely *started*
// was counted as mastered. Fixed to count effective MASTERED.
//
// This test does NOT discriminate the fix, and the comment is here so nobody later reads it as a
// guard: the only consumer is the new-learner gate `ConceptsMastered == 0 && AttemptedDays == 0`,
// and AttemptedDays already excluded anyone with history, so the miscount never changed a
// response. It is worth stating plainly rather than shipping the fix as though a test proved it.
func TestRecommendNext_StartedConceptsAreNotSentToTheDiagnostic(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("erin")
	now := time.Now()

	// Enough evidence to master `a`, left reviewed just now so decay does not apply.
	if err := e.repo.UpsertProgress(&storage.ConceptProgress{
		StudentID: st.ID, ConceptID: "a", Status: "MASTERED",
		AvgResponseTime: 1, LastReviewed: &now, LastAttempted: &now,
	}); err != nil {
		t.Fatal(err)
	}
	// Started, never mastered.
	if err := e.repo.UpsertProgress(&storage.ConceptProgress{
		StudentID: st.ID, ConceptID: "b", Status: "PRACTICING",
		AvgResponseTime: 1, LastReviewed: &now, LastAttempted: &now,
	}); err != nil {
		t.Fatal(err)
	}

	resp, err := e.RecommendNext(st.ID, nil)
	if err != nil {
		t.Fatalf("recommend: %v", err)
	}
	if resp.Primary == nil {
		t.Fatal("no primary recommendation")
	}
	if resp.Primary.Action.Href == "/onboard" {
		t.Error("a learner with attempts and one mastered concept was sent to the diagnostic")
	}
}
