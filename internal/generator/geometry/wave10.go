package geometry

// Wave 10 geometry: altitudes, the power of a point and the section formula.
// The three families of each concept share fact keys so a single authored
// solution schema resolves for every branch.

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave10(reg *generator.Registry) {
	reg.Register("geo.triangle.altitude", &w10AltitudeGen{})
	reg.Register("geo.circle.power_point", &w10PowerPointGen{})
	reg.Register("geo.coord.section", &w10SectionGen{})
	registerWave10Rest(reg)
}

func w10GP(o int, q, a, ex string, facts map[string]string) generator.Problem {
	if facts == nil {
		facts = map[string]string{}
	}
	facts["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: facts, KPOrdinal: o}
}

func w10GI(m map[string]string, k string, v int) {
	m[k] = strconv.Itoa(v)
}

func w10GOrd(ctx generator.GeneratorContext) int {
	o := ctx.KPOrdinal
	if o < 1 || o > 3 {
		o = 1 + rand.Intn(3)
	}
	return o
}

// ---------------------------------------------------------------------------
// 1. Altitudes of a triangle
// ---------------------------------------------------------------------------

type w10AltitudeGen struct{}

func (g *w10AltitudeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10GOrd(ctx)
	base := 2 * (2 + rand.Intn(9))
	height := 2 + rand.Intn(15)
	area := base * height / 2
	facts := map[string]string{}
	w10GI(facts, "base", base)
	w10GI(facts, "height", height)
	w10GI(facts, "area", area)
	switch o {
	case 1:
		w10GI(facts, "offset", 0)
		w10GI(facts, "distance", base)
		return w10GP(o,
			fmt.Sprintf("A triangle has area %d and base %d. What is its altitude to that base?", area, base),
			strconv.Itoa(height),
			fmt.Sprintf("Use A = bh/2. Then h = 2A/b = 2(%d)/%d = %d.", area, base, height),
			facts)
	case 2:
		w10GI(facts, "offset", 0)
		w10GI(facts, "distance", base)
		return w10GP(o,
			fmt.Sprintf("A triangle has base %d and altitude %d. What is its area?", base, height),
			strconv.Itoa(area),
			fmt.Sprintf("A = bh/2 = (%d)(%d)/2 = %d.", base, height, area),
			facts)
	default:
		offset := 1 + rand.Intn(8)
		distance := base + offset
		w10GI(facts, "offset", offset)
		w10GI(facts, "distance", distance)
		return w10GP(o,
			fmt.Sprintf("An altitude to the line containing a triangle's base is drawn from an obtuse vertex. The base is %d and the foot lies %d units beyond the nearer endpoint. What is the distance from the foot to the farther endpoint?", base, offset),
			strconv.Itoa(distance),
			fmt.Sprintf("The foot lies outside the base, so the distance is the base plus the overhang: %d + %d = %d.", base, offset, distance),
			facts)
	}
}

// ---------------------------------------------------------------------------
// 2. Power of a point
// ---------------------------------------------------------------------------

type w10PowerPointGen struct{}

func (g *w10PowerPointGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10GOrd(ctx)
	facts := map[string]string{}
	switch o {
	case 1:
		for {
			a := 2 + rand.Intn(7)
			b := 3 + rand.Intn(9)
			divisors := make([]int, 0, 4)
			for k := 1; k <= b; k++ {
				if b%k == 0 {
					divisors = append(divisors, k)
				}
			}
			k := divisors[rand.Intn(len(divisors))]
			c := a * k
			d := b / k
			if c == d {
				continue
			}
			product := a * b
			w10GI(facts, "product", product)
			w10GI(facts, "a", a)
			w10GI(facts, "b", b)
			w10GI(facts, "c", c)
			w10GI(facts, "d", d)
			facts["kind"] = "chord"
			return w10GP(o,
				fmt.Sprintf("Two chords intersect inside a circle. One chord is split into segments %d and %d. The other is split into segments %d and x. What is x?", a, b, c),
				strconv.Itoa(d),
				fmt.Sprintf("Intersecting chords give %d x %d = %d x x, so x = %d/%d = %d.", a, b, c, product, c, d),
				facts)
		}
	case 2:
		m := 2 + rand.Intn(6)
		n := m + 1 + rand.Intn(6)
		e := m * m
		total := n * n
		tangent := m * n
		product := e * total
		w10GI(facts, "product", product)
		w10GI(facts, "e", e)
		w10GI(facts, "total", total)
		w10GI(facts, "tangent", tangent)
		facts["kind"] = "tangent"
		return w10GP(o,
			fmt.Sprintf("A tangent and a secant are drawn from the same external point. The secant has external part %d and total length %d. What is the tangent length?", e, total),
			strconv.Itoa(tangent),
			fmt.Sprintf("The tangent length t satisfies t^2 = %d x %d = %d, so t = %d.", e, total, product, tangent),
			facts)
	default:
		for tries := 0; tries < 100; tries++ {
			e1 := 2 + rand.Intn(6)
			total1 := e1 + 2 + rand.Intn(8)
			e2 := 2 + rand.Intn(6)
			product := e1 * total1
			if product%e2 != 0 {
				continue
			}
			total2 := product / e2
			internal2 := total2 - e2
			if internal2 <= 0 || internal2 == total1-e1 {
				continue
			}
			w10GI(facts, "product", product)
			facts["kind"] = "secant"
			w10GI(facts, "e1", e1)
			w10GI(facts, "total1", total1)
			w10GI(facts, "e2", e2)
			w10GI(facts, "internal2", internal2)
			return w10GP(o,
				fmt.Sprintf("Two secants are drawn from the same external point. The first has external part %d and total length %d. The second has external part %d and total length %d. What is the internal part of the second secant?", e1, total1, e2, total2),
				strconv.Itoa(internal2),
				fmt.Sprintf("Equal powers give %d x %d = %d x %d = %d. The second secant's internal part is %d - %d = %d.", e1, total1, e2, total2, product, total2, e2, internal2),
				facts)
		}
		// Fixed fallback with different internal lengths.
		e1, total1, e2, total2 := 3, 9, 2, 12
		product := e1 * total1
		internal2 := total2 - e2
		w10GI(facts, "product", product)
		facts["kind"] = "secant"
		w10GI(facts, "e1", e1)
		w10GI(facts, "total1", total1)
		w10GI(facts, "e2", e2)
		w10GI(facts, "internal2", internal2)
		return w10GP(o,
			fmt.Sprintf("Two secants are drawn from the same external point. The first has external part %d and total length %d. The second has external part %d and total length %d. What is the internal part of the second secant?", e1, total1, e2, total2),
			strconv.Itoa(internal2),
			fmt.Sprintf("Equal powers give %d x %d = %d x %d = %d. The second secant's internal part is %d - %d = %d.", e1, total1, e2, total2, product, total2, e2, internal2),
			facts)
	}
}

// ---------------------------------------------------------------------------
// 3. The section formula
// ---------------------------------------------------------------------------

type w10SectionGen struct{}

func w10GSectionPoint(ax, ay, m, n, dx, dy int) (int, int) {
	return ax + m*dx, ay + m*dy
}

func w10GSectionPair(p, q int) string {
	return fmt.Sprintf("(%d, %d)", p, q)
}

func w10GSectionPartWord(n int) string {
	if n == 1 {
		return "part"
	}
	return "parts"
}

func (g *w10SectionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10GOrd(ctx)
	m, n := 1, 1
	if o == 2 {
		m, n = 1, 1
	} else {
		pairs := [][2]int{{1, 2}, {1, 3}, {2, 1}, {2, 3}, {3, 1}, {3, 2}}
		pair := pairs[rand.Intn(len(pairs))]
		m, n = pair[0], pair[1]
	}
	ax := -8 + rand.Intn(17)
	ay := -8 + rand.Intn(17)
	dx := -4 + rand.Intn(9)
	dy := -4 + rand.Intn(9)
	if dx == 0 && dy == 0 {
		dx = 1
	}
	bx := ax + (m+n)*dx
	by := ay + (m+n)*dy
	px, py := w10GSectionPoint(ax, ay, m, n, dx, dy)
	answerPoint := w10GSectionPair(px, py)
	facts := map[string]string{}
	w10GI(facts, "ax", ax)
	w10GI(facts, "ay", ay)
	w10GI(facts, "bx", bx)
	w10GI(facts, "by", by)
	w10GI(facts, "px", px)
	w10GI(facts, "py", py)
	w10GI(facts, "m", m)
	w10GI(facts, "n", n)
	switch o {
	case 1:
		return w10GP(o,
			fmt.Sprintf("Point P divides segment AB in the ratio %d:%d, with A = %s and B = %s. What are the coordinates of P?", m, n, w10GSectionPair(ax, ay), w10GSectionPair(bx, by)),
			answerPoint,
			fmt.Sprintf("The section formula gives P = ((n x1 + m x2)/(m+n), (n y1 + m y2)/(m+n)) = %s.", answerPoint),
			facts)
	case 2:
		return w10GP(o,
			fmt.Sprintf("Find the midpoint of the segment joining %s and %s.", w10GSectionPair(ax, ay), w10GSectionPair(bx, by)),
			answerPoint,
			fmt.Sprintf("The midpoint is ((%d+%d)/2, (%d+%d)/2) = %s.", ax, bx, ay, by, answerPoint),
			facts)
	default:
		answerPoint = w10GSectionPair(bx, by)
		return w10GP(o,
			fmt.Sprintf("Point P = %s divides segment AB in the ratio %d:%d from A to B. If A = %s, what are the coordinates of B?", w10GSectionPair(px, py), m, n, w10GSectionPair(ax, ay)),
			answerPoint,
			fmt.Sprintf("From the section formula, P - A = %d %s of the vector AB. One part is (%d, %d), so B = P + %d %s = %s.", m, w10GSectionPartWord(m), dx, dy, n, w10GSectionPartWord(n), answerPoint),
			facts)
	}
}
