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

func (g *readTableGen) Generate(difficulty float64) generator.Problem {
	scale := int(1 + difficulty*5)
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

func (g *barGraphGen) Generate(difficulty float64) generator.Problem {
	scale := int(1 + difficulty*5)
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

func (g *linePlotGen) Generate(difficulty float64) generator.Problem {
	scale := int(1 + difficulty*5)
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

func (g *meanGen) Generate(difficulty float64) generator.Problem {
	scale := int(1 + difficulty*5)
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

func (g *medianGen) Generate(difficulty float64) generator.Problem {
	scale := int(1 + difficulty*5)
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

func (g *modeGen) Generate(difficulty float64) generator.Problem {
	scale := int(1 + difficulty*5)
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

func (g *rangeGen) Generate(difficulty float64) generator.Problem {
	scale := int(1 + difficulty*5)
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

func (g *sampleSpaceGen) Generate(difficulty float64) generator.Problem {
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

func (g *basicProbGen) Generate(difficulty float64) generator.Problem {
	scale := int(1 + difficulty*5)
	colors := []string{"red", "blue", "green", "yellow", "orange", "purple"}
	counts := make([]int, 3)
	total := 0
	colorNames := make([]string, 3)
	for i := 0; i < 3; i++ {
		colorNames[i] = colors[rand.Intn(len(colors))]
		counts[i] = rand.Intn(max(1, scale)) + 1
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

func (g *complementProbGen) Generate(difficulty float64) generator.Problem {
	scale := int(1 + difficulty*5)
	colors := []string{"red", "blue", "green", "yellow", "orange", "purple"}
	counts := make([]int, 3)
	total := 0
	colorNames := make([]string, 3)
	for i := 0; i < 3; i++ {
		colorNames[i] = colors[rand.Intn(len(colors))]
		counts[i] = rand.Intn(max(1, scale)) + 1
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

func (g *compoundProbGen) Generate(difficulty float64) generator.Problem {
	scale := int(1 + difficulty*5)
	colors := []string{"red", "blue", "green", "yellow"}
	counts := make([]int, 3)
	total := 0
	colorNames := make([]string, 3)
	for i := 0; i < 3; i++ {
		colorNames[i] = colors[rand.Intn(len(colors))]
		counts[i] = rand.Intn(max(1, scale)) + 2
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

func (g *countingGen) Generate(difficulty float64) generator.Problem {
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

func (g *covarianceGen) Generate(difficulty float64) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"What does a positive covariance between two variables indicate?", "they move in the same direction", "Positive covariance means the variables tend to increase together."},
		{"What does a negative covariance between two variables indicate?", "they move in opposite directions", "Negative covariance means one tends to increase while the other decreases."},
		{"What does covariance near zero indicate?", "no linear relationship", "Zero or near-zero covariance suggests no linear relationship between variables."},
		{"Is covariance scale-dependent or scale-invariant?", "scale-dependent", "Covariance depends on the units of measurement; it is not standardized."},
		{"What is the formula for sample covariance?", "Σ(x_i-x̄)(y_i-ȳ)/(n-1)", "Sample covariance divides by (n-1), while population covariance divides by n."},
		{"How is correlation related to covariance?", "correlation = covariance/(σ_x·σ_y)", "Pearson correlation standardizes covariance by dividing by the product of standard deviations."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

type geometricMeanGen struct{}

func (g *geometricMeanGen) Generate(difficulty float64) generator.Problem {
	scale := int(1 + difficulty*5)
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

func (g *harmonicMeanGen) Generate(difficulty float64) generator.Problem {
	scale := int(1 + difficulty*5)
	a := (rand.Intn(max(1, scale)) + 2) * 2
	b := (rand.Intn(max(1, scale)) + 2) * 2
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

func (g *rmsGen) Generate(difficulty float64) generator.Problem {
	scale := int(1 + difficulty*5)
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

func (g *varianceGen) Generate(difficulty float64) generator.Problem {
	scale := int(1 + difficulty*5)
	vals := make([]int, 5)
	sum := 0
	for i := range vals {
		vals[i] = (rand.Intn(max(1, scale))+1)*2
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

func (g *bernoulliGen) Generate(difficulty float64) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"A Bernoulli distribution models an experiment with how many possible outcomes?", "2", "Bernoulli: exactly two outcomes (success/failure)."},
		{"What is the expected value (mean) of a Bernoulli(p) random variable?", "p", "E[X] = p for Bernoulli(p)."},
		{"What is the variance of a Bernoulli(p) random variable?", "p(1-p)", "Var[X] = p(1-p) for Bernoulli(p)."},
		{"What does X ~ Bernoulli(p) mean?", "X takes value 1 with probability p and 0 with probability 1-p", "Bernoulli random variable: P(X=1)=p, P(X=0)=1-p."},
		{"How is a Bernoulli trial different from a Binomial experiment?", "Bernoulli is one trial; Binomial is n trials", "A Binomial(n,p) is the sum of n independent Bernoulli(p) variables."},
		{"What is the probability mass function of Bernoulli(p)?", "P(X=x) = p^x(1-p)^(1-x)", "PMF: P(X=1)=p, P(X=0)=1-p."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

type betaGen struct{}

func (g *betaGen) Generate(difficulty float64) generator.Problem {
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

func (g *binomialDistGen) Generate(difficulty float64) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"What is the expected value of Binomial(n,p)?", "np", "E[X] = np for Binomial(n,p)."},
		{"What is the variance of Binomial(n,p)?", "np(1-p)", "Var[X] = np(1-p) for Binomial(n,p)."},
		{"A Binomial distribution counts the number of ____ in n independent trials.", "successes", "Binomial: count of successes in n independent Bernoulli trials."},
		{"What is P(X=k) for Binomial(n,p)?", "C(n,k) p^k (1-p)^(n-k)", "The binomial probability mass function."},
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

func (g *chiSquareGen) Generate(difficulty float64) generator.Problem {
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
		{"What is the support of the χ² distribution?", "[0, ∞)", "χ² takes only non-negative values."},
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

func (g *exponentialDistGen) Generate(difficulty float64) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"What is the mean of Exponential(λ)?", "1/λ", "The mean of Exp(λ) is 1/λ."},
		{"What is the variance of Exponential(λ)?", "1/λ²", "The variance of Exp(λ) is 1/λ²."},
		{"What is the memoryless property of the Exponential distribution?", "P(X > s+t | X > s) = P(X > t)", "The exponential is the only continuous distribution with the memoryless property."},
		{"What type of process is the Exponential distribution used to model?", "waiting times", "Exponential models waiting times between Poisson events."},
		{"What is the rate parameter λ in an Exponential distribution?", "average number of events per unit time", "λ is the rate parameter; the mean time between events is 1/λ."},
		{"The Exponential distribution is a special case of which distribution?", "Gamma", "Exponential(λ) = Gamma(1, λ)."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

type gammaDistGen struct{}

func (g *gammaDistGen) Generate(difficulty float64) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"What parameters define the Gamma distribution?", "shape (k) and rate (θ or λ)", "Gamma(k, θ) where k is shape and θ is scale (or rate)."},
		{"What is the mean of Gamma(k, θ)?", "kθ", "The mean of Gamma(k, θ) is kθ (shape × scale)."},
		{"Gamma(1, λ) is equivalent to which distribution?", "Exponential(λ)", "Gamma(1, λ) = Exponential(λ)."},
		{"What is the support of the Gamma distribution?", "(0, ∞)", "The Gamma distribution is defined for positive real numbers."},
		{"The sum of k independent Exponential(λ) variables follows what distribution?", "Gamma(k, 1/λ)", "The sum of k i.i.d. Exp(λ) variables is Gamma(k, 1/λ)."},
		{"The χ²(k) distribution is a special case of Gamma with which parameters?", "Gamma(k/2, 2)", "χ²(k) = Gamma(k/2, 2)."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

type geometricDistGen struct{}

func (g *geometricDistGen) Generate(difficulty float64) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"What does a Geometric(p) distribution model?", "number of trials until first success", "Geometric: the number of Bernoulli trials needed to get the first success."},
		{"What is the expected value of Geometric(p)?", "1/p", "E[X] = 1/p for Geometric(p)."},
		{"What is the variance of Geometric(p)?", "(1-p)/p²", "Var[X] = (1-p)/p² for Geometric(p)."},
		{"What is P(X=k) for Geometric(p)?", "(1-p)^(k-1) p", "Geometric PMF: P(X=k) = (1-p)^(k-1)p for k = 1, 2, ..."},
		{"Does the Geometric distribution have the memoryless property?", "yes", "Like the exponential (its continuous analogue), the geometric is memoryless."},
		{"What is the relationship between Geometric and Binomial?", "Geometric = first success; Binomial = count of successes in n trials", "Both involve Bernoulli trials but count different things."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

type hypergeometricGen struct{}

func (g *hypergeometricGen) Generate(difficulty float64) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"How does Hypergeometric differ from Binomial?", "sampling without replacement", "Hypergeometric: sampling without replacement; Binomial: with replacement."},
		{"What parameters define the Hypergeometric distribution?", "N, K, n", "Hypergeometric(N, K, n): population size N, K successes, n draws."},
		{"What is the expected value of Hypergeometric(N,K,n)?", "n(K/N)", "E[X] = nK/N (n × proportion of successes)."},
		{"In Hypergeometric, draws are ____.", "dependent (without replacement)", "Each draw changes the composition of the remaining population."},
		{"What is the support of Hypergeometric?", "max(0, n-(N-K)) to min(n, K)", "The number of successes drawn cannot exceed K or n."},
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

func (g *normalDistGen) Generate(difficulty float64) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"What parameters define the Normal distribution?", "μ (mean) and σ² (variance)", "N(μ, σ²) is parameterized by mean μ and variance σ²."},
		{"What is the 68-95-99.7 rule for Normal distributions?", "68% within 1σ, 95% within 2σ, 99.7% within 3σ", "Empirical rule for the standard normal."},
		{"What is the mean of a standard Normal distribution?", "0", "Standard Normal: N(0,1) with μ=0, σ=1."},
		{"What is the variance of a standard Normal distribution?", "1", "Standard Normal: N(0,1) with μ=0, σ²=1."},
		{"What transformation converts X ∼ N(μ,σ²) to standard Normal?", "Z = (X-μ)/σ", "Standardization: subtract mean, divide by standard deviation."},
		{"The Normal distribution is symmetric about its ____.", "mean", "The Normal distribution is symmetric about μ."},
		{"What is the shape of the Normal distribution?", "bell-shaped", "The Normal distribution has the characteristic bell curve shape."},
		{"What does the Central Limit Theorem say about sample means?", "they approach Normal as n increases", "CLT: sample means of i.i.d. variables approach Normal as n → ∞."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

type poissonGen struct{}

func (g *poissonGen) Generate(difficulty float64) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"What is the mean of Poisson(λ)?", "λ", "E[X] = λ for Poisson(λ)."},
		{"What is the variance of Poisson(λ)?", "λ", "Var[X] = λ for Poisson(λ) (mean = variance)."},
		{"What type of process does Poisson model?", "count of events in a fixed interval", "Poisson models the number of events occurring in a fixed time or space."},
		{"What is P(X=k) for Poisson(λ)?", "e^(-λ) λ^k / k!", "Poisson PMF: P(X=k) = e^(-λ)λ^k/k! for k = 0, 1, 2, ..."},
		{"What is the support of Poisson?", "{0, 1, 2, ...}", "Poisson is defined for non-negative integers."},
		{"What is the relationship between Poisson and Exponential?", "inter-arrival times are Exponential", "In a Poisson process, the time between events follows Exponential(λ)."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

type studentTGen struct{}

func (g *studentTGen) Generate(difficulty float64) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"What parameter defines Student's t distribution?", "degrees of freedom (ν)", "t(ν) is parameterized by degrees of freedom ν."},
		{"As ν → ∞, the t distribution approaches ____.", "standard Normal", "As degrees of freedom increase, t approaches N(0,1)."},
		{"Compared to Normal, Student's t has ____ tails.", "heavier (or fatter)", "The t distribution has heavier tails than Normal, especially for small ν."},
		{"What is the mean of t(ν) for ν > 1?", "0", "The t distribution is symmetric about 0 with mean 0 (for ν > 1)."},
		{"When is the t distribution used instead of Normal?", "when σ is unknown and estimated", "t is used when the population standard deviation is unknown."},
		{"What is the variance of t(ν) for ν > 2?", "ν/(ν-2)", "Var[t(ν)] = ν/(ν-2) for ν > 2."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

type uniformDistGen struct{}

func (g *uniformDistGen) Generate(difficulty float64) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"What is the mean of Uniform(a,b)?", "(a+b)/2", "E[X] = (a+b)/2 for Uniform(a,b)."},
		{"What is the variance of Uniform(a,b)?", "(b-a)²/12", "Var[X] = (b-a)²/12 for Uniform(a,b)."},
		{"What is the probability density function of Uniform(a,b)?", "1/(b-a) for x in [a,b]", "f(x) = 1/(b-a) for a ≤ x ≤ b, 0 otherwise."},
		{"Uniform(0,1) is a special case of which other distribution?", "Beta(1,1)", "Beta(1,1) = Uniform(0,1)."},
		{"In a Uniform distribution, all outcomes are ____.", "equally likely", "Uniform: constant probability density over the interval."},
		{"What is the support of Uniform(a,b)?", "[a, b]", "Uniform is defined on the interval [a, b]."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

type zTableGen struct{}

func (g *zTableGen) Generate(difficulty float64) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"What z-score corresponds to the 50th percentile of N(0,1)?", "0", "The median of standard Normal is 0."},
		{"What z-score approximately corresponds to the 97.5th percentile?", "1.96", "P(Z < 1.96) ≈ 0.975, used for 95% confidence intervals."},
		{"What is the z-score for a 95% confidence interval?", "1.96", "95% CI uses z = 1.96 (approximately 2)."},
		{"What does a z-score measure?", "how many SD from the mean", "z = (x - μ)/σ measures standard deviations from the mean."},
		{"About what percentage of data falls within z = ±1?", "68%", "≈68% of data falls within 1 standard deviation of the mean."},
		{"About what percentage of data falls within z = ±2?", "95%", "≈95% of data falls within 2 standard deviations of the mean."},
		{"About what percentage of data falls within z = ±3?", "99.7%", "≈99.7% of data falls within 3 standard deviations of the mean."},
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

func (g *confidenceGen) Generate(difficulty float64) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"What does a 95% confidence interval mean?", "95% of intervals contain the true parameter", "In repeated sampling, 95% of CIs contain the population parameter."},
		{"What is the formula for a confidence interval for the mean?", "x̄ ± z*(σ/√n)", "CI = sample mean ± margin of error."},
		{"How does sample size affect confidence interval width?", "larger n → narrower interval", "CI width ∝ 1/√n, so larger samples give more precise estimates."},
		{"What happens to CI width when confidence level increases from 95% to 99%?", "it gets wider", "Higher confidence requires a wider interval."},
		{"What is the margin of error in a confidence interval?", "z*(σ/√n)", "Margin of error = critical value × standard error."},
		{"A confidence interval for a proportion uses what formula?", "p̂ ± z*√(p̂(1-p̂)/n)", "CI for proportion p with sample proportion p̂."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

type samplingGen struct{}

func (g *samplingGen) Generate(difficulty float64) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"What is the standard error of the mean?", "σ/√n", "Standard error = population SD / √(sample size)."},
		{"What does the Central Limit Theorem say about the sampling distribution?", "it approaches Normal as n increases", "CLT: the sampling distribution of the mean approaches Normal for large n."},
		{"What is the mean of the sampling distribution of x̄?", "μ (population mean)", "x̄ is an unbiased estimator of μ."},
		{"What is a sampling distribution?", "distribution of a statistic over all possible samples", "The sampling distribution shows how a statistic varies from sample to sample."},
		{"Does a larger sample size reduce bias?", "no (bias is independent of n)", "Bias depends on the sampling method, not sample size."},
		{"What is the difference between population and sample?", "population = all; sample = subset", "The population is the entire group; the sample is a subset drawn from it."},
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

func (g *bayesGen) Generate(difficulty float64) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"What is Bayes' theorem?", "P(A|B) = P(B|A)P(A)/P(B)", "Bayes' theorem relates conditional probabilities."},
		{"In Bayes' theorem, P(A) is called the ____ probability.", "prior", "P(A) is the prior probability before observing B."},
		{"In Bayes' theorem, P(A|B) is called the ____ probability.", "posterior", "P(A|B) is the updated probability after observing B."},
		{"What is the denominator in Bayes' theorem called?", "marginal likelihood (or evidence)", "P(B) normalizes the posterior."},
		{"In medical testing, P(test positive | has disease) is called ____.", "sensitivity (or true positive rate)", "Sensitivity = P(+|disease)."},
		{"In medical testing, P(test negative | no disease) is called ____.", "specificity (or true negative rate)", "Specificity = P(-|no disease)."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

type continuousRVGen struct{}

func (g *continuousRVGen) Generate(difficulty float64) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"What is the probability of a continuous random variable taking any single value?", "0", "P(X = a) = 0 for continuous random variables."},
		{"What function gives probabilities for a continuous random variable via areas?", "probability density function (PDF)", "P(a ≤ X ≤ b) = ∫ₐᵇ f(x)dx."},
		{"What is the cumulative distribution function (CDF)?", "F(x) = P(X ≤ x)", "The CDF gives the probability of being less than or equal to x."},
		{"What is the relationship between PDF and CDF?", "CDF is the integral of PDF", "F(x) = ∫_{-∞}^{x} f(t)dt, and f(x) = F'(x)."},
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

func (g *discreteRVGen) Generate(difficulty float64) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"What is the probability mass function (PMF)?", "P(X = x) for each possible x", "PMF gives the probability of each specific value for a discrete variable."},
		{"For a valid PMF, what must Σ P(X=x) equal?", "1", "All probabilities must sum to 1."},
		{"What is the difference between PMF and PDF?", "PMF for discrete; PDF for continuous", "PMF gives probabilities at specific points; PDF gives density."},
		{"What is the cumulative distribution for a discrete variable?", "F(x) = Σ_{k ≤ x} P(X = k)", "Discrete CDF sums PMF values up to x."},
		{"Can a discrete random variable take infinitely many values?", "yes", "Discrete variables can be countably infinite (e.g., Poisson, Geometric)."},
		{"What is the expected value of a discrete random variable?", "E[X] = Σ x·P(X=x)", "The mean is the probability-weighted sum of all possible values."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

type expectedValueGen struct{}

func (g *expectedValueGen) Generate(difficulty float64) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"What is the expected value of a constant c?", "c", "E[c] = c for any constant c."},
		{"What is E[aX + b] in terms of E[X]?", "a·E[X] + b", "Linearity of expectation: E[aX+b] = aE[X] + b."},
		{"What is the expected value of the sum of random variables?", "E[X+Y] = E[X] + E[Y]", "Linearity: expectation of sum = sum of expectations."},
		{"Is E[XY] always equal to E[X]E[Y]?", "no", "E[XY] = E[X]E[Y] only when X and Y are independent."},
		{"How is variance related to expected value?", "Var[X] = E[X²] - (E[X])²", "Variance = second moment minus squared mean."},
		{"What is E[X²] called?", "second moment", "E[X²] is the second moment about the origin."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}
