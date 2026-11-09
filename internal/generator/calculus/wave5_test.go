package calculus

import (
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func w5TestInt(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer\nfacts: %v", p.Question, key, p.Facts[key], p.Facts)
	}
	return v
}

func w5TestVec(t *testing.T, p generator.Problem, key string) []int {
	t.Helper()
	parts := strings.Split(p.Facts[key], ",")
	out := make([]int, len(parts))
	for i, s := range parts {
		v, err := strconv.Atoi(s)
		if err != nil {
			t.Fatalf("%s: fact %q component %q is not an integer", p.Question, key, s)
		}
		out[i] = v
	}
	return out
}

func w5TestChoose(n, k int) int {
	if k < 0 || k > n {
		return 0
	}
	r := 1
	for i := 1; i <= k; i++ {
		r = r * (n - k + i) / i
	}
	return r
}

func TestWave5CalculusAnswerRecomputation(t *testing.T) {
	kinds := []string{"orientation", "oriented_flux", "differential_form", "wedge_product", "exterior_derivative"}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			for draw := 0; draw < 40; draw++ {
				p := (&w5CalcGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				want := w5TestCalculusAnswer(t, kind, kp, p)
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("%s kp%d: answer = %s, want %d\nquestion: %s\nfacts: %v", kind, kp, p.Answer, want, p.Question, p.Facts)
				}
			}
		}
	}
}

func w5TestCalculusAnswer(t *testing.T, kind string, kp int, p generator.Problem) int {
	t.Helper()
	switch kind {
	case "orientation":
		a := w5TestInt(t, p, "a")
		if kp == 1 {
			return -a
		}
		if kp == 2 {
			return -1
		}
		c := w5TestInt(t, p, "c")
		if c > 0 {
			return 1
		}
		return 0
	case "oriented_flux":
		a, b := w5TestInt(t, p, "a"), w5TestInt(t, p, "b")
		pv, q, rv := w5TestInt(t, p, "p"), w5TestInt(t, p, "q"), w5TestInt(t, p, "r")
		if kp == 1 {
			return -a*pv - b*q + rv
		}
		if kp == 2 {
			return a*pv + b*q - rv
		}
		return -q
	case "differential_form":
		a, b, c := w5TestInt(t, p, "a"), w5TestInt(t, p, "b"), w5TestInt(t, p, "c")
		if kp == 1 {
			return a*w5TestInt(t, p, "p") + b*w5TestInt(t, p, "q") + c*w5TestInt(t, p, "r")
		}
		if kp == 2 {
			u, v := w5TestVec(t, p, "u"), w5TestVec(t, p, "v")
			return a*(u[1]*v[2]-u[2]*v[1]) + b*(u[2]*v[0]-u[0]*v[2]) + c*(u[0]*v[1]-u[1]*v[0])
		}
		return w5TestChoose(w5TestInt(t, p, "n"), w5TestInt(t, p, "k"))
	case "wedge_product":
		a, b, c, d := w5TestInt(t, p, "a"), w5TestInt(t, p, "b"), w5TestInt(t, p, "c"), w5TestInt(t, p, "d")
		det := a*d - b*c
		if kp == 1 {
			return det
		}
		if kp == 2 {
			return -det
		}
		return b*w5TestInt(t, p, "r") - c*w5TestInt(t, p, "q")
	case "exterior_derivative":
		a, b, c := w5TestInt(t, p, "a"), w5TestInt(t, p, "b"), w5TestInt(t, p, "c")
		if kp == 1 {
			return a
		}
		if kp == 2 {
			return b - a
		}
		return a + b + c
	}
	t.Fatalf("unknown kind %q", kind)
	return 0
}

func TestWave5CalculusQuestionDiversity(t *testing.T) {
	kinds := []string{"orientation", "oriented_flux", "differential_form", "wedge_product", "exterior_derivative"}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			seen := map[string]bool{}
			for i := 0; i < 60; i++ {
				p := (&w5CalcGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				seen[p.Question] = true
			}
			if len(seen) < 5 {
				t.Fatalf("%s kp%d: only %d distinct questions in 60 draws, hard floor is 5", kind, kp, len(seen))
			}
			if len(seen) < 10 {
				t.Logf("%s kp%d: only %d distinct questions in 60 draws, below the Wave 5 target of 10", kind, kp, len(seen))
			}
		}
	}
}
