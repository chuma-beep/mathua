package calculus

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

// Wave 23 adds twenty-three concepts from multivariable analysis, vector
// analysis on manifolds and advanced calculus: the Laplacian and Green
// functions, the heat and wave kernels, separation of variables and
// Sturm-Liouville theory, Legendre and Bessel functions, spherical harmonics,
// the delta function, weak derivatives and Sobolev spaces, weak solutions,
// differential forms, Stokes' theorem on manifolds, de Rham cohomology, the
// Hodge star, Christoffel symbols and the covariant derivative, the Riemann
// tensor and connection forms, the gamma function and the asymptotic expansion
// of the exponential integral.
//
// Every kind draws its operands from a per-ordinal range, answers with an
// integer and publishes every operand in Facts. The answer is always computed
// from those operands, never selected from a table of outputs. The ranges of
// different ordinals are disjoint, so two ordinals can never emit the same
// question even when the random source ignores its seed.

func registerWave23(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"calc.multi.laplace_operator":    "laplace",
		"calc.multi.green_function":      "green",
		"calc.multi.heat_kernel":         "heat",
		"calc.multi.wave_equation":       "wave",
		"calc.multi.separation":          "separation",
		"calc.multi.sturm_liouville":     "sturm_liouville",
		"calc.multi.legendre":            "legendre",
		"calc.multi.bessel":              "bessel",
		"calc.multi.spherical_harmonics": "spherical",
		"calc.multi.delta":               "delta",
		"calc.multi.weak_derivative":     "weak_derivative",
		"calc.multi.sobolev":             "sobolev",
		"calc.multi.weak_solution":       "weak_solution",
		"calc.vec.form_integral":         "form_integral",
		"calc.vec.stokes_manifold":       "stokes_manifold",
		"calc.vec.de_rham":               "de_rham",
		"calc.vec.hodge_star":            "hodge_star",
		"calc.vec.christoffel":           "christoffel",
		"calc.vec.covariant_derivative":  "covariant_derivative",
		"calc.vec.riemann_tensor":        "riemann_tensor",
		"calc.vec.connection_form":       "connection_form",
		"calc.adv.gamma_function":        "gamma_function",
		"calc.adv.asymptotic":            "asymptotic",
	} {
		reg.Register(id, &w23Gen{kind: kind})
	}
}

type w23Gen struct{ kind string }

func w23N(lo, hi int) int { return lo + rand.Intn(hi-lo+1) }

func w23I(f map[string]string, k string, v int) { f[k] = strconv.Itoa(v) }

func w23P(kp int, q, a, e string, f map[string]string) generator.Problem {
	f["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: e, Facts: f, KPOrdinal: kp}
}

func w23Fact(n int) int {
	p := 1
	for i := 2; i <= n; i++ {
		p *= i
	}
	return p
}

// w23Sign is the parity of n as 1 or -1.
func w23Sign(n int) int {
	if n%2 == 0 {
		return 1
	}
	return -1
}

func (g *w23Gen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "laplace":
		return w23Laplace(kp)
	case "green":
		return w23Green(kp)
	case "heat":
		return w23Heat(kp)
	case "wave":
		return w23Wave(kp)
	case "separation":
		return w23Separation(kp)
	case "sturm_liouville":
		return w23SturmLiouville(kp)
	case "legendre":
		return w23Legendre(kp)
	case "bessel":
		return w23Bessel(kp)
	case "spherical":
		return w23Spherical(kp)
	case "delta":
		return w23Delta(kp)
	case "weak_derivative":
		return w23WeakDerivative(kp)
	case "sobolev":
		return w23Sobolev(kp)
	case "weak_solution":
		return w23WeakSolution(kp)
	case "form_integral":
		return w23FormIntegral(kp)
	case "stokes_manifold":
		return w23StokesManifold(kp)
	case "de_rham":
		return w23DeRham(kp)
	case "hodge_star":
		return w23HodgeStar(kp)
	case "christoffel":
		return w23Christoffel(kp)
	case "covariant_derivative":
		return w23CovariantDerivative(kp)
	case "riemann_tensor":
		return w23RiemannTensor(kp)
	case "connection_form":
		return w23ConnectionForm(kp)
	case "gamma_function":
		return w23GammaFunction(kp)
	default:
		return w23Asymptotic(kp)
	}
}

// w23Laplace differentiates a quadratic twice in each variable. The Laplacian
// is 2a+2b, so the operands grow with the ordinal and the arithmetic does not.
func w23Laplace(kp int) generator.Problem {
	var a, b int
	switch kp {
	case 1:
		a, b = w23N(1, 4), w23N(1, 4)
	case 2:
		a, b = w23N(5, 9), w23N(5, 9)
	default:
		a, b = w23N(10, 15), w23N(10, 15)
	}
	ans := 2*a + 2*b
	f := map[string]string{}
	w23I(f, "a", a)
	w23I(f, "b", b)
	return w23P(kp,
		fmt.Sprintf("For \\(f(x,y)=%d x^2+%d y^2\\), enter \\(\\Delta f=f_{xx}+f_{yy}\\). (enter a number)", a, b),
		strconv.Itoa(ans),
		fmt.Sprintf("The second derivative in x is %d and the second derivative in y is %d, so the Laplacian is %d+%d=%d.", 2*a, 2*b, 2*a, 2*b, ans),
		f)
}

// w23Green reports the homogeneity exponent of the fundamental solution. In
// dimension d the m-th power of the Laplacian has exponent 2m-d. The common
// branch is the Laplacian itself, where the exponent is 2-d.
func w23Green(kp int) generator.Problem {
	var d, m int
	switch kp {
	case 1:
		d, m = w23N(2, 4), w23N(1, 2)
	case 2:
		d, m = w23N(5, 7), w23N(2, 3)
	default:
		d, m = w23N(8, 10), w23N(3, 5)
	}
	ans := 2*m - d
	op := "the Laplacian"
	if m > 1 {
		op = fmt.Sprintf("the %d-th power of the Laplacian", m)
	}
	f := map[string]string{}
	w23I(f, "d", d)
	w23I(f, "m", m)
	return w23P(kp,
		fmt.Sprintf("In dimension \\(d=%d\\), the fundamental solution of %s scales as \\(G\\sim r^p\\). Enter \\(p\\). (enter a number)", d, op),
		strconv.Itoa(ans),
		fmt.Sprintf("Each Laplacian factor contributes two to the exponent, so \\(p=2m-d=2\\cdot %d-%d=%d\\). At exponent zero the solution is logarithmic, and the exponent is still recorded as 0.", m, d, ans),
		f)
}

// w23Heat reports the exponent of the prefactor (4*pi*t)^(-d/2) for an even
// dimension, because the answer must be an integer.
func w23Heat(kp int) generator.Problem {
	var d int
	switch kp {
	case 1:
		d = 2 * w23N(1, 5)
	case 2:
		d = 2 * w23N(6, 10)
	default:
		d = 2 * w23N(11, 15)
	}
	ans := -d / 2
	f := map[string]string{}
	w23I(f, "d", d)
	return w23P(kp,
		fmt.Sprintf("On \\(\\mathbb{R}^{%d}\\) the heat kernel prefactor is \\((4\\pi t)^{-d/2}\\). Enter the exponent of \\(4\\pi t\\). (enter a number)", d),
		strconv.Itoa(ans),
		fmt.Sprintf("With \\(d=%d\\) the exponent is \\(-d/2=-%d/2=%d\\).", d, d, ans),
		f)
}

// w23Wave reads the speed off a wave equation whose squared speed is a
// perfect square, so the answer is an integer.
func w23Wave(kp int) generator.Problem {
	var c int
	switch kp {
	case 1:
		c = w23N(1, 5)
	case 2:
		c = w23N(6, 10)
	default:
		c = w23N(11, 15)
	}
	k := c * c
	f := map[string]string{}
	w23I(f, "k", k)
	return w23P(kp,
		fmt.Sprintf("The wave equation \\(u_{tt}=%d u_{xx}\\) has speed \\(c\\) with \\(c^2=%d\\). Enter \\(c\\). (enter a number)", k, k),
		strconv.Itoa(c),
		fmt.Sprintf("The coefficient of \\(u_{xx}\\) is \\(c^2=%d\\), and its square root is %d.", k, c),
		f)
}

// w23Separation either reads the eigenvalue index off the separated problem or
// counts the ordinary differential equations that separation produces.
func w23Separation(kp int) generator.Problem {
	var n, speed int
	switch kp {
	case 1:
		n, speed = w23N(1, 5), w23N(2, 6)
	case 2:
		n, speed = w23N(6, 10), w23N(7, 11)
	default:
		n, speed = w23N(11, 15), w23N(12, 16)
	}
	f := map[string]string{}
	if rand.Intn(4) == 0 {
		f["ask"] = "odes"
		w23I(f, "speed", speed)
		return w23P(kp,
			fmt.Sprintf("Separating \\(u_{tt}=%d u_{xx}\\) into \\(u=X(x)T(t)\\) gives one equation for X and one for T. Enter the number of ordinary differential equations. (enter a number)", speed),
			strconv.Itoa(2),
			"The product ansatz splits the partial differential equation into two ordinary differential equations, one in x and one in t.",
			f)
	}
	ans := n * n
	f["ask"] = "eigenvalue"
	w23I(f, "n", n)
	return w23P(kp,
		fmt.Sprintf("For \\(X''+\\lambda X=0\\) with \\(X(0)=X(L)=0\\) the eigenvalues are \\(\\lambda_n=n^2\\pi^2/L^2\\). For \\(n=%d\\), enter \\(\\lambda_n L^2/\\pi^2\\). (enter a number)", n),
		strconv.Itoa(ans),
		fmt.Sprintf("Multiply \\(\\lambda_n\\) by \\(L^2/\\pi^2\\) and the factor cancels, leaving \\(n^2=%d\\).", ans),
		f)
}

// w23SturmLiouville asks for the Dirichlet eigenvalue n^2 on (0,pi), or for
// the count of interior zeros of the n-th eigenfunction.
func w23SturmLiouville(kp int) generator.Problem {
	var n int
	switch kp {
	case 1:
		n = w23N(1, 5)
	case 2:
		n = w23N(6, 10)
	default:
		n = w23N(11, 15)
	}
	f := map[string]string{}
	w23I(f, "n", n)
	if rand.Intn(2) == 0 {
		f["ask"] = "zeros"
		ans := n - 1
		return w23P(kp,
			fmt.Sprintf("For \\(y''+\\lambda y=0\\) on \\((0,\\pi)\\) with \\(y(0)=y(\\pi)=0\\), the n-th eigenfunction has \\(n-1\\) interior zeros. For \\(n=%d\\), enter that count. (enter a number)", n),
			strconv.Itoa(ans),
			fmt.Sprintf("Subtracting one from %d gives %d interior zeros.", n, ans),
			f)
	}
	f["ask"] = "lambda"
	ans := n * n
	return w23P(kp,
		fmt.Sprintf("For \\(y''+\\lambda y=0\\) on \\((0,\\pi)\\) with \\(y(0)=y(\\pi)=0\\), the eigenvalues are \\(\\lambda_n=n^2\\). For \\(n=%d\\), enter \\(\\lambda_n\\). (enter a number)", n),
		strconv.Itoa(ans),
		fmt.Sprintf("The eigenvalue at index %d is %d squared, which is %d.", n, n, ans),
		f)
}

// w23Legendre covers the endpoint values of P_n and the recurrence
// coefficient 2n+1.
func w23Legendre(kp int) generator.Problem {
	var n int
	switch kp {
	case 1:
		n = w23N(1, 5)
	case 2:
		n = w23N(6, 10)
	default:
		n = w23N(11, 15)
	}
	f := map[string]string{}
	w23I(f, "n", n)
	switch rand.Intn(3) {
	case 0:
		f["ask"] = "at_one"
		return w23P(kp,
			fmt.Sprintf("Every Legendre polynomial satisfies \\(P_n(1)=1\\). For \\(n=%d\\), enter \\(P_n(1)\\). (enter a number)", n),
			strconv.Itoa(1),
			"At x equal to one the Legendre polynomials all take the value 1, independently of n.",
			f)
	case 1:
		f["ask"] = "at_minus_one"
		ans := w23Sign(n)
		return w23P(kp,
			fmt.Sprintf("Every Legendre polynomial satisfies \\(P_n(-1)=(-1)^n\\). For \\(n=%d\\), enter \\(P_n(-1)\\) as 1 or -1. (enter a number)", n),
			strconv.Itoa(ans),
			fmt.Sprintf("The index %d is %s, so the endpoint value is %d.", n, w23Word(n), ans),
			f)
	default:
		f["ask"] = "coefficient"
		ans := 2*n + 1
		return w23P(kp,
			fmt.Sprintf("The recurrence \\((n+1)P_{n+1}=(2n+1)xP_n-nP_{n-1}\\) carries the coefficient \\(2n+1\\). For \\(n=%d\\), enter that coefficient. (enter a number)", n),
			strconv.Itoa(ans),
			fmt.Sprintf("Twice %d plus one is %d.", n, ans),
			f)
	}
}

// w23Word names the parity of n for an explanation.
func w23Word(n int) string {
	if n%2 == 0 {
		return "even"
	}
	return "odd"
}

// w23Bessel covers the value at zero, the power of the leading term and the
// factor 2^n in the leading denominator.
func w23Bessel(kp int) generator.Problem {
	f := map[string]string{}
	switch rand.Intn(3) {
	case 0:
		var n int
		switch kp {
		case 1:
			n = w23N(0, 4)
		case 2:
			n = w23N(5, 9)
		default:
			n = w23N(10, 14)
		}
		ans := 0
		if n == 0 {
			ans = 1
		}
		f["ask"] = "at_zero"
		w23I(f, "n", n)
		return w23P(kp,
			fmt.Sprintf("For the Bessel function \\(J_n\\), enter \\(J_n(0)\\) at \\(n=%d\\). (enter a number)", n),
			strconv.Itoa(ans),
			fmt.Sprintf("Only \\(J_0(0)\\) equals one, so at index %d the value is %d.", n, ans),
			f)
	case 1:
		var n int
		switch kp {
		case 1:
			n = w23N(1, 5)
		case 2:
			n = w23N(6, 10)
		default:
			n = w23N(11, 15)
		}
		f["ask"] = "power"
		w23I(f, "n", n)
		return w23P(kp,
			fmt.Sprintf("The leading term of \\(J_n(x)\\) is proportional to \\(x^n\\). For \\(n=%d\\), enter the power of x in that term. (enter a number)", n),
			strconv.Itoa(n),
			fmt.Sprintf("The series for \\(J_n\\) starts at \\(x^n\\), so the power is n=%d.", n),
			f)
	default:
		var n int
		switch kp {
		case 1:
			n = w23N(1, 5)
		case 2:
			n = w23N(6, 10)
		default:
			n = w23N(11, 14)
		}
		ans := 1 << uint(n)
		f["ask"] = "factor"
		w23I(f, "n", n)
		return w23P(kp,
			fmt.Sprintf("The leading term of \\(J_n\\) is \\((x/2)^n/n!\\), whose denominator carries \\(2^n\\). For \\(n=%d\\), enter \\(2^n\\). (enter a number)", n),
			strconv.Itoa(ans),
			fmt.Sprintf("Two to the power %d is %d.", n, ans),
			f)
	}
}

// w23Spherical covers the multiplicity 2l+1, the parity and the inverse
// question that recovers l from the multiplicity.
func w23Spherical(kp int) generator.Problem {
	var l int
	switch kp {
	case 1:
		l = w23N(0, 4)
	case 2:
		l = w23N(5, 9)
	default:
		l = w23N(10, 14)
	}
	f := map[string]string{}
	w23I(f, "l", l)
	switch rand.Intn(3) {
	case 0:
		f["ask"] = "count"
		ans := 2*l + 1
		return w23P(kp,
			fmt.Sprintf("For angular momentum \\(l=%d\\), the magnetic quantum number m takes \\(2l+1\\) values. Enter that count. (enter a number)", l),
			strconv.Itoa(ans),
			fmt.Sprintf("Twice %d plus one is %d.", l, ans),
			f)
	case 1:
		f["ask"] = "parity"
		ans := w23Sign(l)
		return w23P(kp,
			fmt.Sprintf("A spherical harmonic of degree \\(l=%d\\) has parity \\((-1)^l\\) under reflection. Enter that parity as 1 or -1. (enter a number)", l),
			strconv.Itoa(ans),
			fmt.Sprintf("The degree %d is %s, so the parity is %d.", l, w23Word(l), ans),
			f)
	default:
		count := 2*l + 1
		f["ask"] = "degree"
		w23I(f, "count", count)
		return w23P(kp,
			fmt.Sprintf("A spherical harmonic of degree l has \\(2l+1\\) values of m. For %d values, enter l. (enter a number)", count),
			strconv.Itoa(l),
			fmt.Sprintf("Solving \\(2l+1=%d\\) gives \\(l=%d\\).", count, l),
			f)
	}
}

// w23Delta covers the sampling property, the Jacobian factor |a| and the
// exponent in the n-dimensional scaling law.
func w23Delta(kp int) generator.Problem {
	f := map[string]string{}
	switch rand.Intn(3) {
	case 0:
		var c, d, a int
		switch kp {
		case 1:
			c, d, a = w23N(1, 3), w23N(1, 4), w23N(1, 4)
		case 2:
			c, d, a = w23N(4, 6), w23N(5, 9), w23N(5, 8)
		default:
			c, d, a = w23N(7, 9), w23N(10, 15), w23N(9, 12)
		}
		ans := c*a + d
		f["ask"] = "sample"
		w23I(f, "c", c)
		w23I(f, "d", d)
		w23I(f, "a", a)
		return w23P(kp,
			fmt.Sprintf("The delta function samples a function at a point: \\(\\int f(x)\\delta(x-a)\\,dx=f(a)\\). For \\(f(x)=%d x+%d\\) and \\(a=%d\\), enter \\(f(a)\\). (enter a number)", c, d, a),
			strconv.Itoa(ans),
			fmt.Sprintf("Evaluating the linear function at %d gives %d times %d plus %d, which is %d.", a, c, a, d, ans),
			f)
	case 1:
		var a int
		switch kp {
		case 1:
			a = w23N(2, 6)
		case 2:
			a = w23N(7, 11)
		default:
			a = w23N(12, 16)
		}
		if rand.Intn(2) == 0 {
			a = -a
		}
		ans := a
		if ans < 0 {
			ans = -ans
		}
		f["ask"] = "denominator"
		w23I(f, "a", a)
		return w23P(kp,
			fmt.Sprintf("The identity \\(\\delta(ax)=\\delta(x)/|a|\\) has denominator \\(|a|\\). For \\(a=%d\\), enter \\(|a|\\). (enter a number)", a),
			strconv.Itoa(ans),
			fmt.Sprintf("The absolute value of %d is %d.", a, ans),
			f)
	default:
		var e int
		switch kp {
		case 1:
			e = w23N(2, 6)
		case 2:
			e = w23N(7, 11)
		default:
			e = w23N(12, 16)
		}
		f["ask"] = "dimension"
		w23I(f, "e", e)
		return w23P(kp,
			fmt.Sprintf("In \\(\\mathbb{R}^n\\) the delta function scales as \\(\\delta(ax)=|a|^{-e}\\delta(x)\\). For \\(e=%d\\), enter the dimension \\(n\\). (enter a number)", e),
			strconv.Itoa(e),
			fmt.Sprintf("The scaling exponent equals the dimension, so \\(n=e=%d\\).", e),
			f)
	}
}

// w23WeakDerivative covers the coefficient in the derivative of a scaled
// absolute value, the fixed coefficient for x|x| and the derivative of |x| at
// a positive point.
func w23WeakDerivative(kp int) generator.Problem {
	f := map[string]string{}
	switch rand.Intn(3) {
	case 0:
		var a int
		switch kp {
		case 1:
			a = w23N(1, 5)
		case 2:
			a = w23N(6, 10)
		default:
			a = w23N(11, 15)
		}
		ans := 2 * a
		f["ask"] = "scaled"
		w23I(f, "a", a)
		return w23P(kp,
			fmt.Sprintf("The weak derivative of \\(f(x)=%d x|x|\\) is \\(c|x|\\). Enter \\(c\\). (enter a number)", a),
			strconv.Itoa(ans),
			fmt.Sprintf("Differentiating \\(x|x|\\) away from zero gives \\(2|x|\\), and the factor %d doubles it to %d.", a, ans),
			f)
	case 1:
		var b int
		switch kp {
		case 1:
			b = w23N(1, 4)
		case 2:
			b = w23N(5, 8)
		default:
			b = w23N(9, 12)
		}
		f["ask"] = "x_abs"
		w23I(f, "b", b)
		return w23P(kp,
			fmt.Sprintf("The weak derivative of \\(f(x)=x|x|+%d\\) is \\(c|x|\\). Enter \\(c\\). (enter a number)", b),
			strconv.Itoa(2),
			fmt.Sprintf("The constant %d differentiates to zero, and \\(x|x|\\) differentiates to \\(2|x|\\), so the coefficient is 2.", b),
			f)
	default:
		var p int
		switch kp {
		case 1:
			p = w23N(1, 5)
		case 2:
			p = w23N(6, 10)
		default:
			p = w23N(11, 15)
		}
		f["ask"] = "point"
		w23I(f, "p", p)
		return w23P(kp,
			fmt.Sprintf("Enter the derivative of \\(|x|\\) at the positive point \\(x=%d\\). (enter a number)", p),
			strconv.Itoa(1),
			fmt.Sprintf("To the right of zero the absolute value is the identity, so its derivative at %d is 1.", p),
			f)
	}
}

// w23Sobolev covers the order k of H^k, the exponent 2 in the H^1 norm and
// the smallest k with 2k greater than n.
func w23Sobolev(kp int) generator.Problem {
	f := map[string]string{}
	switch rand.Intn(3) {
	case 0:
		var k int
		switch kp {
		case 1:
			k = w23N(1, 6)
		case 2:
			k = w23N(7, 12)
		default:
			k = w23N(13, 18)
		}
		f["ask"] = "order"
		w23I(f, "k", k)
		return w23P(kp,
			fmt.Sprintf("Membership in \\(H^k(\\mathbb{R})\\) requires k weak derivatives in \\(L^2\\). For \\(k=%d\\), enter the number of weak derivatives. (enter a number)", k),
			strconv.Itoa(k),
			fmt.Sprintf("The space \\(H^{%d}\\) imposes exactly %d weak derivatives.", k, k),
			f)
	case 1:
		var n int
		switch kp {
		case 1:
			n = w23N(1, 4)
		case 2:
			n = w23N(5, 8)
		default:
			n = w23N(9, 12)
		}
		f["ask"] = "exponent"
		w23I(f, "n", n)
		return w23P(kp,
			fmt.Sprintf("The Sobolev norm \\(\\|u\\|_{H^1(\\mathbb{R}^{%d})}=\\left(\\sum_{|\\alpha|\\le1}\\|D^\\alpha u\\|_2^2\\right)^{1/p}\\) raises the sum to \\(1/p\\). Enter \\(p\\). (enter a number)", n),
			strconv.Itoa(2),
			"The sum of the squared norms sits under a square root, so it is raised to the power 1/2 and p is 2.",
			f)
	default:
		var n int
		switch kp {
		case 1:
			n = w23N(1, 8)
		case 2:
			n = w23N(9, 16)
		default:
			n = w23N(17, 24)
		}
		ans := n/2 + 1
		f["ask"] = "embedding"
		w23I(f, "n", n)
		return w23P(kp,
			fmt.Sprintf("Enter the smallest integer \\(k\\) with \\(2k>%d\\). (enter a number)", n),
			strconv.Itoa(ans),
			fmt.Sprintf("Dividing %d by two gives %d with remainder, and the next integer is %d.", n, n/2, ans),
			f)
	}
}

// w23WeakSolution covers the boundary coefficient that vanishes, the
// coefficient of the symmetric bilinear form and the number of integration by
// parts steps for an operator of order 2m.
func w23WeakSolution(kp int) generator.Problem {
	f := map[string]string{}
	switch rand.Intn(3) {
	case 0:
		var a int
		switch kp {
		case 1:
			a = w23N(1, 5)
		case 2:
			a = w23N(6, 10)
		default:
			a = w23N(11, 15)
		}
		f["ask"] = "boundary"
		w23I(f, "a", a)
		return w23P(kp,
			fmt.Sprintf("For \\(-\\nabla\\cdot(a\\nabla u)=f\\) with \\(a=%d\\), integration by parts leaves the boundary term \\(a\\,\\partial_\\nu u\\,v\\), which vanishes because \\(v=0\\) on the boundary. Enter its coefficient. (enter a number)", a),
			strconv.Itoa(a),
			fmt.Sprintf("The coefficient of the boundary term is the conductivity %d, and the test function removes the term.", a),
			f)
	case 1:
		var a int
		switch kp {
		case 1:
			a = w23N(1, 5)
		case 2:
			a = w23N(6, 10)
		default:
			a = w23N(11, 15)
		}
		f["ask"] = "bilinear"
		w23I(f, "a", a)
		return w23P(kp,
			fmt.Sprintf("For \\(-a u''=f\\) the weak form is \\(a\\int u'v'=\\int fv\\). With \\(a=%d\\), enter the coefficient of the integral on the left. (enter a number)", a),
			strconv.Itoa(a),
			fmt.Sprintf("The weak form multiplies the symmetric bilinear form by %d, so that is the coefficient.", a),
			f)
	default:
		var m int
		switch kp {
		case 1:
			m = w23N(1, 5)
		case 2:
			m = w23N(6, 10)
		default:
			m = w23N(11, 15)
		}
		f["ask"] = "steps"
		w23I(f, "m", m)
		return w23P(kp,
			fmt.Sprintf("The weak form of a \\(2m\\)-th order equation needs one integration by parts per order. For \\(m=%d\\), enter the number of steps. (enter a number)", m),
			strconv.Itoa(m),
			fmt.Sprintf("There are %d orders in \\(2m\\) with \\(m=%d\\), so %d integrations by parts are needed.", m, m, m),
			f)
	}
}

// w23FormIntegral covers the unit-square integral of a constant form, a
// pullback determinant and the integral over a region of stated area.
func w23FormIntegral(kp int) generator.Problem {
	f := map[string]string{}
	switch rand.Intn(3) {
	case 0:
		var c int
		switch kp {
		case 1:
			c = w23N(1, 5)
		case 2:
			c = w23N(6, 10)
		default:
			c = w23N(11, 15)
		}
		f["ask"] = "square"
		w23I(f, "c", c)
		return w23P(kp,
			fmt.Sprintf("Enter \\(\\int_{[0,1]^2} %d\\,dx\\wedge dy\\). (enter a number)", c),
			strconv.Itoa(c),
			fmt.Sprintf("The unit square has area 1, so the integral of the constant form is %d.", c),
			f)
	case 1:
		var a, b int
		switch kp {
		case 1:
			a, b = w23N(1, 4), w23N(1, 4)
		case 2:
			a, b = w23N(5, 8), w23N(5, 8)
		default:
			a, b = w23N(9, 12), w23N(9, 12)
		}
		ans := a * b
		f["ask"] = "pullback"
		w23I(f, "a", a)
		w23I(f, "b", b)
		return w23P(kp,
			fmt.Sprintf("The map \\((u,v)\\mapsto(%d u,%d v)\\) pulls \\(dx\\wedge dy\\) back to \\(D\\,du\\wedge dv\\). Enter \\(D\\). (enter a number)", a, b),
			strconv.Itoa(ans),
			fmt.Sprintf("The pullback determinant is the product of the diagonal entries, %d times %d, which is %d.", a, b, ans),
			f)
	default:
		var c, area int
		switch kp {
		case 1:
			c, area = w23N(1, 5), w23N(1, 5)
		case 2:
			c, area = w23N(6, 10), w23N(6, 10)
		default:
			c, area = w23N(11, 15), w23N(11, 15)
		}
		ans := c * area
		f["ask"] = "flux"
		w23I(f, "c", c)
		w23I(f, "area", area)
		return w23P(kp,
			fmt.Sprintf("The 2-form \\(\\omega=%d\\,dx\\wedge dy\\) is integrated over a planar region of area %d. Enter \\(\\int\\omega\\). (enter a number)", c, area),
			strconv.Itoa(ans),
			fmt.Sprintf("The integral is the coefficient %d times the area %d, which is %d.", c, area, ans),
			f)
	}
}

// w23StokesManifold covers the dimension of a boundary, the orientation sign
// and the common value of the two sides of Stokes' theorem.
func w23StokesManifold(kp int) generator.Problem {
	var n int
	switch kp {
	case 1:
		n = w23N(2, 6)
	case 2:
		n = w23N(7, 11)
	default:
		n = w23N(12, 16)
	}
	f := map[string]string{}
	switch rand.Intn(3) {
	case 0:
		ans := n - 1
		f["ask"] = "boundary"
		w23I(f, "n", n)
		return w23P(kp,
			fmt.Sprintf("A manifold of dimension \\(n=%d\\) has a boundary of dimension \\(n-1\\). Enter that dimension. (enter a number)", n),
			strconv.Itoa(ans),
			fmt.Sprintf("Subtracting one from %d gives %d.", n, ans),
			f)
	case 1:
		ans := w23Sign(n)
		f["ask"] = "orientation"
		w23I(f, "n", n)
		return w23P(kp,
			fmt.Sprintf("With the outward normal first, the induced orientation of the boundary of \\(\\mathbb{R}^n_+\\) carries the sign \\((-1)^n\\). For \\(n=%d\\), enter that sign as 1 or -1. (enter a number)", n),
			strconv.Itoa(ans),
			fmt.Sprintf("The dimension %d is %s, so the sign is %d.", n, w23Word(n), ans),
			f)
	default:
		var v int
		switch kp {
		case 1:
			v = w23N(1, 5)
		case 2:
			v = w23N(6, 10)
		default:
			v = w23N(11, 15)
		}
		f["ask"] = "equality"
		w23I(f, "v", v)
		w23I(f, "n", n)
		return w23P(kp,
			fmt.Sprintf("Stokes' theorem gives \\(\\int_M d\\omega=\\int_{\\partial M}\\omega\\). For a form whose two integrals both equal %d, enter that common value. (enter a number)", v),
			strconv.Itoa(v),
			"The theorem forces the two integrals to agree, so both carry the stated value.",
			f)
	}
}

// w23DeRham computes Betti numbers on spaces whose homology is known exactly:
// spheres, orientable surfaces of genus g, wedges of circles and the punctured
// plane.
func w23DeRham(kp int) generator.Problem {
	f := map[string]string{}
	space := ""
	switch kp {
	case 1:
		if rand.Intn(2) == 0 {
			space = "sphere"
		} else {
			space = "genus"
		}
	case 2:
		if rand.Intn(2) == 0 {
			space = "sphere"
		} else {
			space = "wedge"
		}
	default:
		if rand.Intn(2) == 0 {
			space = "genus"
		} else {
			space = "punctured"
		}
	}
	f["space"] = space
	switch space {
	case "sphere":
		var n, k int
		switch kp {
		case 1:
			n = w23N(1, 4)
		case 2:
			n = w23N(5, 8)
		default:
			n = w23N(9, 12)
		}
		k = w23N(0, n)
		ans := 0
		if k == 0 || k == n {
			ans = 1
		}
		f["ask"] = "betti"
		w23I(f, "n", n)
		w23I(f, "k", k)
		return w23P(kp,
			fmt.Sprintf("For the sphere \\(S^{%d}\\), enter the Betti number \\(b_{%d}\\). (enter a number)", n, k),
			strconv.Itoa(ans),
			fmt.Sprintf("The homology of \\(S^{%d}\\) is one-dimensional in degrees 0 and %d and zero between them, so \\(b_{%d}=%d\\).", n, n, k, ans),
			f)
	case "genus":
		var g int
		switch kp {
		case 1:
			g = w23N(1, 5)
		case 2:
			g = w23N(6, 10)
		default:
			g = w23N(11, 15)
		}
		ans := 2 * g
		f["ask"] = "b1"
		w23I(f, "g", g)
		return w23P(kp,
			fmt.Sprintf("A closed orientable surface of genus %d has first Betti number \\(2g\\). Enter \\(b_1\\). (enter a number)", g),
			strconv.Itoa(ans),
			fmt.Sprintf("Twice the genus %d is %d.", g, ans),
			f)
	case "wedge":
		var m int
		switch kp {
		case 1:
			m = w23N(2, 6)
		case 2:
			m = w23N(7, 11)
		default:
			m = w23N(12, 16)
		}
		f["ask"] = "b1"
		w23I(f, "m", m)
		return w23P(kp,
			fmt.Sprintf("A wedge of %d circles has first Betti number equal to the number of circles. Enter \\(b_1\\). (enter a number)", m),
			strconv.Itoa(m),
			fmt.Sprintf("Each circle contributes one generator, so there are %d of them.", m),
			f)
	default:
		var k int
		switch kp {
		case 1:
			k = w23N(0, 2)
		case 2:
			k = w23N(3, 5)
		default:
			k = w23N(6, 8)
		}
		ans := 0
		if k <= 1 {
			ans = 1
		}
		f["ask"] = "betti"
		w23I(f, "k", k)
		return w23P(kp,
			fmt.Sprintf("The punctured plane \\(\\mathbb{R}^2\\setminus\\{0\\}\\) deformation retracts to a circle. Enter its Betti number \\(b_{%d}\\). (enter a number)", k),
			strconv.Itoa(ans),
			fmt.Sprintf("A circle has Betti numbers one in degree 0 and 1 and zero above, so \\(b_{%d}=%d\\).", k, ans),
			f)
	}
}

// w23HodgeStar covers the exponent k(n-k) in the square of the Hodge star and
// the sign that exponent produces.
func w23HodgeStar(kp int) generator.Problem {
	var n int
	switch kp {
	case 1:
		n = w23N(2, 4)
	case 2:
		n = w23N(5, 7)
	default:
		n = w23N(8, 10)
	}
	k := w23N(0, n)
	e := k * (n - k)
	f := map[string]string{}
	w23I(f, "n", n)
	w23I(f, "k", k)
	if rand.Intn(2) == 0 {
		f["ask"] = "exponent"
		return w23P(kp,
			fmt.Sprintf("On k-forms in \\(\\mathbb{R}^n\\) the Hodge star satisfies \\(\\star^2=(-1)^{k(n-k)}\\). For \\(n=%d\\) and \\(k=%d\\), enter \\(k(n-k)\\). (enter a number)", n, k),
			strconv.Itoa(e),
			fmt.Sprintf("The product is %d times %d, which is %d.", k, n-k, e),
			f)
	}
	ans := 1
	if e%2 == 1 {
		ans = -1
	}
	f["ask"] = "sign"
	return w23P(kp,
		fmt.Sprintf("For \\(n=%d\\) and \\(k=%d\\), enter the sign \\((-1)^{k(n-k)}\\) as 1 or -1. (enter a number)", n, k),
		strconv.Itoa(ans),
		fmt.Sprintf("The exponent \\(k(n-k)\\) is %d, which is %s, so the sign is %d.", e, w23Word(e), ans),
		f)
}

// w23Christoffel covers the symbol count n^3, the coefficient in the polar
// connection and the number of nonzero symbols for a rotationally symmetric
// metric.
func w23Christoffel(kp int) generator.Problem {
	f := map[string]string{}
	switch rand.Intn(3) {
	case 0:
		var n int
		switch kp {
		case 1:
			n = w23N(2, 6)
		case 2:
			n = w23N(7, 11)
		default:
			n = w23N(12, 16)
		}
		ans := n * n * n
		f["ask"] = "symbols"
		w23I(f, "n", n)
		return w23P(kp,
			fmt.Sprintf("In dimension n the Christoffel symbols \\(\\Gamma^k_{ij}\\) number \\(n^3\\). For \\(n=%d\\), enter the number of symbols. (enter a number)", n),
			strconv.Itoa(ans),
			fmt.Sprintf("Cubing %d gives %d.", n, ans),
			f)
	case 1:
		var m int
		switch kp {
		case 1:
			m = w23N(1, 5)
		case 2:
			m = w23N(6, 10)
		default:
			m = w23N(11, 15)
		}
		ans := -m
		f["ask"] = "polar"
		w23I(f, "m", m)
		return w23P(kp,
			fmt.Sprintf("For the metric \\(ds^2=dr^2+r^{2m}d\\theta^2\\) the connection coefficient is \\(\\Gamma^r_{\\theta\\theta}=-m r^{2m-1}\\). For \\(m=%d\\), enter the coefficient of \\(r^{2m-1}\\). (enter a number)", m),
			strconv.Itoa(ans),
			fmt.Sprintf("The coefficient is \\(-m\\), and with \\(m=%d\\) that is %d.", m, ans),
			f)
	default:
		var a int
		switch kp {
		case 1:
			a = w23N(1, 4)
		case 2:
			a = w23N(5, 8)
		default:
			a = w23N(9, 12)
		}
		f["ask"] = "nonzero"
		w23I(f, "a", a)
		return w23P(kp,
			fmt.Sprintf("For the metric \\(ds^2=dr^2+(%d r)^2d\\theta^2\\), how many Christoffel symbols are nonzero? (enter a number)", a),
			strconv.Itoa(3),
			fmt.Sprintf("The nonzero symbols are \\(\\Gamma^r_{\\theta\\theta}=-a^2r\\) and the pair \\(\\Gamma^\\theta_{r\\theta}=\\Gamma^\\theta_{\\theta r}=1/r\\), so the count is 3 for every value of a, here %d.", a),
			f)
	}
}

// w23CovariantDerivative covers the coefficient in the derivative of the
// radial basis vector, the number of components of a vector field and the
// polar connection coefficient at a point.
func w23CovariantDerivative(kp int) generator.Problem {
	f := map[string]string{}
	switch rand.Intn(3) {
	case 0:
		var phi int
		switch kp {
		case 1:
			phi = w23N(1, 4)
		case 2:
			phi = w23N(5, 8)
		default:
			phi = w23N(9, 12)
		}
		f["ask"] = "basis"
		w23I(f, "phi", phi)
		return w23P(kp,
			fmt.Sprintf("In polar coordinates \\(\\partial_\\theta\\hat e_r=\\hat e_\\theta\\) at every angle. At \\(\\theta=%d\\) radians, enter the coefficient \\(c\\) in \\(\\partial_\\theta\\hat e_r=c\\hat e_\\theta\\). (enter a number)", phi),
			strconv.Itoa(1),
			"The angular derivative of the radial unit vector is the angular unit vector, with coefficient 1 at every point.",
			f)
	case 1:
		var n int
		switch kp {
		case 1:
			n = w23N(2, 6)
		case 2:
			n = w23N(7, 11)
		default:
			n = w23N(12, 16)
		}
		f["ask"] = "components"
		w23I(f, "n", n)
		return w23P(kp,
			fmt.Sprintf("A vector field on a manifold of dimension %d has how many components? (enter a number)", n),
			strconv.Itoa(n),
			fmt.Sprintf("One component for each coordinate direction, so %d in total.", n),
			f)
	default:
		var r int
		switch kp {
		case 1:
			r = w23N(2, 6)
		case 2:
			r = w23N(7, 11)
		default:
			r = w23N(12, 16)
		}
		ans := -r
		f["ask"] = "connection"
		w23I(f, "r", r)
		return w23P(kp,
			fmt.Sprintf("For the polar metric, \\(\\Gamma^r_{\\theta\\theta}=-r\\). At \\(r=%d\\), enter \\(\\Gamma^r_{\\theta\\theta}\\). (enter a number)", r),
			strconv.Itoa(ans),
			fmt.Sprintf("Substituting %d for r gives %d.", r, ans),
			f)
	}
}

// w23RiemannTensor covers the component count n^4, the single independent
// component in two dimensions and the numerator of the sphere curvature.
func w23RiemannTensor(kp int) generator.Problem {
	var n int
	switch kp {
	case 1:
		n = w23N(2, 6)
	case 2:
		n = w23N(7, 11)
	default:
		n = w23N(12, 16)
	}
	f := map[string]string{}
	switch rand.Intn(3) {
	case 0:
		ans := n * n * n * n
		f["ask"] = "components"
		w23I(f, "n", n)
		return w23P(kp,
			fmt.Sprintf("The Riemann tensor \\(R^i{}_{jkl}\\) in dimension n has \\(n^4\\) components. For \\(n=%d\\), enter the number of components. (enter a number)", n),
			strconv.Itoa(ans),
			fmt.Sprintf("Raising %d to the fourth power gives %d.", n, ans),
			f)
	case 1:
		var a int
		switch kp {
		case 1:
			a = w23N(1, 4)
		case 2:
			a = w23N(5, 8)
		default:
			a = w23N(9, 12)
		}
		f["ask"] = "independent"
		w23I(f, "a", a)
		return w23P(kp,
			fmt.Sprintf("In dimension 2 the Riemann tensor has one independent component. For a surface with \\(R_{1212}=%d\\), enter the number of independent components. (enter a number)", a),
			strconv.Itoa(1),
			fmt.Sprintf("All curvature information in two dimensions sits in \\(R_{1212}\\); here that component is %d and the count stays 1.", a),
			f)
	default:
		var r int
		switch kp {
		case 1:
			r = w23N(2, 6)
		case 2:
			r = w23N(7, 11)
		default:
			r = w23N(12, 16)
		}
		f["ask"] = "curvature"
		w23I(f, "radius", r)
		return w23P(kp,
			fmt.Sprintf("The sphere of radius \\(R=%d\\) has sectional curvature \\(1/R^2\\). Enter the numerator of that curvature when it is written over \\(R^2\\). (enter a number)", r),
			strconv.Itoa(1),
			fmt.Sprintf("The curvature is the unit numerator over \\(R^2\\), so the numerator is 1 for every radius, here %d.", r),
			f)
	}
}

// w23ConnectionForm covers the half in the curvature expansion, the two
// Cartan structure equations and the n^2 entries of the connection matrix.
func w23ConnectionForm(kp int) generator.Problem {
	var n int
	switch kp {
	case 1:
		n = w23N(2, 6)
	case 2:
		n = w23N(7, 11)
	default:
		n = w23N(12, 16)
	}
	f := map[string]string{}
	switch rand.Intn(3) {
	case 0:
		f["ask"] = "curvature"
		w23I(f, "n", n)
		return w23P(kp,
			fmt.Sprintf("A connection on a manifold of dimension %d has curvature \\(\\Omega=d\\omega+\\frac{1}{c}[\\omega,\\omega]\\). Enter \\(c\\). (enter a number)", n),
			strconv.Itoa(2),
			"The wedge product of two one-forms is antisymmetric, so the symmetry factor is 1/2 and c is 2.",
			f)
	case 1:
		f["ask"] = "equations"
		w23I(f, "n", n)
		return w23P(kp,
			fmt.Sprintf("Cartan's structure equations on a manifold of dimension %d are the torsion equation and the curvature equation. Enter how many equations there are. (enter a number)", n),
			strconv.Itoa(2),
			"One equation defines torsion and one defines curvature, so there are two.",
			f)
	default:
		ans := n * n
		f["ask"] = "forms"
		w23I(f, "n", n)
		return w23P(kp,
			fmt.Sprintf("A connection on a manifold of dimension %d is a matrix \\(\\omega^i{}_j\\) of one-forms. Enter the number of entries \\(n^2\\). (enter a number)", n),
			strconv.Itoa(ans),
			fmt.Sprintf("The matrix has %d rows and %d columns, which is %d entries.", n, n, ans),
			f)
	}
}

// w23GammaFunction covers the factorial value, the recurrence coefficient and
// the base case of the recurrence.
func w23GammaFunction(kp int) generator.Problem {
	f := map[string]string{}
	switch rand.Intn(3) {
	case 0:
		var n int
		switch kp {
		case 1:
			n = w23N(1, 6)
		case 2:
			n = w23N(7, 12)
		default:
			n = w23N(13, 18)
		}
		ans := w23Fact(n - 1)
		f["ask"] = "value"
		w23I(f, "n", n)
		return w23P(kp,
			fmt.Sprintf("For integer n the gamma function satisfies \\(\\Gamma(n)=(n-1)!\\). For \\(n=%d\\), enter \\(\\Gamma(n)\\). (enter a number)", n),
			strconv.Itoa(ans),
			fmt.Sprintf("The factorial of %d is %d.", n-1, ans),
			f)
	case 1:
		var x int
		switch kp {
		case 1:
			x = w23N(2, 7)
		case 2:
			x = w23N(8, 13)
		default:
			x = w23N(14, 19)
		}
		f["ask"] = "coefficient"
		w23I(f, "x", x)
		return w23P(kp,
			fmt.Sprintf("The recurrence \\(\\Gamma(x+1)=x\\Gamma(x)\\) at \\(x=%d\\) has a coefficient c multiplying \\(\\Gamma(x)\\). Enter c. (enter a number)", x),
			strconv.Itoa(x),
			fmt.Sprintf("The coefficient is the argument itself, so at x=%d it is %d.", x, x),
			f)
	default:
		var k int
		switch kp {
		case 1:
			k = w23N(1, 4)
		case 2:
			k = w23N(5, 8)
		default:
			k = w23N(9, 12)
		}
		f["ask"] = "base"
		w23I(f, "k", k)
		w23I(f, "target", k+1)
		return w23P(kp,
			fmt.Sprintf("Starting from \\(\\Gamma(1)=1\\), the recurrence reaches \\(\\Gamma(%d)\\) after %d steps. Enter the value you start from. (enter a number)", k+1, k),
			strconv.Itoa(1),
			"The recurrence multiplies by successive integers, and the base case is \\(\\Gamma(1)=1\\).",
			f)
	}
}

// w23Asymptotic covers the factorial coefficient of the exponential integral,
// the alternating sign and the index where the terms stop decreasing.
func w23Asymptotic(kp int) generator.Problem {
	f := map[string]string{}
	switch rand.Intn(3) {
	case 0:
		var n int
		switch kp {
		case 1:
			n = w23N(1, 5)
		case 2:
			n = w23N(6, 10)
		default:
			n = w23N(11, 15)
		}
		ans := w23Fact(n)
		f["ask"] = "factorial"
		w23I(f, "n", n)
		return w23P(kp,
			fmt.Sprintf("The asymptotic expansion of the exponential integral is \\(\\mathrm{Ei}(x)\\sim\\frac{e^x}{x}\\sum_{n\\ge0}\\frac{n!}{x^n}\\). For \\(n=%d\\), enter \\(n!\\). (enter a number)", n),
			strconv.Itoa(ans),
			fmt.Sprintf("The factorial of %d is %d.", n, ans),
			f)
	case 1:
		var n int
		switch kp {
		case 1:
			n = w23N(1, 5)
		case 2:
			n = w23N(6, 10)
		default:
			n = w23N(11, 15)
		}
		ans := w23Sign(n)
		f["ask"] = "sign"
		w23I(f, "n", n)
		return w23P(kp,
			fmt.Sprintf("In the alternating expansion of \\(\\mathrm{Ei}(-x)\\) the n-th term carries the sign \\((-1)^n\\). For \\(n=%d\\), enter that sign as 1 or -1. (enter a number)", n),
			strconv.Itoa(ans),
			fmt.Sprintf("The index %d is %s, so the sign is %d.", n, w23Word(n), ans),
			f)
	default:
		var x int
		switch kp {
		case 1:
			x = w23N(2, 6)
		case 2:
			x = w23N(7, 11)
		default:
			x = w23N(12, 16)
		}
		ans := x - 1
		f["ask"] = "smallest"
		w23I(f, "x", x)
		return w23P(kp,
			fmt.Sprintf("For \\(\\sum_{n\\ge0}n!/x^{n+1}\\) consecutive terms have ratio \\((n+1)/x\\). For \\(x=%d\\), enter the smallest n with \\((n+1)/x\\ge1\\). (enter a number)", x),
			strconv.Itoa(ans),
			fmt.Sprintf("The inequality is n+1 at least %d, so the smallest n is %d.", x, ans),
			f)
	}
}
