package realanalysis

// Wave 5 real analysis: the first five review-order concepts on limits of
// sequences and finite Lebesgue-measure examples. Every family publishes the
// operands and the finite measure data its integer answer follows from.

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave5(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"ra.seq.limsup":              "limsup",
		"ra.seq.liminf":              "liminf",
		"ra.seq.cluster":             "cluster",
		"ra.measure.lebesgue":        "lebesgue",
		"ra.measure.simple_integral": "simple_integral",
	} {
		reg.Register(id, &w5RAGen{kind: kind})
	}
}

type w5RAGen struct{ kind string }

func (g *w5RAGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "limsup":
		return w5Limsup(kp)
	case "liminf":
		return w5Liminf(kp)
	case "cluster":
		return w5Cluster(kp)
	case "lebesgue":
		return w5Lebesgue(kp)
	case "simple_integral":
		return w5SimpleIntegral(kp)
	}
	return generator.Problem{Question: "What is 2 + 2? (enter a number)", Answer: "4", Explanation: "2 + 2 = 4."}
}

func w5R(o int, q, a, ex string, facts map[string]string) generator.Problem {
	if facts == nil {
		facts = map[string]string{}
	}
	facts["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: facts, KPOrdinal: o}
}

func w5RI(m map[string]string, k string, v int) { m[k] = strconv.Itoa(v) }

func w5RR(lo, hi int) int {
	if hi <= lo {
		return lo
	}
	return lo + rand.Intn(hi-lo+1)
}

func w5RNZ(lo, hi int) int {
	if hi <= lo {
		return lo
	}
	for {
		v := w5RR(lo, hi)
		if v != 0 {
			return v
		}
	}
}

func w5Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func w5Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func w5ConstSum(a, b int) string {
	if b == 0 {
		return strconv.Itoa(a)
	}
	if a == 0 {
		return strconv.Itoa(b)
	}
	if b < 0 {
		return fmt.Sprintf("%d - %d", a, -b)
	}
	return fmt.Sprintf("%d + %d", a, b)
}

func w5Sub(a, b int) string {
	if b < 0 {
		return fmt.Sprintf("%d-(%d)", a, b)
	}
	return fmt.Sprintf("%d-%d", a, b)
}

func w5Array(vals []int) string {
	parts := make([]string, len(vals))
	for i, v := range vals {
		parts[i] = strconv.Itoa(v)
	}
	return "(" + strings.Join(parts, ", ") + ")"
}

func w5Limsup(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		a, b := w5RR(-4, 6), w5RR(-4, 6)
		w5RI(f, "even", a)
		w5RI(f, "odd", b)
		ans := w5Max(a, b)
		return w5R(kp, fmt.Sprintf("The sequence is \\(a_n=%d\\) for even \\(n\\) and \\(a_n=%d\\) for odd \\(n\\). What is \\(\\limsup a_n\\)? (enter a number)", a, b), strconv.Itoa(ans),
			fmt.Sprintf("The even subsequence has limit %d and the odd subsequence has limit %d. The limsup is their maximum, %d.", a, b, ans), f)
	case 2:
		l := w5RNZ(-5, 5)
		c := w5RR(1, 4)
		w5RI(f, "limit", l)
		w5RI(f, "coefficient", c)
		return w5R(kp, fmt.Sprintf("For \\(a_n=%d+%d/n\\), what is \\(\\limsup a_n\\)? (enter a number)", l, c), strconv.Itoa(l),
			fmt.Sprintf("The sequence converges to %d, so its limsup is %d.", l, l), f)
	default:
		a, b, c := w5RR(-4, 6), w5RR(-4, 6), w5RNZ(-4, 4)
		w5RI(f, "even", a)
		w5RI(f, "odd", b)
		w5RI(f, "shift", c)
		ans := c + w5Max(a, b)
		return w5R(kp, fmt.Sprintf("The sequence is \\(a_n=%s\\) for even \\(n\\) and \\(a_n=%s\\) for odd \\(n\\). What is \\(\\limsup a_n\\)? (enter a number)", w5ConstSum(a, c), w5ConstSum(b, c)), strconv.Itoa(ans),
			fmt.Sprintf("The two constant subsequences have limits %d and %d. Adding the shift %d gives \\(%d+\\max(%d,%d)=%d\\).", a+c, b+c, c, c, a, b, ans), f)
	}
}

func w5Liminf(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		a, b := w5RR(-4, 6), w5RR(-4, 6)
		w5RI(f, "even", a)
		w5RI(f, "odd", b)
		ans := w5Min(a, b)
		return w5R(kp, fmt.Sprintf("The sequence is \\(a_n=%d\\) for even \\(n\\) and \\(a_n=%d\\) for odd \\(n\\). What is \\(\\liminf a_n\\)? (enter a number)", a, b), strconv.Itoa(ans),
			fmt.Sprintf("The even subsequence has limit %d and the odd subsequence has limit %d. The liminf is their minimum, %d.", a, b, ans), f)
	case 2:
		l := w5RNZ(-5, 5)
		c := w5RR(1, 4)
		w5RI(f, "limit", l)
		w5RI(f, "coefficient", c)
		return w5R(kp, fmt.Sprintf("For \\(a_n=%d-%d/n\\), what is \\(\\liminf a_n\\)? (enter a number)", l, c), strconv.Itoa(l),
			fmt.Sprintf("The sequence converges to %d, so its liminf is %d.", l, l), f)
	default:
		a, b, c := w5RR(-4, 6), w5RR(-4, 6), w5RNZ(-4, 4)
		w5RI(f, "even", a)
		w5RI(f, "odd", b)
		w5RI(f, "shift", c)
		ans := c + w5Min(a, b)
		return w5R(kp, fmt.Sprintf("The sequence is \\(a_n=%s\\) for even \\(n\\) and \\(a_n=%s\\) for odd \\(n\\). What is \\(\\liminf a_n\\)? (enter a number)", w5ConstSum(a, c), w5ConstSum(b, c)), strconv.Itoa(ans),
			fmt.Sprintf("The two constant subsequences have limits %d and %d. Adding the shift %d gives \\(%d+\\min(%d,%d)=%d\\).", a+c, b+c, c, c, a, b, ans), f)
	}
}

func w5Cluster(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		k := w5RR(2, 4)
		vals := rand.Perm(11)[:k]
		for i := range vals {
			vals[i] -= 5
		}
		w5RI(f, "count", k)
		f["values"] = strings.Trim(strings.ReplaceAll(fmt.Sprint(vals), " ", ","), "[]")
		return w5R(kp, fmt.Sprintf("The sequence cycles through the %d distinct values %s. How many cluster points does it have? (enter a number)", k, w5Array(vals)), strconv.Itoa(k),
			fmt.Sprintf("Each value repeats infinitely often, so each is a cluster point. There are %d.", k), f)
	case 2:
		a, b := w5RR(-4, 6), w5RR(-4, 6)
		for a == b {
			b = w5RR(-4, 6)
		}
		w5RI(f, "even", a)
		w5RI(f, "odd", b)
		ans := w5Max(a, b)
		return w5R(kp, fmt.Sprintf("For \\(a_{2n}=%d+1/n\\) and \\(a_{2n+1}=%d\\), what is the larger cluster point? (enter a number)", a, b), strconv.Itoa(ans),
			fmt.Sprintf("The even subsequence tends to %d and the odd subsequence is constantly %d. The larger cluster point is %d.", a, b, ans), f)
	default:
		c, d := w5RNZ(-4, 4), w5RR(1, 4)
		w5RI(f, "center", c)
		w5RI(f, "offset", d)
		ans := c - d
		return w5R(kp, fmt.Sprintf("For \\(a_n=%d+(-1)^n %d\\), what is the smaller cluster point? (enter a number)", c, d), strconv.Itoa(ans),
			fmt.Sprintf("Even terms are \\(%d+%d\\) and odd terms are \\(%d-%d\\). The smaller cluster point is %d.", c, d, c, d, ans), f)
	}
}

func w5Lebesgue(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		a := w5RR(-5, 2)
		b := a + w5RR(1, 6)
		w5RI(f, "a", a)
		w5RI(f, "b", b)
		ans := b - a
		return w5R(kp, fmt.Sprintf("What is the Lebesgue measure of the interval \\([%d,%d]\\)? (enter a number)", a, b), strconv.Itoa(ans),
			fmt.Sprintf("The measure of a bounded interval is its length: \\(%s=%d\\).", w5Sub(b, a), ans), f)
	case 2:
		a := w5RR(-5, 1)
		b := a + w5RR(1, 4)
		c := b + w5RR(1, 4)
		d := c + w5RR(1, 4)
		w5RI(f, "a", a)
		w5RI(f, "b", b)
		w5RI(f, "c", c)
		w5RI(f, "d", d)
		ans := (b - a) + (d - c)
		return w5R(kp, fmt.Sprintf("What is the Lebesgue measure of \\([%d,%d]\\cup[%d,%d]\\)? (enter a number)", a, b, c, d), strconv.Itoa(ans),
			fmt.Sprintf("The intervals are disjoint, so add their lengths: \\(%s+%s=%d\\).", w5Sub(b, a), w5Sub(d, c), ans), f)
	default:
		a := w5RR(-5, 0)
		c := a + w5RR(1, 3)
		b := c + w5RR(1, 3)
		d := b + w5RR(1, 3)
		w5RI(f, "a", a)
		w5RI(f, "b", b)
		w5RI(f, "c", c)
		w5RI(f, "d", d)
		ans := d - a
		return w5R(kp, fmt.Sprintf("What is the Lebesgue measure of \\([%d,%d]\\cup[%d,%d]\\) when the intervals overlap? (enter a number)", a, b, c, d), strconv.Itoa(ans),
			fmt.Sprintf("The union is \\([%d,%d]\\), so its measure is \\(%s=%d\\).", a, d, w5Sub(d, a), ans), f)
	}
}

func w5SimpleIntegral(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		a := w5RR(-4, 0)
		b := a + w5RR(1, 4)
		c := b + w5RR(1, 4)
		d := c + w5RR(1, 4)
		c1, c2 := w5RNZ(1, 5), w5RNZ(1, 5)
		w5RI(f, "a", a)
		w5RI(f, "b", b)
		w5RI(f, "c", c)
		w5RI(f, "d", d)
		w5RI(f, "c1", c1)
		w5RI(f, "c2", c2)
		ans := c1*(b-a) + c2*(d-c)
		return w5R(kp, fmt.Sprintf("A simple function equals %d on \\([%d,%d]\\) and %d on \\([%d,%d]\\), and is zero elsewhere. What is its integral? (enter a number)", c1, a, b, c2, c, d), strconv.Itoa(ans),
			fmt.Sprintf("Integrate each constant over its interval: \\(%d(%d)+%d(%d)=%d\\).", c1, b-a, c2, d-c, ans), f)
	case 2:
		c1, c2, c3 := w5RNZ(1, 5), w5RNZ(1, 5), w5RNZ(1, 5)
		w1, w2, w3 := w5RR(1, 4), w5RR(1, 4), w5RR(1, 4)
		w5RI(f, "c1", c1)
		w5RI(f, "c2", c2)
		w5RI(f, "c3", c3)
		w5RI(f, "w1", w1)
		w5RI(f, "w2", w2)
		w5RI(f, "w3", w3)
		ans := c1*w1 + c2*w2 + c3*w3
		return w5R(kp, fmt.Sprintf("A finite counting measure gives weights %d, %d and %d to three points. A simple function has values %d, %d and %d there. What is the integral? (enter a number)", w1, w2, w3, c1, c2, c3), strconv.Itoa(ans),
			fmt.Sprintf("A simple integral is the sum of value times measure: \\(%d\\cdot%d+%d\\cdot%d+%d\\cdot%d=%d\\).", c1, w1, c2, w2, c3, w3, ans), f)
	default:
		a := w5RR(-4, 0)
		b := a + w5RR(1, 4)
		c := b + w5RR(1, 4)
		d := c + w5RR(1, 4)
		w5RI(f, "a", a)
		w5RI(f, "b", b)
		w5RI(f, "c", c)
		w5RI(f, "d", d)
		ans := (b - a) + (d - c)
		return w5R(kp, fmt.Sprintf("What is the integral of the indicator of \\([%d,%d]\\cup[%d,%d]\\) with respect to Lebesgue measure? (enter a number)", a, b, c, d), strconv.Itoa(ans),
			fmt.Sprintf("The indicator integrates to the measure of its support: \\(%s+%s=%d\\).", w5Sub(b, a), w5Sub(d, c), ans), f)
	}
}
