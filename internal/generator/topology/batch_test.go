package topology

// Tests for the 2026-10-18 topology batch.
//
// Every concept is graded through the real router with the grading_type the
// corpus ships. The multiple-choice concepts are checked the way the other
// batches check theirs: the question is parsed into options, the correct answer
// must appear exactly once, and every other option must be rejected.
//
// The six numeric concepts are recomputed from the facts the generator
// publishes: the Euclidean distance and ball endpoints, the boundary/interior
// count, the leading coefficient ratio of the sequence, V - E + F, the genus
// formula chi = 2 - 2g, and the Betti number of each space family.

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

var topoBatchIDs = []string{
	"topo.metric.ball",
	"topo.metric.interior",
	"topo.metric.limit_point",
	"topo.metric.convergence",
	"topo.metric.complete",
	"topo.metric.basis",
	"topo.metric.subspace",
	"topo.metric.dense",
	"topo.invariant.homeomorphism",
	"topo.invariant.euler",
	"topo.invariant.genus",
	"topo.invariant.betti",
}

var topoBatchNumericIDs = []string{
	"topo.metric.ball",
	"topo.metric.interior",
	"topo.metric.convergence",
	"topo.metric.limit_point",
	"topo.metric.complete",
	"topo.metric.basis",
	"topo.metric.subspace",
	"topo.metric.dense",
	"topo.invariant.homeomorphism",
	"topo.invariant.euler",
	"topo.invariant.genus",
	"topo.invariant.betti",
}

type topoCorpusConcept struct {
	ID          string `json:"id"`
	GradingType string `json:"grading_type"`
}

func topoShippedTypes(t *testing.T) map[string]string {
	t.Helper()
	raw, err := os.ReadFile("../../../data/concepts/topology.json")
	if err != nil {
		t.Fatalf("read corpus: %v", err)
	}
	var rows []topoCorpusConcept
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

func topoRegistry(t *testing.T) *generator.Registry {
	t.Helper()
	reg := generator.NewRegistry()
	Register(reg)
	return reg
}

func topoSample(t *testing.T, reg *generator.Registry, id string, i int) generator.Problem {
	t.Helper()
	p, err := reg.GenerateContext(id, generator.GeneratorContext{
		Difficulty: float64(i%20) / 19,
		Seed:       int64(91000 + i),
	})
	if err != nil {
		t.Fatalf("%s: generate: %v", id, err)
	}
	return p
}

var topoOptionLine = regexp.MustCompile(`\n\n([A-E])\) ([^\n]*)`)

func topoOptionsOf(q string) []string {
	var out []string
	for _, m := range topoOptionLine.FindAllStringSubmatch(q, -1) {
		out = append(out, m[2])
	}
	return out
}

func topoFactInt(t *testing.T, p generator.Problem, key string) int {
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
	types := topoShippedTypes(t)
	reg := topoRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	for _, id := range topoBatchIDs {
		typ, ok := types[id]
		if !ok || typ == "" {
			t.Fatalf("%s: no grading_type in the corpus", id)
		}
		if typ != "multiple_choice" && typ != "numeric" {
			t.Fatalf("%s: unexpected grading_type %q", id, typ)
		}
		for i := 0; i < 120; i++ {
			p := topoSample(t, reg, id, i)
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
	types := topoShippedTypes(t)
	reg := topoRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	checked := 0
	for _, id := range topoBatchIDs {
		if types[id] != "multiple_choice" {
			continue
		}
		for i := 0; i < 40; i++ {
			p := topoSample(t, reg, id, i)
			opts := topoOptionsOf(p.Question)
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
		t.Skip("the topology batch is fully numeric; no multiple-choice batch questions remain")
	}
}

// ---------------------------------------------------------------------------
// Independent numeric checks
// ---------------------------------------------------------------------------

func TestBatchBallRecomputesFromFacts(t *testing.T) {
	reg := topoRegistry(t)
	seenDistance, seenEndpoint := false, false
	for i := 0; i < 200; i++ {
		p := topoSample(t, reg, "topo.metric.ball", i)
		answer, err := strconv.ParseFloat(p.Answer, 64)
		if err != nil {
			t.Fatalf("ball: non-numeric answer %q", p.Answer)
		}
		switch p.Facts["kind"] {
		case "distance":
			seenDistance = true
			x1 := topoFactInt(t, p, "x1")
			y1 := topoFactInt(t, p, "y1")
			x2 := topoFactInt(t, p, "x2")
			y2 := topoFactInt(t, p, "y2")
			want := math.Sqrt(float64((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1)))
			if math.Abs(answer-want) > 1e-9 {
				t.Fatalf("ball: distance from (%d,%d) to (%d,%d) answers %v, want %v",
					x1, y1, x2, y2, answer, want)
			}
		case "endpoint":
			seenEndpoint = true
			center := topoFactInt(t, p, "center")
			radius := topoFactInt(t, p, "radius")
			want := center + radius
			if p.Facts["side"] == "lower" {
				want = center - radius
			}
			if answer != float64(want) {
				t.Fatalf("ball: %s endpoint of B(%d,%d) answers %v, want %d",
					p.Facts["side"], center, radius, answer, want)
			}
		default:
			t.Fatalf("ball: unknown kind %q", p.Facts["kind"])
		}
	}
	if !seenDistance || !seenEndpoint {
		t.Fatalf("ball: branches not both exercised (distance=%v endpoint=%v)", seenDistance, seenEndpoint)
	}
}

func TestBatchInteriorRecomputesFromFacts(t *testing.T) {
	reg := topoRegistry(t)
	for i := 0; i < 200; i++ {
		p := topoSample(t, reg, "topo.metric.interior", i)
		want := 0
		if p.Facts["points"] != "" {
			want = len(strings.Split(p.Facts["points"], ","))
		}
		answer, err := strconv.Atoi(p.Answer)
		if err != nil {
			t.Fatalf("interior: non-integer answer %q", p.Answer)
		}
		if answer != want {
			t.Fatalf("interior: %s count with points %q answers %d, want %d",
				p.Facts["kind"], p.Facts["points"], answer, want)
		}
	}
}

func TestBatchConvergenceRecomputesFromFacts(t *testing.T) {
	reg := topoRegistry(t)
	for i := 0; i < 200; i++ {
		p := topoSample(t, reg, "topo.metric.convergence", i)
		a := topoFactInt(t, p, "a")
		c := topoFactInt(t, p, "c")
		if c == 0 {
			t.Fatalf("convergence: zero leading denominator")
		}
		want := float64(a) / float64(c)
		answer, err := strconv.ParseFloat(p.Answer, 64)
		if err != nil {
			t.Fatalf("convergence: non-numeric answer %q", p.Answer)
		}
		if math.Abs(answer-want) > 1e-9 {
			t.Fatalf("convergence: a=%d c=%d answers %v, want %v", a, c, answer, want)
		}
	}
}

func TestBatchEulerRecomputesFromFacts(t *testing.T) {
	reg := topoRegistry(t)
	for i := 0; i < 200; i++ {
		p := topoSample(t, reg, "topo.invariant.euler", i)
		v := topoFactInt(t, p, "v")
		e := topoFactInt(t, p, "e")
		f := topoFactInt(t, p, "f")
		want := v - e + f
		if p.Answer != strconv.Itoa(want) {
			t.Fatalf("euler: V=%d E=%d F=%d answers %s, want %d", v, e, f, p.Answer, want)
		}
	}
}

func TestBatchGenusRecomputesFromFacts(t *testing.T) {
	reg := topoRegistry(t)
	for i := 0; i < 200; i++ {
		p := topoSample(t, reg, "topo.invariant.genus", i)
		var want int
		switch p.Facts["kind"] {
		case "genus":
			chi := topoFactInt(t, p, "chi")
			want = (2 - chi) / 2
		case "chi":
			g := topoFactInt(t, p, "genus")
			want = 2 - 2*g
		default:
			t.Fatalf("genus: unknown kind %q", p.Facts["kind"])
		}
		if p.Answer != strconv.Itoa(want) {
			t.Fatalf("genus: %s with %v answers %s, want %d", p.Facts["kind"], p.Facts, p.Answer, want)
		}
	}
}

func TestBatchBettiRecomputesFromFacts(t *testing.T) {
	reg := topoRegistry(t)
	for i := 0; i < 200; i++ {
		p := topoSample(t, reg, "topo.invariant.betti", i)
		param := topoFactInt(t, p, "param")
		var want int
		switch p.Facts["kind"] {
		case "circle":
			want = 1
		case "surface":
			want = 2 * param
		case "wedge":
			want = param
		case "components":
			want = param
		case "tree":
			want = 0
		case "sphere_euler":
			want = 2
		default:
			t.Fatalf("betti: unknown kind %q", p.Facts["kind"])
		}
		if p.Answer != strconv.Itoa(want) {
			t.Fatalf("betti: %s with param %d answers %s, want %d",
				p.Facts["kind"], param, p.Answer, want)
		}
	}
}

func TestBatchNumericAnswersParseAsDecimals(t *testing.T) {
	reg := topoRegistry(t)
	for _, id := range topoBatchNumericIDs {
		for i := 0; i < 120; i++ {
			p := topoSample(t, reg, id, i)
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
	reg := topoRegistry(t)
	const samples = 150
	for _, id := range topoBatchIDs {
		answers := map[string]bool{}
		questions := map[string]bool{}
		allYesNo := true
		for i := 0; i < samples; i++ {
			p, err := reg.GenerateContext(id, generator.GeneratorContext{Difficulty: 0.12, Seed: int64(92000 + i)})
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
