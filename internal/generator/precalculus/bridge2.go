package precalculus

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

// The 2026 content expansion adds concepts whose three knowledge points ask
// different questions. Each generator below reads KPOrdinal so the practice
// question matches the section the learner is reading.

type precalcBridge2Gen struct{ kind string }

func (g *precalcBridge2Gen) Generate(ctx generator.GeneratorContext) generator.Problem {
	ordinal := precalcBridgeOrdinal(ctx)
	switch g.kind {
	case "poly_graph":
		return polynomialGraphProblem(ordinal)
	case "logistic":
		return logisticProblem(ordinal)
	case "unit_vector":
		return unitVectorProblem(ordinal)
	case "eccentricity":
		return eccentricityProblem(ordinal)
	case "seq_explicit":
		return explicitSequenceProblem(ordinal)
	case "rotated_conic":
		return rotatedConicProblem(ordinal)
	default:
		return vector3DProblem(ordinal)
	}
}

func polynomialGraphProblem(o int) generator.Problem {
	switch o {
	case 1:
		d := 2 + rand.Intn(6)
		lead := 1
		if rand.Intn(2) == 0 {
			lead = -1
		}
		endsUp := 0
		if d%2 == 0 {
			if lead > 0 {
				endsUp = 2
			}
		} else {
			endsUp = 1
		}
		sign := "positive"
		if lead < 0 {
			sign = "negative"
		}
		return precalcBridgeProblem(o,
			fmt.Sprintf("A polynomial has degree %d and a %s leading coefficient. How many of its two ends point upward? (enter a number)", d, sign),
			strconv.Itoa(endsUp),
			fmt.Sprintf("An even degree sends both ends the same way, and an odd degree sends them opposite ways. With degree %d and a %s leading coefficient, %d end(s) point upward.", d, sign, endsUp),
			map[string]string{"degree": strconv.Itoa(d), "answer": strconv.Itoa(endsUp)})
	case 2:
		d := 3 + rand.Intn(6)
		turns := d - 1
		return precalcBridgeProblem(o,
			fmt.Sprintf("What is the largest possible number of turning points of a polynomial of degree %d? (enter a number)", d),
			strconv.Itoa(turns),
			fmt.Sprintf("A polynomial of degree %d has at most %d - 1 = %d turning points.", d, d, turns),
			map[string]string{"degree": strconv.Itoa(d), "answer": strconv.Itoa(turns)})
	default:
		d := 4 + rand.Intn(4)
		a := 1 + rand.Intn(4)
		b := a + 1 + rand.Intn(3)
		m := 1 + rand.Intn(3)
		n := d - m
		return precalcBridgeProblem(o,
			fmt.Sprintf("A polynomial of degree %d has a zero of multiplicity %d at x=%d. If all its remaining zeros are at x=%d, what is the multiplicity there? (enter a number)", d, m, a, b),
			strconv.Itoa(n),
			fmt.Sprintf("The multiplicities add to the degree: %d + n = %d, so n = %d.", m, d, n),
			map[string]string{"degree": strconv.Itoa(d), "answer": strconv.Itoa(n)})
	}
}

func logisticProblem(o int) generator.Problem {
	l := 100 + 20*rand.Intn(21)
	k := 1 + rand.Intn(5)
	a := 1 + rand.Intn(4)
	switch o {
	case 1:
		return precalcBridgeProblem(o,
			fmt.Sprintf("A population follows P(t) = %d / (1 + %d e^(-%dt)). What is the carrying capacity? (enter a number)", l, a, k),
			strconv.Itoa(l),
			fmt.Sprintf("The numerator of the logistic model is the carrying capacity, so it is %d.", l),
			map[string]string{"answer": strconv.Itoa(l), "kind": "capacity"})
	case 2:
		start := l / (1 + a)
		l = start * (1 + a)
		return precalcBridgeProblem(o,
			fmt.Sprintf("A population follows P(t) = %d / (1 + %d e^(-%dt)). What is the initial population P(0)? (enter a number)", l, a, k),
			strconv.Itoa(start),
			fmt.Sprintf("At t=0 the exponential factor is 1, so P(0) = %d / (1 + %d) = %d.", l, a, start),
			map[string]string{"answer": strconv.Itoa(start), "kind": "initial"})
	default:
		l = 2 * (50 + 10*rand.Intn(21))
		half := l / 2
		return precalcBridgeProblem(o,
			fmt.Sprintf("A logistic population has carrying capacity %d. At what population is the growth rate largest? (enter a number)", l),
			strconv.Itoa(half),
			fmt.Sprintf("The logistic curve is steepest halfway to the carrying capacity, at %d / 2 = %d.", l, half),
			map[string]string{"answer": strconv.Itoa(half), "kind": "half"})
	}
}

func unitVectorProblem(o int) generator.Problem {
	triples := [][3]int{{3, 4, 5}, {6, 8, 10}, {5, 12, 13}, {8, 15, 17}, {7, 24, 25}, {20, 21, 29}, {9, 40, 41}, {12, 16, 20}}
	t := triples[rand.Intn(len(triples))]
	switch o {
	case 1:
		return precalcBridgeProblem(o,
			fmt.Sprintf("What is the length of the vector (%d, %d)? (enter a number)", t[0], t[1]),
			strconv.Itoa(t[2]),
			fmt.Sprintf("The length is the square root of %d^2 + %d^2 = %d, which is %d.", t[0], t[1], t[0]*t[0]+t[1]*t[1], t[2]),
			map[string]string{"answer": strconv.Itoa(t[2]), "kind": "length"})
	case 2:
		comp := t[0]
		if rand.Intn(2) == 0 {
			comp = t[1]
		}
		val := precalcDecimal(float64(comp) / float64(t[2]))
		return precalcBridgeProblem(o,
			fmt.Sprintf("What is the %s-component of the unit vector in the direction of (%d, %d)? (enter a decimal)", map[bool]string{true: "x", false: "y"}[comp == t[0]], t[0], t[1]),
			val,
			fmt.Sprintf("Divide the vector by its length %d. The component is %d/%d = %s.", t[2], comp, t[2], val),
			map[string]string{"answer": val, "kind": "component"})
	default:
		quads := [][4]int{{1, 2, 2, 3}, {2, 3, 6, 7}, {3, 4, 12, 13}, {1, 4, 8, 9}, {4, 4, 7, 9}, {2, 6, 9, 11}}
		q := quads[rand.Intn(len(quads))]
		val := precalcDecimal(float64(q[0]) / float64(q[3]))
		return precalcBridgeProblem(o,
			fmt.Sprintf("What is the x-component of the unit vector in the direction of (%d, %d, %d)? (enter a decimal)", q[0], q[1], q[2]),
			val,
			fmt.Sprintf("The length is %d, so the x-component is %d/%d = %s.", q[3], q[0], q[3], val),
			map[string]string{"answer": val, "kind": "three_dimensional"})
	}
}

func eccentricityProblem(o int) generator.Problem {
	ellipses := [][3]int{{5, 4, 3}, {5, 3, 4}, {13, 12, 5}, {13, 5, 12}, {10, 8, 6}, {10, 6, 8}, {25, 24, 7}, {25, 7, 24}}
	hyperbolas := [][3]int{{3, 4, 5}, {5, 12, 13}, {8, 15, 17}, {7, 24, 25}, {4, 3, 5}, {12, 5, 13}, {9, 40, 41}, {20, 21, 29}}
	if o == 2 {
		vals := []string{"0", "0.25", "0.5", "0.75", "1", "1.2", "1.5", "2"}
		v := vals[rand.Intn(len(vals))]
		code := 1
		if v == "1" {
			code = 2
		} else if v == "1.2" || v == "1.5" || v == "2" {
			code = 3
		}
		return precalcBridgeProblem(o,
			fmt.Sprintf("A conic has eccentricity %s. What type is it? Enter 1 for an ellipse, 2 for a parabola and 3 for a hyperbola.", v),
			strconv.Itoa(code),
			fmt.Sprintf("An eccentricity below 1 gives an ellipse, exactly 1 gives a parabola and above 1 gives a hyperbola. Here e = %s.", v),
			map[string]string{"answer": strconv.Itoa(code), "kind": "classify"})
	}
	t := ellipses[rand.Intn(len(ellipses))]
	kind := "ellipse"
	if o == 3 {
		t = hyperbolas[rand.Intn(len(hyperbolas))]
		kind = "hyperbola"
	}
	val := precalcDecimal(float64(t[2]) / float64(t[0]))
	return precalcBridgeProblem(o,
		fmt.Sprintf("A %s has semi-major axis %d and focal distance %d. What is its eccentricity? (enter a decimal)", kind, t[0], t[2]),
		val,
		fmt.Sprintf("Eccentricity is c/a = %d/%d = %s.", t[2], t[0], val),
		map[string]string{"answer": val, "kind": kind})
}

func explicitSequenceProblem(o int) generator.Problem {
	switch o {
	case 1:
		a := 1 + rand.Intn(9)
		d := 1 + rand.Intn(7)
		n := 4 + rand.Intn(8)
		ans := a + (n-1)*d
		return precalcBridgeProblem(o,
			fmt.Sprintf("An arithmetic sequence has first term %d and common difference %d. Use a_n = a_1 + (n-1)d to find a_%d. (enter a number)", a, d, n),
			strconv.Itoa(ans),
			fmt.Sprintf("a_%d = %d + (%d - 1)(%d) = %d.", n, a, n, d, ans),
			map[string]string{"answer": strconv.Itoa(ans), "kind": "arithmetic"})
	case 2:
		a := 1 + rand.Intn(4)
		r := 2 + rand.Intn(3)
		n := 3 + rand.Intn(5)
		ans := a
		for i := 1; i < n; i++ {
			ans *= r
		}
		return precalcBridgeProblem(o,
			fmt.Sprintf("A geometric sequence has first term %d and common ratio %d. Use a_n = a_1 r^(n-1) to find a_%d. (enter a number)", a, r, n),
			strconv.Itoa(ans),
			fmt.Sprintf("a_%d = %d times %d^%d = %d.", n, a, r, n-1, ans),
			map[string]string{"answer": strconv.Itoa(ans), "kind": "geometric"})
	default:
		b := rand.Intn(7) - 3
		c := rand.Intn(9) - 4
		n := 2 + rand.Intn(8)
		ans := n*n + b*n + c
		return precalcBridgeProblem(o,
			fmt.Sprintf("A sequence is defined by a_n = n^2 + (%d)n + (%d). What is a_%d? (enter a number)", b, c, n),
			strconv.Itoa(ans),
			fmt.Sprintf("Substitute n = %d: %d^2 + (%d)(%d) + (%d) = %d.", n, n, b, n, c, ans),
			map[string]string{"answer": strconv.Itoa(ans), "kind": "polynomial"})
	}
}

func rotatedConicProblem(o int) generator.Problem {
	a := 1 + rand.Intn(4)
	c := 1 + rand.Intn(4)
	b := rand.Intn(9) - 4
	switch o {
	case 1:
		disc := b*b - 4*a*c
		return precalcBridgeProblem(o,
			fmt.Sprintf("For the conic %dx^2 + (%d)xy + %dy^2 = 1, compute the discriminant B^2 - 4AC. (enter a number)", a, b, c),
			strconv.Itoa(disc),
			fmt.Sprintf("B^2 - 4AC = (%d)^2 - 4(%d)(%d) = %d.", b, a, c, disc),
			map[string]string{"answer": strconv.Itoa(disc), "kind": "discriminant"})
	case 2:
		disc := b*b - 4*a*c
		code := 1
		if disc == 0 {
			code = 2
		} else if disc > 0 {
			code = 3
		}
		return precalcBridgeProblem(o,
			fmt.Sprintf("For the conic %dx^2 + (%d)xy + %dy^2 + x = 0, the discriminant is %d. What type is it? Enter 1 for an ellipse, 2 for a parabola and 3 for a hyperbola.", a, b, c, disc),
			strconv.Itoa(code),
			fmt.Sprintf("A negative discriminant gives an ellipse, zero gives a parabola and a positive value gives a hyperbola. The discriminant is %d.", disc),
			map[string]string{"answer": strconv.Itoa(code), "kind": "classify"})
	default:
		types := []struct {
			text string
			deg  int
		}{{"0", 0}, {"1", 45}, {"-1", 135}, {"sqrt(3)", 60}, {"-sqrt(3)", 120}, {"sqrt(3)/3", 30}, {"-sqrt(3)/3", 150}}
		t := types[rand.Intn(len(types))]
		return precalcBridgeProblem(o,
			fmt.Sprintf("A rotation removes the xy term when tan(2 theta) = %s. What is 2 theta in degrees, between 0 and 180? (enter a number)", t.text),
			strconv.Itoa(t.deg),
			fmt.Sprintf("The angle 2 theta has tangent %s, so 2 theta = %d degrees in the required range.", t.text, t.deg),
			map[string]string{"answer": strconv.Itoa(t.deg), "kind": "angle"})
	}
}

func vector3DProblem(o int) generator.Problem {
	quads := [][4]int{{1, 2, 2, 3}, {2, 3, 6, 7}, {3, 4, 12, 13}, {1, 4, 8, 9}, {4, 4, 7, 9}, {2, 6, 9, 11}}
	q := quads[rand.Intn(len(quads))]
	switch o {
	case 1:
		return precalcBridgeProblem(o,
			fmt.Sprintf("What is the length of the vector (%d, %d, %d)? (enter a number)", q[0], q[1], q[2]),
			strconv.Itoa(q[3]),
			fmt.Sprintf("The length is the square root of %d^2 + %d^2 + %d^2 = %d, which is %d.", q[0], q[1], q[2], q[0]*q[0]+q[1]*q[1]+q[2]*q[2], q[3]),
			map[string]string{"answer": strconv.Itoa(q[3]), "kind": "length"})
	case 2:
		u := [3]int{rand.Intn(7) - 3, rand.Intn(7) - 3, rand.Intn(7) - 3}
		v := [3]int{rand.Intn(7) - 3, rand.Intn(7) - 3, rand.Intn(7) - 3}
		dot := u[0]*v[0] + u[1]*v[1] + u[2]*v[2]
		return precalcBridgeProblem(o,
			fmt.Sprintf("What is the dot product of (%d, %d, %d) and (%d, %d, %d)? (enter a number)", u[0], u[1], u[2], v[0], v[1], v[2]),
			strconv.Itoa(dot),
			fmt.Sprintf("Multiply matching components and add: %d(%d) + %d(%d) + %d(%d) = %d.", u[0], v[0], u[1], v[1], u[2], v[2], dot),
			map[string]string{"answer": strconv.Itoa(dot), "kind": "dot"})
	default:
		p := [3]int{rand.Intn(5), rand.Intn(5), rand.Intn(5)}
		r := [3]int{p[0] + q[0], p[1] + q[1], p[2] + q[2]}
		return precalcBridgeProblem(o,
			fmt.Sprintf("What is the distance between (%d, %d, %d) and (%d, %d, %d)? (enter a number)", p[0], p[1], p[2], r[0], r[1], r[2]),
			strconv.Itoa(q[3]),
			fmt.Sprintf("The differences are %d, %d and %d, whose squares add to %d. The distance is %d.", q[0], q[1], q[2], q[3]*q[3], q[3]),
			map[string]string{"answer": strconv.Itoa(q[3]), "kind": "distance"})
	}
}

func precalcBridgeProblem(o int, question, answer, explanation string, facts map[string]string) generator.Problem {
	if facts == nil {
		facts = map[string]string{}
	}
	facts["answer"] = answer
	return generator.Problem{Question: question, Answer: answer, Explanation: explanation, Facts: facts, KPOrdinal: o}
}
