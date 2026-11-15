package all

import (
	"math/rand"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

// Wave 2d of the variety pass covers the complex numbers domain. The twenty
// concepts below each drew their three knowledge points from a fixed table of
// six to eight written questions or a numeric range narrow enough to repeat
// inside one practice session. A learner reviewing the concept met the same
// handful of questions at every ordinal.
//
// Each generator now parameterises its branches and honours KPOrdinal, so one
// ordinal alone supplies at least ten distinct problems. The generic shipping
// contract only asks for five in twelve draws, which the old tables passed, so
// this guard holds the wider branches at ten over 200 draws.
func TestWave2DComplexVariety(t *testing.T) {
	concepts := []string{
		"complex.adv.branches",
		"complex.adv.elliptic",
		"complex.adv.gamma",
		"complex.adv.harmonic_conjugate",
		"complex.adv.jensen",
		"complex.adv.logarithm",
		"complex.adv.mittag_leffler",
		"complex.adv.modular_forms",
		"complex.adv.monodromy",
		"complex.adv.schwarz_christoffel",
		"complex.adv.trig_zeros",
		"complex.adv.analytic_continuation",
		"complex.adv.casorati",
		"complex.adv.open_mapping",
		"complex.adv.picard",
		"complex.adv.riemann_surface",
		"complex.adv.roots_unity",
		"complex.adv.weierstrass_factorization",
		"complex.adv.complex_powers",
		"complex.ops.modulus",
	}

	const draws = 200
	reg := generator.NewRegistry()
	Register(reg)

	for _, id := range concepts {
		for kp := 1; kp <= 3; kp++ {
			seen := make(map[string]bool)
			for i := 0; i < draws; i++ {
				p, err := reg.GenerateContext(id, generator.GeneratorContext{
					Difficulty: rand.Float64(),
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
