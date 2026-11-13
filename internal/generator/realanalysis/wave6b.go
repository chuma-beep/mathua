package realanalysis

// Wave 6B real analysis: fifteen concepts from the measure-theory and
// functional-analysis half of the course. Every family draws small integer
// data, so the answer is an integer and Facts carries every operand the
// answer follows from, including the answer itself.

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave6B(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"ra.deriv.cantor_function":          "cantor_function",
		"ra.measure.outer":                  "outer_measure",
		"ra.measure.caratheodory":           "caratheodory",
		"ra.measure.radon_nikodym":          "radon_nikodym",
		"ra.measure.signed":                 "signed",
		"ra.measure.convergence_in_measure": "convergence_in_measure",
		"ra.measure.egorov":                 "egorov",
		"ra.measure.lusin":                  "lusin",
		"ra.metric.cauchy_schwarz":          "cauchy_schwarz",
		"ra.metric.orthogonal_projection":   "orthogonal_projection",
		"ra.metric.uniform_boundedness":     "uniform_boundedness",
		"ra.metric.open_mapping":            "open_mapping",
		"ra.metric.closed_graph":            "closed_graph",
		"ra.metric.hahn_banach":             "hahn_banach",
		"ra.metric.weak_convergence":        "weak_convergence",
	} {
		reg.Register(id, &w6bGen{kind: kind})
	}
}

type w6bGen struct{ kind string }

func (g *w6bGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "cantor_function":
		return w6bCantor(kp)
	case "outer_measure":
		return w6bOuterMeasure(kp)
	case "caratheodory":
		return w6bCaratheodory(kp)
	case "radon_nikodym":
		return w6bRadonNikodym(kp)
	case "signed":
		return w6bSigned(kp)
	case "convergence_in_measure":
		return w6bConvergenceInMeasure(kp)
	case "egorov":
		return w6bEgorov(kp)
	case "lusin":
		return w6bLusin(kp)
	case "cauchy_schwarz":
		return w6bCauchySchwarz(kp)
	case "orthogonal_projection":
		return w6bProjection(kp)
	case "uniform_boundedness":
		return w6bUniformBoundedness(kp)
	case "open_mapping":
		return w6bOpenMapping(kp)
	case "closed_graph":
		return w6bClosedGraph(kp)
	case "hahn_banach":
		return w6bHahnBanach(kp)
	case "weak_convergence":
		return w6bWeakConvergence(kp)
	}
	return generator.Problem{Question: "What is 2 + 2? (enter a number)", Answer: "4", Explanation: "2 + 2 = 4."}
}

func w6bR(o int, q, a, ex string, facts map[string]string) generator.Problem {
	if facts == nil {
		facts = map[string]string{}
	}
	facts["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: facts, KPOrdinal: o}
}

func w6bI(m map[string]string, k string, v int) { m[k] = strconv.Itoa(v) }

func w6bRR(lo, hi int) int {
	if hi <= lo {
		return lo
	}
	return lo + rand.Intn(hi-lo+1)
}

func w6bRNZ(lo, hi int) int {
	if hi <= lo {
		return lo
	}
	for {
		v := w6bRR(lo, hi)
		if v != 0 {
			return v
		}
	}
}

func w6bAbs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func w6bMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func w6bMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func w6bPow(base, exp int) int {
	out := 1
	for i := 0; i < exp; i++ {
		out *= base
	}
	return out
}

func w6bJoin(xs []int) string {
	parts := make([]string, len(xs))
	for i, x := range xs {
		parts[i] = strconv.Itoa(x)
	}
	return strings.Join(parts, ",")
}

func w6bDigits(xs []int) string {
	var b strings.Builder
	for _, x := range xs {
		b.WriteString(strconv.Itoa(x))
	}
	return b.String()
}

func w6bCn(n int) string {
	if n < 0 {
		return fmt.Sprintf("(%d)", n)
	}
	return strconv.Itoa(n)
}

// w6bCantor asks for the numerator of C(m/3^k) over 2^k, for a ternary
// endpoint whose digits are 0 and 2 only. On the Cantor set those digits map
// to the binary digits 0 and 1, so the value is the binary number they spell.
func w6bCantor(kp int) generator.Problem {
	f := map[string]string{}
	var k int
	switch kp {
	case 1:
		k = w6bRR(2, 4)
	case 2:
		k = w6bRR(3, 5)
	default:
		k = w6bRR(4, 6)
	}
	digits := make([]int, k)
	bits := make([]int, k)
	m, a := 0, 0
	for i := 0; i < k; i++ {
		d := 2 * w6bRR(0, 1)
		digits[i] = d
		bits[i] = d / 2
		m = m*3 + d
		a = a*2 + d/2
	}
	den3 := w6bPow(3, k)
	den2 := w6bPow(2, k)
	w6bI(f, "k", k)
	w6bI(f, "m", m)
	w6bI(f, "denominator", den3)
	w6bI(f, "binary_denominator", den2)
	f["ternary_digits"] = w6bJoin(digits)
	f["binary_digits"] = w6bJoin(bits)
	q := fmt.Sprintf("Let \\(C\\) be the Cantor function. The ternary endpoint \\(x=%d/%d\\) has ternary digits 0.%s. If \\(C(x)=A/%d\\), what is \\(A\\)? (enter a number)", m, den3, w6bDigits(digits), den2)
	ex := fmt.Sprintf("On the Cantor set a ternary digit 0 becomes the binary digit 0, and a ternary digit 2 becomes the binary digit 1. The digits %s give the binary digits %s, so \\(C(%d/%d)=%d/%d\\). The numerator is %d.", w6bDigits(digits), w6bDigits(bits), m, den3, a, den2, a)
	return w6bR(kp, q, strconv.Itoa(a), ex, f)
}

func w6bOuterMeasure(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		a := w6bRR(-4, 2)
		b := a + w6bRR(1, 18)
		w6bI(f, "a", a)
		w6bI(f, "b", b)
		ans := b - a
		return w6bR(kp, fmt.Sprintf("What is the outer measure \\(m^*([%d,%d])\\)? (enter a number)", a, b), strconv.Itoa(ans),
			fmt.Sprintf("The outer measure of a bounded interval is its length: \\(%d-%d=%d\\).", b, a, ans), f)
	case 2:
		a := w6bRR(-4, 2)
		b := a + w6bRR(1, 18)
		n := w6bRR(1, 9)
		w6bI(f, "a", a)
		w6bI(f, "b", b)
		w6bI(f, "points", n)
		ans := b - a
		return w6bR(kp, fmt.Sprintf("What is the outer measure of \\([%d,%d]\\) together with a set of %d points? (enter a number)", a, b, n), strconv.Itoa(ans),
			fmt.Sprintf("A finite set is countable, and every countable set has outer measure 0. The union therefore has the measure of the interval, \\(%d-%d=%d\\).", b, a, ans), f)
	default:
		a := w6bRR(-5, 0)
		c := a + w6bRR(1, 4)
		b := c + w6bRR(1, 5)
		d := b + w6bRR(1, 5)
		w6bI(f, "a", a)
		w6bI(f, "b", b)
		w6bI(f, "c", c)
		w6bI(f, "d", d)
		ans := d - a
		return w6bR(kp, fmt.Sprintf("What is \\(m^*([%d,%d]\\cup[%d,%d])\\) when the two intervals overlap? (enter a number)", a, b, c, d), strconv.Itoa(ans),
			fmt.Sprintf("Because \\(%d<%d<%d<%d\\), the union is the single interval \\([%d,%d]\\), whose measure is \\(%d-%d=%d\\).", a, c, b, d, a, d, d, a, ans), f)
	}
}

func w6bCaratheodory(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		a := w6bRR(-4, 0)
		c := a + w6bRR(1, 18)
		b := c + w6bRR(1, 10)
		w6bI(f, "a", a)
		w6bI(f, "b", b)
		w6bI(f, "c", c)
		ans := b - c
		return w6bR(kp, fmt.Sprintf("Carathéodory's criterion writes \\(m^*(E)=m^*(E\\cap A)+m^*(E\\cap A^c)\\). For \\(E=[%d,%d]\\) and \\(A=[%d,%d]\\), what is \\(m^*(E\\cap A^c)\\)? (enter a number)", a, b, a, c), strconv.Itoa(ans),
			fmt.Sprintf("\\(E\\cap A^c=[%d,%d]\\), which has measure \\(%d-%d=%d\\).", c, b, b, c, ans), f)
	case 2:
		a := w6bRR(-5, 0)
		b := a + w6bRR(1, 10)
		c := b + w6bRR(1, 10)
		d := c + w6bRR(1, 10)
		w6bI(f, "a", a)
		w6bI(f, "b", b)
		w6bI(f, "c", c)
		w6bI(f, "d", d)
		ans := (b - a) + (d - c)
		return w6bR(kp, fmt.Sprintf("For the test set \\(E=[%d,%d]\\) and the measurable set \\(A=[%d,%d]\\), what is \\(m^*(E\\cap A^c)\\)? (enter a number)", a, d, b, c), strconv.Itoa(ans),
			fmt.Sprintf("\\(E\\cap A^c\\) is the two-piece set \\([%d,%d]\\cup[%d,%d]\\), so its measure is \\((%d-%d)+(%d-%d)=%d\\).", a, b, c, d, b, a, d, c, ans), f)
	default:
		l1 := w6bRR(3, 12)
		x1 := w6bRR(1, l1-1)
		l2 := w6bRR(3, 12)
		x2 := w6bRR(1, l2-1)
		w6bI(f, "l1", l1)
		w6bI(f, "l2", l2)
		w6bI(f, "x1", x1)
		w6bI(f, "x2", x2)
		ans := (l1 - x1) + (l2 - x2)
		return w6bR(kp, fmt.Sprintf("A test set \\(E\\) is the disjoint union of two intervals of lengths %d and %d. A measurable set \\(A\\) covers %d units of the first interval and %d units of the second. What is \\(m^*(E\\cap A^c)\\)? (enter a number)", l1, l2, x1, x2), strconv.Itoa(ans),
			fmt.Sprintf("The uncovered lengths are \\(%d-%d\\) and \\(%d-%d\\), and they add to \\(%d+%d=%d\\).", l1, x1, l2, x2, l1-x1, l2-x2, ans), f)
	}
}

func w6bRadonNikodym(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		rho := w6bRR(1, 9)
		l := w6bRR(1, 9)
		a := w6bRR(-3, 3)
		w6bI(f, "density", rho)
		w6bI(f, "length", l)
		w6bI(f, "a", a)
		ans := rho * l
		return w6bR(kp, fmt.Sprintf("For the measure \\(\\nu(A)=\\int_A %d\\,dx\\), what is \\(\\nu([%d,%d])\\)? (enter a number)", rho, a, a+l), strconv.Itoa(ans),
			fmt.Sprintf("A constant density multiplies the length: \\(%d\\cdot%d=%d\\).", rho, l, ans), f)
	case 2:
		r1 := w6bRR(1, 6)
		l1 := w6bRR(1, 6)
		r2 := w6bRR(1, 6)
		l2 := w6bRR(1, 6)
		w6bI(f, "density1", r1)
		w6bI(f, "length1", l1)
		w6bI(f, "density2", r2)
		w6bI(f, "length2", l2)
		ans := r1*l1 + r2*l2
		return w6bR(kp, fmt.Sprintf("A density equals %d on an interval of length %d and %d on a disjoint interval of length %d. What is the total mass? (enter a number)", r1, l1, r2, l2), strconv.Itoa(ans),
			fmt.Sprintf("Integrate the step density piece by piece: \\(%d\\cdot%d+%d\\cdot%d=%d\\).", r1, l1, r2, l2, ans), f)
	default:
		r1, r2, r3 := w6bRR(1, 5), w6bRR(1, 5), w6bRR(1, 5)
		l1, l2, l3 := w6bRR(1, 5), w6bRR(1, 5), w6bRR(1, 5)
		w6bI(f, "density1", r1)
		w6bI(f, "length1", l1)
		w6bI(f, "density2", r2)
		w6bI(f, "length2", l2)
		w6bI(f, "density3", r3)
		w6bI(f, "length3", l3)
		ans := r1*l1 + r2*l2 + r3*l3
		return w6bR(kp, fmt.Sprintf("A step density equals %d, %d and %d on three disjoint intervals of lengths %d, %d and %d. What is the total mass? (enter a number)", r1, r2, r3, l1, l2, l3), strconv.Itoa(ans),
			fmt.Sprintf("Add the three pieces: \\(%d\\cdot%d+%d\\cdot%d+%d\\cdot%d=%d\\).", r1, l1, r2, l2, r3, l3, ans), f)
	}
}

func w6bSigned(kp int) generator.Problem {
	f := map[string]string{}
	var masses []int
	switch kp {
	case 1:
		masses = []int{w6bRNZ(-9, 9), w6bRNZ(-9, 9)}
	case 2:
		masses = []int{w6bRNZ(-7, 7), w6bRNZ(-7, 7), w6bRNZ(-7, 7)}
	default:
		masses = []int{w6bRNZ(-9, 9), w6bRNZ(-9, 9), w6bRNZ(-9, 9)}
	}
	ans := 0
	for i, m := range masses {
		w6bI(f, fmt.Sprintf("m%d", i+1), m)
		ans += w6bAbs(m)
	}
	w6bI(f, "atoms", len(masses))
	f["masses"] = w6bJoin(masses)
	parts := make([]string, len(masses))
	for i, m := range masses {
		parts[i] = strconv.Itoa(w6bAbs(m))
	}
	var q string
	if len(masses) == 2 {
		q = fmt.Sprintf("A signed measure is supported on two atoms with masses %d and %d. What is its total variation? (enter a number)", masses[0], masses[1])
	} else {
		q = fmt.Sprintf("A signed measure is supported on three atoms with masses %d, %d and %d. What is its total variation? (enter a number)", masses[0], masses[1], masses[2])
	}
	ex := fmt.Sprintf("The total variation adds the absolute masses: \\(%s=%d\\).", strings.Join(parts, "+"), ans)
	return w6bR(kp, q, strconv.Itoa(ans), ex, f)
}

// w6bConvergenceInMeasure uses f_n(x)=cx/n, whose superlevel set is the
// half-line (en/c, infinity). The draws keep en/c integral so the measure of
// that half-line is an integer.
func w6bConvergenceInMeasure(kp int) generator.Problem {
	f := map[string]string{}
	n := w6bRR(2, 5)
	t := w6bRR(1, 4)
	c := n * t
	q0 := w6bRR(1, 15)
	e := t * q0
	r := w6bRR(1, 15)
	l := q0 + r
	switch kp {
	case 2:
		nn := w6bRR(2, 5)
		cc := w6bRR(2, 5)
		ee := w6bRR(1, 6)
		q := ee * nn / cc
		s := w6bRR(2, 16)
		bigN := q + s
		w6bI(f, "n", nn)
		w6bI(f, "c", cc)
		w6bI(f, "epsilon", ee)
		w6bI(f, "N", bigN)
		w6bI(f, "threshold", q)
		ans := bigN - q
		return w6bR(kp, fmt.Sprintf("On \\(\\{1,\\ldots,%d\\}\\) with counting measure, \\(f_n(k)=%d k/%d\\) and \\(\\epsilon=%d\\). What is the measure of \\(\\{k: f_n(k)>\\epsilon\\}\\)? (enter a number)", bigN, cc, nn, ee), strconv.Itoa(ans),
			fmt.Sprintf("The inequality \\(%d k/%d>%d\\) is \\(k>%d/%d\\), so the first %d integers satisfy it. That leaves \\(%d-%d=%d\\) points.", cc, nn, ee, ee*nn, cc, q, bigN, q, ans), f)
	case 3:
		w6bI(f, "n", n)
		w6bI(f, "c", c)
		w6bI(f, "epsilon", e)
		w6bI(f, "L", l)
		w6bI(f, "threshold", q0)
		return w6bR(kp, fmt.Sprintf("On \\([0,%d]\\), \\(f_n(x)=%d x/%d\\) and \\(\\epsilon=%d\\). What is the measure of \\(\\{x: f_n(x)<\\epsilon\\}\\)? (enter a number)", l, c, n, e), strconv.Itoa(q0),
			fmt.Sprintf("The inequality \\(%d x/%d<%d\\) is \\(x<%d\\), so the set is \\([0,%d)\\), of measure %d.", c, n, e, q0, q0, q0), f)
	default:
		w6bI(f, "n", n)
		w6bI(f, "c", c)
		w6bI(f, "epsilon", e)
		w6bI(f, "L", l)
		w6bI(f, "threshold", q0)
		ans := l - q0
		return w6bR(kp, fmt.Sprintf("On \\([0,%d]\\), \\(f_n(x)=%d x/%d\\) and \\(\\epsilon=%d\\). What is the measure of \\(\\{x: f_n(x)>\\epsilon\\}\\)? (enter a number)", l, c, n, e), strconv.Itoa(ans),
			fmt.Sprintf("The inequality \\(%d x/%d>%d\\) is \\(x>%d\\), so the set is \\((%d,%d]\\), of measure \\(%d-%d=%d\\).", c, n, e, q0, q0, l, l, q0, ans), f)
	}
}

func w6bEgorov(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		total := w6bRR(5, 15)
		removed := w6bRR(1, total-1)
		w6bI(f, "total", total)
		w6bI(f, "exceptional", removed)
		ans := total - removed
		return w6bR(kp, fmt.Sprintf("A set \\(E\\) has measure %d. Egorov's theorem gives a set \\(F\\subset E\\) with \\(m(E\\setminus F)=%d\\) on which the sequence converges uniformly. What is \\(m(F)\\)? (enter a number)", total, removed), strconv.Itoa(ans),
			fmt.Sprintf("Measures subtract on a subset: \\(%d-%d=%d\\).", total, removed, ans), f)
	case 2:
		m1 := w6bRR(3, 10)
		e1 := w6bRR(1, m1-1)
		m2 := w6bRR(1, 8)
		w6bI(f, "m1", m1)
		w6bI(f, "exceptional1", e1)
		w6bI(f, "m2", m2)
		ans := (m1 - e1) + m2
		return w6bR(kp, fmt.Sprintf("A set \\(E\\) is the disjoint union of \\(E_1\\) of measure %d and \\(E_2\\) of measure %d. On \\(E_1\\) an exceptional set of measure %d is removed, and \\(E_2\\) is kept whole. What is the measure of the set that remains? (enter a number)", m1, m2, e1), strconv.Itoa(ans),
			fmt.Sprintf("Keep \\(%d-%d=%d\\) from \\(E_1\\) and all %d from \\(E_2\\), giving \\(%d+%d=%d\\).", m1, e1, m1-e1, m2, m1-e1, m2, ans), f)
	default:
		total := w6bRR(8, 20)
		eps := w6bRR(1, total/2)
		w6bI(f, "total", total)
		w6bI(f, "epsilon", eps)
		ans := total - eps
		return w6bR(kp, fmt.Sprintf("A set \\(E\\) has measure %d. Egorov's theorem gives uniform convergence off a set of measure at most %d. What is the smallest integer lower bound on the measure of the set where convergence is uniform? (enter a number)", total, eps), strconv.Itoa(ans),
			fmt.Sprintf("The uniform set has measure at least \\(%d-%d=%d\\).", total, eps, ans), f)
	}
}

func w6bLusin(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		total := w6bRR(5, 15)
		removed := w6bRR(1, total-1)
		w6bI(f, "total", total)
		w6bI(f, "removed", removed)
		ans := total - removed
		return w6bR(kp, fmt.Sprintf("A measurable set \\(E\\) has measure %d. An open set of measure %d is removed. What is the measure of the closed set that remains? (enter a number)", total, removed), strconv.Itoa(ans),
			fmt.Sprintf("Lusin's theorem keeps a closed subset, and its measure is \\(%d-%d=%d\\).", total, removed, ans), f)
	case 2:
		count := w6bRR(2, 3)
		each := w6bRR(1, 4)
		slack := w6bRR(1, 15)
		total := count*each + slack
		w6bI(f, "total", total)
		w6bI(f, "intervals", count)
		w6bI(f, "each", each)
		w6bI(f, "removed", count*each)
		ans := slack
		return w6bR(kp, fmt.Sprintf("A measurable set of measure %d has %d disjoint open intervals removed, each of length %d. What is the measure of the closed set that remains? (enter a number)", total, count, each), strconv.Itoa(ans),
			fmt.Sprintf("The removed measure is \\(%d\\cdot%d=%d\\), so \\(%d-%d=%d\\) remains.", count, each, count*each, total, count*each, ans), f)
	default:
		r1, r2, r3 := w6bRR(1, 3), w6bRR(1, 3), w6bRR(1, 3)
		slack := w6bRR(1, 15)
		total := r1 + r2 + r3 + slack
		w6bI(f, "total", total)
		w6bI(f, "removed1", r1)
		w6bI(f, "removed2", r2)
		w6bI(f, "removed3", r3)
		ans := total - (r1 + r2 + r3)
		return w6bR(kp, fmt.Sprintf("A measurable set of measure %d has three disjoint open intervals removed, of lengths %d, %d and %d. What is the measure of the closed set that remains? (enter a number)", total, r1, r2, r3), strconv.Itoa(ans),
			fmt.Sprintf("Remove \\(%d+%d+%d=%d\\) from %d, leaving %d.", r1, r2, r3, r1+r2+r3, total, ans), f)
	}
}

func w6bCauchySchwarz(kp int) generator.Problem {
	f := map[string]string{}
	a, b := w6bRNZ(-4, 4), w6bRNZ(-4, 4)
	c, d := w6bRNZ(-4, 4), w6bRNZ(-4, 4)
	w6bI(f, "a", a)
	w6bI(f, "b", b)
	w6bI(f, "c", c)
	w6bI(f, "d", d)
	switch kp {
	case 1:
		ans := a*c + b*d
		return w6bR(kp, fmt.Sprintf("For \\(x=(%d,%d)\\) and \\(y=(%d,%d)\\), what is \\(\\langle x,y\\rangle\\)? (enter a number)", a, b, c, d), strconv.Itoa(ans),
			fmt.Sprintf("The inner product is \\((%d)(%d)+(%d)(%d)=%d\\).", a, c, b, d, ans), f)
	case 2:
		ans := (a+c)*(a+c) + (b+d)*(b+d)
		return w6bR(kp, fmt.Sprintf("For \\(x=(%d,%d)\\) and \\(y=(%d,%d)\\), what is \\(\\|x+y\\|^2\\)? (enter a number)", a, b, c, d), strconv.Itoa(ans),
			fmt.Sprintf("Add the vectors to get \\((%s,%s)\\), then square: \\((%s)^2+(%s)^2=%d\\).", w6bCn(a+c), w6bCn(b+d), w6bCn(a+c), w6bCn(b+d), ans), f)
	default:
		ans := (a*d - b*c) * (a*d - b*c)
		return w6bR(kp, fmt.Sprintf("For \\(x=(%d,%d)\\) and \\(y=(%d,%d)\\), what is \\(\\|x\\|^2\\|y\\|^2-\\langle x,y\\rangle^2\\)? (enter a number)", a, b, c, d), strconv.Itoa(ans),
			fmt.Sprintf("The Gram determinant equals \\((ad-bc)^2\\), here \\((%d\\cdot%d-%d\\cdot%d)^2=%d\\). It vanishes exactly when the vectors are parallel, which is the equality case of Cauchy-Schwarz.", a, d, b, c, ans), f)
	}
}

func w6bProjection(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		a, b := w6bRNZ(-4, 4), w6bRNZ(-4, 4)
		k := w6bRR(2, 16)
		w6bI(f, "a", a)
		w6bI(f, "b", b)
		w6bI(f, "k", k)
		return w6bR(kp, fmt.Sprintf("For \\(u=(%d,%d)\\) and \\(v=(%d,%d)\\), the projection of \\(v\\) onto \\(u\\) is \\(\\frac{v\\cdot u}{u\\cdot u}u\\). What is the coefficient \\(\\frac{v\\cdot u}{u\\cdot u}\\)? (enter a number)", a, b, k*a, k*b), strconv.Itoa(k),
			fmt.Sprintf("Here \\(v=%du\\), so \\(v\\cdot u=%d(u\\cdot u)\\) and the coefficient is %d.", k, k, k), f)
	case 2:
		a, b := w6bRNZ(-4, 4), w6bRNZ(-4, 4)
		c, d := w6bRNZ(-4, 4), w6bRNZ(-4, 4)
		w6bI(f, "a", a)
		w6bI(f, "b", b)
		w6bI(f, "c", c)
		w6bI(f, "d", d)
		ans := a*c + b*d
		return w6bR(kp, fmt.Sprintf("For \\(u=(%d,%d)\\) and \\(v=(%d,%d)\\), the coefficient of the projection of \\(v\\) onto \\(u\\) is \\(\\frac{v\\cdot u}{u\\cdot u}\\). What is the numerator \\(v\\cdot u\\)? (enter a number)", a, b, c, d), strconv.Itoa(ans),
			fmt.Sprintf("The numerator is the inner product \\((%d)(%d)+(%d)(%d)=%d\\).", a, c, b, d, ans), f)
	default:
		a, b := w6bRNZ(-3, 3), w6bRNZ(-3, 3)
		k := w6bRR(2, 5)
		w6bI(f, "a", a)
		w6bI(f, "b", b)
		w6bI(f, "k", k)
		ans := k * k * (a*a + b*b)
		return w6bR(kp, fmt.Sprintf("For \\(u=(%d,%d)\\) and \\(v=(%d,%d)\\), the projection of \\(v\\) onto \\(u\\) is \\(v\\) itself. What is the square of its length? (enter a number)", a, b, k*a, k*b), strconv.Itoa(ans),
			fmt.Sprintf("The projection is \\(v=%du\\), so its squared length is \\(%d^2((%d)^2+(%d)^2)=%d\\).", k, k, a, b, ans), f)
	}
}

func w6bUniformBoundedness(kp int) generator.Problem {
	f := map[string]string{}
	p := w6bRR(2, 15)
	q := w6bRR(2, 15)
	r := w6bRR(2, 15)
	w6bI(f, "bound1", p)
	w6bI(f, "bound2", q)
	w6bI(f, "bound3", r)
	switch kp {
	case 1:
		return w6bR(kp, fmt.Sprintf("Each operator \\(T_n\\) on \\(\\mathbb{R}\\) satisfies \\(|T_n(1)|\\le %d\\) for every \\(n\\). What is the smallest integer \\(B\\) with \\(\\|T_n\\|\\le B\\) for all \\(n\\)? (enter a number)", p), strconv.Itoa(p),
			fmt.Sprintf("In one dimension \\(\\|T_n\\|=|T_n(1)|\\), so the pointwise bound %d is the uniform bound.", p), f)
	case 2:
		ans := w6bMax(p, q)
		return w6bR(kp, fmt.Sprintf("Each diagonal operator \\(T_n(x,y)=(\\alpha_n x,\\beta_n y)\\) on \\(\\mathbb{R}^2\\) satisfies \\(|\\alpha_n|\\le %d\\) and \\(|\\beta_n|\\le %d\\) for every \\(n\\). What is the smallest integer \\(B\\) with \\(\\|T_n\\|\\le B\\) for all \\(n\\)? (enter a number)", p, q), strconv.Itoa(ans),
			fmt.Sprintf("The operator norm is the largest diagonal magnitude, and the pointwise bounds give \\(\\max(%d,%d)=%d\\).", p, q, ans), f)
	default:
		ans := w6bMax(w6bMax(p, q), r)
		return w6bR(kp, fmt.Sprintf("Each diagonal operator \\(T_n\\) on \\(\\mathbb{R}^3\\) has diagonal entries bounded by %d, %d and %d for every \\(n\\). What is the smallest integer \\(B\\) with \\(\\|T_n\\|\\le B\\) for all \\(n\\)? (enter a number)", p, q, r), strconv.Itoa(ans),
			fmt.Sprintf("The uniform bound is the largest of the three pointwise bounds: \\(\\max(%d,%d,%d)=%d\\).", p, q, r, ans), f)
	}
}

func w6bOpenMapping(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		a := w6bRNZ(-8, 8)
		r := w6bRR(1, 4)
		w6bI(f, "a", a)
		w6bI(f, "radius", r)
		ans := w6bAbs(a) * r
		return w6bR(kp, fmt.Sprintf("The map \\(T(x,y)=(%d x,%d y)\\) is open. The image of the ball of radius %d centred at the origin contains a ball of radius \\(s\\) centred at the origin. What is \\(s\\)? (enter a number)", a, a, r), strconv.Itoa(ans),
			fmt.Sprintf("The map scales every vector by \\(|%d|=%d\\), so the image ball has radius \\(%d\\cdot%d=%d\\).", a, w6bAbs(a), w6bAbs(a), r, ans), f)
	case 2:
		a := w6bRNZ(-16, 16)
		b := w6bRNZ(-16, 16)
		w6bI(f, "a", a)
		w6bI(f, "b", b)
		ans := w6bMin(w6bAbs(a), w6bAbs(b))
		return w6bR(kp, fmt.Sprintf("The map \\(T(x,y)=(%d x,%d y)\\) has inverse \\(T^{-1}\\) with \\(\\|T^{-1}\\|=1/m\\). What is \\(m\\)? (enter a number)", a, b), strconv.Itoa(ans),
			fmt.Sprintf("\\(T^{-1}(u,v)=(u/%d,v/%d)\\), so \\(\\|T^{-1}\\|=1/\\min(|%d|,|%d|)=1/%d\\) and \\(m=%d\\).", a, b, a, b, ans, ans), f)
	default:
		a := w6bRNZ(-6, 6)
		b := w6bRNZ(-6, 6)
		r := w6bRR(1, 3)
		w6bI(f, "a", a)
		w6bI(f, "b", b)
		w6bI(f, "radius", r)
		ans := w6bAbs(a*b) * r
		return w6bR(kp, fmt.Sprintf("Let \\(T(x,y)=(%d x,%d y)\\) and \\(S(x,y)=(%d x,%d y)\\). The image of the ball of radius %d under \\(S\\circ T\\) contains a ball of radius \\(s\\). What is \\(s\\)? (enter a number)", a, a, b, b, r), strconv.Itoa(ans),
			fmt.Sprintf("The composite scales by \\(|%d\\cdot%d|=%d\\), so \\(s=%d\\cdot%d=%d\\).", a, b, w6bAbs(a*b), w6bAbs(a*b), r, ans), f)
	}
}

func w6bClosedGraph(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		// The norm is max(|a|,|b|), so the answer space is only as wide as the
		// operand range. A 1..9 range gives nine answers, which the diversity
		// test can miss two of in 200 draws; the wider range keeps the question
		// varied without making the arithmetic harder than the concept needs.
		a := w6bRNZ(-19, 19)
		b := w6bRNZ(-19, 19)
		w6bI(f, "a", a)
		w6bI(f, "b", b)
		ans := w6bMax(w6bAbs(a), w6bAbs(b))
		return w6bR(kp, fmt.Sprintf("For the diagonal operator \\(T(x,y)=(%d x,%d y)\\), what is its operator norm? (enter a number)", a, b), strconv.Itoa(ans),
			fmt.Sprintf("The norm of a diagonal operator is the largest diagonal magnitude, \\(\\max(|%d|,|%d|)=%d\\).", a, b, ans), f)
	case 2:
		a, b := w6bRR(1, 4), w6bRR(1, 4)
		p, q := w6bRR(1, 5), w6bRR(1, 5)
		w6bI(f, "a", a)
		w6bI(f, "b", b)
		w6bI(f, "p", p)
		w6bI(f, "q", q)
		ans := (p + q) + (a*p + b*q)
		return w6bR(kp, fmt.Sprintf("For \\(T(x,y)=(%d x,%d y)\\) and the point \\((p,q)=(%d,%d)\\), the graph norm is \\(\\|(x,Tx)\\|_1=\\|x\\|_1+\\|Tx\\|_1\\). What is the graph norm of \\((%d,%d)\\)? (enter a number)", a, b, p, q, p, q), strconv.Itoa(ans),
			fmt.Sprintf("\\(\\|(%d,%d)\\|_1=%d\\) and \\(\\|T(%d,%d)\\|_1=%d+%d=%d\\). The graph norm is \\(%d+%d=%d\\).", p, q, p+q, p, q, a*p, b*q, a*p+b*q, p+q, a*p+b*q, ans), f)
	default:
		a := w6bRNZ(-9, 9)
		b := w6bRNZ(-9, 9)
		w6bI(f, "a", a)
		w6bI(f, "b", b)
		ans := w6bAbs(a * b)
		return w6bR(kp, fmt.Sprintf("Let \\(T(x,y)=(%d x,%d y)\\) and \\(S(x,y)=(%d x,%d y)\\). The graph of \\(S\\circ T\\) is closed. What is \\(\\|S\\circ T\\|\\)? (enter a number)", a, a, b, b), strconv.Itoa(ans),
			fmt.Sprintf("The composite is scaling by \\(%d\\cdot%d\\), so its norm is \\(|%d\\cdot%d|=%d\\).", a, b, a, b, ans), f)
	}
}

func w6bHahnBanach(kp int) generator.Problem {
	f := map[string]string{}
	a := w6bRNZ(-14, 14)
	b := w6bRNZ(-14, 14)
	p := w6bRR(-14, 14)
	q := w6bRR(-14, 14)
	w6bI(f, "a", a)
	w6bI(f, "b", b)
	w6bI(f, "p", p)
	w6bI(f, "q", q)
	switch kp {
	case 1:
		ans := w6bAbs(a)
		return w6bR(kp, fmt.Sprintf("On \\((\\mathbb{R}^2,\\|\\cdot\\|_1)\\) the functional \\(f(x,y)=%d x\\) has a norm-preserving Hahn-Banach extension. What is its norm? (enter a number)", a), strconv.Itoa(ans),
			fmt.Sprintf("The bound \\(|%d x|\\le %d(|x|+|y|)\\) is attained at \\((1,0)\\), so the norm is \\(|%d|=%d\\).", a, w6bAbs(a), a, ans), f)
	case 2:
		ans := w6bAbs(p)
		return w6bR(kp, fmt.Sprintf("Hahn-Banach separates the point \\((%d,%d)\\) from the hyperplane \\(\\{(x,y):x=0\\}\\). What is the separation distance? (enter a number)", p, q), strconv.Itoa(ans),
			fmt.Sprintf("The nearest point of the hyperplane is \\((0,%d)\\), so the distance is \\(|%d|=%d\\).", q, p, ans), f)
	default:
		ans := w6bAbs(a) + w6bAbs(b)
		return w6bR(kp, fmt.Sprintf("On \\((\\mathbb{R}^2,\\|\\cdot\\|_\\infty)\\) the functional \\(f(x,y)=%d x+%d y\\) has a norm-preserving Hahn-Banach extension. What is its norm? (enter a number)", a, b), strconv.Itoa(ans),
			fmt.Sprintf("On the \\(\\ell^\\infty\\) norm the dual norm is the \\(\\ell^1\\) norm: \\(|%d|+|%d|=%d\\).", a, b, ans), f)
	}
}

func w6bWeakConvergence(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		coords := []int{w6bRR(-9, 9), w6bRR(-9, 9), w6bRR(-9, 9), w6bRR(-9, 9), w6bRR(-9, 9)}
		n := w6bRR(1, 5)
		f["coords"] = w6bJoin(coords)
		w6bI(f, "N", n)
		ans := coords[n-1]
		return w6bR(kp, fmt.Sprintf("In \\(\\ell^2\\) the vector \\(x\\) has coordinates \\((%s)\\). What is \\(\\langle e_{%d},x\\rangle\\)? (enter a number)", w6bJoin(coords), n), strconv.Itoa(ans),
			fmt.Sprintf("The basis vector \\(e_{%d}\\) picks out the %d-th coordinate, which is %d.", n, n, ans), f)
	case 2:
		c := w6bRR(1, 3)
		r := w6bRR(1, 3)
		n := w6bRR(1, 4)
		w6bI(f, "c", c)
		w6bI(f, "ratio", r)
		w6bI(f, "N", n)
		ans := c * w6bPow(r, n-1)
		return w6bR(kp, fmt.Sprintf("In \\(\\mathbb{R}^4\\) the vector \\(x\\) has \\(k\\)-th coordinate \\(%d\\cdot%d^{k-1}\\). What is \\(\\langle e_{%d},x\\rangle\\)? (enter a number)", c, r, n), strconv.Itoa(ans),
			fmt.Sprintf("The functional \\(\\langle e_{%d},\\cdot\\rangle\\) reads the %d-th coordinate: \\(%d\\cdot%d^{%d}=%d\\).", n, n, c, r, n-1, ans), f)
	default:
		a1, a2, a3 := w6bRR(-4, 6), w6bRR(-4, 6), w6bRR(-4, 6)
		w6bI(f, "a1", a1)
		w6bI(f, "a2", a2)
		w6bI(f, "a3", a3)
		ans := a1 + a2 + a3
		return w6bR(kp, fmt.Sprintf("In \\(\\mathbb{R}^3\\) the sequence \\(x^{(n)}=(%d+1/n,\\ %d-1/n,\\ %d)\\) converges weakly to \\(x\\). What is \\(x_1+x_2+x_3\\)? (enter a number)", a1, a2, a3), strconv.Itoa(ans),
			fmt.Sprintf("Weak convergence in finite dimensions is coordinatewise, so \\(x=(%d,%d,%d)\\) and the sum is \\(%d+%d+%d=%d\\).", a1, a2, a3, a1, a2, a3, ans), f)
	}
}
