package complex

import (
	"math"
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
				if p.KPOrdinal != kp {
					t.Fatalf("%s kp%d: problem records ordinal %d", kind, kp, p.KPOrdinal)
				}
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
			// The sum is 0; the product form publishes the product it asks for.
			if prod, ok := p.Facts["product"]; ok {
				v, err := strconv.Atoi(prod)
				if err != nil {
					t.Fatalf("product fact %q is not an integer", prod)
				}
				return v
			}
			return 0
		}
		return w4TestPhi(n)
	case "branches":
		switch kp {
		case 1:
			re := w4TestInt(t, p, "re")
			im := w4TestInt(t, p, "im")
			return int(math.Round(math.Atan2(float64(im), float64(re)) * 180 / math.Pi))
		case 2:
			// A circle meets the cut once; the point form publishes its verdict.
			if v, ok := p.Facts["on_cut"]; ok {
				answer, err := strconv.Atoi(v)
				if err != nil {
					t.Fatalf("on_cut fact %q is not an integer", v)
				}
				return answer
			}
			return 1
		default:
			return 0
		}
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
			return k
		}
		return 2 * k
	case "harmonic_conjugate":
		a := w4TestInt(t, p, "a")
		if kp == 1 {
			return 2 * a
		}
		if kp == 2 {
			return w4TestInt(t, p, "harmonic")
		}
		x0, y0 := w4TestInt(t, p, "x0"), w4TestInt(t, p, "y0")
		return 2 * a * x0 * y0
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
		if kp == 1 {
			return 1
		}
		if kp == 2 {
			return 0
		}
		return w4TestInt(t, p, "applies")
	case "casorati":
		if kp == 1 {
			return w4TestInt(t, p, "coeff")
		}
		if kp == 2 {
			return w4TestInt(t, p, "omit_zero")
		}
		return 1
	case "picard":
		if kp == 1 {
			return 1
		}
		if kp == 2 {
			return w4TestInt(t, p, "omitted")
		}
		return 0
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
		return w4TestInt(t, p, "a")
	case "elliptic":
		if kp == 1 {
			return 2
		}
		if kp == 2 {
			return 0
		}
		return w4TestInt(t, p, "poles")
	case "modular_forms":
		if kp == 1 {
			return 4 * w4TestInt(t, p, "k")
		}
		if kp == 2 {
			return 4*w4TestInt(t, p, "p") + 6*w4TestInt(t, p, "q")
		}
		return 12*w4TestInt(t, p, "k") + 4*w4TestInt(t, p, "m")
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
		if kp == 1 {
			return w4TestInt(t, p, "m")
		}
		if kp == 2 {
			return 2*w4TestInt(t, p, "k")*a - 1
		}
		return 1
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

func TestModulusAndLogarithmRecomputeTheirAnswers(t *testing.T) {
	modulus := &modulusGen{}
	for kp := 1; kp <= 3; kp++ {
		for i := 0; i < 200; i++ {
			p := modulus.Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			if p.KPOrdinal != kp {
				t.Fatalf("modulus kp%d: problem records ordinal %d", kp, p.KPOrdinal)
			}
			want := 0
			switch kp {
			case 1:
				a, b := w4TestInt(t, p, "a"), w4TestInt(t, p, "b")
				want = int(math.Round(math.Sqrt(float64(a*a + b*b))))
				if want*want != a*a+b*b {
					t.Fatalf("modulus kp1: %d+%di has no integer modulus", a, b)
				}
			case 2:
				a, b := w4TestInt(t, p, "a"), w4TestInt(t, p, "b")
				want = a*a + b*b
			default:
				a, b := w4TestInt(t, p, "a"), w4TestInt(t, p, "b")
				c, d := w4TestInt(t, p, "c"), w4TestInt(t, p, "d")
				// Equality holds exactly when z2 is a positive multiple of z1.
				if a*d-b*c == 0 && a*c+b*d > 0 {
					want = 1
				}
			}
			if p.Answer != strconv.Itoa(want) {
				t.Fatalf("modulus kp%d: answer = %s, want %d\nquestion: %s", kp, p.Answer, want, p.Question)
			}
		}
	}

	ln := &logarithmGen{}
	for kp := 1; kp <= 3; kp++ {
		for i := 0; i < 200; i++ {
			p := ln.Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			if p.KPOrdinal != kp {
				t.Fatalf("logarithm kp%d: problem records ordinal %d", kp, p.KPOrdinal)
			}
			want := 0
			switch kp {
			case 1:
				z := w4TestInt(t, p, "z")
				if z < 0 {
					want = 1
				}
			case 2:
				re, im := w4TestInt(t, p, "re"), w4TestInt(t, p, "im")
				ray := w4TestInt(t, p, "ray")
				deg := math.Atan2(float64(im), float64(re)) * 180 / math.Pi
				deg = math.Mod(deg+360, 360)
				if int(math.Round(deg))%360 == ray {
					want = 1
				}
			default:
				sum := w4TestInt(t, p, "sum")
				m := 0
				switch {
				case sum > 180:
					m = 1
				case sum <= -180:
					m = -1
				}
				want = -m
			}
			if p.Answer != strconv.Itoa(want) {
				t.Fatalf("logarithm kp%d: answer = %s, want %d\nquestion: %s\nfacts: %v", kp, p.Answer, want, p.Question, p.Facts)
			}
		}
	}
}
