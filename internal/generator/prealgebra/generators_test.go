package prealgebra

import (
	"math/rand"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func fuzzGen(t *testing.T, gen generator.Generator) {
	t.Helper()
	for i := 0; i < 100; i++ {
		d := rand.Float64()
		p := gen.Generate(d)
		if p.Question == "" || p.Answer == "" || p.Explanation == "" {
			t.Errorf("empty field at difficulty=%.2f: q=%q a=%q e=%q", d, p.Question, p.Answer, p.Explanation)
		}
	}
}

func TestDecCompareGen(t *testing.T)    { fuzzGen(t, &decCompareGen{}) }
func TestDecAddSubGen(t *testing.T)     { fuzzGen(t, &decAddSubGen{}) }
func TestDecMultGen(t *testing.T)       { fuzzGen(t, &decMultGen{}) }
func TestDecDivGen(t *testing.T)        { fuzzGen(t, &decDivGen{}) }
func TestDecFromFracGen(t *testing.T)   { fuzzGen(t, &decFromFracGen{}) }
func TestDecToFracGen(t *testing.T)     { fuzzGen(t, &decToFracGen{}) }
func TestDecRoundGen(t *testing.T)      { fuzzGen(t, &decRoundGen{}) }
func TestPctConceptGen(t *testing.T)    { fuzzGen(t, &pctConceptGen{}) }
func TestPctToDecGen(t *testing.T)      { fuzzGen(t, &pctToDecGen{}) }
func TestPctFromDecGen(t *testing.T)    { fuzzGen(t, &pctFromDecGen{}) }
func TestPctOfNumberGen(t *testing.T)   { fuzzGen(t, &pctOfNumberGen{}) }
func TestPctFindRateGen(t *testing.T)   { fuzzGen(t, &pctFindRateGen{}) }
func TestPctIncreaseGen(t *testing.T)   { fuzzGen(t, &pctIncreaseGen{}) }
func TestPctDiscountGen(t *testing.T)   { fuzzGen(t, &pctDiscountGen{}) }
func TestPctTaxTipGen(t *testing.T)     { fuzzGen(t, &pctTaxTipGen{}) }
func TestAbsValueGen(t *testing.T)      { fuzzGen(t, &absValueGen{}) }
func TestNegOrderOpsGen(t *testing.T)   { fuzzGen(t, &negOrderOpsGen{}) }
func TestRatioConceptGen(t *testing.T)  { fuzzGen(t, &ratioConceptGen{}) }
func TestRatioSimplifyGen(t *testing.T) { fuzzGen(t, &ratioSimplifyGen{}) }
func TestRatioProportionGen(t *testing.T) { fuzzGen(t, &ratioProportionGen{}) }
func TestRatioRateGen(t *testing.T)     { fuzzGen(t, &ratioRateGen{}) }
func TestRatioScaleGen(t *testing.T)    { fuzzGen(t, &ratioScaleGen{}) }
func TestExpNegGen(t *testing.T)        { fuzzGen(t, &expNegGen{}) }
func TestExpZeroGen(t *testing.T)       { fuzzGen(t, &expZeroGen{}) }
func TestSciNotationGen(t *testing.T)   { fuzzGen(t, &sciNotationGen{}) }
func TestSciNotationOpsGen(t *testing.T) { fuzzGen(t, &sciNotationOpsGen{}) }
func TestVarConceptGen(t *testing.T)    { fuzzGen(t, &varConceptGen{}) }
func TestExprEvalGen(t *testing.T)      { fuzzGen(t, &exprEvalGen{}) }
func TestLikeTermsGen(t *testing.T)     { fuzzGen(t, &likeTermsGen{}) }
func TestDistributeGen(t *testing.T)    { fuzzGen(t, &distributeGen{}) }
func TestEqOneStepAddGen(t *testing.T)  { fuzzGen(t, &eqOneStepAddGen{}) }
func TestEqOneStepMultGen(t *testing.T) { fuzzGen(t, &eqOneStepMultGen{}) }
func TestEqTwoStepGen(t *testing.T)     { fuzzGen(t, &eqTwoStepGen{}) }
func TestEqWordGen(t *testing.T)        { fuzzGen(t, &eqWordGen{}) }
func TestIneqOneStepGen(t *testing.T)   { fuzzGen(t, &ineqOneStepGen{}) }
func TestIneqTwoStepGen(t *testing.T)   { fuzzGen(t, &ineqTwoStepGen{}) }
func TestRealConceptGen(t *testing.T)   { fuzzGen(t, &realConceptGen{}) }
func TestRealPropertiesGen(t *testing.T) { fuzzGen(t, &realPropertiesGen{}) }
func TestTypesGen(t *testing.T)         { fuzzGen(t, &typesGen{}) }
