package prealgebra

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func TestNormalizeSciNotationAcceptsUnicodeOperators(t *testing.T) {
	cases := []string{
		"8.97∗10^3",
		"8.97⋅10^3",
		"8.97∙10^3",
		`8.97\ast10^3`,
	}
	for _, in := range cases {
		if got := normalizeSciNotation(in); got != "8.97x10^3" {
			t.Errorf("normalizeSciNotation(%q) = %q, want 8.97x10^3", in, got)
		}
	}
	g := &sciNotationOpsGen{}
	if res := g.Grade("8.97x10^3", "8.97∗10^3"); !res.Correct {
		t.Errorf("Grade rejected Unicode multiplication: %+v", res)
	}
}

// Report 438: a run of signed order-of-operations questions began with 0, as
// in "0 + (-3) x 3", which does not exercise a signed leading term.
func TestNegOrderOpsGenStartsWithNegative(t *testing.T) {
	g := &negOrderOpsGen{}
	for i := 0; i < 500; i++ {
		p := g.Generate(generator.GeneratorContext{Difficulty: rand.Float64()})
		a, _ := strconv.Atoi(p.Facts["a"])
		if a >= 0 {
			t.Fatalf("leading term %d is not negative in %q", a, p.Question)
		}
	}
}

// Report 439: 2/7 was rejected for 2:7 even though the sibling concept
// generator already accepts a slash as a ratio separator.
func TestRatioSimplifyAcceptsSlash(t *testing.T) {
	g := &ratioSimplifyGen{}
	cases := []struct {
		answer string
		want   bool
	}{
		{"2:7", true},
		{"2/7", true},
		{"2 to 7", true},
		{"7:2", false},
		{"10:35", false},
		{"two sevenths", false},
	}
	for _, c := range cases {
		if got := g.Grade("2:7", c.answer).Correct; got != c.want {
			t.Errorf("Grade(2:7, %q) = %v, want %v", c.answer, got, c.want)
		}
	}
}
