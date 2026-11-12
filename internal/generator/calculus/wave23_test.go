package calculus

import (
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

var w23Kinds = map[string]string{
	"calc.multi.laplace_operator":    "laplace",
	"calc.multi.green_function":      "green",
	"calc.multi.heat_kernel":         "heat",
	"calc.multi.wave_equation":       "wave",
	"calc.multi.separation":          "separation",
	"calc.multi.sturm_liouville":     "sturm_liouville",
	"calc.multi.legendre":            "legendre",
	"calc.multi.bessel":              "bessel",
	"calc.multi.spherical_harmonics": "spherical",
	"calc.multi.delta":               "delta",
	"calc.multi.weak_derivative":     "weak_derivative",
	"calc.multi.sobolev":             "sobolev",
	"calc.multi.weak_solution":       "weak_solution",
	"calc.vec.form_integral":         "form_integral",
	"calc.vec.stokes_manifold":       "stokes_manifold",
	"calc.vec.de_rham":               "de_rham",
	"calc.vec.hodge_star":            "hodge_star",
	"calc.vec.christoffel":           "christoffel",
	"calc.vec.covariant_derivative":  "covariant_derivative",
	"calc.vec.riemann_tensor":        "riemann_tensor",
	"calc.vec.connection_form":       "connection_form",
	"calc.adv.gamma_function":        "gamma_function",
	"calc.adv.asymptotic":            "asymptotic",
}

func TestWave23Answers(t *testing.T) {
	reg := generator.NewRegistry()
	registerWave23(reg)
	for id, kind := range w23Kinds {
		gen, err := reg.Get(id)
		if err != nil {
			t.Fatalf("%s: %v", id, err)
		}
		for kp := 1; kp <= 3; kp++ {
			p := gen.Generate(generator.GeneratorContext{KPOrdinal: kp})
			if p.KPOrdinal != kp {
				t.Fatalf("%s kp%d: ordinal = %d", id, kp, p.KPOrdinal)
			}
			if p.Answer == "" || p.Facts["answer"] != p.Answer {
				t.Fatalf("%s kp%d: answer %q facts %q", id, kp, p.Answer, p.Facts["answer"])
			}
			if !strings.HasSuffix(p.Question, "(enter a number)") {
				t.Fatalf("%s kp%d: question does not end with the numeric prompt: %q", id, kp, p.Question)
			}
			want, ok := w23Recompute(kind, p.Facts)
			if !ok {
				t.Fatalf("%s kp%d: no recomputation for facts %v", id, kp, p.Facts)
			}
			if want != p.Answer {
				t.Fatalf("%s kp%d: answer %q, recomputed %q (facts %v)", id, kp, p.Answer, want, p.Facts)
			}
		}
	}
}

func TestWave23Diversity(t *testing.T) {
	reg := generator.NewRegistry()
	registerWave23(reg)
	for id := range w23Kinds {
		gen, err := reg.Get(id)
		if err != nil {
			t.Fatalf("%s: %v", id, err)
		}
		for kp := 1; kp <= 3; kp++ {
			seen := map[string]bool{}
			for i := 0; i < 200; i++ {
				p := gen.Generate(generator.GeneratorContext{KPOrdinal: kp})
				seen[p.Answer] = true
			}
			if len(seen) < 5 {
				t.Fatalf("%s kp%d: only %d distinct answers in 200 draws", id, kp, len(seen))
			}
		}
	}
}

func w23Ati(f map[string]string, k string) int {
	v, _ := strconv.Atoi(f[k])
	return v
}

// w23Recompute derives the answer from the published facts without consulting
// the generator.
func w23Recompute(kind string, f map[string]string) (string, bool) {
	ask := f["ask"]
	switch kind {
	case "laplace":
		return strconv.Itoa(2*w23Ati(f, "a") + 2*w23Ati(f, "b")), true
	case "green":
		return strconv.Itoa(2*w23Ati(f, "m") - w23Ati(f, "d")), true
	case "heat":
		return strconv.Itoa(-w23Ati(f, "d") / 2), true
	case "wave":
		k := w23Ati(f, "k")
		c := 0
		for i := 0; i*i <= k; i++ {
			if i*i == k {
				c = i
			}
		}
		return strconv.Itoa(c), true
	case "separation":
		if ask == "odes" {
			return strconv.Itoa(2), true
		}
		n := w23Ati(f, "n")
		return strconv.Itoa(n * n), true
	case "sturm_liouville":
		n := w23Ati(f, "n")
		if ask == "zeros" {
			return strconv.Itoa(n - 1), true
		}
		return strconv.Itoa(n * n), true
	case "legendre":
		n := w23Ati(f, "n")
		switch ask {
		case "at_one":
			return strconv.Itoa(1), true
		case "at_minus_one":
			return strconv.Itoa(w23Sign(n)), true
		default:
			return strconv.Itoa(2*n + 1), true
		}
	case "bessel":
		n := w23Ati(f, "n")
		switch ask {
		case "at_zero":
			if n == 0 {
				return strconv.Itoa(1), true
			}
			return strconv.Itoa(0), true
		case "power":
			return strconv.Itoa(n), true
		default:
			return strconv.Itoa(1 << uint(n)), true
		}
	case "spherical":
		l := w23Ati(f, "l")
		switch ask {
		case "count":
			return strconv.Itoa(2*l + 1), true
		case "parity":
			return strconv.Itoa(w23Sign(l)), true
		default:
			return strconv.Itoa((w23Ati(f, "count") - 1) / 2), true
		}
	case "delta":
		switch ask {
		case "sample":
			return strconv.Itoa(w23Ati(f, "c")*w23Ati(f, "a") + w23Ati(f, "d")), true
		case "denominator":
			a := w23Ati(f, "a")
			if a < 0 {
				a = -a
			}
			return strconv.Itoa(a), true
		default:
			return strconv.Itoa(w23Ati(f, "e")), true
		}
	case "weak_derivative":
		switch ask {
		case "scaled":
			return strconv.Itoa(2 * w23Ati(f, "a")), true
		case "x_abs":
			return strconv.Itoa(2), true
		default:
			return strconv.Itoa(1), true
		}
	case "sobolev":
		switch ask {
		case "order":
			return strconv.Itoa(w23Ati(f, "k")), true
		case "exponent":
			return strconv.Itoa(2), true
		default:
			return strconv.Itoa(w23Ati(f, "n")/2 + 1), true
		}
	case "weak_solution":
		switch ask {
		case "boundary", "bilinear":
			return strconv.Itoa(w23Ati(f, "a")), true
		default:
			return strconv.Itoa(w23Ati(f, "m")), true
		}
	case "form_integral":
		switch ask {
		case "square":
			return strconv.Itoa(w23Ati(f, "c")), true
		case "pullback":
			return strconv.Itoa(w23Ati(f, "a") * w23Ati(f, "b")), true
		default:
			return strconv.Itoa(w23Ati(f, "c") * w23Ati(f, "area")), true
		}
	case "stokes_manifold":
		switch ask {
		case "boundary":
			return strconv.Itoa(w23Ati(f, "n") - 1), true
		case "orientation":
			return strconv.Itoa(w23Sign(w23Ati(f, "n"))), true
		default:
			return strconv.Itoa(w23Ati(f, "v")), true
		}
	case "de_rham":
		switch f["space"] {
		case "sphere":
			n, k := w23Ati(f, "n"), w23Ati(f, "k")
			if k == 0 || k == n {
				return strconv.Itoa(1), true
			}
			return strconv.Itoa(0), true
		case "genus":
			return strconv.Itoa(2 * w23Ati(f, "g")), true
		case "wedge":
			return strconv.Itoa(w23Ati(f, "m")), true
		default:
			if w23Ati(f, "k") <= 1 {
				return strconv.Itoa(1), true
			}
			return strconv.Itoa(0), true
		}
	case "hodge_star":
		n, k := w23Ati(f, "n"), w23Ati(f, "k")
		e := k * (n - k)
		if ask == "exponent" {
			return strconv.Itoa(e), true
		}
		if e%2 == 0 {
			return strconv.Itoa(1), true
		}
		return strconv.Itoa(-1), true
	case "christoffel":
		switch ask {
		case "symbols":
			n := w23Ati(f, "n")
			return strconv.Itoa(n * n * n), true
		case "polar":
			return strconv.Itoa(-w23Ati(f, "m")), true
		default:
			return strconv.Itoa(3), true
		}
	case "covariant_derivative":
		switch ask {
		case "basis":
			return strconv.Itoa(1), true
		case "components":
			return strconv.Itoa(w23Ati(f, "n")), true
		default:
			return strconv.Itoa(-w23Ati(f, "r")), true
		}
	case "riemann_tensor":
		switch ask {
		case "components":
			n := w23Ati(f, "n")
			return strconv.Itoa(n * n * n * n), true
		case "independent":
			return strconv.Itoa(1), true
		default:
			return strconv.Itoa(1), true
		}
	case "connection_form":
		switch ask {
		case "curvature", "equations":
			return strconv.Itoa(2), true
		default:
			n := w23Ati(f, "n")
			return strconv.Itoa(n * n), true
		}
	case "gamma_function":
		switch ask {
		case "value":
			return strconv.Itoa(w23Fact(w23Ati(f, "n") - 1)), true
		case "coefficient":
			return strconv.Itoa(w23Ati(f, "x")), true
		default:
			return strconv.Itoa(1), true
		}
	case "asymptotic":
		switch ask {
		case "factorial":
			return strconv.Itoa(w23Fact(w23Ati(f, "n"))), true
		case "sign":
			return strconv.Itoa(w23Sign(w23Ati(f, "n"))), true
		default:
			return strconv.Itoa(w23Ati(f, "x") - 1), true
		}
	}
	return "", false
}
