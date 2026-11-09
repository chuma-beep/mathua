package topology

import (
	"strconv"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func w7RemTestInt(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer\nfacts: %v", p.Question, key, p.Facts[key], p.Facts)
	}
	return v
}

func w7RemTestBool(t *testing.T, p generator.Problem, key string) bool {
	t.Helper()
	v, err := strconv.ParseBool(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not a bool", p.Question, key, p.Facts[key])
	}
	return v
}

func w7RemTestPow(a, e int) int {
	r := 1
	for i := 0; i < e; i++ {
		r *= a
	}
	return r
}

func w7RemTestBinom(n, k int) int {
	if k < 0 || k > n {
		return 0
	}
	if k > n-k {
		k = n - k
	}
	r := 1
	for i := 1; i <= k; i++ {
		r = r * (n - k + i) / i
	}
	return r
}

func TestWave7RemainingTopologyAnswerRecomputation(t *testing.T) {
	kinds := []string{"filters", "paracompact", "partition_unity", "function_space", "uniform_space", "limit_point_compact", "continuum", "homotopy_group", "fibration", "cofibration", "simplicial_homology", "singular_homology", "poincare_duality", "jones_polynomial", "seifert_surface", "brouwer"}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			for draw := 0; draw < 40; draw++ {
				p := (&w7TopologyRemainingGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				want := w7RemTestTopologyAnswer(t, kind, kp, p)
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("%s kp%d: answer = %s, want %d\nquestion: %s\nfacts: %v", kind, kp, p.Answer, want, p.Question, p.Facts)
				}
			}
		}
	}
}

func w7RemTestTopologyAnswer(t *testing.T, kind string, kp int, p generator.Problem) int {
	t.Helper()
	switch kind {
	case "filters":
		n, s := w7RemTestInt(t, p, "n"), w7RemTestInt(t, p, "s")
		if kp == 1 {
			return w7RemTestPow(2, n-s)
		}
		if kp == 2 {
			return n
		}
		return 1
	case "paracompact":
		n := w7RemTestInt(t, p, "n")
		if kp == 1 {
			return n
		}
		if kp == 2 {
			return n + w7RemTestInt(t, p, "extra")
		}
		return 1
	case "partition_unity":
		m := w7RemTestInt(t, p, "m")
		if kp == 1 {
			return 1
		}
		if kp == 2 {
			return m
		}
		return 1
	case "function_space":
		a, b := w7RemTestInt(t, p, "a"), w7RemTestInt(t, p, "b")
		if kp == 1 {
			return w7RemTestPow(b, a)
		}
		if kp == 2 {
			return b
		}
		return a
	case "uniform_space":
		n := w7RemTestInt(t, p, "n")
		if kp == 1 {
			return n * n
		}
		if kp == 2 {
			return n
		}
		return n
	case "limit_point_compact":
		if kp == 1 {
			return 0
		}
		if kp == 2 {
			if w7RemTestBool(t, p, "closed") {
				return 1
			}
			return 0
		}
		return 1
	case "continuum":
		n := w7RemTestInt(t, p, "n")
		if kp == 1 || kp == 2 {
			return 1
		}
		return n - 1
	case "homotopy_group":
		if kp == 1 {
			return 1
		}
		if kp == 2 {
			return 0
		}
		return 1
	case "fibration":
		fb, base := w7RemTestInt(t, p, "fiber"), w7RemTestInt(t, p, "base")
		if kp == 1 {
			return fb
		}
		if kp == 2 {
			return fb + base
		}
		return base
	case "cofibration":
		n := w7RemTestInt(t, p, "n")
		if kp == 1 {
			return n + 1
		}
		return 1
	case "simplicial_homology":
		n, k := w7RemTestInt(t, p, "n"), w7RemTestInt(t, p, "k")
		if kp == 1 {
			return w7RemTestBinom(n+1, k+1)
		}
		if kp == 2 {
			return 1
		}
		return n + 1
	case "singular_homology":
		m := w7RemTestInt(t, p, "m")
		if kp == 1 || kp == 2 {
			return m
		}
		return 1 - m
	case "poincare_duality":
		n, k := w7RemTestInt(t, p, "n"), w7RemTestInt(t, p, "k")
		if kp == 1 {
			if k == 0 || k == n {
				return 1
			}
			return 0
		}
		if kp == 2 {
			return n % 2
		}
		return n
	case "jones_polynomial":
		c := w7RemTestInt(t, p, "c")
		if kp == 1 {
			return w7RemTestPow(-2, c-1)
		}
		if kp == 2 {
			return c - 1
		}
		return 1
	case "seifert_surface":
		g, c := w7RemTestInt(t, p, "g"), w7RemTestInt(t, p, "c")
		if kp == 1 {
			return 2 - 2*g - c
		}
		if kp == 2 {
			return 2*g + c
		}
		return c
	case "brouwer":
		n := w7RemTestInt(t, p, "n")
		if kp == 1 {
			return n
		}
		if kp == 2 {
			return 1
		}
		return 0
	}
	t.Fatalf("unknown kind %q", kind)
	return 0
}

func TestWave7RemainingTopologyDiversity(t *testing.T) {
	kinds := []string{"filters", "paracompact", "partition_unity", "function_space", "uniform_space", "limit_point_compact", "continuum", "homotopy_group", "fibration", "cofibration", "simplicial_homology", "singular_homology", "poincare_duality", "jones_polynomial", "seifert_surface", "brouwer"}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			seen := map[string]bool{}
			for i := 0; i < 60; i++ {
				p := (&w7TopologyRemainingGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
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
