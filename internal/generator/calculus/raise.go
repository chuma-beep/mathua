package calculus

// The 2026-10-19 calculus raise: thirty concepts that take the domain from 175
// to the target of 205.
//
// The survey found the domain dense in one-variable technique, so the batch
// deepens the edges of each subdomain instead of repeating it. Limits gain the
// intermediate value theorem as a computation, limit superior and inferior, and
// oscillation as the measure of discontinuity. Derivatives gain the inverse
// hyperbolic rules, the mean value theorem as a physical statement, Newton's
// iteration as a convergence question, convexity and Jensen's inequality,
// marginal analysis for economics, and the epsilon-delta definition of the
// derivative. Series gain operations, the binomial series, the Cauchy product,
// the p-series threshold, Dirichlet's test and the gamma function. Integrals
// gain the numerical error bound, the double integral in polar coordinates, a
// scalar surface integral, the beta function, consumer surplus and Fourier
// coefficients. Multivariable calculus gains Lagrange multipliers with two
// constraints, the implicit and inverse function theorems, iterated limits and
// Euler's homogeneous function theorem. Vector calculus gains the scalar line
// integral, the parametric surface area, a general flux integral and the vector
// identities.
//
// Nineteen concepts are numeric and eleven are multiple choice. Every numeric
// generator publishes the facts its answer follows from, and the tests recompute
// the answer from those facts rather than trusting the author's arithmetic.

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerRaise(reg *generator.Registry) {
	reg.Register("calc.limit.intermediate_value", &ivtValueGen{})
	reg.Register("calc.limit.limsup", &limsupGen{})
	reg.Register("calc.limit.oscillation", &oscillationGen{})

	reg.Register("calc.deriv.inverse_hyperbolic", &inverseHyperbolicGen{})
	reg.Register("calc.deriv.mvt_applications", &mvtApplicationsGen{})
	reg.Register("calc.deriv.newton_convergence", &newtonConvergenceGen{})
	reg.Register("calc.deriv.convex_jensen", &convexJensenGen{})
	reg.Register("calc.deriv.marginal_analysis", &marginalAnalysisGen{})
	reg.Register("calc.deriv.epsilon_delta_deriv", &epsilonDeltaDerivGen{})

	reg.Register("calc.series.operations", &seriesOperationsGen{})
	reg.Register("calc.series.binomial", &binomialSeriesGen{})
	reg.Register("calc.series.cauchy_product", &cauchyProductGen{})
	reg.Register("calc.series.p_series", &pSeriesGen{})
	reg.Register("calc.series.dirichlet_test", &dirichletTestGen{})
	reg.Register("calc.series.gamma", &gammaFunctionGen{})

	reg.Register("calc.integral.numerical_error", &numericalErrorGen{})
	reg.Register("calc.integral.double_polar", &doublePolarGen{})
	reg.Register("calc.integral.surface_scalar", &surfaceScalarGen{})
	reg.Register("calc.integral.beta", &betaFunctionGen{})
	reg.Register("calc.integral.consumer_surplus", &consumerSurplusGen{})
	reg.Register("calc.integral.fourier_coefficients", &fourierCoefficientsGen{})

	reg.Register("calc.multi.lagrange_two", &varietyWrap{id: "calc.multi.lagrange_two", inner: &lagrangeTwoGen{}})
	reg.Register("calc.multi.implicit_function", &implicitFunctionGen{})
	reg.Register("calc.multi.inverse_function", &inverseFunctionGen{})
	reg.Register("calc.multi.iterated_limits", &iteratedLimitsGen{})
	reg.Register("calc.multi.euler_homogeneous", &eulerHomogeneousGen{})

	reg.Register("calc.vec.line_scalar", &lineScalarGen{})
	reg.Register("calc.vec.parametric_surface", &parametricSurfaceGen{})
	reg.Register("calc.vec.flux_general", &fluxGeneralGen{})
	reg.Register("calc.vec.identities", &vectorIdentitiesGen{})
}

// calcDec formats a number with up to six decimal places and no trailing zeros.
func calcDec(v float64) string {
	s := strconv.FormatFloat(v, 'f', 6, 64)
	if strings.Contains(s, ".") {
		s = strings.TrimRight(s, "0")
		s = strings.TrimRight(s, ".")
	}
	if s == "" || s == "-0" {
		s = "0"
	}
	return s
}

func calcGCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

// calcFraction writes num/den in lowest terms, as a whole number when it divides.
func calcFraction(num, den int) string {
	if den < 0 {
		num, den = -num, -den
	}
	g := calcGCD(num, den)
	if g != 0 {
		num, den = num/g, den/g
	}
	if den == 1 {
		return strconv.Itoa(num)
	}
	return strconv.Itoa(num) + "/" + strconv.Itoa(den)
}

func calcFact(n int) int {
	out := 1
	for i := 2; i <= n; i++ {
		out *= i
	}
	return out
}

func calcBinom(n, k int) int {
	if k < 0 || k > n {
		return 0
	}
	return calcFact(n) / (calcFact(k) * calcFact(n-k))
}

type calcMCClaim struct {
	q           string
	correct     string
	distractors []string
	reason      string
}

func calcPickMC(claims []calcMCClaim) generator.Problem {
	c := claims[rand.Intn(len(claims))]
	p := generator.MultipleChoice(c.q, c.correct, c.distractors...)
	p.Explanation = c.reason
	return p
}

// ---------------------------------------------------------------------------
// Limits: the intermediate value theorem as a computation (numeric)
// ---------------------------------------------------------------------------

type ivtValueGen struct{}

func (g *ivtValueGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	root := 2 + rand.Intn(6) // 2..7
	k := root * root
	n := root + 1 + rand.Intn(2)
	answer := strconv.Itoa(root)
	return generator.Problem{
		Question: fmt.Sprintf("The function \\(f(x) = x^{2} - %d\\) is continuous, and \\(f(0) < 0 < f(%d)\\). By the intermediate value theorem it has a root between 0 and %d. What is that root? (enter a number)", k, n, n),
		Answer:   answer,
		Explanation: fmt.Sprintf("The root solves \\(x^{2} = %d\\), so \\(x = %d\\); it lies between 0 and %d, so the intermediate value theorem picks it out.",
			k, root, n),
		Facts: map[string]string{"k": strconv.Itoa(k), "n": strconv.Itoa(n), "answer": answer},
	}
}

// ---------------------------------------------------------------------------
// Limits: limit superior and inferior (multiple choice)
// ---------------------------------------------------------------------------

type limsupGen struct{}

func (g *limsupGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	claims := []calcMCClaim{
		{
			q:           "The limit superior of a bounded sequence is...",
			correct:     "the infimum over n of the supremum of the terms from n onward",
			distractors: []string{"the supremum of the whole sequence", "the largest term of the sequence", "the limit of the largest subsequence"},
			reason:      "It is the eventual upper envelope: the least number that the tail's supremum approaches as n grows.",
		},
		{
			q:           "For the sequence that alternates between -1 and 1, the limit superior and the limit inferior are...",
			correct:     "1 and -1",
			distractors: []string{"both 1", "both 0", "1 and 0"},
			reason:      "Every tail has supremum 1 and infimum -1, so the shrinking bounds are 1 and -1.",
		},
		{
			q:           "A bounded sequence converges exactly when...",
			correct:     "its limit superior equals its limit inferior",
			distractors: []string{"its limit superior is finite", "it has a convergent subsequence", "it is monotone"},
			reason:      "Equality of the two shrinking bounds is the same as the tail of the sequence closing in on one value.",
		},
	}
	return calcPickMC(append(claims, computedThinClaims("calc.limit.limsup")...))
}

// ---------------------------------------------------------------------------
// Limits: oscillation (multiple choice)
// ---------------------------------------------------------------------------

type oscillationGen struct{}

func (g *oscillationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	claims := []calcMCClaim{
		{
			q:           "The oscillation of a function on a set is...",
			correct:     "the supremum of the function on the set minus its infimum on the set",
			distractors: []string{"the largest value of the function on the set", "the number of jumps in the set", "the derivative's largest magnitude"},
			reason:      "Oscillation measures the spread of the values, so it is the difference between the extremes.",
		},
		{
			q:           "A function is continuous at a point exactly when...",
			correct:     "its oscillation on intervals shrinking to the point tends to zero",
			distractors: []string{"it is bounded near the point", "it has a one-sided limit there", "its values near the point are positive"},
			reason:      "The oscillation is small exactly when nearby values are close to the value at the point, which is continuity.",
		},
		{
			q:           "The function \\(\\sin(1/x)\\) near 0 has...",
			correct:     "oscillation 2 on every interval that contains 0",
			distractors: []string{"oscillation 0 on every such interval", "oscillation 1 on every such interval", "finite oscillation only on small intervals"},
			reason:      "The sine takes values arbitrarily close to both 1 and -1 in any neighbourhood of 0, so the spread is 2.",
		},
	}
	return calcPickMC(append(claims, computedThinClaims("calc.limit.oscillation")...))
}

// ---------------------------------------------------------------------------
// Derivatives: the inverse hyperbolic rules (multiple choice)
// ---------------------------------------------------------------------------

type inverseHyperbolicGen struct{}

func (g *inverseHyperbolicGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	claims := []calcMCClaim{
		{
			q:           "The derivative of \\(\\operatorname{arsinh}(x)\\) is...",
			correct:     "1/√(1+x²)",
			distractors: []string{"1/√(x²−1)", "1/(1−x²)", "1/(1+x²)"},
			reason:      "Inverting \\(x=\\sinh y\\) gives \\(dx/dy=\\cosh y\\) and \\(\\cosh^{2}y-\\sinh^{2}y=1\\), so \\(dy/dx=1/\\sqrt{1+x^{2}}\\).",
		},
		{
			q:           "The derivative of \\(\\operatorname{artanh}(x)\\) for \\(|x|<1\\) is...",
			correct:     "1/(1−x²)",
			distractors: []string{"1/(x²−1)", "1/√(1+x²)", "1/(1+x²)"},
			reason:      "\\(dx/dy=\\cosh^{2}y\\) and \\(\\cosh^{2}y=1/(1-\\tanh^{2}y)\\), so \\(dy/dx=1/(1-x^{2})\\).",
		},
		{
			q:           "The derivative of \\(\\operatorname{arcosh}(x)\\) for \\(x>1\\) is...",
			correct:     "1/√(x²−1)",
			distractors: []string{"1/√(1+x²)", "1/√(1−x²)", "1/(1−x²)"},
			reason:      "From \\(x=\\cosh y\\), \\(dx/dy=\\sinh y\\) with \\(\\sinh y=\\sqrt{x^{2}-1}\\) for positive y.",
		},
	}
	return calcPickMC(append(claims, computedThinClaims("calc.deriv.inverse_hyperbolic")...))
}

// ---------------------------------------------------------------------------
// Derivatives: the mean value theorem as a physical statement (numeric)
// ---------------------------------------------------------------------------

type mvtApplicationsGen struct{}

func (g *mvtApplicationsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	hours := 2 + rand.Intn(4) // 2..5
	speed := 20 + 5*rand.Intn(8)
	distance := speed * hours
	answer := strconv.Itoa(speed)
	return generator.Problem{
		Question: fmt.Sprintf("A car travels %d km in %d hours, its position being a differentiable function of time. The mean value theorem guarantees some instant at which its instantaneous speed equals its average speed. What is that speed, in km/h? (enter a number)", distance, hours),
		Answer:   answer,
		Explanation: fmt.Sprintf("The average speed is the distance over the time, %d/%d = %d km/h, and the theorem says the instantaneous speed equals it at some moment.",
			distance, hours, speed),
		Facts: map[string]string{"distance": strconv.Itoa(distance), "hours": strconv.Itoa(hours), "speed": answer},
	}
}

// ---------------------------------------------------------------------------
// Derivatives: Newton's method as a convergence question (numeric)
// ---------------------------------------------------------------------------

type newtonConvergenceGen struct{}

func (g *newtonConvergenceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	roots := []int{2, 3, 5, 7, 11}
	k := roots[rand.Intn(len(roots))]
	starts := []int{1, 2, 3}
	x0 := starts[rand.Intn(len(starts))]
	next := (float64(x0) + float64(k)/float64(x0)) / 2
	answer := calcDec(next)
	return generator.Problem{
		Question: fmt.Sprintf("Newton's method for \\(f(x)=x^{2}-%d\\) uses \\(x_{n+1} = x_n - f(x_n)/f'(x_n)\\). Starting from \\(x_0 = %d\\), what is \\(x_1\\)? (enter a decimal)", k, x0),
		Answer:   answer,
		Explanation: fmt.Sprintf("Here \\(f'(x)=2x\\), so \\(x_1 = x_0 - (x_0^{2}-%d)/(2x_0) = (x_0 + %d/x_0)/2 = %s\\) when rounded to six places.",
			k, k, answer),
		Facts: map[string]string{"k": strconv.Itoa(k), "x0": strconv.Itoa(x0), "answer": answer},
	}
}

// ---------------------------------------------------------------------------
// Derivatives: convexity and Jensen's inequality (multiple choice)
// ---------------------------------------------------------------------------

type convexJensenGen struct{}

func (g *convexJensenGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	claims := []calcMCClaim{
		{
			q:           "A twice-differentiable function on an interval is convex exactly when...",
			correct:     "its second derivative is nonnegative throughout the interval",
			distractors: []string{"its first derivative is nonnegative", "it is increasing", "its second derivative changes sign once"},
			reason:      "The bending upward that convexity describes is exactly the nonnegative second derivative.",
		},
		{
			q:           "Jensen's inequality for a convex function f and weights that sum to one says...",
			correct:     "the value at the weighted average is at most the weighted average of the values",
			distractors: []string{"the value at the weighted average is at least the weighted average of the values", "the weighted average of the values equals the value at the average", "the function is increasing"},
			reason:      "Convexity puts the graph below its chords, which is exactly Jensen's inequality.",
		},
		{
			q:           "A convex function on an open interval...",
			correct:     "is continuous there and lies below every chord between two of its points",
			distractors: []string{"must be differentiable everywhere", "must be bounded", "must be increasing"},
			reason:      "Convexity forces continuity on the open interval, but not differentiability or monotonicity; the chord condition is the definition's geometric face.",
		},
	}
	return calcPickMC(append(claims, computedThinClaims("calc.deriv.convex_jensen")...))
}

// ---------------------------------------------------------------------------
// Derivatives: marginal analysis (numeric)
// ---------------------------------------------------------------------------

type marginalAnalysisGen struct{}

func (g *marginalAnalysisGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := 1 + rand.Intn(3)
	b := 1 + rand.Intn(5)
	c := 1 + rand.Intn(9)
	x0 := 1 + rand.Intn(5)
	marginal := 2*a*x0 + b
	answer := strconv.Itoa(marginal)
	return generator.Problem{
		Question: fmt.Sprintf("A firm's cost of producing \\(x\\) units is \\(C(x) = %dx^{2} + %dx + %d\\). The marginal cost is \\(C'(x)\\). What is the marginal cost at \\(x = %d\\)? (enter a number)", a, b, c, x0),
		Answer:   answer,
		Explanation: fmt.Sprintf("Differentiating gives \\(C'(x) = %dx + %d\\), so at \\(x = %d\\) the marginal cost is %d.",
			2*a, b, x0, marginal),
		Facts: map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "x0": strconv.Itoa(x0), "marginal": answer},
	}
}

// ---------------------------------------------------------------------------
// Derivatives: the epsilon-delta definition (multiple choice)
// ---------------------------------------------------------------------------

type epsilonDeltaDerivGen struct{}

func (g *epsilonDeltaDerivGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	claims := []calcMCClaim{
		{
			q:           "The statement that \\(f'(a) = L\\) means...",
			correct:     "for every ε > 0 there is δ > 0 such that 0 < |x−a| < δ makes the difference quotient within ε of L",
			distractors: []string{"for every ε > 0 there is δ > 0 with |f(x)−f(a)| < ε whenever |x−a| < δ", "f is continuous at a", "the one-sided derivatives exist"},
			reason:      "The derivative's defining limit is the epsilon-delta condition on the difference quotient, not on the function values.",
		},
		{
			q:           "Differentiability of f at a point implies...",
			correct:     "continuity there",
			distractors: []string{"that the derivative is continuous there", "that f is increasing near the point", "that the second derivative exists there"},
			reason:      "The difference quotient tending to a finite limit forces the change in f to tend to zero, which is continuity.",
		},
		{
			q:           "Because the definition of the derivative uses a two-sided limit,...",
			correct:     "the derivative can fail to exist at a corner even though both one-sided derivatives exist",
			distractors: []string{"the derivative exists wherever both one-sided derivatives exist", "the derivative is always linear", "the function must be convex"},
			reason:      "The absolute value has one-sided slopes -1 and 1 at 0; the two-sided limit sees the disagreement and does not exist.",
		},
	}
	return calcPickMC(append(claims, computedThinClaims("calc.deriv.epsilon_delta_deriv")...))
}

// ---------------------------------------------------------------------------
// Series: operations on series (multiple choice)
// ---------------------------------------------------------------------------

type seriesOperationsGen struct{}

func (g *seriesOperationsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	claims := []calcMCClaim{
		{
			q:           "The sum of two convergent series is...",
			correct:     "a convergent series whose sum is the sum of the two sums",
			distractors: []string{"always divergent", "convergent only if both are absolutely convergent", "a series whose sum is the product of the two sums"},
			reason:      "The partial sums add term by term, so the limit of the sums is the sum of the limits.",
		},
		{
			q:           "Multiplying every term of a convergent series by a constant...",
			correct:     "leaves it convergent, with the sum multiplied by the same constant",
			distractors: []string{"may make it diverge", "leaves the sum unchanged", "converges only when the constant is below one"},
			reason:      "Each partial sum scales by the constant, and scaling commutes with the limit.",
		},
		{
			q:           "Inserting parentheses into a convergent series...",
			correct:     "does not change its sum",
			distractors: []string{"always changes the sum", "may make it diverge", "requires absolute convergence to preserve the sum"},
			reason:      "Grouping picks out a subsequence of the partial sums, and a subsequence of a convergent sequence has the same limit.",
		},
	}
	return calcPickMC(append(claims, computedThinClaims("calc.series.operations")...))
}

// ---------------------------------------------------------------------------
// Series: the binomial series (numeric)
// ---------------------------------------------------------------------------

type binomialSeriesGen struct{}

func (g *binomialSeriesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	n := 4 + rand.Intn(5) // 4..8
	k := 1 + rand.Intn(n)
	coeff := calcBinom(n, k)
	answer := strconv.Itoa(coeff)
	return generator.Problem{
		Question:    fmt.Sprintf("In the binomial expansion of \\((1+x)^{%d}\\), what is the coefficient of \\(x^{%d}\\)? (enter a number)", n, k),
		Answer:      answer,
		Explanation: fmt.Sprintf("The coefficient is the binomial coefficient C(%d, %d) = %d.", n, k, coeff),
		Facts:       map[string]string{"n": strconv.Itoa(n), "k": strconv.Itoa(k), "coeff": answer},
	}
}

// ---------------------------------------------------------------------------
// Series: the Cauchy product (numeric)
// ---------------------------------------------------------------------------

type cauchyProductGen struct{}

func (g *cauchyProductGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	n := 2 + rand.Intn(7) // 2..8
	coeff := n + 1
	answer := strconv.Itoa(coeff)
	return generator.Problem{
		Question: fmt.Sprintf("Multiply the two geometric series \\(1+x+x^{2}+\\cdots\\) and \\(1+x+x^{2}+\\cdots\\) as a Cauchy product. What is the coefficient of \\(x^{%d}\\) in the product? (enter a number)", n),
		Answer:   answer,
		Explanation: fmt.Sprintf("The coefficient of \\(x^{%d}\\) is the number of pairs of indices adding to it, 0+%d through %d+0, which is %d.",
			n, n, n, coeff),
		Facts: map[string]string{"n": strconv.Itoa(n), "coeff": answer},
	}
}

// ---------------------------------------------------------------------------
// Series: the p-series (multiple choice)
// ---------------------------------------------------------------------------

type pSeriesGen struct{}

func (g *pSeriesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := rand.Intn(4)
	switch mode {
	case 0:
		// Convergence threshold question
		return calcPickMC([]calcMCClaim{
			{
				q:           "The p-series with terms \\(1/n^{p}\\) converges exactly when...",
				correct:     "p > 1",
				distractors: []string{"p ≥ 1", "p < 1", "p > 0"},
				reason:      "The integral test compares it with the integral of \\(x^{-p}\\), which converges only for \\(p > 1\\).",
			},
			{
				q:           "Which condition on the exponent p guarantees that \\(\\sum_{n=1}^\\infty \\frac{1}{n^p}\\) diverges?",
				correct:     "p ≤ 1",
				distractors: []string{"p > 1", "p ≥ 2", "p > 0"},
				reason:      "A p-series diverges whenever the exponent p is less than or equal to 1, including the harmonic series at p = 1.",
			},
		})
	case 1:
		// Specific radical series convergence/divergence
		type radCase struct {
			term, pStr, exp string
			converges       bool
		}
		cases := []radCase{
			{"\\frac{1}{\\sqrt{n}}", "1/2", "below 1", false},
			{"\\frac{1}{\\sqrt[3]{n}}", "1/3", "below 1", false},
			{"\\frac{1}{\\sqrt[4]{n}}", "1/4", "below 1", false},
			{"\\frac{1}{\\sqrt[3]{n^2}}", "2/3", "below 1", false},
			{"\\frac{1}{\\sqrt[4]{n^3}}", "3/4", "below 1", false},
			{"\\frac{1}{n\\sqrt{n}}", "3/2", "strictly greater than 1", true},
			{"\\frac{1}{n^2\\sqrt{n}}", "5/2", "strictly greater than 1", true},
			{"\\frac{1}{n\\sqrt[3]{n}}", "4/3", "strictly greater than 1", true},
			{"\\frac{1}{\\sqrt[4]{n^5}}", "5/4", "strictly greater than 1", true},
			{"\\frac{1}{n\\sqrt[4]{n^3}}", "7/4", "strictly greater than 1", true},
		}
		c := cases[rand.Intn(len(cases))]
		q := fmt.Sprintf("The series \\(\\sum_{n=1}^\\infty %s\\) ...", c.term)
		if c.converges {
			pProb := generator.MultipleChoice(q,
				fmt.Sprintf("converges, because its exponent p = %s is strictly greater than 1", c.pStr),
				fmt.Sprintf("diverges, because its exponent p = %s is at most 1", c.pStr),
				"diverges, because its terms tend to infinity",
				"converges to 0")
			pProb.Explanation = fmt.Sprintf("The exponent is %s > 1, so the p-series converges.", c.pStr)
			return pProb
		}
		pProb := generator.MultipleChoice(q,
			fmt.Sprintf("diverges, because its exponent p = %s is at most 1", c.pStr),
			fmt.Sprintf("converges, because its exponent p = %s is strictly greater than 1", c.pStr),
			"converges, because its terms tend to 0",
			"converges to 1")
		pProb.Explanation = fmt.Sprintf("The exponent is %s <= 1, so the p-series diverges.", c.pStr)
		return pProb
	case 2:
		// Integer exponent p-series
		p := 2 + rand.Intn(7) // 2 to 8
		q := fmt.Sprintf("Does the p-series \\(\\sum_{n=1}^\\infty \\frac{1}{n^{%d}}\\) converge or diverge?", p)
		ex := fmt.Sprintf("The exponent is p = %d > 1, so by the p-series test the series converges.", p)
		pProb := generator.MultipleChoice(q,
			"converges, because p > 1",
			"diverges, because p ≤ 1",
			"diverges, because the terms tend to zero too slowly",
			"diverges, because it oscillates")
		pProb.Explanation = ex
		return pProb
	default:
		// Which of these series converges / diverges
		divergentTerms := []string{"\\sum \\frac{1}{n}", "\\sum \\frac{1}{\\sqrt{n}}", "\\sum \\frac{1}{\\sqrt[3]{n}}", "\\sum \\frac{1}{n^{0.8}}"}
		convergentTerms := []string{"\\sum \\frac{1}{n^2}", "\\sum \\frac{1}{n^3}", "\\sum \\frac{1}{n^{1.5}}", "\\sum \\frac{1}{n^{2.5}}"}
		if rand.Intn(2) == 0 {
			div := divergentTerms[rand.Intn(len(divergentTerms))]
			q := "Which of these p-series diverges?"
			pProb := generator.MultipleChoice(q, div, convergentTerms[0], convergentTerms[1], convergentTerms[2])
			pProb.Explanation = fmt.Sprintf("The series %s has exponent p <= 1, so it diverges by the p-series test.", div)
			return pProb
		}
		conv := convergentTerms[rand.Intn(len(convergentTerms))]
		q := "Which of these p-series converges?"
		pProb := generator.MultipleChoice(q, conv, divergentTerms[0], divergentTerms[1], divergentTerms[2])
		pProb.Explanation = fmt.Sprintf("The series %s has exponent p > 1, so it converges by the p-series test.", conv)
		return pProb
	}
}

// ---------------------------------------------------------------------------
// Series: Dirichlet's test (multiple choice)
// ---------------------------------------------------------------------------

type dirichletTestGen struct{}

func (g *dirichletTestGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := rand.Intn(3)
	switch mode {
	case 0:
		return calcPickMC([]calcMCClaim{
			{
				q:           "Dirichlet's test says that \\(\\sum a_n b_n\\) converges when...",
				correct:     "the partial sums of a(n) are bounded and b(n) decreases to zero",
				distractors: []string{"a(n) decreases to zero and b(n) is bounded", "both sequences are absolutely summable", "the product a(n)b(n) decreases to zero"},
				reason:      "Bounded partial sums supply the cancellation and a monotone sequence shrinking to zero supplies the size, which is enough for convergence.",
			},
			{
				q:           "Abel's test is the version of Dirichlet's test in which...",
				correct:     "the monotone factor decreases to a positive limit rather than to zero",
				distractors: []string{"the cancelling factor is absolutely summable", "the two factors are interchanged", "the series has positive terms"},
				reason:      "Abel's test handles a factor converging to a nonzero limit by subtracting that limit and applying Dirichlet to the remainder.",
			},
		})
	case 1:
		k := 1 + rand.Intn(5)
		q := fmt.Sprintf("Applied to \\(\\sum_{n=1}^\\infty \\frac{\\sin(%dn)}{n}\\), Dirichlet's test shows that...", k)
		correct := fmt.Sprintf("the series converges, because the partial sums of sin(%dn) are bounded and 1/n decreases to zero", k)
		d1 := fmt.Sprintf("the series diverges, because sin(%dn) does not tend to zero", k)
		d2 := "the series converges absolutely"
		d3 := "the test cannot be applied"
		ex := fmt.Sprintf("The partial sums of sin(%dn) remain bounded by 1/sin(%d/2) and 1/n decreases monotonically to 0, satisfying Dirichlet's test.", k, k)
		p := generator.MultipleChoice(q, correct, d1, d2, d3)
		p.Explanation = ex
		return p
	default:
		k := 1 + rand.Intn(5)
		q := fmt.Sprintf("Applied to \\(\\sum_{n=1}^\\infty \\frac{\\cos(%dn)}{\\sqrt{n}}\\), Dirichlet's test shows that...", k)
		correct := fmt.Sprintf("the series converges, because the partial sums of cos(%dn) are bounded and 1/sqrt(n) decreases to zero", k)
		d1 := fmt.Sprintf("the series diverges, because cos(%dn) oscillates indefinitely", k)
		d2 := "the series converges absolutely"
		d3 := "the test is inconclusive"
		ex := fmt.Sprintf("The partial sums of cos(%dn) are bounded and 1/sqrt(n) decreases monotonically to 0, so the series converges by Dirichlet's test.", k)
		p := generator.MultipleChoice(q, correct, d1, d2, d3)
		p.Explanation = ex
		return p
	}
}

// ---------------------------------------------------------------------------
// Series: the gamma function (numeric)
// ---------------------------------------------------------------------------

type gammaFunctionGen struct{}

func (g *gammaFunctionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	n := randRange(ctx.Difficulty, 3, 13)
	switch rand.Intn(3) {
	case 0:
		// The recurrence is the idea; the factorial is the computation.
		return generator.Problem{
			Question:    fmt.Sprintf("Use \\(\\Gamma(n+1)=n\\Gamma(n)\\). What is \\(\\Gamma(%d)/\\Gamma(%d)\\)?", n+1, n),
			Answer:      strconv.Itoa(n),
			Explanation: fmt.Sprintf("\\(\\Gamma(%d)=%d\\,\\Gamma(%d)\\), so the ratio is %d.", n+1, n, n, n),
			Facts:       map[string]string{"n": strconv.Itoa(n), "shape": "ratio1"},
		}
	case 1:
		return generator.Problem{
			Question:    fmt.Sprintf("Use \\(\\Gamma(n+1)=n\\Gamma(n)\\) twice. What is \\(\\Gamma(%d)/\\Gamma(%d)\\)?", n+2, n),
			Answer:      strconv.Itoa(n * (n + 1)),
			Explanation: fmt.Sprintf("\\(\\Gamma(%d)=(%d)(%d)\\Gamma(%d)\\), so the ratio is %d × %d = %d.", n+2, n+1, n, n, n+1, n, n*(n+1)),
			Facts:       map[string]string{"n": strconv.Itoa(n), "shape": "ratio2"},
		}
	}
	value := calcFact(n - 1)
	answer := strconv.Itoa(value)
	return generator.Problem{
		Question:    fmt.Sprintf("The gamma function satisfies \\(\\Gamma(n) = (n-1)!\\) for positive integers. What is \\(\\Gamma(%d)\\)?", n),
		Answer:      answer,
		Explanation: fmt.Sprintf("\\(\\Gamma(%d) = %d! = %d\\).", n, n-1, value),
		Facts:       map[string]string{"n": strconv.Itoa(n), "value": answer, "shape": "factorial"},
	}
}

// ---------------------------------------------------------------------------
// Integrals: the numerical error bound (numeric)
// ---------------------------------------------------------------------------

type numericalErrorGen struct{}

func (g *numericalErrorGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	m := 12 * (1 + rand.Intn(3)) // 12, 24, 36
	width := 1 + rand.Intn(3)    // 1..3
	n := 2 * (1 + rand.Intn(2))  // 2 or 4
	bound := float64(m*width*width*width) / (12 * float64(n*n))
	answer := calcDec(bound)
	return generator.Problem{
		Question:    fmt.Sprintf("The trapezoid rule on n subintervals has error at most \\(M(b-a)^{3}/(12n^{2})\\), where \\(M\\) bounds \\(|f''|\\). For \\(M = %d\\), \\(b-a = %d\\) and \\(n = %d\\), what is that bound? (enter a decimal)", m, width, n),
		Answer:      answer,
		Explanation: fmt.Sprintf("Substituting gives \\(%d \\cdot %d^{3}/(12 \\cdot %d^{2}) = %s\\).", m, width, n, answer),
		Facts:       map[string]string{"m": strconv.Itoa(m), "width": strconv.Itoa(width), "n": strconv.Itoa(n), "bound": answer},
	}
}

// ---------------------------------------------------------------------------
// Integrals: the double integral in polar coordinates (numeric)
// ---------------------------------------------------------------------------

type doublePolarGen struct{}

func (g *doublePolarGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// The radius alone gave six situations, so the region now varies: a disk, an
	// annulus and a semicircle all have areas that are integer multiples of pi,
	// which keeps the answer typeable as a number.
	r := 2 + rand.Intn(5)
	shape := "disk"
	var multiple int
	switch rand.Intn(3) {
	case 1:
		shape = "annulus"
		multiple = r*r - 1
	case 2:
		shape = "semicircle"
		multiple = r * r
	default:
		multiple = r * r
	}
	answer := strconv.Itoa(multiple)
	var question, explanation string
	switch shape {
	case "annulus":
		question = fmt.Sprintf("Let \\(D\\) be the annulus \\(1 \\le r \\le %d\\) centred at the origin. The area of \\(D\\) is a multiple of \\(\\pi\\). What is that multiple, that is, \\(\\iint_D dA\\) divided by \\(\\pi\\)? (enter a number)", r)
		explanation = fmt.Sprintf("In polar coordinates the area element is \\(r\\,dr\\,d\\theta\\), so the area is \\(\\int_0^{2\\pi}\\int_1^{%d} r\\,dr\\,d\\theta = \\pi \\cdot (%d - 1)\\), and the multiple is %d.", r, r*r, multiple)
	case "semicircle":
		if r%2 == 0 {
			answer = strconv.Itoa(multiple / 2)
		} else {
			answer = fmt.Sprintf("%d/2", multiple)
		}
		question = fmt.Sprintf("Let \\(D\\) be the upper half of the disk of radius %d centred at the origin, so \\(0 \\le \\theta \\le \\pi\\). The area of \\(D\\) is a multiple of \\(\\pi\\). What is that multiple? (enter a number or fraction)", r)
		explanation = fmt.Sprintf("The area element is \\(r\\,dr\\,d\\theta\\), and the angle now runs over half a turn, so the area is \\(\\pi \\cdot %d^2 / 2 = \\pi \\cdot %s\\), giving %s.", r, answer, answer)
	default:
		question = fmt.Sprintf("Let \\(D\\) be the disk of radius %d centred at the origin. The area of \\(D\\) is a multiple of \\(\\pi\\). What is that multiple, that is, \\(\\iint_D dA\\) divided by \\(\\pi\\)? (enter a number)", r)
		explanation = fmt.Sprintf("In polar coordinates the area element is \\(r\\,dr\\,d\\theta\\), so the area is \\(\\int_0^{2\\pi}\\int_0^{%d} r\\,dr\\,d\\theta = \\pi \\cdot %d\\), and the multiple is %d.", r, r*r, multiple)
	}
	return generator.Problem{
		Question:    question,
		Answer:      answer,
		Explanation: explanation,
		Facts:       map[string]string{"r": strconv.Itoa(r), "shape": shape, "answer": answer},
	}
}

// ---------------------------------------------------------------------------
// Integrals: a scalar surface integral (numeric)
// ---------------------------------------------------------------------------

type surfaceScalarGen struct{}

func (g *surfaceScalarGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	r := 1 + rand.Intn(5)
	// Four shapes rather than one, because the radius alone gave five
	// situations: the question never changed, only the number in it.
	switch rand.Intn(4) {
	case 0:
		return generator.Problem{
			Question:    fmt.Sprintf("Let \\(S\\) be the upper hemisphere of radius %d centred at the origin, with \\(z \\ge 0\\). Evaluate \\(\\iint_S z\\,dS\\) and divide by \\(\\pi\\). What is that number? (enter a number)", r),
			Answer:      strconv.Itoa(r * r * r),
			Explanation: fmt.Sprintf("On the hemisphere \\(z = %d\\cos\\varphi\\) and \\(dS = %d^{2}\\sin\\varphi\\,d\\varphi\\,d\\theta\\), so the integral is \\(2\\pi \\cdot %d^{3} \\cdot \\tfrac{1}{2} = \\pi \\cdot %d\\) and the multiple is %d.", r, r, r, r*r*r, r*r*r),
			Facts:       map[string]string{"shape": "hemisphere", "r": strconv.Itoa(r), "answer": strconv.Itoa(r * r * r)},
		}
	case 1:
		return generator.Problem{
			Question:    fmt.Sprintf("A sphere has radius %d. What is its surface area divided by \\(\\pi\\)? (enter a number)", r),
			Answer:      strconv.Itoa(4 * r * r),
			Explanation: fmt.Sprintf("The surface area of a sphere is \\(4\\pi r^{2}\\), so with \\(r = %d\\) the area divided by \\(\\pi\\) is \\(4 \\cdot %d = %d\\).", r, r*r, 4*r*r),
			Facts:       map[string]string{"shape": "sphere", "r": strconv.Itoa(r), "answer": strconv.Itoa(4 * r * r)},
		}
	case 2:
		return generator.Problem{
			Question:    fmt.Sprintf("A disk has radius %d. What is its area divided by \\(\\pi\\)? (enter a number)", r),
			Answer:      strconv.Itoa(r * r),
			Explanation: fmt.Sprintf("The area of a disk is \\(\\pi r^{2}\\), so with \\(r = %d\\) the area divided by \\(\\pi\\) is %d.", r, r*r),
			Facts:       map[string]string{"shape": "disk", "r": strconv.Itoa(r), "answer": strconv.Itoa(r * r)},
		}
	default:
		return generator.Problem{
			Question:    fmt.Sprintf("A right circular cylinder has radius %d and height %d. What is its lateral surface area divided by \\(\\pi\\)? (enter a number)", r, r),
			Answer:      strconv.Itoa(2 * r * r),
			Explanation: fmt.Sprintf("The lateral surface area is \\(2\\pi r h\\), and with \\(r = h = %d\\) that is \\(2\\pi \\cdot %d\\), so the multiple of \\(\\pi\\) is %d.", r, r*r, 2*r*r),
			Facts:       map[string]string{"shape": "cylinder", "r": strconv.Itoa(r), "answer": strconv.Itoa(2 * r * r)},
		}
	}
}

// ---------------------------------------------------------------------------
// Integrals: the beta function (numeric)
// ---------------------------------------------------------------------------

type betaFunctionGen struct{}

func (g *betaFunctionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	m := 2 + rand.Intn(3) // 2..4
	n := 2 + rand.Intn(3) // 2..4
	num := calcFact(m-1) * calcFact(n-1)
	den := calcFact(m + n - 1)
	answer := calcFraction(num, den)
	return generator.Problem{
		Question: fmt.Sprintf("The beta function is \\(B(m,n) = (m-1)!(n-1)!/(m+n-1)!\\). What is \\(B(%d,%d)\\)? (enter a fraction or a whole number)", m, n),
		Answer:   answer,
		Explanation: fmt.Sprintf("\\(B(%d,%d) = %d!\\cdot%d!/(%d-1)! = %d/%d = %s\\) in lowest terms.",
			m, n, m-1, n-1, m+n, num, den, answer),
		Facts: map[string]string{"m": strconv.Itoa(m), "n": strconv.Itoa(n), "answer": answer},
	}
}

// ---------------------------------------------------------------------------
// Integrals: consumer surplus (numeric)
// ---------------------------------------------------------------------------

type consumerSurplusGen struct{}

func (g *consumerSurplusGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	qStar := 2 * (1 + rand.Intn(3)) // 2, 4, 6
	p := 1 + rand.Intn(4)
	a := qStar + p
	surplus := qStar * qStar / 2
	answer := strconv.Itoa(surplus)
	return generator.Problem{
		Question: fmt.Sprintf("Demand is \\(D(q) = %d - q\\) and the market price is \\(p = %d\\), so the quantity sold is where \\(D(q) = p\\). Consumer surplus is the area between the demand curve and the price line. What is it? (enter a number)", a, p),
		Answer:   answer,
		Explanation: fmt.Sprintf("The quantity is \\(q^{*} = %d\\), and consumer surplus is \\(\\int_0^{%d}(%d-q)\\,dq - %d\\cdot%d = %d\\).",
			qStar, qStar, a, p, qStar, surplus),
		Facts: map[string]string{"a": strconv.Itoa(a), "p": strconv.Itoa(p), "qstar": strconv.Itoa(qStar), "surplus": answer},
	}
}

// ---------------------------------------------------------------------------
// Integrals: Fourier coefficients (numeric)
// ---------------------------------------------------------------------------

type fourierCoefficientsGen struct{}

func (g *fourierCoefficientsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	n := randRange(ctx.Difficulty, 1, 2+int(ctx.Difficulty*8)) // 1..10 at full difficulty
	sign := 1.0
	if n%2 == 0 {
		sign = -1.0
	}
	value := sign * 2.0 / float64(n)
	answer := calcDec(value)
	return generator.Problem{
		Question:    fmt.Sprintf("For \\(f(x) = x\\) on \\([-\\pi,\\pi]\\), the sine coefficient is \\(b_n = 2(-1)^{n+1}/n\\). What is \\(b_{%d}\\)? (enter a decimal)", n),
		Answer:      answer,
		Explanation: fmt.Sprintf("With \\(n = %d\\), \\(b_{%d} = 2(-1)^{%d}/%d = %s\\) to six places.", n, n, n+1, n, answer),
		Facts:       map[string]string{"n": strconv.Itoa(n), "b": answer},
	}
}

// ---------------------------------------------------------------------------
// Multivariable calculus: Lagrange multipliers with two constraints (multiple choice)
// ---------------------------------------------------------------------------

type lagrangeTwoGen struct{}

func (g *lagrangeTwoGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return calcPickMC([]calcMCClaim{
		{
			q:           "To find the extrema of \\(f\\) subject to the two constraints \\(g = 0\\) and \\(h = 0\\), the method of Lagrange multipliers solves...",
			correct:     "grad f = λ grad g + μ grad h together with the two constraint equations",
			distractors: []string{"grad f = λ grad g together with h = 0", "grad f = λ grad g × grad h", "grad f = 0 together with the two constraints"},
			reason:      "At a constrained extremum the gradient of f lies in the plane spanned by the two constraint gradients, with one multiplier per constraint.",
		},
		{
			q:           "In the two-constraint Lagrange system for three variables, the number of unknowns is...",
			correct:     "five: the three coordinates and the two multipliers",
			distractors: []string{"three: the coordinates only", "four: the coordinates and one multiplier", "six: the coordinates and three multipliers"},
			reason:      "One multiplier is introduced per constraint, so three coordinates plus two multipliers match the five equations.",
		},
		{
			q:           "Saying \\(\\nabla f\\) lies in the span of \\(\\nabla g\\) and \\(\\nabla h\\) at a constrained extremum means...",
			correct:     "f cannot be increased by any small step that keeps both constraints satisfied",
			distractors: []string{"f is constant along the constraint set", "grad f = 0", "the constraints are linearly dependent"},
			reason:      "If the gradient had a component outside that span, a small admissible step along the constraint intersection would change f, contradicting the extremum.",
		},
		{
			q:           "For two constraints \\(g(x,y,z) = 0\\) and \\(h(x,y,z) = 0\\), the regularity condition required by the Lagrange multiplier method is that...",
			correct:     "grad g and grad h are linearly independent (grad g × grad h ≠ 0)",
			distractors: []string{"grad g and grad h are orthogonal", "grad g and grad h are parallel", "grad g = 0 and grad h = 0"},
			reason:      "Linear independence of the constraint gradients ensures that the tangent spaces intersect transversally and define a smooth 1D curve.",
		},
		{
			q:           "Geometrically, when optimizing \\(f(x,y,z)\\) subject to \\(g(x,y,z) = 0\\) and \\(h(x,y,z) = 0\\), the gradient \\(\\nabla f\\) at an extremum must be...",
			correct:     "orthogonal to the tangent line of the intersection curve of both constraint surfaces",
			distractors: []string{"parallel to the tangent line of the intersection curve", "orthogonal to both surfaces simultaneously", "equal to the zero vector"},
			reason:      "Any directional derivative along the intersection curve must vanish, so grad f is normal to the curve's tangent direction.",
		},
	})
}

// ---------------------------------------------------------------------------
// Multivariable calculus: the implicit function theorem (numeric)
// ---------------------------------------------------------------------------

type implicitFunctionGen struct{}

func (g *implicitFunctionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	x := 1 + 2*rand.Intn(3) // 1, 3, 5
	y := 2 + 2*rand.Intn(4) // 2, 4, 6, 8
	k := x*x + y*y
	slope := calcDec(-float64(x) / float64(y))
	return generator.Problem{
		Question: fmt.Sprintf("The equation \\(x^{2} + y^{2} = %d\\) defines \\(y\\) as a function of \\(x\\) near the point \\((x,y) = (%d,%d)\\). The implicit function theorem gives \\(dy/dx = -F_x/F_y\\). What is \\(dy/dx\\) there? (enter a decimal)", k, x, y),
		Answer:   slope,
		Explanation: fmt.Sprintf("With \\(F = x^{2}+y^{2}-%d\\), \\(F_x = 2x = %d\\) and \\(F_y = 2y = %d\\), so \\(dy/dx = -%d/%d = %s\\).",
			k, 2*x, 2*y, 2*x, 2*y, slope),
		Facts: map[string]string{"x": strconv.Itoa(x), "y": strconv.Itoa(y), "k": strconv.Itoa(k), "slope": slope},
	}
}

// ---------------------------------------------------------------------------
// Multivariable calculus: the inverse function theorem (numeric)
// ---------------------------------------------------------------------------

type inverseFunctionGen struct{}

func (g *inverseFunctionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	x := 1 + rand.Intn(3)
	y := 1 + rand.Intn(3)
	det := 4 * (x*x + y*y)
	answer := strconv.Itoa(det)
	return generator.Problem{
		Question:    fmt.Sprintf("The map \\((u,v) = (x^{2}-y^{2},\\, 2xy)\\) has Jacobian determinant \\(4(x^{2}+y^{2})\\). The inverse function theorem says the map is locally invertible where that determinant is nonzero. What is the determinant at \\((x,y) = (%d,%d)\\)? (enter a number)", x, y),
		Answer:      answer,
		Explanation: fmt.Sprintf("\\(4(%d^{2}+%d^{2}) = %d\\), which is positive, so the map is locally invertible there.", x, y, det),
		Facts:       map[string]string{"x": strconv.Itoa(x), "y": strconv.Itoa(y), "det": answer},
	}
}

// ---------------------------------------------------------------------------
// Multivariable calculus: iterated limits (multiple choice)
// ---------------------------------------------------------------------------

type iteratedLimitsGen struct{}

func (g *iteratedLimitsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	claims := []calcMCClaim{
		{
			q:           "For a function of two variables, the two iterated limits...",
			correct:     "can be different, so the order in which they are taken matters",
			distractors: []string{"are always equal", "always equal the double limit", "exist only when the function is continuous"},
			reason:      "Each iterated limit holds one variable fixed while the other moves, and the two orders need not agree.",
		},
		{
			q:           "If the double limit of f at a point exists and one iterated limit also exists, then...",
			correct:     "the iterated limit equals the double limit",
			distractors: []string{"the other iterated limit need not exist", "the function is continuous there", "the iterated limits are different"},
			reason:      "The double limit controls every path to the point, and an iterated limit is the limit along one special family of paths.",
		},
		{
			q:           "The function \\(f(x,y) = xy/(x^{2}+y^{2})\\) away from the origin shows that...",
			correct:     "both iterated limits can exist and be equal while the double limit does not exist",
			distractors: []string{"the double limit exists whenever the iterated limits agree", "the function is continuous at the origin", "the iterated limits are unequal"},
			reason:      "Along either axis the quotient is 0, so both iterated limits are 0, but along \\(y = x\\) it is \\(1/2\\), so the double limit does not exist.",
		},
	}
	return calcPickMC(append(claims, computedMultiClaims("calc.multi.iterated_limits")...))
}

// ---------------------------------------------------------------------------
// Multivariable calculus: Euler's homogeneous function theorem (numeric)
// ---------------------------------------------------------------------------

type eulerHomogeneousGen struct{}

func (g *eulerHomogeneousGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	x := 1 + rand.Intn(3)
	y := 1 + rand.Intn(3)
	f := x*x*x + 3*x*y*y
	value := 3 * f
	answer := strconv.Itoa(value)
	return generator.Problem{
		Question:    fmt.Sprintf("The function \\(f(x,y) = x^{3} + 3xy^{2}\\) is homogeneous of degree 3, so Euler's theorem gives \\(x f_x + y f_y = 3f\\). What is \\(x f_x + y f_y\\) at \\((x,y) = (%d,%d)\\)? (enter a number)", x, y),
		Answer:      answer,
		Explanation: fmt.Sprintf("\\(f(%d,%d) = %d\\) and the sum is three times it, %d.", x, y, f, value),
		Facts:       map[string]string{"x": strconv.Itoa(x), "y": strconv.Itoa(y), "f": strconv.Itoa(f), "value": answer},
	}
}

// ---------------------------------------------------------------------------
// Vector calculus: the scalar line integral (numeric)
// ---------------------------------------------------------------------------

type lineScalarGen struct{}

func (g *lineScalarGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// Primitive Pythagorean triples from m > n, scaled, so the segment length is a
	// whole number and the integral stays exact.
	m := randRange(ctx.Difficulty, 2, 4+int(ctx.Difficulty*4))
	n := randRange(ctx.Difficulty, 1, m-1)
	a := m*m - n*n
	b := 2 * m * n
	length := m*m + n*n
	k := randRange(ctx.Difficulty, 1, 4)
	a, b, length = k*a, k*b, k*length
	value := float64(a) * float64(length) / 2
	answer := calcDec(value)
	return generator.Problem{
		Question: fmt.Sprintf("Let \\(C\\) be the straight segment from \\((0,0)\\) to \\((%d,%d)\\). Evaluate the scalar line integral \\(\\int_C x\\,ds\\). (enter a decimal)", a, b),
		Answer:   answer,
		Explanation: fmt.Sprintf("The segment has length %d, and x averages %s along it, so the integral is %d times %s = %s.",
			length, decimal(float64(a)/2, 2), length, decimal(float64(a)/2, 2), answer),
		Facts: map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "length": strconv.Itoa(length), "answer": answer},
	}
}

// ---------------------------------------------------------------------------
// Vector calculus: the parametric surface area (numeric)
// ---------------------------------------------------------------------------

type parametricSurfaceGen struct{}

func (g *parametricSurfaceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type leg struct{ q, r, length int }
	legs := []leg{{3, 4, 5}, {6, 8, 10}, {5, 12, 13}}
	l := legs[rand.Intn(len(legs))]
	p := 1 + rand.Intn(4)
	area := p * l.length
	answer := strconv.Itoa(area)
	return generator.Problem{
		Question: fmt.Sprintf("A surface is parametrised by \\(\\mathbf{r}(u,v) = u\\,(%d,0,0) + v\\,(0,%d,%d)\\) over the unit square. Its area is \\(|\\mathbf{r}_u \\times \\mathbf{r}_v|\\). What is the area? (enter a number)", p, l.q, l.r),
		Answer:   answer,
		Explanation: fmt.Sprintf("The cross product of \\((%d,0,0)\\) and \\((0,%d,%d)\\) has length \\(%d\\sqrt{%d^{2}+%d^{2}} = %d \\cdot %d = %d\\).",
			p, l.q, l.r, p, l.q, l.r, p, l.length, area),
		Facts: map[string]string{"p": strconv.Itoa(p), "q": strconv.Itoa(l.q), "r": strconv.Itoa(l.r), "area": answer},
	}
}

// ---------------------------------------------------------------------------
// Vector calculus: a general flux integral (numeric)
// ---------------------------------------------------------------------------

type fluxGeneralGen struct{}

func (g *fluxGeneralGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	s := 1 + rand.Intn(3)
	t := 1 + rand.Intn(4)
	fz := 1 + rand.Intn(5)
	flux := fz * s * t
	answer := strconv.Itoa(flux)
	return generator.Problem{
		Question: fmt.Sprintf("The constant field \\(\\mathbf{F} = (2,3,%d)\\) passes through the rectangle \\(0 \\le x \\le %d\\), \\(0 \\le y \\le %d\\) in the plane \\(z = 0\\), oriented upward. What is the flux \\(\\iint \\mathbf{F}\\cdot\\mathbf{n}\\,dS\\)? (enter a number)", fz, s, t),
		Answer:   answer,
		Explanation: fmt.Sprintf("The upward normal is \\((0,0,1)\\), so only the z-component contributes, and it equals %d times the area \\(%d \\cdot %d = %d\\), giving %d.",
			fz, s, t, s*t, flux),
		Facts: map[string]string{"s": strconv.Itoa(s), "t": strconv.Itoa(t), "fz": strconv.Itoa(fz), "flux": answer},
	}
}

// ---------------------------------------------------------------------------
// Vector calculus: the vector identities (multiple choice)
// ---------------------------------------------------------------------------

type vectorIdentitiesGen struct{}

func (g *vectorIdentitiesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	claims := []calcMCClaim{
		{
			q:           "For any smooth vector field, the divergence of the curl is...",
			correct:     "zero",
			distractors: []string{"the Laplacian of the field", "equal to the curl of the divergence", "zero only for conservative fields"},
			reason:      "The divergence of the curl expands into mixed second partials that cancel pairwise.",
		},
		{
			q:           "The curl of a gradient is...",
			correct:     "zero, so every gradient field is irrotational",
			distractors: []string{"the divergence of the field", "zero only when the field is linear", "equal to the Laplacian"},
			reason:      "Mixed partial derivatives commute, so the antisymmetric combination in the curl of a gradient vanishes.",
		},
		{
			q:           "The two identities, divergence of a curl is zero and curl of a gradient is zero, together mean...",
			correct:     "the two composite operators annihilate every smooth field because mixed partials cancel",
			distractors: []string{"every vector field is a gradient", "every vector field is a curl", "divergence and curl are inverses"},
			reason:      "Both identities come from the equality of mixed partial derivatives, so they hold for every smooth field of the right kind.",
		},
	}
	return calcPickMC(append(claims, computedMultiClaims("calc.vec.identities")...))
}
