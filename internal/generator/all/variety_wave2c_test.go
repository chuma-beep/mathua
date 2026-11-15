package all

import (
	"math/rand"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

// Wave 2c of the variety pass rewrote the algebra generators whose three
// knowledge points each drew from a pool of fewer than ten written questions or
// a narrow numeric range. A learner reviewing one of these concepts met the
// same handful of questions at every ordinal, because the branch either ignored
// which knowledge point the lesson said they were working on or drew from a
// fixed table of six entries.
//
// Every one of these concepts now draws its parameters and branches on
// KPOrdinal, so one ordinal alone supplies at least ten distinct problems. The
// generic shipping contract only asks for five in twelve draws, which the old
// pools passed, so this guard holds the wider branches at ten over 200 draws.
func TestWave2CAlgebraVariety(t *testing.T) {
	concepts := []string{
		"alg.conic.parabola_focus",
		"alg.eq.trinomial",
		"alg.exp.compound",
		"alg.exp.graph",
		"alg.func.composite_domain",
		"alg.func.even_odd",
		"alg.log.model",
		"alg.log.scale",
		"alg.radical.operations",
		"alg.rational.proportions",
		"alg.seq.sigma",
		"alg.systems.nonlinear",
		"alg.conic.ellipse_foci",
		"alg.eq.irrational",
		"alg.factor.diff_squares",
		"alg.func.piecewise",
		"alg.ineq.irrational",
		"alg.radical.simplify",
		"alg.eq.logarithmic",
		"alg.func.sign",
		"alg.ineq.logarithmic",
		"alg.quad.complete_square",
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
