package machinelearning

import (
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"unicode"

	"github.com/chuma-beep/mathua/internal/generator"
)

var convertedMLNumericIDs = []string{
	"ml.linreg.least_squares",
	"ml.linreg.predict",
	"ml.eval.accuracy",
	"ml.eval.precision",
	"ml.eval.recall",
	"ml.reg.ridge_penalty",
	"ml.pca.variance_percent",
	"ml.knn.nearest_distance",
}

func testMLRegistry(t *testing.T) *generator.Registry {
	t.Helper()
	reg := generator.NewRegistry()
	Register(reg)
	return reg
}

func TestMLBatchConversions_CognitiveScaffolding(t *testing.T) {
	reg := testMLRegistry(t)
	for _, id := range convertedMLNumericIDs {
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

func TestMLBatchConversions_DifficultyScaling(t *testing.T) {
	reg := testMLRegistry(t)
	for _, id := range convertedMLNumericIDs {
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
			if _, err := strconv.Atoi(prob.Answer); err != nil {
				t.Fatalf("%s: answer %q is not an integer", id, prob.Answer)
			}
		}
	}
}

func TestMLBatchConversions_ParameterTransparency(t *testing.T) {
	reg := testMLRegistry(t)
	checked := 0
	for _, id := range convertedMLNumericIDs {
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

func TestMLBatchConversions_ZeroOxfordComma(t *testing.T) {
	reg := testMLRegistry(t)
	oxfordComma := regexp.MustCompile(`,\s+and\s+`)
	for _, id := range convertedMLNumericIDs {
		for i := 0; i < 20; i++ {
			p, err := reg.GenerateContext(id, generator.GeneratorContext{
				Difficulty: float64(i%10) / 9.0,
				Seed:       int64(30000 + i),
			})
			if err != nil {
				t.Fatalf("%s: generate failed: %v", id, err)
			}
			if oxfordComma.MatchString(p.Explanation) {
				t.Fatalf("%s: explanation contains Oxford comma: %s", id, p.Explanation)
			}
			if oxfordComma.MatchString(p.Question) {
				t.Fatalf("%s: question contains Oxford comma: %s", id, p.Question)
			}
		}
	}
}

// TestMLBatchAnswerIsNotReadableOffTheQuestion guards the two linear-regression
// generators, where the answer and the printed parameters are the same kind of
// number and can therefore coincide. The least-squares slope used to equal a
// printed coordinate whenever the intercept was 0, was the slope or was its
// negative; a prediction could equal the slope or the input. A percentage equal
// to one of the counts is not this defect, which is why the guard is not applied
// to the metric generators.
func TestMLBatchAnswerIsNotReadableOffTheQuestion(t *testing.T) {
	cases := []struct {
		id       string
		excluded []string
		token    bool
	}{
		{"ml.linreg.least_squares", []string{"b", "y1", "y2"}, true},
		{"ml.linreg.predict", []string{"m", "x"}, false},
	}
	reg := testMLRegistry(t)
	for _, c := range cases {
		for i := 0; i < 500; i++ {
			p, err := reg.GenerateContext(c.id, generator.GeneratorContext{
				Difficulty: rand.Float64(),
				Seed:       int64(40000 + i),
			})
			if err != nil {
				t.Fatalf("%s: generate failed: %v", c.id, err)
			}
			for _, key := range c.excluded {
				if p.Facts[key] == p.Answer {
					t.Fatalf("%s: fact %s = %s is the answer itself: %s", c.id, key, p.Answer, p.Question)
				}
			}
			// A printed coordinate equal to the negative of the answer repeats it
			// as a token: the answer 10 beside the point (1, -10) is readable off
			// the stem even though no coordinate equals 10.
			if c.token && len(p.Answer) >= 2 && stemHasToken(p.Question, p.Answer) {
				t.Fatalf("%s: answer %s appears as a token in the stem: %s", c.id, p.Answer, p.Question)
			}
		}
	}
}

// TestMLBatchVarianceAnswersVary keeps the principal-component question from
// falling back to a fixed pool of eigenvalue ratios, which served seven answers
// however often a learner met the concept.
func TestMLBatchVarianceAnswersVary(t *testing.T) {
	reg := testMLRegistry(t)
	seen := map[string]bool{}
	for i := 0; i < 300; i++ {
		p, err := reg.GenerateContext("ml.pca.variance_percent", generator.GeneratorContext{
			Difficulty: 1,
			Seed:       int64(50000 + i),
		})
		if err != nil {
			t.Fatalf("generate failed: %v", err)
		}
		seen[p.Answer] = true
	}
	if len(seen) < 40 {
		t.Fatalf("explained-variance serves only %d distinct answers at full difficulty, want at least 40", len(seen))
	}
}

// stemHasToken reports whether token appears in stem as a whole alphanumeric
// token, so that a minus sign or a surrounding delimiter cannot hide it.
func stemHasToken(stem, token string) bool {
	fields := strings.FieldsFunc(stem, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
	for _, f := range fields {
		if f == token {
			return true
		}
	}
	return false
}
