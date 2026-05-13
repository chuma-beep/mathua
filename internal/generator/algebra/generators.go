package algebra

import (
	"fmt"
	"math/rand"

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
}

// ---------------------------------------------------------------------------
// Linear
// ---------------------------------------------------------------------------

type slopeGen struct{}

func (g *slopeGen) Generate(difficulty float64) generator.Problem {
	x1 := rand.Intn(10)
	y1 := rand.Intn(10)
	x2 := x1 + rand.Intn(8) + 1
	y2 := y1 + rand.Intn(10) + 1
	if rand.Intn(2) == 0 {
		y2 = y1 - rand.Intn(y1+1)
	}
	dy := y2 - y1
	dx := x2 - x1
	frac := reduce(dy, dx)
	num, den := frac.num, frac.den
	return generator.Problem{
		Question:    fmt.Sprintf("Find the slope between (%d,%d) and (%d,%d).", x1, y1, x2, y2),
		Answer:      fracOrInt(num, den),
		Explanation: fmt.Sprintf("m = (%d - %d)/(%d - %d) = %d/%d = %s.", y2, y1, x2, x1, dy, dx, fracOrInt(num, den)),
	}
}

type slopeInterceptGen struct{}

func (g *slopeInterceptGen) Generate(difficulty float64) generator.Problem {
	m := rand.Intn(5) + 1
	if rand.Intn(2) == 0 {
		m = -m
	}
	b := rand.Intn(10) - 5
	return generator.Problem{
		Question:    fmt.Sprintf("Write the equation of a line with slope %d and y-intercept %d (y = mx + b).", m, b),
		Answer:      formatLinear(m, b),
		Explanation: fmt.Sprintf("y = %dx + %d", m, b),
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
		Question:    fmt.Sprintf("Given y = %s, find y when x = %d.", formatLinear(m, b), x),
		Answer:      fmt.Sprintf("%d", y),
		Explanation: fmt.Sprintf("y = %d(%d) + %d = %d + %d = %d.", m, x, b, m*x, b, y),
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
			Question:    fmt.Sprintf("Is (%d,%d) a solution to %dx + %dy = %d?", x, y, a, b, c),
			Answer:      "yes",
			Explanation: fmt.Sprintf("%d(%d) + %d(%d) = %d + %d = %d. Yes!", a, x, b, y, a*x, b*y, c),
		}
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Is (%d,%d) a solution to %dx + %dy = %d?", x+1, y, a, b, c),
		Answer:      "no",
		Explanation: fmt.Sprintf("%d(%d) + %d(%d) = %d + %d = %d, not %d. No!", a, x+1, b, y, a*(x+1), b*y, a*(x+1)+b*y, c),
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
			Question:    fmt.Sprintf("A line parallel to y = %dx + 3 has slope:", m),
			Answer:      fmt.Sprintf("%d", m),
			Explanation: fmt.Sprintf("Parallel lines have equal slopes. Slope = %d.", m),
		}
	}
	perp := reduce(-1, m)
	return generator.Problem{
		Question:    fmt.Sprintf("A line perpendicular to y = %dx + 3 has slope:", m),
		Answer:      fracOrInt(perp.num, perp.den),
		Explanation: fmt.Sprintf("Perpendicular slopes are negative reciprocals: -1/%d = %s.", m, fracOrInt(perp.num, perp.den)),
	}
}

// ---------------------------------------------------------------------------
// Equations & inequalities
// ---------------------------------------------------------------------------

type multiStepEqGen struct{}

func (g *multiStepEqGen) Generate(difficulty float64) generator.Problem {
	x := rand.Intn(10) + 2
	a := rand.Intn(6) + 2
	b := rand.Intn(10) + 2
	c := a*x + b
	return generator.Problem{
		Question:    fmt.Sprintf("Solve: %dx + %d = %d", a, b, c),
		Answer:      fmt.Sprintf("%d", x),
		Explanation: fmt.Sprintf("%dx + %d = %d -> %dx = %d -> x = %d.", a, b, c, a, c-b, x),
	}
}

type varsBothSidesGen struct{}

func (g *varsBothSidesGen) Generate(difficulty float64) generator.Problem {
	x := rand.Intn(8) + 2
	a := rand.Intn(5) + 2
	c := rand.Intn(5) + 1
	for a == c {
		c++
	}
	b := rand.Intn(10) + 1
	d := (a-c)*x + b
	return generator.Problem{
		Question:    fmt.Sprintf("Solve: %dx + %d = %dx + %d", a, b, c, d),
		Answer:      fmt.Sprintf("%d", x),
		Explanation: fmt.Sprintf("%dx + %d = %dx + %d -> %dx = %d -> x = %d.", a, b, c, d, a-c, d-b, x),
	}
}

type literalEqGen struct{}

func (g *literalEqGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(6) + 2
	b := rand.Intn(5) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("Solve for y: %dx + %dy = z", a, b),
		Answer:      fmt.Sprintf("y = (z - %dx)/%d", a, b),
		Explanation: fmt.Sprintf("%dy = z - %dx -> y = (z - %dx)/%d.", b, a, a, b),
	}
}

type multiStepIneqGen struct{}

func (g *multiStepIneqGen) Generate(difficulty float64) generator.Problem {
	x := rand.Intn(8) + 2
	a := rand.Intn(6) + 2
	b := rand.Intn(10) + 1
	c := a*(x-1) + b + rand.Intn(a)
	return generator.Problem{
		Question:    fmt.Sprintf("Solve: %dx + %d > %d", a, b, c),
		Answer:      fmt.Sprintf("%d", x),
		Explanation: fmt.Sprintf("%dx + %d > %d -> %dx > %d -> x > %d, so x >= %d.", a, b, c, a, c-b, c-b, x),
	}
}

type compoundIneqGen struct{}

func (g *compoundIneqGen) Generate(difficulty float64) generator.Problem {
	x := rand.Intn(8) + 2
	a := x - 1
	b := x + 2
	return generator.Problem{
		Question:    fmt.Sprintf("Solve: %d < x + 1 < %d", a, b),
		Answer:      fmt.Sprintf("%d", x),
		Explanation: fmt.Sprintf("%d < x + 1 < %d -> %d < x < %d, so x = %d.", a, b, a-1, b-1, x),
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
		Question:    fmt.Sprintf("Solve: y = %d,  x + y = %d", y, x+y),
		Answer:      fmt.Sprintf("(%d,%d)", x, y),
		Explanation: fmt.Sprintf("Substitute y=%d: x + %d = %d -> x = %d. Solution: (%d,%d).", y, y, x+y, x, x, y),
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
		Question:    fmt.Sprintf("Solve: %dx + %dy = %d,  %dx + %dy = %d", a, b, eq1, c, b, eq2),
		Answer:      fmt.Sprintf("(%d,%d)", x, y),
		Explanation: fmt.Sprintf("Subtract: %dx = %d -> x = %d; y = %d. Solution: (%d,%d).", a-c, eq1-eq2, x, y, x, y),
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
		Explanation: fmt.Sprintf("x + y = %d, x - y = %d. Add: 2x = %d -> x = %d. y = %d.", sum, diff, a+b, a, b),
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
		Question:    fmt.Sprintf("What is the degree of %dx^%d + 3x + 1?", coeff, exp),
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
		Question:    fmt.Sprintf("Simplify: (%dx^2 + %dx) + (%dx^2 + %dx)", a, c, b, c),
		Answer:      fmt.Sprintf("%dx^2 + %dx", a+b, 2*c),
		Explanation: fmt.Sprintf("%dx^2 + %dx^2 = %dx^2; %dx + %dx = %dx.", a, b, a+b, c, c, 2*c),
	}
}

type polyMultMonoGen struct{}

func (g *polyMultMonoGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(5) + 2
	b := rand.Intn(5) + 2
	c := rand.Intn(5) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("Simplify: %dx(%dx + %d)", a, b, c),
		Answer:      fmt.Sprintf("%dx^2 + %dx", a*b, a*c),
		Explanation: fmt.Sprintf("%dx(%dx) + %dx(%d) = %dx^2 + %dx.", a, b, a, c, a*b, a*c),
	}
}

type polyFoilGen struct{}

func (g *polyFoilGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(6) + 1
	b := rand.Intn(6) + 1
	c := rand.Intn(6) + 1
	d := rand.Intn(6) + 1
	return generator.Problem{
		Question: fmt.Sprintf("Simplify: (%dx + %d)(%dx + %d)", a, b, c, d),
		Answer:   fmt.Sprintf("%dx^2 + %dx + %d", a*c, a*d+b*c, b*d),
		Explanation: fmt.Sprintf("FOIL: %dx^2 + %dx + %dx + %d = %dx^2 + %dx + %d.",
			a*c, a*d, b*c, b*d, a*c, a*d+b*c, b*d),
	}
}

type polySpecialGen struct{}

func (g *polySpecialGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(6) + 1
	b := rand.Intn(5) + 1
	return generator.Problem{
		Question: fmt.Sprintf("Simplify: (%dx + %d)^2", a, b),
		Answer:   fmt.Sprintf("%dx^2 + %dx + %d", a*a, 2*a*b, b*b),
		Explanation: fmt.Sprintf("(%dx)^2 + 2(%dx)(%d) + %d^2 = %dx^2 + %dx + %d.",
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
		Question:    fmt.Sprintf("Factor: %dx + %d", a, b),
		Answer:      fmt.Sprintf("%d(%dx + %d)", f, a/f, b/f),
		Explanation: fmt.Sprintf("GCF is %d: %dx + %d = %d(%dx + %d).", f, a, b, f, a/f, b/f),
	}
}

type factorTrinomialGen struct{}

func (g *factorTrinomialGen) Generate(difficulty float64) generator.Problem {
	r1 := rand.Intn(6) + 1
	r2 := rand.Intn(6) + 1
	b := r1 + r2
	c := r1 * r2
	return generator.Problem{
		Question:    fmt.Sprintf("Factor: x^2 + %dx + %d", b, c),
		Answer:      fmt.Sprintf("(x + %d)(x + %d)", r1, r2),
		Explanation: fmt.Sprintf("Find two numbers that multiply to %d and add to %d: %d and %d.", c, b, r1, r2),
	}
}

type factorDiffSquaresGen struct{}

func (g *factorDiffSquaresGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(8) + 2
	return generator.Problem{
		Question:    fmt.Sprintf("Factor: x^2 - %d", a*a),
		Answer:      fmt.Sprintf("(x + %d)(x - %d)", a, a),
		Explanation: fmt.Sprintf("x^2 - %d = (x + %d)(x - %d).", a*a, a, a),
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
		Question:    fmt.Sprintf("Factor: x^2 + %dx + %d", b, c),
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
	return generator.Problem{
		Question:    fmt.Sprintf("Solve: x^2 + %dx + %d = 0", b, c),
		Answer:      fmt.Sprintf("%d,%d", r1, r2),
		Explanation: fmt.Sprintf("Factors: (x %+d)(x %+d) = 0, so x = %d or x = %d.", -r1, -r2, r1, r2),
	}
}

type quadCompleteSquareGen struct{}

func (g *quadCompleteSquareGen) Generate(difficulty float64) generator.Problem {
	r := rand.Intn(6) + 1
	b := -2 * r
	return generator.Problem{
		Question:    fmt.Sprintf("Complete the square: x^2 + %dx + __ to make a perfect square.", b),
		Answer:      fmt.Sprintf("%d", r*r),
		Explanation: fmt.Sprintf("(b/2)^2 = (%d/2)^2 = %d.", b, r*r),
	}
}

type quadFormulaGen struct{}

func (g *quadFormulaGen) Generate(difficulty float64) generator.Problem {
	r1 := rand.Intn(8) + 1
	r2 := rand.Intn(8) + 1
	b := -(r1 + r2)
	c := r1 * r2
	return generator.Problem{
		Question:    fmt.Sprintf("Solve using quadratic formula: x^2 + %dx + %d = 0", b, c),
		Answer:      fmt.Sprintf("%d,%d", r1, r2),
		Explanation: fmt.Sprintf("x = [-%d +/- sqrt(%d - 4(%d))]/2 = %d or %d.", b, b*b, c, r1, r2),
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
		Question:    fmt.Sprintf("How many real solutions? x^2 + %dx + %d = 0", b, c),
		Answer:      count,
		Explanation: fmt.Sprintf("Discriminant = %d^2 - 4(%d) = %d. %s solutions.", b, c, disc, count),
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
		Question:    fmt.Sprintf("If f(x) = %dx + %d, what is f(%d)?", a, b, x),
		Answer:      fmt.Sprintf("%d", a*x+b),
		Explanation: fmt.Sprintf("f(%d) = %d(%d) + %d = %d.", x, a, x, b, a*x+b),
	}
}

type funcNotationGen struct{}

func (g *funcNotationGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(5) + 2
	x := rand.Intn(5) + 1
	b := rand.Intn(5) + 1
	y := a*x + b
	return generator.Problem{
		Question:    fmt.Sprintf("f(x) = %dx + %d. Find x when f(x) = %d.", a, b, y),
		Answer:      fmt.Sprintf("%d", x),
		Explanation: fmt.Sprintf("%dx + %d = %d -> %dx = %d -> x = %d.", a, b, y, a, y-b, x),
	}
}

type funcEvaluateGen struct{}

func (g *funcEvaluateGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(5) + 2
	x := rand.Intn(5) + 2
	b := rand.Intn(5) + 1
	c := a*x + b
	return generator.Problem{
		Question:    fmt.Sprintf("f(x) = %dx + %d. Find f(%d).", a, b, x),
		Answer:      fmt.Sprintf("%d", c),
		Explanation: fmt.Sprintf("f(%d) = %d(%d) + %d = %d.", x, a, x, b, c),
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
			Question:    fmt.Sprintf("Is y = %s a function? (yes/no)", formatLinear(m, b)),
			Answer:      "yes",
			Explanation: fmt.Sprintf("Linear equations always define functions (passes vertical line test)."),
		}
	}
	c := rand.Intn(5) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("Is x = %d a function? (yes/no)", c),
		Answer:      "no",
		Explanation: fmt.Sprintf("x = %d is a vertical line — it fails the vertical line test.", c),
	}
}

type funcQuadGen struct{}

func (g *funcQuadGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(4) + 1
	b := rand.Intn(5) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("f(x) = %dx^2 + %d. Find f(2).", a, b),
		Answer:      fmt.Sprintf("%d", 4*a+b),
		Explanation: fmt.Sprintf("f(2) = %d(4) + %d = %d + %d = %d.", a, b, 4*a, b, 4*a+b),
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
		Question:    fmt.Sprintf("If f(x) = %d^x, find f(%d).", a, x),
		Answer:      fmt.Sprintf("%d", mathutil.IntPow(a, x)),
		Explanation: fmt.Sprintf("%d^%d = %d.", a, x, mathutil.IntPow(a, x)),
	}
}

type algExpEvaluateGen struct{}

func (g *algExpEvaluateGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(3) + 2
	x := rand.Intn(4) + 2
	return generator.Problem{
		Question:    fmt.Sprintf("Evaluate: %d^%d", a, x),
		Answer:      fmt.Sprintf("%d", mathutil.IntPow(a, x)),
		Explanation: fmt.Sprintf("%d^%d = %d.", a, x, mathutil.IntPow(a, x)),
	}
}

type logConceptGen struct{}

func (g *logConceptGen) Generate(difficulty float64) generator.Problem {
	base := rand.Intn(3) + 2
	exp := rand.Intn(4) + 1
	val := mathutil.IntPow(base, exp)
	return generator.Problem{
		Question:    fmt.Sprintf("Write as a logarithm: %d^%d = %d", base, exp, val),
		Answer:      fmt.Sprintf("log_%d(%d) = %d", base, val, exp),
		Explanation: fmt.Sprintf("log_%d(%d) = %d.", base, val, exp),
	}
}

type logEvaluateGen struct{}

func (g *logEvaluateGen) Generate(difficulty float64) generator.Problem {
	base := rand.Intn(3) + 2
	exp := rand.Intn(4) + 1
	val := mathutil.IntPow(base, exp)
	return generator.Problem{
		Question:    fmt.Sprintf("Evaluate: log_%d(%d)", base, val),
		Answer:      fmt.Sprintf("%d", exp),
		Explanation: fmt.Sprintf("%d^%d = %d, so log_%d(%d) = %d.", base, exp, val, base, val, exp),
	}
}

type logPropertiesGen struct{}

func (g *logPropertiesGen) Generate(difficulty float64) generator.Problem {
	base := rand.Intn(3) + 2
	a := mathutil.IntPow(base, rand.Intn(3)+1)
	b := mathutil.IntPow(base, rand.Intn(3)+2)
	return generator.Problem{
		Question:    fmt.Sprintf("Simplify: log_%d(%d x %d)", base, a, b),
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
		Question:    fmt.Sprintf("Arithmetic sequence: %d, %d, %d, ... Find term %d.", a1, a1+d, a1+2*d, n),
		Answer:      fmt.Sprintf("%d", an),
		Explanation: fmt.Sprintf("a_n = a1 + (n-1)d = %d + (%d)(%d) = %d.", a1, n-1, d, an),
	}
}

type seqGeomGen struct{}

func (g *seqGeomGen) Generate(difficulty float64) generator.Problem {
	a1 := rand.Intn(5) + 2
	r := rand.Intn(3) + 2
	n := rand.Intn(4) + 2
	an := a1 * mathutil.IntPow(r, n-1)
	return generator.Problem{
		Question:    fmt.Sprintf("Geometric sequence: %d, %d, %d, ... Find term %d.", a1, a1*r, a1*r*r, n),
		Answer:      fmt.Sprintf("%d", an),
		Explanation: fmt.Sprintf("a_n = a1 x r^(n-1) = %d x %d^%d = %d.", a1, r, n-1, an),
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
		Question:    fmt.Sprintf("Find sum of arithmetic series: %d + %d + %d + ... (first %d terms).", a1, a1+d, a1+2*d, n),
		Answer:      fmt.Sprintf("%d", sum),
		Explanation: fmt.Sprintf("S_n = n(a1+an)/2 = %d(%d+%d)/2 = %d.", n, a1, an, sum),
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
		Question:    fmt.Sprintf("Find sum of geometric series: %d + %d + %d + ... (first %d terms, r=%d).", a1, a1*r, a1*r*r, n, r),
		Answer:      fmt.Sprintf("%d", sum),
		Explanation: fmt.Sprintf("S_n = a1(1-r^n)/(1-r) = %d(1-%d)/(1-%d) = %d.", a1, rn, r, sum),
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
		Question:    fmt.Sprintf("Is (%d,%d) a solution to y %s %dx + %d?", testX, testY, op, m, b),
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
		Question:    fmt.Sprintf("What is the radius of (x%+d)² + (y%+d)² = %d?", -h, -k, r2),
		Answer:      fmt.Sprintf("%d", r),
		Explanation: fmt.Sprintf("The radius is √%d = %d.", r2, r),
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
		Question:    fmt.Sprintf("What is the %s axis length of (x%+d)²/%d + (y%+d)²/%d = 1?", label, -h, a2, -k, b2),
		Answer:      fmt.Sprintf("%d", val),
		Explanation: fmt.Sprintf("The %s axis length is %d because %s² = %d.", label, val, label[:6], val*val),
	}
}
