package statistics

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

// The variance is a rational. Formatting it with %.0f graded a set whose true
// variance is 50.56 against 51 and printed the same rounded number in the
// working. Every sample must now carry the exact value in lowest terms.
func TestVarianceAnswerIsExact(t *testing.T) {
	g := &varianceGen{}
	for i := 0; i < 200; i++ {
		p := g.Generate(generator.GeneratorContext{Difficulty: rand.Float64(), Seed: rand.Int63()})
		rawParts := strings.SplitN(p.Question, "of: ", 2)
		if len(rawParts) != 2 {
			t.Fatalf("question has no value list: %q", p.Question)
		}
		parts := strings.Split(rawParts[1], ", ")
		if len(parts) != 5 {
			t.Fatalf("question has %d values: %q", len(parts), p.Question)
		}
		sum, sumSq := 0, 0
		for _, part := range parts {
			v, err := strconv.Atoi(part)
			if err != nil {
				t.Fatalf("value %q: %v", part, err)
			}
			sum += v
			sumSq += v * v
		}
		num := 5*sumSq - sum*sum
		den := 25
		div := varianceGCD(num, den)
		num /= div
		den /= div
		want := strconv.Itoa(num)
		if den != 1 {
			want = fmt.Sprintf("%d/%d", num, den)
		}
		if p.Answer != want {
			t.Fatalf("answer = %q, want %q for %q", p.Answer, want, p.Question)
		}
	}
}
