package prealgebra

import (
	"math/rand"
	"strings"
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

func TestDecCompareGen(t *testing.T)      { fuzzGen(t, &decCompareGen{}) }
func TestDecAddSubGen(t *testing.T)       { fuzzGen(t, &decAddSubGen{}) }
func TestDecMultGen(t *testing.T)         { fuzzGen(t, &decMultGen{}) }
func TestDecDivGen(t *testing.T)          { fuzzGen(t, &decDivGen{}) }
func TestDecFromFracGen(t *testing.T)     { fuzzGen(t, &decFromFracGen{}) }
func TestDecToFracGen(t *testing.T)       { fuzzGen(t, &decToFracGen{}) }
func TestDecRoundGen(t *testing.T)        { fuzzGen(t, &decRoundGen{}) }
func TestPctConceptGen(t *testing.T)      { fuzzGen(t, &pctConceptGen{}) }
func TestPctToDecGen(t *testing.T)        { fuzzGen(t, &pctToDecGen{}) }
func TestPctFromDecGen(t *testing.T)      { fuzzGen(t, &pctFromDecGen{}) }
func TestPctOfNumberGen(t *testing.T)     { fuzzGen(t, &pctOfNumberGen{}) }
func TestPctFindRateGen(t *testing.T)     { fuzzGen(t, &pctFindRateGen{}) }
func TestPctIncreaseGen(t *testing.T)     { fuzzGen(t, &pctIncreaseGen{}) }
func TestPctDiscountGen(t *testing.T)     { fuzzGen(t, &pctDiscountGen{}) }
func TestPctTaxTipGen(t *testing.T)       { fuzzGen(t, &pctTaxTipGen{}) }
func TestAbsValueGen(t *testing.T)        { fuzzGen(t, &absValueGen{}) }
func TestNegOrderOpsGen(t *testing.T)     { fuzzGen(t, &negOrderOpsGen{}) }
func TestRatioConceptGen(t *testing.T)    { fuzzGen(t, &ratioConceptGen{}) }
func TestRatioSimplifyGen(t *testing.T)   { fuzzGen(t, &ratioSimplifyGen{}) }
func TestRatioProportionGen(t *testing.T) { fuzzGen(t, &ratioProportionGen{}) }
func TestRatioRateGen(t *testing.T)       { fuzzGen(t, &ratioRateGen{}) }
func TestRatioScaleGen(t *testing.T)      { fuzzGen(t, &ratioScaleGen{}) }
func TestExpNegGen(t *testing.T)          { fuzzGen(t, &expNegGen{}) }
func TestExpZeroGen(t *testing.T)         { fuzzGen(t, &expZeroGen{}) }
func TestSciNotationGen(t *testing.T)     { fuzzGen(t, &sciNotationGen{}) }
func TestSciNotationOpsGen(t *testing.T)  { fuzzGen(t, &sciNotationOpsGen{}) }
func TestVarConceptGen(t *testing.T)      { fuzzGen(t, &varConceptGen{}) }
func TestExprEvalGen(t *testing.T)        { fuzzGen(t, &exprEvalGen{}) }
func TestLikeTermsGen(t *testing.T)       { fuzzGen(t, &likeTermsGen{}) }
func TestDistributeGen(t *testing.T)      { fuzzGen(t, &distributeGen{}) }
func TestEqOneStepAddGen(t *testing.T)    { fuzzGen(t, &eqOneStepAddGen{}) }
func TestEqOneStepMultGen(t *testing.T)   { fuzzGen(t, &eqOneStepMultGen{}) }
func TestEqTwoStepGen(t *testing.T)       { fuzzGen(t, &eqTwoStepGen{}) }
func TestEqWordGen(t *testing.T)          { fuzzGen(t, &eqWordGen{}) }
func TestIneqOneStepGen(t *testing.T)     { fuzzGen(t, &ineqOneStepGen{}) }
func TestIneqTwoStepGen(t *testing.T)     { fuzzGen(t, &ineqTwoStepGen{}) }
func TestRealConceptGen(t *testing.T)     { fuzzGen(t, &realConceptGen{}) }
func TestRealPropertiesGen(t *testing.T)  { fuzzGen(t, &realPropertiesGen{}) }
func TestTypesGen(t *testing.T)           { fuzzGen(t, &typesGen{}) }

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

// TestSciNotationGenBothDirections is the template quality bar for practice
// breadth: the pool must serve scientific->standard (positive and negative
// exponents) and standard->scientific, and every served answer must grade
// true against itself through the generator's own grader.
func TestSciNotationGenBothDirections(t *testing.T) {
	rand.Seed(11)
	gen := &sciNotationGen{}
	seenToStd, seenToSci := false, false
	answers := map[string]bool{}
	for i := 0; i < 150; i++ {
		p := gen.Generate(generator.GeneratorContext{Difficulty: 0.3})
		if p.Question == "" || p.Answer == "" || p.Explanation == "" {
			t.Fatalf("empty field at sample %d: q=%q a=%q e=%q", i, p.Question, p.Answer, p.Explanation)
		}
		answers[p.Answer] = true
		if strings.Contains(p.Answer, "x10^") {
			seenToSci = true
		} else {
			seenToStd = true
		}
		if res := gen.Grade(p.Answer, p.Answer); !res.Correct {
			t.Errorf("answer does not grade against itself: q=%q a=%q fb=%q", p.Question, p.Answer, res.Feedback)
		}
		if res := gen.Grade(p.Answer, ""); res.Correct {
			t.Errorf("empty answer accepted for q=%q", p.Question)
		}
	}
	if len(answers) < 2 {
		t.Errorf("pool non-discriminating: %d distinct answers over 150 samples", len(answers))
	}
	if !seenToStd || !seenToSci {
		t.Errorf("pool covers one direction only: toStd=%v toSci=%v", seenToStd, seenToSci)
	}
	// Spot checks with hand-computed values.
	for _, tc := range []struct {
		expected, answer string
		want             bool
	}{
		{"36000", "36000", true},
		{"0.0057", "0.0057", true},
		{"3.7x10^4", "3.7x10^4", true},
		{"3.7x10^4", "37000", false},
		{"5.7x10^-3", "5.7x10^-3", true},
		{"36000", "35", false},
	} {
		if res := gen.Grade(tc.expected, tc.answer); res.Correct != tc.want {
			t.Errorf("Grade(%q, %q) = %v, want %v (%q)", tc.expected, tc.answer, res.Correct, tc.want, res.Feedback)
		}
	}
}
