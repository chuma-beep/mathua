package grader

import (
	"strings"
	"testing"
)

func TestGrader_Router_CloseAndDefault(t *testing.T) {
	r := NewRouter()
	defer r.Close()

	// Default fallback routes to tuple grader
	res := r.Grade(GradingType("unknown_custom_type"), "(3, 4)", "(3, 4)")
	if !res.Correct {
		t.Errorf("expected default fallback to grade matching tuple correctly, got %+v", res)
	}

	resWrong := r.Grade(GradingType("unknown_custom_type"), "(3, 4)", "(3, 5)")
	if resWrong.Correct {
		t.Errorf("expected default fallback to grade mismatch as incorrect, got %+v", resWrong)
	}
}

func TestGrader_Comparison_InvalidOrEmptyExpected(t *testing.T) {
	r := NewRouter()
	defer r.Close()

	// Empty expected returns empty result
	resEmpty := r.Grade(GradingComparison, "", ">")
	if resEmpty.Correct || resEmpty.Score != 0 {
		t.Errorf("expected empty result for empty expected, got %+v", resEmpty)
	}

	// Invalid expected comparison operator
	resInvalid := r.Grade(GradingComparison, "?", ">")
	if resInvalid.Correct || !strings.Contains(resInvalid.Feedback, "Invalid comparison operator") {
		t.Errorf("expected invalid operator feedback for expected, got %+v", resInvalid)
	}

	// Invalid answer operator
	resInvalidAns := r.Grade(GradingComparison, ">", "?")
	if resInvalidAns.Correct || !strings.Contains(resInvalidAns.Feedback, "Invalid comparison operator") {
		t.Errorf("expected invalid operator feedback for answer, got %+v", resInvalidAns)
	}
}

func TestGrader_Choice_ExpectedSingleLetterInAnswer(t *testing.T) {
	r := NewRouter()
	defer r.Close()

	// Expected is single letter, answer contains that letter
	resMatch := r.Grade(GradingMultipleChoice, "B", "The correct answer is B.")
	if !resMatch.Correct {
		t.Errorf("expected single letter B in answer to match, got %+v", resMatch)
	}

	// Expected single letter not present in answer
	resNoMatch := r.Grade(GradingMultipleChoice, "C", "The choice is D.")
	if resNoMatch.Correct {
		t.Errorf("expected mismatch when letter is absent, got %+v", resNoMatch)
	}

	// Empty expected
	resEmpty := r.Grade(GradingMultipleChoice, "", "A")
	if resEmpty.Correct || resEmpty.Score != 0 {
		t.Errorf("expected empty result for empty expected, got %+v", resEmpty)
	}
}

func TestGrader_Numeric_ScientificAndMixedNumbers(t *testing.T) {
	r := NewRouter()
	defer r.Close()

	cases := []struct {
		expected string
		answer   string
		correct  bool
	}{
		{"0.0015", "1.5e-3", true},
		{"32000", "3.2e4", true},
		{"3.25", "3 and 1/4", true},
		{"-2.25", "-2 1/4", true},
		{"0.75", ".75", true},
		{"-0.5", "-.5", true},
		{"42", "--42", false},
		{"42", "1.2.3", false},
		{"42", "5/", false},
	}

	for _, c := range cases {
		res := r.Grade(GradingNumeric, c.expected, c.answer)
		if res.Correct != c.correct {
			t.Errorf("Grade(GradingNumeric, %q, %q): got correct=%v, want %v (feedback=%s)", c.expected, c.answer, res.Correct, c.correct, res.Feedback)
		}
	}
}

func TestGrader_Tuple_Mismatches(t *testing.T) {
	r := NewRouter()
	defer r.Close()

	// Length mismatch
	resLen := r.Grade(GradingTuple, "(1, 2)", "(1, 2, 3)")
	if resLen.Correct || !strings.Contains(resLen.Feedback, "Expected 1,2, got 1,2,3") {
		t.Errorf("expected length mismatch feedback, got %+v", resLen)
	}

	// Positional mismatch
	resPos := r.Grade(GradingTuple, "(1, 2)", "(1, 9)")
	if resPos.Correct || !strings.Contains(resPos.Feedback, "Mismatch at position 2") {
		t.Errorf("expected positional mismatch feedback, got %+v", resPos)
	}

	// Empty input
	resEmpty := r.Grade(GradingTuple, "", "(1, 2)")
	if resEmpty.Correct || resEmpty.Score != 0 {
		t.Errorf("expected empty result for empty expected tuple, got %+v", resEmpty)
	}
}

func TestGrader_Ordering_Mismatches(t *testing.T) {
	r := NewRouter()
	defer r.Close()

	// Element count mismatch
	resCount := r.Grade(GradingOrdering, "a b c", "a b")
	if resCount.Correct || resCount.Feedback != "Incorrect number of elements" {
		t.Errorf("expected 'Incorrect number of elements', got %+v", resCount)
	}

	// Order mismatch
	resOrder := r.Grade(GradingOrdering, "a b c", "a c b")
	if resOrder.Correct || resOrder.Feedback != "Incorrect order" {
		t.Errorf("expected 'Incorrect order', got %+v", resOrder)
	}

	// Empty input
	resEmpty := r.Grade(GradingOrdering, "", "a b")
	if resEmpty.Correct || resEmpty.Score != 0 {
		t.Errorf("expected empty result for empty expected ordering, got %+v", resEmpty)
	}
}

func TestGrader_Diagnose_OrdinalEdgeCases(t *testing.T) {
	ordCases := []struct {
		n    int
		want string
	}{
		{1, "1st"},
		{2, "2nd"},
		{3, "3rd"},
		{4, "4th"},
		{11, "11th"},
		{12, "12th"},
		{13, "13th"},
		{21, "21st"},
		{22, "22nd"},
		{23, "23rd"},
		{24, "24th"},
		{111, "111th"},
		{112, "112th"},
		{113, "113th"},
		{121, "121st"},
	}

	for _, c := range ordCases {
		if got := ordinal(c.n); got != c.want {
			t.Errorf("ordinal(%d): got %q, want %q", c.n, got, c.want)
		}
	}
}

func TestGrader_Symbolic_ConstEvalMismatch(t *testing.T) {
	r := NewRouter()
	defer r.Close()

	// Expressions evaluated to differing constants
	resMismatch := r.Grade(GradingExpression, "2^3", "9")
	if resMismatch.Correct || resMismatch.Feedback != "Not equivalent" {
		t.Errorf("expected 'Not equivalent' for 2^3 vs 9, got %+v", resMismatch)
	}

	// Constant arithmetic evaluated to matching value
	resMatch := r.Grade(GradingExpression, "10 - 3 * 2", "4")
	if !resMatch.Correct {
		t.Errorf("expected correct for 10 - 3 * 2 vs 4, got %+v", resMatch)
	}
}

// A single-letter expected answer is a symbol far more often than an option
// label, so the prose fallback must not accept a mathematical expression that
// merely contains that letter. Before this rule, expected "i" accepted "-i",
// and every imaginary-unit question marked the opposite sign correct.
func TestGrader_Choice_SymbolAnswerIsNotMatchedByLetter(t *testing.T) {
	r := NewRouter()
	defer r.Close()

	cases := []struct {
		expected, answer string
		want             bool
	}{
		{"i", "i", true},
		{"i", "-i", false},
		{"-i", "i", false},
		{"i", "i^2", false},
		{"x", "x + 1", false},
		{"B", "The correct answer is B.", true},
		{"C", "The choice is D.", false},
	}
	for _, c := range cases {
		got := r.Grade(GradingMultipleChoice, c.expected, c.answer)
		if got.Correct != c.want {
			t.Errorf("Grade(%q, %q).Correct = %v, want %v", c.expected, c.answer, got.Correct, c.want)
		}
	}
}
