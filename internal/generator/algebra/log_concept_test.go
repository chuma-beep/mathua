package algebra

import (
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func TestLogConceptFollowsTheKnowledgePoint(t *testing.T) {
	g := &logConceptGen{}

	definition := g.Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: 1})
	if !strings.Contains(definition.Question, "Write as a logarithm") {
		t.Fatalf("KP1 is not the exponential-to-logarithm prompt: %q", definition.Question)
	}

	for i := 0; i < 30; i++ {
		p := g.Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: 2})
		if strings.Contains(p.Question, "Write as a logarithm") {
			t.Fatalf("KP2 still asks the definition prompt: %q", p.Question)
		}
		assertOneCorrectOption(t, p, optionsOf(p.Question))
	}

	for i := 0; i < 30; i++ {
		p := g.Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: 3})
		if strings.Contains(p.Question, "Write as a logarithm") {
			t.Fatalf("KP3 still asks the definition prompt: %q", p.Question)
		}
		assertOneCorrectOption(t, p, optionsOf(p.Question))
	}
}

func assertOneCorrectOption(t *testing.T, p generator.Problem, options []string) {
	t.Helper()
	if len(options) != 4 {
		t.Fatalf("%q has %d options, want 4", p.Question, len(options))
	}
	seen := 0
	for _, option := range options {
		if strings.TrimSpace(option) == strings.TrimSpace(p.Answer) {
			seen++
		}
	}
	if seen != 1 {
		t.Fatalf("%q has answer %q in %d options: %v", p.Question, p.Answer, seen, options)
	}
}
