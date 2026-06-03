package diagnostic

import (
	"testing"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/generator"
)

type fakeGen struct{}

func (f *fakeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return generator.Problem{
		Question:    "test?",
		Answer:      "test",
		Explanation: "test",
	}
}

func testDAG(t *testing.T) *concepts.DAG {
	t.Helper()
	raw := make([]concepts.Concept, 10)
	for i := range raw {
		raw[i] = concepts.Concept{
			ID:               idFor(i),
			Domain:           "d",
			Prerequisites:    prereqsFor(i),
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 10},
		}
	}
	d, err := concepts.Build(raw)
	if err != nil {
		t.Fatalf("build DAG: %v", err)
	}
	return d
}

func idFor(i int) string {
	letters := "abcdefghij"
	return string(letters[i])
}

func prereqsFor(i int) []string {
	if i == 0 {
		return nil
	}
	return []string{idFor(i - 1)}
}

func mustRegistry(t *testing.T) *generator.Registry {
	t.Helper()
	r := generator.NewRegistry()
	for i := 0; i < 10; i++ {
		if err := r.Register(idFor(i), &fakeGen{}); err != nil {
			t.Fatalf("register: %v", err)
		}
	}
	return r
}

func TestEngine_Start(t *testing.T) {
	e := NewEngine(testDAG(t), mustRegistry(t))
	s := e.Start()
	if s.State != StateProbing {
		t.Errorf("expected probing, got %q", s.State)
	}
	if s.Position != 5 {
		t.Errorf("expected position 5 (midpoint of 10), got %d", s.Position)
	}
}

func TestEngine_NextQuestion(t *testing.T) {
	e := NewEngine(testDAG(t), mustRegistry(t))
	s := e.Start()
	p, cid, err := e.NextQuestion(s)
	if err != nil {
		t.Fatalf("next question: %v", err)
	}
	if p.Question == "" || cid == "" {
		t.Error("expected non-empty question and concept ID")
	}
	if s.LastConceptID != cid {
		t.Errorf("expected LastConceptID=%q, got %q", cid, s.LastConceptID)
	}
}

func TestEngine_RecordAnswer_AdvancesAfterTwoCorrect(t *testing.T) {
	e := NewEngine(testDAG(t), mustRegistry(t))
	s := e.Start()
	initialPos := s.Position
	cid := idFor(initialPos)
	// First correct+fast — not enough to evaluate
	e.RecordAnswer(s, cid, true, true)
	if s.totalCount[cid] != 1 {
		t.Errorf("expected 1 probe, got %d", s.totalCount[cid])
	}
	// Second correct+fast — should advance
	e.RecordAnswer(s, cid, true, true)
	if s.Low <= initialPos {
		t.Errorf("expected low > %d after two correct, got %d", initialPos, s.Low)
	}
}

func TestEngine_RecordAnswer_RetreatsAfterTwoWrong(t *testing.T) {
	e := NewEngine(testDAG(t), mustRegistry(t))
	s := e.Start()
	initialPos := s.Position
	cid := idFor(initialPos)
	e.RecordAnswer(s, cid, false, false)
	e.RecordAnswer(s, cid, false, false)
	if s.High >= initialPos {
		t.Errorf("expected high < %d after two wrong, got %d", initialPos, s.High)
	}
}

func TestEngine_CompletesWhenLowCrossesHigh(t *testing.T) {
	e := NewEngine(testDAG(t), mustRegistry(t))
	s := e.Start()
	// Simulate all concepts wrong to collapse the range quickly
	for i := 0; i < len(s.order) && !e.IsComplete(s); i++ {
		cid := idFor(i)
		e.RecordAnswer(s, cid, false, false)
		e.RecordAnswer(s, cid, false, false)
	}
	if !e.IsComplete(s) {
		t.Errorf("expected complete after all wrong, got state=%s low=%d high=%d asked=%d",
			s.State, s.Low, s.High, s.totalAsked)
	}
}

func TestEngine_FrontierEstimate(t *testing.T) {
	e := NewEngine(testDAG(t), mustRegistry(t))
	s := e.Start()
	if e.FrontierEstimate(s) != 9 {
		t.Errorf("expected frontier 9 (high), got %d", e.FrontierEstimate(s))
	}
	// Two wrong answers at each position collapses the frontier
	for i := 0; i < 5; i++ {
		if s.Position < 0 || s.Position >= len(s.order) {
			break
		}
		cid := idFor(s.Position)
		e.RecordAnswer(s, cid, false, false)
		e.RecordAnswer(s, cid, false, false)
	}
	if e.FrontierEstimate(s) > 4 {
		t.Errorf("expected frontier <= 4 after wrongs at midpoint, got %d", e.FrontierEstimate(s))
	}
}

func TestEngine_AllProbedCompletes(t *testing.T) {
	e := NewEngine(testDAG(t), mustRegistry(t))
	s := e.Start()
	// Probe every concept at least once
	for i := 0; i < len(s.order) && !e.IsComplete(s); i++ {
		cid := idFor(i)
		if s.totalCount[cid] < probesPerConcept {
			e.RecordAnswer(s, cid, true, true)
			e.RecordAnswer(s, cid, true, true)
		}
	}
	if !e.IsComplete(s) {
		t.Errorf("expected complete after all probed, state=%s asked=%d", s.State, s.totalAsked)
	}
}
