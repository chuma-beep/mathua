package machinelearning

// Tests for the 2026-10-18 machine-learning batch.
//
// Every concept is graded through the real router with the grading_type the
// corpus ships, so a generator and its declared type cannot drift apart. The
// multiple-choice concepts are checked the way the algebra batch checks its
// own: the question is parsed into options, the correct answer must appear
// exactly once, and every other option must be rejected by the router.
//
// The seven numeric concepts are recomputed from the facts the generator
// publishes. The sigmoid is evaluated, the absolute error is summed, the
// confusion-matrix cell is picked by the asked-for count, F1 is computed from
// precision and recall, the centroid is averaged, the z-score is divided, and
// the Gini impurity is recomputed from the class counts. A generator that
// quietly changed a table entry would fail here.

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

var mlBatchIDs = []string{
	"ml.activation.sigmoid",
	"ml.logistic.regression",
	"ml.grad.vanishing",
	"ml.loss.mae",
	"ml.eval.confusion",
	"ml.eval.f1",
	"ml.cluster.kmeans",
	"ml.data.scaling",
	"ml.data.onehot",
	"ml.data.imbalance",
	"ml.tree.gini",
	"ml.tree.ensemble",
}

var mlBatchNumericIDs = []string{
	"ml.activation.sigmoid",
	"ml.loss.mae",
	"ml.eval.confusion",
	"ml.eval.f1",
	"ml.cluster.kmeans",
	"ml.data.scaling",
	"ml.tree.gini",
	"ml.logistic.regression",
	"ml.grad.vanishing",
	"ml.data.imbalance",
	"ml.tree.ensemble",
}

type mlCorpusConcept struct {
	ID          string `json:"id"`
	GradingType string `json:"grading_type"`
}

func mlShippedTypes(t *testing.T) map[string]string {
	t.Helper()
	raw, err := os.ReadFile("../../../data/concepts/machine_learning.json")
	if err != nil {
		t.Fatalf("read corpus: %v", err)
	}
	var rows []mlCorpusConcept
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

func mlRegistry(t *testing.T) *generator.Registry {
	t.Helper()
	reg := generator.NewRegistry()
	Register(reg)
	return reg
}

func mlSample(t *testing.T, reg *generator.Registry, id string, i int) generator.Problem {
	t.Helper()
	p, err := reg.GenerateContext(id, generator.GeneratorContext{
		Difficulty: float64(i%20) / 19,
		Seed:       int64(81000 + i),
	})
	if err != nil {
		t.Fatalf("%s: generate: %v", id, err)
	}
	return p
}

var mlOptionLine = regexp.MustCompile(`\n\n([A-E])\) ([^\n]*)`)

func mlOptionsOf(q string) []string {
	var out []string
	for _, m := range mlOptionLine.FindAllStringSubmatch(q, -1) {
		out = append(out, m[2])
	}
	return out
}

func mlFactInt(t *testing.T, p generator.Problem, key string) int {
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
	types := mlShippedTypes(t)
	reg := mlRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	for _, id := range mlBatchIDs {
		typ, ok := types[id]
		if !ok || typ == "" {
			t.Fatalf("%s: no grading_type in the corpus", id)
		}
		if typ != "multiple_choice" && typ != "numeric" {
			t.Fatalf("%s: unexpected grading_type %q", id, typ)
		}
		for i := 0; i < 120; i++ {
			p := mlSample(t, reg, id, i)
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
	types := mlShippedTypes(t)
	reg := mlRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	checked := 0
	for _, id := range mlBatchIDs {
		if types[id] != "multiple_choice" {
			continue
		}
		for i := 0; i < 40; i++ {
			p := mlSample(t, reg, id, i)
			opts := mlOptionsOf(p.Question)
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
		t.Skip("the batch now serves computed count problems")
	}
}

// ---------------------------------------------------------------------------
// Independent numeric checks
// ---------------------------------------------------------------------------

func TestBatchSigmoidRecomputesFromFacts(t *testing.T) {
	reg := mlRegistry(t)
	seenValue, seenLimit := false, false
	for i := 0; i < 300; i++ {
		p := mlSample(t, reg, "ml.activation.sigmoid", i)
		answer, err := strconv.ParseFloat(p.Answer, 64)
		if err != nil {
			t.Fatalf("sigmoid: non-numeric answer %q", p.Answer)
		}
		switch p.Facts["kind"] {
		case "value":
			seenValue = true
			x, err := strconv.ParseFloat(p.Facts["x"], 64)
			if err != nil {
				t.Fatalf("sigmoid: bad x fact %q", p.Facts["x"])
			}
			want := 1 / (1 + math.Exp(-x))
			if math.Abs(answer-want) > 1e-9 {
				t.Fatalf("sigmoid: sigma(%g) answers %v, want %v", x, answer, want)
			}
		case "limit":
			seenLimit = true
			if answer != 1 {
				t.Fatalf("sigmoid: the limit answers %v, want 1", answer)
			}
		default:
			t.Fatalf("sigmoid: unknown kind %q", p.Facts["kind"])
		}
	}
	if !seenValue || !seenLimit {
		t.Fatalf("sigmoid: branches not both exercised (value=%v limit=%v)", seenValue, seenLimit)
	}
}

func TestBatchMAERecomputesFromFacts(t *testing.T) {
	reg := mlRegistry(t)
	for i := 0; i < 200; i++ {
		p := mlSample(t, reg, "ml.loss.mae", i)
		parts := strings.Split(p.Facts["residuals"], ",")
		if len(parts) == 0 {
			t.Fatalf("mae: no residuals")
		}
		sum := 0
		for _, s := range parts {
			n, err := strconv.Atoi(s)
			if err != nil {
				t.Fatalf("mae: bad residual %q", s)
			}
			if n < 0 {
				n = -n
			}
			sum += n
		}
		want := float64(sum) / float64(len(parts))
		answer, err := strconv.ParseFloat(p.Answer, 64)
		if err != nil {
			t.Fatalf("mae: non-numeric answer %q", p.Answer)
		}
		if math.Abs(answer-want) > 1e-9 {
			t.Fatalf("mae: residuals %v answer %v, want %v", parts, answer, want)
		}
	}
}

func TestBatchConfusionRecomputesFromFacts(t *testing.T) {
	reg := mlRegistry(t)
	seen := map[string]bool{}
	for i := 0; i < 250; i++ {
		p := mlSample(t, reg, "ml.eval.confusion", i)
		tp := mlFactInt(t, p, "tp")
		fp := mlFactInt(t, p, "fp")
		fn := mlFactInt(t, p, "fn")
		tn := mlFactInt(t, p, "tn")
		ask := p.Facts["ask"]
		seen[ask] = true
		var want int
		switch ask {
		case "total":
			want = tp + fp + fn + tn
		case "actual_pos":
			want = tp + fn
		case "predicted_pos":
			want = tp + fp
		case "correct":
			want = tp + tn
		case "actual_neg":
			want = fp + tn
		default:
			t.Fatalf("confusion: unknown ask %q", ask)
		}
		if p.Answer != strconv.Itoa(want) {
			t.Fatalf("confusion: %s with tp=%d fp=%d fn=%d tn=%d answers %s, want %d",
				ask, tp, fp, fn, tn, p.Answer, want)
		}
	}
	if len(seen) != 5 {
		t.Fatalf("confusion: expected all five asks in the pool, saw %v", seen)
	}
}

func TestBatchF1RecomputesFromFacts(t *testing.T) {
	reg := mlRegistry(t)
	for i := 0; i < 200; i++ {
		p := mlSample(t, reg, "ml.eval.f1", i)
		precision, err := strconv.ParseFloat(p.Facts["precision"], 64)
		if err != nil {
			t.Fatalf("f1: bad precision fact %q", p.Facts["precision"])
		}
		recall, err := strconv.ParseFloat(p.Facts["recall"], 64)
		if err != nil {
			t.Fatalf("f1: bad recall fact %q", p.Facts["recall"])
		}
		if precision+recall == 0 {
			t.Fatalf("f1: both precision and recall are zero")
		}
		want := 2 * precision * recall / (precision + recall)
		answer, err := strconv.ParseFloat(p.Answer, 64)
		if err != nil {
			t.Fatalf("f1: non-numeric answer %q", p.Answer)
		}
		if math.Abs(answer-want) > 1e-9 {
			t.Fatalf("f1: precision %v recall %v answers %v, want %v", precision, recall, answer, want)
		}
	}
}

func TestBatchKMeansRecomputesFromFacts(t *testing.T) {
	reg := mlRegistry(t)
	for i := 0; i < 200; i++ {
		p := mlSample(t, reg, "ml.cluster.kmeans", i)
		pts := strings.Split(p.Facts["points"], ";")
		if len(pts) == 0 {
			t.Fatalf("kmeans: no points")
		}
		sum := 0
		axis := 0
		if p.Facts["axis"] == "y" {
			axis = 1
		}
		for _, pt := range pts {
			coords := strings.Split(pt, ",")
			if len(coords) != 2 {
				t.Fatalf("kmeans: bad point %q", pt)
			}
			v, err := strconv.Atoi(coords[axis])
			if err != nil {
				t.Fatalf("kmeans: bad coordinate %q", coords[axis])
			}
			sum += v
		}
		want := float64(sum) / float64(len(pts))
		answer, err := strconv.ParseFloat(p.Answer, 64)
		if err != nil {
			t.Fatalf("kmeans: non-numeric answer %q", p.Answer)
		}
		if math.Abs(answer-want) > 1e-9 {
			t.Fatalf("kmeans: %s-coordinate of %v answers %v, want %v", p.Facts["axis"], pts, answer, want)
		}
	}
}

func TestBatchScalingRecomputesFromFacts(t *testing.T) {
	reg := mlRegistry(t)
	for i := 0; i < 200; i++ {
		p := mlSample(t, reg, "ml.data.scaling", i)
		x := mlFactInt(t, p, "x")
		mean := mlFactInt(t, p, "mean")
		sd := mlFactInt(t, p, "sd")
		if sd == 0 {
			t.Fatalf("scaling: standard deviation is zero")
		}
		want := float64(x-mean) / float64(sd)
		answer, err := strconv.ParseFloat(p.Answer, 64)
		if err != nil {
			t.Fatalf("scaling: non-numeric answer %q", p.Answer)
		}
		if math.Abs(answer-want) > 1e-9 {
			t.Fatalf("scaling: (%d - %d)/%d answers %v, want %v", x, mean, sd, answer, want)
		}
	}
}

func TestBatchGiniRecomputesFromFacts(t *testing.T) {
	reg := mlRegistry(t)
	for i := 0; i < 200; i++ {
		p := mlSample(t, reg, "ml.tree.gini", i)
		parts := strings.Split(p.Facts["counts"], ",")
		total := 0
		counts := make([]int, 0, len(parts))
		for _, s := range parts {
			n, err := strconv.Atoi(s)
			if err != nil || n < 0 {
				t.Fatalf("gini: bad count %q", s)
			}
			counts = append(counts, n)
			total += n
		}
		if total == 0 {
			t.Fatalf("gini: node has no examples")
		}
		want := 1.0
		for _, n := range counts {
			prob := float64(n) / float64(total)
			want -= prob * prob
		}
		answer, err := strconv.ParseFloat(p.Answer, 64)
		if err != nil {
			t.Fatalf("gini: non-numeric answer %q", p.Answer)
		}
		if math.Abs(answer-want) > 1e-9 {
			t.Fatalf("gini: counts %v answer %v, want %v", counts, answer, want)
		}
	}
}

func TestBatchNumericAnswersParseAsDecimals(t *testing.T) {
	reg := mlRegistry(t)
	for _, id := range mlBatchNumericIDs {
		for i := 0; i < 120; i++ {
			p := mlSample(t, reg, id, i)
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
	reg := mlRegistry(t)
	const samples = 150
	for _, id := range mlBatchIDs {
		answers := map[string]bool{}
		questions := map[string]bool{}
		allYesNo := true
		for i := 0; i < samples; i++ {
			p, err := reg.GenerateContext(id, generator.GeneratorContext{Difficulty: 0.12, Seed: int64(91000 + i)})
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
