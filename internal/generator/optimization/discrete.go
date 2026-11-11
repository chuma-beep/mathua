package optimization

import (
	"fmt"

	"github.com/chuma-beep/mathua/internal/generator"
)

// Discrete optimization. The instances are small enough to solve exactly, and
// every answer is the optimum, a bound or the value of one table entry.

// ---------------------------------------------------------------------------
// Integer programming
// ---------------------------------------------------------------------------

type integerProgramGen struct{}

func (g *integerProgramGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		whole, frac := limit(1, 9), limit(1, 9)
		return at(kp, fmt.Sprintf("A linear program has the continuous optimum %d.%d. Enter the smallest integer that is not below it.", whole, frac), whole+1,
			fmt.Sprintf("The value lies strictly between %d and %d, so the smallest integer at or above it is %d.", whole, whole+1, whole+1))
	case 2:
		whole, frac := limit(1, 9), limit(1, 9)
		return at(kp, fmt.Sprintf("A solution is %d.%d. Enter its fractional part in tenths, as an integer from 0 to 9.", whole, frac), frac,
			fmt.Sprintf("The fractional part is 0.%d, which is %d tenths, so the solution is not integral.", frac, frac))
	default:
		a, b := limit(2, 9), limit(2, 9)
		x, y := limit(1, 8), limit(1, 8)
		return at(kp, fmt.Sprintf("An integer program maximises %d x + %d y. Enter the objective at the integer point (%d, %d).", a, b, x, y), a*x+b*y,
			fmt.Sprintf("The value is %d·%d + %d·%d = %d.", a, x, b, y, a*x+b*y))
	}
}

// ---------------------------------------------------------------------------
// LP relaxation
// ---------------------------------------------------------------------------

type lpRelaxationGen struct{}

func (g *lpRelaxationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		a, b := limit(2, 6), limit(2, 6)
		c, d := a*limit(2, 8), b*limit(2, 8)
		return at(kp, fmt.Sprintf("Maximize x + y subject to %d x ≤ %d and %d y ≤ %d. Enter the relaxed optimum.", a, c, b, d), c/a+d/b,
			fmt.Sprintf("The constraints give x = %d and y = %d, so the relaxed optimum is %d.", c/a, d/b, c/a+d/b))
	case 2:
		relaxed := limit(10, 40)
		integer := relaxed - limit(1, 5)
		return at(kp, fmt.Sprintf("The LP relaxation optimum is %d and the integer optimum is %d. Enter the integrality gap.", relaxed, integer), relaxed-integer,
			fmt.Sprintf("The gap is %d - %d = %d.", relaxed, integer, relaxed-integer))
	default:
		relaxed, gap := limit(10, 50), limit(1, 7)
		return at(kp, fmt.Sprintf("An integer program has relaxation value %d. Enter the integer optimum when the integrality gap is %d.", relaxed, gap), relaxed-gap,
			fmt.Sprintf("The integer optimum is %d - %d = %d.", relaxed, gap, relaxed-gap))
	}
}

// ---------------------------------------------------------------------------
// Branch and bound
// ---------------------------------------------------------------------------

type branchBoundGen struct{}

func (g *branchBoundGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		whole, frac := limit(1, 9), limit(1, 9)
		return at(kp, fmt.Sprintf("Branch and bound branches on x = %d.%d. Enter the lower branch bound, the largest integer not above the value.", whole, frac), whole,
			fmt.Sprintf("The lower child keeps x ≤ %d.%d, so its integer bound is %d.", whole, frac, whole))
	case 2:
		whole, frac := limit(10, 40), limit(1, 9)
		return at(kp, fmt.Sprintf("A node has the fractional relaxation value %d.%d. Branch and bound rounds the value down for an integer upper bound. Enter that bound.", whole, frac), whole,
			fmt.Sprintf("Rounding %d.%d down gives %d, which the node passes on as its integer upper bound.", whole, frac, whole))
	default:
		incumbent, bound := limit(20, 60), limit(1, 15)
		return at(kp, fmt.Sprintf("An incumbent is %d and a node bound is %d. Enter the amount by which the incumbent beats the bound.", incumbent, incumbent-bound), bound,
			fmt.Sprintf("The incumbent exceeds the bound by %d, so the node can be pruned.", bound))
	}
}

// ---------------------------------------------------------------------------
// Cutting planes
// ---------------------------------------------------------------------------

type cuttingPlaneGen struct{}

func (g *cuttingPlaneGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		a, x, b := limit(1, 5), limit(1, 8), 0
		slack := limit(1, 8)
		b = a*x + slack
		return at(kp, fmt.Sprintf("A cut is %d x ≤ %d and the current point is x = %d. Enter the cut slack.", a, b, x), slack,
			fmt.Sprintf("The left side is %d, so the slack is %d - %d = %d.", a*x, b, a*x, slack))
	case 2:
		whole, frac := limit(1, 8), limit(1, 9)
		return at(kp, fmt.Sprintf("A fractional solution is x = %d.%d. A rounding cut takes the floor of the value. Enter that floor.", whole, frac), whole,
			fmt.Sprintf("The largest integer not above %d.%d is %d.", whole, frac, whole))
	default:
		a, x, b := limit(1, 5), limit(1, 8), 0
		extra := limit(1, 9)
		b = a*x - extra
		return at(kp, fmt.Sprintf("A cut is %d x ≤ %d and the current point is x = %d. Enter the violation, which is the left side minus the right side.", a, b, x), a*x-b,
			fmt.Sprintf("The left side is %d and the right side is %d, so the cut is violated by %d.", a*x, b, a*x-b))
	}
}

// ---------------------------------------------------------------------------
// The knapsack problem
// ---------------------------------------------------------------------------

type knapsackGen struct{}

func (g *knapsackGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		w1, w2, v1, v2 := limit(1, 4), limit(1, 4), limit(2, 8), limit(2, 8)
		capacity := w1 + w2 - limit(0, 1)
		best := 0
		if w1 <= capacity {
			best = maxInt(best, v1)
		}
		if w2 <= capacity {
			best = maxInt(best, v2)
		}
		if w1+w2 <= capacity {
			best = maxInt(best, v1+v2)
		}
		return at(kp, fmt.Sprintf("A knapsack has items (weight %d, value %d) and (weight %d, value %d) with capacity %d. Each item is taken at most once. Enter the maximum value.", w1, v1, w2, v2, capacity), best,
			fmt.Sprintf("The best feasible selection has value %d.", best))
	case 2:
		old, candidate := limit(5, 20), limit(2, 25)
		return at(kp, fmt.Sprintf("A knapsack table entry is old = %d. Including the item gives the alternative %d. Enter the updated entry, the larger of the two.", old, candidate), maxInt(old, candidate),
			fmt.Sprintf("The update takes the larger of %d and %d, which is %d.", old, candidate, maxInt(old, candidate)))
	default:
		weight, value, capacity := limit(2, 5), limit(2, 9), limit(1, 4)
		full := capacity / weight
		remainder := capacity - full*weight
		fraction := remainder * value / weight
		return at(kp, fmt.Sprintf("An item has weight %d and value %d and the capacity is %d. The fractional knapsack takes whole copies and a fraction of the next. Enter the fractional value, rounded down.", weight, value, capacity), full*value+fraction,
			fmt.Sprintf("Whole copies give %d, and the remaining capacity %d adds %d, so the bound is %d.", full*value, remainder, fraction, full*value+fraction))
	}
}

// ---------------------------------------------------------------------------
// The transportation problem
// ---------------------------------------------------------------------------

type transportationGen struct{}

func (g *transportationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		supply, demand := limit(10, 40), limit(1, 9)
		return at(kp, fmt.Sprintf("A transportation table has total supply %d and total demand %d. Enter the amount that must be added to balance the table.", supply, supply-demand), demand,
			fmt.Sprintf("The table is short by %d - %d = %d, so a dummy demand absorbs it.", supply, supply-demand, demand))
	case 2:
		supply, demand := limit(5, 20), limit(5, 20)
		return at(kp, fmt.Sprintf("A north-west corner rule allocates the smaller of supply %d and demand %d. Enter that allocation.", supply, demand), minInt(supply, demand),
			fmt.Sprintf("The allocation is the smaller entry, %d.", minInt(supply, demand)))
	default:
		amount, cost := limit(2, 12), limit(2, 9)
		return at(kp, fmt.Sprintf("An allocation of %d units moves on a route with unit cost %d. Enter the route cost.", amount, cost), amount*cost,
			fmt.Sprintf("The cost is %d·%d = %d.", amount, cost, amount*cost))
	}
}

// ---------------------------------------------------------------------------
// Set cover
// ---------------------------------------------------------------------------

type setCoverGen struct{}

func (g *setCoverGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		covered, missing := limit(1, 30), limit(1, 30)
		total := covered + missing
		return at(kp, fmt.Sprintf("A universe has %d elements and a proposed collection covers %d of them. Enter the number left uncovered.", total, covered), total-covered,
			fmt.Sprintf("The uncovered count is %d - %d = %d, so the collection is not a cover.", total, covered, total-covered))
	case 2:
		a, b := limit(3, 40), limit(3, 40)
		overlap := limit(1, minInt(a, b)-1)
		return at(kp, fmt.Sprintf("Two sets of sizes %d and %d share %d elements. Enter the size of their union.", a, b, overlap), a+b-overlap,
			fmt.Sprintf("The union size is %d + %d - %d = %d.", a, b, overlap, a+b-overlap))
	default:
		cover, vertex := limit(2, 8), limit(2, 8)
		return at(kp, fmt.Sprintf("A set cover has size %d and the corresponding vertex cover has size %d. Enter the difference of the larger minus the smaller.", cover, vertex), absInt(cover-vertex),
			fmt.Sprintf("The difference is |%d - %d| = %d.", cover, vertex, absInt(cover-vertex)))
	}
}

// ---------------------------------------------------------------------------
// The integrality gap
// ---------------------------------------------------------------------------

type integralityGapGen struct{}

func (g *integralityGapGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		relaxed, integer := limit(12, 50), limit(1, 9)
		return at(kp, fmt.Sprintf("A relaxation value is %d and the integer optimum is %d. Enter the integrality gap.", relaxed, relaxed-integer), integer,
			fmt.Sprintf("The gap is %d - %d = %d.", relaxed, relaxed-integer, integer))
	case 2:
		g1, g2 := limit(4, 30), limit(1, 4)
		return at(kp, fmt.Sprintf("Two formulations have integrality gaps %d and %d. Enter the difference of the larger minus the smaller.", g1, g1-g2), g2,
			fmt.Sprintf("The larger gap is %d and the smaller is %d, so the difference is %d.", g1, g1-g2, g2))
	default:
		relaxed, gap := limit(10, 40), limit(1, 6)
		return at(kp, fmt.Sprintf("A relaxation value is %d and the integer optimum is %d lower. Rounding the relaxation up gives an upper bound. Enter the difference between that bound and the integer optimum.", relaxed, gap), gap,
			fmt.Sprintf("The rounded bound is %d and the integer optimum is %d, so the difference is %d.", relaxed, relaxed-gap, gap))
	}
}

func registerDiscrete(reg *generator.Registry) {
	reg.Register("optim.integer_program", &integerProgramGen{})
	reg.Register("optim.lp_relaxation", &lpRelaxationGen{})
	reg.Register("optim.branch_bound", &branchBoundGen{})
	reg.Register("optim.cutting_plane", &cuttingPlaneGen{})
	reg.Register("optim.knapsack", &knapsackGen{})
	reg.Register("optim.transportation", &transportationGen{})
	reg.Register("optim.set_cover", &setCoverGen{})
	reg.Register("optim.integrality_gap", &integralityGapGen{})
}
