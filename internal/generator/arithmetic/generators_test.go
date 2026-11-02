package arithmetic

import (
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func fuzzGen(t *testing.T, gen generator.Generator) {
	t.Helper()
	for i := 0; i < 100; i++ {
		d := rand.Float64()
		p := gen.Generate(generator.GeneratorContext{Difficulty: d})
		if p.Question == "" || p.Answer == "" || p.Explanation == "" {
			t.Errorf("empty field at difficulty=%.2f: q=%q a=%q e=%q", d, p.Question, p.Answer, p.Explanation)
		}
	}
}

func TestAddGen(t *testing.T)        { fuzzGen(t, &addGen{}) }
func TestAddWordGen(t *testing.T)    { fuzzGen(t, &addWordGen{}) }
func TestSubGen(t *testing.T)        { fuzzGen(t, &subGen{}) }
func TestSubBorrowGen(t *testing.T)  { fuzzGen(t, &subBorrowGen{}) }
func TestSubWordGen(t *testing.T)    { fuzzGen(t, &subWordGen{}) }
func TestPlaceValueGen(t *testing.T) { fuzzGen(t, &placeValueGen{}) }

// The stem used to ask for "the value of the digit in the tens place" of 87,
// which reads as 8 and expected 80. It now names the digit and asks for its
// value, so the question has one reading. It also never names a digit that the
// number does not have, which the old draw could do: "the digit 0 is in the
// tens place of 7".
func TestPlaceValueNamesTheDigitItAsksAbout(t *testing.T) {
	for _, label := range []string{"tens", "hundreds", "thousands"} {
		gen := &placeValueGen{max: 9999, label: label}
		for i := 0; i < 200; i++ {
			p := gen.Generate(generator.GeneratorContext{Difficulty: 0.5})
			digit := p.Facts["digit"]
			if digit == "0" {
				t.Fatalf("%s: named digit 0 in %q", label, p.Question)
			}
			want := "the digit " + digit + " is in the " + label + " place"
			if !strings.Contains(p.Question, want) {
				t.Fatalf("%s: question %q does not contain %q", label, p.Question, want)
			}
			if p.Answer != p.Facts["value"] {
				t.Fatalf("%s: answer %q is not the value %q", label, p.Answer, p.Facts["value"])
			}
			if !strings.HasSuffix(p.Question, "?") {
				t.Fatalf("%s: question %q lost its question mark", label, p.Question)
			}
		}
	}
}
func TestRoundGen(t *testing.T)           { fuzzGen(t, &roundGen{}) }
func TestMultConceptGen(t *testing.T)     { fuzzGen(t, &multConceptGen{}) }
func TestMultByGen(t *testing.T)          { fuzzGen(t, &multByGen{}) }
func TestMultTablesGen(t *testing.T)      { fuzzGen(t, &multTablesGen{}) }
func TestMultGen(t *testing.T)            { fuzzGen(t, &multGen{}) }
func TestMultWordGen(t *testing.T)        { fuzzGen(t, &multWordGen{}) }
func TestDivConceptGen(t *testing.T)      { fuzzGen(t, &divConceptGen{}) }
func TestDivBasicGen(t *testing.T)        { fuzzGen(t, &divBasicGen{}) }
func TestDivRemainderGen(t *testing.T)    { fuzzGen(t, &divRemainderGen{}) }
func TestDivLongGen(t *testing.T)         { fuzzGen(t, &divLongGen{}) }
func TestDivWordGen(t *testing.T)         { fuzzGen(t, &divWordGen{}) }
func TestFactorFindGen(t *testing.T)      { fuzzGen(t, &factorFindGen{}) }
func TestPrimeGen(t *testing.T)           { fuzzGen(t, &primeGen{}) }
func TestPrimeFactGen(t *testing.T)       { fuzzGen(t, &primeFactGen{}) }
func TestGcfGen(t *testing.T)             { fuzzGen(t, &gcfGen{}) }
func TestLcmGen(t *testing.T)             { fuzzGen(t, &lcmGen{}) }
func TestCompositeGen(t *testing.T)       { fuzzGen(t, &compositeGen{}) }
func TestExpConceptGen(t *testing.T)      { fuzzGen(t, &expConceptGen{}) }
func TestExpEvalGen(t *testing.T)         { fuzzGen(t, &expEvalGen{}) }
func TestExpProductRuleGen(t *testing.T)  { fuzzGen(t, &expProductRuleGen{}) }
func TestExpQuotientRuleGen(t *testing.T) { fuzzGen(t, &expQuotientRuleGen{}) }
func TestExpPowerRuleGen(t *testing.T)    { fuzzGen(t, &expPowerRuleGen{}) }
func TestSqrtPerfectGen(t *testing.T)     { fuzzGen(t, &sqrtPerfectGen{}) }
func TestSqrtSimplifyGen(t *testing.T)    { fuzzGen(t, &sqrtSimplifyGen{}) }
func TestNegNumberLineGen(t *testing.T)   { fuzzGen(t, &negNumberLineGen{}) }
func TestNegAddSubGen(t *testing.T)       { fuzzGen(t, &negAddSubGen{}) }
func TestNegMultDivGen(t *testing.T)      { fuzzGen(t, &negMultDivGen{}) }
func TestOrderOpsGen(t *testing.T)        { fuzzGen(t, &orderOpsGen{}) }
func TestDecIntroGen(t *testing.T)        { fuzzGen(t, &decIntroGen{}) }

var decZeroTenthsRe = regexp.MustCompile(`Write \d+\.0 as a mixed number\.`)
var decZeroNumerRe = regexp.MustCompile(`\\frac\{0\}\{`)

// A mixed number needs a real fractional part: tenths == 0 produced
// "Write 2.0 as a mixed number" with answer "2 0/10" (degenerate).
func TestDecIntroGen_NonzeroTenths(t *testing.T) {
	g := &decIntroGen{}
	for i := 0; i < 500; i++ {
		d := rand.Float64()
		p := g.Generate(generator.GeneratorContext{Difficulty: d})
		if decZeroTenthsRe.MatchString(p.Question) {
			t.Fatalf("degenerate zero-tenths question at difficulty=%.2f: %q", d, p.Question)
		}
		if decZeroNumerRe.MatchString(p.Explanation) {
			t.Fatalf("zero-numerator fraction in explanation at difficulty=%.2f: %q", d, p.Explanation)
		}
	}
}

func TestFuzz(t *testing.T) {
	reg := generator.NewRegistry()
	Register(reg)
	for _, id := range reg.Concepts() {
		gen, _ := reg.Get(id)
		t.Run(id, func(t *testing.T) {
			for i := 0; i < 100; i++ {
				d := rand.Float64()
				p := gen.Generate(generator.GeneratorContext{Difficulty: d})
				if p.Question == "" || p.Answer == "" || p.Explanation == "" {
					t.Fatalf("empty field at difficulty=%.2f for %s: q=%q a=%q", d, id, p.Question, p.Answer)
				}
			}
		})
	}
}

// The question used to ask for "Q R" and grade it with a regex, so a learner on
// a phone keypad could not type the letter and the report arrived with an empty
// answer. The answer is now the remainder alone, which has to stay a single
// integer below the divisor.
func TestRemainderAnswerIsTheRemainder(t *testing.T) {
	g := &divRemainderGen{}
	for i := 0; i < 200; i++ {
		p := g.Generate(generator.GeneratorContext{Difficulty: rand.Float64(), Seed: rand.Int63()})
		if !strings.Contains(p.Question, "remainder") {
			t.Fatalf("question does not ask for a remainder: %q", p.Question)
		}
		if !regexp.MustCompile(`^\d+$`).MatchString(p.Answer) {
			t.Fatalf("answer %q is not a bare integer", p.Answer)
		}
		divisor, _ := strconv.Atoi(p.Facts["divisor"])
		dividend, _ := strconv.Atoi(p.Facts["dividend"])
		quotient, _ := strconv.Atoi(p.Facts["quotient"])
		remainder, _ := strconv.Atoi(p.Answer)
		if remainder >= divisor || remainder < 1 {
			t.Fatalf("remainder %d out of range for divisor %d", remainder, divisor)
		}
		if divisor*quotient+remainder != dividend {
			t.Fatalf("%d x %d + %d != %d", divisor, quotient, remainder, dividend)
		}
	}
}

// The answer must be fully simplified and must round-trip through the grader.
// sqrt(100) used to be served as "5 sqrt(4)", and 10 was rejected.
func TestSqrtSimplifyAnswerIsFullySimplified(t *testing.T) {
	gen := &sqrtSimplifyGen{}
	nRe := regexp.MustCompile(`\\sqrt\{(\d+)\}`)
	for i := 0; i < 200; i++ {
		p := gen.Generate(generator.GeneratorContext{Difficulty: rand.Float64()})
		m := nRe.FindStringSubmatch(p.Question)
		if m == nil {
			t.Fatalf("question prints no radicand: %q", p.Question)
		}
		n, _ := strconv.Atoi(m[1])
		coef, rad, ok := sqrtParts(p.Answer)
		if !ok {
			t.Fatalf("answer %q is not a coefficient and a radical", p.Answer)
		}
		if coef*coef*rad != n {
			t.Fatalf("answer %q does not equal sqrt(%d)", p.Answer, n)
		}
		for f := 2; f*f <= rad; f++ {
			if rad%(f*f) == 0 {
				t.Fatalf("answer %q leaves a square factor %d inside the radical", p.Answer, f*f)
			}
		}
		if got := gen.Grade(p.Answer, p.Answer); !got.Correct {
			t.Fatalf("answer %q does not grade as itself", p.Answer)
		}
	}
}

func TestSqrtSimplifyAcceptsTheFormsALearnerTypes(t *testing.T) {
	gen := &sqrtSimplifyGen{}
	for _, answer := range []string{"5 sqrt(2)", "5sqrt(2)", "5\u221a2", `5\sqrt{2}`, " sqrt(2) "} {
		expected := "5 sqrt(2)"
		if strings.Contains(answer, "sqrt") && !strings.Contains(answer, "5") {
			expected = "sqrt(2)"
		}
		if got := gen.Grade(expected, answer); !got.Correct {
			t.Errorf("grade(%q, %q) = incorrect", expected, answer)
		}
	}
	if got := gen.Grade("5 sqrt(2)", "5 sqrt(3)"); got.Correct {
		t.Error("a different radical graded correct")
	}
	if got := gen.Grade("5 sqrt(2)", "10"); got.Correct {
		t.Error("10 graded correct for 5 sqrt(2)")
	}
}

// TestSqrtSimplifyReadsTheRadicalMathliveWrites pins the transport shape. The
// math field writes an empty radical as sqrt(), so "3sqrt()2" reached the
// grader and was rejected as malformed (report 233).
func TestSqrtSimplifyReadsTheRadicalMathliveWrites(t *testing.T) {
	cases := map[string]bool{
		"3sqrt()2":   true,
		"3 sqrt(2)":  true,
		"3\\sqrt{2}": true,
		"3\u221a2":   true,
		"seventeen":  false,
	}
	for answer, want := range cases {
		coef, rad, ok := sqrtParts(answer)
		if ok != want {
			t.Fatalf("sqrtParts(%q) ok = %v, want %v", answer, ok, want)
		}
		if want && (coef != 3 || rad != 2) {
			t.Fatalf("sqrtParts(%q) = %d sqrt(%d), want 3 sqrt(2)", answer, coef, rad)
		}
	}
}

func TestAddDoubleSumsAreCorrect(t *testing.T) {
	for _, ensureCarry := range []bool{false, true} {
		gen := &addGen{minA: 10, maxA: 99, minB: 10, maxB: 99, ensureCarry: ensureCarry}
		pattern := regexp.MustCompile(`\\\((\d+) \+ (\d+) = \?\\\)`)
		for i := 0; i < 500; i++ {
			p := gen.Generate(generator.GeneratorContext{Difficulty: 0.5})
			m := pattern.FindStringSubmatch(p.Question)
			if m == nil {
				t.Fatalf("question does not match the addition form: %q", p.Question)
			}
			a, _ := strconv.Atoi(m[1])
			b, _ := strconv.Atoi(m[2])
			if want := strconv.Itoa(a + b); p.Answer != want {
				t.Fatalf("%d + %d answered %q, want %q", a, b, p.Answer, want)
			}
			if p.Facts["sum"] != p.Answer || p.Facts["a"] != m[1] || p.Facts["b"] != m[2] {
				t.Fatalf("facts disagree with the printed sum: %+v for %q", p.Facts, p.Question)
			}
		}
	}
}

func testIsPrime(n int) bool {
	if n < 2 {
		return false
	}
	for d := 2; d*d <= n; d++ {
		if n%d == 0 {
			return false
		}
	}
	return true
}

func TestCompositeProblemCountsComposites(t *testing.T) {
	pattern := regexp.MustCompile(`\((\w)\) (\d+)`)
	for i := 0; i < 200; i++ {
		p := compositeProblem()
		matches := pattern.FindAllStringSubmatch(p.Question, -1)
		if len(matches) != 8 {
			t.Fatalf("expected eight candidates, got %d in %q", len(matches), p.Question)
		}
		composites := 0
		for _, m := range matches {
			n, _ := strconv.Atoi(m[2])
			if !testIsPrime(n) {
				composites++
			}
		}
		if want := strconv.Itoa(composites); p.Answer != want {
			t.Fatalf("question %q answered %q, but %d listed numbers are composite", p.Question, p.Answer, composites)
		}
		if p.Facts["count"] != p.Answer {
			t.Fatalf("count fact %q disagrees with answer %q", p.Facts["count"], p.Answer)
		}
	}
}
