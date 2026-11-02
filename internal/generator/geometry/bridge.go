package geometry

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

type geoBridgeGen struct{ kind string }

func geoOrdinal(ctx generator.GeneratorContext) int {
	if ctx.KPOrdinal >= 1 && ctx.KPOrdinal <= 3 {
		return ctx.KPOrdinal
	}
	return 1 + rand.Intn(3)
}

func geoProb(o int, q, a, e string) generator.Problem {
	return generator.Problem{
		Question:    q,
		Answer:      a,
		Explanation: e,
		Facts:       map[string]string{"answer": a},
		KPOrdinal:   o,
	}
}

func (g *geoBridgeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := geoOrdinal(ctx)
	a := 20 + rand.Intn(131)
	b := 2 + rand.Intn(10)
	x := rand.Intn(9) - 4
	y := rand.Intn(9) - 4
	switch g.kind {
	case "transversal":
		switch o {
		case 1:
			return geoProb(o, fmt.Sprintf("Two parallel lines are cut by a transversal. One interior angle measures %d degrees. What is the measure of its alternate interior angle? (enter a number)", a), strconv.Itoa(a), fmt.Sprintf("Alternate interior angles are equal when the lines are parallel, so the angle is %d degrees.", a))
		case 2:
			ans := 180 - a
			return geoProb(o, fmt.Sprintf("Two parallel lines are cut by a transversal. One interior angle measures %d degrees. What is the same-side interior angle? (enter a number)", a), strconv.Itoa(ans), fmt.Sprintf("Same-side interior angles are supplementary: 180 - %d = %d.", a, ans))
		default:
			return geoProb(o, fmt.Sprintf("Two parallel lines are cut by a transversal. One angle measures %d degrees. What is its corresponding angle? (enter a number)", a), strconv.Itoa(a), fmt.Sprintf("Corresponding angles are equal, so the angle is %d degrees.", a))
		}
	case "congruence":
		switch o {
		case 1:
			p := 20 + rand.Intn(60)
			q := 30 + rand.Intn(50)
			ans := 180 - p - q
			return geoProb(o, fmt.Sprintf("A triangle has angles %d degrees and %d degrees. What is the third angle? (enter a number)", p, q), strconv.Itoa(ans), fmt.Sprintf("The angles sum to 180: 180 - %d - %d = %d.", p, q, ans))
		case 2:
			return geoProb(o, fmt.Sprintf("Two triangles are congruent. One has a side of length %d. What is the corresponding side of the other? (enter a number)", b), strconv.Itoa(b), "Corresponding sides of congruent triangles are equal.")
		default:
			ans := 3*b + 2
			return geoProb(o, fmt.Sprintf("A triangle has sides %d, %d and %d. A congruent triangle is formed. What is its perimeter? (enter a number)", b, b+1, b+1), strconv.Itoa(ans), fmt.Sprintf("Congruent triangles have equal perimeters: %d + %d + %d = %d.", b, b+1, b+1, ans))
		}
	case "pyth_converse":
		triples := [][3]int{{3, 4, 5}, {5, 12, 13}, {6, 8, 10}, {8, 15, 17}, {7, 24, 25}, {9, 12, 15}}
		switch o {
		case 1:
			t := triples[rand.Intn(len(triples))]
			return geoProb(o, fmt.Sprintf("A triangle has sides %d, %d and %d. Is it a right triangle? (enter 1 for yes and 0 for no)", t[0], t[1], t[2]), "1", fmt.Sprintf("The largest square %d^2 equals %d^2 + %d^2, so the triangle is right.", t[2], t[0], t[1]))
		case 2:
			p := 2 + rand.Intn(6)
			q := p + 1 + rand.Intn(3)
			r := q + 1 + rand.Intn(3)
			code := 3
			if r*r == p*p+q*q {
				code = 1
			} else if r*r < p*p+q*q {
				code = 2
			}
			return geoProb(o, fmt.Sprintf("A triangle has sides %d, %d and %d. Classify it: enter 1 for right, 2 for acute and 3 for obtuse.", p, q, r), strconv.Itoa(code), fmt.Sprintf("Compare %d^2 with %d^2 + %d^2. The code is %d.", r, p, q, code))
		default:
			t := triples[rand.Intn(len(triples))]
			return geoProb(o, fmt.Sprintf("A right triangle has legs %d and %d. What is the hypotenuse? (enter a number)", t[0], t[1]), strconv.Itoa(t[2]), fmt.Sprintf("The hypotenuse is %d because %d^2 + %d^2 = %d^2.", t[2], t[0], t[1], t[2]))
		}
	case "special45":
		switch o {
		case 1:
			k := 1 + rand.Intn(9)
			return geoProb(o, fmt.Sprintf("A 45-45-90 triangle has legs of length %d. What is the square of its hypotenuse? (enter a number)", k), strconv.Itoa(2*k*k), fmt.Sprintf("The hypotenuse is %d root 2, so its square is 2 x %d^2 = %d.", k, k, 2*k*k))
		case 2:
			k := 2 * (1 + rand.Intn(5))
			return geoProb(o, fmt.Sprintf("A 45-45-90 triangle has legs of length %d. What is its area? (enter a number)", k), strconv.Itoa(k*k/2), fmt.Sprintf("The area is half the product of the legs: %d x %d / 2 = %d.", k, k, k*k/2))
		default:
			k := 1 + rand.Intn(9)
			return geoProb(o, fmt.Sprintf("A 45-45-90 triangle has hypotenuse squared equal to %d. What is the square of each leg? (enter a number)", 2*k*k), strconv.Itoa(k*k), fmt.Sprintf("Each leg squared is half the hypotenuse squared: %d / 2 = %d.", 2*k*k, k*k))
		}
	case "centroid":
		switch o {
		case 1:
			m := 3 * (1 + rand.Intn(6))
			return geoProb(o, fmt.Sprintf("A median has length %d. The centroid divides it in the ratio 2 to 1 from the vertex. How far is the vertex from the centroid? (enter a number)", m), strconv.Itoa(2*m/3), fmt.Sprintf("Two thirds of %d is %d.", m, 2*m/3))
		case 2:
			p := rand.Intn(9) - 4
			q := rand.Intn(9) - 4
			k := rand.Intn(7) - 3
			r := 3*k - p - q
			return geoProb(o, fmt.Sprintf("A triangle has vertices with x-coordinates %d, %d and %d. What is the x-coordinate of its centroid? (enter a number)", p, q, r), strconv.Itoa(k), fmt.Sprintf("The centroid x-coordinate is the average: (%d + %d + %d)/3 = %d.", p, q, r, k))
		default:
			p := rand.Intn(7) - 3
			q := rand.Intn(7) - 3
			cy := rand.Intn(7) - 3
			ans := 3*cy - p - q
			return geoProb(o, fmt.Sprintf("A triangle has two vertices with y-coordinates %d and %d. Its centroid has y-coordinate %d. What is the y-coordinate of the third vertex? (enter a number)", p, q, cy), strconv.Itoa(ans), fmt.Sprintf("The centroid is the average: 3 x %d = %d + %d + y, so y = %d.", cy, p, q, ans))
		}
	case "translation":
		switch o {
		case 1:
			dx := rand.Intn(9) - 4
			ans := x + dx
			return geoProb(o, fmt.Sprintf("A point (%d, %d) is translated by (%d, %d). What is its new x-coordinate? (enter a number)", x, y, dx, 2), strconv.Itoa(ans), fmt.Sprintf("Add the horizontal component: %d + %d = %d.", x, dx, ans))
		case 2:
			dx := x + 5
			return geoProb(o, fmt.Sprintf("A point (%d, %d) moves to (%d, %d) under a translation. What is the horizontal component of the translation vector? (enter a number)", x, y, x+dx, y), strconv.Itoa(dx), fmt.Sprintf("Subtract the coordinates: %d - %d = %d.", x+dx, x, dx))
		default:
			dy := rand.Intn(9) - 4
			ans := y + dy
			return geoProb(o, fmt.Sprintf("A point (%d, %d) is translated by (%d, %d). What is its new y-coordinate? (enter a number)", x, y, 2, dy), strconv.Itoa(ans), fmt.Sprintf("Add the vertical component: %d + %d = %d.", y, dy, ans))
		}
	case "reflection":
		switch o {
		case 1:
			ans := -y
			return geoProb(o, fmt.Sprintf("A point (%d, %d) is reflected across the x-axis. What is its new y-coordinate? (enter a number)", x, y), strconv.Itoa(ans), fmt.Sprintf("Reflecting across the x-axis changes y to -y: -(%d) = %d.", y, ans))
		case 2:
			ans := -x
			return geoProb(o, fmt.Sprintf("A point (%d, %d) is reflected across the y-axis. What is its new x-coordinate? (enter a number)", x, y), strconv.Itoa(ans), fmt.Sprintf("Reflecting across the y-axis changes x to -x: -(%d) = %d.", x, ans))
		default:
			ans := y
			return geoProb(o, fmt.Sprintf("A point (%d, %d) is reflected across the line y = x. What is its new x-coordinate? (enter a number)", x, y), strconv.Itoa(ans), fmt.Sprintf("Reflecting across y = x swaps the coordinates, so the new x-coordinate is the old y-coordinate %d.", y))
		}
	case "rotation":
		switch o {
		case 1:
			ans := -y
			return geoProb(o, fmt.Sprintf("A point (%d, %d) is rotated 90 degrees counterclockwise about the origin. What is its new x-coordinate? (enter a number)", x, y), strconv.Itoa(ans), fmt.Sprintf("A 90 degree counterclockwise turn sends (x, y) to (-y, x), so the new x-coordinate is %d.", ans))
		case 2:
			ans := -x
			return geoProb(o, fmt.Sprintf("A point (%d, %d) is rotated 180 degrees about the origin. What is its new x-coordinate? (enter a number)", x, y), strconv.Itoa(ans), fmt.Sprintf("A half turn sends (x, y) to (-x, -y), so the new x-coordinate is %d.", ans))
		default:
			ans := y
			return geoProb(o, fmt.Sprintf("A point (%d, %d) is rotated 270 degrees counterclockwise about the origin. What is its new x-coordinate? (enter a number)", x, y), strconv.Itoa(ans), fmt.Sprintf("A 270 degree counterclockwise turn sends (x, y) to (y, -x), so the new x-coordinate is %d.", ans))
		}
	case "cylinder":
		switch o {
		case 1:
			r := 1 + rand.Intn(7)
			h := 2 + rand.Intn(9)
			ans := r * r * h
			return geoProb(o, fmt.Sprintf("A cylinder has radius %d and height %d. What is its volume divided by pi? (enter a number)", r, h), strconv.Itoa(ans), fmt.Sprintf("The volume is pi r squared h, so the coefficient is %d^2 x %d = %d.", r, h, ans))
		case 2:
			r := 1 + rand.Intn(7)
			h := 2 + rand.Intn(9)
			ans := 2 * r * h
			return geoProb(o, fmt.Sprintf("A cylinder has radius %d and height %d. What is its lateral surface area divided by pi? (enter a number)", r, h), strconv.Itoa(ans), fmt.Sprintf("Lateral area is 2 pi r h, so the coefficient is 2 x %d x %d = %d.", r, h, ans))
		default:
			r := 1 + rand.Intn(6)
			h := 2 + rand.Intn(9)
			ans := 2 * r * (r + h)
			return geoProb(o, fmt.Sprintf("A cylinder has radius %d and height %d. What is its total surface area divided by pi? (enter a number)", r, h), strconv.Itoa(ans), fmt.Sprintf("Total surface is 2 pi r (r + h), so the coefficient is 2 x %d x %d = %d.", r, r+h, ans))
		}
	case "cone":
		switch o {
		case 1:
			r := 3 * (1 + rand.Intn(3))
			h := 2 + rand.Intn(8)
			ans := r * r * h / 3
			return geoProb(o, fmt.Sprintf("A cone has radius %d and height %d. What is its volume divided by pi? (enter a number)", r, h), strconv.Itoa(ans), fmt.Sprintf("Volume is one third pi r squared h, so the coefficient is %d^2 x %d / 3 = %d.", r, h, ans))
		case 2:
			triples := [][3]int{{3, 4, 5}, {6, 8, 10}, {5, 12, 13}, {8, 15, 17}, {9, 12, 15}}
			t := triples[rand.Intn(len(triples))]
			return geoProb(o, fmt.Sprintf("A right cone has radius %d and height %d. What is its slant height? (enter a number)", t[0], t[1]), strconv.Itoa(t[2]), fmt.Sprintf("The slant height satisfies %d^2 + %d^2 = %d^2.", t[0], t[1], t[2]))
		default:
			r := 1 + rand.Intn(7)
			l := 2 + rand.Intn(9)
			ans := r * l
			return geoProb(o, fmt.Sprintf("A cone has radius %d and slant height %d. What is its lateral surface area divided by pi? (enter a number)", r, l), strconv.Itoa(ans), fmt.Sprintf("Lateral area is pi r l, so the coefficient is %d x %d = %d.", r, l, ans))
		}
	case "sphere":
		switch o {
		case 1:
			k := 1 + rand.Intn(6)
			r := 3 * k
			ans := 36 * k * k * k
			return geoProb(o, fmt.Sprintf("A sphere has radius %d. What is its volume divided by pi? (enter a number)", r), strconv.Itoa(ans), fmt.Sprintf("Volume is four thirds pi r cubed, so the coefficient is 4 x %d^3 / 3 = %d.", r, ans))
		case 2:
			r := 1 + rand.Intn(9)
			ans := 4 * r * r
			return geoProb(o, fmt.Sprintf("A sphere has radius %d. What is its surface area divided by pi? (enter a number)", r), strconv.Itoa(ans), fmt.Sprintf("Surface area is 4 pi r squared, so the coefficient is 4 x %d^2 = %d.", r, ans))
		default:
			d := 2 + rand.Intn(15)
			return geoProb(o, fmt.Sprintf("A sphere has diameter %d. What is its radius? (enter a number)", d), strconv.Itoa(d/2), fmt.Sprintf("The radius is half the diameter: %d / 2 = %d.", d, d/2))
		}
	case "curved":
		switch o {
		case 1:
			r := 1 + rand.Intn(7)
			h := 2 + rand.Intn(9)
			ans := 2 * r * h
			return geoProb(o, fmt.Sprintf("A cylinder has radius %d and height %d. What is its curved surface area divided by pi? (enter a number)", r, h), strconv.Itoa(ans), fmt.Sprintf("The curved area is 2 pi r h, so the coefficient is 2 x %d x %d = %d.", r, h, ans))
		case 2:
			r := 1 + rand.Intn(7)
			l := 2 + rand.Intn(9)
			ans := r * l
			return geoProb(o, fmt.Sprintf("A cone has radius %d and slant height %d. What is its curved surface area divided by pi? (enter a number)", r, l), strconv.Itoa(ans), fmt.Sprintf("The curved area is pi r l, so the coefficient is %d x %d = %d.", r, l, ans))
		default:
			r := 1 + rand.Intn(9)
			ans := 4 * r * r
			return geoProb(o, fmt.Sprintf("A sphere has radius %d. What is its curved surface area divided by pi? (enter a number)", r), strconv.Itoa(ans), fmt.Sprintf("A sphere has only a curved surface, with coefficient 4 x %d^2 = %d.", r, ans))
		}
	case "regular":
		switch o {
		case 1:
			ap := 2 + rand.Intn(8)
			per := 2 * (2 + rand.Intn(9))
			ans := ap * per / 2
			return geoProb(o, fmt.Sprintf("A regular polygon has apothem %d and perimeter %d. What is its area? (enter a number)", ap, per), strconv.Itoa(ans), fmt.Sprintf("Area is one half apothem times perimeter: %d x %d / 2 = %d.", ap, per, ans))
		case 2:
			ap := 2 + rand.Intn(8)
			per := 2 * (2 + rand.Intn(9))
			area := ap * per / 2
			return geoProb(o, fmt.Sprintf("A regular polygon has area %d and perimeter %d. What is its apothem? (enter a number)", area, per), strconv.Itoa(ap), fmt.Sprintf("From A = a p / 2, a = 2A/p = 2 x %d / %d = %d.", area, per, ap))
		default:
			n := 3 + rand.Intn(8)
			s := 2 + rand.Intn(9)
			return geoProb(o, fmt.Sprintf("A regular polygon has %d sides of length %d. What is its perimeter? (enter a number)", n, s), strconv.Itoa(n*s), fmt.Sprintf("The perimeter is %d x %d = %d.", n, s, n*s))
		}
	case "chordseg":
		switch o {
		case 1:
			p := 2 + rand.Intn(8)
			q := 2 + rand.Intn(8)
			c := 2 + rand.Intn(8)
			ans := p * q / c
			if p*q%c != 0 {
				ans = p * q
			}
			if p*q%c != 0 {
				return geoProb(o, fmt.Sprintf("Two chords intersect. The segments of one chord have lengths %d and %d. What is the product of the segments? (enter a number)", p, q), strconv.Itoa(p*q), fmt.Sprintf("The product is %d x %d = %d.", p, q, p*q))
			}
			return geoProb(o, fmt.Sprintf("Two chords intersect. One chord has segments %d and %d, the other has one segment %d. What is its other segment? (enter a number)", p, q, c), strconv.Itoa(ans), fmt.Sprintf("Intersecting chords give %d x %d = %d x d, so d = %d.", p, q, c, ans))
		case 2:
			p := 2 + rand.Intn(8)
			q := 2 + rand.Intn(8)
			return geoProb(o, fmt.Sprintf("Two chords intersect. The segments of one chord are %d and %d. What is the product of those two lengths? (enter a number)", p, q), strconv.Itoa(p*q), fmt.Sprintf("Multiply the segments: %d x %d = %d.", p, q, p*q))
		default:
			triples := [][3]int{{3, 4, 5}, {5, 12, 13}, {8, 15, 17}, {7, 24, 25}, {20, 21, 29}, {9, 40, 41}}
			t := triples[rand.Intn(len(triples))]
			d, half, r := t[0], t[1], t[2]
			return geoProb(o, fmt.Sprintf("A chord lies %d units from the centre of a circle of radius %d. What is the full chord length? (enter a number)", d, r), strconv.Itoa(2*half), fmt.Sprintf("The half-chord is %d from a right triangle with legs %d and %d and hypotenuse %d, so the full chord is %d.", half, d, half, r, 2*half))
		}
	case "secant":
		switch o {
		case 1:
			outer := 2 + rand.Intn(8)
			whole := outer * (2 + rand.Intn(4))
			t2 := outer * whole
			return geoProb(o, fmt.Sprintf("A tangent of length squared %d meets a secant with external part %d. What is the whole secant length? (enter a number)", t2, outer), strconv.Itoa(whole), fmt.Sprintf("Tangent squared equals external times whole: %d = %d x w, so w = %d.", t2, outer, whole))
		case 2:
			outer := 2 + rand.Intn(8)
			whole := outer * (2 + rand.Intn(4))
			return geoProb(o, fmt.Sprintf("A tangent squared has value %d and meets a secant whose whole length is %d. What is the external part? (enter a number)", outer*whole, whole), strconv.Itoa(outer), fmt.Sprintf("Solve %d = x x %d, so x = %d.", outer*whole, whole, outer))
		default:
			p := 2 + rand.Intn(8)
			q := 2 + rand.Intn(8)
			c := 2 + rand.Intn(8)
			if p*q%c != 0 {
				return geoProb(o, fmt.Sprintf("Two secants from a point have products %d and %d. What is the product of the first secant's segments? (enter a number)", p*q, c), strconv.Itoa(p*q), fmt.Sprintf("The product is %d x %d = %d.", p, q, p*q))
			}
			ans := p * q / c
			return geoProb(o, fmt.Sprintf("Two secants from a point have segment products %d and %d x d. What is d? (enter a number)", p*q, c), strconv.Itoa(ans), fmt.Sprintf("Set %d = %d x d, so d = %d.", p*q, c, ans))
		}
	case "similar":
		switch o {
		case 1:
			criteria := []struct {
				name string
				yes  bool
			}{{"AA", true}, {"SSS", true}, {"SAS", true}, {"SSA", false}, {"AAA", true}, {"ASA", true}}
			perm := rand.Perm(len(criteria))[:5]
			count := 0
			var names []string
			for _, i := range perm {
				names = append(names, criteria[i].name)
				if criteria[i].yes {
					count++
				}
			}
			return geoProb(o, fmt.Sprintf("How many of these criteria guarantee that two triangles are similar: %s, %s, %s, %s and %s? (enter a number)", names[0], names[1], names[2], names[3], names[4]), strconv.Itoa(count), fmt.Sprintf("%d of the listed criteria guarantee similarity. AA, SSS, SAS, AAA and ASA are sufficient; SSA is not.", count))
		case 2:
			k := 2 + rand.Intn(5)
			return geoProb(o, fmt.Sprintf("Two similar figures have a scale factor of %d. A side of the smaller figure is %d. What is the corresponding side of the larger figure? (enter a number)", k, b), strconv.Itoa(k*b), fmt.Sprintf("Multiply by the scale factor: %d x %d = %d.", k, b, k*b))
		default:
			k := 2 + rand.Intn(5)
			side := k * (2 + rand.Intn(5))
			return geoProb(o, fmt.Sprintf("Two similar triangles have scale factor %d. A side of the larger triangle is %d. What is the corresponding side of the smaller triangle? (enter a number)", k, side), strconv.Itoa(side/k), fmt.Sprintf("Divide by the scale factor: %d / %d = %d.", side, k, side/k))
		}
	case "composite":
		switch o {
		case 1:
			l := 2 + rand.Intn(9)
			w := 2 + rand.Intn(9)
			base := 2 + rand.Intn(9)
			h := 2 + rand.Intn(9)
			ans := l*w + base*h/2
			return geoProb(o, fmt.Sprintf("A composite figure is a %d by %d rectangle with a triangle of base %d and height %d on top. What is the total area? (enter a number)", l, w, base, h), strconv.Itoa(ans), fmt.Sprintf("Rectangle area is %d x %d = %d, triangle area is %d x %d / 2 = %d, total %d.", l, w, l*w, base, h, base*h/2, ans))
		case 2:
			l := 2 + rand.Intn(9)
			w := 2 + rand.Intn(9)
			cut := 1 + rand.Intn(l-1)
			hcut := 1 + rand.Intn(w-1)
			ans := l*w + cut*hcut
			return geoProb(o, fmt.Sprintf("An L-shaped figure is a %d by %d rectangle with a %d by %d rectangle added. What is its area? (enter a number)", l, w, cut, hcut), strconv.Itoa(ans), fmt.Sprintf("The two rectangle areas are %d and %d, so the total is %d.", l*w, cut*hcut, ans))
		default:
			l := 3 + rand.Intn(9)
			w := 3 + rand.Intn(9)
			cut := 1 + rand.Intn(l-1)
			hcut := 1 + rand.Intn(w-1)
			ans := l*w - cut*hcut
			return geoProb(o, fmt.Sprintf("A %d by %d rectangle has a %d by %d rectangle removed. What is the remaining area? (enter a number)", l, w, cut, hcut), strconv.Itoa(ans), fmt.Sprintf("Subtract the cut-out: %d - %d = %d.", l*w, cut*hcut, ans))
		}
	case "anglebisector":
		switch o {
		case 1:
			return geoProb(o, fmt.Sprintf("An angle of %d degrees is bisected. What is each half? (enter a number)", 2*a), strconv.Itoa(a), fmt.Sprintf("Divide by two: %d / 2 = %d.", 2*a, a))
		case 2:
			p := 2 + rand.Intn(8)
			q := 2 + rand.Intn(8)
			seg := p * (1 + rand.Intn(4))
			ans := seg * q / p
			return geoProb(o, fmt.Sprintf("An angle bisector divides the opposite side. The adjacent sides are %d and %d, and the segment next to the first side is %d. What is the other segment? (enter a number)", p, q, seg), strconv.Itoa(ans), fmt.Sprintf("The ratio gives %d/%d = x/%d, so x = %d.", p, q, seg, ans))
		default:
			p := 2 + rand.Intn(8)
			q := 2 + rand.Intn(8)
			return geoProb(o, fmt.Sprintf("An angle bisector splits the opposite side into lengths %d and %d. What is the whole side? (enter a number)", p, q), strconv.Itoa(p+q), fmt.Sprintf("Add the two parts: %d + %d = %d.", p, q, p+q))
		}
	case "perpbisector":
		switch o {
		case 1:
			x1 := rand.Intn(11) - 5
			x2 := x1 + 2*(1+rand.Intn(5))
			return geoProb(o, fmt.Sprintf("A segment has endpoints with x-coordinates %d and %d. What is the x-coordinate of its midpoint? (enter a number)", x1, x2), strconv.Itoa((x1+x2)/2), fmt.Sprintf("Average the endpoints: (%d + %d)/2 = %d.", x1, x2, (x1+x2)/2))
		case 2:
			return geoProb(o, fmt.Sprintf("A segment of length %d is bisected. How far is the midpoint from either endpoint? (enter a number)", 2*b), strconv.Itoa(b), fmt.Sprintf("Half the length is %d.", b))
		default:
			triples := [][3]int{{3, 4, 5}, {5, 12, 13}, {8, 15, 17}, {7, 24, 25}, {20, 21, 29}, {9, 40, 41}}
			t := triples[rand.Intn(len(triples))]
			offset, half, ans := t[0], t[1], t[2]
			return geoProb(o, fmt.Sprintf("A point on the perpendicular bisector is %d units from the midpoint of a segment whose half-length is %d. How far is it from an endpoint? (enter a number)", offset, half), strconv.Itoa(ans), fmt.Sprintf("The distance is the hypotenuse of a right triangle with legs %d and %d, so it is %d.", offset, half, ans))
		}
	case "quad":
		switch o {
		case 1:
			p := 40 + rand.Intn(80)
			q := 40 + rand.Intn(80)
			r := 40 + rand.Intn(80)
			ans := 360 - p - q - r
			return geoProb(o, fmt.Sprintf("A quadrilateral has angles %d, %d and %d degrees. What is the fourth angle? (enter a number)", p, q, r), strconv.Itoa(ans), fmt.Sprintf("The angles sum to 360: 360 - %d - %d - %d = %d.", p, q, r, ans))
		case 2:
			return geoProb(o, fmt.Sprintf("A parallelogram has one angle of %d degrees. What is the opposite angle? (enter a number)", a), strconv.Itoa(a), fmt.Sprintf("Opposite angles of a parallelogram are equal, so the angle is %d.", a))
		default:
			return geoProb(o, fmt.Sprintf("A parallelogram has one angle of %d degrees. What is the adjacent angle? (enter a number)", a), strconv.Itoa(180-a), fmt.Sprintf("Adjacent angles are supplementary: 180 - %d = %d.", a, 180-a))
		}
	case "arc":
		switch o {
		case 1:
			return geoProb(o, fmt.Sprintf("A central angle measures %d degrees. What is the measure of its intercepted arc? (enter a number)", a), strconv.Itoa(a), fmt.Sprintf("The arc measure equals the central angle, so it is %d degrees.", a))
		case 2:
			ins := 20 + rand.Intn(70)
			return geoProb(o, fmt.Sprintf("An inscribed angle measures %d degrees. What is the measure of its intercepted arc? (enter a number)", ins), strconv.Itoa(2*ins), fmt.Sprintf("The arc is twice the inscribed angle: 2 x %d = %d.", ins, 2*ins))
		default:
			circ := 360 * (1 + rand.Intn(6))
			angles := []int{45, 60, 90, 120, 180}
			angle := angles[rand.Intn(len(angles))]
			ans := circ * angle / 360
			return geoProb(o, fmt.Sprintf("A circle has circumference %d. What is the length of a %d degree arc? (enter a number)", circ, angle), strconv.Itoa(ans), fmt.Sprintf("The arc is %d/360 of the circumference: %d x %d / 360 = %d.", angle, circ, angle, ans))
		}
	case "cross":
		switch o {
		case 1:
			l := 2 + rand.Intn(9)
			w := 2 + rand.Intn(9)
			return geoProb(o, fmt.Sprintf("A rectangular prism is %d by %d by %d. What is the area of a cross section parallel to its base? (enter a number)", l, w, 2+rand.Intn(9)), strconv.Itoa(l*w), fmt.Sprintf("The slice matches the base, so its area is %d x %d = %d.", l, w, l*w))
		case 2:
			n := 3 + rand.Intn(8)
			return geoProb(o, fmt.Sprintf("A prism has a base that is a regular %d-gon. How many faces does it have? (enter a number)", n), strconv.Itoa(n+2), fmt.Sprintf("A prism has two bases plus %d side faces, so %d + 2 = %d faces.", n, n, n+2))
		default:
			n := 3 + rand.Intn(8)
			return geoProb(o, fmt.Sprintf("A prism has a base that is a regular %d-gon. How many edges does it have? (enter a number)", n), strconv.Itoa(3*n), fmt.Sprintf("Each base has %d edges and there are %d vertical edges, so 3 x %d = %d.", n, n, n, 3*n))
		}
	case "symmetry":
		switch o {
		case 1:
			shapes := []struct {
				name  string
				lines int
			}{{"square", 4}, {"rectangle", 2}, {"equilateral triangle", 3}, {"regular pentagon", 5}, {"isosceles triangle", 1}, {"regular hexagon", 6}, {"regular octagon", 8}, {"rhombus", 2}, {"regular heptagon", 7}}
			s := shapes[rand.Intn(len(shapes))]
			return geoProb(o, fmt.Sprintf("How many lines of symmetry does a %s have? (enter a number)", s.name), strconv.Itoa(s.lines), fmt.Sprintf("A %s has %d line(s) of symmetry.", s.name, s.lines))
		case 2:
			orders := []struct {
				name  string
				order int
			}{{"square", 4}, {"equilateral triangle", 3}, {"regular pentagon", 5}, {"rectangle", 2}, {"regular hexagon", 6}, {"regular octagon", 8}, {"rhombus", 2}, {"regular heptagon", 7}}
			s := orders[rand.Intn(len(orders))]
			return geoProb(o, fmt.Sprintf("What is the order of rotational symmetry of a %s? (enter a number)", s.name), strconv.Itoa(s.order), fmt.Sprintf("A %s repeats under a turn of 360/%d degrees, so the order is %d.", s.name, s.order, s.order))
		default:
			total := []struct {
				name string
				n    int
			}{{"square", 8}, {"rectangle", 4}, {"equilateral triangle", 6}, {"regular pentagon", 10}, {"regular hexagon", 12}, {"regular octagon", 16}, {"regular heptagon", 14}, {"isosceles triangle", 2}, {"rhombus", 4}}
			s := total[rand.Intn(len(total))]
			return geoProb(o, fmt.Sprintf("How many symmetries, including the identity, does a %s have? (enter a number)", s.name), strconv.Itoa(s.n), fmt.Sprintf("Add the reflection and rotation symmetries including the identity; a %s has %d.", s.name, s.n))
		}
	case "euler":
		switch o {
		case 1:
			v := 4 + rand.Intn(9)
			e := v + 4 + rand.Intn(5)
			f := 2 - v + e
			return geoProb(o, fmt.Sprintf("A polyhedron has %d vertices and %d edges. What is V - E + F? (enter a number)", v, e), "2", fmt.Sprintf("Euler's formula gives V - E + F = 2 for every convex polyhedron, and here F = %d.", f))
		case 2:
			v := 4 + rand.Intn(9)
			e := v + 4 + rand.Intn(5)
			ans := 2 - v + e
			return geoProb(o, fmt.Sprintf("A polyhedron has %d vertices and %d edges. How many faces does it have? (enter a number)", v, e), strconv.Itoa(ans), fmt.Sprintf("Rearrange V - E + F = 2 to F = 2 - %d + %d = %d.", v, e, ans))
		default:
			v := 4 + rand.Intn(9)
			f := 4 + rand.Intn(7)
			ans := v + f - 2
			return geoProb(o, fmt.Sprintf("A polyhedron has %d vertices and %d faces. How many edges does it have? (enter a number)", v, f), strconv.Itoa(ans), fmt.Sprintf("From V - E + F = 2, E = %d + %d - 2 = %d.", v, f, ans))
		}
	}
	return geoProb(o, "What is 1 + 1? (enter a number)", "2", "This fallback should not appear.")
}
