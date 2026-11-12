package optimization

import (
	"fmt"

	"github.com/chuma-beep/mathua/internal/generator"
)

// Wave 12 extends the optimization domain to one hundred concepts. Every
// family computes an integer from drawn parameters, so the numeric grader can
// check it and a feed can offer twenty distinct instances of each ordinal.
type trio struct {
	question    string
	answer      int
	explanation string
}

func mk(question string, answer int, explanation string) trio {
	return trio{question: question, answer: answer, explanation: explanation}
}

// w12Gen adapts three family functions to the generator interface.
type w12Gen struct {
	families [3]func() trio
}

func (g *w12Gen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ordinal(ctx)
	out := g.families[kp-1]()
	return at(kp, out.question, out.answer, out.explanation)
}

var wave12Families = map[string][3]func() trio{
	// -----------------------------------------------------------------------
	// Interior point methods and conic optimization
	// -----------------------------------------------------------------------
	"optim.central_path": {
		func() trio {
			mu := limit(1, 40)
			return mk(fmt.Sprintf("Minimize x subject to x ≥ 1 with the barrier -%d log(x - 1). Enter the central point x.", mu), 1+mu,
				fmt.Sprintf("Stationarity gives 1 - %d/(x - 1) = 0, so x - 1 = %d and x = %d.", mu, mu, 1+mu))
		},
		func() trio {
			x, y := limit(2, 12), limit(1, 9)
			return mk(fmt.Sprintf("On the central path the multiplier and the distance to the boundary satisfy λ(x - 1) = μ. At x = %d and λ = %d, enter μ.", x, y), y*(x-1),
				fmt.Sprintf("μ = %d·%d = %d.", y, x-1, y*(x-1)))
		},
		func() trio {
			mu, tol := limit(8, 64), limit(1, 3)
			steps, value := 0, mu
			for value > tol {
				value /= 2
				steps++
			}
			return mk(fmt.Sprintf("A path-following method halves the barrier parameter from %d until it is at most %d. Enter the number of halvings.", mu, tol), steps,
				fmt.Sprintf("%d halved %d times first reaches %d or less.", mu, steps, tol))
		},
	},
	"optim.path_following": {
		func() trio {
			mu, sigma := limit(10, 60), limit(1, 5)
			return mk(fmt.Sprintf("A path-following step reduces the barrier parameter from %d by the factor %d. Enter the new barrier parameter, using integer division.", mu, sigma), mu/sigma,
				fmt.Sprintf("%d / %d = %d.", mu, sigma, mu/sigma))
		},
		func() trio {
			sigma, mu := limit(1, 4), 5*limit(2, 12)
			return mk(fmt.Sprintf("A predictor-corrector step predicts a decrease of σμ with σ = %d and μ = %d. Enter that decrease.", sigma, mu), sigma*mu,
				fmt.Sprintf("The decrease is %d·%d = %d.", sigma, mu, sigma*mu))
		},
		func() trio {
			x, s := limit(2, 9), limit(1, 8)
			return mk(fmt.Sprintf("A centrality residual is x s - μ. At x = %d and s = %d with μ = 1, enter the residual.", x, s), x*s-1,
				fmt.Sprintf("%d·%d - 1 = %d.", x, s, x*s-1))
		},
	},
	"optim.primal_dual_interior": {
		func() trio {
			a, x, b := limit(2, 9), limit(1, 8), limit(1, 40)
			return mk(fmt.Sprintf("A primal residual is a x - b. At a = %d, x = %d and b = %d, enter it.", a, x, b), a*x-b,
				fmt.Sprintf("%d·%d - %d = %d.", a, x, b, a*x-b))
		},
		func() trio {
			x, s := limit(1, 9), limit(1, 9)
			return mk(fmt.Sprintf("A complementarity product is x s. At x = %d and s = %d, enter it.", x, s), x*s,
				fmt.Sprintf("The product is %d·%d = %d, and the method drives it to zero.", x, s, x*s))
		},
		func() trio {
			dx, x, ds, s := limit(1, 4), limit(4, 9), limit(1, 4), limit(4, 9)
			bound := minInt(x/dx, s/ds)
			return mk(fmt.Sprintf("A step must keep dx·t ≤ x and ds·t ≤ s. With dx = %d, x = %d, ds = %d and s = %d, enter the largest integer t, using integer division.", dx, x, ds, s), bound,
				fmt.Sprintf("The bounds are %d/%d = %d and %d/%d = %d, so the largest integer step is %d.", x, dx, x/dx, s, ds, s/ds, bound))
		},
	},
	"optim.newton_system": {
		func() trio {
			a, b := limit(2, 9), limit(2, 30)
			return mk(fmt.Sprintf("Solve %d x = %d for the Newton system. Enter x, using integer division.", a, b), b/a,
				fmt.Sprintf("x = %d / %d = %d.", b, a, b/a))
		},
		func() trio {
			h, g := limit(2, 9), limit(1, 20)
			return mk(fmt.Sprintf("A Newton step solves h s = -g. With h = %d and g = %d, enter the integer step magnitude.", h, g), g/h,
				fmt.Sprintf("|s| = %d / %d = %d.", g, h, g/h))
		},
		func() trio {
			a, b, c, d := limit(1, 6), limit(1, 6), limit(1, 6), limit(1, 6)
			return mk(fmt.Sprintf("A two by two system has the matrix [[%d, %d], [%d, %d]]. Enter its determinant.", a, b, c, d), a*d-b*c,
				fmt.Sprintf("The determinant is %d·%d - %d·%d = %d.", a, d, b, c, a*d-b*c))
		},
	},
	"optim.conic_program": {
		func() trio {
			x, c := limit(1, 12), limit(2, 9)
			return mk(fmt.Sprintf("A conic program minimizes c x over the nonnegative cone with c = %d. At x = %d, enter the objective.", c, x), c*x,
				fmt.Sprintf("The value is %d·%d = %d.", c, x, c*x))
		},
		func() trio {
			x, t := limit(0, 9), limit(0, 9)
			return mk(fmt.Sprintf("A point must satisfy x ≥ 0. At x = %d, enter its cone slack measured as x plus the offset %d.", x, t), x+t,
				fmt.Sprintf("%d + %d = %d; a nonnegative value means the point is in the cone.", x, t, x+t))
		},
		func() trio {
			a, y, c := limit(1, 5), limit(1, 6), limit(2, 9)
			return mk(fmt.Sprintf("Dual feasibility asks that c - A y lies in the dual cone. With one column a = %d, y = %d and c = %d, enter c - a y.", a, y, c), c-a*y,
				fmt.Sprintf("%d - %d·%d = %d.", c, a, y, c-a*y))
		},
	},
	"optim.second_order_cone": {
		func() trio {
			tv, x := limit(3, 12), limit(0, 9)
			return mk(fmt.Sprintf("A point (t, x) lies in the second-order cone when t² ≥ x². For t = %d and x = %d, enter t² - x².", tv, x), tv*tv-x*x,
				fmt.Sprintf("%d² - %d² = %d; a nonnegative value places the point in the cone.", tv, x, tv*tv-x*x))
		},
		func() trio {
			u, v := limit(1, 9), limit(1, 9)
			return mk(fmt.Sprintf("The cone inequality is t ≥ sqrt(u² + v²). For u = %d and v = %d, enter the smallest integer t that satisfies it.", u, v), ceilSqrt(u*u+v*v),
				fmt.Sprintf("sqrt(%d² + %d²) is at most %d, so the smallest integer t is %d.", u, v, ceilSqrt(u*u+v*v), ceilSqrt(u*u+v*v)))
		},
		func() trio {
			s, t, x := limit(1, 5), limit(2, 9), limit(0, 5)
			return mk(fmt.Sprintf("A cone slack is s(t - x) with s = %d, t = %d and x = %d. Enter it.", s, t, x), s*(t-x),
				fmt.Sprintf("%d·(%d - %d) = %d.", s, t, x, s*(t-x)))
		},
	},
	"optim.semidefinite": {
		func() trio {
			a, b, c := limit(1, 9), limit(0, 9), limit(1, 9)
			return mk(fmt.Sprintf("A symmetric two by two matrix has entries %d on the first diagonal, %d off the diagonal and %d on the second diagonal. Enter its trace.", a, b, c), a+c,
				fmt.Sprintf("The trace is %d + %d = %d.", a, c, a+c))
		},
		func() trio {
			a, b, c := limit(1, 9), limit(0, 6), limit(1, 9)
			return mk(fmt.Sprintf("A symmetric matrix [[%d, %d], [%d, %d]] must be positive semidefinite. Enter its determinant.", a, b, b, c), a*c-b*b,
				fmt.Sprintf("The determinant is %d·%d - %d² = %d; a nonnegative value is necessary for positive semidefiniteness.", a, c, b, a*c-b*b))
		},
		func() trio {
			a, b := limit(1, 9), limit(0, 6)
			return mk(fmt.Sprintf("A two by two symmetric matrix [[%d, %d], [%d, %d]] has eigenvalue sum equal to its trace. Enter the larger eigenvalue when the smaller is zero.", a, b, b, a), a+a,
				fmt.Sprintf("The trace is %d + %d = %d, so the other eigenvalue is %d.", a, a, 2*a, 2*a))
		},
	},
	"optim.dual_certificate": {
		func() trio {
			c, a, y := limit(2, 9), limit(1, 5), limit(1, 6)
			return mk(fmt.Sprintf("A dual certificate checks c - a y ≥ 0. With c = %d, a = %d and y = %d, enter c - a y.", c, a, y), c-a*y,
				fmt.Sprintf("%d - %d·%d = %d.", c, a, y, c-a*y))
		},
		func() trio {
			b, y := limit(2, 20), limit(1, 9)
			return mk(fmt.Sprintf("A dual objective is b y with b = %d and y = %d. Enter it.", b, y), b*y,
				fmt.Sprintf("%d·%d = %d; the value is a lower bound.", b, y, b*y))
		},
		func() trio {
			p, d := limit(20, 60), limit(1, 15)
			return mk(fmt.Sprintf("A primal candidate is %d and a dual certificate gives %d. Enter the gap.", p, p-d), d,
				fmt.Sprintf("%d - %d = %d, which no certificate can make negative.", p, p-d, d))
		},
	},
	"optim.infeasibility_certificate": {
		func() trio {
			a, y, b := limit(1, 6), limit(-6, 6), limit(1, 20)
			return mk(fmt.Sprintf("A Farkas certificate needs y A ≥ 0 and y b < 0. With a = %d, y = %d and b = %d, enter y a.", a, y, b), y*a,
				fmt.Sprintf("%d·%d = %d; the certificate also needs y b = %d to be negative.", y, a, y*a, y*b))
		},
		func() trio {
			y, b := limit(1, 9), limit(5, 30)
			return mk(fmt.Sprintf("An infeasibility certificate has y b = %d·%d. Enter its magnitude.", y, b), y*b,
				fmt.Sprintf("The magnitude is %d·%d = %d.", y, b, y*b))
		},
		func() trio {
			a, y := limit(1, 5), limit(1, 8)
			return mk(fmt.Sprintf("A certificate requires the product y a to be at least zero with a = %d and y = %d. Enter y a.", a, y), y*a,
				fmt.Sprintf("%d·%d = %d, which is nonnegative.", y, a, y*a))
		},
	},
	"optim.interior_convergence": {
		func() trio {
			start, tol := 2*limit(10, 40), limit(1, 4)
			value, steps := start, 0
			for value > tol {
				value /= 2
				steps++
			}
			return mk(fmt.Sprintf("An interior method halves the duality measure from %d until it is at most %d. Enter the number of iterations.", start, tol), steps,
				fmt.Sprintf("Halving %d reaches at most %d after %d iterations.", start, tol, steps))
		},
		func() trio {
			gap, factor := limit(30, 90), limit(2, 5)
			return mk(fmt.Sprintf("A method reduces the duality gap by the factor %d per iteration. From a gap of %d, enter the gap after one iteration, using integer division.", factor, gap), gap/factor,
				fmt.Sprintf("%d / %d = %d.", gap, factor, gap/factor))
		},
		func() trio {
			mu, sigma := limit(4, 20), limit(1, 3)
			return mk(fmt.Sprintf("A step sets μ_new = σ μ with σ = %d and μ = %d. Enter the reduction μ - μ_new.", sigma, mu), mu-sigma*mu,
				fmt.Sprintf("%d - %d·%d = %d.", mu, sigma, mu, mu-sigma*mu))
		},
	},
	// -----------------------------------------------------------------------
	// Stochastic optimization
	// -----------------------------------------------------------------------
	"optim.sgd": {
		func() trio {
			x, t, g := limit(1, 10), limit(1, 3), limit(1, 9)
			return mk(fmt.Sprintf("Take one stochastic step x_new = x - t g with x = %d, step size %d and sampled gradient %d. Enter x_new.", x, t, g), x-t*g,
				fmt.Sprintf("%d - %d·%d = %d.", x, t, g, x-t*g))
		},
		func() trio {
			a, x := limit(2, 9), limit(1, 8)
			return mk(fmt.Sprintf("For f(x) = a x² / 2 with a = %d, enter the gradient at x = %d.", a, x), a*x,
				fmt.Sprintf("The gradient is %d·%d = %d.", a, x, a*x))
		},
		func() trio {
			x, t, g := limit(2, 12), limit(1, 2), limit(1, 6)
			first := x - t*g
			second := first - t*g
			return mk(fmt.Sprintf("Apply the same step x_new = x - %d·%d twice from x = %d. Enter the value after two steps.", t, g, x), second,
				fmt.Sprintf("After one step x = %d, and after the second x = %d.", first, second))
		},
	},
	"optim.mini_batch": {
		func() trio {
			g1, g2 := limit(1, 9), limit(1, 9)
			return mk(fmt.Sprintf("A mini-batch of two samples has gradients %d and %d. Enter their average, using integer division.", g1, g2), (g1+g2)/2,
				fmt.Sprintf("(%d + %d)/2 = %d.", g1, g2, (g1+g2)/2))
		},
		func() trio {
			n, batch := limit(20, 90), limit(2, 9)
			return mk(fmt.Sprintf("A dataset has %d samples and the batch size is %d. Enter the number of complete batches, using integer division.", n, batch), n/batch,
				fmt.Sprintf("%d / %d = %d complete batches.", n, batch, n/batch))
		},
		func() trio {
			sigma, n := limit(2, 9), limit(2, 9)
			return mk(fmt.Sprintf("A per-sample gradient has standard deviation %d and the batch size is %d. Enter the batch standard deviation scaled by %d, using integer division.", sigma, n, n), sigma,
				fmt.Sprintf("The batch deviation is %d/sqrt(%d), and scaling by sqrt(%d) returns %d.", sigma, n, n, sigma))
		},
	},
	"optim.variance_reduction": {
		func() trio {
			v, n := limit(10, 60), limit(2, 9)
			return mk(fmt.Sprintf("A gradient estimator has variance %d. Averaging %d independent samples divides it by %d. Enter the new variance, using integer division.", v, n, n), v/n,
				fmt.Sprintf("%d / %d = %d.", v, n, v/n))
		},
		func() trio {
			g, c, e := limit(2, 9), limit(1, 6), limit(1, 4)
			return mk(fmt.Sprintf("A control variate estimator reports g - c + e with g = %d, c = %d and e = %d. Enter the estimate.", g, c, e), g-c+e,
				fmt.Sprintf("%d - %d + %d = %d.", g, c, e, g-c+e))
		},
		func() trio {
			sigma, rho := limit(4, 20), limit(1, 3)
			return mk(fmt.Sprintf("A control variate cuts the variance by the factor %d. From variance %d, enter the reduced variance, using integer division.", rho, sigma), sigma/rho,
				fmt.Sprintf("%d / %d = %d.", sigma, rho, sigma/rho))
		},
	},
	"optim.momentum": {
		func() trio {
			v, beta, g := limit(1, 9), limit(1, 3), limit(1, 9)
			return mk(fmt.Sprintf("A momentum update is v_new = β v + g with β = %d, v = %d and g = %d. Enter v_new.", beta, v, g), beta*v+g,
				fmt.Sprintf("%d·%d + %d = %d.", beta, v, g, beta*v+g))
		},
		func() trio {
			x, t, v := limit(2, 12), limit(1, 3), limit(1, 8)
			return mk(fmt.Sprintf("With momentum velocity %d and step size %d, a step from x = %d lands at x = %d - %d·%d. Enter the new x.", v, t, x, x, t, v), x-t*v,
				fmt.Sprintf("%d - %d·%d = %d.", x, t, v, x-t*v))
		},
		func() trio {
			beta, g := limit(2, 4), limit(1, 9)
			return mk(fmt.Sprintf("The steady-state velocity with constant gradient %d and β = %d is g/(1 - 1/β). Enter the integer value when β = %d.", g, beta, beta), g*beta/(beta-1),
				fmt.Sprintf("%d·%d/(%d - 1) = %d, using integer division.", g, beta, beta, g*beta/(beta-1)))
		},
	},
	"optim.adagrad": {
		func() trio {
			g1, g2 := limit(1, 9), limit(1, 9)
			return mk(fmt.Sprintf("AdaGrad accumulates squared gradients. After gradients %d and %d, enter the accumulated sum.", g1, g2), g1*g1+g2*g2,
				fmt.Sprintf("%d² + %d² = %d.", g1, g2, g1*g1+g2*g2))
		},
		func() trio {
			g := limit(2, 50)
			return mk(fmt.Sprintf("An AdaGrad step uses rate sqrt(G)/2 with G = %d. Enter twice the step magnitude at gradient %d, using integer division.", g*g, g), g,
				fmt.Sprintf("Twice the step is 2·%d/sqrt(%d) = %d.", g, g*g, g))
		},
		func() trio {
			acc, g := limit(4, 25), limit(1, 9)
			return mk(fmt.Sprintf("A coordinate has accumulated squared gradient %d. Add the next squared gradient %d² and enter the new sum.", acc, g), acc+g*g,
				fmt.Sprintf("%d + %d = %d.", acc, g*g, acc+g*g))
		},
	},
	"optim.adam": {
		func() trio {
			m, g, beta := limit(1, 9), limit(1, 9), limit(1, 4)
			return mk(fmt.Sprintf("Adam's first moment updates as m_new = β m + (1 - β) g. With β = %d in quarters, m = %d and g = %d, enter 4·m_new.", beta, m, g), beta*m+(4-beta)*g,
				fmt.Sprintf("4 m_new = %d·%d + (%d)·%d = %d.", beta, m, 4-beta, g, beta*m+(4-beta)*g))
		},
		func() trio {
			m, step := limit(2, 9), limit(1, 4)
			return mk(fmt.Sprintf("Bias correction divides the first moment %d by 1 - β^t with β = 1/2 and t = %d. Enter the corrected value, using integer division.", m, step), m*(1<<step)/((1<<step)-1),
				fmt.Sprintf("%d·%d/(%d - 1) = %d, using integer division.", m, 1<<step, 1<<step, m*(1<<step)/((1<<step)-1)))
		},
		func() trio {
			m, v := limit(1, 9), limit(1, 9)
			return mk(fmt.Sprintf("An Adam update uses the ratio m/sqrt(v). With m = %d and v = %d, enter m scaled by sqrt(v), using integer division.", m, v), m,
				fmt.Sprintf("m/sqrt(%d) times sqrt(%d) returns %d.", v, v, m))
		},
	},
	"optim.learning_rate_schedule": {
		func() trio {
			rate, decay, step := limit(4, 20), limit(1, 3), limit(1, 4)
			value := rate
			for i := 0; i < step; i++ {
				value /= (1 + decay)
			}
			return mk(fmt.Sprintf("A learning rate starts at %d and is divided by %d each step for %d steps. Enter the rate after those steps, using integer division.", rate, 1+decay, step), value,
				fmt.Sprintf("%d divided by %d, %d times, leaves %d.", rate, 1+decay, step, value))
		},
		func() trio {
			base, step := limit(2, 9), limit(1, 5)
			return mk(fmt.Sprintf("A schedule sets the rate to base/(1 + step) with base = %d and step = %d. Enter the rate, using integer division.", base, step), base/(1+step),
				fmt.Sprintf("%d / %d = %d.", base, 1+step, base/(1+step)))
		},
		func() trio {
			rate, factor := limit(6, 24), limit(2, 4)
			return mk(fmt.Sprintf("A schedule multiplies the rate by 1/%d. From %d, enter the decrease, using integer division.", factor, rate), rate-rate/factor,
				fmt.Sprintf("%d - %d/%d = %d.", rate, rate, factor, rate-rate/factor))
		},
	},
	"optim.stochastic_noise": {
		func() trio {
			variance, n := limit(4, 36), limit(2, 9)
			return mk(fmt.Sprintf("A gradient has variance %d and a batch averages %d samples. Enter the reduced variance, using integer division.", variance, n), variance/n,
				fmt.Sprintf("Averaging %d samples divides the variance by %d: %d/%d = %d.", n, n, variance, n, variance/n))
		},
		func() trio {
			v1, v2 := limit(1, 20), limit(1, 20)
			return mk(fmt.Sprintf("Two independent noise variances are %d and %d. Enter their sum.", v1, v2), v1+v2,
				fmt.Sprintf("%d + %d = %d.", v1, v2, v1+v2))
		},
		func() trio {
			signal, noise := limit(4, 20), limit(1, 4)
			return mk(fmt.Sprintf("A signal of strength %d competes with noise %d. Enter the signal-to-noise ratio, using integer division.", signal, noise), signal/noise,
				fmt.Sprintf("%d / %d = %d.", signal, noise, signal/noise))
		},
	},
	"optim.sample_complexity": {
		func() trio {
			sigma, eps := limit(4, 20), limit(1, 4)
			return mk(fmt.Sprintf("A sample bound is σ²/ε² with σ = %d and ε = %d. Enter it.", sigma, eps), (sigma*sigma)/(eps*eps),
				fmt.Sprintf("%d²/%d² = %d.", sigma, eps, (sigma*sigma)/(eps*eps)))
		},
		func() trio {
			delta := limit(2, 30)
			return mk(fmt.Sprintf("A bound needs 1/δ² samples with δ = 1/%d. Enter that number.", delta), delta*delta,
				fmt.Sprintf("1/(1/%d)² = %d² = %d.", delta, delta, delta*delta))
		},
		func() trio {
			s1, s2 := limit(10, 40), limit(1, 9)
			return mk(fmt.Sprintf("One estimator needs %d samples and another needs %d. Enter how many more the first needs.", s1, s2), s1-s2,
				fmt.Sprintf("%d - %d = %d.", s1, s2, s1-s2))
		},
	},
	"optim.online_gradient": {
		func() trio {
			x, g, t := limit(1, 9), limit(1, 9), limit(1, 4)
			return mk(fmt.Sprintf("An online update is x_new = x - g/t with x = %d, g = %d and t = %d. Enter x_new, using integer division.", x, g, t), x-g/t,
				fmt.Sprintf("%d - %d/%d = %d.", x, g, t, x-g/t))
		},
		func() trio {
			losses := []int{limit(1, 5), limit(1, 5), limit(1, 5)}
			return mk(fmt.Sprintf("Three online rounds incur losses %d, %d and %d. Enter the cumulative loss.", losses[0], losses[1], losses[2]), losses[0]+losses[1]+losses[2],
				fmt.Sprintf("%d + %d + %d = %d.", losses[0], losses[1], losses[2], losses[0]+losses[1]+losses[2]))
		},
		func() trio {
			rounds, total := limit(2, 6), limit(1, 5)
			return mk(fmt.Sprintf("A regret bound is sqrt(T) = %d after T = %d rounds, and the observed regret is %d. Enter the slack.", rounds, rounds*rounds, total), rounds-total,
				fmt.Sprintf("The bound %d minus the regret %d leaves %d.", rounds, total, rounds-total))
		},
	},
}

// ceilSqrt returns the smallest integer whose square is at least v.
func ceilSqrt(v int) int {
	n := 0
	for n*n < v {
		n++
	}
	return n
}

// registerWave12 adds the fifty Wave 12 concepts. Each concept is one generator
// backed by its three per-ordinal family functions.
func registerWave12(reg *generator.Registry) {
	for id, families := range wave12Families {
		reg.Register(id, &w12Gen{families: families})
	}
	for id, families := range wave12FamiliesB {
		reg.Register(id, &w12Gen{families: families})
	}
	for id, families := range wave12FamiliesC {
		reg.Register(id, &w12Gen{families: families})
	}
}
