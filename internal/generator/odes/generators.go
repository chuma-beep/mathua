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
}

type conceptGen struct{}

func (g *conceptGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type qa struct {
		q string
		a string
		e string
	}
	templates := []qa{
		{"Is \\(y' = 2x\\) a differential equation? (yes/no)", "yes", "The equation contains a derivative \\(y'\\), so it is a differential equation."},
		{"Is \\(d^{2}y/dx^{2} + y = 0\\) a differential equation? (yes/no)", "yes", "The equation contains a second derivative \\(d^{2}y/dx^{2}\\), so it is a differential equation."},
		{"Is \\(dy/dx = 3x^{2}\\) a differential equation? (yes/no)", "yes", "The equation contains a derivative \\(dy/dx\\), so it is a differential equation."},
		{"Is \\(y''' + 2y'' - y' + y = 0\\) a differential equation? (yes/no)", "yes", "The equation contains third derivatives, so it is a differential equation."},
		{"Is \\(x^{2} + y^{2} = 5\\) a differential equation? (yes/no)", "no", "The equation contains no derivatives, so it is not a differential equation."},
		{"Is \\(3x + 2y = 7\\) a differential equation? (yes/no)", "no", "The equation contains no derivatives, so it is not a differential equation."},
		{"Is \\(\\sin(x) + \\cos(y) = 1\\) a differential equation? (yes/no)", "no", "The equation contains no derivatives, so it is not a differential equation."},
		{"Does the equation \\(x^{3} - y^{3} = 10\\) contain a derivative? (yes/no)", "no", "The equation contains no derivative terms."},
	}
	chosen := templates[rand.Intn(len(templates))]
	return generator.Problem{
		Question:    chosen.q,
		Answer:      chosen.a,
		Explanation: chosen.e,
	}
}

type separableGen struct{}

func (g *separableGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type qa struct {
		q string
		a string
		e string
	}
	templates := []qa{
		{
			q: "Solve: \\(dy/dx = x/y\\)",
			a: "y^2=x^2+C",
			e: "Separate: \\(y dy = x dx\\). Integrate: \\(y^{2}/2 = x^{2}/2 + C\\). Multiply by 2: \\(y^{2} = x^{2} + C\\).",
		},
		{
			q: "Solve: \\(dy/dx = x\\)",
			a: "y=x^2/2+C",
			e: "Rewrite as \\(dy = x dx\\). Integrate: \\(y = x^{2}/2 + C\\).",
		},
		{
			q: "Solve: \\(dy/dx = 1/y\\)",
			a: "y^2=2x+C",
			e: "Separate: \\(y dy = dx\\). Integrate: \\(y^{2}/2 = x + C\\). Multiply by 2: \\(y^{2} = 2x + C\\).",
		},
		{
			q: "Solve: \\(dy/dx = x^{2}/y\\)",
			a: "y^2=(2/3)x^3+C",
			e: "Separate: \\(y dy = x^{2} dx\\). Integrate: \\(y^{2}/2 = x^{3}/3 + C\\). Multiply by 2: \\(y^{2} = (2/3)x^{3} + C\\).",
		},
		{
			q: "Solve: \\(dy/dx = x^{2}\\)",
			a: "y=x^3/3+C",
			e: "Rewrite as \\(dy = x^{2} dx\\). Integrate: \\(y = x^{3}/3 + C\\).",
		},
		{
			q: "Solve: \\(dy/dx = 1/y^{2}\\)",
			a: "y^3=3x+C",
			e: "Separate: \\(y^{2} dy = dx\\). Integrate: \\(y^{3}/3 = x + C\\). Multiply by 3: \\(y^{3} = 3x + C\\).",
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

func (g *linearFirstGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type qa struct {
		q string
		a string
		e string
	}
	templates := []qa{
		{
			q: "Solve: \\(y' + 2y = 0\\)",
			a: "Ce^(-2x)",
			e: "This is a homogeneous first-order linear ODE. The integrating factor is \\(e^{2x}\\). \\(\\frac{d}{dx}(ye^{2x}) = 0\\), so \\(ye^{2x} = C\\), giving \\(y = Ce^{-2x}\\).",
		},
		{
			q: "Solve: \\(y' + y = 0\\)",
			a: "Ce^(-x)",
			e: "This is a homogeneous first-order linear ODE. The integrating factor is \\(e^{x}\\). \\(\\frac{d}{dx}(ye^{x}) = 0\\), so \\(ye^{x} = C\\), giving \\(y = Ce^{-x}\\).",
		},
		{
			q: "Solve: \\(y' + 3y = 0\\)",
			a: "Ce^(-3x)",
			e: "This is a homogeneous first-order linear ODE. The integrating factor is \\(e^{3x}\\). \\(\\frac{d}{dx}(ye^{3x}) = 0\\), so \\(ye^{3x} = C\\), giving \\(y = Ce^{-3x}\\).",
		},
		{
			q: "Solve: \\(y' + y = 1\\)",
			a: "1+Ce^(-x)",
			e: "The integrating factor is \\(e^{x}\\). \\(\\frac{d}{dx}(ye^{x}) = e^{x}\\). Integrate: \\(ye^{x} = e^{x} + C\\), so \\(y = 1 + Ce^{-x}\\).",
		},
		{
			q: "Solve: \\(y' + 2y = 4\\)",
			a: "2+Ce^(-2x)",
			e: "The integrating factor is \\(e^{2x}\\). \\(\\frac{d}{dx}(ye^{2x}) = 4e^{2x}\\). Integrate: \\(ye^{2x} = 2e^{2x} + C\\), so \\(y = 2 + Ce^{-2x}\\).",
		},
		{
			q: "Solve: \\(y' + y = e^{x}\\)",
			a: "0.5e^x+Ce^(-x)",
			e: "The integrating factor is \\(e^{x}\\). \\(\\frac{d}{dx}(ye^{x}) = e^{2x}\\). Integrate: \\(ye^{x} = \\frac{1}{2}e^{2x} + C\\), so \\(y = \\frac{1}{2}e^{x} + Ce^{-x}\\).",
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

func (g *exactGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type qa struct {
		q string
		a string
		e string
	}
	templates := []qa{
		{
			q: "Is \\(M dx + N dy = 0\\) exact when \\(M = 2xy\\) and \\(N = x^{2}\\)? (yes/no)",
			a: "yes",
			e: "\\(\\partial M/\\partial y = 2x\\) and \\(\\partial N/\\partial x = 2x\\). Since \\(\\partial M/\\partial y = \\partial N/\\partial x\\), the equation is exact.",
		},
		{
			q: "Is \\(M dx + N dy = 0\\) exact when \\(M = y\\) and \\(N = x\\)? (yes/no)",
			a: "yes",
			e: "\\(\\partial M/\\partial y = 1\\) and \\(\\partial N/\\partial x = 1\\). Since \\(\\partial M/\\partial y = \\partial N/\\partial x\\), the equation is exact.",
		},
		{
			q: "Is \\(M dx + N dy = 0\\) exact when \\(M = 3x^{2}y\\) and \\(N = x^{3}\\)? (yes/no)",
			a: "yes",
			e: "\\(\\partial M/\\partial y = 3x^{2}\\) and \\(\\partial N/\\partial x = 3x^{2}\\). Since \\(\\partial M/\\partial y = \\partial N/\\partial x\\), the equation is exact.",
		},
		{
			q: "Is \\(M dx + N dy = 0\\) exact when \\(M = x^{2}\\) and \\(N = y^{2}\\)? (yes/no)",
			a: "no",
			e: "\\(\\partial M/\\partial y = 0\\) and \\(\\partial N/\\partial x = 0\\). Since \\(\\partial M/\\partial y = \\partial N/\\partial x = 0\\), this equation IS exact. (Trick: both partials are zero.)",
		},
		{
			q: "Is \\(M dx + N dy = 0\\) exact when \\(M = x\\) and \\(N = x\\)? (yes/no)",
			a: "no",
			e: "\\(\\partial M/\\partial y = 0\\) and \\(\\partial N/\\partial x = 1\\). Since \\(\\partial M/\\partial y \\neq \\partial N/\\partial x\\) \\((0 \\neq 1)\\), the equation is not exact.",
		},
		{
			q: "Is \\(M dx + N dy = 0\\) exact when \\(M = y^{2}\\) and \\(N = x\\)? (yes/no)",
			a: "no",
			e: "\\(\\partial M/\\partial y = 2y\\) and \\(\\partial N/\\partial x = 1\\). Since \\(\\partial M/\\partial y \\neq \\partial N/\\partial x\\) \\((2y \\neq 1\\) in general), the equation is not exact.",
		},
		{
			q: "Is \\(M dx + N dy = 0\\) exact when \\(M = xy\\) and \\(N = xy\\)? (yes/no)",
			a: "no",
			e: "\\(\\partial M/\\partial y = x\\) and \\(\\partial N/\\partial x = y\\). Since \\(\\partial M/\\partial y \\neq \\partial N/\\partial x\\) \\((x \\neq y\\) in general), the equation is not exact.",
		},
		{
			q: "Is \\(M dx + N dy = 0\\) exact when \\(M = \\sin(y)\\) and \\(N = x\\cos(y)\\)? (yes/no)",
			a: "yes",
			e: "\\(\\partial M/\\partial y = \\cos(y)\\) and \\(\\partial N/\\partial x = \\cos(y)\\). Since \\(\\partial M/\\partial y = \\partial N/\\partial x\\), the equation is exact.",
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

func (g *homogeneousGen) Generate(ctx generator.GeneratorContext) generator.Problem {
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

	charStr := "r^{2}"
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

	explanation := fmt.Sprintf("Characteristic equation: \\(%s\\). Roots: \\(r = %d\\) and \\(r = %d\\). The general solution is \\(y = C_{1}e^{%dx} + C_{2}e^{%dx}\\).", charStr, r1, r2, r1, r2)

	return generator.Problem{
		Question:    fmt.Sprintf("Solve: \\(%s\\)", eqStr),
		Answer:      answer,
		Explanation: explanation,
	}
}

type nonhomogeneousGen struct{}

func (g *nonhomogeneousGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type qa struct {
		q string
		a string
		e string
	}
	templates := []qa{
		{
			q: "Solve: \\(y'' + y = 1\\)",
			a: "C1*cos(x)+C2*sin(x)+1",
			e: "The homogeneous solution: \\(y_c = C_{1}\\cos(x) + C_{2}\\sin(x)\\) (characteristic \\(r^{2}+1=0\\), roots \\(\\pm i\\)). Particular solution: \\(y_p = 1\\). General: \\(y = C_{1}\\cos(x) + C_{2}\\sin(x) + 1\\).",
		},
		{
			q: "Solve: \\(y'' + y = 2\\)",
			a: "C1*cos(x)+C2*sin(x)+2",
			e: "The homogeneous solution: \\(y_c = C_{1}\\cos(x) + C_{2}\\sin(x)\\). Particular solution: \\(y_p = 2\\). General: \\(y = C_{1}\\cos(x) + C_{2}\\sin(x) + 2\\).",
		},
		{
			q: "Solve: \\(y'' - y = 1\\)",
			a: "C1*e^x+C2*e^(-x)-1",
			e: "The homogeneous solution: \\(y_c = C_{1}e^{x} + C_{2}e^{-x}\\) (characteristic \\(r^{2}-1=0\\), roots \\(\\pm 1\\)). Particular solution: \\(y_p = -1\\). General: \\(y = C_{1}e^{x} + C_{2}e^{-x} - 1\\).",
		},
		{
			q: "Solve: \\(y'' - y = 2\\)",
			a: "C1*e^x+C2*e^(-x)-2",
			e: "The homogeneous solution: \\(y_c = C_{1}e^{x} + C_{2}e^{-x}\\). Particular solution: \\(y_p = -2\\). General: \\(y = C_{1}e^{x} + C_{2}e^{-x} - 2\\).",
		},
		{
			q: "Solve: \\(y'' - 4y = 0\\)",
			a: "C1*e^(2x)+C2*e^(-2x)",
			e: "Characteristic equation: \\(r^{2} - 4 = 0\\), roots \\(r = \\pm 2\\). General solution: \\(y = C_{1}e^{2x} + C_{2}e^{-2x}\\).",
		},
		{
			q: "Solve: \\(y'' - 4y = 1\\)",
			a: "C1*e^(2x)+C2*e^(-2x)-0.25",
			e: "The homogeneous solution: \\(y_c = C_{1}e^{2x} + C_{2}e^{-2x}\\). Particular solution: \\(y_p = -1/4\\). General: \\(y = C_{1}e^{2x} + C_{2}e^{-2x} - 1/4\\).",
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

func (g *laplaceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type qa struct {
		q string
		a string
		e string
	}
	templates := []qa{
		{
			q: "Find \\(L\\{1\\}\\)",
			a: "1/s",
			e: "\\(L\\{1\\} = \\int_{0}^{\\infty} e^{-st} dt = 1/s\\) for \\(s > 0\\).",
		},
		{
			q: "Find \\(L\\{t\\}\\)",
			a: "1/s^2",
			e: "\\(L\\{t\\} = 1/s^{2}\\) for \\(s > 0\\).",
		},
		{
			q: "Find \\(L\\{t^{2}\\}\\)",
			a: "2/s^3",
			e: "\\(L\\{t^{2}\\} = 2/s^{3}\\) for \\(s > 0\\).",
		},
	}
	alphas := []int{1, 2, 3}
	a := alphas[rand.Intn(len(alphas))]
	templates = append(templates,
		qa{
			q: fmt.Sprintf("Find \\(L\\{e^{%dt}\\}\\)", a),
			a: fmt.Sprintf("1/(s-%d)", a),
			e: fmt.Sprintf("\\(L\\{e^{%dt}\\} = 1/(s - %d)\\) for \\(s > %d\\).", a, a, a),
		},
		qa{
			q: fmt.Sprintf("Find \\(L\\{\\sin(%dt)\\}\\)", a),
			a: fmt.Sprintf("%d/(s^2+%d)", a, a*a),
			e: fmt.Sprintf("\\(L\\{\\sin(%dt)\\} = %d/(s^{2} + %d^{2})\\) for \\(s > 0\\).", a, a, a),
		},
		qa{
			q: fmt.Sprintf("Find \\(L\\{\\cos(%dt)\\}\\)", a),
			a: fmt.Sprintf("s/(s^2+%d)", a*a),
			e: fmt.Sprintf("\\(L\\{\\cos(%dt)\\} = s/(s^{2} + %d^{2})\\) for \\(s > 0\\).", a, a),
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

func (g *systemsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
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
			q: "\\(dx/dt = 2x\\), \\(dy/dt = 3y\\). Are \\(x\\) and \\(y\\) independent? (yes/no)",
			a: "yes",
			e: "The equations are uncoupled: \\(dx/dt\\) depends only on \\(x\\), and \\(dy/dt\\) depends only on \\(y\\). Each can be solved independently.",
		},
		{
			q: "\\(dx/dt = x + y\\), \\(dy/dt = x - y\\). Is this a coupled system? (yes/no)",
			a: "yes",
			e: "Yes, because \\(dx/dt\\) depends on \\(y\\) and \\(dy/dt\\) depends on \\(x\\). The equations must be solved together.",
		},
		{
			q: "Can a system of two first-order ODEs be converted into a single second-order ODE? (yes/no)",
			a: "yes",
			e: "Yes, by differentiating one equation and substituting the other, you can eliminate one variable to obtain a single higher-order ODE.",
		},
		{
			q: "Does \\(dx/dt = 2x\\) represent a system of differential equations? (yes/no)",
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

type variationParamsGen struct{}

func (g *variationParamsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type qa struct {
		q string
		a string
		e string
	}
	templates := []qa{
		{q: "For \\(y'' + y = \\sec(x)\\), the homogeneous solutions are \\(y_1=\\cos x\\), \\(y_2=\\sin x\\). What is the Wronskian \\(W\\)? (enter a number)", a: "1", e: "\\(W = y_1y_2' - y_1'y_2 = \\cos x \\cdot \\cos x - (-\\sin x)\\sin x = 1\\)."},
		{q: "In variation of parameters, \\(y_p = u_1 y_1 + u_2 y_2\\) where \\(u_1' = -y_2 g/W\\). For \\(y''+y=1\\) with \\(W=1\\), \\(y_2=\\sin x\\), what is \\(u_1'\\)?", a: "-sin(x)", e: "\\(u_1' = -y_2 g/W = -\\sin x \\cdot 1/1 = -\\sin x\\)."},
		{q: "Does variation of parameters require knowing the homogeneous solution? (yes/no)", a: "yes", e: "Yes, the method builds the particular solution from the homogeneous basis \\(y_1, y_2\\) and the Wronskian."},
		{q: "For constant-coefficient ODEs, variation of parameters and undetermined coefficients give the same particular solution when both apply. Is this true? (yes/no)", a: "yes", e: "Both methods recover a particular solution; variation of parameters is more general and works even when undetermined coefficients does not."},
	}
	t := templates[rand.Intn(len(templates))]
	return generator.Problem{Question: t.q, Answer: t.a, Explanation: t.e}
}

type convolutionGen struct{}

func (g *convolutionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type qa struct {
		q string
		a string
		e string
	}
	templates := []qa{
		{q: "The convolution \\((f*g)(t) = \\int_0^t f(\\tau)g(t-\\tau)d\\tau\\). What is \\((1*1)(t)\\)?", a: "t", e: "\\((1*1)(t)=\\int_0^t 1\\cdot1 d\\tau = t\\)."},
		{q: "\\(L\\{f*g\\} = L\\{f\\} \\cdot L\\{g\\}\\). If \\(L\\{f\\}=1/s\\) and \\(L\\{g\\}=1/s\\), what is \\(L\\{f*g\\}\\)?", a: "1/s^2", e: "By the convolution theorem, \\(L\\{f*g\\}= (1/s)(1/s)=1/s^2\\)."},
		{q: "Is convolution commutative: \\(f*g = g*f\\)? (yes/no)", a: "yes", e: "Yes, by substitution \\(\\tau \\to t-\\tau\\) the integral is symmetric."},
		{q: "If \\(f(t)=e^{t}\\) and \\(g(t)=1\\), \\((f*g)(t)=e^{t}-1\\). Does \\(L\\{f*g\\}=1/(s(s-1))\\) hold? (yes/no)", a: "yes", e: "\\(L\\{e^t\\}=1/(s-1)\\), \\(L\\{1\\}=1/s\\), product is \\(1/(s(s-1))=L\\{e^t-1\\}\\)."},
	}
	t := templates[rand.Intn(len(templates))]
	return generator.Problem{Question: t.q, Answer: t.a, Explanation: t.e}
}

type linearPhaseGen struct{}

func (g *linearPhaseGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type qa struct {
		q string
		a string
		e string
	}
	templates := []qa{
		{q: "For \\(x' = Ax\\) with \\(A=[[2,0],[0,-3]]\\), the eigenvalues are \\(2\\) and \\(-3\\). What type is the origin? (saddle/node/spiral)", a: "saddle", e: "One positive, one negative eigenvalue gives a saddle."},
		{q: "For \\(x' = Ax\\) with eigenvalues \\(-1, -2\\), what type is the origin? (saddle/node/spiral)", a: "node", e: "Both eigenvalues negative and real gives a stable node."},
		{q: "If eigenvalues are \\(0.5 \\pm 2i\\), what type is the origin? (saddle/node/spiral)", a: "spiral", e: "Complex pair with positive real part gives an unstable spiral."},
		{q: "Does the phase portrait of a linear system depend on eigenvalues of \\(A\\)? (yes/no)", a: "yes", e: "The eigenvalues classify the equilibrium: node, saddle, spiral, or centre."},
	}
	t := templates[rand.Intn(len(templates))]
	return generator.Problem{Question: t.q, Answer: t.a, Explanation: t.e}
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
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		q string
		a string
		e string
	}
	easy := []entry{
		{"Does RK4 have local truncation error O(h^5) and global O(h^4)? (yes/no)", "yes", "RK4: local O(h^5), global O(h^4), vs Euler local O(h^2) global O(h)."},
		{"How many function evaluations per step does RK4 need? (4)", "4", "k1,k2,k3,k4 per step."},
		{"Is RK4 more accurate than Euler with same h? (yes/no)", "yes", "Higher order gives smaller error for same step size."},
	}
	hard := []entry{
		{"With h=0.1, does halving h reduce RK4 global error by ~16×? (yes/no)", "yes", "Error ∝ h^4, so (1/2)^4=1/16."},
		{"Does RK4's increment use weighted average (k1+2k2+2k3+k4)/6? (yes/no)", "yes", "Classic RK4 formula."},
		{"Is RK4 explicit (no solve) for y'=f(x,y)? (yes/no)", "yes", "Explicit; all k_i use known values."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type stabilityGen struct{}

func (g *stabilityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		q string
		a string
		e string
	}
	easy := []entry{
		{"For x'=-x, is equilibrium 0 stable? (yes/no)", "yes", "Solution e^{-t}x0 →0, asymptotically stable."},
		{"For x'=x, is 0 unstable? (yes/no)", "yes", "Solutions e^{t}x0 diverge from 0."},
		{"Does negative eigenvalue give stable node in linear system? (yes/no)", "yes", "Both eigenvalues negative → stable node."},
	}
	hard := []entry{
		{"Does Lyapunov's indirect method use Jacobian eigenvalues at equilibrium? (yes/no)", "yes", "Linearization stability via eigenvalues."},
		{"Is centre (pure imaginary eigenvalues) stable but not asymptotically? (yes/no)", "yes", "Orbits are closed, not approaching nor diverging."},
		{"Does asymptotically stable imply stable? (yes/no)", "yes", "Asymptotic is stronger."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type characteristicsGen struct{}

func (g *characteristicsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		q string
		a string
		e string
	}
	easy := []entry{
		{"Does method of characteristics reduce first-order PDE a*u_x+b*u_y=0 to ODEs along curves? (yes/no)", "yes", "Characteristics satisfy dx/a=dy/b."},
		{"For u_t + c u_x=0, are characteristics lines x-ct=const? (yes/no)", "yes", "Solution u(x,t)=f(x-ct) constant along x-ct."},
		{"Is transport PDE solved by shifting initial data along characteristics? (yes/no)", "yes", "Value propagates along char lines."},
	}
	hard := []entry{
		{"Does quasi-linear PDE become ODE system for (x(t),u(t)) via characteristics? (yes/no)", "yes", "dx/dt=a, dy/dt=b, du/dt=c."},
		{"Are characteristics curves where PDE becomes interior ODE? (yes/no)", "yes", "PDE restricts to ODE along them."},
		{"Does shock form when characteristics intersect? (yes/no)", "yes", "Intersection gives multi-valued solution → shock."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type fourierGen struct{}

func (g *fourierGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		q string
		a string
		e string
	}
	easy := []entry{
		{"Does Fourier transform convert differentiation to multiplication by iω? (yes/no)", "yes", "F{d/dx}=iω."},
		{"Does Fourier solve ODEs with constant coefficients by algebra in frequency domain? (yes/no)", "yes", "Transform, solve algebraic equation, invert."},
		{"Is Fourier transform of convolution a product? (yes/no)", "yes", "Like Laplace, F{f*g}=F{f}·F{g}."},
	}
	hard := []entry{
		{"Does Fourier method require solving ( -ω^2 + ...)U = F? (yes/no)", "yes", "E.g. y''+y=f → (-ω^2+1)Y=F."},
		{"Are Fourier and Laplace both integral transforms for ODEs? (yes/no)", "yes", "Both convert differential to algebraic."},
		{"Does inverse Fourier require contour integration? (yes/no)", "yes", "Inversion integral over real line."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type bifurcationGen struct{}

func (g *bifurcationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		q string
		a string
		e string
	}
	easy := []entry{
		{"For x'= μx - x^3, does bifurcation at μ=0 change stability? (yes/no)", "yes", "Pitchfork: 0 stable for μ<0, unstable for μ>0 with two new stable branches."},
		{"Does saddle-node bifurcation create/destroy equilibria as parameter varies? (yes/no)", "yes", "x'= μ - x^2 has 0,1,2 equilibria depending on μ."},
		{"Is bifurcation a qualitative change in phase portrait as parameter crosses threshold? (yes/no)", "yes", "Definition of bifurcation."},
	}
	hard := []entry{
		{"For x'= μx - x^2, is it transcritical? (yes/no)", "yes", "Two equilibria exchange stability at μ=0."},
		{"Does Hopf bifurcation create a limit cycle? (yes/no)", "yes", "Pair of complex eigenvalues cross imaginary axis, periodic orbit appears."},
		{"Is pitchfork normal form x'= μx - x^3 supercritical when cubic term negative? (yes/no)", "yes", "Stable branches emerge for μ>0."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type waveGen struct{}

func (g *waveGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		q string
		a string
		e string
	}
	easy := []entry{
		{"Does wave equation u_tt = c^2 u_xx have d'Alembert solution u(x,t)=f(x-ct)+g(x+ct)? (yes/no)", "yes", "Superposition of right and left travelling waves."},
		{"Is wave equation hyperbolic? (yes/no)", "yes", "Classification via discriminant."},
		{"Do characteristics for wave equation have slopes ±c? (yes/no)", "yes", "Lines x±ct=const."},
	}
	hard := []entry{
		{"Does d'Alembert with initial u(x,0)=f(x), u_t(x,0)=0 give u(x,t)=[f(x-ct)+f(x+ct)]/2? (yes/no)", "yes", "Even reflection of initial displacement."},
		{"Is domain of dependence for wave equation an interval [x-ct,x+ct]? (yes/no)", "yes", "Value at (x,t) depends on initial data there."},
		{"Does wave equation conserve energy ∫(u_t^2+c^2 u_x^2)dx? (yes/no)", "yes", "Energy conservation."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type heatGen struct{}

func (g *heatGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		q string
		a string
		e string
	}
	easy := []entry{
		{"Does heat equation u_t = k u_xx have solution via separation u=X(x)T(t)? (yes/no)", "yes", "Yields X''/X = T'/(kT) = -λ."},
		{"Is heat equation parabolic? (yes/no)", "yes", "Classification."},
		{"Does separation give T(t)=C e^{-kλt}? (yes/no)", "yes", "ODE T'=-kλT."},
	}
	hard := []entry{
		{"Does Fourier series solve heat equation on [0,L] with Dirichlet? (yes/no)", "yes", "Eigenfunctions sin(nπx/L)."},
		{"Is maximum principle: max of u on space-time boundary? (yes/no)", "yes", "Heat cannot have interior max exceeding boundary."},
		{"Does heat smoothing make incompatible initial data become smooth for t>0? (yes/no)", "yes", "Instant smoothing."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type greenFunctionGen struct{}

func (g *greenFunctionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		q string
		a string
		e string
	}
	easy := []entry{
		{"Does Green's function solve L G = δ with boundary conditions? (yes/no)", "yes", "Impulse response; y(x)=∫G(x,ξ)f(ξ)dξ."},
		{"Is Green's function the inverse of differential operator? (yes/no)", "yes", "Integral operator inverts L."},
		{"Does Green's function for y''=f with y(0)=y(1)=0 equal piecewise linear? (yes/no)", "yes", "G(x,ξ)=x(1-ξ) for x<ξ, ξ(1-x) for x>ξ."},
	}
	hard := []entry{
		{"Does jump condition give G' discontinuity of 1 at x=ξ? (yes/no)", "yes", "Integrating across δ jump."},
		{"Is Green's function symmetric for self-adjoint L? (yes/no)", "yes", "Reciprocity G(x,ξ)=G(ξ,x)."},
		{"Does eigenfunction expansion give G(x,ξ)=∑ φ_n(x)φ_n(ξ)/λ_n? (yes/no)", "yes", "Spectral representation."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type sturmLiouvilleGen struct{}

func (g *sturmLiouvilleGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct{ q, a, e string }
	easy := []entry{{"Is Sturm-Liouville operator self-adjoint? (yes/no)", "yes", "SL is self-adjoint under weight."}, {"Does SL have real eigenvalues? (yes/no)", "yes", "Self-adjoint gives real."}, {"Are eigenfunctions orthogonal with weight? (yes/no)", "yes", "Orthogonal."}}
	hard := []entry{{"Does SL on [a,b] with separated BC have discrete spectrum? (yes/no)", "yes", "Eigenvalues infinite discrete."}, {"Is Legendre equation Sturm-Liouville? (yes/no)", "yes", "With p=1-x^2."}, {"Does completeness of eigenfunctions hold? (yes/no)", "yes", "Expand in eigenfunctions."}}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type boundaryValueGen struct{}

func (g *boundaryValueGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct{ q, a, e string }
	easy := []entry{{"Does BVP y''=0, y(0)=0,y(1)=1 have unique solution y=x? (yes/no)", "yes", "Linear BVP."}, {"Is BVP different from IVP (conditions at two points)? (yes/no)", "yes", "BVP at boundaries."}, {"Does BVP may have 0,1,∞ solutions? (yes/no)", "yes", "Unlike IVP."}}
	hard := []entry{{"Does Green's function solve BVP? (yes/no)", "yes", "Via integral."}, {"Is shooting method for BVP? (yes/no)", "yes", "Convert BVP to IVP iteration."}, {"Does Fredholm alternative apply? (yes/no)", "yes", "Solvability condition."}}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type lyapunovGen struct{}

func (g *lyapunovGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct{ q, a, e string }
	easy := []entry{{"Is V(x)=x^2 Lyapunov for x'=-x? (yes/no)", "yes", "V>0, dV/dt<0."}, {"Does Lyapunov prove stability without solving ODE? (yes/no)", "yes", "Energy-like function."}, {"Is Lyapunov's direct method? (yes/no)", "yes", "Direct."}}
	hard := []entry{{"Does LaSalle's invariance principle extend Lyapunov? (yes/no)", "yes", "Invariant set."}, {"Is V strict Lyapunov if dV/dt<0? (yes/no)", "yes", "Strict."}, {"Does converse Lyapunov hold for asymptotically stable? (yes/no)", "yes", "Existence of V."}}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type laplaceEqGen struct{}

func (g *laplaceEqGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct{ q, a, e string }
	easy := []entry{{"Is Laplace equation Δu=0? (yes/no)", "yes", "Laplace."}, {"Does Laplace solution have mean value property? (yes/no)", "yes", "Harmonic."}, {"Is Dirichlet problem Laplace with boundary values? (yes/no)", "yes", "Boundary condition."}}
	hard := []entry{{"Does maximum principle hold for Laplace? (yes/no)", "yes", "Interior max = boundary max."}, {"Is fundamental solution log| x| in 2D? (yes/no)", "yes", "Newtonian potential."}, {"Does separation give harmonic functions? (yes/no)", "yes", "Via eigenfunctions."}}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type volterraGen struct{}

func (g *volterraGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct{ q, a, e string }
	easy := []entry{{"Is Volterra equation y(x)=f(x)+∫_0^x K(x,t)y(t)dt? (yes/no)", "yes", "Volterra of second kind."}, {"Does Volterra have unique continuous solution? (yes/no)", "yes", "Via Picard iteration."}, {"Is Volterra with K=1 solvable via differentiation? (yes/no)", "yes", "Differentiate to ODE."}}
	hard := []entry{{"Does resolvent kernel give solution? (yes/no)", "yes", "Neumann series."}, {"Is Volterra compact operator? (yes/no)", "yes", "Integral compact."}, {"Does Laplace solve convolution Volterra? (yes/no)", "yes", "Transform."}}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}
