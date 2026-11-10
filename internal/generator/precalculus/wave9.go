package precalculus

// Wave 9 precalculus: the first three review concepts. Every family publishes
// the values its integer answer follows from, and wave9_test.go recomputes the
// answer from those facts.

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave9(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"precalc.func.one_to_one":                 "one_to_one",
		"precalc.poly.complex_factorization":      "complex_factorization",
		"precalc.vector.orthogonal_decomposition": "orthogonal_decomposition",
	} {
		reg.Register(id, &w9PreGen{kind: kind})
	}
}

type w9PreGen struct{ kind string }

func (g *w9PreGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "one_to_one":
		return w9OneToOne(kp)
	case "complex_factorization":
		return w9ComplexFactorization(kp)
	case "orthogonal_decomposition":
		return w9OrthogonalDecomposition(kp)
	}
	return generator.Problem{Question: "What is 2 + 2? (enter a number)", Answer: "4", Explanation: "2 + 2 = 4."}
}

func w9P(o int, q, a, ex string, facts map[string]string) generator.Problem {
	if facts == nil {
		facts = map[string]string{}
	}
	facts["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: facts, KPOrdinal: o}
}

func w9PI(m map[string]string, k string, v int) { m[k] = strconv.Itoa(v) }

func w9PR(lo, hi int) int {
	if hi <= lo {
		return lo
	}
	return lo + rand.Intn(hi-lo+1)
}

func w9PJoin(xs []int) string {
	parts := make([]string, len(xs))
	for i, x := range xs {
		parts[i] = strconv.Itoa(x)
	}
	return strings.Join(parts, ", ")
}

func w9OneToOne(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		kinds := []string{"quadratic", "cubic", "linear", "absolute"}
		kind := kinds[rand.Intn(len(kinds))]
		a := w9PR(1, 4)
		h := w9PR(-3, 3)
		k := w9PR(-4, 4)
		f["kind"] = kind
		w9PI(f, "a", a)
		w9PI(f, "h", h)
		w9PI(f, "k", k)
		var inner string
		if h == 0 {
			inner = "x"
		} else if h > 0 {
			inner = fmt.Sprintf("(x-%d)", h)
		} else {
			inner = fmt.Sprintf("(x+%d)", -h)
		}
		coef := fmt.Sprintf("%d", a)
		if a == 1 {
			coef = ""
		}
		constText := ""
		if k > 0 {
			constText = fmt.Sprintf("+%d", k)
		} else if k < 0 {
			constText = fmt.Sprintf("-%d", -k)
		}
		var expr string
		ans := 0
		switch kind {
		case "quadratic":
			expr = fmt.Sprintf("y=%s%s^2%s", coef, inner, constText)
		case "cubic":
			expr = fmt.Sprintf("y=%s%s^3%s", coef, inner, constText)
			ans = 1
		case "linear":
			expr = fmt.Sprintf("y=%s%s%s", coef, inner, constText)
			ans = 1
		default:
			absInner := "x"
			if h > 0 {
				absInner = fmt.Sprintf("x-%d", h)
			} else if h < 0 {
				absInner = fmt.Sprintf("x+%d", -h)
			}
			expr = fmt.Sprintf("y=%s|%s|%s", coef, absInner, constText)
		}
		return w9P(kp, fmt.Sprintf("Does the graph of %s pass the horizontal-line test on all real numbers? Enter 1 for yes and 0 for no.", expr), strconv.Itoa(ans),
			"The horizontal-line test fails when two inputs share one output. A quadratic and an absolute-value graph reflect across a vertical line, so they fail; a cubic and a nonzero linear function pass.", f)
	case 2:
		ys := []int{w9PR(1, 6), w9PR(1, 6), w9PR(1, 6), w9PR(1, 6)}
		seen := map[int]bool{}
		distinct := 0
		for _, y := range ys {
			if !seen[y] {
				seen[y] = true
				distinct++
			}
		}
		w9PI(f, "distinct", distinct)
		f["ys"] = w9PJoin(ys)
		ans := 0
		if distinct == len(ys) {
			ans = 1
		}
		return w9P(kp, fmt.Sprintf("A table of four input-output pairs has y-values %s. Is the relation one-to-one? Enter 1 for yes and 0 for no.", w9PJoin(ys)), strconv.Itoa(ans),
			fmt.Sprintf("The y-values have %d distinct entries among four pairs, so the relation is %s.", distinct, map[int]string{0: "not one-to-one", 1: "one-to-one"}[ans]), f)
	default:
		lo := w9PR(-5, -1)
		hi := w9PR(1, 5)
		if rand.Intn(2) == 0 {
			lo = w9PR(0, 3)
			hi = lo + w9PR(1, 5)
		}
		w9PI(f, "lo", lo)
		w9PI(f, "hi", hi)
		ans := 0
		if lo >= 0 || hi <= 0 {
			ans = 1
		}
		return w9P(kp, fmt.Sprintf("Is f(x)=x^2 one-to-one on the interval [%d,%d]? Enter 1 for yes and 0 for no.", lo, hi), strconv.Itoa(ans),
			"On an interval that lies entirely on one side of 0, x squared is strictly monotone. An interval that contains both a negative and a positive input has two inputs with the same square.", f)
	}
}

func w9ComplexFactorization(kp int) generator.Problem {
	n := w9PR(1, 3)
	roots := make([]int, n)
	prod := 1
	sum := 0
	for i := range roots {
		roots[i] = w9PR(-3, 3)
		prod *= roots[i]
		sum += roots[i]
	}
	a := w9PR(0, 3)
	b := w9PR(1, 4)
	f := map[string]string{"roots": w9PJoin(roots)}
	w9PI(f, "n", n)
	w9PI(f, "a", a)
	w9PI(f, "b", b)
	w9PI(f, "sum", sum)
	w9PI(f, "product", prod)
	degree := n + 2
	switch kp {
	case 1:
		return w9P(kp, fmt.Sprintf("A monic polynomial has real roots %s and the complex pair %d±%di. What is the sum of all its roots? (enter a number)", w9PJoin(roots), a, b), strconv.Itoa(sum+2*a),
			fmt.Sprintf("Add the real roots and the real parts of the complex pair: %d + 2(%d) = %d.", sum, a, sum+2*a), f)
	case 2:
		return w9P(kp, fmt.Sprintf("A monic polynomial has real roots %s and the complex pair %d±%di. What is the product of all its roots? (enter a number)", w9PJoin(roots), a, b), strconv.Itoa(prod*(a*a+b*b)),
			fmt.Sprintf("The product of the complex pair is %d^2+%d^2=%d. Multiplying by the real roots gives %d.", a, b, a*a+b*b, prod*(a*a+b*b)), f)
	default:
		realWord := "roots"
		if n == 1 {
			realWord = "root"
		}
		return w9P(kp, fmt.Sprintf("A monic polynomial has real roots %s and the complex pair %d±%di. How many linear factors does it have over the complex numbers? (enter a number)", w9PJoin(roots), a, b), strconv.Itoa(degree),
			fmt.Sprintf("The polynomial has %d real %s and one complex pair, so its degree, and therefore its number of linear factors over C, is %d.", n, realWord, degree), f)
	}
}

func w9OrthogonalDecomposition(kp int) generator.Problem {
	a := w9PR(1, 4)
	b := w9PR(1, 4)
	if rand.Intn(2) == 0 {
		b = -b
	}
	k := w9PR(1, 4)
	m := w9PR(1, 4)
	if rand.Intn(2) == 0 {
		m = -m
	}
	vx := k*a - m*b
	vy := k*b + m*a
	f := map[string]string{}
	w9PI(f, "ax", a)
	w9PI(f, "ay", b)
	w9PI(f, "k", k)
	w9PI(f, "m", m)
	w9PI(f, "vx", vx)
	w9PI(f, "vy", vy)
	switch kp {
	case 1:
		parallel := fmt.Sprintf("%d(%d,%d)", k, a, b)
		if k == 1 {
			parallel = fmt.Sprintf("(%d,%d)", a, b)
		} else if k == -1 {
			parallel = fmt.Sprintf("-(%d,%d)", a, b)
		}
		return w9P(kp, fmt.Sprintf("Decompose v=(%d,%d) into a component parallel to u=(%d,%d) and a component perpendicular to u. What is the x-coordinate of the parallel component? (enter a number)", vx, vy, a, b), strconv.Itoa(k*a),
			fmt.Sprintf("The parallel coefficient is (v·u)/(u·u). Here v·u=%d and u·u=%d, so the coefficient is %d and the parallel component is %s. Its x-coordinate is %d.", k*(a*a+b*b), a*a+b*b, k, parallel, k*a), f)
	case 2:
		perp := fmt.Sprintf("%d(%d,%d)", m, -b, a)
		if m == 1 {
			perp = fmt.Sprintf("(%d,%d)", -b, a)
		} else if m == -1 {
			perp = fmt.Sprintf("-(%d,%d)", -b, a)
		}
		return w9P(kp, fmt.Sprintf("Decompose v=(%d,%d) into a component parallel to u=(%d,%d) and a component perpendicular to u. What is the x-coordinate of the perpendicular component? (enter a number)", vx, vy, a, b), strconv.Itoa(-m*b),
			fmt.Sprintf("Subtract the parallel component from v. The perpendicular component is %s, whose x-coordinate is %d.", perp, -m*b), f)
	default:
		par := fmt.Sprintf("%d(%d,%d)", k, a, b)
		if k == 1 {
			par = fmt.Sprintf("(%d,%d)", a, b)
		} else if k == -1 {
			par = fmt.Sprintf("-(%d,%d)", a, b)
		}
		return w9P(kp, fmt.Sprintf("For v=(%d,%d) and u=(%d,%d), what is v·u? (enter a number)", vx, vy, a, b), strconv.Itoa(k*(a*a+b*b)),
			fmt.Sprintf("Because v = %s + (%d)(%d,%d), the second vector is perpendicular to u, so v·u = %d(u·u) = %d.", par, m, -b, a, k, k*(a*a+b*b)), f)
	}
}
