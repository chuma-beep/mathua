package calculus

import (
	"strconv"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func TestWave5RemainingCalculusAnswers(t *testing.T) {
	kinds := []string{
		"general_stokes", "pullback", "closed_exact", "vector_potential", "helmholtz",
		"gauge", "harmonic", "mean_value", "maximum_principle", "poisson", "dirichlet",
		"neumann", "greens_identity", "fundamental_solution", "euler_lagrange",
		"first_variation", "geodesic", "minimal_surface", "surface_measure", "fourier", "convolution",
	}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			for draw := 0; draw < 30; draw++ {
				p := (&w5CalcRemainingGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				want := w5TestRemainingCalculus(t, kind, kp, p)
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("%s kp%d: answer = %s, want %d\nquestion: %s\nfacts: %v", kind, kp, p.Answer, want, p.Question, p.Facts)
				}
			}
		}
	}
}

func w5TestRemainingCalculus(t *testing.T, kind string, kp int, p generator.Problem) int {
	t.Helper()
	switch kind {
	case "general_stokes":
		a, b := w5TestInt(t, p, "a"), w5TestInt(t, p, "b")
		pp, q := w5TestInt(t, p, "p"), w5TestInt(t, p, "q")
		if kp == 1 {
			return a * b
		}
		if kp == 2 {
			return q - pp
		}
		return (q - pp) * a * b
	case "pullback":
		a, b := w5TestInt(t, p, "a"), w5TestInt(t, p, "b")
		if kp == 1 {
			return a
		}
		return a * b
	case "closed_exact":
		a, b := w5TestInt(t, p, "a"), w5TestInt(t, p, "b")
		if kp == 1 {
			if a == b {
				return 1
			}
			return 0
		}
		if kp == 2 {
			return a*w5TestInt(t, p, "x") + b*w5TestInt(t, p, "y")
		}
		return b - a
	case "vector_potential":
		c, x := w5TestInt(t, p, "c"), w5TestInt(t, p, "x")
		if kp == 1 {
			return c * x
		}
		if kp == 2 {
			return c
		}
		return (c + w5TestInt(t, p, "lambda")) * x
	case "helmholtz":
		a, b, c := w5TestInt(t, p, "a"), w5TestInt(t, p, "b"), w5TestInt(t, p, "c")
		if kp == 1 {
			return a + b + c
		}
		if kp == 2 {
			return b - a
		}
		return 1
	case "gauge":
		a, lam := w5TestInt(t, p, "a"), w5TestInt(t, p, "lambda")
		if kp == 1 {
			return a + lam
		}
		if kp == 2 {
			return 0
		}
		return lam * (w5TestInt(t, p, "x") + w5TestInt(t, p, "y") + w5TestInt(t, p, "z"))
	case "harmonic":
		a, b, c := w5TestInt(t, p, "a"), w5TestInt(t, p, "b"), w5TestInt(t, p, "c")
		if kp == 1 {
			return 2 * (a + b + c)
		}
		if kp == 2 {
			return -(a + b)
		}
		return -a
	case "mean_value":
		return w5TestInt(t, p, "a")*w5TestInt(t, p, "x0") + w5TestInt(t, p, "b")*w5TestInt(t, p, "y0") + w5TestInt(t, p, "c")
	case "maximum_principle":
		if kp == 1 {
			return w5TestInt(t, p, "boundary")
		}
		if kp == 2 {
			return 1
		}
		return w5TestInt(t, p, "constant")
	case "poisson":
		a, b := w5TestInt(t, p, "a"), w5TestInt(t, p, "b")
		if kp == 1 || kp == 3 {
			return 2 * (a + b)
		}
		return w5TestInt(t, p, "constant") / 4
	case "dirichlet":
		return w5TestInt(t, p, "boundary")
	case "neumann":
		if kp == 1 {
			return w5TestInt(t, p, "a")
		}
		if kp == 2 {
			return 0
		}
		return w5TestInt(t, p, "c") + w5TestInt(t, p, "shift")
	case "greens_identity":
		if kp == 3 {
			return 1
		}
		return w5TestInt(t, p, "a") * w5TestInt(t, p, "b")
	case "fundamental_solution":
		if kp == 1 {
			return 2 - w5TestInt(t, p, "n")
		}
		return w5TestInt(t, p, "source")
	case "euler_lagrange":
		a := w5TestInt(t, p, "a")
		if kp == 1 {
			return a
		}
		if kp == 2 {
			return a / 2
		}
		if a == 0 {
			return 1
		}
		return 0
	case "first_variation":
		if kp == 1 {
			return w5TestInt(t, p, "a")
		}
		if kp == 2 {
			return w5TestInt(t, p, "slope")
		}
		return 0
	case "geodesic":
		k := w5TestInt(t, p, "scale")
		if kp == 2 {
			return w5TestInt(t, p, "leg1") * k
		}
		return w5TestInt(t, p, "hypotenuse") * k
	case "minimal_surface":
		if kp == 1 {
			return w5TestInt(t, p, "factor") * w5TestInt(t, p, "w") * w5TestInt(t, p, "h")
		}
		if kp == 2 {
			b := w5TestInt(t, p, "b")
			return 1 + b*b
		}
		return 1
	case "surface_measure":
		a, b := w5TestInt(t, p, "a"), w5TestInt(t, p, "b")
		if kp == 1 {
			return a * b
		}
		if kp == 2 {
			return a * b * w5TestInt(t, p, "m") * w5TestInt(t, p, "n")
		}
		return w5TestInt(t, p, "c") * w5TestInt(t, p, "d")
	case "fourier":
		f00, f01 := w5TestInt(t, p, "f00"), w5TestInt(t, p, "f01")
		f10, f11 := w5TestInt(t, p, "f10"), w5TestInt(t, p, "f11")
		if kp == 1 {
			return f00 + f01 + f10 + f11
		}
		if kp == 2 {
			return f00 - f01 + f10 - f11
		}
		return f00 - f01 - f10 + f11
	case "convolution":
		f := [2][2]int{{w5TestInt(t, p, "f00"), w5TestInt(t, p, "f01")}, {w5TestInt(t, p, "f10"), w5TestInt(t, p, "f11")}}
		g := [2][2]int{{w5TestInt(t, p, "g00"), w5TestInt(t, p, "g01")}, {w5TestInt(t, p, "g10"), w5TestInt(t, p, "g11")}}
		h := [2][2]int{}
		for i := 0; i < 2; i++ {
			for j := 0; j < 2; j++ {
				for a := 0; a < 2; a++ {
					for b := 0; b < 2; b++ {
						h[i][j] += f[a][b] * g[(i-a+2)%2][(j-b+2)%2]
					}
				}
			}
		}
		if kp == 1 {
			return h[0][0]
		}
		if kp == 2 {
			return h[1][1]
		}
		return h[0][0] + h[0][1] + h[1][0] + h[1][1]
	}
	t.Fatalf("unknown kind %q", kind)
	return 0
}

func TestWave5RemainingCalculusDiversity(t *testing.T) {
	kinds := []string{
		"general_stokes", "pullback", "closed_exact", "vector_potential", "helmholtz",
		"gauge", "harmonic", "mean_value", "maximum_principle", "poisson", "dirichlet",
		"neumann", "greens_identity", "fundamental_solution", "euler_lagrange",
		"first_variation", "geodesic", "minimal_surface", "surface_measure", "fourier", "convolution",
	}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			seen := map[string]bool{}
			for i := 0; i < 60; i++ {
				p := (&w5CalcRemainingGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
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
