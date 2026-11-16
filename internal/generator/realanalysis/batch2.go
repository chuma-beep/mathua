package realanalysis

// The second real-analysis batch: the rest of limits and continuity, the
// derivative and the integral made rigorous, uniform convergence, and the first
// look at metric spaces.
//
// The first batch built the real line from completeness and proved the theorems
// about sequences, series and the first two-thirds of continuity. This batch
// finishes continuity with connectedness, the classification of discontinuities
// and monotone functions; proves the mean value theorem, Taylor's theorem and
// the Riemann integral with the fundamental theorem, integration by parts and
// improper integrals; treats uniform convergence and the interchange of limits;
// and opens the metric-space picture with open and closed sets, closure,
// completeness, compactness, continuity, the contraction mapping principle,
// Baire's theorem and equivalent metrics.
//
// Ten concepts are numeric, so the definitions can be exercised, and fifteen are
// multiple choice, because the domain is largely about statements and proofs.
// Every numeric generator publishes the parameters its answer follows from, and
// the tests recompute the answer rather than trusting the author's arithmetic.

import (
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func raIntPow(base, exp int) int {
	out := 1
	for i := 0; i < exp; i++ {
		out *= base
	}
	return out
}

// The 14 multiple-choice generators (connectedGen through equivalentGen)
// are defined with elevated parameterization and variety in variety.go.

// ---------------------------------------------------------------------------
// The mean value theorem (numeric)
// ---------------------------------------------------------------------------

type mvtGen struct{}

func (g *mvtGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return raW2HMVT(raW2HOrdinal(ctx))
}

// ---------------------------------------------------------------------------
// The chain rule (numeric)
// ---------------------------------------------------------------------------

type chainGen struct{}

func (g *chainGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return raW2HChain(raW2HOrdinal(ctx))
}

// ---------------------------------------------------------------------------
// L'Hopital's rule (numeric)
// ---------------------------------------------------------------------------

type lhopitalGen struct{}

func (g *lhopitalGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return raW2HLHopital(raW2HOrdinal(ctx))
}

// ---------------------------------------------------------------------------
// The Riemann integral (numeric)
// ---------------------------------------------------------------------------

type riemannGen struct{}

func (g *riemannGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return raW2HRiemann(raW2HOrdinal(ctx))
}

// ---------------------------------------------------------------------------
// The fundamental theorem of calculus (numeric)
// ---------------------------------------------------------------------------

type ftcGen struct{}

func (g *ftcGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := 1 + rand.Intn(6)
	c := rand.Intn(6)
	b := 2 + rand.Intn(8)
	value := float64(a)*float64(b*b)/2 + float64(c*b)
	answer := decimal6(value)
	return generator.Problem{
		Question: "By the fundamental theorem of calculus, the integral from 0 to b of (a·t + c) with respect to t is a·b²/2 + c·b. For a = " +
			strconv.Itoa(a) + ", c = " + strconv.Itoa(c) + " and b = " + strconv.Itoa(b) + ", what is the value? (enter a decimal)",
		Answer:      answer,
		Explanation: "An antiderivative of a·t + c is a·t²/2 + c·t, and evaluating it at b and at 0 gives a·b²/2 + c·b = " + answer + ".",
		Facts: map[string]string{
			"a": strconv.Itoa(a), "c": strconv.Itoa(c), "b": strconv.Itoa(b), "answer": answer,
		},
	}
}

// ---------------------------------------------------------------------------
// Integration by parts (numeric)
// ---------------------------------------------------------------------------

type partsGen struct{}

func (g *partsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return raW2HParts(raW2HOrdinal(ctx))
}

// ---------------------------------------------------------------------------
// Improper integrals (numeric)
// ---------------------------------------------------------------------------

type improperGen struct{}

func (g *improperGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// Every exponent above one is convergent, so the pool is the range rather
	// than a short list of selected values.
	p := 2 + rand.Intn(38)
	answer := decimal6(1 / float64(p-1))
	return generator.Problem{
		Question: "The improper integral from 1 to infinity of x^(−p) converges exactly when p > 1, to 1/(p − 1). For p = " +
			strconv.Itoa(p) + ", what is its value? (enter a decimal)",
		Answer:      answer,
		Explanation: "An antiderivative is x^(1−p)/(1−p), which tends to 0 at infinity when p > 1, leaving 1/(p − 1) = " + answer + ".",
		Facts: map[string]string{
			"p": strconv.Itoa(p), "answer": answer,
		},
	}
}

// ---------------------------------------------------------------------------
// Power series (numeric)
// ---------------------------------------------------------------------------

type powerSeriesGen struct{}

func (g *powerSeriesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return raW2HPowerSeries(raW2HOrdinal(ctx))
}

// ---------------------------------------------------------------------------
// Taylor's theorem (numeric)
// ---------------------------------------------------------------------------

type taylorGen struct{}

func (g *taylorGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return raW2HTaylor(raW2HOrdinal(ctx))
}

// ---------------------------------------------------------------------------
// The contraction mapping principle (numeric)
// ---------------------------------------------------------------------------

type contractionGen struct{}

func (g *contractionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return raW2HContraction(raW2HOrdinal(ctx))
}
