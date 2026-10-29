package statistics

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/grader"
)

var convertedStatNumericIDs = []string{
	"stat.spread.mad",
	"stat.spread.iqr",
	"stat.quartile.value",
	"stat.zscore.value",
	"stat.sampling.se",
	"stat.infer.margin",
	"stat.prob.conditional",
}

func testStatRegistry(t *testing.T) *generator.Registry {
	t.Helper()
	reg := generator.NewRegistry()
	Register(reg)
	return reg
}

func TestStatBatchConversions_CognitiveScaffolding(t *testing.T) {
	reg := testStatRegistry(t)
	for _, id := range convertedStatNumericIDs {
		for i := 0; i < 30; i++ {
			p, err := reg.GenerateContext(id, generator.GeneratorContext{
				Difficulty: float64(i%10) / 9.0,
				Seed:       int64(10000 + i),
			})
			if err != nil {
				t.Fatalf("%s: generate failed: %v", id, err)
			}
			if !strings.Contains(p.Explanation, "Subgoal 1: State governing formula") {
				t.Fatalf("%s: missing Subgoal 1 in explanation: %s", id, p.Explanation)
			}
			if !strings.Contains(p.Explanation, "Subgoal 2: Substitute parameters") {
				t.Fatalf("%s: missing Subgoal 2 in explanation: %s", id, p.Explanation)
			}
			if !strings.Contains(p.Explanation, "Subgoal 3: Evaluate and simplify") {
				t.Fatalf("%s: missing Subgoal 3 in explanation: %s", id, p.Explanation)
			}
		}
	}
}

func TestStatBatchConversions_DifficultyScaling(t *testing.T) {
	reg := testStatRegistry(t)
	g := grader.NewRouter()
	for _, id := range convertedStatNumericIDs {
		gen, err := reg.Get(id)
		if err != nil {
			t.Fatalf("generator not found: %s: %v", id, err)
		}
		for i := 0; i < 50; i++ {
			d := rand.Float64()
			prob := gen.Generate(generator.GeneratorContext{Difficulty: d})
			if prob.Question == "" || prob.Answer == "" || prob.Explanation == "" {
				t.Fatalf("%s: empty problem output at difficulty %.2f", id, d)
			}
			if prob.Facts == nil || prob.Facts["answer"] != prob.Answer {
				t.Fatalf("%s: fact answer missing or mismatch: got %q, want %q", id, prob.Facts["answer"], prob.Answer)
			}
			// Verify grading passes with numeric grader
			res := g.Grade(grader.GradingNumeric, prob.Answer, prob.Answer)
			if !res.Correct {
				t.Fatalf("%s: self-grading failed for answer %q: %s", id, prob.Answer, res.Feedback)
			}
		}
	}
}

func TestStatBatchConversions_ParameterTransparency(t *testing.T) {
	reg := testStatRegistry(t)
	checked := 0
	for _, id := range convertedStatNumericIDs {
		for i := 0; i < 60; i++ {
			p, err := reg.GenerateContext(id, generator.GeneratorContext{
				Difficulty: float64(i%10) / 9.0,
				Seed:       int64(20000 + i),
			})
			if err != nil {
				t.Fatalf("%s: generate failed: %v", id, err)
			}
			for key, raw := range p.Facts {
				if key == "answer" {
					continue
				}
				value, err := strconv.Atoi(raw)
				if err != nil {
					continue
				}
				absVal := value
				if absVal < 0 {
					absVal = -absVal
				}
				pattern := regexp.MustCompile(`(^|\D)` + regexp.QuoteMeta(strconv.Itoa(absVal)) + `(\D|$)`)
				if !pattern.MatchString(p.Question) {
					t.Fatalf("%s: names parameter %s = %d, but question does not print it: %s",
						id, key, value, p.Question)
				}
				checked++
			}
		}
	}
	if checked < 400 {
		t.Fatalf("only %d facts checked", checked)
	}
}

func TestStatBatchConversions_ZeroOxfordComma(t *testing.T) {
	reg := testStatRegistry(t)
	oxfordComma := regexp.MustCompile(`,\s+and\s+`)
	for _, id := range convertedStatNumericIDs {
		for i := 0; i < 20; i++ {
			p, err := reg.GenerateContext(id, generator.GeneratorContext{
				Difficulty: float64(i%10) / 9.0,
				Seed:       int64(30000 + i),
			})
			if err != nil {
				t.Fatalf("%s: generate failed: %v", id, err)
			}
			if oxfordComma.MatchString(p.Question) {
				t.Fatalf("%s: question contains Oxford comma: %s", id, p.Question)
			}
			if oxfordComma.MatchString(p.Answer) {
				t.Fatalf("%s: answer contains Oxford comma: %s", id, p.Answer)
			}
			if oxfordComma.MatchString(p.Explanation) {
				t.Fatalf("%s: explanation contains Oxford comma: %s", id, p.Explanation)
			}
		}
	}
}

func TestStatBatchAnswerIsNotReadableOffTheQuestion(t *testing.T) {
	reg := testStatRegistry(t)
	for _, id := range convertedStatNumericIDs {
		for i := 0; i < 500; i++ {
			p, err := reg.GenerateContext(id, generator.GeneratorContext{
				Difficulty: rand.Float64(),
				Seed:       int64(40000 + i),
			})
			if err != nil {
				t.Fatalf("%s: generate failed: %v", id, err)
			}
			switch id {
			case "stat.zscore.value":
				for _, key := range []string{"mean", "sd", "x"} {
					if p.Facts[key] == p.Answer {
						t.Fatalf("%s: fact %s = %s is the answer itself: %s", id, key, p.Answer, p.Question)
					}
				}
			case "stat.sampling.se":
				for _, key := range []string{"size", "sd"} {
					if p.Facts[key] == p.Answer {
						t.Fatalf("%s: fact %s = %s is the answer itself: %s", id, key, p.Answer, p.Question)
					}
				}
			case "stat.infer.margin":
				for _, key := range []string{"z", "se"} {
					if p.Facts[key] == p.Answer {
						t.Fatalf("%s: fact %s = %s is the answer itself: %s", id, key, p.Answer, p.Question)
					}
				}
			case "stat.prob.conditional":
				bothStr := p.Facts["both"]
				bStr := p.Facts["b"]
				dStr := p.Facts["d"]
				frac1 := fmt.Sprintf("%s/%s", bothStr, dStr)
				frac2 := fmt.Sprintf("%s/%s", bStr, dStr)
				if p.Answer == frac1 || p.Answer == frac2 {
					t.Fatalf("%s: answer %s matches input fraction: %s", id, p.Answer, p.Question)
				}
			case "stat.spread.mad", "stat.spread.iqr", "stat.quartile.value":
				// The answer must not appear as a standalone element in the set brackets \{...\}
				re := regexp.MustCompile(`\\\{([^}]+)\\\}`)
				m := re.FindStringSubmatch(p.Question)
				if m != nil {
					elements := strings.Split(m[1], ",")
					for _, elem := range elements {
						if strings.TrimSpace(elem) == p.Answer {
							t.Fatalf("%s: answer %s is directly present in data set: %s", id, p.Answer, p.Question)
						}
					}
				}
			}
		}
	}
}

func TestStatBatchAnswersVary(t *testing.T) {
	reg := testStatRegistry(t)
	thresholds := map[string]int{
		"stat.spread.mad":       15,
		"stat.spread.iqr":       15,
		"stat.quartile.value":   20,
		"stat.zscore.value":     6,
		"stat.sampling.se":      15,
		"stat.infer.margin":     20,
		"stat.prob.conditional": 20,
	}
	for _, id := range convertedStatNumericIDs {
		seen := map[string]bool{}
		for i := 0; i < 200; i++ {
			p, err := reg.GenerateContext(id, generator.GeneratorContext{
				Difficulty: 1.0,
				Seed:       int64(50000 + i),
			})
			if err != nil {
				t.Fatalf("%s: generate failed: %v", id, err)
			}
			seen[p.Answer] = true
		}
		want := thresholds[id]
		if want == 0 {
			want = 15
		}
		if len(seen) < want {
			t.Fatalf("%s serves only %d distinct answers at full difficulty, want at least %d", id, len(seen), want)
		}
	}
}

func TestZScoreExplanationNamesTheSignRule(t *testing.T) {
	reg := testStatRegistry(t)
	const rule = "A value above the mean has a positive z-score and a value below the mean has a negative z-score."
	for i := 0; i < 200; i++ {
		p, err := reg.GenerateContext("stat.zscore.value", generator.GeneratorContext{
			Difficulty: 0.8,
			Seed:       int64(60000 + i),
		})
		if err != nil {
			t.Fatalf("generate failed: %v", err)
		}
		if !strings.Contains(p.Explanation, rule) {
			t.Fatalf("explanation does not state the sign rule: %s", p.Explanation)
		}
		dist := statProbeMeanSd.FindStringSubmatch(p.Question)
		value := statProbeValue.FindStringSubmatch(p.Question)
		if dist == nil || value == nil {
			t.Fatalf("question does not print mean, deviation and value: %q", p.Question)
		}
		mean, _ := strconv.Atoi(dist[1])
		x, _ := strconv.Atoi(value[1])
		answer, _ := strconv.Atoi(p.Answer)
		if x > mean && answer <= 0 {
			t.Fatalf("value above the mean has a non-positive z-score: %s", p.Question)
		}
		if x < mean && answer >= 0 {
			t.Fatalf("value below the mean has a non-negative z-score: %s", p.Question)
		}
	}
}
