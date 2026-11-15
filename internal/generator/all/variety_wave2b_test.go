package all

import (
	"math/rand"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

// Wave 2b of the variety pass rewrote the differential-equations generators
// whose three knowledge points each drew from a pool of at most nine written
// questions. A learner reviewing one of these concepts met the same handful of
// questions at every ordinal, because the branch ignored which knowledge point
// the lesson said they were working on.
//
// Every one of these concepts now draws its parameters and branches on
// KPOrdinal, so one ordinal alone supplies at least ten distinct problems. The
// generic shipping contract only asks for five in twelve draws, which the old
// pools passed, so this guard holds the wider branches at ten over 200 draws.
func TestWave2BODEVariety(t *testing.T) {
	concepts := []string{
		"ode.dynamical.gradient_flow",
		"ode.dynamical.hartman_grobman",
		"ode.first_order.mixing",
		"ode.first_order.riccati",
		"ode.laplace.value",
		"ode.pde.duhamel",
		"ode.pde.polar_separation",
		"ode.second_order.cauchy_euler",
		"ode.second_order.power_series",
		"ode.second_order.spring",
		"ode.stiff.a_stability",
		"ode.stiff.bdf",
		"ode.stiff.l_stability",
		"ode.stiff.symplectic",
		"ode.stiff.trapezoidal",
		"ode.systems.predator_prey",
		"ode.variational.brachistochrone",
		"ode.variational.euler_lagrange",
		"ode.variational.isoperimetric",
		"ode.variational.legendre",
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
