package statistics

// The numeric statistics concepts: spread, quartiles, z-scores, the standard error and margin of
// error and conditional probability.
//
// Converted to parametric generators scaling with ctx.Difficulty using mathutil primitives
// and structured three-subgoal cognitive derivations in explanations.

import (
	"fmt"
	"math/rand"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/mathutil"
)

func registerNumeric(reg *generator.Registry) {
	reg.Register("stat.spread.mad", &madGen{})
	reg.Register("stat.spread.iqr", &iqrGen{})
	reg.Register("stat.quartile.value", &quartileGen{})
	reg.Register("stat.zscore.value", &zscoreGen{})
	reg.Register("stat.sampling.se", &standardErrorGen{})
	reg.Register("stat.infer.margin", &marginGen{})
	reg.Register("stat.prob.conditional", &conditionalGen{})
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// madGen asks for the mean absolute deviation of an integer data set.
type madGen struct{}

func (g *madGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	maxMu := mathutil.ScaleInt(ctx.Difficulty, 30, 150)
	mu := 20 + rand.Intn(maxMu-19)
	maxMad := mathutil.ScaleInt(ctx.Difficulty, 4, 25)
	mad := 2 + rand.Intn(maxMad-1)

	n := 4
	if ctx.Difficulty >= 0.5 {
		n = 6
	}

	var data []int
	var expSub2, expSub3 string
	if n == 4 {
		delta := rand.Intn(minInt(mad, 4))
		d1 := mad + delta
		d2 := mad - delta
		for d2 <= 0 {
			mad++
			d1 = mad + delta
			d2 = mad - delta
		}
		if mu <= d1 {
			mu = d1 + 5 + rand.Intn(10)
		}
		data = []int{mu - d1, mu - d2, mu + d2, mu + d1}
		for slices.Contains(data, mad) {
			mu += 3
			data = []int{mu - d1, mu - d2, mu + d2, mu + d1}
		}
		sumDev := 2 * (d1 + d2)
		expSub2 = fmt.Sprintf("First find the mean: \\(\\mu = (%d + %d + %d + %d)/4 = %d/4 = %d\\). Next compute each observation's absolute distance from \\(\\mu\\): \\(|%d - %d| = %d\\), \\(|%d - %d| = %d\\), \\(|%d - %d| = %d\\) and \\(|%d - %d| = %d\\).",
			data[0], data[1], data[2], data[3], 4*mu, mu,
			data[0], mu, d1, data[1], mu, d2, data[2], mu, d2, data[3], mu, d1)
		expSub3 = fmt.Sprintf("Sum the absolute deviations and divide by n = 4: \\(\\text{MAD} = (%d + %d + %d + %d)/4 = %d/4 = %d\\).",
			d1, d2, d2, d1, sumDev, mad)
	} else {
		delta := rand.Intn(minInt(mad, 3))
		d1 := mad + delta
		d2 := mad
		d3 := mad - delta
		for d3 <= 0 {
			mad++
			d1 = mad + delta
			d2 = mad
			d3 = mad - delta
		}
		if mu <= d1 {
			mu = d1 + 5 + rand.Intn(10)
		}
		data = []int{mu - d1, mu - d2, mu - d3, mu + d3, mu + d2, mu + d1}
		for slices.Contains(data, mad) {
			mu += 3
			data = []int{mu - d1, mu - d2, mu - d3, mu + d3, mu + d2, mu + d1}
		}
		sumDev := 2 * (d1 + d2 + d3)
		expSub2 = fmt.Sprintf("First find the mean: \\(\\mu = (%d + %d + %d + %d + %d + %d)/6 = %d/6 = %d\\). Next compute each observation's absolute distance from \\(\\mu\\): \\(|%d - %d| = %d\\), \\(|%d - %d| = %d\\), \\(|%d - %d| = %d\\), \\(|%d - %d| = %d\\), \\(|%d - %d| = %d\\) and \\(|%d - %d| = %d\\).",
			data[0], data[1], data[2], data[3], data[4], data[5], 6*mu, mu,
			data[0], mu, d1, data[1], mu, d2, data[2], mu, d3, data[3], mu, d3, data[4], mu, d2, data[5], mu, d1)
		expSub3 = fmt.Sprintf("Sum the absolute deviations and divide by n = 6: \\(\\text{MAD} = (%d + %d + %d + %d + %d + %d)/6 = %d/6 = %d\\).",
			d1, d2, d3, d3, d2, d1, sumDev, mad)
	}

	ans := strconv.Itoa(mad)
	q := fmt.Sprintf("Find the mean absolute deviation of the data set \\(\\{%s\\}\\).", joinInts(data))
	exp := fmt.Sprintf("Subgoal 1: State governing formula. The mean absolute deviation (MAD) measures the average distance of data values from their arithmetic mean: \\(\\text{MAD} = \\frac{1}{n} \\sum_{i=1}^n |x_i - \\mu|\\), where \\(\\mu = \\frac{1}{n} \\sum_{i=1}^n x_i\\).\n\nSubgoal 2: Substitute parameters. %s\n\nSubgoal 3: Evaluate and simplify. %s", expSub2, expSub3)

	return generator.Problem{
		Question:    q,
		Answer:      ans,
		Explanation: exp,
		Facts: map[string]string{
			"answer": ans,
		},
	}
}

// iqrGen asks for the interquartile range of a sorted sample.
type iqrGen struct{}

func (g *iqrGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	n := 8
	if ctx.Difficulty >= 0.5 {
		n = 12
	}
	m := n / 2

	maxStep := mathutil.ScaleInt(ctx.Difficulty, 3, 10)
	base := mathutil.ScaleInt(ctx.Difficulty, 4, 30)

	var vals []int
	var q1, q3, iqr int
	for {
		vals = make([]int, n)
		cur := base + rand.Intn(10)
		for i := 0; i < n; i++ {
			vals[i] = cur
			cur += 1 + rand.Intn(maxStep)
		}
		if (vals[m/2-1]+vals[m/2])%2 != 0 {
			vals[m/2]++
		}
		if (vals[m+m/2-1]+vals[m+m/2])%2 != 0 {
			vals[m+m/2]++
		}
		sort.Ints(vals)
		q1 = (vals[m/2-1] + vals[m/2]) / 2
		q3 = (vals[m+m/2-1] + vals[m+m/2]) / 2
		iqr = q3 - q1
		if iqr > 0 && !slices.Contains(vals, iqr) {
			break
		}
		base += 5
	}

	lower := vals[:m]
	upper := vals[m:]
	ans := strconv.Itoa(iqr)
	q := fmt.Sprintf("Find the interquartile range of the sorted data set \\(\\{%s\\}\\).", joinInts(vals))
	exp := fmt.Sprintf("Subgoal 1: State governing formula. The interquartile range (IQR) measures the statistical dispersion of the middle half of an ordered data set: \\(\\text{IQR} = Q_3 - Q_1\\), where \\(Q_1\\) is the lower quartile and \\(Q_3\\) is the upper quartile.\n\nSubgoal 2: Substitute parameters. The ordered data set contains %d values. Divide the data into equal halves: the lower half is \\(\\{%s\\}\\) with median \\(Q_1 = (%d + %d)/2 = %d\\) and the upper half is \\(\\{%s\\}\\) with median \\(Q_3 = (%d + %d)/2 = %d\\).\n\nSubgoal 3: Evaluate and simplify. Subtract the lower quartile from the upper quartile: \\(\\text{IQR} = %d - %d = %d\\).",
		n, joinInts(lower), vals[m/2-1], vals[m/2], q1,
		joinInts(upper), vals[m+m/2-1], vals[m+m/2], q3,
		q3, q1, iqr)

	return generator.Problem{
		Question:    q,
		Answer:      ans,
		Explanation: exp,
		Facts: map[string]string{
			"answer": ans,
		},
	}
}

// quartileGen asks for either the lower or upper quartile of a sorted sample.
type quartileGen struct{}

func (g *quartileGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	n := 8
	if ctx.Difficulty >= 0.5 {
		n = 12
	}
	m := n / 2

	step := mathutil.ScaleInt(ctx.Difficulty, 2, 6)
	base := mathutil.ScaleInt(ctx.Difficulty, 4, 30)

	var vals []int
	var q1, q3, want int
	which := "lower"
	if rand.Intn(2) == 1 {
		which = "upper"
	}

	for {
		vals = make([]int, n)
		cur := (base/2)*2 + 2
		for i := 0; i < n; i++ {
			vals[i] = cur
			cur += 2 * (1 + rand.Intn(step))
		}
		vals[m/2] = vals[m/2-1] + 2
		vals[m+m/2] = vals[m+m/2-1] + 2
		sort.Ints(vals)

		q1 = (vals[m/2-1] + vals[m/2]) / 2
		q3 = (vals[m+m/2-1] + vals[m+m/2]) / 2
		if which == "lower" {
			want = q1
		} else {
			want = q3
		}

		if !slices.Contains(vals, want) {
			break
		}
		base += 6
	}

	lower := vals[:m]
	upper := vals[m:]
	half := lower
	midA, midB := vals[m/2-1], vals[m/2]
	if which == "upper" {
		half = upper
		midA, midB = vals[m+m/2-1], vals[m+m/2]
	}

	ans := strconv.Itoa(want)
	q := fmt.Sprintf("Find the %s quartile of the sorted data set \\(\\{%s\\}\\).", which, joinInts(vals))
	exp := fmt.Sprintf("Subgoal 1: State governing formula. Quartiles partition an ordered data set into four equal parts. The lower quartile \\(Q_1\\) is the median of the lower half of the data and the upper quartile \\(Q_3\\) is the median of the upper half.\n\nSubgoal 2: Substitute parameters. The ordered data set contains %d values. Divide the data into two equal halves. The %s half comprises the %d values \\(\\{%s\\}\\).\n\nSubgoal 3: Evaluate and simplify. Locate the middle pair of the %s half: %d and %d. Compute their median: \\(Q = (%d + %d)/2 = %d/2 = %d\\).",
		n, which, m, joinInts(half),
		which, midA, midB, midA, midB, midA+midB, want)

	return generator.Problem{
		Question:    q,
		Answer:      ans,
		Explanation: exp,
		Facts: map[string]string{
			"answer": ans,
		},
	}
}

// zscoreGen asks for the z-score of an observed value.
type zscoreGen struct{}

func (g *zscoreGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	maxMu := mathutil.ScaleInt(ctx.Difficulty, 40, 150)
	mean := 30 + rand.Intn(maxMu-29)
	maxSd := mathutil.ScaleInt(ctx.Difficulty, 3, 15)
	sd := 2 + rand.Intn(maxSd-1)

	maxK := 3
	if ctx.Difficulty >= 0.5 {
		maxK = 5
	}
	k := 1 + rand.Intn(maxK)
	if rand.Intn(2) == 0 {
		k = -k
	}

	for k == sd || k == -sd || k == mean || k == -mean {
		sd++
	}

	x := mean + k*sd
	for x == k || x <= 0 {
		mean += 10
		x = mean + k*sd
	}

	ans := strconv.Itoa(k)
	q := fmt.Sprintf("A distribution has mean %d and standard deviation %d. Find the z-score of the value %d.", mean, sd, x)
	exp := fmt.Sprintf("Subgoal 1: State governing formula. The standard score or z-score transforms an observed value \\(x\\) into standard deviation units relative to the distribution mean \\(\\mu\\): \\(z = \\frac{x - \\mu}{\\sigma}\\).\n\nSubgoal 2: Substitute parameters. Given mean \\(\\mu = %d\\), standard deviation \\(\\sigma = %d\\) and value \\(x = %d\\), substitute into the formula: \\(z = \\frac{%d - %d}{%d}\\).\n\nSubgoal 3: Evaluate and simplify. Compute the deviation in the numerator \\(%d - %d = %d\\) and divide by the standard deviation: \\(z = \\frac{%d}{%d} = %d\\). A value above the mean has a positive z-score and a value below the mean has a negative z-score.",
		mean, sd, x, x, mean, sd,
		x, mean, x-mean, x-mean, sd, k)

	return generator.Problem{
		Question:    q,
		Answer:      ans,
		Explanation: exp,
		Facts: map[string]string{
			"mean":   strconv.Itoa(mean),
			"sd":     strconv.Itoa(sd),
			"x":      strconv.Itoa(x),
			"answer": ans,
		},
	}
}

// standardErrorGen asks for the standard error of the mean for a sample of square size.
type standardErrorGen struct{}

func (g *standardErrorGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	maxR := mathutil.ScaleInt(ctx.Difficulty, 4, 12)
	r := 2 + rand.Intn(maxR-1)
	n := r * r

	// The standard error is the whole answer here, so the multiplier is the
	// answer space. It used to stop at 20, which left this concept serving 19
	// distinct answers where the fixed table it replaced served 33.
	maxMult := mathutil.ScaleInt(ctx.Difficulty, 8, 40)
	mult := 2 + rand.Intn(maxMult-1)
	for mult == n || mult == r {
		mult++
	}

	s := r * mult
	ans := strconv.Itoa(mult)

	q := fmt.Sprintf("A sample of size %d has a standard deviation of %d. Find the standard error of the mean.", n, s)
	exp := fmt.Sprintf("Subgoal 1: State governing formula. The standard error of the mean (SE) measures the dispersion of sample means across repeated samples of size n: \\(\\text{SE} = \\frac{\\sigma}{\\sqrt{n}}\\), where \\(\\sigma\\) is the standard deviation and \\(n\\) is the sample size.\n\nSubgoal 2: Substitute parameters. Substitute sample size \\(n = %d\\) and sample standard deviation \\(s = %d\\) into the formula: \\(\\text{SE} = \\frac{%d}{\\sqrt{%d}}\\).\n\nSubgoal 3: Evaluate and simplify. Evaluate the square root in the denominator \\(\\sqrt{%d} = %d\\) and divide: \\(\\text{SE} = \\frac{%d}{%d} = %d\\).",
		n, s, s, n,
		n, r, s, r, mult)

	return generator.Problem{
		Question:    q,
		Answer:      ans,
		Explanation: exp,
		Facts: map[string]string{
			"size":   strconv.Itoa(n),
			"sd":     strconv.Itoa(s),
			"answer": ans,
		},
	}
}

// marginGen asks for the margin of error given a critical value and standard error.
type marginGen struct{}

func (g *marginGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	z := 2
	if rand.Intn(2) == 1 {
		z = 3
	}
	maxSe := mathutil.ScaleInt(ctx.Difficulty, 4, 30)
	se := 2 + rand.Intn(maxSe-1)
	margin := z * se
	ans := strconv.Itoa(margin)

	q := fmt.Sprintf("A confidence interval uses a critical value of %d and a standard error of %d. Find the margin of error.", z, se)
	exp := fmt.Sprintf("Subgoal 1: State governing formula. The margin of error (E) determines the half-width of a confidence interval around a point estimate: \\(E = z^* \\cdot \\text{SE}\\), where \\(z^*\\) is the critical value and \\(\\text{SE}\\) is the standard error.\n\nSubgoal 2: Substitute parameters. Substitute critical value \\(z^* = %d\\) and standard error \\(\\text{SE} = %d\\) into the product: \\(E = %d \\times %d\\).\n\nSubgoal 3: Evaluate and simplify. Multiply the factors to obtain \\(E = %d\\).",
		z, se, z, se, margin)

	return generator.Problem{
		Question:    q,
		Answer:      ans,
		Explanation: exp,
		Facts: map[string]string{
			"z":      strconv.Itoa(z),
			"se":     strconv.Itoa(se),
			"answer": ans,
		},
	}
}

// conditionalGen asks for a conditional probability given joint and marginal probabilities.
type conditionalGen struct{}

func (g *conditionalGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	maxD := mathutil.ScaleInt(ctx.Difficulty, 8, 30)
	d := 6 + rand.Intn(maxD-5)

	var both, b, num, den, gcdVal int
	for {
		b = 2 + rand.Intn(d-2)
		both = 1 + rand.Intn(b-1)
		if b*b != both*d {
			gcdVal = mathutil.GCD(both, b)
			num = both / gcdVal
			den = b / gcdVal
			break
		}
	}

	ans := fmt.Sprintf("%d/%d", num, den)
	q := fmt.Sprintf("For two events, \\(P(A \\text{ and } B) = \\frac{%d}{%d}\\) and \\(P(B) = \\frac{%d}{%d}\\). Find \\(P(A \\mid B)\\).", both, d, b, d)

	var expSub3 string
	if gcdVal > 1 {
		expSub3 = fmt.Sprintf("Cancel the common denominator %d to obtain \\(\\frac{%d}{%d}\\) and simplify by dividing both numerator and denominator by their greatest common divisor %d: \\(P(A \\mid B) = \\frac{%d}{%d}\\).",
			d, both, b, gcdVal, num, den)
	} else {
		expSub3 = fmt.Sprintf("Cancel the common denominator %d to obtain the irreducible fraction \\(P(A \\mid B) = \\frac{%d}{%d}\\).",
			d, num, den)
	}

	exp := fmt.Sprintf("Subgoal 1: State governing formula. The conditional probability of event A given event B is the ratio of their joint probability to the marginal probability of the conditioning event: \\(P(A \\mid B) = \\frac{P(A \\text{ and } B)}{P(B)}\\), where \\(P(B) > 0\\).\n\nSubgoal 2: Substitute parameters. Given joint probability \\(P(A \\text{ and } B) = \\frac{%d}{%d}\\) and condition probability \\(P(B) = \\frac{%d}{%d}\\), substitute them into the formula: \\(P(A \\mid B) = \\frac{%d/%d}{%d/%d}\\).\n\nSubgoal 3: Evaluate and simplify. %s",
		both, d, b, d, both, d, b, d, expSub3)

	return generator.Problem{
		Question:    q,
		Answer:      ans,
		Explanation: exp,
		Facts: map[string]string{
			"both":   strconv.Itoa(both),
			"b":      strconv.Itoa(b),
			"d":      strconv.Itoa(d),
			"answer": ans,
		},
	}
}

// joinInts renders a slice of integers for a question's data set.
func joinInts(v []int) string {
	parts := make([]string, len(v))
	for i, x := range v {
		parts[i] = strconv.Itoa(x)
	}
	return strings.Join(parts, ", ")
}
