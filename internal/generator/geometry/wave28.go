package geometry

// Wave 28 geometry: coordinates, circles, solids, polygons and transformations.
// Every family publishes the facts its solution schema references and answers
// with an integer.

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave28(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"geo.coord.vector":            "vector",
		"geo.coord.dot":               "dot",
		"geo.coord.parametric_line":   "parametric",
		"geo.coord.normal_line":       "normal",
		"geo.coord.equation_circle":   "eqcircle",
		"geo.coord.rotated_conic":     "conic",
		"geo.coord.polar_curve":       "polar",
		"geo.coord.area_polygon":      "shoelace",
		"geo.coord.line_intersection": "intersect",
		"geo.coord.three_d":           "threed",
		"geo.circle.thales":           "thales",
		"geo.circle.arc_sector":       "arcsector",
		"geo.circle.segment_area":     "segment",
		"geo.circle.common_tangents":  "tangents",
		"geo.circle.orthogonal":       "orthogonal",
		"geo.circle.inversion":        "inversion",
		"geo.circle.nine_point":       "ninepoint",
		"geo.circle.euler_line":       "eulerline",
		"geo.solid.tetrahedron":       "tetra",
		"geo.solid.octahedron":        "octa",
		"geo.solid.platonic":          "platonic",
		"geo.solid.oblique_prism":     "oblique",
		"geo.solid.composite_volume":  "composite",
		"geo.solid.cavalieri":         "cavalieri",
		"geo.solid.rotation_volume":   "rotation",
		"geo.solid.truncated_cone":    "truncated",
		"geo.polygon.apothem":         "apothem",
		"geo.polygon.star":            "star",
		"geo.polygon.constructible":   "constructible",
		"geo.polygon.quadrilateral_area": "quadarea",
		"geo.transform.homothety":     "homothety",
		"geo.transform.glide":         "glide",
		"geo.transform.group":         "group",
	} {
		reg.Register(id, &w28Gen{kind: kind})
	}
}

type w28Gen struct{ kind string }

func (g *w28Gen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "vector":
		return w28Vector(kp)
	case "dot":
		return w28Dot(kp)
	case "parametric":
		return w28Parametric(kp)
	case "normal":
		return w28Normal(kp)
	case "eqcircle":
		return w28EqCircle(kp)
	case "conic":
		return w28Conic(kp)
	case "polar":
		return w28Polar(kp)
	case "shoelace":
		return w28Shoelace(kp)
	case "intersect":
		return w28Intersect(kp)
	case "threed":
		return w28ThreeD(kp)
	case "thales":
		return w28Thales(kp)
	case "arcsector":
		return w28ArcSector(kp)
	case "segment":
		return w28Segment(kp)
	case "tangents":
		return w28Tangents(kp)
	case "orthogonal":
		return w28Orthogonal(kp)
	case "inversion":
		return w28Inversion(kp)
	case "ninepoint":
		return w28NinePoint(kp)
	case "eulerline":
		return w28EulerLine(kp)
	case "tetra":
		return w28Tetra(kp)
	case "octa":
		return w28Octa(kp)
	case "platonic":
		return w28Platonic(kp)
	case "oblique":
		return w28Oblique(kp)
	case "composite":
		return w28Composite(kp)
	case "cavalieri":
		return w28Cavalieri(kp)
	case "rotation":
		return w28Rotation(kp)
	case "truncated":
		return w28Truncated(kp)
	case "apothem":
		return w28Apothem(kp)
	case "star":
		return w28Star(kp)
	case "constructible":
		return w28Constructible(kp)
	case "quadarea":
		return w28QuadArea(kp)
	case "homothety":
		return w28Homothety(kp)
	case "glide":
		return w28Glide(kp)
	case "group":
		return w28Group(kp)
	}
	return w28P(kp, "What is 2 + 2? (enter a number)", "4", "2 + 2 = 4.", nil)
}

func w28N(lo, hi int) int { return lo + rand.Intn(hi-lo+1) }

func w28P(o int, q, a, ex string, f map[string]string) generator.Problem {
	if f == nil {
		f = map[string]string{}
	}
	f["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: f, KPOrdinal: o}
}

func w28I(f map[string]string, k string, v int) { f[k] = strconv.Itoa(v) }

func w28GCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func w28Vector(o int) generator.Problem {
	a := w28N(1, 60)
	b := w28N(1, 60)
	sq := a*a + b*b
	f := map[string]string{}
	w28I(f, "a", a)
	w28I(f, "b", b)
	w28I(f, "sq", sq)
	switch o {
	case 1:
		return w28P(o, fmt.Sprintf("A vector has components (%d,%d). What is its squared length?", a, b), strconv.Itoa(sq), "Square and add the components.", f)
	case 2:
		return w28P(o, fmt.Sprintf("A displacement has components (%d,%d). What is its x component?", a, b), strconv.Itoa(a), "Read the first component.", f)
	default:
		return w28P(o, fmt.Sprintf("A displacement has components (%d,%d). What is its y component?", a, b), strconv.Itoa(b), "Read the second component.", f)
	}
}

func w28Dot(o int) generator.Problem {
	a := w28N(1, 30)
	b := w28N(1, 30)
	dot := a * b
	f := map[string]string{}
	w28I(f, "a", a)
	w28I(f, "b", b)
	w28I(f, "dot", dot)
	switch o {
	case 1:
		return w28P(o, fmt.Sprintf("Two vectors are (%d,0) and (%d,0). What is their dot product?", a, b), strconv.Itoa(dot), "Multiply the matching components.", f)
	case 2:
		return w28P(o, fmt.Sprintf("Two vectors are (%d,0) and (%d,0). What is the x component of the first?", a, b), strconv.Itoa(a), "Read the first component.", f)
	default:
		return w28P(o, fmt.Sprintf("Two vectors are (%d,0) and (%d,0). What is the x component of the second?", a, b), strconv.Itoa(b), "Read the second vector.", f)
	}
}

func w28Parametric(o int) generator.Problem {
	x0 := w28N(1, 30)
	dx := w28N(1, 12)
	t := w28N(1, 30)
	x := x0 + dx*t
	f := map[string]string{}
	w28I(f, "x0", x0)
	w28I(f, "y0", 0)
	w28I(f, "dx", dx)
	w28I(f, "dy", 0)
	w28I(f, "t", t)
	w28I(f, "x", x)
	switch o {
	case 1:
		return w28P(o, fmt.Sprintf("A line starts at x=%d with step %d and parameter t=%d. What is x?", x0, dx, t), strconv.Itoa(x), "Add the step multiplied by t.", f)
	case 2:
		return w28P(o, fmt.Sprintf("A line starts at x=%d with step %d. What is the step?", x0, dx), strconv.Itoa(dx), "The step is the parameter coefficient.", f)
	default:
		return w28P(o, fmt.Sprintf("A parametric line uses t=%d with step %d. What is the parameter?", t, dx), strconv.Itoa(t), "The parameter is given.", f)
	}
}

func w28Normal(o int) generator.Problem {
	a := w28N(1, 20)
	b := w28N(1, 20)
	c := w28N(1, 60)
	sq := a*a + b*b
	f := map[string]string{}
	w28I(f, "a", a)
	w28I(f, "b", b)
	w28I(f, "c", c)
	w28I(f, "sq", sq)
	switch o {
	case 1:
		return w28P(o, fmt.Sprintf("A line is %dx+%dy=%d. What is the squared length of its normal vector?", a, b, c), strconv.Itoa(sq), "Square and add the coefficients.", f)
	case 2:
		return w28P(o, fmt.Sprintf("A line is %dx+%dy=%d. What is the x coefficient?", a, b, c), strconv.Itoa(a), "Read the coefficient of x.", f)
	default:
		return w28P(o, fmt.Sprintf("A line is %dx+%dy=%d. What is the constant on the right?", a, b, c), strconv.Itoa(c), "The constant is given.", f)
	}
}

func w28EqCircle(o int) generator.Problem {
	h := w28N(1, 30)
	k := w28N(1, 30)
	r := w28N(1, 30)
	r2 := r * r
	f := map[string]string{}
	w28I(f, "h", h)
	w28I(f, "k", k)
	w28I(f, "r", r)
	w28I(f, "r2", r2)
	switch o {
	case 1:
		return w28P(o, fmt.Sprintf("A circle has centre (%d,%d) and radius %d. What is r squared?", h, k, r), strconv.Itoa(r2), "Square the radius.", f)
	case 2:
		return w28P(o, fmt.Sprintf("A circle has centre (%d,%d) and radius %d. What is the radius?", h, k, r), strconv.Itoa(r), "The radius is given.", f)
	default:
		return w28P(o, fmt.Sprintf("A circle has centre (%d,%d). What is the sum of the centre coordinates?", h, k), strconv.Itoa(h+k), "Add the coordinates.", f)
	}
}

func w28Conic(o int) generator.Problem {
	a := w28N(1, 30)
	c := w28N(1, 30)
	trace := a + c
	f := map[string]string{}
	w28I(f, "a", a)
	w28I(f, "c", c)
	w28I(f, "trace", trace)
	switch o {
	case 1:
		return w28P(o, fmt.Sprintf("A conic has matrix entries a=%d and c=%d on its diagonal. What is the trace?", a, c), strconv.Itoa(trace), "The trace sums the diagonal.", f)
	case 2:
		return w28P(o, fmt.Sprintf("A rotated conic has diagonal entry a=%d. What is a?", a), strconv.Itoa(a), "Read the entry.", f)
	default:
		return w28P(o, fmt.Sprintf("A rotated conic has diagonal entry c=%d. What is c?", c), strconv.Itoa(c), "Read the entry.", f)
	}
}

func w28Polar(o int) generator.Problem {
	twoa := w28N(2, 60)
	r := twoa / 2
	f := map[string]string{}
	w28I(f, "twoa", twoa)
	w28I(f, "r", r)
	switch o {
	case 1:
		return w28P(o, fmt.Sprintf("A polar curve has 2a=%d. What is a?", twoa), strconv.Itoa(r), "Halve 2a.", f)
	case 2:
		return w28P(o, fmt.Sprintf("A polar curve has 2a=%d. What is 2a?", twoa), strconv.Itoa(twoa), "The value is given.", f)
	default:
		return w28P(o, fmt.Sprintf("A polar curve has 2a=%d. What is the coefficient of a in 2a?", twoa), "2", "The coefficient is two.", f)
	}
}

func w28Shoelace(o int) generator.Problem {
	a := w28N(1, 30)
	b := w28N(1, 30)
	twice := a * b
	f := map[string]string{}
	w28I(f, "a", a)
	w28I(f, "b", b)
	w28I(f, "twice", twice)
	switch o {
	case 1:
		return w28P(o, fmt.Sprintf("A polygon has twice its area equal to %d by %d. What is twice the area?", a, b), strconv.Itoa(twice), "Multiply the two shoelace terms.", f)
	case 2:
		return w28P(o, fmt.Sprintf("A polygon has one shoelace term %d. What is it?", a), strconv.Itoa(a), "Read the term.", f)
	default:
		return w28P(o, fmt.Sprintf("A polygon has another shoelace term %d. What is it?", b), strconv.Itoa(b), "Read the term.", f)
	}
}

func w28Intersect(o int) generator.Problem {
	det := w28N(1, 60)
	slope := w28N(1, 30)
	f := map[string]string{}
	w28I(f, "det", det)
	w28I(f, "slope", slope)
	switch o {
	case 1:
		return w28P(o, fmt.Sprintf("Two lines have determinant %d. What is the determinant?", det), strconv.Itoa(det), "The determinant is given.", f)
	case 2:
		return w28P(o, fmt.Sprintf("A line has slope %d. What is the slope?", slope), strconv.Itoa(slope), "The slope is given.", f)
	default:
		return w28P(o, fmt.Sprintf("Two lines meet at one point. How many intersection points are there for slopes %d and %d?", slope, det), "1", "Non-parallel lines meet once.", f)
	}
}

func w28ThreeD(o int) generator.Problem {
	dx := w28N(1, 20)
	dy := w28N(1, 20)
	dz := w28N(1, 20)
	sq := dx*dx + dy*dy + dz*dz
	f := map[string]string{}
	w28I(f, "dx", dx)
	w28I(f, "dy", dy)
	w28I(f, "dz", dz)
	w28I(f, "sq", sq)
	switch o {
	case 1:
		return w28P(o, fmt.Sprintf("A displacement is (%d,%d,%d). What is its squared length?", dx, dy, dz), strconv.Itoa(sq), "Square and add the three components.", f)
	case 2:
		return w28P(o, fmt.Sprintf("A displacement is (%d,%d,%d). What is its x component?", dx, dy, dz), strconv.Itoa(dx), "Read the first component.", f)
	default:
		return w28P(o, fmt.Sprintf("A displacement is (%d,%d,%d). What is its y component?", dx, dy, dz), strconv.Itoa(dy), "Read the second component.", f)
	}
}

func w28Thales(o int) generator.Problem {
	d := 2 * w28N(1, 30)
	r := d / 2
	f := map[string]string{}
	w28I(f, "d", d)
	w28I(f, "r", r)
	switch o {
	case 1:
		return w28P(o, fmt.Sprintf("A circle has diameter %d. What is its radius?", d), strconv.Itoa(r), "Halve the diameter.", f)
	case 2:
		return w28P(o, fmt.Sprintf("A circle has diameter %d. What is the diameter?", d), strconv.Itoa(d), "The diameter is given.", f)
	default:
		return w28P(o, fmt.Sprintf("A circle has diameter %d. The radius is the diameter divided by what number?", d), "2", "A diameter is twice the radius.", f)
	}
}

func w28ArcSector(o int) generator.Problem {
	theta := 2 * w28N(1, 30)
	r := w28N(1, 30)
	arc := theta * r
	sector := theta * r * r / 2
	f := map[string]string{}
	w28I(f, "theta", theta)
	w28I(f, "r", r)
	w28I(f, "arc", arc)
	w28I(f, "sector", sector)
	switch o {
	case 1:
		return w28P(o, fmt.Sprintf("A sector has angle %d radians and radius %d. What is its arc length?", theta, r), strconv.Itoa(arc), "Arc length is radius times angle in radians.", f)
	case 2:
		return w28P(o, fmt.Sprintf("A sector has angle %d radians and radius %d. What is its area?", theta, r), strconv.Itoa(sector), "Sector area is half the radius squared times the angle.", f)
	default:
		return w28P(o, fmt.Sprintf("A sector has angle %d radians. What is the radius if the arc length is %d?", theta, arc), strconv.Itoa(r), "Divide the arc length by the angle.", f)
	}
}

func w28Segment(o int) generator.Problem {
	r := w28N(1, 30)
	theta := 2 * w28N(1, 30)
	sector := theta * r * r / 2
	f := map[string]string{}
	w28I(f, "r", r)
	w28I(f, "theta", theta)
	w28I(f, "sector", sector)
	switch o {
	case 1:
		return w28P(o, fmt.Sprintf("A segment sits inside a sector of area %d. What is the sector area?", sector), strconv.Itoa(sector), "The sector area is given.", f)
	case 2:
		return w28P(o, fmt.Sprintf("A segment sits inside a sector of radius %d. What is the radius?", r), strconv.Itoa(r), "The radius is given.", f)
	default:
		return w28P(o, fmt.Sprintf("A segment sits inside a sector of angle %d radians. What is the angle?", theta), strconv.Itoa(theta), "The angle is given.", f)
	}
}

func w28Tangents(o int) generator.Problem {
	r1 := w28N(1, 20)
	r2 := w28N(1, 20)
	sum := r1 + r2
	d := sum + 2*w28N(1, 20)
	internal := 2
	tangents := 4
	f := map[string]string{}
	w28I(f, "r1", r1)
	w28I(f, "r2", r2)
	w28I(f, "d", d)
	w28I(f, "sum", sum)
	w28I(f, "internal", internal)
	w28I(f, "tangents", tangents)
	switch o {
	case 1:
		return w28P(o, fmt.Sprintf("Two circles have radii %d and %d. What is the sum of the radii?", r1, r2), strconv.Itoa(sum), "Add the radii.", f)
	case 2:
		return w28P(o, fmt.Sprintf("Two separated circles of radii %d and %d have how many common tangents?", r1, r2), strconv.Itoa(tangents), "Separated circles have four common tangents.", f)
	default:
		return w28P(o, fmt.Sprintf("Two circles have centre distance %d. What is the centre distance?", d), strconv.Itoa(d), "The distance is given.", f)
	}
}

func w28Orthogonal(o int) generator.Problem {
	r1 := w28N(1, 20)
	r2 := w28N(1, 20)
	d2 := r1*r1 + r2*r2
	sum := r1 + r2
	f := map[string]string{}
	w28I(f, "r1", r1)
	w28I(f, "r2", r2)
	w28I(f, "d2", d2)
	w28I(f, "sum", sum)
	switch o {
	case 1:
		return w28P(o, fmt.Sprintf("Two circles of radii %d and %d are orthogonal. What is the squared centre distance?", r1, r2), strconv.Itoa(d2), "Orthogonality gives d squared equal to the sum of squared radii.", f)
	case 2:
		return w28P(o, fmt.Sprintf("Two circles of radii %d and %d are orthogonal. What is the sum of the radii?", r1, r2), strconv.Itoa(sum), "Add the radii.", f)
	default:
		return w28P(o, fmt.Sprintf("A circle has radius %d. What is the radius?", r1), strconv.Itoa(r1), "The radius is given.", f)
	}
}

func w28Inversion(o int) generator.Problem {
	k := w28N(1, 20)
	d := w28N(1, 20)
	R := d * k
	r2 := w28N(1, 20)
	image := R * R / d
	f := map[string]string{}
	w28I(f, "R", R)
	w28I(f, "r2", r2)
	w28I(f, "d", d)
	w28I(f, "image", image)
	switch o {
	case 1:
		return w28P(o, fmt.Sprintf("Inversion in a circle of radius %d sends a point at distance %d to the reciprocal distance scaled by R squared. What is the image distance?", R, d), strconv.Itoa(image), "The image distance is R squared over d.", f)
	case 2:
		return w28P(o, fmt.Sprintf("An inverted point has squared radius %d. What is that value?", r2), strconv.Itoa(r2), "The value is given.", f)
	default:
		return w28P(o, fmt.Sprintf("An inversion circle has radius %d. What is the radius?", R), strconv.Itoa(R), "The radius is given.", f)
	}
}

func w28NinePoint(o int) generator.Problem {
	R := 2 * w28N(1, 30)
	r9 := R / 2
	f := map[string]string{}
	w28I(f, "R", R)
	w28I(f, "r9", r9)
	switch o {
	case 1:
		return w28P(o, fmt.Sprintf("A triangle has circumradius %d. What is the radius of its nine-point circle?", R), strconv.Itoa(r9), "The nine-point radius is half the circumradius.", f)
	case 2:
		return w28P(o, fmt.Sprintf("A triangle has circumradius %d. What is the circumradius?", R), strconv.Itoa(R), "The circumradius is given.", f)
	default:
		return w28P(o, fmt.Sprintf("A nine-point circle has radius %d. The circumradius is how many times that?", r9), "2", "The circumradius is twice the nine-point radius.", f)
	}
}

func w28EulerLine(o int) generator.Problem {
	R := w28N(1, 30)
	a := w28N(1, 30)
	b := w28N(1, 30)
	ratio := 2
	f := map[string]string{}
	w28I(f, "R", R)
	w28I(f, "a", a)
	w28I(f, "b", b)
	w28I(f, "ratio", ratio)
	switch o {
	case 1:
		return w28P(o, fmt.Sprintf("On the Euler line the centroid divides the segment from orthocentre to circumcentre in the ratio %d to 1. What is that first number?", R), strconv.Itoa(ratio), "The centroid divides it two to one.", f)
	case 2:
		return w28P(o, fmt.Sprintf("An Euler line triangle has first segment length %d. What is that length?", a), strconv.Itoa(a), "The length is given.", f)
	default:
		return w28P(o, fmt.Sprintf("An Euler line triangle has second segment length %d. What is that length?", b), strconv.Itoa(b), "The length is given.", f)
	}
}

func w28Tetra(o int) generator.Problem {
	a := w28N(1, 30)
	faces := 4
	edges := 6
	f := map[string]string{}
	w28I(f, "a", a)
	w28I(f, "faces", faces)
	w28I(f, "edges", edges)
	switch o {
	case 1:
		return w28P(o, fmt.Sprintf("A regular tetrahedron has edge %d. How many faces does it have?", a), strconv.Itoa(faces), "A tetrahedron has four faces.", f)
	case 2:
		return w28P(o, fmt.Sprintf("A regular tetrahedron has edge %d. What is the edge length?", a), strconv.Itoa(a), "The edge length is given.", f)
	default:
		return w28P(o, fmt.Sprintf("A regular tetrahedron has edge %d. How many edges does it have?", a), strconv.Itoa(edges), "A tetrahedron has six edges.", f)
	}
}

func w28Octa(o int) generator.Problem {
	a := w28N(1, 30)
	v := 6
	e := 12
	fc := 8
	f := map[string]string{}
	w28I(f, "a", a)
	w28I(f, "v", v)
	w28I(f, "e", e)
	w28I(f, "f", fc)
	switch o {
	case 1:
		return w28P(o, fmt.Sprintf("A regular octahedron has edge %d. How many edges does it have?", a), strconv.Itoa(e), "An octahedron has twelve edges.", f)
	case 2:
		return w28P(o, fmt.Sprintf("A regular octahedron has edge %d. How many vertices does it have?", a), strconv.Itoa(v), "An octahedron has six vertices.", f)
	default:
		return w28P(o, fmt.Sprintf("A regular octahedron has edge %d. How many faces does it have?", a), strconv.Itoa(fc), "An octahedron has eight faces.", f)
	}
}

func w28Platonic(o int) generator.Problem {
	faces := []int{4, 6, 8, 12, 20}[rand.Intn(5)]
	a := w28N(1, 30)
	f := map[string]string{}
	w28I(f, "f", faces)
	w28I(f, "a", a)
	switch o {
	case 1:
		return w28P(o, fmt.Sprintf("One Platonic solid with edge %d has %d faces. How many faces does it have?", a, faces), strconv.Itoa(faces), "The face count is given.", f)
	case 2:
		return w28P(o, fmt.Sprintf("A Platonic solid with edge %d is one of how many Platonic solids?", a), "5", "There are five Platonic solids.", f)
	default:
		return w28P(o, fmt.Sprintf("A Platonic solid with edge %d has how many face types?", a), "1", "Every face is the same regular polygon.", f)
	}
}

func w28Oblique(o int) generator.Problem {
	base := w28N(1, 30)
	height := w28N(1, 30)
	volume := base * height
	f := map[string]string{}
	w28I(f, "base", base)
	w28I(f, "height", height)
	w28I(f, "volume", volume)
	switch o {
	case 1:
		return w28P(o, fmt.Sprintf("An oblique prism has base area %d and height %d. What is its volume?", base, height), strconv.Itoa(volume), "Volume is base area times height.", f)
	case 2:
		return w28P(o, fmt.Sprintf("An oblique prism has base area %d. What is the base area?", base), strconv.Itoa(base), "The base area is given.", f)
	default:
		return w28P(o, fmt.Sprintf("An oblique prism has height %d. What is the height?", height), strconv.Itoa(height), "The height is given.", f)
	}
}

func w28Composite(o int) generator.Problem {
	v1 := w28N(1, 30)
	v2 := w28N(1, 30)
	v3 := w28N(1, 30)
	total := v1 + v2 + v3
	f := map[string]string{}
	w28I(f, "v1", v1)
	w28I(f, "v2", v2)
	w28I(f, "v3", v3)
	w28I(f, "total", total)
	w28I(f, "n", 3)
	switch o {
	case 1:
		return w28P(o, fmt.Sprintf("A composite solid is made of parts of volume %d, %d and %d. What is the total volume?", v1, v2, v3), strconv.Itoa(total), "Add the part volumes.", f)
	case 2:
		return w28P(o, fmt.Sprintf("A composite solid has first part volume %d. What is it?", v1), strconv.Itoa(v1), "The volume is given.", f)
	default:
		return w28P(o, fmt.Sprintf("A composite solid of total volume %d has a first part of volume %d. What volume do the remaining parts contribute?", total, v1), strconv.Itoa(v2+v3), "Subtract the first part from the total.", f)
	}
}

func w28Cavalieri(o int) generator.Problem {
	count := w28N(3, 60)
	f := map[string]string{}
	w28I(f, "count", count)
	switch o {
	case 1:
		return w28P(o, fmt.Sprintf("Cavalieri's principle compares %d parallel slices of equal area. How many slices are compared?", count), strconv.Itoa(count), "The slice count is given.", f)
	case 2:
		return w28P(o, fmt.Sprintf("Two solids have equal cross-sectional areas at every one of %d heights. What is the difference in their volumes?", count), "0", "Cavalieri's principle gives equal volumes.", f)
	default:
		return w28P(o, fmt.Sprintf("A solid is cut into %d equal-area slices. How many volumes must agree for the principle to apply?", count), "1", "The volumes agree as one value.", f)
	}
}

func w28Rotation(o int) generator.Problem {
	r := w28N(1, 20)
	h := w28N(1, 30)
	r2 := r * r
	r2h := r2 * h
	f := map[string]string{}
	w28I(f, "r", r)
	w28I(f, "h", h)
	w28I(f, "r2", r2)
	w28I(f, "r2h", r2h)
	w28I(f, "volume", r2h)
	switch o {
	case 1:
		return w28P(o, fmt.Sprintf("A cylinder of radius %d and height %d is written with volume proportional to r squared h. What is r squared h?", r, h), strconv.Itoa(r2h), "Square the radius and multiply by the height.", f)
	case 2:
		return w28P(o, fmt.Sprintf("A solid of revolution has radius %d. What is r squared?", r), strconv.Itoa(r2), "Square the radius.", f)
	default:
		return w28P(o, fmt.Sprintf("A solid of revolution has height %d. What is the height?", h), strconv.Itoa(h), "The height is given.", f)
	}
}

func w28Truncated(o int) generator.Problem {
	R := w28N(2, 20)
	r := w28N(1, R-1)
	h := w28N(1, 30)
	factor := R*R + R*r + r*r
	volume := h * factor
	f := map[string]string{}
	w28I(f, "R", R)
	w28I(f, "r", r)
	w28I(f, "h", h)
	w28I(f, "factor", factor)
	w28I(f, "volume", volume)
	switch o {
	case 1:
		return w28P(o, fmt.Sprintf("A truncated cone has radii %d and %d. What is R squared plus R r plus r squared?", R, r), strconv.Itoa(factor), "Substitute into the three-term factor.", f)
	case 2:
		return w28P(o, fmt.Sprintf("A truncated cone has height %d and factor %d. What is the product of height and factor?", h, factor), strconv.Itoa(volume), "Multiply the height by the factor.", f)
	default:
		return w28P(o, fmt.Sprintf("A truncated cone has height %d. What is the height?", h), strconv.Itoa(h), "The height is given.", f)
	}
}

func w28Apothem(o int) generator.Problem {
	n := w28N(3, 30)
	s := 2 * w28N(1, 20)
	a := w28N(1, 30)
	area := n * s * a / 2
	f := map[string]string{}
	w28I(f, "n", n)
	w28I(f, "s", s)
	w28I(f, "a", a)
	w28I(f, "area", area)
	switch o {
	case 1:
		return w28P(o, fmt.Sprintf("A regular %d-gon has apothem %d. What is the apothem?", n, a), strconv.Itoa(a), "The apothem is given.", f)
	case 2:
		return w28P(o, fmt.Sprintf("A regular polygon has %d sides. How many sides does it have?", n), strconv.Itoa(n), "The side count is given.", f)
	default:
		return w28P(o, fmt.Sprintf("A regular %d-gon has side %d and apothem %d. What is its area?", n, s, a), strconv.Itoa(area), "Area is half the perimeter times the apothem.", f)
	}
}

func w28Star(o int) generator.Problem {
	n := w28N(5, 30)
	k := w28N(2, n-1)
	g := w28GCD(n, k)
	f := map[string]string{}
	w28I(f, "n", n)
	w28I(f, "k", k)
	w28I(f, "g", g)
	switch o {
	case 1:
		return w28P(o, fmt.Sprintf("A star polygon {n/k} has n=%d and step k=%d. What is k?", n, k), strconv.Itoa(k), "The step is given.", f)
	case 2:
		return w28P(o, fmt.Sprintf("A star polygon {n/k} has n=%d vertices. How many vertices does it have?", n), strconv.Itoa(n), "The vertex count is given.", f)
	default:
		return w28P(o, fmt.Sprintf("A star polygon uses n=%d vertices and step k=%d. What is the greatest common divisor of n and k?", n, k), strconv.Itoa(g), "Take the greatest common divisor.", f)
	}
}

func w28Constructible(o int) generator.Problem {
	n := w28N(3, 60)
	primes := w28N(1, 30)
	f := map[string]string{}
	w28I(f, "n", n)
	w28I(f, "primes", primes)
	switch o {
	case 1:
		return w28P(o, fmt.Sprintf("A regular %d-gon is constructible when its odd part is a product of distinct Fermat primes. How many Fermat primes does this statement use at most?", n), "5", "There are five known Fermat primes.", f)
	case 2:
		return w28P(o, fmt.Sprintf("A constructibility test examines a polygon with %d sides. How many sides does it examine?", n), strconv.Itoa(n), "The side count is given.", f)
	default:
		return w28P(o, fmt.Sprintf("A constructibility test for a %d-gon records %d Fermat prime factors. How many factors does it record?", n, primes), strconv.Itoa(primes), "The count is given.", f)
	}
}

func w28QuadArea(o int) generator.Problem {
	a := w28N(1, 30)
	b := w28N(1, 30)
	area := a * b
	f := map[string]string{}
	w28I(f, "a", a)
	w28I(f, "b", b)
	w28I(f, "area", area)
	switch o {
	case 1:
		return w28P(o, fmt.Sprintf("A rectangle has sides %d and %d. What is its area?", a, b), strconv.Itoa(area), "Multiply the sides.", f)
	case 2:
		return w28P(o, fmt.Sprintf("A quadrilateral has first side %d. What is it?", a), strconv.Itoa(a), "The side is given.", f)
	default:
		return w28P(o, fmt.Sprintf("A quadrilateral has second side %d. What is it?", b), strconv.Itoa(b), "The side is given.", f)
	}
}

func w28Homothety(o int) generator.Problem {
	k := w28N(1, 20)
	d := w28N(1, 30)
	image := k * d
	f := map[string]string{}
	w28I(f, "k", k)
	w28I(f, "d", d)
	w28I(f, "image", image)
	switch o {
	case 1:
		return w28P(o, fmt.Sprintf("A homothety of factor %d sends a segment of length %d to what length?", k, d), strconv.Itoa(image), "Multiply the length by the factor.", f)
	case 2:
		return w28P(o, fmt.Sprintf("A homothety has factor %d. What is the factor?", k), strconv.Itoa(k), "The factor is given.", f)
	default:
		return w28P(o, fmt.Sprintf("A homothety moves a point a distance %d. What is the distance?", d), strconv.Itoa(d), "The distance is given.", f)
	}
}

func w28Glide(o int) generator.Problem {
	d := w28N(1, 60)
	n := w28N(1, 30)
	f := map[string]string{}
	w28I(f, "d", d)
	w28I(f, "n", n)
	switch o {
	case 1:
		return w28P(o, fmt.Sprintf("A glide reflection translates by %d before reflecting. What is the translation length?", d), strconv.Itoa(d), "The translation length is given.", f)
	case 2:
		return w28P(o, fmt.Sprintf("A glide reflection is applied %d times. How many applications are there?", n), strconv.Itoa(n), "The count is given.", f)
	default:
		return w28P(o, fmt.Sprintf("A glide reflection translates by %d before reflecting. How many transformations does it combine?", d), "2", "It combines a reflection and a translation.", f)
	}
}

func w28Group(o int) generator.Problem {
	n := w28N(1, 30)
	f := map[string]string{}
	w28I(f, "n", n)
	switch o {
	case 1:
		return w28P(o, fmt.Sprintf("The dihedral group of a regular %d-gon has how many rotations?", n), strconv.Itoa(n), "One rotation per vertex.", f)
	case 2:
		return w28P(o, fmt.Sprintf("The dihedral group of a regular %d-gon has how many elements?", n), strconv.Itoa(2*n), "There are n rotations and n reflections.", f)
	default:
		return w28P(o, fmt.Sprintf("A transformation group on a regular %d-gon has how many element types?", n), "2", "Rotations and reflections are the two types.", f)
	}
}
