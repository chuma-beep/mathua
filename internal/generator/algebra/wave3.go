package algebra

// Wave 2 algebra: eighteen concepts. Every family publishes the values its
// answer follows from, and the tests in wave3_test.go recompute each answer
// from those facts. The three families of a concept share the facts the
// authored solution schema references, so one schema covers all three.

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

// ---------------------------------------------------------------------------
// Shared helpers
// ---------------------------------------------------------------------------

// w3SignChanges counts sign changes in a coefficient sequence, ignoring zeros.
func w3SignChanges(coefs []int) int {
	changes := 0
	prev := 0
	for _, c := range coefs {
		if c == 0 {
			continue
		}
		if prev != 0 && (c > 0) != (prev > 0) {
			changes++
		}
		prev = c
	}
	return changes
}

// w3PolyFromCoefs renders coefficients indexed by degree into a polynomial.
func w3PolyFromCoefs(coefs []int) string {
	terms := make([]polyTerm, 0, len(coefs))
	for d := len(coefs) - 1; d >= 0; d-- {
		terms = append(terms, polyTerm{coefs[d], d})
	}
	return formatPoly(terms...)
}

func w3Abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func w3DivisorCount(n int) int {
	n = w3Abs(n)
	if n == 0 {
		return 0
	}
	count := 0
	for d := 1; d <= n; d++ {
		if n%d == 0 {
			count++
		}
	}
	return count
}

func w3Fib(n int) int {
	a, b := 0, 1
	for i := 0; i < n; i++ {
		a, b = b, a+b
	}
	return a
}

func w3Ordinal(n int) string {
	switch n {
	case 1:
		return "first"
	case 2:
		return "second"
	case 3:
		return "third"
	default:
		return fmt.Sprintf("%dth", n)
	}
}

// w3SignedTerm renders a signed term with a space after the sign and no unit
// coefficient, so a cross term reads "+ 2xy" or "- xy" rather than "+2xy".
func w3SignedTerm(coef int, body string) string {
	mag := w3Abs(coef)
	core := body
	if mag != 1 {
		core = fmt.Sprintf("%d%s", mag, body)
	}
	if coef < 0 {
		return "- " + core
	}
	return "+ " + core
}

// w3RootName names the root that undoes a small integer exponent.
func w3RootName(n int) string {
	switch n {
	case 2:
		return "square"
	case 3:
		return "cube"
	default:
		return fmt.Sprintf("%d-th", n)
	}
}

// ---------------------------------------------------------------------------
// Reciprocal polynomials
// ---------------------------------------------------------------------------

type w3ReciprocalGen struct{}

func (g *w3ReciprocalGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := waveKP(ctx)
	const degree = 4
	coefs := make([]int, degree+1)
	for d := degree; d >= 0; d-- {
		c := rand.Intn(7) - 3
		for c == 0 {
			c = rand.Intn(7) - 3
		}
		coefs[d] = c
	}
	// The leading coefficient is a nonzero multiple of the constant term, so
	// the product of the reciprocal roots is an integer. The first-degree
	// coefficient fixes the sum of the reciprocal roots by Vieta.
	k := rand.Intn(3) + 2 // 2..4
	coefs[degree] = coefs[0] * k
	sum := rand.Intn(7) - 3
	for sum == 0 {
		sum = rand.Intn(7) - 3
	}
	coefs[1] = -sum * coefs[0]
	poly := w3PolyFromCoefs(coefs)

	var question, explanation, answer string
	switch kp {
	case 1:
		answer = strconv.Itoa(coefs[degree])
		question = fmt.Sprintf("The reciprocal polynomial of \\(P(x) = %s\\) is \\(Q(x) = x^{4}P(1/x)\\). What is the constant term of \\(Q\\)? (enter a number)", poly)
		explanation = fmt.Sprintf("Multiplying \\(P(1/x)\\) by \\(x^{4}\\) reverses the coefficients. The constant term of \\(Q\\) is the leading coefficient of \\(P\\), which is %d.", coefs[degree])
	case 2:
		answer = strconv.Itoa(sum)
		question = fmt.Sprintf("For \\(P(x) = %s\\), what is the sum of the reciprocals of its four roots? (enter a number)", poly)
		explanation = fmt.Sprintf("By Vieta, the product of the roots is \\(a_0/a_4\\) because the degree is even, and the sum of the reciprocal roots is \\(-a_1/a_0\\). Here that is \\(-(%d)/%d = %d\\).", coefs[1], coefs[0], sum)
	default:
		answer = strconv.Itoa(k)
		question = fmt.Sprintf("For \\(P(x) = %s\\), what is the product of the reciprocals of its four roots? (enter a number)", poly)
		explanation = fmt.Sprintf("The product of the four roots is \\(a_0/a_4 = %d/%d\\), because the degree is even. The product of the reciprocals is its reciprocal, \\(%d/%d = %d\\).", coefs[0], coefs[degree], coefs[degree], coefs[0], k)
	}

	return generator.Problem{
		Question:    question,
		Answer:      answer,
		Explanation: explanation,
		Facts: map[string]string{
			"poly":                poly,
			"a0":                  strconv.Itoa(coefs[0]),
			"a1":                  strconv.Itoa(coefs[1]),
			"a4":                  strconv.Itoa(coefs[degree]),
			"sum_reciprocals":     strconv.Itoa(sum),
			"product_reciprocals": strconv.Itoa(k),
		},
		KPOrdinal: kp,
	}
}

// ---------------------------------------------------------------------------
// Upper and lower bounds for roots
// ---------------------------------------------------------------------------

type w3PolyBoundsGen struct{}

func (g *w3PolyBoundsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := waveKP(ctx)
	degree := rand.Intn(2) + 3 // 3..4
	coefs := make([]int, degree+1)
	for d := degree; d >= 0; d-- {
		c := rand.Intn(9) - 4 // -4..4
		if c == 0 {
			c = 1 + rand.Intn(4)
		}
		coefs[d] = c
	}
	if kp == 3 {
		coefs[degree] = rand.Intn(4) + 2 // 2..5
	} else {
		coefs[degree] = 1
	}
	poly := w3PolyFromCoefs(coefs)
	leading := coefs[degree]
	constant := coefs[0]
	positiveDivisors := w3DivisorCount(constant)
	candidateTotal := 2 * positiveDivisors * w3DivisorCount(leading)

	maxAbs := 0
	for d := 0; d < degree; d++ {
		if w3Abs(coefs[d]) > maxAbs {
			maxAbs = w3Abs(coefs[d])
		}
	}
	cauchy := fracOrInt(leading+maxAbs, leading)

	var question, explanation, answer string
	switch kp {
	case 1:
		answer = strconv.Itoa(positiveDivisors)
		question = fmt.Sprintf("For \\(P(x) = %s\\), how many positive divisors does the constant term have? Each one is a candidate for a positive integer root. (enter a number)", poly)
		explanation = fmt.Sprintf("The constant term is %d. Its positive divisors number %d, and every integer root of a monic polynomial must divide the constant term.", constant, positiveDivisors)
	case 2:
		answer = cauchy
		question = fmt.Sprintf("For \\(P(x) = %s\\), Cauchy's bound says every real root satisfies \\(|x| < 1 + \\max |a_i|/|a_n|\\). What is that bound? (enter a number)", poly)
		explanation = fmt.Sprintf("The largest coefficient below the leading term has magnitude %d, and the leading coefficient is %d. Cauchy's bound is \\(1 + %d/%d = %s\\), so every real root is smaller than %s in absolute value.", maxAbs, leading, maxAbs, leading, cauchy, cauchy)
	case 3:
		answer = strconv.Itoa(candidateTotal)
		question = fmt.Sprintf("For \\(P(x) = %s\\), how many signed numerator and denominator pairs \\(p/q\\) does the rational root theorem produce before reducing? (enter a number)", poly)
		explanation = fmt.Sprintf("The constant term %d has %d positive divisors and the leading coefficient %d has %d. Each pair gives a signed candidate, so there are \\(2 \\times %d \\times %d = %d\\) candidates.", constant, positiveDivisors, leading, w3DivisorCount(leading), positiveDivisors, w3DivisorCount(leading), candidateTotal)
	}

	return generator.Problem{
		Question:    question,
		Answer:      answer,
		Explanation: explanation,
		Facts: map[string]string{
			"poly":              poly,
			"leading":           strconv.Itoa(leading),
			"constant":          strconv.Itoa(constant),
			"max_coefficient":   strconv.Itoa(maxAbs),
			"positive_divisors": strconv.Itoa(positiveDivisors),
			"cauchy_bound":      cauchy,
			"candidate_total":   strconv.Itoa(candidateTotal),
		},
		KPOrdinal: kp,
	}
}

// ---------------------------------------------------------------------------
// Combining functions
// ---------------------------------------------------------------------------

type w3FuncOperationsGen struct{}

func (g *w3FuncOperationsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := waveKP(ctx)
	a := rand.Intn(5) + 1
	if rand.Intn(2) == 0 {
		a = -a
	}
	c := rand.Intn(5) + 1
	if rand.Intn(2) == 0 {
		c = -c
	}
	b := rand.Intn(9) - 4
	d := rand.Intn(9) - 4
	x := rand.Intn(8) - 3
	fText := formatLine(a, b)
	gText := formatLine(c, d)
	fv := a*x + b
	gv := c*x + d
	sum := fv + gv
	difference := fv - gv
	product := fv * gv

	var question, explanation, answer string
	switch kp {
	case 1:
		answer = strconv.Itoa(sum)
		question = fmt.Sprintf("Let \\(f(x) = %s\\) and \\(g(x) = %s\\). What is \\((f+g)(%d)\\)? (enter a number)", fText, gText, x)
		explanation = fmt.Sprintf("\\(f(%d) = %d\\) and \\(g(%d) = %d\\), so \\((f+g)(%d) = %d + (%d) = %d\\).", x, fv, x, gv, x, fv, gv, sum)
	case 2:
		answer = strconv.Itoa(difference)
		question = fmt.Sprintf("Let \\(f(x) = %s\\) and \\(g(x) = %s\\). What is \\((f-g)(%d)\\)? (enter a number)", fText, gText, x)
		explanation = fmt.Sprintf("\\(f(%d) = %d\\) and \\(g(%d) = %d\\), so \\((f-g)(%d) = %d - (%d) = %d\\).", x, fv, x, gv, x, fv, gv, difference)
	default:
		answer = strconv.Itoa(product)
		question = fmt.Sprintf("Let \\(f(x) = %s\\) and \\(g(x) = %s\\). What is \\((f \\cdot g)(%d)\\)? (enter a number)", fText, gText, x)
		explanation = fmt.Sprintf("\\(f(%d) = %d\\) and \\(g(%d) = %d\\), so \\((f \\cdot g)(%d) = (%d) \\times (%d) = %d\\).", x, fv, x, gv, x, fv, gv, product)
	}

	return generator.Problem{
		Question:    question,
		Answer:      answer,
		Explanation: explanation,
		Facts: map[string]string{
			"f": fText, "g": gText, "x": strconv.Itoa(x),
			"f_value": strconv.Itoa(fv), "g_value": strconv.Itoa(gv),
			"sum": strconv.Itoa(sum), "difference": strconv.Itoa(difference), "product": strconv.Itoa(product),
		},
		KPOrdinal: kp,
	}
}

// ---------------------------------------------------------------------------
// Graphs of inverse functions
// ---------------------------------------------------------------------------

type w3InverseGraphGen struct{}

func (g *w3InverseGraphGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := waveKP(ctx)
	x0 := rand.Intn(9) - 4
	y0 := rand.Intn(9) - 4
	for y0 == x0 {
		y0 = rand.Intn(9) - 4
	}
	shift := y0 - x0
	inverseIntercept := x0 - y0

	var question, explanation, answer string
	switch kp {
	case 1:
		answer = strconv.Itoa(x0)
		question = fmt.Sprintf("The graph of \\(f\\) contains the point \\((%d, %d)\\). What is \\(f^{-1}(%d)\\)? (enter a number)", x0, y0, y0)
		explanation = fmt.Sprintf("The inverse graph swaps the coordinates, so it contains \\((%d, %d)\\). Therefore \\(f^{-1}(%d) = %d\\).", y0, x0, y0, x0)
	case 2:
		answer = strconv.Itoa(y0)
		question = fmt.Sprintf("The graph of \\(f^{-1}\\) contains the point \\((%d, %d)\\). What is \\(f(%d)\\)? (enter a number)", y0, x0, x0)
		explanation = fmt.Sprintf("Swapping the coordinates back puts \\((%d, %d)\\) on the graph of \\(f\\), so \\(f(%d) = %d\\).", x0, y0, x0, y0)
	default:
		answer = strconv.Itoa(inverseIntercept)
		question = fmt.Sprintf("The graph of \\(f\\) is the line of slope 1 through \\((%d, %d)\\). What is the y-intercept of the graph of \\(f^{-1}\\)? (enter a number)", x0, y0)
		explanation = fmt.Sprintf("The line is \\(f(x) = %s\\), because \\(%d - (%d) = %d\\). Its inverse is \\(f^{-1}(x) = %s\\), whose y-intercept is %d.", formatLine(1, shift), y0, x0, shift, formatLine(1, -shift), inverseIntercept)
	}

	return generator.Problem{
		Question:    question,
		Answer:      answer,
		Explanation: explanation,
		Facts: map[string]string{
			"point_x":           strconv.Itoa(x0),
			"point_y":           strconv.Itoa(y0),
			"inverse_value":     strconv.Itoa(x0),
			"original_value":    strconv.Itoa(y0),
			"shift":             strconv.Itoa(shift),
			"inverse_intercept": strconv.Itoa(inverseIntercept),
		},
		KPOrdinal: kp,
	}
}

// ---------------------------------------------------------------------------
// Composition from graphs
// ---------------------------------------------------------------------------

type w3CompositionGraphGen struct{}

func (g *w3CompositionGraphGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := waveKP(ctx)
	a := rand.Intn(3) + 1
	if rand.Intn(2) == 0 {
		a = -a
	}
	c := rand.Intn(3) + 1
	if rand.Intn(2) == 0 {
		c = -c
	}
	b := rand.Intn(7) - 3
	d := rand.Intn(7) - 3
	x := rand.Intn(7) - 3
	fText := formatLine(a, b)
	gText := formatLine(c, d)
	fv := a*x + b
	gv := c*x + d
	fOfG := a*gv + b
	gOfF := c*fv + d
	fOfF := a*fv + b

	var question, explanation, answer string
	switch kp {
	case 1:
		answer = strconv.Itoa(fOfG)
		question = fmt.Sprintf("The functions are \\(f(x) = %s\\) and \\(g(x) = %s\\). What is \\((f \\circ g)(%d)\\)? (enter a number)", fText, gText, x)
		explanation = fmt.Sprintf("Read the inner function first: \\(g(%d) = %d\\). Then \\(f(%d) = %d\\), so \\((f \\circ g)(%d) = %d\\).", x, gv, gv, fOfG, x, fOfG)
	case 2:
		answer = strconv.Itoa(gOfF)
		question = fmt.Sprintf("The functions are \\(f(x) = %s\\) and \\(g(x) = %s\\). What is \\((g \\circ f)(%d)\\)? (enter a number)", fText, gText, x)
		explanation = fmt.Sprintf("Here \\(g\\) is the outer function. First \\(f(%d) = %d\\), then \\(g(%d) = %d\\), so \\((g \\circ f)(%d) = %d\\).", x, fv, fv, gOfF, x, gOfF)
	default:
		answer = strconv.Itoa(fOfF)
		question = fmt.Sprintf("The function is \\(f(x) = %s\\). What is \\((f \\circ f)(%d)\\)? (enter a number)", fText, x)
		explanation = fmt.Sprintf("Apply \\(f\\) twice: \\(f(%d) = %d\\), then \\(f(%d) = %d\\), so \\((f \\circ f)(%d) = %d\\).", x, fv, fv, fOfF, x, fOfF)
	}

	return generator.Problem{
		Question:    question,
		Answer:      answer,
		Explanation: explanation,
		Facts: map[string]string{
			"f": fText, "g": gText, "x": strconv.Itoa(x),
			"f_value": strconv.Itoa(fv), "g_value": strconv.Itoa(gv),
			"f_of_g": strconv.Itoa(fOfG), "g_of_f": strconv.Itoa(gOfF), "f_of_f": strconv.Itoa(fOfF),
		},
		KPOrdinal: kp,
	}
}

// ---------------------------------------------------------------------------
// Exponential regression
// ---------------------------------------------------------------------------

type w3ExpRegressionGen struct{}

func (g *w3ExpRegressionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := waveKP(ctx)
	a := rand.Intn(8) + 2 // 2..9
	b := rand.Intn(2) + 2 // 2..3
	y0 := a
	y1 := a * b
	y2 := a * b * b
	y3 := y2 * b
	y4 := y3 * b
	points := fmt.Sprintf("\\((0, %d)\\), \\((1, %d)\\), \\((2, %d)\\) and \\((3, %d)\\)", y0, y1, y2, y3)

	var question, explanation, answer string
	switch kp {
	case 1:
		answer = strconv.Itoa(b)
		question = fmt.Sprintf("An exponential regression fits \\(y = a b^{x}\\) to the points %s. What is the growth factor \\(b\\)? (enter a number)", points)
		explanation = fmt.Sprintf("Divide consecutive values: \\(%d/%d = %d\\) and \\(%d/%d = %d\\). The constant ratio is the growth factor, so \\(b = %d\\).", y1, y0, b, y2, y1, b, b)
	case 2:
		answer = strconv.Itoa(a)
		question = fmt.Sprintf("An exponential regression fits \\(y = a b^{x}\\) to the points %s. What is the initial value \\(a\\)? (enter a number)", points)
		explanation = fmt.Sprintf("At \\(x = 0\\) the model gives \\(y = a b^{0} = a\\), and the point at \\(x = 0\\) has \\(y = %d\\). So \\(a = %d\\), and the growth factor is \\(b = %d\\).", y0, a, b)
	default:
		answer = strconv.Itoa(y4)
		question = fmt.Sprintf("An exponential regression fits \\(y = a b^{x}\\) to the points %s. What value does it predict at \\(x = 4\\)? (enter a number)", points)
		explanation = fmt.Sprintf("The fit has \\(a = %d\\) and \\(b = %d\\), so at \\(x = 4\\) it predicts \\(%d \\times %d^{4} = %d\\).", a, b, a, b, y4)
	}

	return generator.Problem{
		Question:    question,
		Answer:      answer,
		Explanation: explanation,
		Facts: map[string]string{
			"a": strconv.Itoa(a), "b": strconv.Itoa(b),
			"y0": strconv.Itoa(y0), "y1": strconv.Itoa(y1), "y2": strconv.Itoa(y2),
			"y3": strconv.Itoa(y3), "y4": strconv.Itoa(y4),
		},
		KPOrdinal: kp,
	}
}

// ---------------------------------------------------------------------------
// Domain of logarithmic functions
// ---------------------------------------------------------------------------

type w3LogDomainGen struct{}

func (g *w3LogDomainGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := waveKP(ctx)
	base := rand.Intn(4) + 2 // 2..5
	threshold := rand.Intn(30) + 1

	var question, explanation, answer string
	var endpoint int
	switch kp {
	case 1:
		endpoint = threshold + 1
		answer = strconv.Itoa(endpoint)
		question = fmt.Sprintf("What is the smallest integer in the domain of \\(f(x) = \\log_{%d}(x - %d)\\)? (enter a number)", base, threshold)
		explanation = fmt.Sprintf("The argument must be positive: \\(x - %d > 0\\), so \\(x > %d\\). The smallest integer greater than %d is %d.", threshold, threshold, threshold, endpoint)
	case 2:
		endpoint = threshold - 1
		answer = strconv.Itoa(endpoint)
		question = fmt.Sprintf("What is the largest integer in the domain of \\(f(x) = \\log_{%d}(%d - x)\\)? (enter a number)", base, threshold)
		explanation = fmt.Sprintf("The argument must be positive: \\(%d - x > 0\\), so \\(x < %d\\). The largest integer below %d is %d.", threshold, threshold, threshold, endpoint)
	default:
		root := int(math.Floor(math.Sqrt(float64(threshold))))
		endpoint = root + 1
		answer = strconv.Itoa(endpoint)
		question = fmt.Sprintf("What is the smallest positive integer in the domain of \\(f(x) = \\log_{%d}(x^{2} - %d)\\)? (enter a number)", base, threshold)
		explanation = fmt.Sprintf("The argument must be positive: \\(x^{2} > %d\\). Since \\(%d^{2} = %d\\) is the largest square not exceeding %d, the domain needs \\(x \\ge %d\\), and the smallest positive integer is %d.", threshold, root, root*root, threshold, endpoint, endpoint)
	}

	return generator.Problem{
		Question:    question,
		Answer:      answer,
		Explanation: explanation,
		Facts: map[string]string{
			"base": strconv.Itoa(base), "threshold": strconv.Itoa(threshold), "endpoint": strconv.Itoa(endpoint),
		},
		KPOrdinal: kp,
	}
}

// ---------------------------------------------------------------------------
// Solving for a logarithm's base
// ---------------------------------------------------------------------------

type w3LogBaseEquationGen struct{}

func (g *w3LogBaseEquationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := waveKP(ctx)
	var base, argumentInt int
	var argument, exponent, question, explanation, answer string

	switch kp {
	case 1:
		base = rand.Intn(5) + 2 // 2..6
		e := rand.Intn(2) + 2   // 2..3
		argumentInt = powInt(base, e)
		argument = strconv.Itoa(argumentInt)
		exponent = strconv.Itoa(e)
		question = fmt.Sprintf("Solve \\(\\log_{b}(%d) = %d\\) for the base \\(b\\). (enter a number)", argumentInt, e)
		explanation = fmt.Sprintf("Rewrite the logarithm as a power: \\(b^{%d} = %d\\). Taking the %s root gives \\(b = %d\\).", e, argumentInt, w3RootName(e), base)
	case 2:
		base = rand.Intn(4) + 2 // 2..5
		e := rand.Intn(2) + 2   // 2..3
		denom := powInt(base, e)
		argument = fmt.Sprintf("1/%d", denom)
		exponent = fmt.Sprintf("-%d", e)
		question = fmt.Sprintf("Solve \\(\\log_{b}(1/%d) = -%d\\) for the base \\(b\\). (enter a number)", denom, e)
		explanation = fmt.Sprintf("Rewrite the logarithm as a power: \\(b^{-%d} = 1/%d\\). Taking reciprocals gives \\(b^{%d} = %d\\), so \\(b = %d\\).", e, denom, e, denom, base)
	default:
		argumentInt = rand.Intn(8) + 2 // 2..9
		base = argumentInt * argumentInt
		argument = strconv.Itoa(argumentInt)
		exponent = "1/2"
		question = fmt.Sprintf("Solve \\(\\log_{b}(%d) = 1/2\\) for the base \\(b\\). (enter a number)", argumentInt)
		explanation = fmt.Sprintf("A logarithm of \\(1/2\\) means \\(b^{1/2} = %d\\). Squaring both sides gives \\(b = %d^{2} = %d\\).", argumentInt, argumentInt, base)
	}
	answer = strconv.Itoa(base)

	return generator.Problem{
		Question:    question,
		Answer:      answer,
		Explanation: explanation,
		Facts: map[string]string{
			"base": strconv.Itoa(base), "argument": argument, "exponent": exponent,
		},
		KPOrdinal: kp,
	}
}

// ---------------------------------------------------------------------------
// Arithmetic means between terms
// ---------------------------------------------------------------------------

type w3ArithmeticMeanGen struct{}

func (g *w3ArithmeticMeanGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := waveKP(ctx)
	first := rand.Intn(11) - 5
	difference := rand.Intn(6) + 1

	count, position, answer := 1, 1, 0
	switch kp {
	case 1:
		count = 1
		position = 1
		answer = first + difference
	case 2:
		count = rand.Intn(4) + 1
		position = 1
		answer = difference
	default:
		count = rand.Intn(3) + 2
		position = rand.Intn(count) + 1
		answer = first + position*difference
	}
	last := first + (count+1)*difference
	gap := last - first
	meanValue := first + position*difference
	terms := pluralWord(count, "mean", "means")

	var question, explanation string
	switch kp {
	case 1:
		question = fmt.Sprintf("Insert one arithmetic mean between %d and %d. What is that mean? (enter a number)", first, last)
		explanation = fmt.Sprintf("One mean creates two equal gaps across a span of %d, so the common difference is \\(%d/2 = %d\\). The mean is \\(%d + %d = %d\\).", gap, gap, difference, first, difference, answer)
	case 2:
		question = fmt.Sprintf("Insert %d arithmetic %s between %d and %d. What is the common difference? (enter a number)", count, terms, first, last)
		explanation = fmt.Sprintf("There are \\(%d + 1 = %d\\) equal gaps across a span of %d, so the common difference is \\(%d/%d = %d\\).", count, count+1, gap, gap, count+1, difference)
	default:
		question = fmt.Sprintf("Insert %d arithmetic %s between %d and %d. What is the %s mean? (enter a number)", count, terms, first, last, w3Ordinal(position))
		explanation = fmt.Sprintf("The span is %d across \\(%d + 1 = %d\\) equal gaps, so the common difference is \\(%d/%d = %d\\). The %s mean is \\(%d + %d \\times %d = %d\\).", gap, count, count+1, gap, count+1, difference, w3Ordinal(position), first, position, difference, answer)
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts: map[string]string{
			"first": strconv.Itoa(first), "last": strconv.Itoa(last), "count": strconv.Itoa(count),
			"difference": strconv.Itoa(difference), "position": strconv.Itoa(position), "mean_value": strconv.Itoa(meanValue),
		},
		KPOrdinal: kp,
	}
}

// ---------------------------------------------------------------------------
// Geometric means between terms
// ---------------------------------------------------------------------------

type w3GeometricMeanGen struct{}

func (g *w3GeometricMeanGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := waveKP(ctx)
	first := rand.Intn(4) + 1
	ratio := rand.Intn(2) + 2

	count, position, answer := 1, 1, 0
	switch kp {
	case 1:
		count = 1
		position = 1
		answer = first * ratio
	case 2:
		count = rand.Intn(3) + 1
		position = 1
		answer = ratio
	default:
		count = rand.Intn(3) + 2
		position = rand.Intn(count) + 1
		answer = first * powInt(ratio, position)
	}
	last := first * powInt(ratio, count+1)
	meanValue := first * powInt(ratio, position)
	terms := pluralWord(count, "mean", "means")

	var question, explanation string
	switch kp {
	case 1:
		question = fmt.Sprintf("Insert one geometric mean between %d and %d. What is that mean? (enter a number)", first, last)
		explanation = fmt.Sprintf("One mean creates two equal ratios, so the ratio is \\(\\sqrt{%d/%d} = %d\\). The mean is \\(%d \\times %d = %d\\).", last, first, ratio, first, ratio, answer)
	case 2:
		question = fmt.Sprintf("Insert %d geometric %s between %d and %d. What is the common ratio? (enter a number)", count, terms, first, last)
		explanation = fmt.Sprintf("There are \\(%d + 1 = %d\\) equal ratios between %d and %d, so the common ratio is \\((%d/%d)^{1/%d} = %d\\).", count, count+1, first, last, last, first, count+1, ratio)
	default:
		question = fmt.Sprintf("Insert %d geometric %s between %d and %d. What is the %s mean? (enter a number)", count, terms, first, last, w3Ordinal(position))
		explanation = fmt.Sprintf("The common ratio is \\((%d/%d)^{1/%d} = %d\\). The %s mean is \\(%d \\times %d^{%d} = %d\\).", last, first, count+1, ratio, w3Ordinal(position), first, ratio, position, answer)
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts: map[string]string{
			"first": strconv.Itoa(first), "last": strconv.Itoa(last), "count": strconv.Itoa(count),
			"ratio": strconv.Itoa(ratio), "position": strconv.Itoa(position), "mean_value": strconv.Itoa(meanValue),
		},
		KPOrdinal: kp,
	}
}

// ---------------------------------------------------------------------------
// Fibonacci and Binet's formula
// ---------------------------------------------------------------------------

type w3FibonacciGen struct{}

func (g *w3FibonacciGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := waveKP(ctx)
	n := rand.Intn(10) + 6 // 6..15
	if kp == 3 {
		n = rand.Intn(9) + 5 // 5..13
	}
	fn := w3Fib(n)
	fNext := w3Fib(n + 1)
	fNextNext := w3Fib(n + 2)
	sum := fNextNext - 1

	var question, explanation, answer string
	switch kp {
	case 1:
		answer = strconv.Itoa(fn)
		question = fmt.Sprintf("Using the Fibonacci recurrence \\(F_1 = F_2 = 1\\) and \\(F_k = F_{k-1} + F_{k-2}\\), what is \\(F_{%d}\\)? (enter a number)", n)
		explanation = fmt.Sprintf("Stepping the recurrence from \\(F_1 = 1\\) and \\(F_2 = 1\\) gives \\(F_{%d} = %d\\). Binet's formula produces the same integer.", n, fn)
	case 2:
		answer = strconv.Itoa(fNextNext)
		question = fmt.Sprintf("Two consecutive Fibonacci numbers are \\(F_{%d} = %d\\) and \\(F_{%d} = %d\\). What is \\(F_{%d}\\)? (enter a number)", n, fn, n+1, fNext, n+2)
		explanation = fmt.Sprintf("The recurrence adds the two previous terms: \\(F_{%d} = F_{%d} + F_{%d} = %d + %d = %d\\).", n+2, n+1, n, fNext, fn, fNextNext)
	default:
		answer = strconv.Itoa(sum)
		question = fmt.Sprintf("What is the sum \\(F_1 + F_2 + \\cdots + F_{%d}\\) of the first %d Fibonacci numbers? (enter a number)", n, n)
		explanation = fmt.Sprintf("Use \\(F_1 + F_2 + \\cdots + F_n = F_{n+2} - 1\\). With \\(n = %d\\), \\(F_{%d} = %d\\), so the sum is \\(%d - 1 = %d\\).", n, n+2, fNextNext, fNextNext, sum)
	}

	return generator.Problem{
		Question:    question,
		Answer:      answer,
		Explanation: explanation,
		Facts: map[string]string{
			"n": strconv.Itoa(n), "f_n": strconv.Itoa(fn), "f_next": strconv.Itoa(fNext), "sum": strconv.Itoa(sum),
		},
		KPOrdinal: kp,
	}
}

// ---------------------------------------------------------------------------
// Line and conic intersections
// ---------------------------------------------------------------------------

type w3ConicIntersectionGen struct{}

func (g *w3ConicIntersectionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := waveKP(ctx)
	m := rand.Intn(9) - 4
	for m == 0 {
		m = rand.Intn(9) - 4
	}
	b := rand.Intn(9) - 4
	switch kp {
	case 1:
		if rand.Intn(3) == 0 {
			m = 2 * (rand.Intn(3) - 1) // -2, 0 or 2
			if m == 0 {
				m = 2
			}
			b = -m * m / 4 // tangent, so the discriminant is zero
		}
	case 3:
		for m*m+4*b <= 0 {
			b++
		}
	}
	d := m*m + 4*b
	count := 0
	if d > 0 {
		count = 2
	} else if d == 0 {
		count = 1
	}
	line := formatLinear(m, b)

	var question, explanation, answer string
	switch kp {
	case 1:
		answer = strconv.Itoa(count)
		question = fmt.Sprintf("How many points do the parabola \\(y = x^{2}\\) and the line \\(%s\\) have in common? (enter a number)", line)
		switch count {
		case 0:
			explanation = fmt.Sprintf("Substitute \\(y = %s\\) into \\(y = x^{2}\\) to get \\(%s = 0\\). Its discriminant is \\((%d)^{2} + 4(%d) = %d < 0\\), so there are 0 real intersection points.", line[4:], formatQuadratic(-m, -b), m, b, d)
		case 1:
			explanation = fmt.Sprintf("Substitute \\(y = %s\\) into \\(y = x^{2}\\) to get \\(%s = 0\\). Its discriminant is \\((%d)^{2} + 4(%d) = 0\\), so the line is tangent and there is 1 intersection point.", line[4:], formatQuadratic(-m, -b), m, b)
		default:
			explanation = fmt.Sprintf("Substitute \\(y = %s\\) into \\(y = x^{2}\\) to get \\(%s = 0\\). Its discriminant is \\((%d)^{2} + 4(%d) = %d > 0\\), so there are 2 intersection points.", line[4:], formatQuadratic(-m, -b), m, b, d)
		}
	case 2:
		answer = strconv.Itoa(d)
		question = fmt.Sprintf("The line \\(%s\\) meets the parabola \\(y = x^{2}\\) where \\(%s = 0\\). What is the discriminant of that quadratic? (enter a number)", line, formatQuadratic(-m, -b))
		explanation = fmt.Sprintf("The quadratic is \\(%s\\). Its discriminant is \\(m^{2} + 4b = (%d)^{2} + 4(%d) = %d\\).", formatQuadratic(-m, -b), m, b, d)
	default:
		answer = strconv.Itoa(m)
		question = fmt.Sprintf("The line \\(%s\\) meets the parabola \\(y = x^{2}\\) in two points. What is the sum of their \\(x\\)-coordinates? (enter a number)", line)
		explanation = fmt.Sprintf("The \\(x\\)-coordinates satisfy \\(%s = 0\\). By Vieta, their sum is the coefficient of \\(x\\) with the opposite sign, which is %d.", formatQuadratic(-m, -b), m)
	}

	return generator.Problem{
		Question:    question,
		Answer:      answer,
		Explanation: explanation,
		Facts: map[string]string{
			"m": strconv.Itoa(m), "b": strconv.Itoa(b),
			"discriminant": strconv.Itoa(d), "count": strconv.Itoa(count), "sum_x": strconv.Itoa(m),
		},
		KPOrdinal: kp,
	}
}

// ---------------------------------------------------------------------------
// Translated conic equations
// ---------------------------------------------------------------------------

func w3ShiftedSquare(variable string, shift int) string {
	if shift > 0 {
		return fmt.Sprintf("(%s - %d)^{2}", variable, shift)
	}
	return fmt.Sprintf("(%s + %d)^{2}", variable, -shift)
}

type w3ConicTranslationGen struct{}

func (g *w3ConicTranslationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := waveKP(ctx)
	h := rand.Intn(11) - 5
	for h == 0 {
		h = rand.Intn(11) - 5
	}
	k := rand.Intn(11) - 5
	for k == 0 {
		k = rand.Intn(11) - 5
	}
	a := rand.Intn(5) + 2 // 2..6
	b := rand.Intn(5) + 2
	equation := fmt.Sprintf("\\frac{%s}{%d} + \\frac{%s}{%d} = 1", w3ShiftedSquare("x", h), a*a, w3ShiftedSquare("y", k), b*b)

	var question, explanation, answer string
	switch kp {
	case 1:
		answer = strconv.Itoa(h)
		question = fmt.Sprintf("The ellipse \\(%s\\) has center \\((h, k)\\). What is \\(h\\)? (enter a number)", equation)
		explanation = fmt.Sprintf("The term \\((x - h)^{2}\\) vanishes when \\(x = h\\). The equation shows \\(%s\\), so \\(h = %d\\).", w3ShiftedSquare("x", h), h)
	case 2:
		answer = strconv.Itoa(k)
		question = fmt.Sprintf("The ellipse \\(%s\\) has center \\((h, k)\\). What is \\(k\\)? (enter a number)", equation)
		explanation = fmt.Sprintf("The term \\((y - k)^{2}\\) vanishes when \\(y = k\\). The equation shows \\(%s\\), so \\(k = %d\\).", w3ShiftedSquare("y", k), k)
	default:
		answer = strconv.Itoa(h + k)
		question = fmt.Sprintf("The ellipse \\(%s\\) has center \\((h, k)\\). What is \\(h + k\\)? (enter a number)", equation)
		explanation = fmt.Sprintf("The center is \\((%d, %d)\\), so \\(h + k = (%d) + (%d) = %d\\).", h, k, h, k, h+k)
	}

	return generator.Problem{
		Question:    question,
		Answer:      answer,
		Explanation: explanation,
		Facts: map[string]string{
			"h": strconv.Itoa(h), "k": strconv.Itoa(k),
			"a": strconv.Itoa(a), "b": strconv.Itoa(b), "equation": equation,
		},
		KPOrdinal: kp,
	}
}

// ---------------------------------------------------------------------------
// Tangent lines to conics
// ---------------------------------------------------------------------------

type w3ConicTangentGen struct{}

func (g *w3ConicTangentGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := waveKP(ctx)
	var curve string
	var px, py int
	var slope, intercept, answer string
	param, pval := 0, 0

	switch kp {
	case 1:
		type point struct{ x, y int }
		bases := []point{{3, 4}, {4, 3}, {5, 12}, {12, 5}, {8, 15}, {15, 8}, {7, 24}, {24, 7}}
		p := bases[rand.Intn(len(bases))]
		sx, sy := 1, 1
		if rand.Intn(2) == 0 {
			sx = -1
		}
		if rand.Intn(2) == 0 {
			sy = -1
		}
		px, py = sx*p.x, sy*p.y
		curve = "circle"
		slope = fracOrInt(-px, py)
		intercept = fracOrInt(px*px+py*py, px)
		answer = slope
	case 2:
		a := rand.Intn(8) - 4
		if a == 0 {
			a = 1
		}
		px, py = a, a*a
		curve = "parabola"
		slope = strconv.Itoa(2 * a)
		intercept = strconv.Itoa(-a * a)
		answer = slope
	default:
		p := rand.Intn(3) + 1
		t := rand.Intn(4) + 1
		pval, param = p, t
		px, py = p*t*t, 2*p*t
		curve = "parabola"
		slope = fracOrInt(1, t)
		intercept = strconv.Itoa(-p * t * t)
		answer = intercept
	}

	var question, explanation string
	switch kp {
	case 1:
		question = fmt.Sprintf("The circle \\(x^{2} + y^{2} = %d\\) contains the point \\((%d, %d)\\). What is the slope of the tangent line at that point? (enter a number or fraction)", px*px+py*py, px, py)
		explanation = fmt.Sprintf("The tangent to \\(x^{2} + y^{2} = r^{2}\\) at \\((x_0, y_0)\\) has slope \\(-x_0/y_0\\). Here the slope is %s. The tangent line is \\(%s\\).", slope, formatEquation(px, "x", py, "y", px*px+py*py))
	case 2:
		question = fmt.Sprintf("The parabola \\(y = x^{2}\\) contains the point \\((%d, %d)\\). What is the slope of the tangent line at that point? (enter a number)", px, py)
		explanation = fmt.Sprintf("For \\(y = x^{2}\\), the derivative is \\(2x\\). At \\(x = %d\\) the slope is \\(2 \\times %d = %s\\), and the tangent is \\(y = %s\\).", px, px, slope, formatLine(2*px, -px*px))
	default:
		question = fmt.Sprintf("The parabola \\(y^{2} = 4px\\) contains the point \\((%d, %d)\\). What is the x-intercept of the tangent line at that point? (enter a number)", px, py)
		explanation = fmt.Sprintf("At the parameter \\(t\\) the point is \\((pt^{2}, 2pt)\\), so here \\(t = %d\\) and \\(p = %d\\). The tangent \\(ty = x + pt^{2}\\) meets the x-axis at \\(x = -pt^{2} = %s\\).", param, pval, intercept)
	}

	return generator.Problem{
		Question:    question,
		Answer:      answer,
		Explanation: explanation,
		Facts: map[string]string{
			"curve": curve, "point_x": strconv.Itoa(px), "point_y": strconv.Itoa(py),
			"slope": slope, "intercept": intercept,
		},
		KPOrdinal: kp,
	}
}

// ---------------------------------------------------------------------------
// Solving systems with inverse matrices
// ---------------------------------------------------------------------------

type w3MatrixInverseGen struct{}

func (g *w3MatrixInverseGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := waveKP(ctx)
	var a, b, c, d int
	for {
		a = rand.Intn(7) - 3
		b = rand.Intn(7) - 3
		c = rand.Intn(7) - 3
		d = rand.Intn(7) - 3
		if a*d-b*c != 0 && a != 0 && c != 0 {
			break
		}
	}
	x0 := rand.Intn(9) - 4
	y0 := rand.Intn(9) - 4
	e := a*x0 + b*y0
	f := c*x0 + d*y0
	det := a*d - b*c

	var question, explanation, answer string
	switch kp {
	case 1:
		answer = strconv.Itoa(det)
		question = fmt.Sprintf("The system is \\(%s\\) and \\(%s\\). What is the determinant of its coefficient matrix? (enter a number)", formatEquation(a, "x", b, "y", e), formatEquation(c, "x", d, "y", f))
		explanation = fmt.Sprintf("The coefficient matrix is \\(\\begin{pmatrix}%d & %d\\\\ %d & %d\\end{pmatrix}\\), so its determinant is \\(( %d )( %d ) - ( %d )( %d ) = %d\\).", a, b, c, d, a, d, b, c, det)
	case 2:
		answer = strconv.Itoa(x0)
		question = fmt.Sprintf("The system is \\(%s\\) and \\(%s\\). Using the inverse matrix, what is \\(x\\)? (enter a number)", formatEquation(a, "x", b, "y", e), formatEquation(c, "x", d, "y", f))
		explanation = fmt.Sprintf("The determinant is %d, so \\(A^{-1} = \\frac{1}{%d}\\begin{pmatrix}%d & %d\\\\ %d & %d\\end{pmatrix}\\). Then \\(x = ( %d )( %d ) - ( %d )( %d )\\) divided by %d, which gives %d.", det, det, d, -b, -c, a, d, e, b, f, det, x0)
	default:
		answer = strconv.Itoa(y0)
		question = fmt.Sprintf("The system is \\(%s\\) and \\(%s\\). Using the inverse matrix, what is \\(y\\)? (enter a number)", formatEquation(a, "x", b, "y", e), formatEquation(c, "x", d, "y", f))
		explanation = fmt.Sprintf("The determinant is %d, so \\(A^{-1} = \\frac{1}{%d}\\begin{pmatrix}%d & %d\\\\ %d & %d\\end{pmatrix}\\). Then \\(y = ( %d )( %d ) - ( %d )( %d )\\) divided by %d, which gives %d.", det, det, d, -b, -c, a, a, f, c, e, det, y0)
	}

	return generator.Problem{
		Question:    question,
		Answer:      answer,
		Explanation: explanation,
		Facts: map[string]string{
			"a": strconv.Itoa(a), "b": strconv.Itoa(b), "c": strconv.Itoa(c), "d": strconv.Itoa(d),
			"e": strconv.Itoa(e), "f": strconv.Itoa(f), "det": strconv.Itoa(det),
			"x": strconv.Itoa(x0), "y": strconv.Itoa(y0),
		},
		KPOrdinal: kp,
	}
}

// ---------------------------------------------------------------------------
// Quadratic regression
// ---------------------------------------------------------------------------

type w3QuadraticRegressionGen struct{}

func (g *w3QuadraticRegressionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := waveKP(ctx)
	a := rand.Intn(3) + 1 // 1..3
	b := rand.Intn(5) - 2 // -2..2
	c := rand.Intn(7) - 3 // -3..3
	xs := []int{-2, -1, 0, 1, 2}
	ys := make([]int, len(xs))
	for i, x := range xs {
		ys[i] = a*x*x + b*x + c
	}
	points := ""
	for i := range xs {
		if i > 0 {
			points += ", "
		}
		points += fmt.Sprintf("\\((%d, %d)\\)", xs[i], ys[i])
	}
	prediction := 9*a + 3*b + c
	fit := formatPoly(polyTerm{a, 2}, polyTerm{b, 1}, polyTerm{c, 0})

	var question, explanation, answer string
	switch kp {
	case 1:
		answer = strconv.Itoa(a)
		question = fmt.Sprintf("A quadratic regression fits \\(y = ax^{2} + bx + c\\) to the points %s. What is \\(a\\)? (enter a number)", points)
		explanation = fmt.Sprintf("The points lie exactly on \\(y = %s\\). Substituting them into the normal equations gives \\(a = %d\\), \\(b = %d\\) and \\(c = %d\\).", fit, a, b, c)
	case 2:
		answer = strconv.Itoa(c)
		question = fmt.Sprintf("A quadratic regression fits \\(y = ax^{2} + bx + c\\) to the points %s. What is \\(c\\)? (enter a number)", points)
		explanation = fmt.Sprintf("The point at \\(x = 0\\) is \\((0, %d)\\), so the constant term is \\(c = %d\\). The full fit is \\(y = %s\\).", ys[2], c, fit)
	default:
		answer = strconv.Itoa(prediction)
		question = fmt.Sprintf("A quadratic regression fits \\(y = ax^{2} + bx + c\\) to the points %s. What is the fitted value at \\(x = 3\\)? (enter a number)", points)
		explanation = fmt.Sprintf("The fit is \\(y = %s\\). At \\(x = 3\\) it gives \\(%d \\times 9 + (%d) \\times 3 + (%d) = %d\\).", fit, a, b, c, prediction)
	}

	return generator.Problem{
		Question:    question,
		Answer:      answer,
		Explanation: explanation,
		Facts: map[string]string{
			"a": strconv.Itoa(a), "b": strconv.Itoa(b), "c": strconv.Itoa(c),
			"prediction": strconv.Itoa(prediction), "points": points,
		},
		KPOrdinal: kp,
	}
}

// ---------------------------------------------------------------------------
// Power regression
// ---------------------------------------------------------------------------

type w3PowerRegressionGen struct{}

func (g *w3PowerRegressionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := waveKP(ctx)
	a := rand.Intn(6) + 2 // 2..7
	b := rand.Intn(2) + 2 // 2..3
	y1 := a
	y2 := a * powInt(2, b)
	y3 := a * powInt(3, b)
	prediction := a * powInt(4, b)
	points := fmt.Sprintf("\\((1, %d)\\), \\((2, %d)\\) and \\((3, %d)\\)", y1, y2, y3)

	var question, explanation, answer string
	switch kp {
	case 1:
		answer = strconv.Itoa(b)
		question = fmt.Sprintf("A power regression fits \\(y = a x^{b}\\) to the points %s. What is the exponent \\(b\\)? (enter a number)", points)
		explanation = fmt.Sprintf("Taking logarithms gives \\(\\ln y = \\ln a + b \\ln x\\). The ratio \\(%d/%d = %d\\) is \\(2^{b}\\), and \\(%d/%d = %d\\) is \\(3^{b}\\), so \\(b = %d\\).", y2, y1, y2/y1, y3, y1, y3/y1, b)
	case 2:
		answer = strconv.Itoa(a)
		question = fmt.Sprintf("A power regression fits \\(y = a x^{b}\\) to the points %s. What is the leading coefficient \\(a\\)? (enter a number)", points)
		explanation = fmt.Sprintf("At \\(x = 1\\) the model gives \\(y = a\\), and the first point has \\(y = %d\\). So \\(a = %d\\), and the exponent is \\(b = %d\\).", y1, a, b)
	default:
		answer = strconv.Itoa(prediction)
		question = fmt.Sprintf("A power regression fits \\(y = a x^{b}\\) to the points %s. What is the fitted value at \\(x = 4\\)? (enter a number)", points)
		explanation = fmt.Sprintf("The fit is \\(y = %d x^{%d}\\). At \\(x = 4\\) it gives \\(%d \\times 4^{%d} = %d\\).", a, b, a, b, prediction)
	}

	return generator.Problem{
		Question:    question,
		Answer:      answer,
		Explanation: explanation,
		Facts: map[string]string{
			"a": strconv.Itoa(a), "b": strconv.Itoa(b),
			"y1": strconv.Itoa(y1), "y2": strconv.Itoa(y2), "y3": strconv.Itoa(y3),
			"prediction": strconv.Itoa(prediction), "points": points,
		},
		KPOrdinal: kp,
	}
}

// ---------------------------------------------------------------------------
// Logarithmic regression
// ---------------------------------------------------------------------------

type w3LogarithmicRegressionGen struct{}

func (g *w3LogarithmicRegressionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := waveKP(ctx)
	a := rand.Intn(8) + 1 // 1..8
	b := rand.Intn(5) + 1 // 1..5
	y1 := a
	y2 := a + b
	y3 := a + 2*b
	prediction := a + 3*b
	points := fmt.Sprintf("\\((1, %d)\\), \\((e, %d)\\) and \\((e^{2}, %d)\\)", y1, y2, y3)

	var question, explanation, answer string
	switch kp {
	case 1:
		answer = strconv.Itoa(b)
		question = fmt.Sprintf("A logarithmic regression fits \\(y = a + b \\ln x\\) to the points %s. What is \\(b\\)? (enter a number)", points)
		explanation = fmt.Sprintf("Since \\(\\ln 1 = 0\\), \\(\\ln e = 1\\) and \\(\\ln e^{2} = 2\\), the fitted values increase by \\(b\\) for each unit of \\(\\ln x\\). From %d to %d the increase is %d, so \\(b = %d\\).", y1, y2, b, b)
	case 2:
		answer = strconv.Itoa(a)
		question = fmt.Sprintf("A logarithmic regression fits \\(y = a + b \\ln x\\) to the points %s. What is the intercept \\(a\\)? (enter a number)", points)
		explanation = fmt.Sprintf("At \\(x = 1\\), \\(\\ln x = 0\\), so the model gives \\(y = a\\). The point at \\(x = 1\\) has \\(y = %d\\), so \\(a = %d\\). The coefficient is \\(b = %d\\).", y1, a, b)
	default:
		answer = strconv.Itoa(prediction)
		question = fmt.Sprintf("A logarithmic regression fits \\(y = a + b \\ln x\\) to the points %s. What is the fitted value at \\(x = e^{3}\\)? (enter a number)", points)
		explanation = fmt.Sprintf("The fit is \\(y = %d + %d \\ln x\\). At \\(x = e^{3}\\), \\(\\ln x = 3\\), so the prediction is \\(%d + %d \\times 3 = %d\\).", a, b, a, b, prediction)
	}

	return generator.Problem{
		Question:    question,
		Answer:      answer,
		Explanation: explanation,
		Facts: map[string]string{
			"a": strconv.Itoa(a), "b": strconv.Itoa(b),
			"y1": strconv.Itoa(y1), "y2": strconv.Itoa(y2), "y3": strconv.Itoa(y3),
			"prediction": strconv.Itoa(prediction), "points": points,
		},
		KPOrdinal: kp,
	}
}
