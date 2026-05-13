package complex

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

func TestConcept(t *testing.T)   { fuzzGen(t, &conceptGen{}, grader.GradingMultipleChoice) }
func TestAddSub(t *testing.T)    { fuzzGen(t, &addSubGen{}, grader.GradingMultipleChoice) }
func TestMult(t *testing.T)      { fuzzGen(t, &multGen{}, grader.GradingMultipleChoice) }
func TestConjugate(t *testing.T) { fuzzGen(t, &conjugateGen{}, grader.GradingMultipleChoice) }
func TestDivide(t *testing.T)    { fuzzGen(t, &divideGen{}, grader.GradingMultipleChoice) }
func TestPolar(t *testing.T)     { fuzzGen(t, &polarGen{}, grader.GradingMultipleChoice) }
func TestDeMoivre(t *testing.T)  { fuzzGen(t, &deMoivreGen{}, grader.GradingMultipleChoice) }
func TestRoots(t *testing.T)     { fuzzGen(t, &rootsGen{}, grader.GradingMultipleChoice) }
