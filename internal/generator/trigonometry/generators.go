package trigonometry

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/mathutil"
)

func Register(reg *generator.Registry) {
	reg.Register("trig.basics.radians", &radiansGen{})
	reg.Register("trig.basics.unit_circle", &unitCircleGen{})
	reg.Register("trig.basics.sin_cos_def", &sinCosDefGen{})
	reg.Register("trig.basics.tan_def", &tanDefGen{})
	reg.Register("trig.basics.reciprocal", &reciprocalGen{})
	reg.Register("trig.ident.pythagorean", &pythagoreanIDGen{})
	reg.Register("trig.basics.special_angles", &specialAnglesGen{})
	reg.Register("trig.basics.reference_angle", &referenceAngleGen{})
	reg.Register("trig.graph.sin", &graphSinGen{})
	reg.Register("trig.graph.cos", &graphCosGen{})
	reg.Register("trig.graph.period", &periodGen{})
	reg.Register("trig.adv.inverse", &inverseGen{})
	reg.Register("trig.adv.law_sines", &lawSinesGen{})
	reg.Register("trig.adv.law_cosines", &lawCosinesGen{})

	reg.Register("trig.adv.arctan", &arctanGen{})
	reg.Register("trig.basics.right_triangle", &rightTriangleGen{})
	reg.Register("trig.eq.basic", &trigEqBasicGen{})
	reg.Register("trig.eq.homogeneous", &trigEqHomogeneousGen{})
	reg.Register("trig.hyperbolic.sinh_cosh", &sinhCoshGen{})
	reg.Register("trig.hyperbolic.tanh_coth", &tanhCothGen{})
	reg.Register("trig.ident.identities", &trigIdentGen{})
	reg.Register("trig.ineq.basic", &trigIneqGen{})
}

type radiansGen struct{}

func (g *radiansGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		deg      int
		radHTML  string
		radLatex string
	}
	table := []entry{
		{30, "\\(\\pi/6\\)", "\\pi/6"},
		{45, "\\(\\pi/4\\)", "\\pi/4"},
		{60, "\\(\\pi/3\\)", "\\pi/3"},
		{90, "\\(\\pi/2\\)", "\\pi/2"},
		{180, "\\(\\pi\\)", "\\pi"},
		{270, "\\(3\\pi/2\\)", "3\\pi/2"},
		{360, "\\(2\\pi\\)", "2\\pi"},
	}
	e := table[rand.Intn(len(table))]
	if rand.Intn(2) == 0 {
		return generator.Problem{
			Question:    fmt.Sprintf("Convert \\(%d^{\\circ}\\) to radians.", e.deg),
			Answer:      e.radHTML,
			Explanation: fmt.Sprintf("\\(%d^{\\circ} \\times \\frac{\\pi}{180} = %s\\).", e.deg, e.radHTML),
		}
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Convert %s to degrees.", e.radHTML),
		Answer:      fmt.Sprintf("%d", e.deg),
		Explanation: fmt.Sprintf("%s \\(\\times \\frac{180}{\\pi} = %d^{\\circ}\\).", e.radHTML, e.deg),
	}
}

type unitCircleGen struct{}

var unitCircleAngles = []struct {
	label string
	sin   string
	cos   string
	easy  bool
}{
		{"\\(0^{\\circ}\\)", "0", "1", true},
	{"\\(30^{\\circ}\\)", "\\(1/2\\)", "\\(\\sqrt{3}/2\\)", true},
	{"\\(45^{\\circ}\\)", "\\(\\sqrt{2}/2\\)", "\\(\\sqrt{2}/2\\)", true},
	{"\\(60^{\\circ}\\)", "\\(\\sqrt{3}/2\\)", "\\(1/2\\)", true},
	{"\\(90^{\\circ}\\)", "1", "0", true},
	{"\\(180^{\\circ}\\)", "0", "-1", true},
	{"\\(270^{\\circ}\\)", "-1", "0", true},
	{"\\(360^{\\circ}\\)", "0", "1", true},
	{"\\(120^{\\circ}\\)", "\\(\\sqrt{3}/2\\)", "\\(-1/2\\)", false},
	{"\\(135^{\\circ}\\)", "\\(\\sqrt{2}/2\\)", "\\(-\\sqrt{2}/2\\)", false},
	{"\\(150^{\\circ}\\)", "\\(1/2\\)", "\\(-\\sqrt{3}/2\\)", false},
	{"\\(210^{\\circ}\\)", "\\(-1/2\\)", "\\(-\\sqrt{3}/2\\)", false},
	{"\\(225^{\\circ}\\)", "\\(-\\sqrt{2}/2\\)", "\\(-\\sqrt{2}/2\\)", false},
	{"\\(240^{\\circ}\\)", "\\(-\\sqrt{3}/2\\)", "\\(-1/2\\)", false},
	{"\\(300^{\\circ}\\)", "\\(-\\sqrt{3}/2\\)", "\\(1/2\\)", false},
	{"\\(315^{\\circ}\\)", "\\(-\\sqrt{2}/2\\)", "\\(\\sqrt{2}/2\\)", false},
	{"\\(330^{\\circ}\\)", "\\(-1/2\\)", "\\(\\sqrt{3}/2\\)", false},
}

func (g *unitCircleGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// Easy → pick from quadrant 1 + axes; Hard → pick from any quadrant
	var candidates []struct{ label, sin, cos string }
	for _, a := range unitCircleAngles {
		if ctx.Difficulty > 0.4 || a.easy {
			candidates = append(candidates, struct{ label, sin, cos string }{a.label, a.sin, a.cos})
		}
	}
	a := candidates[rand.Intn(len(candidates))]
	if rand.Intn(2) == 0 {
		return generator.Problem{
			Question:    fmt.Sprintf("What is \\(\\sin(%s)\\)?", a.label),
			Answer:      a.sin,
			Explanation: fmt.Sprintf("On the unit circle, \\(\\sin(%s) = %s\\).", a.label, a.sin),
		}
	}
	return generator.Problem{
		Question:    fmt.Sprintf("What is \\(\\cos(%s)\\)?", a.label),
		Answer:      a.cos,
		Explanation: fmt.Sprintf("On the unit circle, \\(\\cos(%s) = %s\\).", a.label, a.cos),
	}
}

type sinCosDefGen struct{}

type triple struct{ opp, adj, hyp int }

var pythagoreanTriples = []triple{
	{3, 4, 5},
	{5, 12, 13},
	{8, 15, 17},
	{7, 24, 25},
	{9, 40, 41},
	{6, 8, 10},
	{10, 24, 26},
	{12, 16, 20},
}

var easySinTriples = []triple{
	{3, 4, 5},
	{6, 8, 10},
	{5, 12, 13},
	{9, 12, 15},
}

var hardSinTriples = []triple{
	{8, 15, 17},
	{7, 24, 25},
	{9, 40, 41},
	{10, 24, 26},
	{12, 16, 20},
	{15, 8, 17},
	{24, 7, 25},
	{40, 9, 41},
}

func (g *sinCosDefGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	var triples []triple
	if ctx.Difficulty > 0.5 {
		triples = hardSinTriples
	} else {
		triples = easySinTriples
	}
	if ctx.Difficulty > 0.7 && rand.Intn(2) == 0 {
		triples = append(triples, easySinTriples...)
	}
	t := triples[rand.Intn(len(triples))]
	if rand.Intn(2) == 0 {
		gcd := mathutil.GCD(t.opp, t.hyp)
		num := t.opp / gcd
		den := t.hyp / gcd
		return generator.Problem{
			Question:    fmt.Sprintf("In a right triangle with opposite = %d and hypotenuse = %d, what is \\(\\sin(\\theta)\\)?", t.opp, t.hyp),
			Answer:      fmt.Sprintf("%d/%d", num, den),
			Explanation: fmt.Sprintf("sin(θ) = opp/hyp = %d/%d = %d/%d.", t.opp, t.hyp, num, den),
		}
	}
	gcd := mathutil.GCD(t.adj, t.hyp)
	num := t.adj / gcd
	den := t.hyp / gcd
	return generator.Problem{
		Question:    fmt.Sprintf("In a right triangle with adjacent = %d and hypotenuse = %d, what is \\(\\cos(\\theta)\\)?", t.adj, t.hyp),
		Answer:      fmt.Sprintf("%d/%d", num, den),
		Explanation: fmt.Sprintf("cos(θ) = adj/hyp = %d/%d = %d/%d.", t.adj, t.hyp, num, den),
	}
}

type tanDefGen struct{}

func (g *tanDefGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// Half the time use right triangle sides, half use sin/cos ratio
	if rand.Intn(2) == 0 {
		t := pythagoreanTriples[rand.Intn(len(pythagoreanTriples))]
		gcd := mathutil.GCD(t.opp, t.adj)
		num := t.opp / gcd
		den := t.adj / gcd
		return generator.Problem{
			Question:    fmt.Sprintf("In a right triangle with opposite = %d and adjacent = %d, what is \\(\\tan(\\theta)\\)?", t.opp, t.adj),
			Answer:      fmt.Sprintf("%d/%d", num, den),
			Explanation: fmt.Sprintf("tan(θ) = opp/adj = %d/%d = %d/%d.", t.opp, t.adj, num, den),
		}
	}
	t := pythagoreanTriples[rand.Intn(len(pythagoreanTriples))]
	gcdOpp := mathutil.GCD(t.opp, t.hyp)
	gcdAdj := mathutil.GCD(t.adj, t.hyp)
	sinNum := t.opp / gcdOpp
	sinDen := t.hyp / gcdOpp
	cosNum := t.adj / gcdAdj
	cosDen := t.hyp / gcdAdj
	tanGCD := mathutil.GCD(sinNum*cosDen, sinDen*cosNum)
	tanNum := sinNum * cosDen / tanGCD
	tanDen := sinDen * cosNum / tanGCD
	return generator.Problem{
		Question:    fmt.Sprintf("If \\(\\sin(\\theta) = %d/%d\\) and \\(\\cos(\\theta) = %d/%d\\), what is \\(\\tan(\\theta)\\)?", sinNum, sinDen, cosNum, cosDen),
		Answer:      fmt.Sprintf("%d/%d", tanNum, tanDen),
		Explanation: fmt.Sprintf("tan(θ) = sin(θ)/cos(θ) = (%d/%d)/(%d/%d) = %d/%d.", sinNum, sinDen, cosNum, cosDen, tanNum, tanDen),
	}
}

type reciprocalGen struct{}

func (g *reciprocalGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	t := pythagoreanTriples[rand.Intn(len(pythagoreanTriples))]
	funcs := []struct {
		name  string
		num   int
		den   int
		recip string
	}{
		{"\\(\\sin(\\theta)\\)", t.opp, t.hyp, "\\(\\csc(\\theta)\\)"},
		{"\\(\\cos(\\theta)\\)", t.adj, t.hyp, "\\(\\sec(\\theta)\\)"},
		{"\\(\\tan(\\theta)\\)", t.opp, t.adj, "\\(\\cot(\\theta)\\)"},
	}
	fn := funcs[rand.Intn(3)]
	gcd := mathutil.GCD(fn.den, fn.num)
	recipNum := fn.den / gcd
	recipDen := fn.num / gcd
	return generator.Problem{
		Question:    fmt.Sprintf("If %s = %d/%d, what is %s?", fn.name, fn.num, fn.den, fn.recip),
		Answer:      fmt.Sprintf("%d/%d", recipNum, recipDen),
		Explanation: fmt.Sprintf("%s = 1/%s = 1/(%d/%d) = %d/%d.", fn.recip, fn.name, fn.num, fn.den, recipNum, recipDen),
	}
}

type pythagoreanIDGen struct{}

func (g *pythagoreanIDGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	t := pythagoreanTriples[rand.Intn(len(pythagoreanTriples))]
	// Randomly ask for sin from cos or cos from sin
	if rand.Intn(2) == 0 {
		sinGcd := mathutil.GCD(t.opp, t.hyp)
		sinNum := t.opp / sinGcd
		sinDen := t.hyp / sinGcd
		cosGcd := mathutil.GCD(t.adj, t.hyp)
		cosNum := t.adj / cosGcd
		cosDen := t.hyp / cosGcd
		return generator.Problem{
			Question:    fmt.Sprintf("If \\(\\sin(\\theta) = %d/%d\\), find \\(\\cos(\\theta)\\) using \\(\\sin^{2}\\theta + \\cos^{2}\\theta = 1\\) (assume \\(\\theta\\) is acute).", sinNum, sinDen),
			Answer:      fmt.Sprintf("%d/%d", cosNum, cosDen),
			Explanation: fmt.Sprintf("cos²θ = 1 - sin²θ = 1 - (%d/%d)² = 1 - %d/%d = %d/%d, so cos(θ) = √(%d/%d) = %d/%d.", sinNum, sinDen, sinNum*sinNum, sinDen*sinDen, sinDen*sinDen-sinNum*sinNum, sinDen*sinDen, sinDen*sinDen-sinNum*sinNum, sinDen*sinDen, cosNum, cosDen),
		}
	}
	cosGcd := mathutil.GCD(t.adj, t.hyp)
	cosNum := t.adj / cosGcd
	cosDen := t.hyp / cosGcd
	sinGcd := mathutil.GCD(t.opp, t.hyp)
	sinNum := t.opp / sinGcd
	sinDen := t.hyp / sinGcd
	return generator.Problem{
		Question:    fmt.Sprintf("If \\(\\cos(\\theta) = %d/%d\\), find \\(\\sin(\\theta)\\) using \\(\\sin^{2}\\theta + \\cos^{2}\\theta = 1\\) (assume \\(\\theta\\) is acute).", cosNum, cosDen),
		Answer:      fmt.Sprintf("%d/%d", sinNum, sinDen),
		Explanation: fmt.Sprintf("sin²θ = 1 - cos²θ = 1 - (%d/%d)² = 1 - %d/%d = %d/%d, so sin(θ) = √(%d/%d) = %d/%d.", cosNum, cosDen, cosNum*cosNum, cosDen*cosDen, cosDen*cosDen-cosNum*cosNum, cosDen*cosDen, cosDen*cosDen-cosNum*cosNum, cosDen*cosDen, sinNum, sinDen),
	}
}

type specialAnglesGen struct{}

func (g *specialAnglesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		label string
		sin   string
		cos   string
		tan   string
	}
	table := []entry{
		{"\\(30^{\\circ}\\)", "\\(\\frac{1}{2}\\)", "\\(\\frac{\\sqrt{3}}{2}\\)", "\\(\\frac{\\sqrt{3}}{3}\\)"},
		{"\\(45^{\\circ}\\)", "\\(\\frac{\\sqrt{2}}{2}\\)", "\\(\\frac{\\sqrt{2}}{2}\\)", "1"},
		{"\\(60^{\\circ}\\)", "\\(\\frac{\\sqrt{3}}{2}\\)", "\\(\\frac{1}{2}\\)", "\\(\\sqrt{3}\\)"},
	}
	e := table[rand.Intn(3)]
	funcs := []struct {
		name string
		val  string
	}{
		{"sin", e.sin},
		{"cos", e.cos},
		{"tan", e.tan},
	}
	fn := funcs[rand.Intn(3)]
	return generator.Problem{
		Question:    fmt.Sprintf("What is the exact value of %s(%s)?", fn.name, e.label),
		Answer:      fn.val,
		Explanation: fmt.Sprintf("%s(%s) = %s.", fn.name, e.label, fn.val),
	}
}

type referenceAngleGen struct{}

func (g *referenceAngleGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		angle int
		ref   int
	}
	table := []entry{
		{120, 60},
		{135, 45},
		{150, 30},
		{210, 30},
		{225, 45},
		{240, 60},
		{300, 60},
		{315, 45},
		{330, 30},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{
		Question:    fmt.Sprintf("What is the reference angle for %d°?", e.angle),
		Answer:      fmt.Sprintf("%d°", e.ref),
		Explanation: fmt.Sprintf("The reference angle for %d° is the acute angle it makes with the x-axis: %d°.", e.angle, e.ref),
	}
}

type graphSinGen struct{}

func (g *graphSinGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	questions := []struct {
		q, a, e string
	}{
		{"What is \\(\\sin(0)\\)?", "0", "\\(\\sin(0) = 0\\)"},
		{"What is the maximum value of \\(y = \\sin(x)\\)?", "1", "The sine function has a maximum value of 1."},
		{"What is the minimum value of \\(y = \\sin(x)\\)?", "-1", "The sine function has a minimum value of -1."},
		{"What is the period of \\(y = \\sin(x)\\)?", "\\(360^{\\circ}\\)", "The period of \\(\\sin(x)\\) is \\(360^{\\circ}\\) (\\(2\\pi\\) radians)."},
		{"What is the amplitude of \\(y = \\sin(x)\\)?", "1", "The amplitude of \\(\\sin(x)\\) is 1."},
		{"What is \\(\\sin(90^{\\circ})\\)?", "1", "\\(\\sin(90^{\\circ}) = 1\\)"},
		{"What is \\(\\sin(180^{\\circ})\\)?", "0", "\\(\\sin(180^{\\circ}) = 0\\)"},
		{"What is \\(\\sin(270^{\\circ})\\)?", "-1", "\\(\\sin(270^{\\circ}) = -1\\)"},
		{"What is \\(\\sin(360^{\\circ})\\)?", "0", "\\(\\sin(360^{\\circ}) = 0\\)"},
		{"Where does \\(\\sin(x)\\) start on the \\(y\\)-axis?", "0", "\\(\\sin(0) = 0\\)"},
	}
	q := questions[rand.Intn(len(questions))]
	return generator.Problem{
		Question:    q.q,
		Answer:      q.a,
		Explanation: q.e,
	}
}

type graphCosGen struct{}

func (g *graphCosGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	questions := []struct {
		q, a, e string
	}{
		{"What is \\(\\cos(0)\\)?", "1", "\\(\\cos(0) = 1\\)"},
		{"What is the maximum value of \\(y = \\cos(x)\\)?", "1", "The cosine function has a maximum value of 1."},
		{"What is the minimum value of \\(y = \\cos(x)\\)?", "-1", "The cosine function has a minimum value of -1."},
		{"What is the period of \\(y = \\cos(x)\\)?", "\\(360^{\\circ}\\)", "The period of \\(\\cos(x)\\) is \\(360^{\\circ}\\) (\\(2\\pi\\) radians)."},
		{"What is the amplitude of \\(y = \\cos(x)\\)?", "1", "The amplitude of \\(\\cos(x)\\) is 1."},
		{"What is \\(\\cos(90^{\\circ})\\)?", "0", "\\(\\cos(90^{\\circ}) = 0\\)"},
		{"What is \\(\\cos(180^{\\circ})\\)?", "-1", "\\(\\cos(180^{\\circ}) = -1\\)"},
		{"What is \\(\\cos(270^{\\circ})\\)?", "0", "\\(\\cos(270^{\\circ}) = 0\\)"},
		{"What is \\(\\cos(360^{\\circ})\\)?", "1", "\\(\\cos(360^{\\circ}) = 1\\)"},
		{"Where does \\(\\cos(x)\\) start on the \\(y\\)-axis?", "1", "\\(\\cos(0) = 1\\)"},
	}
	q := questions[rand.Intn(len(questions))]
	return generator.Problem{
		Question:    q.q,
		Answer:      q.a,
		Explanation: q.e,
	}
}

type periodGen struct{}

func (g *periodGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	configs := []struct {
		A      int
		B      int
		C      int
		amp    int
		period string
		shift  string
	}{
		{2, 1, 0, 2, "360°", "0"},
		{3, 2, 0, 3, "180°", "0"},
		{1, 3, 0, 1, "120°", "0"},
		{4, 4, 0, 4, "90°", "0"},
		{2, 1, 45, 2, "360°", "45°"},
		{3, 2, 30, 3, "180°", "30°"},
		{1, 3, 60, 1, "120°", "60°"},
	}
	c := configs[rand.Intn(len(configs))]
	opts := []struct {
		q, a, e string
	}{
		{
			fmt.Sprintf("What is the amplitude of \\(y = %d \\sin(%dx)\\)?", c.A, c.B),
			fmt.Sprintf("%d", c.amp),
			fmt.Sprintf("Amplitude = \\(|A| = |%d| = %d\\).", c.A, c.amp),
		},
		{
			fmt.Sprintf("What is the period of \\(y = \\sin(%dx)\\)?", c.B),
			c.period,
			fmt.Sprintf("Period = \\(360^{\\circ}/|B| = 360^{\\circ}/%d = %s\\).", c.B, c.period),
		},
		{
			fmt.Sprintf("What is the phase shift of \\(y = \\sin(%dx + %d)\\)?", c.B, c.C),
			c.shift,
			fmt.Sprintf("Phase shift = \\(-C/B = -%d/%d = -%s\\) (shifted right).", c.C, c.B, c.shift),
		},
	}
	opt := opts[rand.Intn(3)]
	return generator.Problem{
		Question:    opt.q,
		Answer:      opt.a,
		Explanation: opt.e,
	}
}

type inverseGen struct{}

func (g *inverseGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		f      string
		val    string
		result string
		desc   string
	}
	table := []entry{
		{"\\(\\arcsin(0)\\)", "0", "\\(0^{\\circ}\\)", "\\(\\arcsin(0) = 0^{\\circ}\\)"},
		{"\\(\\arcsin(1/2)\\)", "1/2", "\\(30^{\\circ}\\)", "\\(\\arcsin(1/2) = 30^{\\circ}\\)"},
		{"\\(\\arcsin(\\sqrt{2}/2)\\)", "\\(\\sqrt{2}/2\\)", "\\(45^{\\circ}\\)", "\\(\\arcsin(\\sqrt{2}/2) = 45^{\\circ}\\)"},
		{"\\(\\arcsin(\\sqrt{3}/2)\\)", "\\(\\sqrt{3}/2\\)", "\\(60^{\\circ}\\)", "\\(\\arcsin(\\sqrt{3}/2) = 60^{\\circ}\\)"},
		{"\\(\\arcsin(1)\\)", "1", "\\(90^{\\circ}\\)", "\\(\\arcsin(1) = 90^{\\circ}\\)"},
		{"\\(\\arccos(1)\\)", "1", "\\(0^{\\circ}\\)", "\\(\\arccos(1) = 0^{\\circ}\\)"},
		{"\\(\\arccos(\\sqrt{3}/2)\\)", "\\(\\sqrt{3}/2\\)", "\\(30^{\\circ}\\)", "\\(\\arccos(\\sqrt{3}/2) = 30^{\\circ}\\)"},
		{"\\(\\arccos(\\sqrt{2}/2)\\)", "\\(\\sqrt{2}/2\\)", "\\(45^{\\circ}\\)", "\\(\\arccos(\\sqrt{2}/2) = 45^{\\circ}\\)"},
		{"\\(\\arccos(1/2)\\)", "1/2", "\\(60^{\\circ}\\)", "\\(\\arccos(1/2) = 60^{\\circ}\\)"},
		{"\\(\\arccos(0)\\)", "0", "\\(90^{\\circ}\\)", "\\(\\arccos(0) = 90^{\\circ}\\)"},
		{"\\(\\arctan(0)\\)", "0", "\\(0^{\\circ}\\)", "\\(\\arctan(0) = 0^{\\circ}\\)"},
		{"\\(\\arctan(\\sqrt{3}/3)\\)", "\\(\\sqrt{3}/3\\)", "\\(30^{\\circ}\\)", "\\(\\arctan(\\sqrt{3}/3) = 30^{\\circ}\\)"},
		{"\\(\\arctan(1)\\)", "1", "\\(45^{\\circ}\\)", "\\(\\arctan(1) = 45^{\\circ}\\)"},
		{"\\(\\arctan(\\sqrt{3})\\)", "\\(\\sqrt{3}\\)", "\\(60^{\\circ}\\)", "\\(\\arctan(\\sqrt{3}) = 60^{\\circ}\\)"},
		{"\\(\\arcsin(-1/2)\\)", "-1/2", "\\(-30^{\\circ}\\)", "\\(\\arcsin(-1/2) = -30^{\\circ}\\)"},
		{"\\(\\arccos(-1/2)\\)", "-1/2", "\\(120^{\\circ}\\)", "\\(\\arccos(-1/2) = 120^{\\circ}\\)"},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{
		Question:    fmt.Sprintf("What is %s?", e.f),
		Answer:      e.result,
		Explanation: e.desc,
	}
}

type lawSinesGen struct{}

func (g *lawSinesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type triangle struct {
		A, B int
		a    int
	}
	configs := []triangle{
		{30, 45, 10},
		{30, 60, 12},
		{45, 60, 8},
		{30, 30, 15},
		{45, 30, 10},
		{60, 30, 14},
		{30, 90, 10},
		{45, 90, 12},
	}
	c := configs[rand.Intn(len(configs))]
	sinA := math.Sin(float64(c.A) * math.Pi / 180)
	sinB := math.Sin(float64(c.B) * math.Pi / 180)
	b := float64(c.a) * sinB / sinA
	bRounded := math.Round(b*10) / 10
	aDeg := c.A
	bDeg := c.B
	return generator.Problem{
		Question:    fmt.Sprintf("In triangle ABC, \\(A = %d^{\\circ}\\), \\(B = %d^{\\circ}\\), and side \\(a = %d\\). Find side \\(b\\) using the law of sines.", aDeg, bDeg, c.a),
		Answer:      fmt.Sprintf("%.1f", bRounded),
		Explanation: fmt.Sprintf("a/sin(A) = b/sin(B), so b = a·sin(B)/sin(A) = %d·sin(%d°)/sin(%d°) = %.1f.", c.a, bDeg, aDeg, bRounded),
	}
}

type lawCosinesGen struct{}

func (g *lawCosinesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type triangle struct {
		a, b int
		C    int
		c2   int
	}
	configs := []triangle{
		{5, 6, 60, 31},
		{3, 4, 60, 13},
		{7, 8, 60, 57},
		{5, 5, 60, 25},
		{6, 7, 60, 43},
		{4, 9, 60, 61},
		{8, 10, 60, 84},
		{5, 8, 60, 49},
	}
	c := configs[rand.Intn(len(configs))]
	cVal := math.Sqrt(float64(c.c2))
	cRounded := math.Round(cVal*10) / 10
	return generator.Problem{
		Question:    fmt.Sprintf("In triangle ABC, \\(a = %d\\), \\(b = %d\\), and \\(C = %d^{\\circ}\\). Find side \\(c\\) using the law of cosines.", c.a, c.b, c.C),
		Answer:      fmt.Sprintf("%.1f", cRounded),
		Explanation: fmt.Sprintf("c² = a² + b² - 2ab·cos(C) = %d² + %d² - 2(%d)(%d)·cos(%d°) = %d + %d - %d = %d, so c = √%d ≈ %.1f.", c.a, c.b, c.a, c.b, c.C, c.a*c.a, c.b*c.b, c.a*c.b, c.c2, c.c2, cRounded),
	}
}

// arctan: compute arctan values and properties
type arctanGen struct{}

func (g *arctanGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"What is \\(\\arctan(1)\\) in degrees?", "45", "\\(\\tan(45^{\\circ}) = 1\\), so \\(\\arctan(1) = 45^{\\circ}\\)."},
		{"What is \\(\\arctan(0)\\) in degrees?", "0", "\\(\\tan(0^{\\circ}) = 0\\), so \\(\\arctan(0) = 0^{\\circ}\\)."},
		{"What is \\(\\arctan(\\sqrt{3})\\) in degrees?", "60", "\\(\\tan(60^{\\circ}) = \\sqrt{3}\\), so \\(\\arctan(\\sqrt{3}) = 60^{\\circ}\\)."},
		{"What is \\(\\arctan(1/\\sqrt{3})\\) in degrees?", "30", "\\(\\tan(30^{\\circ}) = 1/\\sqrt{3}\\), so \\(\\arctan(1/\\sqrt{3}) = 30^{\\circ}\\)."},
		{"What is the range of \\(\\arctan(x)\\)? (in degrees)", "-90 to 90", "\\(\\arctan(x)\\) returns values in \\((-90^{\\circ}, 90^{\\circ})\\)."},
		{"Is \\(\\arctan(x)\\) an odd function? (yes/no)", "yes", "\\(\\arctan(-x) = -\\arctan(x)\\), so it is odd."},
		{"As \\(x \\to \\infty\\), \\(\\arctan(x)\\) approaches what value in degrees?", "90", "\\(\\lim_{x\\to\\infty} \\arctan(x) = 90^{\\circ}\\)."},
		{"As \\(x \\to -\\infty\\), \\(\\arctan(x)\\) approaches what value in degrees?", "-90", "\\(\\lim_{x\\to-\\infty} \\arctan(x) = -90^{\\circ}\\)."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

// right triangle trig: SOH CAH TOA
type rightTriangleGen struct{}

func (g *rightTriangleGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type triple struct{ opp, adj, hyp int }
	triples := []triple{
		{3, 4, 5}, {4, 3, 5},
		{5, 12, 13}, {12, 5, 13},
		{6, 8, 10}, {8, 6, 10},
		{7, 24, 25}, {24, 7, 25},
		{8, 15, 17}, {15, 8, 17},
		{9, 12, 15}, {12, 9, 15},
	}
	t := triples[rand.Intn(len(triples))]
	angle := rand.Intn(2) // 0 = angle opposite side opp, 1 = angle adjacent to side opp
	opp := t.opp
	adj := t.adj
	if angle == 0 {
		opp = t.adj
	}
	funcs := []string{"sin", "cos", "tan"}
	f := funcs[rand.Intn(len(funcs))]
	var ans string
	switch f {
	case "sin":
		ans = fmt.Sprintf("%d/%d", opp, t.hyp)
	case "cos":
		ans = fmt.Sprintf("%d/%d", adj, t.hyp)
	case "tan":
		ans = fmt.Sprintf("%d/%d", opp, adj)
	}
	return generator.Problem{
		Question:    fmt.Sprintf("In a right triangle with sides %d, %d, %d, what is \\(\\%s\\) of the angle opposite the side of length %d? (as a fraction)", t.opp, t.adj, t.hyp, f, t.opp),
		Answer:      ans,
		Explanation: fmt.Sprintf("SOH CAH TOA: %s(θ) = %s", f, ans),
	}
}

// basic trig equations
type trigEqBasicGen struct{}

func (g *trigEqBasicGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"Solve \\(\\sin(x) = 0\\) for \\(0^{\\circ} \\leq x < 360^{\\circ}\\).", "\\(0^{\\circ},180^{\\circ}\\)", "\\(\\sin(x) = 0\\) when \\(x = 0^{\\circ}\\) or \\(x = 180^{\\circ}\\) in \\([0^{\\circ},360^{\\circ})\\)."},
		{"Solve \\(\\cos(x) = 0\\) for \\(0^{\\circ} \\leq x < 360^{\\circ}\\).", "\\(90^{\\circ},270^{\\circ}\\)", "\\(\\cos(x) = 0\\) when \\(x = 90^{\\circ}\\) or \\(x = 270^{\\circ}\\) in \\([0^{\\circ},360^{\\circ})\\)."},
		{"Solve \\(\\sin(x) = 1\\) for \\(0^{\\circ} \\leq x < 360^{\\circ}\\).", "\\(90^{\\circ}\\)", "\\(\\sin(x) = 1\\) only at \\(x = 90^{\\circ}\\) in \\([0^{\\circ},360^{\\circ})\\)."},
		{"Solve \\(\\cos(x) = 1\\) for \\(0^{\\circ} \\leq x < 360^{\\circ}\\).", "\\(0^{\\circ}\\)", "\\(\\cos(x) = 1\\) only at \\(x = 0^{\\circ}\\) in \\([0^{\\circ},360^{\\circ})\\)."},
		{"Solve \\(\\sin(x) = -1\\) for \\(0^{\\circ} \\leq x < 360^{\\circ}\\).", "\\(270^{\\circ}\\)", "\\(\\sin(x) = -1\\) only at \\(x = 270^{\\circ}\\) in \\([0^{\\circ},360^{\\circ})\\)."},
		{"Solve \\(\\tan(x) = 0\\) for \\(0^{\\circ} \\leq x < 360^{\\circ}\\).", "\\(0^{\\circ},180^{\\circ}\\)", "\\(\\tan(x) = \\sin(x)/\\cos(x)\\), so \\(\\tan(x) = 0\\) when \\(\\sin(x) = 0\\) at \\(x = 0^{\\circ},180^{\\circ}\\)."},
		{"Solve \\(\\sin(x) = 1/2\\) for \\(0^{\\circ} \\leq x < 360^{\\circ}\\).", "\\(30^{\\circ},150^{\\circ}\\)", "\\(\\sin(30^{\\circ}) = 1/2\\) and \\(\\sin(150^{\\circ}) = 1/2\\) in \\([0^{\\circ},360^{\\circ})\\)."},
		{"Solve \\(\\cos(x) = 1/2\\) for \\(0^{\\circ} \\leq x < 360^{\\circ}\\).", "\\(60^{\\circ},300^{\\circ}\\)", "\\(\\cos(60^{\\circ}) = 1/2\\) and \\(\\cos(300^{\\circ}) = 1/2\\) in \\([0^{\\circ},360^{\\circ})\\)."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

// homogeneous trig equations
type trigEqHomogeneousGen struct{}

func (g *trigEqHomogeneousGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"What substitution is used to solve \\(a\\sin(x) + b\\cos(x) = 0\\)?", "\\(\\tan(x) = -b/a\\)", "Divide both sides by \\(\\cos(x)\\): \\(a\\tan(x) + b = 0 \\to \\tan(x) = -b/a\\)."},
		{"Solve \\(\\sin(x) - \\cos(x) = 0\\) for \\(0^{\\circ} \\leq x < 360^{\\circ}\\).", "\\(45^{\\circ},225^{\\circ}\\)", "\\(\\sin(x) = \\cos(x) \\to \\tan(x) = 1 \\to x = 45^{\\circ}, 225^{\\circ}\\)."},
		{"Solve \\(\\sin(x) + \\cos(x) = 0\\) for \\(0^{\\circ} \\leq x < 360^{\\circ}\\).", "\\(135^{\\circ},315^{\\circ}\\)", "\\(\\sin(x) = -\\cos(x) \\to \\tan(x) = -1 \\to x = 135^{\\circ}, 315^{\\circ}\\)."},
		{"Solve \\(\\sqrt{3}\\sin(x) - \\cos(x) = 0\\) for \\(0^{\\circ} \\leq x < 360^{\\circ}\\).", "\\(30^{\\circ},210^{\\circ}\\)", "\\(\\sqrt{3}\\sin(x) = \\cos(x) \\to \\tan(x) = 1/\\sqrt{3} \\to x = 30^{\\circ}, 210^{\\circ}\\)."},
		{"Solve \\(\\sin(x) - \\sqrt{3}\\cos(x) = 0\\) for \\(0^{\\circ} \\leq x < 360^{\\circ}\\).", "\\(60^{\\circ},240^{\\circ}\\)", "\\(\\sin(x) = \\sqrt{3}\\cos(x) \\to \\tan(x) = \\sqrt{3} \\to x = 60^{\\circ}, 240^{\\circ}\\)."},
		{"What is the general method to solve \\(a\\sin(x) + b\\cos(x) = 0\\)?", "divide by \\(\\cos(x)\\)", "Dividing by \\(\\cos(x)\\) gives \\(a\\tan(x) + b = 0\\), which can be solved for \\(x\\)."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

// hyperbolic: sinh and cosh
type sinhCoshGen struct{}

func (g *sinhCoshGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"What is the definition of \\(\\sinh(x)\\)?", "\\((e^{x} - e^{-x})/2\\)", "\\(\\sinh(x) = (e^{x} - e^{-x})/2\\), the odd part of the exponential function."},
		{"What is the definition of \\(\\cosh(x)\\)?", "\\((e^{x} + e^{-x})/2\\)", "\\(\\cosh(x) = (e^{x} + e^{-x})/2\\), the even part of the exponential function."},
		{"What is \\(\\cosh^{2}(x) - \\sinh^{2}(x)\\)?", "1", "\\(\\cosh^{2}(x) - \\sinh^{2}(x) = 1\\) (the hyperbolic analogue of \\(\\cos^{2}+\\sin^{2}=1\\))."},
		{"Is \\(\\sinh(x)\\) an even or odd function?", "odd", "\\(\\sinh(-x) = -\\sinh(x)\\), so \\(\\sinh\\) is odd."},
		{"Is \\(\\cosh(x)\\) an even or odd function?", "even", "\\(\\cosh(-x) = \\cosh(x)\\), so \\(\\cosh\\) is even."},
		{"What is the derivative of \\(\\sinh(x)\\)?", "\\(\\cosh(x)\\)", "\\(d/dx \\sinh(x) = \\cosh(x)\\)."},
		{"What is the derivative of \\(\\cosh(x)\\)?", "\\(\\sinh(x)\\)", "\\(d/dx \\cosh(x) = \\sinh(x)\\)."},
		{"What is the identity relating \\(\\cosh^{2}(x)\\) and \\(\\sinh^{2}(x)\\)?", "\\(\\cosh^{2}(x) - \\sinh^{2}(x) = 1\\)", "\\(\\cosh^{2}(x) - \\sinh^{2}(x) = 1\\) is the fundamental hyperbolic identity."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

// hyperbolic: tanh and coth
type tanhCothGen struct{}

func (g *tanhCothGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"What is the definition of \\(\\tanh(x)\\)?", "\\(\\sinh(x)/\\cosh(x)\\)", "\\(\\tanh(x) = \\sinh(x)/\\cosh(x) = (e^{x} - e^{-x})/(e^{x} + e^{-x})\\)."},
		{"What is the definition of \\(\\coth(x)\\)?", "\\(\\cosh(x)/\\sinh(x)\\)", "\\(\\coth(x) = \\cosh(x)/\\sinh(x) = 1/\\tanh(x)\\)."},
		{"What is \\(\\tanh(0)\\)?", "0", "\\(\\tanh(0) = \\sinh(0)/\\cosh(0) = 0/1 = 0\\)."},
		{"As \\(x \\to \\infty\\), \\(\\tanh(x)\\) approaches what value?", "1", "\\(\\lim_{x\\to\\infty} \\tanh(x) = 1\\) because \\(e^{x}\\) dominates \\(e^{-x}\\)."},
		{"As \\(x \\to -\\infty\\), \\(\\tanh(x)\\) approaches what value?", "-1", "\\(\\lim_{x\\to-\\infty} \\tanh(x) = -1\\) because \\(e^{-x}\\) dominates \\(e^{x}\\)."},
		{"What is the range of \\(\\tanh(x)\\)?", "\\((-1, 1)\\)", "\\(\\tanh(x)\\) maps real numbers to the open interval \\((-1, 1)\\)."},
		{"What is \\(1 - \\tanh^{2}(x)\\)?", "\\(\\sech^{2}(x)\\)", "\\(1 - \\tanh^{2}(x) = \\sech^{2}(x) = 1/\\cosh^{2}(x)\\)."},
		{"Is \\(\\tanh(x)\\) an even or odd function?", "odd", "\\(\\tanh(-x) = -\\tanh(x)\\), so \\(\\tanh\\) is odd."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

// trig identities
type trigIdentGen struct{}

func (g *trigIdentGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"Simplify \\(\\sin^{2}(x) + \\cos^{2}(x)\\).", "1", "\\(\\sin^{2}(x) + \\cos^{2}(x) = 1\\) (the Pythagorean identity)."},
		{"Simplify \\(1 + \\tan^{2}(x)\\).", "\\(\\sec^{2}(x)\\)", "\\(1 + \\tan^{2}(x) = \\sec^{2}(x)\\), derived from \\(\\sin^{2}+\\cos^{2}=1\\) divided by \\(\\cos^{2}\\)."},
		{"Simplify \\(1 + \\cot^{2}(x)\\).", "\\(\\csc^{2}(x)\\)", "\\(1 + \\cot^{2}(x) = \\csc^{2}(x)\\), derived from \\(\\sin^{2}+\\cos^{2}=1\\) divided by \\(\\sin^{2}\\)."},
		{"What is \\(\\sin(2x)\\) in terms of \\(\\sin(x)\\) and \\(\\cos(x)\\)?", "\\(2\\sin(x)\\cos(x)\\)", "\\(\\sin(2x) = 2\\sin(x)\\cos(x)\\) (double angle formula)."},
		{"What is \\(\\cos(2x)\\) in terms of \\(\\cos(x)\\)?", "\\(2\\cos^{2}(x) - 1\\)", "\\(\\cos(2x) = 2\\cos^{2}(x) - 1 = \\cos^{2}(x) - \\sin^{2}(x) = 1 - 2\\sin^{2}(x)\\)."},
		{"What is \\(\\sin(-x)\\) in terms of \\(\\sin(x)\\)?", "\\(-\\sin(x)\\)", "\\(\\sin(-x) = -\\sin(x)\\) (sine is odd)."},
		{"What is \\(\\cos(-x)\\) in terms of \\(\\cos(x)\\)?", "\\(\\cos(x)\\)", "\\(\\cos(-x) = \\cos(x)\\) (cosine is even)."},
		{"Simplify \\(\\sin(x)\\cos(y) + \\cos(x)\\sin(y)\\).", "\\(\\sin(x+y)\\)", "\\(\\sin(x+y) = \\sin(x)\\cos(y) + \\cos(x)\\sin(y)\\) (addition formula)."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

// basic trig inequalities
type trigIneqGen struct{}

func (g *trigIneqGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"For \\(0^{\\circ} < x < 90^{\\circ}\\), is \\(\\sin(x)\\) increasing or decreasing?", "increasing", "\\(\\sin(x)\\) increases from \\(0\\) to \\(1\\) on \\([0^{\\circ}, 90^{\\circ}]\\)."},
		{"For \\(0^{\\circ} < x < 90^{\\circ}\\), is \\(\\cos(x)\\) increasing or decreasing?", "decreasing", "\\(\\cos(x)\\) decreases from \\(1\\) to \\(0\\) on \\([0^{\\circ}, 90^{\\circ}]\\)."},
		{"For \\(0^{\\circ} < x < 90^{\\circ}\\), is \\(\\tan(x)\\) increasing or decreasing?", "increasing", "\\(\\tan(x)\\) increases from \\(0\\) to \\(\\infty\\) on \\([0^{\\circ}, 90^{\\circ})\\)."},
		{"What is the maximum value of \\(\\sin(x)\\)?", "1", "The maximum of \\(\\sin(x)\\) is \\(1\\), achieved at \\(x = 90^{\\circ}\\)."},
		{"What is the minimum value of \\(\\sin(x)\\)?", "-1", "The minimum of \\(\\sin(x)\\) is \\(-1\\), achieved at \\(x = 270^{\\circ}\\)."},
		{"What is the maximum value of \\(\\cos(x)\\)?", "1", "The maximum of \\(\\cos(x)\\) is \\(1\\), achieved at \\(x = 0^{\\circ}\\)."},
		{"What is the minimum value of \\(\\cos(x)\\)?", "-1", "The minimum of \\(\\cos(x)\\) is \\(-1\\), achieved at \\(x = 180^{\\circ}\\)."},
		{"For \\(90^{\\circ} < x < 180^{\\circ}\\), is \\(\\sin(x)\\) positive or negative?", "positive", "\\(\\sin(x)\\) is positive in Quadrant II (\\(90^{\\circ}\\) to \\(180^{\\circ}\\))."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}
