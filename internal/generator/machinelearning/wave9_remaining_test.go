package machinelearning

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

func TestWave9RemainingMLAnswerRecomputation(t *testing.T) {
	for kp := 1; kp <= 3; kp++ {
		for draw := 0; draw < 40; draw++ {
			p := (&w9rSMOTEGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			x1 := w9rTestInt(t, p, "x1")
			x2 := w9rTestInt(t, p, "x2")
			k := w9rTestInt(t, p, "k")
			n := w9rTestInt(t, p, "n")
			m := w9rTestInt(t, p, "m")
			want := []int{(x1 + x2) / 2, k, m + n}[kp-1]
			if p.Answer != strconv.Itoa(want) {
				t.Fatalf("kp%d: answer = %s, want %d\nquestion: %s\nfacts: %v", kp, p.Answer, want, p.Question, p.Facts)
			}
		}
	}
}

func TestWave9RemainingMLQuestionDiversity(t *testing.T) {
	for kp := 1; kp <= 3; kp++ {
		seen := map[string]bool{}
		for i := 0; i < 60; i++ {
			p := (&w9rSMOTEGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			seen[p.Question] = true
		}
		if len(seen) < 5 {
			t.Fatalf("kp%d: only %d distinct questions in 60 draws", kp, len(seen))
		}
		if len(seen) < 10 {
			t.Logf("kp%d: %d distinct questions in 60 draws, below the 10 target", kp, len(seen))
		}
	}
}
