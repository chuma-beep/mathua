package statistics

import (
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func w4TestInt(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer\nfacts: %v", p.Question, key, p.Facts[key], p.Facts)
	}
	return v
}

func w4TestInts(t *testing.T, p generator.Problem, key string) []int {
	t.Helper()
	parts := strings.Split(p.Facts[key], ",")
	out := make([]int, len(parts))
	for i, s := range parts {
		v, err := strconv.Atoi(s)
		if err != nil {
			t.Fatalf("%s: fact %q contains %q, not an integer", p.Question, key, s)
		}
		out[i] = v
	}
	return out
}

func w4TestSum(xs []int) int {
	s := 0
	for _, x := range xs {
		s += x
	}
	return s
}

func w4TestAbs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func w4TestChoose(n, k int) int {
	if k > n-k {
		k = n - k
	}
	r := 1
	for i := 1; i <= k; i++ {
		r = r * (n - k + i) / i
	}
	return r
}

func TestWave4StatisticsAnswerRecomputation(t *testing.T) {
	kinds := []string{
		"bootstrap", "permutation", "wilcoxon", "mann_whitney", "kaplan_meier",
		"credible_interval", "information_criteria", "quantile_regression",
		"poisson_regression", "missing_data",
	}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			for draw := 0; draw < 40; draw++ {
				p := (&w4StatGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				want := w4TestStatisticsAnswer(t, kind, kp, p)
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("%s kp%d: answer = %s, want %d\nquestion: %s\nfacts: %v", kind, kp, p.Answer, want, p.Question, p.Facts)
				}
			}
		}
	}
}

func w4TestStatisticsAnswer(t *testing.T, kind string, kp int, p generator.Problem) int {
	t.Helper()
	switch kind {
	case "bootstrap":
		if kp == 1 {
			xs := w4TestInts(t, p, "sample")
			idx := w4TestInts(t, p, "indices")
			s := 0
			for _, i := range idx {
				s += xs[i-1]
			}
			return s / 4
		}
		if kp == 2 {
			b := w4TestInt(t, p, "B")
			return (b + 39) / 40
		}
		sd := w4TestInt(t, p, "sd")
		return sd * sd
	case "permutation":
		g1, g2 := w4TestInts(t, p, "group1"), w4TestInts(t, p, "group2")
		diff := w4TestSum(g1)/len(g1) - w4TestSum(g2)/len(g2)
		if kp == 1 {
			return diff
		}
		if kp == 2 {
			n, r := w4TestInt(t, p, "n"), w4TestInt(t, p, "group_size")
			return w4TestChoose(n, r)
		}
		stats := w4TestInts(t, p, "permutation_stats")
		c := 0
		for _, s := range stats {
			if w4TestAbs(s) >= w4TestAbs(diff) {
				c++
			}
		}
		return c
	case "wilcoxon":
		diffs := w4TestInts(t, p, "differences")
		type pair struct {
			abs, sign int
		}
		ps := make([]pair, len(diffs))
		for i, d := range diffs {
			ps[i] = pair{w4TestAbs(d), d}
		}
		sort.Slice(ps, func(i, j int) bool { return ps[i].abs < ps[j].abs })
		pos, neg := 0, 0
		for i, q := range ps {
			if q.sign > 0 {
				pos += i + 1
			} else {
				neg += i + 1
			}
		}
		if kp == 1 {
			return pos
		}
		if kp == 2 {
			return neg
		}
		if pos < neg {
			return pos
		}
		return neg
	case "mann_whitney":
		g1, g2 := w4TestInts(t, p, "group1"), w4TestInts(t, p, "group2")
		all := append(append([]int{}, g1...), g2...)
		sort.Ints(all)
		rank := map[int]int{}
		for i, v := range all {
			rank[v] = i + 1
		}
		r1 := 0
		for _, v := range g1 {
			r1 += rank[v]
		}
		u1 := r1 - len(g1)*(len(g1)+1)/2
		u2 := len(g1)*len(g2) - u1
		if kp == 1 {
			return r1
		}
		if kp == 2 {
			return u1
		}
		return u2
	case "kaplan_meier":
		times := w4TestInts(t, p, "times")
		if kp == 1 {
			// At the first event time, four subjects are at risk and one event occurs.
			_ = times
			return 3
		}
		if kp == 2 {
			return 3
		}
		// The survival curve reaches 1/2 at the second event time.
		return times[1]
	case "credible_interval":
		a0, b0 := w4TestInt(t, p, "prior_alpha"), w4TestInt(t, p, "prior_beta")
		s, f := w4TestInt(t, p, "successes"), w4TestInt(t, p, "failures")
		if kp == 1 {
			return a0 + s
		}
		if kp == 2 {
			return b0 + f
		}
		return a0 + s + b0 + f
	case "information_criteria":
		k, dev := w4TestInt(t, p, "parameters"), w4TestInt(t, p, "deviance")
		if kp == 1 {
			return dev + 2*k
		}
		if kp == 2 {
			return 2 * k
		}
		if dev+2*k < dev+2*(k+1) {
			return 1
		}
		return 2
	case "quantile_regression":
		slope, intercept, x := w4TestInt(t, p, "slope"), w4TestInt(t, p, "intercept"), w4TestInt(t, p, "x")
		if kp == 1 {
			return slope*x + intercept
		}
		if kp == 2 {
			return slope
		}
		return 2
	case "poisson_regression":
		a, b, x := w4TestInt(t, p, "intercept"), w4TestInt(t, p, "coefficient"), w4TestInt(t, p, "x")
		if kp == 1 {
			return a + b*x
		}
		if kp == 2 {
			return a
		}
		if b > 0 {
			return 1
		}
		return 0
	case "missing_data":
		if kp == 2 {
			xs := w4TestInts(t, p, "observed")
			return w4TestSum(xs) / len(xs)
		}
		return 1
	}
	t.Fatalf("unknown kind %q", kind)
	return 0
}

func TestWave4StatisticsQuestionDiversity(t *testing.T) {
	kinds := []string{
		"bootstrap", "permutation", "wilcoxon", "mann_whitney", "kaplan_meier",
		"credible_interval", "information_criteria", "quantile_regression",
		"poisson_regression", "missing_data",
	}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			seen := map[string]bool{}
			for i := 0; i < 60; i++ {
				p := (&w4StatGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				seen[p.Question] = true
			}
			if len(seen) < 5 {
				t.Fatalf("%s kp%d: only %d distinct questions in 60 draws", kind, kp, len(seen))
			}
		}
	}
}
