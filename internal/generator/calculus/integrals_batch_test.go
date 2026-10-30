package calculus

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/grader"
)

var integralsBatchIDs = []string{
	"calc.integral.area_between",
	"calc.integral.net_change",
	"calc.integral.mvt",
	"calc.integral.work",
	"calc.integral.center_mass",
	"calc.integral.trapezoid",
	"calc.integral.simpson",
}

func integralsBatchRegistry(t *testing.T) *generator.Registry {
	t.Helper()
	reg := generator.NewRegistry()
	Register(reg)
	return reg
}

func TestIntegralsBatch_ShippedTypesAreNumeric(t *testing.T) {
	shipped := calcShippedTypes(t)
	for _, id := range integralsBatchIDs {
		gt, ok := shipped[id]
		if !ok {
			t.Errorf("concept %s missing from the corpus", id)
			continue
		}
		if gt != "numeric" {
			t.Errorf("concept %s: want grading_type numeric, got %q", id, gt)
		}
	}
}

func TestIntegralsBatch_AnswersParseAsFloats(t *testing.T) {
	reg := integralsBatchRegistry(t)
	for _, id := range integralsBatchIDs {
		for i := 0; i < 50; i++ {
			p, err := reg.GenerateContext(id, generator.GeneratorContext{
				Difficulty: rand.Float64(),
				Seed:       int64(i + 1),
			})
			if err != nil {
				t.Fatalf("concept %s generate: %v", id, err)
			}
			if _, err := strconv.ParseFloat(strings.TrimSpace(p.Answer), 64); err != nil {
				t.Fatalf("concept %s (iter %d): answer %q not parseable float: %v", id, i, p.Answer, err)
			}
		}
	}
}

func TestIntegralsBatch_GradeCorrectly(t *testing.T) {
	reg := integralsBatchRegistry(t)
	router := grader.NewRouter()
	for _, id := range integralsBatchIDs {
		for i := 0; i < 30; i++ {
			p, err := reg.GenerateContext(id, generator.GeneratorContext{
				Difficulty: rand.Float64(),
				Seed:       int64(i + 42),
			})
			if err != nil {
				t.Fatalf("concept %s generate: %v", id, err)
			}
			res := router.Grade(grader.GradingNumeric, p.Answer, p.Answer)
			if !res.Correct {
				t.Fatalf("concept %s: correct answer %q failed grading against itself", id, p.Answer)
			}
		}
	}
}

func TestIntegralsBatch_CognitiveScaffolding(t *testing.T) {
	reg := integralsBatchRegistry(t)
	for _, id := range integralsBatchIDs {
		for i := 0; i < 20; i++ {
			p, err := reg.GenerateContext(id, generator.GeneratorContext{
				Difficulty: rand.Float64(),
				Seed:       int64(i + 100),
			})
			if err != nil {
				t.Fatalf("concept %s generate: %v", id, err)
			}
			if !strings.Contains(p.Explanation, "Subgoal 1:") {
				t.Fatalf("concept %s: explanation missing Subgoal 1: %s", id, p.Explanation)
			}
			if !strings.Contains(p.Explanation, "Subgoal 2:") {
				t.Fatalf("concept %s: explanation missing Subgoal 2: %s", id, p.Explanation)
			}
			if !strings.Contains(p.Explanation, "Subgoal 3:") {
				t.Fatalf("concept %s: explanation missing Subgoal 3: %s", id, p.Explanation)
			}
		}
	}
}

func TestIntegralsBatch_DifficultyScaling(t *testing.T) {
	reg := integralsBatchRegistry(t)
	for _, id := range integralsBatchIDs {
		pLow, err := reg.GenerateContext(id, generator.GeneratorContext{Difficulty: 0.0, Seed: 123})
		if err != nil {
			t.Fatalf("concept %s generate low: %v", id, err)
		}
		pHigh, err := reg.GenerateContext(id, generator.GeneratorContext{Difficulty: 1.0, Seed: 123})
		if err != nil {
			t.Fatalf("concept %s generate high: %v", id, err)
		}
		if pLow.Question == "" || pLow.Answer == "" {
			t.Errorf("concept %s: empty problem at difficulty 0.0", id)
		}
		if pHigh.Question == "" || pHigh.Answer == "" {
			t.Errorf("concept %s: empty problem at difficulty 1.0", id)
		}
	}
}

func TestIntegralsBatch_NumericRejectsNearMisses(t *testing.T) {
	reg := integralsBatchRegistry(t)
	router := grader.NewRouter()
	for _, id := range integralsBatchIDs {
		for i := 0; i < 20; i++ {
			p, err := reg.GenerateContext(id, generator.GeneratorContext{
				Difficulty: rand.Float64(),
				Seed:       int64(i + 200),
			})
			if err != nil {
				t.Fatalf("concept %s generate: %v", id, err)
			}
			val, _ := strconv.ParseFloat(strings.TrimSpace(p.Answer), 64)
			wrong := fmt.Sprintf("%.2f", val+5.0)
			res := router.Grade(grader.GradingNumeric, p.Answer, wrong)
			if res.Correct {
				t.Fatalf("concept %s: near miss %s incorrectly accepted for expected %s", id, wrong, p.Answer)
			}
		}
	}
}
