package algebra

import (
	"fmt"
	"math/rand"
)

// algDistinct returns three distractors for correct that are distinct from it
// and from each other.
//
// generator.MultipleChoice caps at four options and dedupes, so a candidate
// equal to the answer silently drops the question to three options and fails
// TestRaiseMultipleChoiceRejectsEveryDistractor.
// signedConst renders a shift so a stem reads "+ 3" or "- 2" rather than "+ -2".
func signedConst(k int) string {
	if k < 0 {
		return fmt.Sprintf(" - %d", -k)
	}
	return fmt.Sprintf(" + %d", k)
}

func algDistinct(correct string, candidates ...string) []string {
	out := make([]string, 0, 3)
	seen := map[string]bool{correct: true}
	for _, c := range candidates {
		if len(out) == 3 {
			break
		}
		if seen[c] {
			continue
		}
		seen[c] = true
		out = append(out, c)
	}
	for _, c := range []string{"0", "1", "-1", "none of these"} {
		if len(out) == 3 {
			break
		}
		if seen[c] {
			continue
		}
		seen[c] = true
		out = append(out, c)
	}
	return out
}

// computedAlgClaims returns parameterised claims for the concepts whose pools
// held three fixed sentences.
func computedAlgClaims(id string) []algMCClaim {
	switch id {
	case "alg.func.inverse_domain":
		return inverseDomainClaims()
	case "alg.log.graph":
		return logGraphClaims()
	case "alg.seq.induction":
		return inductionClaims()
	case "alg.systems.matrices":
		return matrixClaims()
	}
	return nil
}

func inverseDomainClaims() []algMCClaim {
	a := 1 + rand.Intn(4)
	b := a + 1 + rand.Intn(4)
	c := 1 + rand.Intn(5)
	d := c + 1 + rand.Intn(5)
	base := 2 + rand.Intn(5)
	p := 1 + rand.Intn(6)
	q := 1 + rand.Intn(9)

	rangeAnswer := fmt.Sprintf("[%d, %d]", c, d)
	squareAnswer := fmt.Sprintf("[0, %d]", base*base)
	preimageAnswer := fmt.Sprintf("%d", p)
	return []algMCClaim{
		{
			q:           fmt.Sprintf("If f has domain [%d, %d) and range [%d, %d], what is the domain of its inverse?", a, b, c, d),
			correct:     rangeAnswer,
			distractors: algDistinct(rangeAnswer, fmt.Sprintf("[%d, %d)", a, b), "[0, infinity)", fmt.Sprintf("[%d, %d]", c, c)),
			reason:      "The inverse takes the outputs of f as its inputs, so its domain is the range of f.",
		},
		{
			q:           fmt.Sprintf("If f(x) = x^2 is restricted to the interval [0, %d], what is the domain of its inverse?", base),
			correct:     squareAnswer,
			distractors: algDistinct(squareAnswer, fmt.Sprintf("[0, %d]", base), "(-infinity, infinity)", fmt.Sprintf("[0, %d]", 2*base)),
			reason:      fmt.Sprintf("On that interval the values of x^2 run from 0 to %d, so the inverse is defined exactly there.", base*base),
		},
		{
			q:           fmt.Sprintf("A one-to-one function satisfies f(%d) = %d. What is f^{-1}(%d)?", p, q, q),
			correct:     preimageAnswer,
			distractors: algDistinct(preimageAnswer, fmt.Sprintf("%d", q), "-1", fmt.Sprintf("%d", p+q)),
			reason:      fmt.Sprintf("The inverse sends the output %d back to the input that produced it, which is %d.", q, p),
		},
	}
}

func logGraphClaims() []algMCClaim {
	base := 2 + rand.Intn(4)
	h := 1 + rand.Intn(5)
	k := -3 + rand.Intn(7)
	value := 1 + rand.Intn(4)

	asymptote := fmt.Sprintf("%d", h)
	through := fmt.Sprintf("(1, %d)", k)
	atValue := fmt.Sprintf("%d", powInt(base, value))
	return []algMCClaim{
		{
			q:           fmt.Sprintf("The graph of y = log_%d(x - %d) has a vertical asymptote at which value of x?", base, h),
			correct:     asymptote,
			distractors: algDistinct(asymptote, "0", fmt.Sprintf("%d", base), fmt.Sprintf("%d", -h)),
			reason:      fmt.Sprintf("The logarithm is undefined at x - %d = 0, so the graph blows up there.", h),
		},
		{
			q:           fmt.Sprintf("The graph of y = log_%d(x)%s passes through which point?", base, signedConst(k)),
			correct:     through,
			distractors: algDistinct(through, fmt.Sprintf("(%d, 1)", k), fmt.Sprintf("(0, %d)", k), "(0, 1)"),
			reason:      fmt.Sprintf("At x = 1 the logarithm is 0, so the height is the constant %d.", k),
		},
		{
			q:           fmt.Sprintf("For y = log_%d(x), at which x does y equal %d?", base, value),
			correct:     atValue,
			distractors: algDistinct(atValue, fmt.Sprintf("%d", base*value), fmt.Sprintf("%d", value), fmt.Sprintf("%d", base+value)),
			reason:      fmt.Sprintf("log_%d(x) = %d means x = %d^%d = %d.", base, value, base, value, powInt(base, value)),
		},
	}
}

func inductionClaims() []algMCClaim {
	first := 1 + rand.Intn(5)
	step := 2 * (1 + rand.Intn(3))
	n := 3 + rand.Intn(8)
	m := 2 + rand.Intn(6)

	nthTerm := fmt.Sprintf("%d", first+(n-1)*step)
	nextTerm := fmt.Sprintf("%d", m+1)
	sumFirst := fmt.Sprintf("%d", n*(2*first+(n-1)*step)/2)
	return []algMCClaim{
		{
			q:           fmt.Sprintf("A sequence has a_1 = %d and increases by %d at each step. What is a_%d?", first, step, n),
			correct:     nthTerm,
			distractors: algDistinct(nthTerm, fmt.Sprintf("%d", first+n*step), fmt.Sprintf("%d", first+(n-1))),
			reason:      fmt.Sprintf("The induction step adds the common difference %d exactly %d times, giving %d + %d * %d.", step, n-1, first, n-1, step),
		},
		{
			q:           "In an induction that proves a formula for the sum of the first n terms, the step from m to m + 1 adds which single term?",
			correct:     nextTerm,
			distractors: algDistinct(nextTerm, fmt.Sprintf("%d", m), "m", fmt.Sprintf("%d", m+2)),
			reason:      fmt.Sprintf("Going from the sum of the first %d terms to the first %d adds the term numbered %d.", m, m+1, m+1),
		},
		{
			q:           fmt.Sprintf("An arithmetic sequence starts at %d and increases by %d. What is the sum of its first %d terms?", first, step, n),
			correct:     sumFirst,
			distractors: algDistinct(sumFirst, fmt.Sprintf("%d", n*first), fmt.Sprintf("%d", n*(2*first+(n-1)*step))),
			reason:      fmt.Sprintf("The sum is n(2a + (n-1)d)/2, which with n = %d, a = %d and d = %d is %s.", n, first, step, sumFirst),
		},
	}
}

func matrixClaims() []algMCClaim {
	rows := 2 + rand.Intn(3)
	cols := 2 + rand.Intn(3)
	a := 2 + rand.Intn(4)
	c := 2 + rand.Intn(4)
	rhs1 := 1 + rand.Intn(9)
	rhs2 := 1 + rand.Intn(9)

	columns := fmt.Sprintf("%d", cols+1)
	determinant := fmt.Sprintf("%d", 1-a*c)
	return []algMCClaim{
		{
			q:           fmt.Sprintf("A system of %d equations in %d unknowns has an augmented matrix with how many columns?", rows, cols),
			correct:     columns,
			distractors: algDistinct(columns, fmt.Sprintf("%d", cols), fmt.Sprintf("%d", rows), fmt.Sprintf("%d", rows*cols)),
			reason:      "The augmented matrix holds one column per unknown plus a final column for the constants.",
		},
		{
			q:           fmt.Sprintf("The coefficient matrix of a system with %d equations and %d unknowns has how many rows?", rows, cols),
			correct:     fmt.Sprintf("%d", rows),
			distractors: algDistinct(fmt.Sprintf("%d", rows), fmt.Sprintf("%d", cols), fmt.Sprintf("%d", cols+1), fmt.Sprintf("%d", rows*cols)),
			reason:      "Each equation contributes one row, whatever the number of unknowns.",
		},
		{
			q:           fmt.Sprintf("What is the determinant of the coefficient matrix of x + %d y = %d and %d x + y = %d?", a, rhs1, c, rhs2),
			correct:     determinant,
			distractors: algDistinct(determinant, fmt.Sprintf("%d", a*c), fmt.Sprintf("%d", 1+a*c), "-1"),
			reason:      fmt.Sprintf("The determinant of [[1, %d], [%d, 1]] is 1*1 - %d*%d = %s.", a, c, a, c, determinant),
		},
	}
}

// powInt is base raised to exp for small non-negative exponents.
func powInt(base, exp int) int {
	out := 1
	for i := 0; i < exp; i++ {
		out *= base
	}
	return out
}

// computedDirichletClaims asks the same questions about the Dirichlet function
// with drawn points, because four fixed sentences give a learner the whole
// supply in four questions.
func computedDirichletClaims() []algMCClaim {
	r := 1 + rand.Intn(9)
	s := r + 1 + rand.Intn(5)
	return []algMCClaim{
		{
			q:           fmt.Sprintf("What is D(%d/%d) for the Dirichlet function, 1 on rationals and 0 on irrationals?", r, s),
			correct:     "1",
			reason:      fmt.Sprintf("%d/%d is rational, so D takes the value 1 there.", r, s),
			distractors: []string{"0", "1/2", "2"},
		},
		{
			q:           fmt.Sprintf("Where is the Dirichlet function continuous when tested at x = sqrt(%d)?", r+1),
			correct:     "nowhere",
			reason:      "The Dirichlet function is discontinuous at every real number, rational or irrational.",
			distractors: []string{"everywhere", "at the rationals only", "at the irrationals only"},
		},
		{
			q:           fmt.Sprintf("How many intervals of the form [%d, %d] admit a Riemann integral of the Dirichlet function?", r, s),
			correct:     "none",
			reason:      "On any interval the upper sum is the width and the lower sum is 0, so the two never meet.",
			distractors: []string{"all of them", "only those of integer length", "only those starting at zero"},
		},
	}
}

// computedTransformationClaims adds transformation questions whose shifts
// vary with the draw, instead of the same four fixed shifts.
func computedTransformationClaims() []algMCClaim {
	k := 2 + rand.Intn(7)
	m := 2 + rand.Intn(4)
	return []algMCClaim{
		{
			q:           fmt.Sprintf("Compared with the graph of y = f(x), the graph of y = f(x) + %d is...", k),
			correct:     fmt.Sprintf("shifted up %d units", k),
			distractors: algDistinct(fmt.Sprintf("shifted up %d units", k), fmt.Sprintf("shifted down %d units", k), fmt.Sprintf("shifted right %d units", k), fmt.Sprintf("stretched vertically by a factor of %d", k)),
			reason:      fmt.Sprintf("Adding %d outside the function raises every output by %d, so the graph moves up by that amount.", k, k),
		},
		{
			q:           fmt.Sprintf("The graph of y = f(x - %d) is the graph of y = f(x)...", k),
			correct:     fmt.Sprintf("shifted right %d units", k),
			distractors: algDistinct(fmt.Sprintf("shifted right %d units", k), fmt.Sprintf("shifted left %d units", k), fmt.Sprintf("shifted up %d units", k), fmt.Sprintf("stretched horizontally by a factor of %d", k)),
			reason:      fmt.Sprintf("The input %d now produces the output the argument 0 produced before, so the graph moves right by %d.", k, k),
		},
		{
			q:           fmt.Sprintf("The graph of y = f(%dx) is the graph of y = f(x)...", m),
			correct:     fmt.Sprintf("compressed horizontally by a factor of %d", m),
			distractors: algDistinct(fmt.Sprintf("compressed horizontally by a factor of %d", m), fmt.Sprintf("stretched horizontally by a factor of %d", m), fmt.Sprintf("compressed vertically by a factor of %d", m), fmt.Sprintf("shifted right %d units", m)),
			reason:      fmt.Sprintf("Multiplying the input by %d makes each output occur %d times sooner, which squeezes the graph horizontally.", m, m),
		},
		{
			q:           fmt.Sprintf("The graph of y = f(x) - %d is the graph of y = f(x)...", k),
			correct:     fmt.Sprintf("shifted down %d units", k),
			distractors: algDistinct(fmt.Sprintf("shifted down %d units", k), fmt.Sprintf("shifted up %d units", k), fmt.Sprintf("shifted left %d units", k), fmt.Sprintf("stretched vertically by a factor of %d", k)),
			reason:      fmt.Sprintf("Subtracting %d outside the function lowers every output by %d, so the graph moves down.", k, k),
		},
	}
}
