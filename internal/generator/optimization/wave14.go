package optimization

import (
	"fmt"

	"github.com/chuma-beep/mathua/internal/generator"
)

// Wave 14 extends the optimization domain to two hundred concepts. The
// families reuse the Wave 12 adapter: three integer-valued closures per
// concept, selected by the ordinal, so the numeric grader can check every
// answer.
var wave14Families = map[string][3]func() trio{
	// -----------------------------------------------------------------------
	// Saddle-point methods
	// -----------------------------------------------------------------------
	"optim.minimax": {
		func() trio {
			slope, bound := limit(2, 12), limit(1, 9)
			return mk(fmt.Sprintf("A minimax inner problem maximises %d x over 0 ≤ x ≤ %d. Enter the maximum.", slope, bound), slope*bound,
				fmt.Sprintf("A linear function is largest at the upper bound: %d·%d = %d.", slope, bound, slope*bound))
		},
		func() trio {
			x, y := limit(1, 9), limit(1, 9)
			return mk(fmt.Sprintf("For f(x, y) = x^2 - y^2, enter f at the point (%d, %d).", x, y), x*x-y*y,
				fmt.Sprintf("%d² - %d² = %d, which is the saddle value at that point.", x, y, x*x-y*y))
		},
		func() trio {
			rows, cols := limit(2, 9), limit(2, 9)
			return mk(fmt.Sprintf("A matrix game has %d row strategies and %d column strategies. Enter the number of pure strategy pairs.", rows, cols), rows*cols,
				fmt.Sprintf("%d rows times %d columns is %d pure pairs.", rows, cols, rows*cols))
		},
	},
	"optim.saddle_point": {
		func() trio {
			a, b, x, y := limit(1, 6), limit(1, 6), limit(1, 5), limit(1, 5)
			return mk(fmt.Sprintf("For f(x, y) = %d x^2 - %d y^2, enter f at (%d, %d).", a, b, x, y), a*x*x-b*y*y,
				fmt.Sprintf("%d·%d² - %d·%d² = %d.", a, x, b, y, a*x*x-b*y*y))
		},
		func() trio {
			positive, negative := limit(1, 9), limit(1, 9)
			return mk(fmt.Sprintf("A saddle Hessian has %d positive and %d negative eigenvalues. Enter the count with positive curvature.", positive, negative), positive,
				fmt.Sprintf("The positive eigenvalues are exactly %d of the %d total.", positive, positive+negative))
		},
		func() trio {
			ux, uy := limit(2, 12), limit(2, 12)
			return mk(fmt.Sprintf("At a saddle the horizontal step changes the value by %d and the vertical step by -%d. Enter the sum of the two changes.", ux, uy), ux-uy,
				fmt.Sprintf("%d + (-%d) = %d, and neither single change improves both blocks.", ux, uy, ux-uy))
		},
	},
	"optim.zero_sum_game": {
		func() trio {
			entry, pnum, pden := limit(1, 9), limit(1, 4), limit(2, 5)
			return mk(fmt.Sprintf("A payoff matrix has every entry equal to %d. A mixed strategy weights one column by %d/%d. Enter the expected payoff numerator over the denominator %d.", entry, pnum, pden, pden), entry*pnum,
				fmt.Sprintf("Every column pays %d, so the weighted payoff is %d·%d/%d, with numerator %d.", entry, entry, pnum, pden, entry*pnum))
		},
		func() trio {
			floorValue, cols := limit(2, 12), limit(2, 4)
			return mk(fmt.Sprintf("The row player guarantees %d against each of %d pure columns. Enter the total guarantee across the columns.", floorValue, cols), floorValue*cols,
				fmt.Sprintf("%d per column over %d columns is %d.", floorValue, cols, floorValue*cols))
		},
		func() trio {
			payoff, count := limit(1, 6), limit(2, 5)
			return mk(fmt.Sprintf("A uniform mixture over %d columns each paying %d has expected payoff p/%d. Enter the numerator p.", count, payoff, count), payoff*count,
				fmt.Sprintf("The average is %d·%d/%d, so the numerator is %d.", payoff, count, count, payoff*count))
		},
	},
	"optim.robust_minimax": {
		func() trio {
			nominal, radius, count := limit(5, 20), limit(1, 4), limit(2, 5)
			return mk(fmt.Sprintf("A nominal loss is %d and the uncertainty radius is %d in each of %d independent directions. Enter the worst-case loss.", nominal, radius, count), nominal+radius*count,
				fmt.Sprintf("%d + %d·%d = %d at the worst corner.", nominal, radius, count, nominal+radius*count))
		},
		func() trio {
			low, high := limit(1, 10), limit(2, 8)
			return mk(fmt.Sprintf("A box uncertainty set runs from %d to %d on one coordinate. Enter its width.", low, low+high), high,
				fmt.Sprintf("%d - %d = %d.", low+high, low, high))
		},
		func() trio {
			average, worst, weight := limit(2, 12), limit(2, 12), limit(1, 4)
			return mk(fmt.Sprintf("A robust objective mixes the average loss %d with the worst-case loss %d at weight 1/%d on the worst case. Enter the scaled worst-case term.", average, worst, weight), worst,
				fmt.Sprintf("The worst-case term enters as %d/%d, so its numerator is %d.", worst, weight, worst))
		},
	},
	"optim.primal_dual_saddle": {
		func() trio {
			cost, multiplier, slack := limit(2, 12), limit(1, 5), limit(0, 6)
			return mk(fmt.Sprintf("A Lagrangian is c·x + λ(b - a·x) with c·x = %d, λ = %d and residual b - a·x = %d. Enter its value.", cost, multiplier, slack), cost+multiplier*slack,
				fmt.Sprintf("%d + %d·%d = %d.", cost, multiplier, slack, cost+multiplier*slack))
		},
		func() trio {
			primal, dualValue := limit(5, 20), limit(1, 5)
			return mk(fmt.Sprintf("A convex problem has primal value %d and dual value %d below it. Enter the duality gap.", primal, dualValue), dualValue,
				fmt.Sprintf("The gap is %d - %d = %d.", primal, primal-dualValue, dualValue))
		},
		func() trio {
			constraints, variables := limit(1, 5), limit(1, 6)
			return mk(fmt.Sprintf("A primal-dual method updates %d primal variables and %d multipliers per step. Enter the total coordinates updated.", variables, constraints), variables+constraints,
				fmt.Sprintf("%d + %d = %d coordinates.", variables, constraints, variables+constraints))
		},
	},
	"optim.extragradient": {
		func() trio {
			steps, evaluations := limit(2, 40), 2
			return mk(fmt.Sprintf("An extragradient run takes %d steps and each step evaluates the gradient %d times. Enter the evaluations.", steps, evaluations), steps*evaluations,
				fmt.Sprintf("%d steps times %d evaluations is %d.", steps, evaluations, steps*evaluations))
		},
		func() trio {
			gradient, overshoot := limit(2, 10), limit(1, 4)
			return mk(fmt.Sprintf("A trial point moves %d units along the gradient and the correction removes %d units of overshoot. Enter the corrected move.", gradient, overshoot), gradient-overshoot,
				fmt.Sprintf("%d - %d = %d.", gradient, overshoot, gradient-overshoot))
		},
		func() trio {
			lipschitz, divisor := limit(2, 40), limit(2, 3)
			return mk(fmt.Sprintf("The step bound is 1/L with L = %d. Enter the numerator of a step written as %d/L with the same denominator.", lipschitz, divisor-1), divisor-1,
				fmt.Sprintf("A step of %d/%d stays below 1/L, so the numerator is %d.", divisor-1, lipschitz, divisor-1))
		},
	},
	"optim.optimistic_gradient": {
		func() trio {
			current, previous := limit(2, 12), limit(1, 6)
			return mk(fmt.Sprintf("An optimistic direction is 2g - g_prev with g = %d and g_prev = %d. Enter its value.", current, previous), 2*current-previous,
				fmt.Sprintf("2·%d - %d = %d.", current, previous, 2*current-previous))
		},
		func() trio {
			steps, memory := limit(3, 40), 1
			return mk(fmt.Sprintf("Each of %d optimistic steps keeps %d previous gradient. Enter the previous gradients kept in total.", steps, memory), steps*memory,
				fmt.Sprintf("%d steps times %d vector is %d vectors.", steps, memory, steps*memory))
		},
		func() trio {
			g1, g2 := limit(1, 9), limit(1, 9)
			return mk(fmt.Sprintf("Two consecutive gradients are %d and then %d. Enter the extrapolated value 2g2 - g1.", g1, g2), 2*g2-g1,
				fmt.Sprintf("2·%d - %d = %d.", g2, g1, 2*g2-g1))
		},
	},
	"optim.mirror_descent": {
		func() trio {
			dimension, value := limit(2, 6), limit(1, 9)
			return mk(fmt.Sprintf("A uniform distribution over %d coordinates each holding weight %d has total weight p/%d. Enter the numerator p.", dimension, value, dimension), value*dimension,
				fmt.Sprintf("%d·%d/%d has numerator %d.", value, dimension, dimension, value*dimension))
		},
		func() trio {
			gradient, rate, base := limit(2, 9), limit(1, 3), limit(2, 6)
			return mk(fmt.Sprintf("A dual step adds %d·log(%d) with gradient %d. Enter the dual increment using the integer logarithm value %d.", gradient, base, rate, rate), gradient*rate,
				fmt.Sprintf("%d·%d = %d.", gradient, rate, gradient*rate))
		},
		func() trio {
			entropy, count := limit(2, 9), limit(2, 8)
			return mk(fmt.Sprintf("An entropy term sums %d terms each equal to -%d log(1/%d). Enter the sum of the constant factors before the logarithms.", entropy, count, count), entropy*count,
				fmt.Sprintf("%d terms of %d give %d.", count, entropy, entropy*count))
		},
	},
	"optim.dual_averaging": {
		func() trio {
			gradients, each := limit(2, 9), limit(1, 9)
			return mk(fmt.Sprintf("Dual averaging accumulates %d gradients of value %d. Enter the accumulated sum.", gradients, each), gradients*each,
				fmt.Sprintf("%d·%d = %d.", gradients, each, gradients*each))
		},
		func() trio {
			sum, count, extra := limit(2, 12), limit(2, 6), limit(0, 5)
			total := sum*count + extra
			return mk(fmt.Sprintf("An accumulated dual sum of %d is averaged over %d steps. Enter the integer part of the average.", total, count), sum,
				fmt.Sprintf("%d divided by %d is %d with remainder %d.", total, count, sum, extra))
		},
		func() trio {
			penalty, scale := limit(1, 6), limit(2, 5)
			return mk(fmt.Sprintf("A dual averaging penalty is %d·|x| with |x| = %d. Enter its value.", penalty, scale), penalty*scale,
				fmt.Sprintf("%d·%d = %d.", penalty, scale, penalty*scale))
		},
	},
	"optim.saddle_gap": {
		func() trio {
			upper, lower := limit(8, 25), limit(1, 6)
			return mk(fmt.Sprintf("A candidate pair has best-response values %d and %d. Enter the saddle gap.", upper, upper-lower), lower,
				fmt.Sprintf("%d - %d = %d.", upper, upper-lower, lower))
		},
		func() trio {
			first, second := limit(2, 9), limit(2, 9)
			return mk(fmt.Sprintf("Two single-block solves give primal value %d and dual value %d. Enter the gap.", first, second), first-second,
				fmt.Sprintf("%d - %d = %d.", first, second, first-second))
		},
		func() trio {
			value, tolerance := limit(5, 30), limit(1, 4)
			return mk(fmt.Sprintf("A saddle gap of %d is compared with a tolerance of %d. Enter the slack before the tolerance is met.", value, tolerance), value-tolerance,
				fmt.Sprintf("%d - %d = %d.", value, tolerance, value-tolerance))
		},
	},
	// -----------------------------------------------------------------------
	// Optimal control
	// -----------------------------------------------------------------------
	"optim.optimal_control": {
		func() trio {
			rate, horizon, terminal := limit(2, 9), limit(2, 9), limit(1, 6)
			return mk(fmt.Sprintf("A running cost of %d per unit time accumulates over %d units, plus a terminal cost of %d. Enter the total cost.", rate, horizon, terminal), rate*horizon+terminal,
				fmt.Sprintf("%d·%d + %d = %d.", rate, horizon, terminal, rate*horizon+terminal))
		},
		func() trio {
			steps, controls := limit(2, 9), limit(1, 4)
			return mk(fmt.Sprintf("A horizon is split into %d steps and each step carries %d control parameters. Enter the number of decision variables.", steps, controls), steps*controls,
				fmt.Sprintf("%d·%d = %d parameters.", steps, controls, steps*controls))
		},
		func() trio {
			state, control, coupling := limit(2, 6), limit(2, 6), limit(1, 5)
			return mk(fmt.Sprintf("A state costs %d and a control costs %d per step, with %d coupling units. Enter the total cost of one step.", state, control, coupling), state+control+coupling,
				fmt.Sprintf("%d + %d + %d = %d.", state, control, coupling, state+control+coupling))
		},
	},
	"optim.control_hamiltonian": {
		func() trio {
			run, costate, velocity := limit(2, 9), limit(1, 6), limit(1, 6)
			return mk(fmt.Sprintf("A Hamiltonian is running cost %d plus costate %d times velocity %d. Enter its value.", run, costate, velocity), run+costate*velocity,
				fmt.Sprintf("%d + %d·%d = %d.", run, costate, velocity, run+costate*velocity))
		},
		func() trio {
			cost, multiplier, control := limit(2, 10), limit(1, 5), limit(1, 6)
			return mk(fmt.Sprintf("Minimise %d + %d u over 0 ≤ u ≤ %d. Enter the optimal control.", cost, multiplier, control), control,
				fmt.Sprintf("The control coefficient %d is positive, so the minimum is at u = %d.", multiplier, control))
		},
		func() trio {
			costate, stateDerivative := limit(1, 8), limit(1, 8)
			return mk(fmt.Sprintf("A costate equation subtracts %d times the state derivative %d. Enter the magnitude of the change.", costate, stateDerivative), costate*stateDerivative,
				fmt.Sprintf("%d·%d = %d.", costate, stateDerivative, costate*stateDerivative))
		},
	},
	"optim.hjb": {
		func() trio {
			gradient, timeStep := limit(2, 8), limit(1, 5)
			return mk(fmt.Sprintf("A value function changes by %d per unit state over a step of %d. Enter the decrease in value.", gradient, timeStep), gradient*timeStep,
				fmt.Sprintf("%d·%d = %d.", gradient, timeStep, gradient*timeStep))
		},
		func() trio {
			states, steps := limit(2, 6), limit(2, 6)
			return mk(fmt.Sprintf("A Bellman grid has %d states and %d time steps. Enter the number of grid values.", states, steps), states*steps,
				fmt.Sprintf("%d·%d = %d values.", states, steps, states*steps))
		},
		func() trio {
			terminal, cost := limit(1, 9), limit(1, 9)
			return mk(fmt.Sprintf("A terminal value of %d is compared with a control cost of %d. Enter the optimal remaining cost.", terminal, cost), terminal+cost,
				fmt.Sprintf("%d + %d = %d.", terminal, cost, terminal+cost))
		},
	},
	"optim.lqr": {
		func() trio {
			stateWeight, state, controlWeight, control := limit(1, 5), limit(1, 5), limit(1, 5), limit(1, 5)
			return mk(fmt.Sprintf("A regulator cost is %d x^2 + %d u^2 with x = %d and u = %d. Enter the cost.", stateWeight, controlWeight, state, control), stateWeight*state*state+controlWeight*control*control,
				fmt.Sprintf("%d·%d² + %d·%d² = %d.", stateWeight, state, controlWeight, control, stateWeight*state*state+controlWeight*control*control))
		},
		func() trio {
			gain, state := limit(1, 6), limit(1, 9)
			return mk(fmt.Sprintf("A feedback law sets u = -%d x with x = %d. Enter the signed control.", gain, state), -gain*state,
				fmt.Sprintf("u = -%d·%d = %d.", gain, state, -gain*state))
		},
		func() trio {
			weight, pairs := limit(2, 8), limit(1, 4)
			return mk(fmt.Sprintf("A regulator weights %d state pairs at %d each. Enter the total state penalty.", pairs, weight), pairs*weight,
				fmt.Sprintf("%d·%d = %d.", pairs, weight, pairs*weight))
		},
	},
	"optim.riccati": {
		func() trio {
			previous := limit(4, 40)
			return mk(fmt.Sprintf("A Riccati step divides P^2 by P + 1 with P = %d. Enter the integer part of the quotient.", previous), (previous*previous)/(previous+1),
				fmt.Sprintf("%d divided by %d is %d with the remainder dropped.", previous*previous, previous+1, (previous*previous)/(previous+1)))
		},
		func() trio {
			a, b, r := limit(2, 6), limit(2, 6), limit(2, 6)
			return mk(fmt.Sprintf("A Riccati iteration adds %d²/(P + %d) to a constant term, with P = %d. Enter the integer part of the added term.", b, r, a), (b*b)/(a+r),
				fmt.Sprintf("%d divided by (%d + %d) is %d with the remainder dropped.", b*b, a, r, (b*b)/(a+r)))
		},
		func() trio {
			terminal, steps := limit(1, 6), limit(2, 9)
			return mk(fmt.Sprintf("A Riccati recursion starts from terminal weight %d and runs %d backward steps. Enter the number of matrices it produces.", terminal, steps), steps,
				fmt.Sprintf("One matrix per backward step gives %d.", steps))
		},
	},
	"optim.mpc": {
		func() trio {
			horizon, states, inputs := limit(2, 6), limit(1, 4), limit(1, 3)
			return mk(fmt.Sprintf("A predictive controller optimises %d steps with %d state and %d input variables per step. Enter the decision variables.", horizon, states, inputs), horizon*(states+inputs),
				fmt.Sprintf("%d·(%d + %d) = %d.", horizon, states, inputs, horizon*(states+inputs)))
		},
		func() trio {
			horizon, applied := limit(2, 40), 1
			return mk(fmt.Sprintf("A controller plans %d moves and applies %d before replanning. Enter the moves discarded.", horizon, applied), horizon-applied,
				fmt.Sprintf("%d - %d = %d.", horizon, applied, horizon-applied))
		},
		func() trio {
			constraints, horizon := limit(1, 20), limit(2, 6)
			return mk(fmt.Sprintf("Each of %d steps carries %d bound constraints. Enter the constraints in the program.", horizon, constraints), horizon*constraints,
				fmt.Sprintf("%d·%d = %d.", horizon, constraints, horizon*constraints))
		},
	},
	"optim.shooting": {
		func() trio {
			parameters := limit(3, 40)
			return mk(fmt.Sprintf("Single shooting has %d control parameters and evaluates the cost once per trial. Enter the evaluations for one gradient by forward differences.", parameters), parameters+1,
				fmt.Sprintf("%d + 1 = %d evaluations.", parameters, parameters+1))
		},
		func() trio {
			horizon, stages := limit(2, 9), limit(1, 4)
			return mk(fmt.Sprintf("A shooting simulation integrates %d stages over a horizon of %d. Enter the simulated steps.", stages, horizon), stages*horizon,
				fmt.Sprintf("%d·%d = %d steps.", stages, horizon, stages*horizon))
		},
		func() trio {
			sensitivity, step := limit(2, 9), limit(1, 5)
			return mk(fmt.Sprintf("An early control change of %d is amplified by %d over the horizon. Enter the final state change.", sensitivity, step), sensitivity*step,
				fmt.Sprintf("%d·%d = %d.", sensitivity, step, sensitivity*step))
		},
	},
	"optim.collocation": {
		func() trio {
			intervals, degree := limit(2, 8), limit(1, 3)
			return mk(fmt.Sprintf("A collocation mesh has %d intervals with %d nodes inside each. Enter the interior nodes.", intervals, degree), intervals*degree,
				fmt.Sprintf("%d·%d = %d.", intervals, degree, intervals*degree))
		},
		func() trio {
			intervals, shared := limit(2, 40), limit(1, 2)
			return mk(fmt.Sprintf("A mesh has %d intervals and %d shared nodes between neighbours. Enter the distinct nodes.", intervals, shared+1), intervals+shared,
				fmt.Sprintf("%d + %d = %d nodes.", intervals, shared, intervals+shared))
		},
		func() trio {
			variables, constraints := limit(4, 20), limit(2, 10)
			return mk(fmt.Sprintf("A transcribed problem has %d variables and %d equality constraints. Enter the degrees of freedom.", variables, constraints), variables-constraints,
				fmt.Sprintf("%d - %d = %d.", variables, constraints, variables-constraints))
		},
	},
	"optim.costate": {
		func() trio {
			weight, finalState := limit(1, 6), limit(1, 6)
			return mk(fmt.Sprintf("A terminal cost is %d x^2 with x = %d. Enter the terminal costate 2·%d·x.", weight, finalState, weight), 2*weight*finalState,
				fmt.Sprintf("2·%d·%d = %d.", weight, finalState, 2*weight*finalState))
		},
		func() trio {
			rate, horizon := limit(1, 6), limit(2, 8)
			return mk(fmt.Sprintf("A costate decays at rate %d over %d backward steps. Enter the total decrease.", rate, horizon), rate*horizon,
				fmt.Sprintf("%d·%d = %d.", rate, horizon, rate*horizon))
		},
		func() trio {
			jump, costate := limit(1, 5), limit(1, 9)
			return mk(fmt.Sprintf("An active state constraint adds a jump of %d to a costate of %d. Enter the costate after the jump.", jump, costate), costate+jump,
				fmt.Sprintf("%d + %d = %d.", costate, jump, costate+jump))
		},
	},
	"optim.bang_bang": {
		func() trio {
			lower, upper, switches := limit(0, 4), limit(5, 9), limit(1, 3)
			return mk(fmt.Sprintf("A control is bounded between %d and %d and switches %d times from the lower bound. Enter the final value.", lower, upper, switches), upper,
				fmt.Sprintf("Each switch moves between %d and %d, so the value is %d.", lower, upper, upper))
		},
		func() trio {
			switches, horizon := limit(1, 6), limit(2, 9)
			return mk(fmt.Sprintf("A switching function crosses zero %d times over a horizon of %d intervals. Enter the bang-bang arcs.", switches, horizon), switches+1,
				fmt.Sprintf("%d crossings split the horizon into %d arcs.", switches, switches+1))
		},
		func() trio {
			timeAtBound, count := limit(2, 7), limit(2, 5)
			return mk(fmt.Sprintf("A bang-bang trajectory spends %d units at each of %d bounds. Enter the total time.", timeAtBound, count), timeAtBound*count,
				fmt.Sprintf("%d·%d = %d.", timeAtBound, count, timeAtBound*count))
		},
	},
}

func registerWave14(reg *generator.Registry) {
	for id, families := range wave14Families {
		reg.Register(id, &w12Gen{families: families})
	}
	for id, families := range wave14FamiliesB {
		reg.Register(id, &w12Gen{families: families})
	}
	for id, families := range wave14FamiliesC {
		reg.Register(id, &w12Gen{families: families})
	}
}
