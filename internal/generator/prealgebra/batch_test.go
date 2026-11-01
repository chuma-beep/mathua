package prealgebra

// Tests for the 2026-10-18 prealgebra batch. Every concept is graded through
// the real router with the grading_type the corpus ships, so a generator and
// its declared type cannot drift apart.
//
// The per-concept tests recompute the answer from the published facts rather
// than comparing a string with itself. The branch tests matter most: percent
// change and error have two directions, the formula generator has three
// formulas, the inequality graph has four symbols, and temperature converts in
// both directions. A fact that exists on only one branch would drop the
// authored schema at runtime.

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

var prealgBatchIDs = []string{
	"pct.ops.change",
	"pct.ops.error",
	"pct.ops.commission",
	"pct.ops.markup",
	"pct.ops.simple_interest",
	"ratio.ops.percent_proportion",
	"ratio.ops.better_buy",
	"prealg.expr.translate",
	"prealg.expr.formula",
	"prealg.ineq.word",
	"prealg.ineq.graph",
	"prealg.measure.convert",
	"prealg.measure.metric",
	"prealg.measure.temperature",
}

type prealgCorpusConcept struct {
	ID          string `json:"id"`
	GradingType string `json:"grading_type"`
}

func prealgShippedTypes(t *testing.T) map[string]string {
	t.Helper()
	raw, err := os.ReadFile("../../../data/concepts/prealgebra.json")
	if err != nil {
		t.Fatalf("read corpus: %v", err)
	}
	var rows []prealgCorpusConcept
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

func prealgRegistry(t *testing.T) *generator.Registry {
	t.Helper()
	reg := generator.NewRegistry()
	Register(reg)
	return reg
}

func prealgSample(t *testing.T, reg *generator.Registry, id string, i int) generator.Problem {
	t.Helper()
	p, err := reg.GenerateContext(id, generator.GeneratorContext{
		Difficulty: float64(i%20) / 19,
		Seed:       int64(41000 + i),
	})
	if err != nil {
		t.Fatalf("%s: generate: %v", id, err)
	}
	return p
}

func prealgIntFact(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer", p.Question, key, p.Facts[key])
	}
	return v
}

func prealgFloatFact(t *testing.T, p generator.Problem, key string) float64 {
	t.Helper()
	v, err := strconv.ParseFloat(p.Facts[key], 64)
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not a number", p.Question, key, p.Facts[key])
	}
	return v
}

func prealgNear(a, b float64) bool {
	return math.Abs(a-b) < 0.005
}

func prealgOptionsOf(q string) []string {
	var out []string
	re := regexp.MustCompile(`\n\n([A-E])\) ([^\n]*)`)
	for _, m := range re.FindAllStringSubmatch(q, -1) {
		out = append(out, m[2])
	}
	return out
}

// ---------------------------------------------------------------------------
// Router contracts
// ---------------------------------------------------------------------------

func TestPrealgBatchAnswersGradeWithTheirShippedType(t *testing.T) {
	types := prealgShippedTypes(t)
	reg := prealgRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	for _, id := range prealgBatchIDs {
		typ, ok := types[id]
		if !ok || typ == "" {
			t.Fatalf("%s: no grading_type in the corpus", id)
		}
		for i := 0; i < 120; i++ {
			p := prealgSample(t, reg, id, i)
			if strings.TrimSpace(p.Answer) == "" {
				t.Fatalf("%s: empty answer at sample %d", id, i)
			}
			if p.Facts == nil {
				t.Fatalf("%s: no facts published; its schema would fall back", id)
			}
			if len(p.Question) > 520 || len(p.Answer) > 140 {
				t.Fatalf("%s: question/answer too long (%d/%d): %q / %q",
					id, len(p.Question), len(p.Answer), p.Question, p.Answer)
			}
			res := router.Grade(grader.GradingType(typ), p.Answer, p.Answer)
			if res.Unavailable {
				t.Skipf("%s: grader unavailable (SymPy runtime missing)", id)
			}
			if !res.Correct {
				t.Fatalf("%s: the generator's own answer %q does not grade as %s", id, p.Answer, typ)
			}
		}
	}
}

func TestPrealgBatchMultipleChoiceRejectsEveryDistractor(t *testing.T) {
	types := prealgShippedTypes(t)
	reg := prealgRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	checked := 0
	for _, id := range prealgBatchIDs {
		if types[id] != "multiple_choice" {
			continue
		}
		for i := 0; i < 80; i++ {
			p := prealgSample(t, reg, id, i)
			opts := prealgOptionsOf(p.Question)
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
			if router.Grade(grader.GradingMultipleChoice, p.Answer, "A").Correct {
				t.Fatalf("%s: bare letter A was accepted for %q", id, p.Answer)
			}
			checked++
		}
	}
	if checked == 0 {
		t.Skip("the prealgebra batch now serves computed answers")
	}
}

func TestPrealgBatchExpressionNearMissesAreRejected(t *testing.T) {
	types := prealgShippedTypes(t)
	reg := prealgRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	checked := 0
	for _, id := range prealgBatchIDs {
		if types[id] != "expression" {
			continue
		}
		for i := 0; i < 40; i++ {
			p := prealgSample(t, reg, id, i)
			res := router.Grade(grader.GradingExpression, p.Answer, "("+p.Answer+")+1")
			if res.Unavailable {
				t.Skip("SymPy runtime missing; the expression cases cannot run here")
			}
			if res.Correct {
				t.Fatalf("%s: adding 1 to %q was accepted", id, p.Answer)
			}
			checked++
		}
	}
	if checked == 0 {
		t.Fatal("no expression-typed concepts were checked")
	}
}

func TestPrealgBatchNumericRejectsNearMisses(t *testing.T) {
	types := prealgShippedTypes(t)
	reg := prealgRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	checked := 0
	for _, id := range prealgBatchIDs {
		if types[id] != "numeric" {
			continue
		}
		for i := 0; i < 40; i++ {
			p := prealgSample(t, reg, id, i)
			n, err := strconv.ParseFloat(p.Answer, 64)
			if err != nil {
				t.Fatalf("%s: answer %q is not numeric", id, p.Answer)
			}
			wrong := strconv.FormatFloat(n+1, 'f', -1, 64)
			if router.Grade(grader.GradingNumeric, p.Answer, wrong).Correct {
				t.Fatalf("%s: %s was accepted for %q", id, wrong, p.Answer)
			}
			checked++
		}
	}
	if checked == 0 {
		t.Fatal("no numeric-typed concepts were checked")
	}
}

// ---------------------------------------------------------------------------
// Per-concept mathematics
// ---------------------------------------------------------------------------

func TestPrealgPercentChangeUsesTheOriginalValue(t *testing.T) {
	reg := prealgRegistry(t)
	directions := map[string]int{}
	for i := 0; i < 300; i++ {
		p := prealgSample(t, reg, "pct.ops.change", i)
		old := prealgIntFact(t, p, "old")
		newValue := prealgIntFact(t, p, "new")
		change := prealgIntFact(t, p, "change")
		percent := prealgIntFact(t, p, "percent")
		if change != absInt(newValue-old) {
			t.Fatalf("%s: change %d != |%d-%d|", p.Question, change, newValue, old)
		}
		if change*100 != percent*old {
			t.Fatalf("%s: %d%% of %d is not %d", p.Question, percent, old, change)
		}
		if change == 0 || percent <= 0 || percent >= 100 {
			t.Fatalf("%s: degenerate percent change %d%% / %d", p.Question, percent, change)
		}
		want := "increase"
		if newValue < old {
			want = "decrease"
		}
		if p.Facts["direction"] != want {
			t.Fatalf("%s: direction %q, want %q", p.Question, p.Facts["direction"], want)
		}
		directions[want]++
	}
	if directions["increase"] == 0 || directions["decrease"] == 0 {
		t.Fatalf("percent change did not exercise both directions: %v", directions)
	}
}

func TestPrealgPercentErrorUsesTheAcceptedValue(t *testing.T) {
	reg := prealgRegistry(t)
	signs := map[string]int{}
	for i := 0; i < 300; i++ {
		p := prealgSample(t, reg, "pct.ops.error", i)
		measured := prealgIntFact(t, p, "measured")
		actual := prealgIntFact(t, p, "actual")
		difference := prealgIntFact(t, p, "difference")
		percent := prealgIntFact(t, p, "percent")
		if difference != absInt(measured-actual) {
			t.Fatalf("%s: difference %d != |%d-%d|", p.Question, difference, measured, actual)
		}
		if difference*100 != percent*actual {
			t.Fatalf("%s: %d%% of %d is not %d", p.Question, percent, actual, difference)
		}
		if measured > actual {
			signs["over"]++
		} else if measured < actual {
			signs["under"]++
		} else {
			t.Fatalf("%s: measured equals actual", p.Question)
		}
	}
	if signs["over"] == 0 || signs["under"] == 0 {
		t.Fatalf("percent error did not exercise both directions: %v", signs)
	}
}

func TestPrealgCommissionIsRateTimesSales(t *testing.T) {
	reg := prealgRegistry(t)
	for i := 0; i < 300; i++ {
		p := prealgSample(t, reg, "pct.ops.commission", i)
		sales := prealgFloatFact(t, p, "sales")
		rate := prealgFloatFact(t, p, "rate")
		commission := prealgFloatFact(t, p, "commission")
		if !prealgNear(commission, sales*rate/100) {
			t.Fatalf("%s: commission %.2f != %.2f%% of %.2f", p.Question, commission, rate, sales)
		}
		if commission <= 0 || commission >= sales {
			t.Fatalf("%s: implausible commission %.2f on sales %.2f", p.Question, commission, sales)
		}
	}
}

func TestPrealgMarkupAddsTheMarkup(t *testing.T) {
	reg := prealgRegistry(t)
	for i := 0; i < 300; i++ {
		p := prealgSample(t, reg, "pct.ops.markup", i)
		cost := prealgFloatFact(t, p, "cost")
		rate := prealgFloatFact(t, p, "rate")
		markup := prealgFloatFact(t, p, "markup")
		price := prealgFloatFact(t, p, "price")
		if !prealgNear(markup, cost*rate/100) {
			t.Fatalf("%s: markup %.2f != %.2f%% of %.2f", p.Question, markup, rate, cost)
		}
		if !prealgNear(price, cost+markup) {
			t.Fatalf("%s: price %.2f != cost %.2f + markup %.2f", p.Question, price, cost, markup)
		}
		if price <= cost {
			t.Fatalf("%s: markup did not raise the price", p.Question)
		}
	}
}

func TestPrealgSimpleInterestUsesPrincipalRateTime(t *testing.T) {
	reg := prealgRegistry(t)
	units := map[string]int{}
	for i := 0; i < 300; i++ {
		p := prealgSample(t, reg, "pct.ops.simple_interest", i)
		principal := prealgFloatFact(t, p, "principal")
		rate := prealgFloatFact(t, p, "rate")
		time := prealgFloatFact(t, p, "time")
		interest := prealgFloatFact(t, p, "interest")
		if !prealgNear(interest, principal*rate*time/100) {
			t.Fatalf("%s: interest %.2f != %.2f%% of %.2f for %.0f years", p.Question, interest, rate, principal, time)
		}
		wantUnit := "years"
		if time == 1 {
			wantUnit = "year"
		}
		if p.Facts["time_unit"] != wantUnit {
			t.Fatalf("%s: time unit %q, want %q", p.Question, p.Facts["time_unit"], wantUnit)
		}
		units[wantUnit]++
	}
	if units["year"] == 0 || units["years"] == 0 {
		t.Fatalf("simple interest did not exercise both time units: %v", units)
	}
}

func TestPrealgPercentProportionSolvesTheWhole(t *testing.T) {
	reg := prealgRegistry(t)
	for i := 0; i < 300; i++ {
		p := prealgSample(t, reg, "ratio.ops.percent_proportion", i)
		part := prealgIntFact(t, p, "part")
		percent := prealgIntFact(t, p, "percent")
		whole := prealgIntFact(t, p, "whole")
		if part*100 != percent*whole {
			t.Fatalf("%s: %d/100 != %d/%d", p.Question, part, percent, whole)
		}
		if part >= whole || percent >= 100 {
			t.Fatalf("%s: the unknown whole is not larger than the part", p.Question)
		}
	}
}

func TestPrealgBetterBuyChoosesTheLowerUnitPrice(t *testing.T) {
	reg := prealgRegistry(t)
	for i := 0; i < 300; i++ {
		p := prealgSample(t, reg, "ratio.ops.better_buy", i)
		unit1 := prealgFloatFact(t, p, "unit_1")
		unit2 := prealgFloatFact(t, p, "unit_2")
		if prealgNear(unit1, unit2) {
			t.Fatalf("%s: the two unit prices are equal", p.Question)
		}
		want := (unit1 - unit2) * 100
		if want < 0 {
			want = -want
		}
		got, err := strconv.Atoi(p.Answer)
		if err != nil {
			t.Fatalf("%s: answer %q is not an integer", p.Question, p.Answer)
		}
		if math.Abs(float64(got)-want) > 0.001 {
			t.Fatalf("%s: answer %d, want the difference %v", p.Question, got, want)
		}
	}
}

func TestPrealgTranslateExpressionGradesWithTheRouter(t *testing.T) {
	reg := prealgRegistry(t)
	router := grader.NewRouter()
	defer router.Close()
	answers := map[string]bool{}
	for i := 0; i < 200; i++ {
		p := prealgSample(t, reg, "prealg.expr.translate", i)
		if p.Facts["phrase"] == "" || p.Facts["reason"] == "" {
			t.Fatalf("%s: missing phrase or reason", p.Question)
		}
		res := router.Grade(grader.GradingExpression, p.Answer, p.Answer)
		if res.Unavailable {
			t.Skip("SymPy runtime missing; the expression cases cannot run here")
		}
		if !res.Correct {
			t.Fatalf("%s: answer %q does not grade as an expression", p.Question, p.Answer)
		}
		answers[p.Answer] = true
	}
	if len(answers) < 4 {
		t.Fatalf("translation pool is too narrow: %v", answers)
	}
}

func TestPrealgFormulaEvaluatesEveryShape(t *testing.T) {
	reg := prealgRegistry(t)
	shapes := map[string]int{}
	for i := 0; i < 300; i++ {
		p := prealgSample(t, reg, "prealg.expr.formula", i)
		v1 := prealgIntFact(t, p, "v1")
		v2 := prealgIntFact(t, p, "v2")
		result := prealgIntFact(t, p, "result")
		var want int
		switch p.Facts["formula"] {
		case "P = 2l + 2w":
			want = 2*v1 + 2*v2
		case "A = bh/2":
			want = v1 * v2 / 2
		case "d = rt":
			want = v1 * v2
		default:
			t.Fatalf("%s: unknown formula %q", p.Question, p.Facts["formula"])
		}
		if result != want {
			t.Fatalf("%s: result %d, want %d", p.Question, result, want)
		}
		shapes[p.Facts["formula"]]++
	}
	for _, formula := range []string{"P = 2l + 2w", "A = bh/2", "d = rt"} {
		if shapes[formula] == 0 {
			t.Fatalf("formula pool never produced %q: %v", formula, shapes)
		}
	}
}

func TestPrealgInequalityWordFactsAgree(t *testing.T) {
	reg := prealgRegistry(t)
	contexts := map[string]int{}
	for i := 0; i < 300; i++ {
		p := prealgSample(t, reg, "prealg.ineq.word", i)
		fixed := prealgIntFact(t, p, "fixed")
		rate := prealgIntFact(t, p, "rate")
		limit := prealgIntFact(t, p, "limit")
		bound := prealgIntFact(t, p, "bound")
		answer := prealgIntFact(t, p, "answer")
		if limit != rate*bound+fixed {
			t.Fatalf("%s: limit %d != %d*%d+%d", p.Question, limit, rate, bound, fixed)
		}
		if answer != bound {
			t.Fatalf("%s: answer %d != bound %d", p.Question, answer, bound)
		}
		if p.Facts["context"] == "elevator capacity" {
			if !strings.Contains(p.Question, "each person weighs") || strings.Contains(p.Question, "kg per person") {
				t.Fatalf("%s: elevator wording hides who the per-person weight belongs to", p.Question)
			}
		}
		if p.Facts["relation"] == "at least" {
			if !strings.Contains(p.Facts["inequality"], "\\ge") {
				t.Fatalf("%s: at-least inequality lacks \\ge: %q", p.Question, p.Facts["inequality"])
			}
			if p.Facts["ask"] != "least" {
				t.Fatalf("%s: at-least question asks for %q", p.Question, p.Facts["ask"])
			}
		} else {
			if !strings.Contains(p.Facts["inequality"], "\\le") {
				t.Fatalf("%s: at-most inequality lacks \\le: %q", p.Question, p.Facts["inequality"])
			}
			if p.Facts["ask"] != "greatest" {
				t.Fatalf("%s: at-most question asks for %q", p.Question, p.Facts["ask"])
			}
		}
		contexts[p.Facts["context"]]++
	}
	if len(contexts) < 4 {
		t.Fatalf("inequality word problems did not exercise every context: %v", contexts)
	}
}

func TestPrealgInequalityGraphMatchesTheSymbol(t *testing.T) {
	reg := prealgRegistry(t)
	symbols := map[string]int{}
	for i := 0; i < 400; i++ {
		p := prealgSample(t, reg, "prealg.ineq.graph", i)
		bound := prealgIntFact(t, p, "bound")
		low := prealgIntFact(t, p, "low")
		high := prealgIntFact(t, p, "high")
		symbol := p.Facts["symbol"]
		want := 0
		for x := low; x <= high; x++ {
			ok := false
			switch symbol {
			case ">":
				ok = x > bound
			case ">=":
				ok = x >= bound
			case "<":
				ok = x < bound
			case "<=":
				ok = x <= bound
			}
			if ok {
				want++
			}
		}
		got := prealgIntFact(t, p, "answer")
		if got != want {
			t.Fatalf("%s: answer %d, want %d", p.Question, got, want)
		}
		symbols[symbol]++
	}
	for _, symbol := range []string{">", ">=", "<", "<="} {
		if symbols[symbol] == 0 {
			t.Fatalf("inequality graph never produced %q: %v", symbol, symbols)
		}
	}
}

func TestPrealgCustomaryConversionMultiplies(t *testing.T) {
	reg := prealgRegistry(t)
	conversions := map[string]int{}
	for i := 0; i < 300; i++ {
		p := prealgSample(t, reg, "prealg.measure.convert", i)
		value := prealgFloatFact(t, p, "value")
		factor := prealgFloatFact(t, p, "factor")
		answer := prealgFloatFact(t, p, "answer")
		if !prealgNear(answer, value*factor) {
			t.Fatalf("%s: %s x %s = %s, not %s", p.Question, p.Facts["value"], p.Facts["factor"], p.Facts["value"], p.Facts["answer"])
		}
		if p.Facts["from_unit"] == "" || p.Facts["to_unit"] == "" {
			t.Fatalf("%s: missing unit names", p.Question)
		}
		conversions[p.Facts["from_unit"]+"->"+p.Facts["to_unit"]]++
	}
	if len(conversions) < 4 {
		t.Fatalf("customary conversion pool is too narrow: %v", conversions)
	}
}

func TestPrealgMetricConversionUsesPowersOfTen(t *testing.T) {
	reg := prealgRegistry(t)
	conversions := map[string]int{}
	for i := 0; i < 300; i++ {
		p := prealgSample(t, reg, "prealg.measure.metric", i)
		value := prealgFloatFact(t, p, "value")
		factor := prealgFloatFact(t, p, "factor")
		power := prealgIntFact(t, p, "power")
		answer := prealgFloatFact(t, p, "answer")
		if !prealgNear(answer, value*factor) {
			t.Fatalf("%s: %s x %s = %s, not %s", p.Question, p.Facts["value"], p.Facts["factor"], p.Facts["value"], p.Facts["answer"])
		}
		if factor != math.Pow10(power) {
			t.Fatalf("%s: factor %v != 10^%d", p.Question, factor, power)
		}
		conversions[p.Facts["from_unit"]+"->"+p.Facts["to_unit"]]++
	}
	if len(conversions) < 4 {
		t.Fatalf("metric conversion pool is too narrow: %v", conversions)
	}
}

func TestPrealgTemperatureConversionIsExact(t *testing.T) {
	reg := prealgRegistry(t)
	directions := map[string]int{}
	for i := 0; i < 300; i++ {
		p := prealgSample(t, reg, "prealg.measure.temperature", i)
		value := prealgIntFact(t, p, "value")
		answer := prealgIntFact(t, p, "answer")
		switch p.Facts["direction"] {
		case "CtoF":
			if want := 9*value/5 + 32; answer != want {
				t.Fatalf("%s: %d C = %d F, want %d", p.Question, value, answer, want)
			}
		case "FtoC":
			if want := 5 * (value - 32) / 9; answer != want {
				t.Fatalf("%s: %d F = %d C, want %d", p.Question, value, answer, want)
			}
		default:
			t.Fatalf("%s: unknown direction %q", p.Question, p.Facts["direction"])
		}
		directions[p.Facts["direction"]]++
	}
	if directions["CtoF"] == 0 || directions["FtoC"] == 0 {
		t.Fatalf("temperature conversion did not exercise both directions: %v", directions)
	}
}
