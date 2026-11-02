package odes

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

type forcedOscillationGen struct{}

func (g *forcedOscillationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := odeBridgeOrdinal(ctx)
	s := 2 + rand.Intn(5)
	m := 1 + rand.Intn(4)
	switch o {
	case 1:
		return generator.Problem{
			Question:    fmt.Sprintf("An undamped oscillator has mass %d and stiffness %d. What is its natural frequency? (enter a number)", m, m*s*s),
			Answer:      strconv.Itoa(s),
			Explanation: fmt.Sprintf("The natural frequency is the square root of k/m, which is the square root of %d, equal to %d.", s*s, s),
			KPOrdinal:   o,
			Facts:       map[string]string{"answer": strconv.Itoa(s)},
		}
	case 2:
		ans := 2 / float64(s)
		return generator.Problem{
			Question:    fmt.Sprintf("An oscillator has natural frequency %d. What is the coefficient of pi in its natural period? (enter a decimal)", s),
			Answer:      odeDecimal(ans),
			Explanation: fmt.Sprintf("The period is 2pi divided by %d, so the coefficient of pi is %s.", s, odeDecimal(ans)),
			KPOrdinal:   o,
			Facts:       map[string]string{"answer": odeDecimal(ans)},
		}
	default:
		ans := 2 * s
		return generator.Problem{
			Question:    fmt.Sprintf("At resonance with natural frequency %d, the particular solution contains t sin(%d t) divided by a constant. What is that constant? (enter a number)", s, s),
			Answer:      strconv.Itoa(ans),
			Explanation: fmt.Sprintf("The resonant particular solution carries the factor 1/(2 omega), so the denominator is 2 times %d equals %d.", s, ans),
			KPOrdinal:   o,
			Facts:       map[string]string{"answer": strconv.Itoa(ans)},
		}
	}
}
