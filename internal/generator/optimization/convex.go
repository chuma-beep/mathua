package optimization

import (
	"fmt"

	"github.com/chuma-beep/mathua/internal/generator"
)

// Convex optimization and duality. Every family computes a bound, a primal or
// dual value, a gap or a tableau entry from drawn coefficients.

// ---------------------------------------------------------------------------
// Linear programming standard form
// ---------------------------------------------------------------------------

type lpStandardFormGen struct{}

func (g *lpStandardFormGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		a, b, x, y := limit(1, 6), limit(1, 6), limit(1, 6), limit(1, 6)
		c := a*x + b*y + limit(1, 9)
		return at(kp, fmt.Sprintf("The constraint %d x + %d y ≤ %d is rewritten in standard form with a slack variable s. At (x, y) = (%d, %d), enter s.", a, b, c, x, y), c-a*x-b*y,
			fmt.Sprintf("The slack is c - ax - by = %d - %d - %d = %d.", c, a*x, b*y, c-a*x-b*y))
	case 2:
		a, b, c := limit(1, 5), limit(1, 5), limit(1, 5)
		return at(kp, fmt.Sprintf("A standard-form row reads %d x + %d y + s = %d. Enter the coefficient of x.", a, b, c), a,
			fmt.Sprintf("In standard form each variable has one explicit coefficient, so the coefficient of x is %d.", a))
	default:
		n, m := limit(6, 20), limit(2, 4)
		return at(kp, fmt.Sprintf("A standard-form linear program has %d variables and %d equality constraints, with a full-rank constraint matrix. Enter the number of nonbasic variables.", n, m), n-m,
			fmt.Sprintf("A basic solution sets n - m variables to zero, so %d - %d = %d variables are nonbasic.", n, m, n-m))
	}
}

// ---------------------------------------------------------------------------
// Basic feasible solutions
// ---------------------------------------------------------------------------

type basicSolutionGen struct{}

func (g *basicSolutionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		a, c, d := limit(1, 4), limit(1, 3), limit(1, 4)
		x, y := limit(1, 9), limit(1, 9)
		b1, b2 := a*x, c*x+d*y
		return at(kp, fmt.Sprintf("A basis produces the equations %d x = %d and %d x + %d y = %d. Enter the basic variable x.", a, b1, c, d, b2), x,
			fmt.Sprintf("The first equation gives x = %d / %d = %d, and the second then gives y = %d.", b1, a, x, y))
	case 2:
		a, b, x, y := limit(2, 6), limit(1, 5), limit(1, 9), limit(1, 9)
		slack := limit(0, 9)
		return at(kp, fmt.Sprintf("A basic solution has x = %d and y = %d, and the constraint is %d x + %d y ≤ %d. Enter its slack.", x, y, a, b, a*x+b*y+slack), slack,
			fmt.Sprintf("ax + by = %d·%d + %d·%d = %d, so the slack is %d - %d = %d.", a, x, b, y, a*x+b*y, a*x+b*y+slack, a*x+b*y, slack))
	default:
		n, m := limit(5, 16), limit(2, 3)
		return at(kp, fmt.Sprintf("A standard-form LP has %d variables and %d constraints. Enter the number of basic variables in a basic solution.", n, m), m,
			fmt.Sprintf("A basic solution fixes one basic variable per constraint, so there are %d.", m))
	}
}

// ---------------------------------------------------------------------------
// The simplex pivot
// ---------------------------------------------------------------------------

type simplexPivotGen struct{}

func (g *simplexPivotGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		worst := -limit(1, 9)
		return at(kp, fmt.Sprintf("The reduced costs of three entering candidates are %d, %d and %d. Enter the most negative reduced cost, which the pivoting rule selects.", worst, worst+limit(0, 4), worst+limit(0, 6)), worst,
			fmt.Sprintf("The most negative coefficient is %d, so that column enters the basis.", worst))
	case 2:
		a, b, val := limit(1, 4), limit(1, 5), limit(1, 9)
		return at(kp, fmt.Sprintf("In a ratio test the row has entries %d and %d, and the right-hand side is %d. Enter the ratio for the first entry, rounded down.", a, b, val), val/a,
			fmt.Sprintf("The ratio is %d / %d, which rounds down to %d.", val, a, val/a))
	default:
		pivot, factor, entry := limit(2, 6), limit(1, 4), limit(1, 9)
		return at(kp, fmt.Sprintf("A pivot divides its row by %d. Another row is updated as new = old - %d·pivot_row. Its current entry in the pivot column is %d, so new = old - %d·1. Enter the new entry when old = %d.", pivot, factor, factor, factor, entry), entry-factor,
			fmt.Sprintf("The update is %d - %d·1 = %d.", entry, factor, entry-factor))
	}
}

// ---------------------------------------------------------------------------
// The dual linear program
// ---------------------------------------------------------------------------

type lpDualGen struct{}

func (g *lpDualGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		b := limit(2, 20)
		return at(kp, fmt.Sprintf("A primal maximization has the right-hand side b = (%d, %d). Its dual minimizes bᵀy. Enter the coefficient of y1 in the dual objective.", b, limit(2, 20)), b,
			fmt.Sprintf("The dual objective is %d y1 + ... so the coefficient of y1 is %d.", b, b))
	case 2:
		m := limit(3, 30)
		return at(kp, fmt.Sprintf("A primal problem has %d inequality constraints. Enter the number of dual variables.", m), m,
			fmt.Sprintf("One dual variable belongs to each primal constraint, so there are %d.", m))
	default:
		n, m := limit(2, 6), limit(7, 20)
		return at(kp, fmt.Sprintf("A primal problem has %d variables and %d constraints. Enter the number of dual variables minus the number of primal variables.", n, m), m-n,
			fmt.Sprintf("The dual has %d variables, so %d - %d = %d.", m, m, n, m-n))
	}
}

// ---------------------------------------------------------------------------
// The Lagrange dual function
// ---------------------------------------------------------------------------

type dualFunctionGen struct{}

func (g *dualFunctionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		c, y, a, x, b := limit(2, 9), limit(1, 5), limit(1, 3), limit(1, 8), limit(1, 9)
		return at(kp, fmt.Sprintf("For L(x, y) = %d x - %d(%d x - %d) with x = %d, enter L.", c, y, a, b, x), c*x-y*(a*x-b),
			fmt.Sprintf("L = %d·%d - %d(%d·%d - %d) = %d.", c, x, y, a, x, b, c*x-y*(a*x-b)))
	case 2:
		c, b := limit(2, 9), limit(1, 20)
		y := c
		return at(kp, fmt.Sprintf("The Lagrangian is %d x - y(x - %d) over x ≥ 0. At the dual-optimal y = %d the x coefficient vanishes and the dual value is y b. Enter that value.", c, b, y), y*b,
			fmt.Sprintf("With y = %d the value is %d·%d = %d.", y, y, b, y*b))
	default:
		p, d := limit(10, 60), limit(1, 9)
		return at(kp, fmt.Sprintf("A primal feasible value is %d and a dual feasible value is %d. Enter the duality gap p - d.", p, p-d), d,
			fmt.Sprintf("The gap is %d - %d = %d, which weak duality keeps nonnegative.", p, p-d, d))
	}
}

// ---------------------------------------------------------------------------
// Weak duality
// ---------------------------------------------------------------------------

type weakDualityGen struct{}

func (g *weakDualityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		p, d := limit(20, 90), limit(1, 15)
		return at(kp, fmt.Sprintf("A primal feasible objective is %d and a dual feasible objective is %d. Enter the gap.", p, p-d), d,
			fmt.Sprintf("%d - %d = %d, and weak duality says the gap cannot be negative.", p, p-d, d))
	case 2:
		c, a, y := limit(3, 9), limit(2, 6), limit(1, 4)
		return at(kp, fmt.Sprintf("A dual-feasible y has Aᵀy ≤ c. For a column with a = %d, y = %d and c = %d, enter the reduced cost c - a y.", a, y, c), c-a*y,
			fmt.Sprintf("The reduced cost is %d - %d·%d = %d, which dual feasibility keeps nonnegative.", c, a, y, c-a*y))
	default:
		p, d := limit(30, 80), limit(2, 20)
		return at(kp, fmt.Sprintf("For a minimization, a dual feasible lower bound is %d and a primal feasible value is %d. Enter the gap.", d, d+p), p,
			fmt.Sprintf("The gap is %d - %d = %d.", d+p, d, p))
	}
}

// ---------------------------------------------------------------------------
// Strong duality
// ---------------------------------------------------------------------------

type strongDualityGen struct{}

func (g *strongDualityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		opt := limit(10, 80)
		return at(kp, fmt.Sprintf("A convex problem satisfies Slater's condition. Its dual optimum is %d and strong duality holds. Enter the primal optimum.", opt), opt,
			fmt.Sprintf("Strong duality makes the two optima equal, so the primal optimum is %d.", opt))
	case 2:
		y, b := limit(1, 6), limit(2, 20)
		return at(kp, fmt.Sprintf("At a primal-dual optimum the multiplier y = %d is positive and the constraint is a x - %d ≤ 0. Complementary slackness forces it tight. Enter a x.", y, b), b,
			fmt.Sprintf("A positive multiplier forces zero slack, so a x = %d.", b))
	default:
		p, d := limit(20, 60), limit(1, 8)
		return at(kp, fmt.Sprintf("A primal value is %d and a dual value is %d. Enter the gap as a percentage of the primal value, rounded down.", p, p-d), d*100/p,
			fmt.Sprintf("The gap is %d, which is %d percent of %d.", d, d*100/p, p))
	}
}

// ---------------------------------------------------------------------------
// Complementary slackness
// ---------------------------------------------------------------------------

type complementarySlacknessGen struct{}

func (g *complementarySlacknessGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		x, slack := limit(1, 6), limit(1, 6)
		mult := limit(1, 4)
		return at(kp, fmt.Sprintf("At a candidate optimum x = %d, a constraint has slack %d and multiplier %d. Enter the complementary slackness product x·multiplier + slack·multiplier.", x, slack, mult), (x+slack)*mult,
			fmt.Sprintf("Both terms share the multiplier: (%d + %d)·%d = %d.", x, slack, mult, (x+slack)*mult))
	case 2:
		slack := limit(1, 30)
		return at(kp, fmt.Sprintf("A constraint has positive slack %d and complementary slackness must hold. Enter its multiplier.", slack), 0,
			fmt.Sprintf("Positive slack forces a zero multiplier, because the product slack·multiplier must be zero."))
	default:
		a, y, b, x := limit(1, 4), limit(1, 5), limit(1, 9), limit(1, 6)
		return at(kp, fmt.Sprintf("A multiplier is y = %d, the constraint is %d x ≤ %d and x = %d. Enter y·(%d - %d x).", y, a, b, x, b, a), y*(b-a*x),
			fmt.Sprintf("The slack is %d - %d·%d = %d, so the complementary term is %d·%d = %d.", b, a, x, b-a*x, y, b-a*x, y*(b-a*x)))
	}
}

// ---------------------------------------------------------------------------
// Shadow prices
// ---------------------------------------------------------------------------

type shadowPriceGen struct{}

func (g *shadowPriceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		opt, price, delta := limit(20, 70), limit(1, 6), limit(1, 5)
		return at(kp, fmt.Sprintf("A constraint has shadow price %d and its right-hand side rises by %d. Enter the new optimal objective when the old optimum was %d.", price, delta, opt), opt+price*delta,
			fmt.Sprintf("The objective rises by %d·%d = %d, so the new value is %d.", price, delta, price*delta, opt+price*delta))
	case 2:
		opt, price, delta := limit(15, 60), limit(1, 5), limit(1, 4)
		return at(kp, fmt.Sprintf("A right-hand side falls by %d and its shadow price is %d. Enter the change in the optimal objective, using a negative number for a decrease, when the old optimum was %d.", delta, price, opt), -price*delta,
			fmt.Sprintf("The change is -%d·%d = %d, giving an objective of %d.", price, delta, -price*delta, opt-price*delta))
	default:
		p1, p2 := limit(2, 9), limit(1, 8)
		return at(kp, fmt.Sprintf("One constraint has shadow price %d and another has shadow price %d. Enter the difference of the larger minus the smaller.", p1, p2), absInt(p1-p2),
			fmt.Sprintf("The larger price is %d, the smaller is %d and the difference is %d.", maxInt(p1, p2), minInt(p1, p2), absInt(p1-p2)))
	}
}

// ---------------------------------------------------------------------------
// Sensitivity ranges
// ---------------------------------------------------------------------------

type lpSensitivityGen struct{}

func (g *lpSensitivityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		low, high := limit(1, 9), limit(10, 20)
		return at(kp, fmt.Sprintf("The allowable range for a right-hand side is [%d, %d]. Enter its width.", low, high), high-low,
			fmt.Sprintf("The width is %d - %d = %d.", high, low, high-low))
	case 2:
		c, a, y := limit(2, 8), limit(1, 4), limit(1, 5)
		newC := c + limit(1, 6)
		return at(kp, fmt.Sprintf("A reduced cost is c - a y with c = %d, a = %d and y = %d. After c rises to %d, enter the new reduced cost.", newC, a, y, c), newC-a*y,
			fmt.Sprintf("The new reduced cost is %d - %d·%d = %d. A negative value would make the variable enter the basis.", newC, a, y, newC-a*y))
	default:
		opt, price, delta := limit(10, 50), limit(1, 5), limit(1, 4)
		return at(kp, fmt.Sprintf("A coefficient change raises the optimal objective by shadow price %d times a right-hand change of %d. Enter the new optimum when the old optimum was %d.", price, delta, opt), opt+price*delta,
			fmt.Sprintf("The rise is %d·%d = %d, so the new optimum is %d.", price, delta, price*delta, opt+price*delta))
	}
}

// ---------------------------------------------------------------------------
// Karush-Kuhn-Tucker conditions
// ---------------------------------------------------------------------------

type kktGen struct{}

func (g *kktGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		gradF, mult, gradG := limit(2, 12), limit(1, 5), limit(1, 6)
		return at(kp, fmt.Sprintf("A candidate has gradient of the objective equal to %d and multiplier %d times a constraint gradient of %d. Enter the stationarity residual grad f + mult·grad g.", gradF, mult, gradG), gradF+mult*gradG,
			fmt.Sprintf("The residual is %d + %d·%d = %d; KKT needs it to be zero.", gradF, mult, gradG, gradF+mult*gradG))
	case 2:
		a, k := limit(2, 9), limit(1, 9)
		lambda := k * a
		return at(kp, fmt.Sprintf("Stationarity for a KKT candidate reads %d x + %d = 0, where the second term is a multiplier times a constraint gradient. Enter the integer solution x.", a, lambda), -k,
			fmt.Sprintf("Solving %d x = -%d gives x = %d.", a, lambda, -k))
	default:
		mult, slack := limit(1, 5), limit(0, 6)
		return at(kp, fmt.Sprintf("A KKT candidate has multiplier %d and the matching constraint slack %d. Enter the complementary product.", mult, slack), mult*slack,
			fmt.Sprintf("The product is %d·%d = %d; one of the two must be zero at an optimum.", mult, slack, mult*slack))
	}
}

// ---------------------------------------------------------------------------
// Slater's condition
// ---------------------------------------------------------------------------

type slaterGen struct{}

func (g *slaterGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		a := limit(3, 40)
		return at(kp, fmt.Sprintf("A constraint is x < %d and x must be an integer. Enter the largest integer x that satisfies it strictly.", a), a-1,
			fmt.Sprintf("The largest integer below %d is %d.", a, a-1))
	case 2:
		a, x := limit(2, 15), limit(1, 5)
		return at(kp, fmt.Sprintf("A constraint is %d - x > 0 and the candidate is x = %d. Enter its strict feasibility margin.", a, x), a-x,
			fmt.Sprintf("The margin is %d - %d = %d; a positive value means the point is strictly feasible.", a, x, a-x))
	default:
		left, right := limit(1, 9), limit(10, 20)
		return at(kp, fmt.Sprintf("A convex feasible set contains the two points %d and %d. Enter an integer interior point strictly between them.", left, right), (left+right)/2,
			fmt.Sprintf("The midpoint (%d + %d)/2 = %d lies strictly between the endpoints.", left, right, (left+right)/2))
	}
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

func registerConvex(reg *generator.Registry) {
	reg.Register("optim.lp_standard_form", &lpStandardFormGen{})
	reg.Register("optim.basic_solution", &basicSolutionGen{})
	reg.Register("optim.simplex_pivot", &simplexPivotGen{})
	reg.Register("optim.lp_dual", &lpDualGen{})
	reg.Register("optim.dual_function", &dualFunctionGen{})
	reg.Register("optim.weak_duality", &weakDualityGen{})
	reg.Register("optim.strong_duality", &strongDualityGen{})
	reg.Register("optim.complementary_slackness", &complementarySlacknessGen{})
	reg.Register("optim.shadow_price", &shadowPriceGen{})
	reg.Register("optim.lp_sensitivity", &lpSensitivityGen{})
	reg.Register("optim.kkt", &kktGen{})
	reg.Register("optim.slater", &slaterGen{})
}
