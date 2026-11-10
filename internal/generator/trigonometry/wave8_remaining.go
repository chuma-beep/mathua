package trigonometry

// Wave 8 trigonometry remaining concepts. Every answer is an integer; rational
// results are requested as a numerator, a denominator or a reduced component.

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave8Remaining(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"trig.ident.law_tangents":    "law_tangents",
		"trig.ident.mollweide":       "mollweide",
		"trig.adv.half_side":         "half_side",
		"trig.adv.napier":            "napier",
		"trig.spherical.law_sines":   "spherical_sines",
		"trig.spherical.law_cosines": "spherical_cosines",
		"trig.spherical.excess":      "spherical_excess",
		"trig.spherical.haversine":   "haversine",
		"trig.basics.exsecant":       "exsecant",
		"trig.adv.gudermannian":      "gudermannian",
		"trig.adv.weierstrass":       "weierstrass",
		"trig.adv.chebyshev_second":  "chebyshev_second",
		"trig.adv.ptolemy":           "ptolemy",
		"trig.adv.brahmagupta":       "brahmagupta",
		"trig.adv.machin":            "machin",
		"trig.adv.phasor":            "phasor",
	} {
		reg.Register(id, &w8TrigRemainingGen{kind: kind})
	}
}

type w8TrigRemainingGen struct{ kind string }

func (g *w8TrigRemainingGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "law_tangents":
		return w8LawTangents(kp)
	case "mollweide":
		return w8Mollweide(kp)
	case "half_side":
		return w8HalfSide(kp)
	case "napier":
		return w8Napier(kp)
	case "spherical_sines":
		return w8SphericalSines(kp)
	case "spherical_cosines":
		return w8SphericalCosines(kp)
	case "spherical_excess":
		return w8SphericalExcess(kp)
	case "haversine":
		return w8Haversine(kp)
	case "exsecant":
		return w8Exsecant(kp)
	case "gudermannian":
		return w8Gudermannian(kp)
	case "weierstrass":
		return w8Weierstrass(kp)
	case "chebyshev_second":
		return w8ChebyshevSecond(kp)
	case "ptolemy":
		return w8Ptolemy(kp)
	case "brahmagupta":
		return w8Brahmagupta(kp)
	case "machin":
		return w8Machin(kp)
	case "phasor":
		return w8Phasor(kp)
	}
	return generator.Problem{Question: "What is 2 + 2? (enter a number)", Answer: "4", Explanation: "2 + 2 = 4."}
}

func w8RTrig(o int, q, a, ex string, f map[string]string) generator.Problem {
	f["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: f, KPOrdinal: o}
}

func w8RRightTriple() (int, int, int) {
	triples := [][3]int{{3, 4, 5}, {4, 3, 5}, {5, 12, 13}, {12, 5, 13}, {8, 15, 17}, {15, 8, 17}, {7, 24, 25}, {24, 7, 25}, {9, 40, 41}, {40, 9, 41}}
	t := triples[rand.Intn(len(triples))]
	return t[0], t[1], t[2]
}

func w8RGCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func w8RI(f map[string]string, k string, v int) { f[k] = strconv.Itoa(v) }

func w8LawTangents(kp int) generator.Problem {
	a, b, c := w8RRightTriple()
	f := map[string]string{}
	w8RI(f, "a", a)
	w8RI(f, "b", b)
	w8RI(f, "c", c)
	switch kp {
	case 1:
		return w8RTrig(kp, fmt.Sprintf("For a right triangle with opposite sides %d and %d, the law of tangents gives \\(\\tan((A-B)/2)=(a-b)/(a+b)\\). What is the numerator \\(a-b\\)? (enter a number)", a, b), strconv.Itoa(a-b),
			fmt.Sprintf("\\(a-b=%d-%d=%d\\).", a, b, a-b), f)
	case 2:
		return w8RTrig(kp, fmt.Sprintf("For the same triangle with opposite sides %d and %d, what is the denominator \\(a+b\\)? (enter a number)", a, b), strconv.Itoa(a+b),
			fmt.Sprintf("\\(a+b=%d+%d=%d\\).", a, b, a+b), f)
	default:
		code := 2
		if a > b {
			code = 1
		}
		return w8RTrig(kp, fmt.Sprintf("For the two sides %d and %d, which is opposite the larger angle? Enter 1 for the first and 2 for the second.", a, b), strconv.Itoa(code),
			fmt.Sprintf("The larger side is opposite the larger angle. The answer is %d.", code), f)
	}
}

func w8Mollweide(kp int) generator.Problem {
	a, b, c := w8RRightTriple()
	f := map[string]string{}
	w8RI(f, "a", a)
	w8RI(f, "b", b)
	w8RI(f, "c", c)
	switch kp {
	case 1:
		return w8RTrig(kp, fmt.Sprintf("For a right triangle with legs %d and %d and hypotenuse %d, Mollweide's formula gives \\((a+b)/c\\). What is the numerator \\(a+b\\)? (enter a number)", a, b, c), strconv.Itoa(a+b),
			fmt.Sprintf("\\(a+b=%d+%d=%d\\).", a, b, a+b), f)
	case 2:
		return w8RTrig(kp, fmt.Sprintf("For the triangle with legs %d and %d and hypotenuse %d, what is the denominator \\(c\\) in \\((a+b)/c\\)? (enter a number)", a, b, c), strconv.Itoa(c),
			fmt.Sprintf("The hypotenuse is %d.", c), f)
	default:
		return w8RTrig(kp, fmt.Sprintf("For the triangle with sides %d, %d and %d, what is \\(a-b\\)? (enter a number)", a, b, c), strconv.Itoa(a-b),
			fmt.Sprintf("\\(a-b=%d-%d=%d\\).", a, b, a-b), f)
	}
}

func w8HalfSide(kp int) generator.Problem {
	a, b, c := w8RRightTriple()
	s := (a + b + c) / 2
	r := (a + b - c) / 2
	f := map[string]string{}
	w8RI(f, "a", a)
	w8RI(f, "b", b)
	w8RI(f, "c", c)
	w8RI(f, "s", s)
	w8RI(f, "r", r)
	switch kp {
	case 1:
		return w8RTrig(kp, fmt.Sprintf("For a triangle with sides %d, %d and %d, what is the semiperimeter \\(s\\)? (enter a number)", a, b, c), strconv.Itoa(s),
			fmt.Sprintf("\\(s=(%d+%d+%d)/2=%d\\).", a, b, c, s), f)
	case 2:
		return w8RTrig(kp, fmt.Sprintf("For the triangle with sides %d, %d and %d, what is \\(s-a\\)? (enter a number)", a, b, c), strconv.Itoa(s-a),
			fmt.Sprintf("\\(s-a=%d-%d=%d\\).", s, a, s-a), f)
	default:
		return w8RTrig(kp, fmt.Sprintf("For the right triangle with legs %d and %d and hypotenuse %d, what is the inradius \\(r=(a+b-c)/2\\)? (enter a number)", a, b, c), strconv.Itoa(r),
			fmt.Sprintf("\\(r=(%d+%d-%d)/2=%d\\).", a, b, c, r), f)
	}
}

func w8Napier(kp int) generator.Problem {
	a, b, c := w8RRightTriple()
	f := map[string]string{}
	w8RI(f, "a", a)
	w8RI(f, "b", b)
	w8RI(f, "c", c)
	switch kp {
	case 1:
		return w8RTrig(kp, fmt.Sprintf("Napier's analogy for a right triangle gives \\(\\tan((A-B)/2)=(a-b)/(a+b)\\). For sides %d and %d, what is the numerator \\(a-b\\)? (enter a number)", a, b), strconv.Itoa(a-b),
			fmt.Sprintf("\\(a-b=%d-%d=%d\\).", a, b, a-b), f)
	case 2:
		return w8RTrig(kp, fmt.Sprintf("For the triangle with sides %d and %d, what is the denominator \\(a+b\\)? (enter a number)", a, b), strconv.Itoa(a+b),
			fmt.Sprintf("\\(a+b=%d+%d=%d\\).", a, b, a+b), f)
	default:
		return w8RTrig(kp, fmt.Sprintf("For sides %d, %d and %d, what is \\(a+b+c\\)? (enter a number)", a, b, c), strconv.Itoa(a+b+c),
			fmt.Sprintf("\\(%d+%d+%d=%d\\).", a, b, c, a+b+c), f)
	}
}

func w8SphericalSines(kp int) generator.Problem {
	type sineCase struct {
		aN, aD, AN, AD, BN, BD, rN, rD int
	}
	cases := []sineCase{
		{3, 5, 4, 5, 5, 13, 15, 52},
		{4, 5, 3, 5, 5, 13, 20, 39},
		{3, 5, 5, 13, 3, 5, 117, 125},
		{5, 13, 3, 5, 4, 5, 20, 39},
		{5, 13, 4, 5, 3, 5, 15, 52},
		{4, 5, 5, 13, 5, 13, 4, 5},
		{3, 5, 8, 17, 5, 13, 51, 104},
		{4, 5, 8, 17, 5, 13, 17, 26},
		{3, 5, 7, 25, 5, 13, 75, 91},
		{3, 5, 5, 13, 5, 13, 3, 5},
	}
	c := cases[rand.Intn(len(cases))]
	f := map[string]string{}
	w8RI(f, "sin_a_num", c.aN)
	w8RI(f, "sin_a_den", c.aD)
	w8RI(f, "sin_A_num", c.AN)
	w8RI(f, "sin_A_den", c.AD)
	w8RI(f, "sin_B_num", c.BN)
	w8RI(f, "sin_B_den", c.BD)
	w8RI(f, "result_num", c.rN)
	w8RI(f, "result_den", c.rD)
	q := fmt.Sprintf("In a spherical triangle, \\(\\sin A=%d/%d\\), \\(\\sin a=%d/%d\\) and \\(\\sin B=%d/%d\\). The spherical sine rule gives \\(\\sin b=\\sin B\\sin a/\\sin A\\).", c.AN, c.AD, c.aN, c.aD, c.BN, c.BD)
	switch kp {
	case 1:
		return w8RTrig(kp, q+" What is the numerator of \\(\\sin b\\) in lowest terms? (enter a number)", strconv.Itoa(c.rN),
			fmt.Sprintf("The reduced value is \\(%d/%d\\).", c.rN, c.rD), f)
	case 2:
		return w8RTrig(kp, q+" What is the denominator of \\(\\sin b\\) in lowest terms? (enter a number)", strconv.Itoa(c.rD),
			fmt.Sprintf("The reduced value is \\(%d/%d\\).", c.rN, c.rD), f)
	default:
		return w8RTrig(kp, q+" What is the reduced numerator of \\(\\sin b\\)? (enter a number)", strconv.Itoa(c.rN),
			fmt.Sprintf("The reduced numerator is %d.", c.rN), f)
	}
}

func w8SphericalCosines(kp int) generator.Problem {
	type cosCase struct {
		caN, caD, saN, saD int
		cbN, cbD, sbN, sbD int
		cCN, cCD           int
		rcN, rcD           int
	}
	cases := []cosCase{
		{3, 5, 4, 5, 5, 13, 12, 13, 1, 2, 39, 65},
		{0, 1, 1, 1, 3, 5, 4, 5, 1, 2, 2, 5},
		{3, 5, 4, 5, 3, 5, 4, 5, 0, 1, 9, 25},
		{3, 5, 4, 5, 5, 13, 12, 13, -1, 2, -9, 65},
		{5, 13, 12, 13, 5, 13, 12, 13, 1, 2, 97, 169},
		{0, 1, 1, 1, 5, 13, 12, 13, 0, 1, 0, 1},
		{4, 5, 3, 5, 3, 5, 4, 5, 1, 2, 18, 25},
		{0, 1, 1, 1, 5, 13, 12, 13, 1, 1, 12, 13},
		{3, 5, 4, 5, 0, 1, 1, 1, 1, 2, 2, 5},
		{5, 13, 12, 13, 3, 5, 4, 5, -1, 2, -9, 65},
	}
	c := cases[rand.Intn(len(cases))]
	f := map[string]string{}
	w8RI(f, "cos_a_num", c.caN)
	w8RI(f, "cos_a_den", c.caD)
	w8RI(f, "sin_a_num", c.saN)
	w8RI(f, "sin_a_den", c.saD)
	w8RI(f, "cos_b_num", c.cbN)
	w8RI(f, "cos_b_den", c.cbD)
	w8RI(f, "sin_b_num", c.sbN)
	w8RI(f, "sin_b_den", c.sbD)
	w8RI(f, "cos_C_num", c.cCN)
	w8RI(f, "cos_C_den", c.cCD)
	w8RI(f, "cos_c_num", c.rcN)
	w8RI(f, "cos_c_den", c.rcD)
	q := fmt.Sprintf("In a spherical triangle, cos a=%d/%d, sin a=%d/%d, cos b=%d/%d, sin b=%d/%d and cos C=%d/%d.", c.caN, c.caD, c.saN, c.saD, c.cbN, c.cbD, c.sbN, c.sbD, c.cCN, c.cCD)
	switch kp {
	case 1:
		return w8RTrig(kp, q+" The spherical cosine rule gives cos c. What is the numerator of cos c in lowest terms? (enter a number)", strconv.Itoa(c.rcN),
			fmt.Sprintf("The spherical cosine rule gives cos c = %d/%d.", c.rcN, c.rcD), f)
	case 2:
		return w8RTrig(kp, q+" What is the denominator of cos c in lowest terms? (enter a number)", strconv.Itoa(c.rcD),
			fmt.Sprintf("The reduced value is %d/%d.", c.rcN, c.rcD), f)
	default:
		code := 0
		if c.rcN < 0 {
			code = 1
		}
		return w8RTrig(kp, q+" Is cos c negative? Enter 1 for yes and 0 for no.", strconv.Itoa(code),
			fmt.Sprintf("The reduced numerator is %d, so the answer is %d.", c.rcN, code), f)
	}
}

func w8SphericalExcess(kp int) generator.Problem {
	A := w8TR(50, 100)
	B := w8TR(50, 100)
	C := w8TR(50, 100)
	sum := A + B + C
	if sum <= 180 {
		C += 181 - sum
		sum = A + B + C
	}
	E := sum - 180
	f := map[string]string{}
	w8RI(f, "A", A)
	w8RI(f, "B", B)
	w8RI(f, "C", C)
	w8RI(f, "sum", sum)
	w8RI(f, "excess", E)
	switch kp {
	case 1:
		return w8RTrig(kp, fmt.Sprintf("A spherical triangle has angles %d, %d and %d degrees. What is its spherical excess in degrees? (enter a number)", A, B, C), strconv.Itoa(E),
			fmt.Sprintf("\\(E=A+B+C-180=%d\\).", E), f)
	case 2:
		return w8RTrig(kp, fmt.Sprintf("For the spherical triangle with angles %d, %d and %d, what is \\(A+B+C\\)? (enter a number)", A, B, C), strconv.Itoa(sum),
			fmt.Sprintf("\\(%d+%d+%d=%d\\).", A, B, C, sum), f)
	default:
		return w8RTrig(kp, fmt.Sprintf("For the same triangle with excess %d degrees and \\(A+B=%d\\), what is \\(C\\)? (enter a number)", E, A+B), strconv.Itoa(C),
			fmt.Sprintf("\\(C=180+E-(A+B)=180+%d-%d=%d\\).", E, A+B, C), f)
	}
}

func w8Haversine(kp int) generator.Problem {
	angles := []int{0, 60, 90, 120, 180, 240, 270, 300, 360}
	fourHav := []int{0, 1, 2, 3, 4, 3, 2, 1, 0}
	i := rand.Intn(len(angles))
	theta := angles[i]
	m := w8TR(1, 2)
	f := map[string]string{}
	w8RI(f, "angle", theta)
	w8RI(f, "four_hav", fourHav[i])
	w8RI(f, "scale", m)
	switch kp {
	case 1:
		return w8RTrig(kp, fmt.Sprintf("What is \\(%d\\times4\\operatorname{hav}(%d^\\circ)\\)? (enter a number)", m, theta), strconv.Itoa(m*fourHav[i]),
			fmt.Sprintf("Four times the haversine is %d, so %d times it is %d.", fourHav[i], m, m*fourHav[i]), f)
	case 2:
		return w8RTrig(kp, fmt.Sprintf("For the same angle \\(%d^\\circ\\), what is \\(%d\\times8\\operatorname{hav}(%d^\\circ)\\)? (enter a number)", theta, m, theta), strconv.Itoa(m*2*fourHav[i]),
			fmt.Sprintf("Eight times the haversine is %d, so %d times it is %d.", 2*fourHav[i], m, m*2*fourHav[i]), f)
	default:
		code := 0
		if fourHav[i] >= 2 {
			code = 1
		}
		return w8RTrig(kp, fmt.Sprintf("For \\(%d^\\circ\\), is \\(\\operatorname{hav}(%d^\\circ)\\ge 1/2\\)? Enter 1 for yes and 0 for no.", theta, theta), strconv.Itoa(code),
			fmt.Sprintf("Four times the haversine is %d, so the answer is %d.", fourHav[i], code), f)
	}
}

func w8Exsecant(kp int) generator.Problem {
	a, b, c := w8RRightTriple()
	f := map[string]string{}
	w8RI(f, "a", a)
	w8RI(f, "b", b)
	w8RI(f, "c", c)
	switch kp {
	case 1:
		return w8RTrig(kp, fmt.Sprintf("For a right triangle with adjacent side %d and hypotenuse %d, what is the exsecant numerator \\(c-a\\)? (enter a number)", a, c), strconv.Itoa(c-a),
			fmt.Sprintf("\\(\\operatorname{exsec}\\theta=\\sec\\theta-1=(c-a)/a=%d/%d\\), so the numerator is %d.", c-a, a, c-a), f)
	case 2:
		return w8RTrig(kp, fmt.Sprintf("For the triangle with adjacent side %d and hypotenuse %d, what is the exsecant denominator \\(a\\)? (enter a number)", a, c), strconv.Itoa(a),
			fmt.Sprintf("The denominator is the adjacent side %d.", a), f)
	default:
		return w8RTrig(kp, fmt.Sprintf("For the triangle with adjacent side %d and hypotenuse %d, what is \\(a\\) times the exsecant? (enter a number)", a, c), strconv.Itoa(c-a),
			fmt.Sprintf("\\(a(\\sec\\theta-1)=c-a=%d\\).", c-a), f)
	}
}

func w8Gudermannian(kp int) generator.Problem {
	k := w8TR(2, 20)
	f := map[string]string{}
	w8RI(f, "k", k)
	w8RI(f, "num", k*k-1)
	w8RI(f, "den", 2*k)
	switch kp {
	case 1:
		return w8RTrig(kp, fmt.Sprintf("For \\(x=\\ln %d\\), what is the numerator of \\(\\tan(\\operatorname{gd} x)\\) in lowest terms? (enter a number)", k), strconv.Itoa(k*k-1),
			fmt.Sprintf("\\(\\tan(\\operatorname{gd} x)=\\sinh x=(%d-1/%d)/2=(%d)/%d\\), so the numerator is %d.", k, k, k*k-1, 2*k, k*k-1), f)
	case 2:
		return w8RTrig(kp, fmt.Sprintf("For the same \\(x=\\ln %d\\), what is the denominator of \\(\\tan(\\operatorname{gd} x)\\)? (enter a number)", k), strconv.Itoa(2*k),
			fmt.Sprintf("The denominator is \\(2k=2\\times%d=%d\\).", k, 2*k), f)
	default:
		code := 0
		if k*k-1 > 2*k {
			code = 1
		}
		return w8RTrig(kp, fmt.Sprintf("For \\(x=\\ln %d\\), is \\(\\tan(\\operatorname{gd} x)>1\\)? Enter 1 for yes and 0 for no.", k), strconv.Itoa(code),
			fmt.Sprintf("Compare numerator %d with denominator %d. The answer is %d.", k*k-1, 2*k, code), f)
	}
}

func w8Weierstrass(kp int) generator.Problem {
	t := w8TR(1, 20)
	f := map[string]string{}
	w8RI(f, "t", t)
	w8RI(f, "den", 1+t*t)
	switch kp {
	case 1:
		return w8RTrig(kp, fmt.Sprintf("With \\(t=\\tan(x/2)=%d\\), the Weierstrass substitution gives \\(\\sin x=2t/(1+t^2)\\). What is the numerator \\(2t\\)? (enter a number)", t), strconv.Itoa(2*t),
			fmt.Sprintf("\\(2t=2\\times%d=%d\\).", t, 2*t), f)
	case 2:
		return w8RTrig(kp, fmt.Sprintf("For the same \\(t=%d\\), what is the numerator \\(1-t^2\\) of \\(\\cos x\\)? (enter a number)", t), strconv.Itoa(1-t*t),
			fmt.Sprintf("\\(1-t^2=1-%d^2=%d\\).", t, 1-t*t), f)
	default:
		return w8RTrig(kp, fmt.Sprintf("For the same \\(t=%d\\), what is the common denominator \\(1+t^2\\)? (enter a number)", t), strconv.Itoa(1+t*t),
			fmt.Sprintf("\\(1+t^2=1+%d^2=%d\\).", t, 1+t*t), f)
	}
}

func w8ChebyshevSecond(kp int) generator.Problem {
	n := w8TR(1, 20)
	x := w8TR(-2, 2)
	f := map[string]string{}
	w8RI(f, "n", n)
	w8RI(f, "x", x)
	u := func(n, x int) int {
		if n == 0 {
			return 1
		}
		if n == 1 {
			return 2 * x
		}
		a, b := 1, 2*x
		for k := 2; k <= n; k++ {
			a, b = b, 2*x*b-a
		}
		return b
	}
	switch kp {
	case 1:
		return w8RTrig(kp, fmt.Sprintf("For \\(U_n(x)\\) defined by \\(U_0=1\\), \\(U_1=2x\\) and \\(U_{n+1}=2xU_n-U_{n-1}\\), what is \\(U_{%d}(%d)\\)? (enter a number)", n, x), strconv.Itoa(u(n, x)),
			fmt.Sprintf("Apply the recurrence to reach \\(U_{%d}(%d)=%d\\).", n, x, u(n, x)), f)
	case 2:
		return w8RTrig(kp, fmt.Sprintf("What is the leading coefficient of \\(U_{%d}(x)\\)? (enter a number)", n), strconv.Itoa(1<<n),
			fmt.Sprintf("The leading term of \\(U_n\\) is \\(2^n x^n\\), so the coefficient is \\(2^{%d}=%d\\).", n, 1<<n), f)
	default:
		ans := n + 1
		if n%2 == 1 {
			ans = -(n + 1)
		}
		return w8RTrig(kp, fmt.Sprintf("What is \\(U_{%d}(-1)\\)? (enter a number)", n), strconv.Itoa(ans),
			fmt.Sprintf("\\(U_n(-1)=(-1)^n(n+1)\\). With \\(n=%d\\) the value is %d.", n, ans), f)
	}
}

func w8Ptolemy(kp int) generator.Problem {
	a := w8TR(2, 12)
	b := w8TR(2, 12)
	f := map[string]string{}
	w8RI(f, "a", a)
	w8RI(f, "b", b)
	switch kp {
	case 1:
		return w8RTrig(kp, fmt.Sprintf("A rectangle has sides %d and %d. Ptolemy's theorem says the product of its diagonals is \\(a^2+b^2\\). What is that product? (enter a number)", a, b), strconv.Itoa(a*a+b*b),
			fmt.Sprintf("\\(%d^2+%d^2=%d+%d=%d\\).", a, b, a*a, b*b, a*a+b*b), f)
	case 2:
		return w8RTrig(kp, fmt.Sprintf("For the rectangle with sides %d and %d, what is \\(ab\\)? (enter a number)", a, b), strconv.Itoa(a*b),
			fmt.Sprintf("\\(%d\\times%d=%d\\).", a, b, a*b), f)
	default:
		return w8RTrig(kp, fmt.Sprintf("For the rectangle with sides %d and %d, what is \\((a^2+b^2)-ab\\)? (enter a number)", a, b), strconv.Itoa(a*a+b*b-a*b),
			fmt.Sprintf("\\(%d-%d=%d\\).", a*a+b*b, a*b, a*a+b*b-a*b), f)
	}
}

func w8Brahmagupta(kp int) generator.Problem {
	a := w8TR(2, 12)
	b := w8TR(2, 12)
	f := map[string]string{}
	w8RI(f, "a", a)
	w8RI(f, "b", b)
	switch kp {
	case 1:
		return w8RTrig(kp, fmt.Sprintf("A cyclic quadrilateral has sides %d, %d, %d and %d. Brahmagupta's formula gives its area. What is the area? (enter a number)", a, b, a, b), strconv.Itoa(a*b),
			fmt.Sprintf("The semiperimeter is %d. The area is \\(\\sqrt{%d\\cdot%d\\cdot%d\\cdot%d}=%d\\).", a+b, b, a, a, b, a*b), f)
	case 2:
		return w8RTrig(kp, fmt.Sprintf("For the cyclic quadrilateral with sides %d, %d, %d and %d, what is the semiperimeter? (enter a number)", a, b, a, b), strconv.Itoa(a+b),
			fmt.Sprintf("\\(s=(%d+%d+%d+%d)/2=%d\\).", a, b, a, b, a+b), f)
	default:
		return w8RTrig(kp, fmt.Sprintf("For the cyclic quadrilateral with sides %d, %d, %d and %d, what is \\(s-a\\)? (enter a number)", a, b, a, b), strconv.Itoa(b),
			fmt.Sprintf("\\(s-a=(%d)-%d=%d\\).", a+b, a, b), f)
	}
}

func w8Machin(kp int) generator.Problem {
	type machinForm struct {
		k, a, b, sign int
	}
	forms := []machinForm{{4, 5, 239, -1}, {1, 2, 3, 1}, {2, 3, 7, 1}}
	form := forms[rand.Intn(len(forms))]
	k, a, b, sign := form.k, form.a, form.b, form.sign
	m := w8TR(1, 4)
	op := "+"
	if sign < 0 {
		op = "-"
	}
	f := map[string]string{}
	w8RI(f, "k", k)
	w8RI(f, "a", a)
	w8RI(f, "b", b)
	w8RI(f, "sign", sign)
	w8RI(f, "scale", m)
	switch kp {
	case 1:
		return w8RTrig(kp, fmt.Sprintf("In the Machin-type formula \\(\\pi/4=%d\\arctan(1/%d)%s\\arctan(1/%d)\\), what is %d times the coefficient \\(k\\)? (enter a number)", k, a, op, b, m), strconv.Itoa(m*k),
			fmt.Sprintf("The coefficient of \\(\\arctan(1/%d)\\) is %d.", a, k), f)
	case 2:
		return w8RTrig(kp, fmt.Sprintf("For \\(\\pi/4=%d\\arctan(1/%d)%s\\arctan(1/%d)\\), what is %d times a+b? (enter a number)", k, a, op, b, m), strconv.Itoa(m*(a+b)),
			fmt.Sprintf("\\(a+b=%d+%d=%d\\).", a, b, a+b), f)
	default:
		return w8RTrig(kp, fmt.Sprintf("For the formula \\(\\pi/4=%d\\arctan(1/%d)%s\\arctan(1/%d)\\), what is %d times k*a? (enter a number)", k, a, op, b, m), strconv.Itoa(m*k*a),
			fmt.Sprintf("\\(k\\times a=%d\\times%d=%d\\).", k, a, k*a), f)
	}
}

func w8Phasor(kp int) generator.Problem {
	a, b, r := w8RRightTriple()
	f := map[string]string{}
	w8RI(f, "a", a)
	w8RI(f, "b", b)
	w8RI(f, "r", r)
	switch kp {
	case 1:
		return w8RTrig(kp, fmt.Sprintf("A phasor is \\(%d+%di\\). What is its magnitude? (enter a number)", a, b), strconv.Itoa(r),
			fmt.Sprintf("\\(\\sqrt{%d^2+%d^2}=\\sqrt{%d}=%d\\).", a, b, a*a+b*b, r), f)
	case 2:
		return w8RTrig(kp, fmt.Sprintf("For the phasor \\(%d+%di\\), what is its real part? (enter a number)", a, b), strconv.Itoa(a),
			fmt.Sprintf("The real part is %d.", a), f)
	default:
		return w8RTrig(kp, fmt.Sprintf("For the phasor \\(%d+%di\\), what is its imaginary part? (enter a number)", a, b), strconv.Itoa(b),
			fmt.Sprintf("The imaginary part is %d.", b), f)
	}
}
