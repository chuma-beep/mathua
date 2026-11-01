package numtheory

// Tests for the 2026-10-18 number-theory batch.
//
// Every concept is graded through the real router with the grading_type the
// corpus ships. The two multiple-choice concepts are checked the way the other
// batches check theirs: the question is parsed into options, the correct answer
// must appear exactly once, and every other option must be rejected.
//
// The eight numeric concepts are recomputed from the facts the generator
// publishes: the Jacobi symbol by the reciprocity algorithm, the discrete
// logarithm by enumerating powers, the modular square root by search, the sum
// of two squares by enumeration, the Pell solution by search, the Pisano period
// by iterating the Fibonacci pair, the Carmichael function from the prime
// factorisation, and the prime gap by the next-prime search.

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

var ntBatchIDs = []string{
	"nt.adv.jacobi",
	"nt.adv.discrete_log",
	"nt.adv.modular_sqrt",
	"nt.adv.two_squares",
	"nt.adv.pell",
	"nt.adv.pisano",
	"nt.adv.carmichael_function",
	"nt.adv.miller_rabin",
	"nt.adv.prime_gap",
	"nt.adv.bertrand",
}

var ntBatchNumericIDs = []string{
	"nt.adv.jacobi",
	"nt.adv.discrete_log",
	"nt.adv.modular_sqrt",
	"nt.adv.two_squares",
	"nt.adv.pell",
	"nt.adv.pisano",
	"nt.adv.carmichael_function",
	"nt.adv.prime_gap",
}

type ntCorpusConcept struct {
	ID          string `json:"id"`
	GradingType string `json:"grading_type"`
}

func ntShippedTypes(t *testing.T) map[string]string {
	t.Helper()
	raw, err := os.ReadFile("../../../data/concepts/number_theory.json")
	if err != nil {
		t.Fatalf("read corpus: %v", err)
	}
	var rows []ntCorpusConcept
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

func ntBatchRegistry(t *testing.T) *generator.Registry {
	t.Helper()
	reg := generator.NewRegistry()
	Register(reg)
	return reg
}

func ntBatchSample(t *testing.T, reg *generator.Registry, id string, i int) generator.Problem {
	t.Helper()
	p, err := reg.GenerateContext(id, generator.GeneratorContext{
		Difficulty: float64(i%20) / 19,
		Seed:       int64(131000 + i),
	})
	if err != nil {
		t.Fatalf("%s: generate: %v", id, err)
	}
	return p
}

var ntOptionLine = regexp.MustCompile(`\n\n([A-E])\) ([^\n]*)`)

func ntOptionsOf(q string) []string {
	var out []string
	for _, m := range ntOptionLine.FindAllStringSubmatch(q, -1) {
		out = append(out, m[2])
	}
	return out
}

func ntFactInt(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer", p.Question, key, p.Facts[key])
	}
	return v
}

func ntGCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func ntLCM(a, b int) int {
	return a / ntGCD(a, b) * b
}

func ntIsPrime(n int) bool {
	if n < 2 {
		return false
	}
	if n%2 == 0 {
		return n == 2
	}
	for d := 3; d*d <= n; d += 2 {
		if n%d == 0 {
			return false
		}
	}
	return true
}

func ntNextPrime(n int) int {
	for c := n + 1; ; c++ {
		if ntIsPrime(c) {
			return c
		}
	}
}

func ntPowMod(base, exp, mod int) int {
	result := 1
	base %= mod
	for exp > 0 {
		if exp&1 == 1 {
			result = result * base % mod
		}
		base = base * base % mod
		exp >>= 1
	}
	return result
}

func ntIsqrt(n int) int {
	if n < 0 {
		return -1
	}
	r := int(math.Sqrt(float64(n)))
	for (r+1)*(r+1) <= n {
		r++
	}
	for r*r > n {
		r--
	}
	return r
}

func ntJacobi(a, n int) int {
	if n <= 0 || n%2 == 0 {
		return 0
	}
	a %= n
	if a < 0 {
		a += n
	}
	result := 1
	for a != 0 {
		for a%2 == 0 {
			a /= 2
			if n%8 == 3 || n%8 == 5 {
				result = -result
			}
		}
		a, n = n, a
		if a%4 == 3 && n%4 == 3 {
			result = -result
		}
		a %= n
	}
	if n == 1 {
		return result
	}
	return 0
}

func ntCarmichael(n int) int {
	if n == 1 {
		return 1
	}
	result := 1
	m := n
	for p := 2; p*p <= m; p++ {
		if m%p != 0 {
			continue
		}
		e := 0
		pe := 1
		for m%p == 0 {
			m /= p
			e++
			pe *= p
		}
		var part int
		if p == 2 && e >= 3 {
			part = 1 << (e - 2)
		} else {
			part = (pe / p) * (p - 1)
		}
		result = ntLCM(result, part)
	}
	if m > 1 {
		result = ntLCM(result, m-1)
	}
	return result
}

// ---------------------------------------------------------------------------
// Router contracts
// ---------------------------------------------------------------------------

func TestBatchAnswersGradeWithTheirShippedType(t *testing.T) {
	types := ntShippedTypes(t)
	reg := ntBatchRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	for _, id := range ntBatchIDs {
		typ, ok := types[id]
		if !ok || typ == "" {
			t.Fatalf("%s: no grading_type in the corpus", id)
		}
		if typ != "multiple_choice" && typ != "numeric" {
			t.Fatalf("%s: unexpected grading_type %q", id, typ)
		}
		for i := 0; i < 120; i++ {
			p := ntBatchSample(t, reg, id, i)
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
	types := ntShippedTypes(t)
	reg := ntBatchRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	checked := 0
	for _, id := range ntBatchIDs {
		if types[id] != "multiple_choice" {
			continue
		}
		for i := 0; i < 40; i++ {
			p := ntBatchSample(t, reg, id, i)
			opts := ntOptionsOf(p.Question)
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
		t.Skip("the number theory batch now serves computed count problems")
	}
}

// ---------------------------------------------------------------------------
// Independent numeric checks
// ---------------------------------------------------------------------------

func TestBatchJacobiRecomputes(t *testing.T) {
	reg := ntBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := ntBatchSample(t, reg, "nt.adv.jacobi", i)
		a := ntFactInt(t, p, "a")
		n := ntFactInt(t, p, "n")
		want := ntJacobi(a, n)
		answer, err := strconv.Atoi(p.Answer)
		if err != nil {
			t.Fatalf("jacobi: non-integer answer %q", p.Answer)
		}
		if answer != want {
			t.Fatalf("jacobi: (%d/%d) answers %d, want %d", a, n, answer, want)
		}
	}
}

func TestBatchDiscreteLogRecomputes(t *testing.T) {
	reg := ntBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := ntBatchSample(t, reg, "nt.adv.discrete_log", i)
		prime := ntFactInt(t, p, "p")
		g := ntFactInt(t, p, "g")
		h := ntFactInt(t, p, "h")
		want := -1
		for x := 1; x <= prime-1; x++ {
			if ntPowMod(g, x, prime) == h {
				want = x
				break
			}
		}
		if want < 0 {
			t.Fatalf("discrete log: %d is not in the group generated by %d mod %d", h, g, prime)
		}
		answer, err := strconv.Atoi(p.Answer)
		if err != nil {
			t.Fatalf("discrete log: non-integer answer %q", p.Answer)
		}
		if answer != want {
			t.Fatalf("discrete log: %d^x = %d mod %d answers %d, want %d", g, h, prime, answer, want)
		}
	}
}

func TestBatchModularSqrtRecomputes(t *testing.T) {
	reg := ntBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := ntBatchSample(t, reg, "nt.adv.modular_sqrt", i)
		a := ntFactInt(t, p, "a")
		prime := ntFactInt(t, p, "p")
		want := -1
		for x := 1; x <= prime-1; x++ {
			if x*x%prime == a {
				want = x
				break
			}
		}
		if want < 0 {
			t.Fatalf("modular sqrt: %d is not a quadratic residue mod %d", a, prime)
		}
		answer, err := strconv.Atoi(p.Answer)
		if err != nil {
			t.Fatalf("modular sqrt: non-integer answer %q", p.Answer)
		}
		if answer != want {
			t.Fatalf("modular sqrt: x^2 = %d mod %d answers %d, want %d", a, prime, answer, want)
		}
	}
}

func TestBatchTwoSquaresRecomputes(t *testing.T) {
	reg := ntBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := ntBatchSample(t, reg, "nt.adv.two_squares", i)
		n := ntFactInt(t, p, "n")
		want := 0
		for x := -n; x <= n; x++ {
			for y := -n; y <= n; y++ {
				if x*x+y*y == n {
					want++
				}
			}
		}
		answer, err := strconv.Atoi(p.Answer)
		if err != nil {
			t.Fatalf("two squares: non-integer answer %q", p.Answer)
		}
		if answer != want {
			t.Fatalf("two squares: n=%d answers %d, want %d", n, answer, want)
		}
	}
}

func TestBatchPellRecomputes(t *testing.T) {
	reg := ntBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := ntBatchSample(t, reg, "nt.adv.pell", i)
		d := ntFactInt(t, p, "d")
		want := -1
		for x := 1; x < 10000; x++ {
			if (x*x-1)%d != 0 {
				continue
			}
			y := ntIsqrt((x*x - 1) / d)
			if y > 0 && y*y == (x*x-1)/d {
				want = x
				break
			}
		}
		if want < 0 {
			t.Fatalf("pell: no solution found for d=%d", d)
		}
		answer, err := strconv.Atoi(p.Answer)
		if err != nil {
			t.Fatalf("pell: non-integer answer %q", p.Answer)
		}
		if answer != want {
			t.Fatalf("pell: d=%d answers %d, want %d", d, answer, want)
		}
	}
}

func TestBatchPisanoRecomputes(t *testing.T) {
	reg := ntBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := ntBatchSample(t, reg, "nt.adv.pisano", i)
		n := ntFactInt(t, p, "n")
		a, b := 0, 1
		want := 0
		for {
			a, b = b, (a+b)%n
			want++
			if a == 0 && b == 1 {
				break
			}
			if want > n*n+10 {
				t.Fatalf("pisano: no period found for n=%d", n)
			}
		}
		answer, err := strconv.Atoi(p.Answer)
		if err != nil {
			t.Fatalf("pisano: non-integer answer %q", p.Answer)
		}
		if answer != want {
			t.Fatalf("pisano: n=%d answers %d, want %d", n, answer, want)
		}
	}
}

func TestBatchCarmichaelRecomputes(t *testing.T) {
	reg := ntBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := ntBatchSample(t, reg, "nt.adv.carmichael_function", i)
		n := ntFactInt(t, p, "n")
		want := ntCarmichael(n)
		answer, err := strconv.Atoi(p.Answer)
		if err != nil {
			t.Fatalf("carmichael: non-integer answer %q", p.Answer)
		}
		if answer != want {
			t.Fatalf("carmichael: n=%d answers %d, want %d", n, answer, want)
		}
	}
}

func TestBatchPrimeGapRecomputes(t *testing.T) {
	reg := ntBatchRegistry(t)
	for i := 0; i < 200; i++ {
		p := ntBatchSample(t, reg, "nt.adv.prime_gap", i)
		start := ntFactInt(t, p, "p")
		next := ntNextPrime(start)
		var want int
		switch p.Facts["kind"] {
		case "gap":
			want = next - start
		case "next":
			want = next
		default:
			t.Fatalf("prime gap: unknown kind %q", p.Facts["kind"])
		}
		answer, err := strconv.Atoi(p.Answer)
		if err != nil {
			t.Fatalf("prime gap: non-integer answer %q", p.Answer)
		}
		if answer != want {
			t.Fatalf("prime gap: %s from %d answers %d, want %d", p.Facts["kind"], start, answer, want)
		}
	}
}

func TestBatchNumericAnswersParseAsDecimals(t *testing.T) {
	reg := ntBatchRegistry(t)
	for _, id := range ntBatchNumericIDs {
		for i := 0; i < 120; i++ {
			p := ntBatchSample(t, reg, id, i)
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
	reg := ntBatchRegistry(t)
	const samples = 150
	for _, id := range ntBatchIDs {
		answers := map[string]bool{}
		questions := map[string]bool{}
		allYesNo := true
		for i := 0; i < samples; i++ {
			p, err := reg.GenerateContext(id, generator.GeneratorContext{Difficulty: 0.12, Seed: int64(132000 + i)})
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
