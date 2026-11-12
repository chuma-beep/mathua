package engine

import (
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/diagnostic"
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/scheduler"
	"github.com/chuma-beep/mathua/internal/storage"
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

// A MASTERED row with no review timestamp (legacy or freshly seeded) reads
// as DECAYING, never as fresh mastery: GetProgress applies decay on read,
// and a missing date is the stalest possible signal.
func TestGetProgress_DecaysFromMasteryWhenNoReviewClock(t *testing.T) {
	// A concept mastered through study has no review timestamp until its first
	// review. The decay clock must start at mastery: treating the missing clock
	// as 999 days called every freshly mastered concept due, while the
	// scheduler treated the same row as fresh and never offered it.
	e := testEngine(t)
	st, err := e.CreateStudent("legacy-mastered")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}
	fresh := time.Now().UTC().Add(-time.Hour)
	stale := time.Now().UTC().AddDate(0, 0, -30)
	for id, masteredAt := range map[string]time.Time{"a": fresh, "b": stale} {
		if err := e.repo.UpsertProgress(&storage.ConceptProgress{
			StudentID: st.ID, ConceptID: id, Status: "MASTERED", Streak: 5,
			MasteredAt: &masteredAt,
		}); err != nil {
			t.Fatalf("seed mastery %s: %v", id, err)
		}
	}
	got, err := e.GetProgress(st.ID)
	if err != nil {
		t.Fatalf("get progress: %v", err)
	}
	if got["a"].Status != "MASTERED" {
		t.Errorf("fresh mastery status = %q, want MASTERED", got["a"].Status)
	}
	if got["b"].Status != "DECAYING" {
		t.Errorf("stale mastery status = %q, want DECAYING", got["b"].Status)
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

// Conditionally completed concepts from the diagnostic report seed reachability
// as PRACTICING, unlocking their successors at the frontier. A struggle on the
// successor demotes the provisional prerequisite to LEARNING, which relocks
// dependent progress until remedial practice certifies it.
func TestApplyDiagnosticReport_SuccessorIsEligibleUntilFallbackDemotion(t *testing.T) {
	e := testEngine(t)
	st, err := e.CreateStudent("placed-diagnostic-learner")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}

	report := &diagnostic.DiagnosticReport{
		ConditionallyCompleted: []string{"Concept A"},
	}
	if err := e.ApplyDiagnosticReport(st.ID, report); err != nil {
		t.Fatalf("apply diagnostic report: %v", err)
	}

	// Prerequisite a is seeded PRACTICING with placement and provisional flags.
	pa, err := e.repo.GetProgress(st.ID, "a")
	if err != nil || pa == nil {
		t.Fatalf("get progress a: %v", err)
	}
	if pa.Status != "PRACTICING" || !pa.PlacementSeeded || !pa.IsProvisional {
		t.Fatalf("a progress = %+v, want PRACTICING with placement and provisional flags", pa)
	}

	// Successor b is unlocked and startable.
	child, err := e.EligibilityFor(st.ID, "b")
	if err != nil {
		t.Fatalf("eligibility b: %v", err)
	}
	if !child.Eligible || child.State != scheduler.StateUnlocked {
		t.Fatalf("b = %+v, want eligible unlocked via placement", child)
	}

	// Learner struggles on b: provisional fallback demotes prerequisite a.
	demoted, err := e.CheckProvisionalFallback(st.ID, "b")
	if err != nil {
		t.Fatalf("check fallback: %v", err)
	}
	if len(demoted) == 0 || demoted[0] != "a" {
		t.Fatalf("expected demoted a, got %v", demoted)
	}

	// Demoted prerequisite a is now LEARNING with elevated weakness.
	paAfter, err := e.repo.GetProgress(st.ID, "a")
	if err != nil || paAfter == nil {
		t.Fatalf("get progress a after fallback: %v", err)
	}
	if paAfter.Status != "LEARNING" || paAfter.WeaknessScore != 0.8 {
		t.Fatalf("a progress after demotion = %+v, want LEARNING and weakness 0.8", paAfter)
	}

	// Successor b is now relocked until a is remediated.
	childAfter, err := e.EligibilityFor(st.ID, "b")
	if err != nil {
		t.Fatalf("eligibility b after fallback: %v", err)
	}
	if childAfter.Eligible || childAfter.State == scheduler.StateUnlocked {
		t.Fatalf("b after fallback = %+v, want locked due to demoted prerequisite", childAfter)
	}
}

// chainEngine is a three-level chain, a -> b -> c, which is the shape the
// closure seeding has to cross: probing c leaves two foundations unprobed.
func chainEngine(t *testing.T) *Engine {
	t.Helper()
	d, err := concepts.Build([]concepts.Concept{
		{ID: "a", Label: "Concept A", Domain: "d", GradingType: "numeric", Prerequisites: []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60}},
		{ID: "b", Label: "Concept B", Domain: "d", GradingType: "numeric", Prerequisites: []string{"a"},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60}},
		{ID: "c", Label: "Concept C", Domain: "d", GradingType: "numeric", Prerequisites: []string{"b"},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60}},
	})
	if err != nil {
		t.Fatalf("build chain DAG: %v", err)
	}
	store, err := storage.NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("create store: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	reg := generator.NewRegistry()
	for _, id := range []string{"a", "b", "c"} {
		reg.Register(id, &testGen{answer: "42"})
	}
	return New(store, d, reg, nil, nil)
}

// The diagnostic samples a large prerequisite closure, so the foundations below
// a probed concept are usually never asked. They must still be reachable
// afterwards, or the learner sees an unlocked concept whose definitions they
// cannot study.
func TestApplyGoalResults_SeedsPlacementForUnprobedPrerequisites(t *testing.T) {
	e := chainEngine(t)
	st, err := e.CreateStudent("unprobed-foundations")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}

	session := &diagnostic.Session{
		StudentID: st.ID,
		Attempts:  []diagnostic.Attempt{{ConceptID: "c", Correct: true, Fast: true}},
	}
	if err := e.ApplyGoalResults(st.ID, session); err != nil {
		t.Fatalf("apply placement: %v", err)
	}

	progress, err := e.repo.GetAllProgress(st.ID)
	if err != nil {
		t.Fatalf("get progress: %v", err)
	}
	for _, id := range []string{"a", "b"} {
		p := progress[id]
		if p == nil {
			t.Fatalf("%s has no row, so it stays locked", id)
		}
		if !p.PlacementSeeded {
			t.Errorf("%s PlacementSeeded = false, want true", id)
		}
		if p.Status != "UNSEEN" {
			t.Errorf("%s status = %q, want UNSEEN (placement is not evidence)", id, p.Status)
		}
		if p.LastAttempted != nil {
			t.Errorf("%s LastAttempted is set, but the diagnostic never asked it", id)
		}
	}
	if progress["c"].Status != "PRACTICING" {
		t.Errorf("c status = %q, want PRACTICING", progress["c"].Status)
	}

	for _, id := range []string{"a", "b"} {
		el, err := e.EligibilityFor(st.ID, id)
		if err != nil {
			t.Fatalf("eligibility %s: %v", id, err)
		}
		if !el.Eligible || el.State != scheduler.StateUnlocked {
			t.Errorf("%s = %+v, want eligible unlocked", id, el)
		}
	}
	if el, err := e.EligibilityFor(st.ID, "c"); err != nil || !el.Eligible {
		t.Errorf("c = %+v, %v, want eligible after being probed", el, err)
	}
	for id, p := range progress {
		if p.Status == "MASTERED" {
			t.Errorf("%s was seeded as MASTERED; placement must never write mastery", id)
		}
	}
}

// Seeding must not touch a foundation the learner has already practised, even
// when it sits below a probed concept.
func TestApplyGoalResults_DoesNotReseedPracticedPrerequisite(t *testing.T) {
	e := chainEngine(t)
	st, err := e.CreateStudent("practised-foundation")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}
	if err := e.repo.UpsertProgress(&storage.ConceptProgress{
		StudentID:     st.ID,
		ConceptID:     "a",
		Status:        "PRACTICING",
		Streak:        4,
		BestStreak:    4,
		Attempts:      9,
		SM2EFactor:    2.6,
		WeaknessScore: 0.2,
	}); err != nil {
		t.Fatalf("seed practised row: %v", err)
	}

	session := &diagnostic.Session{
		StudentID: st.ID,
		Attempts:  []diagnostic.Attempt{{ConceptID: "c", Correct: true, Fast: true}},
	}
	if err := e.ApplyGoalResults(st.ID, session); err != nil {
		t.Fatalf("apply placement: %v", err)
	}

	a, err := e.repo.GetProgress(st.ID, "a")
	if err != nil {
		t.Fatalf("get a: %v", err)
	}
	if a.Status != "PRACTICING" || a.Streak != 4 || a.Attempts != 9 || a.PlacementSeeded {
		t.Errorf("practised prerequisite was rewritten: %+v", a)
	}
}
