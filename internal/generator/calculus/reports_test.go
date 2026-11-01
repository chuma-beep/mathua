package calculus

import (
	"strconv"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/mathutil"
)

func TestDefiniteIntegralIsExact(t *testing.T) {
	gen := &integralDefiniteGen{}
	for seed := int64(1); seed <= 200; seed++ {
		p := gen.Generate(generator.GeneratorContext{Difficulty: 0.5, Seed: seed})
		a, _ := strconv.Atoi(p.Facts["a"])
		n, _ := strconv.Atoi(p.Facts["n"])
		upper, _ := strconv.Atoi(p.Facts["upper"])
		got, err := strconv.Atoi(p.Answer)
		if err != nil {
			t.Fatalf("seed %d: bad answer %q: %v", seed, p.Answer, err)
		}
		wantNum := a * mathutil.IntPow(upper, n+1)
		if wantNum%(n+1) != 0 {
			t.Fatalf("seed %d: %d is not divisible by %d", seed, wantNum, n+1)
		}
		if want := wantNum / (n + 1); got != want {
			t.Fatalf("seed %d: answer %d, want %d for %d x^%d on [0,%d]", seed, got, want, a, n, upper)
		}
	}
}
