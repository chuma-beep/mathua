package generator_test

import (
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/solutions"
)

// TestIneqRationalExplanationFollowsTheQuestion pins the shipped schema for
// alg.ineq.rational to both questions its generator asks.
//
// The generator asks for the smallest integer with 1/(x-a) > 0 half the time
// and for the largest integer with 1/(x-a) < 0 the other half. The schema
// narrated only the first, so a learner who was asked the second read the
// solution to a different question, with a-1 explained as the smallest whole
// number above a. Reported in #267 from a diagnostic.
func TestIneqRationalExplanationFollowsTheQuestion(t *testing.T) {
	l, err := solutions.Load(solutionsDir)
	if err != nil {
		t.Fatalf("load shipped schemas: %v", err)
	}
	if _, ok := l.Schema("alg.ineq.rational"); !ok {
		t.Fatal("alg.ineq.rational has no shipped schema")
	}
	reg := generator.NewRegistry()
	reg.SetSolutions(l)
	registerAllDomains(reg)

	seen := map[string]bool{}
	for i := 0; i < 40; i++ {
		seed := int64(4000 + i)
		p, err := reg.GenerateContext("alg.ineq.rational", generator.GeneratorContext{Difficulty: 0.5, Seed: seed})
		if err != nil {
			t.Fatalf("seed %d: generate: %v", seed, err)
		}
		which, other, sign := "smallest", "largest", "> 0"
		if strings.Contains(p.Question, "largest") {
			which, other, sign = "largest", "smallest", "< 0"
		}
		seen[which] = true
		if !strings.Contains(p.Explanation, which) || strings.Contains(p.Explanation, other) || !strings.Contains(p.Explanation, sign) {
			t.Errorf("seed %d asks for the %s integer with 1/(x-a) %s, but the explanation reads:\n%s", seed, which, sign, p.Explanation)
		}
	}
	if len(seen) != 2 {
		t.Fatalf("40 draws covered only %v, both questions must be covered", seen)
	}
}
