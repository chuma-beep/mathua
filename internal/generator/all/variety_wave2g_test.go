package all

import (
	"math/rand"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

// Wave 2g of the variety pass covers fourteen discrete-mathematics concepts
// that answered at least one knowledge point from a fixed claim table, a
// nine-value parameter pair, or a generator that ignored the knowledge-point
// ordinal. A learner reviewing such a knowledge point met the same handful of
// questions, and a generator that ignored the ordinal could not serve the
// section its lesson actually taught.
//
// Every generator below now branches on KPOrdinal and draws its parameters, so
// one ordinal alone supplies at least ten distinct problems. The generic
// shipping contract only asks for five in twelve draws, which the old tables
// passed, so this guard holds the wider branches at ten over 200 draws and
// checks that the ordinal is recorded on every returned problem.
func TestWave2GDiscreteVariety(t *testing.T) {
	concepts := []string{
		"discrete.logic.boolean_functions",
		"discrete.logic.resolution",
		"discrete.proof.cases",
		"discrete.combinatorics.burnside",
		"discrete.combinatorics.pascal",
		"discrete.combinatorics.surjections",
		"discrete.combinatorics.double_counting",
		"discrete.combinatorics.integer_partitions",
		"discrete.combinatorics.permutations",
		"discrete.recurrence.master_theorem",
		"discrete.sequences.recurrence",
		"discrete.coding.entropy",
		"discrete.combinatorics.complement",
		"discrete.graphs.degree_sequence",
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
