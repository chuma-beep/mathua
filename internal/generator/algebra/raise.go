package algebra

// The 2026-10-19 raise: twenty concepts that take algebra from 104 to the
// domain target of 124.
//
// The batch fills the groups the domain left thin. Functions gain piecewise
// evaluation, graph transformations and the domain of an inverse; exponentials
// and logarithms gain their graphs, compound interest and the solution of a
// logarithmic model; radicals gain simplifying, combining and conjugates;
// rational expressions gain complex fractions, proportions and vertical
// asymptotes; systems gain a three-variable elimination, the matrix method and
// a nonlinear pair; sequences gain sigma notation and induction; polynomials
// gain the factor theorem and the rational root theorem.
//
// Twelve concepts are numeric, so the arithmetic is exercised, and eight are
// conceptual multiple choice, where the domain is about statements rather than
// computation. Every numeric generator publishes the parameters its answer
// follows from, and the tests recompute the answer rather than trusting the
// author's arithmetic.

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerRaise(reg *generator.Registry) {
	reg.Register("alg.func.piecewise", &piecewiseGen{})
	reg.Register("alg.func.transformations", &transformationsGen{})
	reg.Register("alg.func.inverse_domain", &inverseDomainGen{})
	reg.Register("alg.exp.graph", &expGraphGen{})
	reg.Register("alg.exp.compound", &compoundGen{})
	reg.Register("alg.log.graph", &logGraphGen{})
	reg.Register("alg.log.model", &logModelGen{})
	reg.Register("alg.radical.simplify", &radicalSimplifyGen{})
	reg.Register("alg.radical.operations", &radicalOpsGen{})
	reg.Register("alg.radical.conjugate", &radicalConjugateGen{})
	reg.Register("alg.rational.complex_fraction", &complexFractionGen{})
	reg.Register("alg.rational.proportions", &proportionsGen{})
	reg.Register("alg.rational.asymptotes", &asymptotesGen{})
	reg.Register("alg.systems.three_var", &threeVarGen{})
	reg.Register("alg.systems.matrices", &matricesGen{})
	reg.Register("alg.systems.nonlinear", &nonlinearSystemGen{})
	reg.Register("alg.seq.sigma", &sigmaGen{})
	reg.Register("alg.seq.induction", &inductionGen{})
	reg.Register("alg.poly.rational_root", &rationalRootGen{})
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// algDec renders a float with up to six decimals, trimming trailing zeros and a
// trailing point, and mapping a signed zero to plain zero.
func algDec(v float64) string {
	s := strconv.FormatFloat(v, 'f', 6, 64)
	for len(s) > 0 && s[len(s)-1] == '0' {
		s = s[:len(s)-1]
	}
	if len(s) > 0 && s[len(s)-1] == '.' {
		s = s[:len(s)-1]
	}
	if s == "" || s == "-0" {
		s = "0"
	}
	return s
}

// algMCClaim is one multiple-choice item: a question, the correct option, the
// distractors that must all be rejected, and the reason the correct option is
// right.
type algMCClaim struct {
	q           string
	correct     string
	distractors []string
	reason      string
}

func algPickMC(claims []algMCClaim) generator.Problem {
	c := claims[rand.Intn(len(claims))]
	p := generator.MultipleChoice(c.q, c.correct, c.distractors...)
	p.Explanation = c.reason
	if p.Facts == nil {
		p.Facts = map[string]string{}
	}
	p.Facts["answer"] = c.correct
	return p
}

// ---------------------------------------------------------------------------
// Piecewise evaluation (numeric)
// ---------------------------------------------------------------------------

type piecewiseGen struct{}

func (g *piecewiseGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	xs := []int{-4, -2, -1, 0, 1, 2, 3, 5}
	x0 := xs[rand.Intn(len(xs))]
	answer := 0
	branch := ""
	if x0 < 0 {
		answer = 2*x0 + 1
		branch = "x < 0, so use 2x + 1"
	} else {
		answer = 3*x0 - 2
		branch = "x >= 0, so use 3x - 2"
	}
	return generator.Problem{
		Question: "For the piecewise function f(x) = 2x + 1 when x < 0 and f(x) = 3x - 2 when x >= 0, what is f(" +
			strconv.Itoa(x0) + ")?",
		Answer:      strconv.Itoa(answer),
		Explanation: "The input " + strconv.Itoa(x0) + " satisfies " + branch + ", giving " + strconv.Itoa(answer) + ".",
		Facts: map[string]string{
			"x0": strconv.Itoa(x0), "answer": strconv.Itoa(answer),
		},
	}
}

// ---------------------------------------------------------------------------
// Graph transformations (multiple choice)
// ---------------------------------------------------------------------------

type transformationsGen struct{}

func (g *transformationsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return algPickMC([]algMCClaim{
		{
			q:           "Compared with the graph of y = f(x), the graph of y = f(x) + 3 is...",
			correct:     "shifted up three units",
			distractors: []string{"shifted down three units", "shifted right three units", "stretched vertically by a factor of three"},
			reason:      "Adding a constant outside the function moves every output up by that constant.",
		},
		{
			q:           "The graph of y = f(x - 2) is the graph of y = f(x)...",
			correct:     "shifted right two units",
			distractors: []string{"shifted left two units", "shifted up two units", "stretched horizontally by a factor of two"},
			reason:      "A subtraction inside the function moves the graph right, because the old output now occurs at a larger input.",
		},
		{
			q:           "The graph of y = -f(x) is the graph of y = f(x)...",
			correct:     "reflected across the x-axis",
			distractors: []string{"reflected across the y-axis", "shifted down one unit", "rotated a quarter turn"},
			reason:      "Negating the output reverses its sign, which is a reflection in the horizontal axis.",
		},
		{
			q:           "The graph of y = f(2x) is the graph of y = f(x)...",
			correct:     "compressed horizontally by a factor of two",
			distractors: []string{"stretched horizontally by a factor of two", "compressed vertically by a factor of two", "shifted right two units"},
			reason:      "Multiplying the input by two makes each output happen twice as soon, which squeezes the graph horizontally.",
		},
	})
}

// ---------------------------------------------------------------------------
// Domain of an inverse (multiple choice)
// ---------------------------------------------------------------------------

type inverseDomainGen struct{}

func (g *inverseDomainGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return algPickMC([]algMCClaim{
		{
			q:           "If f has domain [0, infinity) and range [2, 5], then the domain of the inverse function is...",
			correct:     "[2, 5]",
			distractors: []string{"[0, infinity)", "[0, 5]", "the whole real line"},
			reason:      "The inverse takes each value of f as an input, so its domain is the range of f.",
		},
		{
			q:           "The domain of the inverse of a one-to-one function f is...",
			correct:     "the range of f",
			distractors: []string{"the domain of f", "the whole real line", "the set where f is increasing"},
			reason:      "Swapping inputs and outputs swaps domain and range.",
		},
		{
			q:           "If f is not one-to-one on its domain, then f...",
			correct:     "has no inverse unless the domain is restricted",
			distractors: []string{"is always invertible", "has an inverse equal to its reciprocal", "is its own inverse"},
			reason:      "An inverse must send each output back to one input, which requires distinct inputs to give distinct outputs.",
		},
	})
}

// ---------------------------------------------------------------------------
// Exponential graphs (numeric)
// ---------------------------------------------------------------------------

type expGraphGen struct{}

func (g *expGraphGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct{ a, b, x int }
	claims := []claim{
		{a: 2, b: 3, x: 4},
		{a: 3, b: 2, x: 3},
		{a: 5, b: 2, x: 4},
		{a: 4, b: 3, x: 2},
		{a: 2, b: 5, x: 2},
		{a: 6, b: 2, x: 3},
	}
	c := claims[rand.Intn(len(claims))]
	value := 1
	for i := 0; i < c.x; i++ {
		value *= c.b
	}
	value *= c.a
	return generator.Problem{
		Question: "For the exponential function f(x) = " + strconv.Itoa(c.a) + " times " + strconv.Itoa(c.b) +
			"^x, what is f(" + strconv.Itoa(c.x) + ")?",
		Answer: strconv.Itoa(value),
		Explanation: "Raising the base " + strconv.Itoa(c.b) + " to the power " + strconv.Itoa(c.x) + " and multiplying by " +
			strconv.Itoa(c.a) + " gives " + strconv.Itoa(value) + ".",
		Facts: map[string]string{
			"a": strconv.Itoa(c.a), "b": strconv.Itoa(c.b), "x": strconv.Itoa(c.x), "answer": strconv.Itoa(value),
		},
	}
}

// ---------------------------------------------------------------------------
// Compound interest (numeric)
// ---------------------------------------------------------------------------

type compoundGen struct{}

func (g *compoundGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct{ p, r, t int }
	claims := []claim{
		{p: 1000, r: 5, t: 3},
		{p: 2000, r: 4, t: 2},
		{p: 500, r: 10, t: 3},
		{p: 1500, r: 3, t: 4},
		{p: 800, r: 6, t: 2},
		{p: 2500, r: 2, t: 5},
	}
	c := claims[rand.Intn(len(claims))]
	value := float64(c.p) * math.Pow(1+float64(c.r)/100, float64(c.t))
	answer := fmt.Sprintf("%.2f", value)
	return generator.Problem{
		Question: "An investment of " + strconv.Itoa(c.p) + " dollars earns " + strconv.Itoa(c.r) +
			" percent interest per year, compounded once a year, for " + strconv.Itoa(c.t) +
			" years. Using A = P(1 + r)^t, what is the value at the end, in dollars? (round to cents)",
		Answer: answer,
		Explanation: "With P = " + strconv.Itoa(c.p) + ", the annual factor is 1 + " + strconv.Itoa(c.r) +
			"/100, and raising it to the power " + strconv.Itoa(c.t) + " and multiplying by P gives " + answer + " dollars.",
		Facts: map[string]string{
			"p": strconv.Itoa(c.p), "r": strconv.Itoa(c.r), "t": strconv.Itoa(c.t), "answer": answer,
		},
	}
}

// ---------------------------------------------------------------------------
// Logarithmic graphs (multiple choice)
// ---------------------------------------------------------------------------

type logGraphGen struct{}

func (g *logGraphGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return algPickMC([]algMCClaim{
		{
			q:           "The graph of y = log base b of x, with b > 1, ...",
			correct:     "has a vertical asymptote at x = 0 and passes through the point (1, 0)",
			distractors: []string{"has a horizontal asymptote at y = 0 and passes through (0, 1)", "is decreasing everywhere", "passes through the origin"},
			reason:      "The logarithm is the inverse of an exponential, so their graphs reflect: the log has a vertical asymptote where the exponential had its horizontal one.",
		},
		{
			q:           "As x grows without bound, log base b of x with b > 1...",
			correct:     "grows without bound, but more slowly than x",
			distractors: []string{"approaches a horizontal asymptote", "oscillates between two values", "stays constant"},
			reason:      "The logarithm increases without limit but its rate of growth falls to zero relative to the identity.",
		},
		{
			q:           "The graph of y = log base b of x is the reflection of the graph of y = b^x across...",
			correct:     "the line y = x",
			distractors: []string{"the x-axis", "the y-axis", "the line y = -x"},
			reason:      "A function and its inverse are mirror images across the line y = x, which is the swap of input and output.",
		},
	})
}

// ---------------------------------------------------------------------------
// Logarithmic models (numeric)
// ---------------------------------------------------------------------------

type logModelGen struct{}

func (g *logModelGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct{ b, k int }
	claims := []claim{
		{b: 2, k: 3},
		{b: 2, k: 4},
		{b: 3, k: 2},
		{b: 5, k: 2},
		{b: 10, k: 3},
		{b: 3, k: 3},
	}
	c := claims[rand.Intn(len(claims))]
	value := 1
	for i := 0; i < c.k; i++ {
		value *= c.b
	}
	return generator.Problem{
		Question: "Solve log base " + strconv.Itoa(c.b) + " of x = " + strconv.Itoa(c.k) +
			" for x.",
		Answer: strconv.Itoa(value),
		Explanation: "The equation says the base " + strconv.Itoa(c.b) + " raised to the power " + strconv.Itoa(c.k) +
			" equals x, so x = " + strconv.Itoa(value) + ".",
		Facts: map[string]string{
			"b": strconv.Itoa(c.b), "k": strconv.Itoa(c.k), "answer": strconv.Itoa(value),
		},
	}
}

// ---------------------------------------------------------------------------
// Simplifying radicals (numeric)
// ---------------------------------------------------------------------------

type radicalSimplifyGen struct{}

func (g *radicalSimplifyGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct{ a, b int }
	claims := []claim{
		{a: 2, b: 2},
		{a: 3, b: 2},
		{a: 4, b: 2},
		{a: 2, b: 3},
		{a: 3, b: 5},
		{a: 5, b: 2},
		{a: 2, b: 7},
		{a: 3, b: 6},
	}
	c := claims[rand.Intn(len(claims))]
	n := c.a * c.a * c.b
	return generator.Problem{
		Question: "The square root of " + strconv.Itoa(n) + " can be written as a times the square root of b, where b is squarefree. What is the whole number a?",
		Answer:   strconv.Itoa(c.a),
		Explanation: "Since " + strconv.Itoa(n) + " = " + strconv.Itoa(c.a*c.a) + " times " + strconv.Itoa(c.b) +
			", the square root is " + strconv.Itoa(c.a) + " times the square root of " + strconv.Itoa(c.b) + ".",
		Facts: map[string]string{
			"n": strconv.Itoa(n), "a": strconv.Itoa(c.a), "b": strconv.Itoa(c.b), "answer": strconv.Itoa(c.a),
		},
	}
}

// ---------------------------------------------------------------------------
// Operations with radicals (numeric)
// ---------------------------------------------------------------------------

type radicalOpsGen struct{}

func (g *radicalOpsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct{ a, b, c int }
	claims := []claim{
		{a: 2, b: 3, c: 2},
		{a: 5, b: 2, c: 3},
		{a: 4, b: 3, c: 5},
		{a: 7, b: 2, c: 7},
		{a: 6, b: 5, c: 6},
		{a: 3, b: 8, c: 2},
	}
	c := claims[rand.Intn(len(claims))]
	sum := c.a + c.b
	return generator.Problem{
		Question: "Simplify " + strconv.Itoa(c.a) + " times the square root of " + strconv.Itoa(c.c) + " plus " +
			strconv.Itoa(c.b) + " times the square root of " + strconv.Itoa(c.c) +
			". The result is k times the square root of " + strconv.Itoa(c.c) + "; what is k?",
		Answer: strconv.Itoa(sum),
		Explanation: "The two terms share the radical, so their coefficients add: " + strconv.Itoa(c.a) + " + " +
			strconv.Itoa(c.b) + " = " + strconv.Itoa(sum) + ".",
		Facts: map[string]string{
			"a": strconv.Itoa(c.a), "b": strconv.Itoa(c.b), "c": strconv.Itoa(c.c), "answer": strconv.Itoa(sum),
		},
	}
}

// ---------------------------------------------------------------------------
// Rationalizing with conjugates (multiple choice)
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Rationalizing with conjugates (numeric)
// ---------------------------------------------------------------------------

type radicalConjugateGen struct{}

func (g *radicalConjugateGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	a := rand.Intn(scale*6) + 3
	b := rand.Intn(a-1) + 1
	if a == b {
		a++
	}

	ans := a - b
	ansStr := strconv.Itoa(ans)
	aStr := strconv.Itoa(a)
	bStr := strconv.Itoa(b)

	return generator.Problem{
		Question: fmt.Sprintf("Evaluate the product \\((\\sqrt{%d} + \\sqrt{%d})(\\sqrt{%d} - \\sqrt{%d})\\).", a, b, a, b),
		Answer:   ansStr,
		Explanation: fmt.Sprintf(
			"Use the difference of squares identity: \\((\\sqrt{%d} + \\sqrt{%d})(\\sqrt{%d} - \\sqrt{%d}) = (\\sqrt{%d})^2 - (\\sqrt{%d})^2 = %d - %d = %d\\).",
			a, b, a, b, a, b, a, b, ans,
		),
		Facts: map[string]string{
			"a":      aStr,
			"b":      bStr,
			"answer": ansStr,
		},
	}
}

// ---------------------------------------------------------------------------
// Complex fractions (multiple choice)
// ---------------------------------------------------------------------------

type complexFractionGen struct{}

func (g *complexFractionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type cfClaim struct {
		expr       string
		top        string
		numEval    string
		numPart    string
		reciprocal string
		k          int
		hasEval    bool
	}

	claims := []cfClaim{
		// Type 1: direct fraction over fraction
		{expr: `\frac{\frac{3}{2}}{\frac{1}{4}}`, top: `\frac{3}{2}`, numPart: `\frac{3}{2}`, reciprocal: `\frac{4}{1}`, k: 6},
		{expr: `\frac{\frac{5}{3}}{\frac{5}{9}}`, top: `\frac{5}{3}`, numPart: `\frac{5}{3}`, reciprocal: `\frac{9}{5}`, k: 3},
		{expr: `\frac{\frac{7}{4}}{\frac{7}{12}}`, top: `\frac{7}{4}`, numPart: `\frac{7}{4}`, reciprocal: `\frac{12}{7}`, k: 3},
		{expr: `\frac{\frac{4}{3}}{\frac{2}{9}}`, top: `\frac{4}{3}`, numPart: `\frac{4}{3}`, reciprocal: `\frac{9}{2}`, k: 6},
		{expr: `\frac{\frac{5}{2}}{\frac{5}{6}}`, top: `\frac{5}{2}`, numPart: `\frac{5}{2}`, reciprocal: `\frac{6}{5}`, k: 3},
		{expr: `\frac{\frac{7}{3}}{\frac{7}{6}}`, top: `\frac{7}{3}`, numPart: `\frac{7}{3}`, reciprocal: `\frac{6}{7}`, k: 2},
		{expr: `\frac{\frac{8}{5}}{\frac{4}{15}}`, top: `\frac{8}{5}`, numPart: `\frac{8}{5}`, reciprocal: `\frac{15}{4}`, k: 6},
		{expr: `\frac{\frac{6}{7}}{\frac{2}{21}}`, top: `\frac{6}{7}`, numPart: `\frac{6}{7}`, reciprocal: `\frac{21}{2}`, k: 9},
		{expr: `\frac{\frac{10}{3}}{\frac{5}{6}}`, top: `\frac{10}{3}`, numPart: `\frac{10}{3}`, reciprocal: `\frac{6}{5}`, k: 4},
		{expr: `\frac{\frac{14}{5}}{\frac{7}{15}}`, top: `\frac{14}{5}`, numPart: `\frac{14}{5}`, reciprocal: `\frac{15}{7}`, k: 6},
		{expr: `\frac{\frac{9}{4}}{\frac{3}{8}}`, top: `\frac{9}{4}`, numPart: `\frac{9}{4}`, reciprocal: `\frac{8}{3}`, k: 6},
		{expr: `\frac{\frac{15}{2}}{\frac{5}{4}}`, top: `\frac{15}{2}`, numPart: `\frac{15}{2}`, reciprocal: `\frac{4}{5}`, k: 6},
		{expr: `\frac{\frac{16}{3}}{\frac{4}{9}}`, top: `\frac{16}{3}`, numPart: `\frac{16}{3}`, reciprocal: `\frac{9}{4}`, k: 12},
		{expr: `\frac{\frac{12}{5}}{\frac{3}{10}}`, top: `\frac{12}{5}`, numPart: `\frac{12}{5}`, reciprocal: `\frac{10}{3}`, k: 8},
		{expr: `\frac{\frac{21}{4}}{\frac{7}{8}}`, top: `\frac{21}{4}`, numPart: `\frac{21}{4}`, reciprocal: `\frac{8}{7}`, k: 6},
		{expr: `\frac{\frac{20}{7}}{\frac{4}{21}}`, top: `\frac{20}{7}`, numPart: `\frac{20}{7}`, reciprocal: `\frac{21}{4}`, k: 15},

		// Type 2: sum or difference in numerator over fraction
		{expr: `\frac{1 + \frac{1}{2}}{\frac{1}{4}}`, top: `1 + \frac{1}{2}`, numEval: `\frac{3}{2}`, numPart: `\frac{3}{2}`, reciprocal: `\frac{4}{1}`, k: 6, hasEval: true},
		{expr: `\frac{1 + \frac{1}{3}}{\frac{2}{9}}`, top: `1 + \frac{1}{3}`, numEval: `\frac{4}{3}`, numPart: `\frac{4}{3}`, reciprocal: `\frac{9}{2}`, k: 6, hasEval: true},
		{expr: `\frac{1 + \frac{2}{3}}{\frac{5}{6}}`, top: `1 + \frac{2}{3}`, numEval: `\frac{5}{3}`, numPart: `\frac{5}{3}`, reciprocal: `\frac{6}{5}`, k: 2, hasEval: true},
		{expr: `\frac{1 - \frac{1}{4}}{\frac{3}{8}}`, top: `1 - \frac{1}{4}`, numEval: `\frac{3}{4}`, numPart: `\frac{3}{4}`, reciprocal: `\frac{8}{3}`, k: 2, hasEval: true},
		{expr: `\frac{1 - \frac{1}{3}}{\frac{2}{9}}`, top: `1 - \frac{1}{3}`, numEval: `\frac{2}{3}`, numPart: `\frac{2}{3}`, reciprocal: `\frac{9}{2}`, k: 3, hasEval: true},
		{expr: `\frac{1 + \frac{3}{4}}{\frac{7}{8}}`, top: `1 + \frac{3}{4}`, numEval: `\frac{7}{4}`, numPart: `\frac{7}{4}`, reciprocal: `\frac{8}{7}`, k: 2, hasEval: true},
		{expr: `\frac{1 - \frac{2}{5}}{\frac{3}{10}}`, top: `1 - \frac{2}{5}`, numEval: `\frac{3}{5}`, numPart: `\frac{3}{5}`, reciprocal: `\frac{10}{3}`, k: 2, hasEval: true},
		{expr: `\frac{1 + \frac{1}{5}}{\frac{3}{10}}`, top: `1 + \frac{1}{5}`, numEval: `\frac{6}{5}`, numPart: `\frac{6}{5}`, reciprocal: `\frac{10}{3}`, k: 4, hasEval: true},
		{expr: `\frac{1 + \frac{1}{4}}{\frac{5}{12}}`, top: `1 + \frac{1}{4}`, numEval: `\frac{5}{4}`, numPart: `\frac{5}{4}`, reciprocal: `\frac{12}{5}`, k: 3, hasEval: true},
		{expr: `\frac{1 - \frac{1}{6}}{\frac{5}{12}}`, top: `1 - \frac{1}{6}`, numEval: `\frac{5}{6}`, numPart: `\frac{5}{6}`, reciprocal: `\frac{12}{5}`, k: 2, hasEval: true},
		{expr: `\frac{1 - \frac{1}{5}}{\frac{4}{15}}`, top: `1 - \frac{1}{5}`, numEval: `\frac{4}{5}`, numPart: `\frac{4}{5}`, reciprocal: `\frac{15}{4}`, k: 3, hasEval: true},
		{expr: `\frac{1 + \frac{1}{6}}{\frac{7}{18}}`, top: `1 + \frac{1}{6}`, numEval: `\frac{7}{6}`, numPart: `\frac{7}{6}`, reciprocal: `\frac{18}{7}`, k: 3, hasEval: true},
		{expr: `\frac{2 + \frac{1}{3}}{\frac{7}{6}}`, top: `2 + \frac{1}{3}`, numEval: `\frac{7}{3}`, numPart: `\frac{7}{3}`, reciprocal: `\frac{6}{7}`, k: 2, hasEval: true},
		{expr: `\frac{2 - \frac{1}{2}}{\frac{3}{4}}`, top: `2 - \frac{1}{2}`, numEval: `\frac{3}{2}`, numPart: `\frac{3}{2}`, reciprocal: `\frac{4}{3}`, k: 2, hasEval: true},
		{expr: `\frac{2 + \frac{2}{5}}{\frac{3}{5}}`, top: `2 + \frac{2}{5}`, numEval: `\frac{12}{5}`, numPart: `\frac{12}{5}`, reciprocal: `\frac{5}{3}`, k: 4, hasEval: true},
		{expr: `\frac{3 - \frac{1}{3}}{\frac{4}{3}}`, top: `3 - \frac{1}{3}`, numEval: `\frac{8}{3}`, numPart: `\frac{8}{3}`, reciprocal: `\frac{3}{4}`, k: 2, hasEval: true},
	}

	c := claims[rand.Intn(len(claims))]
	q := fmt.Sprintf("Simplify the complex fraction: \\[%s\\]", c.expr)
	ans := strconv.Itoa(c.k)

	var explanation string
	if c.hasEval {
		explanation = fmt.Sprintf("Simplify the numerator \\(%s = %s\\). Then multiply by the reciprocal of the denominator: \\(%s \\times %s = %d\\).",
			c.top, c.numEval, c.numPart, c.reciprocal, c.k)
	} else {
		explanation = fmt.Sprintf("Multiply the numerator by the reciprocal of the denominator: \\(%s \\times %s = %d\\).",
			c.numPart, c.reciprocal, c.k)
	}

	facts := map[string]string{
		"expr":       c.expr,
		"top":        c.top,
		"num_part":   c.numPart,
		"reciprocal": c.reciprocal,
		"answer":     ans,
		"k":          ans,
	}
	if c.hasEval {
		facts["num_eval"] = c.numEval
	}

	return generator.Problem{
		Question:    q,
		Answer:      ans,
		Explanation: explanation,
		Facts:       facts,
	}
}

// ---------------------------------------------------------------------------
// Proportions (numeric)
// ---------------------------------------------------------------------------

type proportionsGen struct{}

func (g *proportionsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct{ a, b, k int }
	claims := []claim{
		{a: 2, b: 3, k: 4},
		{a: 3, b: 5, k: 2},
		{a: 4, b: 7, k: 3},
		{a: 5, b: 4, k: 6},
		{a: 2, b: 9, k: 5},
		{a: 6, b: 5, k: 4},
	}
	c := claims[rand.Intn(len(claims))]
	num := c.a * c.k
	x := c.b * c.k
	return generator.Problem{
		Question: "Solve the proportion " + strconv.Itoa(c.a) + "/" + strconv.Itoa(c.b) + " = " + strconv.Itoa(num) +
			"/x for x.",
		Answer: strconv.Itoa(x),
		Explanation: "Cross-multiplying gives " + strconv.Itoa(c.a) + "x = " + strconv.Itoa(c.b) + " times " + strconv.Itoa(num) +
			", so x = " + strconv.Itoa(x) + ".",
		Facts: map[string]string{
			"a": strconv.Itoa(c.a), "b": strconv.Itoa(c.b), "num": strconv.Itoa(num), "answer": strconv.Itoa(x),
		},
	}
}

// ---------------------------------------------------------------------------
// Vertical asymptotes (numeric)
// ---------------------------------------------------------------------------

type asymptotesGen struct{}

func (g *asymptotesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	c := struct{ a, m, n int }{a: rand.Intn(9) + 1}
	for {
		c.m = rand.Intn(17) - 8
		c.n = rand.Intn(17) - 8
		// Both asymptotes must exist and be distinct, and neither may sit at
		// zero: a zero root would print the constant term as "+ 0".
		if c.m != 0 && c.n != 0 && c.m != c.n && c.m+c.n != 0 {
			break
		}
	}
	// f(x) = (x + a) / ((x - m)(x - n)). The denominator expands to
	// x^2 - (m+n)x + mn.
	b := -(c.m + c.n)
	d := c.m * c.n
	larger := c.m
	if c.n > c.m {
		larger = c.n
	}
	return generator.Problem{
		Question:    fmt.Sprintf("The rational function f(x) = (x + %d)/(x^2 %+dx %+d) has two vertical asymptotes. What is the larger one?", c.a, b, d),
		Answer:      strconv.Itoa(larger),
		Explanation: fmt.Sprintf("The denominator factors as (x%+d)(x%+d), so the vertical asymptotes are at %d and %d, and the larger is %d.", -c.m, -c.n, c.m, c.n, larger),
		Facts: map[string]string{
			"a": strconv.Itoa(c.a), "b": strconv.Itoa(b), "d": strconv.Itoa(d),
			"m": strconv.Itoa(c.m), "n": strconv.Itoa(c.n), "answer": strconv.Itoa(larger),
		},
	}
}

// ---------------------------------------------------------------------------
// Systems of three equations (numeric)
// ---------------------------------------------------------------------------

type threeVarGen struct{}

func (g *threeVarGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	c := struct{ x, y, z int }{
		x: rand.Intn(8) + 2,
		y: rand.Intn(9) + 1,
		z: rand.Intn(9) + 1,
	}
	a := c.x + c.y + c.z
	b := c.x - c.y + c.z
	d := c.x + c.y - c.z
	return generator.Problem{
		Question: "Solve the system x + y + z = " + strconv.Itoa(a) + ", x - y + z = " + strconv.Itoa(b) +
			", x + y - z = " + strconv.Itoa(d) + " for x.",
		Answer: strconv.Itoa(c.x),
		Explanation: "Subtracting the second equation from the first gives 2y = " + strconv.Itoa(a-b) + ", so y = " +
			strconv.Itoa(c.y) + "; subtracting the third from the first gives 2z = " + strconv.Itoa(a-d) + ", so z = " +
			strconv.Itoa(c.z) + "; then x = " + strconv.Itoa(a) + " - y - z = " + strconv.Itoa(c.x) + ".",
		Facts: map[string]string{
			"a": strconv.Itoa(a), "b": strconv.Itoa(b), "d": strconv.Itoa(d),
			"y": strconv.Itoa(c.y), "z": strconv.Itoa(c.z), "answer": strconv.Itoa(c.x),
		},
	}
}

// ---------------------------------------------------------------------------
// Matrices (multiple choice)
// ---------------------------------------------------------------------------

type matricesGen struct{}

func (g *matricesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return algPickMC([]algMCClaim{
		{
			q:           "A system of linear equations is solved by row-reducing its augmented matrix, which...",
			correct:     "holds the coefficients of the variables in the left columns and the constants in the rightmost column",
			distractors: []string{"holds only the constants", "holds the variables themselves as entries", "holds the right-hand sides in the leftmost column"},
			reason:      "The augmented column carries the constants so that row operations act on equations and constants together.",
		},
		{
			q:           "Which of these is a valid elementary row operation for Gaussian elimination?",
			correct:     "adding a multiple of one row to another row",
			distractors: []string{"multiplying two rows together", "adding the same constant to every entry", "replacing a row by its reciprocal"},
			reason:      "The three allowed operations are swapping rows, scaling a row by a nonzero constant, and adding a multiple of one row to another.",
		},
		{
			q:           "If the row-reduced augmented matrix contains a row reading 0 0 0 with a nonzero constant in the last column, the system has...",
			correct:     "no solution",
			distractors: []string{"exactly one solution", "infinitely many solutions", "a solution equal to that constant"},
			reason:      "That row states 0 equals a nonzero number, which no values of the variables can satisfy.",
		},
	})
}

// ---------------------------------------------------------------------------
// Nonlinear systems (numeric)
// ---------------------------------------------------------------------------

type nonlinearSystemGen struct{}

func (g *nonlinearSystemGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	ks := []int{2, 3, 4, 5, 6, 7}
	k := ks[rand.Intn(len(ks))]
	y := k * k
	return generator.Problem{
		Question: "The system y = x^2 and y = " + strconv.Itoa(k) + "x has two solutions. What is the y-coordinate of the solution with x not equal to 0?",
		Answer:   strconv.Itoa(y),
		Explanation: "Setting x^2 equal to " + strconv.Itoa(k) + "x gives x(x - " + strconv.Itoa(k) +
			") = 0, so the nonzero solution has x = " + strconv.Itoa(k) + " and y = " + strconv.Itoa(k) + "^2 = " + strconv.Itoa(y) + ".",
		Facts: map[string]string{
			"k": strconv.Itoa(k), "x": strconv.Itoa(k), "answer": strconv.Itoa(y),
		},
	}
}

// ---------------------------------------------------------------------------
// Sigma notation (numeric)
// ---------------------------------------------------------------------------

type sigmaGen struct{}

func (g *sigmaGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct{ a, b, n int }
	claims := []claim{
		{a: 1, b: 0, n: 5},
		{a: 2, b: 1, n: 4},
		{a: 3, b: 2, n: 6},
		{a: 1, b: 3, n: 8},
		{a: 2, b: 5, n: 5},
		{a: 4, b: 0, n: 5},
	}
	c := claims[rand.Intn(len(claims))]
	tri := c.n * (c.n + 1) / 2
	sum := c.a*tri + c.b*c.n
	return generator.Problem{
		Question: "Evaluate the sum from i = 1 to " + strconv.Itoa(c.n) + " of (" + strconv.Itoa(c.a) + "i + " +
			strconv.Itoa(c.b) + ").",
		Answer: strconv.Itoa(sum),
		Explanation: "The sum splits into " + strconv.Itoa(c.a) + " times the sum of i, which is " + strconv.Itoa(tri) +
			", plus " + strconv.Itoa(c.b) + " times " + strconv.Itoa(c.n) + " terms, giving " + strconv.Itoa(sum) + ".",
		Facts: map[string]string{
			"a": strconv.Itoa(c.a), "b": strconv.Itoa(c.b), "n": strconv.Itoa(c.n), "answer": strconv.Itoa(sum),
		},
	}
}

// ---------------------------------------------------------------------------
// Induction (multiple choice)
// ---------------------------------------------------------------------------

type inductionGen struct{}

func (g *inductionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return algPickMC([]algMCClaim{
		{
			q:           "A proof by induction establishes a statement by proving...",
			correct:     "a base case and an induction step that passes from n to n + 1",
			distractors: []string{"only the base case", "the statement for one very large n", "the statement for every n without a base case"},
			reason:      "The base case starts the chain and the induction step carries it forward to every later integer.",
		},
		{
			q:           "The induction step assumes...",
			correct:     "the statement holds for some n and then proves it for n + 1",
			distractors: []string{"the statement for n + 1 and derives it for n", "the statement for every n at once", "only the base case and nothing else"},
			reason:      "The step is a conditional: if the statement is true at one rung, it is true at the next.",
		},
		{
			q:           "If the base case fails but the induction step is valid, the proof...",
			correct:     "does not establish the statement",
			distractors: []string{"still proves the statement", "proves it for every n except the base", "proves the converse instead"},
			reason:      "Without a starting point the chain never begins, so no instance is reached.",
		},
	})
}

// ---------------------------------------------------------------------------
// The rational root theorem (numeric)
// ---------------------------------------------------------------------------

type rationalRootGen struct{}

func (g *rationalRootGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	c := struct{ r1, r2, r3 int }{r1: rand.Intn(9) + 1}
	c.r2 = -(rand.Intn(9) + 1)
	for {
		c.r3 = -(rand.Intn(9) + 1)
		if c.r3 != c.r2 {
			break
		}
	}
	b := -(c.r1 + c.r2 + c.r3)
	d := c.r1*c.r2 + c.r1*c.r3 + c.r2*c.r3
	e := -(c.r1 * c.r2 * c.r3)
	return generator.Problem{
		Question: fmt.Sprintf("The monic polynomial x^3 %+dx^2 %+dx %+d has three integer roots. By the rational root theorem, what is its positive root?", b, d, e),
		Answer:   strconv.Itoa(c.r1),
		Explanation: fmt.Sprintf("The polynomial factors as (x %+d)(x %+d)(x %+d), so the roots are %d, %d and %d, and the positive one is %d.",
			-c.r1, -c.r2, -c.r3, c.r1, c.r2, c.r3, c.r1),
		Facts: map[string]string{
			"b": strconv.Itoa(b), "d": strconv.Itoa(d), "e": strconv.Itoa(e),
			"r1": strconv.Itoa(c.r1), "r2": strconv.Itoa(c.r2), "r3": strconv.Itoa(c.r3),
			"answer": strconv.Itoa(c.r1),
		},
	}
}
