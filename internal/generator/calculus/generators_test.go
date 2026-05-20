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
	fuzzGen(t, &derivProductRuleGen{}, grader.GradingExpression)
}
func TestDerivQuotientRule(t *testing.T) {
	fuzzGen(t, &derivQuotientRuleGen{}, grader.GradingExpression)
}
func TestDerivChainRule(t *testing.T) { fuzzGen(t, &derivChainRuleGen{}, grader.GradingExpression) }
func TestDerivTrig(t *testing.T)      { fuzzGen(t, &derivTrigGen{}, grader.GradingExpression) }
func TestDerivExpLog(t *testing.T)    { fuzzGen(t, &derivExpLogGen{}, grader.GradingExpression) }
func TestDerivApplications(t *testing.T) {
	fuzzGen(t, &derivApplicationsGen{}, grader.GradingNumeric)
}
func TestDerivOptimization(t *testing.T) {
	fuzzGen(t, &derivOptimizationGen{}, grader.GradingNumeric)
}

func TestIntegralIndefinite(t *testing.T) {
	fuzzGen(t, &integralIndefiniteGen{}, grader.GradingExpression)
}
func TestIntegralPowerRule(t *testing.T) {
	fuzzGen(t, &integralPowerRuleGen{}, grader.GradingExpression)
}
func TestIntegralSubstitution(t *testing.T) {
	fuzzGen(t, &integralSubstitutionGen{}, grader.GradingExpression)
}
func TestIntegralDefinite(t *testing.T) {
	fuzzGen(t, &integralDefiniteGen{}, grader.GradingNumeric)
}
func TestIntegralFTC(t *testing.T) { fuzzGen(t, &integralFTCGen{}, grader.GradingExpression) }
func TestIntegralAreaBetween(t *testing.T) {
	fuzzGen(t, &integralAreaBetweenGen{}, grader.GradingNumeric)
}
func TestIntegralVolume(t *testing.T) {
	fuzzGen(t, &integralVolumeGen{}, grader.GradingExpression)
}

func TestDerivImplicit(t *testing.T) { fuzzGen(t, &derivImplicitGen{}, grader.GradingExpression) }
func TestDerivRelatedRates(t *testing.T) {
	fuzzGen(t, &derivRelatedRatesGen{}, grader.GradingExpression)
}
func TestIntegralParts(t *testing.T) {
	fuzzGen(t, &integralPartsGen{}, grader.GradingExpression)
}
func TestIntegralPartialFractions(t *testing.T) {
	fuzzGen(t, &integralPartialFractionsGen{}, grader.GradingExpression)
}

func TestRolle(t *testing.T)          { fuzzGen(t, &rolleGen{}, grader.GradingMultipleChoice) }
func TestMVT(t *testing.T)            { fuzzGen(t, &mvtGen{}, grader.GradingMultipleChoice) }
func TestCauchyMVT(t *testing.T)      { fuzzGen(t, &cauchyMVTGen{}, grader.GradingMultipleChoice) }
func TestFermat(t *testing.T)         { fuzzGen(t, &fermatGen{}, grader.GradingMultipleChoice) }
func TestDifferenceQuotient(t *testing.T) {
	fuzzGen(t, &differenceQuotientGen{}, grader.GradingNumeric)
}
func TestDifferential(t *testing.T) { fuzzGen(t, &differentialGen{}, grader.GradingNumeric) }
func TestConvexity(t *testing.T)    { fuzzGen(t, &convexityGen{}, grader.GradingMultipleChoice) }
func TestNonDiff(t *testing.T)      { fuzzGen(t, &nonDiffGen{}, grader.GradingMultipleChoice) }
func TestPartialDeriv(t *testing.T) {
	fuzzGen(t, &partialDerivGen{}, grader.GradingMultipleChoice)
}

func TestArcLength(t *testing.T)           { fuzzGen(t, &arcLengthGen{}, grader.GradingMultipleChoice) }
func TestExpIntegral(t *testing.T)         { fuzzGen(t, &expIntegralGen{}, grader.GradingExpression) }
func TestImproperIntegral(t *testing.T)    { fuzzGen(t, &improperIntegralGen{}, grader.GradingMultipleChoice) }
func TestNumericalIntegral(t *testing.T)   { fuzzGen(t, &numericalIntegralGen{}, grader.GradingMultipleChoice) }
func TestRiemannCriteria(t *testing.T)     { fuzzGen(t, &riemannCriteriaGen{}, grader.GradingMultipleChoice) }
func TestTrigIntegrals(t *testing.T)       { fuzzGen(t, &trigIntegralsGen{}, grader.GradingExpression) }
func TestTrigSubstitution(t *testing.T)    { fuzzGen(t, &trigSubstitutionGen{}, grader.GradingMultipleChoice) }
func TestWeierstrassSub(t *testing.T)      { fuzzGen(t, &weierstrassSubGen{}, grader.GradingMultipleChoice) }

func TestLimitAlgebra(t *testing.T)       { fuzzGen(t, &limitAlgebraGen{}, grader.GradingMultipleChoice) }
func TestAsymptotes(t *testing.T)         { fuzzGen(t, &asymptotesGen{}, grader.GradingMultipleChoice) }
func TestDiscontinuity(t *testing.T)      { fuzzGen(t, &discontinuityGen{}, grader.GradingMultipleChoice) }
func TestIndeterminate(t *testing.T)      { fuzzGen(t, &indeterminateGen{}, grader.GradingMultipleChoice) }
func TestLhopital(t *testing.T)           { fuzzGen(t, &lhopitalGen{}, grader.GradingMultipleChoice) }
func TestBigO(t *testing.T)               { fuzzGen(t, &bigOGen{}, grader.GradingMultipleChoice) }
func TestLittleO(t *testing.T)            { fuzzGen(t, &littleOGen{}, grader.GradingMultipleChoice) }
func TestSqueeze(t *testing.T)            { fuzzGen(t, &squeezeGen{}, grader.GradingMultipleChoice) }
func TestSupremum(t *testing.T)           { fuzzGen(t, &supremumGen{}, grader.GradingMultipleChoice) }
func TestUniformContinuity(t *testing.T)  { fuzzGen(t, &uniformContinuityGen{}, grader.GradingMultipleChoice) }
func TestWeierstrassLimit(t *testing.T)   { fuzzGen(t, &weierstrassLimitGen{}, grader.GradingMultipleChoice) }

func TestSeqConcept(t *testing.T)        { fuzzGen(t, &seqConceptGen{}, grader.GradingMultipleChoice) }
func TestSeqConvergence(t *testing.T)    { fuzzGen(t, &seqConvergenceGen{}, grader.GradingMultipleChoice) }
func TestCauchySeq(t *testing.T)         { fuzzGen(t, &cauchySeqGen{}, grader.GradingMultipleChoice) }
func TestMonotoneSeq(t *testing.T)       { fuzzGen(t, &monotoneSeqGen{}, grader.GradingMultipleChoice) }
func TestEulerSeq(t *testing.T)          { fuzzGen(t, &eulerSeqGen{}, grader.GradingMultipleChoice) }
func TestFunctionSequences(t *testing.T)  { fuzzGen(t, &functionSequencesGen{}, grader.GradingMultipleChoice) }

func TestSeriesConcept(t *testing.T)      { fuzzGen(t, &seriesConceptGen{}, grader.GradingMultipleChoice) }
func TestHarmonicSeries(t *testing.T)     { fuzzGen(t, &harmonicSeriesGen{}, grader.GradingMultipleChoice) }
func TestPositiveTerms(t *testing.T)      { fuzzGen(t, &positiveTermsGen{}, grader.GradingMultipleChoice) }
func TestAlternatingSeries(t *testing.T)  { fuzzGen(t, &alternatingSeriesGen{}, grader.GradingMultipleChoice) }
func TestIntegralTest(t *testing.T)       { fuzzGen(t, &integralTestGen{}, grader.GradingMultipleChoice) }
func TestRootTest(t *testing.T)           { fuzzGen(t, &rootTestGen{}, grader.GradingMultipleChoice) }
func TestPowerSeries(t *testing.T)        { fuzzGen(t, &powerSeriesGen{}, grader.GradingMultipleChoice) }
func TestTaylorSeries(t *testing.T)       { fuzzGen(t, &taylorSeriesGen{}, grader.GradingMultipleChoice) }
func TestSeriesCauchy(t *testing.T)       { fuzzGen(t, &cauchyCriterionSeriesGen{}, grader.GradingMultipleChoice) }
func TestFourierSeries(t *testing.T)      { fuzzGen(t, &fourierSeriesGen{}, grader.GradingMultipleChoice) }
func TestFunctionSeries(t *testing.T)    { fuzzGen(t, &functionSeriesGen{}, grader.GradingMultipleChoice) }
