package geometry

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

func fuzzGenNoSelfGrade(t *testing.T, gen generator.Generator) {
	t.Helper()
	for i := 0; i < 100; i++ {
		p := gen.Generate(rand.Float64())
		if p.Question == "" || p.Answer == "" || p.Explanation == "" {
			t.Errorf("empty field: q=%q a=%q e=%q", p.Question, p.Answer, p.Explanation)
		}
	}
}

func TestPointsLines(t *testing.T)    { fuzzGen(t, &pointsLinesGen{}, grader.GradingMultipleChoice) }
func TestAngleTypes(t *testing.T)     { fuzzGen(t, &angleTypesGen{}, grader.GradingMultipleChoice) }
func TestAngleMeasure(t *testing.T)   { fuzzGen(t, &angleMeasureGen{}, grader.GradingMultipleChoice) }
func TestComplementary(t *testing.T)  { fuzzGen(t, &complementaryGen{}, grader.GradingNumeric) }
func TestVerticalAngles(t *testing.T) { fuzzGen(t, &verticalAnglesGen{}, grader.GradingNumeric) }
func TestTriangleTypes(t *testing.T)  { fuzzGen(t, &triangleTypesGen{}, grader.GradingMultipleChoice) }
func TestTriangleAngles(t *testing.T) { fuzzGen(t, &triangleAnglesGen{}, grader.GradingNumeric) }
func TestTriangleArea(t *testing.T)   { fuzzGen(t, &triangleAreaGen{}, grader.GradingNumeric) }
func TestPythagorean(t *testing.T)    { fuzzGen(t, &pythagoreanGen{}, grader.GradingNumeric) }
func TestQuadTypes(t *testing.T)      { fuzzGen(t, &quadTypesGen{}, grader.GradingMultipleChoice) }
func TestQuadArea(t *testing.T)       { fuzzGen(t, &quadAreaGen{}, grader.GradingNumeric) }
func TestQuadPerim(t *testing.T)      { fuzzGen(t, &quadPerimGen{}, grader.GradingNumeric) }
func TestCircleParts(t *testing.T)    { fuzzGen(t, &circlePartsGen{}, grader.GradingMultipleChoice) }
func TestCircumference(t *testing.T)  { fuzzGen(t, &circumferenceGen{}, grader.GradingNumeric) }
func TestCircleArea(t *testing.T)     { fuzzGen(t, &circleAreaGen{}, grader.GradingNumeric) }
func TestCoordPlot(t *testing.T)      { fuzzGen(t, &coordPlotGen{}, grader.GradingMultipleChoice) }
func TestCoordDistance(t *testing.T)  { fuzzGen(t, &coordDistanceGen{}, grader.GradingNumeric) }
func TestCoordMidpoint(t *testing.T)  { fuzzGen(t, &coordMidpointGen{}, grader.GradingTuple) }
func TestVolume(t *testing.T)         { fuzzGen(t, &volumeGen{}, grader.GradingNumeric) }
func TestSurfaceArea(t *testing.T)    { fuzzGen(t, &surfaceAreaGen{}, grader.GradingNumeric) }
