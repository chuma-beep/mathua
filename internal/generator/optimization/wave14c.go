package optimization

import (
	"fmt"
)

// wave14FamiliesC covers the numerical kernels a solver relies on.
var wave14FamiliesC = map[string][3]func() trio{
	"optim.kkt_solver": {
		func() trio {
			constraints := limit(2, 40)
			return mk(fmt.Sprintf("A KKT system has %d equality constraints. Enter the order of the Schur complement in the multipliers.", constraints), constraints,
				fmt.Sprintf("Eliminating the primal block leaves one equation per constraint, so the order is %d.", constraints))
		},
		func() trio {
			variables, constraints := limit(4, 12), limit(1, 4)
			return mk(fmt.Sprintf("A KKT matrix has %d primal variables and %d constraints. Enter its order.", variables, constraints), variables+constraints,
				fmt.Sprintf("%d + %d = %d.", variables, constraints, variables+constraints))
		},
		func() trio {
			negative, positive := limit(1, 4), limit(1, 6)
			return mk(fmt.Sprintf("A symmetric factorisation reports %d negative and %d positive eigenvalues. Enter the negative count.", negative, positive), negative,
				fmt.Sprintf("The factorisation reports %d negative eigenvalues directly.", negative))
		},
	},
	"optim.inertia_correction": {
		func() trio {
			negative, shift := limit(1, 4), limit(1, 20)
			return mk(fmt.Sprintf("A reduced Hessian has smallest eigenvalue -%d and the correction adds a shift of %d. Enter the shifted smallest eigenvalue.", negative, shift), shift-negative,
				fmt.Sprintf("-%d + %d = %d.", negative, shift, shift-negative))
		},
		func() trio {
			expected, found := limit(2, 20), limit(0, 3)
			return mk(fmt.Sprintf("A trust region expects %d negative eigenvalues and the factorisation reports %d. Enter the missing count.", expected, found), expected-found,
				fmt.Sprintf("%d - %d = %d.", expected, found, expected-found))
		},
		func() trio {
			attempts := limit(2, 40)
			return mk(fmt.Sprintf("An inertia correction refactorises %d times and the first attempt is counted. Enter the shifts tried after the first.", attempts), attempts-1,
				fmt.Sprintf("%d - 1 = %d further shifts.", attempts, attempts-1))
		},
	},
	"optim.regularized_newton": {
		func() trio {
			hessian, damping, step := limit(2, 6), limit(1, 4), limit(1, 5)
			return mk(fmt.Sprintf("A damped Newton step solves (h + μ) s = -g with h = %d, μ = %d and g = %d·(h + μ). Enter the step magnitude.", hessian, damping, step), step,
				fmt.Sprintf("s = %d·%d/%d = %d.", step, hessian+damping, hessian+damping, step))
		},
		func() trio {
			increase, rounds := limit(2, 20), limit(1, 3)
			value := 1
			for i := 0; i < rounds; i++ {
				value *= increase
			}
			return mk(fmt.Sprintf("A damping parameter of 1 is multiplied by %d on each of %d failed steps. Enter the final damping.", increase, rounds), value,
				fmt.Sprintf("%d^%d = %d.", increase, rounds, value))
		},
		func() trio {
			gradient, curvature := limit(2, 9), limit(1, 5)
			return mk(fmt.Sprintf("A Levenberg-Marquardt step is -g/(h + μ) with g = %d·(h + μ) and h + μ = %d. Enter the magnitude.", gradient, curvature), gradient,
				fmt.Sprintf("%d·%d/%d = %d.", gradient, curvature, curvature, gradient))
		},
	},
	"optim.iterative_refinement": {
		func() trio {
			digits, perRound, rounds := limit(12, 16), limit(2, 4), limit(1, 3)
			return mk(fmt.Sprintf("A solve has %d accurate digits and each refinement recovers %d. Enter the digits after %d rounds.", digits, perRound, rounds), digits+perRound*rounds,
				fmt.Sprintf("%d + %d·%d = %d.", digits, perRound, rounds, digits+perRound*rounds))
		},
		func() trio {
			residual, correction := limit(4, 12), limit(1, 5)
			return mk(fmt.Sprintf("A residual of %d is reduced by %d after one correction. Enter the remaining residual.", residual, correction), residual-correction,
				fmt.Sprintf("%d - %d = %d.", residual, correction, residual-correction))
		},
		func() trio {
			rounds, solves := limit(2, 20), limit(1, 3)
			return mk(fmt.Sprintf("Each of %d refinement rounds uses %d solves with the same factorisation. Enter the solves.", rounds, solves), rounds*solves,
				fmt.Sprintf("%d·%d = %d.", rounds, solves, rounds*solves))
		},
	},
	"optim.finite_difference": {
		func() trio {
			variables := limit(2, 40)
			return mk(fmt.Sprintf("A forward-difference gradient needs one extra evaluation per variable for %d variables. Enter the evaluations including the base point.", variables), variables+1,
				fmt.Sprintf("%d + 1 = %d evaluations.", variables, variables+1))
		},
		func() trio {
			variables := limit(2, 40)
			return mk(fmt.Sprintf("A central-difference gradient evaluates the function twice per variable for %d variables. Enter the evaluations.", variables), 2*variables,
				fmt.Sprintf("2·%d = %d.", variables, 2*variables))
		},
		func() trio {
			digits := limit(6, 40)
			return mk(fmt.Sprintf("The truncation and rounding errors of a forward difference balance near 10^-%d. Enter the exponent of the best step.", digits), digits,
				fmt.Sprintf("The balance sits near 10^-%d, so the exponent is %d.", digits, digits))
		},
	},
	"optim.complex_step": {
		func() trio {
			variables, directions := limit(2, 40), 1
			return mk(fmt.Sprintf("A complex-step gradient uses one evaluation per variable for %d variables. Enter the evaluations.", variables), variables*directions,
				fmt.Sprintf("%d·%d = %d evaluations.", variables, directions, variables*directions))
		},
		func() trio {
			digits := limit(10, 40)
			return mk(fmt.Sprintf("A complex step can be as small as 10^-%d before the second-order term dominates. Enter the exponent.", digits), digits,
				fmt.Sprintf("The usable exponent is %d.", digits))
		},
		func() trio {
			variables := limit(2, 40)
			return mk(fmt.Sprintf("A complex-step gradient evaluates %d variables once each. Enter the evaluations.", variables), variables,
				fmt.Sprintf("One evaluation per variable gives %d.", variables))
		},
	},
	"optim.autodiff": {
		func() trio {
			operations := limit(4, 40)
			return mk(fmt.Sprintf("A computation has %d elementary operations and the reverse pass costs about twice that. Enter the reverse operations.", operations), 2*operations,
				fmt.Sprintf("2·%d = %d.", operations, 2*operations))
		},
		func() trio {
			directions, operations := limit(2, 6), limit(4, 15)
			return mk(fmt.Sprintf("A forward-mode sweep computes %d directions over %d operations. Enter the operation count.", directions, operations), directions*operations,
				fmt.Sprintf("%d·%d = %d.", directions, operations, directions*operations))
		},
		func() trio {
			checkpoints, segments := limit(2, 6), limit(1, 4)
			return mk(fmt.Sprintf("Checkpointing splits a graph into %d segments with %d checkpoints each. Enter the stored states.", segments, checkpoints), segments*checkpoints,
				fmt.Sprintf("%d·%d = %d.", segments, checkpoints, segments*checkpoints))
		},
	},
	"optim.jfnk": {
		func() trio {
			outer, inner := limit(2, 6), limit(3, 8)
			return mk(fmt.Sprintf("A Newton-Krylov solve takes %d outer iterations with %d residual evaluations each. Enter the evaluations.", outer, inner), outer*inner,
				fmt.Sprintf("%d·%d = %d.", outer, inner, outer*inner))
		},
		func() trio {
			digits, forcing := limit(3, 30), limit(1, 3)
			return mk(fmt.Sprintf("A forcing term tightens the linear tolerance by %d digits over %d stages. Enter the digits tightened in total.", digits, forcing), digits*forcing,
				fmt.Sprintf("%d·%d = %d.", digits, forcing, digits*forcing))
		},
		func() trio {
			krylov, preconditioner := limit(4, 10), limit(2, 4)
			return mk(fmt.Sprintf("A Krylov solve needs %d iterations and a preconditioner cuts them by a factor of %d. Enter the new iteration count.", krylov, preconditioner), krylov/preconditioner,
				fmt.Sprintf("%d divided by %d is %d.", krylov, preconditioner, krylov/preconditioner))
		},
	},
	"optim.lbfgs": {
		func() trio {
			pairs, vectors := limit(2, 40), 2
			return mk(fmt.Sprintf("Limited memory keeps %d update pairs and each pair needs %d vectors. Enter the vectors stored.", pairs, vectors), pairs*vectors,
				fmt.Sprintf("%d·%d = %d.", pairs, vectors, pairs*vectors))
		},
		func() trio {
			pairs := limit(2, 40)
			return mk(fmt.Sprintf("A two-loop recursion over %d pairs performs two inner products per pair. Enter the inner products.", pairs), 2*pairs,
				fmt.Sprintf("2·%d = %d.", pairs, 2*pairs))
		},
		func() trio {
			pairs, dimension := limit(2, 6), limit(20, 90)
			return mk(fmt.Sprintf("A memory of %d pairs stores %d vectors instead of a dense matrix of order %d. Enter the stored vectors.", pairs, 2, dimension), 2*pairs,
				fmt.Sprintf("Two vectors per pair over %d pairs is %d, independent of the dimension %d.", pairs, 2*pairs, dimension))
		},
	},
	"optim.line_search_interp": {
		func() trio {
			slope, curvature, step := limit(2, 8), limit(1, 4), limit(1, 4)
			return mk(fmt.Sprintf("A quadratic model has slope -%d and curvature %d at the origin. Enter the numerator of the trial step %d/(2·%d).", slope*curvature*step, curvature, slope*curvature*step, curvature), slope,
				fmt.Sprintf("The minimiser is %d/(2·%d), so the numerator is %d.", slope*curvature*step, curvature, slope*curvature*step))
		},
		func() trio {
			trials := limit(1, 40)
			return mk(fmt.Sprintf("A line search evaluates the base point and then %d trial points. Enter the trial evaluations.", trials), trials,
				fmt.Sprintf("Each trial is one evaluation, so the count is %d.", trials))
		},
		func() trio {
			accepted, rejected := limit(2, 6), limit(1, 4)
			return mk(fmt.Sprintf("A line search accepts %d trials and rejects %d. Enter the evaluations after the first.", accepted, rejected), accepted+rejected-1,
				fmt.Sprintf("%d + %d - 1 = %d evaluations after the base point.", accepted, rejected, accepted+rejected-1))
		},
	},
}
