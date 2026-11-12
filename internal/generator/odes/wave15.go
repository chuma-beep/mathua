package odes

// Wave 15 extends differential equations from 66 to 116 concepts. Five new
// subdomains carry the material that follows the ODE survey: PDE theory,
// dynamical systems, asymptotics, stiff solvers and the calculus of variations.
// Every family is integer-valued so the numeric grader can check the answer,
// and every family draws at least two independent parameters so the variety
// floor holds at each ordinal.

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

type w15Trio struct {
	question    string
	answer      int
	explanation string
}

func w15Limit(low, high int) int { return low + rand.Intn(high-low+1) }

func w15Mk(question string, answer int, explanation string) w15Trio {
	return w15Trio{question: question, answer: answer, explanation: explanation}
}

type w15Gen struct {
	families [3]func() w15Trio
}

func (g *w15Gen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	out := g.families[kp-1]()
	return generator.Problem{
		Question:    out.question,
		Answer:      strconv.Itoa(out.answer),
		Explanation: out.explanation,
		KPOrdinal:   kp,
	}
}

func registerWave15(reg *generator.Registry) {
	for id, families := range wave15FamiliesA {
		reg.Register(id, &w15Gen{families: families})
	}
	for id, families := range wave15FamiliesB {
		reg.Register(id, &w15Gen{families: families})
	}
	for id, families := range wave15FamiliesC {
		reg.Register(id, &w15Gen{families: families})
	}
}

// wave15FamiliesA covers the theory of partial differential equations.
var wave15FamiliesA = map[string][3]func() w15Trio{
	"ode.pde.well_posed": {
		func() w15Trio {
			k, m := w15Limit(2, 7), w15Limit(1, 5)
			return w15Mk(fmt.Sprintf("A data error of size 10^-%d is amplified by a factor 10^%d. Enter the exponent of the resulting error.", k, m), m-k,
				fmt.Sprintf("Multiplying by 10^%d adds %d to the exponent, so the error exponent is %d - %d = %d.", m, m, m, k, m-k))
		},
		func() w15Trio {
			n, m := w15Limit(2, 6), w15Limit(0, 2)
			if m >= n {
				m = n - 1
			}
			return w15Mk(fmt.Sprintf("A second-order problem carries %d boundary conditions and the reduced first-order equation can satisfy %d of them. Enter how many boundary conditions the reduced equation cannot satisfy.", n, m), n-m,
				fmt.Sprintf("%d - %d = %d conditions are left unsatisfied.", n, m, n-m))
		},
		func() w15Trio {
			points, intervals := w15Limit(4, 20), w15Limit(2, 9)
			return w15Mk(fmt.Sprintf("A discretisation places %d grid points on each of %d intervals. Enter the total number of grid points.", points, intervals), points*intervals,
				fmt.Sprintf("%d points on each interval times %d intervals gives %d.", points, intervals, points*intervals))
		},
	},
	"ode.pde.energy_method": {
		func() w15Trio {
			scale, halvings := w15Limit(1, 8), w15Limit(1, 3)
			start := 8 * scale
			answer := start
			for i := 0; i < halvings; i++ {
				answer /= 2
			}
			return w15Mk(fmt.Sprintf("The energy of a solution is %d and it halves once per unit of time. Enter the energy after %d units.", start, halvings), answer,
				fmt.Sprintf("Halving %d times gives %d.", halvings, answer))
		},
		func() w15Trio {
			gradient, poincare := w15Limit(3, 12), w15Limit(2, 9)
			return w15Mk(fmt.Sprintf("The Dirichlet integral of the gradient is %d and the Poincare constant is %d. Enter the lower bound on the decay rate.", gradient, poincare), gradient*poincare,
				fmt.Sprintf("%d times %d is %d, so the energy decays at least this fast.", gradient, poincare, gradient*poincare))
		},
		func() w15Trio {
			flux, interior := w15Limit(5, 30), w15Limit(1, 4)
			return w15Mk(fmt.Sprintf("A boundary flux of %d leaves a region whose interior term is %d. Enter the net rate.", flux, interior), flux-interior,
				fmt.Sprintf("%d - %d = %d.", flux, interior, flux-interior))
		},
	},
	"ode.pde.maximum_principle": {
		func() w15Trio {
			bound, margin := w15Limit(2, 20), w15Limit(1, 9)
			return w15Mk(fmt.Sprintf("A comparison function has boundary maximum %d, and a constant %d is added. Enter the resulting bound.", bound, margin), bound+margin,
				fmt.Sprintf("%d + %d = %d.", bound, margin, bound+margin))
		},
		func() w15Trio {
			lower := w15Limit(1, 15)
			upper := lower + w15Limit(1, 9)
			return w15Mk(fmt.Sprintf("Two solutions are ordered with values %d and %d on the boundary. Enter the largest possible gap between them inside.", lower, upper), upper-lower,
				fmt.Sprintf("%d - %d = %d is the gap the comparison principle preserves.", upper, lower, upper-lower))
		},
		func() w15Trio {
			value := w15Limit(1, 20)
			return w15Mk(fmt.Sprintf("A constant solution has boundary value %d. Enter its value at an interior point.", value), value,
				fmt.Sprintf("A constant solution equals %d everywhere.", value))
		},
	},
	"ode.pde.heat_kernel": {
		func() w15Trio {
			t := w15Limit(1, 12)
			return w15Mk(fmt.Sprintf("The heat kernel at time %d has variance 2t. Enter its variance.", t), 2*t,
				fmt.Sprintf("2 times %d is %d.", t, 2*t))
		},
		func() w15Trio {
			width, t := w15Limit(2, 6), w15Limit(1, 5)
			return w15Mk(fmt.Sprintf("At time 1 the kernel width is %d. The width grows as the square root of time. Enter the squared width at time %d.", width, t), width*width*t,
				fmt.Sprintf("The squared width is %d squared times %d, which is %d.", width, t, width*width*t))
		},
		func() w15Trio {
			length, height := w15Limit(2, 6), w15Limit(1, 9)
			total := length * height
			return w15Mk(fmt.Sprintf("Initial heat of %d is spread over a length of %d. Enter the uniform height at long time.", total, length), height,
				fmt.Sprintf("%d divided by %d is %d.", total, length, height))
		},
	},
	"ode.pde.green_identity": {
		func() w15Trio {
			volume, boundary := w15Limit(6, 30), w15Limit(1, 5)
			return w15Mk(fmt.Sprintf("Green's first identity gives volume term %d and boundary term %d. Enter the gradient integral.", volume, boundary), volume-boundary,
				fmt.Sprintf("%d - %d = %d.", volume, boundary, volume-boundary))
		},
		func() w15Trio {
			u, v := w15Limit(3, 20), w15Limit(1, 3)
			return w15Mk(fmt.Sprintf("The second identity has boundary integrand u times the normal derivative of v minus v times the normal derivative of u. With the first product %d and the second %d, enter the integrand.", u, v), u-v,
				fmt.Sprintf("%d - %d = %d.", u, v, u-v))
		},
		func() w15Trio {
			divergence, volume := w15Limit(2, 12), w15Limit(2, 8)
			return w15Mk(fmt.Sprintf("A vector field has constant divergence %d over a region of volume %d. Enter the outward flux.", divergence, volume), divergence*volume,
				fmt.Sprintf("The divergence theorem gives %d times %d, which is %d.", divergence, volume, divergence*volume))
		},
	},
	"ode.pde.harmonic_mean": {
		func() w15Trio {
			samples, mean := w15Limit(4, 12), w15Limit(1, 9)
			total := samples * mean
			return w15Mk(fmt.Sprintf("A circle is sampled at %d equally spaced points whose values sum to %d. Enter the average.", samples, total), mean,
				fmt.Sprintf("%d divided by %d is %d.", total, samples, mean))
		},
		func() w15Trio {
			low := w15Limit(1, 19)
			high := low + w15Limit(1, 9)
			return w15Mk(fmt.Sprintf("A harmonic function has boundary minimum %d and boundary maximum %d. Enter the range.", low, high), high-low,
				fmt.Sprintf("%d - %d = %d.", high, low, high-low))
		},
		func() w15Trio {
			bound, extra := w15Limit(1, 20), w15Limit(1, 9)
			return w15Mk(fmt.Sprintf("A harmonic function is bounded by %d on the boundary and a margin of %d is allowed. Enter the bound.", bound, extra), bound+extra,
				fmt.Sprintf("%d + %d = %d.", bound, extra, bound+extra))
		},
	},
	"ode.pde.polar_separation": {
		func() w15Trio {
			n := w15Limit(2, 9)
			return w15Mk(fmt.Sprintf("The angular separation constant of a periodic solution is n squared. Enter the constant for n equal to %d.", n), n*n,
				fmt.Sprintf("%d squared is %d.", n, n*n))
		},
		func() w15Trio {
			radius, n := w15Limit(2, 3), w15Limit(2, 4)
			value := 1
			for i := 0; i < n; i++ {
				value *= radius
			}
			return w15Mk(fmt.Sprintf("A separated solution is r to the power n times cos(n theta). At r equal to %d, theta equal to 0 and n equal to %d, enter the value.", radius, n), value,
				fmt.Sprintf("%d to the power %d is %d.", radius, n, value))
		},
		func() w15Trio {
			n := w15Limit(1, 8)
			return w15Mk(fmt.Sprintf("The function cos(n theta) has n periods in a full turn. For n equal to %d, enter the number of zeros in one full turn.", n), 2*n,
				fmt.Sprintf("Each period contributes two zeros, so %d periods give %d.", n, 2*n))
		},
	},
	"ode.pde.spherical_harmonics": {
		func() w15Trio {
			degree := w15Limit(0, 9)
			return w15Mk(fmt.Sprintf("Enter the number of spherical harmonics of degree %d.", degree), 2*degree+1,
				fmt.Sprintf("The orders run from -%d to %d, which is %d values.", degree, degree, 2*degree+1))
		},
		func() w15Trio {
			degree := w15Limit(1, 9)
			return w15Mk(fmt.Sprintf("A degree %d multiplet has one harmonic for each allowed order. Enter the number of orders.", degree), 2*degree+1,
				fmt.Sprintf("There are %d orders from -%d to %d.", 2*degree+1, degree, degree))
		},
		func() w15Trio {
			top := w15Limit(0, 6)
			return w15Mk(fmt.Sprintf("The degrees 0 through %d each contribute their full multiplet. Enter the total number of harmonics.", top), (top+1)*(top+1),
				fmt.Sprintf("The sum of the odd counts is (%d + 1) squared, which is %d.", top, (top+1)*(top+1)))
		},
	},
	"ode.pde.weak_solution": {
		func() w15Trio {
			derivatives, moved := w15Limit(2, 4), w15Limit(1, 2)
			return w15Mk(fmt.Sprintf("A weak form transfers %d derivative from the solution to the test function. If the solution starts with %d derivatives, enter how many remain.", moved, derivatives), derivatives-moved,
				fmt.Sprintf("%d - %d = %d.", derivatives, moved, derivatives-moved))
		},
		func() w15Trio {
			coercive, data := w15Limit(2, 5), w15Limit(1, 9)
			norm := coercive * data
			return w15Mk(fmt.Sprintf("A coercive form has constant %d and the data norm is %d. Enter the energy bound.", coercive, norm), data,
				fmt.Sprintf("%d divided by %d is %d.", norm, coercive, data))
		},
		func() w15Trio {
			basis := w15Limit(3, 12)
			return w15Mk(fmt.Sprintf("A Galerkin basis has %d functions. Enter the number of entries in the stiffness matrix.", basis), basis*basis,
				fmt.Sprintf("The matrix is square, so it has %d squared, which is %d, entries.", basis, basis*basis))
		},
	},
	"ode.pde.characteristics_nd": {
		func() w15Trio {
			a, b := w15Limit(1, 9), w15Limit(1, 9)
			return w15Mk(fmt.Sprintf("A constant coefficient vector is (%d, %d). Enter the squared speed along a characteristic.", a, b), a*a+b*b,
				fmt.Sprintf("%d squared plus %d squared is %d.", a, b, a*a+b*b))
		},
		func() w15Trio {
			start, speed, t := w15Limit(1, 9), w15Limit(2, 7), w15Limit(1, 5)
			return w15Mk(fmt.Sprintf("Data start at x equal to %d and travel at speed %d along a characteristic. Enter the position at time %d.", start, speed, t), start+speed*t,
				fmt.Sprintf("%d + %d times %d is %d.", start, speed, t, start+speed*t))
		},
		func() w15Trio {
			a1, a2, a3 := w15Limit(1, 4), w15Limit(1, 4), w15Limit(1, 4)
			b1, b2, b3 := w15Limit(1, 4), w15Limit(1, 4), w15Limit(1, 4)
			return w15Mk(fmt.Sprintf("A surface has normal (%d, %d, %d) and the coefficient vector is (%d, %d, %d). Enter their dot product.", b1, b2, b3, a1, a2, a3), a1*b1+a2*b2+a3*b3,
				fmt.Sprintf("The products sum to %d.", a1*b1+a2*b2+a3*b3))
		},
	},
}

// wave15FamiliesB covers dynamical systems and asymptotics.
var wave15FamiliesB = map[string][3]func() w15Trio{
	"ode.dynamical.vector_field": {
		func() w15Trio {
			a, b := w15Limit(1, 6), w15Limit(1, 6)
			return w15Mk(fmt.Sprintf("The field is x' = %d x and y' = %d y. Enter the squared speed at the point (1, 1).", a, b), a*a+b*b,
				fmt.Sprintf("%d squared plus %d squared is %d.", a, b, a*a+b*b))
		},
		func() w15Trio {
			start, rate, t := w15Limit(1, 9), w15Limit(2, 6), w15Limit(1, 5)
			return w15Mk(fmt.Sprintf("The flow of x' = %d with x(0) = %d is linear in time. Enter x at time %d.", rate, start, t), start+rate*t,
				fmt.Sprintf("%d + %d times %d is %d.", start, rate, t, start+rate*t))
		},
		func() w15Trio {
			first, second := w15Limit(1, 9), w15Limit(1, 9)
			return w15Mk(fmt.Sprintf("A flow waits %d units and then %d more. Enter the total wait.", first, second), first+second,
				fmt.Sprintf("%d + %d = %d.", first, second, first+second))
		},
	},
	"ode.dynamical.nullcline": {
		func() w15Trio {
			a, b := w15Limit(1, 9), w15Limit(1, 9)
			return w15Mk(fmt.Sprintf("The nullclines are x = %d and y = %d. Enter the sum of the coordinates of the fixed point.", a, b), a+b,
				fmt.Sprintf("%d + %d = %d.", a, b, a+b))
		},
		func() w15Trio {
			a, b := w15Limit(1, 9), w15Limit(1, 9)
			return w15Mk(fmt.Sprintf("The nullclines are x = %d and y = %d. Enter the product of the coordinates of the fixed point.", a, b), a*b,
				fmt.Sprintf("%d times %d is %d.", a, b, a*b))
		},
		func() w15Trio {
			a, b := w15Limit(1, 9), w15Limit(1, 9)
			return w15Mk(fmt.Sprintf("A fixed point is at (%d, %d). Enter its squared distance from the origin.", a, b), a*a+b*b,
				fmt.Sprintf("%d squared plus %d squared is %d.", a, b, a*a+b*b))
		},
	},
	"ode.dynamical.linearisation": {
		func() w15Trio {
			a, d := w15Limit(1, 6), w15Limit(1, 6)
			return w15Mk(fmt.Sprintf("A Jacobian has diagonal entries %d and %d. Enter its trace.", a, d), a+d,
				fmt.Sprintf("%d + %d = %d.", a, d, a+d))
		},
		func() w15Trio {
			a, b, c, d := w15Limit(1, 6), w15Limit(1, 4), w15Limit(1, 4), w15Limit(1, 6)
			return w15Mk(fmt.Sprintf("A Jacobian is the matrix with entries %d and %d on the first row and %d and %d on the second. Enter its determinant.", a, b, c, d), a*d-b*c,
				fmt.Sprintf("%d times %d minus %d times %d is %d.", a, d, b, c, a*d-b*c))
		},
		func() w15Trio {
			trace := w15Limit(2, 12)
			return w15Mk(fmt.Sprintf("A 2 by 2 Jacobian has trace %d. Enter the sum of its eigenvalues.", trace), trace,
				fmt.Sprintf("The sum of the eigenvalues equals the trace, which is %d.", trace))
		},
	},
	"ode.dynamical.hartman_grobman": {
		func() w15Trio {
			stable := w15Limit(1, 6)
			return w15Mk(fmt.Sprintf("A hyperbolic fixed point has %d eigenvalues with negative real part. Enter the dimension of its stable manifold.", stable), stable,
				fmt.Sprintf("Each negative eigenvalue contributes one stable direction, so the dimension is %d.", stable))
		},
		func() w15Trio {
			unstable := w15Limit(1, 6)
			return w15Mk(fmt.Sprintf("A hyperbolic fixed point has %d eigenvalues with positive real part. Enter the dimension of its unstable manifold.", unstable), unstable,
				fmt.Sprintf("Each positive eigenvalue contributes one unstable direction, so the dimension is %d.", unstable))
		},
		func() w15Trio {
			zero := w15Limit(1, 6)
			return w15Mk(fmt.Sprintf("A Jacobian has %d zero eigenvalues among its spectrum. Enter the number of eigenvalues on the imaginary axis.", zero), zero,
				fmt.Sprintf("A zero eigenvalue lies on the imaginary axis, so the count is %d.", zero))
		},
	},
	"ode.dynamical.limit_cycle": {
		func() w15Trio {
			period, cycles := w15Limit(2, 6), w15Limit(2, 5)
			return w15Mk(fmt.Sprintf("A limit cycle has period %d. Enter the number of cycles completed in %d units.", period, period*cycles), cycles,
				fmt.Sprintf("%d divided by %d is %d.", period*cycles, period, cycles))
		},
		func() w15Trio {
			multiplier := w15Limit(1, 9)
			return w15Mk(fmt.Sprintf("A return map has multiplier %d tenths. Enter the percentage by which a perturbation contracts in one period.", multiplier), 10-multiplier,
				fmt.Sprintf("The contraction is 1 minus %d tenths, which is %d tenths, or %d percent.", multiplier, 10-multiplier, 10-multiplier))
		},
		func() w15Trio {
			period, count := w15Limit(2, 9), w15Limit(2, 6)
			return w15Mk(fmt.Sprintf("A limit cycle has period %d. Enter the time needed for %d complete cycles.", period, count), period*count,
				fmt.Sprintf("%d times %d is %d.", period, count, period*count))
		},
	},
	"ode.dynamical.poincare_bendixson": {
		func() w15Trio {
			inner := w15Limit(1, 9)
			outer := inner + w15Limit(1, 9)
			return w15Mk(fmt.Sprintf("A trapping annulus has inner radius %d and outer radius %d. Enter its width.", inner, outer), outer-inner,
				fmt.Sprintf("%d - %d = %d.", outer, inner, outer-inner))
		},
		func() w15Trio {
			turns, crossings := w15Limit(2, 9), w15Limit(1, 4)
			return w15Mk(fmt.Sprintf("A trajectory crosses a cross-section %d times per turn. Enter the crossings after %d turns.", crossings, turns), turns*crossings,
				fmt.Sprintf("%d turns times %d crossings is %d.", turns, crossings, turns*crossings))
		},
		func() w15Trio {
			fixed, orbits := w15Limit(1, 9), w15Limit(1, 9)
			return w15Mk(fmt.Sprintf("A trapping region contains %d fixed points and %d closed orbits. Enter the total number of such objects.", fixed, orbits), fixed+orbits,
				fmt.Sprintf("%d + %d = %d.", fixed, orbits, fixed+orbits))
		},
	},
	"ode.dynamical.lyapunov_exponent": {
		func() w15Trio {
			spans, perSpan := w15Limit(2, 6), w15Limit(2, 5)
			return w15Mk(fmt.Sprintf("A separation doubles every %d units. Enter the number of doublings in %d units.", spans, spans*perSpan), perSpan,
				fmt.Sprintf("%d divided by %d is %d.", spans*perSpan, spans, perSpan))
		},
		func() w15Trio {
			start, gain := w15Limit(1, 9), w15Limit(1, 9)
			return w15Mk(fmt.Sprintf("A separation of 10^-%d grows to 10^%d. Enter the exponent gain.", start, gain), start+gain,
				fmt.Sprintf("The exponent moves from -%d to %d, a gain of %d.", start, gain, start+gain))
		},
		func() w15Trio {
			divergence, first, second := w15Limit(3, 12), w15Limit(1, 3), w15Limit(1, 3)
			return w15Mk(fmt.Sprintf("Three exponents sum to the divergence %d. The first two are %d and %d. Enter the third.", divergence, first, second), divergence-first-second,
				fmt.Sprintf("%d - %d - %d = %d.", divergence, first, second, divergence-first-second))
		},
	},
	"ode.dynamical.invariant_manifold": {
		func() w15Trio {
			total := w15Limit(3, 6)
			stable := w15Limit(1, total-1)
			return w15Mk(fmt.Sprintf("A %d-dimensional system has %d negative eigenvalues. Enter the stable manifold dimension.", total, stable), stable,
				fmt.Sprintf("The stable dimension equals the negative eigenvalue count, which is %d.", stable))
		},
		func() w15Trio {
			unstable, stable := w15Limit(1, 4), w15Limit(1, 4)
			return w15Mk(fmt.Sprintf("A hyperbolic point has unstable dimension %d and stable dimension %d. Enter the phase-space dimension.", unstable, stable), unstable+stable,
				fmt.Sprintf("%d + %d = %d.", unstable, stable, unstable+stable))
		},
		func() w15Trio {
			total := w15Limit(3, 8)
			stable := w15Limit(1, total-1)
			return w15Mk(fmt.Sprintf("A stable manifold has dimension %d inside a space of dimension %d. Enter its codimension.", stable, total), total-stable,
				fmt.Sprintf("%d - %d = %d.", total, stable, total-stable))
		},
	},
	"ode.dynamical.hamiltonian": {
		func() w15Trio {
			p, q := w15Limit(1, 6), w15Limit(1, 6)
			return w15Mk(fmt.Sprintf("The Hamiltonian is p squared plus q squared. Enter its value at p equal to %d and q equal to %d.", p, q), p*p+q*q,
				fmt.Sprintf("%d squared plus %d squared is %d.", p, q, p*p+q*q))
		},
		func() w15Trio {
			a, b := w15Limit(1, 9), w15Limit(1, 9)
			return w15Mk(fmt.Sprintf("The Hamiltonian is %d q squared plus %d p squared. Enter its value at q equal to 1 and p equal to 1.", a, b), a+b,
				fmt.Sprintf("%d + %d = %d.", a, b, a+b))
		},
		func() w15Trio {
			degrees := w15Limit(1, 5)
			return w15Mk(fmt.Sprintf("A Hamiltonian system has %d degrees of freedom. Enter the number of first-order equations.", degrees), 2*degrees,
				fmt.Sprintf("Each degree of freedom contributes one position and one momentum equation, so 2 times %d is %d.", degrees, 2*degrees))
		},
	},
	"ode.dynamical.gradient_flow": {
		func() w15Trio {
			x := w15Limit(1, 9)
			return w15Mk(fmt.Sprintf("The potential is x squared and the flow is x' = -2x. Enter the magnitude of the rate of decrease of the potential at x equal to %d.", x), 2*x*x,
				fmt.Sprintf("The derivative of the potential along the flow is -2x squared, whose magnitude at %d is %d.", x, 2*x*x))
		},
		func() w15Trio {
			critical := w15Limit(1, 9)
			return w15Mk(fmt.Sprintf("A potential has %d critical points. Enter the number of fixed points of its gradient flow.", critical), critical,
				fmt.Sprintf("Fixed points of the flow are exactly the critical points, so there are %d.", critical))
		},
		func() w15Trio {
			curvature := w15Limit(1, 9)
			return w15Mk(fmt.Sprintf("The potential is %d x squared. Enter the decay rate of the linearised flow at the origin.", curvature), 2*curvature,
				fmt.Sprintf("Differentiating %d x squared gives 2 times %d, which is %d.", curvature, curvature, 2*curvature))
		},
	},
	"ode.asymptotic.regular": {
		func() w15Trio {
			leading, slope := w15Limit(2, 9), w15Limit(1, 4)
			denominator := w15Limit(2, 3)
			correction := slope * denominator
			return w15Mk(fmt.Sprintf("An expansion is %d plus epsilon times %d. At epsilon equal to 1 over %d, enter the value as an integer.", leading, correction, denominator), leading+slope,
				fmt.Sprintf("The correction contributes %d over %d, which is %d, so the value is %d.", correction, denominator, slope, leading+slope))
		},
		func() w15Trio {
			order := w15Limit(1, 6)
			return w15Mk(fmt.Sprintf("A perturbation series terminates at order %d. Enter the number of non-zero orders including the leading term.", order), order+1,
				fmt.Sprintf("Orders 0 through %d give %d terms.", order, order+1))
		},
		func() w15Trio {
			order := w15Limit(1, 8)
			return w15Mk(fmt.Sprintf("The first correction in a regular expansion is of order epsilon to the power %d. Enter that power.", order), order,
				fmt.Sprintf("The correction is epsilon to the power %d.", order))
		},
	},
	"ode.asymptotic.singular": {
		func() w15Trio {
			original := w15Limit(2, 5)
			reduced := w15Limit(1, original-1)
			return w15Mk(fmt.Sprintf("A singularly perturbed equation has order %d and its reduced equation has order %d. Enter the number of lost derivatives.", original, reduced), original-reduced,
				fmt.Sprintf("%d - %d = %d.", original, reduced, original-reduced))
		},
		func() w15Trio {
			order := w15Limit(2, 8)
			return w15Mk(fmt.Sprintf("A singular perturbation multiplies the derivative of order %d by epsilon. Enter the order of the reduced equation.", order), order-1,
				fmt.Sprintf("%d - 1 = %d.", order, order-1))
		},
		func() w15Trio {
			original := w15Limit(2, 6)
			reduced := w15Limit(1, original-1)
			return w15Mk(fmt.Sprintf("An equation of order %d reduces to order %d in the singular limit. Enter how many boundary conditions it can still satisfy.", original, reduced), reduced,
				fmt.Sprintf("The reduced equation has order %d, so it can satisfy %d boundary conditions.", reduced, reduced))
		},
	},
	"ode.asymptotic.boundary_layer": {
		func() w15Trio {
			exponent, inverse := w15Limit(1, 3), w15Limit(2, 4)
			value := 1
			for i := 0; i < exponent; i++ {
				value *= inverse
			}
			return w15Mk(fmt.Sprintf("A layer has width epsilon to the power %d. For epsilon equal to 1 over %d, enter the width denominator.", exponent, inverse), value,
				fmt.Sprintf("The width is 1 over %d to the power %d, so the denominator is %d.", inverse, exponent, value))
		},
		func() w15Trio {
			inverse := w15Limit(2, 9)
			return w15Mk(fmt.Sprintf("A boundary layer has width epsilon and epsilon equals 1 over %d. Enter how many layer widths fit in an interval of unit length.", inverse), inverse,
				fmt.Sprintf("A unit interval holds %d widths of size 1 over %d.", inverse, inverse))
		},
		func() w15Trio {
			numerator := w15Limit(1, 9)
			inverse := w15Limit(2, 4)
			return w15Mk(fmt.Sprintf("The inner variable is x divided by epsilon. At x equal to %d over %d and epsilon equal to 1 over %d, enter the inner variable.", numerator, inverse, inverse), numerator,
				fmt.Sprintf("Dividing %d over %d by 1 over %d gives %d.", numerator, inverse, inverse, numerator))
		},
	},
	"ode.asymptotic.matched": {
		func() w15Trio {
			outer, inner := w15Limit(1, 9), w15Limit(1, 9)
			common := w15Limit(1, 3)
			return w15Mk(fmt.Sprintf("A composite expansion adds outer %d and inner %d and subtracts their common limit %d. Enter the composite value.", outer, inner, common), outer+inner-common,
				fmt.Sprintf("%d + %d - %d = %d.", outer, inner, common, outer+inner-common))
		},
		func() w15Trio {
			order := w15Limit(1, 6)
			return w15Mk(fmt.Sprintf("Two expansions are matched through order epsilon to the power %d. Enter the number of matched orders including the constant term.", order), order+1,
				fmt.Sprintf("Orders 0 through %d give %d matched terms.", order, order+1))
		},
		func() w15Trio {
			value := w15Limit(1, 9)
			return w15Mk(fmt.Sprintf("Matching sets the inner constant equal to the outer limit %d. Enter the inner constant.", value), value,
				fmt.Sprintf("The matching condition fixes the inner constant at %d.", value))
		},
	},
	"ode.asymptotic.multiple_scales": {
		func() w15Trio {
			slow := w15Limit(1, 4)
			inverse := w15Limit(2, 4)
			t := slow * inverse
			return w15Mk(fmt.Sprintf("The slow time is epsilon times t. At t equal to %d and epsilon equal to 1 over %d, enter the slow time as an integer.", t, inverse), slow,
				fmt.Sprintf("%d divided by %d is %d.", t, inverse, slow))
		},
		func() w15Trio {
			shift := w15Limit(1, 4)
			inverse := w15Limit(2, 3)
			coefficient := shift * inverse
			return w15Mk(fmt.Sprintf("The frequency correction is epsilon times %d. At epsilon equal to 1 over %d, enter the correction as an integer.", coefficient, inverse), shift,
				fmt.Sprintf("%d divided by %d is %d.", coefficient, inverse, shift))
		},
		func() w15Trio {
			rate, horizon := w15Limit(1, 4), w15Limit(1, 4)
			return w15Mk(fmt.Sprintf("An amplitude grows at rate %d per unit of slow time. Enter the exponent after %d units.", 2*rate, horizon), 2*rate*horizon,
				fmt.Sprintf("The amplitude equation gives an exponent of %d times %d, which is %d.", 2*rate, horizon, 2*rate*horizon))
		},
	},
	"ode.asymptotic.wkb": {
		func() w15Trio {
			wave, length := w15Limit(2, 9), w15Limit(2, 9)
			return w15Mk(fmt.Sprintf("The local wave number is %d and the interval has length %d. Enter the accumulated phase.", wave, length), wave*length,
				fmt.Sprintf("%d times %d is %d.", wave, length, wave*length))
		},
		func() w15Trio {
			power := w15Limit(1, 6)
			factor := 1
			for i := 0; i < power; i++ {
				factor *= 2
			}
			return w15Mk(fmt.Sprintf("The wave number increases by a factor of 4 to the power %d. Enter the factor by which the amplitude falls.", power), factor,
				fmt.Sprintf("The amplitude scales as the inverse square root, so the fall is 2 to the power %d, which is %d.", power, factor))
		},
		func() w15Trio {
			wave := w15Limit(2, 9)
			return w15Mk(fmt.Sprintf("The eikonal equation gives the square of the phase derivative as %d. Enter the phase derivative.", wave*wave), wave,
				fmt.Sprintf("The square root of %d squared is %d.", wave, wave))
		},
	},
	"ode.asymptotic.dominant_balance": {
		func() w15Trio {
			low := w15Limit(1, 4)
			high := low + w15Limit(1, 5)
			return w15Mk(fmt.Sprintf("The orders of two terms are epsilon to the power %d and epsilon to the power %d. Enter the numerator of the balance exponent.", low, high), high-low,
				fmt.Sprintf("The exponent numerator is %d - %d, which is %d.", high, low, high-low))
		},
		func() w15Trio {
			quotient := w15Limit(1, 4)
			denominator := w15Limit(2, 3)
			return w15Mk(fmt.Sprintf("A balance gives an exponent of %d over %d. Enter the exponent as an integer.", quotient*denominator, denominator), quotient,
				fmt.Sprintf("%d divided by %d is %d.", quotient*denominator, denominator, quotient))
		},
		func() w15Trio {
			base, exponent := w15Limit(2, 5), w15Limit(2, 3)
			value := 1
			for i := 0; i < exponent; i++ {
				value *= base
			}
			return w15Mk(fmt.Sprintf("Two terms are %d times epsilon to the power %d and epsilon to the power 0. Enter the scale of their balance.", value, exponent), base,
				fmt.Sprintf("The scale is the %d-th root of %d, which is %d.", exponent, value, base))
		},
	},
	"ode.asymptotic.lindstedt": {
		func() w15Trio {
			shift := w15Limit(1, 5)
			inverse := w15Limit(2, 4)
			return w15Mk(fmt.Sprintf("The frequency correction is epsilon times %d. At epsilon equal to 1 over %d, enter the frequency increase as an integer.", shift, inverse), shift,
				fmt.Sprintf("%d divided by %d is %d.", shift, inverse, shift))
		},
		func() w15Trio {
			quotient := w15Limit(1, 5)
			inverse := w15Limit(2, 4)
			coefficient := quotient * inverse
			return w15Mk(fmt.Sprintf("The period increase over 2 pi is %d times epsilon. At epsilon equal to 1 over %d, enter the increase as an integer.", coefficient, inverse), quotient,
				fmt.Sprintf("%d divided by %d is %d.", coefficient, inverse, quotient))
		},
		func() w15Trio {
			amplitude := w15Limit(1, 5)
			return w15Mk(fmt.Sprintf("The frequency shift is three eighths times epsilon times the square of the amplitude. At epsilon equal to 1 and amplitude equal to %d, enter eight thirds of the shift.", amplitude), amplitude*amplitude,
				fmt.Sprintf("Eight thirds of the shift leaves %d squared, which is %d.", amplitude, amplitude*amplitude))
		},
	},
	"ode.asymptotic.averaging": {
		func() w15Trio {
			offset, amplitude := w15Limit(1, 9), w15Limit(1, 9)
			return w15Mk(fmt.Sprintf("The field is %d plus %d cos(t). Enter its average over one period.", offset, amplitude), offset,
				fmt.Sprintf("The cosine averages to zero, so the average is %d.", offset))
		},
		func() w15Trio {
			samples, mean := w15Limit(2, 9), w15Limit(1, 9)
			total := samples * mean
			return w15Mk(fmt.Sprintf("Averaging %d samples gives a total of %d. Enter the mean.", samples, total), mean,
				fmt.Sprintf("%d divided by %d is %d.", total, samples, mean))
		},
		func() w15Trio {
			periods := w15Limit(1, 9)
			return w15Mk(fmt.Sprintf("The average of cos squared over one period is one half. Enter twice that average times %d periods.", periods), periods,
				fmt.Sprintf("Two times one half times %d is %d.", periods, periods))
		},
	},
	"ode.asymptotic.adiabatic": {
		func() w15Trio {
			frequency, action := w15Limit(2, 6), w15Limit(1, 9)
			energy := frequency * action
			return w15Mk(fmt.Sprintf("The action is energy divided by frequency. With energy %d and frequency %d, enter the action.", energy, frequency), action,
				fmt.Sprintf("%d divided by %d is %d.", energy, frequency, action))
		},
		func() w15Trio {
			factor := w15Limit(2, 8)
			return w15Mk(fmt.Sprintf("The frequency doubles %d times at fixed action. Enter the factor by which the energy grows.", factor-1), 1<<(factor-1),
				fmt.Sprintf("Each doubling of frequency doubles the energy, so %d doublings give a factor of %d.", factor-1, 1<<(factor-1)))
		},
		func() w15Trio {
			factor := w15Limit(2, 9)
			return w15Mk(fmt.Sprintf("At fixed action the frequency falls by a factor of %d. Enter the factor by which the energy falls.", factor), factor,
				fmt.Sprintf("Energy is action times frequency, so it falls by the same factor %d.", factor))
		},
	},
}
