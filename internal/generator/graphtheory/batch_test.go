package graphtheory

// Tests for the first graph-theory batch.
//
// Every concept is graded through the real router with the grading_type the
// corpus ships. The ten multiple-choice concepts are checked the way the other
// batches check theirs: the question is parsed into options, the correct answer
// must appear exactly once, and every other option must be rejected.
//
// The ten numeric concepts recompute their answers from the published graph
// parameters: the vertex and edge counts of the complement and product, the
// named graph families, the matching and cover sizes, the chromatic polynomial,
// and the edge- and list-coloring numbers.

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

var graphBatchIDs = []string{
	"graph.isomorphism",
	"graph.complement",
	"graph.subgraph",
	"graph.homomorphism",
	"graph.product",
	"graph.distance",
	"graph.cut_vertex",
	"graph.bridge",
	"graph.connectivity",
	"graph.menger",
	"graph.matching",
	"graph.perfect_matching",
	"graph.tutte",
	"graph.vertex_cover",
	"graph.independent_set",
	"graph.chromatic_polynomial",
	"graph.brooks",
	"graph.edge_coloring",
	"graph.perfect_graph",
	"graph.list_coloring",
}

var graphBatchNumericIDs = []string{
	"graph.complement",
	"graph.product",
	"graph.distance",
	"graph.cut_vertex",
	"graph.bridge",
	"graph.connectivity",
	"graph.matching",
	"graph.vertex_cover",
	"graph.independent_set",
	"graph.chromatic_polynomial",
	"graph.edge_coloring",
	"graph.list_coloring",
}

type graphCorpusConcept struct {
	ID          string `json:"id"`
	GradingType string `json:"grading_type"`
}

func graphShippedTypes(t *testing.T) map[string]string {
	t.Helper()
	raw, err := os.ReadFile("../../../data/concepts/graph_theory.json")
	if err != nil {
		t.Fatalf("read corpus: %v", err)
	}
	var rows []graphCorpusConcept
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

func graphBatchRegistry(t *testing.T) *generator.Registry {
	t.Helper()
	reg := generator.NewRegistry()
	Register(reg)
	return reg
}

func graphBatchSample(t *testing.T, reg *generator.Registry, id string, i int) generator.Problem {
	t.Helper()
	p, err := reg.GenerateContext(id, generator.GeneratorContext{
		Difficulty: float64(i%20) / 19,
		Seed:       int64(161000 + i),
	})
	if err != nil {
		t.Fatalf("%s: generate: %v", id, err)
	}
	return p
}

var graphOptionLine = regexp.MustCompile(`\n\n([A-E])\) ([^\n]*)`)

func graphOptionsOf(q string) []string {
	var out []string
	for _, m := range graphOptionLine.FindAllStringSubmatch(q, -1) {
		out = append(out, m[2])
	}
	return out
}

func graphFactInt(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer", p.Question, key, p.Facts[key])
	}
	return v
}

func graphPow(base, exp int) int {
	result := 1
	for i := 0; i < exp; i++ {
		result *= base
	}
	return result
}

// ---------------------------------------------------------------------------
// Router contracts
// ---------------------------------------------------------------------------

func TestBatchAnswersGradeWithTheirShippedType(t *testing.T) {
	types := graphShippedTypes(t)
	reg := graphBatchRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	for _, id := range graphBatchIDs {
		typ, ok := types[id]
		if !ok || typ == "" {
			t.Fatalf("%s: no grading_type in the corpus", id)
		}
		if typ != "multiple_choice" && typ != "numeric" {
			t.Fatalf("%s: unexpected grading_type %q", id, typ)
		}
		for i := 0; i < 120; i++ {
			p := graphBatchSample(t, reg, id, i)
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
	types := graphShippedTypes(t)
	reg := graphBatchRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	checked := 0
	for _, id := range graphBatchIDs {
		if types[id] != "multiple_choice" {
			continue
		}
		for i := 0; i < 40; i++ {
			p := graphBatchSample(t, reg, id, i)
			opts := graphOptionsOf(p.Question)
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
		t.Skip("the graph-theory batch now serves computed count problems")
	}
}

// ---------------------------------------------------------------------------
// Independent numeric checks
// ---------------------------------------------------------------------------

func TestBatchComplementRecomputes(t *testing.T) {
	reg := graphBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := graphBatchSample(t, reg, "graph.complement", i)
		n := graphFactInt(t, p, "n")
		m := graphFactInt(t, p, "m")
		want := n*(n-1)/2 - m
		if p.Answer != strconv.Itoa(want) {
			t.Fatalf("complement: n=%d m=%d answers %s, want %d", n, m, p.Answer, want)
		}
	}
}

func TestBatchProductRecomputes(t *testing.T) {
	reg := graphBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := graphBatchSample(t, reg, "graph.product", i)
		n1 := graphFactInt(t, p, "n1")
		m1 := graphFactInt(t, p, "m1")
		n2 := graphFactInt(t, p, "n2")
		m2 := graphFactInt(t, p, "m2")
		want := m1*n2 + n1*m2
		if p.Answer != strconv.Itoa(want) {
			t.Fatalf("product: (%d,%d) x (%d,%d) answers %s, want %d", n1, m1, n2, m2, p.Answer, want)
		}
	}
}

func TestBatchDistanceRecomputes(t *testing.T) {
	reg := graphBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := graphBatchSample(t, reg, "graph.distance", i)
		var want int
		switch p.Facts["kind"] {
		case "path":
			want = graphFactInt(t, p, "n") - 1
		case "cycle":
			want = graphFactInt(t, p, "n") / 2
		case "complete":
			want = 1
		case "grid":
			want = graphFactInt(t, p, "rows") + graphFactInt(t, p, "cols") - 2
		default:
			t.Fatalf("distance: unknown kind %q", p.Facts["kind"])
		}
		if p.Answer != strconv.Itoa(want) {
			t.Fatalf("distance: %s n=%v answers %s, want %d", p.Facts["kind"], p.Facts["n"], p.Answer, want)
		}
	}
}

func TestBatchCutVertexRecomputes(t *testing.T) {
	reg := graphBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := graphBatchSample(t, reg, "graph.cut_vertex", i)
		var want int
		switch p.Facts["kind"] {
		case "path":
			want = graphFactInt(t, p, "n") - 2
		case "star":
			want = 1
		case "cycle":
			want = 0
		case "bowtie":
			want = 1
		default:
			t.Fatalf("cut vertex: unknown kind %q", p.Facts["kind"])
		}
		if p.Answer != strconv.Itoa(want) {
			t.Fatalf("cut vertex: %s n=%v answers %s, want %d", p.Facts["kind"], p.Facts["n"], p.Answer, want)
		}
	}
}

func TestBatchBridgeRecomputes(t *testing.T) {
	reg := graphBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := graphBatchSample(t, reg, "graph.bridge", i)
		var want int
		switch p.Facts["kind"] {
		case "path":
			want = graphFactInt(t, p, "n") - 1
		case "cycle":
			want = 0
		case "star":
			want = graphFactInt(t, p, "n")
		case "tree":
			want = graphFactInt(t, p, "n") - 1
		case "join":
			want = 1
		case "chord":
			want = 0
		default:
			t.Fatalf("bridge: unknown kind %q", p.Facts["kind"])
		}
		if p.Answer != strconv.Itoa(want) {
			t.Fatalf("bridge: %s n=%v answers %s, want %d", p.Facts["kind"], p.Facts["n"], p.Answer, want)
		}
	}
}

func TestBatchConnectivityRecomputes(t *testing.T) {
	reg := graphBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := graphBatchSample(t, reg, "graph.connectivity", i)
		var want int
		switch p.Facts["kind"] {
		case "complete":
			want = graphFactInt(t, p, "a") - 1
		case "cycle":
			want = 2
		case "path":
			want = 1
		case "bipartite":
			a := graphFactInt(t, p, "a")
			b := graphFactInt(t, p, "b")
			want = int(math.Min(float64(a), float64(b)))
		case "disjoint":
			want = 0
		default:
			t.Fatalf("connectivity: unknown kind %q", p.Facts["kind"])
		}
		if p.Answer != strconv.Itoa(want) {
			t.Fatalf("connectivity: %s answers %s, want %d", p.Facts["kind"], p.Answer, want)
		}
	}
}

func TestBatchMatchingRecomputes(t *testing.T) {
	reg := graphBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := graphBatchSample(t, reg, "graph.matching", i)
		var want int
		switch p.Facts["kind"] {
		case "path", "cycle":
			want = graphFactInt(t, p, "n") / 2
		case "star":
			want = 1
		default:
			t.Fatalf("matching: unknown kind %q", p.Facts["kind"])
		}
		if p.Answer != strconv.Itoa(want) {
			t.Fatalf("matching: %s n=%v answers %s, want %d", p.Facts["kind"], p.Facts["n"], p.Answer, want)
		}
	}
}

func TestBatchVertexCoverRecomputes(t *testing.T) {
	reg := graphBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := graphBatchSample(t, reg, "graph.vertex_cover", i)
		var want int
		switch p.Facts["kind"] {
		case "path":
			want = graphFactInt(t, p, "a") / 2
		case "cycle":
			want = graphFactInt(t, p, "a") / 2
		case "star":
			want = 1
		case "bipartite":
			a := graphFactInt(t, p, "a")
			b := graphFactInt(t, p, "b")
			want = int(math.Min(float64(a), float64(b)))
		default:
			t.Fatalf("vertex cover: unknown kind %q", p.Facts["kind"])
		}
		if p.Answer != strconv.Itoa(want) {
			t.Fatalf("vertex cover: %s answers %s, want %d", p.Facts["kind"], p.Answer, want)
		}
	}
}

func TestBatchIndependentSetRecomputes(t *testing.T) {
	reg := graphBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := graphBatchSample(t, reg, "graph.independent_set", i)
		var want int
		switch p.Facts["kind"] {
		case "path":
			want = (graphFactInt(t, p, "a") + 1) / 2
		case "cycle":
			want = graphFactInt(t, p, "a") / 2
		case "star":
			want = graphFactInt(t, p, "a")
		case "bipartite":
			a := graphFactInt(t, p, "a")
			b := graphFactInt(t, p, "b")
			want = int(math.Max(float64(a), float64(b)))
		default:
			t.Fatalf("independent set: unknown kind %q", p.Facts["kind"])
		}
		if p.Answer != strconv.Itoa(want) {
			t.Fatalf("independent set: %s answers %s, want %d", p.Facts["kind"], p.Answer, want)
		}
	}
}

func TestBatchChromaticPolynomialRecomputes(t *testing.T) {
	reg := graphBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := graphBatchSample(t, reg, "graph.chromatic_polynomial", i)
		n := graphFactInt(t, p, "n")
		k := graphFactInt(t, p, "k")
		var want int
		switch p.Facts["kind"] {
		case "path", "tree":
			want = k * graphPow(k-1, n-1)
		case "cycle":
			sign := 1
			if n%2 == 1 {
				sign = -1
			}
			want = graphPow(k-1, n) + sign*(k-1)
		case "complete":
			want = 1
			for j := 0; j < n; j++ {
				want *= k - j
			}
		default:
			t.Fatalf("chromatic polynomial: unknown kind %q", p.Facts["kind"])
		}
		if p.Answer != strconv.Itoa(want) {
			t.Fatalf("chromatic polynomial: %s n=%d k=%d answers %s, want %d", p.Facts["kind"], n, k, p.Answer, want)
		}
	}
}

func TestBatchEdgeColoringRecomputes(t *testing.T) {
	reg := graphBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := graphBatchSample(t, reg, "graph.edge_coloring", i)
		n := graphFactInt(t, p, "n")
		var want int
		switch p.Facts["kind"] {
		case "complete":
			if n%2 == 0 {
				want = n - 1
			} else {
				want = n
			}
		case "cycle":
			if n%2 == 0 {
				want = 2
			} else {
				want = 3
			}
		case "path":
			want = 2
		case "star":
			want = n
		default:
			t.Fatalf("edge coloring: unknown kind %q", p.Facts["kind"])
		}
		if p.Answer != strconv.Itoa(want) {
			t.Fatalf("edge coloring: %s n=%d answers %s, want %d", p.Facts["kind"], n, p.Answer, want)
		}
	}
}

func TestBatchListColoringRecomputes(t *testing.T) {
	reg := graphBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := graphBatchSample(t, reg, "graph.list_coloring", i)
		n := graphFactInt(t, p, "n")
		var want int
		switch p.Facts["kind"] {
		case "cycle":
			if n%2 == 0 {
				want = 2
			} else {
				want = 3
			}
		case "path":
			want = 2
		default:
			t.Fatalf("list coloring: unknown kind %q", p.Facts["kind"])
		}
		if p.Answer != strconv.Itoa(want) {
			t.Fatalf("list coloring: %s n=%d answers %s, want %d", p.Facts["kind"], n, p.Answer, want)
		}
	}
}

func TestBatchNumericAnswersParseAsDecimals(t *testing.T) {
	reg := graphBatchRegistry(t)
	for _, id := range graphBatchNumericIDs {
		for i := 0; i < 120; i++ {
			p := graphBatchSample(t, reg, id, i)
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
	reg := graphBatchRegistry(t)
	const samples = 150
	for _, id := range graphBatchIDs {
		answers := map[string]bool{}
		questions := map[string]bool{}
		allYesNo := true
		for i := 0; i < samples; i++ {
			p, err := reg.GenerateContext(id, generator.GeneratorContext{Difficulty: 0.12, Seed: int64(162000 + i)})
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

// TestBatchNumericFactsAppearInTheStem pins the class of defect that shipped
// in prob.dist.order_stats, where the answer key held n = 3 while the stem said
// "let n observations". A learner cannot answer a question that names a
// parameter only in the key. Every numeric fact here has to be readable in the
// question, so the graph it describes is the graph the answer belongs to.
func TestBatchNumericFactsAppearInTheStem(t *testing.T) {
	reg := graphBatchRegistry(t)
	converted := []string{
		"graph.complement", "graph.product", "graph.distance",
		"graph.cut_vertex", "graph.bridge", "graph.matching",
	}
	for _, id := range converted {
		for i := 0; i < 100; i++ {
			p := graphBatchSample(t, reg, id, i)
			for key, raw := range p.Facts {
				if key == "answer" {
					continue
				}
				value, err := strconv.Atoi(raw)
				if err != nil {
					continue
				}
				pattern := regexp.MustCompile(`(^|\D)` + regexp.QuoteMeta(strconv.Itoa(value)) + `(\D|$)`)
				if !pattern.MatchString(p.Question) {
					t.Fatalf("%s: %s names %s = %d, but the question does not print it: %s",
						id, id, key, value, p.Question)
				}
			}
		}
	}
}
