package numerical

// The first numerical-analysis batch: twenty concepts that open the domain with
// error and conditioning, root finding and the first half of interpolation.
//
// The domain is about the gap between an exact mathematical object and the
// finitely many digits a machine can hold and the finitely many steps a method
// can afford. This batch measures that gap (absolute and relative error,
// roundoff, cancellation, conditioning, stability, truncation), closes it by
// iteration (bisection, fixed-point, Newton, secant, false position and the
// orders those methods achieve) and starts reconstructing a function from
// samples (Lagrange interpolation, divided differences, the error formula,
// Chebyshev nodes and splines).
//
// Nine concepts are multiple choice and eleven are numeric. Every numeric
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
// Absolute and relative error (numeric)
// ---------------------------------------------------------------------------

type absRelErrorGen struct{}

func (g *absRelErrorGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	maxTrue := mathutil.ScaleInt(ctx.Difficulty, 500, 4000)
	tVal := (1 + rand.Intn(maxTrue/100)) * 100
	maxDelta := mathutil.ScaleInt(ctx.Difficulty, 4, 25)
	delta := 1 + rand.Intn(maxDelta)
	approx := tVal + delta
	if rand.Intn(2) == 1 {
		approx = tVal - delta
	}

	relErr := float64(delta) / float64(tVal)
	ans := decimal6(relErr)

	q := fmt.Sprintf("A quantity has true value %d and is approximated by %d. What is the relative error |approx − true|/|true|? (enter a decimal with a decimal point, rounded to 6 decimal places)", tVal, approx)
	exp := fmt.Sprintf("Subgoal 1: State governing formula. The relative error measures the magnitude of approximation error relative to the true quantity: \\(\\text{RelErr} = \\frac{|\\text{approx} - \\text{true}|}{|\\text{true}|}\\).\n\nSubgoal 2: Substitute parameters. Given true value %d and approximation %d, substitute into the formula: \\(\\text{RelErr} = \\frac{|%d - %d|}{|%d|}\\).\n\nSubgoal 3: Evaluate and simplify. Compute the absolute difference \\(|%d - %d| = %d\\) and divide by %d: \\(\\text{RelErr} = \\frac{%d}{%d} = %s\\).",
		tVal, approx, approx, tVal, tVal,
		approx, tVal, delta, tVal, delta, tVal, ans)

	return generator.Problem{
		Question:    q,
		Answer:      ans,
		Explanation: exp,
		Facts: map[string]string{
			"true":   strconv.Itoa(tVal),
			"approx": strconv.Itoa(approx),
			"answer": ans,
		},
	}
}

// ---------------------------------------------------------------------------
// Floating-point roundoff and machine epsilon (multiple choice)
// ---------------------------------------------------------------------------

type roundoffGen struct{}

func (g *roundoffGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return roundoffProblem()
}

// ---------------------------------------------------------------------------
// Catastrophic cancellation (numeric)
// ---------------------------------------------------------------------------

type cancellationGen struct{}

func (g *cancellationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	maxK := mathutil.ScaleInt(ctx.Difficulty, 10, 200)
	k := 5 + rand.Intn(maxK-4)
	d := 2 + rand.Intn(10)
	if ((k+1)*d)%2 != 0 {
		d = 2 * d
	}
	x := (k + 1) * d / 2
	y := x - d

	for iter := 0; iter < 5 && (k == x || k == y); iter++ {
		k += 10
		x = (k + 1) * d / 2
		y = x - d
	}

	ans := decimal6(float64(x+y) / float64(d))
	q := fmt.Sprintf("The subtraction x − y loses accuracy when x and y are close: the relative error of the difference is amplified by (|x| + |y|)/|x − y| relative to the inputs. For x = %d and y = %d, what is that amplification factor? (enter a decimal)", x, y)
	exp := fmt.Sprintf("Subgoal 1: State governing formula. Subtraction of nearby quantities suffers from catastrophic cancellation, amplifying input errors by the factor \\(\\kappa = \\frac{|x| + |y|}{|x - y|}\\).\n\nSubgoal 2: Substitute parameters. Substitute \\(x = %d\\) and \\(y = %d\\) into the amplification ratio: \\(\\kappa = \\frac{|%d| + |%d|}{|%d - %d|}\\).\n\nSubgoal 3: Evaluate and simplify. Compute the numerator \\(%d + %d = %d\\) and denominator \\(|%d| = %d\\), then divide: \\(\\kappa = \\frac{%d}{%d} = %s\\).",
		x, y, x, y, x, y,
		x, y, x+y, x-y, d, x+y, d, ans)

	return generator.Problem{
		Question:    q,
		Answer:      ans,
		Explanation: exp,
		Facts: map[string]string{
			"x":      strconv.Itoa(x),
			"y":      strconv.Itoa(y),
			"answer": ans,
		},
	}
}

// ---------------------------------------------------------------------------
// The condition number of a problem (numeric)
// ---------------------------------------------------------------------------

type conditioningGen struct{}

func (g *conditioningGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	maxN := mathutil.ScaleInt(ctx.Difficulty, 6, 25)
	n := 2 + rand.Intn(maxN-1)
	ans := decimal6(float64(n))

	q := fmt.Sprintf("The relative condition number of the function x raised to the power n at a point x is |x·f'(x)|/|f(x)|. For f(x) = x^%d, what is that condition number? (enter a decimal)", n)
	exp := fmt.Sprintf("Subgoal 1: State governing formula. The relative condition number of a differentiable function \\(f(x)\\) measures output sensitivity to input perturbation: \\(\\kappa(x) = \\frac{|x f'(x)|}{|f(x)|}\\).\n\nSubgoal 2: Substitute parameters. For monomial \\(f(x) = x^{%d}\\), compute derivative \\(f'(x) = %d x^{%d}\\) and substitute: \\(\\kappa(x) = \\frac{|x \\cdot %d x^{%d}|}{|x^{%d}|}\\).\n\nSubgoal 3: Evaluate and simplify. Multiply in the numerator to obtain \\(%d |x^{%d}|\\) and divide out \\(|x^{%d}|\\): \\(\\kappa(x) = %d\\).",
		n, n, n-1, n, n-1, n,
		n, n, n, n)

	return generator.Problem{
		Question:    q,
		Answer:      ans,
		Explanation: exp,
		Facts: map[string]string{
			"n":      strconv.Itoa(n),
			"answer": ans,
		},
	}
}

// ---------------------------------------------------------------------------
// Stability of an algorithm (multiple choice)
// ---------------------------------------------------------------------------

type stabilityGen struct{}

func (g *stabilityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return stabilityProblem()
}

// ---------------------------------------------------------------------------
// Taylor truncation error (multiple choice)
// ---------------------------------------------------------------------------

type taylorErrorGen struct{}

func (g *taylorErrorGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return taylorErrorProblem()
}

// ---------------------------------------------------------------------------
// The bisection method (numeric)
// ---------------------------------------------------------------------------

type bisectionGen struct{}

func (g *bisectionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	widths := []float64{1, 2, 4, 8, 16, 32}
	w := widths[rand.Intn(len(widths))]
	wStr := strconv.Itoa(int(w))

	maxK := mathutil.ScaleInt(ctx.Difficulty, 6, 16)
	k := 4 + rand.Intn(maxK-3)
	tol := w / math.Pow(2, float64(k))
	tolStr := strconv.FormatFloat(tol, 'f', -1, 64)

	steps := int(math.Ceil(math.Log2(w / tol)))
	ans := strconv.Itoa(steps)

	for iter := 0; iter < 5 && ans == wStr; iter++ {
		k++
		tol = w / math.Pow(2, float64(k))
		tolStr = strconv.FormatFloat(tol, 'f', -1, 64)
		steps = int(math.Ceil(math.Log2(w / tol)))
		ans = strconv.Itoa(steps)
	}

	q := fmt.Sprintf("Bisection halves the interval at each step. Starting from an interval of width %s, how many steps are needed to make the width at most %s? (enter a whole number)", wStr, tolStr)
	exp := fmt.Sprintf("Subgoal 1: State governing formula. The bisection method halves the search interval at each iteration, giving width \\(W_k = \\frac{W_0}{2^k}\\) after k steps. Requiring \\(W_k \\le \\epsilon\\) implies \\(2^k \\ge \\frac{W_0}{\\epsilon}\\), so \\(k = \\lceil \\log_2(W_0 / \\epsilon) \\rceil\\).\n\nSubgoal 2: Substitute parameters. Substitute initial width \\(W_0 = %s\\) and tolerance \\(\\epsilon = %s\\): \\(\\frac{W_0}{\\epsilon} = \\frac{%s}{%s} = %s\\).\n\nSubgoal 3: Evaluate and simplify. Take the base-two logarithm \\(\\log_2(%s) = %d\\), which establishes that exactly %d steps are required.",
		wStr, tolStr, wStr, tolStr, strconv.FormatFloat(w/tol, 'f', -1, 64),
		strconv.FormatFloat(w/tol, 'f', -1, 64), steps, steps)

	return generator.Problem{
		Question:    q,
		Answer:      ans,
		Explanation: exp,
		Facts: map[string]string{
			"width":  wStr,
			"tol":    tolStr,
			"answer": ans,
		},
	}
}

// ---------------------------------------------------------------------------
// Fixed-point iteration (numeric)
// ---------------------------------------------------------------------------

type fixedPointGen struct{}

func (g *fixedPointGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	maxA := mathutil.ScaleInt(ctx.Difficulty, 3, 10)
	a := 1 + rand.Intn(maxA)
	bd := 2 + rand.Intn(3)
	bn := 1
	if bd > 2 && rand.Intn(2) == 1 {
		bn = bd - 1
	}
	x0 := rand.Intn(4)
	n := 2
	if ctx.Difficulty >= 0.5 {
		n = 3
	}

	x := float64(x0)
	b := float64(bn) / float64(bd)
	var steps []string
	cur := x
	for i := 0; i < n; i++ {
		cur = float64(a) + b*cur
		steps = append(steps, fmt.Sprintf("x_%d = %d + (%d/%d)(%s) = %s", i+1, a, bn, bd, decimal6(x), decimal6(cur)))
		x = cur
	}
	ans := decimal6(x)

	q := fmt.Sprintf("Fixed-point iteration applies x_{k+1} = a + b·x_k with a = %d, b = %d/%d, starting from x_0 = %d. What is x_%d after %d steps? (enter a decimal)",
		a, bn, bd, x0, n, n)
	exp := fmt.Sprintf("Subgoal 1: State governing formula. Fixed-point iteration generates sequential estimates through the map \\(x_{k+1} = g(x_k)\\), where \\(g(x) = a + b x\\).\n\nSubgoal 2: Substitute parameters. Given affine map \\(g(x) = %d + \\frac{%d}{%d} x\\) and initial estimate \\(x_0 = %d\\), evaluate iterates iteratively.\n\nSubgoal 3: Evaluate and simplify. %s, concluding with \\(x_%d = %s\\).",
		a, bn, bd, x0, strings.Join(steps, ", "), n, ans)

	return generator.Problem{
		Question:    q,
		Answer:      ans,
		Explanation: exp,
		Facts: map[string]string{
			"a":      strconv.Itoa(a),
			"bn":     strconv.Itoa(bn),
			"bd":     strconv.Itoa(bd),
			"x0":     strconv.Itoa(x0),
			"n":      strconv.Itoa(n),
			"answer": ans,
		},
	}
}

// ---------------------------------------------------------------------------
// Newton's method (numeric)
// ---------------------------------------------------------------------------

type newtonGen struct{}

func (g *newtonGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	maxC := mathutil.ScaleInt(ctx.Difficulty, 8, 40)
	c := 2 + rand.Intn(maxC-1)
	maxX0 := mathutil.ScaleInt(ctx.Difficulty, 3, 10)
	x0 := 1 + rand.Intn(maxX0)

	for x0*x0 == c {
		c++
	}

	cStr := strconv.Itoa(c)
	x0Str := strconv.Itoa(x0)
	x1 := (float64(x0) + float64(c)/float64(x0)) / 2
	ans := decimal6(x1)

	for iter := 0; iter < 5 && (ans == cStr || ans == x0Str); iter++ {
		c += 3
		cStr = strconv.Itoa(c)
		x1 = (float64(x0) + float64(c)/float64(x0)) / 2
		ans = decimal6(x1)
	}

	q := fmt.Sprintf("One Newton step for f(x) = x² − c is x_1 = (x_0 + c/x_0)/2. For c = %s and x_0 = %s, what is x_1? (enter a decimal)", cStr, x0Str)
	exp := fmt.Sprintf("Subgoal 1: State governing formula. Newton's method for \\(f(x) = x^2 - c\\) simplifies to the Babylonian averaging recurrence \\(x_1 = \\frac{1}{2} \\left(x_0 + \\frac{c}{x_0}\\right)\\).\n\nSubgoal 2: Substitute parameters. Substitute parameter \\(c = %s\\) and initial guess \\(x_0 = %s\\): \\(x_1 = \\frac{1}{2} \\left(%s + \\frac{%s}{%s}\\right)\\).\n\nSubgoal 3: Evaluate and simplify. Compute quotient \\(\\frac{%s}{%s} = %s\\), add %s and divide by 2: \\(x_1 = %s\\).",
		cStr, x0Str, x0Str, cStr, x0Str,
		cStr, x0Str, decimal6(float64(c)/float64(x0)), x0Str, ans)

	return generator.Problem{
		Question:    q,
		Answer:      ans,
		Explanation: exp,
		Facts: map[string]string{
			"c":      cStr,
			"x0":     x0Str,
			"answer": ans,
		},
	}
}

// ---------------------------------------------------------------------------
// The secant method (numeric)
// ---------------------------------------------------------------------------

type secantGen struct{}

func (g *secantGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	maxC := mathutil.ScaleInt(ctx.Difficulty, 6, 30)
	c := float64(2 + rand.Intn(maxC-1))
	x0 := float64(1 + rand.Intn(4))
	x1 := x0 + float64(1+rand.Intn(3))

	f := func(x float64) float64 { return x*x - c }
	for f(x1) == f(x0) {
		x1++
	}
	x2 := x1 - f(x1)*(x1-x0)/(f(x1)-f(x0))
	ans := decimal6(x2)

	for iter := 0; iter < 5 && (ans == decimal6(c) || ans == decimal6(x0) || ans == decimal6(x1)); iter++ {
		c += 2
		x2 = x1 - f(x1)*(x1-x0)/(f(x1)-f(x0))
		ans = decimal6(x2)
	}

	q := fmt.Sprintf("One step of the secant method for f(x) = x² − c uses x_2 = x_1 − f(x_1)(x_1 − x_0)/(f(x_1) − f(x_0)). For c = %s, x_0 = %s and x_1 = %s, what is x_2? (enter a decimal)",
		decimal6(c), decimal6(x0), decimal6(x1))
	exp := fmt.Sprintf("Subgoal 1: State governing formula. The secant method replaces the analytical derivative with the difference quotient across the two previous points: \\(x_2 = x_1 - f(x_1) \\frac{x_1 - x_0}{f(x_1) - f(x_0)}\\).\n\nSubgoal 2: Substitute parameters. For \\(f(x) = x^2 - %s\\), evaluate \\(f(x_0) = f(%s) = %s\\) and \\(f(x_1) = f(%s) = %s\\). Substitute into the secant recurrence: \\(x_2 = %s - (%s) \\frac{%s - %s}{(%s) - (%s)}\\).\n\nSubgoal 3: Evaluate and simplify. Evaluate the finite difference quotient and subtract to obtain \\(x_2 = %s\\).",
		decimal6(c), decimal6(x0), decimal6(f(x0)), decimal6(x1), decimal6(f(x1)),
		decimal6(x1), decimal6(f(x1)), decimal6(x1), decimal6(x0), decimal6(f(x1)), decimal6(f(x0)),
		ans)

	return generator.Problem{
		Question:    q,
		Answer:      ans,
		Explanation: exp,
		Facts: map[string]string{
			"c":      decimal6(c),
			"x0":     decimal6(x0),
			"x1":     decimal6(x1),
			"answer": ans,
		},
	}
}

// ---------------------------------------------------------------------------
// The method of false position (multiple choice)
// ---------------------------------------------------------------------------

type falsePositionGen struct{}

func (g *falsePositionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return falsePositionProblem()
}

// ---------------------------------------------------------------------------
// Order of convergence (numeric)
// ---------------------------------------------------------------------------

type orderGen struct{}

func (g *orderGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	maxP := mathutil.ScaleInt(ctx.Difficulty, 3, 5)
	p := 1 + rand.Intn(maxP)
	a := p + 1 + rand.Intn(3)
	diff := 1 + rand.Intn(2)
	b := a + diff
	c := b + p*diff

	ans := decimal6(float64(c-b) / float64(b-a))

	q := fmt.Sprintf("The errors of an iterative method satisfy \\(e_{k+2} \\approx e_{k+1}^{p}/e_k^{p-1}\\). From three successive errors of \\(10^{-%d}\\), \\(10^{-%d}\\) and \\(10^{-%d}\\), estimate the order \\(p = \\frac{\\log(e_{k+2}/e_{k+1})}{\\log(e_{k+1}/e_k)}\\). (enter a decimal with a decimal point)",
		a, b, c)
	exp := fmt.Sprintf("Subgoal 1: State governing formula. The order of convergence p satisfies \\(p = \\frac{\\log(e_{k+2}/e_{k+1})}{\\log(e_{k+1}/e_k)}\\). For error powers \\(e_i = 10^{-a_i}\\), this logarithmic ratio reduces directly to the exponent difference ratio: \\(p = \\frac{c - b}{b - a}\\).\n\nSubgoal 2: Substitute parameters. Substitute exponents \\(a = %d\\), \\(b = %d\\) and \\(c = %d\\): \\(p = \\frac{%d - %d}{%d - %d}\\).\n\nSubgoal 3: Evaluate and simplify. Compute numerator \\(%d - %d = %d\\) and denominator \\(%d - %d = %d\\), then divide: \\(p = \\frac{%d}{%d} = %s\\).",
		a, b, c, c, b, b, a,
		c, b, c-b, b, a, b-a, c-b, b-a, ans)

	return generator.Problem{
		Question:    q,
		Answer:      ans,
		Explanation: exp,
		Facts: map[string]string{
			"a":      strconv.Itoa(a),
			"b":      strconv.Itoa(b),
			"c":      strconv.Itoa(c),
			"answer": ans,
		},
	}
}

// ---------------------------------------------------------------------------
// Multiple roots and modified Newton (multiple choice)
// ---------------------------------------------------------------------------

type multipleRootGen struct{}

func (g *multipleRootGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return multipleRootProblem()
}

// ---------------------------------------------------------------------------
// Bracketing versus open methods (multiple choice)
// ---------------------------------------------------------------------------

type bracketingGen struct{}

func (g *bracketingGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return bracketingProblem()
}

// ---------------------------------------------------------------------------
// Horner's method (numeric)
// ---------------------------------------------------------------------------

type hornerGen struct{}

func (g *hornerGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	deg := 2
	if ctx.Difficulty >= 0.5 {
		deg = 3
	}
	maxCoeff := mathutil.ScaleInt(ctx.Difficulty, 3, 10)
	coeffs := make([]int, deg+1)
	for i := range coeffs {
		coeffs[i] = rand.Intn(2*maxCoeff+1) - maxCoeff
	}
	if coeffs[deg] == 0 {
		coeffs[deg] = 1 + rand.Intn(maxCoeff)
	}
	maxX := mathutil.ScaleInt(ctx.Difficulty, 2, 5)
	x := 2 + rand.Intn(maxX-1)

	val := 0
	var stepDesc []string
	for k := len(coeffs) - 1; k >= 0; k-- {
		prev := val
		val = val*x + coeffs[k]
		if k == len(coeffs)-1 {
			stepDesc = append(stepDesc, fmt.Sprintf("b_%d = %d", k, val))
		} else {
			stepDesc = append(stepDesc, fmt.Sprintf("b_%d = (%d)(%d) + (%d) = %d", k, prev, x, coeffs[k], val))
		}
	}
	ans := strconv.Itoa(val)

	parts := make([]string, len(coeffs))
	for i, a := range coeffs {
		parts[i] = strconv.Itoa(a)
	}

	q := fmt.Sprintf("Horner's method evaluates a polynomial by nested multiplication. For the polynomial with coefficients (constant first) %s evaluated at x = %d, what is the value p(x)? (enter a decimal)",
		strings.Join(parts, ", "), x)
	exp := fmt.Sprintf("Subgoal 1: State governing formula. Horner's nested evaluation computes \\(p(x) = a_n x^n + \\dots + a_1 x + a_0\\) through intermediate values: \\(b_n = a_n\\) and \\(b_k = b_{k+1} x + a_k\\) for \\(k = n-1, \\dots, 0\\), where \\(p(x) = b_0\\).\n\nSubgoal 2: Substitute parameters. For coefficients \\((%s)\\) and evaluation point \\(x = %d\\), compute nested operations sequentially.\n\nSubgoal 3: Evaluate and simplify. %s, which yields final polynomial value \\(p(%d) = %s\\).",
		strings.Join(parts, ", "), x, strings.Join(stepDesc, ", "), x, ans)

	return generator.Problem{
		Question:    q,
		Answer:      ans,
		Explanation: exp,
		Facts: map[string]string{
			"coeffs": strings.Join(parts, ","),
			"x":      strconv.Itoa(x),
			"answer": ans,
		},
	}
}

// ---------------------------------------------------------------------------
// Lagrange interpolation (numeric)
// ---------------------------------------------------------------------------

type lagrangeGen struct{}

func (g *lagrangeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	maxA := mathutil.ScaleInt(ctx.Difficulty, 2, 5)
	a := 1 + rand.Intn(maxA)
	maxB := mathutil.ScaleInt(ctx.Difficulty, 2, 6)
	b := rand.Intn(2*maxB+1) - maxB
	maxC := mathutil.ScaleInt(ctx.Difficulty, 2, 8)
	c := rand.Intn(2*maxC+1) - maxC
	maxQ := mathutil.ScaleInt(ctx.Difficulty, 3, 6)
	qVal := 3 + rand.Intn(maxQ-2)

	y0 := c
	y1 := a + b + c
	y2 := 4*a + 2*b + c
	val := a*qVal*qVal + b*qVal + c
	ans := strconv.Itoa(val)

	for iter := 0; iter < 5 && (ans == strconv.Itoa(y0) || ans == strconv.Itoa(y1) || ans == strconv.Itoa(y2) || ans == strconv.Itoa(qVal)); iter++ {
		a++
		y1 = a + b + c
		y2 = 4*a + 2*b + c
		val = a*qVal*qVal + b*qVal + c
		ans = strconv.Itoa(val)
	}

	l0 := (qVal - 1) * (qVal - 2) / 2
	l1 := -qVal * (qVal - 2)
	l2 := qVal * (qVal - 1) / 2

	q := fmt.Sprintf("The Lagrange interpolant through the three points (0, %d), (1, %d) and (2, %d) is the quadratic through them. What is the interpolated value at x = %d? (enter a decimal)",
		y0, y1, y2, qVal)
	exp := fmt.Sprintf("Subgoal 1: State governing formula. The quadratic Lagrange interpolant through nodes 0, 1 and 2 is \\(P(x) = y_0 L_0(x) + y_1 L_1(x) + y_2 L_2(x)\\), with cardinal basis polynomials \\(L_0(x) = \\frac{(x-1)(x-2)}{2}\\), \\(L_1(x) = -x(x-2)\\) and \\(L_2(x) = \\frac{x(x-1)}{2}\\).\n\nSubgoal 2: Substitute parameters. At query point \\(x = %d\\), evaluate basis weights: \\(L_0(%d) = %d\\), \\(L_1(%d) = %d\\) and \\(L_2(%d) = %d\\).\n\nSubgoal 3: Evaluate and simplify. Multiply each basis weight by its nodal value: \\(P(%d) = (%d)(%d) + (%d)(%d) + (%d)(%d) = %s\\).",
		qVal, qVal, l0, qVal, l1, qVal, l2,
		qVal, y0, l0, y1, l1, y2, l2, ans)

	return generator.Problem{
		Question:    q,
		Answer:      ans,
		Explanation: exp,
		Facts: map[string]string{
			"y0":     strconv.Itoa(y0),
			"y1":     strconv.Itoa(y1),
			"y2":     strconv.Itoa(y2),
			"q":      strconv.Itoa(qVal),
			"answer": ans,
		},
	}
}

// ---------------------------------------------------------------------------
// Newton's divided differences (numeric)
// ---------------------------------------------------------------------------

type dividedDiffGen struct{}

func (g *dividedDiffGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	maxX0 := mathutil.ScaleInt(ctx.Difficulty, 1, 5)
	x0 := rand.Intn(maxX0 + 1)
	maxDx := mathutil.ScaleInt(ctx.Difficulty, 2, 6)
	dx := 1 + rand.Intn(maxDx)
	x1 := x0 + dx

	maxM := mathutil.ScaleInt(ctx.Difficulty, 4, 30)
	m := 2 + rand.Intn(maxM)
	for m == x0 || m == x1 {
		m++
	}
	maxY0 := mathutil.ScaleInt(ctx.Difficulty, 10, 30)
	y0 := 20 + rand.Intn(maxY0)
	for y0 == m {
		y0 += 5
	}
	y1 := y0 + m*dx
	for y1 == m {
		y0 += 5
		y1 = y0 + m*dx
	}
	ans := decimal6(float64(m))

	q := fmt.Sprintf("The first divided difference of the points (%d, %d) and (%d, %d) is (y_1 − y_0)/(x_1 − x_0). What is it? (enter a decimal)",
		x0, y0, x1, y1)
	exp := fmt.Sprintf("Subgoal 1: State governing formula. The first divided difference between two distinct points \\((x_0, y_0)\\) and \\((x_1, y_1)\\) measures the average rate of change: \\(f[x_0, x_1] = \\frac{y_1 - y_0}{x_1 - x_0}\\).\n\nSubgoal 2: Substitute parameters. Substitute coordinates \\((%d, %d)\\) and \\((%d, %d)\\): \\(f[x_0, x_1] = \\frac{%d - (%d)}{%d - %d}\\).\n\nSubgoal 3: Evaluate and simplify. Compute numerator \\(%d - (%d) = %d\\) and denominator \\(%d - %d = %d\\), then divide: \\(f[x_0, x_1] = \\frac{%d}{%d} = %s\\).",
		x0, y0, x1, y1, y1, y0, x1, x0,
		y1, y0, y1-y0, x1, x0, dx, y1-y0, dx, ans)

	return generator.Problem{
		Question:    q,
		Answer:      ans,
		Explanation: exp,
		Facts: map[string]string{
			"x0":     strconv.Itoa(x0),
			"y0":     strconv.Itoa(y0),
			"x1":     strconv.Itoa(x1),
			"y1":     strconv.Itoa(y1),
			"answer": ans,
		},
	}
}

// ---------------------------------------------------------------------------
// The interpolation error (multiple choice)
// ---------------------------------------------------------------------------

type interpErrorGen struct{}

func (g *interpErrorGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return interpErrorProblem()
}

// ---------------------------------------------------------------------------
// Chebyshev nodes (multiple choice)
// ---------------------------------------------------------------------------

type chebyshevGen struct{}

func (g *chebyshevGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return chebyshevProblem()
}

// ---------------------------------------------------------------------------
// Cubic splines (multiple choice)
// ---------------------------------------------------------------------------

type splinesGen struct{}

func (g *splinesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return splinesProblem()
}

func absInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
