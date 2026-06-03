package complex

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/mathutil"
)

func Register(reg *generator.Registry) {
	reg.Register("complex.basics.concept", &conceptGen{})
	reg.Register("complex.ops.add_sub", &addSubGen{})
	reg.Register("complex.ops.mult", &multGen{})
	reg.Register("complex.ops.conjugate", &conjugateGen{})
	reg.Register("complex.ops.divide", &divideGen{})
	reg.Register("complex.adv.polar", &polarGen{})
	reg.Register("complex.adv.de_moivre", &deMoivreGen{})
	reg.Register("complex.adv.roots", &rootsGen{})

	reg.Register("complex.adv.exponential", &exponentialGen{})
	reg.Register("complex.adv.inequalities", &inequalitiesGen{})
}

func fmtComplex(r, i int) string {
	if i >= 0 {
		return fmt.Sprintf("%d+%di", r, i)
	}
	return fmt.Sprintf("%d%di", r, i)
}

type conceptGen struct{}

func (g *conceptGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := rand.Intn(max(1, scale*8)) + 1
	b := rand.Intn(max(1, scale*8)) + 1
	if rand.Intn(2) == 0 {
		return generator.Problem{
			Question:    fmt.Sprintf("What is the real part of \\(%s\\)?", fmtComplex(a, b)),
			Answer:      strconv.Itoa(a),
			Explanation: fmt.Sprintf("The real part of %s is %d.", fmtComplex(a, b), a),
		}
	}
	return generator.Problem{
		Question:    fmt.Sprintf("What is the imaginary part of \\(%s\\)?", fmtComplex(a, b)),
		Answer:      strconv.Itoa(b),
		Explanation: fmt.Sprintf("The imaginary part of %s is %d.", fmtComplex(a, b), b),
	}
}

type addSubGen struct{}

func (g *addSubGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := rand.Intn(max(1, scale*8)) + 1
	b := rand.Intn(max(1, scale*8)) + 1
	c := rand.Intn(max(1, scale*8)) + 1
	d := rand.Intn(max(1, scale*8)) + 1
	if rand.Intn(2) == 0 {
		r := a + c
		i := b + d
		return generator.Problem{
			Question:    fmt.Sprintf("\\((%s) + (%s) =\\) ?", fmtComplex(a, b), fmtComplex(c, d)),
			Answer:      fmtComplex(r, i),
			Explanation: fmt.Sprintf("(%s) + (%s) = (%d+%d) + (%d+%d)i = %s", fmtComplex(a, b), fmtComplex(c, d), a, c, b, d, fmtComplex(r, i)),
		}
	}
	r := a - c
	i := b - d
	return generator.Problem{
		Question:    fmt.Sprintf("\\((%s) - (%s) =\\) ?", fmtComplex(a, b), fmtComplex(c, d)),
		Answer:      fmtComplex(r, i),
		Explanation: fmt.Sprintf("(%s) - (%s) = (%d-%d) + (%d-%d)i = %s", fmtComplex(a, b), fmtComplex(c, d), a, c, b, d, fmtComplex(r, i)),
	}
}

type multGen struct{}

func (g *multGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := rand.Intn(max(1, scale*5)) + 1
	b := rand.Intn(max(1, scale*5)) + 1
	c := rand.Intn(max(1, scale*5)) + 1
	d := rand.Intn(max(1, scale*5)) + 1
	r := a*c - b*d
	i := a*d + b*c
	return generator.Problem{
		Question:    fmt.Sprintf("\\((%s)(%s) =\\) ?", fmtComplex(a, b), fmtComplex(c, d)),
		Answer:      fmtComplex(r, i),
		Explanation: fmt.Sprintf("(%s)(%s) = (%d)(%d) - (%d)(%d) + [(%d)(%d)+(%d)(%d)]i = %s", fmtComplex(a, b), fmtComplex(c, d), a, c, b, d, a, d, b, c, fmtComplex(r, i)),
	}
}

type conjugateGen struct{}

func (g *conjugateGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := rand.Intn(max(1, scale*8)) + 1
	b := rand.Intn(max(1, scale*8)) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("What is the conjugate of \\(%s\\)?", fmtComplex(a, b)),
		Answer:      fmtComplex(a, -b),
		Explanation: fmt.Sprintf("The conjugate of %s is %s.", fmtComplex(a, b), fmtComplex(a, -b)),
	}
}

type divideGen struct{}

func (g *divideGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	r := rand.Intn(max(1, scale*4)) + 1
	s := rand.Intn(max(1, scale*4)) + 1
	c := rand.Intn(max(1, scale*3)) + 1
	d := rand.Intn(max(1, scale*3)) + 1
	numR := r*c - s*d
	numI := r*d + s*c
	return generator.Problem{
		Question:    fmt.Sprintf("\\((%s) / (%s) =\\) ?", fmtComplex(numR, numI), fmtComplex(c, d)),
		Answer:      fmtComplex(r, s),
		Explanation: fmt.Sprintf("(%s) / (%s) = %s", fmtComplex(numR, numI), fmtComplex(c, d), fmtComplex(r, s)),
	}
}

type polarGen struct{}

func (g *polarGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	if rand.Intn(2) == 0 {
		type angle struct {
			deg     int
			cosReal int
			sinImag int
		}
		angles := []angle{
			{0, 1, 0},
			{90, 0, 1},
			{180, -1, 0},
			{270, 0, -1},
		}
		a := angles[rand.Intn(len(angles))]
		r := rand.Intn(max(1, scale*4)) + 2
		realPart := r * a.cosReal
		imagPart := r * a.sinImag
		return generator.Problem{
			Question:    fmt.Sprintf("Convert \\(%d(\\cos %d^{\\circ} + i\\sin %d^{\\circ})\\) to rectangular form.", r, a.deg, a.deg),
			Answer:      fmtComplex(realPart, imagPart),
			Explanation: fmt.Sprintf("%d(cos %d° + i sin %d°) = %d(%d) + %d(%d)i = %s", r, a.deg, a.deg, r, a.cosReal, r, a.sinImag, fmtComplex(realPart, imagPart)),
		}
	}

	type triple struct{ a, b, r int }
	triples := []triple{
		{3, 4, 5},
		{4, 3, 5},
		{5, 12, 13},
		{12, 5, 13},
		{6, 8, 10},
		{8, 6, 10},
	}
	t := triples[rand.Intn(len(triples))]
	theta := math.Atan2(float64(t.b), float64(t.a)) * 180 / math.Pi
	thetaStr := strconv.Itoa(int(math.Round(theta)))
	return generator.Problem{
		Question:    fmt.Sprintf("Convert \\(%s\\) to polar form.", fmtComplex(t.a, t.b)),
		Answer:      fmt.Sprintf("%d(cos %s° + i sin %s°)", t.r, thetaStr, thetaStr),
		Explanation: fmt.Sprintf("r = √(a²+b²) = √(%d²+%d²) = %d. θ = arctan(b/a) = arctan(%d/%d) ≈ %s°", t.a, t.b, t.r, t.b, t.a, thetaStr),
	}
}

type deMoivreGen struct{}

func (g *deMoivreGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	type stdAngle struct{ deg int }
	angles := []stdAngle{
		{30}, {45}, {60}, {90},
	}
	angle := angles[rand.Intn(len(angles))]
	r := rand.Intn(max(1, scale*4)) + 2
	n := rand.Intn(max(1, scale*2)) + 2

	newR := mathutil.IntPow(r, n)
	newDeg := (angle.deg * n) % 360

	var ans string
	switch newDeg {
	case 0:
		ans = fmtComplex(newR, 0)
	case 90:
		ans = fmtComplex(0, newR)
	case 180:
		ans = fmtComplex(-newR, 0)
	case 270:
		ans = fmtComplex(0, -newR)
	default:
		ans = fmt.Sprintf("%d(cos %d° + i sin %d°)", newR, newDeg, newDeg)
	}

	return generator.Problem{
		Question:    fmt.Sprintf("\\([%d(\\cos %d^{\\circ} + i\\sin %d^{\\circ})]^{%d} =\\) ?", r, angle.deg, angle.deg, n),
		Answer:      ans,
		Explanation: fmt.Sprintf("[%d(cos %d° + i sin %d°)]^%d = %d(cos %d×%d° + i sin %d×%d°) = %d(cos %d° + i sin %d°) = %s", r, angle.deg, angle.deg, n, r, angle.deg, n, angle.deg, n, newR, newDeg, newDeg, ans),
	}
}

type rootsGen struct{}

type exponentialGen struct{}

func (g *exponentialGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	type entry struct {
		theta    int
		realPart int
		imagPart int
	}
	angles := []entry{
		{0, 1, 0},
		{90, 0, 1},
		{180, -1, 0},
		{270, 0, -1},
		{45, 0, 1},
		{60, 0, 1},
	}
	e := angles[rand.Intn(len(angles))]
	r := rand.Intn(max(1, scale*4)) + 2
	if e.theta != 0 && e.theta != 90 && e.theta != 180 && e.theta != 270 {
		return generator.Problem{
			Question:    fmt.Sprintf("Does \\(e^{i\\theta} = \\cos\\theta + i\\sin\\theta\\) hold for \\(\\theta = %d^{\\circ}\\)? (yes/no)", e.theta),
			Answer:      "yes",
			Explanation: "Euler's formula e^(iθ) = cos θ + i sin θ holds for all real θ.",
		}
	}
	realPart := r * e.realPart
	imagPart := r * e.imagPart
	if rand.Intn(2) == 0 {
		return generator.Problem{
			Question:    fmt.Sprintf("Write \\(%d e^{i \\cdot %d^{\\circ}}\\) in rectangular form \\((a+bi)\\).", r, e.theta),
			Answer:      fmtComplex(realPart, imagPart),
			Explanation: fmt.Sprintf("%d(cos %d° + i sin %d°) = %d(%d) + %d(%d)i = %s", r, e.theta, e.theta, r, e.realPart, r, e.imagPart, fmtComplex(realPart, imagPart)),
		}
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Convert \\(%s\\) to polar exponential form using Euler's formula.", fmtComplex(realPart, imagPart)),
		Answer:      fmt.Sprintf("%d e^(i·%d°)", r, e.theta),
		Explanation: fmt.Sprintf("r = √(%d²+%d²) = %d, θ = %d°, so %s = %d e^(i·%d°)", realPart, imagPart, r, e.theta, fmtComplex(realPart, imagPart), r, e.theta),
	}
}

type inequalitiesGen struct{}

func (g *inequalitiesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	type entry struct {
		a, b    int
		desc    string
	}
	entries := []entry{
		{3, 4, "\\(|3+4i|\\)"},
		{5, 12, "\\(|5+12i|\\)"},
		{8, 6, "\\(|8+6i|\\)"},
		{7, 24, "\\(|7+24i|\\)"},
		{9, 12, "\\(|9+12i|\\)"},
	}
	e := entries[rand.Intn(len(entries))]
	if rand.Intn(2) == 0 {
		mag := int(math.Sqrt(float64(e.a*e.a + e.b*e.b)))
		return generator.Problem{
			Question:    fmt.Sprintf("What is %s?", e.desc),
			Answer:      fmt.Sprintf("%d", mag),
			Explanation: fmt.Sprintf("|%d+%di| = √(%d²+%d²) = √%d = %d", e.a, e.b, e.a, e.b, e.a*e.a+e.b*e.b, mag),
		}
	}
	r1 := rand.Intn(max(1, scale*5)) + 1
	i1 := rand.Intn(max(1, scale*5)) + 1
	r2 := rand.Intn(max(1, scale*5)) + 1
	i2 := rand.Intn(max(1, scale*5)) + 1
	// Triangle inequality: |z1+z2| ≤ |z1|+|z2|
	s := fmtComplex(r1, i1)
	t := fmtComplex(r2, i2)
	return generator.Problem{
		Question:    fmt.Sprintf("Let \\(z_{1}=%s\\), \\(z_{2}=%s\\). Which is always true? (enter '\\(\\leq\\)' for \\(|z_{1}+z_{2}| \\leq |z_{1}|+|z_{2}|\\), or '\\(\\geq\\)' for the reverse)", s, t),
		Answer:      "≤",
		Explanation: "The triangle inequality |z₁+z₂| ≤ |z₁|+|z₂| holds for all complex numbers.",
	}
}

func (g *rootsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := rand.Intn(max(1, scale*4)) + 1
	b := rand.Intn(max(1, scale*4)) + 1
	c := a*a - b*b
	d := 2 * a * b

	return generator.Problem{
		Question:    fmt.Sprintf("Find the square roots of \\(%s\\).", fmtComplex(c, d)),
		Answer:      fmt.Sprintf("±%s", fmtComplex(a, b)),
		Explanation: fmt.Sprintf("(%s)² = %s, so the square roots are ±%s.", fmtComplex(a, b), fmtComplex(c, d), fmtComplex(a, b)),
	}
}
