package statistics

// Tests for the 2026-10-18 statistics batch.
//
// Every concept is graded through the real router with the grading_type the
// corpus ships. The multiple-choice concepts are checked the way the other
// batches check theirs: the question is parsed into options, the correct answer
// must appear exactly once, and every other option must be rejected.
//
// The six numeric concepts are recomputed from the facts the generator
// publishes: z and t from the mean, null value, spread and sample size; the
// proportion standard error; chi-square from the observed and expected counts;
// the ANOVA F from the groups; R-squared from the two sums of squares; and the
// histogram count from the values and the half-open bin.

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

var statsBatchIDs = []string{
	"stat.infer.hypothesis",
	"stat.infer.ztest",
	"stat.infer.ttest",
	"stat.infer.pvalue",
	"stat.infer.type_error",
	"stat.infer.proportion",
	"stat.test.chi_square",
	"stat.test.anova",
	"stat.reg.r2",
	"stat.data.histogram",
	"stat.infer.nonparametric",
}

var statsBatchNumericIDs = []string{
	"stat.infer.ztest",
	"stat.infer.ttest",
	"stat.infer.proportion",
	"stat.test.chi_square",
	"stat.test.anova",
	"stat.reg.r2",
	"stat.data.histogram",
}

type statsCorpusConcept struct {
	ID          string `json:"id"`
	GradingType string `json:"grading_type"`
}

func statsShippedTypes(t *testing.T) map[string]string {
	t.Helper()
	raw, err := os.ReadFile("../../../data/concepts/statistics.json")
	if err != nil {
		t.Fatalf("read corpus: %v", err)
	}
	var rows []statsCorpusConcept
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

func statsBatchRegistry(t *testing.T) *generator.Registry {
	t.Helper()
	reg := generator.NewRegistry()
	Register(reg)
	return reg
}

func statsBatchSample(t *testing.T, reg *generator.Registry, id string, i int) generator.Problem {
	t.Helper()
	p, err := reg.GenerateContext(id, generator.GeneratorContext{
		Difficulty: float64(i%20) / 19,
		Seed:       int64(111000 + i),
	})
	if err != nil {
		t.Fatalf("%s: generate: %v", id, err)
	}
	return p
}

var statsOptionLine = regexp.MustCompile(`\n\n([A-E])\) ([^\n]*)`)

func statsOptionsOf(q string) []string {
	var out []string
	for _, m := range statsOptionLine.FindAllStringSubmatch(q, -1) {
		out = append(out, m[2])
	}
	return out
}

func statsFactFloat(t *testing.T, p generator.Problem, key string) float64 {
	t.Helper()
	v, err := strconv.ParseFloat(p.Facts[key], 64)
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not a number", p.Question, key, p.Facts[key])
	}
	return v
}

func statsFactInt(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer", p.Question, key, p.Facts[key])
	}
	return v
}

func statsFactInts(t *testing.T, p generator.Problem, key string) []int {
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

func statsAnovaF(groups [][]int) float64 {
	total := 0
	sum := 0
	for _, g := range groups {
		total += len(g)
		for _, v := range g {
			sum += v
		}
	}
	grand := float64(sum) / float64(total)
	ssb, ssw := 0.0, 0.0
	for _, g := range groups {
		mean := 0.0
		for _, v := range g {
			mean += float64(v)
		}
		mean /= float64(len(g))
		ssb += float64(len(g)) * (mean - grand) * (mean - grand)
		for _, v := range g {
			d := float64(v) - mean
			ssw += d * d
		}
	}
	dfb := float64(len(groups) - 1)
	dfw := float64(total - len(groups))
	return (ssb / dfb) / (ssw / dfw)
}

// ---------------------------------------------------------------------------
// Router contracts
// ---------------------------------------------------------------------------

func TestBatchAnswersGradeWithTheirShippedType(t *testing.T) {
	types := statsShippedTypes(t)
	reg := statsBatchRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	for _, id := range statsBatchIDs {
		typ, ok := types[id]
		if !ok || typ == "" {
			t.Fatalf("%s: no grading_type in the corpus", id)
		}
		if typ != "multiple_choice" && typ != "numeric" {
			t.Fatalf("%s: unexpected grading_type %q", id, typ)
		}
		for i := 0; i < 120; i++ {
			p := statsBatchSample(t, reg, id, i)
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
	types := statsShippedTypes(t)
	reg := statsBatchRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	checked := 0
	for _, id := range statsBatchIDs {
		if types[id] != "multiple_choice" {
			continue
		}
		for i := 0; i < 40; i++ {
			p := statsBatchSample(t, reg, id, i)
			opts := statsOptionsOf(p.Question)
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
		t.Skip("the statistics batch now serves computed answers")
	}
}

// ---------------------------------------------------------------------------
// Independent numeric checks
// ---------------------------------------------------------------------------

func TestBatchZTestRecomputesFromFacts(t *testing.T) {
	reg := statsBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := statsBatchSample(t, reg, "stat.infer.ztest", i)
		xbar := statsFactFloat(t, p, "xbar")
		mu := statsFactFloat(t, p, "mu")
		sigma := statsFactFloat(t, p, "sigma")
		n := statsFactInt(t, p, "n")
		want := (xbar - mu) / (sigma / math.Sqrt(float64(n)))
		answer, err := strconv.ParseFloat(p.Answer, 64)
		if err != nil {
			t.Fatalf("ztest: non-numeric answer %q", p.Answer)
		}
		if math.Abs(answer-want) > 1e-9 {
			t.Fatalf("ztest: xbar=%v mu=%v sigma=%v n=%d answers %v, want %v", xbar, mu, sigma, n, answer, want)
		}
	}
}

func TestBatchTTestRecomputesFromFacts(t *testing.T) {
	reg := statsBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := statsBatchSample(t, reg, "stat.infer.ttest", i)
		xbar := statsFactFloat(t, p, "xbar")
		mu := statsFactFloat(t, p, "mu")
		s := statsFactFloat(t, p, "s")
		n := statsFactInt(t, p, "n")
		want := (xbar - mu) / (s / math.Sqrt(float64(n)))
		answer, err := strconv.ParseFloat(p.Answer, 64)
		if err != nil {
			t.Fatalf("ttest: non-numeric answer %q", p.Answer)
		}
		if math.Abs(answer-want) > 1e-9 {
			t.Fatalf("ttest: xbar=%v mu=%v s=%v n=%d answers %v, want %v", xbar, mu, s, n, answer, want)
		}
	}
}

func TestBatchProportionRecomputesFromFacts(t *testing.T) {
	reg := statsBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := statsBatchSample(t, reg, "stat.infer.proportion", i)
		phat := statsFactFloat(t, p, "p")
		n := statsFactInt(t, p, "n")
		want := math.Sqrt(phat * (1 - phat) / float64(n))
		answer, err := strconv.ParseFloat(p.Answer, 64)
		if err != nil {
			t.Fatalf("proportion: non-numeric answer %q", p.Answer)
		}
		if math.Abs(answer-want) > 1e-9 {
			t.Fatalf("proportion: p=%v n=%d answers %v, want %v", phat, n, answer, want)
		}
	}
}

func TestBatchChiSquareRecomputesFromFacts(t *testing.T) {
	reg := statsBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := statsBatchSample(t, reg, "stat.test.chi_square", i)
		obs := statsFactInts(t, p, "observed")
		exp := statsFactInts(t, p, "expected")
		if len(obs) != len(exp) {
			t.Fatalf("chi-square: %d observed vs %d expected", len(obs), len(exp))
		}
		want := 0.0
		for j := range obs {
			if exp[j] == 0 {
				t.Fatalf("chi-square: zero expected count")
			}
			d := float64(obs[j] - exp[j])
			want += d * d / float64(exp[j])
		}
		answer, err := strconv.ParseFloat(p.Answer, 64)
		if err != nil {
			t.Fatalf("chi-square: non-numeric answer %q", p.Answer)
		}
		if math.Abs(answer-want) > 1e-9 {
			t.Fatalf("chi-square: obs=%v exp=%v answers %v, want %v", obs, exp, answer, want)
		}
	}
}

func TestBatchANOVARecomputesFromFacts(t *testing.T) {
	reg := statsBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := statsBatchSample(t, reg, "stat.test.anova", i)
		groupParts := strings.Split(p.Facts["groups"], ";")
		groups := make([][]int, len(groupParts))
		for g, gp := range groupParts {
			cells := strings.Split(gp, ",")
			groups[g] = make([]int, len(cells))
			for c, s := range cells {
				v, err := strconv.Atoi(strings.TrimSpace(s))
				if err != nil {
					t.Fatalf("anova: bad value %q", s)
				}
				groups[g][c] = v
			}
		}
		want := statsAnovaF(groups)
		answer, err := strconv.ParseFloat(p.Answer, 64)
		if err != nil {
			t.Fatalf("anova: non-numeric answer %q", p.Answer)
		}
		if math.Abs(answer-want) > 1e-9 {
			t.Fatalf("anova: groups %v answers %v, want %v", groups, answer, want)
		}
	}
}

func TestBatchR2RecomputesFromFacts(t *testing.T) {
	reg := statsBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := statsBatchSample(t, reg, "stat.reg.r2", i)
		sst := statsFactFloat(t, p, "sst")
		sse := statsFactFloat(t, p, "sse")
		if sst == 0 {
			t.Fatalf("r2: zero total sum of squares")
		}
		want := 1 - sse/sst
		answer, err := strconv.ParseFloat(p.Answer, 64)
		if err != nil {
			t.Fatalf("r2: non-numeric answer %q", p.Answer)
		}
		if math.Abs(answer-want) > 1e-9 {
			t.Fatalf("r2: sst=%v sse=%v answers %v, want %v", sst, sse, answer, want)
		}
	}
}

func TestBatchHistogramRecomputesFromFacts(t *testing.T) {
	reg := statsBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := statsBatchSample(t, reg, "stat.data.histogram", i)
		values := statsFactInts(t, p, "values")
		low := statsFactInt(t, p, "low")
		high := statsFactInt(t, p, "high")
		want := 0
		for _, v := range values {
			if v >= low && v < high {
				want++
			}
		}
		answer, err := strconv.Atoi(p.Answer)
		if err != nil {
			t.Fatalf("histogram: non-integer answer %q", p.Answer)
		}
		if answer != want {
			t.Fatalf("histogram: values %v in [%d,%d) answers %d, want %d", values, low, high, answer, want)
		}
	}
}

func TestBatchNumericAnswersParseAsDecimals(t *testing.T) {
	reg := statsBatchRegistry(t)
	for _, id := range statsBatchNumericIDs {
		for i := 0; i < 120; i++ {
			p := statsBatchSample(t, reg, id, i)
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
	reg := statsBatchRegistry(t)
	const samples = 150
	for _, id := range statsBatchIDs {
		answers := map[string]bool{}
		questions := map[string]bool{}
		allYesNo := true
		for i := 0; i < samples; i++ {
			p, err := reg.GenerateContext(id, generator.GeneratorContext{Difficulty: 0.12, Seed: int64(112000 + i)})
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
