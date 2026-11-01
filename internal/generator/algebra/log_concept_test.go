package algebra

import (
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func TestLogConceptFollowsTheKnowledgePoint(t *testing.T) {
	g := &logConceptGen{}

	definition := g.Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: 1})
	if !strings.Contains(definition.Question, `\log_`) {
		t.Fatalf("KP1 is not the exponential-to-logarithm prompt: %q", definition.Question)
	}
	if _, err := strconv.Atoi(definition.Answer); err != nil {
		t.Fatalf("KP1 answer %q is not a number", definition.Answer)
	}

	for i := 0; i < 30; i++ {
		p := g.Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: 2})
		if strings.Contains(p.Question, "Write as a logarithm") {
			t.Fatalf("KP2 still asks the definition prompt: %q", p.Question)
		}
		if _, err := strconv.Atoi(p.Answer); err != nil {
			t.Fatalf("KP2 answer %q is not a number: %q", p.Answer, p.Question)
		}
	}

	for i := 0; i < 30; i++ {
		p := g.Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: 3})
		if strings.Contains(p.Question, "Write as a logarithm") {
			t.Fatalf("KP3 still asks the definition prompt: %q", p.Question)
		}
		if _, err := strconv.Atoi(p.Answer); err != nil {
			t.Fatalf("KP3 answer %q is not a number: %q", p.Answer, p.Question)
		}
	}
}
