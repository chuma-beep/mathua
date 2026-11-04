package algebra

import (
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func TestQuadCompleteSquareDoesNotRepeatConsecutively(t *testing.T) {
	gen := &quadCompleteSquareGen{}
	prev := ""
	for i := 0; i < 50; i++ {
		p := gen.Generate(generator.GeneratorContext{})
		if p.Question == prev {
			t.Fatalf("consecutive duplicate question: %q", p.Question)
		}
		if strings.Contains(p.Question, "+ -") {
			t.Fatalf("malformed signed term: %q", p.Question)
		}
		prev = p.Question
	}
}
