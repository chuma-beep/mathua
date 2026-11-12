package abstract

import (
	"strconv"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func wave23IDs() []string {
	return []string{
		"abstract.group.conjugacy", "abstract.group.normalizer",
		"abstract.group.derived_series", "abstract.group.frattini",
		"abstract.group.hall", "abstract.group.exponent",
		"abstract.group.transfer", "abstract.group.cohomology",
		"abstract.group.burnside_lemma", "abstract.group.permutation_character",
		"abstract.rings.groebner", "abstract.rings.regular",
		"abstract.rings.krull", "abstract.rings.integral_closure",
		"abstract.rings.dvr", "abstract.rings.artinian",
		"abstract.rings.local", "abstract.rings.symmetric_algebra",
		"abstract.rings.flat", "abstract.rings.completion",
		"abstract.structures.projective", "abstract.structures.injective",
		"abstract.structures.spectral_sequence", "abstract.structures.chain_homotopy",
		"abstract.structures.derived_functor", "abstract.structures.exact_functor",
		"abstract.structures.direct_limit", "abstract.structures.free_resolution",
		"abstract.field.perfect", "abstract.field.normal",
		"abstract.field.primitive_element", "abstract.field.artin",
		"abstract.field.kummer", "abstract.field.hilbert90",
		"abstract.field.local", "abstract.field.padic",
		"abstract.rep.orthogonality", "abstract.rep.regular",
		"abstract.rep.tensor", "abstract.rep.restriction",
		"abstract.rep.frobenius", "abstract.rep.fourier",
		"abstract.rep.young", "abstract.rep.burnside_pq",
	}
}

func wave23Registry(t *testing.T) *generator.Registry {
	t.Helper()
	reg := generator.NewRegistry()
	Register(reg)
	return reg
}

func wave23Fact(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer\nfacts: %v", p.Question, key, p.Facts[key], p.Facts)
	}
	return v
}

func wave23Factorial(n int) int {
	out := 1
	for i := 2; i <= n; i++ {
		out *= i
	}
	return out
}

func wave23Binom(n, k int) int {
	if k < 0 || k > n {
		return 0
	}
	if k > n-k {
		k = n - k
	}
	out := 1
	for i := 1; i <= k; i++ {
		out = out * (n - k + i) / i
	}
	return out
}

func wave23Pow(a, b int) int {
	out := 1
	for i := 0; i < b; i++ {
		out *= a
	}
	return out
}

func wave23LCM(a, b int) int {
	x, y := a, b
	for y != 0 {
		x, y = y, x%y
	}
	return a / x * b
}

func wave23Totient(n int) int {
	out := n
	for p := 2; p*p <= n; p++ {
		if n%p != 0 {
			continue
		}
		for n%p == 0 {
			n /= p
		}
		out -= out / p
	}
	if n > 1 {
		out -= out / n
	}
	return out
}

func wave23Partitions(n int) int {
	dp := make([]int, n+1)
	dp[0] = 1
	for p := 1; p <= n; p++ {
		for i := p; i <= n; i++ {
			dp[i] += dp[i-p]
		}
	}
	return dp[n]
}

func TestWave23AnswersAreIntegers(t *testing.T) {
	reg := wave23Registry(t)
	for _, id := range wave23IDs() {
		gen, err := reg.Get(id)
		if err != nil {
			t.Fatalf("registry: %v", err)
		}
		for kp := 1; kp <= 3; kp++ {
			for i := 0; i < 8; i++ {
				p := gen.Generate(generator.GeneratorContext{KPOrdinal: kp})
				if _, err := strconv.Atoi(p.Answer); err != nil {
					t.Fatalf("%s kp%d: answer %q is not an integer", id, kp, p.Answer)
				}
				if p.Facts["answer"] != p.Answer {
					t.Fatalf("%s kp%d: facts answer %q does not match %q", id, kp, p.Facts["answer"], p.Answer)
				}
				if p.Explanation == "" {
					t.Fatalf("%s kp%d: empty explanation", id, kp)
				}
				if p.KPOrdinal != kp {
					t.Fatalf("%s kp%d: ordinal %d", id, kp, p.KPOrdinal)
				}
			}
		}
	}
}

func TestWave23QuestionVariety(t *testing.T) {
	reg := wave23Registry(t)
	for _, id := range wave23IDs() {
		gen, err := reg.Get(id)
		if err != nil {
			t.Fatalf("registry: %v", err)
		}
		for kp := 1; kp <= 3; kp++ {
			seen := map[string]bool{}
			for i := 0; i < 12; i++ {
				seen[gen.Generate(generator.GeneratorContext{KPOrdinal: kp}).Question] = true
			}
			if len(seen) < 5 {
				t.Fatalf("%s kp%d: only %d distinct questions in 12 draws", id, kp, len(seen))
			}
		}
	}
}

type wave23Case struct {
	id   string
	want func(t *testing.T, p generator.Problem, kp int) int
}

func TestWave23AnswersRecompute(t *testing.T) {
	reg := wave23Registry(t)
	cases := []wave23Case{
		{"abstract.group.conjugacy", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave23Fact(t, p, "n")
			switch kp {
			case 1:
				return wave23Partitions(n)
			case 2:
				return wave23Binom(n, 2)
			}
			return wave23Factorial(n)
		}},
		{"abstract.group.normalizer", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave23Fact(t, p, "n")
			if kp == 1 {
				return n
			}
			if kp == 2 {
				return wave23Factorial(n) / n
			}
			return 1
		}},
		{"abstract.group.derived_series", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave23Fact(t, p, "n")
			if kp == 1 {
				if n == 4 {
					return 3
				}
				return 2
			}
			if kp == 2 {
				return wave23Factorial(n) / 2
			}
			return 1
		}},
		{"abstract.group.frattini", func(t *testing.T, p generator.Problem, kp int) int {
			pp, k := wave23Fact(t, p, "p"), wave23Fact(t, p, "k")
			if kp == 1 {
				return wave23Pow(pp, k-1)
			}
			if kp == 2 {
				return pp
			}
			return k
		}},
		{"abstract.group.hall", func(t *testing.T, p generator.Problem, kp int) int {
			pp, q := wave23Fact(t, p, "p"), wave23Fact(t, p, "q")
			a, b := wave23Fact(t, p, "a"), wave23Fact(t, p, "b")
			if kp == 1 {
				return wave23Pow(pp, a)
			}
			if kp == 2 {
				return wave23Pow(q, b)
			}
			return 2
		}},
		{"abstract.group.exponent", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave23Fact(t, p, "n")
			if kp == 1 {
				return n
			}
			if kp == 2 {
				out := 1
				for i := 2; i <= n; i++ {
					out = wave23LCM(out, i)
				}
				return out
			}
			return wave23LCM(2, n)
		}},
		{"abstract.group.transfer", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave23Fact(t, p, "p") + 1
			}
			return 2
		}},
		{"abstract.group.cohomology", func(t *testing.T, p generator.Problem, kp int) int {
			r := wave23Fact(t, p, "r")
			if kp == 1 {
				return r
			}
			if kp == 2 {
				return 1
			}
			return 1 - r
		}},
		{"abstract.group.burnside_lemma", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return 1
			}
			if kp == 2 {
				return wave23Fact(t, p, "m")
			}
			return wave23Fact(t, p, "n")
		}},
		{"abstract.group.permutation_character", func(t *testing.T, p generator.Problem, kp int) int {
			n, k := wave23Fact(t, p, "n"), wave23Fact(t, p, "k")
			if kp == 1 {
				return n
			}
			if kp == 2 {
				return wave23Binom(n, k)
			}
			return 1
		}},
		{"abstract.rings.groebner", func(t *testing.T, p generator.Problem, kp int) int {
			a, b := wave23Fact(t, p, "a"), wave23Fact(t, p, "b")
			if kp == 1 {
				return a * b
			}
			if kp == 2 {
				return a
			}
			return 1
		}},
		{"abstract.rings.regular", func(t *testing.T, p generator.Problem, kp int) int {
			d := wave23Fact(t, p, "d")
			if kp == 3 {
				return 1
			}
			return d
		}},
		{"abstract.rings.krull", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave23Fact(t, p, "n")
			}
			return 1
		}},
		{"abstract.rings.integral_closure", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave23Fact(t, p, "n")
			if kp == 1 {
				return 2
			}
			if kp == 2 {
				return n
			}
			return 1
		}},
		{"abstract.rings.dvr", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 2 {
				return wave23Fact(t, p, "p")
			}
			return 1
		}},
		{"abstract.rings.artinian", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave23Fact(t, p, "n")
			if kp == 2 {
				return 1
			}
			return n
		}},
		{"abstract.rings.local", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 2 {
				return wave23Fact(t, p, "p")
			}
			if kp == 3 {
				return wave23Fact(t, p, "n")
			}
			return 1
		}},
		{"abstract.rings.symmetric_algebra", func(t *testing.T, p generator.Problem, kp int) int {
			n, d := wave23Fact(t, p, "n"), wave23Fact(t, p, "d")
			if kp == 1 {
				return wave23Binom(n+d-1, d)
			}
			if kp == 2 {
				return n
			}
			return 1
		}},
		{"abstract.rings.flat", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 3 {
				return 0
			}
			return 1
		}},
		{"abstract.rings.completion", func(t *testing.T, p generator.Problem, kp int) int {
			pp, k := wave23Fact(t, p, "p"), wave23Fact(t, p, "k")
			if kp == 1 {
				return wave23Pow(pp, k)
			}
			if kp == 2 {
				return pp
			}
			return k
		}},
		{"abstract.structures.projective", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return 0
			}
			return wave23Fact(t, p, "n")
		}},
		{"abstract.structures.injective", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave23Fact(t, p, "n")
			}
			return 0
		}},
		{"abstract.structures.spectral_sequence", func(t *testing.T, p generator.Problem, kp int) int {
			r := wave23Fact(t, p, "r")
			if kp == 1 {
				return r * r
			}
			if kp == 2 {
				return 1
			}
			return r
		}},
		{"abstract.structures.chain_homotopy", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave23Fact(t, p, "n")
			if kp == 2 {
				return n - 1
			}
			return n
		}},
		{"abstract.structures.derived_functor", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave23Fact(t, p, "n")
			if kp == 3 {
				return n
			}
			return n + 1
		}},
		{"abstract.structures.exact_functor", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave23Fact(t, p, "n")
			if kp == 1 {
				return n - 1
			}
			if kp == 2 {
				return n - 2
			}
			return n
		}},
		{"abstract.structures.direct_limit", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave23Fact(t, p, "n")
			if kp == 1 {
				return n - 1
			}
			if kp == 2 {
				return n
			}
			return 1
		}},
		{"abstract.structures.free_resolution", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave23Fact(t, p, "n")
			if kp == 1 {
				return n + 1
			}
			if kp == 2 {
				return n
			}
			return 1
		}},
		{"abstract.field.perfect", func(t *testing.T, p generator.Problem, kp int) int {
			pp, e := wave23Fact(t, p, "p"), wave23Fact(t, p, "e")
			if kp == 3 {
				return wave23Pow(pp, e)
			}
			return 1
		}},
		{"abstract.field.normal", func(t *testing.T, p generator.Problem, kp int) int {
			return wave23Totient(wave23Fact(t, p, "n"))
		}},
		{"abstract.field.primitive_element", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 2 {
				return wave23Fact(t, p, "n")
			}
			return 1
		}},
		{"abstract.field.artin", func(t *testing.T, p generator.Problem, kp int) int {
			return wave23Fact(t, p, "n")
		}},
		{"abstract.field.kummer", func(t *testing.T, p generator.Problem, kp int) int {
			return wave23Fact(t, p, "k")
		}},
		{"abstract.field.hilbert90", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave23Fact(t, p, "n")
			if kp == 2 {
				return 1
			}
			return n
		}},
		{"abstract.field.local", func(t *testing.T, p generator.Problem, kp int) int {
			pp, e, ff := wave23Fact(t, p, "p"), wave23Fact(t, p, "e"), wave23Fact(t, p, "f")
			if kp == 1 {
				return e * ff
			}
			if kp == 2 {
				return wave23Pow(pp, ff)
			}
			return e
		}},
		{"abstract.field.padic", func(t *testing.T, p generator.Problem, kp int) int {
			pp, k := wave23Fact(t, p, "p"), wave23Fact(t, p, "k")
			if kp == 2 {
				return pp
			}
			return k
		}},
		{"abstract.rep.orthogonality", func(t *testing.T, p generator.Problem, kp int) int {
			n, c := wave23Fact(t, p, "n"), wave23Fact(t, p, "c")
			if kp == 1 {
				return n
			}
			if kp == 2 {
				return c
			}
			return 1
		}},
		{"abstract.rep.regular", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave23Fact(t, p, "n")
			}
			if kp == 2 {
				return 1
			}
			return wave23Fact(t, p, "d")
		}},
		{"abstract.rep.tensor", func(t *testing.T, p generator.Problem, kp int) int {
			m, n := wave23Fact(t, p, "m"), wave23Fact(t, p, "n")
			if kp == 1 {
				return m * n
			}
			if kp == 2 {
				return m * (m + 1) / 2
			}
			return m * (m - 1) / 2
		}},
		{"abstract.rep.restriction", func(t *testing.T, p generator.Problem, kp int) int {
			n, k := wave23Fact(t, p, "n"), wave23Fact(t, p, "k")
			if kp == 1 {
				return n / k
			}
			if kp == 2 {
				return k
			}
			return 1
		}},
		{"abstract.rep.frobenius", func(t *testing.T, p generator.Problem, kp int) int {
			d, k := wave23Fact(t, p, "d"), wave23Fact(t, p, "k")
			if kp == 1 {
				return d * k
			}
			if kp == 2 {
				return k
			}
			return 1
		}},
		{"abstract.rep.fourier", func(t *testing.T, p generator.Problem, kp int) int {
			return wave23Fact(t, p, "n")
		}},
		{"abstract.rep.young", func(t *testing.T, p generator.Problem, kp int) int {
			n, k := wave23Fact(t, p, "n"), wave23Fact(t, p, "k")
			if kp == 1 {
				return n
			}
			if kp == 2 {
				return wave23Binom(n, k)
			}
			return 2
		}},
		{"abstract.rep.burnside_pq", func(t *testing.T, p generator.Problem, kp int) int {
			pp, q := wave23Fact(t, p, "p"), wave23Fact(t, p, "q")
			a, b := wave23Fact(t, p, "a"), wave23Fact(t, p, "b")
			if kp == 1 {
				return 1
			}
			if kp == 2 {
				return 2
			}
			return wave23Pow(pp, a) * wave23Pow(q, b)
		}},
	}

	if len(cases) != len(wave23IDs()) {
		t.Fatalf("case table has %d entries, want %d", len(cases), len(wave23IDs()))
	}
	for _, tc := range cases {
		gen, err := reg.Get(tc.id)
		if err != nil {
			t.Fatalf("%s: %v", tc.id, err)
		}
		for kp := 1; kp <= 3; kp++ {
			for i := 0; i < 6; i++ {
				p := gen.Generate(generator.GeneratorContext{KPOrdinal: kp})
				want := tc.want(t, p, kp)
				got, err := strconv.Atoi(p.Answer)
				if err != nil {
					t.Fatalf("%s kp%d: answer %q", tc.id, kp, p.Answer)
				}
				if got != want {
					t.Fatalf("%s kp%d: answer %d, want %d\nquestion: %s\nfacts: %v", tc.id, kp, got, want, p.Question, p.Facts)
				}
			}
		}
	}
}
