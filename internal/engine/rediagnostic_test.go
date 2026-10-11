package engine

import (
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/diagnostic"
	"github.com/chuma-beep/mathua/internal/storage"
)

// A second diagnostic must not cost a learner the progress the first one and
// their practice produced.
//
// The onboarding picker now scopes what is assessed, and a learner can reach the
// picker again at any time, so "take it twice" is a supported path rather than a
// curiosity. `ApplyGoalResults` skips any concept that already has a status
// above UNSEEN, which is the behaviour worth pinning: it is the difference
// between re-placement and a silent wipe.

func diagSession(attempts ...diagnostic.Attempt) *diagnostic.Session {
	return &diagnostic.Session{
		ID:       "sess-test",
		Attempts: attempts,
		State:    diagnostic.StateDone,
	}
}

func TestApplyGoalResults_PreservesExistingMastery(t *testing.T) {
	e := testEngine(t)
	st, err := e.CreateStudent("rediag")
	if err != nil {
		t.Fatal(err)
	}

	// The learner has genuinely practised this concept: status, streak and SM-2
	// state that a diagnostic has no business overwriting.
	if err := e.repo.UpsertProgress(&storage.ConceptProgress{
		StudentID:     st.ID,
		ConceptID:     "arith.add.single",
		Status:        "PRACTICING",
		Streak:        4,
		BestStreak:    4,
		Attempts:      9,
		SM2EFactor:    2.6,
		WeaknessScore: 0.2,
	}); err != nil {
		t.Fatal(err)
	}

	err = e.ApplyGoalResults(st.ID, diagSession(diagnostic.Attempt{
		ConceptID: "arith.add.single",
		Correct:   false,
		Timestamp: time.Now().UTC(),
	}))
	if err != nil {
		t.Fatalf("apply: %v", err)
	}

	got, err := e.repo.GetProgress(st.ID, "arith.add.single")
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != "PRACTICING" {
		t.Errorf("status clobbered by a diagnostic miss: got %q, want PRACTICING", got.Status)
	}
	if got.Streak != 4 || got.Attempts != 9 || got.BestStreak != 4 {
		t.Errorf("practice history lost: streak=%d attempts=%d best=%d", got.Streak, got.Attempts, got.BestStreak)
	}
	if got.WeaknessScore != 0.2 {
		t.Errorf("weakness rewritten: got %v, want 0.2", got.WeaknessScore)
	}
}

func TestApplyGoalResults_SeedsOnlyUnseenConcepts(t *testing.T) {
	e := testEngine(t)
	st, err := e.CreateStudent("rediag_seed")
	if err != nil {
		t.Fatal(err)
	}
	if err := e.repo.UpsertProgress(&storage.ConceptProgress{
		StudentID: st.ID, ConceptID: "arith.add.single", Status: "MASTERED", Streak: 12,
	}); err != nil {
		t.Fatal(err)
	}

	now := time.Now().UTC()
	err = e.ApplyGoalResults(st.ID, diagSession(
		diagnostic.Attempt{ConceptID: "arith.add.single", Correct: false, Timestamp: now},
		diagnostic.Attempt{ConceptID: "arith.neg.number_line", Correct: true, Fast: true, Timestamp: now},
	))
	if err != nil {
		t.Fatalf("apply: %v", err)
	}

	mastered, err := e.repo.GetProgress(st.ID, "arith.add.single")
	if err != nil {
		t.Fatal(err)
	}
	if mastered.Status != "MASTERED" {
		t.Errorf("a mastered concept must survive a rediagnostic, got %q", mastered.Status)
	}
	if mastered.PlacementSeeded {
		t.Error("a concept the diagnostic declined to seed must not be marked placement-seeded")
	}

	// The genuinely new concept is seeded, and seeded as reachability rather than
	// evidence: never MASTERED, however well it was answered.
	fresh, err := e.repo.GetProgress(st.ID, "arith.neg.number_line")
	if err != nil {
		t.Fatal(err)
	}
	if fresh == nil {
		t.Fatal("a new concept was not seeded at all")
	}
	if fresh.Status == "MASTERED" {
		t.Error("the diagnostic granted mastery; it must not")
	}
	if !fresh.PlacementSeeded {
		t.Error("a probed concept should be placement-seeded so its successors unlock")
	}
}

func TestApplyGoalResults_LeavesOtherDomainsAlone(t *testing.T) {
	e := testEngine(t)
	st, err := e.CreateStudent("rediag_scope")
	if err != nil {
		t.Fatal(err)
	}
	// Progress from a domain this diagnostic never touches.
	if err := e.repo.UpsertProgress(&storage.ConceptProgress{
		StudentID: st.ID, ConceptID: "geo.basic.points_lines", Status: "PRACTICING", Streak: 3,
	}); err != nil {
		t.Fatal(err)
	}

	if err := e.ApplyGoalResults(st.ID, diagSession(diagnostic.Attempt{
		ConceptID: "arith.add.single", Correct: true, Timestamp: time.Now().UTC(),
	})); err != nil {
		t.Fatal(err)
	}

	other, err := e.repo.GetProgress(st.ID, "geo.basic.points_lines")
	if err != nil {
		t.Fatal(err)
	}
	if other == nil || other.Status != "PRACTICING" || other.Streak != 3 {
		t.Errorf("an unassessed domain was disturbed: %+v", other)
	}
}
