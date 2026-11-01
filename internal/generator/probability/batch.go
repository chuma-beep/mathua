package probability

// The first probability-theory batch: twenty-five concepts that open the domain
// with the axioms, conditional probability, random variables, a few distribution
// results and the classical limit inequalities.
//
// The statistics domain already has sample spaces, basic and compound
// probability, conditional probability, Bayes' theorem, expected value, discrete
// and continuous random variables and the named distributions. This domain
// starts after those and goes deeper: it treats a probability as a measure on a
// sigma-algebra, defines a random variable as a measurable function, works with
// mass and density functions, moments and generating functions, and proves the
// first limit results.
//
// Nine concepts are multiple choice and sixteen are numeric. Every numeric
// generator publishes the parameters its answer follows from, so the tests
// recompute the answer rather than trusting the author's arithmetic.

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/mathutil"
)

// decimal6 formats a probability or other decimal to at most six places, with
// trailing zeros removed, so 0.6400000000000001 becomes "0.64" and 7.0 becomes
// "7". The tests use the same formatter, which keeps the published answer and the
// recomputed one comparable without a tolerance.
func decimal6(v float64) string {
	s := strconv.FormatFloat(v, 'f', 6, 64)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	if s == "" || s == "-0" {
		s = "0"
	}
	return s
}

// ---------------------------------------------------------------------------
// Probability axioms (multiple choice)
// ---------------------------------------------------------------------------

type axiomsGen struct{}

func (g *axiomsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return probabilityAxiomsProblem()
}

// ---------------------------------------------------------------------------
// Sigma-algebras and events (multiple choice)
// ---------------------------------------------------------------------------

type sigmaAlgebraGen struct{}

func (g *sigmaAlgebraGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return sigmaAlgebraProblem()
}

// ---------------------------------------------------------------------------
// Inclusion-exclusion (numeric)
// ---------------------------------------------------------------------------

type inclusionExclusionGen struct{}

var incExcDenomsLow = []int{20, 25, 40, 50}
var incExcDenomsMid = []int{50, 80, 100, 125}
var incExcDenomsHigh = []int{100, 160, 200, 250, 400, 500}

func (g *inclusionExclusionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	tier := mathutil.DifficultyTier(ctx.Difficulty)
	var denoms []int
	switch tier {
	case 0:
		denoms = incExcDenomsLow
	case 1:
		denoms = incExcDenomsMid
	default:
		denoms = incExcDenomsHigh
	}
	n := denoms[rand.Intn(len(denoms))]

	ab := mathutil.ScaleInt(ctx.Difficulty, 3, 20)
	if ab < 2 {
		ab = 2
	}
	if ab >= n/4 {
		ab = n / 5
	}
	if ab < 1 {
		ab = 1
	}

	a := ab + 2 + rand.Intn(max(1, n/3))
	b := ab + 2 + rand.Intn(max(1, n/3))
	if a+b-ab >= n {
		a = ab + (n-ab)/4
		b = ab + (n-ab)/4
	}

	answer := decimal6(float64(a+b-ab) / float64(n))
	question := fmt.Sprintf("A sample space has %d equally likely outcomes. Event A has %d outcomes, event B has %d, and A ∩ B has %d. What is P(A ∪ B)? (enter a decimal)", n, a, b, ab)
	explanation := fmt.Sprintf("Subgoal 1: State governing formula. By inclusion-exclusion P(A ∪ B) = P(A) + P(B) − P(A ∩ B).\nSubgoal 2: Substitute parameters. Substitute total outcomes %d, |A| = %d, |B| = %d and |A ∩ B| = %d.\nSubgoal 3: Evaluate and simplify. P(A ∪ B) = (%d + %d − %d)/%d = %s.", n, a, b, ab, a, b, ab, n, answer)

	return generator.Problem{
		Question:    question,
		Answer:      answer,
		Explanation: explanation,
		Facts: map[string]string{
			"n": strconv.Itoa(n), "a": strconv.Itoa(a), "b": strconv.Itoa(b),
			"ab": strconv.Itoa(ab), "answer": answer,
		},
	}
}

// ---------------------------------------------------------------------------
// Continuity of probability (multiple choice)
// ---------------------------------------------------------------------------

type continuityGen struct{}

func (g *continuityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return continuityProblem()
}

// ---------------------------------------------------------------------------
// The multiplication rule (numeric)
// ---------------------------------------------------------------------------

type multiplicationGen struct{}

var multDenomsLow = []int{12, 16, 20, 25}
var multDenomsMid = []int{30, 40, 50, 60}
var multDenomsHigh = []int{80, 100, 120, 150, 200}

func (g *multiplicationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	tier := mathutil.DifficultyTier(ctx.Difficulty)
	var denoms []int
	switch tier {
	case 0:
		denoms = multDenomsLow
	case 1:
		denoms = multDenomsMid
	default:
		denoms = multDenomsHigh
	}
	n := denoms[rand.Intn(len(denoms))]

	a := mathutil.ScaleInt(ctx.Difficulty, 4, n/2)
	if a < 3 {
		a = 3
	}
	if a >= n {
		a = n - 2
	}

	b := 1 + rand.Intn(a-1)

	answer := decimal6(float64(b) / float64(n))
	question := fmt.Sprintf("An experiment has %d equally likely outcomes. Event A has %d of them, and among those A-outcomes %d also lie in event B. What is P(A ∩ B)? (enter a decimal)", n, a, b)
	explanation := fmt.Sprintf("Subgoal 1: State governing formula. By the multiplication rule P(A ∩ B) = P(A)·P(B | A).\nSubgoal 2: Substitute parameters. P(A) = %d/%d and P(B | A) = %d/%d.\nSubgoal 3: Evaluate and simplify. Multiplying gives P(A ∩ B) = (%d/%d)(%d/%d) = %d/%d = %s.", a, n, b, a, a, n, b, a, b, n, answer)

	return generator.Problem{
		Question:    question,
		Answer:      answer,
		Explanation: explanation,
		Facts: map[string]string{
			"n": strconv.Itoa(n), "a": strconv.Itoa(a), "b": strconv.Itoa(b), "answer": answer,
		},
	}
}

// ---------------------------------------------------------------------------
// Bayes' theorem with priors (numeric)
// ---------------------------------------------------------------------------

type bayesGen struct{}

var bayesDenoms = []int{4, 5, 8, 10, 20}

func (g *bayesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	tA := bayesDenoms[rand.Intn(len(bayesDenoms))]
	tB := bayesDenoms[rand.Intn(len(bayesDenoms))]

	aA := 1 + rand.Intn(tA-1)
	aB := 1 + rand.Intn(tB-1)

	pa := float64(aA) / float64(tA)
	pb := float64(aB) / float64(tB)
	answer := decimal6(pa / (pa + pb))

	question := fmt.Sprintf("Two boxes are equally likely to be chosen. Box one holds %d red marbles out of %d; box two holds %d red out of %d. A red marble is drawn. What is P(box one | red)? (enter a decimal)", aA, tA, aB, tB)
	explanation := fmt.Sprintf("Subgoal 1: State governing formula. By Bayes' theorem with equal priors, P(box one | red) = P(red | box one) / (P(red | box one) + P(red | box two)).\nSubgoal 2: Substitute parameters. Substitute conditional probabilities %d/%d and %d/%d.\nSubgoal 3: Evaluate and simplify. Evaluating the ratio yields (%d/%d)/((%d/%d) + (%d/%d)) = %s.", aA, tA, aB, tB, aA, tA, aA, tA, aB, tB, answer)

	return generator.Problem{
		Question:    question,
		Answer:      answer,
		Explanation: explanation,
		Facts: map[string]string{
			"tA": strconv.Itoa(tA), "aA": strconv.Itoa(aA),
			"tB": strconv.Itoa(tB), "aB": strconv.Itoa(aB), "answer": answer,
		},
	}
}

// ---------------------------------------------------------------------------
// Pairwise and mutual independence (multiple choice)
// ---------------------------------------------------------------------------

type independenceGen struct{}

func (g *independenceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return pairwiseIndependenceProblem()
}

// ---------------------------------------------------------------------------
// Random variables (multiple choice)
// ---------------------------------------------------------------------------

type rvDefinitionGen struct{}

func (g *rvDefinitionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return randomVariableProblem()
}

// ---------------------------------------------------------------------------
// The distribution function (multiple choice)
// ---------------------------------------------------------------------------

type cdfGen struct{}

func (g *cdfGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return cdfProblem()
}

// ---------------------------------------------------------------------------
// The probability mass function (numeric)
// ---------------------------------------------------------------------------

type pmfGen struct{}

func (g *pmfGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		weights []int
		index   int // 1-based
		reason  string
	}
	claims := []claim{
		{weights: []int{2, 4, 6, 8}, index: 4, reason: "The weights sum to 20, so the last one has probability 8/20 = 0.4."},
		{weights: []int{1, 1, 2}, index: 3, reason: "They sum to 4, so the third is 2/4 = 0.5."},
		{weights: []int{3, 5, 2}, index: 1, reason: "They sum to 10, so the first is 3/10 = 0.3."},
		{weights: []int{2, 3, 5}, index: 2, reason: "They sum to 10, so the second is 3/10 = 0.3."},
		{weights: []int{4, 6, 10}, index: 1, reason: "They sum to 20, so the first is 4/20 = 0.2."},
		{weights: []int{5, 3, 2, 10}, index: 4, reason: "They sum to 20, so the fourth is 10/20 = 0.5."},
	}
	c := claims[rand.Intn(len(claims))]
	total := 0
	for _, w := range c.weights {
		total += w
	}
	weight := c.weights[c.index-1]
	answer := decimal6(float64(weight) / float64(total))
	parts := make([]string, len(c.weights))
	for i, w := range c.weights {
		parts[i] = strconv.Itoa(w)
	}
	return generator.Problem{
		Question: "A discrete random variable X has probabilities proportional to the weights " +
			strings.Join(parts, ", ") + " on its outcomes. What is the probability of the outcome with weight " +
			strconv.Itoa(weight) + "? (enter a decimal)",
		Answer:      answer,
		Explanation: "A mass function must sum to one, so each weight is divided by the total: " + c.reason,
		Facts: map[string]string{
			"weights": strings.Join(parts, ","), "index": strconv.Itoa(c.index),
			"total": strconv.Itoa(total), "weight": strconv.Itoa(weight), "answer": answer,
		},
	}
}

// ---------------------------------------------------------------------------
// The probability density function (numeric)
// ---------------------------------------------------------------------------

type pdfGen struct{}

func (g *pdfGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		a      string
		reason string
	}
	claims := []claim{
		{a: "0.5", reason: "Integrating 2x from 0 to 0.5 gives 0.5² = 0.25."},
		{a: "0.3", reason: "0.3² = 0.09."},
		{a: "0.8", reason: "0.8² = 0.64."},
		{a: "0.9", reason: "0.9² = 0.81."},
		{a: "0.2", reason: "0.2² = 0.04."},
		{a: "0.7", reason: "0.7² = 0.49."},
	}
	c := claims[rand.Intn(len(claims))]
	a, _ := strconv.ParseFloat(c.a, 64)
	answer := decimal6(a * a)
	return generator.Problem{
		Question: "A continuous random variable X has density f(x) = 2x on the interval [0, 1]. " +
			"What is P(X ≤ " + c.a + ")? (enter a decimal)",
		Answer:      answer,
		Explanation: "P(X ≤ a) is the integral of 2x from 0 to a, which is a²; " + c.reason,
		Facts:       map[string]string{"a": c.a, "answer": answer},
	}
}

// ---------------------------------------------------------------------------
// Expectation of a function (numeric)
// ---------------------------------------------------------------------------

type expectationGen struct{}

func (g *expectationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		values  []int
		weights []int
		den     int
		reason  string
	}
	claims := []claim{
		{values: []int{0, 1, 2}, weights: []int{2, 3, 5}, den: 10, reason: "0·0.2 + 1·0.3 + 4·0.5 = 2.3."},
		{values: []int{0, 2}, weights: []int{5, 5}, den: 10, reason: "0·0.5 + 4·0.5 = 2."},
		{values: []int{1, 2, 3}, weights: []int{25, 25, 50}, den: 100, reason: "1·0.25 + 4·0.25 + 9·0.5 = 5.75."},
		{values: []int{0, 1, 2, 3}, weights: []int{1, 2, 3, 4}, den: 10, reason: "0 + 0.2 + 4·0.3 + 9·0.4 = 5."},
		{values: []int{1, 4}, weights: []int{6, 4}, den: 10, reason: "1·0.6 + 16·0.4 = 7."},
		{values: []int{0, 1}, weights: []int{2, 8}, den: 10, reason: "0·0.2 + 1·0.8 = 0.8."},
	}
	c := claims[rand.Intn(len(claims))]
	sum := 0
	for i, w := range c.weights {
		sum += c.values[i] * c.values[i] * w
	}
	answer := decimal6(float64(sum) / float64(c.den))
	vals := make([]string, len(c.values))
	wts := make([]string, len(c.weights))
	for i := range c.values {
		vals[i] = strconv.Itoa(c.values[i])
		wts[i] = strconv.Itoa(c.weights[i])
	}
	return generator.Problem{
		Question: "A discrete random variable X takes the values " + strings.Join(vals, ", ") +
			" with probabilities proportional to the weights " + strings.Join(wts, ", ") +
			" (the weights sum to " + strconv.Itoa(c.den) + "). What is E[X²]? (enter a decimal)",
		Answer: answer,
		Explanation: "The law of the unconscious statistician gives E[g(X)] = Σ g(x)p(x), so E[X²] = Σ x²p(x); " +
			c.reason,
		Facts: map[string]string{
			"values": strings.Join(vals, ","), "weights": strings.Join(wts, ","),
			"den": strconv.Itoa(c.den), "answer": answer,
		},
	}
}

// ---------------------------------------------------------------------------
// Variance by the second moment (numeric)
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Moment generating functions (numeric)
// ---------------------------------------------------------------------------

type mgfGen struct{}

func (g *mgfGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		a, b   int
		reason string
	}
	claims := []claim{
		{a: 2, b: 2, reason: "Writing M(t) = e^{μt + σ²t²/2} reads σ² = 2b = 4."},
		{a: 3, b: 1, reason: "σ² = 2b = 2."},
		{a: 5, b: 3, reason: "σ² = 2b = 6."},
		{a: 1, b: 2, reason: "σ² = 2b = 4."},
		{a: 4, b: 2, reason: "σ² = 2b = 4."},
		{a: 0, b: 1, reason: "σ² = 2b = 2."},
	}
	c := claims[rand.Intn(len(claims))]
	answer := strconv.Itoa(2 * c.b)
	return generator.Problem{
		Question: "A random variable has moment generating function M(t) = e^{" + strconv.Itoa(c.a) +
			"t + " + strconv.Itoa(c.b) + "t²}. What is Var(X)? (enter a number)",
		Answer: answer,
		Explanation: "The exponent has the form μt + σ²t²/2, so the coefficient of t² is σ²/2 and σ² = 2b; " +
			c.reason,
		Facts: map[string]string{
			"a": strconv.Itoa(c.a), "b": strconv.Itoa(c.b), "answer": answer,
		},
	}
}

// ---------------------------------------------------------------------------
// Transformations of random variables (numeric)
// ---------------------------------------------------------------------------

type transformGen struct{}

var transformVarsLow = []string{"1", "2", "3", "4", "9"}
var transformVarsMid = []string{"2", "3", "5", "6", "8", "0.5"}
var transformVarsHigh = []string{"0.5", "1.5", "2.5", "3.5", "7", "11"}

func (g *transformGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	tier := mathutil.DifficultyTier(ctx.Difficulty)
	var pool []string
	switch tier {
	case 0:
		pool = transformVarsLow
	case 1:
		pool = transformVarsMid
	default:
		pool = transformVarsHigh
	}
	vStr := pool[rand.Intn(len(pool))]
	c := 2 + rand.Intn(3+tier*2)

	v, _ := strconv.ParseFloat(vStr, 64)
	answer := decimal6(float64(c*c) * v)

	question := fmt.Sprintf("A random variable X has variance %s. What is the variance of the transformed variable Y = %dX? (enter a number)", vStr, c)
	explanation := fmt.Sprintf("Subgoal 1: State governing formula. Scaling a random variable by constant c scales its variance quadratically: Var(cX) = c^2·Var(X).\nSubgoal 2: Substitute parameters. Set scale factor c = %d and base variance Var(X) = %s.\nSubgoal 3: Evaluate and simplify. Var(%dX) = %d^2 · %s = %d · %s = %s.", c, vStr, c, c, vStr, c*c, vStr, answer)

	return generator.Problem{
		Question:    question,
		Answer:      answer,
		Explanation: explanation,
		Facts:       map[string]string{"v": vStr, "c": strconv.Itoa(c), "answer": answer},
	}
}

// ---------------------------------------------------------------------------
// Order statistics (numeric)
// ---------------------------------------------------------------------------

type orderStatsGen struct{}

var orderStatXValues = []string{"0.2", "0.4", "0.5", "0.6", "0.8", "0.9"}

func (g *orderStatsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	tier := mathutil.DifficultyTier(ctx.Difficulty)
	n := 2 + rand.Intn(2+tier*2)
	if n > 6 {
		n = 6
	}
	xStr := orderStatXValues[rand.Intn(len(orderStatXValues))]
	x, _ := strconv.ParseFloat(xStr, 64)
	answer := decimal6(pow(x, n))

	question := fmt.Sprintf("Let %d independent observations be uniform on [0, 1]. What is P(maximum ≤ %s)? (enter a decimal)", n, xStr)
	explanation := fmt.Sprintf("Subgoal 1: State governing formula. The sample maximum X_{(n)} satisfies P(X_{(n)} ≤ x) = (F(x))^n for i.i.d. observations.\nSubgoal 2: Substitute parameters. For %d uniform observations on [0, 1], F(x) = x with threshold x = %s.\nSubgoal 3: Evaluate and simplify. P(maximum ≤ %s) = (%s)^%d = %s.", n, xStr, xStr, xStr, n, answer)

	return generator.Problem{
		Question:    question,
		Answer:      answer,
		Explanation: explanation,
		Facts:       map[string]string{"n": strconv.Itoa(n), "x": xStr, "answer": answer},
	}
}

func pow(x float64, n int) float64 {
	out := 1.0
	for i := 0; i < n; i++ {
		out *= x
	}
	return out
}

// ---------------------------------------------------------------------------
// The memoryless property (multiple choice)
// ---------------------------------------------------------------------------

type memorylessGen struct{}

func (g *memorylessGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return memorylessProblem()
}

// ---------------------------------------------------------------------------
// Markov's inequality (numeric)
// ---------------------------------------------------------------------------

type markovGen struct{}

type markovPair struct {
	mean int
	a    int
}

var markovPairsLow = []markovPair{
	{mean: 8, a: 16}, {mean: 10, a: 20}, {mean: 12, a: 24}, {mean: 10, a: 25}, {mean: 9, a: 12}, {mean: 15, a: 30},
}
var markovPairsMid = []markovPair{
	{mean: 10, a: 40}, {mean: 14, a: 28}, {mean: 16, a: 32}, {mean: 20, a: 25}, {mean: 18, a: 24}, {mean: 25, a: 50},
}
var markovPairsHigh = []markovPair{
	{mean: 12, a: 48}, {mean: 15, a: 60}, {mean: 20, a: 80}, {mean: 30, a: 40}, {mean: 24, a: 32}, {mean: 35, a: 70},
}

func (g *markovGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	tier := mathutil.DifficultyTier(ctx.Difficulty)
	var pool []markovPair
	switch tier {
	case 0:
		pool = markovPairsLow
	case 1:
		pool = markovPairsMid
	default:
		pool = markovPairsHigh
	}
	p := pool[rand.Intn(len(pool))]
	answer := decimal6(float64(p.mean) / float64(p.a))

	question := fmt.Sprintf("A nonnegative random variable has mean %d. By Markov's inequality, what is the upper bound on P(X ≥ %d)? (enter a decimal)", p.mean, p.a)
	explanation := fmt.Sprintf("Subgoal 1: State governing formula. Markov's inequality bounds tail probabilities of nonnegative random variables by P(X ≥ a) ≤ E[X]/a.\nSubgoal 2: Substitute parameters. Substitute expectation E[X] = %d and tail threshold a = %d.\nSubgoal 3: Evaluate and simplify. P(X ≥ %d) ≤ %d/%d = %s.", p.mean, p.a, p.a, p.mean, p.a, answer)

	return generator.Problem{
		Question:    question,
		Answer:      answer,
		Explanation: explanation,
		Facts: map[string]string{
			"mean": strconv.Itoa(p.mean), "a": strconv.Itoa(p.a), "answer": answer,
		},
	}
}

// ---------------------------------------------------------------------------
// Chebyshev's inequality (numeric)
// ---------------------------------------------------------------------------

type chebyshevGen struct{}

type chebyshevPair struct {
	sigma int
	d     int
}

var chebyshevPairsLow = []chebyshevPair{
	{sigma: 1, d: 2}, {sigma: 1, d: 4}, {sigma: 2, d: 4}, {sigma: 2, d: 5}, {sigma: 3, d: 6}, {sigma: 4, d: 5},
}
var chebyshevPairsMid = []chebyshevPair{
	{sigma: 3, d: 10}, {sigma: 5, d: 10}, {sigma: 6, d: 10}, {sigma: 4, d: 8}, {sigma: 5, d: 20}, {sigma: 6, d: 12},
}
var chebyshevPairsHigh = []chebyshevPair{
	{sigma: 8, d: 10}, {sigma: 7, d: 14}, {sigma: 8, d: 20}, {sigma: 9, d: 15}, {sigma: 10, d: 25}, {sigma: 12, d: 20},
}

func (g *chebyshevGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	tier := mathutil.DifficultyTier(ctx.Difficulty)
	var pool []chebyshevPair
	switch tier {
	case 0:
		pool = chebyshevPairsLow
	case 1:
		pool = chebyshevPairsMid
	default:
		pool = chebyshevPairsHigh
	}
	p := pool[rand.Intn(len(pool))]
	answer := decimal6(float64(p.sigma*p.sigma) / float64(p.d*p.d))

	question := fmt.Sprintf("A random variable has standard deviation %d. Chebyshev's inequality bounds P(|X − μ| ≥ %d) by what value? (enter a decimal)", p.sigma, p.d)
	explanation := fmt.Sprintf("Subgoal 1: State governing formula. Chebyshev's inequality bounds departure from the mean by P(|X − μ| ≥ d) ≤ σ^2/d^2.\nSubgoal 2: Substitute parameters. Substitute standard deviation σ = %d and distance d = %d.\nSubgoal 3: Evaluate and simplify. P(|X − μ| ≥ %d) ≤ %d^2/%d^2 = %d/%d = %s.", p.sigma, p.d, p.d, p.sigma, p.d, p.sigma*p.sigma, p.d*p.d, answer)

	return generator.Problem{
		Question:    question,
		Answer:      answer,
		Explanation: explanation,
		Facts: map[string]string{
			"sigma": strconv.Itoa(p.sigma), "d": strconv.Itoa(p.d), "answer": answer,
		},
	}
}

// ---------------------------------------------------------------------------
// The central limit theorem (numeric)
// ---------------------------------------------------------------------------

type cltGen struct{}

type cltCase struct {
	n     int
	mu    int
	sigma int
	x     int
}

var cltCasesLow = []cltCase{
	{n: 16, mu: 5, sigma: 2, x: 88},
	{n: 25, mu: 2, sigma: 1, x: 60},
	{n: 36, mu: 0, sigma: 2, x: 6},
	{n: 64, mu: 0, sigma: 1, x: 8},
	{n: 100, mu: 1, sigma: 2, x: 110},
	{n: 100, mu: 3, sigma: 2, x: 320},
}

var cltCasesMid = []cltCase{
	{n: 25, mu: 4, sigma: 2, x: 110},
	{n: 36, mu: 3, sigma: 1, x: 117},
	{n: 49, mu: 2, sigma: 2, x: 112},
	{n: 64, mu: 2, sigma: 2, x: 144},
	{n: 81, mu: 1, sigma: 3, x: 108},
	{n: 100, mu: 2, sigma: 3, x: 230},
}

var cltCasesHigh = []cltCase{
	{n: 64, mu: 5, sigma: 2, x: 352},
	{n: 100, mu: 5, sigma: 4, x: 560},
	{n: 144, mu: 2, sigma: 2, x: 312},
	{n: 144, mu: 1, sigma: 3, x: 180},
	{n: 100, mu: 4, sigma: 5, x: 475},
	{n: 225, mu: 2, sigma: 2, x: 480},
}

func (g *cltGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	tier := mathutil.DifficultyTier(ctx.Difficulty)
	var pool []cltCase
	switch tier {
	case 0:
		pool = cltCasesLow
	case 1:
		pool = cltCasesMid
	default:
		pool = cltCasesHigh
	}
	c := pool[rand.Intn(len(pool))]
	z := float64(c.x-c.n*c.mu) / (float64(c.sigma) * sqrtInt(c.n))
	answer := decimal6(z)

	question := fmt.Sprintf("The sum S of n independent copies of a variable with mean μ and variance σ² is standardised as z = (x − nμ)/(σ√n). For n = %d, μ = %d, σ = %d and x = %d, what is z? (enter a number)", c.n, c.mu, c.sigma, c.x)
	explanation := fmt.Sprintf("Subgoal 1: State governing formula. The central limit theorem standardises sample sums by z = (x − nμ)/(σ√n).\nSubgoal 2: Substitute parameters. Substitute n = %d, μ = %d, σ = %d and observed sum x = %d.\nSubgoal 3: Evaluate and simplify. Evaluating the standardisation gives z = (%d − %d)/(%d · %d) = %s.", c.n, c.mu, c.sigma, c.x, c.x, c.n*c.mu, c.sigma, int(sqrtInt(c.n)), answer)

	return generator.Problem{
		Question:    question,
		Answer:      answer,
		Explanation: explanation,
		Facts: map[string]string{
			"n": strconv.Itoa(c.n), "mu": strconv.Itoa(c.mu), "sigma": strconv.Itoa(c.sigma),
			"x": strconv.Itoa(c.x), "answer": answer,
		},
	}
}

func sqrtInt(n int) float64 {
	if n <= 0 {
		return 0
	}
	x := float64(n)
	for i := 0; i < 40; i++ {
		x = (x + float64(n)/x) / 2
	}
	return x
}

// ---------------------------------------------------------------------------
// The law of total probability (numeric)
// ---------------------------------------------------------------------------

type totalProbabilityGen struct{}

func (g *totalProbabilityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	tier := mathutil.DifficultyTier(ctx.Difficulty)
	w1 := 1 + rand.Intn(2+tier)
	w2 := 1 + rand.Intn(2+tier)

	t1Pool := []int{2, 4, 5, 10}
	t2Pool := []int{2, 4, 5, 10}
	t1 := t1Pool[rand.Intn(len(t1Pool))]
	t2 := t2Pool[rand.Intn(len(t2Pool))]

	a1 := 1 + rand.Intn(t1-1)
	a2 := 1 + rand.Intn(t2-1)

	total := w1 + w2
	totalProb := float64(w1)/float64(total)*float64(a1)/float64(t1) +
		float64(w2)/float64(total)*float64(a2)/float64(t2)
	answer := decimal6(totalProb)

	question := fmt.Sprintf("Two boxes are chosen with probabilities proportional to the weights %d and %d. Box one yields the event with probability %d/%d, box two with probability %d/%d. What is the total probability of the event? (enter a decimal)", w1, w2, a1, t1, a2, t2)
	explanation := fmt.Sprintf("Subgoal 1: State governing formula. The law of total probability averages conditional probabilities with prior weights: P(E) = (w1/(w1+w2))·P(E | B1) + (w2/(w1+w2))·P(E | B2).\nSubgoal 2: Substitute parameters. Substitute weights %d and %d alongside event rates %d/%d and %d/%d.\nSubgoal 3: Evaluate and simplify. P(E) = (%d/%d)(%d/%d) + (%d/%d)(%d/%d) = %s.", w1, w2, a1, t1, a2, t2, w1, total, a1, t1, w2, total, a2, t2, answer)

	return generator.Problem{
		Question:    question,
		Answer:      answer,
		Explanation: explanation,
		Facts: map[string]string{
			"w1": strconv.Itoa(w1), "a1": strconv.Itoa(a1), "t1": strconv.Itoa(t1),
			"w2": strconv.Itoa(w2), "a2": strconv.Itoa(a2), "t2": strconv.Itoa(t2),
			"answer": answer,
		},
	}
}

// ---------------------------------------------------------------------------
// Joint probability mass functions (numeric)
// ---------------------------------------------------------------------------

type jointPMFGen struct{}

func (g *jointPMFGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	pairs := [][2]int{{0, 0}, {0, 1}, {1, 0}, {1, 1}}
	type claim struct {
		weights []int
		index   int
		reason  string
	}
	claims := []claim{
		{weights: []int{1, 2, 3, 4}, index: 3, reason: "The weights sum to 10 and the pair (1, 0) carries 3, so the mass is 3/10 = 0.3."},
		{weights: []int{2, 2, 1, 5}, index: 4, reason: "They sum to 10 and (1, 1) carries 5, so 5/10 = 0.5."},
		{weights: []int{3, 3, 1, 3}, index: 2, reason: "They sum to 10 and (0, 1) carries 3, so 3/10 = 0.3."},
		{weights: []int{4, 1, 1, 4}, index: 3, reason: "They sum to 10 and (1, 0) carries 1, so 1/10 = 0.1."},
		{weights: []int{5, 2, 2, 1}, index: 1, reason: "They sum to 10 and (0, 0) carries 5, so 5/10 = 0.5."},
		{weights: []int{1, 1, 2, 6}, index: 2, reason: "They sum to 10 and (0, 1) carries 1, so 1/10 = 0.1."},
	}
	c := claims[rand.Intn(len(claims))]
	total := 0
	for _, w := range c.weights {
		total += w
	}
	weight := c.weights[c.index-1]
	pair := pairs[c.index-1]
	answer := decimal6(float64(weight) / float64(total))
	parts := make([]string, len(c.weights))
	for i, w := range c.weights {
		parts[i] = strconv.Itoa(w)
	}
	return generator.Problem{
		Question: "A joint mass function on the four pairs (0,0), (0,1), (1,0) and (1,1) is proportional to the weights " +
			strings.Join(parts, ", ") + " in that order. What is P(X = " + strconv.Itoa(pair[0]) +
			", Y = " + strconv.Itoa(pair[1]) + ")? (enter a decimal)",
		Answer:      answer,
		Explanation: "A joint mass function sums to one over its support, so each weight is divided by the total; " + c.reason,
		Facts: map[string]string{
			"weights": strings.Join(parts, ","), "index": strconv.Itoa(c.index),
			"total": strconv.Itoa(total), "weight": strconv.Itoa(weight),
			"pair": strconv.Itoa(pair[0]) + "," + strconv.Itoa(pair[1]), "answer": answer,
		},
	}
}

// ---------------------------------------------------------------------------
// The Poisson process (numeric)
// ---------------------------------------------------------------------------

type poissonProcessGen struct{}

func (g *poissonProcessGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		m, k   int
		reason string
	}
	claims := []claim{
		{m: 2, k: 1, reason: "e^(-2)·2^1/1! = 0.270671 to six places."},
		{m: 3, k: 2, reason: "e^(-3)·3^2/2! = 0.224042 to six places."},
		{m: 1, k: 0, reason: "e^(-1) = 0.367879 to six places."},
		{m: 4, k: 4, reason: "e^(-4)·4^4/4! = 0.195367 to six places."},
		{m: 2, k: 2, reason: "e^(-2)·2^2/2! = 0.270671 to six places."},
		{m: 3, k: 0, reason: "e^(-3) = 0.049787 to six places."},
	}
	c := claims[rand.Intn(len(claims))]
	answer := decimal6(math.Exp(-float64(c.m)) * pow(float64(c.m), c.k) / factF(c.k))
	return generator.Problem{
		Question: "In a Poisson process the number of arrivals in a fixed interval has mean " +
			strconv.Itoa(c.m) + ". What is the probability of exactly " + strconv.Itoa(c.k) +
			" arrivals? (round to six decimals)",
		Answer: answer,
		Explanation: "The count in an interval is Poisson with mean equal to the rate times the length, so " +
			"P(N = k) = e^(-m)·m^k/k!; " + c.reason,
		Facts: map[string]string{
			"m": strconv.Itoa(c.m), "k": strconv.Itoa(c.k), "answer": answer,
		},
	}
}

func factF(n int) float64 {
	out := 1.0
	for i := 2; i <= n; i++ {
		out *= float64(i)
	}
	return out
}

// ---------------------------------------------------------------------------
// The exponential family (multiple choice)
// ---------------------------------------------------------------------------

type exponentialFamilyGen struct{}

func (g *exponentialFamilyGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return exponentialFamilyProblem()
}

// ---------------------------------------------------------------------------
// The weak law of large numbers (multiple choice)
// ---------------------------------------------------------------------------

type weakLawGen struct{}

func (g *weakLawGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return weakLawProblem()
}
