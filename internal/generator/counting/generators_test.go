package counting

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

func TestCountObjectsGen(t *testing.T)    { fuzzGen(t, &countObjectsGen{}) }
func TestCountCardinalityGen(t *testing.T) { fuzzGen(t, &countCardinalityGen{}) }
func TestCountNumberLineGen(t *testing.T)  { fuzzGen(t, &countNumberLineGen{}) }
func TestCountCompareGen(t *testing.T)     { fuzzGen(t, &countCompareGen{}) }
func TestSkipCountGen(t *testing.T)        { fuzzGen(t, &skipCountGen{}) }
func TestCountObjects20Gen(t *testing.T)   { fuzzGen(t, &countObjects20Gen{}) }
func TestCountOrdinalGen(t *testing.T)     { fuzzGen(t, &countOrdinalGen{}) }
func TestCountBackwardsGen(t *testing.T)   { fuzzGen(t, &countBackwardsGen{}) }
