package statistics

// Wave 6 statistics: multiple comparisons and false discovery rate. The
// generator uses hundredths for p-values so every threshold count is an exact
// integer, and publishes the test count and p-values for wave6_test.go.

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave6(reg *generator.Registry) {
	reg.Register("stat.infer.multiple_comparisons", &w6StatGen{kind: "multiple_comparisons"})
}

type w6StatGen struct{ kind string }

func (g *w6StatGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "multiple_comparisons":
		return w6MultipleComparisons(kp)
	}
	return generator.Problem{Question: "What is 2 + 2? (enter a number)", Answer: "4", Explanation: "2 + 2 = 4."}
}

func w6S(o int, q, a, ex string, facts map[string]string) generator.Problem {
	if facts == nil {
		facts = map[string]string{}
	}
	facts["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: facts, KPOrdinal: o}
}

func w6SI(m map[string]string, k string, v int) { m[k] = strconv.Itoa(v) }

func w6SR(lo, hi int) int {
	if hi <= lo {
		return lo
	}
	return lo + rand.Intn(hi-lo+1)
}

func w6SJoin(xs []int) string {
	parts := make([]string, len(xs))
	for i, x := range xs {
		parts[i] = strconv.Itoa(x)
	}
	return strings.Join(parts, ",")
}

func w6MultipleComparisons(kp int) generator.Problem {
	m := w6SR(4, 9)
	pv := make([]int, m)
	for i := range pv {
		pv[i] = w6SR(1, 20)
	}
	sort.Ints(pv)
	f := map[string]string{}
	w6SI(f, "m", m)
	f["pvalues"] = w6SJoin(pv)
	f["alpha_percent"] = "5"
	switch kp {
	case 1:
		count := 0
		for _, p := range pv {
			if p*m <= 5 {
				count++
			}
		}
		w6SI(f, "bonferroni_count", count)
		return w6S(kp, fmt.Sprintf("A study tests m=%d hypotheses at family-wise alpha=5%%. The p-values in hundredths, sorted, are %s. Using the Bonferroni correction, how many hypotheses are rejected? (enter a number)", m, w6SJoin(pv)), strconv.Itoa(count),
			fmt.Sprintf("The Bonferroni threshold is 5/m percent, so a p-value in hundredths is rejected when p*m <= 5. The count is %d.", count), f)
	case 2:
		alpha := w6SR(1, 10)
		w6SI(f, "alpha_percent", alpha)
		bound := alpha * m
		w6SI(f, "familywise_bound", bound)
		return w6S(kp, fmt.Sprintf("For m=%d independent tests, each at level %d%%, what is the Bonferroni family-wise error bound in percent? (enter a number)", m, alpha), strconv.Itoa(bound),
			fmt.Sprintf("Bonferroni uses %d/m percent per test. Across m=%d tests the bound is %d percent.", alpha, m, bound), f)
	default:
		count := 0
		for i, p := range pv {
			if p*m <= 5*(i+1) {
				count = i + 1
			}
		}
		w6SI(f, "bh_count", count)
		return w6S(kp, fmt.Sprintf("For m=%d tests at target false-discovery rate 5%%, the sorted p-values in hundredths are %s. Using the Benjamini-Hochberg step-up rule, how many hypotheses are rejected? (enter a number)", m, w6SJoin(pv)), strconv.Itoa(count),
			fmt.Sprintf("Compare p_i*m with 5i and take the largest i that satisfies the inequality. The rejection count is %d.", count), f)
	}
}
