package realanalysis

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

func w5TestMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func w5TestMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func TestWave5RealAnalysisAnswerRecomputation(t *testing.T) {
	kinds := []string{"limsup", "liminf", "cluster", "lebesgue", "simple_integral"}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			for draw := 0; draw < 40; draw++ {
				p := (&w5RAGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				want := w5TestRealAnalysisAnswer(t, kind, kp, p)
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("%s kp%d: answer = %s, want %d\nquestion: %s\nfacts: %v", kind, kp, p.Answer, want, p.Question, p.Facts)
				}
			}
		}
	}
}

func w5TestRealAnalysisAnswer(t *testing.T, kind string, kp int, p generator.Problem) int {
	t.Helper()
	switch kind {
	case "limsup", "liminf":
		if kp == 2 {
			return w5TestInt(t, p, "limit")
		}
		a, b := w5TestInt(t, p, "even"), w5TestInt(t, p, "odd")
		if kind == "limsup" {
			if kp == 3 {
				return w5TestInt(t, p, "shift") + w5TestMax(a, b)
			}
			return w5TestMax(a, b)
		}
		if kp == 3 {
			return w5TestInt(t, p, "shift") + w5TestMin(a, b)
		}
		return w5TestMin(a, b)
	case "cluster":
		if kp == 1 {
			vals := strings.Split(p.Facts["values"], ",")
			return len(vals)
		}
		if kp == 2 {
			return w5TestMax(w5TestInt(t, p, "even"), w5TestInt(t, p, "odd"))
		}
		return w5TestInt(t, p, "center") - w5TestInt(t, p, "offset")
	case "lebesgue":
		a, b := w5TestInt(t, p, "a"), w5TestInt(t, p, "b")
		if kp == 1 {
			return b - a
		}
		c, d := w5TestInt(t, p, "c"), w5TestInt(t, p, "d")
		if kp == 2 {
			return (b - a) + (d - c)
		}
		return d - a
	case "simple_integral":
		if kp == 2 {
			return w5TestInt(t, p, "c1")*w5TestInt(t, p, "w1") + w5TestInt(t, p, "c2")*w5TestInt(t, p, "w2") + w5TestInt(t, p, "c3")*w5TestInt(t, p, "w3")
		}
		a, b := w5TestInt(t, p, "a"), w5TestInt(t, p, "b")
		c, d := w5TestInt(t, p, "c"), w5TestInt(t, p, "d")
		if kp == 1 {
			return w5TestInt(t, p, "c1")*(b-a) + w5TestInt(t, p, "c2")*(d-c)
		}
		return (b - a) + (d - c)
	}
	t.Fatalf("unknown kind %q", kind)
	return 0
}

func TestWave5RealAnalysisQuestionDiversity(t *testing.T) {
	kinds := []string{"limsup", "liminf", "cluster", "lebesgue", "simple_integral"}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			seen := map[string]bool{}
			for i := 0; i < 60; i++ {
				p := (&w5RAGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
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
