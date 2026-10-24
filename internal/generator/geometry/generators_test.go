package geometry

import (
	"math/rand"
	"regexp"
	"strconv"
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

func TestPointsLinesGen(t *testing.T)     { fuzzGen(t, &pointsLinesGen{}) }
func TestAngleTypesGen(t *testing.T)      { fuzzGen(t, &angleTypesGen{}) }
func TestAngleMeasureGen(t *testing.T)    { fuzzGen(t, &angleMeasureGen{}) }
func TestComplementaryGen(t *testing.T)   { fuzzGen(t, &complementaryGen{}) }
func TestVerticalAnglesGen(t *testing.T)  { fuzzGen(t, &verticalAnglesGen{}) }
func TestTriangleTypesGen(t *testing.T)   { fuzzGen(t, &triangleTypesGen{}) }
func TestTriangleAnglesGen(t *testing.T)  { fuzzGen(t, &triangleAnglesGen{}) }
func TestTriangleAreaGen(t *testing.T)    { fuzzGen(t, &triangleAreaGen{}) }
func TestPythagoreanGen(t *testing.T)     { fuzzGen(t, &pythagoreanGen{}) }
func TestPythagoreanWordGen(t *testing.T) { fuzzGen(t, &pythagoreanWordGen{}) }
func TestQuadTypesGen(t *testing.T)       { fuzzGen(t, &quadTypesGen{}) }
func TestQuadAreaGen(t *testing.T)        { fuzzGen(t, &quadAreaGen{}) }
func TestQuadPerimGen(t *testing.T)       { fuzzGen(t, &quadPerimGen{}) }
func TestCirclePartsGen(t *testing.T)     { fuzzGen(t, &circlePartsGen{}) }
func TestCircumferenceGen(t *testing.T)   { fuzzGen(t, &circumferenceGen{}) }
func TestCircleAreaGen(t *testing.T)      { fuzzGen(t, &circleAreaGen{}) }
func TestCoordPlotGen(t *testing.T)       { fuzzGen(t, &coordPlotGen{}) }
func TestCoordDistanceGen(t *testing.T)   { fuzzGen(t, &coordDistanceGen{}) }
func TestCoordMidpointGen(t *testing.T)   { fuzzGen(t, &coordMidpointGen{}) }
func TestVolumeGen(t *testing.T)          { fuzzGen(t, &volumeGen{}) }
func TestSurfaceAreaGen(t *testing.T)     { fuzzGen(t, &surfaceAreaGen{}) }
func TestCoordLinesGen(t *testing.T)      { fuzzGen(t, &coordLinesGen{}) }
func TestCoordPolarGen(t *testing.T)      { fuzzGen(t, &coordPolarGen{}) }

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

func TestCircleAnswersCarryTwoDecimalPlaces(t *testing.T) {
	cases := []struct {
		name string
		gen  generator.Generator
		rate float64
	}{
		{"circumference", &circumferenceGen{}, 2},
		{"area", &circleAreaGen{}, 1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for i := 0; i < 100; i++ {
				p := tc.gen.Generate(generator.GeneratorContext{Difficulty: rand.Float64(), Seed: rand.Int63()})
				if !strings.Contains(p.Question, "two decimal places") {
					t.Fatalf("question does not state the precision: %q", p.Question)
				}
				if !strings.Contains(p.Answer, ".") || len(p.Answer)-strings.Index(p.Answer, ".") != 3 {
					t.Fatalf("answer %q does not carry two decimal places", p.Answer)
				}
				got, err := strconv.ParseFloat(p.Answer, 64)
				if err != nil {
					t.Fatalf("answer %q does not parse: %v", p.Answer, err)
				}
				r, err := strconv.Atoi(p.Facts["r"])
				if err != nil {
					t.Fatalf("fact r = %q does not parse", p.Facts["r"])
				}
				want := 3.14 * float64(r*r)
				if tc.rate == 2 {
					want = 2 * 3.14 * float64(r)
				}
				if diff := got - want; diff > 0.005 || diff < -0.005 {
					t.Fatalf("answer %q does not match the exact value %.4f", p.Answer, want)
				}
			}
		})
	}
}

func TestVerticalAngles_ReadOffAvoidance(t *testing.T) {
	gen := &verticalAnglesGen{}
	for i := 0; i < 500; i++ {
		p := gen.Generate(generator.GeneratorContext{
			Difficulty: rand.Float64(),
			Seed:       int64(i*37 + 11),
		})
		ans := strings.TrimSpace(p.Answer)
		// Check that the numerical answer does not appear verbatim in the stem
		pattern := `(^|\D)` + ans + `(\D|$)`
		if matched, _ := regexp.MatchString(pattern, p.Question); matched {
			t.Fatalf("draw %d: vertical angles question %q prints answer %q in stem", i, p.Question, ans)
		}
	}
}
