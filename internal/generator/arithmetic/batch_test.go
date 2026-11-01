package arithmetic

// Tests for the 2026-10-18 arithmetic batch. Every concept is graded through
// the real router with the grading_type the corpus ships, so a generator and its
// declared type cannot drift apart.
//
// The per-concept tests recompute the answer from the published facts rather
// than comparing a string with itself. The branch tests matter most: several
// generators pick between two shapes (multiply/divide, add/subtract,
// full/needed, parenthesised/bare, positive/negative cube), and a fact that
// exists on only one shape would drop the authored schema at runtime.

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

var arithBatchIDs = []string{
	"arith.compare.whole",
	"arith.place.expanded",
	"arith.place.powers_ten",
	"arith.estimate.sum",
	"arith.estimate.product",
	"arith.word.multi_step",
	"arith.div.interpret_remainder",
	"arith.div.zero",
	"arith.neg.word",
	"arith.exp.negative_base",
	"arith.exp.power_product",
	"arith.exp.power_quotient",
	"arith.sqrt.cube",
	"arith.sqrt.estimate",
}

type arithCorpusConcept struct {
	ID          string `json:"id"`
	GradingType string `json:"grading_type"`
}

func arithShippedTypes(t *testing.T) map[string]string {
	t.Helper()
	raw, err := os.ReadFile("../../../data/concepts/arithmetic.json")
	if err != nil {
		t.Fatalf("read corpus: %v", err)
	}
	var rows []arithCorpusConcept
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

func arithRegistry(t *testing.T) *generator.Registry {
	t.Helper()
	reg := generator.NewRegistry()
	Register(reg)
	return reg
}

func arithSample(t *testing.T, reg *generator.Registry, id string, i int) generator.Problem {
	t.Helper()
	p, err := reg.GenerateContext(id, generator.GeneratorContext{
		Difficulty: float64(i%20) / 19,
		Seed:       int64(31000 + i),
	})
	if err != nil {
		t.Fatalf("%s: generate: %v", id, err)
	}
	return p
}

func arithIntFact(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer", p.Question, key, p.Facts[key])
	}
	return v
}

func arithOptionsOf(q string) []string {
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

func TestArithBatchAnswersGradeWithTheirShippedType(t *testing.T) {
	types := arithShippedTypes(t)
	reg := arithRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	for _, id := range arithBatchIDs {
		typ, ok := types[id]
		if !ok || typ == "" {
			t.Fatalf("%s: no grading_type in the corpus", id)
		}
		for i := 0; i < 120; i++ {
			p := arithSample(t, reg, id, i)
			if strings.TrimSpace(p.Answer) == "" {
				t.Fatalf("%s: empty answer at sample %d", id, i)
			}
			if p.Facts == nil {
				t.Fatalf("%s: no facts published; its schema would fall back", id)
			}
			if len(p.Question) > 320 || len(p.Answer) > 80 {
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

func TestArithBatchMultipleChoiceRejectsEveryDistractor(t *testing.T) {
	types := arithShippedTypes(t)
	reg := arithRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	checked := 0
	for _, id := range arithBatchIDs {
		if types[id] != "multiple_choice" {
			continue
		}
		for i := 0; i < 60; i++ {
			p := arithSample(t, reg, id, i)
			opts := arithOptionsOf(p.Question)
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
			// The correct sentence deliberately has no standalone single letter,
			// so a bare option letter cannot hit the grader's letter branch.
			if router.Grade(grader.GradingMultipleChoice, p.Answer, "A").Correct {
				t.Fatalf("%s: bare letter A was accepted for %q", id, p.Answer)
			}
			checked++
		}
	}
	if checked == 0 {
		t.Skip("the arithmetic batch now serves computed answers")
	}
}

func TestArithBatchExpressionNearMissesAreRejected(t *testing.T) {
	types := arithShippedTypes(t)
	reg := arithRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	checked := 0
	for _, id := range arithBatchIDs {
		if types[id] != "expression" {
			continue
		}
		for i := 0; i < 30; i++ {
			p := arithSample(t, reg, id, i)
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

func TestArithBatchTupleOrderMatters(t *testing.T) {
	types := arithShippedTypes(t)
	reg := arithRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	checked := 0
	for _, id := range arithBatchIDs {
		if types[id] != "tuple" {
			continue
		}
		for i := 0; i < 80; i++ {
			p := arithSample(t, reg, id, i)
			parts := strings.Split(p.Answer, ",")
			if len(parts) < 2 {
				t.Fatalf("%s: answer %q is not a tuple", id, p.Answer)
			}
			if parts[0] == parts[1] {
				continue // reordering equal first parts cannot be told apart
			}
			swapped := append([]string(nil), parts...)
			swapped[0], swapped[1] = swapped[1], swapped[0]
			if router.Grade(grader.GradingTuple, p.Answer, strings.Join(swapped, ",")).Correct {
				t.Fatalf("%s: swapped tuple %q was accepted for %q", id, strings.Join(swapped, ","), p.Answer)
			}
			spaced := "(" + strings.Join(parts, ", ") + ")"
			if !router.Grade(grader.GradingTuple, p.Answer, spaced).Correct {
				t.Fatalf("%s: parenthesised tuple %q was rejected for %q", id, spaced, p.Answer)
			}
			checked++
		}
	}
	if checked == 0 {
		t.Fatal("no tuple-typed concepts were checked")
	}
}

func TestArithBatchNumericRejectsNearMisses(t *testing.T) {
	types := arithShippedTypes(t)
	reg := arithRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	checked := 0
	for _, id := range arithBatchIDs {
		if types[id] != "numeric" {
			continue
		}
		for i := 0; i < 40; i++ {
			p := arithSample(t, reg, id, i)
			n, err := strconv.Atoi(p.Answer)
			if err != nil {
				t.Fatalf("%s: answer %q is not an integer", id, p.Answer)
			}
			if router.Grade(grader.GradingNumeric, p.Answer, strconv.Itoa(n+1)).Correct {
				t.Fatalf("%s: %d was accepted for %q", id, n+1, p.Answer)
			}
			checked++
		}
	}
	if checked == 0 {
		t.Fatal("no numeric-typed concepts were checked")
	}
}

func TestArithBatchTuplePartsAreExact(t *testing.T) {
	router := grader.NewRouter()
	defer router.Close()

	cases := []struct {
		answer string
		accept string
		reject string
	}{
		{"4000,300,5", "(4000, 300, 5)", "4000.0,300,5"},
		{"5099,7486,9961", "(5099, 7486, 9961)", "7486,5099,9961"},
		{"7,8", "7,8", "7.0,8"},
	}
	for _, c := range cases {
		if !router.Grade(grader.GradingTuple, c.answer, c.accept).Correct {
			t.Errorf("tuple %q did not accept %q", c.answer, c.accept)
		}
		if router.Grade(grader.GradingTuple, c.answer, c.reject).Correct {
			t.Errorf("tuple %q accepted %q", c.answer, c.reject)
		}
	}
}

func TestArithBatchExpressionEquivalents(t *testing.T) {
	router := grader.NewRouter()
	defer router.Close()

	accept := [][2]string{
		{"9x^2", "3^2x^2"},
		{"8x^3", "2^3x^3"},
		{"x^2/9", "(1/9)x^2"},
		{"x^3/8", "(1/8)x^3"},
	}
	for _, c := range accept {
		res := router.Grade(grader.GradingExpression, c[0], c[1])
		if res.Unavailable {
			t.Skip("SymPy runtime missing; the expression cases cannot run here")
		}
		if !res.Correct {
			t.Errorf("expression %q did not accept the equivalent %q", c[0], c[1])
		}
	}
	reject := [][2]string{
		{"9x^2", "6x^2"},
		{"x^2/9", "x^2/4"},
		{"x^3/8", "x^3/4"},
	}
	for _, c := range reject {
		if router.Grade(grader.GradingExpression, c[0], c[1]).Correct {
			t.Errorf("expression %q accepted the wrong answer %q", c[0], c[1])
		}
	}
}

// ---------------------------------------------------------------------------
// Per-concept mathematics
// ---------------------------------------------------------------------------

func TestArithCompareWholeOrdersTheSameNumbers(t *testing.T) {
	reg := arithRegistry(t)
	for i := 0; i < 200; i++ {
		p := arithSample(t, reg, "arith.compare.whole", i)
		orderedText := strings.Split(p.Answer, ",")
		shownText := strings.Split(p.Facts["numbers"], ", ")
		if len(orderedText) != 3 || len(shownText) != 3 {
			t.Fatalf("%s: malformed tuple %q or numbers %q", p.Question, p.Answer, p.Facts["numbers"])
		}
		ordered := make([]int, 3)
		shown := make([]int, 3)
		for j := range orderedText {
			var err error
			ordered[j], err = strconv.Atoi(orderedText[j])
			if err != nil {
				t.Fatalf("%s: answer part %q is not an integer", p.Question, orderedText[j])
			}
			shown[j], err = strconv.Atoi(shownText[j])
			if err != nil {
				t.Fatalf("%s: presented number %q is not an integer", p.Question, shownText[j])
			}
		}
		if !(ordered[0] < ordered[1] && ordered[1] < ordered[2]) {
			t.Fatalf("%s: answer %q is not ascending", p.Question, p.Answer)
		}
		want := append([]int(nil), shown...)
		for a := 0; a < len(want); a++ {
			for b := a + 1; b < len(want); b++ {
				if want[b] < want[a] {
					want[a], want[b] = want[b], want[a]
				}
			}
		}
		for j := range want {
			if ordered[j] != want[j] {
				t.Fatalf("%s: answer %v is not the sorted copy of %v", p.Question, ordered, shown)
			}
		}
		if shown[0] < shown[1] && shown[1] < shown[2] {
			t.Fatalf("%s: the numbers are already presented in ascending order", p.Question)
		}
	}
}

func TestArithPlaceExpandedPartsSumToTheNumber(t *testing.T) {
	reg := arithRegistry(t)
	partCounts := map[int]int{}
	for i := 0; i < 300; i++ {
		p := arithSample(t, reg, "arith.place.expanded", i)
		n := arithIntFact(t, p, "n")
		parts := strings.Split(p.Answer, ",")
		sum := 0
		prev := 100000
		for _, s := range parts {
			v, err := strconv.Atoi(s)
			if err != nil {
				t.Fatalf("%s: part %q is not an integer", p.Question, s)
			}
			if v <= 0 || v >= prev {
				t.Fatalf("%s: parts %q are not descending nonzero place values", p.Question, p.Answer)
			}
			prev = v
			sum += v
		}
		if sum != n {
			t.Fatalf("%s: parts %q sum to %d, want %d", p.Question, p.Answer, sum, n)
		}
		if len(parts) < 2 || len(parts) > 4 {
			t.Fatalf("%s: got %d parts, want 2..4", p.Question, len(parts))
		}
		partCounts[len(parts)]++
	}
	for _, n := range []int{2, 3, 4} {
		if partCounts[n] == 0 {
			t.Fatalf("expanded form never produced %d parts", n)
		}
	}
}

func TestArithPlacePowersTenShiftByThePower(t *testing.T) {
	reg := arithRegistry(t)
	ops := map[string]int{}
	for i := 0; i < 300; i++ {
		p := arithSample(t, reg, "arith.place.powers_ten", i)
		n := arithIntFact(t, p, "n")
		power := arithIntFact(t, p, "power")
		value := arithIntFact(t, p, "value")
		scale := int(math.Pow10(power))
		switch p.Facts["operation"] {
		case "multiply":
			if n*scale != value {
				t.Fatalf("%s: %d * 10^%d = %d, not %d", p.Question, n, power, n*scale, value)
			}
		case "divide":
			if value*scale != n {
				t.Fatalf("%s: %d / 10^%d = %d, not %d", p.Question, n, power, n/scale, value)
			}
		default:
			t.Fatalf("%s: unknown operation %q", p.Question, p.Facts["operation"])
		}
		ops[p.Facts["operation"]]++
	}
	if ops["multiply"] == 0 || ops["divide"] == 0 {
		t.Fatalf("powers of ten did not exercise both operations: %v", ops)
	}
}

func TestArithEstimateSumRoundsEachNumber(t *testing.T) {
	reg := arithRegistry(t)
	ops := map[string]int{}
	for i := 0; i < 300; i++ {
		p := arithSample(t, reg, "arith.estimate.sum", i)
		a := arithIntFact(t, p, "a")
		b := arithIntFact(t, p, "b")
		ra := arithIntFact(t, p, "rounded_a")
		rb := arithIntFact(t, p, "rounded_b")
		wantA := int(math.Round(float64(a)/100)) * 100
		wantB := int(math.Round(float64(b)/100)) * 100
		if ra != wantA || rb != wantB {
			t.Fatalf("%s: rounded %d,%d to %d,%d, want %d,%d", p.Question, a, b, ra, rb, wantA, wantB)
		}
		var want int
		switch p.Facts["operation"] {
		case "add":
			want = ra + rb
		case "subtract":
			want = ra - rb
			if want <= 0 {
				t.Fatalf("%s: subtraction estimate is not positive: %d", p.Question, want)
			}
		default:
			t.Fatalf("%s: unknown operation %q", p.Question, p.Facts["operation"])
		}
		if got := arithIntFact(t, p, "estimate"); got != want {
			t.Fatalf("%s: estimate = %d, want %d", p.Question, got, want)
		}
		ops[p.Facts["operation"]]++
	}
	if ops["add"] == 0 || ops["subtract"] == 0 {
		t.Fatalf("estimate sums did not exercise both operations: %v", ops)
	}
}

func TestArithEstimateProductRoundsEachFactor(t *testing.T) {
	reg := arithRegistry(t)
	for i := 0; i < 300; i++ {
		p := arithSample(t, reg, "arith.estimate.product", i)
		a := arithIntFact(t, p, "a")
		b := arithIntFact(t, p, "b")
		ra := arithIntFact(t, p, "rounded_a")
		rb := arithIntFact(t, p, "rounded_b")
		if ra != int(math.Round(float64(a)/10))*10 || rb != int(math.Round(float64(b)/10))*10 {
			t.Fatalf("%s: rounded %d,%d to %d,%d", p.Question, a, b, ra, rb)
		}
		if got, want := arithIntFact(t, p, "estimate"), ra*rb; got != want {
			t.Fatalf("%s: estimate = %d, want %d", p.Question, got, want)
		}
	}
}

func TestArithMultiStepSubtractsAfterMultiplying(t *testing.T) {
	reg := arithRegistry(t)
	for i := 0; i < 300; i++ {
		p := arithSample(t, reg, "arith.word.multi_step", i)
		a := arithIntFact(t, p, "a")
		b := arithIntFact(t, p, "b")
		c := arithIntFact(t, p, "c")
		want := a*b - c
		if want <= 0 {
			t.Fatalf("%s: non-positive result %d", p.Question, want)
		}
		if got := arithIntFact(t, p, "result"); got != want {
			t.Fatalf("%s: result = %d, want %d", p.Question, got, want)
		}
		for _, bad := range []string{" 1 pencils", " 1 muffins", " 1 books"} {
			if strings.Contains(p.Question, bad) {
				t.Fatalf("%s: singular count uses a plural noun (%q)", p.Question, bad)
			}
		}
		if c == 1 && !strings.Contains(p.Question, " 1 pencil") &&
			!strings.Contains(p.Question, " 1 muffin") && !strings.Contains(p.Question, " 1 book") {
			t.Fatalf("%s: a removal of one did not use a singular noun", p.Question)
		}
	}
}

func TestArithInterpretRemainderKeepsTheRemainder(t *testing.T) {
	reg := arithRegistry(t)
	modes := map[string]int{}
	for i := 0; i < 300; i++ {
		p := arithSample(t, reg, "arith.div.interpret_remainder", i)
		total := arithIntFact(t, p, "total")
		per := arithIntFact(t, p, "per")
		q := arithIntFact(t, p, "quotient")
		r := arithIntFact(t, p, "remainder")
		if q*per+r != total {
			t.Fatalf("%s: %d*%d+%d != %d", p.Question, q, per, r, total)
		}
		if r <= 0 || r >= per {
			t.Fatalf("%s: remainder %d is not in 1..%d", p.Question, r, per-1)
		}
		want := q
		if p.Facts["mode"] == "needed" {
			want = q + 1
		}
		if got := arithIntFact(t, p, "answer"); got != want {
			t.Fatalf("%s: answer = %d, want %d for mode %q", p.Question, got, want, p.Facts["mode"])
		}
		modes[p.Facts["mode"]]++
	}
	if modes["full"] == 0 || modes["needed"] == 0 {
		t.Fatalf("remainder interpretation did not exercise both modes: %v", modes)
	}
}

func TestArithNegativeBaseRespectsParentheses(t *testing.T) {
	reg := arithRegistry(t)
	shapes := map[string]int{}
	for i := 0; i < 400; i++ {
		p := arithSample(t, reg, "arith.exp.negative_base", i)
		base := arithIntFact(t, p, "base")
		power := arithIntFact(t, p, "power")
		value := arithIntFact(t, p, "value")
		paren := p.Facts["parenthesized"]
		want := -int(math.Pow(float64(base), float64(power)))
		shape := "bare"
		if paren == "true" {
			want = int(math.Pow(float64(-base), float64(power)))
			shape = "parenthesized"
		}
		if value != want {
			t.Fatalf("%s: value = %d, want %d (parenthesized=%s)", p.Question, value, want, paren)
		}
		shapes[shape]++
		shapes[strconv.Itoa(power)]++
	}
	if shapes["bare"] == 0 || shapes["parenthesized"] == 0 {
		t.Fatalf("negative-base powers did not exercise both shapes: %v", shapes)
	}
	if shapes["2"] == 0 || shapes["3"] == 0 {
		t.Fatalf("negative-base powers did not exercise both parities: %v", shapes)
	}
}

func TestArithPowerProductDistributesOverTheCoefficient(t *testing.T) {
	reg := arithRegistry(t)
	for i := 0; i < 300; i++ {
		p := arithSample(t, reg, "arith.exp.power_product", i)
		coef := arithIntFact(t, p, "coefficient")
		power := arithIntFact(t, p, "power")
		want := strconv.Itoa(int(math.Pow(float64(coef), float64(power)))) + "x^" + strconv.Itoa(power)
		if p.Answer != want {
			t.Fatalf("%s: answer %q, want %q", p.Question, p.Answer, want)
		}
	}
}

func TestArithPowerQuotientDistributesOverTheDenominator(t *testing.T) {
	reg := arithRegistry(t)
	for i := 0; i < 300; i++ {
		p := arithSample(t, reg, "arith.exp.power_quotient", i)
		den := arithIntFact(t, p, "denominator")
		power := arithIntFact(t, p, "power")
		want := "x^" + strconv.Itoa(power) + "/" + strconv.Itoa(int(math.Pow(float64(den), float64(power))))
		if p.Answer != want {
			t.Fatalf("%s: answer %q, want %q", p.Question, p.Answer, want)
		}
	}
}

func TestArithCubeRootsAreExact(t *testing.T) {
	reg := arithRegistry(t)
	signs := map[string]int{}
	for i := 0; i < 300; i++ {
		p := arithSample(t, reg, "arith.sqrt.cube", i)
		radicand := arithIntFact(t, p, "radicand")
		root := arithIntFact(t, p, "root")
		if root*root*root != radicand {
			t.Fatalf("%s: %d^3 = %d, not radicand %d", p.Question, root, root*root*root, radicand)
		}
		if radicand < 0 {
			signs["negative"]++
		} else {
			signs["positive"]++
		}
	}
	if signs["negative"] == 0 || signs["positive"] == 0 {
		t.Fatalf("cube roots did not exercise both signs: %v", signs)
	}
}

func TestArithSquareRootBoundsBracketTheRadicand(t *testing.T) {
	reg := arithRegistry(t)
	for i := 0; i < 300; i++ {
		p := arithSample(t, reg, "arith.sqrt.estimate", i)
		n := arithIntFact(t, p, "radicand")
		lower := arithIntFact(t, p, "lower")
		upper := arithIntFact(t, p, "upper")
		if upper != lower+1 {
			t.Fatalf("%s: bounds %d,%d are not consecutive", p.Question, lower, upper)
		}
		if !(lower*lower < n && n < upper*upper) {
			t.Fatalf("%s: %d is not between %d^2 and %d^2", p.Question, n, lower, upper)
		}
		if isPerfectSquare(n) {
			t.Fatalf("%s: radicand %d is a perfect square", p.Question, n)
		}
	}
}

func TestArithIntegerWordProblemsAddTheSignedChange(t *testing.T) {
	reg := arithRegistry(t)
	units := map[string]int{}
	for i := 0; i < 300; i++ {
		p := arithSample(t, reg, "arith.neg.word", i)
		start := arithIntFact(t, p, "start")
		change := arithIntFact(t, p, "change")
		result := arithIntFact(t, p, "result")
		if start+change != result {
			t.Fatalf("%s: %d + (%d) = %d, not %d", p.Question, start, change, start+change, result)
		}
		if change == 0 {
			t.Fatalf("%s: zero change", p.Question)
		}
		if p.Facts["unit"] == "m" && result > 0 {
			t.Fatalf("%s: the diver rose above the surface to %d m", p.Question, result)
		}
		units[p.Facts["unit"]]++
	}
	if units["°C"] == 0 || units["m"] == 0 {
		t.Fatalf("integer word problems did not exercise both contexts: %v", units)
	}
}

func TestArithDivisionByZeroCountsTrueStatements(t *testing.T) {
	reg := arithRegistry(t)
	seen := map[string]bool{}
	for i := 0; i < 120; i++ {
		p := arithSample(t, reg, "arith.div.zero", i)
		for _, letter := range "abcdefgh" {
			if !strings.Contains(p.Question, "("+string(letter)+") ") {
				t.Fatalf("%s: missing item %c", p.Question, letter)
			}
		}
		n, err := strconv.Atoi(p.Answer)
		if err != nil || n < 0 || n > 8 {
			t.Fatalf("%s: answer %q is not a count from 0 to 8", p.Question, p.Answer)
		}
		if p.Facts["count"] != p.Answer {
			t.Fatalf("%s: fact count %q does not match answer %q", p.Question, p.Facts["count"], p.Answer)
		}
		seen[p.Question] = true
	}
	if len(seen) < 20 {
		t.Fatalf("division by zero produced only %d distinct prompts in 120 samples", len(seen))
	}
}
