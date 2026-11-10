package calculus

import (
	"strconv"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func w9rTestInt(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer\nfacts: %v", p.Question, key, p.Facts[key], p.Facts)
	}
	return v
}

func w9rTestFactorial(n int) int {
	r := 1
	for i := 2; i <= n; i++ {
		r *= i
	}
	return r
}

func TestWave9RemainingCalculusAnswers(t *testing.T) {
	kinds := []string{
		"continuity_composition", "continuity_inverse", "continuity_extension",
		"deriv_absolute", "deriv_floor", "deriv_sign", "shadow", "angle",
		"optimization_distance", "optimization_cylinder", "integral_floor",
		"abel_theorem", "asymptotic_series", "euler_maclaurin", "borel", "taubian", "lambert",
	}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			for draw := 0; draw < 40; draw++ {
				p := (&w9rCalcGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				want := w9rTestCalculusAnswer(t, kind, kp, p)
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("%s kp%d: answer = %s, want %d\nquestion: %s\nfacts: %v", kind, kp, p.Answer, want, p.Question, p.Facts)
				}
			}
		}
	}
}

func w9rTestCalculusAnswer(t *testing.T, kind string, kp int, p generator.Problem) int {
	t.Helper()
	switch kind {
	case "continuity_composition":
		a, x0 := w9rTestInt(t, p, "a"), w9rTestInt(t, p, "x0")
		if kp == 1 {
			return (x0 + a) * (x0 + a)
		}
		if kp == 2 {
			return x0*x0 + a
		}
		return x0 + a
	case "continuity_inverse":
		r := w9rTestInt(t, p, "r")
		if kp == 1 {
			return r
		}
		if kp == 2 {
			return r + 1
		}
		return 3 * r * r
	case "continuity_extension":
		a := w9rTestInt(t, p, "a")
		if kp == 1 || kp == 2 {
			return 2 * a
		}
		return 2*a + 1
	case "deriv_absolute":
		if kp == 1 {
			return -1
		}
		if kp == 2 {
			return 1
		}
		return 0
	case "deriv_floor":
		a := w9rTestInt(t, p, "a")
		if kp == 1 || kp == 3 {
			return 0
		}
		return a
	case "deriv_sign":
		if kp == 1 {
			return 0
		}
		return 1
	case "shadow":
		h, H := w9rTestInt(t, p, "h"), w9rTestInt(t, p, "H")
		x, rate := w9rTestInt(t, p, "x"), w9rTestInt(t, p, "rate")
		if kp == 1 {
			return h * rate / (H - h)
		}
		if kp == 2 {
			return h * x / (H - h)
		}
		return 3 * x / 2
	case "angle":
		h, x, rate := w9rTestInt(t, p, "h"), w9rTestInt(t, p, "x"), w9rTestInt(t, p, "rate")
		if kp == 1 {
			return x + h
		}
		if kp == 2 {
			return x*x + h*h
		}
		return h * rate
	case "optimization_distance":
		a := w9rTestInt(t, p, "a")
		if kp == 1 {
			return a / 2
		}
		if kp == 2 {
			return a * a / 2
		}
		return a
	case "optimization_cylinder":
		r := w9rTestInt(t, p, "r")
		if kp == 1 {
			return r
		}
		if kp == 2 {
			return 2 * r
		}
		return 2 * r * r * r
	case "integral_floor":
		n := w9rTestInt(t, p, "n")
		if kp == 1 || kp == 2 {
			return n * (n - 1) / 2
		}
		return n - 1
	case "abel_theorem":
		if kp == 1 {
			return w9rTestInt(t, p, "k") - 1
		}
		if kp == 2 {
			return w9rTestInt(t, p, "k")
		}
		return 2
	case "asymptotic_series":
		n, c := w9rTestInt(t, p, "n"), w9rTestInt(t, p, "c")
		if kp == 1 {
			return c
		}
		return n + 1
	case "euler_maclaurin":
		n := w9rTestInt(t, p, "n")
		if kp == 1 || kp == 2 {
			return n + 1
		}
		return n*n - 1
	case "borel":
		n := w9rTestInt(t, p, "n")
		if kp == 1 {
			return w9rTestFactorial(n)
		}
		if kp == 2 {
			return w9rTestInt(t, p, "k") - 1
		}
		return 0
	case "taubian":
		c, pv := w9rTestInt(t, p, "c"), w9rTestInt(t, p, "p")
		if kp == 1 {
			return c
		}
		if kp == 2 {
			return 0
		}
		if pv == 2 {
			return 1
		}
		return 0
	case "lambert":
		n := w9rTestInt(t, p, "n")
		tau, sigma, sigma2 := 0, 0, 0
		for d := 1; d <= n; d++ {
			if n%d == 0 {
				tau++
				sigma += d
				sigma2 += d * d
			}
		}
		if kp == 1 {
			return tau
		}
		if kp == 2 {
			return sigma
		}
		return sigma2
	}
	t.Fatalf("unknown kind %q", kind)
	return 0
}

func TestWave9RemainingCalculusDiversity(t *testing.T) {
	kinds := []string{
		"continuity_composition", "continuity_inverse", "continuity_extension",
		"deriv_absolute", "deriv_floor", "deriv_sign", "shadow", "angle",
		"optimization_distance", "optimization_cylinder", "integral_floor",
		"abel_theorem", "asymptotic_series", "euler_maclaurin", "borel", "taubian", "lambert",
	}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			seen := map[string]bool{}
			for i := 0; i < 60; i++ {
				p := (&w9rCalcGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				seen[p.Question] = true
			}
			if len(seen) < 5 {
				t.Fatalf("%s kp%d: only %d distinct questions in 60 draws", kind, kp, len(seen))
			}
			if len(seen) < 10 {
				t.Logf("%s kp%d: %d distinct questions in 60 draws (below 10 target)", kind, kp, len(seen))
			}
		}
	}
}
