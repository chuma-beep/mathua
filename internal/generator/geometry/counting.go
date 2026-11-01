package geometry

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

// geoVerdict is one statement in a count question. Every pool is written with
// exactly eight true and eight false statements so the answer spans zero to
// eight without the extremes becoming vanishingly rare.
type geoVerdict struct {
	label string
	yes   bool
	note  string
}

const geoDraw = 8

// countGeoVerdicts asks how many of eight drawn statements are true. It picks
// the answer first, then takes that many true and the rest false statements, so
// all nine answers appear with equal probability.
func countGeoVerdicts(lead string, pool []geoVerdict) generator.Problem {
	var yes, no []geoVerdict
	for _, v := range pool {
		if v.yes {
			yes = append(yes, v)
		} else {
			no = append(no, v)
		}
	}
	count := rand.Intn(geoDraw + 1)
	rand.Shuffle(len(yes), func(i, j int) { yes[i], yes[j] = yes[j], yes[i] })
	rand.Shuffle(len(no), func(i, j int) { no[i], no[j] = no[j], no[i] })

	chosen := make([]geoVerdict, 0, geoDraw)
	chosen = append(chosen, yes[:count]...)
	chosen = append(chosen, no[:geoDraw-count]...)
	rand.Shuffle(len(chosen), func(i, j int) { chosen[i], chosen[j] = chosen[j], chosen[i] })

	question := lead
	notes := make([]string, len(chosen))
	for attempts := 0; attempts < 40; attempts++ {
		parts := make([]string, len(chosen))
		for i, v := range chosen {
			parts[i] = fmt.Sprintf("(%s) %s", string(rune('a'+i)), v.label)
			verdict := "false"
			if v.yes {
				verdict = "true"
			}
			notes[i] = fmt.Sprintf("%s: %s, %s", v.label, verdict, v.note)
		}
		question = fmt.Sprintf("%s %s (enter a number)", lead, strings.Join(parts, " "))
		if !strings.Contains(question, fmt.Sprintf("%d", count)) {
			break
		}
		count = rand.Intn(geoDraw + 1)
		chosen = chosen[:0]
		chosen = append(chosen, yes[:count]...)
		chosen = append(chosen, no[:geoDraw-count]...)
		rand.Shuffle(len(chosen), func(i, j int) { chosen[i], chosen[j] = chosen[j], chosen[i] })
	}

	return generator.Problem{
		Question:    question,
		Answer:      fmt.Sprintf("%d", count),
		Explanation: fmt.Sprintf("Exactly %d of the 8 statements are true. %s", count, strings.Join(notes, "; ")),
		Facts:       map[string]string{"count": fmt.Sprintf("%d", count)},
	}
}

func pointsLinesProblem() generator.Problem {
	return countGeoVerdicts("How many of these statements about points, lines and planes are true?", []geoVerdict{
		{"A point has no length, width or depth", true, "a point is zero-dimensional"},
		{"Exactly one straight line passes through two distinct points", true, "two points determine a unique line"},
		{"A line extends infinitely in both directions", true, "a line has no endpoints"},
		{"A ray has exactly one endpoint", true, "the other direction continues forever"},
		{"A line segment has exactly two endpoints", true, "the endpoints bound a measurable length"},
		{"Collinear points lie on the same straight line", true, "this is the definition of collinear"},
		{"Coplanar points lie in the same plane", true, "this is the definition of coplanar"},
		{"Two distinct lines intersect in at most one point", true, "a second intersection would make them the same line"},
		{"A ray has two endpoints", false, "a ray has one endpoint"},
		{"A line segment extends infinitely in both directions", false, "a segment is bounded by two endpoints"},
		{"A point has length but no width", false, "a point has no size at all"},
		{"Three collinear points always lie in different planes", false, "a line lies in infinitely many planes, not one per point"},
		{"Two distinct points determine infinitely many distinct lines", false, "they determine exactly one line"},
		{"A line has exactly one endpoint", false, "a line has no endpoints"},
		{"Skew lines intersect at a right angle", false, "skew lines never intersect"},
		{"Any three points are always collinear", false, "three points can form a triangle"},
	})
}

func angleTypesProblem() generator.Problem {
	return countGeoVerdicts("How many of these statements about angles are true?", []geoVerdict{
		{"An acute angle measures less than 90 degrees", true, "acute means strictly between 0 and 90"},
		{"A right angle measures exactly 90 degrees", true, "this is the definition of a right angle"},
		{"An obtuse angle measures between 90 and 180 degrees", true, "obtuse is strictly greater than 90 and less than 180"},
		{"A straight angle measures 180 degrees", true, "its sides form a straight line"},
		{"Complementary angles sum to 90 degrees", true, "this is the definition of complementary"},
		{"Supplementary angles sum to 180 degrees", true, "this is the definition of supplementary"},
		{"Vertical angles are equal", true, "opposite angles formed by two lines are congruent"},
		{"Adjacent angles on a straight line sum to 180 degrees", true, "they form a linear pair"},
		{"An acute angle measures more than 90 degrees", false, "that would be obtuse"},
		{"A right angle measures 180 degrees", false, "a right angle measures 90 degrees"},
		{"An obtuse angle measures less than 90 degrees", false, "that would be acute"},
		{"A straight angle measures 360 degrees", false, "a full turn is 360, a straight angle is 180"},
		{"Complementary angles sum to 180 degrees", false, "complementary angles sum to 90"},
		{"Supplementary angles sum to 90 degrees", false, "supplementary angles sum to 180"},
		{"Vertical angles are supplementary in every case", false, "they are equal, and supplementary only when each is 90"},
		{"Adjacent angles on a straight line are always equal", false, "they sum to 180 but need not be equal"},
	})
}

func angleMeasureProblem() generator.Problem {
	switch rand.Intn(4) {
	case 0:
		a := 10 + rand.Intn(71)
		return generator.Problem{
			Question:    fmt.Sprintf("An angle measures %d degrees. How many more degrees are needed to reach a right angle? (enter a number)", a),
			Answer:      fmt.Sprintf("%d", 90-a),
			Explanation: fmt.Sprintf("A right angle is 90 degrees, so 90 - %d = %d.", a, 90-a),
			Facts:       map[string]string{"given": fmt.Sprintf("%d", a), "missing": fmt.Sprintf("%d", 90-a)},
		}
	case 1:
		a := 91 + rand.Intn(80)
		return generator.Problem{
			Question:    fmt.Sprintf("An angle measures %d degrees. How many more degrees are needed to reach a straight angle? (enter a number)", a),
			Answer:      fmt.Sprintf("%d", 180-a),
			Explanation: fmt.Sprintf("A straight angle is 180 degrees, so 180 - %d = %d.", a, 180-a),
			Facts:       map[string]string{"given": fmt.Sprintf("%d", a), "missing": fmt.Sprintf("%d", 180-a)},
		}
	case 2:
		a := 10 + rand.Intn(71)
		return generator.Problem{
			Question:    fmt.Sprintf("Two angles are complementary. One measures %d degrees. What is the other? (enter a number)", a),
			Answer:      fmt.Sprintf("%d", 90-a),
			Explanation: fmt.Sprintf("Complementary angles sum to 90, so the other angle is 90 - %d = %d.", a, 90-a),
			Facts:       map[string]string{"given": fmt.Sprintf("%d", a), "missing": fmt.Sprintf("%d", 90-a)},
		}
	default:
		a := 30 + rand.Intn(121)
		return generator.Problem{
			Question:    fmt.Sprintf("Two angles are supplementary. One measures %d degrees. What is the other? (enter a number)", a),
			Answer:      fmt.Sprintf("%d", 180-a),
			Explanation: fmt.Sprintf("Supplementary angles sum to 180, so the other angle is 180 - %d = %d.", a, 180-a),
			Facts:       map[string]string{"given": fmt.Sprintf("%d", a), "missing": fmt.Sprintf("%d", 180-a)},
		}
	}
}

func triangleTypesProblem() generator.Problem {
	return countGeoVerdicts("How many of these statements about triangles are true?", []geoVerdict{
		{"An equilateral triangle has three equal sides", true, "all three sides have the same length"},
		{"An isosceles triangle has at least two equal sides", true, "this is the defining property"},
		{"A scalene triangle has no equal sides", true, "all three sides have different lengths"},
		{"Every equilateral triangle is isosceles", true, "three equal sides include a pair of equal sides"},
		{"The angles of a triangle sum to 180 degrees", true, "this holds in Euclidean geometry"},
		{"A triangle with sides 3, 4 and 5 is scalene", true, "all three sides differ"},
		{"A triangle with sides 5, 5 and 6 is isosceles", true, "two sides are equal"},
		{"An equilateral triangle has three 60 degree angles", true, "180 divided equally is 60"},
		{"A scalene triangle has two equal sides", false, "scalene means no two sides are equal"},
		{"An isosceles triangle has no equal sides", false, "isosceles requires two equal sides"},
		{"An equilateral triangle has exactly two equal sides", false, "it has three equal sides"},
		{"The angles of a triangle sum to 360 degrees", false, "they sum to 180"},
		{"A triangle with sides 3, 4 and 5 is equilateral", false, "its sides are all different"},
		{"A triangle with sides 5, 5 and 6 is scalene", false, "two sides are equal"},
		{"An equilateral triangle has a right angle", false, "each angle is 60 degrees"},
		{"Every isosceles triangle is equilateral", false, "two equal sides do not force a third"},
	})
}

func quadTypesProblem() generator.Problem {
	return countGeoVerdicts("How many of these statements about quadrilaterals are true?", []geoVerdict{
		{"A square has four equal sides and four right angles", true, "this is the definition of a square"},
		{"A rectangle has four right angles", true, "this is the defining property"},
		{"A rhombus has four equal sides", true, "this is the defining property"},
		{"A parallelogram has both pairs of opposite sides parallel", true, "this is the definition"},
		{"A square is a rectangle", true, "it has four right angles"},
		{"A square is a rhombus", true, "it has four equal sides"},
		{"The diagonals of a rectangle are equal", true, "the diagonals bisect each other and have equal length"},
		{"A parallelogram has opposite angles equal", true, "opposite angles of a parallelogram are congruent"},
		{"A rhombus must have four right angles", false, "a rhombus can have oblique angles"},
		{"A trapezoid has two pairs of parallel sides", false, "a trapezoid has one pair, and two pairs describe a parallelogram"},
		{"Every parallelogram is a square", false, "a parallelogram need not have right angles or equal sides"},
		{"A rectangle has all four sides equal", false, "only a square must have equal sides"},
		{"A square has no right angles", false, "all four angles are right angles"},
		{"A rhombus has no equal sides", false, "all four sides are equal"},
		{"A trapezoid has four equal sides", false, "that would be a rhombus"},
		{"Every rectangle is a square", false, "a rectangle need not have equal sides"},
	})
}

func circlePartsProblem() generator.Problem {
	return countGeoVerdicts("How many of these statements about circle parts are true?", []geoVerdict{
		{"A radius joins the center to a point on the circle", true, "this is the definition of a radius"},
		{"A diameter passes through the center", true, "a diameter is a chord through the center"},
		{"A diameter is twice a radius", true, "this follows from the definition"},
		{"A chord joins two points on the circle", true, "this is the definition of a chord"},
		{"A tangent touches the circle at exactly one point", true, "this is the defining property"},
		{"A secant intersects the circle at two points", true, "this is the defining property"},
		{"The circumference is the distance around the circle", true, "it is the boundary length"},
		{"All radii of a circle have equal length", true, "every point on the circle is equally distant from the center"},
		{"A diameter is half a radius", false, "a diameter is twice a radius"},
		{"A tangent crosses the circle at two points", false, "that is a secant"},
		{"A secant touches the circle at exactly one point", false, "that is a tangent"},
		{"Every chord passes through the center", false, "only a diameter does"},
		{"A radius joins two points on the circle", false, "that is a chord"},
		{"The circumference is the area of the circle", false, "circumference is a length, area is a surface measure"},
		{"A circle has exactly one chord", false, "a circle has infinitely many chords"},
		{"Every chord is a diameter", false, "most chords miss the center"},
	})
}

func coordPlotProblem() generator.Problem {
	x := rand.Intn(11) - 5
	y := rand.Intn(11) - 5
	switch rand.Intn(4) {
	case 0:
		return generator.Problem{
			Question:    fmt.Sprintf("What is the x-coordinate of the point (%d, %d)? (enter a number)", x, y),
			Answer:      fmt.Sprintf("%d", x),
			Explanation: fmt.Sprintf("In (%d, %d) the first entry is the x-coordinate: %d.", x, y, x),
			Facts:       map[string]string{"x": fmt.Sprintf("%d", x), "y": fmt.Sprintf("%d", y), "which": "x"},
		}
	case 1:
		return generator.Problem{
			Question:    fmt.Sprintf("What is the y-coordinate of the point (%d, %d)? (enter a number)", x, y),
			Answer:      fmt.Sprintf("%d", y),
			Explanation: fmt.Sprintf("In (%d, %d) the second entry is the y-coordinate: %d.", x, y, y),
			Facts:       map[string]string{"x": fmt.Sprintf("%d", x), "y": fmt.Sprintf("%d", y), "which": "y"},
		}
	case 2:
		if x == 0 {
			x = 1
		}
		if y == 0 {
			y = -1
		}
		q := 1
		if x < 0 && y > 0 {
			q = 2
		} else if x < 0 && y < 0 {
			q = 3
		} else if x > 0 && y < 0 {
			q = 4
		}
		return generator.Problem{
			Question:    fmt.Sprintf("Which quadrant number contains the point (%d, %d)? (enter 1, 2, 3 or 4)", x, y),
			Answer:      fmt.Sprintf("%d", q),
			Explanation: fmt.Sprintf("The signs of the coordinates put (%d, %d) in quadrant %d.", x, y, q),
			Facts:       map[string]string{"x": fmt.Sprintf("%d", x), "y": fmt.Sprintf("%d", y), "quadrant": fmt.Sprintf("%d", q)},
		}
	default:
		return generator.Problem{
			Question:    fmt.Sprintf("What is the sum of the coordinates of the point (%d, %d)? (enter a number)", x, y),
			Answer:      fmt.Sprintf("%d", x+y),
			Explanation: fmt.Sprintf("Add the coordinates: %d + (%d) = %d.", x, y, x+y),
			Facts:       map[string]string{"x": fmt.Sprintf("%d", x), "y": fmt.Sprintf("%d", y), "sum": fmt.Sprintf("%d", x+y)},
		}
	}
}
