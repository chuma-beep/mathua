package statistics

import (
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

var w7Kinds = map[string]string{
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

func w7Fact(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer\nfacts: %v", p.Question, key, p.Facts[key], p.Facts)
	}
	return v
}

func w7List(t *testing.T, p generator.Problem, key string) []int {
	t.Helper()
	raw := strings.TrimSpace(p.Facts[key])
	if raw == "" {
		t.Fatalf("%s: fact %q is empty", p.Question, key)
	}
	out := []int{}
	for _, part := range strings.Split(raw, ",") {
		v, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			t.Fatalf("%s: bad list value %q in %q", p.Question, part, key)
		}
		out = append(out, v)
	}
	return out
}

func w7Pow(base, exp int) int {
	v := 1
	for i := 0; i < exp; i++ {
		v *= base
	}
	return v
}

func w7Recompute(t *testing.T, p generator.Problem, kind string, kp int) int {
	switch kind {
	case "two_way_table":
		switch kp {
		case 1:
			return w7Fact(t, p, "a") + w7Fact(t, p, "b")
		case 2:
			return w7Fact(t, p, "b") + w7Fact(t, p, "d")
		default:
			return w7Fact(t, p, "a") + w7Fact(t, p, "b") + w7Fact(t, p, "c") + w7Fact(t, p, "d")
		}
	case "stem_leaf":
		vals := w7List(t, p, "values")
		sort.Ints(vals)
		switch kp {
		case 1:
			return len(vals)
		case 2:
			return vals[len(vals)-1] - vals[0]
		default:
			return vals[len(vals)/2]
		}
	case "cumulative_frequency":
		freq := w7List(t, p, "frequencies")
		pick := w7Fact(t, p, "pick")
		switch kp {
		case 1:
			sum := 0
			for i := 0; i < pick; i++ {
				sum += freq[i]
			}
			return sum
		case 2:
			sum := 0
			for _, v := range freq {
				sum += v
			}
			return sum
		default:
			return freq[pick-1]
		}
	case "percentile":
		vals := w7List(t, p, "values")
		n := w7Fact(t, p, "n")
		pr := w7Fact(t, p, "p")
		idx := (pr*n + 99) / 100
		switch kp {
		case 1:
			return idx
		case 2:
			return vals[idx-1]
		default:
			return idx - 1
		}
	case "trimmed_mean":
		vals := w7List(t, p, "values")
		sort.Ints(vals)
		middle := 0
		for i := 2; i < len(vals)-2; i++ {
			middle += vals[i]
		}
		total := 0
		for _, v := range vals {
			total += v
		}
		switch kp {
		case 1:
			return middle
		case 2:
			return total - middle
		default:
			return middle / (len(vals) - 4)
		}
	case "midrange":
		lo, hi := w7Fact(t, p, "min"), w7Fact(t, p, "max")
		switch kp {
		case 1:
			return (lo + hi) / 2
		case 2:
			return 2*w7Fact(t, p, "midrange") - lo
		default:
			return hi - lo
		}
	case "union":
		pa, pb := w7Fact(t, p, "pa"), w7Fact(t, p, "pb")
		inter, union := w7Fact(t, p, "intersection"), w7Fact(t, p, "union")
		switch kp {
		case 1:
			return pa + pb - inter
		case 2:
			return pa + pb - union
		default:
			return union - pa + inter
		}
	case "tree":
		p1 := w7Fact(t, p, "p") * w7Fact(t, p, "q") / 100
		p2 := w7Fact(t, p, "r") * w7Fact(t, p, "s") / 100
		switch kp {
		case 1:
			return p1
		case 2:
			return p1 + p2
		default:
			return 100 - p1
		}
	case "total_probability":
		p1, q := w7Fact(t, p, "p1"), w7Fact(t, p, "q")
		switch kp {
		case 1:
			return q
		case 2:
			return p1 * q / 100
		default:
			return 100 - p1 - w7Fact(t, p, "p2")
		}
	case "permutations":
		n, r := w7Fact(t, p, "n"), w7Fact(t, p, "r")
		v := 1
		for i := 0; i < r; i++ {
			v *= n - i
		}
		return v
	case "combinations":
		n, r := w7Fact(t, p, "n"), w7Fact(t, p, "r")
		v := 1
		for i := 0; i < r; i++ {
			v = v * (n - i) / (i + 1)
		}
		return v
	case "odds":
		a, b := w7Fact(t, p, "favourable"), w7Fact(t, p, "unfavourable")
		switch kp {
		case 1:
			return a * 100 / (a + b)
		case 2:
			return b * 100 / (a + b)
		default:
			return b
		}
	case "weibull":
		k, ratio, value := w7Fact(t, p, "k"), w7Fact(t, p, "ratio"), w7Fact(t, p, "value")
		switch kp {
		case 1:
			return w7Pow(ratio, k)
		case 2:
			for i := 1; i <= 10; i++ {
				if w7Pow(ratio, i) == value {
					return i
				}
			}
			t.Fatalf("no shape fits ratio %d and value %d", ratio, value)
		default:
			for i := 1; i <= 10; i++ {
				if w7Pow(i, k) == value {
					return i
				}
			}
			t.Fatalf("no ratio fits shape %d and value %d", k, value)
		}
	case "cauchy":
		switch kp {
		case 1:
			return 2 * w7Fact(t, p, "scale")
		case 2:
			return w7Fact(t, p, "iqr") / 2
		default:
			return 2 * w7Fact(t, p, "location")
		}
	case "triangular":
		a, b, c := w7Fact(t, p, "a"), w7Fact(t, p, "b"), w7Fact(t, p, "c")
		switch kp {
		case 1:
			return (a + b + c) / 3
		case 2:
			return c - a
		default:
			return c
		}
	case "pareto":
		alpha, xm, mean := w7Fact(t, p, "alpha"), w7Fact(t, p, "xm"), w7Fact(t, p, "mean")
		switch kp {
		case 1:
			return alpha * xm / (alpha - 1)
		case 2:
			return mean * (alpha - 1) / alpha
		default:
			for i := 2; i <= 10; i++ {
				if i*xm%(i-1) == 0 && i*xm/(i-1) == mean {
					return i
				}
			}
			t.Fatalf("no shape fits mean %d and minimum %d", mean, xm)
		}
	case "fisher_exact":
		r1, c1 := w7Fact(t, p, "row1"), w7Fact(t, p, "col1")
		_, c2 := w7Fact(t, p, "row1"), w7Fact(t, p, "col2")
		switch kp {
		case 1:
			return minInt(r1, c1)
		case 2:
			v := r1 - c2
			if v < 0 {
				v = 0
			}
			return v
		default:
			return r1 + w7Fact(t, p, "row2")
		}
	case "mcnemar":
		b, c := w7Fact(t, p, "b"), w7Fact(t, p, "c")
		switch kp {
		case 1:
			return (b - c) * (b - c)
		case 2:
			return b + c
		default:
			return b - c
		}
	case "sign_test":
		vals := w7List(t, p, "differences")
		pos := 0
		for _, v := range vals {
			if v > 0 {
				pos++
			}
		}
		switch kp {
		case 1:
			return pos
		case 2:
			return minInt(pos, len(vals)-pos)
		default:
			return len(vals) - pos
		}
	case "cooks_distance":
		vals := w7List(t, p, "distances")
		switch kp {
		case 1:
			thr := w7Fact(t, p, "threshold_hundredths")
			count := 0
			for _, v := range vals {
				if v > thr {
					count++
				}
			}
			return count
		case 2:
			best := 0
			for i, v := range vals {
				if v > vals[best] {
					best = i
				}
			}
			return best + 1
		default:
			return 400 / w7Fact(t, p, "n")
		}
	case "polynomial":
		d := w7Fact(t, p, "degree")
		switch kp {
		case 1:
			return d + 1
		case 2:
			return d
		default:
			return d + d*(d-1)/2 + 1
		}
	case "ridge":
		sxx, lambda, beta, sxy := w7Fact(t, p, "sxx"), w7Fact(t, p, "lambda"), w7Fact(t, p, "beta"), w7Fact(t, p, "sxy")
		switch kp {
		case 1:
			return sxy / (sxx + lambda)
		case 2:
			return sxy/beta - sxx
		default:
			return sxy/beta - lambda
		}
	}
	t.Fatalf("unknown kind %q", kind)
	return 0
}

func TestWave7StatisticsAnswers(t *testing.T) {
	for _, kind := range w7Kinds {
		for kp := 1; kp <= 3; kp++ {
			for draw := 0; draw < 30; draw++ {
				p := (&w7StatGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				if p.KPOrdinal != kp {
					t.Fatalf("%s kp%d: KPOrdinal = %d", kind, kp, p.KPOrdinal)
				}
				want := w7Recompute(t, p, kind, kp)
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("%s kp%d: answer = %s, want %d\nquestion: %s\nfacts: %v", kind, kp, p.Answer, want, p.Question, p.Facts)
				}
				if !strings.Contains(p.Question, "(enter a number)") {
					t.Fatalf("%s kp%d: question lacks the numeric prompt: %s", kind, kp, p.Question)
				}
			}
		}
	}
}

func TestWave7StatisticsDiversity(t *testing.T) {
	for _, kind := range w7Kinds {
		for kp := 1; kp <= 3; kp++ {
			seen := map[string]bool{}
			for i := 0; i < 200; i++ {
				p := (&w7StatGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				seen[p.Answer] = true
			}
			if len(seen) < 5 {
				t.Fatalf("%s kp%d: only %d distinct answers in 200 draws", kind, kp, len(seen))
			}
		}
	}
}
