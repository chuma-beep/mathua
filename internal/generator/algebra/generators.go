package algebra

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/mathutil"
)

func Register(reg *generator.Registry) {
	reg.Register("alg.linear.slope", &slopeGen{})
	reg.Register("alg.linear.slope_intercept", &slopeInterceptGen{})
	reg.Register("alg.linear.graph", &linearGraphGen{})
	reg.Register("alg.linear.standard_form", &stdFormGen{})
	reg.Register("alg.linear.parallel_perp", &parallelPerpGen{})

	reg.Register("alg.eq.multi_step", &multiStepEqGen{})
	reg.Register("alg.eq.vars_both_sides", &varsBothSidesGen{})
	reg.Register("alg.eq.vars_both_sides.word", &varsBothSidesWordGen{})
	reg.Register("alg.eq.literal", &literalEqGen{})
	reg.Register("alg.ineq.multi_step", &multiStepIneqGen{})
	reg.Register("alg.ineq.multi_step.word", &ineqWordGen{})
	reg.Register("alg.ineq.compound", &compoundIneqGen{})

	reg.Register("alg.systems.substitution", &sysSubstitutionGen{})
	reg.Register("alg.systems.elimination", &sysEliminationGen{})
	reg.Register("alg.systems.word", &sysWordGen{})

	reg.Register("alg.poly.concept", &polyConceptGen{})
	reg.Register("alg.poly.add_sub", &polyAddSubGen{})
	reg.Register("alg.poly.mult_mono", &polyMultMonoGen{})
	reg.Register("alg.poly.foil", &polyFoilGen{})
	reg.Register("alg.poly.special", &polySpecialGen{})

	reg.Register("alg.factor.gcf", &factorGCFGen{})
	reg.Register("alg.factor.trinomial", &factorTrinomialGen{})
	reg.Register("alg.factor.diff_squares", &factorDiffSquaresGen{})
	reg.Register("alg.factor.ac_method", &factorACMethodGen{})

	reg.Register("alg.quad.solve_factor", &quadSolveFactorGen{})
	reg.Register("alg.quad.complete_square", &quadCompleteSquareGen{})
	reg.Register("alg.quad.formula", &quadFormulaGen{})
	reg.Register("alg.quad.discriminant", &quadDiscriminantGen{})

	reg.Register("alg.func.concept", &funcConceptGen{})
	reg.Register("alg.func.notation", &funcNotationGen{})
	reg.Register("alg.func.evaluate", &funcEvaluateGen{})
	reg.Register("alg.func.linear", &funcLinearGen{})
	reg.Register("alg.func.quad", &funcQuadGen{})

	reg.Register("alg.exp.concept", &algExpConceptGen{})
	reg.Register("alg.exp.evaluate", &algExpEvaluateGen{})
	reg.Register("alg.log.concept", &logConceptGen{})
	reg.Register("alg.log.evaluate", &logEvaluateGen{})
	reg.Register("alg.log.properties", &logPropertiesGen{})

	reg.Register("alg.seq.arithmetic", &seqArithGen{})
	reg.Register("alg.seq.geometric", &seqGeomGen{})
	reg.Register("alg.seq.sum_arith", &seqSumArithGen{})
	reg.Register("alg.seq.sum_geo", &seqSumGeoGen{})

	reg.Register("alg.ineq.two_var", &ineqTwoVarGen{})
	reg.Register("alg.conic.circle", &conicCircleGen{})
	reg.Register("alg.conic.ellipse", &conicEllipseGen{})

	reg.Register("alg.conic.parabola", &conicParabolaGen{})
	reg.Register("alg.conic.hyperbola", &conicHyperbolaGen{})
	reg.Register("alg.eq.absolute_value", &eqAbsValGen{})
	reg.Register("alg.eq.binomial", &eqBinomialGen{})
	reg.Register("alg.eq.exponential", &eqExpGen{})
	reg.Register("alg.eq.irrational", &eqIrrationalGen{})
	reg.Register("alg.eq.logarithmic", &eqLogGen{})
	reg.Register("alg.eq.polynomial", &eqPolyGen{})
	reg.Register("alg.eq.rational", &eqRationalGen{})
	reg.Register("alg.eq.trinomial", &eqTrinomialGen{})
	reg.Register("alg.func.absolute_value", &funcAbsValGen{})
	reg.Register("alg.func.composite", &funcCompositeGen{})
	reg.Register("alg.func.dirichlet", &funcDirichletGen{})
	reg.Register("alg.func.even_odd", &funcEvenOddGen{})
	reg.Register("alg.func.inverse", &funcInverseGen{})
	reg.Register("alg.func.monotonicity", &monotonicityGen{})
	reg.Register("alg.func.rational", &funcRationalGen{})
	reg.Register("alg.func.sigmoid", &funcSigmoidGen{})
	reg.Register("alg.func.sign", &funcSignGen{})
	reg.Register("alg.ineq.absolute_value", &ineqAbsValGen{})
	reg.Register("alg.ineq.interval", &ineqIntervalGen{})
	reg.Register("alg.ineq.irrational", &ineqIrrationalGen{})
	reg.Register("alg.ineq.logarithmic", &ineqLogGen{})
	reg.Register("alg.ineq.quadratic", &ineqQuadraticGen{})
	reg.Register("alg.ineq.rational", &ineqRationalGen{})
	reg.Register("alg.ineq.sign_analysis", &signAnalysisGen{})
	reg.Register("alg.ineq.systems", &ineqSystemsGen{})
	reg.Register("alg.poly.division", &polyDivisionGen{})
	reg.Register("alg.poly.monomial", &polyMonomialGen{})
	reg.Register("alg.poly.roots", &polyRootsGen{})
	reg.Register("alg.poly.synthetic_division", &synthDivGen{})
	reg.Register("alg.poly.vieta", &vietaGen{})
	reg.Register("alg.quad.complex", &quadComplexGen{})
	reg.Register("alg.quad.incomplete", &quadIncompleteGen{})
	reg.Register("alg.quad.parametric", &quadParametricGen{})
	reg.Register("alg.quad.quadratic", &quadQuadraticGen{})
	reg.Register("alg.systems.gaussian", &gaussianElimGen{})

	registerBatch(reg)
	registerRaise(reg)
	registerConversions(reg)
}

// ---------------------------------------------------------------------------
// Linear
// ---------------------------------------------------------------------------

type slopeGen struct{}

func (g *slopeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// Higher ctx.Difficulty → wider range, negative slopes more likely
	lim := int(3 + ctx.Difficulty*12)
	x1 := rand.Intn(lim)
	y1 := rand.Intn(lim)
	x2 := x1 + rand.Intn(int(ctx.Difficulty*7)+2) + 1
	y2 := y1 + rand.Intn(lim) + 1
	if rand.Intn(2) == 0 || ctx.Difficulty > 0.6 {
		y2 = y1 - rand.Intn(max(y1, 1)+1)
	}
	dy := y2 - y1
	dx := x2 - x1
	frac := reduce(dy, dx)
	num, den := frac.num, frac.den
	return generator.Problem{
		Question:    fmt.Sprintf("Find the slope between \\((%d,%d)\\) and \\((%d,%d)\\).", x1, y1, x2, y2),
		Answer:      fracOrInt(num, den),
		Explanation: fmt.Sprintf("\\(m = \\frac{%d - %d}{%d - %d} = \\frac{%d}{%d} = %s\\).", y2, y1, x2, x1, dy, dx, fracOrInt(num, den)),
		Facts:       map[string]string{"x1": strconv.Itoa(x1), "y1": strconv.Itoa(y1), "x2": strconv.Itoa(x2), "y2": strconv.Itoa(y2), "dy": strconv.Itoa(dy), "dx": strconv.Itoa(dx), "slope": fracOrInt(num, den)},
	}
}

type slopeInterceptGen struct{}

func (g *slopeInterceptGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// Higher ctx.Difficulty → fractional slopes and larger intercepts
	isFraction := false
	if ctx.KPOrdinal == 3 {
		isFraction = true
	} else if ctx.KPOrdinal == 1 || ctx.KPOrdinal == 2 {
		isFraction = false
	} else if ctx.Difficulty > 0.6 && rand.Intn(2) == 0 {
		isFraction = true
	}
	if isFraction {
		m := rand.Intn(3) + 1
		n := rand.Intn(3) + 2
		if rand.Intn(2) == 0 {
			m = -m
		}
		b := rand.Intn(int(5+ctx.Difficulty*10)) - int(3+ctx.Difficulty*5)
		return generator.Problem{
			Question:    fmt.Sprintf("Write the equation of a line with slope \\(%d/%d\\) and \\(y\\)-intercept %d (\\(y = mx + b\\)).", m, n, b),
			Answer:      fmt.Sprintf("y = (%d/%d)x + %d", m, n, b),
			Explanation: fmt.Sprintf("\\(y = \\frac{%d}{%d}x + %d\\)", m, n, b),
			Facts:       map[string]string{"m_num": strconv.Itoa(m), "m_den": strconv.Itoa(n), "b": strconv.Itoa(b), "which": "fraction"},
		}
	}
	m := rand.Intn(int(1+ctx.Difficulty*5)) + 1
	if rand.Intn(2) == 0 {
		m = -m
	}
	b := rand.Intn(int(5+ctx.Difficulty*8)) - int(3+ctx.Difficulty*4)
	return generator.Problem{
		Question:    fmt.Sprintf("Write the equation of a line with slope %d and \\(y\\)-intercept %d (\\(y = mx + b\\)).", m, b),
		Answer:      formatLinear(m, b),
		Explanation: fmt.Sprintf("\\(y = %dx + %d\\)", m, b),
		Facts:       map[string]string{"m_num": strconv.Itoa(m), "m_den": "1", "b": strconv.Itoa(b), "which": "integer"},
	}
}

type linearGraphGen struct{}

func (g *linearGraphGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	m := rand.Intn(5) + 1
	if rand.Intn(2) == 0 {
		m = -m
	}
	b := rand.Intn(10) - 5
	x := rand.Intn(5) + 1
	y := m*x + b
	return generator.Problem{
		Question:    fmt.Sprintf("Given \\(%s\\), find \\(y\\) when \\(x = %d\\).", formatLinear(m, b), x),
		Answer:      fmt.Sprintf("%d", y),
		Explanation: fmt.Sprintf("\\(y = %d(%d) + %d = %d + %d = %d\\)", m, x, b, m*x, b, y),
		Facts:       map[string]string{"m": strconv.Itoa(m), "b": strconv.Itoa(b), "x": strconv.Itoa(x), "mx": strconv.Itoa(m * x), "y": strconv.Itoa(y)},
	}
}

type stdFormGen struct{}

func (g *stdFormGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := rand.Intn(5) + 2
	b := rand.Intn(5) + 1
	x := rand.Intn(5) + 1
	y := rand.Intn(5) + 1
	total := a*x + b*y
	return generator.Problem{
		Question:    fmt.Sprintf("If \\((%d,%d)\\) satisfies \\(%dx + %dy = c\\), what is \\(c\\)? (enter a number)", x, y, a, b),
		Answer:      strconv.Itoa(total),
		Explanation: fmt.Sprintf("Substitute the point: \\(%d(%d) + %d(%d) = %d\\), so \\(c = %d\\).", a, x, b, y, total, total),
		Facts: map[string]string{
			"x": strconv.Itoa(x), "y": strconv.Itoa(y), "a": strconv.Itoa(a), "b": strconv.Itoa(b),
			"c": strconv.Itoa(total), "ax": strconv.Itoa(a * x), "by": strconv.Itoa(b * y), "total": strconv.Itoa(total),
		},
	}
}

type parallelPerpGen struct{}

func (g *parallelPerpGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	m := rand.Intn(5) + 1
	if rand.Intn(2) == 0 {
		m = -m
	}
	pick := rand.Intn(2)
	if pick == 0 {
		return generator.Problem{
			Question:    fmt.Sprintf("A line parallel to \\(y = %dx + 3\\) has slope:", m),
			Answer:      fmt.Sprintf("%d", m),
			Explanation: fmt.Sprintf("Parallel lines have equal slopes. Slope = %d.", m),
			Facts:       map[string]string{"m": strconv.Itoa(m), "which": "parallel"},
		}
	}
	perp := reduce(-1, m)
	return generator.Problem{
		Question:    fmt.Sprintf("A line perpendicular to \\(y = %dx + 3\\) has slope:", m),
		Answer:      fracOrInt(perp.num, perp.den),
		Explanation: fmt.Sprintf("Perpendicular slopes are negative reciprocals: \\(-\\frac{1}{%d} = %s\\)", m, fracOrInt(perp.num, perp.den)),
		Facts:       map[string]string{"m": strconv.Itoa(m), "which": "perpendicular", "perp": fracOrInt(perp.num, perp.den)},
	}
}

// ---------------------------------------------------------------------------
// Equations & inequalities
// ---------------------------------------------------------------------------

type multiStepEqGen struct{}

func (g *multiStepEqGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	d := ctx.Difficulty
	// Staged ladder: plain ax+b=c → larger/subtraction → parentheses probe.
	if ctx.KPOrdinal == 3 || (ctx.KPOrdinal == 0 && d >= 0.75) {
		a := rand.Intn(4) + 2
		x := rand.Intn(8) + 2
		b := rand.Intn(8) + 1
		inner := x + b
		c := a * inner
		return generator.Problem{
			Question:    fmt.Sprintf("Solve: \\(%d(x + %d) = %d\\)", a, b, c),
			Answer:      fmt.Sprintf("%d", x),
			Explanation: fmt.Sprintf("Given %d(x + %d) = %d. Step 1: divide both sides by %d: x + %d = %d. Step 2: subtract %d: x = %d. Answer: %d", a, b, c, a, b, inner, b, x, x),
			Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "c": strconv.Itoa(c), "inner": strconv.Itoa(inner), "x": strconv.Itoa(x), "which": "distributed"},
		}
	}
	scale := int(1 + d*8)
	if ctx.KPOrdinal == 1 {
		scale = 1
	} else if ctx.KPOrdinal == 2 {
		scale = 4
	}
	x := rand.Intn(scale*2) + 2
	a := rand.Intn(scale) + 2
	b := rand.Intn(scale*2) + 2
	c := a*x + b
	return generator.Problem{
		Question:    fmt.Sprintf("Solve: \\(%dx + %d = %d\\)", a, b, c),
		Answer:      fmt.Sprintf("%d", x),
		Explanation: fmt.Sprintf("Given %dx + %d = %d. Step 1: subtract %d: %dx = %d. Step 2: divide by %d: x = %d. Answer: %d", a, b, c, b, a, c-b, a, x, x),
		Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "c": strconv.Itoa(c), "ax": strconv.Itoa(c - b), "x": strconv.Itoa(x), "which": "axb"},
	}
}

type varsBothSidesGen struct{}

func (g *varsBothSidesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// Staged ladder: small coefficients → larger spread → bigger solutions.
	coefMax := 5
	if ctx.KPOrdinal == 1 {
		coefMax = 4
	} else if ctx.KPOrdinal == 2 {
		coefMax = 7
	} else if ctx.KPOrdinal == 3 {
		coefMax = 10
	} else if ctx.Difficulty >= 0.7 {
		coefMax = 9
	} else if ctx.Difficulty >= 0.45 {
		coefMax = 7
	}
	scale := coefMax
	x := rand.Intn(scale*2) + 2
	a := rand.Intn(scale) + 2
	c := rand.Intn(scale) + 1
	for a == c {
		c++
	}
	b := rand.Intn(scale*2) + 1
	d := (a-c)*x + b
	return generator.Problem{
		Question:    fmt.Sprintf("Solve: \\(%dx + %d = %dx + %d\\)", a, b, c, d),
		Answer:      fmt.Sprintf("%d", x),
		Explanation: fmt.Sprintf("\\(%dx + %d = %dx + %d\\) → \\(%dx = %d\\) → \\(x = %d\\)", a, b, c, d, a-c, d-b, x),
		Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "c": strconv.Itoa(c), "d": strconv.Itoa(d), "diff": strconv.Itoa(a - c), "rhs": strconv.Itoa(d - b), "x": strconv.Itoa(x)},
	}
}

type varsBothSidesWordGen struct{}

// A story whose two offers meet at one number of uses. The meeting point is
// chosen first and the second flat fee is derived from it, so the number the
// learner has to find is an exact integer and the pool discriminates.
func (g *varsBothSidesWordGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	maxPer, maxUses := 9, 8
	if ctx.KPOrdinal == 1 {
		maxPer, maxUses = 6, 5
	} else if ctx.KPOrdinal == 2 {
		maxPer, maxUses = 13, 12
	} else if ctx.KPOrdinal == 3 {
		maxPer, maxUses = 18, 20
	} else {
		if ctx.Difficulty >= 0.45 {
			maxPer, maxUses = 13, 12
		}
		if ctx.Difficulty >= 0.7 {
			maxPer, maxUses = 18, 20
		}
	}
	perA := rand.Intn(maxPer) + 6     // 6..maxPer+5
	perB := perA - (rand.Intn(3) + 1) // 1..3 cheaper per use
	uses := rand.Intn(maxUses) + 3    // 3..maxUses+2
	joinA := rand.Intn(25) + 8        // 8..32
	joinB := joinA + (perA-perB)*uses // B is dearer to join, to compensate
	pairs := []struct{ a, b, unit, plural, thing string }{
		{"Gym A", "Gym B", "visit", "visits", "gym"},
		{"Plan A", "Plan B", "month", "months", "phone plan"},
		{"Studio A", "Studio B", "class", "classes", "studio"},
		{"Club A", "Club B", "round", "rounds", "golf club"},
	}
	p := pairs[rand.Intn(len(pairs))]
	return generator.Problem{
		Question: fmt.Sprintf("%s charges $%d to join and $%d per %s. %s charges $%d to join and $%d per %s. After how many %s do the two %ss cost the same?",
			p.a, joinA, perA, p.unit, p.b, joinB, perB, p.unit, p.plural, p.thing),
		Answer: strconv.Itoa(uses),
		Explanation: fmt.Sprintf("\\(%d + %dn = %d + %dn\\) → \\(%dn = %d\\) → \\(n = %d\\).",
			joinA, perA, joinB, perB, perA-perB, joinB-joinA, uses),
		Facts: map[string]string{
			"a": p.a, "b": p.b, "unit": p.unit, "thing": p.thing,
			"joinA": strconv.Itoa(joinA), "perA": strconv.Itoa(perA),
			"joinB": strconv.Itoa(joinB), "perB": strconv.Itoa(perB),
			"diff": strconv.Itoa(perA - perB), "gap": strconv.Itoa(joinB - joinA),
			"uses": strconv.Itoa(uses),
		},
	}
}

// ineqWordGen: a budget story answered by the greatest whole number that fits.
// The budget is deliberately off the boundary, so "at most" has one integer
// answer instead of sitting exactly on a step.
type ineqWordGen struct{}

func (g *ineqWordGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	maxPer, maxCount := 5, 8
	if ctx.KPOrdinal == 1 {
		maxPer, maxCount = 4, 6
	} else if ctx.KPOrdinal == 2 {
		maxPer, maxCount = 9, 12
	} else if ctx.KPOrdinal == 3 {
		maxPer, maxCount = 15, 20
	} else {
		if ctx.Difficulty >= 0.45 {
			maxPer, maxCount = 9, 12
		}
		if ctx.Difficulty >= 0.7 {
			maxPer, maxCount = 15, 20
		}
	}
	base := rand.Intn(15) + 5        // flat fee 5..19
	per := rand.Intn(maxPer-1) + 2   // 2..maxPer
	count := rand.Intn(maxCount) + 3 // 3..maxCount+2
	slack := rand.Intn(per)          // < per, so the answer is unique
	budget := base + per*count + slack
	room := budget - base // per*count + slack
	stories := []struct{ service, unit string }{
		{"taxi", "mile"},
		{"delivery van", "kilometre"},
		{"moving van", "mile"},
		{"caterer", "guest"},
	}
	s := stories[rand.Intn(len(stories))]
	return generator.Problem{
		Question: fmt.Sprintf("A %s charges $%d plus $%d per %s. You can spend at most $%d. What is the greatest whole number of %ss you can pay for?",
			s.service, base, per, s.unit, budget, s.unit),
		Answer: strconv.Itoa(count),
		Explanation: fmt.Sprintf("\\(%d + %dm \\le %d\\) → \\(%dm \\le %d\\) → \\(m \\le %d\\), and %d %ss would cost %d, more than the budget allows, so \\(m = %d\\).",
			base, per, budget, per, room, room/per, count+1, s.unit, base+per*(count+1), count),
		Facts: map[string]string{
			"service": s.service, "unit": s.unit, "base": strconv.Itoa(base),
			"per": strconv.Itoa(per), "budget": strconv.Itoa(budget),
			"room": strconv.Itoa(room), "quotient": strconv.Itoa(room / per),
			"next": strconv.Itoa(count + 1), "next_cost": strconv.Itoa(base + per*(count+1)),
			"count": strconv.Itoa(count),
		},
	}
}

type literalEqGen struct{}

func (g *literalEqGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := rand.Intn(6) + 2
	b := rand.Intn(5) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("Solve for \\(y\\): \\(%dx + %dy = z\\)", a, b),
		Answer:      fmt.Sprintf("y = (z - %dx)/%d", a, b),
		Explanation: fmt.Sprintf("\\(%dy = z - %dx\\) → \\(y = \\frac{z - %dx}{%d}\\)", b, a, a, b),
		Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b)},
	}
}

type multiStepIneqGen struct{}

// multiStepIneqGen asks for the boundary of a two-step inequality: a x + b > c
// has the solution x > x when c - b is exactly a times the boundary. The earlier
// version added a random remainder to c, so the true boundary was a fraction
// (7x + 4 > 31 gives x > 27/7) while the answer was a whole number, and a learner
// who typed 27/7 or 3 6/7, was marked wrong for being right. The remainder is
// gone and the arithmetic is exact.
func (g *multiStepIneqGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	x := rand.Intn(8) + 2
	a := rand.Intn(6) + 2
	b := rand.Intn(10) + 1
	c := a*x + b
	return generator.Problem{
		Question:    fmt.Sprintf("Solve \\(%dx + %d > %d\\). What is the boundary of the solution set, the number that \\(x\\) must exceed?", a, b, c),
		Answer:      fmt.Sprintf("%d", x),
		Explanation: fmt.Sprintf("\\(%dx + %d > %d\\) → \\(%dx > %d\\) → \\(x > %d\\)", a, b, c, a, a*x, x),
		Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "c": strconv.Itoa(c), "rhs": strconv.Itoa(c - b), "x": strconv.Itoa(x)},
	}
}

type compoundIneqGen struct{}

func (g *compoundIneqGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	x := rand.Intn(8) + 2
	b := rand.Intn(3) + 2
	c := rand.Intn(5) - 2
	// We want a < bx + c < d with exactly one integer solution x.
	// x satisfies: a < bx + c < d
	// x-1 fails: a >= b(x-1) + c
	// x+1 fails: b(x+1) + c >= d
	// So: a = b(x-1) + c + 1, d = bx + c + 1
	a := b*(x-1) + c + 1
	d := b*x + c + 1
	return generator.Problem{
		Question:    fmt.Sprintf("Solve: \\(%d < %dx + %d < %d\\)", a, b, c, d),
		Answer:      fmt.Sprintf("%d", x),
		Explanation: fmt.Sprintf("\\(%d < %dx + %d < %d\\) → \\(%.1f < x < %.1f\\) → \\(x = %d\\)", a, b, c, d, float64(a-c)/float64(b), float64(d-c)/float64(b), x),
		Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "c": strconv.Itoa(c), "d": strconv.Itoa(d), "lo": fmt.Sprintf("%.1f", float64(a-c)/float64(b)), "hi": fmt.Sprintf("%.1f", float64(d-c)/float64(b)), "x": strconv.Itoa(x)},
	}
}

// ---------------------------------------------------------------------------
// Systems
// ---------------------------------------------------------------------------

type sysSubstitutionGen struct{}

func (g *sysSubstitutionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	x := rand.Intn(8) + 2
	y := rand.Intn(8) + 2
	return generator.Problem{
		// "Solve:" alone invited half an answer: a learner found x and wrote
		// "x = 2", which a tuple graded wrong because it wanted the pair
		// (report 147). The stem now says which form to use.
		Question:    fmt.Sprintf("Solve the system \\(y = %d\\), \\(x + y = %d\\). Give both values as \\((x, y)\\).", y, x+y),
		Answer:      fmt.Sprintf("(%d,%d)", x, y),
		Explanation: fmt.Sprintf("Substitute \\(y=%d\\): \\(x + %d = %d\\) → \\(x = %d\\). Solution: \\((%d,%d)\\)", y, y, x+y, x, x, y),
		Facts:       map[string]string{"y": strconv.Itoa(y), "sum": strconv.Itoa(x + y), "x": strconv.Itoa(x)},
	}
}

type sysEliminationGen struct{}

func (g *sysEliminationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	x := rand.Intn(8) + 2
	y := rand.Intn(8) + 2
	a := rand.Intn(5) + 2
	b := rand.Intn(5) + 2
	// The x coefficients must differ, or both equations are the same one and
	// the system has no unique solution: report 229 quoted
	// "6x + 4y = 50, 6x + 4y = 50" and no answer could be right.
	c := a
	for c == a {
		c = rand.Intn(5) + 2
	}
	eq1 := a*x + b*y
	eq2 := c*x + b*y
	return generator.Problem{
		Question:    fmt.Sprintf("Solve the system \\(%dx + %dy = %d\\), \\(%dx + %dy = %d\\). Give both values as \\((x, y)\\).", a, b, eq1, c, b, eq2),
		Answer:      fmt.Sprintf("(%d,%d)", x, y),
		Explanation: fmt.Sprintf("Subtract: \\(%dx = %d\\) → \\(x = %d\\); \\(y = %d\\). Solution: \\((%d,%d)\\)", a-c, eq1-eq2, x, y, x, y),
		Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "c": strconv.Itoa(c), "eq1": strconv.Itoa(eq1), "eq2": strconv.Itoa(eq2), "coef": strconv.Itoa(a - c), "rhs": strconv.Itoa(eq1 - eq2), "x": strconv.Itoa(x), "y": strconv.Itoa(y)},
	}
}

type sysWordGen struct{}

func (g *sysWordGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// The numbers are drawn so the bigger one comes first in the algebra and
	// the answer, which is what the question asks for: "smaller first". Before
	// this the pair was drawn independently, so a draw with a < b printed a
	// negative difference and answered larger-first, and a learner who read the
	// question and wrote the smaller number first was marked wrong (report 218).
	small := rand.Intn(8) + 3
	big := small + rand.Intn(6) + 1
	sum := big + small
	diff := big - small
	return generator.Problem{
		Question:    fmt.Sprintf("Two numbers sum to %d and differ by %d. Find both numbers (smaller first).", sum, diff),
		Answer:      fmt.Sprintf("%d,%d", small, big),
		Explanation: fmt.Sprintf("\\(x + y = %d\\), \\(x - y = %d\\). Add: \\(2x = %d\\) → \\(x = %d\\). \\(y = %d\\)", sum, diff, sum+diff, big, small),
		Facts:       map[string]string{"sum": strconv.Itoa(sum), "diff": strconv.Itoa(diff), "two_x": strconv.Itoa(sum + diff), "big": strconv.Itoa(big), "small": strconv.Itoa(small)},
	}
}

// ---------------------------------------------------------------------------
// Polynomials
// ---------------------------------------------------------------------------

type polyConceptGen struct{}

func (g *polyConceptGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// The schema describes the leading term, so the degree question stays and
	// the coefficient and exponent move. The exponent runs to 9 because above
	// that the printed term stops being readable at a glance.
	coeff := rand.Intn(11) + 2
	exp := rand.Intn(8) + 2
	return generator.Problem{
		Question:    fmt.Sprintf("What is the degree of \\(%dx^{%d} + 3x + 1\\)?", coeff, exp),
		Answer:      fmt.Sprintf("%d", exp),
		Explanation: fmt.Sprintf("The highest exponent is %d, so the degree is %d.", exp, exp),
		Facts:       map[string]string{"coeff": strconv.Itoa(coeff), "exp": strconv.Itoa(exp)},
	}
}

type polyAddSubGen struct{}

func (g *polyAddSubGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := rand.Intn(5) + 2
	b := rand.Intn(5) + 2
	c := rand.Intn(5) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("Simplify: \\((%dx^{2} + %dx) + (%dx^{2} + %dx)\\)", a, c, b, c),
		Answer:      fmt.Sprintf("%dx^2 + %dx", a+b, 2*c),
		Explanation: fmt.Sprintf("\\(%dx^{2} + %dx^{2} = %dx^{2}\\); \\(%dx + %dx = %dx\\)", a, b, a+b, c, c, 2*c),
		Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "c": strconv.Itoa(c), "a2": strconv.Itoa(a + b), "c1": strconv.Itoa(2 * c)},
	}
}

type polyMultMonoGen struct{}

func (g *polyMultMonoGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := rand.Intn(5) + 2
	b := rand.Intn(5) + 2
	c := rand.Intn(5) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("Simplify: \\(%dx(%dx + %d)\\)", a, b, c),
		Answer:      fmt.Sprintf("%dx^2 + %dx", a*b, a*c),
		Explanation: fmt.Sprintf("\\(%dx(%dx) + %dx(%d) = %dx^{2} + %dx\\)", a, b, a, c, a*b, a*c),
		Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "c": strconv.Itoa(c), "ab": strconv.Itoa(a * b), "ac": strconv.Itoa(a * c)},
	}
}

type polyFoilGen struct{}

func (g *polyFoilGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := rand.Intn(6) + 1
	b := rand.Intn(6) + 1
	c := rand.Intn(6) + 1
	d := rand.Intn(6) + 1
	return generator.Problem{
		Question: fmt.Sprintf("Simplify: \\((%dx + %d)(%dx + %d)\\)", a, b, c, d),
		Answer:   fmt.Sprintf("%dx^2 + %dx + %d", a*c, a*d+b*c, b*d),
		Explanation: fmt.Sprintf("FOIL: \\(%dx^{2} + %dx + %dx + %d = %dx^{2} + %dx + %d\\)",
			a*c, a*d, b*c, b*d, a*c, a*d+b*c, b*d),
		Facts: map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "c": strconv.Itoa(c), "d": strconv.Itoa(d), "ac": strconv.Itoa(a * c), "ad": strconv.Itoa(a * d), "bc": strconv.Itoa(b * c), "bd": strconv.Itoa(b * d), "mid": strconv.Itoa(a*d + b*c)},
	}
}

type polySpecialGen struct{}

func (g *polySpecialGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := rand.Intn(6) + 1
	b := rand.Intn(5) + 1
	return generator.Problem{
		Question: fmt.Sprintf("Simplify: \\((%dx + %d)^{2}\\)", a, b),
		Answer:   fmt.Sprintf("%dx^2 + %dx + %d", a*a, 2*a*b, b*b),
		Explanation: fmt.Sprintf("\\((%dx)^{2} + 2(%dx)(%d) + %d^{2} = %dx^{2} + %dx + %d\\)",
			a, a, b, b, a*a, 2*a*b, b*b),
		Facts: map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "a2": strconv.Itoa(a * a), "mid": strconv.Itoa(2 * a * b), "b2": strconv.Itoa(b * b)},
	}
}

// ---------------------------------------------------------------------------
// Factoring
// ---------------------------------------------------------------------------

type factorGCFGen struct{}

func (g *factorGCFGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	f := rand.Intn(5) + 2
	a := f * (rand.Intn(5) + 2)
	b := f * (rand.Intn(5) + 1)
	return generator.Problem{
		Question:    fmt.Sprintf("Factor: \\(%dx + %d\\)", a, b),
		Answer:      fmt.Sprintf("%d(%dx + %d)", f, a/f, b/f),
		Explanation: fmt.Sprintf("GCF is %d: \\(%dx + %d = %d(%dx + %d)\\)", f, a, b, f, a/f, b/f),
		Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "gcf": strconv.Itoa(f), "out_a": strconv.Itoa(a / f), "out_b": strconv.Itoa(b / f)},
	}
}

type factorTrinomialGen struct{}

func (g *factorTrinomialGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	r1 := rand.Intn(6) + 1
	r2 := rand.Intn(6) + 1
	b := r1 + r2
	c := r1 * r2
	return generator.Problem{
		Question:    fmt.Sprintf("Factor: \\(x^{2} + %dx + %d\\)", b, c),
		Answer:      fmt.Sprintf("(x + %d)(x + %d)", r1, r2),
		Explanation: fmt.Sprintf("Find two numbers that multiply to %d and add to %d: %d and %d.", c, b, r1, r2),
		Facts:       map[string]string{"q": formatQuadratic(b, c), "b": strconv.Itoa(b), "c": strconv.Itoa(c), "r1": strconv.Itoa(r1), "r2": strconv.Itoa(r2)},
	}
}

type factorDiffSquaresGen struct{}

func (g *factorDiffSquaresGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := rand.Intn(8) + 2
	return generator.Problem{
		Question:    fmt.Sprintf("Factor: \\(x^{2} - %d\\)", a*a),
		Answer:      fmt.Sprintf("(x + %d)(x - %d)", a, a),
		Explanation: fmt.Sprintf("\\(x^{2} - %d = (x + %d)(x - %d)\\)", a*a, a, a),
		Facts:       map[string]string{"a2": strconv.Itoa(a * a), "a": strconv.Itoa(a)},
	}
}

type factorACMethodGen struct{}

func (g *factorACMethodGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// The concept is "ax^2 + bx + c by the AC method", so the leading
	// coefficient must not be 1, and the printed trinomial must be the product
	// the answer claims. It used to draw the leading coefficient as a sum
	// offset (b = a + r1*r2) while answering (x + r1)(x + r2), so x^2 + 9x + 1
	// was served with (x+1)(x+1), which multiplies out to x^2 + 2x + 1. Report
	// 196 quoted exactly that pair.
	//
	// (ax + q)(x + s) = ax^2 + (as + q)x + qs, and the two numbers the method
	// looks for are as and q: they multiply to ac = a*qs and add to b.
	a := rand.Intn(4) + 2
	q := rand.Intn(5) + 1
	s := rand.Intn(5) + 1
	b := a*s + q
	c := q * s
	ac := a * c
	return generator.Problem{
		Question: fmt.Sprintf("Factor: \\(%dx^{2} + %dx + %d\\)", a, b, c),
		Answer:   fmt.Sprintf("(%dx + %d)(x + %d)", a, q, s),
		Explanation: fmt.Sprintf("Multiply the leading coefficient by the constant: \\(%d \\times %d = %d\\). "+
			"The two numbers that multiply to %d and add to %d are %d and %d. "+
			"Split the middle term and group: \\(%dx^{2} + %dx + %dx + %d = %dx(x + %d) + %d(x + %d) = (%dx + %d)(x + %d)\\).",
			a, c, ac, ac, b, a*s, q, a, a*s, q, c, a, s, q, s, a, q, s),
		Facts: map[string]string{
			"a":  strconv.Itoa(a),
			"b":  strconv.Itoa(b),
			"c":  strconv.Itoa(c),
			"ac": strconv.Itoa(ac),
			"m":  strconv.Itoa(a * s),
			"n":  strconv.Itoa(q),
			"s":  strconv.Itoa(s),
		},
	}
}

// ---------------------------------------------------------------------------
// Quadratics
// ---------------------------------------------------------------------------

type quadSolveFactorGen struct{}

func (g *quadSolveFactorGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	r1 := rand.Intn(8) - 4
	r2 := rand.Intn(8) - 4
	b := -(r1 + r2)
	c := r1 * r2
	q := formatQuadratic(b, c)
	// A repeated root is one number, so the answer is one number. Written as
	// "-1,-1" it asked the learner to repeat themselves, and "-1" was rejected
	// for having the wrong number of parts (report 219).
	answer := fmt.Sprintf("%d,%d", r1, r2)
	explanation := fmt.Sprintf("Factors: \\((x %+d)(x %+d) = 0\\), so \\(x = %d\\) or \\(x = %d\\)", -r1, -r2, r1, r2)
	if r1 == r2 {
		answer = strconv.Itoa(r1)
		explanation = fmt.Sprintf("Factors: \\((x %+d)^2 = 0\\), so \\(x = %d\\) is the only root", -r1, r1)
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Solve: %s = 0", q),
		Answer:      answer,
		Explanation: explanation,
		Facts:       map[string]string{"q": q, "r1": strconv.Itoa(r1), "r2": strconv.Itoa(r2)},
	}
}

type quadCompleteSquareGen struct{}

func (g *quadCompleteSquareGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	r := rand.Intn(6) + 1
	b := -2 * r
	return generator.Problem{
		Question:    fmt.Sprintf("Complete the square: \\(x^{2} + %dx + \\_\\_)\\) to make a perfect square.", b),
		Answer:      fmt.Sprintf("%d", r*r),
		Explanation: fmt.Sprintf("\\(\\left(\\frac{b}{2}\\right)^{2} = \\left(\\frac{%d}{2}\\right)^{2} = %d\\)", b, r*r),
		Facts:       map[string]string{"b": strconv.Itoa(b), "r": strconv.Itoa(r), "r2": strconv.Itoa(r * r)},
	}
}

type quadFormulaGen struct{}

// quadFormulaGen answers with the two roots. A repeated root is one number, and
// the concept grades numerically rather than as a tuple so that the order the
// learner writes the roots in does not decide the mark: a report on 22 October
// gave "5, 7" for the roots of x^2 - 12x + 35 and was marked wrong.
func (g *quadFormulaGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	r1 := rand.Intn(8) + 1
	r2 := rand.Intn(8) + 1
	b := -(r1 + r2)
	c := r1 * r2
	q := formatQuadratic(b, c)
	answer := fmt.Sprintf("%d,%d", r1, r2)
	explanation := fmt.Sprintf("\\(x = \\frac{%d \\pm \\sqrt{%d - 4(%d)}}{2} = %d\\) or %d", b, b*b, c, r1, r2)
	if r1 == r2 {
		// The discriminant is zero, so both roots are the same number and the
		// answer is a single value rather than a list of two of it.
		answer = strconv.Itoa(r1)
		explanation = fmt.Sprintf("\\(x = \\frac{%d \\pm \\sqrt{0}}{2} = %d\\) twice: the root is repeated.", b, r1)
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Solve using quadratic formula: %s = 0", q),
		Answer:      answer,
		Explanation: explanation,
		Facts:       map[string]string{"q": q, "b": strconv.Itoa(b), "c": strconv.Itoa(c), "b2": strconv.Itoa(b * b), "four_c": strconv.Itoa(4 * c), "disc": strconv.Itoa(b*b - 4*c), "r1": strconv.Itoa(r1), "r2": strconv.Itoa(r2)},
	}
}

type quadDiscriminantGen struct{}

func (g *quadDiscriminantGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	b := rand.Intn(10) - 5
	c := rand.Intn(20) - 10
	disc := b*b - 4*c
	count := "two real"
	if disc == 0 {
		count = "one real"
	} else if disc < 0 {
		count = "zero real"
	}
	return generator.Problem{
		Question:    fmt.Sprintf("What is the discriminant of \\(x^{2} + %dx + %d = 0\\)? (enter a number)", b, c),
		Answer:      strconv.Itoa(disc),
		Explanation: fmt.Sprintf("The discriminant is \\(b^{2} - 4ac = %d^{2} - 4(%d) = %d\\), so the equation has %s solutions.", b, c, disc, count),
		Facts: map[string]string{
			"b": strconv.Itoa(b), "c": strconv.Itoa(c), "b2": strconv.Itoa(b * b),
			"four_c": strconv.Itoa(4 * c), "disc": strconv.Itoa(disc), "count": count,
		},
	}
}

// ---------------------------------------------------------------------------
// Functions
// ---------------------------------------------------------------------------

type funcConceptGen struct{}

func (g *funcConceptGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := rand.Intn(5) + 2
	b := rand.Intn(5) + 1
	x := rand.Intn(5) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("If \\(f(x) = %dx + %d\\), what is \\(f(%d)\\)?", a, b, x),
		Answer:      fmt.Sprintf("%d", a*x+b),
		Explanation: fmt.Sprintf("\\(f(%d) = %d(%d) + %d = %d\\)", x, a, x, b, a*x+b),
		Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "x": strconv.Itoa(x), "ax": strconv.Itoa(a * x), "y": strconv.Itoa(a*x + b)},
	}
}

type funcNotationGen struct{}

func (g *funcNotationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := rand.Intn(5) + 2
	x := rand.Intn(5) + 1
	b := rand.Intn(5) + 1
	y := a*x + b
	return generator.Problem{
		Question:    fmt.Sprintf("\\(f(x) = %dx + %d\\). Find \\(x\\) when \\(f(x) = %d\\).", a, b, y),
		Answer:      fmt.Sprintf("%d", x),
		Explanation: fmt.Sprintf("\\(%dx + %d = %d\\) → \\(%dx = %d\\) → \\(x = %d\\)", a, b, y, a, y-b, x),
		Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "y": strconv.Itoa(y), "rhs": strconv.Itoa(y - b), "x": strconv.Itoa(x)},
	}
}

type funcEvaluateGen struct{}

func (g *funcEvaluateGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := rand.Intn(5) + 2
	x := rand.Intn(5) + 2
	b := rand.Intn(5) + 1
	c := a*x + b
	return generator.Problem{
		Question:    fmt.Sprintf("\\(f(x) = %dx + %d\\). Find \\(f(%d)\\).", a, b, x),
		Answer:      fmt.Sprintf("%d", c),
		Explanation: fmt.Sprintf("\\(f(%d) = %d(%d) + %d = %d\\)", x, a, x, b, c),
		Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "x": strconv.Itoa(x), "ax": strconv.Itoa(a * x), "c": strconv.Itoa(c)},
	}
}

type funcLinearGen struct{}

func (g *funcLinearGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	m := rand.Intn(6) + 1
	if rand.Intn(2) == 0 {
		m = -m
	}
	b := rand.Intn(10) - 5
	switch rand.Intn(3) {
	case 0:
		return generator.Problem{
			Question:    fmt.Sprintf("What is the slope of \\(%s\\)? (enter a number)", formatLinear(m, b)),
			Answer:      strconv.Itoa(m),
			Explanation: fmt.Sprintf("In \\(y = mx + b\\), the slope is %d.", m),
			Facts:       map[string]string{"m": strconv.Itoa(m), "b": strconv.Itoa(b), "is_slope": "yes"},
		}
	case 1:
		return generator.Problem{
			Question:    fmt.Sprintf("What is the y-intercept of \\(%s\\)? (enter a number)", formatLinear(m, b)),
			Answer:      strconv.Itoa(b),
			Explanation: fmt.Sprintf("In \\(y = mx + b\\), the y-intercept is %d.", b),
			Facts:       map[string]string{"m": strconv.Itoa(m), "b": strconv.Itoa(b), "is_intercept": "yes"},
		}
	default:
		k := rand.Intn(6) + 1
		return generator.Problem{
			Question:    fmt.Sprintf("If \\(f(x) = %s\\), what is \\(f(%d)\\)? (enter a number)", formatLinear(m, b), k),
			Answer:      strconv.Itoa(m*k + b),
			Explanation: fmt.Sprintf("Substitute \\(x = %d\\): \\(%d(%d) + (%d) = %d\\).", k, m, k, b, m*k+b),
			Facts:       map[string]string{"m": strconv.Itoa(m), "b": strconv.Itoa(b), "x": strconv.Itoa(k), "value": strconv.Itoa(m*k + b)},
		}
	}
}

type funcQuadGen struct{}

func (g *funcQuadGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := rand.Intn(4) + 1
	b := rand.Intn(5) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("\\(f(x) = %dx^{2} + %d\\). Find \\(f(2)\\).", a, b),
		Answer:      fmt.Sprintf("%d", 4*a+b),
		Explanation: fmt.Sprintf("\\(f(2) = %d(4) + %d = %d + %d = %d\\)", a, b, 4*a, b, 4*a+b),
		Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "four_a": strconv.Itoa(4 * a), "y": strconv.Itoa(4*a + b)},
	}
}

// ---------------------------------------------------------------------------
// Exponential & Logarithmic
// ---------------------------------------------------------------------------

type algExpConceptGen struct{}

func (g *algExpConceptGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := rand.Intn(8) + 2
	x := rand.Intn(6) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("If \\(f(x) = %d^{x}\\), find \\(f(%d)\\).", a, x),
		Answer:      fmt.Sprintf("%d", mathutil.IntPow(a, x)),
		Explanation: fmt.Sprintf("\\(%d^{%d} = %d\\)", a, x, mathutil.IntPow(a, x)),
		Facts:       map[string]string{"base": strconv.Itoa(a), "x": strconv.Itoa(x), "val": strconv.Itoa(mathutil.IntPow(a, x))},
	}
}

type algExpEvaluateGen struct{}

func (g *algExpEvaluateGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := rand.Intn(8) + 3
	x := rand.Intn(6) + 2
	return generator.Problem{
		Question:    fmt.Sprintf("Evaluate: \\(%d^{%d}\\)", a, x),
		Answer:      fmt.Sprintf("%d", mathutil.IntPow(a, x)),
		Explanation: fmt.Sprintf("\\(%d^{%d} = %d\\)", a, x, mathutil.IntPow(a, x)),
		Facts:       map[string]string{"base": strconv.Itoa(a), "x": strconv.Itoa(x), "val": strconv.Itoa(mathutil.IntPow(a, x))},
	}
}

type logConceptGen struct{}

// The concept has three knowledge points: the definition, the basic
// identities, and the logarithmic function. The generator used to return the
// conversion prompt on every rung, so a learner who reached the identities or
// the function still saw only "write as a logarithm".
func (g *logConceptGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch ctx.KPOrdinal {
	case 2:
		return logIdentityProblem()
	case 3:
		return logarithmicFunctionProblem()
	default:
		return logDefinitionProblem()
	}
}

func logDefinitionProblem() generator.Problem {
	base := rand.Intn(3) + 2
	exp := rand.Intn(4) + 1
	val := mathutil.IntPow(base, exp)
	return generator.Problem{
		Question:    fmt.Sprintf("If \\(%d^{%d} = %d\\), what is \\(\\log_{%d}(%d)\\)? (enter a number)", base, exp, val, base, val),
		Answer:      strconv.Itoa(exp),
		Explanation: fmt.Sprintf("A logarithm answers the exponent: \\(%d^{%d} = %d\\), so \\(\\log_{%d}(%d) = %d\\).", base, exp, val, base, val, exp),
		Facts:       map[string]string{"base": strconv.Itoa(base), "exp": strconv.Itoa(exp), "val": strconv.Itoa(val)},
	}
}

func logIdentityProblem() generator.Problem {
	base := rand.Intn(8) + 2
	switch rand.Intn(3) {
	case 0:
		return generator.Problem{
			Question:    fmt.Sprintf("What is \\(\\log_{%d}(1)\\)? (enter a number)", base),
			Answer:      "0",
			Explanation: fmt.Sprintf("\\(%d^{0} = 1\\), so the logarithm is 0.", base),
			Facts:       map[string]string{"base": strconv.Itoa(base), "value": "0"},
		}
	case 1:
		return generator.Problem{
			Question:    fmt.Sprintf("What is \\(\\log_{%d}(%d)\\)? (enter a number)", base, base),
			Answer:      "1",
			Explanation: fmt.Sprintf("\\(%d^{1} = %d\\), so the logarithm is 1.", base, base),
			Facts:       map[string]string{"base": strconv.Itoa(base), "value": "1"},
		}
	default:
		value := rand.Intn(8) + 2
		for value == base {
			value = rand.Intn(8) + 2
		}
		return generator.Problem{
			Question:    fmt.Sprintf("Simplify \\(%d^{\\log_{%d}(%d)}\\). (enter a number)", base, base, value),
			Answer:      strconv.Itoa(value),
			Explanation: fmt.Sprintf("The exponential and logarithm are inverses, so \\(%d^{\\log_{%d}(%d)} = %d\\).", base, base, value, value),
			Facts:       map[string]string{"base": strconv.Itoa(base), "value": strconv.Itoa(value)},
		}
	}
}

func logarithmicFunctionProblem() generator.Problem {
	base := rand.Intn(8) + 2
	k := rand.Intn(4) + 1
	x := mathutil.IntPow(base, k)
	return generator.Problem{
		Question:    fmt.Sprintf("Solve \\(\\log_{%d}(x) = %d\\) for \\(x\\). (enter a number)", base, k),
		Answer:      strconv.Itoa(x),
		Explanation: fmt.Sprintf("Rewrite the logarithm as an exponential: \\(x = %d^{%d} = %d\\).", base, k, x),
		Facts:       map[string]string{"base": strconv.Itoa(base), "exponent": strconv.Itoa(k), "x": strconv.Itoa(x)},
	}
}

type logEvaluateGen struct{}

func (g *logEvaluateGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	base := rand.Intn(9) + 2
	exp := rand.Intn(6) + 1
	val := mathutil.IntPow(base, exp)
	return generator.Problem{
		Question:    fmt.Sprintf("Evaluate: \\(\\log_{%d}(%d)\\)", base, val),
		Answer:      fmt.Sprintf("%d", exp),
		Explanation: fmt.Sprintf("\\(%d^{%d} = %d\\), so \\(\\log_{%d}(%d) = %d\\)", base, exp, val, base, val, exp),
		Facts:       map[string]string{"base": strconv.Itoa(base), "exp": strconv.Itoa(exp), "val": strconv.Itoa(val)},
	}
}

type logPropertiesGen struct{}

func (g *logPropertiesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	base := rand.Intn(3) + 2
	a := mathutil.IntPow(base, rand.Intn(3)+1)
	b := mathutil.IntPow(base, rand.Intn(3)+2)
	return generator.Problem{
		Question:    fmt.Sprintf("Simplify: \\(\\log_{%d}(%d \\times %d)\\)", base, a, b),
		Answer:      fmt.Sprintf("log_%d(%d) + log_%d(%d)", base, a, base, b),
		Explanation: fmt.Sprintf("log_b(xy) = log_b(x) + log_b(y)."),
		Facts:       map[string]string{"base": strconv.Itoa(base), "a": strconv.Itoa(a), "b": strconv.Itoa(b)},
	}
}

// ---------------------------------------------------------------------------
// Sequences
// ---------------------------------------------------------------------------

type seqArithGen struct{}

func (g *seqArithGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a1 := rand.Intn(10) + 1
	d := rand.Intn(5) + 2
	n := rand.Intn(5) + 3
	an := a1 + (n-1)*d
	return generator.Problem{
		Question:    fmt.Sprintf("Arithmetic sequence: \\(%d, %d, %d, \\ldots\\) Find term %d.", a1, a1+d, a1+2*d, n),
		Answer:      fmt.Sprintf("%d", an),
		Explanation: fmt.Sprintf("\\(a_{n} = a_{1} + (n-1)d = %d + (%d)(%d) = %d\\)", a1, n-1, d, an),
		Facts:       map[string]string{"a1": strconv.Itoa(a1), "d": strconv.Itoa(d), "n": strconv.Itoa(n), "nm1": strconv.Itoa(n - 1), "an": strconv.Itoa(an)},
	}
}

type seqGeomGen struct{}

func (g *seqGeomGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a1 := rand.Intn(5) + 2
	r := rand.Intn(3) + 2
	n := rand.Intn(4) + 2
	an := a1 * mathutil.IntPow(r, n-1)
	return generator.Problem{
		Question:    fmt.Sprintf("Geometric sequence: \\(%d, %d, %d, \\ldots\\) Find term %d.", a1, a1*r, a1*r*r, n),
		Answer:      fmt.Sprintf("%d", an),
		Explanation: fmt.Sprintf("\\(a_{n} = a_{1} \\times r^{n-1} = %d \\times %d^{%d} = %d\\)", a1, r, n-1, an),
		Facts:       map[string]string{"a1": strconv.Itoa(a1), "r": strconv.Itoa(r), "n": strconv.Itoa(n), "nm1": strconv.Itoa(n - 1), "an": strconv.Itoa(an)},
	}
}

type seqSumArithGen struct{}

func (g *seqSumArithGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a1 := rand.Intn(10) + 1
	d := rand.Intn(5) + 2
	n := rand.Intn(6) + 2
	an := a1 + (n-1)*d
	sum := n * (a1 + an) / 2
	return generator.Problem{
		Question:    fmt.Sprintf("Find sum of arithmetic series: \\(%d + %d + %d + \\cdots\\) (first %d terms).", a1, a1+d, a1+2*d, n),
		Answer:      fmt.Sprintf("%d", sum),
		Explanation: fmt.Sprintf("\\(S_{n} = \\frac{n(a_{1}+a_{n})}{2} = \\frac{%d(%d+%d)}{2} = %d\\)", n, a1, an, sum),
		Facts:       map[string]string{"n": strconv.Itoa(n), "a1": strconv.Itoa(a1), "an": strconv.Itoa(an), "ends": strconv.Itoa(a1 + an), "sum": strconv.Itoa(sum)},
	}
}

type seqSumGeoGen struct{}

func (g *seqSumGeoGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a1 := rand.Intn(5) + 2
	r := rand.Intn(3) + 2
	// Sum of first n terms: a1(1-r^n)/(1-r)
	// Keep values small
	n := rand.Intn(4) + 2
	rn := mathutil.IntPow(r, n)
	sum := a1 * (1 - rn) / (1 - r)
	return generator.Problem{
		Question:    fmt.Sprintf("Find sum of geometric series: \\(%d + %d + %d + \\cdots\\) (first %d terms, \\(r=%d\\)).", a1, a1*r, a1*r*r, n, r),
		Answer:      fmt.Sprintf("%d", sum),
		Explanation: fmt.Sprintf("\\(S_{n} = \\frac{a_{1}(1-r^{n})}{1-r} = \\frac{%d(1-%d)}{1-%d} = %d\\)", a1, rn, r, sum),
		Facts:       map[string]string{"a1": strconv.Itoa(a1), "r": strconv.Itoa(r), "n": strconv.Itoa(n), "rn": strconv.Itoa(rn), "sum": strconv.Itoa(sum)},
	}
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

type frac struct{ num, den int }

func reduce(num, den int) frac {
	if num == 0 {
		return frac{0, 1}
	}
	g := mathutil.GCD(mathutil.Abs(num), den)
	n := num / g
	d := den / g
	if d < 0 {
		n, d = -n, -d
	}
	return frac{n, d}
}

func fracOrInt(num, den int) string {
	f := reduce(num, den)
	if f.den == 1 {
		return fmt.Sprintf("%d", f.num)
	}
	return fmt.Sprintf("%d/%d", f.num, f.den)
}

func formatLinear(m, b int) string {
	if b == 0 {
		return fmt.Sprintf("y = %dx", m)
	}
	if b < 0 {
		return fmt.Sprintf("y = %dx - %d", m, -b)
	}
	return fmt.Sprintf("y = %dx + %d", m, b)
}

// formatLinearExpr formats "ax + b" with proper signs, e.g. "3x - 5"
func formatLinearExpr(a, b int) string {
	if b == 0 {
		return fmt.Sprintf("%dx", a)
	}
	if b < 0 {
		return fmt.Sprintf("%dx - %d", a, -b)
	}
	return fmt.Sprintf("%dx + %d", a, b)
}

func formatQuadratic(b, c int) string {
	var parts []string
	parts = append(parts, "x^2")
	if b != 0 {
		if b == 1 {
			parts = append(parts, "+ x")
		} else if b == -1 {
			parts = append(parts, "- x")
		} else if b > 0 {
			parts = append(parts, fmt.Sprintf("+ %dx", b))
		} else {
			parts = append(parts, fmt.Sprintf("- %dx", -b))
		}
	}
	if c != 0 {
		if c > 0 {
			parts = append(parts, fmt.Sprintf("+ %d", c))
		} else {
			parts = append(parts, fmt.Sprintf("- %d", -c))
		}
	}
	return strings.Join(parts, " ")
}

type ineqTwoVarGen struct{}

func (g *ineqTwoVarGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	m := rand.Intn(4) + 1
	b := rand.Intn(6) - 3
	op := "<"
	if rand.Intn(2) == 0 {
		op = ">"
	}
	testX := rand.Intn(5) + 1
	boundary := m*testX + b
	answer := boundary - 1
	adjective := "largest"
	if op == ">" {
		answer = boundary + 1
		adjective = "smallest"
	}
	return generator.Problem{
		Question:    fmt.Sprintf("For \\(y %s %dx + %d\\) with \\(x = %d\\), what is the %s integer value of \\(y\\) that satisfies the inequality? (enter a number)", op, m, b, testX, adjective),
		Answer:      strconv.Itoa(answer),
		Explanation: fmt.Sprintf("At \\(x = %d\\) the boundary is \\(y = %d\\), so the %s integer satisfying \\(y %s %d\\) is %d.", testX, boundary, adjective, op, boundary, answer),
		Facts:       map[string]string{"x": strconv.Itoa(testX), "m": strconv.Itoa(m), "b": strconv.Itoa(b), "op": op, "boundary": strconv.Itoa(boundary), "answer": strconv.Itoa(answer)},
	}
}

type conicCircleGen struct{}

func (g *conicCircleGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	h := rand.Intn(5) - 2
	k := rand.Intn(5) - 2
	r := rand.Intn(4) + 2
	r2 := r * r
	return generator.Problem{
		Question:    fmt.Sprintf("What is the radius of \\((x%+d)^{2} + (y%+d)^{2} = %d\\)?", -h, -k, r2),
		Answer:      fmt.Sprintf("%d", r),
		Explanation: fmt.Sprintf("The radius is \\(\\sqrt{%d} = %d\\)", r2, r),
		Facts:       map[string]string{"r2": strconv.Itoa(r2), "r": strconv.Itoa(r)},
	}
}

type conicEllipseGen struct{}

func (g *conicEllipseGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	h := rand.Intn(5) - 2
	k := rand.Intn(5) - 2
	a := rand.Intn(3) + 2
	b := rand.Intn(2) + 1
	for a <= b {
		b = rand.Intn(2) + 1
	}
	a2 := a * a
	b2 := b * b
	askMajor := rand.Intn(2) == 0
	label := "semi-minor"
	val := b
	if askMajor {
		label = "semi-major"
		val = a
	}
	return generator.Problem{
		Question:    fmt.Sprintf("What is the %s axis length of \\(\\frac{(x%+d)^{2}}{%d} + \\frac{(y%+d)^{2}}{%d} = 1\\)?", label, -h, a2, -k, b2),
		Answer:      fmt.Sprintf("%d", val),
		Explanation: fmt.Sprintf("The %s axis length is %d because \\(%s^{2} = %d\\)", label, val, label[:6], val*val),
		Facts:       map[string]string{"label": label, "val": strconv.Itoa(val), "val2": strconv.Itoa(val * val)},
	}
}

// ============= conic =============

type conicParabolaGen struct{}

func (g *conicParabolaGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	h := rand.Intn(5) - 2
	k := rand.Intn(5) - 2
	p := rand.Intn(3) + 1
	if rand.Intn(2) == 0 {
		p = -p
	}
	if rand.Intn(2) == 0 {
		return generator.Problem{
			Question:    fmt.Sprintf("Parabola: \\((x%+d)^{2} = %d(y%+d)\\). Which direction does it open?", -h, 4*p, -k),
			Answer:      map[bool]string{true: "up", false: "down"}[p > 0],
			Explanation: fmt.Sprintf("If 4p > 0, opens up; if 4p < 0, opens down. Here 4p = %d.", 4*p),
			Facts:       map[string]string{"p": strconv.Itoa(p), "four_p": strconv.Itoa(4 * p), "is_vertical_axis": "yes"},
		}
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Parabola: \\((y%+d)^{2} = %d(x%+d)\\). Which direction does it open?", -k, 4*p, -h),
		Answer:      map[bool]string{true: "right", false: "left"}[p > 0],
		Explanation: fmt.Sprintf("If 4p > 0, opens right; if 4p < 0, opens left. Here 4p = %d.", 4*p),
		Facts:       map[string]string{"p": strconv.Itoa(p), "four_p": strconv.Itoa(4 * p), "is_horizontal_axis": "yes"},
	}
}

type conicHyperbolaGen struct{}

func (g *conicHyperbolaGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := rand.Intn(3) + 2
	b := rand.Intn(2) + 1
	h := rand.Intn(3) - 1
	k := rand.Intn(3) - 1
	if rand.Intn(2) == 0 {
		return generator.Problem{
			Question:    fmt.Sprintf("Hyperbola: \\(\\frac{(x%+d)^{2}}{%d} - \\frac{(y%+d)^{2}}{%d} = 1\\). Which axis is transverse?", -h, a*a, -k, b*b),
			Answer:      "horizontal (x-axis)",
			Explanation: fmt.Sprintf("x term positive → horizontal transverse axis."),
			Facts:       map[string]string{"is_horizontal": "yes"},
		}
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Hyperbola: \\(\\frac{(y%+d)^{2}}{%d} - \\frac{(x%+d)^{2}}{%d} = 1\\). Which axis is transverse?", -k, a*a, -h, b*b),
		Answer:      "vertical (y-axis)",
		Explanation: fmt.Sprintf("y term positive → vertical transverse axis."),
		Facts:       map[string]string{"is_vertical": "yes"},
	}
}

// ============= equations =============

type eqAbsValGen struct{}

// absValueTerm writes the inside of the absolute value without the
// "x + -1" that a bare %+d-style concatenation produces.
func absValueTerm(b int) string {
	switch {
	case b > 0:
		return fmt.Sprintf("x + %d", b)
	case b < 0:
		return fmt.Sprintf("x - %d", -b)
	default:
		return "x"
	}
}

func (g *eqAbsValGen) AnswerFormat() generator.AnswerFormat {
	return generator.AnswerFormat{Hint: "Both solutions separated by a comma, e.g. -3, 5", InputMode: "text"}
}

func (g *eqAbsValGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := rand.Intn(8) + 1
	b := rand.Intn(10) - 5
	term := absValueTerm(b)
	return generator.Problem{
		Question:    fmt.Sprintf("Solve: \\(|%s| = %d\\)", term, a),
		Answer:      fmt.Sprintf("%d,%d", a-b, -a-b),
		Explanation: fmt.Sprintf("\\(%s = %d\\) or \\(%s = -%d\\) → \\(x = %d\\) or \\(x = %d\\)", term, a, term, a, a-b, -a-b),
		Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "r1": strconv.Itoa(a - b), "r2": strconv.Itoa(-a - b)},
	}
}

type eqBinomialGen struct{}

func (g *eqBinomialGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	r := rand.Intn(6) + 1
	k := rand.Intn(4) + 2
	rhs := mathutil.IntPow(k, r)
	return generator.Problem{
		Question:    fmt.Sprintf("Solve: \\(x^{%d} = %d\\)", r, rhs),
		Answer:      fmt.Sprintf("%d", k),
		Explanation: fmt.Sprintf("\\(x^{%d} = %d\\) → \\(x = %d^{\\frac{1}{%d}} = %d\\)", r, rhs, rhs, r, k),
		Facts:       map[string]string{"r": strconv.Itoa(r), "rhs": strconv.Itoa(rhs), "k": strconv.Itoa(k)},
	}
}

type eqExpGen struct{}

func (g *eqExpGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	b := rand.Intn(4) + 2
	e := rand.Intn(4) + 1
	p := rand.Intn(3) + 1
	rhs := mathutil.IntPow(b, e+p)
	return generator.Problem{
		Question:    fmt.Sprintf("Solve: \\(%d^{x-%d} = %d\\)", b, p, rhs),
		Answer:      fmt.Sprintf("%d", e+p),
		Explanation: fmt.Sprintf("\\(%d^{x-%d} = %d^{%d}\\) → \\(x-%d = %d\\) → \\(x = %d\\)", b, p, b, e+p, p, e, e+p),
		Facts:       map[string]string{"base": strconv.Itoa(b), "p": strconv.Itoa(p), "rhs": strconv.Itoa(rhs), "e": strconv.Itoa(e), "x": strconv.Itoa(e + p)},
	}
}

type extraneousRootsGen struct{}

func (g *extraneousRootsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"Can squaring both sides of an equation introduce extraneous roots? (yes/no)", "yes", "Squaring is not reversible: if a=b then a²=b², but b²=a² does not imply a=b (e.g. a=2, b=-2). So squaring can introduce extraneous roots, always check candidates in the original equation."},
		{"Solve \\(\\sqrt{x} = -2\\). How many solutions? (enter a number)", "0", "The principal square root is always \\(\\geq 0\\), so \\(\\sqrt{x} = -2\\) has no solution."},
		{"The candidates when solving \\(\\sqrt{x+3} = x-3\\) are \\(x = 1\\) and \\(x = 6\\). Which one satisfies the original equation? (enter a number)", "6", "Check: \\(\\sqrt{1+3} = 2\\) but \\(1-3 = -2\\), so \\(x=1\\) is extraneous; \\(\\sqrt{6+3} = 3\\) and \\(6-3 = 3\\), so \\(x=6\\) works. Extraneous roots satisfy the squared equation but not the original."},
		{"If you multiply both sides of an equation by \\((x - 2)\\), should you check whether \\(x = 2\\) is extraneous? (yes/no)", "yes", "If you multiply by something that equals 0 for some x, you may introduce extraneous solutions, always substitute candidates back into the original equation."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question: e.question, Answer: e.answer, Explanation: e.exp,
	}
}

type eqIrrationalGen struct{}

func (g *eqIrrationalGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	x := rand.Intn(8) + 2
	k := x * x
	return generator.Problem{
		Question:    fmt.Sprintf("Solve: \\(\\sqrt{x} = %d\\)", x),
		Answer:      fmt.Sprintf("%d", k),
		Explanation: fmt.Sprintf("Square both sides: \\(x = %d^{2} = %d\\)", x, k),
		Facts:       map[string]string{"x": strconv.Itoa(x), "k": strconv.Itoa(k)},
	}
}

type eqLogGen struct{}

func (g *eqLogGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	b := rand.Intn(3) + 2
	e := rand.Intn(3) + 1
	v := mathutil.IntPow(b, e)
	return generator.Problem{
		Question:    fmt.Sprintf("Solve: \\(\\log_{%d}(x) = %d\\)", b, e),
		Answer:      fmt.Sprintf("%d", v),
		Explanation: fmt.Sprintf("\\(\\log_{%d}(x) = %d\\) → \\(x = %d^{%d} = %d\\)", b, e, b, e, v),
		Facts:       map[string]string{"base": strconv.Itoa(b), "e": strconv.Itoa(e), "v": strconv.Itoa(v)},
	}
}

type eqPolyGen struct{}

func (g *eqPolyGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	r1 := rand.Intn(6) - 3
	r2 := rand.Intn(6) - 3
	// (x - r1)(x - r2) = 0
	b := -(r1 + r2)
	c := r1 * r2
	q := formatQuadratic(b, c)
	return generator.Problem{
		Question:    fmt.Sprintf("Solve polynomial: %s = 0", q),
		Answer:      fmt.Sprintf("%d,%d", r1, r2),
		Explanation: fmt.Sprintf("\\((x %+d)(x %+d) = 0\\) → \\(x = %d\\) or \\(x = %d\\)", -r1, -r2, r1, r2),
		Facts:       map[string]string{"q": q, "r1": strconv.Itoa(r1), "r2": strconv.Itoa(r2)},
	}
}

type eqRationalGen struct{}

func (g *eqRationalGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	x := rand.Intn(8) + 2
	a := rand.Intn(5) + 1
	b := rand.Intn(5) + 1
	rhs := a * x / b
	return generator.Problem{
		Question:    fmt.Sprintf("Solve: \\(%dx/%d = %d\\)", a, b, rhs),
		Answer:      fmt.Sprintf("%d", x),
		Explanation: fmt.Sprintf("Multiply both sides by %d: \\(%dx = %d\\) → \\(x = %d\\)", b, a, rhs*b, x),
		Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "rhs": strconv.Itoa(rhs), "ax": strconv.Itoa(rhs * b), "x": strconv.Itoa(x)},
	}
}

type eqTrinomialGen struct{}

func (g *eqTrinomialGen) AnswerFormat() generator.AnswerFormat {
	return generator.AnswerFormat{
		Hint:      "The value of x, e.g. -6",
		InputMode: "numeric",
	}
}

func (g *eqTrinomialGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	r := rand.Intn(6) + 2
	// x^2 + 2r x + r^2 = 0
	return generator.Problem{
		Question:    fmt.Sprintf("Solve for \\(x\\): \\(x^{2} + %dx + %d = 0\\). Enter the value of \\(x\\).", 2*r, r*r),
		Answer:      fmt.Sprintf("%d", -r),
		Explanation: fmt.Sprintf("\\((x + %d)^{2} = 0\\) → \\(x = %d\\)", r, -r),
		Facts:       map[string]string{"r": strconv.Itoa(r), "coef": strconv.Itoa(2 * r), "const": strconv.Itoa(r * r), "x": strconv.Itoa(-r)},
	}
}

// ============= functions =============

type funcAbsValGen struct{}

func (g *funcAbsValGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// The authored walkthrough fixes the shape as f(x) = |x| evaluated at x, so
	// the widening is the range of x rather than a second family of functions.
	x := rand.Intn(81) - 40
	return generator.Problem{
		Question:    fmt.Sprintf("If \\(f(x) = |x|\\), what is \\(f(%d)\\)?", x),
		Answer:      fmt.Sprintf("%d", mathutil.Abs(x)),
		Explanation: fmt.Sprintf("\\(|%d| = %d\\)", x, mathutil.Abs(x)),
		Facts:       map[string]string{"x": strconv.Itoa(x), "abs": strconv.Itoa(mathutil.Abs(x))},
	}
}

type funcCompositeGen struct{}

func (g *funcCompositeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	x := rand.Intn(5) + 1
	a := rand.Intn(4) + 1
	b := rand.Intn(4) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("If \\(f(x) = %dx\\) and \\(g(x) = x + %d\\), what is \\(f(g(%d))\\)?", a, b, x),
		Answer:      fmt.Sprintf("%d", a*(x+b)),
		Explanation: fmt.Sprintf("\\(g(%d) = %d+%d = %d\\). \\(f(g(%d)) = f(%d) = %d \\times %d = %d\\)", x, x, b, x+b, x, x+b, a, x+b, a*(x+b)),
		Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "x": strconv.Itoa(x), "gx": strconv.Itoa(x + b), "y": strconv.Itoa(a * (x + b))},
	}
}

type funcDirichletGen struct{}

func (g *funcDirichletGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch rand.Intn(4) {
	case 0:
		// Rational input evaluations
		type ratSpec struct {
			val string
		}
		rats := []ratSpec{
			{"1/2"}, {"-1/2"}, {"2/3"}, {"-3/4"}, {"5/7"}, {"-7/9"},
			{"4/5"}, {"-1/3"}, {"8/5"}, {"-5/6"}, {"22/7"}, {"0"},
			{"1"}, {"-4"}, {"12"}, {"3/8"}, {"-9/2"}, {"7/10"},
		}
		r := rats[rand.Intn(len(rats))]
		return generator.Problem{
			Question:    fmt.Sprintf("For the Dirichlet function D(x) (1 on rationals, 0 on irrationals), what is D(%s)? (enter a number)", r.val),
			Answer:      "1",
			Explanation: fmt.Sprintf("%s is a rational number, so by definition D(%s) = 1.", r.val, r.val),
			Facts:       map[string]string{"input": r.val, "kind": "rational", "answer": "1"},
		}
	case 1:
		// Irrational input evaluations
		type irratSpec struct {
			val string
		}
		irrats := []irratSpec{
			{"\\sqrt{2}"}, {"\\sqrt{3}"}, {"\\sqrt{5}"}, {"\\pi"}, {"e"},
			{"\\sqrt{7}"}, {"\\pi/2"}, {"\\pi/4"}, {"\\sqrt{11}"}, {"\\ln(2)"},
			{"1 + \\sqrt{2}"}, {"\\sqrt{10}"}, {"2\\pi"}, {"\\sqrt{6}"},
		}
		ir := irrats[rand.Intn(len(irrats))]
		return generator.Problem{
			Question:    fmt.Sprintf("For the Dirichlet function D(x) (1 on rationals, 0 on irrationals), what is D(%s)? (enter a number)", ir.val),
			Answer:      "0",
			Explanation: fmt.Sprintf("%s is an irrational number, so by definition D(%s) = 0.", ir.val, ir.val),
			Facts:       map[string]string{"input": ir.val, "kind": "irrational", "answer": "0"},
		}
	case 2:
		// Scaled combination a * D(x) + b
		a := rand.Intn(6) + 2
		b := rand.Intn(6) + 1
		isRat := rand.Intn(2) == 0
		if isRat {
			rats := []string{"3/5", "-2/7", "5/2", "1/4", "-4/3", "7/8"}
			val := rats[rand.Intn(len(rats))]
			ans := strconv.Itoa(a + b)
			return generator.Problem{
				Question:    fmt.Sprintf("Let D(x) be the Dirichlet function. What is %d D(%s) + %d? (enter a number)", a, val, b),
				Answer:      ans,
				Explanation: fmt.Sprintf("Since %s is rational, D(%s) = 1. Thus %d(1) + %d = %s.", val, val, a, b, ans),
				Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "input": val, "answer": ans},
			}
		}
		irrats := []string{"\\sqrt{3}", "\\pi", "e", "\\sqrt{5}", "\\pi/3"}
		val := irrats[rand.Intn(len(irrats))]
		ans := strconv.Itoa(b)
		return generator.Problem{
			Question:    fmt.Sprintf("Let D(x) be the Dirichlet function. What is %d D(%s) + %d? (enter a number)", a, val, b),
			Answer:      ans,
			Explanation: fmt.Sprintf("Since %s is irrational, D(%s) = 0. Thus %d(0) + %d = %s.", val, val, a, b, ans),
			Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "input": val, "answer": ans},
		}
	default:
		// Conceptual properties
		type propSpec struct {
			q   string
			ans string
			exp string
		}
		props := []propSpec{
			{"Is the Dirichlet function continuous anywhere? (yes/no)", "no", "The Dirichlet function is discontinuous at every real point because both rationals and irrationals are dense."},
			{"Is the Dirichlet function periodic? (yes/no)", "yes", "The Dirichlet function is periodic with every rational number as a period."},
			{"Does the Dirichlet function have a fundamental (smallest positive) period? (yes/no)", "no", "There is no smallest positive rational number, so no fundamental period exists."},
			{"Is the Dirichlet function Riemann integrable on [0, 1]? (yes/no)", "no", "Upper Darboux sums equal 1 and lower sums equal 0, so it is not Riemann integrable."},
			{"Is the Dirichlet function Lebesgue integrable on [0, 1]? (yes/no)", "yes", "The rationals have Lebesgue measure 0, so the function equals 0 almost everywhere and is Lebesgue integrable."},
			{"What is the Lebesgue integral of the Dirichlet function over [0, 1]? (enter a number)", "0", "Since rationals have measure 0, the Lebesgue integral over [0, 1] is exactly 0."},
			{"Is the Dirichlet function an even function satisfying D(-x) = D(x)? (yes/no)", "yes", "-x is rational if and only if x is rational, so D(-x) = D(x)."},
			{"Is the Dirichlet function bounded on the real line? (yes/no)", "yes", "The range of D(x) is {0, 1}, which is bounded."},
		}
		p := props[rand.Intn(len(props))]
		return generator.Problem{
			Question:    p.q,
			Answer:      p.ans,
			Explanation: p.exp,
			Facts:       map[string]string{"question": p.q, "answer": p.ans},
		}
	}
}

type funcDomainGen struct{}

func (g *funcDomainGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"What is excluded from the domain of \\(f(x) = 1/x\\)? (enter a number)", "0", "Division by zero is undefined, so \\(x\\) cannot be 0. The domain is all nonzero reals."},
		{"What is the smallest integer in the domain of \\(f(x) = \\sqrt{x}\\)? (enter a number)", "0", "The square root of a negative number is not real, so the domain is \\(x \\geq 0\\); the smallest integer in it is 0."},
		{"What is the smallest integer in the domain of \\(f(x) = \\ln(x)\\)? (enter a number)", "1", "The natural log of zero or negative numbers is undefined, so the domain is \\(x > 0\\); the smallest integer in it is 1."},
		{"What is excluded from the domain of \\(f(x) = 1/(x-2)\\)? (enter a number)", "2", "The denominator is zero at \\(x = 2\\), so 2 is excluded from the domain."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question: e.question, Answer: e.answer, Explanation: e.exp,
	}
}

type funcEvenOddGen struct{}

func (g *funcEvenOddGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"Is \\(f(x) = x^{2}\\) even or odd?", "even", "\\(f(-x) = (-x)^{2} = x^{2} = f(x)\\), so it is even."},
		{"Is \\(f(x) = x^{3}\\) even or odd?", "odd", "\\(f(-x) = (-x)^{3} = -x^{3} = -f(x)\\), so it is odd."},
		{"Is \\(f(x) = \\sin(x)\\) even or odd?", "odd", "\\(\\sin(-x) = -\\sin(x)\\), so sine is odd."},
		{"Is \\(f(x) = \\cos(x)\\) even or odd?", "even", "\\(\\cos(-x) = \\cos(x)\\), so cosine is even."},
		{"Is \\(f(x) = x^{2} + 1\\) even or odd?", "even", "\\(f(-x) = (-x)^{2}+1 = x^{2}+1 = f(x)\\), so it is even."},
		{"Is \\(f(x) = |x|\\) even or odd?", "even", "\\(|-x| = |x|\\), so absolute value is even."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question: e.question, Answer: e.answer, Explanation: e.exp,
	}
}

type graphAnalysisGen struct{}

func (g *graphAnalysisGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"For \\(f(x) = 3x + 2\\), what is \\(f(0)\\)? (enter a number)", "2", "The y-intercept is f(0), the output when x = 0. Here \\(f(0) = 3(0)+2 = 2\\)."},
		{"What is the x-intercept of \\(f(x) = 2x - 4\\)? (enter a number)", "2", "X-intercepts are the real roots/solutions of f(x) = 0: \\(2x-4 = 0\\) gives \\(x = 2\\)."},
		{"If \\(f'(x) > 0\\) on an interval, is \\(f\\) increasing there? (yes/no)", "yes", "A positive derivative means the function is increasing on that interval."},
		{"Does \\(f(x) = 1/(x-2)\\) have a vertical asymptote at \\(x = 2\\)? (yes/no)", "yes", "As \\(x \\to 2\\), \\(f(x) \\to \\pm\\infty\\): a vertical asymptote at x = a means the function approaches ±∞ at that x."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question: e.question, Answer: e.answer, Explanation: e.exp,
	}
}

type funcInverseGen struct{}

func (g *funcInverseGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := rand.Intn(5) + 2
	b := rand.Intn(10) - 5
	return generator.Problem{
		Question:    fmt.Sprintf("If \\(f(x) = %dx + %d\\), what is \\(f^{-1}(x)\\)?", a, b),
		Answer:      fmt.Sprintf("(x %+d)/%d", -b, a),
		Explanation: fmt.Sprintf("\\(y = %dx + %d\\) → \\(x = \\frac{y %+d}{%d}\\) → \\(f^{-1}(x) = \\frac{x %+d}{%d}\\)", a, b, -b, a, -b, a),
		Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "neg_b": strconv.Itoa(-b)},
	}
}

type monotonicityGen struct{}

func (g *monotonicityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// The authored walkthrough describes f(x) = ax and asks which way it goes,
	// so the slope is the only free parameter. Both signs must stay reachable,
	// otherwise the question is answerable without reading it.
	a := rand.Intn(25) + 1
	if rand.Intn(2) == 0 {
		a = -a
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Is \\(f(x) = %dx\\) increasing or decreasing on \\(\\mathbb{R}\\)?", a),
		Answer:      map[bool]string{true: "increasing", false: "decreasing"}[a > 0],
		Explanation: fmt.Sprintf("Slope = %d. %s slope means the function is %s.", a, map[bool]string{true: "Positive", false: "Negative"}[a > 0], map[bool]string{true: "increasing", false: "decreasing"}[a > 0]),
		Facts:       map[string]string{"a": strconv.Itoa(a), "dir": map[bool]string{true: "increasing", false: "decreasing"}[a > 0], "sign": map[bool]string{true: "positive", false: "negative"}[a > 0]},
	}
}

type funcRationalGen struct{}

func (g *funcRationalGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// The authored walkthrough describes a/(x - b), so only a and b are free.
	// b reaches further than a because the asymptote is what the learner reads.
	a := rand.Intn(9) + 1
	b := rand.Intn(15) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("What is the vertical asymptote of \\(f(x) = \\frac{%d}{x-%d}\\)?", a, b),
		Answer:      fmt.Sprintf("%d", b),
		Explanation: fmt.Sprintf("The denominator is zero at x = %d, so there is a vertical asymptote at x = %d.", b, b),
		Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b)},
	}
}

type funcSigmoidGen struct{}

func (g *funcSigmoidGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// Every question is evaluated at a point where the sigmoid is exactly rational, so
	// the answer is an exact fraction rather than a rounded decimal and symbolic grading
	// compares it without a tolerance.
	ks := []int{2, 3, 4, 5, 6, 9, 11}
	k := ks[rand.Intn(len(ks))]
	switch rand.Intn(3) {
	case 0:
		return generator.Problem{
			Question: fmt.Sprintf("Evaluate \\(\\sigma(x)\\) at \\(x = \\ln %d\\), where \\(\\sigma(x) = 1/(1+e^{-x})\\).", k),
			Answer:   fmt.Sprintf("%d/%d", k, k+1),
			Explanation: fmt.Sprintf("\\(e^{-\\ln %d} = 1/%d\\), so \\(\\sigma(\\ln %d) = 1/(1+1/%d) = %d/%d\\).",
				k, k, k, k, k, k+1),
			Facts: map[string]string{"k": strconv.Itoa(k), "x": "ln " + strconv.Itoa(k)},
		}
	case 1:
		return generator.Problem{
			Question: fmt.Sprintf("Evaluate \\(\\sigma(x)\\) at \\(x = -\\ln %d\\), where \\(\\sigma(x) = 1/(1+e^{-x})\\).", k),
			Answer:   fmt.Sprintf("1/%d", k+1),
			Explanation: fmt.Sprintf("\\(e^{\\ln %d} = %d\\), so \\(\\sigma(-\\ln %d) = 1/(1+%d) = 1/%d\\).",
				k, k, k, k, k+1),
			Facts: map[string]string{"k": strconv.Itoa(k), "x": "-ln " + strconv.Itoa(k)},
		}
	default:
		return generator.Problem{
			Question: fmt.Sprintf("What is the slope of \\(\\sigma(kx)\\) at \\(x = 0\\), where \\(\\sigma(u) = 1/(1+e^{-u})\\) and \\(k = %d\\)?", k),
			Answer:   fmt.Sprintf("%d/4", k),
			Explanation: fmt.Sprintf("\\(\\sigma'(u) = \\sigma(u)(1-\\sigma(u))\\) and \\(\\sigma(0) = 1/2\\), so the chain rule gives \\(k \\cdot (1/2)(1/2) = %d/4\\).",
				k),
			Facts: map[string]string{"k": strconv.Itoa(k), "x": "0"},
		}
	}
}

type funcSignGen struct{}

func (g *funcSignGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	x := rand.Intn(10) - 5
	if x == 0 {
		x = 3
	}
	ans := "positive"
	if x < 0 {
		ans = "negative"
	}
	return generator.Problem{
		Question:    fmt.Sprintf("What is the sign of \\(f(x) = x\\) at \\(x = %d\\)?", x),
		Answer:      ans,
		Explanation: fmt.Sprintf("At x = %d, f(x) = %d, which is %s.", x, x, ans),
		Facts:       map[string]string{"x": strconv.Itoa(x), "ans": ans},
	}
}

// ============= inequalities =============

type ineqAbsValGen struct{}

func (g *ineqAbsValGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := rand.Intn(18) + 3
	switch rand.Intn(3) {
	case 0:
		return generator.Problem{
			Question:    fmt.Sprintf("What is the largest integer \\(x\\) with \\(|x| < %d\\)? (enter a number)", a),
			Answer:      fmt.Sprintf("%d", a-1),
			Explanation: fmt.Sprintf("\\(|x| < %d\\) means \\(-%d < x < %d\\), so the largest integer is %d", a, a, a, a-1),
			Facts:       map[string]string{"a": strconv.Itoa(a), "x": strconv.Itoa(a - 1)},
		}
	case 1:
		return generator.Problem{
			Question:    fmt.Sprintf("What is the largest integer \\(x\\) with \\(|x| \\leq %d\\)? (enter a number)", a),
			Answer:      fmt.Sprintf("%d", a),
			Explanation: fmt.Sprintf("\\(|x| \\leq %d\\) means \\(-%d \\leq x \\leq %d\\), and %d is included, so the largest integer is %d", a, a, a, a, a),
			Facts:       map[string]string{"a": strconv.Itoa(a), "x": strconv.Itoa(a)},
		}
	default:
		return generator.Problem{
			Question:    fmt.Sprintf("What is the smallest integer \\(x\\) with \\(|x| > %d\\)? (enter a number)", a),
			Answer:      fmt.Sprintf("%d", a+1),
			Explanation: fmt.Sprintf("\\(|x| > %d\\) excludes every value from \\(-%d\\) to %d, so the smallest integer above the gap is %d", a, a, a, a+1),
			Facts:       map[string]string{"a": strconv.Itoa(a), "x": strconv.Itoa(a + 1)},
		}
	}
}

type ineqIntervalGen struct{}

func (g *ineqIntervalGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := rand.Intn(8) + 1
	b := a + rand.Intn(5) + 2
	leftClosed := rand.Intn(2) == 0
	rightClosed := rand.Intn(2) == 0
	leftBracket := "("
	if leftClosed {
		leftBracket = "["
	}
	rightBracket := ")"
	if rightClosed {
		rightBracket = "]"
	}
	notation := fmt.Sprintf("%s%d, %d%s", leftBracket, a, b, rightBracket)
	smallest := rand.Intn(2) == 0
	answer := a
	extreme := "smallest"
	if smallest {
		if !leftClosed {
			answer = a + 1
		}
	} else {
		extreme = "largest"
		answer = b
		if !rightClosed {
			answer = b - 1
		}
	}
	return generator.Problem{
		Question: fmt.Sprintf("What is the %s integer in the interval \\(%s\\)? (enter a number)", extreme, notation),
		Answer:   strconv.Itoa(answer),
		Explanation: fmt.Sprintf("The interval \\(%s\\) contains integers from %d to %d, so its %s integer is %d.", notation, func() int {
			if leftClosed {
				return a
			}
			return a + 1
		}(), func() int {
			if rightClosed {
				return b
			}
			return b - 1
		}(), extreme, answer),
		Facts: map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "notation": notation, "answer": strconv.Itoa(answer)},
	}
}

type ineqIrrationalGen struct{}

func (g *ineqIrrationalGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	x := rand.Intn(8) + 2
	return generator.Problem{
		Question:    fmt.Sprintf("What is the smallest integer \\(x\\) with \\(\\sqrt{x} > %d\\)? (enter a number)", x),
		Answer:      fmt.Sprintf("%d", x*x+1),
		Explanation: fmt.Sprintf("\\(\\sqrt{x} > %d\\) → \\(x > %d^{2} = %d\\) (x ≥ 0 implied), so the smallest integer is %d", x, x, x*x, x*x+1),
		Facts:       map[string]string{"x": strconv.Itoa(x), "x2": strconv.Itoa(x * x), "ans": strconv.Itoa(x*x + 1)},
	}
}

type ineqLogGen struct{}

func (g *ineqLogGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	b := rand.Intn(3) + 2
	e := rand.Intn(3) + 1
	v := mathutil.IntPow(b, e)
	return generator.Problem{
		Question:    fmt.Sprintf("What is the smallest integer \\(x\\) with \\(\\log_{%d}(x) > %d\\)? (enter a number)", b, e),
		Answer:      fmt.Sprintf("%d", v+1),
		Explanation: fmt.Sprintf("\\(\\log_{%d}(x) > %d\\) → \\(x > %d^{%d} = %d\\) (base > 1 preserves inequality), so the smallest integer is %d", b, e, b, e, v, v+1),
		Facts:       map[string]string{"base": strconv.Itoa(b), "e": strconv.Itoa(e), "v": strconv.Itoa(v), "x": strconv.Itoa(v + 1)},
	}
}

type ineqQuadraticGen struct{}

func (g *ineqQuadraticGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	r := rand.Intn(18) + 2
	// The endpoint case is the one learners get wrong: a strict inequality
	// excludes the root, an inclusive one includes it.
	if rand.Intn(2) == 0 {
		return generator.Problem{
			Question:    fmt.Sprintf("What is the largest integer \\(x\\) with \\(x^{2} - %d < 0\\)? (enter a number)", r*r),
			Answer:      fmt.Sprintf("%d", r-1),
			Explanation: fmt.Sprintf("\\(x^{2} - %d < 0\\) → \\((x-%d)(x+%d) < 0\\) → \\(-%d < x < %d\\), so the largest integer is %d", r*r, r, r, r, r, r-1),
			Facts:       map[string]string{"r2": strconv.Itoa(r * r), "r": strconv.Itoa(r), "x": strconv.Itoa(r - 1)},
		}
	}
	return generator.Problem{
		Question:    fmt.Sprintf("What is the largest integer \\(x\\) with \\(x^{2} - %d \\leq 0\\)? (enter a number)", r*r),
		Answer:      fmt.Sprintf("%d", r),
		Explanation: fmt.Sprintf("\\(x^{2} \\leq %d\\) allows \\(x = %d\\) as well as every value between \\(-%d\\) and %d, so the largest integer is %d", r*r, r, r, r, r),
		Facts:       map[string]string{"r2": strconv.Itoa(r * r), "r": strconv.Itoa(r), "x": strconv.Itoa(r)},
	}
}

type ineqRationalGen struct{}

func (g *ineqRationalGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := rand.Intn(18) + 2
	if rand.Intn(2) == 0 {
		return generator.Problem{
			Question:    fmt.Sprintf("What is the smallest integer \\(x\\) with \\(\\frac{1}{x-%d} > 0\\)? (enter a number)", a),
			Answer:      fmt.Sprintf("%d", a+1),
			Explanation: fmt.Sprintf("The numerator 1 is always positive, so \\(\\frac{1}{x-%d} > 0\\) when \\(x-%d > 0\\) → \\(x > %d\\); the smallest integer is %d", a, a, a, a+1),
			Facts: map[string]string{
				"a": strconv.Itoa(a), "ans": strconv.Itoa(a + 1),
				"extreme": "smallest", "relation": ">",
				"positive_branch": "yes", "positive_answer": "yes",
			},
		}
	}
	return generator.Problem{
		Question:    fmt.Sprintf("What is the largest integer \\(x\\) with \\(\\frac{1}{x-%d} < 0\\)? (enter a number)", a),
		Answer:      fmt.Sprintf("%d", a-1),
		Explanation: fmt.Sprintf("A negative reciprocal needs a negative denominator, so \\(x-%d < 0\\) → \\(x < %d\\); the largest integer is %d", a, a, a-1),
		Facts: map[string]string{
			"a": strconv.Itoa(a), "ans": strconv.Itoa(a - 1),
			"extreme": "largest", "relation": "<",
			"negative_branch": "yes", "negative_answer": "yes",
		},
	}
}

type signAnalysisGen struct{}

func (g *signAnalysisGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	r := rand.Intn(5) + 1
	switch rand.Intn(4) {
	case 0:
		// production: name the zeros (unique answers).
		return generator.Problem{
			Question:    fmt.Sprintf("What are the zeros of \\(f(x) = (x+%d)(x-%d)\\)? (enter like a,b with a<b)", r, r),
			Answer:      fmt.Sprintf("%d,%d", -r, r),
			Explanation: fmt.Sprintf("\\(x+%d = 0\\) gives \\(x = %d\\); \\(x-%d = 0\\) gives \\(x = %d\\).", r, -r, r, r),
			Facts:       map[string]string{"r": strconv.Itoa(r), "lo": strconv.Itoa(-r), "hi": strconv.Itoa(r), "is_zeros": "yes"},
		}
	case 1:
		// production: sign between the roots (unique answer).
		return generator.Problem{
			Question:    fmt.Sprintf("Analyze the sign of \\(f(x) = (x+%d)(x-%d)\\) for \\(-%d < x < %d\\).", r, r, r, r),
			Answer:      "negative",
			Explanation: fmt.Sprintf("For \\(-%d < x < %d\\): \\((x+%d)\\) is positive but \\((x-%d)\\) is negative, product is negative.", r, r, r, r),
			Facts:       map[string]string{"r": strconv.Itoa(r), "is_between": "yes"},
		}
	default:
		if rand.Intn(2) == 0 {
			return generator.Problem{
				Question:    fmt.Sprintf("Analyze the sign of \\(f(x) = (x+%d)(x-%d)\\) for \\(x < -%d\\).", r, r, r),
				Answer:      "positive",
				Explanation: fmt.Sprintf("For \\(x < -%d\\): both \\((x+%d)\\) and \\((x-%d)\\) are negative, product is positive.", r, r, r),
				Facts:       map[string]string{"r": strconv.Itoa(r), "is_below": "yes"},
			}
		}
		return generator.Problem{
			Question:    fmt.Sprintf("Analyze the sign of \\(f(x) = (x+%d)(x-%d)\\) for \\(x > %d\\).", r, r, r),
			Answer:      "positive",
			Explanation: fmt.Sprintf("For \\(x > %d\\): both \\((x+%d)\\) and \\((x-%d)\\) are positive, product is positive.", r, r, r),
			Facts:       map[string]string{"r": strconv.Itoa(r), "is_above": "yes"},
		}
	}
}

type ineqSystemsGen struct{}

func (g *ineqSystemsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := rand.Intn(4) + 2
	b := rand.Intn(5) - 2
	x := rand.Intn(4) + 1
	y := a*x + b + rand.Intn(3)
	satisfies := y > a*x+b
	ans := map[bool]string{true: "yes", false: "no"}[satisfies]
	return generator.Problem{
		Question:    fmt.Sprintf("Is \\((%d,%d)\\) a solution to \\(y > %dx + %d\\)?", x, y, a, b),
		Answer:      ans,
		Explanation: fmt.Sprintf("At \\(x=%d\\): RHS = \\(%d(%d)+%d = %d\\). \\(y=%d %s %d\\), so %s.", x, a, x, b, a*x+b, y, map[bool]string{true: ">", false: "≤"}[satisfies], a*x+b, ans),
		Facts:       map[string]string{"x": strconv.Itoa(x), "y": strconv.Itoa(y), "a": strconv.Itoa(a), "b": strconv.Itoa(b), "rhs": strconv.Itoa(a*x + b), "ans": ans},
	}
}

// ============= polynomials =============

type polyDivisionGen struct{}

func (g *polyDivisionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	r := rand.Intn(6) + 1
	a := rand.Intn(4) + 1
	// (ax^2 + r*x) / x = ax + r
	return generator.Problem{
		Question:    fmt.Sprintf("Simplify: \\((%dx^{2} + %dx) / x\\)", a, r*a),
		Answer:      fmt.Sprintf("%dx + %d", a, r*a),
		Explanation: fmt.Sprintf("Divide each term: \\(%dx^{2}/x = %dx\\), \\(%dx/x = %d\\)", a, a, r*a, r*a),
		Facts:       map[string]string{"a": strconv.Itoa(a), "r": strconv.Itoa(r), "ra": strconv.Itoa(r * a)},
	}
}

type polyMonomialGen struct{}

func (g *polyMonomialGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := rand.Intn(5) + 1
	n := rand.Intn(4) + 2
	return generator.Problem{
		Question:    fmt.Sprintf("What is the degree of \\(%dx^{%d}\\)?", a, n),
		Answer:      fmt.Sprintf("%d", n),
		Explanation: fmt.Sprintf("The monomial \\(%dx^{%d}\\) has degree %d.", a, n, n),
		Facts:       map[string]string{"a": strconv.Itoa(a), "n": strconv.Itoa(n)},
	}
}

type polyRootsGen struct{}

func (g *polyRootsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	r1 := rand.Intn(6) - 3
	r2 := rand.Intn(6) - 3
	b := -(r1 + r2)
	c := r1 * r2
	// The stem comes from formatQuadratic, which prints "x^2 - 2x - 3". The
	// old "x^{2} + %dx + %d" format printed "x^2 + -2x + 0" whenever a
	// coefficient was negative or zero, which is what report 224 quoted.
	q := formatQuadratic(b, c)
	answer := fmt.Sprintf("%d,%d", r1, r2)
	explanation := fmt.Sprintf("\\((x %+d)(x %+d) = 0\\) \u2192 roots are %d and %d", -r1, -r2, r1, r2)
	if r1 == r2 {
		// A double root is one root, not two: publishing "2,2" rejected the
		// learner who wrote 2, the same defect the quadratic solver had.
		answer = fmt.Sprintf("%d", r1)
		explanation = fmt.Sprintf("\\((x %+d)^{2} = 0\\), so %d is a double root", -r1, r1)
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Find the roots of \\(%s = 0\\).", q),
		Answer:      answer,
		Explanation: explanation,
		Facts:       map[string]string{"q": q, "b": strconv.Itoa(b), "c": strconv.Itoa(c), "r1": strconv.Itoa(r1), "r2": strconv.Itoa(r2)},
	}
}

type synthDivGen struct{}

func (g *synthDivGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	r := rand.Intn(6) - 3
	a := rand.Intn(4) + 1
	b := rand.Intn(5) - 2
	// (ax^2 + bx - ar) / (x - r) = ax + (ar + b)
	return generator.Problem{
		Question:    fmt.Sprintf("Use synthetic division to divide \\((%dx^{2} + %dx)\\) by \\((x - %d)\\). What is the quotient?", a, b, r),
		Answer:      fmt.Sprintf("%dx + %d", a, b+a*r),
		Explanation: fmt.Sprintf("Synthetic division by %d gives coefficients %d and %d → \\(%dx + %d\\)", r, a, b+a*r, a, b+a*r),
		Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "r": strconv.Itoa(r), "out_a": strconv.Itoa(a), "out_b": strconv.Itoa(b + a*r)},
	}
}

type vietaGen struct{}

func (g *vietaGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	r1 := rand.Intn(6) + 1
	r2 := rand.Intn(6) + 1
	b := -(r1 + r2)
	c := r1 * r2
	return generator.Problem{
		Question:    fmt.Sprintf("For \\(x^{2} + %dx + %d = 0\\), what is the sum of the roots?", b, c),
		Answer:      fmt.Sprintf("%d", r1+r2),
		Explanation: fmt.Sprintf("By Vieta: sum of roots = %d. Roots are %d and %d, sum = %d", -b, r1, r2, r1+r2),
		Facts:       map[string]string{"q": formatQuadratic(b, c), "b": strconv.Itoa(b), "c": strconv.Itoa(c), "sum": strconv.Itoa(r1 + r2), "neg_b": strconv.Itoa(-b), "r1": strconv.Itoa(r1), "r2": strconv.Itoa(r2)},
	}
}

// ============= quadratics =============

type quadComplexGen struct{}

func (g *quadComplexGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	b := rand.Intn(max(5, scale*2)) + 2
	if scale > 3 {
		b = rand.Intn(12) + 6
	}
	c := (b*b)/4 + 1
	disc := b*b - 4*c
	if rand.Intn(3) == 0 {
		// production: discriminant value and complex-solution count (unique answers).
		if rand.Intn(2) == 0 {
			return generator.Problem{
				Question:    fmt.Sprintf("What is the discriminant of \\(x^{2} + %dx + %d = 0\\)? (enter a number)", b, c),
				Answer:      fmt.Sprintf("%d", disc),
				Explanation: fmt.Sprintf("Discriminant = \\(%d^{2} - 4(%d) = %d - %d = %d < 0\\) → two complex (non-real) solutions.", b, c, b*b, 4*c, disc),
				Facts:       map[string]string{"b": strconv.Itoa(b), "c": strconv.Itoa(c), "b2": strconv.Itoa(b * b), "four_c": strconv.Itoa(4 * c), "disc": strconv.Itoa(disc), "want_disc": "yes"},
			}
		}
		return generator.Problem{
			Question:    fmt.Sprintf("How many complex (non-real) solutions does \\(x^{2} + %dx + %d = 0\\) have? (enter a number)", b, c),
			Answer:      "2",
			Explanation: fmt.Sprintf("Discriminant = \\(%d^{2} - 4(%d) = %d < 0\\), so both solutions are complex (non-real).", b, c, disc),
			Facts:       map[string]string{"b": strconv.Itoa(b), "c": strconv.Itoa(c), "b2": strconv.Itoa(b * b), "four_c": strconv.Itoa(4 * c), "disc": strconv.Itoa(disc), "want_complex": "yes"},
		}
	}
	return generator.Problem{
		Question:    fmt.Sprintf("How many real solutions does \\(x^{2} + %dx + %d = 0\\) have?", b, c),
		Answer:      "0",
		Explanation: fmt.Sprintf("Discriminant = \\(%d^{2} - 4(%d) = %d - %d = %d < 0\\) → two complex (non-real) solutions.", b, c, b*b, 4*c, disc),
		Facts:       map[string]string{"b": strconv.Itoa(b), "c": strconv.Itoa(c), "b2": strconv.Itoa(b * b), "four_c": strconv.Itoa(4 * c), "disc": strconv.Itoa(disc), "want_real": "yes"},
	}
}

type quadIncompleteGen struct{}

func (g *quadIncompleteGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	k := rand.Intn(18) + 2
	if rand.Intn(2) == 0 {
		return generator.Problem{
			Question:    fmt.Sprintf("Solve: \\(x^{2} - %d = 0\\)", k*k),
			Answer:      fmt.Sprintf("%d,%d", k, -k),
			Explanation: fmt.Sprintf("\\(x^{2} = %d\\) → \\(x = \\pm \\sqrt{%d} = \\pm %d\\)", k*k, k*k, k),
			Facts:       map[string]string{"k2": strconv.Itoa(k * k), "k": strconv.Itoa(k)},
		}
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Solve: \\(x^{2} = %d\\)", k*k),
		Answer:      fmt.Sprintf("%d,%d", k, -k),
		Explanation: fmt.Sprintf("Taking the square root of both sides gives \\(x = \\pm %d\\), and %d is a perfect square, so both roots are integers.", k, k*k),
		Facts:       map[string]string{"k2": strconv.Itoa(k * k), "k": strconv.Itoa(k)},
	}
}

// quadParametricGen asks for the positive k that makes x² + kx + p² a perfect
// square, so the discriminant is k² - 4p².
//
// The parameter has to appear in the printed equation. Before 23 October 2026
// the question showed x² + 2px + p² and asked for k, which was not in it, and
// the explanation "solved" a discriminant that is already zero for every p.
// Nine learners reported the same missing k that day, which is the whole reason
// this comment names the failure rather than the shape.
type quadParametricGen struct{}

func (g *quadParametricGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	p := rand.Intn(29) + 2
	k := 2 * p
	// Three ways of asking the same discriminant condition, so the answer set is
	// not one number per parameter: k from the constant, the constant from k, and
	// the negative middle coefficient.
	switch rand.Intn(3) {
	case 0:
		return generator.Problem{
			Question:    fmt.Sprintf("For what positive value of \\(k\\) does \\(x^{2} + kx + %d = 0\\) have exactly one solution?", p*p),
			Answer:      strconv.Itoa(k),
			Explanation: fmt.Sprintf("Exactly one solution means the discriminant is zero: \\(k^{2} - 4(%d) = 0\\), so \\(k^{2} = %d\\) and the positive value is \\(k = %d\\).", p*p, k*k, k),
			Facts:       map[string]string{"shape": "askK", "coef": strconv.Itoa(k), "const": strconv.Itoa(p * p), "k": strconv.Itoa(k), "answer": strconv.Itoa(k)},
		}
	case 1:
		return generator.Problem{
			Question:    fmt.Sprintf("For what positive value of \\(c\\) does \\(x^{2} + %dx + c = 0\\) have exactly one solution?", k),
			Answer:      strconv.Itoa(p * p),
			Explanation: fmt.Sprintf("The discriminant is \\(%d^{2} - 4c\\), and it must be zero: \\(4c = %d\\), so \\(c = %d\\).", k, k*k, p*p),
			Facts:       map[string]string{"shape": "askC", "coef": strconv.Itoa(k), "const": strconv.Itoa(p * p), "k": strconv.Itoa(k), "answer": strconv.Itoa(p * p)},
		}
	default:
		return generator.Problem{
			Question:    fmt.Sprintf("For what positive value of \\(k\\) does \\(x^{2} - kx + %d = 0\\) have exactly one solution?", p*p),
			Answer:      strconv.Itoa(k),
			Explanation: fmt.Sprintf("The discriminant is \\(k^{2} - 4(%d)\\), which is zero when \\(k = %d\\) for the positive choice.", p*p, k),
			Facts:       map[string]string{"shape": "askKNeg", "coef": strconv.Itoa(-k), "const": strconv.Itoa(p * p), "k": strconv.Itoa(k), "answer": strconv.Itoa(k)},
		}
	}
}

type quadQuadraticGen struct{}

func (g *quadQuadraticGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// Gradable numeric forms: the standard-form / formula / discriminant facts
	// live in the Explanation (never graded); the graded answer is a number.
	switch rand.Intn(4) {
	case 0:
		b := rand.Intn(5) + 1
		c := rand.Intn(5) + 1
		return generator.Problem{
			Question:    fmt.Sprintf("For \\(x^{2} + %dx + %d = 0\\) in standard form \\(ax^{2} + bx + c = 0\\), what is \\(b\\)? (enter a number)", b, c),
			Answer:      fmt.Sprintf("%d", b),
			Explanation: fmt.Sprintf("Standard form is \\(ax^{2} + bx + c = 0\\) with \\(a \\neq 0\\); here \\(a=1\\), \\(b=%d\\), \\(c=%d\\)", b, c),
			Facts:       map[string]string{"b": strconv.Itoa(b), "c": strconv.Itoa(c), "is_standard": "yes"},
		}
	case 1:
		r1 := rand.Intn(6) + 2
		r2 := rand.Intn(6) + 2
		b := -(r1 + r2)
		c := r1 * r2
		lo := min(r1, r2)
		return generator.Problem{
			Question:    fmt.Sprintf("Solve with the quadratic formula: %s = 0. What is the smaller root? (enter a number)", formatQuadratic(b, c)),
			Answer:      fmt.Sprintf("%d", lo),
			Explanation: fmt.Sprintf("\\(x = [-b \\pm \\sqrt{b^{2}-4ac}]/(2a)\\) gives \\(x = %d\\) or \\(x = %d\\); the smaller root is %d. The formula solves \\(ax^{2} + bx + c = 0\\).", r1, r2, lo),
			Facts:       map[string]string{"q": formatQuadratic(b, c), "r1": strconv.Itoa(r1), "r2": strconv.Itoa(r2), "lo": strconv.Itoa(lo)},
		}
	case 2:
		b := rand.Intn(5) + 2
		c := rand.Intn(5) + 1
		disc := b*b - 4*c
		return generator.Problem{
			Question:    fmt.Sprintf("For \\(x^{2} + %dx + %d = 0\\), what is the discriminant \\(b^{2} - 4ac\\)? (enter a number)", b, c),
			Answer:      fmt.Sprintf("%d", disc),
			Explanation: fmt.Sprintf("The discriminant \\(b^{2}-4ac = %d\\) determines the number and type of solutions: \\(> 0\\): two real; \\(= 0\\): one real; \\(< 0\\): two complex.", disc),
			Facts:       map[string]string{"b": strconv.Itoa(b), "c": strconv.Itoa(c), "disc": strconv.Itoa(disc)},
		}
	default:
		return generator.Problem{
			Question:    "The graph of a quadratic function is a ____.",
			Answer:      "parabola",
			Explanation: "\\(f(x) = ax^{2} + bx + c\\) graphs as a parabola.",
			Facts:       map[string]string{"parabola": "yes"},
		}
	}
}

// ============= systems =============

type systemsConceptGen struct{}

func (g *systemsConceptGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"Is \\((2,3)\\) a solution of the system \\(x + y = 5\\), \\(x - y = -1\\)? (yes/no)", "yes", "To solve a system is to find values satisfying all equations simultaneously: \\(2+3 = 5\\) and \\(2-3 = -1\\), so yes."},
		{"Two distinct parallel lines intersect in how many points? (enter a number)", "0", "A system of two linear equations has 0 solutions (parallel distinct lines), 1 (intersecting lines), or infinitely many (the same line)."},
		{"Is the system \\(x + y = 3\\), \\(x - y = 1\\) consistent? (yes/no)", "yes", "A consistent system has at least one solution: \\(x=2\\), \\(y=1\\) satisfies both equations, so yes."},
		{"How many solutions does an inconsistent system have? (enter a number)", "0", "An inconsistent system has no solution (parallel lines that never intersect)."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question: e.question, Answer: e.answer, Explanation: e.exp,
	}
}

type gaussianElimGen struct{}

func (g *gaussianElimGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// 2x2 system with integer solution
	x := rand.Intn(5) - 2
	y := rand.Intn(5) - 2
	a := rand.Intn(4) + 1
	b := rand.Intn(4) + 1
	c := rand.Intn(4) + 1
	d := rand.Intn(4) + 1
	// Ensure determinant ≠ 0
	for a*d-b*c == 0 {
		d++
	}
	e1 := a*x + b*y
	e2 := c*x + d*y
	if rand.Intn(2) == 0 {
		return generator.Problem{
			Question:    fmt.Sprintf("Solve using elimination: \\(%dx + %dy = %d\\), \\(%dx + %dy = %d\\)", a, b, e1, c, d, e2),
			Answer:      fmt.Sprintf("(%d,%d)", x, y),
			Explanation: fmt.Sprintf("Solution: \\(x=%d\\), \\(y=%d\\)", x, y),
			Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "c": strconv.Itoa(c), "d": strconv.Itoa(d), "e1": strconv.Itoa(e1), "e2": strconv.Itoa(e2), "x": strconv.Itoa(x), "y": strconv.Itoa(y), "is_solve": "yes"},
		}
	}
	// Row echelon form question
	return generator.Problem{
		Question:    fmt.Sprintf("Put the system \\(%dx + %dy = %d\\), \\(%dx + %dy = %d\\) into row echelon form.", a, b, e1, c, d, e2),
		Answer:      fmt.Sprintf("[[%d,%d,%d],[0,%d,%d]]", a, b, e1, d-c*b/a, e2-c*e1/a),
		Explanation: fmt.Sprintf("The augmented matrix [[%d,%d,%d],[%d,%d,%d]] can be reduced.", a, b, e1, c, d, e2),
		Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "c": strconv.Itoa(c), "d": strconv.Itoa(d), "e1": strconv.Itoa(e1), "e2": strconv.Itoa(e2), "is_echelon": "yes"},
	}
}

type injectivityGen struct{}

func (g *injectivityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		question string
		answer   string
		reason   string
	}
	tableEasy := []entry{
		{"Is f(x) = 2x + 3 injective? (yes/no)", "yes", "Strictly monotonic: distinct inputs give distinct outputs."},
		{"Is f(x) = x^2 injective on all of R? (yes/no)", "no", "f(-1) = f(1) = 1."},
		{"Is f(x) = x^3 bijective from R to R? (yes/no)", "yes", "Strictly increasing and unbounded both ways."},
		{"A function passes inputs through one-to-one mapping. At most how many preimages can one output have? (enter a number)", "1", "Injective means at most one."},
	}
	tableHard := []entry{
		{"Is f(x) = e^x surjective onto R? (yes/no)", "no", "Range is (0, inf): negatives never hit."},
		{"If f is bijective, how many inverses does each output have? (enter a number)", "1", "Exactly one: bijective is one-to-one and onto."},
		{"Is the composition of two injective functions injective? (yes/no)", "yes", "Distinct inputs stay distinct through both maps."},
	}
	pool := tableEasy
	if scale > 3 {
		pool = append(tableEasy, tableHard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.reason}
}

type powerFuncGen struct{}

func (g *powerFuncGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		question string
		answer   string
		reason   string
	}
	tableEasy := []entry{
		{"What is 2^5? (enter a number)", "32", "2·2·2·2·2 = 32."},
		{"For x^a with a > 0, is the function increasing for x > 0? (yes/no)", "yes", "Positive exponent: larger base gives larger value."},
		{"What is 10^3? (enter a number)", "1000", "1 followed by 3 zeros."},
		{"Is x^2 an even function? (yes/no)", "yes", "(-x)^2 = x^2."},
	}
	tableHard := []entry{
		{"What is 3^4? (enter a number)", "81", "3·3·3·3 = 81."},
		{"For 0 < a < 1, is x^a increasing or decreasing for x > 0? (type increasing or decreasing)", "increasing", "Still increasing, just slower (concave)."},
		{"What is (-2)^3? (enter a number)", "-8", "Odd power preserves the sign."},
	}
	pool := tableEasy
	if scale > 3 {
		pool = append(tableEasy, tableHard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.reason}
}

type irrationalFuncGen struct{}

func (g *irrationalFuncGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		question string
		answer   string
		reason   string
	}
	tableEasy := []entry{
		{"What is sqrt(49)? (enter a number)", "7", "7·7 = 49."},
		{"What is cbrt(-8)? (enter a number)", "-2", "(-2)^3 = -8; odd roots take negatives."},
		{"Is sqrt(x) defined at x = -1 over the reals? (yes/no)", "no", "No real square of a negative."},
		{"Does sqrt(x^2) = |x| hold for all real x? (yes/no)", "yes", "Principal root is nonnegative by definition."},
	}
	tableHard := []entry{
		{"What is sqrt(144)? (enter a number)", "12", "12·12 = 144."},
		{"What is cbrt(27)? (enter a number)", "3", "3·3·3 = 27."},
		{"Is 4th-root(x) defined for negative x over the reals? (yes/no)", "no", "Even roots need nonnegative radicands."},
	}
	pool := tableEasy
	if scale > 3 {
		pool = append(tableEasy, tableHard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.reason}
}
