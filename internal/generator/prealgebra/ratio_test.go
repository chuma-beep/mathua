package prealgebra

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

// Report 351: the prompt asks for simplest form, but the answer was stored
// unreduced, so a learner who simplified could be marked wrong.
func TestRatioConceptAnswerIsSimplified(t *testing.T) {
	g := &ratioConceptGen{}
	for i := 0; i < 200; i++ {
		p := g.Generate(generator.GeneratorContext{Difficulty: rand.Float64(), Seed: rand.Int63()})
		var a, b int
		if _, err := fmt.Sscanf(p.Question, "Write the ratio of %d to %d in simplest form.", &a, &b); err != nil {
			t.Fatalf("question %q: %v", p.Question, err)
		}
		div := ratioGCD(a, b)
		want := fmt.Sprintf("%d:%d", a/div, b/div)
		if p.Answer != want {
			t.Fatalf("answer = %q, want %q for %q", p.Answer, want, p.Question)
		}
	}
}

// The report answered 5/2 for 5:2 and was told the format was invalid. The
// three forms a learner writes are the same ratio; an unreduced pair is not
// the requested simplest form.
func TestRatioConceptAcceptsEquivalentForms(t *testing.T) {
	g := &ratioConceptGen{}
	for _, answer := range []string{"5:2", "5/2", "5 to 2", "5 / 2"} {
		if !g.Grade("5:2", answer).Correct {
			t.Errorf("grade(%q) = incorrect, want correct", answer)
		}
	}
	if g.Grade("5:2", "2:5").Correct {
		t.Error("the inverse ratio graded correct")
	}
	if g.Grade("3:2", "6:4").Correct {
		t.Error("an unreduced equivalent graded correct for a simplest-form question")
	}
}
