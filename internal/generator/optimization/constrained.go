package optimization

import (
	"fmt"

	"github.com/chuma-beep/mathua/internal/generator"
)

// Constrained optimization. The families compute steps, multipliers, residual
// penalties and gaps from small integer instances so every answer is exact.

// ---------------------------------------------------------------------------
// Equality-constrained optimization
// ---------------------------------------------------------------------------

type equalityLagrangeGen struct{}

func (g *equalityLagrangeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		c := 2 * limit(1, 40)
		return at(kp, fmt.Sprintf("Minimize x^2 + y^2 subject to x + y = %d. Enter the optimal x.", c), c/2,
			fmt.Sprintf("The Lagrange condition gives 2x = 2y, so x = y and 2x = %d, hence x = %d.", c, c/2))
	case 2:
		x, a, mult := limit(1, 9), limit(1, 9), limit(1, 6)
		return at(kp, fmt.Sprintf("For f(x) = x^2 and the constraint x - %d = 0, the stationarity residual is 2x + λ. At x = %d and λ = %d, enter the residual.", a, x, mult), 2*x+mult,
			fmt.Sprintf("The residual is 2·%d + %d = %d; a constrained optimum needs it to vanish.", x, mult, 2*x+mult))
	default:
		c := 2 * limit(1, 40)
		return at(kp, fmt.Sprintf("Minimize x^2 + y^2 subject to x + y = %d. Enter the Lagrange multiplier λ at the optimum.", c), c,
			fmt.Sprintf("Stationarity gives 2x = 2y = λ, so λ = x + y = %d.", c))
	}
}

// ---------------------------------------------------------------------------
// Active and inactive constraints
// ---------------------------------------------------------------------------

type activeConstraintGen struct{}

func (g *activeConstraintGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		a, b, x, y := limit(1, 5), limit(1, 5), limit(1, 8), limit(1, 8)
		slack := limit(0, 6)
		return at(kp, fmt.Sprintf("The constraint is %d x + %d y ≤ %d and the candidate is (%d, %d). Enter its slack.", a, b, a*x+b*y+slack, x, y), slack,
			fmt.Sprintf("ax + by = %d, so the slack is %d - %d = %d.", a*x+b*y, a*x+b*y+slack, a*x+b*y, slack))
	case 2:
		x, y := limit(1, 6), limit(1, 6)
		s1 := limit(1, 5) + x + y
		s2 := limit(1, 5) + 2*x + y
		answer := minInt(s1-(x+y), s2-(2*x+y))
		return at(kp, fmt.Sprintf("At (%d, %d) the constraints x + y ≤ %d and 2x + y ≤ %d are tested. Enter the smaller slack, which measures the tightest constraint.", x, y, s1, s2), answer,
			fmt.Sprintf("The slacks are %d and %d, so the smaller is %d.", s1-(x+y), s2-(2*x+y), answer))
	default:
		a, x := limit(2, 9), limit(1, 7)
		return at(kp, fmt.Sprintf("The constraint %d x ≤ b is active at x = %d. Enter b.", a, x), a*x,
			fmt.Sprintf("Active means equality, so b = %d·%d = %d.", a, x, a*x))
	}
}

// ---------------------------------------------------------------------------
// Quadratic penalty method
// ---------------------------------------------------------------------------

type penaltyGen struct{}

func (g *penaltyGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		x, a, mu := limit(1, 6), limit(2, 10), limit(1, 5)
		violation := x - a
		answer := x*x + mu*violation*violation
		return at(kp, fmt.Sprintf("A penalty method evaluates f(x) + μ g(x)^2 with f(x) = x^2, g(x) = x - %d and μ = %d. At x = %d, enter the penalized value.", a, mu, x), answer,
			fmt.Sprintf("x^2 = %d and g = %d, so the value is %d + %d·%d = %d.", x*x, violation, x*x, mu, violation*violation, answer))
	case 2:
		a := 2 * limit(1, 30)
		return at(kp, fmt.Sprintf("Minimize x^2 + (x - %d)^2 with μ = 1. Enter the minimizer.", a), a/2,
			fmt.Sprintf("The derivative 2x + 2(x - %d) vanishes at x = %d.", a, a/2))
	default:
		a := limit(4, 40)
		return at(kp, fmt.Sprintf("Minimize x^2 subject to x = %d with the penalty μ(x - %d)^2. The minimizer is %d μ/(1 + μ), so the constraint violation is %d/(1 + μ). Enter the smallest integer μ that makes the violation at most 1.", a, a, a, a), a-1,
			fmt.Sprintf("The violation %d/(1 + μ) is at most 1 when μ ≥ %d, so the smallest integer is %d.", a, a-1, a-1))
	}
}

// ---------------------------------------------------------------------------
// Logarithmic barrier method
// ---------------------------------------------------------------------------

type barrierGen struct{}

func (g *barrierGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		mu := limit(1, 40)
		return at(kp, fmt.Sprintf("Minimize x subject to x ≥ 1 with the barrier -μ log(x - 1). The central point is x = 1 + μ. For μ = %d, enter x.", mu), 1+mu,
			fmt.Sprintf("Setting the derivative 1 - μ/(x - 1) to zero gives x - 1 = %d, so x = %d.", mu, 1+mu))
	case 2:
		x, lambda := limit(2, 10), limit(1, 8)
		return at(kp, fmt.Sprintf("On the central path the multiplier and the distance to the boundary satisfy λ(x - 1) = μ. At x = %d and λ = %d, enter μ.", x, lambda), lambda*(x-1),
			fmt.Sprintf("μ = %d·%d = %d.", lambda, x-1, lambda*(x-1)))
	default:
		x, lambda := limit(2, 12), limit(1, 9)
		return at(kp, fmt.Sprintf("A barrier iterate has x = %d and multiplier λ = %d. Enter the central-path product λ(x - 1).", x, lambda), lambda*(x-1),
			fmt.Sprintf("The product is %d·%d = %d, which equals the barrier parameter on the central path.", lambda, x-1, lambda*(x-1)))
	}
}

// ---------------------------------------------------------------------------
// Augmented Lagrangian
// ---------------------------------------------------------------------------

type augmentedLagrangianGen struct{}

func (g *augmentedLagrangianGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		lambda, rho, gvalue := limit(1, 9), limit(1, 5), limit(-6, 6)
		return at(kp, fmt.Sprintf("An augmented Lagrangian update is λ_new = λ + ρ g(x). At λ = %d, ρ = %d and g = %d, enter λ_new.", lambda, rho, gvalue), lambda+rho*gvalue,
			fmt.Sprintf("λ_new = %d + %d·%d = %d.", lambda, rho, gvalue, lambda+rho*gvalue))
	case 2:
		x, lambda, rho, gvalue := limit(1, 9), limit(1, 5), 2*limit(1, 4), limit(-5, 5)
		answer := x*x + lambda*gvalue + (rho/2)*gvalue*gvalue
		return at(kp, fmt.Sprintf("Evaluate f + λ g + (ρ/2) g^2 with f = x^2, x = %d, λ = %d, ρ = %d and g = %d.", x, lambda, rho, gvalue), answer,
			fmt.Sprintf("The value is %d + %d·%d + %d·%d = %d.", x*x, lambda, gvalue, rho/2, gvalue*gvalue, answer))
	default:
		gvalue := limit(-30, 30)
		return at(kp, fmt.Sprintf("The augmented term is (ρ/2) g^2. Enter 2/ρ times that term when g = %d.", gvalue), gvalue*gvalue,
			fmt.Sprintf("(2/ρ)·(ρ/2)g^2 = g^2 = %d.", gvalue*gvalue))
	}
}

// ---------------------------------------------------------------------------
// Projected gradient descent
// ---------------------------------------------------------------------------

type projectedGradientGen struct{}

func (g *projectedGradientGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		low, high := limit(1, 5), limit(6, 15)
		x := limit(-6, 20)
		answer := minInt(maxInt(x, low), high)
		return at(kp, fmt.Sprintf("Project x = %d onto the interval [%d, %d]. Enter the projected value.", x, low, high), answer,
			fmt.Sprintf("The projection clamps below at %d and above at %d, so %d becomes %d.", low, high, x, answer))
	case 2:
		low, high := limit(0, 4), limit(8, 15)
		x, grad, t := limit(1, 12), limit(1, 6), limit(1, 3)
		raw := x - t*grad
		answer := minInt(maxInt(raw, low), high)
		return at(kp, fmt.Sprintf("Take a projected gradient step from x = %d with gradient %d and step size %d, projecting onto [%d, %d]. Enter the new point.", x, grad, t, low, high), answer,
			fmt.Sprintf("The raw step gives %d, which projects to %d.", raw, answer))
	default:
		low, high := limit(0, 4), limit(8, 15)
		x := limit(16, 30)
		return at(kp, fmt.Sprintf("Project x = %d onto [%d, %d] and enter the distance the projection moves the point.", x, low, high), x-high,
			fmt.Sprintf("The projection is %d, so the move is %d - %d = %d.", high, x, high, x-high))
	}
}

// ---------------------------------------------------------------------------
// Active-set methods
// ---------------------------------------------------------------------------

type activeSetGen struct{}

func (g *activeSetGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		x, y := limit(1, 8), limit(1, 8)
		b1, b2 := x+limit(0, 5), y+limit(0, 5)
		answer := maxInt(x-b1, y-b2)
		return at(kp, fmt.Sprintf("At (%d, %d) the working set tests x ≤ %d and y ≤ %d. Enter the largest violation, a positive number when a constraint is exceeded.", x, y, b1, b2), answer,
			fmt.Sprintf("The candidate exceeds the first bound by %d and the second by %d, so the largest violation is %d.", x-b1, y-b2, answer))
	case 2:
		c := 2 * limit(1, 40)
		return at(kp, fmt.Sprintf("An active-set step solves the equality-constrained subproblem: minimize x^2 + y^2 subject to x + y = %d. Enter x.", c), c/2,
			fmt.Sprintf("The two variables are equal at the optimum, so x = %d/2 = %d.", c, c/2))
	default:
		a, b, x, y := limit(1, 5), limit(1, 5), limit(1, 8), limit(1, 8)
		return at(kp, fmt.Sprintf("On a face with y fixed at %d, the objective is f(x, y) = x^2 + %d x + %d y. Enter the partial derivative in the free direction x at x = %d.", y, a, b, x), 2*x+a,
			fmt.Sprintf("The partial in x is 2·%d + %d = %d; zero would make the point optimal on the face.", x, a, 2*x+a))
	}
}

// ---------------------------------------------------------------------------
// The trust-region subproblem
// ---------------------------------------------------------------------------

type trustRegionSubproblemGen struct{}

func (g *trustRegionSubproblemGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		a, grad, delta := limit(2, 9), limit(1, 9), limit(1, 6)
		answer := -delta
		if grad < 0 {
			answer = delta
		}
		return at(kp, fmt.Sprintf("Minimize (%d/2)x^2 + %d x over |x| ≤ %d. The unconstrained minimizer lies outside the region, so the solution is on the boundary. Enter x.", a, grad, delta), answer,
			fmt.Sprintf("The unconstrained minimizer is -%d/%d, outside the radius, so the boundary point opposite the gradient is %d.", grad, a, answer))
	case 2:
		a1, k := limit(1, 6), limit(1, 6)
		grad := -k * (a1 + 1)
		return at(kp, fmt.Sprintf("Minimize (%d/2)x^2 + %d x with multiplier λ = 1. The solution is x = -g/(a + λ). Enter x.", a1, grad), k,
			fmt.Sprintf("x = -(%d)/(%d + 1) = %d.", grad, a1, k))
	default:
		p, q := limit(1, 8), limit(1, 8)
		return at(kp, fmt.Sprintf("A boundary solution has components %d and %d. Enter its squared trust-region radius.", p, q), p*p+q*q,
			fmt.Sprintf("The squared norm is %d^2 + %d^2 = %d, which the radius must match.", p, q, p*p+q*q))
	}
}

// ---------------------------------------------------------------------------
// The KKT matrix
// ---------------------------------------------------------------------------

type kktMatrixGen struct{}

func (g *kktMatrixGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		h1, h2 := limit(2, 9), limit(2, 9)
		return at(kp, fmt.Sprintf("A KKT matrix has Hessian block diag(%d, %d), one equality row (1, 1) and a zero corner. Enter its trace.", h1, h2), h1+h2,
			fmt.Sprintf("The trace is the sum of the diagonal entries: %d + %d + 0 = %d.", h1, h2, h1+h2))
	case 2:
		h1, h2 := limit(2, 9), limit(2, 9)
		return at(kp, fmt.Sprintf("The KKT matrix is [[%d, 0, 1], [0, %d, 1], [1, 1, 0]]. Enter its determinant.", h1, h2), -(h1 + h2),
			fmt.Sprintf("Expanding along the last column gives -(%d + %d) = %d, so the matrix is nonsingular.", h1, h2, -(h1+h2)))
	default:
		x1, x2 := limit(1, 9), limit(1, 9)
		return at(kp, fmt.Sprintf("A KKT system has equality row (1, 1) with right-hand side %d. Enter the first primal component x1 when the second component is %d.", x1+x2, x2), x1,
			fmt.Sprintf("The equality fixes x1 + x2 = %d, so x1 = %d - %d = %d.", x1+x2, x1+x2, x2, x1))
	}
}

// ---------------------------------------------------------------------------
// The duality gap
// ---------------------------------------------------------------------------

type dualityGapGen struct{}

func (g *dualityGapGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		p, gap := limit(30, 90), limit(1, 15)
		return at(kp, fmt.Sprintf("A primal value is %d and the gap is %d. Enter the dual value.", p, gap), p-gap,
			fmt.Sprintf("gap = primal - dual, so dual = %d - %d = %d.", p, gap, p-gap))
	case 2:
		k := limit(2, 40)
		return at(kp, fmt.Sprintf("A stopping rule halves the gap each iteration. Starting from a gap of %d and stopping at a tolerance of 1, enter the number of halvings needed.", 1<<k), k,
			fmt.Sprintf("%d halved %d times reaches 1, because 2^%d = %d.", 1<<k, k, k, 1<<k))
	default:
		g1, g2 := limit(5, 40), limit(1, 5)
		return at(kp, fmt.Sprintf("Two candidate pairs have gaps %d and %d. Enter the difference of the larger minus the smaller.", g1, g1-g2), g2,
			fmt.Sprintf("The larger gap is %d and the smaller is %d, so the difference is %d.", g1, g1-g2, g2))
	}
}

func registerConstrained(reg *generator.Registry) {
	reg.Register("optim.equality_lagrange", &equalityLagrangeGen{})
	reg.Register("optim.active_constraint", &activeConstraintGen{})
	reg.Register("optim.penalty", &penaltyGen{})
	reg.Register("optim.barrier", &barrierGen{})
	reg.Register("optim.augmented_lagrangian", &augmentedLagrangianGen{})
	reg.Register("optim.projected_gradient", &projectedGradientGen{})
	reg.Register("optim.active_set", &activeSetGen{})
	reg.Register("optim.trust_region_subproblem", &trustRegionSubproblemGen{})
	reg.Register("optim.kkt_matrix", &kktMatrixGen{})
	reg.Register("optim.duality_gap", &dualityGapGen{})
}
