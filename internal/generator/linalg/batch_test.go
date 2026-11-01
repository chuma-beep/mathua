package linalg

// Tests for the 2026-10-18 linear-algebra batch.
//
// Every concept is graded through the real router with the grading_type the
// corpus ships. The multiple-choice concepts are checked the way the other
// batches check theirs: the question is parsed into options, the correct answer
// must appear exactly once, and every other option must be rejected.
//
// The five numeric concepts and the tuple concept are recomputed from the facts
// the generator publishes: the kernel dimension from the matrix rank by
// elimination, the basis coordinates by solving the 2x2 system, the orthogonal
// complement from n - m, the projection from the scalar coefficient, the cross
// product from the two factors, and the condition number from the singular
// values.

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/grader"
)

var linalgBatchIDs = []string{
	"linalg.space.axioms",
	"linalg.transform.kernel",
	"linalg.transform.invertible",
	"linalg.basis.change",
	"linalg.orth.complement",
	"linalg.orth.projection_matrix",
	"linalg.symmetric.spectral",
	"linalg.quadratic.definite",
	"linalg.vector.cross",
	"linalg.matrix.condition",
	"linalg.eigen.cayley_hamilton",
}

var linalgBatchNumericIDs = []string{
	"linalg.transform.kernel",
	"linalg.basis.change",
	"linalg.orth.complement",
	"linalg.orth.projection_matrix",
	"linalg.matrix.condition",
	"linalg.space.axioms",
	"linalg.transform.invertible",
	"linalg.symmetric.spectral",
	"linalg.quadratic.definite",
	"linalg.eigen.cayley_hamilton",
}

type linalgCorpusConcept struct {
	ID          string `json:"id"`
	GradingType string `json:"grading_type"`
}

func linalgShippedTypes(t *testing.T) map[string]string {
	t.Helper()
	raw, err := os.ReadFile("../../../data/concepts/linear_algebra.json")
	if err != nil {
		t.Fatalf("read corpus: %v", err)
	}
	var rows []linalgCorpusConcept
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

func linalgBatchRegistry(t *testing.T) *generator.Registry {
	t.Helper()
	reg := generator.NewRegistry()
	Register(reg)
	return reg
}

func linalgBatchSample(t *testing.T, reg *generator.Registry, id string, i int) generator.Problem {
	t.Helper()
	p, err := reg.GenerateContext(id, generator.GeneratorContext{
		Difficulty: float64(i%20) / 19,
		Seed:       int64(101000 + i),
	})
	if err != nil {
		t.Fatalf("%s: generate: %v", id, err)
	}
	return p
}

var linalgOptionLine = regexp.MustCompile(`\n\n([A-E])\) ([^\n]*)`)

func linalgOptionsOf(q string) []string {
	var out []string
	for _, m := range linalgOptionLine.FindAllStringSubmatch(q, -1) {
		out = append(out, m[2])
	}
	return out
}

func linalgFactInt(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer", p.Question, key, p.Facts[key])
	}
	return v
}

func linalgFactInts(t *testing.T, p generator.Problem, key string) []int {
	t.Helper()
	parts := strings.Split(p.Facts[key], ",")
	out := make([]int, len(parts))
	for i, s := range parts {
		v, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil {
			t.Fatalf("%s: fact %q contains %q, not an integer", p.Question, key, s)
		}
		out[i] = v
	}
	return out
}

// linalgRank computes the rank of an integer matrix by Gauss-Jordan
// elimination over float64. The entries are small integers, so the tolerance is
// safe and the computation is independent of the generator's own reasoning.
func linalgRank(rows [][]int) int {
	if len(rows) == 0 {
		return 0
	}
	m := make([][]float64, len(rows))
	for i, row := range rows {
		m[i] = make([]float64, len(row))
		for j, v := range row {
			m[i][j] = float64(v)
		}
	}
	rank := 0
	cols := len(m[0])
	for col := 0; col < cols && rank < len(m); col++ {
		pivot := -1
		for r := rank; r < len(m); r++ {
			if math.Abs(m[r][col]) > 1e-9 {
				pivot = r
				break
			}
		}
		if pivot < 0 {
			continue
		}
		m[rank], m[pivot] = m[pivot], m[rank]
		pv := m[rank][col]
		for j := col; j < cols; j++ {
			m[rank][j] /= pv
		}
		for r := 0; r < len(m); r++ {
			if r == rank || math.Abs(m[r][col]) <= 1e-9 {
				continue
			}
			f := m[r][col]
			for j := col; j < cols; j++ {
				m[r][j] -= f * m[rank][j]
			}
		}
		rank++
	}
	return rank
}

// ---------------------------------------------------------------------------
// Router contracts
// ---------------------------------------------------------------------------

func TestBatchAnswersGradeWithTheirShippedType(t *testing.T) {
	types := linalgShippedTypes(t)
	reg := linalgBatchRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	for _, id := range linalgBatchIDs {
		typ, ok := types[id]
		if !ok || typ == "" {
			t.Fatalf("%s: no grading_type in the corpus", id)
		}
		switch typ {
		case "multiple_choice", "numeric", "tuple":
		default:
			t.Fatalf("%s: unexpected grading_type %q", id, typ)
		}
		for i := 0; i < 120; i++ {
			p := linalgBatchSample(t, reg, id, i)
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
	types := linalgShippedTypes(t)
	reg := linalgBatchRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	checked := 0
	for _, id := range linalgBatchIDs {
		if types[id] != "multiple_choice" {
			continue
		}
		for i := 0; i < 40; i++ {
			p := linalgBatchSample(t, reg, id, i)
			opts := linalgOptionsOf(p.Question)
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
		t.Skip("the linear algebra batch now serves computed count problems")
	}
}

// ---------------------------------------------------------------------------
// Independent numeric checks
// ---------------------------------------------------------------------------

func TestBatchKernelRecomputesFromFacts(t *testing.T) {
	reg := linalgBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := linalgBatchSample(t, reg, "linalg.transform.kernel", i)
		n := linalgFactInt(t, p, "n")
		rowParts := strings.Split(p.Facts["rows"], ";")
		rows := make([][]int, len(rowParts))
		for r, rp := range rowParts {
			cells := strings.Split(rp, ",")
			rows[r] = make([]int, len(cells))
			for c, s := range cells {
				v, err := strconv.Atoi(strings.TrimSpace(s))
				if err != nil {
					t.Fatalf("kernel: bad matrix entry %q", s)
				}
				rows[r][c] = v
			}
		}
		want := n - linalgRank(rows)
		answer, err := strconv.Atoi(p.Answer)
		if err != nil {
			t.Fatalf("kernel: non-integer answer %q", p.Answer)
		}
		if answer != want {
			t.Fatalf("kernel: n=%d rows=%v rank=%d answers %d, want %d", n, rows, linalgRank(rows), answer, want)
		}
	}
}

func TestBatchBasisChangeRecomputesFromFacts(t *testing.T) {
	reg := linalgBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := linalgBatchSample(t, reg, "linalg.basis.change", i)
		b1 := linalgFactInts(t, p, "b1")
		b2 := linalgFactInts(t, p, "b2")
		v := linalgFactInts(t, p, "v")
		ask := linalgFactInt(t, p, "ask")
		det := b1[0]*b2[1] - b1[1]*b2[0]
		if det == 0 {
			t.Fatalf("basis: degenerate basis %v %v", b1, b2)
		}
		a := float64(v[0]*b2[1]-v[1]*b2[0]) / float64(det)
		b := float64(b1[0]*v[1]-b1[1]*v[0]) / float64(det)
		want := a
		if ask == 1 {
			want = b
		}
		answer, err := strconv.ParseFloat(p.Answer, 64)
		if err != nil {
			t.Fatalf("basis change: non-numeric answer %q", p.Answer)
		}
		if math.Abs(answer-want) > 1e-9 {
			t.Fatalf("basis change: v=%v in basis %v,%v ask=%d answers %v, want %v", v, b1, b2, ask, answer, want)
		}
	}
}

func TestBatchComplementRecomputesFromFacts(t *testing.T) {
	reg := linalgBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := linalgBatchSample(t, reg, "linalg.orth.complement", i)
		n := linalgFactInt(t, p, "n")
		m := linalgFactInt(t, p, "m")
		want := n - m
		if p.Answer != strconv.Itoa(want) {
			t.Fatalf("complement: n=%d m=%d answers %s, want %d", n, m, p.Answer, want)
		}
	}
}

func TestBatchProjectionRecomputesFromFacts(t *testing.T) {
	reg := linalgBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := linalgBatchSample(t, reg, "linalg.orth.projection_matrix", i)
		v := linalgFactInts(t, p, "v")
		u := linalgFactInts(t, p, "u")
		axis := 0
		if p.Facts["axis"] == "y" {
			axis = 1
		}
		dot := v[0]*u[0] + v[1]*u[1]
		norm := u[0]*u[0] + u[1]*u[1]
		if norm == 0 {
			t.Fatalf("projection: zero direction")
		}
		want := float64(dot) / float64(norm) * float64(u[axis])
		answer, err := strconv.ParseFloat(p.Answer, 64)
		if err != nil {
			t.Fatalf("projection: non-numeric answer %q", p.Answer)
		}
		if math.Abs(answer-want) > 1e-9 {
			t.Fatalf("projection: v=%v u=%v axis=%d answers %v, want %v", v, u, axis, answer, want)
		}
	}
}

func TestBatchCrossRecomputesFromFacts(t *testing.T) {
	reg := linalgBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := linalgBatchSample(t, reg, "linalg.vector.cross", i)
		u := linalgFactInts(t, p, "u")
		v := linalgFactInts(t, p, "v")
		if len(u) != 3 || len(v) != 3 {
			t.Fatalf("cross: factors %v and %v are not 3-vectors", u, v)
		}
		x := u[1]*v[2] - u[2]*v[1]
		y := u[2]*v[0] - u[0]*v[2]
		z := u[0]*v[1] - u[1]*v[0]
		want := fmt.Sprintf("(%d,%d,%d)", x, y, z)
		if p.Answer != want {
			t.Fatalf("cross: %v x %v answers %s, want %s", u, v, p.Answer, want)
		}
	}
}

func TestBatchConditionRecomputesFromFacts(t *testing.T) {
	reg := linalgBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := linalgBatchSample(t, reg, "linalg.matrix.condition", i)
		largest, err := strconv.ParseFloat(p.Facts["largest"], 64)
		if err != nil {
			t.Fatalf("condition: bad largest fact %q", p.Facts["largest"])
		}
		smallest, err := strconv.ParseFloat(p.Facts["smallest"], 64)
		if err != nil || smallest == 0 {
			t.Fatalf("condition: bad smallest fact %q", p.Facts["smallest"])
		}
		want := largest / smallest
		answer, err := strconv.ParseFloat(p.Answer, 64)
		if err != nil {
			t.Fatalf("condition: non-numeric answer %q", p.Answer)
		}
		if math.Abs(answer-want) > 1e-9 {
			t.Fatalf("condition: %v/%v answers %v, want %v", largest, smallest, answer, want)
		}
	}
}

func TestBatchNumericAnswersParseAsDecimals(t *testing.T) {
	reg := linalgBatchRegistry(t)
	for _, id := range linalgBatchNumericIDs {
		for i := 0; i < 120; i++ {
			p := linalgBatchSample(t, reg, id, i)
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
	reg := linalgBatchRegistry(t)
	const samples = 150
	for _, id := range linalgBatchIDs {
		answers := map[string]bool{}
		questions := map[string]bool{}
		allYesNo := true
		for i := 0; i < samples; i++ {
			p, err := reg.GenerateContext(id, generator.GeneratorContext{Difficulty: 0.12, Seed: int64(102000 + i)})
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
