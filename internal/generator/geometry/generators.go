package geometry

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/chuma-beep/mathua/internal/generator"
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

	// stub generators for newly-added concepts
	reg.Register("geo.coord.lines", &generator.Stub{ConceptID: "geo.coord.lines"})
	reg.Register("geo.coord.polar", &generator.Stub{ConceptID: "geo.coord.polar"})
}

type pointsLinesGen struct{}

func (g *pointsLinesGen) Generate(difficulty float64) generator.Problem {
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

func (g *angleTypesGen) Generate(difficulty float64) generator.Problem {
	deg := rand.Intn(181)
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

func (g *angleMeasureGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(179) + 1
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

func (g *complementaryGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(89) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("Two angles are complementary. One is %d degrees. What is the other?", a),
		Answer:      fmt.Sprintf("%d", 90-a),
		Explanation: fmt.Sprintf("Complementary angles sum to 90. 90 - %d = %d.", a, 90-a),
	}
}

type verticalAnglesGen struct{}

func (g *verticalAnglesGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(179) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("Two lines intersect. One angle is %d degrees. What is the vertical angle?", a),
		Answer:      fmt.Sprintf("%d", a),
		Explanation: fmt.Sprintf("Vertical angles are equal: %d = %d.", a, a),
	}
}

type triangleTypesGen struct{}

func (g *triangleTypesGen) Generate(difficulty float64) generator.Problem {
	sides := []int{rand.Intn(8) + 2, rand.Intn(8) + 2, rand.Intn(8) + 2}
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

func (g *triangleAnglesGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(90) + 30
	b := rand.Intn(180-a-20) + 20
	c := 180 - a - b
	return generator.Problem{
		Question:    fmt.Sprintf("A triangle has angles %d and %d degrees. What is the third angle?", a, b),
		Answer:      fmt.Sprintf("%d", c),
		Explanation: fmt.Sprintf("Triangle angles sum to 180. 180 - %d - %d = %d.", a, b, c),
	}
}

type triangleAreaGen struct{}

func (g *triangleAreaGen) Generate(difficulty float64) generator.Problem {
	base := (rand.Intn(20) + 1) * 2
	height := rand.Intn(20) + 1
	area := base * height / 2
	return generator.Problem{
		Question:    fmt.Sprintf("Triangle: base = %d, height = %d. Find the area.", base, height),
		Answer:      fmt.Sprintf("%d", area),
		Explanation: fmt.Sprintf("Area = base x height / 2 = %d x %d / 2 = %d.", base, height, area),
	}
}

type pythagoreanGen struct{}

func (g *pythagoreanGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(12) + 3
	b := rand.Intn(12) + 3
	cSq := a*a + b*b
	c := int(math.Sqrt(float64(cSq)))
	// Ensure perfect square
	for c*c != cSq {
		a = rand.Intn(12) + 3
		b = rand.Intn(12) + 3
		cSq = a*a + b*b
		c = int(math.Sqrt(float64(cSq)))
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Right triangle: legs = %d and %d. Find the hypotenuse c.", a, b),
		Answer:      fmt.Sprintf("%d", c),
		Explanation: fmt.Sprintf("c^2 = %d^2 + %d^2 = %d + %d = %d, so c = sqrt(%d) = %d.", a, b, a*a, b*b, c*c, c*c, c),
	}
}

type quadTypesGen struct{}

func (g *quadTypesGen) Generate(difficulty float64) generator.Problem {
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

func (g *quadAreaGen) Generate(difficulty float64) generator.Problem {
	base := rand.Intn(20) + 5
	height := rand.Intn(20) + 5
	shape := []string{"rectangle", "parallelogram"}[rand.Intn(2)]
	return generator.Problem{
		Question:    fmt.Sprintf("Find the area of a %s with base %d and height %d.", shape, base, height),
		Answer:      fmt.Sprintf("%d", base*height),
		Explanation: fmt.Sprintf("Area = base x height = %d x %d = %d.", base, height, base*height),
	}
}

type quadPerimGen struct{}

func (g *quadPerimGen) Generate(difficulty float64) generator.Problem {
	w := rand.Intn(30) + 5
	h := rand.Intn(30) + 5
	return generator.Problem{
		Question:    fmt.Sprintf("Rectangle: width = %d, height = %d. Find the perimeter.", w, h),
		Answer:      fmt.Sprintf("%d", 2*(w+h)),
		Explanation: fmt.Sprintf("Perimeter = 2 x (width + height) = 2 x (%d + %d) = %d.", w, h, 2*(w+h)),
	}
}

type circlePartsGen struct{}

func (g *circlePartsGen) Generate(difficulty float64) generator.Problem {
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

func (g *circumferenceGen) Generate(difficulty float64) generator.Problem {
	r := rand.Intn(20) + 5
	pi := 3.14
	circ := math.Round(2 * pi * float64(r))
	return generator.Problem{
		Question:    fmt.Sprintf("Circle radius = %d. Find the circumference (use pi = 3.14).", r),
		Answer:      fmt.Sprintf("%.0f", circ),
		Explanation: fmt.Sprintf("C = 2 x pi x r = 2 x 3.14 x %d = %.0f.", r, 2*3.14*float64(r)),
	}
}

type circleAreaGen struct{}

func (g *circleAreaGen) Generate(difficulty float64) generator.Problem {
	r := rand.Intn(10) + 3
	area := math.Round(3.14 * float64(r*r))
	return generator.Problem{
		Question:    fmt.Sprintf("Circle radius = %d. Find the area (use pi = 3.14).", r),
		Answer:      fmt.Sprintf("%.0f", area),
		Explanation: fmt.Sprintf("A = pi x r^2 = 3.14 x %d^2 = 3.14 x %d = %.0f.", r, r*r, 3.14*float64(r*r)),
	}
}

type coordPlotGen struct{}

func (g *coordPlotGen) Generate(difficulty float64) generator.Problem {
	x := rand.Intn(11) - 5
	y := rand.Intn(11) - 5
	return generator.Problem{
		Question:    fmt.Sprintf("What quadrant is the point (%d, %d) in?", x, y),
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

func (g *coordDistanceGen) Generate(difficulty float64) generator.Problem {
	dx := rand.Intn(9) + 1
	dy := rand.Intn(9) + 1
	dist := int(math.Round(math.Sqrt(float64(dx*dx + dy*dy))))
	for dist*dist != dx*dx+dy*dy {
		dx = rand.Intn(8) + 1
		dy = rand.Intn(8) + 1
		dist = int(math.Round(math.Sqrt(float64(dx*dx + dy*dy))))
	}
	x1 := rand.Intn(5)
	y1 := rand.Intn(5)
	return generator.Problem{
		Question:    fmt.Sprintf("Find the distance between (%d,%d) and (%d,%d).", x1, y1, x1+dx, y1+dy),
		Answer:      fmt.Sprintf("%d", dist),
		Explanation: fmt.Sprintf("sqrt((%d-%d)^2 + (%d-%d)^2) = sqrt(%d + %d) = sqrt(%d) = %d.", x1+dx, x1, y1+dy, y1, dx*dx, dy*dy, dx*dx+dy*dy, dist),
	}
}

type coordMidpointGen struct{}

func (g *coordMidpointGen) Generate(difficulty float64) generator.Problem {
	x1 := (rand.Intn(9) + 1) * 2
	y1 := (rand.Intn(9) + 1) * 2
	x2 := x1 + (rand.Intn(5)+1)*2
	y2 := y1 + (rand.Intn(5)+1)*2
	return generator.Problem{
		Question:    fmt.Sprintf("Find the midpoint of (%d,%d) and (%d,%d).", x1, y1, x2, y2),
		Answer:      fmt.Sprintf("(%.1f, %.1f)", float64(x1+x2)/2, float64(y1+y2)/2),
		Explanation: fmt.Sprintf("Midpoint = ((%d+%d)/2, (%d+%d)/2) = (%.1f, %.1f).", x1, x2, y1, y2, float64(x1+x2)/2, float64(y1+y2)/2),
	}
}

type volumeGen struct{}

func (g *volumeGen) Generate(difficulty float64) generator.Problem {
	l := rand.Intn(10) + 2
	w := rand.Intn(10) + 2
	h := rand.Intn(10) + 2
	return generator.Problem{
		Question:    fmt.Sprintf("Rectangular prism: length=%d width=%d height=%d. Find volume.", l, w, h),
		Answer:      fmt.Sprintf("%d", l*w*h),
		Explanation: fmt.Sprintf("V = l x w x h = %d x %d x %d = %d.", l, w, h, l*w*h),
	}
}

type surfaceAreaGen struct{}

func (g *surfaceAreaGen) Generate(difficulty float64) generator.Problem {
	l := rand.Intn(10) + 2
	w := rand.Intn(10) + 2
	h := rand.Intn(10) + 2
	sa := 2 * (l*w + l*h + w*h)
	return generator.Problem{
		Question:    fmt.Sprintf("Rectangular prism: length=%d width=%d height=%d. Find surface area.", l, w, h),
		Answer:      fmt.Sprintf("%d", sa),
		Explanation: fmt.Sprintf("SA = 2(lw + lh + wh) = 2(%d + %d + %d) = %d.", l*w, l*h, w*h, sa),
	}
}
