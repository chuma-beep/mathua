package abstract

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

func TestGroupDef(t *testing.T)      { fuzzGen(t, &groupDefGen{}, grader.GradingMultipleChoice) }
func TestGroupExamples(t *testing.T) { fuzzGen(t, &groupExamplesGen{}, grader.GradingMultipleChoice) }
func TestSubgroup(t *testing.T)      { fuzzGen(t, &subgroupGen{}, grader.GradingMultipleChoice) }
func TestRing(t *testing.T)          { fuzzGen(t, &ringGen{}, grader.GradingMultipleChoice) }
func TestHomomorphism(t *testing.T)  { fuzzGen(t, &homomorphismGen{}, grader.GradingMultipleChoice) }
