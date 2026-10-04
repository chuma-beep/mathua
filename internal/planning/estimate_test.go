package planning

import (
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/mastery"
	"github.com/chuma-beep/mathua/internal/storage"
)

func ptrTime(t time.Time) *time.Time { return &t }

func testDAG(t *testing.T) *concepts.DAG {
	t.Helper()
	raw := []concepts.Concept{
		{ID: "a", Label: "A", Domain: "math", GradingType: "numeric", MasteryThreshold: concepts.MasteryThreshold{Streak: 4, AvgTimeSeconds: 10}},
		{ID: "b", Label: "B", Domain: "math", GradingType: "numeric", Prerequisites: []string{"a"}, MasteryThreshold: concepts.MasteryThreshold{Streak: 4, AvgTimeSeconds: 20}},
		{ID: "c.word", Label: "C word", Domain: "words", GradingType: "numeric", Prerequisites: []string{"b"}, MasteryThreshold: concepts.MasteryThreshold{Streak: 2, AvgTimeSeconds: 30}},
	}
	dag, err := concepts.Build(raw)
	if err != nil {
		t.Fatalf("build dag: %v", err)
	}
	return dag
}

func estimatePlanner(t *testing.T) *Planner {
	t.Helper()
	p := New(testDAG(t))
	if err := func() error {
		p.byID["t1"] = &Course{ID: "t1", Name: "T1", Targets: []string{"c.word"}}
		return nil
	}(); err != nil {
		t.Fatal(err)
	}
	return p
}

// Expected earned XP collapses to remaining-streak x base: accuracy cancels
// (more attempts at a lower hit rate), so estimates can't be gamed by rate.
func TestEstimateXPCollapse(t *testing.T) {
	p := estimatePlanner(t)
	path, err := p.PrerequisitesOf([]string{"b"})
	if err != nil {
		t.Fatal(err)
	}
	lo := EstimateWorkload(path, nil, map[string]AttemptStats{"a": {Correct: 1, Failed: 4}}, EstimateOpts{DailyXPRate: 50})
	hi := EstimateWorkload(path, nil, map[string]AttemptStats{"a": {Correct: 9, Failed: 1}}, EstimateOpts{DailyXPRate: 50})
	var xa, xb float64
	for _, c := range lo.Costs {
		if c.ConceptID == "a" {
			xa = c.XP
		}
	}
	for _, c := range hi.Costs {
		if c.ConceptID == "a" {
			xb = c.XP
		}
	}
	if xa != xb {
		t.Errorf("XP for a differs by accuracy: %v vs %v (must collapse to streak x base)", xa, xb)
	}
	if xa != 4*1 {
		t.Errorf("XP for a = %v, want 4*1=4", xa)
	}
}

// Word problems price at the multistep base.
func TestEstimateWordBase(t *testing.T) {
	p := estimatePlanner(t)
	path, err := p.PrerequisitesOf([]string{"c.word"})
	if err != nil {
		t.Fatal(err)
	}
	est := EstimateWorkload(path, nil, nil, EstimateOpts{DailyXPRate: 50})
	for _, c := range est.Costs {
		if c.ConceptID == "c.word" && c.XP != 2*1 {
			t.Errorf("c.word XP = %v, want 2*1=2", c.XP)
		}
	}
}

// Anti-gaming floor: don't-knows add re-orientation time and never shrink
// the total below the no-dontknow estimate.
func TestEstimateDontKnowFloor(t *testing.T) {
	p := estimatePlanner(t)
	path, err := p.PrerequisitesOf([]string{"a"})
	if err != nil {
		t.Fatal(err)
	}
	base := EstimateWorkload(path, nil, map[string]AttemptStats{"a": {Correct: 2, Failed: 2}}, EstimateOpts{DailyXPRate: 50})
	gamed := EstimateWorkload(path, nil, map[string]AttemptStats{"a": {Correct: 2, Failed: 2, DontKnow: 10}}, EstimateOpts{DailyXPRate: 50})
	if gamed.XPRemaining < base.XPRemaining {
		t.Errorf("dontknows shrank XP: %v < %v", gamed.XPRemaining, base.XPRemaining)
	}
	if gamed.TimeMinRemaining < base.TimeMinRemaining {
		t.Errorf("dontknows shrank time: %v < %v", gamed.TimeMinRemaining, base.TimeMinRemaining)
	}
	// And the time floor holds at full-threshold pace per remaining streak.
	var tm float64
	for _, c := range gamed.Costs {
		if c.ConceptID == "a" {
			tm = c.TimeMin
		}
	}
	if tm < 4*10.0/60.0 {
		t.Errorf("a time %v below floor 4*10s", tm)
	}
}

// Quizzes: one per 150 XP of learning effort, time-only.
func TestEstimateQuizzes(t *testing.T) {
	p := estimatePlanner(t)
	path, err := p.PrerequisitesOf([]string{"c.word"})
	if err != nil {
		t.Fatal(err)
	}
	est := EstimateWorkload(path, nil, nil, EstimateOpts{DailyXPRate: 50})
	// XP: a=4 + b=4 + c.word=2 = 10 -> 1 quiz.
	if est.QuizzesAhead != 1 {
		t.Errorf("quizzes = %d, want 1 for 110 XP", est.QuizzesAhead)
	}
	if est.AssessmentMin != 5*30.0/60.0 {
		t.Errorf("assessment min = %v, want 2.5", est.AssessmentMin)
	}
	if est.Lines.Assessment != 0 || est.Lines.Diagnostic != 0 {
		t.Errorf("assessment/diagnostic lines must be 0 XP debt, got %+v", est.Lines)
	}
}

// Deadline mode: required rate + feasibility.
func TestEstimateDeadline(t *testing.T) {
	p := estimatePlanner(t)
	path, err := p.PrerequisitesOf([]string{"b"})
	if err != nil {
		t.Fatal(err)
	}
	// XP: 8 over 10 days, no rest -> 0.8/day, feasible under cap 50.
	est := EstimateWorkload(path, nil, nil, EstimateOpts{DailyXPRate: 50, DeadlineDays: 10, FeasibleCap: 50})
	if est.RequiredPerDay != 0.8 {
		t.Errorf("required/day = %v, want 0.8", est.RequiredPerDay)
	}
	if !est.Feasible {
		t.Errorf("should be feasible under cap 50")
	}
	tight := EstimateWorkload(path, nil, nil, EstimateOpts{DailyXPRate: 50, DeadlineDays: 2, FeasibleCap: 3})
	if tight.Feasible {
		t.Errorf("4/day over 2 days should be infeasible under cap 3, got required=%v", tight.RequiredPerDay)
	}
}

// Rest days reduce the effective rate.
func TestEstimateRestDays(t *testing.T) {
	p := estimatePlanner(t)
	path, err := p.PrerequisitesOf([]string{"b"})
	if err != nil {
		t.Fatal(err)
	}
	full := EstimateWorkload(path, nil, nil, EstimateOpts{DailyXPRate: 70})
	rest := EstimateWorkload(path, nil, nil, EstimateOpts{DailyXPRate: 70, RestDaysPerWeek: 2})
	if rest.Days <= full.Days {
		t.Errorf("rest days should extend estimate: %v vs %v", rest.Days, full.Days)
	}
}

// Sampler is deterministic, domain-spread, and capped.
func TestSampleClosure(t *testing.T) {
	p := estimatePlanner(t)
	path, err := p.PrerequisitesOf([]string{"c.word"})
	if err != nil {
		t.Fatal(err)
	}
	a := SampleClosure(path, 2)
	b := SampleClosure(path, 2)
	if len(a) != 2 || len(b) != 2 || a[0] != b[0] || a[1] != b[1] {
		t.Errorf("sampler not deterministic: %v vs %v", a, b)
	}
	// math: a,b + words: c.word -> first two span both domains.
	if got := SampleClosure(path, 10); len(got) != 3 {
		t.Errorf("uncapped sample = %v, want all 3", got)
	}
}

// Destinations resolve to union closures with shared prereqs counted once.
func TestDestinationUnion(t *testing.T) {
	p := estimatePlanner(t)
	p.SetDestinations([]*Destination{{ID: "d", Name: "D", Courses: []string{"t1"}}})
	path, err := p.PathForDestination("d")
	if err != nil {
		t.Fatal(err)
	}
	if len(path.Concepts) != 3 {
		t.Errorf("union has %d concepts, want 3 (a,b,c.word)", len(path.Concepts))
	}
	if _, err := p.PathForDestination("nope"); err == nil {
		t.Errorf("expected error for unknown destination")
	}
}

// The estimator and the scheduler must agree on when a mastered concept needs a
// review pass.
//
// They did not. `defaultDecayDays` here was 30 while `scheduler.effectiveState`
// and the engine used 14, so a plan priced reviews on a cadence the app would not
// honour — and reviews are priced on the same path the deadline is computed from
// (ADR-019). The constant is now shared; this pins the agreement at the boundary
// so a future divergence is a failing test rather than a silently optimistic date.
func TestEstimate_AgreesWithSchedulerOnDecay(t *testing.T) {
	p := estimatePlanner(t)
	path, err := p.PrerequisitesOf([]string{"a"})
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		name      string
		daysSince float64
		wantDue   bool
	}{
		{"fresh just inside the window", mastery.DecayDays - 1, false},
		{"due exactly at the window", mastery.DecayDays, true},
		{"due past the window", mastery.DecayDays + 10, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			reviewed := time.Now().UTC().Add(-time.Duration(tc.daysSince*24) * time.Hour)
			// NextReviewDue is deliberately in the future so the SM-2 schedule is not
			// what makes it due — this is purely the decay window under test.
			progress := map[string]*storage.ConceptProgress{
				"a": {Status: "MASTERED", LastReviewed: &reviewed, NextReviewDue: ptrTime(time.Now().UTC().AddDate(0, 0, 90))},
			}
			est := EstimateWorkload(path, progress, map[string]AttemptStats{}, EstimateOpts{})
			if got := est.ReviewsDue; (got > 0) != tc.wantDue {
				t.Errorf("ReviewsDue = %d with %.0f days since review, want due=%v", got, tc.daysSince, tc.wantDue)
			}
			// The same fact read through the shared constant the scheduler uses.
			if got := mastery.EffectiveStatus(mastery.StatusMastered, tc.daysSince, mastery.DecayDays); (got == "DECAYING") != tc.wantDue {
				t.Errorf("EffectiveStatus(%.0f days) = %q, want decaying=%v", tc.daysSince, got, tc.wantDue)
			}
		})
	}
}
