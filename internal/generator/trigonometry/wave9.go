package trigonometry

// Wave 9 trigonometry: the twenty-two concepts that take the domain to one
// hundred and twenty. Every answer is an integer, every operand is published
// in Facts so wave9_test.go can recompute the answer independently, and every
// ordinal draws wide enough ranges to keep five distinct questions in twelve
// draws.

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave9(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"trig.eq.factoring":            "factoring",
		"trig.eq.reciprocal":           "reciprocal",
		"trig.eq.mixed_angles":         "mixed_angles",
		"trig.eq.extraneous":           "extraneous",
		"trig.eq.absolute":             "absolute",
		"trig.func.higher_derivatives": "higher_derivatives",
		"trig.func.taylor":             "taylor",
		"trig.func.power_integrals":    "power_integrals",
		"trig.graph.damping":           "damping",
		"trig.graph.envelope":          "envelope",
		"trig.graph.beats":             "beats",
		"trig.hyperbolic.catenary":     "catenary",
		"trig.hyperbolic.gudermannian": "gudermannian",
		"trig.hyperbolic.addition":     "hyperbolic_addition",
		"trig.hyperbolic.integral":     "hyperbolic_integral",
		"trig.spherical.right":         "spherical_right",
		"trig.spherical.polar":         "spherical_polar",
		"trig.spherical.polygon":       "spherical_polygon",
		"trig.basics.gradian":          "gradian",
		"trig.basics.signed":           "signed",
		"trig.ident.conditional":       "conditional",
		"trig.ident.tangent_triple":    "tangent_triple",
	} {
		reg.Register(id, &w9TrigGen{kind: kind})
	}
}

type w9TrigGen struct{ kind string }

func (g *w9TrigGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "factoring":
		return w9Factoring(kp)
	case "reciprocal":
		return w9Reciprocal(kp)
	case "mixed_angles":
		return w9MixedAngles(kp)
	case "extraneous":
		return w9Extraneous(kp)
	case "absolute":
		return w9Absolute(kp)
	case "higher_derivatives":
		return w9HigherDerivatives(kp)
	case "taylor":
		return w9Taylor(kp)
	case "power_integrals":
		return w9PowerIntegrals(kp)
	case "damping":
		return w9Damping(kp)
	case "envelope":
		return w9Envelope(kp)
	case "beats":
		return w9Beats(kp)
	case "catenary":
		return w9Catenary(kp)
	case "gudermannian":
		return w9Gudermannian(kp)
	case "hyperbolic_addition":
		return w9HyperbolicAddition(kp)
	case "hyperbolic_integral":
		return w9HyperbolicIntegral(kp)
	case "spherical_right":
		return w9SphericalRight(kp)
	case "spherical_polar":
		return w9SphericalPolar(kp)
	case "spherical_polygon":
		return w9SphericalPolygon(kp)
	case "gradian":
		return w9Gradian(kp)
	case "signed":
		return w9Signed(kp)
	case "conditional":
		return w9Conditional(kp)
	case "tangent_triple":
		return w9TangentTriple(kp)
	}
	return generator.Problem{Question: "What is 2 + 2? (enter a number)", Answer: "4", Explanation: "2 + 2 = 4."}
}

func w9P(o int, q, a, ex string, f map[string]string) generator.Problem {
	f["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: f, KPOrdinal: o}
}

func w9I(f map[string]string, k string, v int) { f[k] = strconv.Itoa(v) }

func w9R(lo, hi int) int {
	if hi <= lo {
		return lo
	}
	return lo + rand.Intn(hi-lo+1)
}

func w9Fact(n int) int {
	out := 1
	for i := 2; i <= n; i++ {
		out *= i
	}
	return out
}

func w9Degrees(rad float64) int { return int(math.Round(rad * 180 / math.Pi)) }

func w9Factoring(kp int) generator.Problem {
	switch kp {
	case 1:
		n := w9R(3, 12)
		f := map[string]string{}
		w9I(f, "N", n)
		return w9P(kp,
			fmt.Sprintf("The equation (sin x)(cos x) = 0 is solved on [0, %d degrees). How many solutions are there? (enter a number)", 360*n),
			strconv.Itoa(4*n),
			fmt.Sprintf("sin x = 0 gives %d solutions and cos x = 0 gives %d more, for %d in total.", 2*n, 2*n, 4*n), f)
	case 2:
		n := w9R(5, 20)
		f := map[string]string{}
		w9I(f, "N", n)
		return w9P(kp,
			fmt.Sprintf("The equation (tan x)(sin x) = 0 is solved on [0, %d degrees). How many solutions are there? (enter a number)", 360*n),
			strconv.Itoa(2*n),
			fmt.Sprintf("Both factors vanish at the same points, 0 and 180 degrees in each turn, so there are %d solutions.", 2*n), f)
	default:
		n := w9R(5, 20)
		f := map[string]string{}
		w9I(f, "N", n)
		return w9P(kp,
			fmt.Sprintf("The equation (sin x - 1)(cos x - 1) = 0 is solved on [0, %d degrees). How many solutions are there? (enter a number)", 360*n),
			strconv.Itoa(n),
			fmt.Sprintf("sin x = 1 and cos x = 1 coincide at the start of each turn, so the %d turns give %d solutions.", n, n), f)
	}
}

func w9Reciprocal(kp int) generator.Problem {
	triples := [][3]int{{3, 4, 5}, {5, 12, 13}, {8, 15, 17}, {7, 24, 25}, {9, 40, 41}, {20, 21, 29}}
	t := triples[rand.Intn(len(triples))]
	f := map[string]string{}
	w9I(f, "a", t[0])
	w9I(f, "b", t[1])
	w9I(f, "h", t[2])
	switch kp {
	case 1:
		return w9P(kp,
			fmt.Sprintf("If sec(theta) = %d/%d, what is the numerator of cos(theta) in lowest terms? (enter a number)", t[2], t[0]),
			strconv.Itoa(t[0]),
			fmt.Sprintf("cos(theta) = 1/sec(theta) = %d/%d, so the numerator is %d.", t[0], t[2], t[0]), f)
	case 2:
		return w9P(kp,
			fmt.Sprintf("If csc(theta) = %d/%d, what is the numerator of sin(theta) in lowest terms? (enter a number)", t[2], t[1]),
			strconv.Itoa(t[1]),
			fmt.Sprintf("sin(theta) = 1/csc(theta) = %d/%d, so the numerator is %d.", t[1], t[2], t[1]), f)
	default:
		return w9P(kp,
			fmt.Sprintf("If cot(theta) = %d/%d, what is the numerator of tan(theta) in lowest terms? (enter a number)", t[0], t[1]),
			strconv.Itoa(t[1]),
			fmt.Sprintf("tan(theta) = 1/cot(theta) = %d/%d, so the numerator is %d.", t[1], t[0], t[1]), f)
	}
}

func w9MixedAngles(kp int) generator.Problem {
	m := w9R(2, 7)
	n := w9R(3, 9)
	f := map[string]string{}
	w9I(f, "m", m)
	w9I(f, "N", n)
	switch kp {
	case 1:
		return w9P(kp,
			fmt.Sprintf("How many solutions does sin(%dx) = 0 have on [0, %d degrees)? (enter a number)", m, 360*n),
			strconv.Itoa(2*m*n),
			fmt.Sprintf("sin(%dx) = 0 when %dx = 180k, so x = 180k/%d. There are %d such points.", m, m, m, 2*m*n), f)
	case 2:
		return w9P(kp,
			fmt.Sprintf("How many solutions does cos(%dx) = 0 have on [0, %d degrees)? (enter a number)", m, 360*n),
			strconv.Itoa(2*m*n),
			fmt.Sprintf("cos(%dx) = 0 when %dx = 90 + 180k, which gives %d solutions in the interval.", m, m, 2*m*n), f)
	default:
		return w9P(kp,
			fmt.Sprintf("How many solutions does sin(%dx) = 1 have on [0, %d degrees)? (enter a number)", m, 360*n),
			strconv.Itoa(m*n),
			fmt.Sprintf("sin(%dx) = 1 when %dx = 90 + 360k, which gives %d solutions.", m, m, m*n), f)
	}
}

func w9Extraneous(kp int) generator.Problem {
	n := w9R(3, 12)
	f := map[string]string{}
	w9I(f, "N", n)
	fn := "sin"
	if kp == 2 {
		fn = "cos"
	}
	if kp == 3 {
		fn = "tan"
	}
	f["fn"] = fn
	return w9P(kp,
		fmt.Sprintf("Squaring %s x = c, with 0 < c < 1, gives %s^2 x = c^2 on [0, %d degrees). The original equation has %d solutions. How many extraneous solutions does squaring add? (enter a number)", fn, fn, 360*n, 2*n),
		strconv.Itoa(2*n),
		fmt.Sprintf("The squared equation also admits %s x = -c, which doubles the solutions to %d; %d of those are extraneous.", fn, 4*n, 2*n), f)
}

func w9Absolute(kp int) generator.Problem {
	n := w9R(3, 15)
	f := map[string]string{}
	w9I(f, "N", n)
	switch kp {
	case 1:
		return w9P(kp,
			fmt.Sprintf("How many solutions does |sin x| = 1/2 have on [0, %d degrees)? (enter a number)", 360*n),
			strconv.Itoa(4*n),
			fmt.Sprintf("Each turn contributes four solutions, so %d turns give %d.", n, 4*n), f)
	case 2:
		return w9P(kp,
			fmt.Sprintf("How many solutions does |sin x| = 1 have on [0, %d degrees)? (enter a number)", 360*n),
			strconv.Itoa(2*n),
			fmt.Sprintf("|sin x| = 1 at the top and bottom of each turn, so %d turns give %d.", n, 2*n), f)
	default:
		return w9P(kp,
			fmt.Sprintf("How many solutions does |tan x| = 1 have on [0, %d degrees)? (enter a number)", 180*n),
			strconv.Itoa(2*n),
			fmt.Sprintf("tan x = 1 and tan x = -1 each occur once per 180 degrees, so the interval holds %d.", 2*n), f)
	}
}

func w9HigherDerivatives(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		a := w9R(2, 7)
		w9I(f, "a", a)
		ans := -a * a * a
		return w9P(kp,
			fmt.Sprintf("For f(x) = sin(%dx), the third derivative at x = 0 is -a^3 with a = %d. Give f'''(0). (enter a number)", a, a),
			strconv.Itoa(ans),
			fmt.Sprintf("Differentiating three times gives -%d^3 = %d.", a, ans), f)
	case 2:
		a := w9R(1, 6)
		w9I(f, "a", a)
		ans := a * a * a * a
		return w9P(kp,
			fmt.Sprintf("For f(x) = cos(%dx), the fourth derivative at x = 0 is a^4 with a = %d. Give f''''(0). (enter a number)", a, a),
			strconv.Itoa(ans),
			fmt.Sprintf("The fourth derivative of cos(ax) is a^4 cos(ax), which is %d at zero.", ans), f)
	default:
		a := w9R(2, 12)
		w9I(f, "a", a)
		ans := 2 * a
		return w9P(kp,
			fmt.Sprintf("For f(x) = x sin(%dx), the second derivative at x = 0 is 2a with a = %d. Give f''(0). (enter a number)", a, a),
			strconv.Itoa(ans),
			fmt.Sprintf("The product rule and the chain rule leave 2a, which is %d.", ans), f)
	}
}

func w9Taylor(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		n := w9R(1, 6)
		w9I(f, "n", n)
		ans := w9Fact(2*n + 1)
		return w9P(kp,
			fmt.Sprintf("The Taylor series of sin x has term x^(2n+1)/(2n+1)!. Give the denominator for n = %d. (enter a number)", n),
			strconv.Itoa(ans),
			fmt.Sprintf("(2n+1)! = %d! = %d.", 2*n+1, ans), f)
	case 2:
		n := w9R(2, 7)
		w9I(f, "n", n)
		ans := w9Fact(2 * n)
		return w9P(kp,
			fmt.Sprintf("The Taylor series of cos x has term x^(2n)/(2n)!. Give the denominator for n = %d. (enter a number)", n),
			strconv.Itoa(ans),
			fmt.Sprintf("(2n)! = %d! = %d.", 2*n, ans), f)
	default:
		a := w9R(2, 8)
		w9I(f, "a", a)
		ans := a * a * a
		return w9P(kp,
			fmt.Sprintf("The coefficient of x^3 in sin(%dx) is -a^3/6 with a = %d. Give the numerator of its magnitude. (enter a number)", a, a),
			strconv.Itoa(ans),
			fmt.Sprintf("The magnitude is %d^3/6, so the numerator is %d.", a, ans), f)
	}
}

func w9DoubleFact(n int) int {
	out := 1
	for i := n; i > 0; i -= 2 {
		out *= i
	}
	return out
}

func w9PowerIntegrals(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		n := 1 + 2*rand.Intn(6)
		w9I(f, "n", n)
		ans := w9DoubleFact(n - 1)
		return w9P(kp,
			fmt.Sprintf("For odd n, the integral from 0 to pi/2 of sin^n x dx has denominator (n-1)!!. Give the denominator for n = %d. (enter a number)", n),
			strconv.Itoa(ans),
			fmt.Sprintf("(%d-1)!! = %d.", n, ans), f)
	case 2:
		n := 2 + 2*rand.Intn(6)
		w9I(f, "n", n)
		ans := w9DoubleFact(n)
		return w9P(kp,
			fmt.Sprintf("For even n, the integral from 0 to pi/2 of sin^n x dx equals ((n-1)!!/n!!)(pi/2). Give the denominator n!! for n = %d. (enter a number)", n),
			strconv.Itoa(ans),
			fmt.Sprintf("%d!! = %d.", n, ans), f)
	default:
		n := 4 + 2*rand.Intn(6)
		w9I(f, "n", n)
		ans := w9DoubleFact(n - 1)
		return w9P(kp,
			fmt.Sprintf("With I_n the integral from 0 to pi/2 of sin^n x dx, I_n = ((n-1)/n) I_(n-2). Give the numerator of the rational factor of I_%d after reducing the product of all factors. (enter a number)", n),
			strconv.Itoa(ans),
			fmt.Sprintf("The rational factor of I_%d is %d/%d, so the numerator is %d.", n, ans, w9DoubleFact(n), ans), f)
	}
}

func w9Damping(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		h := w9R(2, 5)
		m := w9R(2, 9)
		w9I(f, "h", h)
		w9I(f, "m", m)
		return w9P(kp,
			fmt.Sprintf("A damped wave loses half its amplitude every %d seconds. How many half-lives pass in %d seconds? (enter a number)", h, h*m),
			strconv.Itoa(m),
			fmt.Sprintf("%d seconds divided by the %d second half-life is %d.", h*m, h, m), f)
	case 2:
		m := w9R(2, 8)
		w9I(f, "m", m)
		ans := 1 << m
		return w9P(kp,
			fmt.Sprintf("After m half-lives the amplitude is A/2^m. Give the denominator when m = %d. (enter a number)", m),
			strconv.Itoa(ans),
			fmt.Sprintf("2^%d = %d.", m, ans), f)
	default:
		k := w9R(1, 6)
		t := w9R(2, 9)
		w9I(f, "k", k)
		w9I(f, "t", t)
		ans := k * t
		return w9P(kp,
			fmt.Sprintf("A damped envelope is e^(-kt). Give the exponent at t = %d when k = %d. (enter a number)", t, k),
			strconv.Itoa(ans),
			fmt.Sprintf("kt = %d times %d = %d.", k, t, ans), f)
	}
}

func w9Envelope(kp int) generator.Problem {
	a := w9R(3, 12)
	b := w9R(1, 7)
	f := map[string]string{}
	w9I(f, "a", a)
	w9I(f, "b", b)
	switch kp {
	case 1:
		ans := a + b
		return w9P(kp,
			fmt.Sprintf("A wave has envelope A(t) = %d + %d cos t. Give its maximum. (enter a number)", a, b),
			strconv.Itoa(ans),
			fmt.Sprintf("The maximum of cos t is 1, so the maximum is %d + %d = %d.", a, b, ans), f)
	case 2:
		ans := a - b
		return w9P(kp,
			fmt.Sprintf("A wave has envelope A(t) = %d + %d cos t. Give its minimum. (enter a number)", a, b),
			strconv.Itoa(ans),
			fmt.Sprintf("The minimum of cos t is -1, so the minimum is %d - %d = %d.", a, b, ans), f)
	default:
		ans := 2 * b
		return w9P(kp,
			fmt.Sprintf("A wave has envelope A(t) = %d + %d cos t. Give the spread between its maximum and minimum. (enter a number)", a, b),
			strconv.Itoa(ans),
			fmt.Sprintf("The spread is (a+b) - (a-b) = 2b = %d.", ans), f)
	}
}

func w9Beats(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		f1 := w9R(200, 600)
		d := w9R(2, 12)
		w9I(f, "f1", f1)
		w9I(f, "f2", f1+d)
		return w9P(kp,
			fmt.Sprintf("Two tones of %d Hz and %d Hz sound together. Give the beat frequency in Hz. (enter a number)", f1, f1+d),
			strconv.Itoa(d),
			fmt.Sprintf("The beat frequency is the difference, %d Hz.", d), f)
	case 2:
		rate := w9R(2, 12)
		secs := w9R(2, 10)
		w9I(f, "f", rate)
		w9I(f, "m", secs)
		ans := rate * secs
		return w9P(kp,
			fmt.Sprintf("A beat frequency of %d Hz lasts for %d seconds. How many beats occur? (enter a number)", rate, secs),
			strconv.Itoa(ans),
			fmt.Sprintf("%d beats per second times %d seconds is %d.", rate, secs, ans), f)
	default:
		base := w9R(300, 500)
		d := w9R(2, 9)
		w9I(f, "f1", base)
		w9I(f, "f2", base+d)
		ans := 10 * d
		return w9P(kp,
			fmt.Sprintf("Two tones of %d Hz and %d Hz sound together for 10 seconds. How many beats occur? (enter a number)", base, base+d),
			strconv.Itoa(ans),
			fmt.Sprintf("The beat frequency is %d Hz, so 10 seconds hold %d beats.", d, ans), f)
	}
}

func w9Catenary(kp int) generator.Problem {
	a := w9R(2, 9)
	f := map[string]string{}
	w9I(f, "a", a)
	switch kp {
	case 1:
		ans := int(math.Round(float64(a) * math.Cosh(1)))
		return w9P(kp,
			fmt.Sprintf("A chain hangs as y = %d cosh(x/%d). Give its height at x = %d, rounded to the nearest unit. (enter a number)", a, a, a),
			strconv.Itoa(ans),
			fmt.Sprintf("y = %d cosh(1) = %.3f, which rounds to %d.", a, float64(a)*math.Cosh(1), ans), f)
	case 2:
		ans := int(math.Round(float64(a) * math.Cosh(0.5)))
		return w9P(kp,
			fmt.Sprintf("A chain hangs as y = %d cosh(x/%d). Give its height at x = %d/2, rounded to the nearest unit. (enter a number)", a, a, a),
			strconv.Itoa(ans),
			fmt.Sprintf("y = %d cosh(0.5) = %.3f, which rounds to %d.", a, float64(a)*math.Cosh(0.5), ans), f)
	default:
		ans := int(math.Round(2 * float64(a) * math.Sinh(1)))
		return w9P(kp,
			fmt.Sprintf("The arc length of y = %d cosh(x/%d) from x = -%d to x = %d is 2a sinh(1). Give the length rounded to the nearest unit. (enter a number)", a, a, a, a),
			strconv.Itoa(ans),
			fmt.Sprintf("2 times %d times sinh(1) is %.3f, which rounds to %d.", a, 2*float64(a)*math.Sinh(1), ans), f)
	}
}

func w9Gudermannian(kp int) generator.Problem {
	x := w9R(1, 6)
	f := map[string]string{}
	w9I(f, "x", x)
	switch kp {
	case 1:
		ans := int(math.Round(math.Sinh(float64(x))))
		return w9P(kp,
			fmt.Sprintf("The Gudermannian angle satisfies tan(gd x) = sinh x. Give sinh %d rounded to the nearest integer. (enter a number)", x),
			strconv.Itoa(ans),
			fmt.Sprintf("sinh %d is about %.3f, which rounds to %d.", x, math.Sinh(float64(x)), ans), f)
	case 2:
		gd := 2*math.Atan(math.Exp(float64(x))) - math.Pi/2
		ans := w9Degrees(gd)
		return w9P(kp,
			fmt.Sprintf("Give the Gudermannian angle gd(%d) in degrees, rounded to the nearest degree. (enter a number)", x),
			strconv.Itoa(ans),
			fmt.Sprintf("gd(%d) is about %d degrees.", x, ans), f)
	default:
		k := w9R(2, 12)
		w9I(f, "x", k)
		return w9P(kp,
			fmt.Sprintf("If sinh x = %d, what is tan(gd x)? (enter a number)", k),
			strconv.Itoa(k),
			fmt.Sprintf("tan(gd x) equals sinh x, so the answer is %d.", k), f)
	}
}

func w9HyperbolicAddition(kp int) generator.Problem {
	triples := [][2]int{{3, 5}, {5, 13}, {8, 17}, {7, 25}, {9, 41}, {20, 29}}
	if kp == 3 {
		one := triples[rand.Intn(len(triples))]
		two := triples[rand.Intn(len(triples))]
		f := map[string]string{}
		w9I(f, "s1", one[0])
		w9I(f, "c1", one[1])
		w9I(f, "s2", two[0])
		w9I(f, "c2", two[1])
		ans := one[0]*two[1] + one[1]*two[0]
		return w9P(kp,
			fmt.Sprintf("If sinh x = %d, cosh x = %d, sinh y = %d and cosh y = %d, give sinh(x+y). (enter a number)", one[0], one[1], two[0], two[1]),
			strconv.Itoa(ans),
			fmt.Sprintf("sinh(x+y) = sinh x cosh y + cosh x sinh y = %d.", ans), f)
	}
	t := triples[rand.Intn(len(triples))]
	f := map[string]string{}
	w9I(f, "s", t[0])
	w9I(f, "c", t[1])
	if kp == 1 {
		ans := 2 * t[0] * t[1]
		return w9P(kp,
			fmt.Sprintf("If sinh x = %d and cosh x = %d, give sinh(2x). (enter a number)", t[0], t[1]),
			strconv.Itoa(ans),
			fmt.Sprintf("sinh(2x) = 2 sinh x cosh x = %d.", ans), f)
	}
	ans := t[0]*t[0] + t[1]*t[1]
	return w9P(kp,
		fmt.Sprintf("If sinh x = %d and cosh x = %d, give cosh(2x). (enter a number)", t[0], t[1]),
		strconv.Itoa(ans),
		fmt.Sprintf("cosh(2x) = cosh^2 x + sinh^2 x = %d.", ans), f)
}

func w9HyperbolicIntegral(kp int) generator.Problem {
	a := w9R(1, 6)
	f := map[string]string{}
	w9I(f, "a", a)
	switch kp {
	case 1:
		ans := int(math.Round(math.Sinh(float64(a))))
		return w9P(kp,
			fmt.Sprintf("The integral from 0 to %d of cosh x dx equals sinh %d. Give its value rounded to the nearest integer. (enter a number)", a, a),
			strconv.Itoa(ans),
			fmt.Sprintf("sinh %d is about %.3f, which rounds to %d.", a, math.Sinh(float64(a)), ans), f)
	case 2:
		ans := int(math.Round(math.Cosh(float64(a)) - 1))
		return w9P(kp,
			fmt.Sprintf("The integral from 0 to %d of sinh x dx equals cosh %d - 1. Give its value rounded to the nearest integer. (enter a number)", a, a),
			strconv.Itoa(ans),
			fmt.Sprintf("cosh %d - 1 is about %.3f, which rounds to %d.", a, math.Cosh(float64(a))-1, ans), f)
	default:
		k := w9R(2, 12)
		w9I(f, "a", k)
		return w9P(kp,
			fmt.Sprintf("The integral of sinh(%dx) dx is cosh(%dx)/%d + C. Give the denominator. (enter a number)", k, k, k),
			strconv.Itoa(k),
			fmt.Sprintf("The chain rule divides by %d, so the denominator is %d.", k, k), f)
	}
}

func w9SphericalRight(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		a := w9R(40, 80)
		b := w9R(40, 80)
		w9I(f, "A", a)
		w9I(f, "B", b)
		ans := a + b
		return w9P(kp,
			fmt.Sprintf("A right spherical triangle has angles A = %d degrees and B = %d degrees besides the right angle. Give A+B in degrees. (enter a number)", a, b),
			strconv.Itoa(ans),
			fmt.Sprintf("%d + %d = %d.", a, b, ans), f)
	case 2:
		a := w9R(30, 80)
		b := w9R(30, 80)
		w9I(f, "a", a)
		w9I(f, "b", b)
		c := math.Acos(math.Cos(float64(a)*math.Pi/180) * math.Cos(float64(b)*math.Pi/180))
		ans := w9Degrees(c)
		return w9P(kp,
			fmt.Sprintf("In a right spherical triangle, cos c = cos a cos b. Give the hypotenuse c in degrees for a = %d and b = %d, rounded to the nearest degree. (enter a number)", a, b),
			strconv.Itoa(ans),
			fmt.Sprintf("cos c is about %.3f, so c is about %d degrees.", math.Cos(c), ans), f)
	default:
		a := w9R(30, 70)
		b := w9R(30, 80)
		w9I(f, "a", a)
		w9I(f, "B", b)
		val := math.Cos(float64(a)*math.Pi/180) * math.Sin(float64(b)*math.Pi/180)
		ans := w9Degrees(math.Acos(val))
		return w9P(kp,
			fmt.Sprintf("In a right spherical triangle, cos A = cos a sin B. Give A in degrees for a = %d and B = %d, rounded to the nearest degree. (enter a number)", a, b),
			strconv.Itoa(ans),
			fmt.Sprintf("cos A is about %.3f, so A is about %d degrees.", val, ans), f)
	}
}

func w9SphericalPolar(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		a := w9R(30, 120)
		w9I(f, "a", a)
		ans := 180 - a
		return w9P(kp,
			fmt.Sprintf("The polar triangle of a spherical triangle has angle 180 - a opposite the side a. Give that angle for a = %d degrees. (enter a number)", a),
			strconv.Itoa(ans),
			fmt.Sprintf("180 - %d = %d.", a, ans), f)
	case 2:
		a := w9R(30, 120)
		w9I(f, "A", a)
		ans := 180 - a
		return w9P(kp,
			fmt.Sprintf("The polar triangle has side 180 - A opposite the angle A. Give that side for A = %d degrees. (enter a number)", a),
			strconv.Itoa(ans),
			fmt.Sprintf("180 - %d = %d.", a, ans), f)
	default:
		a := w9R(40, 80)
		b := w9R(40, 80)
		c := w9R(40, 80)
		w9I(f, "A", a)
		w9I(f, "B", b)
		w9I(f, "C", c)
		ans := 540 - (a + b + c)
		return w9P(kp,
			fmt.Sprintf("A spherical triangle has angles %d, %d and %d degrees. Its polar triangle has sides 180-A, 180-B and 180-C. Give the sum of those three sides in degrees. (enter a number)", a, b, c),
			strconv.Itoa(ans),
			fmt.Sprintf("540 - (%d + %d + %d) = %d.", a, b, c, ans), f)
	}
}

func w9SphericalPolygon(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		n := 4
		a := w9R(100, 140)
		w9I(f, "n", n)
		w9I(f, "A", a)
		ans := n*a - (n-2)*180
		return w9P(kp,
			fmt.Sprintf("A spherical quadrilateral has four equal angles of %d degrees. Its excess in degrees is nA - (n-2)180 with n = 4. Give the excess. (enter a number)", a),
			strconv.Itoa(ans),
			fmt.Sprintf("4 times %d minus 360 is %d.", a, ans), f)
	case 2:
		n := 5
		a := w9R(115, 140)
		w9I(f, "n", n)
		w9I(f, "A", a)
		ans := n*a - (n-2)*180
		return w9P(kp,
			fmt.Sprintf("A spherical pentagon has five equal angles of %d degrees. Its excess in degrees is nA - (n-2)180 with n = 5. Give the excess. (enter a number)", a),
			strconv.Itoa(ans),
			fmt.Sprintf("5 times %d minus 540 is %d.", a, ans), f)
	default:
		n := 6
		a := w9R(125, 150)
		w9I(f, "n", n)
		w9I(f, "A", a)
		ans := n*a - (n-2)*180
		return w9P(kp,
			fmt.Sprintf("A spherical hexagon has six equal angles of %d degrees. Its excess in degrees is nA - (n-2)180 with n = 6. Give the excess. (enter a number)", a),
			strconv.Itoa(ans),
			fmt.Sprintf("6 times %d minus 720 is %d.", a, ans), f)
	}
}

func w9Gradian(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		k := w9R(1, 12)
		w9I(f, "k", k)
		ans := 9 * k
		return w9P(kp,
			fmt.Sprintf("Convert %d gradians to degrees, using 100 gradians = 90 degrees. (enter a number)", 10*k),
			strconv.Itoa(ans),
			fmt.Sprintf("Each 10 gradians is 9 degrees, so %d gradians is %d degrees.", 10*k, ans), f)
	case 2:
		k := w9R(1, 15)
		w9I(f, "k", k)
		ans := 10 * k
		return w9P(kp,
			fmt.Sprintf("Convert %d degrees to gradians, using 90 degrees = 100 gradians. (enter a number)", 9*k),
			strconv.Itoa(ans),
			fmt.Sprintf("Each 9 degrees is 10 gradians, so %d degrees is %d gradians.", 9*k, ans), f)
	default:
		k := w9R(1, 6)
		w9I(f, "k", k)
		deg := 45 * k
		ans := 10 * deg / 9
		return w9P(kp,
			fmt.Sprintf("Convert %d degrees to gradians. (enter a number)", deg),
			strconv.Itoa(ans),
			fmt.Sprintf("%d times 10/9 is %d gradians.", deg, ans), f)
	}
}

func w9Signed(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		k := w9R(1, 35) * 10
		w9I(f, "k", k)
		ans := 360 - k
		return w9P(kp,
			fmt.Sprintf("An angle measures -%d degrees. Give its positive coterminal angle in [0, 360). (enter a number)", k),
			strconv.Itoa(ans),
			fmt.Sprintf("Add one full turn: 360 - %d = %d.", k, ans), f)
	case 2:
		k := w9R(1, 35) * 10
		w9I(f, "k", k)
		return w9P(kp,
			fmt.Sprintf("A clockwise rotation of %d degrees has signed measure -k. Give the magnitude k. (enter a number)", k),
			strconv.Itoa(k),
			fmt.Sprintf("The magnitude is %d degrees.", k), f)
	default:
		k := w9R(37, 100) * 10
		w9I(f, "k", k)
		ans := (360 - k%360) % 360
		return w9P(kp,
			fmt.Sprintf("An angle measures -%d degrees. Give its smallest positive coterminal angle in [0, 360). (enter a number)", k),
			strconv.Itoa(ans),
			fmt.Sprintf("Subtract whole turns and add one if needed: the result is %d degrees.", ans), f)
	}
}

func w9Conditional(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		a := w9R(10, 80)
		w9I(f, "A", a)
		ans := 90 - a
		return w9P(kp,
			fmt.Sprintf("If A + B = 90 degrees and A = %d degrees, give B. (enter a number)", a),
			strconv.Itoa(ans),
			fmt.Sprintf("B = 90 - %d = %d.", a, ans), f)
	case 2:
		a := w9R(30, 90)
		b := w9R(30, 90)
		w9I(f, "A", a)
		w9I(f, "B", b)
		ans := 180 - a - b
		return w9P(kp,
			fmt.Sprintf("In a triangle A + B + C = 180 degrees. If A = %d and B = %d, give C. (enter a number)", a, b),
			strconv.Itoa(ans),
			fmt.Sprintf("C = 180 - %d - %d = %d.", a, b, ans), f)
	default:
		b := w9R(10, 80)
		w9I(f, "B", b)
		ans := 90 - b
		return w9P(kp,
			fmt.Sprintf("If A + B = 90 degrees and B + C = 90 degrees, then A = C. Give A when B = %d degrees. (enter a number)", b),
			strconv.Itoa(ans),
			fmt.Sprintf("A = 90 - %d = %d.", b, ans), f)
	}
}

func w9TangentTriple(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		n := w9R(2, 8)
		w9I(f, "n", n)
		ans := 3*n*n - 1
		return w9P(kp,
			fmt.Sprintf("If tan(theta) = 1/n and tan(3theta) = (3n^2-1)/(n(n^2-3)), give the numerator for n = %d. (enter a number)", n),
			strconv.Itoa(ans),
			fmt.Sprintf("3 times %d squared minus 1 is %d.", n, ans), f)
	case 2:
		n := w9R(2, 8)
		w9I(f, "n", n)
		ans := n * (n*n - 3)
		return w9P(kp,
			fmt.Sprintf("If tan(theta) = 1/n and tan(3theta) = (3n^2-1)/(n(n^2-3)), give the denominator for n = %d. (enter a number)", n),
			strconv.Itoa(ans),
			fmt.Sprintf("%d times (%d squared minus 3) is %d.", n, n, ans), f)
	default:
		n := w9R(2, 7)
		w9I(f, "n", n)
		ans := 3*n - n*n*n
		return w9P(kp,
			fmt.Sprintf("If tan(theta) = n, then tan(3theta) = (3n-n^3)/(1-3n^2). Give the numerator for n = %d. (enter a number)", n),
			strconv.Itoa(ans),
			fmt.Sprintf("3 times %d minus %d cubed is %d.", n, n, ans), f)
	}
}
