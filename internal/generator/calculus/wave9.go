package calculus

// Wave 9 calculus: five review concepts, three from Calculus I and two from
// Calculus II. Every answer is an integer and every family publishes its facts.

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave9(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"calc.limit.asymptotic_equivalence": "asymptotic_equivalence",
		"calc.deriv.elasticity":             "elasticity",
		"calc.integral.leibniz":             "leibniz",
		"calc.series.cesaro":                "cesaro",
		"calc.series.dirichlet_series":      "dirichlet_series",
	} {
		reg.Register(id, &w9CalcGen{kind: kind})
	}
}

type w9CalcGen struct{ kind string }

func (g *w9CalcGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "asymptotic_equivalence":
		return w9Asymptotic(kp)
	case "elasticity":
		return w9Elasticity(kp)
	case "leibniz":
		return w9Leibniz(kp)
	case "cesaro":
		return w9Cesaro(kp)
	case "dirichlet_series":
		return w9Dirichlet(kp)
	}
	return generator.Problem{Question: "What is 2 + 2? (enter a number)", Answer: "4", Explanation: "2 + 2 = 4."}
}

func w9C(o int, q, a, ex string, facts map[string]string) generator.Problem {
	if facts == nil {
		facts = map[string]string{}
	}
	facts["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: facts, KPOrdinal: o}
}

func w9CI(m map[string]string, k string, v int) { m[k] = strconv.Itoa(v) }

func w9CR(lo, hi int) int {
	if hi <= lo {
		return lo
	}
	return lo + rand.Intn(hi-lo+1)
}

func w9CIntPow(b, e int) int {
	r := 1
	for i := 0; i < e; i++ {
		r *= b
	}
	return r
}

func w9Asymptotic(kp int) generator.Problem {
	a := w9CR(2, 12)
	b := w9CR(1, a-1)
	for a%b != 0 {
		b = w9CR(1, a-1)
	}
	p := w9CR(2, 10)
	c := w9CR(-5, 5)
	f := map[string]string{}
	w9CI(f, "a", a)
	w9CI(f, "b", b)
	w9CI(f, "p", p)
	w9CI(f, "c", c)
	g := fmt.Sprintf("%d x^%d", b, p)
	if b == 1 {
		g = fmt.Sprintf("x^%d", p)
	}
	switch kp {
	case 1:
		return w9C(kp, fmt.Sprintf("For f(x)=%d x^%d + (%d) x^%d and g(x)=%s, what is the limit of f(x)/g(x) as x tends to infinity? (enter a number)", a, p, c, p-1, g), strconv.Itoa(a/b),
			fmt.Sprintf("The leading terms dominate, so the ratio tends to %d/%d = %d.", a, b, a/b), f)
	case 2:
		return w9C(kp, fmt.Sprintf("For f(x)=%d x^%d + (%d) x^%d, what is the limit of f(x)/x^%d as x tends to infinity? (enter a number)", a, p, c, p-1, p), strconv.Itoa(a),
			fmt.Sprintf("Divide every term by x^%d. The lower term vanishes, leaving %d.", p, a), f)
	default:
		q := p + w9CR(1, 3)
		w9CI(f, "q", q)
		return w9C(kp, fmt.Sprintf("Compare x^%d and x^%d as x tends to infinity. Enter 1 if the first grows faster, -1 if the second grows faster, and 0 if they have the same order.", p, q), "-1",
			fmt.Sprintf("The larger exponent %d dominates, so x^%d grows faster than x^%d. The answer is -1.", q, q, p), f)
	}
}

func w9Elasticity(kp int) generator.Problem {
	k := w9CR(1, 6)
	a := w9CR(2, 9)
	f := map[string]string{}
	w9CI(f, "k", k)
	w9CI(f, "a", a)
	switch kp {
	case 1:
		return w9C(kp, fmt.Sprintf("Demand is q=%d p^(-%d). What is the price elasticity of demand? (enter a number)", a, k), strconv.Itoa(k),
			fmt.Sprintf("For q=A p^(-k), the elasticity is the positive constant k. Here k=%d.", k), f)
	case 2:
		code := 3
		if k > 1 {
			code = 1
		} else if k < 1 {
			code = 2
		}
		return w9C(kp, fmt.Sprintf("Demand is q=%d p^(-%d). Enter 1 if demand is elastic, 2 if it is inelastic, and 3 if it is unit elastic.", a, k), strconv.Itoa(code),
			fmt.Sprintf("Elasticity is %d. Demand is elastic above 1, inelastic below 1, and unit elastic at 1.", k), f)
	default:
		sign := 0
		if k < 1 {
			sign = 1
		} else if k > 1 {
			sign = -1
		}
		return w9C(kp, fmt.Sprintf("Demand is q=%d p^(-%d). As price increases, does revenue rise (enter 1), fall (enter -1), or stay constant (enter 0)?", a, k), strconv.Itoa(sign),
			fmt.Sprintf("Revenue is R=pq=A p^(1-k). Its derivative has the sign of 1-k. With k=%d, the sign is %d.", k, sign), f)
	}
}

func w9Leibniz(kp int) generator.Problem {
	a := w9CR(1, 6)
	b := w9CR(1, 5)
	f := map[string]string{}
	w9CI(f, "a", a)
	w9CI(f, "b", b)
	bx := fmt.Sprintf("%d a x", b)
	if b == 1 {
		bx = "a x"
	}
	switch kp {
	case 1:
		return w9C(kp, fmt.Sprintf("Let F(a)=∫_0^1 (a^2 x^2 + %s) dx. What is 6F'(a) at a=%d? (enter a number)", bx, a), strconv.Itoa(4*a+3*b),
			fmt.Sprintf("Differentiate under the integral sign and integrate: 6F'(a)=4a+%d. At a=%d this is %d.", 3*b, a, 4*a+3*b), f)
	case 2:
		return w9C(kp, fmt.Sprintf("Let F(a)=∫_0^1 (a^2 x^2 + %s) dx. What is 6F(a) at a=%d? (enter a number)", bx, a), strconv.Itoa(2*a*a+3*a*b),
			fmt.Sprintf("Integrating first gives 6F(a)=2a^2+%d a. At a=%d this is %d.", 3*b, a, 2*a*a+3*a*b), f)
	default:
		return w9C(kp, fmt.Sprintf("For f(x,a)=a^2 x^2 + %s, what is the partial derivative with respect to a at x=1 and a=%d? (enter a number)", bx, a), strconv.Itoa(2*a+b),
			fmt.Sprintf("The partial derivative is 2a x + %d. At x=1 and a=%d it is %d.", b, a, 2*a+b), f)
	}
}

func w9Cesaro(kp int) generator.Problem {
	switch kp {
	case 1:
		c := w9CR(1, 6)
		n := 2*w9CR(1, 4) + 1
		f := map[string]string{}
		w9CI(f, "c", c)
		w9CI(f, "n", n)
		return w9C(kp, fmt.Sprintf("For the sequence a_k=%d+k, what is the Cesaro mean of the first %d terms? (enter a number)", c, n), strconv.Itoa(c+(n+1)/2),
			fmt.Sprintf("The average of the first %d terms is %d+(%d+1)/2=%d.", n, c, n, c+(n+1)/2), f)
	case 2:
		c := w9CR(1, 6)
		n := []int{3, 5, 7}[rand.Intn(3)]
		d := n * w9CR(1, 4)
		f := map[string]string{}
		w9CI(f, "c", c)
		w9CI(f, "d", d)
		w9CI(f, "n", n)
		return w9C(kp, fmt.Sprintf("For a_k=%d+(-1)^k(%d), what is the Cesaro mean of the first %d terms? (enter a number)", c, d, n), strconv.Itoa(c-d/n),
			fmt.Sprintf("The alternating perturbation contributes -%d to the sum of an odd number of terms, so the mean is %d-%d/%d=%d.", d, c, d, n, c-d/n), f)
	default:
		c := w9CR(1, 8)
		d := w9CR(1, 6)
		f := map[string]string{}
		w9CI(f, "c", c)
		w9CI(f, "d", d)
		return w9C(kp, fmt.Sprintf("For a_k=%d+%d/k, what is the limit of the Cesaro means of the first n terms as n tends to infinity? (enter a number)", c, d), strconv.Itoa(c),
			fmt.Sprintf("The average of %d/k over k=1..n tends to 0, so the Cesaro means tend to %d.", d, c), f)
	}
}

func w9Dirichlet(kp int) generator.Problem {
	switch kp {
	case 1:
		n := w9CR(2, 6)
		p := w9CR(1, 4)
		f := map[string]string{}
		w9CI(f, "n", n)
		w9CI(f, "p", p)
		return w9C(kp, fmt.Sprintf("In the Dirichlet series with a_n=n^%d, what is the coefficient of %d^(-s)? (enter a number)", p, n), strconv.Itoa(w9CIntPow(n, p)),
			fmt.Sprintf("The coefficient at n=%d is n^%d=%d.", n, p, w9CIntPow(n, p)), f)
	case 2:
		p := w9CR(1, 24)
		f := map[string]string{}
		w9CI(f, "p", p)
		return w9C(kp, fmt.Sprintf("For the Dirichlet series with a_n=n^%d, what is the abscissa of convergence? (enter a number)", p), strconv.Itoa(p+1),
			fmt.Sprintf("The series is Σ n^(%d-s), which converges exactly when s>%d. The abscissa is %d.", p, p+1, p+1), f)
	default:
		c := w9CR(1, 6)
		n := w9CR(2, 6)
		f := map[string]string{}
		w9CI(f, "c", c)
		w9CI(f, "N", n)
		coef := fmt.Sprintf("%d n^(-s)", c)
		if c == 1 {
			coef = "n^(-s)"
		}
		return w9C(kp, fmt.Sprintf("For the finite Dirichlet series Σ_{n=1}^{%d} %s, what is its value at s=0? (enter a number)", n, coef), strconv.Itoa(c*n),
			fmt.Sprintf("At s=0 every term equals %d. There are %d terms, so the value is %d.", c, n, c*n), f)
	}
}
