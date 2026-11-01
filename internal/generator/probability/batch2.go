package probability

// The second probability-theory batch: twenty-five concepts that close the
// domain with joint and conditional distributions, covariance and correlation,
// the modes of convergence and their limit theorems, characteristic functions and
// stochastic processes.
//
// The first batch ended at the classical limit inequalities and the central limit
// theorem. This batch starts where those leave off: it conditions one variable on
// another, measures linear dependence, separates the modes of convergence, and
// introduces the transform that proves the limit theorems and the processes that
// run in time.
//
// Twelve concepts are multiple choice and thirteen are numeric. Every numeric
// generator publishes the parameters its answer follows from, so the tests
// recompute the answer rather than trusting the author's arithmetic.

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/mathutil"
)

// ---------------------------------------------------------------------------
// Marginal distributions (numeric)
// ---------------------------------------------------------------------------

type marginalGen struct{}

func (g *marginalGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	tier := mathutil.DifficultyTier(ctx.Difficulty)
	var totalOptions []int
	switch tier {
	case 0:
		totalOptions = []int{10, 20}
	case 1:
		totalOptions = []int{20, 25, 40, 50}
	default:
		totalOptions = []int{50, 80, 100}
	}
	total := totalOptions[rand.Intn(len(totalOptions))]

	wX := 2 + rand.Intn(total-3)
	if wX >= total-1 {
		wX = total - 2
	}
	w10 := 1 + rand.Intn(wX-1)
	w11 := wX - w10
	rem := total - wX
	w00 := 1 + rand.Intn(rem-1)
	w01 := rem - w00

	weights := []int{w00, w01, w10, w11}
	answer := decimal6(float64(w10+w11) / float64(total))
	parts := make([]string, len(weights))
	for i, w := range weights {
		parts[i] = strconv.Itoa(w)
	}

	p := generator.Problem{
		Question: "A joint mass function on the four pairs (0,0), (0,1), (1,0) and (1,1) is proportional to the weights " +
			strings.Join(parts, ", ") + " in that order. What is the marginal P(X = 1)? (enter a decimal)",
		Answer: answer,
		Facts: map[string]string{
			"weights": strings.Join(parts, ","),
			"total":   strconv.Itoa(total),
			"answer":  answer,
		},
	}
	p.Explanation = fmt.Sprintf("Subgoal 1: State governing formula. The marginal probability P(X = 1) is found by summing joint probabilities across all values of Y: P(X = 1) = P(X = 1, Y = 0) + P(X = 1, Y = 1) = (w10 + w11) / T.\nSubgoal 2: Substitute parameters. The weights corresponding to X = 1 are w10 = %d and w11 = %d out of total weight T = %d.\nSubgoal 3: Evaluate and simplify. P(X = 1) = (%d + %d) / %d = %d/%d = %s.",
		w10, w11, total, w10, w11, total, w10+w11, total, answer)
	return p
}

// ---------------------------------------------------------------------------
// Joint densities (numeric)
// ---------------------------------------------------------------------------

type jointDensityGen struct{}

func (g *jointDensityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	tier := mathutil.DifficultyTier(ctx.Difficulty)
	var dims []int
	switch tier {
	case 0:
		dims = []int{2, 4, 5, 10}
	case 1:
		dims = []int{4, 5, 8, 10, 20}
	default:
		dims = []int{5, 8, 10, 16, 20, 25, 40, 50}
	}
	a := dims[rand.Intn(len(dims))]
	b := dims[rand.Intn(len(dims))]
	x := 1 + rand.Intn(a)
	y := 1 + rand.Intn(b)

	answer := decimal6(float64(x) / float64(a) * float64(y) / float64(b))
	p := generator.Problem{
		Question: fmt.Sprintf("A pair (X, Y) is uniform on the rectangle from 0 to %d in x and 0 to %d in y. What is P(X ≤ %d and Y ≤ %d)? (enter a decimal)",
			a, b, x, y),
		Answer: answer,
		Facts: map[string]string{
			"a":      strconv.Itoa(a),
			"b":      strconv.Itoa(b),
			"x":      strconv.Itoa(x),
			"y":      strconv.Itoa(y),
			"answer": answer,
		},
	}
	p.Explanation = fmt.Sprintf("Subgoal 1: State governing formula. For a uniform distribution on the rectangle [0, a] × [0, b], coordinates X and Y are independent with uniform marginals, so P(X ≤ x and Y ≤ y) = (x/a) · (y/b).\nSubgoal 2: Substitute parameters. Substitute bounds a = %d and b = %d with evaluated thresholds x = %d and y = %d.\nSubgoal 3: Evaluate and simplify. P(X ≤ %d and Y ≤ %d) = (%d/%d) · (%d/%d) = %d/%d = %s.",
		a, b, x, y, x, y, x, a, y, b, x*y, a*b, answer)
	return p
}

// ---------------------------------------------------------------------------
// Independent random variables (multiple choice)
// ---------------------------------------------------------------------------

type rvIndependenceGen struct{}

func (g *rvIndependenceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return rvIndependenceProblem()
}

// ---------------------------------------------------------------------------
// Conditional expectation (numeric)
// ---------------------------------------------------------------------------

type conditionalGen struct{}

func (g *conditionalGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	tier := mathutil.DifficultyTier(ctx.Difficulty)
	y := rand.Intn(2)
	var colSums []int
	switch tier {
	case 0:
		colSums = []int{2, 4, 5, 10}
	case 1:
		colSums = []int{4, 5, 8, 10, 20}
	default:
		colSums = []int{10, 20, 25, 40, 50}
	}
	targetSum := colSums[rand.Intn(len(colSums))]
	otherSum := colSums[rand.Intn(len(colSums))]

	wChosen1 := 1 + rand.Intn(targetSum-1)
	wChosen0 := targetSum - wChosen1

	wOther1 := 1 + rand.Intn(otherSum-1)
	wOther0 := otherSum - wOther1

	// weights order: (0,0), (0,1), (1,0), (1,1)
	var weights []int
	if y == 0 {
		weights = []int{wChosen0, wOther0, wChosen1, wOther1}
	} else {
		weights = []int{wOther0, wChosen0, wOther1, wChosen1}
	}

	num := weights[2+y]
	den := weights[y] + num
	answer := decimal6(float64(num) / float64(den))

	parts := make([]string, len(weights))
	for i, w := range weights {
		parts[i] = strconv.Itoa(w)
	}

	p := generator.Problem{
		Question: "A joint mass function on the four pairs (0,0), (0,1), (1,0) and (1,1) is proportional to the weights " +
			strings.Join(parts, ", ") + " in that order. What is E[X | Y = " + strconv.Itoa(y) + "]? (enter a decimal)",
		Answer: answer,
		Facts: map[string]string{
			"weights": strings.Join(parts, ","),
			"y":       strconv.Itoa(y),
			"answer":  answer,
		},
	}
	p.Explanation = fmt.Sprintf("Subgoal 1: State governing formula. Conditioning on Y = %d fixes the column Y = %d. For a binary indicator variable X ∈ {0, 1}, the conditional expectation is E[X | Y = %d] = P(X = 1 | Y = %d) = w_{1%d} / (w_{0%d} + w_{1%d}).\nSubgoal 2: Substitute parameters. For column Y = %d, the weights are w_{0%d} = %d and w_{1%d} = %d.\nSubgoal 3: Evaluate and simplify. E[X | Y = %d] = %d / (%d + %d) = %d/%d = %s.",
		y, y, y, y, y, y, y, y, y, weights[y], y, num, y, num, weights[y], num, num, den, answer)
	return p
}

// ---------------------------------------------------------------------------
// The law of total expectation (numeric)
// ---------------------------------------------------------------------------

type totalExpectationGen struct{}

func (g *totalExpectationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	tier := mathutil.DifficultyTier(ctx.Difficulty)
	var denomOptions []int
	switch tier {
	case 0:
		denomOptions = []int{2, 4, 5, 10}
	case 1:
		denomOptions = []int{4, 5, 8, 10, 20}
	default:
		denomOptions = []int{5, 8, 10, 16, 20, 25}
	}
	pd := denomOptions[rand.Intn(len(denomOptions))]
	pn := 1 + rand.Intn(pd-1)
	for mathutil.GCD(pn, pd) > 1 {
		pn = 1 + rand.Intn(pd-1)
	}

	maxMean := mathutil.ScaleInt(ctx.Difficulty, 6, 25)
	ca := rand.Intn(maxMean + 1)
	cb := rand.Intn(maxMean + 1)
	for ca == cb {
		cb = rand.Intn(maxMean + 1)
	}

	p := float64(pn) / float64(pd)
	answer := decimal6((1-p)*float64(ca) + p*float64(cb))

	pb := generator.Problem{
		Question: fmt.Sprintf("A variable Y takes the values zero and one with P(Y = 1) = %d/%d. Given Y, the conditional mean of X is %d when Y is zero and %d when Y is one. What is E[X]? (enter a decimal)",
			pn, pd, ca, cb),
		Answer: answer,
		Facts: map[string]string{
			"ca":     strconv.Itoa(ca),
			"cb":     strconv.Itoa(cb),
			"pn":     strconv.Itoa(pn),
			"pd":     strconv.Itoa(pd),
			"answer": answer,
		},
	}
	pb.Explanation = fmt.Sprintf("Subgoal 1: State governing formula. By the law of total expectation, E[X] = E[E[X | Y]] = P(Y = 0) · E[X | Y = 0] + P(Y = 1) · E[X | Y = 1].\nSubgoal 2: Substitute parameters. The probabilities are P(Y = 0) = %d/%d and P(Y = 1) = %d/%d with conditional expectations E[X | Y = 0] = %d and E[X | Y = 1] = %d.\nSubgoal 3: Evaluate and simplify. E[X] = (%d/%d)(%d) + (%d/%d)(%d) = %s.",
		pd-pn, pd, pn, pd, ca, cb, pd-pn, pd, ca, pn, pd, cb, answer)
	return pb
}

// ---------------------------------------------------------------------------
// The law of total variance (numeric)
// ---------------------------------------------------------------------------

type totalVarianceGen struct{}

func (g *totalVarianceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	maxM := mathutil.ScaleInt(ctx.Difficulty, 5, 20)
	maxV := mathutil.ScaleInt(ctx.Difficulty, 4, 15)

	m0 := rand.Intn(maxM + 1)
	m1 := rand.Intn(maxM + 1)
	v0 := 1 + rand.Intn(maxV)
	v1 := 1 + rand.Intn(maxV)
	for m0 == m1 && v0 == v1 {
		m1 = rand.Intn(maxM + 1)
	}

	between := float64((m0-m1)*(m0-m1)) / 4.0
	within := float64(v0+v1) / 2.0
	answer := decimal6(within + between)

	p := generator.Problem{
		Question: fmt.Sprintf("A variable Y takes the values zero and one with equal probability. Given Y, the conditional mean of X is %d or %d and the conditional variance is %d or %d. What is Var(X)? (enter a decimal)",
			m0, m1, v0, v1),
		Answer: answer,
		Facts: map[string]string{
			"m0":     strconv.Itoa(m0),
			"m1":     strconv.Itoa(m1),
			"v0":     strconv.Itoa(v0),
			"v1":     strconv.Itoa(v1),
			"answer": answer,
		},
	}
	p.Explanation = fmt.Sprintf("Subgoal 1: State governing formula. The law of total variance (Eve's law) partitions variance as Var(X) = E[Var(X | Y)] + Var(E[X | Y]).\nSubgoal 2: Substitute parameters. Since Y takes states 0 and 1 with probability 1/2 each, E[Var(X | Y)] = (%d + %d)/2 = %s and Var(E[X | Y]) = (%d - %d)²/4 = %s.\nSubgoal 3: Evaluate and simplify. Var(X) = %s + %s = %s.",
		v0, v1, decimal6(within), m0, m1, decimal6(between), decimal6(within), decimal6(between), answer)
	return p
}

// ---------------------------------------------------------------------------
// Covariance (numeric)
// ---------------------------------------------------------------------------

type covarianceGen struct{}

func (g *covarianceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	tier := mathutil.DifficultyTier(ctx.Difficulty)
	var totalOptions []int
	switch tier {
	case 0:
		totalOptions = []int{10, 20}
	case 1:
		totalOptions = []int{20, 25, 40, 50}
	default:
		totalOptions = []int{50, 80, 100}
	}
	total := totalOptions[rand.Intn(len(totalOptions))]

	var w00, w01, w10, w11, num int
	for attempts := 0; attempts < 100; attempts++ {
		w11 = 1 + rand.Intn(total/2)
		rem := total - w11
		w10 = 1 + rand.Intn(rem/2)
		rem2 := rem - w10
		w01 = 1 + rand.Intn(rem2)
		w00 = rem2 - w01
		if w00 >= 1 {
			num = w11*total - (w10+w11)*(w01+w11)
			if num != 0 {
				break
			}
		}
	}
	if w00 < 1 || num == 0 {
		w00, w01, w10, w11 = 1, 1, 1, total-3
		num = w11*total - (w10+w11)*(w01+w11)
	}

	weights := []int{w00, w01, w10, w11}
	answer := decimal6(float64(num) / float64(total*total))
	parts := make([]string, len(weights))
	for i, w := range weights {
		parts[i] = strconv.Itoa(w)
	}

	p := generator.Problem{
		Question: "A joint mass function on the four pairs (0,0), (0,1), (1,0) and (1,1) is proportional to the weights " +
			strings.Join(parts, ", ") + " in that order. What is Cov(X, Y) = E[XY] − E[X]E[Y]? (enter a decimal)",
		Answer: answer,
		Facts: map[string]string{
			"weights": strings.Join(parts, ","),
			"total":   strconv.Itoa(total),
			"answer":  answer,
		},
	}
	p.Explanation = fmt.Sprintf("Subgoal 1: State governing formula. For indicator variables, Cov(X, Y) = E[XY] - E[X]E[Y] = (w11/T) - ((w10 + w11)/T)((w01 + w11)/T) = (w11 · T - (w10 + w11)(w01 + w11)) / T².\nSubgoal 2: Substitute parameters. Total weight T = %d, joint weight w11 = %d, marginal weight for X is %d and marginal weight for Y is %d.\nSubgoal 3: Evaluate and simplify. Cov(X, Y) = (%d · %d - %d · %d) / %d² = (%d - %d) / %d = %d/%d = %s.",
		total, w11, w10+w11, w01+w11, w11, total, w10+w11, w01+w11, total, w11*total, (w10+w11)*(w01+w11), total*total, num, total*total, answer)
	return p
}

// ---------------------------------------------------------------------------
// Correlation (numeric)
// ---------------------------------------------------------------------------

type correlationGen struct{}

func (g *correlationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	tier := mathutil.DifficultyTier(ctx.Difficulty)
	var sxOptions, syOptions []int
	switch tier {
	case 0:
		sxOptions = []int{2, 4, 5}
		syOptions = []int{2, 4, 5}
	case 1:
		sxOptions = []int{2, 4, 5, 8, 10}
		syOptions = []int{2, 4, 5, 8, 10}
	default:
		sxOptions = []int{4, 5, 8, 10, 20}
		syOptions = []int{4, 5, 8, 10, 20}
	}
	sx := sxOptions[rand.Intn(len(sxOptions))]
	sy := syOptions[rand.Intn(len(syOptions))]

	maxCov := sx * sy
	cov := 1 + rand.Intn(maxCov)
	if rand.Intn(3) == 0 {
		cov = -cov
	}

	answer := decimal6(float64(cov) / float64(sx*sy))
	p := generator.Problem{
		Question: fmt.Sprintf("X and Y have covariance %d and standard deviations %d and %d. What is their correlation? (enter a decimal)",
			cov, sx, sy),
		Answer: answer,
		Facts: map[string]string{
			"cov":    strconv.Itoa(cov),
			"sx":     strconv.Itoa(sx),
			"sy":     strconv.Itoa(sy),
			"answer": answer,
		},
	}
	p.Explanation = fmt.Sprintf("Subgoal 1: State governing formula. The Pearson correlation coefficient is defined as ρ = Cov(X, Y) / (σX · σY).\nSubgoal 2: Substitute parameters. Substitute Cov(X, Y) = %d, σX = %d and σY = %d.\nSubgoal 3: Evaluate and simplify. ρ = %d / (%d · %d) = %d/%d = %s.",
		cov, sx, sy, cov, sx, sy, cov, sx*sy, answer)
	return p
}

// ---------------------------------------------------------------------------
// Convergence in probability (multiple choice)
// ---------------------------------------------------------------------------

type convergenceInProbabilityGen struct{}

func (g *convergenceInProbabilityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return convergenceInProbabilityProblem()
}

// ---------------------------------------------------------------------------
// Convergence in distribution (multiple choice)
// ---------------------------------------------------------------------------

type convergenceInDistributionGen struct{}

func (g *convergenceInDistributionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return convergenceInDistributionProblem()
}

// ---------------------------------------------------------------------------
// Almost-sure convergence (multiple choice)
// ---------------------------------------------------------------------------

type almostSureGen struct{}

func (g *almostSureGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return almostSureProblem()
}

// ---------------------------------------------------------------------------
// The strong law of large numbers (multiple choice)
// ---------------------------------------------------------------------------

type strongLawGen struct{}

func (g *strongLawGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return strongLawProblem()
}

// ---------------------------------------------------------------------------
// The Borel-Cantelli lemmas (multiple choice)
// ---------------------------------------------------------------------------

type borelCantelliGen struct{}

func (g *borelCantelliGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return borelCantelliProblem()
}

// ---------------------------------------------------------------------------
// Slutsky's theorem (multiple choice)
// ---------------------------------------------------------------------------

type slutskyGen struct{}

func (g *slutskyGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return slutskyProblem()
}

// ---------------------------------------------------------------------------
// The delta method (numeric)
// ---------------------------------------------------------------------------

type deltaMethodGen struct{}

func (g *deltaMethodGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	tier := mathutil.DifficultyTier(ctx.Difficulty)
	fn := "square"
	if rand.Intn(2) == 1 {
		fn = "sqrt"
	}

	var mu, sig int
	var answer string
	var label string
	var derivStr, evalDerivStr string

	if fn == "square" {
		label = "x squared"
		switch tier {
		case 0:
			mu = 1 + rand.Intn(4)
			sig = 1 + rand.Intn(3)
		case 1:
			mu = 3 + rand.Intn(6)
			sig = 2 + rand.Intn(4)
		default:
			mu = 5 + rand.Intn(10)
			sig = 3 + rand.Intn(6)
		}
		derivVal := 2 * mu
		asympVar := float64(4 * mu * mu * sig * sig)
		answer = decimal6(asympVar)
		derivStr = "g'(x) = 2x"
		evalDerivStr = fmt.Sprintf("g'(%d) = 2(%d) = %d", mu, mu, derivVal)
	} else {
		label = "the square root of x"
		var muChoices []int
		switch tier {
		case 0:
			muChoices = []int{1, 4}
			sig = 1 + rand.Intn(4)
		case 1:
			muChoices = []int{1, 2, 4, 5, 8, 10}
			sig = 2 + rand.Intn(5)
		default:
			muChoices = []int{4, 5, 8, 10, 16, 20, 25}
			sig = 2 + rand.Intn(8)
		}
		mu = muChoices[rand.Intn(len(muChoices))]
		asympVar := float64(sig*sig) / float64(4*mu)
		answer = decimal6(asympVar)
		derivStr = "g'(x) = 1 / (2√x)"
		evalDerivStr = fmt.Sprintf("g'(%d) = 1 / (2√%d)", mu, mu)
	}

	p := generator.Problem{
		Question: fmt.Sprintf("A sample mean of n observations has mean μ = %d and variance σ² = %d. The delta method gives the asymptotic variance of g(sample mean) as σ² times the square of the derivative of g at μ. For g(x) equal to %s, what is that asymptotic variance? (enter a decimal)",
			mu, sig*sig, label),
		Answer: answer,
		Facts: map[string]string{
			"fn":     fn,
			"mu":     strconv.Itoa(mu),
			"sigma":  strconv.Itoa(sig),
			"answer": answer,
		},
	}
	p.Explanation = fmt.Sprintf("Subgoal 1: State governing formula. The delta method states that the asymptotic variance of g(X̄_n) is [g'(μ)]² · σ².\nSubgoal 2: Substitute parameters. For g(x) equal to %s, the derivative is %s, which evaluated at μ = %d gives %s with base variance σ² = %d.\nSubgoal 3: Evaluate and simplify. The asymptotic variance is [g'(%d)]² · %d = %s.",
		label, derivStr, mu, evalDerivStr, sig*sig, mu, sig*sig, answer)
	return p
}

// ---------------------------------------------------------------------------
// Characteristic functions (multiple choice)
// ---------------------------------------------------------------------------

type characteristicGen struct{}

func (g *characteristicGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return characteristicFunctionProblem()
}

// ---------------------------------------------------------------------------
// Moments from the characteristic function (numeric)
// ---------------------------------------------------------------------------

type characteristicMomentsGen struct{}

func (g *characteristicMomentsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	tier := mathutil.DifficultyTier(ctx.Difficulty)
	var denomOptions []int
	switch tier {
	case 0:
		denomOptions = []int{2, 4, 5, 10}
	case 1:
		denomOptions = []int{8, 16, 20, 25, 40, 50}
	default:
		denomOptions = []int{80, 100, 125, 200, 250, 400, 500, 1000}
	}
	pd := denomOptions[rand.Intn(len(denomOptions))]
	pn := 1 + rand.Intn(pd-1)
	for mathutil.GCD(pn, pd) > 1 {
		pn = 1 + rand.Intn(pd-1)
	}

	p := float64(pn) / float64(pd)
	answer := decimal6(p * (1 - p))

	pb := generator.Problem{
		Question: fmt.Sprintf("A random variable X takes the value zero with probability 1 − p and one with probability p, so its characteristic function is (1 − p) + p e^{it}. For p = %d/%d, what is Var(X)? (enter a decimal)",
			pn, pd),
		Answer: answer,
		Facts: map[string]string{
			"pn":     strconv.Itoa(pn),
			"pd":     strconv.Itoa(pd),
			"answer": answer,
		},
	}
	pb.Explanation = fmt.Sprintf("Subgoal 1: State governing formula. Derivatives of the characteristic function at the origin evaluate moments: E[X] = -i φ'(0) = p and E[X²] = -φ''(0) = p, giving Var(X) = E[X²] - (E[X])² = p(1 - p).\nSubgoal 2: Substitute parameters. Substitute success probability p = %d/%d, so failure probability is 1 - p = %d/%d.\nSubgoal 3: Evaluate and simplify. Var(X) = (%d/%d)(%d/%d) = %d/%d = %s.",
		pn, pd, pd-pn, pd, pn, pd, pd-pn, pd, pn*(pd-pn), pd*pd, answer)
	return pb
}

// ---------------------------------------------------------------------------
// Characteristic functions of sums (multiple choice)
// ---------------------------------------------------------------------------

type characteristicSumsGen struct{}

func (g *characteristicSumsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return characteristicSumsProblem()
}

// ---------------------------------------------------------------------------
// The continuity theorem (multiple choice)
// ---------------------------------------------------------------------------

type continuityTheoremGen struct{}

func (g *continuityTheoremGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return continuityTheoremProblem()
}

// ---------------------------------------------------------------------------
// The central limit theorem by characteristic functions (multiple choice)
// ---------------------------------------------------------------------------

type characteristicCLTGen struct{}

func (g *characteristicCLTGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return characteristicCLTProblem()
}

// ---------------------------------------------------------------------------
// Markov chains (numeric)
// ---------------------------------------------------------------------------

type markovChainGen struct{}

func (g *markovChainGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	tier := mathutil.DifficultyTier(ctx.Difficulty)
	var p01Val, p11Val float64
	switch tier {
	case 0:
		k := 1 + rand.Intn(9)
		p01Val = float64(k) / 10.0
		m := 1 + rand.Intn(9)
		p11Val = float64(m) / 10.0
	case 1:
		k := 1 + rand.Intn(19)
		p01Val = float64(k) * 0.05
		m := 1 + rand.Intn(19)
		p11Val = float64(m) * 0.05
	default:
		k := 2 + rand.Intn(96)
		p01Val = float64(k) / 100.0
		m := 2 + rand.Intn(96)
		p11Val = float64(m) / 100.0
	}
	p00Val := 1.0 - p01Val

	p00 := decimal6(p00Val)
	p01 := decimal6(p01Val)
	p11 := decimal6(p11Val)

	f := func(s string) float64 { v, _ := strconv.ParseFloat(s, 64); return v }
	v00, v01, v11 := f(p00), f(p01), f(p11)
	answer := decimal6(v00*v01 + v01*v11)

	p := generator.Problem{
		Question: fmt.Sprintf("A two-state Markov chain has P(state 0 to state 0) = %s, P(state 0 to state 1) = %s and P(state 1 to state 1) = %s. What is the two-step probability P(X_2 = 1 | X_0 = 0)? (enter a decimal)",
			p00, p01, p11),
		Answer: answer,
		Facts: map[string]string{
			"p00":    p00,
			"p01":    p01,
			"p11":    p11,
			"answer": answer,
		},
	}
	p.Explanation = fmt.Sprintf("Subgoal 1: State governing formula. By the Chapman-Kolmogorov equations, the two-step transition probability sums over the intermediate state X_1: P(X_2 = 1 | X_0 = 0) = P(0 → 0)P(0 → 1) + P(0 → 1)P(1 → 1) = P_00 · P_01 + P_01 · P_11.\nSubgoal 2: Substitute parameters. Substitute transition probabilities P_00 = %s, P_01 = %s and P_11 = %s.\nSubgoal 3: Evaluate and simplify. P(X_2 = 1 | X_0 = 0) = (%s)(%s) + (%s)(%s) = %s + %s = %s.",
		p00, p01, p11, p00, p01, p01, p11, decimal6(v00*v01), decimal6(v01*v11), answer)
	return p
}

// ---------------------------------------------------------------------------
// Transition matrices (numeric)
// ---------------------------------------------------------------------------

type transitionMatrixGen struct{}

func (g *transitionMatrixGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	tier := mathutil.DifficultyTier(ctx.Difficulty)
	var pi0Val, p01Val, p11Val float64
	switch tier {
	case 0:
		k := 1 + rand.Intn(9)
		pi0Val = float64(k) / 10.0
		a := 1 + rand.Intn(9)
		p01Val = float64(a) / 10.0
		b := 1 + rand.Intn(9)
		p11Val = float64(b) / 10.0
	case 1:
		k := 1 + rand.Intn(19)
		pi0Val = float64(k) * 0.05
		a := 1 + rand.Intn(19)
		p01Val = float64(a) * 0.05
		b := 1 + rand.Intn(19)
		p11Val = float64(b) * 0.05
	default:
		k := 2 + rand.Intn(96)
		pi0Val = float64(k) / 100.0
		a := 2 + rand.Intn(96)
		p01Val = float64(a) / 100.0
		b := 2 + rand.Intn(96)
		p11Val = float64(b) / 100.0
	}
	pi1Val := 1.0 - pi0Val

	pi0 := decimal6(pi0Val)
	pi1 := decimal6(pi1Val)
	p01 := decimal6(p01Val)
	p11 := decimal6(p11Val)

	f := func(s string) float64 { v, _ := strconv.ParseFloat(s, 64); return v }
	vPi0, vPi1, v01, v11 := f(pi0), f(pi1), f(p01), f(p11)
	answer := decimal6(vPi0*v01 + vPi1*v11)

	p := generator.Problem{
		Question: fmt.Sprintf("A chain starts with P(X_0 = 0) = %s and P(X_0 = 1) = %s. Its one-step probabilities are P(state 0 to state 1) = %s and P(state 1 to state 1) = %s. What is P(X_1 = 1)? (enter a decimal)",
			pi0, pi1, p01, p11),
		Answer: answer,
		Facts: map[string]string{
			"pi0":    pi0,
			"pi1":    pi1,
			"p01":    p01,
			"p11":    p11,
			"answer": answer,
		},
	}
	p.Explanation = fmt.Sprintf("Subgoal 1: State governing formula. The state distribution after one transition is obtained by multiplying the initial distribution by the transition matrix: P(X_1 = 1) = P(X_0 = 0)P(0 → 1) + P(X_0 = 1)P(1 → 1) = π_0 · P_01 + π_1 · P_11.\nSubgoal 2: Substitute parameters. Substitute initial probabilities π_0 = %s and π_1 = %s with transition rates P_01 = %s and P_11 = %s.\nSubgoal 3: Evaluate and simplify. P(X_1 = 1) = (%s)(%s) + (%s)(%s) = %s + %s = %s.",
		pi0, pi1, p01, p11, pi0, p01, pi1, p11, decimal6(vPi0*v01), decimal6(vPi1*v11), answer)
	return p
}

// ---------------------------------------------------------------------------
// Stationary distributions (numeric)
// ---------------------------------------------------------------------------

type stationaryGen struct{}

func (g *stationaryGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	tier := mathutil.DifficultyTier(ctx.Difficulty)
	type pair struct{ up, down string }
	var pool []pair
	switch tier {
	case 0:
		pool = []pair{
			{"0.3", "0.7"}, {"0.5", "0.5"}, {"0.4", "0.6"}, {"0.1", "0.9"},
			{"0.2", "0.8"}, {"0.6", "0.4"}, {"0.7", "0.3"}, {"0.8", "0.2"},
			{"0.1", "0.4"}, {"0.2", "0.3"}, {"0.3", "0.2"}, {"0.4", "0.1"},
		}
	case 1:
		pool = []pair{
			{"0.2", "0.6"}, {"0.6", "0.2"}, {"0.1", "0.7"}, {"0.3", "0.5"},
			{"0.5", "0.3"}, {"0.7", "0.1"}, {"0.1", "0.3"}, {"0.3", "0.1"},
			{"0.15", "0.85"}, {"0.35", "0.65"}, {"0.45", "0.55"}, {"0.25", "0.75"},
		}
	default:
		pool = []pair{
			{"0.25", "0.55"}, {"0.05", "0.75"}, {"0.45", "0.35"}, {"0.65", "0.15"},
			{"0.12", "0.68"}, {"0.28", "0.52"}, {"0.36", "0.44"}, {"0.48", "0.32"},
			{"0.16", "0.64"}, {"0.32", "0.48"}, {"0.64", "0.16"}, {"0.72", "0.08"},
		}
	}
	selected := pool[rand.Intn(len(pool))]
	p01, p10 := selected.up, selected.down

	f := func(s string) float64 { v, _ := strconv.ParseFloat(s, 64); return v }
	up, down := f(p01), f(p10)
	sumVal := up + down
	answer := decimal6(up / sumVal)

	p := generator.Problem{
		Question: fmt.Sprintf("A two-state Markov chain has P(state 0 to state 1) = %s and P(state 1 to state 0) = %s. What is the stationary probability of state 1? (enter a decimal)",
			p01, p10),
		Answer: answer,
		Facts: map[string]string{
			"p01":    p01,
			"p10":    p10,
			"answer": answer,
		},
	}
	p.Explanation = fmt.Sprintf("Subgoal 1: State governing formula. The stationary distribution balances flow between states: π_0 · P_01 = π_1 · P_10. Since π_0 + π_1 = 1, substituting π_0 = 1 - π_1 yields π_1 = P_01 / (P_01 + P_10).\nSubgoal 2: Substitute parameters. Substitute transition probabilities P_01 = %s and P_10 = %s.\nSubgoal 3: Evaluate and simplify. π_1 = %s / (%s + %s) = %s / %s = %s.",
		p01, p10, p01, p01, p10, p01, decimal6(sumVal), answer)
	return p
}

// ---------------------------------------------------------------------------
// Random walks (numeric)
// ---------------------------------------------------------------------------

type randomWalkGen struct{}

func (g *randomWalkGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	tier := mathutil.DifficultyTier(ctx.Difficulty)
	var denomOptions []int
	var n int
	switch tier {
	case 0:
		denomOptions = []int{2, 4, 5}
		n = 5 + rand.Intn(20)
	case 1:
		denomOptions = []int{4, 5, 8, 10}
		n = 20 + rand.Intn(45)
	default:
		denomOptions = []int{8, 10, 16, 20}
		n = 50 + rand.Intn(150)
	}
	pd := denomOptions[rand.Intn(len(denomOptions))]
	pn := 1 + rand.Intn(pd-1)
	for mathutil.GCD(pn, pd) > 1 {
		pn = 1 + rand.Intn(pd-1)
	}

	p := float64(pn) / float64(pd)
	answer := decimal6(4 * float64(n) * p * (1 - p))

	pb := generator.Problem{
		Question: fmt.Sprintf("A random walk takes a step of plus one with probability p = %d/%d and minus one otherwise. What is Var(S_n) after n = %d steps? (enter a decimal)",
			pn, pd, n),
		Answer: answer,
		Facts: map[string]string{
			"n":      strconv.Itoa(n),
			"pn":     strconv.Itoa(pn),
			"pd":     strconv.Itoa(pd),
			"answer": answer,
		},
	}
	pb.Explanation = fmt.Sprintf("Subgoal 1: State governing formula. A single step X_1 takes values ±1 with probability p and 1 - p, having variance Var(X_1) = 1 - (2p - 1)² = 4p(1 - p). For n independent identically distributed steps, Var(S_n) = n · Var(X_1) = 4np(1 - p).\nSubgoal 2: Substitute parameters. Substitute number of steps n = %d with step success probability p = %d/%d and failure probability 1 - p = %d/%d.\nSubgoal 3: Evaluate and simplify. Var(S_n) = 4(%d)(%d/%d)(%d/%d) = %s.",
		n, pn, pd, pd-pn, pd, n, pn, pd, pd-pn, pd, answer)
	return pb
}

// ---------------------------------------------------------------------------
// Martingales (multiple choice)
// ---------------------------------------------------------------------------

type martingaleGen struct{}

func (g *martingaleGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return martingaleProblem()
}
