package calculus

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

// Wave 21 adds six calculus concepts: the reciprocal rule, the Leibniz rule,
// the linear-approximation error, Stolz-Cesaro, recursive sequences and
// integrals of periodic functions.

func registerWave21(reg *generator.Registry) {
	reg.Register("calc.deriv.reciprocal", &w21Gen{kind: "reciprocal"})
	reg.Register("calc.deriv.leibniz_rule", &w21Gen{kind: "leibniz"})
	reg.Register("calc.deriv.linear_error", &w21Gen{kind: "linear_error"})
	reg.Register("calc.limit.stolz", &w21Gen{kind: "stolz"})
	reg.Register("calc.seq.recursive", &w21Gen{kind: "recursive"})
	reg.Register("calc.integral.periodic", &w21Gen{kind: "periodic"})
}

type w21Gen struct{ kind string }

func w21N(lo, hi int) int { return lo + rand.Intn(hi-lo+1) }

func w21P(o int, q, a, ex string, facts map[string]string) generator.Problem {
	facts["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: facts, KPOrdinal: o}
}

func w21I(f map[string]string, k string, v int) { f[k] = strconv.Itoa(v) }

func w21Binom(n, k int) int {
	if k < 0 || k > n {
		return 0
	}
	r := 1
	for i := 1; i <= k; i++ {
		r = r * (n - i + 1) / i
	}
	return r
}

func (g *w21Gen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "reciprocal":
		return w21Reciprocal(kp)
	case "leibniz":
		return w21Leibniz(kp)
	case "linear_error":
		return w21LinearError(kp)
	case "stolz":
		return w21Stolz(kp)
	case "recursive":
		return w21Recursive(kp)
	default:
		return w21Periodic(kp)
	}
}

func w21Reciprocal(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		n := w21N(1, 30)
		w21I(f, "n", n)
		return w21P(kp,
			fmt.Sprintf("What is the magnitude of the coefficient in the derivative of 1/x^%d? (enter a number)", n),
			strconv.Itoa(n),
			fmt.Sprintf("The derivative is -%d/x^%d, so the coefficient has magnitude %d.", n, n+1, n), f)
	case 2:
		a := w21N(1, 20)
		w21I(f, "a", a)
		return w21P(kp,
			fmt.Sprintf("For g(x)=x, the derivative of 1/g at x=%d is -1/a^2. What is a^2? (enter a number)", a),
			strconv.Itoa(a*a),
			fmt.Sprintf("At x=%d the squared denominator is %d.", a, a*a), f)
	default:
		q := w21N(1, 30)
		w21I(f, "q", q)
		return w21P(kp,
			fmt.Sprintf("At a point where u=1 and u'=%d, the derivative of 1/u is -q. What is q? (enter a number)", q),
			strconv.Itoa(q),
			fmt.Sprintf("The reciprocal rule gives -u'/u^2 = -%d, so q=%d.", q, q), f)
	}
}

func w21Leibniz(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		n := w21N(1, 20)
		w21I(f, "n", n)
		return w21P(kp,
			fmt.Sprintf("How many terms does the %d-th derivative of a product have? (enter a number)", n),
			strconv.Itoa(n+1),
			fmt.Sprintf("The orders split as 0 through %d, giving %d terms.", n, n+1), f)
	case 2:
		n := w21N(2, 15)
		k := w21N(1, n-1)
		w21I(f, "n", n)
		w21I(f, "k", k)
		return w21P(kp,
			fmt.Sprintf("In the %d-th derivative of a product, what is the coefficient of the term with %d derivatives on the second factor? (enter a number)", n, k),
			strconv.Itoa(w21Binom(n, k)),
			fmt.Sprintf("The coefficient is C(%d,%d), which is %d.", n, k, w21Binom(n, k)), f)
	default:
		d := w21N(1, 20)
		w21I(f, "d", d)
		return w21P(kp,
			fmt.Sprintf("One factor is a polynomial of degree %d. How many terms survive in a high-order Leibniz expansion? (enter a number)", d),
			strconv.Itoa(d+1),
			fmt.Sprintf("The polynomial contributes derivatives of order 0 through %d, so %d terms survive.", d, d+1), f)
	}
}

func w21LinearError(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		m := w21N(1, 15)
		d := w21N(1, 12)
		w21I(f, "m", m)
		w21I(f, "d", d)
		return w21P(kp,
			fmt.Sprintf("A function has second-derivative bound %d and the step is %d. The error bound is M d^2/2. What is the bound? (enter a number)", 2*m, d),
			strconv.Itoa(m*d*d),
			fmt.Sprintf("M d^2/2 = %d * %d^2 / 2 = %d.", 2*m, d, m*d*d), f)
	case 2:
		bigM := 2 * w21N(1, 10)
		d := w21N(1, 12)
		w21I(f, "m", bigM)
		w21I(f, "d", d)
		return w21P(kp,
			fmt.Sprintf("The second-derivative bound is %d and the step is %d. What is the linear-approximation error bound? (enter a number)", bigM, d),
			strconv.Itoa(bigM*d*d/2),
			fmt.Sprintf("%d * %d^2 / 2 = %d.", bigM, d, bigM*d*d/2), f)
	default:
		k := w21N(2, 12)
		w21I(f, "k", k)
		return w21P(kp,
			fmt.Sprintf("The step is multiplied by %d. By what factor does the linear-approximation error bound change? (enter a number)", k),
			strconv.Itoa(k*k),
			fmt.Sprintf("The error scales with the square of the step, so the factor is %d^2 = %d.", k, k*k), f)
	}
}

func w21Stolz(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		p := w21N(2, 20)
		w21I(f, "p", p)
		return w21P(kp,
			fmt.Sprintf("For a_n=n^p and b_n=n^(p+1), the difference ratio tends to p/(p+1). What is p+1? p is %d. (enter a number)", p),
			strconv.Itoa(p+1),
			fmt.Sprintf("The numerator exponent is %d, so the denominator exponent is %d.", p, p+1), f)
	case 2:
		k := w21N(1, 20)
		w21I(f, "k", k)
		return w21P(kp,
			fmt.Sprintf("For b_n=n^%d, the difference ratio tends to k+1. What is k+1? (enter a number)", k),
			strconv.Itoa(k+1),
			fmt.Sprintf("The exponent of the denominator is %d, so the limit is %d.", k, k+1), f)
	default:
		p := w21N(2, 9)
		q := p + w21N(1, 6)
		w21I(f, "p", p)
		w21I(f, "q", q)
		return w21P(kp,
			fmt.Sprintf("For a_n=%d^n and b_n=%d^n, the difference ratio tends to q/p. What is q? (enter a number)", q, p),
			strconv.Itoa(q),
			fmt.Sprintf("The numerator base is %d, so q=%d.", q, q), f)
	}
}

func w21Recursive(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		k := 1 << w21N(1, 10)
		w21I(f, "k", k)
		steps := 0
		for v := k; v >= 1; v /= 2 {
			steps++
		}
		// steps counts k, k/2, ... down to 1; halvings needed to drop below 1.
		return w21P(kp,
			fmt.Sprintf("The sequence starts at %d and halves each step. How many halvings are needed to drop below 1? (enter a number)", k),
			strconv.Itoa(steps),
			fmt.Sprintf("Halving %d repeatedly: %d halvings reach 1 and one more drops below it.", k, steps-1), f)
	case 2:
		a := w21N(1, 30)
		w21I(f, "a", a)
		return w21P(kp,
			fmt.Sprintf("The recursion is x(n+1)=(x(n)+%d)/2. What is the limit? (enter a number)", a),
			strconv.Itoa(a),
			fmt.Sprintf("The fixed point solves x=(x+%d)/2, so x=%d.", a, a), f)
	default:
		c := w21N(1, 30)
		w21I(f, "c", c)
		return w21P(kp,
			fmt.Sprintf("The recursion is x(n+1)=x(n)/2+%d. What is the limit? (enter a number)", c),
			strconv.Itoa(2*c),
			fmt.Sprintf("The fixed point is c/(1-1/2)=2c, which is %d.", 2*c), f)
	}
}

func w21Periodic(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		i := w21N(1, 30)
		k := w21N(2, 20)
		w21I(f, "i", i)
		w21I(f, "k", k)
		return w21P(kp,
			fmt.Sprintf("A periodic function has integral %d over one period. What is its integral over %d periods? (enter a number)", i, k),
			strconv.Itoa(k*i),
			fmt.Sprintf("Each period contributes %d, so %d periods contribute %d.", i, k, k*i), f)
	case 2:
		i := w21N(2, 40)
		j := w21N(1, i-1)
		w21I(f, "i", i)
		w21I(f, "j", j)
		return w21P(kp,
			fmt.Sprintf("Over one full period the integral is %d and over the first half it is %d. What is the integral over the second half? (enter a number)", i, j),
			strconv.Itoa(i-j),
			fmt.Sprintf("The halves add to the whole, so the second half contributes %d-%d=%d.", i, j, i-j), f)
	default:
		t := w21N(1, 10)
		q := w21N(1, 12)
		i := t * q
		w21I(f, "i", i)
		w21I(f, "t", t)
		return w21P(kp,
			fmt.Sprintf("A periodic function has period %d and integral %d over that period. What is its average value? (enter a number)", t, i),
			strconv.Itoa(q),
			fmt.Sprintf("The average is the integral divided by the period: %d/%d=%d.", i, t, q), f)
	}
}
