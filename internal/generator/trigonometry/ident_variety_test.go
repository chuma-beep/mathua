package trigonometry

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func TestElevatedTrigVariety(t *testing.T) {
	rand.Seed(42)
	reg := generator.NewRegistry()
	Register(reg)
	conceptIDs := []string{
		"trig.ident.product_to_sum",
		"trig.ident.sum_to_product",
		"trig.ident.pythagorean_variants",
		"trig.eq.general_solution",
		"trig.ineq.advanced",
		"trig.hyperbolic.identities",
	}

	for _, id := range conceptIDs {
		gen, err := reg.Get(id)
		if err != nil {
			t.Fatalf("concept %s not registered: %v", id, err)
		}

		stems := make(map[string]bool)
		answers := make(map[string]bool)
		for i := 0; i < 50; i++ {
			ctx := generator.GeneratorContext{Difficulty: 0.5}
			prob := gen.Generate(ctx)

			if strings.TrimSpace(prob.Question) == "" {
				t.Errorf("%s: empty question at draw %d", id, i)
			}
			if strings.TrimSpace(prob.Answer) == "" {
				t.Errorf("%s: empty answer at draw %d", id, i)
			}
			if strings.TrimSpace(prob.Explanation) == "" {
				t.Errorf("%s: empty explanation at draw %d", id, i)
			}

			answers[prob.Answer] = true

			stems[prob.Question] = true
		}

		if len(stems) < 10 {
			t.Errorf("%s: variety too low, got %d distinct stems across 50 draws", id, len(stems))
		}
		if len(answers) < 5 {
			t.Errorf("%s: answer variety too low, got %d distinct counts", id, len(answers))
		}
	}
}
