package precalculus

import (
	"strconv"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func wave25IDs() []string {
	return []string{
		"precalc.func.domain", "precalc.func.even_odd", "precalc.func.monotonic",
		"precalc.func.extrema", "precalc.func.end_behavior", "precalc.func.period",
		"precalc.func.inverse_graph", "precalc.func.sign_chart",
		"precalc.func.absolute",
		"precalc.exp.half_life", "precalc.exp.fitting", "precalc.exp.base_change",
		"precalc.exp.log_properties", "precalc.exp.richter", "precalc.exp.ph",
		"precalc.exp.effective_rate", "precalc.exp.inequality",
		"precalc.analytic.midpoint", "precalc.analytic.slope",
		"precalc.analytic.parallel", "precalc.analytic.line_forms",
		"precalc.analytic.point_line", "precalc.analytic.circle_complete",
		"precalc.trig.radians", "precalc.trig.unit_circle",
		"precalc.trig.law_sines", "precalc.trig.law_cosines",
		"precalc.trig.identities", "precalc.trig.inverse",
		"precalc.trig.graphs", "precalc.trig.harmonic",
		"precalc.seq.telescoping", "precalc.seq.recurrence",
		"precalc.seq.fibonacci", "precalc.seq.monotone",
		"precalc.seq.geometric_mean", "precalc.seq.annuity",
		"precalc.seq.amortization", "precalc.seq.series_bound",
	}
}

func wave25Registry(t *testing.T) *generator.Registry {
	t.Helper()
	reg := generator.NewRegistry()
	Register(reg)
	return reg
}

func wave25Fact(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer\nfacts: %v", p.Question, key, p.Facts[key], p.Facts)
	}
	return v
}

func TestWave25AnswersAreIntegers(t *testing.T) {
	reg := wave25Registry(t)
	for _, id := range wave25IDs() {
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

func TestWave25QuestionVariety(t *testing.T) {
	reg := wave25Registry(t)
	for _, id := range wave25IDs() {
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

type wave25Case struct {
	id   string
	want func(t *testing.T, p generator.Problem, kp int) int
}

func TestWave25AnswersRecompute(t *testing.T) {
	reg := wave25Registry(t)
	cases := []wave25Case{
		{"precalc.func.domain", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave25Fact(t, p, "n")
			if kp == 1 {
				return 2
			}
			if kp == 2 {
				return n
			}
			return n + 1
		}},
		{"precalc.func.even_odd", func(t *testing.T, p generator.Problem, kp int) int {
			n, m := wave25Fact(t, p, "n"), wave25Fact(t, p, "m")
			if kp == 1 {
				return 1 - n%2
			}
			return 1 - (n+m)%2
		}},
		{"precalc.func.monotonic", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave25Fact(t, p, "n")
			if kp == 1 {
				return n + 1
			}
			if kp == 2 {
				return 1
			}
			return n
		}},
		{"precalc.func.extrema", func(t *testing.T, p generator.Problem, kp int) int {
			return wave25Fact(t, p, "n") - 1
		}},
		{"precalc.func.end_behavior", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave25Fact(t, p, "n")
			if kp == 1 {
				return (n + 1) % 2
			}
			if kp == 2 {
				return 1
			}
			return 2
		}},
		{"precalc.func.period", func(t *testing.T, p generator.Problem, kp int) int {
			return wave25Fact(t, p, "b")
		}},
		{"precalc.func.inverse_graph", func(t *testing.T, p generator.Problem, kp int) int {
			a, b := wave25Fact(t, p, "a"), wave25Fact(t, p, "b")
			if kp == 1 {
				return b
			}
			if kp == 2 {
				return a
			}
			return wave25Fact(t, p, "n")
		}},
		{"precalc.func.sign_chart", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave25Fact(t, p, "n")
			if kp == 1 {
				return n + 1
			}
			if kp == 2 {
				return 2
			}
			return n
		}},
		{"precalc.func.absolute", func(t *testing.T, p generator.Problem, kp int) int {
			a, n := wave25Fact(t, p, "a"), wave25Fact(t, p, "n")
			if kp == 1 {
				return 2
			}
			if kp == 2 {
				return a
			}
			return n
		}},
		{"precalc.exp.half_life", func(t *testing.T, p generator.Problem, kp int) int {
			k := wave25Fact(t, p, "k")
			if kp == 1 {
				return w25Pow(2, k)
			}
			if kp == 2 {
				return k
			}
			return 1
		}},
		{"precalc.exp.fitting", func(t *testing.T, p generator.Problem, kp int) int {
			k := wave25Fact(t, p, "k")
			if kp == 3 {
				return w25Pow(2, k)
			}
			return k
		}},
		{"precalc.exp.base_change", func(t *testing.T, p generator.Problem, kp int) int {
			b, k := wave25Fact(t, p, "b"), wave25Fact(t, p, "k")
			if kp == 2 {
				return b
			}
			return k
		}},
		{"precalc.exp.log_properties", func(t *testing.T, p generator.Problem, kp int) int {
			n, k := wave25Fact(t, p, "n"), wave25Fact(t, p, "k")
			if kp == 1 {
				return n
			}
			return k
		}},
		{"precalc.exp.richter", func(t *testing.T, p generator.Problem, kp int) int {
			m, d := wave25Fact(t, p, "m"), wave25Fact(t, p, "d")
			if kp == 1 {
				return m - wave25Fact(t, p, "r")
			}
			if kp == 3 {
				return d
			}
			return m
		}},
		{"precalc.exp.ph", func(t *testing.T, p generator.Problem, kp int) int {
			pv := wave25Fact(t, p, "p")
			if kp == 3 {
				return wave25Fact(t, p, "d")
			}
			return pv
		}},
		{"precalc.exp.effective_rate", func(t *testing.T, p generator.Problem, kp int) int {
			return wave25Fact(t, p, "n")
		}},
		{"precalc.exp.inequality", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return 0
			}
			if kp == 2 {
				return 1
			}
			return 1
		}},
		{"precalc.analytic.midpoint", func(t *testing.T, p generator.Problem, kp int) int {
			a, b := wave25Fact(t, p, "a"), wave25Fact(t, p, "b")
			if kp == 2 {
				return 2
			}
			return a + b
		}},
		{"precalc.analytic.slope", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave25Fact(t, p, "rise")
			}
			if kp == 2 {
				return wave25Fact(t, p, "run")
			}
			return 0
		}},
		{"precalc.analytic.parallel", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 3 {
				return 0
			}
			if kp == 1 {
				return -1
			}
			return -1
		}},
		{"precalc.analytic.line_forms", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return 2
			}
			return wave25Fact(t, p, "b")
		}},
		{"precalc.analytic.point_line", func(t *testing.T, p generator.Problem, kp int) int {
			a, b, d := wave25Fact(t, p, "a"), wave25Fact(t, p, "b"), wave25Fact(t, p, "d")
			if kp == 1 {
				return d
			}
			return a*a + b*b
		}},
		{"precalc.analytic.circle_complete", func(t *testing.T, p generator.Problem, kp int) int {
			h, k, r := wave25Fact(t, p, "h"), wave25Fact(t, p, "k"), wave25Fact(t, p, "r")
			if kp == 1 {
				return h + k
			}
			if kp == 2 {
				return r * r
			}
			return r
		}},
		{"precalc.trig.radians", func(t *testing.T, p generator.Problem, kp int) int {
			d := wave25Fact(t, p, "d")
			if kp == 1 {
				return d
			}
			if kp == 2 {
				return 180
			}
			return d / 15
		}},
		{"precalc.trig.unit_circle", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return 1
			}
			return 2
		}},
		{"precalc.trig.law_sines", func(t *testing.T, p generator.Problem, kp int) int {
			r := wave25Fact(t, p, "r")
			if kp == 2 {
				return 1
			}
			return 2 * r
		}},
		{"precalc.trig.law_cosines", func(t *testing.T, p generator.Problem, kp int) int {
			a, b := wave25Fact(t, p, "a"), wave25Fact(t, p, "b")
			if kp == 1 {
				return a*a + b*b
			}
			if kp == 2 {
				return 2
			}
			return 0
		}},
		{"precalc.trig.identities", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return 1
			}
			if kp == 2 {
				return 3
			}
			return 0
		}},
		{"precalc.trig.inverse", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return 2
			}
			if kp == 2 {
				return 1
			}
			return 2
		}},
		{"precalc.trig.graphs", func(t *testing.T, p generator.Problem, kp int) int {
			a, b, d := wave25Fact(t, p, "a"), wave25Fact(t, p, "b"), wave25Fact(t, p, "d")
			if kp == 1 {
				return a
			}
			if kp == 2 {
				return b
			}
			return d
		}},
		{"precalc.trig.harmonic", func(t *testing.T, p generator.Problem, kp int) int {
			w := wave25Fact(t, p, "w")
			if kp == 2 {
				return 2
			}
			return w
		}},
		{"precalc.seq.telescoping", func(t *testing.T, p generator.Problem, kp int) int {
			n := wave25Fact(t, p, "n")
			if kp == 1 {
				return n + 1
			}
			if kp == 2 {
				return n
			}
			return 2
		}},
		{"precalc.seq.recurrence", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return wave25Fact(t, p, "r")
			}
			if kp == 2 {
				return wave25Fact(t, p, "c")
			}
			return wave25Fact(t, p, "n")
		}},
		{"precalc.seq.fibonacci", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 2 {
				return 1
			}
			return 2
		}},
		{"precalc.seq.monotone", func(t *testing.T, p generator.Problem, kp int) int {
			if kp == 1 {
				return 1
			}
			if kp == 2 {
				return 1
			}
			return 0
		}},
		{"precalc.seq.geometric_mean", func(t *testing.T, p generator.Problem, kp int) int {
			a, b := wave25Fact(t, p, "a"), wave25Fact(t, p, "b")
			if kp == 1 {
				return a * a * b * b
			}
			if kp == 2 {
				return a * b
			}
			return wave25Fact(t, p, "n")
		}},
		{"precalc.seq.annuity", func(t *testing.T, p generator.Problem, kp int) int {
			return wave25Fact(t, p, "n")
		}},
		{"precalc.seq.amortization", func(t *testing.T, p generator.Problem, kp int) int {
			n, m := wave25Fact(t, p, "n"), wave25Fact(t, p, "m")
			if kp == 1 {
				return n
			}
			if kp == 2 {
				return 12 * n
			}
			return m * n
		}},
		{"precalc.seq.series_bound", func(t *testing.T, p generator.Problem, kp int) int {
			n, r := wave25Fact(t, p, "n"), wave25Fact(t, p, "r")
			if kp == 2 {
				return r - 1
			}
			return n
		}},
	}

	if len(cases) != len(wave25IDs()) {
		t.Fatalf("case table has %d entries, want %d", len(cases), len(wave25IDs()))
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
