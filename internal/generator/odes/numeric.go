package odes

// The numeric differential-equations concepts: an initial value problem
// evaluated at a point, the sum of a characteristic equation's roots, a non-zero equilibrium,
// the Wronskian of two powers, the number of Euler steps between two points, the Laplace transform
// of a power evaluated at a point and the discriminant that classifies a second-order PDE.
//
// Converted to parametric generators scaling with ctx.Difficulty using mathutil primitives
// and structured three-subgoal cognitive derivations in explanations.

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/mathutil"
)

func registerNumeric(reg *generator.Registry) {
	reg.Register("ode.ivp.evaluate", &ivpEvaluateGen{})
	reg.Register("ode.characteristic.sum", &characteristicSumGen{})
	reg.Register("ode.equilibrium.nonzero", &equilibriumGen{})
	reg.Register("ode.wronskian", &wronskianGen{})
	reg.Register("ode.numerical.steps", &eulerStepsGen{})
	reg.Register("ode.laplace.value", &laplaceValueGen{})
	reg.Register("ode.pde.discriminant", &pdeDiscriminantGen{})
}

// ivpEvaluateGen integrates a linear forcing and evaluates the solution at a point.
type ivpEvaluateGen struct{}

func (g *ivpEvaluateGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	maxA := mathutil.ScaleInt(ctx.Difficulty, 3, 8)
	aHalf := 1 + rand.Intn(maxA)
	a := 2 * aHalf // even so a/2 is whole
	maxB := mathutil.ScaleInt(ctx.Difficulty, 4, 12)
	b := rand.Intn(2*maxB+1) - maxB
	if b == 0 {
		b = 1
	}
	maxC := mathutil.ScaleInt(ctx.Difficulty, 5, 15)
	c := rand.Intn(2*maxC+1) - maxC
	if c == 0 {
		c = 2
	}
	maxT := mathutil.ScaleInt(ctx.Difficulty, 2, 6)
	t := 1 + rand.Intn(maxT)

	v := aHalf*t*t + b*t + c
	ans := strconv.Itoa(v)
	q := fmt.Sprintf("Solve \\(y' = %dx %+d\\) with \\(y(0) = %d\\) and find \\(y(%d)\\).", a, b, c, t)
	exp := fmt.Sprintf("Subgoal 1: State governing formula. The solution to a first-order linear differential equation y' = f(x) is obtained by indefinite integration y(x) = ∫ f(x) dx + C, with integration constant C determined by the initial condition y(0).\n\nSubgoal 2: Substitute parameters. Integrate the right side to obtain y(x) = (%d/2)x^2 %+dx + C = %dx^2 %+dx + C. Applying the initial condition y(0) = %d gives C = %d, so the particular solution is y(x) = %dx^2 %+dx %+d.\n\nSubgoal 3: Evaluate and simplify. Substitute x = %d into the particular solution: y(%d) = %d(%d)^2 %+d(%d) %+d = %d.",
		a, b, aHalf, b, c, c, aHalf, b, c, t, t, aHalf, t, b, t, c, v)

	return generator.Problem{
		Question:    q,
		Answer:      ans,
		Explanation: exp,
		Facts: map[string]string{
			"a":      strconv.Itoa(a),
			"b":      strconv.Itoa(b),
			"c":      strconv.Itoa(c),
			"t":      strconv.Itoa(t),
			"answer": ans,
		},
	}
}

// characteristicSumGen uses Vieta's formulas to compute the sum of roots.
type characteristicSumGen struct{}

func (g *characteristicSumGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	maxA := mathutil.ScaleInt(ctx.Difficulty, 2, 6)
	a := 2 + rand.Intn(maxA)
	maxK := mathutil.ScaleInt(ctx.Difficulty, 3, 8)
	k := rand.Intn(2*maxK+1) - maxK
	if k == 0 {
		k = 2
	}
	b := a * k
	maxC := mathutil.ScaleInt(ctx.Difficulty, 5, 20)
	c := 1 + rand.Intn(maxC)
	sum := -b / a
	ans := strconv.Itoa(sum)

	q := fmt.Sprintf("For the equation \\(%dy'' %+dy' %+dy = 0\\), what is the sum of the roots of the characteristic equation?", a, b, c)
	exp := fmt.Sprintf("Subgoal 1: State governing formula. For a second-order linear homogeneous differential equation with constant coefficients ar^2 + br + c = 0, Vieta's formulas state that the sum of the characteristic roots satisfies r₁ + r₂ = -b/a.\n\nSubgoal 2: Substitute parameters. The characteristic polynomial for the differential equation is %dr^2 %+dr %+d = 0 with leading coefficient a = %d and linear coefficient b = %d.\n\nSubgoal 3: Evaluate and simplify. Compute the root sum: r₁ + r₂ = -(%d)/%d = %d.",
		a, b, c, a, b, b, a, sum)

	return generator.Problem{
		Question:    q,
		Answer:      ans,
		Explanation: exp,
		Facts: map[string]string{
			"a":      strconv.Itoa(a),
			"b":      strconv.Itoa(b),
			"c":      strconv.Itoa(c),
			"answer": ans,
		},
	}
}

// equilibriumGen asks for the non-zero equilibrium of a logistic form whose
// capacity is not printed as the answer.
//
// It used to draw the capacity a and ask for it, so every one of 500 sampled
// draws printed its own answer in the stem: the question tested copying, not
// solving. Writing the model as y' = a y - b y^2 keeps the same equilibrium set
// and makes the non-zero one a/b, a division the learner has to perform.
type equilibriumGen struct{}

func (g *equilibriumGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	maxB := mathutil.ScaleInt(ctx.Difficulty, 4, 12)
	b := 2 + rand.Intn(maxB)
	maxQ := mathutil.ScaleInt(ctx.Difficulty, 8, 30)
	q := 1 + rand.Intn(maxQ)
	a := b * q
	ans := strconv.Itoa(q)

	question := fmt.Sprintf("The population model \\(y' = %dy - %dy^{2}\\) has two equilibria. What is the non-zero one?", a, b)
	exp := fmt.Sprintf("Subgoal 1: State governing formula. An equilibrium point y* of an autonomous differential equation y' = f(y) satisfies f(y*) = 0.\n\nSubgoal 2: Substitute parameters. Factor the right side: %dy - %dy^2 = y(%d - %dy). The product vanishes at y = 0 and where %d - %dy = 0.\n\nSubgoal 3: Evaluate and simplify. Solve %d - %dy = 0 to get y = %d/%d = %d, the non-zero equilibrium.",
		a, b, a, b, a, b, a, b, a, b, q)

	return generator.Problem{
		Question:    question,
		Answer:      ans,
		Explanation: exp,
		Facts: map[string]string{
			"a":      strconv.Itoa(a),
			"b":      strconv.Itoa(b),
			"answer": ans,
		},
	}
}

// wronskianGen evaluates the Wronskian of two power functions at a point.
type wronskianGen struct{}

func (g *wronskianGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	maxM := mathutil.ScaleInt(ctx.Difficulty, 2, 4)
	m := rand.Intn(maxM)
	maxDelta := mathutil.ScaleInt(ctx.Difficulty, 3, 5)
	delta := 1 + rand.Intn(maxDelta)
	n := m + delta
	// The evaluation point starts at two: at x = 1 the answer is just n - m,
	// which throws away the exponent and collapses distinct pairs onto one value.
	maxX := mathutil.ScaleInt(ctx.Difficulty, 2, 3)
	x := 2 + rand.Intn(maxX)
	expVal := m + n - 1
	v := delta * powInt(x, expVal)
	ans := strconv.Itoa(v)

	q := fmt.Sprintf("Find the Wronskian \\(W(x^{%d}, x^{%d})\\) evaluated at \\(x = %d\\).", m, n, x)
	exp := fmt.Sprintf("Subgoal 1: State governing formula. The Wronskian of two differentiable functions f(x) and g(x) is defined by the determinant W(f, g)(x) = f(x)g'(x) - f'(x)g(x).\n\nSubgoal 2: Substitute parameters. Differentiating f(x) = x^{%d} gives f'(x) = %d x^{%d}. Differentiating g(x) = x^{%d} gives g'(x) = %d x^{%d}. The Wronskian is W(x) = x^{%d}(%d x^{%d}) - (%d x^{%d})x^{%d} = (%d - %d)x^{%d} = %d x^{%d}.\n\nSubgoal 3: Evaluate and simplify. Evaluate the expression at x = %d: W(%d) = %d(%d)^{%d} = %d.",
		m, m, m-1, n, n, n-1, m, n, n-1, m, m-1, n, n, m, expVal, delta, expVal, x, x, delta, x, expVal, v)

	return generator.Problem{
		Question:    q,
		Answer:      ans,
		Explanation: exp,
		Facts: map[string]string{
			"m":      strconv.Itoa(m),
			"n":      strconv.Itoa(n),
			"x":      strconv.Itoa(x),
			"answer": ans,
		},
	}
}

// eulerStepsGen asks how many steps a step size needs to cross an interval.
type eulerStepsGen struct{}

func (g *eulerStepsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	maxH := mathutil.ScaleInt(ctx.Difficulty, 2, 6)
	h := 1 + rand.Intn(maxH)
	maxSteps := mathutil.ScaleInt(ctx.Difficulty, 6, 30)
	steps := 2 + rand.Intn(maxSteps)
	maxA := mathutil.ScaleInt(ctx.Difficulty, 4, 15)
	a := rand.Intn(maxA + 1)
	b := a + h*steps
	ans := strconv.Itoa(steps)

	q := fmt.Sprintf("Euler's method runs from \\(x = %d\\) to \\(x = %d\\) with step size \\(h = %d\\). How many steps does it take?", a, b, h)
	exp := fmt.Sprintf("Subgoal 1: State governing formula. The number of uniform discretization steps N across an interval [a, b] with step size h is given by N = (b - a)/h.\n\nSubgoal 2: Substitute parameters. Set lower bound a = %d, upper bound b = %d and step size h = %d into the step count formula.\n\nSubgoal 3: Evaluate and simplify. Compute N = (%d - %d)/%d = %d/%d = %d steps.",
		a, b, h, b, a, h, b-a, h, steps)

	return generator.Problem{
		Question:    q,
		Answer:      ans,
		Explanation: exp,
		Facts: map[string]string{
			"a":      strconv.Itoa(a),
			"b":      strconv.Itoa(b),
			"h":      strconv.Itoa(h),
			"answer": ans,
		},
	}
}

// laplaceValueGen evaluates the Laplace transform of a monomial at s = 1.
type laplaceValueGen struct{}

func (g *laplaceValueGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// The power is drawn from a range that widens with difficulty. It used to
	// span four values in total, so a learner reviewing this concept saw the
	// same factorial asked four ways and nothing else.
	minN := mathutil.ScaleInt(ctx.Difficulty, 2, 5)
	maxN := mathutil.ScaleInt(ctx.Difficulty, 5, 9)
	n := minN + rand.Intn(maxN-minN+1)
	f := 1
	for i := 2; i <= n; i++ {
		f *= i
	}
	ans := strconv.Itoa(f)

	q := fmt.Sprintf("Evaluate the Laplace transform \\(L\\{t^{%d}\\}(s)\\) at \\(s = 1\\).", n)
	exp := fmt.Sprintf("Subgoal 1: State governing formula. The Laplace transform of the power function t^n for a non-negative integer n is L{t^n}(s) = n! / s^{n+1} for s > 0.\n\nSubgoal 2: Substitute parameters. Substitute power n = %d and evaluation point s = 1 into the transform: L{t^%d}(1) = %d! / 1^{%d}.\n\nSubgoal 3: Evaluate and simplify. Compute the factorial %d! = %d and divide by 1: %d / 1 = %d.",
		n, n, n, n+1, n, f, f, f)

	return generator.Problem{
		Question:    q,
		Answer:      ans,
		Explanation: exp,
		Facts: map[string]string{
			"n":      strconv.Itoa(n),
			"s":      "1",
			"answer": ans,
		},
	}
}

// pdeDiscriminantGen computes B^2 - 4AC for a second-order linear PDE.
type pdeDiscriminantGen struct{}

func (g *pdeDiscriminantGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	maxA := mathutil.ScaleInt(ctx.Difficulty, 3, 10)
	a := 1 + rand.Intn(maxA)
	maxB := mathutil.ScaleInt(ctx.Difficulty, 4, 12)
	b := rand.Intn(2*maxB+1) - maxB
	maxC := mathutil.ScaleInt(ctx.Difficulty, 3, 10)
	c := 1 + rand.Intn(maxC)
	d := b*b - 4*a*c
	ans := strconv.Itoa(d)

	q := fmt.Sprintf("For the equation \\(%du_{xx} %+du_{xy} %+du_{yy} = 0\\), find the discriminant \\(B^2 - 4AC\\).", a, b, c)
	exp := fmt.Sprintf("Subgoal 1: State governing formula. For a second-order linear PDE in two variables A u_{xx} + B u_{xy} + C u_{yy} = 0, the discriminant is defined as Δ = B^2 - 4AC.\n\nSubgoal 2: Substitute parameters. Identify coefficients A = %d, B = %d and C = %d from the equation.\n\nSubgoal 3: Evaluate and simplify. Compute Δ = (%d)^2 - 4(%d)(%d) = %d - %d = %d.",
		a, b, c, b, a, c, b*b, 4*a*c, d)

	return generator.Problem{
		Question:    q,
		Answer:      ans,
		Explanation: exp,
		Facts: map[string]string{
			"a":      strconv.Itoa(a),
			"b":      strconv.Itoa(b),
			"c":      strconv.Itoa(c),
			"answer": ans,
		},
	}
}

func smallInt(lo, hi int) int {
	if hi < lo {
		hi = lo
	}
	return lo + rand.Intn(hi-lo+1)
}

func powInt(b, e int) int {
	r := 1
	for i := 0; i < e; i++ {
		r *= b
	}
	return r
}
