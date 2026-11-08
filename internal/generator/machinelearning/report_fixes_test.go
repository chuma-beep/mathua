package machinelearning

import (
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

// A question must not contain its own answer. One backprop entry used to state
// the dependency and then ask which half of it was true.
func TestBackpropQuestionDoesNotStateItsAnswer(t *testing.T) {
	g := &backpropGen{}
	for i := 0; i < 500; i++ {
		p := g.Generate(generator.GeneratorContext{})
		if p.Answer == "" {
			t.Fatalf("empty answer for %q", p.Question)
		}
		if strings.Contains(strings.ToLower(p.Question), strings.ToLower(p.Answer)) {
			t.Fatalf("question states its answer %q: %q", p.Answer, p.Question)
		}
	}
}
