package discrete

import (
	"strconv"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func wave27IDs() []string {
	return []string{
		"discrete.alg.quicksort", "discrete.alg.merge_sort",
		"discrete.alg.binary_search", "discrete.alg.dijkstra",
		"discrete.alg.huffman", "discrete.alg.backtracking",
		"discrete.alg.branch_bound",
		"discrete.coding.entropy", "discrete.coding.convolutional",
		"discrete.coding.ldpc", "discrete.coding.bch",
		"discrete.coding.singleton",
		"discrete.complexity.space", "discrete.complexity.hierarchy",
		"discrete.complexity.conp", "discrete.complexity.pspace",
		"discrete.complexity.approximation", "discrete.complexity.fpt",
		"discrete.seq.homogeneous", "discrete.seq.characteristic",
		"discrete.seq.nonhomogeneous", "discrete.seq.zeckendorf",
		"discrete.seq.continued_fraction", "discrete.seq.ackermann",
		"discrete.sets.cardinality", "discrete.sets.cartesian",
		"discrete.sets.equivalence_classes", "discrete.sets.lattice",
		"discrete.proof.invariant", "discrete.proof.diagonal",
		"discrete.proof.descent", "discrete.proof.constructive",
		"discrete.proof.nonconstructive",
	}
}

func wave27Registry(t *testing.T) *generator.Registry {
	t.Helper()
	reg := generator.NewRegistry()
	Register(reg)
	return reg
}

func wave27Fact(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer\nfacts: %v", p.Question, key, p.Facts[key], p.Facts)
	}
	return v
}

func TestWave27AnswersAreIntegers(t *testing.T) {
	reg := wave27Registry(t)
	for _, id := range wave27IDs() {
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

func TestWave27QuestionVariety(t *testing.T) {
	reg := wave27Registry(t)
	for _, id := range wave27IDs() {
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

type wave27Case struct {
	id   string
	want func(t *testing.T, p generator.Problem, kp int) int
}

func TestWave27AnswersRecompute(t *testing.T) {
	reg := wave27Registry(t)
	cases := []wave27Case{
		{"discrete.alg.quicksort", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave27Fact(t, p, "n")
			if kp == 1 {
				return n * (n - 1) / 2
			}
			if kp == 2 {
				return w27CeilLog2(n)
			}
			return n
		}},
		{"discrete.alg.merge_sort", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave27Fact(t, p, "n")
			if kp == 1 {
				return w27CeilLog2(n)
			}
			if kp == 2 {
				return n - 1
			}
			return n
		}},
		{"discrete.alg.binary_search", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave27Fact(t, p, "n")
			if kp == 3 {
				return n - 1
			}
			return w27CeilLog2(n)
		}},
		{"discrete.alg.dijkstra", func(t *testing.T, p generator.Problem, kp int) int {
			n, m := wave27Fact(t, p, "n"), wave27Fact(t, p, "m")
			if kp == 2 {
				return m
			}
			return n
		}},
		{"discrete.alg.huffman", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave27Fact(t, p, "n")
			if kp == 1 {
				return n
			}
			if kp == 2 {
				return n - 1
			}
			return 2*n - 1
		}},
		{"discrete.alg.backtracking", func(t *testing.T, p generator.Problem, kp int) int {
			n, k := wave27Fact(t, p, "n"), wave27Fact(t, p, "k")
			if kp == 1 {
				return n
			}
			if kp == 2 {
				return w27Pow(k, n)
			}
			return k
		}},
		{"discrete.alg.branch_bound", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave27Fact(t, p, "n")
			if kp == 3 {
				return 1
			}
			return n
		}},
		{"discrete.coding.entropy", func(t *testing.T, p generator.Problem, kp int) int {
			k := wave27Fact(t, p, "k")
			if kp == 1 {
				return k
			}
			if kp == 2 {
				return w27CeilLog2(k)
			}
			return 1
		}},
		{"discrete.coding.convolutional", func(t *testing.T, p generator.Problem, kp int) int {
			k := wave27Fact(t, p, "k")
			if kp == 3 {
				return w27Pow(2, k)
			}
			return k
		}},
		{"discrete.coding.ldpc", func(t *testing.T, p generator.Problem, kp int) int {
			n, m := wave27Fact(t, p, "n"), wave27Fact(t, p, "m")
			if kp == 2 {
				return n - m
			}
			return m
		}},
		{"discrete.coding.bch", func(t *testing.T, p generator.Problem, kp int) int {
			tt, n := wave27Fact(t, p, "t"), wave27Fact(t, p, "n")
			if kp == 1 {
				return tt
			}
			if kp == 2 {
				return 2 * tt
			}
			return n
		}},
		{"discrete.coding.singleton", func(t *testing.T, p generator.Problem, kp int) int {
			n, k := wave27Fact(t, p, "n"), wave27Fact(t, p, "k")
			if kp == 1 {
				return n - k + 1
			}
			if kp == 2 {
				return n - k
			}
			return 1
		}},
		{"discrete.complexity.space", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave27Fact(t, p, "n")
			if kp == 2 {
				return w27CeilLog2(n)
			}
			return n
		}},
		{"discrete.complexity.hierarchy", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 2 {
				return 2
			}
			return 1
		}},
		{"discrete.complexity.conp", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave27Fact(t, p, "n")
			}
			if kp == 2 {
				return 1
			}
			return 2
		}},
		{"discrete.complexity.pspace", func(t *testing.T, p generator.Problem, kp int) int {
			k := wave27Fact(t, p, "k")
			if kp == 1 {
				return k
			}
			if kp == 2 {
				return 2
			}
			return k + 1
		}},
		{"discrete.complexity.approximation", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave27Fact(t, p, "r")
			}
			if kp == 2 {
				return 2
			}
			return wave27Fact(t, p, "n")
		}},
		{"discrete.complexity.fpt", func(t *testing.T, p generator.Problem, kp int) int {
			k := wave27Fact(t, p, "k")
			if kp == 3 {
				return 1
			}
			return k
		}},
		{"discrete.seq.homogeneous", func(t *testing.T, p generator.Problem, kp int) int {
			return wave27Fact(t, p, "k")
		}},
		{"discrete.seq.characteristic", func(t *testing.T, p generator.Problem, kp int) int {
			r, m := wave27Fact(t, p, "r"), wave27Fact(t, p, "m")
			if kp == 1 {
				return r
			}
			if kp == 2 {
				return m - 1
			}
			return r
		}},
		{"discrete.seq.nonhomogeneous", func(t *testing.T, p generator.Problem, kp int) int {
			k, d := wave27Fact(t, p, "k"), wave27Fact(t, p, "d")
			if kp == 1 {
				return 1
			}
			if kp == 2 {
				return k
			}
			return d + 1
		}},
		{"discrete.seq.zeckendorf", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave27Fact(t, p, "n")
			if kp == 2 {
				return 0
			}
			if kp == 3 {
				return 1
			}
			return n
		}},
		{"discrete.seq.continued_fraction", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave27Fact(t, p, "n")
			if kp == 3 {
				return 1
			}
			return n
		}},
		{"discrete.seq.ackermann", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return 2
			}
			if kp == 2 {
				return 2
			}
			return 0
		}},
		{"discrete.sets.cardinality", func(t *testing.T, p generator.Problem, kp int) int {
			n, m := wave27Fact(t, p, "n"), wave27Fact(t, p, "m")
			if kp == 1 {
				return w27Pow(2, n)
			}
			if kp == 2 {
				return n * m
			}
			return n
		}},
		{"discrete.sets.cartesian", func(t *testing.T, p generator.Problem, kp int) int {
			n, m := wave27Fact(t, p, "n"), wave27Fact(t, p, "m")
			if kp == 1 {
				return n * m
			}
			if kp == 2 {
				return 2
			}
			return w27Pow(2, n*m)
		}},
		{"discrete.sets.equivalence_classes", func(t *testing.T, p generator.Problem, kp int) int {
			k := wave27Fact(t, p, "k")
			if kp == 3 {
				return 1
			}
			return k
		}},
		{"discrete.sets.lattice", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave27Fact(t, p, "n")
			if kp == 1 {
				return n
			}
			return 2
		}},
		{"discrete.proof.invariant", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave27Fact(t, p, "n")
			if kp == 1 {
				return 2 * n
			}
			if kp == 2 {
				return n + 1
			}
			return n
		}},
		{"discrete.proof.diagonal", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave27Fact(t, p, "n")
			if kp == 3 {
				return 1
			}
			return n
		}},
		{"discrete.proof.descent", func(t *testing.T, p generator.Problem, kp int) int {
			k := wave27Fact(t, p, "k")
			if kp == 2 {
				return 0
			}
			return k
		}},
		{"discrete.proof.constructive", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave27Fact(t, p, "n")
			if kp == 3 {
				return 1
			}
			return n
		}},
		{"discrete.proof.nonconstructive", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return 0
			}
			return 1
		}},
	}

	if len(cases) != len(wave27IDs()) {
		t.Fatalf("case table has %d entries, want %d", len(cases), len(wave27IDs()))
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
