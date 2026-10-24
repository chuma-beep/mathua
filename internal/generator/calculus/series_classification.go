package calculus

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

// Series classification is the natural question for these concepts, and a
// classification answer cannot be typed into a maths field, so the
// yes/no shapes are multiple choice. Every pool mixes answers in both
// directions: an easy pool where every answer is "yes" lets a learner score
// without reading the question, which the pool discrimination test rejects.

func (g *seriesConceptGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := randRange(ctx.Difficulty, 1, 9)
	d := 2 * randRange(ctx.Difficulty, 1, 3) // even, so the partial sum is an integer
	n := randRange(ctx.Difficulty, 3, 12)
	if rand.Intn(2) == 0 {
		sum := n * (2*a + (n-1)*d) / 2
		return generator.Problem{
			Question:    fmt.Sprintf("An arithmetic sequence starts at %d with common difference %d. What is its partial sum S_%d? (enter a number)", a, d, n),
			Answer:      strconv.Itoa(sum),
			Explanation: fmt.Sprintf("S_%d = %d/2 times (2·%d + %d·%d) = %d. A partial sum is the sum of the first %d terms, and for an arithmetic sequence it is the term count times the average of the first and last terms.", n, n, a, n-1, d, sum, n),
			Facts:       map[string]string{"a": strconv.Itoa(a), "d": strconv.Itoa(d), "n": strconv.Itoa(n), "sum": strconv.Itoa(sum)},
		}
	}
	if rand.Intn(2) == 0 {
		r := generator.MultipleChoice("Does a series converge exactly when its sequence of partial sums converges?", "yes", "no")
		r.Explanation = "The sum of a series is defined as the limit of its partial sums. If that limit exists and is finite the series converges, so question and answer are the same statement."
		return r
	}
	r := generator.MultipleChoice("Does every series converge?", "no", "yes")
	r.Explanation = "No. The harmonic series has terms tending to zero and still diverges, so neither a small term nor a long list of them forces convergence."
	return r
}

func (g *harmonicSeriesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	p := randRange(ctx.Difficulty, 1, 6)
	if rand.Intn(3) == 0 {
		converges := p > 1
		answer, other := "yes", "no"
		if !converges {
			answer, other = "no", "yes"
		}
		r := generator.MultipleChoice(fmt.Sprintf("Does the p-series \\(\\sum 1/n^{%d}\\) converge?", p), answer, other)
		if converges {
			r.Explanation = fmt.Sprintf("With p = %d > 1 the terms shrink fast enough for the sum to stay bounded, so the p-series converges.", p)
		} else {
			r.Explanation = "With p = 1 this is the harmonic series. Its terms tend to zero, yet the partial sums grow like the natural logarithm, so it diverges. p = 1 is the boundary case."
		}
		return r
	}
	if rand.Intn(3) == 0 {
		n := randRange(ctx.Difficulty, 2, 4)
		h := []string{"", "", "3/2", "11/6", "25/12"}[n]
		return generator.Problem{
			Question:    fmt.Sprintf("What is the harmonic number H_%d? (enter a fraction)", n),
			Answer:      h,
			Explanation: fmt.Sprintf("H_%d is the sum of 1/k for k from 1 to %d, which is %s. These partial sums keep growing, slowly, without bound.", n, n, h),
			Facts:       map[string]string{"n": strconv.Itoa(n), "h": h},
		}
	}
	if rand.Intn(2) == 0 {
		c := randRange(ctx.Difficulty, 2, 9)
		answer, other := "yes", "no"
		if p <= 1 {
			answer, other = "no", "yes"
		}
		r := generator.MultipleChoice(fmt.Sprintf("Does \\(\\sum %d/n^{%d}\\) converge?", c, p), answer, other)
		r.Explanation = "A constant factor does not change convergence, so this behaves exactly like the p-series with the same exponent."
		return r
	}
	if rand.Intn(2) == 0 {
		r := generator.MultipleChoice("Do the harmonic numbers H_n stay bounded as n grows?", "no", "yes")
		r.Explanation = "No. H_n grows like the natural logarithm, so it increases without bound even though each added term gets smaller."
		return r
	}
	r := generator.MultipleChoice("If the terms of a series tend to zero, must the series converge?", "no", "yes")
	r.Explanation = "No. Terms tending to zero is necessary but not sufficient, and the harmonic series is the standard counterexample."
	return r
}

func (g *alternatingSeriesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	p := randRange(ctx.Difficulty, 1, 3)
	switch rand.Intn(3) {
	case 0:
		r := generator.MultipleChoice(fmt.Sprintf("Does \\(\\sum (-1)^{n}/n^{%d}\\) converge?", p), "yes", "no")
		r.Explanation = fmt.Sprintf("The magnitudes 1/n^%d decrease to zero, so the alternating series test applies and the series converges.", p)
		return r
	case 1:
		r := generator.MultipleChoice(fmt.Sprintf("Does \\(\\sum (-1)^{n}n^{%d}\\) converge?", p), "no", "yes")
		r.Explanation = fmt.Sprintf("The terms grow in magnitude like n^%d, so they do not tend to zero. A series whose terms fail to tend to zero diverges by the nth-term test.", p)
		return r
	}
	if rand.Intn(2) == 0 {
		n := randRange(ctx.Difficulty, 4, 20)
		r := generator.MultipleChoice(fmt.Sprintf("Does \\(\\sum (-1)^{n}/(n+%d)\\) converge?", n), "yes", "no")
		r.Explanation = "The magnitudes 1/(n+%d) decrease to zero, so the alternating series test applies."
		return r
	}
	if rand.Intn(2) == 0 {
		n := randRange(ctx.Difficulty, 4, 20)
		p := randRange(ctx.Difficulty, 1, 3)
		den := intPow(n+1, p)
		answer := fmt.Sprintf("1/%d", den)
		if den == 1 {
			answer = "1"
		}
		return generator.Problem{
			Question:    fmt.Sprintf("Stopping \\(\\sum (-1)^{n}/n^{%d}\\) after the term with n = %d leaves an error bounded by which first omitted term? (enter a fraction)", p, n),
			Answer:      answer,
			Explanation: fmt.Sprintf("For a convergent alternating series with decreasing magnitudes the error is at most the first omitted term, \\(1/%d^{%d}\\).", n+1, p),
			Facts:       map[string]string{"n": fmt.Sprintf("%d", n), "p": fmt.Sprintf("%d", p)},
		}
	}
	r := generator.MultipleChoice("Does the alternating series test require the magnitudes to decrease to zero?", "yes", "no")
	r.Explanation = "Yes. Two conditions are needed: the magnitudes decrease monotonically, and they tend to zero. Losing either one leaves the test silent."
	return r
}

func (g *integralTestGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	p := randRange(ctx.Difficulty, 1, 6)
	if rand.Intn(2) == 0 {
		answer, other := "yes", "no"
		if p <= 1 {
			answer, other = "no", "yes"
		}
		r := generator.MultipleChoice(fmt.Sprintf("Using the integral test, does \\(\\sum 1/n^{%d}\\) converge?", p), answer, other)
		if p > 1 {
			r.Explanation = fmt.Sprintf("The integral of 1/x^%d from 1 to infinity is finite, so the series converges. The function is positive, continuous and decreasing on the interval the test needs.", p)
		} else {
			r.Explanation = "The integral of 1/x from 1 to infinity diverges, so the series diverges. The integral test transfers the question to an improper integral that can be evaluated directly."
		}
		return r
	}
	if rand.Intn(2) == 0 {
		c := randRange(ctx.Difficulty, 2, 7)
		r := generator.MultipleChoice(fmt.Sprintf("Does \\(\\sum \\frac{1}{(n+%d)^{%d}}\\) converge?", c, p), "yes", "no")
		r.Explanation = fmt.Sprintf("Shifting the index does not change convergence, so this behaves like the p-series with p = %d, which converges because %d > 1.", p, p)
		return r
	}
	// The integral itself, which is the quantity the test actually evaluates.
	q := fmt.Sprintf("Evaluate \\(\\int_{1}^{\\infty} \\frac{1}{x^{%d}}\\,dx\\). (enter a fraction or a number)", p)
	if p == 1 {
		r := generator.MultipleChoice("Does \\(\\int_{1}^{\\infty} \\frac{1}{x}\\,dx\\) converge?", "no", "yes")
		r.Explanation = "The antiderivative is the natural logarithm, which grows without bound, so the improper integral diverges. That is exactly why the harmonic series diverges."
		return r
	}
	answer := fmt.Sprintf("1/%d", p-1)
	if p == 2 {
		answer = "1"
	}
	return generator.Problem{
		Question:    q,
		Answer:      answer,
		Explanation: fmt.Sprintf("The antiderivative is \\(x^{%d}/(%d)\\), and at infinity it vanishes, while at 1 it contributes 1/(%d).", 1-p, 1-p, p-1),
		Facts:       map[string]string{"p": strconv.Itoa(p)},
	}
}

func (g *nthTermGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	c := randRange(ctx.Difficulty, 2, 5)
	k := randRange(ctx.Difficulty, 1, 4)
	a := c * k
	b := randRange(ctx.Difficulty, 1, 9)
	d := randRange(ctx.Difficulty, 1, 9)
	switch rand.Intn(3) {
	case 0:
		return generator.Problem{
			Question:    fmt.Sprintf("What is \\(\\lim_{n\\to\\infty} \\frac{%dn+%d}{%dn+%d}\\)? (enter a number)", a, b, c, d),
			Answer:      strconv.Itoa(k),
			Explanation: fmt.Sprintf("Both numerator and denominator grow like n, so the limit is the ratio of the leading coefficients, %d/%d = %d.", a, c, k),
			Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "c": strconv.Itoa(c), "d": strconv.Itoa(d), "l": strconv.Itoa(k)},
		}
	case 1:
		r := generator.MultipleChoice("If the terms of a series tend to a nonzero limit, does the series diverge?", "yes", "no")
		r.Explanation = "Yes. Convergence forces the terms to tend to zero, so a nonzero limit is enough to conclude divergence by the nth-term test."
		return r
	}
	r := generator.MultipleChoice(fmt.Sprintf("Does \\(\\sum \\frac{%dn+%d}{%dn+%d}\\) converge?", a, b, c, d), "no", "yes")
	r.Explanation = fmt.Sprintf("The terms tend to %d, not zero, so the series diverges by the nth-term test regardless of how the early terms behave.", k)
	return r
}

func (g *absoluteConvergenceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	p := randRange(ctx.Difficulty, 1, 4)
	options := []string{"absolutely convergent", "conditionally convergent", "divergent"}
	switch rand.Intn(4) {
	case 0:
		// p >= 2: the absolute series is a convergent p-series.
		r := generator.MultipleChoice(fmt.Sprintf("Classify \\(\\sum (-1)^{n}/n^{%d}\\).", p+1), options[0], options[1], options[2])
		r.Explanation = fmt.Sprintf("The absolute series is the p-series with p = %d > 1, which converges, so the original series converges absolutely.", p+1)
		return r
	case 1:
		r := generator.MultipleChoice("Classify \\(\\sum (-1)^{n}/n\\).", options[1], options[0], options[2])
		r.Explanation = "The series converges by the alternating series test, but the absolute series is harmonic and diverges. Convergent without absolute convergence is conditional convergence."
		return r
	case 2:
		r := generator.MultipleChoice(fmt.Sprintf("Classify \\(\\sum (-1)^{n}n^{%d}\\).", p), options[2], options[0], options[1])
		r.Explanation = fmt.Sprintf("The terms grow in magnitude like n^%d, so they do not tend to zero and the series diverges. Neither absolute nor conditional convergence applies.", p)
		return r
	}
	r := generator.MultipleChoice("Does absolute convergence imply convergence?", "yes", "no")
	r.Explanation = "Yes. If the series of absolute values converges then the original series converges, so absolute convergence is the stronger condition of the two."
	return r
}
