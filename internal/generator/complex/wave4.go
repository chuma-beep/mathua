package complex

// Wave 4 complex analysis: twenty concepts. Every family asks for an integer
// invariant or a small exact value, and publishes the parameters in Facts so
// wave4_test.go can recompute the answer.

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave4(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"complex.adv.roots_unity":               "roots_unity",
		"complex.adv.branches":                  "branches",
		"complex.adv.complex_powers":            "complex_powers",
		"complex.adv.trig_zeros":                "trig_zeros",
		"complex.adv.harmonic_conjugate":        "harmonic_conjugate",
		"complex.adv.poisson_kernel":            "poisson_kernel",
		"complex.adv.rouche":                    "rouche",
		"complex.adv.open_mapping":              "open_mapping",
		"complex.adv.casorati":                  "casorati",
		"complex.adv.picard":                    "picard",
		"complex.adv.analytic_continuation":     "analytic_continuation",
		"complex.adv.monodromy":                 "monodromy",
		"complex.adv.riemann_surface":           "riemann_surface",
		"complex.adv.elliptic":                  "elliptic",
		"complex.adv.modular_forms":             "modular_forms",
		"complex.adv.gamma":                     "gamma",
		"complex.adv.weierstrass_factorization": "weierstrass_factorization",
		"complex.adv.mittag_leffler":            "mittag_leffler",
		"complex.adv.schwarz_christoffel":       "schwarz_christoffel",
		"complex.adv.jensen":                    "jensen",
	} {
		reg.Register(id, &w4ComplexGen{kind: kind})
	}
}

type w4ComplexGen struct{ kind string }

func (g *w4ComplexGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "roots_unity":
		return w4RootsUnity(kp)
	case "branches":
		return w4Branches(kp)
	case "complex_powers":
		return w4ComplexPowers(kp)
	case "trig_zeros":
		return w4TrigZeros(kp)
	case "harmonic_conjugate":
		return w4HarmonicConjugate(kp)
	case "poisson_kernel":
		return w4PoissonKernel(kp)
	case "rouche":
		return w4Rouche(kp)
	case "open_mapping":
		return w4OpenMapping(kp)
	case "casorati":
		return w4Casorati(kp)
	case "picard":
		return w4Picard(kp)
	case "analytic_continuation":
		return w4AnalyticContinuation(kp)
	case "monodromy":
		return w4Monodromy(kp)
	case "riemann_surface":
		return w4RiemannSurface(kp)
	case "elliptic":
		return w4Elliptic(kp)
	case "modular_forms":
		return w4ModularForms(kp)
	case "gamma":
		return w4Gamma(kp)
	case "weierstrass_factorization":
		return w4Weierstrass(kp)
	case "mittag_leffler":
		return w4MittagLeffler(kp)
	case "schwarz_christoffel":
		return w4SchwarzChristoffel(kp)
	case "jensen":
		return w4Jensen(kp)
	}
	return generator.Problem{Question: "What is 2 + 2? (enter a number)", Answer: "4", Explanation: "2 + 2 = 4."}
}

func w4C(o int, q, a, ex string, facts map[string]string) generator.Problem {
	if facts == nil {
		facts = map[string]string{}
	}
	facts["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: facts, KPOrdinal: o}
}

func w4CI(m map[string]string, k string, v int) { m[k] = strconv.Itoa(v) }

func w4CR(lo, hi int) int {
	if hi <= lo {
		return lo
	}
	return lo + rand.Intn(hi-lo+1)
}

func w4Phi(n int) int {
	c := 0
	for k := 1; k <= n; k++ {
		if gcdInt(k, n) == 1 {
			c++
		}
	}
	return c
}

func w4Fact(n int) int {
	r := 1
	for i := 2; i <= n; i++ {
		r *= i
	}
	return r
}

func w4RootsUnity(kp int) generator.Problem {
	n := w4CR(3, 9)
	f := map[string]string{}
	w4CI(f, "n", n)
	switch kp {
	case 1:
		return w4C(kp, fmt.Sprintf("How many distinct n-th roots of unity are there for n=%d? (enter a number)", n), strconv.Itoa(n),
			fmt.Sprintf("The equation z^{%d}=1 has exactly %d distinct roots in the complex plane.", n, n), f)
	case 2:
		return w4C(kp, fmt.Sprintf("For n=%d, what is the sum of all n-th roots of unity? (enter a number)", n), "0",
			"The roots are the vertices of a regular polygon centred at the origin, so their vector sum is 0.", f)
	default:
		p := w4Phi(n)
		return w4C(kp, fmt.Sprintf("For n=%d, how many n-th roots of unity are primitive? (enter a number)", n), strconv.Itoa(p),
			fmt.Sprintf("A primitive root has order exactly %d. Their number is Euler's totient phi(%d)=%d.", n, n, p), f)
	}
}

func w4Branches(kp int) generator.Problem {
	n := w4CR(2, 7)
	f := map[string]string{}
	w4CI(f, "n", n)
	switch kp {
	case 1:
		return w4C(kp, fmt.Sprintf("What is the principal argument of -%d, in degrees? (enter a number)", n), "180",
			"The point is on the negative real axis. Its principal argument is 180 degrees.", f)
	case 2:
		return w4C(kp, fmt.Sprintf("On the circle |z|=%d, how many points lie on the principal branch cut of Log? (enter a number)", n), "1",
			"The principal logarithm is cut along the negative real axis. That ray meets the circle at one point.", f)
	default:
		return w4C(kp, fmt.Sprintf("For z=-%d+0.1i and z=-%d-0.1i, is the principal logarithm continuous across the cut? Enter 1 for yes and 0 for no.", n, n), "0",
			"The principal arguments approach +pi and -pi from the two sides, so the values jump by 2 pi i.", f)
	}
}

func w4ComplexPowers(kp int) generator.Problem {
	n := w4CR(2, 9)
	f := map[string]string{}
	w4CI(f, "n", n)
	re, im := 0, 0
	for k := 0; k <= n; k++ {
		c := 1
		for i := 1; i <= k; i++ {
			c = c * (n - i + 1) / i
		}
		switch k % 4 {
		case 0:
			re += c
		case 1:
			im += c
		case 2:
			re -= c
		case 3:
			im -= c
		}
	}
	w4CI(f, "real", re)
	w4CI(f, "imag", im)
	switch kp {
	case 1:
		return w4C(kp, fmt.Sprintf("What is the real part of (1+i)^{%d}? (enter a number)", n), strconv.Itoa(re),
			fmt.Sprintf("Use (1+i)^2=2i and expand. The real part is %d.", re), f)
	case 2:
		return w4C(kp, fmt.Sprintf("What is the imaginary part of (1+i)^{%d}? (enter a number)", n), strconv.Itoa(im),
			fmt.Sprintf("Use (1+i)^2=2i and expand. The imaginary part is %d.", im), f)
	default:
		return w4C(kp, fmt.Sprintf("What is the squared modulus of (1+i)^{%d}? (enter a number)", n), strconv.Itoa(1<<n),
			fmt.Sprintf("The modulus of 1+i is sqrt(2), so the squared modulus of the power is 2^{%d}=%d.", n, 1<<n), f)
	}
}

func w4TrigZeros(kp int) generator.Problem {
	k := w4CR(1, 6)
	f := map[string]string{}
	w4CI(f, "k", k)
	switch kp {
	case 1:
		return w4C(kp, fmt.Sprintf("How many zeros does sin z have in the interval [-%d pi, %d pi]? (enter a number)", k, k), strconv.Itoa(2*k+1),
			fmt.Sprintf("The zeros are the integer multiples of pi. There are %d of them in that interval.", 2*k+1), f)
	case 2:
		return w4C(kp, fmt.Sprintf("For sin(%d z), the period is 2 pi divided by %d. What is %d times that period, measured in units of pi? (enter a number)", k, k, k), "2",
			fmt.Sprintf("The period is 2 pi/%d. Multiplying by %d gives 2 pi, which is 2 units of pi.", k, k), f)
	default:
		return w4C(kp, fmt.Sprintf("How many zeros does cos z have in the interval [0, %d pi]? (enter a number)", 2*k), strconv.Itoa(2*k),
			fmt.Sprintf("The zeros of cosine are pi/2 plus integer multiples of pi. There are %d in the interval.", 2*k), f)
	}
}

func w4HarmonicConjugate(kp int) generator.Problem {
	a := w4CR(1, 6)
	f := map[string]string{}
	w4CI(f, "a", a)
	switch kp {
	case 1:
		return w4C(kp, fmt.Sprintf("For u=%d(x^2-y^2), what is the coefficient of xy in a harmonic conjugate v? (enter a number)", a), strconv.Itoa(2*a),
			fmt.Sprintf("The Cauchy-Riemann equations give v=2*%d xy plus a constant, so the coefficient is %d.", a, 2*a), f)
	case 2:
		return w4C(kp, fmt.Sprintf("Is u=%d(x^2-y^2) harmonic? Enter 1 for yes and 0 for no.", a), "1",
			fmt.Sprintf("The Laplacian is %d(2-2)=0, so the function is harmonic.", a), f)
	default:
		return w4C(kp, fmt.Sprintf("For u=%d(x^2-y^2), a harmonic conjugate is v=2*%d xy. What is v(1,2)? (enter a number)", a, a), strconv.Itoa(4*a),
			fmt.Sprintf("Substitute x=1 and y=2: v=2*%d*1*2=%d.", a, 4*a), f)
	}
}

func w4PoissonKernel(kp int) generator.Problem {
	rs := [][2]int{{1, 3}, {1, 2}, {3, 5}, {2, 3}, {5, 7}, {3, 4}, {4, 5}}
	r := rs[w4CR(0, len(rs)-1)]
	c := w4CR(1, 5)
	f := map[string]string{}
	w4CI(f, "r_num", r[0])
	w4CI(f, "r_den", r[1])
	w4CI(f, "scale", c)
	p0 := (r[1] + r[0]) / (r[1] - r[0])
	switch kp {
	case 1:
		return w4C(kp, fmt.Sprintf("What is the Poisson kernel P_r(0) for r=%d/%d? (enter a number)", r[0], r[1]), strconv.Itoa(p0),
			fmt.Sprintf("At theta=0, P_r=(1-r^2)/(1-r)^2=(1+r)/(1-r)=%d.", p0), f)
	case 2:
		return w4C(kp, fmt.Sprintf("For r=%d/%d, what value does the Poisson integral of the constant boundary value %d approach? (enter a number)", r[0], r[1], c), strconv.Itoa(c),
			fmt.Sprintf("The Poisson kernel is an approximate identity, so the integral tends to the boundary value %d.", c), f)
	default:
		return w4C(kp, fmt.Sprintf("For r=%d/%d, what is the integral of %d P_r over a full turn? (enter a number)", r[0], r[1], c), strconv.Itoa(c),
			fmt.Sprintf("The Poisson kernel integrates to 1, so scaling it by %d gives integral %d.", c, c), f)
	}
}

func w4Rouche(kp int) generator.Problem {
	n := w4CR(4, 7)
	m := w4CR(1, 2)
	a := w4CR(2, 4)
	f := map[string]string{}
	w4CI(f, "n", n)
	w4CI(f, "m", m)
	w4CI(f, "a", a)
	switch kp {
	case 1:
		return w4C(kp, fmt.Sprintf("On |z|=1, does |z^{%d}| < |%d z^{%d}|? Enter 1 for yes and 0 for no.", n, a, m), "1",
			fmt.Sprintf("On the unit circle, |z^{%d}|=1 while |%d z^{%d}|=%d, which is larger.", n, a, m, a), f)
	case 2:
		return w4C(kp, fmt.Sprintf("How many zeros does z^{%d}+%d z^{%d} have inside |z|<1? (enter a number)", n, a, m), strconv.Itoa(m),
			fmt.Sprintf("Rouche's theorem compares the function with %d z^{%d}, which has a zero of order %d at 0.", a, m, m), f)
	default:
		return w4C(kp, fmt.Sprintf("Rouche's theorem says z^{%d}+%d z^{%d} has the same number of zeros inside |z|<1 as which monomial? Enter its exponent.", n, a, m), strconv.Itoa(m),
			fmt.Sprintf("The dominant term on the circle is %d z^{%d}, so the zero count is %d.", a, m, m), f)
	}
}

func w4OpenMapping(kp int) generator.Problem {
	a := w4CR(1, 7)
	f := map[string]string{}
	w4CI(f, "a", a)
	switch kp {
	case 1:
		return w4C(kp, fmt.Sprintf("For the nonconstant analytic function f(z)=z^2+%d, is the image of an open set open? Enter 1 for yes and 0 for no.", a), "1",
			"The open mapping theorem says every nonconstant analytic function maps open sets to open sets.", f)
	case 2:
		return w4C(kp, fmt.Sprintf("For the constant function f(z)=%d, is the image of the open unit disk open? Enter 1 for yes and 0 for no.", a), "0",
			"A constant map sends the disk to a single point, and a point is not open in the plane.", f)
	default:
		return w4C(kp, fmt.Sprintf("Does the open mapping theorem apply to f(z)=z^3+%d z? Enter 1 for yes and 0 for no.", a), "1",
			"The function is nonconstant and analytic, so the theorem applies.", f)
	}
}

func w4Casorati(kp int) generator.Problem {
	a := w4CR(1, 7)
	f := map[string]string{}
	w4CI(f, "a", a)
	switch kp {
	case 1:
		return w4C(kp, fmt.Sprintf("The Laurent expansion of f(z)=e^{%d/z} has coefficient of z^{-1} equal to what? (enter a number)", a), strconv.Itoa(a),
			fmt.Sprintf("The expansion contains the term %d/z, so the coefficient of z^{-1} is %d.", a, a), f)
	case 2:
		return w4C(kp, fmt.Sprintf("Does e^{%d/z} omit the value 0? Enter 1 for yes and 0 for no.", a), "1",
			"The exponential is never zero, so 0 is an omitted value.", f)
	default:
		return w4C(kp, fmt.Sprintf("By Casorati-Weierstrass, is the image of a punctured neighbourhood of 0 under e^{%d/z} dense in the plane? Enter 1 for yes and 0 for no.", a), "1",
			"The theorem says the image is dense for an essential singularity.", f)
	}
}

func w4Picard(kp int) generator.Problem {
	a := w4CR(1, 7)
	f := map[string]string{}
	w4CI(f, "a", a)
	switch kp {
	case 1:
		return w4C(kp, fmt.Sprintf("For e^{%d/z}, Picard's great theorem says an analytic function near an essential singularity omits at most how many values? (enter a number)", a), "1",
			"At most one value can be omitted in every punctured neighbourhood.", f)
	case 2:
		return w4C(kp, fmt.Sprintf("How many values does e^{%d/z} omit near 0? (enter a number)", a), "1",
			"The exponential omits only 0.", f)
	default:
		return w4C(kp, fmt.Sprintf("Can e^{%d/z} omit two distinct values near 0? Enter 1 for yes and 0 for no.", a), "0",
			"Picard's theorem allows at most one omitted value, and the exponential already omits 0.", f)
	}
}

func w4AnalyticContinuation(kp int) generator.Problem {
	a := w4CR(2, 8)
	f := map[string]string{}
	w4CI(f, "a", a)
	switch kp {
	case 1:
		return w4C(kp, fmt.Sprintf("What is the radius of convergence around 0 of 1/(1-z/%d)? (enter a number)", a), strconv.Itoa(a),
			fmt.Sprintf("The nearest singularity is the pole at z=%d, so the radius is %d.", a, a), f)
	case 2:
		return w4C(kp, fmt.Sprintf("The analytic continuation of 1/(1-z/%d) is evaluated at z=%d. What is the value? (enter a number)", a, a+1), strconv.Itoa(-a),
			fmt.Sprintf("1/(1-(%d)/%d)=1/(-1/%d)=-%d.", a+1, a, a, a), f)
	default:
		return w4C(kp, fmt.Sprintf("Can 1/(1-z/%d) be continued analytically through every point of the circle |z|=%d except z=%d? Enter 1 for yes and 0 for no.", a, a, a), "1",
			fmt.Sprintf("The only obstruction on the circle is the pole at z=%d.", a), f)
	}
}

func w4Monodromy(kp int) generator.Problem {
	m := w4CR(1, 6)
	n := w4CR(1, 6)
	f := map[string]string{}
	w4CI(f, "m", m)
	w4CI(f, "n", n)
	switch kp {
	case 1:
		return w4C(kp, fmt.Sprintf("When log(z^%d) is continued once around the origin, by how many times 2 pi i does it change? (enter a number)", m), strconv.Itoa(m),
			fmt.Sprintf("log(z^%d)=%d log z plus a constant, so one turn adds %d times 2 pi i.", m, m, m), f)
	case 2:
		return w4C(kp, fmt.Sprintf("After %d turns around the origin, how many times 2 pi i does log z change? (enter a number)", n), strconv.Itoa(n),
			fmt.Sprintf("Each turn adds 2 pi i, so %d turns add %d times 2 pi i.", n, n), f)
	default:
		return w4C(kp, fmt.Sprintf("Does log(z^%d) return to its original branch after one loop around the origin? Enter 1 for yes and 0 for no.", m), "0",
			"The branch changes by a nonzero multiple of 2 pi i, so monodromy is nontrivial.", f)
	}
}

func w4RiemannSurface(kp int) generator.Problem {
	n := w4CR(2, 8)
	f := map[string]string{}
	w4CI(f, "n", n)
	switch kp {
	case 1:
		return w4C(kp, fmt.Sprintf("How many sheets does the Riemann surface of z^{1/%d} have? (enter a number)", n), strconv.Itoa(n),
			fmt.Sprintf("Continuing z^{1/%d} once around 0 multiplies it by a primitive %d-th root of unity, so %d sheets are needed.", n, n, n), f)
	case 2:
		return w4C(kp, fmt.Sprintf("For z^{1/%d}, at which real point is the branch point? (enter the coordinate)", n), "0",
			"The branch point of z^{1/n} is the origin.", f)
	default:
		return w4C(kp, fmt.Sprintf("For the local model w^{%d}=z near the branch point, what is the denominator of the exponent 1/%d? (enter a number)", n, n), strconv.Itoa(n),
			fmt.Sprintf("The local model is an %d-sheeted cover, so the denominator is %d.", n, n), f)
	}
}

func w4Elliptic(kp int) generator.Problem {
	c := w4CR(1, 6)
	f := map[string]string{}
	w4CI(f, "c", c)
	switch kp {
	case 1:
		return w4C(kp, fmt.Sprintf("How many independent periods does the nonconstant elliptic function wp(z)+%d have? (enter a number)", c), "2",
			"An elliptic function is doubly periodic, so its period lattice has rank 2.", f)
	case 2:
		return w4C(kp, fmt.Sprintf("What is the sum of the residues of wp(z)+%d in one fundamental parallelogram? (enter a number)", c), "0",
			"The sum of residues is zero because the boundary integrals cancel in pairs.", f)
	default:
		return w4C(kp, fmt.Sprintf("Counting multiplicity, what is the minimum number of poles of a nonconstant elliptic function such as wp(z)+%d in one fundamental parallelogram? (enter a number)", c), "2",
			"A nonconstant elliptic function has at least two poles counted with multiplicity.", f)
	}
}

func w4ModularForms(kp int) generator.Problem {
	k := w4CR(1, 6)
	f := map[string]string{}
	w4CI(f, "k", k)
	switch kp {
	case 1:
		return w4C(kp, fmt.Sprintf("What is the weight of E4^{%d}? (enter a number)", k), strconv.Itoa(4*k),
			fmt.Sprintf("Each E4 has weight 4, so the %d-th power has weight %d.", k, 4*k), f)
	case 2:
		return w4C(kp, fmt.Sprintf("What is the weight of E6^{%d}? (enter a number)", k), strconv.Itoa(6*k),
			fmt.Sprintf("Each E6 has weight 6, so the %d-th power has weight %d.", k, 6*k), f)
	default:
		return w4C(kp, fmt.Sprintf("What is the weight of Delta^{%d}? (enter a number)", k), strconv.Itoa(12*k),
			fmt.Sprintf("Delta has weight 12, so the %d-th power has weight %d.", k, 12*k), f)
	}
}

func w4Gamma(kp int) generator.Problem {
	n := w4CR(3, 9)
	k := w4CR(1, 6)
	f := map[string]string{}
	w4CI(f, "n", n)
	w4CI(f, "k", k)
	switch kp {
	case 1:
		return w4C(kp, fmt.Sprintf("What is Gamma(%d)? (enter a number)", n), strconv.Itoa(w4Fact(n-1)),
			fmt.Sprintf("For a positive integer n, Gamma(n)=(n-1)!. Here Gamma(%d)=%d.", n, w4Fact(n-1)), f)
	case 2:
		return w4C(kp, fmt.Sprintf("How many poles does Gamma have in the set {0,-1,...,-%d}? (enter a number)", k), strconv.Itoa(k+1),
			fmt.Sprintf("Gamma has a simple pole at every nonpositive integer. The set contains %d of them.", k+1), f)
	default:
		return w4C(kp, fmt.Sprintf("What is the order of the pole of Gamma at z=-%d? (enter a number)", k), "1",
			"Every pole of Gamma is simple, so the order is 1.", f)
	}
}

func w4Weierstrass(kp int) generator.Problem {
	a := w4CR(1, 7)
	f := map[string]string{}
	w4CI(f, "a", a)
	zeroWord := "zeros"
	if 2*a-1 == 1 {
		zeroWord = "zero"
	}
	switch kp {
	case 1:
		return w4C(kp, fmt.Sprintf("What is the order of the zero of sin(%d z) at z=0? (enter a number)", a), "1",
			fmt.Sprintf("sin(%d z) behaves like %d z near 0, so the zero is simple.", a, a), f)
	case 2:
		return w4C(kp, fmt.Sprintf("How many zeros does sin(%d z) have in |z|<pi? (enter a number)", a), strconv.Itoa(2*a-1),
			fmt.Sprintf("The zeros are k pi/%d with |k|<%d. The integer k satisfying that inequality gives %d %s.", a, a, 2*a-1, zeroWord), f)
	default:
		return w4C(kp, fmt.Sprintf("For sin(%d z), the elementary factor for the zero at 0 is z. Enter 1 for yes and 0 for no.", a), "1",
			"The first Weierstrass elementary factor for a zero at 0 is z.", f)
	}
}

func w4MittagLeffler(kp int) generator.Problem {
	k := w4CR(1, 6)
	f := map[string]string{}
	w4CI(f, "k", k)
	switch kp {
	case 1:
		return w4C(kp, fmt.Sprintf("How many poles does csc z have in [-%d pi, %d pi]? (enter a number)", k, k), strconv.Itoa(2*k+1),
			fmt.Sprintf("The poles are the integer multiples of pi. There are %d of them in the interval.", 2*k+1), f)
	case 2:
		res := 1
		if k%2 == 1 {
			res = -1
		}
		return w4C(kp, fmt.Sprintf("What is the residue of csc z at z=%d pi? (enter a number)", k), strconv.Itoa(res),
			fmt.Sprintf("The residue of csc z at an integer multiple k pi is (-1)^k. Here k=%d, so the residue is %d.", k, res), f)
	default:
		return w4C(kp, fmt.Sprintf("For the doubled interval [-%d pi, %d pi], how many poles does csc z have? (enter a number)", 2*k, 2*k), strconv.Itoa(4*k+1),
			fmt.Sprintf("The integer multiples of pi in the interval give %d poles.", 4*k+1), f)
	}
}

func w4SchwarzChristoffel(kp int) generator.Problem {
	n := w4CR(3, 8)
	m := w4CR(2, 7)
	f := map[string]string{}
	w4CI(f, "n", n)
	w4CI(f, "m", m)
	switch kp {
	case 1:
		return w4C(kp, fmt.Sprintf("A Schwarz-Christoffel map sends the upper half-plane to a polygon with %d sides. What is the sum of its interior angles in units of pi? (enter a number)", n), strconv.Itoa(n-2),
			fmt.Sprintf("A polygon with %d sides has interior angles summing to (%d-2) pi, which is %d units of pi.", n, n, n-2), f)
	case 2:
		return w4C(kp, fmt.Sprintf("For an interior angle pi/%d at a vertex, the Schwarz-Christoffel exponent is alpha/pi-1. What is the denominator of that exponent when written with denominator %d? (enter a number)", m, m), strconv.Itoa(m),
			fmt.Sprintf("The exponent is 1/%d-1=(1-%d)/%d, so its denominator is %d.", m, m, m, m), f)
	default:
		return w4C(kp, fmt.Sprintf("For a polygon with %d vertices, what is the sum of the Schwarz-Christoffel exponents alpha/pi-1? (enter a number)", n), "-2",
			"The sum of the exterior angles is 2 pi, so the sum of the exponents is -2.", f)
	}
}

func w4Jensen(kp int) generator.Problem {
	a := w4CR(1, 6)
	r := w4CR(1, 6)
	f := map[string]string{}
	w4CI(f, "a", a)
	w4CI(f, "r", r)
	zeroWord := "zeros"
	if a == 1 {
		zeroWord = "zero"
	}
	switch kp {
	case 1:
		return w4C(kp, fmt.Sprintf("How many zeros does z^{%d}(z-1) have inside |z|<1/2? (enter a number)", a), strconv.Itoa(a),
			fmt.Sprintf("The factor z^{%d} contributes %d %s at 0, while z-1 has no zero inside the small disk.", a, a, zeroWord), f)
	case 2:
		return w4C(kp, fmt.Sprintf("If an analytic function has no zeros in |z|<%d, does Jensen's formula reduce to log|f(0)| = the average of log|f(%d e^{i theta})|? Enter 1 for yes and 0 for no.", r, r), "1",
			"With no zeros the sum over zeros is empty, so the formula reduces to the mean-value identity.", f)
	default:
		return w4C(kp, fmt.Sprintf("How many zeros does z^{%d}(z-1) have inside |z|<2? (enter a number)", a), strconv.Itoa(a+1),
			fmt.Sprintf("The factor z^{%d} contributes %d %s and z-1 contributes one more, all inside the disk of radius 2.", a, a, zeroWord), f)
	}
}
