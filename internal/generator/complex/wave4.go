package complex

// Wave 4 complex analysis: twenty concepts. Every family asks for an integer
// invariant or a small exact value, and publishes the parameters in Facts so
// wave4_test.go can recompute the answer.
//
// Each knowledge point draws its own parameters and, where the mathematics
// offers more than one honest question, its own wording. The three ordinals
// therefore pose different problems instead of the same three problems under
// three labels, and one ordinal alone supplies well over ten distinct
// questions.

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
	return generator.Problem{Question: "What is 2 + 2? (enter a number)", Answer: "4", Explanation: "2 + 2 = 4.", KPOrdinal: kp}
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

// w4Point writes a complex point without decoration the reader does not need:
// 3, 4i, -4i, 3+4i or 3-4i.
func w4Point(re, im int) string {
	switch {
	case im == 0:
		return strconv.Itoa(re)
	case re == 0:
		if im > 0 {
			return fmt.Sprintf("%di", im)
		}
		return fmt.Sprintf("-%di", -im)
	default:
		return fmtComplex(re, im)
	}
}

// w4SinZ writes sin z or sin(a z), keeping a unit coefficient implicit.
func w4SinZ(a int) string {
	if a == 1 {
		return "sin z"
	}
	return fmt.Sprintf("sin(%d z)", a)
}

// w4ZPow writes z or z^a, keeping a unit exponent implicit.
func w4ZPow(a int) string {
	if a == 1 {
		return "z"
	}
	return fmt.Sprintf("z^{%d}", a)
}

// w4PiMul writes pi or a pi, keeping a unit coefficient implicit.
func w4PiMul(a int) string {
	if a == 1 {
		return "pi"
	}
	return fmt.Sprintf("%d pi", a)
}

func w4RootsUnity(kp int) generator.Problem {
	n := w4CR(3, 16)
	f := map[string]string{}
	w4CI(f, "n", n)
	switch kp {
	case 1:
		if rand.Intn(2) == 0 {
			return w4C(kp, fmt.Sprintf("How many distinct n-th roots of unity are there for n=%d? (enter a number)", n), strconv.Itoa(n),
				fmt.Sprintf("The equation z^{%d}=1 has exactly %d distinct roots in the complex plane.", n, n), f)
		}
		return w4C(kp, fmt.Sprintf("The %d-th roots of unity are the vertices of a regular polygon. How many vertices does that polygon have? (enter a number)", n), strconv.Itoa(n),
			fmt.Sprintf("Each of the %d roots is one vertex, and the roots are equally spaced around the unit circle.", n), f)
	case 2:
		if rand.Intn(2) == 0 {
			return w4C(kp, fmt.Sprintf("For n=%d, what is the sum of all n-th roots of unity? (enter a number)", n), "0",
				"The roots are the vertices of a regular polygon centred at the origin, so their vector sum is 0.", f)
		}
		product := 1
		if n%2 == 0 {
			product = -1
		}
		w4CI(f, "product", product)
		return w4C(kp, fmt.Sprintf("For n=%d, what is the product of all n-th roots of unity? (enter a number)", n), strconv.Itoa(product),
			fmt.Sprintf("The roots are the roots of z^{%d}-1, so their product is (-1)^{%d+1}=%d.", n, n, product), f)
	default:
		p := w4Phi(n)
		return w4C(kp, fmt.Sprintf("For n=%d, how many n-th roots of unity are primitive? (enter a number)", n), strconv.Itoa(p),
			fmt.Sprintf("A primitive root has order exactly %d. Their number is Euler's totient phi(%d)=%d.", n, n, p), f)
	}
}

// w4Dirs lists points on the eight compass directions of the complex plane and
// the principal argument of each, in degrees. Principal arguments lie in
// (-180, 180].
var w4Dirs = []struct{ re, im, deg int }{
	{1, 0, 0},
	{1, 1, 45},
	{0, 1, 90},
	{-1, 1, 135},
	{-1, 0, 180},
	{-1, -1, -135},
	{0, -1, -90},
	{1, -1, -45},
}

func w4Branches(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		d := w4Dirs[rand.Intn(len(w4Dirs))]
		r := w4CR(2, 9)
		re, im := r*d.re, r*d.im
		w4CI(f, "n", r)
		w4CI(f, "re", re)
		w4CI(f, "im", im)
		return w4C(kp, fmt.Sprintf("What is the principal argument of %s, in degrees? (enter a number)", w4Point(re, im)), strconv.Itoa(d.deg),
			fmt.Sprintf("The point %s lies at %d degrees from the positive real axis, inside the principal range (-180, 180], so its principal argument is %d degrees.", w4Point(re, im), d.deg, d.deg), f)
	case 2:
		if rand.Intn(2) == 0 {
			r := w4CR(2, 14)
			w4CI(f, "n", r)
			return w4C(kp, fmt.Sprintf("On the circle |z|=%d, how many points lie on the principal branch cut of Log? (enter a number)", r), "1",
				"The principal logarithm is cut along the negative real axis. That ray meets the circle at exactly one point, z=-"+strconv.Itoa(r)+".", f)
		}
		d := w4Dirs[rand.Intn(len(w4Dirs))]
		r := w4CR(2, 9)
		re, im := r*d.re, r*d.im
		onCut := 0
		if d.deg == 180 {
			onCut = 1
		}
		w4CI(f, "n", r)
		w4CI(f, "re", re)
		w4CI(f, "im", im)
		w4CI(f, "on_cut", onCut)
		place := "does not lie on"
		if onCut == 1 {
			place = "lies on"
		}
		return w4C(kp, fmt.Sprintf("For the principal Log, is the point %s on the branch cut? Enter 1 for yes and 0 for no.", w4Point(re, im)), strconv.Itoa(onCut),
			fmt.Sprintf("The cut is the negative real axis. The point %s %s that ray.", w4Point(re, im), place), f)
	default:
		r := w4CR(2, 14)
		k := w4CR(1, 5)
		w4CI(f, "n", r)
		w4CI(f, "k", k)
		return w4C(kp, fmt.Sprintf("For z=-%d+%di and z=-%d-%di, is the principal logarithm continuous across the cut? Enter 1 for yes and 0 for no.", r, k, r, k), "0",
			"The principal arguments approach +180 degrees from above the cut and -180 degrees from below, so the values jump by 2 pi i and the function is not continuous across it.", f)
	}
}

func w4ComplexPowers(kp int) generator.Problem {
	n := w4CR(1, 14)
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
	k := w4CR(1, 14)
	f := map[string]string{}
	w4CI(f, "k", k)
	switch kp {
	case 1:
		return w4C(kp, fmt.Sprintf("How many zeros does sin z have in the interval [-%s, %s]? (enter a number)", w4PiMul(k), w4PiMul(k)), strconv.Itoa(2*k+1),
			fmt.Sprintf("The zeros are the integer multiples of pi. There are %d of them in that interval.", 2*k+1), f)
	case 2:
		period := "2 pi"
		if k > 1 {
			period = fmt.Sprintf("2 pi/%d", k)
		}
		return w4C(kp, fmt.Sprintf("How many complete periods of %s fit in the interval [0, 2 pi]? (enter a number)", w4SinZ(k)), strconv.Itoa(k),
			fmt.Sprintf("The period of %s is %s, so exactly %d of them fit in an interval of length 2 pi.", w4SinZ(k), period, k), f)
	default:
		return w4C(kp, fmt.Sprintf("How many zeros does cos z have in the interval [0, %d pi]? (enter a number)", 2*k), strconv.Itoa(2*k),
			fmt.Sprintf("The zeros of cosine are pi/2 plus integer multiples of pi. There are %d in the interval.", 2*k), f)
	}
}

// w4HarmCase is one candidate for the harmonicity check. Every expression takes
// exactly one integer parameter, and the Laplacian quoted in why decides the
// answer.
type w4HarmCase struct {
	expr     string
	harmonic bool
	why      string
}

var w4HarmCases = []w4HarmCase{
	{"%d(x^2-y^2)", true, "u_xx=2a and u_yy=-2a, so the Laplacian is zero."},
	{"%d(x^2+y^2)", false, "u_xx=2a and u_yy=2a, so the Laplacian is 4a, which is not zero."},
	{"%dxy", true, "u_xx=0 and u_yy=0, so the Laplacian is zero."},
	{"%d(x^3-3xy^2)", true, "u_xx=6ax and u_yy=-6ax, so the Laplacian is zero."},
	{"%d(x^3+y^3)", false, "u_xx=6ax and u_yy=6ay, so the Laplacian is 6a(x+y), not identically zero."},
	{"%dx^2y", false, "u_xx=2ay, so the Laplacian is not identically zero."},
	{"%de^x sin y", true, "u_xx=a e^x sin y and u_yy=-a e^x sin y, so the Laplacian is zero."},
	{"%de^x cos y", true, "u_xx=a e^x cos y and u_yy=-a e^x cos y, so the Laplacian is zero."},
	{"%dx^2y^2", false, "u_xx=2ay^2 and u_yy=2ax^2, so the Laplacian is not identically zero."},
}

func w4HarmonicConjugate(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		a := w4CR(1, 14)
		w4CI(f, "a", a)
		return w4C(kp, fmt.Sprintf("For u=%d(x^2-y^2), what is the coefficient of xy in a harmonic conjugate v? (enter a number)", a), strconv.Itoa(2*a),
			fmt.Sprintf("The Cauchy-Riemann equations give v=2*%d xy plus a constant, so the coefficient is %d.", a, 2*a), f)
	case 2:
		a := w4CR(1, 12)
		c := w4HarmCases[rand.Intn(len(w4HarmCases))]
		expr := fmt.Sprintf(c.expr, a)
		w4CI(f, "a", a)
		if c.harmonic {
			w4CI(f, "harmonic", 1)
		} else {
			w4CI(f, "harmonic", 0)
		}
		answer := "0"
		if c.harmonic {
			answer = "1"
		}
		return w4C(kp, fmt.Sprintf("Is u=%s harmonic? Enter 1 for yes and 0 for no.", expr), answer,
			fmt.Sprintf("Write a=%d. Then %s", a, c.why), f)
	default:
		a := w4CR(1, 6)
		x0 := w4CR(1, 4)
		y0 := w4CR(1, 4)
		v := 2 * a * x0 * y0
		w4CI(f, "a", a)
		w4CI(f, "x0", x0)
		w4CI(f, "y0", y0)
		return w4C(kp, fmt.Sprintf("For u=%d(x^2-y^2), a harmonic conjugate is v=2*%d xy. What is v(%d,%d)? (enter a number)", a, a, x0, y0), strconv.Itoa(v),
			fmt.Sprintf("Substitute x=%d and y=%d: v=2*%d*%d*%d=%d.", x0, y0, a, x0, y0, v), f)
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

// w4AnalyticForms are entire nonconstant functions with one integer parameter.
var w4AnalyticForms = []string{
	"z^2+%d", "z^3+%dz", "e^z+%d", "sin z+%d", "cos z+%d", "z^4-%dz^2", "e^{%dz}",
}

// w4NonAnalyticForms are functions that fail the hypothesis of the open mapping
// theorem because they are not analytic.
var w4NonAnalyticForms = []string{"conjugate z", "|z|^2", "Re z", "Im z+%d"}

func w4OpenMapping(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		a := w4CR(1, 12)
		expr := fmt.Sprintf(w4AnalyticForms[rand.Intn(len(w4AnalyticForms))], a)
		w4CI(f, "a", a)
		w4CI(f, "applies", 1)
		return w4C(kp, fmt.Sprintf("For the nonconstant analytic function f(z)=%s, is the image of an open set open? Enter 1 for yes and 0 for no.", expr), "1",
			"The open mapping theorem says every nonconstant analytic function maps open sets to open sets.", f)
	case 2:
		a := w4CR(1, 9)
		if rand.Intn(2) == 0 {
			a = -a
		}
		w4CI(f, "a", a)
		return w4C(kp, fmt.Sprintf("For the constant function f(z)=%d, is the image of the open unit disk open? Enter 1 for yes and 0 for no.", a), "0",
			fmt.Sprintf("A constant map sends the disk to the single point %d, and a point is not open in the plane.", a), f)
	default:
		a := w4CR(1, 12)
		applies := 0
		forms := w4NonAnalyticForms
		why := "The function is not analytic, so the hypothesis of the open mapping theorem fails."
		if rand.Intn(2) == 0 {
			forms = w4AnalyticForms
			applies = 1
			why = "The function is nonconstant and analytic, so the theorem applies."
		}
		expr := fmt.Sprintf(forms[rand.Intn(len(forms))], a)
		w4CI(f, "a", a)
		w4CI(f, "applies", applies)
		return w4C(kp, fmt.Sprintf("Does the open mapping theorem apply to f(z)=%s on the unit disk? Enter 1 for yes and 0 for no.", expr), strconv.Itoa(applies), why, f)
	}
}

// w4SingCase is one function with an isolated singularity at 0. omitZero says
// whether 0 is missing from its image near 0.
type w4SingCase struct {
	expr     string
	omitZero bool
	why      string
}

var w4CasoratiCases = []w4SingCase{
	{"e^{%d/z}", true, "The exponential is never zero."},
	{"%de^{1/z}", true, "A nonzero multiple of the exponential is never zero."},
	{"e^{1/(%dz)}", true, "The exponential is never zero."},
	{"e^{%d/z}+1", false, "It attains zero where e^{%d/z}=-1."},
	{"e^{%d/z}-1", false, "It attains zero where e^{%d/z}=1."},
	{"e^{%d/z}+i", false, "It attains zero where e^{%d/z}=-i."},
	{"sin(%d/z)", false, "It attains zero wherever %d/z is an integer multiple of pi."},
	{"cos(%d/z)", false, "It attains zero wherever %d/z is pi/2 plus an integer multiple of pi."},
	{"sinh(%d/z)", false, "It attains zero wherever %d/z is an integer multiple of pi i."},
}

func w4Casorati(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		a := w4CR(1, 14)
		w4CI(f, "a", a)
		w4CI(f, "coeff", a)
		if rand.Intn(2) == 0 {
			return w4C(kp, fmt.Sprintf("The Laurent expansion of f(z)=e^{%d/z} about 0 has coefficient of z^{-1} equal to what? (enter a number)", a), strconv.Itoa(a),
				fmt.Sprintf("The expansion contains the term %d/z, so the coefficient of z^{-1} is %d.", a, a), f)
		}
		return w4C(kp, fmt.Sprintf("The Laurent expansion of f(z)=%de^{1/z} about 0 has coefficient of z^{-1} equal to what? (enter a number)", a), strconv.Itoa(a),
			fmt.Sprintf("The expansion contains the term %d/z, so the coefficient of z^{-1} is %d.", a, a), f)
	case 2:
		a := w4CR(1, 14)
		c := w4CasoratiCases[rand.Intn(len(w4CasoratiCases))]
		expr := fmt.Sprintf(c.expr, a)
		w4CI(f, "a", a)
		if c.omitZero {
			w4CI(f, "omit_zero", 1)
		} else {
			w4CI(f, "omit_zero", 0)
		}
		answer := "0"
		if c.omitZero {
			answer = "1"
		}
		return w4C(kp, fmt.Sprintf("Does f(z)=%s omit the value 0 near 0? Enter 1 for yes and 0 for no.", expr), answer, c.why, f)
	default:
		a := w4CR(1, 14)
		c := w4CasoratiCases[rand.Intn(len(w4CasoratiCases))]
		expr := fmt.Sprintf(c.expr, a)
		w4CI(f, "a", a)
		return w4C(kp, fmt.Sprintf("By Casorati-Weierstrass, is the image of a punctured neighbourhood of 0 under f(z)=%s dense in the plane? Enter 1 for yes and 0 for no.", expr), "1",
			"The function has an essential singularity at 0, so Casorati-Weierstrass says its image on every punctured neighbourhood is dense.", f)
	}
}

// w4PicardCases lists functions near an isolated singularity and how many
// values each omits there.
type w4PicardCase struct {
	expr    string
	omitted int
	why     string
}

var w4PicardCases = []w4PicardCase{
	{"e^{%d/z}", 1, "The exponential omits only 0."},
	{"%de^{1/z}", 1, "A nonzero multiple of the exponential omits only 0."},
	{"e^{1/(%dz)}", 1, "The exponential omits only 0."},
	{"e^{%d/z}+1", 1, "The shifted exponential omits only 1."},
	{"e^{%d/z}-1", 1, "The shifted exponential omits only -1."},
	{"e^{%d/z}+i", 1, "The shifted exponential omits only -i."},
	{"sin(%d/z)", 0, "The sine is surjective on the plane, so it omits no value."},
	{"cos(%d/z)", 0, "The cosine is surjective on the plane, so it omits no value."},
	{"sinh(%d/z)", 0, "The hyperbolic sine is surjective on the plane, so it omits no value."},
	{"cosh(%d/z)", 0, "The hyperbolic cosine is surjective on the plane, so it omits no value."},
}

func w4Picard(kp int) generator.Problem {
	a := w4CR(1, 14)
	c := w4PicardCases[rand.Intn(len(w4PicardCases))]
	expr := fmt.Sprintf(c.expr, a)
	f := map[string]string{}
	w4CI(f, "a", a)
	w4CI(f, "omitted", c.omitted)
	switch kp {
	case 1:
		return w4C(kp, fmt.Sprintf("Picard's great theorem says that near the essential singularity of f(z)=%s, the function omits at most how many values? (enter a number)", expr), "1",
			"At most one value can be omitted in every punctured neighbourhood of an essential singularity.", f)
	case 2:
		return w4C(kp, fmt.Sprintf("How many values does f(z)=%s omit near 0? (enter a number)", expr), strconv.Itoa(c.omitted), c.why, f)
	default:
		return w4C(kp, fmt.Sprintf("Can f(z)=%s omit two distinct values near 0? Enter 1 for yes and 0 for no.", expr), "0",
			"Picard's theorem allows at most one omitted value near an essential singularity, so two cannot be omitted.", f)
	}
}

func w4AnalyticContinuation(kp int) generator.Problem {
	a := w4CR(2, 15)
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
			fmt.Sprintf("The only obstruction on the circle is the pole at z=%d, so every other point of the circle is a regular point of the continuation.", a), f)
	}
}

func w4Monodromy(kp int) generator.Problem {
	m := w4CR(1, 14)
	n := w4CR(1, 14)
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
			fmt.Sprintf("One loop adds %d times 2 pi i, which is nonzero, so monodromy is nontrivial.", m), f)
	}
}

// w4RootExpr writes z^{1/n} or (z-a)^{1/n}, omitting a zero shift.
func w4RootExpr(a, n int) string {
	if a == 0 {
		return fmt.Sprintf("z^{1/%d}", n)
	}
	return fmt.Sprintf("(z-%d)^{1/%d}", a, n)
}

func w4RiemannSurface(kp int) generator.Problem {
	n := w4CR(2, 15)
	a := w4CR(0, 9)
	f := map[string]string{}
	w4CI(f, "n", n)
	w4CI(f, "a", a)
	switch kp {
	case 1:
		return w4C(kp, fmt.Sprintf("How many sheets does the Riemann surface of %s have? (enter a number)", w4RootExpr(a, n)), strconv.Itoa(n),
			fmt.Sprintf("Continuing %s once around its branch point multiplies it by a primitive %d-th root of unity, so %d sheets are needed.", w4RootExpr(a, n), n, n), f)
	case 2:
		return w4C(kp, fmt.Sprintf("For %s, at which real point is the branch point? (enter the coordinate)", w4RootExpr(a, n)), strconv.Itoa(a),
			fmt.Sprintf("The function %s is singular at the point where its argument vanishes, namely z=%d.", w4RootExpr(a, n), a), f)
	default:
		model := fmt.Sprintf("w^{%d}=z", n)
		where := "the origin"
		if a != 0 {
			model = fmt.Sprintf("w^{%d}=z-%d", n, a)
			where = fmt.Sprintf("z=%d", a)
		}
		return w4C(kp, fmt.Sprintf("For the local model %s near the branch point %s, the projection to z is how many-to-one away from the branch point? (enter a number)", model, where), strconv.Itoa(n),
			fmt.Sprintf("The local model is an %d-sheeted cover, so the projection is %d-to-one away from the branch point.", n, n), f)
	}
}

// w4EllCases lists nonconstant elliptic functions and the order of their pole in
// one fundamental parallelogram.
type w4EllCase struct {
	expr  string
	poles int
}

var w4EllCases = []w4EllCase{
	{"wp(z)+%d", 2},
	{"wp(z)^2+%d", 4},
	{"wp'(z)^2+%d", 6},
	{"wp(z)wp'(z)+%d", 5},
	{"%dwp(z)", 2},
	{"wp(z)^3+%d", 6},
}

func w4Elliptic(kp int) generator.Problem {
	a := w4CR(1, 12)
	c := w4EllCases[rand.Intn(len(w4EllCases))]
	expr := fmt.Sprintf(c.expr, a)
	f := map[string]string{}
	w4CI(f, "a", a)
	w4CI(f, "poles", c.poles)
	switch kp {
	case 1:
		return w4C(kp, fmt.Sprintf("How many independent periods does the nonconstant elliptic function %s have? (enter a number)", expr), "2",
			"An elliptic function is doubly periodic, so its period lattice has rank 2.", f)
	case 2:
		return w4C(kp, fmt.Sprintf("What is the sum of the residues of %s in one fundamental parallelogram? (enter a number)", expr), "0",
			"The sum of the residues is zero because the boundary integrals cancel in pairs.", f)
	default:
		return w4C(kp, fmt.Sprintf("Counting multiplicity, how many poles does %s have in one fundamental parallelogram? (enter a number)", expr), strconv.Itoa(c.poles),
			fmt.Sprintf("wp has a double pole at each lattice point and wp' has a triple pole, so this combination has a pole of order %d.", c.poles), f)
	}
}

// w4WeightForm writes a modular form such as E4, E6 or Delta and omits a unit
// exponent.
func w4WeightForm(base string, n int) string {
	if n == 1 {
		return base
	}
	return fmt.Sprintf("%s^{%d}", base, n)
}

func w4ModularForms(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		k := w4CR(1, 14)
		w4CI(f, "k", k)
		expr := w4WeightForm("E4", k)
		why := "E4 has weight 4."
		if k > 1 {
			why = fmt.Sprintf("Each E4 has weight 4, so the %d-th power has weight %d.", k, 4*k)
		}
		return w4C(kp, fmt.Sprintf("What is the weight of %s? (enter a number)", expr), strconv.Itoa(4*k), why, f)
	case 2:
		p := w4CR(0, 4)
		q := w4CR(1, 12)
		w4CI(f, "p", p)
		w4CI(f, "q", q)
		expr := w4WeightForm("E6", q)
		if p > 0 {
			expr = w4WeightForm("E4", p) + " " + w4WeightForm("E6", q)
		}
		weight := 4*p + 6*q
		return w4C(kp, fmt.Sprintf("What is the weight of %s? (enter a number)", expr), strconv.Itoa(weight),
			fmt.Sprintf("E4 has weight 4 and E6 has weight 6, so the product has weight 4*%d+6*%d=%d.", p, q, weight), f)
	default:
		k := w4CR(1, 12)
		m := w4CR(0, 3)
		w4CI(f, "k", k)
		w4CI(f, "m", m)
		expr := w4WeightForm("Delta", k)
		if m > 0 {
			expr = w4WeightForm("Delta", k) + " " + w4WeightForm("E4", m)
		}
		weight := 12*k + 4*m
		return w4C(kp, fmt.Sprintf("What is the weight of %s? (enter a number)", expr), strconv.Itoa(weight),
			fmt.Sprintf("Delta has weight 12 and E4 has weight 4, so the product has weight 12*%d+4*%d=%d.", k, m, weight), f)
	}
}

func w4Gamma(kp int) generator.Problem {
	n := w4CR(1, 12)
	k := w4CR(1, 14)
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
	a := w4CR(1, 12)
	f := map[string]string{}
	w4CI(f, "a", a)
	switch kp {
	case 1:
		m := w4CR(1, 3)
		w4CI(f, "m", m)
		expr := w4SinZ(a)
		if m > 1 {
			expr = fmt.Sprintf("%s^%d", w4SinZ(a), m)
		}
		return w4C(kp, fmt.Sprintf("What is the order of the zero of %s at z=0? (enter a number)", expr), strconv.Itoa(m),
			fmt.Sprintf("%s has a simple zero at 0, so its %d-th power has a zero of order %d.", w4SinZ(a), m, m), f)
	case 2:
		k := w4CR(1, 4)
		w4CI(f, "k", k)
		count := 2*k*a - 1
		zeroWord := "zeros"
		if count == 1 {
			zeroWord = "zero"
		}
		disk := "pi"
		if k > 1 {
			disk = fmt.Sprintf("%d pi", k)
		}
		return w4C(kp, fmt.Sprintf("How many zeros does %s have in |z|<%s? (enter a number)", w4SinZ(a), disk), strconv.Itoa(count),
			fmt.Sprintf("The zeros are m pi/%d with |m|<%d, so there are %d %s.", a, k*a, count, zeroWord), f)
	default:
		return w4C(kp, fmt.Sprintf("For %s, the elementary factor for the zero at 0 is z. Enter 1 for yes and 0 for no.", w4SinZ(a)), "1",
			"The first Weierstrass elementary factor for a zero at 0 is z.", f)
	}
}

func w4MittagLeffler(kp int) generator.Problem {
	k := w4CR(1, 14)
	f := map[string]string{}
	w4CI(f, "k", k)
	switch kp {
	case 1:
		return w4C(kp, fmt.Sprintf("How many poles does csc z have in [-%s, %s]? (enter a number)", w4PiMul(k), w4PiMul(k)), strconv.Itoa(2*k+1),
			fmt.Sprintf("The poles are the integer multiples of pi. There are %d of them in the interval.", 2*k+1), f)
	case 2:
		res := 1
		if k%2 == 1 {
			res = -1
		}
		return w4C(kp, fmt.Sprintf("What is the residue of csc z at z=%s? (enter a number)", w4PiMul(k)), strconv.Itoa(res),
			fmt.Sprintf("The residue of csc z at an integer multiple k pi is (-1)^k. Here k=%d, so the residue is %d.", k, res), f)
	default:
		return w4C(kp, fmt.Sprintf("For the doubled interval [-%s, %s], how many poles does csc z have? (enter a number)", w4PiMul(2*k), w4PiMul(2*k)), strconv.Itoa(4*k+1),
			fmt.Sprintf("The integer multiples of pi in the interval give %d poles.", 4*k+1), f)
	}
}

func w4SchwarzChristoffel(kp int) generator.Problem {
	n := w4CR(3, 15)
	m := w4CR(2, 13)
	f := map[string]string{}
	w4CI(f, "n", n)
	w4CI(f, "m", m)
	switch kp {
	case 1:
		return w4C(kp, fmt.Sprintf("A Schwarz-Christoffel map sends the upper half-plane to a polygon with %d sides. What is the sum of its interior angles in units of pi? (enter a number)", n), strconv.Itoa(n-2),
			fmt.Sprintf("A polygon with %d sides has interior angles summing to (%d-2) pi, which is %d units of pi.", n, n, n-2), f)
	case 2:
		return w4C(kp, fmt.Sprintf("For an interior angle pi/%d at a vertex, the Schwarz-Christoffel exponent is alpha/pi-1=(1-%d)/%d. What is its denominator? (enter a number)", m, m, m), strconv.Itoa(m),
			fmt.Sprintf("The exponent is 1/%d-1=(1-%d)/%d, and that fraction is in lowest terms, so its denominator is %d.", m, m, m, m), f)
	default:
		return w4C(kp, fmt.Sprintf("For a polygon with %d vertices, what is the sum of the Schwarz-Christoffel exponents alpha/pi-1? (enter a number)", n), "-2",
			"The sum of the exterior angles is 2 pi, so the sum of the exponents is -2.", f)
	}
}

func w4Jensen(kp int) generator.Problem {
	a := w4CR(1, 12)
	f := map[string]string{}
	w4CI(f, "a", a)
	switch kp {
	case 1:
		zeroWord := "zeros"
		if a == 1 {
			zeroWord = "zero"
		}
		return w4C(kp, fmt.Sprintf("How many zeros does %s(z-1) have inside |z|<1/2? (enter a number)", w4ZPow(a)), strconv.Itoa(a),
			fmt.Sprintf("The factor %s contributes %d %s at 0, while z-1 has no zero inside the small disk.", w4ZPow(a), a, zeroWord), f)
	case 2:
		r := w4CR(1, 14)
		w4CI(f, "r", r)
		return w4C(kp, fmt.Sprintf("If an analytic function has no zeros in |z|<%d, does Jensen's formula reduce to log|f(0)| = the average of log|f(%d e^{i theta})|? Enter 1 for yes and 0 for no.", r, r), "1",
			"With no zeros the sum over zeros is empty, so the formula reduces to the mean-value identity.", f)
	default:
		b := w4CR(2, 6)
		w4CI(f, "b", b)
		return w4C(kp, fmt.Sprintf("How many zeros does %s(z-1) have inside |z|<%d? (enter a number)", w4ZPow(a), b), strconv.Itoa(a+1),
			fmt.Sprintf("The factor %s contributes %d zeros and z-1 contributes one more, all inside the disk of radius %d.", w4ZPow(a), a, b), f)
	}
}
