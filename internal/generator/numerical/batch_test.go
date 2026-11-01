package numerical

// Tests for the first numerical-analysis batch.
//
// Every concept is graded through the real router with the grading_type the
// corpus ships. The batch now serves computed answers, so the multiple-choice
// check skips when no recognition concept remains.
//
// The eleven numeric concepts recompute their answers from the published
// parameters: the relative error from the true and approximate values, the
// cancellation amplification from the two numbers, the condition number from the
// exponent, the bisection step count from the width and tolerance, a fixed-point
// iterate, a Newton and a secant step, the order from three successive errors,
// a Horner evaluation, a Lagrange interpolation and a divided difference.

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

type numaCorpusConcept struct {
	ID          string `json:"id"`
	GradingType string `json:"grading_type"`
}

var numaBatchIDs = []string{
	"numa.errors.absolute",
	"numa.errors.roundoff",
	"numa.errors.cancellation",
	"numa.errors.conditioning",
	"numa.errors.stability",
	"numa.errors.taylor",
	"numa.roots.bisection",
	"numa.roots.fixed_point",
	"numa.roots.newton",
	"numa.roots.secant",
	"numa.roots.false_position",
	"numa.roots.order",
	"numa.roots.multiple",
	"numa.roots.bracketing",
	"numa.roots.horner",
	"numa.interp.lagrange",
	"numa.interp.divided_differences",
	"numa.interp.error",
	"numa.interp.chebyshev",
	"numa.interp.splines",
}

var numaMCIDs = []string{}

func numaShippedTypes(t *testing.T) map[string]string {
	t.Helper()
	raw, err := os.ReadFile("../../../data/concepts/numerical_analysis.json")
	if err != nil {
		t.Fatalf("read corpus: %v", err)
	}
	var rows []numaCorpusConcept
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

func numaRegistry(t *testing.T) *generator.Registry {
	t.Helper()
	reg := generator.NewRegistry()
	Register(reg)
	return reg
}

func numaSample(t *testing.T, reg *generator.Registry, id string, i int) generator.Problem {
	t.Helper()
	p, err := reg.GenerateContext(id, generator.GeneratorContext{
		Difficulty: float64(i%20) / 19,
		Seed:       int64(271000 + i),
	})
	if err != nil {
		t.Fatalf("%s: generate: %v", id, err)
	}
	return p
}

var numaOptionLine = regexp.MustCompile(`\n\n([A-E])\) ([^\n]*)`)

func numaOptionsOf(q string) []string {
	var out []string
	for _, m := range numaOptionLine.FindAllStringSubmatch(q, -1) {
		out = append(out, m[2])
	}
	return out
}

func numaFactInt(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer", p.Question, key, p.Facts[key])
	}
	return v
}

func numaFactFloat(t *testing.T, p generator.Problem, key string) float64 {
	t.Helper()
	v, err := strconv.ParseFloat(p.Facts[key], 64)
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not a number", p.Question, key, p.Facts[key])
	}
	return v
}

func numaIntList(t *testing.T, p generator.Problem, key string) []int {
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

func TestNumaAnswersGradeWithTheirShippedType(t *testing.T) {
	types := numaShippedTypes(t)
	reg := numaRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	for _, id := range numaBatchIDs {
		typ, ok := types[id]
		if !ok || typ == "" {
			t.Fatalf("%s: no grading_type in the corpus", id)
		}
		if typ != "multiple_choice" && typ != "numeric" {
			t.Fatalf("%s: unexpected grading_type %q", id, typ)
		}
		for i := 0; i < 120; i++ {
			p := numaSample(t, reg, id, i)
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

func TestNumaMultipleChoiceRejectsEveryDistractor(t *testing.T) {
	if len(numaMCIDs) == 0 {
		t.Skip("the numerical batch now serves computed count problems")
	}
	types := numaShippedTypes(t)
	reg := numaRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	for _, id := range numaMCIDs {
		if types[id] != "multiple_choice" {
			t.Fatalf("%s: expected multiple_choice, corpus says %q", id, types[id])
		}
		for i := 0; i < 40; i++ {
			p := numaSample(t, reg, id, i)
			opts := numaOptionsOf(p.Question)
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

func TestNumaAbsoluteErrorRecomputes(t *testing.T) {
	reg := numaRegistry(t)
	for i := 0; i < 200; i++ {
		p := numaSample(t, reg, "numa.errors.absolute", i)
		tr := numaFactInt(t, p, "true")
		ap := numaFactInt(t, p, "approx")
		want := decimal6(math.Abs(float64(ap-tr)) / math.Abs(float64(tr)))
		if p.Answer != want {
			t.Fatalf("absolute error: true=%d approx=%d answers %s, want %s", tr, ap, p.Answer, want)
		}
	}
}

func TestNumaCancellationRecomputes(t *testing.T) {
	reg := numaRegistry(t)
	for i := 0; i < 200; i++ {
		p := numaSample(t, reg, "numa.errors.cancellation", i)
		x := numaFactInt(t, p, "x")
		y := numaFactInt(t, p, "y")
		want := decimal6(float64(absInt(x)+absInt(y)) / float64(absInt(x-y)))
		if p.Answer != want {
			t.Fatalf("cancellation: x=%d y=%d answers %s, want %s", x, y, p.Answer, want)
		}
	}
}

func TestNumaConditioningRecomputes(t *testing.T) {
	reg := numaRegistry(t)
	for i := 0; i < 200; i++ {
		p := numaSample(t, reg, "numa.errors.conditioning", i)
		n := numaFactInt(t, p, "n")
		if p.Answer != decimal6(float64(n)) {
			t.Fatalf("conditioning: n=%d answers %s, want %s", n, p.Answer, decimal6(float64(n)))
		}
	}
}

func TestNumaBisectionRecomputes(t *testing.T) {
	reg := numaRegistry(t)
	for i := 0; i < 200; i++ {
		p := numaSample(t, reg, "numa.roots.bisection", i)
		w := numaFactFloat(t, p, "width")
		tol := numaFactFloat(t, p, "tol")
		want := strconv.Itoa(int(math.Ceil(math.Log2(w / tol))))
		if p.Answer != want {
			t.Fatalf("bisection: width=%v tol=%v answers %s, want %s", w, tol, p.Answer, want)
		}
	}
}

func TestNumaFixedPointRecomputes(t *testing.T) {
	reg := numaRegistry(t)
	for i := 0; i < 200; i++ {
		p := numaSample(t, reg, "numa.roots.fixed_point", i)
		a := numaFactInt(t, p, "a")
		bn := numaFactInt(t, p, "bn")
		bd := numaFactInt(t, p, "bd")
		x0 := numaFactInt(t, p, "x0")
		n := numaFactInt(t, p, "n")
		x := float64(x0)
		b := float64(bn) / float64(bd)
		for k := 0; k < n; k++ {
			x = float64(a) + b*x
		}
		if p.Answer != decimal6(x) {
			t.Fatalf("fixed point: a=%d b=%d/%d x0=%d n=%d answers %s, want %s", a, bn, bd, x0, n, p.Answer, decimal6(x))
		}
	}
}

func TestNumaNewtonRecomputes(t *testing.T) {
	reg := numaRegistry(t)
	for i := 0; i < 200; i++ {
		p := numaSample(t, reg, "numa.roots.newton", i)
		c := numaFactFloat(t, p, "c")
		x0 := numaFactFloat(t, p, "x0")
		want := decimal6((x0 + c/x0) / 2)
		if p.Answer != want {
			t.Fatalf("newton: c=%v x0=%v answers %s, want %s", c, x0, p.Answer, want)
		}
	}
}

func TestNumaSecantRecomputes(t *testing.T) {
	reg := numaRegistry(t)
	for i := 0; i < 200; i++ {
		p := numaSample(t, reg, "numa.roots.secant", i)
		c := numaFactFloat(t, p, "c")
		x0 := numaFactFloat(t, p, "x0")
		x1 := numaFactFloat(t, p, "x1")
		f := func(x float64) float64 { return x*x - c }
		want := decimal6(x1 - f(x1)*(x1-x0)/(f(x1)-f(x0)))
		if p.Answer != want {
			t.Fatalf("secant: c=%v x0=%v x1=%v answers %s, want %s", c, x0, x1, p.Answer, want)
		}
	}
}

func TestNumaOrderRecomputes(t *testing.T) {
	reg := numaRegistry(t)
	for i := 0; i < 200; i++ {
		p := numaSample(t, reg, "numa.roots.order", i)
		a := numaFactInt(t, p, "a")
		b := numaFactInt(t, p, "b")
		c := numaFactInt(t, p, "c")
		want := decimal6(float64(c-b) / float64(b-a))
		if p.Answer != want {
			t.Fatalf("order: a=%d b=%d c=%d answers %s, want %s", a, b, c, p.Answer, want)
		}
	}
}

func TestNumaHornerRecomputes(t *testing.T) {
	reg := numaRegistry(t)
	for i := 0; i < 200; i++ {
		p := numaSample(t, reg, "numa.roots.horner", i)
		coeffs := numaIntList(t, p, "coeffs")
		x := numaFactInt(t, p, "x")
		val := 0
		for k := len(coeffs) - 1; k >= 0; k-- {
			val = val*x + coeffs[k]
		}
		if p.Answer != strconv.Itoa(val) {
			t.Fatalf("horner: coeffs %v x=%d answers %s, want %d", coeffs, x, p.Answer, val)
		}
	}
}

func TestNumaLagrangeRecomputes(t *testing.T) {
	reg := numaRegistry(t)
	for i := 0; i < 200; i++ {
		p := numaSample(t, reg, "numa.interp.lagrange", i)
		y0 := numaFactInt(t, p, "y0")
		y1 := numaFactInt(t, p, "y1")
		y2 := numaFactInt(t, p, "y2")
		q := float64(numaFactInt(t, p, "q"))
		l0 := (q - 1) * (q - 2) / 2
		l1 := -q * (q - 2)
		l2 := q * (q - 1) / 2
		val := float64(y0)*l0 + float64(y1)*l1 + float64(y2)*l2
		if p.Answer != strconv.Itoa(int(math.Round(val))) {
			t.Fatalf("lagrange: y=%d/%d/%d q=%v answers %s, want %v", y0, y1, y2, q, p.Answer, val)
		}
	}
}

func TestNumaDividedDifferenceRecomputes(t *testing.T) {
	reg := numaRegistry(t)
	for i := 0; i < 200; i++ {
		p := numaSample(t, reg, "numa.interp.divided_differences", i)
		x0 := numaFactInt(t, p, "x0")
		y0 := numaFactInt(t, p, "y0")
		x1 := numaFactInt(t, p, "x1")
		y1 := numaFactInt(t, p, "y1")
		want := decimal6(float64(y1-y0) / float64(x1-x0))
		if p.Answer != want {
			t.Fatalf("divided difference: (%d,%d),(%d,%d) answers %s, want %s", x0, y0, x1, y1, p.Answer, want)
		}
	}
}
