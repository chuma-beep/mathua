package fractions

import (
	"strconv"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func w8TestInt(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer\nfacts: %v", p.Question, key, p.Facts[key], p.Facts)
	}
	return v
}

func w8TestGCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func TestWave8FractionAnswers(t *testing.T) {
	for _, kind := range []string{"set", "same_numerator", "common_denominator"} {
		for kp := 1; kp <= 3; kp++ {
			for draw := 0; draw < 40; draw++ {
				p := (&w8FracGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				want := w8TestFractionAnswer(t, kind, kp, p)
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("%s kp%d: answer %s, want %d\nquestion: %s\nfacts: %v", kind, kp, p.Answer, want, p.Question, p.Facts)
				}
			}
		}
	}
}

func w8TestFractionAnswer(t *testing.T, kind string, kp int, p generator.Problem) int {
	t.Helper()
	switch kind {
	case "set":
		n, num, den := w8TestInt(t, p, "n"), w8TestInt(t, p, "num"), w8TestInt(t, p, "den")
		red := (n / den) * num
		switch kp {
		case 1:
			return red
		case 2:
			return n - red
		default:
			return n / den
		}
	case "same_numerator":
		a, b, c := w8TestInt(t, p, "a"), w8TestInt(t, p, "b"), w8TestInt(t, p, "c")
		switch kp {
		case 1:
			if b < c {
				return 1
			}
			return 2
		case 2:
			if b < c {
				return b
			}
			return c
		default:
			if a < b {
				return 1
			}
			return 0
		}
	case "common_denominator":
		a, b := w8TestInt(t, p, "a"), w8TestInt(t, p, "b")
		g := w8TestGCD(a, b)
		switch kp {
		case 1:
			return a / g
		case 2:
			return b / g
		default:
			return g
		}
	}
	t.Fatalf("unknown kind %q", kind)
	return 0
}

func TestWave8FractionDiversity(t *testing.T) {
	for _, kind := range []string{"set", "same_numerator", "common_denominator"} {
		for kp := 1; kp <= 3; kp++ {
			seen := map[string]bool{}
			for i := 0; i < 60; i++ {
				p := (&w8FracGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				seen[p.Question] = true
			}
			if len(seen) < 5 {
				t.Fatalf("%s kp%d: only %d distinct questions in 60 draws", kind, kp, len(seen))
			}
			if len(seen) < 10 {
				t.Logf("%s kp%d: %d distinct questions in 60 draws", kind, kp, len(seen))
			}
		}
	}
}
