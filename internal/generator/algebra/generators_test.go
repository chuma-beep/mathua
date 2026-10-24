package algebra

import (
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/grader"
)

func fuzzGen(t *testing.T, gen generator.Generator) {
	t.Helper()
	gr := grader.NewRouter()
	for i := 0; i < 100; i++ {
		d := rand.Float64()
		p := gen.Generate(generator.GeneratorContext{Difficulty: d})
		if p.Question == "" || p.Answer == "" || p.Explanation == "" {
			t.Errorf("empty field at difficulty=%.2f: q=%q a=%q e=%q", d, p.Question, p.Answer, p.Explanation)
		}
		res := gr.Grade(grader.GradingNumeric, p.Answer, p.Answer)
		if !res.Correct {
			res2 := gr.Grade(grader.GradingMultipleChoice, p.Answer, p.Answer)
			if !res2.Correct {
				t.Errorf("self-grade failed for %q (answer=%q)", p.Question, p.Answer)
			}
		}
	}
}

func TestSlopeGen(t *testing.T)             { fuzzGen(t, &slopeGen{}) }
func TestSlopeInterceptGen(t *testing.T)    { fuzzGen(t, &slopeInterceptGen{}) }
func TestMultiStepEqGen(t *testing.T)       { fuzzGen(t, &multiStepEqGen{}) }
func TestVarsBothSidesGen(t *testing.T)     { fuzzGen(t, &varsBothSidesGen{}) }
func TestVarsBothSidesWordGen(t *testing.T) { fuzzGen(t, &varsBothSidesWordGen{}) }
func TestIneqWordGen(t *testing.T)          { fuzzGen(t, &ineqWordGen{}) }
func TestStdFormGen(t *testing.T)           { fuzzGen(t, &stdFormGen{}) }
func TestParallelPerpGen(t *testing.T)      { fuzzGen(t, &parallelPerpGen{}) }

func TestDifficultyScaling(t *testing.T) {
	gen := &multiStepEqGen{}
	// Low difficulty should produce smaller numbers
	low := gen.Generate(generator.GeneratorContext{Difficulty: 0.1})
	high := gen.Generate(generator.GeneratorContext{Difficulty: 1.0})
	if low.Question == "" || high.Question == "" {
		t.Fatal("empty questions")
	}
	t.Logf("low=%.2f: %s", 0.1, low.Question)
	t.Logf("high=%.2f: %s", 1.0, high.Question)
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
func TestInjectivityGen(t *testing.T)    { fuzzGen(t, &injectivityGen{}) }
func TestPowerFuncGen(t *testing.T)      { fuzzGen(t, &powerFuncGen{}) }
func TestIrrationalFuncGen(t *testing.T) { fuzzGen(t, &irrationalFuncGen{}) }

// Nine reports on 23 October 2026 said the same thing: the question asked for k
// and the equation had no k in it. The answer was unrelated to the printed
// quadratic as well. This pins the equation a learner reads against the answer.
func TestQuadParametricAsksForAParameterItPrints(t *testing.T) {
	// The generator asks the same discriminant condition three ways. Each shape
	// must print the quantity it asks for, and the facts must agree with the
	// answer, so the check switches on the recorded shape instead of assuming
	// one of them.
	for i := 0; i < 100; i++ {
		p := (&quadParametricGen{}).Generate(generator.GeneratorContext{})
		constant, err := strconv.Atoi(p.Facts["const"])
		if err != nil {
			t.Fatalf("const fact %q is not an integer", p.Facts["const"])
		}
		answer, err := strconv.Atoi(p.Answer)
		if err != nil {
			t.Fatalf("answer %q is not an integer", p.Answer)
		}
		if _, err := strconv.Atoi(p.Facts["coef"]); err != nil {
			t.Fatalf("coef fact %q is not an integer", p.Facts["coef"])
		}
		switch p.Facts["shape"] {
		case "askK":
			if !strings.Contains(p.Question, "kx") {
				t.Fatalf("question does not print the parameter it asks for: %q", p.Question)
			}
			if answer*answer != 4*constant {
				t.Fatalf("k=%d gives discriminant %d for x^2 + kx + %d, want zero", answer, answer*answer-4*constant, constant)
			}
			if p.Facts["coef"] != p.Answer || p.Facts["k"] != p.Answer {
				t.Fatalf("facts disagree with the answer: %v", p.Facts)
			}
			if !strings.Contains(p.Question, strconv.Itoa(constant)) {
				t.Fatalf("constant %d is not printed in %q", constant, p.Question)
			}
		case "askC":
			if !strings.Contains(p.Question, p.Facts["coef"]+"x") {
				t.Fatalf("middle coefficient is not printed in %q", p.Question)
			}
			if p.Answer != p.Facts["const"] {
				t.Fatalf("answer %q is not the constant %q", p.Answer, p.Facts["const"])
			}
			coef, _ := strconv.Atoi(p.Facts["coef"])
			if coef*coef != 4*answer {
				t.Fatalf("coef=%d gives discriminant %d for x^2 + %dx + c, want zero", coef, coef*coef-4*answer, coef)
			}
		case "askKNeg":
			if !strings.Contains(p.Question, "- kx") {
				t.Fatalf("question does not print the negative parameter it asks for: %q", p.Question)
			}
			if answer*answer != 4*constant {
				t.Fatalf("k=%d gives discriminant %d for x^2 - kx + %d, want zero", answer, answer*answer-4*constant, constant)
			}
			if p.Facts["coef"] != "-"+p.Answer {
				t.Fatalf("coef fact %q is not the negative of the answer %q", p.Facts["coef"], p.Answer)
			}
		default:
			t.Fatalf("unknown shape %q", p.Facts["shape"])
		}
	}
}

// The boundary of a two-step inequality has to be exact. The earlier generator
// added a remainder to the constant, so the true solution was 27/7 while the
// answer field held a whole number and the explanation printed the boundary as
// an integer, the shape of two live reports.
func TestMultiStepInequalityHasAnExactBoundary(t *testing.T) {
	gen := &multiStepIneqGen{}
	for i := 0; i < 200; i++ {
		p := gen.Generate(generator.GeneratorContext{})
		a, _ := strconv.Atoi(p.Facts["a"])
		b, _ := strconv.Atoi(p.Facts["b"])
		c, _ := strconv.Atoi(p.Facts["c"])
		want, _ := strconv.Atoi(p.Answer)
		if a == 0 {
			t.Fatalf("zero coefficient in %q", p.Question)
		}
		if (c-b)%a != 0 {
			t.Fatalf("%q: c-b = %d is not a multiple of %d, so the boundary is a fraction", p.Question, c-b, a)
		}
		if got := (c - b) / a; got != want {
			t.Fatalf("%q: boundary %d, answer %d", p.Question, got, want)
		}
		if strings.Contains(p.Explanation, "ge") {
			t.Fatalf("%q: explanation still claims a greater-or-equal bound: %q", p.Question, p.Explanation)
		}
		// "Solve:" alone invited "x > 5", which the numeric grader cannot read
		// (reports 228 and 226). The question has to name the number it wants.
		if !strings.Contains(p.Question, "boundary of the solution set") {
			t.Fatalf("%q does not ask for the boundary", p.Question)
		}
	}
}

func TestSystemsWordAnswerIsSmallerFirst(t *testing.T) {
	g := &sysWordGen{}
	for i := 0; i < 200; i++ {
		p := g.Generate(generator.GeneratorContext{Difficulty: rand.Float64(), Seed: rand.Int63()})
		parts := strings.Split(p.Answer, ",")
		if len(parts) != 2 {
			t.Fatalf("answer %q is not a pair: %q", p.Answer, p.Question)
		}
		small, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
		big, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err1 != nil || err2 != nil {
			t.Fatalf("answer %q is not two integers: %q", p.Answer, p.Question)
		}
		// The question asks for the smaller number first, so the answer has to
		// agree with it or a learner who reads it is marked wrong.
		if small >= big {
			t.Fatalf("answer %q is not smaller first: %q", p.Answer, p.Question)
		}
		sum, _ := strconv.Atoi(p.Facts["sum"])
		diff, _ := strconv.Atoi(p.Facts["diff"])
		if small+big != sum {
			t.Fatalf("answer %q does not sum to %d: %q", p.Answer, sum, p.Question)
		}
		if big-small != diff {
			t.Fatalf("answer %q does not differ by %d: %q", p.Answer, diff, p.Question)
		}
		if diff <= 0 {
			t.Fatalf("difference %d is not positive: %q", diff, p.Question)
		}
	}
}

func TestQuadraticRootsAreOneNumberWhenRepeated(t *testing.T) {
	g := &quadSolveFactorGen{}
	repeated := 0
	for i := 0; i < 400; i++ {
		p := g.Generate(generator.GeneratorContext{Difficulty: rand.Float64(), Seed: rand.Int63()})
		r1, _ := strconv.Atoi(p.Facts["r1"])
		r2, _ := strconv.Atoi(p.Facts["r2"])
		parts := strings.Split(p.Answer, ",")
		if r1 == r2 {
			repeated++
			if len(parts) != 1 {
				t.Fatalf("repeated root %d answered as %q: %q", r1, p.Answer, p.Question)
			}
			if parts[0] != strconv.Itoa(r1) {
				t.Fatalf("answer %q is not the repeated root %d: %q", p.Answer, r1, p.Question)
			}
			continue
		}
		if len(parts) != 2 {
			t.Fatalf("two distinct roots answered as %q: %q", p.Answer, p.Question)
		}
	}
	if repeated == 0 {
		t.Fatal("no repeated root was drawn, so the case was never checked")
	}
}

func TestSystemQuestionsAskForThePair(t *testing.T) {
	// A tuple answers with both values, so the stem has to say so: "Solve:"
	// alone let a learner stop at x = 2 and be marked wrong (report 147).
	pair := regexp.MustCompile(`^\((\d+),(\d+)\)$`)

	for i := 0; i < 100; i++ {
		p := (&sysSubstitutionGen{}).Generate(generator.GeneratorContext{})
		if !strings.Contains(p.Question, "Give both values as") {
			t.Fatalf("substitution stem does not ask for the pair: %q", p.Question)
		}
		m := pair.FindStringSubmatch(p.Answer)
		if m == nil {
			t.Fatalf("substitution answer %q is not a pair", p.Answer)
		}
		y, _ := strconv.Atoi(p.Facts["y"])
		sum, _ := strconv.Atoi(p.Facts["sum"])
		x, _ := strconv.Atoi(m[1])
		if got := x + y; got != sum {
			t.Fatalf("x + y = %d, want the sum in the stem %d", got, sum)
		}
	}

	for i := 0; i < 100; i++ {
		p := (&sysEliminationGen{}).Generate(generator.GeneratorContext{})
		if !strings.Contains(p.Question, "Give both values as") {
			t.Fatalf("elimination stem does not ask for the pair: %q", p.Question)
		}
		m := pair.FindStringSubmatch(p.Answer)
		if m == nil {
			t.Fatalf("elimination answer %q is not a pair", p.Answer)
		}
		a, _ := strconv.Atoi(p.Facts["a"])
		b, _ := strconv.Atoi(p.Facts["b"])
		eq1, _ := strconv.Atoi(p.Facts["eq1"])
		x, _ := strconv.Atoi(m[1])
		y, _ := strconv.Atoi(m[2])
		if got := a*x + b*y; got != eq1 {
			t.Fatalf("%dx + %dy = %d, want %d", a, b, got, eq1)
		}
		c, _ := strconv.Atoi(p.Facts["c"])
		eq2, _ := strconv.Atoi(p.Facts["eq2"])
		if c == a {
			// Both equations identical is an underdetermined system, and
			// report 229 quoted exactly that: "6x + 4y = 50, 6x + 4y = 50".
			t.Fatalf("both equations use x coefficient %d: %q", a, p.Question)
		}
		if got := c*x + b*y; got != eq2 {
			t.Fatalf("%dx + %dy = %d, want %d", c, b, got, eq2)
		}
	}
}

// The stem used to be built with "x^{2} + %dx + %d", which printed
// "x^2 + -2x + 0" for negative or zero coefficients (report 224).
func TestPolyRootsPrintsReadableCoefficients(t *testing.T) {
	gen := &polyRootsGen{}
	for i := 0; i < 300; i++ {
		p := gen.Generate(generator.GeneratorContext{})
		if strings.Contains(p.Question, "+ -") || strings.Contains(p.Question, "+ 0") {
			t.Fatalf("unreadable quadratic: %q", p.Question)
		}
		b, _ := strconv.Atoi(p.Facts["b"])
		c, _ := strconv.Atoi(p.Facts["c"])
		r1, _ := strconv.Atoi(p.Facts["r1"])
		r2, _ := strconv.Atoi(p.Facts["r2"])
		if r1 != r2 && p.Answer != p.Facts["r1"]+","+p.Facts["r2"] {
			t.Fatalf("answer %q does not list both roots %d and %d", p.Answer, r1, r2)
		}
		if r1 == r2 && p.Answer != p.Facts["r1"] {
			t.Fatalf("double root %d answered as %q", r1, p.Answer)
		}
		for _, r := range []int{r1, r2} {
			if r*r+b*r+c != 0 {
				t.Fatalf("root %d does not solve x^2 + %dx + %d", r, b, c)
			}
		}
	}
}

// Report 196: the AC-method generator served x^2 + 9x + 1 with (x+1)(x+1),
// which multiplies out to x^2 + 2x + 1. The concept is "ax^2 + bx + c", so the
// leading coefficient must also not be 1: otherwise the "AC" step is a no-op.
func TestACMethodTrinomialMultipliesOut(t *testing.T) {
	q := regexp.MustCompile(`^Factor: \\\((\d+)x\^\{2\} \+ (\d+)x \+ (\d+)\\\)$`)
	ans := regexp.MustCompile(`^\((\d+)x \+ (\d+)\)\(x \+ (\d+)\)$`)
	g := &factorACMethodGen{}
	wideLead := 0
	for i := 0; i < 300; i++ {
		p := g.Generate(generator.GeneratorContext{Seed: int64(i)})
		qm := q.FindStringSubmatch(p.Question)
		if qm == nil {
			t.Fatalf("question not in ax^2 + bx + c form: %q", p.Question)
		}
		am := ans.FindStringSubmatch(p.Answer)
		if am == nil {
			t.Fatalf("answer not in factored form: %q", p.Answer)
		}
		a, _ := strconv.Atoi(qm[1])
		b, _ := strconv.Atoi(qm[2])
		c, _ := strconv.Atoi(qm[3])
		fa, _ := strconv.Atoi(am[1])
		fq, _ := strconv.Atoi(am[2])
		fs, _ := strconv.Atoi(am[3])
		if a != fa {
			t.Fatalf("question leads with %d but the answer leads with %d: %q", a, fa, p.Question)
		}
		if a > 1 {
			wideLead++
		}
		// (ax + q)(x + s) = ax^2 + (as + q)x + qs
		if b != a*fs+fq || c != fq*fs {
			t.Fatalf("printed trinomial does not multiply out: %q = %q gives b=%d c=%d, want b=%d c=%d",
				p.Question, p.Answer, a*fs+fq, fq*fs, b, c)
		}
		// The AC step needs a pair multiplying to ac and adding to b.
		if a*fs+fq != b || (a*fs)*fq != a*c {
			t.Fatalf("AC pair as=%d and q=%d do not fit ac=%d and b=%d", a*fs, fq, a*c, b)
		}
	}
	if wideLead != 300 {
		t.Fatalf("only %d of 300 draws had a leading coefficient above 1", wideLead)
	}
}
