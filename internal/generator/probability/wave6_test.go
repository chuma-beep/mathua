package probability

import (
	"math"
	"strconv"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func w6TestInt(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer\nfacts: %v", p.Question, key, p.Facts[key], p.Facts)
	}
	return v
}

func w6TestExt4(p0, p1, p2 int) int {
	if p2 == 0 {
		if p1 == 8 {
			return 0
		}
		return 4 * p0 / (8 - p1)
	}
	a := float64(p2)
	b := float64(p1 - 8)
	c := float64(p0)
	disc := b*b - 4*a*c
	if disc < 0 {
		disc = 0
	}
	r1 := (-b - math.Sqrt(disc)) / (2 * a)
	r2 := (-b + math.Sqrt(disc)) / (2 * a)
	lo := math.Min(r1, r2)
	return int(math.Round(4 * lo))
}

func TestWave6ProbabilityAnswerRecomputation(t *testing.T) {
	kinds := []string{"branching", "brownian", "chernoff"}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			for draw := 0; draw < 30; draw++ {
				p := (&w6ProbGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				want := w6TestProbAnswer(t, kind, kp, p)
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("%s kp%d: answer = %s, want %d\nquestion: %s\nfacts: %v", kind, kp, p.Answer, want, p.Question, p.Facts)
				}
			}
		}
	}
}

func w6TestProbAnswer(t *testing.T, kind string, kp int, p generator.Problem) int {
	t.Helper()
	switch kind {
	case "branching":
		p0 := w6TestInt(t, p, "p0")
		p1 := w6TestInt(t, p, "p1")
		p2 := w6TestInt(t, p, "p2")
		if p0+p1+p2 != 8 {
			t.Fatalf("offspring probabilities do not sum to 8: %d+%d+%d", p0, p1, p2)
		}
		if kp == 1 {
			return p1 + 2*p2
		}
		if kp == 2 {
			return w6TestExt4(p0, p1, p2)
		}
		return 64*p0 + 8*p1*p0 + p2*p0*p0
	case "brownian":
		if kp == 1 {
			t0 := w6TestInt(t, p, "t")
			c := w6TestInt(t, p, "c")
			return c * c * t0
		}
		if kp == 2 {
			return w6TestInt(t, p, "s")
		}
		return w6TestInt(t, p, "b")
	case "chernoff":
		k := w6TestInt(t, p, "p_num")
		m := w6TestInt(t, p, "trials")
		if kp == 2 {
			return m * k
		}
		r := 1
		for i := 0; i < m; i++ {
			r *= 8 + k
		}
		return r
	}
	t.Fatalf("unknown kind %q", kind)
	return 0
}

func TestWave6ProbabilityDiversity(t *testing.T) {
	kinds := []string{"branching", "brownian", "chernoff"}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			seen := map[string]bool{}
			for i := 0; i < 60; i++ {
				p := (&w6ProbGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				seen[p.Question] = true
			}
			if len(seen) < 5 {
				t.Fatalf("%s kp%d: only %d distinct questions in 60 draws", kind, kp, len(seen))
			}
			if len(seen) < 10 {
				t.Logf("%s kp%d: %d distinct questions in 60 draws, below the 10 target", kind, kp, len(seen))
			}
		}
	}
}
