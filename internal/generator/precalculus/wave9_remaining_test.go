package precalculus

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

func w9rTestChoose(n, k int) int {
	if k < 0 || k > n {
		return 0
	}
	r := 1
	for i := 1; i <= k; i++ {
		r = r * (n - k + i) / i
	}
	return r
}

func w9rTestPow3(n int) int {
	r := 1
	for i := 0; i < n; i++ {
		r *= 3
	}
	return r
}

func TestWave9RemainingPrecalcAnswers(t *testing.T) {
	kinds := []string{
		"inverse_composition", "piecewise_operations", "area_polygon", "sign_chart",
		"asymptote_intersection", "doubling", "logistic_inverse", "summation_properties",
		"power_sums", "conic_parametric", "conic_polar", "polar_symmetry",
		"polar_intersection", "cycloid", "vector_work", "direction_cosines",
		"nonlinear_inequality", "partial_repeated", "partial_quadratic", "multinomial",
	}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			for draw := 0; draw < 40; draw++ {
				p := (&w9rPreGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				want := w9rTestPrecalcAnswer(t, kind, kp, p)
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("%s kp%d: answer = %s, want %d\nquestion: %s\nfacts: %v", kind, kp, p.Answer, want, p.Question, p.Facts)
				}
			}
		}
	}
}

func w9rTestPrecalcAnswer(t *testing.T, kind string, kp int, p generator.Problem) int {
	t.Helper()
	switch kind {
	case "inverse_composition":
		m, b := w9rTestInt(t, p, "m"), w9rTestInt(t, p, "b")
		n, c := w9rTestInt(t, p, "n"), w9rTestInt(t, p, "c")
		x := w9rTestInt(t, p, "x")
		if kp == 1 {
			return m*(n*x+c) + b
		}
		if kp == 2 {
			return n*(m*x+b) + c
		}
		if m*n == 1 && m*c+b == 0 {
			return 1
		}
		return 0
	case "piecewise_operations":
		x := w9rTestInt(t, p, "x")
		fv := w9rTestInt(t, p, "a")*x + w9rTestInt(t, p, "b")
		if x >= 0 {
			fv = w9rTestInt(t, p, "c")*x + w9rTestInt(t, p, "d")
		}
		gv := w9rTestInt(t, p, "e")*x + w9rTestInt(t, p, "f")
		if x >= 1 {
			gv = w9rTestInt(t, p, "g")*x + w9rTestInt(t, p, "h")
		}
		if kp == 1 {
			return fv + gv
		}
		if kp == 2 {
			return fv * gv
		}
		if x < 0 {
			return 1
		}
		if x < 1 {
			return 2
		}
		return 3
	case "area_polygon":
		if kp == 1 {
			return w9rTestInt(t, p, "a") * w9rTestInt(t, p, "b") / 2
		}
		if kp == 2 {
			return w9rTestInt(t, p, "w") * w9rTestInt(t, p, "h")
		}
		return w9rTestInt(t, p, "a") * w9rTestInt(t, p, "b")
	case "sign_chart":
		a, b, x := w9rTestInt(t, p, "a"), w9rTestInt(t, p, "b"), w9rTestInt(t, p, "x")
		if kp == 1 {
			return x - a
		}
		if kp == 2 {
			return x - b
		}
		if (x-a)*(x-b) < 0 {
			return -1
		}
		return 1
	case "asymptote_intersection":
		a, k := w9rTestInt(t, p, "a"), w9rTestInt(t, p, "k")
		x0 := w9rTestInt(t, p, "x0")
		if kp == 1 {
			return -k / a
		}
		if kp == 2 {
			return 0
		}
		return k / (x0 - a)
	case "doubling":
		A, T := w9rTestInt(t, p, "A"), w9rTestInt(t, p, "T")
		k, m := w9rTestInt(t, p, "k"), w9rTestInt(t, p, "m")
		if kp == 1 {
			r := A
			for i := 0; i < k; i++ {
				r *= 2
			}
			return r
		}
		if kp == 2 {
			return k * T
		}
		return m
	case "logistic_inverse":
		z := w9rTestInt(t, p, "z")
		if kp == 1 {
			return z
		}
		if kp == 2 {
			return 1 << z
		}
		return (1 << z) + 1
	case "summation_properties":
		a, b := w9rTestInt(t, p, "a"), w9rTestInt(t, p, "b")
		c, d := w9rTestInt(t, p, "c"), w9rTestInt(t, p, "d")
		n := w9rTestInt(t, p, "n")
		sumK := n * (n + 1) / 2
		if kp == 1 {
			return a*sumK + b*n
		}
		if kp == 2 {
			return c * a * sumK
		}
		return (a+c)*sumK + (b+d)*n
	case "power_sums":
		n := w9rTestInt(t, p, "n")
		if kp == 1 {
			return n * (n + 1) / 2
		}
		if kp == 2 {
			return n * (n + 1) * (2*n + 1) / 6
		}
		s := n * (n + 1) / 2
		return s * s
	case "conic_parametric":
		if kp == 1 {
			return w9rTestInt(t, p, "h") + w9rTestInt(t, p, "a")
		}
		if kp == 2 {
			return w9rTestInt(t, p, "k") + w9rTestInt(t, p, "b")
		}
		return w9rTestInt(t, p, "h") - w9rTestInt(t, p, "a")
	case "conic_polar":
		en, ed, d := w9rTestInt(t, p, "e_num"), w9rTestInt(t, p, "e_den"), w9rTestInt(t, p, "d")
		if kp == 1 {
			if en == ed {
				return 1
			}
			if en > ed {
				return 2
			}
			return 3
		}
		if kp == 2 {
			return ed + en
		}
		return d
	case "polar_symmetry":
		typ := w9rTestInt(t, p, "type")
		if kp == 1 || kp == 2 {
			if typ == 2 {
				return 0
			}
			return 1
		}
		if typ == 1 || typ == 3 {
			return 2
		}
		return 1
	case "polar_intersection":
		if kp == 1 {
			return 2
		}
		if kp == 2 {
			return 60
		}
		return w9rTestInt(t, p, "a")
	case "cycloid":
		a := w9rTestInt(t, p, "a")
		if kp == 1 || kp == 2 {
			return 2 * a
		}
		return 2
	case "vector_work":
		pv, q := w9rTestInt(t, p, "p"), w9rTestInt(t, p, "q")
		a, b := w9rTestInt(t, p, "a"), w9rTestInt(t, p, "b")
		work := pv*a + q*b
		if kp == 1 {
			return work
		}
		if kp == 2 {
			return work / 5
		}
		return 2 * work
	case "direction_cosines":
		a, b, r := w9rTestInt(t, p, "a"), w9rTestInt(t, p, "b"), w9rTestInt(t, p, "r")
		if kp == 1 {
			return b
		}
		if kp == 2 {
			return r * r
		}
		return a
	case "nonlinear_inequality":
		r := w9rTestInt(t, p, "r")
		px, py := w9rTestInt(t, p, "px"), w9rTestInt(t, p, "py")
		dist := px*px + py*py
		if kp == 1 {
			if dist <= r*r {
				return 1
			}
			return 0
		}
		if kp == 2 {
			if dist < r*r {
				return 1
			}
			return 0
		}
		return dist
	case "partial_repeated":
		A, B := w9rTestInt(t, p, "A"), w9rTestInt(t, p, "B")
		n0 := w9rTestInt(t, p, "n0")
		if kp == 1 {
			return B
		}
		if kp == 2 {
			return A
		}
		return A + n0
	case "partial_quadratic":
		n1, n0 := w9rTestInt(t, p, "n1"), w9rTestInt(t, p, "n0")
		pv, q := w9rTestInt(t, p, "p"), w9rTestInt(t, p, "q")
		if kp == 1 {
			return n1 + n0
		}
		if kp == 2 {
			return n0
		}
		return pv*pv - 4*q
	case "multinomial":
		n, a, b := w9rTestInt(t, p, "n"), w9rTestInt(t, p, "a"), w9rTestInt(t, p, "b")
		c := w9rTestInt(t, p, "c")
		if kp == 1 {
			num := 1
			for i := 2; i <= n; i++ {
				num *= i
			}
			den := 1
			for _, x := range []int{a, b, c} {
				for i := 2; i <= x; i++ {
					den *= i
				}
			}
			return num / den
		}
		if kp == 2 {
			return w9rTestChoose(n+2, 2)
		}
		return w9rTestPow3(n)
	}
	t.Fatalf("unknown kind %q", kind)
	return 0
}

func TestWave9RemainingPrecalcDiversity(t *testing.T) {
	kinds := []string{
		"inverse_composition", "piecewise_operations", "area_polygon", "sign_chart",
		"asymptote_intersection", "doubling", "logistic_inverse", "summation_properties",
		"power_sums", "conic_parametric", "conic_polar", "polar_symmetry",
		"polar_intersection", "cycloid", "vector_work", "direction_cosines",
		"nonlinear_inequality", "partial_repeated", "partial_quadratic", "multinomial",
	}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			seen := map[string]bool{}
			for i := 0; i < 60; i++ {
				p := (&w9rPreGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
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
