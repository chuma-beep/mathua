package odes

// The 2026-10-18 differential-equations batch: ten concepts that take the
// domain from 32 to the target of 42.
//
// The survey found the domain already covering the standard first-order and
// second-order methods, Laplace transforms, systems, numerical methods, PDEs,
// Sturm-Liouville theory and integral equations. Several obvious additions were
// already taught and were not repeated: integrating factors in the linear and
// exact lessons, the y = vx substitution in the homogeneous lesson,
// undetermined coefficients and the modification rule in the nonhomogeneous
// lesson, separation of variables and d'Alembert in the heat and wave lessons,
// linearization in the stability lesson, and Picard and Fredholm in the
// boundary and Volterra lessons.
//
// The batch fills real gaps: Bernoulli equations, Cauchy-Euler equations,
// reduction of order, power-series solutions, the matrix exponential, Heaviside
// steps and impulses, the spring-mass system, mixing problems, predator-prey
// systems, and Heun's method.
//
// Five concepts are numeric and five are multiple choice. Every numeric
// generator publishes the facts its answer follows from, so the tests recompute
// the answer from those facts.

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerBatch(reg *generator.Registry) {
	reg.Register("ode.first_order.bernoulli", &bernoulliGen{})
	reg.Register("ode.second_order.cauchy_euler", &cauchyEulerGen{})
	reg.Register("ode.second_order.reduction", &reductionGen{})
	reg.Register("ode.second_order.power_series", &powerSeriesGen{})
	reg.Register("ode.systems.matrix_exp", &matrixExpGen{})
	reg.Register("ode.transforms.step", &stepGen{})
	reg.Register("ode.second_order.spring", &springGen{})
	reg.Register("ode.first_order.mixing", &mixingGen{})
	reg.Register("ode.systems.predator_prey", &predatorPreyGen{})
	reg.Register("ode.numerical.heun", &heunGen{})
}

// ---------------------------------------------------------------------------
// 1. Bernoulli equations
// ---------------------------------------------------------------------------

type bernoulliGen struct{}

func (g *bernoulliGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return bernoulliProblem()
}

// ---------------------------------------------------------------------------
// 2. Cauchy-Euler equations (numeric)
// ---------------------------------------------------------------------------

type cauchyEulerGen struct{}

func (g *cauchyEulerGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := odeBridgeOrdinal(ctx)
	// Draw the two integer characteristic roots directly. With roots p and q
	// the characteristic equation m^2 + (a-1)m + b = 0 factors as (m-p)(m-q),
	// so a = 1-p-q and b = p*q, and the discriminant (p-q)^2 is never negative.
	p := w4OR(-6, 2)
	q := p + 1 + w4OR(1, 8)
	a, b := 1-p-q, p*q
	facts := map[string]string{
		"a": strconv.Itoa(a),
		"b": strconv.Itoa(b),
		"p": strconv.Itoa(p),
		"q": strconv.Itoa(q),
	}
	stem := fmt.Sprintf("For the Cauchy-Euler equation x^2 y'' + (%d) x y' + (%d) y = 0, substitute y = x^m.", a, b)
	characteristic := fmt.Sprintf("m(m-1) + (%d) m + (%d) = 0", a, b)
	factorised := fmt.Sprintf("(m - (%d))(m - (%d))", p, q)
	switch o {
	case 1:
		facts["ask"] = "larger"
		return odeBridgeProblem(o,
			stem+" What is the larger root of "+characteristic+"? (enter a number)",
			strconv.Itoa(q),
			fmt.Sprintf("The characteristic equation is m^2 + (%d) m + (%d) = %s, so the larger root is %d.", a-1, b, factorised, q),
			facts)
	case 2:
		facts["ask"] = "smaller"
		return odeBridgeProblem(o,
			stem+" What is the smaller root of "+characteristic+"? (enter a number)",
			strconv.Itoa(p),
			fmt.Sprintf("The characteristic equation is m^2 + (%d) m + (%d) = %s, so the smaller root is %d.", a-1, b, factorised, p),
			facts)
	default:
		facts["ask"] = "sum"
		return odeBridgeProblem(o,
			stem+" The roots are the exponents in y = x^m. What is the sum of the two roots of "+characteristic+"? (enter a number)",
			strconv.Itoa(p+q),
			fmt.Sprintf("The characteristic equation is m^2 + (%d) m + (%d) = %s, whose roots %d and %d sum to %d.", a-1, b, factorised, p, q, p+q),
			facts)
	}
}

// ---------------------------------------------------------------------------
// 3. Reduction of order
// ---------------------------------------------------------------------------

type reductionGen struct{}

func (g *reductionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return reductionProblem()
}

// ---------------------------------------------------------------------------
// 4. Power-series solutions (numeric)
// ---------------------------------------------------------------------------

type powerSeriesGen struct{}

func (g *powerSeriesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := odeBridgeOrdinal(ctx)
	facts := map[string]string{}
	switch o {
	case 1:
		k, y0 := w4OR(2, 7), w4OR(1, 3)
		answer := float64(y0*k*k) / 2
		facts["kind"] = "first"
		facts["k"] = strconv.Itoa(k)
		facts["y0"] = strconv.Itoa(y0)
		return odeBridgeProblem(o,
			fmt.Sprintf("For y' = %d y with y(0) = %d, write the solution as a power series sum a_n x^n. What is a_2? (enter a decimal)", k, y0),
			strconv.FormatFloat(answer, 'f', -1, 64),
			fmt.Sprintf("The solution is %d e^(%d x) = %d sum (%d x)^n/n!, so a_2 = %d times %d squared over 2 = %s.",
				y0, k, y0, k, y0, k, strconv.FormatFloat(answer, 'f', -1, 64)),
			facts)
	case 2:
		k := w4OR(1, 12)
		answer := -float64(k*k) / 2
		facts["kind"] = "cos"
		facts["k"] = strconv.Itoa(k)
		facts["y0"] = "1"
		return odeBridgeProblem(o,
			fmt.Sprintf("For y'' + %d y = 0 with y(0) = 1 and y'(0) = 0, the solution is cos(%d x) = sum a_n x^n. What is a_2? (enter a decimal)", k*k, k),
			strconv.FormatFloat(answer, 'f', -1, 64),
			fmt.Sprintf("The cosine series is 1 - (%d x)^2/2 + ..., so a_2 = -%d/2 = %s.", k, k*k, strconv.FormatFloat(answer, 'f', -1, 64)),
			facts)
	default:
		k := w4OR(1, 12)
		answer := float64(k*k) / 2
		facts["kind"] = "cosh"
		facts["k"] = strconv.Itoa(k)
		facts["y0"] = "1"
		return odeBridgeProblem(o,
			fmt.Sprintf("For y'' - %d y = 0 with y(0) = 1 and y'(0) = 0, the solution is cosh(%d x) = sum a_n x^n. What is a_2? (enter a decimal)", k*k, k),
			strconv.FormatFloat(answer, 'f', -1, 64),
			fmt.Sprintf("The cosh series is 1 + (%d x)^2/2 + ..., so a_2 = %d/2 = %s.", k, k*k, strconv.FormatFloat(answer, 'f', -1, 64)),
			facts)
	}
}

// ---------------------------------------------------------------------------
// 5. The matrix exponential
// ---------------------------------------------------------------------------

type matrixExpGen struct{}

func (g *matrixExpGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return matrixExpProblem()
}

// ---------------------------------------------------------------------------
// 6. Heaviside steps and impulses
// ---------------------------------------------------------------------------

type stepGen struct{}

func (g *stepGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return stepProblem()
}

// ---------------------------------------------------------------------------
// 7. The spring-mass system (numeric)
// ---------------------------------------------------------------------------

type springGen struct{}

// odeQuarterTurn returns cos and sin of q times pi over 2 for an integer q.
// Every evaluation time in this generator is a quarter turn of the oscillation,
// so the solution is an integer without any trigonometry in floating point.
func odeQuarterTurn(q int) (int, int) {
	switch ((q % 4) + 4) % 4 {
	case 0:
		return 1, 0
	case 1:
		return 0, 1
	case 2:
		return -1, 0
	default:
		return 0, -1
	}
}

func (g *springGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := odeBridgeOrdinal(ctx)
	var omega, y0, v0, turns int
	var label string
	switch o {
	case 1:
		s := w4OR(1, 6)
		omega = 2 * s
		y0 = w4OR(0, 5)
		v0 = omega * w4OR(-3, 3)
		turns = s
		label = "pi/4"
	case 2:
		s := w4OR(1, 4)
		omega = 3 * s
		y0 = w4OR(0, 5)
		v0 = omega * w4OR(-3, 3)
		turns = s
		label = "pi/6"
	default:
		s := w4OR(1, 6)
		omega = s
		y0 = w4OR(0, 6)
		v0 = omega * w4OR(-3, 3)
		turns = 2 * s
		label = "pi"
	}
	cosine, sine := odeQuarterTurn(turns)
	amplitude := v0 / omega
	answer := y0*cosine + amplitude*sine
	return odeBridgeProblem(o,
		fmt.Sprintf("Solve y'' + %d y = 0 with y(0) = %d and y'(0) = %d. What is y(%s)? (enter a number)",
			omega*omega, y0, v0, label),
		strconv.Itoa(answer),
		fmt.Sprintf("The solution is %d cos(%d t) + (%d) sin(%d t). At t = %s the cosine equals %d and the sine equals %d, so y = %d.",
			y0, omega, amplitude, omega, label, cosine, sine, answer),
		map[string]string{
			"omega2": strconv.Itoa(omega * omega),
			"y0":     strconv.Itoa(y0),
			"v0":     strconv.Itoa(v0),
			"t":      label,
		})
}

// ---------------------------------------------------------------------------
// 8. Mixing problems (numeric)
// ---------------------------------------------------------------------------

type mixingGen struct{}

func (g *mixingGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := odeBridgeOrdinal(ctx)
	// The volume is a multiple of 20 and the concentration a multiple of 5
	// hundredths, so V*c is an integer and the grader never sees float noise.
	volume := 20 * w4OR(1, 20)
	rate := w4OR(2, 20)
	hundredths := []int{10, 20, 25, 40, 50}[rand.Intn(5)]
	answer := volume * hundredths / 100
	grams := hundredths * 10
	facts := map[string]string{
		"v":      strconv.Itoa(volume),
		"rate":   strconv.Itoa(rate),
		"c":      odeDecimal(float64(hundredths) / 100),
		"answer": strconv.Itoa(answer),
	}
	switch o {
	case 1:
		return odeBridgeProblem(o,
			fmt.Sprintf("A tank holds %d litres of water. Brine with %s kg/L enters at %d L/min and the well-mixed solution leaves at the same rate. What is the limiting amount of salt in the tank? (enter a number)",
				volume, odeDecimal(float64(hundredths)/100), rate),
			strconv.Itoa(answer),
			fmt.Sprintf("At equilibrium inflow equals outflow: %d times %s equals %d times Q over %d, so Q = %d times %s = %d.",
				rate, odeDecimal(float64(hundredths)/100), rate, volume, volume, odeDecimal(float64(hundredths)/100), answer),
			facts)
	case 2:
		initial := w4OR(1, 30)
		return odeBridgeProblem(o,
			fmt.Sprintf("A tank holds %d litres and starts with %d kg of salt dissolved in it. Brine with %s kg/L enters at %d L/min and the well-mixed solution leaves at the same rate. What is the limiting amount of salt in the tank? (enter a number)",
				volume, initial, odeDecimal(float64(hundredths)/100), rate),
			strconv.Itoa(answer),
			fmt.Sprintf("The initial %d kg washes out or grows towards the equilibrium, which depends only on the inflow: Q = %d times %s = %d.",
				initial, volume, odeDecimal(float64(hundredths)/100), answer),
			facts)
	default:
		return odeBridgeProblem(o,
			fmt.Sprintf("Brine with %d grams per litre enters a %d-litre tank at %d L/min and the well-mixed solution leaves at the same rate. What is the limiting amount of salt in the tank, in kilograms? (enter a number)",
				grams, volume, rate),
			strconv.Itoa(answer),
			fmt.Sprintf("%d grams per litre is %s kg/L, so the equilibrium is Q = %d times %s = %d kg.",
				grams, odeDecimal(float64(hundredths)/100), volume, odeDecimal(float64(hundredths)/100), answer),
			facts)
	}
}

// ---------------------------------------------------------------------------
// 9. Predator-prey systems (numeric)
// ---------------------------------------------------------------------------

type predatorPreyGen struct{}

func (g *predatorPreyGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := odeBridgeOrdinal(ctx)
	// Draw the two coordinates of the coexistence equilibrium, then choose the
	// coefficients that produce them: x* = c/d and y* = a/b exactly.
	d := w4OR(2, 5)
	xStar := w4OR(1, 8)
	c := d * xStar
	b := w4OR(1, 5)
	yStar := w4OR(1, 8)
	a := b * yStar
	facts := map[string]string{
		"a": strconv.Itoa(a), "b": strconv.Itoa(b),
		"c": strconv.Itoa(c), "d": strconv.Itoa(d),
	}
	system := fmt.Sprintf("For the Lotka-Volterra system x' = ax - bxy, y' = -cy + dxy with a = %d, b = %d, c = %d and d = %d", a, b, c, d)
	switch o {
	case 1:
		facts["ask"] = "x"
		return odeBridgeProblem(o,
			system+", what is the x-coordinate c/d of the nonzero equilibrium? (enter a number)",
			strconv.Itoa(xStar),
			fmt.Sprintf("Setting y' = 0 gives x = c/d = %d/%d = %d.", c, d, xStar),
			facts)
	case 2:
		facts["ask"] = "y"
		return odeBridgeProblem(o,
			system+", what is the y-coordinate a/b of the nonzero equilibrium? (enter a number)",
			strconv.Itoa(yStar),
			fmt.Sprintf("Setting x' = 0 gives y = a/b = %d/%d = %d.", a, b, yStar),
			facts)
	default:
		facts["ask"] = "x"
		return odeBridgeProblem(o,
			system+", the nonzero equilibrium is the coexistence state. What is the prey population there? (enter a number)",
			strconv.Itoa(xStar),
			fmt.Sprintf("At coexistence y' = 0, so the prey population is x = c/d = %d/%d = %d.", c, d, xStar),
			facts)
	}
}

// ---------------------------------------------------------------------------
// 10. Heun's method (numeric)
// ---------------------------------------------------------------------------

type heunGen struct{}

func (g *heunGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	functions := []string{"y", "double", "square", "tplusy"}
	name := functions[rand.Intn(len(functions))]
	y0 := float64(1 + rand.Intn(3))
	t0 := float64(rand.Intn(3))
	steps := []float64{0.1, 0.05, 0.2}
	h := steps[rand.Intn(len(steps))]

	slope := func(tv, y float64) float64 {
		switch name {
		case "y":
			return y
		case "double":
			return 2 * y
		case "square":
			return y * y
		default:
			return tv + y
		}
	}
	predictor := y0 + h*slope(t0, y0)
	value := y0 + h/2*(slope(t0, y0)+slope(t0+h, predictor))
	// Round before printing: a binary float renders 3.32 as 3.3200000000000003.
	answer := strconv.FormatFloat(math.Round(value*1e6)/1e6, 'f', -1, 64)

	symbol := map[string]string{
		"y": "y", "double": "2y", "square": "y^2", "tplusy": "t + y",
	}[name]
	return generator.Problem{
		Question: fmt.Sprintf("Use one step of Heun's method with h = %s to estimate y(%s) for y' = %s with y(%s) = %s. (enter a decimal)",
			trimFloat(h), trimFloat(t0+h), symbol, trimFloat(t0), trimFloat(y0)),
		Answer: answer,
		Explanation: fmt.Sprintf("The predictor is %s. The slopes at the start and at the predictor are %s and %s, so the corrector is %s.",
			trimFloat(predictor), trimFloat(slope(t0, y0)), trimFloat(slope(t0+h, predictor)), answer),
		Facts: map[string]string{
			"function": name,
			"y0":       trimFloat(y0),
			"t0":       trimFloat(t0),
			"h":        trimFloat(h),
			"answer":   answer,
		},
	}
}

// trimFloat renders a float without a trailing .0 so the stem reads 0 and 0.1
// rather than 0.000000 and 0.100000.
func trimFloat(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}
