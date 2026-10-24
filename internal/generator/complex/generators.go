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
	reg.Register("complex.ops.modulus", &modulusGen{})
	reg.Register("complex.ops.argument", &argumentGen{})
	reg.Register("complex.adv.euler_identity", &eulerIdentityGen{})
	reg.Register("complex.adv.cauchy_riemann", &cauchyRiemannGen{})
	reg.Register("complex.adv.residue", &residueGen{})
	reg.Register("complex.adv.analytic", &analyticGen{})
	reg.Register("complex.series.laurent", &laurentGen{})
	reg.Register("complex.integral.contour", &contourGen{})
	reg.Register("complex.adv.harmonic", &harmonicGen{})
	reg.Register("complex.topology.riemann_sphere", &riemannSphereGen{})
	reg.Register("complex.adv.conformal", &conformalGen{})
	reg.Register("complex.integral.cauchy_goursat", &cauchyGoursatGen{})
	reg.Register("complex.series.power", &powerSeriesGen{})
	reg.Register("complex.adv.mobius", &mobiusTransformGen{})
	reg.Register("complex.adv.schwarz_lemma", &schwarzLemmaGen{})
	reg.Register("complex.adv.logarithm", &logarithmGen{})

	registerBatch(reg)
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

	if rand.Intn(2) == 0 {
		return generator.Problem{
			Question:    fmt.Sprintf("Convert \\(%d(\\cos %d^{\\circ} + i\\sin %d^{\\circ})\\) to rectangular form.", r, a.deg, a.deg),
			Answer:      fmtComplex(realPart, imagPart),
			Explanation: fmt.Sprintf("%d(cos %d° + i sin %d°) = %d(%d) + %d(%d)i = %s", r, a.deg, a.deg, r, a.cosReal, r, a.sinImag, fmtComplex(realPart, imagPart)),
		}
	}

	return generator.Problem{
		Question:    fmt.Sprintf("Convert \\(%s\\) to polar form.", fmtComplex(realPart, imagPart)),
		Answer:      fmt.Sprintf("%d(cos %d° + i sin %d°)", r, a.deg, a.deg),
		Explanation: fmt.Sprintf("r = √(a²+b²) = √((%d)²+(%d)²) = %d. θ = %d°", realPart, imagPart, r, a.deg),
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
	if rand.Intn(2) == 0 {
		return generator.Problem{
			Question:    fmt.Sprintf("What is the modulus of \\(%s\\)? (enter a number)", fmtComplex(realPart, imagPart)),
			Answer:      fmt.Sprintf("%d", r),
			Explanation: fmt.Sprintf("Modulus r = √(%d²+%d²) = %d.", realPart, imagPart, r),
		}
	}
	return generator.Problem{
		Question:    fmt.Sprintf("What is the argument in degrees of \\(%s\\)? (enter a number)", fmtComplex(realPart, imagPart)),
		Answer:      fmt.Sprintf("%d", e.theta),
		Explanation: fmt.Sprintf("The point sits at angle %d°: r = %d, θ = %d°.", e.theta, r, e.theta),
	}
}

type inequalitiesGen struct{}

func (g *inequalitiesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	type entry struct {
		a, b int
		desc string
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
		Question:    fmt.Sprintf("Let z1=%s, z2=%s. Does the triangle inequality use less-or-equal, that is |z1+z2| is less-or-equal to |z1|+|z2|? (yes/no)", s, t),
		Answer:      "yes",
		Explanation: "The triangle inequality |z1+z2| <= |z1|+|z2| holds for all complex numbers.",
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

type modulusGen struct{}

func (g *modulusGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type triple struct{ a, b, r int }
	tableEasy := []triple{
		{3, 4, 5}, {5, 12, 13}, {8, 6, 10}, {7, 24, 25}, {9, 12, 15},
	}
	tableHard := []triple{
		{9, 40, 41}, {12, 35, 37}, {20, 21, 29},
	}
	pool := tableEasy
	if scale > 3 {
		pool = append(tableEasy, tableHard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{
		Question:    fmt.Sprintf("What is \\(|%s|\\)?", fmtComplex(e.a, e.b)),
		Answer:      fmt.Sprintf("%d", e.r),
		Explanation: fmt.Sprintf("|%s| = √(%d²+%d²) = √%d = %d", fmtComplex(e.a, e.b), e.a, e.b, e.a*e.a+e.b*e.b, e.r),
	}
}

type argumentGen struct{}

func (g *argumentGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	k := rand.Intn(max(1, int(1+ctx.Difficulty*4))) + 1

	type angleCase struct {
		zLatex      string
		deg         int
		quadrant    string
		explanation string
	}

	var cases []angleCase

	kStr := strconv.Itoa(k)
	if k == 1 {
		kStr = ""
	}
	cases = append(cases,
		angleCase{
			zLatex:      strconv.Itoa(k),
			deg:         0,
			quadrant:    "positive real axis",
			explanation: fmt.Sprintf("The complex number \\(z = %d\\) lies on the positive real axis. Therefore, \\(\\arg(z) = 0^\\circ\\).", k),
		},
		angleCase{
			zLatex:      kStr + "i",
			deg:         90,
			quadrant:    "positive imaginary axis",
			explanation: fmt.Sprintf("The complex number \\(z = %si\\) lies on the positive imaginary axis. Therefore, \\(\\arg(z) = 90^\\circ\\).", kStr),
		},
		angleCase{
			zLatex:      fmt.Sprintf("-%d", k),
			deg:         180,
			quadrant:    "negative real axis",
			explanation: fmt.Sprintf("The complex number \\(z = -%d\\) lies on the negative real axis. Therefore, \\(\\arg(z) = 180^\\circ\\).", k),
		},
		angleCase{
			zLatex:      "-" + kStr + "i",
			deg:         270,
			quadrant:    "negative imaginary axis",
			explanation: fmt.Sprintf("The complex number \\(z = -%si\\) lies on the negative imaginary axis. Therefore, \\(\\arg(z) = 270^\\circ\\).", kStr),
		},
	)

	var z45, z135, z225, z315 string
	if k == 1 {
		z45 = "1 + i"
		z135 = "-1 + i"
		z225 = "-1 - i"
		z315 = "1 - i"
	} else {
		z45 = fmt.Sprintf("%d + %di", k, k)
		z135 = fmt.Sprintf("-%d + %di", k, k)
		z225 = fmt.Sprintf("-%d - %di", k, k)
		z315 = fmt.Sprintf("%d - %di", k, k)
	}

	cases = append(cases,
		angleCase{
			zLatex:      z45,
			deg:         45,
			quadrant:    "quadrant 1",
			explanation: fmt.Sprintf("The complex number \\(z = %s\\) lies in quadrant 1 with equal real and imaginary parts. The reference angle is \\(\\arctan(1) = 45^\\circ\\), so \\(\\arg(z) = 45^\\circ\\).", z45),
		},
		angleCase{
			zLatex:      z135,
			deg:         135,
			quadrant:    "quadrant 2",
			explanation: fmt.Sprintf("The complex number \\(z = %s\\) lies in quadrant 2. The reference angle is \\(\\arctan(1) = 45^\\circ\\), so \\(\\arg(z) = 180^\\circ - 45^\\circ = 135^\\circ\\).", z135),
		},
		angleCase{
			zLatex:      z225,
			deg:         225,
			quadrant:    "quadrant 3",
			explanation: fmt.Sprintf("The complex number \\(z = %s\\) lies in quadrant 3. The reference angle is \\(\\arctan(1) = 45^\\circ\\), so \\(\\arg(z) = 180^\\circ + 45^\\circ = 225^\\circ\\).", z225),
		},
		angleCase{
			zLatex:      z315,
			deg:         315,
			quadrant:    "quadrant 4",
			explanation: fmt.Sprintf("The complex number \\(z = %s\\) lies in quadrant 4. The reference angle is \\(\\arctan(1) = 45^\\circ\\), so \\(\\arg(z) = 360^\\circ - 45^\\circ = 315^\\circ\\).", z315),
		},
	)

	if ctx.Difficulty > 0.3 {
		var sqrt3Term, kTerm string
		if k == 1 {
			sqrt3Term = "\\sqrt{3}"
			kTerm = "1"
		} else {
			sqrt3Term = fmt.Sprintf("%d\\sqrt{3}", k)
			kTerm = strconv.Itoa(k)
		}
		var iTerm string
		if k == 1 {
			iTerm = "i"
		} else {
			iTerm = fmt.Sprintf("%di", k)
		}

		z30 := fmt.Sprintf("%s + %s", sqrt3Term, iTerm)
		z60 := fmt.Sprintf("%s + %s i", kTerm, sqrt3Term)
		z120 := fmt.Sprintf("-%s + %s i", kTerm, sqrt3Term)
		z150 := fmt.Sprintf("-%s + %s", sqrt3Term, iTerm)
		z210 := fmt.Sprintf("-%s - %s", sqrt3Term, iTerm)
		z240 := fmt.Sprintf("-%s - %s i", kTerm, sqrt3Term)
		z300 := fmt.Sprintf("%s - %s i", kTerm, sqrt3Term)
		z330 := fmt.Sprintf("%s - %s", sqrt3Term, iTerm)

		cases = append(cases,
			angleCase{
				zLatex:      z30,
				deg:         30,
				quadrant:    "quadrant 1",
				explanation: fmt.Sprintf("The complex number \\(z = %s\\) lies in quadrant 1 with imaginary-to-real ratio \\(1/\\sqrt{3}\\). The reference angle is \\(\\arctan(1/\\sqrt{3}) = 30^\\circ\\), so \\(\\arg(z) = 30^\\circ\\).", z30),
			},
			angleCase{
				zLatex:      z60,
				deg:         60,
				quadrant:    "quadrant 1",
				explanation: fmt.Sprintf("The complex number \\(z = %s\\) lies in quadrant 1 with imaginary-to-real ratio \\(\\sqrt{3}\\). The reference angle is \\(\\arctan(\\sqrt{3}) = 60^\\circ\\), so \\(\\arg(z) = 60^\\circ\\).", z60),
			},
			angleCase{
				zLatex:      z120,
				deg:         120,
				quadrant:    "quadrant 2",
				explanation: fmt.Sprintf("The complex number \\(z = %s\\) lies in quadrant 2. The reference angle is \\(\\arctan(\\sqrt{3}) = 60^\\circ\\), so \\(\\arg(z) = 180^\\circ - 60^\\circ = 120^\\circ\\).", z120),
			},
			angleCase{
				zLatex:      z150,
				deg:         150,
				quadrant:    "quadrant 2",
				explanation: fmt.Sprintf("The complex number \\(z = %s\\) lies in quadrant 2. The reference angle is \\(\\arctan(1/\\sqrt{3}) = 30^\\circ\\), so \\(\\arg(z) = 180^\\circ - 30^\\circ = 150^\\circ\\).", z150),
			},
			angleCase{
				zLatex:      z210,
				deg:         210,
				quadrant:    "quadrant 3",
				explanation: fmt.Sprintf("The complex number \\(z = %s\\) lies in quadrant 3. The reference angle is \\(\\arctan(1/\\sqrt{3}) = 30^\\circ\\), so \\(\\arg(z) = 180^\\circ + 30^\\circ = 210^\\circ\\).", z210),
			},
			angleCase{
				zLatex:      z240,
				deg:         240,
				quadrant:    "quadrant 3",
				explanation: fmt.Sprintf("The complex number \\(z = %s\\) lies in quadrant 3. The reference angle is \\(\\arctan(\\sqrt{3}) = 60^\\circ\\), so \\(\\arg(z) = 180^\\circ + 60^\\circ = 240^\\circ\\).", z240),
			},
			angleCase{
				zLatex:      z300,
				deg:         300,
				quadrant:    "quadrant 4",
				explanation: fmt.Sprintf("The complex number \\(z = %s\\) lies in quadrant 4. The reference angle is \\(\\arctan(\\sqrt{3}) = 60^\\circ\\), so \\(\\arg(z) = 360^\\circ - 60^\\circ = 300^\\circ\\).", z300),
			},
			angleCase{
				zLatex:      z330,
				deg:         330,
				quadrant:    "quadrant 4",
				explanation: fmt.Sprintf("The complex number \\(z = %s\\) lies in quadrant 4. The reference angle is \\(\\arctan(1/\\sqrt{3}) = 30^\\circ\\), so \\(\\arg(z) = 360^\\circ - 30^\\circ = 330^\\circ\\).", z330),
			},
		)
	}

	c := cases[rand.Intn(len(cases))]
	q := fmt.Sprintf("What is the principal argument \\(\\arg(z)\\) in degrees for \\(z = %s\\)? (enter an integer in \\([0, 360)\\))", c.zLatex)
	ans := strconv.Itoa(c.deg)

	return generator.Problem{
		Question:    q,
		Answer:      ans,
		Explanation: c.explanation,
		Facts: map[string]string{
			"z":        c.zLatex,
			"arg_deg":  strconv.Itoa(c.deg),
			"quadrant": c.quadrant,
		},
	}
}

type laurentGen struct{}

func (g *laurentGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type qa struct {
		q string
		a string
		e string
	}
	easy := []qa{
		{q: "For f(z)=1/z, what is a_{-1} at 0? (enter a number)", a: "1", e: "1/z = z^{-1} is its own Laurent series; principal part is 1/z."},
		{q: "For f(z)=1/z^2, what is the residue at 0? (enter a number)", a: "0", e: "The z^{-1} slot is empty, so a_{-1}=0."},
		{q: "Does a Laurent series include negative powers for poles? (yes/no)", a: "yes", e: "Taylor has only non-negative powers; Laurent adds the principal part."},
		{q: "Is 1/z analytic at 0 (a removable singularity)? (yes/no)", a: "no", e: "1/z blows up at 0; it is a pole, not removable."},
		{q: "For f(z)=1/z^3, what is a_{-1} at 0? (enter a number)", a: "0", e: "Only the z^{-3} term is nonzero; a_{-1}=0."},
		{q: "Is the coefficient a_{-1} the residue? (yes/no)", a: "yes", e: "Residue is Laurent a_{-1}."},
	}
	hard := []qa{
		{q: "Is classification: pole (finite principal part), essential (infinite), removable (none)? (yes/no)", a: "yes", e: "Laurent principal part classifies singularities."},
		{q: "Does e^{1/z} have infinitely many negative terms at 0? (yes/no)", a: "yes", e: "Essential singularity: Laurent has infinite principal part."},
		{q: "Is annulus 0<|z|<1 the domain for Laurent of 1/(z(1-z))? (yes/no)", a: "yes", e: "Poles at 0,1; annuli determined by radii to nearest singularities."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type contourGen struct{}

func (g *contourGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type qa struct {
		q string
		a string
		e string
	}
	easy := []qa{
		{q: "What is the contour integral over |z|=1 of dz/z, divided by 2πi? (enter a number)", a: "1", e: "Residue 1 at 0 inside, so the integral is 2πi; divided by 2πi gives 1."},
		{q: "What is the contour integral over |z|=1 of dz/z^2, divided by 2πi? (enter a number)", a: "0", e: "Residue 0 at the double pole, so the integral is 0."},
		{q: "Is a contour integral independent of path in a simply connected analytic domain? (yes/no)", a: "yes", e: "Cauchy's theorem: integral around null-homotopic loop is 0."},
		{q: "Does the ML-estimate bound |integral| by M times L? (yes/no)", a: "yes", e: "ML inequality: sup |f| times length."},
		{q: "Is the integral of dz/z around |z|=1 equal to 0? (yes/no)", a: "no", e: "It is 2πi, nonzero: the pole at 0 is enclosed."},
	}
	hard := []qa{
		{q: "Does Cauchy integral formula give f(a)=∮ f(z)/(z-a) dz /2πi? (yes/no)", a: "yes", e: "Cauchy's formula for analytic f."},
		{q: "Is real integral ∫_{-∞}^{∞} dx/(x^2+1)=π via contour? (yes/no)", a: "yes", e: "Semicircle contour with poles at ±i gives π."},
		{q: "Does deformation of contour not cross singularities preserve integral? (yes/no)", a: "yes", e: "Homotopy invariance."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type riemannSphereGen struct{}

func (g *riemannSphereGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type qa struct {
		q string
		a string
		e string
	}
	easy := []qa{
		{q: "How many points are added to C to form the Riemann sphere? (enter a number)", a: "1", e: "One point at infinity compactifies C."},
		{q: "Under w=1/z, where does 0 go? (enter infinity for the point at infinity)", a: "infinity", e: "1/0 is the north pole: infinity."},
		{q: "Is the extended complex plane compact? (yes/no)", a: "yes", e: "One-point compactification of C."},
		{q: "Is C itself compact? (yes/no)", a: "no", e: "C is not compact; adding infinity compactifies it."},
		{q: "Does stereographic projection send circles on the sphere to circles or lines? (yes/no)", a: "yes", e: "Circles correspond, lines being circles through infinity."},
		{q: "Under w=1/z, where does infinity go? (enter a number)", a: "0", e: "Large |z| gives small |w|; infinity maps to 0."},
	}
	hard := []qa{
		{q: "Is meromorphic on sphere = rational function? (yes/no)", a: "yes", e: "Meromorphic on compact sphere must be rational."},
		{q: "Does Möbius transformation correspond to rotation of sphere? (yes/no)", a: "yes", e: "Automorphisms of sphere are Möbius."},
		{q: "Is ∞ a point where 1/z has a pole? (yes/no)", a: "yes", e: "Behaviour at ∞ via chart w=1/z."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type cauchyGoursatGen struct{}

func (g *cauchyGoursatGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type qa struct {
		q string
		a string
		e string
	}
	easy := []qa{
		{q: "What is the integral of z^2 around any closed contour? (enter a number)", a: "0", e: "z^2 is entire with antiderivative z^3/3, so closed loops give 0."},
		{q: "What is the integral of 1/z around |z|=1, divided by 2πi? (enter a number)", a: "1", e: "Interior pole of residue 1 gives 2πi; divided by 2πi is 1."},
		{q: "Does Cauchy-Goursat state the integral is 0 for f analytic inside simple closed C? (yes/no)", a: "yes", e: "Integral around null-homotopic loop zero."},
		{q: "Does the theorem give integral zero when f has a pole inside C? (yes/no)", a: "no", e: "The analyticity hypothesis fails; e.g. 1/z gives 2πi."},
		{q: "Did Goursat remove the continuity-of-f' assumption? (yes/no)", a: "yes", e: "Goursat needs only analyticity."},
	}
	hard := []qa{
		{q: "Does Cauchy-Goursat imply existence of antiderivative locally? (yes/no)", a: "yes", e: "Integral independent of path."},
		{q: "Is triangle version used to prove analytic ⇒ infinitely differentiable? (yes/no)", a: "yes", e: "Goursat lemma."},
		{q: "Does proof subdivide triangle into 4 smaller triangles? (yes/no)", a: "yes", e: "Goursat's bisection."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type powerSeriesGen struct{}

func (g *powerSeriesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type qa struct {
		q string
		a string
		e string
	}
	easy := []qa{
		{q: "What is the radius of convergence of 1/(1-z) at 0? (enter a number)", a: "1", e: "Pole at z=1, distance 1 from 0."},
		{q: "What is the z^2 coefficient of e^z at 0? (enter a fraction like 1/2)", a: "1/2", e: "Taylor term is z^2/2, so the coefficient is 1/2."},
		{q: "Does uniform convergence allow termwise differentiation inside the radius? (yes/no)", a: "yes", e: "Inside radius, differentiate term by term."},
		{q: "Does sum z^n converge at z=2? (yes/no)", a: "no", e: "Terms 2^n do not tend to 0; radius is 1."},
		{q: "Is the radius of 1/(1-z) equal to 2? (yes/no)", a: "no", e: "The pole at distance 1 gives radius 1, not 2."},
	}
	hard := []qa{
		{q: "Is power series of e^z entire (radius ∞)? (yes/no)", a: "yes", e: "No finite singularity."},
		{q: "Does Hadamard formula give 1/R=limsup |a_n|^{1/n}? (yes/no)", a: "yes", e: "Root test."},
		{q: "Is analytic continuation possible beyond disk via overlapping disks? (yes/no)", a: "yes", e: "Analytic continuation."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

// ----- complex logarithm -----

type logarithmGen struct{}

func (g *logarithmGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		question string
		answer   string
		reason   string
	}
	tableEasy := []entry{
		{"Does e^w = 0 have solutions? (yes/no)", "no", "The complex exponential never vanishes: |e^w| > 0 always, so 0 has no logarithm."},
		{"What is the principal Log(1)? (enter a number)", "0", "ln|1| + i*Arg(1) = 0 + 0 = 0."},
		{"Is the complex logarithm single-valued without a branch cut? (yes/no)", "no", "Values differ by 2*pi*i: a cut selects one branch."},
		{"Does Log(e^w) = w hold for all w? (yes/no)", "no", "Only on the strip -pi < Im(w) <= pi; otherwise subtract 2*k*pi*i."},
	}
	tableHard := []entry{
		{"Log(-1) + Log(-1) vs Log(1): the correction term is a multiple of pi. Which multiple of pi? (enter like 2pi)", "2pi", "i*pi + i*pi = 2*pi*i vs 0: differs by 2*pi*i."},
		{"Principal Log(-1) equals? (enter like pi*i)", "pi*i", "ln|-1| + i*Arg(-1) = 0 + pi*i."},
	}
	pool := tableEasy
	if scale > 3 {
		pool = append(tableEasy, tableHard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.reason}
}
