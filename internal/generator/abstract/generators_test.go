package abstract

import (
	"math/rand"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func fuzzGen(t *testing.T, gen generator.Generator) {
	t.Helper()
	for i := 0; i < 100; i++ {
		d := rand.Float64()
		p := gen.Generate(d)
		if p.Question == "" || p.Answer == "" || p.Explanation == "" {
			t.Errorf("empty field at difficulty=%.2f: q=%q a=%q e=%q", d, p.Question, p.Answer, p.Explanation)
		}
	}
}

func TestGroupDefGen(t *testing.T)       { fuzzGen(t, &groupDefGen{}) }
func TestGroupExamplesGen(t *testing.T)  { fuzzGen(t, &groupExamplesGen{}) }
func TestSubgroupGen(t *testing.T)       { fuzzGen(t, &subgroupGen{}) }
func TestRingGen(t *testing.T)           { fuzzGen(t, &ringGen{}) }
func TestHomomorphismGen(t *testing.T)   { fuzzGen(t, &homomorphismGen{}) }
func TestFieldGen(t *testing.T)          { fuzzGen(t, &fieldGen{}) }
func TestModuleGen(t *testing.T)         { fuzzGen(t, &moduleGen{}) }
