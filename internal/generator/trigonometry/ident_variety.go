package trigonometry

import (
	"github.com/chuma-beep/mathua/internal/generator"
)

// distinctIdentDistractors returns three options that differ from the answer and from
// each other, falling back to diverse expressions when collisions occur.
func distinctIdentDistractors(correct string, candidates ...string) []string {
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
		"0",
		"1",
		"-1",
		"undefined",
		"no real solution",
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
// 1. Product to sum (trig.ident.product_to_sum)
// ---------------------------------------------------------------------------

type productToSumGen struct{}

func (g *productToSumGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return productToSumProblem()
}

// ---------------------------------------------------------------------------
// 2. Sum to product (trig.ident.sum_to_product)
// ---------------------------------------------------------------------------

type sumToProductGen struct{}

func (g *sumToProductGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return sumToProductProblem()
}

// ---------------------------------------------------------------------------
// 3. Pythagorean variants (trig.ident.pythagorean_variants)
// ---------------------------------------------------------------------------

type pythagoreanVariantsGen struct{}

func (g *pythagoreanVariantsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return pythagoreanVariantProblem()
}

// ---------------------------------------------------------------------------
// 4. General solutions (trig.eq.general_solution)
// ---------------------------------------------------------------------------

type generalSolutionGen struct{}

func (g *generalSolutionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return generalSolutionProblem()
}

// ---------------------------------------------------------------------------
// 5. Advanced inequalities (trig.ineq.advanced)
// ---------------------------------------------------------------------------

type advancedIneqGen struct{}

func (g *advancedIneqGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return advancedInequalityProblem()
}

// ---------------------------------------------------------------------------
// 6. Hyperbolic identities (trig.hyperbolic.identities)
// ---------------------------------------------------------------------------

type hyperbolicIdentGen struct{}

func (g *hyperbolicIdentGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return hyperbolicIdentityProblem()
}

// ---------------------------------------------------------------------------
// 7. Hyperbolic derivatives (trig.hyperbolic.derivative)
// ---------------------------------------------------------------------------

type hyperbolicDerivGen struct{}

func (g *hyperbolicDerivGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return hyperbolicDerivativeProblem()
}
