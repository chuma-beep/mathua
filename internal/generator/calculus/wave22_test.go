package calculus

import (
	"math"
	"strconv"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

var w22Kinds = map[string]string{
	"calc.seq.limit":                    "limit",
	"calc.seq.subsequence":              "subsequence",
	"calc.series.riemann_zeta":          "zeta",
	"calc.series.eta":                   "eta",
	"calc.series.bernoulli":             "bernoulli",
	"calc.series.double":                "double",
	"calc.series.stirling":              "stirling",
	"calc.series.abel_summation":        "abel",
	"calc.series.power_uniqueness":      "power_uniqueness",
	"calc.series.grouping":              "grouping",
	"calc.series.generating":            "generating",
	"calc.series.arithmetico_geometric": "arithmetico_geometric",
	"calc.integral.cross_sections":      "cross_sections",
	"calc.integral.parametric_length":   "parametric_length",
	"calc.integral.parametric_area":     "parametric_area",
	"calc.integral.beta_gamma":          "beta_gamma",
	"calc.integral.unbounded":           "unbounded",
	"calc.integral.pumping":             "pumping",
	"calc.integral.simpson_38":          "simpson_38",
	"calc.integral.stieltjes":           "stieltjes",
	"calc.integral.holder":              "holder",
	"calc.integral.jensen":              "jensen",
	"calc.integral.radial_mass":         "radial_mass",
}

func TestWave22Answers(t *testing.T) {
	reg := generator.NewRegistry()
	registerWave22(reg)
	for id, kind := range w22Kinds {
		for kp := 1; kp <= 3; kp++ {
			gen, err := reg.Get(id)
			if err != nil {
				t.Fatalf("%s: %v", id, err)
			}
			p := gen.Generate(generator.GeneratorContext{KPOrdinal: kp})
			if p.KPOrdinal != kp {
				t.Fatalf("%s kp%d: ordinal = %d", id, kp, p.KPOrdinal)
			}
			if p.Answer == "" || p.Facts["answer"] != p.Answer {
				t.Fatalf("%s kp%d: answer %q facts %q", id, kp, p.Answer, p.Facts["answer"])
			}
			want, ok := w22Recompute(kind, kp, p.Facts)
			if !ok {
				t.Fatalf("%s kp%d: no recomputation", id, kp)
			}
			if want != p.Answer {
				t.Fatalf("%s kp%d: answer %q, recomputed %q (facts %v)", id, kp, p.Answer, want, p.Facts)
			}
		}
	}
}

func TestWave22Diversity(t *testing.T) {
	reg := generator.NewRegistry()
	registerWave22(reg)
	for id := range w22Kinds {
		gen, err := reg.Get(id)
		if err != nil {
			t.Fatalf("%s: %v", id, err)
		}
		seen := map[string]bool{}
		for i := 0; i < 200; i++ {
			p := gen.Generate(generator.GeneratorContext{})
			seen[p.Answer] = true
		}
		if len(seen) < 5 {
			t.Fatalf("%s: only %d distinct answers in 200 draws", id, len(seen))
		}
	}
}

func w22Ati(f map[string]string, k string) int {
	v, _ := strconv.Atoi(f[k])
	return v
}

func w22Recompute(kind string, kp int, f map[string]string) (string, bool) {
	switch kind {
	case "limit":
		switch kp {
		case 1:
			return strconv.Itoa(w22Ati(f, "n0") + 1), true
		case 2:
			return strconv.Itoa(w22Ati(f, "a") + w22Ati(f, "b")), true
		default:
			return strconv.Itoa(w22Ati(f, "k")), true
		}
	case "subsequence":
		switch kp {
		case 1:
			return strconv.Itoa(2 * w22Ati(f, "c")), true
		case 2:
			return strconv.Itoa(w22Ati(f, "b") - w22Ati(f, "a") + 1), true
		default:
			return strconv.Itoa(w22Ati(f, "k")), true
		}
	case "zeta":
		switch kp {
		case 1:
			return strconv.Itoa(1 << uint(w22Ati(f, "s"))), true
		case 2:
			return strconv.Itoa(w22Ati(f, "p") - 1), true
		default:
			return strconv.Itoa(2 * w22Ati(f, "k")), true
		}
	case "eta":
		switch kp {
		case 1:
			return strconv.Itoa((1 << uint(w22Ati(f, "s")-1)) - 1), true
		case 2:
			return strconv.Itoa(2 * w22Ati(f, "k")), true
		default:
			dens := []int{1, 2, 6, 12, 60, 60, 420, 840, 2520, 2520}
			return strconv.Itoa(dens[w22Ati(f, "n")-1]), true
		}
	case "bernoulli":
		switch kp {
		case 1:
			dens := []int{6, 30, 42, 30, 66, 2730, 6, 510, 798, 330, 138, 2730}
			return strconv.Itoa(dens[w22Ati(f, "k")-1]), true
		case 2:
			return strconv.Itoa(2 * w22Ati(f, "k")), true
		default:
			dens := []int{2, 12, 0, 720, 0, 30240, 0, 1209600, 0, 47900160, 0, 1307674368000}
			return strconv.Itoa(dens[w22Ati(f, "n")-1]), true
		}
	case "double":
		switch kp {
		case 1:
			return strconv.Itoa((w22Ati(f, "a") - 1) * (w22Ati(f, "b") - 1)), true
		case 2:
			return strconv.Itoa(w22Ati(f, "p") + w22Ati(f, "q")), true
		default:
			zetaDen := map[int]int{2: 6, 4: 90, 6: 945, 8: 9450}
			return strconv.Itoa(zetaDen[w22Ati(f, "p")] * zetaDen[w22Ati(f, "q")]), true
		}
	case "stirling":
		switch kp {
		case 1:
			return strconv.Itoa(2 * w22Ati(f, "n")), true
		case 2:
			digits := map[int]int{5: 3, 10: 7, 15: 13, 20: 19, 25: 26, 30: 33, 40: 48, 50: 65, 75: 110, 100: 158}
			return strconv.Itoa(digits[w22Ati(f, "n")]), true
		default:
			return strconv.Itoa(12 * w22Ati(f, "n")), true
		}
	case "abel":
		switch kp {
		case 1:
			return strconv.Itoa(w22Ati(f, "c") * w22Ati(f, "n")), true
		case 2:
			n := w22Ati(f, "n")
			return strconv.Itoa(n * (n + 1) / 2), true
		default:
			return strconv.Itoa(w22Ati(f, "m") * w22Ati(f, "b")), true
		}
	case "power_uniqueness":
		switch kp {
		case 1:
			a, k := w22Ati(f, "a"), w22Ati(f, "k")
			ans := 1
			for i := 0; i < k; i++ {
				ans *= a
			}
			return strconv.Itoa(ans), true
		case 2:
			c, n := w22Ati(f, "c"), w22Ati(f, "n")
			ans := 1
			for i := 0; i < n; i++ {
				ans *= c
			}
			return strconv.Itoa(ans), true
		default:
			return strconv.Itoa(w22Ati(f, "k")), true
		}
	case "grouping":
		switch kp {
		case 1:
			k := w22Ati(f, "k")
			return strconv.Itoa((2*k - 1) * (2 * k)), true
		case 2:
			return strconv.Itoa(4*w22Ati(f, "n") - 1), true
		default:
			return strconv.Itoa(1 << uint(w22Ati(f, "k"))), true
		}
	case "generating":
		switch kp {
		case 1:
			return strconv.Itoa(w22Ati(f, "c")), true
		case 2:
			m, k := w22Ati(f, "m"), w22Ati(f, "k")
			return strconv.Itoa(w22Binom(k+m-1, m-1)), true
		default:
			return strconv.Itoa(w22Fib(w22Ati(f, "k"))), true
		}
	case "arithmetico_geometric":
		switch kp {
		case 1:
			return strconv.Itoa(2*w22Ati(f, "a") + 2*w22Ati(f, "d")), true
		case 2:
			return strconv.Itoa(3*w22Ati(f, "A") + 3*w22Ati(f, "D")), true
		default:
			return strconv.Itoa(4*w22Ati(f, "A") + 4*w22Ati(f, "D")), true
		}
	case "cross_sections":
		switch kp {
		case 1:
			s, h := w22Ati(f, "s"), w22Ati(f, "h")
			return strconv.Itoa(s * s * h), true
		case 2:
			b, h := w22Ati(f, "b"), w22Ati(f, "h")
			return strconv.Itoa(b * b * h / 3), true
		default:
			c, a := w22Ati(f, "c"), w22Ati(f, "a")
			return strconv.Itoa(c * a * a * a / 3), true
		}
	case "parametric_length":
		switch kp {
		case 1:
			a, b := w22Ati(f, "a"), w22Ati(f, "b")
			return strconv.Itoa(int(math.Round(math.Sqrt(float64(a*a + b*b))))), true
		case 2:
			return strconv.Itoa(w22Ati(f, "b")), true
		default:
			a, b, c := w22Ati(f, "a"), w22Ati(f, "b"), w22Ati(f, "c")
			return strconv.Itoa(int(math.Round(math.Sqrt(float64(a*a + b*b + c*c))))), true
		}
	case "parametric_area":
		switch kp {
		case 1:
			a := w22Ati(f, "a")
			return strconv.Itoa(a * a * a / 3), true
		case 2:
			a := w22Ati(f, "a")
			return strconv.Itoa(2 * a * a * a / 3), true
		default:
			return strconv.Itoa(w22Ati(f, "a") * w22Ati(f, "b") / 4), true
		}
	case "beta_gamma":
		switch kp {
		case 1:
			p, q := w22Ati(f, "p"), w22Ati(f, "q")
			num := w22Fact(p-1) * w22Fact(q-1)
			den := w22Fact(p + q - 1)
			return strconv.Itoa(den / w22GCD(num, den)), true
		case 2:
			return strconv.Itoa(w22Fact(w22Ati(f, "k") - 1)), true
		default:
			return strconv.Itoa(w22Ati(f, "p")), true
		}
	case "unbounded":
		switch kp {
		case 1:
			return strconv.Itoa(w22Ati(f, "k")), true
		case 2:
			return strconv.Itoa(2 * w22Ati(f, "m")), true
		default:
			return strconv.Itoa(w22Ati(f, "k")), true
		}
	case "pumping":
		switch kp {
		case 1:
			r, h := w22Ati(f, "r"), w22Ati(f, "h")
			return strconv.Itoa(r * r * h * h / 2), true
		case 2:
			l, w, h := w22Ati(f, "l"), w22Ati(f, "w"), w22Ati(f, "h")
			return strconv.Itoa(l * w * h * h / 2), true
		default:
			r, h, e := w22Ati(f, "r"), w22Ati(f, "h"), w22Ati(f, "extra")
			return strconv.Itoa(r * r * h * (h/2 + e)), true
		}
	case "simpson_38":
		switch kp {
		case 1:
			return strconv.Itoa(w22Ati(f, "f0") + 3*w22Ati(f, "f1") + 3*w22Ati(f, "f2") + w22Ati(f, "f3")), true
		case 2:
			sum := 0
			weights := []int{1, 3, 3, 2, 3, 3, 1}
			for i, w := range weights {
				sum += w * w22Ati(f, "f"+strconv.Itoa(i))
			}
			return strconv.Itoa(sum), true
		default:
			return strconv.Itoa((w22Ati(f, "b") - w22Ati(f, "a")) / 3), true
		}
	case "stieltjes":
		switch kp {
		case 1:
			n := w22Ati(f, "n")
			return strconv.Itoa(n * (n + 1) / 2), true
		case 2:
			return strconv.Itoa(w22Ati(f, "c1") + 2*w22Ati(f, "c2") + 3*w22Ati(f, "c3")), true
		default:
			a := w22Ati(f, "a")
			return strconv.Itoa(a * a), true
		}
	case "holder":
		switch kp {
		case 1:
			return strconv.Itoa(w22Ati(f, "p") - 1), true
		case 2:
			return strconv.Itoa(w22Ati(f, "a")*w22Ati(f, "p") + 1), true
		default:
			return strconv.Itoa(2*w22Ati(f, "p") - 1), true
		}
	case "jensen":
		switch kp {
		case 1:
			a, b := w22Ati(f, "a"), w22Ati(f, "b")
			m := (b - a) / 2
			return strconv.Itoa(m * m), true
		case 2:
			l := w22Ati(f, "l")
			return strconv.Itoa(l * l / 12), true
		default:
			a, b := w22Ati(f, "a"), w22Ati(f, "b")
			m := (b - a) / 3
			return strconv.Itoa(2 * m * m), true
		}
	case "radial_mass":
		switch kp {
		case 1:
			r, rho := w22Ati(f, "r"), w22Ati(f, "rho")
			return strconv.Itoa(rho * r * r), true
		case 2:
			r, k := w22Ati(f, "r"), w22Ati(f, "k")
			return strconv.Itoa(2 * k * r * r * r / 3), true
		default:
			r, k := w22Ati(f, "r"), w22Ati(f, "k")
			return strconv.Itoa(k * r * r * r / 3), true
		}
	}
	return "", false
}
