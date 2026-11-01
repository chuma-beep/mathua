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
	type claim struct {
		a, b   int
		answer string
		reason string
	}
	claims := []claim{
		{a: -2, b: 2, answer: "2", reason: "The characteristic equation is m^2 - 3m + 2 = 0 with roots 1 and 2; the larger is 2."},
		{a: 1, b: -1, answer: "1", reason: "The characteristic equation is m^2 - 1 = 0 with roots -1 and 1; the larger is 1."},
		{a: -3, b: 4, answer: "2", reason: "The characteristic equation is m^2 - 4m + 4 = 0, a double root at 2."},
		{a: 2, b: -6, answer: "2", reason: "The characteristic equation is m^2 + m - 6 = (m+3)(m-2), so the larger root is 2."},
		{a: -5, b: 9, answer: "3", reason: "The characteristic equation is m^2 - 6m + 9 = 0, a double root at 3."},
		{a: 4, b: 2, answer: "-1", reason: "The characteristic equation is m^2 + 3m + 2 = (m+1)(m+2), so the larger root is -1."},
	}
	c := claims[rand.Intn(len(claims))]
	return generator.Problem{
		Question: "For the Cauchy-Euler equation x^2 y'' + (" + strconv.Itoa(c.a) +
			") x y' + (" + strconv.Itoa(c.b) +
			") y = 0, substitute y = x^m. What is the larger root of m(m-1) + (" +
			strconv.Itoa(c.a) + ") m + (" + strconv.Itoa(c.b) + ") = 0? (enter a number)",
		Answer:      c.answer,
		Explanation: c.reason,
		Facts:       map[string]string{"a": strconv.Itoa(c.a), "b": strconv.Itoa(c.b), "answer": c.answer},
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
	type claim struct {
		question string
		answer   string
		kind     string
		k        int
		y0       int
		reason   string
	}
	claims := []claim{
		{
			question: "For y' = 2y with y(0) = 1, write the solution as a power series sum a_n x^n. What is a_2? (enter a number)",
			answer:   "2", kind: "first", k: 2, y0: 1,
			reason: "The solution is e^(2x) = sum (2x)^n/n!, so a_2 = 2^2/2 = 2.",
		},
		{
			question: "For y' = 3y with y(0) = 1, write the solution as a power series sum a_n x^n. What is a_2? (enter a decimal)",
			answer:   "4.5", kind: "first", k: 3, y0: 1,
			reason: "The solution is e^(3x), so a_2 = 3^2/2 = 4.5.",
		},
		{
			question: "For y' = y with y(0) = 2, write the solution as a power series sum a_n x^n. What is a_2? (enter a number)",
			answer:   "1", kind: "first", k: 1, y0: 2,
			reason: "The solution is 2e^x = 2 sum x^n/n!, so a_2 = 2/2 = 1.",
		},
		{
			question: "For y'' + y = 0 with y(0) = 1 and y'(0) = 0, the solution is cos x = sum a_n x^n. What is a_2? (enter a decimal)",
			answer:   "-0.5", kind: "cos", k: 1, y0: 1,
			reason: "The cosine series is 1 - x^2/2 + ..., so a_2 = -1/2 = -0.5.",
		},
		{
			question: "For y'' - y = 0 with y(0) = 1 and y'(0) = 0, the solution is cosh x = sum a_n x^n. What is a_2? (enter a decimal)",
			answer:   "0.5", kind: "cosh", k: 1, y0: 1,
			reason: "The cosh series is 1 + x^2/2 + ..., so a_2 = 1/2 = 0.5.",
		},
		{
			question: "For y' = 4y with y(0) = 1, write the solution as a power series sum a_n x^n. What is a_2? (enter a number)",
			answer:   "8", kind: "first", k: 4, y0: 1,
			reason: "The solution is e^(4x), so a_2 = 4^2/2 = 8.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	return generator.Problem{
		Question:    c.question,
		Answer:      c.answer,
		Explanation: c.reason,
		Facts: map[string]string{
			"kind": c.kind, "k": strconv.Itoa(c.k), "y0": strconv.Itoa(c.y0), "answer": c.answer,
		},
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

func (g *springGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		omega2   int
		y0, v0   int
		tLabel   string
		reason   string
	}
	claims := []claim{
		{
			question: "Solve y'' + 4y = 0 with y(0) = 0 and y'(0) = 2. What is y(pi/4)? (enter a number)",
			answer:   "1", omega2: 4, y0: 0, v0: 2, tLabel: "pi/4",
			reason: "The solution is sin(2t), and sin(pi/2) = 1.",
		},
		{
			question: "Solve y'' + y = 0 with y(0) = 1 and y'(0) = 0. What is y(pi)? (enter a number)",
			answer:   "-1", omega2: 1, y0: 1, v0: 0, tLabel: "pi",
			reason: "The solution is cos(t), and cos(pi) = -1.",
		},
		{
			question: "Solve y'' + 4y = 0 with y(0) = 2 and y'(0) = 0. What is y(pi/4)? (enter a number)",
			answer:   "0", omega2: 4, y0: 2, v0: 0, tLabel: "pi/4",
			reason: "The solution is 2cos(2t), and cos(pi/2) = 0.",
		},
		{
			question: "Solve y'' + 9y = 0 with y(0) = 0 and y'(0) = 3. What is y(pi/6)? (enter a number)",
			answer:   "1", omega2: 9, y0: 0, v0: 3, tLabel: "pi/6",
			reason: "The solution is sin(3t), and sin(pi/2) = 1.",
		},
		{
			question: "Solve y'' + 9y = 0 with y(0) = 1 and y'(0) = 0. What is y(pi/6)? (enter a number)",
			answer:   "0", omega2: 9, y0: 1, v0: 0, tLabel: "pi/6",
			reason: "The solution is cos(3t), and cos(pi/2) = 0.",
		},
		{
			question: "Solve y'' + 4y = 0 with y(0) = 0 and y'(0) = 4. What is y(pi/4)? (enter a number)",
			answer:   "2", omega2: 4, y0: 0, v0: 4, tLabel: "pi/4",
			reason: "The solution is 2sin(2t), and sin(pi/2) = 1, so the value is 2.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	return generator.Problem{
		Question:    c.question,
		Answer:      c.answer,
		Explanation: c.reason,
		Facts: map[string]string{
			"omega2": strconv.Itoa(c.omega2), "y0": strconv.Itoa(c.y0),
			"v0": strconv.Itoa(c.v0), "t": c.tLabel, "answer": c.answer,
		},
	}
}

// ---------------------------------------------------------------------------
// 8. Mixing problems (numeric)
// ---------------------------------------------------------------------------

type mixingGen struct{}

func (g *mixingGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		v, rate  int
		c        float64
		reason   string
	}
	claims := []claim{
		{
			question: "A tank holds 100 litres of water. Brine with 0.2 kg/L enters at 5 L/min and the well-mixed solution leaves at 5 L/min. What is the limiting amount of salt in the tank? (enter a number)",
			answer:   "20", v: 100, rate: 5, c: 0.2,
			reason: "At equilibrium inflow equals outflow: 5(0.2) = 5Q/100, so Q = 100(0.2) = 20.",
		},
		{
			question: "A tank holds 100 litres of water. Brine with 0.1 kg/L enters at 10 L/min and the well-mixed solution leaves at 10 L/min. What is the limiting amount of salt in the tank? (enter a number)",
			answer:   "10", v: 100, rate: 10, c: 0.1,
			reason: "The equilibrium is Q = V c = 100(0.1) = 10.",
		},
		{
			question: "A tank holds 200 litres of water. Brine with 0.5 kg/L enters at 10 L/min and the well-mixed solution leaves at 10 L/min. What is the limiting amount of salt in the tank? (enter a number)",
			answer:   "100", v: 200, rate: 10, c: 0.5,
			reason: "The equilibrium is Q = V c = 200(0.5) = 100.",
		},
		{
			question: "A tank holds 50 litres of water. Brine with 0.4 kg/L enters at 5 L/min and the well-mixed solution leaves at 5 L/min. What is the limiting amount of salt in the tank? (enter a number)",
			answer:   "20", v: 50, rate: 5, c: 0.4,
			reason: "The equilibrium is Q = V c = 50(0.4) = 20.",
		},
		{
			question: "A tank holds 100 litres of water. Brine with 0.3 kg/L enters at 2 L/min and the well-mixed solution leaves at 2 L/min. What is the limiting amount of salt in the tank? (enter a number)",
			answer:   "30", v: 100, rate: 2, c: 0.3,
			reason: "The equilibrium is Q = V c = 100(0.3) = 30; the flow rate cancels.",
		},
		{
			question: "A tank holds 400 litres of water. Brine with 0.1 kg/L enters at 20 L/min and the well-mixed solution leaves at 20 L/min. What is the limiting amount of salt in the tank? (enter a number)",
			answer:   "40", v: 400, rate: 20, c: 0.1,
			reason: "The equilibrium is Q = V c = 400(0.1) = 40.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	return generator.Problem{
		Question:    c.question,
		Answer:      c.answer,
		Explanation: c.reason,
		Facts: map[string]string{
			"v": strconv.Itoa(c.v), "rate": strconv.Itoa(c.rate),
			"c": strconv.FormatFloat(c.c, 'g', -1, 64), "answer": c.answer,
		},
	}
}

// ---------------------------------------------------------------------------
// 9. Predator-prey systems (numeric)
// ---------------------------------------------------------------------------

type predatorPreyGen struct{}

func (g *predatorPreyGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		a, b, c  int
		d        int
		ask      string
		reason   string
	}
	claims := []claim{
		{
			question: "For the Lotka-Volterra system x' = ax - bxy, y' = -cy + dxy with a = 2, b = 1, c = 3 and d = 1, what is the x-coordinate c/d of the nonzero equilibrium? (enter a number)",
			answer:   "3", a: 2, b: 1, c: 3, d: 1, ask: "x",
			reason: "Setting y' = 0 gives x = c/d = 3/1 = 3.",
		},
		{
			question: "For the Lotka-Volterra system x' = ax - bxy, y' = -cy + dxy with a = 4, b = 2, c = 2 and d = 1, what is the y-coordinate a/b of the nonzero equilibrium? (enter a number)",
			answer:   "2", a: 4, b: 2, c: 2, d: 1, ask: "y",
			reason: "Setting x' = 0 gives y = a/b = 4/2 = 2.",
		},
		{
			question: "For the Lotka-Volterra system x' = ax - bxy, y' = -cy + dxy with a = 1, b = 1, c = 5 and d = 1, what is the x-coordinate c/d of the nonzero equilibrium? (enter a number)",
			answer:   "5", a: 1, b: 1, c: 5, d: 1, ask: "x",
			reason: "Setting y' = 0 gives x = c/d = 5.",
		},
		{
			question: "For the Lotka-Volterra system x' = ax - bxy, y' = -cy + dxy with a = 6, b = 3, c = 4 and d = 2, what is the x-coordinate c/d of the nonzero equilibrium? (enter a number)",
			answer:   "2", a: 6, b: 3, c: 4, d: 2, ask: "x",
			reason: "Setting y' = 0 gives x = c/d = 4/2 = 2.",
		},
		{
			question: "For the Lotka-Volterra system x' = ax - bxy, y' = -cy + dxy with a = 3, b = 1, c = 1 and d = 1, what is the y-coordinate a/b of the nonzero equilibrium? (enter a number)",
			answer:   "3", a: 3, b: 1, c: 1, d: 1, ask: "y",
			reason: "Setting x' = 0 gives y = a/b = 3.",
		},
		{
			question: "For the Lotka-Volterra system x' = ax - bxy, y' = -cy + dxy with a = 5, b = 1, c = 10 and d = 2, what is the x-coordinate c/d of the nonzero equilibrium? (enter a number)",
			answer:   "5", a: 5, b: 1, c: 10, d: 2, ask: "x",
			reason: "Setting y' = 0 gives x = c/d = 10/2 = 5.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	return generator.Problem{
		Question:    c.question,
		Answer:      c.answer,
		Explanation: c.reason,
		Facts: map[string]string{
			"a": strconv.Itoa(c.a), "b": strconv.Itoa(c.b),
			"c": strconv.Itoa(c.c), "d": strconv.Itoa(c.d), "ask": c.ask, "answer": c.answer,
		},
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
