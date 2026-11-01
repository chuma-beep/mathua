package statistics

import (
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func TestMeanAnswerIsExact(t *testing.T) {
	gen := &meanGen{}
	for seed := int64(1); seed <= 200; seed++ {
		p := gen.Generate(generator.GeneratorContext{Difficulty: 0.5, Seed: seed})
		sum, _ := strconv.Atoi(p.Facts["sum"])
		n, _ := strconv.Atoi(p.Facts["n"])
		if p.Answer != p.Facts["answer"] {
			t.Fatalf("seed %d: answer %q does not match fact %q", seed, p.Answer, p.Facts["answer"])
		}
		if strings.Contains(p.Answer, ".") {
			t.Fatalf("seed %d: answer %q is a rounded decimal", seed, p.Answer)
		}
		if sum%n == 0 {
			if p.Answer != strconv.Itoa(sum/n) {
				t.Fatalf("seed %d: exact mean should be %d, got %q", seed, sum/n, p.Answer)
			}
			continue
		}
		num, den, ok := strings.Cut(p.Answer, "/")
		if !ok {
			t.Fatalf("seed %d: non-integer mean %q is not a fraction", seed, p.Answer)
		}
		a, _ := strconv.Atoi(num)
		b, _ := strconv.Atoi(den)
		if a*n != sum*b {
			t.Fatalf("seed %d: %q does not equal %d/%d", seed, p.Answer, sum, n)
		}
	}
}
