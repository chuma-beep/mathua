package statistics

// Wave 7 statistics: data summaries, counting rules, further distributions,
// exact tests and regression diagnostics. Every answer is an integer and every
// operand is published in Facts so wave7_test.go can recompute it.

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave7(reg *generator.Registry) {
	kinds := map[string]string{
		"stat.data.two_way_table":        "two_way_table",
		"stat.data.stem_leaf":            "stem_leaf",
		"stat.data.cumulative_frequency": "cumulative_frequency",
		"stat.data.percentile":           "percentile",
		"stat.data.trimmed_mean":         "trimmed_mean",
		"stat.data.midrange":             "midrange",
		"stat.prob.union":                "union",
		"stat.prob.tree":                 "tree",
		"stat.prob.total_probability":    "total_probability",
		"stat.prob.permutations":         "permutations",
		"stat.prob.combinations":         "combinations",
		"stat.prob.odds":                 "odds",
		"stat.dist.weibull":              "weibull",
		"stat.dist.cauchy":               "cauchy",
		"stat.dist.triangular":           "triangular",
		"stat.dist.pareto":               "pareto",
		"stat.infer.fisher_exact":        "fisher_exact",
		"stat.infer.mcnemar":             "mcnemar",
		"stat.infer.sign_test":           "sign_test",
		"stat.reg.cooks_distance":        "cooks_distance",
		"stat.reg.polynomial":            "polynomial",
		"stat.reg.ridge":                 "ridge",
	}
	for id, kind := range kinds {
		reg.Register(id, &w7StatGen{kind: kind})
	}
}

type w7StatGen struct{ kind string }

func (g *w7StatGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "two_way_table":
		return w7TwoWayTable(kp)
	case "stem_leaf":
		return w7StemLeaf(kp)
	case "cumulative_frequency":
		return w7CumulativeFrequency(kp)
	case "percentile":
		return w7Percentile(kp)
	case "trimmed_mean":
		return w7TrimmedMean(kp)
	case "midrange":
		return w7Midrange(kp)
	case "union":
		return w7Union(kp)
	case "tree":
		return w7Tree(kp)
	case "total_probability":
		return w7TotalProbability(kp)
	case "permutations":
		return w7Permutations(kp)
	case "combinations":
		return w7Combinations(kp)
	case "odds":
		return w7Odds(kp)
	case "weibull":
		return w7Weibull(kp)
	case "cauchy":
		return w7Cauchy(kp)
	case "triangular":
		return w7Triangular(kp)
	case "pareto":
		return w7Pareto(kp)
	case "fisher_exact":
		return w7FisherExact(kp)
	case "mcnemar":
		return w7McNemar(kp)
	case "sign_test":
		return w7SignTest(kp)
	case "cooks_distance":
		return w7CooksDistance(kp)
	case "polynomial":
		return w7Polynomial(kp)
	case "ridge":
		return w7Ridge(kp)
	}
	return generator.Problem{Question: "What is 2 + 2? (enter a number)", Answer: "4", Explanation: "2 + 2 = 4."}
}

func w7S(o int, q, a, ex string, facts map[string]string) generator.Problem {
	if facts == nil {
		facts = map[string]string{}
	}
	facts["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: facts, KPOrdinal: o}
}

func w7I(m map[string]string, k string, v int) { m[k] = strconv.Itoa(v) }

func w7R(lo, hi int) int {
	if hi <= lo {
		return lo
	}
	return lo + rand.Intn(hi-lo+1)
}

func w7Join(xs []int) string {
	parts := make([]string, len(xs))
	for i, x := range xs {
		parts[i] = strconv.Itoa(x)
	}
	return strings.Join(parts, ",")
}

func w7TwoWayTable(kp int) generator.Problem {
	a, b := w7R(5, 20), w7R(5, 20)
	c, d := w7R(5, 20), w7R(5, 20)
	row1, row2 := a+b, c+d
	col1, col2 := a+c, b+d
	total := row1 + row2
	f := map[string]string{}
	for k, v := range map[string]int{"a": a, "b": b, "c": c, "d": d, "row1": row1, "row2": row2, "col1": col1, "col2": col2, "total": total} {
		w7I(f, k, v)
	}
	switch kp {
	case 1:
		return w7S(kp, fmt.Sprintf("A two-way table has a=%d and b=%d in the first row, and c=%d and d=%d in the second row. What is the first-row total? (enter a number)", a, b, c, d), strconv.Itoa(row1),
			fmt.Sprintf("The first-row total is a+b = %d+%d = %d.", a, b, row1), f)
	case 2:
		return w7S(kp, fmt.Sprintf("A two-way table has a=%d and b=%d in the first row, and c=%d and d=%d in the second row. What is the second-column total? (enter a number)", a, b, c, d), strconv.Itoa(col2),
			fmt.Sprintf("The second-column total is b+d = %d+%d = %d.", b, d, col2), f)
	default:
		return w7S(kp, fmt.Sprintf("A two-way table has a=%d and b=%d in the first row, and c=%d and d=%d in the second row. What is the grand total? (enter a number)", a, b, c, d), strconv.Itoa(total),
			fmt.Sprintf("The grand total is the sum of the four cells: %d+%d+%d+%d = %d.", a, b, c, d, total), f)
	}
}

func w7StemLeaf(kp int) generator.Problem {
	n := w7R(10, 20)
	vals := make([]int, n)
	for i := range vals {
		vals[i] = w7R(10, 99)
	}
	sort.Ints(vals)
	parts := []string{}
	for _, v := range vals {
		parts = append(parts, fmt.Sprintf("%d|%d", v/10, v%10))
	}
	f := map[string]string{"values": w7Join(vals), "stems": strings.Join(parts, " ")}
	w7I(f, "n", n)
	w7I(f, "min", vals[0])
	w7I(f, "max", vals[n-1])
	median := vals[n/2]
	w7I(f, "median", median)
	switch kp {
	case 1:
		return w7S(kp, fmt.Sprintf("A stem-and-leaf plot holds the values %s. How many values are there? (enter a number)", w7Join(vals)), strconv.Itoa(n),
			fmt.Sprintf("Count the leaves across all stems: there are %d values.", n), f)
	case 2:
		return w7S(kp, fmt.Sprintf("A stem-and-leaf plot holds the values %s. What is the range? (enter a number)", w7Join(vals)), strconv.Itoa(vals[n-1]-vals[0]),
			fmt.Sprintf("The range is the maximum minus the minimum: %d-%d = %d.", vals[n-1], vals[0], vals[n-1]-vals[0]), f)
	default:
		return w7S(kp, fmt.Sprintf("A stem-and-leaf plot holds the values %s. What is the median? (enter a number)", w7Join(vals)), strconv.Itoa(median),
			fmt.Sprintf("With %d ordered values the median is the value in position %d, which is %d.", n, n/2+1, median), f)
	}
}

func w7CumulativeFrequency(kp int) generator.Problem {
	k := w7R(4, 5)
	freq := make([]int, k)
	for i := range freq {
		freq[i] = w7R(2, 12)
	}
	cum := make([]int, k)
	run := 0
	for i, v := range freq {
		run += v
		cum[i] = run
	}
	pick := w7R(1, k-1)
	f := map[string]string{"frequencies": w7Join(freq), "cumulative": w7Join(cum)}
	for i, v := range freq {
		w7I(f, fmt.Sprintf("f%d", i+1), v)
	}
	for i, v := range cum {
		w7I(f, fmt.Sprintf("c%d", i+1), v)
	}
	w7I(f, "classes", k)
	w7I(f, "pick", pick)
	switch kp {
	case 1:
		return w7S(kp, fmt.Sprintf("The class frequencies are %s. What is the cumulative frequency after the first %d classes? (enter a number)", w7Join(freq), pick), strconv.Itoa(cum[pick-1]),
			fmt.Sprintf("Add the first %d frequencies: the cumulative total is %d.", pick, cum[pick-1]), f)
	case 2:
		return w7S(kp, fmt.Sprintf("The class frequencies are %s. What is the total frequency? (enter a number)", w7Join(freq)), strconv.Itoa(run),
			fmt.Sprintf("The total is the sum of every class frequency, which is %d.", run), f)
	default:
		before := 0
		if pick > 1 {
			before = cum[pick-2]
		}
		return w7S(kp, fmt.Sprintf("The cumulative frequencies after classes 1 to %d are %s. What is the frequency of class %d alone? (enter a number)", k, w7Join(cum), pick), strconv.Itoa(freq[pick-1]),
			fmt.Sprintf("Subtract the previous cumulative total: %d-%d = %d.", cum[pick-1], before, freq[pick-1]), f)
	}
}

func w7Percentile(kp int) generator.Problem {
	n := w7R(10, 20)
	vals := make([]int, n)
	for i := range vals {
		vals[i] = w7R(10, 99)
	}
	sort.Ints(vals)
	ps := []int{10, 25, 50, 75, 90}
	p := ps[rand.Intn(len(ps))]
	idx := (p*n + 99) / 100
	if idx < 1 {
		idx = 1
	}
	if idx > n {
		idx = n
	}
	value := vals[idx-1]
	f := map[string]string{"values": w7Join(vals)}
	w7I(f, "n", n)
	w7I(f, "p", p)
	w7I(f, "rank", idx)
	w7I(f, "value", value)
	switch kp {
	case 1:
		return w7S(kp, fmt.Sprintf("A sorted data set has %d values. Using the nearest-rank method, what is the position of the %dth percentile? (enter a number)", n, p), strconv.Itoa(idx),
			fmt.Sprintf("The nearest-rank position is ceil(%d/100 * %d) = %d.", p, n, idx), f)
	case 2:
		return w7S(kp, fmt.Sprintf("A sorted data set has %d values: %s. Using the nearest-rank method, what is the %dth percentile? (enter a number)", n, w7Join(vals), p), strconv.Itoa(value),
			fmt.Sprintf("The position is %d, so the %dth percentile is %d.", idx, p, value), f)
	default:
		return w7S(kp, fmt.Sprintf("A sorted data set has %d values: %s. Using the nearest-rank method, how many values lie below the %dth percentile? (enter a number)", n, w7Join(vals), p), strconv.Itoa(idx-1),
			fmt.Sprintf("The percentile sits at position %d, so %d values lie below it.", idx, idx-1), f)
	}
}

func w7TrimmedMean(kp int) generator.Problem {
	base := w7R(5, 40)
	vals := make([]int, 9)
	for i := range vals {
		vals[i] = base + 2*i
	}
	sumRetained := 5*base + 40
	sumRemoved := 4*base + 32
	mean := base + 8
	f := map[string]string{"values": w7Join(vals)}
	for k, v := range map[string]int{"base": base, "sum_retained": sumRetained, "sum_removed": sumRemoved, "mean": mean} {
		w7I(f, k, v)
	}
	switch kp {
	case 1:
		return w7S(kp, fmt.Sprintf("The values are %s. Remove the two smallest and the two largest. What is the sum of the five remaining values? (enter a number)", w7Join(vals)), strconv.Itoa(sumRetained),
			fmt.Sprintf("The five middle values sum to %d.", sumRetained), f)
	case 2:
		return w7S(kp, fmt.Sprintf("The values are %s. Remove the two smallest and the two largest. What is the sum of the four removed values? (enter a number)", w7Join(vals)), strconv.Itoa(sumRemoved),
			fmt.Sprintf("The two smallest plus the two largest sum to %d.", sumRemoved), f)
	default:
		return w7S(kp, fmt.Sprintf("The values are %s. Remove the two smallest and the two largest. What is the mean of the remaining values? (enter a number)", w7Join(vals)), strconv.Itoa(mean),
			fmt.Sprintf("The five remaining values sum to %d, so their mean is %d.", sumRetained, mean), f)
	}
}

func w7Midrange(kp int) generator.Problem {
	min := w7R(10, 60)
	max := min + 2*w7R(1, 20)
	mid := (min + max) / 2
	f := map[string]string{}
	for k, v := range map[string]int{"min": min, "max": max, "midrange": mid, "range": max - min} {
		w7I(f, k, v)
	}
	switch kp {
	case 1:
		return w7S(kp, fmt.Sprintf("A data set has minimum %d and maximum %d. What is its midrange? (enter a number)", min, max), strconv.Itoa(mid),
			fmt.Sprintf("The midrange is (%d+%d)/2 = %d.", min, max, mid), f)
	case 2:
		return w7S(kp, fmt.Sprintf("A data set has minimum %d and midrange %d. What is its maximum? (enter a number)", min, mid), strconv.Itoa(max),
			fmt.Sprintf("The maximum is 2*%d-%d = %d.", mid, min, max), f)
	default:
		return w7S(kp, fmt.Sprintf("A data set has minimum %d and maximum %d. What is its range? (enter a number)", min, max), strconv.Itoa(max-min),
			fmt.Sprintf("The range is %d-%d = %d.", max, min, max-min), f)
	}
}

func w7Union(kp int) generator.Problem {
	pa := w7R(20, 60)
	pb := w7R(20, 60)
	inter := w7R(5, minInt(pa, pb)-5)
	union := pa + pb - inter
	f := map[string]string{}
	for k, v := range map[string]int{"pa": pa, "pb": pb, "intersection": inter, "union": union} {
		w7I(f, k, v)
	}
	switch kp {
	case 1:
		return w7S(kp, fmt.Sprintf("P(A)=%d%% and P(B)=%d%%, with P(A and B)=%d%%. What is P(A or B), in percent? (enter a number)", pa, pb, inter), strconv.Itoa(union),
			fmt.Sprintf("The addition rule gives %d+%d-%d = %d percent.", pa, pb, inter, union), f)
	case 2:
		return w7S(kp, fmt.Sprintf("P(A)=%d%%, P(B)=%d%% and P(A or B)=%d%%. What is P(A and B), in percent? (enter a number)", pa, pb, union), strconv.Itoa(inter),
			fmt.Sprintf("Rearrange the addition rule: %d+%d-%d = %d percent.", pa, pb, union, inter), f)
	default:
		return w7S(kp, fmt.Sprintf("P(A)=%d%%, P(A and B)=%d%% and P(A or B)=%d%%. What is P(B), in percent? (enter a number)", pa, inter, union), strconv.Itoa(pb),
			fmt.Sprintf("Rearrange the addition rule: %d-%d+%d = %d percent.", union, pa, inter, pb), f)
	}
}

func w7Tree(kp int) generator.Problem {
	p := w7R(2, 9) * 10
	q := w7R(2, 9) * 10
	r := w7R(2, 9) * 10
	s := w7R(2, 9) * 10
	path1 := p * q / 100
	path2 := r * s / 100
	total := path1 + path2
	complement := 100 - path1
	f := map[string]string{}
	for k, v := range map[string]int{"p": p, "q": q, "r": r, "s": s, "path1": path1, "path2": path2, "total": total, "complement": complement} {
		w7I(f, k, v)
	}
	switch kp {
	case 1:
		return w7S(kp, fmt.Sprintf("A two-stage experiment takes a first branch with probability %d%% and then a second branch with probability %d%%. What is the probability of that path, in percent? (enter a number)", p, q), strconv.Itoa(path1),
			fmt.Sprintf("Multiply along the path: %d/100 * %d/100 = %d/100, or %d percent.", p, q, path1, path1), f)
	case 2:
		return w7S(kp, fmt.Sprintf("A tree has one path with probabilities %d%% and %d%%, and a second path with probabilities %d%% and %d%%. What is the combined probability of the two paths, in percent? (enter a number)", p, q, r, s), strconv.Itoa(total),
			fmt.Sprintf("The paths contribute %d and %d percent, so together they contribute %d percent.", path1, path2, total), f)
	default:
		return w7S(kp, fmt.Sprintf("A tree has one path with probabilities %d%% and %d%%. What is the probability of not taking that path, in percent? (enter a number)", p, q), strconv.Itoa(complement),
			fmt.Sprintf("The path has probability %d percent, so its complement is %d percent.", path1, complement), f)
	}
}

func w7TotalProbability(kp int) generator.Problem {
	p1 := w7R(1, 4) * 10
	p2 := w7R(1, 4) * 10
	p3 := w7R(1, 4) * 10
	q := w7R(2, 8) * 10
	sum := p1 + p2 + p3
	if sum > 100 {
		p3 = 100 - p1 - p2
		if p3 <= 0 {
			p1, p2, p3 = 20, 30, 50
		}
		sum = p1 + p2 + p3
	}
	if sum != 100 {
		p3 = 100 - p1 - p2
	}
	contrib1 := p1 * q / 100
	f := map[string]string{}
	for k, v := range map[string]int{"p1": p1, "p2": p2, "p3": p3, "q": q, "contribution1": contrib1} {
		w7I(f, k, v)
	}
	switch kp {
	case 1:
		return w7S(kp, fmt.Sprintf("A partition has probabilities %d%%, %d%% and %d%%. Within each part the event has conditional probability %d%%. What is the total probability of the event, in percent? (enter a number)", p1, p2, p3, q), strconv.Itoa(q),
			fmt.Sprintf("Every part contributes p_i*q/100, and the parts sum to 100 percent, so the total is %d percent.", q), f)
	case 2:
		return w7S(kp, fmt.Sprintf("A partition has probabilities %d%%, %d%% and %d%%. Within the first part the event has conditional probability %d%%. What is that part's contribution to the total probability, in percent? (enter a number)", p1, p2, p3, q), strconv.Itoa(contrib1),
			fmt.Sprintf("The first part contributes %d/100 * %d = %d percent.", p1, q, contrib1), f)
	default:
		return w7S(kp, fmt.Sprintf("A partition has three parts. The first is %d%% and the second is %d%%. What is the third part, in percent? (enter a number)", p1, p2), strconv.Itoa(p3),
			fmt.Sprintf("A partition sums to 100 percent, so the third part is 100-%d-%d = %d percent.", p1, p2, p3), f)
	}
}

func w7Permutations(kp int) generator.Problem {
	n := w7R(5, 12)
	f := map[string]string{}
	w7I(f, "n", n)
	switch kp {
	case 1:
		v := n * (n - 1)
		w7I(f, "r", 2)
		w7I(f, "value", v)
		return w7S(kp, fmt.Sprintf("How many ordered arrangements of 2 items can be chosen from %d distinct items? (enter a number)", n), strconv.Itoa(v),
			fmt.Sprintf("P(%d,2) = %d*%d = %d.", n, n, n-1, v), f)
	case 2:
		v := n * (n - 1) * (n - 2)
		w7I(f, "r", 3)
		w7I(f, "value", v)
		return w7S(kp, fmt.Sprintf("How many ordered arrangements of 3 items can be chosen from %d distinct items? (enter a number)", n), strconv.Itoa(v),
			fmt.Sprintf("P(%d,3) = %d*%d*%d = %d.", n, n, n-1, n-2, v), f)
	default:
		v := n * (n - 1) * (n - 2) * (n - 3)
		w7I(f, "r", 4)
		w7I(f, "value", v)
		return w7S(kp, fmt.Sprintf("How many ordered arrangements of 4 items can be chosen from %d distinct items? (enter a number)", n), strconv.Itoa(v),
			fmt.Sprintf("P(%d,4) = %d*%d*%d*%d = %d.", n, n, n-1, n-2, n-3, v), f)
	}
}

func w7Combinations(kp int) generator.Problem {
	n := w7R(6, 15)
	f := map[string]string{}
	w7I(f, "n", n)
	switch kp {
	case 1:
		v := n * (n - 1) / 2
		w7I(f, "r", 2)
		w7I(f, "value", v)
		return w7S(kp, fmt.Sprintf("How many unordered selections of 2 items can be made from %d distinct items? (enter a number)", n), strconv.Itoa(v),
			fmt.Sprintf("C(%d,2) = %d*%d/2 = %d.", n, n, n-1, v), f)
	case 2:
		v := n * (n - 1) * (n - 2) / 6
		w7I(f, "r", 3)
		w7I(f, "value", v)
		return w7S(kp, fmt.Sprintf("How many unordered selections of 3 items can be made from %d distinct items? (enter a number)", n), strconv.Itoa(v),
			fmt.Sprintf("C(%d,3) = %d*%d*%d/6 = %d.", n, n, n-1, n-2, v), f)
	default:
		v := n * (n - 1) * (n - 2) * (n - 3) / 24
		w7I(f, "r", 4)
		w7I(f, "value", v)
		return w7S(kp, fmt.Sprintf("How many unordered selections of 4 items can be made from %d distinct items? (enter a number)", n), strconv.Itoa(v),
			fmt.Sprintf("C(%d,4) = %d*%d*%d*%d/24 = %d.", n, n, n-1, n-2, n-3, v), f)
	}
}

func w7Odds(kp int) generator.Problem {
	pairs := [][2]int{{1, 3}, {2, 3}, {1, 4}, {3, 7}, {2, 8}, {4, 6}, {1, 9}, {5, 5}, {3, 2}, {7, 3}, {9, 1}, {2, 2}, {3, 3}, {1, 1}, {4, 1}, {6, 4}}
	pair := pairs[rand.Intn(len(pairs))]
	a, b := pair[0], pair[1]
	prob := a * 100 / (a + b)
	against := b * 100 / (a + b)
	f := map[string]string{}
	for k, v := range map[string]int{"favourable": a, "unfavourable": b, "probability": prob, "against": against} {
		w7I(f, k, v)
	}
	switch kp {
	case 1:
		return w7S(kp, fmt.Sprintf("The odds in favour of an event are %d to %d. What is its probability, in percent? (enter a number)", a, b), strconv.Itoa(prob),
			fmt.Sprintf("The probability is %d/(%d+%d) = %d percent.", a, a, b, prob), f)
	case 2:
		return w7S(kp, fmt.Sprintf("The odds in favour of an event are %d to %d. What is the probability against it, in percent? (enter a number)", a, b), strconv.Itoa(against),
			fmt.Sprintf("The probability against is the complement, %d percent.", against), f)
	default:
		return w7S(kp, fmt.Sprintf("An event has probability %d%% and %d favourable outcomes. How many unfavourable outcomes are there? (enter a number)", prob, a), strconv.Itoa(b),
			fmt.Sprintf("The odds are %d to b with %d/(%d+b) = %d percent, so b = %d.", a, a, a, prob, b), f)
	}
}

func w7Weibull(kp int) generator.Problem {
	k := w7R(2, 6)
	ratio := w7R(2, 8)
	value := 1
	for i := 0; i < k; i++ {
		value *= ratio
	}
	f := map[string]string{}
	for key, v := range map[string]int{"k": k, "ratio": ratio, "value": value} {
		w7I(f, key, v)
	}
	switch kp {
	case 1:
		return w7S(kp, fmt.Sprintf("In a Weibull model (t/lambda)=%d and the shape is k=%d. What is (t/lambda)^k? (enter a number)", ratio, k), strconv.Itoa(value),
			fmt.Sprintf("Raise the ratio to the shape: %d^%d = %d.", ratio, k, value), f)
	case 2:
		return w7S(kp, fmt.Sprintf("In a Weibull model (t/lambda)=%d and (t/lambda)^k=%d. What is the shape k? (enter a number)", ratio, value), strconv.Itoa(k),
			fmt.Sprintf("The ratio %d raised to k gives %d, so k = %d.", ratio, value, k), f)
	default:
		return w7S(kp, fmt.Sprintf("In a Weibull model the shape is k=%d and (t/lambda)^k=%d. What is t/lambda? (enter a number)", k, value), strconv.Itoa(ratio),
			fmt.Sprintf("Take the kth root: %d^(1/%d) = %d.", value, k, ratio), f)
	}
}

func w7Cauchy(kp int) generator.Problem {
	location := w7R(0, 20)
	scale := w7R(1, 10)
	iqr := 2 * scale
	f := map[string]string{}
	for k, v := range map[string]int{"location": location, "scale": scale, "iqr": iqr, "q1": location - scale, "q3": location + scale, "sum_quartiles": 2 * location} {
		w7I(f, k, v)
	}
	switch kp {
	case 1:
		return w7S(kp, fmt.Sprintf("A Cauchy distribution has scale %d. What is its interquartile range? (enter a number)", scale), strconv.Itoa(iqr),
			fmt.Sprintf("For a Cauchy distribution the interquartile range is 2*scale = 2*%d = %d.", scale, iqr), f)
	case 2:
		return w7S(kp, fmt.Sprintf("A Cauchy distribution has interquartile range %d. What is its scale? (enter a number)", iqr), strconv.Itoa(scale),
			fmt.Sprintf("The scale is half the interquartile range: %d/2 = %d.", iqr, scale), f)
	default:
		return w7S(kp, fmt.Sprintf("A Cauchy distribution has location %d and scale %d. What is the sum of its two quartiles? (enter a number)", location, scale), strconv.Itoa(2*location),
			fmt.Sprintf("The quartiles are %d and %d, whose sum is %d.", location-scale, location+scale, 2*location), f)
	}
}

func w7Triangular(kp int) generator.Problem {
	a := w7R(1, 10)
	b := a + w7R(2, 12)
	delta := 1
	for i := 0; i < 12; i++ {
		d := w7R(1, 12)
		if (a+b+d)%3 == 0 {
			delta = d
			break
		}
	}
	c := b + delta
	mean := (a + b + c) / 3
	f := map[string]string{}
	for k, v := range map[string]int{"a": a, "b": b, "c": c, "mean": mean, "range": c - a, "mode": c} {
		w7I(f, k, v)
	}
	switch kp {
	case 1:
		return w7S(kp, fmt.Sprintf("A triangular distribution has minimum %d, maximum %d and mode %d. What is its mean? (enter a number)", a, c, c), strconv.Itoa(mean),
			fmt.Sprintf("The mean is (a+b+c)/3 = (%d+%d+%d)/3 = %d.", a, b, c, mean), f)
	case 2:
		return w7S(kp, fmt.Sprintf("A triangular distribution has minimum %d and maximum %d. What is the width of its support? (enter a number)", a, c), strconv.Itoa(c-a),
			fmt.Sprintf("The support width is %d-%d = %d.", c, a, c-a), f)
	default:
		return w7S(kp, fmt.Sprintf("A triangular distribution has minimum %d, maximum %d and mean %d. What is its mode? (enter a number)", a, c, mean), strconv.Itoa(c),
			fmt.Sprintf("The mode of this triangular distribution is the peak at %d.", c), f)
	}
}

func w7Pareto(kp int) generator.Problem {
	alpha := w7R(2, 6)
	var xm int
	switch alpha {
	case 2:
		xm = w7R(2, 20)
	case 3:
		xm = 2 * w7R(2, 15)
	case 4:
		xm = 3 * w7R(2, 10)
	case 5:
		xm = 4 * w7R(2, 10)
	default:
		xm = 5 * w7R(2, 10)
	}
	mean := alpha * xm / (alpha - 1)
	f := map[string]string{}
	for k, v := range map[string]int{"alpha": alpha, "xm": xm, "mean": mean} {
		w7I(f, k, v)
	}
	switch kp {
	case 1:
		return w7S(kp, fmt.Sprintf("A Pareto distribution has shape alpha=%d and minimum x_m=%d. What is its mean? (enter a number)", alpha, xm), strconv.Itoa(mean),
			fmt.Sprintf("The mean is alpha*x_m/(alpha-1) = %d*%d/%d = %d.", alpha, xm, alpha-1, mean), f)
	case 2:
		return w7S(kp, fmt.Sprintf("A Pareto distribution has shape alpha=%d and mean %d. What is its minimum x_m? (enter a number)", alpha, mean), strconv.Itoa(xm),
			fmt.Sprintf("Rearrange the mean: x_m = mean*(alpha-1)/alpha = %d*%d/%d = %d.", mean, alpha-1, alpha, xm), f)
	default:
		return w7S(kp, fmt.Sprintf("A Pareto distribution has minimum x_m=%d and mean %d. What is its shape alpha? (enter a number)", xm, mean), strconv.Itoa(alpha),
			fmt.Sprintf("With mean %d and minimum %d, the ratio is %d/2, which gives alpha = %d.", mean, xm, mean/xm, alpha), f)
	}
}

func w7FisherExact(kp int) generator.Problem {
	row1 := w7R(8, 20)
	row2 := w7R(8, 20)
	total := row1 + row2
	col1 := w7R(8, total-8)
	col2 := total - col1
	maxA := minInt(row1, col1)
	minA := row1 - col2
	if minA < 0 {
		minA = 0
	}
	f := map[string]string{}
	for k, v := range map[string]int{"row1": row1, "row2": row2, "col1": col1, "col2": col2, "total": total, "max_a": maxA, "min_a": minA} {
		w7I(f, k, v)
	}
	switch kp {
	case 1:
		return w7S(kp, fmt.Sprintf("A 2x2 table has first row total %d, second row total %d, first column total %d and second column total %d. In Fisher's exact test, what is the largest possible value of the top-left cell? (enter a number)", row1, row2, col1, col2), strconv.Itoa(maxA),
			fmt.Sprintf("The top-left cell cannot exceed either its row total or its column total, so its maximum is min(%d,%d) = %d.", row1, col1, maxA), f)
	case 2:
		return w7S(kp, fmt.Sprintf("A 2x2 table has first row total %d, second row total %d, first column total %d and second column total %d. In Fisher's exact test, what is the smallest possible value of the top-left cell? (enter a number)", row1, row2, col1, col2), strconv.Itoa(minA),
			fmt.Sprintf("The top-left cell must be at least row1-col2 = %d-%d = %d.", row1, col2, minA), f)
	default:
		return w7S(kp, fmt.Sprintf("A 2x2 table has first row total %d, second row total %d, first column total %d and second column total %d. What is the grand total? (enter a number)", row1, row2, col1, col2), strconv.Itoa(total),
			fmt.Sprintf("The grand total is the sum of the row totals: %d+%d = %d, which matches the column totals.", row1, row2, total), f)
	}
}

func w7McNemar(kp int) generator.Problem {
	s := w7R(3, 7)
	t := w7R(0, 3)
	d := t * s
	if (s*s+d)%2 != 0 {
		d += s
	}
	b := (s*s + d) / 2
	c := (s*s - d) / 2
	f := map[string]string{}
	for k, v := range map[string]int{"b": b, "c": c, "discordant": b + c, "difference": b - c, "numerator": (b - c) * (b - c)} {
		w7I(f, k, v)
	}
	switch kp {
	case 1:
		return w7S(kp, fmt.Sprintf("A McNemar table has discordant counts b=%d and c=%d. What is the numerator (b-c)^2 of the test statistic? (enter a number)", b, c), strconv.Itoa((b-c)*(b-c)),
			fmt.Sprintf("The difference is %d, whose square is %d.", b-c, (b-c)*(b-c)), f)
	case 2:
		return w7S(kp, fmt.Sprintf("A McNemar table has discordant counts b=%d and c=%d. What is the number of discordant pairs? (enter a number)", b, c), strconv.Itoa(b+c),
			fmt.Sprintf("The discordant pairs are b+c = %d+%d = %d.", b, c, b+c), f)
	default:
		return w7S(kp, fmt.Sprintf("A McNemar table has discordant counts b=%d and c=%d. What is b-c? (enter a number)", b, c), strconv.Itoa(b-c),
			fmt.Sprintf("Subtract the discordant counts: %d-%d = %d.", b, c, b-c), f)
	}
}

func w7SignTest(kp int) generator.Problem {
	n := w7R(8, 15)
	vals := make([]int, n)
	k := 0
	for i := range vals {
		if rand.Intn(2) == 0 {
			vals[i] = -w7R(1, 9)
		} else {
			vals[i] = w7R(1, 9)
			k++
		}
	}
	f := map[string]string{"differences": w7Join(vals)}
	for key, v := range map[string]int{"n": n, "positive": k, "negative": n - k, "statistic": minInt(k, n-k)} {
		w7I(f, key, v)
	}
	switch kp {
	case 1:
		return w7S(kp, fmt.Sprintf("The differences in a paired sample are %s. How many are positive? (enter a number)", w7Join(vals)), strconv.Itoa(k),
			fmt.Sprintf("Count the values above zero: there are %d.", k), f)
	case 2:
		return w7S(kp, fmt.Sprintf("The differences in a paired sample are %s. Using the sign test, what is the smaller of the positive and negative counts? (enter a number)", w7Join(vals)), strconv.Itoa(minInt(k, n-k)),
			fmt.Sprintf("There are %d positive and %d negative differences, so the smaller count is %d.", k, n-k, minInt(k, n-k)), f)
	default:
		return w7S(kp, fmt.Sprintf("The differences in a paired sample are %s. How many are negative? (enter a number)", w7Join(vals)), strconv.Itoa(n-k),
			fmt.Sprintf("There are %d differences in total and %d are positive, so %d are negative.", n, k, n-k), f)
	}
}

func w7CooksDistance(kp int) generator.Problem {
	ns := []int{10, 16, 20, 25, 40, 50, 80, 100}
	n := ns[rand.Intn(len(ns))]
	threshold := 400 / n
	vals := make([]int, 8)
	for i := range vals {
		vals[i] = w7R(1, maxInt(2*threshold, 20))
	}
	count, maxIdx := 0, 0
	for i, v := range vals {
		if v > threshold {
			count++
		}
		if v > vals[maxIdx] {
			maxIdx = i
		}
	}
	f := map[string]string{"distances": w7Join(vals)}
	for key, v := range map[string]int{"n": n, "threshold_hundredths": threshold, "exceeding": count, "max_index": maxIdx + 1} {
		w7I(f, key, v)
	}
	switch kp {
	case 1:
		return w7S(kp, fmt.Sprintf("A regression has n=%d observations. Eight Cook's distances, in hundredths, are %s. Using the 4/n cutoff, how many exceed it? (enter a number)", n, w7Join(vals)), strconv.Itoa(count),
			fmt.Sprintf("The cutoff is 4/%d = %.2f, and %d of the eight distances exceed it.", n, 4.0/float64(n), count), f)
	case 2:
		return w7S(kp, fmt.Sprintf("Eight Cook's distances, in hundredths, are %s. Which observation has the largest distance? Give its position from 1. (enter a number)", w7Join(vals)), strconv.Itoa(maxIdx+1),
			fmt.Sprintf("The largest distance is %d at position %d.", vals[maxIdx], maxIdx+1), f)
	default:
		return w7S(kp, fmt.Sprintf("A regression has n=%d observations. What is the 4/n Cook's distance cutoff, in hundredths? (enter a number)", n), strconv.Itoa(threshold),
			fmt.Sprintf("The cutoff is 4/%d = %.2f, which is %d hundredths.", n, 4.0/float64(n), threshold), f)
	}
}

func w7Polynomial(kp int) generator.Problem {
	d := w7R(2, 6)
	coeffs := d + 1
	interactions := d * (d - 1) / 2
	f := map[string]string{}
	for k, v := range map[string]int{"degree": d, "coefficients": coeffs, "predictors": d, "with_interactions": d + interactions} {
		w7I(f, k, v)
	}
	switch kp {
	case 1:
		return w7S(kp, fmt.Sprintf("A polynomial regression of degree %d in one variable is fitted. How many coefficients does the model have? (enter a number)", d), strconv.Itoa(coeffs),
			fmt.Sprintf("A degree-%d polynomial has %d coefficients, from the constant term to the degree-%d term.", d, coeffs, d), f)
	case 2:
		return w7S(kp, fmt.Sprintf("A polynomial regression of degree %d in one variable is fitted. How many predictor terms carry a power of the variable? (enter a number)", d), strconv.Itoa(d),
			fmt.Sprintf("The powers run from 1 to %d, so there are %d predictor terms.", d, d), f)
	default:
		return w7S(kp, fmt.Sprintf("A polynomial regression of degree %d is fitted in %d variables and includes every pairwise interaction. How many parameters are there, counting the intercept? (enter a number)", d, d), strconv.Itoa(d+interactions+1),
			fmt.Sprintf("There are %d main effects, %d pairwise interactions and one intercept, so %d parameters.", d, interactions, d+interactions+1), f)
	}
}

func w7Ridge(kp int) generator.Problem {
	sxx := w7R(20, 60)
	lambda := w7R(5, 30)
	beta := w7R(2, 9)
	sxy := beta * (sxx + lambda)
	f := map[string]string{}
	for k, v := range map[string]int{"sxx": sxx, "lambda": lambda, "beta": beta, "sxy": sxy} {
		w7I(f, k, v)
	}
	switch kp {
	case 1:
		return w7S(kp, fmt.Sprintf("In a standardised ridge regression, Sxx=%d, lambda=%d and Sxy=%d. The ridge estimate is Sxy/(Sxx+lambda). What is it? (enter a number)", sxx, lambda, sxy), strconv.Itoa(beta),
			fmt.Sprintf("The estimate is %d/(%d+%d) = %d.", sxy, sxx, lambda, beta), f)
	case 2:
		return w7S(kp, fmt.Sprintf("In a standardised ridge regression, Sxx=%d, Sxy=%d and the ridge estimate is %d. What is lambda? (enter a number)", sxx, sxy, beta), strconv.Itoa(lambda),
			fmt.Sprintf("Rearrange Sxy/(Sxx+lambda)=beta: lambda = %d/%d-%d = %d.", sxy, beta, sxx, lambda), f)
	default:
		return w7S(kp, fmt.Sprintf("In a standardised ridge regression, lambda=%d, Sxy=%d and the ridge estimate is %d. What is Sxx? (enter a number)", lambda, sxy, beta), strconv.Itoa(sxx),
			fmt.Sprintf("Rearrange: Sxx = %d/%d-%d = %d.", sxy, beta, lambda, sxx), f)
	}
}
