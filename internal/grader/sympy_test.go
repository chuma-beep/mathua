package grader

import (
	"os/exec"
	"testing"
)

func sympyAvailable(t *testing.T) bool {
	t.Helper()
	if err := exec.Command("python3", "-c", "import sympy").Run(); err != nil {
		t.Log("SymPy not available, skipping test")
		return false
	}
	if findSymPyService() == "" {
		t.Log("sympy_service.py not found, skipping test")
		return false
	}
	return true
}

func TestSympyGrade_EquivalentExpressions(t *testing.T) {
	if !sympyAvailable(t) {
		t.Skip("SymPy not available")
	}
	r := NewRouter()
	tests := []struct {
		expected, answer string
	}{
		{"x^2+2x+1", "(x+1)^2"},
		{"6x^2", "6*x^2"},
		{"(x+3)(x+4)", "x^2+7x+12"},
		{"e^x", "exp(x)"},
		{"2e^(2x)", "2*exp(2*x)"},
		{"cos(x)", "cos(x)"},
		{"3(2x+1)^2", "12x^2+12x+3"},
		{"y = (z - 2x)/3", "(z-2*x)/3"},
		{"xe^x-e^x+C", "e^x*(x-1)+C"},
		{"(x+5)(x-5)", "x^2-25"},
		{"-2", "-2"},
	}
	for _, tc := range tests {
		res := r.Grade(GradingExpression, tc.expected, tc.answer)
		if !res.Correct {
			t.Errorf("expected %q ≈ %q to be correct, got feedback: %q", tc.expected, tc.answer, res.Feedback)
		}
	}
}

func TestSympyGrade_NonequivalentExpressions(t *testing.T) {
	if !sympyAvailable(t) {
		t.Skip("SymPy not available")
	}
	r := NewRouter()
	tests := []struct {
		expected, answer string
	}{
		{"2x+3", "2x+4"},
		{"x^2", "x^3"},
		{"cos(x)", "sin(x)"},
	}
	for _, tc := range tests {
		res := r.Grade(GradingExpression, tc.expected, tc.answer)
		if res.Correct {
			t.Errorf("expected %q ≠ %q to be incorrect", tc.expected, tc.answer)
		}
	}
}

func TestSympyGrade_ParseError(t *testing.T) {
	if !sympyAvailable(t) {
		t.Skip("SymPy not available")
	}
	r := NewRouter()
	res := r.Grade(GradingExpression, "x^2", "invalid(((")
	if res.Correct {
		t.Error("expected parse error to be incorrect")
	}
	if res.Feedback == "" || res.Feedback == "Incorrect" {
		t.Errorf("expected specific parse error feedback, got %q", res.Feedback)
	}
}

func TestSympyGrade_ExpressionAndPolynomialBothRoute(t *testing.T) {
	if !sympyAvailable(t) {
		t.Skip("SymPy not available")
	}
	r := NewRouter()
	tests := []struct {
		gt               GradingType
		expected, answer string
	}{
		{GradingExpression, "x^2+2x+1", "(x+1)^2"},
		{GradingPolynomial, "6x^2", "6*x^2"},
	}
	for _, tc := range tests {
		res := r.Grade(tc.gt, tc.expected, tc.answer)
		if !res.Correct {
			t.Errorf("expected %s(%q ≈ %q) to be correct", tc.gt, tc.expected, tc.answer)
		}
	}
}

func TestSympyGrade_MultipleRequests(t *testing.T) {
	if !sympyAvailable(t) {
		t.Skip("SymPy not available")
	}
	r := NewRouter()
	for i := 0; i < 5; i++ {
		res := r.Grade(GradingExpression, "x^2+2x+1", "(x+1)^2")
		if !res.Correct {
			t.Errorf("request %d: expected correct, got %q", i, res.Feedback)
		}
	}
}
