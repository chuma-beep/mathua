package statistics

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

// The remaining descriptive and inferential concepts each have three knowledge
// points that ask for different work: a definition or direct reading, a formula
// applied to given numbers, and an interpretation or a count over examples.
// Every generator reads KPOrdinal so the question matches the section shown.

func statFormat(v float64) string {
	s := strconv.FormatFloat(v, 'f', 2, 64)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	if s == "-0" {
		return "0"
	}
	return s
}

type measurementLevelsGen struct{}

func (g *measurementLevelsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := bridgeOrdinal(ctx)
	type item struct {
		id    string
		level int
	}
	pool := []item{
		{"eye colour", 1},
		{"country of birth", 1},
		{"telephone number", 1},
		{"satisfaction from poor to excellent", 2},
		{"education level from primary to doctoral", 2},
		{"finishing position in a race", 2},
		{"temperature in degrees Celsius", 3},
		{"calendar year", 3},
		{"score on a standardised IQ test", 3},
		{"height in centimetres", 4},
		{"annual income in dollars", 4},
		{"number of siblings", 4},
	}
	switch o {
	case 1:
		x := pool[rand.Intn(len(pool))]
		return generator.Problem{
			Question:    fmt.Sprintf("A variable records %s. Which level of measurement is it? (enter 1 nominal, 2 ordinal, 3 interval, 4 ratio)", x.id),
			Answer:      strconv.Itoa(x.level),
			Explanation: fmt.Sprintf("The variable is measured at level %d on the nominal-ordinal-interval-ratio scale.", x.level),
			Facts:       map[string]string{"kind": "level", "answer": strconv.Itoa(x.level)},
		}
	case 2:
		perm := rand.Perm(len(pool))[:6]
		parts := make([]string, 0, 6)
		count := 0
		for i, j := range perm {
			parts = append(parts, fmt.Sprintf("(%c) %s", 'a'+i, pool[j].id))
			if pool[j].level >= 2 {
				count++
			}
		}
		return generator.Problem{
			Question:    fmt.Sprintf("How many of these variables are measured at least at the ordinal level? %s (enter a number)", strings.Join(parts, " ")),
			Answer:      strconv.Itoa(count),
			Explanation: fmt.Sprintf("Exactly %d of the six variables are ordinal, interval or ratio; the rest are nominal.", count),
			Facts:       map[string]string{"kind": "at_least_ordinal", "answer": strconv.Itoa(count)},
		}
	default:
		perm := rand.Perm(len(pool))[:6]
		parts := make([]string, 0, 6)
		count := 0
		for i, j := range perm {
			parts = append(parts, fmt.Sprintf("(%c) %s", 'a'+i, pool[j].id))
			if pool[j].level == 4 {
				count++
			}
		}
		return generator.Problem{
			Question:    fmt.Sprintf("How many of these variables are measured at the ratio level? %s (enter a number)", strings.Join(parts, " ")),
			Answer:      strconv.Itoa(count),
			Explanation: fmt.Sprintf("Exactly %d of the six variables have a true zero and equal intervals, which makes them ratio variables.", count),
			Facts:       map[string]string{"kind": "ratio", "answer": strconv.Itoa(count)},
		}
	}
}

type twoSampleTGen struct{}

func (g *twoSampleTGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := bridgeOrdinal(ctx)
	pairs := [][2]int{{3, 4}, {6, 8}, {5, 12}, {8, 15}, {9, 12}, {7, 24}}
	switch o {
	case 1:
		n1 := 4 + rand.Intn(17)
		n2 := 4 + rand.Intn(17)
		df := n1 + n2 - 2
		return generator.Problem{
			Question:    fmt.Sprintf("A two-sample t test uses samples of size %d and %d. What are the pooled degrees of freedom? (enter a number)", n1, n2),
			Answer:      strconv.Itoa(df),
			Explanation: fmt.Sprintf("Pooled degrees of freedom = n1 + n2 - 2 = %d + %d - 2 = %d.", n1, n2, df),
			Facts:       map[string]string{"kind": "df", "answer": strconv.Itoa(df)},
		}
	case 2:
		p := pairs[rand.Intn(len(pairs))]
		n1 := []int{4, 9, 16, 25}[rand.Intn(4)]
		n2 := []int{4, 9, 16, 25}[rand.Intn(4)]
		s1 := p[0] * int(math.Sqrt(float64(n1)))
		s2 := p[1] * int(math.Sqrt(float64(n2)))
		se := math.Sqrt(float64(p[0]*p[0] + p[1]*p[1]))
		return generator.Problem{
			Question:    fmt.Sprintf("Two independent samples have sizes %d and %d and standard deviations %d and %d. What is the standard error of the difference in means? (enter a number)", n1, n2, s1, s2),
			Answer:      statFormat(se),
			Explanation: fmt.Sprintf("SE = sqrt(s1^2/n1 + s2^2/n2) = sqrt(%d/%d + %d/%d) = sqrt(%d + %d) = %s.", s1*s1, n1, s2*s2, n2, p[0]*p[0], p[1]*p[1], statFormat(se)),
			Facts:       map[string]string{"kind": "se", "answer": statFormat(se)},
		}
	default:
		p := pairs[rand.Intn(len(pairs))]
		se := math.Sqrt(float64(p[0]*p[0] + p[1]*p[1]))
		k := 1 + rand.Intn(3)
		diff := int(math.Round(se)) * k
		return generator.Problem{
			Question:    fmt.Sprintf("A two-sample test finds a difference in means of %d and a standard error of %s. What is the t statistic? (enter a number)", diff, statFormat(se)),
			Answer:      strconv.Itoa(k),
			Explanation: fmt.Sprintf("t = (difference in means) / SE = %d / %s = %d.", diff, statFormat(se), k),
			Facts:       map[string]string{"kind": "t", "answer": strconv.Itoa(k)},
		}
	}
}

type pairedTGen struct{}

func (g *pairedTGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := bridgeOrdinal(ctx)
	type pair struct{ before, after int }
	pairs := make([]pair, 5)
	m := -2 + rand.Intn(7)
	for i := range pairs {
		offset := []int{-2, -1, 0, 1, 2}[i]
		pairs[i] = pair{10 + 3*i + rand.Intn(3), 10 + 3*i + rand.Intn(3) + m + offset}
	}
	diffs := make([]int, len(pairs))
	sum := 0
	for i, p := range pairs {
		diffs[i] = p.after - p.before
		sum += diffs[i]
	}
	mean := float64(sum) / float64(len(diffs))
	switch o {
	case 1:
		parts := make([]string, len(pairs))
		for i, p := range pairs {
			parts[i] = fmt.Sprintf("(%d, %d)", p.before, p.after)
		}
		return generator.Problem{
			Question:    fmt.Sprintf("Five units are measured before and after a treatment as %s. What is the mean difference after minus before? (enter a number)", strings.Join(parts, ", ")),
			Answer:      statFormat(mean),
			Explanation: fmt.Sprintf("The differences are %v, whose mean is %s.", diffs, statFormat(mean)),
			Facts:       map[string]string{"kind": "mean_diff", "answer": statFormat(mean)},
		}
	case 2:
		ss := 0.0
		for _, d := range diffs {
			ss += (float64(d) - mean) * (float64(d) - mean)
		}
		return generator.Problem{
			Question:    fmt.Sprintf("The differences in a paired study are %v. What is the sum of squared deviations from their mean? (enter a number)", diffs),
			Answer:      statFormat(ss),
			Explanation: fmt.Sprintf("The mean difference is %s. Squaring each deviation and adding gives %s.", statFormat(mean), statFormat(ss)),
			Facts:       map[string]string{"kind": "ss", "answer": statFormat(ss)},
		}
	default:
		variance := 0.0
		for _, d := range diffs {
			variance += (float64(d) - mean) * (float64(d) - mean)
		}
		variance /= float64(len(diffs) - 1)
		sd := math.Sqrt(variance)
		t := mean / (sd / math.Sqrt(float64(len(diffs))))
		return generator.Problem{
			Question:    fmt.Sprintf("The differences in a paired study are %v. What is the paired t statistic? (enter a decimal to two places)", diffs),
			Answer:      statFormat(t),
			Explanation: fmt.Sprintf("The mean difference is %s and the standard deviation is %s, so t = %s / (%s / sqrt(5)) = %s.", statFormat(mean), statFormat(sd), statFormat(mean), statFormat(sd), statFormat(t)),
			Facts:       map[string]string{"kind": "t", "answer": statFormat(t)},
		}
	}
}

type proportionZGen struct{}

func (g *proportionZGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := bridgeOrdinal(ctx)
	switch o {
	case 1:
		n := []int{20, 25, 40, 50, 100, 200}[rand.Intn(6)]
		x := 1 + rand.Intn(n-1)
		phat := float64(x) / float64(n)
		return generator.Problem{
			Question:    fmt.Sprintf("A sample of %d has %d successes. What is the sample proportion? (enter a decimal)", n, x),
			Answer:      statFormat(phat),
			Explanation: fmt.Sprintf("The sample proportion is x / n = %d / %d = %s.", x, n, statFormat(phat)),
			Facts:       map[string]string{"kind": "phat", "answer": statFormat(phat)},
		}
	case 2:
		type seCase struct {
			p  float64
			n  int
			se float64
		}
		cases := []seCase{{0.5, 100, 0.05}, {0.2, 100, 0.04}, {0.8, 100, 0.04}, {0.5, 400, 0.025}, {0.2, 400, 0.02}, {0.5, 900, 1.0 / 60.0}}
		c := cases[rand.Intn(len(cases))]
		return generator.Problem{
			Question:    fmt.Sprintf("A one-proportion z test uses null proportion %.1f and sample size %d. What is the standard error under the null? (enter a decimal)", c.p, c.n),
			Answer:      statFormat(c.se),
			Explanation: fmt.Sprintf("SE = sqrt(p0(1-p0)/n) = sqrt(%.1f x %.1f / %d) = %s.", c.p, 1-c.p, c.n, statFormat(c.se)),
			Facts:       map[string]string{"kind": "se", "answer": statFormat(c.se)},
		}
	default:
		p0 := []float64{0.2, 0.25, 0.4, 0.5, 0.6, 0.75}[rand.Intn(6)]
		n := []int{100, 200, 400}[rand.Intn(3)]
		se := math.Sqrt(p0 * (1 - p0) / float64(n))
		delta := []float64{0.03, -0.04, 0.05, -0.06}[rand.Intn(4)]
		phat := p0 + delta
		z := (phat - p0) / se
		return generator.Problem{
			Question:    fmt.Sprintf("A test has null proportion %.2f, sample size %d and sample proportion %.2f. What is the z statistic? (enter a decimal to two places)", p0, n, phat),
			Answer:      statFormat(z),
			Explanation: fmt.Sprintf("SE = sqrt(%.2f x %.2f / %d) = %s, so z = (%.2f - %.2f) / %s = %s.", p0, 1-p0, n, statFormat(se), phat, p0, statFormat(se), statFormat(z)),
			Facts:       map[string]string{"kind": "z", "answer": statFormat(z)},
		}
	}
}

type residualGen struct{}

func (g *residualGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := bridgeOrdinal(ctx)
	m := 1 + rand.Intn(4)
	b := -5 + rand.Intn(11)
	switch o {
	case 1:
		x := 1 + rand.Intn(8)
		pred := m*x + b
		e := -4 + rand.Intn(9)
		obs := pred + e
		return generator.Problem{
			Question:    fmt.Sprintf("A fitted line is y = %dx + %d. An observation at x = %d has y = %d. What is its residual? (enter a number)", m, b, x, obs),
			Answer:      strconv.Itoa(e),
			Explanation: fmt.Sprintf("The predicted value is %d, so the residual is observed minus predicted: %d - %d = %d.", pred, obs, pred, e),
			Facts:       map[string]string{"kind": "residual", "answer": strconv.Itoa(e)},
		}
	case 2:
		xs := rand.Perm(8)[:3]
		es := []int{-3 + rand.Intn(7), -3 + rand.Intn(7), -3 + rand.Intn(7)}
		parts := make([]string, 3)
		sse := 0
		for i, x0 := range xs {
			x := x0 + 1
			pred := m*x + b
			parts[i] = fmt.Sprintf("(%d, %d)", x, pred+es[i])
			sse += es[i] * es[i]
		}
		return generator.Problem{
			Question:    fmt.Sprintf("A fitted line is y = %dx + %d. The observations are %s. What is the residual sum of squares? (enter a number)", m, b, strings.Join(parts, ", ")),
			Answer:      strconv.Itoa(sse),
			Explanation: fmt.Sprintf("The residuals are %v. Their squares add to %d.", es, sse),
			Facts:       map[string]string{"kind": "sse", "answer": strconv.Itoa(sse)},
		}
	default:
		xs := rand.Perm(9)[:5]
		es := []int{-3 + rand.Intn(7), -3 + rand.Intn(7), -3 + rand.Intn(7), -3 + rand.Intn(7), -3 + rand.Intn(7)}
		parts := make([]string, 5)
		positive := 0
		for i, x0 := range xs {
			x := x0 + 1
			pred := m*x + b
			parts[i] = fmt.Sprintf("(%d, %d)", x, pred+es[i])
			if es[i] > 0 {
				positive++
			}
		}
		return generator.Problem{
			Question:    fmt.Sprintf("A fitted line is y = %dx + %d. The observations are %s. How many residuals are positive? (enter a number)", m, b, strings.Join(parts, ", ")),
			Answer:      strconv.Itoa(positive),
			Explanation: fmt.Sprintf("The residuals are %v, and %d of them are positive.", es, positive),
			Facts:       map[string]string{"kind": "positive", "answer": strconv.Itoa(positive)},
		}
	}
}

type outlierGen struct{}

func (g *outlierGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := bridgeOrdinal(ctx)
	q1 := 5 + rand.Intn(16)
	iqr := 2 * (2 + rand.Intn(7))
	q3 := q1 + iqr
	switch o {
	case 1:
		return generator.Problem{
			Question:    fmt.Sprintf("A data set has first quartile %d and third quartile %d. What is the interquartile range? (enter a number)", q1, q3),
			Answer:      strconv.Itoa(iqr),
			Explanation: fmt.Sprintf("IQR = Q3 - Q1 = %d - %d = %d.", q3, q1, iqr),
			Facts:       map[string]string{"kind": "iqr", "answer": strconv.Itoa(iqr)},
		}
	case 2:
		upper := q3 + 3*iqr/2
		return generator.Problem{
			Question:    fmt.Sprintf("A data set has first quartile %d and third quartile %d. What is the upper 1.5 IQR fence? (enter a number)", q1, q3),
			Answer:      strconv.Itoa(upper),
			Explanation: fmt.Sprintf("The IQR is %d, so the upper fence is %d + 1.5 x %d = %d.", iqr, q3, iqr, upper),
			Facts:       map[string]string{"kind": "upper_fence", "answer": strconv.Itoa(upper)},
		}
	default:
		base := q1 + iqr/2
		vals := []int{q1 - iqr, q1, base, q3, q3 + iqr, q3 + 2*iqr, q3 + 3*iqr}
		rand.Shuffle(len(vals), func(i, j int) { vals[i], vals[j] = vals[j], vals[i] })
		low := q1 - 3*iqr/2
		high := q3 + 3*iqr/2
		count := 0
		for _, v := range vals {
			if v < low || v > high {
				count++
			}
		}
		return generator.Problem{
			Question:    fmt.Sprintf("A data set has first quartile %d, third quartile %d and values %v. Using the 1.5 IQR rule, how many values are outside the fences? (enter a number)", q1, q3, vals),
			Answer:      strconv.Itoa(count),
			Explanation: fmt.Sprintf("The fences are %d and %d, and %d values lie outside them.", low, high, count),
			Facts:       map[string]string{"kind": "outliers", "answer": strconv.Itoa(count)},
		}
	}
}

type normalApproxGen struct{}

func (g *normalApproxGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := bridgeOrdinal(ctx)
	type approxCase struct {
		n    int
		p    float64
		mean float64
		sd   float64
	}
	cases := []approxCase{{100, 0.5, 50, 5}, {100, 0.2, 20, 4}, {100, 0.8, 80, 4}, {400, 0.5, 200, 10}, {400, 0.2, 80, 8}, {400, 0.8, 320, 8}}
	c := cases[rand.Intn(len(cases))]
	switch o {
	case 1:
		return generator.Problem{
			Question:    fmt.Sprintf("A binomial count has %d trials and success probability %.1f. What is its mean? (enter a number)", c.n, c.p),
			Answer:      statFormat(c.mean),
			Explanation: fmt.Sprintf("Mean = np = %d x %.1f = %s.", c.n, c.p, statFormat(c.mean)),
			Facts:       map[string]string{"kind": "mean", "answer": statFormat(c.mean)},
		}
	case 2:
		return generator.Problem{
			Question:    fmt.Sprintf("A binomial count has %d trials and success probability %.1f. What is its standard deviation? (enter a number)", c.n, c.p),
			Answer:      statFormat(c.sd),
			Explanation: fmt.Sprintf("Variance = np(1-p) = %d x %.1f x %.1f = %s, so the standard deviation is %s.", c.n, c.p, 1-c.p, statFormat(c.sd*c.sd), statFormat(c.sd)),
			Facts:       map[string]string{"kind": "sd", "answer": statFormat(c.sd)},
		}
	default:
		k := int(c.mean) + rand.Intn(3) - 1
		z := (float64(k) + 0.5 - c.mean) / c.sd
		return generator.Problem{
			Question:    fmt.Sprintf("A binomial count has %d trials and success probability %.1f. Using the continuity correction, what z value approximates P(X <= %d)? (enter a decimal to two places)", c.n, c.p, k),
			Answer:      statFormat(z),
			Explanation: fmt.Sprintf("z = (k + 0.5 - np) / sd = (%d + 0.5 - %s) / %s = %s.", k, statFormat(c.mean), statFormat(c.sd), statFormat(z)),
			Facts:       map[string]string{"kind": "z", "answer": statFormat(z)},
		}
	}
}

type empiricalRuleGen struct{}

func (g *empiricalRuleGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := bridgeOrdinal(ctx)
	switch o {
	case 1:
		type ruleCase struct {
			k    int
			pct  float64
			word string
		}
		cases := []ruleCase{{1, 68, "inside"}, {2, 95, "inside"}, {3, 99.7, "inside"}, {1, 32, "outside"}, {2, 5, "outside"}, {3, 0.3, "outside"}}
		c := cases[rand.Intn(len(cases))]
		return generator.Problem{
			Question:    fmt.Sprintf("For a normal distribution, what percent of observations lie %s %d standard deviation(s) of the mean? (enter a number)", c.word, c.k),
			Answer:      statFormat(c.pct),
			Explanation: fmt.Sprintf("The empirical rule gives %s percent %s %d standard deviation(s) of the mean.", statFormat(c.pct), c.word, c.k),
			Facts:       map[string]string{"kind": "percent", "answer": statFormat(c.pct)},
		}
	case 2:
		mean := 40 + rand.Intn(41)
		sd := 2 + rand.Intn(9)
		k := 1 + rand.Intn(3)
		lower := mean - k*sd
		upper := mean + k*sd
		if rand.Intn(2) == 0 {
			return generator.Problem{
				Question:    fmt.Sprintf("A normal distribution has mean %d and standard deviation %d. What is the lower endpoint of the interval within %d standard deviation(s) of the mean? (enter a number)", mean, sd, k),
				Answer:      strconv.Itoa(lower),
				Explanation: fmt.Sprintf("The lower endpoint is mean - %d x sd = %d - %d = %d.", k, mean, k*sd, lower),
				Facts:       map[string]string{"kind": "lower", "answer": strconv.Itoa(lower)},
			}
		}
		return generator.Problem{
			Question:    fmt.Sprintf("A normal distribution has mean %d and standard deviation %d. What is the upper endpoint of the interval within %d standard deviation(s) of the mean? (enter a number)", mean, sd, k),
			Answer:      strconv.Itoa(upper),
			Explanation: fmt.Sprintf("The upper endpoint is mean + %d x sd = %d + %d = %d.", k, mean, k*sd, upper),
			Facts:       map[string]string{"kind": "upper", "answer": strconv.Itoa(upper)},
		}
	default:
		n := 1000 * (1 + rand.Intn(3))
		k := 1 + rand.Intn(3)
		inside := true
		if rand.Intn(2) == 0 {
			inside = false
		}
		pct := 68.0
		if k == 2 {
			pct = 95.0
		} else if k == 3 {
			pct = 99.7
		}
		if !inside {
			pct = 100 - pct
		}
		count := int(math.Round(float64(n) * pct / 100))
		word := "inside"
		if !inside {
			word = "outside"
		}
		return generator.Problem{
			Question:    fmt.Sprintf("A normal sample has %d observations. About how many lie %s %d standard deviation(s) of the mean? (enter a number)", n, word, k),
			Answer:      strconv.Itoa(count),
			Explanation: fmt.Sprintf("About %s percent of %d is %d observations.", statFormat(pct), n, count),
			Facts:       map[string]string{"kind": "count", "answer": strconv.Itoa(count)},
		}
	}
}

type skewGen struct{}

func (g *skewGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := bridgeOrdinal(ctx)
	type item struct {
		text string
		skew int
	}
	pool := []item{
		{"House prices, where a few mansions are far above the rest.", 1},
		{"Annual incomes, where a few earners make much more than most.", 1},
		{"Waiting times, where most waits are short and a few are very long.", 1},
		{"Exam scores with a strict ceiling, where most students score near the top.", -1},
		{"Ages at retirement, where most people retire near the same age and a few retire very early.", -1},
		{"Shoe sizes in a mixed adult sample, clustered around a central size.", 0},
		{"Heights in a large adult sample, with similar tails on both sides.", 0},
		{"Daily sales, where most days are modest and a few sale days are huge.", 1},
	}
	switch o {
	case 1:
		x := pool[rand.Intn(len(pool))]
		return generator.Problem{
			Question:    fmt.Sprintf("%s Which direction is the skew? (enter 1 for right skew, -1 for left skew, or 0 for symmetric)", x.text),
			Answer:      strconv.Itoa(x.skew),
			Explanation: fmt.Sprintf("The longer tail points %s, so the skew value is %d.", map[int]string{1: "right", -1: "left", 0: "in neither direction"}[x.skew], x.skew),
			Facts:       map[string]string{"kind": "direction", "answer": strconv.Itoa(x.skew)},
		}
	case 2:
		x := pool[rand.Intn(len(pool))]
		answer := x.skew
		return generator.Problem{
			Question:    fmt.Sprintf("%s Which is larger: the mean or the median? (enter 1 for the mean, -1 for the median, or 0 if they are about equal)", x.text),
			Answer:      strconv.Itoa(answer),
			Explanation: fmt.Sprintf("A %s distribution pulls the mean toward its longer tail, so the answer is %d.", map[int]string{1: "right-skewed", -1: "left-skewed", 0: "symmetric"}[x.skew], answer),
			Facts:       map[string]string{"kind": "mean_vs_median", "answer": strconv.Itoa(answer)},
		}
	default:
		perm := rand.Perm(len(pool))[:6]
		parts := make([]string, 6)
		count := 0
		for i, j := range perm {
			parts[i] = fmt.Sprintf("(%c) %s", 'a'+i, pool[j].text)
			if pool[j].skew == 1 {
				count++
			}
		}
		return generator.Problem{
			Question:    fmt.Sprintf("How many of these distributions are right-skewed? %s (enter a number)", strings.Join(parts, " ")),
			Answer:      strconv.Itoa(count),
			Explanation: fmt.Sprintf("Exactly %d of the six descriptions have a longer right tail.", count),
			Facts:       map[string]string{"kind": "count_right", "answer": strconv.Itoa(count)},
		}
	}
}

type regressionSlopeGen struct{}

func (g *regressionSlopeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := bridgeOrdinal(ctx)
	type seCase struct{ s, sxx, se int }
	cases := []seCase{{2, 4, 1}, {3, 9, 1}, {4, 16, 1}, {6, 4, 3}, {10, 25, 2}, {12, 36, 2}, {8, 16, 2}, {15, 9, 5}}
	c := cases[rand.Intn(len(cases))]
	switch o {
	case 1:
		return generator.Problem{
			Question:    fmt.Sprintf("A regression has residual standard deviation %d and sum of squared x deviations %d. What is the standard error of the slope? (enter a number)", c.s, c.sxx),
			Answer:      strconv.Itoa(c.se),
			Explanation: fmt.Sprintf("SE(b1) = s / sqrt(Sxx) = %d / sqrt(%d) = %d.", c.s, c.sxx, c.se),
			Facts:       map[string]string{"kind": "se", "answer": strconv.Itoa(c.se)},
		}
	case 2:
		k := 1 + rand.Intn(3)
		b1 := k * c.se
		return generator.Problem{
			Question:    fmt.Sprintf("A regression slope is %d and its standard error is %d. What is the t statistic for testing a zero slope? (enter a number)", b1, c.se),
			Answer:      strconv.Itoa(k),
			Explanation: fmt.Sprintf("t = b1 / SE(b1) = %d / %d = %d.", b1, c.se, k),
			Facts:       map[string]string{"kind": "t", "answer": strconv.Itoa(k)},
		}
	default:
		n := 5 + rand.Intn(26)
		df := n - 2
		return generator.Problem{
			Question:    fmt.Sprintf("A simple linear regression uses %d observations. How many degrees of freedom does the slope test have? (enter a number)", n),
			Answer:      strconv.Itoa(df),
			Explanation: fmt.Sprintf("A simple regression spends one degree of freedom on the intercept and one on the slope, leaving n - 2 = %d - 2 = %d.", n, df),
			Facts:       map[string]string{"kind": "df", "answer": strconv.Itoa(df)},
		}
	}
}
