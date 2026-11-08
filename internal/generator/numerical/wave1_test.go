package numerical

import (
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func TestBrentBisectionGrammar(t *testing.T) {
	for i := 0; i < 400; i++ {
		p := (&brentGen{}).Generate(generator.GeneratorContext{KPOrdinal: 1})
		if strings.Contains(p.Question, "1 bisection steps") {
			t.Fatalf("single step rendered as plural: %q", p.Question)
		}
	}
}

func TestBackwardErrorExplanationNamesTheArithmetic(t *testing.T) {
	for _, kp := range []int{1, 2, 3} {
		p := (&backwardErrorGen{}).Generate(generator.GeneratorContext{KPOrdinal: kp})
		if !strings.Contains(p.Explanation, p.Answer) {
			t.Fatalf("KP%d explanation does not include the answer %s: %q", kp, p.Answer, p.Explanation)
		}
	}
}
