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

type modulusGen struct{}

func (g *modulusGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	entries := []struct{ a, b, r int }{
		{3, 4, 5}, {5, 12, 13}, {8, 6, 10}, {7, 24, 25}, {9, 12, 15},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    fmt.Sprintf("What is \\(|%s|\\)?", fmtComplex(e.a, e.b)),
		Answer:      fmt.Sprintf("%d", e.r),
		Explanation: fmt.Sprintf("|%s| = √(%d²+%d²) = √%d = %d", fmtComplex(e.a, e.b), e.a, e.b, e.a*e.a+e.b*e.b, e.r),
	}
}

type argumentGen struct{}

func (g *argumentGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type qa struct {
		q string
		a string
		e string
	}
	templates := []qa{
		{q: "What is \\(\\arg(1)\\) in degrees? (enter a number)", a: "0", e: "1 lies on the positive real axis, so arg = 0°."},
		{q: "What is \\(\\arg(i)\\) in degrees? (enter a number)", a: "90", e: "i lies on the positive imaginary axis, so arg = 90°."},
		{q: "What is \\(\\arg(-1)\\) in degrees? (enter a number)", a: "180", e: "-1 lies on the negative real axis, so arg = 180°."},
		{q: "What is \\(\\arg(-i)\\) in degrees? (enter a number)", a: "270", e: "-i lies on the negative imaginary axis, so arg = 270° (or -90°)."},
		{q: "Does \\(\\arg(z_1 z_2) = \\arg(z_1)+\\arg(z_2)\\) hold mod \\(360^{\\circ}\\)? (yes/no)", a: "yes", e: "Arguments add when multiplying complex numbers (mod 360°)."},
	}
	t := templates[rand.Intn(len(templates))]
	return generator.Problem{Question: t.q, Answer: t.a, Explanation: t.e}
}

type eulerIdentityGen struct{}

func (g *eulerIdentityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	templates := []struct{ q, a, e string }{
		{"What is \\(e^{i\\pi} + 1\\)? (enter a number)", "0", "Euler's identity: e^{iπ} + 1 = 0."},
		{"What is \\(e^{i\\pi}\\)?", "-1", "Euler's identity: e^{iπ} = -1."},
		{"Does \\(e^{i\\pi} = -1\\) hold? (yes/no)", "yes", "Yes, Euler's identity gives e^{iπ} = cos π + i sin π = -1."},
		{"What is \\(e^{i\\pi/2}\\)?", "i", "e^{iπ/2} = cos(π/2) + i sin(π/2) = i."},
	}
	t := templates[rand.Intn(len(templates))]
	return generator.Problem{Question: t.q, Answer: t.a, Explanation: t.e}
}

type cauchyRiemannGen struct{}

func (g *cauchyRiemannGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type qa struct {
		q string
		a string
		e string
	}
	templates := []qa{
		{q: "For \\(f(z)=z^{2}\\), do the Cauchy-Riemann equations hold everywhere? (yes/no)", a: "yes", e: "f(z)=z² is entire; u=x²-y², v=2xy satisfy u_x=v_y and u_y=-v_x everywhere."},
		{q: "For \\(f(z)=\\bar{z}\\), do the Cauchy-Riemann equations hold? (yes/no)", a: "no", e: "f(z)=x-iy has u_x=1, v_y=-1, so u_x≠v_y; not analytic anywhere."},
		{q: "If \\(f\\) is analytic, must \\(u_x = v_y\\) and \\(u_y = -v_x\\) hold? (yes/no)", a: "yes", e: "These are the Cauchy-Riemann equations characterizing analytic functions."},
		{q: "Does analyticity imply the Cauchy-Riemann equations? (yes/no)", a: "yes", e: "Yes, analytic functions satisfy CR everywhere in their domain."},
	}
	t := templates[rand.Intn(len(templates))]
	return generator.Problem{Question: t.q, Answer: t.a, Explanation: t.e}
}

type residueGen struct{}

func (g *residueGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type qa struct {
		q string
		a string
		e string
	}
	templates := []qa{
		{q: "For \\(f(z)=1/z\\), what is \\(\\operatorname{Res}(f,0)\\)? (enter a number)", a: "1", e: "The Laurent coefficient of 1/z at 0 is 1, so Res(f,0)=1."},
		{q: "If \\(f\\) has a simple pole at \\(z_0\\) with numerator \\(g(z_0)\\neq0\\), is \\(\\operatorname{Res}(f,z_0)=g(z_0)/h'(z_0)\\) when \\(f=g/h\\)? (yes/no)", a: "yes", e: "For a simple pole, Res = g(z₀)/h'(z₀)."},
		{q: "Does \\(\\oint_C f(z)dz = 2\\pi i \\sum \\operatorname{Res}(f, z_k)\\) hold for \\(f\\) analytic inside \\(C\\) except at poles? (yes/no)", a: "yes", e: "This is the residue theorem."},
		{q: "For \\(f(z)=1/z^{2}\\), what is \\(\\operatorname{Res}(f,0)\\)? (enter a number)", a: "0", e: "The coefficient of 1/z in 1/z² is 0."},
	}
	t := templates[rand.Intn(len(templates))]
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
		{q: "Does analytic imply continuous? (yes/no)", a: "yes", e: "Differentiable implies continuous."},
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
		{q: "Does 1/z have Laurent series ∑_{n=-1}^{∞} 0·z^n with a_{-1}=1 at 0? (yes/no)", a: "yes", e: "1/z = z^{-1} is its own Laurent series; principal part is 1/z."},
		{q: "Is annulus 0<|z|<1 the domain for Laurent of 1/(z(1-z))? (yes/no)", a: "yes", e: "Poles at 0,1; annuli determined by radii to nearest singularities."},
		{q: "Does Laurent include negative powers for poles/essential singularities? (yes/no)", a: "yes", e: "Taylor has only ≥0; Laurent adds principal part."},
	}
	hard := []qa{
		{q: "Is classification: pole (finite principal part), essential (infinite), removable (none)? (yes/no)", a: "yes", e: "Laurent principal part classifies singularities."},
		{q: "Does e^{1/z} have infinitely many negative terms at 0? (yes/no)", a: "yes", e: "Essential singularity: Laurent has infinite principal part."},
		{q: "Is coefficient a_{-1} the residue? (yes/no)", a: "yes", e: "Residue is Laurent a_{-1}."},
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
		{q: "Does ∮_{|z|=1} dz/z = 2πi? (yes/no)", a: "yes", e: "Residue 1 at 0 inside, contour integral 2πi."},
		{q: "Is contour integral independent of path in simply connected analytic domain? (yes/no)", a: "yes", e: "Cauchy's theorem: integral around null-homotopic loop is 0."},
		{q: "Does ML-estimate bound |∮ f| ≤ M·L? (yes/no)", a: "yes", e: "ML inequality: sup |f| times length."},
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
		{q: "Is real part of analytic function harmonic? (yes/no)", a: "yes", e: "u = Re f satisfies Laplace Δu=0 via CR."},
		{q: "Does harmonic mean value property hold: u(a)=avg on circle? (yes/no)", a: "yes", e: "Mean value characterizes harmonic."},
		{q: "Is constant function harmonic? (yes/no)", a: "yes", e: "Δc=0."},
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
		{q: "Is Riemann sphere C∪{∞} via stereographic projection? (yes/no)", a: "yes", e: "Sphere S^2 minus north pole ≅ C; add ∞ at pole."},
		{q: "Does stereographic map circles on sphere to circles/lines in C? (yes/no)", a: "yes", e: "Circles correspond."},
		{q: "Is extended complex plane compact? (yes/no)", a: "yes", e: "One-point compactification of C."},
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
		{q: "Is conformal map angle-preserving where f'≠0? (yes/no)", a: "yes", e: "Analytic with nonzero derivative preserves angles."},
		{q: "Does w=z^2 double angles at 0? (yes/no)", a: "yes", e: "Derivative 0 at 0, not conformal there."},
		{q: "Is Möbius transformation conformal? (yes/no)", a: "yes", e: "Möbius analytic with nonzero derivative."},
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
		{q: "Does Cauchy-Goursat state ∮_C f=0 for f analytic inside simple closed C? (yes/no)", a: "yes", e: "Integral around null-homotopic loop zero."},
		{q: "Does Goursat improve Cauchy by removing continuity of f'? (yes/no)", a: "yes", e: "Goursat needs only analyticity."},
		{q: "Is f(z)=1/z analytic inside |z|=1 except at 0, so theorem fails? (yes/no)", a: "yes", e: "Singularity inside, integral 2πi."},
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
		{q: "Does analytic function have power series converging in disk to nearest singularity? (yes/no)", a: "yes", e: "Radius = distance to singularity."},
		{q: "Is radius of 1/(1-z) equal to 1? (yes/no)", a: "yes", e: "Singularity at z=1, distance 1 from 0."},
		{q: "Does uniform convergence allow termwise differentiation? (yes/no)", a: "yes", e: "Inside radius."},
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
		{q: "Is Möbius map w=(az+b)/(cz+d) with ad-bc≠0 bijective on Riemann sphere? (yes/no)", a: "yes", e: "Automorphism of sphere."},
		{q: "Does Möbius map circles/lines to circles/lines? (yes/no)", a: "yes", e: "Circle-preserving."},
		{q: "Is cross-ratio preserved by Möbius? (yes/no)", a: "yes", e: "Möbius preserves cross-ratio."},
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
		{q: "Does Schwarz lemma state: if f analytic on unit disk, f(0)=0, |f(z)|≤1 then |f(z)|≤|z|? (yes/no)", a: "yes", e: "Schwarz lemma."},
		{q: "Does equality |f(z0)|=|z0| imply f(z)=e^{iθ}z? (yes/no)", a: "yes", e: "Rotation."},
		{q: "Does Schwarz give |f'(0)|≤1? (yes/no)", a: "yes", e: "Derivative bound."},
	}
	hard := []qa{
		{q: "Is Schwarz-Pick: |f(z1)-f(z2)|/|1-\\bar{f(z2)}f(z1)| ≤ |z1-z2|/|1-\\bar{z2}z1|? (yes/no)", a: "yes", e: "Invariant form."},
		{q: "Does Schwarz lemma prove fundamental theorem of algebra? (no)", a: "no", e: "Different."},
		{q: "Is automorphisms of disk are Möbius (z-a)/(1-\\bar{a}z)? (yes/no)", a: "yes", e: "Disk automorphisms."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}
