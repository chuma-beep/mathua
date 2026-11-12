// Package optimization generates the Wave 11 Optimization domain. Every
// family computes a number, so no family asks a learner to name a theorem and
// no family has a fixed yes or no answer.
package optimization

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func ordinal(ctx generator.GeneratorContext) int {
	if ctx.KPOrdinal >= 1 && ctx.KPOrdinal <= 3 {
		return ctx.KPOrdinal
	}
	return rand.Intn(3) + 1
}

func at(kp int, question string, answer int, explanation string) generator.Problem {
	value := strconv.Itoa(answer)
	return generator.Problem{
		Question:    question,
		Answer:      value,
		Explanation: explanation,
		Facts:       map[string]string{"answer": value},
		KPOrdinal:   kp,
	}
}

func limit(low, high int) int { return low + rand.Intn(high-low+1) }

// ---------------------------------------------------------------------------
// Anatomy of an optimization problem
// ---------------------------------------------------------------------------

type anatomyGen struct{}

func (g *anatomyGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		a, b := limit(2, 9), limit(2, 9)
		x, y := limit(1, 6), limit(1, 6)
		answer := a*x*x + b*y
		return at(kp, fmt.Sprintf("A model maximises f(x, y) = %d x^2 + %d y. Evaluate the objective at (x, y) = (%d, %d). (enter a number)", a, b, x, y), answer,
			fmt.Sprintf("Substitute both values: %d·%d² + %d·%d = %d.", a, x, b, y, answer))
	case 2:
		a, b := limit(2, 9), limit(2, 9)
		x, y := limit(1, 7), limit(1, 7)
		c := limit(5, 40)
		answer := a*x + b*y
		return at(kp, fmt.Sprintf("A model has the constraint %d x + %d y ≤ %d. For the point (%d, %d), by how much does the left side fall short of the limit? Use a negative number when the point violates the constraint. (enter a number)", a, b, c, x, y), c-(a*x+b*y),
			fmt.Sprintf("The left side is %d·%d + %d·%d = %d, so the slack is %d - %d = %d.", a, x, b, y, answer, c, answer, c-answer))
	default:
		a, b, c := limit(2, 8), limit(2, 8), limit(1, 9)
		x, y, z := limit(1, 5), limit(1, 5), limit(1, 5)
		answer := a*x + b*y + c*z
		return at(kp, fmt.Sprintf("A production plan chooses x, y and z to maximise %d x + %d y + %d z. Evaluate the objective at (%d, %d, %d). (enter a number)", a, b, c, x, y, z), answer,
			fmt.Sprintf("Each coefficient multiplies its own variable: %d·%d + %d·%d + %d·%d = %d.", a, x, b, y, c, z, answer))
	}
}

// ---------------------------------------------------------------------------
// Feasible and infeasible sets
// ---------------------------------------------------------------------------

type feasibleSetGen struct{}

func (g *feasibleSetGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		a, b := limit(2, 9), limit(2, 9)
		x, y := limit(1, 8), limit(1, 8)
		c := limit(10, 60)
		return at(kp, fmt.Sprintf("Test the point (%d, %d) against the single constraint %d x + %d y ≤ %d by computing the slack. (enter a number)", x, y, a, b, c), c-(a*x+b*y),
			fmt.Sprintf("The slack is %d - (%d·%d + %d·%d) = %d. A negative slack means the point is infeasible.", c, a, x, b, y, c-a*x-b*y))
	case 2:
		low := limit(1, 10)
		width := limit(4, 30)
		return at(kp, fmt.Sprintf("A feasible set in one variable is defined by 2x ≥ %d and 2x ≤ %d. How wide is the feasible interval in units of x? (enter a number)", 2*low, 2*(low+width)), width,
			fmt.Sprintf("Divide both bounds by 2: the interval is [%d, %d], whose width is %d.", low, low+width, width))
	default:
		a, b := limit(2, 9), limit(10, 60)
		return at(kp, fmt.Sprintf("A feasible set is x ≥ 0 and %d x ≤ %d with x a whole number. How many integer values of x are feasible? (enter a number)", a, b), b/a+1,
			fmt.Sprintf("The largest feasible integer is ⌊%d / %d⌋ = %d, and x = 0 is allowed, so there are %d values.", b, a, b/a, b/a+1))
	}
}

// ---------------------------------------------------------------------------
// Local and global minima
// ---------------------------------------------------------------------------

type localGlobalGen struct{}

func (g *localGlobalGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		h, a := limit(1, 8), limit(1, 6)
		p, q := limit(1, 10), limit(1, 10)
		diff := a * ((p-h)*(p-h) - (q-h)*(q-h))
		return at(kp, fmt.Sprintf("For f(x) = %d (x - %d)^2, compute f(%d) - f(%d). (enter a number)", a, h, p, q), diff,
			fmt.Sprintf("The squares are %d and %d, so the difference is %d·(%d - %d) = %d.", (p-h)*(p-h), (q-h)*(q-h), a, (p-h)*(p-h), (q-h)*(q-h), diff))
	case 2:
		h, k := limit(1, 9), limit(1, 12)
		answer := k
		return at(kp, fmt.Sprintf("A function has the form f(x) = 4 (x - %d)^2 + %d. What is the global minimum value? (enter a number)", h, k), answer,
			fmt.Sprintf("The square is never negative, so the smallest value is the constant %d, reached at x = %d.", k, h))
	default:
		a, b, c := limit(1, 4), limit(1, 4), limit(1, 4)
		base := limit(1, 6)
		values := []int{a, b, c}
		best := values[0]
		for _, v := range values {
			if v < best {
				best = v
			}
		}
		return at(kp, fmt.Sprintf("A function takes the values %d, %d and %d at three candidate points numbered 1 to 3. What is the smallest value among them? (enter a number)", a+base, b+base, c+base), best+base,
			fmt.Sprintf("The smallest listed value is %d, so that candidate is the best of the three.", best+base))
	}
}

// ---------------------------------------------------------------------------
// Convex sets
// ---------------------------------------------------------------------------

type convexSetGen struct{}

func (g *convexSetGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		mx, my := limit(2, 12), limit(2, 12)
		dx, dy := limit(1, 6), limit(1, 6)
		x1, y1 := mx-dx, my-dy
		x2, y2 := mx+dx, my+dy
		return at(kp, fmt.Sprintf("Take the convex combination of (%d, %d) and (%d, %d) with weight 0.5 on the first point. What is the sum of the two coordinates of the result? (enter a number)", x1, y1, x2, y2), mx+my,
			fmt.Sprintf("The midpoint is ((%d+%d)/2, (%d+%d)/2) = (%d, %d). Its coordinate sum is %d.", x1, x2, y1, y2, mx, my, mx+my))
	case 2:
		a, b := limit(2, 9), limit(2, 9)
		x, y := limit(1, 8), limit(1, 8)
		return at(kp, fmt.Sprintf("The line %d x + %d y = c separates a set from a point. For the point (%d, %d), what is the value of the left side? (enter a number)", a, b, x, y), a*x+b*y,
			fmt.Sprintf("Substitute the point: %d·%d + %d·%d = %d.", a, x, b, y, a*x+b*y))
	default:
		a, b := limit(1, 9), limit(1, 9)
		x, y := limit(1, 9), limit(1, 9)
		return at(kp, fmt.Sprintf("A set consists of every point on one side of %d x + %d y ≤ 12. Compute the left side at (%d, %d) and report how much room the point has to the boundary. (enter a number)", a, b, x, y), 12-(a*x+b*y),
			fmt.Sprintf("The margin is 12 - (%d·%d + %d·%d) = %d. A non-negative margin places the point inside.", a, x, b, y, 12-a*x-b*y))
	}
}

// ---------------------------------------------------------------------------
// Strict convexity and uniqueness
// ---------------------------------------------------------------------------

type strictConvexityGen struct{}

func (g *strictConvexityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		a := limit(2, 40)
		return at(kp, fmt.Sprintf("A twice differentiable function has constant second derivative %d on an interval. What is its second derivative at the midpoint? (enter a number)", a), a,
			fmt.Sprintf("A constant second derivative stays %d at every point of the interval, and a positive value makes the function strictly convex.", a))
	case 2:
		h, k := limit(1, 12), limit(1, 12)
		return at(kp, fmt.Sprintf("For f(x) = 7 (x - %d)^2 + %d, at which x does the unique minimiser occur? (enter a number)", h, k), h,
			fmt.Sprintf("The squared term is zero only at x = %d, and it is positive everywhere else, so the minimiser is unique.", h))
	default:
		a, gap := limit(1, 9), limit(1, 9)
		b := a + 2*gap
		mid := a + gap
		answer := -gap * gap
		return at(kp, fmt.Sprintf("For f(x) = x^2 on the interval from %d to %d, compute f(midpoint) minus the average of the endpoint values. (enter a number)", a, b), answer,
			fmt.Sprintf("The midpoint is %d, so f(mid) = %d and the endpoint average is (%d + %d)/2 = %d. The difference is %d.", mid, mid*mid, a*a, b*b, (a*a+b*b)/2, answer))
	}
}

// ---------------------------------------------------------------------------
// Subgradients
// ---------------------------------------------------------------------------

type subgradientGen struct{}

func (g *subgradientGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		a, x := limit(1, 9), limit(1, 9)
		g := limit(-4, 4)
		margin := g*(x-a) - absInt(x-a)
		return at(kp, fmt.Sprintf("For f(x) = |x - %d| and the point x = %d, test g = %d by computing g(x - a) - |x - a|. (enter a number)", a, x, g), margin,
			fmt.Sprintf("The supporting-line test is g(x-a) - |x-a| = %d·%d - %d = %d. A non-positive value means g is a valid subgradient.", g, x-a, absInt(x-a), margin))
	case 2:
		a, c := limit(1, 9), limit(1, 6)
		return at(kp, fmt.Sprintf("For f(x) = %d |x - %d|, what is the largest subgradient at x = %d? (enter a number)", c, a, a), c,
			fmt.Sprintf("At the kink the subdifferential is the interval [-%d, %d], so the largest member is %d.", c, c, c))
	default:
		a := limit(1, 10)
		x := a + limit(1, 8)
		return at(kp, fmt.Sprintf("For f(x) = |x - %d| at the point x = %d, the supporting line has slope 1. What is its intercept, f(x) - 1·x? (enter a number)", a, x), absInt(x-a)-x,
			fmt.Sprintf("The value is %d and the intercept is %d - %d = %d.", absInt(x-a), absInt(x-a), x, absInt(x-a)-x))
	}
}

// ---------------------------------------------------------------------------
// Epigraphs and sublevel sets
// ---------------------------------------------------------------------------

type epigraphGen struct{}

func (g *epigraphGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		x, t := limit(-5, 5), limit(-5, 15)
		return at(kp, fmt.Sprintf("For f(x) = x^2, how far above the graph is the point (%d, %d)? Compute t - x^2. (enter a number)", x, t), t-x*x,
			fmt.Sprintf("At x = %d the graph value is %d, so t - f(x) = %d. A non-negative result places the point in the epigraph.", x, x*x, t-x*x))
	case 2:
		r := limit(4, 100)
		answer := 0
		for (answer+1)*(answer+1) <= r {
			answer++
		}
		return at(kp, fmt.Sprintf("The sublevel set is x^2 ≤ %d. What is the largest whole number x in it? (enter a number)", r), answer,
			fmt.Sprintf("The largest square not exceeding %d is %d·%d, so the largest whole x is %d.", r, answer, answer, answer))
	default:
		a, b := limit(2, 9), limit(2, 9)
		return at(kp, fmt.Sprintf("For f(x) = %d x^2 + %d x, compute the second difference f(1) + f(-1) - 2 f(0). (enter a number)", a, b), 2*a,
			fmt.Sprintf("f(1) + f(-1) = %d + %d and 2 f(0) = 0, so the second difference is %d.", a+b, a-b, 2*a))
	}
}

// ---------------------------------------------------------------------------
// Coercivity and existence
// ---------------------------------------------------------------------------

type coercivityGen struct{}

func (g *coercivityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		a := limit(2, 40)
		return at(kp, fmt.Sprintf("For f(x) = %d x^2 + 3x - 7, what is the limit of f(x)/x^2 as |x| grows without bound? (enter a number)", a), a,
			fmt.Sprintf("Dividing by x^2 leaves %d + 3/x - 7/x^2, whose limit is %d. A positive limit confirms coercivity.", a, a))
	case 2:
		a, c := limit(1, 6), limit(1, 10)
		h := limit(1, 8)
		answer := c
		return at(kp, fmt.Sprintf("For f(x) = %d (x - %d)^2 + %d, what is the minimum value? (enter a number)", a, h, c), answer,
			fmt.Sprintf("The square vanishes at x = %d and is positive elsewhere, so the minimum is %d.", h, c))
	default:
		a, k, m := limit(1, 6), limit(1, 10), limit(1, 30)
		x := 0
		for a*x*x+k <= m {
			x++
		}
		return at(kp, fmt.Sprintf("For f(x) = %d x^2 + %d with x ≥ 0, what is the smallest whole x for which f(x) exceeds %d? (enter a number)", a, k, m), x,
			fmt.Sprintf("The smallest x with %d x^2 + %d > %d is %d.", a, k, m, x))
	}
}

func absInt(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

// ---------------------------------------------------------------------------
// First-order optimality
// ---------------------------------------------------------------------------

type firstOrderGen struct{}

func (g *firstOrderGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		a, b := limit(1, 9), limit(1, 9)
		x, y := limit(1, 6), limit(1, 6)
		return at(kp, fmt.Sprintf("For f(x, y) = %d x^2 + %d y^2, what is the first component of the gradient at (%d, %d)? (enter a number)", a, b, x, y), 2*a*x,
			fmt.Sprintf("The partial derivative in x is %d x, which at x = %d equals %d.", 2*a, x, 2*a*x))
	case 2:
		a, k := limit(1, 8), limit(1, 10)
		b := 2 * a * k
		return at(kp, fmt.Sprintf("For f(x) = %d x^2 - %d x, at which x does the gradient vanish? (enter a number)", a, b), k,
			fmt.Sprintf("Setting %d x - %d = 0 gives x = %d.", 2*a, b, k))
	default:
		gx, gy := limit(-6, 6), limit(-6, 6)
		dx, dy := limit(-4, 4), limit(-4, 4)
		return at(kp, fmt.Sprintf("At a point the gradient is (%d, %d). For the direction (%d, %d), what is the directional derivative? (enter a number)", gx, gy, dx, dy), gx*dx+gy*dy,
			fmt.Sprintf("The directional derivative is the dot product %d·%d + %d·%d = %d. A negative value means the direction descends.", gx, dx, gy, dy, gx*dx+gy*dy))
	}
}

// ---------------------------------------------------------------------------
// Second-order optimality and saddle points
// ---------------------------------------------------------------------------

type secondOrderGen struct{}

func (g *secondOrderGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		a, b := limit(1, 9), limit(1, 9)
		return at(kp, fmt.Sprintf("For f(x, y) = %d x^2 + %d y^2, what is the Hessian entry in the x row and x column at any point? (enter a number)", a, b), 2*a,
			fmt.Sprintf("The second partial in x is %d, so that Hessian entry is %d everywhere.", 2*a, 2*a))
	case 2:
		a, b := limit(1, 9), limit(1, 9)
		det := 4 * a * b
		return at(kp, fmt.Sprintf("At a stationary point of f(x, y) = %d x^2 + %d y^2, compute the determinant of the Hessian. (enter a number)", a, b), det,
			fmt.Sprintf("The Hessian is diagonal with entries %d and %d, so its determinant is %d.", 2*a, 2*b, det))
	default:
		a, b, x, y := limit(1, 6), limit(1, 6), limit(1, 5), limit(1, 5)
		answer := a*x*x - b*y*y
		return at(kp, fmt.Sprintf("For f(x, y) = %d x^2 - %d y^2, evaluate f at the point (%d, %d). (enter a number)", a, b, x, y), answer,
			fmt.Sprintf("Substitute both coordinates: %d·%d² - %d·%d² = %d. The Hessian has one positive and one negative eigenvalue, so the origin is a saddle.", a, x, b, y, answer))
	}
}

func Register(reg *generator.Registry) {
	registerFoundations(reg)
	registerUnconstrained(reg)
	registerConvex(reg)
	registerConstrained(reg)
	registerDiscrete(reg)
	registerWave12(reg)
}

func registerFoundations(reg *generator.Registry) {
	reg.Register("optim.problem.anatomy", &anatomyGen{})
	reg.Register("optim.feasible.set", &feasibleSetGen{})
	reg.Register("optim.local_global", &localGlobalGen{})
	reg.Register("optim.convex_set", &convexSetGen{})
	reg.Register("optim.strict_convexity", &strictConvexityGen{})
	reg.Register("optim.subgradient", &subgradientGen{})
	reg.Register("optim.epigraph", &epigraphGen{})
	reg.Register("optim.coercivity", &coercivityGen{})
	reg.Register("optim.first_order", &firstOrderGen{})
	reg.Register("optim.second_order", &secondOrderGen{})
}
