package calculus

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
			t.Errorf("empty field: q=%q a=%q e=%q", p.Question, p.Answer, p.Explanation)
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

func TestLimitConcept(t *testing.T)    { fuzzGen(t, &limitConceptGen{}, grader.GradingNumeric) }
func TestLimitNumeric(t *testing.T)    { fuzzGen(t, &limitNumericGen{}, grader.GradingNumeric) }
func TestLimitProperties(t *testing.T) { fuzzGen(t, &limitPropertiesGen{}, grader.GradingNumeric) }
func TestLimitInfinity(t *testing.T)   { fuzzGen(t, &limitInfinityGen{}, grader.GradingNumeric) }
func TestLimitContinuity(t *testing.T) {
	fuzzGen(t, &limitContinuityGen{}, grader.GradingMultipleChoice)
}

func TestDerivConcept(t *testing.T)   { fuzzGen(t, &derivConceptGen{}, grader.GradingNumeric) }
func TestDerivPowerRule(t *testing.T) { fuzzGen(t, &derivPowerRuleGen{}, grader.GradingPolynomial) }
func TestDerivSumRule(t *testing.T)   { fuzzGen(t, &derivSumRuleGen{}, grader.GradingPolynomial) }
func TestDerivProductRule(t *testing.T) {
	fuzzGen(t, &derivProductRuleGen{}, grader.GradingMultipleChoice)
}
func TestDerivQuotientRule(t *testing.T) {
	fuzzGen(t, &derivQuotientRuleGen{}, grader.GradingNumeric)
}
func TestDerivChainRule(t *testing.T) { fuzzGen(t, &derivChainRuleGen{}, grader.GradingMultipleChoice) }
func TestDerivTrig(t *testing.T)      { fuzzGen(t, &derivTrigGen{}, grader.GradingMultipleChoice) }
func TestDerivExpLog(t *testing.T)    { fuzzGen(t, &derivExpLogGen{}, grader.GradingMultipleChoice) }
func TestDerivApplications(t *testing.T) {
	fuzzGen(t, &derivApplicationsGen{}, grader.GradingNumeric)
}
func TestDerivOptimization(t *testing.T) {
	fuzzGen(t, &derivOptimizationGen{}, grader.GradingNumeric)
}

func TestIntegralIndefinite(t *testing.T) {
	fuzzGen(t, &integralIndefiniteGen{}, grader.GradingMultipleChoice)
}
func TestIntegralPowerRule(t *testing.T) {
	fuzzGen(t, &integralPowerRuleGen{}, grader.GradingPolynomial)
}
func TestIntegralSubstitution(t *testing.T) {
	fuzzGen(t, &integralSubstitutionGen{}, grader.GradingMultipleChoice)
}
func TestIntegralDefinite(t *testing.T) {
	fuzzGen(t, &integralDefiniteGen{}, grader.GradingNumeric)
}
func TestIntegralFTC(t *testing.T) { fuzzGen(t, &integralFTCGen{}, grader.GradingMultipleChoice) }
func TestIntegralAreaBetween(t *testing.T) {
	fuzzGen(t, &integralAreaBetweenGen{}, grader.GradingNumeric)
}
func TestIntegralVolume(t *testing.T) {
	fuzzGen(t, &integralVolumeGen{}, grader.GradingMultipleChoice)
}

func TestDerivImplicit(t *testing.T) { fuzzGen(t, &derivImplicitGen{}, grader.GradingNumeric) }
func TestDerivRelatedRates(t *testing.T) {
	fuzzGen(t, &derivRelatedRatesGen{}, grader.GradingMultipleChoice)
}
func TestIntegralParts(t *testing.T) {
	fuzzGen(t, &integralPartsGen{}, grader.GradingMultipleChoice)
}
func TestIntegralPartialFractions(t *testing.T) {
	fuzzGen(t, &integralPartialFractionsGen{}, grader.GradingMultipleChoice)
}
