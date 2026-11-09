package odes

// Wave 4 ODE: twenty concepts. Each family publishes the numbers its answer
// follows from, and wave4_test.go recomputes every answer from those facts.

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave4(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"ode.first_order.homogeneous_substitution": "homogeneous_substitution",
		"ode.first_order.clairaut":                 "clairaut",
		"ode.first_order.riccati":                  "riccati",
		"ode.transforms.inverse_laplace":           "inverse_laplace",
		"ode.transforms.derivative":                "laplace_derivative",
		"ode.transforms.periodic":                  "periodic",
		"ode.transforms.initial_final":             "initial_final",
		"ode.second_order.frobenius":               "frobenius",
		"ode.second_order.bessel":                  "bessel",
		"ode.second_order.legendre":                "legendre",
		"ode.numerical.multistep":                  "multistep",
		"ode.numerical.rkf":                        "rkf",
		"ode.numerical.implicit":                   "implicit_euler",
		"ode.boundary.shooting":                    "shooting",
		"ode.boundary.finite_difference":           "finite_difference",
		"ode.integral.fredholm":                    "fredholm",
		"ode.basics.delay":                         "delay",
		"ode.nonlinear.perturbation":               "perturbation",
		"ode.pde.duhamel":                          "duhamel",
		"ode.pde.images":                           "images",
	} {
		reg.Register(id, &w4ODEGen{kind: kind})
	}
}

type w4ODEGen struct{ kind string }

func (g *w4ODEGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "homogeneous_substitution":
		return w4Homogeneous(kp)
	case "clairaut":
		return w4Clairaut(kp)
	case "riccati":
		return w4Riccati(kp)
	case "inverse_laplace":
		return w4InverseLaplace(kp)
	case "laplace_derivative":
		return w4LaplaceDerivative(kp)
	case "periodic":
		return w4Periodic(kp)
	case "initial_final":
		return w4InitialFinal(kp)
	case "frobenius":
		return w4Frobenius(kp)
	case "bessel":
		return w4Bessel(kp)
	case "legendre":
		return w4Legendre(kp)
	case "multistep":
		return w4Multistep(kp)
	case "rkf":
		return w4RKF(kp)
	case "implicit_euler":
		return w4ImplicitEuler(kp)
	case "shooting":
		return w4Shooting(kp)
	case "finite_difference":
		return w4FiniteDifference(kp)
	case "fredholm":
		return w4Fredholm(kp)
	case "delay":
		return w4Delay(kp)
	case "perturbation":
		return w4Perturbation(kp)
	case "duhamel":
		return w4Duhamel(kp)
	case "images":
		return w4Images(kp)
	}
	return generator.Problem{Question: "What is 2 + 2? (enter a number)", Answer: "4", Explanation: "2 + 2 = 4."}
}

func w4OR(lo, hi int) int {
	if hi <= lo {
		return lo
	}
	return lo + rand.Intn(hi-lo+1)
}

func w4O(o int, q, a, ex string, facts map[string]string) generator.Problem {
	if facts == nil {
		facts = map[string]string{}
	}
	facts["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: facts, KPOrdinal: o}
}

func w4OI(f map[string]string, k string, v int) { f[k] = strconv.Itoa(v) }

func w4Homogeneous(kp int) generator.Problem {
	a, b, k := w4OR(1, 5), 2, w4OR(1, 5)
	f := map[string]string{}
	w4OI(f, "a", a)
	w4OI(f, "b", b)
	w4OI(f, "y1", k)
	switch kp {
	case 1:
		return w4O(kp, fmt.Sprintf("For \\(y'=(%d x + %d y)/x\\), after writing \\(v=y/x\\), what is the coefficient of \\(v\\) in the numerator?", a, b), strconv.Itoa(b),
			fmt.Sprintf("Substitute \\(y=vx\\): the numerator is \\(%dx+%d(vx)=x(%d+%dv)\\), so the coefficient is %d.", a, b, a, b, b), f)
	case 2:
		return w4O(kp, fmt.Sprintf("For \\(y'=(%d x + %d y)/x\\), after \\(v=y/x\\), what is the constant term in \\(x v'=(a+(b-1)v)\\)?", a, b), strconv.Itoa(a),
			fmt.Sprintf("The equation reduces to \\(x v'=%d+(%d-1)v\\), whose constant term is %d.", a, b, a), f)
	default:
		ans := 2 * (2*k + a)
		return w4O(kp, fmt.Sprintf("Solve \\(y'=(%d x+%d y)/x\\) with \\(y(1)=%d\\). What is \\(y(2)\\)? (enter a number)", a, b, k), strconv.Itoa(ans),
			fmt.Sprintf("With \\(b=2\\), the solution is \\(y=x(Cx-%d)\\). The condition \\(y(1)=%d\\) gives \\(C=%d\\), so \\(y(2)=%d\\).", a, k, a+k, ans), f)
	}
}

func w4Clairaut(kp int) generator.Problem {
	c := w4OR(1, 5)
	m := w4OR(1, 5)
	x0 := 2 * m
	f := map[string]string{}
	w4OI(f, "c", c)
	w4OI(f, "x0", x0)
	switch kp {
	case 1:
		ans := c*x0 + c*c
		return w4O(kp, fmt.Sprintf("For the Clairaut equation \\(y=xy'+y'^2\\), the general solution is \\(y=Cx+C^2\\). What is \\(y(%d)\\) when \\(C=%d\\)?", x0, c), strconv.Itoa(ans),
			fmt.Sprintf("Substitute \\(x=%d\\) and \\(C=%d\\): \\(y=%d\\cdot%d+%d^2=%d\\).", x0, c, c, x0, c, ans), f)
	case 2:
		ans := -(x0 * x0) / 4
		return w4O(kp, fmt.Sprintf("For \\(y=xy'+y'^2\\), the singular solution is \\(y=-x^2/4\\). What is \\(y(%d)\\)?", x0), strconv.Itoa(ans),
			fmt.Sprintf("At \\(x=%d\\), \\(-x^2/4=-%d/4=%d\\).", x0, x0*x0, ans), f)
	default:
		return w4O(kp, fmt.Sprintf("For the singular solution \\(y=-x^2/4\\), what is \\(y+x^2/4\\) at \\(x=%d\\)?", x0), "0",
			"The singular solution satisfies \\(y+x^2/4=0\\) at every point.", f)
	}
}

func w4Riccati(kp int) generator.Problem {
	a := w4OR(1, 7)
	f := map[string]string{}
	w4OI(f, "a", a)
	lead := fmt.Sprintf("%d y^2", a)
	if a == 1 {
		lead = "y^2"
	}
	switch kp {
	case 1:
		return w4O(kp, fmt.Sprintf("For \\(y'=%s-%d\\), what constant particular solution \\(y=k\\) works?", lead, a), "1",
			fmt.Sprintf("\\(y=1\\) gives \\(y'=0\\) and \\(%d\\cdot1^2-%d=0\\).", a, a), f)
	case 2:
		return w4O(kp, fmt.Sprintf("For \\(y'=%d y^2-%d\\), after \\(y=1+1/u\\), what is the coefficient of \\(u\\) in the resulting linear equation?", a, a), strconv.Itoa(2*a),
			fmt.Sprintf("The substitution gives \\(u'+%du=-%d\\), so the coefficient is %d.", 2*a, a, 2*a), f)
	default:
		return w4O(kp, fmt.Sprintf("For \\(y'=%d y^2-%d\\), the substitution \\(y=1+1/u\\) gives \\(u'+%du=-%d\\). What is the forcing constant on the right?", a, a, 2*a, a), strconv.Itoa(-a),
			fmt.Sprintf("The linear equation is \\(u'+%du=-%d\\), so the constant is %d.", 2*a, a, -a), f)
	}
}

func w4InverseLaplace(kp int) generator.Problem {
	a := w4OR(1, 5)
	b := a + w4OR(1, 4)
	A := w4OR(1, 4)
	B := w4OR(1, 4)
	c := w4OR(1, 4)
	f := map[string]string{}
	w4OI(f, "a", a)
	w4OI(f, "b", b)
	w4OI(f, "A", A)
	w4OI(f, "B", B)
	w4OI(f, "c", c)
	switch kp {
	case 1:
		return w4O(kp, fmt.Sprintf("For \\(F(s)=1/(s-%d)\\), what is \\(f(0)\\)?", a), "1",
			fmt.Sprintf("\\(f(t)=e^{%dt}\\), so \\(f(0)=1\\).", a), f)
	case 2:
		return w4O(kp, fmt.Sprintf("For \\(F(s)=%d/(s-%d)+%d/(s-%d)\\), what is \\(f(0)\\)?", A, a, B, b), strconv.Itoa(A+B),
			fmt.Sprintf("The inverse is \\(%de^{%dt}+%de^{%dt}\\), whose value at 0 is %d.", A, a, B, b, A+B), f)
	default:
		t := c + w4OR(1, 2)
		w4OI(f, "t", t)
		return w4O(kp, fmt.Sprintf("For \\(F(s)=e^{-%ds}/s\\), what is \\(f(%d)\\)?", c, t), "1",
			fmt.Sprintf("The inverse is the unit step at \\(t=%d\\), so it equals 1 after the step.", c), f)
	}
}

func w4LaplaceDerivative(kp int) generator.Problem {
	a := w4OR(1, 4)
	y0 := a * w4OR(1, 4)
	f := map[string]string{}
	w4OI(f, "a", a)
	w4OI(f, "y0", y0)
	switch kp {
	case 1:
		return w4O(kp, fmt.Sprintf("For \\(y(0)=%d\\), what is the constant term in \\(\\mathcal{L}\\{y'\\}\\)?", y0), strconv.Itoa(-y0),
			fmt.Sprintf("\\(\\mathcal{L}\\{y'\\}=sY(s)-y(0)=sY(s)-%d\\).", y0), f)
	case 2:
		return w4O(kp, fmt.Sprintf("For \\(y(0)=%d\\), what is the coefficient of \\(s\\) in \\(\\mathcal{L}\\{y''\\}\\)?", y0), strconv.Itoa(-y0),
			fmt.Sprintf("\\(\\mathcal{L}\\{y''\\}=s^2Y(s)-s y(0)-y'(0)\\), so the coefficient is \\(-y(0)=%d\\).", -y0), f)
	default:
		ans := y0 / a
		return w4O(kp, fmt.Sprintf("Solve \\(y'+%d y=0\\) with \\(y(0)=%d\\). What is \\(Y(0)\\), where \\(Y=\\mathcal{L}\\{y\\}\\)?", a, y0), strconv.Itoa(ans),
			fmt.Sprintf("The transform is \\(Y(s)=%d/(s+%d)\\), so \\(Y(0)=%d/%d=%d\\).", y0, a, y0, a, ans), f)
	}
}

func w4Units(n int) string {
	if n == 1 {
		return "1 unit"
	}
	return fmt.Sprintf("%d units", n)
}

func w4Periodic(kp int) generator.Problem {
	T := w4OR(2, 6)
	duty := w4OR(1, T-1)
	high := w4OR(1, 5)
	f := map[string]string{}
	w4OI(f, "period", T)
	w4OI(f, "duty", duty)
	w4OI(f, "high", high)
	switch kp {
	case 1:
		return w4O(kp, fmt.Sprintf("A periodic function repeats every %d units. What is the length of three periods?", T), strconv.Itoa(3*T),
			fmt.Sprintf("One period is %d units, so three periods span %d units.", T, 3*T), f)
	case 2:
		t := w4OR(0, T-1)
		ans := 0
		if t < duty {
			ans = high
		}
		w4OI(f, "t", t)
		return w4O(kp, fmt.Sprintf("On one period, a function equals %d for the first %s and 0 afterwards. What is its value at \\(t=%d\\)?", high, w4Units(duty), t), strconv.Itoa(ans),
			fmt.Sprintf("The time %d lies in the %s part of the period, so the value is %d.", t, map[bool]string{true: "high", false: "zero"}[t < duty], ans), f)
	default:
		ans := high * duty
		return w4O(kp, fmt.Sprintf("A periodic function equals %d for %s and 0 for the remaining %s of its period. What is the period length times its average value?", high, w4Units(duty), w4Units(T-duty)), strconv.Itoa(ans),
			fmt.Sprintf("The average is \\(%d\\cdot%d/%d\\), so the period times the average is %d.", high, duty, T, ans), f)
	}
}

func w4InitialFinal(kp int) generator.Problem {
	b := w4OR(1, 7)
	a := b * w4OR(1, 3)
	k := w4OR(1, 4)
	f := map[string]string{}
	w4OI(f, "a", a)
	w4OI(f, "b", b)
	w4OI(f, "k", k)
	switch kp {
	case 1:
		return w4O(kp, fmt.Sprintf("For \\(Y(s)=(%d s+%d)/(s+%d)\\), what is the initial value \\(y(0)\\)?", a, b, b), strconv.Itoa(a),
			fmt.Sprintf("The initial value theorem gives \\(\\lim_{s\\to\\infty}sY(s)=%d\\).", a), f)
	case 2:
		return w4O(kp, fmt.Sprintf("For \\(Y(s)=%d/(s+%d)\\), what is the final value \\(y(\\infty)\\)?", a*k, b), strconv.Itoa(a*k/b),
			fmt.Sprintf("The final value is \\(\\lim_{s\\to0}sY(s)=%d/%d=%d\\).", a*k, b, a*k/b), f)
	default:
		stable := 1
		if b < 0 {
			stable = 0
		}
		return w4O(kp, fmt.Sprintf("Does the final value theorem apply to \\(Y(s)=1/(s+%d)\\)? Enter 1 for yes and 0 for no.", b), strconv.Itoa(stable),
			fmt.Sprintf("The pole is at \\(s=-%d\\), which lies in the left half-plane, so the final value theorem applies. The answer is %d.", b, stable), f)
	}
}

func w4Frobenius(kp int) generator.Problem {
	p0 := w4OR(-3, 3)
	q0 := w4OR(-3, 3)
	f := map[string]string{}
	w4OI(f, "p0", p0)
	w4OI(f, "q0", q0)
	sum := 1 - p0
	prod := q0
	disc := (p0-1)*(p0-1) - 4*q0
	switch kp {
	case 1:
		return w4O(kp, fmt.Sprintf("For a Frobenius equation with \\(p_0=%d\\), what is the sum of the indicial roots?", p0), strconv.Itoa(sum),
			fmt.Sprintf("The indicial equation is \\(r(r-1)+p_0r+q_0=0\\), whose root sum is \\(1-p_0=%d\\).", sum), f)
	case 2:
		return w4O(kp, fmt.Sprintf("For a Frobenius equation with \\(q_0=%d\\), what is the product of the indicial roots?", q0), strconv.Itoa(prod),
			fmt.Sprintf("The constant term of the indicial equation is \\(q_0=%d\\), which is the product of the roots.", q0), f)
	default:
		return w4O(kp, fmt.Sprintf("For a Frobenius equation with \\(p_0=%d\\) and \\(q_0=%d\\), what is the discriminant of the indicial equation?", p0, q0), strconv.Itoa(disc),
			fmt.Sprintf("The discriminant is \\((p_0-1)^2-4q_0=(%d)^2-4(%d)=%d\\).", p0-1, q0, disc), f)
	}
}

func w4Bessel(kp int) generator.Problem {
	nu := w4OR(0, 7)
	k := w4OR(1, 3)
	f := map[string]string{}
	w4OI(f, "nu", nu)
	w4OI(f, "k", k)
	switch kp {
	case 1:
		return w4O(kp, fmt.Sprintf("For Bessel's equation with parameter \\(\\nu=%d\\), what is the order of the equation?", nu), strconv.Itoa(nu),
			fmt.Sprintf("The parameter in \\(x^2y''+xy'+(x^2-\\nu^2)y=0\\) is the order, here %d.", nu), f)
	case 2:
		den := 2 * k * (nu + k)
		return w4O(kp, fmt.Sprintf("In the Bessel recurrence, the coefficient \\(a_{2k}\\) has denominator \\(2k(\\nu+k)\\). For \\(\\nu=%d\\) and \\(k=%d\\), what is that denominator?", nu, k), strconv.Itoa(den),
			fmt.Sprintf("\\(2\\cdot%d\\cdot(%d+%d)=%d\\).", k, nu, k, den), f)
	default:
		val := 0
		if nu == 0 {
			val = 1
		}
		return w4O(kp, fmt.Sprintf("What is \\(J_{%d}(0)\\)?", nu), strconv.Itoa(val),
			fmt.Sprintf("\\(J_0(0)=1\\) and every Bessel function of positive integer order vanishes at 0. For order %d the value is %d.", nu, val), f)
	}
}

func w4Legendre(kp int) generator.Problem {
	n := w4OR(0, 4)
	x := 1
	if rand.Intn(2) == 0 {
		x = -1
	}
	f := map[string]string{}
	w4OI(f, "n", n)
	w4OI(f, "x", x)
	val := 1
	if x == -1 && n%2 == 1 {
		val = -1
	}
	switch kp {
	case 1:
		return w4O(kp, fmt.Sprintf("What is \\(P_{%d}(%d)\\)?", n, x), strconv.Itoa(val),
			fmt.Sprintf("\\(P_n(1)=1\\) and \\(P_n(-1)=(-1)^n\\). For \\(n=%d\\) and \\(x=%d\\), the value is %d.", n, x, val), f)
	case 2:
		return w4O(kp, fmt.Sprintf("For Legendre polynomials, what is \\((2n+1)\\int_{-1}^{1}P_n(x)^2\\,dx\\) when \\(n=%d\\)?", n), "2",
			"The orthogonality integral is \\(2/(2n+1)\\), so multiplying by \\(2n+1\\) gives 2.", f)
	default:
		return w4O(kp, fmt.Sprintf("In the Legendre recurrence, what is the coefficient \\(2n+1\\) when \\(n=%d\\)?", n), strconv.Itoa(2*n+1),
			fmt.Sprintf("\\(2\\cdot%d+1=%d\\).", n, 2*n+1), f)
	}
}

func w4Multistep(kp int) generator.Problem {
	y := w4OR(1, 5)
	fn := w4OR(1, 5)
	fp := w4OR(1, 5)
	f := map[string]string{}
	w4OI(f, "y", y)
	w4OI(f, "fn", fn)
	w4OI(f, "fprev", fp)
	switch kp {
	case 1:
		ans := y + 3*fn - fp
		return w4O(kp, fmt.Sprintf("Use the two-step Adams-Bashforth method \\(y_{n+1}=y_n+\\frac{h}{2}(3f_n-f_{n-1})\\) with \\(h=2\\), \\(y_n=%d\\), \\(f_n=%d\\) and \\(f_{n-1}=%d\\). What is \\(y_{n+1}\\)?", y, fn, fp), strconv.Itoa(ans),
			fmt.Sprintf("\\(y_{n+1}=%d+3\\cdot%d-%d=%d\\).", y, fn, fp, ans), f)
	case 2:
		h := w4OR(1, 9)
		w4OI(f, "h", h)
		return w4O(kp, fmt.Sprintf("For the two-step Adams-Bashforth method applied with step size h=%d, what is its order?", h), "2",
			"The two-step Adams-Bashforth method is second-order accurate.", f)
	default:
		r := w4OR(-3, 3)
		w4OI(f, "root", r)
		stable := 0
		if r >= -1 && r <= 1 {
			stable = 1
		}
		return w4O(kp, fmt.Sprintf("A multistep method has a characteristic root \\(r=%d\\). Does it satisfy the zero-stability root condition? Enter 1 for yes and 0 for no.", r), strconv.Itoa(stable),
			fmt.Sprintf("Zero stability requires every root to satisfy \\(|r|\\le1\\), with roots on the unit circle simple. For \\(r=%d\\) the answer is %d.", r, stable), f)
	}
}

func w4RKF(kp int) generator.Problem {
	a := w4OR(1, 10)
	b := a + w4OR(1, 5)
	f := map[string]string{}
	w4OI(f, "low", a)
	w4OI(f, "high", b)
	switch kp {
	case 1:
		h := w4OR(1, 9)
		w4OI(f, "h", h)
		w4OI(f, "low", 4)
		w4OI(f, "high", 5)
		return w4O(kp, fmt.Sprintf("A Runge-Kutta-Fehlberg pair embeds a fourth-order method in a fifth-order method. For a trial step h=%d, what is the sum of the two orders?", h), "9",
			"The embedded pair compares orders 4 and 5, whose sum is 9.", f)
	case 2:
		return w4O(kp, fmt.Sprintf("The fourth-order estimate is %d and the fifth-order estimate is %d. What is the absolute local error estimate?", a, b), strconv.Itoa(b-a),
			fmt.Sprintf("\\(|%d-%d|=%d\\).", b, a, b-a), f)
	default:
		tol := w4OR(1, 9)
		w4OI(f, "tol", tol)
		return w4O(kp, fmt.Sprintf("In the step-size control law \\(h_{\\text{new}}=h(\\text{tol}/\\text{err})^{1/5}\\) with tol=%d, what is the exponent denominator?", tol), "5",
			"The exponent is \\(1/5\\), so the denominator is 5.", f)
	}
}

func w4ImplicitEuler(kp int) generator.Problem {
	h := w4OR(1, 3)
	lambda := w4OR(1, 4)
	k := w4OR(1, 4)
	y0 := (1 + h*lambda) * k
	f := map[string]string{}
	w4OI(f, "h", h)
	w4OI(f, "lambda", lambda)
	w4OI(f, "y0", y0)
	switch kp {
	case 1:
		return w4O(kp, fmt.Sprintf("Apply implicit Euler to \\(y'=-%d y\\) with \\(h=%d\\) and \\(y_n=%d\\). What is \\(y_{n+1}=y_n/(1+h\\lambda)\\)?", lambda, h, y0), strconv.Itoa(k),
			fmt.Sprintf("\\(y_{n+1}=%d/(1+%d\\cdot%d)=%d\\).", y0, h, lambda, k), f)
	case 2:
		return w4O(kp, fmt.Sprintf("For \\(y'=-%d y\\), \\(h=%d\\) and \\(y_n=%d\\), solve the implicit step for \\(y_{n+1}\\).", lambda, h, y0), strconv.Itoa(k),
			fmt.Sprintf("The implicit equation is \\(y_{n+1}=%d-%d\\cdot%d\\cdot y_{n+1}\\), so \\(y_{n+1}=%d\\).", y0, h, lambda, k), f)
	default:
		stable := 0
		if h*lambda > 0 {
			stable = 1
		}
		return w4O(kp, fmt.Sprintf("For implicit Euler applied to \\(y'=-%d y\\) with \\(h=%d\\), is the step stable? Enter 1 for yes and 0 for no.", lambda, h), strconv.Itoa(stable),
			fmt.Sprintf("The amplification factor is \\(1/(1+h\\lambda)\\), which has magnitude below 1 when \\(h\\lambda>0\\). Here \\(h\\lambda=%d\\), so the answer is %d.", h*lambda, stable), f)
	}
}

func w4Shooting(kp int) generator.Problem {
	a := w4OR(1, 5)
	b := w4OR(1, 5)
	m := w4OR(1, 5)
	f := map[string]string{}
	w4OI(f, "a", a)
	w4OI(f, "b", b)
	w4OI(f, "m", m)
	switch kp {
	case 1:
		return w4O(kp, fmt.Sprintf("For \\(y''=0\\), \\(y(0)=%d\\) and \\(y(1)=%d\\), what initial slope \\(m=y'(0)\\) hits the boundary?", a, b), strconv.Itoa(b-a),
			fmt.Sprintf("The solution is \\(y=%d+mx\\). At \\(x=1\\), \\(m=%d-%d=%d\\).", a, b, a, b-a), f)
	case 2:
		res := a + m - b
		return w4O(kp, fmt.Sprintf("For \\(y''=0\\), \\(y(0)=%d\\), \\(y(1)=%d\\) and trial slope \\(m=%d\\), what is the residual \\(y(1)-%d\\)?", a, b, m, b), strconv.Itoa(res),
			fmt.Sprintf("The trial solution is \\(y=%d+%dx\\), so \\(y(1)-%d=%d\\).", a, m, b, res), f)
	default:
		return w4O(kp, fmt.Sprintf("For the linear shooting problem \\(y''=0\\), \\(y(0)=%d\\), \\(y(1)=%d\\), how many correction updates are needed after the first trial?", a, b), "1",
			"For a linear boundary-value problem the update is exact, so one correction is enough.", f)
	}
}

func w4FiniteDifference(kp int) generator.Problem {
	n := w4OR(2, 9)
	k := w4OR(1, 4)
	f := map[string]string{}
	w4OI(f, "n", n)
	w4OI(f, "k", k)
	switch kp {
	case 1:
		return w4O(kp, fmt.Sprintf("A finite-difference grid uses \\(N=%d\\) equal intervals on \\([0,1]\\). How many interior grid points are there?", n), strconv.Itoa(n-1),
			fmt.Sprintf("The interior points are \\(x_1\\) through \\(x_{N-1}\\), so there are \\(N-1=%d\\).", n-1), f)
	case 2:
		return w4O(kp, fmt.Sprintf("For \\(y''=0\\) on \\([0,1]\\) with \\(N=%d\\) intervals, \\(y(0)=0\\) and \\(y(1)=%d\\), what is the boundary value \\(y_N\\)?", n, n*k), strconv.Itoa(n*k),
			fmt.Sprintf("The right boundary is \\(y(1)=%d\\).", n*k), f)
	default:
		i := w4OR(1, n-1)
		ans := i * k
		w4OI(f, "i", i)
		return w4O(kp, fmt.Sprintf("For \\(y''=0\\) with \\(N=%d\\), \\(y(0)=0\\) and \\(y(1)=%d\\), what is \\(y_{%d}\\)?", n, n*k, i), strconv.Itoa(ans),
			fmt.Sprintf("The exact linear solution is \\(y=%dx\\), so \\(y_{%d}=%d\\).", k, i, ans), f)
	}
}

func w4Fredholm(kp int) generator.Problem {
	a := w4OR(1, 9)
	b := w4OR(1, 4)
	lambda := w4OR(1, 2)
	f := map[string]string{}
	w4OI(f, "a", a)
	w4OI(f, "b", b)
	w4OI(f, "lambda", lambda)
	switch kp {
	case 1:
		return w4O(kp, fmt.Sprintf("For the Fredholm kernel \\(K(x,t)=%dxt\\), what is \\(K(1,1)\\)?", a), strconv.Itoa(a),
			fmt.Sprintf("\\(K(1,1)=%d\\cdot1\\cdot1=%d\\).", a, a), f)
	case 2:
		return w4O(kp, fmt.Sprintf("For the separable kernel \\(K(x,t)=%d x\\cdot%d t\\), what is the product of the two constant factors?", a, b), strconv.Itoa(a*b),
			fmt.Sprintf("\\(%d\\cdot%d=%d\\).", a, b, a*b), f)
	default:
		den := 1 - lambda*a*b
		return w4O(kp, fmt.Sprintf("For the rank-one Fredholm kernel \\(K(x,t)=%d x\\cdot%d t\\) and \\(\\lambda=%d\\), what is the denominator \\(1-\\lambda ab\\) of the resolvent?", a, b, lambda), strconv.Itoa(den),
			fmt.Sprintf("\\(1-%d\\cdot%d\\cdot%d=%d\\).", lambda, a, b, den), f)
	}
}

func w4Delay(kp int) generator.Problem {
	a := w4OR(1, 4)
	tau := w4OR(1, 3)
	phi := w4OR(1, 5)
	f := map[string]string{}
	w4OI(f, "a", a)
	w4OI(f, "tau", tau)
	w4OI(f, "phi", phi)
	switch kp {
	case 1:
		return w4O(kp, fmt.Sprintf("A delay equation has history \\(y(t)=%d\\) for \\(t\\le0\\). What is \\(y(-1)\\)?", phi), strconv.Itoa(phi),
			fmt.Sprintf("The history value is %d on the whole past interval.", phi), f)
	case 2:
		ans := 1 - a*tau
		return w4O(kp, fmt.Sprintf("For \\(y'=-%d y(t-%d)\\) with \\(y(t)=1\\) on \\([-%d,0]\\), what is \\(y(%d)\\)?", a, tau, tau, tau), strconv.Itoa(ans),
			fmt.Sprintf("On \\([0,%d]\\), \\(y'=-%d\\), so \\(y(t)=1-%dt\\). At \\(t=%d\\), \\(y=%d\\).", tau, a, a, tau, ans), f)
	default:
		return w4O(kp, fmt.Sprintf("For the delay equation \\(y'=-%d y(t-%d)\\), what is the coefficient \\(a\\) in the characteristic equation \\(\\lambda=-a e^{-\\lambda\\tau}\\)?", a, tau), strconv.Itoa(a),
			fmt.Sprintf("The equation is \\(y'=-%d y(t-%d)\\), so the coefficient is %d.", a, tau, a), f)
	}
}

func w4Perturbation(kp int) generator.Problem {
	eps := []int{1, 2, 3}[rand.Intn(3)]
	c := w4OR(1, 5)
	f := map[string]string{}
	w4OI(f, "eps_tenths", eps)
	w4OI(f, "c", c)
	switch kp {
	case 1:
		return w4O(kp, fmt.Sprintf("For \\(x=%d+\\epsilon x^2\\), what is the leading-order value \\(x_0\\)?", c), strconv.Itoa(c),
			fmt.Sprintf("At \\(\\epsilon=0\\), the equation is \\(x=%d\\), so \\(x_0=%d\\).", c, c), f)
	case 2:
		return w4O(kp, fmt.Sprintf("For \\(x=%d+\\epsilon x^2\\), what is the first correction \\(x_1\\)?", c), strconv.Itoa(c*c),
			fmt.Sprintf("Write \\(x=%d+\\epsilon x_1+\\cdots\\). The right side is \\(%d+\\epsilon %d^2\\), so \\(x_1=%d\\).", c, c, c, c*c), f)
	default:
		return w4O(kp, fmt.Sprintf("For \\(x=%d+\\epsilon x^2\\) with \\(\\epsilon=0.%d\\), is the perturbation parameter in the usual small range \\(|\\epsilon|<1\\)? Enter 1 for yes and 0 for no.", c, eps), "1",
			fmt.Sprintf("\\(0.%d<1\\), so the perturbation series is being used in its usual range.", eps), f)
	}
}

func w4Duhamel(kp int) generator.Problem {
	a := w4OR(1, 9)
	f := map[string]string{}
	w4OI(f, "a", a)
	switch kp {
	case 1:
		return w4O(kp, fmt.Sprintf("For \\(y'+%d y=\\delta(t)\\), what is the impulse response at \\(t=0\\)?", a), "1",
			fmt.Sprintf("The impulse response is \\(e^{-%dt}\\), whose value at 0 is 1.", a), f)
	case 2:
		if a == 1 {
			return w4O(kp, fmt.Sprintf("For \\(y'+1 y=1\\) with zero initial condition, what is the limiting value of \\(y(t)\\) as \\(t\\to\\infty\\)?"), "1",
				"The limiting value is 1/1, which is 1.", f)
		}
		return w4O(kp, fmt.Sprintf("For \\(y'+%d y=1\\) with zero initial condition, what is %d times the limiting value of \\(y(t)\\) as \\(t\\to\\infty\\)?", a, a), "1",
			fmt.Sprintf("The limiting value is \\(1/%d\\), so %d times it is 1.", a, a), f)
	default:
		return w4O(kp, fmt.Sprintf("For the Duhamel convolution \\(y(t)=\\int_0^t e^{-%d(t-s)}f(s)\\,ds\\), what is \\(y(0)\\)?", a), "0",
			"The integral has zero width at \\(t=0\\), so \\(y(0)=0\\).", f)
	}
}

func w4Images(kp int) generator.Problem {
	x0 := w4OR(1, 5)
	y0 := w4OR(1, 5)
	f := map[string]string{}
	w4OI(f, "x0", x0)
	w4OI(f, "y0", y0)
	switch kp {
	case 1:
		return w4O(kp, fmt.Sprintf("A point source sits at \\((%d,%d)\\) above the line \\(y=0\\). What is the y-coordinate of its image source?", x0, y0), strconv.Itoa(-y0),
			fmt.Sprintf("The image is reflected across \\(y=0\\), so its y-coordinate is \\(-%d\\).", y0), f)
	case 2:
		return w4O(kp, fmt.Sprintf("For a source at \\((%d,%d)\\) and zero boundary data on \\(y=0\\), what is the potential at \\((%d,0)\\)?", x0, y0, x0), "0",
			"The boundary condition sets the potential to zero on the line \\(y=0\\).", f)
	default:
		return w4O(kp, fmt.Sprintf("For a source at \\((%d,%d)\\) and its image \\((%d,-%d)\\), what is the vertical distance between the two sources?", x0, y0, x0, y0), strconv.Itoa(2*y0),
			fmt.Sprintf("The vertical separation is \\(%d-(-%d)=%d\\).", y0, y0, 2*y0), f)
	}
}
