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
	reg.Register("stat.func.gaussian", &gaussianFuncGen{})
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
	registerNumeric(reg)
	registerBatch(reg)
}

type readTableGen struct{}

func (g *readTableGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	numRows := rand.Intn(max(1, scale/2)) + 3
	numCols := rand.Intn(max(1, scale/2)) + 2
	// Labels are drawn without replacement: two rows both called Mangoes makes
	// the table unanswerable, because the question asks for one row's total.
	items := []string{"Apples", "Oranges", "Bananas", "Grapes", "Mangoes", "Cherries", "Pears", "Plums"}
	rand.Shuffle(len(items), func(i, j int) { items[i], items[j] = items[j], items[i] })
	rowLabels := append([]string(nil), items[:numRows]...)
	units := []string{"Week 1", "Week 2", "Week 3", "Jan", "Feb", "Mar", "Mon", "Tue", "Wed"}
	rand.Shuffle(len(units), func(i, j int) { units[i], units[j] = units[j], units[i] })
	colLabels := append([]string(nil), units[:numCols]...)
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
	// A pipe table, not padded columns. The question is markdown rendered through
	// remark-gfm, which collapses the runs of spaces that made the grid line up,
	// so a padded table arrived as one line of run-together words: the report of
	// "Feb Mar Mangoes 8" under an unanswerable question.
	var sb strings.Builder
	// The corner cell stays empty: the first column holds row labels, and a
	// header that starts with a label shifts every column one cell left.
	sb.WriteString("|  |")
	for _, label := range colLabels {
		sb.WriteString(" " + label + " |")
	}
	sb.WriteString("\n| --- |")
	for range colLabels {
		sb.WriteString(" --- |")
	}
	sb.WriteString("\n")
	for i := 0; i < numRows; i++ {
		sb.WriteString("| " + rowLabels[i] + " |")
		for j := 0; j < numCols; j++ {
			sb.WriteString(" " + strconv.Itoa(table[i][j]) + " |")
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

// distinctColors returns n different colour names.
//
// The bag questions used to draw each colour independently, so a bag could read
// "2 blue, 11 blue, 11 yellow". That was not only odd prose: pick selected one
// of the two blue groups, so the stored answer was the probability of that group
// while a reader naturally sums both, and a correct answer was marked wrong.
func distinctColors(n int) []string {
	pool := []string{"red", "blue", "green", "yellow", "orange", "purple"}
	rand.Shuffle(len(pool), func(i, j int) { pool[i], pool[j] = pool[j], pool[i] })
	if n > len(pool) {
		n = len(pool)
	}
	return pool[:n]
}

type basicProbGen struct{}

func (g *basicProbGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	colorNames := distinctColors(3)
	counts := make([]int, 3)
	total := 0
	for i := 0; i < 3; i++ {
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
	colorNames := distinctColors(3)
	counts := make([]int, 3)
	total := 0
	for i := 0; i < 3; i++ {
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
	colorNames := distinctColors(3)
	counts := make([]int, 3)
	total := 0
	for i := 0; i < 3; i++ {
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
	switch rand.Intn(5) {
	case 0:
		type covDefItem struct {
			ex  int
			ey  int
			exy int
			cov int
		}
		items := []covDefItem{
			{5, 6, 42, 12},
			{4, 7, 38, 10},
			{8, 5, 55, 15},
			{6, 9, 68, 14},
			{3, 10, 48, 18},
			{7, 8, 76, 20},
			{9, 4, 52, 16},
			{10, 5, 75, 25},
		}
		it := items[rand.Intn(len(items))]
		prod := it.ex * it.ey
		return generator.Problem{
			Question: fmt.Sprintf("Random variables \\(X\\) and \\(Y\\) satisfy \\(E[X] = %d\\), \\(E[Y] = %d\\) and \\(E[XY] = %d\\). Calculate the covariance \\(\\operatorname{Cov}(X, Y)\\).", it.ex, it.ey, it.exy),
			Answer:   fmt.Sprintf("%d", it.cov),
			Explanation: fmt.Sprintf("By definition, \\(\\operatorname{Cov}(X, Y) = E[XY] - E[X]E[Y] = %d - (%d \\cdot %d) = %d - %d = %d\\).",
				it.exy, it.ex, it.ey, it.exy, prod, it.cov),
		}
	case 1:
		type varSumItem struct {
			vx  int
			vy  int
			cov int
			v   int
		}
		items := []varSumItem{
			{15, 20, 6, 47},
			{18, 24, 5, 52},
			{25, 30, 8, 71},
			{12, 16, 4, 36},
			{32, 28, 10, 80},
			{40, 35, 12, 99},
			{22, 18, 7, 54},
			{50, 45, 15, 125},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X\\) and \\(Y\\) have variances \\(\\operatorname{Var}(X) = %d\\) and \\(\\operatorname{Var}(Y) = %d\\), with \\(\\operatorname{Cov}(X, Y) = %d\\). Determine \\(\\operatorname{Var}(X + Y)\\).", it.vx, it.vy, it.cov),
			Answer:   fmt.Sprintf("%d", it.v),
			Explanation: fmt.Sprintf("The variance of a sum is \\(\\operatorname{Var}(X + Y) = \\operatorname{Var}(X) + \\operatorname{Var}(Y) + 2\\operatorname{Cov}(X, Y) = %d + %d + 2(%d) = %d\\).",
				it.vx, it.vy, it.cov, it.v),
		}
	case 2:
		type varDiffItem struct {
			vx  int
			vy  int
			cov int
			v   int
		}
		items := []varDiffItem{
			{20, 15, 4, 27},
			{30, 25, 6, 43},
			{40, 35, 10, 55},
			{28, 22, 5, 40},
			{50, 30, 12, 56},
			{45, 40, 15, 55},
			{36, 24, 8, 44},
			{60, 50, 20, 70},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("Given \\(\\operatorname{Var}(X) = %d\\), \\(\\operatorname{Var}(Y) = %d\\) and \\(\\operatorname{Cov}(X, Y) = %d\\), compute \\(\\operatorname{Var}(X - Y)\\).", it.vx, it.vy, it.cov),
			Answer:   fmt.Sprintf("%d", it.v),
			Explanation: fmt.Sprintf("The variance of a difference is \\(\\operatorname{Var}(X - Y) = \\operatorname{Var}(X) + \\operatorname{Var}(Y) - 2\\operatorname{Cov}(X, Y) = %d + %d - 2(%d) = %d\\).",
				it.vx, it.vy, it.cov, it.v),
		}
	case 3:
		type scaleItem struct {
			cov int
			a   int
			b   int
			ans int
		}
		items := []scaleItem{
			{4, 3, 2, 24},
			{5, 2, 4, 40},
			{6, 3, 3, 54},
			{7, 4, 2, 56},
			{8, 2, 5, 80},
			{3, 5, 4, 60},
			{9, 2, 3, 54},
			{10, 3, 4, 120},
		}
		it := items[rand.Intn(len(items))]
		ab := it.a * it.b
		return generator.Problem{
			Question: fmt.Sprintf("If \\(\\operatorname{Cov}(X, Y) = %d\\), calculate \\(\\operatorname{Cov}(%dX, %dY)\\).", it.cov, it.a, it.b),
			Answer:   fmt.Sprintf("%d", it.ans),
			Explanation: fmt.Sprintf("By the bilinear property of covariance, \\(\\operatorname{Cov}(%dX, %dY) = (%d \\cdot %d)\\operatorname{Cov}(X, Y) = %d \\cdot %d = %d\\).",
				it.a, it.b, it.a, it.b, ab, it.cov, it.ans),
		}
	default:
		type selfCovItem struct {
			v int
		}
		items := []selfCovItem{
			{9},
			{16},
			{25},
			{36},
			{49},
			{64},
			{81},
			{100},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("A random variable \\(X\\) has variance \\(\\operatorname{Var}(X) = %d\\). What is the covariance \\(\\operatorname{Cov}(X, X)\\)?", it.v),
			Answer:   fmt.Sprintf("%d", it.v),
			Explanation: fmt.Sprintf("The covariance of any random variable with itself equals its variance: \\(\\operatorname{Cov}(X, X) = \\operatorname{Var}(X) = %d\\).",
				it.v),
		}
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
		Question:    fmt.Sprintf("Find the geometric mean of %d and %d.", a, b),
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
		Question:    fmt.Sprintf("Find the harmonic mean of %d and %d.", a, b),
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
		Question:    fmt.Sprintf("Find the root mean square (RMS) of %d and %d.", a, b),
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
	switch rand.Intn(5) {
	case 0:
		type param struct {
			n int
			p float64
		}
		params := []param{
			{20, 0.3},
			{50, 0.4},
			{40, 0.5},
			{100, 0.7},
			{80, 0.25},
			{60, 0.6},
		}
		item := params[rand.Intn(len(params))]
		expected := int(math.Round(float64(item.n) * item.p))
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X_1, X_2, \\dots, X_{%d}\\) be independent Bernoulli random variables each with success probability \\(p = %.2f\\). Compute the expected total number of successes \\(E[\\sum_{i=1}^{%d} X_i]\\).", item.n, item.p, item.n),
			Answer:   fmt.Sprintf("%d", expected),
			Explanation: fmt.Sprintf("By linearity of expectation, the expected value of a sum of independent Bernoulli variables is \\(E[\\sum X_i] = %d \\times %.2f = %d\\).",
				item.n, item.p, expected),
		}
	case 1:
		type varParam struct {
			n int
			p float64
			v int
		}
		varParams := []varParam{
			{40, 0.5, 10},
			{100, 0.5, 25},
			{100, 0.2, 16},
			{100, 0.1, 9},
			{50, 0.2, 8},
			{25, 0.2, 4},
		}
		vp := varParams[rand.Intn(len(varParams))]
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X_1, X_2, \\dots, X_{%d}\\) be independent Bernoulli random variables with parameter \\(p = %.1f\\). Find the variance \\(\\operatorname{Var}(\\sum_{i=1}^{%d} X_i)\\).", vp.n, vp.p, vp.n),
			Answer:   fmt.Sprintf("%d", vp.v),
			Explanation: fmt.Sprintf("For independent Bernoulli trials, the variance of the sum is \\(\\operatorname{Var}(\\sum X_i) = %d \\times p(1-p) = %d \\times %.1f \\times %.1f = %d\\).",
				vp.n, vp.n, vp.p, 1-vp.p, vp.v),
		}
	case 2:
		pPct := 10 * (1 + rand.Intn(8))
		p := float64(pPct) / 100.0
		failPct := 100 - pPct
		return generator.Problem{
			Question: fmt.Sprintf("A trial modeled as a Bernoulli random variable succeeds with probability \\(p = %.2f\\). What is the probability (expressed as an integer percentage from 0 to 100) that the trial fails?", p),
			Answer:   fmt.Sprintf("%d", failPct),
			Explanation: fmt.Sprintf("The probability of failure is \\(q = 1 - p = 1 - %.2f = %.2f\\), which equals %d percent.",
				p, 1-p, failPct),
		}
	case 3:
		// The probability is printed as a fraction because the odds it implies
		// are then exact. Rounded to two decimals, 2/3 reads as 0.67, and a
		// learner who divides the printed numbers gets 2.03 for an answer of 2.
		oddsList := []int{1, 2, 3, 4, 9}
		k := oddsList[rand.Intn(len(oddsList))]
		return generator.Problem{
			Question: fmt.Sprintf("A Bernoulli trial has success probability \\(p = \\frac{%d}{%d}\\). Compute the odds in favor of success, defined as the ratio \\(\\frac{p}{1-p}\\).", k, k+1),
			Answer:   fmt.Sprintf("%d", k),
			Explanation: fmt.Sprintf("\\(1-p = \\frac{1}{%d}\\), so \\(\\frac{p}{1-p} = \\frac{%d}{%d} \\div \\frac{1}{%d} = %d\\).",
				k+1, k, k+1, k+1, k),
		}
	default:
		type scaleParam struct {
			a   int
			p   float64
			ans int
		}
		scales := []scaleParam{
			{5, 0.4, 10},
			{10, 0.3, 30},
			{4, 0.5, 8},
			{6, 0.5, 18},
			{5, 0.2, 5},
			{10, 0.8, 80},
		}
		sp := scales[rand.Intn(len(scales))]
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X \\sim \\operatorname{Bernoulli}(%.1f)\\) and define \\(Y = %d X\\). Find the second moment \\(E[Y^2]\\).", sp.p, sp.a),
			Answer:   fmt.Sprintf("%d", sp.ans),
			Explanation: fmt.Sprintf("Because \\(X \\in \\{0, 1\\}\\), \\(X^2 = X\\) and \\(E[X^2] = p = %.1f\\). Therefore \\(E[Y^2] = %d^2 E[X^2] = %d \\times %.1f = %d\\).",
				sp.p, sp.a, sp.a*sp.a, sp.p, sp.ans),
		}
	}
}

type betaGen struct{}

func (g *betaGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch rand.Intn(5) {
	case 0:
		type meanParam struct {
			a, b, pct int
		}
		params := []meanParam{
			{3, 7, 30},
			{1, 3, 25},
			{7, 13, 35},
			{3, 2, 60},
			{9, 16, 36},
			{1, 1, 50},
			{4, 1, 80},
		}
		p := params[rand.Intn(len(params))]
		return generator.Problem{
			Question: fmt.Sprintf("For a random variable \\(X \\sim \\operatorname{Beta}(\\alpha = %d, \\beta = %d)\\), compute the expected value \\(E[X]\\) expressed as an integer percentage from 0 to 100.", p.a, p.b),
			Answer:   fmt.Sprintf("%d", p.pct),
			Explanation: fmt.Sprintf("The mean of a Beta distribution is \\(E[X] = \\frac{\\alpha}{\\alpha + \\beta} = \\frac{%d}{%d} = %.2f\\), which corresponds to %d percent.",
				p.a, p.a+p.b, float64(p.pct)/100.0, p.pct),
		}
	case 1:
		type modeParam struct {
			a, b, pct int
		}
		params := []modeParam{
			{4, 8, 30},
			{2, 4, 25},
			{3, 4, 40},
			{8, 14, 35},
			{5, 2, 80},
			{3, 3, 50},
		}
		p := params[rand.Intn(len(params))]
		num := p.a - 1
		denom := p.a + p.b - 2
		return generator.Problem{
			Question: fmt.Sprintf("For a random variable \\(X \\sim \\operatorname{Beta}(\\alpha = %d, \\beta = %d)\\) with \\(\\alpha, \\beta > 1\\), compute the mode of the distribution expressed as an integer percentage from 0 to 100.", p.a, p.b),
			Answer:   fmt.Sprintf("%d", p.pct),
			Explanation: fmt.Sprintf("The mode of a Beta distribution with parameters greater than 1 is \\(\\operatorname{Mode} = \\frac{\\alpha - 1}{\\alpha + \\beta - 2} = \\frac{%d}{%d} = %.2f\\), which equals %d percent.",
				num, denom, float64(p.pct)/100.0, p.pct),
		}
	case 2:
		a0 := 1 + rand.Intn(6)
		b0 := 1 + rand.Intn(6)
		s := 2 + rand.Intn(10)
		f := 2 + rand.Intn(10)
		aPost := a0 + s
		return generator.Problem{
			Question: fmt.Sprintf("A Bayesian prior is modeled as \\(\\operatorname{Beta}(\\alpha = %d, \\beta = %d)\\). An experiment observes %d successes and %d failures in independent Bernoulli trials. What is the updated parameter \\(\\alpha_{\\text{post}}\\) of the posterior Beta distribution?", a0, b0, s, f),
			Answer:   fmt.Sprintf("%d", aPost),
			Explanation: fmt.Sprintf("Because the Beta prior is conjugate to the Bernoulli likelihood, the posterior parameter adds the observed successes: \\(\\alpha_{\\text{post}} = \\alpha_0 + s = %d + %d = %d\\).",
				a0, s, aPost),
		}
	case 3:
		a := 2 + rand.Intn(8)
		b := 2 + rand.Intn(8)
		total := a + b
		return generator.Problem{
			Question: fmt.Sprintf("In Bayesian statistics with a \\(\\operatorname{Beta}(\\alpha = %d, \\beta = %d)\\) prior, compute the effective prior pseudo-sample size \\(\\alpha + \\beta\\).", a, b),
			Answer:   fmt.Sprintf("%d", total),
			Explanation: fmt.Sprintf("The sum of the shape parameters represents the effective prior sample size: \\(\\alpha + \\beta = %d + %d = %d\\).",
				a, b, total),
		}
	default:
		k := 2 + rand.Intn(6)
		b := 1 + rand.Intn(5)
		a := k * b
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X \\sim \\operatorname{Beta}(\\alpha = %d, \\beta = %d)\\). Find the ratio of expected values \\(\\frac{E[X]}{E[1 - X]}\\).", a, b),
			Answer:   fmt.Sprintf("%d", k),
			Explanation: fmt.Sprintf("Because \\(E[X] = \\frac{\\alpha}{\\alpha + \\beta}\\) and \\(E[1 - X] = \\frac{\\beta}{\\alpha + \\beta}\\), their ratio is \\(\\frac{E[X]}{E[1 - X]} = \\frac{\\alpha}{\\beta} = \\frac{%d}{%d} = %d\\).",
				a, b, k),
		}
	}
}

type binomialDistGen struct{}

func (g *binomialDistGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch rand.Intn(5) {
	case 0:
		type item struct {
			n int
			p float64
		}
		items := []item{
			{20, 0.25},
			{50, 0.40},
			{60, 0.15},
			{80, 0.35},
			{45, 0.20},
			{30, 0.40},
			{75, 0.60},
			{120, 0.05},
		}
		it := items[rand.Intn(len(items))]
		mean := int(math.Round(float64(it.n) * it.p))
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X \\sim \\operatorname{Binomial}(%d, %.2f)\\). Compute the expected value \\(E[X]\\).", it.n, it.p),
			Answer:   fmt.Sprintf("%d", mean),
			Explanation: fmt.Sprintf("The expected value of a binomial random variable is \\(E[X] = n p = %d \\times %.2f = %d\\).",
				it.n, it.p, mean),
		}
	case 1:
		type varItem struct {
			n int
			p float64
			v int
		}
		items := []varItem{
			{36, 0.50, 9},
			{64, 0.50, 16},
			{100, 0.50, 25},
			{144, 0.50, 36},
			{25, 0.20, 4},
			{50, 0.20, 8},
			{75, 0.20, 12},
			{100, 0.20, 16},
			{100, 0.10, 9},
			{200, 0.10, 18},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X \\sim \\operatorname{Binomial}(%d, %.2f)\\). Calculate the variance \\(\\operatorname{Var}(X)\\).", it.n, it.p),
			Answer:   fmt.Sprintf("%d", it.v),
			Explanation: fmt.Sprintf("The variance of a binomial random variable is \\(\\operatorname{Var}(X) = n p (1 - p) = %d \\times %.2f \\times %.2f = %d\\).",
				it.n, it.p, 1-it.p, it.v),
		}
	case 2:
		type sdItem struct {
			n  int
			p  float64
			v  int
			sd int
		}
		items := []sdItem{
			{36, 0.50, 9, 3},
			{64, 0.50, 16, 4},
			{100, 0.50, 25, 5},
			{144, 0.50, 36, 6},
			{25, 0.20, 4, 2},
			{100, 0.20, 16, 4},
			{225, 0.20, 36, 6},
			{100, 0.10, 9, 3},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("A random variable \\(X\\) follows a binomial distribution \\(\\operatorname{Binomial}(%d, %.2f)\\). Find the standard deviation \\(\\sigma\\).", it.n, it.p),
			Answer:   fmt.Sprintf("%d", it.sd),
			Explanation: fmt.Sprintf("The standard deviation is \\(\\sigma = \\sqrt{n p (1 - p)} = \\sqrt{%d \\times %.2f \\times %.2f} = \\sqrt{%d} = %d\\).",
				it.n, it.p, 1-it.p, it.v, it.sd),
		}
	case 3:
		type combItem struct {
			n int
			k int
			c int
		}
		items := []combItem{
			{5, 2, 10},
			{6, 2, 15},
			{6, 3, 20},
			{7, 2, 21},
			{7, 3, 35},
			{8, 2, 28},
			{8, 3, 56},
			{9, 2, 36},
			{10, 2, 45},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("In a binomial experiment with \\(n = %d\\) trials, how many distinct sequences contain exactly \\(k = %d\\) successes?", it.n, it.k),
			Answer:   fmt.Sprintf("%d", it.c),
			Explanation: fmt.Sprintf("The number of sequences is given by the binomial coefficient \\(\\binom{n}{k} = \\binom{%d}{%d} = \\frac{%d!}{%d!(%d - %d)!} = %d\\).",
				it.n, it.k, it.n, it.k, it.n, it.k, it.c),
		}
	default:
		type sumItem struct {
			n1 int
			n2 int
			p  float64
			e  int
		}
		items := []sumItem{
			{30, 20, 0.40, 20},
			{45, 35, 0.25, 20},
			{25, 15, 0.50, 20},
			{40, 60, 0.30, 30},
			{70, 50, 0.15, 18},
			{50, 50, 0.10, 10},
		}
		it := items[rand.Intn(len(items))]
		nTotal := it.n1 + it.n2
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X_1 \\sim \\operatorname{Binomial}(%d, %.2f)\\) and \\(X_2 \\sim \\operatorname{Binomial}(%d, %.2f)\\) be independent random variables. Find the expected value of their sum \\(E[X_1 + X_2]\\).", it.n1, it.p, it.n2, it.p),
			Answer:   fmt.Sprintf("%d", it.e),
			Explanation: fmt.Sprintf("By linearity of expectation, \\(E[X_1 + X_2] = E[X_1] + E[X_2] = (n_1 + n_2)p = (%d + %d) \\times %.2f = %d \\times %.2f = %d\\).",
				it.n1, it.n2, it.p, nTotal, it.p, it.e),
		}
	}
}

type chiSquareGen struct{}

func (g *chiSquareGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch rand.Intn(5) {
	case 0:
		kValues := []int{3, 5, 7, 9, 11, 13, 15, 20}
		k := kValues[rand.Intn(len(kValues))]
		v := 2 * k
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X \\sim \\chi^2(%d)\\) follow a chi-square distribution with %d degrees of freedom. Calculate the variance \\(\\operatorname{Var}(X)\\).", k, k),
			Answer:   fmt.Sprintf("%d", v),
			Explanation: fmt.Sprintf("The variance of a chi-square distribution with \\(k\\) degrees of freedom is \\(\\operatorname{Var}(X) = 2k = 2 \\times %d = %d\\).",
				k, v),
		}
	case 1:
		type sdItem struct {
			k  int
			v  int
			sd int
		}
		items := []sdItem{
			{2, 4, 2},
			{8, 16, 4},
			{18, 36, 6},
			{32, 64, 8},
			{50, 100, 10},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("A random variable \\(X\\) follows a chi-square distribution with %d degrees of freedom: \\(X \\sim \\chi^2(%d)\\). Find the standard deviation \\(\\sigma\\).", it.k, it.k),
			Answer:   fmt.Sprintf("%d", it.sd),
			Explanation: fmt.Sprintf("For a chi-square distribution, the variance is \\(\\operatorname{Var}(X) = 2k = 2 \\times %d = %d\\). The standard deviation is \\(\\sigma = \\sqrt{2k} = \\sqrt{%d} = %d\\).",
				it.k, it.v, it.v, it.sd),
		}
	case 2:
		k1List := []int{4, 6, 8, 10, 12}
		k2List := []int{5, 7, 9, 11, 15}
		k1 := k1List[rand.Intn(len(k1List))]
		k2 := k2List[rand.Intn(len(k2List))]
		sumK := k1 + k2
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X_1 \\sim \\chi^2(%d)\\) and \\(X_2 \\sim \\chi^2(%d)\\) be independent chi-square random variables. Find the expected value of their sum \\(E[X_1 + X_2]\\).", k1, k2),
			Answer:   fmt.Sprintf("%d", sumK),
			Explanation: fmt.Sprintf("By independence and linearity of expectation, \\(E[X_1 + X_2] = E[X_1] + E[X_2] = k_1 + k_2 = %d + %d = %d\\).",
				k1, k2, sumK),
		}
	case 3:
		kValues := []int{4, 6, 8, 10, 12, 16, 25}
		k := kValues[rand.Intn(len(kValues))]
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(Z_1, Z_2, \\dots, Z_{%d}\\) be independent standard normal random variables. What is the degrees of freedom parameter \\(k\\) of the variable \\(W = \\sum_{i=1}^{%d} Z_i^2\\)?", k, k),
			Answer:   fmt.Sprintf("%d", k),
			Explanation: fmt.Sprintf("The sum of \\(k\\) independent squared standard normal variables follows a chi-square distribution \\(\\chi^2(k)\\). Here %d variables are summed, so \\(k = %d\\).",
				k, k),
		}
	default:
		kValues := []int{3, 4, 5, 6, 7, 8, 10}
		k := kValues[rand.Intn(len(kValues))]
		m2 := k*k + 2*k
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X \\sim \\chi^2(%d)\\). Calculate the second moment \\(E[X^2]\\).", k),
			Answer:   fmt.Sprintf("%d", m2),
			Explanation: fmt.Sprintf("Using the variance formula \\(\\operatorname{Var}(X) = E[X^2] - (E[X])^2\\), we find \\(E[X^2] = \\operatorname{Var}(X) + (E[X])^2 = 2(%d) + (%d)^2 = %d + %d = %d\\).",
				k, k, 2*k, k*k, m2),
		}
	}
}

type exponentialDistGen struct{}

func (g *exponentialDistGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch rand.Intn(5) {
	case 0:
		type rateItem struct {
			m      int
			lambda int
		}
		items := []rateItem{
			{5, 12},
			{6, 10},
			{10, 6},
			{12, 5},
			{15, 4},
			{20, 3},
			{30, 2},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("Customer arrivals follow an exponential distribution with a mean waiting time of %d minutes. What is the arrival rate \\(\\lambda\\) in events per hour?", it.m),
			Answer:   fmt.Sprintf("%d", it.lambda),
			Explanation: fmt.Sprintf("With 60 minutes per hour and mean \\(E[X] = %d\\) minutes, the hourly rate parameter is \\(\\lambda = \\frac{60}{%d} = %d\\).",
				it.m, it.m, it.lambda),
		}
	case 1:
		type memItem struct {
			mu int
			s  int
		}
		items := []memItem{
			{4, 10},
			{5, 8},
			{6, 12},
			{8, 20},
			{10, 15},
			{12, 24},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("A component lifetime follows an exponential distribution with mean \\(E[X] = %d\\) hours. Given that the component has operated without failure for %d hours, what is the expected remaining lifetime until failure?", it.mu, it.s),
			Answer:   fmt.Sprintf("%d", it.mu),
			Explanation: fmt.Sprintf("By the memoryless property of the exponential distribution, prior operating time does not affect remaining lifetime: \\(E[X - %d \\mid X > %d] = E[X] = %d\\).",
				it.s, it.s, it.mu),
		}
	case 2:
		l1List := []int{2, 3, 4, 5, 6}
		l2List := []int{3, 4, 5, 7, 8}
		l1 := l1List[rand.Intn(len(l1List))]
		l2 := l2List[rand.Intn(len(l2List))]
		sumL := l1 + l2
		return generator.Problem{
			Question: fmt.Sprintf("Two independent sensors monitor a machine. The time until sensor 1 triggers follows \\(\\operatorname{Exp}(\\lambda_1)\\) with rate \\(\\lambda_1 = %d\\) per day, and sensor 2 follows \\(\\operatorname{Exp}(\\lambda_2)\\) with rate \\(\\lambda_2 = %d\\) per day. What is the rate parameter \\(\\lambda_{\\min}\\) of the first sensor trigger time \\(T = \\min(X_1, X_2)\\)?", l1, l2),
			Answer:   fmt.Sprintf("%d", sumL),
			Explanation: fmt.Sprintf("The minimum of independent exponential variables is exponentially distributed with rate equal to the sum of the rates: \\(\\lambda_{\\min} = \\lambda_1 + \\lambda_2 = %d + %d = %d\\).",
				l1, l2, sumL),
		}
	case 3:
		muList := []int{3, 4, 5, 6, 7, 8, 10}
		mu := muList[rand.Intn(len(muList))]
		m2 := 2 * mu * mu
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X \\sim \\operatorname{Exp}(\\lambda)\\) have mean \\(E[X] = %d\\). Calculate the second moment \\(E[X^2]\\).", mu),
			Answer:   fmt.Sprintf("%d", m2),
			Explanation: fmt.Sprintf("For an exponential random variable, \\(\\operatorname{Var}(X) = \\mu^2 = %d^2 = %d\\). Therefore, \\(E[X^2] = \\operatorname{Var}(X) + (E[X])^2 = %d + %d = %d\\).",
				mu, mu*mu, mu*mu, mu*mu, m2),
		}
	default:
		type sumItem struct {
			n  int
			mu int
		}
		items := []sumItem{
			{4, 5},
			{5, 6},
			{6, 7},
			{8, 5},
			{10, 4},
			{12, 3},
		}
		it := items[rand.Intn(len(items))]
		total := it.n * it.mu
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X_1, X_2, \\dots, X_{%d}\\) be independent exponential random variables each having mean \\(E[X_i] = %d\\). Find the expected value of their sum \\(E[\\sum_{i=1}^{%d} X_i]\\).", it.n, it.mu, it.n),
			Answer:   fmt.Sprintf("%d", total),
			Explanation: fmt.Sprintf("By linearity of expectation, the expected value of the sum is \\(E[\\sum_{i=1}^{%d} X_i] = %d \\times E[X_1] = %d \\times %d = %d\\).",
				it.n, it.n, it.n, it.mu, total),
		}
	}
}

type gammaDistGen struct{}

func (g *gammaDistGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch rand.Intn(5) {
	case 0:
		type meanItem struct {
			k     int
			theta int
		}
		items := []meanItem{
			{3, 4},
			{4, 5},
			{5, 2},
			{6, 3},
			{7, 4},
			{8, 5},
		}
		it := items[rand.Intn(len(items))]
		mean := it.k * it.theta
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X \\sim \\operatorname{Gamma}(k, \\theta)\\) have shape parameter \\(k = %d\\) and scale parameter \\(\\theta = %d\\). Calculate the expected value \\(E[X]\\).", it.k, it.theta),
			Answer:   fmt.Sprintf("%d", mean),
			Explanation: fmt.Sprintf("The expected value of a Gamma distribution with shape \\(k\\) and scale \\(\\theta\\) is \\(E[X] = k \\theta = %d \\times %d = %d\\).",
				it.k, it.theta, mean),
		}
	case 1:
		type varItem struct {
			k     int
			theta int
		}
		items := []varItem{
			{2, 3},
			{3, 4},
			{4, 3},
			{5, 2},
			{6, 4},
			{2, 5},
		}
		it := items[rand.Intn(len(items))]
		v := it.k * it.theta * it.theta
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X \\sim \\operatorname{Gamma}(k, \\theta)\\) have shape \\(k = %d\\) and scale \\(\\theta = %d\\). Calculate the variance \\(\\operatorname{Var}(X)\\).", it.k, it.theta),
			Answer:   fmt.Sprintf("%d", v),
			Explanation: fmt.Sprintf("The variance of a Gamma distribution with shape \\(k\\) and scale \\(\\theta\\) is \\(\\operatorname{Var}(X) = k \\theta^2 = %d \\times %d^2 = %d \\times %d = %d\\).",
				it.k, it.theta, it.k, it.theta*it.theta, v),
		}
	case 2:
		type sdItem struct {
			k     int
			theta int
			sd    int
		}
		items := []sdItem{
			{4, 3, 6},
			{4, 5, 10},
			{9, 2, 6},
			{9, 4, 12},
			{16, 2, 8},
			{16, 3, 12},
		}
		it := items[rand.Intn(len(items))]
		v := it.k * it.theta * it.theta
		return generator.Problem{
			Question: fmt.Sprintf("A random variable \\(X\\) follows a Gamma distribution with shape \\(k = %d\\) and scale \\(\\theta = %d\\). Find the standard deviation \\(\\sigma\\).", it.k, it.theta),
			Answer:   fmt.Sprintf("%d", it.sd),
			Explanation: fmt.Sprintf("The variance is \\(\\operatorname{Var}(X) = k \\theta^2 = %d \\times %d^2 = %d\\). Therefore, the standard deviation is \\(\\sigma = \\sqrt{k \\theta^2} = \\sqrt{%d} = %d\\).",
				it.k, it.theta, v, v, it.sd),
		}
	case 3:
		k1List := []int{2, 3, 4, 5}
		k2List := []int{3, 4, 5, 6}
		thetaList := []int{2, 3, 4}
		k1 := k1List[rand.Intn(len(k1List))]
		k2 := k2List[rand.Intn(len(k2List))]
		th := thetaList[rand.Intn(len(thetaList))]
		sumK := k1 + k2
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X_1 \\sim \\operatorname{Gamma}(%d, %d)\\) and \\(X_2 \\sim \\operatorname{Gamma}(%d, %d)\\) be independent random variables sharing scale parameter \\(\\theta = %d\\). What is the shape parameter \\(k_{\\text{sum}}\\) of their sum \\(S = X_1 + X_2\\)?", k1, th, k2, th, th),
			Answer:   fmt.Sprintf("%d", sumK),
			Explanation: fmt.Sprintf("The sum of independent Gamma variables sharing the same scale parameter has shape parameter equal to the sum of the individual shapes: \\(k_{\\text{sum}} = k_1 + k_2 = %d + %d = %d\\).",
				k1, k2, sumK),
		}
	default:
		type modeItem struct {
			k     int
			theta int
		}
		items := []modeItem{
			{4, 3},
			{5, 2},
			{6, 4},
			{3, 5},
			{7, 3},
			{8, 2},
		}
		it := items[rand.Intn(len(items))]
		mode := (it.k - 1) * it.theta
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X\\) follow a Gamma distribution with shape parameter \\(k = %d\\) and scale parameter \\(\\theta = %d\\). Find the mode of \\(X\\).", it.k, it.theta),
			Answer:   fmt.Sprintf("%d", mode),
			Explanation: fmt.Sprintf("For shape parameter \\(k \\ge 1\\), the mode of \\(\\operatorname{Gamma}(k, \\theta)\\) is \\((k - 1)\\theta = (%d - 1) \\times %d = %d \\times %d = %d\\).",
				it.k, it.theta, it.k-1, it.theta, mode),
		}
	}
}

type geometricDistGen struct{}

func (g *geometricDistGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch rand.Intn(5) {
	case 0:
		type item struct {
			p   float64
			exp int
		}
		items := []item{
			{0.50, 2},
			{0.25, 4},
			{0.20, 5},
			{0.10, 10},
			{0.05, 20},
			{0.04, 25},
			{0.02, 50},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X\\) follow a geometric distribution with success probability \\(p = %.2f\\), counting trials until the first success. Find the expected value \\(E[X]\\).", it.p),
			Answer:   fmt.Sprintf("%d", it.exp),
			Explanation: fmt.Sprintf("The expected value of a geometric random variable is \\(E[X] = \\frac{1}{p} = \\frac{1}{%.2f} = %d\\).",
				it.p, it.exp),
		}
	case 1:
		mList := []int{2, 3, 4, 5, 6, 10}
		m := mList[rand.Intn(len(mList))]
		v := m * (m - 1)
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X \\sim \\operatorname{Geometric}(p)\\) with success probability \\(p = \\frac{1}{%d}\\). Compute the variance \\(\\operatorname{Var}(X)\\).", m),
			Answer:   fmt.Sprintf("%d", v),
			Explanation: fmt.Sprintf("The variance of a geometric distribution is \\(\\operatorname{Var}(X) = \\frac{1-p}{p^2} = \\frac{1 - 1/%d}{(1/%d)^2} = \\frac{%d/%d}{1/%d} = %d \\times %d = %d\\).",
				m, m, m-1, m, m*m, m-1, m, v),
		}
	case 2:
		type cdfItem struct {
			p   float64
			k   int
			pct int
		}
		items := []cdfItem{
			{0.50, 1, 50},
			{0.50, 2, 75},
			{0.20, 1, 20},
			{0.20, 2, 36},
			{0.10, 1, 10},
			{0.10, 2, 19},
			{0.60, 2, 84},
			{0.70, 2, 91},
			{0.80, 2, 96},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("In independent Bernoulli trials each having success probability \\(p = %.2f\\), what is the probability (expressed as an integer percentage from 0 to 100) that the first success occurs on or before trial \\(k = %d\\)?", it.p, it.k),
			Answer:   fmt.Sprintf("%d", it.pct),
			Explanation: fmt.Sprintf("The cumulative probability is \\(P(X \\le %d) = 1 - (1-p)^{%d} = 1 - (%.2f)^{%d} = 1 - %.2f = %.2f\\), which equals %d percent.",
				it.k, it.k, 1-it.p, it.k, float64(100-it.pct)/100.0, float64(it.pct)/100.0, it.pct),
		}
	case 3:
		type ratioItem struct {
			p float64
			r int
		}
		items := []ratioItem{
			{0.50, 2},
			{0.75, 4},
			{0.80, 5},
			{0.90, 10},
			{0.95, 20},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X \\sim \\operatorname{Geometric}(p)\\) with success probability \\(p = %.2f\\). Compute the probability ratio \\(\\frac{P(X = 1)}{P(X = 2)}\\).", it.p),
			Answer:   fmt.Sprintf("%d", it.r),
			Explanation: fmt.Sprintf("For a geometric distribution, \\(P(X = 1) = p\\) and \\(P(X = 2) = (1-p)p\\). The ratio is \\(\\frac{P(X = 1)}{P(X = 2)} = \\frac{1}{1-p} = \\frac{1}{1 - %.2f} = \\frac{1}{%.2f} = %d\\).",
				it.p, 1-it.p, it.r),
		}
	default:
		type memItem struct {
			p float64
			s int
			e int
		}
		items := []memItem{
			{0.20, 5, 5},
			{0.25, 4, 4},
			{0.10, 8, 10},
			{0.50, 3, 2},
			{0.05, 12, 20},
			{0.04, 15, 25},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("A geometric random variable \\(X\\) has success probability \\(p = %.2f\\). Given that the first %d trials resulted in failure, what is the expected number of additional trials needed until the first success?", it.p, it.s),
			Answer:   fmt.Sprintf("%d", it.e),
			Explanation: fmt.Sprintf("By the memoryless property of the geometric distribution, prior failures do not affect future trials: \\(E[X - %d \\mid X > %d] = E[X] = \\frac{1}{p} = \\frac{1}{%.2f} = %d\\).",
				it.s, it.s, it.p, it.e),
		}
	}
}

type hypergeometricGen struct{}

func (g *hypergeometricGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch rand.Intn(5) {
	case 0:
		type expItem struct {
			nTotal int
			k      int
			n      int
		}
		items := []expItem{
			{20, 8, 5},
			{30, 12, 10},
			{40, 16, 15},
			{50, 20, 25},
			{24, 6, 12},
			{36, 9, 16},
			{50, 15, 20},
			{60, 24, 15},
		}
		it := items[rand.Intn(len(items))]
		mean := it.n * it.k / it.nTotal
		return generator.Problem{
			Question: fmt.Sprintf("An urn contains %d total items including %d marked items. If a random sample of %d items is drawn without replacement, find the expected number of marked items \\(E[X]\\).", it.nTotal, it.k, it.n),
			Answer:   fmt.Sprintf("%d", mean),
			Explanation: fmt.Sprintf("The expected number of marked items in hypergeometric sampling without replacement is \\(E[X] = n \\frac{K}{N} = %d \\times \\frac{%d}{%d} = %d\\).",
				it.n, it.k, it.nTotal, mean),
		}
	case 1:
		type failItem struct {
			nTotal int
			k      int
			n      int
		}
		items := []failItem{
			{20, 8, 5},
			{30, 12, 10},
			{40, 16, 15},
			{50, 20, 25},
			{60, 24, 15},
		}
		it := items[rand.Intn(len(items))]
		nonDef := it.nTotal - it.k
		expNonDef := it.n * nonDef / it.nTotal
		return generator.Problem{
			Question: fmt.Sprintf("A batch of %d items contains %d defective items and %d non-defective items. If %d items are sampled without replacement, what is the expected number of non-defective items in the sample?", it.nTotal, it.k, nonDef, it.n),
			Answer:   fmt.Sprintf("%d", expNonDef),
			Explanation: fmt.Sprintf("The expected number of non-defective items is \\(E[X_{\\text{non-def}}] = n \\frac{N - K}{N} = %d \\times \\frac{%d}{%d} = %d\\).",
				it.n, nonDef, it.nTotal, expNonDef),
		}
	case 2:
		type maxItem struct {
			nTotal int
			k      int
			n      int
			max    int
		}
		items := []maxItem{
			{25, 4, 8, 4},
			{30, 5, 9, 5},
			{20, 9, 6, 6},
			{40, 6, 12, 6},
			{35, 14, 8, 8},
			{50, 7, 15, 7},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("From a population of \\(N = %d\\) containing \\(K = %d\\) target items, a random sample of \\(n = %d\\) is selected without replacement. What is the theoretical maximum number of target items that can appear in the sample?", it.nTotal, it.k, it.n),
			Answer:   fmt.Sprintf("%d", it.max),
			Explanation: fmt.Sprintf("The number of observed target items cannot exceed the total target items in the population \\(K = %d\\) or the sample size \\(n = %d\\): \\(\\max(X) = \\min(n, K) = \\min(%d, %d) = %d\\).",
				it.k, it.n, it.n, it.k, it.max),
		}
	case 3:
		type minItem struct {
			nTotal int
			k      int
			n      int
			min    int
		}
		items := []minItem{
			{20, 15, 8, 3},
			{25, 20, 9, 4},
			{30, 24, 10, 4},
			{35, 28, 12, 5},
			{40, 32, 14, 6},
		}
		it := items[rand.Intn(len(items))]
		nonTarget := it.nTotal - it.k
		return generator.Problem{
			Question: fmt.Sprintf("A committee of %d members contains %d seniors and %d juniors. If a subcommittee of %d members is chosen at random without replacement, what is the minimum number of seniors that must be on the subcommittee?", it.nTotal, it.k, nonTarget, it.n),
			Answer:   fmt.Sprintf("%d", it.min),
			Explanation: fmt.Sprintf("Even if all %d juniors were selected, the remaining seats must be filled by seniors: \\(\\min(X) = \\max(0, n - (N - K)) = %d - (%d - %d) = %d - %d = %d\\).",
				nonTarget, it.n, it.nTotal, it.k, it.n, nonTarget, it.min),
		}
	default:
		type combItem struct {
			nTotal int
			n      int
			comb   int
		}
		items := []combItem{
			{6, 2, 15},
			{7, 2, 21},
			{7, 3, 35},
			{8, 2, 28},
			{8, 3, 56},
			{9, 2, 36},
			{10, 2, 45},
			{10, 3, 120},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("In a hypergeometric experiment with total population size \\(N = %d\\), how many distinct unordered samples of size \\(n = %d\\) can be drawn without replacement?", it.nTotal, it.n),
			Answer:   fmt.Sprintf("%d", it.comb),
			Explanation: fmt.Sprintf("The total number of possible samples is given by the combination formula \\(\\binom{N}{n} = \\binom{%d}{%d} = \\frac{%d!}{%d!(%d - %d)!} = %d\\).",
				it.nTotal, it.n, it.nTotal, it.n, it.nTotal, it.n, it.comb),
		}
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
	switch rand.Intn(5) {
	case 0:
		type sdItem struct {
			lambda int
			sd     int
		}
		items := []sdItem{
			{4, 2},
			{9, 3},
			{16, 4},
			{25, 5},
			{36, 6},
			{49, 7},
			{64, 8},
			{81, 9},
			{100, 10},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("A random variable \\(X\\) follows a Poisson distribution \\(\\operatorname{Poisson}(%d)\\). Find the standard deviation \\(\\sigma\\).", it.lambda),
			Answer:   fmt.Sprintf("%d", it.sd),
			Explanation: fmt.Sprintf("For a Poisson distribution, the variance equals the parameter: \\(\\operatorname{Var}(X) = \\lambda = %d\\). Therefore, the standard deviation is \\(\\sigma = \\sqrt{\\lambda} = \\sqrt{%d} = %d\\).",
				it.lambda, it.lambda, it.sd),
		}
	case 1:
		rates := []int{2, 3, 4, 5, 6, 7, 8}
		times := []int{2, 3, 4, 5}
		lambda := rates[rand.Intn(len(rates))]
		t := times[rand.Intn(len(times))]
		total := lambda * t
		return generator.Problem{
			Question: fmt.Sprintf("Events occur according to a Poisson process at an average rate of %d events per hour. How many events are expected to occur over a %d-hour interval?", lambda, t),
			Answer:   fmt.Sprintf("%d", total),
			Explanation: fmt.Sprintf("For a Poisson process with rate \\(\\lambda = %d\\), the expected number of events in \\(t = %d\\) hours is \\(E[N(t)] = \\lambda t = %d \\times %d = %d\\).",
				lambda, t, lambda, t, total),
		}
	case 2:
		l1List := []int{3, 5, 7, 10, 12, 15}
		l2List := []int{4, 6, 8, 11, 14, 20}
		l1 := l1List[rand.Intn(len(l1List))]
		l2 := l2List[rand.Intn(len(l2List))]
		sum := l1 + l2
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X_1 \\sim \\operatorname{Poisson}(%d)\\) and \\(X_2 \\sim \\operatorname{Poisson}(%d)\\) be independent random variables. Find the variance of their sum \\(\\operatorname{Var}(X_1 + X_2)\\).", l1, l2),
			Answer:   fmt.Sprintf("%d", sum),
			Explanation: fmt.Sprintf("Since \\(X_1\\) and \\(X_2\\) are independent, \\(\\operatorname{Var}(X_1 + X_2) = \\operatorname{Var}(X_1) + \\operatorname{Var}(X_2) = \\lambda_1 + \\lambda_2 = %d + %d = %d\\).",
				l1, l2, sum),
		}
	case 3:
		type ratioItem struct {
			k      int
			lambda int
			ratio  int
		}
		items := []ratioItem{
			{2, 6, 3},
			{2, 8, 4},
			{2, 10, 5},
			{3, 12, 4},
			{3, 15, 5},
			{4, 20, 5},
			{4, 24, 6},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X \\sim \\operatorname{Poisson}(%d)\\). Find the probability ratio \\(\\frac{P(X = %d)}{P(X = %d)}\\).", it.lambda, it.k, it.k-1),
			Answer:   fmt.Sprintf("%d", it.ratio),
			Explanation: fmt.Sprintf("From the Poisson PMF, consecutive terms satisfy \\(\\frac{P(X = k)}{P(X = k-1)} = \\frac{\\lambda}{k}\\). Here \\(\\frac{P(X = %d)}{P(X = %d)} = \\frac{%d}{%d} = %d\\).",
				it.k, it.k-1, it.lambda, it.k, it.ratio),
		}
	default:
		lambdas := []int{2, 3, 4, 5, 6, 7, 8, 10}
		lam := lambdas[rand.Intn(len(lambdas))]
		m2 := lam + lam*lam
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X \\sim \\operatorname{Poisson}(%d)\\). Compute the second moment \\(E[X^2]\\).", lam),
			Answer:   fmt.Sprintf("%d", m2),
			Explanation: fmt.Sprintf("Using \\(\\operatorname{Var}(X) = E[X^2] - (E[X])^2\\), we find \\(E[X^2] = \\lambda + \\lambda^2 = %d + %d = %d\\).",
				lam, lam*lam, m2),
		}
	}
}

type studentTGen struct{}

func (g *studentTGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch rand.Intn(5) {
	case 0:
		type oneSampleItem struct {
			n int
		}
		items := []oneSampleItem{
			{15},
			{20},
			{25},
			{30},
			{36},
			{45},
			{50},
			{60},
		}
		it := items[rand.Intn(len(items))]
		df := it.n - 1
		return generator.Problem{
			Question: fmt.Sprintf("A researcher conducts a one-sample t-test using a random sample of \\(n = %d\\) observations. How many degrees of freedom \\(\\nu\\) does the corresponding Student's t-distribution have?", it.n),
			Answer:   fmt.Sprintf("%d", df),
			Explanation: fmt.Sprintf("For a single sample of size \\(n\\), the degrees of freedom are \\(\\nu = n - 1 = %d - 1 = %d\\).",
				it.n, df),
		}
	case 1:
		type twoSampleItem struct {
			n1 int
			n2 int
		}
		items := []twoSampleItem{
			{15, 20},
			{12, 18},
			{25, 25},
			{16, 24},
			{30, 20},
			{14, 21},
			{22, 18},
			{35, 25},
		}
		it := items[rand.Intn(len(items))]
		df := it.n1 + it.n2 - 2
		return generator.Problem{
			Question: fmt.Sprintf("An independent two-sample t-test compares two treatment groups with sample sizes \\(n_1 = %d\\) and \\(n_2 = %d\\). Assuming pooled variance, calculate the total degrees of freedom \\(\\nu\\).", it.n1, it.n2),
			Answer:   fmt.Sprintf("%d", df),
			Explanation: fmt.Sprintf("Under pooled variance, the degrees of freedom are \\(\\nu = n_1 + n_2 - 2 = %d + %d - 2 = %d\\).",
				it.n1, it.n2, df),
		}
	case 2:
		type meanItem struct {
			nu int
		}
		items := []meanItem{
			{3},
			{5},
			{8},
			{12},
			{15},
			{20},
			{28},
			{35},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(T \\sim t(%d)\\) follow a Student's t-distribution with \\(\\nu = %d\\) degrees of freedom. What is the expected value \\(E[T]\\)?", it.nu, it.nu),
			Answer:   "0",
			Explanation: fmt.Sprintf("Because the Student's t-distribution is symmetric about zero and \\(\\nu = %d > 1\\), the expected value is \\(E[T] = 0\\).",
				it.nu),
		}
	case 3:
		type tStatItem struct {
			n    int
			xbar int
			mu0  int
			s    int
			t    int
		}
		items := []tStatItem{
			{16, 58, 50, 8, 4},
			{25, 115, 100, 15, 5},
			{36, 84, 75, 18, 3},
			{16, 46, 40, 12, 2},
			{64, 72, 60, 32, 3},
			{25, 130, 110, 20, 5},
			{36, 95, 85, 12, 5},
			{16, 66, 54, 16, 3},
		}
		it := items[rand.Intn(len(items))]
		diff := it.xbar - it.mu0
		se := diff / it.t
		return generator.Problem{
			Question: fmt.Sprintf("A sample of size \\(n = %d\\) has sample mean \\(\\bar{x} = %d\\) and sample standard deviation \\(s = %d\\). Test the hypothesis \\(H_0: \\mu = %d\\) by calculating the sample t-statistic \\(t\\).", it.n, it.xbar, it.s, it.mu0),
			Answer:   fmt.Sprintf("%d", it.t),
			Explanation: fmt.Sprintf("The sample standard error is \\(SE = \\frac{s}{\\sqrt{n}} = \\frac{%d}{\\sqrt{%d}} = %d\\). The test statistic is \\(t = \\frac{\\bar{x} - \\mu_0}{SE} = \\frac{%d - %d}{%d} = %d\\).",
				it.s, it.n, se, it.xbar, it.mu0, se, it.t),
		}
	default:
		type pairedItem struct {
			n int
		}
		items := []pairedItem{
			{18},
			{25},
			{32},
			{40},
			{50},
			{15},
			{28},
			{45},
		}
		it := items[rand.Intn(len(items))]
		df := it.n - 1
		return generator.Problem{
			Question: fmt.Sprintf("A paired-samples experiment measures before-and-after responses from \\(n = %d\\) participants. What are the degrees of freedom \\(\\nu\\) for the paired t-test?", it.n),
			Answer:   fmt.Sprintf("%d", df),
			Explanation: fmt.Sprintf("A paired t-test reduces the problem to %d difference scores, yielding degrees of freedom \\(\\nu = n - 1 = %d - 1 = %d\\).",
				it.n, it.n, df),
		}
	}
}

type uniformDistGen struct{}

func (g *uniformDistGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch rand.Intn(5) {
	case 0:
		type meanItem struct {
			a int
			b int
		}
		items := []meanItem{
			{4, 16},
			{10, 34},
			{6, 26},
			{12, 48},
			{5, 35},
			{8, 40},
			{14, 50},
		}
		it := items[rand.Intn(len(items))]
		mean := (it.a + it.b) / 2
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X \\sim \\operatorname{Uniform}(%d, %d)\\) be a continuous uniform random variable. Find the expected value \\(E[X]\\).", it.a, it.b),
			Answer:   fmt.Sprintf("%d", mean),
			Explanation: fmt.Sprintf("The expected value of a uniform random variable on interval \\([a, b]\\) is \\(E[X] = \\frac{a + b}{2} = \\frac{%d + %d}{2} = %d\\).",
				it.a, it.b, mean),
		}
	case 1:
		type varItem struct {
			a int
			b int
		}
		items := []varItem{
			{0, 12},
			{2, 14},
			{1, 19},
			{0, 24},
			{4, 16},
			{5, 11},
			{10, 40},
			{7, 25},
		}
		it := items[rand.Intn(len(items))]
		diff := it.b - it.a
		v := diff * diff / 12
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X \\sim \\operatorname{Uniform}(%d, %d)\\). Calculate the variance \\(\\operatorname{Var}(X)\\).", it.a, it.b),
			Answer:   fmt.Sprintf("%d", v),
			Explanation: fmt.Sprintf("The variance of a uniform distribution on \\([a, b]\\) is \\(\\operatorname{Var}(X) = \\frac{(b - a)^2}{12} = \\frac{(%d - %d)^2}{12} = \\frac{%d}{12} = %d\\).",
				it.b, it.a, diff*diff, v),
		}
	case 2:
		type cdfItem struct {
			a   int
			b   int
			x   int
			pct int
		}
		items := []cdfItem{
			{0, 20, 8, 40},
			{10, 50, 20, 25},
			{5, 25, 15, 50},
			{0, 50, 35, 70},
			{20, 60, 50, 75},
			{10, 35, 30, 80},
			{0, 40, 24, 60},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("A random variable \\(X\\) is uniformly distributed on the interval \\([%d, %d]\\). What is the probability \\(P(X \\le %d)\\) expressed as an integer percentage from 0 to 100?", it.a, it.b, it.x),
			Answer:   fmt.Sprintf("%d", it.pct),
			Explanation: fmt.Sprintf("For a uniform distribution on \\([a, b]\\), the cumulative probability is \\(P(X \\le %d) = \\frac{%d - %d}{%d - %d} = \\frac{%d}{%d} = %.2f\\), which equals %d percent.",
				it.x, it.x, it.a, it.b, it.a, it.x-it.a, it.b-it.a, float64(it.pct)/100.0, it.pct),
		}
	case 3:
		type subItem struct {
			a   int
			b   int
			c   int
			d   int
			pct int
		}
		items := []subItem{
			{0, 100, 25, 65, 40},
			{10, 60, 20, 45, 50},
			{0, 50, 15, 30, 30},
			{20, 70, 30, 60, 60},
			{0, 80, 16, 56, 50},
			{0, 100, 15, 85, 70},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X \\sim \\operatorname{Uniform}(%d, %d)\\). Find the probability \\(P(%d \\le X \\le %d)\\) expressed as an integer percentage from 0 to 100.", it.a, it.b, it.c, it.d),
			Answer:   fmt.Sprintf("%d", it.pct),
			Explanation: fmt.Sprintf("The probability that \\(X\\) falls in sub-interval \\([%d, %d]\\) is \\(P(%d \\le X \\le %d) = \\frac{%d - %d}{%d - %d} = \\frac{%d}{%d} = %.2f\\), which equals %d percent.",
				it.c, it.d, it.c, it.d, it.d, it.c, it.b, it.a, it.d-it.c, it.b-it.a, float64(it.pct)/100.0, it.pct),
		}
	default:
		type m2Item struct {
			a  int
			b  int
			m2 int
		}
		items := []m2Item{
			{0, 6, 12},
			{0, 9, 27},
			{0, 12, 48},
			{0, 15, 75},
			{3, 6, 21},
			{2, 8, 28},
			{3, 9, 39},
			{4, 10, 52},
		}
		it := items[rand.Intn(len(items))]
		sumSq := it.a*it.a + it.a*it.b + it.b*it.b
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X \\sim \\operatorname{Uniform}(%d, %d)\\). Compute the second moment \\(E[X^2]\\).", it.a, it.b),
			Answer:   fmt.Sprintf("%d", it.m2),
			Explanation: fmt.Sprintf("The second moment of a continuous uniform distribution on \\([a, b]\\) is \\(E[X^2] = \\frac{a^2 + ab + b^2}{3} = \\frac{%d}{3} = %d\\).",
				sumSq, it.m2),
		}
	}
}

type zTableGen struct{}

func (g *zTableGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch rand.Intn(5) {
	case 0:
		type zScoreItem struct {
			x     int
			mu    int
			sigma int
			z     int
		}
		items := []zScoreItem{
			{85, 70, 5, 3},
			{116, 100, 8, 2},
			{42, 50, 4, -2},
			{135, 120, 5, 3},
			{65, 80, 5, -3},
			{148, 100, 12, 4},
			{76, 70, 3, 2},
			{55, 65, 5, -2},
		}
		it := items[rand.Intn(len(items))]
		diff := it.x - it.mu
		return generator.Problem{
			Question: fmt.Sprintf("A normally distributed population has mean \\(\\mu = %d\\) and standard deviation \\(\\sigma = %d\\). Calculate the z-score for an observed value \\(x = %d\\).", it.mu, it.sigma, it.x),
			Answer:   fmt.Sprintf("%d", it.z),
			Explanation: fmt.Sprintf("The z-score standardizes the observation: \\(z = \\frac{x - \\mu}{\\sigma} = \\frac{%d - %d}{%d} = \\frac{%d}{%d} = %d\\).",
				it.x, it.mu, it.sigma, diff, it.sigma, it.z),
		}
	case 1:
		type rawItem struct {
			mu    int
			sigma int
			z     int
			x     int
		}
		items := []rawItem{
			{100, 15, 2, 130},
			{75, 8, 3, 99},
			{50, 6, -2, 38},
			{80, 10, 3, 110},
			{120, 12, 2, 144},
			{60, 5, -3, 45},
			{90, 7, 4, 118},
			{70, 9, -2, 52},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("A test has normal distribution with mean \\(\\mu = %d\\) and standard deviation \\(\\sigma = %d\\). Find the raw score \\(x\\) corresponding to \\(z = %d\\).", it.mu, it.sigma, it.z),
			Answer:   fmt.Sprintf("%d", it.x),
			Explanation: fmt.Sprintf("Unstandardizing the score gives \\(x = \\mu + z\\sigma = %d + (%d \\cdot %d) = %d + %d = %d\\).",
				it.mu, it.z, it.sigma, it.mu, it.z*it.sigma, it.x),
		}
	case 2:
		type empItem struct {
			mu    int
			sigma int
			k     int
			pct   int
		}
		items := []empItem{
			{100, 15, 1, 68},
			{50, 5, 2, 95},
			{70, 10, 1, 68},
			{120, 20, 2, 95},
			{80, 8, 1, 68},
			{90, 12, 2, 95},
			{60, 6, 1, 68},
			{110, 14, 2, 95},
		}
		it := items[rand.Intn(len(items))]
		low := it.mu - it.k*it.sigma
		high := it.mu + it.k*it.sigma
		return generator.Problem{
			Question: fmt.Sprintf("A variable follows a normal distribution with mean \\(\\mu = %d\\) and standard deviation \\(\\sigma = %d\\). According to the empirical rule, what integer percentage of observations lies between %d and %d?", it.mu, it.sigma, low, high),
			Answer:   fmt.Sprintf("%d", it.pct),
			Explanation: fmt.Sprintf("The interval \\([%d, %d]\\) corresponds to \\(\\mu \\pm %d\\sigma\\). By the empirical rule, approximately %d\\%% of data falls within %d standard deviations of the mean.",
				low, high, it.k, it.pct, it.k),
		}
	case 3:
		type centItem struct {
			query string
		}
		items := []centItem{
			{"mean"},
			{"median"},
			{"mode"},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("For a standard normal distribution \\(Z \\sim \\mathcal{N}(0, 1)\\), what is the %s value?", it.query),
			Answer:   "0",
			Explanation: fmt.Sprintf("By symmetry and definition of the standard normal distribution, the %s is exactly \\(z = 0\\).",
				it.query),
		}
	default:
		type stdNormItem struct {
			property string
			val      int
			symbol   string
		}
		items := []stdNormItem{
			{"variance", 1, "\\sigma^2"},
			{"standard deviation", 1, "\\sigma"},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("What is the %s of the standard normal distribution \\(Z \\sim \\mathcal{N}(0, 1)\\)?", it.property),
			Answer:   fmt.Sprintf("%d", it.val),
			Explanation: fmt.Sprintf("By definition of the standard normal distribution, the %s is \\(%s = %d\\).",
				it.property, it.symbol, it.val),
		}
	}
}

// ----- stat.infer.* (inference) -----

type confidenceGen struct{}

func (g *confidenceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch rand.Intn(5) {
	case 0:
		type marginItem struct {
			n     int
			sigma int
			z     int
			e     int
		}
		items := []marginItem{
			{25, 15, 2, 6},
			{36, 24, 2, 8},
			{64, 40, 2, 10},
			{100, 30, 3, 9},
			{100, 50, 2, 10},
			{144, 36, 2, 6},
			{144, 48, 3, 12},
			{225, 45, 2, 6},
		}
		it := items[rand.Intn(len(items))]
		se := it.e / it.z
		return generator.Problem{
			Question: fmt.Sprintf("A random sample of size \\(n = %d\\) is drawn from a population with standard deviation \\(\\sigma = %d\\). Using critical value \\(z = %d\\), calculate the margin of error \\(E\\).", it.n, it.sigma, it.z),
			Answer:   fmt.Sprintf("%d", it.e),
			Explanation: fmt.Sprintf("The margin of error for a confidence interval is \\(E = z \\cdot \\frac{\\sigma}{\\sqrt{n}} = %d \\cdot \\frac{%d}{\\sqrt{%d}} = %d \\cdot %d = %d\\).",
				it.z, it.sigma, it.n, it.z, se, it.e),
		}
	case 1:
		type sampleSizeItem struct {
			sigma int
			z     int
			e     int
			ratio int
			n     int
		}
		items := []sampleSizeItem{
			{10, 2, 4, 5, 25},
			{15, 2, 5, 6, 36},
			{20, 2, 5, 8, 64},
			{12, 2, 3, 8, 64},
			{30, 2, 6, 10, 100},
			{16, 2, 4, 8, 64},
			{24, 3, 6, 12, 144},
			{25, 2, 5, 10, 100},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("A researcher requires a margin of error of \\(E = %d\\) with critical value \\(z = %d\\) from a population with standard deviation \\(\\sigma = %d\\). Determine the required sample size \\(n\\).", it.e, it.z, it.sigma),
			Answer:   fmt.Sprintf("%d", it.n),
			Explanation: fmt.Sprintf("The required sample size formula is \\(n = \\left(\\frac{z\\sigma}{E}\\right)^2 = \\left(\\frac{%d \\cdot %d}{%d}\\right)^2 = %d^2 = %d\\).",
				it.z, it.sigma, it.e, it.ratio, it.n),
		}
	case 2:
		type midpointItem struct {
			lower int
			upper int
			mean  int
		}
		items := []midpointItem{
			{64, 76, 70},
			{42, 58, 50},
			{115, 145, 130},
			{82, 98, 90},
			{23, 37, 30},
			{180, 220, 200},
			{75, 105, 90},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("A confidence interval for the population mean is computed as \\([%d, %d]\\). Find the sample mean point estimate \\(\\bar{x}\\).", it.lower, it.upper),
			Answer:   fmt.Sprintf("%d", it.mean),
			Explanation: fmt.Sprintf("The sample mean is the midpoint of the confidence interval: \\(\\bar{x} = \\frac{L + U}{2} = \\frac{%d + %d}{2} = \\frac{%d}{2} = %d\\).",
				it.lower, it.upper, it.lower+it.upper, it.mean),
		}
	case 3:
		type widthItem struct {
			lower  int
			upper  int
			margin int
		}
		items := []widthItem{
			{38, 62, 12},
			{85, 105, 10},
			{14, 26, 6},
			{55, 75, 10},
			{120, 150, 15},
			{63, 81, 9},
			{210, 250, 20},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("A symmetric confidence interval for a population parameter is given by \\([%d, %d]\\). What is the margin of error \\(E\\)?", it.lower, it.upper),
			Answer:   fmt.Sprintf("%d", it.margin),
			Explanation: fmt.Sprintf("The margin of error is half the interval width: \\(E = \\frac{U - L}{2} = \\frac{%d - %d}{2} = \\frac{%d}{2} = %d\\).",
				it.upper, it.lower, it.upper-it.lower, it.margin),
		}
	default:
		type totalWidthItem struct {
			e int
			w int
		}
		items := []totalWidthItem{
			{7, 14},
			{12, 24},
			{9, 18},
			{15, 30},
			{8, 16},
			{22, 44},
			{16, 32},
			{25, 50},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("A survey report establishes a confidence interval with a margin of error \\(E = %d\\). What is the total width \\(W\\) of this confidence interval?", it.e),
			Answer:   fmt.Sprintf("%d", it.w),
			Explanation: fmt.Sprintf("The total width of a confidence interval spanning \\([\\bar{x} - E, \\bar{x} + E]\\) is \\(W = 2E = 2 \\cdot %d = %d\\).",
				it.e, it.w),
		}
	}
}

type samplingGen struct{}

func (g *samplingGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch rand.Intn(5) {
	case 0:
		type seItem struct {
			n     int
			sigma int
			se    int
		}
		items := []seItem{
			{16, 20, 5},
			{25, 35, 7},
			{36, 42, 7},
			{64, 48, 6},
			{100, 80, 8},
			{144, 72, 6},
			{25, 15, 3},
			{36, 18, 3},
		}
		it := items[rand.Intn(len(items))]
		sqrtN := it.sigma / it.se
		return generator.Problem{
			Question: fmt.Sprintf("A population has standard deviation \\(\\sigma = %d\\). If a random sample of size \\(n = %d\\) is taken, calculate the standard error of the mean \\(\\sigma_{\\bar{x}}\\).", it.sigma, it.n),
			Answer:   fmt.Sprintf("%d", it.se),
			Explanation: fmt.Sprintf("The standard error of the mean is \\(\\sigma_{\\bar{x}} = \\frac{\\sigma}{\\sqrt{n}} = \\frac{%d}{\\sqrt{%d}} = \\frac{%d}{%d} = %d\\).",
				it.sigma, it.n, it.sigma, sqrtN, it.se),
		}
	case 1:
		type varItem struct {
			varPop int
			n      int
			v      int
		}
		items := []varItem{
			{120, 5, 24},
			{200, 8, 25},
			{180, 12, 15},
			{90, 6, 15},
			{144, 9, 16},
			{300, 15, 20},
			{160, 10, 16},
			{240, 16, 15},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("A random variable has population variance \\(\\sigma^2 = %d\\). A random sample of size \\(n = %d\\) is drawn. What is the variance of the sample mean \\(\\operatorname{Var}(\\bar{X})\\)?", it.varPop, it.n),
			Answer:   fmt.Sprintf("%d", it.v),
			Explanation: fmt.Sprintf("The variance of the sample mean is \\(\\operatorname{Var}(\\bar{X}) = \\frac{\\sigma^2}{n} = \\frac{%d}{%d} = %d\\).",
				it.varPop, it.n, it.v),
		}
	case 2:
		type meanItem struct {
			mu int
			n  int
		}
		items := []meanItem{
			{45, 25},
			{72, 36},
			{110, 50},
			{85, 100},
			{28, 16},
			{135, 64},
			{54, 40},
			{96, 60},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("A population has mean \\(\\mu = %d\\). If random samples of size \\(n = %d\\) are repeatedly collected, what is the expected value of the sample mean \\(E[\\bar{X}]\\)?", it.mu, it.n),
			Answer:   fmt.Sprintf("%d", it.mu),
			Explanation: fmt.Sprintf("The sample mean is an unbiased estimator of the population mean, so \\(E[\\bar{X}] = \\mu = %d\\).",
				it.mu),
		}
	case 3:
		type sampleSizeReqItem struct {
			sigma int
			se    int
			n     int
		}
		items := []sampleSizeReqItem{
			{24, 4, 36},
			{30, 5, 36},
			{40, 5, 64},
			{50, 5, 100},
			{36, 4, 81},
			{48, 6, 64},
			{60, 5, 144},
			{42, 6, 49},
		}
		it := items[rand.Intn(len(items))]
		ratio := it.sigma / it.se
		return generator.Problem{
			Question: fmt.Sprintf("A population has standard deviation \\(\\sigma = %d\\). What sample size \\(n\\) is needed so that the standard error of the mean equals \\(\\sigma_{\\bar{x}} = %d\\)?", it.sigma, it.se),
			Answer:   fmt.Sprintf("%d", it.n),
			Explanation: fmt.Sprintf("Solving \\(\\sigma_{\\bar{x}} = \\frac{\\sigma}{\\sqrt{n}}\\) gives \\(n = \\left(\\frac{\\sigma}{\\sigma_{\\bar{x}}}\\right)^2 = \\left(\\frac{%d}{%d}\\right)^2 = %d^2 = %d\\).",
				it.sigma, it.se, ratio, it.n),
		}
	default:
		type totalItem struct {
			n   int
			mu  int
			sum int
		}
		items := []totalItem{
			{8, 15, 120},
			{12, 20, 240},
			{10, 35, 350},
			{6, 45, 270},
			{15, 14, 210},
			{9, 25, 225},
			{16, 30, 480},
			{7, 50, 350},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("A sample of size \\(n = %d\\) is taken from a distribution with mean \\(\\mu = %d\\). Find the expected value of the sample total \\(E\\left[\\sum_{i=1}^n X_i\\right]\\).", it.n, it.mu),
			Answer:   fmt.Sprintf("%d", it.sum),
			Explanation: fmt.Sprintf("By linearity of expectation, the expected total of observations is \\(E\\left[\\sum_{i=1}^n X_i\\right] = n\\mu = %d \\cdot %d = %d\\).",
				it.n, it.mu, it.sum),
		}
	}
}

// ----- stat.prob.* (additional probability) -----

type bayesGen struct{}

func (g *bayesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch rand.Intn(5) {
	case 0:
		type postItem struct {
			pA      int
			pBA     int
			pBAc    int
			postPct int
		}
		items := []postItem{
			{50, 80, 20, 80},
			{50, 60, 20, 75},
			{50, 90, 30, 75},
			{40, 90, 40, 60},
			{25, 60, 20, 50},
			{10, 90, 10, 50},
			{50, 60, 40, 60},
			{25, 90, 10, 75},
		}
		it := items[rand.Intn(len(items))]
		num := float64(it.pBA*it.pA) / 10000.0
		den := num + float64(it.pBAc*(100-it.pA))/10000.0
		return generator.Problem{
			Question: fmt.Sprintf("In a binary hypothesis test, event \\(A\\) has prior probability \\(P(A) = %d\\%%\\). Evidence \\(B\\) occurs with likelihood \\(P(B \\mid A) = %d\\%%\\) and false alarm probability \\(P(B \\mid A^c) = %d\\%%\\). Compute the posterior probability \\(P(A \\mid B)\\) as an integer percentage.", it.pA, it.pBA, it.pBAc),
			Answer:   fmt.Sprintf("%d", it.postPct),
			Explanation: fmt.Sprintf("By Bayes theorem, \\(P(A \\mid B) = \\frac{P(B \\mid A)P(A)}{P(B \\mid A)P(A) + P(B \\mid A^c)P(A^c)} = \\frac{%.3f}{%.3f} = %d\\%%\\).",
				num, den, it.postPct),
		}
	case 1:
		type totalItem struct {
			pA1      int
			pA2      int
			pBA1     int
			pBA2     int
			totalPct int
		}
		items := []totalItem{
			{60, 40, 30, 80, 50},
			{70, 30, 40, 90, 55},
			{50, 50, 24, 66, 45},
			{80, 20, 25, 75, 35},
			{40, 60, 50, 20, 32},
			{30, 70, 80, 40, 52},
			{50, 50, 36, 84, 60},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(A_1\\) and \\(A_2\\) form a partition of the sample space with \\(P(A_1) = %d\\%%\\) and \\(P(A_2) = %d\\%%\\). If \\(P(B \\mid A_1) = %d\\%%\\) and \\(P(B \\mid A_2) = %d\\%%\\), calculate the total probability \\(P(B)\\) as an integer percentage.", it.pA1, it.pA2, it.pBA1, it.pBA2),
			Answer:   fmt.Sprintf("%d", it.totalPct),
			Explanation: fmt.Sprintf("By the law of total probability, \\(P(B) = P(B \\mid A_1)P(A_1) + P(B \\mid A_2)P(A_2) = (%d\\%% \\cdot %d\\%%) + (%d\\%% \\cdot %d\\%%) = %d\\%%\\).",
				it.pBA1, it.pA1, it.pBA2, it.pA2, it.totalPct),
		}
	case 2:
		type oddsItem struct {
			priorOdds int
			lr        int
			postOdds  int
		}
		items := []oddsItem{
			{2, 3, 6},
			{3, 4, 12},
			{1, 5, 5},
			{4, 2, 8},
			{5, 3, 15},
			{2, 7, 14},
			{6, 3, 18},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("A hypothesis has prior odds of %d to 1 in favor of \\(A\\). New observation \\(B\\) produces a likelihood ratio \\(\\frac{P(B \\mid A)}{P(B \\mid A^c)} = %d\\). What are the updated posterior odds in favor of \\(A\\)?", it.priorOdds, it.lr),
			Answer:   fmt.Sprintf("%d", it.postOdds),
			Explanation: fmt.Sprintf("In odds form, Bayes rule states \\(\\text{Posterior Odds} = \\text{Prior Odds} \\cdot \\text{Likelihood Ratio} = %d \\cdot %d = %d\\).",
				it.priorOdds, it.lr, it.postOdds),
		}
	case 3:
		type fpItem struct {
			n    int
			spec int
			fp   int
		}
		items := []fpItem{
			{2000, 95, 100},
			{5000, 98, 100},
			{1000, 96, 40},
			{4000, 97, 120},
			{3000, 99, 30},
			{8000, 95, 400},
			{6000, 98, 120},
			{10000, 99, 100},
		}
		it := items[rand.Intn(len(items))]
		fpRate := 100 - it.spec
		return generator.Problem{
			Question: fmt.Sprintf("A diagnostic screening test has specificity \\(%d\\%%\\). In a control group of %d healthy individuals, how many false positive results are expected?", it.spec, it.n),
			Answer:   fmt.Sprintf("%d", it.fp),
			Explanation: fmt.Sprintf("The false positive rate is \\(100\\%% - \\text{specificity} = 100\\%% - %d\\%% = %d\\%%\\). For %d healthy individuals, expected false positives are \\(%d \\cdot %d\\%% = %d\\).",
				it.spec, fpRate, it.n, it.n, fpRate, it.fp),
		}
	default:
		type compItem struct {
			k int
		}
		items := []compItem{
			{75},
			{80},
			{60},
			{70},
			{85},
			{92},
			{65},
			{88},
		}
		it := items[rand.Intn(len(items))]
		comp := 100 - it.k
		return generator.Problem{
			Question: fmt.Sprintf("After conditioning on observation \\(B\\), the posterior probability is \\(P(A \\mid B) = %d\\%%\\). What is the probability of the complementary hypothesis \\(P(A^c \\mid B)\\) as an integer percentage?", it.k),
			Answer:   fmt.Sprintf("%d", comp),
			Explanation: fmt.Sprintf("Conditional probabilities satisfy the complement rule: \\(P(A^c \\mid B) = 100\\%% - P(A \\mid B) = 100\\%% - %d\\%% = %d\\%%\\).",
				it.k, comp),
		}
	}
}

type continuousRVGen struct{}

func (g *continuousRVGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch rand.Intn(5) {
	case 0:
		type pointItem struct {
			c int
		}
		items := []pointItem{
			{3},
			{7},
			{12},
			{25},
			{40},
			{65},
			{80},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X\\) be a continuous random variable with probability density function \\(f(x)\\). What is the exact value of \\(P(X = %d)\\)?", it.c),
			Answer:   "0",
			Explanation: fmt.Sprintf("For every continuous random variable, the probability of taking any single exact point is zero: \\(P(X = %d) = \\int_{%d}^{%d} f(x)\\,dx = 0\\).",
				it.c, it.c, it.c),
		}
	case 1:
		type triItem struct {
			b   int
			b2  int
			x0  int
			pct int
		}
		items := []triItem{
			{10, 100, 5, 25},
			{10, 100, 6, 36},
			{10, 100, 7, 49},
			{10, 100, 8, 64},
			{10, 100, 9, 81},
			{20, 400, 10, 25},
			{20, 400, 14, 49},
			{20, 400, 16, 64},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("A continuous random variable \\(X\\) has probability density function \\(f(x) = \\frac{2x}{%d}\\) on the interval \\([0, %d]\\). Calculate the cumulative probability \\(P(X \\le %d)\\) as an integer percentage.", it.b2, it.b, it.x0),
			Answer:   fmt.Sprintf("%d", it.pct),
			Explanation: fmt.Sprintf("Integrating the density function yields \\(P(X \\le %d) = \\int_{0}^{%d} \\frac{2x}{%d}\\,dx = \\left[ \\frac{x^2}{%d} \\right]_{0}^{%d} = \\frac{%d}{%d} = %d\\%%\\).",
				it.x0, it.x0, it.b2, it.b2, it.x0, it.x0*it.x0, it.b2, it.pct),
		}
	case 2:
		type cdfDiffItem struct {
			a  int
			b  int
			fa int
			fb int
		}
		items := []cdfDiffItem{
			{2, 6, 25, 75},
			{1, 5, 15, 80},
			{3, 8, 30, 90},
			{2, 7, 20, 85},
			{0, 4, 10, 70},
			{4, 9, 35, 95},
			{1, 6, 18, 68},
			{2, 8, 22, 72},
		}
		it := items[rand.Intn(len(items))]
		diff := it.fb - it.fa
		return generator.Problem{
			Question: fmt.Sprintf("A continuous random variable \\(X\\) has cumulative distribution function \\(F(x)\\). If \\(F(%d) = %d\\%%\\) and \\(F(%d) = %d\\%%\\), what is \\(P(%d \\le X \\le %d)\\) as an integer percentage?", it.a, it.fa, it.b, it.fb, it.a, it.b),
			Answer:   fmt.Sprintf("%d", diff),
			Explanation: fmt.Sprintf("By the fundamental property of cumulative distribution functions, \\(P(%d \\le X \\le %d) = F(%d) - F(%d) = %d\\%% - %d\\%% = %d\\%%\\).",
				it.a, it.b, it.b, it.a, it.fb, it.fa, diff),
		}
	case 3:
		type tailItem struct {
			x   int
			pct int
		}
		items := []tailItem{
			{5, 60},
			{8, 75},
			{10, 85},
			{4, 70},
			{12, 90},
			{7, 65},
			{9, 82},
			{15, 95},
		}
		it := items[rand.Intn(len(items))]
		comp := 100 - it.pct
		return generator.Problem{
			Question: fmt.Sprintf("For a continuous random variable \\(X\\) with CDF \\(F(x)\\), the probability \\(P(X \\le %d) = %d\\%%\\). Find the tail probability \\(P(X > %d)\\) as an integer percentage.", it.x, it.pct, it.x),
			Answer:   fmt.Sprintf("%d", comp),
			Explanation: fmt.Sprintf("The complement rule for cumulative distributions gives \\(P(X > %d) = 100\\%% - P(X \\le %d) = 100\\%% - %d\\%% = %d\\%%\\).",
				it.x, it.x, it.pct, comp),
		}
	default:
		type unifItem struct {
			a   int
			b   int
			c   int
			d   int
			pct int
		}
		items := []unifItem{
			{0, 50, 10, 30, 40},
			{0, 40, 8, 28, 50},
			{10, 60, 25, 40, 30},
			{0, 100, 25, 75, 50},
			{20, 70, 35, 55, 40},
			{0, 25, 5, 15, 40},
			{10, 50, 18, 38, 50},
			{0, 20, 6, 14, 40},
		}
		it := items[rand.Intn(len(items))]
		lenSub := it.d - it.c
		lenTot := it.b - it.a
		return generator.Problem{
			Question: fmt.Sprintf("A continuous variable \\(X\\) is uniformly distributed over \\([%d, %d]\\). What is the probability \\(P(%d \\le X \\le %d)\\) as an integer percentage?", it.a, it.b, it.c, it.d),
			Answer:   fmt.Sprintf("%d", it.pct),
			Explanation: fmt.Sprintf("For a continuous uniform distribution, \\(P(%d \\le X \\le %d) = \\frac{%d - %d}{%d - %d} = \\frac{%d}{%d} = %d\\%%\\).",
				it.c, it.d, it.d, it.c, it.b, it.a, lenSub, lenTot, it.pct),
		}
	}
}

type discreteRVGen struct{}

func (g *discreteRVGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch rand.Intn(5) {
	case 0:
		type missingItem struct {
			p1 int
			p2 int
			p3 int
			p4 int
		}
		items := []missingItem{
			{15, 25, 35, 25},
			{20, 30, 40, 10},
			{10, 45, 25, 20},
			{22, 18, 30, 30},
			{12, 28, 45, 15},
			{35, 15, 35, 15},
			{16, 24, 20, 40},
			{25, 25, 30, 20},
		}
		it := items[rand.Intn(len(items))]
		sum3 := it.p1 + it.p2 + it.p3
		return generator.Problem{
			Question: fmt.Sprintf("A discrete random variable \\(X\\) has support \\(\\{1, 2, 3, 4\\}\\). If \\(P(X = 1) = %d\\%%\\), \\(P(X = 2) = %d\\%%\\) and \\(P(X = 3) = %d\\%%\\), determine \\(P(X = 4)\\) as an integer percentage.", it.p1, it.p2, it.p3),
			Answer:   fmt.Sprintf("%d", it.p4),
			Explanation: fmt.Sprintf("Probabilities of a discrete distribution must sum to 100\\%%: \\(P(X = 4) = 100\\%% - (%d\\%% + %d\\%% + %d\\%%) = 100\\%% - %d\\%% = %d\\%%\\).",
				it.p1, it.p2, it.p3, sum3, it.p4),
		}
	case 1:
		type cdfItem struct {
			p1  int
			p2  int
			p3  int
			p4  int
			k   int
			cdf int
		}
		items := []cdfItem{
			{20, 35, 25, 20, 2, 55},
			{15, 25, 40, 20, 2, 40},
			{30, 30, 20, 20, 3, 80},
			{10, 40, 30, 20, 3, 80},
			{25, 35, 15, 25, 2, 60},
			{18, 32, 26, 24, 2, 50},
			{12, 28, 35, 25, 3, 75},
			{22, 38, 25, 15, 2, 60},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("A discrete random variable \\(X\\) has probabilities \\(P(X = 1) = %d\\%%\\), \\(P(X = 2) = %d\\%%\\), \\(P(X = 3) = %d\\%%\\) and \\(P(X = 4) = %d\\%%\\). Compute the cumulative probability \\(F(%d) = P(X \\le %d)\\) as an integer percentage.", it.p1, it.p2, it.p3, it.p4, it.k, it.k),
			Answer:   fmt.Sprintf("%d", it.cdf),
			Explanation: fmt.Sprintf("The CDF is the sum of probabilities up to %d: \\(F(%d) = P(X \\le %d) = %d\\%%\\).",
				it.k, it.k, it.k, it.cdf),
		}
	case 2:
		type tailItem struct {
			p3 int
			p4 int
		}
		items := []tailItem{
			{30, 15},
			{25, 35},
			{40, 20},
			{35, 25},
			{20, 15},
			{50, 10},
			{15, 25},
			{32, 18},
		}
		it := items[rand.Intn(len(items))]
		tail := it.p3 + it.p4
		return generator.Problem{
			Question: fmt.Sprintf("A random variable takes values in \\(\\{1, 2, 3, 4\\}\\). If \\(P(X = 3) = %d\\%%\\) and \\(P(X = 4) = %d\\%%\\), calculate \\(P(X \\ge 3)\\) as an integer percentage.", it.p3, it.p4),
			Answer:   fmt.Sprintf("%d", tail),
			Explanation: fmt.Sprintf("For discrete values, \\(P(X \\ge 3) = P(X = 3) + P(X = 4) = %d\\%% + %d\\%% = %d\\%%\\).",
				it.p3, it.p4, tail),
		}
	case 3:
		type weightItem struct {
			w1      int
			w2      int
			w3      int
			targetX int
			targetW int
			ansPct  int
		}
		items := []weightItem{
			{2, 3, 5, 3, 5, 50},
			{3, 7, 10, 2, 7, 35},
			{4, 9, 12, 3, 12, 48},
			{1, 3, 6, 2, 3, 30},
			{5, 15, 30, 2, 15, 30},
			{7, 8, 10, 1, 7, 28},
			{2, 5, 13, 2, 5, 25},
			{6, 14, 30, 3, 30, 60},
		}
		it := items[rand.Intn(len(items))]
		totalW := it.w1 + it.w2 + it.w3
		return generator.Problem{
			Question: fmt.Sprintf("A discrete random variable \\(X \\in \\{1, 2, 3\\}\\) has probabilities proportional to weights \\(w_1 = %d\\), \\(w_2 = %d\\) and \\(w_3 = %d\\). Calculate \\(P(X = %d)\\) as an integer percentage.", it.w1, it.w2, it.w3, it.targetX),
			Answer:   fmt.Sprintf("%d", it.ansPct),
			Explanation: fmt.Sprintf("The total weight is \\(W = %d + %d + %d = %d\\). The probability is \\(P(X = %d) = \\frac{%d}{%d} = %d\\%%\\).",
				it.w1, it.w2, it.w3, totalW, it.targetX, it.targetW, totalW, it.ansPct),
		}
	default:
		type unifItem struct {
			n   int
			m   int
			pct int
		}
		items := []unifItem{
			{10, 3, 30},
			{10, 7, 70},
			{20, 9, 45},
			{20, 13, 65},
			{25, 11, 44},
			{25, 16, 64},
			{50, 21, 42},
			{50, 37, 74},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("A fair integer spinner is equally likely to land on any integer from 1 to %d. What is the probability that the result lands within a designated subset of %d distinct outcomes as an integer percentage?", it.n, it.m),
			Answer:   fmt.Sprintf("%d", it.pct),
			Explanation: fmt.Sprintf("In a discrete uniform distribution on %d elements, any subset of size %d has probability \\(\\frac{%d}{%d} = %d\\%%\\).",
				it.n, it.m, it.m, it.n, it.pct),
		}
	}
}

type expectedValueGen struct{}

func (g *expectedValueGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch rand.Intn(5) {
	case 0:
		type linItem struct {
			ex  int
			a   int
			b   int
			ans int
		}
		items := []linItem{
			{4, 3, 5, 17},
			{6, 2, 7, 19},
			{5, 4, 3, 23},
			{8, 3, 6, 30},
			{7, 5, 4, 39},
			{9, 2, 8, 26},
			{10, 3, 12, 42},
			{12, 4, 5, 53},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X\\) be a random variable with expected value \\(E[X] = %d\\). Calculate \\(E[%dX + %d]\\).", it.ex, it.a, it.b),
			Answer:   fmt.Sprintf("%d", it.ans),
			Explanation: fmt.Sprintf("By linearity of expectation, \\(E[%dX + %d] = %dE[X] + %d = (%d \\cdot %d) + %d = %d\\).",
				it.a, it.b, it.a, it.b, it.a, it.ex, it.b, it.ans),
		}
	case 1:
		type combItem struct {
			ex  int
			ey  int
			a   int
			b   int
			ans int
		}
		items := []combItem{
			{5, 8, 2, 3, 34},
			{4, 7, 3, 2, 26},
			{6, 9, 4, 3, 51},
			{10, 4, 2, 5, 40},
			{8, 6, 3, 4, 48},
			{7, 12, 5, 2, 59},
			{9, 5, 2, 6, 48},
			{3, 11, 4, 3, 45},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("Random variables \\(X\\) and \\(Y\\) have expected values \\(E[X] = %d\\) and \\(E[Y] = %d\\). Compute \\(E[%dX + %dY]\\).", it.ex, it.ey, it.a, it.b),
			Answer:   fmt.Sprintf("%d", it.ans),
			Explanation: fmt.Sprintf("By linearity of expectation, \\(E[%dX + %dY] = %dE[X] + %dE[Y] = (%d \\cdot %d) + (%d \\cdot %d) = %d + %d = %d\\).",
				it.a, it.b, it.a, it.b, it.a, it.ex, it.b, it.ey, it.a*it.ex, it.b*it.ey, it.ans),
		}
	case 2:
		type momentItem struct {
			ex  int
			vx  int
			ex2 int
		}
		items := []momentItem{
			{3, 7, 16},
			{4, 9, 25},
			{5, 11, 36},
			{6, 13, 49},
			{7, 15, 64},
			{8, 17, 81},
			{10, 25, 125},
			{12, 20, 164},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("A random variable \\(X\\) has expected value \\(E[X] = %d\\) and variance \\(\\operatorname{Var}(X) = %d\\). Calculate the second moment \\(E[X^2]\\).", it.ex, it.vx),
			Answer:   fmt.Sprintf("%d", it.ex2),
			Explanation: fmt.Sprintf("Using the variance formula \\(\\operatorname{Var}(X) = E[X^2] - (E[X])^2\\), we find \\(E[X^2] = \\operatorname{Var}(X) + (E[X])^2 = %d + %d^2 = %d + %d = %d\\).",
				it.vx, it.ex, it.vx, it.ex*it.ex, it.ex2),
		}
	case 3:
		type discreteItem struct {
			x1   int
			x2   int
			x3   int
			p1   float64
			p2   float64
			p3   float64
			mean int
		}
		items := []discreteItem{
			{2, 6, 10, 0.25, 0.50, 0.25, 6},
			{1, 5, 9, 0.20, 0.60, 0.20, 5},
			{3, 7, 11, 0.30, 0.40, 0.30, 7},
			{4, 8, 16, 0.25, 0.50, 0.25, 9},
			{0, 5, 10, 0.30, 0.40, 0.30, 5},
			{1, 6, 11, 0.40, 0.20, 0.40, 6},
			{2, 7, 12, 0.10, 0.80, 0.10, 7},
			{4, 10, 16, 0.30, 0.40, 0.30, 10},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("A discrete random variable \\(X\\) takes values %d, %d and %d with probabilities %.2f, %.2f and %.2f respectively. Calculate the expected value \\(E[X]\\).", it.x1, it.x2, it.x3, it.p1, it.p2, it.p3),
			Answer:   fmt.Sprintf("%d", it.mean),
			Explanation: fmt.Sprintf("The expected value is the probability-weighted sum \\(E[X] = \\sum x_i P(X = x_i) = (%d \\cdot %.2f) + (%d \\cdot %.2f) + (%d \\cdot %.2f) = %d\\).",
				it.x1, it.p1, it.x2, it.p2, it.x3, it.p3, it.mean),
		}
	default:
		type prodItem struct {
			ex   int
			ey   int
			prod int
		}
		items := []prodItem{
			{4, 7, 28},
			{5, 9, 45},
			{6, 8, 48},
			{3, 12, 36},
			{7, 8, 56},
			{9, 6, 54},
			{11, 5, 55},
			{12, 7, 84},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X\\) and \\(Y\\) be independent random variables with \\(E[X] = %d\\) and \\(E[Y] = %d\\). Find the expected value of their product \\(E[XY]\\).", it.ex, it.ey),
			Answer:   fmt.Sprintf("%d", it.prod),
			Explanation: fmt.Sprintf("For independent random variables, the expected value of the product is the product of their expectations: \\(E[XY] = E[X]E[Y] = %d \\cdot %d = %d\\).",
				it.ex, it.ey, it.prod),
		}
	}
}

type gaussianFuncGen struct{}

func (g *gaussianFuncGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch rand.Intn(5) {
	case 0:
		type peakItem struct {
			a     int
			mu    int
			denom int
		}
		items := []peakItem{
			{5, 12, 18},
			{8, 25, 32},
			{3, 7, 8},
			{10, 40, 50},
			{4, 15, 18},
			{6, 30, 72},
			{9, 18, 50},
			{7, 22, 32},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("The Gaussian curve is given by \\(f(x) = %d \\exp\\left(-\\frac{(x - %d)^2}{%d}\\right)\\). At what value of \\(x\\) does the function reach its maximum peak?", it.a, it.mu, it.denom),
			Answer:   fmt.Sprintf("%d", it.mu),
			Explanation: fmt.Sprintf("A Gaussian function attains its global maximum at the center of symmetry \\(x = \\mu = %d\\).",
				it.mu),
		}
	case 1:
		type heightItem struct {
			a     int
			mu    int
			denom int
		}
		items := []heightItem{
			{24, 10, 18},
			{45, 15, 32},
			{80, 20, 50},
			{15, 8, 8},
			{60, 25, 72},
			{36, 12, 18},
			{50, 30, 98},
			{90, 16, 50},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("A bell-shaped Gaussian function is defined as \\(f(x) = %d \\exp\\left(-\\frac{(x - %d)^2}{%d}\\right)\\). What is the maximum peak height \\(f(%d)\\)?", it.a, it.mu, it.denom, it.mu),
			Answer:   fmt.Sprintf("%d", it.a),
			Explanation: fmt.Sprintf("Evaluating at the mean gives \\(f(%d) = %d \\exp(0) = %d \\cdot 1 = %d\\).",
				it.mu, it.a, it.a, it.a),
		}
	case 2:
		type convItem struct {
			v1 int
			v2 int
		}
		items := []convItem{
			{12, 13},
			{20, 16},
			{35, 29},
			{48, 52},
			{18, 31},
			{27, 54},
			{60, 84},
			{45, 55},
		}
		it := items[rand.Intn(len(items))]
		sumV := it.v1 + it.v2
		return generator.Problem{
			Question: fmt.Sprintf("Two independent Gaussian random variables have variances \\(\\sigma_1^2 = %d\\) and \\(\\sigma_2^2 = %d\\). Under convolution, what is the variance of their sum \\(\\sigma_{\\text{sum}}^2\\)?", it.v1, it.v2),
			Answer:   fmt.Sprintf("%d", sumV),
			Explanation: fmt.Sprintf("The convolution of two independent Gaussians yields a Gaussian with variance equal to the sum of variances: \\(\\sigma_{\\text{sum}}^2 = \\sigma_1^2 + \\sigma_2^2 = %d + %d = %d\\).",
				it.v1, it.v2, sumV),
		}
	case 3:
		type inflItem struct {
			sigma int
		}
		items := []inflItem{
			{3},
			{5},
			{8},
			{12},
			{15},
			{7},
			{9},
			{14},
		}
		it := items[rand.Intn(len(items))]
		dist := 2 * it.sigma
		return generator.Problem{
			Question: fmt.Sprintf("A Gaussian distribution has standard deviation \\(\\sigma = %d\\). What is the distance between its two inflection points at \\(x = \\mu - \\sigma\\) and \\(x = \\mu + \\sigma\\)?", it.sigma),
			Answer:   fmt.Sprintf("%d", dist),
			Explanation: fmt.Sprintf("The distance between the inflection points is \\((\\mu + \\sigma) - (\\mu - \\sigma) = 2\\sigma = 2 \\cdot %d = %d\\).",
				it.sigma, dist),
		}
	default:
		type denomItem struct {
			mu    int
			d     int
			sigma int
		}
		items := []denomItem{
			{5, 8, 2},
			{10, 18, 3},
			{15, 32, 4},
			{20, 50, 5},
			{25, 72, 6},
			{30, 98, 7},
			{35, 128, 8},
			{40, 200, 10},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("A Gaussian curve has the form \\(f(x) = C \\exp\\left(-\\frac{(x - %d)^2}{%d}\\right)\\), where the denominator represents \\(2\\sigma^2\\). Determine the standard deviation \\(\\sigma\\).", it.mu, it.d),
			Answer:   fmt.Sprintf("%d", it.sigma),
			Explanation: fmt.Sprintf("Equating \\(2\\sigma^2 = %d\\) gives \\(\\sigma^2 = %d\\), so \\(\\sigma = \\sqrt{%d} = %d\\).",
				it.d, it.d/2, it.d/2, it.sigma),
		}
	}
}
