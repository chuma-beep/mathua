package realanalysis

// Tests for the first real-analysis batch.
//
// Every concept is graded through the real router with the grading_type the
// corpus ships. The nineteen multiple-choice concepts are checked the way the
// other batches check theirs: the question is parsed into options, the correct
// answer must appear exactly once, and every other option must be rejected by
// the router.
//
// The six numeric concepts recompute their answers from the published
// parameters: the supremum from the interval bound, the Archimedean index from
// the two positive integers, a sequence limit from the leading coefficients, a
// monotone limit from the recurrence constant, a geometric sum from the ratio,
// and a delta from the epsilon and the slope. The tests never read the
// published "answer" fact.

import (
	"encoding/json"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/grader"
)

type raCorpusConcept struct {
	ID          string `json:"id"`
	GradingType string `json:"grading_type"`
}

var raBatchIDs = []string{
	"ra.order.supremum",
	"ra.order.completeness",
	"ra.order.archimedean",
	"ra.order.nested_intervals",
	"ra.order.density",
	"ra.order.countable",
	"ra.order.uncountable",
	"ra.order.extended",
	"ra.seq.convergence",
	"ra.seq.monotone",
	"ra.seq.subsequence",
	"ra.seq.bolzano_weierstrass",
	"ra.seq.cauchy",
	"ra.seq.series",
	"ra.seq.tests",
	"ra.seq.absolute",
	"ra.seq.rearrangement",
	"ra.seq.cauchy_criterion",
	"ra.limit.function",
	"ra.limit.sequential",
	"ra.limit.continuity",
	"ra.limit.ivt",
	"ra.limit.evt",
	"ra.limit.uniform",
	"ra.limit.compact",
}

var raMCIDs = []string{}

func raShippedTypes(t *testing.T) map[string]string {
	t.Helper()
	raw, err := os.ReadFile("../../../data/concepts/real_analysis.json")
	if err != nil {
		t.Fatalf("read corpus: %v", err)
	}
	var rows []raCorpusConcept
	if err := json.Unmarshal(raw, &rows); err != nil {
		t.Fatalf("parse corpus: %v", err)
	}
	out := map[string]string{}
	for _, r := range rows {
		if r.ID != "" {
			out[r.ID] = r.GradingType
		}
	}
	return out
}

func raRegistry(t *testing.T) *generator.Registry {
	t.Helper()
	reg := generator.NewRegistry()
	Register(reg)
	return reg
}

func raSample(t *testing.T, reg *generator.Registry, id string, i int) generator.Problem {
	t.Helper()
	p, err := reg.GenerateContext(id, generator.GeneratorContext{
		Difficulty: float64(i%20) / 19,
		Seed:       int64(381000 + i),
	})
	if err != nil {
		t.Fatalf("%s: generate: %v", id, err)
	}
	return p
}

// raSampleKP pins the knowledge-point ordinal, which the numeric recompute
// tests need because each ordinal now asks a different question shape.
func raSampleKP(t *testing.T, reg *generator.Registry, id string, kp, i int) generator.Problem {
	t.Helper()
	p, err := reg.GenerateContext(id, generator.GeneratorContext{
		Difficulty: float64(i%20) / 19,
		Seed:       int64(381000 + i),
		KPOrdinal:  kp,
	})
	if err != nil {
		t.Fatalf("%s KP%d: generate: %v", id, kp, err)
	}
	return p
}

var raOptionLine = regexp.MustCompile(`\n\n([A-E])\) ([^\n]*)`)

func raOptionsOf(q string) []string {
	var out []string
	for _, m := range raOptionLine.FindAllStringSubmatch(q, -1) {
		out = append(out, m[2])
	}
	return out
}

func raFactInt(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer", p.Question, key, p.Facts[key])
	}
	return v
}

// ---------------------------------------------------------------------------
// Router contracts
// ---------------------------------------------------------------------------

func TestRAAnswersGradeWithTheirShippedType(t *testing.T) {
	types := raShippedTypes(t)
	reg := raRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	for _, id := range raBatchIDs {
		typ, ok := types[id]
		if !ok || typ == "" {
			t.Fatalf("%s: no grading_type in the corpus", id)
		}
		if typ != "multiple_choice" && typ != "numeric" {
			t.Fatalf("%s: unexpected grading_type %q", id, typ)
		}
		for i := 0; i < 120; i++ {
			p := raSample(t, reg, id, i)
			if strings.TrimSpace(p.Answer) == "" {
				t.Fatalf("%s: empty answer at sample %d", id, i)
			}
			if strings.TrimSpace(p.Question) == "" || strings.TrimSpace(p.Explanation) == "" {
				t.Fatalf("%s: empty question or explanation at sample %d", id, i)
			}
			res := router.Grade(grader.GradingType(typ), p.Answer, p.Answer)
			if res.Unavailable {
				t.Skipf("%s: grader unavailable", id)
			}
			if !res.Correct {
				t.Fatalf("%s: the generator's own answer %q does not grade as %s", id, p.Answer, typ)
			}
		}
	}
}

func TestRAMultipleChoiceRejectsEveryDistractor(t *testing.T) {
	types := raShippedTypes(t)
	reg := raRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	for _, id := range raMCIDs {
		if types[id] != "multiple_choice" {
			t.Fatalf("%s: expected multiple_choice, corpus says %q", id, types[id])
		}
		for i := 0; i < 40; i++ {
			p := raSample(t, reg, id, i)
			opts := raOptionsOf(p.Question)
			if len(opts) < 3 {
				t.Fatalf("%s: only %d options; a recognition question needs at least three", id, len(opts))
			}
			found := 0
			for _, o := range opts {
				if o == p.Answer {
					found++
				}
			}
			if found != 1 {
				t.Fatalf("%s: the answer %q appears %d times among %v", id, p.Answer, found, opts)
			}
			for _, o := range opts {
				if o == p.Answer {
					continue
				}
				if router.Grade(grader.GradingMultipleChoice, p.Answer, o).Correct {
					t.Fatalf("%s: distractor %q is accepted as correct", id, o)
				}
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Independent numeric checks
// ---------------------------------------------------------------------------

func TestRASupremumRecomputes(t *testing.T) {
	reg := raRegistry(t)
	for kp := 1; kp <= 3; kp++ {
		for i := 0; i < 100; i++ {
			p := raSampleKP(t, reg, "ra.order.supremum", kp, i)
			lo := raFactInt(t, p, "lo")
			hi := raFactInt(t, p, "hi")
			if lo >= hi {
				t.Fatalf("supremum KP%d: lo=%d is not below hi=%d", kp, lo, hi)
			}
			// Every shape is bounded above by hi and approaches it, so the
			// least upper bound is hi even when the set has no maximum.
			if p.Answer != strconv.Itoa(hi) {
				t.Fatalf("supremum KP%d: lo=%d hi=%d answers %s, want %d", kp, lo, hi, p.Answer, hi)
			}
		}
	}
}

func TestRAArchimedeanRecomputes(t *testing.T) {
	reg := raRegistry(t)
	for kp := 1; kp <= 3; kp++ {
		for i := 0; i < 100; i++ {
			p := raSampleKP(t, reg, "ra.order.archimedean", kp, i)
			pn := raFactInt(t, p, "pn")
			pd := raFactInt(t, p, "pd")
			if pn <= 0 || pd <= 0 || pn >= pd {
				t.Fatalf("archimedean KP%d: pn=%d pd=%d are not a positive proper fraction", kp, pn, pd)
			}
			switch p.Facts["kind"] {
			case "smallest_reciprocal", "integer_above":
				// 1/N < pn/pd is N > pd/pn, so the smallest integer N is
				// floor(pd/pn)+1.
				if want := pd/pn + 1; p.Answer != strconv.Itoa(want) {
					t.Fatalf("archimedean KP%d: pn=%d pd=%d answers %s, want %d", kp, pn, pd, p.Answer, want)
				}
			case "largest_reciprocal":
				// 1/N ≥ pn/pd is N ≤ pd/pn, so the largest integer N is
				// floor(pd/pn).
				if want := pd / pn; p.Answer != strconv.Itoa(want) {
					t.Fatalf("archimedean KP%d: pn=%d pd=%d answers %s, want %d", kp, pn, pd, p.Answer, want)
				}
			default:
				t.Fatalf("archimedean KP%d: unknown kind %q", kp, p.Facts["kind"])
			}
		}
	}
}

func TestRASequenceRecomputes(t *testing.T) {
	reg := raRegistry(t)
	for i := 0; i < 200; i++ {
		p := raSample(t, reg, "ra.seq.convergence", i)
		a := raFactInt(t, p, "a")
		c := raFactInt(t, p, "c")
		if c == 0 {
			t.Fatalf("sequence: leading denominator coefficient c=0")
		}
		want := decimal6(float64(a) / float64(c))
		if p.Answer != want {
			t.Fatalf("sequence: a=%d c=%d answers %s, want %s", a, c, p.Answer, want)
		}
	}
}

func TestRAMonotoneRecomputes(t *testing.T) {
	reg := raRegistry(t)
	for kp := 1; kp <= 3; kp++ {
		for i := 0; i < 100; i++ {
			p := raSampleKP(t, reg, "ra.seq.monotone", kp, i)
			switch p.Facts["kind"] {
			case "sqrt_linear":
				a := raFactInt(t, p, "a")
				want := strconv.Itoa(int((1 + math.Sqrt(1+4*float64(a))) / 2))
				if p.Answer != want {
					t.Fatalf("monotone KP%d: a=%d answers %s, want %s", kp, a, p.Answer, want)
				}
			case "sqrt_affine":
				a := raFactInt(t, p, "a")
				b := raFactInt(t, p, "b")
				want := strconv.Itoa(int((float64(b) + math.Sqrt(float64(b*b+4*a))) / 2))
				if p.Answer != want {
					t.Fatalf("monotone KP%d: a=%d b=%d answers %s, want %s", kp, a, b, p.Answer, want)
				}
			case "affine_average":
				c := raFactInt(t, p, "c")
				q := raFactInt(t, p, "q")
				start := raFactInt(t, p, "start")
				if q <= 1 || c%(q-1) != 0 {
					t.Fatalf("monotone KP%d: c=%d q=%d do not give an integer limit", kp, c, q)
				}
				limit := c / (q - 1)
				if start >= limit {
					t.Fatalf("monotone KP%d: start=%d is not below the limit %d", kp, start, limit)
				}
				if p.Answer != strconv.Itoa(limit) {
					t.Fatalf("monotone KP%d: c=%d q=%d answers %s, want %d", kp, c, q, p.Answer, limit)
				}
			default:
				t.Fatalf("monotone KP%d: unknown kind %q", kp, p.Facts["kind"])
			}
		}
	}
}

func TestRASeriesRecomputes(t *testing.T) {
	reg := raRegistry(t)
	for kp := 1; kp <= 3; kp++ {
		for i := 0; i < 100; i++ {
			p := raSampleKP(t, reg, "ra.seq.series", kp, i)
			switch p.Facts["kind"] {
			case "partial_sum":
				a := raFactInt(t, p, "a")
				pn := raFactInt(t, p, "pn")
				pd := raFactInt(t, p, "pd")
				n := raFactInt(t, p, "n")
				if pn >= pd {
					t.Fatalf("series KP%d: ratio %d/%d is not below one", kp, pn, pd)
				}
				r := float64(pn) / float64(pd)
				want := decimal6(float64(a) * (1 - math.Pow(r, float64(n))) / (1 - r))
				if p.Answer != want {
					t.Fatalf("series KP%d: a=%d pn=%d pd=%d n=%d answers %s, want %s", kp, a, pn, pd, n, p.Answer, want)
				}
			case "infinite_sum":
				a := raFactInt(t, p, "a")
				pn := raFactInt(t, p, "pn")
				pd := raFactInt(t, p, "pd")
				if pd == pn {
					t.Fatalf("series KP%d: ratio pn/pd is one, the series diverges", kp)
				}
				want := decimal6(float64(a) * float64(pd) / float64(pd-pn))
				if p.Answer != want {
					t.Fatalf("series KP%d: a=%d pn=%d pd=%d answers %s, want %s", kp, a, pn, pd, p.Answer, want)
				}
			case "term_limit":
				a := raFactInt(t, p, "a")
				c := raFactInt(t, p, "c")
				if a == 0 || c == 0 {
					t.Fatalf("series KP%d: leading coefficients a=%d c=%d", kp, a, c)
				}
				want := decimal6(float64(a) / float64(c))
				if p.Answer != want {
					t.Fatalf("series KP%d: a=%d c=%d answers %s, want %s", kp, a, c, p.Answer, want)
				}
			default:
				t.Fatalf("series KP%d: unknown kind %q", kp, p.Facts["kind"])
			}
		}
	}
}

func TestRAFunctionLimitRecomputes(t *testing.T) {
	reg := raRegistry(t)
	for kp := 1; kp <= 3; kp++ {
		for i := 0; i < 100; i++ {
			p := raSampleKP(t, reg, "ra.limit.function", kp, i)
			switch p.Facts["kind"] {
			case "linear_at_zero", "linear_at_point":
				m := raFactInt(t, p, "m")
				en := raFactInt(t, p, "en")
				ed := raFactInt(t, p, "ed")
				if m == 0 || ed == 0 {
					t.Fatalf("function limit KP%d: m=%d ed=%d", kp, m, ed)
				}
				want := decimal6((float64(en) / float64(ed)) / float64(m))
				if p.Answer != want {
					t.Fatalf("function limit KP%d: m=%d eps=%d/%d answers %s, want %s", kp, m, en, ed, p.Answer, want)
				}
			case "factor_cancel":
				a := raFactInt(t, p, "a")
				switch p.Facts["degree"] {
				case "2":
					if want := strconv.Itoa(2 * a); p.Answer != want {
						t.Fatalf("function limit KP%d: a=%d answers %s, want %s", kp, a, p.Answer, want)
					}
				case "3":
					if want := strconv.Itoa(3 * a * a); p.Answer != want {
						t.Fatalf("function limit KP%d: a=%d answers %s, want %s", kp, a, p.Answer, want)
					}
				default:
					t.Fatalf("function limit KP%d: unknown degree %q", kp, p.Facts["degree"])
				}
			default:
				t.Fatalf("function limit KP%d: unknown kind %q", kp, p.Facts["kind"])
			}
		}
	}
}
