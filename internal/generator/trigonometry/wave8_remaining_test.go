package trigonometry

import (
	"strconv"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func w8RTInt(t *testing.T, p generator.Problem, k string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[k])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer: %v", p.Question, k, p.Facts[k], p.Facts)
	}
	return v
}

func w8RUGH(n int) int {
	r := 1
	for i := 2; i <= n; i++ {
		r *= i
	}
	return r
}

func w8RChebU(n, x int) int {
	if n == 0 {
		return 1
	}
	if n == 1 {
		return 2 * x
	}
	a, b := 1, 2*x
	for k := 2; k <= n; k++ {
		a, b = b, 2*x*b-a
	}
	return b
}

func TestWave8RemainingTrig(t *testing.T) {
	kinds := []string{"law_tangents", "mollweide", "half_side", "napier", "spherical_sines", "spherical_cosines", "spherical_excess", "haversine", "exsecant", "gudermannian", "weierstrass", "chebyshev_second", "ptolemy", "brahmagupta", "machin", "phasor"}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			for draw := 0; draw < 40; draw++ {
				p := (&w8TrigRemainingGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				want := w8RemTrigAnswer(t, kind, kp, p)
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("%s kp%d: answer %s want %d\nquestion: %s\nfacts: %v", kind, kp, p.Answer, want, p.Question, p.Facts)
				}
			}
		}
	}
}

func w8RemTrigAnswer(t *testing.T, kind string, kp int, p generator.Problem) int {
	t.Helper()
	switch kind {
	case "law_tangents":
		a, b := w8RTInt(t, p, "a"), w8RTInt(t, p, "b")
		if kp == 1 {
			return a - b
		}
		if kp == 2 {
			return a + b
		}
		if a > b {
			return 1
		}
		return 2
	case "mollweide":
		a, b, c := w8RTInt(t, p, "a"), w8RTInt(t, p, "b"), w8RTInt(t, p, "c")
		if kp == 1 {
			return a + b
		}
		if kp == 2 {
			return c
		}
		return a - b
	case "half_side":
		a, b, c := w8RTInt(t, p, "a"), w8RTInt(t, p, "b"), w8RTInt(t, p, "c")
		s := (a + b + c) / 2
		if kp == 1 {
			return s
		}
		if kp == 2 {
			return s - a
		}
		return (a + b - c) / 2
	case "napier":
		a, b, c := w8RTInt(t, p, "a"), w8RTInt(t, p, "b"), w8RTInt(t, p, "c")
		if kp == 1 {
			return a - b
		}
		if kp == 2 {
			return a + b
		}
		return a + b + c
	case "spherical_sines":
		if kp == 2 {
			return w8RTInt(t, p, "result_den")
		}
		return w8RTInt(t, p, "result_num")
	case "spherical_cosines":
		if kp == 1 {
			return w8RTInt(t, p, "cos_c_num")
		}
		if kp == 2 {
			return w8RTInt(t, p, "cos_c_den")
		}
		if w8RTInt(t, p, "cos_c_num") < 0 {
			return 1
		}
		return 0
	case "spherical_excess":
		if kp == 1 {
			return w8RTInt(t, p, "excess")
		}
		if kp == 2 {
			return w8RTInt(t, p, "sum")
		}
		return w8RTInt(t, p, "C")
	case "haversine":
		fh := w8RTInt(t, p, "four_hav")
		m := w8RTInt(t, p, "scale")
		if kp == 1 {
			return m * fh
		}
		if kp == 2 {
			return m * 2 * fh
		}
		if fh >= 2 {
			return 1
		}
		return 0
	case "exsecant":
		a, c := w8RTInt(t, p, "a"), w8RTInt(t, p, "c")
		if kp == 2 {
			return a
		}
		return c - a
	case "gudermannian":
		num, den := w8RTInt(t, p, "num"), w8RTInt(t, p, "den")
		if kp == 1 {
			return num
		}
		if kp == 2 {
			return den
		}
		if num > den {
			return 1
		}
		return 0
	case "weierstrass":
		tv := w8RTInt(t, p, "t")
		if kp == 1 {
			return 2 * tv
		}
		if kp == 2 {
			return 1 - tv*tv
		}
		return w8RTInt(t, p, "den")
	case "chebyshev_second":
		n, x := w8RTInt(t, p, "n"), w8RTInt(t, p, "x")
		if kp == 1 {
			return w8RChebU(n, x)
		}
		if kp == 2 {
			return 1 << n
		}
		ans := n + 1
		if n%2 == 1 {
			ans = -(n + 1)
		}
		return ans
	case "ptolemy":
		a, b := w8RTInt(t, p, "a"), w8RTInt(t, p, "b")
		if kp == 1 {
			return a*a + b*b
		}
		if kp == 2 {
			return a * b
		}
		return a*a + b*b - a*b
	case "brahmagupta":
		a, b := w8RTInt(t, p, "a"), w8RTInt(t, p, "b")
		if kp == 1 {
			return a * b
		}
		if kp == 2 {
			return a + b
		}
		return b
	case "machin":
		k, a, b := w8RTInt(t, p, "k"), w8RTInt(t, p, "a"), w8RTInt(t, p, "b")
		m := w8RTInt(t, p, "scale")
		if kp == 1 {
			return m * k
		}
		if kp == 2 {
			return m * (a + b)
		}
		return m * k * a
	case "phasor":
		if kp == 1 {
			return w8RTInt(t, p, "r")
		}
		if kp == 2 {
			return w8RTInt(t, p, "a")
		}
		return w8RTInt(t, p, "b")
	}
	t.Fatalf("unknown kind %s", kind)
	return 0
}

func TestWave8RemainingTrigDiversity(t *testing.T) {
	kinds := []string{"law_tangents", "mollweide", "half_side", "napier", "spherical_sines", "spherical_cosines", "spherical_excess", "haversine", "exsecant", "gudermannian", "weierstrass", "chebyshev_second", "ptolemy", "brahmagupta", "machin", "phasor"}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			seen := map[string]bool{}
			for i := 0; i < 60; i++ {
				p := (&w8TrigRemainingGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				seen[p.Question] = true
			}
			if len(seen) < 5 {
				t.Fatalf("%s kp%d only %d distinct questions", kind, kp, len(seen))
			}
			if len(seen) < 10 {
				t.Logf("below target: %s kp%d = %d distinct in 60", kind, kp, len(seen))
			}
		}
	}
}
