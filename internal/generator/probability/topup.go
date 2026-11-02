package probability

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/chuma-beep/mathua/internal/generator"
)

type geometricDistributionGen struct{}

func (g *geometricDistributionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := probBridgeOrdinal(ctx)
	ps := []float64{0.05, 0.1, 0.15, 0.2, 0.25, 0.4, 0.5, 0.75}
	p := ps[rand.Intn(len(ps))]
	switch o {
	case 1:
		k := 1 + rand.Intn(5)
		ans := math.Pow(1-p, float64(k-1)) * p
		return generator.Problem{
			Question:    fmt.Sprintf("A geometric random variable has success probability %s. What is P(X = %d)? (enter a decimal)", probDecimal4(p), k),
			Answer:      probDecimal4(ans),
			Explanation: fmt.Sprintf("The first %d trials fail and the next succeeds: (%s)^%d times %s equals %s.", k-1, probDecimal4(1-p), k-1, probDecimal4(p), probDecimal4(ans)),
			KPOrdinal:   o,
			Facts:       map[string]string{"answer": probDecimal4(ans)},
		}
	case 2:
		ans := 1 / p
		return generator.Problem{
			Question:    fmt.Sprintf("A geometric random variable has success probability %s. What is its mean? (enter a decimal)", probDecimal4(p)),
			Answer:      probDecimal4(ans),
			Explanation: fmt.Sprintf("The mean of a geometric count is 1/p, so 1 divided by %s equals %s.", probDecimal4(p), probDecimal4(ans)),
			KPOrdinal:   o,
			Facts:       map[string]string{"answer": probDecimal4(ans)},
		}
	default:
		ans := (1 - p) / (p * p)
		return generator.Problem{
			Question:    fmt.Sprintf("A geometric random variable has success probability %s. What is its variance? (enter a decimal)", probDecimal4(p)),
			Answer:      probDecimal4(ans),
			Explanation: fmt.Sprintf("The variance is (1-p)/p^2: %s divided by %s squared equals %s.", probDecimal4(1-p), probDecimal4(p), probDecimal4(ans)),
			KPOrdinal:   o,
			Facts:       map[string]string{"answer": probDecimal4(ans)},
		}
	}
}
