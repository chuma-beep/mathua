package calculus

import (
	"fmt"
	"math/rand"
	"strconv"
)

// computedThinClaims adds drawn claims to concepts whose pools held a handful
// of fixed sentences. Each generator keeps its authored claims and gains three
// computed ones, so the supply grows without contradicting what the fixed
// sentences teach.
func computedThinClaimsRaw(id string) []calcMCClaim {
	switch id {
	case "calc.limit.limsup":
		c := 1 + rand.Intn(5)
		// k must differ from c, or the distractor k repeats the answer and the
		// question collapses to three options after MultipleChoice dedupes.
		k := 2 + rand.Intn(5)
		for k == c {
			k = 2 + rand.Intn(5)
		}
		return []calcMCClaim{
			{
				q:           fmt.Sprintf("For the sequence a_n = %d + (-1)^n, what is the limit superior?", c),
				correct:     strconv.Itoa(c + 1),
				distractors: []string{strconv.Itoa(c - 1), strconv.Itoa(c), strconv.Itoa(c + 2)},
				reason:      fmt.Sprintf("The even terms are %d and the odd terms are %d, so every tail has supremum %d.", c+1, c-1, c+1),
			},
			{
				q:           fmt.Sprintf("For the sequence a_n = %d/(n+1) + %d, what is the limit superior?", k, c),
				correct:     strconv.Itoa(c),
				distractors: []string{strconv.Itoa(c + 1), strconv.Itoa(c - 1), strconv.Itoa(k)},
				reason:      fmt.Sprintf("The fraction shrinks to zero, so the tail suprema decrease to %d.", c),
			},
			{
				q:           "For the sequence a_n = (-1)^n(1 + 1/n), what is the limit superior?",
				correct:     "1",
				distractors: []string{"-1", "0", "2"},
				reason:      "The even terms approach 1 from above and the odd terms approach -1 from below, so the eventual supremum is 1.",
			},
		}
	case "calc.limit.oscillation":
		delta := 1 + rand.Intn(3)
		// A jump of one would make the distractor "1" repeat the answer.
		jump := 2 + rand.Intn(4)
		// A slope of one would make the answer equal to the width distractor.
		slope := 2 + rand.Intn(3)
		return []calcMCClaim{
			{
				q:           fmt.Sprintf("A step function jumps by %d at a point. What is its oscillation there?", jump),
				correct:     strconv.Itoa(jump),
				distractors: []string{"0", strconv.Itoa(2 * jump), "1"},
				reason:      fmt.Sprintf("The one-sided limits differ by %d, and oscillation is that difference for a jump.", jump),
			},
			{
				q:           fmt.Sprintf("On an interval of width %d, what is the oscillation of f(x) = x?", 2*delta),
				correct:     strconv.Itoa(2 * delta),
				distractors: []string{strconv.Itoa(delta), "0", strconv.Itoa(2*delta + 1)},
				reason:      "A linear function of slope one attains the supremum and infimum at the endpoints, so oscillation equals the width.",
			},
			{
				q:           fmt.Sprintf("On an interval of width %d, what is the oscillation of f(x) = %dx?", delta, slope),
				correct:     strconv.Itoa(slope * delta),
				distractors: []string{strconv.Itoa(delta), strconv.Itoa(slope), strconv.Itoa(slope*delta + 1)},
				reason:      fmt.Sprintf("The values differ by the slope times the width, which is %d times %d.", slope, delta),
			},
		}
	case "calc.series.operations":
		s := 2 + rand.Intn(8)
		t := 2 + rand.Intn(8)
		k := 2 + rand.Intn(3)
		base := 2 + rand.Intn(4)
		return []calcMCClaim{
			{
				q:           fmt.Sprintf("If a series sums to %d and another sums to %d, what does their term-by-term sum converge to?", s, t),
				correct:     strconv.Itoa(s + t),
				distractors: []string{strconv.Itoa(s + t + 1), strconv.Itoa(s + t - 1), strconv.Itoa(s*t + 2)},
				reason:      fmt.Sprintf("Limits add, so the partial sums converge to %d plus %d.", s, t),
			},
			{
				q:           fmt.Sprintf("If a series sums to %d, what does multiplying every term by %d give?", s, k),
				correct:     strconv.Itoa(k * s),
				distractors: []string{strconv.Itoa(k*s + 1), strconv.Itoa(k + s), strconv.Itoa(s)},
				reason:      fmt.Sprintf("Every partial sum is multiplied by %d, and the limit inherits the factor.", k),
			},
			{
				q:           fmt.Sprintf("What is the sum of the geometric series with ratio 1/%d?", base),
				correct:     fmt.Sprintf("1/%d", base-1),
				distractors: []string{fmt.Sprintf("1/%d", base), fmt.Sprintf("%d/%d", base, base+1), fmt.Sprintf("1/%d", base+1)},
				reason:      fmt.Sprintf("A geometric series with first term 1 and ratio r sums to 1/(1-r), which here is 1/%d.", base-1),
			},
		}
	case "calc.series.p_series":
		p := 1 + rand.Intn(5)
		k := 2 + rand.Intn(4)
		converges := "yes"
		if p == 1 {
			converges = "no"
		}
		return []calcMCClaim{
			{
				q:           fmt.Sprintf("Does the series with terms 1/n^%d converge?", p),
				correct:     converges,
				distractors: []string{otherVerdict(converges), "only for finitely many terms", "it converges to zero"},
				reason:      fmt.Sprintf("A p-series converges exactly when p exceeds 1, and here p is %d.", p),
			},
			{
				q:           fmt.Sprintf("What is the sum of the series with terms 1/%d^n, starting at n = 1?", k),
				correct:     fmt.Sprintf("1/%d", k-1),
				distractors: []string{fmt.Sprintf("1/%d", k), fmt.Sprintf("1/%d", k+1), strconv.Itoa(k)},
				reason:      fmt.Sprintf("The first term is 1/%d and the ratio is 1/%d, so the sum is 1/(%d-1).", k, k, k),
			},
			{
				q:           fmt.Sprintf("Does the series with terms 1/n^%d converge?", 2*p),
				correct:     "yes",
				distractors: []string{"no", "only when n is even", "only for the first few terms"},
				reason:      fmt.Sprintf("Doubling the exponent gives p equal to %d, which exceeds 1.", 2*p),
			},
		}
	case "calc.series.dirichlet_test":
		p := 1 + rand.Intn(5)
		q := 1 + rand.Intn(3)
		return []calcMCClaim{
			{
				q:           fmt.Sprintf("Does the series with terms (-1)^n/n^%d converge?", p),
				correct:     "yes",
				distractors: []string{"no", "only if the terms decrease faster", "only for even n"},
				reason:      fmt.Sprintf("The alternating signs have bounded partial sums and 1/n^%d decreases to zero, which is Dirichlet's test.", p),
			},
			{
				q:           fmt.Sprintf("Does the series with terms (-1)^n n^%d converge?", q),
				correct:     "no",
				distractors: []string{"yes", "yes, absolutely", "only when n is odd"},
				reason:      fmt.Sprintf("The terms n^%d grow rather than shrink to zero, so the series cannot converge.", q),
			},
			{
				q:           "In Dirichlet's test, the second factor must decrease monotonically to zero. Which series satisfies that?",
				correct:     fmt.Sprintf("the one with terms 1/n^%d", p),
				distractors: []string{fmt.Sprintf("the one with terms n^%d", q), "the one with terms (-1)^n", "the one with terms 1 + 1/n"},
				reason:      "Dirichlet needs the non-oscillating factor to tend to zero, which a negative power of n does and a positive power does not.",
			},
		}
	case "calc.deriv.convex_jensen":
		p := 2 + rand.Intn(3)
		a := 2 + rand.Intn(3)
		second := p * (p - 1)
		for i := 0; i < p-2; i++ {
			second *= a
		}
		return []calcMCClaim{
			{
				q:           fmt.Sprintf("What is the second derivative of x^%d at x = %d?", p, a),
				correct:     strconv.Itoa(second),
				distractors: []string{strconv.Itoa(second + 1), strconv.Itoa(p * a), strconv.Itoa(second - 2)},
				reason:      fmt.Sprintf("The first derivative is %d x^%d and the second is %d x^%d, which at x = %d is %d.", p, p-1, second, p-2, a, second),
			},
			{
				q:           fmt.Sprintf("Is x^%d convex on the whole real line?", 2+p%2),
				correct:     "yes",
				distractors: []string{"no", "only for positive x", "only between -1 and 1"},
				reason:      "An even power has a non-negative second derivative everywhere, which is exactly the convexity condition.",
			},
			{
				q:           "Is x^3 convex on the whole real line?",
				correct:     "no",
				distractors: []string{"yes", "only for positive x", "only for negative x"},
				reason:      "The second derivative 6x is negative for negative x, so the graph is concave there and convexity fails on the whole line.",
			},
		}
	case "calc.deriv.epsilon_delta_deriv":
		c := 2 + rand.Intn(5)
		a := 1 + rand.Intn(5)
		return []calcMCClaim{
			{
				q:           fmt.Sprintf("For f(x) = %dx and epsilon = 1, which delta proves continuity at every point?", c),
				correct:     fmt.Sprintf("1/%d", c),
				distractors: []string{strconv.Itoa(c), "1", fmt.Sprintf("1/%d", c+1)},
				reason:      fmt.Sprintf("The slope is %d, so an input change of delta moves the value by %d delta, and %d delta = 1 gives delta = 1/%d.", c, c, c, c),
			},
			{
				q:           fmt.Sprintf("For f(x) = x^2 at x = %d and epsilon = 1, which delta works?", a),
				correct:     fmt.Sprintf("1/%d", 2*a+1),
				distractors: []string{fmt.Sprintf("1/%d", 2*a), "1", fmt.Sprintf("1/%d", a)},
				reason:      fmt.Sprintf("The difference near x = %d is at most (2*%d + 1) delta, so the smallest safe choice is 1/%d.", a, a, 2*a+1),
			},
			{
				q:           "Which function is uniformly continuous on the whole real line?",
				correct:     "a linear function with a fixed slope",
				distractors: []string{"x^2", "e^x", "1/x"},
				reason:      "A fixed slope gives a single delta that works everywhere, while the others need a smaller delta as the point moves.",
			},
		}
	case "calc.deriv.inverse_hyperbolic":
		a := 2 + rand.Intn(4)
		return []calcMCClaim{
			{
				q:           fmt.Sprintf("What is the derivative of arsinh(x/%d)?", a),
				correct:     fmt.Sprintf("1/sqrt(x^2 + %d)", a*a),
				distractors: []string{fmt.Sprintf("1/sqrt(x^2 - %d)", a*a), fmt.Sprintf("%d/sqrt(x^2 + 1)", a), fmt.Sprintf("1/(x^2 + %d)", a)},
				reason:      fmt.Sprintf("By the chain rule the derivative is (1/%d) divided by sqrt((x/%d)^2 + 1), which simplifies to 1/sqrt(x^2 + %d).", a, a, a*a),
			},
			{
				q:           fmt.Sprintf("What is the derivative of artanh(x/%d)?", a),
				correct:     fmt.Sprintf("%d/(%d - x^2)", a, a*a),
				distractors: []string{fmt.Sprintf("%d/(%d + x^2)", a, a*a), fmt.Sprintf("1/(%d - x^2)", a*a), "1/(1 + x^2)"},
				reason:      fmt.Sprintf("The derivative is (1/%d) divided by (1 - (x/%d)^2), which is %d/(%d - x^2).", a, a, a, a*a),
			},
			{
				q:           fmt.Sprintf("What is the derivative of arcosh(x/%d) for x above %d?", a, a),
				correct:     fmt.Sprintf("1/sqrt(x^2 - %d)", a*a),
				distractors: []string{fmt.Sprintf("1/sqrt(x^2 + %d)", a*a), fmt.Sprintf("1/sqrt(%d - x^2)", a*a), fmt.Sprintf("%d/x", a)},
				reason:      fmt.Sprintf("The chain rule gives (1/%d) divided by sqrt((x/%d)^2 - 1), which is 1/sqrt(x^2 - %d).", a, a, a*a),
			},
		}
	}
	return nil
}

// otherVerdict flips a yes/no answer so a pool contains both.
func otherVerdict(answer string) string {
	if answer == "yes" {
		return "no"
	}
	return "yes"
}

// computedThinClaims adds drawn claims and guarantees that every one ships four
// distinct options. generator.MultipleChoice dedupes option text, so a
// parameter that coincides with a distractor silently shortens the question to
// three choices, which the batch test rejects.
func computedThinClaims(id string) []calcMCClaim {
	claims := computedThinClaimsRaw(id)
	for i := range claims {
		claims[i].distractors = thinDistractors(claims[i].correct, claims[i].distractors...)
	}
	return claims
}

// thinDistractors drops anything that repeats the answer or another option and
// tops the list up to three from a neutral pool.
func thinDistractors(correct string, candidates ...string) []string {
	seen := map[string]bool{correct: true}
	out := make([]string, 0, 3)
	for _, candidate := range candidates {
		if candidate == "" || seen[candidate] {
			continue
		}
		seen[candidate] = true
		out = append(out, candidate)
		if len(out) == 3 {
			return out
		}
	}
	for _, spare := range []string{"0", "-1", "1", "2", "3", "none of these"} {
		if seen[spare] {
			continue
		}
		seen[spare] = true
		out = append(out, spare)
		if len(out) == 3 {
			break
		}
	}
	return out
}
