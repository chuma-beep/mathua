package complex

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

func TestConceptGen(t *testing.T)      { fuzzGen(t, &conceptGen{}) }
func TestAddSubGen(t *testing.T)       { fuzzGen(t, &addSubGen{}) }
func TestMultGen(t *testing.T)         { fuzzGen(t, &multGen{}) }
func TestConjugateGen(t *testing.T)    { fuzzGen(t, &conjugateGen{}) }
func TestDivideGen(t *testing.T)       { fuzzGen(t, &divideGen{}) }
func TestPolarGen(t *testing.T)        { fuzzGen(t, &polarGen{}) }
func TestDeMoivreGen(t *testing.T)     { fuzzGen(t, &deMoivreGen{}) }
func TestRootsGen(t *testing.T)        { fuzzGen(t, &rootsGen{}) }
func TestExponentialGen(t *testing.T)  { fuzzGen(t, &exponentialGen{}) }
func TestInequalitiesGen(t *testing.T) { fuzzGen(t, &inequalitiesGen{}) }

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
