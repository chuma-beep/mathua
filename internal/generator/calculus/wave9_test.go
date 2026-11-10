package calculus

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

func w9TestPow(b, e int) int {
	r := 1
	for i := 0; i < e; i++ {
		r *= b
	}
	return r
}

func TestWave9CalcAnswerRecomputation(t *testing.T) {
	for _, kind := range []string{"asymptotic_equivalence", "elasticity", "leibniz", "cesaro", "dirichlet_series"} {
		for kp := 1; kp <= 3; kp++ {
			for draw := 0; draw < 40; draw++ {
				p := (&w9CalcGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				want := w9CalcAnswer(t, kind, kp, p)
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("%s kp%d: answer = %s, want %d\nquestion: %s\nfacts: %v", kind, kp, p.Answer, want, p.Question, p.Facts)
				}
			}
		}
	}
}

func w9CalcAnswer(t *testing.T, kind string, kp int, p generator.Problem) int {
	t.Helper()
	switch kind {
	case "asymptotic_equivalence":
		a, b := w9TestInt(t, p, "a"), w9TestInt(t, p, "b")
		if kp == 1 {
			return a / b
		}
		if kp == 2 {
			return a
		}
		// q > p in this family, so the second power grows faster.
		_ = w9TestInt(t, p, "q")
		return -1
	case "elasticity":
		k := w9TestInt(t, p, "k")
		if kp == 1 {
			return k
		}
		if kp == 2 {
			if k > 1 {
				return 1
			}
			if k < 1 {
				return 2
			}
			return 3
		}
		if k < 1 {
			return 1
		}
		if k > 1 {
			return -1
		}
		return 0
	case "leibniz":
		a, b := w9TestInt(t, p, "a"), w9TestInt(t, p, "b")
		if kp == 1 {
			return 4*a + 3*b
		}
		if kp == 2 {
			return 2*a*a + 3*a*b
		}
		return 2*a + b
	case "cesaro":
		if kp == 1 {
			return w9TestInt(t, p, "c") + (w9TestInt(t, p, "n")+1)/2
		}
		if kp == 2 {
			return w9TestInt(t, p, "c") - w9TestInt(t, p, "d")/w9TestInt(t, p, "n")
		}
		return w9TestInt(t, p, "c")
	case "dirichlet_series":
		if kp == 1 {
			return w9TestPow(w9TestInt(t, p, "n"), w9TestInt(t, p, "p"))
		}
		if kp == 2 {
			return w9TestInt(t, p, "p") + 1
		}
		return w9TestInt(t, p, "c") * w9TestInt(t, p, "N")
	}
	t.Fatalf("unknown kind %q", kind)
	return 0
}

func TestWave9CalcQuestionDiversity(t *testing.T) {
	for _, kind := range []string{"asymptotic_equivalence", "elasticity", "leibniz", "cesaro", "dirichlet_series"} {
		for kp := 1; kp <= 3; kp++ {
			seen := map[string]bool{}
			for i := 0; i < 60; i++ {
				p := (&w9CalcGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				seen[p.Question] = true
			}
			if len(seen) < 5 {
				t.Fatalf("%s kp%d: only %d distinct questions in 60 draws", kind, kp, len(seen))
			}
		}
	}
}
