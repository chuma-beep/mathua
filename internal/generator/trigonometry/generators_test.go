package trigonometry

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

func TestRadians(t *testing.T)        { fuzzGen(t, &radiansGen{}, grader.GradingMultipleChoice) }
func TestUnitCircle(t *testing.T)     { fuzzGen(t, &unitCircleGen{}, grader.GradingMultipleChoice) }
func TestSinCosDef(t *testing.T)      { fuzzGen(t, &sinCosDefGen{}, grader.GradingNumeric) }
func TestTanDef(t *testing.T)         { fuzzGen(t, &tanDefGen{}, grader.GradingNumeric) }
func TestReciprocal(t *testing.T)     { fuzzGen(t, &reciprocalGen{}, grader.GradingNumeric) }
func TestPythagoreanID(t *testing.T)  { fuzzGen(t, &pythagoreanIDGen{}, grader.GradingNumeric) }
func TestSpecialAngles(t *testing.T)  { fuzzGen(t, &specialAnglesGen{}, grader.GradingMultipleChoice) }
func TestReferenceAngle(t *testing.T) { fuzzGen(t, &referenceAngleGen{}, grader.GradingMultipleChoice) }
func TestGraphSin(t *testing.T)       { fuzzGen(t, &graphSinGen{}, grader.GradingMultipleChoice) }
func TestGraphCos(t *testing.T)       { fuzzGen(t, &graphCosGen{}, grader.GradingMultipleChoice) }
func TestPeriod(t *testing.T)         { fuzzGen(t, &periodGen{}, grader.GradingMultipleChoice) }
func TestInverse(t *testing.T)        { fuzzGen(t, &inverseGen{}, grader.GradingMultipleChoice) }
func TestLawSines(t *testing.T)       { fuzzGen(t, &lawSinesGen{}, grader.GradingNumeric) }
func TestLawCosines(t *testing.T)     { fuzzGen(t, &lawCosinesGen{}, grader.GradingNumeric) }
func TestArctan(t *testing.T)         { fuzzGen(t, &arctanGen{}, grader.GradingMultipleChoice) }
func TestRightTriangle(t *testing.T)  { fuzzGen(t, &rightTriangleGen{}, grader.GradingNumeric) }
func TestTrigEqBasic(t *testing.T)    { fuzzGen(t, &trigEqBasicGen{}, grader.GradingMultipleChoice) }
func TestTrigEqHomogeneous(t *testing.T) { fuzzGen(t, &trigEqHomogeneousGen{}, grader.GradingMultipleChoice) }
func TestSinhCosh(t *testing.T)       { fuzzGen(t, &sinhCoshGen{}, grader.GradingMultipleChoice) }
func TestTanhCoth(t *testing.T)       { fuzzGen(t, &tanhCothGen{}, grader.GradingMultipleChoice) }
func TestTrigIdent(t *testing.T)      { fuzzGen(t, &trigIdentGen{}, grader.GradingMultipleChoice) }
func TestTrigIneq(t *testing.T)       { fuzzGen(t, &trigIneqGen{}, grader.GradingMultipleChoice) }
