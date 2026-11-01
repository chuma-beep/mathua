package algebra

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

// algVerdict is one statement in a count question. The generators in this
// package that used to ask "which of these is true" now ask how many of a
// fixed set are true, so the answer is a number the learner can enter rather
// than an option to recognise.
type algVerdict struct {
	label string
	yes   bool
	note  string
}

const algDraw = 8

// countAlgVerdicts draws a count from 0 to 8, takes that many true statements
// and the rest false, shuffles them, and asks for the count. Drawing the count
// first keeps every value 0..8 reachable; sampling eight from a pool with eight
// of each would make 0 and 8 vanishingly rare.
func countAlgVerdicts(lead string, pool []algVerdict) generator.Problem {
	yes := make([]algVerdict, 0, algDraw)
	no := make([]algVerdict, 0, algDraw)
	for _, v := range pool {
		if v.yes {
			yes = append(yes, v)
		} else {
			no = append(no, v)
		}
	}
	var chosen []algVerdict
	var count int
	question := ""
	for attempt := 0; attempt < 40; attempt++ {
		count = rand.Intn(algDraw + 1)
		if count > len(yes) {
			count = len(yes)
		}
		if algDraw-count > len(no) {
			continue
		}
		rand.Shuffle(len(yes), func(i, j int) { yes[i], yes[j] = yes[j], yes[i] })
		rand.Shuffle(len(no), func(i, j int) { no[i], no[j] = no[j], no[i] })
		chosen = append(append([]algVerdict{}, yes[:count]...), no[:algDraw-count]...)
		rand.Shuffle(len(chosen), func(i, j int) { chosen[i], chosen[j] = chosen[j], chosen[i] })
		var parts []string
		for i, v := range chosen {
			parts = append(parts, fmt.Sprintf("(%c) %s", 'a'+i, v.label))
		}
		question = fmt.Sprintf("%s %s (enter a number)", lead, strings.Join(parts, " "))
		if !strings.Contains(question, fmt.Sprintf("%d", count)) {
			break
		}
	}
	var notes []string
	for _, v := range chosen {
		verdict := "false"
		if v.yes {
			verdict = "true"
		}
		notes = append(notes, fmt.Sprintf("%s: %s, %s", v.label, verdict, v.note))
	}
	return generator.Problem{
		Question:    question,
		Answer:      fmt.Sprintf("%d", count),
		Explanation: fmt.Sprintf("Exactly %d of the %d statements are true. %s", count, algDraw, strings.Join(notes, "; ")),
		Facts:       map[string]string{"count": fmt.Sprintf("%d", count)},
	}
}

func inductionProblem() generator.Problem {
	pool := []algVerdict{
		{"A base case plus an induction step proves the statement for every integer from the base onward.", true, "the base starts the chain and the step carries it forward"},
		{"The induction step assumes the statement for n and proves it for n + 1.", true, "the step is a conditional from one rung to the next"},
		{"The base case is the starting point of the chain.", true, "without it no instance is reached"},
		{"If the base case fails, the proof does not establish the statement.", true, "the chain never begins"},
		{"Strong induction may assume the statement for every value below n.", true, "the stronger hypothesis is equivalent to ordinary induction"},
		{"Well-ordering and induction are equivalent principles on the natural numbers.", true, "each can prove the other"},
		{"If P(1) holds and P(n) implies P(n + 1) for every n, then P(5) holds.", true, "apply the step four times"},
		{"An induction proof must establish the step for an arbitrary n, not one chosen value.", true, "a single instance does not carry the chain"},
		{"An induction proof needs no base case.", false, "without a base the chain has no starting point"},
		{"The induction step proves P(n) from P(n + 1).", false, "the direction is from n to n + 1"},
		{"Checking one very large n proves the statement for all n.", false, "one instance says nothing about the others"},
		{"The base case alone proves the statement for every n.", false, "the step is what carries it forward"},
		{"Induction applies directly to the real numbers in their usual order.", false, "it needs a discrete successor or well-ordering"},
		{"Strong induction needs infinitely many base cases.", false, "one base case plus the strong step suffices"},
		{"A valid base case proves the statement even if the induction step is invalid.", false, "the step is required to reach later values"},
		{"Induction can be used to disprove a statement about all natural numbers.", false, "it proves universally quantified statements, it does not refute them"},
	}
	return countAlgVerdicts("How many of these statements about mathematical induction are true?", pool)
}

func conicClassifyProblem() generator.Problem {
	targets := []string{"circle", "ellipse", "parabola", "hyperbola"}
	target := targets[rand.Intn(len(targets))]
	makeEquation := func(class string) string {
		switch class {
		case "circle":
			r := rand.Intn(7) + 2
			return fmt.Sprintf("x^2 + y^2 = %d", r*r)
		case "ellipse":
			p := rand.Intn(7) + 2
			q := rand.Intn(7) + 2
			for q == p {
				q = rand.Intn(7) + 2
			}
			return fmt.Sprintf("x^2/%d + y^2/%d = 1", p*p, q*q)
		case "parabola":
			a := rand.Intn(3) + 1
			if rand.Intn(2) == 0 {
				a = -a
			}
			b := rand.Intn(7) - 3
			c := rand.Intn(11) - 5
			return "y = " + formatPoly(polyTerm{a, 2}, polyTerm{b, 1}, polyTerm{c, 0})
		default:
			p := rand.Intn(7) + 2
			q := rand.Intn(7) + 2
			return fmt.Sprintf("x^2/%d - y^2/%d = 1", p*p, q*q)
		}
	}
	pool := make([]algVerdict, 0, 16)
	for i := 0; i < algDraw; i++ {
		pool = append(pool, algVerdict{makeEquation(target), true, fmt.Sprintf("this equation has the form of a %s", target)})
	}
	for i := 0; i < algDraw; i++ {
		other := targets[rand.Intn(len(targets))]
		for other == target {
			other = targets[rand.Intn(len(targets))]
		}
		pool = append(pool, algVerdict{makeEquation(other), false, fmt.Sprintf("this equation has the form of a %s", other)})
	}
	article := "a"
	if target == "ellipse" {
		article = "an"
	}
	return countAlgVerdicts(fmt.Sprintf("How many of these equations represent %s %s?", article, target), pool)
}
