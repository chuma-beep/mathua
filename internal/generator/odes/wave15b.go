package odes

// Wave 15, second half: stiff solvers and the calculus of variations.

import "fmt"

// wave15FamiliesC covers stiff numerical methods and variational problems.
var wave15FamiliesC = map[string][3]func() w15Trio{
	"ode.stiff.ratio": {
		func() w15Trio {
			slow, ratio := w15Limit(1, 5), w15Limit(2, 9)
			return w15Mk(fmt.Sprintf("The fastest rate is %d and the slowest is %d. Enter the stiffness ratio.", slow*ratio, slow), ratio,
				fmt.Sprintf("%d divided by %d is %d.", slow*ratio, slow, ratio))
		},
		func() w15Trio {
			short, long := w15Limit(2, 5), w15Limit(2, 5)
			return w15Mk(fmt.Sprintf("The shortest time scale is 1 over %d and the longest is %d. Enter their ratio.", short, long), short*long,
				fmt.Sprintf("The ratio is %d divided by 1 over %d, which is %d.", long, short, short*long))
		},
		func() w15Trio {
			steps := w15Limit(2, 20)
			return w15Mk(fmt.Sprintf("An explicit method takes %d fast steps for every slow step. Enter the number of steps that the slow time scale alone would allow.", steps), steps-1,
				fmt.Sprintf("%d - 1 = %d wasted steps.", steps, steps-1))
		},
	},
	"ode.stiff.a_stability": {
		func() w15Trio {
			rate := w15Limit(1, 9)
			return w15Mk(fmt.Sprintf("Backward Euler has stability function 1 over (1 - z). At z equal to -%d, enter the denominator.", rate), 1+rate,
				fmt.Sprintf("1 - (-%d) = %d.", rate, 1+rate))
		},
		func() w15Trio {
			order := w15Limit(1, 5)
			return w15Mk(fmt.Sprintf("An A-stable method has stability function 1 over (1 - z) to the power %d. At z equal to -1, enter the denominator.", order), 1<<order,
				fmt.Sprintf("The denominator is 2 to the power %d, which is %d.", order, 1<<order))
		},
		func() w15Trio {
			rate := w15Limit(1, 9)
			return w15Mk(fmt.Sprintf("The test parameter is z equal to -%d. Enter the real part of z as a positive number.", rate), rate,
				fmt.Sprintf("The real part is -%d, whose magnitude is %d.", rate, rate))
		},
	},
	"ode.stiff.l_stability": {
		func() w15Trio {
			rate := w15Limit(2, 9)
			return w15Mk(fmt.Sprintf("Backward Euler at z equal to -%d gives a factor of 1 over %d. Enter the denominator.", rate, 1+rate), 1+rate,
				fmt.Sprintf("1 - (-%d) = %d.", rate, 1+rate))
		},
		func() w15Trio {
			steps := w15Limit(1, 5)
			return w15Mk(fmt.Sprintf("At z equal to -1 the factor is one half, and %d steps multiply the factors. Enter the denominator of the result.", steps), 1<<steps,
				fmt.Sprintf("One half to the power %d has denominator 2 to the power %d, which is %d.", steps, steps, 1<<steps))
		},
		func() w15Trio {
			rate := w15Limit(2, 6)
			return w15Mk(fmt.Sprintf("The trapezoidal rule gives (2 - m) over (2 + m) at z equal to -m. For m equal to %d, enter the magnitude of the numerator.", rate), rate-2,
				fmt.Sprintf("The numerator is 2 - %d, whose magnitude is %d.", rate, rate-2))
		},
	},
	"ode.stiff.trapezoidal": {
		func() w15Trio {
			rate := w15Limit(2, 9)
			return w15Mk(fmt.Sprintf("The trapezoidal rule at z equal to -2k gives a factor whose magnitude is (k - 1) over (k + 1). For k equal to %d, enter the numerator.", rate), rate-1,
				fmt.Sprintf("The numerator is %d - 1, which is %d.", rate, rate-1))
		},
		func() w15Trio {
			rate := w15Limit(2, 9)
			return w15Mk(fmt.Sprintf("At z equal to -2k the trapezoidal rule gives a factor whose magnitude is (k - 1) over (k + 1). For k equal to %d, enter the denominator.", rate), rate+1,
				fmt.Sprintf("The denominator is %d + 1, which is %d.", rate, rate+1))
		},
		func() w15Trio {
			length := w15Limit(1, 9)
			return w15Mk(fmt.Sprintf("The trapezoidal rule is applied over %d subintervals of unit length. Enter the total length.", length), length,
				fmt.Sprintf("%d unit subintervals cover a length of %d.", length, length))
		},
	},
	"ode.stiff.bdf": {
		func() w15Trio {
			order := w15Limit(1, 5)
			return w15Mk(fmt.Sprintf("A BDF method of order %d uses the new value plus the previous values. Enter the total number of values in the formula.", order), order+1,
				fmt.Sprintf("The new value plus %d past values gives %d values.", order, order+1))
		},
		func() w15Trio {
			order := w15Limit(1, 5)
			return w15Mk(fmt.Sprintf("The BDF family is zero-stable for orders up to six. A formula of order %d is stable. Enter six minus that order.", order), 6-order,
				fmt.Sprintf("6 - %d = %d.", order, 6-order))
		},
		func() w15Trio {
			order := w15Limit(2, 6)
			return w15Mk(fmt.Sprintf("A BDF method of order %d uses the previous %d values. Enter the number of past values.", order, order), order,
				fmt.Sprintf("The method uses %d past values.", order))
		},
	},
	"ode.stiff.rosenbrock": {
		func() w15Trio {
			stages := w15Limit(2, 7)
			return w15Mk(fmt.Sprintf("A Rosenbrock method has %d stages and each stage solves a linear system. Enter the number of solves per step.", stages), stages,
				fmt.Sprintf("One solve per stage gives %d solves.", stages))
		},
		func() w15Trio {
			size := w15Limit(2, 9)
			return w15Mk(fmt.Sprintf("A dense Jacobian has size %d. Enter the number of entries.", size), size*size,
				fmt.Sprintf("%d squared is %d.", size, size*size))
		},
		func() w15Trio {
			stages, steps := w15Limit(2, 5), w15Limit(2, 9)
			return w15Mk(fmt.Sprintf("A Rosenbrock method has %d stages and takes %d steps. Enter the total number of linear solves.", stages, steps), stages*steps,
				fmt.Sprintf("%d times %d is %d.", stages, steps, stages*steps))
		},
	},
	"ode.stiff.stability_region": {
		func() w15Trio {
			quotient := w15Limit(1, 4)
			rate := w15Limit(2, 6)
			reach := quotient * rate
			return w15Mk(fmt.Sprintf("A method is stable on the negative real axis out to %d and the mode rate is %d. Enter the largest stable step as an integer.", reach, rate), quotient,
				fmt.Sprintf("%d divided by %d is %d.", reach, rate, quotient))
		},
		func() w15Trio {
			total, left := w15Limit(4, 9), w15Limit(1, 4)
			return w15Mk(fmt.Sprintf("A stability region covers %d sectors, and %d of them lie in the left half-plane. Enter the number in the right half-plane.", total, left), total-left,
				fmt.Sprintf("%d - %d = %d.", total, left, total-left))
		},
		func() w15Trio {
			real, imaginary := w15Limit(1, 9), w15Limit(1, 9)
			return w15Mk(fmt.Sprintf("A stability boundary crosses the real axis at -%d and the imaginary axis at %d. Enter the sum of the two magnitudes.", real, imaginary), real+imaginary,
				fmt.Sprintf("%d + %d = %d.", real, imaginary, real+imaginary))
		},
	},
	"ode.stiff.step_control": {
		func() w15Trio {
			order := w15Limit(1, 6)
			return w15Mk(fmt.Sprintf("A controller with a method of order %d raises the error ratio to the power 1 over the order plus one. Enter the denominator of that exponent.", order), order+1,
				fmt.Sprintf("The denominator is %d + 1, which is %d.", order, order+1))
		},
		func() w15Trio {
			attempts := w15Limit(3, 20)
			rejected := w15Limit(1, 2)
			return w15Mk(fmt.Sprintf("A solver attempts %d steps and rejects %d. Enter the accepted count.", attempts, rejected), attempts-rejected,
				fmt.Sprintf("%d - %d = %d.", attempts, rejected, attempts-rejected))
		},
		func() w15Trio {
			halvings := w15Limit(1, 5)
			return w15Mk(fmt.Sprintf("A step is halved %d times. Enter the denominator of the resulting step as a multiple of the original.", halvings), 1<<halvings,
				fmt.Sprintf("The denominator is 2 to the power %d, which is %d.", halvings, 1<<halvings))
		},
	},
	"ode.stiff.dae_index": {
		func() w15Trio {
			index := w15Limit(2, 7)
			return w15Mk(fmt.Sprintf("An index %d constraint must be differentiated before it can be solved for the algebraic variable. Enter the number of differentiations.", index), index-1,
				fmt.Sprintf("%d - 1 = %d.", index, index-1))
		},
		func() w15Trio {
			differential := w15Limit(2, 9)
			algebraic := w15Limit(1, 5)
			return w15Mk(fmt.Sprintf("A DAE has %d differential variables and %d algebraic variables. Enter the dimension of the differential state.", differential, algebraic), differential,
				fmt.Sprintf("The differential state has dimension %d.", differential))
		},
		func() w15Trio {
			differentiated := w15Limit(1, 7)
			return w15Mk(fmt.Sprintf("Differentiating the constraint %d times reduces the index to one. Enter the original index.", differentiated), differentiated+1,
				fmt.Sprintf("Each differentiation lowers the index by one, so the original index is %d + 1, which is %d.", differentiated, differentiated+1))
		},
	},
	"ode.stiff.symplectic": {
		func() w15Trio {
			steps := w15Limit(1, 6)
			return w15Mk(fmt.Sprintf("Stormer-Verlet uses one force evaluation per step. Enter the number of evaluations over %d steps.", steps), steps,
				fmt.Sprintf("One per step over %d steps is %d.", steps, steps))
		},
		func() w15Trio {
			dimension := w15Limit(1, 6)
			return w15Mk(fmt.Sprintf("A symplectic form on phase space has even dimension 2n, and n equals %d. Enter n.", dimension), dimension,
				fmt.Sprintf("The half-dimension is %d.", dimension))
		},
		func() w15Trio {
			bound := w15Limit(1, 9)
			return w15Mk(fmt.Sprintf("A symplectic method keeps the energy error between -%d and %d. Enter the peak-to-peak amplitude.", bound, bound), 2*bound,
				fmt.Sprintf("%d - (-%d) = %d.", bound, bound, 2*bound))
		},
	},
	"ode.variational.functional": {
		func() w15Trio {
			integrand, length := w15Limit(2, 9), w15Limit(2, 9)
			return w15Mk(fmt.Sprintf("A functional integrates the constant %d over an interval of length %d. Enter its value.", integrand, length), integrand*length,
				fmt.Sprintf("%d times %d is %d.", integrand, length, integrand*length))
		},
		func() w15Trio {
			amplitude, inverse := w15Limit(1, 4), w15Limit(2, 3)
			change := amplitude * inverse
			return w15Mk(fmt.Sprintf("A variation has amplitude %d and parameter 1 over %d. Enter the change as an integer.", change, inverse), amplitude,
				fmt.Sprintf("%d divided by %d is %d.", change, inverse, amplitude))
		},
		func() w15Trio {
			derivatives := w15Limit(1, 7)
			return w15Mk(fmt.Sprintf("A Lagrangian depends on derivatives up to order %d. Enter the order of the Euler-Lagrange equation.", derivatives), 2*derivatives,
				fmt.Sprintf("The equation has order twice %d, which is %d.", derivatives, 2*derivatives))
		},
	},
	"ode.variational.euler_lagrange": {
		func() w15Trio {
			slope := w15Limit(1, 9)
			return w15Mk(fmt.Sprintf("The Lagrangian is half the square of the slope. At slope %d, enter the derivative of the Lagrangian with respect to the slope.", slope), slope,
				fmt.Sprintf("Differentiating half the square gives the slope itself, which is %d.", slope))
		},
		func() w15Trio {
			coefficient := w15Limit(1, 9)
			return w15Mk(fmt.Sprintf("The Lagrangian is %d times the square of the slope. Enter its second derivative with respect to the slope.", coefficient), 2*coefficient,
				fmt.Sprintf("The second derivative is twice %d, which is %d.", coefficient, 2*coefficient))
		},
		func() w15Trio {
			coefficient := w15Limit(1, 9)
			return w15Mk(fmt.Sprintf("The Lagrangian is %d times x times the square of the slope. At x equal to 1 and slope equal to 1, enter its derivative with respect to the slope.", coefficient), 2*coefficient,
				fmt.Sprintf("The derivative is twice %d times x times the slope, which is %d at the given point.", coefficient, 2*coefficient))
		},
	},
	"ode.variational.natural_bc": {
		func() w15Trio {
			endpoints := w15Limit(2, 6)
			free := w15Limit(0, 2)
			if free >= endpoints {
				free = endpoints - 1
			}
			return w15Mk(fmt.Sprintf("A problem has %d endpoints and %d of them are free. Enter the number of essential boundary conditions.", endpoints, free), endpoints-free,
				fmt.Sprintf("%d - %d = %d fixed endpoints, each giving one essential condition.", endpoints, free, endpoints-free))
		},
		func() w15Trio {
			endpoints := w15Limit(2, 6)
			fixed := w15Limit(1, endpoints-1)
			return w15Mk(fmt.Sprintf("A problem has %d endpoints and %d of them are fixed. Enter the number of natural boundary conditions.", endpoints, fixed), endpoints-fixed,
				fmt.Sprintf("%d - %d = %d free endpoints, each giving one natural condition.", endpoints, fixed, endpoints-fixed))
		},
		func() w15Trio {
			free := w15Limit(1, 6)
			return w15Mk(fmt.Sprintf("A variational problem has %d free endpoints. Enter the number of natural conditions.", free), free,
				fmt.Sprintf("Each free endpoint supplies one natural condition, so there are %d.", free))
		},
	},
	"ode.variational.legendre": {
		func() w15Trio {
			coefficient := w15Limit(1, 9)
			return w15Mk(fmt.Sprintf("The function is %d times x squared. Enter its derivative at x equal to 1.", coefficient), 2*coefficient,
				fmt.Sprintf("The derivative is twice %d times x, which is %d at x equal to 1.", coefficient, 2*coefficient))
		},
		func() w15Trio {
			coefficient := w15Limit(1, 9)
			return w15Mk(fmt.Sprintf("The function is %d times x squared and the conjugate variable is twice %d. Enter the Legendre transform.", coefficient, coefficient), coefficient,
				fmt.Sprintf("For this quadratic the transform equals the function value at the conjugate point, which is %d.", coefficient))
		},
		func() w15Trio {
			slope := w15Limit(1, 5)
			return w15Mk(fmt.Sprintf("The function is x squared, so the Legendre transform is the square of the conjugate variable divided by four. At conjugate variable %d, enter the transform.", 2*slope), slope*slope,
				fmt.Sprintf("The transform is (2 times %d) squared over four, which is %d.", slope, slope*slope))
		},
	},
	"ode.variational.hamilton_equations": {
		func() w15Trio {
			p, q := w15Limit(1, 6), w15Limit(1, 6)
			return w15Mk(fmt.Sprintf("The Hamiltonian is p squared plus q squared. Enter it at p equal to %d and q equal to %d.", p, q), p*p+q*q,
				fmt.Sprintf("%d squared plus %d squared is %d.", p, q, p*p+q*q))
		},
		func() w15Trio {
			coefficient, momentum := w15Limit(2, 4), w15Limit(1, 9)
			return w15Mk(fmt.Sprintf("The Hamiltonian is the square of p, so q' equals twice p. At p equal to %d and with the coefficient %d, enter q'.", momentum, coefficient), coefficient*momentum,
				fmt.Sprintf("%d times %d is %d.", coefficient, momentum, coefficient*momentum))
		},
		func() w15Trio {
			coefficient, position := w15Limit(2, 4), w15Limit(1, 9)
			return w15Mk(fmt.Sprintf("The Hamiltonian is the square of q, so the magnitude of p' equals twice q. At q equal to %d and with the coefficient %d, enter that magnitude.", position, coefficient), coefficient*position,
				fmt.Sprintf("%d times %d is %d.", coefficient, position, coefficient*position))
		},
	},
	"ode.variational.noether": {
		func() w15Trio {
			mass, velocity := w15Limit(1, 6), w15Limit(1, 9)
			return w15Mk(fmt.Sprintf("Translation symmetry conserves momentum mass times velocity. With mass %d and velocity %d, enter the momentum.", mass, velocity), mass*velocity,
				fmt.Sprintf("%d times %d is %d.", mass, velocity, mass*velocity))
		},
		func() w15Trio {
			kinetic, potential := w15Limit(1, 9), w15Limit(1, 9)
			return w15Mk(fmt.Sprintf("Time translation conserves energy. With kinetic energy %d and potential energy %d, enter the energy.", kinetic, potential), kinetic+potential,
				fmt.Sprintf("%d + %d = %d.", kinetic, potential, kinetic+potential))
		},
		func() w15Trio {
			x, vy := w15Limit(1, 9), w15Limit(1, 9)
			return w15Mk(fmt.Sprintf("Rotation symmetry conserves angular momentum. At x equal to %d, y equal to 0 and vertical velocity %d with unit mass, enter the angular momentum.", x, vy), x*vy,
				fmt.Sprintf("%d times %d is %d.", x, vy, x*vy))
		},
	},
	"ode.variational.geodesic": {
		func() w15Trio {
			a, b := w15Limit(1, 9), w15Limit(1, 9)
			return w15Mk(fmt.Sprintf("In flat coordinates a geodesic runs from (0, 0) to (%d, %d). Enter its squared length.", a, b), a*a+b*b,
				fmt.Sprintf("%d squared plus %d squared is %d.", a, b, a*a+b*b))
		},
		func() w15Trio {
			radius := w15Limit(2, 9)
			return w15Mk(fmt.Sprintf("In polar coordinates the angular metric component is r squared. At r equal to %d, enter that component.", radius), radius*radius,
				fmt.Sprintf("%d squared is %d.", radius, radius*radius))
		},
		func() w15Trio {
			components := w15Limit(2, 6)
			return w15Mk(fmt.Sprintf("A diagonal metric in %d coordinates has one independent entry per coordinate. Enter the number of independent diagonal entries.", components), components,
				fmt.Sprintf("One entry per coordinate gives %d entries.", components))
		},
	},
	"ode.variational.brachistochrone": {
		func() w15Trio {
			scale := w15Limit(1, 9)
			return w15Mk(fmt.Sprintf("The cycloid has y equal to a times (1 - cos theta). At theta equal to pi and a equal to %d, enter y.", scale), 2*scale,
				fmt.Sprintf("1 - cos pi is 2, so y is twice %d, which is %d.", scale, 2*scale))
		},
		func() w15Trio {
			scale := w15Limit(1, 9)
			return w15Mk(fmt.Sprintf("The cycloid has y equal to a times (1 - cos theta). At theta equal to pi over 2 and a equal to %d, enter y.", scale), scale,
				fmt.Sprintf("1 - cos pi over 2 is 1, so y equals %d.", scale))
		},
		func() w15Trio {
			scale := w15Limit(2, 9)
			return w15Mk(fmt.Sprintf("The descent time scales as the square root of a. If a is multiplied by %d squared, enter the factor by which the time grows.", scale), scale,
				fmt.Sprintf("The square root of %d squared is %d.", scale, scale))
		},
	},
	"ode.variational.isoperimetric": {
		func() w15Trio {
			radius := w15Limit(1, 9)
			return w15Mk(fmt.Sprintf("A circle has radius %d. Enter its area divided by pi times its radius.", radius), radius,
				fmt.Sprintf("The area is pi times %d squared, and dividing by pi times %d leaves %d.", radius, radius, radius))
		},
		func() w15Trio {
			quotient := w15Limit(1, 9)
			return w15Mk(fmt.Sprintf("A circle of perimeter 2 pi k has radius k. For perimeter equal to 2 pi times %d, enter the radius.", quotient), quotient,
				fmt.Sprintf("The radius is the factor %d.", quotient))
		},
		func() w15Trio {
			side := w15Limit(2, 9)
			return w15Mk(fmt.Sprintf("A square with perimeter 4s has side s. For perimeter %d, enter the side.", 4*side), side,
				fmt.Sprintf("%d divided by 4 is %d.", 4*side, side))
		},
	},
	"ode.variational.second_variation": {
		func() w15Trio {
			first, second := w15Limit(1, 9), w15Limit(1, 9)
			return w15Mk(fmt.Sprintf("The second variation integrand is P times the square of eta prime plus Q times the square of eta. With P equal to %d, Q equal to %d, eta equal to 1 and eta prime equal to 1, enter its value.", first, second), first+second,
				fmt.Sprintf("%d + %d = %d.", first, second, first+second))
		},
		func() w15Trio {
			length := w15Limit(3, 12)
			conjugate := w15Limit(1, length-1)
			return w15Mk(fmt.Sprintf("An extremal runs to %d and has a conjugate point at %d. Enter the distance from the conjugate point to the endpoint.", length, conjugate), length-conjugate,
				fmt.Sprintf("%d - %d = %d.", length, conjugate, length-conjugate))
		},
		func() w15Trio {
			coefficient := w15Limit(1, 9)
			return w15Mk(fmt.Sprintf("The second variation has Q equal to %d and eta equal to 2. Enter Q times the square of eta.", coefficient), 4*coefficient,
				fmt.Sprintf("%d times 2 squared is %d.", coefficient, 4*coefficient))
		},
	},
}
