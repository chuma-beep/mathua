package algebra

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

// ---------------------------------------------------------------------------
// Piecewise function graphs
// ---------------------------------------------------------------------------

type piecewiseGraphGen struct{}

func (g *piecewiseGraphGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = rand.Intn(3) + 1
	}

	boundary := rand.Intn(5) + 1 // 1..5
	m1 := rand.Intn(4) + 1
	m2 := rand.Intn(4) + 1
	b1 := rand.Intn(11) - 5
	b2 := rand.Intn(11) - 5

	var x, value int
	var branch string
	switch kp {
	case 1:
		x = boundary - 1 - rand.Intn(3)
		value = m1*x + b1
		branch = "left"
	case 2:
		x = boundary + rand.Intn(3)
		value = m2*x + b2
		branch = "right"
	default:
		x = boundary
		value = m2*x + b2
		branch = "right"
	}

	question := fmt.Sprintf("For f(x) = %s when x < %d, and f(x) = %s when x ≥ %d, what is f(%d)? (enter a number)",
		formatLine(m1, b1), boundary, formatLine(m2, b2), boundary, x)
	explanation := fmt.Sprintf("The input %d satisfies x %s %d, so use the %s branch: %s = %d.",
		x, map[bool]string{true: "<", false: "≥"}[branch == "left"], boundary,
		branch, formatLine(map[string]int{"left": m1, "right": m2}[branch], map[string]int{"left": b1, "right": b2}[branch]), value)

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(value),
		Explanation: explanation,
		Facts: map[string]string{
			"x": strconv.Itoa(x), "value": strconv.Itoa(value), "branch": branch,
			"boundary": strconv.Itoa(boundary),
		},
		KPOrdinal: kp,
	}
}

// formatLine renders y = mx + b without a leading plus sign.
func formatLine(m, b int) string {
	var coef string
	switch m {
	case 1:
		coef = "x"
	case -1:
		coef = "-x"
	default:
		coef = fmt.Sprintf("%dx", m)
	}
	if b == 0 {
		return coef
	}
	if b < 0 {
		return fmt.Sprintf("%s - %d", coef, -b)
	}
	return fmt.Sprintf("%s + %d", coef, b)
}

// ---------------------------------------------------------------------------
// Solving systems by graphing
// ---------------------------------------------------------------------------

type systemsGraphingGen struct{}

func (g *systemsGraphingGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = rand.Intn(3) + 1
	}

	x0 := rand.Intn(9) - 4
	y0 := rand.Intn(9) - 4
	m1 := rand.Intn(4) + 1
	m2 := -rand.Intn(4) - 1
	if rand.Intn(2) == 0 {
		m1, m2 = -m1, -m2
	}
	b1 := y0 - m1*x0
	b2 := y0 - m2*x0

	var answer int
	var asked string
	switch kp {
	case 1:
		answer = x0
		asked = "x-coordinate"
	case 2:
		answer = y0
		asked = "y-coordinate"
	default:
		answer = x0 + y0
		asked = "sum of the coordinates"
	}

	question := fmt.Sprintf("The lines y = %s and y = %s intersect at one point. What is the %s of the intersection? (enter a number)",
		formatLine(m1, b1), formatLine(m2, b2), asked)
	explanation := fmt.Sprintf("Set the expressions equal: %s = %s gives x = %d. Substituting back gives y = %d, so the point is (%d, %d).",
		formatLine(m1, b1), formatLine(m2, b2), x0, y0, x0, y0)

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts: map[string]string{
			"x": strconv.Itoa(x0), "y": strconv.Itoa(y0),
			"sum": strconv.Itoa(x0 + y0),
			"m1":  strconv.Itoa(m1), "b1": strconv.Itoa(b1),
			"m2": strconv.Itoa(m2), "b2": strconv.Itoa(b2),
		},
		KPOrdinal: kp,
	}
}

// ---------------------------------------------------------------------------
// Cramer's rule
// ---------------------------------------------------------------------------

type cramerGen struct{}

func (g *cramerGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = rand.Intn(3) + 1
	}

	x0 := rand.Intn(9) - 4
	y0 := rand.Intn(9) - 4
	var a, b, c, d int
	for {
		a = rand.Intn(7) - 3
		b = rand.Intn(7) - 3
		c = rand.Intn(7) - 3
		d = rand.Intn(7) - 3
		if a*d-b*c != 0 && a != 0 && c != 0 {
			break
		}
	}
	e := a*x0 + b*y0
	f := c*x0 + d*y0

	var answer int
	var asked string
	switch kp {
	case 1:
		answer = a*d - b*c
		asked = "determinant of the coefficient matrix"
	case 2:
		answer = x0
		asked = "value of x"
	default:
		answer = y0
		asked = "value of y"
	}

	question := fmt.Sprintf("For the system %s and %s, what is the %s? (enter a number)",
		formatEquation(a, "x", b, "y", e), formatEquation(c, "x", d, "y", f), asked)
	explanation := fmt.Sprintf("The coefficient determinant is (%d)(%d) - (%d)(%d) = %d. Solving gives x = %d and y = %d.",
		a, d, b, c, a*d-b*c, x0, y0)

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts: map[string]string{
			"a": strconv.Itoa(a), "b": strconv.Itoa(b), "c": strconv.Itoa(c), "d": strconv.Itoa(d),
			"e": strconv.Itoa(e), "f": strconv.Itoa(f),
			"det": strconv.Itoa(a*d - b*c), "x": strconv.Itoa(x0), "y": strconv.Itoa(y0),
		},
		KPOrdinal: kp,
	}
}

// formatEquation renders ax + by = c with readable signs.
func formatEquation(a int, xv string, b int, yv string, c int) string {
	var out string
	switch a {
	case 1:
		out = xv
	case -1:
		out = "-" + xv
	default:
		out = fmt.Sprintf("%d%s", a, xv)
	}
	if b != 0 {
		mag := b
		switch mag {
		case 1:
			out += " + " + yv
		case -1:
			out += " - " + yv
		default:
			if mag < 0 {
				out += fmt.Sprintf(" - %d%s", -mag, yv)
			} else {
				out += fmt.Sprintf(" + %d%s", mag, yv)
			}
		}
	}
	return fmt.Sprintf("%s = %d", out, c)
}

// formatFactor renders a divisor of the form x - a without a double sign.
func formatFactor(a int) string {
	switch {
	case a == 0:
		return "x"
	case a > 0:
		return fmt.Sprintf("x - %d", a)
	default:
		return fmt.Sprintf("x + %d", -a)
	}
}

// ---------------------------------------------------------------------------
// The remainder theorem
// ---------------------------------------------------------------------------

type polyRemainderGen struct{}

func (g *polyRemainderGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = rand.Intn(3) + 1
	}

	degree := 2 + rand.Intn(3) // 2..4
	terms := make([]polyTerm, 0, degree+1)
	coefs := make([]int, degree+1)
	for d := degree; d >= 0; d-- {
		c := rand.Intn(7) - 3
		if d == degree && c == 0 {
			c = 2
		}
		coefs[d] = c
		terms = append(terms, polyTerm{c, d})
	}
	poly := formatPoly(terms...)
	a := rand.Intn(9) - 4
	remainder := 0
	for d, c := range coefs {
		remainder += c * powInt(a, d)
	}

	var question string
	divisor := formatFactor(a)
	switch kp {
	case 1:
		question = fmt.Sprintf("What is the remainder when f(x) = %s is divided by %s? (enter a number)", poly, divisor)
	case 2:
		question = fmt.Sprintf("For f(x) = %s, what is f(%d)? (enter a number)", poly, a)
	default:
		question = fmt.Sprintf("By the remainder theorem, dividing f(x) = %s by %s leaves what remainder? (enter a number)", poly, divisor)
	}
	explanation := fmt.Sprintf("The remainder theorem says the remainder is f(%d). Substituting gives %d.", a, remainder)

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(remainder),
		Explanation: explanation,
		Facts: map[string]string{
			"a": strconv.Itoa(a), "poly": poly, "remainder": strconv.Itoa(remainder),
		},
		KPOrdinal: kp,
	}
}
