package probability

import (
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

type varianceGen struct{}

func (g *varianceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	n := 2 + rand.Intn(2)
	values := make([]int, n)
	weights := make([]int, n)
	for i := 0; i < n; i++ {
		values[i] = rand.Intn(6)
	}
	for i := 0; i < n-1; i++ {
		weights[i] = 1 + rand.Intn(8)
	}
	// The last weight closes the total to a round number, which keeps the
	// arithmetic exact in integers.
	used := 0
	for i := 0; i < n-1; i++ {
		used += weights[i]
	}
	weights[n-1] = 10 - used%10
	if weights[n-1] <= 0 {
		weights[n-1] += 10
	}
	den := 0
	for _, w := range weights {
		den += w
	}

	sw, sw2 := 0, 0
	for i := 0; i < n; i++ {
		sw += values[i] * weights[i]
		sw2 += values[i] * values[i] * weights[i]
	}
	num := den*sw2 - sw*sw
	answer := decimal6(float64(num) / float64(den*den))

	vals := make([]string, n)
	wts := make([]string, n)
	for i := 0; i < n; i++ {
		vals[i] = strconv.Itoa(values[i])
		wts[i] = strconv.Itoa(weights[i])
	}
	return generator.Problem{
		Question: "A discrete random variable X takes the values " + strings.Join(vals, ", ") +
			" with probabilities proportional to the weights " + strings.Join(wts, ", ") +
			" (the weights sum to " + strconv.Itoa(den) + "). What is Var(X) = E[X²] − (E[X])²? (enter a decimal)",
		Answer: answer,
		Explanation: "E[X] is the weighted average " + strconv.Itoa(sw) + "/" + strconv.Itoa(den) +
			", and E[X²] is " + strconv.Itoa(sw2) + "/" + strconv.Itoa(den) +
			". Subtracting the square of the mean gives " + answer + ".",
		Facts: map[string]string{
			"values":  strings.Join(vals, ","),
			"weights": strings.Join(wts, ","),
			"den":     strconv.Itoa(den),
			"answer":  answer,
		},
	}
}
