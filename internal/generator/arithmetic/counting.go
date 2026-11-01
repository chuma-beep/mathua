package arithmetic

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

// arithVerdict is one statement or number in a count question. The point of
// these questions is that the learner produces the count instead of picking a
// yes/no option, so the pool always holds both true and false cases.
type arithVerdict struct {
	label string
	yes   bool
	note  string
}

const arithDraw = 8

// countArithVerdicts draws the answer first, then builds a list with exactly
// that many true cases and the rest false. Drawing the answer first keeps all
// nine counts equally likely; sampling eight from a balanced pool would make
// zero and eight vanishingly rare.
func countArithVerdicts(lead string, pool []arithVerdict) generator.Problem {
	var yes, no []arithVerdict
	for _, v := range pool {
		if v.yes {
			yes = append(yes, v)
		} else {
			no = append(no, v)
		}
	}
	var count int
	var chosen []arithVerdict
	var question string
	for attempt := 0; attempt < 40; attempt++ {
		count = rand.Intn(arithDraw + 1)
		if count > len(yes) {
			count = len(yes)
		}
		if arithDraw-count > len(no) {
			count = arithDraw - len(no)
		}
		rand.Shuffle(len(yes), func(i, j int) { yes[i], yes[j] = yes[j], yes[i] })
		rand.Shuffle(len(no), func(i, j int) { no[i], no[j] = no[j], no[i] })
		chosen = append(append([]arithVerdict{}, yes[:count]...), no[:arithDraw-count]...)
		rand.Shuffle(len(chosen), func(i, j int) { chosen[i], chosen[j] = chosen[j], chosen[i] })
		parts := make([]string, len(chosen))
		for i, v := range chosen {
			parts[i] = fmt.Sprintf("(%c) %s", 'a'+i, v.label)
		}
		question = fmt.Sprintf("%s %s (enter a number)", lead, strings.Join(parts, " "))
		if !strings.Contains(question, fmt.Sprint(count)) {
			break
		}
	}
	notes := make([]string, len(chosen))
	for i, v := range chosen {
		verdict := "false"
		if v.yes {
			verdict = "true"
		}
		notes[i] = fmt.Sprintf("%s: %s. %s", v.label, verdict, v.note)
	}
	return generator.Problem{
		Question:    question,
		Answer:      fmt.Sprint(count),
		Explanation: fmt.Sprintf("Exactly %d of the %d statements are true. %s", count, arithDraw, strings.Join(notes, " ")),
		Facts:       map[string]string{"count": fmt.Sprint(count)},
	}
}

func primeProblem() generator.Problem {
	pool := []arithVerdict{
		{label: "2", yes: true, note: "2 has no divisors other than 1 and itself."},
		{label: "3", yes: true, note: "3 has no divisors other than 1 and itself."},
		{label: "5", yes: true, note: "5 has no divisors other than 1 and itself."},
		{label: "7", yes: true, note: "7 has no divisors other than 1 and itself."},
		{label: "11", yes: true, note: "11 has no divisors other than 1 and itself."},
		{label: "13", yes: true, note: "13 has no divisors other than 1 and itself."},
		{label: "17", yes: true, note: "17 has no divisors other than 1 and itself."},
		{label: "19", yes: true, note: "19 has no divisors other than 1 and itself."},
		{label: "4", yes: false, note: "4 = 2 x 2, so it is composite."},
		{label: "6", yes: false, note: "6 = 2 x 3, so it is composite."},
		{label: "8", yes: false, note: "8 = 2 x 4, so it is composite."},
		{label: "9", yes: false, note: "9 = 3 x 3, so it is composite."},
		{label: "10", yes: false, note: "10 = 2 x 5, so it is composite."},
		{label: "12", yes: false, note: "12 = 3 x 4, so it is composite."},
		{label: "14", yes: false, note: "14 = 2 x 7, so it is composite."},
		{label: "15", yes: false, note: "15 = 3 x 5, so it is composite."},
	}
	return countArithVerdicts("How many of these numbers are prime?", pool)
}

func compositeProblem() generator.Problem {
	pool := []arithVerdict{
		{label: "23", yes: true, note: "23 = 2 x 6 + 1 and has no divisor up to 4, so it is prime."},
		{label: "29", yes: true, note: "29 has no divisor up to 5, so it is prime."},
		{label: "31", yes: true, note: "31 has no divisor up to 5, so it is prime."},
		{label: "37", yes: true, note: "37 has no divisor up to 6, so it is prime."},
		{label: "41", yes: true, note: "41 has no divisor up to 6, so it is prime."},
		{label: "43", yes: true, note: "43 has no divisor up to 6, so it is prime."},
		{label: "47", yes: true, note: "47 has no divisor up to 6, so it is prime."},
		{label: "53", yes: true, note: "53 has no divisor up to 7, so it is prime."},
		{label: "16", yes: false, note: "16 = 2 x 8, so it is composite."},
		{label: "18", yes: false, note: "18 = 2 x 9, so it is composite."},
		{label: "20", yes: false, note: "20 = 4 x 5, so it is composite."},
		{label: "21", yes: false, note: "21 = 3 x 7, so it is composite."},
		{label: "22", yes: false, note: "22 = 2 x 11, so it is composite."},
		{label: "24", yes: false, note: "24 = 4 x 6, so it is composite."},
		{label: "25", yes: false, note: "25 = 5 x 5, so it is composite."},
		{label: "26", yes: false, note: "26 = 2 x 13, so it is composite."},
	}
	return countArithVerdicts("How many of these numbers are composite?", pool)
}

func divisionByZeroProblem() generator.Problem {
	pool := []arithVerdict{
		{label: "Division by zero is undefined.", yes: true, note: "No real number times zero gives a nonzero dividend."},
		{label: "0 divided by 5 is 0.", yes: true, note: "5 x 0 = 0, so the quotient is 0."},
		{label: "No number times 0 is 7.", yes: true, note: "Every product with 0 is 0, never 7."},
		{label: "7 divided by 0 is undefined.", yes: true, note: "The equation 0 x q = 7 has no solution."},
		{label: "0 divided by 0 is undefined.", yes: true, note: "Every q solves 0 x q = 0, so no single quotient is determined."},
		{label: "For a nonzero, 0/a = 0.", yes: true, note: "a x 0 = 0, so the quotient is 0."},
		{label: "0*q = 9 has no solution.", yes: true, note: "The left side is always 0, never 9."},
		{label: "0*q = 0 has every solution.", yes: true, note: "Every real q satisfies it."},
		{label: "5 divided by 0 is 0.", yes: false, note: "0 x 0 = 0, not 5, so the quotient cannot be 0."},
		{label: "0 divided by 0 is 0.", yes: false, note: "Every number solves 0 x q = 0, so the quotient is not determined."},
		{label: "0 divided by 7 is undefined.", yes: false, note: "7 x 0 = 0, so 0 divided by 7 equals 0."},
		{label: "Division by zero is infinity.", yes: false, note: "Infinity is not a real number, so the expression stays undefined."},
		{label: "Every x/0 is 0.", yes: false, note: "A nonzero dividend divided by zero has no value."},
		{label: "0 divided by 5 is undefined.", yes: false, note: "5 x 0 = 0, so 0 divided by 5 equals 0."},
		{label: "8 divided by 0 is 8.", yes: false, note: "0 x 8 = 0, not 8, so the quotient cannot be 8."},
		{label: "0*q = 0 has no solution.", yes: false, note: "Every real q solves it."},
	}
	return countArithVerdicts("How many division statements are true?", pool)
}
