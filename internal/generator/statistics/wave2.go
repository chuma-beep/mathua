package statistics

// The 2026-11-09 Wave 2 statistics batch: fourteen concepts that extend the
// inference and distribution spine of the domain.
//
// Every concept has three question families, one per knowledge point. The
// first family asks for a direct computation, the second applies the same
// relation in the reverse direction or to a second quantity, and the third
// interprets or combines what the first two produced. Numeric families publish
// the operands in Problem.Facts so a test can recompute the answer, and the two
// classification families publish the probabilities they compare.

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave2(reg *generator.Registry) {
	reg.Register("stat.data.weighted_mean", &w2WeightedMeanGen{})
	reg.Register("stat.adv.coefficient_variation", &w2CoefficientVariationGen{})
	reg.Register("stat.prob.independence", &w2IndependenceGen{})
	reg.Register("stat.prob.mutually_exclusive", &w2MutuallyExclusiveGen{})
	reg.Register("stat.dist.negative_binomial", &w2NegativeBinomialGen{})
	reg.Register("stat.dist.multinomial", &w2MultinomialGen{})
	reg.Register("stat.dist.lognormal", &w2LognormalGen{})
	reg.Register("stat.infer.ci_mean_t", &w2CIMeanTGen{})
	reg.Register("stat.infer.ci_proportion", &w2CIProportionGen{})
	reg.Register("stat.infer.power", &w2PowerGen{})
	reg.Register("stat.infer.effect_size", &w2EffectSizeGen{})
	reg.Register("stat.reg.multiple", &w2MultipleRegressionGen{})
	reg.Register("stat.reg.prediction_interval", &w2PredictionIntervalGen{})
	reg.Register("stat.test.kruskal_wallis", &w2KruskalWallisGen{})
}

// ---------------------------------------------------------------------------
// Shared helpers
// ---------------------------------------------------------------------------

func w2Ordinal(ctx generator.GeneratorContext) int {
	if ctx.KPOrdinal >= 1 && ctx.KPOrdinal <= 3 {
		return ctx.KPOrdinal
	}
	return 1 + rand.Intn(3)
}

func w2Done(p generator.Problem, ctx generator.GeneratorContext) generator.Problem {
	p.KPOrdinal = ctx.KPOrdinal
	return p
}

func w2GCD(a, b int) int {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func w2Fraction(num, den int) string {
	if den == 0 {
		return "0"
	}
	if den < 0 {
		num, den = -num, -den
	}
	g := w2GCD(num, den)
	if g == 0 {
		return "0"
	}
	num /= g
	den /= g
	if den == 1 {
		return strconv.Itoa(num)
	}
	return fmt.Sprintf("%d/%d", num, den)
}

func w2Choose(n, k int) int {
	if k < 0 || k > n {
		return 0
	}
	if k > n-k {
		k = n - k
	}
	out := 1
	for i := 1; i <= k; i++ {
		out = out * (n - k + i) / i
	}
	return out
}

// w2Decimal writes a rounded value with trailing zeros removed, so 12.5000
// becomes 12.5 and 12.0000 becomes 12. The question states the rounding rule.
func w2Decimal(x float64, places int) string {
	s := strconv.FormatFloat(x, 'f', places, 64)
	if strings.Contains(s, ".") {
		s = strings.TrimRight(s, "0")
		s = strings.TrimRight(s, ".")
	}
	if s == "-0" || s == "" {
		return "0"
	}
	return s
}

func w2IntStrings(vals []int) []string {
	out := make([]string, len(vals))
	for i, v := range vals {
		out[i] = strconv.Itoa(v)
	}
	return out
}

// w2Noun writes a count with the singular or plural noun that agrees with it.
func w2Noun(n int, one, many string) string {
	if n == 1 {
		return "1 " + one
	}
	return strconv.Itoa(n) + " " + many
}

// w2JoinAnd joins a prose list without an Oxford comma: "a, b and c".
func w2JoinAnd(parts []string) string {
	switch len(parts) {
	case 0:
		return ""
	case 1:
		return parts[0]
	case 2:
		return parts[0] + " and " + parts[1]
	default:
		return strings.Join(parts[:len(parts)-1], ", ") + " and " + parts[len(parts)-1]
	}
}

// w2FactInts stores a numeric list in a fact without confusing the separator
// with the comma a decimal uses.
func w2FactInts(vals []int) string {
	return strings.Join(w2IntStrings(vals), ";")
}

func w2FactGroups(groups [][]int) string {
	parts := make([]string, len(groups))
	for i, group := range groups {
		parts[i] = strings.Join(w2IntStrings(group), ",")
	}
	return strings.Join(parts, "|")
}

func w2GroupsText(groups [][]int) string {
	parts := make([]string, len(groups))
	for i, group := range groups {
		parts[i] = fmt.Sprintf("group %d: %s", i+1, w2JoinAnd(w2IntStrings(group)))
	}
	return strings.Join(parts, "; ")
}

func w2LatexFrac(num, den int) string {
	if den == 1 {
		return strconv.Itoa(num)
	}
	return fmt.Sprintf("\\frac{%d}{%d}", num, den)
}

// ---------------------------------------------------------------------------
// 1. Weighted mean
// ---------------------------------------------------------------------------

type w2WeightedMeanGen struct{}

func w2WeightedCase() ([]int, []int, int) {
	for attempt := 0; attempt < generator.MaxDrawAttempts; attempt++ {
		n := 3 + rand.Intn(2)
		weights := make([]int, n)
		sumWeights := 0
		for i := range weights {
			weights[i] = 1 + rand.Intn(5)
			sumWeights += weights[i]
		}
		mean := 60 + rand.Intn(31)
		values := make([]int, n)
		remaining := mean * sumWeights
		for i := 0; i < n-1; i++ {
			values[i] = 40 + rand.Intn(61)
			remaining -= values[i] * weights[i]
		}
		lastWeight := weights[n-1]
		if remaining%lastWeight != 0 {
			continue
		}
		last := remaining / lastWeight
		if last < 0 || last > 100 {
			continue
		}
		values[n-1] = last
		return values, weights, mean
	}
	return []int{60, 100, 80}, []int{1, 2, 1}, 85
}

func (g *w2WeightedMeanGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	ordinal := w2Ordinal(ctx)
	switch ordinal {
	case 1:
		values, weights, mean := w2WeightedCase()
		terms := make([]string, len(values))
		for i := range values {
			terms[i] = fmt.Sprintf("%d \\times %d", values[i], weights[i])
		}
		return w2Done(generator.Problem{
			Question: fmt.Sprintf("A student scores %s on %d tests, with weights %s. What is the weighted mean?",
				w2JoinAnd(w2IntStrings(values)), len(values), w2JoinAnd(w2IntStrings(weights))),
			Answer: strconv.Itoa(mean),
			Explanation: fmt.Sprintf("The weighted mean divides the weighted total by the total weight: \\(\\frac{%s}{%s} = \\frac{%d}{%d} = %d\\).",
				strings.Join(terms, " + "), w2SumStrings(weights), w2WeightedTotal(values, weights), w2SumInts(weights), mean),
			Facts: map[string]string{
				"values":  w2FactInts(values),
				"weights": w2FactInts(weights),
				"n":       strconv.Itoa(len(values)),
			},
		}, ctx)

	case 2:
		values, weights, mean := w2WeightedCase()
		missing := rand.Intn(len(values))
		knownValues := make([]int, 0, len(values)-1)
		knownWeights := make([]int, 0, len(weights)-1)
		for i := range values {
			if i == missing {
				continue
			}
			knownValues = append(knownValues, values[i])
			knownWeights = append(knownWeights, weights[i])
		}
		knownParts := make([]string, len(knownValues))
		for i := range knownValues {
			knownParts[i] = fmt.Sprintf("%d with weight %d", knownValues[i], knownWeights[i])
		}
		return w2Done(generator.Problem{
			Question: fmt.Sprintf("A weighted mean of %d is required across %d components. The known scores are %s. The remaining component has weight %d. What score is needed on it? (enter a number)",
				mean, len(values), w2JoinAnd(knownParts), weights[missing]),
			Answer: strconv.Itoa(values[missing]),
			Explanation: fmt.Sprintf("Let \\(x\\) be the missing score. The weighted mean equation is \\(\\frac{%s + %s}{%d} = %d\\). The known weighted total is %d, so \\(%s = %d \\times %d - %d = %d\\) and \\(x = %d\\).",
				w2KnownWeightedTerms(knownValues, knownWeights), w2CoeffVar(weights[missing], "x"), w2SumInts(weights), mean,
				w2KnownWeightedTotal(knownValues, knownWeights),
				w2CoeffVar(weights[missing], "x"), mean, w2SumInts(weights), w2KnownWeightedTotal(knownValues, knownWeights),
				mean*w2SumInts(weights)-w2KnownWeightedTotal(knownValues, knownWeights),
				values[missing]),
			Facts: map[string]string{
				"known_values":   w2FactInts(knownValues),
				"known_weights":  w2FactInts(knownWeights),
				"missing_weight": strconv.Itoa(weights[missing]),
				"target":         strconv.Itoa(mean),
			},
		}, ctx)

	default:
		values, weights, mean := w2WeightedCase()
		unweighted := float64(w2SumInts(values)) / float64(len(values))
		answer := w2Decimal(float64(mean)-unweighted, 2)
		return w2Done(generator.Problem{
			Question: fmt.Sprintf("The scores %s carry weights %s. What is the weighted mean minus the unweighted mean? (enter a decimal; a negative value means the weighted mean is smaller)",
				w2JoinAnd(w2IntStrings(values)), w2JoinAnd(w2IntStrings(weights))),
			Answer: answer,
			Explanation: fmt.Sprintf("The weighted mean is \\(\\frac{%s}{%d} = %d\\). The unweighted mean is \\(\\frac{%s}{%d} = %s\\). The difference is \\(%d - %s = %s\\).",
				w2WeightedTermText(values, weights), w2SumInts(weights), mean,
				strings.Join(w2IntStrings(values), " + "), len(values), w2Decimal(unweighted, 2),
				mean, w2Decimal(unweighted, 2), answer),
			Facts: map[string]string{
				"values":  w2FactInts(values),
				"weights": w2FactInts(weights),
				"n":       strconv.Itoa(len(values)),
			},
		}, ctx)
	}
}

func w2SumInts(vals []int) int {
	total := 0
	for _, v := range vals {
		total += v
	}
	return total
}

func w2SumStrings(vals []int) string {
	return strings.Join(w2IntStrings(vals), " + ")
}

func w2WeightedTotal(values, weights []int) int {
	total := 0
	for i := range values {
		total += values[i] * weights[i]
	}
	return total
}

func w2KnownWeightedTotal(values, weights []int) int {
	return w2WeightedTotal(values, weights)
}

func w2KnownWeightedTerms(values, weights []int) string {
	parts := make([]string, len(values))
	for i := range values {
		parts[i] = fmt.Sprintf("%d \\times %d", values[i], weights[i])
	}
	return strings.Join(parts, " + ")
}

func w2WeightedTermText(values, weights []int) string {
	parts := make([]string, len(values))
	for i := range values {
		parts[i] = fmt.Sprintf("%d \\times %d", values[i], weights[i])
	}
	return strings.Join(parts, " + ")
}

// ---------------------------------------------------------------------------
// 2. Coefficient of variation
// ---------------------------------------------------------------------------

type w2CoefficientVariationGen struct{}

func w2CVCase() (int, int, int) {
	for attempt := 0; attempt < generator.MaxDrawAttempts; attempt++ {
		mean := 10 * (5 + rand.Intn(16))
		cv := 5 * (1 + rand.Intn(12))
		if mean*cv%100 != 0 {
			continue
		}
		sd := mean * cv / 100
		if sd <= 0 {
			continue
		}
		return mean, sd, cv
	}
	return 100, 15, 15
}

func (g *w2CoefficientVariationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	ordinal := w2Ordinal(ctx)
	switch ordinal {
	case 1:
		mean, sd, cv := w2CVCase()
		return w2Done(generator.Problem{
			Question: fmt.Sprintf("A sample has mean %d and standard deviation %d. What is the coefficient of variation as a percentage? (enter a number)", mean, sd),
			Answer:   strconv.Itoa(cv),
			Explanation: fmt.Sprintf("The coefficient of variation is the standard deviation divided by the mean, written as a percentage: \\(CV = \\frac{%d}{%d} \\times 100\\%% = %d\\%%\\).",
				sd, mean, cv),
			Facts: map[string]string{
				"mean": strconv.Itoa(mean),
				"sd":   strconv.Itoa(sd),
			},
		}, ctx)

	case 2:
		mean, sd, cv := w2CVCase()
		return w2Done(generator.Problem{
			Question: fmt.Sprintf("A sample has mean %d and coefficient of variation %d%%. What is its standard deviation? (enter a number)", mean, cv),
			Answer:   strconv.Itoa(sd),
			Explanation: fmt.Sprintf("For the coefficient of variation, rearrange \\(CV = \\frac{s}{\\bar{x}} \\times 100\\) to \\(s = \\frac{CV \\cdot \\bar{x}}{100} = \\frac{%d \\times %d}{100} = %d\\).",
				cv, mean, sd),
			Facts: map[string]string{
				"mean":       strconv.Itoa(mean),
				"cv_percent": strconv.Itoa(cv),
			},
		}, ctx)

	default:
		meanA, sdA, cvA := w2CVCase()
		meanB, sdB, cvB := w2CVCase()
		for cvA == cvB {
			meanB, sdB, cvB = w2CVCase()
		}
		ratio := float64(cvA) / float64(cvB)
		answer := w2Decimal(ratio, 2)
		return w2Done(generator.Problem{
			Question: fmt.Sprintf("Sample A has mean %d and standard deviation %d. Sample B has mean %d and standard deviation %d. What is the ratio of the coefficient of variation of A to that of B? (round to two decimal places)",
				meanA, sdA, meanB, sdB),
			Answer: answer,
			Explanation: fmt.Sprintf("The coefficient of variation of A is \\(\\frac{%d}{%d} \\times 100\\%% = %d\\%%\\) and that of B is \\(\\frac{%d}{%d} \\times 100\\%% = %d\\%%\\). Their ratio is \\(\\frac{%d}{%d} = %s\\).",
				sdA, meanA, cvA, sdB, meanB, cvB, cvA, cvB, answer),
			Facts: map[string]string{
				"mean_a": strconv.Itoa(meanA),
				"sd_a":   strconv.Itoa(sdA),
				"mean_b": strconv.Itoa(meanB),
				"sd_b":   strconv.Itoa(sdB),
			},
		}, ctx)
	}
}

// ---------------------------------------------------------------------------
// 3. Independent events (classification)
// ---------------------------------------------------------------------------

type w2IndependenceGen struct{}

const (
	w2Independent          = "Independent"
	w2NotIndependent       = "Not independent"
	w2MutuallyExclusive    = "Mutually exclusive"
	w2NotMutuallyExclusive = "Not mutually exclusive"
)

func w2Decision(question, correct, other, explanation string, facts map[string]string) generator.Problem {
	p := generator.MultipleChoice(question, correct, other)
	p.Explanation = explanation
	p.Facts = facts
	return p
}

func (g *w2IndependenceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	ordinal := w2Ordinal(ctx)
	switch ordinal {
	case 1:
		a := 1 + rand.Intn(9)
		b := 1 + rand.Intn(9)
		pabNum := a * b
		correct := w2Independent
		if rand.Intn(2) == 0 {
			maxNum := minInt(a, b) * 10
			pabNum = 1 + rand.Intn(maxNum)
			for pabNum == a*b {
				pabNum = 1 + rand.Intn(maxNum)
			}
			correct = w2NotIndependent
		}
		product := w2Fraction(a*b, 100)
		intersection := w2Fraction(pabNum, 100)
		explanation := fmt.Sprintf("For independence, \\(P(A \\text{ and } B)\\) must equal \\(P(A)P(B) = %s \\times %s = %s\\). The stated intersection is %s, so the events are %s.",
			w2LatexFrac(a, 10), w2LatexFrac(b, 10), product, intersection, strings.ToLower(correct))
		return w2Done(w2Decision(
			fmt.Sprintf("Let \\(P(A) = %s\\), \\(P(B) = %s\\) and \\(P(A \\text{ and } B) = %s\\). Are A and B independent?",
				w2LatexFrac(a, 10), w2LatexFrac(b, 10), w2LatexFrac(pabNum, 100)),
			correct, w2NotIndependentIf(correct), explanation,
			map[string]string{
				"p_a_num":  strconv.Itoa(a),
				"p_a_den":  "10",
				"p_b_num":  strconv.Itoa(b),
				"p_b_den":  "10",
				"p_ab_num": strconv.Itoa(pabNum),
				"p_ab_den": "100",
			}), ctx)

	case 2:
		a := 1 + rand.Intn(9)
		conditional := a
		correct := w2Independent
		if rand.Intn(2) == 0 {
			conditional = 1 + rand.Intn(9)
			for conditional == a {
				conditional = 1 + rand.Intn(9)
			}
			correct = w2NotIndependent
		}
		explanation := fmt.Sprintf("Independence means the conditioning does not change the probability: \\(P(A \\mid B) = P(A)\\). Here \\(P(A) = %s\\) and \\(P(A \\mid B) = %s\\), so the events are %s.",
			w2LatexFrac(a, 10), w2LatexFrac(conditional, 10), strings.ToLower(correct))
		return w2Done(w2Decision(
			fmt.Sprintf("For events A and B, \\(P(A) = %s\\) and \\(P(A \\mid B) = %s\\). Are A and B independent?",
				w2LatexFrac(a, 10), w2LatexFrac(conditional, 10)),
			correct, w2NotIndependentIf(correct), explanation,
			map[string]string{
				"p_a_num":         strconv.Itoa(a),
				"p_a_den":         "10",
				"p_a_given_b_num": strconv.Itoa(conditional),
				"p_a_given_b_den": "10",
			}), ctx)

	default:
		var a, b, union int
		var independent bool
		for attempt := 0; attempt < generator.MaxDrawAttempts; attempt++ {
			a = 1 + rand.Intn(9)
			b = 1 + rand.Intn(9)
			independent = rand.Intn(2) == 0
			if independent {
				if (a*b)%10 != 0 {
					continue
				}
				union = a + b - a*b/10
			} else {
				low := maxInt(a, b)
				high := minInt(10, a+b)
				if high < low {
					continue
				}
				union = low + rand.Intn(high-low+1)
			}
			intersection := a + b - union
			if intersection < 0 || intersection > minInt(a, b) {
				continue
			}
			isIndependent := 10*intersection == a*b
			if isIndependent == independent {
				break
			}
		}
		correct := w2NotIndependent
		if 10*(a+b-union) == a*b {
			correct = w2Independent
		}
		intersection := a + b - union
		explanation := fmt.Sprintf("First recover the intersection: \\(P(A \\text{ and } B) = P(A) + P(B) - P(A \\text{ or } B) = %s + %s - %s = %s\\). Compare it with \\(P(A)P(B) = %s\\). They %s, so the events are %s.",
			w2LatexFrac(a, 10), w2LatexFrac(b, 10), w2LatexFrac(union, 10),
			w2LatexFrac(intersection, 10), w2LatexFrac(a*b, 100),
			w2MatchVerb(10*intersection == a*b), strings.ToLower(correct))
		return w2Done(w2Decision(
			fmt.Sprintf("For events A and B, \\(P(A) = %s\\), \\(P(B) = %s\\) and \\(P(A \\text{ or } B) = %s\\). Are A and B independent?",
				w2LatexFrac(a, 10), w2LatexFrac(b, 10), w2LatexFrac(union, 10)),
			correct, w2NotIndependentIf(correct), explanation,
			map[string]string{
				"p_a_num":      strconv.Itoa(a),
				"p_a_den":      "10",
				"p_b_num":      strconv.Itoa(b),
				"p_b_den":      "10",
				"p_a_or_b_num": strconv.Itoa(union),
				"p_a_or_b_den": "10",
			}), ctx)
	}
}

func w2NotIndependentIf(correct string) string {
	if correct == w2Independent {
		return w2NotIndependent
	}
	return w2Independent
}

func w2MatchVerb(equal bool) string {
	if equal {
		return "agree"
	}
	return "differ"
}

// ---------------------------------------------------------------------------
// 4. Mutually exclusive events (classification)
// ---------------------------------------------------------------------------

type w2MutuallyExclusiveGen struct{}

func (g *w2MutuallyExclusiveGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	ordinal := w2Ordinal(ctx)
	switch ordinal {
	case 1:
		a := 1 + rand.Intn(9)
		b := 1 + rand.Intn(9)
		pab := 0
		correct := w2MutuallyExclusive
		if rand.Intn(2) == 0 {
			pab = 1 + rand.Intn(minInt(a, b))
			correct = w2NotMutuallyExclusive
		}
		explanation := fmt.Sprintf("Events are mutually exclusive exactly when \\(P(A \\text{ and } B) = 0\\). Here that probability is \\(%s\\), so the events are %s.",
			w2LatexFrac(pab, 10), strings.ToLower(correct))
		return w2Done(w2Decision(
			fmt.Sprintf("Events A and B have \\(P(A) = %s\\), \\(P(B) = %s\\) and \\(P(A \\text{ and } B) = %s\\). Are A and B mutually exclusive?",
				w2LatexFrac(a, 10), w2LatexFrac(b, 10), w2LatexFrac(pab, 10)),
			correct, w2OtherExclusive(correct), explanation,
			map[string]string{
				"p_a_num":  strconv.Itoa(a),
				"p_a_den":  "10",
				"p_b_num":  strconv.Itoa(b),
				"p_b_den":  "10",
				"p_ab_num": strconv.Itoa(pab),
				"p_ab_den": "10",
			}), ctx)

	case 2:
		var a, b, union int
		var exclusive bool
		for attempt := 0; attempt < generator.MaxDrawAttempts; attempt++ {
			a = 1 + rand.Intn(9)
			b = 1 + rand.Intn(9)
			exclusive = rand.Intn(2) == 0
			if exclusive {
				if a+b > 10 {
					continue
				}
				union = a + b
			} else {
				low := maxInt(a, b)
				high := minInt(10, a+b-1)
				if high < low {
					continue
				}
				union = low + rand.Intn(high-low+1)
			}
			if (union == a+b) == exclusive {
				break
			}
		}
		correct := w2NotMutuallyExclusive
		if union == a+b {
			correct = w2MutuallyExclusive
		}
		explanation := fmt.Sprintf("If A and B were mutually exclusive then \\(P(A \\text{ or } B)\\) would equal \\(P(A) + P(B) = %s + %s = %s\\). The stated union is \\(%s\\), so the events are %s.",
			w2LatexFrac(a, 10), w2LatexFrac(b, 10), w2LatexFrac(a+b, 10),
			w2LatexFrac(union, 10), strings.ToLower(correct))
		return w2Done(w2Decision(
			fmt.Sprintf("Events A and B have \\(P(A) = %s\\), \\(P(B) = %s\\) and \\(P(A \\text{ or } B) = %s\\). Are A and B mutually exclusive?",
				w2LatexFrac(a, 10), w2LatexFrac(b, 10), w2LatexFrac(union, 10)),
			correct, w2OtherExclusive(correct), explanation,
			map[string]string{
				"p_a_num":      strconv.Itoa(a),
				"p_a_den":      "10",
				"p_b_num":      strconv.Itoa(b),
				"p_b_den":      "10",
				"p_a_or_b_num": strconv.Itoa(union),
				"p_a_or_b_den": "10",
			}), ctx)

	default:
		a := 1 + rand.Intn(9)
		b := 1 + rand.Intn(9)
		possible := a+b <= 10
		correct := "Not possible"
		other := "Possible"
		if possible {
			correct = "Possible"
			other = "Not possible"
		}
		explanation := fmt.Sprintf("Mutually exclusive events satisfy \\(P(A) + P(B) \\le 1\\). Here \\(%s + %s = %s\\), so such events are %s.",
			w2LatexFrac(a, 10), w2LatexFrac(b, 10), w2LatexFrac(a+b, 10), strings.ToLower(correct))
		return w2Done(w2Decision(
			fmt.Sprintf("Can events A and B with \\(P(A) = %s\\) and \\(P(B) = %s\\) be mutually exclusive?",
				w2LatexFrac(a, 10), w2LatexFrac(b, 10)),
			correct, other, explanation,
			map[string]string{
				"p_a_num": strconv.Itoa(a),
				"p_a_den": "10",
				"p_b_num": strconv.Itoa(b),
				"p_b_den": "10",
			}), ctx)
	}
}

func w2OtherExclusive(correct string) string {
	if correct == w2MutuallyExclusive {
		return w2NotMutuallyExclusive
	}
	return w2MutuallyExclusive
}

// ---------------------------------------------------------------------------
// 5. Negative binomial distribution
// ---------------------------------------------------------------------------

type w2NegativeBinomialGen struct{}

func (g *w2NegativeBinomialGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	ordinal := w2Ordinal(ctx)
	switch ordinal {
	case 1:
		pDen := 2 + rand.Intn(9)
		r := 1 + rand.Intn(6)
		mean := r * pDen
		return w2Done(generator.Problem{
			Question: fmt.Sprintf("Each trial succeeds with probability \\(%s\\). What is the expected number of trials needed to obtain %d successes? (enter a number)",
				w2LatexFrac(1, pDen), r),
			Answer: strconv.Itoa(mean),
			Explanation: fmt.Sprintf("For the negative binomial count of trials, \\(E[X] = \\frac{r}{p} = \\frac{%d}{%s} = %d \\times %d = %d\\).",
				r, w2LatexFrac(1, pDen), r, pDen, mean),
			Facts: map[string]string{
				"r":     strconv.Itoa(r),
				"p_num": "1",
				"p_den": strconv.Itoa(pDen),
			},
		}, ctx)

	case 2:
		pDen := 2 + rand.Intn(5)
		r := 1 + rand.Intn(4)
		k := r + rand.Intn(5)
		num := w2Choose(k-1, r-1) * w2PowInt(pDen-1, k-r)
		den := w2PowInt(pDen, k)
		answer := w2Fraction(num, den)
		return w2Done(generator.Problem{
			Question: fmt.Sprintf("Each trial succeeds with probability \\(%s\\). What is the probability that exactly %d trials are needed to obtain %d successes? (enter a fraction)",
				w2LatexFrac(1, pDen), k, r),
			Answer: answer,
			Explanation: fmt.Sprintf("The negative binomial probability is \\(P(X = %d) = \\binom{%d}{%d} p^{%d}(1-p)^{%d} = %d \\times \\left(%s\\right)^{%d}\\left(%s\\right)^{%d} = %s\\).",
				k, k-1, r-1, r, k-r, w2Choose(k-1, r-1),
				w2LatexFrac(1, pDen), r, w2LatexFrac(pDen-1, pDen), k-r, answer),
			Facts: map[string]string{
				"r":     strconv.Itoa(r),
				"k":     strconv.Itoa(k),
				"p_num": "1",
				"p_den": strconv.Itoa(pDen),
			},
		}, ctx)

	default:
		pDen := 2 + rand.Intn(9)
		r := 1 + rand.Intn(6)
		failures := r * (pDen - 1)
		return w2Done(generator.Problem{
			Question: fmt.Sprintf("Each trial succeeds with probability \\(%s\\). What is the expected number of failures before success number %d? (enter a number)",
				w2LatexFrac(1, pDen), r),
			Answer: strconv.Itoa(failures),
			Explanation: fmt.Sprintf("For the negative binomial distribution, the expected number of failures is \\(\\frac{r(1-p)}{p} = \\frac{%d \\times %s}{%s} = %d \\times (%d - 1) = %d\\).",
				r, w2LatexFrac(pDen-1, pDen), w2LatexFrac(1, pDen), r, pDen, failures),
			Facts: map[string]string{
				"r":     strconv.Itoa(r),
				"p_num": "1",
				"p_den": strconv.Itoa(pDen),
			},
		}, ctx)
	}
}

func w2PowInt(base, exp int) int {
	out := 1
	for i := 0; i < exp; i++ {
		out *= base
	}
	return out
}

// ---------------------------------------------------------------------------
// 6. Multinomial distribution
// ---------------------------------------------------------------------------

type w2MultinomialGen struct{}

func w2MultinomialCounts(n int) []int {
	for attempt := 0; attempt < generator.MaxDrawAttempts; attempt++ {
		c1 := rand.Intn(n + 1)
		c2 := rand.Intn(n - c1 + 1)
		c3 := n - c1 - c2
		if c1+c2+c3 == n {
			return []int{c1, c2, c3}
		}
	}
	return []int{n, 0, 0}
}

func (g *w2MultinomialGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	ordinal := w2Ordinal(ctx)
	switch ordinal {
	case 1:
		n := 4 + rand.Intn(5)
		counts := w2MultinomialCounts(n)
		for counts[0] == 0 || counts[1] == 0 || counts[2] == 0 {
			counts = w2MultinomialCounts(n)
		}
		coeff := w2Choose(n, counts[0]) * w2Choose(n-counts[0], counts[1])
		return w2Done(generator.Problem{
			Question: fmt.Sprintf("A multinomial experiment has %d trials and three categories. How many distinct sequences produce counts %s? (enter a number)",
				n, w2JoinAnd(w2IntStrings(counts))),
			Answer: strconv.Itoa(coeff),
			Explanation: fmt.Sprintf("The multinomial coefficient counts the sequences: \\(\\frac{%d!}{%d!\\,%d!\\,%d!} = %d\\).",
				n, counts[0], counts[1], counts[2], coeff),
			Facts: map[string]string{
				"n":      strconv.Itoa(n),
				"counts": w2FactInts(counts),
			},
		}, ctx)

	case 2:
		n := 50 * (1 + rand.Intn(8))
		pNum := 1 + rand.Intn(9)
		expected := n * pNum / 10
		return w2Done(generator.Problem{
			Question: fmt.Sprintf("A multinomial experiment has %d trials. One category has probability \\(%s\\). What is the expected count in that category? (enter a number)",
				n, w2LatexFrac(pNum, 10)),
			Answer: strconv.Itoa(expected),
			Explanation: fmt.Sprintf("For a multinomial experiment, the expected count is \\(n p = %d \\times %s = %d\\).",
				n, w2LatexFrac(pNum, 10), expected),
			Facts: map[string]string{
				"n":     strconv.Itoa(n),
				"p_num": strconv.Itoa(pNum),
				"p_den": "10",
			},
		}, ctx)

	default:
		n := 3 + rand.Intn(4)
		counts := w2MultinomialCounts(n)
		coeff := w2Choose(n, counts[0]) * w2Choose(n-counts[0], counts[1])
		den := w2PowInt(3, n)
		answer := w2Fraction(coeff, den)
		return w2Done(generator.Problem{
			Question: fmt.Sprintf("Each of three categories is equally likely. What is the probability of counts %s in %d trials? (enter a fraction)",
				w2JoinAnd(w2IntStrings(counts)), n),
			Answer: answer,
			Explanation: fmt.Sprintf("With equal probabilities \\(p = \\frac{1}{3}\\), the multinomial probability is \\(\\frac{%d!}{%d!\\,%d!\\,%d!}\\left(\\frac{1}{3}\\right)^{%d} = \\frac{%d}{%d} = %s\\).",
				n, counts[0], counts[1], counts[2], n, coeff, den, answer),
			Facts: map[string]string{
				"n":      strconv.Itoa(n),
				"counts": w2FactInts(counts),
				"p_num":  "1",
				"p_den":  "3",
			},
		}, ctx)
	}
}

// ---------------------------------------------------------------------------
// 7. Lognormal distribution
// ---------------------------------------------------------------------------

type w2LognormalGen struct{}

func (g *w2LognormalGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	ordinal := w2Ordinal(ctx)
	mu := math.Round((0.1+rand.Float64()*1.9)*100) / 100
	sigma := math.Round((0.1+rand.Float64()*0.8)*100) / 100
	muText := w2Decimal(mu, 2)
	sigmaText := w2Decimal(sigma, 2)
	switch ordinal {
	case 1:
		answer := w2Decimal(math.Exp(mu), 4)
		return w2Done(generator.Problem{
			Question: fmt.Sprintf("The natural logarithm of \\(X\\) is normally distributed with mean %s and standard deviation %s. What is the median of \\(X\\)? (round to four decimal places)",
				muText, sigmaText),
			Answer: answer,
			Explanation: fmt.Sprintf("For a lognormal variable the median is \\(e^{\\mu} = e^{%s} = %s\\). The median is not affected by \\(\\sigma\\) because \\(P(X \\le e^{\\mu}) = 0.5\\).",
				muText, answer),
			Facts: map[string]string{"mu": muText, "sigma": sigmaText},
		}, ctx)

	case 2:
		answer := w2Decimal(math.Exp(mu+sigma*sigma/2), 4)
		return w2Done(generator.Problem{
			Question: fmt.Sprintf("The natural logarithm of \\(X\\) is normally distributed with mean %s and standard deviation %s. What is the mean of \\(X\\)? (round to four decimal places)",
				muText, sigmaText),
			Answer: answer,
			Explanation: fmt.Sprintf("The lognormal mean is \\(e^{\\mu + \\sigma^2/2} = e^{%s + %s/2} = %s\\).",
				muText, w2Decimal(sigma*sigma, 4), answer),
			Facts: map[string]string{"mu": muText, "sigma": sigmaText},
		}, ctx)

	default:
		answer := w2Decimal(math.Exp(sigma*sigma/2), 4)
		return w2Done(generator.Problem{
			Question: fmt.Sprintf("The natural logarithm of \\(X\\) is normally distributed with mean %s and standard deviation %s. What is the ratio of the mean of \\(X\\) to its median? (round to four decimal places)",
				muText, sigmaText),
			Answer: answer,
			Explanation: fmt.Sprintf("For a lognormal variable, the mean is \\(e^{\\mu + \\sigma^2/2}\\) and the median is \\(e^{\\mu}\\), so their ratio is \\(e^{\\sigma^2/2} = e^{%s/2} = %s\\).",
				w2Decimal(sigma*sigma, 4), answer),
			Facts: map[string]string{"mu": muText, "sigma": sigmaText},
		}, ctx)
	}
}

// ---------------------------------------------------------------------------
// 8. Confidence interval for a mean with t
// ---------------------------------------------------------------------------

type w2TEntry struct {
	df int
	t  float64
}

// w2TTable is the two-sided 95% critical value of Student's t. The generator
// chooses only the listed degrees of freedom, so no interpolation rule is
// needed: a question always uses the printed value for its own df.
var w2TTable = []w2TEntry{
	{1, 12.706}, {2, 4.303}, {3, 3.182}, {4, 2.776}, {5, 2.571},
	{6, 2.447}, {7, 2.365}, {8, 2.306}, {9, 2.262}, {10, 2.228},
	{12, 2.179}, {15, 2.131}, {20, 2.086}, {25, 2.060}, {30, 2.042},
	{40, 2.021}, {60, 2.000}, {120, 1.980},
}

type w2CIMeanTGen struct{}

func (g *w2CIMeanTGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	ordinal := w2Ordinal(ctx)
	entry := w2TTable[rand.Intn(len(w2TTable))]
	n := entry.df + 1
	s := 2 + rand.Intn(19)
	xbar := 20 + rand.Intn(81)
	tText := w2Decimal(entry.t, 3)
	se := float64(s) / math.Sqrt(float64(n))
	margin := entry.t * se
	switch ordinal {
	case 1:
		return w2Done(generator.Problem{
			Question: fmt.Sprintf("A sample of %d observations has standard deviation %d. What is the standard error \\(s/\\sqrt{n}\\) used in the t interval? (round to three decimal places)",
				n, s),
			Answer: w2Decimal(se, 3),
			Explanation: fmt.Sprintf("The standard error is \\(\\frac{s}{\\sqrt{n}} = \\frac{%d}{\\sqrt{%d}} = %s\\). The t critical value for %d degrees of freedom is then applied to this quantity.",
				s, n, w2Decimal(se, 3), entry.df),
			Facts: map[string]string{
				"xbar":   strconv.Itoa(xbar),
				"s":      strconv.Itoa(s),
				"n":      strconv.Itoa(n),
				"df":     strconv.Itoa(entry.df),
				"t_crit": tText,
			},
		}, ctx)

	case 2:
		return w2Done(generator.Problem{
			Question: fmt.Sprintf("A sample of %d observations has mean %d and standard deviation %d. For 95%% confidence the t critical value with %d degrees of freedom is %s. What is the margin of error? (round to three decimal places)",
				n, xbar, s, entry.df, tText),
			Answer: w2Decimal(margin, 3),
			Explanation: fmt.Sprintf("The standard error is the sample standard deviation over the square root of n, and the margin is \\(t^* \\frac{s}{\\sqrt{n}} = %s \\times \\frac{%d}{\\sqrt{%d}} = %s \\times %s = %s\\).",
				tText, s, n, tText, w2Decimal(se, 4), w2Decimal(margin, 3)),
			Facts: map[string]string{
				"xbar":   strconv.Itoa(xbar),
				"s":      strconv.Itoa(s),
				"n":      strconv.Itoa(n),
				"df":     strconv.Itoa(entry.df),
				"t_crit": tText,
			},
		}, ctx)

	default:
		lower := float64(xbar) - margin
		return w2Done(generator.Problem{
			Question: fmt.Sprintf("A sample of %d observations has mean %d and standard deviation %d. For 95%% confidence the t critical value with %d degrees of freedom is %s. What is the lower endpoint of the confidence interval? (round to three decimal places)",
				n, xbar, s, entry.df, tText),
			Answer: w2Decimal(lower, 3),
			Explanation: fmt.Sprintf("The standard error is the sample standard deviation over the square root of n, and the margin is \\(%s \\times \\frac{%d}{\\sqrt{%d}} = %s\\), so the lower endpoint is \\(%d - %s = %s\\).",
				tText, s, n, w2Decimal(margin, 3), xbar, w2Decimal(margin, 3), w2Decimal(lower, 3)),
			Facts: map[string]string{
				"xbar":   strconv.Itoa(xbar),
				"s":      strconv.Itoa(s),
				"n":      strconv.Itoa(n),
				"df":     strconv.Itoa(entry.df),
				"t_crit": tText,
			},
		}, ctx)
	}
}

// ---------------------------------------------------------------------------
// 9. Confidence interval for a proportion
// ---------------------------------------------------------------------------

type w2CIProportionGen struct{}

func (g *w2CIProportionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	ordinal := w2Ordinal(ctx)
	phatNum := 2 + rand.Intn(7)
	phatDen := 10
	phat := float64(phatNum) / float64(phatDen)
	sizes := []int{50, 64, 96, 100, 150, 200, 300, 400, 500}
	n := sizes[rand.Intn(len(sizes))]
	type zChoice struct {
		z    float64
		text string
	}
	choices := []zChoice{{1.645, "1.645"}, {1.96, "1.96"}, {2.576, "2.576"}}
	zc := choices[rand.Intn(len(choices))]
	se := math.Sqrt(phat * (1 - phat) / float64(n))
	margin := zc.z * se
	switch ordinal {
	case 1:
		return w2Done(generator.Problem{
			Question: fmt.Sprintf("A poll of %d respondents finds a sample proportion of \\(%s\\). What is the standard error of the sample proportion? (round to four decimal places)",
				n, w2LatexFrac(phatNum, phatDen)),
			Answer: w2Decimal(se, 4),
			Explanation: fmt.Sprintf("The standard error is \\(\\sqrt{\\frac{\\hat p(1-\\hat p)}{n}} = \\sqrt{\\frac{%s \\times %s}{%d}} = %s\\).",
				w2LatexFrac(phatNum, phatDen), w2LatexFrac(phatDen-phatNum, phatDen), n,
				w2Decimal(se, 4)),
			Facts: map[string]string{
				"phat_num": strconv.Itoa(phatNum),
				"phat_den": strconv.Itoa(phatDen),
				"n":        strconv.Itoa(n),
				"z_crit":   zc.text,
			},
		}, ctx)

	case 2:
		return w2Done(generator.Problem{
			Question: fmt.Sprintf("A poll of %d respondents finds a sample proportion of \\(%s\\). Using the critical value \\(z^* = %s\\), what is the margin of error? (round to four decimal places)",
				n, w2LatexFrac(phatNum, phatDen), zc.text),
			Answer: w2Decimal(margin, 4),
			Explanation: fmt.Sprintf("The standard error is the square root of the Bernoulli variance over n, and the margin is \\(z^* \\sqrt{\\frac{\\hat p(1-\\hat p)}{n}} = %s \\times \\sqrt{\\frac{%s}{%d}} = %s \\times %s = %s\\).",
				zc.text, w2LatexFrac(phatNum*(phatDen-phatNum), phatDen*phatDen), n,
				zc.text, w2Decimal(se, 4), w2Decimal(margin, 4)),
			Facts: map[string]string{
				"phat_num": strconv.Itoa(phatNum),
				"phat_den": strconv.Itoa(phatDen),
				"n":        strconv.Itoa(n),
				"z_crit":   zc.text,
			},
		}, ctx)

	default:
		lower := phat - margin
		return w2Done(generator.Problem{
			Question: fmt.Sprintf("A poll of %d respondents finds a sample proportion of \\(%s\\). Using the critical value \\(z^* = %s\\), what is the lower endpoint of the confidence interval? (round to four decimal places)",
				n, w2LatexFrac(phatNum, phatDen), zc.text),
			Answer: w2Decimal(lower, 4),
			Explanation: fmt.Sprintf("The standard error is the square root of the Bernoulli variance over n, and the margin is \\(%s \\times \\sqrt{\\frac{%s}{%d}} = %s\\), so the lower endpoint is \\(%s - %s = %s\\).",
				zc.text, w2LatexFrac(phatNum*(phatDen-phatNum), phatDen*phatDen), n,
				w2Decimal(margin, 4), w2Decimal(phat, 4), w2Decimal(margin, 4), w2Decimal(lower, 4)),
			Facts: map[string]string{
				"phat_num": strconv.Itoa(phatNum),
				"phat_den": strconv.Itoa(phatDen),
				"n":        strconv.Itoa(n),
				"z_crit":   zc.text,
			},
		}, ctx)
	}
}

// ---------------------------------------------------------------------------
// 10. Power of a hypothesis test
// ---------------------------------------------------------------------------

type w2PowerGen struct{}

// w2TailTable is the upper-tail probability P(Z > z) for the z values this
// generator can produce. A question states the critical value, so the learner
// either reads the same small table or uses one of the exact half-plane cases.
var w2TailTable = map[int]string{
	-2: "0.9772",
	-1: "0.8413",
	0:  "0.5",
}

func (g *w2PowerGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	ordinal := w2Ordinal(ctx)
	switch ordinal {
	case 1:
		beta := 0.05 + 0.05*float64(rand.Intn(8))
		answer := w2Decimal(1-beta, 2)
		return w2Done(generator.Problem{
			Question: fmt.Sprintf("A hypothesis test has probability \\(\\beta = %s\\) of a Type II error. What is the power of the test? (enter a decimal)",
				w2Decimal(beta, 2)),
			Answer: answer,
			Explanation: fmt.Sprintf("Power is the probability of rejecting a false null hypothesis: \\(1 - \\beta = 1 - %s = %s\\).",
				w2Decimal(beta, 2), answer),
			Facts: map[string]string{"beta": w2Decimal(beta, 2)},
		}, ctx)

	case 2:
		power := 0.50 + 0.05*float64(rand.Intn(10))
		answer := w2Decimal(1-power, 2)
		return w2Done(generator.Problem{
			Question: fmt.Sprintf("A hypothesis test has power %s. What is the probability of a Type II error? (enter a decimal)",
				w2Decimal(power, 2)),
			Answer: answer,
			Explanation: fmt.Sprintf("Because power \\(= 1 - \\beta\\), the Type II error probability is \\(\\beta = 1 - %s = %s\\).",
				w2Decimal(power, 2), answer),
			Facts: map[string]string{"power": w2Decimal(power, 2)},
		}, ctx)

	default:
		mu0 := float64(40 + 10*rand.Intn(6))
		se := float64([]int{1, 2, 4}[rand.Intn(3)])
		zAlpha := float64([]int{1, 2}[rand.Intn(2)])
		zAlt := []int{-2, -1, 0}[rand.Intn(3)]
		delta := se * (zAlpha - float64(zAlt))
		mu1 := mu0 + delta
		critical := mu0 + zAlpha*se
		answer := w2TailTable[zAlt]
		return w2Done(generator.Problem{
			Question: fmt.Sprintf("A test of \\(H_0: \\mu = %s\\) rejects when the sample mean exceeds %s. The standard error is %s and the true mean is %s. What is the power of the test? (round to four decimal places)",
				w2Decimal(mu0, 0), w2Decimal(critical, 0), w2Decimal(se, 0), w2Decimal(mu1, 0)),
			Answer: answer,
			Explanation: fmt.Sprintf("Under the true mean the standardized critical value is \\(z = \\frac{%s - %s}{%s} = %d\\). The power is \\(P(Z > %d)\\), which the standard normal table gives as %s.",
				w2Decimal(critical, 0), w2Decimal(mu1, 0), w2Decimal(se, 0), zAlt, zAlt, answer),
			Facts: map[string]string{
				"null_mean": w2Decimal(mu0, 0),
				"alt_mean":  w2Decimal(mu1, 0),
				"se":        w2Decimal(se, 0),
				"critical":  w2Decimal(critical, 0),
				"z_alt":     strconv.Itoa(zAlt),
			},
		}, ctx)
	}
}

// ---------------------------------------------------------------------------
// 11. Effect size
// ---------------------------------------------------------------------------

type w2EffectSizeGen struct{}

func (g *w2EffectSizeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	ordinal := w2Ordinal(ctx)
	switch ordinal {
	case 1:
		mu := 40 + rand.Intn(41)
		s := 2 + rand.Intn(9)
		diff := 1 + rand.Intn(10)
		if rand.Intn(2) == 0 {
			diff = -diff
		}
		xbar := mu + diff
		answer := w2Decimal(float64(diff)/float64(s), 3)
		return w2Done(generator.Problem{
			Question: fmt.Sprintf("A sample has mean %d, a hypothesized mean of %d and standard deviation %d. What is Cohen's d? (round to three decimal places)",
				xbar, mu, s),
			Answer: answer,
			Explanation: fmt.Sprintf("Cohen's d is the mean difference in standard deviation units: \\(d = \\frac{\\bar{x} - \\mu}{s} = \\frac{%d - %d}{%d} = %s\\).",
				xbar, mu, s, answer),
			Facts: map[string]string{
				"xbar": strconv.Itoa(xbar),
				"mu":   strconv.Itoa(mu),
				"s":    strconv.Itoa(s),
			},
		}, ctx)

	case 2:
		dChoices := []float64{0.2, 0.4, 0.5, 0.6, 0.8, 1.0, 1.2, 1.5, 2.0}
		d := dChoices[rand.Intn(len(dChoices))]
		s := 2 + rand.Intn(9)
		answer := w2Decimal(d*float64(s), 3)
		return w2Done(generator.Problem{
			Question: fmt.Sprintf("A study reports Cohen's d = %s with standard deviation %d. By how many units do the two means differ? (round to three decimal places)",
				w2Decimal(d, 2), s),
			Answer: answer,
			Explanation: fmt.Sprintf("For Cohen's d, rearrange \\(d = \\frac{\\bar{x}_1 - \\bar{x}_2}{s}\\) to get the mean difference \\(d s = %s \\times %d = %s\\).",
				w2Decimal(d, 2), s, answer),
			Facts: map[string]string{
				"d": w2Decimal(d, 2),
				"s": strconv.Itoa(s),
			},
		}, ctx)

	default:
		n := 10 * (1 + rand.Intn(4))
		s1 := 3 + rand.Intn(8)
		s2 := 3 + rand.Intn(8)
		mean1 := 40 + rand.Intn(41)
		mean2 := 40 + rand.Intn(41)
		for mean1 == mean2 {
			mean2 = 40 + rand.Intn(41)
		}
		pooled := math.Sqrt((float64(n-1)*float64(s1*s1) + float64(n-1)*float64(s2*s2)) / float64(2*n-2))
		answer := w2Decimal((float64(mean1)-float64(mean2))/pooled, 3)
		return w2Done(generator.Problem{
			Question: fmt.Sprintf("Two groups each have %d observations. Group 1 has mean %d and standard deviation %d; group 2 has mean %d and standard deviation %d. What is Cohen's d using the pooled standard deviation? (round to three decimal places)",
				n, mean1, s1, mean2, s2),
			Answer: answer,
			Explanation: fmt.Sprintf("For Cohen's d, the pooled variance is \\(\\frac{(%d-1)%d^2 + (%d-1)%d^2}{%d} = %s\\), so the pooled standard deviation is %s. Then \\(d = \\frac{%d - %d}{%s} = %s\\).",
				n, s1, n, s2, 2*n-2, w2Decimal((float64(n-1)*float64(s1*s1)+float64(n-1)*float64(s2*s2))/float64(2*n-2), 2),
				w2Decimal(pooled, 3), mean1, mean2, w2Decimal(pooled, 3), answer),
			Facts: map[string]string{
				"mean1": strconv.Itoa(mean1),
				"mean2": strconv.Itoa(mean2),
				"s1":    strconv.Itoa(s1),
				"s2":    strconv.Itoa(s2),
				"n1":    strconv.Itoa(n),
				"n2":    strconv.Itoa(n),
			},
		}, ctx)
	}
}

// ---------------------------------------------------------------------------
// 12. Multiple regression
// ---------------------------------------------------------------------------

type w2MultipleRegressionGen struct{}

func w2SignTerm(coef int, variable string) string {
	if coef == 1 {
		return "+ " + variable
	}
	if coef == -1 {
		return "- " + variable
	}
	if coef >= 0 {
		return fmt.Sprintf("+ %d%s", coef, variable)
	}
	return fmt.Sprintf("- %d%s", -coef, variable)
}

// w2CoeffVar writes a coefficient and its variable, omitting a coefficient of
// one so a term never reads "1x".
func w2CoeffVar(coef int, variable string) string {
	if coef == 1 {
		return variable
	}
	if coef == -1 {
		return "-" + variable
	}
	return fmt.Sprintf("%d%s", coef, variable)
}

func w2RegressionExpression(intercept, b1, b2 int) string {
	return fmt.Sprintf("%d %s %s", intercept, w2SignTerm(b1, "x_1"), w2SignTerm(b2, "x_2"))
}

func (g *w2MultipleRegressionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	ordinal := w2Ordinal(ctx)
	b0 := -10 + rand.Intn(21)
	b1 := -5 + rand.Intn(11)
	b2 := -5 + rand.Intn(11)
	for b1 == 0 {
		b1 = -5 + rand.Intn(11)
	}
	for b2 == 0 {
		b2 = -5 + rand.Intn(11)
	}
	x1 := 1 + rand.Intn(10)
	x2 := 1 + rand.Intn(10)
	predicted := b0 + b1*x1 + b2*x2
	switch ordinal {
	case 1:
		return w2Done(generator.Problem{
			Question: fmt.Sprintf("A multiple regression model is \\(\\hat{y} = %s\\). What is the predicted value when \\(x_1 = %d\\) and \\(x_2 = %d\\)? (enter a number)",
				w2RegressionExpression(b0, b1, b2), x1, x2),
			Answer: strconv.Itoa(predicted),
			Explanation: fmt.Sprintf("In the multiple regression model, substitute both predictors: \\(\\hat{y} = %d %s %s = %d\\).",
				b0, w2SignTerm(b1, fmt.Sprintf("(%d)", x1)), w2SignTerm(b2, fmt.Sprintf("(%d)", x2)), predicted),
			Facts: map[string]string{
				"intercept": strconv.Itoa(b0),
				"b1":        strconv.Itoa(b1),
				"b2":        strconv.Itoa(b2),
				"x1":        strconv.Itoa(x1),
				"x2":        strconv.Itoa(x2),
			},
		}, ctx)

	case 2:
		observed := predicted + (-6 + rand.Intn(13))
		residual := observed - predicted
		return w2Done(generator.Problem{
			Question: fmt.Sprintf("A multiple regression model is \\(\\hat{y} = %s\\). For \\(x_1 = %d\\) and \\(x_2 = %d\\) the observed response is %d. What is the residual? (enter a number)",
				w2RegressionExpression(b0, b1, b2), x1, x2, observed),
			Answer: strconv.Itoa(residual),
			Explanation: fmt.Sprintf("In the multiple regression model, the fitted value is \\(\\hat{y} = %d\\). The residual is \\(y - \\hat{y} = %d - %d = %d\\).",
				predicted, observed, predicted, residual),
			Facts: map[string]string{
				"intercept": strconv.Itoa(b0),
				"b1":        strconv.Itoa(b1),
				"b2":        strconv.Itoa(b2),
				"x1":        strconv.Itoa(x1),
				"x2":        strconv.Itoa(x2),
				"observed":  strconv.Itoa(observed),
			},
		}, ctx)

	default:
		r2 := (40 + 5*rand.Intn(11))
		n := []int{20, 25, 30, 40, 50, 60}[rand.Intn(6)]
		predictors := 2 + rand.Intn(4)
		adjusted := 1 - (1-float64(r2)/100)*float64(n-1)/float64(n-predictors-1)
		answer := w2Decimal(adjusted, 4)
		return w2Done(generator.Problem{
			Question: fmt.Sprintf("A multiple regression with %d predictors and %d observations reports \\(R^2 = %s\\). What is the adjusted \\(R^2\\)? (round to four decimal places)",
				predictors, n, w2Decimal(float64(r2)/100, 2)),
			Answer: answer,
			Explanation: fmt.Sprintf("For the multiple regression, the adjusted value is \\(1 - (1 - R^2)\\frac{n-1}{n-k-1} = 1 - %s \\times \\frac{%d}{%d} = %s\\).",
				w2Decimal(1-float64(r2)/100, 2), n-1, n-predictors-1, answer),
			Facts: map[string]string{
				"r2":         w2Decimal(float64(r2)/100, 2),
				"n":          strconv.Itoa(n),
				"predictors": strconv.Itoa(predictors),
			},
		}, ctx)
	}
}

// ---------------------------------------------------------------------------
// 13. Prediction intervals
// ---------------------------------------------------------------------------

type w2PredictionIntervalGen struct{}

func (g *w2PredictionIntervalGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	ordinal := w2Ordinal(ctx)
	switch ordinal {
	case 1:
		tCrit := []float64{1.96, 2.0, 2.228, 2.306, 2.776}[rand.Intn(5)]
		sPred := math.Round((0.5+rand.Float64()*9.5)*10) / 10
		half := tCrit * sPred
		return w2Done(generator.Problem{
			Question: fmt.Sprintf("A prediction interval for a new response uses critical value %s and standard error of prediction %s. What is the half-width of the interval? (round to two decimal places)",
				w2Decimal(tCrit, 3), w2Decimal(sPred, 1)),
			Answer: w2Decimal(half, 2),
			Explanation: fmt.Sprintf("The prediction interval half-width is \\(t^* s_{pred} = %s \\times %s = %s\\).",
				w2Decimal(tCrit, 3), w2Decimal(sPred, 1), w2Decimal(half, 2)),
			Facts: map[string]string{
				"t_crit": w2Decimal(tCrit, 3),
				"s_pred": w2Decimal(sPred, 1),
			},
		}, ctx)

	case 2:
		tCrit := []float64{1.96, 2.0, 2.228, 2.306, 2.776}[rand.Intn(5)]
		sPred := math.Round((0.5+rand.Float64()*9.5)*10) / 10
		yhat := 20 + rand.Intn(81)
		half := tCrit * sPred
		upper := float64(yhat) + half
		return w2Done(generator.Problem{
			Question: fmt.Sprintf("A regression predicts %d for a new observation. The prediction interval uses critical value %s and standard error of prediction %s. What is its upper endpoint? (round to two decimal places)",
				yhat, w2Decimal(tCrit, 3), w2Decimal(sPred, 1)),
			Answer: w2Decimal(upper, 2),
			Explanation: fmt.Sprintf("The prediction interval half-width is \\(%s \\times %s = %s\\), so the upper endpoint is \\(%d + %s = %s\\).",
				w2Decimal(tCrit, 3), w2Decimal(sPred, 1), w2Decimal(half, 2),
				yhat, w2Decimal(half, 2), w2Decimal(upper, 2)),
			Facts: map[string]string{
				"yhat":   strconv.Itoa(yhat),
				"t_crit": w2Decimal(tCrit, 3),
				"s_pred": w2Decimal(sPred, 1),
			},
		}, ctx)

	default:
		pairs := [][2]int{{3, 4}, {5, 12}, {8, 15}, {7, 24}, {9, 40}, {20, 21}, {6, 8}}
		pair := pairs[rand.Intn(len(pairs))]
		answer := w2Decimal(math.Sqrt(float64(pair[0]*pair[0]+pair[1]*pair[1])), 2)
		return w2Done(generator.Problem{
			Question: fmt.Sprintf("The standard error of prediction for a new response combines the residual standard error %d and the standard error of the fitted mean %d. What is \\(s_{pred} = \\sqrt{s_e^2 + s_{mean}^2}\\)? (round to two decimal places)",
				pair[0], pair[1]),
			Answer: answer,
			Explanation: fmt.Sprintf("For a prediction interval, \\(s_{pred} = \\sqrt{%d^2 + %d^2} = \\sqrt{%d + %d} = \\sqrt{%d} = %s\\).",
				pair[0], pair[1], pair[0]*pair[0], pair[1]*pair[1], pair[0]*pair[0]+pair[1]*pair[1], answer),
			Facts: map[string]string{
				"s_e":    strconv.Itoa(pair[0]),
				"s_mean": strconv.Itoa(pair[1]),
			},
		}, ctx)
	}
}

// ---------------------------------------------------------------------------
// 14. Kruskal-Wallis test
// ---------------------------------------------------------------------------

type w2KruskalWallisGen struct{}

func w2KruskalGroups(k int) [][]int {
	total := 0
	sizes := make([]int, k)
	for i := range sizes {
		sizes[i] = 2 + rand.Intn(3)
		total += sizes[i]
	}
	pool := rand.Perm(30)
	groups := make([][]int, k)
	next := 0
	for i := 0; i < k; i++ {
		groups[i] = make([]int, sizes[i])
		for j := range groups[i] {
			groups[i][j] = pool[next] + 1
			next++
		}
	}
	return groups
}

func w2RankSums(groups [][]int) []int {
	type item struct {
		value int
		group int
	}
	var items []item
	for gi, group := range groups {
		for _, v := range group {
			items = append(items, item{v, gi})
		}
	}
	for i := 1; i < len(items); i++ {
		for j := i; j > 0 && items[j].value < items[j-1].value; j-- {
			items[j], items[j-1] = items[j-1], items[j]
		}
	}
	sums := make([]int, len(groups))
	for rank, it := range items {
		sums[it.group] += rank + 1
	}
	return sums
}

func w2KruskalH(groups [][]int) float64 {
	sums := w2RankSums(groups)
	n := 0
	for _, group := range groups {
		n += len(group)
	}
	term := 0.0
	for i, group := range groups {
		term += float64(sums[i]*sums[i]) / float64(len(group))
	}
	return 12.0/float64(n*(n+1))*term - 3.0*float64(n+1)
}

func (g *w2KruskalWallisGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	ordinal := w2Ordinal(ctx)
	switch ordinal {
	case 1:
		k := 3 + rand.Intn(2)
		groups := w2KruskalGroups(k)
		groupIndex := 1 + rand.Intn(k)
		sums := w2RankSums(groups)
		answer := sums[groupIndex-1]
		return w2Done(generator.Problem{
			Question: fmt.Sprintf("A Kruskal-Wallis test ranks all observations from smallest to largest. The samples are %s. What is the sum of the ranks of group %d? (enter a number)",
				w2GroupsText(groups), groupIndex),
			Answer: strconv.Itoa(answer),
			Explanation: fmt.Sprintf("For the Kruskal-Wallis test, ordering all %d observations gives ranks 1 through %d. The values in group %d occupy ranks whose sum is %d.",
				w2GroupTotal(groups), w2GroupTotal(groups), groupIndex, answer),
			Facts: map[string]string{
				"groups":      w2FactGroups(groups),
				"k":           strconv.Itoa(k),
				"group_index": strconv.Itoa(groupIndex),
			},
		}, ctx)

	case 2:
		k := 3 + rand.Intn(2)
		groups := w2KruskalGroups(k)
		answer := w2Decimal(w2KruskalH(groups), 2)
		return w2Done(generator.Problem{
			Question: fmt.Sprintf("A Kruskal-Wallis test compares %d independent samples. All observations are distinct. The samples are %s. What is the test statistic \\(H\\)? (round to two decimal places)",
				k, w2GroupsText(groups)),
			Answer: answer,
			Explanation: fmt.Sprintf("For the Kruskal-Wallis test, rank all observations from 1 to %d and sum the ranks in each group. The statistic is \\(H = \\frac{12}{N(N+1)}\\sum \\frac{R_i^2}{n_i} - 3(N+1)\\), which evaluates to %s.",
				w2GroupTotal(groups), answer),
			Facts: map[string]string{
				"groups": w2FactGroups(groups),
				"k":      strconv.Itoa(k),
			},
		}, ctx)

	default:
		k := 3 + rand.Intn(6)
		n := 10 + rand.Intn(41)
		return w2Done(generator.Problem{
			Question: fmt.Sprintf("A Kruskal-Wallis test compares %d independent samples with a total of %d observations. How many degrees of freedom does the chi-square approximation use? (enter a number)",
				k, n),
			Answer: strconv.Itoa(k - 1),
			Explanation: fmt.Sprintf("For the Kruskal-Wallis test, the chi-square approximation has \\(k - 1\\) degrees of freedom, where \\(k\\) is the number of groups: \\(%d - 1 = %d\\).",
				k, k-1),
			Facts: map[string]string{
				"k": strconv.Itoa(k),
				"n": strconv.Itoa(n),
			},
		}, ctx)
	}
}

func w2GroupTotal(groups [][]int) int {
	total := 0
	for _, group := range groups {
		total += len(group)
	}
	return total
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
