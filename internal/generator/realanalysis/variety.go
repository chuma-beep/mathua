package realanalysis

import (
	"github.com/chuma-beep/mathua/internal/generator"
)

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
// 1. Connectedness (ra.limit.connected)
// ---------------------------------------------------------------------------

type connectedGen struct{}

func (g *connectedGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return connectedStatementsProblem()
}

// ---------------------------------------------------------------------------
// 2. Classification of discontinuities (ra.limit.discontinuity)
// ---------------------------------------------------------------------------

type discontinuityGen struct{}

func (g *discontinuityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return discontinuityStatementsProblem()
}

// ---------------------------------------------------------------------------
// 3. Monotone functions (ra.limit.monotone_fn)
// ---------------------------------------------------------------------------

type monotoneFnGen struct{}

func (g *monotoneFnGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return monotoneStatementsProblem()
}

// ---------------------------------------------------------------------------
// 4. Darboux's theorem (ra.deriv.darboux)
// ---------------------------------------------------------------------------

type darbouxGen struct{}

func (g *darbouxGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return darbouxStatementsProblem()
}

// ---------------------------------------------------------------------------
// 5. Uniform convergence (ra.deriv.uniform_conv)
// ---------------------------------------------------------------------------

type uniformConvGen struct{}

func (g *uniformConvGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return uniformConvStatementsProblem()
}

// ---------------------------------------------------------------------------
// 6. Interchange of limits (ra.deriv.interchange)
// ---------------------------------------------------------------------------

type interchangeGen struct{}

func (g *interchangeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return interchangeStatementsProblem()
}

// ---------------------------------------------------------------------------
// 7. Open and closed sets (ra.metric.open_closed)
// ---------------------------------------------------------------------------

type openClosedGen struct{}

func (g *openClosedGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return openClosedProblem()
}

// ---------------------------------------------------------------------------
// 8. Closure and boundary (ra.metric.closure)
// ---------------------------------------------------------------------------

type closureGen struct{}

func (g *closureGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return closureProblem()
}

// ---------------------------------------------------------------------------
// 9. Convergence in metric spaces (ra.metric.convergence)
// ---------------------------------------------------------------------------

type metricConvGen struct{}

func (g *metricConvGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return metricConvergenceProblem()
}

// ---------------------------------------------------------------------------
// 10. Complete metric spaces (ra.metric.complete)
// ---------------------------------------------------------------------------

type completeGen struct{}

func (g *completeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return completeProblem()
}

// ---------------------------------------------------------------------------
// 11. Compact metric spaces (ra.metric.compact)
// ---------------------------------------------------------------------------

type metricCompactGen struct{}

func (g *metricCompactGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return metricCompactStatementsProblem()
}

// ---------------------------------------------------------------------------
// 12. Continuity in metric spaces (ra.metric.continuity)
// ---------------------------------------------------------------------------

type metricContinuityGen struct{}

func (g *metricContinuityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return metricContinuityStatementsProblem()
}

// ---------------------------------------------------------------------------
// 13. Baire category theorem (ra.metric.baire)
// ---------------------------------------------------------------------------

type baireGen struct{}

func (g *baireGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return baireStatementsProblem()
}

// ---------------------------------------------------------------------------
// 14. Equivalent metrics (ra.metric.equivalent)
// ---------------------------------------------------------------------------

type equivalentGen struct{}

func (g *equivalentGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return equivalentStatementsProblem()
}
