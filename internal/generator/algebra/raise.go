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
	kp := waveKP(ctx)
	// The function stays the two-piece rule the authored schema describes.
	// The ordinal chooses which side of the boundary the learner meets: KP1
	// mixes both branches, KP2 stays on the negative branch and KP3 stays on
	// the nonnegative branch. Each range is wider than ten inputs.
	var x0 int
	switch kp {
	case 2:
		x0 = -(rand.Intn(30) + 1)
	case 3:
		x0 = rand.Intn(31)
	default:
		x0 = rand.Intn(21) - 10
	}
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
		KPOrdinal: kp,
	}
}

// ---------------------------------------------------------------------------
// Graph transformations (multiple choice)
// ---------------------------------------------------------------------------

type transformationsGen struct{}

func (g *transformationsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	x0 := rand.Intn(9) - 4
	y0 := rand.Intn(9) - 4
	k := rand.Intn(8) + 2
	c := rand.Intn(5) + 2
	mode := rand.Intn(6)
	var question, explanation, answer string
	var coord string
	switch mode {
	case 0:
		coord = "y"
		answer = strconv.Itoa(y0 + k)
		question = fmt.Sprintf("If \\((%d, %d)\\) lies on \\(y = f(x)\\), what is the y-coordinate of the corresponding point on \\(y = f(x) + %d\\)? (enter a number)", x0, y0, k)
		explanation = fmt.Sprintf("Adding %d outside the function adds %d to every y-coordinate: \\(%d + %d = %d\\).", k, k, y0, k, y0+k)
	case 1:
		coord = "x"
		answer = strconv.Itoa(x0 + k)
		question = fmt.Sprintf("If \\((%d, %d)\\) lies on \\(y = f(x)\\), what is the x-coordinate of the corresponding point on \\(y = f(x - %d)\\)? (enter a number)", x0, y0, k)
		explanation = fmt.Sprintf("Replacing x by x - %d shifts the graph right by %d: \\(%d + %d = %d\\).", k, k, x0, k, x0+k)
	case 2:
		coord = "y"
		answer = strconv.Itoa(c * y0)
		question = fmt.Sprintf("If \\((%d, %d)\\) lies on \\(y = f(x)\\), what is the y-coordinate of the corresponding point on \\(y = %d f(x)\\)? (enter a number)", x0, y0, c)
		explanation = fmt.Sprintf("The vertical scale factor multiplies every y-coordinate by %d: \\(%d \\times %d = %d\\).", c, c, y0, c*y0)
	case 3:
		x0 = (rand.Intn(7) - 3) * c
		coord = "x"
		answer = strconv.Itoa(x0 / c)
		question = fmt.Sprintf("If \\((%d, %d)\\) lies on \\(y = f(x)\\), what is the x-coordinate of the corresponding point on \\(y = f(%d x)\\)? (enter a number)", x0, y0, c)
		explanation = fmt.Sprintf("The horizontal scale factor %d divides every x-coordinate: \\(%d / %d = %d\\).", c, x0, c, x0/c)
	case 4:
		coord = "y"
		answer = strconv.Itoa(-y0)
		question = fmt.Sprintf("If \\((%d, %d)\\) lies on \\(y = f(x)\\), what is the y-coordinate of the corresponding point on \\(y = -f(x)\\)? (enter a number)", x0, y0)
		explanation = fmt.Sprintf("The reflection across the x-axis negates every y-coordinate: \\(-%d = %d\\).", y0, -y0)
	default:
		coord = "x"
		answer = strconv.Itoa(-x0)
		question = fmt.Sprintf("If \\((%d, %d)\\) lies on \\(y = f(x)\\), what is the x-coordinate of the corresponding point on \\(y = f(-x)\\)? (enter a number)", x0, y0)
		explanation = fmt.Sprintf("The reflection across the y-axis negates every x-coordinate: \\(-%d = %d\\).", x0, -x0)
	}
	return generator.Problem{
		Question:    question,
		Answer:      answer,
		Explanation: explanation,
		Facts:       map[string]string{"x": strconv.Itoa(x0), "y": strconv.Itoa(y0), "mode": strconv.Itoa(mode), "answer": answer, "coordinate": coord},
	}
}

// ---------------------------------------------------------------------------
// Domain of an inverse (multiple choice)
// ---------------------------------------------------------------------------

type inverseDomainGen struct{}

func (g *inverseDomainGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := rand.Intn(4)
	switch mode {
	case 0:
		a := rand.Intn(5) - 4
		b := a + rand.Intn(4) + 1
		c := b + rand.Intn(4) + 1
		d := c + rand.Intn(4) + 2
		return generator.Problem{
			Question:    fmt.Sprintf("A one-to-one function \\(f\\) has domain \\([%d, %d]\\) and range \\([%d, %d]\\). What is the lower endpoint of the domain of \\(f^{-1}\\)? (enter a number)", a, b, c, d),
			Answer:      strconv.Itoa(c),
			Explanation: fmt.Sprintf("The domain of \\(f^{-1}\\) is the range of \\(f\\), so it is \\([%d, %d]\\) and its lower endpoint is %d.", c, d, c),
			Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "c": strconv.Itoa(c), "d": strconv.Itoa(d)},
		}
	case 1:
		a := rand.Intn(7) + 1
		b := rand.Intn(6) + 1
		for b == a {
			b = rand.Intn(6) + 1
		}
		return generator.Problem{
			Question:    fmt.Sprintf("The one-to-one function \\(f(x) = \\sqrt{x - %d} + %d\\) has domain \\([%d, \\infty)\\). What is the lower endpoint of the domain of \\(f^{-1}\\)? (enter a number)", a, b, a),
			Answer:      strconv.Itoa(b),
			Explanation: fmt.Sprintf("The square root is at least 0, so the range of \\(f\\) is \\([%d, \\infty)\\), which is the domain of \\(f^{-1}\\).", b),
			Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b)},
		}
	case 2:
		a := rand.Intn(6) + 1
		b := rand.Intn(6) + 1
		for b == a {
			b = rand.Intn(6) + 1
		}
		return generator.Problem{
			Question:    fmt.Sprintf("The quadratic \\(f(x) = (x - %d)^2 + %d\\) with domain \\([%d, \\infty)\\) has an inverse. What is the lower endpoint of the domain of \\(f^{-1}\\)? (enter a number)", a, b, a),
			Answer:      strconv.Itoa(b),
			Explanation: fmt.Sprintf("The minimum value of \\(f\\) is %d, so its range is \\([%d, \\infty)\\) and that is the domain of \\(f^{-1}\\).", b, b),
			Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b)},
		}
	default:
		a := rand.Intn(5) + 1
		b := rand.Intn(5) + 1
		for b == a {
			b = rand.Intn(5) + 1
		}
		return generator.Problem{
			Question:    fmt.Sprintf("The increasing function \\(f(x) = e^{x - %d} + %d\\) has domain all real numbers. What is the lower endpoint of the domain of \\(f^{-1}\\)? (enter a number)", a, b),
			Answer:      strconv.Itoa(b),
			Explanation: fmt.Sprintf("Since \\(e^{x - %d} > 0\\), the range of \\(f\\) is \\((%d, \\infty)\\), so the lower endpoint is %d.", a, b, b),
			Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b)},
		}
	}
}

// ---------------------------------------------------------------------------
// Exponential graphs (numeric)
// ---------------------------------------------------------------------------

type expGraphGen struct{}

func (g *expGraphGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := waveKP(ctx)
	// The function keeps the shape f(x) = a b^x the authored schema names.
	// KP1 reads a small exponent, KP2 reads a larger one so the constant
	// factor per unit step is visible, and KP3 uses a base above five where
	// the growth is unmistakable.
	a := rand.Intn(8) + 2
	var b, x int
	switch kp {
	case 2:
		b = rand.Intn(4) + 2
		x = rand.Intn(2) + 4
	case 3:
		b = rand.Intn(4) + 6
		x = rand.Intn(2) + 2
	default:
		b = rand.Intn(4) + 2
		x = rand.Intn(2) + 2
	}
	value := a * algPowInt(b, x)
	return generator.Problem{
		Question: "For the exponential function f(x) = " + strconv.Itoa(a) + " times " + strconv.Itoa(b) +
			"^x, what is f(" + strconv.Itoa(x) + ")?",
		Answer: strconv.Itoa(value),
		Explanation: "Raising the base " + strconv.Itoa(b) + " to the power " + strconv.Itoa(x) + " and multiplying by " +
			strconv.Itoa(a) + " gives " + strconv.Itoa(value) + ".",
		Facts: map[string]string{
			"a": strconv.Itoa(a), "b": strconv.Itoa(b), "x": strconv.Itoa(x), "answer": strconv.Itoa(value),
		},
		KPOrdinal: kp,
	}
}

// ---------------------------------------------------------------------------
// Compound interest (numeric)
// ---------------------------------------------------------------------------

type compoundGen struct{}

func (g *compoundGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := waveKP(ctx)
	// The formula stays annual compounding, A = P(1 + r)^t, as the authored
	// schema states. The ordinal moves along the same curve: KP1 draws a
	// short horizon and a small principal, KP2 a medium one and KP3 a long
	// one, so no ordinal repeats another's parameter set.
	var p, t int
	switch kp {
	case 2:
		p = (rand.Intn(30) + 20) * 100 // 2000..4900
		t = rand.Intn(3) + 5           // 5..7
	case 3:
		p = (rand.Intn(50) + 50) * 100 // 5000..9900
		t = rand.Intn(3) + 8           // 8..10
	default:
		p = (rand.Intn(15) + 5) * 100 // 500..1900
		t = rand.Intn(3) + 2          // 2..4
	}
	r := rand.Intn(8) + 2 // 2..9 percent
	value := float64(p) * math.Pow(1+float64(r)/100, float64(t))
	answer := fmt.Sprintf("%.2f", value)
	return generator.Problem{
		Question: "An investment of " + strconv.Itoa(p) + " dollars earns " + strconv.Itoa(r) +
			" percent interest per year, compounded once a year, for " + strconv.Itoa(t) +
			" years. Using A = P(1 + r)^t, what is the value at the end, in dollars? (round to cents)",
		Answer: answer,
		Explanation: "With P = " + strconv.Itoa(p) + ", the annual factor is 1 + " + strconv.Itoa(r) +
			"/100, and raising it to the power " + strconv.Itoa(t) + " and multiplying by P gives " + answer + " dollars.",
		Facts: map[string]string{
			"p": strconv.Itoa(p), "r": strconv.Itoa(r), "t": strconv.Itoa(t), "answer": answer,
		},
		KPOrdinal: kp,
	}
}

// ---------------------------------------------------------------------------
// Logarithmic graphs (multiple choice)
// ---------------------------------------------------------------------------

type logGraphGen struct{}

func (g *logGraphGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	b := rand.Intn(8) + 2
	h := rand.Intn(8) + 1
	switch rand.Intn(5) {
	case 0:
		if rand.Intn(2) == 0 {
			return generator.Problem{
				Question:    fmt.Sprintf("What is the x-coordinate of the vertical asymptote of \\(y = \\log_{%d}(x - %d)\\)? (enter a number)", b, h),
				Answer:      strconv.Itoa(h),
				Explanation: fmt.Sprintf("The argument must be positive: \\(x - %d > 0\\), so the asymptote is \\(x = %d\\).", h, h),
				Facts:       map[string]string{"base": strconv.Itoa(b), "shift": strconv.Itoa(h), "answer": strconv.Itoa(h)},
			}
		}
		return generator.Problem{
			Question:    fmt.Sprintf("What is the x-coordinate of the vertical asymptote of \\(y = \\log_{%d}(x + %d)\\)? (enter a number)", b, h),
			Answer:      strconv.Itoa(-h),
			Explanation: fmt.Sprintf("Set the argument to zero: \\(x + %d = 0\\), so the asymptote is \\(x = %d\\).", h, -h),
			Facts:       map[string]string{"base": strconv.Itoa(b), "shift": strconv.Itoa(-h), "answer": strconv.Itoa(-h)},
		}
	case 1:
		return generator.Problem{
			Question:    fmt.Sprintf("What is the x-coordinate of the x-intercept of \\(y = \\log_{%d}(x - %d)\\)? (enter a number)", b, h),
			Answer:      strconv.Itoa(h + 1),
			Explanation: fmt.Sprintf("The logarithm is zero when its argument is 1: \\(x - %d = 1\\), so \\(x = %d\\).", h, h+1),
			Facts:       map[string]string{"base": strconv.Itoa(b), "shift": strconv.Itoa(h), "answer": strconv.Itoa(h + 1)},
		}
	case 2:
		return generator.Problem{
			Question:    fmt.Sprintf("What is the lower endpoint of the domain of \\(f(x) = \\log_{%d}(x - %d)\\)? (enter a number)", b, h),
			Answer:      strconv.Itoa(h),
			Explanation: fmt.Sprintf("The argument must be strictly positive: \\(x - %d > 0\\), so the domain is \\((%d, \\infty)\\).", h, h),
			Facts:       map[string]string{"base": strconv.Itoa(b), "shift": strconv.Itoa(h), "answer": strconv.Itoa(h)},
		}
	case 3:
		k := rand.Intn(3) + 1
		x := h + algPowInt(b, k)
		return generator.Problem{
			Question:    fmt.Sprintf("Solve \\(\\log_{%d}(x - %d) = %d\\) for \\(x\\). (enter a number)", b, h, k),
			Answer:      strconv.Itoa(x),
			Explanation: fmt.Sprintf("Rewrite as an exponential: \\(x - %d = %d^{%d} = %d\\), so \\(x = %d\\).", h, b, k, algPowInt(b, k), x),
			Facts:       map[string]string{"base": strconv.Itoa(b), "shift": strconv.Itoa(h), "exponent": strconv.Itoa(k), "answer": strconv.Itoa(x)},
		}
	default:
		k := rand.Intn(4) + 1
		val := algPowInt(b, k)
		return generator.Problem{
			Question:    fmt.Sprintf("What is \\(\\log_{%d}(%d)\\)? (enter a number)", b, val),
			Answer:      strconv.Itoa(k),
			Explanation: fmt.Sprintf("\\(%d^{%d} = %d\\), so the logarithm is %d.", b, k, val, k),
			Facts:       map[string]string{"base": strconv.Itoa(b), "value": strconv.Itoa(val), "answer": strconv.Itoa(k)},
		}
	}
}

func algPowInt(base, exp int) int {
	out := 1
	for i := 0; i < exp; i++ {
		out *= base
	}
	return out
}

// ---------------------------------------------------------------------------
// Logarithmic models (numeric)
// ---------------------------------------------------------------------------

type logModelGen struct{}

func (g *logModelGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := waveKP(ctx)
	// The equation stays log_b(x) = k. The ordinal selects the exponent band:
	// KP1 rewrites a one- or two-step logarithm, KP2 a three- or four-step
	// one and KP3 a five- or six-step one, so the three ordinals do not share
	// a single fixed table.
	b := rand.Intn(7) + 2 // 2..8
	var k int
	switch kp {
	case 2:
		k = rand.Intn(2) + 3 // 3..4
	case 3:
		k = rand.Intn(2) + 5 // 5..6
	default:
		k = rand.Intn(2) + 1 // 1..2
	}
	value := algPowInt(b, k)
	return generator.Problem{
		Question: "Solve log base " + strconv.Itoa(b) + " of x = " + strconv.Itoa(k) +
			" for x. (enter a number)",
		Answer: strconv.Itoa(value),
		Explanation: "The equation says the base " + strconv.Itoa(b) + " raised to the power " + strconv.Itoa(k) +
			" equals x, so x = " + strconv.Itoa(value) + ".",
		Facts: map[string]string{
			"b": strconv.Itoa(b), "k": strconv.Itoa(k), "answer": strconv.Itoa(value),
		},
		KPOrdinal: kp,
	}
}

// ---------------------------------------------------------------------------
// Simplifying radicals (numeric)
// ---------------------------------------------------------------------------

type radicalSimplifyGen struct{}

func (g *radicalSimplifyGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := waveKP(ctx)
	// The squarefree factor b is drawn from a band per ordinal, and a covers a
	// wide range, so n = a^2 b is a fresh number rather than one of eight.
	squarefree := []int{2, 3, 5, 6, 7, 10, 11, 13, 14, 15, 17}
	var a, b int
	switch kp {
	case 2:
		a = rand.Intn(8) + 9 // 9..16
		b = squarefree[rand.Intn(7)]
	case 3:
		a = rand.Intn(11) + 2 // 2..12
		b = squarefree[rand.Intn(6)+5]
	default:
		a = rand.Intn(7) + 2 // 2..8
		b = squarefree[rand.Intn(5)]
	}
	n := a * a * b
	return generator.Problem{
		Question:    fmt.Sprintf("The square root of %d can be written as \\(a\\sqrt{b}\\), where \\(b\\) is squarefree. What is the whole number \\(a\\)? (enter a number)", n),
		Answer:      strconv.Itoa(a),
		Explanation: fmt.Sprintf("Since %d = %d times %d, the square root is \\(%d\\sqrt{%d}\\).", n, a*a, b, a, b),
		Facts: map[string]string{
			"n": strconv.Itoa(n), "a": strconv.Itoa(a), "b": strconv.Itoa(b), "answer": strconv.Itoa(a),
		},
		KPOrdinal: kp,
	}
}

// ---------------------------------------------------------------------------
// Operations with radicals (numeric)
// ---------------------------------------------------------------------------

type radicalOpsGen struct{}

func (g *radicalOpsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := waveKP(ctx)
	// The operation stays "add two like radicals". The ordinal chooses a band
	// of squarefree radicands and a coefficient range, so the three ordinals
	// pose different families instead of replaying six fixed pairs.
	bands := [3][]int{
		{2, 3, 5, 6},
		{7, 10, 11, 13, 14},
		{15, 17, 19, 21, 22},
	}
	band := bands[kp-1]
	c := band[rand.Intn(len(band))]
	var a, b int
	switch kp {
	case 2:
		a = rand.Intn(12) + 3 // 3..14
		b = rand.Intn(12) + 3
	case 3:
		a = rand.Intn(15) + 2 // 2..16
		b = rand.Intn(15) + 2
	default:
		a = rand.Intn(9) + 2 // 2..10
		b = rand.Intn(9) + 2
	}
	sum := a + b
	return generator.Problem{
		Question: "Simplify " + strconv.Itoa(a) + " times the square root of " + strconv.Itoa(c) + " plus " +
			strconv.Itoa(b) + " times the square root of " + strconv.Itoa(c) +
			". The result is k times the square root of " + strconv.Itoa(c) + "; what is k? (enter a number)",
		Answer: strconv.Itoa(sum),
		Explanation: "The two terms share the radical, so their coefficients add: " + strconv.Itoa(a) + " + " +
			strconv.Itoa(b) + " = " + strconv.Itoa(sum) + ".",
		Facts: map[string]string{
			"a": strconv.Itoa(a), "b": strconv.Itoa(b), "c": strconv.Itoa(c), "answer": strconv.Itoa(sum),
		},
		KPOrdinal: kp,
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
	kp := waveKP(ctx)
	// The proportion keeps the shape a/b = num/x the authored schema names.
	// The scale k sets num = ak and x = bk, so every draw divides exactly; the
	// ordinal uses a different band of scales and a wider pair of ratios.
	var a, b, k, lo, hi int
	switch kp {
	case 2:
		lo, hi = 3, 12
		k = rand.Intn(3) + 5 // 5..7
	case 3:
		lo, hi = 2, 12
		k = rand.Intn(4) + 8 // 8..11
	default:
		lo, hi = 2, 10
		k = rand.Intn(3) + 2 // 2..4
	}
	a = rand.Intn(hi-lo+1) + lo
	b = rand.Intn(hi-lo+1) + lo
	// A proportion with a equal to b is the trivial 1 = num/x, so redraw.
	for b == a {
		b = rand.Intn(hi-lo+1) + lo
	}
	num := a * k
	x := b * k
	return generator.Problem{
		Question: "Solve the proportion " + strconv.Itoa(a) + "/" + strconv.Itoa(b) + " = " + strconv.Itoa(num) +
			"/x for x. (enter a number)",
		Answer: strconv.Itoa(x),
		Explanation: "Cross-multiplying gives " + strconv.Itoa(a) + "x = " + strconv.Itoa(b) + " times " + strconv.Itoa(num) +
			", so x = " + strconv.Itoa(x) + ".",
		Facts: map[string]string{
			"a": strconv.Itoa(a), "b": strconv.Itoa(b), "num": strconv.Itoa(num), "answer": strconv.Itoa(x),
		},
		KPOrdinal: kp,
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
		if c.m != 0 && c.n != 0 && c.m != c.n && c.m+c.n != 0 && c.a != -c.m && c.a != -c.n {
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
	switch rand.Intn(3) {
	case 0:
		eqs := rand.Intn(3) + 2
		vars := rand.Intn(3) + 2
		entries := eqs * (vars + 1)
		return generator.Problem{
			Question:    fmt.Sprintf("How many entries does the augmented matrix for a linear system of %d equations in %d variables have? (enter a number)", eqs, vars),
			Answer:      strconv.Itoa(entries),
			Explanation: fmt.Sprintf("The matrix has %d rows and %d columns, so it has \\(%d \\times %d = %d\\) entries.", eqs, vars+1, eqs, vars+1, entries),
			Facts:       map[string]string{"equations": strconv.Itoa(eqs), "variables": strconv.Itoa(vars), "entries": strconv.Itoa(entries)},
		}
	case 1:
		vars := rand.Intn(3) + 3
		rank := vars - 1
		free := vars - rank
		return generator.Problem{
			Question:    fmt.Sprintf("A consistent linear system in %d variables has rank %d. How many free variables does it have? (enter a number)", vars, rank),
			Answer:      strconv.Itoa(free),
			Explanation: fmt.Sprintf("The number of free variables is the number of variables minus the rank: \\(%d - %d = %d\\).", vars, rank, free),
			Facts:       map[string]string{"variables": strconv.Itoa(vars), "rank": strconv.Itoa(rank), "free": strconv.Itoa(free)},
		}
	default:
		v := rand.Intn(15) - 7
		rowIdx := rand.Intn(3)
		vars := []string{"x", "y", "z"}
		row := fmt.Sprintf("[1  0  0 | %d]", v)
		if rowIdx == 1 {
			row = fmt.Sprintf("[0  1  0 | %d]", v)
		} else if rowIdx == 2 {
			row = fmt.Sprintf("[0  0  1 | %d]", v)
		}
		return generator.Problem{
			Question:    fmt.Sprintf("In reduced row echelon form, row %d of a system in \\((x, y, z)\\) is %s. What is the value of \\(%s\\)? (enter a number)", rowIdx+1, row, vars[rowIdx]),
			Answer:      strconv.Itoa(v),
			Explanation: fmt.Sprintf("The row translates directly to \\(%s = %d\\).", vars[rowIdx], v),
			Facts:       map[string]string{"row": strconv.Itoa(rowIdx + 1), "variable": vars[rowIdx], "value": strconv.Itoa(v)},
		}
	}
}

// ---------------------------------------------------------------------------
// Nonlinear systems (numeric)
// ---------------------------------------------------------------------------

type nonlinearSystemGen struct{}

func (g *nonlinearSystemGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := waveKP(ctx)
	// The system stays y = x^2 and y = kx, whose nonzero solution is x = k.
	// The ordinal draws k from a different band so the three branches do not
	// share one fixed list of six slopes.
	var k int
	switch kp {
	case 2:
		k = rand.Intn(12) + 14 // 14..25
	case 3:
		k = rand.Intn(15) + 26 // 26..40
	default:
		k = rand.Intn(12) + 2 // 2..13
	}
	y := k * k
	return generator.Problem{
		Question: "The system y = x^2 and y = " + strconv.Itoa(k) + "x has two solutions. What is the y-coordinate of the solution with x not equal to 0? (enter a number)",
		Answer:   strconv.Itoa(y),
		Explanation: "Setting x^2 equal to " + strconv.Itoa(k) + "x gives x(x - " + strconv.Itoa(k) +
			") = 0, so the nonzero solution has x = " + strconv.Itoa(k) + " and y = " + strconv.Itoa(k) + "^2 = " + strconv.Itoa(y) + ".",
		Facts: map[string]string{
			"k": strconv.Itoa(k), "x": strconv.Itoa(k), "answer": strconv.Itoa(y),
		},
		KPOrdinal: kp,
	}
}

// ---------------------------------------------------------------------------
// Sigma notation (numeric)
// ---------------------------------------------------------------------------

type sigmaGen struct{}

func (g *sigmaGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := waveKP(ctx)
	// The term stays ai + b. The ordinal chooses a band of upper limits: KP1
	// reads a short sum, KP2 a medium one and KP3 a long one, so the three
	// branches never replay the same six rows.
	var n int
	switch kp {
	case 2:
		n = rand.Intn(5) + 8 // 8..12
	case 3:
		n = rand.Intn(6) + 13 // 13..18
	default:
		n = rand.Intn(5) + 3 // 3..7
	}
	a := rand.Intn(6) + 1 // 1..6
	b := rand.Intn(9) + 1 // 1..9
	tri := n * (n + 1) / 2
	sum := a*tri + b*n
	return generator.Problem{
		Question: "Evaluate the sum from i = 1 to " + strconv.Itoa(n) + " of (" + strconv.Itoa(a) + "i + " +
			strconv.Itoa(b) + "). (enter a number)",
		Answer: strconv.Itoa(sum),
		Explanation: "The sum splits into " + strconv.Itoa(a) + " times the sum of i, which is " + strconv.Itoa(tri) +
			", plus " + strconv.Itoa(b) + " times " + strconv.Itoa(n) + " terms, giving " + strconv.Itoa(sum) + ".",
		Facts: map[string]string{
			"a": strconv.Itoa(a), "b": strconv.Itoa(b), "n": strconv.Itoa(n), "answer": strconv.Itoa(sum),
		},
		KPOrdinal: kp,
	}
}

// ---------------------------------------------------------------------------
// Induction (multiple choice)
// ---------------------------------------------------------------------------

type inductionGen struct{}

func (g *inductionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return inductionProblem()
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
		Question: fmt.Sprintf("The monic polynomial \\(x^3 %+d x^2 %+d x %+d\\) has three integer roots. By the rational root theorem, what is its positive root?", b, d, e),
		Answer:   strconv.Itoa(c.r1),
		Explanation: fmt.Sprintf("The polynomial factors as \\((x %+d)(x %+d)(x %+d)\\), so the roots are %d, %d and %d, and the positive one is %d.",
			-c.r1, -c.r2, -c.r3, c.r1, c.r2, c.r3, c.r1),
		Facts: map[string]string{
			"b": strconv.Itoa(b), "d": strconv.Itoa(d), "e": strconv.Itoa(e),
			"r1": strconv.Itoa(c.r1), "r2": strconv.Itoa(c.r2), "r3": strconv.Itoa(c.r3),
			"answer": strconv.Itoa(c.r1),
		},
	}
}
