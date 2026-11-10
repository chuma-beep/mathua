package fractions

import (
	"strconv"
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

func TestWave9FractionAnswerRecomputation(t *testing.T) {
	for kp := 1; kp <= 3; kp++ {
		for draw := 0; draw < 40; draw++ {
			p := (&w9MixedAddCarryGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			b, d := w9TestInt(t, p, "b"), w9TestInt(t, p, "d")
			l := lcm(b, d)
			fracNum := w9TestInt(t, p, "a")*l/b + w9TestInt(t, p, "c")*l/d
			carry := fracNum / l
			rem := fracNum % l
			whole := w9TestInt(t, p, "w1") + w9TestInt(t, p, "w2") + carry
			want := []int{carry, whole, rem}[kp-1]
			if p.Answer != strconv.Itoa(want) {
				t.Fatalf("kp%d: answer = %s, want %d\nquestion: %s\nfacts: %v", kp, p.Answer, want, p.Question, p.Facts)
			}
		}
	}
}

func TestWave9FractionQuestionDiversity(t *testing.T) {
	for kp := 1; kp <= 3; kp++ {
		seen := map[string]bool{}
		for i := 0; i < 60; i++ {
			p := (&w9MixedAddCarryGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			seen[p.Question] = true
		}
		if len(seen) < 5 {
			t.Fatalf("kp%d: only %d distinct questions in 60 draws", kp, len(seen))
		}
	}
}
