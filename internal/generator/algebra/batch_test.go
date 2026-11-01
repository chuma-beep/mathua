package algebra

// Tests for the 2026-10-18 batch. Every concept is graded through the real
// router with the grading_type the corpus ships for it, so a generator and its
// declared type cannot drift apart.
//
// The identity tests are the part that matters for the three rational
// expressions: they evaluate the generated numerator, denominator and answer at
// sample points away from the excluded values, so a cancellation that silently
// changed the value (or dropped a restriction) fails here rather than in front
// of a learner. The factoring and radical tests expand or check the generated
// factors instead of comparing strings.

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/grader"
)

var batchIDs = []string{
	"alg.linear.two_points",
	"alg.linear.intercepts",
	"alg.factor.grouping",
	"alg.factor.cubes",
	"alg.quad.vertex_form",
	"alg.quad.word_problem",
	"alg.poly.end_behavior",
	"alg.rational.simplify",
	"alg.rational.mult_div",
	"alg.rational.add_sub",
	"alg.radical.rational_exponents",
	"alg.radical.rationalize",
	"alg.func.variation",
	"alg.log.change_base",
	"alg.conic.classify",
}

type corpusConcept struct {
	ID          string `json:"id"`
	GradingType string `json:"grading_type"`
}

func shippedTypes(t *testing.T) map[string]string {
	t.Helper()
	raw, err := os.ReadFile("../../../data/concepts/algebra.json")
	if err != nil {
		t.Fatalf("read corpus: %v", err)
	}
	var rows []corpusConcept
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

func batchRegistry(t *testing.T) *generator.Registry {
	t.Helper()
	reg := generator.NewRegistry()
	Register(reg)
	return reg
}

func sample(t *testing.T, reg *generator.Registry, id string, i int) generator.Problem {
	t.Helper()
	p, err := reg.GenerateContext(id, generator.GeneratorContext{
		Difficulty: float64(i%20) / 19,
		Seed:       int64(21000 + i),
	})
	if err != nil {
		t.Fatalf("%s: generate: %v", id, err)
	}
	return p
}

var optionLine = regexp.MustCompile(`\n\n([A-E])\) ([^\n]*)`)

func optionsOf(q string) []string {
	var out []string
	for _, m := range optionLine.FindAllStringSubmatch(q, -1) {
		out = append(out, m[2])
	}
	return out
}

// ---------------------------------------------------------------------------
// Router contracts
// ---------------------------------------------------------------------------

func TestBatchAnswersGradeWithTheirShippedType(t *testing.T) {
	types := shippedTypes(t)
	reg := batchRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	for _, id := range batchIDs {
		typ, ok := types[id]
		if !ok || typ == "" {
			t.Fatalf("%s: no grading_type in the corpus", id)
		}
		for i := 0; i < 120; i++ {
			p := sample(t, reg, id, i)
			if strings.TrimSpace(p.Answer) == "" {
				t.Fatalf("%s: empty answer at sample %d", id, i)
			}
			if p.Facts == nil {
				t.Fatalf("%s: no facts published; its schema would fall back", id)
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

func TestBatchMultipleChoiceRejectsEveryDistractor(t *testing.T) {
	types := shippedTypes(t)
	reg := batchRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	checked := 0
	for _, id := range batchIDs {
		if types[id] != "multiple_choice" {
			continue
		}
		for i := 0; i < 40; i++ {
			p := sample(t, reg, id, i)
			opts := optionsOf(p.Question)
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
		t.Skip("the algebra batch now serves computed numeric answers")
	}
}

func TestBatchExpressionNearMissesAreRejected(t *testing.T) {
	types := shippedTypes(t)
	reg := batchRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	checked := 0
	for _, id := range batchIDs {
		if types[id] != "expression" {
			continue
		}
		for i := 0; i < 30; i++ {
			p := sample(t, reg, id, i)
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

func TestBatchTupleOrderMatters(t *testing.T) {
	types := shippedTypes(t)
	reg := batchRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	checked := 0
	for _, id := range batchIDs {
		if types[id] != "tuple" {
			continue
		}
		for i := 0; i < 60; i++ {
			p := sample(t, reg, id, i)
			parts := strings.Split(p.Answer, ",")
			if len(parts) != 2 {
				t.Fatalf("%s: answer %q is not a pair", id, p.Answer)
			}
			if parts[0] == parts[1] {
				continue // reordering a palindromic pair cannot be told apart
			}
			swapped := parts[1] + "," + parts[0]
			if router.Grade(grader.GradingTuple, p.Answer, swapped).Correct {
				t.Fatalf("%s: swapped pair %q was accepted for %q", id, swapped, p.Answer)
			}
			spaced := "(" + parts[0] + ", " + parts[1] + ")"
			if !router.Grade(grader.GradingTuple, p.Answer, spaced).Correct {
				t.Fatalf("%s: parenthesised pair %q was rejected for %q", id, spaced, p.Answer)
			}
			checked++
		}
	}
	if checked == 0 {
		t.Fatal("no tuple-typed concepts were checked")
	}
}

func TestBatchNumericRejectsNearMisses(t *testing.T) {
	types := shippedTypes(t)
	reg := batchRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	checked := 0
	for _, id := range batchIDs {
		if types[id] != "numeric" {
			continue
		}
		for i := 0; i < 30; i++ {
			p := sample(t, reg, id, i)
			wrong := "1"
			if !strings.Contains(p.Answer, "/") {
				n, err := strconv.Atoi(p.Answer)
				if err != nil {
					t.Fatalf("%s: answer %q is neither an integer nor a fraction", id, p.Answer)
				}
				wrong = strconv.Itoa(n + 1)
			}
			if router.Grade(grader.GradingNumeric, p.Answer, wrong).Correct {
				t.Fatalf("%s: %q was accepted for %q", id, wrong, p.Answer)
			}
			checked++
		}
	}
	if checked == 0 {
		t.Fatal("no numeric-typed concepts were checked")
	}
}

// ---------------------------------------------------------------------------
// Evaluation helpers
// ---------------------------------------------------------------------------

// evalPolyAt evaluates the plain polynomial strings the batch publishes:
// "3x^2 - 4x + 5". It is deliberately a separate implementation from the
// formatter, so a formatting bug is measured by arithmetic rather than by
// comparing the string with itself.
func evalPolyAt(expr string, x float64) float64 {
	s := strings.ReplaceAll(expr, " ", "")
	s = strings.ReplaceAll(s, "-", "+-")
	total := 0.0
	for _, term := range strings.Split(s, "+") {
		if term == "" {
			continue
		}
		if i := strings.Index(term, "x"); i >= 0 {
			coef := 1.0
			switch term[:i] {
			case "", "+":
			case "-":
				coef = -1
			default:
				coef, _ = strconv.ParseFloat(term[:i], 64)
			}
			exp := 1
			if j := strings.Index(term, "^"); j >= 0 {
				exp, _ = strconv.Atoi(term[j+1:])
			}
			total += coef * math.Pow(x, float64(exp))
			continue
		}
		v, _ := strconv.ParseFloat(term, 64)
		total += v
	}
	return total
}

// evalRationalAt evaluates "(numerator)/(denominator)" with plain polynomial
// parts. A string with no top-level slash is evaluated as a polynomial.
func evalRationalAt(expr string, x float64) float64 {
	depth := 0
	for i, r := range expr {
		switch r {
		case '(':
			depth++
		case ')':
			depth--
		case '/':
			if depth == 0 {
				left := strings.Trim(strings.TrimSpace(expr[:i]), "()")
				right := strings.Trim(strings.TrimSpace(expr[i+1:]), "()")
				return evalPolyAt(left, x) / evalPolyAt(right, x)
			}
		}
	}
	return evalPolyAt(expr, x)
}

func intFact(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer", p.Question, key, p.Facts[key])
	}
	return v
}

func fractionFact(t *testing.T, p generator.Problem, key string) (int, int) {
	t.Helper()
	s := p.Facts[key]
	if n, d, ok := strings.Cut(s, "/"); ok {
		num, err := strconv.Atoi(n)
		if err != nil {
			t.Fatalf("%s: fact %q = %q is not a fraction", p.Question, key, s)
		}
		den, err := strconv.Atoi(d)
		if err != nil {
			t.Fatalf("%s: fact %q = %q is not a fraction", p.Question, key, s)
		}
		return num, den
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not numeric", p.Question, key, s)
	}
	return v, 1
}

var identityPoints = []float64{-7, -5, -3, -2, -1, 0, 1, 2, 3, 4, 5, 7}

func closeEnough(a, b float64) bool {
	return math.Abs(a-b) < 1e-6
}

// ---------------------------------------------------------------------------
// Per-concept mathematics
// ---------------------------------------------------------------------------

func TestBatchTwoPointsAreCollinear(t *testing.T) {
	reg := batchRegistry(t)
	for i := 0; i < 200; i++ {
		p := sample(t, reg, "alg.linear.two_points", i)
		m := intFact(t, p, "m")
		b := intFact(t, p, "b")
		x1, y1 := intFact(t, p, "x1"), intFact(t, p, "y1")
		x2, y2 := intFact(t, p, "x2"), intFact(t, p, "y2")
		if dx := x2 - x1; dx == 0 || (y2-y1) != m*dx {
			t.Fatalf("%s: slope %d does not fit the two points (%d,%d) and (%d,%d)", p.Question, m, x1, y1, x2, y2)
		}
		if b != y1-m*x1 {
			t.Fatalf("%s: intercept %d does not fit slope %d and (%d,%d)", p.Question, b, m, x1, y1)
		}
		for _, x := range []float64{-2, 0, 3, 7} {
			if !closeEnough(evalPolyAt(p.Answer, x), float64(m)*x+float64(b)) {
				t.Fatalf("%s: answer %q does not evaluate to %dx+%d", p.Question, p.Answer, m, b)
			}
		}
	}
}

func TestBatchInterceptsSolveTheEquation(t *testing.T) {
	reg := batchRegistry(t)
	for i := 0; i < 200; i++ {
		p := sample(t, reg, "alg.linear.intercepts", i)
		cx, cy, rhs := intFact(t, p, "coef_x"), intFact(t, p, "coef_y"), intFact(t, p, "rhs")
		xi, yi := intFact(t, p, "x_int"), intFact(t, p, "y_int")
		if cx*xi != rhs {
			t.Fatalf("%s: x-intercept %d does not solve %dx = %d", p.Question, xi, cx, rhs)
		}
		if cy*yi != rhs {
			t.Fatalf("%s: y-intercept %d does not solve %dy = %d", p.Question, yi, cy, rhs)
		}
		if p.Answer != strconv.Itoa(xi)+","+strconv.Itoa(yi) {
			t.Fatalf("%s: answer %q does not match the intercepts (%d,%d)", p.Question, p.Answer, xi, yi)
		}
	}
}

func TestBatchGroupingExpandsToThePolynomial(t *testing.T) {
	reg := batchRegistry(t)
	for i := 0; i < 200; i++ {
		p := sample(t, reg, "alg.factor.grouping", i)
		for _, x := range []float64{-3, -1, 0, 2, 5} {
			got := evalPolyAt(p.Facts["binomial"], x) * evalPolyAt(p.Facts["quadratic"], x)
			if !closeEnough(got, evalPolyAt(p.Facts["poly"], x)) {
				t.Fatalf("%s: (%s)(%s) does not expand to %s", p.Question,
					p.Facts["binomial"], p.Facts["quadratic"], p.Facts["poly"])
			}
		}
		if p.Answer != p.Facts["binomial"] {
			t.Fatalf("%s: answer %q is not the shared binomial %q", p.Question, p.Answer, p.Facts["binomial"])
		}
	}
}

func TestBatchCubesCoverBothSigns(t *testing.T) {
	reg := batchRegistry(t)
	seen := map[string]int{}
	for i := 0; i < 200; i++ {
		p := sample(t, reg, "alg.factor.cubes", i)
		kind := p.Facts["kind"]
		seen[kind]++
		if !strings.Contains(p.Question, "quadratic factor") {
			t.Fatalf("%s: the question does not say which factor it wants", p.Question)
		}
		if p.Answer != p.Facts["quadratic"] {
			t.Fatalf("%s: answer %q is not the quadratic factor %q", p.Question, p.Answer, p.Facts["quadratic"])
		}
		if p.Answer == p.Facts["linear"] {
			t.Fatalf("%s: the quadratic factor equals the linear factor", p.Question)
		}
		for _, x := range []float64{-4, -1, 0, 2, 5} {
			got := evalPolyAt(p.Facts["linear"], x) * evalPolyAt(p.Facts["quadratic"], x)
			if !closeEnough(got, evalPolyAt(p.Facts["poly"], x)) {
				t.Fatalf("%s: (%s)(%s) does not expand to %s", p.Question,
					p.Facts["linear"], p.Facts["quadratic"], p.Facts["poly"])
			}
		}
	}
	if seen["sum"] == 0 || seen["difference"] == 0 {
		t.Fatalf("both signs of cubes must be generated, got %v", seen)
	}
}

func TestBatchVertexFormMatchesThePolynomial(t *testing.T) {
	reg := batchRegistry(t)
	for i := 0; i < 200; i++ {
		p := sample(t, reg, "alg.quad.vertex_form", i)
		h, k := intFact(t, p, "h"), intFact(t, p, "k")
		if p.Answer != strconv.Itoa(h)+","+strconv.Itoa(k) {
			t.Fatalf("%s: answer %q is not the vertex (%d,%d)", p.Question, p.Answer, h, k)
		}
		for _, x := range []float64{-4, -1, 0, 2, 6} {
			want := math.Pow(x-float64(h), 2) + float64(k)
			if !closeEnough(evalPolyAt(p.Facts["quad"], x), want) {
				t.Fatalf("%s: %s is not (x-%d)^2+%d", p.Question, p.Facts["quad"], h, k)
			}
		}
	}
}

func TestBatchQuadraticWordProblemIsConsistent(t *testing.T) {
	reg := batchRegistry(t)
	for i := 0; i < 200; i++ {
		p := sample(t, reg, "alg.quad.word_problem", i)
		w := intFact(t, p, "w")
		m, d := intFact(t, p, "m"), intFact(t, p, "d")
		length, area := intFact(t, p, "length"), intFact(t, p, "area")
		if p.Answer != strconv.Itoa(w) {
			t.Fatalf("%s: answer %q is not the width %d", p.Question, p.Answer, w)
		}
		if length != m*w+d || area != w*length {
			t.Fatalf("%s: length %d and area %d do not fit width %d, m %d, d %d",
				p.Question, length, area, w, m, d)
		}
		if !closeEnough(evalPolyAt(p.Facts["quadratic"], float64(w)), 0) {
			t.Fatalf("%s: width %d is not a root of %s", p.Question, w, p.Facts["quadratic"])
		}
	}
}

func TestBatchEndBehaviorFollowsDegreeAndSign(t *testing.T) {
	reg := batchRegistry(t)
	seen := map[string]bool{}
	for i := 0; i < 200; i++ {
		p := sample(t, reg, "alg.poly.end_behavior", i)
		degree, lead := intFact(t, p, "degree"), intFact(t, p, "leading")
		even := degree%2 == 0
		want := 1
		switch {
		case even && lead > 0:
			want = 2
		case even && lead < 0:
			want = 0
		}
		if p.Answer != fmt.Sprintf("%d", want) {
			t.Fatalf("%s: degree %d with leading %d answered %q, want %d", p.Question, degree, lead, p.Answer, want)
		}
		seen[p.Answer] = true
	}
	if len(seen) != 3 {
		t.Fatalf("all three up-end counts should be generated, got %v", seen)
	}
}

func checkIdentity(t *testing.T, p generator.Problem, eval func(x float64) (float64, bool), answer string, label string) {
	t.Helper()
	tested := 0
	for _, x := range identityPoints {
		want, ok := eval(x)
		if !ok {
			continue
		}
		if !closeEnough(evalRationalAt(answer, x), want) {
			t.Fatalf("%s: %s at x=%g is %g, want %g (answer %q)",
				p.Question, label, x, evalRationalAt(answer, x), want, answer)
		}
		tested++
	}
	if tested == 0 {
		t.Fatalf("%s: no sample point avoided the excluded values", p.Question)
	}
}

func TestBatchRationalSimplifyPreservesValue(t *testing.T) {
	reg := batchRegistry(t)
	for i := 0; i < 200; i++ {
		p := sample(t, reg, "alg.rational.simplify", i)
		num, den := p.Facts["numerator"], p.Facts["denominator"]
		checkIdentity(t, p, func(x float64) (float64, bool) {
			d := evalPolyAt(den, x)
			if math.Abs(d) < 1e-9 {
				return 0, false
			}
			return evalPolyAt(num, x) / d, true
		}, p.Answer, "numerator/denominator")
		if p.Facts["valid_for"] == "" {
			t.Fatalf("%s: no domain restriction published", p.Question)
		}
	}
}

func TestBatchRationalMultDivPreservesValue(t *testing.T) {
	reg := batchRegistry(t)
	seen := map[string]bool{}
	for i := 0; i < 200; i++ {
		p := sample(t, reg, "alg.rational.mult_div", i)
		num1, den1 := p.Facts["num1"], p.Facts["den1"]
		num2, den2 := p.Facts["num2"], p.Facts["den2"]
		divide := p.Facts["operation"] == "dividing"
		seen[p.Facts["operation"]] = true
		checkIdentity(t, p, func(x float64) (float64, bool) {
			d1 := evalPolyAt(den1, x)
			d2 := evalPolyAt(den2, x)
			if math.Abs(d1) < 1e-9 || math.Abs(d2) < 1e-9 || math.Abs(evalPolyAt(num2, x)) < 1e-9 {
				return 0, false
			}
			left := evalPolyAt(num1, x) / d1
			if divide {
				return left / (evalPolyAt(num2, x) / d2), true
			}
			return left * (evalPolyAt(num2, x) / d2), true
		}, p.Answer, "the operation")
		if !strings.Contains(p.Facts["valid_for"], "x \\ne") {
			t.Fatalf("%s: the restriction is not a pair of excluded values: %q", p.Question, p.Facts["valid_for"])
		}
	}
	if !seen["multiplying"] || !seen["dividing"] {
		t.Fatalf("both operations must be generated, got %v", seen)
	}
}

func TestBatchRationalAddSubPreservesValue(t *testing.T) {
	reg := batchRegistry(t)
	seen := map[string]bool{}
	for i := 0; i < 200; i++ {
		p := sample(t, reg, "alg.rational.add_sub", i)
		a, m, n := intFact(t, p, "a"), intFact(t, p, "m"), intFact(t, p, "n")
		plus := p.Facts["operation"] == "adding"
		seen[p.Facts["operation"]] = true
		checkIdentity(t, p, func(x float64) (float64, bool) {
			if closeEnough(x, float64(a)) || closeEnough(x, float64(-a)) {
				return 0, false
			}
			left := float64(m) / (x - float64(a))
			right := float64(n) / (x + float64(a))
			if plus {
				return left + right, true
			}
			return left - right, true
		}, p.Answer, "the sum")
	}
	if !seen["adding"] || !seen["subtracting"] {
		t.Fatalf("both operations must be generated, got %v", seen)
	}
}

func TestBatchRationalExponentsAreExact(t *testing.T) {
	reg := batchRegistry(t)
	negatives := 0
	for i := 0; i < 200; i++ {
		p := sample(t, reg, "alg.radical.rational_exponents", i)
		base, index, power := intFact(t, p, "base"), intFact(t, p, "index"), intFact(t, p, "power")
		root, magnitude := intFact(t, p, "root"), intFact(t, p, "magnitude")
		if pow(root, index) != base || pow(root, power) != magnitude {
			t.Fatalf("%s: base %d and magnitude %d do not fit root %d with index %d and power %d",
				p.Question, base, magnitude, root, index, power)
		}
		negative := strings.HasPrefix(p.Facts["exponent"], "-")
		num, den := fractionFact(t, p, "value")
		if negative {
			negatives++
			// value = 1/magnitude, exactly.
			if num*magnitude != den {
				t.Fatalf("%s: value %q is not 1/%d", p.Question, p.Facts["value"], magnitude)
			}
		} else if num != magnitude || den != 1 {
			t.Fatalf("%s: value %q is not the integer %d", p.Question, p.Facts["value"], magnitude)
		}
	}
	if negatives == 0 {
		t.Fatal("the negative-exponent branch never appeared")
	}
}

func TestBatchRationalizeIsSquarefree(t *testing.T) {
	reg := batchRegistry(t)
	for i := 0; i < 200; i++ {
		p := sample(t, reg, "alg.radical.rationalize", i)
		k, a, b := intFact(t, p, "k"), intFact(t, p, "a"), intFact(t, p, "radicand")
		if k != a*b {
			t.Fatalf("%s: k=%d is not a*b=%d*%d", p.Question, k, a, b)
		}
		for _, sq := range []int{4, 9, 25, 49} {
			if b%sq == 0 {
				t.Fatalf("%s: radicand %d is not squarefree (divisible by %d)", p.Question, b, sq)
			}
		}
		// a*sqrt(b) squared must equal k^2/b.
		if a*a*b*b != k*k {
			t.Fatalf("%s: a=%d, b=%d does not square to k=%d", p.Question, a, b, k)
		}
		if p.Answer != strconv.Itoa(a)+","+strconv.Itoa(b) {
			t.Fatalf("%s: answer %q is not (%d,%d)", p.Question, p.Answer, a, b)
		}
	}
}

func TestBatchVariationKeepsTheConstant(t *testing.T) {
	reg := batchRegistry(t)
	seen := map[string]bool{}
	for i := 0; i < 200; i++ {
		p := sample(t, reg, "alg.func.variation", i)
		kind := p.Facts["kind"]
		seen[kind] = true
		k := intFact(t, p, "constant")
		x1, y1 := intFact(t, p, "x1"), intFact(t, p, "y1")
		x2, value := intFact(t, p, "x2"), intFact(t, p, "value")
		if p.Answer != strconv.Itoa(value) {
			t.Fatalf("%s: answer %q is not %d", p.Question, p.Answer, value)
		}
		if kind == "directly" {
			if k*x1 != y1 || k*x2 != value {
				t.Fatalf("%s: direct variation with k=%d does not fit (%d,%d) and x=%d -> %d",
					p.Question, k, x1, y1, x2, value)
			}
		} else {
			if x1*y1 != k || x2*value != k {
				t.Fatalf("%s: inverse variation with k=%d does not fit (%d,%d) and x=%d -> %d",
					p.Question, k, x1, y1, x2, value)
			}
		}
	}
	if !seen["directly"] || !seen["inversely"] {
		t.Fatalf("both kinds of variation must be generated, got %v", seen)
	}
}

func TestBatchChangeBaseIsExact(t *testing.T) {
	reg := batchRegistry(t)
	for i := 0; i < 200; i++ {
		p := sample(t, reg, "alg.log.change_base", i)
		base, arg := intFact(t, p, "base"), intFact(t, p, "argument")
		a, m, n := intFact(t, p, "a"), intFact(t, p, "m"), intFact(t, p, "n")
		if pow(a, m) != base || pow(a, n) != arg {
			t.Fatalf("%s: base %d and argument %d are not powers %d^%d and %d^%d",
				p.Question, base, arg, a, m, a, n)
		}
		if n%m == 0 {
			t.Fatalf("%s: n=%d is a multiple of m=%d, so the answer is a whole log", p.Question, n, m)
		}
		num, den := fractionFact(t, p, "value")
		if num*m != n*den {
			t.Fatalf("%s: value %q is not %d/%d", p.Question, p.Facts["value"], n, m)
		}
	}
}

func TestBatchConicClassesMatchTheEquations(t *testing.T) {
	reg := batchRegistry(t)
	seen := map[string]bool{}
	for i := 0; i < 200; i++ {
		p := sample(t, reg, "alg.conic.classify", i)
		class := ""
		switch {
		case strings.Contains(p.Question, "represent a circle"):
			class = "circle"
		case strings.Contains(p.Question, "represent an ellipse"):
			class = "ellipse"
		case strings.Contains(p.Question, "represent a parabola"):
			class = "parabola"
		case strings.Contains(p.Question, "represent a hyperbola"):
			class = "hyperbola"
		}
		if class == "" {
			t.Fatalf("%s: no conic class in the question", p.Question)
		}
		seen[class] = true
		want := 0
		switch class {
		case "circle":
			want = strings.Count(p.Question, "x^2 + y^2 =")
		case "ellipse":
			want = strings.Count(p.Question, " + y^2/")
		case "parabola":
			want = strings.Count(p.Question, "y = ")
		case "hyperbola":
			want = strings.Count(p.Question, " - y^2/")
		}
		if p.Answer != fmt.Sprintf("%d", want) {
			t.Fatalf("%s: answer %q does not match %d %s equations", p.Question, p.Answer, want, class)
		}
	}
	if len(seen) != 4 {
		t.Fatalf("all four conic classes should be generated, got %v", seen)
	}
}

// TestBatchQuestionsFitTheTimeBudget checks the surface a threshold measures:
// short questions and generation far below a learner's first second. The
// thresholds themselves are set per type in the corpus, at the times the survey
// recorded.
func TestBatchQuestionsFitTheTimeBudget(t *testing.T) {
	reg := batchRegistry(t)
	start := time.Now()
	longest := 0
	for _, id := range batchIDs {
		for i := 0; i < 100; i++ {
			p := sample(t, reg, id, i)
			if n := len([]rune(p.Question)); n > longest {
				longest = n
			}
		}
	}
	if longest > 420 {
		t.Fatalf("a question is %d runes long; that no longer fits a 15-45s item", longest)
	}
	if elapsed := time.Since(start); elapsed > 2*time.Second {
		t.Fatalf("generating 1500 problems took %s", elapsed)
	}
	t.Logf("1500 generated items, longest %d runes, %s total", longest, time.Since(start).Round(time.Millisecond))
}

func pow(a, b int) int {
	v := 1
	for i := 0; i < b; i++ {
		v *= a
	}
	return v
}
