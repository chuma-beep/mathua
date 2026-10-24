package calculus

// The calculus concepts still missing after the multivariable and vector blocks: the average value
// of a function, volumes from known cross-sections, changing the order of integration, the tangent
// plane as a linear approximation, extrema on a closed region, and the two polar-coordinate
// computations.
//
// Every answer is an integer by construction: the average value uses a length that cancels, the
// cross-section volume uses a side whose cube is divisible by three, the double integrals use even
// widths, and the polar slope is evaluated at angles whose cotangent is whole.
import (
	"fmt"
	"math/rand"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerRemaining(reg *generator.Registry) {
	reg.Register("calc.integral.average_value", &averageValueGen{})
	reg.Register("calc.integral.cross_section_volume", &crossSectionGen{})
	reg.Register("calc.multi.change_order", &changeOrderGen{})
	reg.Register("calc.multi.tangent_plane", &tangentPlaneGen{})
	reg.Register("calc.multi.closed_extrema", &closedExtremaGen{})
	reg.Register("calc.polar.area", &polarAreaGen{})
	reg.Register("calc.polar.slope", &polarSlopeGen{})
}

// averageValueGen divides the integral by the interval's length, with an even coefficient so the
// result is whole.
type averageValueGen struct{}

func (g *averageValueGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	b := 2 * (1 + rand.Intn(3+int(ctx.Difficulty*4))) // the interval [0, b]
	slope := 2 * (1 + rand.Intn(3+int(ctx.Difficulty*4)))
	intercept := 1 + rand.Intn(5+int(ctx.Difficulty*6))
	avg := slope*b/2 + intercept
	return generator.Problem{
		Question: fmt.Sprintf("Find the average value of \\(f(x) = %dx + %d\\) on the interval \\([0, %d]\\).", slope, intercept, b),
		Answer:   fmt.Sprintf("%d", avg),
		Explanation: fmt.Sprintf("The average value is the integral divided by the length of the interval: (1/%d) times the area under the line, which is %d. That gives %d.",
			b, b*(slope*b/2+intercept), avg),
		Facts: map[string]string{"average": fmt.Sprintf("%d", avg), "width": fmt.Sprintf("%d", b)},
	}
}

// crossSectionGen integrates the area of a known cross-section, with a side that makes the cube
// divisible by three.
type crossSectionGen struct{}

func (g *crossSectionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// n is a multiple of six so every cross-section below gives an integer
	// volume: n^3/3, 4n^3/3, n^3/6 and n^3/12 all divide. One square section
	// left the concept with a single question however the side moved.
	n := 6 * (1 + rand.Intn(1+int(ctx.Difficulty*2)))
	shape := []string{"square", "double square", "triangle", "half square"}[rand.Intn(4)]
	var v int
	var section, areaText string
	switch shape {
	case "square":
		v = n * n * n / 3
		section = "a square with side \\(s(x) = x\\)"
		areaText = "s(x)^2 = x^2"
	case "double square":
		v = 4 * n * n * n / 3
		section = "a square with side \\(s(x) = 2x\\)"
		areaText = "s(x)^2 = 4x^2"
	case "triangle":
		v = n * n * n / 6
		section = "an isosceles right triangle with legs \\(s(x) = x\\)"
		areaText = "s(x)^2/2 = x^2/2"
	default:
		v = n * n * n / 12
		section = "a square with side \\(s(x) = x/2\\)"
		areaText = "s(x)^2 = x^2/4"
	}
	return generator.Problem{
		Question: fmt.Sprintf("A solid has, for each \\(x\\) between 0 and %d, %s perpendicular to the x-axis. What is the volume of the solid?", n, section),
		Answer:   fmt.Sprintf("%d", v),
		Explanation: fmt.Sprintf("Each cross-section has area %s, so the volume is the integral of that area along the axis: integrating over \\(x\\) from 0 to %d gives %d.",
			areaText, n, v),
		Facts: map[string]string{"volume": fmt.Sprintf("%d", v), "shape": shape, "n": fmt.Sprintf("%d", n)},
	}
}

// changeOrderGen integrates over a triangle, where swapping the order is the point of the exercise.
type changeOrderGen struct{}

func (g *changeOrderGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	n := 2 * (1 + rand.Intn(3+int(ctx.Difficulty*4)))
	v := n * n * n / 2
	return generator.Problem{
		Question: fmt.Sprintf("Evaluate \\(\\iint_R (x + y)\\,dA\\) where \\(R\\) is the triangle \\(0 \\le y \\le x \\le %d\\).", n),
		Answer:   fmt.Sprintf("%d", v),
		Explanation: fmt.Sprintf("Integrating in y first runs from 0 to x, giving %dx^2/2; integrating that in x from 0 to %d gives %d^3/2 = %d. The other order gives the same value, which is Fubini's theorem, but needs the region written as type II limits.",
			n, n, n, v),
		Facts: map[string]string{"value": fmt.Sprintf("%d", v)},
	}
}

// tangentPlaneGen evaluates the linear approximation, which is what a tangent plane is for.
type tangentPlaneGen struct{}

func (g *tangentPlaneGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := 1 + rand.Intn(3+int(ctx.Difficulty*3))
	b := 1 + rand.Intn(3+int(ctx.Difficulty*3))
	p := 1 + rand.Intn(3+int(ctx.Difficulty*3))
	q := 1 + rand.Intn(3+int(ctx.Difficulty*3))
	// f(x,y) = a x^2 + b y^2, its tangent plane at (p,q), evaluated at (p+1, q+1)
	f := a*p*p + b*q*q
	v := f + 2*a*p + 2*b*q
	return generator.Problem{
		Question: fmt.Sprintf("For \\(f(x,y) = %dx^2 + %dy^2\\), use the tangent plane at \\((%d,%d)\\) to estimate \\(f(%d,%d)\\).", a, b, p, q, p+1, q+1),
		Answer:   fmt.Sprintf("%d", v),
		Explanation: fmt.Sprintf("At (%d,%d) the value is %d, and the partial derivatives are %d and %d, so a step of one in each direction adds those: %d + %d + %d = %d. The estimate is exact here because the function is quadratic and the step is one unit.",
			p, q, f, 2*a*p, 2*b*q, f, 2*a*p, 2*b*q, v),
		Facts: map[string]string{"estimate": fmt.Sprintf("%d", v), "value": fmt.Sprintf("%d", f)},
	}
}

// closedExtremaGen tests finding extrema on a closed region by checking the boundary.
type closedExtremaGen struct{}

func (g *closedExtremaGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	c := 25 + rand.Intn(20+int(ctx.Difficulty*20))
	a := 1 + rand.Intn(3+int(ctx.Difficulty*2)) // 1 to 4
	b := 1 + rand.Intn(3+int(ctx.Difficulty*2)) // 1 to 4
	maxCoeff := a
	if b > maxCoeff {
		maxCoeff = b
	}

	if rand.Float64() < 0.5 {
		// Minimum of a downward paraboloid on x^2 + y^2 <= 4
		minVal := c - 4*maxCoeff
		return generator.Problem{
			Question: fmt.Sprintf("Find the minimum value of \\(f(x,y) = %d - %dx^2 - %dy^2\\) on the closed disc \\(x^2 + y^2 \\le 4\\).", c, a, b),
			Answer:   fmt.Sprintf("%d", minVal),
			Explanation: fmt.Sprintf("The interior critical point (0,0) is a local maximum with value %d. By the Extreme Value Theorem, the minimum occurs on the boundary circle \\(x^2 + y^2 = 4\\). Parameterizing the boundary gives \\(f(\\theta) = %d - 4(%d\\cos^2\\theta + %d\\sin^2\\theta)\\). The minimum occurs where the larger coefficient %d is multiplied by 4, yielding %d - 4(%d) = %d.",
				c, c, a, b, maxCoeff, c, maxCoeff, minVal),
			Facts: map[string]string{"minimum": fmt.Sprintf("%d", minVal), "constant": fmt.Sprintf("%d", c)},
		}
	}

	// Maximum of an upward paraboloid on x^2 + y^2 <= 4
	maxVal := c + 4*maxCoeff
	return generator.Problem{
		Question: fmt.Sprintf("Find the maximum value of \\(f(x,y) = %dx^2 + %dy^2 + %d\\) on the closed disc \\(x^2 + y^2 \\le 4\\).", a, b, c),
		Answer:   fmt.Sprintf("%d", maxVal),
		Explanation: fmt.Sprintf("The interior critical point (0,0) is a local minimum with value %d. By the Extreme Value Theorem, the maximum occurs on the boundary circle \\(x^2 + y^2 = 4\\). On the boundary, the maximum is reached when the larger coefficient %d takes the full radius squared of 4, giving 4(%d) + %d = %d.",
			c, maxCoeff, maxCoeff, c, maxVal),
		Facts: map[string]string{"maximum": fmt.Sprintf("%d", maxVal), "constant": fmt.Sprintf("%d", c)},
	}
}

// polarAreaGen uses the sector area integral, with an angular width that keeps it whole.
type polarAreaGen struct{}

func (g *polarAreaGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	r := 2 + rand.Intn(4+int(ctx.Difficulty*6))
	k := 2 * (1 + rand.Intn(3+int(ctx.Difficulty*3))) // the angular width, in radians
	v := r * r * k / 2
	return generator.Problem{
		Question: fmt.Sprintf("Find the area of the region swept by \\(r = %d\\) as \\(\\theta\\) runs from 0 to %d radians.", r, k),
		Answer:   fmt.Sprintf("%d", v),
		Explanation: fmt.Sprintf("The area in polar coordinates is (1/2) times the integral of r^2 dtheta, and here r is constant: (1/2)(%d^2)(%d) = %d. Half the square of the radius times the angle is also the fraction of the circle's area that the angle covers.",
			r, k, v),
		Facts: map[string]string{"area": fmt.Sprintf("%d", v)},
	}
}

// polarSlopeGen evaluates dy/dx for a circle given in polar form, at angles whose cotangent is
// whole.
type polarSlopeGen struct{}

func (g *polarSlopeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// r = 5 is a circle: dy/dx = -cot(theta). Choosing theta with tan(theta) = b/a
	// from a Pythagorean triple keeps the slope an exact decimal while varying the
	// angle.
	pairs := [][2]int{{1, 1}, {3, 4}, {4, 3}, {5, 12}, {12, 5}, {8, 15}, {15, 8}, {7, 24}, {24, 7}}
	p := pairs[rand.Intn(len(pairs))]
	a, b := p[0], p[1]
	if rand.Intn(2) == 0 {
		a, b = -a, -b
	}
	value := -float64(a) / float64(b)
	answer := decimal(value, 4)
	return generator.Problem{
		Question: fmt.Sprintf("The curve \\(r = 5\\) is a circle in polar coordinates. Find \\(\\frac{dy}{dx}\\) at \\(\\theta\\) where \\(\\tan\\theta = %d/%d\\). (enter a decimal)", b, a),
		Answer:   answer,
		Explanation: fmt.Sprintf("For a polar curve, dy/dx is (r' sin t + r cos t)/(r' cos t - r sin t). With r constant its derivative is zero, leaving -cot t = -%d/%d = %s. That agrees with the tangent to the circle at that point.",
			a, b, answer),
		Facts: map[string]string{"slope": answer},
	}
}
