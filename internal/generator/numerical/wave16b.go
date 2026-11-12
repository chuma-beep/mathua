package numerical

// Wave 16, part two: numerical partial differential equations and nonlinear
// systems. The families follow the same integer-valued contract as part one.

import "fmt"

var wave16FamiliesB = map[string][3]func() w16Trio{
	"numa.pde.cfl": {
		func() w16Trio {
			c := w16Limit(1, 6)
			k := w16Limit(1, 8)
			dx := c * k
			return w16Mk(fmt.Sprintf("An explicit advection scheme runs at speed %d with spacing %d. What is the largest stable time step, in the same units?", c, dx), k,
				fmt.Sprintf("The Courant number must not exceed one, so the largest step is %d/%d = %d.", dx, c, k))
		},
		func() w16Trio {
			c := w16Limit(1, 6)
			dt := w16Limit(1, 5)
			return w16Mk(fmt.Sprintf("A signal travels at speed %d and the time step is %d. How far does it travel in one step?", c, dt), c*dt,
				fmt.Sprintf("Distance is speed times time, so %d · %d = %d.", c, dt, c*dt))
		},
		func() w16Trio {
			c := w16Limit(1, 5)
			dt := w16Limit(1, 4)
			k := w16Limit(2, 9)
			length := c * dt * k
			return w16Mk(fmt.Sprintf("A signal travels at speed %d with time step %d over a distance of %d. How many steps does it need?", c, dt, length), k,
				fmt.Sprintf("Each step covers %d · %d = %d units, so %d/%d = %d steps are needed.", c, dt, c*dt, length, c*dt, k))
		},
	},
	"numa.pde.explicit_heat": {
		func() w16Trio {
			dx := w16Limit(1, 3)
			k := w16Limit(1, 5)
			dt := w16Limit(1, 4)
			alpha := dx * dx * k
			return w16Mk(fmt.Sprintf("An explicit heat scheme has diffusivity %d, spacing %d and time step %d. What is the ratio r = alpha dt / dx squared?", alpha, dx, dt), k*dt,
				fmt.Sprintf("The spacing squared is %d and alpha dt is %d, so r = %d · %d = %d.", dx*dx, alpha*dt, k, dt, k*dt))
		},
		func() w16Trio {
			k := w16Limit(1, 5)
			m := w16Limit(1, 5)
			alpha := k
			dx := 2 * k * m
			return w16Mk(fmt.Sprintf("An explicit heat scheme has diffusivity %d and spacing %d. What is the largest stable time step, given dt <= dx squared / (2 alpha)?", alpha, dx), 2*k*m*m,
				fmt.Sprintf("The bound is (%d squared)/(2 · %d) = %d.", dx, alpha, 2*k*m*m))
		},
		func() w16Trio {
			dt := w16Limit(1, 5)
			k := w16Limit(2, 10)
			total := dt * k
			return w16Mk(fmt.Sprintf("A simulation runs to time %d with a time step of %d. How many steps does it take?", total, dt), k,
				fmt.Sprintf("The step count is %d/%d = %d.", total, dt, k))
		},
	},
	"numa.pde.implicit_heat": {
		func() w16Trio {
			n := w16Limit(5, 30)
			return w16Mk(fmt.Sprintf("An implicit heat solve uses %d nodes including both boundaries. How many interior unknowns does the linear system have?", n), n-2,
				fmt.Sprintf("The two boundary values are known, so %d - 2 = %d unknowns remain.", n, n-2))
		},
		func() w16Trio {
			n := w16Limit(5, 30)
			return w16Mk(fmt.Sprintf("The implicit heat matrix for %d nodes is tridiagonal. How many entries can be nonzero?", n), 3*n-2,
				fmt.Sprintf("The diagonal holds %d entries and the two off-diagonals hold %d, giving %d.", n, 2*(n-1), 3*n-2))
		},
		func() w16Trio {
			dt := w16Limit(1, 5)
			k := w16Limit(2, 12)
			total := dt * k
			return w16Mk(fmt.Sprintf("An implicit scheme advances to time %d with step %d. How many linear systems must be solved?", total, dt), k,
				fmt.Sprintf("One system is solved per step, so %d/%d = %d systems.", total, dt, k))
		},
	},
	"numa.pde.crank_nicolson": {
		func() w16Trio {
			n := w16Limit(10, 60)
			m := w16Limit(3, 12)
			return w16Mk(fmt.Sprintf("A Crank-Nicolson run takes %d steps with %d unknowns per step. How many unknown values are computed in total?", n, m), n*m,
				fmt.Sprintf("%d steps times %d unknowns gives %d values.", n, m, n*m))
		},
		func() w16Trio {
			r := w16Limit(2, 4)
			p := w16Limit(2, 3)
			factor := 1
			for i := 0; i < p; i++ {
				factor *= r
			}
			return w16Mk(fmt.Sprintf("A scheme of order %d in time refines the step by a factor of %d. By what factor does the time error fall?", p, r), factor,
				fmt.Sprintf("The error scales like the step to the power %d, so it falls by %d to the power %d = %d.", p, r, p, factor))
		},
		func() w16Trio {
			n := w16Limit(10, 80)
			return w16Mk(fmt.Sprintf("A Crank-Nicolson run takes %d steps. How many linear systems does it solve?", n), n,
				fmt.Sprintf("Each step solves one system, so %d steps solve %d systems.", n, n))
		},
	},
	"numa.pde.method_of_lines": {
		func() w16Trio {
			n := w16Limit(6, 40)
			return w16Mk(fmt.Sprintf("The method of lines uses %d spatial nodes including both boundaries. How many ordinary differential equations make up the semi-discrete system?", n), n-2,
				fmt.Sprintf("The two boundary nodes are fixed, leaving %d - 2 = %d equations.", n, n-2))
		},
		func() w16Trio {
			n := w16Limit(6, 40)
			return w16Mk(fmt.Sprintf("A second-difference discretisation on %d nodes produces a tridiagonal Jacobian. How many entries can be nonzero?", n), 3*n-2,
				fmt.Sprintf("The diagonal has %d entries and the off-diagonals have %d, giving %d.", n, 2*(n-1), 3*n-2))
		},
		func() w16Trio {
			alpha := w16Limit(1, 4)
			k := w16Limit(1, 4)
			dx := 2 * alpha * k
			return w16Mk(fmt.Sprintf("A diffusion term has coefficient %d and the grid spacing is %d. What is the largest explicit stable time step, dt <= dx squared / (2 alpha)?", alpha, dx), 2*alpha*k*k,
				fmt.Sprintf("The bound is (%d squared)/(2 · %d) = %d.", dx, alpha, 2*alpha*k*k))
		},
	},
	"numa.pde.upwind": {
		func() w16Trio {
			c := 2 * w16Limit(1, 5)
			dx := w16Limit(1, 5)
			courant := 20 * w16Limit(0, 4)
			return w16Mk(fmt.Sprintf("Upwind differencing runs at speed %d with spacing %d and Courant number %d percent. What is the artificial diffusion coefficient c dx (1 - C) / 2?", c, dx, courant),
				c*dx*(100-courant)/200,
				fmt.Sprintf("The coefficient is %d · %d · (%d)/200 = %d.", c, dx, 100-courant, c*dx*(100-courant)/200))
		},
		func() w16Trio {
			i := w16Limit(2, 20)
			sign := 1
			if w16Limit(0, 1) == 0 {
				sign = -1
			}
			return w16Mk(fmt.Sprintf("An upwind scheme at node %d has velocity %+d. Which neighbouring index i + s supplies the upwind value? Enter s.", i, sign), sign,
				fmt.Sprintf("Positive velocity takes the value from the left, so s = -1; negative velocity takes it from the right, so s = +1. Here s = %d.", sign))
		},
		func() w16Trio {
			c := w16Limit(1, 6)
			k := w16Limit(1, 8)
			dx := c * k
			return w16Mk(fmt.Sprintf("An upwind scheme runs at speed %d with spacing %d. What is the largest stable time step?", c, dx), k,
				fmt.Sprintf("The Courant bound gives %d/%d = %d.", dx, c, k))
		},
	},
	"numa.pde.finite_element": {
		func() w16Trio {
			e := w16Limit(2, 30)
			return w16Mk(fmt.Sprintf("A one-dimensional mesh has %d elements with linear basis functions. How many nodes does it have?", e), e+1,
				fmt.Sprintf("Linear elements share nodes, so %d elements give %d nodes.", e, e+1))
		},
		func() w16Trio {
			p := w16Limit(1, 6)
			return w16Mk(fmt.Sprintf("A finite element uses polynomials of degree %d on a triangle. How many basis functions does the element carry, and hence how many entries does its local matrix have?", p), (p+1)*(p+2)/2,
				fmt.Sprintf("A degree-%d triangle element has (%d)(%d)/2 = %d basis functions, so its local matrix has %d entries.", p, p+1, p+2, (p+1)*(p+2)/2, (p+1)*(p+2)/2))
		},
		func() w16Trio {
			e := w16Limit(2, 30)
			return w16Mk(fmt.Sprintf("A one-dimensional linear finite element mesh has %d elements. How many nonzeros does the assembled global matrix have?", e), 3*e+1,
				fmt.Sprintf("The mesh has %d nodes and a tridiagonal matrix, so it holds 3(%d) - 2 = %d nonzeros.", e+1, e+1, 3*(e+1)-2))
		},
	},
	"numa.pde.stiffness_matrix": {
		func() w16Trio {
			e := w16Limit(2, 20)
			l := w16Limit(2, 6)
			return w16Mk(fmt.Sprintf("A mesh has %d elements and each element matrix has %d by %d entries. How many entries are assembled in total, counting repeats?", e, l, l), e*l*l,
				fmt.Sprintf("%d elements times %d entries each gives %d.", e, l*l, e*l*l))
		},
		func() w16Trio {
			n := w16Limit(4, 30)
			return w16Mk(fmt.Sprintf("A one-dimensional stiffness matrix on %d nodes is tridiagonal. How many entries can be nonzero?", n), 3*n-2,
				fmt.Sprintf("The diagonal holds %d entries and the off-diagonals hold %d, giving %d.", n, 2*(n-1), 3*n-2))
		},
		func() w16Trio {
			r := w16Limit(1, 6)
			factor := 1
			for i := 0; i < r; i++ {
				factor *= 4
			}
			return w16Mk(fmt.Sprintf("The stiffness matrix condition number grows like h to the power minus two. The mesh is refined %d times, halving h each time. By what factor does the condition number grow?", r), factor,
				fmt.Sprintf("Halving h multiplies the condition number by four, so %d refinements multiply it by %d.", r, factor))
		},
	},
	"numa.pde.mesh": {
		func() w16Trio {
			e := w16Limit(2, 20)
			return w16Mk(fmt.Sprintf("A mesh has %d elements and every element is halved uniformly. How many elements does the refined mesh have?", e), 2*e,
				fmt.Sprintf("Each element becomes two, so %d · 2 = %d.", e, 2*e))
		},
		func() w16Trio {
			r := w16Limit(2, 4)
			p := w16Limit(2, 3)
			factor := 1
			for i := 0; i < p; i++ {
				factor *= r
			}
			return w16Mk(fmt.Sprintf("A method of order %d uses a mesh refined by a factor of %d. By what factor does the error fall?", p, r), factor,
				fmt.Sprintf("The error scales like h to the power %d, so the factor is %d to the power %d = %d.", p, r, p, factor))
		},
		func() w16Trio {
			e := w16Limit(2, 10)
			r := w16Limit(1, 4)
			total := e
			for i := 0; i < r; i++ {
				total *= 2
			}
			return w16Mk(fmt.Sprintf("A mesh starts with %d elements and is uniformly refined %d times, doubling the element count each time. How many elements does it end with?", e, r), total,
				fmt.Sprintf("%d · 2 to the power %d = %d.", e, r, total))
		},
	},
	"numa.pde.spectral": {
		func() w16Trio {
			n := w16Limit(4, 30)
			return w16Mk(fmt.Sprintf("A spectral method keeps %d Fourier modes. How many multiplications does a derivative cost in coefficient space?", n), n,
				fmt.Sprintf("The derivative multiplies each coefficient by its wave number, so it costs %d multiplications.", n))
		},
		func() w16Trio {
			levels := w16Limit(2, 10)
			n := 1 << levels
			return w16Mk(fmt.Sprintf("A spectral transform uses %d points, a power of two equal to 2^%d. The fast transform costs N log2 N. How many operations is that?", n, levels), n*levels,
				fmt.Sprintf("The cost is %d · %d = %d.", n, levels, n*levels))
		},
		func() w16Trio {
			p := w16Limit(1, 12)
			return w16Mk(fmt.Sprintf("A spectral basis uses polynomials of degree up to %d. How many coefficients are needed for one variable?", p), p+1,
				fmt.Sprintf("Degrees 0 through %d give %d coefficients.", p, p+1))
		},
	},
	"numa.nonlinear.newton_system": {
		func() w16Trio {
			n := w16Limit(2, 12)
			return w16Mk(fmt.Sprintf("Newton's method is applied to a system of %d equations in %d unknowns. How many entries does the Jacobian have?", n, n), n*n,
				fmt.Sprintf("The Jacobian is %d by %d, so it has %d entries.", n, n, n*n))
		},
		func() w16Trio {
			n := w16Limit(2, 8)
			return w16Mk(fmt.Sprintf("A dense Newton step solves a %d by %d linear system by elimination. How many operations does the solve cost, in the usual cubic count?", n, n), n*n*n,
				fmt.Sprintf("Dense elimination costs about n cubed, here %d.", n*n*n))
		},
		func() w16Trio {
			e := w16Limit(2, 9)
			return w16Mk(fmt.Sprintf("Near a simple root the Newton error satisfies e_next = C e squared with C = 1. If the current error is %d, what is the next error?", e), e*e,
				fmt.Sprintf("Squaring the error gives %d squared = %d.", e, e*e))
		},
	},
	"numa.nonlinear.jacobian_cost": {
		func() w16Trio {
			n := w16Limit(2, 20)
			return w16Mk(fmt.Sprintf("A forward-difference Jacobian is formed for %d unknowns. How many function evaluations does it need?", n), n+1,
				fmt.Sprintf("One evaluation for the base point plus one per unknown gives %d + 1 = %d.", n, n+1))
		},
		func() w16Trio {
			n := w16Limit(2, 20)
			return w16Mk(fmt.Sprintf("A central-difference Jacobian is formed for %d unknowns. How many function evaluations does it need?", n), 2*n,
				fmt.Sprintf("Each unknown needs one evaluation on each side, so 2 · %d = %d.", n, 2*n))
		},
		func() w16Trio {
			b := w16Limit(2, 6)
			return w16Mk(fmt.Sprintf("A Jacobian has bandwidth %d and its columns are coloured so that columns sharing a row are not perturbed together. How many colours are needed?", b), b,
				fmt.Sprintf("A band of width %d needs %d colours, one for each offset.", b, b))
		},
	},
	"numa.nonlinear.broyden": {
		func() w16Trio {
			n := w16Limit(2, 12)
			return w16Mk(fmt.Sprintf("Broyden's method maintains an approximate Jacobian for a system of %d unknowns. How many entries does the update store?", n), n*n,
				fmt.Sprintf("The approximation is a full %d by %d matrix, so it stores %d entries.", n, n, n*n))
		},
		func() w16Trio {
			n := w16Limit(2, 9)
			return w16Mk(fmt.Sprintf("A dense Newton step costs n cubed and a Broyden step costs n squared for a system of order %d. By what factor is one Broyden step cheaper?", n), n,
				fmt.Sprintf("The ratio is n cubed over n squared, which is n = %d.", n))
		},
		func() w16Trio {
			n := w16Limit(2, 15)
			return w16Mk(fmt.Sprintf("A Newton step evaluates the full Jacobian for %d unknowns, needing %d evaluations, while a Broyden step needs one. How many evaluations does Broyden save per step?", n, n+1), n,
				fmt.Sprintf("Newton needs %d + 1 evaluations and Broyden needs 1, so the saving is %d.", n, n))
		},
	},
	"numa.nonlinear.secant_condition": {
		func() w16Trio {
			s := w16Limit(1, 8)
			k := w16Limit(2, 12)
			y := s * k
			return w16Mk(fmt.Sprintf("A secant step changes x by %d and the residual by %d. What slope does the secant condition require?", s, y), k,
				fmt.Sprintf("The slope is %d/%d = %d.", y, s, k))
		},
		func() w16Trio {
			n := w16Limit(2, 10)
			return w16Mk(fmt.Sprintf("The secant condition is one scalar equation for an approximate Jacobian of order %d. How many entries of the matrix remain free?", n), n*n-1,
				fmt.Sprintf("The matrix has %d entries and the condition fixes one combination, leaving %d.", n*n, n*n-1))
		},
		func() w16Trio {
			n := w16Limit(2, 12)
			return w16Mk(fmt.Sprintf("A rank-one secant update is applied to an approximate Jacobian of order %d. How many entries does the update touch?", n), n*n,
				fmt.Sprintf("A rank-one update is a full %d by %d matrix, so it touches %d entries.", n, n, n*n))
		},
	},
	"numa.nonlinear.continuation": {
		func() w16Trio {
			h := w16Limit(1, 5)
			k := w16Limit(2, 12)
			length := h * k
			return w16Mk(fmt.Sprintf("A continuation runs over a parameter interval of length %d with step size %d. How many steps does it take?", length, h), k,
				fmt.Sprintf("The step count is %d/%d = %d.", length, h, k))
		},
		func() w16Trio {
			s := w16Limit(2, 12)
			return w16Mk(fmt.Sprintf("A continuation takes %d predictor-corrector steps, each costing two linear solves. How many linear solves does it perform?", s), 2*s,
				fmt.Sprintf("Two solves per step times %d steps gives %d.", s, 2*s))
		},
		func() w16Trio {
			s := w16Limit(2, 12)
			return w16Mk(fmt.Sprintf("A continuation takes %d steps and records the solution at every step and at the starting point. How many solution points does it record?", s), s+1,
				fmt.Sprintf("The start plus %d steps gives %d points.", s, s+1))
		},
	},
	"numa.nonlinear.trust_region": {
		func() w16Trio {
			r := w16Limit(1, 8)
			return w16Mk(fmt.Sprintf("A trust-region method has radius %d and the agreement ratio exceeds 0.75, so the radius doubles. What is the new radius?", r), 2*r,
				fmt.Sprintf("Doubling %d gives %d.", r, 2*r))
		},
		func() w16Trio {
			k := w16Limit(1, 6)
			r := 4 * k
			return w16Mk(fmt.Sprintf("A trust-region method has radius %d and the agreement ratio falls below 0.25, so the radius is quartered. What is the new radius?", r), k,
				fmt.Sprintf("A quarter of %d is %d.", r, k))
		},
		func() w16Trio {
			k := w16Limit(2, 4)
			actual := w16Limit(1, 10)
			predicted := actual * k
			return w16Mk(fmt.Sprintf("A trust-region step gives an actual reduction of %d and a predicted reduction of %d. What is the agreement ratio in percent?", actual, predicted), 100/k,
				fmt.Sprintf("The ratio is %d/%d = 1/%d, which is %d percent.", actual, predicted, k, 100/k))
		},
	},
	"numa.nonlinear.levenberg_marquardt": {
		func() w16Trio {
			n := w16Limit(2, 12)
			return w16Mk(fmt.Sprintf("A Levenberg-Marquardt system has %d unknowns and the damping adds a multiple of the identity. How many diagonal entries are shifted?", n), n,
				fmt.Sprintf("The identity is %d by %d, so it shifts %d diagonal entries.", n, n, n))
		},
		func() w16Trio {
			lambda := w16Limit(1, 20)
			return w16Mk(fmt.Sprintf("A Levenberg-Marquardt step fails and the damping parameter %d is multiplied by ten. What is the new value?", lambda), 10*lambda,
				fmt.Sprintf("Ten times %d is %d.", lambda, 10*lambda))
		},
		func() w16Trio {
			k := w16Limit(1, 20)
			lambda := 10 * k
			return w16Mk(fmt.Sprintf("A Levenberg-Marquardt step succeeds and the damping parameter %d is divided by ten. What is the new value?", lambda), k,
				fmt.Sprintf("Dividing %d by ten gives %d.", lambda, k))
		},
	},
	"numa.nonlinear.residual_norm": {
		func() w16Trio {
			k := w16Limit(1, 9)
			return w16Mk(fmt.Sprintf("A residual vector is (%d, %d). What is its Euclidean norm?", 3*k, 4*k), 5*k,
				fmt.Sprintf("The vector is %d times (3, 4), and the 3-4-5 triangle gives norm %d.", k, 5*k))
		},
		func() w16Trio {
			residual := w16Limit(1, 10)
			k := w16Limit(2, 5)
			rhs := residual * k
			return w16Mk(fmt.Sprintf("A residual has norm %d and the right-hand side has norm %d. What is the relative residual in percent?", residual, rhs), 100/k,
				fmt.Sprintf("The ratio is %d/%d = 1/%d, which is %d percent.", residual, rhs, k, 100/k))
		},
		func() w16Trio {
			e := w16Limit(1, 12)
			return w16Mk(fmt.Sprintf("A stopping test asks for a relative residual below 10 to the power minus %d. How many significant digits does that require?", e), e,
				fmt.Sprintf("A tolerance of 10^-%d fixes %d significant digits.", e, e))
		},
	},
	"numa.nonlinear.armijo": {
		func() w16Trio {
			alpha := w16Limit(1, 5)
			g := 5 * w16Limit(1, 6)
			return w16Mk(fmt.Sprintf("An Armijo test uses c = 0.2, step alpha = %d and directional derivative %d. What decrease does the condition require? (enter the value c alpha g)", alpha, g), alpha*(g/5),
				fmt.Sprintf("The requirement is 0.2 · %d · %d = %d.", alpha, g, alpha*(g/5)))
		},
		func() w16Trio {
			k := w16Limit(1, 8)
			return w16Mk(fmt.Sprintf("A backtracking line search halves the step %d times before the Armijo condition holds. How many function evaluations does it perform, counting the initial trial?", k), k+1,
				fmt.Sprintf("The initial trial plus %d backtracks gives %d evaluations.", k, k+1))
		},
		func() w16Trio {
			cPercent := []int{10, 20, 25, 50}[w16Limit(0, 3)]
			alpha := w16Limit(1, 8)
			return w16Mk(fmt.Sprintf("An Armijo test uses c = 0.%02d and accepts step alpha = %d. What percentage of the directional derivative does the required decrease represent?", cPercent, alpha), cPercent*alpha,
				fmt.Sprintf("The fraction is 0.%02d · %d = %d percent.", cPercent, alpha, cPercent*alpha))
		},
	},
	"numa.nonlinear.quadratic_convergence": {
		func() w16Trio {
			d := w16Limit(1, 8)
			return w16Mk(fmt.Sprintf("An iteration with quadratic convergence currently has %d correct decimal digits. How many digits does it have after one more step?", d), 2*d,
				fmt.Sprintf("Quadratic convergence doubles the digits, so %d becomes %d.", d, 2*d))
		},
		func() w16Trio {
			d := w16Limit(1, 4)
			k := w16Limit(1, 5)
			p := d << k
			return w16Mk(fmt.Sprintf("Quadratic convergence starts with %d correct digits. How many steps are needed to reach %d correct digits?", d, p), k,
				fmt.Sprintf("Each step doubles the digits: %d, %d, ... after %d steps reaches %d.", d, 2*d, k, p))
		},
		func() w16Trio {
			d := w16Limit(1, 5)
			k := w16Limit(1, 6)
			exp := d << k
			return w16Mk(fmt.Sprintf("Quadratic convergence has error 10 to the power minus %d. After %d steps, what is the exponent of the error? (enter a positive number)", d, k), exp,
				fmt.Sprintf("The exponent doubles each step, so after %d steps it is %d · 2^%d = %d.", k, d, k, exp))
		},
	},
}
