package odes

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

func TestConceptGen(t *testing.T)        { fuzzGen(t, &conceptGen{}) }
func TestSeparableGen(t *testing.T)      { fuzzGen(t, &separableGen{}) }
func TestLinearFirstGen(t *testing.T)    { fuzzGen(t, &linearFirstGen{}) }
func TestExactGen(t *testing.T)          { fuzzGen(t, &exactGen{}) }
func TestHomogeneousGen(t *testing.T)    { fuzzGen(t, &homogeneousGen{}) }
func TestNonhomogeneousGen(t *testing.T) { fuzzGen(t, &nonhomogeneousGen{}) }
func TestLaplaceGen(t *testing.T)        { fuzzGen(t, &laplaceGen{}) }
func TestSystemsGen(t *testing.T)        { fuzzGen(t, &systemsGen{}) }

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
