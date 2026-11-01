package algebra

import (
	"strconv"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/mathutil"
)

func TestExponentialEquationAnswerMatchesRightHandSide(t *testing.T) {
	gen := &eqExpGen{}
	for seed := int64(1); seed <= 200; seed++ {
		p := gen.Generate(generator.GeneratorContext{Difficulty: 0.5, Seed: seed})
		base, err := strconv.Atoi(p.Facts["base"])
		if err != nil {
			t.Fatalf("seed %d: bad base: %v", seed, err)
		}
		rhs, err := strconv.Atoi(p.Facts["rhs"])
		if err != nil {
			t.Fatalf("seed %d: bad rhs: %v", seed, err)
		}
		x, err := strconv.Atoi(p.Answer)
		if err != nil {
			t.Fatalf("seed %d: bad answer %q: %v", seed, p.Answer, err)
		}
		if got := mathutil.IntPow(base, x); got != rhs {
			t.Fatalf("seed %d: %d^%d = %d, want %d", seed, base, x, got, rhs)
		}
	}
}
