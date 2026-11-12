package optimization

import (
	"fmt"
)

// Wave 13, nonsmooth optimization. Three integer families per concept, in the
// same adapter shape as the other Wave 13 files.
var wave13FamiliesC = map[string][3]func() trio{
	"optim.proximal_point": {
		func() trio {
			value, penalty := limit(3, 30), limit(1, 9)
			return mk(fmt.Sprintf("The proximal operator of the penalty %d at the value %d, with value above penalty, returns v - λ. Enter it.", penalty, value), value-penalty,
				fmt.Sprintf("%d - %d = %d.", value, penalty, value-penalty))
		},
		func() trio {
			value, shift := 2*limit(1, 15), 2*limit(1, 15)
			return mk(fmt.Sprintf("With step one the proximal point is (v + a)/2 with v = %d and a = %d. Enter it.", value, shift), (value+shift)/2,
				fmt.Sprintf("(%d + %d)/2 = %d.", value, shift, (value+shift)/2))
		},
		func() trio {
			steps, cost := limit(2, 20), limit(1, 9)
			return mk(fmt.Sprintf("A proximal method takes %d steps at %d units of work each. Enter the work.", steps, cost), steps*cost,
				fmt.Sprintf("%d times %d is %d.", steps, cost, steps*cost))
		},
	},
	"optim.moreau_envelope": {
		func() trio {
			distance := 2 * limit(1, 30)
			return mk(fmt.Sprintf("The Moreau envelope with parameter one adds half the squared distance %d². Enter the added amount.", distance), distance*distance/2,
				fmt.Sprintf("%d/2 = %d.", distance*distance, distance*distance/2))
		},
		func() trio {
			first, second := limit(2, 40), limit(2, 40)
			return mk(fmt.Sprintf("A one-Lipschitz envelope takes the values %d and %d. Enter the largest possible difference.", first, second), absInt(first-second),
				fmt.Sprintf("|%d - %d| = %d.", first, second, absInt(first-second)))
		},
		func() trio {
			distance := limit(1, 20)
			return mk(fmt.Sprintf("The envelope with parameter two adds half of two times the squared distance %d². Enter the added amount.", distance), distance*distance,
				fmt.Sprintf("2·%d/2 = %d.", distance*distance, distance*distance))
		},
	},
	"optim.dual_smoothing": {
		func() trio {
			parameter, errorSize := limit(1, 20), limit(1, 12)
			return mk(fmt.Sprintf("Smoothing adds a term μ·ε with μ = %d and ε = %d. Enter the added term.", parameter, errorSize), parameter*errorSize,
				fmt.Sprintf("%d times %d is %d.", parameter, errorSize, parameter*errorSize))
		},
		func() trio {
			steps, parameter, left := limit(1, 15), limit(1, 9), limit(1, 10)
			start := steps*parameter + left
			return mk(fmt.Sprintf("A smoothed gap starts at %d and falls by a smoothing of %d each step. After %d steps, enter the gap.", start, parameter, steps), left,
				fmt.Sprintf("%d - %d·%d = %d.", start, steps, parameter, left))
		},
		func() trio {
			first, second := 2*limit(1, 20), 2*limit(1, 20)
			return mk(fmt.Sprintf("Two smoothed values %d and %d are averaged. Enter the average.", first, second), (first+second)/2,
				fmt.Sprintf("(%d + %d)/2 = %d.", first, second, (first+second)/2))
		},
	},
	"optim.smoothing_parameter": {
		func() trio {
			base, halvings := limit(1, 12), limit(1, 8)
			return mk(fmt.Sprintf("A smoothing parameter of %d is halved %d times. Enter the final parameter.", base<<halvings, halvings), base,
				fmt.Sprintf("%d divided by %d is %d.", base<<halvings, 1<<halvings, base))
		},
		func() trio {
			ratio, parameter := limit(1, 20), limit(1, 9)
			return mk(fmt.Sprintf("A smoothing error is ε/μ with ε = %d and μ = %d. Enter the error.", ratio*parameter, parameter), ratio,
				fmt.Sprintf("%d divided by %d is %d.", ratio*parameter, parameter, ratio))
		},
		func() trio {
			ratio, steps := limit(1, 20), limit(2, 12)
			return mk(fmt.Sprintf("A schedule sets μ_t = μ0/t with μ0 = %d and t = %d. Enter μ_t.", ratio*steps, steps), ratio,
				fmt.Sprintf("%d divided by %d is %d.", ratio*steps, steps, ratio))
		},
	},
	"optim.huber_loss": {
		func() trio {
			residual := 2 * limit(1, 25)
			return mk(fmt.Sprintf("Inside the Huber band the loss is r²/2. For residual %d, enter it.", residual), residual*residual/2,
				fmt.Sprintf("%d/2 = %d.", residual*residual, residual*residual/2))
		},
		func() trio {
			delta, excess := 2*limit(1, 8), limit(1, 9)
			return mk(fmt.Sprintf("Outside the band the Huber loss is δ·k + δ²/2 with δ = %d and excess %d. Enter it.", delta, excess), delta*excess+delta*delta/2,
				fmt.Sprintf("%d·%d + %d/2 = %d.", delta, excess, delta*delta, delta*excess+delta*delta/2))
		},
		func() trio {
			first, second := 2*limit(1, 12), 2*limit(1, 12)
			return mk(fmt.Sprintf("Two residuals %d and %d both lie inside the band. Enter half the sum of their squares.", first, second), (first*first+second*second)/2,
				fmt.Sprintf("(%d + %d)/2 = %d.", first*first, second*second, (first*first+second*second)/2))
		},
	},
	"optim.hinge_loss": {
		func() trio {
			margin := -limit(1, 25)
			return mk(fmt.Sprintf("A hinge loss is max(0, 1 - ywx) with margin ywx = %d. Enter the loss.", margin), 1-margin,
				fmt.Sprintf("1 - (%d) = %d.", margin, 1-margin))
		},
		func() trio {
			terms, each := limit(2, 20), limit(1, 9)
			return mk(fmt.Sprintf("%d hinge terms each contribute %d. Enter their sum.", terms, each), terms*each,
				fmt.Sprintf("%d times %d is %d.", terms, each, terms*each))
		},
		func() trio {
			threshold, margin := limit(2, 20), limit(-20, -1)
			return mk(fmt.Sprintf("A hinge loss is max(0, threshold - margin) with threshold %d and margin %d. Enter it.", threshold, margin), threshold-margin,
				fmt.Sprintf("%d - (%d) = %d.", threshold, margin, threshold-margin))
		},
	},
	"optim.subgradient_method": {
		func() trio {
			ratio, step := limit(1, 20), limit(1, 9)
			return mk(fmt.Sprintf("A subgradient step is η/g with η = %d and g = %d. Enter it.", step, ratio*step), ratio,
				fmt.Sprintf("%d divided by %d is %d.", step, ratio*step, ratio))
		},
		func() trio {
			start, steps, size := limit(2, 60), limit(1, 9), limit(1, 6)
			return mk(fmt.Sprintf("A constant-step subgradient method starts at %d and subtracts %d for %d steps. Enter the final iterate.", start, size, steps), start-steps*size,
				fmt.Sprintf("%d - %d·%d = %d.", start, steps, size, start-steps*size))
		},
		func() trio {
			ratio, norm := limit(1, 20), limit(1, 12)
			return mk(fmt.Sprintf("A Polyak step is (f - f*)/||g||² with f - f* = %d and squared norm %d. Enter it.", ratio*norm, norm), ratio,
				fmt.Sprintf("%d divided by %d is %d.", ratio*norm, norm, ratio))
		},
	},
	"optim.polyak_step": {
		func() trio {
			ratio, norm := limit(1, 25), limit(1, 10)
			return mk(fmt.Sprintf("A Polyak step uses gap %d and squared norm %d. Enter the step.", ratio*norm, norm), ratio,
				fmt.Sprintf("%d divided by %d is %d.", ratio*norm, norm, ratio))
		},
		func() trio {
			target, norm := limit(1, 30), limit(1, 10)
			return mk(fmt.Sprintf("A target decrease of %d with squared norm %d gives the step. Enter it.", target*norm, norm), target,
				fmt.Sprintf("%d divided by %d is %d.", target*norm, norm, target))
		},
		func() trio {
			current, best := limit(20, 80), limit(1, 18)
			return mk(fmt.Sprintf("A current value is %d and the best known value is %d. Enter the gap.", current, best), current-best,
				fmt.Sprintf("%d - %d = %d.", current, best, current-best))
		},
	},
	"optim.bundle_method": {
		func() trio {
			cuts, length := limit(2, 15), limit(2, 12)
			return mk(fmt.Sprintf("A bundle stores %d subgradients of length %d. Enter the stored entries.", cuts, length), cuts*length,
				fmt.Sprintf("%d times %d is %d.", cuts, length, cuts*length))
		},
		func() trio {
			cuts := limit(2, 30)
			return mk(fmt.Sprintf("A bundle master solves a quadratic program with one variable per cut, and there are %d cuts. Enter the variables.", cuts), cuts,
				fmt.Sprintf("One variable per cut gives %d.", cuts))
		},
		func() trio {
			added, iterations := limit(1, 9), limit(2, 20)
			return mk(fmt.Sprintf("%d cuts are added per iteration over %d iterations. Enter the added cuts.", added, iterations), added*iterations,
				fmt.Sprintf("%d times %d is %d.", added, iterations, added*iterations))
		},
	},
	"optim.nonsmooth_kkt": {
		func() trio {
			first, second := limit(-20, 20), limit(1, 20)
			return mk(fmt.Sprintf("A nonsmooth stationarity condition adds subgradients %d and %d. Enter the sum.", first, second), first+second,
				fmt.Sprintf("%d + %d = %d.", first, second, first+second))
		},
		func() trio {
			gradient, multiplier := limit(2, 30), limit(1, 12)
			return mk(fmt.Sprintf("A bound multiplier %d is added to the gradient %d. Enter the KKT residual.", multiplier, gradient), gradient+multiplier,
				fmt.Sprintf("%d + %d = %d.", gradient, multiplier, gradient+multiplier))
		},
		func() trio {
			first, second := limit(1, 20), limit(1, 20)
			return mk(fmt.Sprintf("Two active constraints carry multipliers %d and %d. Enter their sum.", first, second), first+second,
				fmt.Sprintf("%d + %d = %d.", first, second, first+second))
		},
	},
}
