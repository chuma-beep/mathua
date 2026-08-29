package quiz

import (
	"testing"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/generator"
)

type fakeGen struct {
	difficulty float64
}

func (f *fakeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	f.difficulty = ctx.Difficulty
	return generator.Problem{Question: "q?", Answer: "a", Explanation: "e"}
}

func miniDAG(t *testing.T) *concepts.DAG {
	t.Helper()
	d, err := concepts.Build([]concepts.Concept{
		{ID: "a", Domain: "d", Subdomain: "s1", Prerequisites: []string{}, MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 10}},
		{ID: "b", Domain: "d", Subdomain: "s2", Prerequisites: []string{}, MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 10}},
		{ID: "c", Domain: "d", Subdomain: "s3", Prerequisites: []string{}, MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 10}},
	})
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	return d
}

func TestEngine_NextQuestion_UsesDifficultyMap(t *testing.T) {
	d := miniDAG(t)
	g := &fakeGen{}
	r := generator.NewRegistry()
	for _, id := range []string{"a", "b", "c"} {
		if err := r.Register(id, g); err != nil {
			t.Fatalf("register: %v", err)
		}
	}
	e := NewEngine(d, r)
	sess := e.Start(d.Order())
	// Weak concept "a" → easy difficulty 0.35 (per weakness map); strong "b" → 0.8.
	sess.Difficulties["a"] = 0.35
	sess.Difficulties["b"] = 0.8
	_, cid, err := e.NextQuestion(sess)
	if err != nil {
		t.Fatalf("next: %v", err)
	}
	if cid == "" {
		t.Fatal("expected a question")
	}
	if g.difficulty != 0.35 {
		t.Errorf("expected difficulty 0.35 for weak a, got %f", g.difficulty)
	}
	// Fallback when no map entry: 0.8 target.
	sess.Order = []*concepts.Concept{d.Concept("b")}
	sess.Index = 0
	if _, _, err := e.NextQuestion(sess); err != nil {
		t.Fatalf("next b: %v", err)
	}
	if g.difficulty != 0.8 {
		t.Errorf("expected fallback difficulty 0.8, got %f", g.difficulty)
	}
}

func TestPickQuizConcepts_DiverseSubdomains(t *testing.T) {
	d := miniDAG(t)
	picked := PickQuizConcepts(d, map[string]float64{"a": 0.8, "b": 0.7, "c": 0.6})
	if len(picked) != 3 {
		t.Fatalf("expected 3 picked, got %d", len(picked))
	}
	seen := map[string]bool{}
	for _, c := range picked {
		if seen[c.Subdomain] {
			t.Errorf("duplicate subdomain in picked set: %+v", picked)
		}
		seen[c.Subdomain] = true
	}
}

func TestPickQuizConcepts_FallbackOnNoWeak(t *testing.T) {
	d := miniDAG(t)
	picked := PickQuizConcepts(d, map[string]float64{})
	if len(picked) == 0 {
		t.Fatal("expected fallback picks with no weak concepts")
	}
	if len(picked) > 5 {
		t.Errorf("expected at most 5 picks, got %d", len(picked))
	}
}

func TestEngine_RecordAnswerAdvances(t *testing.T) {
	d := miniDAG(t)
	r := generator.NewRegistry()
	_ = r.Register("a", &fakeGen{})
	e := NewEngine(d, r)
	sess := e.Start(d.Order())
	e.RecordAnswer(sess, "a", true)
	if sess.Index != 1 {
		t.Errorf("expected index advance to 1, got %d", sess.Index)
	}
	if e.IsComplete(sess) {
		t.Error("expected not complete after 1 of 3")
	}
}
