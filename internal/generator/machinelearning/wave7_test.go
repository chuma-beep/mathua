package machinelearning

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

func w7TestAbs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func TestWave7MLAnswerRecomputation(t *testing.T) {
	kinds := []string{"lasso", "gradient_boosting"}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			for draw := 0; draw < 40; draw++ {
				p := (&w7MLGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				want := w7TestMLAnswer(t, kind, kp, p)
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("%s kp%d: answer = %s, want %d\nquestion: %s\nfacts: %v", kind, kp, p.Answer, want, p.Question, p.Facts)
				}
			}
		}
	}
}

func w7TestMLAnswer(t *testing.T, kind string, kp int, p generator.Problem) int {
	t.Helper()
	switch kind {
	case "lasso":
		lambda := w7TestInt(t, p, "lambda")
		if kp == 1 {
			return lambda * (w7TestAbs(w7TestInt(t, p, "w1")) + w7TestAbs(w7TestInt(t, p, "w2")))
		}
		if kp == 2 {
			c := w7TestInt(t, p, "c")
			if w7TestAbs(c) <= lambda {
				return 0
			}
			if c > 0 {
				return c - lambda
			}
			return c + lambda
		}
		x, w, y := w7TestInt(t, p, "x"), w7TestInt(t, p, "w"), w7TestInt(t, p, "y")
		r := y - x*w
		return r*r + 2*lambda*w7TestAbs(w)
	case "gradient_boosting":
		if kp == 1 {
			return w7TestInt(t, p, "y") - w7TestInt(t, p, "pred")
		}
		if kp == 2 {
			y, pred, den := w7TestInt(t, p, "y"), w7TestInt(t, p, "pred"), w7TestInt(t, p, "den")
			return pred + (y-pred)/den
		}
		return w7TestInt(t, p, "base") + (w7TestInt(t, p, "r1")+w7TestInt(t, p, "r2"))/2
	}
	t.Fatalf("unknown kind %q", kind)
	return 0
}

func TestWave7MLDiversity(t *testing.T) {
	kinds := []string{"lasso", "gradient_boosting"}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			seen := map[string]bool{}
			for i := 0; i < 60; i++ {
				p := (&w7MLGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
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
