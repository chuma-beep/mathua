package algebra

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

func waveKP(ctx generator.GeneratorContext) int {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = rand.Intn(3) + 1
	}
	return kp
}

func binomCoef(n, k int) int {
	if k < 0 || k > n {
		return 0
	}
	if k > n-k {
		k = n - k
	}
	out := 1
	for i := 1; i <= k; i++ {
		out = out * (n - k + i) / i
	}
	return out
}

// ---------------------------------------------------------------------------
// The binomial theorem
// ---------------------------------------------------------------------------

type binomialGen struct{}

func (g *binomialGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := waveKP(ctx)
	n := rand.Intn(4) + 3 // 3..6
	a := rand.Intn(4) + 1 // 1..4

	var question, explanation string
	var answer int
	k := 0
	coefficient := 0
	switch kp {
	case 1:
		k = rand.Intn(n-1) + 1
		coefficient = binomCoef(n, k) * powInt(a, n-k)
		answer = coefficient
		question = fmt.Sprintf("In the expansion of \\((x + %d)^{%d}\\), what is the coefficient of \\(x^{%d}\\)? (enter a number)", a, n, k)
		explanation = fmt.Sprintf("The \\(x^{%d}\\) term is \\(\\binom{%d}{%d}x^{%d}%d^{%d}\\), so the coefficient is %d × %d = %d.",
			k, n, k, k, a, n-k, binomCoef(n, k), powInt(a, n-k), coefficient)
	case 2:
		constant := powInt(a, n)
		answer = constant
		question = fmt.Sprintf("In the expansion of \\((x + %d)^{%d}\\), what is the constant term? (enter a number)", a, n)
		explanation = fmt.Sprintf("The constant term is the \\(x^0\\) term, which is \\(%d^{%d} = %d\\).", a, n, constant)
	default:
		total := powInt(a+1, n)
		answer = total
		question = fmt.Sprintf("What is the sum of all coefficients in the expansion of \\((x + %d)^{%d}\\)? (enter a number)", a, n)
		explanation = fmt.Sprintf("Substituting \\(x = 1\\) gives \\((%d + %d)^{%d} = %d\\), which is the sum of the coefficients.", a, 1, n, total)
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts: map[string]string{
			"n":           strconv.Itoa(n),
			"a":           strconv.Itoa(a),
			"k":           strconv.Itoa(k),
			"coefficient": strconv.Itoa(coefficient),
			"constant":    strconv.Itoa(powInt(a, n)),
			"sum":         strconv.Itoa(powInt(a+1, n)),
		},
		KPOrdinal: kp,
	}
}

// ---------------------------------------------------------------------------
// Infinite geometric series
// ---------------------------------------------------------------------------

type infiniteGeometricGen struct{}

func (g *infiniteGeometricGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := waveKP(ctx)
	d := rand.Intn(3) + 2 // ratio 1/d, d = 2..4
	a := (rand.Intn(5) + 1) * (d - 1)
	sum := a * d / (d - 1)

	var question, explanation string
	var answer int
	switch kp {
	case 1:
		answer = sum
		question = fmt.Sprintf("An infinite geometric series has first term %d and common ratio 1/%d. What is its sum? (enter a number)", a, d)
		explanation = fmt.Sprintf("The ratio satisfies \\(|r| < 1\\), so the sum is \\(a/(1-r) = %d/(1 - 1/%d) = %d\\).", a, d, sum)
	case 2:
		answer = d
		question = fmt.Sprintf("An infinite geometric series has first term %d and sum %d. Its common ratio is 1/d. What is d? (enter a number)", a, sum)
		explanation = fmt.Sprintf("Since \\(%d/(1 - 1/d) = %d\\), solving gives \\(d = %d\\).", a, sum, d)
	default:
		answer = a
		question = fmt.Sprintf("An infinite geometric series with common ratio 1/%d has sum %d. What is its first term? (enter a number)", d, sum)
		explanation = fmt.Sprintf("Use \\(a = S(1-r) = %d(1 - 1/%d) = %d\\).", sum, d, a)
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts: map[string]string{
			"first": strconv.Itoa(a), "denominator": strconv.Itoa(d), "sum": strconv.Itoa(sum),
		},
		KPOrdinal: kp,
	}
}

// ---------------------------------------------------------------------------
// Exponential function through two points
// ---------------------------------------------------------------------------

type expFromPointsGen struct{}

func (g *expFromPointsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := waveKP(ctx)
	a := rand.Intn(8) + 2 // 2..9
	b := rand.Intn(3) + 2 // 2..4

	var question, explanation string
	var answer int
	switch kp {
	case 1:
		answer = b
		question = fmt.Sprintf("An exponential function \\(y = a b^x\\) passes through \\((0, %d)\\) and \\((1, %d)\\). What is b? (enter a number)", a, a*b)
		explanation = fmt.Sprintf("At \\(x = 0\\), \\(y = a = %d\\). At \\(x = 1\\), \\(y = ab = %d\\), so \\(b = %d\\).", a, a*b, b)
	case 2:
		answer = a
		question = fmt.Sprintf("An exponential function \\(y = a b^x\\) passes through \\((1, %d)\\) and \\((2, %d)\\). What is a? (enter a number)", a*b, a*b*b)
		explanation = fmt.Sprintf("Divide the values: \\(b = %d/%d = %d\\). Then \\(a = %d/%d = %d\\).", a*b*b, a*b, b, a*b, b, a)
	default:
		answer = a * b * b * b
		question = fmt.Sprintf("An exponential function \\(y = a b^x\\) passes through \\((0, %d)\\) and \\((1, %d)\\). What is y when x = 3? (enter a number)", a, a*b)
		explanation = fmt.Sprintf("Here \\(a = %d\\) and \\(b = %d\\), so \\(y(3) = %d × %d^3 = %d\\).", a, b, a, b, a*b*b*b)
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts: map[string]string{
			"a": strconv.Itoa(a), "b": strconv.Itoa(b), "y3": strconv.Itoa(a * b * b * b),
		},
		KPOrdinal: kp,
	}
}

// ---------------------------------------------------------------------------
// Logarithmic scale
// ---------------------------------------------------------------------------

type logScaleGen struct{}

func (g *logScaleGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := waveKP(ctx)
	k := rand.Intn(6) + 1 // 1..6

	var question, explanation string
	var answer int
	switch kp {
	case 1:
		answer = k
		question = fmt.Sprintf("Two measurements are %d times apart. How many orders of magnitude separate them? (enter a number)", powInt(10, k))
		explanation = fmt.Sprintf("An order of magnitude is a factor of 10, and \\(10^{%d} = %d\\), so the difference is %d orders of magnitude.", k, powInt(10, k), k)
	case 2:
		answer = powInt(10, k)
		question = fmt.Sprintf("A quantity increases by %d decades. By what factor does it grow? (enter a number)", k)
		explanation = fmt.Sprintf("One decade is a factor of 10, so %d decades is a factor of \\(10^{%d} = %d\\).", k, k, powInt(10, k))
	default:
		answer = k
		question = fmt.Sprintf("What is \\(\\log_{10}(%d)\\)? (enter a number)", powInt(10, k))
		explanation = fmt.Sprintf("\\(\\log_{10}(%d) = %d\\) because \\(10^{%d} = %d\\).", powInt(10, k), k, k, powInt(10, k))
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts:       map[string]string{"exponent": strconv.Itoa(k), "value": strconv.Itoa(powInt(10, k))},
		KPOrdinal:   kp,
	}
}

// ---------------------------------------------------------------------------
// Least-squares regression line
// ---------------------------------------------------------------------------

type linearRegressionGen struct{}

func (g *linearRegressionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := waveKP(ctx)
	m := rand.Intn(4) + 1 // 1..4
	b := rand.Intn(7) - 3 // -3..3
	amp := rand.Intn(2) + 1
	xs := []int{-2, -1, 1, 2}
	noise := []int{1, -1, -1, 1}
	ys := make([]int, 4)
	for i, x := range xs {
		ys[i] = m*x + b + amp*noise[i]
	}

	var answer int
	var question, explanation string
	switch kp {
	case 1:
		answer = m
		question = fmt.Sprintf("A least-squares line is fitted to the points \\((-2, %d)\\), \\((-1, %d)\\), \\((1, %d)\\) and \\((2, %d)\\). What is its slope? (enter a number)",
			ys[0], ys[1], ys[2], ys[3])
		explanation = fmt.Sprintf("The x-values are symmetric and the residuals cancel in pairs, so the least-squares slope is the underlying slope %d. Direct calculation gives \\(m = %d\\).", m, m)
	case 2:
		answer = b
		question = fmt.Sprintf("A least-squares line is fitted to the points \\((-2, %d)\\), \\((-1, %d)\\), \\((1, %d)\\) and \\((2, %d)\\). What is its y-intercept? (enter a number)",
			ys[0], ys[1], ys[2], ys[3])
		explanation = fmt.Sprintf("With the slope %d, the line passes through the mean point \\((0, %d)\\), so the intercept is %d.", m, b, b)
	default:
		answer = 3*m + b
		question = fmt.Sprintf("A least-squares line is fitted to the points \\((-2, %d)\\), \\((-1, %d)\\), \\((1, %d)\\) and \\((2, %d)\\). What value does the line predict at x = 3? (enter a number)",
			ys[0], ys[1], ys[2], ys[3])
		explanation = fmt.Sprintf("The fitted line is \\(y = %dx + %d\\), so at \\(x = 3\\) it predicts %d.", m, b, answer)
	}

	facts := map[string]string{"slope": strconv.Itoa(m), "intercept": strconv.Itoa(b), "prediction": strconv.Itoa(3*m + b)}
	for i := range xs {
		facts[fmt.Sprintf("x%d", i+1)] = strconv.Itoa(xs[i])
		facts[fmt.Sprintf("y%d", i+1)] = strconv.Itoa(ys[i])
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts:       facts,
		KPOrdinal:   kp,
	}
}

// ---------------------------------------------------------------------------
// Parabola focus and directrix
// ---------------------------------------------------------------------------

type parabolaFocusGen struct{}

func (g *parabolaFocusGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := waveKP(ctx)
	p := rand.Intn(5) + 1

	var question, explanation string
	var answer int
	switch kp {
	case 1:
		answer = p
		question = fmt.Sprintf("The parabola \\(y = x^2/%d\\) has its vertex at the origin. What is the y-coordinate of its focus? (enter a number)", 4*p)
		explanation = fmt.Sprintf("For \\(y = x^2/(4p)\\), the focus is \\((0, p)\\). Here \\(4p = %d\\), so \\(p = %d\\).", 4*p, p)
	case 2:
		answer = -p
		question = fmt.Sprintf("The parabola \\(y = x^2/%d\\) opens upward. What is the y-coordinate of its directrix? (enter a number)", 4*p)
		explanation = fmt.Sprintf("The focus is \\((0, %d)\\), and the directrix is the same distance below the vertex, so \\(y = %d\\).", p, -p)
	default:
		answer = 4 * p
		question = fmt.Sprintf("The parabola \\(y = x^2/%d\\) has focus \\((0, %d)\\). What is the length of its latus rectum? (enter a number)", 4*p, p)
		explanation = fmt.Sprintf("The latus rectum has length \\(4p\\). With \\(p = %d\\), the length is %d.", p, 4*p)
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts:       map[string]string{"p": strconv.Itoa(p), "focus": strconv.Itoa(p), "directrix": strconv.Itoa(-p), "latus": strconv.Itoa(4 * p)},
		KPOrdinal:   kp,
	}
}

// ---------------------------------------------------------------------------
// Ellipse foci
// ---------------------------------------------------------------------------

type ellipseFociGen struct{}

func (g *ellipseFociGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := waveKP(ctx)
	type triple struct{ a, b, c int }
	pairs := []triple{{5, 3, 4}, {5, 4, 3}, {10, 6, 8}, {10, 8, 6}, {13, 5, 12}, {13, 12, 5}, {17, 8, 15}, {17, 15, 8}}
	t := pairs[rand.Intn(len(pairs))]

	var question, explanation string
	var answer int
	switch kp {
	case 1:
		answer = t.c
		question = fmt.Sprintf("The ellipse \\(x^2/%d + y^2/%d = 1\\) has foci on the x-axis. What is c in the focus \\((c, 0)\\)? (enter a number)", t.a*t.a, t.b*t.b)
		explanation = fmt.Sprintf("For an ellipse, \\(c^2 = a^2 - b^2 = %d - %d = %d\\), so \\(c = %d\\).", t.a*t.a, t.b*t.b, t.c*t.c, t.c)
	case 2:
		answer = 2 * t.a
		question = fmt.Sprintf("The ellipse \\(x^2/%d + y^2/%d = 1\\) has foci \\((±c, 0)\\). What is the sum of the distances from any point on the ellipse to the two foci? (enter a number)", t.a*t.a, t.b*t.b)
		explanation = fmt.Sprintf("The defining property of an ellipse gives a constant sum of \\(2a\\). Here \\(a = %d\\), so the sum is %d.", t.a, 2*t.a)
	default:
		answer = t.c
		question = fmt.Sprintf("The ellipse \\(x^2/%d + y^2/%d = 1\\) has a focus at \\((c, 0)\\) with \\(c > 0\\). What is that x-coordinate? (enter a number)", t.a*t.a, t.b*t.b)
		explanation = fmt.Sprintf("\\(c^2 = a^2 - b^2 = %d - %d = %d\\), and \\(c > 0\\), so the right focus is at \\(x = %d\\).", t.a*t.a, t.b*t.b, t.c*t.c, t.c)
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts:       map[string]string{"a": strconv.Itoa(t.a), "b": strconv.Itoa(t.b), "c": strconv.Itoa(t.c), "sum": strconv.Itoa(2 * t.a)},
		KPOrdinal:   kp,
	}
}

// ---------------------------------------------------------------------------
// Absolute-value inequalities on a number line
// ---------------------------------------------------------------------------

type absoluteInequalityGraphGen struct{}

func (g *absoluteInequalityGraphGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := waveKP(ctx)
	h := rand.Intn(11) - 5 // -5..5
	r := rand.Intn(5) + 1  // 1..5

	centerText := "x"
	switch {
	case h > 0:
		centerText = fmt.Sprintf("x - %d", h)
	case h < 0:
		centerText = fmt.Sprintf("x + %d", -h)
	}

	var question, explanation string
	var answer int
	switch kp {
	case 1:
		answer = h
		question = fmt.Sprintf("The inequality \\(|%s| ≤ %d\\) has a solution interval. What is its center? (enter a number)", centerText, r)
		explanation = fmt.Sprintf("The inequality means the distance from \\(x\\) to %d is at most %d, so the interval is \\([%d, %d]\\) and its center is %d.", h, r, h-r, h+r, h)
	case 2:
		answer = r
		question = fmt.Sprintf("The inequality \\(|%s| ≤ %d\\) has a solution interval. What is its radius? (enter a number)", centerText, r)
		explanation = fmt.Sprintf("The bound %d is the distance from the center %d to either endpoint, so the radius is %d.", r, h, r)
	default:
		answer = 2*r + 1
		question = fmt.Sprintf("How many integers satisfy the inequality \\(|%s| ≤ %d\\)? (enter a number)", centerText, r)
		explanation = fmt.Sprintf("The interval is \\([%d, %d]\\), which contains \\(2(%d) + 1 = %d\\) integers.", h-r, h+r, r, 2*r+1)
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts:       map[string]string{"center": strconv.Itoa(h), "radius": strconv.Itoa(r), "count": strconv.Itoa(2*r + 1)},
		KPOrdinal:   kp,
	}
}

// ---------------------------------------------------------------------------
// Zeros and multiplicity
// ---------------------------------------------------------------------------

type zerosMultiplicityGen struct{}

func (g *zerosMultiplicityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := waveKP(ctx)
	nRoots := rand.Intn(2) + 2 // 2..3
	used := map[int]bool{}
	roots := make([]int, 0, nRoots)
	for len(roots) < nRoots {
		r := rand.Intn(7) - 3 // -3..3
		if used[r] {
			continue
		}
		used[r] = true
		roots = append(roots, r)
	}
	mults := make([]int, nRoots)
	degree := 0
	for i := range roots {
		mults[i] = rand.Intn(3) + 1 // 1..3
		degree += mults[i]
	}

	parts := make([]string, nRoots)
	for i, r := range roots {
		factor := "x"
		switch {
		case r > 0:
			factor = fmt.Sprintf("(x - %d)", r)
		case r < 0:
			factor = fmt.Sprintf("(x + %d)", -r)
		}
		if mults[i] > 1 {
			factor = fmt.Sprintf("%s^%d", factor, mults[i])
		}
		parts[i] = factor
	}
	var polyBuilder strings.Builder
	for i, part := range parts {
		if i > 0 {
			polyBuilder.WriteString(" ")
		}
		polyBuilder.WriteString(part)
	}
	poly := polyBuilder.String()

	var question, explanation string
	var answer int
	askedRoot, askedMult := 0, 0
	switch kp {
	case 1:
		i := rand.Intn(nRoots)
		answer = mults[i]
		askedRoot, askedMult = roots[i], mults[i]
		question = fmt.Sprintf("The polynomial \\(f(x) = %s\\) has a zero at \\(x = %d\\). What is its multiplicity? (enter a number)", poly, roots[i])
		explanation = fmt.Sprintf("The factor \\((x - %d)\\) appears %d %s, so the zero \\(x = %d\\) has multiplicity %d.", roots[i], mults[i], pluralWord(mults[i], "time", "times"), roots[i], mults[i])
	case 2:
		answer = degree
		question = fmt.Sprintf("What is the degree of \\(f(x) = %s\\)? (enter a number)", poly)
		explanation = fmt.Sprintf("Add the exponents: %s = %d, so the degree is %d.", joinInts(mults, " + "), degree, degree)
	default:
		answer = nRoots
		question = fmt.Sprintf("How many distinct zeros does \\(f(x) = %s\\) have? (enter a number)", poly)
		explanation = fmt.Sprintf("The distinct zeros are %s, so there are %d.", joinInts(roots, ", "), nRoots)
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts: map[string]string{
			"degree": strconv.Itoa(degree), "distinct": strconv.Itoa(nRoots),
			"multiplicities": joinInts(mults, ","), "roots": joinInts(roots, ","),
			"asked_root": strconv.Itoa(askedRoot), "asked_multiplicity": strconv.Itoa(askedMult),
		},
		KPOrdinal: kp,
	}
}

func joinInts(nums []int, sep string) string {
	var out strings.Builder
	for i, n := range nums {
		if i > 0 {
			out.WriteString(sep)
		}
		out.WriteString(strconv.Itoa(n))
	}
	return out.String()
}
