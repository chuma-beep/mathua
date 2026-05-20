package calculus

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

func TestLimitConceptGen(t *testing.T)       { fuzzGen(t, &limitConceptGen{}) }
func TestLimitNumericGen(t *testing.T)        { fuzzGen(t, &limitNumericGen{}) }
func TestLimitPropertiesGen(t *testing.T)     { fuzzGen(t, &limitPropertiesGen{}) }
func TestLimitInfinityGen(t *testing.T)       { fuzzGen(t, &limitInfinityGen{}) }
func TestLimitContinuityGen(t *testing.T)     { fuzzGen(t, &limitContinuityGen{}) }
func TestDerivConceptGen(t *testing.T)        { fuzzGen(t, &derivConceptGen{}) }
func TestDerivPowerRuleGen(t *testing.T)      { fuzzGen(t, &derivPowerRuleGen{}) }
func TestDerivSumRuleGen(t *testing.T)        { fuzzGen(t, &derivSumRuleGen{}) }
func TestDerivProductRuleGen(t *testing.T)    { fuzzGen(t, &derivProductRuleGen{}) }
func TestDerivQuotientRuleGen(t *testing.T)   { fuzzGen(t, &derivQuotientRuleGen{}) }
func TestDerivChainRuleGen(t *testing.T)      { fuzzGen(t, &derivChainRuleGen{}) }
func TestDerivTrigGen(t *testing.T)           { fuzzGen(t, &derivTrigGen{}) }
func TestDerivExpLogGen(t *testing.T)         { fuzzGen(t, &derivExpLogGen{}) }
func TestDerivApplicationsGen(t *testing.T)   { fuzzGen(t, &derivApplicationsGen{}) }
func TestDerivOptimizationGen(t *testing.T)   { fuzzGen(t, &derivOptimizationGen{}) }
func TestIntegralIndefiniteGen(t *testing.T)  { fuzzGen(t, &integralIndefiniteGen{}) }
func TestIntegralPowerRuleGen(t *testing.T)   { fuzzGen(t, &integralPowerRuleGen{}) }
func TestIntegralSubstitutionGen(t *testing.T) { fuzzGen(t, &integralSubstitutionGen{}) }
func TestIntegralDefiniteGen(t *testing.T)    { fuzzGen(t, &integralDefiniteGen{}) }
func TestIntegralFTCGen(t *testing.T)         { fuzzGen(t, &integralFTCGen{}) }
func TestIntegralAreaBetweenGen(t *testing.T) { fuzzGen(t, &integralAreaBetweenGen{}) }
func TestIntegralVolumeGen(t *testing.T)      { fuzzGen(t, &integralVolumeGen{}) }
func TestDerivImplicitGen(t *testing.T)       { fuzzGen(t, &derivImplicitGen{}) }
func TestDerivRelatedRatesGen(t *testing.T)   { fuzzGen(t, &derivRelatedRatesGen{}) }
func TestIntegralPartsGen(t *testing.T)       { fuzzGen(t, &integralPartsGen{}) }
func TestIntegralPartialFractionsGen(t *testing.T) { fuzzGen(t, &integralPartialFractionsGen{}) }
func TestRolleGen(t *testing.T)               { fuzzGen(t, &rolleGen{}) }
func TestMvtGen(t *testing.T)                 { fuzzGen(t, &mvtGen{}) }
func TestCauchyMVTGen(t *testing.T)           { fuzzGen(t, &cauchyMVTGen{}) }
func TestFermatGen(t *testing.T)              { fuzzGen(t, &fermatGen{}) }
func TestDifferenceQuotientGen(t *testing.T)  { fuzzGen(t, &differenceQuotientGen{}) }
func TestDifferentialGen(t *testing.T)        { fuzzGen(t, &differentialGen{}) }
func TestConvexityGen(t *testing.T)           { fuzzGen(t, &convexityGen{}) }
func TestNonDiffGen(t *testing.T)             { fuzzGen(t, &nonDiffGen{}) }
func TestPartialDerivGen(t *testing.T)        { fuzzGen(t, &partialDerivGen{}) }
func TestArcLengthGen(t *testing.T)           { fuzzGen(t, &arcLengthGen{}) }
func TestExpIntegralGen(t *testing.T)         { fuzzGen(t, &expIntegralGen{}) }
func TestImproperIntegralGen(t *testing.T)    { fuzzGen(t, &improperIntegralGen{}) }
func TestNumericalIntegralGen(t *testing.T)   { fuzzGen(t, &numericalIntegralGen{}) }
func TestRiemannCriteriaGen(t *testing.T)     { fuzzGen(t, &riemannCriteriaGen{}) }
func TestTrigIntegralsGen(t *testing.T)       { fuzzGen(t, &trigIntegralsGen{}) }
func TestTrigSubstitutionGen(t *testing.T)    { fuzzGen(t, &trigSubstitutionGen{}) }
func TestWeierstrassSubGen(t *testing.T)      { fuzzGen(t, &weierstrassSubGen{}) }
func TestLimitAlgebraGen(t *testing.T)        { fuzzGen(t, &limitAlgebraGen{}) }
func TestAsymptotesGen(t *testing.T)          { fuzzGen(t, &asymptotesGen{}) }
func TestBigOGen(t *testing.T)                { fuzzGen(t, &bigOGen{}) }
func TestDiscontinuityGen(t *testing.T)       { fuzzGen(t, &discontinuityGen{}) }
func TestIndeterminateGen(t *testing.T)       { fuzzGen(t, &indeterminateGen{}) }
func TestLhopitalGen(t *testing.T)            { fuzzGen(t, &lhopitalGen{}) }
func TestLittleOGen(t *testing.T)             { fuzzGen(t, &littleOGen{}) }
func TestSqueezeGen(t *testing.T)             { fuzzGen(t, &squeezeGen{}) }
func TestSupremumGen(t *testing.T)            { fuzzGen(t, &supremumGen{}) }
func TestUniformContinuityGen(t *testing.T)   { fuzzGen(t, &uniformContinuityGen{}) }
func TestWeierstrassLimitGen(t *testing.T)    { fuzzGen(t, &weierstrassLimitGen{}) }
func TestSeqConceptGen(t *testing.T)          { fuzzGen(t, &seqConceptGen{}) }
func TestSeqConvergenceGen(t *testing.T)      { fuzzGen(t, &seqConvergenceGen{}) }
func TestCauchySeqGen(t *testing.T)           { fuzzGen(t, &cauchySeqGen{}) }
func TestMonotoneSeqGen(t *testing.T)         { fuzzGen(t, &monotoneSeqGen{}) }
func TestEulerSeqGen(t *testing.T)            { fuzzGen(t, &eulerSeqGen{}) }
func TestFunctionSequencesGen(t *testing.T)   { fuzzGen(t, &functionSequencesGen{}) }
func TestSeriesConceptGen(t *testing.T)       { fuzzGen(t, &seriesConceptGen{}) }
func TestHarmonicSeriesGen(t *testing.T)      { fuzzGen(t, &harmonicSeriesGen{}) }
func TestPositiveTermsGen(t *testing.T)       { fuzzGen(t, &positiveTermsGen{}) }
func TestAlternatingSeriesGen(t *testing.T)   { fuzzGen(t, &alternatingSeriesGen{}) }
func TestIntegralTestGen(t *testing.T)        { fuzzGen(t, &integralTestGen{}) }
func TestRootTestGen(t *testing.T)            { fuzzGen(t, &rootTestGen{}) }
func TestPowerSeriesGen(t *testing.T)         { fuzzGen(t, &powerSeriesGen{}) }
func TestTaylorSeriesGen(t *testing.T)        { fuzzGen(t, &taylorSeriesGen{}) }
func TestCauchyCriterionSeriesGen(t *testing.T) { fuzzGen(t, &cauchyCriterionSeriesGen{}) }
func TestFourierSeriesGen(t *testing.T)       { fuzzGen(t, &fourierSeriesGen{}) }
func TestFunctionSeriesGen(t *testing.T)      { fuzzGen(t, &functionSeriesGen{}) }
