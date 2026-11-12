package realanalysis

// Wave 6a real analysis: sixteen seeded integer families across order, sequence,
// limit, derivative and inverse-function concepts. Every family publishes the
// operands its integer answer follows from in Facts, every stem asks for a
// number, and every draw is built so the arithmetic can be checked by hand.

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave6A(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"ra.order.dedekind":         "dedekind",
		"ra.order.well_ordering":    "well_ordering",
		"ra.order.cardinal":         "cardinal",
		"ra.seq.root_test":          "root_test",
		"ra.seq.integral_test":      "integral_test",
		"ra.seq.cesaro":             "cesaro",
		"ra.seq.dirichlet":          "dirichlet",
		"ra.limit.oscillation":      "oscillation",
		"ra.limit.semicontinuity":   "semicontinuity",
		"ra.limit.dini":             "dini",
		"ra.limit.uniform_limit":    "uniform_limit",
		"ra.deriv.rolle":            "rolle",
		"ra.deriv.convex":           "convex",
		"ra.deriv.jensen":           "jensen",
		"ra.deriv.lipschitz":        "lipschitz",
		"ra.deriv.inverse_function": "inverse_function",
	} {
		reg.Register(id, &w6aGen{kind: kind})
	}
}

type w6aGen struct{ kind string }

func (g *w6aGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "dedekind":
		return w6aNoReadOff(func() generator.Problem { return w6aDedekind(kp) })
	case "well_ordering":
		return w6aNoReadOff(func() generator.Problem { return w6aWellOrdering(kp) })
	case "cardinal":
		return w6aNoReadOff(func() generator.Problem { return w6aCardinal(kp) })
	case "root_test":
		return w6aNoReadOff(func() generator.Problem { return w6aRootTest(kp) })
	case "integral_test":
		return w6aNoReadOff(func() generator.Problem { return w6aIntegralTest(kp) })
	case "cesaro":
		return w6aNoReadOff(func() generator.Problem { return w6aCesaro(kp) })
	case "dirichlet":
		return w6aNoReadOff(func() generator.Problem { return w6aDirichlet(kp) })
	case "oscillation":
		return w6aNoReadOff(func() generator.Problem { return w6aOscillation(kp) })
	case "semicontinuity":
		return w6aNoReadOff(func() generator.Problem { return w6aSemicontinuity(kp) })
	case "dini":
		return w6aNoReadOff(func() generator.Problem { return w6aDini(kp) })
	case "uniform_limit":
		return w6aNoReadOff(func() generator.Problem { return w6aUniformLimit(kp) })
	case "rolle":
		return w6aNoReadOff(func() generator.Problem { return w6aRolle(kp) })
	case "convex":
		return w6aNoReadOff(func() generator.Problem { return w6aConvex(kp) })
	case "jensen":
		return w6aNoReadOff(func() generator.Problem { return w6aJensen(kp) })
	case "lipschitz":
		return w6aNoReadOff(func() generator.Problem { return w6aLipschitz(kp) })
	case "inverse_function":
		return w6aNoReadOff(func() generator.Problem { return w6aInverseFunction(kp) })
	}
	return generator.Problem{Question: "What is 2 + 2? (enter a number)", Answer: "4", Explanation: "2 + 2 = 4.", KPOrdinal: kp}
}

// w6aNoReadOff redraws a bounded number of times when the integer answer would
// appear in the stem as a standalone token. A redraw is cheaper than a stem
// that hands the learner the number, and the ceiling keeps a family that can
// never satisfy the check from looping forever.
func w6aNoReadOff(build func() generator.Problem) generator.Problem {
	var p generator.Problem
	for i := 0; i < generator.MaxDrawAttempts; i++ {
		p = build()
		if !w6aLeaks(p.Question, p.Answer) {
			return p
		}
	}
	return p
}

func w6aLeaks(stem, ans string) bool {
	if ans == "" {
		return false
	}
	if len(ans) == 1 && ans[0] >= '0' && ans[0] <= '9' {
		return false
	}
	for i := 0; i+len(ans) <= len(stem); {
		j := strings.Index(stem[i:], ans)
		if j < 0 {
			return false
		}
		j += i
		before := j == 0 || !w6aWordByte(stem[j-1])
		after := j+len(ans) == len(stem) || !w6aWordByte(stem[j+len(ans)])
		if before && after {
			return true
		}
		i = j + 1
	}
	return false
}

func w6aWordByte(b byte) bool {
	return b >= '0' && b <= '9' || b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z' || b == '_'
}

func w6aR(kp int, q, a, ex string, facts map[string]string) generator.Problem {
	if facts == nil {
		facts = map[string]string{}
	}
	facts["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: facts, KPOrdinal: kp}
}

func w6aRI(m map[string]string, k string, v int) { m[k] = strconv.Itoa(v) }

func w6aRR(lo, hi int) int {
	if hi <= lo {
		return lo
	}
	return lo + rand.Intn(hi-lo+1)
}

func w6aRNZ(lo, hi int) int {
	if hi <= lo {
		return lo
	}
	for {
		if v := w6aRR(lo, hi); v != 0 {
			return v
		}
	}
}

func w6aMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func w6aMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func w6aAbs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// w6aDedekind: a rational cut below r. The least integer upper bound is
// ceil(r) and the greatest integer inside the cut is floor(r); both are asked
// for directly or combined, so the answer stays an integer.
func w6aDedekind(kp int) generator.Problem {
	f := map[string]string{}
	k := w6aRR(-4, 11)
	switch kp {
	case 1:
		num := 2*k + 1
		w6aRI(f, "k", k)
		w6aRI(f, "numerator", num)
		w6aRI(f, "denominator", 2)
		ans := k + 1
		return w6aR(kp, fmt.Sprintf("A Dedekind cut is the set of all rationals \\(q\\) with \\(q<\\frac{%d}{2}\\). What is the least integer upper bound of the cut? (enter a number)", num), strconv.Itoa(ans),
			fmt.Sprintf("The cut lies below \\(\\frac{%d}{2}=%d+\\frac{1}{2}\\), so the least integer at or above the cut is %d.", num, k, ans), f)
	case 2:
		num := 3*k + 1
		w6aRI(f, "k", k)
		w6aRI(f, "numerator", num)
		w6aRI(f, "denominator", 3)
		ans := k
		return w6aR(kp, fmt.Sprintf("A Dedekind cut is the set of all rationals \\(q\\) with \\(q<\\frac{%d}{3}\\). What is the greatest integer that belongs to the cut? (enter a number)", num), strconv.Itoa(ans),
			fmt.Sprintf("The cut lies below \\(\\frac{%d}{3}=%d+\\frac{1}{3}\\), so the greatest integer below the cut is %d.", num, k, ans), f)
	default:
		num := 4*k + 3
		w6aRI(f, "k", k)
		w6aRI(f, "numerator", num)
		w6aRI(f, "denominator", 4)
		ans := 2*k + 1
		return w6aR(kp, fmt.Sprintf("A Dedekind cut is the set of all rationals \\(q\\) with \\(q<\\frac{%d}{4}\\). What is the sum of the greatest integer in the cut and the least integer upper bound of the cut? (enter a number)", num), strconv.Itoa(ans),
			fmt.Sprintf("The cut lies below \\(\\frac{%d}{4}=%d+\\frac{3}{4}\\). Its greatest integer is %d and its least integer upper bound is %d, so the sum is \\(%d+%d=%d\\).", num, k, k, k+1, k, k+1, ans), f)
	}
}

// w6aWellOrdering: a finite well-order is a listing, so ranks are positions in
// the listing. Values are drawn above twenty so a rank can never coincide with
// a printed value.
func w6aWellOrdering(kp int) generator.Problem {
	f := map[string]string{}
	n := w6aRR(8, 12)
	perm := rand.Perm(26)[:n]
	vals := make([]int, n)
	parts := make([]string, n)
	for i := range vals {
		vals[i] = perm[i] + 20
		parts[i] = strconv.Itoa(vals[i])
	}
	order := strings.Join(parts, ", ")
	w6aRI(f, "count", n)
	f["order"] = order
	switch kp {
	case 1:
		i := rand.Intn(n)
		w6aRI(f, "target", vals[i])
		ans := i + 1
		return w6aR(kp, fmt.Sprintf("A finite set is well ordered as \\(v_1\\prec v_2\\prec\\cdots\\prec v_{%d}\\), with values \\((%s)\\). What is the 1-indexed rank of the value %d? (enter a number)", n, order, vals[i]), strconv.Itoa(ans),
			fmt.Sprintf("The values are listed in well-order, so %d is the %d-th value and its rank is %d.", vals[i], ans, ans), f)
	case 2:
		imax := 0
		for i, v := range vals {
			if v > vals[imax] {
				imax = i
			}
		}
		w6aRI(f, "largest", vals[imax])
		ans := imax + 1
		return w6aR(kp, fmt.Sprintf("A finite set is well ordered as \\(v_1\\prec\\cdots\\prec v_{%d}\\), with values \\((%s)\\). What is the 1-indexed rank of the greatest value? (enter a number)", n, order), strconv.Itoa(ans),
			fmt.Sprintf("The greatest value is %d, and it is the %d-th value in the listing, so its rank is %d.", vals[imax], ans, ans), f)
	default:
		i := rand.Intn(n - 1)
		j := i + 1 + rand.Intn(n-i-1)
		w6aRI(f, "first", vals[i])
		w6aRI(f, "second", vals[j])
		ans := j - i - 1
		return w6aR(kp, fmt.Sprintf("A finite set is well ordered as \\(v_1\\prec\\cdots\\prec v_{%d}\\), with values \\((%s)\\). How many elements lie strictly between the values %d and %d in this well-order? (enter a number)", n, order, vals[i], vals[j]), strconv.Itoa(ans),
			fmt.Sprintf("The values %d and %d occupy ranks %d and %d, so \\(%d-%d-1=%d\\) values lie between them.", vals[i], vals[j], i+1, j+1, j+1, i+1, ans), f)
	}
}

func w6aCardinal(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		n := w6aRR(1, 10)
		w6aRI(f, "n", n)
		ans := 1 << n
		return w6aR(kp, fmt.Sprintf("A set \\(A\\) has %d elements. How many elements does its power set have? (enter a number)", n), strconv.Itoa(ans),
			fmt.Sprintf("Each element is either included or excluded, so the power set has \\(2^{%d}=%d\\) elements.", n, ans), f)
	case 2:
		a, b := w6aRR(2, 9), w6aRR(2, 9)
		w6aRI(f, "a", a)
		w6aRI(f, "b", b)
		ans := a * b
		return w6aR(kp, fmt.Sprintf("Finite sets \\(A\\) and \\(B\\) have %d and %d elements. What is the cardinality of \\(A\\times B\\)? (enter a number)", a, b), strconv.Itoa(ans),
			fmt.Sprintf("Every element of \\(A\\) pairs with every element of \\(B\\): \\(%d\\cdot%d=%d\\).", a, b, ans), f)
	default:
		a, b := w6aRR(3, 9), w6aRR(3, 9)
		overlap := w6aRR(1, w6aMin(a, b)-1)
		w6aRI(f, "a", a)
		w6aRI(f, "b", b)
		w6aRI(f, "overlap", overlap)
		ans := a + b - overlap
		return w6aR(kp, fmt.Sprintf("Finite sets \\(A\\) and \\(B\\) have %d and %d elements, and their intersection has %d elements. What is the cardinality of \\(A\\cup B\\)? (enter a number)", a, b, overlap), strconv.Itoa(ans),
			fmt.Sprintf("Inclusion-exclusion gives \\(%d+%d-%d=%d\\).", a, b, overlap, ans), f)
	}
}

// w6aRootTest: the limiting root of a geometric term and the radius of a
// geometric power series. The square m^2 is written unsimplified in the stem so
// the answer is a computation rather than a copy.
func w6aRootTest(kp int) generator.Problem {
	f := map[string]string{}
	m := w6aRR(2, 12)
	sq := m * m
	w6aRI(f, "base", m)
	w6aRI(f, "limit", sq)
	switch kp {
	case 1:
		return w6aR(kp, fmt.Sprintf("For the sequence \\(a_n=(%d^2)^n\\), what is \\(\\lim |a_n|^{1/n}\\)? (enter a number)", m), strconv.Itoa(sq),
			fmt.Sprintf("\\(|a_n|^{1/n}=%d^2=%d\\), so the limiting root is %d.", m, sq, sq), f)
	case 2:
		return w6aR(kp, fmt.Sprintf("For the power series \\(\\sum_{n\\ge0}(x/(%d^2))^n\\), what is the radius of convergence? (enter a number)", m), strconv.Itoa(sq),
			fmt.Sprintf("The series is geometric with ratio \\(x/%d\\), so it converges exactly for \\(|x|<%d\\). The radius of convergence is %d.", sq, sq, sq), f)
	default:
		return w6aR(kp, fmt.Sprintf("For the sequence \\(a_n=(%d^2)^n\\), what is \\(\\lim |a_{n+1}/a_n|\\)? (enter a number)", m), strconv.Itoa(sq),
			fmt.Sprintf("The ratio is \\(|a_{n+1}/a_n|=%d^2=%d\\), so the ratio test limit is %d.", m, sq, sq), f)
	}
}

// w6aIntegralTest: integral-test denominators for p-series tails. Answers are
// denominators or thresholds, never the integral value itself, so every answer
// is an integer.
func w6aIntegralTest(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		n := w6aRR(2, 9)
		w6aRI(f, "n", n)
		ans := 2 * n * n
		return w6aR(kp, fmt.Sprintf("The integral test bounds the tail \\(\\sum_{n>N}1/n^3\\) by \\(\\int_N^\\infty x^{-3}\\,dx\\). For \\(N=%d\\) this bound is \\(1/k\\). What is \\(k\\)? (enter a number)", n), strconv.Itoa(ans),
			fmt.Sprintf("\\(\\int_N^\\infty x^{-3}\\,dx=\\frac{1}{2N^2}\\), so for \\(N=%d\\) the denominator is \\(2\\cdot%d^2=%d\\).", n, n, ans), f)
	case 2:
		p := w6aRR(3, 14)
		w6aRI(f, "p", p)
		ans := p - 1
		return w6aR(kp, fmt.Sprintf("For \\(p=%d\\), the integral \\(\\int_1^\\infty x^{-p}\\,dx\\) equals \\(1/k\\). What is \\(k\\)? (enter a number)", p), strconv.Itoa(ans),
			fmt.Sprintf("\\(\\int_1^\\infty x^{-p}\\,dx=\\frac{1}{p-1}=\\frac{1}{%d}\\), so \\(k=%d\\).", ans, ans), f)
	default:
		l := w6aRR(2, 20)
		w6aRI(f, "bound_denominator", l)
		ans := l + 1
		return w6aR(kp, fmt.Sprintf("For the p-series with \\(p=2\\), the tail bound is \\(\\int_N^\\infty x^{-2}\\,dx=1/N\\). What is the smallest integer \\(N\\) for which this bound is less than \\(1/%d\\)? (enter a number)", l), strconv.Itoa(ans),
			fmt.Sprintf("\\(1/N<1/%d\\) exactly when \\(N>%d\\), so the smallest integer is %d.", l, l, ans), f)
	}
}

func w6aCesaro(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		c := w6aRR(10, 20)
		n := w6aRR(3, 9)
		w6aRI(f, "value", c)
		w6aRI(f, "terms", n)
		ans := c * n
		return w6aR(kp, fmt.Sprintf("For the constant sequence \\(a_n=%d\\), the Cesaro mean of the first \\(N\\) terms is \\(\\sigma_N=(a_1+\\cdots+a_N)/N\\). What is \\(N\\sigma_N\\) for \\(N=%d\\)? (enter a number)", c, n), strconv.Itoa(ans),
			fmt.Sprintf("A constant sequence has \\(\\sigma_N=%d\\), so \\(N\\sigma_N=%d\\cdot%d=%d\\).", c, n, c, ans), f)
	case 2:
		a := w6aRR(10, 20)
		d := w6aRR(2, 9)
		n := 3 + 2*rand.Intn(4)
		w6aRI(f, "a", a)
		w6aRI(f, "d", d)
		w6aRI(f, "terms", n)
		ans := a + d*(n-1)/2
		return w6aR(kp, fmt.Sprintf("For the arithmetic progression \\(a_n=%d+(n-1)%d\\), what is the Cesaro mean of the first \\(N\\) terms for \\(N=%d\\)? (enter a number)", a, d, n), strconv.Itoa(ans),
			fmt.Sprintf("The mean of an arithmetic progression is its midpoint value: \\(%d+%d\\cdot\\frac{%d-1}{2}=%d\\).", a, d, n, ans), f)
	default:
		a := w6aRR(10, 20)
		d := w6aRR(2, 9)
		n := w6aRR(3, 9)
		w6aRI(f, "a", a)
		w6aRI(f, "d", d)
		w6aRI(f, "terms", n)
		ans := n*a + d*n*(n-1)/2
		return w6aR(kp, fmt.Sprintf("For the arithmetic progression \\(a_n=%d+(n-1)%d\\), what is \\(N\\) times the Cesaro mean of the first \\(N\\) terms for \\(N=%d\\)? (enter a number)", a, d, n), strconv.Itoa(ans),
			fmt.Sprintf("The mean is \\(%d+%d\\cdot\\frac{%d-1}{2}=%d\\). Multiplying by \\(N=%d\\) gives \\(%d\\cdot%d=%d\\).", a, d, n, a+d*(n-1)/2, n, n, a+d*(n-1)/2, ans), f)
	}
}

// w6aDirichlet: bounds from the Dirichlet and alternating-series tests. The
// constants are chosen so the answer is a distance or a denominator, not a
// printed operand.
func w6aDirichlet(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		m := w6aRR(2, 9)
		k := w6aRR(2, 9)
		w6aRI(f, "partial_bound", m)
		w6aRI(f, "first_term", k)
		ans := m * k
		return w6aR(kp, fmt.Sprintf("The partial sums of \\((a_n)\\) satisfy \\(|\\sum_{n\\le N}a_n|\\le %d\\) for every \\(N\\), and \\((b_n)\\) decreases to 0 with \\(b_1=%d\\). Dirichlet's test bounds \\(|\\sum a_n b_n|\\) by \\(M b_1\\). What is this bound? (enter a number)", m, k), strconv.Itoa(ans),
			fmt.Sprintf("The bound is \\(M b_1=%d\\cdot%d=%d\\).", m, k, ans), f)
	case 2:
		n := w6aRR(2, 9)
		w6aRI(f, "terms", n)
		ans := (n + 1) * (n + 1)
		return w6aR(kp, fmt.Sprintf("For the alternating series \\(\\sum(-1)^{n+1}/n^2\\), the error after \\(N\\) terms is at most \\(1/(N+1)^2\\). For \\(N=%d\\), write this bound as \\(1/k\\). What is \\(k\\)? (enter a number)", n), strconv.Itoa(ans),
			fmt.Sprintf("The first omitted term is \\(1/(%d+1)^2=1/%d\\), so \\(k=%d\\).", n, ans, ans), f)
	default:
		n := 4 + 2*rand.Intn(5)
		b := w6aRR(50, 90)
		d := w6aRR(3, 6)
		w6aRI(f, "terms", n)
		w6aRI(f, "constant", b)
		w6aRI(f, "difference", d)
		ans := (n / 2) * d
		return w6aR(kp, fmt.Sprintf("For \\(a_n=(-1)^{n+1}\\) and \\(b_n=%d-(n-1)%d\\), what is \\(\\sum_{n=1}^{N}a_n b_n\\) for \\(N=%d\\)? (enter a number)", b, d, n), strconv.Itoa(ans),
			fmt.Sprintf("Group consecutive terms: \\((b_1-b_2)+(b_3-b_4)+\\cdots\\). Each pair cancels the constant and leaves the difference %d, and there are \\(%d/2=%d\\) pairs, so the sum is \\(%d\\cdot%d=%d\\).", d, n, n/2, n/2, d, ans), f)
	}
}

// w6aOscillation: the oscillation is the supremum minus the infimum of the
// values a function takes on the interval.
func w6aOscillation(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		a := w6aRR(10, 15)
		b := w6aRR(40, 50)
		w6aRI(f, "value_left", a)
		w6aRI(f, "value_right", b)
		ans := b - a
		return w6aR(kp, fmt.Sprintf("Let \\(f(x)=%d\\) for \\(x<1\\) and \\(f(x)=%d\\) for \\(x\\ge1\\). What is the oscillation of \\(f\\) on \\([0,2]\\)? (enter a number)", a, b), strconv.Itoa(ans),
			fmt.Sprintf("On \\([0,2]\\) the values are %d and %d, so the oscillation is \\(%d-%d=%d\\).", a, b, b, a, ans), f)
	case 2:
		a := w6aRR(1, 9)
		b := w6aRR(10, 19)
		c := w6aRR(30, 39)
		w6aRI(f, "value_low", a)
		w6aRI(f, "value_mid", b)
		w6aRI(f, "value_high", c)
		ans := c - a
		return w6aR(kp, fmt.Sprintf("Let \\(f(x)=%d\\) on \\([0,1)\\), \\(f(x)=%d\\) on \\([1,2)\\) and \\(f(x)=%d\\) on \\([2,3]\\). What is the oscillation of \\(f\\) on \\([0,3]\\)? (enter a number)", a, b, c), strconv.Itoa(ans),
			fmt.Sprintf("The largest value is %d and the smallest is %d, so the oscillation is \\(%d-%d=%d\\).", c, a, c, a, ans), f)
	default:
		a := w6aRR(2, 5)
		b := w6aRR(10, 15)
		w6aRI(f, "a", a)
		w6aRI(f, "b", b)
		ans := b*b - a*a
		return w6aR(kp, fmt.Sprintf("For \\(f(x)=x^2\\) on \\([%d,%d]\\), what is the oscillation of \\(f\\) on that interval? (enter a number)", a, b), strconv.Itoa(ans),
			fmt.Sprintf("\\(f\\) increases on the interval, so the supremum is \\(%d^2=%d\\) and the infimum is \\(%d^2=%d\\). The oscillation is \\(%d-%d=%d\\).", b, b*b, a, a*a, b*b, a*a, ans), f)
	}
}

// w6aSemicontinuity: a jump at c. Lower semicontinuity caps f(c) by the
// smaller one-sided limit, upper semicontinuity floors it by the larger one.
func w6aSemicontinuity(kp int) generator.Problem {
	f := map[string]string{}
	m := w6aRNZ(-3, 3)
	c := w6aRR(1, 6)
	b := w6aRR(20, 30)
	d := w6aRR(40, 50)
	left := m*c + b
	right := m*c + d
	w6aRI(f, "slope", m)
	w6aRI(f, "point", c)
	w6aRI(f, "left_offset", b)
	w6aRI(f, "right_offset", d)
	w6aRI(f, "left_limit", left)
	w6aRI(f, "right_limit", right)
	switch kp {
	case 1:
		ans := w6aMin(left, right)
		return w6aR(kp, fmt.Sprintf("A function equals \\(%dx+%d\\) for \\(x<%d\\) and \\(%dx+%d\\) for \\(x>%d\\). It is lower semicontinuous at %d. What is the largest permitted value of \\(f(%d)\\)? (enter a number)", m, b, c, m, d, c, c, c), strconv.Itoa(ans),
			fmt.Sprintf("The left limit is \\(%d\\cdot%d+%d=%d\\) and the right limit is \\(%d\\cdot%d+%d=%d\\). Lower semicontinuity needs \\(f(%d)\\le\\min(%d,%d)=%d\\), so the largest permitted value is %d.", m, c, b, left, m, c, d, right, c, left, right, ans, ans), f)
	case 2:
		ans := w6aMax(left, right)
		return w6aR(kp, fmt.Sprintf("A function equals \\(%dx+%d\\) for \\(x<%d\\) and \\(%dx+%d\\) for \\(x>%d\\). It is upper semicontinuous at %d. What is the smallest permitted value of \\(f(%d)\\)? (enter a number)", m, b, c, m, d, c, c, c), strconv.Itoa(ans),
			fmt.Sprintf("The left limit is \\(%d\\cdot%d+%d=%d\\) and the right limit is \\(%d\\cdot%d+%d=%d\\). Upper semicontinuity needs \\(f(%d)\\ge\\max(%d,%d)=%d\\), so the smallest permitted value is %d.", m, c, b, left, m, c, d, right, c, left, right, ans, ans), f)
	default:
		ans := d - b
		return w6aR(kp, fmt.Sprintf("A function equals \\(%dx+%d\\) for \\(x<%d\\) and \\(%dx+%d\\) for \\(x>%d\\), and is upper semicontinuous at %d. At the smallest permitted \\(f(%d)\\), what is \\(f(%d)-f(%d^-)\\), where the second term is the left limit? (enter a number)", m, b, c, m, d, c, c, c, c, c), strconv.Itoa(ans),
			fmt.Sprintf("Upper semicontinuity gives \\(f(%d)\\ge\\max(%d,%d)=%d\\). The left limit is %d, so the jump is \\(%d-%d=%d\\).", c, left, right, right, left, right, left, ans), f)
	}
}

func w6aDini(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		n := w6aRR(2, 6)
		k := w6aRR(10, 19)
		a := k * n
		w6aRI(f, "index", n)
		w6aRI(f, "endpoint", a)
		return w6aR(kp, fmt.Sprintf("For \\(f_n(x)=x/%d\\) on \\([0,%d]\\), what is the sup norm \\(\\|f_n\\|_\\infty\\)? (enter a number)", n, a), strconv.Itoa(k),
			fmt.Sprintf("\\(f_n\\) increases on the interval, so its maximum is at \\(x=%d\\): \\(%d/%d=%d\\).", a, a, n, k), f)
	case 2:
		n := w6aRR(2, 6)
		k := w6aRR(10, 19)
		a := k * n * n
		w6aRI(f, "index", n)
		w6aRI(f, "endpoint", a)
		return w6aR(kp, fmt.Sprintf("For \\(f_n(x)=x/(%d^2)\\) on \\([0,%d]\\), what is the sup norm \\(\\|f_n\\|_\\infty\\)? (enter a number)", n, a), strconv.Itoa(k),
			fmt.Sprintf("\\(f_n\\) increases and reaches \\(%d/%d^2=%d\\) at the endpoint.", a, n, k), f)
	default:
		a := w6aRR(3, 20)
		w6aRI(f, "endpoint", a)
		ans := a + 1
		return w6aR(kp, fmt.Sprintf("For \\(f_n(x)=x/n\\) on \\([0,%d]\\), what is the smallest \\(N\\) with \\(\\|f_N\\|_\\infty<1\\)? (enter a number)", a), strconv.Itoa(ans),
			fmt.Sprintf("The sup norm is \\(%d/N\\), and \\(%d/N<1\\) exactly when \\(N>%d\\). The smallest integer is %d.", a, a, a, ans), f)
	}
}

func w6aUniformLimit(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		c := w6aRR(20, 30)
		n := w6aRR(2, 6)
		k := w6aRR(10, 19)
		a := k * n
		w6aRI(f, "level", c)
		w6aRI(f, "index", n)
		w6aRI(f, "endpoint", a)
		return w6aR(kp, fmt.Sprintf("For \\(f_n(x)=%d+x/%d\\) on \\([0,%d]\\), the pointwise limit is the constant %d. What is \\(\\|f_n-%d\\|_\\infty\\)? (enter a number)", c, n, a, c, c), strconv.Itoa(k),
			fmt.Sprintf("The difference is \\(x/%d\\), which is largest at \\(x=%d\\): \\(%d/%d=%d\\).", n, a, a, n, k), f)
	case 2:
		c := w6aRR(20, 30)
		n := w6aRR(2, 6)
		k := w6aRR(10, 19)
		a := k * n
		w6aRI(f, "level", c)
		w6aRI(f, "index", n)
		w6aRI(f, "endpoint", a)
		return w6aR(kp, fmt.Sprintf("For \\(f_n(x)=%d-x/%d\\) on \\([0,%d]\\), the pointwise limit is the constant %d. What is \\(\\|f_n-%d\\|_\\infty\\)? (enter a number)", c, n, a, c, c), strconv.Itoa(k),
			fmt.Sprintf("The difference is \\(-x/%d\\), whose size is largest at \\(x=%d\\): \\(%d/%d=%d\\).", n, a, a, n, k), f)
	default:
		a := w6aRR(3, 20)
		w6aRI(f, "endpoint", a)
		ans := 2 * a
		return w6aR(kp, fmt.Sprintf("For \\(f_n(x)=x/n\\) on \\([0,%d]\\), what is the smallest \\(N\\) with \\(\\|f_N\\|_\\infty\\le\\frac{1}{2}\\)? (enter a number)", a), strconv.Itoa(ans),
			fmt.Sprintf("The sup norm is \\(%d/N\\), and \\(%d/N\\le\\frac{1}{2}\\) exactly when \\(N\\ge%d\\). The smallest integer is %d.", a, a, ans, ans), f)
	}
}

func w6aRolle(kp int) generator.Problem {
	f := map[string]string{}
	a := w6aRR(-6, 6)
	b := a + 2*w6aRR(1, 4)
	c := (a + b) / 2
	w6aRI(f, "root1", a)
	w6aRI(f, "root2", b)
	switch kp {
	case 1:
		return w6aR(kp, fmt.Sprintf("The quadratic \\(f(x)=(x-%d)(x-%d)\\) has roots %d and %d. Rolle's theorem gives one interior critical point. What is it? (enter a number)", a, b, a, b), strconv.Itoa(c),
			fmt.Sprintf("A quadratic is symmetric about the midpoint of its roots: \\(x=\\frac{%d+%d}{2}=%d\\), and %d lies strictly between the roots.", a, b, c, c), f)
	case 2:
		k := w6aRR(2, 5)
		w6aRI(f, "scale", k)
		return w6aR(kp, fmt.Sprintf("The quadratic \\(f(x)=%d(x-%d)(x-%d)\\) has roots %d and %d. Where is its interior critical point? (enter a number)", k, a, b, a, b), strconv.Itoa(c),
			fmt.Sprintf("The factor %d rescales the values but not the location of the critical point: \\(x=\\frac{%d+%d}{2}=%d\\).", k, a, b, c), f)
	default:
		half := (b - a) / 2
		ans := -half * half
		return w6aR(kp, fmt.Sprintf("The quadratic \\(f(x)=(x-%d)(x-%d)\\) has roots %d and %d. What is the value of \\(f\\) at its interior critical point? (enter a number)", a, b, a, b), strconv.Itoa(ans),
			fmt.Sprintf("The critical point is \\(\\frac{%d+%d}{2}=%d\\). Then \\(f(%d)=(%d-%d)(%d-%d)=(%d)(%d)=%d\\).", a, b, c, c, c, a, c, b, c-a, c-b, ans), f)
	}
}

// w6aConvex: gaps for x^2. The gap at the midpoint is (a-b)^2/4, so equal
// parity gives an integer gap and four times the gap is always an integer.
func w6aConvex(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		a := w6aRR(-8, 8)
		b := a + 2*w6aRR(1, 5)
		c := (a + b) / 2
		half := (b - a) / 2
		ans := half * half
		w6aRI(f, "a", a)
		w6aRI(f, "b", b)
		return w6aR(kp, fmt.Sprintf("For \\(f(x)=x^2\\), the Jensen gap at the midpoint \\(m\\) of %d and %d is \\(\\frac{f(%d)+f(%d)}{2}-f(m)\\). What is this gap? (enter a number)", a, b, a, b), strconv.Itoa(ans),
			fmt.Sprintf("\\(m=%d\\), so the gap is \\(\\frac{%d+%d}{2}-%d=%d\\), which also equals \\((\\frac{%d-%d}{2})^2=%d\\).", c, a*a, b*b, c*c, ans, a, b, ans), f)
	case 2:
		a := w6aRR(-8, 8)
		b := a + w6aRR(1, 7)
		ans := (a - b) * (a - b)
		w6aRI(f, "a", a)
		w6aRI(f, "b", b)
		return w6aR(kp, fmt.Sprintf("For \\(f(x)=x^2\\), the Jensen gap at the midpoint of %d and %d is \\(\\frac{(a-b)^2}{4}\\). What is four times the gap? (enter a number)", a, b), strconv.Itoa(ans),
			fmt.Sprintf("The gap is \\(\\frac{(%d-%d)^2}{4}\\), so four times it is \\((%d-%d)^2=%d\\).", a, b, a, b, ans), f)
	default:
		a := w6aRR(-8, 8)
		b := a + 2*w6aRR(1, 5)
		ans := (a*a + b*b) / 2
		w6aRI(f, "a", a)
		w6aRI(f, "b", b)
		return w6aR(kp, fmt.Sprintf("For \\(f(x)=x^2\\), what is the average \\(\\frac{f(%d)+f(%d)}{2}\\)? (enter a number)", a, b), strconv.Itoa(ans),
			fmt.Sprintf("The average is \\(\\frac{%d+%d}{2}=\\frac{%d}{2}=%d\\).", a*a, b*b, a*a+b*b, ans), f)
	}
}

// w6aJensen: the weighted Jensen gap of x^2 is pq(a-b)^2/d^2 for weights p/d
// and q/d. Each ordinal uses a different weight pair, and a-b is a multiple of
// d, so the gap is an integer.
func w6aJensen(kp int) generator.Problem {
	f := map[string]string{}
	p, q, d := 1, 1, 2
	switch kp {
	case 2:
		p, q, d = 1, 2, 3
	case 3:
		p, q, d = 2, 3, 5
	}
	k := w6aRR(1, 9)
	a := w6aRR(-5, 5)
	b := a + d*k
	w6aRI(f, "weight_a", p)
	w6aRI(f, "weight_b", q)
	w6aRI(f, "denominator", d)
	w6aRI(f, "a", a)
	w6aRI(f, "b", b)
	ans := p * q * k * k
	return w6aR(kp, fmt.Sprintf("For \\(f(x)=x^2\\) and weights \\(\\frac{%d}{%d}\\) and \\(\\frac{%d}{%d}\\), the Jensen gap at %d and %d is \\(\\frac{%d}{%d}f(%d)+\\frac{%d}{%d}f(%d)-f(\\frac{%d}{%d}%d+\\frac{%d}{%d}%d)\\). What is this gap? (enter a number)", p, d, q, d, a, b, p, d, a, q, d, b, p, d, a, q, d, b), strconv.Itoa(ans),
		fmt.Sprintf("For \\(f(x)=x^2\\) the gap is \\(\\frac{pq}{d^2}(a-b)^2\\). Here \\(a-b=%d\\), so the gap is \\(\\frac{%d\\cdot%d}{%d^2}\\cdot%d^2=%d\\).", a-b, p, q, d, k, ans), f)
}

func w6aLipschitz(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		sign := 1
		if rand.Intn(2) == 0 {
			sign = -1
		}
		s := sign * w6aRR(3, 9)
		p0 := w6aRR(20, 30)
		p2 := p0 + 2*s
		w6aRI(f, "f0", p0)
		w6aRI(f, "f2", p2)
		ans := w6aAbs(s)
		return w6aR(kp, fmt.Sprintf("A linear function satisfies \\(f(0)=%d\\) and \\(f(2)=%d\\). What is its best Lipschitz constant on \\(\\mathbb{R}\\)? (enter a number)", p0, p2), strconv.Itoa(ans),
			fmt.Sprintf("The slope is \\(\\frac{%d-%d}{2}=%d\\), so \\(|f(x)-f(y)|=%d|x-y|\\) and the best constant is %d.", p2, p0, s, ans, ans), f)
	case 2:
		a := w6aRR(2, 9)
		w6aRI(f, "endpoint", a)
		ans := 2 * a
		return w6aR(kp, fmt.Sprintf("What is the best Lipschitz constant of \\(f(x)=x^2\\) on \\([0,%d]\\)? (enter a number)", a), strconv.Itoa(ans),
			fmt.Sprintf("\\(|f'(x)|=2x\\le 2\\cdot%d=%d\\) on the interval, and the bound is attained at \\(x=%d\\).", a, ans, a), f)
	default:
		a := w6aRR(1, 6)
		b := w6aRR(8, 20)
		w6aRI(f, "a", a)
		w6aRI(f, "b", b)
		ans := 2 * b
		return w6aR(kp, fmt.Sprintf("What is the best Lipschitz constant of \\(f(x)=x^2\\) on \\([-%d,%d]\\)? (enter a number)", a, b), strconv.Itoa(ans),
			fmt.Sprintf("\\(|f'(x)|=2|x|\\), and the largest value on \\([-%d,%d]\\) is at \\(x=%d\\): \\(2\\cdot%d=%d\\).", a, b, b, b, ans), f)
	}
}

// w6aInverseFunction: g'(y0)=1/f'(x0). The questions ask for the denominator k,
// so every answer is f'(x0), an integer by construction.
func w6aInverseFunction(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		s := w6aRNZ(-9, 9)
		for s > -3 && s < 3 {
			s = w6aRNZ(-9, 9)
		}
		p0 := w6aRR(20, 30)
		p2 := p0 + 2*s
		w6aRI(f, "f0", p0)
		w6aRI(f, "f2", p2)
		return w6aR(kp, fmt.Sprintf("A linear function satisfies \\(f(0)=%d\\) and \\(f(2)=%d\\). Let \\(g=f^{-1}\\). At \\(y_0=f(0)\\), \\(g'(y_0)=1/k\\). What is \\(k\\)? (enter a number)", p0, p2), strconv.Itoa(s),
			fmt.Sprintf("The slope is \\(\\frac{%d-%d}{2}=%d\\), so \\(f'(0)=%d\\) and \\(g'(y_0)=\\frac{1}{%d}\\). Hence \\(k=%d\\).", p2, p0, s, s, s, s), f)
	case 2:
		c := w6aRR(20, 30)
		x0 := w6aRR(2, 9)
		w6aRI(f, "c", c)
		w6aRI(f, "x0", x0)
		ans := 2 * x0
		return w6aR(kp, fmt.Sprintf("On \\(x\\ge0\\), \\(f(x)=x^2+%d\\). Let \\(g=f^{-1}\\). At \\(y_0=f(%d)\\), \\(g'(y_0)=1/k\\). What is \\(k\\)? (enter a number)", c, x0), strconv.Itoa(ans),
			fmt.Sprintf("\\(f'(x)=2x\\), so \\(f'(%d)=%d\\) and \\(g'(y_0)=\\frac{1}{%d}\\). Hence \\(k=%d\\).", x0, ans, ans, ans), f)
	default:
		a := w6aRR(1, 4)
		c := w6aRR(2, 9)
		x0 := w6aRR(2, 5)
		d := w6aRR(-5, 5)
		w6aRI(f, "a", a)
		w6aRI(f, "c", c)
		w6aRI(f, "x0", x0)
		w6aRI(f, "d", d)
		ans := 3*a*x0*x0 + c
		return w6aR(kp, fmt.Sprintf("Let \\(f(x)=%dx^3+%dx+%d\\) and \\(g=f^{-1}\\). At \\(y_0=f(%d)\\), \\(g'(y_0)=1/k\\). What is \\(k\\)? (enter a number)", a, c, d, x0), strconv.Itoa(ans),
			fmt.Sprintf("\\(f'(x)=3\\cdot%d\\,x^2+%d\\), so \\(f'(%d)=3\\cdot%d\\cdot%d^2+%d=%d\\) and \\(g'(y_0)=\\frac{1}{%d}\\). Hence \\(k=%d\\).", a, c, x0, a, x0, c, ans, ans, ans), f)
	}
}
