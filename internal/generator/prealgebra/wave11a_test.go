package prealgebra

import (
	"math"
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

type w11aCase struct {
	name string
	gen  generator.Generator
	keys []string
}

func w11aCases() []w11aCase {
	return []w11aCase{
		{"dec.basics.number_line", &w11aNumberLineGen{}, []string{"value", "base", "tenths", "answer"}},
		{"dec.ops.compare_order", &w11aCompareOrderGen{}, []string{"threshold", "count", "answer"}},
		{"dec.convert.repeating", &w11aRepeatingGen{}, []string{"digit", "denominator", "answer"}},
		{"dec.ops.word_multi_step", &w11aWordMultiStepGen{}, []string{"items", "price", "fee", "total", "answer"}},
		{"pct.basics.benchmarks", &w11aBenchmarkGen{}, []string{"percent", "number", "answer"}},
		{"pct.ops.depreciation", &w11aDepreciationGen{}, []string{"value", "rate", "answer"}},
		{"pct.ops.concentration", &w11aConcentrationGen{}, []string{"solute", "solution", "percent", "answer"}},
		{"pct.ops.tip_split", &w11aTipSplitGen{}, []string{"bill", "tip_percent", "people", "answer"}},
		{"ratio.ops.similar_figures", &w11aSimilarFiguresGen{}, []string{"a", "b", "c", "answer"}},
		{"ratio.ops.recipe", &w11aRecipeGen{}, []string{"servings", "other_servings", "cups", "answer"}},
		{"ratio.ops.mixture", &w11aMixtureGen{}, []string{"part_a", "part_b", "total", "answer"}},
		{"ratio.ops.speed", &w11aSpeedGen{}, []string{"rate", "time", "distance", "answer"}},
		{"ratio.ops.density", &w11aDensityGen{}, []string{"density", "volume", "mass", "answer"}},
		{"prealg.expr.substitute_negative", &w11aSubstituteNegativeGen{}, []string{"a", "b", "x", "answer"}},
		{"prealg.expr.factor_common", &w11aFactorCommonGen{}, []string{"a", "b", "gcf", "answer"}},
		{"prealg.var.independent_dependent", &w11aIndependentDependentGen{}, []string{"rate", "start", "input", "answer"}},
	}
}

func w11aTestInt(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	raw, ok := p.Facts[key]
	if !ok {
		t.Fatalf("%s: missing fact %q\nfacts: %v", p.Question, key, p.Facts)
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer\nfacts: %v", p.Question, key, raw, p.Facts)
	}
	return v
}

func w11aTestFloat(t *testing.T, p generator.Problem, key string) float64 {
	t.Helper()
	raw, ok := p.Facts[key]
	if !ok {
		t.Fatalf("%s: missing fact %q\nfacts: %v", p.Question, key, p.Facts)
	}
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not numeric\nfacts: %v", p.Question, key, raw, p.Facts)
	}
	return v
}

func w11aTestWant(t *testing.T, name string, kp int, p generator.Problem) int {
	t.Helper()
	switch name {
	case "dec.basics.number_line":
		value := w11aTestFloat(t, p, "value")
		base := w11aTestInt(t, p, "base")
		tenths := w11aTestInt(t, p, "tenths")
		if got := int(math.Round((value - float64(base)) * 10)); got != tenths {
			t.Fatalf("number line: value %v and base %d give %d tenths, fact says %d\nfacts: %v", value, base, got, tenths, p.Facts)
		}
		if kp == 2 {
			if value >= 0 {
				return int(math.Floor(value + 0.5))
			}
			return int(math.Ceil(value - 0.5))
		}
		return tenths

	case "dec.ops.compare_order":
		threshold := w11aTestFloat(t, p, "threshold")
		raw, ok := p.Facts["values"]
		if !ok || raw == "" {
			t.Fatalf("compare: missing fact values\nfacts: %v", p.Facts)
		}
		count := 0
		for _, part := range strings.Split(raw, ",") {
			v, err := strconv.ParseFloat(strings.TrimSpace(part), 64)
			if err != nil {
				t.Fatalf("compare: value %q is not numeric", part)
			}
			switch kp {
			case 1, 2:
				if v > threshold {
					count++
				}
			case 3:
				if v < threshold {
					count++
				}
			}
		}
		return count

	case "dec.convert.repeating":
		digits, ok := p.Facts["digit"]
		if !ok {
			t.Fatalf("repeating: missing fact digit")
		}
		wantDen := 0
		for range digits {
			wantDen = wantDen*10 + 9
		}
		if got := w11aTestInt(t, p, "denominator"); got != wantDen {
			t.Fatalf("repeating: denominator %d, want %d for %q", got, wantDen, digits)
		}
		num, err := strconv.Atoi(digits)
		if err != nil {
			t.Fatalf("repeating: digit %q is not an integer", digits)
		}
		return num

	case "dec.ops.word_multi_step":
		items := w11aTestInt(t, p, "items")
		price := w11aTestInt(t, p, "price")
		fee := w11aTestInt(t, p, "fee")
		total := w11aTestInt(t, p, "total")
		if kp == 3 {
			return total - items*price - fee
		}
		if total != items*price+fee {
			t.Fatalf("multi-step: total %d is not items*price+fee = %d", total, items*price+fee)
		}
		return total

	case "pct.basics.benchmarks":
		percent := w11aTestInt(t, p, "percent")
		number := w11aTestInt(t, p, "number")
		if number*percent%100 != 0 {
			t.Fatalf("benchmarks: %d%% of %d is not an integer", percent, number)
		}
		return number * percent / 100

	case "pct.ops.depreciation":
		value := w11aTestInt(t, p, "value")
		rate := w11aTestInt(t, p, "rate")
		switch kp {
		case 1:
			return value * (100 - rate) / 100
		case 2:
			return value * rate / 100
		default:
			return value * (100 - rate) * (100 - rate) / 10000
		}

	case "pct.ops.concentration":
		solute := w11aTestInt(t, p, "solute")
		solution := w11aTestInt(t, p, "solution")
		percent := w11aTestInt(t, p, "percent")
		if kp == 3 {
			if solute*100%percent != 0 {
				t.Fatalf("concentration: %d mL at %d%% is not a whole final volume", solute, percent)
			}
			final := solute * 100 / percent
			return final - solution
		}
		if solute != solution*percent/100 {
			t.Fatalf("concentration: solute %d does not match %d mL at %d%%", solute, solution, percent)
		}
		if kp == 1 {
			return percent
		}
		return solute

	case "pct.ops.tip_split":
		bill := w11aTestInt(t, p, "bill")
		percent := w11aTestInt(t, p, "tip_percent")
		people := w11aTestInt(t, p, "people")
		if bill*percent%100 != 0 {
			t.Fatalf("tip: %d%% of %d is not a whole number of dollars", percent, bill)
		}
		tip := bill * percent / 100
		switch kp {
		case 1:
			return tip
		case 2:
			return (bill + tip) / people
		default:
			return tip / people
		}

	case "ratio.ops.similar_figures":
		a := w11aTestInt(t, p, "a")
		b := w11aTestInt(t, p, "b")
		c := w11aTestInt(t, p, "c")
		if b*c%a != 0 {
			t.Fatalf("similar figures: %d x %d is not divisible by %d", b, c, a)
		}
		return b * c / a

	case "ratio.ops.recipe":
		servings := w11aTestInt(t, p, "servings")
		other := w11aTestInt(t, p, "other_servings")
		cups := w11aTestInt(t, p, "cups")
		if kp == 3 {
			given := w11aTestInt(t, p, "given_cups")
			if given*servings%cups != 0 {
				t.Fatalf("recipe: %d cups does not give a whole number of servings", given)
			}
			return given * servings / cups
		}
		if cups*other%servings != 0 {
			t.Fatalf("recipe: %d cups for %d servings is not exact", cups, other)
		}
		return cups * other / servings

	case "ratio.ops.mixture":
		partA := w11aTestInt(t, p, "part_a")
		partB := w11aTestInt(t, p, "part_b")
		total := w11aTestInt(t, p, "total")
		if partA+partB != total {
			t.Fatalf("mixture: %d + %d is not %d", partA, partB, total)
		}
		switch kp {
		case 1:
			return partA
		case 2:
			return partB
		default:
			return total
		}

	case "ratio.ops.speed":
		rate := w11aTestInt(t, p, "rate")
		time := w11aTestInt(t, p, "time")
		distance := w11aTestInt(t, p, "distance")
		if rate*time != distance {
			t.Fatalf("speed: %d x %d is not %d", rate, time, distance)
		}
		switch kp {
		case 1:
			return distance
		case 2:
			return time
		default:
			return rate
		}

	case "ratio.ops.density":
		density := w11aTestInt(t, p, "density")
		volume := w11aTestInt(t, p, "volume")
		mass := w11aTestInt(t, p, "mass")
		if density*volume != mass {
			t.Fatalf("density: %d x %d is not %d", density, volume, mass)
		}
		switch kp {
		case 1:
			return mass
		case 2:
			return volume
		default:
			return density
		}

	case "prealg.expr.substitute_negative":
		a := w11aTestInt(t, p, "a")
		b := w11aTestInt(t, p, "b")
		x := w11aTestInt(t, p, "x")
		if x >= 0 {
			t.Fatalf("substitute: x = %d is not negative", x)
		}
		return a*x + b

	case "prealg.expr.factor_common":
		a := w11aTestInt(t, p, "a")
		b := w11aTestInt(t, p, "b")
		gcf := w11aTestInt(t, p, "gcf")
		if w11aGCD(a, b) != gcf {
			t.Fatalf("factor: gcd(%d, %d) = %d, fact says %d", a, b, w11aGCD(a, b), gcf)
		}
		return gcf

	case "prealg.var.independent_dependent":
		rate := w11aTestInt(t, p, "rate")
		start := w11aTestInt(t, p, "start")
		input := w11aTestInt(t, p, "input")
		output := w11aTestInt(t, p, "output")
		if kp == 2 {
			if (output-start)%rate != 0 {
				t.Fatalf("independent: output %d is not reachable with rate %d", output, rate)
			}
			got := (output - start) / rate
			if rate*got+start != output {
				t.Fatalf("independent: input %d does not produce output %d", got, output)
			}
			return got
		}
		if rate*input+start != output {
			t.Fatalf("independent: %d x %d + %d is not %d", rate, input, start, output)
		}
		return output
	}
	t.Fatalf("unknown kind %q", name)
	return 0
}

func TestWave11AAnswers(t *testing.T) {
	for _, tc := range w11aCases() {
		for kp := 1; kp <= 3; kp++ {
			rand.Seed(int64(20261112 + kp))
			for draw := 0; draw < 60; draw++ {
				p := tc.gen.Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				if p.KPOrdinal != kp {
					t.Fatalf("%s kp%d: KPOrdinal = %d", tc.name, kp, p.KPOrdinal)
				}
				if !strings.HasSuffix(p.Question, "(enter a number)") {
					t.Fatalf("%s kp%d: question does not end with the numeric prompt: %q", tc.name, kp, p.Question)
				}
				if p.Facts["answer"] != p.Answer {
					t.Fatalf("%s kp%d: fact answer = %q, answer = %q", tc.name, kp, p.Facts["answer"], p.Answer)
				}
				for _, key := range tc.keys {
					if _, ok := p.Facts[key]; !ok {
						t.Fatalf("%s kp%d: missing fact %q\nfacts: %v", tc.name, kp, key, p.Facts)
					}
				}
				if w11aPrintsAnswer(p.Question, p.Answer) {
					t.Fatalf("%s kp%d: the stem prints the answer %q: %q", tc.name, kp, p.Answer, p.Question)
				}
				want := w11aTestWant(t, tc.name, kp, p)
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("%s kp%d: answer = %s, want %d\nquestion: %s\nfacts: %v", tc.name, kp, p.Answer, want, p.Question, p.Facts)
				}
			}
		}
	}
}

// TestWave11AOrdinalsDiffer pins the requirement that the three ordinals ask
// different questions at the same random seed.
func TestWave11AOrdinalsDiffer(t *testing.T) {
	for _, tc := range w11aCases() {
		for seed := int64(1); seed <= 8; seed++ {
			questions := [3]string{}
			for kp := 1; kp <= 3; kp++ {
				rand.Seed(seed)
				questions[kp-1] = tc.gen.Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp}).Question
			}
			if questions[0] == questions[1] || questions[0] == questions[2] || questions[1] == questions[2] {
				t.Fatalf("%s seed %d: the three ordinals repeat a question\nkp1: %s\nkp2: %s\nkp3: %s", tc.name, seed, questions[0], questions[1], questions[2])
			}
		}
	}
}

func TestWave11ADiversity(t *testing.T) {
	for idx, tc := range w11aCases() {
		rand.Seed(int64(20261112 + idx))
		overall := map[string]bool{}
		perKP := [4]map[string]bool{{}, {}, {}, {}}
		for i := 0; i < 200; i++ {
			kp := i%3 + 1
			p := tc.gen.Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			overall[p.Answer] = true
			perKP[kp][p.Answer] = true
		}
		if len(overall) < 8 {
			t.Fatalf("%s: only %d distinct answers in 200 draws", tc.name, len(overall))
		}
		for kp := 1; kp <= 3; kp++ {
			if len(perKP[kp]) < 5 {
				t.Fatalf("%s kp%d: only %d distinct answers in 200 draws", tc.name, kp, len(perKP[kp]))
			}
			rand.Seed(int64(20261110 + kp))
			short := map[string]bool{}
			for i := 0; i < 12; i++ {
				p := tc.gen.Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				short[p.Answer] = true
			}
			if len(short) < 5 {
				t.Fatalf("%s kp%d: only %d distinct answers in 12 draws", tc.name, kp, len(short))
			}
		}
	}
}
