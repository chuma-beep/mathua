package fractions

import (
	"strconv"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func w9rTestInt(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer\nfacts: %v", p.Question, key, p.Facts[key], p.Facts)
	}
	return v
}

func TestWave9RemainingFractionAnswerRecomputation(t *testing.T) {
	for draw := 0; draw < 40; draw++ {
		for kp := 1; kp <= 3; kp++ {
			p := (&w9rMixedOrderGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			w1 := w9rTestInt(t, p, "w1")
			a := w9rTestInt(t, p, "a")
			b := w9rTestInt(t, p, "b")
			w2 := w9rTestInt(t, p, "w2")
			c := w9rTestInt(t, p, "c")
			d := w9rTestInt(t, p, "d")
			left := (w1*b + a) * d
			right := (w2*d + c) * b
			if left == right {
				t.Fatalf("tie in mixed-order comparison: %v", p.Facts)
			}
			want := 0
			switch kp {
			case 1:
				want = 2
				if left > right {
					want = 1
				}
			case 2:
				want = w2
				if left > right {
					want = w1
				}
			default:
				if a*d > c*b {
					want = 1
				}
			}
			if p.Answer != strconv.Itoa(want) {
				t.Fatalf("mixed-order kp%d: answer = %s, want %d\nquestion: %s\nfacts: %v", kp, p.Answer, want, p.Question, p.Facts)
			}

			q := (&w9rPercentOfFractionGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			pp := w9rTestInt(t, q, "p")
			aa := w9rTestInt(t, q, "a")
			bb := w9rTestInt(t, q, "b")
			num := pp * aa
			den := 100 * bb
			want = 0
			switch kp {
			case 1:
				want = num
			case 2:
				want = den
			default:
				if 2*num > den {
					want = 1
				}
			}
			if q.Answer != strconv.Itoa(want) {
				t.Fatalf("percent kp%d: answer = %s, want %d\nquestion: %s\nfacts: %v", kp, q.Answer, want, q.Question, q.Facts)
			}
		}
	}
}

func TestWave9RemainingFractionQuestionDiversity(t *testing.T) {
	gens := []struct {
		name string
		gen  generator.Generator
	}{
		{"mixed order", &w9rMixedOrderGen{}},
		{"percent of fraction", &w9rPercentOfFractionGen{}},
	}
	for _, tc := range gens {
		for kp := 1; kp <= 3; kp++ {
			seen := map[string]bool{}
			for i := 0; i < 60; i++ {
				p := tc.gen.Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				seen[p.Question] = true
			}
			if len(seen) < 5 {
				t.Fatalf("%s kp%d: only %d distinct questions in 60 draws", tc.name, kp, len(seen))
			}
			if len(seen) < 10 {
				t.Logf("%s kp%d: %d distinct questions in 60 draws, below the 10 target", tc.name, kp, len(seen))
			}
		}
	}
}
