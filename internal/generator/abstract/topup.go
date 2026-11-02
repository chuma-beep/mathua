package abstract

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

type cauchyGroupGen struct{}

func cauchyPrimes(n int) []int {
	out := []int{}
	for p := 2; p <= n; p++ {
		if n%p != 0 {
			continue
		}
		prime := true
		for d := 2; d*d <= p; d++ {
			if p%d == 0 {
				prime = false
				break
			}
		}
		if prime {
			out = append(out, p)
			for n%p == 0 {
				n /= p
			}
		}
	}
	return out
}

func (g *cauchyGroupGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := abstractBridgeOrdinal(ctx)
	switch o {
	case 1:
		p := cauchyPrimes(13 + rand.Intn(30))
		prime := p[rand.Intn(len(p))]
		ans := prime - 1
		return generator.Problem{
			Question:    fmt.Sprintf("A cyclic group has order %d. How many of its elements have order %d? (enter a number)", prime, prime),
			Answer:      strconv.Itoa(ans),
			Explanation: fmt.Sprintf("A cyclic group of prime order %d has %d non-identity elements, and each has order %d.", prime, ans, prime),
			KPOrdinal:   o,
			Facts:       map[string]string{"answer": strconv.Itoa(ans)},
		}
	case 2:
		n := 6 + rand.Intn(55)
		ps := cauchyPrimes(n)
		ans := len(ps)
		return generator.Problem{
			Question:    fmt.Sprintf("A group has order %d. Cauchy's theorem guarantees an element of prime order for each distinct prime divisor. How many such prime orders are guaranteed? (enter a number)", n),
			Answer:      strconv.Itoa(ans),
			Explanation: fmt.Sprintf("The distinct prime divisors of %d are %v, so %d prime orders are guaranteed.", n, ps, ans),
			KPOrdinal:   o,
			Facts:       map[string]string{"answer": strconv.Itoa(ans)},
		}
	default:
		n := 6 + rand.Intn(55)
		ps := cauchyPrimes(n)
		ans := ps[len(ps)-1]
		return generator.Problem{
			Question:    fmt.Sprintf("A group has order %d. What is the largest prime order that Cauchy's theorem guarantees? (enter a number)", n),
			Answer:      strconv.Itoa(ans),
			Explanation: fmt.Sprintf("The prime divisors of %d are %v, and the largest is %d.", n, ps, ans),
			KPOrdinal:   o,
			Facts:       map[string]string{"answer": strconv.Itoa(ans)},
		}
	}
}
