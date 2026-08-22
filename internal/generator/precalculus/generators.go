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
)

func Register(reg *generator.Registry) {
	for _, e := range all {
		reg.Register(e.id, e.gen)
	}
}

var all = []struct {
	id  string
	gen generator.Generator
}{
	{"precalc.vector.basics", &vectorBasicsGen{}},
	{"precalc.vector.ops", &vectorOpsGen{}},
	{"precalc.vector.dot", &vectorDotGen{}},
	{"precalc.polar.coordinates", &polarCoordsGen{}},
	{"precalc.func.transformations", &funcTransformationsGen{}},
	{"precalc.polar.graph", &polarGraphGen{}},
	{"precalc.parametric.graph", &parametricGraphGen{}},
	{"precalc.func.piecewise", &piecewiseGen{}},
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
	h := 1 + rand.Intn(2+int(ctx.Difficulty*5))
	k := 1 + rand.Intn(2+int(ctx.Difficulty*5))
	a := 2 + rand.Intn(3)
	all := []struct {
		form, label string
	}{
		{fmt.Sprintf("f(x-%d)", h), fmt.Sprintf("shift right by %d", h)},
		{fmt.Sprintf("f(x+%d)", h), fmt.Sprintf("shift left by %d", h)},
		{fmt.Sprintf("f(x)+%d", k), fmt.Sprintf("shift up by %d", k)},
		{fmt.Sprintf("f(x)-%d", k), fmt.Sprintf("shift down by %d", k)},
		{"-f(x)", "reflect across the x-axis"},
		{"f(-x)", "reflect across the y-axis"},
		{fmt.Sprintf("%df(x)", a), fmt.Sprintf("vertical stretch by a factor of %d", a)},
	}
	idx := rand.Intn(len(all))
	distractorIdx := make([]int, 0, len(all)-1)
	for i := range all {
		if i != idx {
			distractorIdx = append(distractorIdx, i)
		}
	}
	rand.Shuffle(len(distractorIdx), func(i, j int) { distractorIdx[i], distractorIdx[j] = distractorIdx[j], distractorIdx[i] })
	ds := []string{all[distractorIdx[0]].label, all[distractorIdx[1]].label, all[distractorIdx[2]].label}
	return mcProblem(
		fmt.Sprintf("Let \\(f\\) be any function. The graph of \\(g(x) = %s\\) is obtained from the graph of \\(f\\) by which single transformation?", all[idx].form),
		all[idx].label, ds)
}

type polarGraphGen struct{}

func (g *polarGraphGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := 1 + rand.Intn(2+int(ctx.Difficulty*4))
	deg := []string{"30^{\\circ}", "60^{\\circ}", "90^{\\circ}", "120^{\\circ}", "150^{\\circ}"}[rand.Intn(5)]
	circleAxis := fmt.Sprintf("a circle of diameter \\(|a|=%d\\) passing through the pole, centered on the polar axis", 2*a)
	circlePerp := fmt.Sprintf("a circle of diameter \\(|a|=%d\\) passing through the pole, centered on the line \\(\\theta=90^{\\circ}\\)", 2*a)
	line := fmt.Sprintf("a straight line through the pole making an angle \\(%s\\) with the polar axis", deg)
	centerCircle := fmt.Sprintf("a circle centered at the pole with radius \\(%d\\)", a)
	type variant struct{ form, correct, wrong1, wrong2, wrong3 string }
	v := []variant{
		{fmt.Sprintf("r = %d", a), centerCircle, line, circleAxis, circlePerp},
		{fmt.Sprintf("\\theta = %s", deg), line, centerCircle, circleAxis, circlePerp},
		{fmt.Sprintf("r = %d\\cos\\theta", a), circleAxis, circlePerp, centerCircle, line},
		{fmt.Sprintf("r = %d\\sin\\theta", a), circlePerp, circleAxis, centerCircle, line},
	}[rand.Intn(4)]
	return mcProblem(fmt.Sprintf("In polar coordinates \\((r,\\theta)\\), what shape does \\(%s\\) describe?", v.form),
		v.correct, []string{v.wrong1, v.wrong2, v.wrong3})
}

type parametricGraphGen struct{}

func (g *parametricGraphGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	h := 1 + rand.Intn(3)
	k := 1 + rand.Intn(3)
	type variant struct{ form, correct, d1, d2, d3 string }
	v := []variant{
		{fmt.Sprintf("x = t,\\ y = t^2"), "a parabola opening upward",
			"a straight line", "a parabola opening rightward", "a parabola opening downward"},
		{fmt.Sprintf("x = t^2,\\ y = t"), "a parabola opening rightward",
			"a straight line", "a parabola opening upward", "a parabola opening downward"},
		{fmt.Sprintf("x = t + %d,\\ y = t^2 + %d", h, k), "a parabola opening upward",
			"a parabola opening rightward", "a straight line", "a circle"},
		{fmt.Sprintf("x = t^2 - %d,\\ y = t + %d", k, h), "a parabola opening rightward",
			"a parabola opening upward", "a straight line", "a hyperbola"},
	}[rand.Intn(4)]
	return mcProblem(fmt.Sprintf("The curve given by \\(%s\\) (with \\(t\\) ranging over the reals) is best described as:", v.form),
		v.correct, []string{v.d1, v.d2, v.d3})
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
	ex := fmt.Sprintf("The term is \\(\\binom{%d}{%d} (%+d)^{%d} x^{%d}\\), giving coefficient \\(%s\\).", n, k, c, n-k, k, ans)
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
	ex := fmt.Sprintf("LHS at \\(n=%d\\): %s. RHS: substituting \\(n=%d\\) gives the same value \\(%s\\).", n, value, n, value)
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
	for q == p {
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
	parab := fmt.Sprintf("y = x^2 %+d", c)
	line := fmt.Sprintf("y = %dx %+d", m, bTotal)
	qText := fmt.Sprintf("Solve the system \\[%s\\\\ %s\\]\nEnter the solution with the smaller \\(x\\)-coordinate in the form \\((x, y)\\).", parab, line)
	ex := fmt.Sprintf("Substituting: \\(x^2 %+d = %dx %+d\\Rightarrow x^2 -%dx %+d = 0\\), so \\(x\\in\\{%d, %d\\}\\). For \\(x=%d\\), \\(y=%d\\).", c, m, bTotal, m, -bTotal, p, q, smaller, ySmall)
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
		q := fmt.Sprintf("The parabola \\(y^2 = %sx^\\) has its focus on the x-axis at \\((p, 0)\\). Find \\(p\\).", coefStr(fourP))
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
	if rand.Intn(2) == 0 {
		a, b = b, a
	}
	if ctx.Difficulty > 0.5 && rand.Intn(2) == 0 {
		q := fmt.Sprintf("An ellipse is given by \\(\\frac{x^2}{%d} + \\frac{y^2}{%d} = 1\\). Find the distance from the center to each focus \\(c\\).", a*a, b*b)
		ex := fmt.Sprintf("\\(c^2 = a^2-b^2 = %d-%d = %d\\), so \\(c = %d\\).", maxInt(a*a, b*b), minInt(a*a, b*b), c*c, c)
		return generator.Problem{Question: q, Answer: strconv.Itoa(c), Explanation: ex}
	}
	q := fmt.Sprintf("An ellipse has semi-axes \\(a=%d\\) and \\(b=%d\\) (with \\(a>b\\)). Find \\(c\\), the distance from center to focus, where \\(c^2=a^2-b^2\\).", a, b)
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
	a1 := rand.Intn(2*(1+int(ctx.Difficulty*5))+1) - (1 + int(ctx.Difficulty*5))
	d := 0
	for d == 0 {
		d = rand.Intn(2*(1+int(ctx.Difficulty*4))+1) - (1 + int(ctx.Difficulty*4))
	}
	n := 5 + rand.Intn(10+int(ctx.Difficulty*15))
	an := a1 + (n-1)*d
	q := fmt.Sprintf("An arithmetic sequence starts at \\(%d\\) with common difference \\(%+d\\). Find term \\(a_{%d}\\).", a1, d, n)
	ex := fmt.Sprintf("\\(a_n=a_1+(n-1)d=%d+%d\\cdot%d=%d\\).", a1, n-1, d, an)
	return generator.Problem{Question: q, Answer: strconv.Itoa(an), Explanation: ex}
}

type geoSeqGen struct{}

func (g *geoSeqGen) Grade(expected, userAnswer string) grader.Result {
	return gradeNum(expected, userAnswer, 1e-6)
}

func (g *geoSeqGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	r := []int{2, 3, -2}[rand.Intn(3)]
	maxN := map[int]int{2: 10, 3: 6, -2: 8}[r]
	n := 2 + rand.Intn(maxN-1)
	a1 := 1 + rand.Intn(3+int(ctx.Difficulty*4))
	an := a1
	for i := 1; i < n; i++ {
		an *= r
	}
	q := fmt.Sprintf("A geometric sequence starts at \\(%d\\) with common ratio \\(%d\\). Find term \\(a_{%d}\\).", a1, r, n)
	ex := fmt.Sprintf("\\(a_n=a_1 r^{n-1}=%d\\cdot(%d)^{%d}=%d\\).", a1, r, n-1, an)
	return generator.Problem{Question: q, Answer: strconv.Itoa(an), Explanation: ex}
}

type summationGen struct{}

func (g *summationGen) Grade(expected, userAnswer string) grader.Result {
	return gradeNum(expected, userAnswer, 1e-6)
}

func (g *summationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	n := 4 + rand.Intn(3+int(ctx.Difficulty*8))
	a := 1 + rand.Intn(3)
	b := rand.Intn(2*3+1) - 3
	total := 0
	for i := 1; i <= n; i++ {
		total += a*i + b
	}
	q := fmt.Sprintf("Evaluate \\(\\sum_{i=1}^{%d} (%si %+d)\\).", n, coefStr(a), b)
	ex := fmt.Sprintf("\\(\\sum (ai+b)=a\\cdot\\frac{n(n+1)}{2}+bn=%d\\cdot\\frac{%d\\cdot%d}{2}+(%+d)\\cdot%d=%d\\).", a, n, n+1, b, n, total)
	return generator.Problem{Question: q, Answer: strconv.Itoa(total), Explanation: ex}
}

type arithSeriesGen struct{}

func (g *arithSeriesGen) Grade(expected, userAnswer string) grader.Result {
	return gradeNum(expected, userAnswer, 1e-6)
}

func (g *arithSeriesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a1 := rand.Intn(11) - 5
	d := 0
	for d == 0 {
		d = rand.Intn(9) - 4
	}
	n := 4 + rand.Intn(8)
	an := a1 + (n-1)*d
	sum := n * (a1 + an) / 2
	q := fmt.Sprintf("Find the sum of the first \\(%d\\) terms of the arithmetic sequence \\(a_1=%d\\), \\(d=%+d\\).", n, a1, d)
	ex := fmt.Sprintf("\\(S_n=\\frac{n(a_1+a_n)}{2}=\\frac{%d(%d %+d)}{2}=%d\\).", n, a1, (n-1)*d, sum)
	return generator.Problem{Question: q, Answer: strconv.Itoa(sum), Explanation: ex}
}

type geoSeriesGen struct{}

func (g *geoSeriesGen) Grade(expected, userAnswer string) grader.Result {
	return gradeNum(expected, userAnswer, 1e-6)
}

func (g *geoSeriesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	r := []int{2, 3, -2}[rand.Intn(3)]
	n := 3 + rand.Intn(5)
	a1 := 1 + rand.Intn(4)
	sum := a1
	term := a1
	for i := 1; i < n; i++ {
		term *= r
		sum += term
	}
	q := fmt.Sprintf("Find the sum of the first \\(%d\\) terms of the geometric sequence \\(a_1=%d\\), \\(r=%d\\).", n, a1, r)
	ex := fmt.Sprintf("\\(S_n=a_1\\frac{r^{n}-1}{r-1}=%d\\cdot\\frac{(%d)^{%d}-1}{%d}=%d\\).", a1, r, n, r-1, sum)
	return generator.Problem{Question: q, Answer: strconv.Itoa(sum), Explanation: ex}
}

type geomInfiniteGen struct{}

func (g *geomInfiniteGen) Grade(expected, userAnswer string) grader.Result {
	return gradeNum(expected, userAnswer, 1e-6)
}

func (g *geomInfiniteGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type frac struct{ p, q int }
	options := []frac{{1, 2}, {-1, 2}, {1, 3}, {1, 4}, {-1, 4}, {2, 3}}
	f := options[rand.Intn(len(options))]
	gap := f.q - f.p
	a1 := gap * (1 + rand.Intn(2+int(ctx.Difficulty*4)))
	sum := a1 * f.q / gap
	qText := fmt.Sprintf("An infinite geometric series has first term \\(%d\\) and common ratio \\(%s\\). Give its sum.", a1, fracStr(a1, f))
	ex := fmt.Sprintf("Since \\(|r|<1\\), the series converges: \\(S=\\frac{a_1}{1-r}=\\frac{%d}{1-%s}=%d\\).", a1, fracStr(0, f), sum)
	return generator.Problem{Question: qText, Answer: strconv.Itoa(sum), Explanation: ex}
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
