package all

import (
	"math/rand"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

// Wave 1 of the variety pass widened the fourteen knowledge points that could
// produce fewer than five distinct questions. The generic shipping contract
// only asks for five in twelve draws, so a regression back to a narrow pool
// would still pass it. This guard holds the widened branches at ten.
func TestWave1VarietyFloor(t *testing.T) {
	cases := []struct {
		concept string
		kp      int
	}{
		{"calc.integral.simpson", 1},
		{"calc.integral.simpson", 2},
		{"calc.integral.simpson", 3},
		{"calc.integral.area_between", 2},
		{"calc.integral.area_between", 3},
		{"calc.integral.trapezoid", 1},
		{"calc.integral.trapezoid", 3},
		{"calc.deriv.difference_quotient", 3},
		{"alg.eq.multi_step", 1},
		{"calc.deriv.cauchy_mvt", 2},
		{"calc.integral.center_mass", 3},
		{"nt.basics.prime_inf", 1},
		{"nt.basics.prime_inf", 2},
		{"nt.basics.prime_inf", 3},
	}

	reg := Registry()
	for _, c := range cases {
		seen := make(map[string]bool)
		for i := 0; i < 200; i++ {
			p, err := reg.GenerateContext(c.concept, generator.GeneratorContext{
				Difficulty: 0.2 + 0.8*rand.Float64(),
				KPOrdinal:  c.kp,
			})
			if err != nil {
				t.Fatalf("%s KP%d: %v", c.concept, c.kp, err)
			}
			seen[canonicalProblem(p.Question)] = true
		}
		if len(seen) < 10 {
			t.Errorf("%s KP%d: %d distinct questions, want at least 10", c.concept, c.kp, len(seen))
		}
	}
}
