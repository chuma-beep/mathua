package statistics

// Wave 4 statistics: ten concepts. Every family asks for an integer count,
// statistic or coded conclusion and publishes its operands in Facts.

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave4(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"stat.infer.bootstrap":          "bootstrap",
		"stat.test.permutation":         "permutation",
		"stat.test.wilcoxon":            "wilcoxon",
		"stat.test.mann_whitney":        "mann_whitney",
		"stat.infer.kaplan_meier":       "kaplan_meier",
		"stat.infer.credible_interval":  "credible_interval",
		"stat.reg.information_criteria": "information_criteria",
		"stat.reg.quantile":             "quantile_regression",
		"stat.reg.poisson":              "poisson_regression",
		"stat.data.missing":             "missing_data",
	} {
		reg.Register(id, &w4StatGen{kind: kind})
	}
}

type w4StatGen struct{ kind string }

func (g *w4StatGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "bootstrap":
		return w4Bootstrap(kp)
	case "permutation":
		return w4Permutation(kp)
	case "wilcoxon":
		return w4Wilcoxon(kp)
	case "mann_whitney":
		return w4MannWhitney(kp)
	case "kaplan_meier":
		return w4KaplanMeier(kp)
	case "credible_interval":
		return w4CredibleInterval(kp)
	case "information_criteria":
		return w4InformationCriteria(kp)
	case "quantile_regression":
		return w4QuantileRegression(kp)
	case "poisson_regression":
		return w4PoissonRegression(kp)
	case "missing_data":
		return w4MissingData(kp)
	}
	return generator.Problem{Question: "What is 2 + 2? (enter a number)", Answer: "4", Explanation: "2 + 2 = 4."}
}

func w4S(o int, q, a, ex string, facts map[string]string) generator.Problem {
	if facts == nil {
		facts = map[string]string{}
	}
	facts["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: facts, KPOrdinal: o}
}

func w4SI(m map[string]string, k string, v int) { m[k] = strconv.Itoa(v) }

func w4SR(lo, hi int) int {
	if hi <= lo {
		return lo
	}
	return lo + rand.Intn(hi-lo+1)
}

func w4Abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func w4SSum(xs []int) int {
	s := 0
	for _, x := range xs {
		s += x
	}
	return s
}

func w4Bootstrap(kp int) generator.Problem {
	xs := []int{w4SR(1, 9), w4SR(1, 9), w4SR(1, 9), w4SR(1, 9)}
	f := map[string]string{}
	f["sample"] = fmt.Sprintf("%d,%d,%d,%d", xs[0], xs[1], xs[2], xs[3])
	switch kp {
	case 1:
		var idx []int
		var vals []int
		for {
			idx = []int{w4SR(0, 3), w4SR(0, 3), w4SR(0, 3), w4SR(0, 3)}
			vals = []int{xs[idx[0]], xs[idx[1]], xs[idx[2]], xs[idx[3]]}
			if w4SSum(vals)%4 == 0 {
				break
			}
		}
		f["indices"] = fmt.Sprintf("%d,%d,%d,%d", idx[0]+1, idx[1]+1, idx[2]+1, idx[3]+1)
		ans := w4SSum(vals) / 4
		return w4S(kp, fmt.Sprintf("A bootstrap resample of the sample %v uses indices %v. What is the resample mean? (enter a number)", xs, []int{idx[0] + 1, idx[1] + 1, idx[2] + 1, idx[3] + 1}), strconv.Itoa(ans),
			fmt.Sprintf("Take the four indexed values %v, add them and divide by 4. The mean is %d.", vals, ans), f)
	case 2:
		b := w4SR(100, 400)
		f["B"] = strconv.Itoa(b)
		idx := (b + 39) / 40
		return w4S(kp, fmt.Sprintf("For B=%d bootstrap replicates, the nearest-rank 2.5th percentile is which 1-based order statistic? (enter a number)", b), strconv.Itoa(idx),
			fmt.Sprintf("Use ceil(0.025 B) = ceil(%d/40) = %d.", b, idx), f)
	default:
		v := w4SR(1, 9)
		f["sd"] = strconv.Itoa(v)
		return w4S(kp, fmt.Sprintf("A bootstrap distribution has standard deviation %d. What is its variance? (enter a number)", v), strconv.Itoa(v*v),
			fmt.Sprintf("Variance is the square of the standard deviation: %d squared is %d.", v, v*v), f)
	}
}

func w4Permutation(kp int) generator.Problem {
	var g1, g2 []int
	for {
		g1 = []int{w4SR(1, 9), w4SR(1, 9), w4SR(1, 9)}
		g2 = []int{w4SR(1, 9), w4SR(1, 9), w4SR(1, 9)}
		if w4SSum(g1)%3 == 0 && w4SSum(g2)%3 == 0 {
			break
		}
	}
	f := map[string]string{}
	f["group1"] = fmt.Sprintf("%d,%d,%d", g1[0], g1[1], g1[2])
	f["group2"] = fmt.Sprintf("%d,%d,%d", g2[0], g2[1], g2[2])
	diff := w4SSum(g1)/3 - w4SSum(g2)/3
	switch kp {
	case 1:
		return w4S(kp, fmt.Sprintf("The observed difference in means between %v and %v is mean(group1) minus mean(group2). What is it? (enter a number)", g1, g2), strconv.Itoa(diff),
			fmt.Sprintf("The first mean is %d and the second is %d, so the difference is %d.", w4SSum(g1)/3, w4SSum(g2)/3, diff), f)
	case 2:
		n := w4SR(5, 8)
		r := w4SR(2, n-2)
		c := 1
		for i := 1; i <= r; i++ {
			c = c * (n - r + i) / i
		}
		f["n"] = strconv.Itoa(n)
		f["group_size"] = strconv.Itoa(r)
		return w4S(kp, fmt.Sprintf("With %d observations split into a group of %d and a group of %d, how many distinct label assignments are there? (enter a number)", n, r, n-r), strconv.Itoa(c),
			fmt.Sprintf("Choose the positions of the smaller group: C(%d,%d)=%d.", n, r, c), f)
	default:
		stats := []int{diff, diff + 1, -diff, diff - 1, 0, 2, -2}
		f["permutation_stats"] = fmt.Sprintf("%d,%d,%d,%d,%d,%d,%d", stats[0], stats[1], stats[2], stats[3], stats[4], stats[5], stats[6])
		count := 0
		for _, s := range stats {
			if w4Abs(s) >= w4Abs(diff) {
				count++
			}
		}
		return w4S(kp, fmt.Sprintf("The observed difference is %d. Seven permuted differences are %v. How many are at least as extreme in absolute value? (enter a number)", diff, stats), strconv.Itoa(count),
			fmt.Sprintf("Count the entries whose absolute value is at least %d. There are %d.", w4Abs(diff), count), f)
	}
}

func w4Wilcoxon(kp int) generator.Problem {
	// Paired differences with distinct absolute values avoid tie handling.
	base := []int{1, 2, 3, 4, 5}
	diffs := make([]int, 5)
	for i := 0; i < 5; i++ {
		diffs[i] = base[i]
		if rand.Intn(2) == 0 {
			diffs[i] = -diffs[i]
		}
	}
	type pair struct {
		abs, sign int
	}
	ps := make([]pair, 5)
	for i, d := range diffs {
		a := d
		if a < 0 {
			a = -a
		}
		ps[i] = pair{a, d}
	}
	sort.Slice(ps, func(i, j int) bool { return ps[i].abs < ps[j].abs })
	pos, neg := 0, 0
	for rank, p := range ps {
		if p.sign > 0 {
			pos += rank + 1
		} else {
			neg += rank + 1
		}
	}
	f := map[string]string{}
	f["differences"] = fmt.Sprintf("%d,%d,%d,%d,%d", diffs[0], diffs[1], diffs[2], diffs[3], diffs[4])
	w4SI(f, "positive_rank_sum", pos)
	w4SI(f, "negative_rank_sum", neg)
	switch kp {
	case 1:
		return w4S(kp, fmt.Sprintf("For paired differences %v, rank the absolute values and add the ranks of the positive differences. What is the sum? (enter a number)", diffs), strconv.Itoa(pos),
			fmt.Sprintf("The positive signed-rank sum is %d.", pos), f)
	case 2:
		return w4S(kp, fmt.Sprintf("For paired differences %v, what is the sum of the ranks of the negative differences? (enter a number)", diffs), strconv.Itoa(neg),
			fmt.Sprintf("The negative signed-rank sum is %d.", neg), f)
	default:
		stat := pos
		if neg < stat {
			stat = neg
		}
		return w4S(kp, fmt.Sprintf("For paired differences %v, what is the Wilcoxon signed-rank test statistic, the smaller of the two signed-rank sums? (enter a number)", diffs), strconv.Itoa(stat),
			fmt.Sprintf("The positive sum is %d and the negative sum is %d, so the smaller is %d.", pos, neg, stat), f)
	}
}

func w4MannWhitney(kp int) generator.Problem {
	n1, n2 := 3, 4
	var g1 []int
	var g2 []int
	for {
		g1 = []int{w4SR(1, 12), w4SR(1, 12), w4SR(1, 12)}
		g2 = []int{w4SR(1, 12), w4SR(1, 12), w4SR(1, 12), w4SR(1, 12)}
		seen := map[int]bool{}
		ok := true
		for _, v := range append(append([]int{}, g1...), g2...) {
			if seen[v] {
				ok = false
				break
			}
			seen[v] = true
		}
		if ok {
			break
		}
	}
	type item struct {
		v  int
		g1 bool
	}
	var pooled []item
	for _, v := range g1 {
		pooled = append(pooled, item{v, true})
	}
	for _, v := range g2 {
		pooled = append(pooled, item{v, false})
	}
	sort.Slice(pooled, func(i, j int) bool { return pooled[i].v < pooled[j].v })
	r1 := 0
	for rank, it := range pooled {
		if it.g1 {
			r1 += rank + 1
		}
	}
	u1 := r1 - n1*(n1+1)/2
	u2 := n1*n2 - u1
	f := map[string]string{}
	f["group1"] = fmt.Sprintf("%d,%d,%d", g1[0], g1[1], g1[2])
	f["group2"] = fmt.Sprintf("%d,%d,%d,%d", g2[0], g2[1], g2[2], g2[3])
	w4SI(f, "rank_sum1", r1)
	w4SI(f, "u1", u1)
	w4SI(f, "u2", u2)
	switch kp {
	case 1:
		return w4S(kp, fmt.Sprintf("Pool the two samples %v and %v and rank them from smallest to largest. What is the sum of ranks for group 1? (enter a number)", g1, g2), strconv.Itoa(r1),
			fmt.Sprintf("Group 1 has rank sum %d.", r1), f)
	case 2:
		return w4S(kp, fmt.Sprintf("For the samples %v and %v with pooled rank sum %d for group 1, what is the Mann-Whitney U for group 1? (enter a number)", g1, g2, r1), strconv.Itoa(u1),
			fmt.Sprintf("U1 = R1 - n1(n1+1)/2 = %d - 6 = %d.", r1, u1), f)
	default:
		return w4S(kp, fmt.Sprintf("For the samples %v and %v, what is U2 = n1 n2 - U1? (enter a number)", g1, g2), strconv.Itoa(u2),
			fmt.Sprintf("U2 = 3 times 4 minus %d = %d.", u1, u2), f)
	}
}

func w4KaplanMeier(kp int) generator.Problem {
	// Four observed event times, one censored observation, shifted together.
	shift := w4SR(0, 5)
	f := map[string]string{}
	f["times"] = fmt.Sprintf("%d,%d,%d,%d", shift+2, shift+4, shift+6, shift+8)
	f["censored"] = strconv.Itoa(shift + 5)
	switch kp {
	case 1:
		// S(2) = 3/4, ask the numerator.
		return w4S(kp, fmt.Sprintf("For event times %s with one censoring at %s, what is the numerator of S(%d) when the survival probability is written with denominator 4? (enter a number)", f["times"], f["censored"], shift+2), "3",
			fmt.Sprintf("At time %d, four subjects are at risk and one event occurs, so the survival probability is 3/4. The numerator is 3.", shift+2), f)
	case 2:
		return w4S(kp, fmt.Sprintf("For event times %s with one censoring at %s, how many subjects are at risk just before time %d? (enter a number)", f["times"], f["censored"], shift+4), "3",
			fmt.Sprintf("One event has occurred at time %d, so three subjects remain at risk before time %d.", shift+2, shift+4), f)
	default:
		return w4S(kp, fmt.Sprintf("For event times %s with one censoring at %s, what is the median survival time? (enter a number)", f["times"], f["censored"]), strconv.Itoa(shift+4),
			fmt.Sprintf("The survival curve first falls to 0.5 or below at time %d, so the median is %d.", shift+4, shift+4), f)
	}
}

func w4CredibleInterval(kp int) generator.Problem {
	a0, b0 := w4SR(1, 4), w4SR(1, 4)
	s, fcount := w4SR(2, 7), w4SR(2, 7)
	f := map[string]string{}
	w4SI(f, "prior_alpha", a0)
	w4SI(f, "prior_beta", b0)
	w4SI(f, "successes", s)
	w4SI(f, "failures", fcount)
	switch kp {
	case 1:
		return w4S(kp, fmt.Sprintf("A Beta(%d,%d) prior is updated with %d successes. What is the posterior alpha? (enter a number)", a0, b0, s), strconv.Itoa(a0+s),
			fmt.Sprintf("The posterior alpha is prior alpha plus successes: %d+%d=%d.", a0, s, a0+s), f)
	case 2:
		return w4S(kp, fmt.Sprintf("A Beta(%d,%d) prior is updated with %d failures. What is the posterior beta? (enter a number)", a0, b0, fcount), strconv.Itoa(b0+fcount),
			fmt.Sprintf("The posterior beta is prior beta plus failures: %d+%d=%d.", b0, fcount, b0+fcount), f)
	default:
		return w4S(kp, fmt.Sprintf("For the posterior Beta(%d,%d), what is the denominator of the posterior mean alpha/(alpha+beta)? (enter a number)", a0+s, b0+fcount), strconv.Itoa(a0+s+b0+fcount),
			fmt.Sprintf("The posterior mean has denominator %d+%d=%d.", a0+s, b0+fcount, a0+s+b0+fcount), f)
	}
}

func w4InformationCriteria(kp int) generator.Problem {
	k := w4SR(2, 6)
	dev := w4SR(10, 40)
	f := map[string]string{}
	w4SI(f, "parameters", k)
	w4SI(f, "deviance", dev)
	switch kp {
	case 1:
		return w4S(kp, fmt.Sprintf("A model has %d parameters and -2 log likelihood %d. What is AIC? (enter a number)", k, dev), strconv.Itoa(dev+2*k),
			fmt.Sprintf("AIC = -2 log L + 2k = %d + %d = %d.", dev, 2*k, dev+2*k), f)
	case 2:
		return w4S(kp, fmt.Sprintf("For a model with %d parameters, what is the AIC penalty term 2k? (enter a number)", k), strconv.Itoa(2*k),
			fmt.Sprintf("The AIC penalty is twice the number of parameters: 2 times %d = %d.", k, 2*k), f)
	default:
		a, b := dev+2*k, dev+2*(k+1)
		f["aic_small"] = strconv.Itoa(a)
		f["aic_large"] = strconv.Itoa(b)
		return w4S(kp, fmt.Sprintf("Two nested models have AIC %d and %d. Which has the better AIC? Enter 1 for the first and 2 for the second.", a, b), "1",
			fmt.Sprintf("Lower AIC is better. The first model has %d, which is lower than %d.", a, b), f)
	}
}

func w4QuantileRegression(kp int) generator.Problem {
	slope := w4SR(1, 4)
	intercept := w4SR(0, 6)
	x := w4SR(1, 5)
	f := map[string]string{}
	w4SI(f, "slope", slope)
	w4SI(f, "intercept", intercept)
	w4SI(f, "x", x)
	switch kp {
	case 1:
		return w4S(kp, fmt.Sprintf("A quantile regression line is y=%d x + %d. What is the predicted value at x=%d? (enter a number)", slope, intercept, x), strconv.Itoa(slope*x+intercept),
			fmt.Sprintf("Substitute x=%d: %d times %d plus %d = %d.", x, slope, x, intercept, slope*x+intercept), f)
	case 2:
		return w4S(kp, fmt.Sprintf("A quantile regression line is y=%d x + %d. What is its slope? (enter a number)", slope, intercept), strconv.Itoa(slope),
			fmt.Sprintf("The coefficient of x is the slope, which is %d.", slope), f)
	default:
		return w4S(kp, fmt.Sprintf("For the quantile regression line y=%d x + %d at x=%d, which quantile is estimated by minimizing the sum of absolute residuals? Enter 2 for the median and 3 for the upper quartile.", slope, intercept, x), "2",
			"The check loss with tau=1/2 is half the absolute residual, so the median is the minimizer.", f)
	}
}

func w4PoissonRegression(kp int) generator.Problem {
	a := w4SR(-2, 2)
	b := w4SR(1, 3)
	x := w4SR(1, 4)
	f := map[string]string{}
	w4SI(f, "intercept", a)
	w4SI(f, "coefficient", b)
	w4SI(f, "x", x)
	switch kp {
	case 1:
		return w4S(kp, fmt.Sprintf("A Poisson regression has log(rate)=%d+%d x. What is the linear predictor at x=%d? (enter a number)", a, b, x), strconv.Itoa(a+b*x),
			fmt.Sprintf("Substitute x=%d: %d+%d times %d=%d.", x, a, b, x, a+b*x), f)
	case 2:
		return w4S(kp, fmt.Sprintf("A Poisson regression has log(rate)=%d+%d x. At x=0 the rate is e^a. What is a? (enter a number)", a, b), strconv.Itoa(a),
			fmt.Sprintf("At x=0 the log rate is the intercept %d, so the rate is e to the %d.", a, a), f)
	default:
		return w4S(kp, fmt.Sprintf("In log(rate)=%d+%d x, is the rate increasing or decreasing in x? Enter 1 for increasing and 0 for decreasing.", a, b), "1",
			fmt.Sprintf("The coefficient %d is positive, so the log rate and the rate increase with x.", b), f)
	}
}

func w4MissingData(kp int) generator.Problem {
	switch kp {
	case 1:
		n := w4SR(3, 8)
		return w4S(kp, fmt.Sprintf("A study measures %d variables and the missingness is independent of both observed and unobserved values. Enter 1 for MCAR, 2 for MAR and 3 for MNAR.", n), "1",
			"MCAR is the mechanism in which missingness is independent of the data values.", map[string]string{"variables": strconv.Itoa(n)})
	case 2:
		m := w4SR(3, 8)
		step := 2 * w4SR(1, 3)
		vals := []int{m, m + step, m + 2*step, m + 3*step}
		mean := m + (3*step)/2
		f := map[string]string{"observed": fmt.Sprintf("%d,%d,%d,%d", vals[0], vals[1], vals[2], vals[3])}
		return w4S(kp, fmt.Sprintf("A variable has observed values %v and one missing value. What is the mean imputation? (enter a number)", vals), strconv.Itoa(mean),
			fmt.Sprintf("The observed mean is %d, so the missing value is filled with %d.", mean, mean), f)
	default:
		n := w4SR(2, 6)
		return w4S(kp, fmt.Sprintf("In a sample of %d observations, high values are preferentially missing. Does mean imputation tend to understate the population mean? Enter 1 for yes and 0 for no.", n), "1",
			"Dropping high values pulls the observed mean down, so mean imputation understates the true mean.", map[string]string{"n": strconv.Itoa(n)})
	}
}
