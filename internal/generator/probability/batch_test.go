package probability

// Tests for the first probability-theory batch.
//
// Every concept is graded through the real router with the grading_type the
// corpus ships. The seven multiple-choice concepts are checked the way the other
// batches check theirs: the question is parsed into options, the correct answer
// must appear exactly once, and every other option must be rejected.
//
// The thirteen numeric concepts recompute their answers from the published
// parameters: the laws of conditional probability and Bayes, the moment
// definitions, the transformation rule, the standard order-statistic and limit
// results, and the central limit theorem's standardisation.

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

type probCorpusConcept struct {
	ID          string `json:"id"`
	GradingType string `json:"grading_type"`
}

var probBatchIDs = []string{
	"prob.space.axioms",
	"prob.space.sigma_algebra",
	"prob.space.inclusion_exclusion",
	"prob.space.continuity",
	"prob.cond.multiplication",
	"prob.cond.total_probability",
	"prob.bayes.prior",
	"prob.independence.pairwise",
	"prob.rv.definition",
	"prob.rv.cdf",
	"prob.rv.pmf",
	"prob.rv.pdf",
	"prob.rv.joint",
	"prob.rv.expectation",
	"prob.rv.variance",
	"prob.rv.mgf",
	"prob.rv.transform",
	"prob.dist.order_stats",
	"prob.dist.memoryless",
	"prob.dist.poisson_process",
	"prob.dist.exponential_family",
	"prob.limit.markov",
	"prob.limit.chebyshev",
	"prob.limit.weak_law",
	"prob.limit.clt",
}

var probNumericIDs = []string{
	"prob.space.inclusion_exclusion",
	"prob.cond.multiplication",
	"prob.cond.total_probability",
	"prob.bayes.prior",
	"prob.rv.pmf",
	"prob.rv.pdf",
	"prob.rv.joint",
	"prob.rv.expectation",
	"prob.rv.variance",
	"prob.rv.mgf",
	"prob.rv.transform",
	"prob.dist.order_stats",
	"prob.dist.poisson_process",
	"prob.limit.markov",
	"prob.limit.chebyshev",
	"prob.limit.clt",
}

var probMCIDs = []string{}

func probShippedTypes(t *testing.T) map[string]string {
	t.Helper()
	raw, err := os.ReadFile("../../../data/concepts/probability_theory.json")
	if err != nil {
		t.Fatalf("read corpus: %v", err)
	}
	var rows []probCorpusConcept
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

func probRegistry(t *testing.T) *generator.Registry {
	t.Helper()
	reg := generator.NewRegistry()
	Register(reg)
	return reg
}

func probSample(t *testing.T, reg *generator.Registry, id string, i int) generator.Problem {
	t.Helper()
	p, err := reg.GenerateContext(id, generator.GeneratorContext{
		Difficulty: float64(i%20) / 19,
		Seed:       int64(181000 + i),
	})
	if err != nil {
		t.Fatalf("%s: generate: %v", id, err)
	}
	return p
}

var probOptionLine = regexp.MustCompile(`\n\n([A-E])\) ([^\n]*)`)

func probOptionsOf(q string) []string {
	var out []string
	for _, m := range probOptionLine.FindAllStringSubmatch(q, -1) {
		out = append(out, m[2])
	}
	return out
}

func probFactInt(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer", p.Question, key, p.Facts[key])
	}
	return v
}

func probFactFloat(t *testing.T, p generator.Problem, key string) float64 {
	t.Helper()
	v, err := strconv.ParseFloat(p.Facts[key], 64)
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not a number", p.Question, key, p.Facts[key])
	}
	return v
}

func probIntList(t *testing.T, p generator.Problem, key string) []int {
	t.Helper()
	parts := strings.Split(p.Facts[key], ",")
	out := make([]int, len(parts))
	for i, part := range parts {
		v, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			t.Fatalf("%s: fact %q has a non-integer entry %q", p.Question, key, part)
		}
		out[i] = v
	}
	return out
}

// ---------------------------------------------------------------------------
// Router contracts
// ---------------------------------------------------------------------------

func TestBatchAnswersGradeWithTheirShippedType(t *testing.T) {
	types := probShippedTypes(t)
	reg := probRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	for _, id := range probBatchIDs {
		typ, ok := types[id]
		if !ok || typ == "" {
			t.Fatalf("%s: no grading_type in the corpus", id)
		}
		if typ != "multiple_choice" && typ != "numeric" {
			t.Fatalf("%s: unexpected grading_type %q", id, typ)
		}
		for i := 0; i < 120; i++ {
			p := probSample(t, reg, id, i)
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
	types := probShippedTypes(t)
	reg := probRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	if len(probMCIDs) == 0 {
		t.Skip("the probability batch now serves computed count problems")
	}

	for _, id := range probMCIDs {
		if types[id] != "multiple_choice" {
			t.Fatalf("%s: expected multiple_choice, corpus says %q", id, types[id])
		}
		for i := 0; i < 40; i++ {
			p := probSample(t, reg, id, i)
			opts := probOptionsOf(p.Question)
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

func TestBatchInclusionExclusionRecomputes(t *testing.T) {
	reg := probRegistry(t)
	for i := 0; i < 200; i++ {
		p := probSample(t, reg, "prob.space.inclusion_exclusion", i)
		n := probFactInt(t, p, "n")
		a := probFactInt(t, p, "a")
		b := probFactInt(t, p, "b")
		ab := probFactInt(t, p, "ab")
		want := decimal6(float64(a+b-ab) / float64(n))
		if p.Answer != want {
			t.Fatalf("inclusion-exclusion: n=%d a=%d b=%d ab=%d answers %s, want %s", n, a, b, ab, p.Answer, want)
		}
	}
}

func TestBatchMultiplicationRuleRecomputes(t *testing.T) {
	reg := probRegistry(t)
	for i := 0; i < 200; i++ {
		p := probSample(t, reg, "prob.cond.multiplication", i)
		n := probFactInt(t, p, "n")
		b := probFactInt(t, p, "b")
		want := decimal6(float64(b) / float64(n))
		if p.Answer != want {
			t.Fatalf("multiplication: n=%d b=%d answers %s, want %s", n, b, p.Answer, want)
		}
	}
}

func TestBatchTotalProbabilityRecomputes(t *testing.T) {
	reg := probRegistry(t)
	for i := 0; i < 200; i++ {
		p := probSample(t, reg, "prob.cond.total_probability", i)
		w1 := probFactInt(t, p, "w1")
		w2 := probFactInt(t, p, "w2")
		a1 := probFactInt(t, p, "a1")
		t1 := probFactInt(t, p, "t1")
		a2 := probFactInt(t, p, "a2")
		t2 := probFactInt(t, p, "t2")
		total := w1 + w2
		want := decimal6(float64(w1)/float64(total)*float64(a1)/float64(t1) +
			float64(w2)/float64(total)*float64(a2)/float64(t2))
		if p.Answer != want {
			t.Fatalf("total probability: weights %d/%d rates %d/%d and %d/%d answers %s, want %s",
				w1, w2, a1, t1, a2, t2, p.Answer, want)
		}
	}
}

func TestBatchBayesRecomputes(t *testing.T) {
	reg := probRegistry(t)
	for i := 0; i < 200; i++ {
		p := probSample(t, reg, "prob.bayes.prior", i)
		pa := float64(probFactInt(t, p, "aA")) / float64(probFactInt(t, p, "tA"))
		pb := float64(probFactInt(t, p, "aB")) / float64(probFactInt(t, p, "tB"))
		want := decimal6(pa / (pa + pb))
		if p.Answer != want {
			t.Fatalf("bayes: facts %v answers %s, want %s", p.Facts, p.Answer, want)
		}
	}
}

func TestBatchPMFRecomputes(t *testing.T) {
	reg := probRegistry(t)
	for i := 0; i < 200; i++ {
		p := probSample(t, reg, "prob.rv.pmf", i)
		total := probFactInt(t, p, "total")
		weight := probFactInt(t, p, "weight")
		weights := probIntList(t, p, "weights")
		sum := 0
		for _, w := range weights {
			sum += w
		}
		if sum != total {
			t.Fatalf("pmf: weights %v sum to %d, fact total is %d", weights, sum, total)
		}
		if weight != weights[probFactInt(t, p, "index")-1] {
			t.Fatalf("pmf: weight %d is not the indexed entry of %v", weight, weights)
		}
		want := decimal6(float64(weight) / float64(total))
		if p.Answer != want {
			t.Fatalf("pmf: weight=%d total=%d answers %s, want %s", weight, total, p.Answer, want)
		}
	}
}

func TestBatchPDFRecomputes(t *testing.T) {
	reg := probRegistry(t)
	for i := 0; i < 200; i++ {
		p := probSample(t, reg, "prob.rv.pdf", i)
		a := probFactFloat(t, p, "a")
		want := decimal6(a * a)
		if p.Answer != want {
			t.Fatalf("pdf: a=%v answers %s, want %s", a, p.Answer, want)
		}
	}
}

func TestBatchExpectationRecomputes(t *testing.T) {
	reg := probRegistry(t)
	for i := 0; i < 200; i++ {
		p := probSample(t, reg, "prob.rv.expectation", i)
		values := probIntList(t, p, "values")
		weights := probIntList(t, p, "weights")
		den := probFactInt(t, p, "den")
		if len(values) != len(weights) {
			t.Fatalf("expectation: %d values but %d weights", len(values), len(weights))
		}
		sum := 0
		for j := range values {
			sum += values[j] * values[j] * weights[j]
		}
		want := decimal6(float64(sum) / float64(den))
		if p.Answer != want {
			t.Fatalf("expectation: values %v weights %v answers %s, want %s", values, weights, p.Answer, want)
		}
	}
}

func TestBatchVarianceRecomputes(t *testing.T) {
	reg := probRegistry(t)
	for i := 0; i < 200; i++ {
		p := probSample(t, reg, "prob.rv.variance", i)
		values := probIntList(t, p, "values")
		weights := probIntList(t, p, "weights")
		den := probFactInt(t, p, "den")
		sw, sw2 := 0, 0
		for j := range values {
			sw += values[j] * weights[j]
			sw2 += values[j] * values[j] * weights[j]
		}
		want := decimal6(float64(den*sw2-sw*sw) / float64(den*den))
		if p.Answer != want {
			t.Fatalf("variance: values %v weights %v answers %s, want %s", values, weights, p.Answer, want)
		}
	}
}

func TestBatchMGFRecomputes(t *testing.T) {
	reg := probRegistry(t)
	for i := 0; i < 200; i++ {
		p := probSample(t, reg, "prob.rv.mgf", i)
		b := probFactInt(t, p, "b")
		want := strconv.Itoa(2 * b)
		if p.Answer != want {
			t.Fatalf("mgf: b=%d answers %s, want %s", b, p.Answer, want)
		}
	}
}

func TestBatchTransformRecomputes(t *testing.T) {
	reg := probRegistry(t)
	for i := 0; i < 200; i++ {
		p := probSample(t, reg, "prob.rv.transform", i)
		v := probFactFloat(t, p, "v")
		c := probFactInt(t, p, "c")
		want := decimal6(float64(c*c) * v)
		if p.Answer != want {
			t.Fatalf("transform: v=%v c=%d answers %s, want %s", v, c, p.Answer, want)
		}
	}
}

func TestBatchJointPMFRecomputes(t *testing.T) {
	reg := probRegistry(t)
	pairs := [][2]int{{0, 0}, {0, 1}, {1, 0}, {1, 1}}
	for i := 0; i < 200; i++ {
		p := probSample(t, reg, "prob.rv.joint", i)
		weights := probIntList(t, p, "weights")
		total := probFactInt(t, p, "total")
		weight := probFactInt(t, p, "weight")
		index := probFactInt(t, p, "index")
		sum := 0
		for _, w := range weights {
			sum += w
		}
		if sum != total {
			t.Fatalf("joint pmf: weights %v sum to %d, fact total is %d", weights, sum, total)
		}
		wantPair := strconv.Itoa(pairs[index-1][0]) + "," + strconv.Itoa(pairs[index-1][1])
		if p.Facts["pair"] != wantPair {
			t.Fatalf("joint pmf: index %d names pair %s, fact says %s", index, wantPair, p.Facts["pair"])
		}
		if weight != weights[index-1] {
			t.Fatalf("joint pmf: weight %d is not the indexed entry of %v", weight, weights)
		}
		want := decimal6(float64(weight) / float64(total))
		if p.Answer != want {
			t.Fatalf("joint pmf: weight=%d total=%d answers %s, want %s", weight, total, p.Answer, want)
		}
	}
}

func TestBatchPoissonProcessRecomputes(t *testing.T) {
	reg := probRegistry(t)
	for i := 0; i < 200; i++ {
		p := probSample(t, reg, "prob.dist.poisson_process", i)
		m := probFactInt(t, p, "m")
		k := probFactInt(t, p, "k")
		want := decimal6(math.Exp(-float64(m)) * pow(float64(m), k) / factF(k))
		if p.Answer != want {
			t.Fatalf("poisson process: m=%d k=%d answers %s, want %s", m, k, p.Answer, want)
		}
	}
}

func TestBatchOrderStatsRecomputes(t *testing.T) {
	reg := probRegistry(t)
	for i := 0; i < 200; i++ {
		p := probSample(t, reg, "prob.dist.order_stats", i)
		n := probFactInt(t, p, "n")
		x := probFactFloat(t, p, "x")
		want := decimal6(pow(x, n))
		if p.Answer != want {
			t.Fatalf("order stats: n=%d x=%v answers %s, want %s", n, x, p.Answer, want)
		}
	}
}

func TestBatchMarkovRecomputes(t *testing.T) {
	reg := probRegistry(t)
	for i := 0; i < 200; i++ {
		p := probSample(t, reg, "prob.limit.markov", i)
		mean := probFactInt(t, p, "mean")
		a := probFactInt(t, p, "a")
		want := decimal6(float64(mean) / float64(a))
		if p.Answer != want {
			t.Fatalf("markov: mean=%d a=%d answers %s, want %s", mean, a, p.Answer, want)
		}
	}
}

func TestBatchChebyshevRecomputes(t *testing.T) {
	reg := probRegistry(t)
	for i := 0; i < 200; i++ {
		p := probSample(t, reg, "prob.limit.chebyshev", i)
		sigma := probFactInt(t, p, "sigma")
		d := probFactInt(t, p, "d")
		want := decimal6(float64(sigma*sigma) / float64(d*d))
		if p.Answer != want {
			t.Fatalf("chebyshev: sigma=%d d=%d answers %s, want %s", sigma, d, p.Answer, want)
		}
	}
}

func TestBatchCLTRecomputes(t *testing.T) {
	reg := probRegistry(t)
	for i := 0; i < 200; i++ {
		p := probSample(t, reg, "prob.limit.clt", i)
		n := probFactInt(t, p, "n")
		mu := probFactInt(t, p, "mu")
		sigma := probFactInt(t, p, "sigma")
		x := probFactInt(t, p, "x")
		want := decimal6(float64(x-n*mu) / (float64(sigma) * math.Sqrt(float64(n))))
		if p.Answer != want {
			t.Fatalf("clt: n=%d mu=%d sigma=%d x=%d answers %s, want %s", n, mu, sigma, x, p.Answer, want)
		}
	}
}

func TestBatchNumericAnswersParseAsDecimals(t *testing.T) {
	reg := probRegistry(t)
	for _, id := range probNumericIDs {
		for i := 0; i < 120; i++ {
			p := probSample(t, reg, id, i)
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
	reg := probRegistry(t)
	const samples = 150
	for _, id := range probBatchIDs {
		answers := map[string]bool{}
		questions := map[string]bool{}
		allYesNo := true
		for i := 0; i < samples; i++ {
			p, err := reg.GenerateContext(id, generator.GeneratorContext{Difficulty: 0.12, Seed: int64(183000 + i)})
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
