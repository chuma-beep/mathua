package statistics

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

// The statistics batch used to ask recognition questions: a definition with
// three paraphrases beside it. These builders replace them with count problems
// over a pool of statements, so the answer has to be computed from the
// statements rather than matched against a remembered phrase. Each pool holds
// sixteen statements, eight true and eight false, and the answer is drawn
// uniformly from zero to eight before the statements are chosen.

type statVerdict struct {
	label string
	yes   bool
	note  string
}

const statDraw = 8

func countStatVerdicts(lead string, pool []statVerdict) generator.Problem {
	var yes, no []statVerdict
	for _, v := range pool {
		if v.yes {
			yes = append(yes, v)
		} else {
			no = append(no, v)
		}
	}

	question := ""
	chosen := make([]statVerdict, 0, statDraw)
	count := 0
	for attempt := 0; attempt < 40; attempt++ {
		count = rand.Intn(statDraw + 1)
		yesCopy := append([]statVerdict(nil), yes...)
		noCopy := append([]statVerdict(nil), no...)
		rand.Shuffle(len(yesCopy), func(i, j int) { yesCopy[i], yesCopy[j] = yesCopy[j], yesCopy[i] })
		rand.Shuffle(len(noCopy), func(i, j int) { noCopy[i], noCopy[j] = noCopy[j], noCopy[i] })
		chosen = append(chosen[:0], yesCopy[:count]...)
		chosen = append(chosen, noCopy[:statDraw-count]...)
		rand.Shuffle(len(chosen), func(i, j int) { chosen[i], chosen[j] = chosen[j], chosen[i] })

		parts := make([]string, len(chosen))
		for i, v := range chosen {
			parts[i] = fmt.Sprintf("(%c) %s", 'a'+i, v.label)
		}
		question = fmt.Sprintf("%s %s (enter a number)", lead, strings.Join(parts, " "))
		if !strings.Contains(question, fmt.Sprintf("%d", count)) {
			break
		}
	}

	notes := make([]string, len(chosen))
	for i, v := range chosen {
		verdict := "false"
		if v.yes {
			verdict = "true"
		}
		notes[i] = fmt.Sprintf("%s: %s. %s", v.label, verdict, v.note)
	}

	return generator.Problem{
		Question:    question,
		Answer:      fmt.Sprintf("%d", count),
		Explanation: fmt.Sprintf("Exactly %d of the %d statements are true. %s", count, statDraw, strings.Join(notes, " ")),
		Facts:       map[string]string{"count": fmt.Sprintf("%d", count)},
	}
}

func hypothesisProblem() generator.Problem {
	pool := []statVerdict{
		{"The null hypothesis is the default claim that the test examines.", true, "It is the baseline the data are compared against."},
		{"The alternative hypothesis is the claim accepted when the null is rejected.", true, "It states the departure the test has power to detect."},
		{"The significance level alpha caps the Type I error rate.", true, "Alpha is fixed before the data arrive."},
		{"A one-sided test concentrates its rejection region in one tail.", true, "This buys power in that direction and ignores the other."},
		{"The test statistic standardises the estimate against the null value.", true, "A z or t statistic measures the gap in standard-error units."},
		{"A p-value below alpha leads to rejecting the null.", true, "The result falls in the rejection region at that level."},
		{"The null hypothesis often states that there is no effect or no difference.", true, "This is the usual baseline in a significance test."},
		{"Hypotheses are chosen before the data are examined.", true, "Choosing after seeing the data invalidates the p-value."},
		{"The null hypothesis is the claim the researcher hopes to prove.", false, "That is usually the alternative; the null is the baseline."},
		{"The alternative hypothesis is always two-sided.", false, "It can be one-sided when only one direction matters."},
		{"Alpha is the probability that the null hypothesis is true.", false, "Alpha is an error rate chosen by the researcher, not a posterior probability."},
		{"The p-value is the probability that the null hypothesis is true.", false, "It is computed assuming the null and measures compatibility with it."},
		{"A larger alpha makes rejecting the null harder.", false, "A larger alpha widens the rejection region and makes rejection easier."},
		{"The test statistic measures the size of the sample.", false, "It measures the distance between the estimate and the null value."},
		{"A one-sided test doubles the Type I error rate.", false, "It keeps alpha fixed and moves the whole rejection region to one tail."},
		{"Hypotheses may be changed after seeing the data.", false, "That is p-hacking and destroys the meaning of the p-value."},
	}
	return countStatVerdicts("How many of these statements about hypothesis tests are true?", pool)
}

func pValueProblem() generator.Problem {
	pool := []statVerdict{
		{"A p-value is the probability under the null of data at least as extreme as observed.", true, "The calculation assumes the null hypothesis."},
		{"A small p-value is evidence against the null hypothesis.", true, "Unlikely data count against the model that made them unlikely."},
		{"A p-value below alpha leads to rejecting the null.", true, "The observed result lies in the rejection region."},
		{"A two-sided p-value counts both tails of the null distribution.", true, "Departures in either direction count as evidence."},
		{"A p-value is not the probability that the null is true.", true, "It conditions on the null rather than assigning it a probability."},
		{"A result with p equal to 0.20 at alpha 0.05 leads to failing to reject the null.", true, "The p-value lies above the threshold."},
		{"The p-value is computed from the test statistic and its null distribution.", true, "It is a tail probability of that distribution."},
		{"A smaller p-value is stronger evidence against the null.", true, "More extreme data are less compatible with the null."},
		{"A p-value is the probability that the null hypothesis is true.", false, "That is a common misinterpretation; the p-value conditions on the null."},
		{"A p-value is the probability that the alternative hypothesis is true.", false, "The p-value says nothing directly about the alternative."},
		{"A p-value above alpha proves that the null is true.", false, "Failing to reject is not the same as proving the null."},
		{"A two-sided p-value counts only the upper tail.", false, "Both tails contribute to a two-sided p-value."},
		{"A p-value measures the size of the effect.", false, "A tiny effect in a large sample can still produce a small p-value."},
		{"A result with p equal to 0.03 at alpha 0.05 is not statistically significant.", false, "It is below alpha and therefore significant at that level."},
		{"The p-value is the same quantity as alpha.", false, "Alpha is chosen in advance; the p-value is computed from the data."},
		{"A large p-value proves that the alternative is false.", false, "It only means the data are compatible with the null."},
	}
	return countStatVerdicts("How many of these statements about p-values are true?", pool)
}

func typeErrorProblem() generator.Problem {
	pool := []statVerdict{
		{"A Type I error rejects a true null hypothesis.", true, "This is the false alarm, occurring at rate alpha."},
		{"A Type II error fails to reject a false null hypothesis.", true, "This is the missed detection."},
		{"Power is one minus the Type II error rate.", true, "Power measures the chance of detecting a real effect."},
		{"Lowering alpha while holding the sample size fixed raises the Type II error rate.", true, "A stricter threshold makes missed detections more likely."},
		{"Power increases when the sample size increases.", true, "More data shrink the standard error."},
		{"Power increases when the effect size increases.", true, "A larger signal is easier to distinguish from noise."},
		{"The Type I error rate is alpha.", true, "Alpha is defined as the maximum false alarm rate."},
		{"Power is the probability of rejecting a false null hypothesis.", true, "This is the complement of the Type II error rate."},
		{"A Type I error fails to reject a false null hypothesis.", false, "That is a Type II error."},
		{"A Type II error rejects a true null hypothesis.", false, "That is a Type I error."},
		{"Power is equal to alpha.", false, "Power is one minus beta, not the significance level."},
		{"Lowering alpha lowers both kinds of error.", false, "It trades fewer Type I errors for more Type II errors."},
		{"Power falls when the sample size increases.", false, "More data increase power."},
		{"Power falls when the effect size increases.", false, "A larger effect is easier to detect."},
		{"The Type I error rate is beta.", false, "Beta is the Type II error rate."},
		{"Power is the probability of accepting a true null hypothesis.", false, "Power concerns rejecting a false null."},
	}
	return countStatVerdicts("How many of these statements about errors and power are true?", pool)
}

func nonparametricProblem() generator.Problem {
	pool := []statVerdict{
		{"The sign test uses the binomial distribution with probability one half under the null.", true, "Each pair is equally likely to favour either side."},
		{"The sign test uses only the signs of the differences.", true, "Discarding magnitudes makes it robust to outliers."},
		{"Rank-based tests are robust to outliers.", true, "A few extreme values cannot move ranks far."},
		{"The Wilcoxon signed-rank test uses the ranks of the absolute differences.", true, "Ranks retain some magnitude information."},
		{"Nonparametric tests make fewer distributional assumptions.", true, "They do not require a normal population."},
		{"The Mann-Whitney test compares two independent samples by rank.", true, "It is the rank analogue of the two-sample t test."},
		{"A nonparametric test usually has less power when its parametric assumptions hold.", true, "Fewer assumptions cost efficiency when the model is correct."},
		{"Permutation tests rely on exchangeability under the null.", true, "Labels are shuffled to build the null distribution."},
		{"The sign test uses the normal distribution under the null.", false, "It uses the binomial distribution with probability one half."},
		{"The sign test uses the magnitudes of the differences.", false, "It uses only their signs."},
		{"Rank-based tests require a normal population.", false, "They are designed to avoid that requirement."},
		{"The Wilcoxon signed-rank test uses only the signs.", false, "It ranks the absolute differences."},
		{"Nonparametric tests make no assumptions at all.", false, "They still assume independent observations and a meaningful null."},
		{"A nonparametric test is always more powerful than a parametric test.", false, "It usually loses power when the parametric model holds."},
		{"The Mann-Whitney test is a paired test.", false, "It compares two independent samples."},
		{"Permutation tests require normally distributed data.", false, "They require exchangeability, not normality."},
	}
	return countStatVerdicts("How many of these statements about nonparametric tests are true?", pool)
}
