package algebra

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

// turningPointsGen asks about the relationship between a polynomial's degree
// and the number of turning points. Every family publishes the same two facts,
// degree and turning_points, so one written solution covers all three.
type turningPointsGen struct{}

func (g *turningPointsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = rand.Intn(3) + 1
	}

	var question, explanation string
	var degree, turning, answer int

	switch kp {
	case 1:
		degree = rand.Intn(8) + 2 // 2..9
		turning = degree - 1
		answer = turning
		question = fmt.Sprintf("A polynomial has degree %d. What is the greatest possible number of turning points? (enter a number)", degree)
		explanation = fmt.Sprintf("The derivative of a degree-%d polynomial has degree %d, so \\(f'(x) = 0\\) has at most %d real %s. Each simple solution is a turning point, so there are at most %d turning %s.", degree, degree-1, degree-1, pluralWord(degree-1, "solution", "solutions"), turning, pluralWord(turning, "point", "points"))

	case 2:
		degree = rand.Intn(5) + 2 // 2..6
		lead := rand.Intn(8) - 4
		if lead == 0 {
			lead = 3
		}
		terms := []polyTerm{{lead, degree}}
		for d := degree - 1; d >= 0; d-- {
			terms = append(terms, polyTerm{rand.Intn(9) - 4, d})
		}
		poly := formatPoly(terms...)
		turning = degree - 1
		answer = turning
		question = fmt.Sprintf("What is the greatest possible number of turning points of \\(f(x) = %s\\)? (enter a number)", poly)
		explanation = fmt.Sprintf("The degree is %d, so the derivative has degree %d and \\(f'(x) = 0\\) has at most %d real %s. The lower-order terms cannot raise that ceiling.", degree, degree-1, turning, pluralWord(turning, "solution", "solutions"))

	default:
		turning = rand.Intn(7) + 1 // 1..7
		intercept := rand.Intn(11) - 5
		degree = turning + 1
		answer = degree
		question = fmt.Sprintf("A polynomial has exactly %d turning points and its graph passes through \\((0, %d)\\). What is the least possible degree? (enter a number)", turning, intercept)
		explanation = fmt.Sprintf("A polynomial with %d turning points needs a derivative with at least %d real %s, so its degree is at least %d. The y-intercept does not affect the degree.", turning, turning, pluralWord(turning, "zero", "zeros"), degree)
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts: map[string]string{
			"degree":            strconv.Itoa(degree),
			"derivative_degree": strconv.Itoa(degree - 1),
			"turning_points":    strconv.Itoa(turning),
		},
		KPOrdinal: kp,
	}
}

// pluralWord picks the singular or plural noun for a count.
func pluralWord(n int, one, many string) string {
	if n == 1 {
		return one
	}
	return many
}
