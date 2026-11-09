package odes

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

func w4TestFact(n, k int) int {
	r := 1
	for i := 2; i <= n; i++ {
		r *= i
	}
	return r
}

func TestWave4ODEAnswerRecomputation(t *testing.T) {
	kinds := []string{
		"homogeneous_substitution", "clairaut", "riccati", "inverse_laplace", "laplace_derivative",
		"periodic", "initial_final", "frobenius", "bessel", "legendre", "multistep", "rkf",
		"implicit_euler", "shooting", "finite_difference", "fredholm", "delay", "perturbation",
		"duhamel", "images",
	}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			for draw := 0; draw < 40; draw++ {
				p := (&w4ODEGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				want := w4TestODEAnswer(t, kind, kp, p)
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("%s kp%d: answer = %s, want %d\nquestion: %s\nfacts: %v", kind, kp, p.Answer, want, p.Question, p.Facts)
				}
			}
		}
	}
}

func w4TestODEAnswer(t *testing.T, kind string, kp int, p generator.Problem) int {
	t.Helper()
	switch kind {
	case "homogeneous_substitution":
		a, b, y1 := w4TestInt(t, p, "a"), w4TestInt(t, p, "b"), w4TestInt(t, p, "y1")
		if kp == 1 {
			return b
		}
		if kp == 2 {
			return a
		}
		return 2 * (2*y1 + a)
	case "clairaut":
		c, x0 := w4TestInt(t, p, "c"), w4TestInt(t, p, "x0")
		if kp == 1 {
			return c*x0 + c*c
		}
		if kp == 2 {
			return -(x0 * x0) / 4
		}
		return 0
	case "riccati":
		a := w4TestInt(t, p, "a")
		if kp == 1 {
			return 1
		}
		if kp == 2 {
			return 2 * a
		}
		return -a
	case "inverse_laplace":
		a, b := w4TestInt(t, p, "a"), w4TestInt(t, p, "b")
		A, B := w4TestInt(t, p, "A"), w4TestInt(t, p, "B")
		if kp == 1 {
			return 1
		}
		if kp == 2 {
			return A + B
		}
		_ = a
		_ = b
		return 1
	case "laplace_derivative":
		a, y0 := w4TestInt(t, p, "a"), w4TestInt(t, p, "y0")
		if kp == 1 || kp == 2 {
			return -y0
		}
		return y0 / a
	case "periodic":
		T, duty, high := w4TestInt(t, p, "period"), w4TestInt(t, p, "duty"), w4TestInt(t, p, "high")
		if kp == 1 {
			return 3 * T
		}
		if kp == 2 {
			if w4TestInt(t, p, "t") < duty {
				return high
			}
			return 0
		}
		return high * duty
	case "initial_final":
		a, b, k := w4TestInt(t, p, "a"), w4TestInt(t, p, "b"), w4TestInt(t, p, "k")
		if kp == 1 {
			return a
		}
		if kp == 2 {
			return a * k / b
		}
		if b > 0 {
			return 1
		}
		return 0
	case "frobenius":
		p0, q0 := w4TestInt(t, p, "p0"), w4TestInt(t, p, "q0")
		if kp == 1 {
			return 1 - p0
		}
		if kp == 2 {
			return q0
		}
		return (p0-1)*(p0-1) - 4*q0
	case "bessel":
		nu, k := w4TestInt(t, p, "nu"), w4TestInt(t, p, "k")
		if kp == 1 {
			return nu
		}
		if kp == 2 {
			return 2 * k * (nu + k)
		}
		if nu == 0 {
			return 1
		}
		return 0
	case "legendre":
		n, x := w4TestInt(t, p, "n"), w4TestInt(t, p, "x")
		if kp == 1 {
			if x == 1 || n%2 == 0 {
				return 1
			}
			return -1
		}
		if kp == 2 {
			return 2
		}
		return 2*n + 1
	case "multistep":
		y, fn, fp := w4TestInt(t, p, "y"), w4TestInt(t, p, "fn"), w4TestInt(t, p, "fprev")
		if kp == 1 {
			return y + 3*fn - fp
		}
		if kp == 2 {
			return 2
		}
		r := w4TestInt(t, p, "root")
		if r >= -1 && r <= 1 {
			return 1
		}
		return 0
	case "rkf":
		low, high := w4TestInt(t, p, "low"), w4TestInt(t, p, "high")
		if kp == 1 {
			return low + high
		}
		if kp == 2 {
			return high - low
		}
		return 5
	case "implicit_euler":
		h, lambda, y0 := w4TestInt(t, p, "h"), w4TestInt(t, p, "lambda"), w4TestInt(t, p, "y0")
		if kp == 1 || kp == 2 {
			return y0 / (1 + h*lambda)
		}
		if h*lambda > 0 {
			return 1
		}
		return 0
	case "shooting":
		a, b, m := w4TestInt(t, p, "a"), w4TestInt(t, p, "b"), w4TestInt(t, p, "m")
		if kp == 1 {
			return b - a
		}
		if kp == 2 {
			return a + m - b
		}
		return 1
	case "finite_difference":
		n, k := w4TestInt(t, p, "n"), w4TestInt(t, p, "k")
		if kp == 1 {
			return n - 1
		}
		if kp == 2 {
			return n * k
		}
		return w4TestInt(t, p, "i") * k
	case "fredholm":
		a, b, lambda := w4TestInt(t, p, "a"), w4TestInt(t, p, "b"), w4TestInt(t, p, "lambda")
		if kp == 1 {
			return a
		}
		if kp == 2 {
			return a * b
		}
		return 1 - lambda*a*b
	case "delay":
		a, tau, phi := w4TestInt(t, p, "a"), w4TestInt(t, p, "tau"), w4TestInt(t, p, "phi")
		if kp == 1 {
			return phi
		}
		if kp == 2 {
			return 1 - a*tau
		}
		return a
	case "perturbation":
		c := w4TestInt(t, p, "c")
		if kp == 1 {
			return c
		}
		if kp == 2 {
			return c * c
		}
		return 1
	case "duhamel":
		if kp == 1 || kp == 2 {
			return 1
		}
		return 0
	case "images":
		y0 := w4TestInt(t, p, "y0")
		if kp == 1 {
			return -y0
		}
		if kp == 2 {
			return 0
		}
		return 2 * y0
	}
	t.Fatalf("unknown kind %q", kind)
	return 0
}

func TestWave4ODEQuestionDiversity(t *testing.T) {
	kinds := []string{
		"homogeneous_substitution", "clairaut", "riccati", "inverse_laplace", "laplace_derivative",
		"periodic", "initial_final", "frobenius", "bessel", "legendre", "multistep", "rkf",
		"implicit_euler", "shooting", "finite_difference", "fredholm", "delay", "perturbation",
		"duhamel", "images",
	}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			seen := map[string]bool{}
			for i := 0; i < 60; i++ {
				p := (&w4ODEGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				seen[p.Question] = true
			}
			if len(seen) < 5 {
				t.Fatalf("%s kp%d: only %d distinct questions in 60 draws", kind, kp, len(seen))
			}
		}
	}
}
