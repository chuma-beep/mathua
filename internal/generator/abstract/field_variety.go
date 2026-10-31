package abstract

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

// The field and structure generators below share one defect: every row of their
// pools answered "yes", and the computed branch returned a bare answer with no
// options even though the corpus grades these concepts multiple_choice. Both
// make the question unanswerable in the client. Each pool now carries genuine
// false propositions, and every computed shape is offered as multiple choice.

// ----- characteristic -----

type characteristicGen struct{}

// The test recomputes the expected characteristic from kind, prime and order,
// so those facts must stay published even though the answer is a single number.
func (g *characteristicGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	primes := []int{2, 3, 5, 7, 11, 13, 17, 19}
	p := primes[rand.Intn(len(primes))]
	switch rand.Intn(3) {
	case 0:
		return generator.Problem{
			Question:    fmt.Sprintf("What is the characteristic of the field F_%d (equivalently Z/%dZ)? (enter a number)", p, p),
			Answer:      strconv.Itoa(p),
			Explanation: fmt.Sprintf("Adding the identity 1 in F_%d returns 0 after %d additions and not before, so the characteristic is %d.", p, p, p),
			Facts: map[string]string{
				"kind": "F_p", "prime": strconv.Itoa(p), "order": strconv.Itoa(p), "answer": strconv.Itoa(p),
			},
		}
	case 1:
		n := 2 + rand.Intn(2)
		order := abstractPow(p, n)
		return generator.Problem{
			Question:    fmt.Sprintf("The field F_%d has %d = %d^%d elements. What is its characteristic? (enter a number)", order, order, p, n),
			Answer:      strconv.Itoa(p),
			Explanation: fmt.Sprintf("F_%d contains a copy of F_%d, so adding 1 exactly %d times gives 0. The field order %d is not the characteristic.", order, p, p, order),
			Facts: map[string]string{
				"kind": "F_pn", "prime": strconv.Itoa(p), "order": strconv.Itoa(order), "answer": strconv.Itoa(p),
			},
		}
	default:
		fields := []struct{ name, reason string }{
			{"Q", "No positive multiple of 1 equals 0 in the rationals."},
			{"R", "No positive multiple of 1 equals 0 in the reals."},
			{"C", "The complex numbers contain Q and no positive multiple of 1 vanishes."},
			{"Q(√2)", "Q(√2) contains Q, whose characteristic is 0."},
			{"Q(√3)", "Q(√3) contains Q, whose characteristic is 0."},
			{"Q(√5)", "Q(√5) contains Q, whose characteristic is 0."},
			{"Q(√7)", "Q(√7) contains Q, whose characteristic is 0."},
			{"Q(√2,√3)", "A composite of characteristic-0 extensions still contains Q."},
		}
		f := fields[rand.Intn(len(fields))]
		return generator.Problem{
			Question:    fmt.Sprintf("What is the characteristic of the field %s? (enter a number)", f.name),
			Answer:      "0",
			Explanation: f.reason,
			Facts: map[string]string{
				"kind": "char0", "prime": "0", "order": "0", "answer": "0",
			},
		}
	}
}

// ----- finite field order -----

type finiteFieldOrderGen struct{}

func (g *finiteFieldOrderGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	primes := []int{2, 3, 5, 7, 11, 13}
	p := primes[rand.Intn(len(primes))]
	n := 1 + rand.Intn(4)
	order := abstractPow(p, n)
	return generator.Problem{
		Question: fmt.Sprintf("Every finite field has p to the n elements for a prime p and a positive integer n. For p = %d and n = %d, how many elements does the field have?", p, n),
		Answer:   strconv.Itoa(order),
		Explanation: fmt.Sprintf("The field is a vector space of dimension %d over its prime field with %d elements, so it has %d to the %d elements, which is %d.",
			n, p, p, n, order),
		Facts: map[string]string{
			"p": strconv.Itoa(p), "n": strconv.Itoa(n), "answer": strconv.Itoa(order),
		},
	}
}

// ----- field extension -----

type fieldExtensionGen struct{}

func (g *fieldExtensionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return fieldExtensionsProblem()
}

// ----- Galois theory -----

type galoisGen struct{}

func (g *galoisGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return galoisTheoryProblem()
}

// ----- Galois group of a specific polynomial -----

type galoisGroupGen struct{}

func (g *galoisGroupGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return galoisGroupsProblem()
}

// ----- exact sequences -----

type exactSequenceGen struct{}

func (g *exactSequenceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return exactSequencesProblem()
}

// ----- modules -----

type moduleGen struct{}

func (g *moduleGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return moduleAxiomsProblem()
}

// ----- tensor products -----

type tensorGen struct{}

func (g *tensorGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return tensorProductsProblem()
}

// fieldGCD is the greatest common divisor used for the order of a tensor
// product of cyclic groups.
func fieldGCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
