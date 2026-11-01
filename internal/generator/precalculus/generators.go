package precalculus

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/grader"
	"github.com/chuma-beep/mathua/internal/mathutil"
)

func Register(reg *generator.Registry) {
	for _, e := range all {
		reg.Register(e.id, e.gen)
	}
	registerBatch(reg)
}

var all = []struct {
	id  string
	gen generator.Generator
}{
	{"precalc.vector.basics", &vectorBasicsGen{}},
	{"precalc.poly.remainder_theorem", &remainderTheoremGen{}},
	{"precalc.poly.zeros_multiplicity", &zerosMultiplicityGen{}},
	{"precalc.poly.descartes", &descartesGen{}},
	{"precalc.poly.complex_conjugate", &complexConjugateGen{}},
	{"precalc.rational.oblique_asymptote", &obliqueAsymptoteGen{}},
	{"precalc.log.equation", &logEquationGen{}},
	{"precalc.exp.compound_interest", &compoundInterestGen{}},
	{"precalc.func.composition_value", &compositionValueGen{}},
	{"precalc.analytic.distance", &distanceGen{}},
	{"precalc.poly.rational_roots", &rationalRootsGen{}},
	{"precalc.rational.asymptote", &rationalAsymptoteGen{}},
	{"precalc.exp.equation", &expEquationGen{}},
	{"precalc.log.evaluate", &logEvaluateGen{}},
	{"precalc.exp.growth_model", &growthModelGen{}},
	{"precalc.func.inverse_value", &inverseValueGen{}},
	{"precalc.conic.circle_radius", &circleRadiusGen{}},
	{"precalc.vector.ops", &vectorOpsGen{}},
	{"precalc.vector.dot", &vectorDotGen{}},
	{"precalc.polar.coordinates", &polarCoordsGen{}},
	{"precalc.func.transformations", &funcTransformationsGen{}},
	{"precalc.polar.graph", &polarGraphGen{}},
	{"precalc.parametric.graph", &parametricGraphGen{}},
	{"precalc.func.piecewise", &piecewiseGen{}},
	{"precalc.func.floor_ceiling", &floorCeilingGen{}},
	{"precalc.func.heaviside", &heavisideGen{}},
	{"precalc.func.arith_combine", &arithCombineGen{}},
	{"precalc.parametric.equations", &parametricEqGen{}},
	{"precalc.binomial_theorem", &binomialTheoremGen{}},
	{"precalc.induction", &inductionGen{}},
	{"precalc.systems.nonlinear", &nonlinearSystemsGen{}},
	{"precalc.partial_fractions", &partialFractionsGen{}},
	{"precalc.conic.parabola", &conicParabolaGen{}},
	{"precalc.conic.ellipse", &conicEllipseGen{}},
	{"precalc.conic.hyperbola", &conicHyperbolaGen{}},
	{"precalc.seq.arithmetic", &arithSeqGen{}},
	{"precalc.seq.geometric", &geoSeqGen{}},
	{"precalc.seq.summation", &summationGen{}},
	{"precalc.series.arithmetic", &arithSeriesGen{}},
	{"precalc.series.geometric", &geoSeriesGen{}},
	{"precalc.series.geom_infinite", &geomInfiniteGen{}},
}

func canonPoly(coeffs map[int]int) string {
	if len(coeffs) == 0 {
		return "0"
	}
	exps := make([]int, 0, len(coeffs))
	for e := range coeffs {
		if coeffs[e] != 0 {
			exps = append(exps, e)
		}
	}
	if len(exps) == 0 {
		return "0"
	}
	sort.Sort(sort.Reverse(sort.IntSlice(exps)))
	var b strings.Builder
	for i, e := range exps {
		c := coeffs[e]
		neg := c < 0
		if neg {
			c = -c
		}
		if i == 0 {
			if neg {
				b.WriteString("-")
			}
		} else if neg {
			b.WriteString(" - ")
		} else {
			b.WriteString(" + ")
		}
		switch {
		case e == 0:
			b.WriteString(strconv.Itoa(c))
		case c != 1:
			b.WriteString(strconv.Itoa(c))
			b.WriteString("x")
		default:
			b.WriteString("x")
		}
		if e > 1 {
			b.WriteString("^")
			b.WriteString(strconv.Itoa(e))
		}
	}
	return strings.TrimSpace(b.String())
}

func gradeNum(expected, user string, tol float64) grader.Result {
	e, err1 := strconv.ParseFloat(strings.TrimSpace(expected), 64)
	uStr := strings.TrimSpace(user)
	uStr = strings.TrimSuffix(uStr, ".0")
	u, err2 := strconv.ParseFloat(uStr, 64)
	if err1 != nil || err2 != nil {
		return grader.Result{Correct: false, Score: 0, Feedback: "Answer must be a number"}
	}
	if math.Abs(e-u) <= tol*math.Max(1, math.Abs(e)) {
		return grader.Result{Correct: true, Score: 1}
	}
	return grader.Result{Correct: false, Score: 0, Feedback: "Incorrect value"}
}

func parseVec(s string) (float64, float64, bool) {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "(")
	s = strings.TrimSuffix(s, ")")
	parts := strings.Split(s, ",")
	if len(parts) != 2 {
		return 0, 0, false
	}
	x, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	y, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	return x, y, err1 == nil && err2 == nil
}

func gradeVec(expected, user string) grader.Result {
	ex, ey, ok1 := parseVec(expected)
	ux, uy, ok2 := parseVec(user)
	if !ok1 || !ok2 {
		return grader.Result{Correct: false, Score: 0, Feedback: "Answer must be of the form (x, y)"}
	}
	const tol = 1e-6
	if math.Abs(ex-ux) <= tol && math.Abs(ey-uy) <= tol {
		return grader.Result{Correct: true, Score: 1}
	}
	return grader.Result{Correct: false, Score: 0, Feedback: "Incorrect components"}
}

type vec struct{ x, y int }

func (v vec) String() string { return fmt.Sprintf("(%d, %d)", v.x, v.y) }

var triples = [][2]int{{3, 4}, {6, 8}, {5, 12}, {8, 15}, {7, 24}, {9, 12}, {20, 21}}

type vectorBasicsGen struct{}

func (g *vectorBasicsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	t := triples[rand.Intn(len(triples))]
	k := 1 + rand.Intn(1+int(ctx.Difficulty*5))
	v := vec{t[0] * k, t[1] * k}
	mag := k * isqrt(t[0]*t[0]+t[1]*t[1])
	ans := strconv.Itoa(mag)
	return generator.Problem{
		Question:    fmt.Sprintf("Find the magnitude of the vector \\(\\langle %d, %d \\rangle\\).", v.x, v.y),
		Answer:      ans,
		Explanation: fmt.Sprintf("\\(\\|\\mathbf{v}\\| = \\sqrt{%d^2 + %d^2} = %d\\).", v.x, v.y, mag),
	}
}

func isqrt(n int) int {
	r := int(math.Sqrt(float64(n)))
	for r*r > n {
		r--
	}
	for (r+1)*(r+1) <= n {
		r++
	}
	return r
}

type vectorOpsGen struct{}

// A vector answer is a pair in parentheses. Declared numeric, the concept put a
// keypad in front of it that cannot type a bracket or a minus sign in the shape
// the grader expects.
func (g *vectorOpsGen) AnswerFormat() generator.AnswerFormat {
	return generator.AnswerFormat{Hint: "Vector as (x, y), e.g. (3, -2)", InputMode: "text"}
}

func (g *vectorOpsGen) Grade(expected, userAnswer string) grader.Result {
	return gradeVec(expected, userAnswer)
}

func (g *vectorOpsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := 1 + rand.Intn(1+int(ctx.Difficulty*8))
	u := vec{(rand.Intn(2*scale+1) - scale), (rand.Intn(2*scale+1) - scale)}
	w := vec{(rand.Intn(2*scale+1) - scale), (rand.Intn(2*scale+1) - scale)}
	var res vec
	var q, ex string
	switch rand.Intn(3) {
	case 0:
		res = vec{u.x + w.x, u.y + w.y}
		q = fmt.Sprintf("\\(\\mathbf{u} = \\langle %d, %d \\rangle\\), \\(\\mathbf{w} = \\langle %d, %d \\rangle\\). Find \\(\\mathbf{u} + \\mathbf{w}\\).", u.x, u.y, w.x, w.y)
		ex = "\\(\\langle u_x+w_x,\\ u_y+w_y \\rangle\\)."
	case 1:
		res = vec{u.x - w.x, u.y - w.y}
		q = fmt.Sprintf("\\(\\mathbf{u} = \\langle %d, %d \\rangle\\), \\(\\mathbf{w} = \\langle %d, %d \\rangle\\). Find \\(\\mathbf{u} - \\mathbf{w}\\).", u.x, u.y, w.x, w.y)
		ex = "\\(\\langle u_x-w_x,\\ u_y-w_y \\rangle\\)."
	default:
		k := 1 + rand.Intn(5)
		res = vec{k * u.x, k * u.y}
		q = fmt.Sprintf("Given \\(\\mathbf{u} = \\langle %d, %d \\rangle\\), find \\(%d\\,\\mathbf{u}\\).", u.x, u.y, k)
		ex = fmt.Sprintf("\\(%d\\langle u_x, u_y \\rangle = \\langle %d\\,u_x, %d\\,u_y \\rangle\\).", k, k, k)
	}
	ans := res.String()
	return generator.Problem{
		Question:    q,
		Answer:      ans,
		Explanation: "Componentwise arithmetic: " + ex,
	}
}

type vectorDotGen struct{}

func (g *vectorDotGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := 1 + rand.Intn(1+int(ctx.Difficulty*6))
	a := vec{rand.Intn(2*scale+1) - scale, rand.Intn(2*scale+1) - scale}
	b := vec{rand.Intn(2*scale+1) - scale, rand.Intn(2*scale+1) - scale}
	dot := a.x*b.x + a.y*b.y
	ans := strconv.Itoa(dot)
	return generator.Problem{
		Question:    fmt.Sprintf("Compute the dot product \\(\\langle %d, %d \\rangle \\cdot \\langle %d, %d \\rangle\\).", a.x, a.y, b.x, b.y),
		Answer:      ans,
		Explanation: fmt.Sprintf("\\(%d\\cdot%d + %d\\cdot%d = %d\\).", a.x, b.x, a.y, b.y, dot),
	}
}

type polarCoordsGen struct{}

func (g *polarCoordsGen) Grade(expected, userAnswer string) grader.Result {
	return gradeNum(expected, userAnswer, 0.02)
}

func (g *polarCoordsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type cand struct {
		q, ans, ex string
	}
	var cands []cand
	axis := rand.Intn(4)
	r := 1 + rand.Intn(3+int(ctx.Difficulty*6))
	switch axis {
	case 0:
		cands = append(cands,
			cand{fmt.Sprintf("Convert the point \\((%d, 0)\\) to polar coordinates \\((r,\\theta)\\) with \\(r>0\\) and \\(0^{\\circ} \\le \\theta < 360^{\\circ}\\). Give \\(r\\).", r),
				strconv.Itoa(r), fmt.Sprintf("\\(r=\\sqrt{x^2+y^2}=%d\\).", r)},
			cand{fmt.Sprintf("Convert the point \\((%d, 0)\\) to polar coordinates. Give \\(\\theta\\) in degrees.", r),
				"0", "\\(\\theta=0^{\\circ}\\) on the positive x-axis."})
	case 1:
		cands = append(cands,
			cand{fmt.Sprintf("Convert the point \\((0, %d)\\) to polar coordinates. Give \\(r\\).", r),
				strconv.Itoa(r), fmt.Sprintf("\\(r=\\sqrt{0+%d^2}=%d\\).", r, r)},
			cand{fmt.Sprintf("Convert the point \\((0, %d)\\) to polar coordinates. Give \\(\\theta\\) in degrees.", r),
				"90", "\\(\\theta=90^{\\circ}\\) on the positive y-axis."})
	case 2:
		cands = append(cands,
			cand{fmt.Sprintf("Convert the point \\((- %d, 0)\\) to polar coordinates. Give \\(\\theta\\) in degrees.", r),
				"180", "\\((-r,0)\\) lies on the negative x-axis."})
	default:
		cands = append(cands,
			cand{fmt.Sprintf("Convert the point \\((0, -%d)\\) to polar coordinates. Give \\(\\theta\\) in degrees.", r),
				"270", "\\((0,-r)\\) lies on the negative y-axis."})
	}
	if ctx.Difficulty > 0.35 {
		deg := []int{45, 135, 225, 315}[rand.Intn(4)]
		mag := float64(1+rand.Intn(2+int(ctx.Difficulty*4))) * math.Sqrt2
		x := mag * math.Cos(float64(deg)*math.Pi/180)
		y := mag * math.Sin(float64(deg)*math.Pi/180)
		askR := rand.Intn(2) == 0
		if askR {
			cands = append(cands, cand{
				fmt.Sprintf("A point has Cartesian coordinates \\((%.4f, %.4f)\\). Find its polar distance \\(r\\), rounded to 2 decimals.", round2(x), round2(y)),
				fmt.Sprintf("%.2f", mag),
				"\\(r=\\sqrt{x^2+y^2}\\), which recovers the original radius.",
			})
		} else {
			cands = append(cands, cand{
				fmt.Sprintf("A point has Cartesian coordinates \\((%.4f, %.4f)\\). Give its polar angle \\(\\theta\\) in degrees, \\(0^{\\circ} \\le \\theta < 360^{\\circ}\\).", round2(x), round2(y)),
				strconv.Itoa(deg),
				fmt.Sprintf("The point lies on the ray at \\(%d^{\\circ}\\) from the polar axis.", deg),
			})
		}
	}
	c := cands[rand.Intn(len(cands))]
	return generator.Problem{Question: c.q, Answer: c.ans, Explanation: c.ex}
}

func round2(f float64) float64 { return math.Round(f*100) / 100 }

func mcProblem(q string, correct string, distractors []string) generator.Problem {
	opts := []string{correct}
	for _, d := range distractors {
		if d != correct {
			opts = append(opts, d)
		}
	}
	rand.Shuffle(len(opts), func(i, j int) { opts[i], opts[j] = opts[j], opts[i] })
	letters := []string{"A", "B", "C", "D", "E"}
	var qb strings.Builder
	qb.WriteString(q)
	for i, o := range opts {
		fmt.Fprintf(&qb, "\n\n%s) %s", letters[i], o)
	}
	return generator.Problem{
		Question:    qb.String(),
		Answer:      correct,
		Explanation: fmt.Sprintf("The correct choice is: %s.", correct),
	}
}

type funcTransformationsGen struct{}

func (g *funcTransformationsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return transformationPointProblem(ctx)
}

type polarGraphGen struct{}

func (g *polarGraphGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return polarCircleProblem()
}

type parametricGraphGen struct{}

func (g *parametricGraphGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return parametricParabolaProblem()
}

type piecewiseGen struct{}

func (g *piecewiseGen) Grade(expected, userAnswer string) grader.Result {
	return gradeNum(expected, userAnswer, 1e-6)
}

func (g *piecewiseGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	b1 := -(1 + rand.Intn(3))
	b2 := 1 + rand.Intn(3)
	a1 := 1 + rand.Intn(1+int(ctx.Difficulty*3))
	a2 := 1 + rand.Intn(1+int(ctx.Difficulty*3))
	c1 := rand.Intn(2*3+1) - 3
	a3 := 1 + rand.Intn(1+int(ctx.Difficulty*3))
	c3 := rand.Intn(2*3+1) - 3

	evalAt := func(x int) int {
		switch {
		case x < b1:
			return a1*x + c1
		case x > b2:
			return a3*x + c3
		default:
			return a2*x*x - 1
		}
	}
	region := rand.Intn(3)
	var x int
	switch region {
	case 0:
		x = b1 - 1 - rand.Intn(3)
	case 1:
		x = b2 + 1 + rand.Intn(3)
	default:
		x = b1 + rand.Intn(b2-b1+1)
	}
	val := evalAt(x)
	q := fmt.Sprintf("Let\n\\[f(x)=\\begin{cases}%dx %+d & x<%d \\\\ %dx^2-1 & %d \\le x \\le %d \\\\ %dx %+d & x>%d\\end{cases}\\]\nFind \\(f(%d)\\).",
		a1, c1, b1, a2, b1, b2, a3, c3, b2, x)
	ans := strconv.Itoa(val)
	ex := ""
	switch {
	case x < b1:
		ex = fmt.Sprintf("\\(%d<%d\\), so use the first branch: \\(f(%d)=%d\\cdot%d %+d = %d\\).", x, b1, x, a1, x, c1, val)
	case x > b2:
		ex = fmt.Sprintf("\\(%d>%d\\), so use the third branch: \\(f(%d)=%d\\cdot%d %+d = %d\\).", x, b2, x, a3, x, c3, val)
	default:
		ex = fmt.Sprintf("\\(%d \\le %d \\le %d\\), so use the middle branch: \\(f(%d)=%d\\cdot%d^2-1 = %d\\).", b1, x, b2, x, a2, x, val)
	}
	return generator.Problem{Question: q, Answer: ans, Explanation: ex}
}

type arithCombineGen struct{}

func (g *arithCombineGen) Grade(expected, userAnswer string) grader.Result {
	return gradeExpr(expected, userAnswer)
}

func (g *arithCombineGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := 1 + rand.Intn(1+int(ctx.Difficulty*4))
	fc := map[int]int{2: 1 + rand.Intn(scale), 1: rand.Intn(2*scale+1) - scale, 0: rand.Intn(2*scale+1) - scale}
	gc := map[int]int{2: 1 + rand.Intn(scale), 1: rand.Intn(2*scale+1) - scale, 0: rand.Intn(2*scale+1) - scale}
	op := rand.Intn(3)
	out := map[int]int{}
	names := [3]string{"+", "-", ""}
	switch op {
	case 0:
		for e := range fc {
			out[e] = fc[e] + gc[e]
		}
	case 1:
		for e := range fc {
			out[e] = fc[e] - gc[e]
		}
	default:
		for e1, c1 := range fc {
			for e2, c2 := range gc {
				out[e1+e2] += c1 * c2
			}
		}
	}
	fS := polyWithX(canonPoly(fc))
	gS := polyWithX(canonPoly(gc))
	name := names[op]
	if op == 2 {
		name = "\\cdot"
	}
	q := fmt.Sprintf("Let \\(f(x) = %s\\) and \\(g(x) = %s\\). Express \\((f%s g)(x)\\) in expanded form.",
		fS, gS, name)
	ans := canonPoly(out)
	ex := fmt.Sprintf("Combine like powers of \\(x\\): \\((f%s g)(x) = %s\\).", name, ans)
	return generator.Problem{Question: q, Answer: ans, Explanation: ex}
}

func polyWithX(p string) string {
	if p == "0" {
		return "0"
	}
	return p
}

type parametricEqGen struct{}

func (g *parametricEqGen) Grade(expected, userAnswer string) grader.Result {
	return gradeExpr(expected, userAnswer)
}

func (g *parametricEqGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	h := rand.Intn(2*3+1) - 3
	k := rand.Intn(2*3+1) - 3
	a := 1 + rand.Intn(1+int(ctx.Difficulty*2))
	q := fmt.Sprintf("Eliminate the parameter: \\(x = t %+d,\\ y = %dt^2 %+d\\). Write \\(y\\) as a function of \\(x\\) in expanded form.", h, a, k)
	t := map[string]string{"t": fmt.Sprintf("(x %+d)", -h)}
	exInner := fmt.Sprintf("%s^2 %+d", t["t"], k)
	ans := canonPoly(expandShiftedSquare(a, -h, k))
	ex := fmt.Sprintf("Since \\(t = x %+d\\), substitute: \\(y = %d\\left(%s\\right)^2 %+d = %s\\).", -h, a, t["t"], k, ans)
	_ = exInner
	return generator.Problem{Question: q, Answer: ans, Explanation: ex}
}

func expandShiftedSquare(a, s, k int) map[int]int {
	return map[int]int{2: a, 1: 2 * a * s, 0: a*s*s + k}
}

type binomialTheoremGen struct{}

func binom(n, k int) int {
	if k < 0 || k > n {
		return 0
	}
	if k > n-k {
		k = n - k
	}
	res := 1
	for i := 0; i < k; i++ {
		res = res * (n - i) / (i + 1)
	}
	return res
}

func (g *binomialTheoremGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	n := 3 + rand.Intn(2+int(ctx.Difficulty*5))
	k := rand.Intn(n + 1)
	c := 0
	for c == 0 {
		c = rand.Intn(2*(1+int(ctx.Difficulty*3))+1) - (1 + int(ctx.Difficulty*3))
	}
	coef := binom(n, k)
	for i := 0; i < n-k; i++ {
		coef *= c
	}
	q := fmt.Sprintf("Find the coefficient of \\(x^{%d}\\) in the expansion of \\((x %+d)^{%d}\\).", k, c, n)
	ans := strconv.Itoa(coef)
	ex := fmt.Sprintf("The term is \\(\\binom{%d}{%d} (%+d)^{%d} x^{%d}\\), giving coefficient %s.", n, k, c, n-k, k, ans)
	return generator.Problem{Question: q, Answer: ans, Explanation: ex}
}

type inductionGen struct{}

func (g *inductionGen) Grade(expected, userAnswer string) grader.Result {
	return gradeNum(expected, userAnswer, 1e-6)
}

func (g *inductionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	n := 3 + rand.Intn(1+int(ctx.Difficulty*4))
	type claim struct{ q, formula string }
	claims := []claim{
		{"\\(P(n):\\ 1+3+5+\\cdots+(2n-1)=n^2\\)", fmt.Sprintf("sum of first %d odd numbers = %d", n, n*n)},
		{"\\(P(n):\\ 1+2+\\cdots+n=\\frac{n(n+1)}{2}\\)", fmt.Sprintf("sum of first %d integers = %d", n, n*(n+1)/2)},
		{"\\(P(n):\\ 3+6+9+\\cdots+3n=\\frac{3n(n+1)}{2}\\)", fmt.Sprintf("sum of first %d multiples of 3 = %d", n, 3*n*(n+1)/2)},
	}
	cl := claims[rand.Intn(len(claims))]
	value := strings.Split(cl.formula, "= ")[1]
	q := cl.q + fmt.Sprintf(" To start the induction, verify \\(P(%d)\\): what common value do both sides equal?", n)
	ex := fmt.Sprintf("LHS at \\(n=%d\\): %s. RHS: substituting \\(n=%d\\) gives the same value %s.", n, value, n, value)
	return generator.Problem{Question: q, Answer: value, Explanation: ex}
}

type nonlinearSystemsGen struct{}

func (g *nonlinearSystemsGen) Grade(expected, userAnswer string) grader.Result {
	return gradeVec(expected, userAnswer)
}

func (g *nonlinearSystemsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := 1 + rand.Intn(1+int(ctx.Difficulty*3))
	p := rand.Intn(2*scale+1) - scale
	q := rand.Intn(2*scale+1) - scale
	// q == p would make one root a double root and the "smaller" one ambiguous;
	// p + q == 0 would make the line's coefficient zero, and a line written
	// "y = 0x - 1" is not the situation the question describes.
	for q == p || p+q == 0 {
		q = rand.Intn(2*scale+1) - scale
	}
	m := p + q
	b := q * p
	c := rand.Intn(2*scale+1) - scale
	bTotal := b + c
	ySmall := m*p + bTotal
	smaller := p
	if q < smaller {
		smaller = q
		ySmall = m*q + bTotal
	}
	ans := fmt.Sprintf("(%d, %d)", smaller, ySmall)
	// Both equations used to print their constant with %+d and the line its
	// coefficient as %dx, so a draw could read "y = x^2 +0" beside "y = -1x +0"
	// and the system arrived at the learner as noise. Zero terms are omitted and
	// a unit coefficient prints bare.
	parab := "y = x^{2}"
	if c > 0 {
		parab += fmt.Sprintf(" + %d", c)
	} else if c < 0 {
		parab += fmt.Sprintf(" - %d", -c)
	}
	line := "y = "
	switch m {
	case 1:
		line += "x"
	case -1:
		line += "-x"
	default:
		line += fmt.Sprintf("%dx", m)
	}
	if bTotal > 0 {
		line += fmt.Sprintf(" + %d", bTotal)
	} else if bTotal < 0 {
		line += fmt.Sprintf(" - %d", -bTotal)
	}
	qText := fmt.Sprintf("Solve the system \\[\\begin{cases}%s\\\\%s\\end{cases}\\]\nEnter the solution with the smaller \\(x\\)-coordinate in the form \\((x, y)\\).", parab, line)
	ex := fmt.Sprintf("Substitute the line into the parabola and collect the terms: the \\(x\\) values are %d and %d. At the smaller one, \\(x=%d\\), the line gives \\(y=%d\\), which is the point asked for.", p, q, smaller, ySmall)
	return generator.Problem{Question: qText, Answer: ans, Explanation: ex}
}

type partialFractionsGen struct{}

func gradeExpr(expected, userAnswer string) grader.Result {
	router := grader.NewRouter()
	res := router.Grade(grader.GradingExpression, expected, userAnswer)
	if res.Correct || !strings.Contains(strings.ToLower(res.Feedback), "error") {
		return res
	}
	norm := func(s string) string {
		s = strings.ReplaceAll(s, " ", "")
		return strings.ToLower(s)
	}
	if norm(expected) == norm(userAnswer) {
		return grader.Result{Correct: true, Score: 1}
	}
	return grader.Result{Correct: false, Score: 0, Feedback: "Incorrect expression"}
}

func (g *partialFractionsGen) Grade(expected, userAnswer string) grader.Result {
	router := grader.NewRouter()
	res := router.Grade(grader.GradingExpression, expected, userAnswer)
	if res.Correct || !strings.Contains(strings.ToLower(res.Feedback), "error") {
		return res
	}
	norm := func(s string) string {
		s = strings.ReplaceAll(s, " ", "")
		return strings.ToLower(s)
	}
	if norm(expected) == norm(userAnswer) {
		return grader.Result{Correct: true, Score: 1}
	}
	return grader.Result{Correct: false, Score: 0, Feedback: "Incorrect decomposition"}
}

func (g *partialFractionsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := 1 + rand.Intn(1+int(ctx.Difficulty*4))
	r1 := rand.Intn(2*scale+1) - scale
	r2 := rand.Intn(2*scale+1) - scale
	for r2 == r1 {
		r2 = rand.Intn(2*scale+1) - scale
	}
	A := 0
	for A == 0 {
		A = rand.Intn(2*scale+1) - scale
	}
	B := 0
	for B == 0 || B == A {
		B = rand.Intn(2*scale+1) - scale
	}
	num := map[int]int{}
	num[1] = A + B
	num[0] = -(A*r2 + B*r1)

	factor := func(root int) string {
		if root < 0 {
			return fmt.Sprintf("(x + %d)", -root)
		}
		return fmt.Sprintf("(x - %d)", root)
	}
	piece := func(coef, root int) string {
		cs := strconv.Itoa(coef)
		return cs + "/" + factor(root)
	}
	ans := piece(A, r1)
	t2 := piece(B, r2)
	if strings.HasPrefix(t2, "-") {
		ans += " - " + strings.TrimPrefix(t2, "-")
	} else {
		ans += " + " + t2
	}
	q := fmt.Sprintf("Decompose into partial fractions: \\(\\frac{%s}{%s%s}\\). Write the answer in the form \\(%s\\) with constants \\(A\\), \\(B\\).",
		canonPoly(num), factor(r1), factor(r2), "\\frac{A}{x-r_1}+\\frac{B}{x-r_2}")
	ex := fmt.Sprintf("Cover-up method: setting \\(x=%d\\) kills the \\(%s\\) factor and yields \\(A=%d\\); setting \\(x=%d\\) yields \\(B=%d\\). So the decomposition is \\(%s\\).",
		r1, factor(r2), A, r2, B, ans)
	return generator.Problem{Question: q, Answer: ans, Explanation: ex}
}

var conicPairs = [][3]int{{3, 4, 5}, {6, 8, 10}, {5, 12, 13}, {9, 12, 15}, {8, 15, 17}, {12, 16, 20}, {7, 24, 25}}

type conicParabolaGen struct{}

func (g *conicParabolaGen) Grade(expected, userAnswer string) grader.Result {
	return gradeNum(expected, userAnswer, 1e-6)
}

func (g *conicParabolaGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	p := 1 + rand.Intn(2+int(ctx.Difficulty*4))
	neg := rand.Intn(2) == 0
	if neg {
		p = -p
	}
	fourP := 4 * p
	if rand.Intn(2) == 0 {
		ans := strconv.Itoa(p)
		q := fmt.Sprintf("The parabola \\(y^2 = %sx\\) has its focus on the x-axis at \\((p, 0)\\). Find \\(p\\).", coefStr(fourP))
		ex := fmt.Sprintf("Comparing with \\(y^2 = 4px\\): \\(4p = %+d\\), so \\(p = %d\\).", fourP, p)
		return generator.Problem{Question: q, Answer: ans, Explanation: ex}
	}
	ans := strconv.Itoa(p)
	q := fmt.Sprintf("The parabola \\(x^2 = %sy\\) has its focus on the y-axis at \\((0, p)\\). Find \\(p\\).", coefStr(fourP))
	ex := fmt.Sprintf("Comparing with \\(x^2 = 4py\\): \\(4p = %+d\\), so \\(p = %d\\).", fourP, p)
	return generator.Problem{Question: q, Answer: ans, Explanation: ex}
}

func coefStr(c int) string {
	switch c {
	case 1:
		return ""
	case -1:
		return "-"
	default:
		return strconv.Itoa(c)
	}
}

type conicEllipseGen struct{}

func (g *conicEllipseGen) Grade(expected, userAnswer string) grader.Result {
	return gradeNum(expected, userAnswer, 1e-6)
}

func (g *conicEllipseGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	t := conicPairs[rand.Intn(len(conicPairs))]
	a := t[0]
	b := t[1]
	c := t[2]
	// The report text promises a is the semi-major axis. Swapping the pair
	// produced questions such as a=6, b=8 with the words "with a>b", which
	// report 262 quoted back verbatim. Keep the invariant the text states.
	if a < b {
		a, b = b, a
	}
	if ctx.Difficulty > 0.5 && rand.Intn(2) == 0 {
		xDen, yDen := a*a, b*b
		if rand.Intn(2) == 0 {
			xDen, yDen = yDen, xDen
		}
		q := fmt.Sprintf("An ellipse is given by \\(\\frac{x^2}{%d} + \\frac{y^2}{%d} = 1\\). Find the distance from the center to each focus \\(c\\).", xDen, yDen)
		ex := fmt.Sprintf("\\(c^2 = a^2-b^2 = %d-%d = %d\\), so \\(c = %d\\).", maxInt(a*a, b*b), minInt(a*a, b*b), c*c, c)
		return generator.Problem{Question: q, Answer: strconv.Itoa(c), Explanation: ex}
	}
	var q string
	if rand.Intn(2) == 0 {
		q = fmt.Sprintf("An ellipse has semi-axes \\(a=%d\\) and \\(b=%d\\) (with \\(a>b\\)). Find \\(c\\), the distance from center to focus, where \\(c^2=a^2-b^2\\).", a, b)
	} else {
		q = fmt.Sprintf("An ellipse has semi-major axis \\(a=%d\\) and semi-minor axis \\(b=%d\\). Find \\(c\\), the distance from center to focus, where \\(c^2=a^2-b^2\\).", a, b)
	}
	ex := fmt.Sprintf("\\(c=\\sqrt{%d-%d}=\\sqrt{%d}=%d\\).", a*a, b*b, c*c, c)
	return generator.Problem{Question: q, Answer: strconv.Itoa(c), Explanation: ex}
}

type conicHyperbolaGen struct{}

func (g *conicHyperbolaGen) Grade(expected, userAnswer string) grader.Result {
	return gradeNum(expected, userAnswer, 0.02)
}

func (g *conicHyperbolaGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	t := conicPairs[rand.Intn(len(conicPairs))]
	switch rand.Intn(3) {
	case 0:
		q := fmt.Sprintf("For the hyperbola \\(\\frac{x^2}{%d} - \\frac{y^2}{%d} = 1\\), find \\(c\\), the distance from center to focus \\(c^2=a^2+b^2\\).", t[0]*t[0], t[1]*t[1])
		ex := fmt.Sprintf("\\(c=\\sqrt{%d+%d}=\\sqrt{%d}=%d\\).", t[0]*t[0], t[1]*t[1], t[2]*t[2], t[2])
		return generator.Problem{Question: q, Answer: strconv.Itoa(t[2]), Explanation: ex}
	case 1:
		mult := 1 + rand.Intn(3+int(ctx.Difficulty*3))
		b := t[1] * mult
		a := t[0] * mult
		val := float64(b) / float64(a)
		q := fmt.Sprintf("For \\(\\frac{x^2}{%d} - \\frac{y^2}{%d} = 1\\), find the positive slope of the asymptotes \\(y=\\pm mx\\), rounded to 2 decimals.", a*a, b*b)
		ex := fmt.Sprintf("Asymptotes have \\(m=b/a=%d/%d=%.2f\\).", b, a, val)
		return generator.Problem{Question: q, Answer: fmt.Sprintf("%.2f", val), Explanation: ex}
	default:
		q := fmt.Sprintf("For \\(\\frac{x^2}{%d} - \\frac{y^2}{%d} = 1\\), find \\(a\\), the length of the semi-transverse axis.", t[0]*t[0], t[1]*t[1])
		ex := fmt.Sprintf("The \\(x^2\\) term is positive, so the transverse axis is horizontal with \\(a^2=%d\\): \\(a=%d\\).", t[0]*t[0], t[0])
		return generator.Problem{Question: q, Answer: strconv.Itoa(t[0]), Explanation: ex}
	}
}

type arithSeqGen struct{}

func (g *arithSeqGen) Grade(expected, userAnswer string) grader.Result {
	return gradeNum(expected, userAnswer, 1e-6)
}

func (g *arithSeqGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	tier := mathutil.DifficultyTier(ctx.Difficulty)
	var a1, d, n, an int
	var q, ex string
	for {
		a1Range := 5 + tier*10
		a1 = rand.Intn(2*a1Range+1) - a1Range
		dRange := 4 + tier*6
		d = 0
		for d == 0 {
			d = rand.Intn(2*dRange+1) - dRange
		}
		nMin := 5 + tier*5
		nMax := 12 + tier*15
		n = smallInt(nMin, nMax)
		if n < 5 {
			n = 5
		}
		an = a1 + (n-1)*d
		ansStr := strconv.Itoa(an)
		q = fmt.Sprintf("An arithmetic sequence starts at %d with common difference %+d. Find term \\(a_{%d}\\).", a1, d, n)
		if an == a1 || an == d || an == -d || an == n || containsWord(q, ansStr) {
			continue
		}
		ex = fmt.Sprintf("The nth term of an arithmetic sequence is given by the formula \\(a_n = a_1 + (n - 1)d\\). Substitute first term \\(a_1 = %d\\), common difference \\(d = %+d\\) and term index \\(n = %d\\) into the formula: \\(a_{%d} = %d + (%d - 1)(%+d)\\). Simplify the product \\((%d)(%+d) = %d\\) and add to obtain \\(a_{%d} = %d + (%d) = %d\\).",
			a1, d, n, n, a1, n, d, n-1, d, (n-1)*d, n, a1, (n-1)*d, an)
		break
	}
	return generator.Problem{
		Question:    q,
		Answer:      strconv.Itoa(an),
		Explanation: ex,
		Facts: map[string]string{
			"a1":     strconv.Itoa(a1),
			"d":      strconv.Itoa(d),
			"n":      strconv.Itoa(n),
			"answer": strconv.Itoa(an),
		},
	}
}

type geoSeqGen struct{}

func (g *geoSeqGen) Grade(expected, userAnswer string) grader.Result {
	return gradeNum(expected, userAnswer, 1e-6)
}

func (g *geoSeqGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	tier := mathutil.DifficultyTier(ctx.Difficulty)
	ratioPool := []int{2, 3, -2}
	if tier == 1 {
		ratioPool = []int{2, 3, -2, -3, 4}
	} else if tier == 2 {
		ratioPool = []int{2, 3, -2, -3, 4, -4, 5}
	}
	var a1, r, n, an, powerVal int
	var q, ex string
	for {
		r = ratioPool[rand.Intn(len(ratioPool))]
		maxN := 8
		if abs(r) == 2 {
			maxN = 6 + tier*2
		} else if abs(r) == 3 {
			maxN = 4 + tier*2
		} else {
			maxN = 3 + tier
		}
		n = 3 + rand.Intn(maxN-2)
		a1 = 1 + rand.Intn(2+tier*3)
		if tier > 0 && rand.Intn(2) == 0 {
			a1 = -a1
		}
		powerVal = 1
		for i := 1; i < n; i++ {
			powerVal *= r
		}
		an = a1 * powerVal
		ansStr := strconv.Itoa(an)
		q = fmt.Sprintf("A geometric sequence starts at %d with common ratio %d. Find term \\(a_{%d}\\).", a1, r, n)
		if an == a1 || an == r || an == -r || an == n || containsWord(q, ansStr) {
			continue
		}
		ex = fmt.Sprintf("The nth term of a geometric sequence is governed by \\(a_n = a_1 r^{n-1}\\). Substitute initial term \\(a_1 = %d\\), common ratio \\(r = %d\\) and index \\(n = %d\\): \\(a_{%d} = %d \\times (%d)^{%d - 1} = %d \\times (%d)^{%d}\\). Evaluate the exponent \\((%d)^{%d} = %d\\) and multiply to yield \\(a_{%d} = %d \\times %d = %d\\).",
			a1, r, n, n, a1, r, n, a1, r, n-1, r, n-1, powerVal, n, a1, powerVal, an)
		break
	}
	return generator.Problem{
		Question:    q,
		Answer:      strconv.Itoa(an),
		Explanation: ex,
		Facts: map[string]string{
			"a1":     strconv.Itoa(a1),
			"r":      strconv.Itoa(r),
			"n":      strconv.Itoa(n),
			"answer": strconv.Itoa(an),
		},
	}
}

type summationGen struct{}

func (g *summationGen) Grade(expected, userAnswer string) grader.Result {
	return gradeNum(expected, userAnswer, 1e-6)
}

func (g *summationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	tier := mathutil.DifficultyTier(ctx.Difficulty)
	var n, a, b, total, triVal int
	var q, ex string
	for {
		nMin := 4 + tier*4
		nMax := 8 + tier*8
		n = smallInt(nMin, nMax)
		a = 2 + rand.Intn(2+tier*3)
		b = 0
		for b == 0 {
			bRange := 3 + tier*4
			b = rand.Intn(2*bRange+1) - bRange
		}
		triVal = n * (n + 1) / 2
		total = a*triVal + b*n
		ansStr := strconv.Itoa(total)
		q = fmt.Sprintf("Evaluate \\(\\sum_{i=1}^{%d} (%si %+d)\\).", n, coefStr(a), b)
		if total == n || total == a || total == b || total == -b || containsWord(q, ansStr) {
			continue
		}
		ex = fmt.Sprintf("By the linearity of summation, \\(\\sum_{i=1}^n (ai + b) = a \\frac{n(n+1)}{2} + bn\\). Substitute \\(n = %d\\), coefficient \\(a = %d\\) and constant \\(b = %d\\): \\(\\sum_{i=1}^{%d} (%si %+d) = %d \\times \\frac{%d(%d)}{2} + (%d)(%d)\\). Compute the triangular number \\(\\frac{%d(%d)}{2} = %d\\) and linear term \\((%d)(%d) = %d\\), giving \\(%d(%d) + (%d) = %d + (%d) = %d\\).",
			n, a, b, n, coefStr(a), b, a, n, n+1, b, n, n, n+1, triVal, b, n, b*n, a, triVal, b*n, a*triVal, b*n, total)
		break
	}
	return generator.Problem{
		Question:    q,
		Answer:      strconv.Itoa(total),
		Explanation: ex,
		Facts: map[string]string{
			"n":      strconv.Itoa(n),
			"a":      strconv.Itoa(a),
			"b":      strconv.Itoa(b),
			"answer": strconv.Itoa(total),
		},
	}
}

type arithSeriesGen struct{}

func (g *arithSeriesGen) Grade(expected, userAnswer string) grader.Result {
	return gradeNum(expected, userAnswer, 1e-6)
}

func (g *arithSeriesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	tier := mathutil.DifficultyTier(ctx.Difficulty)
	var n, a1, d, an, sum, bracketVal int
	var q, ex string
	for {
		nMin := 4 + tier*4
		nMax := 8 + tier*10
		n = smallInt(nMin, nMax)
		a1Range := 6 + tier*8
		a1 = rand.Intn(2*a1Range+1) - a1Range
		dRange := 4 + tier*5
		d = 0
		for d == 0 {
			d = rand.Intn(2*dRange+1) - dRange
		}
		an = a1 + (n-1)*d
		sum = n * (a1 + an) / 2
		bracketVal = 2*a1 + (n-1)*d
		ansStr := strconv.Itoa(sum)
		q = fmt.Sprintf("Find the sum of the first %d terms of the arithmetic sequence \\(a_1=%d\\), \\(d=%+d\\).", n, a1, d)
		if sum == n || sum == a1 || sum == d || sum == -d || containsWord(q, ansStr) {
			continue
		}
		ex = fmt.Sprintf("The sum of the first n terms of an arithmetic series is \\(S_n = \\frac{n}{2}[2a_1 + (n - 1)d]\\). Substitute \\(n = %d\\), \\(a_1 = %d\\) and \\(d = %+d\\): \\(S_{%d} = \\frac{%d}{2}[2(%d) + (%d - 1)(%+d)]\\). Evaluate the bracketed terms \\(2(%d) + (%d)(%+d) = %d + (%d) = %d\\) and multiply by \\(\\frac{%d}{2}\\) to get \\(S_{%d} = %d\\).",
			n, a1, d, n, n, a1, n, d, a1, n-1, d, 2*a1, (n-1)*d, bracketVal, n, n, sum)
		break
	}
	return generator.Problem{
		Question:    q,
		Answer:      strconv.Itoa(sum),
		Explanation: ex,
		Facts: map[string]string{
			"n":      strconv.Itoa(n),
			"a1":     strconv.Itoa(a1),
			"d":      strconv.Itoa(d),
			"answer": strconv.Itoa(sum),
		},
	}
}

type geoSeriesGen struct{}

func (g *geoSeriesGen) Grade(expected, userAnswer string) grader.Result {
	return gradeNum(expected, userAnswer, 1e-6)
}

func (g *geoSeriesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	tier := mathutil.DifficultyTier(ctx.Difficulty)
	ratioPool := []int{2, 3, -2}
	if tier == 1 {
		ratioPool = []int{2, 3, -2, -3, 4}
	} else if tier == 2 {
		ratioPool = []int{2, 3, -2, -3, 4, 5}
	}
	var a1, r, n, sum, rPow int
	var q, ex string
	for {
		r = ratioPool[rand.Intn(len(ratioPool))]
		maxN := 6
		if abs(r) == 2 {
			maxN = 5 + tier*2
		} else if abs(r) == 3 {
			maxN = 4 + tier
		} else {
			maxN = 3 + tier
		}
		n = 3 + rand.Intn(maxN-2)
		a1 = 1 + rand.Intn(3+tier*2)
		rPow = 1
		for i := 0; i < n; i++ {
			rPow *= r
		}
		sum = a1 * (rPow - 1) / (r - 1)
		ansStr := strconv.Itoa(sum)
		q = fmt.Sprintf("Find the sum of the first %d terms of the geometric sequence \\(a_1=%d\\), \\(r=%d\\).", n, a1, r)
		if sum == n || sum == a1 || sum == r || sum == -r || containsWord(q, ansStr) {
			continue
		}
		ex = fmt.Sprintf("The sum of a finite geometric series with common ratio \\(r \\neq 1\\) is \\(S_n = a_1 \\frac{r^n - 1}{r - 1}\\). Substitute \\(a_1 = %d\\), \\(r = %d\\) and \\(n = %d\\): \\(S_{%d} = %d \\times \\frac{(%d)^{%d} - 1}{%d - 1}\\). Evaluate power \\((%d)^{%d} = %d\\), yielding \\(S_{%d} = %d \\times \\frac{%d - 1}{%d} = %d \\times \\frac{%d}{%d} = %d\\).",
			a1, r, n, n, a1, r, n, r, r, n, rPow, n, a1, rPow, r-1, a1, rPow-1, r-1, sum)
		break
	}
	return generator.Problem{
		Question:    q,
		Answer:      strconv.Itoa(sum),
		Explanation: ex,
		Facts: map[string]string{
			"n":      strconv.Itoa(n),
			"a1":     strconv.Itoa(a1),
			"r":      strconv.Itoa(r),
			"answer": strconv.Itoa(sum),
		},
	}
}

type geomInfiniteGen struct{}

func (g *geomInfiniteGen) Grade(expected, userAnswer string) grader.Result {
	return gradeNum(expected, userAnswer, 1e-6)
}

func (g *geomInfiniteGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	tier := mathutil.DifficultyTier(ctx.Difficulty)
	type frac struct{ p, q int }
	options := []frac{{1, 2}, {-1, 2}, {1, 3}, {1, 4}}
	if tier == 1 {
		options = []frac{{1, 2}, {-1, 2}, {1, 3}, {-1, 3}, {2, 3}, {1, 4}, {-1, 4}, {3, 4}}
	} else if tier == 2 {
		options = []frac{{1, 2}, {-1, 2}, {1, 3}, {-1, 3}, {2, 3}, {-2, 3}, {1, 4}, {-1, 4}, {3, 4}, {-3, 4}, {1, 5}, {2, 5}, {3, 5}, {4, 5}}
	}
	var a1, sum, gap int
	var f frac
	var qText, ex string
	for {
		f = options[rand.Intn(len(options))]
		gap = f.q - f.p
		k := 1 + rand.Intn(2+tier*3)
		a1 = gap * k
		sum = a1 * f.q / gap
		ansStr := strconv.Itoa(sum)
		qText = fmt.Sprintf("An infinite geometric series has first term %d and common ratio \\(%s\\). Give its sum.", a1, fracStr(0, f))
		if sum == a1 || sum == f.p || sum == -f.p || sum == f.q || containsWord(qText, ansStr) {
			continue
		}
		ex = fmt.Sprintf("An infinite geometric series converges if and only if \\(|r| < 1\\), with sum given by \\(S_\\infty = \\frac{a_1}{1 - r}\\). Verify convergence with \\(|%s| < 1\\) and substitute \\(a_1 = %d\\) and \\(r = %s\\): \\(S_\\infty = \\frac{%d}{1 - (%s)}\\). Simplify the denominator to \\(1 - (%s) = \\frac{%d}{%d}\\) and multiply by reciprocal: \\(%d \\times \\frac{%d}{%d} = %d\\).",
			fracStr(0, f), a1, fracStr(0, f), a1, fracStr(0, f), fracStr(0, f), gap, f.q, a1, f.q, gap, sum)
		break
	}
	return generator.Problem{
		Question:    qText,
		Answer:      strconv.Itoa(sum),
		Explanation: ex,
		Facts: map[string]string{
			"a1":     strconv.Itoa(a1),
			"answer": strconv.Itoa(sum),
		},
	}
}

func coefStrSigned(n int, denom int) string {
	v := float64(n) / float64(denom)
	if v == float64(int(v)) {
		return strconv.Itoa(int(v))
	}
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.4f", v), "0"), ".")
}

func coefFrac(n int) string {
	return coefStrSigned(n, 1)
}

func coefFracRat(n int, f struct{ p, q int }) string {
	num := n * f.p
	return coefStrSigned(num, f.q)
}

func fracStr(_ int, f struct{ p, q int }) string {
	sign := ""
	if f.p < 0 {
		sign = "-"
	}
	if f.q == 1 {
		return sign + strconv.Itoa(abs(f.p))
	}
	return fmt.Sprintf("%s\\frac{%d}{%d}", sign, abs(f.p), abs(f.q))
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

type floorCeilingGen struct{}

func (g *floorCeilingGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		question string
		answer   string
		reason   string
	}
	tableEasy := []entry{
		{"What is floor(3.7)? (enter a number)", "3", "Greatest integer at most 3.7."},
		{"What is ceil(3.2)? (enter a number)", "4", "Least integer at least 3.2."},
		{"What is floor(-2.3)? (enter a number)", "-3", "Greatest integer at most -2.3 (more negative)."},
		{"What is ceil(-2.8)? (enter a number)", "-2", "Least integer at least -2.8 (toward zero)."},
		{"Is floor(x) <= x always? (yes/no)", "yes", "Floor rounds down by definition."},
	}
	tableHard := []entry{
		{"What is floor(sqrt(50))? (enter a number)", "7", "7^2=49 <= 50 < 64 = 8^2."},
		{"What is ceil(pi)? (enter a number)", "4", "Smallest integer at least 3.14159."},
		{"How many integers n satisfy floor(n/2) = 3? (enter a number)", "2", "n = 6 and n = 7."},
	}
	pool := tableEasy
	if scale > 3 {
		pool = append(tableEasy, tableHard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.reason}
}

type heavisideGen struct{}

func (g *heavisideGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		question string
		answer   string
		reason   string
	}
	tableEasy := []entry{
		{"What is H(3) for the Heaviside step? (enter a number)", "1", "Positive argument: step is up."},
		{"What is H(-2) for the Heaviside step? (enter a number)", "0", "Negative argument: step is down."},
		{"Is the Heaviside step right-continuous at 0? (yes/no)", "yes", "Right limit equals the (usual) value 1."},
		{"Does H model an on/off switch at 0? (yes/no)", "yes", "0 before, 1 after: the canonical switch."},
	}
	tableHard := []entry{
		{"What is H(100) for the Heaviside step? (enter a number)", "1", "Positive argument: step is up."},
		{"Is the Heaviside step differentiable at 0? (yes/no)", "no", "Jump discontinuity: no derivative there (distributionally a delta)."},
		{"What is H(10) * H(-10)? (enter a number)", "0", "1 times 0."},
	}
	pool := tableEasy
	if scale > 3 {
		pool = append(tableEasy, tableHard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.reason}
}
