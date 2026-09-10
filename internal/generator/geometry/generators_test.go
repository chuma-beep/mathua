package geometry

import (
	"math/rand"
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

func TestPointsLinesGen(t *testing.T)    { fuzzGen(t, &pointsLinesGen{}) }
func TestAngleTypesGen(t *testing.T)     { fuzzGen(t, &angleTypesGen{}) }
func TestAngleMeasureGen(t *testing.T)   { fuzzGen(t, &angleMeasureGen{}) }
func TestComplementaryGen(t *testing.T)  { fuzzGen(t, &complementaryGen{}) }
func TestVerticalAnglesGen(t *testing.T) { fuzzGen(t, &verticalAnglesGen{}) }
func TestTriangleTypesGen(t *testing.T)  { fuzzGen(t, &triangleTypesGen{}) }
func TestTriangleAnglesGen(t *testing.T) { fuzzGen(t, &triangleAnglesGen{}) }
func TestTriangleAreaGen(t *testing.T)   { fuzzGen(t, &triangleAreaGen{}) }
func TestPythagoreanGen(t *testing.T)    { fuzzGen(t, &pythagoreanGen{}) }
func TestQuadTypesGen(t *testing.T)      { fuzzGen(t, &quadTypesGen{}) }
func TestQuadAreaGen(t *testing.T)       { fuzzGen(t, &quadAreaGen{}) }
func TestQuadPerimGen(t *testing.T)      { fuzzGen(t, &quadPerimGen{}) }
func TestCirclePartsGen(t *testing.T)    { fuzzGen(t, &circlePartsGen{}) }
func TestCircumferenceGen(t *testing.T)  { fuzzGen(t, &circumferenceGen{}) }
func TestCircleAreaGen(t *testing.T)     { fuzzGen(t, &circleAreaGen{}) }
func TestCoordPlotGen(t *testing.T)      { fuzzGen(t, &coordPlotGen{}) }
func TestCoordDistanceGen(t *testing.T)  { fuzzGen(t, &coordDistanceGen{}) }
func TestCoordMidpointGen(t *testing.T)  { fuzzGen(t, &coordMidpointGen{}) }
func TestVolumeGen(t *testing.T)         { fuzzGen(t, &volumeGen{}) }
func TestSurfaceAreaGen(t *testing.T)    { fuzzGen(t, &surfaceAreaGen{}) }
func TestCoordLinesGen(t *testing.T)     { fuzzGen(t, &coordLinesGen{}) }
func TestCoordPolarGen(t *testing.T)     { fuzzGen(t, &coordPolarGen{}) }

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
