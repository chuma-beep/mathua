package machinelearning

import (
	"strconv"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func wave22IDs() []string {
	return []string{
		"ml.prob.gaussian_process", "ml.prob.variational", "ml.prob.mcmc",
		"ml.prob.hmc", "ml.prob.bayes_opt", "ml.prob.dirichlet",
		"ml.prob.latent", "ml.prob.gmm", "ml.prob.bayesian_nn",
		"ml.kernel.gram", "ml.kernel.rbf", "ml.kernel.polynomial",
		"ml.kernel.mercer", "ml.kernel.rkhs", "ml.kernel.gp_regression",
		"ml.kernel.kpca", "ml.kernel.mmd", "ml.kernel.multiple",
		"ml.theory.vc", "ml.theory.shattering", "ml.theory.rademacher",
		"ml.theory.pac", "ml.theory.no_free_lunch", "ml.theory.erm",
		"ml.theory.concentration", "ml.theory.stability",
		"ml.theory.sample_complexity",
		"ml.robust.adversarial", "ml.robust.certified", "ml.robust.poisoning",
		"ml.robust.distribution_shift", "ml.robust.domain_adaptation",
		"ml.robust.ood", "ml.robust.fairness",
		"ml.robust.differential_privacy", "ml.robust.federated",
		"ml.causal.dag", "ml.causal.confounding", "ml.causal.do",
		"ml.causal.iv", "ml.causal.propensity", "ml.causal.did",
		"ml.causal.counterfactual", "ml.causal.backdoor", "ml.causal.ate",
	}
}

func wave22Registry(t *testing.T) *generator.Registry {
	t.Helper()
	reg := generator.NewRegistry()
	Register(reg)
	return reg
}

func wave22Fact(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer\nfacts: %v", p.Question, key, p.Facts[key], p.Facts)
	}
	return v
}

func TestWave22AnswersAreIntegers(t *testing.T) {
	reg := wave22Registry(t)
	for _, id := range wave22IDs() {
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

func TestWave22QuestionVariety(t *testing.T) {
	reg := wave22Registry(t)
	for _, id := range wave22IDs() {
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

type wave22Case struct {
	id   string
	want func(t *testing.T, p generator.Problem, kp int) int
}

func TestWave22AnswersRecompute(t *testing.T) {
	reg := wave22Registry(t)
	cases := []wave22Case{
		{"ml.prob.gaussian_process", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave22Fact(t, p, "n")
			if kp == 1 {
				return n * n
			}
			if kp == 2 {
				return n * (n + 1) / 2
			}
			return n
		}},
		{"ml.prob.variational", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return 2
			}
			if kp == 2 {
				return 0
			}
			return wave22Fact(t, p, "m")
		}},
		{"ml.prob.mcmc", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave22Fact(t, p, "n")
			if kp == 1 {
				return n * n
			}
			if kp == 2 {
				return 1
			}
			return n
		}},
		{"ml.prob.hmc", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave22Fact(t, p, "l")
			}
			return 1
		}},
		{"ml.prob.bayes_opt", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave22Fact(t, p, "n")
			if kp == 2 {
				return n + 1
			}
			return n
		}},
		{"ml.prob.dirichlet", func(t *testing.T, p generator.Problem, kp int) int {
			k, a := wave22Fact(t, p, "k"), wave22Fact(t, p, "a")
			if kp == 1 {
				return k
			}
			if kp == 2 {
				return k * a
			}
			return a
		}},
		{"ml.prob.latent", func(t *testing.T, p generator.Problem, kp int) int {
			d, m := wave22Fact(t, p, "d"), wave22Fact(t, p, "m")
			if kp == 2 {
				return d * m
			}
			return d
		}},
		{"ml.prob.gmm", func(t *testing.T, p generator.Problem, kp int) int {
			k, d := wave22Fact(t, p, "k"), wave22Fact(t, p, "d")
			if kp == 1 {
				return k * (2*d + 1)
			}
			if kp == 2 {
				return k - 1
			}
			return 1
		}},
		{"ml.prob.bayesian_nn", func(t *testing.T, p generator.Problem, kp int) int {
			m := wave22Fact(t, p, "m")
			if kp == 2 {
				return 2
			}
			return m
		}},
		{"ml.kernel.gram", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave22Fact(t, p, "n")
			if kp == 1 {
				return n * n
			}
			if kp == 2 {
				return n * (n + 1) / 2
			}
			return n
		}},
		{"ml.kernel.rbf", func(t *testing.T, p generator.Problem, kp int) int {
			s, d2 := wave22Fact(t, p, "sigma"), wave22Fact(t, p, "d2")
			if kp == 1 {
				return 1
			}
			if kp == 2 {
				return 2 * s * s
			}
			return d2
		}},
		{"ml.kernel.polynomial", func(t *testing.T, p generator.Problem, kp int) int {
			n, d := wave22Fact(t, p, "n"), wave22Fact(t, p, "d")
			if kp == 1 {
				return d
			}
			if kp == 2 {
				return w22Binom(n+d-1, d)
			}
			return 1
		}},
		{"ml.kernel.mercer", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave22Fact(t, p, "n")
			if kp == 1 {
				return n
			}
			if kp == 2 {
				return 1
			}
			return n
		}},
		{"ml.kernel.rkhs", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave22Fact(t, p, "n")
			if kp == 1 {
				return n
			}
			if kp == 2 {
				return n * n
			}
			return 1
		}},
		{"ml.kernel.gp_regression", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave22Fact(t, p, "n")
			if kp == 1 {
				return n
			}
			if kp == 2 {
				return n * n
			}
			return wave22Fact(t, p, "s")
		}},
		{"ml.kernel.kpca", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave22Fact(t, p, "n")
			if kp == 1 {
				return n * n
			}
			return n
		}},
		{"ml.kernel.mmd", func(t *testing.T, p generator.Problem, kp int) int {
			m, n := wave22Fact(t, p, "m"), wave22Fact(t, p, "n")
			if kp == 1 {
				return m + n
			}
			if kp == 2 {
				return m*m + n*n + 2*m*n
			}
			return 0
		}},
		{"ml.kernel.multiple", func(t *testing.T, p generator.Problem, kp int) int {
			k, n := wave22Fact(t, p, "k"), wave22Fact(t, p, "n")
			if kp == 1 {
				return k * n * n
			}
			if kp == 2 {
				return 1
			}
			return k
		}},
		{"ml.theory.vc", func(t *testing.T, p generator.Problem, kp int) int {
			d := wave22Fact(t, p, "d")
			if kp == 1 {
				return 2
			}
			if kp == 2 {
				return d + 1
			}
			return 4
		}},
		{"ml.theory.shattering", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave22Fact(t, p, "n")
			if kp == 3 {
				return n
			}
			return 1 << n
		}},
		{"ml.theory.rademacher", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave22Fact(t, p, "n")
			if kp == 1 {
				return 1 << n
			}
			if kp == 2 {
				return 2
			}
			return 2
		}},
		{"ml.theory.pac", func(t *testing.T, p generator.Problem, kp int) int {
			n, e := wave22Fact(t, p, "n"), wave22Fact(t, p, "e")
			if kp == 1 {
				return w22CeilLog2(n)
			}
			if kp == 2 {
				return e
			}
			return 2
		}},
		{"ml.theory.no_free_lunch", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave22Fact(t, p, "n")
			if kp == 1 {
				return 1 << n
			}
			if kp == 2 {
				return 1
			}
			return n
		}},
		{"ml.theory.erm", func(t *testing.T, p generator.Problem, kp int) int {
			n, k := wave22Fact(t, p, "n"), wave22Fact(t, p, "k")
			if kp == 1 {
				return k
			}
			if kp == 2 {
				return n - k
			}
			return n
		}},
		{"ml.theory.concentration", func(t *testing.T, p generator.Problem, kp int) int {
			n, e := wave22Fact(t, p, "n"), wave22Fact(t, p, "e")
			if kp == 1 {
				return 2 * n * e * e
			}
			if kp == 2 {
				return 2
			}
			return n
		}},
		{"ml.theory.stability", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave22Fact(t, p, "n")
			if kp == 1 {
				return 1
			}
			if kp == 2 {
				return 2
			}
			return n
		}},
		{"ml.theory.sample_complexity", func(t *testing.T, p generator.Problem, kp int) int {
			n, e := wave22Fact(t, p, "n"), wave22Fact(t, p, "e")
			if kp == 2 {
				return e
			}
			return w22CeilLog2(n)
		}},
		{"ml.robust.adversarial", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave22Fact(t, p, "e")
			}
			if kp == 2 {
				return wave22Fact(t, p, "d")
			}
			return 1
		}},
		{"ml.robust.certified", func(t *testing.T, p generator.Problem, kp int) int {
			k, l := wave22Fact(t, p, "k"), wave22Fact(t, p, "l")
			if kp == 1 {
				return k
			}
			if kp == 2 {
				return k * l
			}
			return k
		}},
		{"ml.robust.poisoning", func(t *testing.T, p generator.Problem, kp int) int {
			n, k := wave22Fact(t, p, "n"), wave22Fact(t, p, "k")
			if kp == 1 {
				return k
			}
			if kp == 2 {
				return n - k
			}
			return n
		}},
		{"ml.robust.distribution_shift", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave22Fact(t, p, "n")
			if kp == 1 {
				return 0
			}
			return n
		}},
		{"ml.robust.domain_adaptation", func(t *testing.T, p generator.Problem, kp int) int {
			m, n := wave22Fact(t, p, "m"), wave22Fact(t, p, "n")
			if kp == 1 {
				return m + n
			}
			if kp == 2 {
				return m
			}
			return 0
		}},
		{"ml.robust.ood", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave22Fact(t, p, "n")
			if kp == 1 {
				return n * n
			}
			if kp == 2 {
				return n - 1
			}
			return n
		}},
		{"ml.robust.fairness", func(t *testing.T, p generator.Problem, kp int) int {
			k, g := wave22Fact(t, p, "k"), wave22Fact(t, p, "g")
			if kp == 1 {
				return g
			}
			if kp == 2 {
				return 2 * k
			}
			return k
		}},
		{"ml.robust.differential_privacy", func(t *testing.T, p generator.Problem, kp int) int {
			k, e, s := wave22Fact(t, p, "k"), wave22Fact(t, p, "e"), wave22Fact(t, p, "s")
			if kp == 1 {
				return k * e
			}
			if kp == 2 {
				return s
			}
			return k
		}},
		{"ml.robust.federated", func(t *testing.T, p generator.Problem, kp int) int {
			k, r := wave22Fact(t, p, "k"), wave22Fact(t, p, "r")
			if kp == 1 {
				return k * r
			}
			if kp == 2 {
				return k
			}
			return r
		}},
		{"ml.causal.dag", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave22Fact(t, p, "n")
			if kp == 1 {
				return n * (n - 1)
			}
			if kp == 2 {
				return 0
			}
			return n
		}},
		{"ml.causal.confounding", func(t *testing.T, p generator.Problem, kp int) int {
			pth, k := wave22Fact(t, p, "p"), wave22Fact(t, p, "k")
			if kp == 1 {
				return pth
			}
			if kp == 2 {
				return k
			}
			return 2
		}},
		{"ml.causal.do", func(t *testing.T, p generator.Problem, kp int) int {
			k, n := wave22Fact(t, p, "k"), wave22Fact(t, p, "n")
			if kp == 1 {
				return 3
			}
			if kp == 2 {
				return k
			}
			return n - k
		}},
		{"ml.causal.iv", func(t *testing.T, p generator.Problem, kp int) int {
			k := wave22Fact(t, p, "k")
			if kp == 1 {
				return 2
			}
			if kp == 2 {
				return k
			}
			return 3
		}},
		{"ml.causal.propensity", func(t *testing.T, p generator.Problem, kp int) int {
			return wave22Fact(t, p, "n")
		}},
		{"ml.causal.did", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave22Fact(t, p, "n")
			if kp == 1 {
				return 4
			}
			if kp == 2 {
				return 2
			}
			return 2 * n
		}},
		{"ml.causal.counterfactual", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave22Fact(t, p, "n")
			if kp == 1 {
				return 2 * n
			}
			return n
		}},
		{"ml.causal.backdoor", func(t *testing.T, p generator.Problem, kp int) int {
			pth := wave22Fact(t, p, "p")
			if kp == 2 {
				return 0
			}
			return pth
		}},
		{"ml.causal.ate", func(t *testing.T, p generator.Problem, kp int) int {
			nt, nc := wave22Fact(t, p, "nt"), wave22Fact(t, p, "nc")
			if kp == 1 {
				return nt + nc
			}
			if kp == 2 {
				return 2
			}
			return nt
		}},
	}

	if len(cases) != len(wave22IDs()) {
		t.Fatalf("case table has %d entries, want %d", len(cases), len(wave22IDs()))
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
