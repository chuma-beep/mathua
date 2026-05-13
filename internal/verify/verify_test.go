package verify

import (
	"testing"
)

func TestEval(t *testing.T) {
	tests := []struct {
		expr     string
		expected float64
	}{
		{"2+2", 4},
		{"3*4", 12},
		{"10-3", 7},
		{"20/4", 5},
		{"2**3", 8},
		{"3+4*2", 11},
		{"(3+4)*2", 14},
		{"25/100*80", 20},
		{"100-35", 65},
	}
	for _, tt := range tests {
		got, err := Eval(tt.expr)
		if err != nil {
			t.Errorf("Eval(%q) error: %v", tt.expr, err)
			continue
		}
		if got != tt.expected {
			t.Errorf("Eval(%q) = %g, want %g", tt.expr, got, tt.expected)
		}
	}
}

func TestCheckExpr(t *testing.T) {
	if err := CheckExpr("2+2", "4"); err != nil {
		t.Errorf("CheckExpr(2+2, 4) failed: %v", err)
	}
	if err := CheckExpr("3*4", "12"); err != nil {
		t.Errorf("CheckExpr(3*4, 12) failed: %v", err)
	}
	if err := CheckExpr("2+2", "5"); err == nil {
		t.Error("CheckExpr(2+2, 5) should have failed")
	}
}

func TestExtractExpr(t *testing.T) {
	tests := []struct {
		question string
		expected string
		ok       bool
	}{
		{"5 + 3 = ?", "5+3", true},
		{"5+3=?", "5+3", true},
		{"12 - 7 = ?", "12-(7)", true},
		{"4 x 3 = ?", "4*3", true},
		{"8 * 2 = ?", "8*2", true},
		{"What is 25% of 80?", "25/100*80", true},
		{"random text with no pattern", "", false},
	}
	for _, tt := range tests {
		got, ok := ExtractExpr(tt.question)
		if ok != tt.ok {
			t.Errorf("ExtractExpr(%q) ok=%v, want %v", tt.question, ok, tt.ok)
			continue
		}
		if ok && got != tt.expected {
			t.Errorf("ExtractExpr(%q) = %q, want %q", tt.question, got, tt.expected)
		}
	}
}
