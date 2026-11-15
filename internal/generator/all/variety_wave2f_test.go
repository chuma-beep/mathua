package all

import (
	"math/rand"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

// Wave 2f of the variety pass covers the narrow Calculus I concepts. Each of
// the sixteen concepts below answered at least one knowledge point from a fixed
// table or a numeric range of nine values or fewer, so a learner reviewing that
// knowledge point met the same handful of questions.
//
// Every generator now branches on KPOrdinal and draws its parameters, so one
// ordinal alone supplies at least ten distinct problems. The generic shipping
// contract only asks for five in twelve draws, which the old tables passed, so
// this guard holds the wider branches at ten over 200 draws.
func TestWave2FCalc1Variety(t *testing.T) {
	concepts := []string{
		"calc.deriv.cauchy_mvt",
		"calc.deriv.differential",
		"calc.deriv.hyperbolic",
		"calc.deriv.normal_line",
		"calc.deriv.optimization_cylinder",
		"calc.deriv.optimization_distance",
		"calc.deriv.optimization_fence",
		"calc.deriv.related_rates_cone",
		"calc.deriv.related_rates_ladder",
		"calc.deriv.rolle",
		"calc.integral.absolute_value",
		"calc.integral.floor",
		"calc.integral.total_distance",
		"calc.limit.absolute_value",
		"calc.limit.complex_fraction",
		"calc.limit.continuity_extension",
	}

	const draws = 200
	reg := Registry()

	for _, id := range concepts {
		for kp := 1; kp <= 3; kp++ {
			seen := make(map[string]bool)
			for i := 0; i < draws; i++ {
				p, err := reg.GenerateContext(id, generator.GeneratorContext{
					Difficulty: 0.2 + 0.8*rand.Float64(),
					KPOrdinal:  kp,
				})
				if err != nil {
					t.Fatalf("%s KP%d: %v", id, kp, err)
				}
				// The ordinal is part of the contract: a problem served for
				// knowledge point 2 must say so, or the alignment check in
				// contract_test.go cannot see a mis-filed question.
				if p.KPOrdinal != kp {
					t.Fatalf("%s KP%d: problem records ordinal %d", id, kp, p.KPOrdinal)
				}
				seen[canonicalProblem(p.Question)] = true
			}
			if len(seen) < 10 {
				t.Errorf("%s KP%d: %d distinct questions in %d draws, want at least 10", id, kp, len(seen), draws)
			}
			t.Logf("%s KP%d: %d distinct questions", id, kp, len(seen))
		}
	}
}
