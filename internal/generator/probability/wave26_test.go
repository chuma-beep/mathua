package probability

import (
	"strconv"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func wave26IDs() []string {
	return []string{
		"prob.stoch.quadratic_variation", "prob.stoch.ito_formula",
		"prob.stoch.ito_isometry", "prob.stoch.girsanov",
		"prob.stoch.fokker_planck", "prob.stoch.martingale_representation",
		"prob.stoch.levy_characterisation",
		"prob.stat.estimator", "prob.stat.bias", "prob.stat.consistency",
		"prob.stat.sufficiency", "prob.stat.cramer_rao", "prob.stat.mle",
		"prob.stat.confidence",
		"prob.mv.normal", "prob.mv.conditioning", "prob.mv.covariance",
		"prob.mv.dirichlet", "prob.mv.copula", "prob.mv.wishart",
		"prob.mv.elliptical",
		"prob.comb.coupon", "prob.comb.birthday", "prob.comb.occupancy",
		"prob.comb.records", "prob.comb.permutation", "prob.comb.hashing",
		"prob.comb.bins",
		"prob.rel.hazard", "prob.rel.survival", "prob.rel.weibull",
		"prob.rel.function", "prob.rel.mtbf", "prob.rel.series",
		"prob.rel.parallel",
	}
}

func wave26Registry(t *testing.T) *generator.Registry {
	t.Helper()
	reg := generator.NewRegistry()
	Register(reg)
	return reg
}

func wave26Fact(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer\nfacts: %v", p.Question, key, p.Facts[key], p.Facts)
	}
	return v
}

func wave26Factorial(n int) int {
	out := 1
	for i := 2; i <= n; i++ {
		out *= i
	}
	return out
}

func TestWave26AnswersAreIntegers(t *testing.T) {
	reg := wave26Registry(t)
	for _, id := range wave26IDs() {
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

func TestWave26QuestionVariety(t *testing.T) {
	reg := wave26Registry(t)
	for _, id := range wave26IDs() {
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

type wave26Case struct {
	id   string
	want func(t *testing.T, p generator.Problem, kp int) int
}

func TestWave26AnswersRecompute(t *testing.T) {
	reg := wave26Registry(t)
	cases := []wave26Case{
		{"prob.stoch.quadratic_variation", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave26Fact(t, p, "t")
			}
			return wave26Fact(t, p, "n")
		}},
		{"prob.stoch.ito_formula", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return 1
			}
			return 2
		}},
		{"prob.stoch.ito_isometry", func(t *testing.T, p generator.Problem, kp int) int {
			c, tt := wave26Fact(t, p, "c"), wave26Fact(t, p, "t")
			if kp == 1 {
				return c * c
			}
			if kp == 2 {
				return tt
			}
			return c * c * tt
		}},
		{"prob.stoch.girsanov", func(t *testing.T, p generator.Problem, kp int) int {
			th := wave26Fact(t, p, "theta")
			if kp == 1 {
				return 1
			}
			if kp == 2 {
				return 2
			}
			return th
		}},
		{"prob.stoch.fokker_planck", func(t *testing.T, p generator.Problem, kp int) int {
			d := wave26Fact(t, p, "d")
			if kp == 1 {
				return d
			}
			if kp == 2 {
				return 2
			}
			return d * d
		}},
		{"prob.stoch.martingale_representation", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave26Fact(t, p, "n")
			if kp == 3 {
				return wave26Fact(t, p, "t")
			}
			return n
		}},
		{"prob.stoch.levy_characterisation", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave26Fact(t, p, "t")
			}
			if kp == 2 {
				return 1
			}
			return 0
		}},
		{"prob.stat.estimator", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave26Fact(t, p, "n")
			if kp == 3 {
				return n - 1
			}
			return n
		}},
		{"prob.stat.bias", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave26Fact(t, p, "c")
			}
			if kp == 2 {
				return wave26Fact(t, p, "n")
			}
			return 0
		}},
		{"prob.stat.consistency", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return 1
			}
			if kp == 2 {
				return 0
			}
			return wave26Fact(t, p, "n")
		}},
		{"prob.stat.sufficiency", func(t *testing.T, p generator.Problem, kp int) int {
			n, k := wave26Fact(t, p, "n"), wave26Fact(t, p, "k")
			if kp == 1 {
				return k
			}
			if kp == 2 {
				return n - k
			}
			return 1
		}},
		{"prob.stat.cramer_rao", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave26Fact(t, p, "n")
			}
			if kp == 2 {
				return wave26Fact(t, p, "info") * wave26Fact(t, p, "n")
			}
			return 1
		}},
		{"prob.stat.mle", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave26Fact(t, p, "k")
			}
			if kp == 2 {
				return wave26Fact(t, p, "n")
			}
			return 1
		}},
		{"prob.stat.confidence", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave26Fact(t, p, "level")
			}
			if kp == 2 {
				return 2
			}
			return wave26Fact(t, p, "n")
		}},
		{"prob.mv.normal", func(t *testing.T, p generator.Problem, kp int) int {
			d := wave26Fact(t, p, "d")
			if kp == 2 {
				return d * (d + 1) / 2
			}
			return d
		}},
		{"prob.mv.conditioning", func(t *testing.T, p generator.Problem, kp int) int {
			d, k := wave26Fact(t, p, "d"), wave26Fact(t, p, "k")
			if kp == 1 {
				return d - k
			}
			if kp == 2 {
				return k
			}
			return 2
		}},
		{"prob.mv.covariance", func(t *testing.T, p generator.Problem, kp int) int {
			d := wave26Fact(t, p, "d")
			if kp == 1 {
				return d * d
			}
			if kp == 2 {
				return d
			}
			return d * (d - 1) / 2
		}},
		{"prob.mv.dirichlet", func(t *testing.T, p generator.Problem, kp int) int {
			k, a := wave26Fact(t, p, "k"), wave26Fact(t, p, "a")
			if kp == 1 {
				return k
			}
			if kp == 2 {
				return k * a
			}
			return a
		}},
		{"prob.mv.copula", func(t *testing.T, p generator.Problem, kp int) int {
			return wave26Fact(t, p, "d")
		}},
		{"prob.mv.wishart", func(t *testing.T, p generator.Problem, kp int) int {
			d := wave26Fact(t, p, "d")
			if kp == 1 {
				return wave26Fact(t, p, "nu")
			}
			if kp == 2 {
				return d
			}
			return d * (d + 1) / 2
		}},
		{"prob.mv.elliptical", func(t *testing.T, p generator.Problem, kp int) int {
			d := wave26Fact(t, p, "d")
			if kp == 1 {
				return d
			}
			if kp == 2 {
				return d * (d + 1) / 2
			}
			return 1
		}},
		{"prob.comb.coupon", func(t *testing.T, p generator.Problem, kp int) int {
			return wave26Fact(t, p, "n")
		}},
		{"prob.comb.birthday", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave26Fact(t, p, "n")
			if kp == 1 {
				return n * (n - 1) / 2
			}
			if kp == 2 {
				return 2
			}
			return n
		}},
		{"prob.comb.occupancy", func(t *testing.T, p generator.Problem, kp int) int {
			m, n := wave26Fact(t, p, "m"), wave26Fact(t, p, "n")
			if kp == 2 {
				return m
			}
			return n
		}},
		{"prob.comb.records", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave26Fact(t, p, "n")
			if kp == 3 {
				return 1
			}
			return n
		}},
		{"prob.comb.permutation", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave26Fact(t, p, "n")
			if kp == 1 {
				return wave26Factorial(n)
			}
			if kp == 2 {
				return n
			}
			return 1
		}},
		{"prob.comb.hashing", func(t *testing.T, p generator.Problem, kp int) int {
			m := wave26Fact(t, p, "m")
			if kp == 2 {
				return wave26Fact(t, p, "n")
			}
			return m * (m - 1) / 2
		}},
		{"prob.comb.bins", func(t *testing.T, p generator.Problem, kp int) int {
			m, n := wave26Fact(t, p, "m"), wave26Fact(t, p, "n")
			if kp == 1 {
				return (m + n - 1) / n
			}
			if kp == 2 {
				return m
			}
			return n
		}},
		{"prob.rel.hazard", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave26Fact(t, p, "n")
			if kp == 2 {
				return 1
			}
			return n
		}},
		{"prob.rel.survival", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave26Fact(t, p, "n")
			}
			if kp == 2 {
				return 1
			}
			return wave26Fact(t, p, "t")
		}},
		{"prob.rel.weibull", func(t *testing.T, p generator.Problem, kp int) int {
			k := wave26Fact(t, p, "k")
			if kp == 1 {
				return k
			}
			if kp == 2 {
				return 1
			}
			return k - 1
		}},
		{"prob.rel.function", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return 1
			}
			if kp == 2 {
				return 0
			}
			return wave26Fact(t, p, "t")
		}},
		{"prob.rel.mtbf", func(t *testing.T, p generator.Problem, kp int) int {
			n, k := wave26Fact(t, p, "n"), wave26Fact(t, p, "k")
			if kp == 1 {
				return n
			}
			return k
		}},
		{"prob.rel.series", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave26Fact(t, p, "n")
			if kp == 2 {
				return 1
			}
			return n
		}},
		{"prob.rel.parallel", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave26Fact(t, p, "n")
			if kp == 2 {
				return 1
			}
			return n
		}},
	}

	if len(cases) != len(wave26IDs()) {
		t.Fatalf("case table has %d entries, want %d", len(cases), len(wave26IDs()))
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
