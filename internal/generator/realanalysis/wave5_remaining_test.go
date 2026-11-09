package realanalysis

import (
	"strconv"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func w5RTestAbs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func w5RTestMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func TestWave5RemainingRealAnswers(t *testing.T) {
	kinds := []string{
		"double_series", "fourier_convergence", "measurable_sets", "measurable_functions",
		"lebesgue_integral", "monotone_convergence", "fatou", "dominated_convergence",
		"product", "fubini_tonelli", "lp_space", "holder", "minkowski", "banach",
		"hilbert", "bounded_operator", "riesz_representation", "absolute_continuity",
		"lebesgue_differentiation",
	}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			for draw := 0; draw < 30; draw++ {
				p := (&w5RARemainingGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				want := w5TestRemainingReal(t, kind, kp, p)
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("%s kp%d: answer = %s, want %d\nquestion: %s\nfacts: %v", kind, kp, p.Answer, want, p.Question, p.Facts)
				}
			}
		}
	}
}

func w5TestRemainingReal(t *testing.T, kind string, kp int, p generator.Problem) int {
	t.Helper()
	switch kind {
	case "double_series":
		m, n := w5TestInt(t, p, "m"), w5TestInt(t, p, "n")
		if kp == 2 {
			return (m * (m + 1) / 2) * (n * (n + 1) / 2)
		}
		return n*m*(m+1)/2 + m*n*(n+1)/2
	case "fourier_convergence":
		if kp == 1 {
			return w5TestInt(t, p, "coefficient")
		}
		if kp == 2 {
			n := w5TestInt(t, p, "index")
			return 2*n + 1
		}
		return w5TestInt(t, p, "a0") + w5TestInt(t, p, "a1")
	case "measurable_sets":
		if kp == 1 {
			return 1 << w5TestInt(t, p, "atoms")
		}
		if kp == 2 {
			n := w5TestInt(t, p, "n")
			r := w5TestInt(t, p, "residue")
			c := 0
			for i := 1; i <= n; i++ {
				if i%3 == r {
					c++
				}
			}
			return c
		}
		return w5TestInt(t, p, "m1") + w5TestInt(t, p, "m2") + w5TestInt(t, p, "m3")
	case "measurable_functions":
		if kp == 1 {
			x := w5TestInt(t, p, "x")
			thr := w5TestInt(t, p, "threshold")
			if x >= thr {
				return 1
			}
			return 0
		}
		if kp == 2 {
			return w5TestInt(t, p, "scale")*w5TestInt(t, p, "x") + 1
		}
		x := w5TestInt(t, p, "x")
		a := w5TestInt(t, p, "a") * x
		b := w5TestInt(t, p, "b") * x
		return w5RTestMax(a, b)
	case "lebesgue_integral":
		if kp == 1 {
			return w5TestInt(t, p, "c1")*w5TestInt(t, p, "l1") + w5TestInt(t, p, "c2")*w5TestInt(t, p, "l2")
		}
		if kp == 2 {
			return w5TestInt(t, p, "c") * w5TestInt(t, p, "length")
		}
		return w5TestInt(t, p, "c1")*w5TestInt(t, p, "w1") + w5TestInt(t, p, "c2")*w5TestInt(t, p, "w2") + w5TestInt(t, p, "c3")*w5TestInt(t, p, "w3")
	case "monotone_convergence":
		if kp == 1 {
			return w5TestInt(t, p, "limit")
		}
		if kp == 2 {
			return 1
		}
		return w5TestInt(t, p, "value") * w5TestInt(t, p, "measure")
	case "fatou":
		if kp == 1 {
			return w5TestInt(t, p, "pointwise")
		}
		if kp == 2 {
			return 1
		}
		return w5TestInt(t, p, "c")
	case "dominated_convergence":
		if kp == 1 {
			return w5TestInt(t, p, "limit")
		}
		if kp == 2 {
			return w5TestInt(t, p, "bound")
		}
		return 1
	case "product":
		if kp == 1 {
			return w5TestInt(t, p, "a") * w5TestInt(t, p, "b")
		}
		if kp == 2 {
			return w5TestInt(t, p, "a") * w5TestInt(t, p, "b")
		}
		return w5TestInt(t, p, "m") * w5TestInt(t, p, "n")
	case "fubini_tonelli":
		a, b, c := w5TestInt(t, p, "a"), w5TestInt(t, p, "b"), w5TestInt(t, p, "c")
		if kp == 1 || kp == 2 {
			return c * a * b
		}
		s := 0
		for i := 1; i <= a; i++ {
			for j := 1; j <= b; j++ {
				s += i * c
			}
		}
		return s
	case "lp_space":
		a, b := w5TestInt(t, p, "a"), w5TestInt(t, p, "b")
		if kp == 1 {
			return w5RTestAbs(a) + w5RTestAbs(b)
		}
		if kp == 2 {
			return a*a + b*b
		}
		return w5RTestMax(w5RTestAbs(a), w5RTestAbs(b))
	case "holder":
		a, b, c, d := w5TestInt(t, p, "a"), w5TestInt(t, p, "b"), w5TestInt(t, p, "c"), w5TestInt(t, p, "d")
		if kp == 1 {
			return a*c + b*d
		}
		if kp == 2 {
			return 2
		}
		return (a*a + b*b) * (c*c + d*d)
	case "minkowski":
		a, b, c, d := w5TestInt(t, p, "a"), w5TestInt(t, p, "b"), w5TestInt(t, p, "c"), w5TestInt(t, p, "d")
		if kp == 1 {
			return (a+c)*(a+c) + (b+d)*(b+d)
		}
		if kp == 2 {
			return w5RTestAbs(a+c) + w5RTestAbs(b+d)
		}
		return 3
	case "banach":
		if kp == 1 {
			return 1
		}
		if kp == 2 {
			return w5TestInt(t, p, "limit")
		}
		return w5TestInt(t, p, "dimension")
	case "hilbert":
		a, b, c, d := w5TestInt(t, p, "a"), w5TestInt(t, p, "b"), w5TestInt(t, p, "c"), w5TestInt(t, p, "d")
		dot := a*c + b*d
		if kp == 1 {
			return dot
		}
		if kp == 2 {
			if dot == 0 {
				return 1
			}
			return 0
		}
		return a*a + b*b
	case "bounded_operator":
		a, b := w5TestInt(t, p, "a"), w5TestInt(t, p, "b")
		if kp == 1 {
			return w5RTestMax(w5RTestAbs(a), w5RTestAbs(b))
		}
		if kp == 2 {
			return 1
		}
		return a * w5TestInt(t, p, "x")
	case "riesz_representation":
		a, b := w5TestInt(t, p, "a"), w5TestInt(t, p, "b")
		if kp == 1 {
			return a
		}
		if kp == 2 {
			return a*w5TestInt(t, p, "x") + b*w5TestInt(t, p, "y")
		}
		return a*a + b*b
	case "absolute_continuity":
		a, c, d := w5TestInt(t, p, "a"), w5TestInt(t, p, "c"), w5TestInt(t, p, "d")
		if kp == 1 {
			return a * (d - c)
		}
		if kp == 2 {
			return 0
		}
		return w5RTestAbs(a) * (d - c)
	case "lebesgue_differentiation":
		a, b, c, d := w5TestInt(t, p, "a"), w5TestInt(t, p, "b"), w5TestInt(t, p, "c"), w5TestInt(t, p, "d")
		if kp == 1 {
			return a*(c+d)/2 + b
		}
		if kp == 2 {
			return a
		}
		return a*c + b
	}
	t.Fatalf("unknown kind %q", kind)
	return 0
}

func TestWave5RemainingRealDiversity(t *testing.T) {
	kinds := []string{
		"double_series", "fourier_convergence", "measurable_sets", "measurable_functions",
		"lebesgue_integral", "monotone_convergence", "fatou", "dominated_convergence",
		"product", "fubini_tonelli", "lp_space", "holder", "minkowski", "banach",
		"hilbert", "bounded_operator", "riesz_representation", "absolute_continuity",
		"lebesgue_differentiation",
	}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			seen := map[string]bool{}
			for i := 0; i < 60; i++ {
				p := (&w5RARemainingGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				seen[p.Question] = true
			}
			if len(seen) < 5 {
				t.Fatalf("%s kp%d: only %d distinct questions in 60 draws", kind, kp, len(seen))
			}
			if len(seen) < 10 {
				t.Logf("%s kp%d: %d distinct questions, below the 10-stem design target", kind, kp, len(seen))
			}
		}
	}
}
