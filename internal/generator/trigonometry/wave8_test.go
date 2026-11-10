package trigonometry

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

func w8TestCheb(n, x int) int {
	if n == 0 {
		return 1
	}
	if n == 1 {
		return x
	}
	a, b := 1, x
	for k := 2; k <= n; k++ {
		a, b = b, 2*x*b-a
	}
	return b
}

func TestWave8TrigAnswers(t *testing.T) {
	kinds := []string{"versine", "chebyshev", "harmonic", "arctan_sum"}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			for draw := 0; draw < 40; draw++ {
				p := (&w8TrigGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				want := w8TestTrigAnswer(t, kind, kp, p)
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("%s kp%d: answer %s, want %d\nquestion: %s\nfacts: %v", kind, kp, p.Answer, want, p.Question, p.Facts)
				}
			}
		}
	}
}

func w8TestTrigAnswer(t *testing.T, kind string, kp int, p generator.Problem) int {
	t.Helper()
	switch kind {
	case "versine":
		s, c, d := w8TestInt(t, p, "sin_num"), w8TestInt(t, p, "cos_num"), w8TestInt(t, p, "den")
		switch kp {
		case 1:
			return d - c
		case 2:
			return d - s
		default:
			return 2*d - s - c
		}
	case "chebyshev":
		n, x := w8TestInt(t, p, "n"), w8TestInt(t, p, "x")
		switch kp {
		case 1:
			return w8TestCheb(n, x)
		case 2:
			if n%2 == 0 {
				return 1
			}
			return -1
		default:
			return 1 << (n - 1)
		}
	case "harmonic":
		switch kp {
		case 1:
			return w8TestInt(t, p, "r")
		case 2:
			return w8TestInt(t, p, "a")
		default:
			return w8TestInt(t, p, "b")
		}
	case "arctan_sum":
		a, b := w8TestInt(t, p, "a"), w8TestInt(t, p, "b")
		switch kp {
		case 1:
			return a + b
		case 2:
			return 1 - a*b
		default:
			if a*b < 1 {
				return 0
			}
			return 1
		}
	}
	t.Fatalf("unknown kind %q", kind)
	return 0
}

func TestWave8TrigDiversity(t *testing.T) {
	for _, kind := range []string{"versine", "chebyshev", "harmonic", "arctan_sum"} {
		for kp := 1; kp <= 3; kp++ {
			seen := map[string]bool{}
			for i := 0; i < 60; i++ {
				p := (&w8TrigGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
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
