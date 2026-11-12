package optimization

import "fmt"

var wave12FamiliesB = map[string][3]func() trio{
	// -----------------------------------------------------------------------
	// Multiobjective optimization
	// -----------------------------------------------------------------------
	"optim.pareto_dominance": {
		func() trio {
			a1, a2 := limit(2, 9), limit(2, 9)
			b1, b2 := a1+limit(1, 3), a2+limit(1, 3)
			return mk(fmt.Sprintf("Point A has objectives (%d, %d) and point B has (%d, %d), both minimized. Enter how much larger B is than A in the first objective.", a1, a2, b1, b2), b1-a1,
				fmt.Sprintf("B is worse in both objectives, by %d in the first and %d in the second, so A dominates B.", b1-a1, b2-a2))
		},
		func() trio {
			a1, a2, b1, b2 := limit(1, 6), limit(1, 6), limit(1, 6), limit(1, 6)
			worse := 0
			if b1 > a1 {
				worse++
			}
			if b2 > a2 {
				worse++
			}
			return mk(fmt.Sprintf("A has objectives (%d, %d) and B has (%d, %d), both minimized. Enter how many objectives B is worse in.", a1, a2, b1, b2), worse,
				fmt.Sprintf("B is worse in %d of the two objectives, so dominance needs both.", worse))
		},
		func() trio {
			a1, a2, b1, b2 := limit(1, 9), limit(1, 9), limit(1, 9), limit(1, 9)
			sumA, sumB := a1+a2, b1+b2
			return mk(fmt.Sprintf("A has objectives (%d, %d) and B has (%d, %d). Enter B's objective sum minus A's.", a1, a2, b1, b2), sumB-sumA,
				fmt.Sprintf("%d + %d - (%d + %d) = %d.", b1, b2, a1, a2, sumB-sumA))
		},
	},
	"optim.pareto_front": {
		func() trio {
			a1, a2 := limit(1, 8), limit(1, 8)
			b1, b2 := limit(1, 9), limit(1, 9)
			nondominated := 0
			if !(a1 <= b1 && a2 <= b2 && (a1 < b1 || a2 < b2)) {
				nondominated++
			}
			if !(b1 <= a1 && b2 <= a2 && (b1 < a1 || b2 < a2)) {
				nondominated++
			}
			return mk(fmt.Sprintf("Two points (%d, %d) and (%d, %d) are minimized. Enter how many of them lie on the Pareto front.", a1, a2, b1, b2), nondominated,
				fmt.Sprintf("A point is on the front when the other does not dominate it; here %d of the two qualify.", nondominated))
		},
		func() trio {
			a1, a2, b1 := limit(1, 6), limit(1, 6), limit(1, 6)
			b2 := a2 + limit(1, 4)
			return mk(fmt.Sprintf("A front point is (%d, %d). A new point is (%d, %d). Enter how much worse the new second objective is.", a1, a2, b1, b2), b2-a2,
				fmt.Sprintf("%d - %d = %d.", b2, a2, b2-a2))
		},
		func() trio {
			points1 := []int{limit(1, 5), limit(1, 5), limit(1, 5)}
			points2 := []int{points1[0] + limit(0, 4), points1[1] + limit(0, 4), points1[2] + limit(0, 4)}
			dominated := 0
			for i := range points1 {
				if points2[i] > points1[i] {
					dominated++
				}
			}
			return mk(fmt.Sprintf("Three points have first-objective values %d, %d and %d, and the compared set is uniformly larger. Enter how many of the three are dominated.", points1[0], points1[1], points1[2]), dominated,
				fmt.Sprintf("%d of the three are beaten by the uniformly larger set.", dominated))
		},
	},
	"optim.weighted_sum": {
		func() trio {
			w1, w2, f1, f2 := limit(1, 5), limit(1, 5), limit(1, 9), limit(1, 9)
			return mk(fmt.Sprintf("A weighted sum uses weights %d and %d on objectives %d and %d. Enter the weighted value.", w1, w2, f1, f2), w1*f1+w2*f2,
				fmt.Sprintf("%d·%d + %d·%d = %d.", w1, f1, w2, f2, w1*f1+w2*f2))
		},
		func() trio {
			w1, w2 := limit(1, 7), limit(1, 7)
			total := w1 + w2
			return mk(fmt.Sprintf("Weights %d and %d are normalized to sum to one. Enter the numerator of the first normalized weight when both are scaled by %d.", w1, w2, total), w1,
				fmt.Sprintf("Scaling by %d leaves %d/%d, whose numerator is %d.", total, w1, total, w1))
		},
		func() trio {
			a1, a2 := limit(1, 9), limit(1, 9)
			b1, b2 := limit(1, 9), limit(1, 9)
			w1, w2 := limit(1, 3), limit(1, 3)
			va, vb := w1*a1+w2*a2, w1*b1+w2*b2
			best := 0
			if vb < va {
				best = 1
			}
			return mk(fmt.Sprintf("With weights %d and %d, A scores %d and B scores %d. Enter 0 when A is better and 1 when B is better.", w1, w2, va, vb), best,
				fmt.Sprintf("%d versus %d selects %s.", va, vb, map[bool]string{true: "B", false: "A"}[best == 1]))
		},
	},
	"optim.epsilon_constraint": {
		func() trio {
			limitValue, actual := limit(4, 15), limit(2, 12)
			return mk(fmt.Sprintf("An epsilon constraint allows the second objective at most %d. The candidate reaches %d. Enter the slack.", limitValue, actual), limitValue-actual,
				fmt.Sprintf("%d - %d = %d; a nonnegative value is feasible.", limitValue, actual, limitValue-actual))
		},
		func() trio {
			a1, a2, eps := limit(1, 6), limit(1, 6), limit(1, 8)
			return mk(fmt.Sprintf("A problem minimizes the first objective with the second at most %d. At objectives (%d, %d), enter the first objective plus the second.", eps, a1, a2), a1+a2,
				fmt.Sprintf("The value is %d + %d = %d, and feasibility needs the second not to exceed %d.", a1, a2, a1+a2, eps))
		},
		func() trio {
			eps, step := limit(6, 20), limit(1, 3)
			return mk(fmt.Sprintf("Tightening epsilon from %d by %d per solve. After two solves, enter the new epsilon.", eps, step), eps-2*step,
				fmt.Sprintf("%d - 2·%d = %d.", eps, step, eps-2*step))
		},
	},
	"optim.goal_programming": {
		func() trio {
			target, reached := limit(10, 40), limit(1, 9)
			return mk(fmt.Sprintf("A goal is %d and the plan reaches %d. Enter the shortfall.", target, reached), target-reached,
				fmt.Sprintf("%d - %d = %d.", target, reached, target-reached))
		},
		func() trio {
			d1, d2, w1, w2 := limit(1, 9), limit(1, 9), limit(1, 4), limit(1, 4)
			return mk(fmt.Sprintf("Two deviations %d and %d are weighted by %d and %d. Enter the weighted deviation.", d1, d2, w1, w2), w1*d1+w2*d2,
				fmt.Sprintf("%d·%d + %d·%d = %d.", w1, d1, w2, d2, w1*d1+w2*d2))
		},
		func() trio {
			d1, d2, d3 := limit(1, 6), limit(1, 6), limit(1, 6)
			return mk(fmt.Sprintf("Three goals are missed by %d, %d and %d. Enter the total deviation.", d1, d2, d3), d1+d2+d3,
				fmt.Sprintf("%d + %d + %d = %d.", d1, d2, d3, d1+d2+d3))
		},
	},
	"optim.tradeoff_rate": {
		func() trio {
			df1, df2 := limit(1, 6), limit(1, 9)
			return mk(fmt.Sprintf("Increasing the first objective by %d allows the second to fall by %d. Enter the trade-off rate df2/df1, using integer division.", df1, df2), df2/df1,
				fmt.Sprintf("%d / %d = %d.", df2, df1, df2/df1))
		},
		func() trio {
			p1, p2, q1, q2 := limit(1, 5), limit(1, 9), limit(6, 9), limit(1, 4)
			return mk(fmt.Sprintf("A front moves from (%d, %d) to (%d, %d). Enter the change in the second objective.", p1, p2, q1, q2), q2-p2,
				fmt.Sprintf("%d - %d = %d; the sign says whether the second objective improved.", q2, p2, q2-p2))
		},
		func() trio {
			slope, dx := limit(2, 6), limit(1, 4)
			return mk(fmt.Sprintf("A trade-off slope of %d applies to a change of %d in the first objective. Enter the resulting change in the second.", slope, dx), slope*dx,
				fmt.Sprintf("%d·%d = %d.", slope, dx, slope*dx))
		},
	},
	"optim.kkt_multiobjective": {
		func() trio {
			w1, g1, w2, g2 := limit(1, 4), limit(1, 6), limit(1, 4), limit(1, 6)
			return mk(fmt.Sprintf("A weighted stationarity condition is w1 g1 + w2 g2. With weights %d and %d and gradients %d and %d, enter the residual.", w1, w2, g1, g2), w1*g1+w2*g2,
				fmt.Sprintf("%d·%d + %d·%d = %d.", w1, g1, w2, g2, w1*g1+w2*g2))
		},
		func() trio {
			w1, w2, w3 := limit(1, 4), limit(1, 4), limit(1, 4)
			return mk(fmt.Sprintf("Pareto stationarity uses nonnegative weights summing to one. With unnormalized weights %d, %d and %d, enter their sum.", w1, w2, w3), w1+w2+w3,
				fmt.Sprintf("%d + %d + %d = %d, which normalizes the three weights.", w1, w2, w3, w1+w2+w3))
		},
		func() trio {
			m1, g1, m2, g2 := limit(1, 5), limit(1, 6), limit(1, 5), limit(1, 6)
			return mk(fmt.Sprintf("Stationarity with multipliers m1 = %d and m2 = %d and gradients %d and %d gives m1 g1 + m2 g2. Enter the magnitude when both terms are positive.", m1, m2, g1, g2), m1*g1+m2*g2,
				fmt.Sprintf("%d·%d + %d·%d = %d.", m1, g1, m2, g2, m1*g1+m2*g2))
		},
	},
	"optim.pareto_optimality": {
		func() trio {
			w1, g1, w2, g2 := limit(1, 3), limit(1, 8), limit(1, 3), limit(1, 8)
			return mk(fmt.Sprintf("A point is Pareto optimal when some positive weights make the weighted gradient vanish. At weights %d and %d and gradients %d and %d, enter the weighted gradient.", w1, w2, g1, g2), w1*g1+w2*g2,
				fmt.Sprintf("%d·%d + %d·%d = %d, which must be zero at a supported optimum.", w1, g1, w2, g2, w1*g1+w2*g2))
		},
		func() trio {
			w1, w2, w3 := limit(1, 9), limit(1, 9), limit(1, 9)
			return mk(fmt.Sprintf("Supporting weights %d, %d and %d are all positive. Enter their product.", w1, w2, w3), w1*w2*w3,
				fmt.Sprintf("%d·%d·%d = %d, and positivity of every weight supports the point.", w1, w2, w3, w1*w2*w3))
		},
		func() trio {
			w1, w2 := limit(1, 9), limit(1, 9)
			return mk(fmt.Sprintf("Supporting weights %d and %d are both positive. Enter their product.", w1, w2), w1*w2,
				fmt.Sprintf("%d·%d = %d, and positivity of both weights supports the point.", w1, w2, w1*w2))
		},
	},
	"optim.nondominated_sort": {
		func() trio {
			first, second := limit(2, 9), limit(2, 9)
			return mk(fmt.Sprintf("The first front holds %d points and the second holds %d. Enter their total.", first, second), first+second,
				fmt.Sprintf("%d + %d = %d.", first, second, first+second))
		},
		func() trio {
			n := limit(2, 40)
			return mk(fmt.Sprintf("A set has %d points and every pair is compared once. Enter the number of comparisons.", n), n*(n-1)/2,
				fmt.Sprintf("%d·%d/2 = %d comparisons.", n, n-1, n*(n-1)/2))
		},
		func() trio {
			size, distance := limit(2, 9), limit(1, 9)
			return mk(fmt.Sprintf("%d points on a front each have crowding distance %d. Enter their sum.", size, distance), size*distance,
				fmt.Sprintf("%d·%d = %d.", size, distance, size*distance))
		},
	},
	"optim.crowding_distance": {
		func() trio {
			left, right := limit(1, 8), limit(9, 20)
			return mk(fmt.Sprintf("A point has neighbours at objective values %d and %d. Enter its crowding distance, the gap between them.", left, right), right-left,
				fmt.Sprintf("%d - %d = %d.", right, left, right-left))
		},
		func() trio {
			g1, g2, range1, range2 := limit(1, 6), limit(1, 6), limit(7, 12), limit(7, 12)
			return mk(fmt.Sprintf("Normalized gaps %d/%d and %d/%d are added. Enter the crowding distance scaled by %d, using integer division.", g1, range1, g2, range2, range1*range2), g1*range2+g2*range1,
				fmt.Sprintf("Scaling by %d clears both denominators, giving %d.", range1*range2, g1*range2+g2*range1))
		},
		func() trio {
			d1, d2, d3 := limit(1, 9), limit(1, 9), limit(1, 9)
			return mk(fmt.Sprintf("Three points on a front have crowding distances %d, %d and %d. Enter their average, using integer division.", d1, d2, d3), (d1+d2+d3)/3,
				fmt.Sprintf("(%d + %d + %d)/3 = %d.", d1, d2, d3, (d1+d2+d3)/3))
		},
	},
	// -----------------------------------------------------------------------
	// Large-scale optimization
	// -----------------------------------------------------------------------
	"optim.sparsity": {
		func() trio {
			n, zero := limit(6, 20), limit(1, 5)
			return mk(fmt.Sprintf("A vector of length %d has %d zero entries. Enter the number of nonzeros.", n, zero), n-zero,
				fmt.Sprintf("%d - %d = %d nonzeros.", n, zero, n-zero))
		},
		func() trio {
			a, b, c := limit(1, 5), -limit(1, 5), limit(1, 5)
			return mk(fmt.Sprintf("A sparse vector has entries %d, %d and %d. Enter its l1 norm.", a, b, c), absInt(a)+absInt(b)+absInt(c),
				fmt.Sprintf("%d + %d + %d = %d.", absInt(a), absInt(b), absInt(c), absInt(a)+absInt(b)+absInt(c)))
		},
		func() trio {
			n, nnz := limit(20, 90), limit(2, 9)
			return mk(fmt.Sprintf("A matrix has %d entries and %d nonzeros. Enter the percentage of entries that are zero, using integer division.", n, nnz), (n-nnz)*100/n,
				fmt.Sprintf("(%d - %d)·100/%d = %d percent.", n, nnz, n, (n-nnz)*100/n))
		},
	},
	"optim.first_order_method": {
		func() trio {
			x, g, t := limit(1, 9), limit(1, 9), limit(1, 3)
			return mk(fmt.Sprintf("A first-order step is x - t g with x = %d, g = %d and t = %d. Enter the new x.", x, g, t), x-t*g,
				fmt.Sprintf("%d - %d·%d = %d.", x, t, g, x-t*g))
		},
		func() trio {
			cond, factor := limit(2, 8), limit(2, 6)
			return mk(fmt.Sprintf("A first-order method needs about %d times the iterations of a second-order method because the condition number is %d. Enter the slowdown, using integer division.", factor-1, cond), (factor - 1),
				fmt.Sprintf("The condition number %d costs a factor of about %d.", cond, factor-1))
		},
		func() trio {
			initial, tol := limit(20, 90), limit(1, 4)
			start, steps := initial, 0
			for start > tol {
				start /= 2
				steps++
			}
			return mk(fmt.Sprintf("A first-order rate halves the error from %d until it is at most %d. Enter the number of iterations.", initial, tol), steps,
				fmt.Sprintf("From %d, %d halvings reach %d or less.", 2*start, steps+1, tol))
		},
	},
	"optim.dual_coordinate": {
		func() trio {
			y, a, x, c := limit(1, 6), limit(1, 5), limit(1, 8), limit(2, 9)
			return mk(fmt.Sprintf("A dual coordinate step moves y toward (c - a x)/a. With y = %d, a = %d, x = %d and c = %d, enter the integer target.", y, a, x, c), (c-a*x)/a,
				fmt.Sprintf("(%d - %d·%d)/%d = %d.", c, a, x, a, (c-a*x)/a))
		},
		func() trio {
			y1, y2 := limit(1, 9), limit(1, 9)
			return mk(fmt.Sprintf("Two dual coordinates change by %d and %d. Enter the l1 change.", y1, y2), y1+y2,
				fmt.Sprintf("%d + %d = %d.", y1, y2, y1+y2))
		},
		func() trio {
			n, chosen := limit(4, 20), limit(1, 3)
			return mk(fmt.Sprintf("A dual problem has %d coordinates and %d are updated per sweep. Enter the number of sweeps needed to touch every coordinate once, using integer division rounded up.", n, chosen), (n+chosen-1)/chosen,
				fmt.Sprintf("ceil(%d/%d) = %d sweeps.", n, chosen, (n+chosen-1)/chosen))
		},
	},
	"optim.decomposition": {
		func() trio {
			parts, size := limit(2, 6), limit(3, 9)
			return mk(fmt.Sprintf("A problem splits into %d subproblems of %d variables each. Enter the total number of subproblem variables.", parts, size), parts*size,
				fmt.Sprintf("%d·%d = %d.", parts, size, parts*size))
		},
		func() trio {
			total, coupled := limit(10, 40), limit(1, 6)
			return mk(fmt.Sprintf("A decomposition has %d variables and %d coupling variables. Enter the number that each subproblem solves locally.", total, coupled), total-coupled,
				fmt.Sprintf("%d - %d = %d.", total, coupled, total-coupled))
		},
		func() trio {
			s1, s2, s3 := limit(2, 9), limit(2, 9), limit(2, 9)
			return mk(fmt.Sprintf("Three subproblems report objective values %d, %d and %d. Enter their sum.", s1, s2, s3), s1+s2+s3,
				fmt.Sprintf("%d + %d + %d = %d.", s1, s2, s3, s1+s2+s3))
		},
	},
	"optim.admm": {
		func() trio {
			x1, x2 := limit(1, 9), limit(1, 9)
			return mk(fmt.Sprintf("An ADMM x-update averages two local values %d and %d. Enter the average, using integer division.", x1, x2), (x1+x2)/2,
				fmt.Sprintf("(%d + %d)/2 = %d.", x1, x2, (x1+x2)/2))
		},
		func() trio {
			y, rho, residual := limit(1, 6), limit(1, 5), limit(1, 8)
			return mk(fmt.Sprintf("The dual update is y_new = y + ρ r with y = %d, ρ = %d and residual r = %d. Enter y_new.", y, rho, residual), y+rho*residual,
				fmt.Sprintf("%d + %d·%d = %d.", y, rho, residual, y+rho*residual))
		},
		func() trio {
			primal, dual := limit(1, 9), limit(1, 9)
			return mk(fmt.Sprintf("A stopping rule tests primal residual %d plus dual residual %d against a tolerance. Enter their sum.", primal, dual), primal+dual,
				fmt.Sprintf("%d + %d = %d.", primal, dual, primal+dual))
		},
	},
	"optim.consensus": {
		func() trio {
			a, b, c := limit(1, 9), limit(1, 9), limit(1, 9)
			return mk(fmt.Sprintf("Three agents hold values %d, %d and %d. Enter their average, using integer division.", a, b, c), (a+b+c)/3,
				fmt.Sprintf("(%d + %d + %d)/3 = %d.", a, b, c, (a+b+c)/3))
		},
		func() trio {
			a, b := limit(1, 9), limit(1, 9)
			return mk(fmt.Sprintf("Two agents differ by values %d and %d. Enter their disagreement.", a, b), absInt(a-b),
				fmt.Sprintf("|%d - %d| = %d.", a, b, absInt(a-b)))
		},
		func() trio {
			initial, factor := limit(10, 60), limit(2, 4)
			start, steps := initial, 0
			for start > 1 {
				start /= factor
				steps++
			}
			return mk(fmt.Sprintf("Consensus contracts the disagreement from %d by the factor %d each round. Enter the number of rounds until it reaches 1, using integer division.", initial, factor), steps,
				fmt.Sprintf("%d shrinks to 1 in %d rounds.", initial, steps))
		},
	},
	"optim.stochastic_averaging": {
		func() trio {
			avg, sample, n := limit(1, 9), limit(1, 9), limit(2, 8)
			return mk(fmt.Sprintf("A running average over %d samples is %d. A new sample %d arrives. Enter the new average numerator, avg·%d + sample.", n, avg, sample, n), avg*n+sample,
				fmt.Sprintf("%d·%d + %d = %d.", avg, n, sample, avg*n+sample))
		},
		func() trio {
			n := limit(2, 40)
			return mk(fmt.Sprintf("A running average weights the newest sample by 1/%d. Enter that weight scaled by %d, using integer division.", n, n*n), n,
				fmt.Sprintf("Scaling 1/%d by %d gives %d.", n, n*n, n))
		},
		func() trio {
			a, b, c, d := limit(1, 6), limit(1, 6), limit(1, 6), limit(1, 6)
			return mk(fmt.Sprintf("Four samples %d, %d, %d and %d are averaged. Enter their sum.", a, b, c, d), a+b+c+d,
				fmt.Sprintf("%d + %d + %d + %d = %d.", a, b, c, d, a+b+c+d))
		},
	},
	"optim.subsampling": {
		func() trio {
			n, fraction := limit(40, 200), limit(2, 5)
			return mk(fmt.Sprintf("A dataset has %d rows and a subsample keeps one row in %d. Enter the subsample size, using integer division.", n, fraction), n/fraction,
				fmt.Sprintf("%d / %d = %d.", n, fraction, n/fraction))
		},
		func() trio {
			full, sub := limit(20, 60), limit(2, 9)
			return mk(fmt.Sprintf("A full pass costs %d units and a subsampled pass costs %d. Enter the saving.", full, sub), full-sub,
				fmt.Sprintf("%d - %d = %d.", full, sub, full-sub))
		},
		func() trio {
			n, k := limit(2, 6), limit(2, 5)
			return mk(fmt.Sprintf("A subsample of 1/%d of %d rows is repeated %d times. Enter the total rows processed.", k, n*n, k), n*n,
				fmt.Sprintf("%d rows per pass, times %d, covers %d rows.", n*n/k, k, n*n))
		},
	},
	"optim.matrix_free": {
		func() trio {
			n, cost := limit(10, 50), limit(2, 6)
			return mk(fmt.Sprintf("A matrix-vector product with a dense %d by %d matrix costs %d operations per entry. Enter the total cost.", n, n, cost), n*n*cost,
				fmt.Sprintf("%d·%d·%d = %d.", n, n, cost, n*n*cost))
		},
		func() trio {
			a, x := limit(2, 9), limit(1, 9)
			return mk(fmt.Sprintf("A Hessian-vector product for f(x) = a x² / 2 returns a x. With a = %d and x = %d, enter the product.", a, x), a*x,
				fmt.Sprintf("%d·%d = %d.", a, x, a*x))
		},
		func() trio {
			dense, structured := limit(20, 60), limit(2, 9)
			return mk(fmt.Sprintf("A dense representation stores %d entries and a structured one stores %d. Enter the memory saved.", dense, structured), dense-structured,
				fmt.Sprintf("%d - %d = %d.", dense, structured, dense-structured))
		},
	},
	"optim.preconditioning": {
		func() trio {
			before, after := limit(20, 90), limit(2, 9)
			return mk(fmt.Sprintf("A preconditioner lowers the condition number from %d to %d. Enter the improvement factor, using integer division.", before, after), before/after,
				fmt.Sprintf("%d / %d = %d.", before, after, before/after))
		},
		func() trio {
			diag, value := limit(2, 9), limit(1, 9)
			return mk(fmt.Sprintf("A diagonal preconditioner uses the entry %d. At value %d, enter the scaled value, using integer division.", diag, value), value/diag,
				fmt.Sprintf("%d / %d = %d.", value, diag, value/diag))
		},
		func() trio {
			iterations, factor := limit(20, 80), limit(2, 5)
			return mk(fmt.Sprintf("Preconditioning cuts the iteration count by the factor %d. From %d iterations, enter the new count, using integer division.", factor, iterations), iterations/factor,
				fmt.Sprintf("%d / %d = %d.", iterations, factor, iterations/factor))
		},
	},
}
