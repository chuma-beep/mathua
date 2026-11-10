package precalculus

// Wave 9 remaining precalculus: twenty concepts. Every family is an exact
// integer computation and publishes the operands it uses.

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

var w9rPureNumber = regexp.MustCompile(`\\\((-?\d+)\\\)`)
var w9rUnitCoefficient = regexp.MustCompile(`(^|[^0-9])1 ?([xyzuvtr])([^A-Za-z]|$)`)
var w9rUnitParen = regexp.MustCompile(`(^|[^0-9])1\(`)
var w9rSpaces = regexp.MustCompile(`[ \t]+`)

// w9rClean removes rendering defects that a learner should not have to read.
func w9rClean(s string) string {
	s = w9rPureNumber.ReplaceAllString(s, "$1")
	s = strings.ReplaceAll(s, "+-", "-")
	s = strings.ReplaceAll(s, "--", "+")
	s = w9rSpaces.ReplaceAllString(s, " ")
	s = w9rUnitCoefficient.ReplaceAllString(s, "${1}${2}${3}")
	s = w9rUnitParen.ReplaceAllString(s, "${1}(")
	return s
}

func registerWave9Remaining(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"precalc.func.inverse_composition":        "inverse_composition",
		"precalc.func.piecewise_operations":       "piecewise_operations",
		"precalc.analytic.area_polygon":           "area_polygon",
		"precalc.rational.sign_chart":             "sign_chart",
		"precalc.rational.asymptote_intersection": "asymptote_intersection",
		"precalc.exp.doubling":                    "doubling",
		"precalc.exp.logistic_inverse":            "logistic_inverse",
		"precalc.seq.summation_properties":        "summation_properties",
		"precalc.seq.power_sums":                  "power_sums",
		"precalc.conic.parametric":                "conic_parametric",
		"precalc.conic.polar":                     "conic_polar",
		"precalc.polar.symmetry":                  "polar_symmetry",
		"precalc.polar.intersection":              "polar_intersection",
		"precalc.parametric.cycloid":              "cycloid",
		"precalc.vector.work":                     "vector_work",
		"precalc.vector.direction_cosines":        "direction_cosines",
		"precalc.systems.nonlinear_inequality":    "nonlinear_inequality",
		"precalc.partial_fractions.repeated":      "partial_repeated",
		"precalc.partial_fractions.quadratic":     "partial_quadratic",
		"precalc.binomial.multinomial":            "multinomial",
	} {
		reg.Register(id, &w9rPreGen{kind: kind})
	}
}

type w9rPreGen struct{ kind string }

func (g *w9rPreGen) Generate(ctx generator.GeneratorContext) (prob generator.Problem) {
	defer func() {
		prob.Question = w9rClean(prob.Question)
		prob.Explanation = w9rClean(prob.Explanation)
	}()
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "inverse_composition":
		return w9rInverseComposition(kp)
	case "piecewise_operations":
		return w9rPiecewiseOperations(kp)
	case "area_polygon":
		return w9rAreaPolygon(kp)
	case "sign_chart":
		return w9rSignChart(kp)
	case "asymptote_intersection":
		return w9rAsymptoteIntersection(kp)
	case "doubling":
		return w9rDoubling(kp)
	case "logistic_inverse":
		return w9rLogisticInverse(kp)
	case "summation_properties":
		return w9rSummationProperties(kp)
	case "power_sums":
		return w9rPowerSums(kp)
	case "conic_parametric":
		return w9rConicParametric(kp)
	case "conic_polar":
		return w9rConicPolar(kp)
	case "polar_symmetry":
		return w9rPolarSymmetry(kp)
	case "polar_intersection":
		return w9rPolarIntersection(kp)
	case "cycloid":
		return w9rCycloid(kp)
	case "vector_work":
		return w9rVectorWork(kp)
	case "direction_cosines":
		return w9rDirectionCosines(kp)
	case "nonlinear_inequality":
		return w9rNonlinearInequality(kp)
	case "partial_repeated":
		return w9rPartialRepeated(kp)
	case "partial_quadratic":
		return w9rPartialQuadratic(kp)
	case "multinomial":
		return w9rMultinomial(kp)
	}
	return generator.Problem{Question: "What is 2 + 2? (enter a number)", Answer: "4", Explanation: "2 + 2 = 4."}
}

func w9rP(o int, q, a, ex string, facts map[string]string) generator.Problem {
	if facts == nil {
		facts = map[string]string{}
	}
	facts["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: facts, KPOrdinal: o}
}

func w9rPI(m map[string]string, k string, v int) { m[k] = strconv.Itoa(v) }

func w9rPR(lo, hi int) int {
	if hi <= lo {
		return lo
	}
	return lo + rand.Intn(hi-lo+1)
}

func w9rPNZ(lo, hi int) int {
	for {
		v := w9rPR(lo, hi)
		if v != 0 {
			return v
		}
	}
}

func w9rChoose(n, k int) int {
	if k < 0 || k > n {
		return 0
	}
	r := 1
	for i := 1; i <= k; i++ {
		r = r * (n - k + i) / i
	}
	return r
}

func w9rInverseComposition(kp int) generator.Problem {
	x := 2
	m := w9rPR(1, 4)
	b := w9rPR(-4, 4)
	n, c := w9rPR(1, 3), w9rPR(-4, 4)
	if kp == 3 && rand.Intn(2) == 0 {
		n, c = 1, -b
	}
	f := map[string]string{}
	w9rPI(f, "m", m)
	w9rPI(f, "b", b)
	w9rPI(f, "n", n)
	w9rPI(f, "c", c)
	w9rPI(f, "x", x)
	switch kp {
	case 1:
		ans := m*(n*x+c) + b
		return w9rP(kp, fmt.Sprintf("For \\(f(x)=%dx+%d\\) and \\(g(x)=%dx+%d\\), what is \\(f(g(%d))\\)? (enter a number)", m, b, n, c, x), strconv.Itoa(ans),
			fmt.Sprintf("\\(g(%d)=%d\\), so \\(f(g(%d))=%d\\).", x, n*x+c, x, ans), f)
	case 2:
		ans := n*(m*x+b) + c
		return w9rP(kp, fmt.Sprintf("For \\(f(x)=%dx+%d\\) and \\(g(x)=%dx+%d\\), what is \\(g(f(%d))\\)? (enter a number)", m, b, n, c, x), strconv.Itoa(ans),
			fmt.Sprintf("\\(f(%d)=%d\\), so \\(g(f(%d))=%d\\).", x, m*x+b, x, ans), f)
	default:
		ok := 0
		if m*n == 1 && m*c+b == 0 {
			ok = 1
		}
		return w9rP(kp, fmt.Sprintf("Are \\(f(x)=%dx+%d\\) and \\(g(x)=%dx+%d\\) inverse functions? Enter 1 for yes and 0 for no.", m, b, n, c), strconv.Itoa(ok),
			fmt.Sprintf("They are inverses exactly when \\(mn=1\\) and \\(mc+b=0\\). Here \\(mn=%d\\) and \\(mc+b=%d\\), so the answer is %d.", m*n, m*c+b, ok), f)
	}
}

func w9rPiecewiseOperations(kp int) generator.Problem {
	x := w9rPR(-8, 8)
	a, b, c, d := w9rPNZ(-3, 3), w9rPR(-4, 4), w9rPNZ(-3, 3), w9rPR(-4, 4)
	e, f2, g, h := w9rPNZ(-3, 3), w9rPR(-4, 4), w9rPNZ(-3, 3), w9rPR(-4, 4)
	f := map[string]string{}
	for key, v := range map[string]int{"a": a, "b": b, "c": c, "d": d, "e": e, "f": f2, "g": g, "h": h, "x": x} {
		w9rPI(f, key, v)
	}
	fv := a*x + b
	gv := e*x + f2
	if x >= 0 {
		fv = c*x + d
	}
	if x >= 1 {
		gv = g*x + h
	}
	switch kp {
	case 1:
		return w9rP(kp, fmt.Sprintf("For the piecewise functions \\(f\\) and \\(g\\) defined by \\(f(x)=%dx+%d\\) for \\(x<0\\), \\(f(x)=%dx+%d\\) for \\(x\\ge0\\), \\(g(x)=%dx+%d\\) for \\(x<1\\) and \\(g(x)=%dx+%d\\) for \\(x\\ge1\\), what is \\(f(%d)+g(%d)\\)? (enter a number)", a, b, c, d, e, f2, g, h, x, x), strconv.Itoa(fv+gv),
			fmt.Sprintf("At \\(x=%d\\), \\(f(x)=%d\\) and \\(g(x)=%d\\). Their sum is %d.", x, fv, gv, fv+gv), f)
	case 2:
		return w9rP(kp, fmt.Sprintf("For the same piecewise functions, what is \\(f(%d)g(%d)\\)? (enter a number)", x, x), strconv.Itoa(fv*gv),
			fmt.Sprintf("At \\(x=%d\\), \\(f(x)=%d\\) and \\(g(x)=%d\\). Their product is %d.", x, fv, gv, fv*gv), f)
	default:
		code := 3
		if x < 0 {
			code = 1
		} else if x < 1 {
			code = 2
		}
		return w9rP(kp, fmt.Sprintf("For the same piecewise functions, which branch pair applies at \\(x=%d\\)? Enter 1 for \\(x<0\\), 2 for \\(0\\le x<1\\) and 3 for \\(x\\ge1\\).", x), strconv.Itoa(code),
			fmt.Sprintf("The value \\(x=%d\\) lies in branch %d.", x, code), f)
	}
}

func w9rAreaPolygon(kp int) generator.Problem {
	a := w9rPR(2, 9)
	b := w9rPR(2, 9)
	w, h := w9rPR(2, 9), w9rPR(2, 9)
	if a%2 == 1 {
		b = b * 2
	}
	f := map[string]string{}
	w9rPI(f, "a", a)
	w9rPI(f, "b", b)
	w9rPI(f, "w", w)
	w9rPI(f, "h", h)
	switch kp {
	case 1:
		return w9rP(kp, fmt.Sprintf("A triangle has vertices \\((0,0)\\), \\((%d,0)\\) and \\((0,%d)\\). What is its area? (enter a number)", a, b), strconv.Itoa(a*b/2),
			fmt.Sprintf("The legs have lengths %d and %d, so the area is \\(%d\\cdot%d/2=%d\\).", a, b, a, b, a*b/2), f)
	case 2:
		return w9rP(kp, fmt.Sprintf("A rectangle has vertices \\((0,0)\\), \\((%d,0)\\), \\((%d,%d)\\) and \\((0,%d)\\). What is its area? (enter a number)", w, w, h, h), strconv.Itoa(w*h),
			fmt.Sprintf("The rectangle has width %d and height %d, so its area is %d.", w, h, w*h), f)
	default:
		return w9rP(kp, fmt.Sprintf("For the triangle with vertices \\((0,0)\\), \\((%d,0)\\) and \\((0,%d)\\), what is twice its area (the shoelace determinant)? (enter a number)", a, b), strconv.Itoa(a*b),
			fmt.Sprintf("Twice the area is the absolute shoelace determinant \\(%d\\cdot%d=%d\\).", a, b, a*b), f)
	}
}

func w9rSignChart(kp int) generator.Problem {
	a := w9rPR(-3, 3)
	b := w9rPR(-3, 3)
	for b == a {
		b = w9rPR(-3, 3)
	}
	x := w9rPR(-5, 5)
	for x == a || x == b {
		x = w9rPR(-5, 5)
	}
	f := map[string]string{}
	w9rPI(f, "a", a)
	w9rPI(f, "b", b)
	w9rPI(f, "x", x)
	switch kp {
	case 1:
		return w9rP(kp, fmt.Sprintf("For \\(r(x)=(x-%d)/(x-%d)\\), what is the numerator at \\(x=%d\\)? (enter a number)", a, b, x), strconv.Itoa(x-a),
			fmt.Sprintf("The numerator is \\(x-%d=%d\\).", a, x-a), f)
	case 2:
		return w9rP(kp, fmt.Sprintf("For \\(r(x)=(x-%d)/(x-%d)\\), what is the denominator at \\(x=%d\\)? (enter a number)", a, b, x), strconv.Itoa(x-b),
			fmt.Sprintf("The denominator is \\(x-%d=%d\\), which is nonzero at the chosen test point.", b, x-b), f)
	default:
		sign := 1
		if (x-a)*(x-b) < 0 {
			sign = -1
		}
		return w9rP(kp, fmt.Sprintf("For \\(r(x)=(x-%d)/(x-%d)\\), what is the sign of \\(r(%d)\\)? Enter 1 for positive and -1 for negative.", a, b, x), strconv.Itoa(sign),
			fmt.Sprintf("The numerator is %d and the denominator is %d, so the sign is %d.", x-a, x-b, sign), f)
	}
}

func w9rAsymptoteIntersection(kp int) generator.Problem {
	m := w9rPR(1, 4)
	c := w9rPR(-4, 4)
	a := []int{-2, -1, 1, 2}[rand.Intn(4)]
	pad := w9rPR(1, 3)
	x0 := a + pad
	k := a * w9rPR(1, 4) * pad
	f := map[string]string{}
	w9rPI(f, "m", m)
	w9rPI(f, "c", c)
	w9rPI(f, "a", a)
	w9rPI(f, "k", k)
	w9rPI(f, "x0", x0)
	switch kp {
	case 1:
		off := -k / a
		return w9rP(kp, fmt.Sprintf("For \\(f(x)=%dx+%d+%d/(x-%d)\\), what is the vertical offset of \\(f(0)\\) from the oblique asymptote \\(y=%dx+%d\\)? (enter a number)", m, c, k, a, m, c), strconv.Itoa(off),
			fmt.Sprintf("The offset is \\(f(0)-c=%d/(0-%d)=%d\\).", k, a, off), f)
	case 2:
		return w9rP(kp, fmt.Sprintf("How many finite points satisfy \\(f(x)=%dx+%d\\) for \\(f(x)=%dx+%d+%d/(x-%d)\\)? (enter a number)", m, c, m, c, k, a), "0",
			"The equation reduces to the nonzero constant numerator over \\(x-a\\), so there is no finite solution.", f)
	default:
		off := k / (x0 - a)
		return w9rP(kp, fmt.Sprintf("For \\(f(x)=%dx+%d+%d/(x-%d)\\), how far above the asymptote is \\(f(%d)\\)? (enter a number)", m, c, k, a, x0), strconv.Itoa(off),
			fmt.Sprintf("The offset is \\(%d/(%d-%d)=%d\\).", k, x0, a, off), f)
	}
}

func w9rDoubling(kp int) generator.Problem {
	A := w9rPR(2, 12)
	T := w9rPR(2, 10)
	k := w9rPR(1, 6)
	m := w9rPR(1, 7)
	f := map[string]string{}
	w9rPI(f, "A", A)
	w9rPI(f, "T", T)
	w9rPI(f, "k", k)
	w9rPI(f, "m", m)
	switch kp {
	case 1:
		ans := A * (1 << k)
		return w9rP(kp, fmt.Sprintf("A quantity starts at %d and doubles every %d years. What is its value after %d years? (enter a number)", A, T, k*T), strconv.Itoa(ans),
			fmt.Sprintf("%d years is %d doubling periods, so the value is \\(%d\\cdot2^{%d}=%d\\).", k*T, k, A, k, ans), f)
	case 2:
		return w9rP(kp, fmt.Sprintf("A quantity starts at %d and doubles every %d years. How many years pass during %d doublings? (enter a number)", A, T, k), strconv.Itoa(k*T),
			fmt.Sprintf("%d doublings take \\(%d\\cdot%d=%d\\) years.", k, k, T, k*T), f)
	default:
		return w9rP(kp, fmt.Sprintf("A quantity starts at %d and doubles every %d years. How many doublings are needed to reach %d times its starting value? (enter a number)", A, T, 1<<m), strconv.Itoa(m),
			fmt.Sprintf("The factor is \\(2^{%d}\\), so %d doublings are needed.", m, m), f)
	}
}

func w9rLogisticInverse(kp int) generator.Problem {
	z := w9rPR(1, 12)
	f := map[string]string{}
	w9rPI(f, "z", z)
	p := 1 << z
	switch kp {
	case 1:
		return w9rP(kp, fmt.Sprintf("For the logistic probability \\(p=2^{%d}/(2^{%d}+1)\\), what is the logit \\(\\ln(p/(1-p))\\) divided by \\(\\ln 2\\)? (enter a number)", z, z), strconv.Itoa(z),
			fmt.Sprintf("The odds are \\(2^{%d}\\), so the base-2 logit is %d.", z, z), f)
	case 2:
		return w9rP(kp, fmt.Sprintf("The base-2 logit is %d. If the logistic probability has the form \\(2^{z}/(2^{z}+1)\\), what is its numerator? (enter a number)", z), strconv.Itoa(p),
			fmt.Sprintf("The numerator is \\(2^{%d}=%d\\).", z, p), f)
	default:
		return w9rP(kp, fmt.Sprintf("The base-2 logit is %d. If the logistic probability has the form \\(2^{z}/(2^{z}+1)\\), what is its denominator? (enter a number)", z), strconv.Itoa(p+1),
			fmt.Sprintf("The denominator is \\(2^{%d}+1=%d\\).", z, p+1), f)
	}
}

func w9rSummationProperties(kp int) generator.Problem {
	a := w9rPR(1, 5)
	b := w9rPR(-4, 5)
	c := w9rPR(1, 4)
	d := w9rPR(-4, 5)
	n := w9rPR(2, 10)
	sumK := n * (n + 1) / 2
	f := map[string]string{}
	for key, v := range map[string]int{"a": a, "b": b, "c": c, "d": d, "n": n} {
		w9rPI(f, key, v)
	}
	switch kp {
	case 1:
		ans := a*sumK + b*n
		return w9rP(kp, fmt.Sprintf("What is \\(\\sum_{k=1}^{%d}(%d k%+d)\\)? (enter a number)", n, a, b), strconv.Itoa(ans),
			fmt.Sprintf("The sum is \\(%d\\cdot\\frac{%d\\cdot%d}{2}+%d\\cdot%d=%d\\).", a, n, n+1, b, n, ans), f)
	case 2:
		ans := c * a * sumK
		return w9rP(kp, fmt.Sprintf("What is \\(\\sum_{k=1}^{%d}%d\\cdot(%d k)\\)? (enter a number)", n, c, a), strconv.Itoa(ans),
			fmt.Sprintf("A constant multiple factors out: \\(%d\\cdot%d\\cdot\\frac{%d\\cdot%d}{2}=%d\\).", c, a, n, n+1, ans), f)
	default:
		ans := (a+c)*sumK + (b+d)*n
		return w9rP(kp, fmt.Sprintf("What is \\(\\sum_{k=1}^{%d}(%d k%+d+%d k%+d)\\)? (enter a number)", n, a, b, c, d), strconv.Itoa(ans),
			fmt.Sprintf("Combine like terms: \\((%d+%d)\\frac{%d\\cdot%d}{2}+(%d+%d)%d=%d\\).", a, c, n, n+1, b, d, n, ans), f)
	}
}

func w9rPowerSums(kp int) generator.Problem {
	n := w9rPR(2, 12)
	s1 := n * (n + 1) / 2
	f := map[string]string{}
	w9rPI(f, "n", n)
	switch kp {
	case 1:
		return w9rP(kp, fmt.Sprintf("What is \\(\\sum_{k=1}^{%d} k\\)? (enter a number)", n), strconv.Itoa(s1),
			fmt.Sprintf("The sum is \\(%d\\cdot%d/2=%d\\).", n, n+1, s1), f)
	case 2:
		ans := n * (n + 1) * (2*n + 1) / 6
		return w9rP(kp, fmt.Sprintf("What is \\(\\sum_{k=1}^{%d} k^2\\)? (enter a number)", n), strconv.Itoa(ans),
			fmt.Sprintf("Use \\(n(n+1)(2n+1)/6\\): the value is %d.", ans), f)
	default:
		return w9rP(kp, fmt.Sprintf("What is \\(\\sum_{k=1}^{%d} k^3\\)? (enter a number)", n), strconv.Itoa(s1*s1),
			fmt.Sprintf("The sum of cubes is the square of the sum: \\(%d^2=%d\\).", s1, s1*s1), f)
	}
}

func w9rConicParametric(kp int) generator.Problem {
	h := w9rPR(-3, 3)
	k := w9rPR(-3, 3)
	a := w9rPR(1, 5)
	b := w9rPR(1, 5)
	f := map[string]string{}
	for key, v := range map[string]int{"h": h, "k": k, "a": a, "b": b} {
		w9rPI(f, key, v)
	}
	switch kp {
	case 1:
		return w9rP(kp, fmt.Sprintf("For \\(x=%d+%d\\cos t\\) and \\(y=%d+%d\\sin t\\), what is \\(x\\) at \\(t=0\\)? (enter a number)", h, a, k, b), strconv.Itoa(h+a),
			fmt.Sprintf("\\(\\cos 0=1\\), so \\(x=%d+%d=%d\\).", h, a, h+a), f)
	case 2:
		return w9rP(kp, fmt.Sprintf("For \\(x=%d+%d\\cos t\\) and \\(y=%d+%d\\sin t\\), what is \\(y\\) at \\(t=\\pi/2\\)? (enter a number)", h, a, k, b), strconv.Itoa(k+b),
			fmt.Sprintf("\\(\\sin(\\pi/2)=1\\), so \\(y=%d+%d=%d\\).", k, b, k+b), f)
	default:
		return w9rP(kp, fmt.Sprintf("For \\(x=%d+%d\\cos t\\) and \\(y=%d+%d\\sin t\\), what is \\(x\\) at \\(t=\\pi\\)? (enter a number)", h, a, k, b), strconv.Itoa(h-a),
			fmt.Sprintf("\\(\\cos\\pi=-1\\), so \\(x=%d-%d=%d\\).", h, a, h-a), f)
	}
}

func w9rConicPolar(kp int) generator.Problem {
	en := []int{1, 1, 2, 3}[rand.Intn(4)]
	ed := 1
	if en == 1 && rand.Intn(2) == 0 {
		ed = 2
	}
	d := w9rPR(2, 9)
	f := map[string]string{}
	w9rPI(f, "e_num", en)
	w9rPI(f, "e_den", ed)
	w9rPI(f, "d", d)
	switch kp {
	case 1:
		code := 1
		if en > ed {
			code = 2
		} else if en < ed {
			code = 3
		}
		return w9rP(kp, fmt.Sprintf("A conic has polar equation \\(r=%d/(%d+%d\\cos\\theta)\\). Is it a parabola, hyperbola or ellipse? Enter 1 for parabola, 2 for hyperbola and 3 for ellipse.", en*d, ed, en), strconv.Itoa(code),
			fmt.Sprintf("The eccentricity is \\(e=%d/%d=%d.%d\\), so the conic is %s.", en, ed, en/ed, (en%ed)*10/ed, map[int]string{1: "a parabola", 2: "a hyperbola", 3: "an ellipse"}[code]), f)
	case 2:
		return w9rP(kp, fmt.Sprintf("For the polar equation \\(r=%d/(%d+%d\\cos\\theta)\\), what is the value of the denominator at \\(\\theta=0\\)? (enter a number)", en*d, ed, en), strconv.Itoa(ed+en),
			fmt.Sprintf("At \\(\\theta=0\\), \\(\\cos\\theta=1\\), so the denominator is \\(%d+%d=%d\\).", ed, en, ed+en), f)
	default:
		return w9rP(kp, fmt.Sprintf("A conic has polar equation \\(r=%d/(%d+%d\\cos\\theta)\\). What is the directrix distance \\(d\\)? (enter a number)", en*d, ed, en), strconv.Itoa(d),
			fmt.Sprintf("The constant numerator is \\(ed=%d\\cdot%d=%d\\), so \\(d=%d\\).", en, d, en*d, d), f)
	}
}

func w9rPolarSymmetry(kp int) generator.Problem {
	typ := w9rPR(1, 3)
	k := w9rPR(1, 4)
	f := map[string]string{}
	w9rPI(f, "type", typ)
	w9rPI(f, "k", k)
	eq := fmt.Sprintf("r=%d", k)
	if typ == 2 {
		eq = fmt.Sprintf("r=%d\\sin\\theta", k)
	} else if typ == 3 {
		eq = fmt.Sprintf("r=%d\\cos 2\\theta", k)
	}
	switch kp {
	case 1:
		sym := 1
		if typ == 2 {
			sym = 0
		}
		return w9rP(kp, fmt.Sprintf("For the polar curve \\(%s\\), is the curve symmetric about the x-axis? Enter 1 for yes and 0 for no.", eq), strconv.Itoa(sym),
			fmt.Sprintf("The symmetry test gives code %d for this curve.", sym), f)
	case 2:
		test := 1
		if typ == 2 {
			test = 0
		}
		return w9rP(kp, fmt.Sprintf("For the polar curve \\(%s\\), does replacing \\(\\theta\\) by \\(-\\theta\\) leave the equation unchanged? Enter 1 for yes and 0 for no.", eq), strconv.Itoa(test),
			fmt.Sprintf("The substitution gives code %d for this curve.", test), f)
	default:
		count := 1
		if typ == 1 || typ == 3 {
			count = 2
		}
		return w9rP(kp, fmt.Sprintf("For the polar curve \\(%s\\), how many coordinate-axis symmetries does it have? (enter a number)", eq), strconv.Itoa(count),
			fmt.Sprintf("The coordinate-axis symmetry count is %d.", count), f)
	}
}

func w9rPolarIntersection(kp int) generator.Problem {
	a := w9rPR(1, 12)
	f := map[string]string{}
	w9rPI(f, "a", a)
	switch kp {
	case 1:
		return w9rP(kp, fmt.Sprintf("How many intersection points do \\(r=%d\\) and \\(r=2\\cdot%d\\cos\\theta\\) have? (enter a number)", a, a), "2",
			"\\(\\cos\\theta=1/2\\) gives two angles, so there are two intersection points.", f)
	case 2:
		return w9rP(kp, fmt.Sprintf("For \\(r=%d\\) and \\(r=2\\cdot%d\\cos\\theta\\), what is the intersection angle in degrees between 0 and 90? (enter a number)", a, a), "60",
			"\\(\\cos\\theta=1/2\\) gives \\(\\theta=60^\\circ\\) in the first quadrant.", f)
	default:
		return w9rP(kp, fmt.Sprintf("For \\(r=%d\\) and \\(r=2\\cdot%d\\cos\\theta\\), what is the common radius at the first-quadrant intersection? (enter a number)", a, a), strconv.Itoa(a),
			fmt.Sprintf("The circle \\(r=%d\\) fixes the common radius at %d.", a, a), f)
	}
}

func w9rCycloid(kp int) generator.Problem {
	a := w9rPR(1, 12)
	f := map[string]string{}
	w9rPI(f, "a", a)
	switch kp {
	case 1:
		return w9rP(kp, fmt.Sprintf("A cycloid has \\(x=%d(t-\\sin t)\\), \\(y=%d(1-\\cos t)\\). What is its maximum height? (enter a number)", a, a), strconv.Itoa(2*a),
			fmt.Sprintf("The maximum of \\(1-\\cos t\\) is 2, so the height is \\(2\\cdot%d=%d\\).", a, 2*a), f)
	case 2:
		return w9rP(kp, fmt.Sprintf("For the cycloid \\(x=%d(t-\\sin t)\\), \\(y=%d(1-\\cos t)\\), what is \\(y\\) at \\(t=\\pi\\)? (enter a number)", a, a), strconv.Itoa(2*a),
			fmt.Sprintf("\\(\\cos\\pi=-1\\), so \\(y=%d(1-(-1))=%d\\).", a, 2*a), f)
	default:
		return w9rP(kp, fmt.Sprintf("For the cycloid \\(x=%d(t-\\sin t)\\), \\(y=%d(1-\\cos t)\\), the period is how many multiples of \\(\\pi %d\\)? (enter a number)", a, a, a), "2",
			"The period is \\(2\\pi a\\), so the multiplier is 2.", f)
	}
}

func w9rVectorWork(kp int) generator.Problem {
	p, q := w9rPR(-4, 5), w9rPR(-4, 5)
	a, b := 3, 4
	work := p*a + q*b
	for work%5 != 0 {
		p, q = w9rPR(-4, 5), w9rPR(-4, 5)
		work = p*a + q*b
	}
	f := map[string]string{}
	for key, v := range map[string]int{"p": p, "q": q, "a": a, "b": b} {
		w9rPI(f, key, v)
	}
	switch kp {
	case 1:
		return w9rP(kp, fmt.Sprintf("A constant force \\(F=(%d,%d)\\) acts through displacement \\(d=(%d,%d)\\). What is the work? (enter a number)", p, q, a, b), strconv.Itoa(work),
			fmt.Sprintf("Work is the dot product \\(%d\\cdot%d+%d\\cdot%d=%d\\).", p, a, q, b, work), f)
	case 2:
		return w9rP(kp, fmt.Sprintf("For \\(F=(%d,%d)\\) and \\(d=(%d,%d)\\), what is the scalar projection of \\(F\\) on \\(d\\)? (enter a number)", p, q, a, b), strconv.Itoa(work/5),
			fmt.Sprintf("The dot product is %d and \\(\\|d\\|=5\\), so the projection is \\(%d/5=%d\\).", work, work, work/5), f)
	default:
		return w9rP(kp, fmt.Sprintf("For \\(F=(%d,%d)\\) and \\(d=(%d,%d)\\), what is the work if the displacement is doubled? (enter a number)", p, q, a, b), strconv.Itoa(2*work),
			fmt.Sprintf("Doubling displacement doubles the dot product: \\(2(%d)=%d\\).", work, 2*work), f)
	}
}

func w9rDirectionCosines(kp int) generator.Problem {
	triples := [][4]int{{1, 2, 2, 3}, {2, 3, 6, 7}, {3, 4, 12, 13}}
	t := triples[rand.Intn(len(triples))]
	scale := w9rPR(1, 4)
	a, b, c, r := t[0]*scale, t[1]*scale, t[2]*scale, t[3]*scale
	f := map[string]string{}
	for key, v := range map[string]int{"a": a, "b": b, "c": c, "r": r} {
		w9rPI(f, key, v)
	}
	switch kp {
	case 1:
		return w9rP(kp, fmt.Sprintf("For the vector with components \\((%d,%d,%d)\\) and magnitude %d, what is the numerator of \\(\\cos\\beta\\) when written over %d? (enter a number)", a, b, c, r, r), strconv.Itoa(b),
			fmt.Sprintf("The second direction cosine is \\(%d/%d\\), so its numerator is %d.", b, r, b), f)
	case 2:
		return w9rP(kp, fmt.Sprintf("For the vector direction with components \\((%d,%d,%d)\\), what is the sum of the squares of the direction cosines times %d squared? (enter a number)", a, b, c, r), strconv.Itoa(r*r),
			fmt.Sprintf("The squares sum to 1, so multiplied by \\(r^2=%d\\) the result is %d.", r*r, r*r), f)
	default:
		return w9rP(kp, fmt.Sprintf("A vector has magnitude %d and first direction cosine \\(%d/%d\\). What is its first component? (enter a number)", r, a, r), strconv.Itoa(a),
			fmt.Sprintf("The first component is \\(%d\\cdot(%d/%d)=%d\\).", r, a, r, a), f)
	}
}

func w9rNonlinearInequality(kp int) generator.Problem {
	r := w9rPR(2, 8)
	px, py := w9rPR(-8, 8), w9rPR(-8, 8)
	dist := px*px + py*py
	f := map[string]string{}
	w9rPI(f, "r", r)
	w9rPI(f, "px", px)
	w9rPI(f, "py", py)
	switch kp {
	case 1:
		inside := 0
		if dist <= r*r {
			inside = 1
		}
		return w9rP(kp, fmt.Sprintf("Does the point \\((%d,%d)\\) satisfy \\(x^2+y^2\\le %d^2\\)? Enter 1 for yes and 0 for no.", px, py, r), strconv.Itoa(inside),
			fmt.Sprintf("The squared distance is \\(%d\\), compared with \\(r^2=%d\\). The answer is %d.", dist, r*r, inside), f)
	case 2:
		strict := 0
		if dist < r*r {
			strict = 1
		}
		return w9rP(kp, fmt.Sprintf("Does the point \\((%d,%d)\\) satisfy the strict inequality \\(x^2+y^2< %d^2\\)? Enter 1 for yes and 0 for no.", px, py, r), strconv.Itoa(strict),
			fmt.Sprintf("The squared distance is %d, compared with %d; the strict inequality holds? %d.", dist, r*r, strict), f)
	default:
		return w9rP(kp, fmt.Sprintf("What is the squared distance from \\((%d,%d)\\) to the origin? (enter a number)", px, py), strconv.Itoa(dist),
			fmt.Sprintf("The squared distance is \\((%d)^2+(%d)^2=%d\\).", px, py, dist), f)
	}
}

func w9rPartialRepeated(kp int) generator.Problem {
	a := 0
	for a == 0 {
		a = w9rPR(-3, 3)
	}
	A := w9rPR(1, 5)
	if rand.Intn(2) == 0 {
		A = -A
	}
	B := w9rPR(-5, 5)
	n0 := B - A*a
	f := map[string]string{}
	w9rPI(f, "a", a)
	w9rPI(f, "A", A)
	w9rPI(f, "B", B)
	w9rPI(f, "n0", n0)
	switch kp {
	case 1:
		return w9rP(kp, fmt.Sprintf("The decomposition \\((%d x+%d)/(x-%d)^2=%d/(x-%d)+C/(x-%d)^2\\) has what value of \\(C\\)? (enter a number)", A, n0, a, A, a, a), strconv.Itoa(B),
			fmt.Sprintf("The numerator is \\(%d(x-%d)+C\\). Matching the constant gives \\(%d=%d(-%d)+C\\), so \\(C=%d\\).", A, a, n0, A, a, B), f)
	case 2:
		return w9rP(kp, fmt.Sprintf("The decomposition \\((A x+%d)/(x-%d)^2=A/(x-%d)+%d/(x-%d)^2\\) has what value of \\(A\\)? (enter a number)", n0, a, a, B, a), strconv.Itoa(A),
			fmt.Sprintf("Matching the constant gives \\(%d=A(-%d)+%d\\), so \\(A=%d\\).", n0, a, B, A), f)
	default:
		val := A + n0
		return w9rP(kp, fmt.Sprintf("For the numerator \\(A x+%d\\) in the decomposition with \\(A=%d\\), what is its value at \\(x=1\\)? (enter a number)", n0, A), strconv.Itoa(val),
			fmt.Sprintf("At \\(x=1\\) the numerator is \\(%d+%d=%d\\).", A, n0, val), f)
	}
}

func w9rPartialQuadratic(kp int) generator.Problem {
	p := w9rPNZ(-3, 3)
	q := w9rPR(1, 8)
	n1 := w9rPR(-5, 5)
	n0 := w9rPR(-5, 5)
	f := map[string]string{}
	for key, v := range map[string]int{"p": p, "q": q, "n1": n1, "n0": n0} {
		w9rPI(f, key, v)
	}
	switch kp {
	case 1:
		return w9rP(kp, fmt.Sprintf("For the numerator \\(%d x+%d\\) over the irreducible quadratic \\(x^2+%d x+%d\\), what is the numerator value at \\(x=1\\)? (enter a number)", n1, n0, p, q), strconv.Itoa(n1+n0),
			fmt.Sprintf("At \\(x=1\\) the numerator is \\(%d+%d=%d\\).", n1, n0, n1+n0), f)
	case 2:
		return w9rP(kp, fmt.Sprintf("For the numerator \\(%d x+%d\\) over \\(x^2+%d x+%d\\), what is the numerator value at \\(x=0\\)? (enter a number)", n1, n0, p, q), strconv.Itoa(n0),
			fmt.Sprintf("At \\(x=0\\) the numerator is %d.", n0), f)
	default:
		return w9rP(kp, fmt.Sprintf("For the irreducible quadratic \\(x^2+%d x+%d\\), what is its discriminant? (enter a number)", p, q), strconv.Itoa(p*p-4*q),
			fmt.Sprintf("The discriminant is \\(%d^2-4\\cdot%d=%d\\).", p, q, p*p-4*q), f)
	}
}

func w9rMultinomial(kp int) generator.Problem {
	n := w9rPR(2, 6)
	a := w9rPR(0, n)
	b := w9rPR(0, n-a)
	c := n - a - b
	coef := 1
	for i := 2; i <= n; i++ {
		coef *= i
	}
	den := 1
	for i := 2; i <= a; i++ {
		den *= i
	}
	for i := 2; i <= b; i++ {
		den *= i
	}
	for i := 2; i <= c; i++ {
		den *= i
	}
	coef /= den
	f := map[string]string{}
	for key, v := range map[string]int{"n": n, "a": a, "b": b, "c": c} {
		w9rPI(f, key, v)
	}
	switch kp {
	case 1:
		return w9rP(kp, fmt.Sprintf("What is the coefficient of \\(x^{%d}y^{%d}z^{%d}\\) in \\((x+y+z)^{%d}\\)? (enter a number)", a, b, c, n), strconv.Itoa(coef),
			fmt.Sprintf("The coefficient is \\(%d!/(%d!%d!%d!)=%d\\).", n, a, b, c, coef), f)
	case 2:
		return w9rP(kp, fmt.Sprintf("How many terms does the expansion of \\((x+y+z)^{%d}\\) have? (enter a number)", n), strconv.Itoa(w9rChoose(n+2, 2)),
			fmt.Sprintf("The number of monomials of degree %d in three variables is \\(\\binom{%d}{2}=%d\\).", n, n+2, w9rChoose(n+2, 2)), f)
	default:
		return w9rP(kp, fmt.Sprintf("What is \\((1+1+1)^{%d}\\)? (enter a number)", n), strconv.Itoa(w9rPow3(n)),
			fmt.Sprintf("The sum of all coefficients is \\(3^{%d}=%d\\).", n, w9rPow3(n)), f)
	}
}

func w9rPow3(n int) int {
	r := 1
	for i := 0; i < n; i++ {
		r *= 3
	}
	return r
}
