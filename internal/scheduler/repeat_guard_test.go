package scheduler

import (
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/mastery"
)

// seededSnapshots puts every concept in the practising state so the shelf has candidates. A
// learner with no attempts at all is offered the diagnostic, which is a different decision.
func seededSnapshots(dag *concepts.DAG, now time.Time) map[string]*ConceptSnapshot {
	snaps := make(map[string]*ConceptSnapshot, dag.Count())
	yesterday := now.AddDate(0, 0, -1)
	for id := range dag.Concepts() {
		snaps[id] = &ConceptSnapshot{
			Status:         mastery.StatusLearning,
			RequiredStreak: 3,
			TimeThreshold:  60,
			LastAttempted:  yesterday,
			LastReviewed:   yesterday,
		}
	}
	return snaps
}

// The repeat guard exists because a thin concept asked five times in one sitting returns the
// same situation minutes apart. Measured over a simulated month, one hundred concepts were
// asked more than once in a day and the worst was asked five times from five situations.
func TestRecommend_AvoidsTheConceptJustServed(t *testing.T) {
	dag, err := concepts.LoadDir("../../data/concepts")
	if err != nil {
		t.Fatalf("load concepts: %v", err)
	}
	s := New(dag)
	now := time.Date(2026, 10, 24, 9, 0, 0, 0, time.UTC)
	snaps := seededSnapshots(dag, now)

	ranked := s.conceptRecommendations(RecommendInput{Now: now, Snapshots: snaps}, now, 0)
	if len(ranked) < 2 {
		t.Fatalf("only %d candidates, adjacency cannot be distinguished", len(ranked))
	}
	last := ranked[0].ConceptID

	after := s.conceptRecommendations(RecommendInput{
		Now:           now.Add(time.Minute),
		Snapshots:     snaps,
		LastConceptID: last,
	}, now.Add(time.Minute), 0)
	if len(after) == 0 {
		t.Fatalf("the demotion emptied the ranked list")
	}
	if after[0].ConceptID == last {
		t.Fatalf("ranked %s first again immediately after serving it", last)
	}
}

func TestRecommend_StopsAtTheDailyCap(t *testing.T) {
	dag, err := concepts.LoadDir("../../data/concepts")
	if err != nil {
		t.Fatalf("load concepts: %v", err)
	}
	s := New(dag)
	now := time.Date(2026, 10, 24, 9, 0, 0, 0, time.UTC)
	snaps := seededSnapshots(dag, now)

	// The guard is asserted on the ranked list rather than on the composed shelf, because the
	// shelf also carries a review card and a diagnostic card, which are different decisions.
	unbounded := s.conceptRecommendations(RecommendInput{Now: now, Snapshots: snaps}, now, 0)
	if len(unbounded) < 2 {
		t.Fatalf("only %d candidates, the guard cannot be distinguished", len(unbounded))
	}
	capped := unbounded[0].ConceptID

	guarded := s.conceptRecommendations(RecommendInput{
		Now:          now,
		Snapshots:    snaps,
		ServedToday:  map[string]int{capped: 2},
		ServedPerDay: 2,
	}, now, 0)
	if len(guarded) == 0 {
		t.Fatalf("the cap emptied the ranked list instead of moving to another concept")
	}
	for _, r := range guarded {
		if r.ConceptID == capped {
			t.Fatalf("%s supplied two items today and is still ranked", capped)
		}
	}

	// Raising the cap brings it back, which proves the guard is the cap and not an unrelated
	// change in ranking.
	lifted := s.conceptRecommendations(RecommendInput{
		Now:          now,
		Snapshots:    snaps,
		ServedToday:  map[string]int{capped: 2},
		ServedPerDay: 5,
	}, now, 0)
	found := false
	for _, r := range lifted {
		if r.ConceptID == capped {
			found = true
		}
	}
	if !found {
		t.Fatalf("raising the cap did not restore %s, so the exclusion is not the cap", capped)
	}
}
