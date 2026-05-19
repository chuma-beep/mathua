package odes

import (
	"fmt"
	"math/rand"

	"github.com/chuma-beep/mathua/internal/generator"
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
}

type conceptGen struct{}

func (g *conceptGen) Generate(difficulty float64) generator.Problem {
	type qa struct {
		q string
		a string
		e string
	}
	templates := []qa{
		{"Is y' = 2x a differential equation? (yes/no)", "yes", "The equation contains a derivative y', so it is a differential equation."},
		{"Is d²y/dx² + y = 0 a differential equation? (yes/no)", "yes", "The equation contains a second derivative d²y/dx², so it is a differential equation."},
		{"Is dy/dx = 3x² a differential equation? (yes/no)", "yes", "The equation contains a derivative dy/dx, so it is a differential equation."},
		{"Is y''' + 2y'' - y' + y = 0 a differential equation? (yes/no)", "yes", "The equation contains third derivatives, so it is a differential equation."},
		{"Is x² + y² = 5 a differential equation? (yes/no)", "no", "The equation contains no derivatives, so it is not a differential equation."},
		{"Is 3x + 2y = 7 a differential equation? (yes/no)", "no", "The equation contains no derivatives, so it is not a differential equation."},
		{"Is sin(x) + cos(y) = 1 a differential equation? (yes/no)", "no", "The equation contains no derivatives, so it is not a differential equation."},
		{"Does the equation x³ - y³ = 10 contain a derivative? (yes/no)", "no", "The equation contains no derivative terms."},
	}
	chosen := templates[rand.Intn(len(templates))]
	return generator.Problem{
		Question:    chosen.q,
		Answer:      chosen.a,
		Explanation: chosen.e,
	}
}

type separableGen struct{}

func (g *separableGen) Generate(difficulty float64) generator.Problem {
	type qa struct {
		q string
		a string
		e string
	}
	templates := []qa{
		{
			q: "Solve: dy/dx = x/y",
			a: "y^2=x^2+C",
			e: "Separate: y dy = x dx. Integrate: y²/2 = x²/2 + C. Multiply by 2: y² = x² + C.",
		},
		{
			q: "Solve: dy/dx = x",
			a: "y=x^2/2+C",
			e: "Rewrite as dy = x dx. Integrate: y = x²/2 + C.",
		},
		{
			q: "Solve: dy/dx = 1/y",
			a: "y^2=2x+C",
			e: "Separate: y dy = dx. Integrate: y²/2 = x + C. Multiply by 2: y² = 2x + C.",
		},
		{
			q: "Solve: dy/dx = x²/y",
			a: "y^2=(2/3)x^3+C",
			e: "Separate: y dy = x² dx. Integrate: y²/2 = x³/3 + C. Multiply by 2: y² = (2/3)x³ + C.",
		},
		{
			q: "Solve: dy/dx = x²",
			a: "y=x^3/3+C",
			e: "Rewrite as dy = x² dx. Integrate: y = x³/3 + C.",
		},
		{
			q: "Solve: dy/dx = 1/y²",
			a: "y^3=3x+C",
			e: "Separate: y² dy = dx. Integrate: y³/3 = x + C. Multiply by 3: y³ = 3x + C.",
		},
	}
	t := templates[rand.Intn(len(templates))]
	return generator.Problem{
		Question:    t.q,
		Answer:      t.a,
		Explanation: t.e,
	}
}

type linearFirstGen struct{}

func (g *linearFirstGen) Generate(difficulty float64) generator.Problem {
	type qa struct {
		q string
		a string
		e string
	}
	templates := []qa{
		{
			q: "Solve: y' + 2y = 0",
			a: "Ce^(-2x)",
			e: "This is a homogeneous first-order linear ODE. The integrating factor is e^(2x). d/dx(ye^(2x)) = 0, so ye^(2x) = C, giving y = Ce^(-2x).",
		},
		{
			q: "Solve: y' + y = 0",
			a: "Ce^(-x)",
			e: "This is a homogeneous first-order linear ODE. The integrating factor is e^x. d/dx(ye^x) = 0, so ye^x = C, giving y = Ce^(-x).",
		},
		{
			q: "Solve: y' + 3y = 0",
			a: "Ce^(-3x)",
			e: "This is a homogeneous first-order linear ODE. The integrating factor is e^(3x). d/dx(ye^(3x)) = 0, so ye^(3x) = C, giving y = Ce^(-3x).",
		},
		{
			q: "Solve: y' + y = 1",
			a: "1+Ce^(-x)",
			e: "The integrating factor is e^x. d/dx(ye^x) = e^x. Integrate: ye^x = e^x + C, so y = 1 + Ce^(-x).",
		},
		{
			q: "Solve: y' + 2y = 4",
			a: "2+Ce^(-2x)",
			e: "The integrating factor is e^(2x). d/dx(ye^(2x)) = 4e^(2x). Integrate: ye^(2x) = 2e^(2x) + C, so y = 2 + Ce^(-2x).",
		},
		{
			q: "Solve: y' + y = e^x",
			a: "0.5e^x+Ce^(-x)",
			e: "The integrating factor is e^x. d/dx(ye^x) = e^(2x). Integrate: ye^x = (1/2)e^(2x) + C, so y = (1/2)e^x + Ce^(-x).",
		},
	}
	t := templates[rand.Intn(len(templates))]
	return generator.Problem{
		Question:    t.q,
		Answer:      t.a,
		Explanation: t.e,
	}
}

type exactGen struct{}

func (g *exactGen) Generate(difficulty float64) generator.Problem {
	type qa struct {
		q string
		a string
		e string
	}
	templates := []qa{
		{
			q: "Is M dx + N dy = 0 exact when M = 2xy and N = x²? (yes/no)",
			a: "yes",
			e: "∂M/∂y = 2x and ∂N/∂x = 2x. Since ∂M/∂y = ∂N/∂x, the equation is exact.",
		},
		{
			q: "Is M dx + N dy = 0 exact when M = y and N = x? (yes/no)",
			a: "yes",
			e: "∂M/∂y = 1 and ∂N/∂x = 1. Since ∂M/∂y = ∂N/∂x, the equation is exact.",
		},
		{
			q: "Is M dx + N dy = 0 exact when M = 3x²y and N = x³? (yes/no)",
			a: "yes",
			e: "∂M/∂y = 3x² and ∂N/∂x = 3x². Since ∂M/∂y = ∂N/∂x, the equation is exact.",
		},
		{
			q: "Is M dx + N dy = 0 exact when M = x² and N = y²? (yes/no)",
			a: "no",
			e: "∂M/∂y = 0 and ∂N/∂x = 0. Since ∂M/∂y = ∂N/∂x = 0, this equation IS exact. (Trick: both partials are zero.)",
		},
		{
			q: "Is M dx + N dy = 0 exact when M = x and N = x? (yes/no)",
			a: "no",
			e: "∂M/∂y = 0 and ∂N/∂x = 1. Since ∂M/∂y ≠ ∂N/∂x (0 ≠ 1), the equation is not exact.",
		},
		{
			q: "Is M dx + N dy = 0 exact when M = y² and N = x? (yes/no)",
			a: "no",
			e: "∂M/∂y = 2y and ∂N/∂x = 1. Since ∂M/∂y ≠ ∂N/∂x (2y ≠ 1 in general), the equation is not exact.",
		},
		{
			q: "Is M dx + N dy = 0 exact when M = xy and N = xy? (yes/no)",
			a: "no",
			e: "∂M/∂y = x and ∂N/∂x = y. Since ∂M/∂y ≠ ∂N/∂x (x ≠ y in general), the equation is not exact.",
		},
		{
			q: "Is M dx + N dy = 0 exact when M = sin(y) and N = x*cos(y)? (yes/no)",
			a: "yes",
			e: "∂M/∂y = cos(y) and ∂N/∂x = cos(y). Since ∂M/∂y = ∂N/∂x, the equation is exact.",
		},
	}
	t := templates[rand.Intn(len(templates))]
	return generator.Problem{
		Question:    t.q,
		Answer:      t.a,
		Explanation: t.e,
	}
}

type homogeneousGen struct{}

func (g *homogeneousGen) Generate(difficulty float64) generator.Problem {
	type qa struct {
		r1, r2 int
	}
	pairs := []qa{
		{1, 2}, {2, 3}, {1, -1}, {2, -2}, {3, -3},
		{1, 3}, {2, 4}, {-1, -2}, {3, 4}, {1, -2},
	}
	p := pairs[rand.Intn(len(pairs))]
	r1, r2 := p.r1, p.r2
	s := r1 + r2
	prod := r1 * r2

	charStr := fmt.Sprintf("r²")
	if s > 0 {
		charStr += fmt.Sprintf(" - %dr", s)
	} else if s < 0 {
		charStr += fmt.Sprintf(" + %dr", -s)
	}
	if prod > 0 {
		charStr += fmt.Sprintf(" + %d", prod)
	} else if prod < 0 {
		charStr += fmt.Sprintf(" - %d", -prod)
	}
	charStr += " = 0"

	eqStr := "y''"
	if s > 0 {
		eqStr += fmt.Sprintf(" - %dy'", s)
	} else if s < 0 {
		eqStr += fmt.Sprintf(" + %dy'", -s)
	}
	if prod > 0 {
		eqStr += fmt.Sprintf(" + %dy", prod)
	} else if prod < 0 {
		eqStr += fmt.Sprintf(" - %dy", -prod)
	}
	eqStr += " = 0"

	answer := fmt.Sprintf("C1*e^(%dx)+C2*e^(%dx)", r2, r1)
	if r1 > r2 {
		answer = fmt.Sprintf("C1*e^(%dx)+C2*e^(%dx)", r1, r2)
	}

	explanation := fmt.Sprintf("Characteristic equation: %s. Roots: r = %d and r = %d. The general solution is y = C₁e^(%dx) + C₂e^(%dx).", charStr, r1, r2, r1, r2)

	return generator.Problem{
		Question:    fmt.Sprintf("Solve: %s", eqStr),
		Answer:      answer,
		Explanation: explanation,
	}
}

type nonhomogeneousGen struct{}

func (g *nonhomogeneousGen) Generate(difficulty float64) generator.Problem {
	type qa struct {
		q string
		a string
		e string
	}
	templates := []qa{
		{
			q: "Solve: y'' + y = 1",
			a: "C1*cos(x)+C2*sin(x)+1",
			e: "The homogeneous solution: y_c = C₁cos(x) + C₂sin(x) (characteristic r²+1=0, roots ±i). Particular solution: y_p = 1. General: y = C₁cos(x) + C₂sin(x) + 1.",
		},
		{
			q: "Solve: y'' + y = 2",
			a: "C1*cos(x)+C2*sin(x)+2",
			e: "The homogeneous solution: y_c = C₁cos(x) + C₂sin(x). Particular solution: y_p = 2. General: y = C₁cos(x) + C₂sin(x) + 2.",
		},
		{
			q: "Solve: y'' - y = 1",
			a: "C1*e^x+C2*e^(-x)-1",
			e: "The homogeneous solution: y_c = C₁e^x + C₂e^(-x) (characteristic r²-1=0, roots ±1). Particular solution: y_p = -1. General: y = C₁e^x + C₂e^(-x) - 1.",
		},
		{
			q: "Solve: y'' - y = 2",
			a: "C1*e^x+C2*e^(-x)-2",
			e: "The homogeneous solution: y_c = C₁e^x + C₂e^(-x). Particular solution: y_p = -2. General: y = C₁e^x + C₂e^(-x) - 2.",
		},
		{
			q: "Solve: y'' - 4y = 0",
			a: "C1*e^(2x)+C2*e^(-2x)",
			e: "Characteristic equation: r² - 4 = 0, roots r = ±2. General solution: y = C₁e^(2x) + C₂e^(-2x).",
		},
		{
			q: "Solve: y'' - 4y = 1",
			a: "C1*e^(2x)+C2*e^(-2x)-0.25",
			e: "The homogeneous solution: y_c = C₁e^(2x) + C₂e^(-2x). Particular solution: y_p = -1/4. General: y = C₁e^(2x) + C₂e^(-2x) - 1/4.",
		},
	}
	t := templates[rand.Intn(len(templates))]
	return generator.Problem{
		Question:    t.q,
		Answer:      t.a,
		Explanation: t.e,
	}
}

type laplaceGen struct{}

func (g *laplaceGen) Generate(difficulty float64) generator.Problem {
	type qa struct {
		q string
		a string
		e string
	}
	templates := []qa{
		{
			q: "Find L{1}",
			a: "1/s",
			e: "L{1} = ∫₀^∞ e^(-st) dt = 1/s for s > 0.",
		},
		{
			q: "Find L{t}",
			a: "1/s^2",
			e: "L{t} = 1/s² for s > 0.",
		},
		{
			q: "Find L{t²}",
			a: "2/s^3",
			e: "L{t²} = 2/s³ for s > 0.",
		},
	}
	alphas := []int{1, 2, 3}
	a := alphas[rand.Intn(len(alphas))]
	templates = append(templates,
		qa{
			q: fmt.Sprintf("Find L{e^(%dt)}", a),
			a: fmt.Sprintf("1/(s-%d)", a),
			e: fmt.Sprintf("L{e^(%dt)} = 1/(s - %d) for s > %d.", a, a, a),
		},
		qa{
			q: fmt.Sprintf("Find L{sin(%dt)}", a),
			a: fmt.Sprintf("%d/(s^2+%d)", a, a*a),
			e: fmt.Sprintf("L{sin(%dt)} = %d/(s² + %d²) for s > 0.", a, a, a),
		},
		qa{
			q: fmt.Sprintf("Find L{cos(%dt)}", a),
			a: fmt.Sprintf("s/(s^2+%d)", a*a),
			e: fmt.Sprintf("L{cos(%dt)} = s/(s² + %d²) for s > 0.", a, a),
		},
	)

	t := templates[rand.Intn(len(templates))]
	return generator.Problem{
		Question:    t.q,
		Answer:      t.a,
		Explanation: t.e,
	}
}

type systemsGen struct{}

func (g *systemsGen) Generate(difficulty float64) generator.Problem {
	type qa struct {
		q string
		a string
		e string
	}
	templates := []qa{
		{
			q: "A system of 2 first-order ODEs represents how many unknown functions? (enter a number)",
			a: "2",
			e: "Each first-order ODE typically describes the derivative of one unknown function, so 2 equations means 2 unknown functions.",
		},
		{
			q: "A system of 3 first-order ODEs represents how many unknown functions? (enter a number)",
			a: "3",
			e: "Each first-order ODE typically describes the derivative of one unknown function, so 3 equations means 3 unknown functions.",
		},
		{
			q: "dx/dt = 2x, dy/dt = 3y. Are x and y independent? (yes/no)",
			a: "yes",
			e: "The equations are uncoupled: dx/dt depends only on x, and dy/dt depends only on y. Each can be solved independently.",
		},
		{
			q: "dx/dt = x + y, dy/dt = x - y. Is this a coupled system? (yes/no)",
			a: "yes",
			e: "Yes, because dx/dt depends on y and dy/dt depends on x. The equations must be solved together.",
		},
		{
			q: "Can a system of two first-order ODEs be converted into a single second-order ODE? (yes/no)",
			a: "yes",
			e: "Yes, by differentiating one equation and substituting the other, you can eliminate one variable to obtain a single higher-order ODE.",
		},
		{
			q: "Does dx/dt = 2x represent a system of differential equations? (yes/no)",
			a: "no",
			e: "This is a single differential equation, not a system. A system requires at least two equations.",
		},
	}
	t := templates[rand.Intn(len(templates))]
	return generator.Problem{
		Question:    t.q,
		Answer:      t.a,
		Explanation: t.e,
	}
}
