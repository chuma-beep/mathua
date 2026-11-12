package optimization

import "fmt"

var wave12FamiliesC = map[string][3]func() trio{
	// -----------------------------------------------------------------------
	// Structured optimization
	// -----------------------------------------------------------------------
	"optim.least_squares": {
		func() trio {
			a, k := limit(2, 9), limit(1, 9)
			return mk(fmt.Sprintf("A scalar least-squares problem minimizes (a x - b)² with a = %d and b = %d. Enter its solution, using integer division.", a, a*k), k,
				fmt.Sprintf("The normal equation gives a x = b, so x = %d/%d = %d.", a*k, a, k))
		},
		func() trio {
			a, x, b := limit(2, 6), limit(1, 8), limit(1, 30)
			return mk(fmt.Sprintf("A fit uses a = %d, x = %d and b = %d. Enter the residual b - a x.", a, x, b), b-a*x,
				fmt.Sprintf("%d - %d·%d = %d.", b, a, x, b-a*x))
		},
		func() trio {
			a, x, b := limit(1, 4), limit(1, 5), limit(1, 9)
			r := b - a*x
			return mk(fmt.Sprintf("A least-squares objective is the squared residual. With a = %d, x = %d and b = %d, enter the residual squared.", a, x, b), r*r,
				fmt.Sprintf("The residual is %d, so its square is %d.", r, r*r))
		},
	},
	"optim.regularized_least_squares": {
		func() trio {
			a, lam, k := limit(2, 6), limit(1, 4), limit(1, 8)
			return mk(fmt.Sprintf("A ridge problem minimizes (a x - b)² + λ x² with a = %d, λ = %d and b = %d. Enter its solution, using integer division.", a, lam, (a+lam)*k), k,
				fmt.Sprintf("The normal equation gives (a + λ) x = b, so x = %d/%d = %d.", (a+lam)*k, a+lam, k))
		},
		func() trio {
			lam, x := limit(1, 6), limit(1, 9)
			return mk(fmt.Sprintf("A ridge penalty is λ x² with λ = %d and x = %d. Enter the penalty.", lam, x), lam*x*x,
				fmt.Sprintf("%d·%d² = %d.", lam, x, lam*x*x))
		},
		func() trio {
			a, x, b, lam := limit(1, 4), limit(1, 6), limit(1, 9), limit(1, 3)
			r := b - a*x
			return mk(fmt.Sprintf("A ridge objective is (a x - b)² + λ x² with a = %d, x = %d, b = %d and λ = %d. Enter its value.", a, x, b, lam), r*r+lam*x*x,
				fmt.Sprintf("The squared residual is %d and the penalty is %d, so the total is %d.", r*r, lam*x*x, r*r+lam*x*x))
		},
	},
	"optim.lasso": {
		func() trio {
			z, lam := limit(-9, 9), limit(1, 4)
			answer := 0
			if z > lam {
				answer = z - lam
			} else if z < -lam {
				answer = z + lam
			}
			return mk(fmt.Sprintf("A lasso soft threshold shrinks z = %d toward zero by λ = %d. Enter the thresholded value.", z, lam), answer,
				fmt.Sprintf("The result is sign(z)·max(|z| - λ, 0), which here is %d.", answer))
		},
		func() trio {
			lam, x := limit(1, 5), limit(1, 9)
			return mk(fmt.Sprintf("A lasso penalty is λ|x| with λ = %d and x = %d. Enter the penalty.", lam, x), lam*x,
				fmt.Sprintf("%d·%d = %d.", lam, x, lam*x))
		},
		func() trio {
			z, x, lam := limit(-8, 8), limit(-8, 8), limit(1, 3)
			answer := (x-z)*(x-z) + lam*absInt(x)
			return mk(fmt.Sprintf("A lasso objective is (x - z)² + λ|x| with z = %d, x = %d and λ = %d. Enter its value.", z, x, lam), answer,
				fmt.Sprintf("The squared gap is %d and the penalty is %d, so the total is %d.", (x-z)*(x-z), lam*absInt(x), answer))
		},
	},
	"optim.proximal_operator": {
		func() trio {
			v, lam := limit(1, 9), limit(1, 3)
			return mk(fmt.Sprintf("The proximal operator of f(x) = x²/2 is v/(1 + λ). With v = %d and λ = %d, enter the result, using integer division.", v, lam), v/(1+lam),
				fmt.Sprintf("%d/%d = %d.", v, 1+lam, v/(1+lam)))
		},
		func() trio {
			v, lam := limit(-9, 9), limit(1, 4)
			answer := 0
			if v > lam {
				answer = v - lam
			} else if v < -lam {
				answer = v + lam
			}
			return mk(fmt.Sprintf("The proximal operator of λ|x| is the soft threshold. With v = %d and λ = %d, enter the result.", v, lam), answer,
				fmt.Sprintf("Soft thresholding gives %d.", answer))
		},
		func() trio {
			v, lam := limit(2, 12), limit(1, 4)
			prox := v / (1 + lam)
			return mk(fmt.Sprintf("A proximal step moves v = %d to %d. Enter the distance moved.", v, prox), v-prox,
				fmt.Sprintf("%d - %d = %d.", v, prox, v-prox))
		},
	},
	"optim.proximal_gradient": {
		func() trio {
			x, t, g := limit(1, 9), limit(1, 2), limit(1, 6)
			grad := x - t*g
			return mk(fmt.Sprintf("A proximal gradient step first takes x - t g with x = %d, t = %d and g = %d. Enter the intermediate value.", x, t, g), grad,
				fmt.Sprintf("%d - %d·%d = %d.", x, t, g, grad))
		},
		func() trio {
			x, t, g, lam := limit(2, 12), limit(1, 2), limit(1, 5), limit(1, 2)
			grad := x - t*g
			answer := 0
			if grad > lam {
				answer = grad - lam
			} else if grad < -lam {
				answer = grad + lam
			}
			return mk(fmt.Sprintf("A proximal gradient step takes x - t g then soft thresholds by λ. With x = %d, t = %d, g = %d and λ = %d, enter the result.", x, t, g, lam), answer,
				fmt.Sprintf("The gradient step gives %d, and soft thresholding gives %d.", grad, answer))
		},
		func() trio {
			x, t, g := limit(2, 10), limit(1, 2), limit(1, 6)
			return mk(fmt.Sprintf("A gradient step moves x = %d to x - t g with t = %d and g = %d. Enter the distance moved.", x, t, g), t*g,
				fmt.Sprintf("The move is %d·%d = %d.", t, g, t*g))
		},
	},
	"optim.soft_threshold": {
		func() trio {
			z, lam := limit(-9, 9), limit(1, 5)
			answer := 0
			if z > lam {
				answer = z - lam
			} else if z < -lam {
				answer = z + lam
			}
			return mk(fmt.Sprintf("Apply the soft threshold to z = %d with λ = %d. Enter the result.", z, lam), answer,
				fmt.Sprintf("The result is sign(z)·max(|z| - λ, 0), which here is %d.", answer))
		},
		func() trio {
			z, lam := limit(1, 9), limit(1, 6)
			return mk(fmt.Sprintf("A soft threshold shrinks z = %d by λ = %d. Enter the amount removed.", z, lam), lam,
				fmt.Sprintf("The value falls by λ = %d while it stays positive.", lam))
		},
		func() trio {
			z, lam := limit(2, 40), limit(1, 3)
			return mk(fmt.Sprintf("A soft threshold is applied to both components of (%d, %d) with λ = %d. Enter the sum of the results.", z, z, lam), 2*(z-lam),
				fmt.Sprintf("Each component becomes %d, so the sum is %d.", z-lam, 2*(z-lam)))
		},
	},
	"optim.orthogonal_projection": {
		func() trio {
			x, low, high := limit(-6, 6), limit(0, 3), limit(7, 12)
			return mk(fmt.Sprintf("Project x = %d onto [%d, %d]. Enter the projected value.", x, low, high), minInt(maxInt(x, low), high),
				fmt.Sprintf("Clamping to the interval gives %d.", minInt(maxInt(x, low), high)))
		},
		func() trio {
			x, high := limit(15, 25), limit(4, 10)
			return mk(fmt.Sprintf("Project x = %d onto (-infinity, %d]. Enter the distance moved.", x, high), x-high,
				fmt.Sprintf("%d - %d = %d.", x, high, x-high))
		},
		func() trio {
			a, b, x, y := limit(1, 4), limit(1, 4), limit(1, 6), limit(1, 6)
			limitValue := a*x + b*y - limit(1, 9)
			return mk(fmt.Sprintf("A halfspace is a x + b y ≤ %d. At (%d, %d) with a = %d and b = %d, enter the violation of the boundary a x + b y - %d.", limitValue, x, y, a, b, limitValue), a*x+b*y-limitValue,
				fmt.Sprintf("%d + %d - %d = %d; a positive value must be projected away.", a*x, b*y, limitValue, a*x+b*y-limitValue))
		},
	},
	"optim.quadratic_program": {
		func() trio {
			a, k := limit(2, 9), limit(1, 9)
			return mk(fmt.Sprintf("Minimize (a/2)x² + b x with a = %d and b = %d. Enter the unconstrained minimizer, using integer division.", a, a*k), -k,
				fmt.Sprintf("The gradient a x + b vanishes at x = -%d/%d = %d.", a*k, a, -k))
		},
		func() trio {
			a, x, b := limit(2, 9), limit(1, 9), limit(1, 20)
			return mk(fmt.Sprintf("A quadratic program has stationarity residual a x + b with a = %d, x = %d and b = %d. Enter it.", a, x, b), a*x+b,
				fmt.Sprintf("%d·%d + %d = %d.", a, x, b, a*x+b))
		},
		func() trio {
			a, x, b := limit(2, 6), limit(1, 8), limit(1, 12)
			value := (a*x*x)/2 + b*x
			return mk(fmt.Sprintf("A quadratic objective is (a/2)x² + b x with a = %d, x = %d and b = %d. Enter its value, using integer division.", a, x, b), value,
				fmt.Sprintf("(%d·%d²)/2 + %d·%d = %d.", a, x, b, x, value))
		},
	},
	"optim.trust_region_global": {
		func() trio {
			a, k, delta := limit(2, 9), limit(-6, 6), limit(1, 5)
			answer := minInt(maxInt(-k, -delta), delta)
			return mk(fmt.Sprintf("Minimize (%d/2)x² + %d x over |x| ≤ %d. Enter x.", a, a*k, delta), answer,
				fmt.Sprintf("The unconstrained minimizer is %d, and the ball keeps %d.", -k, answer))
		},
		func() trio {
			actual, predicted := limit(1, 9), limit(1, 9)
			ratio := actual * 4 / predicted
			return mk(fmt.Sprintf("The trust-region ratio is actual/predicted with actual = %d and predicted = %d. Enter it scaled by 4, using integer division.", actual, predicted), ratio,
				fmt.Sprintf("4·%d/%d = %d.", actual, predicted, ratio))
		},
		func() trio {
			radius, factor := limit(2, 8), limit(2, 4)
			return mk(fmt.Sprintf("A rejected step shrinks the radius from %d by the factor %d. Enter the new radius.", radius, factor), radius/factor,
				fmt.Sprintf("%d / %d = %d.", radius, factor, radius/factor))
		},
	},
	"optim.nonsmooth_optimality": {
		func() trio {
			c := limit(2, 40)
			return mk(fmt.Sprintf("For f(x) = %d|x| the subdifferential at x = 0 is [-%d, %d]. Enter the number of integers in that interval.", c, c, c), 2*c+1,
				fmt.Sprintf("The integers from -%d to %d number %d.", c, c, 2*c+1))
		},
		func() trio {
			c := limit(2, 40)
			return mk(fmt.Sprintf("For f(x) = %d|x - a| the subdifferential at the kink is [-%d, %d]. Enter the width of that interval.", c, c, c), 2*c,
				fmt.Sprintf("%d - (-%d) = %d.", c, c, 2*c))
		},
		func() trio {
			a, c := limit(1, 6), limit(1, 6)
			upper := c + 1
			return mk(fmt.Sprintf("For f(x) = |x - %d| + %d x the subdifferential at the kink is [%d, %d]. Enter its largest member.", a, c, c-1, c+1), upper,
				fmt.Sprintf("The largest member is %d.", upper))
		},
	},
}
