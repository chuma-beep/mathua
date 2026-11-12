package complex

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

// Wave 19, remaining families: polynomials, residue bookkeeping and special
// functions. Every answer is an integer or an integer count so the grader and
// the test share one representation.

func registerWave19Rest(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"complex.poly.fta":               "poly_fta",
		"complex.poly.vieta_sum":         "poly_vieta_sum",
		"complex.poly.vieta_product":     "poly_vieta_product",
		"complex.poly.unity_sum":         "poly_unity_sum",
		"complex.poly.cyclotomic":        "poly_cyclotomic",
		"complex.poly.chebyshev":         "poly_chebyshev",
		"complex.poly.discriminant":      "poly_discriminant",
		"complex.poly.conjugate_pairs":   "poly_conjugate_pairs",
		"complex.poly.multiplicity":      "poly_multiplicity",
		"complex.poly.resultant":         "poly_resultant",
		"complex.meth.simple_pole":       "meth_simple_pole",
		"complex.meth.higher_pole":       "meth_higher_pole",
		"complex.meth.partial_fractions": "meth_partial_fractions",
		"complex.meth.residue_sum":       "meth_residue_sum",
		"complex.meth.winding":           "meth_winding",
		"complex.meth.pole_count":        "meth_pole_count",
		"complex.meth.zero_count":        "meth_zero_count",
		"complex.meth.order":             "meth_order",
		"complex.meth.jordan":            "meth_jordan",
		"complex.meth.branch_cut":        "meth_branch_cut",
		"complex.spec.gamma_half":        "spec_gamma_half",
		"complex.spec.beta":              "spec_beta",
		"complex.spec.zeta_even":         "spec_zeta_even",
		"complex.spec.theta":             "spec_theta",
		"complex.spec.bessel":            "spec_bessel",
		"complex.spec.airy":              "spec_airy",
		"complex.spec.hypergeometric":    "spec_hypergeometric",
		"complex.spec.polygamma":         "spec_polygamma",
		"complex.spec.erf":               "spec_erf",
		"complex.spec.polylog":           "spec_polylog",
	} {
		reg.Register(id, &w19RestGen{kind: kind})
	}
}

type w19RestGen struct{ kind string }

func (g *w19RestGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := ctx.KPOrdinal
	if o < 1 || o > 3 {
		o = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "poly_fta":
		return w19PolyFTA(o)
	case "poly_vieta_sum":
		return w19PolyVietaSum(o)
	case "poly_vieta_product":
		return w19PolyVietaProduct(o)
	case "poly_unity_sum":
		return w19PolyUnitySum(o)
	case "poly_cyclotomic":
		return w19PolyCyclotomic(o)
	case "poly_chebyshev":
		return w19PolyChebyshev(o)
	case "poly_discriminant":
		return w19PolyDiscriminant(o)
	case "poly_conjugate_pairs":
		return w19PolyConjugatePairs(o)
	case "poly_multiplicity":
		return w19PolyMultiplicity(o)
	case "poly_resultant":
		return w19PolyResultant(o)
	case "meth_simple_pole":
		return w19MethSimplePole(o)
	case "meth_higher_pole":
		return w19MethHigherPole(o)
	case "meth_partial_fractions":
		return w19MethPartialFractions(o)
	case "meth_residue_sum":
		return w19MethResidueSum(o)
	case "meth_winding":
		return w19MethWinding(o)
	case "meth_pole_count":
		return w19MethPoleCount(o)
	case "meth_zero_count":
		return w19MethZeroCount(o)
	case "meth_order":
		return w19MethOrder(o)
	case "meth_jordan":
		return w19MethJordan(o)
	case "meth_branch_cut":
		return w19MethBranchCut(o)
	case "spec_gamma_half":
		return w19SpecGammaHalf(o)
	case "spec_beta":
		return w19SpecBeta(o)
	case "spec_zeta_even":
		return w19SpecZetaEven(o)
	case "spec_theta":
		return w19SpecTheta(o)
	case "spec_bessel":
		return w19SpecBessel(o)
	case "spec_airy":
		return w19SpecAiry(o)
	case "spec_hypergeometric":
		return w19SpecHypergeometric(o)
	case "spec_polygamma":
		return w19SpecPolygamma(o)
	case "spec_erf":
		return w19SpecErf(o)
	case "spec_polylog":
		return w19SpecPolylog(o)
	}
	return generator.Problem{Question: "What is 2 + 2? (enter a number)", Answer: "4", Explanation: "2 + 2 = 4."}
}

func w19PolyFTA(o int) generator.Problem {
	d := w19N(2, 12)
	f := map[string]string{}
	w19I(f, "d", d)
	switch o {
	case 1:
		return w19P(o, fmt.Sprintf("How many complex roots, counted with multiplicity, does a degree %d polynomial have? (enter an integer)", d),
			strconv.Itoa(d), "The fundamental theorem gives one root per degree.", f)
	case 2:
		return w19P(o, fmt.Sprintf("What is the least number of roots a nonconstant polynomial of degree %d must have? (enter an integer)", d),
			"1", "Every nonconstant complex polynomial has at least one root.", f)
	default:
		return w19P(o, fmt.Sprintf("A degree %d polynomial factors into how many linear factors over the complex numbers? (enter an integer)", d),
			strconv.Itoa(d), "Each linear factor accounts for one degree.", f)
	}
}

func w19PolyVietaSum(o int) generator.Problem {
	f := map[string]string{}
	switch o {
	case 1:
		b := w19N(-12, 12)
		c := w19N(-12, 12)
		w19I(f, "b", b)
		w19I(f, "c", c)
		return w19P(o, fmt.Sprintf("For z^2 + (%d)z + (%d) = 0, what is the sum of the roots? (enter an integer)", b, c),
			strconv.Itoa(-b), "The root sum is minus the coefficient of z.", f)
	case 2:
		a := w19N(2, 9)
		b := a * w19N(-6, 6)
		c := w19N(-12, 12)
		w19I(f, "a", a)
		w19I(f, "b", b)
		w19I(f, "c", c)
		return w19P(o, fmt.Sprintf("For %dz^2 + (%d)z + (%d) = 0, what is the sum of the roots? (enter an integer)", a, b, c),
			strconv.Itoa(-b/a), "The root sum is -b/a.", f)
	default:
		n := w19N(2, 20)
		w19I(f, "n", n)
		return w19P(o, fmt.Sprintf("What is the sum of all %d-th roots of unity? (enter an integer)", n),
			"0", "The roots are the roots of z^n-1, whose z^{n-1} coefficient is zero.", f)
	}
}

func w19PolyVietaProduct(o int) generator.Problem {
	f := map[string]string{}
	switch o {
	case 1:
		b := w19N(-12, 12)
		c := w19N(-12, 12)
		w19I(f, "b", b)
		w19I(f, "c", c)
		return w19P(o, fmt.Sprintf("For z^2 + (%d)z + (%d) = 0, what is the product of the roots? (enter an integer)", b, c),
			strconv.Itoa(c), "For a monic quadratic the root product is the constant term.", f)
	case 2:
		b := w19N(-9, 9)
		c := w19N(-9, 9)
		d := w19N(-9, 9)
		w19I(f, "b", b)
		w19I(f, "c", c)
		w19I(f, "d", d)
		return w19P(o, fmt.Sprintf("For z^3 + (%d)z^2 + (%d)z + (%d) = 0, what is the product of the roots? (enter an integer)", b, c, d),
			strconv.Itoa(-d), "For a monic cubic the product is minus the constant term.", f)
	default:
		n := w19N(2, 20)
		w19I(f, "n", n)
		prod := 1
		if n%2 == 0 {
			prod = -1
		}
		return w19P(o, fmt.Sprintf("What is the product of all %d-th roots of unity? (enter an integer)", n),
			strconv.Itoa(prod), "The product is (-1)^{n+1} from the constant term of z^n-1.", f)
	}
}

func w19PolyUnitySum(o int) generator.Problem {
	f := map[string]string{}
	switch o {
	case 1:
		n := w19N(2, 24)
		w19I(f, "n", n)
		return w19P(o, fmt.Sprintf("How many n-th roots of unity are there for n = %d? (enter an integer)", n),
			strconv.Itoa(n), "There is one root for each residue modulo n.", f)
	case 2:
		n := w19N(2, 24)
		w19I(f, "n", n)
		return w19P(o, fmt.Sprintf("What is the sum of the first powers of all %d-th roots of unity? (enter an integer)", n),
			"0", "A complete set of roots sums to zero.", f)
	default:
		n := w19N(2, 12)
		j := w19N(1, 30)
		w19I(f, "n", n)
		w19I(f, "j", j)
		ans := 0
		if j%n == 0 {
			ans = n
		}
		return w19P(o, fmt.Sprintf("What is the sum of the %d-th powers of all %d-th roots of unity? (enter an integer)", j, n),
			strconv.Itoa(ans), "The sum is n when n divides the exponent and zero otherwise.", f)
	}
}

func w19Totient(n int) int {
	count := 0
	for k := 1; k <= n; k++ {
		if w19GCD(k, n) == 1 {
			count++
		}
	}
	return count
}

func w19PolyCyclotomic(o int) generator.Problem {
	f := map[string]string{}
	switch o {
	case 1:
		n := w19N(2, 60)
		w19I(f, "n", n)
		return w19P(o, fmt.Sprintf("What is the degree of the %d-th cyclotomic polynomial? (enter an integer)", n),
			strconv.Itoa(w19Totient(n)), "The degree is Euler's totient of n.", f)
	case 2:
		p := w19Prime(0, 3, 97, false)
		w19I(f, "p", p)
		return w19P(o, fmt.Sprintf("What is the degree of the %d-th cyclotomic polynomial for the prime %d? (enter an integer)", p, p),
			strconv.Itoa(p-1), "For a prime p the degree is p-1.", f)
	default:
		n := w19N(2, 60)
		w19I(f, "n", n)
		return w19P(o, fmt.Sprintf("What is the sum of the degrees of the first %d cyclotomic polynomials? (enter an integer)", n),
			strconv.Itoa(n), "The degrees sum to n because every root of unity has one order.", f)
	}
}

func w19PolyChebyshev(o int) generator.Problem {
	n := w19N(1, 20)
	f := map[string]string{}
	w19I(f, "n", n)
	switch o {
	case 1:
		return w19P(o, fmt.Sprintf("What is T_%d(1) for the Chebyshev polynomial of the first kind? (enter an integer)", n),
			"1", "T_n(cos theta)=cos(n theta), and cos(0)=1.", f)
	case 2:
		ans := 1
		if n%2 == 1 {
			ans = -1
		}
		return w19P(o, fmt.Sprintf("What is T_%d(-1)? (enter an integer)", n),
			strconv.Itoa(ans), "The value is (-1)^n.", f)
	default:
		ans := 0
		if n%2 == 0 {
			ans = 1
			if (n/2)%2 == 1 {
				ans = -1
			}
		}
		return w19P(o, fmt.Sprintf("What is T_%d(0)? (enter an integer)", n),
			strconv.Itoa(ans), "Odd degree gives zero, and even degree gives (-1)^{n/2}.", f)
	}
}

func w19PolyDiscriminant(o int) generator.Problem {
	f := map[string]string{}
	switch o {
	case 1:
		a := w19N(1, 6)
		b := w19N(-10, 10)
		c := w19N(-10, 10)
		w19I(f, "a", a)
		w19I(f, "b", b)
		w19I(f, "c", c)
		return w19P(o, fmt.Sprintf("For %dz^2 + (%d)z + (%d) = 0, what is the discriminant? (enter an integer)", a, b, c),
			strconv.Itoa(b*b-4*a*c), "The discriminant is b^2-4ac.", f)
	case 2:
		k := w19N(1, 30)
		w19I(f, "k", k)
		return w19P(o, fmt.Sprintf("Does z^2 - (%d)z + (%d) have a double root? (enter 1 for yes, 0 for no)", 2*k, k*k),
			"1", "The discriminant is (2k)^2-4k^2=0, so the root is double.", f)
	default:
		a := 1
		b := w19N(-10, 10)
		c := w19N(-10, 10)
		w19I(f, "a", a)
		w19I(f, "b", b)
		w19I(f, "c", c)
		disc := b*b - 4*a*c
		roots := 0
		if disc > 0 {
			roots = 2
		} else if disc == 0 {
			roots = 1
		}
		return w19P(o, fmt.Sprintf("How many distinct real roots does z^2 + (%d)z + (%d) have? (enter an integer)", b, c),
			strconv.Itoa(roots), "The sign of the discriminant gives two, one or no real roots.", f)
	}
}

func w19PolyConjugatePairs(o int) generator.Problem {
	f := map[string]string{}
	switch o {
	case 1:
		d := w19N(3, 12)
		r := w19N(1, d-2)
		w19I(f, "d", d)
		w19I(f, "r", r)
		return w19P(o, fmt.Sprintf("A real polynomial of degree %d has %d real roots. How many nonreal roots does it have? (enter an integer)", d, r),
			strconv.Itoa(d-r), "The remaining roots are nonreal and come in pairs.", f)
	case 2:
		d := 2*w19N(1, 15) + 1
		w19I(f, "d", d)
		return w19P(o, fmt.Sprintf("What is the least number of nonreal roots an odd-degree real polynomial of degree %d can have? (enter an integer)", d),
			"1", "An odd degree forces at least one real root, and the rest pair up.", f)
	default:
		r := w19N(0, 4)
		pairs := w19N(1, 5)
		d := r + 2*pairs
		w19I(f, "d", d)
		w19I(f, "r", r)
		return w19P(o, fmt.Sprintf("A real polynomial of degree %d has %d real roots. How many conjugate pairs of nonreal roots does it have? (enter an integer)", d, r),
			strconv.Itoa((d-r)/2), "Divide the nonreal count by two.", f)
	}
}

func w19PolyMultiplicity(o int) generator.Problem {
	f := map[string]string{}
	switch o {
	case 1:
		k := w19N(1, 30)
		w19I(f, "k", k)
		return w19P(o, fmt.Sprintf("What is the multiplicity of the root 1 of (z-1)^%d? (enter an integer)", k),
			strconv.Itoa(k), "The exponent is the multiplicity.", f)
	case 2:
		d := w19N(3, 12)
		k := w19N(2, d)
		w19I(f, "d", d)
		w19I(f, "k", k)
		return w19P(o, fmt.Sprintf("A degree %d polynomial has one root of multiplicity %d and all other roots simple. How many distinct roots does it have? (enter an integer)", d, k),
			strconv.Itoa(d-k+1), "One repeated root accounts for k of the degree.", f)
	default:
		d := w19N(3, 15)
		w19I(f, "d", d)
		return w19P(o, fmt.Sprintf("For a degree %d polynomial, what is the sum of the multiplicities of all its roots? (enter an integer)", d),
			strconv.Itoa(d), "The multiplicities add to the degree.", f)
	}
}

func w19PolyResultant(o int) generator.Problem {
	f := map[string]string{}
	switch o {
	case 1:
		a, b, c := w19N(-7, 7), w19N(-7, 7), w19N(-7, 7)
		w19I(f, "a", a)
		w19I(f, "b", b)
		w19I(f, "c", c)
		return w19P(o, fmt.Sprintf("What is the resultant of (z-%d)(z-%d) and (z-%d)? (enter an integer)", a, b, c),
			strconv.Itoa((a-c)*(b-c)), "The resultant is the product of all root differences.", f)
	case 2:
		a, b := w19N(-9, 9), w19N(-9, 9)
		w19I(f, "a", a)
		w19I(f, "b", b)
		return w19P(o, fmt.Sprintf("What is the resultant of (z-%d) and (z-%d)? (enter an integer)", a, b),
			strconv.Itoa(a-b), "For two linear factors the resultant is the single root difference.", f)
	default:
		a, b, c := w19N(-5, 5), w19N(-5, 5), w19N(-5, 5)
		if rand.Intn(2) == 0 {
			c = a
		}
		w19I(f, "a", a)
		w19I(f, "b", b)
		w19I(f, "c", c)
		shared := 0
		if c == a || c == b {
			shared = 1
		}
		return w19P(o, fmt.Sprintf("Do (z-%d)(z-%d) and (z-%d) share a root? (enter 1 for yes, 0 for no)", a, b, c),
			strconv.Itoa(shared), "A shared root makes the resultant zero.", f)
	}
}

func w19MethSimplePole(o int) generator.Problem {
	f := map[string]string{}
	switch o {
	case 1:
		a := w19N(-9, 9)
		b := w19N(1, 9)
		w19I(f, "a", a)
		w19I(f, "b", b)
		return w19P(o, fmt.Sprintf("What is the residue at z=%d of (z+%d)/(z-%d)? (enter an integer)", b, a, b),
			strconv.Itoa(a+b), "For a simple pole the residue is the numerator evaluated at the pole.", f)
	case 2:
		c := w19N(-9, 9)
		d := w19N(1, 9)
		w19I(f, "c", c)
		w19I(f, "d", d)
		return w19P(o, fmt.Sprintf("What is the residue at z=%d of (z^2+%d)/(z-%d)? (enter an integer)", d, c, d),
			strconv.Itoa(d*d+c), "Evaluate the numerator at the pole.", f)
	default:
		a, b := w19N(-9, 9), w19N(-9, 9)
		for b == a {
			b = w19N(-9, 9)
		}
		w19I(f, "a", a)
		w19I(f, "b", b)
		return w19P(o, fmt.Sprintf("How many simple poles does 1/((z-%d)(z-%d)) have? (enter an integer)", a, b),
			"2", "The two distinct linear factors give two simple poles.", f)
	}
}

func w19MethHigherPole(o int) generator.Problem {
	f := map[string]string{}
	switch o {
	case 1:
		m := w19N(1, 30)
		w19I(f, "m", m)
		return w19P(o, fmt.Sprintf("What is the order of the pole at z=0 of z^(-%d)? (enter an integer)", m),
			strconv.Itoa(m), "The order is the exponent of the denominator.", f)
	case 2:
		a := w19N(-9, 9)
		w19I(f, "a", a)
		return w19P(o, fmt.Sprintf("What is the residue at z=%d of 1/(z-%d)? (enter an integer)", a, a),
			"1", "A simple pole with unit numerator has residue one.", f)
	default:
		m := w19N(2, 30)
		w19I(f, "m", m)
		return w19P(o, fmt.Sprintf("For a pole of order %d, which Laurent exponent carries the residue? (enter the exponent)", m),
			strconv.Itoa(-1), "The residue is always the coefficient of (z-c)^{-1}.", f)
	}
}

func w19MethPartialFractions(o int) generator.Problem {
	a := w19N(-8, 8)
	b := w19N(-8, 8)
	for b == a {
		b = w19N(-8, 8)
	}
	coefA := w19N(-8, 8)
	coefB := w19N(-8, 8)
	f := map[string]string{}
	w19I(f, "a", a)
	w19I(f, "b", b)
	w19I(f, "coefA", coefA)
	w19I(f, "coefB", coefB)
	n1 := coefA + coefB
	n0 := -(coefA*b + coefB*a)
	switch o {
	case 1:
		return w19P(o, fmt.Sprintf("In the partial fractions of (%dz + (%d))/((z-%d)(z-%d)), what is the coefficient of 1/(z-%d)? (enter an integer)", n1, n0, a, b, a),
			strconv.Itoa(coefA), "The coefficient is the residue at that pole.", f)
	case 2:
		return w19P(o, fmt.Sprintf("In the partial fractions of (%dz + (%d))/((z-%d)(z-%d)), what is the coefficient of 1/(z-%d)? (enter an integer)", n1, n0, a, b, b),
			strconv.Itoa(coefB), "The coefficient is the residue at the second pole.", f)
	default:
		return w19P(o, fmt.Sprintf("In the partial fractions of (%dz + (%d))/((z-%d)(z-%d)), what is the sum of the two coefficients? (enter an integer)", n1, n0, a, b),
			strconv.Itoa(coefA+coefB), "The sum is the leading coefficient of the numerator.", f)
	}
}

func w19MethResidueSum(o int) generator.Problem {
	a := w19N(-9, 9)
	b := w19N(-9, 9)
	f := map[string]string{}
	w19I(f, "a", a)
	w19I(f, "b", b)
	switch o {
	case 1:
		return w19P(o, fmt.Sprintf("Two finite residues are %d and %d. What is their sum? (enter an integer)", a, b),
			strconv.Itoa(a+b), "Add the finite residues.", f)
	case 2:
		return w19P(o, fmt.Sprintf("Two finite residues are %d and %d. What is the residue at infinity? (enter an integer)", a, b),
			strconv.Itoa(-(a + b)), "The residue at infinity balances the finite sum.", f)
	default:
		return w19P(o, fmt.Sprintf("Two finite residues are %d and %d. What is the sum of all residues including infinity? (enter an integer)", a, b),
			"0", "Residues sum to zero on the sphere.", f)
	}
}

func w19MethWinding(o int) generator.Problem {
	k := w19N(1, 15)
	f := map[string]string{}
	w19I(f, "k", k)
	switch o {
	case 1:
		return w19P(o, fmt.Sprintf("What is the winding number of the unit circle about 0, traversed %d time(s) counterclockwise? (enter an integer)", k),
			strconv.Itoa(k), "Each counterclockwise turn adds one.", f)
	case 2:
		return w19P(o, fmt.Sprintf("What is the winding number of a circle traversed %d times clockwise? (enter an integer)", k),
			strconv.Itoa(-k), "Clockwise turns count negative.", f)
	default:
		return w19P(o, fmt.Sprintf("A contour is traversed %d times clockwise and then %d times counterclockwise. What is its winding number? (enter an integer)", k, k),
			"0", "The two contributions cancel.", f)
	}
}

func w19MethPoleCount(o int) generator.Problem {
	f := map[string]string{}
	switch o {
	case 1:
		n := w19N(2, 12)
		w19I(f, "n", n)
		return w19P(o, fmt.Sprintf("How many poles does 1/(z^%d-1) have inside a circle of radius 3? (enter an integer)", n),
			strconv.Itoa(n), "All n roots of z^n=1 lie on the unit circle.", f)
	case 2:
		r := w19N(2, 40)
		w19I(f, "r", r)
		return w19P(o, fmt.Sprintf("How many poles does 1/(z^2+1) have inside |z|<%d? (enter an integer)", r),
			"2", "The two roots are i and -i, both inside every radius greater than one.", f)
	default:
		r := w19N(2, 40)
		w19I(f, "r", r)
		return w19P(o, fmt.Sprintf("How many poles does 1/(z^2+1) have outside |z|>%d? (enter an integer)", r),
			"0", "Both poles lie on the unit circle, inside every radius greater than one.", f)
	}
}

func w19MethZeroCount(o int) generator.Problem {
	f := map[string]string{}
	switch o {
	case 1:
		n := w19N(2, 12)
		w19I(f, "n", n)
		return w19P(o, fmt.Sprintf("How many zeros does a degree %d polynomial have inside a contour containing all its roots? (enter an integer)", n),
			strconv.Itoa(n), "Count with multiplicity.", f)
	case 2:
		n := w19N(2, 12)
		r := w19N(2, 8)
		w19I(f, "n", n)
		w19I(f, "r", r)
		return w19P(o, fmt.Sprintf("How many zeros does z^%d-1 have inside |z|<%d? (enter an integer)", n, r),
			strconv.Itoa(n), "All n roots lie on the unit circle.", f)
	default:
		k := w19N(1, 9)
		w19I(f, "k", k)
		return w19P(o, fmt.Sprintf("How many zeros does e^(%dz) have in the complex plane? (enter an integer)", k),
			"0", "The exponential never vanishes.", f)
	}
}

func w19MethOrder(o int) generator.Problem {
	f := map[string]string{}
	switch o {
	case 1:
		k := w19N(1, 8)
		w19I(f, "k", k)
		return w19P(o, fmt.Sprintf("What is the order of the zero at z=0 of z^%d? (enter an integer)", k),
			strconv.Itoa(k), "The exponent is the order.", f)
	case 2:
		k1, k2 := w19N(1, 5), w19N(1, 5)
		w19I(f, "k1", k1)
		w19I(f, "k2", k2)
		return w19P(o, fmt.Sprintf("A function has a zero of order %d and another of order %d. What is the order of their product at the same point? (enter an integer)", k1, k2),
			strconv.Itoa(k1+k2), "Orders add under multiplication.", f)
	default:
		k1 := w19N(4, 9)
		k2 := w19N(1, k1-1)
		w19I(f, "k1", k1)
		w19I(f, "k2", k2)
		return w19P(o, fmt.Sprintf("A function has a zero of order %d and a pole of order %d at the same point. What is the order of the quotient? (enter an integer)", k1, k2),
			strconv.Itoa(k1-k2), "Orders subtract under division.", f)
	}
}

func w19MethJordan(o int) generator.Problem {
	f := map[string]string{}
	switch o {
	case 1:
		dn, nn := w19N(2, 6), w19N(0, 1)
		w19I(f, "dn", dn)
		w19I(f, "nn", nn)
		return w19P(o, fmt.Sprintf("The denominator has degree %d and the numerator degree %d. What is the degree gap? (enter an integer)", dn, nn),
			strconv.Itoa(dn-nn), "Subtract the numerator degree from the denominator degree.", f)
	case 2:
		k := w19N(1, 9)
		w19I(f, "k", k)
		return w19P(o, fmt.Sprintf("For a Jordan contour with denominator z^2+%d and numerator 1, what is the degree gap? (enter an integer)", k),
			"2", "The gap is 2-0.", f)
	default:
		k := w19N(1, 9)
		w19I(f, "k", k)
		return w19P(o, fmt.Sprintf("How many poles of 1/(z^2+%d) lie in the upper half-plane? (enter an integer)", k),
			"1", "Only the positive imaginary root lies above the real axis.", f)
	}
}

func w19MethBranchCut(o int) generator.Problem {
	f := map[string]string{}
	k := w19N(1, 40)
	w19I(f, "k", k)
	switch o {
	case 1:
		return w19P(o, fmt.Sprintf("How many branch points does log(z^%d) have on the Riemann sphere? (enter an integer)", k),
			"2", "The logarithm branches at zero and at infinity.", f)
	case 2:
		return w19P(o, fmt.Sprintf("How many branch points does sqrt(z^%d) have on the Riemann sphere? (enter an integer)", k),
			"2", "Zero and infinity are the branch points.", f)
	default:
		return w19P(o, fmt.Sprintf("How many branch cuts does the principal logarithm of z^%d need? (enter an integer)", k),
			"1", "One cut joins the two branch points.", f)
	}
}

func w19Fact(n int) int64 {
	out := int64(1)
	for i := 2; i <= n; i++ {
		out *= int64(i)
	}
	return out
}

func w19GCD64(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func w19SpecGammaHalf(o int) generator.Problem {
	f := map[string]string{}
	switch o {
	case 1:
		n := w19N(1, 10)
		w19I(f, "n", n)
		return w19P(o, fmt.Sprintf("What is Gamma(%d) for the integer %d? (enter an integer)", n+1, n+1),
			strconv.FormatInt(w19Fact(n), 10), "Gamma(n+1)=n!.", f)
	case 2:
		n := w19N(1, 8)
		w19I(f, "n", n)
		num := w19Fact(2 * n)
		den := int64(1)
		for i := 0; i < n; i++ {
			den *= 4
		}
		den *= w19Fact(n)
		g := w19GCD64(num, den)
		return w19P(o, fmt.Sprintf("The value Gamma(%d/2) equals sqrt(pi) times a rational number. What is that rational number's denominator in lowest terms? (enter an integer)", 2*n+1),
			strconv.FormatInt(den/g, 10), "The rational factor is (2n)!/(4^n n!), reduced.", f)
	default:
		n := w19N(1, 12)
		w19I(f, "n", n)
		return w19P(o, fmt.Sprintf("What is Gamma(%d)? (enter an integer)", n),
			strconv.FormatInt(w19Fact(n-1), 10), "Gamma(n)=(n-1)!.", f)
	}
}

func w19SpecBeta(o int) generator.Problem {
	m := w19N(2, 9)
	n := w19N(2, 9)
	f := map[string]string{}
	w19I(f, "m", m)
	w19I(f, "n", n)
	num := w19Fact(m-1) * w19Fact(n-1)
	den := w19Fact(m + n - 1)
	g := w19GCD64(num, den)
	switch o {
	case 1:
		return w19P(o, fmt.Sprintf("What is the numerator of B(%d,%d) before reduction, as a factorial product? (enter an integer)", m, n),
			strconv.FormatInt(num, 10), "The numerator is (m-1)!(n-1)!.", f)
	case 2:
		return w19P(o, fmt.Sprintf("What is the denominator of B(%d,%d) before reduction, as a factorial? (enter an integer)", m, n),
			strconv.FormatInt(den, 10), "The denominator is (m+n-1)!.", f)
	default:
		return w19P(o, fmt.Sprintf("What is the denominator of B(%d,%d) in lowest terms? (enter an integer)", m, n),
			strconv.FormatInt(den/g, 10), "Reduce the factorial ratio by its gcd.", f)
	}
}

func w19SpecZetaEven(o int) generator.Problem {
	f := map[string]string{}
	switch o {
	case 1:
		k := w19N(1, 8)
		m := w19N(1, 5)
		w19I(f, "k", k)
		w19I(f, "m", m)
		den := []int64{6, 90, 945, 9450, 93555, 638512875, 18243225, 325641566250}[k-1]
		return w19P(o, fmt.Sprintf("If zeta(%d) is pi^%d divided by an integer d, what is %d*d? (enter a number)", 2*k, 2*k, m),
			strconv.FormatInt(int64(m)*den, 10), "Multiply the Bernoulli denominator by the stated factor.", f)
	case 2:
		k := w19N(1, 9)
		w19I(f, "k", k)
		return w19P(o, fmt.Sprintf("If zeta(0) is written with denominator %d, what is its numerator? (enter an integer)", 2*k),
			strconv.Itoa(-k), "zeta(0) is -1/2, so the numerator is -k over the denominator 2k.", f)
	default:
		k := w19N(1, 15)
		w19I(f, "k", k)
		return w19P(o, fmt.Sprintf("What is zeta(-%d)? (enter an integer)", 2*k),
			"0", "Negative even integers are trivial zeros.", f)
	}
}

func w19SpecTheta(o int) generator.Problem {
	f := map[string]string{}
	switch o {
	case 1:
		k := w19N(1, 30)
		w19I(f, "k", k)
		return w19P(o, fmt.Sprintf("How many terms of the theta series have |n| <= %d? (enter an integer)", k),
			strconv.Itoa(2*k+1), "The indices run from -k to k.", f)
	case 2:
		n := w19N(1, 12)
		w19I(f, "n", n)
		return w19P(o, fmt.Sprintf("In sum q^(n^2), which exponent comes from n = %d? (enter an integer)", n),
			strconv.Itoa(n*n), "Square the index.", f)
	default:
		n := w19N(1, 200)
		w19I(f, "n", n)
		one, three := w19DivisorCounts(n)
		return w19P(o, fmt.Sprintf("What is the coefficient of q^%d in the theta series of the square lattice? (enter an integer)", n),
			strconv.Itoa(4*(one-three)), "The coefficient is r2(n).", f)
	}
}

func w19SpecBessel(o int) generator.Problem {
	f := map[string]string{}
	switch o {
	case 1:
		n := w19N(0, 15)
		w19I(f, "n", n)
		ans := 0
		if n == 0 {
			ans = 1
		}
		return w19P(o, fmt.Sprintf("What is J_%d(0)? (enter an integer)", n),
			strconv.Itoa(ans), "Only J_0 is nonzero at the origin.", f)
	case 2:
		n := w19N(1, 15)
		w19I(f, "n", n)
		return w19P(o, fmt.Sprintf("What is J_(-%d)(0)? (enter an integer)", n),
			"0", "Every negative order vanishes at zero as well.", f)
	default:
		n := w19N(1, 15)
		w19I(f, "n", n)
		den := int64(1)
		for i := 0; i < n; i++ {
			den *= 2
		}
		den *= w19Fact(n)
		return w19P(o, fmt.Sprintf("What is the denominator of the leading coefficient of J_%d(z)? (enter a number)", n),
			strconv.FormatInt(den, 10), "The leading term is z^n/(2^n n!).", f)
	}
}

func w19SpecAiry(o int) generator.Problem {
	f := map[string]string{}
	k := w19N(1, 40)
	w19I(f, "k", k)
	switch o {
	case 1:
		return w19P(o, fmt.Sprintf("What is the order of the differential equation y'' = x^%d y? (enter an integer)", k),
			"2", "The highest derivative is the second.", f)
	case 2:
		return w19P(o, fmt.Sprintf("What is the dimension of the solution space of y'' = x^%d y? (enter an integer)", k),
			"2", "A second-order linear equation has a two-dimensional solution space.", f)
	default:
		return w19P(o, fmt.Sprintf("Is the Wronskian of Ai and Bi zero at x = %d? (enter 1 for yes, 0 for no)", k),
			"0", "The Wronskian is a nonzero constant, so it never vanishes.", f)
	}
}

func w19SpecHypergeometric(o int) generator.Problem {
	f := map[string]string{}
	switch o {
	case 1:
		k := w19N(2, 9)
		w19I(f, "k", k)
		return w19P(o, fmt.Sprintf("What is 1F0(1;;1-1/%d)? (enter an integer)", k),
			strconv.Itoa(k), "1F0(1;;z)=1/(1-z), which is k at z=1-1/k.", f)
	case 2:
		p := w19N(1, 9)
		w19I(f, "p", p)
		return w19P(o, fmt.Sprintf("How many numerator parameters does a pFq function have when p = %d? (enter an integer)", p),
			strconv.Itoa(p), "The parameters before the semicolon are numerator parameters.", f)
	default:
		k := w19N(1, 9)
		w19I(f, "k", k)
		return w19P(o, fmt.Sprintf("What is the %d-th derivative of 0F0(;;z) at zero? (enter an integer)", k),
			"1", "Every derivative of exp at zero is one.", f)
	}
}

func w19SpecPolygamma(o int) generator.Problem {
	f := map[string]string{}
	switch o {
	case 1:
		n := w19N(1, 9)
		w19I(f, "n", n)
		return w19P(o, fmt.Sprintf("What is the denominator of psi(%d)-psi(1) as a harmonic number? (enter an integer)", n+1),
			strconv.Itoa(n), "The difference is H_n, whose denominator can be n.", f)
	case 2:
		k := w19N(1, 9)
		w19I(f, "k", k)
		return w19P(o, fmt.Sprintf("What is %d times (psi(%d)-psi(%d+1))? (enter an integer)", k, k, k),
			"-1", "The difference is -1/k, so k times it is -1.", f)
	default:
		n := w19N(0, 30)
		w19I(f, "n", n)
		return w19P(o, fmt.Sprintf("What is the order of the pole of psi at z=-%d? (enter an integer)", n),
			"1", "Every nonpositive integer is a simple pole.", f)
	}
}

func w19SpecErf(o int) generator.Problem {
	f := map[string]string{}
	k := w19N(1, 40)
	w19I(f, "k", k)
	switch o {
	case 1:
		return w19P(o, fmt.Sprintf("What is erf(0) for the Gaussian e^(-%d x^2)? (enter an integer)", k),
			"0", "The integral from zero to zero vanishes.", f)
	case 2:
		return w19P(o, fmt.Sprintf("What is the limit of erf(%dx) as x tends to infinity? (enter an integer)", k),
			"1", "The error function is normalised to total mass one.", f)
	default:
		return w19P(o, fmt.Sprintf("The integral of e^(-%d x^2) over the real line equals c*sqrt(pi/%d). What is c? (enter an integer)", k, k),
			"1", "The Gaussian integral is sqrt(pi/k).", f)
	}
}

func w19SpecPolylog(o int) generator.Problem {
	f := map[string]string{}
	switch o {
	case 1:
		k := w19N(1, 30)
		w19I(f, "k", k)
		return w19P(o, fmt.Sprintf("What is Li_%d(0)? (enter an integer)", k),
			"0", "Every polylogarithm vanishes at the origin.", f)
	case 2:
		k := w19N(2, 9)
		w19I(f, "k", k)
		return w19P(o, fmt.Sprintf("Li_1(1-1/%d) equals c*log %d. What is c? (enter an integer)", k, k),
			"1", "Li_1(z)=-log(1-z), so the value is log k.", f)
	default:
		k := w19N(1, 8)
		m := w19N(1, 5)
		w19I(f, "k", k)
		w19I(f, "m", m)
		den := []int64{6, 90, 945, 9450, 93555, 638512875, 18243225, 325641566250}[k-1]
		return w19P(o, fmt.Sprintf("If Li_%d(1) is pi^%d divided by an integer d, what is %d*d? (enter a number)", 2*k, 2*k, m),
			strconv.FormatInt(int64(m)*den, 10), "Li_s(1)=zeta(s), so multiply the zeta denominator.", f)
	}
}
