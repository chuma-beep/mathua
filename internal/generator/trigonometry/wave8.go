package trigonometry

// Wave 8 trigonometry: the first four review concepts. All answers are
// integers and every operand is published in Facts so wave8_test.go can
// recompute the answer independently.

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave8(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"trig.basics.versine":        "versine",
		"trig.adv.chebyshev":         "chebyshev",
		"trig.adv.harmonic_addition": "harmonic",
		"trig.adv.arctan_sum":        "arctan_sum",
	} {
		reg.Register(id, &w8TrigGen{kind: kind})
	}
}

type w8TrigGen struct{ kind string }

func (g *w8TrigGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "versine":
		return w8Versine(kp)
	case "chebyshev":
		return w8Chebyshev(kp)
	case "harmonic":
		return w8Harmonic(kp)
	case "arctan_sum":
		return w8ArctanSum(kp)
	}
	return generator.Problem{Question: "What is 2 + 2? (enter a number)", Answer: "4", Explanation: "2 + 2 = 4."}
}

func w8T(o int, q, a, ex string, f map[string]string) generator.Problem {
	f["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: f, KPOrdinal: o}
}

func w8TI(f map[string]string, k string, v int) { f[k] = strconv.Itoa(v) }

func w8TR(lo, hi int) int {
	if hi <= lo {
		return lo
	}
	return lo + rand.Intn(hi-lo+1)
}

func w8Versine(kp int) generator.Problem {
	// Pythagorean triples supply exact sine and cosine values.
	triples := [][3]int{{3, 4, 5}, {4, 3, 5}, {5, 12, 13}, {12, 5, 13}, {8, 15, 17}, {15, 8, 17}, {7, 24, 25}, {24, 7, 25}, {9, 40, 41}, {40, 9, 41}, {20, 21, 29}, {21, 20, 29}, {12, 35, 37}, {35, 12, 37}, {28, 45, 53}, {45, 28, 53}, {11, 60, 61}, {60, 11, 61}, {16, 63, 65}, {63, 16, 65}, {33, 56, 65}, {56, 33, 65}, {13, 84, 85}, {84, 13, 85}}
	t := triples[rand.Intn(len(triples))]
	sinNum, cosNum, den := t[1], t[0], t[2]
	f := map[string]string{}
	w8TI(f, "sin_num", sinNum)
	w8TI(f, "cos_num", cosNum)
	w8TI(f, "den", den)
	switch kp {
	case 1:
		ans := den - cosNum
		return w8T(kp, fmt.Sprintf("For an angle with \\(\\cos\\theta=%d/%d\\), what is \\(%d(1-\\cos\\theta)\\)? (enter a number)", cosNum, den, den), strconv.Itoa(ans),
			fmt.Sprintf("The versine is \\(1-\\cos\\theta=(%d-%d)/%d\\), so %d times it is %d.", den, cosNum, den, den, ans), f)
	case 2:
		ans := den - sinNum
		return w8T(kp, fmt.Sprintf("For the same angle with \\(\\sin\\theta=%d/%d\\), what is \\(%d(1-\\sin\\theta)\\)? (enter a number)", sinNum, den, den), strconv.Itoa(ans),
			fmt.Sprintf("The coversine is \\(1-\\sin\\theta=(%d-%d)/%d\\), so %d times it is %d.", den, sinNum, den, den, ans), f)
	default:
		ans := 2*den - sinNum - cosNum
		return w8T(kp, fmt.Sprintf("For an angle with \\(\\sin\\theta=%d/%d\\) and \\(\\cos\\theta=%d/%d\\), what is %d times the sum versine plus coversine? (enter a number)", sinNum, den, cosNum, den, den), strconv.Itoa(ans),
			fmt.Sprintf("The sum is \\(2-(\\sin\\theta+\\cos\\theta)=(2\\cdot%d-%d-%d)/%d\\), so %d times it is %d.", den, sinNum, cosNum, den, den, ans), f)
	}
}

func w8Chebyshev(kp int) generator.Problem {
	n := w8TR(2, 21)
	x := w8TR(-2, 3)
	f := map[string]string{}
	w8TI(f, "n", n)
	w8TI(f, "x", x)
	switch kp {
	case 1:
		ans := w8ChebT(n, x)
		return w8T(kp, fmt.Sprintf("For \\(T_n(x)\\) defined by \\(T_0=1\\), \\(T_1=x\\) and \\(T_{n+1}=2xT_n-T_{n-1}\\), what is \\(T_{%d}(%d)\\)? (enter a number)", n, x), strconv.Itoa(ans),
			fmt.Sprintf("Apply the recurrence up to \\(T_{%d}\\). At \\(x=%d\\) the value is %d.", n, x, ans), f)
	case 2:
		ans := 1
		if n%2 == 1 {
			ans = -1
		}
		return w8T(kp, fmt.Sprintf("What is \\(T_{%d}(-1)\\)? (enter a number)", n), strconv.Itoa(ans),
			fmt.Sprintf("For every \\(n\\), \\(T_n(-1)=(-1)^n\\). With \\(n=%d\\) the value is %d.", n, ans), f)
	default:
		ans := 1 << (n - 1)
		return w8T(kp, fmt.Sprintf("What is the leading coefficient of \\(T_{%d}(x)\\)? (enter a number)", n), strconv.Itoa(ans),
			fmt.Sprintf("The leading term of \\(T_n\\) is \\(2^{n-1}x^n\\). For \\(n=%d\\) the coefficient is \\(2^{%d}=%d\\).", n, n-1, ans), f)
	}
}

func w8ChebT(n, x int) int {
	if n == 0 {
		return 1
	}
	if n == 1 {
		return x
	}
	a, b := 1, x
	for k := 2; k <= n; k++ {
		a, b = b, 2*x*b-a
	}
	return b
}

func w8Harmonic(kp int) generator.Problem {
	triples := [][3]int{{3, 4, 5}, {4, 3, 5}, {5, 12, 13}, {12, 5, 13}, {8, 15, 17}, {15, 8, 17}, {7, 24, 25}, {24, 7, 25}, {9, 40, 41}, {40, 9, 41}, {20, 21, 29}, {21, 20, 29}, {12, 35, 37}, {35, 12, 37}, {28, 45, 53}, {45, 28, 53}, {11, 60, 61}, {60, 11, 61}, {16, 63, 65}, {63, 16, 65}, {33, 56, 65}, {56, 33, 65}, {13, 84, 85}, {84, 13, 85}}
	t := triples[rand.Intn(len(triples))]
	a, b, r := t[0], t[1], t[2]
	f := map[string]string{}
	w8TI(f, "a", a)
	w8TI(f, "b", b)
	w8TI(f, "r", r)
	switch kp {
	case 1:
		return w8T(kp, fmt.Sprintf("Write \\(%d\\sin x+%d\\cos x\\) as \\(R\\sin(x+\\phi)\\). What is \\(R\\)? (enter a number)", a, b), strconv.Itoa(r),
			fmt.Sprintf("\\(R^2=a^2+b^2=%d^2+%d^2=%d\\), so \\(R=%d\\).", a, b, a*a+b*b, r), f)
	case 2:
		return w8T(kp, fmt.Sprintf("For \\(%d\\sin x+%d\\cos x=R\\sin(x+\\phi)\\), what is \\(R\\cos\\phi\\)? (enter a number)", a, b), strconv.Itoa(a),
			fmt.Sprintf("Expand \\(R\\sin(x+\\phi)\\): the sine coefficient is \\(R\\cos\\phi=%d\\).", a), f)
	default:
		return w8T(kp, fmt.Sprintf("For \\(%d\\sin x+%d\\cos x\\), what is the value at \\(x=0\\)? (enter a number)", a, b), strconv.Itoa(b),
			fmt.Sprintf("At \\(x=0\\), \\(\\sin 0=0\\) and \\(\\cos 0=1\\), so the value is %d.", b), f)
	}
}

func w8ArctanSum(kp int) generator.Problem {
	a := w8TR(0, 4)
	b := w8TR(0, 4)
	for a == 0 && b == 0 {
		b = w8TR(0, 4)
	}
	f := map[string]string{}
	w8TI(f, "a", a)
	w8TI(f, "b", b)
	switch kp {
	case 1:
		return w8T(kp, fmt.Sprintf("For \\(\\arctan %d+\\arctan %d\\), the tangent of the sum is \\((a+b)/(1-ab)\\). What is the numerator \\(a+b\\)? (enter a number)", a, b), strconv.Itoa(a+b),
			fmt.Sprintf("The numerator is \\(%d+%d=%d\\).", a, b, a+b), f)
	case 2:
		return w8T(kp, fmt.Sprintf("For \\(\\arctan %d+\\arctan %d\\), what is the denominator \\(1-ab\\)? (enter a number)", a, b), strconv.Itoa(1-a*b),
			fmt.Sprintf("The denominator is \\(1-%d\\cdot%d=%d\\).", a, b, 1-a*b), f)
	default:
		code := 0
		if a*b >= 1 {
			code = 1
		}
		return w8T(kp, fmt.Sprintf("For \\(\\arctan %d+\\arctan %d\\), enter 0 if the sum lies strictly between \\(-\\pi/2\\) and \\(\\pi/2\\), and 1 otherwise.", a, b), strconv.Itoa(code),
			fmt.Sprintf("The tangent-addition denominator is \\(1-ab=1-%d=%d\\). The sum stays strictly inside \\((-\\pi/2,\\pi/2)\\) exactly when \\(ab<1\\). The code is %d.", a*b, 1-a*b, code), f)
	}
}
