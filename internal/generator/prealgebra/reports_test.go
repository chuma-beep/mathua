package prealgebra

import (
	"strconv"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func TestDecimalComparisonAlwaysMatchesTheSymbol(t *testing.T) {
	gen := &decCompareGen{}
	for seed := int64(1); seed <= 3000; seed++ {
		p := gen.Generate(generator.GeneratorContext{Difficulty: 0.5, Seed: seed})
		a, err := strconv.ParseFloat(p.Facts["a"], 64)
		if err != nil {
			t.Fatalf("seed %d: bad a: %v", seed, err)
		}
		b, err := strconv.ParseFloat(p.Facts["b"], 64)
		if err != nil {
			t.Fatalf("seed %d: bad b: %v", seed, err)
		}
		switch p.Answer {
		case "<":
			if !(a < b) {
				t.Fatalf("seed %d: %v < %v is false", seed, a, b)
			}
		case ">":
			if !(a > b) {
				t.Fatalf("seed %d: %v > %v is false", seed, a, b)
			}
		case "=":
			if a != b {
				t.Fatalf("seed %d: %v = %v is false", seed, a, b)
			}
		default:
			t.Fatalf("seed %d: unknown symbol %q", seed, p.Answer)
		}
	}
}
