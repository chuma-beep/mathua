package fractions

import (
	"math/rand"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/grader"
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
	}
}

func TestFracConcept(t *testing.T)      { fuzzGen(t, &fracConceptGen{}, grader.GradingNumeric) }
func TestFracParts(t *testing.T)        { fuzzGen(t, &fracPartsGen{}, grader.GradingNumeric) }
func TestFracNumberLine(t *testing.T)   { fuzzGen(t, &fracNumberLineGen{}, grader.GradingNumeric) }
func TestFracEquivalent(t *testing.T)   { fuzzGen(t, &fracEquivalentGen{}, grader.GradingNumeric) }
func TestFracSimplify(t *testing.T)     { fuzzGen(t, &fracSimplifyGen{}, grader.GradingNumeric) }
func TestFracCompare(t *testing.T)      { fuzzGen(t, &fracCompareGen{}, grader.GradingComparison) }
func TestFracBenchmark(t *testing.T)    { fuzzGen(t, &fracBenchmarkGen{}, grader.GradingMultipleChoice) }
func TestFracToDecimal(t *testing.T)    { fuzzGen(t, &fracToDecimalGen{}, grader.GradingNumeric) }
func TestFracAddSame(t *testing.T)      { fuzzGen(t, &fracOpSameDenGen{op: "+"}, grader.GradingNumeric) }
func TestFracSubSame(t *testing.T)      { fuzzGen(t, &fracOpSameDenGen{op: "-"}, grader.GradingNumeric) }
func TestFracAddDiff(t *testing.T)      { fuzzGen(t, &fracOpDiffDenGen{op: "+"}, grader.GradingNumeric) }
func TestFracSubDiff(t *testing.T)      { fuzzGen(t, &fracOpDiffDenGen{op: "-"}, grader.GradingNumeric) }
func TestFracAddWord(t *testing.T)      { fuzzGen(t, &fracWordGen{op: "+"}, grader.GradingNumeric) }
func TestFracMult(t *testing.T)         { fuzzGen(t, &fracMultGen{wholeMul: false}, grader.GradingNumeric) }
func TestFracMultWhole(t *testing.T)    { fuzzGen(t, &fracMultGen{wholeMul: true}, grader.GradingNumeric) }
func TestFracDiv(t *testing.T)          { fuzzGen(t, &fracDivGen{wholeDiv: false}, grader.GradingNumeric) }
func TestFracDivWhole(t *testing.T)     { fuzzGen(t, &fracDivGen{wholeDiv: true}, grader.GradingNumeric) }
func TestFracMixedConvert(t *testing.T) {
	gen := &fracMixedConvertGen{}
	for i := 0; i < 100; i++ {
		p := gen.Generate(rand.Float64())
		if p.Question == "" || p.Answer == "" || p.Explanation == "" {
			t.Errorf("empty field")
		}
	}
}
func TestFracMixedAdd(t *testing.T)     { fuzzGen(t, &fracMixedOpGen{op: "+"}, grader.GradingNumeric) }
func TestFracMixedSub(t *testing.T)     { fuzzGen(t, &fracMixedOpGen{op: "-"}, grader.GradingNumeric) }
func TestFracMixedMult(t *testing.T)    { fuzzGen(t, &fracMixedMultGen{}, grader.GradingNumeric) }
