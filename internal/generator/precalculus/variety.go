package precalculus

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

// conicText prints a conic without the clutter a format string produces: no
// unit coefficients, no zero terms, and a leading positive term with no sign.
func conicText(a, b, c, d, e, f int) string {
	parts := []string{}
	add := func(coef int, suffix string) {
		if coef == 0 {
			return
		}
		core := suffix
		mag := coef
		if mag < 0 {
			mag = -mag
		}
		if mag != 1 || suffix == "" {
			core = strconv.Itoa(mag) + suffix
		}
		if len(parts) == 0 {
			if coef < 0 {
				parts = append(parts, "-"+core)
			} else {
				parts = append(parts, core)
			}
			return
		}
		if coef < 0 {
			parts = append(parts, " - "+core)
		} else {
			parts = append(parts, " + "+core)
		}
	}
	add(a, "x^2")
	add(b, "xy")
	add(c, "y^2")
	add(d, "x")
	add(e, "y")
	add(f, "")
	return strings.Join(parts, "")
}

type decompositionGen struct{}

type generalConicGen struct{}

func shiftString(a int) string {
	if a >= 0 {
		return "x+" + strconv.Itoa(a)
	}
	return "x-" + strconv.Itoa(-a)
}

// decompositionGen draws the inner function and asks for it back, publishing
// the drawn parameters so the batch test can rebuild the expected expression
// rather than compare against a table.
func (g *decompositionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch rand.Intn(6) {
	case 0:
		a := rand.Intn(9) + 1
		inner := shiftString(a)
		return generator.Problem{
			Question:    fmt.Sprintf("Write h(x) = (%s)^2 as f(g(x)) with f(x) = x^2. What is g(x)? (enter an expression)", inner),
			Answer:      inner,
			Explanation: fmt.Sprintf("The inner function is whatever is squared, so g(x) = %s.", inner),
			Facts:       map[string]string{"kind": "square_shift", "a": strconv.Itoa(a), "answer": inner},
		}
	case 1:
		a := rand.Intn(9) + 1
		inner := fmt.Sprintf("x^2+%d", a)
		return generator.Problem{
			Question:    fmt.Sprintf("Write h(x) = sqrt(%s) as f(g(x)) with f(x) = sqrt(x). What is g(x)? (enter an expression)", inner),
			Answer:      inner,
			Explanation: fmt.Sprintf("The inner function is the radicand, so g(x) = %s.", inner),
			Facts:       map[string]string{"kind": "sqrt", "a": strconv.Itoa(a), "answer": inner},
		}
	case 2:
		a := rand.Intn(9) + 1
		inner := shiftString(-a)
		return generator.Problem{
			Question:    fmt.Sprintf("Write h(x) = 1/(%s) as f(g(x)) with f(x) = 1/x. What is g(x)? (enter an expression)", inner),
			Answer:      inner,
			Explanation: fmt.Sprintf("The inner function is the denominator, so g(x) = %s.", inner),
			Facts:       map[string]string{"kind": "reciprocal", "a": strconv.Itoa(-a), "answer": inner},
		}
	case 3:
		k := rand.Intn(7) + 2
		inner := fmt.Sprintf("%d*x", k)
		return generator.Problem{
			Question:    fmt.Sprintf("Write h(x) = e^(%s) as f(g(x)) with f(x) = e^x. What is g(x)? (enter an expression)", inner),
			Answer:      inner,
			Explanation: fmt.Sprintf("The exponent is the inner function, so g(x) = %s.", inner),
			Facts:       map[string]string{"kind": "exp", "k": strconv.Itoa(k), "answer": inner},
		}
	case 4:
		a := rand.Intn(9) + 1
		inner := shiftString(a)
		return generator.Problem{
			Question:    fmt.Sprintf("Write h(x) = (%s)^3 as f(g(x)) with f(x) = x^3. What is g(x)? (enter an expression)", inner),
			Answer:      inner,
			Explanation: fmt.Sprintf("The inner function is whatever is cubed, so g(x) = %s.", inner),
			Facts:       map[string]string{"kind": "cube", "a": strconv.Itoa(a), "answer": inner},
		}
	default:
		k := rand.Intn(4) + 2
		inner := fmt.Sprintf("x^%d", k)
		return generator.Problem{
			Question:    fmt.Sprintf("Write h(x) = ln(%s) as f(g(x)) with f(x) = ln(x). What is g(x)? (enter an expression)", inner),
			Answer:      inner,
			Explanation: fmt.Sprintf("The argument of the logarithm is the inner function, so g(x) = %s.", inner),
			Facts:       map[string]string{"kind": "log", "k": strconv.Itoa(k), "answer": inner},
		}
	}
}

// generalConicGen prints a conic from drawn coefficients and asks for the
// discriminant, so the classification follows from the arithmetic rather than
// from five remembered examples.
func (g *generalConicGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := rand.Intn(3) + 1
	c := rand.Intn(3) + 1
	if rand.Intn(3) == 0 {
		c = -c
	}
	b := rand.Intn(4)
	d := rand.Intn(7) - 3
	e := rand.Intn(7) - 3
	f := rand.Intn(7) - 3

	disc := b*b - 4*a*c
	kind := "ellipse"
	switch {
	case disc > 0:
		kind = "hyperbola"
	case disc == 0:
		kind = "parabola"
	case a == c && b == 0:
		kind = "circle"
	}

	return generator.Problem{
		Question: fmt.Sprintf("For the conic %s = 0, what is the discriminant B^2 - 4AC? (enter a number)", conicText(a, b, c, d, e, f)),
		Answer:   strconv.Itoa(disc),
		Explanation: fmt.Sprintf(
			"Here A = %d, B = %d and C = %d. B^2 is %d and 4AC is %d, so the discriminant is %d, which makes the conic %s.",
			a, b, c, b*b, 4*a*c, disc, kind),
		Facts: map[string]string{
			"a": strconv.Itoa(a), "b": strconv.Itoa(b), "c": strconv.Itoa(c), "answer": strconv.Itoa(disc),
		},
	}
}
