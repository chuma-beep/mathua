package topology

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

func w7TestBool(t *testing.T, p generator.Problem, key string) bool {
	t.Helper()
	v, err := strconv.ParseBool(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not a bool", p.Question, key, p.Facts[key])
	}
	return v
}

func w7TestPhi(n int) int {
	c := 0
	for k := 1; k <= n; k++ {
		if w7TestGCD(k, n) == 1 {
			c++
		}
	}
	return c
}

func w7TestGCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func TestWave7TopologyAnswerRecomputation(t *testing.T) {
	kinds := []string{"nets", "sequential_compact", "deck", "cup_product"}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			for draw := 0; draw < 40; draw++ {
				p := (&w7TopologyGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				want := w7TestTopologyAnswer(t, kind, kp, p)
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("%s kp%d: answer = %s, want %d\nquestion: %s\nfacts: %v", kind, kp, p.Answer, want, p.Question, p.Facts)
				}
			}
		}
	}
}

func w7TestTopologyAnswer(t *testing.T, kind string, kp int, p generator.Problem) int {
	t.Helper()
	switch kind {
	case "nets":
		if kp == 1 {
			return w7TestInt(t, p, "n")
		}
		if kp == 2 {
			return 2
		}
		return w7TestInt(t, p, "c")
	case "sequential_compact":
		if kp == 1 {
			return w7TestInt(t, p, "m")
		}
		if kp == 2 {
			return 3
		}
		if w7TestBool(t, p, "closed") {
			return 1
		}
		return 0
	case "deck":
		if kp == 1 {
			return w7TestInt(t, p, "n")
		}
		if kp == 2 {
			return w7TestPhi(w7TestInt(t, p, "n"))
		}
		return w7TestInt(t, p, "m")
	case "cup_product":
		if kp == 1 {
			return w7TestInt(t, p, "p") + w7TestInt(t, p, "q")
		}
		if kp == 2 {
			if (w7TestInt(t, p, "p")*w7TestInt(t, p, "q"))%2 == 0 {
				return 1
			}
			return -1
		}
		return w7TestInt(t, p, "n") + 1
	}
	t.Fatalf("unknown kind %q", kind)
	return 0
}

func TestWave7TopologyDiversity(t *testing.T) {
	kinds := []string{"nets", "sequential_compact", "deck", "cup_product"}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			seen := map[string]bool{}
			for i := 0; i < 60; i++ {
				p := (&w7TopologyGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
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
