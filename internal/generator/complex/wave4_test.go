package complex

import (
	"strconv"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func w4TestInt(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer\nfacts: %v", p.Question, key, p.Facts[key], p.Facts)
	}
	return v
}

func w4TestPhi(n int) int {
	c := 0
	for k := 1; k <= n; k++ {
		if w4TestGCD(k, n) == 1 {
			c++
		}
	}
	return c
}

func w4TestGCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func w4TestFact(n int) int {
	r := 1
	for i := 2; i <= n; i++ {
		r *= i
	}
	return r
}

func TestWave4ComplexAnswerRecomputation(t *testing.T) {
	kinds := []string{
		"roots_unity", "branches", "complex_powers", "trig_zeros", "harmonic_conjugate",
		"poisson_kernel", "rouche", "open_mapping", "casorati", "picard",
		"analytic_continuation", "monodromy", "riemann_surface", "elliptic", "modular_forms",
		"gamma", "weierstrass_factorization", "mittag_leffler", "schwarz_christoffel", "jensen",
	}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			for draw := 0; draw < 40; draw++ {
				p := (&w4ComplexGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				want := w4TestComplexAnswer(t, kind, kp, p)
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("%s kp%d: answer = %s, want %d\nquestion: %s\nfacts: %v", kind, kp, p.Answer, want, p.Question, p.Facts)
				}
			}
		}
	}
}

func w4TestComplexAnswer(t *testing.T, kind string, kp int, p generator.Problem) int {
	t.Helper()
	switch kind {
	case "roots_unity":
		n := w4TestInt(t, p, "n")
		if kp == 1 {
			return n
		}
		if kp == 2 {
			return 0
		}
		return w4TestPhi(n)
	case "branches":
		if kp == 1 {
			return 180
		}
		if kp == 2 {
			return 1
		}
		return 0
	case "complex_powers":
		n := w4TestInt(t, p, "n")
		re, im := 0, 0
		for k := 0; k <= n; k++ {
			c := 1
			for i := 1; i <= k; i++ {
				c = c * (n - i + 1) / i
			}
			switch k % 4 {
			case 0:
				re += c
			case 1:
				im += c
			case 2:
				re -= c
			case 3:
				im -= c
			}
		}
		if kp == 1 {
			return re
		}
		if kp == 2 {
			return im
		}
		return 1 << n
	case "trig_zeros":
		k := w4TestInt(t, p, "k")
		if kp == 1 {
			return 2*k + 1
		}
		if kp == 2 {
			return 2
		}
		return 2 * k
	case "harmonic_conjugate":
		a := w4TestInt(t, p, "a")
		if kp == 1 {
			return 2 * a
		}
		if kp == 2 {
			return 1
		}
		return 4 * a
	case "poisson_kernel":
		num, den, scale := w4TestInt(t, p, "r_num"), w4TestInt(t, p, "r_den"), w4TestInt(t, p, "scale")
		if kp == 1 {
			return (den + num) / (den - num)
		}
		return scale
	case "rouche":
		if kp == 1 {
			return 1
		}
		return w4TestInt(t, p, "m")
	case "open_mapping":
		if kp == 2 {
			return 0
		}
		return 1
	case "casorati":
		if kp == 1 {
			return w4TestInt(t, p, "a")
		}
		return 1
	case "picard":
		if kp == 3 {
			return 0
		}
		return 1
	case "analytic_continuation":
		a := w4TestInt(t, p, "a")
		if kp == 1 {
			return a
		}
		if kp == 2 {
			return -a
		}
		return 1
	case "monodromy":
		m, n := w4TestInt(t, p, "m"), w4TestInt(t, p, "n")
		if kp == 1 {
			return m
		}
		if kp == 2 {
			return n
		}
		return 0
	case "riemann_surface":
		n := w4TestInt(t, p, "n")
		if kp == 1 || kp == 3 {
			return n
		}
		return 0
	case "elliptic":
		if kp == 1 || kp == 3 {
			return 2
		}
		return 0
	case "modular_forms":
		k := w4TestInt(t, p, "k")
		if kp == 1 {
			return 4 * k
		}
		if kp == 2 {
			return 6 * k
		}
		return 12 * k
	case "gamma":
		n, k := w4TestInt(t, p, "n"), w4TestInt(t, p, "k")
		if kp == 1 {
			return w4TestFact(n - 1)
		}
		if kp == 2 {
			return k + 1
		}
		return 1
	case "weierstrass_factorization":
		a := w4TestInt(t, p, "a")
		if kp == 1 || kp == 3 {
			return 1
		}
		return 2*a - 1
	case "mittag_leffler":
		k := w4TestInt(t, p, "k")
		if kp == 1 {
			return 2*k + 1
		}
		if kp == 2 {
			if k%2 == 1 {
				return -1
			}
			return 1
		}
		return 4*k + 1
	case "schwarz_christoffel":
		if kp == 1 {
			return w4TestInt(t, p, "n") - 2
		}
		if kp == 2 {
			return w4TestInt(t, p, "m")
		}
		return -2
	case "jensen":
		a := w4TestInt(t, p, "a")
		if kp == 1 {
			return a
		}
		if kp == 2 {
			return 1
		}
		return a + 1
	}
	t.Fatalf("unknown kind %q", kind)
	return 0
}

func TestWave4ComplexQuestionDiversity(t *testing.T) {
	kinds := []string{
		"roots_unity", "branches", "complex_powers", "trig_zeros", "harmonic_conjugate",
		"poisson_kernel", "rouche", "open_mapping", "casorati", "picard",
		"analytic_continuation", "monodromy", "riemann_surface", "elliptic", "modular_forms",
		"gamma", "weierstrass_factorization", "mittag_leffler", "schwarz_christoffel", "jensen",
	}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			seen := map[string]bool{}
			for i := 0; i < 60; i++ {
				p := (&w4ComplexGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				seen[p.Question] = true
			}
			if len(seen) < 5 {
				t.Fatalf("%s kp%d: only %d distinct questions in 60 draws", kind, kp, len(seen))
			}
		}
	}
}
