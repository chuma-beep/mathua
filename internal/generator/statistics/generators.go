package statistics

import (
	"fmt"
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
}

type readTableGen struct{}

func (g *readTableGen) Generate(difficulty float64) generator.Problem {
	numRows := rand.Intn(2) + 3
	numCols := rand.Intn(2) + 2
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
			table[i][j] = rand.Intn(20) + 1
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
		values[i] = rand.Intn(30) + 5
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
	n := rand.Intn(3) + 4
	xs := make([]int, n)
	ys := make([]int, n)
	startX := rand.Intn(5) + 1
	startY := rand.Intn(20) + 5
	for i := 0; i < n; i++ {
		xs[i] = startX + i
		if i == 0 {
			ys[i] = startY
		} else {
			delta := rand.Intn(7) - 3
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
	n := rand.Intn(4) + 4
	vals := make([]int, n)
	sum := 0
	for i := 0; i < n; i++ {
		vals[i] = rand.Intn(40) + 1
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
	n := rand.Intn(5) + 5
	vals := make([]int, n)
	for i := 0; i < n; i++ {
		vals[i] = rand.Intn(50) + 1
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
	n := rand.Intn(5) + 5
	base := rand.Intn(20) + 1
	modeVal := base + rand.Intn(10)
	vals := make([]int, 0, n)
	modeCount := rand.Intn(3) + 2
	for i := 0; i < modeCount; i++ {
		vals = append(vals, modeVal)
	}
	otherCount := n - modeCount
	for i := 0; i < otherCount; i++ {
		v := rand.Intn(30) + 1
		for v == modeVal {
			v = rand.Intn(30) + 1
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
	n := rand.Intn(5) + 5
	vals := make([]int, n)
	minVal := rand.Intn(20) + 1
	maxVal := minVal + rand.Intn(30) + 10
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
	colors := []string{"red", "blue", "green", "yellow", "orange", "purple"}
	counts := make([]int, 3)
	total := 0
	colorNames := make([]string, 3)
	for i := 0; i < 3; i++ {
		colorNames[i] = colors[rand.Intn(len(colors))]
		counts[i] = rand.Intn(5) + 1
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
		Question:    fmt.Sprintf("A bag contains %s. What is P(%s)?", strings.Join(bagParts, ", "), colorNames[pick]),
		Answer:      fmt.Sprintf("%d/%d", fav/gcd, total/gcd),
		Explanation: fmt.Sprintf("P(%s) = %d/%d = %d/%d", colorNames[pick], fav, total, fav/gcd, total/gcd),
	}
}

type complementProbGen struct{}

func (g *complementProbGen) Generate(difficulty float64) generator.Problem {
	colors := []string{"red", "blue", "green", "yellow", "orange", "purple"}
	counts := make([]int, 3)
	total := 0
	colorNames := make([]string, 3)
	for i := 0; i < 3; i++ {
		colorNames[i] = colors[rand.Intn(len(colors))]
		counts[i] = rand.Intn(5) + 1
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
		Question:    fmt.Sprintf("A bag contains %s. What is P(not %s)?", strings.Join(bagParts, ", "), colorNames[pick]),
		Answer:      fmt.Sprintf("%d/%d", other/gcd, total/gcd),
		Explanation: fmt.Sprintf("P(not %s) = 1 - P(%s) = 1 - %d/%d = %d/%d", colorNames[pick], colorNames[pick], counts[pick], total, other/gcd, total/gcd),
	}
}

type compoundProbGen struct{}

func (g *compoundProbGen) Generate(difficulty float64) generator.Problem {
	colors := []string{"red", "blue", "green", "yellow"}
	counts := make([]int, 3)
	total := 0
	colorNames := make([]string, 3)
	for i := 0; i < 3; i++ {
		colorNames[i] = colors[rand.Intn(len(colors))]
		counts[i] = rand.Intn(4) + 2
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
		Question:    fmt.Sprintf("A bag contains %s. You draw one marble, then another without replacement. What is P(%s and then %s)?", strings.Join(bagParts, ", "), n1, n2),
		Answer:      fmt.Sprintf("%d/%d", num/gcd, den/gcd),
		Explanation: fmt.Sprintf("P(%s then %s) = %d/%d x %d/%d = %d/%d = %d/%d", n1, n2, c1, total, c2, total-1, num, den, num/gcd, den/gcd),
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
