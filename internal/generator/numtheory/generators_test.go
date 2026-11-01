package numtheory

import (
	"math/rand"
	"strconv"
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

func TestDivisibilityGen(t *testing.T) { fuzzGen(t, &divisibilityGen{}) }
func TestGcdEuclideanGen(t *testing.T) { fuzzGen(t, &gcdEuclideanGen{}) }
func TestModularGen(t *testing.T)      { fuzzGen(t, &modularGen{}) }
func TestCongruenceGen(t *testing.T)   { fuzzGen(t, &congruenceGen{}) }
func TestFermatLittleGen(t *testing.T) { fuzzGen(t, &fermatLittleGen{}) }
func TestEulerPhiGen(t *testing.T)     { fuzzGen(t, &eulerPhiGen{}) }
func TestDiophantineGen(t *testing.T)  { fuzzGen(t, &diophantineGen{}) }
func TestCryptoGen(t *testing.T)       { fuzzGen(t, &cryptoGen{}) }

func TestFuzz(t *testing.T) {
	reg := generator.NewRegistry()
	Register(reg)
	for _, id := range reg.Concepts() {
		gen, _ := reg.Get(id)
		t.Run(id, func(t *testing.T) {
			for i := 0; i < 100; i++ {
				d := rand.Float64()
				p := gen.Generate(generator.GeneratorContext{Difficulty: d})
				if p.Question == "" || p.Answer == "" || p.Explanation == "" {
					t.Fatalf("empty field at difficulty=%.2f for %s: q=%q a=%q", d, id, p.Question, p.Answer)
				}
			}
		})
	}
}

func TestDiophantineAnswerSolvesTheEquation(t *testing.T) {
	g := &diophantineGen{}
	for i := 0; i < 100; i++ {
		p := g.Generate(generator.GeneratorContext{Difficulty: rand.Float64()})
		c, errC := strconv.Atoi(p.Facts["c"])
		x, errX := strconv.Atoi(p.Facts["x"])
		y, errY := strconv.Atoi(p.Answer)
		if errC != nil || errX != nil || errY != nil {
			t.Fatalf("non-integer diophantine data: c=%q x=%q answer=%q", p.Facts["c"], p.Facts["x"], p.Answer)
		}
		if 3*x+5*y != c {
			t.Fatalf("answer y=%d with x=%d does not solve 3x+5y=%d", y, x, c)
		}
	}
}
