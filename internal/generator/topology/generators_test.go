package topology

import (
	"math/rand"
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

func TestMetricGen(t *testing.T)     { fuzzGen(t, &metricGen{}) }
func TestOpenClosedGen(t *testing.T) { fuzzGen(t, &openClosedGen{}) }
func TestContinuousGen(t *testing.T) { fuzzGen(t, &continuousGen{}) }
