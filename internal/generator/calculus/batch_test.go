package calculus

// Tests for the 2026-10-18 calculus batch.
//
// Every concept is graded through the real router with the grading_type the
// corpus ships. The batch has no multiple-choice concepts: seven numeric, one
// expression, one symbolic and one tuple. Each generator publishes the facts
// its answer follows from, and the tests recompute the answer from those facts:
// the hyperbolic rule, the parametric derivatives, the finite-difference method,
// the polar curve, the three curve derivatives, the two frame derivatives, the
// gradient and Hessian, the plate dimensions, the centroid data, and the mass
// distribution.

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/grader"
)

var calcBatchIDs = []string{
	"calc.deriv.hyperbolic",
	"calc.deriv.parametric_second",
	"calc.deriv.finite_difference",
	"calc.polar.arc_length",
	"calc.vec.torsion",
	"calc.vec.frenet",
	"calc.multi.taylor2",
	"calc.integral.hydrostatic",
	"calc.integral.pappus",
	"calc.integral.moment_inertia",
}

var calcBatchNumericIDs = []string{
	"calc.deriv.parametric_second",
	"calc.deriv.finite_difference",
	"calc.polar.arc_length",
	"calc.vec.torsion",
	"calc.multi.taylor2",
	"calc.integral.hydrostatic",
	"calc.integral.moment_inertia",
}

type calcCorpusConcept struct {
	ID          string `json:"id"`
	GradingType string `json:"grading_type"`
}

func calcShippedTypes(t *testing.T) map[string]string {
	t.Helper()
	paths, err := filepath.Glob("../../../data/concepts/*.json")
	if err != nil {
		t.Fatalf("list corpus: %v", err)
	}
	out := map[string]string{}
	for _, path := range paths {
		if filepath.Base(path) == "enrichment.json" {
			continue
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read corpus %s: %v", path, err)
		}
		var rows []calcCorpusConcept
		if err := json.Unmarshal(raw, &rows); err != nil {
			t.Fatalf("parse corpus %s: %v", path, err)
		}
		for _, r := range rows {
			if strings.HasPrefix(r.ID, "calc.") {
				out[r.ID] = r.GradingType
			}
		}
	}
	return out
}

func calcBatchRegistry(t *testing.T) *generator.Registry {
	t.Helper()
	reg := generator.NewRegistry()
	Register(reg)
	return reg
}

func calcBatchSample(t *testing.T, reg *generator.Registry, id string, i int) generator.Problem {
	t.Helper()
	p, err := reg.GenerateContext(id, generator.GeneratorContext{
		Difficulty: float64(i%20) / 19,
		Seed:       int64(121000 + i),
	})
	if err != nil {
		t.Fatalf("%s: generate: %v", id, err)
	}
	return p
}

func calcFactFloat(t *testing.T, p generator.Problem, key string) float64 {
	t.Helper()
	v, err := strconv.ParseFloat(p.Facts[key], 64)
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not a number", p.Question, key, p.Facts[key])
	}
	return v
}

func calcFactInt(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer", p.Question, key, p.Facts[key])
	}
	return v
}

func calcFactInts(t *testing.T, p generator.Problem, key string) []int {
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

func calcCross(a, b []int) []int {
	return []int{
		a[1]*b[2] - a[2]*b[1],
		a[2]*b[0] - a[0]*b[2],
		a[0]*b[1] - a[1]*b[0],
	}
}

func calcDot(a, b []int) int {
	return a[0]*b[0] + a[1]*b[1] + a[2]*b[2]
}

func calcNorm(a []int) float64 {
	return math.Sqrt(float64(calcDot(a, a)))
}

// ---------------------------------------------------------------------------
// Router contracts
// ---------------------------------------------------------------------------

func TestBatchAnswersGradeWithTheirShippedType(t *testing.T) {
	types := calcShippedTypes(t)
	reg := calcBatchRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	for _, id := range calcBatchIDs {
		typ, ok := types[id]
		if !ok || typ == "" {
			t.Fatalf("%s: no grading_type in the corpus", id)
		}
		switch typ {
		case "numeric", "expression", "symbolic", "tuple":
		default:
			t.Fatalf("%s: unexpected grading_type %q", id, typ)
		}
		for i := 0; i < 120; i++ {
			p := calcBatchSample(t, reg, id, i)
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

// ---------------------------------------------------------------------------
// Independent checks
// ---------------------------------------------------------------------------

func TestBatchHyperbolicMatchesRule(t *testing.T) {
	reg := calcBatchRegistry(t)
	want := map[string]string{
		"sinh":       "cosh(x)",
		"cosh":       "sinh(x)",
		"tanh":       "1/cosh(x)^2",
		"sinh_chain": "3*cosh(3*x)",
		"cosh_chain": "2*x*sinh(x^2)",
		"tanh_chain": "(1/2)/cosh(x/2)^2",
	}
	seen := map[string]bool{}
	for i := 0; i < 300; i++ {
		p := calcBatchSample(t, reg, "calc.deriv.hyperbolic", i)
		rule := p.Facts["rule"]
		seen[rule] = true
		w, ok := want[rule]
		if !ok {
			t.Fatalf("hyperbolic: unknown rule %q", rule)
		}
		if p.Answer != w {
			t.Fatalf("hyperbolic: rule %s answers %q, want %q", rule, p.Answer, w)
		}
	}
	if len(seen) != len(want) {
		t.Fatalf("hyperbolic: expected all %d rules, saw %v", len(want), seen)
	}
}

func TestBatchParametricSecondRecomputes(t *testing.T) {
	reg := calcBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := calcBatchSample(t, reg, "calc.deriv.parametric_second", i)
		xp := calcFactFloat(t, p, "xp")
		yp := calcFactFloat(t, p, "yp")
		xpp := calcFactFloat(t, p, "xpp")
		ypp := calcFactFloat(t, p, "ypp")
		if xp == 0 {
			t.Fatalf("parametric second: x' is zero")
		}
		want := (xp*ypp - yp*xpp) / (xp * xp * xp)
		answer := calcFactFloat(t, p, "answer")
		if math.Abs(answer-want) > 1e-9 {
			t.Fatalf("parametric second: x'=%v y'=%v x''=%v y''=%v answers %v, want %v", xp, yp, xpp, ypp, answer, want)
		}
	}
}

func TestBatchFiniteDifferenceRecomputes(t *testing.T) {
	reg := calcBatchRegistry(t)
	for i := 0; i < 250; i++ {
		p := calcBatchSample(t, reg, "calc.deriv.finite_difference", i)
		x := calcFactFloat(t, p, "x")
		h := calcFactFloat(t, p, "h")
		if h == 0 {
			t.Fatalf("finite difference: zero step")
		}
		f := func(v float64) float64 {
			if p.Facts["function"] == "cube" {
				return v * v * v
			}
			return v * v
		}
		var want float64
		switch p.Facts["method"] {
		case "forward":
			want = (f(x+h) - f(x)) / h
		case "backward":
			want = (f(x) - f(x-h)) / h
		case "central":
			want = (f(x+h) - f(x-h)) / (2 * h)
		default:
			t.Fatalf("finite difference: unknown method %q", p.Facts["method"])
		}
		answer := calcFactFloat(t, p, "answer")
		if math.Abs(answer-want) > 1e-9 {
			t.Fatalf("finite difference: %s for %s at x=%v h=%v answers %v, want %v",
				p.Facts["method"], p.Facts["function"], x, h, answer, want)
		}
	}
}

func TestBatchPolarArcLengthRecomputes(t *testing.T) {
	reg := calcBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := calcBatchSample(t, reg, "calc.polar.arc_length", i)
		a := float64(calcFactInt(t, p, "a"))
		var want float64
		switch p.Facts["kind"] {
		case "sec":
			if p.Facts["b0"] != "pi/4" {
				t.Fatalf("polar arc length: unexpected sec upper limit %q", p.Facts["b0"])
			}
			want = a * math.Tan(math.Pi/4)
		case "cardioid":
			switch p.Facts["b0"] {
			case "pi":
				want = 4
			case "2*pi":
				want = 8
			default:
				t.Fatalf("polar arc length: unexpected cardioid upper limit %q", p.Facts["b0"])
			}
		default:
			t.Fatalf("polar arc length: unknown kind %q", p.Facts["kind"])
		}
		answer := calcFactFloat(t, p, "answer")
		if math.Abs(answer-want) > 1e-9 {
			t.Fatalf("polar arc length: %s a=%v b=%q answers %v, want %v", p.Facts["kind"], a, p.Facts["b0"], answer, want)
		}
	}
}

func TestBatchTorsionRecomputes(t *testing.T) {
	reg := calcBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := calcBatchSample(t, reg, "calc.vec.torsion", i)
		rp := calcFactInts(t, p, "rp")
		rpp := calcFactInts(t, p, "rpp")
		rppp := calcFactInts(t, p, "rppp")
		cross := calcCross(rp, rpp)
		normSq := calcDot(cross, cross)
		if normSq == 0 {
			t.Fatalf("torsion: r' x r'' is zero, curve is not regular")
		}
		want := float64(calcDot(cross, rppp)) / float64(normSq)
		answer := calcFactFloat(t, p, "answer")
		if math.Abs(answer-want) > 1e-9 {
			t.Fatalf("torsion: rp=%v rpp=%v rppp=%v answers %v, want %v", rp, rpp, rppp, answer, want)
		}
	}
}

func TestBatchFrenetRecomputes(t *testing.T) {
	reg := calcBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := calcBatchSample(t, reg, "calc.vec.frenet", i)
		rp := calcFactInts(t, p, "rp")
		rpp := calcFactInts(t, p, "rpp")
		cross := calcCross(rp, rpp)
		norm := calcNorm(cross)
		if norm == 0 {
			t.Fatalf("frenet: r' x r'' is zero, binormal is undefined")
		}
		binormal := []int{int(math.Round(float64(cross[0]) / norm)), int(math.Round(float64(cross[1]) / norm)), int(math.Round(float64(cross[2]) / norm))}
		tangent := rp
		normal := []int{
			binormal[1]*tangent[2] - binormal[2]*tangent[1],
			binormal[2]*tangent[0] - binormal[0]*tangent[2],
			binormal[0]*tangent[1] - binormal[1]*tangent[0],
		}
		want := "(" + strings.Join([]string{
			strconv.Itoa(tangent[0]), strconv.Itoa(tangent[1]), strconv.Itoa(tangent[2]),
			strconv.Itoa(normal[0]), strconv.Itoa(normal[1]), strconv.Itoa(normal[2]),
			strconv.Itoa(binormal[0]), strconv.Itoa(binormal[1]), strconv.Itoa(binormal[2]),
		}, ",") + ")"
		if p.Answer != want {
			t.Fatalf("frenet: rp=%v rpp=%v answers %s, want %s", rp, rpp, p.Answer, want)
		}
	}
}

func TestBatchTaylor2Recomputes(t *testing.T) {
	reg := calcBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := calcBatchSample(t, reg, "calc.multi.taylor2", i)
		f0 := calcFactFloat(t, p, "f0")
		fx := calcFactFloat(t, p, "fx")
		fy := calcFactFloat(t, p, "fy")
		fxx := calcFactFloat(t, p, "fxx")
		fxy := calcFactFloat(t, p, "fxy")
		fyy := calcFactFloat(t, p, "fyy")
		dx := calcFactFloat(t, p, "dx")
		dy := calcFactFloat(t, p, "dy")
		want := f0 + fx*dx + fy*dy + 0.5*(fxx*dx*dx+2*fxy*dx*dy+fyy*dy*dy)
		answer := calcFactFloat(t, p, "answer")
		if math.Abs(answer-want) > 1e-9 {
			t.Fatalf("taylor2: facts %v answers %v, want %v", p.Facts, answer, want)
		}
	}
}

func TestBatchHydrostaticRecomputes(t *testing.T) {
	reg := calcBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := calcBatchSample(t, reg, "calc.integral.hydrostatic", i)
		w := float64(calcFactInt(t, p, "w"))
		a := float64(calcFactInt(t, p, "a"))
		b := float64(calcFactInt(t, p, "b"))
		want := 10 * w * (b*b - a*a) / 2
		answer := calcFactFloat(t, p, "answer")
		if math.Abs(answer-want) > 1e-9 {
			t.Fatalf("hydrostatic: w=%v a=%v b=%v answers %v, want %v", w, a, b, answer, want)
		}
	}
}

func TestBatchPappusRecomputes(t *testing.T) {
	reg := calcBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := calcBatchSample(t, reg, "calc.integral.pappus", i)
		area := calcFactInt(t, p, "area")
		d := calcFactInt(t, p, "d")
		want := fmt.Sprintf("%d*pi", 2*area*d)
		if p.Answer != want {
			t.Fatalf("pappus: area=%d d=%d answers %s, want %s", area, d, p.Answer, want)
		}
	}
}

func TestBatchMomentInertiaRecomputes(t *testing.T) {
	reg := calcBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := calcBatchSample(t, reg, "calc.integral.moment_inertia", i)
		var want float64
		switch p.Facts["kind"] {
		case "rod_center":
			m := float64(calcFactInt(t, p, "m"))
			l := float64(calcFactInt(t, p, "l"))
			want = m * l * l / 12
		case "rod_end":
			m := float64(calcFactInt(t, p, "m"))
			l := float64(calcFactInt(t, p, "l"))
			want = m * l * l / 3
		case "disk":
			m := float64(calcFactInt(t, p, "m"))
			r := float64(calcFactInt(t, p, "r"))
			want = 0.5 * m * r * r
		case "points":
			masses := calcFactInts(t, p, "masses")
			distances := calcFactInts(t, p, "distances")
			if len(masses) != len(distances) {
				t.Fatalf("moment of inertia: %d masses and %d distances", len(masses), len(distances))
			}
			for j := range masses {
				want += float64(masses[j] * distances[j] * distances[j])
			}
		default:
			t.Fatalf("moment of inertia: unknown kind %q", p.Facts["kind"])
		}
		answer := calcFactFloat(t, p, "answer")
		if math.Abs(answer-want) > 1e-9 {
			t.Fatalf("moment of inertia: %s with %v answers %v, want %v", p.Facts["kind"], p.Facts, answer, want)
		}
	}
}

func TestBatchNumericAnswersParseAsDecimals(t *testing.T) {
	reg := calcBatchRegistry(t)
	for _, id := range calcBatchNumericIDs {
		for i := 0; i < 120; i++ {
			p := calcBatchSample(t, reg, id, i)
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
	reg := calcBatchRegistry(t)
	const samples = 150
	for _, id := range calcBatchIDs {
		answers := map[string]bool{}
		questions := map[string]bool{}
		for i := 0; i < samples; i++ {
			p, err := reg.GenerateContext(id, generator.GeneratorContext{Difficulty: 0.12, Seed: int64(122000 + i)})
			if err != nil {
				t.Fatalf("%s: generate: %v", id, err)
			}
			answers[strings.ToLower(strings.TrimSpace(p.Answer))] = true
			questions[p.Question] = true
		}
		if len(answers) < 2 {
			t.Fatalf("%s: only %d distinct answers over %d samples", id, len(answers), samples)
		}
		if len(questions) < 2 {
			t.Fatalf("%s: only %d distinct questions over %d samples", id, len(questions), samples)
		}
	}
}
