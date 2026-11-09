package statistics

import (
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func w6TestInt(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer\nfacts: %v", p.Question, key, p.Facts[key], p.Facts)
	}
	return v
}

func w6TestPValues(t *testing.T, p generator.Problem) []int {
	t.Helper()
	var out []int
	for _, part := range strings.Split(p.Facts["pvalues"], ",") {
		v, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			t.Fatalf("bad p-value %q", part)
		}
		out = append(out, v)
	}
	return out
}

func TestWave6StatisticsAnswerRecomputation(t *testing.T) {
	for kp := 1; kp <= 3; kp++ {
		for draw := 0; draw < 30; draw++ {
			p := (&w6StatGen{kind: "multiple_comparisons"}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			m := w6TestInt(t, p, "m")
			pv := w6TestPValues(t, p)
			want := 0
			switch kp {
			case 1:
				for _, x := range pv {
					if x*m <= 5 {
						want++
					}
				}
			case 2:
				alpha := w6TestInt(t, p, "alpha_percent")
				want = alpha * m
			default:
				for i, x := range pv {
					if x*m <= 5*(i+1) {
						want = i + 1
					}
				}
			}
			if p.Answer != strconv.Itoa(want) {
				t.Fatalf("multiple_comparisons kp%d: answer = %s, want %d\nquestion: %s\nfacts: %v", kp, p.Answer, want, p.Question, p.Facts)
			}
		}
	}
}

func TestWave6StatisticsDiversity(t *testing.T) {
	for kp := 1; kp <= 3; kp++ {
		seen := map[string]bool{}
		for i := 0; i < 60; i++ {
			p := (&w6StatGen{kind: "multiple_comparisons"}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			seen[p.Question] = true
		}
		if len(seen) < 5 {
			t.Fatalf("multiple_comparisons kp%d: only %d distinct questions in 60 draws", kp, len(seen))
		}
		if len(seen) < 10 {
			t.Logf("multiple_comparisons kp%d: %d distinct questions in 60 draws, below the 10 target", kp, len(seen))
		}
	}
}
