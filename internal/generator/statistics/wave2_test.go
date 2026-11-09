package statistics

// Independent fact-based checks for the Wave 2 statistics batch.
//
// Every concept is drawn at all three knowledge-point ordinals over 300 seeds.
// The test parses Problem.Facts, recomputes the answer from the printed
// operands, and compares. It never reads a fact that already holds the answer,
// so a generator that copies a value into both Answer and Facts fails here.

import (
	"math"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

var wave2ConceptIDs = []string{
	"stat.data.weighted_mean",
	"stat.adv.coefficient_variation",
	"stat.prob.independence",
	"stat.prob.mutually_exclusive",
	"stat.dist.negative_binomial",
	"stat.dist.multinomial",
	"stat.dist.lognormal",
	"stat.infer.ci_mean_t",
	"stat.infer.ci_proportion",
	"stat.infer.power",
	"stat.infer.effect_size",
	"stat.reg.multiple",
	"stat.reg.prediction_interval",
	"stat.test.kruskal_wallis",
}

type w2Rat struct {
	num int
	den int
}

func w2RatAdd(a, b w2Rat) w2Rat {
	return w2Rat{a.num*b.den + b.num*a.den, a.den * b.den}
}

func w2RatSub(a, b w2Rat) w2Rat {
	return w2Rat{a.num*b.den - b.num*a.den, a.den * b.den}
}

func w2RatEq(a, b w2Rat) bool {
	return a.num*b.den == b.num*a.den
}

func w2TestGCD(a, b int) int {
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

func w2TestFraction(num, den int) string {
	if den < 0 {
		num, den = -num, -den
	}
	g := w2TestGCD(num, den)
	if g == 0 {
		return "0"
	}
	num /= g
	den /= g
	if den == 1 {
		return strconv.Itoa(num)
	}
	return strconv.Itoa(num) + "/" + strconv.Itoa(den)
}

func w2TestDecimal(x float64, places int) string {
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

func w2TestInt(t *testing.T, facts map[string]string, key string) int {
	t.Helper()
	v, err := strconv.Atoi(facts[key])
	if err != nil {
		t.Fatalf("fact %q is not an integer: %q", key, facts[key])
	}
	return v
}

func w2TestFloat(t *testing.T, facts map[string]string, key string) float64 {
	t.Helper()
	v, err := strconv.ParseFloat(facts[key], 64)
	if err != nil {
		t.Fatalf("fact %q is not a number: %q", key, facts[key])
	}
	return v
}

func w2TestIntList(t *testing.T, facts map[string]string, key string) []int {
	t.Helper()
	raw := facts[key]
	if raw == "" {
		t.Fatalf("fact %q is empty", key)
	}
	parts := strings.Split(raw, ";")
	out := make([]int, len(parts))
	for i, part := range parts {
		v, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			t.Fatalf("fact %q has a non-integer member %q", key, part)
		}
		out[i] = v
	}
	return out
}

func w2TestGroups(t *testing.T, facts map[string]string) [][]int {
	t.Helper()
	raw := facts["groups"]
	if raw == "" {
		t.Fatal("groups fact is empty")
	}
	groupParts := strings.Split(raw, "|")
	groups := make([][]int, len(groupParts))
	for i, part := range groupParts {
		members := strings.Split(part, ",")
		groups[i] = make([]int, len(members))
		for j, member := range members {
			v, err := strconv.Atoi(strings.TrimSpace(member))
			if err != nil {
				t.Fatalf("group %d has a non-integer member %q", i+1, member)
			}
			groups[i][j] = v
		}
	}
	return groups
}

func w2TestChoose(n, k int) int {
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

func w2TestPow(base, exp int) int {
	out := 1
	for i := 0; i < exp; i++ {
		out *= base
	}
	return out
}

func w2TestRankSums(groups [][]int) []int {
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
	sort.Slice(items, func(i, j int) bool { return items[i].value < items[j].value })
	sums := make([]int, len(groups))
	for rank, it := range items {
		sums[it.group] += rank + 1
	}
	return sums
}

func w2TestKruskalH(groups [][]int) float64 {
	sums := w2TestRankSums(groups)
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

func w2Recompute(t *testing.T, id string, ordinal int, facts map[string]string) string {
	t.Helper()
	switch id {
	case "stat.data.weighted_mean":
		switch ordinal {
		case 1:
			values := w2TestIntList(t, facts, "values")
			weights := w2TestIntList(t, facts, "weights")
			if len(values) != len(weights) {
				t.Fatalf("values and weights differ in length")
			}
			total, weight := 0, 0
			for i := range values {
				total += values[i] * weights[i]
				weight += weights[i]
			}
			if total%weight != 0 {
				t.Fatalf("weighted mean is not an integer: %d/%d", total, weight)
			}
			return strconv.Itoa(total / weight)
		case 2:
			known := w2TestIntList(t, facts, "known_values")
			weights := w2TestIntList(t, facts, "known_weights")
			missingWeight := w2TestInt(t, facts, "missing_weight")
			target := w2TestInt(t, facts, "target")
			if len(known) != len(weights) {
				t.Fatalf("known values and weights differ in length")
			}
			totalWeight := missingWeight
			knownTotal := 0
			for i := range known {
				knownTotal += known[i] * weights[i]
				totalWeight += weights[i]
			}
			numerator := target*totalWeight - knownTotal
			if numerator%missingWeight != 0 {
				t.Fatalf("missing score is not an integer: %d/%d", numerator, missingWeight)
			}
			return strconv.Itoa(numerator / missingWeight)
		default:
			values := w2TestIntList(t, facts, "values")
			weights := w2TestIntList(t, facts, "weights")
			total, weight := 0, 0
			for i := range values {
				total += values[i] * weights[i]
				weight += weights[i]
			}
			weighted := float64(total) / float64(weight)
			unweighted := float64(sumIntsW2(values)) / float64(len(values))
			return w2TestDecimal(weighted-unweighted, 2)
		}

	case "stat.adv.coefficient_variation":
		switch ordinal {
		case 1:
			mean := w2TestFloat(t, facts, "mean")
			sd := w2TestFloat(t, facts, "sd")
			return strconv.Itoa(int(math.Round(sd / mean * 100)))
		case 2:
			mean := w2TestFloat(t, facts, "mean")
			cv := w2TestFloat(t, facts, "cv_percent")
			return strconv.Itoa(int(math.Round(mean * cv / 100)))
		default:
			meanA := w2TestFloat(t, facts, "mean_a")
			sdA := w2TestFloat(t, facts, "sd_a")
			meanB := w2TestFloat(t, facts, "mean_b")
			sdB := w2TestFloat(t, facts, "sd_b")
			cvA := math.Round(sdA / meanA * 100)
			cvB := math.Round(sdB / meanB * 100)
			return w2TestDecimal(cvA/cvB, 2)
		}

	case "stat.prob.independence":
		switch ordinal {
		case 1:
			pa := w2Rat{w2TestInt(t, facts, "p_a_num"), w2TestInt(t, facts, "p_a_den")}
			pb := w2Rat{w2TestInt(t, facts, "p_b_num"), w2TestInt(t, facts, "p_b_den")}
			pab := w2Rat{w2TestInt(t, facts, "p_ab_num"), w2TestInt(t, facts, "p_ab_den")}
			if w2RatEq(pab, w2Rat{pa.num * pb.num, pa.den * pb.den}) {
				return "Independent"
			}
			return "Not independent"
		case 2:
			pa := w2Rat{w2TestInt(t, facts, "p_a_num"), w2TestInt(t, facts, "p_a_den")}
			cond := w2Rat{w2TestInt(t, facts, "p_a_given_b_num"), w2TestInt(t, facts, "p_a_given_b_den")}
			if w2RatEq(pa, cond) {
				return "Independent"
			}
			return "Not independent"
		default:
			pa := w2Rat{w2TestInt(t, facts, "p_a_num"), w2TestInt(t, facts, "p_a_den")}
			pb := w2Rat{w2TestInt(t, facts, "p_b_num"), w2TestInt(t, facts, "p_b_den")}
			union := w2Rat{w2TestInt(t, facts, "p_a_or_b_num"), w2TestInt(t, facts, "p_a_or_b_den")}
			intersection := w2RatSub(w2RatAdd(pa, pb), union)
			if w2RatEq(intersection, w2Rat{pa.num * pb.num, pa.den * pb.den}) {
				return "Independent"
			}
			return "Not independent"
		}

	case "stat.prob.mutually_exclusive":
		switch ordinal {
		case 1:
			if w2TestInt(t, facts, "p_ab_num") == 0 {
				return "Mutually exclusive"
			}
			return "Not mutually exclusive"
		case 2:
			pa := w2Rat{w2TestInt(t, facts, "p_a_num"), w2TestInt(t, facts, "p_a_den")}
			pb := w2Rat{w2TestInt(t, facts, "p_b_num"), w2TestInt(t, facts, "p_b_den")}
			union := w2Rat{w2TestInt(t, facts, "p_a_or_b_num"), w2TestInt(t, facts, "p_a_or_b_den")}
			if w2RatEq(union, w2RatAdd(pa, pb)) {
				return "Mutually exclusive"
			}
			return "Not mutually exclusive"
		default:
			pa := w2Rat{w2TestInt(t, facts, "p_a_num"), w2TestInt(t, facts, "p_a_den")}
			pb := w2Rat{w2TestInt(t, facts, "p_b_num"), w2TestInt(t, facts, "p_b_den")}
			sum := w2RatAdd(pa, pb)
			if sum.num <= sum.den {
				return "Possible"
			}
			return "Not possible"
		}

	case "stat.dist.negative_binomial":
		pNum := w2TestInt(t, facts, "p_num")
		pDen := w2TestInt(t, facts, "p_den")
		r := w2TestInt(t, facts, "r")
		switch ordinal {
		case 1:
			if r*pDen%pNum != 0 {
				t.Fatalf("expected trials is not an integer")
			}
			return strconv.Itoa(r * pDen / pNum)
		case 2:
			k := w2TestInt(t, facts, "k")
			num := w2TestChoose(k-1, r-1) * w2TestPow(pNum, r) * w2TestPow(pDen-pNum, k-r)
			den := w2TestPow(pDen, k)
			return w2TestFraction(num, den)
		default:
			if r*(pDen-pNum)%pNum != 0 {
				t.Fatalf("expected failures is not an integer")
			}
			return strconv.Itoa(r * (pDen - pNum) / pNum)
		}

	case "stat.dist.multinomial":
		n := w2TestInt(t, facts, "n")
		switch ordinal {
		case 1:
			counts := w2TestIntList(t, facts, "counts")
			coeff := w2TestChoose(n, counts[0]) * w2TestChoose(n-counts[0], counts[1])
			return strconv.Itoa(coeff)
		case 2:
			pNum := w2TestInt(t, facts, "p_num")
			pDen := w2TestInt(t, facts, "p_den")
			if n*pNum%pDen != 0 {
				t.Fatalf("expected count is not an integer")
			}
			return strconv.Itoa(n * pNum / pDen)
		default:
			counts := w2TestIntList(t, facts, "counts")
			pNum := w2TestInt(t, facts, "p_num")
			pDen := w2TestInt(t, facts, "p_den")
			coeff := w2TestChoose(n, counts[0]) * w2TestChoose(n-counts[0], counts[1])
			num := coeff * w2TestPow(pNum, n)
			den := w2TestPow(pDen, n)
			return w2TestFraction(num, den)
		}

	case "stat.dist.lognormal":
		mu := w2TestFloat(t, facts, "mu")
		sigma := w2TestFloat(t, facts, "sigma")
		switch ordinal {
		case 1:
			return w2TestDecimal(math.Exp(mu), 4)
		case 2:
			return w2TestDecimal(math.Exp(mu+sigma*sigma/2), 4)
		default:
			return w2TestDecimal(math.Exp(sigma*sigma/2), 4)
		}

	case "stat.infer.ci_mean_t":
		xbar := w2TestFloat(t, facts, "xbar")
		s := w2TestFloat(t, facts, "s")
		n := w2TestFloat(t, facts, "n")
		df := w2TestInt(t, facts, "df")
		tCrit := w2TestFloat(t, facts, "t_crit")
		if want, ok := w2TestTTable[df]; !ok || math.Abs(want-tCrit) > 1e-9 {
			t.Fatalf("t critical value %v does not match the published table for df %d", tCrit, df)
		}
		margin := tCrit * (s / math.Sqrt(n))
		switch ordinal {
		case 1:
			return w2TestDecimal(s/math.Sqrt(n), 3)
		case 2:
			return w2TestDecimal(margin, 3)
		default:
			return w2TestDecimal(xbar-margin, 3)
		}

	case "stat.infer.ci_proportion":
		phatNum := w2TestFloat(t, facts, "phat_num")
		phatDen := w2TestFloat(t, facts, "phat_den")
		n := w2TestFloat(t, facts, "n")
		z := w2TestFloat(t, facts, "z_crit")
		phat := phatNum / phatDen
		margin := z * math.Sqrt(phat*(1-phat)/n)
		switch ordinal {
		case 1:
			return w2TestDecimal(math.Sqrt(phat*(1-phat)/n), 4)
		case 2:
			return w2TestDecimal(margin, 4)
		default:
			return w2TestDecimal(phat-margin, 4)
		}

	case "stat.infer.power":
		switch ordinal {
		case 1:
			return w2TestDecimal(1-w2TestFloat(t, facts, "beta"), 2)
		case 2:
			return w2TestDecimal(1-w2TestFloat(t, facts, "power"), 2)
		default:
			nullMean := w2TestFloat(t, facts, "null_mean")
			altMean := w2TestFloat(t, facts, "alt_mean")
			se := w2TestFloat(t, facts, "se")
			critical := w2TestFloat(t, facts, "critical")
			z := math.Round((critical - altMean) / se)
			if math.Abs((critical-altMean)/se-z) > 1e-9 {
				t.Fatalf("power case does not standardize to a whole z: %v", (critical-altMean)/se)
			}
			if got := w2TestInt(t, facts, "z_alt"); float64(got) != z {
				t.Fatalf("z_alt fact %d does not match recomputed %v", got, z)
			}
			tail, ok := map[int]string{-2: "0.9772", -1: "0.8413", 0: "0.5"}[int(z)]
			if !ok {
				t.Fatalf("no published tail probability for z = %v", z)
			}
			_ = nullMean
			return tail
		}

	case "stat.infer.effect_size":
		switch ordinal {
		case 1:
			xbar := w2TestFloat(t, facts, "xbar")
			mu := w2TestFloat(t, facts, "mu")
			s := w2TestFloat(t, facts, "s")
			return w2TestDecimal((xbar-mu)/s, 3)
		case 2:
			d := w2TestFloat(t, facts, "d")
			s := w2TestFloat(t, facts, "s")
			return w2TestDecimal(d*s, 3)
		default:
			mean1 := w2TestFloat(t, facts, "mean1")
			mean2 := w2TestFloat(t, facts, "mean2")
			s1 := w2TestFloat(t, facts, "s1")
			s2 := w2TestFloat(t, facts, "s2")
			n1 := w2TestFloat(t, facts, "n1")
			n2 := w2TestFloat(t, facts, "n2")
			pooled := math.Sqrt(((n1-1)*s1*s1 + (n2-1)*s2*s2) / (n1 + n2 - 2))
			return w2TestDecimal((mean1-mean2)/pooled, 3)
		}

	case "stat.reg.multiple":
		switch ordinal {
		case 1, 2:
			b0 := w2TestFloat(t, facts, "intercept")
			b1 := w2TestFloat(t, facts, "b1")
			b2 := w2TestFloat(t, facts, "b2")
			x1 := w2TestFloat(t, facts, "x1")
			x2 := w2TestFloat(t, facts, "x2")
			predicted := b0 + b1*x1 + b2*x2
			if ordinal == 1 {
				return w2TestDecimal(predicted, 0)
			}
			observed := w2TestFloat(t, facts, "observed")
			return w2TestDecimal(observed-predicted, 0)
		default:
			r2 := w2TestFloat(t, facts, "r2")
			n := w2TestFloat(t, facts, "n")
			predictors := w2TestFloat(t, facts, "predictors")
			return w2TestDecimal(1-(1-r2)*(n-1)/(n-predictors-1), 4)
		}

	case "stat.reg.prediction_interval":
		switch ordinal {
		case 1:
			tCrit := w2TestFloat(t, facts, "t_crit")
			sPred := w2TestFloat(t, facts, "s_pred")
			return w2TestDecimal(tCrit*sPred, 2)
		case 2:
			yhat := w2TestFloat(t, facts, "yhat")
			tCrit := w2TestFloat(t, facts, "t_crit")
			sPred := w2TestFloat(t, facts, "s_pred")
			return w2TestDecimal(yhat+tCrit*sPred, 2)
		default:
			se := w2TestFloat(t, facts, "s_e")
			sm := w2TestFloat(t, facts, "s_mean")
			return w2TestDecimal(math.Sqrt(se*se+sm*sm), 2)
		}

	case "stat.test.kruskal_wallis":
		switch ordinal {
		case 1:
			groups := w2TestGroups(t, facts)
			index := w2TestInt(t, facts, "group_index")
			sums := w2TestRankSums(groups)
			return strconv.Itoa(sums[index-1])
		case 2:
			groups := w2TestGroups(t, facts)
			return w2TestDecimal(w2TestKruskalH(groups), 2)
		default:
			return strconv.Itoa(w2TestInt(t, facts, "k") - 1)
		}
	}
	t.Fatalf("no recomputation is written for %s", id)
	return ""
}

func sumIntsW2(vals []int) int {
	total := 0
	for _, v := range vals {
		total += v
	}
	return total
}

// The two-sided 95% t critical values the generator is allowed to choose. The
// test checks the published fact against its own copy of that table rather than
// trusting a value the generator wrote.
var w2TestTTable = map[int]float64{
	1: 12.706, 2: 4.303, 3: 3.182, 4: 2.776, 5: 2.571,
	6: 2.447, 7: 2.365, 8: 2.306, 9: 2.262, 10: 2.228,
	12: 2.179, 15: 2.131, 20: 2.086, 25: 2.060, 30: 2.042,
	40: 2.021, 60: 2.000, 120: 1.980,
}

func TestWave2AnswersRecomputeFromFacts(t *testing.T) {
	reg := testStatRegistry(t)
	checked := 0
	for _, id := range wave2ConceptIDs {
		for ordinal := 1; ordinal <= 3; ordinal++ {
			for seed := 1; seed <= 300; seed++ {
				p, err := reg.GenerateContext(id, generator.GeneratorContext{
					Difficulty: 0.5,
					Seed:       int64(seed),
					KPOrdinal:  ordinal,
				})
				if err != nil {
					t.Fatalf("%s: generate failed: %v", id, err)
				}
				if p.Question == "" || p.Answer == "" || p.Explanation == "" {
					t.Fatalf("%s kp%d seed %d: empty problem field", id, ordinal, seed)
				}
				if p.KPOrdinal != ordinal {
					t.Fatalf("%s kp%d seed %d: problem records kp%d", id, ordinal, seed, p.KPOrdinal)
				}
				want := w2Recompute(t, id, ordinal, p.Facts)
				if p.Answer != want {
					t.Fatalf("%s kp%d seed %d: answer %q does not recompute to %q\nquestion: %s\nfacts: %v",
						id, ordinal, seed, p.Answer, want, p.Question, p.Facts)
				}
				checked++
			}
		}
	}
	if checked != len(wave2ConceptIDs)*3*300 {
		t.Fatalf("checked %d problems, want %d", checked, len(wave2ConceptIDs)*3*300)
	}
}

func TestWave2ExplanationsStateTheMathematics(t *testing.T) {
	reg := testStatRegistry(t)
	cues := map[string]string{
		"stat.data.weighted_mean":        "weighted mean",
		"stat.adv.coefficient_variation": "coefficient of variation",
		"stat.prob.independence":         "independent",
		"stat.prob.mutually_exclusive":   "mutually exclusive",
		"stat.dist.negative_binomial":    "negative binomial",
		"stat.dist.multinomial":          "multinomial",
		"stat.dist.lognormal":            "lognormal",
		"stat.infer.ci_mean_t":           "standard error",
		"stat.infer.ci_proportion":       "standard error",
		"stat.infer.power":               "power",
		"stat.infer.effect_size":         "cohen's d",
		"stat.reg.multiple":              "multiple regression",
		"stat.reg.prediction_interval":   "prediction interval",
		"stat.test.kruskal_wallis":       "kruskal-wallis",
	}
	for _, id := range wave2ConceptIDs {
		for ordinal := 1; ordinal <= 3; ordinal++ {
			for seed := 1; seed <= 40; seed++ {
				p, err := reg.GenerateContext(id, generator.GeneratorContext{
					Difficulty: 0.5,
					Seed:       int64(1000 + seed),
					KPOrdinal:  ordinal,
				})
				if err != nil {
					t.Fatalf("%s: generate failed: %v", id, err)
				}
				if !strings.Contains(strings.ToLower(p.Explanation), cues[id]) {
					t.Fatalf("%s kp%d: explanation does not name the concept: %s", id, ordinal, p.Explanation)
				}
			}
		}
	}
}

func TestWave2OrdinalsAskDifferentQuestions(t *testing.T) {
	reg := testStatRegistry(t)
	for _, id := range wave2ConceptIDs {
		questions := make([]string, 3)
		for ordinal := 1; ordinal <= 3; ordinal++ {
			p, err := reg.GenerateContext(id, generator.GeneratorContext{
				Difficulty: 0.5,
				Seed:       4242,
				KPOrdinal:  ordinal,
			})
			if err != nil {
				t.Fatalf("%s: generate failed: %v", id, err)
			}
			questions[ordinal-1] = p.Question
		}
		if questions[0] == questions[1] || questions[1] == questions[2] || questions[0] == questions[2] {
			t.Fatalf("%s: two ordinals share a seeded question", id)
		}
	}
}

func TestWave2QuestionDiversity(t *testing.T) {
	reg := testStatRegistry(t)
	for _, id := range wave2ConceptIDs {
		for ordinal := 1; ordinal <= 3; ordinal++ {
			problems, err := reg.BatchGenerateContext(id, 12, generator.GeneratorContext{
				Difficulty: 0.5,
				Seed:       int64(5000 + ordinal),
				KPOrdinal:  ordinal,
			})
			if err != nil {
				t.Fatalf("%s: batch generate failed: %v", id, err)
			}
			seen := map[string]bool{}
			for _, p := range problems {
				seen[p.Question] = true
			}
			if len(seen) < 5 {
				t.Fatalf("%s kp%d: only %d distinct questions in 12 draws", id, ordinal, len(seen))
			}
		}
	}
}
