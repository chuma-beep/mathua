package geometry

import (
	"fmt"
	"math/rand"
	"strconv"

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
	reg.Register("geo.triangle.pythagorean.word", &pythagoreanWordGen{})
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
	registerSimilarity(reg)
	registerTriangles(reg)
	registerAngleFacts(reg)
	registerFinal(reg)
}

type pointsLinesGen struct{}

func (g *pointsLinesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	items := []string{"point", "line", "ray", "line segment"}
	// Predicates, not noun phrases, because the stem supplies the noun. The
	// earlier stem was `What is: %s?` with `has two endpoints`, which asked
	// "What is: has two endpoints?" — a sentence with its subject missing, on a
	// screen where the learner had one field and no way to say "the question,
	// not me, is broken".
	descs := []string{
		"is an exact location with no size",
		"goes forever in two directions",
		"has one endpoint and goes forever",
		"has two endpoints",
	}
	i := rand.Intn(4)
	// The answer is a phrase, so it has to be recognisable rather than typeable:
	// a learner who knows the term was still stuck behind a math field, which
	// dropped the space in "line segment" and compared the wreckage to the
	// authored string. The options carry the phrase intact.
	others := make([]string, 0, len(items)-1)
	for j, it := range items {
		if j != i {
			others = append(others, it)
		}
	}
	p := generator.MultipleChoice(
		fmt.Sprintf("Which geometric term describes something that %s?", descs[i]),
		items[i], others...)
	p.Explanation = fmt.Sprintf("A %s is a figure that %s.", items[i], descs[i])
	p.Facts = map[string]string{"term": items[i], "desc": descs[i]}
	return p
}

type angleTypesGen struct{}

func (g *angleTypesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	if rand.Intn(3) == 0 {
		// production: degrees short of a right angle (unique numeric answer).
		a := rand.Intn(60) + 10 // 10..69
		return generator.Problem{
			Question:    fmt.Sprintf("An angle measures %d degrees. How many degrees short of a right angle is it? (enter a number)", a),
			Answer:      fmt.Sprintf("%d", 90-a),
			Explanation: fmt.Sprintf("A right angle is 90 degrees, so 90 - %d = %d.", a, 90-a),
			Facts:       map[string]string{"given": strconv.Itoa(a), "missing": strconv.Itoa(90 - a)},
		}
	}
	type entry struct {
		deg int
		cls string
	}
	tableEasy := []entry{
		{30, "acute"},
		{45, "acute"},
		{90, "right"},
		{120, "obtuse"},
		{135, "obtuse"},
		{180, "straight"},
	}
	tableHard := []entry{
		{89, "acute"},
		{91, "obtuse"},
		{179, "obtuse"},
		{10, "acute"},
		{150, "obtuse"},
	}
	pool := tableEasy
	if scale > 3 {
		pool = append(tableEasy, tableHard...)
	}
	e := pool[rand.Intn(len(pool))]
	// "obtuse" and "straight angle" cannot be typed into a math field, so the
	// four classifications are offered as options and the phrase stays intact.
	classes := []string{"acute", "right", "obtuse", "straight"}
	others := make([]string, 0, len(classes)-1)
	for _, c := range classes {
		if c != e.cls {
			others = append(others, c)
		}
	}
	p := generator.MultipleChoice(
		fmt.Sprintf("An angle of %d degrees is classified as:", e.deg),
		e.cls, others...)
	p.Explanation = fmt.Sprintf("Acute < 90, right = 90, obtuse > 90 and < 180, straight = 180. %d is %s.", e.deg, e.cls)
	p.Facts = map[string]string{"given": strconv.Itoa(e.deg), "deg": strconv.Itoa(e.deg), "cls": e.cls}
	return p
}

type angleMeasureGen struct{}

func (g *angleMeasureGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	if rand.Intn(3) == 0 {
		// production: degrees to reach a straight angle (unique numeric answer).
		a := rand.Intn(80) + 91 // 91..170
		return generator.Problem{
			Question:    fmt.Sprintf("An angle measures %d degrees. How many more degrees are needed to reach a straight angle? (enter a number)", a),
			Answer:      fmt.Sprintf("%d", 180-a),
			Explanation: fmt.Sprintf("A straight angle is 180 degrees, so 180 - %d = %d.", a, 180-a),
			Facts:       map[string]string{"given": strconv.Itoa(a), "missing": strconv.Itoa(180 - a)},
		}
	}
	type entry struct {
		deg int
		cls string
	}
	tableEasy := []entry{
		{25, "acute angle"},
		{60, "acute angle"},
		{90, "right angle"},
		{110, "obtuse angle"},
		{150, "obtuse angle"},
		{180, "straight angle"},
	}
	tableHard := []entry{
		{89, "acute angle"},
		{91, "obtuse angle"},
		{179, "obtuse angle"},
		{5, "acute angle"},
		{160, "obtuse angle"},
	}
	pool := tableEasy
	if scale > 3 {
		pool = append(tableEasy, tableHard...)
	}
	e := pool[rand.Intn(len(pool))]
	// The answer is a phrase, and a phrase cannot survive the math field: report
	// #200 showed "line segment" arriving as one word. Offer the classes instead.
	classes := []string{"acute angle", "right angle", "obtuse angle", "straight angle"}
	distractors := make([]string, 0, len(classes)-1)
	for _, c := range classes {
		if c != e.cls {
			distractors = append(distractors, c)
		}
	}
	p := generator.MultipleChoice(fmt.Sprintf("An angle of %d degrees is a:", e.deg), e.cls, distractors...)
	p.Explanation = fmt.Sprintf("%d degrees is %s.", e.deg, e.cls)
	p.Facts = map[string]string{"given": strconv.Itoa(e.deg), "deg": strconv.Itoa(e.deg), "cls": e.cls}
	return p
}

type complementaryGen struct{}

func (g *complementaryGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := rand.Intn(max(1, scale*15)) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("Two angles are complementary. One is %d degrees. What is the other?", a),
		Answer:      fmt.Sprintf("%d", 90-a),
		Explanation: fmt.Sprintf("Complementary angles sum to 90. 90 - %d = %d.", a, 90-a),
		Facts:       map[string]string{"given": strconv.Itoa(a), "missing": strconv.Itoa(90 - a)},
	}
}

type verticalAnglesGen struct{}

func (g *verticalAnglesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// An angle adjacent to angle 1 has measure d degrees. Find the vertical angle to angle 1.
	// This tests supplementary adjacent angles and vertical angle equality without reading
	// off the answer from the stem.
	d := 20 + rand.Intn(141) // 20 to 160
	if d == 90 {
		d = 85
	}
	ans := 180 - d
	return generator.Problem{
		Question:    fmt.Sprintf("Two lines intersect. An angle adjacent to \\(\\angle 1\\) measures %d degrees. What is the measure of the vertical angle to \\(\\angle 1\\) (in degrees)?", d),
		Answer:      fmt.Sprintf("%d", ans),
		Explanation: fmt.Sprintf("Adjacent angles along a straight line are supplementary (sum to 180 degrees), so \\(\\angle 1 = 180 - %d = %d^\\circ\\). Vertical angles are equal, so the vertical angle also measures %d degrees.", d, ans, ans),
		Facts: map[string]string{
			"adjacent": strconv.Itoa(d),
			"a":        strconv.Itoa(ans),
			"vertical": strconv.Itoa(ans),
		},
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
	// The three kinds are names, so they are offered rather than typed.
	p := generator.MultipleChoice(
		fmt.Sprintf("A triangle has sides %d, %d, %d. By side length, it is:", sides[0], sides[1], sides[2]),
		ans, "scalene", "isosceles", "equilateral",
	)
	p.Explanation = "All 3 equal = equilateral. Two equal = isosceles. None equal = scalene."
	p.Facts = map[string]string{"s1": strconv.Itoa(sides[0]), "s2": strconv.Itoa(sides[1]), "s3": strconv.Itoa(sides[2]), "kind": ans}
	return p
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
		Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "c": strconv.Itoa(c), "sum": "180"},
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
		Facts:       map[string]string{"base": strconv.Itoa(base), "height": strconv.Itoa(height), "product": strconv.Itoa(base * height), "area": strconv.Itoa(area)},
	}
}

type pythagoreanGen struct{}

func (g *pythagoreanGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	type triple struct{ a, b, c int }
	tableEasy := []triple{
		{3, 4, 5},
		{6, 8, 10},
		{5, 12, 13},
		{9, 12, 15},
	}
	tableHard := []triple{
		{8, 15, 17},
		{7, 24, 25},
		{20, 21, 29},
	}
	if rand.Intn(3) == 0 {
		// production: missing leg from hypotenuse and one leg (unique answer).
		pool := tableEasy
		if scale > 3 {
			pool = append(tableEasy, tableHard...)
		}
		t := pool[rand.Intn(len(pool))]
		known, unknown := t.a, t.b
		if rand.Intn(2) == 0 {
			known, unknown = t.b, t.a
		}
		return generator.Problem{
			Question:    fmt.Sprintf("Right triangle: hypotenuse \\(= %d\\) and one leg \\(= %d\\). Find the other leg. (enter a number)", t.c, known),
			Answer:      fmt.Sprintf("%d", unknown),
			Explanation: fmt.Sprintf("\\(\\text{leg}^{2} = %d^{2} - %d^{2} = %d - %d = %d\\), so \\(\\text{leg} = \\sqrt{%d} = %d\\).", t.c, known, t.c*t.c, known*known, unknown*unknown, unknown*unknown, unknown),
			Facts:       map[string]string{"a": strconv.Itoa(t.a), "b": strconv.Itoa(t.b), "c": strconv.Itoa(t.c), "a_sq": strconv.Itoa(t.a * t.a), "b_sq": strconv.Itoa(t.b * t.b), "leg_sq": strconv.Itoa(t.a*t.a - t.b*t.b), "leg": strconv.Itoa(unknown), "c_sq": strconv.Itoa(t.c * t.c)},
		}
	}
	pool := tableEasy
	if scale > 3 {
		pool = append(tableEasy, tableHard...)
	}
	t := pool[rand.Intn(len(pool))]
	return generator.Problem{
		Question:    fmt.Sprintf("Right triangle: legs \\(= %d\\) and %d. Find the hypotenuse \\(c\\).", t.a, t.b),
		Answer:      fmt.Sprintf("%d", t.c),
		Explanation: fmt.Sprintf("\\(c^{2} = %d^{2} + %d^{2} = %d + %d = %d\\), so \\(c = \\sqrt{%d} = %d\\).", t.a, t.b, t.a*t.a, t.b*t.b, t.c*t.c, t.c*t.c, t.c),
		Facts:       map[string]string{"a": strconv.Itoa(t.a), "b": strconv.Itoa(t.b), "c": strconv.Itoa(t.c), "a_sq": strconv.Itoa(t.a * t.a), "b_sq": strconv.Itoa(t.b * t.b), "c_sq": strconv.Itoa(t.c * t.c)},
	}
}

// pythagoreanWordGen: the theorem as a story. The numbers come from the same
// triples the parent uses, so every square root is exact and the answer is an
// integer. The two flavours -- missing hypotenuse, missing leg -- publish the
// same facts, so one authored solution schema covers both.
type pythagoreanWordGen struct{}

func (g *pythagoreanWordGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type triple struct{ a, b, c int }
	triples := []triple{{3, 4, 5}, {6, 8, 10}, {5, 12, 13}, {9, 12, 15}}
	if ctx.Difficulty >= 0.45 {
		triples = append(triples, triple{8, 15, 17}, triple{7, 24, 25})
	}
	if ctx.Difficulty >= 0.7 {
		triples = append(triples, triple{20, 21, 29}, triple{12, 16, 20})
	}
	t := triples[rand.Intn(len(triples))]
	leg, other, hyp := t.a, t.b, t.c
	if rand.Intn(2) == 0 {
		// The ladder's foot and the height it reaches are the two legs.
		return generator.Problem{
			Question: fmt.Sprintf("A ladder leans against a wall. Its foot is %d feet from the wall and it reaches %d feet up the wall. How long is the ladder?",
				leg, other),
			Answer: strconv.Itoa(hyp),
			Explanation: fmt.Sprintf("\\(%d^2 + %d^2 = %d + %d = %d\\), so the ladder is \\(\\sqrt{%d} = %d\\) feet long.",
				leg, other, leg*leg, other*other, hyp*hyp, hyp*hyp, hyp),
			Facts: map[string]string{
				"story":   "A ladder leans against a wall.",
				"rule":    "the square of the hypotenuse is the sum of the squares of the legs",
				"known":   fmt.Sprintf("the two legs, %d feet and %d feet", leg, other),
				"side":    "ladder (the hypotenuse)",
				"unknown": strconv.Itoa(hyp),
			},
		}
	}
	// The ladder's length is the hypotenuse and its foot distance is one leg.
	return generator.Problem{
		Question: fmt.Sprintf("A ladder %d feet long leans against a wall. Its foot is %d feet from the wall. How far up the wall does it reach?",
			hyp, leg),
		Answer: strconv.Itoa(other),
		Explanation: fmt.Sprintf("\\(%d^2 - %d^2 = %d - %d = %d\\), so the height is \\(\\sqrt{%d} = %d\\) feet.",
			hyp, leg, hyp*hyp, leg*leg, other*other, other*other, other),
		Facts: map[string]string{
			"story":   "A ladder leans against a wall.",
			"rule":    "a leg is the square root of the hypotenuse squared minus the other leg squared",
			"known":   fmt.Sprintf("the hypotenuse, %d feet, and one leg, %d feet", hyp, leg),
			"side":    "height reached on the wall (a leg)",
			"unknown": strconv.Itoa(other),
		},
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
	others := make([]string, 0, len(types)-1)
	for _, t := range types {
		if t.name != types[i].name {
			others = append(others, t.name)
		}
	}
	p := generator.MultipleChoice(fmt.Sprintf("A quadrilateral with %s is a:", types[i].desc), types[i].name, others...)
	p.Explanation = types[i].desc + " describes a " + types[i].name + "."
	p.Facts = map[string]string{"name": types[i].name, "desc": types[i].desc}
	return p
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
		Facts:       map[string]string{"shape": shape, "base": strconv.Itoa(base), "height": strconv.Itoa(height), "area": strconv.Itoa(base * height)},
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
		Facts:       map[string]string{"w": strconv.Itoa(w), "h": strconv.Itoa(h), "sum": strconv.Itoa(w + h), "perimeter": strconv.Itoa(2 * (w + h))},
	}
}

type circlePartsGen struct{}

func (g *circlePartsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	parts := []string{"radius", "diameter", "circumference", "chord"}
	descs := []string{"distance from center to edge", "distance across the circle through center", "distance around the circle", "a line from one edge to another (not through center)"}
	i := rand.Intn(4)
	others := make([]string, 0, len(parts)-1)
	for _, part := range parts {
		if part != parts[i] {
			others = append(others, part)
		}
	}
	p := generator.MultipleChoice(fmt.Sprintf("The %s is the:", descs[i]), parts[i], others...)
	p.Explanation = descs[i] + " = " + parts[i] + "."
	p.Facts = map[string]string{"term": parts[i], "desc": descs[i]}
	return p
}

type circumferenceGen struct{}

func (g *circumferenceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	r := rand.Intn(max(1, scale*4)) + 5
	pi := 3.14
	// Two decimal places, because the learner computes 2 x 3.14 x r exactly and
	// rounding the answer to a whole number rejected 69.08 for r = 11 (report
	// 213, 23 October 2026).
	circ := 2 * pi * float64(r)
	return generator.Problem{
		Question:    fmt.Sprintf("Circle radius = %d. Find the circumference to two decimal places (use pi = 3.14).", r),
		Answer:      fmt.Sprintf("%.2f", circ),
		Explanation: fmt.Sprintf("\\(C = 2 \\pi r = 2 \\times 3.14 \\times %d = %.2f\\).", r, circ),
		Facts:       map[string]string{"r": strconv.Itoa(r), "circ": fmt.Sprintf("%.2f", circ), "twice_r": strconv.Itoa(2 * r)},
	}
}

type circleAreaGen struct{}

func (g *circleAreaGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	r := rand.Intn(max(1, scale*2)) + 3
	// Same precision rule as the circumference: the learner computes
	// 3.14 x r x r exactly, so the answer carries two decimal places.
	area := 3.14 * float64(r*r)
	return generator.Problem{
		Question:    fmt.Sprintf("Circle radius = %d. Find the area to two decimal places (use pi = 3.14).", r),
		Answer:      fmt.Sprintf("%.2f", area),
		Explanation: fmt.Sprintf("\\(A = \\pi r^{2} = 3.14 \\times %d^{2} = 3.14 \\times %d = %.2f\\).", r, r*r, area),
		Facts:       map[string]string{"r": strconv.Itoa(r), "r_sq": strconv.Itoa(r * r), "area": fmt.Sprintf("%.2f", area)},
	}
}

type coordPlotGen struct{}

func (g *coordPlotGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	if rand.Intn(3) == 0 {
		// production: read one coordinate off a plotted point (unique answer).
		x := rand.Intn(11) - 5
		y := rand.Intn(11) - 5
		if rand.Intn(2) == 0 {
			return generator.Problem{
				Question:    fmt.Sprintf("What is the x-coordinate of the point \\((%d, %d)\\)? (enter a number)", x, y),
				Answer:      fmt.Sprintf("%d", x),
				Explanation: fmt.Sprintf("In (%d,%d) the first entry is the x-coordinate: %d.", x, y, x),
				Facts:       map[string]string{"x": strconv.Itoa(x), "y": strconv.Itoa(y), "which": "x"},
			}
		}
		return generator.Problem{
			Question:    fmt.Sprintf("What is the y-coordinate of the point \\((%d, %d)\\)? (enter a number)", x, y),
			Answer:      fmt.Sprintf("%d", y),
			Explanation: fmt.Sprintf("In (%d,%d) the second entry is the y-coordinate: %d.", x, y, y),
			Facts:       map[string]string{"x": strconv.Itoa(x), "y": strconv.Itoa(y), "which": "y"},
		}
	}
	type entry struct{ x, y int }
	tableEasy := []entry{
		{2, 3},
		{-2, 3},
		{-2, -3},
		{2, -3},
		{0, 4},
		{3, 0},
	}
	tableHard := []entry{
		{5, 7},
		{-6, 2},
		{-4, -8},
		{7, -1},
		{0, -5},
		{-3, 0},
	}
	pool := tableEasy
	if scale > 3 {
		pool = append(tableEasy, tableHard...)
	}
	e := pool[rand.Intn(len(pool))]
	// A quadrant is a name, not something a learner can type in a math field.
	p := generator.MultipleChoice(
		fmt.Sprintf("What quadrant is the point \\((%d, %d)\\) in?", e.x, e.y),
		quad(e.x, e.y),
		"Quadrant I", "Quadrant II", "Quadrant III", "Quadrant IV", "on an axis",
	)
	p.Explanation = fmt.Sprintf("(%d,%d) is in %s.", e.x, e.y, quad(e.x, e.y))
	p.Facts = map[string]string{"x": strconv.Itoa(e.x), "y": strconv.Itoa(e.y), "quadrant": quad(e.x, e.y)}
	return p
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
	type triple struct{ dx, dy, dist int }
	tableEasy := []triple{
		{3, 4, 5},
		{6, 8, 10},
		{5, 12, 13},
		{9, 12, 15},
	}
	tableHard := []triple{
		{8, 15, 17},
		{7, 24, 25},
		{20, 21, 29},
	}
	if rand.Intn(3) == 0 {
		// production: axis-aligned distance (unique numeric answer).
		x1 := rand.Intn(9) - 4
		y := rand.Intn(9) - 4
		dx := rand.Intn(8) + 2 // 2..9
		return generator.Problem{
			Question:    fmt.Sprintf("Find the distance between \\((%d,%d)\\) and \\((%d,%d)\\). (enter a number)", x1, y, x1+dx, y),
			Answer:      fmt.Sprintf("%d", dx),
			Explanation: fmt.Sprintf("Same y-coordinate, so the distance is \\(|%d-%d| = %d\\).", x1+dx, x1, dx),
			Facts:       map[string]string{"x1": strconv.Itoa(x1), "x2": strconv.Itoa(x1 + dx), "y": strconv.Itoa(y), "dist": strconv.Itoa(dx), "dx": strconv.Itoa(dx), "dy": "0", "dx_sq": strconv.Itoa(dx * dx), "dy_sq": "0", "sum_sq": strconv.Itoa(dx * dx)},
		}
	}
	pool := tableEasy
	if scale > 3 {
		pool = append(tableEasy, tableHard...)
	}
	t := pool[rand.Intn(len(pool))]
	x1 := rand.Intn(5)
	y1 := rand.Intn(5)
	return generator.Problem{
		Question:    fmt.Sprintf("Find the distance between \\((%d,%d)\\) and \\((%d,%d)\\).", x1, y1, x1+t.dx, y1+t.dy),
		Answer:      fmt.Sprintf("%d", t.dist),
		Explanation: fmt.Sprintf("\\(\\sqrt{(%d-%d)^{2} + (%d-%d)^{2}} = \\sqrt{%d + %d} = \\sqrt{%d} = %d\\).", x1+t.dx, x1, y1+t.dy, y1, t.dx*t.dx, t.dy*t.dy, t.dx*t.dx+t.dy*t.dy, t.dist),
		Facts:       map[string]string{"x1": strconv.Itoa(x1), "y1": strconv.Itoa(y1), "dx": strconv.Itoa(t.dx), "dy": strconv.Itoa(t.dy), "dx_sq": strconv.Itoa(t.dx * t.dx), "dy_sq": strconv.Itoa(t.dy * t.dy), "sum_sq": strconv.Itoa(t.dx*t.dx + t.dy*t.dy), "dist": strconv.Itoa(t.dist)},
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
		Facts:       map[string]string{"x1": strconv.Itoa(x1), "y1": strconv.Itoa(y1), "x2": strconv.Itoa(x2), "y2": strconv.Itoa(y2), "mid_x": fmt.Sprintf("%.1f", float64(x1+x2)/2), "mid_y": fmt.Sprintf("%.1f", float64(y1+y2)/2)},
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
		Facts:       map[string]string{"l": strconv.Itoa(l), "w": strconv.Itoa(w), "h": strconv.Itoa(h), "volume": strconv.Itoa(l * w * h)},
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
		Facts:       map[string]string{"l": strconv.Itoa(l), "w": strconv.Itoa(w), "h": strconv.Itoa(h), "lw": strconv.Itoa(l * w), "lh": strconv.Itoa(l * h), "wh": strconv.Itoa(w * h), "sum": strconv.Itoa(l*w + l*h + w*h), "sa": strconv.Itoa(sa)},
	}
}

type coordLinesGen struct{}

func (g *coordLinesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	if rand.Intn(3) == 0 {
		// production: slope of a line through two points (unique numeric answer).
		run := rand.Intn(4) + 1
		slope := rand.Intn(4) + 1
		rise := run * slope
		x1 := rand.Intn(7) - 3
		y1 := rand.Intn(7) - 3
		return generator.Problem{
			Question:    fmt.Sprintf("What is the slope of the line through \\((%d,%d)\\) and \\((%d,%d)\\)? (enter a number)", x1, y1, x1+run, y1+rise),
			Answer:      fmt.Sprintf("%d", slope),
			Explanation: fmt.Sprintf("\\(\\text{Slope} = \\frac{%d-%d}{%d-%d} = \\frac{%d}{%d} = %d\\).", y1+rise, y1, x1+run, x1, rise, run, slope),
			Facts:       map[string]string{"x1": strconv.Itoa(x1), "y1": strconv.Itoa(y1), "x2": strconv.Itoa(x1 + run), "y2": strconv.Itoa(y1 + rise), "rise": strconv.Itoa(rise), "run": strconv.Itoa(run), "slope": strconv.Itoa(slope), "dx": strconv.Itoa(run), "dy": strconv.Itoa(rise)},
		}
	}
	x1 := rand.Intn(scale*2+4) - 3
	y1 := rand.Intn(scale*2+4) - 3
	x2 := x1 + rand.Intn(scale*2+2) + 1
	y2 := y1 + rand.Intn(scale*2+2) + 1
	dx := x2 - x1
	dy := y2 - y1
	if scale <= 3 {
		// easy: integer slopes only; resample the rise until it divides evenly.
		for dy%dx != 0 {
			y2 = y1 + rand.Intn(scale*2+2) + 1
			dy = y2 - y1
		}
	}
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
		Facts:       map[string]string{"x1": strconv.Itoa(x1), "y1": strconv.Itoa(y1), "x2": strconv.Itoa(x2), "y2": strconv.Itoa(y2), "dx": strconv.Itoa(dx), "dy": strconv.Itoa(dy), "slope": answer},
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
			Facts:       map[string]string{"r": strconv.Itoa(e.r), "theta": strconv.Itoa(e.theta), "x": strconv.Itoa(e.x), "y": strconv.Itoa(e.y)},
		}
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Convert \\((%d,%d)\\) to polar coordinates \\((r > 0, 0 \\leq \\theta < 360)\\).", e.x, e.y),
		Answer:      fmt.Sprintf("(%d,%d°)", e.r, e.theta),
		Explanation: fmt.Sprintf("\\(r = \\sqrt{%d^{2}+%d^{2}} = %d\\), \\(\\theta = \\arctan\\left(\\frac{%d}{%d}\\right) = %d^{\\circ} \\to (%d,%d^{\\circ})\\)", e.x, e.y, e.r, e.y, e.x, e.theta, e.r, e.theta),
		Facts:       map[string]string{"r": strconv.Itoa(e.r), "theta": strconv.Itoa(e.theta), "x": strconv.Itoa(e.x), "y": strconv.Itoa(e.y)},
	}
}
