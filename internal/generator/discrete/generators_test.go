package discrete

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

func TestPropositions(t *testing.T) { fuzzGen(t, &propositionsGen{}, grader.GradingMultipleChoice) }
func TestConnectives(t *testing.T)  { fuzzGen(t, &connectivesGen{}, grader.GradingMultipleChoice) }
func TestTruthTables(t *testing.T)  { fuzzGen(t, &truthTablesGen{}, grader.GradingMultipleChoice) }
func TestQuantifiers(t *testing.T)  { fuzzGen(t, &quantifiersGen{}, grader.GradingMultipleChoice) }
func TestSetOps(t *testing.T)       { fuzzGen(t, &setOpsGen{}, grader.GradingMultipleChoice) }
func TestVenn(t *testing.T)         { fuzzGen(t, &vennGen{}, grader.GradingMultipleChoice) }
func TestPermutations(t *testing.T) { fuzzGen(t, &permutationsGen{}, grader.GradingNumeric) }
func TestCombinations(t *testing.T) { fuzzGen(t, &combinationsGen{}, grader.GradingNumeric) }
func TestPascal(t *testing.T)       { fuzzGen(t, &pascalGen{}, grader.GradingMultipleChoice) }
func TestGraphBasics(t *testing.T)  { fuzzGen(t, &graphBasicsGen{}, grader.GradingNumeric) }
func TestGraphPaths(t *testing.T)   { fuzzGen(t, &graphPathsGen{}, grader.GradingMultipleChoice) }
func TestGraphTrees(t *testing.T)   { fuzzGen(t, &treesGen{}, grader.GradingNumeric) }
func TestRecurrence(t *testing.T)   { fuzzGen(t, &recurrenceGen{}, grader.GradingNumeric) }
func TestInduction(t *testing.T)    { fuzzGen(t, &inductionGen{}, grader.GradingMultipleChoice) }
