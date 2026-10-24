package calculus

import (
	"fmt"
	"math/rand"

	"github.com/chuma-beep/mathua/internal/generator"
)

// conceptClaim is a multiple-choice claim whose answer is short enough to be
// read as a word, a verdict or a small number.
type conceptClaim struct {
	question    string
	correct     string
	reason      string
	distractors []string
}

// conceptDistinct keeps three distractors that differ from the answer and from
// each other.
//
// generator.MultipleChoice dedupes options and stops at four letters, so one
// collision silently turns a recognition question into a three-way guess.
func conceptDistinct(correct string, candidates ...string) []string {
	out := make([]string, 0, 3)
	seen := map[string]bool{correct: true}
	for _, c := range candidates {
		if c == "" || seen[c] {
			continue
		}
		seen[c] = true
		out = append(out, c)
		if len(out) == 3 {
			return out
		}
	}
	for _, spare := range []string{"yes", "no", "none of these", "0", "1", "2"} {
		if len(out) == 3 {
			break
		}
		if !seen[spare] {
			seen[spare] = true
			out = append(out, spare)
		}
	}
	return out
}

// conceptPickMC renders one claim of the pool as a four-option question.
func conceptPickMC(claims []conceptClaim) generator.Problem {
	c := claims[rand.Intn(len(claims))]
	p := generator.MultipleChoice(c.question, c.correct, conceptDistinct(c.correct, c.distractors...)...)
	p.Explanation = c.reason
	return p
}

// opposite returns the other verdict of a yes/no claim, and an empty string for
// anything else so a numeric answer never receives a yes/no distractor.
func opposite(answer string) string {
	switch answer {
	case "yes":
		return "no"
	case "no":
		return "yes"
	}
	return ""
}

// computedConceptClaims builds claims from drawn parameters. A concept whose
// table holds four fixed sentences gives a learner the whole supply in four
// questions, which is why each of these returns three computed claims that the
// table could not state.
func computedConceptClaims(id string) []conceptClaim {
	switch id {
	case "calc.integral.riemann_criteria":
		p := 2 + rand.Intn(5)
		a := 2 + rand.Intn(6)
		b := a + 1 + rand.Intn(6)
		jumps := 2 + rand.Intn(7)
		d := 2 + rand.Intn(3)
		return []conceptClaim{
			{
				question:    fmt.Sprintf("Is f(x) = x^{%d} Riemann integrable on [0, 1]?", p),
				correct:     "yes",
				reason:      fmt.Sprintf("x^{%d} is continuous on the closed interval [0, 1], so it is Riemann integrable.", p),
				distractors: []string{"no"},
			},
			{
				question:    fmt.Sprintf("Is f(x) = 1/(x - %d)^{%d} Riemann integrable on [%d, %d]?", a, d, a, b),
				correct:     "no",
				reason:      fmt.Sprintf("It is unbounded near x = %d, so no Riemann sum can converge.", a),
				distractors: []string{"yes"},
			},
			{
				question:    fmt.Sprintf("Is a bounded function with exactly %d discontinuities on [0, 1] Riemann integrable?", jumps),
				correct:     "yes",
				reason:      fmt.Sprintf("Finitely many discontinuities, here %d, cannot prevent integrability for a bounded function.", jumps),
				distractors: []string{"no"},
			},
		}
	case "calc.limit.little_o":
		p := 1 + rand.Intn(5)
		d := 1 + rand.Intn(3)
		c := 2 + rand.Intn(7)
		return []conceptClaim{
			{
				question:    fmt.Sprintf("Is x^{%d} = o(x^{%d}) as x → ∞?", p, p+d),
				correct:     "yes",
				reason:      fmt.Sprintf("x^{%d}/x^{%d} = 1/x^{%d} → 0, which is exactly what little-o requires.", p, p+d, d),
				distractors: []string{"no"},
			},
			{
				question:    fmt.Sprintf("Is x^{%d} = o(x^{%d}) as x → ∞?", p+d, p),
				correct:     "no",
				reason:      fmt.Sprintf("The ratio x^{%d}/x^{%d} = x^{%d} grows without bound, so it is not bounded, let alone tending to 0.", p+d, p, d),
				distractors: []string{"yes"},
			},
			{
				question:    fmt.Sprintf("If f = o(g) as x → ∞, must f = O(g) as x → ∞ for g(x) = x^{%d}?", c),
				correct:     "yes",
				reason:      "A ratio tending to 0 is bounded, and bounded ratio is the definition of big-O.",
				distractors: []string{"no"},
			},
		}
	case "calc.limit.weierstrass":
		p := 2 + rand.Intn(4)
		a := 1 + rand.Intn(6)
		b := a + 1 + rand.Intn(6)
		return []conceptClaim{
			{
				question:    fmt.Sprintf("Does the extreme value theorem apply to f(x) = x^{%d} on [%d, %d]?", p, a, b),
				correct:     "yes",
				reason:      fmt.Sprintf("x^{%d} is continuous and [%d, %d] is closed and bounded, so the maximum and minimum are attained.", p, a, b),
				distractors: []string{"no"},
			},
			{
				question:    fmt.Sprintf("Does the extreme value theorem apply to f(x) = 1/(x - %d) on (%d, %d]?", a, a, b),
				correct:     "no",
				reason:      fmt.Sprintf("The interval (%d, %d] is not closed, and the function is unbounded near %d.", a, b, a),
				distractors: []string{"yes"},
			},
			{
				question:    fmt.Sprintf("Does the extreme value theorem apply to f(x) = |x - %d| on [%d, %d]?", a, a, b),
				correct:     "yes",
				reason:      fmt.Sprintf("|x - %d| is continuous on the closed bounded interval [%d, %d], so it attains both extremes.", a, a, b),
				distractors: []string{"no"},
			},
		}
	case "calc.seq.cauchy":
		n := 2 + rand.Intn(8)
		p := 1 + rand.Intn(3)
		return []conceptClaim{
			{
				question:    fmt.Sprintf("Is the sequence a_n = 1/n^{%d} Cauchy?", p),
				correct:     "yes",
				reason:      fmt.Sprintf("1/n^{%d} converges to 0, and every convergent sequence is Cauchy.", p),
				distractors: []string{"no"},
			},
			{
				question:    fmt.Sprintf("Is the sequence a_n = n^{%d} Cauchy?", p),
				correct:     "no",
				reason:      fmt.Sprintf("Terms of n^{%d} separate without bound, so no tail stays within a fixed epsilon.", p),
				distractors: []string{"yes"},
			},
			{
				question:    fmt.Sprintf("Is the sequence a_n = (%d n + 1)/n Cauchy?", n),
				correct:     "yes",
				reason:      fmt.Sprintf("(%d n + 1)/n tends to %d, so it converges and is therefore Cauchy.", n, n),
				distractors: []string{"no"},
			},
		}
	case "calc.seq.bounded":
		n := 2 + rand.Intn(8)
		p := 1 + rand.Intn(3)
		return []conceptClaim{
			{
				question:    fmt.Sprintf("Is the sequence a_n = (-1)^n/%d bounded?", n),
				correct:     "yes",
				reason:      fmt.Sprintf("Every term lies between -1/%d and 1/%d.", n, n),
				distractors: []string{"no"},
			},
			{
				question:    fmt.Sprintf("Is the sequence a_n = n^{%d} bounded?", p),
				correct:     "no",
				reason:      fmt.Sprintf("n^{%d} grows without bound, so no constant bounds it.", p),
				distractors: []string{"yes"},
			},
			{
				question:    fmt.Sprintf("Is the sequence a_n = %d sin(n)/n bounded?", n),
				correct:     "yes",
				reason:      fmt.Sprintf("|%d sin(n)/n| ≤ %d, so the sequence is bounded and in fact tends to 0.", n, n),
				distractors: []string{"no"},
			},
		}
	case "calc.series.function_series":
		p := 2 + rand.Intn(4)
		r := 1 + rand.Intn(4)
		m := 2 + rand.Intn(3)
		return []conceptClaim{
			{
				question:    fmt.Sprintf("Does Σ x^n/n^{%d} converge uniformly on [-%d, %d]?", p, r, r),
				correct:     "yes",
				reason:      fmt.Sprintf("On [-%d, %d] the terms are bounded by %d^n/n^{%d}, whose sum converges, so the M-test applies.", r, r, r, p),
				distractors: []string{"no"},
			},
			{
				question:    "Does Σ x^n converge uniformly on [0, 1]?",
				correct:     "no",
				reason:      "The partial sums converge pointwise but the remainder near x = 1 stays large, so the convergence is not uniform.",
				distractors: []string{"yes"},
			},
			{
				question:    fmt.Sprintf("Does Σ 1/n^{%d} converge, as the M-test needs?", m),
				correct:     "yes",
				reason:      fmt.Sprintf("It is a p-series with exponent %d > 1, so it converges.", m),
				distractors: []string{"no"},
			},
		}
	case "calc.integral.p_test":
		p := 1 + rand.Intn(6)
		q := 1 + rand.Intn(3)
		lower := 2 + rand.Intn(7)
		verdict, because := "yes", "converges"
		if p == 1 {
			verdict, because = "no", "diverges"
		}
		return []conceptClaim{
			{
				question:    fmt.Sprintf("Does the integral of x^{-%d} over [%d, ∞) converge?", p, lower),
				correct:     verdict,
				reason:      fmt.Sprintf("The p-integral converges exactly when the exponent exceeds 1; here it is %d, so it %s.", p, because),
				distractors: []string{opposite(verdict)},
			},
			{
				question:    fmt.Sprintf("Does the integral of x^{-%d} over (0, 1] converge?", q),
				correct:     "yes",
				reason:      fmt.Sprintf("At the endpoint 0 the integral converges when the exponent is below 1; here it is %d.", q),
				distractors: []string{"no"},
			},
			{
				question:    fmt.Sprintf("Does the integral of x^{-%d} over (0, 1] converge?", q+2),
				correct:     "no",
				reason:      fmt.Sprintf("An exponent of %d exceeds 1, and near 0 that makes the integral diverge.", q+2),
				distractors: []string{"yes"},
			},
		}
	}
	return nil
}

// computedCurveSketchClaims adds curve-sketching questions whose numbers vary
// with the draw. The fixed table answers the same four questions every time.
func computedCurveSketchClaims() []conceptClaim {
	n := 2 + rand.Intn(5)
	a := 1 + rand.Intn(6)
	claims := []conceptClaim{
		{
			question: fmt.Sprintf("Where is the inflection point of \\(f(x)=x^{3}+%dx^{2}\\)?", 3*a),
			correct:  fmt.Sprintf("-%d", a),
			distractors: conceptDistinct(fmt.Sprintf("-%d", a),
				"0", fmt.Sprintf("%d", a), fmt.Sprintf("%d", 3*a)),
			reason: fmt.Sprintf("\\(f''(x) = 6x + %d\\), which vanishes at \\(x=%d\\) and changes sign there.", 6*a, -a),
		},
		{
			question: fmt.Sprintf("Is \\(f(x)=x^{%d}\\) increasing on \\((0, \\infty)\\)?", n),
			correct:  "yes",
			distractors: conceptDistinct("yes",
				"no"),
			reason: fmt.Sprintf("\\(f'(x) = %dx^{%d} > 0\\) for \\(x>0\\), so the function increases there.", n, n-1),
		},
		{
			question: fmt.Sprintf("Does \\(f(x)=1/(x-%d)\\) have a vertical asymptote at \\(x=%d\\)?", a, a),
			correct:  "yes",
			distractors: conceptDistinct("yes",
				"no"),
			reason: fmt.Sprintf("As \\(x \\to %d\\) the denominator tends to zero and the values grow without bound, so \\(x=%d\\) is a vertical asymptote.", a, a),
		},
	}
	return claims
}
