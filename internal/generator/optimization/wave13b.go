package optimization

import (
	"fmt"
)

// Wave 13, conic applications and decomposition. The families follow the same
// adapter as wave13.go: three closures, one per ordinal, each computing an
// integer from drawn parameters.
var wave13FamiliesB = map[string][3]func() trio{
	// -----------------------------------------------------------------------
	// Conic applications
	// -----------------------------------------------------------------------
	"optim.lp_as_socp": {
		func() trio {
			constraints := limit(2, 40)
			return mk(fmt.Sprintf("A linear program has %d inequality constraints. Embedding one second-order cone per constraint, enter the number of cones.", constraints), constraints,
				fmt.Sprintf("Each inequality becomes one cone, so there are %d cones.", constraints))
		},
		func() trio {
			free := limit(2, 30)
			return mk(fmt.Sprintf("A cone constraint carries %d free coordinates plus one extra. Enter its dimension.", free), free+1,
				fmt.Sprintf("%d + 1 = %d.", free, free+1))
		},
		func() trio {
			base, scenarios := limit(2, 20), limit(2, 20)
			return mk(fmt.Sprintf("A robust linear program keeps its %d base constraints and adds one per uncertainty scenario, with %d scenarios. Enter the total constraints.", base, scenarios), base+scenarios,
				fmt.Sprintf("%d + %d = %d.", base, scenarios, base+scenarios))
		},
	},
	"optim.qp_as_socp": {
		func() trio {
			vars := limit(2, 30)
			return mk(fmt.Sprintf("A convex quadratic program has %d variables. Its second-order cone formulation uses a cone of dimension n + 2. Enter that dimension.", vars), vars+2,
				fmt.Sprintf("%d + 2 = %d.", vars, vars+2))
		},
		func() trio {
			n := limit(2, 30)
			return mk(fmt.Sprintf("A symmetric matrix of order %d has n(n+1)/2 independent entries. Enter them.", n), n*(n+1)/2,
				fmt.Sprintf("%d·%d/2 = %d.", n, n+1, n*(n+1)/2))
		},
		func() trio {
			cones := limit(2, 30)
			return mk(fmt.Sprintf("A quadratic program has %d linear constraints, each lifted to one cone. Enter the cones.", cones), cones,
				fmt.Sprintf("Each of the %d constraints needs one cone.", cones))
		},
	},
	"optim.robust_lp": {
		func() trio {
			uncertain := limit(1, 20)
			return mk(fmt.Sprintf("A box uncertainty set in %d coefficients has 2^k vertices. Enter its vertex count.", uncertain), 1<<uncertain,
				fmt.Sprintf("2^%d = %d vertices.", uncertain, 1<<uncertain))
		},
		func() trio {
			base, radius, coefficient := limit(2, 30), limit(1, 9), limit(2, 20)
			return mk(fmt.Sprintf("A worst-case value is a + r|b| with a = %d, r = %d and b = %d. Enter it.", base, radius, coefficient), base+radius*coefficient,
				fmt.Sprintf("%d + %d·%d = %d.", base, radius, coefficient, base+radius*coefficient))
		},
		func() trio {
			radius, norm := limit(1, 9), limit(2, 20)
			return mk(fmt.Sprintf("A robust constraint shifts its right-hand side by r·||a|| with r = %d and ||a|| = %d. Enter the shift.", radius, norm), radius*norm,
				fmt.Sprintf("%d times %d is %d.", radius, norm, radius*norm))
		},
	},
	"optim.portfolio_socp": {
		func() trio {
			assets, cap := limit(2, 20), limit(10, 40)
			return mk(fmt.Sprintf("A portfolio holds %d assets, each capped at %d percent. Enter the total cap in percent.", assets, cap), assets*cap,
				fmt.Sprintf("%d times %d is %d percent.", assets, cap, assets*cap))
		},
		func() trio {
			assets, scale := limit(2, 12), limit(1, 9)
			return mk(fmt.Sprintf("Equal weights over %d assets with diagonal risk %d each have a risk of n·m when the variance is n·m². Enter that risk.", assets, scale*scale), assets*scale,
				fmt.Sprintf("The square root of %d·%d is %d.", assets, scale*scale, assets*scale))
		},
		func() trio {
			ratio, budget := limit(1, 20), limit(2, 9)
			return mk(fmt.Sprintf("A return target and a risk budget trade off along a slope of r/s with r = %d and s = %d. Enter the slope.", ratio*budget, budget), ratio,
				fmt.Sprintf("%d divided by %d is %d.", ratio*budget, budget, ratio))
		},
	},
	"optim.sdp_relaxation": {
		func() trio {
			n := limit(2, 30)
			return mk(fmt.Sprintf("A semidefinite relaxation of order %d has n(n+1)/2 variables. Enter them.", n), n*(n+1)/2,
				fmt.Sprintf("%d·%d/2 = %d.", n, n+1, n*(n+1)/2))
		},
		func() trio {
			binary := limit(2, 25)
			return mk(fmt.Sprintf("A semidefinite relaxation of a problem with %d binary variables uses a matrix of that order. Enter its entries.", binary), binary*binary,
				fmt.Sprintf("%d² = %d entries.", binary, binary*binary))
		},
		func() trio {
			bound, order := limit(1, 9), limit(2, 15)
			return mk(fmt.Sprintf("A trace bound of %d times the identity of order %d gives a total trace. Enter it.", bound, order), bound*order,
				fmt.Sprintf("%d·%d = %d.", bound, order, bound*order))
		},
	},
	"optim.max_cut_sdp": {
		func() trio {
			nodes := limit(2, 25)
			return mk(fmt.Sprintf("A complete graph on %d nodes has n(n-1)/2 edges. Enter the edges.", nodes), nodes*(nodes-1)/2,
				fmt.Sprintf("%d·%d/2 = %d.", nodes, nodes-1, nodes*(nodes-1)/2))
		},
		func() trio {
			left, right := limit(2, 20), limit(2, 20)
			return mk(fmt.Sprintf("A cut separates %d nodes from %d nodes in a complete graph. Enter the crossing edges.", left, right), left*right,
				fmt.Sprintf("%d times %d is %d.", left, right, left*right))
		},
		func() trio {
			relaxed, rounded := limit(20, 80), limit(2, 15)
			return mk(fmt.Sprintf("A semidefinite bound is %d and rounding gives %d. Enter the gap.", relaxed, relaxed-rounded), rounded,
				fmt.Sprintf("%d - %d = %d.", relaxed, relaxed-rounded, rounded))
		},
	},
	"optim.matrix_completion": {
		func() trio {
			rows, cols, known := limit(2, 20), limit(2, 20), limit(1, 9)
			return mk(fmt.Sprintf("A %d by %d matrix has %d observed entries. Enter the missing entries.", rows, cols, known), rows*cols-known,
				fmt.Sprintf("%d·%d - %d = %d.", rows, cols, known, rows*cols-known))
		},
		func() trio {
			rank, rows, cols := limit(1, 6), limit(2, 20), limit(2, 20)
			return mk(fmt.Sprintf("A rank-%d completion of a %d by %d matrix stores r(m + n) factor entries. Enter them.", rank, rows, cols), rank*(rows+cols),
				fmt.Sprintf("%d·(%d + %d) = %d.", rank, rows, cols, rank*(rows+cols)))
		},
		func() trio {
			rows, perRow := limit(2, 20), limit(1, 9)
			return mk(fmt.Sprintf("A matrix with %d rows has %d observed entries per row. Enter the observed entries.", rows, perRow), rows*perRow,
				fmt.Sprintf("%d times %d is %d.", rows, perRow, rows*perRow))
		},
	},
	"optim.certificate_check": {
		func() trio {
			primal, dual := limit(20, 90), limit(2, 18)
			return mk(fmt.Sprintf("A primal value is %d and a dual value is %d. Enter the duality gap.", primal, primal-dual), dual,
				fmt.Sprintf("%d - %d = %d.", primal, primal-dual, dual))
		},
		func() trio {
			multipliers := limit(2, 30)
			return mk(fmt.Sprintf("A certificate lists %d multipliers. Enter how many must be nonnegative.", multipliers), multipliers,
				fmt.Sprintf("All %d inequality multipliers are nonnegative.", multipliers))
		},
		func() trio {
			primal, dual := limit(30, 80), limit(1, 25)
			return mk(fmt.Sprintf("A certificate checks weak duality with primal %d and dual %d. Enter the slack.", dual+primal, dual), primal,
				fmt.Sprintf("%d - %d = %d, which weak duality keeps nonnegative.", dual+primal, dual, primal))
		},
	},
	"optim.conic_duality": {
		func() trio {
			dim := limit(2, 40)
			return mk(fmt.Sprintf("A cone in R^%d and its dual both have dimension n. Enter their combined dimension.", dim), 2*dim,
				fmt.Sprintf("2·%d = %d.", dim, 2*dim))
		},
		func() trio {
			cones, dim := limit(2, 15), limit(2, 15)
			return mk(fmt.Sprintf("A problem has %d self-dual cones of dimension %d. Enter the combined dimension.", cones, dim), cones*dim,
				fmt.Sprintf("%d times %d is %d.", cones, dim, cones*dim))
		},
		func() trio {
			closed, halvings := limit(1, 20), limit(1, 8)
			return mk(fmt.Sprintf("A conic gap of %d is halved %d times. Enter the final gap.", closed<<halvings, halvings), closed,
				fmt.Sprintf("%d divided by %d is %d.", closed<<halvings, 1<<halvings, closed))
		},
	},
	"optim.cone_membership": {
		func() trio {
			axis, norm := limit(3, 25), limit(1, 20)
			limitAxis := axis + norm
			return mk(fmt.Sprintf("A point lies on the axis at %d with a perpendicular norm of %d. Enter its second-order-cone slack x0 - ||x||.", limitAxis, norm), axis,
				fmt.Sprintf("%d - %d = %d, which is nonnegative so the point is in the cone.", limitAxis, norm, axis))
		},
		func() trio {
			axis, norm := limit(3, 20), limit(1, 15)
			return mk(fmt.Sprintf("A second-order-cone point has axis %d and squared perpendicular norm %d². Enter the squared slack t² - ||x||².", axis, norm), axis*axis-norm*norm,
				fmt.Sprintf("%d - %d = %d.", axis*axis, norm*norm, axis*axis-norm*norm))
		},
		func() trio {
			dim := limit(2, 30)
			return mk(fmt.Sprintf("A second-order cone of dimension %d has one axis coordinate. Enter the number of free coordinates.", dim), dim-1,
				fmt.Sprintf("%d - 1 = %d.", dim, dim-1))
		},
	},

	// -----------------------------------------------------------------------
	// Decomposition methods
	// -----------------------------------------------------------------------
	"optim.dual_decomposition": {
		func() trio {
			problems, each := limit(2, 15), limit(2, 12)
			return mk(fmt.Sprintf("Dual decomposition splits a problem into %d subproblems of %d variables. Enter the total variables.", problems, each), problems*each,
				fmt.Sprintf("%d times %d is %d.", problems, each, problems*each))
		},
		func() trio {
			price, step, residual := limit(2, 20), limit(1, 5), limit(-6, 6)
			return mk(fmt.Sprintf("A dual price %d is updated by step·residual with step %d and residual %d. Enter the new price.", price, step, residual), price+step*residual,
				fmt.Sprintf("%d + %d·(%d) = %d.", price, step, residual, price+step*residual))
		},
		func() trio {
			problems, rounds := limit(2, 15), limit(2, 12)
			return mk(fmt.Sprintf("%d subproblems are each solved %d times. Enter the solves.", problems, rounds), problems*rounds,
				fmt.Sprintf("%d times %d is %d.", problems, rounds, problems*rounds))
		},
	},
	"optim.primal_decomposition": {
		func() trio {
			blocks, size := limit(2, 15), limit(2, 12)
			return mk(fmt.Sprintf("Primal decomposition splits a problem into %d blocks of size %d. Enter the total variables.", blocks, size), blocks*size,
				fmt.Sprintf("%d times %d is %d.", blocks, size, blocks*size))
		},
		func() trio {
			coupling, blocks := limit(2, 20), limit(2, 12)
			return mk(fmt.Sprintf("A problem has %d coupling constraints per block and %d blocks. Enter the coupling rows.", coupling, blocks), coupling*blocks,
				fmt.Sprintf("%d times %d is %d.", coupling, blocks, coupling*blocks))
		},
		func() trio {
			blocks, sweeps := limit(2, 12), limit(2, 10)
			return mk(fmt.Sprintf("A decomposition performs %d sweeps over %d blocks. Enter the block updates.", sweeps, blocks), sweeps*blocks,
				fmt.Sprintf("%d times %d is %d.", sweeps, blocks, sweeps*blocks))
		},
	},
	"optim.consensus_admm": {
		func() trio {
			agents := limit(2, 30)
			return mk(fmt.Sprintf("Consensus ADMM runs over %d agents, each holding a local copy. Enter the copies.", agents), agents,
				fmt.Sprintf("Each of the %d agents keeps one copy.", agents))
		},
		func() trio {
			first, second := 2*limit(1, 15), 2*limit(1, 15)
			return mk(fmt.Sprintf("Two agents hold %d and %d, and consensus averages them. Enter the average.", first, second), (first+second)/2,
				fmt.Sprintf("(%d + %d)/2 = %d.", first, second, (first+second)/2))
		},
		func() trio {
			dual, update := limit(1, 20), limit(-5, 5)
			return mk(fmt.Sprintf("A consensus dual variable %d receives the update %d. Enter its new value.", dual, update), dual+update,
				fmt.Sprintf("%d + (%d) = %d.", dual, update, dual+update))
		},
	},
	"optim.alternating_direction": {
		func() trio {
			iterations := limit(2, 30)
			return mk(fmt.Sprintf("ADMM performs two updates per iteration. Over %d iterations, enter the updates.", iterations), 2*iterations,
				fmt.Sprintf("2·%d = %d.", iterations, 2*iterations))
		},
		func() trio {
			rho, residual := 2*limit(1, 6), limit(1, 9)
			return mk(fmt.Sprintf("The augmented term is (ρ/2)r² with ρ = %d and residual %d. Enter it.", rho, residual), rho*residual*residual/2,
				fmt.Sprintf("%d·%d/2 = %d.", rho, residual*residual, rho*residual*residual/2))
		},
		func() trio {
			primal, dual := limit(2, 30), limit(2, 30)
			return mk(fmt.Sprintf("A primal residual is %d and a dual residual is %d. Enter the larger.", primal, dual), maxInt(primal, dual),
				fmt.Sprintf("The larger of %d and %d is %d.", primal, dual, maxInt(primal, dual)))
		},
	},
	"optim.block_coordinate": {
		func() trio {
			blocks, size := limit(2, 15), limit(2, 12)
			return mk(fmt.Sprintf("Block coordinate descent works on %d blocks of %d variables. Enter the total variables.", blocks, size), blocks*size,
				fmt.Sprintf("%d times %d is %d.", blocks, size, blocks*size))
		},
		func() trio {
			blocks, sweeps := limit(2, 12), limit(2, 12)
			return mk(fmt.Sprintf("A sweep updates %d blocks, and %d sweeps are taken. Enter the updates.", blocks, sweeps), blocks*sweeps,
				fmt.Sprintf("%d times %d is %d.", blocks, sweeps, blocks*sweeps))
		},
		func() trio {
			blocks, size := limit(2, 15), limit(2, 12)
			return mk(fmt.Sprintf("A block method has %d blocks of size %d. Enter the variables.", blocks, size), blocks*size,
				fmt.Sprintf("%d times %d is %d.", blocks, size, blocks*size))
		},
	},
	"optim.jacobi_iteration": {
		func() trio {
			equations, sweeps := limit(2, 25), limit(2, 10)
			return mk(fmt.Sprintf("A Jacobi solve on %d equations takes %d sweeps. Enter the equation updates.", equations, sweeps), equations*sweeps,
				fmt.Sprintf("%d times %d is %d.", equations, sweeps, equations*sweeps))
		},
		func() trio {
			offDiagonal, diagonal, step := limit(1, 20), limit(1, 9), limit(1, 6)
			right := offDiagonal + step*diagonal
			return mk(fmt.Sprintf("A Jacobi update is (b - s)/d with off-diagonal sum %d, diagonal %d and right-hand side %d. Enter the update.", offDiagonal, diagonal, right), step,
				fmt.Sprintf("(%d - %d)/%d = %d.", right, offDiagonal, diagonal, step))
		},
		func() trio {
			unknowns := limit(2, 30)
			return mk(fmt.Sprintf("A Jacobi sweep reads %d old values and writes %d new ones. Enter the total touches.", unknowns, unknowns), 2*unknowns,
				fmt.Sprintf("2·%d = %d.", unknowns, 2*unknowns))
		},
	},
	"optim.gauss_seidel": {
		func() trio {
			unknowns := limit(2, 40)
			return mk(fmt.Sprintf("Gauss-Seidel updates its %d unknowns in sequence. Enter the updates per sweep.", unknowns), unknowns,
				fmt.Sprintf("Each of the %d unknowns is updated once.", unknowns))
		},
		func() trio {
			unknowns, sweeps := limit(2, 20), limit(2, 12)
			return mk(fmt.Sprintf("%d unknowns are updated over %d sweeps. Enter the total updates.", unknowns, sweeps), unknowns*sweeps,
				fmt.Sprintf("%d times %d is %d.", unknowns, sweeps, unknowns*sweeps))
		},
		func() trio {
			n := limit(2, 30)
			return mk(fmt.Sprintf("A lower-triangular matrix of order %d has n(n+1)/2 nonzeros. Enter them.", n), n*(n+1)/2,
				fmt.Sprintf("%d·%d/2 = %d.", n, n+1, n*(n+1)/2))
		},
	},
	"optim.column_generation": {
		func() trio {
			rows := limit(2, 40)
			return mk(fmt.Sprintf("A restricted master has %d rows, each with one dual price. Enter the prices.", rows), rows,
				fmt.Sprintf("One price per row gives %d.", rows))
		},
		func() trio {
			iterations, columns := limit(2, 20), limit(1, 9)
			return mk(fmt.Sprintf("Each of %d iterations adds %d columns. Enter the added columns.", iterations, columns), iterations*columns,
				fmt.Sprintf("%d times %d is %d.", iterations, columns, iterations*columns))
		},
		func() trio {
			master, added := limit(2, 20), limit(1, 9)
			return mk(fmt.Sprintf("A restricted master of %d columns is extended by %d. Enter the total columns.", master, added), master+added,
				fmt.Sprintf("%d + %d = %d.", master, added, master+added))
		},
	},
	"optim.dantzig_wolfe": {
		func() trio {
			proposals := limit(2, 30)
			return mk(fmt.Sprintf("A Dantzig-Wolfe master has %d proposals, each with one weight. Enter the weights.", proposals), proposals,
				fmt.Sprintf("One weight per proposal gives %d.", proposals))
		},
		func() trio {
			iterations, columns := limit(2, 20), limit(1, 9)
			return mk(fmt.Sprintf("Each of %d iterations adds %d columns. Enter the added columns.", iterations, columns), iterations*columns,
				fmt.Sprintf("%d times %d is %d.", iterations, columns, iterations*columns))
		},
		func() trio {
			objective, reduced := limit(5, 60), limit(1, 20)
			return mk(fmt.Sprintf("A master objective is %d and a subproblem returns reduced cost %d. Enter their sum.", objective, reduced), objective+reduced,
				fmt.Sprintf("%d + %d = %d.", objective, reduced, objective+reduced))
		},
	},
	"optim.benders": {
		func() trio {
			iterations := limit(2, 30)
			return mk(fmt.Sprintf("Benders adds one cut per iteration. After %d iterations, enter the cuts.", iterations), iterations,
				fmt.Sprintf("Each iteration adds one cut, so there are %d.", iterations))
		},
		func() trio {
			vars, subproblems := limit(2, 20), limit(2, 15)
			return mk(fmt.Sprintf("A Benders master has %d variables and drives %d subproblems. Enter the total.", vars, subproblems), vars+subproblems,
				fmt.Sprintf("%d + %d = %d.", vars, subproblems, vars+subproblems))
		},
		func() trio {
			upper, lower := limit(20, 90), limit(2, 18)
			return mk(fmt.Sprintf("An upper bound is %d and a lower bound is %d. Enter the gap.", upper, upper-lower), lower,
				fmt.Sprintf("%d - %d = %d.", upper, upper-lower, lower))
		},
	},
}
