package algebra

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

func TestSlope(t *testing.T)          { fuzzGen(t, &slopeGen{}, grader.GradingNumeric) }
func TestSlopeIntercept(t *testing.T) { fuzzGen(t, &slopeInterceptGen{}, grader.GradingExpression) }
func TestLinearGraph(t *testing.T)    { fuzzGen(t, &linearGraphGen{}, grader.GradingNumeric) }
func TestStandardForm(t *testing.T)   { fuzzGen(t, &stdFormGen{}, grader.GradingMultipleChoice) }
func TestParallelPerp(t *testing.T)   { fuzzGen(t, &parallelPerpGen{}, grader.GradingNumeric) }

func TestMultiStepEq(t *testing.T)   { fuzzGen(t, &multiStepEqGen{}, grader.GradingNumeric) }
func TestVarsBothSides(t *testing.T) { fuzzGen(t, &varsBothSidesGen{}, grader.GradingNumeric) }
func TestLiteralEq(t *testing.T)     { fuzzGen(t, &literalEqGen{}, grader.GradingExpression) }
func TestMultiStepIneq(t *testing.T) { fuzzGen(t, &multiStepIneqGen{}, grader.GradingNumeric) }
func TestCompoundIneq(t *testing.T)  { fuzzGen(t, &compoundIneqGen{}, grader.GradingNumeric) }

func TestSysSubstitution(t *testing.T) { fuzzGen(t, &sysSubstitutionGen{}, grader.GradingTuple) }
func TestSysElimination(t *testing.T)  { fuzzGen(t, &sysEliminationGen{}, grader.GradingTuple) }
func TestSysWord(t *testing.T)         { fuzzGen(t, &sysWordGen{}, grader.GradingTuple) }

func TestPolyConcept(t *testing.T)  { fuzzGen(t, &polyConceptGen{}, grader.GradingNumeric) }
func TestPolyAddSub(t *testing.T)   { fuzzGen(t, &polyAddSubGen{}, grader.GradingPolynomial) }
func TestPolyMultMono(t *testing.T) { fuzzGen(t, &polyMultMonoGen{}, grader.GradingPolynomial) }
func TestPolyFoil(t *testing.T)     { fuzzGen(t, &polyFoilGen{}, grader.GradingPolynomial) }
func TestPolySpecial(t *testing.T)  { fuzzGen(t, &polySpecialGen{}, grader.GradingPolynomial) }

func TestFactorGCF(t *testing.T)       { fuzzGen(t, &factorGCFGen{}, grader.GradingExpression) }
func TestFactorTrinomial(t *testing.T) { fuzzGen(t, &factorTrinomialGen{}, grader.GradingExpression) }
func TestFactorDiffSquares(t *testing.T) {
	fuzzGen(t, &factorDiffSquaresGen{}, grader.GradingExpression)
}
func TestFactorACMethod(t *testing.T) { fuzzGen(t, &factorACMethodGen{}, grader.GradingExpression) }

func TestQuadSolveFactor(t *testing.T) { fuzzGen(t, &quadSolveFactorGen{}, grader.GradingTuple) }
func TestQuadCompleteSquare(t *testing.T) {
	fuzzGen(t, &quadCompleteSquareGen{}, grader.GradingNumeric)
}
func TestQuadFormula(t *testing.T)      { fuzzGen(t, &quadFormulaGen{}, grader.GradingTuple) }
func TestQuadDiscriminant(t *testing.T) { fuzzGen(t, &quadDiscriminantGen{}, grader.GradingMultipleChoice) }

func TestFuncConcept(t *testing.T)  { fuzzGen(t, &funcConceptGen{}, grader.GradingNumeric) }
func TestFuncNotation(t *testing.T) { fuzzGen(t, &funcNotationGen{}, grader.GradingNumeric) }
func TestFuncEvaluate(t *testing.T) { fuzzGen(t, &funcEvaluateGen{}, grader.GradingNumeric) }
func TestFuncLinear(t *testing.T)   { fuzzGen(t, &funcLinearGen{}, grader.GradingMultipleChoice) }
func TestFuncQuad(t *testing.T)     { fuzzGen(t, &funcQuadGen{}, grader.GradingNumeric) }

func TestAlgExpConcept(t *testing.T)  { fuzzGen(t, &algExpConceptGen{}, grader.GradingNumeric) }
func TestAlgExpEvaluate(t *testing.T) { fuzzGen(t, &algExpEvaluateGen{}, grader.GradingNumeric) }
func TestLogConcept(t *testing.T)     { fuzzGen(t, &logConceptGen{}, grader.GradingMultipleChoice) }
func TestLogEvaluate(t *testing.T)    { fuzzGen(t, &logEvaluateGen{}, grader.GradingNumeric) }
func TestLogProperties(t *testing.T)  { fuzzGen(t, &logPropertiesGen{}, grader.GradingExpression) }

func TestSeqArith(t *testing.T)    { fuzzGen(t, &seqArithGen{}, grader.GradingNumeric) }
func TestSeqGeom(t *testing.T)     { fuzzGen(t, &seqGeomGen{}, grader.GradingNumeric) }
func TestSeqSumArith(t *testing.T) { fuzzGen(t, &seqSumArithGen{}, grader.GradingNumeric) }
func TestSeqSumGeo(t *testing.T)   { fuzzGen(t, &seqSumGeoGen{}, grader.GradingNumeric) }

func TestIneqTwoVar(t *testing.T)   { fuzzGen(t, &ineqTwoVarGen{}, grader.GradingMultipleChoice) }
func TestConicCircle(t *testing.T)  { fuzzGen(t, &conicCircleGen{}, grader.GradingNumeric) }
func TestConicEllipse(t *testing.T) { fuzzGen(t, &conicEllipseGen{}, grader.GradingNumeric) }
