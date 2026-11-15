package trigonometry

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/grader"
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
	reg.Register("trig.func.arcsine", &arcsineFuncGen{})
	reg.Register("trig.func.arccosine", &arccosineFuncGen{})
	reg.Register("trig.func.arctangent", &arctangentFuncGen{})
	reg.Register("trig.func.arccotangent", &arccotangentFuncGen{})

	// Concepts added by the 2026 content expansion.
	reg.Register("trig.eq.multiple_angle", &trigBridgeGen{kind: "multiple_angle"})
	reg.Register("trig.basics.elevation_depression", &trigBridgeGen{kind: "elevation"})
	reg.Register("trig.graph.sinusoidal_model", &trigBridgeGen{kind: "sinusoidal_model"})
	reg.Register("trig.adv.inverse_composition", &trigBridgeGen{kind: "inverse_composition"})
	reg.Register("trig.adv.vector_components", &trigBridgeGen{kind: "vector_components"})
	reg.Register("trig.ident.even_odd", &trigBridgeGen{kind: "even_odd"})
	reg.Register("trig.graph.frequency", &trigBridgeGen{kind: "frequency"})
	reg.Register("trig.eq.using_identities", &trigBridgeGen{kind: "using_identities"})
	reg.Register("trig.graph.harmonic_motion", &trigBridgeGen{kind: "harmonic_motion"})
	reg.Register("trig.basics.segment_area", &trigBridgeGen{kind: "segment_area"})
	registerComputed(reg)
	registerRaise(reg)
	registerFuncRatios(reg)
	registerWave8(reg)
	registerWave8Remaining(reg)
	registerWave9(reg)
}

type radiansGen struct{}

// trig.basics.radians is a multiple-choice concept, but this generator served a
// bare question whose answer was "\\(\\pi/2\\)", delimiters included, and no
// options to tap. A learner typed pi/2, which is right, and the choice grader
// compared it with the delimited string and marked it wrong (report 119). The
// options are plain text so they render in a button, and the answer is the
// option's own text so a tap resolves to it.
func (g *radiansGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return angleConversionProblem()
}

// others returns every value except one, for use as distractors.
func others(values []string, except string) []string {
	out := make([]string, 0, len(values))
	for _, v := range values {
		if v != except {
			out = append(out, v)
		}
	}
	return out
}

type unitCircleGen struct{}

// The angle is stored as its degree value rather than as a ready-made label. The label carries
// its own math delimiters because the solution schema quotes it in prose ("the angle [[label]] is
// marked on the unit circle"), but the question and explanation put the angle *inside* a math
// region, where a delimited label would nest one region in another -- `\(\sin(\(180^{\circ}\))\)`
// -- and the canonicaliser rendered that as `$\sin((180^{\circ}$))`, a broken question. Building
// the bare angle for the math and the delimited label for the schema keeps both correct.
var unitCircleAngles = []struct {
	deg  int
	sin  string
	cos  string
	easy bool
}{
	{0, "0", "1", true},
	{30, "\\(1/2\\)", "\\(\\sqrt{3}/2\\)", true},
	{45, "\\(\\sqrt{2}/2\\)", "\\(\\sqrt{2}/2\\)", true},
	{60, "\\(\\sqrt{3}/2\\)", "\\(1/2\\)", true},
	{90, "1", "0", true},
	{180, "0", "-1", true},
	{270, "-1", "0", true},
	{360, "0", "1", true},
	{120, "\\(\\sqrt{3}/2\\)", "\\(-1/2\\)", false},
	{135, "\\(\\sqrt{2}/2\\)", "\\(-\\sqrt{2}/2\\)", false},
	{150, "\\(1/2\\)", "\\(-\\sqrt{3}/2\\)", false},
	{210, "\\(-1/2\\)", "\\(-\\sqrt{3}/2\\)", false},
	{225, "\\(-\\sqrt{2}/2\\)", "\\(-\\sqrt{2}/2\\)", false},
	{240, "\\(-\\sqrt{3}/2\\)", "\\(-1/2\\)", false},
	{300, "\\(-\\sqrt{3}/2\\)", "\\(1/2\\)", false},
	{315, "\\(-\\sqrt{2}/2\\)", "\\(\\sqrt{2}/2\\)", false},
	{330, "\\(-1/2\\)", "\\(\\sqrt{3}/2\\)", false},
}

// bare strips a value's own math delimiters so it can be embedded inside another math region.
// Some unit-circle values are written delimited ("\(\sqrt{2}/2\)") and some bare ("-1"); the
// question and explanation put the value inside `\(...\)`, where a delimited one would nest.
func bare(v string) string {
	return strings.TrimSuffix(strings.TrimPrefix(v, "\\("), "\\)")
}

func (g *unitCircleGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return unitCircleFactsProblem()
}

// unitCirclePlainValues are the answers a unit-circle question can have, written
// the way a learner reads them. They are option labels, so they carry no math
// delimiters: the option row renders plain text.
var unitCirclePlainValues = []string{
	"-1", "-√3/2", "-√2/2", "-1/2", "0", "1/2", "√2/2", "√3/2", "1",
}

// plainUnitValue turns a table value into its plain label. The table writes
// radicals as LaTeX because the question embeds them inside a math region.
func plainUnitValue(v string) string {
	v = bare(v)
	v = strings.ReplaceAll(v, "\\sqrt{2}", "√2")
	v = strings.ReplaceAll(v, "\\sqrt{3}", "√3")
	return v
}

// unitCircleDistractors returns three other values from the same family, so the
// choice tests which coordinate the angle gives rather than which symbol looks
// familiar.
func unitCircleDistractors(want string) []string {
	pool := make([]string, 0, len(unitCirclePlainValues))
	for _, v := range unitCirclePlainValues {
		if v != want {
			pool = append(pool, v)
		}
	}
	rand.Shuffle(len(pool), func(i, j int) { pool[i], pool[j] = pool[j], pool[i] })
	return pool[:3]
}

// unitCircleProblem serves the value as options because a radical and a fraction
// slash have no key on a phone keypad: report 28 was a learner who could not
// enter the answer at all. The answer text is the plain value the learner reads.
// Facts["val"] keeps its math delimiters because the solution schema quotes it
// inside prose ("its vertical coordinate is [[val]]").
func unitCircleProblem(angle, question, value, fn, label string) generator.Problem {
	plain := plainUnitValue(value)
	p := generator.MultipleChoice(question, plain, unitCircleDistractors(plain)...)
	p.Explanation = fmt.Sprintf("On the unit circle, \\(%s(%s) = %s\\).", fn, angle, bare(value))
	p.Facts = map[string]string{"label": label, "val": value, "fn": fn}
	return p
}

type sinCosDefGen struct{}

type triple struct{ opp, adj, hyp int }

var pythagoreanTriples = []triple{
	{3, 4, 5},
	{5, 12, 13},
	{8, 15, 17},
	{7, 24, 25},
	{9, 40, 41},
	{20, 21, 29},
	{12, 35, 37},
	{28, 45, 53},
	{11, 60, 61},
	{13, 84, 85},
	{16, 63, 65},
	{33, 56, 65},
	{36, 77, 85},
	{39, 80, 89},
}

func (g *sinCosDefGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	maxTriples := 3 + int(ctx.Difficulty*float64(len(pythagoreanTriples)-3))
	if maxTriples > len(pythagoreanTriples) {
		maxTriples = len(pythagoreanTriples)
	}
	t := pythagoreanTriples[rand.Intn(maxTriples)]
	k := 1 + rand.Intn(1+int(ctx.Difficulty*4))
	if rand.Intn(2) == 0 {
		opp := t.opp * k
		hyp := t.hyp * k
		gcdVal := mathutil.GCD(opp, hyp)
		num := opp / gcdVal
		den := hyp / gcdVal
		return generator.Problem{
			Question:    fmt.Sprintf("In a right triangle, the leg opposite the acute angle \\(\\theta\\) is %d units long and the hypotenuse is %d units long. What is \\(\\sin(\\theta)\\)?", opp, hyp),
			Answer:      fmt.Sprintf("%d/%d", num, den),
			Explanation: fmt.Sprintf("The sine of an angle is defined as the ratio of opposite leg to hypotenuse: \\(\\sin(\\theta) = \\frac{\\text{opposite}}{\\text{hypotenuse}}\\). Substitute the given side lengths %d and %d: \\(\\sin(\\theta) = \\frac{%d}{%d}\\). Simplifying the fraction by dividing numerator and denominator by their greatest common divisor yields %d/%d.", opp, hyp, opp, hyp, num, den),
			Facts:       map[string]string{"opp": strconv.Itoa(opp), "hyp": strconv.Itoa(hyp), "num": strconv.Itoa(num), "den": strconv.Itoa(den)},
		}
	}
	adj := t.adj * k
	hyp := t.hyp * k
	gcdVal := mathutil.GCD(adj, hyp)
	num := adj / gcdVal
	den := hyp / gcdVal
	return generator.Problem{
		Question:    fmt.Sprintf("In a right triangle, the leg adjacent to the acute angle \\(\\theta\\) is %d units long and the hypotenuse is %d units long. What is \\(\\cos(\\theta)\\)?", adj, hyp),
		Answer:      fmt.Sprintf("%d/%d", num, den),
		Explanation: fmt.Sprintf("The cosine of an angle is defined as the ratio of adjacent leg to hypotenuse: \\(\\cos(\\theta) = \\frac{\\text{adjacent}}{\\text{hypotenuse}}\\). Substitute the given side lengths %d and %d: \\(\\cos(\\theta) = \\frac{%d}{%d}\\). Simplifying the fraction by dividing numerator and denominator by their greatest common divisor yields %d/%d.", adj, hyp, adj, hyp, num, den),
		Facts:       map[string]string{"adj": strconv.Itoa(adj), "hyp": strconv.Itoa(hyp), "num": strconv.Itoa(num), "den": strconv.Itoa(den)},
	}
}

type tanDefGen struct{}

func (g *tanDefGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	maxTriples := 3 + int(ctx.Difficulty*float64(len(pythagoreanTriples)-3))
	if maxTriples > len(pythagoreanTriples) {
		maxTriples = len(pythagoreanTriples)
	}
	t := pythagoreanTriples[rand.Intn(maxTriples)]
	// Half the time use right triangle sides, half use sin/cos ratio
	if rand.Intn(2) == 0 {
		k := 1 + rand.Intn(1+int(ctx.Difficulty*4))
		opp := t.opp * k
		adj := t.adj * k
		gcdVal := mathutil.GCD(opp, adj)
		num := opp / gcdVal
		den := adj / gcdVal
		return generator.Problem{
			Question:    fmt.Sprintf("In a right triangle, the leg opposite the acute angle \\(\\theta\\) is %d units long and the leg adjacent to \\(\\theta\\) is %d units long. What is \\(\\tan(\\theta)\\)?", opp, adj),
			Answer:      fmt.Sprintf("%d/%d", num, den),
			Explanation: fmt.Sprintf("The tangent of an angle is defined as the ratio of opposite leg to adjacent leg: \\(\\tan(\\theta) = \\frac{\\text{opposite}}{\\text{adjacent}}\\). Substitute the given side lengths %d and %d: \\(\\tan(\\theta) = \\frac{%d}{%d}\\). Simplifying the fraction by dividing numerator and denominator by their greatest common divisor yields %d/%d.", opp, adj, opp, adj, num, den),
			Facts:       map[string]string{"opp": strconv.Itoa(opp), "adj": strconv.Itoa(adj), "num": strconv.Itoa(num), "den": strconv.Itoa(den)},
		}
	}
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
		Question:    fmt.Sprintf("For an acute angle \\(\\theta\\), if \\(\\sin(\\theta) = %d/%d\\) and \\(\\cos(\\theta) = %d/%d\\), what is \\(\\tan(\\theta)\\)?", sinNum, sinDen, cosNum, cosDen),
		Answer:      fmt.Sprintf("%d/%d", tanNum, tanDen),
		Explanation: fmt.Sprintf("The tangent of an angle is defined as the quotient of sine and cosine: \\(\\tan(\\theta) = \\frac{\\sin(\\theta)}{\\cos(\\theta)}\\). Substitute the given ratios: \\(\\tan(\\theta) = \\frac{%d/%d}{%d/%d} = \\frac{%d}{%d} \\times \\frac{%d}{%d}\\). Multiplying and reducing to lowest terms yields %d/%d.", sinNum, sinDen, cosNum, cosDen, sinNum, sinDen, cosDen, cosNum, tanNum, tanDen),
		Facts:       map[string]string{"sin_num": strconv.Itoa(sinNum), "sin_den": strconv.Itoa(sinDen), "cos_num": strconv.Itoa(cosNum), "cos_den": strconv.Itoa(cosDen), "num": strconv.Itoa(tanNum), "den": strconv.Itoa(tanDen)},
	}
}

type reciprocalGen struct{}

func (g *reciprocalGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	maxTriples := 3 + int(ctx.Difficulty*float64(len(pythagoreanTriples)-3))
	if maxTriples > len(pythagoreanTriples) {
		maxTriples = len(pythagoreanTriples)
	}
	t := pythagoreanTriples[rand.Intn(maxTriples)]
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
	gcdVal := mathutil.GCD(fn.den, fn.num)
	recipNum := fn.den / gcdVal
	recipDen := fn.num / gcdVal
	return generator.Problem{
		Question:    fmt.Sprintf("If %s = %d/%d, what is %s?", fn.name, fn.num, fn.den, fn.recip),
		Answer:      fmt.Sprintf("%d/%d", recipNum, recipDen),
		Explanation: fmt.Sprintf("The reciprocal of a trigonometric function is defined as one divided by that function: %s = 1/%s. Substitute the given fraction %d/%d: 1/(%d/%d) = %d/%d. Simplifying to lowest terms yields %d/%d.", fn.recip, fn.name, fn.num, fn.den, fn.num, fn.den, fn.den, fn.num, recipNum, recipDen),
		Facts:       map[string]string{"name": fn.name, "num": strconv.Itoa(fn.num), "den": strconv.Itoa(fn.den), "recip": fn.recip, "recip_num": strconv.Itoa(recipNum), "recip_den": strconv.Itoa(recipDen)},
	}
}

type pythagoreanIDGen struct{}

func (g *pythagoreanIDGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	maxTriples := 3 + int(ctx.Difficulty*float64(len(pythagoreanTriples)-3))
	if maxTriples > len(pythagoreanTriples) {
		maxTriples = len(pythagoreanTriples)
	}
	t := pythagoreanTriples[rand.Intn(maxTriples)]
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
			Explanation: fmt.Sprintf("By the Pythagorean trigonometric identity, \\(\\sin^{2}\\theta + \\cos^{2}\\theta = 1\\). Substitute \\(\\sin(\\theta) = %d/%d\\) and solve for cosine: \\(\\cos^{2}\\theta = 1 - (%d/%d)^{2} = 1 - %d/%d = %d/%d\\). Taking the principal square root for acute \\(\\theta\\) and simplifying gives \\(\\cos(\\theta) = %d/%d\\).", sinNum, sinDen, sinNum, sinDen, sinNum*sinNum, sinDen*sinDen, sinDen*sinDen-sinNum*sinNum, sinDen*sinDen, cosNum, cosDen),
			Facts:       map[string]string{"known": strconv.Itoa(sinNum) + "/" + strconv.Itoa(sinDen), "known_num": strconv.Itoa(sinNum), "known_den": strconv.Itoa(sinDen), "want": "$\\cos(\\theta)$", "sq_num": strconv.Itoa(sinNum * sinNum), "sq_den": strconv.Itoa(sinDen * sinDen), "rem_num": strconv.Itoa(sinDen*sinDen - sinNum*sinNum), "num": strconv.Itoa(cosNum), "den": strconv.Itoa(cosDen)},
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
		Explanation: fmt.Sprintf("By the Pythagorean trigonometric identity, \\(\\sin^{2}\\theta + \\cos^{2}\\theta = 1\\). Substitute \\(\\cos(\\theta) = %d/%d\\) and solve for sine: \\(\\sin^{2}\\theta = 1 - (%d/%d)^{2} = 1 - %d/%d = %d/%d\\). Taking the principal square root for acute \\(\\theta\\) and simplifying gives \\(\\sin(\\theta) = %d/%d\\).", cosNum, cosDen, cosNum, cosDen, cosNum*cosNum, cosDen*cosDen, cosDen*cosDen-cosNum*cosNum, cosDen*cosDen, sinNum, sinDen),
		Facts:       map[string]string{"known": strconv.Itoa(cosNum) + "/" + strconv.Itoa(cosDen), "known_num": strconv.Itoa(cosNum), "known_den": strconv.Itoa(cosDen), "want": "$\\sin(\\theta)$", "sq_num": strconv.Itoa(cosNum * cosNum), "sq_den": strconv.Itoa(cosDen * cosDen), "rem_num": strconv.Itoa(cosDen*cosDen - cosNum*cosNum), "num": strconv.Itoa(sinNum), "den": strconv.Itoa(sinDen)},
	}
}

// trig.basics.special_angles is a multiple-choice concept that served its value
// as LaTeX with no options to tap: report 28 arrived from a learner who had no
// way to enter sqrt(3)/2 that the choice grader would accept. The values are
// plain text now, so a tap resolves to the answer and a typed square root still
// reads as itself.
type specialAnglesGen struct{}

func (g *specialAnglesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return specialAngleProblem()
}

// trig.basics.reference_angle is also a multiple-choice concept, and it answered
// "60\u00b0" with no options, so a learner had nothing to tap and a typed 60 did
// not match the degree sign. The three acute answers are options now.
type referenceAngleGen struct{}

func (g *referenceAngleGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return referenceAngleProblem()
}

type graphSinGen struct{}

func (g *graphSinGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return sineGraphProblem()
}

// graphNumberValues and graphDegreeValues are the answer families a sine or
// cosine graph question draws from. Options stay inside one family so a period
// question never offers "-1" beside "360°".
var graphNumberValues = []string{"-1", "0", "1", "2"}
var graphDegreeValues = []string{"90°", "180°", "270°", "360°"}

// graphTrigProblem serves the value as options because a degree sign has no key
// on a phone keypad, and only one of these questions asks for a number a learner
// could type anyway.
func graphTrigProblem(question, answer, explanation string) generator.Problem {
	plain := strings.ReplaceAll(bare(answer), "^{\\circ}", "°")
	p := generator.MultipleChoice(question, plain, graphValueDistractors(plain)...)
	p.Explanation = explanation
	return p
}

func graphValueDistractors(want string) []string {
	pool := graphNumberValues
	if strings.HasSuffix(want, "°") {
		pool = graphDegreeValues
	}
	out := make([]string, 0, len(pool))
	for _, v := range pool {
		if v != want {
			out = append(out, v)
		}
	}
	return out
}

type graphCosGen struct{}

func (g *graphCosGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return cosineGraphProblem()
}

type periodGen struct{}

func (g *periodGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// B is chosen from divisors of 360 so the period in degrees is an exact integer.
	bDivisors := []int{1, 2, 3, 4, 5, 6, 8, 9, 10, 12}
	b := bDivisors[rand.Intn(len(bDivisors))]
	a := rand.Intn(7) + 2
	period := 360 / b

	fn := "sin"
	if rand.Intn(2) == 1 {
		fn = "cos"
	}

	q := fmt.Sprintf("For the trigonometric function \\(y = %d \\%s(%dx)\\), what is its period in degrees? (enter an integer)",
		a, fn, b)
	exp := fmt.Sprintf("The period of a standard \\(\\%s\\) function is \\(360^{\\circ}\\). The frequency multiplier \\(B = %d\\) compresses the period by a factor of %d: \\(360^{\\circ} / %d = %d^{\\circ}\\).",
		fn, b, b, b, period)

	return generator.Problem{
		Question:    q,
		Answer:      strconv.Itoa(period),
		Explanation: exp,
		Facts: map[string]string{
			"A":           strconv.Itoa(a),
			"B":           strconv.Itoa(b),
			"fn":          fn,
			"base_period": "360",
			"answer":      strconv.Itoa(period),
		},
	}
}

type inverseGen struct{}

func (g *inverseGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return inverseFunctionProblem()
}

// inverseAngleValues are the principal angles an inverse-function question can
// return, including the negative arcsin value.
var inverseAngleValues = []string{"-30°", "0°", "30°", "45°", "60°", "90°", "120°"}

func inverseAngleDistractors(want string) []string {
	pool := make([]string, 0, len(inverseAngleValues))
	for _, v := range inverseAngleValues {
		if v != want {
			pool = append(pool, v)
		}
	}
	rand.Shuffle(len(pool), func(i, j int) { pool[i], pool[j] = pool[j], pool[i] })
	return pool[:3]
}

type lawSinesGen struct{}

func (g *lawSinesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type anglePair struct {
		A, B int
	}
	pairs := []anglePair{
		{30, 45}, {30, 60}, {45, 60}, {30, 30},
		{45, 30}, {60, 30}, {30, 90}, {45, 90},
		{60, 45}, {60, 90}, {45, 45},
	}
	p := pairs[rand.Intn(len(pairs))]
	a := 6 + rand.Intn(10+int(ctx.Difficulty*30))
	sinA := math.Sin(float64(p.A) * math.Pi / 180.0)
	sinB := math.Sin(float64(p.B) * math.Pi / 180.0)
	b := float64(a) * sinB / sinA
	bRounded := math.Round(b*10.0) / 10.0
	bStr := fmt.Sprintf("%.1f", bRounded)
	return generator.Problem{
		Question:    fmt.Sprintf("In triangle ABC, \\(A = %d^{\\circ}\\), \\(B = %d^{\\circ}\\) and side \\(a = %d\\). Find side \\(b\\) using the law of sines.", p.A, p.B, a),
		Answer:      bStr,
		Explanation: fmt.Sprintf("By the law of sines, the ratio of a side to the sine of its opposite angle is constant: \\(\\frac{a}{\\sin(A)} = \\frac{b}{\\sin(B)}\\). Rearrange to solve for side b: \\(b = \\frac{a \\sin(B)}{\\sin(A)} = \\frac{%d \\sin(%d^{\\circ})}{\\sin(%d^{\\circ})}\\). Evaluating with trigonometric values and rounding to one decimal place yields %s.", a, p.B, p.A, bStr),
		Facts:       map[string]string{"a": strconv.Itoa(a), "A": strconv.Itoa(p.A), "B": strconv.Itoa(p.B), "b": bStr},
	}
}

type lawCosinesGen struct{}

func (g *lawCosinesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	C := 60
	a := 3 + rand.Intn(6+int(ctx.Difficulty*15))
	b := a + 1 + rand.Intn(5+int(ctx.Difficulty*10))
	crossTerm := a * b // 2 * a * b * cos(60) = a * b
	c2 := a*a + b*b - crossTerm
	cVal := math.Sqrt(float64(c2))
	cRounded := math.Round(cVal*10.0) / 10.0
	cStr := fmt.Sprintf("%.1f", cRounded)
	return generator.Problem{
		Question:    fmt.Sprintf("In triangle ABC, \\(a = %d\\), \\(b = %d\\) and \\(C = %d^{\\circ}\\). Find side \\(c\\) using the law of cosines.", a, b, C),
		Answer:      cStr,
		Explanation: fmt.Sprintf("By the law of cosines, the unknown side satisfies \\(c^{2} = a^{2} + b^{2} - 2ab\\cos(C)\\). Substitute the given parameters: \\(c^{2} = %d^{2} + %d^{2} - 2(%d)(%d)\\cos(%d^{\\circ}) = %d + %d - %d = %d\\). Taking the square root and rounding to one decimal place yields \\(c = \\sqrt{%d} \\approx %s\\).", a, b, a, b, C, a*a, b*b, crossTerm, c2, c2, cStr),
		Facts: map[string]string{
			"a":      strconv.Itoa(a),
			"b":      strconv.Itoa(b),
			"C":      strconv.Itoa(C),
			"a_sq":   strconv.Itoa(a * a),
			"b_sq":   strconv.Itoa(b * b),
			"two_ab": strconv.Itoa(crossTerm),
			"c_sq":   strconv.Itoa(c2),
			"c":      cStr,
		},
	}
}

// trigQA is one fixed written question with its answer and its reasoning. The
// variety pass keeps the strongest written questions from a pool and places
// parameterised forms beside them, so a branch is never a table of eight.
type trigQA struct {
	question string
	answer   string
	exp      string
}

// arctan: compute arctan values and properties
type arctanGen struct{}

// arctanExactArgs are the arguments whose arctangent is a whole number of
// degrees. There are only finitely many, so this half of the supply is a table;
// the reciprocal and telescoping identities below add the parameterised half.
var arctanExactArgs = []struct{ arg, deg string }{
	{"0", "0"},
	{"1", "45"},
	{"-1", "-45"},
	{"\\sqrt{3}", "60"},
	{"-\\sqrt{3}", "-60"},
	{"1/\\sqrt{3}", "30"},
	{"-1/\\sqrt{3}", "-30"},
	{"\\sqrt{3}/3", "30"},
	{"-\\sqrt{3}/3", "-30"},
}

func (g *arctanGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := trigOrdinal(ctx)
	switch o {
	case 1:
		if rand.Intn(3) > 0 {
			e := arctanExactArgs[rand.Intn(len(arctanExactArgs))]
			return trigProb(o,
				fmt.Sprintf("What is \\(\\arctan(%s)\\) in degrees?", e.arg),
				e.deg,
				fmt.Sprintf("\\(\\tan(%s^{\\circ}) = %s\\), so \\(\\arctan(%s) = %s^{\\circ}\\).", e.deg, e.arg, e.arg, e.deg))
		}
		// arctan(x) + arctan(1/x) is a right angle for every positive x,
		// because the two arguments are complementary slopes.
		k := 2 + rand.Intn(19)
		return trigProb(o,
			fmt.Sprintf("For \\(x = %d > 0\\), what is \\(\\arctan(x) + \\arctan(1/x)\\) in degrees?", k),
			"90",
			fmt.Sprintf("The two arguments are reciprocal slopes, so the angles are complementary: \\(\\arctan(%d) + \\arctan(1/%d) = 90^{\\circ}\\).", k, k))
	case 2:
		written := []trigQA{
			{"What is the range of \\(\\arctan(x)\\)? (in degrees)", "-90 to 90", "\\(\\arctan(x)\\) returns values in \\((-90^{\\circ}, 90^{\\circ})\\)."},
			{"Is \\(\\arctan(x)\\) an odd function? (yes/no)", "yes", "\\(\\arctan(-x) = -\\arctan(x)\\), so it is odd."},
			{"As \\(x \\to \\infty\\), \\(\\arctan(x)\\) approaches what value in degrees?", "90", "\\(\\lim_{x\\to\\infty} \\arctan(x) = 90^{\\circ}\\)."},
			{"As \\(x \\to -\\infty\\), \\(\\arctan(x)\\) approaches what value in degrees?", "-90", "\\(\\lim_{x\\to-\\infty} \\arctan(x) = -90^{\\circ}\\)."},
			{"Is \\(\\arctan(x)\\) increasing on its whole domain? (yes/no)", "yes", "Its derivative \\(1/(1+x^2)\\) is positive everywhere, so the graph never turns back."},
			{"How many horizontal asymptotes does \\(y = \\arctan(x)\\) have? (enter a number)", "2", "The graph approaches \\(y = 90^{\\circ}\\) and \\(y = -90^{\\circ}\\), so it has two."},
			{"What is \\(\\arctan(0)\\) in degrees?", "0", "\\(\\tan(0^{\\circ}) = 0\\), so \\(\\arctan(0) = 0^{\\circ}\\)."},
			{"What is the derivative of \\(\\arctan(x)\\) at \\(x = 0\\)?", "1", "The derivative is \\(1/(1+x^2)\\), which equals 1 at the origin."},
		}
		if rand.Intn(2) == 0 {
			e := written[rand.Intn(len(written))]
			return trigProb(o, e.question, e.answer, e.exp)
		}
		// tan(arctan(x)) returns x on the whole real line, unlike the
		// compositions that only return the principal value.
		k := 2 + rand.Intn(23)
		return trigProb(o,
			fmt.Sprintf("For \\(x = %d\\), what is \\(\\tan(\\arctan(x))\\)? (enter a number)", k),
			strconv.Itoa(k),
			fmt.Sprintf("\\(\\arctan\\) returns an angle whose tangent is the input, so \\(\\tan(\\arctan(%d)) = %d\\).", k, k))
	default:
		written := []trigQA{
			{"What is \\(\\arctan(1/2) + \\arctan(1/3)\\) in degrees?", "45", "The tangent-addition formula gives \\((1/2+1/3)/(1-1/6) = 1\\), so the sum is \\(45^{\\circ}\\)."},
			{"What is \\(\\arctan(1) + \\arctan(2) + \\arctan(3)\\) in degrees?", "180", "The first two terms add to \\(108.43^{\\circ}\\), and the third brings the total to a straight angle."},
			{"What is \\(\\arctan(1/2) + \\arctan(1/5) + \\arctan(1/8)\\) in degrees?", "45", "The tangent-addition formula applied twice leaves a tangent of 1, so the three angles total \\(45^{\\circ}\\)."},
			{"What is \\(\\arctan(\\sqrt{3}) - \\arctan(1/\\sqrt{3})\\) in degrees?", "30", "\\(60^{\\circ} - 30^{\\circ} = 30^{\\circ}\\)."},
			{"What is \\(\\lim_{x\\to 0} \\arctan(x)/x\\)?", "1", "Near the origin \\(\\arctan(x)\\) matches \\(x\\), so the ratio tends to 1."},
			{"What is the derivative of \\(\\arctan(x)\\) at \\(x = 1\\)? Give an exact fraction.", "1/2", "The derivative is \\(1/(1+x^2)\\), which equals \\(1/2\\) at \\(x = 1\\)."},
			{"If \\(\\arctan(x) = 30^{\\circ}\\), what is \\(x\\)?", "1/\\sqrt{3}", "\\(x = \\tan(30^{\\circ}) = 1/\\sqrt{3}\\)."},
			{"How many solutions does \\(\\arctan(x) = 45^{\\circ}\\) have? (enter a number)", "1", "\\(\\arctan\\) is one-to-one, so the horizontal line meets the graph once."},
		}
		if rand.Intn(2) == 0 {
			e := written[rand.Intn(len(written))]
			return trigProb(o, e.question, e.answer, e.exp)
		}
		k := 1 + rand.Intn(19)
		return trigProb(o,
			fmt.Sprintf("For \\(x = -%d < 0\\), what is \\(\\arctan(x) + \\arctan(1/x)\\) in degrees?", k),
			"-90",
			fmt.Sprintf("For a negative argument both principal angles are negative, and reciprocal slopes are still complementary: \\(\\arctan(-%d) + \\arctan(-1/%d) = -90^{\\circ}\\).", k, k))
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
		Facts:       map[string]string{"fn": f, "opp": strconv.Itoa(t.opp), "adj": strconv.Itoa(adj), "hyp": strconv.Itoa(t.hyp), "val": ans},
	}
}

// trigDegreeList renders a comma separated list of degree answers the way an
// explanation reads it: 30,150 becomes 30° and 150°.
func trigDegreeList(list string) string {
	parts := strings.Split(list, ",")
	for i, p := range parts {
		parts[i] = p + "^{\\circ}"
	}
	if len(parts) == 1 {
		return parts[0]
	}
	return strings.Join(parts[:len(parts)-1], ", ") + " and " + parts[len(parts)-1]
}

// basic trig equations
type trigEqBasicGen struct{}

func (g *trigEqBasicGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := trigOrdinal(ctx)
	type solution struct{ equation, answer string }
	pool := []solution{
		{"\\sin(x) = 0", "0,180"},
		{"\\sin(x) = 1", "90"},
		{"\\sin(x) = -1", "270"},
		{"\\sin(x) = 1/2", "30,150"},
		{"\\sin(x) = -1/2", "210,330"},
		{"\\sin(x) = \\sqrt{2}/2", "45,135"},
		{"\\sin(x) = -\\sqrt{2}/2", "225,315"},
		{"\\sin(x) = \\sqrt{3}/2", "60,120"},
		{"\\sin(x) = -\\sqrt{3}/2", "240,300"},
		// A doubled argument compresses two periods into the first turn, so the
		// solution set is a different one even for the same right-hand side.
		{"\\sin(2x) = 0", "0,90,180,270"},
		{"\\sin(2x) = 1", "45,225"},
		{"\\sin(2x) = -1", "135,315"},
		{"\\sin(2x) = 1/2", "15,75,195,255"},
		{"\\sin(2x) = -1/2", "105,165,285,345"},
		{"\\sin(3x) = 0", "0,60,120,180,240,300"},
		{"\\sin(3x) = 1", "30,150,270"},
	}
	switch o {
	case 2:
		pool = []solution{
			{"\\cos(x) = 0", "90,270"},
			{"\\cos(x) = 1", "0"},
			{"\\cos(x) = -1", "180"},
			{"\\cos(x) = 1/2", "60,300"},
			{"\\cos(x) = -1/2", "120,240"},
			{"\\cos(x) = \\sqrt{2}/2", "45,315"},
			{"\\cos(x) = -\\sqrt{2}/2", "135,225"},
			{"\\cos(x) = \\sqrt{3}/2", "30,330"},
			{"\\cos(x) = -\\sqrt{3}/2", "150,210"},
			{"\\cos(2x) = 0", "45,135,225,315"},
			{"\\cos(2x) = 1", "0,180"},
			{"\\cos(2x) = -1", "90,270"},
			{"\\cos(2x) = 1/2", "30,150,210,330"},
			{"\\cos(2x) = -1/2", "60,120,240,300"},
			{"\\cos(3x) = 0", "30,90,150,210,270,330"},
			{"\\cos(3x) = -1", "60,180,300"},
		}
	case 3:
		pool = []solution{
			{"\\tan(x) = 0", "0,180"},
			{"\\tan(x) = 1", "45,225"},
			{"\\tan(x) = -1", "135,315"},
			{"\\tan(x) = \\sqrt{3}", "60,240"},
			{"\\tan(x) = -\\sqrt{3}", "120,300"},
			{"\\tan(x) = 1/\\sqrt{3}", "30,210"},
			{"\\tan(x) = -1/\\sqrt{3}", "150,330"},
			{"\\tan(2x) = 0", "0,90,180,270"},
		}
		// A value outside the range of the function has no solution, which is a
		// different question from solving one that does.
		if rand.Intn(3) == 0 {
			type impossible struct{ equation string }
			out := []impossible{
				{"\\sin(x) = 2"}, {"\\sin(x) = -2"}, {"\\sin(x) = 3/2"},
				{"\\cos(x) = 2"}, {"\\cos(x) = -2"}, {"\\cos(x) = 3/2"},
			}
			e := out[rand.Intn(len(out))]
			return generator.Problem{
				Question:    fmt.Sprintf("How many solutions does \\(%s\\) have for \\(0^{\\circ} \\leq x < 360^{\\circ}\\)? (enter a number)", e.equation),
				Answer:      "0",
				Explanation: fmt.Sprintf("The sine and cosine functions never leave \\([-1, 1]\\), so \\(%s\\) is never reached and there are no solutions.", e.equation),
				Facts:       map[string]string{"equation": e.equation, "answer": "0"},
				KPOrdinal:   o,
			}
		}
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{
		Question: fmt.Sprintf("Solve \\(%s\\) for \\(0^{\\circ} \\leq x < 360^{\\circ}\\).", e.equation),
		Answer:   e.answer,
		Explanation: fmt.Sprintf("On the first turn, \\(%s\\) holds at \\(x = %s\\). These are the angles the unit circle assigns that value, and no other angle in \\([0^{\\circ}, 360^{\\circ})\\) shares it.",
			e.equation, trigDegreeList(e.answer)),
		Facts:     map[string]string{"equation": e.equation, "answer": e.answer},
		KPOrdinal: o,
	}
}

// homogeneous trig equations
type trigEqHomogeneousGen struct{}

func (g *trigEqHomogeneousGen) Grade(expected, userAnswer string) grader.Result {
	eParts := strings.Split(expected, ",")
	aParts := strings.Split(strings.Trim(strings.TrimSpace(userAnswer), "()[]"), ",")
	if len(eParts) > 1 && len(aParts) == len(eParts) {
		matched := true
		used := make([]bool, len(eParts))
		for _, a := range aParts {
			found := false
			aTrim := strings.TrimSpace(a)
			for i, e := range eParts {
				if !used[i] && strings.TrimSpace(e) == aTrim {
					used[i] = true
					found = true
					break
				}
			}
			if !found {
				matched = false
				break
			}
		}
		if matched {
			return grader.Result{Correct: true, Score: 1}
		}
	}
	router := grader.NewRouter()
	return router.Grade(grader.GradingTuple, expected, userAnswer)
}

// trigTurnList lists the angles a tangent value repeats at inside a whole
// number of turns. The tangent period is 180 degrees, so an interval of
// 360*turns degrees holds 2*turns solutions.
func trigTurnList(base, turns int) string {
	parts := make([]string, 0, 2*turns)
	for j := 0; j < 2*turns; j++ {
		parts = append(parts, strconv.Itoa(base+180*j))
	}
	return strings.Join(parts, ",")
}

// homogeneousEntry is one linear equation in sine and cosine, with the tangent
// value it reduces to and the first-turn angle that tangent starts from.
type homogeneousEntry struct {
	equation string
	slope    string
	base     int
}

var homogeneousEquations = []homogeneousEntry{
	{"\\sin(x) - \\cos(x) = 0", "1", 45},
	{"2\\sin(x) - 2\\cos(x) = 0", "1", 45},
	{"\\sin(x) + \\cos(x) = 0", "-1", 135},
	{"3\\sin(x) + 3\\cos(x) = 0", "-1", 135},
	{"\\sqrt{3}\\sin(x) - \\cos(x) = 0", "1/\\sqrt{3}", 30},
	{"3\\sin(x) - \\sqrt{3}\\cos(x) = 0", "1/\\sqrt{3}", 30},
	{"\\sin(x) - \\sqrt{3}\\cos(x) = 0", "\\sqrt{3}", 60},
	{"\\sqrt{3}\\sin(x) - 3\\cos(x) = 0", "\\sqrt{3}", 60},
	{"\\sin(x) + \\sqrt{3}\\cos(x) = 0", "-\\sqrt{3}", 120},
	{"\\sqrt{3}\\sin(x) + \\cos(x) = 0", "-1/\\sqrt{3}", 150},
	{"\\sin(x) = 0", "0", 0},
	{"2\\sin(x) = 0", "0", 0},
}

var homogeneousMethods = []trigQA{
	{"What substitution is used to solve \\(a\\sin(x) + b\\cos(x) = 0\\)?", "tan(x) = -b/a", "Divide both sides by \\(\\cos(x)\\): \\(a\\tan(x) + b = 0 \\to \\tan(x) = -b/a\\)."},
	{"What is the general method to solve \\(a\\sin(x) + b\\cos(x) = 0\\)?", "divide by cos(x)", "Dividing by \\(\\cos(x)\\) gives \\(a\\tan(x) + b = 0\\), which can be solved for \\(x\\)."},
}

var homogeneousGeneral = []trigQA{
	{"Solve \\(\\tan(x) = 1\\) for all real \\(x\\) in degrees. Enter the general solution in the form a + 180k.", "45 + 180k", "\\(\\tan\\) repeats every \\(180^{\\circ}\\), so every solution is \\(45^{\\circ}\\) plus a whole number of half-turns."},
	{"Solve \\(\\tan(x) = \\sqrt{3}\\) for all real \\(x\\) in degrees. Enter the general solution in the form a + 180k.", "60 + 180k", "The first solution is \\(60^{\\circ}\\) and the tangent period adds \\(180^{\\circ}\\) each time."},
	{"Solve \\(\\tan(x) = 1/\\sqrt{3}\\) for all real \\(x\\) in degrees. Enter the general solution in the form a + 180k.", "30 + 180k", "The first solution is \\(30^{\\circ}\\) and the tangent period adds \\(180^{\\circ}\\) each time."},
	{"Solve \\(\\tan(x) = -1\\) for all real \\(x\\) in degrees. Enter the general solution in the form a + 180k.", "135 + 180k", "The first positive solution is \\(135^{\\circ}\\), and every half-turn repeats it."},
	{"Solve \\(\\tan(x) = -\\sqrt{3}\\) for all real \\(x\\) in degrees. Enter the general solution in the form a + 180k.", "120 + 180k", "The first positive solution is \\(120^{\\circ}\\), and every half-turn repeats it."},
	{"Solve \\(\\tan(x) = -1/\\sqrt{3}\\) for all real \\(x\\) in degrees. Enter the general solution in the form a + 180k.", "150 + 180k", "The first positive solution is \\(150^{\\circ}\\), and every half-turn repeats it."},
	{"Solve \\(\\sin(x) = 0\\) for all real \\(x\\) in degrees. Enter the general solution in the form 180k.", "180k", "The sine vanishes at every whole number of half-turns."},
	{"Solve \\(\\cos(x) = 0\\) for all real \\(x\\) in degrees. Enter the general solution in the form a + 180k.", "90 + 180k", "The cosine vanishes at every quarter-turn followed by whole half-turns."},
}

func (g *trigEqHomogeneousGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := trigOrdinal(ctx)
	turns := 1
	switch o {
	case 2:
		turns = 2
	case 3:
		turns = 3
		if rand.Intn(2) == 0 {
			e := homogeneousGeneral[rand.Intn(len(homogeneousGeneral))]
			return trigProb(o, e.question, e.answer, e.exp)
		}
	}
	if o == 1 && rand.Intn(4) == 0 {
		e := homogeneousMethods[rand.Intn(len(homogeneousMethods))]
		return trigProb(o, e.question, e.answer, e.exp)
	}
	e := homogeneousEquations[rand.Intn(len(homogeneousEquations))]
	answer := trigTurnList(e.base, turns)
	limit := 360 * turns
	return generator.Problem{
		Question: fmt.Sprintf("Solve \\(%s\\) for \\(0^{\\circ} \\leq x < %d^{\\circ}\\). Enter degrees separated by commas.", e.equation, limit),
		Answer:   answer,
		Explanation: fmt.Sprintf("Dividing by \\(\\cos(x)\\) turns \\(%s\\) into \\(\\tan(x) = %s\\), whose first-turn angle is \\(%d^{\\circ}\\). The tangent repeats every \\(180^{\\circ}\\), so \\([0^{\\circ}, %d^{\\circ})\\) holds the solutions \\(x = %s\\).",
			e.equation, e.slope, e.base, limit, trigDegreeList(answer)),
		Facts:     map[string]string{"equation": e.equation, "slope": e.slope, "base": strconv.Itoa(e.base), "turns": strconv.Itoa(turns), "answer": answer},
		KPOrdinal: o,
	}
}

// hyperbolic: sinh and cosh
type sinhCoshGen struct{}

func (g *sinhCoshGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := trigOrdinal(ctx)
	switch o {
	case 1:
		written := []trigQA{
			{"What is sinh(0)? (enter a number)", "0", "sinh(0) = (e^{0} - e^{-0})/2 = (1 - 1)/2 = 0; \\(\\sinh(x) = (e^{x} - e^{-x})/2\\), the odd part of the exponential function."},
			{"What is cosh(0)? (enter a number)", "1", "cosh(0) = (e^{0} + e^{-0})/2 = (1 + 1)/2 = 1; \\(\\cosh(x) = (e^{x} + e^{-x})/2\\), the even part of the exponential function."},
			{"What is \\(\\cosh^{2}(x) - \\sinh^{2}(x)\\)?", "1", "\\(\\cosh^{2}(x) - \\sinh^{2}(x) = 1\\) (the hyperbolic analogue of \\(\\cos^{2}+\\sin^{2}=1\\))."},
			{"Is \\(\\sinh(x)\\) an even or odd function?", "odd", "\\(\\sinh(-x) = -\\sinh(x)\\), so \\(\\sinh\\) is odd."},
			{"Is \\(\\cosh(x)\\) an even or odd function?", "even", "\\(\\cosh(-x) = \\cosh(x)\\), so \\(\\cosh\\) is even."},
			{"What is the derivative of \\(\\sinh(x)\\)?", "cosh(x)", "\\(d/dx \\sinh(x) = \\cosh(x)\\)."},
			{"What is the derivative of \\(\\cosh(x)\\)?", "sinh(x)", "\\(d/dx \\cosh(x) = \\sinh(x)\\)."},
			{"What is \\(\\sinh(-x)\\) in terms of \\(\\sinh(x)\\)?", "-sinh(x)", "\\(\\sinh\\) is odd, so the sign comes out."},
			{"What is \\(\\cosh(-x)\\) in terms of \\(\\cosh(x)\\)?", "cosh(x)", "\\(\\cosh\\) is even, so the value is unchanged."},
		}
		if rand.Intn(2) == 0 {
			e := written[rand.Intn(len(written))]
			return trigProb(o, e.question, e.answer, e.exp)
		}
		// sinh(ln k) = (k - 1/k)/2, which is rational for every whole k.
		k := 2 + rand.Intn(11)
		return trigProb(o,
			fmt.Sprintf("What is \\(\\sinh(\\ln %d)\\)? Give an exact fraction.", k),
			trigFrac(k*k-1, 2*k),
			fmt.Sprintf("\\(\\sinh(\\ln %d) = (e^{\\ln %d} - e^{-\\ln %d})/2 = (%d - 1/%d)/2 = %s\\).", k, k, k, k, k, trigFrac(k*k-1, 2*k)))
	case 2:
		written := []trigQA{
			{"What is \\(\\cosh^{2}(x) - \\sinh^{2}(x)\\)?", "1", "\\(\\cosh^{2}(x) - \\sinh^{2}(x) = 1\\) for every real \\(x\\)."},
			{"What is the second derivative of \\(\\sinh(x)\\)?", "sinh(x)", "Differentiating \\(\\cosh(x)\\) returns \\(\\sinh(x)\\)."},
			{"What is the second derivative of \\(\\cosh(x)\\)?", "cosh(x)", "Differentiating \\(\\sinh(x)\\) returns \\(\\cosh(x)\\)."},
			{"What is the indefinite integral of \\(\\sinh(x)\\)?", "cosh(x)", "Since \\(d/dx \\cosh(x) = \\sinh(x)\\), the integral is \\(\\cosh(x) + c\\)."},
			{"What is the indefinite integral of \\(\\cosh(x)\\)?", "sinh(x)", "Since \\(d/dx \\sinh(x) = \\cosh(x)\\), the integral is \\(\\sinh(x) + c\\)."},
			{"What is \\(\\sinh(2x)\\) in terms of \\(\\sinh(x)\\) and \\(\\cosh(x)\\)?", "2*sinh(x)*cosh(x)", "\\(\\sinh(2x) = 2\\sinh(x)\\cosh(x)\\)."},
			{"What is \\(\\cosh(2x)\\) in terms of \\(\\sinh(x)\\) and \\(\\cosh(x)\\)?", "cosh(x)^2 + sinh(x)^2", "\\(\\cosh(2x) = \\cosh^{2}(x) + \\sinh^{2}(x)\\)."},
		}
		if rand.Intn(2) == 0 {
			e := written[rand.Intn(len(written))]
			return trigProb(o, e.question, e.answer, e.exp)
		}
		// The chain rule on sinh(kx) and on cosh(kx) carries the coefficient.
		k := 2 + rand.Intn(11)
		if rand.Intn(2) == 0 {
			return trigProb(o,
				fmt.Sprintf("What is the derivative of \\(\\sinh(%dx)\\)?", k),
				fmt.Sprintf("%d cosh(%dx)", k, k),
				fmt.Sprintf("The chain rule multiplies by the inner coefficient %d, so the derivative is \\(%d\\cosh(%dx)\\).", k, k, k))
		}
		return trigProb(o,
			fmt.Sprintf("What is the derivative of \\(\\cosh(%dx)\\)?", k),
			fmt.Sprintf("%d sinh(%dx)", k, k),
			fmt.Sprintf("The chain rule multiplies by the inner coefficient %d, so the derivative is \\(%d\\sinh(%dx)\\).", k, k, k))
	default:
		written := []trigQA{
			{"What is \\(\\cosh(2x)\\) in terms of \\(\\cosh(x)\\) alone?", "2*cosh(x)^2 - 1", "\\(\\cosh(2x) = 2\\cosh^{2}(x) - 1\\)."},
			{"What is \\(\\sinh(x+y)\\) in terms of the hyperbolic functions of \\(x\\) and \\(y\\)?", "sinh(x)cosh(y) + cosh(x)sinh(y)", "\\(\\sinh(x+y) = \\sinh(x)\\cosh(y) + \\cosh(x)\\sinh(y)\\)."},
			{"What is \\(\\cosh(x+y)\\) in terms of the hyperbolic functions of \\(x\\) and \\(y\\)?", "cosh(x)cosh(y) + sinh(x)sinh(y)", "\\(\\cosh(x+y) = \\cosh(x)\\cosh(y) + \\sinh(x)\\sinh(y)\\)."},
		}
		if rand.Intn(4) == 0 {
			e := written[rand.Intn(len(written))]
			return trigProb(o, e.question, e.answer, e.exp)
		}
		// cosh(ln k) is rational, and so is cosh(2x) built from the pair
		// sinh(ln k) and cosh(ln k).
		k := 2 + rand.Intn(11)
		switch rand.Intn(2) {
		case 0:
			return trigProb(o,
				fmt.Sprintf("What is \\(\\cosh(\\ln %d)\\)? Give an exact fraction.", k),
				trigFrac(k*k+1, 2*k),
				fmt.Sprintf("\\(\\cosh(\\ln %d) = (%d + 1/%d)/2 = %s\\).", k, k, k, trigFrac(k*k+1, 2*k)))
		default:
			sinhPart := trigFrac(k*k-1, 2*k)
			coshPart := trigFrac(k*k+1, 2*k)
			return trigProb(o,
				fmt.Sprintf("If \\(\\sinh x = %s\\) and \\(\\cosh x = %s\\), what is \\(\\cosh(2x)\\)? Give an exact fraction.", sinhPart, coshPart),
				trigFrac(k*k*k*k+1, 2*k*k),
				fmt.Sprintf("\\(\\cosh(2x) = \\cosh^{2}(x) + \\sinh^{2}(x) = %s\\).", trigFrac(k*k*k*k+1, 2*k*k)))
		}
	}
}

// hyperbolic: tanh and coth
type tanhCothGen struct{}

func (g *tanhCothGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := trigOrdinal(ctx)
	switch o {
	case 1:
		written := []trigQA{
			{"What is the definition of \\(\\tanh(x)\\)?", "sinh(x)/cosh(x)", "\\(\\tanh(x) = \\sinh(x)/\\cosh(x) = (e^{x} - e^{-x})/(e^{x} + e^{-x})\\)."},
			{"What is \\(\\tanh(0)\\)?", "0", "\\(\\tanh(0) = \\sinh(0)/\\cosh(0) = 0/1 = 0\\)."},
			{"As \\(x \\to \\infty\\), \\(\\tanh(x)\\) approaches what value?", "1", "\\(\\lim_{x\\to\\infty} \\tanh(x) = 1\\) because \\(e^{x}\\) dominates \\(e^{-x}\\)."},
			{"As \\(x \\to -\\infty\\), \\(\\tanh(x)\\) approaches what value?", "-1", "\\(\\lim_{x\\to-\\infty} \\tanh(x) = -1\\) because \\(e^{-x}\\) dominates \\(e^{x}\\)."},
			{"How many horizontal asymptotes does \\(y = \\tanh(x)\\) have? (enter a number)", "2", "\\(\\tanh(x)\\) maps real numbers to the open interval \\((-1, 1)\\), with horizontal asymptotes \\(y = 1\\) and \\(y = -1\\)."},
			{"Is \\(\\tanh(x)\\) an even or odd function?", "odd", "\\(\\tanh(-x) = -\\tanh(x)\\), so \\(\\tanh\\) is odd."},
			{"What is \\(\\tanh(-x)\\) in terms of \\(\\tanh(x)\\)?", "-tanh(x)", "\\(\\tanh\\) is odd, so the sign comes out."},
			{"What is the limit of \\(\\tanh(x)\\) as \\(x \\to 0\\)?", "0", "The hyperbolic tangent is continuous at the origin and \\(\\tanh(0) = 0\\)."},
		}
		if rand.Intn(2) == 0 {
			e := written[rand.Intn(len(written))]
			return trigProb(o, e.question, e.answer, e.exp)
		}
		// tanh(ln k) = (k^2-1)/(k^2+1), which is rational for every whole k.
		k := 2 + rand.Intn(11)
		return trigProb(o,
			fmt.Sprintf("What is \\(\\tanh(\\ln %d)\\)? Give an exact fraction.", k),
			trigFrac(k*k-1, k*k+1),
			fmt.Sprintf("\\(\\tanh(\\ln %d) = (%d - 1/%d)/(%d + 1/%d) = %s\\).", k, k, k, k, k, trigFrac(k*k-1, k*k+1)))
	case 2:
		written := []trigQA{
			{"What is the definition of \\(\\coth(x)\\)?", "cosh(x)/sinh(x)", "\\(\\coth(x) = \\cosh(x)/\\sinh(x) = 1/\\tanh(x)\\)."},
			{"What is \\(\\coth(x)\\) at \\(x = 0\\)?", "undefined", "\\(\\sinh(0) = 0\\), so \\(\\coth\\) has a vertical asymptote at the origin and no value there."},
			{"At which \\(x\\) in \\((-1, 1)\\) is \\(\\coth(x)\\) undefined? (enter a number)", "0", "\\(\\coth\\) is undefined only where \\(\\sinh(x) = 0\\), which is the single point \\(x = 0\\) in that interval."},
			{"Is \\(\\coth(x)\\) an even or odd function?", "odd", "\\(\\coth(-x) = -\\coth(x)\\), so \\(\\coth\\) is odd."},
			{"As \\(x \\to \\infty\\), \\(\\coth(x)\\) approaches what value?", "1", "Both \\(\\cosh(x)\\) and \\(\\sinh(x)\\) grow like \\(e^{x}/2\\), so their ratio tends to 1."},
			{"As \\(x \\to 0^{+}\\), \\(\\coth(x)\\) approaches what value?", "+infinity", "The numerator tends to 1 while the denominator tends to zero from above."},
			{"What is \\(1 - \\coth^{2}(x)\\)?", "-csch(x)^2", "\\(\\coth^{2}(x) - 1 = \\operatorname{csch}^{2}(x)\\), so \\(1 - \\coth^{2}(x) = -\\operatorname{csch}^{2}(x)\\)."},
		}
		if rand.Intn(2) == 0 {
			e := written[rand.Intn(len(written))]
			return trigProb(o, e.question, e.answer, e.exp)
		}
		k := 2 + rand.Intn(11)
		return trigProb(o,
			fmt.Sprintf("What is \\(\\coth(\\ln %d)\\)? Give an exact fraction.", k),
			trigFrac(k*k+1, k*k-1),
			fmt.Sprintf("\\(\\coth(\\ln %d) = (%d + 1/%d)/(%d - 1/%d) = %s\\).", k, k, k, k, k, trigFrac(k*k+1, k*k-1)))
	default:
		written := []trigQA{
			{"What is \\(1 - \\tanh^{2}(x)\\)?", "sech(x)^2", "\\(1 - \\tanh^{2}(x) = \\operatorname{sech}^{2}(x) = 1/\\cosh^{2}(x)\\)."},
			{"What is \\(\\tanh(x+y)\\) in terms of \\(\\tanh(x)\\) and \\(\\tanh(y)\\)?", "(tanh(x)+tanh(y))/(1+tanh(x)tanh(y))", "\\(\\tanh(x+y) = \\frac{\\tanh(x)+\\tanh(y)}{1+\\tanh(x)\\tanh(y)}\\)."},
			{"What is \\(\\coth(x+y)\\) in terms of \\(\\coth(x)\\) and \\(\\coth(y)\\)?", "(coth(x)coth(y)+1)/(coth(x)+coth(y))", "\\(\\coth(x+y) = \\frac{\\coth(x)\\coth(y)+1}{\\coth(x)+\\coth(y)}\\)."},
			{"What is the derivative of \\(\\tanh(x)\\)?", "sech(x)^2", "\\(d/dx \\tanh(x) = \\operatorname{sech}^{2}(x)\\)."},
			{"What is the derivative of \\(\\coth(x)\\)?", "-csch(x)^2", "\\(d/dx \\coth(x) = -\\operatorname{csch}^{2}(x)\\)."},
		}
		if rand.Intn(4) == 0 {
			e := written[rand.Intn(len(written))]
			return trigProb(o, e.question, e.answer, e.exp)
		}
		// The chain rule on tanh(kx) and on coth(kx) carries the coefficient and
		// the sign, so both are exact written answers.
		k := 2 + rand.Intn(11)
		if rand.Intn(2) == 0 {
			return trigProb(o,
				fmt.Sprintf("What is the derivative of \\(\\tanh(%dx)\\)?", k),
				fmt.Sprintf("%d sech(%dx)^2", k, k),
				fmt.Sprintf("The chain rule multiplies by the inner coefficient %d, so the derivative is \\(%d\\operatorname{sech}^{2}(%dx)\\).", k, k, k))
		}
		return trigProb(o,
			fmt.Sprintf("What is the derivative of \\(\\coth(%dx)\\)?", k),
			fmt.Sprintf("-%d csch(%dx)^2", k, k),
			fmt.Sprintf("The chain rule multiplies by the inner coefficient %d and the derivative of coth is negative, so the answer is \\(-%d\\operatorname{csch}^{2}(%dx)\\).", k, k, k))
	}
}

// trig identities
type trigIdentGen struct{}

func (g *trigIdentGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := trigOrdinal(ctx)
	switch o {
	case 1:
		written := []trigQA{
			{"Simplify \\(\\sin^{2}(x) + \\cos^{2}(x)\\).", "1", "\\(\\sin^{2}(x) + \\cos^{2}(x) = 1\\) (the Pythagorean identity)."},
			{"Simplify \\(1 + \\cot^{2}(x)\\).", "csc(x)^2", "\\(1 + \\cot^{2}(x) = \\csc^{2}(x)\\), derived from \\(\\sin^{2}+\\cos^{2}=1\\) divided by \\(\\sin^{2}\\)."},
			{"Simplify \\(1 + \\tan^{2}(x)\\).", "sec(x)^2", "\\(1 + \\tan^{2}(x) = \\sec^{2}(x)\\), derived from \\(\\sin^{2}+\\cos^{2}=1\\) divided by \\(\\cos^{2}\\)."},
			{"What is 1 + tan(x)^2 at x = 45 degrees? (enter a number)", "2", "At 45 degrees, tan = 1, so 1 + tan^2 = 1 + 1 = 2; in general \\(1 + \\tan^{2}(x) = \\sec^{2}(x)\\), derived from \\(\\sin^{2}+\\cos^{2}=1\\) divided by \\(\\cos^{2}\\)."},
		}
		if rand.Intn(3) == 0 {
			e := written[rand.Intn(len(written))]
			return trigProb(o, e.question, e.answer, e.exp)
		}
		// A Pythagorean triple supplies an exact sine, so the identity returns
		// an exact rational square rather than a rounded decimal.
		a, b, c := pythagoreanPair(9)
		return trigProb(o,
			fmt.Sprintf("If \\(\\sin(x) = %d/%d\\) and \\(x\\) is acute, what is \\(\\cos^{2}(x)\\)? Give an exact fraction.", a, c),
			trigFrac(b*b, c*c),
			fmt.Sprintf("\\(\\cos^{2}(x) = 1 - \\sin^{2}(x) = 1 - %d/%d = %d/%d\\), which reduces to %s.",
				a*a, c*c, c*c-a*a, c*c, trigFrac(b*b, c*c)))
	case 2:
		written := []trigQA{
			{"What is \\(\\sin(-x)\\) in terms of \\(\\sin(x)\\)?", "-sin(x)", "\\(\\sin(-x) = -\\sin(x)\\) (sine is odd)."},
			{"What is \\(\\cos(-x)\\) in terms of \\(\\cos(x)\\)?", "cos(x)", "\\(\\cos(-x) = \\cos(x)\\) (cosine is even)."},
			{"What is \\(\\sin(2x)\\) in terms of \\(\\sin(x)\\) and \\(\\cos(x)\\)?", "2*sin(x)*cos(x)", "\\(\\sin(2x) = 2\\sin(x)\\cos(x)\\) (double angle formula)."},
			{"What is cos(2x) at x = 90 degrees? (enter a number)", "-1", "cos(180 degrees) = -1; in general \\(\\cos(2x) = 2\\cos^{2}(x) - 1 = \\cos^{2}(x) - \\sin^{2}(x) = 1 - 2\\sin^{2}(x)\\)."},
			{"What is \\(\\cos(2x)\\) in terms of \\(\\cos(x)\\) alone?", "2*cos(x)^2 - 1", "\\(\\cos(2x) = 2\\cos^{2}(x) - 1\\)."},
		}
		if rand.Intn(3) == 0 {
			e := written[rand.Intn(len(written))]
			return trigProb(o, e.question, e.answer, e.exp)
		}
		a, b, c := pythagoreanPair(9)
		return trigProb(o,
			fmt.Sprintf("If \\(\\sin(x) = %d/%d\\) and \\(x\\) is acute, what is \\(\\sin(2x)\\)? Give an exact fraction.", a, c),
			trigFrac(2*a*b, c*c),
			fmt.Sprintf("\\(\\sin(2x) = 2\\sin(x)\\cos(x)\\), and the triple gives \\(\\cos(x) = %d/%d\\), so the product is \\(2 \\cdot %d/%d \\cdot %d/%d = %s\\).",
				b, c, a, c, b, c, trigFrac(2*a*b, c*c)))
	default:
		written := []trigQA{
			{"Simplify \\(\\sin(x)\\cos(y) + \\cos(x)\\sin(y)\\).", "sin(x+y)", "\\(\\sin(x+y) = \\sin(x)\\cos(y) + \\cos(x)\\sin(y)\\) (addition formula)."},
			{"Simplify \\(\\cos(x)\\cos(y) - \\sin(x)\\sin(y)\\).", "cos(x+y)", "\\(\\cos(x+y) = \\cos(x)\\cos(y) - \\sin(x)\\sin(y)\\) (addition formula)."},
			{"Simplify \\(\\sin(x)\\cos(y) - \\cos(x)\\sin(y)\\).", "sin(x-y)", "\\(\\sin(x-y) = \\sin(x)\\cos(y) - \\cos(x)\\sin(y)\\)."},
			{"Simplify \\(\\cos(x)\\cos(y) + \\sin(x)\\sin(y)\\).", "cos(x-y)", "\\(\\cos(x-y) = \\cos(x)\\cos(y) + \\sin(x)\\sin(y)\\)."},
			{"What is \\(\\tan(2x)\\) in terms of \\(\\tan(x)\\)?", "2*tan(x)/(1-tan(x)^2)", "\\(\\tan(2x) = \\frac{2\\tan(x)}{1-\\tan^{2}(x)}\\)."},
		}
		if rand.Intn(3) == 0 {
			e := written[rand.Intn(len(written))]
			return trigProb(o, e.question, e.answer, e.exp)
		}
		a, b, c := pythagoreanPair(9)
		return trigProb(o,
			fmt.Sprintf("If \\(\\sin(x) = %d/%d\\) and \\(x\\) is acute, what is \\(\\cos(2x)\\)? Give an exact fraction.", a, c),
			trigFrac(b*b-a*a, c*c),
			fmt.Sprintf("\\(\\cos(2x) = \\cos^{2}(x) - \\sin^{2}(x) = %d/%d - %d/%d = %s\\).",
				b*b, c*c, a*a, c*c, trigFrac(b*b-a*a, c*c)))
	}
}

// basic trig inequalities
type trigIneqGen struct{}

func (g *trigIneqGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := trigOrdinal(ctx)
	switch o {
	case 1:
		written := []trigQA{
			{"What is the maximum value of \\(\\sin(x)\\)?", "1", "The maximum of \\(\\sin(x)\\) is 1, achieved at \\(x = 90^{\\circ}\\)."},
			{"What is the minimum value of \\(\\sin(x)\\)?", "-1", "The minimum of \\(\\sin(x)\\) is -1, achieved at \\(x = 270^{\\circ}\\)."},
			{"What is the maximum value of \\(\\cos(x)\\)?", "1", "The maximum of \\(\\cos(x)\\) is 1, achieved at \\(x = 0^{\\circ}\\)."},
			{"What is the minimum value of \\(\\cos(x)\\)?", "-1", "The minimum of \\(\\cos(x)\\) is -1, achieved at \\(x = 180^{\\circ}\\)."},
		}
		if rand.Intn(2) == 0 {
			e := written[rand.Intn(len(written))]
			return trigProb(o, e.question, e.answer, e.exp)
		}
		// A vertical shift or a coefficient moves the extreme value by an
		// amount the learner has to read off the expression.
		k := 1 + rand.Intn(20)
		switch rand.Intn(4) {
		case 0:
			return trigProb(o,
				fmt.Sprintf("What is the maximum value of \\(y = %d + \\sin(x)\\)?", k),
				strconv.Itoa(k+1),
				fmt.Sprintf("The sine never exceeds 1, so the maximum is \\(%d + 1 = %d\\).", k, k+1))
		case 1:
			return trigProb(o,
				fmt.Sprintf("What is the minimum value of \\(y = %d + \\cos(x)\\)?", k),
				strconv.Itoa(k-1),
				fmt.Sprintf("The cosine never drops below -1, so the minimum is \\(%d - 1 = %d\\).", k, k-1))
		case 2:
			return trigProb(o,
				fmt.Sprintf("What is the maximum value of \\(y = %d\\sin(x)\\)?", k),
				strconv.Itoa(k),
				fmt.Sprintf("The positive coefficient scales the peak of the sine, so the maximum is %d.", k))
		default:
			return trigProb(o,
				fmt.Sprintf("What is the minimum value of \\(y = -%d\\cos(x)\\)?", k),
				strconv.Itoa(-k),
				fmt.Sprintf("The cosine reaches 1 at the origin, so \\(-%d\\cos(x)\\) reaches \\(-%d\\) there.", k, k))
		}
	case 2:
		fn := []string{"\\sin", "\\cos", "\\tan"}[rand.Intn(3)]
		angle := 1 + rand.Intn(359)
		for angle%90 == 0 {
			angle = 1 + rand.Intn(359)
		}
		positive := false
		switch fn {
		case "\\sin":
			positive = angle < 180
		case "\\cos":
			positive = angle < 90 || angle > 270
		default:
			positive = angle < 90 || (angle > 180 && angle < 270)
		}
		word := "negative"
		if positive {
			word = "positive"
		}
		return generator.Problem{
			Question:    fmt.Sprintf("For \\(x = %d^{\\circ}\\), is \\(%s(x)\\) positive or negative?", angle, fn),
			Answer:      word,
			Explanation: fmt.Sprintf("An angle of \\(%d^{\\circ}\\) terminates in quadrant %d, where \\(%s\\) is %s. The sign of a trigonometric function is fixed by the quadrant, not by the size of the angle.", angle, angle/90+1, strings.TrimPrefix(fn, "\\"), word),
			Facts:       map[string]string{"fn": fn, "angle": strconv.Itoa(angle), "answer": word},
			KPOrdinal:   o,
		}
	default:
		type monotone struct {
			fn       string
			lo, hi   int
			behavior string
		}
		pool := []monotone{
			{"\\sin", 0, 90, "increasing"},
			{"\\sin", 90, 270, "decreasing"},
			{"\\sin", 270, 360, "increasing"},
			{"\\sin", 0, 180, "neither"},
			{"\\cos", 0, 180, "decreasing"},
			{"\\cos", 180, 360, "increasing"},
			{"\\cos", 0, 90, "decreasing"},
			{"\\cos", 90, 180, "decreasing"},
			{"\\cos", 0, 360, "neither"},
			{"\\tan", 0, 90, "increasing"},
			{"\\tan", 90, 180, "increasing"},
			{"\\tan", 180, 270, "increasing"},
		}
		e := pool[rand.Intn(len(pool))]
		return generator.Problem{
			Question:    fmt.Sprintf("For \\(%d^{\\circ} < x < %d^{\\circ}\\), is \\(%s(x)\\) increasing, decreasing or neither?", e.lo, e.hi, e.fn),
			Answer:      e.behavior,
			Explanation: fmt.Sprintf("On \\((%d^{\\circ}, %d^{\\circ})\\), \\(%s(x)\\) is %s. The sign of the derivative decides the direction, and a curve that turns inside the interval is neither increasing nor decreasing over the whole of it.", e.lo, e.hi, e.fn, e.behavior),
			Facts:       map[string]string{"fn": e.fn, "lo": strconv.Itoa(e.lo), "hi": strconv.Itoa(e.hi), "answer": e.behavior},
			KPOrdinal:   o,
		}
	}
}

// ----- P3.1 trig functions wave: each function as a real function of a real
// variable (domain, parity, asymptotes, limits, derivatives, integrals).
// Easy pools are fully typed (no yes/no recognition); hard rows unlock at
// scale > 3. All answers are ASCII-safe short forms.

type funcEntry struct {
	question string
	answer   string
	exp      string
}

// functionChoices gives the options for the prose answers in the function pools.
// A word cannot survive the math field: report #198 recorded a whole sentence
// with every space stripped. Numeric answers are left as typed input.
var functionChoices = map[string][]string{
	"odd":           {"odd", "even", "neither"},
	"even":          {"odd", "even", "neither"},
	"neither":       {"odd", "even", "neither"},
	"increasing":    {"increasing", "decreasing", "neither"},
	"decreasing":    {"increasing", "decreasing", "neither"},
	"all reals":     {"all reals", "[-1, 1]", "[0, pi]", "(-pi/2, pi/2)"},
	"[-pi/2,pi/2]":  {"[-pi/2,pi/2]", "[-1, 1]", "[0, pi]", "all reals"},
	"[-pi/2, pi/2]": {"[-pi/2, pi/2]", "[-1, 1]", "[0, pi]", "all reals"},
	"[0, pi]":       {"[0, pi]", "[-1, 1]", "all reals", "(-pi/2, pi/2)"},
	"(0, pi)":       {"(0, pi)", "[-1, 1]", "all reals", "[0, pi]"},
}

func funcPool(ctx generator.GeneratorContext, easy, hard []funcEntry) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	if options, ok := functionChoices[e.answer]; ok {
		p := generator.MultipleChoice(strings.ReplaceAll(e.question, " (type one word)", ""), e.answer, options...)
		p.Explanation = e.exp
		return p
	}
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.exp}
}

// sine as a real function
type sineFuncGen struct{}

func (g *sineFuncGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	easy := []funcEntry{
		{"What is the period of \\(y = \\sin(x)\\)? (type like 2pi)", "2pi", "The sine wave repeats every \\(2\\pi\\) radians."},
		{"What is the maximum value of \\(y = \\sin(x)\\)? (enter a number)", "1", "The sine function oscillates between -1 and 1; the maximum is 1."},
		{"What is the minimum value of \\(y = \\sin(x)\\)? (enter a number)", "-1", "The sine function oscillates between -1 and 1; the minimum is -1."},
		{"What is \\(\\sin(0)\\)? (enter a number)", "0", "\\(\\sin(0) = 0\\): the wave starts at the origin."},
		{"Is \\(\\sin(x)\\) odd, even, or neither? (type one word)", "odd", "\\(\\sin(-x) = -\\sin(x)\\), so sine is odd."},
		{"What is the derivative of \\(\\sin(x)\\)?", "cos(x)", "\\(d/dx \\sin(x) = \\cos(x)\\)."},
		{"What is the indefinite integral of \\(\\sin(x)\\)?", "-cos(x)", "Since \\(d/dx(-\\cos(x)) = \\sin(x)\\), the integral is \\(-\\cos(x) + c\\)."},
		{"What is the limit of \\(\\sin(x)/x\\) as \\(x \\to 0\\)? (enter a number)", "1", "Near the origin the sine curve matches \\(y = x\\), so the ratio tends to 1."},
	}
	hard := []funcEntry{
		{"What is \\(\\sin(\\pi/2)\\)? (enter a number)", "1", "\\(\\sin(\\pi/2) = 1\\): the first maximum of the wave."},
		{"What is the second derivative of \\(\\sin(x)\\)?", "-sin(x)", "Differentiating twice: \\(\\cos(x) \\to -\\sin(x)\\)."},
		{"What is the first term of the Maclaurin series of \\(\\sin(x)\\)?", "x", "\\(\\sin(x) = x - x^3/3! + \\cdots\\): only odd powers appear."},
		{"What is \\(\\sin(\\pi)\\)? (enter a number)", "0", "\\(\\sin(\\pi) = 0\\): a root at every \\(n\\pi\\)."},
	}
	return funcPool(ctx, easy, hard)
}

// cosine as a real function
type cosineFuncGen struct{}

func (g *cosineFuncGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	easy := []funcEntry{
		{"What is the period of \\(y = \\cos(x)\\)? (type like 2pi)", "2pi", "The cosine wave repeats every \\(2\\pi\\) radians."},
		{"What is the maximum value of \\(y = \\cos(x)\\)? (enter a number)", "1", "The cosine function oscillates between -1 and 1; the maximum is 1."},
		{"What is the minimum value of \\(y = \\cos(x)\\)? (enter a number)", "-1", "The cosine function oscillates between -1 and 1; the minimum is -1."},
		{"What is \\(\\cos(0)\\)? (enter a number)", "1", "\\(\\cos(0) = 1\\): the wave starts at its peak."},
		{"Is \\(\\cos(x)\\) odd, even, or neither? (type one word)", "even", "\\(\\cos(-x) = \\cos(x)\\), so cosine is even."},
		{"What is the derivative of \\(\\cos(x)\\)?", "-sin(x)", "\\(d/dx \\cos(x) = -\\sin(x)\\)."},
		{"What is the indefinite integral of \\(\\cos(x)\\)?", "sin(x)", "Since \\(d/dx(\\sin(x)) = \\cos(x)\\), the integral is \\(\\sin(x) + c\\)."},
		{"What is the limit of \\((1 - \\cos(x))/x\\) as \\(x \\to 0\\)? (enter a number)", "0", "Near the origin \\(1 - \\cos(x)\\) vanishes faster than \\(x\\), so the ratio tends to 0."},
	}
	hard := []funcEntry{
		{"What is \\(\\cos(\\pi)\\)? (enter a number)", "-1", "\\(\\cos(\\pi) = -1\\): the first minimum of the wave."},
		{"What is the second derivative of \\(\\cos(x)\\)?", "-cos(x)", "Differentiating twice: \\(-\\sin(x) \\to -\\cos(x)\\)."},
		{"What is the constant term of the Maclaurin series of \\(\\cos(x)\\)? (enter a number)", "1", "\\(\\cos(x) = 1 - x^2/2! + \\cdots\\): only even powers appear."},
		{"Evaluate the integral of \\(\\cos(x)\\) from 0 to \\(\\pi/2\\). (enter a number)", "1", "\\(\\sin(\\pi/2) - \\sin(0) = 1 - 0 = 1\\)."},
	}
	return funcPool(ctx, easy, hard)
}

// tangent as a real function
type tangentFuncGen struct{}

func (g *tangentFuncGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	easy := []funcEntry{
		{"What is the period of \\(y = \\tan(x)\\)? (type like pi)", "pi", "Tangent repeats every \\(\\pi\\) radians, with one branch per period."},
		{"What is \\(\\tan(0)\\)? (enter a number)", "0", "\\(\\tan(0) = \\sin(0)/\\cos(0) = 0/1 = 0\\)."},
		{"Is \\(\\tan(x)\\) odd, even, or neither? (type one word)", "odd", "\\(\\tan(-x) = -\\tan(x)\\), so tangent is odd."},
		{"What is the derivative of \\(\\tan(x)\\)?", "sec^2(x)", "\\(d/dx \\tan(x) = \\sec^{2}(x)\\)."},
		{"What is the indefinite integral of \\(\\tan(x)\\)?", "-ln|cos(x)|", "The integral is \\(-\\ln|\\cos(x)| + c\\)."},
		{"At which \\(x\\) in \\([0, \\pi)\\) is \\(\\tan(x)\\) undefined? (type like pi/2)", "pi/2", "Cosine vanishes at \\(\\pi/2\\), so tangent has a vertical asymptote there."},
		{"What is the range of \\(\\tan(x)\\)? (type: all reals)", "all reals", "Each branch of tangent covers every real value."},
		{"Is \\(\\tan(x)\\) increasing or decreasing on each branch? (type one word)", "increasing", "Tangent is strictly increasing on every interval \\((-\\pi/2 + k\\pi, \\pi/2 + k\\pi)\\)."},
	}
	hard := []funcEntry{
		{"What is \\(\\tan(\\pi/4)\\)? (enter a number)", "1", "\\(\\tan(\\pi/4) = \\sin(\\pi/4)/\\cos(\\pi/4) = 1\\)."},
		{"What is \\(\\tan(-\\pi/4)\\)? (enter a number)", "-1", "Tangent is odd, so \\(\\tan(-\\pi/4) = -\\tan(\\pi/4) = -1\\)."},
		{"What is the limit of \\(\\tan(x)/x\\) as \\(x \\to 0\\)? (enter a number)", "1", "Near the origin the tangent curve matches \\(y = x\\), so the ratio tends to 1."},
		{"What is the limit of \\(\\tan(x)\\) as \\(x \\to \\pi/2\\) from the left? (type +infinity or -infinity)", "+infinity", "From the left the cosine is positive and tends to zero, so tangent grows without bound."},
	}
	return funcPool(ctx, easy, hard)
}

// cotangent as a real function
type cotangentFuncGen struct{}

func (g *cotangentFuncGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	zeros := []funcEntry{
		{"Name a point in \\([0, \\pi]\\) where \\(\\cot(x)\\) is undefined (type 0 or pi)", "0", "Sine vanishes at 0, so cotangent has a vertical asymptote there."},
		{"Name a point in \\([0, \\pi]\\) where \\(\\cot(x)\\) is undefined (type 0 or pi)", "pi", "Sine vanishes at \\(\\pi\\), so cotangent has a vertical asymptote there."},
	}
	easy := []funcEntry{
		{"What is the period of \\(y = \\cot(x)\\)? (type like pi)", "pi", "Cotangent repeats every \\(\\pi\\) radians, with one branch per period."},
		{"Is \\(\\cot(x)\\) odd, even, or neither? (type one word)", "odd", "\\(\\cot(-x) = -\\cot(x)\\), so cotangent is odd."},
		{"What is the derivative of \\(\\cot(x)\\)?", "-csc^2(x)", "\\(d/dx \\cot(x) = -\\csc^{2}(x)\\)."},
		{"What is the indefinite integral of \\(\\cot(x)\\)?", "ln|sin(x)|", "The integral is \\(\\ln|\\sin(x)| + c\\)."},
		{"What is \\(\\cot(\\pi/2)\\)? (enter a number)", "0", "\\(\\cot(\\pi/2) = \\cos(\\pi/2)/\\sin(\\pi/2) = 0/1 = 0\\)."},
		{"What is the range of \\(\\cot(x)\\)? (type: all reals)", "all reals", "Each branch of cotangent covers every real value."},
	}
	easy = append(easy, zeros...)
	hard := []funcEntry{
		{"What is \\(\\cot(\\pi/4)\\)? (enter a number)", "1", "\\(\\cot(\\pi/4) = \\cos(\\pi/4)/\\sin(\\pi/4) = 1\\)."},
		{"What is the limit of \\(x\\cot(x)\\) as \\(x \\to 0\\)? (enter a number)", "1", "Since \\(\\cot(x) = \\cos(x)/\\sin(x)\\) and \\(\\sin(x)/x \\to 1\\), the product tends to 1."},
		{"Is \\(\\cot(x)\\) increasing or decreasing on each branch? (type one word)", "decreasing", "Cotangent is strictly decreasing on every interval \\((k\\pi, \\pi + k\\pi)\\)."},
	}
	return funcPool(ctx, easy, hard)
}

// secant as a real function
type secantFuncGen struct{}

func (g *secantFuncGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	recips := []funcEntry{
		{"If \\(\\cos(\\theta) = 1/2\\), what is \\(\\sec(\\theta)\\)? (enter a number)", "2", "Secant is the reciprocal of cosine: \\(1/(1/2) = 2\\)."},
		{"If \\(\\cos(\\theta) = 1/3\\), what is \\(\\sec(\\theta)\\)? (enter a number)", "3", "Secant is the reciprocal of cosine: \\(1/(1/3) = 3\\)."},
		{"If \\(\\cos(\\theta) = 1/4\\), what is \\(\\sec(\\theta)\\)? (enter a number)", "4", "Secant is the reciprocal of cosine: \\(1/(1/4) = 4\\)."},
		{"If \\(\\cos(\\theta) = 1/5\\), what is \\(\\sec(\\theta)\\)? (enter a number)", "5", "Secant is the reciprocal of cosine: \\(1/(1/5) = 5\\)."},
	}
	easy := []funcEntry{
		{"What is \\(\\sec(0)\\)? (enter a number)", "1", "\\(\\sec(0) = 1/\\cos(0) = 1/1 = 1\\): the least positive value."},
		{"What is the period of \\(y = \\sec(x)\\)? (type like 2pi)", "2pi", "Secant inherits the \\(2\\pi\\) period of cosine."},
		{"Is \\(\\sec(x)\\) odd, even, or neither? (type one word)", "even", "\\(\\sec(-x) = \\sec(x)\\), so secant is even."},
		{"What is the derivative of \\(\\sec(x)\\)?", "sec(x)tan(x)", "\\(d/dx \\sec(x) = \\sec(x)\\tan(x)\\)."},
		{"At which \\(x\\) in \\([0, \\pi]\\) is \\(\\sec(x)\\) undefined? (type like pi/2)", "pi/2", "Cosine vanishes at \\(\\pi/2\\), so secant has a vertical asymptote there."},
	}
	easy = append(easy, recips...)
	hard := []funcEntry{
		{"What is \\(\\sec(\\pi/3)\\)? (enter a number)", "2", "\\(\\sec(\\pi/3) = 1/\\cos(\\pi/3) = 1/(1/2) = 2\\)."},
		{"What is \\(\\sec^{2}(x) - \\tan^{2}(x)\\)? (enter a number)", "1", "Dividing \\(\\sin^{2}+\\cos^{2}=1\\) by \\(\\cos^{2}\\) gives \\(\\sec^{2} - \\tan^{2} = 1\\)."},
		{"What is the indefinite integral of \\(\\sec(x)\\)?", "ln|sec(x)+tan(x)|", "The integral is \\(\\ln|\\sec(x)+\\tan(x)| + c\\)."},
	}
	return funcPool(ctx, easy, hard)
}

// cosecant as a real function
type cosecantFuncGen struct{}

func (g *cosecantFuncGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	recips := []funcEntry{
		{"If \\(\\sin(\\theta) = 1/2\\), what is \\(\\csc(\\theta)\\)? (enter a number)", "2", "Cosecant is the reciprocal of sine: \\(1/(1/2) = 2\\)."},
		{"If \\(\\sin(\\theta) = 1/3\\), what is \\(\\csc(\\theta)\\)? (enter a number)", "3", "Cosecant is the reciprocal of sine: \\(1/(1/3) = 3\\)."},
		{"If \\(\\sin(\\theta) = 1/4\\), what is \\(\\csc(\\theta)\\)? (enter a number)", "4", "Cosecant is the reciprocal of sine: \\(1/(1/4) = 4\\)."},
	}
	easy := []funcEntry{
		{"What is \\(\\csc(\\pi/2)\\)? (enter a number)", "1", "\\(\\csc(\\pi/2) = 1/\\sin(\\pi/2) = 1/1 = 1\\): the least positive value."},
		{"What is the period of \\(y = \\csc(x)\\)? (type like 2pi)", "2pi", "Cosecant inherits the \\(2\\pi\\) period of sine."},
		{"Is \\(\\csc(x)\\) odd, even, or neither? (type one word)", "odd", "\\(\\csc(-x) = -\\csc(x)\\), so cosecant is odd."},
		{"What is the derivative of \\(\\csc(x)\\)?", "-csc(x)cot(x)", "\\(d/dx \\csc(x) = -\\csc(x)\\cot(x)\\)."},
		{"At which \\(x\\) in \\([0, \\pi]\\) is \\(\\csc(x)\\) undefined? (type 0 or pi)", "0", "Sine vanishes at 0, so cosecant has a vertical asymptote there."},
	}
	easy = append(easy, recips...)
	hard := []funcEntry{
		{"What is \\(\\csc(\\pi/6)\\)? (enter a number)", "2", "\\(\\csc(\\pi/6) = 1/\\sin(\\pi/6) = 1/(1/2) = 2\\)."},
		{"What is \\(\\csc^{2}(x) - \\cot^{2}(x)\\)? (enter a number)", "1", "Dividing \\(\\sin^{2}+\\cos^{2}=1\\) by \\(\\sin^{2}\\) gives \\(\\csc^{2} - \\cot^{2} = 1\\)."},
		{"What is the indefinite integral of \\(\\csc(x)\\)?", "-ln|csc(x)+cot(x)|", "The integral is \\(-\\ln|\\csc(x)+\\cot(x)| + c\\)."},
	}
	return funcPool(ctx, easy, hard)
}

// arcsine as a real function
type arcsineFuncGen struct{}

func (g *arcsineFuncGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return arcsineProblem()
}

// arccosine as a real function
type arccosineFuncGen struct{}

func (g *arccosineFuncGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return arccosineProblem()
}

// arctangent as a real function
type arctangentFuncGen struct{}

func (g *arctangentFuncGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return arctangentProblem()
}

// arccotangent as a real function
type arccotangentFuncGen struct{}

func (g *arccotangentFuncGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return arccotangentProblem()
}
