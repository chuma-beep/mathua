package statistics

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/mathutil"
)

func Register(reg *generator.Registry) {
	reg.Register("stat.data.read_table", &readTableGen{})
	reg.Register("stat.data.bar_graph", &barGraphGen{})
	reg.Register("stat.data.line_plot", &linePlotGen{})
	reg.Register("stat.data.mean", &meanGen{})
	reg.Register("stat.data.median", &medianGen{})
	reg.Register("stat.data.mode", &modeGen{})
	reg.Register("stat.data.range", &rangeGen{})
	reg.Register("stat.prob.sample_space", &sampleSpaceGen{})
	reg.Register("stat.prob.basic", &basicProbGen{})
	reg.Register("stat.prob.complement", &complementProbGen{})
	reg.Register("stat.prob.compound", &compoundProbGen{})
	reg.Register("stat.prob.counting", &countingGen{})

	reg.Register("stat.adv.covariance", &covarianceGen{})
	reg.Register("stat.adv.geometric_mean", &geometricMeanGen{})
	reg.Register("stat.adv.harmonic_mean", &harmonicMeanGen{})
	reg.Register("stat.adv.rms", &rmsGen{})
	reg.Register("stat.adv.variance", &varianceGen{})
	reg.Register("stat.dist.bernoulli", &bernoulliGen{})
	reg.Register("stat.dist.beta", &betaGen{})
	reg.Register("stat.dist.binomial", &binomialDistGen{})
	reg.Register("stat.dist.chi_square", &chiSquareGen{})
	reg.Register("stat.dist.exponential", &exponentialDistGen{})
	reg.Register("stat.dist.gamma", &gammaDistGen{})
	reg.Register("stat.dist.geometric", &geometricDistGen{})
	reg.Register("stat.dist.hypergeometric", &hypergeometricGen{})
	reg.Register("stat.dist.normal", &normalDistGen{})
	reg.Register("stat.dist.poisson", &poissonGen{})
	reg.Register("stat.dist.student_t", &studentTGen{})
	reg.Register("stat.dist.uniform", &uniformDistGen{})
	reg.Register("stat.dist.z_table", &zTableGen{})
	reg.Register("stat.infer.confidence", &confidenceGen{})
	reg.Register("stat.infer.sampling", &samplingGen{})
	reg.Register("stat.prob.bayes", &bayesGen{})
	reg.Register("stat.prob.continuous_rv", &continuousRVGen{})
	reg.Register("stat.prob.discrete_rv", &discreteRVGen{})
	reg.Register("stat.prob.expected_value", &expectedValueGen{})
}

type readTableGen struct{}

func (g *readTableGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	numRows := rand.Intn(max(1, scale/2)) + 3
	numCols := rand.Intn(max(1, scale/2)) + 2
	rowLabels := make([]string, numRows)
	for i := 0; i < numRows; i++ {
		items := []string{"Apples", "Oranges", "Bananas", "Grapes", "Mangoes", "Cherries", "Pears", "Plums"}
		rowLabels[i] = items[rand.Intn(len(items))]
	}
	colLabels := make([]string, numCols)
	for i := 0; i < numCols; i++ {
		units := []string{"Week 1", "Week 2", "Week 3", "Jan", "Feb", "Mar", "Mon", "Tue", "Wed"}
		colLabels[i] = units[rand.Intn(len(units))]
	}
	table := make([][]int, numRows)
	for i := 0; i < numRows; i++ {
		table[i] = make([]int, numCols)
		for j := 0; j < numCols; j++ {
			table[i][j] = rand.Intn(max(1, scale*4)) + 1
		}
	}
	pickRow := rand.Intn(numRows)
	sum := 0
	for j := 0; j < numCols; j++ {
		sum += table[pickRow][j]
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-12s", ""))
	for j := 0; j < numCols; j++ {
		sb.WriteString(fmt.Sprintf("%-10s", colLabels[j]))
	}
	sb.WriteString("\n")
	for i := 0; i < numRows; i++ {
		sb.WriteString(fmt.Sprintf("%-12s", rowLabels[i]))
		for j := 0; j < numCols; j++ {
			sb.WriteString(fmt.Sprintf("%-10d", table[i][j]))
		}
		sb.WriteString("\n")
	}
	return generator.Problem{
		Question: fmt.Sprintf("According to the table, how many %s are there in total?\n\n%s", rowLabels[pickRow], sb.String()),
		Answer:   strconv.Itoa(sum),
		Explanation: fmt.Sprintf("Sum of %s: %s = %d", rowLabels[pickRow], strings.Join(func() []string {
			vals := make([]string, numCols)
			for j, v := range table[pickRow] {
				vals[j] = strconv.Itoa(v)
			}
			return vals
		}(), " + "), sum),
	}
}

type barGraphGen struct{}

func (g *barGraphGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	n := rand.Intn(3) + 3
	labels := make([]string, n)
	values := make([]int, n)
	items := []string{"Cats", "Dogs", "Birds", "Fish", "Hamsters", "Rabbits", "Turtles", "Snakes"}
	chosen := make(map[string]bool)
	for i := 0; i < n; i++ {
		var label string
		for {
			label = items[rand.Intn(len(items))]
			if !chosen[label] {
				chosen[label] = true
				break
			}
		}
		labels[i] = label
		values[i] = rand.Intn(max(1, scale*6)) + 5
	}
	maxIdx := 0
	for i := 1; i < n; i++ {
		if values[i] > values[maxIdx] {
			maxIdx = i
		}
	}
	dataParts := make([]string, n)
	for i := 0; i < n; i++ {
		dataParts[i] = fmt.Sprintf("%s: %d", labels[i], values[i])
	}
	return generator.Problem{
		Question:    fmt.Sprintf("In a bar graph showing pet popularity (%s), which pet has the highest value?", strings.Join(dataParts, ", ")),
		Answer:      labels[maxIdx],
		Explanation: fmt.Sprintf("%s has the highest value of %d.", labels[maxIdx], values[maxIdx]),
	}
}

type linePlotGen struct{}

func (g *linePlotGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	n := rand.Intn(3) + 4
	xs := make([]int, n)
	ys := make([]int, n)
	startX := rand.Intn(5) + 1
	startY := rand.Intn(max(1, scale*4)) + 5
	for i := 0; i < n; i++ {
		xs[i] = startX + i
		if i == 0 {
			ys[i] = startY
		} else {
			delta := rand.Intn(max(1, scale)) - 3
			ys[i] = ys[i-1] + delta
			if ys[i] < 0 {
				ys[i] = 0
			}
		}
	}
	pick := rand.Intn(n)
	points := make([]string, n)
	for i := 0; i < n; i++ {
		points[i] = fmt.Sprintf("(%d,%d)", xs[i], ys[i])
	}
	return generator.Problem{
		Question:    fmt.Sprintf("In the line plot with points %s, what is the y-value when x = %d?", strings.Join(points, ", "), xs[pick]),
		Answer:      strconv.Itoa(ys[pick]),
		Explanation: fmt.Sprintf("When x = %d, the corresponding y-value is %d.", xs[pick], ys[pick]),
	}
}

type meanGen struct{}

func (g *meanGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	n := rand.Intn(4) + 4
	vals := make([]int, n)
	sum := 0
	for i := 0; i < n; i++ {
		vals[i] = rand.Intn(max(1, scale*8)) + 1
		sum += vals[i]
	}
	mean := float64(sum) / float64(n)
	valStrs := make([]string, n)
	for i, v := range vals {
		valStrs[i] = strconv.Itoa(v)
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Find the mean of: %s", strings.Join(valStrs, ", ")),
		Answer:      fmt.Sprintf("%.1f", mean),
		Explanation: fmt.Sprintf("Mean = (%s) / %d = %d / %d = %.1f", strings.Join(valStrs, " + "), n, sum, n, mean),
	}
}

type medianGen struct{}

func (g *medianGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	n := rand.Intn(5) + 5
	vals := make([]int, n)
	for i := 0; i < n; i++ {
		vals[i] = rand.Intn(max(1, scale*10)) + 1
	}
	sorted := make([]int, n)
	copy(sorted, vals)
	sort.Ints(sorted)
	var median string
	if n%2 == 1 {
		median = strconv.Itoa(sorted[n/2])
	} else {
		mid1, mid2 := sorted[n/2-1], sorted[n/2]
		avg := float64(mid1+mid2) / 2.0
		median = fmt.Sprintf("%.1f", avg)
	}
	valStrs := make([]string, n)
	for i, v := range vals {
		valStrs[i] = strconv.Itoa(v)
	}
	sortedStrs := make([]string, n)
	for i, v := range sorted {
		sortedStrs[i] = strconv.Itoa(v)
	}
	exp := fmt.Sprintf("Sorted: %s. ", strings.Join(sortedStrs, ", "))
	if n%2 == 1 {
		exp += fmt.Sprintf("Middle value (position %d) is %s.", n/2+1, median)
	} else {
		exp += fmt.Sprintf("Average of two middle values (%d and %d) = %s.", sorted[n/2-1], sorted[n/2], median)
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Find the median of: %s", strings.Join(valStrs, ", ")),
		Answer:      median,
		Explanation: exp,
	}
}

type modeGen struct{}

func (g *modeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	n := rand.Intn(5) + 5
	base := rand.Intn(max(1, scale*4)) + 1
	modeVal := base + rand.Intn(10)
	vals := make([]int, 0, n)
	modeCount := rand.Intn(3) + 2
	for i := 0; i < modeCount; i++ {
		vals = append(vals, modeVal)
	}
	otherCount := n - modeCount
	for i := 0; i < otherCount; i++ {
		v := rand.Intn(max(1, scale*6)) + 1
		for v == modeVal {
			v = rand.Intn(max(1, scale*6)) + 1
		}
		vals = append(vals, v)
	}
	rand.Shuffle(len(vals), func(i, j int) { vals[i], vals[j] = vals[j], vals[i] })
	valStrs := make([]string, n)
	for i, v := range vals {
		valStrs[i] = strconv.Itoa(v)
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Find the mode of: %s", strings.Join(valStrs, ", ")),
		Answer:      strconv.Itoa(modeVal),
		Explanation: fmt.Sprintf("%d appears %d times, more than any other value.", modeVal, modeCount),
	}
}

type rangeGen struct{}

func (g *rangeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	n := rand.Intn(5) + 5
	vals := make([]int, n)
	minVal := rand.Intn(max(1, scale*4)) + 1
	maxVal := minVal + rand.Intn(max(1, scale*6)) + 10
	for i := 0; i < n; i++ {
		vals[i] = rand.Intn(maxVal-minVal+1) + minVal
	}
	actualMin, actualMax := vals[0], vals[0]
	for _, v := range vals {
		if v < actualMin {
			actualMin = v
		}
		if v > actualMax {
			actualMax = v
		}
	}
	r := actualMax - actualMin
	valStrs := make([]string, n)
	for i, v := range vals {
		valStrs[i] = strconv.Itoa(v)
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Find the range of: %s", strings.Join(valStrs, ", ")),
		Answer:      strconv.Itoa(r),
		Explanation: fmt.Sprintf("Range = max - min = %d - %d = %d", actualMax, actualMin, r),
	}
}

type sampleSpaceGen struct{}

func (g *sampleSpaceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scenarios := []struct {
		question string
		count    int
		expl     string
	}{
		{"rolling a standard six-sided die", 6, "A die has 6 faces: 1, 2, 3, 4, 5, 6"},
		{"flipping a coin", 2, "A coin has 2 outcomes: Heads, Tails"},
		{"rolling a die and flipping a coin", 12, "6 die outcomes x 2 coin outcomes = 12"},
		{"selecting a card from a standard 52-card deck", 52, "There are 52 unique cards in a standard deck"},
		{"spinning a spinner with 4 equal sections labeled A-D", 4, "4 sections: A, B, C, D"},
	}
	pick := rand.Intn(len(scenarios))
	return generator.Problem{
		Question:    fmt.Sprintf("List all possible outcomes when %s. How many total outcomes are there?", scenarios[pick].question),
		Answer:      strconv.Itoa(scenarios[pick].count),
		Explanation: scenarios[pick].expl,
	}
}

type basicProbGen struct{}

func (g *basicProbGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	colors := []string{"red", "blue", "green", "yellow", "orange", "purple"}
	counts := make([]int, 3)
	total := 0
	colorNames := make([]string, 3)
	for i := 0; i < 3; i++ {
		colorNames[i] = colors[rand.Intn(len(colors))]
		counts[i] = rand.Intn(max(4, scale*4)) + 1
		total += counts[i]
	}
	pick := rand.Intn(3)
	fav := counts[pick]
	gcd := mathutil.GCD(fav, total)
	bagParts := make([]string, 3)
	for i := 0; i < 3; i++ {
		bagParts[i] = fmt.Sprintf("%d %s", counts[i], colorNames[i])
	}
	return generator.Problem{
		Question:    fmt.Sprintf("A bag contains %s. What is \\(P(%s)\\)?", strings.Join(bagParts, ", "), colorNames[pick]),
		Answer:      fmt.Sprintf("%d/%d", fav/gcd, total/gcd),
		Explanation: fmt.Sprintf("\\(P(%s) = \\frac{%d}{%d} = \\frac{%d}{%d}\\)", colorNames[pick], fav, total, fav/gcd, total/gcd),
	}
}

type complementProbGen struct{}

func (g *complementProbGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	colors := []string{"red", "blue", "green", "yellow", "orange", "purple"}
	counts := make([]int, 3)
	total := 0
	colorNames := make([]string, 3)
	for i := 0; i < 3; i++ {
		colorNames[i] = colors[rand.Intn(len(colors))]
		counts[i] = rand.Intn(max(4, scale*4)) + 1
		total += counts[i]
	}
	pick := rand.Intn(3)
	other := total - counts[pick]
	gcd := mathutil.GCD(other, total)
	bagParts := make([]string, 3)
	for i := 0; i < 3; i++ {
		bagParts[i] = fmt.Sprintf("%d %s", counts[i], colorNames[i])
	}
	return generator.Problem{
		Question:    fmt.Sprintf("A bag contains %s. What is \\(P(\\text{not }%s)\\)?", strings.Join(bagParts, ", "), colorNames[pick]),
		Answer:      fmt.Sprintf("%d/%d", other/gcd, total/gcd),
		Explanation: fmt.Sprintf("\\(P(\\text{not }%s) = 1 - P(%s) = 1 - \\frac{%d}{%d} = \\frac{%d}{%d}\\)", colorNames[pick], colorNames[pick], counts[pick], total, other/gcd, total/gcd),
	}
}

type compoundProbGen struct{}

func (g *compoundProbGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	colors := []string{"red", "blue", "green", "yellow"}
	counts := make([]int, 3)
	total := 0
	colorNames := make([]string, 3)
	for i := 0; i < 3; i++ {
		colorNames[i] = colors[rand.Intn(len(colors))]
		counts[i] = rand.Intn(max(4, scale*4)) + 2
		total += counts[i]
	}
	pick1 := rand.Intn(3)
	pick2 := (pick1 + 1 + rand.Intn(2)) % 3
	c1, c2 := counts[pick1], counts[pick2]
	n1, n2 := colorNames[pick1], colorNames[pick2]
	num := c1 * c2
	den := total * (total - 1)
	gcd := mathutil.GCD(num, den)
	bagParts := make([]string, 3)
	for i := 0; i < 3; i++ {
		bagParts[i] = fmt.Sprintf("%d %s", counts[i], colorNames[i])
	}
	return generator.Problem{
		Question:    fmt.Sprintf("A bag contains %s. You draw one marble, then another without replacement. What is \\(P(%s \\text{ and then } %s)\\)?", strings.Join(bagParts, ", "), n1, n2),
		Answer:      fmt.Sprintf("%d/%d", num/gcd, den/gcd),
		Explanation: fmt.Sprintf("\\(P(%s \\text{ then } %s) = \\frac{%d}{%d} \\times \\frac{%d}{%d} = \\frac{%d}{%d} = \\frac{%d}{%d}\\)", n1, n2, c1, total, c2, total-1, num, den, num/gcd, den/gcd),
	}
}

type countingGen struct{}

func (g *countingGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	categories := []struct {
		name    string
		options []string
	}{
		{"shirts", []string{"3", "4", "5", "6", "7"}},
		{"pants", []string{"2", "3", "4", "5"}},
		{"hats", []string{"2", "3", "4"}},
		{"shoes", []string{"2", "3"}},
	}
	n := rand.Intn(2) + 2
	items := make([]struct{ name, count string }, n)
	chosen := make(map[int]bool)
	total := 1
	for i := 0; i < n; i++ {
		var idx int
		for {
			idx = rand.Intn(len(categories))
			if !chosen[idx] {
				chosen[idx] = true
				break
			}
		}
		opt := categories[idx].options[rand.Intn(len(categories[idx].options))]
		items[i] = struct{ name, count string }{categories[idx].name, opt}
		v, _ := strconv.Atoi(opt)
		total *= v
	}
	parts := make([]string, n)
	for i, item := range items {
		parts[i] = fmt.Sprintf("%s %s", item.count, item.name)
	}
	multParts := make([]string, n)
	for i, item := range items {
		multParts[i] = item.count
	}
	return generator.Problem{
		Question:    fmt.Sprintf("You have %s. How many different outfit combinations can you make?", strings.Join(parts, ", ")),
		Answer:      strconv.Itoa(total),
		Explanation: fmt.Sprintf("Total combinations = %s = %d", strings.Join(multParts, " x "), total),
	}
}

// ----- stat.adv.* (advanced measures) -----

type covarianceGen struct{}

func (g *covarianceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"Do two variables that tend to increase together have positive covariance? (yes/no)", "yes", "Positive covariance means the variables tend to increase together."},
		{"Do two variables that move in opposite directions have positive covariance? (yes/no)", "no", "Negative covariance means one tends to increase while the other decreases."},
		{"What does covariance near zero indicate?", "no linear relationship", "Zero or near-zero covariance suggests no linear relationship between variables."},
		{"Is covariance scale-dependent or scale-invariant?", "scale-dependent", "Covariance depends on the units of measurement; it is not standardized."},
		{"For a sample of size 11, what is the denominator of the sample covariance?", "10", "Sample covariance divides by (n-1) = 10, while population covariance divides by n."},
		{"Is correlation equal to covariance divided by the product of the standard deviations? (yes/no)", "yes", "Pearson correlation standardizes covariance by dividing by the product of standard deviations."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

type geometricMeanGen struct{}

func (g *geometricMeanGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := rand.Intn(max(1, scale*2)) + 1
	b := rand.Intn(max(1, scale*2)) + 1
	p := float64(a * b)
	gm := int(math.Sqrt(p))
	for gm*gm != a*b {
		a = rand.Intn(10) + 1
		b = rand.Intn(10) + 1
		p = float64(a * b)
		gm = int(math.Sqrt(p))
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Find the geometric mean of \\(%d\\) and \\(%d\\).", a, b),
		Answer:      fmt.Sprintf("%d", gm),
		Explanation: fmt.Sprintf("\\(\\text{Geometric mean} = \\sqrt{%d \\times %d} = \\sqrt{%d} = %d\\)", a, b, a*b, gm),
	}
}

type harmonicMeanGen struct{}

func (g *harmonicMeanGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := (rand.Intn(max(6, scale*4)) + 2) * 2
	b := (rand.Intn(max(6, scale*4)) + 2) * 2
	numer := 2 * a * b
	denom := a + b
	gcd := mathutil.GCD(numer, denom)
	return generator.Problem{
		Question:    fmt.Sprintf("Find the harmonic mean of \\(%d\\) and \\(%d\\).", a, b),
		Answer:      fmt.Sprintf("%d/%d", numer/gcd, denom/gcd),
		Explanation: fmt.Sprintf("\\(H = \\frac{2}{\\frac{1}{%d} + \\frac{1}{%d}} = \\frac{2}{\\frac{%d}{%d}} = \\frac{%d}{%d} = \\frac{%d}{%d}\\)", a, b, a+b, a*b, 2*a*b, a+b, numer/gcd, denom/gcd),
	}
}

type rmsGen struct{}

func (g *rmsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := rand.Intn(max(1, scale*2)) + 1
	b := rand.Intn(max(1, scale*2)) + 1
	sq := float64(a*a + b*b)
	p := int(sq / 2)
	rms := int(math.Sqrt(float64(p)))
	for rms*rms != p {
		a = rand.Intn(8) + 1
		b = a
		sq = float64(a*a + b*b)
		p = int(sq / 2)
		rms = int(math.Sqrt(float64(p)))
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Find the root mean square (RMS) of \\(%d\\) and \\(%d\\).", a, b),
		Answer:      fmt.Sprintf("%d", rms),
		Explanation: fmt.Sprintf("\\(\\text{RMS} = \\sqrt{\\frac{%d^{2}+%d^{2}}{2}} = \\sqrt{\\frac{%d+%d}{2}} = \\sqrt{%d} = %d\\)", a, b, a*a, b*b, p, rms),
	}
}

type varianceGen struct{}

func (g *varianceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	vals := make([]int, 5)
	sum := 0
	for i := range vals {
		vals[i] = (rand.Intn(max(4, scale*4)) + 1) * 2
		sum += vals[i]
	}
	mean := float64(sum) / float64(len(vals))
	var ss float64
	for _, v := range vals {
		d := float64(v) - mean
		ss += d * d
	}
	popVar := ss / float64(len(vals))
	valStrs := make([]string, len(vals))
	for i, v := range vals {
		valStrs[i] = strconv.Itoa(v)
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Find the population variance of: %s", strings.Join(valStrs, ", ")),
		Answer:      fmt.Sprintf("%.0f", popVar),
		Explanation: fmt.Sprintf("\\(\\text{Mean} = \\frac{%d}{%d} = %.0f\\). \\(\\text{Variance} = \\frac{\\Sigma(x_i - \\bar{x})^{2}}{n} = \\frac{%.0f}{%d} = %.0f\\)", sum, len(vals), mean, ss, len(vals), popVar),
	}
}

// ----- stat.dist.* (probability distributions) -----

type bernoulliGen struct{}

func (g *bernoulliGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"A Bernoulli distribution models an experiment with how many possible outcomes?", "2", "Bernoulli: exactly two outcomes (success/failure)."},
		{"What is the expected value (mean) of a Bernoulli(p) random variable?", "p", "E[X] = p for Bernoulli(p)."},
		{"What is the variance of a Bernoulli(p) random variable?", "p(1-p)", "Var[X] = p(1-p) for Bernoulli(p)."},
		{"What does X ~ Bernoulli(p) mean? Does X = 1 occur with probability p? (yes/no)", "yes", "Bernoulli random variable: P(X=1)=p, P(X=0)=1-p."},
		{"A Binomial(7,p) experiment consists of how many Bernoulli trials?", "7", "A Binomial(n,p) is the sum of n independent Bernoulli(p) variables."},
		{"For X ~ Bernoulli(0.3), what is P(X = 1) as a decimal?", "0.3", "PMF: P(X=1)=p=0.3, P(X=0)=1-p=0.7."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

type betaGen struct{}

func (g *betaGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"What is the support of the Beta distribution?", "[0, 1]", "The Beta distribution is defined on the interval [0, 1]."},
		{"What parameters define the Beta distribution?", "α (alpha) and β (beta)", "Beta(α,β) is parameterized by two shape parameters α > 0 and β > 0."},
		{"What is Beta(1,1) equivalent to?", "Uniform(0,1)", "Beta(1,1) = Uniform(0,1) (a flat prior in Bayesian statistics)."},
		{"The Beta distribution is the conjugate prior for which likelihood?", "Bernoulli/Binomial", "Beta is conjugate to Bernoulli and Binomial likelihoods."},
		{"What is the mean of Beta(α,β)?", "α/(α+β)", "Mean of Beta(α,β) = α/(α+β)."},
		{"The Beta distribution generalizes which simpler distribution?", "Uniform", "Beta(1,1) = Uniform, and other parameter choices create various shapes on [0,1]."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

type binomialDistGen struct{}

func (g *binomialDistGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"What is the expected value of Binomial(n,p)?", "np", "E[X] = np for Binomial(n,p)."},
		{"What is the variance of Binomial(n,p)?", "np(1-p)", "Var[X] = np(1-p) for Binomial(n,p)."},
		{"A Binomial distribution counts the number of ____ in n independent trials.", "successes", "Binomial: count of successes in n independent Bernoulli trials."},
		{"In the Binomial PMF, what is C(5,2), the number of ways to get 2 successes in 5 trials?", "10", "The binomial PMF is C(n,k) p^k (1-p)^(n-k); here C(5,2) = 10."},
		{"How many outcomes does a single Binomial trial have?", "2", "Each trial is Bernoulli with 2 outcomes (success/failure)."},
		{"What condition must hold between trials in a Binomial experiment?", "independence", "Binomial trials must be independent of each other."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

type chiSquareGen struct{}

func (g *chiSquareGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"What is the mean of χ²(k)?", "k", "The mean of a χ² distribution with k degrees of freedom is k."},
		{"What is the variance of χ²(k)?", "2k", "The variance of χ²(k) is 2k."},
		{"The χ² distribution is the sum of squared independent ____ random variables.", "standard normal", "If Z_i ∼ N(0,1) i.i.d., then ΣZ_i² ∼ χ²(k)."},
		{"What parameter determines the shape of a χ² distribution?", "degrees of freedom (k)", "The degrees of freedom parameter k controls the shape."},
		{"Does a chi-square random variable take negative values? (yes/no)", "no", "Chi-square takes only non-negative values (support x >= 0)."},
		{"The χ² distribution is used in which common statistical test?", "goodness of fit", "The χ² test is used for goodness of fit and independence tests."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

type exponentialDistGen struct{}

func (g *exponentialDistGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"Events occur at rate 4 per hour. What is the mean waiting time in hours as a decimal?", "0.25", "The mean of Exp(rate) is 1/rate = 1/4."},
		{"For an Exponential(rate=5) variable, what is the variance of the waiting time as a decimal?", "0.04", "The variance of Exp(rate) is 1/rate^2 = 1/25."},
		{"Does the Exponential distribution have the memoryless property? (yes/no)", "yes", "The exponential is the only continuous distribution with the memoryless property: past waiting tells nothing about the future."},
		{"What type of process is the Exponential distribution used to model?", "waiting times", "Exponential models waiting times between Poisson events."},
		{"The mean waiting time is 0.5 hours. What is the rate parameter per hour?", "2", "The rate is the average number of events per unit time; rate = 1/mean = 2."},
		{"The Exponential distribution is a special case of which distribution?", "Gamma", "Exponential(rate) = Gamma(shape=1, rate)."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

type gammaDistGen struct{}

func (g *gammaDistGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"How many parameters define the Gamma distribution?", "2", "Gamma is defined by shape k and a rate (or scale) parameter."},
		{"If X has Gamma(shape=3, scale=2), what is the mean of X?", "6", "The mean of Gamma is shape times scale = 3 times 2."},
		{"Is Gamma(shape=1) the same as an Exponential distribution? (yes/no)", "yes", "Gamma(shape=1, rate) = Exponential(rate)."},
		{"Can a Gamma random variable take negative values? (yes/no)", "no", "The Gamma distribution is defined only for positive real numbers."},
		{"The sum of 5 independent Exponential(rate=2) variables follows a Gamma distribution with what shape?", "5", "The sum of k independent Exp(rate) variables is Gamma(shape=k, rate)."},
		{"Is chi-square with k degrees of freedom a special case of the Gamma distribution? (yes/no)", "yes", "Chi-square(k) = Gamma(shape=k/2, scale=2)."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

type geometricDistGen struct{}

func (g *geometricDistGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"What does a Geometric(p) distribution model?", "number of trials until first success", "Geometric: the number of Bernoulli trials needed to get the first success."},
		{"What is the expected value of Geometric(p)?", "1/p", "E[X] = 1/p for Geometric(p)."},
		{"For a Geometric(p=0.5) variable, what is the variance?", "2", "Var[X] = (1-p)/p^2 = 0.5/0.25 for p = 0.5."},
		{"What is P(X=k) for Geometric(p)?", "(1-p)^(k-1) p", "Geometric PMF: P(X=k) = (1-p)^(k-1)p for k = 1, 2, ..."},
		{"Does the Geometric distribution have the memoryless property?", "yes", "Like the exponential (its continuous analogue), the geometric is memoryless."},
		{"Does Geometric count trials to the first success while Binomial counts successes in n trials? (yes/no)", "yes", "Both involve Bernoulli trials but count different things."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

type hypergeometricGen struct{}

func (g *hypergeometricGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"How does Hypergeometric differ from Binomial?", "sampling without replacement", "Hypergeometric: sampling without replacement; Binomial: with replacement."},
		{"How many parameters define the Hypergeometric distribution? (enter a number)", "3", "Hypergeometric(N, K, n): population size N, K successes, n draws: three parameters."},
		{"Hypergeometric(N=10,K=4,n=5): what is the expected value? (enter a number)", "2", "E[X] = nK/N = 5·4/10 = 2 (n × proportion of successes)."},
		{"In Hypergeometric, draws are ____.", "dependent (without replacement)", "Each draw changes the composition of the remaining population."},
		{"Is the support of Hypergeometric finite? (yes/no)", "yes", "The number of successes drawn cannot exceed K or n: finitely many values."},
		{"As N → ∞ with K/N fixed, Hypergeometric approaches which distribution?", "Binomial", "When N is large, sampling without replacement approximates Binomial."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

type normalDistGen struct{}

func (g *normalDistGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"How many parameters define the Normal distribution?", "2", "N(mean, variance) is parameterized by its mean and its variance."},
		{"About what proportion of Normal data falls within 2 SD of the mean? Answer as a decimal.", "0.95", "About 0.68 within 1 SD, 0.95 within 2 SD, 0.997 within 3 SD."},
		{"What is the mean of a standard Normal distribution?", "0", "Standard Normal: N(0,1) with μ=0, σ=1."},
		{"What is the variance of a standard Normal distribution?", "1", "Standard Normal: N(0,1) with μ=0, σ²=1."},
		{"If X = 70 comes from a Normal population with mean 50 and SD 10, what is the z-score?", "2", "Standardization: z = (x - mean)/SD = 30/10."},
		{"The Normal distribution is symmetric about its ____.", "mean", "The Normal distribution is symmetric about μ."},
		{"What is the shape of the Normal distribution?", "bell curve", "The Normal distribution has the characteristic bell curve shape."},
		{"Do sample means approach a Normal distribution as the sample size grows? (yes/no)", "yes", "CLT: sample means of independent identical variables approach Normal as n grows."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

type poissonGen struct{}

func (g *poissonGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"A Poisson process averages 3 events per hour. How many events on average in 2 hours?", "6", "Mean = rate times time = 3 times 2."},
		{"For X ~ Poisson(4), what is the variance?", "4", "For Poisson, variance = mean."},
		{"Does a Poisson model count events in a fixed interval? (yes/no)", "yes", "Poisson models the number of events occurring in a fixed time or space."},
		{"Is P(X=k) = e^(-L) L^k / k! the Poisson PMF? (yes/no)", "yes", "Poisson PMF: P(X=k) = e^(-L)L^k/k! for k = 0, 1, 2, ..."},
		{"Is 3 a possible value of a Poisson random variable? (yes/no)", "yes", "Poisson is defined for non-negative integers 0, 1, 2, ..."},
		{"What is the relationship between Poisson and Exponential?", "inter-arrival times are Exponential", "In a Poisson process, the time between events follows Exponential(rate)."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

type studentTGen struct{}

func (g *studentTGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"Is Student's t defined by its degrees of freedom? (yes/no)", "yes", "t(ν) is parameterized by degrees of freedom ν."},
		{"As ν → ∞, the t distribution approaches ____.", "standard Normal", "As degrees of freedom increase, t approaches N(0,1)."},
		{"Compared to Normal, does Student's t have heavier tails? (yes/no)", "yes", "The t distribution has heavier tails than Normal, especially for small ν."},
		{"What is the mean of t(ν) for ν > 1?", "0", "The t distribution is symmetric about 0 with mean 0 (for ν > 1)."},
		{"Is the t distribution used instead of Normal when sigma is unknown? (yes/no)", "yes", "t is used when the population standard deviation is unknown."},
		{"The t(nu) variance is finite when nu exceeds which integer? (enter a number)", "2", "Var[t(ν)] = ν/(ν-2) for ν > 2."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

type uniformDistGen struct{}

func (g *uniformDistGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"What is the mean of Uniform(a,b)?", "(a+b)/2", "E[X] = (a+b)/2 for Uniform(a,b)."},
		{"For Uniform(0,12), what is the variance?", "12", "Var = (b-a)^2/12 = 144/12."},
		{"For Uniform(2,6), what is the density value as a decimal?", "0.25", "Density = 1/(b-a) = 1/4 on the interval, 0 otherwise."},
		{"Is Uniform(0,1) the same as Beta(1,1)? (yes/no)", "yes", "Beta(1,1) = Uniform(0,1)."},
		{"In a Uniform distribution, all outcomes are ____.", "equally likely", "Uniform: constant probability density over the interval."},
		{"For Uniform(3,9), what is the length of the support interval?", "6", "Uniform is defined on the interval from a to b, here length 9-3."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

type zTableGen struct{}

func (g *zTableGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"What z-score corresponds to the 50th percentile of N(0,1)?", "0", "The median of standard Normal is 0."},
		{"What z-score approximately corresponds to the 97.5th percentile?", "1.96", "P(Z < 1.96) ≈ 0.975, used for 95% confidence intervals."},
		{"What is the z-score for a 95% confidence interval?", "1.96", "95% CI uses z = 1.96 (approximately 2)."},
		{"Does a z-score measure how many SD a value is from the mean? (yes/no)", "yes", "z = (x - mean)/SD measures standard deviations from the mean."},
		{"About what proportion of data falls within z = +-1? Answer as a decimal.", "0.68", "About 0.68 of data falls within 1 standard deviation of the mean."},
		{"About what proportion of data falls within z = +-2? Answer as a decimal.", "0.95", "About 0.95 of data falls within 2 standard deviations of the mean."},
		{"About what proportion of data falls within z = +-3? Answer as a decimal.", "0.997", "About 0.997 of data falls within 3 standard deviations of the mean."},
		{"A z-score of -2 means the value is ____.", "2 SD below the mean", "Negative z: below mean; positive z: above mean."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

// ----- stat.infer.* (inference) -----

type confidenceGen struct{}

func (g *confidenceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"Does a 95% confidence interval mean that 95% of such intervals contain the true parameter? (yes/no)", "yes", "In repeated sampling, 95% of CIs contain the population parameter."},
		{"A sample of n=100 has SD 20. With z=2, what is the margin of error?", "4", "CI = sample mean plus/minus margin; margin = z times SD/sqrt(n) = 2 times 2."},
		{"Does a larger sample give a narrower confidence interval? (yes/no)", "yes", "CI width goes like 1/sqrt(n), so larger samples give more precise estimates."},
		{"What happens to CI width when confidence level increases from 95% to 99%?", "it gets wider", "Higher confidence requires a wider interval."},
		{"A 95% CI is 10 plus/minus 2. What is the margin of error?", "2", "Margin of error = critical value times standard error."},
		{"For 50 of 100 voters with z=2, what is the margin of error as a decimal?", "0.1", "Margin = z times sqrt(phat(1-phat)/n) = 2 times 0.05."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

type samplingGen struct{}

func (g *samplingGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"The population SD is 12 and n=36. What is the standard error of the mean?", "2", "Standard error = population SD / sqrt(sample size) = 12/6."},
		{"Does the sampling distribution of the mean approach Normal as n grows? (yes/no)", "yes", "CLT: the sampling distribution of the mean approaches Normal for large n."},
		{"The population mean is 50. What is the mean of the sampling distribution of xbar?", "50", "xbar is an unbiased estimator of the population mean."},
		{"Is a sampling distribution the distribution of a statistic over all possible samples? (yes/no)", "yes", "The sampling distribution shows how a statistic varies from sample to sample."},
		{"Does a larger sample size reduce bias? (yes/no)", "no", "Bias depends on the sampling method, not sample size; bias is independent of n."},
		{"Is a sample a subset of the population? (yes/no)", "yes", "The population is the entire group; the sample is a subset drawn from it."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

// ----- stat.prob.* (additional probability) -----

type bayesGen struct{}

func (g *bayesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"Is P(A|B) = P(B|A)P(A)/P(B) Bayes theorem? (yes/no)", "yes", "Bayes theorem relates conditional probabilities."},
		{"In Bayes' theorem, P(A) is called the ____ probability.", "prior", "P(A) is the prior probability before observing B."},
		{"In Bayes' theorem, P(A|B) is called the ____ probability.", "posterior", "P(A|B) is the updated probability after observing B."},
		{"What is the denominator in Bayes' theorem called?", "marginal likelihood", "P(B) normalizes the posterior; it is called the marginal likelihood or evidence."},
		{"In medical testing, P(test positive | has disease) is called ____.", "sensitivity", "Sensitivity = P(positive|disease), the true positive rate."},
		{"In medical testing, P(test negative | no disease) is called ____.", "specificity", "Specificity = P(negative|no disease), the true negative rate."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

type continuousRVGen struct{}

func (g *continuousRVGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"What is the probability of a continuous random variable taking any single value?", "0", "P(X = a) = 0 for continuous random variables."},
		{"What function gives probabilities for a continuous random variable via areas?", "probability density function (PDF)", "P(a ≤ X ≤ b) = ∫ₐᵇ f(x)dx."},
		{"Is the CDF equal to P(X <= x), the probability of being less than or equal to x? (yes/no)", "yes", "The CDF gives the probability of being less than or equal to x."},
		{"Is the CDF the integral of the PDF? (yes/no)", "yes", "CDF(x) is the integral of the PDF up to x, and PDF is the derivative of the CDF."},
		{"For a valid PDF, what must ∫_{-∞}^{∞} f(x)dx equal?", "1", "Total probability under the PDF must be 1."},
		{"Is the PDF of a continuous variable always ≤ 1?", "no", "The PDF can be > 1 as long as the total area is 1."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

type discreteRVGen struct{}

func (g *discreteRVGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"Does the PMF give P(X = x) for each possible x? (yes/no)", "yes", "PMF gives the probability of each specific value for a discrete variable."},
		{"For a valid PMF, what must Σ P(X=x) equal?", "1", "All probabilities must sum to 1."},
		{"Is the PMF for discrete variables and the PDF for continuous ones? (yes/no)", "yes", "PMF gives probabilities at specific points; PDF gives density."},
		{"Is the discrete CDF the sum of PMF values up to x? (yes/no)", "yes", "Discrete CDF sums PMF values up to x."},
		{"Can a discrete random variable take infinitely many values?", "yes", "Discrete variables can be countably infinite (e.g., Poisson, Geometric)."},
		{"X takes values 1 and 3 with probability 0.5 each. What is E[X]?", "2", "The mean is the probability-weighted sum: 1 times 0.5 plus 3 times 0.5."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

type expectedValueGen struct{}

func (g *expectedValueGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"What is the expected value of a constant c?", "c", "E[c] = c for any constant c."},
		{"If E[X] = 5, what is E[2X + 1]?", "11", "Linearity of expectation: E[aX+b] = aE[X] + b = 2 times 5 plus 1."},
		{"If E[X] = 3 and E[Y] = 4, what is E[X + Y]?", "7", "Linearity: expectation of sum = sum of expectations."},
		{"Is E[XY] always equal to E[X]E[Y]?", "no", "E[XY] = E[X]E[Y] only when X and Y are independent."},
		{"If E[X] = 2 and E[X^2] = 8, what is Var[X]?", "4", "Variance = second moment minus squared mean = 8 - 4."},
		{"What is E[X²] called?", "second moment", "E[X²] is the second moment about the origin."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}
