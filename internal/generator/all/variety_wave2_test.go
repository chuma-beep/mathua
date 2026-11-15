package all

import (
	"math/rand"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

// Wave 2a of the variety pass rewrote the graph-theory generators whose three
// knowledge points each drew from one fixed pool of six written questions. A
// learner who met such a concept met six questions in total, whatever ordinal
// the lesson said they were working on, and the seventh question repeated the
// first.
//
// Every one of these concepts now draws the graph, the matrix size or the
// numeric parameter, and branches on KPOrdinal, so one ordinal alone supplies at
// least ten distinct problems. The generic shipping contract only asks for five
// in twelve draws, which the old pools passed, so this guard holds the wider
// branches at ten over 200 draws.
func TestWave2GraphVariety(t *testing.T) {
	concepts := []string{
		"graph.adjacency",
		"graph.chromatic_polynomial",
		"graph.connectivity",
		"graph.degeneracy",
		"graph.edge_coloring",
		"graph.genus",
		"graph.hamiltonian",
		"graph.independent_set",
		"graph.laplacian",
		"graph.list_coloring",
		"graph.matrix_tree",
		"graph.max_flow",
		"graph.moore_bound",
		"graph.multiway_cut",
		"graph.outerplanar",
		"graph.ramsey",
		"graph.spectral_gap",
		"graph.spectral_radius",
		"graph.tournament",
		"graph.treewidth",
		"graph.turan",
		"graph.vertex_cover",
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
		}
	}
}
