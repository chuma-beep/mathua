package statistics

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

// statTrim renders a chi-square value without trailing zeros, so an exact zero
// reads 0 rather than 0.000000.
func statTrim(v float64) string {
	return strconv.FormatFloat(math.Round(v*1e6)/1e6, 'f', -1, 64)
}

type sampleSpaceGen struct{}

func (g *sampleSpaceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch rand.Intn(5) {
	case 0:
		sides := 4 + rand.Intn(17)
		return generator.Problem{
			Question:    fmt.Sprintf("A die has %d faces. How many outcomes are in the sample space when rolling it once?", sides),
			Answer:      strconv.Itoa(sides),
			Explanation: fmt.Sprintf("Each face is one outcome, so the sample space has %d elements.", sides),
		}
	case 1:
		flips := 1 + rand.Intn(5)
		n := 1 << flips
		return generator.Problem{
			Question:    fmt.Sprintf("A coin is flipped %d times. How many outcomes are in the sample space?", flips),
			Answer:      strconv.Itoa(n),
			Explanation: fmt.Sprintf("Each flip doubles the count, so %d flips give 2^%d = %d ordered outcomes.", flips, flips, n),
		}
	case 2:
		sides := 4 + rand.Intn(9)
		n := 2 * sides
		return generator.Problem{
			Question:    fmt.Sprintf("A die with %d sides is rolled and a coin is flipped. How many outcomes are in the sample space?", sides),
			Answer:      strconv.Itoa(n),
			Explanation: fmt.Sprintf("The die gives %d outcomes and the coin gives 2, so by the product rule there are %d.", sides, n),
		}
	case 3:
		suits := 2 + rand.Intn(3)
		ranks := 6 + rand.Intn(8)
		n := suits * ranks
		return generator.Problem{
			Question:    fmt.Sprintf("A deck has %d suits and %d ranks. How many distinct cards are there?", suits, ranks),
			Answer:      strconv.Itoa(n),
			Explanation: fmt.Sprintf("Every suit carries every rank, so the deck holds %d times %d = %d cards.", suits, ranks, n),
		}
	default:
		sections := 3 + rand.Intn(10)
		return generator.Problem{
			Question:    fmt.Sprintf("A spinner has %d equally likely sections. How many outcomes are in the sample space of one spin?", sections),
			Answer:      strconv.Itoa(sections),
			Explanation: fmt.Sprintf("Each section is one outcome, so there are %d.", sections),
		}
	}
}

type chiSquareTestGen struct{}

func (g *chiSquareTestGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	k := 2 + rand.Intn(4)
	expected := 10
	observed := make([]int, k)
	total := 0
	for i := 0; i < k-1; i++ {
		d := rand.Intn(9) - 4
		if d == 0 {
			d = 3
		}
		observed[i] = expected + d
		total += observed[i]
	}
	observed[k-1] = expected*k - total

	chi := 0.0
	parts := make([]string, k)
	for i := 0; i < k; i++ {
		diff := observed[i] - expected
		chi += float64(diff*diff) / float64(expected)
		parts[i] = fmt.Sprintf("%d", diff*diff)
	}
	obs := ""
	for i, o := range observed {
		if i > 0 {
			obs += ", "
		}
		obs += strconv.Itoa(o)
	}
	return generator.Problem{
		Question: fmt.Sprintf("A goodness-of-fit test compares observed counts %s with %d expected counts of %d each. What is chi-square = sum (O - E)^2/E?", obs, k, expected),
		Answer:   statTrim(chi),
		Explanation: fmt.Sprintf(
			"Subtract the expected count from each observed count, square it and divide by %d. The squared differences are %s, so the sum is %s.",
			expected, joinStrings(parts), statTrim(chi)),
		Facts: map[string]string{
			"k":        strconv.Itoa(k),
			"expected": repeatInt(expected, k),
			"observed": obs,
			"answer":   statTrim(chi),
		},
	}
}

// repeatInt renders the same count k times, which is what the batch test reads
// to recompute chi-square independently.
func repeatInt(v, k int) string {
	out := ""
	for i := 0; i < k; i++ {
		if i > 0 {
			out += ","
		}
		out += strconv.Itoa(v)
	}
	return out
}

// joinStrings joins with commas for the explanation text.
func joinStrings(parts []string) string {
	out := ""
	for i, p := range parts {
		if i > 0 {
			out += ", "
		}
		out += p
	}
	return out
}
