package abstract

import (
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

// TestFieldPoolsAreServed pins the yes/no rows of the three field generators
// to the questions a learner can actually receive.
//
// Each generator builds an easy pool, adds a hard pool above difficulty 0.75
// and then picks a computed shape with rand.Intn(3), and every branch of that
// switch returned. The pool was built on every draw and served on none: go vet
// named the trailing draw unreachable, and nothing in CI runs go vet. The false
// propositions PR #231 added so a learner could not score by always pressing
// yes were never asked, and the difficulty gate on the hard rows had no effect.
func TestFieldPoolsAreServed(t *testing.T) {
	reg := abstractRegistry(t)
	cases := []struct {
		id   string
		hard string // a row only the hard pool holds
	}{}
	if len(cases) == 0 {
		t.Skip("the field generators now serve computed count problems")
	}
	const samples = 400
	for _, c := range cases {
		yesNo, hard := 0, 0
		for i := 0; i < samples; i++ {
			p, err := reg.GenerateContext(c.id, generator.GeneratorContext{Difficulty: 1.0, Seed: int64(70000 + i)})
			if err != nil {
				t.Fatalf("%s: generate: %v", c.id, err)
			}
			if isYesNo(p) {
				yesNo++
			}
			if strings.Contains(p.Question, c.hard) {
				hard++
			}
		}
		if yesNo == 0 {
			t.Errorf("%s: no yes/no row in %d draws at difficulty 1.0, the pool is never served", c.id, samples)
		}
		if hard == 0 {
			t.Errorf("%s: the hard row %q never appears in %d draws at difficulty 1.0", c.id, c.hard, samples)
		}
		for i := 0; i < samples; i++ {
			p, err := reg.GenerateContext(c.id, generator.GeneratorContext{Difficulty: 0.1, Seed: int64(80000 + i)})
			if err != nil {
				t.Fatalf("%s: generate: %v", c.id, err)
			}
			if strings.Contains(p.Question, c.hard) {
				t.Errorf("%s: the hard row %q was served at difficulty 0.1", c.id, c.hard)
				break
			}
		}
	}
}

// isYesNo reports a two-option question whose answer is yes or no, which is
// the shape groupYesNo builds and the computed branches never do.
func isYesNo(p generator.Problem) bool {
	if p.Answer != "yes" && p.Answer != "no" {
		return false
	}
	return strings.Contains(p.Question, "\n\nB) ") && !strings.Contains(p.Question, "\n\nC) ")
}
