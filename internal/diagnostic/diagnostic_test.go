package diagnostic

import (
	"testing"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/generator"
)

type fakeGen struct{}

func (f *fakeGen) Generate(difficulty float64) generator.Problem {
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
}

func TestEngine_RecordAnswer_Advances(t *testing.T) {
	e := NewEngine(testDAG(t), mustRegistry(t))
	s := e.Start()
	initialPos := s.Position
	e.RecordAnswer(s, idFor(initialPos), true, true)
	if s.Position <= initialPos {
		t.Errorf("expected position > %d after correct+fast, got %d", initialPos, s.Position)
	}
	if s.Low != initialPos+1 {
		t.Errorf("expected low=%d, got %d", initialPos+1, s.Low)
	}
}

func TestEngine_RecordAnswer_Retreats(t *testing.T) {
	e := NewEngine(testDAG(t), mustRegistry(t))
	s := e.Start()
	initialPos := s.Position
	e.RecordAnswer(s, idFor(initialPos), false, false)
	if s.Position >= initialPos {
		t.Errorf("expected position < %d after incorrect, got %d", initialPos, s.Position)
	}
}

func TestEngine_CompletesAfter3Consecutive(t *testing.T) {
	e := NewEngine(testDAG(t), mustRegistry(t))
	s := e.Start()
	for _, cid := range []string{idFor(s.Position)} {
		e.RecordAnswer(s, cid, true, true)
		if e.IsComplete(s) {
			t.Logf("completed early at attempt %d: consecutive=%d low=%d high=%d pos=%d",
				len(s.Attempts), s.ConsecutiveOK, s.Low, s.High, s.Position)
			break
		}
	}
	// Keep answering correctly
	for i := 0; i < 10 && !e.IsComplete(s); i++ {
		if s.Position >= s.Low && s.Position <= s.High {
			cid := idFor(s.Position)
			e.RecordAnswer(s, cid, true, true)
		}
	}
	if !e.IsComplete(s) {
		t.Error("expected session to complete")
	}
}

func TestEngine_FrontierEstimate(t *testing.T) {
	e := NewEngine(testDAG(t), mustRegistry(t))
	s := e.Start()
	if e.FrontierEstimate(s) != 9 {
		t.Errorf("expected frontier 9 (high), got %d", e.FrontierEstimate(s))
	}
	e.RecordAnswer(s, idFor(s.Position), false, false)
	if e.FrontierEstimate(s) > 4 {
		t.Errorf("expected frontier <= 4 after wrong at midpoint, got %d", e.FrontierEstimate(s))
	}
}

func TestEngine_CrossesBounds(t *testing.T) {
	e := NewEngine(testDAG(t), mustRegistry(t))
	s := e.Start()
	s.Position = 0
	s.Low = 0
	s.High = 0
	s.ConsecutiveOK = 3
	e.RecordAnswer(s, idFor(0), true, true)
	if !e.IsComplete(s) {
		t.Error("expected complete when consecutive OK >= 3")
	}
}
