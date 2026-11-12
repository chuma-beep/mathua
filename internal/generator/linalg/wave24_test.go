package linalg

import (
	"strconv"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func wave24IDs() []string {
	return []string{
		"linalg.ml.tensor_space", "linalg.ml.symmetric_square",
		"linalg.ml.exterior_square", "linalg.ml.wedge",
		"linalg.ml.alternating", "linalg.ml.pfaffian",
		"linalg.ml.bilinear_rank", "linalg.ml.quadratic_rank",
		"linalg.ml.inertia", "linalg.ml.polarization",
		"linalg.func.dual", "linalg.func.annihilator",
		"linalg.func.quotient", "linalg.func.exact",
		"linalg.func.adjoint", "linalg.func.riesz",
		"linalg.func.invariant", "linalg.func.cyclic",
		"linalg.func.minimal_polynomial", "linalg.func.direct_sum",
		"linalg.itr.jacobi", "linalg.itr.gauss_seidel", "linalg.itr.sor",
		"linalg.itr.residual", "linalg.itr.convergence",
		"linalg.itr.krylov", "linalg.itr.arnoldi", "linalg.itr.gmres",
		"linalg.itr.cg", "linalg.itr.preconditioner",
		"linalg.stab.forward_error", "linalg.stab.backward_error",
		"linalg.stab.machine_epsilon", "linalg.stab.growth",
		"linalg.stab.pivoting", "linalg.stab.condition_error",
		"linalg.stab.residual_bound", "linalg.stab.componentwise",
		"linalg.stab.perturbation", "linalg.stab.rank_revealing",
	}
}

func wave24Registry(t *testing.T) *generator.Registry {
	t.Helper()
	reg := generator.NewRegistry()
	Register(reg)
	return reg
}

func wave24Fact(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer\nfacts: %v", p.Question, key, p.Facts[key], p.Facts)
	}
	return v
}

func TestWave24AnswersAreIntegers(t *testing.T) {
	reg := wave24Registry(t)
	for _, id := range wave24IDs() {
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

func TestWave24QuestionVariety(t *testing.T) {
	reg := wave24Registry(t)
	for _, id := range wave24IDs() {
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

type wave24Case struct {
	id   string
	want func(t *testing.T, p generator.Problem, kp int) int
}

func TestWave24AnswersRecompute(t *testing.T) {
	reg := wave24Registry(t)
	cases := []wave24Case{
		{"linalg.ml.tensor_space", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 3 {
				return 1
			}
			return wave24Fact(t, p, "m") * wave24Fact(t, p, "n")
		}},
		{"linalg.ml.symmetric_square", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave24Fact(t, p, "n")
			if kp == 2 {
				return n
			}
			return n * (n + 1) / 2
		}},
		{"linalg.ml.exterior_square", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave24Fact(t, p, "n")
			if kp == 3 {
				return 0
			}
			return n * (n - 1) / 2
		}},
		{"linalg.ml.wedge", func(t *testing.T, p generator.Problem, kp int) int {
			n, k := wave24Fact(t, p, "n"), wave24Fact(t, p, "k")
			if kp == 1 {
				return w24Binom(n, 2)
			}
			if kp == 2 {
				return w24Binom(n, k)
			}
			return 0
		}},
		{"linalg.ml.alternating", func(t *testing.T, p generator.Problem, kp int) int {
			n, k := wave24Fact(t, p, "n"), wave24Fact(t, p, "k")
			if kp == 1 {
				return w24Binom(n, k)
			}
			if kp == 2 {
				return w24Fact(k)
			}
			return 1
		}},
		{"linalg.ml.pfaffian", func(t *testing.T, p generator.Problem, kp int) int {
			m := wave24Fact(t, p, "m")
			if kp == 1 {
				return 2 * wave24Fact(t, p, "k")
			}
			if kp == 2 {
				out := 1
				for i := 1; i <= 2*m-1; i += 2 {
					out *= i
				}
				return out
			}
			return 2
		}},
		{"linalg.ml.bilinear_rank", func(t *testing.T, p generator.Problem, kp int) int {
			n, r := wave24Fact(t, p, "n"), wave24Fact(t, p, "r")
			if kp == 2 {
				return n - r
			}
			if kp == 3 {
				return 1
			}
			return r
		}},
		{"linalg.ml.quadratic_rank", func(t *testing.T, p generator.Problem, kp int) int {
			n, r := wave24Fact(t, p, "n"), wave24Fact(t, p, "r")
			if kp == 2 {
				return 2
			}
			if kp == 3 {
				return n - r
			}
			return r
		}},
		{"linalg.ml.inertia", func(t *testing.T, p generator.Problem, kp int) int {
			pv, q := wave24Fact(t, p, "p"), wave24Fact(t, p, "q")
			if kp == 1 {
				return pv + q
			}
			if kp == 2 {
				return pv - q
			}
			return pv
		}},
		{"linalg.ml.polarization", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return 4
			}
			if kp == 2 {
				return 2
			}
			return 1
		}},
		{"linalg.func.dual", func(t *testing.T, p generator.Problem, kp int) int {
			return wave24Fact(t, p, "n")
		}},
		{"linalg.func.annihilator", func(t *testing.T, p generator.Problem, kp int) int {
			n, k := wave24Fact(t, p, "n"), wave24Fact(t, p, "k")
			if kp == 2 {
				return 0
			}
			if kp == 3 {
				return k
			}
			return n - k
		}},
		{"linalg.func.quotient", func(t *testing.T, p generator.Problem, kp int) int {
			n, k := wave24Fact(t, p, "n"), wave24Fact(t, p, "k")
			if kp == 1 {
				return n - k
			}
			return k
		}},
		{"linalg.func.exact", func(t *testing.T, p generator.Problem, kp int) int {
			a, c := wave24Fact(t, p, "a"), wave24Fact(t, p, "c")
			if kp == 1 {
				return 0
			}
			if kp == 2 {
				return a + c
			}
			return 3
		}},
		{"linalg.func.adjoint", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave24Fact(t, p, "n")
			if kp == 2 {
				return n * n
			}
			if kp == 3 {
				return 0
			}
			return n
		}},
		{"linalg.func.riesz", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 2 {
				return wave24Fact(t, p, "n")
			}
			return 1
		}},
		{"linalg.func.invariant", func(t *testing.T, p generator.Problem, kp int) int {
			n, k := wave24Fact(t, p, "n"), wave24Fact(t, p, "k")
			if kp == 1 {
				return k
			}
			if kp == 2 {
				return 1
			}
			return n - k
		}},
		{"linalg.func.cyclic", func(t *testing.T, p generator.Problem, kp int) int {
			d, n := wave24Fact(t, p, "d"), wave24Fact(t, p, "n")
			if kp == 3 {
				return n
			}
			return d
		}},
		{"linalg.func.minimal_polynomial", func(t *testing.T, p generator.Problem, kp int) int {
			n, k := wave24Fact(t, p, "n"), wave24Fact(t, p, "k")
			if kp == 1 {
				return n
			}
			if kp == 2 {
				return k
			}
			return 1
		}},
		{"linalg.func.direct_sum", func(t *testing.T, p generator.Problem, kp int) int {
			n, m := wave24Fact(t, p, "n"), wave24Fact(t, p, "m")
			if kp == 1 {
				return n - m
			}
			if kp == 2 {
				return 0
			}
			return m
		}},
		{"linalg.itr.jacobi", func(t *testing.T, p generator.Problem, kp int) int {
			n, k := wave24Fact(t, p, "n"), wave24Fact(t, p, "k")
			if kp == 1 {
				return n
			}
			if kp == 2 {
				return k
			}
			return 1
		}},
		{"linalg.itr.gauss_seidel", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave24Fact(t, p, "n")
			if kp == 1 {
				return n
			}
			if kp == 2 {
				return 2
			}
			return 1
		}},
		{"linalg.itr.sor", func(t *testing.T, p generator.Problem, kp int) int {
			n, k := wave24Fact(t, p, "n"), wave24Fact(t, p, "k")
			if kp == 1 {
				return 1
			}
			if kp == 2 {
				return k
			}
			return n
		}},
		{"linalg.itr.residual", func(t *testing.T, p generator.Problem, kp int) int {
			n, m := wave24Fact(t, p, "n"), wave24Fact(t, p, "m")
			if kp == 2 {
				return m
			}
			if kp == 3 {
				return 1
			}
			return n
		}},
		{"linalg.itr.convergence", func(t *testing.T, p generator.Problem, kp int) int {
			k, d := wave24Fact(t, p, "k"), wave24Fact(t, p, "d")
			if kp == 2 {
				return d * k
			}
			return k
		}},
		{"linalg.itr.krylov", func(t *testing.T, p generator.Problem, kp int) int {
			k := wave24Fact(t, p, "k")
			if kp == 2 {
				return k - 1
			}
			return k
		}},
		{"linalg.itr.arnoldi", func(t *testing.T, p generator.Problem, kp int) int {
			k := wave24Fact(t, p, "k")
			if kp == 2 {
				return k - 1
			}
			return k
		}},
		{"linalg.itr.gmres", func(t *testing.T, p generator.Problem, kp int) int {
			return wave24Fact(t, p, "k")
		}},
		{"linalg.itr.cg", func(t *testing.T, p generator.Problem, kp int) int {
			k := wave24Fact(t, p, "k")
			if kp == 3 {
				return 1
			}
			return k
		}},
		{"linalg.itr.preconditioner", func(t *testing.T, p generator.Problem, kp int) int {
			k := wave24Fact(t, p, "k")
			if kp == 3 {
				return 1
			}
			return k
		}},
		{"linalg.stab.forward_error", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave24Fact(t, p, "n")
			if kp == 3 {
				return 2
			}
			return n
		}},
		{"linalg.stab.backward_error", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave24Fact(t, p, "n")
			if kp == 1 {
				return n
			}
			if kp == 2 {
				return 1
			}
			return n * n
		}},
		{"linalg.stab.machine_epsilon", func(t *testing.T, p generator.Problem, kp int) int {
			b := wave24Fact(t, p, "b")
			if kp == 1 {
				return b - 1
			}
			if kp == 2 {
				return int(float64(b) * 0.30103)
			}
			return b
		}},
		{"linalg.stab.growth", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave24Fact(t, p, "n")
			if kp == 3 {
				return 1
			}
			return n - 1
		}},
		{"linalg.stab.pivoting", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave24Fact(t, p, "n")
			if kp == 1 {
				return n - 1
			}
			if kp == 2 {
				return 1
			}
			return w24Fact(n)
		}},
		{"linalg.stab.condition_error", func(t *testing.T, p generator.Problem, kp int) int {
			k, d := wave24Fact(t, p, "k"), wave24Fact(t, p, "d")
			if kp == 1 {
				return k
			}
			return d
		}},
		{"linalg.stab.residual_bound", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave24Fact(t, p, "n")
			if kp == 1 {
				return n
			}
			if kp == 2 {
				return 1
			}
			return 2
		}},
		{"linalg.stab.componentwise", func(t *testing.T, p generator.Problem, kp int) int {
			return wave24Fact(t, p, "n")
		}},
		{"linalg.stab.perturbation", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave24Fact(t, p, "n")
			if kp == 1 {
				return 1
			}
			return n
		}},
		{"linalg.stab.rank_revealing", func(t *testing.T, p generator.Problem, kp int) int {
			n, r := wave24Fact(t, p, "n"), wave24Fact(t, p, "r")
			if kp == 1 {
				return r
			}
			if kp == 2 {
				return 1
			}
			return n - r
		}},
	}

	if len(cases) != len(wave24IDs()) {
		t.Fatalf("case table has %d entries, want %d", len(cases), len(wave24IDs()))
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
