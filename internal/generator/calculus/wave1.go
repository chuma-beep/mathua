package calculus

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

// ---------------------------------------------------------------------------
// Limits of piecewise functions
// ---------------------------------------------------------------------------

type limitPiecewiseGen struct{}

func (g *limitPiecewiseGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = rand.Intn(3) + 1
	}

	boundary := rand.Intn(5) + 1
	m1 := rand.Intn(4) + 1
	m2 := rand.Intn(4) + 1
	b1 := rand.Intn(11) - 5
	leftValue := m1*boundary + b1

	var question string
	var answer int
	var explanation string
	switch kp {
	case 1:
		answer = leftValue
		question = fmt.Sprintf("For f(x) = %s when x < %d, and f(x) = %s when x ≥ %d, what is the left-hand limit as x approaches %d? (enter a number)",
			formatCalcLine(m1, b1), boundary, formatCalcLine(m2, b1+3), boundary, boundary)
		explanation = fmt.Sprintf("Below x = %d the rule is %s, so the left-hand limit is %d.", boundary, formatCalcLine(m1, b1), leftValue)
	case 2:
		answer = m2*boundary + b1 + 3
		question = fmt.Sprintf("For f(x) = %s when x < %d, and f(x) = %s when x ≥ %d, what is the right-hand limit as x approaches %d? (enter a number)",
			formatCalcLine(m1, b1), boundary, formatCalcLine(m2, b1+3), boundary, boundary)
		explanation = fmt.Sprintf("Above x = %d the rule is %s, so the right-hand limit is %d.", boundary, formatCalcLine(m2, b1+3), answer)
	default:
		answer = leftValue - m2*boundary
		question = fmt.Sprintf("For f(x) = %s when x < %d, and f(x) = %s + c when x ≥ %d, what value of c makes f continuous at x = %d? (enter a number)",
			formatCalcLine(m1, b1), boundary, formatCalcLine(m2, 0), boundary, boundary)
		explanation = fmt.Sprintf("The left limit is %d. The right branch is %s + c, so continuity needs c = %d - %d = %d.",
			leftValue, formatCalcLine(m2, 0), leftValue, m2*boundary, answer)
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts: map[string]string{
			"boundary": strconv.Itoa(boundary), "left": strconv.Itoa(leftValue),
			"right": strconv.Itoa(m2*boundary + b1 + 3), "answer": strconv.Itoa(answer),
		},
		KPOrdinal: kp,
	}
}

func formatCalcLine(m, b int) string {
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
// Intervals of continuity
// ---------------------------------------------------------------------------

type continuityIntervalGen struct{}

func (g *continuityIntervalGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = rand.Intn(3) + 1
	}

	var question, explanation string
	var answer int
	switch kp {
	case 1:
		a := rand.Intn(9) + 1
		answer = a
		question = fmt.Sprintf("The function f(x) = 1/(x - %d) is continuous everywhere except at one point. At what x is it discontinuous? (enter a number)", a)
		explanation = fmt.Sprintf("The denominator x - %d is zero at x = %d, so the function is undefined and discontinuous there.", a, a)

	case 2:
		a := rand.Intn(9) + 1
		answer = -a
		question = fmt.Sprintf("The function f(x) = √(x + %d) is continuous on an interval [c, ∞). What is the left endpoint c? (enter a number)", a)
		explanation = fmt.Sprintf("The radicand x + %d is nonnegative when x ≥ %d, so the interval starts at %d.", a, answer, answer)

	default:
		a := rand.Intn(5) + 1
		b := a + rand.Intn(5) + 1
		answer = b
		question = fmt.Sprintf("The function f(x) = 1/((x - %d)(x - %d)) has two points of discontinuity. What is the larger x? (enter a number)", a, b)
		explanation = fmt.Sprintf("The denominator is zero at x = %d and x = %d, so the larger point of discontinuity is %d.", a, b, b)
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts:       map[string]string{"answer": strconv.Itoa(answer)},
		KPOrdinal:   kp,
	}
}
