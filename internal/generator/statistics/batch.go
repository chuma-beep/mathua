package statistics

// The 2026-10-18 statistics batch: eleven concepts that take the domain from 44
// to the target of 55.
//
// The domain already computes means, spreads, probabilities and distributions,
// and it has sampling, confidence intervals, standard error and margin of
// error. What it did not have is the testing framework those tools feed: null
// and alternative hypotheses, z and t statistics, p-values, Type I and II
// errors, proportions, goodness of fit and ANOVA. The batch adds that spine,
// plus the coefficient of determination, histograms, and the nonparametric
// alternatives that do not assume normality.
//
// Several obvious additions were already taught and were not repeated: standard
// deviation and the correlation coefficient live in the variance/covariance
// lesson, sample size in confidence intervals, five-number summaries in the
// quartiles lesson, the regression line in machine learning, and the central
// limit theorem in sampling distributions.
//
// Six concepts are numeric and five are multiple choice. Every numeric
// generator publishes the facts its answer follows from, so the tests recompute
// the statistic (z, t, standard error, chi-square, F, R-squared, or a histogram
// count) instead of comparing a string with itself.

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerBatch(reg *generator.Registry) {
	reg.Register("stat.infer.hypothesis", &hypothesisGen{})
	reg.Register("stat.infer.ztest", &zTestGen{})
	reg.Register("stat.infer.ttest", &tTestGen{})
	reg.Register("stat.infer.pvalue", &pValueGen{})
	reg.Register("stat.infer.type_error", &typeErrorGen{})
	reg.Register("stat.infer.proportion", &proportionSEGen{})
	reg.Register("stat.test.chi_square", &chiSquareTestGen{})
	reg.Register("stat.test.anova", &anovaGen{})
	reg.Register("stat.reg.r2", &r2Gen{})
	reg.Register("stat.data.histogram", &histogramGen{})
	reg.Register("stat.infer.nonparametric", &nonparametricGen{})
}

// ---------------------------------------------------------------------------
// 1. Null and alternative hypotheses
// ---------------------------------------------------------------------------

type hypothesisGen struct{}

func (g *hypothesisGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question    string
		correct     string
		distractors []string
		reason      string
	}
	claims := []claim{
		{
			question:    "The null hypothesis is...",
			correct:     "the default claim the test looks for evidence against",
			distractors: []string{"the claim the researcher hopes to prove", "always identical to the alternative", "the observed effect size"},
			reason:      "The null is the baseline; a small p-value is evidence against it, not proof that the alternative is true.",
		},
		{
			question:    "The alternative hypothesis is...",
			correct:     "the claim accepted when the null is rejected",
			distractors: []string{"the default claim", "always two-sided", "the significance level"},
			reason:      "The alternative states the departure the test has power to detect, and it is chosen before seeing the data.",
		},
		{
			question:    "A one-sided test is appropriate when...",
			correct:     "only departures in one direction matter",
			distractors: []string{"both directions matter equally", "the sample is small", "the population is normal"},
			reason:      "Putting all the rejection region in one tail buys power there, at the cost of ignoring the other direction.",
		},
		{
			question:    "The significance level alpha is...",
			correct:     "the maximum Type I error rate the test allows",
			distractors: []string{"the probability the null is true", "the power of the test", "the observed p-value"},
			reason:      "Alpha is fixed before the data arrive and caps the chance of rejecting a true null hypothesis.",
		},
		{
			question:    "The test statistic measures...",
			correct:     "how far the estimate lies from the null value in standard-error units",
			distractors: []string{"the probability the null is true", "the sample mean only", "the size of the population"},
			reason:      "A z or t statistic standardises the gap between the estimate and the null value, so it can be compared with a reference distribution.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	p := generator.MultipleChoice(c.question, c.correct, c.distractors...)
	p.Explanation = c.reason
	return p
}

// ---------------------------------------------------------------------------
// 2. The one-sample z statistic (numeric)
// ---------------------------------------------------------------------------

type zTestGen struct{}

func (g *zTestGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	nVals := []int{4, 9, 16, 25, 36, 64, 100}
	seVals := []int{1, 2, 4, 5}
	zSteps := []int{-7, -6, -5, -4, -3, -2, -1, 1, 2, 3, 4, 5, 6, 7}

	var question, answer, reason string
	var xbar, mu, sigma float64
	var n int

	// Bounded rather than open ended: see generator.MaxDrawAttempts. The
	// explanation is built on every draw so that the last iteration still
	// returns a complete problem if no draw avoids printing its answer.
	for attempt := 0; attempt < generator.MaxDrawAttempts; attempt++ {
		n = nVals[rand.Intn(len(nVals))]
		sqrtN := int(math.Round(math.Sqrt(float64(n))))

		se := seVals[rand.Intn(len(seVals))]
		sigma = float64(se * sqrtN)

		step := zSteps[rand.Intn(len(zSteps))]
		zVal := float64(step) * 0.5

		mu = float64((rand.Intn(15) + 5) * 10)
		xbar = mu + zVal*float64(se)

		answer = strconv.FormatFloat(zVal, 'f', -1, 64)
		qType := "a number"
		if strings.Contains(answer, ".") {
			qType = "a decimal"
		}
		question = fmt.Sprintf("A sample mean is %s, the null mean is %s, the population standard deviation is %s and the sample size is %d. What is the z statistic (xbar - mu)/(sigma/sqrt(n))? (enter %s)",
			strconv.FormatFloat(xbar, 'f', -1, 64),
			strconv.FormatFloat(mu, 'f', -1, 64),
			strconv.FormatFloat(sigma, 'f', -1, 64),
			n,
			qType,
		)

		xbarStr := strconv.FormatFloat(xbar, 'f', -1, 64)
		muStr := strconv.FormatFloat(mu, 'f', -1, 64)
		sigmaStr := strconv.FormatFloat(sigma, 'f', -1, 64)
		nStr := strconv.Itoa(n)

		reason = fmt.Sprintf("The standard error is %s/sqrt(%d) = %d, so z = (%s - %s)/%d = %s.",
			sigmaStr,
			n,
			se,
			xbarStr,
			muStr,
			se,
			answer,
		)

		if answer != xbarStr && answer != muStr && answer != sigmaStr && answer != nStr &&
			answer != "-"+xbarStr && answer != "-"+muStr && !strings.Contains(question, answer) {
			break
		}
	}

	return generator.Problem{
		Question:    question,
		Answer:      answer,
		Explanation: reason,
		Facts: map[string]string{
			"xbar":   strconv.FormatFloat(xbar, 'g', -1, 64),
			"mu":     strconv.FormatFloat(mu, 'g', -1, 64),
			"sigma":  strconv.FormatFloat(sigma, 'g', -1, 64),
			"n":      strconv.Itoa(n),
			"answer": answer,
		},
	}
}

// ---------------------------------------------------------------------------
// 3. The one-sample t statistic (numeric)
// ---------------------------------------------------------------------------

type tTestGen struct{}

func (g *tTestGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		xbar     float64
		mu       float64
		s        float64
		n        int
		reason   string
	}
	claims := []claim{
		{
			question: "A sample mean is 52, the null mean is 50, the sample standard deviation is 4 and the sample size is 4. What is the t statistic (xbar - mu)/(s/sqrt(n))? (enter a number)",
			answer:   "1", xbar: 52, mu: 50, s: 4, n: 4,
			reason: "The standard error is 4/sqrt(4) = 2, so t = (52 - 50)/2 = 1.",
		},
		{
			question: "A sample mean is 54, the null mean is 50, the sample standard deviation is 6 and the sample size is 9. What is the t statistic (xbar - mu)/(s/sqrt(n))? (enter a number)",
			answer:   "2", xbar: 54, mu: 50, s: 6, n: 9,
			reason: "The standard error is 6/3 = 2, so t = 4/2 = 2.",
		},
		{
			question: "A sample mean is 47, the null mean is 50, the sample standard deviation is 4 and the sample size is 16. What is the t statistic (xbar - mu)/(s/sqrt(n))? (enter a number)",
			answer:   "-3", xbar: 47, mu: 50, s: 4, n: 16,
			reason: "The standard error is 4/4 = 1, so t = -3.",
		},
		{
			question: "A sample mean is 12, the null mean is 10, the sample standard deviation is 5 and the sample size is 25. What is the t statistic (xbar - mu)/(s/sqrt(n))? (enter a number)",
			answer:   "2", xbar: 12, mu: 10, s: 5, n: 25,
			reason: "The standard error is 5/5 = 1, so t = 2.",
		},
		{
			question: "A sample mean is 8, the null mean is 10, the sample standard deviation is 6 and the sample size is 9. What is the t statistic (xbar - mu)/(s/sqrt(n))? (enter a number)",
			answer:   "-1", xbar: 8, mu: 10, s: 6, n: 9,
			reason: "The standard error is 6/3 = 2, so t = -2/2 = -1.",
		},
		{
			question: "A sample mean is 21, the null mean is 20, the sample standard deviation is 4 and the sample size is 16. What is the t statistic (xbar - mu)/(s/sqrt(n))? (enter a number)",
			answer:   "1", xbar: 21, mu: 20, s: 4, n: 16,
			reason: "The standard error is 4/4 = 1, so t = 1.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	return generator.Problem{
		Question:    c.question,
		Answer:      c.answer,
		Explanation: c.reason,
		Facts: map[string]string{
			"xbar": strconv.FormatFloat(c.xbar, 'g', -1, 64),
			"mu":   strconv.FormatFloat(c.mu, 'g', -1, 64),
			"s":    strconv.FormatFloat(c.s, 'g', -1, 64),
			"n":    strconv.Itoa(c.n), "answer": c.answer,
		},
	}
}

// ---------------------------------------------------------------------------
// 4. p-values and significance
// ---------------------------------------------------------------------------

type pValueGen struct{}

func (g *pValueGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question    string
		correct     string
		distractors []string
		reason      string
	}
	claims := []claim{
		{
			question:    "A p-value is...",
			correct:     "the probability, under the null, of data at least as extreme as observed",
			distractors: []string{"the probability the null is true", "the probability the alternative is true", "the size of the effect"},
			reason:      "The p-value is computed assuming the null is true; it measures compatibility with the null, not the chance that the null is true.",
		},
		{
			question:    "A result with p = 0.03 tested at alpha = 0.05 is...",
			correct:     "statistically significant",
			distractors: []string{"not significant", "proof that the null is true", "a Type II error"},
			reason:      "A p-value below alpha falls in the rejection region, so the null is rejected at that level.",
		},
		{
			question:    "A small p-value means...",
			correct:     "the observed data are unlikely under the null",
			distractors: []string{"the null is certainly false", "the effect is necessarily large", "the sample is biased"},
			reason:      "Smallness is evidence against the null, but a tiny effect in a large sample can also produce a small p-value.",
		},
		{
			question:    "A two-sided p-value counts...",
			correct:     "both tails of the null distribution",
			distractors: []string{"only the upper tail", "only the lower tail", "the whole distribution"},
			reason:      "A two-sided alternative treats departures in either direction as evidence, so both tails contribute.",
		},
		{
			question:    "A result with p = 0.20 tested at alpha = 0.05 means...",
			correct:     "fail to reject the null",
			distractors: []string{"reject the null", "the alternative is proven", "the test is invalid"},
			reason:      "A p-value above alpha leaves the null in place; failing to reject is not the same as proving the null.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	p := generator.MultipleChoice(c.question, c.correct, c.distractors...)
	p.Explanation = c.reason
	return p
}

// ---------------------------------------------------------------------------
// 5. Type I and II errors and power
// ---------------------------------------------------------------------------

type typeErrorGen struct{}

func (g *typeErrorGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question    string
		correct     string
		distractors []string
		reason      string
	}
	claims := []claim{
		{
			question:    "A Type I error is...",
			correct:     "rejecting a true null hypothesis",
			distractors: []string{"failing to reject a false null", "accepting a true null", "rejecting a false null"},
			reason:      "The Type I error is the false alarm: the null is true but the test rejects it, at rate alpha.",
		},
		{
			question:    "A Type II error is...",
			correct:     "failing to reject a false null hypothesis",
			distractors: []string{"rejecting a true null", "accepting a true null", "rejecting a false null"},
			reason:      "The Type II error is the missed detection: the null is false but the test does not reject it.",
		},
		{
			question:    "The power of a test is...",
			correct:     "the probability of rejecting a false null",
			distractors: []string{"the probability of a Type I error", "the significance level", "the probability the null is true"},
			reason:      "Power is one minus the Type II error rate, so it measures the chance of detecting a real effect.",
		},
		{
			question:    "Lowering alpha while holding the sample size fixed...",
			correct:     "reduces Type I errors and raises Type II errors",
			distractors: []string{"reduces both kinds of error", "raises both kinds of error", "has no effect on either"},
			reason:      "A stricter threshold makes false alarms rarer but missed detections more common, unless the sample grows.",
		},
		{
			question:    "Power increases when...",
			correct:     "the sample size increases",
			distractors: []string{"alpha decreases", "the effect size decreases", "the variance increases"},
			reason:      "More data shrink the standard error, so a real effect is easier to distinguish from noise.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	p := generator.MultipleChoice(c.question, c.correct, c.distractors...)
	p.Explanation = c.reason
	return p
}

// ---------------------------------------------------------------------------
// 6. Standard error of a proportion (numeric)
// ---------------------------------------------------------------------------

type proportionSEGen struct{}

func (g *proportionSEGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		p        float64
		n        int
		reason   string
	}
	claims := []claim{
		{
			question: "A poll of 100 respondents gives p-hat = 0.5. What is the standard error sqrt(p-hat(1 - p-hat)/n)? (enter a decimal)",
			answer:   "0.05", p: 0.5, n: 100,
			reason: "0.5(0.5)/100 = 0.0025, whose square root is 0.05.",
		},
		{
			question: "A poll of 400 respondents gives p-hat = 0.5. What is the standard error sqrt(p-hat(1 - p-hat)/n)? (enter a decimal)",
			answer:   "0.025", p: 0.5, n: 400,
			reason: "0.25/400 = 0.000625, whose square root is 0.025.",
		},
		{
			question: "A poll of 100 respondents gives p-hat = 0.2. What is the standard error sqrt(p-hat(1 - p-hat)/n)? (enter a decimal)",
			answer:   "0.04", p: 0.2, n: 100,
			reason: "0.2(0.8)/100 = 0.0016, whose square root is 0.04.",
		},
		{
			question: "A poll of 100 respondents gives p-hat = 0.8. What is the standard error sqrt(p-hat(1 - p-hat)/n)? (enter a decimal)",
			answer:   "0.04", p: 0.8, n: 100,
			reason: "0.8(0.2)/100 = 0.0016, whose square root is 0.04.",
		},
		{
			question: "A poll of 25 respondents gives p-hat = 0.2. What is the standard error sqrt(p-hat(1 - p-hat)/n)? (enter a decimal)",
			answer:   "0.08", p: 0.2, n: 25,
			reason: "0.2(0.8)/25 = 0.0064, whose square root is 0.08.",
		},
		{
			question: "A poll of 64 respondents gives p-hat = 0.5. What is the standard error sqrt(p-hat(1 - p-hat)/n)? (enter a decimal)",
			answer:   "0.0625", p: 0.5, n: 64,
			reason: "0.25/64 = 0.00390625, whose square root is 0.0625.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	return generator.Problem{
		Question:    c.question,
		Answer:      c.answer,
		Explanation: c.reason,
		Facts: map[string]string{
			"p": strconv.FormatFloat(c.p, 'g', -1, 64), "n": strconv.Itoa(c.n), "answer": c.answer,
		},
	}
}

// ---------------------------------------------------------------------------
// 7. The chi-square goodness-of-fit test (numeric)
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// 8. One-way ANOVA (numeric)
// ---------------------------------------------------------------------------

type anovaGen struct{}

func (g *anovaGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		groups   [][]int
		reason   string
	}
	claims := []claim{
		{
			question: "A one-way ANOVA compares the groups 1, 3 and 3, 5 and 5, 7. What is the F statistic? (enter a number)",
			answer:   "4", groups: [][]int{{1, 3}, {3, 5}, {5, 7}},
			reason: "The group means are 2, 4 and 6 with grand mean 4; SSB = 16 on 2 degrees of freedom and SSW = 6 on 3, so F = (16/2)/(6/3) = 4.",
		},
		{
			question: "A one-way ANOVA compares the groups 2, 4 and 6, 8. What is the F statistic? (enter a number)",
			answer:   "8", groups: [][]int{{2, 4}, {6, 8}},
			reason: "The group means are 3 and 7 with grand mean 5; SSB = 16 on 1 degree of freedom and SSW = 4 on 2, so F = 16/(4/2) = 8.",
		},
		{
			question: "A one-way ANOVA compares the groups 1, 2, 3 and 2, 3, 4 and 3, 4, 5. What is the F statistic? (enter a number)",
			answer:   "3", groups: [][]int{{1, 2, 3}, {2, 3, 4}, {3, 4, 5}},
			reason: "The group means are 2, 3 and 4 with grand mean 3; SSB = 6 on 2 degrees of freedom and SSW = 6 on 6, so F = 3/1 = 3.",
		},
		{
			question: "A one-way ANOVA compares the groups 1, 3 and 3, 5. What is the F statistic? (enter a number)",
			answer:   "2", groups: [][]int{{1, 3}, {3, 5}},
			reason: "The group means are 2 and 4 with grand mean 3; SSB = 4 on 1 degree of freedom and SSW = 4 on 2, so F = 4/2 = 2.",
		},
		{
			question: "A one-way ANOVA compares the groups 0, 2 and 2, 4 and 4, 6. What is the F statistic? (enter a number)",
			answer:   "4", groups: [][]int{{0, 2}, {2, 4}, {4, 6}},
			reason: "The group means are 1, 3 and 5 with grand mean 3; SSB = 16 on 2 degrees of freedom and SSW = 6 on 3, so F = 8/2 = 4.",
		},
		{
			question: "A one-way ANOVA compares the groups 1, 5 and 3, 7. What is the F statistic? (enter a decimal)",
			answer:   "0.5", groups: [][]int{{1, 5}, {3, 7}},
			reason: "The group means are 3 and 5 with grand mean 4; SSB = 4 on 1 degree of freedom and SSW = 16 on 2, so F = 4/8 = 0.5.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	parts := make([]string, len(c.groups))
	for i, group := range c.groups {
		cells := make([]string, len(group))
		for j, v := range group {
			cells[j] = strconv.Itoa(v)
		}
		parts[i] = strings.Join(cells, ",")
	}
	return generator.Problem{
		Question:    c.question,
		Answer:      c.answer,
		Explanation: c.reason,
		Facts: map[string]string{
			"groups": strings.Join(parts, ";"), "answer": c.answer,
		},
	}
}

// ---------------------------------------------------------------------------
// 9. The coefficient of determination (numeric)
// ---------------------------------------------------------------------------

type r2Gen struct{}

func (g *r2Gen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		sst, sse int
		reason   string
	}
	claims := []claim{
		{
			question: "A regression has total sum of squares 100 and residual sum of squares 25. What is R^2 = 1 - SSE/SST? (enter a decimal)",
			answer:   "0.75", sst: 100, sse: 25,
			reason: "1 - 25/100 = 0.75.",
		},
		{
			question: "A regression has total sum of squares 80 and residual sum of squares 8. What is R^2 = 1 - SSE/SST? (enter a decimal)",
			answer:   "0.9", sst: 80, sse: 8,
			reason: "1 - 8/80 = 0.9.",
		},
		{
			question: "A regression has total sum of squares 50 and residual sum of squares 10. What is R^2 = 1 - SSE/SST? (enter a decimal)",
			answer:   "0.8", sst: 50, sse: 10,
			reason: "1 - 10/50 = 0.8.",
		},
		{
			question: "A regression has total sum of squares 200 and residual sum of squares 50. What is R^2 = 1 - SSE/SST? (enter a decimal)",
			answer:   "0.75", sst: 200, sse: 50,
			reason: "1 - 50/200 = 0.75.",
		},
		{
			question: "A regression has total sum of squares 40 and residual sum of squares 4. What is R^2 = 1 - SSE/SST? (enter a decimal)",
			answer:   "0.9", sst: 40, sse: 4,
			reason: "1 - 4/40 = 0.9.",
		},
		{
			question: "A regression has total sum of squares 25 and residual sum of squares 16. What is R^2 = 1 - SSE/SST? (enter a decimal)",
			answer:   "0.36", sst: 25, sse: 16,
			reason: "1 - 16/25 = 0.36.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	return generator.Problem{
		Question:    c.question,
		Answer:      c.answer,
		Explanation: c.reason,
		Facts: map[string]string{
			"sst": strconv.Itoa(c.sst), "sse": strconv.Itoa(c.sse), "answer": c.answer,
		},
	}
}

// ---------------------------------------------------------------------------
// 10. Histograms (numeric)
// ---------------------------------------------------------------------------

type histogramGen struct{}

func (g *histogramGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question  string
		answer    string
		values    []int
		low, high int
		reason    string
	}
	claims := []claim{
		{
			question: "A histogram of the values 5, 12, 14, 19, 22, 25 uses bins of width 10. How many values fall in the bin [10, 20)? (enter a number)",
			answer:   "3", values: []int{5, 12, 14, 19, 22, 25}, low: 10, high: 20,
			reason: "The values 12, 14 and 19 lie in [10, 20); 5 is below and 22, 25 are above.",
		},
		{
			question: "A histogram of the values 1, 3, 5, 7, 9 uses bins of width 4. How many values fall in the bin [4, 8)? (enter a number)",
			answer:   "2", values: []int{1, 3, 5, 7, 9}, low: 4, high: 8,
			reason: "The values 5 and 7 lie in [4, 8); 1, 3 are below and 9 is above.",
		},
		{
			question: "A histogram of the values 15, 16, 18, 21, 23, 25 uses bins of width 5. How many values fall in the bin [15, 20)? (enter a number)",
			answer:   "3", values: []int{15, 16, 18, 21, 23, 25}, low: 15, high: 20,
			reason: "The values 15, 16 and 18 lie in [15, 20); 21, 23 and 25 are above.",
		},
		{
			question: "A histogram of the values 2, 4, 6, 8, 10, 12 uses bins of width 6. How many values fall in the bin [4, 10)? (enter a number)",
			answer:   "3", values: []int{2, 4, 6, 8, 10, 12}, low: 4, high: 10,
			reason: "The values 4, 6 and 8 lie in [4, 10); 10 is excluded by the half-open bin.",
		},
		{
			question: "A histogram of the values 0, 5, 10, 15, 20, 25 uses bins of width 15. How many values fall in the bin [5, 20)? (enter a number)",
			answer:   "3", values: []int{0, 5, 10, 15, 20, 25}, low: 5, high: 20,
			reason: "The values 5, 10 and 15 lie in [5, 20); 0 is below and 20, 25 are outside.",
		},
		{
			question: "A histogram of the values 3, 8, 13, 18, 23, 28 uses bins of width 10. How many values fall in the bin [0, 10)? (enter a number)",
			answer:   "2", values: []int{3, 8, 13, 18, 23, 28}, low: 0, high: 10,
			reason: "The values 3 and 8 lie in [0, 10); the rest are 10 or above.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	vals := make([]string, len(c.values))
	for i, v := range c.values {
		vals[i] = strconv.Itoa(v)
	}
	return generator.Problem{
		Question:    c.question,
		Answer:      c.answer,
		Explanation: c.reason,
		Facts: map[string]string{
			"values": strings.Join(vals, ","), "low": strconv.Itoa(c.low),
			"high": strconv.Itoa(c.high), "answer": c.answer,
		},
	}
}

// ---------------------------------------------------------------------------
// 11. Nonparametric tests
// ---------------------------------------------------------------------------

type nonparametricGen struct{}

func (g *nonparametricGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question    string
		correct     string
		distractors []string
		reason      string
	}
	claims := []claim{
		{
			question:    "The sign test uses which distribution under the null hypothesis?",
			correct:     "the binomial distribution with p = 1/2",
			distractors: []string{"the normal distribution with mean zero", "the chi-square distribution with n - 1 degrees of freedom", "the Poisson distribution with mean n"},
			reason:      "Under the null, each pair is equally likely to favour either side, so the number of positive signs is binomial with p = 1/2.",
		},
		{
			question:    "The sign test ignores...",
			correct:     "the magnitudes of the differences, using only their signs",
			distractors: []string{"the signs of the differences", "the number of pairs", "whether the observations are paired"},
			reason:      "Discarding magnitudes makes the test robust to outliers but less powerful when the magnitudes carry information.",
		},
		{
			question:    "A rank-based test is preferred when...",
			correct:     "the data are ordinal or clearly non-normal",
			distractors: []string{"the data are exactly normal", "the sample is enormous", "the variances are known"},
			reason:      "Rank tests do not assume a normal population, so they are the safer choice when that assumption is doubtful.",
		},
		{
			question:    "The Wilcoxon signed-rank test differs from the sign test by...",
			correct:     "using the ranks of the absolute differences",
			distractors: []string{"using only the signs", "ignoring the pairing", "requiring a normal population"},
			reason:      "Ranks keep some information about magnitude while remaining robust, which usually gives more power than signs alone.",
		},
		{
			question:    "A nonparametric test typically has...",
			correct:     "less power when its parametric assumptions hold",
			distractors: []string{"more power in every case", "no null hypothesis", "no test statistic"},
			reason:      "The trade for fewer assumptions is a loss of efficiency when the parametric model is actually correct.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	p := generator.MultipleChoice(c.question, c.correct, c.distractors...)
	p.Explanation = c.reason
	return p
}
