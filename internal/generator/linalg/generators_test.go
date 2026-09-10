package linalg

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

func TestVectorConceptGen(t *testing.T)    { fuzzGen(t, &vectorConceptGen{}) }
func TestVectorAddGen(t *testing.T)        { fuzzGen(t, &vectorAddGen{}) }
func TestVectorDotGen(t *testing.T)        { fuzzGen(t, &vectorDotGen{}) }
func TestMatrixConceptGen(t *testing.T)    { fuzzGen(t, &matrixConceptGen{}) }
func TestMatrixAddGen(t *testing.T)        { fuzzGen(t, &matrixAddGen{}) }
func TestMatrixMultGen(t *testing.T)       { fuzzGen(t, &matrixMultGen{}) }
func TestMatrixIdentityGen(t *testing.T)   { fuzzGen(t, &matrixIdentityGen{}) }
func TestDet2x2Gen(t *testing.T)           { fuzzGen(t, &det2x2Gen{}) }
func TestDet3x3Gen(t *testing.T)           { fuzzGen(t, &det3x3Gen{}) }
func TestSystemsMatrixGen(t *testing.T)    { fuzzGen(t, &systemsMatrixGen{}) }
func TestCramerGen(t *testing.T)           { fuzzGen(t, &cramerGen{}) }
func TestEigenConceptGen(t *testing.T)     { fuzzGen(t, &eigenConceptGen{}) }
func TestEigenComputeGen(t *testing.T)     { fuzzGen(t, &eigenComputeGen{}) }
func TestTransformationsGen(t *testing.T)  { fuzzGen(t, &transformationsGen{}) }
func TestSpanGen(t *testing.T)             { fuzzGen(t, &spanGen{}) }
func TestBasisGen(t *testing.T)            { fuzzGen(t, &basisGen{}) }
func TestDiagonalizationGen(t *testing.T)  { fuzzGen(t, &diagonalizationGen{}) }
func TestRankGen(t *testing.T)             { fuzzGen(t, &rankGen{}) }
func TestCosineSimilarityGen(t *testing.T) { fuzzGen(t, &cosineSimilarityGen{}) }
func TestParametricGen(t *testing.T)       { fuzzGen(t, &parametricGen{}) }

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
