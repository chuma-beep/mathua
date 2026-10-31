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
	case "calc.integral.reduction":
		n := 2 + rand.Intn(5)
		m := 3 + rand.Intn(4)
		return []conceptClaim{
			{
				question: fmt.Sprintf("Which reduction formula is correct for ∫x^%d e^x dx?", n),
				correct:  fmt.Sprintf("x^%d e^x - %d∫x^%d e^x dx", n, n, n-1),
				reason:   "Integrate by parts with u=x^n and dv=e^x dx.",
				distractors: []string{
					fmt.Sprintf("x^%d e^x + %d∫x^%d e^x dx", n, n, n-1),
					fmt.Sprintf("%d x^%d e^x - ∫x^%d e^x dx", n, n, n-1),
					fmt.Sprintf("x^%d e^x - ∫x^%d e^x dx", n, n),
				},
			},
			{
				question: fmt.Sprintf("Which reduction formula is correct for ∫(ln x)^%d dx?", n),
				correct:  fmt.Sprintf("x(ln x)^%d - %d∫(ln x)^%d dx", n, n, n-1),
				reason:   "Integrate by parts with u=(ln x)^n and dv=dx.",
				distractors: []string{
					fmt.Sprintf("x(ln x)^%d + %d∫(ln x)^%d dx", n, n, n-1),
					fmt.Sprintf("(ln x)^%d/(%d)", n+1, n+1),
					fmt.Sprintf("x(ln x)^%d - %d∫(ln x)^%d dx", n-1, n, n),
				},
			},
			{
				question: fmt.Sprintf("For n=%d, which reduction formula is correct for ∫sin^n x dx?", m),
				correct:  fmt.Sprintf("-(1/%d)sin^%d x cos x + (%d/%d)∫sin^%d x dx", m, m-1, m-1, m, m-2),
				reason:   "Rewrite sin^n x as sin^(n-1)x sin x and integrate by parts.",
				distractors: []string{
					fmt.Sprintf("(1/%d)sin^%d x cos x + (%d/%d)∫sin^%d x dx", m, m-1, m-1, m, m-2),
					fmt.Sprintf("-(1/%d)sin^%d x cos x - (%d/%d)∫sin^%d x dx", m, m-1, m-1, m, m-2),
					fmt.Sprintf("-(1/%d)sin^%d x cos x + (%d/%d)∫sin^%d x dx", m, m, m-1, m, m),
				},
			},
			{
				question: fmt.Sprintf("For n=%d, which reduction formula is correct for ∫cos^n x dx?", m),
				correct:  fmt.Sprintf("(1/%d)cos^%d x sin x + (%d/%d)∫cos^%d x dx", m, m-1, m-1, m, m-2),
				reason:   "Rewrite cos^n x as cos^(n-1)x cos x and integrate by parts.",
				distractors: []string{
					fmt.Sprintf("-(1/%d)cos^%d x sin x + (%d/%d)∫cos^%d x dx", m, m-1, m-1, m, m-2),
					fmt.Sprintf("(1/%d)cos^%d x sin x - (%d/%d)∫cos^%d x dx", m, m-1, m-1, m, m-2),
					fmt.Sprintf("(1/%d)cos^%d x sin x + (%d/%d)∫cos^%d x dx", m, m, m-1, m, m),
				},
			},
			{
				question: fmt.Sprintf("For n=%d, which reduction formula is correct for ∫sec^n x dx?", m),
				correct:  fmt.Sprintf("(1/%d)sec^%d x tan x + (%d/%d)∫sec^%d x dx", m-1, m-2, m-2, m-1, m-2),
				reason:   "Split off sec^2 x and integrate by parts.",
				distractors: []string{
					fmt.Sprintf("(1/%d)sec^%d x tan x + (%d/%d)∫sec^%d x dx", m-1, m-2, m-1, m-1, m-2),
					fmt.Sprintf("(1/%d)sec^%d x tan x + (%d/%d)∫sec^%d x dx", m, m-1, m-1, m, m-2),
					fmt.Sprintf("(1/%d)sec^%d x tan x + (%d/%d)∫sec^%d x dx", m-1, m-1, m-2, m-1, m-2),
				},
			},
		}
	case "calc.integral.strategy":
		a := 2 + rand.Intn(6)
		return []conceptClaim{
			{
				question:    "Which method is best for ∫x e^(x^2) dx?",
				correct:     "substitution",
				reason:      "Let u=x^2, so du=2x dx.",
				distractors: []string{"integration by parts", "partial fractions", "trigonometric substitution"},
			},
			{
				question:    fmt.Sprintf("Which method is best for ∫1/(x^2+%d) dx?", a),
				correct:     "trigonometric substitution",
				reason:      "Let x=a tan θ to use 1+tan^2 θ=sec^2 θ.",
				distractors: []string{"substitution", "partial fractions", "integration by parts"},
			},
			{
				question:    "Which method is best for ∫x^2/(x+1) dx?",
				correct:     "polynomial division",
				reason:      "Divide x^2 by x+1 first, then integrate the quotient and remainder.",
				distractors: []string{"partial fractions", "substitution", "integration by parts"},
			},
			{
				question:    fmt.Sprintf("Which method is best for ∫1/(x^2-%d) dx?", a),
				correct:     "partial fractions",
				reason:      "Factor x^2-a^2 and split into two linear denominators.",
				distractors: []string{"trigonometric substitution", "substitution", "integration by parts"},
			},
			{
				question:    "Which method is best for ∫x ln x dx?",
				correct:     "integration by parts",
				reason:      "Take u=ln x and dv=x dx.",
				distractors: []string{"substitution", "partial fractions", "trigonometric substitution"},
			},
			{
				question:    "Which method is best for ∫sin^3 x cos^2 x dx?",
				correct:     "Pythagorean identity",
				reason:      "Rewrite sin^3 x as sin x(1-cos^2 x) and substitute u=cos x.",
				distractors: []string{"integration by parts", "partial fractions", "trigonometric substitution"},
			},
		}
	case "calc.integral.trig_integrals":
		m := 3 + 2*rand.Intn(2)
		nEven := 2 + 2*rand.Intn(3)
		nOdd := 3 + 2*rand.Intn(2)
		mEven := 2 + 2*rand.Intn(3)
		return []conceptClaim{
			{
				question:    fmt.Sprintf("For ∫sin^%d x cos^%d x dx, which substitution works?", m, nEven),
				correct:     "u=cos x",
				reason:      "The sine power is odd, so save one sin x for du.",
				distractors: []string{"u=sin x", "half-angle identities", "integration by parts"},
			},
			{
				question:    fmt.Sprintf("For ∫sin^%d x cos^%d x dx, which substitution works?", mEven, nOdd),
				correct:     "u=sin x",
				reason:      "The cosine power is odd, so save one cos x for du.",
				distractors: []string{"u=cos x", "half-angle identities", "integration by parts"},
			},
			{
				question:    fmt.Sprintf("For ∫sin^%d x cos^%d x dx, which approach works?", mEven, nEven),
				correct:     "half-angle identities",
				reason:      "Both powers are even, so lower the degree with sin^2 x=(1-cos 2x)/2.",
				distractors: []string{"u=sin x", "u=cos x", "integration by parts"},
			},
			{
				question:    fmt.Sprintf("For ∫sin^%d x dx with n odd, what is the first step?", m),
				correct:     "save one sin x and substitute u=cos x",
				reason:      "Then sin^2 x becomes 1-u^2.",
				distractors: []string{"save one cos x and substitute u=sin x", "use a half-angle identity", "integrate by parts"},
			},
			{
				question:    fmt.Sprintf("For ∫cos^%d x dx with n even, what is the first step?", nEven),
				correct:     "use a half-angle identity",
				reason:      "cos^2 x=(1+cos 2x)/2 lowers the power.",
				distractors: []string{"substitute u=sin x", "substitute u=cos x", "integrate by parts"},
			},
		}
	case "calc.deriv.implicit_second":
		r := 2 + rand.Intn(6)
		k := 2 + rand.Intn(8)
		return []conceptClaim{
			{
				question:    fmt.Sprintf("If x^2+y^2=%d, what is y''?", r),
				correct:     fmt.Sprintf("-%d/y^3", r),
				reason:      "Differentiate implicitly twice and use x^2+y^2=r^2.",
				distractors: []string{fmt.Sprintf("-%d/y^2", r), fmt.Sprintf("%d/y^3", r), "-2x/y^3"},
			},
			{
				question:    fmt.Sprintf("If x^2-y^2=%d, what is y''?", r),
				correct:     fmt.Sprintf("-%d/y^3", r),
				reason:      "Differentiate twice and use x^2-y^2=r^2.",
				distractors: []string{fmt.Sprintf("-%d/y^2", r), fmt.Sprintf("%d/y^3", r), "2x/y^3"},
			},
			{
				question:    fmt.Sprintf("If xy=%d, what is y''?", k),
				correct:     fmt.Sprintf("2*%d/x^3", k),
				reason:      "From y=k/x, differentiate twice.",
				distractors: []string{fmt.Sprintf("-2*%d/x^3", k), fmt.Sprintf("%d/x^3", k), fmt.Sprintf("2*%d/x^2", k)},
			},
			{
				question:    fmt.Sprintf("If x^2+y^2=%d, what is y' at (0,%d)?", r, r),
				correct:     "0",
				reason:      "At the top of the circle the tangent is horizontal.",
				distractors: []string{"1", "-1", "undefined"},
			},
			{
				question:    fmt.Sprintf("If xy=%d, what is y' at x=1?", k),
				correct:     fmt.Sprintf("-%d", k),
				reason:      "y=k/x, so y'=-k/x^2 and at x=1 it is -k.",
				distractors: []string{fmt.Sprintf("%d", k), "0", fmt.Sprintf("-%d", k*k)},
			},
		}
	case "calc.integral.convergence":
		p := 2 + rand.Intn(4)
		q := 2 + rand.Intn(3)
		a := 2 + rand.Intn(6)
		return []conceptClaim{
			{
				question:    fmt.Sprintf("Does ∫_1^∞ 1/x^%d dx converge?", p),
				correct:     "yes",
				reason:      fmt.Sprintf("p=%d is greater than 1, so the p-integral converges.", p),
				distractors: []string{"no"},
			},
			{
				question:    fmt.Sprintf("Does ∫_0^1 1/x^%d dx converge?", q),
				correct:     "no",
				reason:      fmt.Sprintf("For the interval (0,1], convergence needs p<1, and p=%d.", q),
				distractors: []string{"yes"},
			},
			{
				question:    fmt.Sprintf("Does ∫_1^∞ e^(-%d x) dx converge?", a),
				correct:     "yes",
				reason:      fmt.Sprintf("The integral equals 1/%d, a finite value.", a),
				distractors: []string{"no"},
			},
			{
				question:    fmt.Sprintf("Does ∫_0^1 1/(1-x)^%d dx converge?", q),
				correct:     "no",
				reason:      fmt.Sprintf("The singularity at x=1 behaves like 1/t^%d with p=%d, so it diverges.", q, q),
				distractors: []string{"yes"},
			},
			{
				question:    fmt.Sprintf("Does ∫_2^∞ 1/(x(ln x)^%d) dx converge?", p),
				correct:     "yes",
				reason:      fmt.Sprintf("The substitution u=ln x gives a p-integral with p=%d>1.", p),
				distractors: []string{"no"},
			},
		}
	case "calc.limit.continuity":
		a := 2 + rand.Intn(6)
		a2 := a * a
		return []conceptClaim{
			{
				question:    fmt.Sprintf("What type of discontinuity does f(x)=(x^2-%d)/(x-%d) have at x=%d?", a2, a, a),
				correct:     "removable",
				reason:      "The factor x-a cancels, leaving a hole at x=a.",
				distractors: []string{"jump", "infinite", "continuous"},
			},
			{
				question:    fmt.Sprintf("What type of discontinuity does f(x)=1/(x-%d) have at x=%d?", a, a),
				correct:     "infinite",
				reason:      "The denominator tends to zero while the numerator stays 1.",
				distractors: []string{"removable", "jump", "continuous"},
			},
			{
				question:    fmt.Sprintf("What type of discontinuity does f(x)=|x-%d|/(x-%d) have at x=%d?", a, a, a),
				correct:     "jump",
				reason:      "The one-sided limits are -1 and 1.",
				distractors: []string{"removable", "infinite", "continuous"},
			},
			{
				question:    fmt.Sprintf("What type of discontinuity does f(x)=1/(x-%d)^2 have at x=%d?", a, a),
				correct:     "infinite",
				reason:      "The function grows without bound on both sides.",
				distractors: []string{"removable", "jump", "continuous"},
			},
			{
				question:    fmt.Sprintf("Where is f(x)=x^2+%d continuous?", a),
				correct:     "everywhere",
				reason:      "Every polynomial is continuous on the real line.",
				distractors: []string{fmt.Sprintf("only at x=%d", a), "nowhere", "only at x=0"},
			},
		}
	case "calc.series.root_test":
		k := 2 + rand.Intn(5)
		a := 1 + rand.Intn(4)
		b := 2 + rand.Intn(4)
		if a == b {
			b++
		}
		return []conceptClaim{
			{
				question:    fmt.Sprintf("Apply the root test to ∑(n/(%d n+1))^n. What is L?", k),
				correct:     fmt.Sprintf("1/%d", k),
				reason:      "The n-th root tends to 1/k, which is less than 1.",
				distractors: []string{"1", fmt.Sprintf("%d", k), "0"},
			},
			{
				question:    "Apply the root test to ∑(n/(n+1))^(n^2). What is L?",
				correct:     "1/e",
				reason:      "The n-th root is (n/(n+1))^n, which tends to 1/e.",
				distractors: []string{"1", "e", "0"},
			},
			{
				question:    "Apply the root test to ∑(1/n)^n. What is L?",
				correct:     "0",
				reason:      "The n-th root is 1/n, which tends to 0.",
				distractors: []string{"1", "e", "infinite"},
			},
			{
				question:    fmt.Sprintf("Apply the root test to ∑(%d n/(%d n+1))^n. What is L?", a, b),
				correct:     fmt.Sprintf("%d/%d", a, b),
				reason:      "The n-th root tends to a/b, the ratio of the leading coefficients.",
				distractors: []string{"1", fmt.Sprintf("%d/%d", b, a), "0"},
			},
			{
				question:    "If the root test gives L=1, what can you conclude?",
				correct:     "the test is inconclusive",
				reason:      "The root test only decides convergence when L is not 1.",
				distractors: []string{"the series converges", "the series diverges", "the series has sum 1"},
			},
		}
	case "calc.integral.irrational":
		a := 2 + rand.Intn(6)
		a2 := a * a
		return []conceptClaim{
			{
				question:    fmt.Sprintf("What is ∫dx/sqrt(%d-x^2)?", a2),
				correct:     fmt.Sprintf("arcsin(x/%d)+C", a),
				reason:      "This is the standard arcsine form with a=%d.",
				distractors: []string{fmt.Sprintf("arcsin(x/%d)+C", a2), fmt.Sprintf("arctan(x/%d)+C", a), fmt.Sprintf("(1/%d)arcsin(x/%d)+C", a, a)},
			},
			{
				question:    fmt.Sprintf("What is ∫dx/(x^2+%d)?", a2),
				correct:     fmt.Sprintf("(1/%d)arctan(x/%d)+C", a, a),
				reason:      "Factor out a^2 and use the arctangent form.",
				distractors: []string{fmt.Sprintf("arctan(x/%d)+C", a), fmt.Sprintf("(1/%d)arctan(x/%d)+C", a, a2), fmt.Sprintf("(1/%d)arcsin(x/%d)+C", a, a)},
			},
			{
				question:    fmt.Sprintf("What is ∫dx/sqrt(x^2+%d)?", a2),
				correct:     fmt.Sprintf("ln|x+sqrt(x^2+%d)|+C", a2),
				reason:      "This is the standard logarithmic form.",
				distractors: []string{fmt.Sprintf("ln|x+sqrt(x^2+%d)|+C", a), fmt.Sprintf("arcsin(x/%d)+C", a), fmt.Sprintf("(1/%d)ln|x+sqrt(x^2+%d)|+C", a, a2)},
			},
			{
				question: fmt.Sprintf("What is ∫sqrt(%d-x^2) dx?", a2),
				correct:  fmt.Sprintf("(x/2)sqrt(%d-x^2)+(%d/2)arcsin(x/%d)+C", a2, a2, a),
				reason:   "Use x=a sin θ and the double-angle identity.",
				distractors: []string{
					fmt.Sprintf("(x/2)sqrt(%d-x^2)+(%d/2)arcsin(x/%d)+C", a2, a, a),
					fmt.Sprintf("(x/2)sqrt(%d-x^2)+(%d/2)arctan(x/%d)+C", a2, a2, a),
					fmt.Sprintf("x sqrt(%d-x^2)+%d arcsin(x/%d)+C", a2, a, a),
				},
			},
			{
				question:    fmt.Sprintf("What is ∫dx/(x sqrt(x^2-%d))?", a2),
				correct:     fmt.Sprintf("(1/%d)arcsec(x/%d)+C", a, a),
				reason:      "Use x=a sec θ, so sqrt(x^2-a^2)=a tan θ.",
				distractors: []string{fmt.Sprintf("arcsec(x/%d)+C", a), fmt.Sprintf("(1/%d)arcsin(x/%d)+C", a, a), fmt.Sprintf("(1/%d)arctan(x/%d)+C", a, a)},
			},
		}
	case "calc.series.absolute":
		p := 2 + rand.Intn(3)
		k := 2 + rand.Intn(6)
		return []conceptClaim{
			{
				question:    fmt.Sprintf("Does ∑(-1)^n/n^%d converge absolutely, conditionally, or diverge?", p),
				correct:     "absolutely",
				reason:      fmt.Sprintf("The absolute series is a p-series with p=%d>1.", p),
				distractors: []string{"conditionally", "diverges"},
			},
			{
				question:    "Does ∑(-1)^n/n converge absolutely, conditionally, or diverge?",
				correct:     "conditionally",
				reason:      "It converges by the alternating series test, but the harmonic series diverges.",
				distractors: []string{"absolutely", "diverges"},
			},
			{
				question:    "Does ∑(-1)^n/n^2 converge absolutely, conditionally, or diverge?",
				correct:     "absolutely",
				reason:      "The absolute series is the convergent p-series with p=2.",
				distractors: []string{"conditionally", "diverges"},
			},
			{
				question:    "Does ∑(-1)^n/sqrt(n) converge absolutely, conditionally, or diverge?",
				correct:     "conditionally",
				reason:      "The alternating series converges, but the p-series with p=1/2 diverges.",
				distractors: []string{"absolutely", "diverges"},
			},
			{
				question:    fmt.Sprintf("Does ∑(-1)^n/(n^2+%d) converge absolutely, conditionally, or diverge?", k),
				correct:     "absolutely",
				reason:      "The absolute series is comparable to 1/n^2.",
				distractors: []string{"conditionally", "diverges"},
			},
		}
	case "calc.limit.numeric":
		a := 2 + rand.Intn(5)
		k := 2 + rand.Intn(5)
		c := 2 + rand.Intn(4)
		m := 2 + rand.Intn(4)
		num := c * m
		return []conceptClaim{
			{
				question:    fmt.Sprintf("What is lim_{x→%d} (x^2-%d)/(x-%d)?", a, a*a, a),
				correct:     fmt.Sprintf("%d", 2*a),
				reason:      "Factor the numerator as (x-a)(x+a) and cancel x-a.",
				distractors: []string{fmt.Sprintf("%d", a), fmt.Sprintf("%d", a*a), "0"},
			},
			{
				question:    fmt.Sprintf("What is lim_{x→%d} (x^3-%d)/(x-%d)?", a, a*a*a, a),
				correct:     fmt.Sprintf("%d", 3*a*a),
				reason:      "Factor x^3-a^3 and cancel x-a.",
				distractors: []string{fmt.Sprintf("%d", a*a), fmt.Sprintf("%d", a*a*a), "0"},
			},
			{
				question:    fmt.Sprintf("What is lim_{x→0} sin(%d x)/x?", k),
				correct:     fmt.Sprintf("%d", k),
				reason:      "Use sin(kx)/x = k·sin(kx)/(kx) and the standard limit.",
				distractors: []string{"1", fmt.Sprintf("%d", k*k), "0"},
			},
			{
				question:    fmt.Sprintf("What is lim_{x→0} (e^(%d x)-1)/x?", k),
				correct:     fmt.Sprintf("%d", k),
				reason:      "This is the derivative of e^(kx) at x=0.",
				distractors: []string{"1", "0", fmt.Sprintf("%d", k*k)},
			},
			{
				question:    fmt.Sprintf("What is lim_{x→∞} (%d x^2+3)/(%d x^2+1)?", num, c),
				correct:     fmt.Sprintf("%d", m),
				reason:      "The ratio of leading coefficients is a/c.",
				distractors: []string{fmt.Sprintf("%d", c), fmt.Sprintf("%d", num), "0"},
			},
		}
	case "calc.deriv.optimization":
		k := 2 + rand.Intn(5)
		return []conceptClaim{
			{
				question:    fmt.Sprintf("A rectangle has perimeter %d. What is its maximum area?", 4*k),
				correct:     fmt.Sprintf("%d", k*k),
				reason:      "A square with side k has the largest area for a fixed perimeter.",
				distractors: []string{fmt.Sprintf("%d", 2*k), fmt.Sprintf("%d", 4*k), fmt.Sprintf("%d", k)},
			},
			{
				question:    fmt.Sprintf("A square has area %d. What is its minimum perimeter?", k*k),
				correct:     fmt.Sprintf("%d", 4*k),
				reason:      "The side is k, so the perimeter is 4k.",
				distractors: []string{fmt.Sprintf("%d", 2*k), fmt.Sprintf("%d", k*k), fmt.Sprintf("%d", 8*k)},
			},
			{
				question:    fmt.Sprintf("An open box is made from a square sheet of side %d by cutting squares of side x from each corner. What x maximizes the volume?", 6*k),
				correct:     fmt.Sprintf("%d", k),
				reason:      "Volume is x(6k-2x)^2, whose maximum is at x=k.",
				distractors: []string{fmt.Sprintf("%d", 2*k), fmt.Sprintf("%d", 3*k), fmt.Sprintf("%d", 6*k)},
			},
			{
				question:    fmt.Sprintf("A fence along a river encloses a rectangle using %d meters of fence on three sides. What is the maximum area?", 4*k),
				correct:     fmt.Sprintf("%d", 2*k*k),
				reason:      "If the side along the river is 2k, the other sides are k each.",
				distractors: []string{fmt.Sprintf("%d", k*k), fmt.Sprintf("%d", 4*k*k), fmt.Sprintf("%d", 2*k)},
			},
			{
				question:    fmt.Sprintf("Two positive numbers sum to %d. What is their maximum product?", 2*k),
				correct:     fmt.Sprintf("%d", k*k),
				reason:      "The product is largest when the numbers are equal.",
				distractors: []string{fmt.Sprintf("%d", 2*k), fmt.Sprintf("%d", 4*k), fmt.Sprintf("%d", k)},
			},
		}
	}
	return nil
}
