package precalculus

import (
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func w9TestInt(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer\nfacts: %v", p.Question, key, p.Facts[key], p.Facts)
	}
	return v
}

func w9TestInts(s string) []int {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]int, len(parts))
	for i, part := range parts {
		out[i], _ = strconv.Atoi(strings.TrimSpace(part))
	}
	return out
}

func TestWave9PreAnswerRecomputation(t *testing.T) {
	for _, kind := range []string{"one_to_one", "complex_factorization", "orthogonal_decomposition"} {
		for kp := 1; kp <= 3; kp++ {
			for draw := 0; draw < 40; draw++ {
				p := (&w9PreGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				want := w9PreAnswer(t, kind, kp, p)
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("%s kp%d: answer = %s, want %d\nquestion: %s\nfacts: %v", kind, kp, p.Answer, want, p.Question, p.Facts)
				}
			}
		}
	}
}

func w9PreAnswer(t *testing.T, kind string, kp int, p generator.Problem) int {
	t.Helper()
	switch kind {
	case "one_to_one":
		if kp == 1 {
			switch p.Facts["kind"] {
			case "cubic", "linear":
				return 1
			default:
				return 0
			}
		}
		if kp == 2 {
			ys := w9TestInts(p.Facts["ys"])
			seen := map[int]bool{}
			for _, y := range ys {
				seen[y] = true
			}
			if len(seen) == len(ys) {
				return 1
			}
			return 0
		}
		lo, hi := w9TestInt(t, p, "lo"), w9TestInt(t, p, "hi")
		if lo >= 0 || hi <= 0 {
			return 1
		}
		return 0
	case "complex_factorization":
		roots := w9TestInts(p.Facts["roots"])
		a, b := w9TestInt(t, p, "a"), w9TestInt(t, p, "b")
		sum, prod := 0, 1
		for _, r := range roots {
			sum += r
			prod *= r
		}
		switch kp {
		case 1:
			return sum + 2*a
		case 2:
			return prod * (a*a + b*b)
		default:
			return len(roots) + 2
		}
	case "orthogonal_decomposition":
		ax, ay := w9TestInt(t, p, "ax"), w9TestInt(t, p, "ay")
		k, m := w9TestInt(t, p, "k"), w9TestInt(t, p, "m")
		switch kp {
		case 1:
			return k * ax
		case 2:
			return -m * ay
		default:
			return k * (ax*ax + ay*ay)
		}
	}
	t.Fatalf("unknown kind %q", kind)
	return 0
}

func TestWave9PreQuestionDiversity(t *testing.T) {
	for _, kind := range []string{"one_to_one", "complex_factorization", "orthogonal_decomposition"} {
		for kp := 1; kp <= 3; kp++ {
			seen := map[string]bool{}
			for i := 0; i < 60; i++ {
				p := (&w9PreGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				seen[p.Question] = true
			}
			if len(seen) < 5 {
				t.Fatalf("%s kp%d: only %d distinct questions in 60 draws", kind, kp, len(seen))
			}
		}
	}
}
