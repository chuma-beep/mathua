package abstract

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

// groupChoice builds a multiple choice problem with an explanation. The
// generator package's helper returns no explanation, and a batch test fails on
// an empty one.
func groupChoice(question, correct, explanation string, distractors ...string) generator.Problem {
	p := generator.MultipleChoice(question, correct, distractors...)
	p.Explanation = explanation
	return p
}

// otherAnswer flips a yes/no answer so a pool always contains both.
func otherAnswer(a string) string {
	if a == "yes" {
		return "no"
	}
	return "yes"
}

// groupYesNo turns a yes/no row into a multiple choice question. The concept is
// graded multiple_choice, so a bare "yes" with no options cannot be answered in
// the client.
func groupYesNo(question, answer, explanation string) generator.Problem {
	clean := strings.ReplaceAll(question, " (yes/no)", "")
	clean = strings.ReplaceAll(clean, "(yes/no)", "")
	return groupChoice(clean, answer, explanation, otherAnswer(answer))
}

// ----- p-groups -----

type pGroupGen struct{}

// The concrete branch names groups whose centers are independently tabulated in
// the package test, so the generator cannot grade itself into being right. The
// least branch is the general fact: for a group of order p^k the center is a
// positive multiple of p, and the least such multiple p is attained.
func (g *pGroupGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type concrete struct {
		group, question, reason string
		prime, center           int
	}
	concretes := []concrete{
		{"Q_8", "The quaternion group Q_8 has order 8. What is the order of its center? (enter a number)", "The center of Q_8 is {1, -1}, of order 2.", 2, 2},
		{"D_4", "The dihedral group D_4 has order 8. What is the order of its center? (enter a number)", "Only the identity and the half-turn commute with everything, so the center has order 2.", 2, 2},
		{"D_8", "The dihedral group D_8 of order 16 has what order of center? (enter a number)", "A dihedral group of order 2n has center of order 2 when n is even, which includes D_8.", 2, 2},
		{"Z_9", "The cyclic group Z_9 is abelian. What is the order of its center? (enter a number)", "Every element of an abelian group is central, so the center is the whole group.", 3, 9},
		{"V_4", "The Klein four-group V_4 has what order of center? (enter a number)", "V_4 is abelian, so its center is all four elements.", 2, 4},
		{"Heisenberg(3)", "The Heisenberg group over F_3 has order 27. What is the order of its center? (enter a number)", "The center is the scalar matrices, three of them.", 3, 3},
	}
	if rand.Intn(3) == 0 {
		c := concretes[rand.Intn(len(concretes))]
		return generator.Problem{
			Question:    c.question,
			Answer:      strconv.Itoa(c.center),
			Explanation: c.reason,
			Facts: map[string]string{
				"kind": "concrete", "group": c.group,
				"prime": strconv.Itoa(c.prime), "answer": strconv.Itoa(c.center),
			},
		}
	}

	primes := []int{2, 3, 5, 7, 11, 13}
	p := primes[rand.Intn(len(primes))]
	k := 3 + rand.Intn(2)
	order := 1
	for i := 0; i < k; i++ {
		order *= p
	}
	return generator.Problem{
		Question:    fmt.Sprintf("A finite group G has order %d = %d^%d. Its center is nontrivial and |Z(G)| is divisible by %d. What is the least possible value of |Z(G)|? (enter a number)", order, p, k, p),
		Answer:      strconv.Itoa(p),
		Explanation: fmt.Sprintf("The class equation forces |Z(G)| to be a positive multiple of %d, and a group of order %d with center of order %d exists, so the least value is %d.", p, order, p, p),
		Facts: map[string]string{
			"kind": "least", "prime": strconv.Itoa(p),
			"order": strconv.Itoa(order), "answer": strconv.Itoa(p),
		},
	}
}
