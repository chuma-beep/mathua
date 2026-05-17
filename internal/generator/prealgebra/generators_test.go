package prealgebra

import (
	"math/rand"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/grader"
	"github.com/chuma-beep/mathua/internal/verify"
)

func fuzzGen(t *testing.T, gen generator.Generator, gtype grader.GradingType) {
	t.Helper()
	gr := grader.NewRouter()
	for i := 0; i < 100; i++ {
		p := gen.Generate(rand.Float64())
		if p.Question == "" || p.Answer == "" || p.Explanation == "" {
			t.Errorf("empty field")
		}
		res := gr.Grade(gtype, p.Answer, p.Answer)
		if !res.Correct {
			t.Errorf("self-grade failed: gtype=%s a=%q", gtype, p.Answer)
		}
		if expr, ok := verify.ExtractExpr(p.Question); ok {
			if err := verify.CheckExpr(expr, p.Answer); err != nil {
				t.Errorf("%s -> %s: %v", p.Question, p.Answer, err)
			}
		}
	}
}

func fuzzNoSelf(t *testing.T, gen generator.Generator) {
	t.Helper()
	for i := 0; i < 100; i++ {
		p := gen.Generate(rand.Float64())
		if p.Question == "" || p.Answer == "" || p.Explanation == "" {
			t.Errorf("empty field")
		}
	}
}

func fuzzGrade(t *testing.T, gen generator.Generator) {
	t.Helper()
	gg, ok := gen.(generator.GradedGenerator)
	if !ok {
		t.Fatal("not a GradedGenerator")
	}
	for i := 0; i < 100; i++ {
		p := gen.Generate(rand.Float64())
		if p.Question == "" || p.Answer == "" || p.Explanation == "" {
			t.Errorf("empty field")
		}
		res := gg.Grade(p.Answer, p.Answer)
		if !res.Correct {
			t.Errorf("self-grade via Grade() failed: a=%q", p.Answer)
		}
	}
}

func TestDecCompare(t *testing.T)  { fuzzGen(t, &decCompareGen{}, grader.GradingComparison) }
func TestDecAdd(t *testing.T)      { fuzzGen(t, &decAddSubGen{op: "+"}, grader.GradingNumeric) }
func TestDecSub(t *testing.T)      { fuzzGen(t, &decAddSubGen{op: "-"}, grader.GradingNumeric) }
func TestDecMult(t *testing.T)     { fuzzGen(t, &decMultGen{}, grader.GradingNumeric) }
func TestDecDiv(t *testing.T)      { fuzzGen(t, &decDivGen{}, grader.GradingNumeric) }
func TestDecFromFrac(t *testing.T) { fuzzGen(t, &decFromFracGen{}, grader.GradingNumeric) }
func TestDecToFrac(t *testing.T)   { fuzzGen(t, &decToFracGen{}, grader.GradingNumeric) }
func TestDecRound(t *testing.T)    { fuzzGen(t, &decRoundGen{}, grader.GradingNumeric) }

func TestPctConcept(t *testing.T)  { fuzzGen(t, &pctConceptGen{}, grader.GradingNumeric) }
func TestPctToDec(t *testing.T)    { fuzzGen(t, &pctToDecGen{}, grader.GradingNumeric) }
func TestPctFromDec(t *testing.T)  { fuzzGen(t, &pctFromDecGen{}, grader.GradingNumeric) }
func TestPctOfNumber(t *testing.T) { fuzzGen(t, &pctOfNumberGen{}, grader.GradingNumeric) }
func TestPctFindRate(t *testing.T) { fuzzGen(t, &pctFindRateGen{}, grader.GradingNumeric) }
func TestPctIncrease(t *testing.T) { fuzzGen(t, &pctIncreaseGen{}, grader.GradingNumeric) }
func TestPctDiscount(t *testing.T) { fuzzGen(t, &pctDiscountGen{}, grader.GradingNumeric) }
func TestPctTaxTip(t *testing.T)   { fuzzGen(t, &pctTaxTipGen{}, grader.GradingNumeric) }

func TestAbsValue(t *testing.T)    { fuzzGen(t, &absValueGen{}, grader.GradingNumeric) }
func TestNegOrderOps(t *testing.T) { fuzzGen(t, &negOrderOpsGen{}, grader.GradingNumeric) }

func TestRatioConcept(t *testing.T)    { fuzzGrade(t, &ratioConceptGen{}) }
func TestRatioSimplify(t *testing.T)   { fuzzGrade(t, &ratioSimplifyGen{}) }
func TestRatioProportion(t *testing.T) { fuzzGen(t, &ratioProportionGen{}, grader.GradingNumeric) }
func TestRatioRate(t *testing.T)       { fuzzGen(t, &ratioRateGen{}, grader.GradingNumeric) }
func TestRatioScale(t *testing.T)      { fuzzGen(t, &ratioScaleGen{}, grader.GradingNumeric) }

func TestExpNeg(t *testing.T)         { fuzzGen(t, &expNegGen{}, grader.GradingNumeric) }
func TestExpZero(t *testing.T)        { fuzzGen(t, &expZeroGen{}, grader.GradingNumeric) }
func TestSciNotation(t *testing.T)    { fuzzGen(t, &sciNotationGen{}, grader.GradingNumeric) }
func TestSciNotationOps(t *testing.T) { fuzzGrade(t, &sciNotationOpsGen{}) }

func TestVarConcept(t *testing.T)    { fuzzGen(t, &varConceptGen{}, grader.GradingNumeric) }
func TestExprEval(t *testing.T)      { fuzzGen(t, &exprEvalGen{}, grader.GradingNumeric) }
func TestLikeTerms(t *testing.T)     { fuzzGen(t, &likeTermsGen{}, grader.GradingPolynomial) }
func TestDistribute(t *testing.T)    { fuzzGen(t, &distributeGen{}, grader.GradingPolynomial) }
func TestEqOneStepAdd(t *testing.T)  { fuzzGen(t, &eqOneStepAddGen{}, grader.GradingNumeric) }
func TestEqOneStepMult(t *testing.T) { fuzzGen(t, &eqOneStepMultGen{}, grader.GradingNumeric) }
func TestEqTwoStep(t *testing.T)     { fuzzGen(t, &eqTwoStepGen{}, grader.GradingNumeric) }
func TestEqWord(t *testing.T)        { fuzzGen(t, &eqWordGen{}, grader.GradingNumeric) }
func TestIneqOneStep(t *testing.T)   { fuzzGen(t, &ineqOneStepGen{}, grader.GradingNumeric) }
func TestIneqTwoStep(t *testing.T)   { fuzzGen(t, &ineqTwoStepGen{}, grader.GradingNumeric) }
