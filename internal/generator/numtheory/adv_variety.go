package numtheory

import (
	"github.com/chuma-beep/mathua/internal/generator"
)

// These concepts are graded multiple choice, so every shape carries options.
// The conceptual rows keep the theorem statements the concept is about, and the
// computed rows ask for a small quantity the learner must actually work out.

func advChoice(question, correct, explanation string, distractors ...string) generator.Problem {
	p := generator.MultipleChoice(question, correct, distractors...)
	p.Explanation = explanation
	return p
}

type chebyshevGen struct{}

func (g *chebyshevGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return chebyshevProblem()
}

type continuedFractionsGen struct{}

func (g *continuedFractionsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return continuedFractionsProblem()
}

type dirichletGen struct{}

func (g *dirichletGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return dirichletProblem()
}

type ellipticCurveGen struct{}

func (g *ellipticCurveGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return ellipticCurveProblem()
}

func otherAnswer(a string) string {
	if a == "yes" {
		return "no"
	}
	return "yes"
}

func primePowerCount(x int) int {
	count := 0
	for n := 2; n <= x; n++ {
		if isPrimePower(n) {
			count++
		}
	}
	return count
}

func isPrimePower(n int) bool {
	for p := 2; p*p <= n; p++ {
		if n%p == 0 {
			m := n
			for m%p == 0 {
				m /= p
			}
			return m == 1
		}
	}
	return n >= 2
}

func smallestPrimeInAP(a, m int) int {
	for n := a; ; n += m {
		if n < 2 {
			continue
		}
		if isPrimePower(n) && isPrimeNumber(n) {
			return n
		}
	}
}

func isPrimeNumber(n int) bool {
	if n < 2 {
		return false
	}
	for d := 2; d*d <= n; d++ {
		if n%d == 0 {
			return false
		}
	}
	return true
}
