package prealgebra

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func TestFractionEquationAnswerIsCorrect(t *testing.T) {
	gen := &prealgBridgeGen{kind: "eq_fractions"}
	pattern := regexp.MustCompile(`\((\d+)/(\d+)\) x = (\d+)`)

	for i := 0; i < 300; i++ {
		p := gen.Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: 3})
		match := pattern.FindStringSubmatch(p.Question)
		if match == nil {
			t.Fatalf("question does not match the expected form: %q", p.Question)
		}
		a, _ := strconv.Atoi(match[1])
		b, _ := strconv.Atoi(match[2])
		rhs, _ := strconv.Atoi(match[3])

		num := rhs * b
		den := a
		g := ratioGCD(num, den)
		num /= g
		den /= g
		want := strconv.Itoa(num)
		if den != 1 {
			want = fmt.Sprintf("%d/%d", num, den)
		}
		if p.Answer != want {
			t.Fatalf("question %q answers %q, want %q", p.Question, p.Answer, want)
		}
	}
}
