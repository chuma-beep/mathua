package generator

import (
	"testing"
)

type stubGen struct {
	fn func(difficulty float64) Problem
}

func (s *stubGen) Generate(ctx GeneratorContext) Problem {
	return s.fn(ctx.Difficulty)
}

func TestRegistry_RegisterAndGenerate(t *testing.T) {
	r := NewRegistry()
	g := &stubGen{fn: func(d float64) Problem {
		return Problem{Question: "q", Answer: "a", Explanation: "e"}
	}}
	if err := r.Register("test.concept", g); err != nil {
		t.Fatalf("unexpected register error: %v", err)
	}
	p, err := r.Generate("test.concept", 0.5)
	if err != nil {
		t.Fatalf("unexpected generate error: %v", err)
	}
	if p.Question != "q" || p.Answer != "a" || p.Explanation != "e" {
		t.Errorf("unexpected problem: %+v", p)
	}
}

func TestRegistry_DuplicateRegistration(t *testing.T) {
	r := NewRegistry()
	g := &stubGen{fn: func(d float64) Problem { return Problem{} }}
	_ = r.Register("c", g)
	err := r.Register("c", g)
	if err == nil {
		t.Error("expected error for duplicate registration")
	}
}

func TestRegistry_NilGenerator(t *testing.T) {
	r := NewRegistry()
	err := r.Register("c", nil)
	if err == nil {
		t.Error("expected error for nil generator")
	}
}

func TestRegistry_EmptyConceptID(t *testing.T) {
	r := NewRegistry()
	g := &stubGen{fn: func(d float64) Problem { return Problem{} }}
	err := r.Register("", g)
	if err == nil {
		t.Error("expected error for empty concept ID")
	}
}

func TestRegistry_UnknownConcept(t *testing.T) {
	r := NewRegistry()
	_, err := r.Generate("nonexistent", 0.5)
	if err == nil {
		t.Error("expected error for unknown concept")
	}
}

func TestRegistry_Has(t *testing.T) {
	r := NewRegistry()
	if r.Has("c") {
		t.Error("expected Has to return false before registration")
	}
	g := &stubGen{fn: func(d float64) Problem { return Problem{} }}
	_ = r.Register("c", g)
	if !r.Has("c") {
		t.Error("expected Has to return true after registration")
	}
}

func TestRegistry_Count(t *testing.T) {
	r := NewRegistry()
	if r.Count() != 0 {
		t.Errorf("expected count 0, got %d", r.Count())
	}
	g := &stubGen{fn: func(d float64) Problem { return Problem{} }}
	_ = r.Register("a", g)
	_ = r.Register("b", g)
	if r.Count() != 2 {
		t.Errorf("expected count 2, got %d", r.Count())
	}
}

func TestRegistry_ClampDifficulty(t *testing.T) {
	r := NewRegistry()
	g := &stubGen{fn: func(d float64) Problem {
		return Problem{Question: "ok", Answer: "ok"}
	}}
	_ = r.Register("c", g)
	if _, err := r.Generate("c", -0.5); err != nil {
		t.Errorf("expected no error for clamped low difficulty: %v", err)
	}
	if _, err := r.Generate("c", 1.5); err != nil {
		t.Errorf("expected no error for clamped high difficulty: %v", err)
	}
}

func TestRegistry_MultipleGenerators(t *testing.T) {
	r := NewRegistry()
	g1 := &stubGen{fn: func(d float64) Problem { return Problem{Question: "q1"} }}
	g2 := &stubGen{fn: func(d float64) Problem { return Problem{Question: "q2"} }}
	_ = r.Register("a", g1)
	_ = r.Register("b", g2)
	p1, _ := r.Generate("a", 0)
	p2, _ := r.Generate("b", 0.5)
	if p1.Question != "q1" {
		t.Errorf("expected q1, got %q", p1.Question)
	}
	if p2.Question != "q2" {
		t.Errorf("expected q2, got %q", p2.Question)
	}
}

func TestProblem_Fields(t *testing.T) {
	p := Problem{
		Question:    "What is 2+2?",
		Answer:      "4",
		Explanation: "2 + 2 = 4",
	}
	if p.Question == "" || p.Answer == "" || p.Explanation == "" {
		t.Error("problem fields should not be empty")
	}
}
