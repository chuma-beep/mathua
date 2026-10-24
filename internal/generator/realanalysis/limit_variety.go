package realanalysis

import (
	"fmt"
	"math/rand"
)

// computedLimitClaims returns parameterised claims about limits, connectedness
// and monotone functions. The fixed claim pools these generators shipped with
// held three sentences, so a learner met the whole supply in three questions.
func computedLimitClaims(id string) []mcClaim {
	a := 1 + rand.Intn(8)
	b := a + 1 + rand.Intn(6)
	c := b + 2 + rand.Intn(5)
	switch id {
	case "ra.limit.connected":
		return []mcClaim{
			{
				q:           fmt.Sprintf("Is the open interval (%d, %d) connected?", a, b),
				correct:     "yes",
				distractors: []string{"no", "only if it is closed", "only if it is bounded"},
				reason:      "Every interval of the line is connected, and openness plays no part in that.",
			},
			{
				q:           fmt.Sprintf("Is the union [%d, %d] ∪ [%d, %d] connected?", a, b, c, c+a),
				correct:     "no",
				distractors: []string{"yes", "yes, because it is closed", "yes, because it is bounded"},
				reason:      fmt.Sprintf("The gap between %d and %d splits the union into two nonempty pieces separated in the subspace.", b, c),
			},
			{
				q:           fmt.Sprintf("Is the line with the point %d removed connected?", a),
				correct:     "no",
				distractors: []string{"yes", "only if the point is irrational", "yes, because removing a point cannot split a line"},
				reason:      fmt.Sprintf("The removed point separates the two rays, each of which is open and closed in the complement."),
			},
		}
	case "ra.limit.discontinuity":
		k := 1 + rand.Intn(5)
		return []mcClaim{
			{
				q:           fmt.Sprintf("The function that equals %d for x < %d and equals %d for x ≥ %d has at x = %d...", -k, a, k, a, a),
				correct:     "a jump discontinuity",
				distractors: []string{"a removable discontinuity", "an essential discontinuity", "no discontinuity"},
				reason:      fmt.Sprintf("The one-sided limits are %d and %d, both finite and different, which is a jump.", -k, k),
			},
			{
				q:           fmt.Sprintf("The function (x − %d)/(x − %d), undefined at x = %d, has there...", a, a, a),
				correct:     "a removable discontinuity",
				distractors: []string{"a jump discontinuity", "an essential discontinuity", "no discontinuity"},
				reason:      fmt.Sprintf("The quotient equals 1 wherever it is defined, so the limit is 1 and redefining the single value removes the gap."),
			},
			{
				q:           fmt.Sprintf("The function 1/(x − %d)%d has at x = %d...", a, k, a),
				correct:     "an essential discontinuity",
				distractors: []string{"a removable discontinuity", "a jump discontinuity", "no discontinuity"},
				reason:      fmt.Sprintf("The function is unbounded near %d, so neither one-sided limit is finite.", a),
			},
		}
	case "ra.limit.monotone_fn":
		a = 2 + rand.Intn(7)
		p := 2 + rand.Intn(4)
		n := 2 + rand.Intn(6)
		return []mcClaim{
			{
				q:           fmt.Sprintf("The increasing step function on [0, %d] with %d jumps is...", b, n),
				correct:     "Riemann integrable",
				distractors: []string{"not Riemann integrable", "continuous", "differentiable everywhere"},
				reason:      fmt.Sprintf("Its %d jump discontinuities form a finite set, so the function is Riemann integrable.", n),
			},
			{
				q:           fmt.Sprintf("The increasing function x^%d is continuous on [0, %d]. Can it have uncountably many discontinuities there?", p, b),
				correct:     "no",
				distractors: []string{"yes", "only if the exponent is even", "only on a closed interval"},
				reason:      "A monotone function has at most countably many discontinuities, and this one has none at all.",
			},
			{
				q:           fmt.Sprintf("If f increases on [0, %d] and f(x) ≤ %d for every x there, what is true of the left limit at %d?", b, b*a, b),
				correct:     fmt.Sprintf("it exists and is at most %d", b*a),
				distractors: []string{"it need not exist", "it equals f(b)", fmt.Sprintf("it is at least %d", b*a)},
				reason:      "Monotonicity makes the left limit the supremum of the values, and the bound limits that supremum.",
			},
		}
	}
	return nil
}

// computedDerivClaims returns parameterised claims about Darboux's theorem,
// exchanging limits with integrals and derivatives, and uniformity.
func computedDerivClaims(id string) []mcClaim {
	a := 2 + rand.Intn(7)
	p := 2 + rand.Intn(3)
	r := 1 + rand.Intn(4)
	switch id {
	case "ra.deriv.darboux":
		return []mcClaim{
			{
				q:           fmt.Sprintf("Does the derivative of x^%d on [%d, %d] take every value between its endpoint values?", p, a, a+r),
				correct:     "yes",
				distractors: []string{"no", "only when the derivative is continuous", "only when the interval is open"},
				reason:      "Derivatives have the intermediate value property by Darboux's theorem, continuous or not.",
			},
			{
				q:           "Is every function with the intermediate value property a derivative?",
				correct:     "no",
				distractors: []string{"yes", "yes, when it is bounded", "only on a closed interval"},
				reason:      "The intermediate value property is necessary for a derivative and not sufficient.",
			},
			{
				q:           fmt.Sprintf("The function x^2 sin(1/x) extended by 0 at the origin is differentiable everywhere on [−%d, %d]. Is its derivative continuous there?", a, a),
				correct:     "no",
				distractors: []string{"yes", "yes, because the function is differentiable", "only at the origin"},
				reason:      "The derivative oscillates without limit near the origin, so it is discontinuous there while still having the intermediate value property.",
			},
		}
	case "ra.deriv.interchange":
		n := 2 + rand.Intn(6)
		return []mcClaim{
			{
				q:           fmt.Sprintf("If f_n → f uniformly on [0, %d], does the integral of f equal the limit of the integrals?", a),
				correct:     "yes",
				distractors: []string{"no", "only if every f_n is continuous", "only for monotone sequences"},
				reason:      "Uniform convergence bounds the error in the integral by the sup norm times the length of the interval.",
			},
			{
				q:           fmt.Sprintf("Does f_n(x) = x^%d converge uniformly on [0, 1]?", n),
				correct:     "no",
				distractors: []string{"yes", "yes, because the limit is continuous", "only when the exponent is even"},
				reason:      fmt.Sprintf("The limit jumps from 0 to 1 at x = 1, so the sup distance stays 1 however large n is."),
			},
			{
				q:           "If f_n → f pointwise on [0, 1] and every f_n is continuous, must the limit of the integrals equal the integral of the limit?",
				correct:     "no",
				distractors: []string{"yes", "yes, when the limit is bounded", "only when the f_n are increasing"},
				reason:      "Pointwise convergence gives no bound on the sup distance, so mass can escape to a point.",
			},
		}
	case "ra.deriv.uniform_conv":
		return []mcClaim{
			{
				q:           fmt.Sprintf("Does f_n(x) = x/n converge uniformly on [0, %d]?", a),
				correct:     "yes",
				distractors: []string{"no", "only pointwise", "only when the interval is open"},
				reason:      fmt.Sprintf("The sup distance is %d/n, which tends to zero.", a),
			},
			{
				q:           fmt.Sprintf("Does f_n(x) = x^n converge uniformly on [0, %d] for %d < 1?", r, r),
				correct:     "yes",
				distractors: []string{"no", "only pointwise", "only when the ratio equals 1"},
				reason:      "On a closed interval strictly inside the unit disc the sup norm is r^n, which tends to zero.",
			},
			{
				q:           fmt.Sprintf("If every f_n is continuous and f_n → f uniformly, is f continuous?"),
				correct:     "yes",
				distractors: []string{"no", "only if the f_n are bounded", "only if the domain is compact"},
				reason:      "Uniform convergence preserves continuity, which is the content of the uniform limit theorem.",
			},
		}
	}
	return nil
}
