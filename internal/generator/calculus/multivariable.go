package calculus

// The multivariable and vector-calculus block.
//
// `calculus.advanced` already holds the multivariable *derivatives* -- partials, the gradient, the
// Jacobian, Clairaut -- so what the domain was missing is what you integrate over a region or along
// a curve: multiple integrals, the arc length of a parametrised curve, curvature, line integrals,
// and the Jacobian evaluated as a change of measure. Those are the concepts a third calculus course
// is built from, and their absence is why a hundred of the domain's concepts have no course path.
//
// Every concept here is `numeric` and every answer is an exact integer by construction: the
// parameters are chosen so that the curvature, the arc length, the Jacobian determinant and the
// multiple integrals come out whole. That is also what keeps the easy pool discriminating, since
// each parameter still varies.
import (
	"fmt"
	"math/rand"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerMultivariable(reg *generator.Registry) {
	reg.Register("calc.vec.parametric_derivative", &parametricDerivativeGen{})
	reg.Register("calc.vec.arc_length", &arcLengthParametricGen{})
	reg.Register("calc.vec.curvature", &curvatureGen{})
	reg.Register("calc.vec.line_integral", &lineIntegralGen{})
	reg.Register("calc.vec.jacobian_value", &jacobianValueGen{})
	reg.Register("calc.multi.double_integral_rect", &doubleIntegralRectGen{})
	reg.Register("calc.multi.triple_integral_box", &tripleIntegralBoxGen{})
}

// pickMV returns a random int in [0, n), and never panics for n < 1, so a difficulty of 0 yields a
// one-element range rather than a panic inside rand.Intn.
func pickMV(n int) int {
	if n < 1 {
		n = 1
	}
	return rand.Intn(n)
}

// scaleMV widens an integer range with difficulty: `floor` is the easy-pool width and the range
// grows by up to `extra` more at difficulty 1.
func scaleMV(difficulty float64, floor, extra int) int {
	return floor + int(difficulty*float64(extra))
}

// parametricDerivativeGen asks for dy/dx on x = t^2, y = t^3 at an even t, where it is 3t/2 = 3k.

// arcLengthParametricGen uses a Pythagorean direction, so the length of the segment comes out whole.
type arcLengthParametricGen struct{}

func (g *arcLengthParametricGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// Three directions at every difficulty: the practice feed asks for three questions, so a
	// bank of two cannot fill one request.
	easy := [][3]int{{3, 4, 5}, {6, 8, 10}, {9, 12, 15}}
	hard := [][3]int{{5, 12, 13}, {8, 15, 17}, {7, 24, 25}}
	triples := hard
	if ctx.Difficulty < 0.45 {
		triples = easy
	}
	tr := triples[pickMV(len(triples))]
	return generator.Problem{
		Question: fmt.Sprintf("Find the arc length of the curve \\(x = %dt,\\ y = %dt\\) for \\(0 \\le t \\le 1\\).", tr[0], tr[1]),
		Answer:   fmt.Sprintf("%d", tr[2]),
		Explanation: fmt.Sprintf("Arc length is the integral of sqrt((dx/dt)^2 + (dy/dt)^2), and here the integrand is the constant sqrt(%d^2 + %d^2) = %d, so over an interval of length 1 the length is %d.",
			tr[0], tr[1], tr[2], tr[2]),
		Facts: map[string]string{"dx": fmt.Sprintf("%d", tr[0]), "dy": fmt.Sprintf("%d", tr[1]),
			"length": fmt.Sprintf("%d", tr[2])},
	}
}

// curvatureGen asks for the curvature of y = a x^2 at the origin, where y' = 0 and so the curvature
// is just |y”| = 2a.

// lineIntegralGen integrates the field (x, y) -- the gradient of (x^2 + y^2)/2 -- along a straight
// segment from the origin. Both offsets are even, so the potential difference is whole.
type lineIntegralGen struct{}

func (g *lineIntegralGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := 2 * (1 + pickMV(scaleMV(ctx.Difficulty, 2, 4)))
	b := 2 * (1 + pickMV(scaleMV(ctx.Difficulty, 2, 4)))
	v := (a*a + b*b) / 2
	return generator.Problem{
		Question: fmt.Sprintf("Evaluate the line integral \\(\\int x\\,dx + y\\,dy\\) along the straight segment from \\((0,0)\\) to \\((%d,%d)\\).", a, b),
		Answer:   fmt.Sprintf("%d", v),
		Explanation: fmt.Sprintf("The field is the gradient of (x^2 + y^2)/2, so the integral is the potential difference between the endpoints: (%d^2 + %d^2)/2 = %d.",
			a, b, v),
		Facts: map[string]string{"x1": fmt.Sprintf("%d", a), "y1": fmt.Sprintf("%d", b),
			"value": fmt.Sprintf("%d", v)},
	}
}

// jacobianValueGen evaluates the polar Jacobian determinant, which is r, at a sample radius.
type jacobianValueGen struct{}

func (g *jacobianValueGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	r := 1 + pickMV(scaleMV(ctx.Difficulty, 3, 5))
	return generator.Problem{
		Question: fmt.Sprintf("For the change of variables \\(x = r\\cos\\theta,\\ y = r\\sin\\theta\\), evaluate the Jacobian determinant \\(\\frac{\\partial(x,y)}{\\partial(r,\\theta)}\\) at \\(r = %d\\).", r),
		Answer:   fmt.Sprintf("%d", r),
		Explanation: fmt.Sprintf("The matrix of partials has determinant (cos t)(r cos t) - (-r sin t)(sin t) = r(sin^2 t + cos^2 t) = r, so at r = %d it is %d -- which is also why an area element in polar coordinates carries a factor of r.",
			r, r),
		Facts: map[string]string{"r": fmt.Sprintf("%d", r), "jacobian": fmt.Sprintf("%d", r)},
	}
}

// doubleIntegralRectGen integrates x + y over a rectangle; a is even, so ab(a+b)/2 is whole.
type doubleIntegralRectGen struct{}

func (g *doubleIntegralRectGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := 2 * (1 + pickMV(scaleMV(ctx.Difficulty, 2, 3)))
	b := 1 + pickMV(scaleMV(ctx.Difficulty, 2, 4))
	v := a * b * (a + b) / 2
	return generator.Problem{
		Question: fmt.Sprintf("Evaluate \\(\\int_0^{%d}\\!\\int_0^{%d} (x + y)\\,dy\\,dx\\).", a, b),
		Answer:   fmt.Sprintf("%d", v),
		Explanation: fmt.Sprintf("Integrate in y first: [xy + y^2/2] from 0 to %d is %dx + %d^2/2. Integrating that in x from 0 to %d gives %d.",
			b, b, b, a, v),
		Facts: map[string]string{"a": fmt.Sprintf("%d", a), "b": fmt.Sprintf("%d", b),
			"value": fmt.Sprintf("%d", v)},
	}
}

// tripleIntegralBoxGen integrates x over a box; a is even, so a^2 b c / 2 is whole.
type tripleIntegralBoxGen struct{}

func (g *tripleIntegralBoxGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := 2 * (1 + pickMV(scaleMV(ctx.Difficulty, 2, 3)))
	b := 1 + pickMV(scaleMV(ctx.Difficulty, 2, 3))
	c := 1 + pickMV(scaleMV(ctx.Difficulty, 2, 3))
	v := a * a * b * c / 2
	return generator.Problem{
		Question: fmt.Sprintf("Evaluate \\(\\int_0^{%d}\\!\\int_0^{%d}\\!\\int_0^{%d} x\\,dz\\,dy\\,dx\\).", a, b, c),
		Answer:   fmt.Sprintf("%d", v),
		Explanation: fmt.Sprintf("The integrand depends only on x, so the two inner integrals contribute their lengths %d and %d, leaving %d times the integral of x from 0 to %d, which is %d * %d = %d.",
			c, b, c*b, a, c*b, a*a/2, v),
		Facts: map[string]string{"a": fmt.Sprintf("%d", a), "b": fmt.Sprintf("%d", b),
			"c": fmt.Sprintf("%d", c), "value": fmt.Sprintf("%d", v)},
	}
}
