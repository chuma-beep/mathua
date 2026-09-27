package trigonometry

import (
	"math/rand"
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

func TestUnitCircleDifficulty(t *testing.T) {
	gen := &unitCircleGen{}
	// Low difficulty should prefer quadrant I angles
	lowAngles := make(map[string]int)
	for i := 0; i < 50; i++ {
		p := gen.Generate(generator.GeneratorContext{Difficulty: 0.1})
		lowAngles[p.Question]++
	}
	highAngles := make(map[string]int)
	for i := 0; i < 50; i++ {
		p := gen.Generate(generator.GeneratorContext{Difficulty: 1.0})
		highAngles[p.Question]++
	}
	t.Logf("low difficulty produced %d distinct questions", len(lowAngles))
	t.Logf("high difficulty produced %d distinct questions", len(highAngles))
}

func TestSinCosDefDifficulty(t *testing.T) {
	gen := &sinCosDefGen{}
	for i := 0; i < 20; i++ {
		p := gen.Generate(generator.GeneratorContext{Difficulty: 0.1})
		if p.Question == "" {
			t.Fatal("empty question at low difficulty")
		}
	}
	for i := 0; i < 20; i++ {
		p := gen.Generate(generator.GeneratorContext{Difficulty: 0.9})
		if p.Question == "" {
			t.Fatal("empty question at high difficulty")
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

func TestDifficultyScaling(t *testing.T) {
	reg := generator.NewRegistry()
	Register(reg)
	for _, id := range reg.Concepts() {
		gen, _ := reg.Get(id)
		lowQs := make(map[string]bool)
		highQs := make(map[string]bool)
		for i := 0; i < 30; i++ {
			pLow := gen.Generate(generator.GeneratorContext{Difficulty: 0.1})
			pHigh := gen.Generate(generator.GeneratorContext{Difficulty: 0.9})
			if pLow.Question == "" || pHigh.Question == "" {
				t.Fatalf("empty question for %s at low/high difficulty", id)
			}
			lowQs[pLow.Question] = true
			highQs[pHigh.Question] = true
		}
		// For generators that scale, high difficulty should produce at least as much variety
		// At minimum, ensure both produce output without panic; log if identical (unscaled)
		if len(lowQs) == 1 && len(highQs) == 1 {
			for k := range lowQs {
				if highQs[k] {
					t.Logf("generator %s appears unscaled: low and high produce same single question %q", id, k)
				}
			}
		}
	}
}

// TestFuncGensDiscriminate is the P3.1 quality bar (mirrors
// TestPolynomialRingGenDiscriminates): every new trig-function easy pool must
// discriminate (>=2 distinct answers) and must contain production
// (typed-answer) questions, never yes/no-only, over 150 samples at
// Difficulty 0.12.
func TestFuncGensDiscriminate(t *testing.T) {
	rand.Seed(31)
	gens := map[string]generator.Generator{
		"trig.func.sine":         &sineFuncGen{},
		"trig.func.cosine":       &cosineFuncGen{},
		"trig.func.tangent":      &tangentFuncGen{},
		"trig.func.cotangent":    &cotangentFuncGen{},
		"trig.func.secant":       &secantFuncGen{},
		"trig.func.cosecant":     &cosecantFuncGen{},
		"trig.func.arcsine":      &arcsineFuncGen{},
		"trig.func.arccosine":    &arccosineFuncGen{},
		"trig.func.arctangent":   &arctangentFuncGen{},
		"trig.func.arccotangent": &arccotangentFuncGen{},
	}
	const n = 150
	for id, gen := range gens {
		answers := map[string]bool{}
		typed := 0
		for i := 0; i < n; i++ {
			p := gen.Generate(generator.GeneratorContext{Difficulty: 0.12})
			if p.Question == "" || p.Answer == "" || p.Explanation == "" {
				t.Fatalf("%s: empty field at sample %d: q=%q a=%q e=%q", id, i, p.Question, p.Answer, p.Explanation)
			}
			answers[strings.ToLower(strings.TrimSpace(p.Answer))] = true
			if !strings.Contains(p.Question, "(yes/no)") {
				typed++
			}
		}
		if len(answers) < 2 {
			t.Errorf("%s: easy pool non-discriminating: %d distinct answers over %d samples", id, len(answers), n)
		}
		if typed == 0 {
			t.Errorf("%s: easy pool has no production (typed-answer) questions: recognition-only", id)
		}
	}
}

// TestFuncGensHardGating checks the scale>3 hard-pool gate: hard-only answer
// forms must never appear at Difficulty 0.12 but must appear at high
// difficulty.
func TestFuncGensHardGating(t *testing.T) {
	rand.Seed(77)
	gens := map[string]generator.Generator{
		"trig.func.sine":         &sineFuncGen{},
		"trig.func.cosine":       &cosineFuncGen{},
		"trig.func.tangent":      &tangentFuncGen{},
		"trig.func.cotangent":    &cotangentFuncGen{},
		"trig.func.secant":       &secantFuncGen{},
		"trig.func.cosecant":     &cosecantFuncGen{},
		"trig.func.arcsine":      &arcsineFuncGen{},
		"trig.func.arccosine":    &arccosineFuncGen{},
		"trig.func.arctangent":   &arctangentFuncGen{},
		"trig.func.arccotangent": &arccotangentFuncGen{},
	}
	// One hard-only answer form per generator (absent from its easy table).
	hardOnly := map[string]string{
		"trig.func.sine":         "-sin(x)",
		"trig.func.cosine":       "-cos(x)",
		"trig.func.tangent":      "+infinity",
		"trig.func.cotangent":    "decreasing",
		"trig.func.secant":       "ln|sec(x)+tan(x)|",
		"trig.func.cosecant":     "-ln|csc(x)+cot(x)|",
		"trig.func.arcsine":      "-pi/2",
		"trig.func.arccosine":    "pi",
		"trig.func.arctangent":   "-pi/4",
		"trig.func.arccotangent": "3pi/4",
	}
	for id, gen := range gens {
		for i := 0; i < 150; i++ {
			p := gen.Generate(generator.GeneratorContext{Difficulty: 0.12})
			if p.Answer == hardOnly[id] {
				t.Fatalf("%s: hard-only answer %q leaked into easy pool", id, p.Answer)
			}
		}
		seen := false
		for i := 0; i < 600; i++ {
			p := gen.Generate(generator.GeneratorContext{Difficulty: 0.9})
			if p.Answer == hardOnly[id] {
				seen = true
				break
			}
		}
		if !seen {
			t.Errorf("%s: hard-only answer %q never appeared at high difficulty", id, hardOnly[id])
		}
	}
}

// TestFuncGensASCIIAndSelfGrade requires ASCII-safe answers that grade
// against themselves through the declared multiple_choice router.
func TestFuncGensASCIIAndSelfGrade(t *testing.T) {
	rand.Seed(99)
	reg := generator.NewRegistry()
	Register(reg)
	r := grader.NewRouter()
	for _, id := range []string{
		"trig.func.sine", "trig.func.cosine", "trig.func.tangent",
		"trig.func.cotangent", "trig.func.secant", "trig.func.cosecant",
		"trig.func.arcsine", "trig.func.arccosine", "trig.func.arctangent",
		"trig.func.arccotangent",
	} {
		gen, err := reg.Get(id)
		if err != nil {
			t.Fatalf("registered concept %q not retrievable: %v", id, err)
		}
		forms := map[string]bool{}
		for i := 0; i < 300; i++ {
			d := 0.12
			if i%2 == 0 {
				d = 0.9
			}
			p := gen.Generate(generator.GeneratorContext{Difficulty: d})
			forms[p.Answer] = true
			for j, c := range []byte(p.Answer) {
				_ = j
				if c > 127 {
					t.Errorf("%s: non-ASCII byte in answer %q", id, p.Answer)
					break
				}
			}
		}
		for a := range forms {
			if res := r.Grade(grader.GradingMultipleChoice, a, a); !res.Correct {
				t.Errorf("%s: answer form %q does not grade against itself", id, a)
			}
		}
	}
}
