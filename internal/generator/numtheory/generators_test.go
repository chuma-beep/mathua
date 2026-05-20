package numtheory

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

func TestDivisibilityGen(t *testing.T)  { fuzzGen(t, &divisibilityGen{}) }
func TestGcdEuclideanGen(t *testing.T)  { fuzzGen(t, &gcdEuclideanGen{}) }
func TestModularGen(t *testing.T)       { fuzzGen(t, &modularGen{}) }
func TestCongruenceGen(t *testing.T)    { fuzzGen(t, &congruenceGen{}) }
func TestFermatLittleGen(t *testing.T)  { fuzzGen(t, &fermatLittleGen{}) }
func TestEulerPhiGen(t *testing.T)      { fuzzGen(t, &eulerPhiGen{}) }
func TestDiophantineGen(t *testing.T)   { fuzzGen(t, &diophantineGen{}) }
func TestCryptoGen(t *testing.T)        { fuzzGen(t, &cryptoGen{}) }
