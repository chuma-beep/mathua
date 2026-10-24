package odes

// Tests for the 2026-10-18 differential-equations batch.
//
// Every concept is graded through the real router with the grading_type the
// corpus ships. The four multiple-choice concepts are checked the way the other
// batches check theirs: the question is parsed into options, the correct answer
// must appear exactly once, and every other option must be rejected.
//
// The six numeric concepts recompute their answers from the published facts:
// the Cauchy-Euler coefficients, the series kind and rate, the spring constants
// and initial conditions, the tank volume and concentration, the Lotka-Volterra
// coefficients, and the Heun step.

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

var odeBatchIDs = []string{
	"ode.first_order.bernoulli",
	"ode.second_order.cauchy_euler",
	"ode.second_order.reduction",
	"ode.second_order.power_series",
	"ode.systems.matrix_exp",
	"ode.transforms.step",
	"ode.second_order.spring",
	"ode.first_order.mixing",
	"ode.systems.predator_prey",
	"ode.numerical.heun",
}

var odeBatchNumericIDs = []string{
	"ode.second_order.cauchy_euler",
	"ode.second_order.power_series",
	"ode.second_order.spring",
	"ode.first_order.mixing",
	"ode.systems.predator_prey",
	"ode.numerical.heun",
}

type odeCorpusConcept struct {
	ID          string `json:"id"`
	GradingType string `json:"grading_type"`
}

func odeShippedTypes(t *testing.T) map[string]string {
	t.Helper()
	raw, err := os.ReadFile("../../../data/concepts/differential_equations.json")
	if err != nil {
		t.Fatalf("read corpus: %v", err)
	}
	var rows []odeCorpusConcept
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

func odeBatchRegistry(t *testing.T) *generator.Registry {
	t.Helper()
	reg := generator.NewRegistry()
	Register(reg)
	return reg
}

func odeBatchSample(t *testing.T, reg *generator.Registry, id string, i int) generator.Problem {
	t.Helper()
	p, err := reg.GenerateContext(id, generator.GeneratorContext{
		Difficulty: float64(i%20) / 19,
		Seed:       int64(141000 + i),
	})
	if err != nil {
		t.Fatalf("%s: generate: %v", id, err)
	}
	return p
}

var odeOptionLine = regexp.MustCompile(`\n\n([A-E])\) ([^\n]*)`)

func odeOptionsOf(q string) []string {
	var out []string
	for _, m := range odeOptionLine.FindAllStringSubmatch(q, -1) {
		out = append(out, m[2])
	}
	return out
}

func odeFactInt(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer", p.Question, key, p.Facts[key])
	}
	return v
}

func odeFactFloat(t *testing.T, p generator.Problem, key string) float64 {
	t.Helper()
	v, err := strconv.ParseFloat(p.Facts[key], 64)
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not a number", p.Question, key, p.Facts[key])
	}
	return v
}

// ---------------------------------------------------------------------------
// Router contracts
// ---------------------------------------------------------------------------

func TestBatchAnswersGradeWithTheirShippedType(t *testing.T) {
	types := odeShippedTypes(t)
	reg := odeBatchRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	for _, id := range odeBatchIDs {
		typ, ok := types[id]
		if !ok || typ == "" {
			t.Fatalf("%s: no grading_type in the corpus", id)
		}
		if typ != "multiple_choice" && typ != "numeric" {
			t.Fatalf("%s: unexpected grading_type %q", id, typ)
		}
		for i := 0; i < 120; i++ {
			p := odeBatchSample(t, reg, id, i)
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
	types := odeShippedTypes(t)
	reg := odeBatchRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	checked := 0
	for _, id := range odeBatchIDs {
		if types[id] != "multiple_choice" {
			continue
		}
		for i := 0; i < 40; i++ {
			p := odeBatchSample(t, reg, id, i)
			opts := odeOptionsOf(p.Question)
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
		t.Fatal("no multiple_choice concepts were checked")
	}
}

// ---------------------------------------------------------------------------
// Independent numeric checks
// ---------------------------------------------------------------------------

func TestBatchCauchyEulerRecomputes(t *testing.T) {
	reg := odeBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := odeBatchSample(t, reg, "ode.second_order.cauchy_euler", i)
		a := float64(odeFactInt(t, p, "a"))
		b := float64(odeFactInt(t, p, "b"))
		bb := a - 1
		disc := bb*bb - 4*b
		if disc < 0 {
			t.Fatalf("cauchy-euler: complex roots for a=%v b=%v", a, b)
		}
		root1 := (-bb + math.Sqrt(disc)) / 2
		root2 := (-bb - math.Sqrt(disc)) / 2
		want := math.Max(root1, root2)
		answer, err := strconv.ParseFloat(p.Answer, 64)
		if err != nil {
			t.Fatalf("cauchy-euler: non-numeric answer %q", p.Answer)
		}
		if math.Abs(answer-want) > 1e-9 {
			t.Fatalf("cauchy-euler: a=%v b=%v answers %v, want %v", a, b, answer, want)
		}
	}
}

func TestBatchPowerSeriesRecomputes(t *testing.T) {
	reg := odeBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := odeBatchSample(t, reg, "ode.second_order.power_series", i)
		k := float64(odeFactInt(t, p, "k"))
		y0 := float64(odeFactInt(t, p, "y0"))
		var want float64
		switch p.Facts["kind"] {
		case "first":
			want = y0 * k * k / 2
		case "cos":
			want = -k * k / 2
		case "cosh":
			want = k * k / 2
		default:
			t.Fatalf("power series: unknown kind %q", p.Facts["kind"])
		}
		answer, err := strconv.ParseFloat(p.Answer, 64)
		if err != nil {
			t.Fatalf("power series: non-numeric answer %q", p.Answer)
		}
		if math.Abs(answer-want) > 1e-9 {
			t.Fatalf("power series: %s k=%v y0=%v answers %v, want %v", p.Facts["kind"], k, y0, answer, want)
		}
	}
}

func TestBatchSpringRecomputes(t *testing.T) {
	reg := odeBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := odeBatchSample(t, reg, "ode.second_order.spring", i)
		omega := math.Sqrt(float64(odeFactInt(t, p, "omega2")))
		y0 := float64(odeFactInt(t, p, "y0"))
		v0 := float64(odeFactInt(t, p, "v0"))
		var tval float64
		switch p.Facts["t"] {
		case "pi/4":
			tval = math.Pi / 4
		case "pi/6":
			tval = math.Pi / 6
		case "pi":
			tval = math.Pi
		default:
			t.Fatalf("spring: unknown time %q", p.Facts["t"])
		}
		want := y0*math.Cos(omega*tval) + (v0/omega)*math.Sin(omega*tval)
		answer, err := strconv.ParseFloat(p.Answer, 64)
		if err != nil {
			t.Fatalf("spring: non-numeric answer %q", p.Answer)
		}
		if math.Abs(answer-want) > 1e-9 {
			t.Fatalf("spring: omega2=%v y0=%v v0=%v t=%s answers %v, want %v",
				omega*omega, y0, v0, p.Facts["t"], answer, want)
		}
	}
}

func TestBatchMixingRecomputes(t *testing.T) {
	reg := odeBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := odeBatchSample(t, reg, "ode.first_order.mixing", i)
		v := float64(odeFactInt(t, p, "v"))
		c := odeFactFloat(t, p, "c")
		want := v * c
		answer, err := strconv.ParseFloat(p.Answer, 64)
		if err != nil {
			t.Fatalf("mixing: non-numeric answer %q", p.Answer)
		}
		if math.Abs(answer-want) > 1e-9 {
			t.Fatalf("mixing: v=%v c=%v answers %v, want %v", v, c, answer, want)
		}
	}
}

func TestBatchPredatorPreyRecomputes(t *testing.T) {
	reg := odeBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := odeBatchSample(t, reg, "ode.systems.predator_prey", i)
		a := odeFactFloat(t, p, "a")
		b := odeFactFloat(t, p, "b")
		c := odeFactFloat(t, p, "c")
		d := odeFactFloat(t, p, "d")
		var want float64
		switch p.Facts["ask"] {
		case "x":
			want = c / d
		case "y":
			want = a / b
		default:
			t.Fatalf("predator-prey: unknown ask %q", p.Facts["ask"])
		}
		answer, err := strconv.ParseFloat(p.Answer, 64)
		if err != nil {
			t.Fatalf("predator-prey: non-numeric answer %q", p.Answer)
		}
		if math.Abs(answer-want) > 1e-9 {
			t.Fatalf("predator-prey: a=%v b=%v c=%v d=%v ask=%s answers %v, want %v",
				a, b, c, d, p.Facts["ask"], answer, want)
		}
	}
}

func TestBatchHeunRecomputes(t *testing.T) {
	reg := odeBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := odeBatchSample(t, reg, "ode.numerical.heun", i)
		y0 := odeFactFloat(t, p, "y0")
		t0 := odeFactFloat(t, p, "t0")
		h := odeFactFloat(t, p, "h")
		f := func(tv, y float64) float64 {
			switch p.Facts["function"] {
			case "y":
				return y
			case "double":
				return 2 * y
			case "square":
				return y * y
			case "tplusy":
				return tv + y
			default:
				t.Fatalf("heun: unknown function %q", p.Facts["function"])
				return 0
			}
		}
		predictor := y0 + h*f(t0, y0)
		want := y0 + h/2*(f(t0, y0)+f(t0+h, predictor))
		answer, err := strconv.ParseFloat(p.Answer, 64)
		if err != nil {
			t.Fatalf("heun: non-numeric answer %q", p.Answer)
		}
		// The published answer is rounded to six decimal places so it reads as
		// 1.052563 rather than 1.052562500000001, so compare at that precision.
		if math.Abs(answer-want) > 1e-6 {
			t.Fatalf("heun: function=%s y0=%v t0=%v h=%v answers %v, want %v",
				p.Facts["function"], y0, t0, h, answer, want)
		}
	}
}

func TestBatchNumericAnswersParseAsDecimals(t *testing.T) {
	reg := odeBatchRegistry(t)
	for _, id := range odeBatchNumericIDs {
		for i := 0; i < 120; i++ {
			p := odeBatchSample(t, reg, id, i)
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
	reg := odeBatchRegistry(t)
	const samples = 150
	for _, id := range odeBatchIDs {
		answers := map[string]bool{}
		questions := map[string]bool{}
		allYesNo := true
		for i := 0; i < samples; i++ {
			p, err := reg.GenerateContext(id, generator.GeneratorContext{Difficulty: 0.12, Seed: int64(142000 + i)})
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
