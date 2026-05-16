package arithmetic

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

func TestAddSingle(t *testing.T) {
	fuzzGen(t, &addGen{minA: 1, maxA: 9, minB: 1, maxB: 9}, grader.GradingNumeric)
}
func TestAddDouble(t *testing.T) {
	fuzzGen(t, &addGen{minA: 10, maxA: 99, minB: 10, maxB: 99}, grader.GradingNumeric)
}
func TestAddCarry(t *testing.T) {
	fuzzGen(t, &addGen{minA: 10, maxA: 99, minB: 10, maxB: 99, ensureCarry: true}, grader.GradingNumeric)
}
func TestAddTriple(t *testing.T) {
	fuzzGen(t, &addGen{minA: 100, maxA: 999, minB: 100, maxB: 999}, grader.GradingNumeric)
}
func TestAddWord(t *testing.T) { fuzzGen(t, &addWordGen{}, grader.GradingNumeric) }

func TestSubSingle(t *testing.T) {
	fuzzGen(t, &subGen{minA: 5, maxA: 9, minB: 1, maxB: 4}, grader.GradingNumeric)
}
func TestSubDouble(t *testing.T) {
	fuzzGen(t, &subGen{minA: 20, maxA: 99, minB: 10, maxB: 50}, grader.GradingNumeric)
}
func TestSubBorrow(t *testing.T) { fuzzGen(t, &subBorrowGen{}, grader.GradingNumeric) }
func TestSubWord(t *testing.T)   { fuzzGen(t, &subWordGen{}, grader.GradingNumeric) }

func TestPlaceTens(t *testing.T) {
	fuzzGen(t, &placeValueGen{max: 99, label: "tens"}, grader.GradingNumeric)
}
func TestPlaceHundreds(t *testing.T) {
	fuzzGen(t, &placeValueGen{max: 999, label: "hundreds"}, grader.GradingNumeric)
}
func TestPlaceThousands(t *testing.T) {
	fuzzGen(t, &placeValueGen{max: 9999, label: "thousands"}, grader.GradingNumeric)
}

func TestRoundTens(t *testing.T)      { fuzzGen(t, &roundGen{to: 10}, grader.GradingNumeric) }
func TestRoundHundreds(t *testing.T)  { fuzzGen(t, &roundGen{to: 100}, grader.GradingNumeric) }
func TestRoundThousands(t *testing.T) { fuzzGen(t, &roundGen{to: 1000}, grader.GradingNumeric) }

func TestMultConcept(t *testing.T)  { fuzzGen(t, &multConceptGen{}, grader.GradingNumeric) }
func TestMultBy2_5_10(t *testing.T) { fuzzGen(t, &multByGen{}, grader.GradingNumeric) }
func TestMultTables(t *testing.T)   { fuzzGen(t, &multTablesGen{}, grader.GradingNumeric) }
func TestMultDouble(t *testing.T) {
	fuzzGen(t, &multGen{minA: 10, maxA: 99, minB: 2, maxB: 9}, grader.GradingNumeric)
}
func TestMultTriple(t *testing.T) {
	fuzzGen(t, &multGen{minA: 100, maxA: 999, minB: 2, maxB: 9}, grader.GradingNumeric)
}
func TestMultWord(t *testing.T) { fuzzGen(t, &multWordGen{}, grader.GradingNumeric) }

func TestDivConcept(t *testing.T) { fuzzGen(t, &divConceptGen{}, grader.GradingNumeric) }
func TestDivBasic(t *testing.T)   { fuzzGen(t, &divBasicGen{}, grader.GradingNumeric) }
func TestDivRemainder(t *testing.T) {
	gen := &divRemainderGen{}
	for i := 0; i < 100; i++ {
		p := gen.Generate(rand.Float64())
		if p.Question == "" || p.Answer == "" || p.Explanation == "" {
			t.Errorf("empty field: q=%q a=%q e=%q", p.Question, p.Answer, p.Explanation)
		}
	}
}
func TestDivLong(t *testing.T) { fuzzGen(t, &divLongGen{}, grader.GradingNumeric) }
func TestDivWord(t *testing.T) { fuzzGen(t, &divWordGen{}, grader.GradingNumeric) }

func TestFactorFind(t *testing.T) { fuzzGen(t, &factorFindGen{}, grader.GradingOrdering) }
func TestPrime(t *testing.T)      { fuzzGen(t, &primeGen{}, grader.GradingMultipleChoice) }
func TestPrimeFact(t *testing.T)  { fuzzGen(t, &primeFactGen{}, grader.GradingOrdering) }
func TestGCF(t *testing.T)        { fuzzGen(t, &gcfGen{}, grader.GradingNumeric) }
func TestLCM(t *testing.T)        { fuzzGen(t, &lcmGen{}, grader.GradingNumeric) }
func TestComposite(t *testing.T)  { fuzzGen(t, &compositeGen{}, grader.GradingMultipleChoice) }

func TestExpConcept(t *testing.T)      { fuzzGen(t, &expConceptGen{}, grader.GradingNumeric) }
func TestExpEval(t *testing.T)         { fuzzGen(t, &expEvalGen{}, grader.GradingNumeric) }
func TestExpProductRule(t *testing.T)  { fuzzGen(t, &expProductRuleGen{}, grader.GradingExpression) }
func TestExpQuotientRule(t *testing.T) { fuzzGen(t, &expQuotientRuleGen{}, grader.GradingExpression) }
func TestExpPowerRule(t *testing.T)    { fuzzGen(t, &expPowerRuleGen{}, grader.GradingExpression) }

func TestSqrtPerfect(t *testing.T)  { fuzzGen(t, &sqrtPerfectGen{}, grader.GradingNumeric) }
func TestSqrtSimplify(t *testing.T) { fuzzGen(t, &sqrtSimplifyGen{}, grader.GradingPolynomial) }

func TestNegNumberLine(t *testing.T) { fuzzGen(t, &negNumberLineGen{}, grader.GradingNumeric) }
func TestNegAddSub(t *testing.T)     { fuzzGen(t, &negAddSubGen{}, grader.GradingNumeric) }
func TestNegMultDiv(t *testing.T)    { fuzzGen(t, &negMultDivGen{}, grader.GradingNumeric) }

func TestOrderOpsBasic(t *testing.T) {
	fuzzGen(t, &orderOpsGen{parens: false, exponents: false}, grader.GradingNumeric)
}
func TestOrderOpsFull(t *testing.T) {
	fuzzGen(t, &orderOpsGen{parens: true, exponents: false}, grader.GradingNumeric)
}
func TestOrderOpsNested(t *testing.T) {
	fuzzGen(t, &orderOpsGen{parens: true, exponents: true}, grader.GradingNumeric)
}

func TestDecIntro(t *testing.T) { fuzzGen(t, &decIntroGen{}, grader.GradingNumeric) }
