package precalculus

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func precalcBridgeOrdinal(ctx generator.GeneratorContext) int {
	if ctx.KPOrdinal >= 1 && ctx.KPOrdinal <= 3 {
		return ctx.KPOrdinal
	}
	return 1 + rand.Intn(3)
}

func averageRateProblem(ordinal int, question, answer, explanation string, facts map[string]string) generator.Problem {
	if facts == nil {
		facts = map[string]string{}
	}
	facts["answer"] = answer
	return generator.Problem{
		Question:    question,
		Answer:      answer,
		Explanation: explanation,
		Facts:       facts,
		KPOrdinal:   ordinal,
	}
}

type averageRateGen struct{}

func (g *averageRateGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	ordinal := precalcBridgeOrdinal(ctx)
	switch ordinal {
	case 1:
		m := 1 + rand.Intn(9)
		c := rand.Intn(13) - 6
		a := rand.Intn(5)
		b := a + 1 + rand.Intn(5)
		fa := m*a + c
		fb := m*b + c
		return averageRateProblem(ordinal,
			fmt.Sprintf("For f(x) = %dx + %d, what is the average rate of change on [%d, %d]? (enter a number)", m, c, a, b),
			strconv.Itoa(m),
			fmt.Sprintf("f(%d) = %d and f(%d) = %d, so the average rate is (%d - %d)/(%d - %d) = %d.", a, fa, b, fb, fb, fa, b, a, m),
			map[string]string{"kind": "linear", "m": strconv.Itoa(m), "a": strconv.Itoa(a), "b": strconv.Itoa(b)},
		)
	case 2:
		c := rand.Intn(11) - 5
		a := rand.Intn(5)
		b := a + 1 + rand.Intn(5)
		fa := a*a + c
		fb := b*b + c
		answer := a + b
		return averageRateProblem(ordinal,
			fmt.Sprintf("For f(x) = x^2 + %d, what is the average rate of change on [%d, %d]? (enter a number)", c, a, b),
			strconv.Itoa(answer),
			fmt.Sprintf("f(%d) = %d and f(%d) = %d, so the average rate is (%d - %d)/(%d - %d) = %d.", a, fa, b, fb, fb, fa, b, a, answer),
			map[string]string{"kind": "quadratic", "a": strconv.Itoa(a), "b": strconv.Itoa(b)},
		)
	default:
		c := rand.Intn(9) - 4
		a := rand.Intn(4)
		b := a + 1 + rand.Intn(4)
		fa := a*a*a + c
		fb := b*b*b + c
		answer := a*a + a*b + b*b
		return averageRateProblem(ordinal,
			fmt.Sprintf("For f(x) = x^3 + %d, what is the average rate of change on [%d, %d]? (enter a number)", c, a, b),
			strconv.Itoa(answer),
			fmt.Sprintf("f(%d) = %d and f(%d) = %d. The quotient simplifies to a^2 + ab + b^2 = %d.", a, fa, b, fb, answer),
			map[string]string{"kind": "cubic", "a": strconv.Itoa(a), "b": strconv.Itoa(b)},
		)
	}
}
