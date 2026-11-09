package realanalysis

// Wave 5 real analysis, remaining nineteen concepts. The generators use finite
// measure spaces, finite grids, polynomials and small sequences so every answer
// is an integer and every assumption can be checked in the Facts.

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave5Remaining(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"ra.seq.double_series":              "double_series",
		"ra.seq.fourier_convergence":        "fourier_convergence",
		"ra.measure.measurable_sets":        "measurable_sets",
		"ra.measure.measurable_functions":   "measurable_functions",
		"ra.measure.lebesgue_integral":      "lebesgue_integral",
		"ra.measure.monotone_convergence":   "monotone_convergence",
		"ra.measure.fatou":                  "fatou",
		"ra.measure.dominated_convergence":  "dominated_convergence",
		"ra.measure.product":                "product",
		"ra.measure.fubini_tonelli":         "fubini_tonelli",
		"ra.metric.lp_space":                "lp_space",
		"ra.metric.holder":                  "holder",
		"ra.metric.minkowski":               "minkowski",
		"ra.metric.banach":                  "banach",
		"ra.metric.hilbert":                 "hilbert",
		"ra.metric.bounded_operator":        "bounded_operator",
		"ra.metric.riesz_representation":    "riesz_representation",
		"ra.deriv.absolute_continuity":      "absolute_continuity",
		"ra.deriv.lebesgue_differentiation": "lebesgue_differentiation",
	} {
		reg.Register(id, &w5RARemainingGen{kind: kind})
	}
}

type w5RARemainingGen struct{ kind string }

func (g *w5RARemainingGen) Generate(ctx generator.GeneratorContext) (p generator.Problem) {
	defer func() {
		p.Question = w5Clean(p.Question)
		p.Explanation = w5Clean(p.Explanation)
	}()
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "double_series":
		return w5DoubleSeries(kp)
	case "fourier_convergence":
		return w5FourierConvergence(kp)
	case "measurable_sets":
		return w5MeasurableSets(kp)
	case "measurable_functions":
		return w5MeasurableFunctions(kp)
	case "lebesgue_integral":
		return w5LebesgueIntegral(kp)
	case "monotone_convergence":
		return w5MonotoneConvergence(kp)
	case "fatou":
		return w5Fatou(kp)
	case "dominated_convergence":
		return w5DominatedConvergence(kp)
	case "product":
		return w5ProductMeasure(kp)
	case "fubini_tonelli":
		return w5Fubini(kp)
	case "lp_space":
		return w5LpSpace(kp)
	case "holder":
		return w5Holder(kp)
	case "minkowski":
		return w5Minkowski(kp)
	case "banach":
		return w5Banach(kp)
	case "hilbert":
		return w5Hilbert(kp)
	case "bounded_operator":
		return w5BoundedOperator(kp)
	case "riesz_representation":
		return w5Riesz(kp)
	case "absolute_continuity":
		return w5AbsoluteContinuity(kp)
	case "lebesgue_differentiation":
		return w5LebesgueDifferentiation(kp)
	}
	return generator.Problem{Question: "What is 2 + 2? (enter a number)", Answer: "4", Explanation: "2 + 2 = 4."}
}

var w5PureNumber = regexp.MustCompile(`\\\((-?\d+)\\\)`)
var w5UnitCoefficient = regexp.MustCompile(`(^|[^0-9])1 ([a-zA-Z])`)
var w5UnitParen = regexp.MustCompile(`(^|[^0-9])1\(`)

// w5Clean removes presentation artefacts that are not part of the mathematics:
// pure numbers wrapped as math, unit coefficients and doubled signs.
func w5Clean(s string) string {
	s = w5PureNumber.ReplaceAllString(s, "$1")
	s = strings.ReplaceAll(s, "+-", "-")
	s = strings.ReplaceAll(s, "--", "+")
	s = w5UnitCoefficient.ReplaceAllString(s, "${1}${2}")
	s = w5UnitParen.ReplaceAllString(s, "${1}(")
	s = strings.ReplaceAll(s, "  ", " ")
	return s
}

func w5Abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func w5DoubleSeries(kp int) generator.Problem {
	m, n := w5RR(2, 5), w5RR(2, 5)
	f := map[string]string{}
	w5RI(f, "m", m)
	w5RI(f, "n", n)
	switch kp {
	case 1:
		s := n*m*(m+1)/2 + m*n*(n+1)/2
		return w5R(kp, fmt.Sprintf("For \\(a_{ij}=i+j\\) on \\(1\\le i\\le %d\\), \\(1\\le j\\le %d\\), what is the finite double sum? (enter a number)", m, n), strconv.Itoa(s),
			fmt.Sprintf("Sum the row sums: the total is %d.", s), f)
	case 2:
		s := (m * (m + 1) / 2) * (n * (n + 1) / 2)
		return w5R(kp, fmt.Sprintf("For \\(a_{ij}=ij\\) on \\(1\\le i\\le %d\\), \\(1\\le j\\le %d\\), what is the finite double sum? (enter a number)", m, n), strconv.Itoa(s),
			fmt.Sprintf("The sum factors as \\((1+\\cdots+%d)(1+\\cdots+%d)=%d\\).", m, n, s), f)
	default:
		s := 0
		for i := 1; i <= m; i++ {
			for j := 1; j <= n; j++ {
				s += i + j
			}
		}
		return w5R(kp, fmt.Sprintf("For \\(a_{ij}=i+j\\) on \\(1\\le i\\le %d\\), \\(1\\le j\\le %d\\), what is the sum of absolute values? (enter a number)", m, n), strconv.Itoa(s),
			fmt.Sprintf("All terms are positive here, so the absolute sum is the same as the ordinary sum, %d.", s), f)
	}
}

func w5FourierConvergence(kp int) generator.Problem {
	c, n, p := w5RR(1, 5), w5RR(1, 10), w5RR(1, 3)
	f := map[string]string{}
	w5RI(f, "coefficient", c)
	w5RI(f, "index", n)
	w5RI(f, "power", p)
	switch kp {
	case 1:
		return w5R(kp, fmt.Sprintf("A Fourier coefficient has the form \\(a_n=a/n^p\\). For \\(a=%d\\), \\(n=%d\\) and \\(p=%d\\), what is \\(a_n n^p\\)? (enter a number)", c, n, p), strconv.Itoa(c),
			fmt.Sprintf("Multiplying by \\(n^p\\) cancels the denominator and leaves the numerator %d.", c), f)
	case 2:
		return w5R(kp, fmt.Sprintf("For the Dirichlet kernel \\(D_n(0)\\), what is its value at zero when \\(n=%d\\)? (enter a number)", n), strconv.Itoa(2*n+1),
			fmt.Sprintf("\\(D_n(0)=2n+1=2(%d)+1=%d\\).", n, 2*n+1), f)
	default:
		a := w5RR(-3, 3)
		b := w5RR(-3, 3)
		w5RI(f, "a0", a)
		w5RI(f, "a1", b)
		return w5R(kp, fmt.Sprintf("A finite Fourier polynomial has \\(f(0)=a_0+a_1\\) with \\(a_0=%d\\) and \\(a_1=%d\\). What is \\(f(0)\\)? (enter a number)", a, b), strconv.Itoa(a+b),
			fmt.Sprintf("Substitute the coefficients: \\(%d+%d=%d\\).", a, b, a+b), f)
	}
}

func w5MeasurableSets(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		k := w5RR(1, 10)
		w5RI(f, "atoms", k)
		return w5R(kp, fmt.Sprintf("A finite sigma-algebra is generated by a partition into %d atoms. How many measurable sets does it contain? (enter a number)", k), strconv.Itoa(1<<k),
			fmt.Sprintf("Each atom is either included or excluded, so the count is \\(2^{%d}=%d\\).", k, 1<<k), f)
	case 2:
		n := w5RR(3, 9)
		r := w5RR(0, 2)
		count := 0
		for i := 1; i <= n; i++ {
			if i%3 == r {
				count++
			}
		}
		w5RI(f, "n", n)
		w5RI(f, "residue", r)
		verb := fmt.Sprintf("there are %d", count)
		if count == 1 {
			verb = "there is 1"
		}
		return w5R(kp, fmt.Sprintf("For \\(f(i)=i\\bmod 3\\) on \\(\\{1,\\ldots,%d\\}\\), how many points lie in the preimage of the residue \\(%d\\)? (enter a number)", n, r), strconv.Itoa(count),
			fmt.Sprintf("Count the integers in the range with that residue: %s.", verb), f)
	case 3:
		a, b, c := w5RR(1, 4), w5RR(1, 4), w5RR(1, 4)
		w5RI(f, "m1", a)
		w5RI(f, "m2", b)
		w5RI(f, "m3", c)
		return w5R(kp, fmt.Sprintf("Three disjoint measurable sets have measures %d, %d and %d. What is the measure of their union? (enter a number)", a, b, c), strconv.Itoa(a+b+c),
			fmt.Sprintf("Additivity gives \\(%d+%d+%d=%d\\).", a, b, c, a+b+c), f)
	}
	return generator.Problem{}
}

func w5MeasurableFunctions(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		x := w5RR(1, 8)
		thr := w5RR(1, 8)
		val := 0
		if x >= thr {
			val = 1
		}
		w5RI(f, "x", x)
		w5RI(f, "threshold", thr)
		return w5R(kp, fmt.Sprintf("For the indicator \\(1_{[%d,\\infty)}(x)\\), what is its value at \\(x=%d\\)? (enter a number)", thr, x), strconv.Itoa(val),
			fmt.Sprintf("The point is %s the threshold, so the indicator is %d.", map[int]string{0: "below", 1: "at or above"}[val], val), f)
	case 2:
		x := w5RR(1, 6)
		a := w5RR(1, 4)
		w5RI(f, "x", x)
		w5RI(f, "scale", a)
		return w5R(kp, fmt.Sprintf("For the measurable function \\(f(x)=%d x\\) and the continuous function \\(g(y)=y+1\\), what is \\(g(f(%d))\\)? (enter a number)", a, x), strconv.Itoa(a*x+1),
			fmt.Sprintf("\\(f(%d)=%d\\), and adding 1 gives %d.", x, a*x, a*x+1), f)
	default:
		x := w5RR(1, 6)
		a, b := w5RR(-4, 4), w5RR(-4, 4)
		ans := a * x
		if b*x > ans {
			ans = b * x
		}
		w5RI(f, "x", x)
		w5RI(f, "a", a)
		w5RI(f, "b", b)
		return w5R(kp, fmt.Sprintf("For measurable functions \\(f(x)=%d x\\) and \\(g(x)=%d x\\), what is \\(\\max(f,g)\\) at \\(x=%d\\)? (enter a number)", a, b, x), strconv.Itoa(ans),
			fmt.Sprintf("Compare %d and %d; the larger is %d.", a*x, b*x, ans), f)
	}
}

func w5LebesgueIntegral(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		c1, c2 := w5RR(1, 5), w5RR(1, 5)
		l1, l2 := w5RR(1, 4), w5RR(1, 4)
		w5RI(f, "c1", c1)
		w5RI(f, "c2", c2)
		w5RI(f, "l1", l1)
		w5RI(f, "l2", l2)
		return w5R(kp, fmt.Sprintf("A step function equals %d on a set of measure %d and %d on a disjoint set of measure %d. What is its integral? (enter a number)", c1, l1, c2, l2), strconv.Itoa(c1*l1+c2*l2),
			fmt.Sprintf("The integral is \\(%d\\cdot%d+%d\\cdot%d=%d\\).", c1, l1, c2, l2, c1*l1+c2*l2), f)
	case 2:
		c := w5RR(1, 6)
		l := w5RR(1, 6)
		w5RI(f, "c", c)
		w5RI(f, "length", l)
		return w5R(kp, fmt.Sprintf("For the constant function %d on an interval of length %d, what is the Lebesgue integral? (enter a number)", c, l), strconv.Itoa(c*l),
			fmt.Sprintf("The Riemann and Lebesgue integrals agree here: \\(%d\\cdot%d=%d\\).", c, l, c*l), f)
	default:
		c1, c2, c3 := w5RR(1, 5), w5RR(1, 5), w5RR(1, 5)
		w1, w2, w3 := w5RR(1, 3), w5RR(1, 3), w5RR(1, 3)
		w5RI(f, "c1", c1)
		w5RI(f, "c2", c2)
		w5RI(f, "c3", c3)
		w5RI(f, "w1", w1)
		w5RI(f, "w2", w2)
		w5RI(f, "w3", w3)
		return w5R(kp, fmt.Sprintf("Under counting measure with weights %d, %d and %d, a function has values %d, %d and %d. What is its integral? (enter a number)", w1, w2, w3, c1, c2, c3), strconv.Itoa(c1*w1+c2*w2+c3*w3),
			fmt.Sprintf("The integral is \\(%d\\cdot%d+%d\\cdot%d+%d\\cdot%d=%d\\).", c1, w1, c2, w2, c3, w3, c1*w1+c2*w2+c3*w3), f)
	}
}

func w5MonotoneConvergence(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		c := w5RR(1, 12)
		w5RI(f, "limit", c)
		return w5R(kp, fmt.Sprintf("An increasing sequence of functions converges pointwise to the constant %d on a set of measure 1. What is the limit of the integrals? (enter a number)", c), strconv.Itoa(c),
			fmt.Sprintf("Monotone convergence gives the integral of the limit, \\(%d\\cdot1=%d\\).", c, c), f)
	case 2:
		n := w5RR(2, 11)
		w5RI(f, "n", n)
		return w5R(kp, fmt.Sprintf("The indicators \\(1_{[0,1-1/%d]}\\) increase to \\(1_{[0,1]}\\). What is the limit of their integrals? (enter a number)", n), "1",
			"The limiting indicator is 1 on [0,1], whose integral is 1. The answer is 1.", f)
	default:
		c, m := w5RR(1, 5), w5RR(1, 5)
		w5RI(f, "value", c)
		w5RI(f, "measure", m)
		return w5R(kp, fmt.Sprintf("Functions increasing to the constant %d on a set of measure %d have integrals tending to what value? (enter a number)", c, m), strconv.Itoa(c*m),
			fmt.Sprintf("The limit integral is \\(%d\\cdot%d=%d\\).", c, m, c*m), f)
	}
}

func w5Fatou(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		b := w5RR(1, 10)
		w5RI(f, "pointwise", b)
		return w5R(kp, fmt.Sprintf("Fatou's lemma gives the lower bound \\(\\liminf\\int f_n\\ge\\int\\liminf f_n\\). If the right-hand side is %d, what is the lower bound on the left-hand side? (enter a number)", b), strconv.Itoa(b),
			fmt.Sprintf("The lower bound is the integral of the pointwise liminf, %d.", b), f)
	case 2:
		n := w5RR(2, 11)
		w5RI(f, "n", n)
		return w5R(kp, fmt.Sprintf("For \\(f_n=%d\\) on \\([0,1/%d]\\) and zero elsewhere, what is \\(\\int f_n\\)? (enter a number)", n, n), "1",
			fmt.Sprintf("The rectangle has height %d and width 1/%d, so the integral is 1.", n, n), f)
	default:
		c := w5RR(1, 6)
		w5RI(f, "c", c)
		return w5R(kp, fmt.Sprintf("For the nonnegative constant function %d on \\([0,1]\\), what is its integral? (enter a number)", c), strconv.Itoa(c),
			fmt.Sprintf("The integral is \\(%d\\cdot1=%d\\).", c, c), f)
	}
}

func w5DominatedConvergence(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		c := w5RR(1, 12)
		w5RI(f, "limit", c)
		return w5R(kp, fmt.Sprintf("Functions \\(f_n\\) converge pointwise to the constant %d on \\([0,1]\\) and are dominated by an integrable function. What is the limit of their integrals? (enter a number)", c), strconv.Itoa(c),
			fmt.Sprintf("Dominated convergence gives the integral of the limit, \\(%d\\cdot1=%d\\).", c, c), f)
	case 2:
		g := w5RR(1, 8)
		w5RI(f, "bound", g)
		return w5R(kp, fmt.Sprintf("If \\(|f_n|\\le g\\) and \\(\\int g=%d\\), what is the largest possible value of \\(\\int|f_n|\\)? (enter a number)", g), strconv.Itoa(g),
			fmt.Sprintf("The bound gives \\(\\int|f_n|\\le\\int g=%d\\).", g), f)
	default:
		n := w5RR(2, 11)
		w5RI(f, "n", n)
		return w5R(kp, fmt.Sprintf("The functions \\(1_{[0,1-1/%d]}\\) are dominated by 1 and converge to 1. What is the limit of their integrals? (enter a number)", n), "1",
			"The limit function is 1 on [0,1], whose integral is 1. The answer is 1.", f)
	}
}

func w5ProductMeasure(kp int) generator.Problem {
	a, b := w5RR(1, 5), w5RR(1, 5)
	f := map[string]string{}
	w5RI(f, "a", a)
	w5RI(f, "b", b)
	switch kp {
	case 1:
		return w5R(kp, fmt.Sprintf("A rectangle has side lengths %d and %d. What is its product measure? (enter a number)", a, b), strconv.Itoa(a*b),
			fmt.Sprintf("The product measure is \\(%d\\cdot%d=%d\\).", a, b, a*b), f)
	case 2:
		return w5R(kp, fmt.Sprintf("A finite product sigma-algebra is generated by %d sets in the first factor and %d in the second. How many measurable rectangles does it have? (enter a number)", a, b), strconv.Itoa(a*b),
			fmt.Sprintf("Each first-factor set pairs with each second-factor set, giving \\(%d\\cdot%d=%d\\).", a, b, a*b), f)
	default:
		m, n := w5RR(1, 4), w5RR(1, 4)
		w5RI(f, "m", m)
		w5RI(f, "n", n)
		return w5R(kp, fmt.Sprintf("A product set has measure %d in the first factor and %d in the second. What is its product measure? (enter a number)", m, n), strconv.Itoa(m*n),
			fmt.Sprintf("The product measure is \\(%d\\cdot%d=%d\\).", m, n, m*n), f)
	}
}

func w5Fubini(kp int) generator.Problem {
	a, b, c := w5RR(1, 4), w5RR(1, 4), w5RR(1, 4)
	f := map[string]string{}
	w5RI(f, "a", a)
	w5RI(f, "b", b)
	w5RI(f, "c", c)
	switch kp {
	case 1:
		return w5R(kp, fmt.Sprintf("For the constant function %d on a rectangle of area %d, what is the iterated integral? (enter a number)", c, a*b), strconv.Itoa(c*a*b),
			fmt.Sprintf("The integral is \\(%d\\cdot%d=%d\\).", c, a*b, c*a*b), f)
	case 2:
		return w5R(kp, fmt.Sprintf("For the constant function %d on a rectangle of area %d, what is the integral after swapping the order? (enter a number)", c, a*b), strconv.Itoa(c*a*b),
			fmt.Sprintf("Fubini gives the same value, \\(%d\\cdot%d=%d\\).", c, a*b, c*a*b), f)
	default:
		s := 0
		for i := 1; i <= a; i++ {
			for j := 1; j <= b; j++ {
				s += i * c
			}
		}
		return w5R(kp, fmt.Sprintf("For \\(f(i,j)=%d i\\) on the grid \\(1\\le i\\le %d\\), \\(1\\le j\\le %d\\), what is the finite double sum? (enter a number)", c, a, b), strconv.Itoa(s),
			fmt.Sprintf("Sum the rows: the total is %d.", s), f)
	}
}

func w5LpSpace(kp int) generator.Problem {
	a, b := w5RNZ(-5, 5), w5RNZ(-5, 5)
	f := map[string]string{}
	w5RI(f, "a", a)
	w5RI(f, "b", b)
	switch kp {
	case 1:
		return w5R(kp, fmt.Sprintf("For the vector \\((%d,%d)\\), what is its \\(\\ell^1\\) norm? (enter a number)", a, b), strconv.Itoa(w5Abs(a)+w5Abs(b)),
			fmt.Sprintf("The \\(\\ell^1\\) norm is \\(|%d|+|%d|=%d\\).", a, b, w5Abs(a)+w5Abs(b)), f)
	case 2:
		return w5R(kp, fmt.Sprintf("For the vector \\((%d,%d)\\), what is the square of its \\(\\ell^2\\) norm? (enter a number)", a, b), strconv.Itoa(a*a+b*b),
			fmt.Sprintf("The square is \\((%d)^2+(%d)^2=%d\\).", a, b, a*a+b*b), f)
	default:
		return w5R(kp, fmt.Sprintf("For the vector \\((%d,%d)\\), what is its \\(\\ell^\\infty\\) norm? (enter a number)", a, b), strconv.Itoa(w5Max(w5Abs(a), w5Abs(b))),
			fmt.Sprintf("The maximum absolute component is %d.", w5Max(w5Abs(a), w5Abs(b))), f)
	}
}

func w5Holder(kp int) generator.Problem {
	a, b := w5RNZ(-4, 4), w5RNZ(-4, 4)
	c, d := w5RNZ(-4, 4), w5RNZ(-4, 4)
	f := map[string]string{}
	w5RI(f, "a", a)
	w5RI(f, "b", b)
	w5RI(f, "c", c)
	w5RI(f, "d", d)
	switch kp {
	case 1:
		return w5R(kp, fmt.Sprintf("For \\(x=(%d,%d)\\) and \\(y=(%d,%d)\\), what is the dot product? (enter a number)", a, b, c, d), strconv.Itoa(a*c+b*d),
			fmt.Sprintf("The dot product is \\((%d)(%d)+(%d)(%d)=%d\\).", a, c, b, d, a*c+b*d), f)
	case 2:
		return w5R(kp, fmt.Sprintf("For \\(x=(%d,%d)\\) and \\(y=%d x\\), what value of \\(c\\) makes \\(y=c x\\)? (enter a number)", a, b, 2), "2",
			"The second vector is twice the first, so the equality constant is 2.", f)
	default:
		left := (a*a + b*b) * (c*c + d*d)
		return w5R(kp, fmt.Sprintf("For \\(x=(%d,%d)\\) and \\(y=(%d,%d)\\), what is \\(\\|x\\|^2\\|y\\|^2\\)? (enter a number)", a, b, c, d), strconv.Itoa(left),
			fmt.Sprintf("Multiply the squared norms: \\(((%d)^2+(%d)^2)((%d)^2+(%d)^2)=%d\\).", a, b, c, d, left), f)
	}
}

func w5Minkowski(kp int) generator.Problem {
	a, b := w5RNZ(-4, 4), w5RNZ(-4, 4)
	c, d := w5RNZ(-4, 4), w5RNZ(-4, 4)
	f := map[string]string{}
	w5RI(f, "a", a)
	w5RI(f, "b", b)
	w5RI(f, "c", c)
	w5RI(f, "d", d)
	switch kp {
	case 1:
		return w5R(kp, fmt.Sprintf("For \\(x=(%d,%d)\\) and \\(y=(%d,%d)\\), what is the square of \\(\\|x+y\\|_2\\)? (enter a number)", a, b, c, d), strconv.Itoa((a+c)*(a+c)+(b+d)*(b+d)),
			fmt.Sprintf("Add the vectors and square the components: the result is %d.", (a+c)*(a+c)+(b+d)*(b+d)), f)
	case 2:
		return w5R(kp, fmt.Sprintf("For \\(x=(%d,%d)\\) and \\(y=(%d,%d)\\), what is \\(\\|x+y\\|_1\\)? (enter a number)", a, b, c, d), strconv.Itoa(w5Abs(a+c)+w5Abs(b+d)),
			fmt.Sprintf("The \\(\\ell^1\\) norm of the sum is %d.", w5Abs(a+c)+w5Abs(b+d)), f)
	default:
		return w5R(kp, fmt.Sprintf("For \\(x=(%d,%d)\\) and \\(y=%d x\\), what value of \\(c\\) makes \\(y=c x\\)? (enter a number)", a, b, 3), "3",
			"The second vector is three times the first, so the equality constant is 3.", f)
	}
}

func w5Banach(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		n := w5RR(1, 10)
		w5RI(f, "n", n)
		return w5R(kp, fmt.Sprintf("Is \\(\\mathbb{R}^{%d}\\) complete with the Euclidean metric? Enter 1 for yes and 0 for no.", n), "1",
			"Every finite-dimensional Euclidean space is complete. The answer is 1.", f)
	case 2:
		l := w5RNZ(-5, 5)
		w5RI(f, "limit", l)
		return w5R(kp, fmt.Sprintf("The Cauchy sequence \\(a_n=%d+1/n\\) has what limit? (enter a number)", l), strconv.Itoa(l),
			fmt.Sprintf("The sequence converges to %d.", l), f)
	default:
		n := w5RR(2, 11)
		w5RI(f, "dimension", n)
		return w5R(kp, fmt.Sprintf("What is the dimension of the finite-dimensional Banach space \\(\\mathbb{R}^{%d}\\)? (enter a number)", n), strconv.Itoa(n),
			fmt.Sprintf("The dimension is %d.", n), f)
	}
}

func w5Hilbert(kp int) generator.Problem {
	a, b := w5RNZ(-4, 4), w5RNZ(-4, 4)
	c, d := w5RNZ(-4, 4), w5RNZ(-4, 4)
	f := map[string]string{}
	w5RI(f, "a", a)
	w5RI(f, "b", b)
	w5RI(f, "c", c)
	w5RI(f, "d", d)
	switch kp {
	case 1:
		return w5R(kp, fmt.Sprintf("For \\(x=(%d,%d)\\) and \\(y=(%d,%d)\\), what is \\(\\langle x,y\\rangle\\)? (enter a number)", a, b, c, d), strconv.Itoa(a*c+b*d),
			fmt.Sprintf("The inner product is \\((%d)(%d)+(%d)(%d)=%d\\).", a, c, b, d, a*c+b*d), f)
	case 2:
		return w5R(kp, fmt.Sprintf("For \\(x=(%d,%d)\\) and \\(y=(%d,%d)\\), are they orthogonal? Enter 1 for yes and 0 for no.", a, b, c, d), strconv.Itoa(func() int {
			if a*c+b*d == 0 {
				return 1
			}
			return 0
		}()),
			"They are orthogonal exactly when the dot product is zero.", f)
	default:
		return w5R(kp, fmt.Sprintf("For \\(x=(%d,%d)\\), what is the square of its Hilbert norm? (enter a number)", a, b), strconv.Itoa(a*a+b*b),
			fmt.Sprintf("The squared norm is \\((%d)^2+(%d)^2=%d\\).", a, b, a*a+b*b), f)
	}
}

func w5BoundedOperator(kp int) generator.Problem {
	a, b := w5RNZ(-5, 5), w5RNZ(-5, 5)
	x, y := w5RNZ(-4, 4), w5RNZ(-4, 4)
	f := map[string]string{}
	w5RI(f, "a", a)
	w5RI(f, "b", b)
	w5RI(f, "x", x)
	w5RI(f, "y", y)
	switch kp {
	case 1:
		return w5R(kp, fmt.Sprintf("For the diagonal operator \\(T(x,y)=(%d x,%d y)\\), what is its operator norm? (enter a number)", a, b), strconv.Itoa(w5Max(w5Abs(a), w5Abs(b))),
			fmt.Sprintf("The operator norm is the largest diagonal magnitude, %d.", w5Max(w5Abs(a), w5Abs(b))), f)
	case 2:
		return w5R(kp, fmt.Sprintf("Is the finite matrix operator \\(T(x,y)=(%d x,%d y)\\) bounded on \\(\\mathbb{R}^2\\)? Enter 1 for yes and 0 for no.", a, b), "1",
			"Every finite matrix defines a bounded linear operator. The answer is 1.", f)
	default:
		return w5R(kp, fmt.Sprintf("For \\(T(x,y)=(%d x,%d y)\\), what is the first component of \\(T(%d,%d)\\)? (enter a number)", a, b, x, y), strconv.Itoa(a*x),
			fmt.Sprintf("The first component is \\((%d)(%d)=%d\\).", a, x, a*x), f)
	}
}

func w5Riesz(kp int) generator.Problem {
	a, b := w5RNZ(-4, 4), w5RNZ(-4, 4)
	x, y := w5RNZ(-3, 3), w5RNZ(-3, 3)
	f := map[string]string{}
	w5RI(f, "a", a)
	w5RI(f, "b", b)
	w5RI(f, "x", x)
	w5RI(f, "y", y)
	switch kp {
	case 1:
		return w5R(kp, fmt.Sprintf("The functional \\(f(x,y)=%d x+%d y\\) is represented by which vector? What is its first component? (enter a number)", a, b), strconv.Itoa(a),
			fmt.Sprintf("The representative is \\((%d,%d)\\), so the first component is %d.", a, b, a), f)
	case 2:
		return w5R(kp, fmt.Sprintf("For the functional \\(f(x,y)=%d x+%d y\\) and the point \\((%d,%d)\\), what is \\(f(%d,%d)\\)? (enter a number)", a, b, x, y, x, y), strconv.Itoa(a*x+b*y),
			fmt.Sprintf("Substitute the point: \\((%d)(%d)+(%d)(%d)=%d\\).", a, x, b, y, a*x+b*y), f)
	default:
		return w5R(kp, fmt.Sprintf("The functional \\(f(x,y)=%d x+%d y\\) is represented by \\((%d,%d)\\). What is the square of the representative's norm? (enter a number)", a, b, a, b), strconv.Itoa(a*a+b*b),
			fmt.Sprintf("The squared norm is \\((%d)^2+(%d)^2=%d\\).", a, b, a*a+b*b), f)
	}
}

func w5AbsoluteContinuity(kp int) generator.Problem {
	a, b, c, d := w5RNZ(-6, 6), w5RNZ(-6, 6), w5RNZ(-4, 4), w5RNZ(-4, 4)
	if c > d {
		c, d = d, c
	}
	f := map[string]string{}
	w5RI(f, "a", a)
	w5RI(f, "b", b)
	w5RI(f, "c", c)
	w5RI(f, "d", d)
	switch kp {
	case 1:
		return w5R(kp, fmt.Sprintf("For \\(f(x)=%d x+%d\\), what is \\(f(%d)-f(%d)\\)? (enter a number)", a, b, d, c), strconv.Itoa(a*(d-c)),
			fmt.Sprintf("The constant cancels and the difference is \\((%d)((%d)-(%d))=%d\\).", a, d, c, a*(d-c)), f)
	case 2:
		return w5R(kp, fmt.Sprintf("On the interval \\([%d,%d]\\), the Cantor function is continuous and nondecreasing, but its derivative is zero almost everywhere. Is it absolutely continuous there? Enter 1 for yes and 0 for no.", c, d), "0",
			"A function whose derivative is zero almost everywhere but which is not constant fails absolute continuity. The answer is 0.", f)
	default:
		return w5R(kp, fmt.Sprintf("For \\(f(x)=%d x+%d\\) on the interval \\([%d,%d]\\), what is its total variation? (enter a number)", a, b, c, d), strconv.Itoa(w5Abs(a)*(d-c)),
			fmt.Sprintf("A linear function has total variation \\(|%d|(%d-%d)=%d\\).", a, d, c, w5Abs(a)*(d-c)), f)
	}
}

func w5LebesgueDifferentiation(kp int) generator.Problem {
	a, b, c, d := w5RNZ(-6, 6), w5RNZ(-6, 6), w5RNZ(-4, 4), w5RNZ(-4, 4)
	if c > d {
		c, d = d, c
	}
	if (c+d)%2 != 0 {
		d++
	}
	f := map[string]string{}
	w5RI(f, "a", a)
	w5RI(f, "b", b)
	w5RI(f, "c", c)
	w5RI(f, "d", d)
	switch kp {
	case 1:
		return w5R(kp, fmt.Sprintf("For \\(f(x)=%d x+%d\\), what is its average over \\([%d,%d]\\)? (enter a number)", a, b, c, d), strconv.Itoa(a*(c+d)/2+b),
			fmt.Sprintf("The average of a linear function is its midpoint value, \\((%d)((%d)+(%d))/2+(%d)=%d\\).", a, c, d, b, a*(c+d)/2+b), f)
	case 2:
		return w5R(kp, fmt.Sprintf("For \\(f(t)=%d\\), what is the derivative of \\(F(x)=\\int_0^x f(t)\\,dt\\)? (enter a number)", a), strconv.Itoa(a),
			fmt.Sprintf("The fundamental theorem gives \\(F'(x)=f(x)=%d\\).", a), f)
	default:
		return w5R(kp, fmt.Sprintf("For \\(f(x)=%d x+%d\\), what is the value at the Lebesgue point \\(x=%d\\)? (enter a number)", a, b, c), strconv.Itoa(a*c+b),
			fmt.Sprintf("The value is \\((%d)(%d)+(%d)=%d\\).", a, c, b, a*c+b), f)
	}
}
