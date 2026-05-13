package statistics

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

func TestReadTable(t *testing.T)      { fuzzGen(t, &readTableGen{}, grader.GradingNumeric) }
func TestBarGraph(t *testing.T)       { fuzzGen(t, &barGraphGen{}, grader.GradingMultipleChoice) }
func TestLinePlot(t *testing.T)       { fuzzGen(t, &linePlotGen{}, grader.GradingNumeric) }
func TestMean(t *testing.T)           { fuzzGen(t, &meanGen{}, grader.GradingNumeric) }
func TestMedian(t *testing.T)         { fuzzGen(t, &medianGen{}, grader.GradingNumeric) }
func TestMode(t *testing.T)           { fuzzGen(t, &modeGen{}, grader.GradingNumeric) }
func TestRange(t *testing.T)          { fuzzGen(t, &rangeGen{}, grader.GradingNumeric) }
func TestSampleSpace(t *testing.T)    { fuzzGen(t, &sampleSpaceGen{}, grader.GradingNumeric) }
func TestBasicProb(t *testing.T)      { fuzzGen(t, &basicProbGen{}, grader.GradingNumeric) }
func TestComplementProb(t *testing.T) { fuzzGen(t, &complementProbGen{}, grader.GradingNumeric) }
func TestCompoundProb(t *testing.T)   { fuzzGen(t, &compoundProbGen{}, grader.GradingNumeric) }
func TestCounting(t *testing.T)       { fuzzGen(t, &countingGen{}, grader.GradingNumeric) }
