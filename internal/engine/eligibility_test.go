package engine

import (
	"testing"

	"github.com/chuma-beep/mathua/internal/diagnostic"
	"github.com/chuma-beep/mathua/internal/scheduler"
)

func TestEligibilityFor_FreshLearner(t *testing.T) {
	e := testEngine(t)
	st, err := e.CreateStudent("placement-learner")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}

	root, err := e.EligibilityFor(st.ID, "a")
	if err != nil {
		t.Fatalf("eligibility a: %v", err)
	}
	if !root.Eligible || root.State != scheduler.StateUnlocked {
		t.Errorf("a = %+v, want eligible unlocked (no prerequisites)", root)
	}

	child, err := e.EligibilityFor(st.ID, "b")
	if err != nil {
		t.Fatalf("eligibility b: %v", err)
	}
	if child.Eligible || child.State != scheduler.StateLocked {
		t.Errorf("b = %+v, want locked (a unseen)", child)
	}
	if child.Reason != "missing_prerequisites" {
		t.Errorf("b.reason = %q, want missing_prerequisites", child.Reason)
	}
	if len(child.Prerequisites) != 1 || child.Prerequisites[0].Met {
		t.Errorf("b.prerequisites = %+v, want a marked unmet", child.Prerequisites)
	}
}

// Placement seeds reachability, never mastery: a diagnostic-correct concept
// unlocks its successor for starting while still reporting zero mastered
// concepts. This is the property the whole post-diagnostic flow rests on.
func TestApplyGoalResults_PlacementUnlocksWithoutMastery(t *testing.T) {
	e := testEngine(t)
	st, err := e.CreateStudent("placed-learner")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}
	session := &diagnostic.Session{
		StudentID: st.ID,
		Attempts:  []diagnostic.Attempt{{ConceptID: "a", Correct: true, Fast: true}},
	}
	if err := e.ApplyGoalResults(st.ID, session); err != nil {
		t.Fatalf("apply placement: %v", err)
	}

	// The seeded row is placement, not mastery.
	progress, err := e.repo.GetAllProgress(st.ID)
	if err != nil {
		t.Fatalf("get progress: %v", err)
	}
	ap := progress["a"]
	if ap == nil || !ap.PlacementSeeded {
		t.Fatalf("a progress = %+v, want PlacementSeeded", ap)
	}
	if ap.Status == "MASTERED" {
		t.Errorf("placement wrote mastery: status = %q", ap.Status)
	}

	// The successor is now startable.
	child, err := e.EligibilityFor(st.ID, "b")
	if err != nil {
		t.Fatalf("eligibility b: %v", err)
	}
	if !child.Eligible || child.State != scheduler.StateUnlocked {
		t.Errorf("b = %+v, want eligible unlocked via placement", child)
	}
}
