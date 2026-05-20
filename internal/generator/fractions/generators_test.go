package fractions

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

func TestFracConceptGen(t *testing.T)      { fuzzGen(t, &fracConceptGen{}) }
func TestFracPartsGen(t *testing.T)        { fuzzGen(t, &fracPartsGen{}) }
func TestFracNumberLineGen(t *testing.T)   { fuzzGen(t, &fracNumberLineGen{}) }
func TestFracEquivalentGen(t *testing.T)   { fuzzGen(t, &fracEquivalentGen{}) }
func TestFracSimplifyGen(t *testing.T)     { fuzzGen(t, &fracSimplifyGen{}) }
func TestFracCompareGen(t *testing.T)      { fuzzGen(t, &fracCompareGen{}) }
func TestFracBenchmarkGen(t *testing.T)    { fuzzGen(t, &fracBenchmarkGen{}) }
func TestFracToDecimalGen(t *testing.T)    { fuzzGen(t, &fracToDecimalGen{}) }
func TestFracOpSameDenGen(t *testing.T)    { fuzzGen(t, &fracOpSameDenGen{}) }
func TestFracOpDiffDenGen(t *testing.T)    { fuzzGen(t, &fracOpDiffDenGen{}) }
func TestFracWordGen(t *testing.T)         { fuzzGen(t, &fracWordGen{}) }
func TestFracMultGen(t *testing.T)         { fuzzGen(t, &fracMultGen{}) }
func TestFracDivGen(t *testing.T)          { fuzzGen(t, &fracDivGen{}) }
func TestFracMixedConvertGen(t *testing.T) { fuzzGen(t, &fracMixedConvertGen{}) }
func TestFracMixedOpGen(t *testing.T)      { fuzzGen(t, &fracMixedOpGen{}) }
func TestFracMixedMultGen(t *testing.T)    { fuzzGen(t, &fracMixedMultGen{}) }
