package abstract

import (
	"strconv"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func w7TestInt(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer\nfacts: %v", p.Question, key, p.Facts[key], p.Facts)
	}
	return v
}

func w7TestGCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func TestWave7AbstractAnswerRecomputation(t *testing.T) {
	kinds := []string{"semidirect", "representation", "symmetric", "ext"}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			for draw := 0; draw < 40; draw++ {
				p := (&w7AbstractGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				want := w7TestAbstractAnswer(t, kind, kp, p)
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("%s kp%d: answer = %s, want %d\nquestion: %s\nfacts: %v", kind, kp, p.Answer, want, p.Question, p.Facts)
				}
			}
		}
	}
}

func w7TestAbstractAnswer(t *testing.T, kind string, kp int, p generator.Problem) int {
	t.Helper()
	switch kind {
	case "semidirect":
		n := w7TestInt(t, p, "n")
		if kp == 1 {
			return 2 * n
		}
		if kp == 2 {
			a := w7TestInt(t, p, "a")
			b := w7TestInt(t, p, "b")
			h1 := w7TestInt(t, p, "h1")
			if h1 == 0 {
				return (a + b) % n
			}
			return (a - b + n) % n
		}
		if n%2 == 0 {
			return 2
		}
		return 1
	case "representation":
		if kp == 1 {
			return w7TestInt(t, p, "order")
		}
		if kp == 2 {
			return w7TestInt(t, p, "a") + w7TestInt(t, p, "b") + w7TestInt(t, p, "c")
		}
		return w7TestInt(t, p, "n") - 2
	case "symmetric":
		a, b, c := w7TestInt(t, p, "a"), w7TestInt(t, p, "b"), w7TestInt(t, p, "c")
		if kp == 1 {
			return a*b + a*c + b*c
		}
		if kp == 2 {
			return a*a + b*b + c*c
		}
		return a * b * c
	case "ext":
		m, n := w7TestInt(t, p, "m"), w7TestInt(t, p, "n")
		if kp == 1 || kp == 2 {
			return w7TestGCD(m, n)
		}
		return w7TestInt(t, p, "dimA") + w7TestInt(t, p, "dimC")
	}
	t.Fatalf("unknown kind %q", kind)
	return 0
}

func TestWave7AbstractDiversity(t *testing.T) {
	kinds := []string{"semidirect", "representation", "symmetric", "ext"}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			seen := map[string]bool{}
			for i := 0; i < 60; i++ {
				p := (&w7AbstractGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
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
