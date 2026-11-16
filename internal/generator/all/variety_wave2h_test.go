package all

import (
	"math/rand"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

// Wave 2h of the variety pass covers fourteen real-analysis concepts that
// answered every knowledge point from a fixed table of six precomputed
// claims, from a narrow parameter pair, or from a generator that ignored the
// knowledge-point ordinal. A learner reviewing such a knowledge point met the
// same six questions and then saw them repeat, and a generator that ignored
// the ordinal could not serve the section its lesson actually taught.
//
// Every generator below now branches on KPOrdinal and draws its operands, so
// one ordinal alone supplies at least ten distinct problems. The generic
// shipping contract only asks for five in twelve draws, which the old tables
// passed, so this guard holds the wider branches at ten over 200 draws and
// checks that the ordinal is recorded on every returned problem.
func TestWave2HRealAnalysisVariety(t *testing.T) {
	concepts := []string{
		"ra.order.supremum",
		"ra.order.archimedean",
		"ra.seq.monotone",
		"ra.seq.series",
		"ra.limit.function",
		"ra.deriv.mvt",
		"ra.deriv.chain",
		"ra.deriv.taylor",
		"ra.deriv.lhopital",
		"ra.deriv.riemann",
		"ra.deriv.parts",
		"ra.deriv.power_series",
		"ra.metric.contraction",
		"ra.measure.cantor_set",
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
