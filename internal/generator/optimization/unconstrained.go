package optimization

import (
	"fmt"

	"github.com/chuma-beep/mathua/internal/generator"
)

// This file covers unconstrained optimization. Every answer is an integer
// because the numeric grading type compares a parsed number, and every family
// draws at least three parameters so a feed can offer twenty distinct
// instances of the same ordinal.

// ---------------------------------------------------------------------------
// The Hessian matrix
// ---------------------------------------------------------------------------

type hessianGen struct{}

func (g *hessianGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		a, b, c := limit(2, 9), limit(1, 6), limit(2, 9)
		return at(kp, fmt.Sprintf("For f(x, y) = %d x^2 + %d x y + %d y^2, add the two diagonal entries of the Hessian. (enter a number)", a, b, c), 2*a+2*c,
			fmt.Sprintf("The diagonal entries are 2·%d = %d and 2·%d = %d, so the sum is %d.", a, 2*a, c, 2*c, 2*a+2*c))
	case 2:
		a, b, c := limit(2, 9), limit(1, 6), limit(2, 9)
		return at(kp, fmt.Sprintf("For f(x, y) = %d x^2 + %d x y + %d y^2, the Hessian is constant. Enter its determinant.", a, b, c), 4*a*c-b*b,
			fmt.Sprintf("The Hessian is [[%d, %d], [%d, %d]], so its determinant is 4·%d·%d - %d² = %d.", 2*a, b, b, 2*c, a, c, b, 4*a*c-b*b))
	default:
		a, b, c := limit(2, 7), limit(1, 5), limit(2, 7)
		return at(kp, fmt.Sprintf("For f(x, y) = %d x^2 + %d x y + %d y^2 the Hessian determinant is 4ac - b². Enter ac - b².", a, b, c), a*c-b*b,
			fmt.Sprintf("ac - b² = %d·%d - %d² = %d. A positive value means the Hessian is positive definite.", a, c, b, a*c-b*b))
	}
}

// ---------------------------------------------------------------------------
// A gradient descent step
// ---------------------------------------------------------------------------

type gradientStepGen struct{}

func (g *gradientStepGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		a, b, x, t := limit(2, 9), limit(1, 12), limit(1, 9), limit(1, 3)
		next := x - t*(a*x+b)
		return at(kp, fmt.Sprintf("For f'(x) = %d x + %d, take one step from x = %d with step size %d: x_new = x - %d f'(x). Enter x_new.", a, b, x, t, t), next,
			fmt.Sprintf("The gradient is %d·%d + %d = %d, so x_new = %d - %d·%d = %d.", a, x, b, a*x+b, x, t, a*x+b, next))
	case 2:
		a, b, x, t := limit(2, 6), limit(1, 8), limit(1, 6), limit(1, 2)
		first := x - t*(a*x+b)
		second := first - t*(a*first+b)
		return at(kp, fmt.Sprintf("Repeat the step x_new = x - %d(%d x + %d) twice, starting at x = %d. Enter the value after two steps.", t, a, b, x), second,
			fmt.Sprintf("After one step x = %d. The next gradient is %d, so x = %d - %d·%d = %d.", first, a*first+b, first, t, a*first+b, second))
	default:
		a, b, x := limit(2, 9), limit(1, 12), limit(1, 9)
		step := x - (a*x + b)
		return at(kp, fmt.Sprintf("With f'(x) = %d x + %d and x = %d, the step x - f'(x) is taken. Enter the signed distance moved, x - x_new.", a, b, x), x-step,
			fmt.Sprintf("The gradient is %d, so the move is %d. x - x_new = %d - %d = %d.", a*x+b, a*x+b, x, step, a*x+b))
	}
}

// ---------------------------------------------------------------------------
// Step size and convergence
// ---------------------------------------------------------------------------

type stepSizeGen struct{}

func (g *stepSizeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		a, t, x := limit(2, 5), limit(1, 4), limit(2, 9)
		next := absInt(1-t*a) * x
		return at(kp, fmt.Sprintf("Gradient descent on f(x) = %d x^2 / 2 with step size %d multiplies the distance to zero by |1 - %d·%d|. Starting at x = %d, enter the new distance to zero.", a, t, t, a, x), next,
			fmt.Sprintf("The contraction factor is |1 - %d| = %d, so the distance becomes %d·%d = %d.", t*a, absInt(1-t*a), absInt(1-t*a), x, next))
	case 2:
		a, x := 2, 1<<limit(5, 24)
		steps := 0
		for v := x; v >= 2; v /= a {
			steps++
		}
		return at(kp, fmt.Sprintf("With step size 1 on f(x) = x^2, the iterate halves each step. Starting at x = %d, how many steps until the distance to zero first drops below 2?", x), steps,
			fmt.Sprintf("After k steps the distance is %d/2^k, which is below 2 when 2^k > %d, so k = %d.", x, x/2, steps))
	default:
		a, t := limit(1, 5), limit(1, 5)
		return at(kp, fmt.Sprintf("For f(x) = %d x^2 / 2 with step size %d, enter the contraction factor |1 - t a| as an integer.", a, t), absInt(1-t*a),
			fmt.Sprintf("The factor is |1 - %d·%d| = %d. A factor below 1 contracts, a factor of exactly 1 stalls and a larger factor diverges.", t, a, absInt(1-t*a)))
	}
}

// ---------------------------------------------------------------------------
// Newton's method for minimization
// ---------------------------------------------------------------------------

type newtonOptimGen struct{}

func (g *newtonOptimGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		a, k, x := limit(2, 9), limit(1, 9), limit(1, 9)
		return at(kp, fmt.Sprintf("For f(x) = %d x^2 / 2 + %d x, one Newton step from x = %d lands on the minimizer. Enter that minimizer.", a, a*k, x), x-k,
			fmt.Sprintf("The Newton step is -f'(x)/f''(x) = -(%d x + %d)/%d, which at x = %d equals -%d.", a, a*k, a, x, k))
	case 2:
		a, g := limit(2, 9), limit(1, 12)
		step := (a * g)
		return at(kp, fmt.Sprintf("A quadratic has second derivative %d and gradient %d at the current point. Enter the magnitude of the Newton step.", a, step), step/a,
			fmt.Sprintf("The Newton step magnitude is |gradient| / f'' = %d / %d = %d.", step, a, step/a))
	default:
		a, x := limit(2, 9), limit(1, 9)
		grad := a * x
		return at(kp, fmt.Sprintf("A quadratic has second derivative %d and gradient %d at x = %d. Compare the Newton step with a gradient step of size 1: enter Newton step magnitude divided by gradient step magnitude.", a, grad, x), a,
			fmt.Sprintf("Newton moves by %d/%d = %d while the gradient step moves by %d, so the ratio is %d.", grad, a, x, grad, a))
	}
}

// ---------------------------------------------------------------------------
// The BFGS update
// ---------------------------------------------------------------------------

type bfgsGen struct{}

func (g *bfgsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		s, y := limit(2, 9), limit(1, 6)
		return at(kp, fmt.Sprintf("In one dimension the BFGS inverse-Hessian update from H = 1 with step s = %d and gradient change y = %d gives H_new = s / y. Enter H_new.", s*y, y), s,
			fmt.Sprintf("H_new = %d / %d = %d. The inverse Hessian is the ratio of the step to the gradient change.", s*y, y, s))
	case 2:
		a, b := limit(2, 9), limit(2, 9)
		return at(kp, fmt.Sprintf("With B = diag(%d, %d) and step s = (1, 1), the secant condition gives y = B s. Enter the sum of the two entries of y.", a, b), a+b,
			fmt.Sprintf("B s = (%d, %d), so the sum is %d.", a, b, a+b))
	default:
		a, c, s1 := limit(2, 9), limit(2, 9), limit(1, 6)
		return at(kp, fmt.Sprintf("With B = diag(%d, %d) and s = (%d, %d), enter the first entry of B s.", a, c, s1, limit(1, 6)), a*s1,
			fmt.Sprintf("The first entry of B s is %d·%d = %d.", a, s1, a*s1))
	}
}

// ---------------------------------------------------------------------------
// Coordinate descent
// ---------------------------------------------------------------------------

type coordinateDescentGen struct{}

func (g *coordinateDescentGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		a, k, y := limit(1, 6), limit(1, 8), limit(1, 8)
		return at(kp, fmt.Sprintf("For f(x, y) = %d x^2 + %d x y + y^2, hold y = %d and minimize over x. Enter the new x.", a, 2*a*k, y), -k*y,
			fmt.Sprintf("The derivative in x is %d x + %d y, so x = -%d·%d / (2·%d) = %d.", 2*a, 2*a*k, 2*a*k, y, a, -k*y))
	case 2:
		a, b, c := limit(2, 9), limit(1, 6), limit(2, 9)
		x, y := limit(1, 9), limit(1, 9)
		px, py := 2*a*x+b*y, 2*c*y+b*x
		answer := py
		if absInt(px) > absInt(py) {
			answer = px
		}
		return at(kp, fmt.Sprintf("For f(x, y) = %d x^2 + %d x y + %d y^2 at (%d, %d), coordinate descent picks the coordinate with the larger absolute partial derivative. Enter that partial derivative.", a, b, c, x, y), answer,
			fmt.Sprintf("The partials are %d and %d, so the chosen value is %d.", px, py, answer))
	default:
		a, b, c := limit(2, 9), limit(1, 6), limit(2, 9)
		x, y := limit(1, 9), limit(1, 9)
		return at(kp, fmt.Sprintf("For f(x, y) = %d x^2 + %d x y + %d y^2 at (%d, %d), enter the partial derivative in x.", a, b, c, x, y), 2*a*x+b*y,
			fmt.Sprintf("The partial in x is 2·%d·%d + %d·%d = %d.", a, x, b, y, 2*a*x+b*y))
	}
}

// ---------------------------------------------------------------------------
// Conjugate directions
// ---------------------------------------------------------------------------

type conjugateDirectionGen struct{}

func (g *conjugateDirectionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		a, c := limit(2, 9), limit(2, 9)
		p1, p2, q1, q2 := limit(1, 6), limit(1, 6), limit(1, 6), limit(1, 6)
		return at(kp, fmt.Sprintf("With A = diag(%d, %d), p = (%d, %d) and q = (%d, %d), enter pᵀAq.", a, c, p1, p2, q1, q2), a*p1*q1+c*p2*q2,
			fmt.Sprintf("pᵀAq = %d·%d·%d + %d·%d·%d = %d. Conjugate directions make this zero.", a, p1, q1, c, p2, q2, a*p1*q1+c*p2*q2))
	case 2:
		a, c, p2 := limit(2, 9), limit(2, 9), limit(1, 6)
		return at(kp, fmt.Sprintf("With A = diag(%d, %d) and p = (1, %d), a direction q = (-%d·%d, %d) is A-conjugate to p. Enter the first entry of q.", a, c, p2, c, p2, a), -c*p2,
			fmt.Sprintf("A-conjugacy needs %d·1·q1 + %d·%d·%d = 0, so q1 = -%d·%d = %d.", a, c, p2, a, c, p2, -c*p2))
	default:
		n, k := limit(3, 9), limit(1, 3)
		return at(kp, fmt.Sprintf("A conjugate-direction method runs in R^%d. After %d exact steps on a quadratic, enter the dimension of the remaining A-orthogonal complement.", n, k), n-k,
			fmt.Sprintf("Each step removes one direction, so %d - %d = %d directions remain.", n, k, n-k))
	}
}

// ---------------------------------------------------------------------------
// The Armijo condition
// ---------------------------------------------------------------------------

type armijoGen struct{}

func (g *armijoGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		a, x := 2*limit(1, 5), limit(2, 9)
		return at(kp, fmt.Sprintf("For f(x) = %d x^2 / 2 with x = %d, gradient descent with step size 1 moves by -f'(x). Enter the decrease in f: f(x) - f(x - f'(x)).", a, x), (a*a/2)*(2-a)*x*x,
			fmt.Sprintf("The move is %d, so the decrease is (%d/2)(2 - %d)·%d² = %d.", a*x, a, a, x, (a*a/2)*(2-a)*x*x))
	case 2:
		a, x := 2*limit(1, 5), limit(1, 8)
		return at(kp, fmt.Sprintf("An Armijo test uses c = 0.5, step size 1 and f(x) = %d x^2 / 2, so the required decrease is c·1·f'(x)². At x = %d, enter that required decrease.", a, x), (a*a/2)*x*x,
			fmt.Sprintf("f'(x) = %d, so the required decrease is 0.5·%d² = %d.", a*x, a*x, (a*a/2)*x*x))
	default:
		a, x, t := 2*limit(1, 4), limit(1, 8), limit(1, 3)
		value := (a / 2) * (x - t*a*x) * (x - t*a*x)
		return at(kp, fmt.Sprintf("For f(x) = %d x^2 / 2, x = %d and step size %d, enter 2·f(x - t f'(x)).", a, x, t), 2*value,
			fmt.Sprintf("x_new = %d - %d·%d = %d, so 2 f(x_new) = %d.", x, t, a*x, x-t*a*x, 2*value))
	}
}

// ---------------------------------------------------------------------------
// The Wolfe conditions
// ---------------------------------------------------------------------------

type wolfeGen struct{}

func (g *wolfeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		a, x, t := limit(2, 6), limit(1, 8), limit(1, 3)
		slope := -a * x * a * x
		after := -a * (x - t*a*x) * a * x
		return at(kp, fmt.Sprintf("For f(x) = %d x^2 / 2 and p = -f'(x), enter the curvature ratio (new slope · p) / (old slope · p), that is %d / %d.", a, after, slope), after/slope,
			fmt.Sprintf("The new slope is -%d·%d·%d, the old slope is -%d·%d, and the ratio is %d.", a, x-t*a*x, a*x, a*x, a*x, after/slope))
	case 2:
		a, x := limit(2, 6), limit(1, 8)
		return at(kp, fmt.Sprintf("A Wolfe test with c2 = 0.5 requires the new slope to be at least -0.5 times the old slope magnitude. For f(x) = %d x^2 at x = %d with step 1, enter the old slope magnitude.", a, x), a*x,
			fmt.Sprintf("f'(x) = %d, so the magnitude is %d.", a*x, a*x))
	default:
		a, x, t := limit(2, 6), limit(1, 8), limit(1, 3)
		ratio := (x - t*a*x) * (x - t*a*x) * 100 / (x * x)
		return at(kp, fmt.Sprintf("For f(x) = x^2 with x = %d and step size %d, enter the squared-distance ratio after the step as a percentage of the original squared distance.", x, t), ratio,
			fmt.Sprintf("The new point is %d, so the ratio is %d²/%d² = %d percent.", x-t*a*x, x-t*a*x, x, ratio))
	}
}

// ---------------------------------------------------------------------------
// The trust-region ratio
// ---------------------------------------------------------------------------

type trustRegionRatioGen struct{}

func (g *trustRegionRatioGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	switch kp {
	case 1:
		a, x, s := limit(1, 5), limit(1, 6), limit(1, 4)
		actual := a * (x*x - (x-s)*(x-s))
		return at(kp, fmt.Sprintf("For f(x) = %d x^2 at x = %d, a step of %d changes the function. Enter the actual reduction f(x) - f(x + s).", a, x, s), actual,
			fmt.Sprintf("f(%d) - f(%d) = %d - %d = %d.", x, x-s, a*x*x, a*(x-s)*(x-s), actual))
	case 2:
		grad, step, curv := limit(2, 9), limit(1, 6), limit(2, 8)
		return at(kp, fmt.Sprintf("A quadratic model has gradient %d and curvature %d. For a step s = %d, enter 2 times the predicted reduction, 2(g s + 0.5 B s²).", grad, curv, step), 2*grad*step+curv*step*step,
			fmt.Sprintf("g s + 0.5 B s² = %d·%d + 0.5·%d·%d² = %d, so twice it is %d.", grad, step, curv, step, grad*step+curv*step*step/2, 2*grad*step+curv*step*step))
	default:
		r, actual, predicted := limit(2, 6), limit(2, 9), limit(2, 9)
		factor := 4
		if 4*actual >= 3*predicted {
			factor = 2
		}
		return at(kp, fmt.Sprintf("A trust region has radius %d. The ratio of actual to predicted reduction decides the next radius: double it when the ratio is at least 0.75, otherwise quarter it. For actual reduction %d and predicted reduction %d, enter the new radius.", r, actual, predicted), r*factor,
			fmt.Sprintf("The ratio is %d/%d, so the radius is multiplied by %d and becomes %d.", actual, predicted, factor, r*factor))
	}
}

func registerUnconstrained(reg *generator.Registry) {
	reg.Register("optim.hessian", &hessianGen{})
	reg.Register("optim.gradient_step", &gradientStepGen{})
	reg.Register("optim.step_size", &stepSizeGen{})
	reg.Register("optim.newton_optim", &newtonOptimGen{})
	reg.Register("optim.bfgs", &bfgsGen{})
	reg.Register("optim.coordinate_descent", &coordinateDescentGen{})
	reg.Register("optim.conjugate_direction", &conjugateDirectionGen{})
	reg.Register("optim.armijo", &armijoGen{})
	reg.Register("optim.wolfe", &wolfeGen{})
	reg.Register("optim.trust_region_ratio", &trustRegionRatioGen{})
}
