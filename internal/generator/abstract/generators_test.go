package abstract

import (
	"math/rand"
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

func TestGroupDefGen(t *testing.T)       { fuzzGen(t, &groupDefGen{}) }
func TestGroupExamplesGen(t *testing.T)  { fuzzGen(t, &groupExamplesGen{}) }
func TestSubgroupGen(t *testing.T)       { fuzzGen(t, &subgroupGen{}) }
func TestRingGen(t *testing.T)           { fuzzGen(t, &ringGen{}) }
func TestHomomorphismGen(t *testing.T)   { fuzzGen(t, &homomorphismGen{}) }
func TestFieldGen(t *testing.T)          { fuzzGen(t, &fieldGen{}) }
func TestModuleGen(t *testing.T)         { fuzzGen(t, &moduleGen{}) }

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

