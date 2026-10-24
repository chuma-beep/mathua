package trigonometry

import (
	"fmt"
	"math/rand"
	"strconv"
)

// computedTrigClaims returns claims whose answers are derived from drawn
// parameters rather than recited.
//
// Each of these concepts shipped a pool of three fixed sentences, so a learner
// met its entire supply in three questions. The identity pools in particular
// described the general identity and never applied it, which is the difference
// between recognising a formula and using one.
// computedTrigClaims returns the computed claims for one concept, with every
// distractor list made safe for generator.MultipleChoice.
//
// MultipleChoice dedupes and caps at four options, so a distractor equal to the
// correct answer or to another distractor silently drops a question to three
// options and fails the batch test. Parameter collisions do this at particular
// draws rather than always, so a probe over 200 draws found five such cases in
// the hyperbolic identities before this guard existed.
func computedTrigClaims(id string) []trigMCClaim {
	claims := computedTrigClaimsRaw(id)
	for i := range claims {
		claims[i].distractors = trigDistinct(claims[i].correct, claims[i].distractors...)
	}
	return claims
}

// trigDistinct drops repeats of the correct answer and of each other, then tops
// up from a pool of neutral values so exactly three distractors survive.
func trigDistinct(correct string, candidates ...string) []string {
	seen := map[string]bool{correct: true}
	out := make([]string, 0, len(candidates))
	for _, c := range candidates {
		if c == "" || seen[c] {
			continue
		}
		seen[c] = true
		out = append(out, c)
	}
	for _, top := range []string{"0", "1", "-1", "none of these"} {
		if len(out) >= 3 {
			break
		}
		if !seen[top] {
			seen[top] = true
			out = append(out, top)
		}
	}
	return out
}

// degSigned renders an angle so a negative value reads as cos(-10) rather than
// cos -10, which looks like a subtraction.
func degSigned(v int) string {
	if v < 0 {
		return fmt.Sprintf("(%d)", v)
	}
	return strconv.Itoa(v)
}

func computedTrigClaimsRaw(id string) []trigMCClaim {
	switch id {
	case "trig.ident.product_to_sum":
		a := 10 + 5*rand.Intn(10)
		b := 5 + 5*rand.Intn(8)
		if a < b {
			a, b = b, a
		}
		sum, diff := a+b, a-b
		return []trigMCClaim{
			{
				q:       fmt.Sprintf("Write sin %d cos %d as a sum.", a, b),
				correct: fmt.Sprintf("(sin %d + sin %d)/2", sum, diff),
				distractors: []string{
					fmt.Sprintf("(sin %d - sin %d)/2", sum, diff),
					fmt.Sprintf("(cos %d - cos %d)/2", sum, diff),
					fmt.Sprintf("(cos %d + cos %d)/2", diff, sum),
				},
				reason: fmt.Sprintf("Adding the sum and difference formulas for sine cancels the cosines: sin %d cos %d is half of sin %d plus sin %d.", a, b, sum, diff),
			},
			{
				q:       fmt.Sprintf("Write cos %d cos %d as a sum.", a, b),
				correct: fmt.Sprintf("(cos %d + cos %d)/2", diff, sum),
				distractors: []string{
					fmt.Sprintf("(cos %d - cos %d)/2", diff, sum),
					fmt.Sprintf("(sin %d + sin %d)/2", sum, diff),
					fmt.Sprintf("2 cos %d cos %d", sum, diff),
				},
				reason: fmt.Sprintf("Both cosine formulas contain the product with the same sign, so adding them gives cos %d cos %d as half of cos %d plus cos %d.", a, b, diff, sum),
			},
			{
				q:       fmt.Sprintf("Write sin %d sin %d as a difference of cosines.", a, b),
				correct: fmt.Sprintf("(cos %d - cos %d)/2", diff, sum),
				distractors: []string{
					fmt.Sprintf("(cos %d - cos %d)/2", sum, diff),
					fmt.Sprintf("(cos %d + cos %d)/2", diff, sum),
					fmt.Sprintf("(sin %d - sin %d)/2", sum, diff),
				},
				reason: fmt.Sprintf("Subtracting the cosine formulas leaves twice the product of the sines with opposite signs, so sin %d sin %d is half of cos %d minus cos %d.", a, b, diff, sum),
			},
		}

	case "trig.ident.sum_to_product":
		p := 15 + 5*rand.Intn(10)
		q := 5 + 5*rand.Intn(5)
		// Equal halves would make "2 cos p cos q" and "2 cos q cos p" the same
		// option, and MultipleChoice dedupes them down to three.
		if q == p {
			q += 5
		}
		a, b := p+q, p-q
		return []trigMCClaim{
			{
				q:       fmt.Sprintf("Write sin %s + sin %s as a product.", degSigned(a), degSigned(b)),
				correct: fmt.Sprintf("2 sin %d cos %d", p, q),
				distractors: []string{
					fmt.Sprintf("2 cos %d cos %d", p, q),
					fmt.Sprintf("2 sin %d sin %d", p, q),
					fmt.Sprintf("2 sin %d cos %d", q, p),
				},
				reason: fmt.Sprintf("A sum of sines is twice the sine of the mean times the cosine of the half-difference: the mean is %d and the half-difference is %d.", p, q),
			},
			{
				q:       fmt.Sprintf("Write cos %s + cos %s as a product.", degSigned(a), degSigned(b)),
				correct: fmt.Sprintf("2 cos %d cos %d", p, q),
				distractors: []string{
					fmt.Sprintf("2 sin %d sin %d", p, q),
					fmt.Sprintf("2 cos %d cos %d", q, p),
					fmt.Sprintf("2 cos %d sin %d", p, q),
				},
				reason: fmt.Sprintf("Cosines add to twice the cosine of the mean times the cosine of the half-difference, and here those are %d and %d.", p, q),
			},
			{
				q:       fmt.Sprintf("Write cos %s - cos %s as a product.", degSigned(a), degSigned(b)),
				correct: fmt.Sprintf("-2 sin %d sin %d", p, q),
				distractors: []string{
					fmt.Sprintf("2 sin %d sin %d", p, q),
					fmt.Sprintf("-2 cos %d cos %d", p, q),
					fmt.Sprintf("-2 sin %d sin %d", q, p),
				},
				reason: fmt.Sprintf("A difference of cosines reverses the sign of the product of sines, so cos %d minus cos %d is minus twice sin %d sin %d.", a, b, p, q),
			},
		}

	case "trig.ident.pythagorean_variants":
		k := 2 + rand.Intn(8)
		d := 3 + rand.Intn(6)
		n := 1 + rand.Intn(d-1)
		return []trigMCClaim{
			{
				q:           fmt.Sprintf("If tan t = %d, what is sec^2 t?", k),
				correct:     strconv.Itoa(k*k + 1),
				distractors: []string{strconv.Itoa(k*k - 1), strconv.Itoa(k*k + 2), strconv.Itoa(2 * k)},
				reason:      fmt.Sprintf("tan^2 t + 1 = sec^2 t, so sec^2 t is %d squared plus one.", k),
			},
			{
				q:           fmt.Sprintf("If cot t = %d, what is csc^2 t?", k),
				correct:     strconv.Itoa(k*k + 1),
				distractors: []string{strconv.Itoa(k*k - 1), strconv.Itoa(k*k + 2), strconv.Itoa(2 * k)},
				reason:      fmt.Sprintf("1 + cot^2 t = csc^2 t, so csc^2 t is %d squared plus one.", k),
			},
			{
				q:           fmt.Sprintf("If sin t = %d/%d, what is cos^2 t?", n, d),
				correct:     trigFrac(d*d-n*n, d*d),
				distractors: []string{trigFrac(n*n, d*d), trigFrac(d*d+n*n, d*d), trigFrac(d-n, d)},
				reason:      fmt.Sprintf("cos^2 t is 1 minus sin^2 t, and sin^2 t is %d/%d, so cos^2 t is %s.", n*n, d*d, trigFrac(d*d-n*n, d*d)),
			},
		}

	case "trig.hyperbolic.identities":
		c := 2 + rand.Intn(7)
		q := 2 + rand.Intn(4)
		p := 1 + rand.Intn(q-1)
		return []trigMCClaim{
			{
				q:           fmt.Sprintf("If cosh x = %d, what is sinh^2 x?", c),
				correct:     strconv.Itoa(c*c - 1),
				distractors: []string{strconv.Itoa(c*c + 1), strconv.Itoa(c * c), strconv.Itoa(c - 1)},
				reason:      fmt.Sprintf("cosh^2 x - sinh^2 x = 1, so sinh^2 x is %d squared minus one.", c),
			},
			{
				q:           fmt.Sprintf("If sinh x = %d, what is cosh^2 x?", c),
				correct:     strconv.Itoa(c*c + 1),
				distractors: []string{strconv.Itoa(c*c - 1), strconv.Itoa(c * c), strconv.Itoa(2 * c)},
				reason:      fmt.Sprintf("The fundamental identity rearranges to cosh^2 x equal to sinh^2 x plus one, here %d squared plus one.", c),
			},
			{
				q:           fmt.Sprintf("If tanh x = %d/%d, what is sech^2 x?", p, q),
				correct:     trigFrac(q*q-p*p, q*q),
				distractors: []string{trigFrac(p*p, q*q), trigFrac(q*q+p*p, q*q), trigFrac(q-p, q)},
				reason:      fmt.Sprintf("1 - tanh^2 x = sech^2 x, and tanh^2 x is %d/%d, so sech^2 x is %s.", p*p, q*q, trigFrac(q*q-p*p, q*q)),
			},
		}

	case "trig.hyperbolic.derivative":
		a := 2 + rand.Intn(7)
		return []trigMCClaim{
			{
				q:       fmt.Sprintf("What is the derivative of sinh(%dx)?", a),
				correct: fmt.Sprintf("%d cosh(%dx)", a, a),
				distractors: []string{
					fmt.Sprintf("%d sinh(%dx)", a, a),
					fmt.Sprintf("cosh(%dx)", a),
					fmt.Sprintf("-%d cosh(%dx)", a, a),
				},
				reason: fmt.Sprintf("The chain rule multiplies the derivative of the outer sinh by the slope %d, and the derivative of sinh is cosh with no sign change.", a),
			},
			{
				q:       fmt.Sprintf("What is the derivative of cosh(%dx)?", a),
				correct: fmt.Sprintf("%d sinh(%dx)", a, a),
				distractors: []string{
					fmt.Sprintf("-%d sinh(%dx)", a, a),
					fmt.Sprintf("%d cosh(%dx)", a, a),
					fmt.Sprintf("sinh(%dx)", a),
				},
				reason: fmt.Sprintf("Cosine hyperbolic differentiates to sine hyperbolic with a plus sign, and the chain rule contributes the factor %d.", a),
			},
			{
				q:       fmt.Sprintf("What is the derivative of tanh(%dx)?", a),
				correct: fmt.Sprintf("%d sech^2(%dx)", a, a),
				distractors: []string{
					fmt.Sprintf("-%d sech^2(%dx)", a, a),
					fmt.Sprintf("%d tanh^2(%dx)", a, a),
					fmt.Sprintf("sech^2(%dx)", a),
				},
				reason: fmt.Sprintf("The derivative of tanh is sech^2, and the inner slope %d multiplies it.", a),
			},
		}

	case "trig.eq.general_solution":
		return []trigMCClaim{
			{
				q:           "The general solution of sin t = 1/2, in degrees, is...",
				correct:     "t = 30 + 360k or t = 150 + 360k",
				distractors: []string{"t = 30 + 180k", "t = 150 + 360k only", "t = 30 + 360k only"},
				reason:      "Sine equals one half at 30 and at its mirror 150 within the first turn, and each repeats every full turn.",
			},
			{
				q:           "The general solution of cos t = 1/2, in degrees, is...",
				correct:     "t = 60 + 360k or t = 300 + 360k",
				distractors: []string{"t = 60 + 180k", "t = 120 + 360k or t = 240 + 360k", "t = 60 + 360k only"},
				reason:      "Cosine is symmetric about the horizontal axis, so the two solutions in a turn are 60 and 300, each repeating every turn.",
			},
			{
				q:           "The general solution of tan t = 1, in degrees, is...",
				correct:     "t = 45 + 180k",
				distractors: []string{"t = 45 + 360k", "t = 45 + 360k or t = 225 + 360k", "t = 135 + 180k"},
				reason:      "Tangent has period 180 degrees, so one solution generates the rest by adding whole half-turns.",
			},
			{
				q:           "The general solution of sin t = -1/2, in degrees, is...",
				correct:     "t = 210 + 360k or t = 330 + 360k",
				distractors: []string{"t = 30 + 360k or t = 150 + 360k", "t = 210 + 180k", "t = 330 + 360k only"},
				reason:      "Sine is negative in the third and fourth quadrants, at 210 and 330, and each repeats every full turn.",
			},
			{
				q:           "The general solution of cos t = -1/2, in degrees, is...",
				correct:     "t = 120 + 360k or t = 240 + 360k",
				distractors: []string{"t = 60 + 360k or t = 300 + 360k", "t = 120 + 180k", "t = 240 + 360k only"},
				reason:      "Cosine is negative in the second and third quadrants, at 120 and 240, each repeating every turn.",
			},
			{
				q:           "The general solution of sin t = 0, in degrees, is...",
				correct:     "t = 180k",
				distractors: []string{"t = 90 + 180k", "t = 90k", "t = 360k"},
				reason:      "Sine vanishes at every whole turn, 0, 180, 360 and so on, which is exactly the multiples of 180.",
			},
			{
				q:           "The general solution of cos t = 0, in degrees, is...",
				correct:     "t = 90 + 180k",
				distractors: []string{"t = 180k", "t = 90 + 360k", "t = 45 + 90k"},
				reason:      "Cosine vanishes at the quarter turns 90 and 270 and again every half-turn after that.",
			},
			{
				q:           "The general solution of tan t = -1, in degrees, is...",
				correct:     "t = 135 + 180k",
				distractors: []string{"t = 45 + 180k", "t = 135 + 360k", "t = 315 + 360k only"},
				reason:      "Tangent repeats every 180 degrees, and it is negative in the second quadrant, so 135 alone generates the family.",
			},
		}

	case "trig.ineq.advanced":
		return []trigMCClaim{
			{
				q:           "The solution of sin t > 1/2 in the interval from 0 to 360 degrees is...",
				correct:     "30 < t < 150",
				distractors: []string{"0 < t < 30 or 150 < t < 360", "30 < t < 330", "150 < t < 360"},
				reason:      "Sine rises past one half at 30, stays above it across the top of the circle, and falls back below at the mirror angle 150.",
			},
			{
				q:           "The solution of sin t > 0 in the interval from 0 to 360 degrees is...",
				correct:     "0 < t < 180",
				distractors: []string{"90 < t < 270", "0 < t < 90", "180 < t < 360"},
				reason:      "Sine is positive exactly in the first and second quadrants, which run from 0 to 180 degrees.",
			},
			{
				q:           "The solution of cos t < 0 in the interval from 0 to 360 degrees is...",
				correct:     "90 < t < 270",
				distractors: []string{"0 < t < 90 or 270 < t < 360", "270 < t < 360", "0 < t < 180"},
				reason:      "Cosine is negative in the second and third quadrants, running from 90 to 270 degrees.",
			},
			{
				q:           "The solution of cos t < 1/2 in the interval from 0 to 360 degrees is...",
				correct:     "60 < t < 300",
				distractors: []string{"0 < t < 60 or 300 < t < 360", "90 < t < 270", "60 < t < 120"},
				reason:      "Cosine falls below one half just after 60 and stays there until it returns to one half at 300.",
			},
			{
				q:           "The solution of tan t >= 1 in the interval from 0 to 180 degrees is...",
				correct:     "45 <= t < 90",
				distractors: []string{"0 <= t <= 45", "45 <= t <= 90", "90 < t < 180"},
				reason:      "Tangent reaches one at 45 and grows without bound as the angle approaches 90, which is excluded because tangent is undefined there.",
			},
			{
				q:           "The solution of tan t < 0 in the interval from 0 to 180 degrees is...",
				correct:     "90 < t < 180",
				distractors: []string{"0 < t < 90", "0 < t < 45", "135 < t < 180"},
				reason:      "Tangent is negative in the second quadrant, where sine is positive and cosine is negative.",
			},
			{
				q:           "The solution of sin t < -1/2 in the interval from 0 to 360 degrees is...",
				correct:     "210 < t < 330",
				distractors: []string{"150 < t < 210", "0 < t < 30 or 150 < t < 360", "210 < t < 300"},
				reason:      "Sine is below minus one half between the third quadrant angle 210 and the fourth quadrant angle 330.",
			},
			{
				q:           "The solution of cos t > 1/2 in the interval from 0 to 360 degrees is...",
				correct:     "0 <= t < 60 or 300 < t <= 360",
				distractors: []string{"60 < t < 300", "0 <= t < 90", "300 < t <= 360 only"},
				reason:      "Cosine exceeds one half near angle zero on both sides, up to 60 and again after 300, and the interval wraps through 0.",
			},
		}
	}
	return nil
}
