package odes

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/grader"
)

func Register(reg *generator.Registry) {
	reg.Register("ode.basics.concept", &conceptGen{})
	reg.Register("ode.first_order.separable", &separableGen{})
	reg.Register("ode.first_order.linear", &linearFirstGen{})
	reg.Register("ode.first_order.exact", &exactGen{})
	reg.Register("ode.second_order.homogeneous", &homogeneousGen{})
	reg.Register("ode.second_order.nonhomogeneous", &nonhomogeneousGen{})
	reg.Register("ode.transforms.laplace", &laplaceGen{})
	reg.Register("ode.adv.systems", &systemsGen{})
	reg.Register("ode.second_order.variation_params", &variationParamsGen{})
	reg.Register("ode.transforms.convolution", &convolutionGen{})
	reg.Register("ode.systems.linear_phase", &linearPhaseGen{})
	reg.Register("ode.numerical.euler", &eulerGen{})
	reg.Register("ode.numerical.rk4", &rk4Gen{})
	reg.Register("ode.stability", &stabilityGen{})
	reg.Register("ode.pde.characteristics", &characteristicsGen{})
	reg.Register("ode.transforms.fourier", &fourierGen{})
	reg.Register("ode.nonlinear.bifurcation", &bifurcationGen{})
	reg.Register("ode.pde.wave", &waveGen{})
	reg.Register("ode.pde.heat", &heatGen{})
	reg.Register("ode.green.function", &greenFunctionGen{})
	reg.Register("ode.sturm_liouville", &sturmLiouvilleGen{})
	reg.Register("ode.boundary_value", &boundaryValueGen{})
	reg.Register("ode.nonlinear.lyapunov", &lyapunovGen{})
	reg.Register("ode.pde.laplace_eq", &laplaceEqGen{})
	reg.Register("ode.integral.volterra", &volterraGen{})
	reg.Register("ode.basics.slope_field", &slopeFieldGen{})
	reg.Register("ode.first_order.logistic", &logisticGen{})
	reg.Register("ode.basics.existence_uniqueness", &existenceUniquenessGen{})
	reg.Register("ode.second_order.forced", &forcedOscillationGen{})
	registerNumeric(reg)
	registerBatch(reg)
	registerWave4(reg)
	registerWave15(reg)
}

type conceptGen struct{}

func (g *conceptGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := rand.Intn(5)
	switch mode {
	case 0:
		// Order of differential equation y^(n) + a y' + b y = 0
		n := smallInt(1, 5+int(ctx.Difficulty*3))
		a := smallInt(2, 6)
		b := smallInt(1, 9)
		derivStr := "y'"
		switch n {
		case 1:
			derivStr = "y'"
		case 2:
			derivStr = "y''"
		case 3:
			derivStr = "y'''"
		default:
			derivStr = fmt.Sprintf("y^{(%d)}", n)
		}
		return generator.Problem{
			Question: fmt.Sprintf("What is the order of the differential equation \\(%s + %dy' + %dy = 0\\)?", derivStr, a, b),
			Answer:   fmt.Sprintf("%d", n),
			Explanation: fmt.Sprintf(
				"The order of a differential equation is the order of the highest derivative appearing in the equation. Here the highest derivative is \\(%s\\), so the order is %d.",
				derivStr, n,
			),
		}

	case 1:
		// Degree (exponent of highest derivative)
		n := smallInt(1, 3)
		deg := smallInt(2, 5)
		a := smallInt(2, 6)
		derivStr := "y'"
		if n == 2 {
			derivStr = "y''"
		} else if n == 3 {
			derivStr = "y'''"
		}
		return generator.Problem{
			Question: fmt.Sprintf("What is the degree (the power of the highest-order derivative) of the differential equation \\((%s)^{%d} + %dy = 0\\)?", derivStr, deg, a),
			Answer:   fmt.Sprintf("%d", deg),
			Explanation: fmt.Sprintf(
				"The degree of a differential equation is the exponent of the highest-order derivative when written in polynomial form. The highest-order derivative is \\(%s\\) raised to power %d, giving degree %d.",
				derivStr, deg, deg,
			),
		}

	case 2:
		// Number of essential arbitrary constants in n-th order ODE general solution
		n := smallInt(1, 6+int(ctx.Difficulty*3))
		return generator.Problem{
			Question: fmt.Sprintf("How many essential arbitrary constants appear in the general solution of a linear ordinary differential equation of order %d?", n),
			Answer:   fmt.Sprintf("%d", n),
			Explanation: fmt.Sprintf(
				"An ordinary differential equation of order %d requires %d successive integrations to solve, introducing exactly %d essential arbitrary constants in its general solution.",
				n, n, n,
			),
		}

	case 3:
		// Order in Leibniz notation d^k y / dx^k
		k := smallInt(2, 6)
		c := smallInt(2, 8)
		return generator.Problem{
			Question: fmt.Sprintf("What is the order of the differential equation \\(\\frac{d^{%d}y}{dx^{%d}} + %dy = 0\\)?", k, k, c),
			Answer:   fmt.Sprintf("%d", k),
			Explanation: fmt.Sprintf(
				"In Leibniz notation, \\(\\frac{d^{%d}y}{dx^{%d}}\\) denotes the %d-th derivative of \\(y\\) with respect to \\(x\\). Therefore the equation has order %d.",
				k, k, k, k,
			),
		}

	default:
		// State-space system dimension
		n := smallInt(2, 6)
		a := smallInt(2, 7)
		return generator.Problem{
			Question: fmt.Sprintf("Converting the scalar %d-th order differential equation \\(y^{(%d)} + %dy = 0\\) into an equivalent system of first-order differential equations produces a state vector of what dimension?", n, n, a),
			Answer:   fmt.Sprintf("%d", n),
			Explanation: fmt.Sprintf(
				"Defining state variables \\(x_1 = y, x_2 = y', \\dots, x_%d = y^{(%d)}\\) converts the %d-th order scalar ODE into a first-order system with a state vector of dimension %d.",
				n, n-1, n, n,
			),
		}
	}
}

type separableGen struct{}

func (g *separableGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := rand.Intn(5)
	switch mode {
	case 0:
		// y dy/dx = a x with y(0) = y0, evaluate [y(x1)]^2
		a := smallInt(1, 4)
		y0 := smallInt(1, 5)
		x1 := smallInt(1, 4)
		val := a*x1*x1 + y0*y0
		return generator.Problem{
			Question: fmt.Sprintf("Solve the separable initial value problem \\(y \\frac{dy}{dx} = %d x\\) with \\(y(0) = %d\\) (for \\(y > 0\\)). What is the value of \\([y(%d)]^2\\)?", a, y0, x1),
			Answer:   fmt.Sprintf("%d", val),
			Explanation: fmt.Sprintf(
				"Separating variables gives \\(y dy = %d x dx\\). Integrating both sides yields \\(\\frac{y^2}{2} = \\frac{%d x^2}{2} + C\\), so \\(y^2 = %d x^2 + 2C\\). Using initial condition \\(y(0) = %d\\) gives \\(2C = (%d)^2 = %d\\). Thus \\(y(x)^2 = %d x^2 + %d\\). Evaluating at \\(x = %d\\) gives %d.",
				a, a, a, y0, y0, y0*y0, a, y0*y0, x1, val,
			),
		}

	case 1:
		// dy/dx = (2 alpha x + b) y with y(0) = y0 -> y(x) = A exp(alpha x^2 + b x)
		alpha := smallInt(1, 3)
		b := smallInt(1, 4)
		y0 := smallInt(1, 6)
		coeff := 2 * alpha
		return generator.Problem{
			Question: fmt.Sprintf("The separable initial value problem \\(\\frac{dy}{dx} = (%d x + %d) y\\) with \\(y(0) = %d\\) has solution of the form \\(y(x) = A \\exp(%d x^2 + %d x)\\). What is the value of the coefficient \\(A\\)?", coeff, b, y0, alpha, b),
			Answer:   fmt.Sprintf("%d", y0),
			Explanation: fmt.Sprintf(
				"Separating variables yields \\(\\frac{dy}{y} = (%d x + %d) dx\\). Integrating gives \\(\\ln|y| = %d x^2 + %d x + C\\), so \\(y(x) = A \\exp(%d x^2 + %d x)\\) where \\(A = e^C\\). Applying initial condition \\(y(0) = %d\\) gives %d.",
				coeff, b, alpha, b, alpha, b, y0, y0,
			),
		}

	case 2:
		// dy/dx = 3 alpha x^2 + 2 beta x + c with y(0) = y0, evaluate y(x1)
		alpha := smallInt(1, 3)
		beta := smallInt(1, 3)
		c := smallInt(1, 5)
		y0 := smallInt(1, 5)
		x1 := smallInt(1, 3)
		aCoeff := 3 * alpha
		bCoeff := 2 * beta
		val := alpha*x1*x1*x1 + beta*x1*x1 + c*x1 + y0
		return generator.Problem{
			Question: fmt.Sprintf("Solve the separable initial value problem \\(\\frac{dy}{dx} = %d x^2 + %d x + %d\\) with \\(y(0) = %d\\). What is the value of \\(y(%d)\\)?", aCoeff, bCoeff, c, y0, x1),
			Answer:   fmt.Sprintf("%d", val),
			Explanation: fmt.Sprintf(
				"Separating variables gives \\(dy = (%d x^2 + %d x + %d) dx\\). Integrating both sides yields \\(y(x) = %d x^3 + %d x^2 + %d x + C\\). Applying \\(y(0) = %d\\) gives \\(C = %d\\). Evaluating at \\(x = %d\\) gives %d.",
				aCoeff, bCoeff, c, alpha, beta, c, y0, y0, x1, val,
			),
		}

	case 3:
		// Autonomous separable logistic/quadratic: dy/dx = a y - b y^2 = 0 non-zero equilibrium
		b := smallInt(2, 5)
		k := smallInt(2, 6+int(ctx.Difficulty*2))
		a := b * k
		return generator.Problem{
			Question: fmt.Sprintf("For the autonomous separable differential equation \\(\\frac{dy}{dx} = %d y - %d y^2\\), what is the non-zero equilibrium solution \\(y\\)?", a, b),
			Answer:   fmt.Sprintf("%d", k),
			Explanation: fmt.Sprintf(
				"Equilibrium solutions occur where the derivative vanishes: \\(\\frac{dy}{dx} = y(%d - %d y) = 0\\). Setting the non-zero factor to zero yields \\(%d - %d y = 0\\), so \\(y = \\frac{%d}{%d} = %d\\).",
				a, b, a, b, a, b, k,
			),
		}

	default:
		// y^2 dy/dx = a with y(0) = y0, evaluate [y(x1)]^3
		a := smallInt(1, 4)
		y0 := smallInt(1, 4)
		x1 := smallInt(1, 3)
		val := 3*a*x1 + y0*y0*y0
		return generator.Problem{
			Question: fmt.Sprintf("Solve the separable differential equation \\(y^2 \\frac{dy}{dx} = %d\\) subject to initial condition \\(y(0) = %d\\). What is the value of \\([y(%d)]^3\\)?", a, y0, x1),
			Answer:   fmt.Sprintf("%d", val),
			Explanation: fmt.Sprintf(
				"Separating variables gives \\(y^2 dy = %d dx\\). Integrating both sides gives \\(\\frac{y^3}{3} = %d x + C\\), so \\(y^3 = %d x + 3C\\). With \\(y(0) = %d\\), we find \\(3C = (%d)^3 = %d\\), yielding \\(y(x)^3 = %d x + %d\\). Evaluating at \\(x = %d\\) gives %d.",
				a, a, 3*a, y0, y0, y0*y0*y0, 3*a, y0*y0*y0, x1, val,
			),
		}
	}
}

type linearFirstGen struct{}

func (g *linearFirstGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	for {
		mode := rand.Intn(5)
		switch mode {
		case 0:
			// Integrating factor evaluation at x0: y' + (a/x) y = Q(x)
			a := smallInt(2, 4)
			x0 := smallInt(2, 4)
			// I(x) = exp(int a/x dx) = x^a
			val := 1
			for i := 0; i < a; i++ {
				val *= x0
			}
			ansStr := fmt.Sprintf("%d", val)
			q := fmt.Sprintf("For the first-order linear differential equation \\(y' + \\frac{%d}{x}y = x^2\\) (with \\(x > 0\\)), what is the value of the integrating factor \\(I(x) = \\exp\\left(\\int P(x)dx\\right)\\) evaluated at \\(x = %d\\)?", a, x0)
			if strings.Contains(q, ansStr) {
				continue
			}
			return generator.Problem{
				Question: q,
				Answer:   ansStr,
				Explanation: fmt.Sprintf(
					"The coefficient is \\(P(x) = \\frac{%d}{x}\\). The integrating factor is \\(I(x) = \\exp\\left(\\int \\frac{%d}{x} dx\\right) = \\exp(%d \\ln x) = x^{%d}\\). Evaluating at \\(x = %d\\) gives %d.",
					a, a, a, a, x0, val,
				),
			}

		case 1:
			// IVP integration constant C in y' + a y = b with y(0) = y0
			a := smallInt(2, 5)
			k := smallInt(2, 4)
			b := a * k
			y0 := smallInt(1, 8)
			constC := y0 - k
			if constC == 0 {
				continue
			}
			ansStr := fmt.Sprintf("%d", constC)
			q := fmt.Sprintf("For the first-order linear initial value problem \\(y' + %d y = %d\\) with \\(y(0) = %d\\), what is the value of the integration constant \\(C\\) in the general solution \\(y(x) = y_p + C e^{-%d x}\\)?", a, b, y0, a)
			if strings.Contains(q, ansStr) || strings.Contains(q, strings.TrimPrefix(ansStr, "-")) {
				continue
			}
			return generator.Problem{
				Question: q,
				Answer:   ansStr,
				Explanation: fmt.Sprintf(
					"The differential equation has particular solution \\(y_p = \\frac{%d}{%d} = %d\\) and homogeneous solution \\(y_h = C e^{-%d x}\\). Setting \\(x = 0\\) gives \\(y(0) = %d + C = %d\\), so \\(C = %d - %d = %d\\).",
					b, a, k, a, k, y0, y0, k, constC,
				),
			}

		case 2:
			// Long-term equilibrium limit lim_{x -> inf} y(x) for y' + a y = b (a > 0)
			a := smallInt(2, 6)
			k := smallInt(1, 8+int(ctx.Difficulty*3))
			b := a * k
			ansStr := fmt.Sprintf("%d", k)
			q := fmt.Sprintf("For the stable first-order linear differential equation \\(y' + %d y = %d\\), what is the long-term equilibrium limit \\(\\lim_{x \\to \\infty} y(x)\\)?", a, b)
			if strings.Contains(q, ansStr) {
				continue
			}
			return generator.Problem{
				Question: q,
				Answer:   ansStr,
				Explanation: fmt.Sprintf(
					"The general solution is \\(y(x) = %d + C e^{-%d x}\\). Because the damping coefficient %d is positive, the transient exponential term \\(C e^{-%d x} \\to 0\\) as \\(x \\to \\infty\\). Thus the steady-state equilibrium limit is %d.",
					k, a, a, a, k,
				),
			}

		case 3:
			// Homogeneous solution evaluation: x y' - k y = 0 with y(1) = y0, evaluate at x1
			k := smallInt(1, 3)
			y0 := smallInt(1, 4)
			x1 := smallInt(2, 3)
			xPow := 1
			for i := 0; i < k; i++ {
				xPow *= x1
			}
			val := y0 * xPow
			ansStr := fmt.Sprintf("%d", val)
			q := fmt.Sprintf("Solve the initial value problem \\(x y' - %d y = 0\\) with \\(y(1) = %d\\) for \\(x > 0\\). What is the value of \\(y(%d)\\)?", k, y0, x1)
			if strings.Contains(q, ansStr) {
				continue
			}
			return generator.Problem{
				Question: q,
				Answer:   ansStr,
				Explanation: fmt.Sprintf(
					"Separating variables yields \\(\\frac{dy}{y} = \\frac{%d}{x} dx\\), giving \\(\\ln|y| = %d \\ln x + C\\) and general solution \\(y(x) = C x^{%d}\\). Applying initial condition \\(y(1) = %d\\) yields \\(C(1)^{%d} = %d\\), so \\(y(x) = %d x^{%d}\\). Evaluating at \\(x = %d\\) gives %d.",
					k, k, k, y0, k, y0, y0, k, x1, val,
				),
			}

		default:
			// Particular solution polynomial constant coefficient B in y' + a y = b x + c
			a := smallInt(2, 4)
			k := smallInt(1, 4)
			b := a * k
			m := smallInt(1, 5)
			c := k + a*m
			ansStr := fmt.Sprintf("%d", m)
			q := fmt.Sprintf("The differential equation \\(y' + %d y = %d x + %d\\) possesses a linear particular solution of the form \\(y_p(x) = A x + B\\). What is the value of the constant coefficient \\(B\\)?", a, b, c)
			if strings.Contains(q, ansStr) {
				continue
			}
			return generator.Problem{
				Question: q,
				Answer:   ansStr,
				Explanation: fmt.Sprintf(
					"Substituting \\(y_p = A x + B\\) into the differential equation gives \\(A + %d(A x + B) = %d A x + (A + %d B) = %d x + %d\\). Equating coefficients of \\(x\\) yields \\(%d A = %d\\), so \\(A = %d\\). Equating constant terms gives \\(%d + %d B = %d\\), so \\(%d B = %d\\), which yields \\(B = %d\\).",
					a, a, a, b, c, a, b, k, k, a, c, a, c-k, m,
				),
			}
		}
	}
}

type exactGen struct{}

func (g *exactGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	for {
		mode := rand.Intn(5)
		switch mode {
		case 0:
			// Exactness condition: partial M / partial y = partial N / partial x
			// M = (k x^2 y + b x) dx, N = (c x^3 + d y) dy
			// dM/dy = k x^2, dN/dx = 3 c x^2 => k = 3 c
			c := smallInt(2, 6+int(ctx.Difficulty*2))
			b := smallInt(1, 8)
			d := smallInt(1, 8)
			kVal := 3 * c
			ansStr := fmt.Sprintf("%d", kVal)
			q := fmt.Sprintf("What value of the constant \\(k\\) makes the differential equation \\((k x^2 y + %d x) dx + (%d x^3 + %d y) dy = 0\\) exact?", b, c, d)
			if strings.Contains(q, ansStr) {
				continue
			}
			return generator.Problem{
				Question: q,
				Answer:   ansStr,
				Explanation: fmt.Sprintf(
					"A first-order differential equation \\(M dx + N dy = 0\\) is exact when \\(\\frac{\\partial M}{\\partial y} = \\frac{\\partial N}{\\partial x}\\). Here \\(\\frac{\\partial M}{\\partial y} = k x^2\\) and \\(\\frac{\\partial N}{\\partial x} = 3(%d) x^2 = %d x^2\\). Equating partial derivatives gives \\(k = %d\\).",
					c, kVal, kVal,
				),
			}

		case 1:
			// Partial derivative evaluation dM/dy at (x0, y0)
			a := smallInt(2, 5)
			b := smallInt(2, 6)
			x0 := smallInt(1, 3)
			y0 := smallInt(1, 3)
			val := 2 * a * x0 * y0
			ansStr := fmt.Sprintf("%d", val)
			q := fmt.Sprintf("For the differential equation \\(M dx + N dy = 0\\) with \\(M(x, y) = %d x y^2 + %d x^2\\), what is the value of the partial derivative \\(\\frac{\\partial M}{\\partial y}\\) evaluated at the point \\((%d, %d)\\)?", a, b, x0, y0)
			if strings.Contains(q, ansStr) {
				continue
			}
			return generator.Problem{
				Question: q,
				Answer:   ansStr,
				Explanation: fmt.Sprintf(
					"Differentiating \\(M(x, y) = %d x y^2 + %d x^2\\) with respect to \\(y\\) treats \\(x\\) as constant, giving \\(\\frac{\\partial M}{\\partial y} = 2(%d) x y = %d x y\\). Evaluating at \\((%d, %d)\\) gives \\(%d(%d)(%d) = %d\\).",
					a, b, a, 2*a, x0, y0, 2*a, x0, y0, val,
				),
			}

		case 2:
			// Potential function evaluation Psi(x0, y0) with Psi(0, 0) = 0
			a := smallInt(1, 4)
			b := smallInt(1, 4)
			c := smallInt(1, 4)
			x0 := smallInt(1, 3)
			y0 := smallInt(1, 3)
			val := a*x0*x0*y0 + b*x0*x0 + c*y0*y0
			ansStr := fmt.Sprintf("%d", val)
			q := fmt.Sprintf("The exact differential equation \\((%d x y + %d x) dx + (%d x^2 + %d y) dy = 0\\) possesses a potential function \\(\\Psi(x, y)\\) such that \\(d\\Psi = M dx + N dy\\) and \\(\\Psi(0, 0) = 0\\). Evaluate \\(\\Psi(%d, %d)\\).", 2*a, 2*b, a, 2*c, x0, y0)
			if strings.Contains(q, ansStr) {
				continue
			}
			return generator.Problem{
				Question: q,
				Answer:   ansStr,
				Explanation: fmt.Sprintf(
					"Integrating \\(M(x, y)\\) with respect to \\(x\\) yields \\(\\Psi(x, y) = %d x^2 y + %d x^2 + g(y)\\). Differentiating with respect to \\(y\\) gives \\(\\frac{\\partial \\Psi}{\\partial y} = %d x^2 + g'(y) = %d x^2 + %d y\\), meaning \\(g'(y) = %d y\\) and \\(g(y) = %d y^2\\). With \\(\\Psi(0,0)=0\\), the potential function is \\(\\Psi(x, y) = %d x^2 y + %d x^2 + %d y^2\\). Evaluating at \\((%d, %d)\\) gives %d.",
					a, b, a, a, 2*c, 2*c, c, a, b, c, x0, y0, val,
				),
			}

		case 3:
			// Integrating factor exponent p for mu(x) = x^p
			p := smallInt(2, 5+int(ctx.Difficulty*2))
			coeff := p + 1
			ansStr := fmt.Sprintf("%d", p)
			q := fmt.Sprintf("The differential equation \\((%d y) dx + x dy = 0\\) is not exact, but multiplying by an integrating factor of the form \\(\\mu(x) = x^p\\) transforms it into an exact equation. What is the value of the exponent \\(p\\)?", coeff)
			if strings.Contains(q, ansStr) {
				continue
			}
			return generator.Problem{
				Question: q,
				Answer:   ansStr,
				Explanation: fmt.Sprintf(
					"Computing \\(\\frac{1}{N}\\left(\\frac{\\partial M}{\\partial y} - \\frac{\\partial N}{\\partial x}\\right) = \\frac{%d - 1}{x} = \\frac{%d}{x}\\). The integrating factor is \\(\\mu(x) = \\exp\\left(\\int \\frac{%d}{x} dx\\right) = x^{%d}\\), giving exponent \\(p = %d\\).",
					coeff, p, p, p, p,
				),
			}

		default:
			// First integral IVP constant C
			a := smallInt(1, 3)
			b := smallInt(1, 3)
			c := smallInt(1, 3)
			x0 := smallInt(1, 3)
			y0 := smallInt(1, 3)
			constC := a*x0*x0 + b*x0*y0 + c*y0*y0
			ansStr := fmt.Sprintf("%d", constC)
			q := fmt.Sprintf("The exact differential equation \\((%d x + %d y) dx + (%d x + %d y) dy = 0\\) has first integral \\(%d x^2 + %d x y + %d y^2 = C\\). Given the initial condition \\(y(%d) = %d\\), what is the value of the constant \\(C\\)?", 2*a, b, b, 2*c, a, b, c, x0, y0)
			if strings.Contains(q, ansStr) {
				continue
			}
			return generator.Problem{
				Question: q,
				Answer:   ansStr,
				Explanation: fmt.Sprintf(
					"Substituting the initial condition \\(x = %d\\) and \\(y = %d\\) into the first integral yields \\(C = %d(%d)^2 + %d(%d)(%d) + %d(%d)^2 = %d + %d + %d = %d\\).",
					x0, y0, a, x0, b, x0, y0, c, y0, a*x0*x0, b*x0*y0, c*y0*y0, constC,
				),
			}
		}
	}
}

type homogeneousGen struct{}

func (g *homogeneousGen) AnswerFormat() generator.AnswerFormat {
	return generator.AnswerFormat{
		Hint:      "Select the option or enter the answer",
		InputMode: "text",
	}
}

func (g *homogeneousGen) Grade(expected, userAnswer string) grader.Result {
	normExp := strings.TrimSpace(expected)
	normUser := strings.TrimSpace(userAnswer)
	if normExp == normUser {
		return grader.Result{Correct: true, Score: 1}
	}
	r := grader.NewRouter()
	if strings.Contains(expected, "\\(y(x)") {
		return r.Grade(grader.GradingMultipleChoice, expected, userAnswer)
	}
	return r.Grade(grader.GradingNumeric, expected, userAnswer)
}

func (g *homogeneousGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	for {
		mode := rand.Intn(5)
		switch mode {
		case 0:
			// Characteristic equation roots for y'' - (r1+r2) y' + r1*r2 y = 0
			r1 := smallInt(1, 5)
			r2 := r1 + smallInt(1, 5)
			sumR := r1 + r2
			prodR := r1 * r2
			ansStr := fmt.Sprintf("%d", r2)
			q := fmt.Sprintf("What is the larger characteristic root \\(r\\) of the differential equation \\(y'' - %d y' + %d y = 0\\)?", sumR, prodR)
			if strings.Contains(q, ansStr) {
				continue
			}
			ex := fmt.Sprintf(
				"The characteristic equation is \\(r^2 - %d r + %d = (r - %d)(r - %d) = 0\\). The roots are \\(r = %d\\) and \\(r = %d\\), so the larger root is %d.",
				sumR, prodR, r1, r2, r1, r2, r2,
			)
			return generator.Problem{
				Question:    q,
				Answer:      ansStr,
				Explanation: ex,
				Facts: map[string]string{
					"r1":     strconv.Itoa(r1),
					"r2":     strconv.Itoa(r2),
					"answer": ansStr,
				},
			}

		case 1:
			// Repeated root general solution y(x) = (c1 + c2 x) e^{r0 x}
			r0 := smallInt(1, 5)
			bCoeff := 2 * r0
			cSq := r0 * r0

			correct := fmt.Sprintf(`\(y(x) = (c_1 + c_2 x) e^{%d x}\)`, r0)
			dist1 := fmt.Sprintf(`\(y(x) = c_1 e^{%d x} + c_2 e^{-%d x}\)`, r0, r0)
			dist2 := fmt.Sprintf(`\(y(x) = c_1 \cos(%d x) + c_2 \sin(%d x)\)`, r0, r0)
			dist3 := fmt.Sprintf(`\(y(x) = c_1 e^{%d x} + c_2 x^2 e^{%d x}\)`, r0, r0)

			choices := []string{correct, dist1, dist2, dist3}
			rand.Shuffle(len(choices), func(i, j int) { choices[i], choices[j] = choices[j], choices[i] })

			q := fmt.Sprintf("Which expression is the general solution of the differential equation \\(y'' - %d y' + %d y = 0\\)?\n\nA) %s\n\nB) %s\n\nC) %s\n\nD) %s",
				bCoeff, cSq, choices[0], choices[1], choices[2], choices[3])

			ex := fmt.Sprintf(
				"The characteristic equation is \\(r^2 - %d r + %d = (r - %d)^2 = 0\\), giving repeated root \\(r = %d\\). For a repeated real root, the linearly independent solutions are \\(e^{%d x}\\) and \\(x e^{%d x}\\), yielding general solution %s.",
				bCoeff, cSq, r0, r0, r0, r0, correct,
			)
			return generator.Problem{
				Question:    q,
				Answer:      correct,
				Explanation: ex,
				Facts: map[string]string{
					"r0":     strconv.Itoa(r0),
					"answer": correct,
				},
			}

		case 2:
			// Underdamped complex conjugate angular frequency beta for y'' + 2 alpha y' + (alpha^2 + beta^2) y = 0
			alpha := smallInt(1, 4)
			beta := smallInt(2, 7)
			bCoeff := 2 * alpha
			cVal := alpha*alpha + beta*beta
			discNeg := 4 * beta * beta
			ansStr := fmt.Sprintf("%d", beta)
			q := fmt.Sprintf("The underdamped differential equation \\(y'' + %d y' + %d y = 0\\) has complex characteristic roots \\(r = -\\alpha \\pm i \\beta\\). What is the oscillatory angular frequency \\(\\beta\\)?", bCoeff, cVal)
			if strings.Contains(q, ansStr) {
				continue
			}
			ex := fmt.Sprintf(
				"The characteristic equation is \\(r^2 + %d r + %d = 0\\). By the quadratic formula, \\(r = \\frac{-%d \\pm \\sqrt{(%d)^2 - 4(%d)}}{2} = \\frac{-%d \\pm \\sqrt{-%d}}{2} = -%d \\pm %d i\\). Thus the angular frequency is %d.",
				bCoeff, cVal, bCoeff, bCoeff, cVal, bCoeff, discNeg, alpha, beta, beta,
			)
			return generator.Problem{
				Question:    q,
				Answer:      ansStr,
				Explanation: ex,
			}

		case 3:
			// Symbolic general solution for distinct real roots r1, r2: y(x) = c_1 e^{r1 x} + c_2 e^{r2 x}
			r1 := smallInt(1, 4)
			r2 := r1 + smallInt(1, 3)
			sumR := r1 + r2
			prodR := r1 * r2

			correct := fmt.Sprintf(`\(y(x) = c_1 e^{%d x} + c_2 e^{%d x}\)`, r1, r2)
			dist1 := fmt.Sprintf(`\(y(x) = c_1 e^{-%d x} + c_2 e^{-%d x}\)`, r1, r2)
			dist2 := fmt.Sprintf(`\(y(x) = (c_1 + c_2 x) e^{%d x}\)`, r1)
			dist3 := fmt.Sprintf(`\(y(x) = c_1 \cos(%d x) + c_2 \sin(%d x)\)`, r1, r2)

			choices := []string{correct, dist1, dist2, dist3}
			rand.Shuffle(len(choices), func(i, j int) { choices[i], choices[j] = choices[j], choices[i] })

			q := fmt.Sprintf("Which expression is the general solution of the differential equation \\(y'' - %d y' + %d y = 0\\)?\n\nA) %s\n\nB) %s\n\nC) %s\n\nD) %s",
				sumR, prodR, choices[0], choices[1], choices[2], choices[3])

			ex := fmt.Sprintf(
				"The characteristic equation is \\(r^2 - %d r + %d = (r - %d)(r - %d) = 0\\), with distinct real roots \\(r_1 = %d\\) and \\(r_2 = %d\\). Therefore the general solution is %s.",
				sumR, prodR, r1, r2, r1, r2, correct,
			)
			return generator.Problem{
				Question:    q,
				Answer:      correct,
				Explanation: ex,
				Facts: map[string]string{
					"r1":     strconv.Itoa(r1),
					"r2":     strconv.Itoa(r2),
					"answer": correct,
				},
			}

		default:
			// Harmonic oscillator amplitude A = sqrt(y0^2 + (v0/omega)^2) with Pythagorean pairs
			type pair struct {
				y0, k, a int
			}
			pairs := []pair{
				{3, 4, 5},
				{4, 3, 5},
				{6, 8, 10},
				{8, 6, 10},
				{5, 12, 13},
				{12, 5, 13},
				{7, 24, 25},
				{24, 7, 25},
				{8, 15, 17},
				{15, 8, 17},
				{9, 12, 15},
				{12, 9, 15},
				{10, 24, 26},
				{24, 10, 26},
				{15, 20, 25},
				{20, 15, 25},
				{20, 21, 29},
				{21, 20, 29},
			}
			p := pairs[rand.Intn(len(pairs))]
			omega := smallInt(2, 4)
			omegaSq := omega * omega
			v0 := omega * p.k
			ansStr := fmt.Sprintf("%d", p.a)
			q := fmt.Sprintf("The undamped oscillator \\(y'' + %d y = 0\\) satisfies initial displacement \\(y(0) = %d\\) and velocity \\(y'(0) = %d\\). What is the amplitude \\(A\\) of the resulting motion \\(y(t) = A \\cos(\\omega t - \\phi)\\)?", omegaSq, p.y0, v0)
			if strings.Contains(q, ansStr) {
				continue
			}
			return generator.Problem{
				Question: q,
				Answer:   ansStr,
				Explanation: fmt.Sprintf(
					"The natural frequency is \\(\\omega = \\sqrt{%d} = %d\\). The solution is \\(y(t) = %d \\cos(%d t) + %d \\sin(%d t)\\). The amplitude is \\(A = \\sqrt{(%d)^2 + (%d)^2} = \\sqrt{%d + %d} = \\sqrt{%d} = %d\\).",
					omegaSq, omega, p.y0, omega, p.k, omega, p.y0, p.k, p.y0*p.y0, p.k*p.k, p.a*p.a, p.a,
				),
			}
		}
	}
}

type nonhomogeneousGen struct{}

func (g *nonhomogeneousGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	for attempt := 0; attempt < 200; attempt++ {
		mode := rand.Intn(5)
		var prob generator.Problem
		switch mode {
		case 0:
			// Constant particular solution y_p for y'' + b y' + c y = K
			b := smallInt(1, 5)
			c := smallInt(2, 6)
			m := smallInt(1, 7+int(ctx.Difficulty*3))
			kVal := c * m
			prob = generator.Problem{
				Question: fmt.Sprintf("Find the constant particular solution \\(y_p\\) for the differential equation \\(y'' + %dy' + %dy = %d\\).", b, c, kVal),
				Answer:   fmt.Sprintf("%d", m),
				Explanation: fmt.Sprintf(
					"For a constant particular solution \\(y_p = A\\), the derivatives vanish: \\(y_p' = 0\\) and \\(y_p'' = 0\\). Substituting gives \\(0 + 0 + %d A = %d\\), which yields \\(A = \\frac{%d}{%d} = %d\\).",
					c, kVal, kVal, c, m,
				),
			}

		case 1:
			// Undamped resonance angular frequency omega_0 for a y'' + b y = cos(omega t)
			omega0 := smallInt(2, 9)
			a := smallInt(1, 4)
			b := a * omega0 * omega0
			prob = generator.Problem{
				Question: fmt.Sprintf("At what driving angular frequency \\(\\omega\\) does pure resonance occur for the driven undamped oscillator \\(%d y'' + %d y = \\cos(\\omega t)\\)?", a, b),
				Answer:   fmt.Sprintf("%d", omega0),
				Explanation: fmt.Sprintf(
					"Pure resonance occurs when the driving frequency matches the natural frequency \\(\\omega_0 = \\sqrt{\\frac{%d}{%d}} = %d\\).",
					b, a, omega0,
				),
			}

		case 2:
			// Exponential particular solution amplitude A for y'' - a^2 y = K e^(b x)
			a := smallInt(1, 3)
			b := smallInt(4, 6)
			m := smallInt(1, 5)
			aSq := a * a
			bSq := b * b
			denom := bSq - aSq
			kVal := m * denom
			prob = generator.Problem{
				Question: fmt.Sprintf("The differential equation \\(y'' - %d y = %d e^{%d x}\\) has a particular solution of the form \\(y_p(x) = A e^{%d x}\\). What is the value of the amplitude coefficient \\(A\\)?", aSq, kVal, b, b),
				Answer:   fmt.Sprintf("%d", m),
				Explanation: fmt.Sprintf(
					"Computing derivatives gives \\(y_p' = %d A e^{%d x}\\) and \\(y_p'' = %d A e^{%d x}\\). Substituting into the differential equation yields \\((%d - %d) A e^{%d x} = %d A e^{%d x} = %d e^{%d x}\\). Solving for \\(A\\) gives \\(A = \\frac{%d}{%d} = %d\\).",
					b, b, bSq, b, bSq, aSq, b, denom, b, kVal, b, kVal, denom, m,
				),
			}

		case 3:
			// Polynomial particular solution linear coefficient A_1 for y'' + c y = K x
			c := smallInt(2, 5)
			m := smallInt(1, 6)
			kVal := c * m
			prob = generator.Problem{
				Question: fmt.Sprintf("The differential equation \\(y'' + %d y = %d x\\) has a particular solution \\(y_p(x) = A_1 x + A_0\\). What is the value of the linear coefficient \\(A_1\\)?", c, kVal),
				Answer:   fmt.Sprintf("%d", m),
				Explanation: fmt.Sprintf(
					"For \\(y_p(x) = A_1 x + A_0\\), the second derivative is \\(y_p'' = 0\\). Substituting into the differential equation gives \\(0 + %d(A_1 x + A_0) = %d A_1 x + %d A_0 = %d x\\). Equating linear coefficients yields \\(%d A_1 = %d\\), so \\(A_1 = %d\\).",
					c, c, c, kVal, c, kVal, m,
				),
			}

		default:
			// Initial value problem particular solution constant C_1 in y(x) = C_1 cos(omega x) + y_p with y(0) = y_0
			omega := smallInt(2, 5)
			m := smallInt(1, 5)
			y0 := smallInt(1, 9)
			if y0 == m {
				continue
			}
			omegaSq := omega * omega
			kVal := m * omegaSq
			constC1 := y0 - m
			prob = generator.Problem{
				Question: fmt.Sprintf("For the initial value problem \\(y'' + %d y = %d\\) with initial conditions \\(y(0) = %d\\) and \\(y'(0) = 0\\), the general solution is \\(y(x) = C_1 \\cos(%d x) + y_p\\) where \\(y_p\\) is a constant particular solution. What is the value of the constant \\(C_1\\)?", omegaSq, kVal, y0, omega),
				Answer:   fmt.Sprintf("%d", constC1),
				Explanation: fmt.Sprintf(
					"The constant particular solution is \\(y_p = \\frac{%d}{%d} = %d\\). The general solution with zero initial velocity is \\(y(x) = C_1 \\cos(%d x) + %d\\). Applying initial condition \\(y(0) = C_1 + %d = %d\\) yields \\(C_1 = %d - %d = %d\\).",
					kVal, omegaSq, m, omega, m, m, y0, y0, m, constC1,
				),
			}
		}

		cleanAns := strings.TrimPrefix(prob.Answer, "-")
		if strings.Contains(prob.Question, prob.Answer) || (cleanAns != "" && strings.Contains(prob.Question, cleanAns)) {
			continue
		}
		return prob
	}
	// Fallback that is provably free of read-off
	return generator.Problem{
		Question:    "Find the constant particular solution \\(y_p\\) for the differential equation \\(y'' + 4y' + 3y = 21\\).",
		Answer:      "7",
		Explanation: "Setting derivatives to 0 gives 3*y_p = 21, hence y_p = 7.",
	}
}

type laplaceGen struct{}

func (g *laplaceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := rand.Intn(5)
	switch mode {
	case 0:
		a := smallInt(1, 4)
		d := smallInt(1, 3)
		m := smallInt(1, 4+int(ctx.Difficulty*3))
		s0 := a + d
		k := m * d
		return generator.Problem{
			Question:    fmt.Sprintf("Evaluate the Laplace transform \\(\\mathcal{L}\\{%d e^{%d t}\\}(s)\\) at \\(s = %d\\).", k, a, s0),
			Answer:      fmt.Sprintf("%d", m),
			Explanation: fmt.Sprintf("The Laplace transform is \\(\\mathcal{L}\\{e^{at}\\}(s) = 1/(s - a)\\). Thus \\(\\mathcal{L}\\{%d e^{%d t}\\}(s) = %d / (s - %d)\\). At \\(s = %d\\), the value is %d / %d = %d.", k, a, k, a, s0, k, d, m),
		}
	case 1:
		s0 := smallInt(2, 4)
		m := smallInt(1, 4+int(ctx.Difficulty*3))
		k := m * s0 * s0
		return generator.Problem{
			Question:    fmt.Sprintf("Evaluate the Laplace transform \\(\\mathcal{L}\\{%d t\\}(s)\\) at \\(s = %d\\).", k, s0),
			Answer:      fmt.Sprintf("%d", m),
			Explanation: fmt.Sprintf("The Laplace transform of \\(t\\) is \\(\\mathcal{L}\\{t\\}(s) = 1/s^2\\). Thus \\(\\mathcal{L}\\{%d t\\}(%d) = %d / %d^2 = %d\\).", k, s0, k, s0, m),
		}
	case 2:
		s0 := smallInt(2, 5)
		m := smallInt(2, 6+int(ctx.Difficulty*4))
		k := m * s0
		return generator.Problem{
			Question:    fmt.Sprintf("Evaluate the Laplace transform \\(\\mathcal{L}\\{%d\\}(s)\\) at \\(s = %d\\).", k, s0),
			Answer:      fmt.Sprintf("%d", m),
			Explanation: fmt.Sprintf("The Laplace transform of a constant \\(k\\) is \\(k/s\\). Evaluating at \\(s = %d\\) gives %d / %d = %d.", s0, k, s0, m),
		}
	case 3:
		a := smallInt(1, 3)
		d := smallInt(1, 3)
		m := smallInt(1, 4+int(ctx.Difficulty*3))
		s0 := a + d
		k := m * d * d
		return generator.Problem{
			Question:    fmt.Sprintf("By the first shifting theorem, evaluate \\(\\mathcal{L}\\{%d t e^{%d t}\\}(s)\\) at \\(s = %d\\).", k, a, s0),
			Answer:      fmt.Sprintf("%d", m),
			Explanation: fmt.Sprintf("By the first shifting theorem, \\(\\mathcal{L}\\{t e^{at}\\}(s) = 1/(s - a)^2\\). For \\(s = %d\\) and \\(a = %d\\), \\(s - a = %d\\). Thus \\(\\mathcal{L}\\{%d t e^{%d t}\\}(%d) = %d / (%d)^2 = %d\\).", s0, a, d, k, a, s0, k, d, m),
		}
	default:
		s0 := smallInt(2, 5)
		Y0 := smallInt(2, 6+int(ctx.Difficulty*3))
		y0 := smallInt(1, 5)
		ans := s0*Y0 - y0
		return generator.Problem{
			Question:    fmt.Sprintf("Given \\(y(0) = %d\\) and \\(Y(%d) = %d\\), compute the value of \\(\\mathcal{L}\\{y'\\}(%d) = s Y(s) - y(0)\\).", y0, s0, Y0, s0),
			Answer:      fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("By the derivative property of the Laplace transform, \\(\\mathcal{L}\\{y'\\}(s) = s Y(s) - y(0)\\). Substituting \\(s = %d\\), \\(Y(%d) = %d\\) and \\(y(0) = %d\\) gives (%d)(%d) - %d = %d.", s0, s0, Y0, y0, s0, Y0, y0, ans),
		}
	}
}

type systemsGen struct{}

func (g *systemsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := rand.Intn(5)
	switch mode {
	case 0:
		// Trace of coefficient matrix A: tr(A) = a + d
		a := smallInt(1, 6)
		b := smallInt(1, 6)
		c := smallInt(1, 6)
		d := smallInt(1, 6)
		tr := a + d
		return generator.Problem{
			Question: fmt.Sprintf("For the linear system of differential equations \\(x' = %d x + %d y\\) and \\(y' = %d x + %d y\\), what is the trace \\(\\operatorname{tr}(A)\\) of the coefficient matrix \\(A\\)?", a, b, c, d),
			Answer:   fmt.Sprintf("%d", tr),
			Explanation: fmt.Sprintf(
				"The coefficient matrix is \\(A = \\begin{pmatrix} %d & %d \\\\ %d & %d \\end{pmatrix}\\). The trace is the sum of the diagonal entries: \\(\\operatorname{tr}(A) = %d + %d = %d\\).",
				a, b, c, d, a, d, tr,
			),
		}

	case 1:
		// Determinant of coefficient matrix A: det(A) = a*d - b*c
		a := smallInt(1, 5)
		b := smallInt(1, 5)
		c := smallInt(1, 5)
		d := smallInt(1, 5)
		detA := a*d - b*c
		return generator.Problem{
			Question: fmt.Sprintf("For the linear system \\(x' = %d x + %d y\\) and \\(y' = %d x + %d y\\), what is the determinant \\(\\det(A)\\) of the coefficient matrix \\(A\\)?", a, b, c, d),
			Answer:   fmt.Sprintf("%d", detA),
			Explanation: fmt.Sprintf(
				"The coefficient matrix is \\(A = \\begin{pmatrix} %d & %d \\\\ %d & %d \\end{pmatrix}\\). Its determinant is \\(\\det(A) = (%d)(%d) - (%d)(%d) = %d - %d = %d\\).",
				a, b, c, d, a, d, b, c, a*d, b*c, detA,
			),
		}

	case 2:
		// Triangular matrix eigenvalues: max(lambda1, lambda2)
		l1 := smallInt(1, 5)
		l2 := l1 + smallInt(1, 4)
		b := smallInt(1, 6)
		maxL := l2
		return generator.Problem{
			Question: fmt.Sprintf("For the coupled system \\(x' = %d x + %d y\\) and \\(y' = %d y\\), what is the largest eigenvalue \\(\\lambda_{\\max}\\) of the system matrix?", l1, b, l2),
			Answer:   fmt.Sprintf("%d", maxL),
			Explanation: fmt.Sprintf(
				"The coefficient matrix \\(A = \\begin{pmatrix} %d & %d \\\\ 0 & %d \\end{pmatrix}\\) is upper triangular, so its eigenvalues are the diagonal entries %d and %d. The largest eigenvalue is %d.",
				l1, b, l2, l1, l2, maxL,
			),
		}

	case 3:
		// Uncoupled system initial velocity squared: [x'(0)]^2 + [y'(0)]^2
		a := smallInt(1, 4)
		b := smallInt(1, 4)
		x0 := smallInt(1, 3)
		y0 := smallInt(1, 3)
		vx := a * x0
		vy := b * y0
		vSq := vx*vx + vy*vy
		return generator.Problem{
			Question: fmt.Sprintf("For the system \\(x' = %d x\\) and \\(y' = %d y\\) with initial conditions \\(x(0) = %d\\) and \\(y(0) = %d\\), what is the squared initial speed \\([x'(0)]^2 + [y'(0)]^2\\)?", a, b, x0, y0),
			Answer:   fmt.Sprintf("%d", vSq),
			Explanation: fmt.Sprintf(
				"At \\(t = 0\\), the velocity components are \\(x'(0) = %d(%d) = %d\\) and \\(y'(0) = %d(%d) = %d\\). The squared initial speed is \\((%d)^2 + (%d)^2 = %d + %d = %d\\).",
				a, x0, vx, b, y0, vy, vx, vy, vx*vx, vy*vy, vSq,
			),
		}

	default:
		// Phase portrait discriminant D = tau^2 - 4*Delta
		a := smallInt(1, 4)
		b := smallInt(1, 3)
		c := smallInt(1, 3)
		d := smallInt(1, 4)
		tau := a + d
		detA := a*d - b*c
		disc := tau*tau - 4*detA
		return generator.Problem{
			Question: fmt.Sprintf("For the planar linear system with coefficient matrix \\(A = \\begin{pmatrix} %d & %d \\\\ %d & %d \\end{pmatrix}\\), what is the discriminant \\(D = \\tau^2 - 4\\Delta\\) that classifies the phase portrait?", a, b, c, d),
			Answer:   fmt.Sprintf("%d", disc),
			Explanation: fmt.Sprintf(
				"The trace is \\(\\tau = %d + %d = %d\\) and determinant is \\(\\Delta = (%d)(%d) - (%d)(%d) = %d - %d = %d\\). The discriminant is \\(D = \\tau^2 - 4\\Delta = (%d)^2 - 4(%d) = %d - %d = %d\\).",
				a, d, tau, a, d, b, c, a*d, b*c, detA, tau, detA, tau*tau, 4*detA, disc,
			),
		}
	}
}

type variationParamsGen struct{}

func (g *variationParamsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := rand.Intn(5)
	switch mode {
	case 0:
		// Wronskian of exponential basis y1 = e^(a x), y2 = e^(b x) evaluated at x = 0
		a := smallInt(1, 4)
		b := a + smallInt(1, 4)
		w0 := b - a
		return generator.Problem{
			Question: fmt.Sprintf("For the fundamental solutions \\(y_1(x) = e^{%d x}\\) and \\(y_2(x) = e^{%d x}\\), what is the value of the Wronskian \\(W(y_1, y_2)\\) evaluated at \\(x = 0\\)?", a, b),
			Answer:   fmt.Sprintf("%d", w0),
			Explanation: fmt.Sprintf(
				"The Wronskian is \\(W(x) = y_1 y_2' - y_1' y_2 = e^{%d x}(%d e^{%d x}) - (%d e^{%d x})e^{%d x} = (%d - %d)e^{%d x} = %d e^{%d x}\\). Evaluating at \\(x = 0\\) yields %d.",
				a, b, b, a, a, b, b, a, a+b, w0, a+b, w0,
			),
		}

	case 1:
		// Wronskian of trigonometric basis y1 = cos(omega x), y2 = sin(omega x)
		omega := smallInt(2, 8+int(ctx.Difficulty*2))
		omegaSq := omega * omega
		return generator.Problem{
			Question: fmt.Sprintf("For the differential equation \\(y'' + %d y = g(x)\\), the homogeneous basis is \\(y_1(x) = \\cos(%d x)\\) and \\(y_2(x) = \\sin(%d x)\\). What is the constant value of their Wronskian \\(W(y_1, y_2)\\)?", omegaSq, omega, omega),
			Answer:   fmt.Sprintf("%d", omega),
			Explanation: fmt.Sprintf(
				"Computing the Wronskian yields \\(W(x) = y_1 y_2' - y_1' y_2 = \\cos(%d x)(%d \\cos(%d x)) - (-%d \\sin(%d x))\\sin(%d x) = %d(\\cos^2(%d x) + \\sin^2(%d x)) = %d\\).",
				omega, omega, omega, omega, omega, omega, omega, omega, omega, omega,
			),
		}

	case 2:
		// Variation of parameters derivative u_2'(x) for y'' + omega^2 y = k sec(omega x)
		omega := smallInt(2, 5)
		m := smallInt(1, 6)
		kVal := m * omega
		omegaSq := omega * omega
		return generator.Problem{
			Question: fmt.Sprintf("In applying variation of parameters to \\(y'' + %d y = %d \\sec(%d x)\\) with \\(y_1 = \\cos(%d x)\\), \\(y_2 = \\sin(%d x)\\) and Wronskian \\(W = %d\\), what is the constant value of the derivative \\(u_2'(x) = \\frac{y_1(x) g(x)}{W(x)}\\)?", omegaSq, kVal, omega, omega, omega, omega),
			Answer:   fmt.Sprintf("%d", m),
			Explanation: fmt.Sprintf(
				"Substituting gives \\(u_2'(x) = \\frac{\\cos(%d x) \\cdot %d \\sec(%d x)}{%d} = \\frac{%d}{%d} = %d\\), because \\(\\cos(%d x)\\sec(%d x) = 1\\).",
				omega, kVal, omega, omega, kVal, omega, m, omega, omega,
			),
		}

	case 3:
		// Particular solution value y_p(0) for y'' + omega^2 y = K
		omega := smallInt(2, 5)
		m := smallInt(1, 6)
		omegaSq := omega * omega
		kVal := m * omegaSq
		return generator.Problem{
			Question: fmt.Sprintf("For the differential equation \\(y'' + %d y = %d\\), variation of parameters gives a particular solution \\(y_p(x)\\). What is the value of \\(y_p(0)\\)?", omegaSq, kVal),
			Answer:   fmt.Sprintf("%d", m),
			Explanation: fmt.Sprintf(
				"The homogeneous solutions are \\(y_1 = \\cos(%d x)\\) and \\(y_2 = \\sin(%d x)\\) with Wronskian \\(W = %d\\). The particular solution is \\(y_p = \\frac{%d}{%d} = %d\\), so \\(y_p(0) = %d\\).",
				omega, omega, omega, kVal, omegaSq, m, m,
			),
		}

	default:
		// Wronskian of Cauchy-Euler basis y1 = x, y2 = x^2 at x = x0
		x0 := smallInt(2, 7)
		ans := x0 * x0
		return generator.Problem{
			Question: fmt.Sprintf("For the differential equation with fundamental solutions \\(y_1(x) = x\\) and \\(y_2(x) = x^2\\) (for \\(x > 0\\)), what is the value of the Wronskian \\(W(y_1, y_2)\\) evaluated at \\(x = %d\\)?", x0),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf(
				"The Wronskian is \\(W(x) = y_1 y_2' - y_1' y_2 = x(2x) - (1)(x^2) = 2x^2 - x^2 = x^2\\). Evaluating at \\(x = %d\\) yields %d.",
				x0, ans,
			),
		}
	}
}

type convolutionGen struct{}

func (g *convolutionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := rand.Intn(5)
	switch mode {
	case 0:
		k := smallInt(1, 3+int(ctx.Difficulty*2))
		a := 2 * k
		b := smallInt(1, 4+int(ctx.Difficulty*3))
		tVal := smallInt(1, 3+int(ctx.Difficulty*2))
		ans := k * b * tVal * tVal
		return generator.Problem{
			Question:    fmt.Sprintf("Compute the convolution value \\((f * g)(%d)\\) where \\(f(t) = %d\\) and \\(g(t) = %d t\\).", tVal, a, b),
			Answer:      fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("The continuous convolution is \\((f * g)(t) = \\int_0^t f(\\tau) g(t - \\tau) d\\tau = \\int_0^t %d \\cdot %d (t - \\tau) d\\tau = %d \\left[ t\\tau - \\frac{\\tau^2}{2} \\right]_0^t = %d \\cdot \\frac{t^2}{2}\\). Evaluating at \\(t = %d\\) gives %d.", a, b, a*b, a*b, tVal, ans),
		}
	case 1:
		k := smallInt(1, 3+int(ctx.Difficulty*2))
		a := 2 * k
		b := smallInt(1, 4+int(ctx.Difficulty*3))
		tVal := smallInt(1, 3+int(ctx.Difficulty*2))
		ans := k * b * tVal * tVal
		return generator.Problem{
			Question:    fmt.Sprintf("Compute the convolution value \\((f * g)(%d)\\) where \\(f(t) = %d t\\) and \\(g(t) = %d\\).", tVal, a, b),
			Answer:      fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("The convolution integral is \\((f * g)(t) = \\int_0^t a \\tau \\cdot b d\\tau = %d \\int_0^t \\tau d\\tau = %d \\frac{t^2}{2}\\). Evaluating at \\(t = %d\\) gives %d.", a*b, a*b, tVal, ans),
		}
	case 2:
		a := smallInt(2, 5+int(ctx.Difficulty*3))
		b := smallInt(2, 5+int(ctx.Difficulty*3))
		tVal := smallInt(2, 6+int(ctx.Difficulty*4))
		ans := a * b * tVal
		return generator.Problem{
			Question:    fmt.Sprintf("Compute the convolution \\((f * g)(%d)\\) where \\(f(t) = %d\\) and \\(g(t) = %d\\).", tVal, a, b),
			Answer:      fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("The convolution of two constant functions is \\((f * g)(t) = \\int_0^t %d \\cdot %d d\\tau = %d t\\). At \\(t = %d\\), the value is %d.", a, b, a*b, tVal, ans),
		}
	case 3:
		k := smallInt(1, 2+int(ctx.Difficulty*2))
		a := 6 * k
		m := smallInt(1, 3+int(ctx.Difficulty*2))
		tVal := smallInt(1, 3)
		ans := k * m * tVal * tVal * tVal
		return generator.Problem{
			Question:    fmt.Sprintf("Compute the convolution value \\((f * g)(%d)\\) where \\(f(t) = %d t\\) and \\(g(t) = %d t\\).", tVal, a, m),
			Answer:      fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("The convolution integral is \\((f * g)(t) = \\int_0^t %d \\tau \\cdot %d (t - \\tau) d\\tau = %d \\left[ \\frac{t \\tau^2}{2} - \\frac{\\tau^3}{3} \\right]_0^t = %d \\frac{t^3}{6} = %d t^3\\). Evaluating at \\(t = %d\\) gives %d.", a, m, a*m, a*m, k*m, tVal, ans),
		}
	default:
		k := smallInt(1, 2+int(ctx.Difficulty*2))
		a := 3 * k
		m := smallInt(1, 3+int(ctx.Difficulty*2))
		tVal := smallInt(1, 3)
		ans := k * m * tVal * tVal * tVal
		return generator.Problem{
			Question:    fmt.Sprintf("Compute the convolution value \\((f * g)(%d)\\) where \\(f(t) = %d t^2\\) and \\(g(t) = %d\\).", tVal, a, m),
			Answer:      fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("The convolution integral is \\((f * g)(t) = \\int_0^t %d \\tau^2 \\cdot %d d\\tau = %d \\left[ \\frac{\\tau^3}{3} \\right]_0^t = %d t^3\\). Evaluating at \\(t = %d\\) gives %d.", a, m, a*m, k*m, tVal, ans),
		}
	}
}

type linearPhaseGen struct{}

func (g *linearPhaseGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := rand.Intn(5)
	switch mode {
	case 0:
		a := smallInt(-5, 5)
		d := smallInt(-5, 5)
		b := smallInt(-4, 4)
		c := smallInt(-4, 4)
		ans := a*d - b*c
		return generator.Problem{
			Question:    fmt.Sprintf("For the planar linear autonomous system \\(x' = Ax\\) with \\(A = \\begin{pmatrix} %d & %d \\\\ %d & %d \\end{pmatrix}\\), compute the determinant \\(\\det(A)\\).", a, b, c, d),
			Answer:      fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("The determinant of matrix \\(A\\) is \\(\\det(A) = (%d)(%d) - (%d)(%d) = %d\\).", a, d, b, c, ans),
		}
	case 1:
		a := smallInt(-6, 6)
		d := smallInt(-6, 6)
		ans := a + d
		b := smallInt(-4, 4)
		c := smallInt(-4, 4)
		return generator.Problem{
			Question:    fmt.Sprintf("For the planar linear system \\(x' = Ax\\) with \\(A = \\begin{pmatrix} %d & %d \\\\ %d & %d \\end{pmatrix}\\), compute the trace \\(\\operatorname{tr}(A)\\).", a, b, c, d),
			Answer:      fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("The trace is the sum of diagonal elements: \\(\\operatorname{tr}(A) = %d + (%d) = %d\\).", a, d, ans),
		}
	case 2:
		a := smallInt(-4, 4)
		d := smallInt(-4, 4)
		b := smallInt(-3, 3)
		c := smallInt(-3, 3)
		tau := a + d
		delta := a*d - b*c
		ans := tau*tau - 4*delta
		return generator.Problem{
			Question:    fmt.Sprintf("For the planar linear autonomous system \\(x' = Ax\\) with \\(A = \\begin{pmatrix} %d & %d \\\\ %d & %d \\end{pmatrix}\\), compute the characteristic discriminant \\(\\mathcal{D} = \\tau^2 - 4\\Delta\\), where \\(\\tau = \\operatorname{tr}(A)\\) and \\(\\Delta = \\det(A)\\).", a, b, c, d),
			Answer:      fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("The trace is \\(\\tau = %d\\) and the determinant is \\(\\Delta = %d\\). The discriminant is \\(\\tau^2 - 4\\Delta = (%d)^2 - 4(%d) = %d\\).", tau, delta, tau, delta, ans),
		}
	case 3:
		l1 := smallInt(-6, 6)
		l2 := smallInt(-6, 6)
		for l2 == l1 {
			l2 = smallInt(-6, 6)
		}
		b := smallInt(1, 5)
		ans := max(l1, l2)
		return generator.Problem{
			Question:    fmt.Sprintf("For the planar system \\(x' = \\begin{pmatrix} %d & %d \\\\ 0 & %d \\end{pmatrix} x\\), what is the larger eigenvalue \\(\\lambda_{\\max}\\)?", l1, b, l2),
			Answer:      fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("For an upper triangular matrix, the eigenvalues are the diagonal entries %d and %d. The larger eigenvalue is %d.", l1, l2, ans),
		}
	default:
		d := smallInt(-6, 6)
		for d == 0 {
			d = smallInt(-6, 6)
		}
		b := smallInt(-4, 4)
		c := smallInt(-4, 4)
		ans := -d
		return generator.Problem{
			Question:    fmt.Sprintf("For the planar system \\(x' = \\begin{pmatrix} k & %d \\\\ %d & %d \\end{pmatrix} x\\), find the parameter value \\(k\\) that yields zero trace \\(\\operatorname{tr}(A) = 0\\).", b, c, d),
			Answer:      fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("The trace is \\(\\operatorname{tr}(A) = k + (%d) = 0\\), which yields \\(k = %d\\).", d, ans),
		}
	}
}

type eulerGen struct{}

func (g *eulerGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*3)
	hChoices := []float64{0.1, 0.2, 0.5}
	h := hChoices[rand.Intn(len(hChoices))]
	x0 := 0.0
	y0 := float64(rand.Intn(max(1, scale*2)) + 1)
	// simple ODE y' = y, exact y = y0*e^x, Euler: y1 = y0 + h*y0
	y1 := y0 + h*y0
	_ = x0
	return generator.Problem{
		Question:    fmt.Sprintf("Use one step of Euler's method with \\(h=%.1f\\) for \\(y'=y\\), \\(y(0)=%.0f\\). Estimate \\(y(%.1f)\\).", h, y0, h),
		Answer:      fmt.Sprintf("%.4f", y1),
		Explanation: fmt.Sprintf("Euler: \\(y_1 = y_0 + h f(x_0,y_0) = %.0f + %.1f\\cdot%.0f = %.4f\\).", y0, h, y0, y1),
	}
}

type rk4Gen struct{}

func (g *rk4Gen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := rand.Intn(5)
	switch mode {
	case 0:
		// Stage evaluations across N steps
		n := smallInt(2, 6+int(ctx.Difficulty*10))
		ans := 4 * n
		return generator.Problem{
			Question: fmt.Sprintf("How many stage function evaluations does the classic 4th-order Runge-Kutta method (RK4) require across %d steps?", n),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf(
				"Each step of classic RK4 evaluates the slope at 4 stages (k_1, k_2, k_3 and k_4). For %d steps, the total number of evaluations is 4 * %d = %d.",
				n, n, ans,
			),
		}

	case 1:
		// Global error reduction factor under step size division
		kChoices := []int{2, 3, 4, 5}
		k := kChoices[rand.Intn(len(kChoices))]
		ans := k * k * k * k
		return generator.Problem{
			Question: fmt.Sprintf(
				"If the step size \\(h\\) in the fourth-order Runge-Kutta method (RK4) is reduced by a factor of %d, by what theoretical factor is the global truncation error reduced?",
				k,
			),
			Answer: fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf(
				"RK4 has a global truncation error of \\(\\mathcal{O}(h^4)\\). Dividing the step size by %d reduces the global error by a factor of \\(%d^4 = %d\\).",
				k, k, ans,
			),
		}

	case 2:
		// First stage slope k_1 = f(x_0, y_0)
		a := smallInt(1, 4)
		b := smallInt(1, 4)
		x0 := smallInt(0, 3)
		y0 := smallInt(1, 4)
		k1 := a*x0 + b*y0
		return generator.Problem{
			Question: fmt.Sprintf(
				"For the initial value problem \\(y' = %dx + %dy\\) with \\(y(%d) = %d\\), compute the first stage slope \\(k_1\\) in the classic RK4 method.",
				a, b, x0, y0,
			),
			Answer: fmt.Sprintf("%d", k1),
			Explanation: fmt.Sprintf(
				"In the classic RK4 method, the first stage slope is evaluated at the initial point: \\(k_1 = f(x_0, y_0) = %d(%d) + %d(%d) = %d\\).",
				a, x0, b, y0, k1,
			),
		}

	case 3:
		// Truncation error order
		if rand.Float64() < 0.5 {
			return generator.Problem{
				Question:    "What is the integer order \\(p\\) of the global truncation error \\(\\mathcal{O}(h^p)\\) for the classic Runge-Kutta method?",
				Answer:      "4",
				Explanation: "The classic fourth-order Runge-Kutta method achieves a global truncation error of order \\(\\mathcal{O}(h^4)\\).",
			}
		}
		return generator.Problem{
			Question:    "What is the integer order \\(p\\) of the local truncation error \\(\\mathcal{O}(h^p)\\) for the classic fourth-order Runge-Kutta method?",
			Answer:      "5",
			Explanation: "Classic RK4 has local truncation error \\(\\mathcal{O}(h^5)\\), which accumulates over \\(\\mathcal{O}(1/h)\\) steps to produce a global error of \\(\\mathcal{O}(h^4)\\).",
		}

	default:
		// Sum of intermediate stage weights in butcher increment
		return generator.Problem{
			Question:    "In the classic RK4 increment formula \\(\\frac{h}{6}(k_1 + w_2 k_2 + w_3 k_3 + k_4)\\), what is the sum of the inner weights \\(w_2 + w_3\\)?",
			Answer:      "4",
			Explanation: "The classic RK4 butcher tableau weights the mid-point slopes at 2 each (\\(w_2 = 2\\) and \\(w_3 = 2\\)), so \\(w_2 + w_3 = 2 + 2 = 4\\).",
		}
	}
}

type stabilityGen struct{}

func (g *stabilityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := rand.Intn(5)
	switch mode {
	case 0:
		// Real equilibria count for x' = mu - x^2
		mus := []int{-3, -2, -1, 1, 2, 4}
		mu := mus[rand.Intn(len(mus))]
		count := 0
		reason := fmt.Sprintf("Because %d < 0, no real solution exists and there are 0 real equilibria.", mu)
		if mu > 0 {
			count = 2
			reason = fmt.Sprintf("Because %d > 0, there are 2 real equilibria at \\(\\pm\\sqrt{%d}\\).", mu, mu)
		}
		return generator.Problem{
			Question: fmt.Sprintf("For the autonomous differential equation \\(x' = %d - x^2\\), how many real equilibrium points exist?", mu),
			Answer:   fmt.Sprintf("%d", count),
			Explanation: fmt.Sprintf(
				"Setting \\(x' = 0\\) gives \\(x^2 = %d\\). %s",
				mu, reason,
			),
		}

	case 1:
		// Value of non-zero equilibrium point
		a := smallInt(2, 6+int(ctx.Difficulty*4))
		return generator.Problem{
			Question: fmt.Sprintf("Find the non-zero equilibrium point \\(x^*\\) of the flow equation \\(x' = %dx - x^2\\).", a),
			Answer:   fmt.Sprintf("%d", a),
			Explanation: fmt.Sprintf(
				"Setting \\(x' = x(%d - x) = 0\\) yields equilibria at \\(x = 0\\) and \\(x = %d\\). The non-zero equilibrium is %d.",
				a, a, a,
			),
		}

	case 2:
		// Linearized derivative at stable equilibrium
		a := smallInt(2, 6+int(ctx.Difficulty*4))
		return generator.Problem{
			Question: fmt.Sprintf("For the flow \\(x' = %dx - x^2\\), evaluate the derivative \\(f'(x^*)\\) at the stable equilibrium \\(x^* = %d\\).", a, a),
			Answer:   fmt.Sprintf("%d", -a),
			Explanation: fmt.Sprintf(
				"The flow function is \\(f(x) = %dx - x^2\\) with derivative \\(f'(x) = %d - 2x\\). At \\(x^* = %d\\), \\(f'(%d) = %d - 2(%d) = %d\\). The negative derivative confirms asymptotic stability.",
				a, a, a, a, a, a, -a,
			),
		}

	case 3:
		// Number of asymptotically stable equilibria
		a := smallInt(2, 6+int(ctx.Difficulty*4))
		return generator.Problem{
			Question: fmt.Sprintf("How many asymptotically stable equilibrium points does the differential equation \\(x' = %dx - x^2\\) possess?", a),
			Answer:   "1",
			Explanation: fmt.Sprintf(
				"The equilibria are \\(x = 0\\) and \\(x = %d\\). The derivative \\(f'(x) = %d - 2x\\) gives \\(f'(0) = %d > 0\\) (unstable) and \\(f'(%d) = -%d < 0\\) (stable), so exactly 1 equilibrium is asymptotically stable.",
				a, a, a, a, a,
			),
		}

	default:
		// 2D linear system trace
		p := smallInt(1, 5)
		q := smallInt(1, 5)
		tr := -(p + q)
		return generator.Problem{
			Question: fmt.Sprintf("For the 2D linear system \\(x' = -%dx\\) and \\(y' = -%dy\\) governing decay to the origin, what is the trace of the system matrix \\(A\\)?", p, q),
			Answer:   fmt.Sprintf("%d", tr),
			Explanation: fmt.Sprintf(
				"The system matrix is diagonal with eigenvalues \\(\\lambda_1 = -%d\\) and \\(\\lambda_2 = -%d\\). The trace is the sum of the eigenvalues: \\(\\operatorname{tr}(A) = -%d + (-%d) = %d\\).",
				p, q, p, q, tr,
			),
		}
	}
}

type characteristicsGen struct{}

func (g *characteristicsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := rand.Intn(5)
	switch mode {
	case 0:
		// Propagation speed dx/dt for transport equation u_t + c u_x = 0
		c := smallInt(2, 6+int(ctx.Difficulty*4))
		return generator.Problem{
			Question: fmt.Sprintf("For the first-order linear PDE \\(u_t + %du_x = 0\\), what is the characteristic propagation speed \\(\\frac{dx}{dt}\\)?", c),
			Answer:   fmt.Sprintf("%d", c),
			Explanation: fmt.Sprintf(
				"The characteristic equations for the advection equation \\(u_t + %du_x = 0\\) are \\(\\frac{dt}{1} = \\frac{dx}{%d} = \\frac{du}{0}\\). Along the characteristic curves, \\(\\frac{dx}{dt} = %d\\).",
				c, c, c,
			),
		}

	case 1:
		// Evaluation of quadratic initial condition u(x, 0) = x^2
		c := smallInt(2, 4)
		t0 := smallInt(1, 3)
		delta := smallInt(1, 5)
		x0 := c*t0 + delta
		ans := delta * delta
		return generator.Problem{
			Question: fmt.Sprintf("Consider the transport equation \\(u_t + %du_x = 0\\) with initial profile \\(u(x, 0) = x^2\\). Evaluate \\(u(%d, %d)\\).", c, x0, t0),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf(
				"Along the characteristic curves \\(x - %dt = C\\), the solution is constant: \\(u(x, t) = u(x - %dt, 0) = (x - %dt)^2\\). Substituting \\((x, t) = (%d, %d)\\) gives \\(u(%d, %d) = (%d - %d \\cdot %d)^2 = %d^2 = %d\\).",
				c, c, c, x0, t0, x0, t0, x0, c, t0, delta, ans,
			),
		}

	case 2:
		// Linear initial condition u(x, 0) = m*x + k
		c := smallInt(2, 4)
		m := smallInt(2, 5)
		k := smallInt(1, 9)
		t0 := smallInt(1, 3)
		x0 := smallInt(2, 8) + c*t0
		xBase := x0 - c*t0
		ans := m*xBase + k
		return generator.Problem{
			Question: fmt.Sprintf("For the advection equation \\(u_t + %du_x = 0\\) with initial condition \\(u(x, 0) = %dx + %d\\), evaluate \\(u(%d, %d)\\).", c, m, k, x0, t0),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf(
				"The general solution is \\(u(x, t) = u(x - %dt, 0) = %d(x - %dt) + %d\\). Evaluating at \\((x, t) = (%d, %d)\\) yields \\(u(%d, %d) = %d(%d - %d \\cdot %d) + %d = %d(%d) + %d = %d\\).",
				c, m, c, k, x0, t0, x0, t0, m, x0, c, t0, k, m, xBase, k, ans,
			),
		}

	case 3:
		// Characteristic ray origin at t = 0
		c := smallInt(2, 5)
		t1 := smallInt(2, 5)
		xOrigin := smallInt(1, 10)
		x1 := xOrigin + c*t1
		return generator.Problem{
			Question: fmt.Sprintf("A characteristic curve for the PDE \\(u_t + %du_x = 0\\) passes through the point \\((x, t) = (%d, %d)\\). At what coordinate \\(x\\) does this characteristic curve intersect the initial axis \\(t = 0\\)?", c, x1, t1),
			Answer:   fmt.Sprintf("%d", xOrigin),
			Explanation: fmt.Sprintf(
				"The characteristic curves satisfy \\(\\frac{dx}{dt} = %d\\), which integrates to \\(x - %dt = x_0\\). Using the point \\((%d, %d)\\), the initial position is \\(x_0 = %d - %d \\cdot %d = %d\\).",
				c, c, x1, t1, x1, c, t1, xOrigin,
			),
		}

	default:
		// Inhomogeneous advection with zero initial state u(x, 0) = 0
		c := smallInt(2, 6)
		k := smallInt(2, 7)
		t0 := smallInt(2, 6)
		x0 := smallInt(1, 12)
		ans := k * t0
		return generator.Problem{
			Question: fmt.Sprintf("Consider the inhomogeneous transport equation \\(u_t + %du_x = %d\\) with initial condition \\(u(x, 0) = 0\\). What is the value of \\(u(%d, %d)\\)?", c, k, x0, t0),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf(
				"Along the characteristics \\(\\frac{dx}{dt} = %d\\), the transport equation simplifies to \\(\\frac{du}{dt} = %d\\). Integrating from \\(t = 0\\) with \\(u(x - %dt, 0) = 0\\) yields \\(u(x, t) = %dt\\). Evaluating at \\(t = %d\\) gives \\(u(%d, %d) = %d \\cdot %d = %d\\).",
				c, k, c, k, t0, x0, t0, k, t0, ans,
			),
		}
	}
}

type fourierGen struct{}

func (g *fourierGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := rand.Intn(5)
	switch mode {
	case 0:
		// Algebraic multiplier evaluation at specific angular frequency omega
		k := smallInt(12, 30)
		w := smallInt(1, 4)
		pw := k - w*w
		return generator.Problem{
			Question: fmt.Sprintf("Applying the Fourier transform to \\(y'' + %dy = f(t)\\) produces the algebraic relation \\(P(\\omega) \\hat{y}(\\omega) = \\hat{f}(\\omega)\\), where \\(P(\\omega) = %d - \\omega^2\\). Evaluate \\(P(%d)\\).", k, k, w),
			Answer:   fmt.Sprintf("%d", pw),
			Explanation: fmt.Sprintf(
				"The Fourier transform replaces \\(\\frac{d^2}{dt^2}\\) with \\((i\\omega)^2 = -\\omega^2\\). Therefore \\(P(\\omega) = %d - \\omega^2\\). At \\(\\omega = %d\\), \\(P(%d) = %d - %d^2 = %d - %d = %d\\).",
				k, w, w, k, w, k, w*w, pw,
			),
		}

	case 1:
		// Squared transfer denominator for first-order system y' + a*y = f(t)
		a := smallInt(2, 6)
		w := smallInt(1, 8)
		denom := w*w + a*a
		return generator.Problem{
			Question: fmt.Sprintf("For the relaxation system \\(y' + %dy = f(t)\\), the squared denominator of the frequency response function is \\(D(\\omega) = \\omega^2 + %d^2\\). Evaluate \\(D(%d)\\).", a, a, w),
			Answer:   fmt.Sprintf("%d", denom),
			Explanation: fmt.Sprintf(
				"The Fourier transform of the operator is \\(i\\omega + %d\\). The squared denominator is \\(|i\\omega + %d|^2 = \\omega^2 + %d^2 = \\omega^2 + %d\\). Evaluating at \\(\\omega = %d\\) gives \\(%d^2 + %d = %d + %d = %d\\).",
				a, a, a, a*a, w, w, a*a, w*w, a*a, denom,
			),
		}

	case 2:
		// Resonant frequency of undamped oscillator
		w0 := smallInt(2, 9)
		k := w0 * w0
		return generator.Problem{
			Question: fmt.Sprintf("At what positive angular frequency \\(\\omega > 0\\) does the undamped oscillator \\(y'' + %dy = f(t)\\) experience resonance?", k),
			Answer:   fmt.Sprintf("%d", w0),
			Explanation: fmt.Sprintf(
				"The Fourier transform gives \\((-\\omega^2 + %d)\\hat{y}(\\omega) = \\hat{f}(\\omega)\\). Resonance occurs when the transfer denominator vanishes: \\(\\omega^2 = %d\\), giving \\(\\omega = %d\\).",
				k, k, w0,
			),
		}

	case 3:
		// Fourth-order operator frequency factor
		c := smallInt(1, 20)
		w := smallInt(1, 3)
		w4 := w * w * w * w
		ans := w4 + c
		return generator.Problem{
			Question: fmt.Sprintf("For the fourth-order differential equation \\(y^{(4)} + %dy = f(t)\\), the Fourier transform yields \\((\\omega^4 + %d)\\hat{y}(\\omega) = \\hat{f}(\\omega)\\). Evaluate \\(\\omega^4 + %d\\) at \\(\\omega = %d\\).", c, c, c, w),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf(
				"Because \\((i\\omega)^4 = i^4 \\omega^4 = \\omega^4\\), the operator multiplier is \\(\\omega^4 + %d\\). Substituting \\(\\omega = %d\\) gives \\(%d^4 + %d = %d + %d = %d\\).",
				c, w, w, c, w4, c, ans,
			),
		}

	default:
		// Plancherel total energy
		energy := smallInt(5, 50)
		return generator.Problem{
			Question: fmt.Sprintf("A signal \\(y(t)\\) has total time-domain energy \\(\\int_{-\\infty}^\\infty |y(t)|^2 dt = %d\\). What is the total spectral energy \\(\\frac{1}{2\\pi}\\int_{-\\infty}^\\infty |\\hat{y}(\\omega)|^2 d\\omega\\)?", energy),
			Answer:   fmt.Sprintf("%d", energy),
			Explanation: fmt.Sprintf(
				"By Plancherel's theorem for the Fourier transform, the total energy in the time domain is identical to the total energy in the frequency domain: \\(\\frac{1}{2\\pi}\\int_{-\\infty}^\\infty |\\hat{y}(\\omega)|^2 d\\omega = \\int_{-\\infty}^\\infty |y(t)|^2 dt = %d\\).",
				energy,
			),
		}
	}
}

type bifurcationGen struct{}

func (g *bifurcationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := rand.Intn(5)
	switch mode {
	case 0:
		// Real equilibria count for pitchfork bifurcation x' = mu*x - x^3
		mus := []int{-4, -3, -2, -1, 1, 2, 3, 5}
		mu := mus[rand.Intn(len(mus))]
		count := 1
		reason := fmt.Sprintf("Because \\(\\mu = %d < 0\\), the equation \\(x(\\mu - x^2) = 0\\) has only the origin \\(x = 0\\) as a real solution (1 real equilibrium).", mu)
		if mu > 0 {
			count = 3
			reason = fmt.Sprintf("Because \\(\\mu = %d > 0\\), the equation \\(x(\\mu - x^2) = 0\\) yields three real equilibria at \\(x = 0\\) and \\(x = \\pm\\sqrt{%d}\\).", mu, mu)
		}
		return generator.Problem{
			Question: fmt.Sprintf("For the pitchfork bifurcation system \\(x' = %dx - x^3\\), how many real equilibrium points exist?", mu),
			Answer:   fmt.Sprintf("%d", count),
			Explanation: fmt.Sprintf(
				"Setting \\(x' = x(%d - x^2) = 0\\) gives equilibria. %s",
				mu, reason,
			),
		}

	case 1:
		// Non-zero equilibrium point in transcritical bifurcation x' = ax - x^2
		a := smallInt(2, 6+int(ctx.Difficulty*4))
		return generator.Problem{
			Question: fmt.Sprintf("In the transcritical bifurcation system \\(x' = %dx - x^2\\), find the value of the non-zero equilibrium point \\(x^*\\).", a),
			Answer:   fmt.Sprintf("%d", a),
			Explanation: fmt.Sprintf(
				"Factoring \\(x' = x(%d - x) = 0\\) gives the trivial equilibrium \\(x = 0\\) and the non-zero branch \\(x^* = %d\\).",
				a, a,
			),
		}

	case 2:
		// Linearized derivative at the non-zero equilibrium of transcritical flow
		a := smallInt(2, 6+int(ctx.Difficulty*4))
		return generator.Problem{
			Question: fmt.Sprintf("For the transcritical flow \\(x' = %dx - x^2\\), evaluate the linearized derivative \\(f'(x^*)\\) at the non-zero equilibrium \\(x^* = %d\\).", a, a),
			Answer:   fmt.Sprintf("%d", -a),
			Explanation: fmt.Sprintf(
				"The flow function \\(f(x) = %dx - x^2\\) has derivative \\(f'(x) = %d - 2x\\). At \\(x^* = %d\\), \\(f'(%d) = %d - 2(%d) = %d\\). The negative derivative indicates exchange of stability.",
				a, a, a, a, a, a, -a,
			),
		}

	case 3:
		// Critical saddle-node parameter r_c for x' = r + 2kx - x^2
		k := smallInt(1, 5)
		twoK := 2 * k
		rc := -(k * k)
		return generator.Problem{
			Question: fmt.Sprintf("At what critical parameter value \\(r_c\\) does a saddle-node bifurcation occur for \\(x' = r + %dx - x^2\\)?", twoK),
			Answer:   fmt.Sprintf("%d", rc),
			Explanation: fmt.Sprintf(
				"Rewriting \\(f(x) = r + %dx - x^2 = r + %d - (x - %d)^2\\), equilibria satisfy \\((x - %d)^2 = r + %d\\). Real solutions appear when the discriminant \\(r + %d = 0\\), yielding critical parameter \\(r_c = %d\\).",
				twoK, k*k, k, k, k*k, k*k, rc,
			),
		}

	default:
		// Squared amplitude of non-trivial equilibria in pitchfork bifurcation
		k := smallInt(2, 9)
		return generator.Problem{
			Question: fmt.Sprintf("In the supercritical pitchfork bifurcation \\(x' = %dx - x^3\\), what is the squared magnitude \\((x^*)^2\\) of the non-trivial equilibrium points?", k),
			Answer:   fmt.Sprintf("%d", k),
			Explanation: fmt.Sprintf(
				"Setting \\(x' = x(%d - x^2) = 0\\) for non-zero equilibria gives \\(x^2 = %d\\). Therefore \\((x^*)^2 = %d\\).",
				k, k, k,
			),
		}
	}
}

type waveGen struct{}

func (g *waveGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := rand.Intn(5)
	switch mode {
	case 0:
		// Propagation speed c from u_tt = c^2 u_xx
		c := smallInt(2, 6+int(ctx.Difficulty*4))
		c2 := c * c
		return generator.Problem{
			Question: fmt.Sprintf("For the one-dimensional wave equation \\(u_{tt} = %du_{xx}\\), what is the propagation speed \\(c > 0\\) of traveling wave disturbances?", c2),
			Answer:   fmt.Sprintf("%d", c),
			Explanation: fmt.Sprintf(
				"The wave equation is written as \\(u_{tt} = c^2 u_{xx}\\). Setting \\(c^2 = %d\\) gives the positive wave speed \\(c = \\sqrt{%d} = %d\\).",
				c2, c2, c,
			),
		}

	case 1:
		// d'Alembert with quadratic initial displacement u(x, 0) = x^2 and u_t(x, 0) = 0
		c := smallInt(1, 3)
		x0 := smallInt(1, 5)
		t0 := smallInt(1, 4)
		ans := x0*x0 + c*c*t0*t0
		return generator.Problem{
			Question: fmt.Sprintf("Consider the wave equation \\(u_{tt} = %du_{xx}\\) with initial conditions \\(u(x, 0) = x^2\\) and \\(u_t(x, 0) = 0\\). Evaluate \\(u(%d, %d)\\).", c*c, x0, t0),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf(
				"By d'Alembert's formula with zero initial velocity, \\(u(x, t) = \\frac{1}{2}[(x - %dt)^2 + (x + %dt)^2] = x^2 + %dt^2\\). Substituting \\((x, t) = (%d, %d)\\) gives \\(u(%d, %d) = %d^2 + %d(%d^2) = %d + %d = %d\\).",
				c, c, c*c, x0, t0, x0, t0, x0, c*c, t0, x0*x0, c*c*t0*t0, ans,
			),
		}

	case 2:
		// Constant initial velocity u_t(x, 0) = v0 and u(x, 0) = 0
		c := smallInt(1, 4)
		v0 := smallInt(2, 7)
		x0 := smallInt(1, 8)
		t0 := smallInt(1, 5)
		ans := v0 * t0
		return generator.Problem{
			Question: fmt.Sprintf("Consider the wave equation \\(u_{tt} = %du_{xx}\\) with initial conditions \\(u(x, 0) = 0\\) and constant initial velocity \\(u_t(x, 0) = %d\\). Evaluate \\(u(%d, %d)\\).", c*c, v0, x0, t0),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf(
				"Applying d'Alembert's velocity integral: \\(u(x, t) = \\frac{1}{2(%d)}\\int_{x - %dt}^{x + %dt} %d\\,ds = \\frac{%d}{2(%d)}(2 \\cdot %dt) = %dt\\). At \\(t = %d\\), \\(u(%d, %d) = %d \\cdot %d = %d\\).",
				c, c, c, v0, v0, c, c, v0, t0, x0, t0, v0, t0, ans,
			),
		}

	case 3:
		// Length of domain of dependence interval [x - ct, x + ct]
		c := smallInt(2, 5)
		t0 := smallInt(2, 5)
		length := 2 * c * t0
		return generator.Problem{
			Question: fmt.Sprintf("For the wave equation \\(u_{tt} = %du_{xx}\\), what is the length of the domain of dependence interval \\([x - ct, x + ct]\\) on the initial axis at time \\(t = %d\\)?", c*c, t0),
			Answer:   fmt.Sprintf("%d", length),
			Explanation: fmt.Sprintf(
				"The domain of dependence on the line \\(t = 0\\) for a point \\((x, t)\\) is the interval \\([x - ct, x + ct]\\). With wave velocity \\(c = %d\\) and time \\(t = %d\\), the length of this interval is \\(2ct = 2(%d)(%d) = %d\\).",
				c, t0, c, t0, length,
			),
		}

	default:
		// Linear initial displacement u(x, 0) = mx + k and u_t(x, 0) = 0
		c := smallInt(1, 4)
		m := smallInt(2, 5)
		k := smallInt(1, 9)
		x0 := smallInt(1, 8)
		t0 := smallInt(1, 4)
		ans := m*x0 + k
		return generator.Problem{
			Question: fmt.Sprintf("Consider the wave equation \\(u_{tt} = %du_{xx}\\) with initial conditions \\(u(x, 0) = %dx + %d\\) and \\(u_t(x, 0) = 0\\). Evaluate \\(u(%d, %d)\\).", c*c, m, k, x0, t0),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf(
				"Using d'Alembert's formula with zero initial velocity gives \\(u(x, t) = \\frac{1}{2}[(%d(x - %dt) + %d) + (%d(x + %dt) + %d)] = %dx + %d\\). At \\((x, t) = (%d, %d)\\), \\(u(%d, %d) = %d(%d) + %d = %d\\).",
				m, c, k, m, c, k, m, k, x0, t0, x0, t0, m, x0, k, ans,
			),
		}
	}
}

type heatGen struct{}

func (g *heatGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := rand.Intn(5)
	switch mode {
	case 0:
		// Spatial Dirichlet eigenvalue lambda_n on [0, pi]
		n := smallInt(1, 5+int(ctx.Difficulty*4))
		ans := n * n
		return generator.Problem{
			Question: fmt.Sprintf("For the heat equation on the spatial domain \\([0, \\pi]\\) with Dirichlet boundary conditions \\(u(0, t) = u(\\pi, t) = 0\\), what is the spatial eigenvalue \\(\\lambda_%d\\) corresponding to eigenfunction \\(\\sin(%dx)\\)?", n, n),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf(
				"The spatial eigenvalue problem is \\(X'' + \\lambda X = 0\\) with \\(X(0) = X(\\pi) = 0\\). The eigenfunctions are \\(X_n(x) = \\sin(nx)\\) with eigenvalues \\(\\lambda_n = n^2\\). For \\(n = %d\\), \\(\\lambda_%d = %d^2 = %d\\).",
				n, n, n, ans,
			),
		}

	case 1:
		// Temporal decay rate beta = alpha * lambda_n
		alpha := smallInt(1, 5)
		n := smallInt(1, 4)
		lambda := n * n
		beta := alpha * lambda
		return generator.Problem{
			Question: fmt.Sprintf("For the diffusion equation \\(u_t = %du_{xx}\\) on \\([0, \\pi]\\) with Dirichlet ends, what is the exponential decay rate \\(\\beta\\) in the temporal factor \\(e^{-\\beta t}\\) for harmonic mode \\(n = %d\\)?", alpha, n),
			Answer:   fmt.Sprintf("%d", beta),
			Explanation: fmt.Sprintf(
				"Separation of variables gives \\(T'(t) = -\\alpha \\lambda_n T(t)\\). On \\([0, \\pi]\\), the \\(n = %d\\) spatial eigenvalue is \\(\\lambda_%d = %d^2 = %d\\). The decay rate is \\(\\beta = %d \\cdot %d = %d\\).",
				n, n, n, lambda, alpha, lambda, beta,
			),
		}

	case 2:
		// Steady-state linear temperature profile
		L := smallInt(2, 8)
		T0 := smallInt(0, 30)
		s := smallInt(2, 6)
		TL := T0 + s*L
		x0 := smallInt(1, L-1)
		ans := T0 + s*x0
		return generator.Problem{
			Question: fmt.Sprintf("A conductive rod of length \\(L = %d\\) satisfies the steady-state equation \\(u_{xx} = 0\\) with boundary temperatures \\(u(0) = %d\\) and \\(u(%d) = %d\\). Evaluate the temperature \\(u(%d)\\).", L, T0, L, TL, x0),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf(
				"Steady-state heat distribution satisfies \\(u_{xx} = 0\\), producing the linear profile \\(u(x) = T_0 + \\frac{T_L - T_0}{L}x = %d + \\frac{%d - %d}{%d}x = %d + %dx\\). At \\(x = %d\\), \\(u(%d) = %d + %d(%d) = %d\\).",
				T0, TL, T0, L, T0, s, x0, x0, T0, s, x0, ans,
			),
		}

	case 3:
		// Count of interior nodal points
		n := smallInt(2, 8)
		nodalCount := n - 1
		return generator.Problem{
			Question: fmt.Sprintf("How many interior nodal points (coordinates where \\(u(x, t) = 0\\) strictly inside the open interval \\(0 < x < \\pi\\)) does the spatial harmonic \\(\\sin(%dx)\\) possess?", n),
			Answer:   fmt.Sprintf("%d", nodalCount),
			Explanation: fmt.Sprintf(
				"The eigenfunction \\(\\sin(%dx) = 0\\) occurs when \\(%dx = k\\pi\\) for integer \\(k\\). Inside the open domain \\(0 < x < \\pi\\), valid indices are \\(k = 1, 2, \\dots, %d\\). There are %d interior nodal points.",
				n, n, nodalCount, nodalCount,
			),
		}

	default:
		// Maximum principle bound
		L := smallInt(2, 8)
		tMax := smallInt(5, 15)
		M := smallInt(40, 95)
		return generator.Problem{
			Question: fmt.Sprintf("A heat diffusion process on the region \\([0, %d] \\times [0, %d]\\) has temperatures bounded by \\(M = %d\\) along the initial and spatial boundaries. By the maximum principle, what is the maximum possible temperature in the interior?", L, tMax, M),
			Answer:   fmt.Sprintf("%d", M),
			Explanation: fmt.Sprintf(
				"By the weak maximum principle for parabolic equations, the maximum of \\(u(x, t)\\) across a space-time domain is achieved on the parabolic boundary. Therefore the maximum interior temperature cannot exceed the boundary upper bound of %d.",
				M,
			),
		}
	}
}

type greenFunctionGen struct{}

func (g *greenFunctionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := rand.Intn(5)
	switch mode {
	case 0:
		// Scaled Green's function L * G(x, xi) = min(x, xi) * (L - max(x, xi))
		L := smallInt(6, 12)
		x := smallInt(1, L-1)
		xi := smallInt(1, L-1)
		mn := x
		mx := xi
		if xi < x {
			mn = xi
			mx = x
		}
		ans := mn * (L - mx)
		return generator.Problem{
			Question: fmt.Sprintf("For the boundary value problem \\(-y'' = f(x)\\) on \\([0, %d]\\) with Dirichlet conditions \\(y(0) = y(%d) = 0\\), the Green's function satisfies \\(L \\cdot G(x, \\xi) = \\min(x, \\xi)(L - \\max(x, \\xi))\\). Evaluate the integer quantity \\(%d \\cdot G(%d, %d)\\).", L, L, L, x, xi),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf(
				"For \\(-y'' = f(x)\\) on \\([0, %d]\\), \\(G(x, \\xi) = \\frac{\\min(x, \\xi)(%d - \\max(x, \\xi))}{%d}\\). With \\(x = %d\\) and \\(\\xi = %d\\), the minimum is %d and the maximum is %d. Thus \\(%d \\cdot G(%d, %d) = %d(%d - %d) = %d \\cdot %d = %d\\).",
				L, L, L, x, xi, mn, mx, L, x, xi, mn, L, mx, mn, L-mx, ans,
			),
		}

	case 1:
		// Derivative jump discontinuity magnitude
		return generator.Problem{
			Question:    "For the differential operator \\(L[y] = y''\\) governed by \\(y''(x) = \\delta(x - \\xi)\\), what is the jump discontinuity magnitude \\(\\left|\\left.\\frac{dG}{dx}\\right|_{x=\\xi^+} - \\left.\\frac{dG}{dx}\\right|_{x=\\xi^-}\\right|\\) across the point source \\(x = \\xi\\)?",
			Answer:      "1",
			Explanation: "Integrating the equation \\(G''(x, \\xi) = \\delta(x - \\xi)\\) over the infinitesimal interval \\((\\xi - \\epsilon, \\xi + \\epsilon)\\) yields \\(G'(\\xi^+) - G'(\\xi^-) = \\int_{\\xi-\\epsilon}^{\\xi+\\epsilon} \\delta(x-\\xi)\\,dx = 1\\). The jump magnitude is exactly 1.",
		}

	case 2:
		// Midpoint deflection under constant load: y(m) = c * m^2
		m := smallInt(1, 4)
		L := 2 * m
		c := smallInt(1, 5)
		k := 2 * c
		ans := c * m * m
		return generator.Problem{
			Question: fmt.Sprintf("A structural rod governed by \\(-y'' = %d\\) on \\([0, %d]\\) with boundary values \\(y(0) = y(%d) = 0\\) has deflection given by \\(y(x) = \\int_0^{%d} G(x, \\xi)(%d)\\,d\\xi\\). Evaluate the midpoint deflection \\(y(%d)\\).", k, L, L, L, k, m),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf(
				"Integrating \\(-y'' = %d\\) twice with \\(y(0) = y(%d) = 0\\) yields \\(y(x) = \\frac{%d}{2}x(%d - x) = %d x(%d - x)\\). Evaluating at the midpoint \\(x = %d\\) gives \\(y(%d) = %d(%d)(%d - %d) = %d(%d)(%d) = %d\\).",
				k, L, k, L, c, L, m, m, c, m, L, m, c, m, m, ans,
			),
		}

	case 3:
		// First-order causal impulse jump
		a := smallInt(1, 5)
		k := smallInt(2, 9)
		t0 := smallInt(1, 5)
		return generator.Problem{
			Question: fmt.Sprintf("For the relaxation system \\(y' + %dy = %d \\delta(t - %d)\\) with initial rest state \\(y(0) = 0\\), what is the response value \\(y(%d^+)\\) immediately after the impulse?", a, k, t0, t0),
			Answer:   fmt.Sprintf("%d", k),
			Explanation: fmt.Sprintf(
				"Integrating \\(y' + %dy = %d\\delta(t - %d)\\) across the impulse yields a jump of \\(y(%d^+) - y(%d^-) = %d\\). Because the system is initially at rest with \\(y(%d^-) = 0\\), the post-impulse state is \\(y(%d^+) = %d\\).",
				a, k, t0, t0, t0, k, t0, t0, k,
			),
		}

	default:
		// Self-adjoint reciprocity relation
		x1 := smallInt(1, 4)
		xi1 := smallInt(5, 9)
		val := smallInt(2, 25)
		return generator.Problem{
			Question: fmt.Sprintf("A self-adjoint boundary value problem has symmetric Green's function \\(G(x, \\xi)\\). If \\(G(%d, %d) = %d\\), what is the value of \\(G(%d, %d)\\)?", x1, xi1, val, xi1, x1),
			Answer:   fmt.Sprintf("%d", val),
			Explanation: fmt.Sprintf(
				"By the reciprocity theorem for self-adjoint differential operators, the Green's function is symmetric: \\(G(x, \\xi) = G(\\xi, x)\\). Therefore \\(G(%d, %d) = G(%d, %d) = %d\\).",
				xi1, x1, x1, xi1, val,
			),
		}
	}
}

type sturmLiouvilleGen struct{}

func (g *sturmLiouvilleGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := rand.Intn(5)
	switch mode {
	case 0:
		// Dirichlet eigenvalues lambda_n = n^2 on [0, pi]
		n := smallInt(1, 5+int(ctx.Difficulty*4))
		ans := n * n
		return generator.Problem{
			Question: fmt.Sprintf("For the regular Sturm-Liouville eigenvalue problem \\(-y'' = \\lambda y\\) on \\([0, \\pi]\\) with Dirichlet boundary conditions \\(y(0) = y(\\pi) = 0\\), what is the eigenvalue \\(\\lambda_%d\\)?", n),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf(
				"The eigenfunctions are \\(y_n(x) = \\sin(nx)\\), yielding eigenvalues \\(\\lambda_n = n^2\\). For harmonic mode \\(n = %d\\), \\(\\lambda_%d = %d^2 = %d\\).",
				n, n, n, ans,
			),
		}

	case 1:
		// Neumann boundary condition principal eigenvalue lambda_0 = 0
		return generator.Problem{
			Question:    "For the regular Sturm-Liouville problem \\(-y'' = \\lambda y\\) on \\([0, \\pi]\\) with homogeneous Neumann boundary conditions \\(y'(0) = y'(\\pi) = 0\\), what is the principal (smallest) eigenvalue \\(\\lambda_0\\)?",
			Answer:      "0",
			Explanation: "With Neumann boundary conditions \\(y'(0) = y'(\\pi) = 0\\), the constant function \\(y_0(x) = 1\\) satisfies \\(-y_0'' = 0 = 0 \\cdot y_0\\), giving the smallest eigenvalue \\(\\lambda_0 = 0\\).",
		}

	case 2:
		// Orthogonality integral between distinct eigenfunctions
		m := smallInt(1, 4)
		n := m + smallInt(1, 4)
		return generator.Problem{
			Question: fmt.Sprintf("Two distinct eigenfunctions \\(y_{%d}(x) = \\sin(%dx)\\) and \\(y_{%d}(x) = \\sin(%dx)\\) solve the Sturm-Liouville problem on \\([0, \\pi]\\) with unit weight. Evaluate the inner product integral \\(\\int_0^\\pi \\sin(%dx)\\sin(%dx)\\,dx\\).", m, m, n, n, m, n),
			Answer:   "0",
			Explanation: fmt.Sprintf(
				"Eigenfunctions corresponding to distinct eigenvalues of a regular self-adjoint Sturm-Liouville problem are strictly orthogonal with respect to the weight function: \\(\\int_0^\\pi \\sin(%dx)\\sin(%dx)\\,dx = 0\\).",
				m, n,
			),
		}

	case 3:
		// Sturm oscillation theorem: interior zeros count n - 1
		n := smallInt(2, 9)
		zeros := n - 1
		return generator.Problem{
			Question: fmt.Sprintf("By the Sturm oscillation theorem, how many interior zeros in the open interval \\((0, \\pi)\\) does the \\(n = %d\\)-th eigenfunction of a regular Sturm-Liouville problem possess?", n),
			Answer:   fmt.Sprintf("%d", zeros),
			Explanation: fmt.Sprintf(
				"The Sturm oscillation theorem establishes that the eigenfunction corresponding to the \\(n\\)-th eigenvalue (ordered \\(\\lambda_1 < \\lambda_2 < \\dots\\)) has exactly \\(n - 1\\) zeros in the interior of the domain. For \\(n = %d\\), there are \\(%d - 1 = %d\\) interior zeros.",
				n, n, zeros,
			),
		}

	default:
		// Scaled interval Dirichlet eigenvalue lambda_1 = k^2
		k := smallInt(2, 6)
		ans := k * k
		return generator.Problem{
			Question: fmt.Sprintf("For the Dirichlet Sturm-Liouville problem \\(-y'' = \\lambda y\\) on the interval \\([0, \\pi/%d]\\) with \\(y(0) = y(\\pi/%d) = 0\\), what is the first eigenvalue \\(\\lambda_1\\)?", k, k),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf(
				"The boundary condition \\(y(\\pi/%d) = \\sin(\\sqrt{\\lambda}\\pi/%d) = 0\\) requires \\(\\sqrt{\\lambda}\\pi/%d = n\\pi\\), giving \\(\\sqrt{\\lambda} = n(%d)\\). For \\(n = 1\\), \\(\\lambda_1 = %d^2 = %d\\).",
				k, k, k, k, k, ans,
			),
		}
	}
}

type boundaryValueGen struct{}

func (g *boundaryValueGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := rand.Intn(5)
	switch mode {
	case 0:
		// Linear BVP y'' = 0 on [0, L] with y(0) = y0, y(L) = yL
		L := smallInt(2, 6)
		y0 := smallInt(1, 10)
		m := smallInt(1, 5)
		yL := y0 + m*L
		x0 := smallInt(1, L-1)
		ans := y0 + m*x0
		return generator.Problem{
			Question: fmt.Sprintf("Consider the two-point boundary value problem \\(y'' = 0\\) on \\([0, %d]\\) with boundary values \\(y(0) = %d\\) and \\(y(%d) = %d\\). Evaluate \\(y(%d)\\).", L, y0, L, yL, x0),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf(
				"Integrating \\(y'' = 0\\) twice yields the linear solution \\(y(x) = C_1 x + C_2\\). From \\(y(0) = %d\\), \\(C_2 = %d\\). From \\(y(%d) = %d\\), \\(C_1(%d) + %d = %d\\), giving \\(C_1 = %d\\). Thus \\(y(x) = %dx + %d\\). At \\(x = %d\\), \\(y(%d) = %d(%d) + %d = %d\\).",
				y0, y0, L, yL, L, y0, yL, m, m, y0, x0, x0, m, x0, y0, ans,
			),
		}

	case 1:
		// Nonhomogeneous BVP y'' = 2c on [0, 2m] with y(0) = y(2m) = 0
		m := smallInt(1, 4)
		L := 2 * m
		c := smallInt(1, 5)
		k := 2 * c
		ans := -c * m * m
		return generator.Problem{
			Question: fmt.Sprintf("For the two-point boundary value problem \\(y'' = %d\\) on \\([0, %d]\\) with Dirichlet conditions \\(y(0) = y(%d) = 0\\), evaluate the midpoint value \\(y(%d)\\).", k, L, L, m),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf(
				"Integrating \\(y'' = %d\\) gives \\(y(x) = %d x^2 + C_1 x + C_2\\). From \\(y(0) = 0\\), \\(C_2 = 0\\). From \\(y(%d) = 0\\), \\(%d(%d^2) + C_1(%d) = 0\\), yielding \\(C_1 = -%d(%d) = -%d\\). Thus \\(y(x) = %d x(x - %d)\\). At \\(x = %d\\), \\(y(%d) = %d(%d)(%d - %d) = %d\\).",
				k, c, L, c, L, L, c, L, c*L, c, L, m, m, c, m, m, L, ans,
			),
		}

	case 2:
		// Resonant frequency non-uniqueness mode omega = n
		n := smallInt(1, 8)
		return generator.Problem{
			Question: fmt.Sprintf("For what positive integer \\(\\omega > 0\\) does the boundary value problem \\(y'' + \\omega^2 y = 0\\) on \\([0, \\pi]\\) with \\(y(0) = y(\\pi) = 0\\) possess non-trivial solutions for harmonic mode \\(n = %d\\)?", n),
			Answer:   fmt.Sprintf("%d", n),
			Explanation: fmt.Sprintf(
				"The general solution is \\(y(x) = A\\cos(\\omega x) + B\\sin(\\omega x)\\). The boundary condition \\(y(0) = 0\\) sets \\(A = 0\\). Then \\(y(\\pi) = B\\sin(\\omega\\pi) = 0\\) requires \\(\\omega\\pi = n\\pi\\) for non-zero solutions. Thus \\(\\omega = n = %d\\).",
				n,
			),
		}

	case 3:
		// Mixed Dirichlet/Neumann BVP y(0) = y0 and y'(L) = s
		s := smallInt(2, 6)
		y0 := smallInt(1, 9)
		L := smallInt(3, 8)
		x0 := smallInt(1, L-1)
		ans := s*x0 + y0
		return generator.Problem{
			Question: fmt.Sprintf("Consider the mixed boundary value problem \\(y'' = 0\\) on \\([0, %d]\\) with Dirichlet condition \\(y(0) = %d\\) and Neumann condition \\(y'(%d) = %d\\). Evaluate \\(y(%d)\\).", L, y0, L, s, x0),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf(
				"Integrating \\(y'' = 0\\) gives \\(y'(x) = C_1\\). From \\(y'(%d) = %d\\), \\(C_1 = %d\\). Integrating again gives \\(y(x) = %dx + C_2\\). With \\(y(0) = %d\\), \\(C_2 = %d\\). Thus \\(y(x) = %dx + %d\\). At \\(x = %d\\), \\(y(%d) = %d(%d) + %d = %d\\).",
				L, s, s, s, y0, y0, s, y0, x0, x0, s, x0, y0, ans,
			),
		}

	default:
		// Positive definite BVP solution count = 1
		k2 := smallInt(1, 9)
		return generator.Problem{
			Question: fmt.Sprintf("For the homogeneous boundary value problem \\(y'' - %dy = 0\\) on \\([0, 1]\\) with boundary conditions \\(y(0) = 0\\) and \\(y(1) = 0\\), how many solutions exist?", k2),
			Answer:   "1",
			Explanation: fmt.Sprintf(
				"Because the operator \\(-y'' + %dy\\) is strictly positive definite, the boundary value problem has trivial kernel. The only solution is the zero solution \\(y(x) = 0\\), meaning exactly 1 solution exists.",
				k2,
			),
		}
	}
}

type lyapunovGen struct{}

func (g *lyapunovGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := rand.Intn(5)
	switch mode {
	case 0:
		// 1D orbital derivative dot(V) = -2 k x^2
		k := smallInt(1, 5)
		x0 := smallInt(1, 4)
		ans := -2 * k * x0 * x0
		return generator.Problem{
			Question: fmt.Sprintf("For the dynamical system \\(x' = -%dx\\) with candidate Lyapunov function \\(V(x) = x^2\\), evaluate the orbital time derivative \\(\\dot{V}(%d)\\).", k, x0),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf(
				"The orbital derivative is \\(\\dot{V}(x) = \\frac{dV}{dx}x' = 2x(-%dx) = -%dx^2\\). Evaluating at \\(x = %d\\) yields \\(\\dot{V}(%d) = -%d(%d^2) = -%d(%d) = %d\\).",
				k, 2*k, x0, x0, 2*k, x0, 2*k, x0*x0, ans,
			),
		}

	case 1:
		// 2D linear system orbital derivative dot(V)
		a := smallInt(1, 4)
		b := smallInt(1, 4)
		x0 := smallInt(1, 3)
		y0 := smallInt(1, 3)
		ans := -2*a*x0*x0 - 2*b*y0*y0
		return generator.Problem{
			Question: fmt.Sprintf("For the 2D linear system \\(x' = -%dx\\) and \\(y' = -%dy\\) with Lyapunov function \\(V(x, y) = x^2 + y^2\\), evaluate the orbital time derivative \\(\\dot{V}(%d, %d)\\).", a, b, x0, y0),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf(
				"The orbital time derivative is \\(\\dot{V} = \\nabla V \\cdot (x', y') = 2x(-%dx) + 2y(-%dy) = -%dx^2 - %dy^2\\). Substituting \\((x, y) = (%d, %d)\\) gives \\(\\dot{V}(%d, %d) = -%d(%d^2) - %d(%d^2) = -%d(%d) - %d(%d) = %d\\).",
				a, b, 2*a, 2*b, x0, y0, x0, y0, 2*a, x0, 2*b, y0, 2*a, x0*x0, 2*b, y0*y0, ans,
			),
		}

	case 2:
		// Basin radius r where dot(V) < 0 for x' = -a*x + b*x^3
		b := smallInt(1, 3)
		r := smallInt(2, 5)
		a := b * r * r
		return generator.Problem{
			Question: fmt.Sprintf("For the nonlinear system \\(x' = -%dx + %dx^3\\) with Lyapunov function \\(V(x) = x^2\\), the origin is strictly asymptotically stable on the symmetric open interval \\((-r, r)\\) where \\(\\dot{V}(x) < 0\\) for \\(x \\ne 0\\). What is the positive radius \\(r > 0\\)?", a, b),
			Answer:   fmt.Sprintf("%d", r),
			Explanation: fmt.Sprintf(
				"The orbital derivative is \\(\\dot{V}(x) = 2x(-%dx + %dx^3) = 2x^2(%dx^2 - %d)\\). For \\(x \\ne 0\\), strict negativity \\(\\dot{V}(x) < 0\\) requires \\(%dx^2 - %d < 0\\), which simplifies to \\(x^2 < %d\\), or \\(|x| < %d\\). Thus the basin radius is \\(r = %d\\).",
				a, b, b, a, b, a, r*r, r, r,
			),
		}

	case 3:
		// Quadratic Lyapunov energy evaluation V(x0, y0) = a*x0^2 + b*y0^2
		a := smallInt(1, 5)
		b := smallInt(1, 5)
		x0 := smallInt(1, 4)
		y0 := smallInt(1, 4)
		ans := a*x0*x0 + b*y0*y0
		return generator.Problem{
			Question: fmt.Sprintf("Given the positive definite Lyapunov energy function \\(V(x, y) = %dx^2 + %dy^2\\), evaluate the energy level \\(V(%d, %d)\\).", a, b, x0, y0),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf(
				"Substituting the coordinates \\((x, y) = (%d, %d)\\) into the quadratic form gives \\(V(%d, %d) = %d(%d^2) + %d(%d^2) = %d(%d) + %d(%d) = %d + %d = %d\\).",
				x0, y0, x0, y0, a, x0, b, y0, a, x0*x0, b, y0*y0, a*x0*x0, b*y0*y0, ans,
			),
		}

	default:
		// Conservative Hamiltonian oscillator dot(E) = 0
		k := smallInt(2, 9)
		x0 := smallInt(1, 5)
		y0 := smallInt(1, 5)
		return generator.Problem{
			Question: fmt.Sprintf("For the undamped harmonic oscillator \\(x' = y\\) and \\(y' = -%dx\\) with total mechanical energy \\(E(x, y) = \\frac{%d}{2}x^2 + \\frac{1}{2}y^2\\), evaluate the orbital time derivative \\(\\dot{E}(%d, %d)\\).", k, k, x0, y0),
			Answer:   "0",
			Explanation: fmt.Sprintf(
				"Calculating the orbital derivative along trajectories: \\(\\dot{E} = \\frac{\\partial E}{\\partial x}x' + \\frac{\\partial E}{\\partial y}y' = (%dx)(y) + (y)(-%dx) = %dxy - %dxy = 0\\). Total energy is conserved along orbits, so \\(\\dot{E} = 0\\).",
				k, k, k, k,
			),
		}
	}
}

type laplaceEqGen struct{}

func (g *laplaceEqGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := rand.Intn(5)
	switch mode {
	case 0:
		// Mean Value Property at circle center
		R := smallInt(2, 6)
		x0 := smallInt(1, 5)
		y0 := smallInt(1, 5)
		v := smallInt(3, 25)
		return generator.Problem{
			Question: fmt.Sprintf("A function \\(u(x, y)\\) is harmonic on a disk of radius \\(R = %d\\) centered at \\((%d, %d)\\). If its average boundary value satisfies \\(\\frac{1}{2\\pi R}\\oint_{\\partial B} u\\,ds = %d\\), what is the center value \\(u(%d, %d)\\)?", R, x0, y0, v, x0, y0),
			Answer:   fmt.Sprintf("%d", v),
			Explanation: fmt.Sprintf(
				"By the Gauss Mean Value Property for harmonic functions, the value at the center of a disk equals the average over its boundary circle: \\(u(%d, %d) = \\frac{1}{2\\pi R}\\oint_{\\partial B} u\\,ds = %d\\).",
				x0, y0, v,
			),
		}

	case 1:
		// Harmonic cubic polynomial coefficient b for u = a*x^3 + b*x*y^2
		a := smallInt(1, 5)
		b := -3 * a
		return generator.Problem{
			Question: fmt.Sprintf("For what integer coefficient \\(b\\) is the polynomial \\(u(x, y) = %dx^3 + bx y^2\\) a harmonic function (\\(\\Delta u = 0\\))?", a),
			Answer:   fmt.Sprintf("%d", b),
			Explanation: fmt.Sprintf(
				"Calculating second partial derivatives: \\(u_{xx} = 6(%d)x = %dx\\) and \\(u_{yy} = 2bx\\). Requiring \\(\\Delta u = u_{xx} + u_{yy} = %dx + 2bx = 0\\) yields \\(2b = -%d\\), so \\(b = %d\\).",
				a, 6*a, 6*a, 6*a, b,
			),
		}

	case 2:
		// Harmonic polynomial evaluation u(x, y) = x^2 - y^2
		x0 := smallInt(3, 8)
		y0 := smallInt(1, x0-1)
		ans := x0*x0 - y0*y0
		return generator.Problem{
			Question: fmt.Sprintf("Evaluate the harmonic polynomial \\(u(x, y) = x^2 - y^2\\) at the point \\((%d, %d)\\).", x0, y0),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf(
				"The polynomial is harmonic because \\(\\Delta u = u_{xx} + u_{yy} = 2 - 2 = 0\\). Evaluating at \\((%d, %d)\\) gives \\(u(%d, %d) = %d^2 - %d^2 = %d - %d = %d\\).",
				x0, y0, x0, y0, x0, y0, x0*x0, y0*y0, ans,
			),
		}

	case 3:
		// Elliptic maximum principle interior bound
		M := smallInt(15, 60)
		return generator.Problem{
			Question: fmt.Sprintf("A harmonic function \\(u(x, y)\\) on the unit square satisfies boundary temperatures bounded by \\(u \\le %d\\). By the maximum principle, what is the maximum possible temperature at the interior center \\((0.5, 0.5)\\)?", M),
			Answer:   fmt.Sprintf("%d", M),
			Explanation: fmt.Sprintf(
				"By the strong maximum principle for elliptic equations, a non-constant harmonic function achieves its maximum strictly on the boundary. Therefore the interior temperature cannot exceed the boundary upper bound of %d.",
				M,
			),
		}

	default:
		// Total outward normal flux of harmonic function vanishes
		return generator.Problem{
			Question:    "For a smooth harmonic function \\(u(x, y)\\) on a bounded planar domain \\(\\Omega\\), evaluate the total outward boundary normal flux \\(\\oint_{\\partial \\Omega} \\frac{\\partial u}{\\partial n}\\,ds\\).",
			Answer:      "0",
			Explanation: "By Green's first identity (the divergence theorem), \\(\\oint_{\\partial \\Omega} \\frac{\\partial u}{\\partial n}\\,ds = \\iint_\\Omega \\Delta u\\,dA\\). Because \\(u\\) is harmonic (\\(\\Delta u = 0\\)), the integral vanishes identically, yielding 0.",
		}
	}
}

type volterraGen struct{}

func (g *volterraGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := rand.Intn(5)
	switch mode {
	case 0:
		// Initial value y(0) = y0
		y0 := smallInt(2, 20)
		c := smallInt(1, 9)
		return generator.Problem{
			Question: fmt.Sprintf("Consider the Volterra integral equation of the second kind \\(y(x) = %d + %dx + \\int_0^x (x - t)y(t)\\,dt\\). Evaluate the initial value \\(y(0)\\).", y0, c),
			Answer:   fmt.Sprintf("%d", y0),
			Explanation: fmt.Sprintf(
				"Evaluating at the lower integration limit \\(x = 0\\), the definite integral \\(\\int_0^0 (0 - t)y(t)\\,dt = 0\\). Therefore \\(y(0) = %d + %d(0) + 0 = %d\\).",
				y0, c, y0,
			),
		}

	case 1:
		// Initial derivative y'(0) = m + c * k
		m := smallInt(1, 8)
		c := smallInt(1, 5)
		k := smallInt(1, 6)
		ans := m + c*k
		return generator.Problem{
			Question: fmt.Sprintf("For the Volterra integral equation \\(y(x) = %d + %dx + \\int_0^x %dy(t)\\,dt\\), evaluate the derivative value \\(y'(0)\\).", k, m, c),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf(
				"Differentiating both sides with respect to \\(x\\) via the Leibniz integral rule gives \\(y'(x) = %d + %dy(x)\\). Evaluating at \\(x = 0\\) with \\(y(0) = %d\\) yields \\(y'(0) = %d + %d(%d) = %d + %d = %d\\).",
				m, c, k, m, c, k, m, c*k, ans,
			),
		}

	case 2:
		// Direct polynomial evaluation (a + c)*x0 + b
		a := smallInt(1, 5)
		b := smallInt(1, 9)
		c := smallInt(1, 5)
		x0 := smallInt(1, 6)
		ans := (a+c)*x0 + b
		return generator.Problem{
			Question: fmt.Sprintf("For the integral equation \\(y(x) = %dx + %d + \\int_0^x %d\\,dt\\), evaluate \\(y(%d)\\).", a, b, c, x0),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf(
				"Evaluating the definite integral gives \\(\\int_0^x %d\\,dt = %dx\\). Thus \\(y(x) = %dx + %d + %dx = %dx + %d\\). Substituting \\(x = %d\\) yields \\(y(%d) = %d(%d) + %d = %d + %d = %d\\).",
				c, c, a, b, c, a+c, b, x0, x0, a+c, x0, b, (a+c)*x0, b, ans,
			),
		}

	case 3:
		// Equivalent ODE growth rate constant a = k
		y0 := smallInt(1, 9)
		k := smallInt(2, 8)
		return generator.Problem{
			Question: fmt.Sprintf("Differentiating the Volterra integral equation \\(y(x) = %d + \\int_0^x %dy(t)\\,dt\\) converts it into the initial value problem \\(y' = ay\\) with \\(y(0) = %d\\). What is the value of coefficient \\(a\\)?", y0, k, y0),
			Answer:   fmt.Sprintf("%d", k),
			Explanation: fmt.Sprintf(
				"By the fundamental theorem of calculus, differentiating \\(y(x) = %d + \\int_0^x %dy(t)\\,dt\\) yields \\(y'(x) = %dy(x)\\). Therefore the growth coefficient is \\(a = %d\\).",
				y0, k, k, k,
			),
		}

	default:
		// Second derivative y''(0) = c * k
		k := smallInt(2, 7)
		c := smallInt(2, 6)
		ans := c * k
		return generator.Problem{
			Question: fmt.Sprintf("Consider the Volterra integral equation with displacement kernel \\(y(x) = %d + \\int_0^x %d(x - t)y(t)\\,dt\\). Evaluate the second derivative \\(y''(0)\\).", k, c),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf(
				"Differentiating twice with respect to \\(x\\) via the Leibniz rule: \\(y'(x) = %d\\int_0^x y(t)\\,dt\\) and \\(y''(x) = %dy(x)\\). At \\(x = 0\\), the integral equation gives \\(y(0) = %d\\). Therefore \\(y''(0) = %dy(0) = %d(%d) = %d\\).",
				c, c, k, c, c, k, ans,
			),
		}
	}
}
