package all

import (
	"math/rand"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

// Wave 2e of the variety pass covers the narrow trigonometry concepts. Each of
// the seventeen concepts below answered every knowledge point from a fixed
// table of six to nine written questions, or from a numeric range of eight or
// fewer values, so a learner reviewing the concept met the same handful of
// questions at every ordinal.
//
// Every generator now branches on KPOrdinal and draws its parameters, so one
// ordinal alone supplies at least ten distinct problems. The generic shipping
// contract only asks for five in twelve draws, which the old tables passed, so
// this guard holds the wider branches at ten over 200 draws.
func TestWave2ETrigVariety(t *testing.T) {
	concepts := []string{
		"trig.adv.arctan",
		"trig.adv.heron_area",
		"trig.adv.vector_projection",
		"trig.eq.basic",
		"trig.eq.homogeneous",
		"trig.eq.reciprocal",
		"trig.eq.using_identities",
		"trig.func.power_integrals",
		"trig.func.taylor",
		"trig.hyperbolic.catenary",
		"trig.hyperbolic.inverse",
		"trig.hyperbolic.sinh_cosh",
		"trig.hyperbolic.tanh_coth",
		"trig.ident.identities",
		"trig.ident.power_reduction",
		"trig.ident.tangent_triple",
		"trig.ineq.basic",
	}

	const draws = 200
	reg := generator.NewRegistry()
	Register(reg)

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
