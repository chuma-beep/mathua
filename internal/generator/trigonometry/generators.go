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

	// stub generators for newly-added concepts
	reg.Register("trig.adv.arctan", &generator.Stub{ConceptID: "trig.adv.arctan"})
	reg.Register("trig.basics.right_triangle", &generator.Stub{ConceptID: "trig.basics.right_triangle"})
	reg.Register("trig.eq.basic", &generator.Stub{ConceptID: "trig.eq.basic"})
	reg.Register("trig.eq.homogeneous", &generator.Stub{ConceptID: "trig.eq.homogeneous"})
	reg.Register("trig.hyperbolic.sinh_cosh", &generator.Stub{ConceptID: "trig.hyperbolic.sinh_cosh"})
	reg.Register("trig.hyperbolic.tanh_coth", &generator.Stub{ConceptID: "trig.hyperbolic.tanh_coth"})
	reg.Register("trig.ident.identities", &generator.Stub{ConceptID: "trig.ident.identities"})
	reg.Register("trig.ineq.basic", &generator.Stub{ConceptID: "trig.ineq.basic"})
}

type radiansGen struct{}

func (g *radiansGen) Generate(difficulty float64) generator.Problem {
	type entry struct {
		deg      int
		radHTML  string
		radLatex string
	}
	table := []entry{
		{30, "π/6", "\\pi/6"},
		{45, "π/4", "\\pi/4"},
		{60, "π/3", "\\pi/3"},
		{90, "π/2", "\\pi/2"},
		{180, "π", "\\pi"},
		{270, "3π/2", "3\\pi/2"},
		{360, "2π", "2\\pi"},
	}
	e := table[rand.Intn(len(table))]
	if rand.Intn(2) == 0 {
		return generator.Problem{
			Question:    fmt.Sprintf("Convert %d° to radians.", e.deg),
			Answer:      e.radHTML,
			Explanation: fmt.Sprintf("%d° × π/180 = %s.", e.deg, e.radHTML),
		}
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Convert %s to degrees.", e.radHTML),
		Answer:      fmt.Sprintf("%d", e.deg),
		Explanation: fmt.Sprintf("%s × 180/π = %d°.", e.radHTML, e.deg),
	}
}

type unitCircleGen struct{}

func (g *unitCircleGen) Generate(difficulty float64) generator.Problem {
	type angle struct {
		label string
		sin   string
		cos   string
	}
	angles := []angle{
		{"0°", "0", "1"},
		{"30°", "1/2", "√3/2"},
		{"45°", "√2/2", "√2/2"},
		{"60°", "√3/2", "1/2"},
		{"90°", "1", "0"},
		{"120°", "√3/2", "-1/2"},
		{"135°", "√2/2", "-√2/2"},
		{"150°", "1/2", "-√3/2"},
		{"180°", "0", "-1"},
		{"210°", "-1/2", "-√3/2"},
		{"225°", "-√2/2", "-√2/2"},
		{"240°", "-√3/2", "-1/2"},
		{"270°", "-1", "0"},
		{"300°", "-√3/2", "1/2"},
		{"315°", "-√2/2", "√2/2"},
		{"330°", "-1/2", "√3/2"},
		{"360°", "0", "1"},
	}
	a := angles[rand.Intn(len(angles))]
	if rand.Intn(2) == 0 {
		return generator.Problem{
			Question:    fmt.Sprintf("What is sin(%s)?", a.label),
			Answer:      a.sin,
			Explanation: fmt.Sprintf("On the unit circle, sin(%s) = %s.", a.label, a.sin),
		}
	}
	return generator.Problem{
		Question:    fmt.Sprintf("What is cos(%s)?", a.label),
		Answer:      a.cos,
		Explanation: fmt.Sprintf("On the unit circle, cos(%s) = %s.", a.label, a.cos),
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

func (g *sinCosDefGen) Generate(difficulty float64) generator.Problem {
	t := pythagoreanTriples[rand.Intn(len(pythagoreanTriples))]
	if rand.Intn(2) == 0 {
		gcd := mathutil.GCD(t.opp, t.hyp)
		num := t.opp / gcd
		den := t.hyp / gcd
		return generator.Problem{
			Question:    fmt.Sprintf("In a right triangle with opposite = %d and hypotenuse = %d, what is sin(θ)?", t.opp, t.hyp),
			Answer:      fmt.Sprintf("%d/%d", num, den),
			Explanation: fmt.Sprintf("sin(θ) = opp/hyp = %d/%d = %d/%d.", t.opp, t.hyp, num, den),
		}
	}
	gcd := mathutil.GCD(t.adj, t.hyp)
	num := t.adj / gcd
	den := t.hyp / gcd
	return generator.Problem{
		Question:    fmt.Sprintf("In a right triangle with adjacent = %d and hypotenuse = %d, what is cos(θ)?", t.adj, t.hyp),
		Answer:      fmt.Sprintf("%d/%d", num, den),
		Explanation: fmt.Sprintf("cos(θ) = adj/hyp = %d/%d = %d/%d.", t.adj, t.hyp, num, den),
	}
}

type tanDefGen struct{}

func (g *tanDefGen) Generate(difficulty float64) generator.Problem {
	// Half the time use right triangle sides, half use sin/cos ratio
	if rand.Intn(2) == 0 {
		t := pythagoreanTriples[rand.Intn(len(pythagoreanTriples))]
		gcd := mathutil.GCD(t.opp, t.adj)
		num := t.opp / gcd
		den := t.adj / gcd
		return generator.Problem{
			Question:    fmt.Sprintf("In a right triangle with opposite = %d and adjacent = %d, what is tan(θ)?", t.opp, t.adj),
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
		Question:    fmt.Sprintf("If sin(θ) = %d/%d and cos(θ) = %d/%d, what is tan(θ)?", sinNum, sinDen, cosNum, cosDen),
		Answer:      fmt.Sprintf("%d/%d", tanNum, tanDen),
		Explanation: fmt.Sprintf("tan(θ) = sin(θ)/cos(θ) = (%d/%d)/(%d/%d) = %d/%d.", sinNum, sinDen, cosNum, cosDen, tanNum, tanDen),
	}
}

type reciprocalGen struct{}

func (g *reciprocalGen) Generate(difficulty float64) generator.Problem {
	t := pythagoreanTriples[rand.Intn(len(pythagoreanTriples))]
	funcs := []struct {
		name  string
		num   int
		den   int
		recip string
	}{
		{"sin(θ)", t.opp, t.hyp, "csc(θ)"},
		{"cos(θ)", t.adj, t.hyp, "sec(θ)"},
		{"tan(θ)", t.opp, t.adj, "cot(θ)"},
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

func (g *pythagoreanIDGen) Generate(difficulty float64) generator.Problem {
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
			Question:    fmt.Sprintf("If sin(θ) = %d/%d, find cos(θ) using sin²θ + cos²θ = 1 (assume θ is acute).", sinNum, sinDen),
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
		Question:    fmt.Sprintf("If cos(θ) = %d/%d, find sin(θ) using sin²θ + cos²θ = 1 (assume θ is acute).", cosNum, cosDen),
		Answer:      fmt.Sprintf("%d/%d", sinNum, sinDen),
		Explanation: fmt.Sprintf("sin²θ = 1 - cos²θ = 1 - (%d/%d)² = 1 - %d/%d = %d/%d, so sin(θ) = √(%d/%d) = %d/%d.", cosNum, cosDen, cosNum*cosNum, cosDen*cosDen, cosDen*cosDen-cosNum*cosNum, cosDen*cosDen, cosDen*cosDen-cosNum*cosNum, cosDen*cosDen, sinNum, sinDen),
	}
}

type specialAnglesGen struct{}

func (g *specialAnglesGen) Generate(difficulty float64) generator.Problem {
	type entry struct {
		label string
		sin   string
		cos   string
		tan   string
	}
	table := []entry{
		{"30°", "1/2", "√3/2", "√3/3"},
		{"45°", "√2/2", "√2/2", "1"},
		{"60°", "√3/2", "1/2", "√3"},
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

func (g *referenceAngleGen) Generate(difficulty float64) generator.Problem {
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

func (g *graphSinGen) Generate(difficulty float64) generator.Problem {
	questions := []struct {
		q, a, e string
	}{
		{"What is sin(0)?", "0", "sin(0) = 0"},
		{"What is the maximum value of y = sin(x)?", "1", "The sine function has a maximum value of 1."},
		{"What is the minimum value of y = sin(x)?", "-1", "The sine function has a minimum value of -1."},
		{"What is the period of y = sin(x)?", "360°", "The period of sin(x) is 360° (2π radians)."},
		{"What is the amplitude of y = sin(x)?", "1", "The amplitude of sin(x) is 1."},
		{"What is sin(90°)?", "1", "sin(90°) = 1"},
		{"What is sin(180°)?", "0", "sin(180°) = 0"},
		{"What is sin(270°)?", "-1", "sin(270°) = -1"},
		{"What is sin(360°)?", "0", "sin(360°) = 0"},
		{"Where does sin(x) start on the y-axis?", "0", "sin(0) = 0"},
	}
	q := questions[rand.Intn(len(questions))]
	return generator.Problem{
		Question:    q.q,
		Answer:      q.a,
		Explanation: q.e,
	}
}

type graphCosGen struct{}

func (g *graphCosGen) Generate(difficulty float64) generator.Problem {
	questions := []struct {
		q, a, e string
	}{
		{"What is cos(0)?", "1", "cos(0) = 1"},
		{"What is the maximum value of y = cos(x)?", "1", "The cosine function has a maximum value of 1."},
		{"What is the minimum value of y = cos(x)?", "-1", "The cosine function has a minimum value of -1."},
		{"What is the period of y = cos(x)?", "360°", "The period of cos(x) is 360° (2π radians)."},
		{"What is the amplitude of y = cos(x)?", "1", "The amplitude of cos(x) is 1."},
		{"What is cos(90°)?", "0", "cos(90°) = 0"},
		{"What is cos(180°)?", "-1", "cos(180°) = -1"},
		{"What is cos(270°)?", "0", "cos(270°) = 0"},
		{"What is cos(360°)?", "1", "cos(360°) = 1"},
		{"Where does cos(x) start on the y-axis?", "1", "cos(0) = 1"},
	}
	q := questions[rand.Intn(len(questions))]
	return generator.Problem{
		Question:    q.q,
		Answer:      q.a,
		Explanation: q.e,
	}
}

type periodGen struct{}

func (g *periodGen) Generate(difficulty float64) generator.Problem {
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
			fmt.Sprintf("What is the amplitude of y = %d sin(%dx)?", c.A, c.B),
			fmt.Sprintf("%d", c.amp),
			fmt.Sprintf("Amplitude = |A| = |%d| = %d.", c.A, c.amp),
		},
		{
			fmt.Sprintf("What is the period of y = sin(%dx)?", c.B),
			c.period,
			fmt.Sprintf("Period = 360°/|B| = 360°/%d = %s.", c.B, c.period),
		},
		{
			fmt.Sprintf("What is the phase shift of y = sin(%dx + %d)?", c.B, c.C),
			c.shift,
			fmt.Sprintf("Phase shift = -C/B = -%d/%d = -%s (shifted right).", c.C, c.B, c.shift),
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

func (g *inverseGen) Generate(difficulty float64) generator.Problem {
	type entry struct {
		f      string
		val    string
		result string
		desc   string
	}
	table := []entry{
		{"arcsin(0)", "0", "0°", "arcsin(0) = 0°"},
		{"arcsin(1/2)", "1/2", "30°", "arcsin(1/2) = 30°"},
		{"arcsin(√2/2)", "√2/2", "45°", "arcsin(√2/2) = 45°"},
		{"arcsin(√3/2)", "√3/2", "60°", "arcsin(√3/2) = 60°"},
		{"arcsin(1)", "1", "90°", "arcsin(1) = 90°"},
		{"arccos(1)", "1", "0°", "arccos(1) = 0°"},
		{"arccos(√3/2)", "√3/2", "30°", "arccos(√3/2) = 30°"},
		{"arccos(√2/2)", "√2/2", "45°", "arccos(√2/2) = 45°"},
		{"arccos(1/2)", "1/2", "60°", "arccos(1/2) = 60°"},
		{"arccos(0)", "0", "90°", "arccos(0) = 90°"},
		{"arctan(0)", "0", "0°", "arctan(0) = 0°"},
		{"arctan(√3/3)", "√3/3", "30°", "arctan(√3/3) = 30°"},
		{"arctan(1)", "1", "45°", "arctan(1) = 45°"},
		{"arctan(√3)", "√3", "60°", "arctan(√3) = 60°"},
		{"arcsin(-1/2)", "-1/2", "-30°", "arcsin(-1/2) = -30°"},
		{"arccos(-1/2)", "-1/2", "120°", "arccos(-1/2) = 120°"},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{
		Question:    fmt.Sprintf("What is %s?", e.f),
		Answer:      e.result,
		Explanation: e.desc,
	}
}

type lawSinesGen struct{}

func (g *lawSinesGen) Generate(difficulty float64) generator.Problem {
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
		Question:    fmt.Sprintf("In triangle ABC, A = %d°, B = %d°, and side a = %d. Find side b using the law of sines.", aDeg, bDeg, c.a),
		Answer:      fmt.Sprintf("%.1f", bRounded),
		Explanation: fmt.Sprintf("a/sin(A) = b/sin(B), so b = a·sin(B)/sin(A) = %d·sin(%d°)/sin(%d°) = %.1f.", c.a, bDeg, aDeg, bRounded),
	}
}

type lawCosinesGen struct{}

func (g *lawCosinesGen) Generate(difficulty float64) generator.Problem {
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
		Question:    fmt.Sprintf("In triangle ABC, a = %d, b = %d, and C = %d°. Find side c using the law of cosines.", c.a, c.b, c.C),
		Answer:      fmt.Sprintf("%.1f", cRounded),
		Explanation: fmt.Sprintf("c² = a² + b² - 2ab·cos(C) = %d² + %d² - 2(%d)(%d)·cos(%d°) = %d + %d - %d = %d, so c = √%d ≈ %.1f.", c.a, c.b, c.a, c.b, c.C, c.a*c.a, c.b*c.b, c.a*c.b, c.c2, c.c2, cRounded),
	}
}
