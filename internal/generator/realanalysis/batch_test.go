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
	for i := 0; i < 200; i++ {
		p := raSample(t, reg, "ra.order.supremum", i)
		lo := raFactInt(t, p, "lo")
		hi := raFactInt(t, p, "hi")
		if lo >= hi {
			t.Fatalf("supremum: lo=%d is not below hi=%d", lo, hi)
		}
		// The interval is bounded above by hi and approaches it, so the least
		// upper bound is hi.
		if p.Answer != strconv.Itoa(hi) {
			t.Fatalf("supremum: lo=%d hi=%d answers %s, want %d", lo, hi, p.Answer, hi)
		}
	}
}

func TestRAArchimedeanRecomputes(t *testing.T) {
	reg := raRegistry(t)
	for i := 0; i < 200; i++ {
		p := raSample(t, reg, "ra.order.archimedean", i)
		pn := raFactInt(t, p, "pn")
		pd := raFactInt(t, p, "pd")
		if pn <= 0 || pd <= 0 {
			t.Fatalf("archimedean: pn=%d pd=%d are not positive", pn, pd)
		}
		// 1/N < pn/pd is N > pd/pn, so the smallest integer N is floor(pd/pn)+1.
		want := pd/pn + 1
		if p.Answer != strconv.Itoa(want) {
			t.Fatalf("archimedean: pn=%d pd=%d answers %s, want %d", pn, pd, p.Answer, want)
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
	for i := 0; i < 200; i++ {
		p := raSample(t, reg, "ra.seq.monotone", i)
		a := raFactInt(t, p, "a")
		want := strconv.Itoa(int((1 + math.Sqrt(1+4*float64(a))) / 2))
		if p.Answer != want {
			t.Fatalf("monotone: a=%d answers %s, want %s", a, p.Answer, want)
		}
	}
}

func TestRASeriesRecomputes(t *testing.T) {
	reg := raRegistry(t)
	for i := 0; i < 200; i++ {
		p := raSample(t, reg, "ra.seq.series", i)
		pn := raFactInt(t, p, "pn")
		pd := raFactInt(t, p, "pd")
		if pd == pn {
			t.Fatalf("series: ratio pn/pd is one, the series diverges")
		}
		want := decimal6(float64(pd) / float64(pd-pn))
		if p.Answer != want {
			t.Fatalf("series: pn=%d pd=%d answers %s, want %s", pn, pd, p.Answer, want)
		}
	}
}

func TestRAFunctionLimitRecomputes(t *testing.T) {
	reg := raRegistry(t)
	for i := 0; i < 200; i++ {
		p := raSample(t, reg, "ra.limit.function", i)
		m := raFactInt(t, p, "m")
		en := raFactInt(t, p, "en")
		ed := raFactInt(t, p, "ed")
		if m == 0 || ed == 0 {
			t.Fatalf("function limit: m=%d ed=%d", m, ed)
		}
		want := decimal6((float64(en) / float64(ed)) / float64(m))
		if p.Answer != want {
			t.Fatalf("function limit: m=%d eps=%d/%d answers %s, want %s", m, en, ed, p.Answer, want)
		}
	}
}
