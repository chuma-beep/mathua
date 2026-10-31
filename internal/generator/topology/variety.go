package topology

import "github.com/chuma-beep/mathua/internal/generator"

// distinctDistractors ensures option candidates do not collide with the correct
// answer or with each other, providing fallbacks when collisions occur.
func distinctDistractors(correct string, candidates ...string) []string {
	seen := map[string]bool{correct: true}
	out := make([]string, 0, 3)
	for _, cand := range candidates {
		if cand != "" && !seen[cand] {
			seen[cand] = true
			out = append(out, cand)
			if len(out) == 3 {
				break
			}
		}
	}
	fallbacks := []string{
		"not enough information is provided",
		"only on compact subspaces",
		"none of the other choices holds",
		"holds only for finite sets",
		"diverges everywhere",
	}
	for _, fb := range fallbacks {
		if len(out) == 3 {
			break
		}
		if !seen[fb] {
			seen[fb] = true
			out = append(out, fb)
		}
	}
	return out
}

// ---------------------------------------------------------------------------
// 1. Hausdorff Spaces (topo.hausdorff)
// ---------------------------------------------------------------------------

type hausdorffGen struct{}

func (g *hausdorffGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return hausdorffProblem()
}

// ---------------------------------------------------------------------------
// 2. Quotient Spaces (topo.quotient)
// ---------------------------------------------------------------------------

type quotientGen struct{}

func (g *quotientGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return quotientProblem()
}

// ---------------------------------------------------------------------------
// 3. Product Spaces (topo.product)
// ---------------------------------------------------------------------------

type productGen struct{}

func (g *productGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return productProblem()
}

// ---------------------------------------------------------------------------
// 4. Urysohn Lemma (topo.urysohn)
// ---------------------------------------------------------------------------

type urysohnGen struct{}

func (g *urysohnGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return urysohnProblem()
}

// ---------------------------------------------------------------------------
// 5. Metrization Theorems (topo.metrization)
// ---------------------------------------------------------------------------

type metrizationGen struct{}

func (g *metrizationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return metrizationProblem()
}

// ---------------------------------------------------------------------------
// 6. Path-Connectedness (topo.connected.path)
// ---------------------------------------------------------------------------

type pathConnectedGen struct{}

func (g *pathConnectedGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return pathConnectedProblem()
}

// ---------------------------------------------------------------------------
// 7. Local Compactness (topo.compact.local)
// ---------------------------------------------------------------------------

type localCompactGen struct{}

func (g *localCompactGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return localCompactProblem()
}

// ---------------------------------------------------------------------------
// 8. Tietze Extension Theorem (topo.separation.tietze)
// ---------------------------------------------------------------------------

type tietzeGen struct{}

func (g *tietzeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return tietzeProblem()
}

// ---------------------------------------------------------------------------
// 9. Homotopy and Fundamental Group (topo.homotopy)
// ---------------------------------------------------------------------------

type homotopyGen struct{}

func (g *homotopyGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return homotopyProblem()
}

// ---------------------------------------------------------------------------
// 10. Topological Manifolds (topo.manifold)
// ---------------------------------------------------------------------------

type manifoldGen struct{}

func (g *manifoldGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return manifoldProblem()
}

// ---------------------------------------------------------------------------
// 11. Local Path-Connectedness (topo.connected.local_path)
// ---------------------------------------------------------------------------

type localPathGen struct{}

func (g *localPathGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return localPathProblem()
}

// ---------------------------------------------------------------------------
// 12. One-Point Compactification (topo.compact.one_point)
// ---------------------------------------------------------------------------

type onePointGen struct{}

func (g *onePointGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return onePointProblem()
}

// ---------------------------------------------------------------------------
// 13. Complete Regularity (topo.separation.completely_regular)
// ---------------------------------------------------------------------------

type completelyRegularGen struct{}

func (g *completelyRegularGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return completelyRegularProblem()
}
