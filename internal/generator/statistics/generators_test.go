package statistics

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

func TestReadTableGen(t *testing.T)       { fuzzGen(t, &readTableGen{}) }
func TestBarGraphGen(t *testing.T)        { fuzzGen(t, &barGraphGen{}) }
func TestLinePlotGen(t *testing.T)        { fuzzGen(t, &linePlotGen{}) }
func TestMeanGen(t *testing.T)            { fuzzGen(t, &meanGen{}) }
func TestMedianGen(t *testing.T)          { fuzzGen(t, &medianGen{}) }
func TestModeGen(t *testing.T)            { fuzzGen(t, &modeGen{}) }
func TestRangeGen(t *testing.T)           { fuzzGen(t, &rangeGen{}) }
func TestSampleSpaceGen(t *testing.T)     { fuzzGen(t, &sampleSpaceGen{}) }
func TestBasicProbGen(t *testing.T)       { fuzzGen(t, &basicProbGen{}) }
func TestComplementProbGen(t *testing.T)  { fuzzGen(t, &complementProbGen{}) }
func TestCompoundProbGen(t *testing.T)    { fuzzGen(t, &compoundProbGen{}) }
func TestCountingGen(t *testing.T)        { fuzzGen(t, &countingGen{}) }
func TestCovarianceGen(t *testing.T)      { fuzzGen(t, &covarianceGen{}) }
func TestGeometricMeanGen(t *testing.T)   { fuzzGen(t, &geometricMeanGen{}) }
func TestHarmonicMeanGen(t *testing.T)    { fuzzGen(t, &harmonicMeanGen{}) }
func TestRmsGen(t *testing.T)             { fuzzGen(t, &rmsGen{}) }
func TestVarianceGen(t *testing.T)        { fuzzGen(t, &varianceGen{}) }
func TestBernoulliGen(t *testing.T)       { fuzzGen(t, &bernoulliGen{}) }
func TestBetaGen(t *testing.T)            { fuzzGen(t, &betaGen{}) }
func TestBinomialDistGen(t *testing.T)    { fuzzGen(t, &binomialDistGen{}) }
func TestChiSquareGen(t *testing.T)       { fuzzGen(t, &chiSquareGen{}) }
func TestExponentialDistGen(t *testing.T) { fuzzGen(t, &exponentialDistGen{}) }
func TestGammaDistGen(t *testing.T)       { fuzzGen(t, &gammaDistGen{}) }
func TestGeometricDistGen(t *testing.T)   { fuzzGen(t, &geometricDistGen{}) }
func TestHypergeometricGen(t *testing.T)  { fuzzGen(t, &hypergeometricGen{}) }
func TestNormalDistGen(t *testing.T)      { fuzzGen(t, &normalDistGen{}) }
func TestPoissonGen(t *testing.T)         { fuzzGen(t, &poissonGen{}) }
func TestStudentTGen(t *testing.T)        { fuzzGen(t, &studentTGen{}) }
func TestUniformDistGen(t *testing.T)     { fuzzGen(t, &uniformDistGen{}) }
func TestZTableGen(t *testing.T)          { fuzzGen(t, &zTableGen{}) }
func TestConfidenceGen(t *testing.T)      { fuzzGen(t, &confidenceGen{}) }
func TestSamplingGen(t *testing.T)        { fuzzGen(t, &samplingGen{}) }
func TestBayesGen(t *testing.T)           { fuzzGen(t, &bayesGen{}) }
func TestContinuousRVGen(t *testing.T)    { fuzzGen(t, &continuousRVGen{}) }
func TestDiscreteRVGen(t *testing.T)      { fuzzGen(t, &discreteRVGen{}) }
func TestExpectedValueGen(t *testing.T)   { fuzzGen(t, &expectedValueGen{}) }
