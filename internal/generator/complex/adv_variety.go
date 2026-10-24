package complex

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

// complexChoice builds a multiple choice problem with an explanation. These
// concepts are graded multiple_choice, and MultipleChoice returns no
// explanation, so every call site supplies one.
func complexChoice(question, correct, explanation string, distractors ...string) generator.Problem {
	p := generator.MultipleChoice(question, correct, distractors...)
	p.Explanation = explanation
	return p
}

type analyticGen struct{}

func (g *analyticGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	n := 2 + rand.Intn(4)
	a := 1 + rand.Intn(5)
	b := 1 + rand.Intn(5)
	switch rand.Intn(4) {
	case 0:
		return complexChoice(
			fmt.Sprintf("Is \\(f(z)=z^{%d}+%sz\\) analytic at \\(z=%d\\)?", n, unitCoef(b), a),
			"yes",
			"Polynomials are entire: they have a complex derivative at every point, so the answer is yes everywhere.",
			"no",
		)
	case 1:
		return complexChoice(
			"Is \\(f(z)=\\bar{z}\\) analytic anywhere?",
			"no",
			"Conjugating reverses orientation, and the Cauchy-Riemann equations fail at every point, so it is nowhere analytic.",
			"yes",
		)
	case 2:
		return complexChoice(
			fmt.Sprintf("Is \\(f(z)=\\frac{1}{z-%d}\\) analytic at \\(z=%d\\)?", a, a),
			"no",
			fmt.Sprintf("The function is undefined at \\(z=%d\\), so it cannot be differentiable there. It is analytic at every other point.", a),
			"yes",
		)
	default:
		return complexChoice(
			fmt.Sprintf("Is \\(f(z)=\\frac{1}{z-%d}\\) analytic at \\(z=%d\\)?", a, b),
			"yes",
			fmt.Sprintf("The only singular point is \\(z=%d\\), and %d is not that point, so the quotient rule applies there.", a, b),
			"no",
		)
	}
}

type cauchyRiemannGen struct{}

func (g *cauchyRiemannGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	k := 1 + rand.Intn(5)
	x := 1 + rand.Intn(4)
	switch rand.Intn(3) {
	case 0:
		return complexChoice(
			fmt.Sprintf("For \\(f(z)=z^{2}\\) written as \\(u+iv\\) with \\(u=x^{2}-y^{2}\\), what is \\(u_x\\) at \\((%d,%d)\\)?", x, k),
			strconv.Itoa(2*x),
			fmt.Sprintf("\\(u_x=2x\\), and at \\(x=%d\\) that is %d. Its partner \\(v_y=2x\\) matches, which is why the function is analytic.", x, 2*x),
			strconv.Itoa(2*k), strconv.Itoa(x), strconv.Itoa(x*x),
		)
	case 1:
		return complexChoice(
			"For \\(f=u+iv\\) with \\(u=x^{2}+y^{2}\\), do the Cauchy-Riemann equations hold anywhere?",
			"no",
			"\\(u_x=2x\\) and \\(v_y=-2x\\) would have to agree, forcing \\(x=0\\); \\(u_y=2y\\) and \\(v_x=2y\\) would then force \\(y=0\\), and at the origin the derivatives do not exist in the complex sense.",
			"yes",
		)
	default:
		return complexChoice(
			"Which pair of equations is the Cauchy-Riemann system for \\(f=u+iv\\)?",
			"u_x = v_y and u_y = -v_x",
			"Differentiability in the complex sense forces the two partials along the real and imaginary directions to agree, which is exactly this pair.",
			"u_x = v_x and u_y = v_y",
			"u_x = -v_y and u_y = v_x",
		)
	}
}

type conformalGen struct{}

func (g *conformalGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := 2 + rand.Intn(5)
	n := 2 + rand.Intn(3)
	switch rand.Intn(6) {
	case 0:
		return complexChoice(
			fmt.Sprintf("Is \\(w=z^{%d}\\) conformal at \\(z=%d\\)?", n, a),
			"yes",
			fmt.Sprintf("\\(f'(z)=%dz^{%d}\\), which is %d at \\(z=%d\\). A non-zero derivative preserves angles, so the map is conformal.", n, n-1, n*intPowLocal(a, n-1), a),
			"no",
		)
	case 1:
		return complexChoice(
			fmt.Sprintf("Is \\(w=z^{%d}\\) conformal at \\(z=0\\)?", n),
			"no",
			fmt.Sprintf("\\(f'(0)=0\\) for \\(n=%d\\), so angles are multiplied by %d rather than preserved. Conformality needs a non-zero derivative.", n, n),
			"yes",
		)
	case 2:
		return complexChoice(
			"Is \\(w=e^{z}\\) conformal everywhere?",
			"yes",
			"\\(f'(z)=e^{z}\\) never vanishes, so the map preserves angles at every point.",
			"no",
		)
	case 3:
		return complexChoice(
			fmt.Sprintf("Is \\(w=z^{2}+%dz\\) conformal at \\(z=%d\\)?", a, a),
			"yes",
			fmt.Sprintf("\\(f'(z)=2z+%d\\), which is %d at \\(z=%d\\), so the derivative is non-zero and the map is conformal there.", a, 3*a, a),
			"no",
		)
	case 4:
		return complexChoice(
			fmt.Sprintf("Is \\(w=1/z\\) conformal at \\(z=%d\\)?", a),
			"yes",
			fmt.Sprintf("\\(f'(z)=-1/z^{2}\\), which is not zero at \\(z=%d\\), so angles are preserved there.", a),
			"no",
		)
	default:
		return complexChoice(
			fmt.Sprintf("Is \\(w=\\bar{z}\\) conformal at \\(z=%d\\)?", a),
			"no",
			"Conjugation reverses orientation, so it cannot preserve angles at any point, whatever the value of \\(z\\).",
			"yes",
		)
	}
}

// intPowLocal is a small integer power used to print derivative values.
func intPowLocal(base, exp int) int {
	out := 1
	for i := 0; i < exp; i++ {
		out *= base
	}
	return out
}

type eulerIdentityGen struct{}

func (g *eulerIdentityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	k := 1 + rand.Intn(4)
	values := map[int]string{1: "i", 2: "-1", 3: "-i", 4: "1"}
	angles := map[int]string{1: `\pi/2`, 2: `\pi`, 3: `3\pi/2`, 4: `2\pi`}
	answer := values[k]
	others := []string{}
	for j := 1; j <= 4; j++ {
		if j != k {
			others = append(others, values[j])
		}
	}
	return complexChoice(
		fmt.Sprintf("Use Euler's formula to evaluate \\(e^{i%s}\\).", angles[k]),
		answer,
		fmt.Sprintf("\\(e^{i\\theta}=\\cos\\theta+i\\sin\\theta\\), and at \\(\\theta=%s\\) the cosine and sine give %s.", angles[k], answer),
		others...,
	)
}

// unitCoef renders a coefficient for a term x, dropping a leading one.
func unitCoef(c int) string {
	if c == 1 {
		return ""
	}
	return strconv.Itoa(c)
}

type harmonicGen struct{}

func (g *harmonicGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := 1 + rand.Intn(4)
	switch rand.Intn(4) {
	case 0:
		return complexChoice(
			fmt.Sprintf("Is \\(u=x^{2}-y^{2}\\) harmonic at \\((%d,%d)\\)?", a, a),
			"yes",
			"\\(u_{xx}=2\\) and \\(u_{yy}=-2\\), so the Laplacian is zero everywhere. It is the real part of \\(z^{2}\\).",
			"no",
		)
	case 1:
		return complexChoice(
			"Is \\(u=x^{2}+y^{2}\\) harmonic?",
			"no",
			"\\(u_{xx}+u_{yy}=2+2=4\\), which is not zero, so it satisfies no harmonic equation.",
			"yes",
		)
	case 2:
		return complexChoice(
			"Is \\(u=x^{3}-3xy^{2}\\) harmonic?",
			"yes",
			"It is the real part of \\(z^{3}\\), and real parts of analytic functions are harmonic. Checking directly: \\(6x-6x=0\\).",
			"no",
		)
	default:
		return complexChoice(
			fmt.Sprintf("What is the Laplacian of \\(u=x^{2}+y^{2}\\) at \\((%d,%d)\\)?", a, a),
			"4",
			"\\(u_{xx}=2\\) and \\(u_{yy}=2\\), so the Laplacian is 4 at every point, not just this one.",
			"0", "2", "8",
		)
	}
}

type mobiusTransformGen struct{}

func (g *mobiusTransformGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := 1 + rand.Intn(3)
	b := 1 + rand.Intn(6)
	c := 1 + rand.Intn(3)
	d := 1 + rand.Intn(6)
	det := a*d - b*c
	for det == 0 {
		d = 1 + rand.Intn(6)
		det = a*d - b*c
	}
	return complexChoice(
		fmt.Sprintf("For the Möbius map \\(f(z)=\\frac{%sz+%d}{%sz+%d}\\), what is \\(f(0)\\)?", unitCoef(a), b, unitCoef(c), d),
		fmt.Sprintf("%d/%d", b, d),
		fmt.Sprintf("Substituting \\(z=0\\) leaves the constant term %d over the constant term %d. The determinant %d is non-zero, so the map is invertible.", b, d, det),
		fmt.Sprintf("%d/%d", a, c), strconv.Itoa(b+d), fmt.Sprintf("-%d/%d", b, d),
	)
}

type residueGen struct{}

func (g *residueGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := 1 + rand.Intn(4)
	k := 1 + rand.Intn(5)
	switch rand.Intn(3) {
	case 0:
		return complexChoice(
			fmt.Sprintf("What is the residue of \\(f(z)=\\frac{%d}{z-%d}\\) at \\(z=%d\\)?", k, a, a),
			strconv.Itoa(k),
			"The residue of \\(c/(z-a)\\) is the numerator c, so the answer is the coefficient itself.",
			"0", strconv.Itoa(a), strconv.Itoa(k*a),
		)
	case 1:
		return complexChoice(
			fmt.Sprintf("What is the residue of \\(f(z)=\\frac{1}{(z-%d)^{2}}\\) at \\(z=%d\\)?", a, a),
			"0",
			"A pole of order two contributes nothing to the \\(1/(z-a)\\) coefficient, so the residue is zero.",
			"1", strconv.Itoa(a), strconv.Itoa(2*a),
		)
	default:
		return complexChoice(
			fmt.Sprintf("What is the residue of \\(f(z)=\\frac{z}{z-%d}\\) at \\(z=%d\\)?", a, a),
			strconv.Itoa(a),
			fmt.Sprintf("Rewrite \\(z=(z-%d)+%d\\), so the function is \\(1+%d/(z-%d)\\) and the residue is %d.", a, a, a, a, a),
			"1", "0", strconv.Itoa(-a),
		)
	}
}
