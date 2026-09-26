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
	scale := int(1 + ctx.Difficulty*4)
	type qa struct {
		q string
		a string
		e string
	}
	tableEasy := []qa{
		{q: "What is \\(\\arg(1)\\) in degrees? (enter a number)", a: "0", e: "1 lies on the positive real axis, so arg = 0°."},
		{q: "What is \\(\\arg(i)\\) in degrees? (enter a number)", a: "90", e: "i lies on the positive imaginary axis, so arg = 90°."},
		{q: "What is \\(\\arg(-1)\\) in degrees? (enter a number)", a: "180", e: "-1 lies on the negative real axis, so arg = 180°."},
		{q: "What is \\(\\arg(-i)\\) in degrees? (enter a number)", a: "270", e: "-i lies on the negative imaginary axis, so arg = 270° (or -90°)."},
		{q: "Does \\(\\arg(z_1 z_2) = \\arg(z_1)+\\arg(z_2)\\) hold mod \\(360^{\\circ}\\)? (yes/no)", a: "yes", e: "Arguments add when multiplying complex numbers (mod 360°)."},
		{q: "Is \\(\\arg(z)\\) always between 0 and 90 degrees? (yes/no)", a: "no", e: "Arguments span the full circle: arg(-1) = 180°, outside 0-90°."},
	}
	tableHard := []qa{
		{q: "What is \\(\\arg(1+i)\\) in degrees? (enter a number)", a: "45", e: "1+i is in the first quadrant with b/a=1, so arg = 45°."},
		{q: "What is \\(\\arg(-1+i)\\) in degrees? (enter a number)", a: "135", e: "-1+i is in the second quadrant, 180°-45° = 135°."},
		{q: "What is \\(\\arg(1-i)\\) in degrees? (enter a number)", a: "315", e: "1-i is in the fourth quadrant, 360°-45° = 315°."},
	}
	pool := tableEasy
	if scale > 3 {
		pool = append(tableEasy, tableHard...)
	}
	t := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: t.q, Answer: t.a, Explanation: t.e}
}

type eulerIdentityGen struct{}

func (g *eulerIdentityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		q string
		a string
		e string
	}
	tableEasy := []entry{
		{"What is \\(e^{i\\pi} + 1\\)? (enter a number)", "0", "Euler's identity: e^{iπ} + 1 = 0."},
		{"What is \\(e^{i\\pi}\\)?", "-1", "Euler's identity: e^{iπ} = -1."},
		{"Does \\(e^{i\\pi} = -1\\) hold? (yes/no)", "yes", "Yes, Euler's identity gives e^{iπ} = cos π + i sin π = -1."},
		{"What is \\(e^{i\\pi/2}\\)?", "i", "e^{iπ/2} = cos(π/2) + i sin(π/2) = i."},
		{"Is \\(e^{i\\pi} = 1\\)? (yes/no)", "no", "e^{iπ} = -1, not 1."},
	}
	tableHard := []entry{
		{"What is \\(e^{2i\\pi}\\)?", "1", "Full turn: cos(2π) + i sin(2π) = 1."},
		{"What is \\(|e^{i\\theta}|\\) for real \\(\\theta\\)? (enter a number)", "1", "e^{iθ} lies on the unit circle, so its modulus is 1."},
		{"Does \\(e^{i\\theta}\\) stay on the unit circle for real \\(\\theta\\)? (yes/no)", "yes", "|e^{iθ}| = 1 for all real θ."},
	}
	pool := tableEasy
	if scale > 3 {
		pool = append(tableEasy, tableHard...)
	}
	t := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: t.q, Answer: t.a, Explanation: t.e}
}

type cauchyRiemannGen struct{}

func (g *cauchyRiemannGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type qa struct {
		q string
		a string
		e string
	}
	tableEasy := []qa{
		{q: "For \\(f(z)=z^{2}\\) with \\(u=x^{2}-y^{2}\\), what is \\(u_{x}\\) at \\((1,1)\\)? (enter a number)", a: "2", e: "u_x = 2x, so at (1,1) it is 2."},
		{q: "For \\(f(z)=\\bar{z}\\) with \\(u=x\\), what is \\(u_{x}\\)? (enter a number)", a: "1", e: "u = x differentiates to u_x = 1."},
		{q: "For \\(f(z)=z^{2}\\), do the Cauchy-Riemann equations hold everywhere? (yes/no)", a: "yes", e: "f(z)=z² is entire; u=x²-y², v=2xy satisfy u_x=v_y and u_y=-v_x everywhere."},
		{q: "For \\(f(z)=\\bar{z}\\), do the Cauchy-Riemann equations hold? (yes/no)", a: "no", e: "f(z)=x-iy has u_x=1, v_y=-1, so u_x≠v_y; not analytic anywhere."},
		{q: "If \\(f\\) is analytic, must \\(u_x = v_y\\) hold? (yes/no)", a: "yes", e: "This is the first Cauchy-Riemann equation characterizing analytic functions."},
		{q: "For \\(f(z)=z\\) with \\(v=y\\), what is \\(v_{y}\\)? (enter a number)", a: "1", e: "v = y differentiates to v_y = 1."},
	}
	tableHard := []qa{
		{q: "For \\(f(z)=e^{z}\\) with \\(u=e^{x}\\cos y\\), what is \\(u_{x}\\) at \\((0,0)\\)? (enter a number)", a: "1", e: "u_x = e^x cos y, so at (0,0) it is 1."},
		{q: "For \\(f(z)=|z|^{2}\\) with \\(u=x^{2}+y^{2}\\), do CR hold at \\((1,0)\\)? (yes/no)", a: "no", e: "u_x = 2x = 2 but v_y = 0 there, so u_x≠v_y."},
		{q: "Does analyticity on a domain imply CR at every point of it? (yes/no)", a: "yes", e: "Analytic functions satisfy CR everywhere in their domain."},
	}
	pool := tableEasy
	if scale > 3 {
		pool = append(tableEasy, tableHard...)
	}
	t := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: t.q, Answer: t.a, Explanation: t.e}
}

type residueGen struct{}

func (g *residueGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type qa struct {
		q string
		a string
		e string
	}
	tableEasy := []qa{
		{q: "For \\(f(z)=1/z\\), what is \\(\\operatorname{Res}(f,0)\\)? (enter a number)", a: "1", e: "The Laurent coefficient of 1/z at 0 is 1, so Res(f,0)=1."},
		{q: "If \\(f\\) has a simple pole at \\(z_0\\) with numerator \\(g(z_0)\\neq0\\), is \\(\\operatorname{Res}(f,z_0)=g(z_0)/h'(z_0)\\) when \\(f=g/h\\)? (yes/no)", a: "yes", e: "For a simple pole, Res = g(z₀)/h'(z₀)."},
		{q: "Does \\(\\oint_C f(z)dz = 2\\pi i \\sum \\operatorname{Res}(f, z_k)\\) hold for \\(f\\) analytic inside \\(C\\) except at poles? (yes/no)", a: "yes", e: "This is the residue theorem."},
		{q: "For \\(f(z)=1/z^{2}\\), what is \\(\\operatorname{Res}(f,0)\\)? (enter a number)", a: "0", e: "The coefficient of 1/z in 1/z² is 0."},
		{q: "Is the residue of \\(1/z\\) at 0 equal to 0? (yes/no)", a: "no", e: "The residue is 1, the coefficient of 1/z."},
	}
	tableHard := []qa{
		{q: "For \\(f(z)=e^{z}/z\\), what is \\(\\operatorname{Res}(f,0)\\)? (enter a number)", a: "1", e: "e^z/z = 1/z + 1 + ..., so a_{-1} = 1."},
		{q: "For \\(f(z)=1/(z(z-1))\\), what is \\(\\operatorname{Res}(f,0)\\)? (enter a number)", a: "-1", e: "Simple pole with h(z)=z(z-1), h'(0)=-1, so Res = 1/(-1) = -1."},
		{q: "Is the sum of all residues of a rational function on the sphere 0? (yes/no)", a: "yes", e: "Residues including infinity sum to zero."},
	}
	pool := tableEasy
	if scale > 3 {
		pool = append(tableEasy, tableHard...)
	}
	t := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: t.q, Answer: t.a, Explanation: t.e}
}

type analyticGen struct{}

func (g *analyticGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type qa struct {
		q string
		a string
		e string
	}
	easy := []qa{
		{q: "Is f(z)=z^2 analytic everywhere? (yes/no)", a: "yes", e: "Polynomials are entire; CR holds everywhere."},
		{q: "Is f(z)=conj(z) analytic anywhere? (yes/no)", a: "no", e: "Conjugate fails CR everywhere."},
		{q: "For f(z)=z^2, what is f'(0)? (enter a number)", a: "0", e: "f'(z)=2z, so f'(0)=0."},
		{q: "For f(z)=2z+1, what is f'(0)? (enter a number)", a: "2", e: "f'(z)=2 everywhere, so f'(0)=2."},
		{q: "Does analytic imply continuous? (yes/no)", a: "yes", e: "Differentiable implies continuous."},
		{q: "For f(z)=z^3, what is f'(1)? (enter a number)", a: "3", e: "f'(z)=3z^2, so f'(1)=3."},
	}
	hard := []qa{
		{q: "Is f(z)=|z|^2 analytic only at 0? (yes/no)", a: "yes", e: "CR only at origin for |z|^2."},
		{q: "Does analytic on domain imply infinitely differentiable? (yes/no)", a: "yes", e: "Complex analytic ⇒ holomorphic ⇒ C∞."},
		{q: "Is Liouville's theorem: bounded entire ⇒ constant? (yes/no)", a: "yes", e: "Bounded entire functions are constant."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
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

type harmonicGen struct{}

func (g *harmonicGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type qa struct {
		q string
		a string
		e string
	}
	easy := []qa{
		{q: "For u=x^2-y^2, what is u_xx? (enter a number)", a: "2", e: "u_x=2x, so u_xx=2."},
		{q: "For u=xy, what is the Laplacian u_xx+u_yy? (enter a number)", a: "0", e: "u_xx=0 and u_yy=0, so the sum is 0."},
		{q: "Is the real part of an analytic function harmonic? (yes/no)", a: "yes", e: "u = Re f satisfies Laplace Δu=0 via CR."},
		{q: "Is u=x^3 harmonic? (yes/no)", a: "no", e: "u_xx=6x and u_yy=0; the sum 6x is not identically 0."},
		{q: "Is a constant function harmonic? (yes/no)", a: "yes", e: "Δc=0."},
		{q: "For u=2xy, what is u_xy? (enter a number)", a: "2", e: "u_x=2y, so u_xy=2."},
	}
	hard := []qa{
		{q: "Does harmonic conjugate v exist locally for harmonic u? (yes/no)", a: "yes", e: "On simply connected domain, harmonic has conjugate making u+iv analytic."},
		{q: "Does maximum principle hold: harmonic non-constant has no interior max? (yes/no)", a: "yes", e: "Interior max forces constancy."},
		{q: "Is Poisson integral formula for harmonic on disk? (yes/no)", a: "yes", e: "Harmonic extension of boundary values."},
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

type conformalGen struct{}

func (g *conformalGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type qa struct {
		q string
		a string
		e string
	}
	easy := []qa{
		{q: "For w=z^2, what is w'(1)? (enter a number)", a: "2", e: "w'=2z, so w'(1)=2."},
		{q: "For w=e^z, what is w'(0)? (enter a number)", a: "1", e: "w'=e^z, so w'(0)=1."},
		{q: "Is w=z^2 conformal at 0? (yes/no)", a: "no", e: "w'(0)=0 doubles angles there, so not conformal at 0."},
		{q: "Is a Mobius transformation conformal? (yes/no)", a: "yes", e: "Mobius analytic with nonzero derivative."},
		{q: "Is a conformal map angle-preserving where f' is nonzero? (yes/no)", a: "yes", e: "Analytic with nonzero derivative preserves angles."},
		{q: "For w=3z, what is w'(2)? (enter a number)", a: "3", e: "w'=3 everywhere, so w'(2)=3."},
	}
	hard := []qa{
		{q: "Does Riemann mapping map simply connected domain to unit disk conformally? (yes/no)", a: "yes", e: "Riemann mapping theorem."},
		{q: "Is conformal at point equivalent to complex differentiable with nonzero derivative? (yes/no)", a: "yes", e: "Definition."},
		{q: "Does conformal preserve orientation? (yes/no)", a: "yes", e: "Positive Jacobian."},
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

type mobiusTransformGen struct{}

func (g *mobiusTransformGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type qa struct {
		q string
		a string
		e string
	}
	easy := []qa{
		{q: "For w=(z-1)/(z+1), what is w(1)? (enter a number)", a: "0", e: "(1-1)/(1+1)=0."},
		{q: "For w=1/z, what is w(-1)? (enter a number)", a: "-1", e: "1/(-1)=-1."},
		{q: "Is w=1/z a Mobius transformation? (yes/no)", a: "yes", e: "a=0,b=1,c=1,d=0 with ad-bc=-1."},
		{q: "Is a Mobius map fixed by prescribing images of only 2 points? (yes/no)", a: "no", e: "Three points are needed; the action is 3-transitive."},
		{q: "Does a Mobius map send circles and lines to circles and lines? (yes/no)", a: "yes", e: "Circle-preserving on the sphere."},
		{q: "For w=2z, what is w(3)? (enter a number)", a: "6", e: "2 times 3 is 6."},
	}
	hard := []qa{
		{q: "Does group of Möbius is PSL(2,C)? (yes/no)", a: "yes", e: "Projective linear group."},
		{q: "Can any 3 points be mapped to 0,1,∞ by Möbius? (yes/no)", a: "yes", e: "3-transitive."},
		{q: "Is w=1/z Möbius? (yes/no)", a: "yes", e: "a=0,b=1,c=1,d=0."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type schwarzLemmaGen struct{}

func (g *schwarzLemmaGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type qa struct {
		q string
		a string
		e string
	}
	easy := []qa{
		{q: "For f(z)=z/2 on the unit disk, what is |f'(0)|? (enter a fraction like 1/2)", a: "1/2", e: "f'(z)=1/2 everywhere, so |f'(0)|=1/2."},
		{q: "For f(z)=z^2, what is f'(0)? (enter a number)", a: "0", e: "f'(z)=2z vanishes at 0."},
		{q: "Does Schwarz lemma give |f'(0)| at most 1 when f(0)=0 and |f| at most 1? (yes/no)", a: "yes", e: "Derivative bound of the lemma."},
		{q: "Does the lemma apply when f(0)=1? (yes/no)", a: "no", e: "The hypothesis f(0)=0 fails."},
		{q: "Does equality |f(z0)|=|z0| force a rotation? (yes/no)", a: "yes", e: "Equality case is f(z)=e^{iθ}z."},
	}
	hard := []qa{
		{q: "Is Schwarz-Pick: |f(z1)-f(z2)|/|1-\\bar{f(z2)}f(z1)| ≤ |z1-z2|/|1-\\bar{z2}z1|? (yes/no)", a: "yes", e: "Invariant form."},
		{q: "Is the extremal for |f'(0)| at most 1 attained by a rotation? (yes/no)", a: "yes", e: "Rotations give |f'(0)|=1."},
		{q: "Is automorphisms of disk are Möbius (z-a)/(1-\\bar{a}z)? (yes/no)", a: "yes", e: "Disk automorphisms."},
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
