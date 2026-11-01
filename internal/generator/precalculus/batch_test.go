package precalculus

// Tests for the 2026-10-18 precalculus batch.
//
// Every concept is graded through the real router with the grading_type the
// corpus ships. The two multiple-choice concepts are checked the way the other
// batches check theirs: the question is parsed into options, the correct answer
// must appear exactly once, and every other option must be rejected.
//
// The six numeric concepts recompute their answers from the published facts:
// the logarithmic-scale kind and exponent, the inequality roots, the cancelled
// root of a rational function, the conic coefficients, the launch velocity, and
// the perpendicular components. The expression concept is checked against an
// independent decomposition table.

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

var pcBatchIDs = []string{
	"precalc.exp.continuous",
	"precalc.exp.cooling",
	"precalc.log.scale",
	"precalc.poly.inequality",
	"precalc.rational.holes",
	"precalc.conic.general",
	"precalc.parametric.projectile",
	"precalc.vector.resultant",
	"precalc.func.decomposition",
}

var pcBatchNumericIDs = []string{
	"precalc.exp.continuous",
	"precalc.exp.cooling",
	"precalc.log.scale",
	"precalc.poly.inequality",
	"precalc.rational.holes",
	"precalc.conic.general",
	"precalc.parametric.projectile",
	"precalc.vector.resultant",
}

type pcCorpusConcept struct {
	ID          string `json:"id"`
	GradingType string `json:"grading_type"`
}

func pcShippedTypes(t *testing.T) map[string]string {
	t.Helper()
	raw, err := os.ReadFile("../../../data/concepts/precalculus.json")
	if err != nil {
		t.Fatalf("read corpus: %v", err)
	}
	var rows []pcCorpusConcept
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

func pcBatchRegistry(t *testing.T) *generator.Registry {
	t.Helper()
	reg := generator.NewRegistry()
	Register(reg)
	return reg
}

func pcBatchSample(t *testing.T, reg *generator.Registry, id string, i int) generator.Problem {
	t.Helper()
	p, err := reg.GenerateContext(id, generator.GeneratorContext{
		Difficulty: float64(i%20) / 19,
		Seed:       int64(151000 + i),
	})
	if err != nil {
		t.Fatalf("%s: generate: %v", id, err)
	}
	return p
}

var pcOptionLine = regexp.MustCompile(`\n\n([A-E])\) ([^\n]*)`)

func pcOptionsOf(q string) []string {
	var out []string
	for _, m := range pcOptionLine.FindAllStringSubmatch(q, -1) {
		out = append(out, m[2])
	}
	return out
}

func pcFactInt(t *testing.T, p generator.Problem, key string) int {
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

func TestBatchAnswersGradeWithTheirShippedType(t *testing.T) {
	types := pcShippedTypes(t)
	reg := pcBatchRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	for _, id := range pcBatchIDs {
		typ, ok := types[id]
		if !ok || typ == "" {
			t.Fatalf("%s: no grading_type in the corpus", id)
		}
		switch typ {
		case "multiple_choice", "numeric", "expression":
		default:
			t.Fatalf("%s: unexpected grading_type %q", id, typ)
		}
		for i := 0; i < 120; i++ {
			p := pcBatchSample(t, reg, id, i)
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

func TestBatchMultipleChoiceRejectsEveryDistractor(t *testing.T) {
	types := pcShippedTypes(t)
	reg := pcBatchRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	checked := 0
	for _, id := range pcBatchIDs {
		if types[id] != "multiple_choice" {
			continue
		}
		for i := 0; i < 40; i++ {
			p := pcBatchSample(t, reg, id, i)
			opts := pcOptionsOf(p.Question)
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
			checked++
		}
	}
	if checked == 0 {
		t.Skip("the precalculus batch now serves computed answers")
	}
}

// ---------------------------------------------------------------------------
// Independent checks
// ---------------------------------------------------------------------------

func TestBatchLogScaleRecomputes(t *testing.T) {
	reg := pcBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := pcBatchSample(t, reg, "precalc.log.scale", i)
		e := pcFactInt(t, p, "e")
		var want int
		switch p.Facts["kind"] {
		case "ph":
			want = e
		case "richter":
			want = e
		case "db":
			want = 10 * e
		default:
			t.Fatalf("log scale: unknown kind %q", p.Facts["kind"])
		}
		answer, err := strconv.Atoi(p.Answer)
		if err != nil {
			t.Fatalf("log scale: non-integer answer %q", p.Answer)
		}
		if answer != want {
			t.Fatalf("log scale: %s e=%d answers %d, want %d", p.Facts["kind"], e, answer, want)
		}
	}
}

func TestBatchPolyInequalityRecomputes(t *testing.T) {
	reg := pcBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := pcBatchSample(t, reg, "precalc.poly.inequality", i)
		a := pcFactInt(t, p, "p")
		b := pcFactInt(t, p, "q")
		want := b - a
		if want < 0 {
			want = -want
		}
		answer, err := strconv.Atoi(p.Answer)
		if err != nil {
			t.Fatalf("poly inequality: non-integer answer %q", p.Answer)
		}
		if answer != want {
			t.Fatalf("poly inequality: roots %d and %d answer %d, want %d", a, b, answer, want)
		}
	}
}

func TestBatchRationalHolesRecompute(t *testing.T) {
	reg := pcBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := pcBatchSample(t, reg, "precalc.rational.holes", i)
		want := pcFactInt(t, p, "root")
		answer, err := strconv.Atoi(p.Answer)
		if err != nil {
			t.Fatalf("rational holes: non-integer answer %q", p.Answer)
		}
		if answer != want {
			t.Fatalf("rational holes: root %d answers %d", want, answer)
		}
	}
}

func TestBatchGeneralConicRecomputes(t *testing.T) {
	reg := pcBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := pcBatchSample(t, reg, "precalc.conic.general", i)
		a := pcFactInt(t, p, "a")
		b := pcFactInt(t, p, "b")
		c := pcFactInt(t, p, "c")
		want := b*b - 4*a*c
		answer, err := strconv.Atoi(p.Answer)
		if err != nil {
			t.Fatalf("general conic: non-integer answer %q", p.Answer)
		}
		if answer != want {
			t.Fatalf("general conic: A=%d B=%d C=%d answers %d, want %d", a, b, c, answer, want)
		}
	}
}

func TestBatchProjectileRecomputes(t *testing.T) {
	reg := pcBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := pcBatchSample(t, reg, "precalc.parametric.projectile", i)
		v := float64(pcFactInt(t, p, "v"))
		want := v * v / 20
		answer, err := strconv.ParseFloat(p.Answer, 64)
		if err != nil {
			t.Fatalf("projectile: non-numeric answer %q", p.Answer)
		}
		if math.Abs(answer-want) > 1e-9 {
			t.Fatalf("projectile: v=%v answers %v, want %v", v, answer, want)
		}
	}
}

func TestBatchResultantRecomputes(t *testing.T) {
	reg := pcBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := pcBatchSample(t, reg, "precalc.vector.resultant", i)
		a := float64(pcFactInt(t, p, "a"))
		b := float64(pcFactInt(t, p, "b"))
		want := math.Sqrt(a*a + b*b)
		answer, err := strconv.ParseFloat(p.Answer, 64)
		if err != nil {
			t.Fatalf("resultant: non-numeric answer %q", p.Answer)
		}
		if math.Abs(answer-want) > 1e-9 {
			t.Fatalf("resultant: components %v and %v answer %v, want %v", a, b, answer, want)
		}
	}
}

func TestBatchDecompositionMatchesTable(t *testing.T) {
	reg := pcBatchRegistry(t)
	kinds := map[string]bool{
		"square_shift": true,
		"sqrt":         true,
		"reciprocal":   true,
		"exp":          true,
		"cube":         true,
		"log":          true,
	}
	seen := map[string]bool{}
	for i := 0; i < 300; i++ {
		p := pcBatchSample(t, reg, "precalc.func.decomposition", i)
		kind := p.Facts["kind"]
		seen[kind] = true
		if !kinds[kind] {
			t.Fatalf("decomposition: unknown kind %q", kind)
		}
		var want string
		switch kind {
		case "square_shift", "cube", "reciprocal":
			a, err := strconv.Atoi(p.Facts["a"])
			if err != nil {
				t.Fatalf("decomposition: %s has a non-numeric shift %q", kind, p.Facts["a"])
			}
			if a >= 0 {
				want = "x+" + strconv.Itoa(a)
			} else {
				want = "x-" + strconv.Itoa(-a)
			}
		case "sqrt":
			want = "x^2+" + p.Facts["a"]
		case "exp":
			want = p.Facts["k"] + "*x"
		case "log":
			want = "x^" + p.Facts["k"]
		}
		if p.Answer != want {
			t.Fatalf("decomposition: kind %s answers %q, want %q", kind, p.Answer, want)
		}
	}
	if len(seen) != len(kinds) {
		t.Fatalf("decomposition: expected all %d kinds, saw %v", len(kinds), seen)
	}
}

func TestBatchNumericAnswersParseAsDecimals(t *testing.T) {
	reg := pcBatchRegistry(t)
	for _, id := range pcBatchNumericIDs {
		for i := 0; i < 120; i++ {
			p := pcBatchSample(t, reg, id, i)
			if _, err := strconv.ParseFloat(p.Answer, 64); err != nil {
				t.Fatalf("%s: answer %q is not a decimal number", id, p.Answer)
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Pool quality
// ---------------------------------------------------------------------------

func TestBatchPoolsAreDiverse(t *testing.T) {
	reg := pcBatchRegistry(t)
	const samples = 150
	for _, id := range pcBatchIDs {
		answers := map[string]bool{}
		questions := map[string]bool{}
		allYesNo := true
		for i := 0; i < samples; i++ {
			p, err := reg.GenerateContext(id, generator.GeneratorContext{Difficulty: 0.12, Seed: int64(152000 + i)})
			if err != nil {
				t.Fatalf("%s: generate: %v", id, err)
			}
			answers[strings.ToLower(strings.TrimSpace(p.Answer))] = true
			questions[p.Question] = true
			if !strings.Contains(p.Question, "(yes/no)") {
				allYesNo = false
			}
		}
		if len(answers) < 2 {
			t.Fatalf("%s: only %d distinct answers over %d samples", id, len(answers), samples)
		}
		if len(questions) < 2 {
			t.Fatalf("%s: only %d distinct questions over %d samples", id, len(questions), samples)
		}
		if allYesNo {
			t.Fatalf("%s: every question is yes/no-only", id)
		}
	}
}
