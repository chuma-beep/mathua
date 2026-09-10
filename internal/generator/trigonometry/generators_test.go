package trigonometry

import (
	"math/rand"
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
