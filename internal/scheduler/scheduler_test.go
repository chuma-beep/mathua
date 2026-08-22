package scheduler

import (
	"math"
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

func TestNext_FiltersMasteredConcepts(t *testing.T) {
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
		t.Errorf("expected nil (all mastered), got concept %s", next.Concept.ID)
	}
}

func TestNext_AllowsDecayingConcept(t *testing.T) {
	d := miniDAG(t)
	s := New(d)
	old := time.Now().UTC().AddDate(0, 0, -20)
	now := time.Now().UTC()
	snap := map[string]*ConceptSnapshot{
		"a": {Status: mastery.StatusMastered, LastReviewed: old, RequiredStreak: 3},
		"b": {Status: mastery.StatusMastered, LastReviewed: now, RequiredStreak: 3},
	}
	next := s.Next(snap, "", 0, 0)
	if next == nil {
		t.Fatal("expected decaying concept to be selected")
	}
	if !next.IsReview {
		t.Error("decaying should be review")
	}
}

func TestNext_PrereqGating(t *testing.T) {
	d := miniDAG(t)
	s := New(d)
	now := time.Now().UTC()
	snap := map[string]*ConceptSnapshot{
		"a": {Status: mastery.StatusLearning, Streak: 1, LastAttempted: now, RequiredStreak: 3},
	}
	for i := 0; i < 20; i++ {
		next := s.Next(snap, "", 0, 0)
		if next == nil {
			t.Fatal("expected root concept while prereqs unmastered")
		}
		if next.Concept.ID != "a" {
			t.Fatalf("prereq gating violated: returned %s though a is unmastered", next.Concept.ID)
		}
	}
}

func TestComputePriority_UnseenBaseline(t *testing.T) {
	if got := computePriority(nil, time.Now().UTC(), false); math.Abs(got-0.2) > 1e-9 {
		t.Errorf("expected 0.2 baseline for unseen, got %f", got)
	}
}

func TestComputePriority_Weights(t *testing.T) {
	now := time.Now().UTC()
	snap := &ConceptSnapshot{
		LastAttempted:  now.Add(-48 * time.Hour),
		Streak:         1,
		RequiredStreak: 4,
		WeaknessScore:  0.5,
	}
	want := 0.5*2.0 + 0.2*(1.0-0.25) + 0.3*0.5
	if got := computePriority(snap, now, false); math.Abs(got-want) > 1e-9 {
		t.Errorf("expected %f, got %f", want, got)
	}
}

func TestComputePriority_DecayingBoost(t *testing.T) {
	now := time.Now().UTC()
	snap := &ConceptSnapshot{LastAttempted: now.Add(-24 * time.Hour), Streak: 2, RequiredStreak: 4}
	plain := computePriority(snap, now, false)
	decay := computePriority(snap, now, true)
	if decay-plain != 5.0 {
		t.Errorf("expected +5.0 decay boost, got %f", decay-plain)
	}
}

func TestEffectiveState_NextReviewDuePassed(t *testing.T) {
	now := time.Now().UTC()
	past := now.Add(-time.Hour)
	future := now.Add(24 * time.Hour)

	status, isReview := effectiveState(&ConceptSnapshot{
		Status:        mastery.StatusMastered,
		LastReviewed:  now.Add(-24 * time.Hour),
		NextReviewDue: &past,
	}, now)
	if status != mastery.StatusMastered {
		t.Errorf("freshly-mastered should not be DECAYING, got %s", status)
	}
	if !isReview {
		t.Error("due review date in the past should mark IsReview")
	}

	_, isReview = effectiveState(&ConceptSnapshot{
		Status:        mastery.StatusMastered,
		LastReviewed:  now,
		NextReviewDue: &future,
	}, now)
	if isReview {
		t.Error("future review due date must not mark IsReview")
	}
}

func TestSelectWithBalance_ForcesReviewWhenUnderRatio(t *testing.T) {
	newTop := Candidate{Concept: &concepts.Concept{ID: "n"}, Priority: 10}
	reviewLow := Candidate{Concept: &concepts.Concept{ID: "r"}, IsReview: true, Priority: 1}

	got := selectWithBalance([]Candidate{newTop, reviewLow}, 0, 9)
	if !got.IsReview {
		t.Error("under 30%% review ratio, a review candidate must be preferred over higher-priority new")
	}
}

func TestSelectWithBalance_AtRatioPicksTopPriority(t *testing.T) {
	newTop := Candidate{Concept: &concepts.Concept{ID: "n"}, Priority: 10}
	reviewLow := Candidate{Concept: &concepts.Concept{ID: "r"}, IsReview: true, Priority: 1}

	got := selectWithBalance([]Candidate{newTop, reviewLow}, 3, 6)
	if got.IsReview || got.Concept.ID != "n" {
		t.Errorf("at target ratio the top-priority candidate should win, got %+v", got.Concept.ID)
	}
}

func TestSelectWithBalance_NoReviewsFallsBackToTop(t *testing.T) {
	a := Candidate{Concept: &concepts.Concept{ID: "a"}, Priority: 10}
	b := Candidate{Concept: &concepts.Concept{ID: "b"}, Priority: 5}
	got := selectWithBalance([]Candidate{a, b}, 0, 99)
	if got.Concept.ID != "a" {
		t.Errorf("without review candidates the top candidate should win, got %s", got.Concept.ID)
	}
}

func TestNextReview_ReturnsOnlyReviewCandidates(t *testing.T) {
	d := miniDAG(t)
	s := New(d)
	old := time.Now().UTC().AddDate(0, 0, -20)
	snap := map[string]*ConceptSnapshot{
		"a": {Status: mastery.StatusMastered, LastReviewed: old, RequiredStreak: 3},
		"b": {Status: mastery.StatusMastered, LastReviewed: old, RequiredStreak: 3},
	}
	for i := 0; i < 10; i++ {
		nx := s.NextReview(snap, "")
		if nx == nil {
			t.Fatal("expected a review candidate")
		}
		if !nx.IsReview {
			t.Fatalf("NextReview returned non-review %s", nx.Concept.ID)
		}
		if nx.Concept.ID != "a" && nx.Concept.ID != "b" {
			t.Fatalf("unexpected review pick %s; only a/b are decaying", nx.Concept.ID)
		}
	}
}

func TestNextReview_SkipsPrev(t *testing.T) {
	d := miniDAG(t)
	s := New(d)
	old := time.Now().UTC().AddDate(0, 0, -20)
	snap := map[string]*ConceptSnapshot{
		"a": {Status: mastery.StatusMastered, LastReviewed: old, RequiredStreak: 3},
		"b": {Status: mastery.StatusMastered, LastReviewed: old, RequiredStreak: 3},
	}
	nx := s.NextReview(snap, "a")
	if nx == nil {
		t.Fatal("expected a review candidate")
	}
	if nx.Concept.ID == "a" {
		t.Error("NextReview must skip prevConceptID")
	}
}

func TestNextReview_NilWhenNothingDecayed(t *testing.T) {
	d := miniDAG(t)
	s := New(d)
	nx := s.NextReview(map[string]*ConceptSnapshot{}, "")
	if nx != nil {
		t.Errorf("expected nil with no decayed concepts, got %s", nx.Concept.ID)
	}
}

func TestDaysSince_ZeroTimeSentinel(t *testing.T) {
	if got := daysSince(time.Time{}); got != 999 {
		t.Errorf("expected 999 sentinel for zero time, got %f", got)
	}
}
