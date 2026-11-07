package prealgebra

import "testing"

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
