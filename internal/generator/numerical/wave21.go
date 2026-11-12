package numerical

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

// Wave 21 adds one numerical-analysis concept, Anderson acceleration. The
// answers count history depth, damping and saved iterations.

func registerWave21(reg *generator.Registry) {
	reg.Register("numa.optimization.anderson", &w21AndersonGen{})
}

func w21N(lo, hi int) int { return lo + rand.Intn(hi-lo+1) }

func w21P(o int, q, a, ex string, facts map[string]string) generator.Problem {
	facts["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: facts, KPOrdinal: o}
}

type w21AndersonGen struct{}

func (g *w21AndersonGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch kp {
	case 1:
		m := w21N(1, 30)
		return w21P(kp,
			fmt.Sprintf("Anderson acceleration stores the last %d residuals. How many residuals are kept? (enter a number)", m),
			strconv.Itoa(m),
			fmt.Sprintf("The history depth is the number of stored residuals, so %d are kept.", m),
			map[string]string{"m": strconv.Itoa(m)})
	case 2:
		beta := 10 * w21N(1, 9)
		plain := 100 - beta
		return w21P(kp,
			fmt.Sprintf("A damping factor of %d percent is applied to the accelerated iterate. What percentage weight does the plain iterate receive? (enter a number)", beta),
			strconv.Itoa(plain),
			fmt.Sprintf("The two weights sum to 100, so the plain iterate receives %d percent.", plain),
			map[string]string{"beta": strconv.Itoa(beta)})
	default:
		n := w21N(10, 100)
		s := w21N(20, 60)
		saved := (n*s + 99) / 100
		return w21P(kp,
			fmt.Sprintf("Plain iteration needs %d steps. Acceleration removes %d percent of them. How many steps are saved? (round up, enter a number)", n, s),
			strconv.Itoa(saved),
			fmt.Sprintf("%d percent of %d is %d.%02d, which rounds up to %d steps.", s, n, n*s/100, n*s%100, saved),
			map[string]string{"n": strconv.Itoa(n), "s": strconv.Itoa(s)})
	}
}
