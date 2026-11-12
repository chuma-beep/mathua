package realanalysis

import (
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

var w6aKinds = []string{
	"dedekind", "well_ordering", "cardinal", "root_test", "integral_test",
	"cesaro", "dirichlet", "oscillation", "semicontinuity", "dini",
	"uniform_limit", "rolle", "convex", "jensen", "lipschitz", "inverse_function",
}

func w6aTestInt(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer\nfacts: %v", p.Question, key, p.Facts[key], p.Facts)
	}
	return v
}

func w6aTestOrder(t *testing.T, p generator.Problem) []int {
	t.Helper()
	raw := p.Facts["order"]
	parts := strings.Split(raw, ", ")
	vals := make([]int, len(parts))
	for i, s := range parts {
		v, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil {
			t.Fatalf("%s: order %q is not a list of integers", p.Question, raw)
		}
		vals[i] = v
	}
	return vals
}

func w6aTestIndex(t *testing.T, vals []int, target int) int {
	t.Helper()
	for i, v := range vals {
		if v == target {
			return i
		}
	}
	t.Fatalf("value %d is not in the order %v", target, vals)
	return -1
}

func TestWave6ARealAnswers(t *testing.T) {
	for _, kind := range w6aKinds {
		for kp := 1; kp <= 3; kp++ {
			for draw := 0; draw < 30; draw++ {
				p := (&w6aGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				if p.KPOrdinal != kp {
					t.Fatalf("%s kp%d: KPOrdinal = %d", kind, kp, p.KPOrdinal)
				}
				if !strings.HasSuffix(p.Question, "(enter a number)") {
					t.Fatalf("%s kp%d: question does not ask for a number: %q", kind, kp, p.Question)
				}
				if p.Facts["answer"] != p.Answer {
					t.Fatalf("%s kp%d: fact answer = %q, answer = %q", kind, kp, p.Facts["answer"], p.Answer)
				}
				want := w6aTestAnswer(t, kind, kp, p)
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("%s kp%d: answer = %s, want %d\nquestion: %s\nfacts: %v", kind, kp, p.Answer, want, p.Question, p.Facts)
				}
			}
		}
	}
}

func w6aTestAnswer(t *testing.T, kind string, kp int, p generator.Problem) int {
	t.Helper()
	switch kind {
	case "dedekind":
		k := w6aTestInt(t, p, "k")
		switch kp {
		case 1:
			return k + 1
		case 2:
			return k
		}
		return 2*k + 1
	case "well_ordering":
		vals := w6aTestOrder(t, p)
		switch kp {
		case 1:
			return w6aTestIndex(t, vals, w6aTestInt(t, p, "target")) + 1
		case 2:
			best := vals[0]
			for _, v := range vals {
				if v > best {
					best = v
				}
			}
			return w6aTestIndex(t, vals, best) + 1
		}
		i := w6aTestIndex(t, vals, w6aTestInt(t, p, "first"))
		j := w6aTestIndex(t, vals, w6aTestInt(t, p, "second"))
		if i > j {
			i, j = j, i
		}
		return j - i - 1
	case "cardinal":
		switch kp {
		case 1:
			return 1 << w6aTestInt(t, p, "n")
		case 2:
			return w6aTestInt(t, p, "a") * w6aTestInt(t, p, "b")
		}
		return w6aTestInt(t, p, "a") + w6aTestInt(t, p, "b") - w6aTestInt(t, p, "overlap")
	case "root_test":
		base := w6aTestInt(t, p, "base")
		return base * base
	case "integral_test":
		switch kp {
		case 1:
			n := w6aTestInt(t, p, "n")
			return 2 * n * n
		case 2:
			return w6aTestInt(t, p, "p") - 1
		}
		return w6aTestInt(t, p, "bound_denominator") + 1
	case "cesaro":
		if kp == 1 {
			return w6aTestInt(t, p, "value") * w6aTestInt(t, p, "terms")
		}
		a := w6aTestInt(t, p, "a")
		d := w6aTestInt(t, p, "d")
		n := w6aTestInt(t, p, "terms")
		if kp == 2 {
			return a + d*(n-1)/2
		}
		return n*a + d*n*(n-1)/2
	case "dirichlet":
		switch kp {
		case 1:
			return w6aTestInt(t, p, "partial_bound") * w6aTestInt(t, p, "first_term")
		case 2:
			n := w6aTestInt(t, p, "terms")
			return (n + 1) * (n + 1)
		}
		return w6aTestInt(t, p, "terms") / 2 * w6aTestInt(t, p, "difference")
	case "oscillation":
		switch kp {
		case 1:
			return w6aTestInt(t, p, "value_right") - w6aTestInt(t, p, "value_left")
		case 2:
			return w6aTestInt(t, p, "value_high") - w6aTestInt(t, p, "value_low")
		}
		a := w6aTestInt(t, p, "a")
		b := w6aTestInt(t, p, "b")
		return b*b - a*a
	case "semicontinuity":
		left := w6aTestInt(t, p, "left_limit")
		right := w6aTestInt(t, p, "right_limit")
		switch kp {
		case 1:
			if left < right {
				return left
			}
			return right
		case 2:
			if left > right {
				return left
			}
			return right
		}
		return w6aTestInt(t, p, "right_offset") - w6aTestInt(t, p, "left_offset")
	case "dini":
		a := w6aTestInt(t, p, "endpoint")
		switch kp {
		case 1:
			return a / w6aTestInt(t, p, "index")
		case 2:
			n := w6aTestInt(t, p, "index")
			return a / (n * n)
		}
		return a + 1
	case "uniform_limit":
		a := w6aTestInt(t, p, "endpoint")
		if kp == 3 {
			return 2 * a
		}
		return a / w6aTestInt(t, p, "index")
	case "rolle":
		a := w6aTestInt(t, p, "root1")
		b := w6aTestInt(t, p, "root2")
		if kp == 3 {
			half := (b - a) / 2
			return -half * half
		}
		return (a + b) / 2
	case "convex":
		a := w6aTestInt(t, p, "a")
		b := w6aTestInt(t, p, "b")
		switch kp {
		case 1:
			half := (b - a) / 2
			return half * half
		case 2:
			return (a - b) * (a - b)
		}
		return (a*a + b*b) / 2
	case "jensen":
		pw := w6aTestInt(t, p, "weight_a")
		qw := w6aTestInt(t, p, "weight_b")
		d := w6aTestInt(t, p, "denominator")
		a := w6aTestInt(t, p, "a")
		b := w6aTestInt(t, p, "b")
		k := (b - a) / d
		return pw * qw * k * k
	case "lipschitz":
		switch kp {
		case 1:
			slope := (w6aTestInt(t, p, "f2") - w6aTestInt(t, p, "f0")) / 2
			if slope < 0 {
				slope = -slope
			}
			return slope
		case 2:
			return 2 * w6aTestInt(t, p, "endpoint")
		}
		return 2 * w6aTestInt(t, p, "b")
	case "inverse_function":
		switch kp {
		case 1:
			return (w6aTestInt(t, p, "f2") - w6aTestInt(t, p, "f0")) / 2
		case 2:
			return 2 * w6aTestInt(t, p, "x0")
		}
		a := w6aTestInt(t, p, "a")
		c := w6aTestInt(t, p, "c")
		x0 := w6aTestInt(t, p, "x0")
		return 3*a*x0*x0 + c
	}
	t.Fatalf("unknown kind %q", kind)
	return 0
}

func TestWave6ADiversity(t *testing.T) {
	if len(w6aKinds) != 16 {
		t.Fatalf("expected 16 kinds, got %d", len(w6aKinds))
	}
	for _, kind := range w6aKinds {
		seen := map[string]bool{}
		for i := 0; i < 200; i++ {
			kp := 1 + i%3
			p := (&w6aGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			seen[p.Answer] = true
		}
		if len(seen) < 8 {
			answers := make([]string, 0, len(seen))
			for a := range seen {
				answers = append(answers, a)
			}
			t.Fatalf("%s: only %d distinct answers in 200 draws: %v", kind, len(seen), answers)
		}
	}
}
