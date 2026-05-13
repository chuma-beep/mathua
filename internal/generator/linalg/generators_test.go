package linalg

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

func TestVectorConcept(t *testing.T)  { fuzzGen(t, &vectorConceptGen{}, grader.GradingMultipleChoice) }
func TestVectorAdd(t *testing.T)      { fuzzGen(t, &vectorAddGen{}, grader.GradingMultipleChoice) }
func TestVectorDot(t *testing.T)      { fuzzGen(t, &vectorDotGen{}, grader.GradingNumeric) }
func TestMatrixConcept(t *testing.T)  { fuzzGen(t, &matrixConceptGen{}, grader.GradingMultipleChoice) }
func TestMatrixAdd(t *testing.T)      { fuzzGen(t, &matrixAddGen{}, grader.GradingMultipleChoice) }
func TestMatrixMult(t *testing.T)     { fuzzGen(t, &matrixMultGen{}, grader.GradingMultipleChoice) }
func TestMatrixIdentity(t *testing.T) { fuzzGen(t, &matrixIdentityGen{}, grader.GradingMultipleChoice) }
func TestDet2x2(t *testing.T)         { fuzzGen(t, &det2x2Gen{}, grader.GradingNumeric) }
func TestDet3x3(t *testing.T)         { fuzzGen(t, &det3x3Gen{}, grader.GradingNumeric) }
func TestSystemsMatrix(t *testing.T)  { fuzzGen(t, &systemsMatrixGen{}, grader.GradingMultipleChoice) }
func TestCramer(t *testing.T)         { fuzzGen(t, &cramerGen{}, grader.GradingNumeric) }
func TestEigenConcept(t *testing.T)   { fuzzGen(t, &eigenConceptGen{}, grader.GradingMultipleChoice) }
func TestEigenCompute(t *testing.T)   { fuzzGen(t, &eigenComputeGen{}, grader.GradingMultipleChoice) }
func TestTransformations(t *testing.T) {
	fuzzGen(t, &transformationsGen{}, grader.GradingMultipleChoice)
}
func TestSpan(t *testing.T)  { fuzzGen(t, &spanGen{}, grader.GradingMultipleChoice) }
func TestBasis(t *testing.T) { fuzzGen(t, &basisGen{}, grader.GradingMultipleChoice) }
