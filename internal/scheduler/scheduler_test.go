package scheduler

import (
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/mastery"
)

// SM2 tests
func TestSM2_Default(t *testing.T) {
	sm := DefaultSM2()
	if sm.EFactor != 2.5 || sm.Repetitions != 0 || sm.Interval != 0 {
		t.Errorf("unexpected default SM2: %+v", sm)
	}
}

func TestSM2_PerfectReview(t *testing.T) {
	prev := SM2{Repetitions: 0, EFactor: 2.5, Interval: 0}
	next := ComputeSM2(prev, 5)
	if next.Repetitions != 1 || next.Interval != 1 {
		t.Errorf("expected rep=1 interval=1, got %+v", next)
	}
}

func TestSM2_SecondReview(t *testing.T) {
	prev := SM2{Repetitions: 1, EFactor: 2.5, Interval: 1}
	next := ComputeSM2(prev, 5)
	if next.Interval != 6 {
		t.Errorf("expected interval=6, got %d", next.Interval)
	}
}

func TestSM2_ThirdReview(t *testing.T) {
	prev := SM2{Repetitions: 2, EFactor: 2.5, Interval: 6}
	next := ComputeSM2(prev, 4)
	if next.Repetitions != 3 || next.Interval != 15 {
		t.Errorf("expected rep=3 interval=15, got %+v", next)
	}
}

func TestSM2_FailedReview(t *testing.T) {
	prev := SM2{Repetitions: 3, EFactor: 2.5, Interval: 15}
	next := ComputeSM2(prev, 2)
	if next.Repetitions != 0 || next.Interval != 1 {
		t.Errorf("expected rep=0 interval=1, got %+v", next)
	}
}

func TestSM2_EFactorFloor(t *testing.T) {
	prev := SM2{Repetitions: 1, EFactor: 1.31, Interval: 6}
	for i := 0; i < 10; i++ {
		prev = ComputeSM2(prev, 0)
	}
	if prev.EFactor < 1.3 {
		t.Errorf("EFactor below 1.3: %f", prev.EFactor)
	}
}

// Scheduler tests

func miniDAG(t *testing.T) *concepts.DAG {
	t.Helper()
	d, err := concepts.Build([]concepts.Concept{
		{ID: "a", Domain: "d", Prerequisites: []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 10}},
		{ID: "b", Domain: "d", Prerequisites: []string{"a"},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 10}},
		{ID: "c", Domain: "d", Prerequisites: []string{"a"},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 10}},
		{ID: "d", Domain: "d", Prerequisites: []string{"b", "c"},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 10}},
	})
	if err != nil {
		t.Fatalf("build mini DAG: %v", err)
	}
	return d
}

func TestNext_BootstrapSingleConcept(t *testing.T) {
	d := miniDAG(t)
	s := New(d)
	snap := map[string]*ConceptSnapshot{}
	next := s.Next(snap, "", 0, 0)
	if next == nil {
		t.Fatal("expected a next concept")
	}
	if next.Concept.ID != "a" {
		t.Errorf("expected a, got %s", next.Concept.ID)
	}
	if next.IsReview {
		t.Error("expected not review")
	}
}

func TestNext_OnlyAvailableConcepts(t *testing.T) {
	d := miniDAG(t)
	s := New(d)
	now := time.Now().UTC()
	snap := map[string]*ConceptSnapshot{
		"a": {Status: mastery.StatusMastered, LastReviewed: now, RequiredStreak: 3},
	}
	next := s.Next(snap, "", 0, 0)
	if next == nil {
		t.Fatal("expected a next concept")
	}
	id := next.Concept.ID
	if id != "b" && id != "c" {
		t.Errorf("expected b or c, got %s", id)
	}
}

func TestNext_SkipsPrevConcept(t *testing.T) {
	d := miniDAG(t)
	s := New(d)
	now := time.Now().UTC()
	snap := map[string]*ConceptSnapshot{
		"a": {Status: mastery.StatusMastered, LastReviewed: now, RequiredStreak: 3},
	}
	next := s.Next(snap, "b", 0, 0)
	if next == nil {
		t.Fatal("expected a next concept")
	}
	if next.Concept.ID == "b" {
		t.Error("expected to skip b as prevConceptID")
	}
}

func TestNext_NothingAvailable(t *testing.T) {
	d := miniDAG(t)
	s := New(d)
	now := time.Now().UTC()
	snap := map[string]*ConceptSnapshot{
		"a": {Status: mastery.StatusMastered, LastReviewed: now, RequiredStreak: 3},
		"b": {Status: mastery.StatusMastered, LastReviewed: now, RequiredStreak: 3},
		"c": {Status: mastery.StatusMastered, LastReviewed: now, RequiredStreak: 3},
		"d": {Status: mastery.StatusMastered, LastReviewed: now, RequiredStreak: 3},
	}
	next := s.Next(snap, "", 0, 0)
	if next != nil {
		t.Errorf("expected nil, got %s", next.Concept.ID)
	}
}

func TestNext_DecayingPriority(t *testing.T) {
	d := miniDAG(t)
	s := New(d)
	old := time.Now().UTC().AddDate(0, 0, -20)
	now := time.Now().UTC()
	snap := map[string]*ConceptSnapshot{
		"a": {Status: mastery.StatusMastered, LastReviewed: now, RequiredStreak: 3},
		"b": {Status: mastery.StatusMastered, LastReviewed: old, RequiredStreak: 3},
	}
	next := s.Next(snap, "", 0, 0)
	if next == nil {
		t.Fatal("expected a next concept")
	}
	if !next.IsReview && next.Concept.ID == "b" {
		t.Error("decaying concept b should be marked as review")
	}
}

func TestNext_70_30_BalancePrefersReview(t *testing.T) {
	d := miniDAG(t)
	s := New(d)
	old := time.Now().UTC().AddDate(0, 0, -20)
	now := time.Now().UTC()
	snap := map[string]*ConceptSnapshot{
		"a": {Status: mastery.StatusMastered, LastReviewed: now, RequiredStreak: 3},
		"b": {Status: mastery.StatusMastered, LastReviewed: old, RequiredStreak: 3},
	}
	next := s.Next(snap, "", 0, 3)
	if next == nil {
		t.Fatal("expected a next concept")
	}
	if !next.IsReview {
		t.Error("expected review after 3 new questions")
	}
}

func TestNext_FallthroughWhenBlocked(t *testing.T) {
	d := miniDAG(t)
	s := New(d)
	snap := map[string]*ConceptSnapshot{}
	next := s.Next(snap, "a", 0, 0)
	if next == nil {
		t.Fatal("expected fallthrough when only concept blocked")
	}
	if next.Concept.ID != "a" {
		t.Errorf("expected fallthrough to a, got %s", next.Concept.ID)
	}
}
