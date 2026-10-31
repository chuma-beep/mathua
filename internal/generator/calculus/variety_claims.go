package calculus

// Variety claims for concepts whose generator table held only a handful of
// fixed sentences or one shape. The ratchet measures distinct question stems in
// 60 draws, so a table of five rows gives a learner the whole supply in five
// questions. Each case below draws its parameters, which makes the stem differ
// from draw to draw without changing what the concept asks.
//
// These claims are multiple choice on purpose. Some of the concepts ship as
// symbolic, and the calculus package already renders symbolic concepts as
// verdict questions, for example the Riemann integrability criteria.

import (
	"fmt"
	"math/rand"

	"github.com/chuma-beep/mathua/internal/generator"
)

// varietyWrap mixes computed claims with the concept's original generator. Two
// draws in three use a computed claim; the original still appears so the
// concept keeps the shape its lesson describes.
type varietyWrap struct {
	id    string
	inner generator.Generator
}

func (g *varietyWrap) Generate(ctx generator.GeneratorContext) generator.Problem {
	if claims := computedVarietyClaims(g.id); len(claims) > 0 && rand.Intn(3) != 0 {
		return conceptPickMC(claims)
	}
	return g.inner.Generate(ctx)
}

// computedVarietyClaims builds one draw's pool of claims for a thin concept.
func computedVarietyClaims(id string) []conceptClaim {
	switch id {
	case "calc.deriv.second_test":
		a := 2 + rand.Intn(6)
		b := 2 + rand.Intn(6)
		k := 1 + rand.Intn(5)
		return []conceptClaim{
			{
				question:    fmt.Sprintf("For f(x)=%dx^2+%dx, is the critical point a local minimum?", a, b),
				correct:     "yes",
				reason:      fmt.Sprintf("f''(x)=%d is positive, so the parabola opens upward and the critical point is a minimum.", 2*a),
				distractors: []string{"no"},
			},
			{
				question:    fmt.Sprintf("For f(x)=-%dx^2+%dx, is the critical point a local maximum?", a, b),
				correct:     "yes",
				reason:      fmt.Sprintf("f''(x)=-%d is negative, so the parabola opens downward and the critical point is a maximum.", 2*a),
				distractors: []string{"no"},
			},
			{
				question:    fmt.Sprintf("If f'(x0)=0 and f''(x0)=%d, what does the second derivative test identify at x0?", k),
				correct:     "local minimum",
				reason:      fmt.Sprintf("A positive second derivative %d means the graph is concave up at x0, so x0 is a local minimum.", k),
				distractors: []string{"local maximum", "inflection point", "none of these"},
			},
			{
				question:    fmt.Sprintf("If f'(x0)=0 and f''(x0)=-%d, what does the second derivative test identify at x0?", k),
				correct:     "local maximum",
				reason:      fmt.Sprintf("A negative second derivative -%d means the graph is concave down at x0, so x0 is a local maximum.", k),
				distractors: []string{"local minimum", "inflection point", "none of these"},
			},
			{
				question:    "If f'(x0)=0 and f''(x0)=0, does the second derivative test decide the type of critical point?",
				correct:     "no",
				reason:      "A zero second derivative leaves the test inconclusive, so another method is needed.",
				distractors: []string{"yes"},
			},
		}
	case "calc.limit.big_o":
		p := 1 + rand.Intn(5)
		d := 1 + rand.Intn(4)
		return []conceptClaim{
			{
				question:    fmt.Sprintf("Is x^{%d} = O(x^{%d}) as x → ∞?", p, p+d),
				correct:     "yes",
				reason:      fmt.Sprintf("x^{%d}/x^{%d} = 1/x^{%d}, which is bounded, so the smaller power is big-O of the larger.", p, p+d, d),
				distractors: []string{"no"},
			},
			{
				question:    fmt.Sprintf("Is x^{%d} = O(x^{%d}) as x → ∞?", p+d, p),
				correct:     "no",
				reason:      fmt.Sprintf("x^{%d}/x^{%d} = x^{%d} grows without bound, so it cannot be bounded by a constant.", p+d, p, d),
				distractors: []string{"yes"},
			},
			{
				question:    fmt.Sprintf("Is %dx^{%d}+x^{%d} = O(x^{%d}) as x → ∞?", p, p, p+d, p+d),
				correct:     "yes",
				reason:      fmt.Sprintf("Both terms are bounded by a multiple of x^{%d}, so the sum is O(x^{%d}).", p+d, p+d),
				distractors: []string{"no"},
			},
			{
				question:    "Is sin(x) = O(1) as x → ∞?",
				correct:     "yes",
				reason:      "The sine function stays between -1 and 1, so a constant bounds it.",
				distractors: []string{"no"},
			},
			{
				question:    fmt.Sprintf("Is x^{%d} = O(1) as x → ∞?", p),
				correct:     "no",
				reason:      fmt.Sprintf("x^{%d} grows without bound, so no constant bounds it.", p),
				distractors: []string{"yes"},
			},
		}
	case "calc.limit.uniform_continuity":
		p := 2 + rand.Intn(4)
		a := 1 + rand.Intn(5)
		b := a + 1 + rand.Intn(5)
		k := 2 + rand.Intn(5)
		return []conceptClaim{
			{
				question:    fmt.Sprintf("Is f(x)=x^{%d} uniformly continuous on [%d, %d]?", p, a, b),
				correct:     "yes",
				reason:      fmt.Sprintf("x^{%d} is continuous on the closed bounded interval [%d, %d], and continuity on a compact set is uniform.", p, a, b),
				distractors: []string{"no"},
			},
			{
				question:    "Is f(x)=x^2 uniformly continuous on [0, ∞)?",
				correct:     "no",
				reason:      "The gap between x^2 and (x+δ)^2 grows with x, so one δ cannot serve every point.",
				distractors: []string{"yes"},
			},
			{
				question:    fmt.Sprintf("Is f(x)=1/(x-%d) uniformly continuous on (%d, %d]?", a, a, b),
				correct:     "no",
				reason:      fmt.Sprintf("The function is unbounded near x=%d, so it cannot be uniformly continuous there.", a),
				distractors: []string{"yes"},
			},
			{
				question:    fmt.Sprintf("Is f(x)=sin(%dx) uniformly continuous on the real line?", k),
				correct:     "yes",
				reason:      fmt.Sprintf("Its derivative %d cos(%dx) is bounded, so the function is Lipschitz and therefore uniformly continuous.", k, k),
				distractors: []string{"no"},
			},
			{
				question:    "Is f(x)=x uniformly continuous on the real line?",
				correct:     "yes",
				reason:      "The identity function is Lipschitz with constant 1, so δ=ε works everywhere.",
				distractors: []string{"no"},
			},
		}
	case "calc.seq.function_sequences":
		k := 2 + rand.Intn(6)
		return []conceptClaim{
			{
				question:    fmt.Sprintf("Does f_n(x)=x^n converge uniformly on [0, 1/%d]?", k),
				correct:     "yes",
				reason:      fmt.Sprintf("On [0, 1/%d] the supremum of x^n is 1/%d^n, which tends to 0.", k, k),
				distractors: []string{"no"},
			},
			{
				question:    "Does f_n(x)=x^n converge uniformly on [0, 1]?",
				correct:     "no",
				reason:      "The pointwise limit jumps at x=1, and the supremum of the remainder stays at 1.",
				distractors: []string{"yes"},
			},
			{
				question:    "Does f_n(x)=x/n converge uniformly on [0, 1]?",
				correct:     "yes",
				reason:      "The supremum of |x/n| is 1/n, which tends to 0.",
				distractors: []string{"no"},
			},
			{
				question:    "Does f_n(x)=n x e^{-n x} converge uniformly on [0, 1]?",
				correct:     "no",
				reason:      "The maximum is 1/e at x=1/n, so the supremum does not tend to 0.",
				distractors: []string{"yes"},
			},
			{
				question:    "Does f_n(x)=x^n converge pointwise on [0, 1]?",
				correct:     "yes",
				reason:      "For each fixed x in [0, 1), x^n tends to 0, and at x=1 it stays 1.",
				distractors: []string{"no"},
			},
		}
	case "calc.seq.monotone":
		k := 2 + rand.Intn(8)
		return []conceptClaim{
			{
				question:    fmt.Sprintf("Is the sequence a_n=n/(n+%d) increasing?", k),
				correct:     "yes",
				reason:      fmt.Sprintf("a_{n+1}-a_n = %d/((n+1+%d)(n+%d)) is positive.", k, k, k),
				distractors: []string{"no"},
			},
			{
				question:    fmt.Sprintf("Is the sequence a_n=%d/n decreasing?", k),
				correct:     "yes",
				reason:      fmt.Sprintf("Each term %d/n is larger than %d/(n+1).", k, k),
				distractors: []string{"no"},
			},
			{
				question:    "Is the sequence a_n=(-1)^n/n monotone?",
				correct:     "no",
				reason:      "The sign alternates, so the sequence moves up and down rather than in one direction.",
				distractors: []string{"yes"},
			},
			{
				question:    "Is every bounded monotone sequence convergent?",
				correct:     "yes",
				reason:      "The monotone convergence theorem guarantees a limit for a bounded monotone sequence.",
				distractors: []string{"no"},
			},
			{
				question:    fmt.Sprintf("Is the sequence a_n=(n+%d)/n decreasing?", k),
				correct:     "yes",
				reason:      fmt.Sprintf("(n+%d)/n = 1+%d/n, and %d/n decreases with n.", k, k, k),
				distractors: []string{"no"},
			},
		}
	case "calc.series.cauchy_criterion":
		p := 2 + rand.Intn(5)
		return []conceptClaim{
			{
				question:    fmt.Sprintf("Does Σ 1/n^{%d} satisfy the Cauchy criterion?", p),
				correct:     "yes",
				reason:      fmt.Sprintf("It is a p-series with exponent %d > 1, so its tail sums tend to 0.", p),
				distractors: []string{"no"},
			},
			{
				question:    "Does Σ 1/n satisfy the Cauchy criterion?",
				correct:     "no",
				reason:      "The harmonic series diverges, so its tail sums do not tend to 0.",
				distractors: []string{"yes"},
			},
			{
				question:    "Does a series converge exactly when its tail sums tend to zero?",
				correct:     "yes",
				reason:      "That is the Cauchy criterion for series, and it is equivalent to convergence in the real numbers.",
				distractors: []string{"no"},
			},
			{
				question:    "Does Σ (-1)^n/n converge?",
				correct:     "yes",
				reason:      "It is the alternating harmonic series, which converges by the alternating series test.",
				distractors: []string{"no"},
			},
			{
				question:    fmt.Sprintf("Does Σ 1/(n^{%d}+1) satisfy the Cauchy criterion?", p),
				correct:     "yes",
				reason:      fmt.Sprintf("The terms are bounded by 1/n^{%d}, whose series converges, so the comparison test gives convergence.", p),
				distractors: []string{"no"},
			},
		}
	case "calc.series.fourier":
		p := 1 + rand.Intn(5)
		return []conceptClaim{
			{
				question:    "Is the Fourier series of an odd function a sine series?",
				correct:     "yes",
				reason:      "All cosine coefficients vanish for an odd function, so only sine terms remain.",
				distractors: []string{"no"},
			},
			{
				question:    "Is the Fourier series of an even function a cosine series?",
				correct:     "yes",
				reason:      "All sine coefficients vanish for an even function, so only cosine terms remain.",
				distractors: []string{"no"},
			},
			{
				question:    "Are the cosine coefficients of an odd function zero?",
				correct:     "yes",
				reason:      "An odd function times a cosine is odd, and its integral over a symmetric interval is zero.",
				distractors: []string{"no"},
			},
			{
				question:    "Are the sine coefficients of an even function zero?",
				correct:     "yes",
				reason:      "An even function times a sine is odd, and its integral over a symmetric interval is zero.",
				distractors: []string{"no"},
			},
			{
				question:    fmt.Sprintf("Does the Fourier series of f(x)=x^{%d} on (-π, π) contain sine terms?", 2*p),
				correct:     "no",
				reason:      fmt.Sprintf("x^{%d} is even, so every sine coefficient vanishes.", 2*p),
				distractors: []string{"yes"},
			},
			{
				question:    fmt.Sprintf("Does the Fourier series of f(x)=x^{%d} on (-π, π) contain cosine terms?", 2*p+1),
				correct:     "no",
				reason:      fmt.Sprintf("x^{%d} is odd, so every cosine coefficient vanishes.", 2*p+1),
				distractors: []string{"yes"},
			},
		}
	case "calc.series.power":
		p := 1 + rand.Intn(5)
		a := 1 + rand.Intn(5)
		return []conceptClaim{
			{
				question:    fmt.Sprintf("What is the radius of convergence of Σ x^n/n^{%d}?", p),
				correct:     "1",
				reason:      fmt.Sprintf("The ratio test gives |x|, so the series converges for |x|<1 regardless of the exponent %d.", p),
				distractors: []string{"0", "2", "infinity"},
			},
			{
				question:    fmt.Sprintf("What is the radius of convergence of Σ (x-%d)^n/n?", a),
				correct:     "1",
				reason:      fmt.Sprintf("The shift to x-%d does not change the ratio limit, which is 1.", a),
				distractors: []string{"0", "2", "infinity"},
			},
			{
				question:    "What is the radius of convergence of Σ n! x^n?",
				correct:     "0",
				reason:      "The ratio is (n+1)|x|, which exceeds 1 for every nonzero x, so only x=0 converges.",
				distractors: []string{"1", "2", "infinity"},
			},
			{
				question:    "What is the radius of convergence of Σ x^n/n!?",
				correct:     "infinity",
				reason:      "The ratio is |x|/(n+1), which tends to 0 for every x, so the series converges everywhere.",
				distractors: []string{"0", "1", "2"},
			},
			{
				question:    "What is the radius of convergence of Σ (x/2)^n?",
				correct:     "2",
				reason:      "The ratio is |x|/2, which is below 1 exactly when |x|<2.",
				distractors: []string{"0", "1", "infinity"},
			},
		}
	case "calc.series.ratio_test":
		r := 1 + rand.Intn(9)
		p := 1 + rand.Intn(4)
		return []conceptClaim{
			{
				question:    fmt.Sprintf("If lim |a_{n+1}/a_n| = %d/10 < 1, does Σ a_n converge?", r),
				correct:     "yes",
				reason:      fmt.Sprintf("A ratio limit of %d/10 is below 1, which is exactly the convergence case of the ratio test.", r),
				distractors: []string{"no"},
			},
			{
				question:    fmt.Sprintf("If lim |a_{n+1}/a_n| = %d/10 > 1, does Σ a_n converge?", 10+r),
				correct:     "no",
				reason:      fmt.Sprintf("A ratio limit of %d/10 exceeds 1, so the terms grow and the series diverges.", 10+r),
				distractors: []string{"yes"},
			},
			{
				question:    "If lim |a_{n+1}/a_n| = 1, is the ratio test conclusive?",
				correct:     "no",
				reason:      "A ratio limit of 1 is the inconclusive case; the series may converge or diverge.",
				distractors: []string{"yes"},
			},
			{
				question:    "Does Σ n!/n^n converge by the ratio test?",
				correct:     "yes",
				reason:      "The ratio tends to 1/e, which is below 1, so the series converges.",
				distractors: []string{"no"},
			},
			{
				question:    fmt.Sprintf("Does Σ n^{%d} pass the ratio test?", p),
				correct:     "no",
				reason:      "The ratio tends to 1, which is the inconclusive case, so the ratio test gives no answer.",
				distractors: []string{"yes"},
			},
		}
	case "calc.series.taylor":
		n := 2 + rand.Intn(5)
		k := 1 + rand.Intn(3)
		nFact := calcFact(n)
		kFact := calcFact(2*k + 1)
		sign := ""
		otherSign := "-"
		if k%2 == 1 {
			sign = "-"
			otherSign = ""
		}
		return []conceptClaim{
			{
				question:    fmt.Sprintf("What is the coefficient of x^{%d} in the Taylor series of e^x at 0?", n),
				correct:     fmt.Sprintf("1/%d", nFact),
				reason:      fmt.Sprintf("The coefficient is 1/n! with n=%d, which is 1/%d.", n, nFact),
				distractors: []string{fmt.Sprintf("1/%d", n+1), fmt.Sprintf("1/%d", n+2), fmt.Sprintf("%d", n)},
			},
			{
				question:    fmt.Sprintf("What is the coefficient of x^{%d} in the Taylor series of sin(x) at 0?", 2*k+1),
				correct:     fmt.Sprintf("%s1/%d", sign, kFact),
				reason:      fmt.Sprintf("The coefficient is (-1)^k/(2k+1)! with k=%d, which is %s1/%d.", k, sign, kFact),
				distractors: []string{fmt.Sprintf("%s1/%d", otherSign, kFact), "1", fmt.Sprintf("1/%d", 2*k+1)},
			},
			{
				question:    "What is the Taylor series of 1/(1-x) at 0?",
				correct:     "sum x^n",
				reason:      "The geometric series 1+x+x^2+... equals 1/(1-x) for |x|<1.",
				distractors: []string{"sum x^n/n!", "sum (-1)^n x^n", "sum x^n/n"},
			},
			{
				question:    "What is the constant term of the Taylor series of cos(x) at 0?",
				correct:     "1",
				reason:      "cos(0)=1, so the constant term is 1.",
				distractors: []string{"0", "-1", "1/2"},
			},
			{
				question:    "What is the coefficient of x^2 in the Taylor series of cos(x) at 0?",
				correct:     "-1/2",
				reason:      "The x^2 coefficient is -1/2!, which is -1/2.",
				distractors: []string{"1/2", "-1", "1"},
			},
		}
	case "calc.limit.indeterminate":
		forms := []struct {
			form    string
			verdict string
			reason  string
		}{
			{"0/0", "yes", "Both the numerator and denominator tend to zero, so the quotient can approach different values along different paths."},
			{"∞/∞", "yes", "Both parts grow without bound, and their ratio depends on the relative rates."},
			{"0·∞", "yes", "One factor shrinks while the other grows, and the product can approach any value."},
			{"∞-∞", "yes", "The two large terms can cancel to different degrees, so the difference is not determined."},
			{"0^0", "yes", "The base and the exponent compete, and the limit depends on how each approaches its target."},
			{"1^∞", "yes", "The base approaches one while the exponent grows, so the result depends on the rates."},
			{"∞+∞", "no", "Both terms grow without bound, so their sum also grows without bound."},
			{"0/1", "no", "The denominator tends to 1, so the quotient tends to 0."},
		}
		f := forms[rand.Intn(len(forms))]
		return []conceptClaim{
			{
				question:    fmt.Sprintf("Is %s an indeterminate form?", f.form),
				correct:     f.verdict,
				reason:      f.reason,
				distractors: []string{opposite(f.verdict)},
			},
		}
	case "calc.multi.lagrange_two":
		s := 2 * (2 + rand.Intn(6))
		r := 1 + rand.Intn(6)
		return []conceptClaim{
			{
				question:    fmt.Sprintf("Maximize x+y subject to x^2+y^2=%d. What is the maximum?", r*r),
				correct:     fmt.Sprintf("%d√2", r),
				reason:      fmt.Sprintf("Lagrange gives x=y=±%d/√2, so x+y reaches %d√2.", r, r),
				distractors: []string{fmt.Sprintf("%d", 2*r), fmt.Sprintf("%d", r), fmt.Sprintf("%d/2", r)},
			},
			{
				question:    fmt.Sprintf("Minimize x^2+y^2 subject to x+y=%d. What is the minimum?", s),
				correct:     fmt.Sprintf("%d", s*s/2),
				reason:      fmt.Sprintf("The closest point is x=y=%d, so x^2+y^2 = %d.", s/2, s*s/2),
				distractors: []string{fmt.Sprintf("%d", s*s), fmt.Sprintf("%d", s*s/4), fmt.Sprintf("%d", s)},
			},
			{
				question:    fmt.Sprintf("Maximize xy subject to x+y=%d. What is the maximum?", s),
				correct:     fmt.Sprintf("%d", s*s/4),
				reason:      fmt.Sprintf("The product is largest at x=y=%d, giving %d.", s/2, s*s/4),
				distractors: []string{fmt.Sprintf("%d", s*s/2), fmt.Sprintf("%d", s), fmt.Sprintf("%d", s*s)},
			},
			{
				question:    "Maximize x+y subject to x^2+y^2=2. What is the maximum?",
				correct:     "2",
				reason:      "The maximum is √2·√2 = 2, attained at x=y=1.",
				distractors: []string{"1", "√2", "4"},
			},
		}
	case "calc.seq.convergence":
		k := 2 + rand.Intn(8)
		p := 1 + rand.Intn(3)
		return []conceptClaim{
			{
				question:    fmt.Sprintf("Does the sequence a_n=n/(n+%d) converge?", k),
				correct:     "yes",
				reason:      fmt.Sprintf("Dividing by n gives 1/(1+%d/n), which tends to 1.", k),
				distractors: []string{"no"},
			},
			{
				question:    "Does the sequence a_n=(-1)^n converge?",
				correct:     "no",
				reason:      "The terms alternate between 1 and -1, so there is no single limit.",
				distractors: []string{"yes"},
			},
			{
				question:    "Does the sequence a_n=(1+1/n)^n converge?",
				correct:     "yes",
				reason:      "It converges to e, the base of the natural logarithm.",
				distractors: []string{"no"},
			},
			{
				question:    fmt.Sprintf("Does the sequence a_n=n^{%d} converge?", p),
				correct:     "no",
				reason:      fmt.Sprintf("n^{%d} grows without bound, so the sequence diverges.", p),
				distractors: []string{"yes"},
			},
			{
				question:    fmt.Sprintf("Does the sequence a_n=%d/n converge?", k),
				correct:     "yes",
				reason:      fmt.Sprintf("%d/n tends to 0 as n grows.", k),
				distractors: []string{"no"},
			},
		}
	}
	return nil
}
