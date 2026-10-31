package realanalysis

// The first real-analysis batch: twenty-five concepts that open the domain with
// the completeness of the real line, the theory of sequences and series, and the
// first two-thirds of limits and continuity.
//
// Real analysis is where the informal calculus is made rigorous. This batch
// starts from the least upper bound property, draws out the Archimedean
// property, density, countability and the uncountability of the reals, and then
// proves the theorems about sequences, series, limits and continuity that the
// calculus course used. Nineteen concepts are multiple choice, because the
// domain is largely about statements and their proofs, and six are numeric, so
// that the definitions can be exercised as well as quoted.
//
// Every numeric generator publishes the parameters its answer follows from, so
// the tests recompute the answer rather than trusting the author's arithmetic.

import (
	"math"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

func decimal6(v float64) string { return decimalPlaces(v, 6) }

// decimal4 keeps the answer at the precision a learner is likely to type.
// Report 290 answered 8.3891 for 8.389056 and was rejected because the key
// carried six decimals and the numeric tolerance was narrower than a
// four-decimal rounding.
func decimal4(v float64) string { return decimalPlaces(v, 4) }

func decimalPlaces(v float64, places int) string {
	s := strconv.FormatFloat(v, 'f', places, 64)
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

// mcClaim is one multiple-choice item: a question, the correct option, the
// distractors that must all be rejected, and the reason the correct option is
// right.
type mcClaim struct {
	q           string
	correct     string
	distractors []string
	reason      string
}

func pickMC(claims []mcClaim) generator.Problem {
	c := claims[rand.Intn(len(claims))]
	p := generator.MultipleChoice(c.q, c.correct, c.distractors...)
	p.Explanation = c.reason
	return p
}

// ---------------------------------------------------------------------------
// Supremum and infimum (numeric)
// ---------------------------------------------------------------------------

type supremumGen struct{}

func (g *supremumGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		lo, hi int
		set    string
	}
	claims := []claim{
		{lo: 0, hi: 1, set: "the open interval of all x with 0 < x < 1"},
		{lo: -2, hi: 3, set: "the half-open interval of all x with −2 < x ≤ 3"},
		{lo: 1, hi: 5, set: "the half-open interval of all x with 1 ≤ x < 5"},
		{lo: 2, hi: 4, set: "the open interval of all x with 2 < x < 4"},
		{lo: 0, hi: 10, set: "the closed interval of all x with 0 ≤ x ≤ 10"},
		{lo: -5, hi: -1, set: "the open interval of all x with −5 < x < −1"},
	}
	c := claims[rand.Intn(len(claims))]
	answer := strconv.Itoa(c.hi)
	return generator.Problem{
		Question: "What is the supremum of " + c.set + "? (enter a decimal)",
		Answer:   answer,
		Explanation: "The supremum is the least upper bound: it is an upper bound, and no smaller number is. Every " +
			c.set + " is at most " + strconv.Itoa(c.hi) + ", and numbers just below it are exceeded by elements of the set, so the least upper bound is " +
			strconv.Itoa(c.hi) + ".",
		Facts: map[string]string{
			"lo": strconv.Itoa(c.lo), "hi": strconv.Itoa(c.hi), "answer": answer,
		},
	}
}

// ---------------------------------------------------------------------------
// The Archimedean property (numeric)
// ---------------------------------------------------------------------------

type archimedeanGen struct{}

func (g *archimedeanGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		pn, pd int
	}
	claims := []claim{
		{pn: 1, pd: 100},
		{pn: 1, pd: 10},
		{pn: 2, pd: 5},
		{pn: 1, pd: 1000},
		{pn: 3, pd: 4},
		{pn: 1, pd: 7},
	}
	c := claims[rand.Intn(len(claims))]
	n := c.pd/c.pn + 1
	if c.pd%c.pn == 0 {
		n = c.pd/c.pn + 1
	} else {
		n = int(math.Floor(float64(c.pd)/float64(c.pn))) + 1
	}
	answer := strconv.Itoa(n)
	return generator.Problem{
		Question: "By the Archimedean property there is a positive integer N with 1/N < ε. For ε = " +
			strconv.Itoa(c.pn) + "/" + strconv.Itoa(c.pd) + ", what is the smallest such N? (enter a whole number)",
		Answer:      answer,
		Explanation: "The inequality 1/N < ε is the same as N > 1/ε, so the smallest integer N is the floor of 1/ε plus one; here 1/ε = " + strconv.Itoa(c.pd) + "/" + strconv.Itoa(c.pn) + ", giving N = " + answer + ".",
		Facts: map[string]string{
			"pn": strconv.Itoa(c.pn), "pd": strconv.Itoa(c.pd), "answer": answer,
		},
	}
}

// ---------------------------------------------------------------------------
// Convergence of sequences (numeric)
// ---------------------------------------------------------------------------

type sequenceGen struct{}

func (g *sequenceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	cDenoms := []int{2, 4, 5, 8, 10}
	var a, b, c, d int
	var answer, question string
	// Bounded rather than open ended: see generator.MaxDrawAttempts. If no draw
	// avoids printing its own answer the last one is returned, which is a
	// readoff leak, not a hung request.
	for attempt := 0; attempt < generator.MaxDrawAttempts; attempt++ {
		c = cDenoms[rand.Intn(len(cDenoms))]
		a = rand.Intn(40) + 1
		b = rand.Intn(15) + 1
		d = rand.Intn(15) + 1
		answer = decimal6(float64(a) / float64(c))
		question = "The sequence with terms (" + strconv.Itoa(a) + "n + " + strconv.Itoa(b) + ")/(" +
			strconv.Itoa(c) + "n + " + strconv.Itoa(d) + ") converges as n tends to infinity. What is its limit? (enter a decimal)"
		if answer != strconv.Itoa(a) && answer != strconv.Itoa(b) && answer != strconv.Itoa(c) && answer != strconv.Itoa(d) && !strings.Contains(question, answer) {
			break
		}
	}
	return generator.Problem{
		Question:    question,
		Answer:      answer,
		Explanation: "Dividing the numerator and denominator by n, the ratio of the leading coefficients is the limit, because the constant terms become negligible; here that ratio is " + strconv.Itoa(a) + "/" + strconv.Itoa(c) + " = " + answer + ".",
		Facts: map[string]string{
			"a": strconv.Itoa(a), "b": strconv.Itoa(b), "c": strconv.Itoa(c), "d": strconv.Itoa(d), "answer": answer,
		},
	}
}

// ---------------------------------------------------------------------------
// Monotone convergence (numeric)
// ---------------------------------------------------------------------------

type monotoneGen struct{}

func (g *monotoneGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// a = L(L-1) makes the fixed point of L = sqrt(a + L) equal to L.
	claims := []int{2, 6, 12, 20, 30, 42}
	a := claims[rand.Intn(len(claims))]
	l := int((1 + math.Sqrt(1+4*float64(a))) / 2)
	answer := strconv.Itoa(l)
	return generator.Problem{
		Question: "The increasing sequence x_1 = 0, x_{n+1} = sqrt(a + x_n) is bounded above and converges. Its limit L satisfies L = sqrt(a + L). For a = " +
			strconv.Itoa(a) + ", what is L? (enter a decimal)",
		Answer:      answer,
		Explanation: "Squaring the fixed-point equation gives L² − L − a = 0, whose positive root is (1 + sqrt(1 + 4a))/2; the increasing sequence starting at zero converges to that root.",
		Facts: map[string]string{
			"a": strconv.Itoa(a), "answer": answer,
		},
	}
}

// ---------------------------------------------------------------------------
// Series and partial sums (numeric)
// ---------------------------------------------------------------------------

type seriesGen struct{}

func (g *seriesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		pn, pd int
	}
	claims := []claim{
		{pn: 1, pd: 2},
		{pn: 1, pd: 3},
		{pn: 1, pd: 4},
		{pn: 2, pd: 3},
		{pn: 1, pd: 5},
		{pn: 3, pd: 5},
	}
	c := claims[rand.Intn(len(claims))]
	r := float64(c.pn) / float64(c.pd)
	answer := decimal6(1 / (1 - r))
	return generator.Problem{
		Question: "The geometric series with first term one and ratio r = " + strconv.Itoa(c.pn) + "/" + strconv.Itoa(c.pd) +
			" converges to 1/(1 − r). What is its sum? (enter a decimal)",
		Answer:      answer,
		Explanation: "A geometric series with ratio of absolute value less than one converges to the first term divided by one minus the ratio, which is the limit of its partial sums.",
		Facts: map[string]string{
			"pn": strconv.Itoa(c.pn), "pd": strconv.Itoa(c.pd), "answer": answer,
		},
	}
}

// ---------------------------------------------------------------------------
// The limit of a function (numeric)
// ---------------------------------------------------------------------------

type functionLimitGen struct{}

func (g *functionLimitGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		m, en, ed int
	}
	claims := []claim{
		{m: 3, en: 1, ed: 1},
		{m: 2, en: 1, ed: 1},
		{m: 4, en: 1, ed: 1},
		{m: 5, en: 2, ed: 1},
		{m: 3, en: 2, ed: 1},
		{m: 2, en: 3, ed: 1},
	}
	c := claims[rand.Intn(len(claims))]
	eps := float64(c.en) / float64(c.ed)
	answer := decimal6(eps / float64(c.m))
	return generator.Problem{
		Question: "For the linear function f(x) = " + strconv.Itoa(c.m) + "x − 1 and its limit at x = 2, the largest δ that makes |f(x) − f(2)| < ε whenever |x − 2| < δ is ε divided by the slope. For ε = " +
			strconv.Itoa(c.en) + "/" + strconv.Itoa(c.ed) + ", what is δ? (enter a decimal)",
		Answer:      answer,
		Explanation: "Because the function is linear, |f(x) − f(2)| = m|x − 2|, so the condition is m|x − 2| < ε, that is |x − 2| < ε/m; the largest such δ is exactly ε/m.",
		Facts: map[string]string{
			"m": strconv.Itoa(c.m), "en": strconv.Itoa(c.en), "ed": strconv.Itoa(c.ed), "answer": answer,
		},
	}
}

// ---------------------------------------------------------------------------
// The completeness axiom (multiple choice)
// ---------------------------------------------------------------------------

type completenessGen struct{}

func (g *completenessGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return completenessProblem()
}

// ---------------------------------------------------------------------------
// Nested intervals (multiple choice)
// ---------------------------------------------------------------------------

type nestedIntervalsGen struct{}

func (g *nestedIntervalsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return nestedIntervalsProblem()
}

// ---------------------------------------------------------------------------
// Density of the rationals (multiple choice)
// ---------------------------------------------------------------------------

type densityGen struct{}

func (g *densityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return densityProblem()
}

// ---------------------------------------------------------------------------
// Countable sets (multiple choice)
// ---------------------------------------------------------------------------

type countableGen struct{}

func (g *countableGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return countableProblem()
}

// ---------------------------------------------------------------------------
// The uncountability of the reals (multiple choice)
// ---------------------------------------------------------------------------

type uncountableGen struct{}

func (g *uncountableGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return uncountableProblem()
}

// ---------------------------------------------------------------------------
// The extended real line (multiple choice)
// ---------------------------------------------------------------------------

type extendedGen struct{}

func (g *extendedGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return extendedStatementsProblem()
}

// ---------------------------------------------------------------------------
// Subsequences (multiple choice)
// ---------------------------------------------------------------------------

type subsequenceGen struct{}

func (g *subsequenceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return subsequenceProblem()
}

// ---------------------------------------------------------------------------
// The Bolzano-Weierstrass theorem (multiple choice)
// ---------------------------------------------------------------------------

type bolzanoWeierstrassGen struct{}

func (g *bolzanoWeierstrassGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return bolzanoProblem()
}

// ---------------------------------------------------------------------------
// Cauchy sequences (multiple choice)
// ---------------------------------------------------------------------------

type cauchyGen struct{}

func (g *cauchyGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return cauchyProblem()
}

// ---------------------------------------------------------------------------
// Comparison and ratio tests (multiple choice)
// ---------------------------------------------------------------------------

type ratioTestGen struct{}

func (g *ratioTestGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return seriesTestsProblem()
}

// ---------------------------------------------------------------------------
// Absolute and conditional convergence (multiple choice)
// ---------------------------------------------------------------------------

type absoluteGen struct{}

func (g *absoluteGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return absoluteProblem()
}

// ---------------------------------------------------------------------------
// Rearrangement of series (multiple choice)
// ---------------------------------------------------------------------------

type rearrangementGen struct{}

func (g *rearrangementGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return rearrangementProblem()
}

// ---------------------------------------------------------------------------
// The Cauchy criterion for series (multiple choice)
// ---------------------------------------------------------------------------

type cauchyCriterionGen struct{}

func (g *cauchyCriterionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return cauchyCriterionStatementsProblem()
}

// ---------------------------------------------------------------------------
// The sequential criterion (multiple choice)
// ---------------------------------------------------------------------------

type sequentialGen struct{}

func (g *sequentialGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return sequentialProblem()
}

// ---------------------------------------------------------------------------
// Continuity and its algebra (multiple choice)
// ---------------------------------------------------------------------------

type continuityGen struct{}

func (g *continuityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return continuityProblem()
}

// ---------------------------------------------------------------------------
// The intermediate value theorem (multiple choice)
// ---------------------------------------------------------------------------

type ivtGen struct{}

func (g *ivtGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return ivtProblem()
}

// ---------------------------------------------------------------------------
// The extreme value theorem (multiple choice)
// ---------------------------------------------------------------------------

type evtGen struct{}

func (g *evtGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return evtProblem()
}

// ---------------------------------------------------------------------------
// Uniform continuity (multiple choice)
// ---------------------------------------------------------------------------

type uniformGen struct{}

func (g *uniformGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return uniformProblem()
}

// ---------------------------------------------------------------------------
// Compactness and Heine-Borel (multiple choice)
// ---------------------------------------------------------------------------

type compactGen struct{}

func (g *compactGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return compactStatementsProblem()
}
