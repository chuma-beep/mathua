package geometry

// Wave 10 geometry, part two. Every branch publishes a, b, c and answer so the
// shared authored schema resolves for all three knowledge points.

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave10Rest(reg *generator.Registry) {
	reg.Register("geo.triangle.orthocenter", &w10gOrthocenterGen{})
	reg.Register("geo.triangle.circumcenter", &w10gCircumcenterGen{})
	reg.Register("geo.triangle.incenter", &w10gIncenterGen{})
	reg.Register("geo.triangle.cevian", &w10gCevaGen{})
	reg.Register("geo.triangle.apollonius", &w10gApolloniusGen{})
	reg.Register("geo.circle.sagitta", &w10gSagittaGen{})
	reg.Register("geo.circle.radical_axis", &w10gRadicalAxisGen{})
	reg.Register("geo.polygon.diagonal_count", &w10gDiagonalGen{})
	reg.Register("geo.polygon.tessellation", &w10gTessellationGen{})
	reg.Register("geo.coord.sphere_equation", &w10gSphereGen{})
	reg.Register("geo.transform.composition", &w10gCompositionGen{})
	reg.Register("geo.transform.isometry", &w10gIsometryGen{})
	reg.Register("geo.solid.hemisphere", &w10gHemisphereGen{})
	reg.Register("geo.solid.density", &w10gDensityGen{})
}

func w10gP(o int, q, a, ex string, facts map[string]string) generator.Problem {
	if facts == nil {
		facts = map[string]string{}
	}
	facts["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: facts, KPOrdinal: o}
}

func w10gI(m map[string]string, k string, v int) { m[k] = strconv.Itoa(v) }

type w10gOrthocenterGen struct{}

func (g *w10gOrthocenterGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10GOrd(ctx)
	for {
		b := 6 + rand.Intn(10)
		cx := 2 + rand.Intn(b-3)
		cy := 3 + rand.Intn(10)
		if (cx*(b-cx))%cy != 0 {
			continue
		}
		hy := cx * (b - cx) / cy
		facts := map[string]string{}
		w10gI(facts, "a", cx)
		w10gI(facts, "b", cy)
		w10gI(facts, "c", b)
		w10gI(facts, "hx", cx)
		w10gI(facts, "hy", hy)
		switch o {
		case 1:
			return w10gP(o, fmt.Sprintf("A triangle has vertices A(0, 0), B(%d, 0) and C(%d, %d). What is the x-coordinate of its orthocenter?", b, cx, cy), strconv.Itoa(cx), fmt.Sprintf("The altitude from C is vertical, so the orthocenter has x = %d.", cx), facts)
		case 2:
			return w10gP(o, fmt.Sprintf("A triangle has vertices A(0, 0), B(%d, 0) and C(%d, %d). What is the y-coordinate of its orthocenter?", b, cx, cy), strconv.Itoa(hy), fmt.Sprintf("The altitude from A meets x = %d at y = %d.", cx, hy), facts)
		default:
			return w10gP(o, fmt.Sprintf("A triangle has vertices A(0, 0), B(%d, 0) and C(%d, %d). What is the sum of the coordinates of its orthocenter?", b, cx, cy), strconv.Itoa(cx+hy), fmt.Sprintf("The orthocenter is (%d, %d), so the sum is %d.", cx, hy, cx+hy), facts)
		}
	}
}

type w10gCircumcenterGen struct{}

func (g *w10gCircumcenterGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10GOrd(ctx)
	s := 1 + rand.Intn(10)
	b, c, hyp := 6*s, 8*s, 10*s
	x, y := b/2, c/2
	r := hyp / 2
	facts := map[string]string{}
	w10gI(facts, "a", b)
	w10gI(facts, "b", c)
	w10gI(facts, "c", r)
	w10gI(facts, "x", x)
	w10gI(facts, "y", y)
	switch o {
	case 1:
		return w10gP(o, fmt.Sprintf("A right triangle has vertices (0, 0), (%d, 0) and (0, %d). What is the x-coordinate of its circumcenter?", b, c), strconv.Itoa(x), fmt.Sprintf("The circumcenter is the midpoint of the hypotenuse, so its x-coordinate is %d/2 = %d.", b, x), facts)
	case 2:
		return w10gP(o, fmt.Sprintf("A right triangle has legs %d and %d. What is its circumradius?", b, c), strconv.Itoa(r), fmt.Sprintf("The hypotenuse is %d and the circumradius is half of it: %d.", hyp, r), facts)
	default:
		return w10gP(o, fmt.Sprintf("A right triangle has vertices (0, 0), (%d, 0) and (0, %d). What is the sum of the coordinates of its circumcenter?", b, c), strconv.Itoa(x+y), fmt.Sprintf("The circumcenter is (%d, %d), so the sum is %d.", x, y, x+y), facts)
	}
}

type w10gIncenterGen struct{}

func (g *w10gIncenterGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10GOrd(ctx)
	k := 1 + rand.Intn(10)
	b, c, hyp := 3*k, 4*k, 5*k
	r := (b + c - hyp) / 2
	facts := map[string]string{}
	w10gI(facts, "a", b)
	w10gI(facts, "b", c)
	w10gI(facts, "c", hyp)
	w10gI(facts, "r", r)
	switch o {
	case 1:
		return w10gP(o, fmt.Sprintf("A right triangle has legs %d and %d and hypotenuse %d. What is its inradius?", b, c, hyp), strconv.Itoa(r), fmt.Sprintf("For a right triangle r = (a + b - c)/2 = (%d + %d - %d)/2 = %d.", b, c, hyp, r), facts)
	case 2:
		return w10gP(o, fmt.Sprintf("A right triangle with legs %d and %d has incenter at (%d, %d). What is the sum of the incenter coordinates?", b, c, r, r), strconv.Itoa(2*r), fmt.Sprintf("The incenter lies r from each leg, so its coordinates sum to %d + %d = %d.", r, r, 2*r), facts)
	default:
		return w10gP(o, fmt.Sprintf("A right triangle has area %d and semiperimeter %d. What is its inradius?", b*c/2, (b+c+hyp)/2), strconv.Itoa(r), fmt.Sprintf("Use A = rs, so r = A/s = %d/%d = %d.", b*c/2, (b+c+hyp)/2, r), facts)
	}
}

type w10gCevaGen struct{}

func (g *w10gCevaGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10GOrd(ctx)
	x, y, z := 1+rand.Intn(5), 1+rand.Intn(5), 1+rand.Intn(5)
	product := x * y * z
	facts := map[string]string{}
	w10gI(facts, "a", x)
	w10gI(facts, "b", y)
	w10gI(facts, "c", z)
	w10gI(facts, "product", product)
	switch o {
	case 1:
		return w10gP(o, fmt.Sprintf("Three cevians have ratios %d, %d and %d. What is their Ceva product?", x, y, z), strconv.Itoa(product), fmt.Sprintf("Multiply the three ratios: %d x %d x %d = %d.", x, y, z, product), facts)
	case 2:
		return w10gP(o, fmt.Sprintf("Three cevians have product %d. Two ratios are %d and %d. What is the third ratio?", product, x, y), strconv.Itoa(z), fmt.Sprintf("Divide the product by the two known ratios: %d / (%d x %d) = %d.", product, x, y, z), facts)
	default:
		need := 1.0 / float64(x*y)
		return w10gP(o, fmt.Sprintf("Two cevian ratios are %d and %d. What third ratio makes the Ceva product equal 1?", x, y), strconv.FormatFloat(need, 'f', -1, 64), fmt.Sprintf("The product is %d x %d x r = 1, so r = 1/%d = %s.", x, y, x*y, strconv.FormatFloat(need, 'f', -1, 64)), facts)
	}
}

type w10gApolloniusGen struct{}

func (g *w10gApolloniusGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10GOrd(ctx)
	triples := [][3]int{{3, 4, 5}, {5, 12, 13}, {6, 8, 10}, {7, 24, 25}, {8, 15, 17}, {9, 12, 15}, {9, 40, 41}, {12, 16, 20}, {12, 35, 37}, {20, 21, 29}, {28, 45, 53}}
	t := triples[rand.Intn(len(triples))]
	d, h, s := t[0], t[1], t[2]
	facts := map[string]string{}
	w10gI(facts, "a", d)
	w10gI(facts, "b", h)
	w10gI(facts, "c", s)
	switch o {
	case 1:
		return w10gP(o, fmt.Sprintf("An isosceles triangle has base %d and equal sides %d. What is the length of the median to the base?", 2*d, s), strconv.Itoa(h), fmt.Sprintf("The median splits the base into %d. By the Pythagorean theorem its length is %d.", d, h), facts)
	case 2:
		return w10gP(o, fmt.Sprintf("An isosceles triangle has half-base %d and median %d. What is each equal side?", d, h), strconv.Itoa(s), fmt.Sprintf("The side is the hypotenuse of a right triangle with legs %d and %d: %d.", d, h, s), facts)
	default:
		return w10gP(o, fmt.Sprintf("For an isosceles triangle with half-base %d and equal side %d, what is the square of the median?", d, s), strconv.Itoa(s*s-d*d), fmt.Sprintf("Apollonius gives m^2 = s^2 - d^2 = %d - %d = %d.", s*s, d*d, s*s-d*d), facts)
	}
}

type w10gSagittaGen struct{}

func (g *w10gSagittaGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10GOrd(ctx)
	rows := [][3]int{{5, 1, 3}, {5, 2, 4}, {10, 2, 6}, {13, 1, 5}, {13, 8, 12}, {17, 2, 8}, {17, 9, 15}, {25, 1, 7}, {25, 18, 24}, {29, 8, 20}, {29, 9, 21}, {41, 1, 9}, {41, 32, 40}}
	row := rows[rand.Intn(len(rows))]
	r, s, a := row[0], row[1], row[2]
	facts := map[string]string{}
	w10gI(facts, "a", a)
	w10gI(facts, "b", s)
	w10gI(facts, "c", r)
	switch o {
	case 1:
		return w10gP(o, fmt.Sprintf("A circle has a chord of length %d and a sagitta of %d. What is its radius?", 2*a, s), strconv.Itoa(r), fmt.Sprintf("Use R = (a^2 + s^2)/(2s) with half-chord %d: (%d + %d)/(%d) = %d.", a, a*a, s*s, 2*s, r), facts)
	case 2:
		return w10gP(o, fmt.Sprintf("A circle has radius %d and a chord of length %d. What is the sagitta?", r, 2*a), strconv.Itoa(s), fmt.Sprintf("The distance from the center to the chord is %d. The sagitta is %d - %d = %d.", r-s, r, r-s, s), facts)
	default:
		return w10gP(o, fmt.Sprintf("A circle has radius %d and sagitta %d. What is the chord length?", r, s), strconv.Itoa(2*a), fmt.Sprintf("Half the chord is sqrt(R^2 - (R-s)^2) = %d, so the chord is %d.", a, 2*a), facts)
	}
}

type w10gRadicalAxisGen struct{}

func (g *w10gRadicalAxisGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10GOrd(ctx)
	d := 2 * (2 + rand.Intn(6))
	r := 3 + rand.Intn(8)
	intersections := 2
	if d == 2*r {
		intersections = 1
	} else if d > 2*r {
		intersections = 0
	}
	x := d / 2
	facts := map[string]string{}
	w10gI(facts, "a", d)
	w10gI(facts, "b", r)
	w10gI(facts, "c", intersections)
	w10gI(facts, "x", x)
	switch o {
	case 1:
		return w10gP(o, fmt.Sprintf("Two circles have equal radii %d and centers (%d, 0) apart. What is the x-coordinate of their radical axis?", r, d), strconv.Itoa(x), fmt.Sprintf("Equal radii put the radical axis halfway between the centers: %d/2 = %d.", d, x), facts)
	case 2:
		return w10gP(o, fmt.Sprintf("Two equal circles have radius %d and centers %d apart. At the left center, what is the difference of the two powers?", r, d), strconv.Itoa(d*d), fmt.Sprintf("The power difference is 2dx - d^2 at x = 0, which is %d.", d*d), facts)
	default:
		return w10gP(o, fmt.Sprintf("Two circles have radius %d and centers %d apart. How many intersection points do they have?", r, d), strconv.Itoa(intersections), fmt.Sprintf("Compare the center distance %d with twice the radius %d: the circles meet in %d point(s).", d, 2*r, intersections), facts)
	}
}

type w10gDiagonalGen struct{}

func (g *w10gDiagonalGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10GOrd(ctx)
	n := 5 + rand.Intn(16)
	diag := n * (n - 3) / 2
	segments := n * (n - 1) / 2
	facts := map[string]string{}
	w10gI(facts, "a", n)
	w10gI(facts, "b", diag)
	w10gI(facts, "c", segments)
	switch o {
	case 1:
		return w10gP(o, fmt.Sprintf("How many diagonals does a convex %d-gon have?", n), strconv.Itoa(diag), fmt.Sprintf("Use n(n-3)/2 = %d x %d / 2 = %d.", n, n-3, diag), facts)
	case 2:
		return w10gP(o, fmt.Sprintf("A convex polygon has %d diagonals. How many sides does it have?", diag), strconv.Itoa(n), fmt.Sprintf("Solve n(n-3)/2 = %d; the positive solution is n = %d.", diag, n), facts)
	default:
		return w10gP(o, fmt.Sprintf("How many segments join pairs of vertices in a convex %d-gon?", n), strconv.Itoa(segments), fmt.Sprintf("Use n(n-1)/2 = %d x %d / 2 = %d.", n, n-1, segments), facts)
	}
}

type w10gTessellationGen struct{}

func (g *w10gTessellationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10GOrd(ctx)
	n := 3 + rand.Intn(10)
	if o == 3 {
		n = []int{3, 4, 6}[rand.Intn(3)]
	}
	sum := (n - 2) * 180
	angle := sum / n
	used := 1 + rand.Intn(2*n/(n-2)-1)
	remain := 360 - used*angle
	facts := map[string]string{}
	w10gI(facts, "a", n)
	w10gI(facts, "b", sum)
	w10gI(facts, "c", angle)
	w10gI(facts, "used", used)
	switch o {
	case 1:
		return w10gP(o, fmt.Sprintf("What is the sum of the interior angles of a %d-gon?", n), strconv.Itoa(sum), fmt.Sprintf("Use (n-2)180 = (%d - 2)180 = %d.", n, sum), facts)
	case 2:
		return w10gP(o, fmt.Sprintf("What is each interior angle of a regular %d-gon?", n), strconv.Itoa(angle), fmt.Sprintf("Divide the angle sum by n: %d/%d = %d.", sum, n, angle), facts)
	default:
		return w10gP(o, fmt.Sprintf("A vertex already has %d regular %d-gons meeting there. How many degrees remain before the angles fill 360?", used, n), strconv.Itoa(remain), fmt.Sprintf("Each interior angle is %d degrees. After %d copies, the remaining angle is 360 - %d x %d = %d.", angle, used, used, angle, remain), facts)
	}
}

func w10gSquareTerm(term string, center int) string {
	if center == 0 {
		return term + "^2"
	}
	return fmt.Sprintf("(%s - %d)^2", term, center)
}

type w10gSphereGen struct{}

func (g *w10gSphereGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10GOrd(ctx)
	h, k, l := rand.Intn(9), rand.Intn(9), rand.Intn(9)
	r := 2 + rand.Intn(14)
	px, py, pz := h+r, k, l
	if rand.Intn(2) == 0 {
		px, py, pz = h+1, k+2, l+2
	}
	dist2 := (px-h)*(px-h) + (py-k)*(py-k) + (pz-l)*(pz-l)
	facts := map[string]string{}
	w10gI(facts, "a", h)
	w10gI(facts, "b", k)
	w10gI(facts, "c", l)
	w10gI(facts, "r", r)
	w10gI(facts, "dist2", dist2)
	w10gI(facts, "px", px)
	w10gI(facts, "py", py)
	w10gI(facts, "pz", pz)
	switch o {
	case 1:
		return w10gP(o, fmt.Sprintf("A sphere has equation %s + %s + %s = %d. What is its radius?", w10gSquareTerm("x", h), w10gSquareTerm("y", k), w10gSquareTerm("z", l), r*r), strconv.Itoa(r), fmt.Sprintf("The right side is r^2 = %d, so r = %d.", r*r, r), facts)
	case 2:
		return w10gP(o, fmt.Sprintf("A sphere has equation %s + %s + %s = %d. What is the x-coordinate of its center?", w10gSquareTerm("x", h), w10gSquareTerm("y", k), w10gSquareTerm("z", l), r*r), strconv.Itoa(h), fmt.Sprintf("The center is (%d, %d, %d), so its x-coordinate is %d.", h, k, l, h), facts)
	default:
		return w10gP(o, fmt.Sprintf("A sphere has center (%d, %d, %d) and radius %d. What is the squared distance from the center to (%d, %d, %d)?", h, k, l, r, px, py, pz), strconv.Itoa(dist2), fmt.Sprintf("Add the squared coordinate differences: %d.", dist2), facts)
	}
}

type w10gCompositionGen struct{}

func (g *w10gCompositionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10GOrd(ctx)
	x, y := -6+rand.Intn(13), -6+rand.Intn(13)
	tx, ty := -5+rand.Intn(11), -5+rand.Intn(11)
	facts := map[string]string{}
	w10gI(facts, "a", x)
	w10gI(facts, "b", y)
	w10gI(facts, "c", tx)
	w10gI(facts, "tx", tx)
	w10gI(facts, "ty", ty)
	switch o {
	case 1:
		return w10gP(o, fmt.Sprintf("Translate (%d, %d) by (%d, %d) and then by (%d, %d). What is the sum of the final coordinates?", x, y, tx, ty, -tx, -ty), strconv.Itoa(x+y), fmt.Sprintf("The translations cancel, leaving (%d, %d), whose coordinates sum to %d.", x, y, x+y), facts)
	case 2:
		return w10gP(o, fmt.Sprintf("Rotate (%d, %d) by 90 degrees counterclockwise and then translate by (%d, %d). What is the final x-coordinate?", x, y, tx, ty), strconv.Itoa(-y+tx), fmt.Sprintf("Rotation gives (%d, %d); adding the translation gives x = %d.", -y, x, -y+tx), facts)
	default:
		return w10gP(o, fmt.Sprintf("After translating by (%d, %d) a point has y-coordinate %d. What was its original y-coordinate?", tx, ty, y+ty), strconv.Itoa(y), fmt.Sprintf("Subtract the translation: %d - (%d) = %d.", y+ty, ty, y), facts)
	}
}

type w10gIsometryGen struct{}

func (g *w10gIsometryGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10GOrd(ctx)
	x1, y1 := -5+rand.Intn(11), -5+rand.Intn(11)
	x2, y2 := -5+rand.Intn(11), -5+rand.Intn(11)
	dx, dy := x2-x1, y2-y1
	dist2 := dx*dx + dy*dy
	facts := map[string]string{}
	w10gI(facts, "a", x1)
	w10gI(facts, "b", y1)
	w10gI(facts, "c", dist2)
	w10gI(facts, "x2", x2)
	w10gI(facts, "y2", y2)
	switch o {
	case 1:
		return w10gP(o, fmt.Sprintf("A translation moves (%d, %d) to (%d, %d). What is the squared distance between the points?", x1, y1, x2, y2), strconv.Itoa(dist2), fmt.Sprintf("Add the squared differences: (%d)^2 + (%d)^2 = %d.", dx, dy, dist2), facts)
	case 2:
		area := 3 + rand.Intn(20)
		w10gI(facts, "area", area)
		return w10gP(o, fmt.Sprintf("A rotation by 180 degrees is applied to a figure of area %d. What is the area of the image?", area), strconv.Itoa(area), "A rotation is an isometry, so area is unchanged.", facts)
	default:
		return w10gP(o, fmt.Sprintf("Reflect the point (%d, %d) across the x-axis. What is the y-coordinate of its image?", x1, y1), strconv.Itoa(-y1), fmt.Sprintf("Reflection across the x-axis changes y to -y, so the image has y = %d.", -y1), facts)
	}
}

type w10gHemisphereGen struct{}

func (g *w10gHemisphereGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10GOrd(ctx)
	r := 2 + rand.Intn(14)
	facts := map[string]string{}
	w10gI(facts, "a", r)
	w10gI(facts, "b", r)
	w10gI(facts, "c", r)
	switch o {
	case 1:
		v := 2 * r * r * r
		return w10gP(o, fmt.Sprintf("Use pi = 3. What is the volume of a hemisphere of radius %d?", r), strconv.Itoa(v), fmt.Sprintf("The volume is (2/3) x 3 x %d^3 = %d.", r, v), facts)
	case 2:
		v := 6 * r * r
		return w10gP(o, fmt.Sprintf("Use pi = 3. What is the curved surface area of a hemisphere of radius %d?", r), strconv.Itoa(v), fmt.Sprintf("The curved area is 2 x 3 x %d^2 = %d.", r, v), facts)
	default:
		v := 9 * r * r
		return w10gP(o, fmt.Sprintf("Use pi = 3. What is the total surface area of a hemisphere of radius %d?", r), strconv.Itoa(v), fmt.Sprintf("The total area is 3 x 3 x %d^2 = %d.", r, v), facts)
	}
}

type w10gDensityGen struct{}

func (g *w10gDensityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10GOrd(ctx)
	density := 2 + rand.Intn(9)
	volume := 2 + rand.Intn(19)
	mass := density * volume
	facts := map[string]string{}
	w10gI(facts, "a", mass)
	w10gI(facts, "b", volume)
	w10gI(facts, "c", density)
	switch o {
	case 1:
		return w10gP(o, fmt.Sprintf("A solid has mass %d grams and volume %d cubic centimetres. What is its density in grams per cubic centimetre?", mass, volume), strconv.Itoa(density), fmt.Sprintf("Divide mass by volume: %d/%d = %d.", mass, volume, density), facts)
	case 2:
		return w10gP(o, fmt.Sprintf("A material has density %d grams per cubic centimetre and volume %d cubic centimetres. What is its mass?", density, volume), strconv.Itoa(mass), fmt.Sprintf("Multiply density by volume: %d x %d = %d.", density, volume, mass), facts)
	default:
		return w10gP(o, fmt.Sprintf("A solid has mass %d grams and density %d grams per cubic centimetre. What is its volume?", mass, density), strconv.Itoa(volume), fmt.Sprintf("Divide mass by density: %d/%d = %d.", mass, density, volume), facts)
	}
}
