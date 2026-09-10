package geometry

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/mathutil"
)

func Register(reg *generator.Registry) {
	reg.Register("geo.basic.points_lines", &pointsLinesGen{})
	reg.Register("geo.basic.angles", &angleTypesGen{})
	reg.Register("geo.basic.angle_measure", &angleMeasureGen{})
	reg.Register("geo.basic.complementary", &complementaryGen{})
	reg.Register("geo.basic.vertical", &verticalAnglesGen{})
	reg.Register("geo.triangle.types", &triangleTypesGen{})
	reg.Register("geo.triangle.angles", &triangleAnglesGen{})
	reg.Register("geo.triangle.area", &triangleAreaGen{})
	reg.Register("geo.triangle.pythagorean", &pythagoreanGen{})
	reg.Register("geo.quad.types", &quadTypesGen{})
	reg.Register("geo.quad.area", &quadAreaGen{})
	reg.Register("geo.quad.perimeter", &quadPerimGen{})
	reg.Register("geo.circle.parts", &circlePartsGen{})
	reg.Register("geo.circle.circumference", &circumferenceGen{})
	reg.Register("geo.circle.area", &circleAreaGen{})
	reg.Register("geo.coord.plot", &coordPlotGen{})
	reg.Register("geo.coord.distance", &coordDistanceGen{})
	reg.Register("geo.coord.midpoint", &coordMidpointGen{})
	reg.Register("geo.solid.volume_rect", &volumeGen{})
	reg.Register("geo.solid.surface_area", &surfaceAreaGen{})

	reg.Register("geo.coord.lines", &coordLinesGen{})
	reg.Register("geo.coord.polar", &coordPolarGen{})
}

type pointsLinesGen struct{}

func (g *pointsLinesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	items := []string{"point", "line", "ray", "line segment"}
	descs := []string{"an exact location with no size", "goes forever in two directions", "has one endpoint and goes forever", "has two endpoints"}
	i := rand.Intn(4)
	return generator.Problem{
		Question:    fmt.Sprintf("What is: %s?", descs[i]),
		Answer:      items[i],
		Explanation: fmt.Sprintf("A %s is %s.", items[i], descs[i]),
	}
}

type angleTypesGen struct{}

func (g *angleTypesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	deg := rand.Intn(max(1, scale*30))
	ans := "straight"
	if deg < 90 {
		ans = "acute"
	} else if deg == 90 {
		ans = "right"
	} else if deg < 180 {
		ans = "obtuse"
	}
	return generator.Problem{
		Question:    fmt.Sprintf("An angle of %d degrees is classified as:", deg),
		Answer:      ans,
		Explanation: fmt.Sprintf("Acute < 90, right = 90, obtuse > 90 and < 180, straight = 180. %d is %s.", deg, ans),
	}
}

type angleMeasureGen struct{}

func (g *angleMeasureGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := rand.Intn(max(1, scale*30)) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("An angle of %d degrees is a:", a),
		Answer:      classAngle(a),
		Explanation: fmt.Sprintf("%d degrees is %s.", a, classAngle(a)),
	}
}

func classAngle(d int) string {
	switch {
	case d < 90:
		return "acute angle"
	case d == 90:
		return "right angle"
	case d < 180:
		return "obtuse angle"
	default:
		return "straight angle"
	}
}

type complementaryGen struct{}

func (g *complementaryGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := rand.Intn(max(1, scale*15)) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("Two angles are complementary. One is %d degrees. What is the other?", a),
		Answer:      fmt.Sprintf("%d", 90-a),
		Explanation: fmt.Sprintf("Complementary angles sum to 90. 90 - %d = %d.", a, 90-a),
	}
}

type verticalAnglesGen struct{}

func (g *verticalAnglesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := rand.Intn(max(1, scale*30)) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("Two lines intersect. One angle is %d degrees. What is the vertical angle?", a),
		Answer:      fmt.Sprintf("%d", a),
		Explanation: fmt.Sprintf("Vertical angles are equal: %d = %d.", a, a),
	}
}

type triangleTypesGen struct{}

func (g *triangleTypesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	sides := []int{rand.Intn(max(1, scale*2)) + 2, rand.Intn(max(1, scale*2)) + 2, rand.Intn(max(1, scale*2)) + 2}
	ans := "scalene"
	if sides[0] == sides[1] && sides[1] == sides[2] {
		ans = "equilateral"
	} else if sides[0] == sides[1] || sides[1] == sides[2] || sides[0] == sides[2] {
		ans = "isosceles"
	}
	return generator.Problem{
		Question:    fmt.Sprintf("A triangle has sides %d, %d, %d. By side length, it is:", sides[0], sides[1], sides[2]),
		Answer:      ans,
		Explanation: fmt.Sprintf("All 3 equal = equilateral. Two equal = isosceles. None equal = scalene."),
	}
}

type triangleAnglesGen struct{}

func (g *triangleAnglesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := rand.Intn(max(1, scale*15)) + 30
	b := rand.Intn(max(1, scale*15)) + 20
	c := 180 - a - b
	return generator.Problem{
		Question:    fmt.Sprintf("A triangle has angles %d and %d degrees. What is the third angle?", a, b),
		Answer:      fmt.Sprintf("%d", c),
		Explanation: fmt.Sprintf("Triangle angles sum to 180. 180 - %d - %d = %d.", a, b, c),
	}
}

type triangleAreaGen struct{}

func (g *triangleAreaGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	base := (rand.Intn(max(1, scale*4)) + 1) * 2
	height := rand.Intn(max(1, scale*4)) + 1
	area := base * height / 2
	return generator.Problem{
		Question:    fmt.Sprintf("Triangle: base = %d, height = %d. Find the area.", base, height),
		Answer:      fmt.Sprintf("%d", area),
		Explanation: fmt.Sprintf("\\(\\text{Area} = \\frac{\\text{base} \\times \\text{height}}{2} = \\frac{%d \\times %d}{2} = %d\\).", base, height, area),
	}
}

type pythagoreanGen struct{}

func (g *pythagoreanGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := rand.Intn(max(1, scale*3)) + 3
	b := rand.Intn(max(1, scale*3)) + 3
	cSq := a*a + b*b
	c := int(math.Sqrt(float64(cSq)))
	for c*c != cSq {
		a = rand.Intn(max(1, scale*3)) + 3
		b = rand.Intn(max(1, scale*3)) + 3
		cSq = a*a + b*b
		c = int(math.Sqrt(float64(cSq)))
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Right triangle: legs \\(= %d\\) and \\(%d\\). Find the hypotenuse \\(c\\).", a, b),
		Answer:      fmt.Sprintf("%d", c),
		Explanation: fmt.Sprintf("\\(c^{2} = %d^{2} + %d^{2} = %d + %d = %d\\), so \\(c = \\sqrt{%d} = %d\\).", a, b, a*a, b*b, c*c, c*c, c),
	}
}

type quadTypesGen struct{}

func (g *quadTypesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	types := []struct{ name, desc string }{
		{"square", "all sides equal and all angles 90 degrees"},
		{"rectangle", "opposite sides equal and all angles 90 degrees"},
		{"rhombus", "all sides equal, angles not necessarily 90"},
		{"parallelogram", "opposite sides parallel and equal"},
		{"trapezoid", "one pair of parallel sides"},
	}
	i := rand.Intn(len(types))
	return generator.Problem{
		Question:    fmt.Sprintf("A quadrilateral with %s is a:", types[i].desc),
		Answer:      types[i].name,
		Explanation: types[i].desc + " describes a " + types[i].name + ".",
	}
}

type quadAreaGen struct{}

func (g *quadAreaGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	base := rand.Intn(max(1, scale*4)) + 5
	height := rand.Intn(max(1, scale*4)) + 5
	shape := []string{"rectangle", "parallelogram"}[rand.Intn(2)]
	return generator.Problem{
		Question:    fmt.Sprintf("Find the area of a %s with base %d and height %d.", shape, base, height),
		Answer:      fmt.Sprintf("%d", base*height),
		Explanation: fmt.Sprintf("\\(\\text{Area} = \\text{base} \\times \\text{height} = %d \\times %d = %d\\).", base, height, base*height),
	}
}

type quadPerimGen struct{}

func (g *quadPerimGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	w := rand.Intn(max(1, scale*5)) + 5
	h := rand.Intn(max(1, scale*5)) + 5
	return generator.Problem{
		Question:    fmt.Sprintf("Rectangle: width = %d, height = %d. Find the perimeter.", w, h),
		Answer:      fmt.Sprintf("%d", 2*(w+h)),
		Explanation: fmt.Sprintf("\\(\\text{Perimeter} = 2(\\text{width} + \\text{height}) = 2(%d + %d) = %d\\).", w, h, 2*(w+h)),
	}
}

type circlePartsGen struct{}

func (g *circlePartsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	parts := []string{"radius", "diameter", "circumference", "chord"}
	descs := []string{"distance from center to edge", "distance across the circle through center", "distance around the circle", "a line from one edge to another (not through center)"}
	i := rand.Intn(4)
	return generator.Problem{
		Question:    fmt.Sprintf("The %s is the:", descs[i]),
		Answer:      parts[i],
		Explanation: descs[i] + " = " + parts[i] + ".",
	}
}

type circumferenceGen struct{}

func (g *circumferenceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	r := rand.Intn(max(1, scale*4)) + 5
	pi := 3.14
	circ := math.Round(2 * pi * float64(r))
	return generator.Problem{
		Question:    fmt.Sprintf("Circle radius = %d. Find the circumference (use pi = 3.14).", r),
		Answer:      fmt.Sprintf("%.0f", circ),
		Explanation: fmt.Sprintf("\\(C = 2 \\pi r = 2 \\times 3.14 \\times %d = %.0f\\).", r, 2*3.14*float64(r)),
	}
}

type circleAreaGen struct{}

func (g *circleAreaGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	r := rand.Intn(max(1, scale*2)) + 3
	area := math.Round(3.14 * float64(r*r))
	return generator.Problem{
		Question:    fmt.Sprintf("Circle radius = %d. Find the area (use pi = 3.14).", r),
		Answer:      fmt.Sprintf("%.0f", area),
		Explanation: fmt.Sprintf("\\(A = \\pi r^{2} = 3.14 \\times %d^{2} = 3.14 \\times %d = %.0f\\).", r, r*r, 3.14*float64(r*r)),
	}
}

type coordPlotGen struct{}

func (g *coordPlotGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	x := rand.Intn(max(1, scale*2)) - 5
	y := rand.Intn(max(1, scale*2)) - 5
	return generator.Problem{
		Question:    fmt.Sprintf("What quadrant is the point \\((%d, %d)\\) in?", x, y),
		Answer:      quad(x, y),
		Explanation: fmt.Sprintf("(%d,%d) is in %s.", x, y, quad(x, y)),
	}
}

func quad(x, y int) string {
	switch {
	case x > 0 && y > 0:
		return "Quadrant I"
	case x < 0 && y > 0:
		return "Quadrant II"
	case x < 0 && y < 0:
		return "Quadrant III"
	case x > 0 && y < 0:
		return "Quadrant IV"
	default:
		return "on an axis"
	}
}

type coordDistanceGen struct{}

func (g *coordDistanceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	dx := rand.Intn(max(1, scale*4)) + 1
	dy := rand.Intn(max(1, scale*4)) + 1
	dist := int(math.Round(math.Sqrt(float64(dx*dx + dy*dy))))
	for dist*dist != dx*dx+dy*dy {
		dx = rand.Intn(max(1, scale*4)) + 1
		dy = rand.Intn(max(1, scale*4)) + 1
		dist = int(math.Round(math.Sqrt(float64(dx*dx + dy*dy))))
	}
	x1 := rand.Intn(5)
	y1 := rand.Intn(5)
	return generator.Problem{
		Question:    fmt.Sprintf("Find the distance between \\((%d,%d)\\) and \\((%d,%d)\\).", x1, y1, x1+dx, y1+dy),
		Answer:      fmt.Sprintf("%d", dist),
		Explanation: fmt.Sprintf("\\(\\sqrt{(%d-%d)^{2} + (%d-%d)^{2}} = \\sqrt{%d + %d} = \\sqrt{%d} = %d\\).", x1+dx, x1, y1+dy, y1, dx*dx, dy*dy, dx*dx+dy*dy, dist),
	}
}

type coordMidpointGen struct{}

func (g *coordMidpointGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	x1 := (rand.Intn(max(1, scale*2)) + 1) * 2
	y1 := (rand.Intn(max(1, scale*2)) + 1) * 2
	x2 := x1 + (rand.Intn(5)+1)*2
	y2 := y1 + (rand.Intn(5)+1)*2
	return generator.Problem{
		Question:    fmt.Sprintf("Find the midpoint of \\((%d,%d)\\) and \\((%d,%d)\\).", x1, y1, x2, y2),
		Answer:      fmt.Sprintf("(%.1f, %.1f)", float64(x1+x2)/2, float64(y1+y2)/2),
		Explanation: fmt.Sprintf("\\(\\text{Midpoint} = \\left(\\frac{%d+%d}{2}, \\frac{%d+%d}{2}\\right) = (%.1f, %.1f)\\).", x1, x2, y1, y2, float64(x1+x2)/2, float64(y1+y2)/2),
	}
}

type volumeGen struct{}

func (g *volumeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	l := rand.Intn(max(1, scale*2)) + 2
	w := rand.Intn(max(1, scale*2)) + 2
	h := rand.Intn(max(1, scale*2)) + 2
	return generator.Problem{
		Question:    fmt.Sprintf("Rectangular prism: length=%d width=%d height=%d. Find volume.", l, w, h),
		Answer:      fmt.Sprintf("%d", l*w*h),
		Explanation: fmt.Sprintf("\\(V = l \\times w \\times h = %d \\times %d \\times %d = %d\\).", l, w, h, l*w*h),
	}
}

type surfaceAreaGen struct{}

func (g *surfaceAreaGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	l := rand.Intn(max(1, scale*2)) + 2
	w := rand.Intn(max(1, scale*2)) + 2
	h := rand.Intn(max(1, scale*2)) + 2
	sa := 2 * (l*w + l*h + w*h)
	return generator.Problem{
		Question:    fmt.Sprintf("Rectangular prism: length=%d width=%d height=%d. Find surface area.", l, w, h),
		Answer:      fmt.Sprintf("%d", sa),
		Explanation: fmt.Sprintf("\\(\\text{SA} = 2(lw + lh + wh) = 2(%d + %d + %d) = %d\\).", l*w, l*h, w*h, sa),
	}
}

type coordLinesGen struct{}

func (g *coordLinesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	x1 := rand.Intn(max(1, scale)) - 2
	y1 := rand.Intn(max(1, scale)) - 2
	x2 := x1 + rand.Intn(max(1, scale)) + 1
	y2 := y1 + rand.Intn(max(1, scale)) + 1
	dx := x2 - x1
	dy := y2 - y1
	// slope = dy/dx
	gcd := mathutil.GCD(dy, dx)
	mNum := dy / gcd
	mDen := dx / gcd
	if mDen < 0 {
		mNum = -mNum
		mDen = -mDen
	}
	// y - y1 = m(x - x1), solve for y-intercept
	// y = mx + b, b = y1 - m*x1
	// b = y1 - (dy/dx)*x1 = (y1*dx - dy*x1)/dx
	bNum := y1*dx - dy*x1
	bDen := dx
	bGcd := mathutil.GCD(mathutil.Abs(bNum), mathutil.Abs(bDen))
	bNum /= bGcd
	bDen /= bGcd
	if bDen < 0 {
		bNum = -bNum
		bDen = -bDen
	}

	var answer string
	if mDen == 1 {
		if bDen == 1 {
			answer = fmt.Sprintf("y=%dx+%d", mNum, bNum)
		} else {
			answer = fmt.Sprintf("y=%dx+%d/%d", mNum, bNum, bDen)
		}
	} else {
		if bDen == 1 {
			answer = fmt.Sprintf("y=%d/%dx+%d", mNum, mDen, bNum)
		} else {
			answer = fmt.Sprintf("y=%d/%dx+%d/%d", mNum, mDen, bNum, bDen)
		}
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Find the equation of the line through \\((%d,%d)\\) and \\((%d,%d)\\).", x1, y1, x2, y2),
		Answer:      answer,
		Explanation: fmt.Sprintf("\\(\\text{Slope} = \\frac{%d-%d}{%d-%d} = \\frac{%d}{%d}\\). Line through \\((%d,%d)\\): \\(y - %d = \\frac{%d}{%d}(x - %d) \\to %s\\)", y2, y1, x2, x1, dy, dx, x1, y1, y1, dy, dx, x1, answer),
	}
}

type coordPolarGen struct{}

func (g *coordPolarGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		r, theta int
		x, y     int
	}
	entries := []entry{
		{1, 0, 1, 0},
		{1, 90, 0, 1},
		{1, 180, -1, 0},
		{1, 270, 0, -1},
		{2, 0, 2, 0},
		{3, 90, 0, 3},
		{2, 180, -2, 0},
		{4, 270, 0, -4},
		{2, 45, 1, 1},    // approx (√2, √2) → 1,1
		{3, 45, 2, 2},    // approx
		{2, 135, -1, 1},  // approx
		{2, 225, -1, -1}, // approx
		{2, 315, 1, -1},  // approx
	}
	e := entries[rand.Intn(len(entries))]
	if rand.Intn(2) == 0 {
		return generator.Problem{
			Question:    fmt.Sprintf("Convert \\((r=%d, \\theta=%d^{\\circ})\\) to rectangular coordinates.", e.r, e.theta),
			Answer:      fmt.Sprintf("(%d,%d)", e.x, e.y),
			Explanation: fmt.Sprintf("\\(x = %d \\cdot \\cos(%d^{\\circ}) = %d\\), \\(y = %d \\cdot \\sin(%d^{\\circ}) = %d \\to (%d,%d)\\)", e.r, e.theta, e.x, e.r, e.theta, e.y, e.x, e.y),
		}
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Convert \\((%d,%d)\\) to polar coordinates \\((r > 0, 0 \\leq \\theta < 360)\\).", e.x, e.y),
		Answer:      fmt.Sprintf("(%d,%d°)", e.r, e.theta),
		Explanation: fmt.Sprintf("\\(r = \\sqrt{%d^{2}+%d^{2}} = %d\\), \\(\\theta = \\arctan\\left(\\frac{%d}{%d}\\right) = %d^{\\circ} \\to (%d,%d^{\\circ})\\)", e.x, e.y, e.r, e.y, e.x, e.theta, e.r, e.theta),
	}
}
