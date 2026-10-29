package linalg

import (
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/grader"
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

func TestVectorConceptGen(t *testing.T) { fuzzGen(t, &vectorConceptGen{}) }

func TestVectorConceptGenFollowsTheKnowledgePoint(t *testing.T) {
	g := &vectorConceptGen{}

	for i := 0; i < 40; i++ {
		p := g.Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: 1})
		if !strings.Contains(p.Question, "magnitude of vector") {
			t.Fatalf("KP1 question is not a magnitude prompt: %q", p.Question)
		}
	}

	for i := 0; i < 40; i++ {
		p := g.Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: 2})
		if !strings.Contains(p.Question, "points in the direction of") {
			t.Fatalf("KP2 question does not ask for a component from a direction: %q", p.Question)
		}
		x, _ := strconv.Atoi(p.Facts["x"])
		scale, _ := strconv.Atoi(p.Facts["scale"])
		if want := strconv.Itoa(x * scale); p.Answer != want {
			t.Fatalf("KP2 answer %q does not match x*scale %q: %q", p.Answer, want, p.Question)
		}
	}

	for i := 0; i < 40; i++ {
		p := g.Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: 3})
		if !strings.Contains(p.Question, `\mathbf{u}`) || !strings.Contains(p.Question, `\mathbf{v}`) {
			t.Fatalf("KP3 question does not ask about two named vectors: %q", p.Question)
		}
		u1, _ := strconv.Atoi(p.Facts["u1"])
		v1, _ := strconv.Atoi(p.Facts["v1"])
		if want := strconv.Itoa(u1 + v1); p.Answer != want {
			t.Fatalf("KP3 answer %q does not match u1+v1 %q: %q", p.Answer, want, p.Question)
		}
	}
}
func TestVectorAddGen(t *testing.T)        { fuzzGen(t, &vectorAddGen{}) }
func TestVectorDotGen(t *testing.T)        { fuzzGen(t, &vectorDotGen{}) }
func TestMatrixConceptGen(t *testing.T)    { fuzzGen(t, &matrixConceptGen{}) }
func TestMatrixAddGen(t *testing.T)        { fuzzGen(t, &matrixAddGen{}) }
func TestMatrixMultGen(t *testing.T)       { fuzzGen(t, &matrixMultGen{}) }
func TestMatrixIdentityGen(t *testing.T)   { fuzzGen(t, &matrixIdentityGen{}) }
func TestDet2x2Gen(t *testing.T)           { fuzzGen(t, &det2x2Gen{}) }
func TestDet3x3Gen(t *testing.T)           { fuzzGen(t, &det3x3Gen{}) }
func TestSystemsMatrixGen(t *testing.T)    { fuzzGen(t, &systemsMatrixGen{}) }
func TestCramerGen(t *testing.T)           { fuzzGen(t, &cramerGen{}) }
func TestEigenConceptGen(t *testing.T)     { fuzzGen(t, &eigenConceptGen{}) }
func TestEigenComputeGen(t *testing.T)     { fuzzGen(t, &eigenComputeGen{}) }
func TestTransformationsGen(t *testing.T)  { fuzzGen(t, &transformationsGen{}) }
func TestSpanGen(t *testing.T)             { fuzzGen(t, &spanGen{}) }
func TestBasisGen(t *testing.T)            { fuzzGen(t, &basisGen{}) }
func TestDiagonalizationGen(t *testing.T)  { fuzzGen(t, &diagonalizationGen{}) }
func TestRankGen(t *testing.T)             { fuzzGen(t, &rankGen{}) }
func TestCosineSimilarityGen(t *testing.T) { fuzzGen(t, &cosineSimilarityGen{}) }
func TestParametricGen(t *testing.T)       { fuzzGen(t, &parametricGen{}) }

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

func TestParametricEquationsUseSignedTerms(t *testing.T) {
	g := &parametricGen{}
	r := grader.NewRouter()
	sawNegative := false
	for i := 0; i < 300; i++ {
		p := g.Generate(generator.GeneratorContext{Difficulty: rand.Float64()})
		if strings.Contains(p.Question, "+ -") || strings.Contains(p.Answer, "+ -") || strings.Contains(p.Explanation, "+ -") {
			t.Fatalf("parametric problem writes a plus followed by a negative term: q=%q a=%q e=%q", p.Question, p.Answer, p.Explanation)
		}
		if strings.Contains(p.Answer, "-") {
			sawNegative = true
		}
		if res := r.Grade(grader.GradingTuple, p.Answer, p.Answer); !res.Correct {
			t.Fatalf("parametric answer %q does not grade as its own tuple: %v", p.Answer, res.Feedback)
		}
	}
	if !sawNegative {
		t.Fatal("300 draws produced no negative direction or starting value; the sign test did not exercise the defect")
	}
	if hint := g.AnswerFormat().Hint; !strings.Contains(hint, "x=") || !strings.Contains(hint, "y=") {
		t.Fatalf("answer hint does not show both coordinate equations: %q", hint)
	}
}

func TestRandomMatrixCoefficientsAreNonsingular(t *testing.T) {
	for i := 0; i < 500; i++ {
		scale := 1 + rand.Intn(5)
		a, b, c, d := randomMatrixCoefficients(scale)
		if a*d-b*c == 0 {
			t.Fatalf("singular matrix from scale %d: [[%d,%d],[%d,%d]]", scale, a, b, c, d)
		}
	}
}
