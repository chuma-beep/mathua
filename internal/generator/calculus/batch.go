package calculus

// The 2026-10-18 calculus batch: ten concepts that take the domain from 165 to
// the target of 175.
//
// The survey found the domain already dense: limits, every derivative and
// integral technique, sequences and series, multivariable and vector calculus.
// Several obvious additions were already taught and were not repeated: Taylor's
// theorem with remainder and the Weierstrass M-test in the series lessons,
// numerical-integration error bounds, the Hessian in the partial-derivatives
// lesson, hyperbolic integrals in the integration tables, and surface area of
// graphs and revolutions. The batch fills ten genuine gaps instead: hyperbolic
// derivatives, the second parametric derivative, finite differences, polar arc
// length, torsion, the Frenet frame, second-order multivariable Taylor,
// hydrostatic force, Pappus's centroid theorem, and moment of inertia.
//
// The wiring is deliberately scoped. Section 5.1 measured the sink-targeting
// policy and rejected it for calculus: the sink set is most of the domain, so
// calc1 and calc2 would stop being stages. The three courses therefore gain
// only the new concepts that belong to their stage, and the existing curated
// reach is left intact.
//
// Seven concepts are numeric, one is an expression, one is a symbolic multiple
// of pi, and one is a tuple. Every generator publishes the facts its answer
// follows from, so the tests recompute the answer from those facts.

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerBatch(reg *generator.Registry) {
	reg.Register("calc.deriv.hyperbolic", &hyperbolicDerivGen{})
	reg.Register("calc.deriv.parametric_second", &parametricSecondGen{})
	reg.Register("calc.deriv.finite_difference", &finiteDifferenceGen{})
	reg.Register("calc.polar.arc_length", &polarArcLengthGen{})
	reg.Register("calc.vec.torsion", &torsionGen{})
	reg.Register("calc.vec.frenet", &frenetGen{})
	reg.Register("calc.multi.taylor2", &taylor2Gen{})
	reg.Register("calc.integral.hydrostatic", &hydrostaticGen{})
	reg.Register("calc.integral.pappus", &pappusGen{})
	reg.Register("calc.integral.moment_inertia", &momentInertiaGen{})
}

// ---------------------------------------------------------------------------
// 1. Derivatives of hyperbolic functions (expression)
// ---------------------------------------------------------------------------

type hyperbolicDerivGen struct{}

func (g *hyperbolicDerivGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		rule     string
		reason   string
	}
	claims := []claim{
		{
			question: "Differentiate f(x) = sinh(x) with respect to x. (enter an expression)",
			answer:   "cosh(x)", rule: "sinh",
			reason: "The derivative of sinh is cosh: d/dx sinh(x) = cosh(x).",
		},
		{
			question: "Differentiate f(x) = cosh(x) with respect to x. (enter an expression)",
			answer:   "sinh(x)", rule: "cosh",
			reason: "The derivative of cosh is sinh, with no sign change: d/dx cosh(x) = sinh(x).",
		},
		{
			question: "Differentiate f(x) = tanh(x) with respect to x. (enter an expression)",
			answer:   "1/cosh(x)^2", rule: "tanh",
			reason: "d/dx tanh(x) = sech^2(x) = 1/cosh(x)^2.",
		},
		{
			question: "Differentiate f(x) = sinh(3x) with respect to x. (enter an expression)",
			answer:   "3*cosh(3*x)", rule: "sinh_chain",
			reason: "The chain rule gives cosh(3x) times the derivative 3 of the inner function.",
		},
		{
			question: "Differentiate f(x) = cosh(x^2) with respect to x. (enter an expression)",
			answer:   "2*x*sinh(x^2)", rule: "cosh_chain",
			reason: "The chain rule gives sinh(x^2) times the derivative 2x of the inner function.",
		},
		{
			question: "Differentiate f(x) = tanh(x/2) with respect to x. (enter an expression)",
			answer:   "(1/2)/cosh(x/2)^2", rule: "tanh_chain",
			reason: "The derivative is sech^2(x/2) times the inner derivative 1/2, giving (1/2)/cosh(x/2)^2.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	return generator.Problem{
		Question:    c.question,
		Answer:      c.answer,
		Explanation: c.reason,
		Facts:       map[string]string{"rule": c.rule, "answer": c.answer},
	}
}

// ---------------------------------------------------------------------------
// 2. The second derivative of a parametric curve (numeric)
// ---------------------------------------------------------------------------

type parametricSecondGen struct{}

func (g *parametricSecondGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		xp, yp   int
		xpp, ypp int
		reason   string
	}
	claims := []claim{
		{
			question: "At a point on a parametric curve, x' = 1, y' = 2, x'' = 0 and y'' = 2. What is d^2y/dx^2 = (x'y'' - y'x'')/(x')^3? (enter a number)",
			answer:   "2", xp: 1, yp: 2, xpp: 0, ypp: 2,
			reason: "The numerator is 1*2 - 2*0 = 2, and the denominator is 1, so the second derivative is 2.",
		},
		{
			question: "At a point on a parametric curve, x' = 2, y' = 3, x'' = 2 and y'' = 6. What is d^2y/dx^2 = (x'y'' - y'x'')/(x')^3? (enter a decimal)",
			answer:   "0.75", xp: 2, yp: 3, xpp: 2, ypp: 6,
			reason: "The numerator is 2*6 - 3*2 = 6, and the denominator is 8, so the second derivative is 0.75.",
		},
		{
			question: "At a point on a parametric curve, x' = -1, y' = 0, x'' = 0 and y'' = -1. What is d^2y/dx^2 = (x'y'' - y'x'')/(x')^3? (enter a number)",
			answer:   "-1", xp: -1, yp: 0, xpp: 0, ypp: -1,
			reason: "The numerator is (-1)(-1) - 0 = 1, and the denominator is -1, so the second derivative is -1.",
		},
		{
			question: "At a point on a parametric curve, x' = 1, y' = 12, x'' = 0 and y'' = 12. What is d^2y/dx^2 = (x'y'' - y'x'')/(x')^3? (enter a number)",
			answer:   "12", xp: 1, yp: 12, xpp: 0, ypp: 12,
			reason: "The numerator is 12 and the denominator is 1, so the second derivative is 12.",
		},
		{
			question: "At a point on a parametric curve, x' = 1, y' = 2, x'' = 1 and y'' = 4. What is d^2y/dx^2 = (x'y'' - y'x'')/(x')^3? (enter a number)",
			answer:   "2", xp: 1, yp: 2, xpp: 1, ypp: 4,
			reason: "The numerator is 1*4 - 2*1 = 2, and the denominator is 1, so the second derivative is 2.",
		},
		{
			question: "At a point on a parametric curve, x' = 4, y' = 1, x'' = 2 and y'' = 0. What is d^2y/dx^2 = (x'y'' - y'x'')/(x')^3? (enter a decimal)",
			answer:   "-0.03125", xp: 4, yp: 1, xpp: 2, ypp: 0,
			reason: "The numerator is 0 - 2 = -2, and the denominator is 64, so the second derivative is -0.03125.",
		},
		{
			question: "At a point on a parametric curve, x' = 1, y' = 1, x'' = 0 and y'' = 3. What is d^2y/dx^2 = (x'y'' - y'x'')/(x')^3? (enter a number)",
			answer:   "3", xp: 1, yp: 1, xpp: 0, ypp: 3,
			reason: "The numerator is 1*3 - 1*0 = 3, and the denominator is 1, so the second derivative is 3.",
		},
		{
			question: "At a point on a parametric curve, x' = 1, y' = 2, x'' = 0 and y'' = 5. What is d^2y/dx^2 = (x'y'' - y'x'')/(x')^3? (enter a number)",
			answer:   "5", xp: 1, yp: 2, xpp: 0, ypp: 5,
			reason: "The numerator is 1*5 - 2*0 = 5, and the denominator is 1, so the second derivative is 5.",
		},
		{
			question: "At a point on a parametric curve, x' = 1, y' = 1, x'' = 0 and y'' = 7. What is d^2y/dx^2 = (x'y'' - y'x'')/(x')^3? (enter a number)",
			answer:   "7", xp: 1, yp: 1, xpp: 0, ypp: 7,
			reason: "The numerator is 1*7 - 1*0 = 7, and the denominator is 1, so the second derivative is 7.",
		},
		{
			question: "At a point on a parametric curve, x' = 1, y' = 3, x'' = 0 and y'' = 1. What is d^2y/dx^2 = (x'y'' - y'x'')/(x')^3? (enter a number)",
			answer:   "1", xp: 1, yp: 3, xpp: 0, ypp: 1,
			reason: "The numerator is 1*1 - 3*0 = 1, and the denominator is 1, so the second derivative is 1.",
		},
		{
			question: "At a point on a parametric curve, x' = 1, y' = 3, x'' = 1 and y'' = 9. What is d^2y/dx^2 = (x'y'' - y'x'')/(x')^3? (enter a number)",
			answer:   "6", xp: 1, yp: 3, xpp: 1, ypp: 9,
			reason: "The numerator is 1*9 - 3*1 = 6, and the denominator is 1, so the second derivative is 6.",
		},
		{
			question: "At a point on a parametric curve, x' = 1, y' = 2, x'' = 2 and y'' = 8. What is d^2y/dx^2 = (x'y'' - y'x'')/(x')^3? (enter a number)",
			answer:   "4", xp: 1, yp: 2, xpp: 2, ypp: 8,
			reason: "The numerator is 1*8 - 2*2 = 4, and the denominator is 1, so the second derivative is 4.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	return generator.Problem{
		Question:    c.question,
		Answer:      c.answer,
		Explanation: c.reason,
		Facts: map[string]string{
			"xp": strconv.Itoa(c.xp), "yp": strconv.Itoa(c.yp),
			"xpp": strconv.Itoa(c.xpp), "ypp": strconv.Itoa(c.ypp), "answer": c.answer,
		},
	}
}

// ---------------------------------------------------------------------------
// 3. Finite differences (numeric)
// ---------------------------------------------------------------------------

type finiteDifferenceGen struct{}

func (g *finiteDifferenceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		method   string
		function string
		x, h     float64
		reason   string
	}
	claims := []claim{
		{
			question: "Use the forward difference (f(x+h) - f(x))/h for f(x) = x^2 at x = 3 with h = 0.1. What is the estimate? (enter a decimal)",
			answer:   "6.1", method: "forward", function: "square", x: 3, h: 0.1,
			reason: "f(3.1) = 9.61 and f(3) = 9, so (9.61 - 9)/0.1 = 6.1.",
		},
		{
			question: "Use the backward difference (f(x) - f(x-h))/h for f(x) = x^2 at x = 3 with h = 0.1. What is the estimate? (enter a decimal)",
			answer:   "5.9", method: "backward", function: "square", x: 3, h: 0.1,
			reason: "f(3) = 9 and f(2.9) = 8.41, so (9 - 8.41)/0.1 = 5.9.",
		},
		{
			question: "Use the central difference (f(x+h) - f(x-h))/(2h) for f(x) = x^2 at x = 3 with h = 0.1. What is the estimate? (enter a number)",
			answer:   "6", method: "central", function: "square", x: 3, h: 0.1,
			reason: "The central difference averages the two one-sided estimates, (9.61 - 8.41)/0.2 = 6.",
		},
		{
			question: "Use the forward difference (f(x+h) - f(x))/h for f(x) = x^3 at x = 2 with h = 0.1. What is the estimate? (enter a decimal)",
			answer:   "12.61", method: "forward", function: "cube", x: 2, h: 0.1,
			reason: "f(2.1) = 9.261 and f(2) = 8, so (9.261 - 8)/0.1 = 12.61.",
		},
		{
			question: "Use the central difference (f(x+h) - f(x-h))/(2h) for f(x) = x^3 at x = 2 with h = 0.1. What is the estimate? (enter a decimal)",
			answer:   "12.01", method: "central", function: "cube", x: 2, h: 0.1,
			reason: "f(2.1) = 9.261 and f(1.9) = 6.859, so (9.261 - 6.859)/0.2 = 12.01.",
		},
		{
			question: "Use the forward difference (f(x+h) - f(x))/h for f(x) = x^2 at x = 4 with h = 0.5. What is the estimate? (enter a decimal)",
			answer:   "8.5", method: "forward", function: "square", x: 4, h: 0.5,
			reason: "f(4.5) = 20.25 and f(4) = 16, so (20.25 - 16)/0.5 = 8.5.",
		},
		{
			question: "Use the backward difference (f(x) - f(x-h))/h for f(x) = x^2 at x = 4 with h = 0.5. What is the estimate? (enter a decimal)",
			answer:   "7.5", method: "backward", function: "square", x: 4, h: 0.5,
			reason: "f(4) = 16 and f(3.5) = 12.25, so (16 - 12.25)/0.5 = 7.5.",
		},
		{
			question: "Use the central difference (f(x+h) - f(x-h))/(2h) for f(x) = x^2 at x = 1 with h = 0.25. What is the estimate? (enter a number)",
			answer:   "2", method: "central", function: "square", x: 1, h: 0.25,
			reason: "f(1.25) = 1.5625 and f(0.75) = 0.5625, so (1.5625 - 0.5625)/0.5 = 2.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	return generator.Problem{
		Question:    c.question,
		Answer:      c.answer,
		Explanation: c.reason,
		Facts: map[string]string{
			"method": c.method, "function": c.function,
			"x": strconv.FormatFloat(c.x, 'g', -1, 64),
			"h": strconv.FormatFloat(c.h, 'g', -1, 64), "answer": c.answer,
		},
	}
}

// ---------------------------------------------------------------------------
// 4. Arc length in polar coordinates (numeric)
// ---------------------------------------------------------------------------

type polarArcLengthGen struct{}

func (g *polarArcLengthGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		kind     string
		a        int
		bLabel   string
		reason   string
	}
	claims := []claim{
		{
			question: "The polar curve r = sec(theta) is traced from theta = 0 to theta = pi/4. What is its arc length L = integral of sqrt(r^2 + (dr/dtheta)^2) dtheta? (enter a number)",
			answer:   "1", kind: "sec", a: 1, bLabel: "pi/4",
			reason: "For r = sec(theta), r^2 + (r')^2 = sec^2(theta), so L = tan(theta) from 0 to pi/4 = 1.",
		},
		{
			question: "The polar curve r = 2*sec(theta) is traced from theta = 0 to theta = pi/4. What is its arc length L = integral of sqrt(r^2 + (dr/dtheta)^2) dtheta? (enter a number)",
			answer:   "2", kind: "sec", a: 2, bLabel: "pi/4",
			reason: "The integrand is 2*sec^2(theta), so L = 2 tan(theta) from 0 to pi/4 = 2.",
		},
		{
			question: "The polar curve r = 3*sec(theta) is traced from theta = 0 to theta = pi/4. What is its arc length L = integral of sqrt(r^2 + (dr/dtheta)^2) dtheta? (enter a number)",
			answer:   "3", kind: "sec", a: 3, bLabel: "pi/4",
			reason: "The integrand is 3*sec^2(theta), so L = 3 tan(theta) from 0 to pi/4 = 3.",
		},
		{
			question: "The polar curve r = 4*sec(theta) is traced from theta = 0 to theta = pi/4. What is its arc length L = integral of sqrt(r^2 + (dr/dtheta)^2) dtheta? (enter a number)",
			answer:   "4", kind: "sec", a: 4, bLabel: "pi/4",
			reason: "The integrand is 4*sec^2(theta), so L = 4 tan(theta) from 0 to pi/4 = 4.",
		},
		{
			question: "The cardioid r = 1 + cos(theta) is traced from theta = 0 to theta = pi. What is its arc length L = integral of sqrt(r^2 + (dr/dtheta)^2) dtheta? (enter a number)",
			answer:   "4", kind: "cardioid", a: 1, bLabel: "pi",
			reason: "The integrand simplifies to 2*cos(theta/2), whose integral from 0 to pi is 4.",
		},
		{
			question: "The cardioid r = 1 + cos(theta) is traced from theta = 0 to theta = 2*pi. What is its arc length L = integral of sqrt(r^2 + (dr/dtheta)^2) dtheta? (enter a number)",
			answer:   "8", kind: "cardioid", a: 1, bLabel: "2*pi",
			reason: "The full cardioid has arc length 8: twice the half-cardioid length 4.",
		},
		{
			question: "The polar curve r = 5*sec(theta) is traced from theta = 0 to theta = pi/4. What is its arc length L = integral of sqrt(r^2 + (dr/dtheta)^2) dtheta? (enter a number)",
			answer:   "5", kind: "sec", a: 5, bLabel: "pi/4",
			reason: "The integrand is 5*sec^2(theta), so L = 5 tan(theta) from 0 to pi/4 = 5.",
		},
		{
			question: "The polar curve r = 6*sec(theta) is traced from theta = 0 to theta = pi/4. What is its arc length L = integral of sqrt(r^2 + (dr/dtheta)^2) dtheta? (enter a number)",
			answer:   "6", kind: "sec", a: 6, bLabel: "pi/4",
			reason: "The integrand is 6*sec^2(theta), so L = 6 tan(theta) from 0 to pi/4 = 6.",
		},
		{
			question: "The polar curve r = 7*sec(theta) is traced from theta = 0 to theta = pi/4. What is its arc length L = integral of sqrt(r^2 + (dr/dtheta)^2) dtheta? (enter a number)",
			answer:   "7", kind: "sec", a: 7, bLabel: "pi/4",
			reason: "The integrand is 7*sec^2(theta), so L = 7 tan(theta) from 0 to pi/4 = 7.",
		},
		{
			question: "The polar curve r = 8*sec(theta) is traced from theta = 0 to theta = pi/4. What is its arc length L = integral of sqrt(r^2 + (dr/dtheta)^2) dtheta? (enter a number)",
			answer:   "8", kind: "sec", a: 8, bLabel: "pi/4",
			reason: "The integrand is 8*sec^2(theta), so L = 8 tan(theta) from 0 to pi/4 = 8.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	return generator.Problem{
		Question:    c.question,
		Answer:      c.answer,
		Explanation: c.reason,
		Facts: map[string]string{
			"kind": c.kind, "a": strconv.Itoa(c.a), "b0": c.bLabel, "answer": c.answer,
		},
	}
}

// ---------------------------------------------------------------------------
// 5. Torsion of a space curve (numeric)
// ---------------------------------------------------------------------------

type torsionGen struct{}

func (g *torsionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		rp, rpp  [3]int
		rppp     [3]int
		reason   string
	}
	claims := []claim{
		{
			question: "A regular space curve has r' = (1, 0, 0), r'' = (0, 2, 0) and r''' = (0, 0, 6) at a point. What is its torsion tau = ((r' x r'') . r''')/|r' x r''|^2? (enter a number)",
			answer:   "3", rp: [3]int{1, 0, 0}, rpp: [3]int{0, 2, 0}, rppp: [3]int{0, 0, 6},
			reason: "r' x r'' = (0, 0, 2), its squared length is 4, and the dot product with r''' is 12, so tau = 12/4 = 3.",
		},
		{
			question: "A regular space curve has r' = (0, 1, 1), r'' = (-1, 0, 0) and r''' = (0, -1, 0) at a point. What is its torsion tau = ((r' x r'') . r''')/|r' x r''|^2? (enter a decimal)",
			answer:   "0.5", rp: [3]int{0, 1, 1}, rpp: [3]int{-1, 0, 0}, rppp: [3]int{0, -1, 0},
			reason: "r' x r'' = (0, -1, 1), its squared length is 2, and the dot product with r''' is 1, so tau = 0.5.",
		},
		{
			question: "A regular space curve has r' = (0, 1, 0), r'' = (-1, 0, 0) and r''' = (0, -1, 0) at a point. What is its torsion tau = ((r' x r'') . r''')/|r' x r''|^2? (enter a number)",
			answer:   "0", rp: [3]int{0, 1, 0}, rpp: [3]int{-1, 0, 0}, rppp: [3]int{0, -1, 0},
			reason: "r' x r'' = (0, 0, 1), and its dot product with r''' is 0, so the torsion is 0: the curve is planar.",
		},
		{
			question: "A regular space curve has r' = (1, 0, 0), r'' = (0, 4, 0) and r''' = (0, 0, 6) at a point. What is its torsion tau = ((r' x r'') . r''')/|r' x r''|^2? (enter a decimal)",
			answer:   "1.5", rp: [3]int{1, 0, 0}, rpp: [3]int{0, 4, 0}, rppp: [3]int{0, 0, 6},
			reason: "r' x r'' = (0, 0, 4), its squared length is 16, and the dot product with r''' is 24, so tau = 24/16 = 1.5.",
		},
		{
			question: "A regular space curve has r' = (1, 0, 0), r'' = (0, 2, 0) and r''' = (0, 0, 12) at a point. What is its torsion tau = ((r' x r'') . r''')/|r' x r''|^2? (enter a number)",
			answer:   "6", rp: [3]int{1, 0, 0}, rpp: [3]int{0, 2, 0}, rppp: [3]int{0, 0, 12},
			reason: "r' x r'' = (0, 0, 2), its squared length is 4, and the dot product with r''' is 24, so tau = 6.",
		},
		{
			question: "A regular space curve has r' = (0, 2, 1), r'' = (-2, 0, 0) and r''' = (0, -2, 0) at a point. What is its torsion tau = ((r' x r'') . r''')/|r' x r''|^2? (enter a decimal)",
			answer:   "0.2", rp: [3]int{0, 2, 1}, rpp: [3]int{-2, 0, 0}, rppp: [3]int{0, -2, 0},
			reason: "r' x r'' = (0, -2, 4), its squared length is 20, and the dot product with r''' is 4, so tau = 0.2.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	return generator.Problem{
		Question:    c.question,
		Answer:      c.answer,
		Explanation: c.reason,
		Facts: map[string]string{
			"rp":     strconv.Itoa(c.rp[0]) + "," + strconv.Itoa(c.rp[1]) + "," + strconv.Itoa(c.rp[2]),
			"rpp":    strconv.Itoa(c.rpp[0]) + "," + strconv.Itoa(c.rpp[1]) + "," + strconv.Itoa(c.rpp[2]),
			"rppp":   strconv.Itoa(c.rppp[0]) + "," + strconv.Itoa(c.rppp[1]) + "," + strconv.Itoa(c.rppp[2]),
			"answer": c.answer,
		},
	}
}

// ---------------------------------------------------------------------------
// 6. The Frenet frame (tuple)
// ---------------------------------------------------------------------------

type frenetGen struct{}

func frenetCrossInt(a, b [3]int) [3]int {
	return [3]int{
		a[1]*b[2] - a[2]*b[1],
		a[2]*b[0] - a[0]*b[2],
		a[0]*b[1] - a[1]*b[0],
	}
}

func frenetNormInt(a [3]int) int {
	return int(math.Sqrt(float64(a[0]*a[0] + a[1]*a[1] + a[2]*a[2])))
}

func frenetString(v [3]int) string {
	return fmt.Sprintf("(%d, %d, %d)", v[0], v[1], v[2])
}

// frenetGen asks for the whole frame, not just the binormal: the concept is
// named for the frame, and a question that returns one of its three vectors
// tests a third of the concept.
func (g *frenetGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		rp, rpp [3]int
	}
	claims := []claim{
		{rp: [3]int{0, 1, 0}, rpp: [3]int{-1, 0, 0}},
		{rp: [3]int{1, 0, 0}, rpp: [3]int{0, 2, 0}},
		{rp: [3]int{1, 0, 0}, rpp: [3]int{0, -2, 0}},
		{rp: [3]int{1, 0, 0}, rpp: [3]int{0, 0, 2}},
		{rp: [3]int{0, 0, 1}, rpp: [3]int{-1, 0, 0}},
		{rp: [3]int{1, 0, 0}, rpp: [3]int{0, 0, -2}},
		{rp: [3]int{0, 1, 0}, rpp: [3]int{0, 0, 2}},
		{rp: [3]int{0, 1, 0}, rpp: [3]int{0, 0, -2}},
		{rp: [3]int{0, 0, 1}, rpp: [3]int{1, 0, 0}},
	}
	c := claims[rand.Intn(len(claims))]
	// Every r' in the table is a unit axis vector, so T is r' itself.
	tangent := c.rp
	cross := frenetCrossInt(c.rp, c.rpp)
	norm := frenetNormInt(cross)
	binormal := [3]int{cross[0] / norm, cross[1] / norm, cross[2] / norm}
	normal := frenetCrossInt(binormal, tangent)
	answer := fmt.Sprintf("(%d,%d,%d,%d,%d,%d,%d,%d,%d)",
		tangent[0], tangent[1], tangent[2],
		normal[0], normal[1], normal[2],
		binormal[0], binormal[1], binormal[2])
	return generator.Problem{
		Question: fmt.Sprintf("For a regular space curve with r' = %s and r'' = %s at a point, the Frenet frame is (T, N, B) with T = r'/|r'|, B = (r' x r'')/|r' x r''| and N = B x T. Enter the nine components in the order (Tx, Ty, Tz, Nx, Ny, Nz, Bx, By, Bz).",
			frenetString(c.rp), frenetString(c.rpp)),
		Answer: answer,
		Explanation: fmt.Sprintf("r' is already a unit vector, so T = %s. r' x r'' = %s has length %d, so B = %s. Then N = B x T = %s, which points the same way as the component of r'' perpendicular to T.",
			frenetString(tangent), frenetString(cross), norm, frenetString(binormal), frenetString(normal)),
		Facts: map[string]string{
			"rp":     strconv.Itoa(c.rp[0]) + "," + strconv.Itoa(c.rp[1]) + "," + strconv.Itoa(c.rp[2]),
			"rpp":    strconv.Itoa(c.rpp[0]) + "," + strconv.Itoa(c.rpp[1]) + "," + strconv.Itoa(c.rpp[2]),
			"answer": answer,
		},
	}
}

// ---------------------------------------------------------------------------
// 7. Second-order Taylor approximation (numeric)
// ---------------------------------------------------------------------------

type taylor2Gen struct{}

func (g *taylor2Gen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		f0, fx   float64
		fy       float64
		fxx, fxy float64
		fyy      float64
		dx, dy   float64
		reason   string
	}
	claims := []claim{
		{
			question: "At the expansion point, f = 5, fx = 2, fy = 4, fxx = 2, fxy = 0 and fyy = 2. For dx = 0.1 and dy = 0.1, what is the second-order Taylor approximation f + fx dx + fy dy + (1/2)(fxx dx^2 + 2 fxy dx dy + fyy dy^2)? (enter a decimal)",
			answer:   "5.62", f0: 5, fx: 2, fy: 4, fxx: 2, fxy: 0, fyy: 2, dx: 0.1, dy: 0.1,
			reason: "The linear part is 0.2 + 0.4 = 0.6 and the Hessian part is 0.5(0.02 + 0 + 0.02) = 0.02, so the value is 5.62.",
		},
		{
			question: "At the expansion point, f = 0, fx = 0, fy = 0, fxx = 2, fxy = 0 and fyy = 2. For dx = 1 and dy = 1, what is the second-order Taylor approximation f + fx dx + fy dy + (1/2)(fxx dx^2 + 2 fxy dx dy + fyy dy^2)? (enter a number)",
			answer:   "2", f0: 0, fx: 0, fy: 0, fxx: 2, fxy: 0, fyy: 2, dx: 1, dy: 1,
			reason: "The linear part vanishes and the Hessian part is 0.5(2 + 0 + 2) = 2.",
		},
		{
			question: "At the expansion point, f = 2, fx = 3, fy = 1, fxx = 2, fxy = 1 and fyy = 0. For dx = 0.1 and dy = 0.2, what is the second-order Taylor approximation f + fx dx + fy dy + (1/2)(fxx dx^2 + 2 fxy dx dy + fyy dy^2)? (enter a decimal)",
			answer:   "2.53", f0: 2, fx: 3, fy: 1, fxx: 2, fxy: 1, fyy: 0, dx: 0.1, dy: 0.2,
			reason: "The linear part is 0.3 + 0.2 = 0.5 and the Hessian part is 0.5(0.02 + 0.04 + 0) = 0.03, so the value is 2.53.",
		},
		{
			question: "At the expansion point, f = 3, fx = 4, fy = -2, fxx = 2, fxy = 0 and fyy = -2. For dx = 0.1 and dy = 0.1, what is the second-order Taylor approximation f + fx dx + fy dy + (1/2)(fxx dx^2 + 2 fxy dx dy + fyy dy^2)? (enter a decimal)",
			answer:   "3.2", f0: 3, fx: 4, fy: -2, fxx: 2, fxy: 0, fyy: -2, dx: 0.1, dy: 0.1,
			reason: "The linear part is 0.4 - 0.2 = 0.2 and the Hessian part is 0.5(0.02 - 0.02) = 0, so the value is 3.2.",
		},
		{
			question: "At the expansion point, f = 2, fx = 2, fy = 1, fxx = 0, fxy = 1 and fyy = 0. For dx = 0.1 and dy = 0.1, what is the second-order Taylor approximation f + fx dx + fy dy + (1/2)(fxx dx^2 + 2 fxy dx dy + fyy dy^2)? (enter a decimal)",
			answer:   "2.31", f0: 2, fx: 2, fy: 1, fxx: 0, fxy: 1, fyy: 0, dx: 0.1, dy: 0.1,
			reason: "The linear part is 0.2 + 0.1 = 0.3 and the Hessian part is 0.5(0 + 0.02 + 0) = 0.01, so the value is 2.31.",
		},
		{
			question: "At the expansion point, f = 5, fx = 4, fy = 2, fxx = 2, fxy = 0 and fyy = 2. For dx = -0.1 and dy = 0.2, what is the second-order Taylor approximation f + fx dx + fy dy + (1/2)(fxx dx^2 + 2 fxy dx dy + fyy dy^2)? (enter a decimal)",
			answer:   "5.05", f0: 5, fx: 4, fy: 2, fxx: 2, fxy: 0, fyy: 2, dx: -0.1, dy: 0.2,
			reason: "The linear part is -0.4 + 0.4 = 0 and the Hessian part is 0.5(0.02 + 0 + 0.08) = 0.05, so the value is 5.05.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	format := func(v float64) string { return strconv.FormatFloat(v, 'g', -1, 64) }
	return generator.Problem{
		Question:    c.question,
		Answer:      c.answer,
		Explanation: c.reason,
		Facts: map[string]string{
			"f0": format(c.f0), "fx": format(c.fx), "fy": format(c.fy),
			"fxx": format(c.fxx), "fxy": format(c.fxy), "fyy": format(c.fyy),
			"dx": format(c.dx), "dy": format(c.dy), "answer": c.answer,
		},
	}
}

// ---------------------------------------------------------------------------
// 8. Hydrostatic force (numeric)
// ---------------------------------------------------------------------------

type hydrostaticGen struct{}

func (g *hydrostaticGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		w, a, b  int
		reason   string
	}
	claims := []claim{
		{
			question: "A vertical rectangular plate is 2 metres wide, with its top edge at depth 0 and its bottom edge at depth 3. Taking rho*g = 10, what is the hydrostatic force F = rho*g * integral from a to b of w*y dy? (enter a number)",
			answer:   "90", w: 2, a: 0, b: 3,
			reason: "The integral is 2 * (3^2 - 0^2)/2 = 9, so F = 10 * 9 = 90.",
		},
		{
			question: "A vertical rectangular plate is 3 metres wide, with its top edge at depth 1 and its bottom edge at depth 3. Taking rho*g = 10, what is the hydrostatic force F = rho*g * integral from a to b of w*y dy? (enter a number)",
			answer:   "120", w: 3, a: 1, b: 3,
			reason: "The integral is 3 * (9 - 1)/2 = 12, so F = 10 * 12 = 120.",
		},
		{
			question: "A vertical rectangular plate is 4 metres wide, with its top edge at depth 2 and its bottom edge at depth 4. Taking rho*g = 10, what is the hydrostatic force F = rho*g * integral from a to b of w*y dy? (enter a number)",
			answer:   "240", w: 4, a: 2, b: 4,
			reason: "The integral is 4 * (16 - 4)/2 = 24, so F = 10 * 24 = 240.",
		},
		{
			question: "A vertical rectangular plate is 1 metre wide, with its top edge at depth 0 and its bottom edge at depth 4. Taking rho*g = 10, what is the hydrostatic force F = rho*g * integral from a to b of w*y dy? (enter a number)",
			answer:   "80", w: 1, a: 0, b: 4,
			reason: "The integral is 1 * (16 - 0)/2 = 8, so F = 10 * 8 = 80.",
		},
		{
			question: "A vertical rectangular plate is 5 metres wide, with its top edge at depth 0 and its bottom edge at depth 2. Taking rho*g = 10, what is the hydrostatic force F = rho*g * integral from a to b of w*y dy? (enter a number)",
			answer:   "100", w: 5, a: 0, b: 2,
			reason: "The integral is 5 * (4 - 0)/2 = 10, so F = 10 * 10 = 100.",
		},
		{
			question: "A vertical rectangular plate is 2 metres wide, with its top edge at depth 3 and its bottom edge at depth 5. Taking rho*g = 10, what is the hydrostatic force F = rho*g * integral from a to b of w*y dy? (enter a number)",
			answer:   "160", w: 2, a: 3, b: 5,
			reason: "The integral is 2 * (25 - 9)/2 = 16, so F = 10 * 16 = 160.",
		},
	}
	// Draw the plate instead of reciting four of them. The depth gap is even so
	// the integral w*(b^2-a^2)/2 stays an integer.
	w := 1 + rand.Intn(5)
	a := rand.Intn(3)
	b := a + 2 + 2*rand.Intn(2)
	integral := w * (b*b - a*a) / 2
	claims = append(claims,
		claim{
			question: fmt.Sprintf("A vertical rectangular plate is %d metres wide, with its top edge at depth %d and its bottom edge at depth %d. Taking rho*g = 10, what is the hydrostatic force? (enter a number)", w, a, b),
			answer:   strconv.Itoa(10 * integral),
			w:        w, a: a, b: b,
			reason: fmt.Sprintf("The integral is %d * (%d^2 - %d^2)/2 = %d, so F = 10 * %d = %d.", w, b, a, integral, integral, 10*integral),
		},
	)
	c := claims[rand.Intn(len(claims))]
	return generator.Problem{
		Question:    c.question,
		Answer:      c.answer,
		Explanation: c.reason,
		Facts: map[string]string{
			"w": strconv.Itoa(c.w), "a": strconv.Itoa(c.a), "b": strconv.Itoa(c.b), "answer": c.answer,
		},
	}
}

// ---------------------------------------------------------------------------
// 9. Pappus's centroid theorem (symbolic)
// ---------------------------------------------------------------------------

type pappusGen struct{}

func (g *pappusGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		area, d  int
		reason   string
	}
	claims := []claim{
		{
			question: "A plane region of area 6 has its centroid 5 units from an external axis that does not cut the region. Pappus's centroid theorem gives the volume swept as V = 2*pi*d*A. What is V? (enter an exact multiple of pi, such as 60*pi)",
			answer:   "60*pi", area: 6, d: 5,
			reason: "V = 2*pi*5*6 = 60*pi.",
		},
		{
			question: "A plane region of area 4 has its centroid 3 units from an external axis that does not cut the region. Pappus's centroid theorem gives the volume swept as V = 2*pi*d*A. What is V? (enter an exact multiple of pi, such as 24*pi)",
			answer:   "24*pi", area: 4, d: 3,
			reason: "V = 2*pi*3*4 = 24*pi.",
		},
		{
			question: "A plane region of area 6 has its centroid 2 units from an external axis that does not cut the region. Pappus's centroid theorem gives the volume swept as V = 2*pi*d*A. What is V? (enter an exact multiple of pi, such as 24*pi)",
			answer:   "24*pi", area: 6, d: 2,
			reason: "V = 2*pi*2*6 = 24*pi.",
		},
		{
			question: "A plane region of area 9 has its centroid 4 units from an external axis that does not cut the region. Pappus's centroid theorem gives the volume swept as V = 2*pi*d*A. What is V? (enter an exact multiple of pi, such as 72*pi)",
			answer:   "72*pi", area: 9, d: 4,
			reason: "V = 2*pi*4*9 = 72*pi.",
		},
		{
			question: "A plane region of area 10 has its centroid 1 unit from an external axis that does not cut the region. Pappus's centroid theorem gives the volume swept as V = 2*pi*d*A. What is V? (enter an exact multiple of pi, such as 20*pi)",
			answer:   "20*pi", area: 10, d: 1,
			reason: "V = 2*pi*1*10 = 20*pi.",
		},
		{
			question: "A plane region of area 8 has its centroid 3 units from an external axis that does not cut the region. Pappus's centroid theorem gives the volume swept as V = 2*pi*d*A. What is V? (enter an exact multiple of pi, such as 48*pi)",
			answer:   "48*pi", area: 8, d: 3,
			reason: "V = 2*pi*3*8 = 48*pi.",
		},
		{
			question: "A plane region of area 3 has its centroid 1 unit from an external axis that does not cut the region. Pappus's centroid theorem gives the volume swept as V = 2*pi*d*A. What is V? (enter an exact multiple of pi, such as 6*pi)",
			answer:   "6*pi", area: 3, d: 1,
			reason: "V = 2*pi*1*3 = 6*pi.",
		},
		{
			question: "A plane region of area 4 has its centroid 2 units from an external axis that does not cut the region. Pappus's centroid theorem gives the volume swept as V = 2*pi*d*A. What is V? (enter an exact multiple of pi, such as 16*pi)",
			answer:   "16*pi", area: 4, d: 2,
			reason: "V = 2*pi*2*4 = 16*pi.",
		},
		{
			question: "A plane region of area 5 has its centroid 3 units from an external axis that does not cut the region. Pappus's centroid theorem gives the volume swept as V = 2*pi*d*A. What is V? (enter an exact multiple of pi, such as 30*pi)",
			answer:   "30*pi", area: 5, d: 3,
			reason: "V = 2*pi*3*5 = 30*pi.",
		},
		{
			question: "A plane region of area 1 has its centroid 6 units from an external axis that does not cut the region. Pappus's centroid theorem gives the volume swept as V = 2*pi*d*A. What is V? (enter an exact multiple of pi, such as 12*pi)",
			answer:   "12*pi", area: 1, d: 6,
			reason: "V = 2*pi*6*1 = 12*pi.",
		},
		{
			question: "A plane region of area 7 has its centroid 3 units from an external axis that does not cut the region. Pappus's centroid theorem gives the volume swept as V = 2*pi*d*A. What is V? (enter an exact multiple of pi, such as 42*pi)",
			answer:   "42*pi", area: 7, d: 3,
			reason: "V = 2*pi*3*7 = 42*pi.",
		},
		{
			question: "A plane region of area 9 has its centroid 2 units from an external axis that does not cut the region. Pappus's centroid theorem gives the volume swept as V = 2*pi*d*A. What is V? (enter an exact multiple of pi, such as 36*pi)",
			answer:   "36*pi", area: 9, d: 2,
			reason: "V = 2*pi*2*9 = 36*pi.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	return generator.Problem{
		Question:    c.question,
		Answer:      c.answer,
		Explanation: c.reason,
		Facts: map[string]string{
			"area": strconv.Itoa(c.area), "d": strconv.Itoa(c.d), "answer": c.answer,
		},
	}
}

// ---------------------------------------------------------------------------
// 10. Moment of inertia (numeric)
// ---------------------------------------------------------------------------

type momentInertiaGen struct{}

func (g *momentInertiaGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question  string
		answer    string
		kind      string
		m, l      int
		r         int
		masses    []int
		distances []int
		reason    string
	}
	claims := []claim{
		{
			question: "A uniform rod of mass 12 and length 2 rotates about an axis through its centre, perpendicular to the rod. What is its moment of inertia M*L^2/12? (enter a number)",
			answer:   "4", kind: "rod_center", m: 12, l: 2,
			reason: "M*L^2/12 = 12*4/12 = 4.",
		},
		{
			question: "A uniform rod of mass 12 and length 1 rotates about an axis through one end, perpendicular to the rod. What is its moment of inertia M*L^2/3? (enter a number)",
			answer:   "4", kind: "rod_end", m: 12, l: 1,
			reason: "M*L^2/3 = 12*1/3 = 4.",
		},
		{
			question: "A uniform rod of mass 3 and length 2 rotates about an axis through its centre, perpendicular to the rod. What is its moment of inertia M*L^2/12? (enter a number)",
			answer:   "1", kind: "rod_center", m: 3, l: 2,
			reason: "M*L^2/12 = 3*4/12 = 1.",
		},
		{
			question: "A uniform rod of mass 6 and length 2 rotates about an axis through one end, perpendicular to the rod. What is its moment of inertia M*L^2/3? (enter a number)",
			answer:   "8", kind: "rod_end", m: 6, l: 2,
			reason: "M*L^2/3 = 6*4/3 = 8.",
		},
		{
			question: "Point masses 2 and 3 sit at distances 1 and 2 from an axis. What is the moment of inertia sum m*r^2? (enter a number)",
			answer:   "14", kind: "points", masses: []int{2, 3}, distances: []int{1, 2},
			reason: "2*1^2 + 3*2^2 = 2 + 12 = 14.",
		},
		{
			question: "A uniform disk of mass 4 and radius 2 rotates about its central axis. What is its moment of inertia (1/2)M*R^2? (enter a number)",
			answer:   "8", kind: "disk", m: 4, r: 2,
			reason: "(1/2)M*R^2 = 0.5*4*4 = 8.",
		},
		{
			question: "A uniform rod of mass 6 and length 2 rotates about an axis through its centre, perpendicular to the rod. What is its moment of inertia M*L^2/12? (enter a number)",
			answer:   "2", kind: "rod_center", m: 6, l: 2,
			reason: "M*L^2/12 = 6*4/12 = 2.",
		},
		{
			question: "A uniform rod of mass 9 and length 1 rotates about an axis through one end, perpendicular to the rod. What is its moment of inertia M*L^2/3? (enter a number)",
			answer:   "3", kind: "rod_end", m: 9, l: 1,
			reason: "M*L^2/3 = 9*1/3 = 3.",
		},
		{
			question: "A uniform rod of mass 12 and length 4 rotates about an axis through its centre, perpendicular to the rod. What is its moment of inertia M*L^2/12? (enter a number)",
			answer:   "16", kind: "rod_center", m: 12, l: 4,
			reason: "M*L^2/12 = 12*16/12 = 16.",
		},
		{
			question: "A uniform rod of mass 15 and length 2 rotates about an axis through one end, perpendicular to the rod. What is its moment of inertia M*L^2/3? (enter a number)",
			answer:   "20", kind: "rod_end", m: 15, l: 2,
			reason: "M*L^2/3 = 15*4/3 = 20.",
		},
		{
			question: "A uniform disk of mass 4 and radius 3 rotates about its central axis. What is its moment of inertia (1/2)M*R^2? (enter a number)",
			answer:   "18", kind: "disk", m: 4, r: 3,
			reason: "(1/2)M*R^2 = 0.5*4*9 = 18.",
		},
		{
			question: "A uniform disk of mass 6 and radius 2 rotates about its central axis. What is its moment of inertia (1/2)M*R^2? (enter a number)",
			answer:   "12", kind: "disk", m: 6, r: 2,
			reason: "(1/2)M*R^2 = 0.5*6*4 = 12.",
		},
		{
			question: "Point masses 1 and 2 sit at distances 2 and 3 from an axis. What is the moment of inertia sum m*r^2? (enter a number)",
			answer:   "22", kind: "points", masses: []int{1, 2}, distances: []int{2, 3},
			reason: "1*2^2 + 2*3^2 = 4 + 18 = 22.",
		},
		{
			question: "Point masses 3 and 4 sit at distances 1 and 2 from an axis. What is the moment of inertia sum m*r^2? (enter a number)",
			answer:   "19", kind: "points", masses: []int{3, 4}, distances: []int{1, 2},
			reason: "3*1^2 + 4*2^2 = 3 + 16 = 19.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	facts := map[string]string{"kind": c.kind, "answer": c.answer}
	switch c.kind {
	case "rod_center", "rod_end":
		facts["m"] = strconv.Itoa(c.m)
		facts["l"] = strconv.Itoa(c.l)
	case "disk":
		facts["m"] = strconv.Itoa(c.m)
		facts["r"] = strconv.Itoa(c.r)
	case "points":
		masses := make([]string, len(c.masses))
		distances := make([]string, len(c.distances))
		for i := range c.masses {
			masses[i] = strconv.Itoa(c.masses[i])
			distances[i] = strconv.Itoa(c.distances[i])
		}
		facts["masses"] = strings.Join(masses, ",")
		facts["distances"] = strings.Join(distances, ",")
	}
	return generator.Problem{
		Question:    c.question,
		Answer:      c.answer,
		Explanation: c.reason,
		Facts:       facts,
	}
}
