package statistics

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

func statTopupDecimal(v float64) string {
	s := strconv.FormatFloat(v, 'f', 4, 64)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	if s == "-0" || s == "" {
		s = "0"
	}
	return s
}

type logisticRegressionGen struct{}

func (g *logisticRegressionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := bridgeOrdinal(ctx)
	ps := []float64{0.1, 0.2, 0.25, 0.5, 0.75, 0.8, 0.9}
	p := ps[rand.Intn(len(ps))]
	switch o {
	case 1:
		ans := math.Log(p / (1 - p))
		return generator.Problem{
			Question:    fmt.Sprintf("A binary outcome has success probability %s. What is its log-odds? (enter a decimal)", statTopupDecimal(p)),
			Answer:      statTopupDecimal(ans),
			Explanation: fmt.Sprintf("The odds are p/(1-p) = %s, and the natural logarithm of that value is %s.", statTopupDecimal(p/(1-p)), statTopupDecimal(ans)),
			KPOrdinal:   o,
			Facts:       map[string]string{"answer": statTopupDecimal(ans)},
		}
	case 2:
		z := float64(-2 + rand.Intn(5))
		ans := 1 / (1 + math.Exp(-z))
		return generator.Problem{
			Question:    fmt.Sprintf("A logistic model produces the linear score %s. What probability does the sigmoid give? (enter a decimal)", statTopupDecimal(z)),
			Answer:      statTopupDecimal(ans),
			Explanation: fmt.Sprintf("The sigmoid 1/(1+e^(-z)) at z = %s is %s.", statTopupDecimal(z), statTopupDecimal(ans)),
			KPOrdinal:   o,
			Facts:       map[string]string{"answer": statTopupDecimal(ans)},
		}
	default:
		ans := p / (1 - p)
		return generator.Problem{
			Question:    fmt.Sprintf("A binary outcome has success probability %s. What are the odds of success? (enter a decimal)", statTopupDecimal(p)),
			Answer:      statTopupDecimal(ans),
			Explanation: fmt.Sprintf("The odds are p/(1-p): %s divided by %s equals %s.", statTopupDecimal(p), statTopupDecimal(1-p), statTopupDecimal(ans)),
			KPOrdinal:   o,
			Facts:       map[string]string{"answer": statTopupDecimal(ans)},
		}
	}
}
