package discrete

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

func TestPropositionsGen(t *testing.T)   { fuzzGen(t, &propositionsGen{}) }
func TestConnectivesGen(t *testing.T)    { fuzzGen(t, &connectivesGen{}) }
func TestTruthTablesGen(t *testing.T)    { fuzzGen(t, &truthTablesGen{}) }
func TestQuantifiersGen(t *testing.T)    { fuzzGen(t, &quantifiersGen{}) }
func TestSetOpsGen(t *testing.T)         { fuzzGen(t, &setOpsGen{}) }
func TestVennGen(t *testing.T)           { fuzzGen(t, &vennGen{}) }
func TestPermutationsGen(t *testing.T)   { fuzzGen(t, &permutationsGen{}) }
func TestCombinationsGen(t *testing.T)   { fuzzGen(t, &combinationsGen{}) }
func TestPascalGen(t *testing.T)         { fuzzGen(t, &pascalGen{}) }
func TestGraphBasicsGen(t *testing.T)    { fuzzGen(t, &graphBasicsGen{}) }
func TestGraphPathsGen(t *testing.T)     { fuzzGen(t, &graphPathsGen{}) }
func TestTreesGen(t *testing.T)          { fuzzGen(t, &treesGen{}) }
func TestRecurrenceGen(t *testing.T)     { fuzzGen(t, &recurrenceGen{}) }
func TestInductionGen(t *testing.T)      { fuzzGen(t, &inductionGen{}) }
func TestBinomialTheoremGen(t *testing.T) { fuzzGen(t, &binomialTheoremGen{}) }
