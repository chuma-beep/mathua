package algebra

import (
	"fmt"
	"math/rand"
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
	reg.Register("alg.eq.literal", &literalEqGen{})
	reg.Register("alg.ineq.multi_step", &multiStepIneqGen{})
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
	reg.Register("alg.eq.extraneous_roots", &extraneousRootsGen{})
	reg.Register("alg.eq.irrational", &eqIrrationalGen{})
	reg.Register("alg.eq.logarithmic", &eqLogGen{})
	reg.Register("alg.eq.polynomial", &eqPolyGen{})
	reg.Register("alg.eq.rational", &eqRationalGen{})
	reg.Register("alg.eq.trinomial", &eqTrinomialGen{})
	reg.Register("alg.func.absolute_value", &funcAbsValGen{})
	reg.Register("alg.func.composite", &funcCompositeGen{})
	reg.Register("alg.func.dirichlet", &funcDirichletGen{})
	reg.Register("alg.func.domain", &funcDomainGen{})
	reg.Register("alg.func.even_odd", &funcEvenOddGen{})
	reg.Register("alg.func.graph_analysis", &graphAnalysisGen{})
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
	reg.Register("alg.systems.concept", &systemsConceptGen{})
	reg.Register("alg.systems.gaussian", &gaussianElimGen{})
}

// ---------------------------------------------------------------------------
// Linear
// ---------------------------------------------------------------------------

type slopeGen struct{}

func (g *slopeGen) Generate(difficulty float64) generator.Problem {
	// Higher difficulty → wider range, negative slopes more likely
	lim := int(3 + difficulty*12)
	x1 := rand.Intn(lim)
	y1 := rand.Intn(lim)
	x2 := x1 + rand.Intn(int(difficulty*7)+2) + 1
	y2 := y1 + rand.Intn(lim) + 1
	if rand.Intn(2) == 0 || difficulty > 0.6 {
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
	}
}

type slopeInterceptGen struct{}

func (g *slopeInterceptGen) Generate(difficulty float64) generator.Problem {
	// Higher difficulty → fractional slopes and larger intercepts
	if difficulty > 0.6 && rand.Intn(2) == 0 {
		m := rand.Intn(3) + 1
		n := rand.Intn(3) + 2
		if rand.Intn(2) == 0 { m = -m }
		b := rand.Intn(int(5+difficulty*10)) - int(3+difficulty*5)
		return generator.Problem{
			Question:    fmt.Sprintf("Write the equation of a line with slope \\(%d/%d\\) and \\(y\\)-intercept %d (\\(y = mx + b\\)).", m, n, b),
			Answer:      fmt.Sprintf("y = (%d/%d)x + %d", m, n, b),
			Explanation: fmt.Sprintf("\\(y = \\frac{%d}{%d}x + %d\\)", m, n, b),
		}
	}
	m := rand.Intn(int(1+difficulty*5)) + 1
	if rand.Intn(2) == 0 {
		m = -m
	}
	b := rand.Intn(int(5+difficulty*8)) - int(3+difficulty*4)
	return generator.Problem{
		Question:    fmt.Sprintf("Write the equation of a line with slope %d and \\(y\\)-intercept %d (\\(y = mx + b\\)).", m, b),
		Answer:      formatLinear(m, b),
		Explanation: fmt.Sprintf("\\(y = %dx + %d\\)", m, b),
	}
}

type linearGraphGen struct{}

func (g *linearGraphGen) Generate(difficulty float64) generator.Problem {
	m := rand.Intn(5) + 1
	if rand.Intn(2) == 0 {
		m = -m
	}
	b := rand.Intn(10) - 5
	x := rand.Intn(5) + 1
	y := m*x + b
	return generator.Problem{
		Question:    fmt.Sprintf("Given \\(y = %s\\), find \\(y\\) when \\(x = %d\\).", formatLinear(m, b), x),
		Answer:      fmt.Sprintf("%d", y),
		Explanation: fmt.Sprintf("\\(y = %d(%d) + %d = %d + %d = %d\\)", m, x, b, m*x, b, y),
	}
}

type stdFormGen struct{}

func (g *stdFormGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(5) + 2
	b := rand.Intn(5) + 1
	x := rand.Intn(5) + 1
	y := rand.Intn(5) + 1
	c := a*x + b*y
	if rand.Intn(2) == 0 {
		return generator.Problem{
			Question:    fmt.Sprintf("Is \\((%d,%d)\\) a solution to \\(%dx + %dy = %d\\)?", x, y, a, b, c),
			Answer:      "yes",
			Explanation: fmt.Sprintf("\\(%d(%d) + %d(%d) = %d + %d = %d\\). Yes!", a, x, b, y, a*x, b*y, c),
		}
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Is \\((%d,%d)\\) a solution to \\(%dx + %dy = %d\\)?", x+1, y, a, b, c),
		Answer:      "no",
		Explanation: fmt.Sprintf("\\(%d(%d) + %d(%d) = %d + %d = %d\\), not \\(%d\\). No!", a, x+1, b, y, a*(x+1), b*y, a*(x+1)+b*y, c),
	}
}

type parallelPerpGen struct{}

func (g *parallelPerpGen) Generate(difficulty float64) generator.Problem {
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
		}
	}
	perp := reduce(-1, m)
	return generator.Problem{
		Question:    fmt.Sprintf("A line perpendicular to \\(y = %dx + 3\\) has slope:", m),
		Answer:      fracOrInt(perp.num, perp.den),
		Explanation: fmt.Sprintf("Perpendicular slopes are negative reciprocals: \\(-\\frac{1}{%d} = %s\\)", m, fracOrInt(perp.num, perp.den)),
	}
}

// ---------------------------------------------------------------------------
// Equations & inequalities
// ---------------------------------------------------------------------------

type multiStepEqGen struct{}

func (g *multiStepEqGen) Generate(difficulty float64) generator.Problem {
	scale := int(1 + difficulty*8)
	x := rand.Intn(scale*2) + 2
	a := rand.Intn(scale) + 2
	b := rand.Intn(scale*2) + 2
	c := a*x + b
	return generator.Problem{
		Question:    fmt.Sprintf("Solve: \\(%dx + %d = %d\\)", a, b, c),
		Answer:      fmt.Sprintf("%d", x),
		Explanation: fmt.Sprintf("\\(%dx + %d = %d\\) → \\(%dx = %d\\) → \\(x = %d\\)", a, b, c, a, c-b, x),
	}
}

type varsBothSidesGen struct{}

func (g *varsBothSidesGen) Generate(difficulty float64) generator.Problem {
	scale := int(1 + difficulty*7)
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
	}
}

type literalEqGen struct{}

func (g *literalEqGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(6) + 2
	b := rand.Intn(5) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("Solve for \\(y\\): \\(%dx + %dy = z\\)", a, b),
		Answer:      fmt.Sprintf("y = (z - %dx)/%d", a, b),
		Explanation: fmt.Sprintf("\\(%dy = z - %dx\\) → \\(y = \\frac{z - %dx}{%d}\\)", b, a, a, b),
	}
}

type multiStepIneqGen struct{}

func (g *multiStepIneqGen) Generate(difficulty float64) generator.Problem {
	x := rand.Intn(8) + 2
	a := rand.Intn(6) + 2
	b := rand.Intn(10) + 1
	c := a*(x-1) + b + rand.Intn(a)
	return generator.Problem{
		Question:    fmt.Sprintf("Solve: \\(%dx + %d > %d\\)", a, b, c),
		Answer:      fmt.Sprintf("%d", x),
		Explanation: fmt.Sprintf("\\(%dx + %d > %d\\) → \\(%dx > %d\\) → \\(x > %d\\), so \\(x \\ge %d\\)", a, b, c, a, c-b, c-b, x),
	}
}

type compoundIneqGen struct{}

func (g *compoundIneqGen) Generate(difficulty float64) generator.Problem {
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
	}
}

// ---------------------------------------------------------------------------
// Systems
// ---------------------------------------------------------------------------

type sysSubstitutionGen struct{}

func (g *sysSubstitutionGen) Generate(difficulty float64) generator.Problem {
	x := rand.Intn(8) + 2
	y := rand.Intn(8) + 2
	return generator.Problem{
		Question:    fmt.Sprintf("Solve: \\(y = %d\\),  \\(x + y = %d\\)", y, x+y),
		Answer:      fmt.Sprintf("(%d,%d)", x, y),
		Explanation: fmt.Sprintf("Substitute \\(y=%d\\): \\(x + %d = %d\\) → \\(x = %d\\). Solution: \\((%d,%d)\\)", y, y, x+y, x, x, y),
	}
}

type sysEliminationGen struct{}

func (g *sysEliminationGen) Generate(difficulty float64) generator.Problem {
	x := rand.Intn(8) + 2
	y := rand.Intn(8) + 2
	a := rand.Intn(5) + 2
	b := rand.Intn(5) + 2
	c := rand.Intn(5) + 2
	eq1 := a*x + b*y
	eq2 := c*x + b*y
	return generator.Problem{
		Question:    fmt.Sprintf("Solve: \\(%dx + %dy = %d\\),  \\(%dx + %dy = %d\\)", a, b, eq1, c, b, eq2),
		Answer:      fmt.Sprintf("(%d,%d)", x, y),
		Explanation: fmt.Sprintf("Subtract: \\(%dx = %d\\) → \\(x = %d\\); \\(y = %d\\). Solution: \\((%d,%d)\\)", a-c, eq1-eq2, x, y, x, y),
	}
}

type sysWordGen struct{}

func (g *sysWordGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(8) + 3
	b := rand.Intn(8) + 3
	sum := a + b
	diff := a - b
	return generator.Problem{
		Question:    fmt.Sprintf("Two numbers sum to %d and differ by %d. Find both numbers (smaller first).", sum, diff),
		Answer:      fmt.Sprintf("%d,%d", b, a),
		Explanation: fmt.Sprintf("\\(x + y = %d\\), \\(x - y = %d\\). Add: \\(2x = %d\\) → \\(x = %d\\). \\(y = %d\\)", sum, diff, a+b, a, b),
	}
}

// ---------------------------------------------------------------------------
// Polynomials
// ---------------------------------------------------------------------------

type polyConceptGen struct{}

func (g *polyConceptGen) Generate(difficulty float64) generator.Problem {
	coeff := rand.Intn(6) + 2
	exp := rand.Intn(3) + 2
	return generator.Problem{
		Question:    fmt.Sprintf("What is the degree of \\(%dx^{%d} + 3x + 1\\)?", coeff, exp),
		Answer:      fmt.Sprintf("%d", exp),
		Explanation: fmt.Sprintf("The highest exponent is %d, so the degree is %d.", exp, exp),
	}
}

type polyAddSubGen struct{}

func (g *polyAddSubGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(5) + 2
	b := rand.Intn(5) + 2
	c := rand.Intn(5) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("Simplify: \\((%dx^{2} + %dx) + (%dx^{2} + %dx)\\)", a, c, b, c),
		Answer:      fmt.Sprintf("%dx^2 + %dx", a+b, 2*c),
		Explanation: fmt.Sprintf("\\(%dx^{2} + %dx^{2} = %dx^{2}\\); \\(%dx + %dx = %dx\\)", a, b, a+b, c, c, 2*c),
	}
}

type polyMultMonoGen struct{}

func (g *polyMultMonoGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(5) + 2
	b := rand.Intn(5) + 2
	c := rand.Intn(5) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("Simplify: \\(%dx(%dx + %d)\\)", a, b, c),
		Answer:      fmt.Sprintf("%dx^2 + %dx", a*b, a*c),
		Explanation: fmt.Sprintf("\\(%dx(%dx) + %dx(%d) = %dx^{2} + %dx\\)", a, b, a, c, a*b, a*c),
	}
}

type polyFoilGen struct{}

func (g *polyFoilGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(6) + 1
	b := rand.Intn(6) + 1
	c := rand.Intn(6) + 1
	d := rand.Intn(6) + 1
	return generator.Problem{
		Question: fmt.Sprintf("Simplify: \\((%dx + %d)(%dx + %d)\\)", a, b, c, d),
		Answer:   fmt.Sprintf("%dx^2 + %dx + %d", a*c, a*d+b*c, b*d),
		Explanation: fmt.Sprintf("FOIL: \\(%dx^{2} + %dx + %dx + %d = %dx^{2} + %dx + %d\\)",
			a*c, a*d, b*c, b*d, a*c, a*d+b*c, b*d),
	}
}

type polySpecialGen struct{}

func (g *polySpecialGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(6) + 1
	b := rand.Intn(5) + 1
	return generator.Problem{
		Question: fmt.Sprintf("Simplify: \\((%dx + %d)^{2}\\)", a, b),
		Answer:   fmt.Sprintf("%dx^2 + %dx + %d", a*a, 2*a*b, b*b),
		Explanation: fmt.Sprintf("\\((%dx)^{2} + 2(%dx)(%d) + %d^{2} = %dx^{2} + %dx + %d\\)",
			a, a, b, b, a*a, 2*a*b, b*b),
	}
}

// ---------------------------------------------------------------------------
// Factoring
// ---------------------------------------------------------------------------

type factorGCFGen struct{}

func (g *factorGCFGen) Generate(difficulty float64) generator.Problem {
	f := rand.Intn(5) + 2
	a := f * (rand.Intn(5) + 2)
	b := f * (rand.Intn(5) + 1)
	return generator.Problem{
		Question:    fmt.Sprintf("Factor: \\(%dx + %d\\)", a, b),
		Answer:      fmt.Sprintf("%d(%dx + %d)", f, a/f, b/f),
		Explanation: fmt.Sprintf("GCF is \\(%d\\): \\(%dx + %d = %d(%dx + %d)\\)", f, a, b, f, a/f, b/f),
	}
}

type factorTrinomialGen struct{}

func (g *factorTrinomialGen) Generate(difficulty float64) generator.Problem {
	r1 := rand.Intn(6) + 1
	r2 := rand.Intn(6) + 1
	b := r1 + r2
	c := r1 * r2
	return generator.Problem{
		Question:    fmt.Sprintf("Factor: \\(x^{2} + %dx + %d\\)", b, c),
		Answer:      fmt.Sprintf("(x + %d)(x + %d)", r1, r2),
		Explanation: fmt.Sprintf("Find two numbers that multiply to %d and add to %d: %d and %d.", c, b, r1, r2),
	}
}

type factorDiffSquaresGen struct{}

func (g *factorDiffSquaresGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(8) + 2
	return generator.Problem{
		Question:    fmt.Sprintf("Factor: \\(x^{2} - %d\\)", a*a),
		Answer:      fmt.Sprintf("(x + %d)(x - %d)", a, a),
		Explanation: fmt.Sprintf("\\(x^{2} - %d = (x + %d)(x - %d)\\)", a*a, a, a),
	}
}

type factorACMethodGen struct{}

func (g *factorACMethodGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(4) + 2
	r1 := rand.Intn(4) + 1
	r2 := rand.Intn(4) + 1
	b := a + r1*r2
	c := r1 * r2
	return generator.Problem{
		Question:    fmt.Sprintf("Factor: \\(x^{2} + %dx + %d\\)", b, c),
		Answer:      fmt.Sprintf("(x + %d)(x + %d)", r1, r2),
		Explanation: fmt.Sprintf("Find numbers multiplying to %d and adding to %d: %d and %d.", c, b, r1, r2),
	}
}

// ---------------------------------------------------------------------------
// Quadratics
// ---------------------------------------------------------------------------

type quadSolveFactorGen struct{}

func (g *quadSolveFactorGen) Generate(difficulty float64) generator.Problem {
	r1 := rand.Intn(8) - 4
	r2 := rand.Intn(8) - 4
	b := -(r1 + r2)
	c := r1 * r2
	q := formatQuadratic(b, c)
	return generator.Problem{
		Question:    fmt.Sprintf("Solve: %s = 0", q),
		Answer:      fmt.Sprintf("%d,%d", r1, r2),
		Explanation: fmt.Sprintf("Factors: \\((x %+d)(x %+d) = 0\\), so \\(x = %d\\) or \\(x = %d\\)", -r1, -r2, r1, r2),
	}
}

type quadCompleteSquareGen struct{}

func (g *quadCompleteSquareGen) Generate(difficulty float64) generator.Problem {
	r := rand.Intn(6) + 1
	b := -2 * r
	return generator.Problem{
		Question:    fmt.Sprintf("Complete the square: \\(x^{2} + %dx + \\_\\_)\\) to make a perfect square.", b),
		Answer:      fmt.Sprintf("%d", r*r),
		Explanation: fmt.Sprintf("\\(\\left(\\frac{b}{2}\\right)^{2} = \\left(\\frac{%d}{2}\\right)^{2} = %d\\)", b, r*r),
	}
}

type quadFormulaGen struct{}

func (g *quadFormulaGen) Generate(difficulty float64) generator.Problem {
	r1 := rand.Intn(8) + 1
	r2 := rand.Intn(8) + 1
	b := -(r1 + r2)
	c := r1 * r2
	q := formatQuadratic(b, c)
	return generator.Problem{
		Question:    fmt.Sprintf("Solve using quadratic formula: %s = 0", q),
		Answer:      fmt.Sprintf("%d,%d", r1, r2),
		Explanation: fmt.Sprintf("\\(x = \\frac{%d \\pm \\sqrt{%d - 4(%d)}}{2} = %d\\) or \\(%d\\)", b, b*b, c, r1, r2),
	}
}

type quadDiscriminantGen struct{}

func (g *quadDiscriminantGen) Generate(difficulty float64) generator.Problem {
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
		Question:    fmt.Sprintf("How many real solutions? \\(x^{2} + %dx + %d = 0\\)", b, c),
		Answer:      count,
		Explanation: fmt.Sprintf("Discriminant: \\(%d^{2} - 4(%d) = %d\\). %s solutions.", b, c, disc, count),
	}
}

// ---------------------------------------------------------------------------
// Functions
// ---------------------------------------------------------------------------

type funcConceptGen struct{}

func (g *funcConceptGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(5) + 2
	b := rand.Intn(5) + 1
	x := rand.Intn(5) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("If \\(f(x) = %dx + %d\\), what is \\(f(%d)\\)?", a, b, x),
		Answer:      fmt.Sprintf("%d", a*x+b),
		Explanation: fmt.Sprintf("\\(f(%d) = %d(%d) + %d = %d\\)", x, a, x, b, a*x+b),
	}
}

type funcNotationGen struct{}

func (g *funcNotationGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(5) + 2
	x := rand.Intn(5) + 1
	b := rand.Intn(5) + 1
	y := a*x + b
	return generator.Problem{
		Question:    fmt.Sprintf("\\(f(x) = %dx + %d\\). Find \\(x\\) when \\(f(x) = %d\\).", a, b, y),
		Answer:      fmt.Sprintf("%d", x),
		Explanation: fmt.Sprintf("\\(%dx + %d = %d\\) → \\(%dx = %d\\) → \\(x = %d\\)", a, b, y, a, y-b, x),
	}
}

type funcEvaluateGen struct{}

func (g *funcEvaluateGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(5) + 2
	x := rand.Intn(5) + 2
	b := rand.Intn(5) + 1
	c := a*x + b
	return generator.Problem{
		Question:    fmt.Sprintf("\\(f(x) = %dx + %d\\). Find \\(f(%d)\\).", a, b, x),
		Answer:      fmt.Sprintf("%d", c),
		Explanation: fmt.Sprintf("\\(f(%d) = %d(%d) + %d = %d\\)", x, a, x, b, c),
	}
}

type funcLinearGen struct{}

func (g *funcLinearGen) Generate(difficulty float64) generator.Problem {
	if rand.Intn(2) == 0 {
		m := rand.Intn(6) - 3
		if m == 0 {
			m = 1
		}
		b := rand.Intn(10) - 5
		return generator.Problem{
			Question:    fmt.Sprintf("Is \\(y = %s\\) a function? (yes/no)", formatLinear(m, b)),
			Answer:      "yes",
			Explanation: fmt.Sprintf("Linear equations always define functions (passes vertical line test)."),
		}
	}
	c := rand.Intn(5) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("Is \\(x = %d\\) a function? (yes/no)", c),
		Answer:      "no",
		Explanation: fmt.Sprintf("x = %d is a vertical line — it fails the vertical line test.", c),
	}
}

type funcQuadGen struct{}

func (g *funcQuadGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(4) + 1
	b := rand.Intn(5) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("\\(f(x) = %dx^{2} + %d\\). Find \\(f(2)\\).", a, b),
		Answer:      fmt.Sprintf("%d", 4*a+b),
		Explanation: fmt.Sprintf("\\(f(2) = %d(4) + %d = %d + %d = %d\\)", a, b, 4*a, b, 4*a+b),
	}
}

// ---------------------------------------------------------------------------
// Exponential & Logarithmic
// ---------------------------------------------------------------------------

type algExpConceptGen struct{}

func (g *algExpConceptGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(3) + 2
	x := rand.Intn(4) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("If \\(f(x) = %d^{x}\\), find \\(f(%d)\\).", a, x),
		Answer:      fmt.Sprintf("%d", mathutil.IntPow(a, x)),
		Explanation: fmt.Sprintf("\\(%d^{%d} = %d\\)", a, x, mathutil.IntPow(a, x)),
	}
}

type algExpEvaluateGen struct{}

func (g *algExpEvaluateGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(3) + 2
	x := rand.Intn(4) + 2
	return generator.Problem{
		Question:    fmt.Sprintf("Evaluate: \\(%d^{%d}\\)", a, x),
		Answer:      fmt.Sprintf("%d", mathutil.IntPow(a, x)),
		Explanation: fmt.Sprintf("\\(%d^{%d} = %d\\)", a, x, mathutil.IntPow(a, x)),
	}
}

type logConceptGen struct{}

func (g *logConceptGen) Generate(difficulty float64) generator.Problem {
	base := rand.Intn(3) + 2
	exp := rand.Intn(4) + 1
	val := mathutil.IntPow(base, exp)
	return generator.Problem{
		Question:    fmt.Sprintf("Write as a logarithm: \\(%d^{%d} = %d\\)", base, exp, val),
		Answer:      fmt.Sprintf("log_%d(%d) = %d", base, val, exp),
		Explanation: fmt.Sprintf("\\(\\log_{%d}(%d) = %d\\)", base, val, exp),
	}
}

type logEvaluateGen struct{}

func (g *logEvaluateGen) Generate(difficulty float64) generator.Problem {
	base := rand.Intn(3) + 2
	exp := rand.Intn(4) + 1
	val := mathutil.IntPow(base, exp)
	return generator.Problem{
		Question:    fmt.Sprintf("Evaluate: \\(\\log_{%d}(%d)\\)", base, val),
		Answer:      fmt.Sprintf("%d", exp),
		Explanation: fmt.Sprintf("\\(%d^{%d} = %d\\), so \\(\\log_{%d}(%d) = %d\\)", base, exp, val, base, val, exp),
	}
}

type logPropertiesGen struct{}

func (g *logPropertiesGen) Generate(difficulty float64) generator.Problem {
	base := rand.Intn(3) + 2
	a := mathutil.IntPow(base, rand.Intn(3)+1)
	b := mathutil.IntPow(base, rand.Intn(3)+2)
	return generator.Problem{
		Question:    fmt.Sprintf("Simplify: \\(\\log_{%d}(%d \\times %d)\\)", base, a, b),
		Answer:      fmt.Sprintf("log_%d(%d) + log_%d(%d)", base, a, base, b),
		Explanation: fmt.Sprintf("log_b(xy) = log_b(x) + log_b(y)."),
	}
}

// ---------------------------------------------------------------------------
// Sequences
// ---------------------------------------------------------------------------

type seqArithGen struct{}

func (g *seqArithGen) Generate(difficulty float64) generator.Problem {
	a1 := rand.Intn(10) + 1
	d := rand.Intn(5) + 2
	n := rand.Intn(5) + 3
	an := a1 + (n-1)*d
	return generator.Problem{
		Question:    fmt.Sprintf("Arithmetic sequence: \\(%d, %d, %d, \\ldots\\) Find term \\(%d\\).", a1, a1+d, a1+2*d, n),
		Answer:      fmt.Sprintf("%d", an),
		Explanation: fmt.Sprintf("\\(a_{n} = a_{1} + (n-1)d = %d + (%d)(%d) = %d\\)", a1, n-1, d, an),
	}
}

type seqGeomGen struct{}

func (g *seqGeomGen) Generate(difficulty float64) generator.Problem {
	a1 := rand.Intn(5) + 2
	r := rand.Intn(3) + 2
	n := rand.Intn(4) + 2
	an := a1 * mathutil.IntPow(r, n-1)
	return generator.Problem{
		Question:    fmt.Sprintf("Geometric sequence: \\(%d, %d, %d, \\ldots\\) Find term \\(%d\\).", a1, a1*r, a1*r*r, n),
		Answer:      fmt.Sprintf("%d", an),
		Explanation: fmt.Sprintf("\\(a_{n} = a_{1} \\times r^{n-1} = %d \\times %d^{%d} = %d\\)", a1, r, n-1, an),
	}
}

type seqSumArithGen struct{}

func (g *seqSumArithGen) Generate(difficulty float64) generator.Problem {
	a1 := rand.Intn(10) + 1
	d := rand.Intn(5) + 2
	n := rand.Intn(6) + 2
	an := a1 + (n-1)*d
	sum := n * (a1 + an) / 2
	return generator.Problem{
		Question:    fmt.Sprintf("Find sum of arithmetic series: \\(%d + %d + %d + \\cdots\\) (first \\(%d\\) terms).", a1, a1+d, a1+2*d, n),
		Answer:      fmt.Sprintf("%d", sum),
		Explanation: fmt.Sprintf("\\(S_{n} = \\frac{n(a_{1}+a_{n})}{2} = \\frac{%d(%d+%d)}{2} = %d\\)", n, a1, an, sum),
	}
}

type seqSumGeoGen struct{}

func (g *seqSumGeoGen) Generate(difficulty float64) generator.Problem {
	a1 := rand.Intn(5) + 2
	r := rand.Intn(3) + 2
	// Sum of first n terms: a1(1-r^n)/(1-r)
	// Keep values small
	n := rand.Intn(4) + 2
	rn := mathutil.IntPow(r, n)
	sum := a1 * (1 - rn) / (1 - r)
	return generator.Problem{
		Question:    fmt.Sprintf("Find sum of geometric series: \\(%d + %d + %d + \\cdots\\) (first \\(%d\\) terms, \\(r=%d\\)).", a1, a1*r, a1*r*r, n, r),
		Answer:      fmt.Sprintf("%d", sum),
		Explanation: fmt.Sprintf("\\(S_{n} = \\frac{a_{1}(1-r^{n})}{1-r} = \\frac{%d(1-%d)}{1-%d} = %d\\)", a1, rn, r, sum),
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

func (g *ineqTwoVarGen) Generate(difficulty float64) generator.Problem {
	m := rand.Intn(4) + 1
	b := rand.Intn(6) - 3
	op := ">"
	if rand.Intn(2) == 0 {
		op = "<"
	}
	testX := rand.Intn(5) + 1
	testY := m*testX + b + rand.Intn(4) - 2
	expected := m*testX + b
	satisfies := testY > expected
	if op == "<" {
		satisfies = testY < expected
	}
	ans := "no"
	if satisfies {
		ans = "yes"
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Is \\((%d,%d)\\) a solution to \\(y %s %dx + %d\\)?", testX, testY, op, m, b),
		Answer:      ans,
		Explanation: fmt.Sprintf("At x=%d: y should be %s %d. y=%d %s %d = %t.", testX, op, expected, testY, op, expected, satisfies),
	}
}

type conicCircleGen struct{}

func (g *conicCircleGen) Generate(difficulty float64) generator.Problem {
	h := rand.Intn(5) - 2
	k := rand.Intn(5) - 2
	r := rand.Intn(4) + 2
	r2 := r * r
	return generator.Problem{
		Question:    fmt.Sprintf("What is the radius of \\((x%+d)^{2} + (y%+d)^{2} = %d\\)?", -h, -k, r2),
		Answer:      fmt.Sprintf("%d", r),
		Explanation: fmt.Sprintf("The radius is \\(\\sqrt{%d} = %d\\)", r2, r),
	}
}

type conicEllipseGen struct{}

func (g *conicEllipseGen) Generate(difficulty float64) generator.Problem {
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
	}
}

// ============= conic =============

type conicParabolaGen struct{}

func (g *conicParabolaGen) Generate(difficulty float64) generator.Problem {
	h := rand.Intn(5) - 2
	k := rand.Intn(5) - 2
	p := rand.Intn(3) + 1
	if rand.Intn(2) == 0 { p = -p }
	if rand.Intn(2) == 0 {
		return generator.Problem{
			Question:    fmt.Sprintf("Parabola: \\((x%+d)^{2} = %d(y%+d)\\). Which direction does it open?", -h, 4*p, -k),
			Answer:      map[bool]string{true: "up", false: "down"}[p > 0],
			Explanation: fmt.Sprintf("If 4p > 0, opens up; if 4p < 0, opens down. Here 4p = %d.", 4*p),
		}
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Parabola: \\((y%+d)^{2} = %d(x%+d)\\). Which direction does it open?", -k, 4*p, -h),
		Answer:      map[bool]string{true: "right", false: "left"}[p > 0],
		Explanation: fmt.Sprintf("If 4p > 0, opens right; if 4p < 0, opens left. Here 4p = %d.", 4*p),
	}
}

type conicHyperbolaGen struct{}

func (g *conicHyperbolaGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(3) + 2
	b := rand.Intn(2) + 1
	h := rand.Intn(3) - 1
	k := rand.Intn(3) - 1
	if rand.Intn(2) == 0 {
		return generator.Problem{
			Question:    fmt.Sprintf("Hyperbola: \\(\\frac{(x%+d)^{2}}{%d} - \\frac{(y%+d)^{2}}{%d} = 1\\). Which axis is transverse?", -h, a*a, -k, b*b),
			Answer:      "horizontal (x-axis)",
			Explanation: fmt.Sprintf("x term positive → horizontal transverse axis."),
		}
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Hyperbola: \\(\\frac{(y%+d)^{2}}{%d} - \\frac{(x%+d)^{2}}{%d} = 1\\). Which axis is transverse?", -k, a*a, -h, b*b),
		Answer:      "vertical (y-axis)",
		Explanation: fmt.Sprintf("y term positive → vertical transverse axis."),
	}
}

// ============= equations =============

type eqAbsValGen struct{}

func (g *eqAbsValGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(8) + 1
	b := rand.Intn(10) - 5
	return generator.Problem{
		Question:    fmt.Sprintf("Solve: \\(|x + %d| = %d\\)", b, a),
		Answer:      fmt.Sprintf("%d,%d", a-b, -a-b),
		Explanation: fmt.Sprintf("\\(x + %d = %d\\) or \\(x + %d = -%d\\) → \\(x = %d\\) or \\(x = %d\\)", b, a, b, a, a-b, -a-b),
	}
}

type eqBinomialGen struct{}

func (g *eqBinomialGen) Generate(difficulty float64) generator.Problem {
	r := rand.Intn(6) + 1
	k := rand.Intn(4) + 2
	rhs := mathutil.IntPow(k, r)
	return generator.Problem{
		Question:    fmt.Sprintf("Solve: \\(x^{%d} = %d\\)", r, rhs),
		Answer:      fmt.Sprintf("%d", k),
		Explanation: fmt.Sprintf("\\(x^{%d} = %d\\) → \\(x = %d^{\\frac{1}{%d}} = %d\\)", r, rhs, rhs, r, k),
	}
}

type eqExpGen struct{}

func (g *eqExpGen) Generate(difficulty float64) generator.Problem {
	b := rand.Intn(4) + 2
	e := rand.Intn(4) + 1
	p := rand.Intn(3) + 1
	rhs := mathutil.IntPow(b, e+p)
	return generator.Problem{
		Question:    fmt.Sprintf("Solve: \\(%d^{x-%d} = %d\\)", b, p, rhs),
		Answer:      fmt.Sprintf("%d", e+p),
		Explanation: fmt.Sprintf("\\(%d^{x-%d} = %d^{%d}\\) → \\(x-%d = %d\\) → \\(x = %d\\)", b, p, b, e+p, p, e, e+p),
	}
}

type extraneousRootsGen struct{}

func (g *extraneousRootsGen) Generate(difficulty float64) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"Why must you check solutions when squaring both sides of an equation?", "squaring can introduce extraneous roots", "Squaring is not reversible: if a=b then a²=b², but not vice versa."},
		{"Solve \\(\\sqrt{x} = -2\\). How many solutions?", "0", "The principal square root is always \\(\\geq 0\\), so \\(\\sqrt{x} = -2\\) has no solution."},
		{"When solving \\(\\sqrt{x+3} = x-3\\), what should you check after finding candidates?", "substitute back into original equation", "Extraneous roots satisfy the squared equation but not the original."},
		{"Multiplying both sides by a variable expression can introduce extraneous roots. (true/false)", "true", "If you multiply by something that equals 0 for some x, you may introduce extraneous solutions."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question, Answer: e.answer, Explanation: e.exp,
	}
}

type eqIrrationalGen struct{}

func (g *eqIrrationalGen) Generate(difficulty float64) generator.Problem {
	x := rand.Intn(8) + 2
	k := x * x
	return generator.Problem{
		Question:    fmt.Sprintf("Solve: \\(\\sqrt{x} = %d\\)", x),
		Answer:      fmt.Sprintf("%d", k),
		Explanation: fmt.Sprintf("Square both sides: \\(x = %d^{2} = %d\\)", x, k),
	}
}

type eqLogGen struct{}

func (g *eqLogGen) Generate(difficulty float64) generator.Problem {
	b := rand.Intn(3) + 2
	e := rand.Intn(3) + 1
	v := mathutil.IntPow(b, e)
	return generator.Problem{
		Question:    fmt.Sprintf("Solve: \\(\\log_{%d}(x) = %d\\)", b, e),
		Answer:      fmt.Sprintf("%d", v),
		Explanation: fmt.Sprintf("\\(\\log_{%d}(x) = %d\\) → \\(x = %d^{%d} = %d\\)", b, e, b, e, v),
	}
}

type eqPolyGen struct{}

func (g *eqPolyGen) Generate(difficulty float64) generator.Problem {
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
	}
}

type eqRationalGen struct{}

func (g *eqRationalGen) Generate(difficulty float64) generator.Problem {
	x := rand.Intn(8) + 2
	a := rand.Intn(5) + 1
	b := rand.Intn(5) + 1
	rhs := a * x / b
	return generator.Problem{
		Question:    fmt.Sprintf("Solve: \\(%dx/%d = %d\\)", a, b, rhs),
		Answer:      fmt.Sprintf("%d", x),
		Explanation: fmt.Sprintf("Multiply both sides by %d: \\(%dx = %d\\) → \\(x = %d\\)", b, a, rhs*b, x),
	}
}

type eqTrinomialGen struct{}

func (g *eqTrinomialGen) Generate(difficulty float64) generator.Problem {
	r := rand.Intn(6) + 2
	// x^2 + 2r x + r^2 = 0
	return generator.Problem{
		Question:    fmt.Sprintf("Solve: \\(x^{2} + %dx + %d = 0\\)", 2*r, r*r),
		Answer:      fmt.Sprintf("%d", -r),
		Explanation: fmt.Sprintf("\\((x + %d)^{2} = 0\\) → \\(x = %d\\)", r, -r),
	}
}

// ============= functions =============

type funcAbsValGen struct{}

func (g *funcAbsValGen) Generate(difficulty float64) generator.Problem {
	x := rand.Intn(10) - 5
	return generator.Problem{
		Question:    fmt.Sprintf("If \\(f(x) = |x|\\), what is \\(f(%d)\\)?", x),
		Answer:      fmt.Sprintf("%d", mathutil.Abs(x)),
		Explanation: fmt.Sprintf("\\(|%d| = %d\\)", x, mathutil.Abs(x)),
	}
}

type funcCompositeGen struct{}

func (g *funcCompositeGen) Generate(difficulty float64) generator.Problem {
	x := rand.Intn(5) + 1
	a := rand.Intn(4) + 1
	b := rand.Intn(4) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("If \\(f(x) = %dx\\) and \\(g(x) = x + %d\\), what is \\(f(g(%d))\\)?", a, b, x),
		Answer:      fmt.Sprintf("%d", a*(x+b)),
		Explanation: fmt.Sprintf("\\(g(%d) = %d+%d = %d\\). \\(f(g(%d)) = f(%d) = %d \\times %d = %d\\)", x, x, b, x+b, x, x+b, a, x+b, a*(x+b)),
	}
}

type funcDirichletGen struct{}

func (g *funcDirichletGen) Generate(difficulty float64) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"What is the Dirichlet function?", "D(x) = 1 if x is rational, 0 if x is irrational", "The Dirichlet function is nowhere continuous."},
		{"Is the Dirichlet function continuous anywhere?", "no", "The Dirichlet function is discontinuous at every point."},
		{"Is the Dirichlet function periodic?", "yes", "The Dirichlet function is periodic with any rational period."},
		{"Is the Dirichlet function Riemann integrable?", "no", "The Dirichlet function is not Riemann integrable on any interval."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question, Answer: e.answer, Explanation: e.exp,
	}
}

type funcDomainGen struct{}

func (g *funcDomainGen) Generate(difficulty float64) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"What is the domain of \\(f(x) = 1/x\\)?", "\\(x \\neq 0\\)", "Division by zero is undefined, so \\(x\\) cannot be \\(0\\)."},
		{"What is the domain of \\(f(x) = \\sqrt{x}\\)?", "\\(x \\geq 0\\)", "The square root of a negative number is not real."},
		{"What is the domain of \\(f(x) = \\ln(x)\\)?", "\\(x > 0\\)", "The natural log of zero or negative numbers is undefined."},
		{"What is the domain of \\(f(x) = 1/(x-2)\\)?", "\\(x \\neq 2\\)", "The denominator is zero at \\(x = 2\\)."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question, Answer: e.answer, Explanation: e.exp,
	}
}

type funcEvenOddGen struct{}

func (g *funcEvenOddGen) Generate(difficulty float64) generator.Problem {
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
		Question:    e.question, Answer: e.answer, Explanation: e.exp,
	}
}

type graphAnalysisGen struct{}

func (g *graphAnalysisGen) Generate(difficulty float64) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"What does the y-intercept tell you about a function?", "the output when x = 0", "The y-intercept is f(0), the value at x = 0."},
		{"On a graph, what are the x-intercepts?", "where f(x) = 0", "X-intercepts are the real roots/solutions of f(x) = 0."},
		{"What does f'(x) > 0 on an interval tell you?", "f is increasing", "A positive derivative means the function is increasing."},
		{"What does a vertical asymptote indicate?", "the function approaches ±∞ at that x", "As x → a, f(x) → ±∞ at a vertical asymptote x = a."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question, Answer: e.answer, Explanation: e.exp,
	}
}

type funcInverseGen struct{}

func (g *funcInverseGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(5) + 2
	b := rand.Intn(10) - 5
	return generator.Problem{
		Question:    fmt.Sprintf("If \\(f(x) = %dx + %d\\), what is \\(f^{-1}(x)\\)?", a, b),
		Answer:      fmt.Sprintf("(x %+d)/%d", -b, a),
		Explanation: fmt.Sprintf("\\(y = %dx + %d\\) → \\(x = \\frac{y %+d}{%d}\\) → \\(f^{-1}(x) = \\frac{x %+d}{%d}\\)", a, b, -b, a, -b, a),
	}
}

type monotonicityGen struct{}

func (g *monotonicityGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(5) + 1
	if rand.Intn(2) == 0 { a = -a }
	return generator.Problem{
		Question:    fmt.Sprintf("Is \\(f(x) = %dx\\) increasing or decreasing on \\(\\mathbb{R}\\)?", a),
		Answer:      map[bool]string{true: "increasing", false: "decreasing"}[a > 0],
		Explanation: fmt.Sprintf("Slope = \\(%d\\). %s slope means the function is %s.", a, map[bool]string{true: "Positive", false: "Negative"}[a > 0], map[bool]string{true: "increasing", false: "decreasing"}[a > 0]),
	}
}

type funcRationalGen struct{}

func (g *funcRationalGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(4) + 1
	b := rand.Intn(4) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("What is the vertical asymptote of \\(f(x) = \\frac{%d}{x-%d}\\)?", a, b),
		Answer:      fmt.Sprintf("%d", b),
		Explanation: fmt.Sprintf("The denominator is zero at x = %d, so there is a vertical asymptote at x = %d.", b, b),
	}
}

type funcSigmoidGen struct{}

func (g *funcSigmoidGen) Generate(difficulty float64) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"What is the range of the sigmoid function \\(\\sigma(x) = 1/(1+e^{-x})\\)?", "\\((0, 1)\\)", "The sigmoid maps all real numbers to the open interval \\((0,1)\\)."},
		{"What is \\(\\sigma(0)\\) for the sigmoid function?", "0.5", "\\(\\sigma(0) = 1/(1+e^{0}) = 1/2 = 0.5\\)."},
		{"The sigmoid function is commonly used as an ____ function in neural networks.", "activation", "The sigmoid is a popular activation function that introduces non-linearity."},
		{"As \\(x \\to \\infty\\), \\(\\sigma(x)\\) approaches ____.", "1", "As \\(x \\to \\infty\\), \\(e^{-x} \\to 0\\), so \\(\\sigma(x) \\to 1\\)."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question, Answer: e.answer, Explanation: e.exp,
	}
}

type funcSignGen struct{}

func (g *funcSignGen) Generate(difficulty float64) generator.Problem {
	x := rand.Intn(10) - 5
	if x == 0 { x = 3 }
	ans := "positive"
	if x < 0 { ans = "negative" }
	return generator.Problem{
		Question:    fmt.Sprintf("What is the sign of \\(f(x) = x\\) at \\(x = %d\\)?", x),
		Answer:      ans,
		Explanation: fmt.Sprintf("At x = %d, f(x) = %d, which is %s.", x, x, ans),
	}
}

// ============= inequalities =============

type ineqAbsValGen struct{}

func (g *ineqAbsValGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(5) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("Solve: \\(|x| < %d\\)", a),
		Answer:      fmt.Sprintf("-%d < x < %d", a, a),
		Explanation: fmt.Sprintf("\\(|x| < %d\\) means \\(-%d < x < %d\\)", a, a, a),
	}
}

type ineqIntervalGen struct{}

func (g *ineqIntervalGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(5) - 3
	b := a + rand.Intn(5) + 2
	types := []struct {
		q string
		a string
	}{
		{fmt.Sprintf("Write in interval notation: \\(%d \\leq x < %d\\).", a, b), fmt.Sprintf("[%d,%d)", a, b)},
		{fmt.Sprintf("Write in interval notation: \\(x > %d\\).", a), fmt.Sprintf("(%d,\\infty)", a)},
		{fmt.Sprintf("Write in interval notation: \\(x \\leq %d\\).", b), fmt.Sprintf("(-\\infty,%d]", b)},
		{fmt.Sprintf("Write in interval notation: \\(%d < x < %d\\).", a, b), fmt.Sprintf("(%d,%d)", a, b)},
	}
	e := types[rand.Intn(len(types))]
	return generator.Problem{
		Question: e.q, Answer: e.a,
		Explanation: fmt.Sprintf("The interval notation for this set is %s.", e.a),
	}
}

type ineqIrrationalGen struct{}

func (g *ineqIrrationalGen) Generate(difficulty float64) generator.Problem {
	x := rand.Intn(8) + 2
	return generator.Problem{
		Question:    fmt.Sprintf("Solve: \\(\\sqrt{x} > %d\\)", x),
		Answer:      fmt.Sprintf("x > %d", x*x),
		Explanation: fmt.Sprintf("\\(\\sqrt{x} > %d\\) → \\(x > %d^{2} = %d\\) (x ≥ 0 implied)", x, x, x*x),
	}
}

type ineqLogGen struct{}

func (g *ineqLogGen) Generate(difficulty float64) generator.Problem {
	b := rand.Intn(3) + 2
	e := rand.Intn(3) + 1
	v := mathutil.IntPow(b, e)
	return generator.Problem{
		Question:    fmt.Sprintf("Solve: \\(\\log_{%d}(x) > %d\\)", b, e),
		Answer:      fmt.Sprintf("x > %d", v),
		Explanation: fmt.Sprintf("\\(\\log_{%d}(x) > %d\\) → \\(x > %d^{%d} = %d\\) (base > 1 preserves inequality)", b, e, b, e, v),
	}
}

type ineqQuadraticGen struct{}

func (g *ineqQuadraticGen) Generate(difficulty float64) generator.Problem {
	r := rand.Intn(5) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("Solve: \\(x^{2} - %d < 0\\)", r*r),
		Answer:      fmt.Sprintf("-%d < x < %d", r, r),
		Explanation: fmt.Sprintf("\\(x^{2} - %d < 0\\) → \\((x-%d)(x+%d) < 0\\) → \\(-%d < x < %d\\)", r*r, r, r, r, r),
	}
}

type ineqRationalGen struct{}

func (g *ineqRationalGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(5) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("Solve: \\(\\frac{1}{x-%d} > 0\\)", a),
		Answer:      fmt.Sprintf("x > %d", a),
		Explanation: fmt.Sprintf("numerator is always positive 1, so \\(\\frac{1}{x-%d} > 0\\) when \\(x-%d > 0\\) → \\(x > %d\\)", a, a, a),
	}
}

type signAnalysisGen struct{}

func (g *signAnalysisGen) Generate(difficulty float64) generator.Problem {
	r := rand.Intn(5) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("Analyze the sign of \\(f(x) = (x+%d)(x-%d)\\) for \\(x < -%d\\).", r, r, r),
		Answer:      "positive",
		Explanation: fmt.Sprintf("For \\(x < -%d\\): both \\((x+%d)\\) and \\((x-%d)\\) are negative, product is positive.", r, r, r),
	}
}

type ineqSystemsGen struct{}

func (g *ineqSystemsGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(4) + 2
	b := rand.Intn(5) - 2
	x := rand.Intn(4) + 1
	y := a*x + b + rand.Intn(3)
	satisfies := y > a*x + b
	ans := map[bool]string{true: "yes", false: "no"}[satisfies]
	return generator.Problem{
		Question:    fmt.Sprintf("Is \\((%d,%d)\\) a solution to \\(y > %dx + %d\\)?", x, y, a, b),
		Answer:      ans,
		Explanation: fmt.Sprintf("At \\(x=%d\\): RHS = \\(%d(%d)+%d = %d\\). \\(y=%d %s %d\\), so %s.", x, a, x, b, a*x+b, y, map[bool]string{true: ">", false: "≤"}[satisfies], a*x+b, ans),
	}
}

// ============= polynomials =============

type polyDivisionGen struct{}

func (g *polyDivisionGen) Generate(difficulty float64) generator.Problem {
	r := rand.Intn(6) + 1
	a := rand.Intn(4) + 1
	// (ax^2 + r*x) / x = ax + r
	return generator.Problem{
		Question:    fmt.Sprintf("Simplify: \\((%dx^{2} + %dx) / x\\)", a, r*a),
		Answer:      fmt.Sprintf("%dx + %d", a, r*a),
		Explanation: fmt.Sprintf("Divide each term: \\(%dx^{2}/x = %dx\\), \\(%dx/x = %d\\)", a, a, r*a, r*a),
	}
}

type polyMonomialGen struct{}

func (g *polyMonomialGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(5) + 1
	n := rand.Intn(4) + 2
	return generator.Problem{
		Question:    fmt.Sprintf("What is the degree of \\(%dx^{%d}\\)?", a, n),
		Answer:      fmt.Sprintf("%d", n),
		Explanation: fmt.Sprintf("The monomial \\(%dx^{%d}\\) has degree %d.", a, n, n),
	}
}

type polyRootsGen struct{}

func (g *polyRootsGen) Generate(difficulty float64) generator.Problem {
	r1 := rand.Intn(6) - 3
	r2 := rand.Intn(6) - 3
	b := -(r1 + r2)
	c := r1 * r2
	return generator.Problem{
		Question:    fmt.Sprintf("Find the roots of \\(x^{2} + %dx + %d = 0\\).", b, c),
		Answer:      fmt.Sprintf("%d,%d", r1, r2),
		Explanation: fmt.Sprintf("\\((x %+d)(x %+d) = 0\\) → roots are \\(%d\\) and \\(%d\\)", -r1, -r2, r1, r2),
	}
}

type synthDivGen struct{}

func (g *synthDivGen) Generate(difficulty float64) generator.Problem {
	r := rand.Intn(6) - 3
	a := rand.Intn(4) + 1
	b := rand.Intn(5) - 2
	// (ax^2 + bx - ar) / (x - r) = ax + (ar + b)
	return generator.Problem{
		Question:    fmt.Sprintf("Use synthetic division to divide \\((%dx^{2} + %dx)\\) by \\((x - %d)\\). What is the quotient?", a, b, r),
		Answer:      fmt.Sprintf("%dx + %d", a, b+a*r),
		Explanation: fmt.Sprintf("Synthetic division by \\(%d\\) gives coefficients %d and %d → \\(%dx + %d\\)", r, a, b+a*r, a, b+a*r),
	}
}

type vietaGen struct{}

func (g *vietaGen) Generate(difficulty float64) generator.Problem {
	r1 := rand.Intn(6) + 1
	r2 := rand.Intn(6) + 1
	b := -(r1 + r2)
	c := r1 * r2
	return generator.Problem{
		Question:    fmt.Sprintf("For \\(x^{2} + %dx + %d = 0\\), what is the sum of the roots?", b, c),
		Answer:      fmt.Sprintf("%d", r1+r2),
		Explanation: fmt.Sprintf("By Vieta: sum of roots = \\(%d\\). Roots are \\(%d\\) and \\(%d\\), sum = \\(%d\\)", -b, r1, r2, r1+r2),
	}
}

// ============= quadratics =============

type quadComplexGen struct{}

func (g *quadComplexGen) Generate(difficulty float64) generator.Problem {
	b := rand.Intn(6) + 2
	c := (b*b)/4 + 1
	return generator.Problem{
		Question:    fmt.Sprintf("How many real solutions does \\(x^{2} + %dx + %d = 0\\) have?", b, c),
		Answer:      "0",
		Explanation: fmt.Sprintf("Discriminant = \\(%d^{2} - 4(%d) = %d - %d = %d < 0\\) → two complex (non-real) solutions.", b, c, b*b, 4*c, b*b-4*c),
	}
}

type quadIncompleteGen struct{}

func (g *quadIncompleteGen) Generate(difficulty float64) generator.Problem {
	k := rand.Intn(5) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("Solve: \\(x^{2} - %d = 0\\)", k*k),
		Answer:      fmt.Sprintf("%d,%d", k, -k),
		Explanation: fmt.Sprintf("\\(x^{2} = %d\\) → \\(x = \\pm \\sqrt{%d} = \\pm %d\\)", k*k, k*k, k),
	}
}

type quadParametricGen struct{}

func (g *quadParametricGen) Generate(difficulty float64) generator.Problem {
	p := rand.Intn(5) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("For what value(s) of \\(k\\) does \\(x^{2} + %dx + %d = 0\\) have exactly one solution?", 2*p, p*p),
		Answer:      fmt.Sprintf("%d", p),
		Explanation: fmt.Sprintf("Discriminant = 0: \\((%d)^{2} - 4(%d) = %d - %d = 0\\) → \\(k = %d\\)", 2*p, p*p, 4*p*p, 4*p*p, p),
	}
}

type quadQuadraticGen struct{}

func (g *quadQuadraticGen) Generate(difficulty float64) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"What is the standard form of a quadratic equation?", "\\(ax^{2} + bx + c = 0\\)", "Standard form: \\(ax^{2} + bx + c = 0\\) where \\(a \\neq 0\\)."},
		{"What is the quadratic formula?", "\\(x = [-b \\pm \\sqrt{b^{2}-4ac}]/(2a)\\)", "The quadratic formula solves \\(ax^{2} + bx + c = 0\\)."},
		{"What does the discriminant determine?", "the number and type of solutions", "\\(b^{2}-4ac > 0\\): two real; \\(= 0\\): one real; \\(< 0\\): two complex."},
		{"The graph of a quadratic function is a ____.", "parabola", "\\(f(x) = ax^{2} + bx + c\\) graphs as a parabola."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question: e.question, Answer: e.answer, Explanation: e.exp,
	}
}

// ============= systems =============

type systemsConceptGen struct{}

func (g *systemsConceptGen) Generate(difficulty float64) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"What does it mean to solve a system of equations?", "find values satisfying all equations simultaneously", "The solution makes every equation in the system true."},
		{"How many solutions can a system of two linear equations have?", "0, 1, or infinitely many", "Two lines can intersect once, be parallel (0), or be the same line (∞)."},
		{"What is a consistent system?", "a system with at least one solution", "Consistent: at least one solution exists."},
		{"What is an inconsistent system?", "a system with no solution", "Inconsistent: parallel lines that never intersect."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question: e.question, Answer: e.answer, Explanation: e.exp,
	}
}

type gaussianElimGen struct{}

func (g *gaussianElimGen) Generate(difficulty float64) generator.Problem {
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
		}
	}
	// Row echelon form question
	return generator.Problem{
		Question:    fmt.Sprintf("Put the system \\(%dx + %dy = %d\\), \\(%dx + %dy = %d\\) into row echelon form.", a, b, e1, c, d, e2),
		Answer:      fmt.Sprintf("[[%d,%d,%d],[0,%d,%d]]", a, b, e1, d-c*b/a, e2-c*e1/a),
		Explanation: fmt.Sprintf("The augmented matrix [[%d,%d,%d],[%d,%d,%d]] can be reduced.", a, b, e1, c, d, e2),
	}
}
