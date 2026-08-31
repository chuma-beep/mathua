package arithmetic

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

func TestAddGen(t *testing.T)            { fuzzGen(t, &addGen{}) }
func TestAddWordGen(t *testing.T)        { fuzzGen(t, &addWordGen{}) }
func TestSubGen(t *testing.T)            { fuzzGen(t, &subGen{}) }
func TestSubBorrowGen(t *testing.T)      { fuzzGen(t, &subBorrowGen{}) }
func TestSubWordGen(t *testing.T)        { fuzzGen(t, &subWordGen{}) }
func TestPlaceValueGen(t *testing.T)     { fuzzGen(t, &placeValueGen{}) }
func TestRoundGen(t *testing.T)          { fuzzGen(t, &roundGen{}) }
func TestMultConceptGen(t *testing.T)    { fuzzGen(t, &multConceptGen{}) }
func TestMultByGen(t *testing.T)         { fuzzGen(t, &multByGen{}) }
func TestMultTablesGen(t *testing.T)     { fuzzGen(t, &multTablesGen{}) }
func TestMultGen(t *testing.T)           { fuzzGen(t, &multGen{}) }
func TestMultWordGen(t *testing.T)       { fuzzGen(t, &multWordGen{}) }
func TestDivConceptGen(t *testing.T)     { fuzzGen(t, &divConceptGen{}) }
func TestDivBasicGen(t *testing.T)       { fuzzGen(t, &divBasicGen{}) }
func TestDivRemainderGen(t *testing.T)   { fuzzGen(t, &divRemainderGen{}) }
func TestDivLongGen(t *testing.T)        { fuzzGen(t, &divLongGen{}) }
func TestDivWordGen(t *testing.T)        { fuzzGen(t, &divWordGen{}) }
func TestFactorFindGen(t *testing.T)     { fuzzGen(t, &factorFindGen{}) }
func TestPrimeGen(t *testing.T)          { fuzzGen(t, &primeGen{}) }
func TestPrimeFactGen(t *testing.T)      { fuzzGen(t, &primeFactGen{}) }
func TestGcfGen(t *testing.T)            { fuzzGen(t, &gcfGen{}) }
func TestLcmGen(t *testing.T)            { fuzzGen(t, &lcmGen{}) }
func TestCompositeGen(t *testing.T)      { fuzzGen(t, &compositeGen{}) }
func TestExpConceptGen(t *testing.T)     { fuzzGen(t, &expConceptGen{}) }
func TestExpEvalGen(t *testing.T)        { fuzzGen(t, &expEvalGen{}) }
func TestExpProductRuleGen(t *testing.T) { fuzzGen(t, &expProductRuleGen{}) }
func TestExpQuotientRuleGen(t *testing.T) { fuzzGen(t, &expQuotientRuleGen{}) }
func TestExpPowerRuleGen(t *testing.T)   { fuzzGen(t, &expPowerRuleGen{}) }
func TestSqrtPerfectGen(t *testing.T)    { fuzzGen(t, &sqrtPerfectGen{}) }
func TestSqrtSimplifyGen(t *testing.T)   { fuzzGen(t, &sqrtSimplifyGen{}) }
func TestNegNumberLineGen(t *testing.T)  { fuzzGen(t, &negNumberLineGen{}) }
func TestNegAddSubGen(t *testing.T)      { fuzzGen(t, &negAddSubGen{}) }
func TestNegMultDivGen(t *testing.T)     { fuzzGen(t, &negMultDivGen{}) }
func TestOrderOpsGen(t *testing.T)       { fuzzGen(t, &orderOpsGen{}) }
func TestDecIntroGen(t *testing.T)       { fuzzGen(t, &decIntroGen{}) }

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

