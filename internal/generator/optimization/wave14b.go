package optimization

import (
	"fmt"
)

// wave14FamiliesB covers learning-oriented optimisation and integer models.
var wave14FamiliesB = map[string][3]func() trio{
	// -----------------------------------------------------------------------
	// Learning-oriented optimisation
	// -----------------------------------------------------------------------
	"optim.implicit_diff": {
		func() trio {
			slope, hidden := limit(2, 9), limit(2, 9)
			return mk(fmt.Sprintf("A layer solves %d y + %d = 0. Enter the integer solution y.", slope, -slope*hidden), hidden,
				fmt.Sprintf("y = -(%d)/%d = %d.", -slope*hidden, slope, hidden))
		},
		func() trio {
			jacobian, rhs := limit(2, 9), limit(2, 9)
			return mk(fmt.Sprintf("An implicit derivative solves %d z = %d. Enter z.", jacobian, jacobian*rhs), rhs,
				fmt.Sprintf("z = %d/%d = %d.", jacobian*rhs, jacobian, rhs))
		},
		func() trio {
			iterations, rhs := limit(2, 40), limit(1, 6)
			return mk(fmt.Sprintf("A forward solve ran %d iterations and the backward pass reuses one factorisation for %d right-hand sides. Enter the backward solves.", iterations, rhs), rhs,
				fmt.Sprintf("One factorisation serves %d right-hand sides, so the count is %d regardless of the %d forward iterations.", rhs, rhs, iterations))
		},
	},
	"optim.bilevel": {
		func() trio {
			outer, inner := limit(2, 6), limit(2, 6)
			return mk(fmt.Sprintf("A bilevel search tries %d outer candidates and solves %d inner options for each. Enter the inner solves.", outer, inner), outer*inner,
				fmt.Sprintf("%d·%d = %d.", outer, inner, outer*inner))
		},
		func() trio {
			innerBest, penalty := limit(4, 15), limit(1, 6)
			return mk(fmt.Sprintf("An inner problem returns its optimum %d and the outer adds a penalty of %d. Enter the outer objective.", innerBest, penalty), innerBest+penalty,
				fmt.Sprintf("%d + %d = %d.", innerBest, penalty, innerBest+penalty))
		},
		func() trio {
			constraints, variables := limit(1, 4), limit(2, 6)
			return mk(fmt.Sprintf("An inner problem has %d equality constraints and %d variables. Enter the multiplier count in its single-level form.", constraints, variables), constraints,
				fmt.Sprintf("One multiplier per constraint gives %d.", constraints))
		},
	},
	"optim.meta_gradient": {
		func() trio {
			steps, parameters := limit(2, 8), limit(2, 9)
			return mk(fmt.Sprintf("A meta-gradient unrolls %d training steps with %d parameters. Enter the stored activation units.", steps, parameters), steps*parameters,
				fmt.Sprintf("%d steps times %d parameters is %d.", steps, parameters, steps*parameters))
		},
		func() trio {
			rate, loss := limit(1, 5), limit(2, 9)
			return mk(fmt.Sprintf("A meta-objective changes by %d per unit learning rate and the rate changes by %d. Enter the first-order change.", rate, loss), rate*loss,
				fmt.Sprintf("%d·%d = %d.", rate, loss, rate*loss))
		},
		func() trio {
			steps, checkpoint := limit(4, 16), limit(1, 4)
			return mk(fmt.Sprintf("An unrolled pass over %d steps checkpoints every %d steps. Enter the checkpoints stored.", steps, checkpoint), steps/checkpoint,
				fmt.Sprintf("%d divided by %d is %d checkpoints.", steps, checkpoint, steps/checkpoint))
		},
	},
	"optim.hypergradient": {
		func() trio {
			hessian, gradient := limit(2, 9), limit(2, 9)
			return mk(fmt.Sprintf("A hypergradient solves h·z = g with h = %d and g = %d. Enter z.", hessian, hessian*gradient), gradient,
				fmt.Sprintf("z = %d/%d = %d.", hessian*gradient, hessian, gradient))
		},
		func() trio {
			parameters, products := limit(2, 9), limit(1, 4)
			return mk(fmt.Sprintf("A hypergradient needs %d Hessian-vector products for each of %d validation batches. Enter the products.", parameters, products), parameters*products,
				fmt.Sprintf("%d·%d = %d.", parameters, products, parameters*products))
		},
		func() trio {
			rate, weight := limit(1, 5), limit(2, 9)
			return mk(fmt.Sprintf("A regularisation weight %d is reduced by %d per hypergradient step. Enter the new weight.", weight, rate), weight-rate,
				fmt.Sprintf("%d - %d = %d.", weight, rate, weight-rate))
		},
	},
	"optim.diff_layers": {
		func() trio {
			value, low, high := limit(-6, 12), 0, limit(1, 9)
			clamped := value
			if clamped < low {
				clamped = low
			}
			if clamped > high {
				clamped = high
			}
			return mk(fmt.Sprintf("A projection layer onto [%d, %d] receives %d. Enter the output.", low, high, value), clamped,
				fmt.Sprintf("Clamping %d to [%d, %d] gives %d.", value, low, high, clamped))
		},
		func() trio {
			inputs, outputs := limit(1, 5), limit(1, 6)
			return mk(fmt.Sprintf("A layer solves a problem with %d inputs and %d outputs. Enter the Jacobian entries.", inputs, outputs), inputs*outputs,
				fmt.Sprintf("%d·%d = %d entries.", inputs, outputs, inputs*outputs))
		},
		func() trio {
			regulariser, scale := limit(1, 6), limit(2, 9)
			return mk(fmt.Sprintf("A layer adds a regulariser %d times the squared norm %d. Enter the added term.", regulariser, scale), regulariser*scale,
				fmt.Sprintf("%d·%d = %d.", regulariser, scale, regulariser*scale))
		},
	},
	"optim.influence": {
		func() trio {
			gradient, curvature := limit(1, 9), limit(2, 9)
			return mk(fmt.Sprintf("An influence estimate is -g/h with g = %d and h = %d. Enter it as an integer.", gradient*curvature, curvature), -gradient,
				fmt.Sprintf("-(%d·%d)/%d = %d.", gradient, curvature, curvature, -gradient))
		},
		func() trio {
			points, products := limit(2, 8), limit(1, 5)
			return mk(fmt.Sprintf("An influence study scores %d training points and uses %d Hessian-vector products each. Enter the products.", points, products), points*products,
				fmt.Sprintf("%d·%d = %d.", points, products, points*products))
		},
		func() trio {
			weight, change := limit(2, 9), limit(1, 4)
			return mk(fmt.Sprintf("Removing a point of weight %d changes the loss by %d per unit influence. Enter the total change.", weight, change), weight*change,
				fmt.Sprintf("%d·%d = %d.", weight, change, weight*change))
		},
	},
	"optim.sam": {
		func() trio {
			point, radius := limit(2, 9), limit(1, 4)
			return mk(fmt.Sprintf("A sharpness step perturbs a parameter %d by a radius of %d along its gradient. Enter the perturbed value.", point, radius), point+radius,
				fmt.Sprintf("%d + %d = %d.", point, radius, point+radius))
		},
		func() trio {
			steps, gradients := limit(2, 40), 2
			return mk(fmt.Sprintf("Each of %d sharpness steps evaluates the gradient %d times. Enter the evaluations.", steps, gradients), steps*gradients,
				fmt.Sprintf("%d·%d = %d.", steps, gradients, steps*gradients))
		},
		func() trio {
			loss, increase := limit(3, 12), limit(1, 4)
			return mk(fmt.Sprintf("The worst loss inside a ball is %d and the clean loss is %d. Enter the sharpness gap.", loss+increase, increase), increase,
				fmt.Sprintf("%d - %d = %d.", loss+increase, loss, increase))
		},
	},
	"optim.adversarial_training": {
		func() trio {
			input, epsilon, sign := limit(2, 9), limit(1, 3), limit(-1, 1)
			if sign == 0 {
				sign = 1
			}
			return mk(fmt.Sprintf("A one-step attack moves an input %d by %d times the sign %d. Enter the perturbed input.", input, epsilon, sign), input+epsilon*sign,
				fmt.Sprintf("%d + %d·%d = %d.", input, epsilon, sign, input+epsilon*sign))
		},
		func() trio {
			steps, batch := limit(2, 6), limit(1, 5)
			return mk(fmt.Sprintf("An attack takes %d ascent steps over a batch of %d examples. Enter the gradient evaluations.", steps, batch), steps*batch,
				fmt.Sprintf("%d·%d = %d.", steps, batch, steps*batch))
		},
		func() trio {
			clean, drop := limit(6, 15), limit(1, 4)
			return mk(fmt.Sprintf("A robust model scores %d on clean data and %d below that on perturbed data. Enter the perturbed score.", clean, drop), clean-drop,
				fmt.Sprintf("%d - %d = %d.", clean, drop, clean-drop))
		},
	},
	"optim.dro": {
		func() trio {
			loss, worst, radius := limit(2, 9), limit(2, 9), limit(1, 3)
			return mk(fmt.Sprintf("A robust risk is (1 - r)L + r W with L = %d, W = %d and r = 1/%d. Enter the numerator of the weighted risk before dividing by %d.", loss, worst, radius, radius), (radius-1)*loss+worst,
				fmt.Sprintf("(%d - 1)·%d + %d = %d, and dividing by %d gives the risk.", radius, loss, worst, (radius-1)*loss+worst, radius))
		},
		func() trio {
			radius, shift := limit(2, 6), limit(1, 5)
			return mk(fmt.Sprintf("A Wasserstein radius of %d is scaled by %d. Enter the effective penalty.", radius, shift), radius*shift,
				fmt.Sprintf("%d·%d = %d.", radius, shift, radius*shift))
		},
		func() trio {
			average, gap := limit(2, 9), limit(1, 6)
			return mk(fmt.Sprintf("An empirical risk is %d and the worst-case risk is %d higher. Enter the worst-case risk.", average, gap), average+gap,
				fmt.Sprintf("%d + %d = %d.", average, gap, average+gap))
		},
	},
	"optim.gradient_penalty": {
		func() trio {
			loss, weight, gradient := limit(2, 9), limit(1, 4), limit(1, 5)
			return mk(fmt.Sprintf("A penalised objective is L + λ||g||² with L = %d, λ = %d and ||g|| = %d. Enter it.", loss, weight, gradient), loss+weight*gradient*gradient,
				fmt.Sprintf("%d + %d·%d² = %d.", loss, weight, gradient, loss+weight*gradient*gradient))
		},
		func() trio {
			points, cost := limit(2, 9), limit(1, 4)
			return mk(fmt.Sprintf("A penalty is evaluated at %d sampled points costing %d each. Enter the evaluations.", points, cost), points*cost,
				fmt.Sprintf("%d·%d = %d.", points, cost, points*cost))
		},
		func() trio {
			slope, weight := limit(2, 20), limit(1, 3)
			return mk(fmt.Sprintf("A linear model has slope %d and a penalty weight of %d. Enter the penalty on the slope.", slope, weight), weight*slope*slope,
				fmt.Sprintf("%d·%d² = %d.", weight, slope, weight*slope*slope))
		},
	},
	// -----------------------------------------------------------------------
	// Integer models
	// -----------------------------------------------------------------------
	"optim.branch_cut": {
		func() trio {
			depth := limit(2, 40)
			total := 1
			for i := 0; i < depth; i++ {
				total *= 2
			}
			return mk(fmt.Sprintf("A full binary branch tree has depth %d. Enter its node count.", depth), total-1,
				fmt.Sprintf("A full tree of depth %d has 2^%d - 1 = %d nodes.", depth, depth, total-1))
		},
		func() trio {
			cuts, nodes := limit(1, 4), limit(2, 9)
			return mk(fmt.Sprintf("Each of %d nodes adds %d cuts. Enter the cuts added.", nodes, cuts), cuts*nodes,
				fmt.Sprintf("%d·%d = %d.", cuts, nodes, cuts*nodes))
		},
		func() trio {
			bound, incumbent := limit(8, 20), limit(1, 6)
			return mk(fmt.Sprintf("A node bound is %d and the incumbent is %d below it. Enter the incumbent value.", bound, incumbent), bound-incumbent,
				fmt.Sprintf("%d - %d = %d.", bound, incumbent, bound-incumbent))
		},
	},
	"optim.lagrangian_relaxation": {
		func() trio {
			cost, multiplier, residual := limit(2, 9), limit(1, 5), limit(1, 6)
			return mk(fmt.Sprintf("A Lagrangian is cost %d plus multiplier %d times residual %d. Enter its value.", cost, multiplier, residual), cost+multiplier*residual,
				fmt.Sprintf("%d + %d·%d = %d.", cost, multiplier, residual, cost+multiplier*residual))
		},
		func() trio {
			relaxed, penalty := limit(6, 15), limit(1, 4)
			return mk(fmt.Sprintf("A relaxed value is %d and a feasibility repair raises the cost by %d. Enter the repaired cost.", relaxed, penalty), relaxed+penalty,
				fmt.Sprintf("%d + %d = %d.", relaxed, penalty, relaxed+penalty))
		},
		func() trio {
			constraints, blocks := limit(1, 3), limit(2, 20)
			return mk(fmt.Sprintf("Relaxing %d coupling constraints leaves %d independent blocks. Enter the total subproblems.", constraints, blocks), blocks,
				fmt.Sprintf("The %d remaining blocks are solved independently, so the count is %d.", blocks, blocks))
		},
	},
	"optim.subgradient_lagrangian": {
		func() trio {
			multiplier, step, residual := limit(1, 5), limit(1, 4), limit(1, 6)
			return mk(fmt.Sprintf("A dual update is λ + t·r with λ = %d, t = %d and r = %d. Enter the new multiplier.", multiplier, step, residual), multiplier+step*residual,
				fmt.Sprintf("%d + %d·%d = %d.", multiplier, step, residual, multiplier+step*residual))
		},
		func() trio {
			start, decrease := limit(8, 20), limit(1, 5)
			return mk(fmt.Sprintf("A dual value rises from %d by %d in one step. Enter the new value.", start, decrease), start+decrease,
				fmt.Sprintf("%d + %d = %d.", start, decrease, start+decrease))
		},
		func() trio {
			iterations, perStep := limit(2, 9), limit(1, 3)
			return mk(fmt.Sprintf("A multiplier search runs %d iterations and uses %d subgradient evaluations each. Enter the evaluations.", iterations, perStep), iterations*perStep,
				fmt.Sprintf("%d·%d = %d.", iterations, perStep, iterations*perStep))
		},
	},
	"optim.gomory_cut": {
		func() trio {
			whole, fraction := limit(1, 8), limit(1, 9)
			return mk(fmt.Sprintf("A tableau right-hand side is %d.%d. Enter its integer part.", whole, fraction), whole,
				fmt.Sprintf("The integer part of %d.%d is %d.", whole, fraction, whole))
		},
		func() trio {
			basic, fractional := limit(2, 20), limit(1, 3)
			return mk(fmt.Sprintf("A basis has %d basic variables and %d of them are fractional. Enter the integer basics.", basic, fractional), basic-fractional,
				fmt.Sprintf("%d - %d = %d.", basic, fractional, basic-fractional))
		},
		func() trio {
			rounds, cuts := limit(1, 4), limit(1, 5)
			return mk(fmt.Sprintf("A solver adds %d cuts per round over %d rounds. Enter the cuts added.", cuts, rounds), cuts*rounds,
				fmt.Sprintf("%d·%d = %d.", cuts, rounds, cuts*rounds))
		},
	},
	"optim.cutting_stock": {
		func() trio {
			pieces, length, roll := limit(2, 6), limit(2, 6), 10
			return mk(fmt.Sprintf("An order needs %d pieces of length %d from rolls of length %d. Enter the rolls needed when there is no waste.", pieces, length, roll), (pieces*length+roll-1)/roll,
				fmt.Sprintf("%d·%d = %d units, and %d rolls hold %d.", pieces, length, pieces*length, (pieces*length+roll-1)/roll, ((pieces*length+roll-1)/roll)*roll))
		},
		func() trio {
			pieces, perPattern := limit(2, 20), limit(1, 3)
			return mk(fmt.Sprintf("A pattern yields %d pieces and %d patterns are used. Enter the pieces produced.", perPattern, pieces), pieces*perPattern,
				fmt.Sprintf("%d·%d = %d.", pieces, perPattern, pieces*perPattern))
		},
		func() trio {
			demand, perRoll, rolls := limit(4, 12), limit(2, 5), limit(2, 5)
			return mk(fmt.Sprintf("A stock of %d rolls holds %d pieces each and the demand is %d. Enter the spare pieces.", rolls, perRoll, demand), rolls*perRoll-demand,
				fmt.Sprintf("%d·%d - %d = %d.", rolls, perRoll, demand, rolls*perRoll-demand))
		},
	},
	"optim.facility_location": {
		func() trio {
			fixed, open, assign, customers := limit(3, 9), limit(1, 3), limit(1, 4), limit(2, 6)
			return mk(fmt.Sprintf("%d facilities cost %d each to open and %d customers cost %d each to assign. Enter the total cost.", open, fixed, customers, assign), open*fixed+customers*assign,
				fmt.Sprintf("%d·%d + %d·%d = %d.", open, fixed, customers, assign, open*fixed+customers*assign))
		},
		func() trio {
			customers, open := limit(4, 12), limit(1, 4)
			return mk(fmt.Sprintf("%d customers are served by %d open facilities. Enter the assignments.", customers, open), customers,
				fmt.Sprintf("Every customer gets exactly one assignment, so the count is %d.", customers))
		},
		func() trio {
			open, closed := limit(1, 20), limit(1, 4)
			return mk(fmt.Sprintf("A network has %d open and %d closed facilities. Enter the total facilities.", open, closed), open+closed,
				fmt.Sprintf("%d + %d = %d.", open, closed, open+closed))
		},
	},
	"optim.tsp_relaxation": {
		func() trio {
			cities := limit(3, 40)
			return mk(fmt.Sprintf("A tour visits %d cities and returns to the start. Enter its edge count.", cities), cities,
				fmt.Sprintf("A cycle through %d cities has %d edges.", cities, cities))
		},
		func() trio {
			cities := limit(3, 40)
			return mk(fmt.Sprintf("The degree relaxation gives two edges per city for %d cities. Enter the total edge endpoints.", cities), 2*cities,
				fmt.Sprintf("2·%d = %d endpoints, which is %d edges.", cities, 2*cities, cities))
		},
		func() trio {
			cycles, each := limit(2, 4), limit(2, 20)
			return mk(fmt.Sprintf("A relaxed solution splits %d cities into %d disjoint cycles of %d cities. Enter the edges on the cycles.", cycles*each, each, each), cycles*each,
				fmt.Sprintf("%d cycles of %d cities carry %d edges in total.", cycles, each, cycles*each))
		},
	},
	"optim.quadratic_assignment": {
		func() trio {
			flow, distance, pairs := limit(1, 5), limit(1, 5), limit(2, 5)
			return mk(fmt.Sprintf("A quadratic assignment has %d pairs with flow %d and distance %d. Enter its cost.", pairs, flow, distance), pairs*flow*distance,
				fmt.Sprintf("%d·%d·%d = %d.", pairs, flow, distance, pairs*flow*distance))
		},
		func() trio {
			facilities := limit(2, 40)
			return mk(fmt.Sprintf("An assignment of %d facilities to %d locations is written as a permutation matrix. Enter the number of ones.", facilities, facilities), facilities,
				fmt.Sprintf("A permutation matrix of order %d has %d ones.", facilities, facilities))
		},
		func() trio {
			pairs, variables := limit(2, 4), limit(2, 20)
			return mk(fmt.Sprintf("A linearisation uses %d pair variables for each of %d location pairs. Enter the variables.", pairs, variables), pairs*variables,
				fmt.Sprintf("%d·%d = %d.", pairs, variables, pairs*variables))
		},
	},
	"optim.big_m": {
		func() trio {
			upper, rhs := limit(6, 15), limit(1, 5)
			return mk(fmt.Sprintf("A constraint activity is bounded by %d and its right-hand side is %d. Enter the smallest valid M.", upper, rhs), upper-rhs,
				fmt.Sprintf("%d - %d = %d.", upper, rhs, upper-rhs))
		},
		func() trio {
			value, offset := limit(1, 9), limit(1, 5)
			return mk(fmt.Sprintf("A switched constraint reads x + M z ≤ b with x = %d and M = %d. Enter the left-hand side when z = 1 and b is the bound %d.", value, offset, value+offset), value+offset,
				fmt.Sprintf("%d + %d = %d, which meets the bound.", value, offset, value+offset))
		},
		func() trio {
			constraints, binaries := limit(2, 6), limit(1, 4)
			return mk(fmt.Sprintf("A model has %d switched constraints and %d binary switches. Enter the binary variables.", constraints, binaries), binaries,
				fmt.Sprintf("One binary serves every switched row, so the count is %d.", binaries))
		},
	},
	"optim.sos1": {
		func() trio {
			sets, members := limit(1, 20), limit(2, 9)
			return mk(fmt.Sprintf("A model has %d SOS1 sets of %d members each. Enter the total members.", sets, members), sets*members,
				fmt.Sprintf("%d·%d = %d.", sets, members, sets*members))
		},
		func() trio {
			alternatives := limit(3, 40)
			return mk(fmt.Sprintf("A choice among %d alternatives needs %d binaries in a big-M model and one SOS1 set. Enter the binaries the set saves.", alternatives, alternatives), alternatives-1,
				fmt.Sprintf("%d binaries are replaced by one set, saving %d.", alternatives, alternatives-1))
		},
		func() trio {
			breakpoints := limit(2, 40)
			return mk(fmt.Sprintf("A piecewise linear function has %d breakpoints. Enter the segments in its SOS1 formulation.", breakpoints), breakpoints-1,
				fmt.Sprintf("%d breakpoints give %d segments.", breakpoints, breakpoints-1))
		},
	},
}
