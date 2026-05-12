package counting

import (
	"math/rand"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/grader"
)

func fuzzGen(t *testing.T, gen generator.Generator, gradingType grader.GradingType, n int) {
	t.Helper()
	gr := grader.NewRouter()
	for i := 0; i < n; i++ {
		p := gen.Generate(rand.Float64())
		if p.Question == "" {
			t.Error("empty question")
		}
		if p.Answer == "" {
			t.Error("empty answer")
		}
		if p.Explanation == "" {
			t.Error("empty explanation")
		}
		res := gr.Grade(gradingType, p.Answer, p.Answer)
		if !res.Correct {
			t.Errorf("grader could not validate own answer: q=%q a=%q", p.Question, p.Answer)
		}
	}
}

func TestCountObjects(t *testing.T)      { fuzzGen(t, &countObjectsGen{}, grader.GradingNumeric, 100) }
func TestCountCardinality(t *testing.T)  { fuzzGen(t, &countCardinalityGen{}, grader.GradingNumeric, 100) }
func TestCountNumberLine(t *testing.T)   { fuzzGen(t, &countNumberLineGen{}, grader.GradingNumeric, 100) }
func TestCountCompare(t *testing.T)      { fuzzGen(t, &countCompareGen{}, grader.GradingComparison, 100) }
func TestSkip2(t *testing.T)             { fuzzGen(t, &skipCountGen{step: 2, max: 20}, grader.GradingNumeric, 100) }
func TestSkip5(t *testing.T)             { fuzzGen(t, &skipCountGen{step: 5, max: 50}, grader.GradingNumeric, 100) }
func TestSkip10(t *testing.T)            { fuzzGen(t, &skipCountGen{step: 10, max: 100}, grader.GradingNumeric, 100) }
func TestCountObjects20(t *testing.T)    { fuzzGen(t, &countObjects20Gen{}, grader.GradingNumeric, 100) }
func TestCountOrdinal(t *testing.T)      { fuzzGen(t, &countOrdinalGen{}, grader.GradingMultipleChoice, 100) }
func TestCountBackwards(t *testing.T)    { fuzzGen(t, &countBackwardsGen{}, grader.GradingNumeric, 100) }
