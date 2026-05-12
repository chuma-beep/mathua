package grader

import (
	"testing"
)

// ---------------------------------------------------------------------------
// Numeric grader
// ---------------------------------------------------------------------------

func TestNumeric_Grade_IntegerMatch(t *testing.T) {
	r := NewRouter()
	for _, tc := range []struct{ expected, answer string }{
		{"42", "42"},
		{"-7", "-7"},
		{"0", "0"},
	} {
		res := r.Grade(GradingNumeric, tc.expected, tc.answer)
		if !res.Correct {
			t.Errorf("expected %q == %q to be correct", tc.expected, tc.answer)
		}
	}
}

func TestNumeric_Grade_ZeroVariants(t *testing.T) {
	r := NewRouter()
	for _, tc := range []struct{ expected, answer string }{
		{"0", "-0"},
		{"0", "+0"},
		{"+0", "0"},
		{"-0", "0"},
	} {
		res := r.Grade(GradingNumeric, tc.expected, tc.answer)
		if !res.Correct {
			t.Errorf("expected %q == %q to be correct", tc.expected, tc.answer)
		}
	}
}

func TestNumeric_Grade_Decimals(t *testing.T) {
	r := NewRouter()
	for _, tc := range []struct{ expected, answer string }{
		{"3.14", "3.14"},
		{"0.5", ".5"},
		{"0.5", "0.50"},
		{"-0.75", "-.75"},
	} {
		res := r.Grade(GradingNumeric, tc.expected, tc.answer)
		if !res.Correct {
			t.Errorf("expected %q == %q to be correct", tc.expected, tc.answer)
		}
	}
}

func TestNumeric_Grade_FloatTolerance(t *testing.T) {
	r := NewRouter()
	res := r.Grade(GradingNumeric, "0.3", "0.30000000001")
	if !res.Correct {
		t.Error("expected 0.3 ≈ 0.30000000001 to be correct")
	}
}

func TestNumeric_Grade_Fractions(t *testing.T) {
	r := NewRouter()
	for _, tc := range []struct{ expected, answer string }{
		{"3/4", "3/4"},
		{"3/4", "6/8"},
		{"3/4", "0.75"},
		{"-1/2", "-0.5"},
		{"1/3", "0.3333333333"},
	} {
		res := r.Grade(GradingNumeric, tc.expected, tc.answer)
		if !res.Correct {
			t.Errorf("expected %q == %q to be correct", tc.expected, tc.answer)
		}
	}
}

func TestNumeric_Grade_MixedNumbers(t *testing.T) {
	r := NewRouter()
	res := r.Grade(GradingNumeric, "1 1/2", "1.5")
	if !res.Correct {
		t.Error("expected 1 1/2 == 1.5 to be correct")
	}
	res = r.Grade(GradingNumeric, "1 1/2", "3/2")
	if !res.Correct {
		t.Error("expected 1 1/2 == 3/2 to be correct")
	}
}

func TestNumeric_Grade_ScientificNotation(t *testing.T) {
	r := NewRouter()
	for _, tc := range []struct{ expected, answer string }{
		{"100", "1e2"},
		{"100", "1E2"},
		{"0.001", "1e-3"},
		{"150", "1.5e2"},
	} {
		res := r.Grade(GradingNumeric, tc.expected, tc.answer)
		if !res.Correct {
			t.Errorf("expected %q == %q to be correct", tc.expected, tc.answer)
		}
	}
}

func TestNumeric_Grade_WrongAnswer(t *testing.T) {
	r := NewRouter()
	res := r.Grade(GradingNumeric, "5", "6")
	if res.Correct {
		t.Error("expected 5 != 6 to be incorrect")
	}
}

func TestNumeric_Grade_EmptyAnswer(t *testing.T) {
	r := NewRouter()
	res := r.Grade(GradingNumeric, "5", "")
	if res.Correct {
		t.Error("expected empty answer to be incorrect")
	}
	if res.Feedback == "" {
		t.Error("expected feedback for empty answer")
	}
}

func TestNumeric_Grade_EmptyExpected(t *testing.T) {
	r := NewRouter()
	res := r.Grade(GradingNumeric, "", "5")
	if res.Correct {
		t.Error("expected empty expected to be incorrect")
	}
}

func TestNumeric_Grade_Whitespace(t *testing.T) {
	r := NewRouter()
	res := r.Grade(GradingNumeric, "  42  ", "42")
	if !res.Correct {
		t.Error("expected whitespace-normalised match")
	}
}

func TestNumeric_Grade_Commas(t *testing.T) {
	r := NewRouter()
	res := r.Grade(GradingNumeric, "1,000", "1000")
	if !res.Correct {
		t.Error("expected comma-normalised match")
	}
}

func TestNumeric_Grade_MinusZeroVsZero(t *testing.T) {
	r := NewRouter()
	res := r.Grade(GradingNumeric, "-0", "0")
	if !res.Correct {
		t.Error("expected -0 == 0 to be correct")
	}
}

// ---------------------------------------------------------------------------
// Multiple choice grader
// ---------------------------------------------------------------------------

func TestChoice_Grade_ExactMatch(t *testing.T) {
	r := NewRouter()
	res := r.Grade(GradingMultipleChoice, "Option B", "Option B")
	if !res.Correct {
		t.Error("expected exact match")
	}
}

func TestChoice_Grade_CaseInsensitive(t *testing.T) {
	r := NewRouter()
	res := r.Grade(GradingMultipleChoice, "Option B", "option b")
	if !res.Correct {
		t.Error("expected case-insensitive match")
	}
}

func TestChoice_Grade_LetterOnly(t *testing.T) {
	r := NewRouter()
	res := r.Grade(GradingMultipleChoice, "Option B", "B")
	if !res.Correct {
		t.Error("expected single-letter answer B to match option starting with B")
	}
	res = r.Grade(GradingMultipleChoice, "option b", "b")
	if !res.Correct {
		t.Error("expected single-letter b to match b")
	}
}

func TestChoice_Grade_LetterOnly_ExpectedLetter(t *testing.T) {
	r := NewRouter()
	res := r.Grade(GradingMultipleChoice, "B", "Option B")
	if !res.Correct {
		t.Error("expected single-letter expected B to match full Option B")
	}
}

func TestChoice_Grade_WrongLetter(t *testing.T) {
	r := NewRouter()
	res := r.Grade(GradingMultipleChoice, "Option B", "A")
	if res.Correct {
		t.Error("expected wrong letter to be incorrect")
	}
}

func TestChoice_Grade_Empty(t *testing.T) {
	r := NewRouter()
	res := r.Grade(GradingMultipleChoice, "Option B", "")
	if res.Correct {
		t.Error("expected empty answer to be incorrect")
	}
}

// ---------------------------------------------------------------------------
// Comparison grader
// ---------------------------------------------------------------------------

func TestComparison_Grade_Correct(t *testing.T) {
	r := NewRouter()
	for _, tc := range []struct{ expected, answer string }{
		{">", ">"},
		{"<", "<"},
		{"=", "="},
		{">=", ">="},
		{"<=", "<="},
		{"!=", "!="},
		{"==", "=="},
	} {
		res := r.Grade(GradingComparison, tc.expected, tc.answer)
		if !res.Correct {
			t.Errorf("expected %q == %q to be correct", tc.expected, tc.answer)
		}
	}
}

func TestComparison_Grade_Wrong(t *testing.T) {
	r := NewRouter()
	res := r.Grade(GradingComparison, ">", "<")
	if res.Correct {
		t.Error("expected wrong comparison to be incorrect")
	}
}

func TestComparison_Grade_Empty(t *testing.T) {
	r := NewRouter()
	res := r.Grade(GradingComparison, ">", "")
	if res.Correct {
		t.Error("expected empty answer to be incorrect")
	}
}

func TestComparison_Grade_InvalidOperator(t *testing.T) {
	r := NewRouter()
	res := r.Grade(GradingComparison, ">", "?")
	if res.Correct {
		t.Error("expected invalid operator to be incorrect")
	}
}

// ---------------------------------------------------------------------------
// Ordering grader
// ---------------------------------------------------------------------------

func TestOrdering_Grade_Correct(t *testing.T) {
	r := NewRouter()
	for _, tc := range []struct{ expected, answer string }{
		{"1,2,3", "1, 2, 3"},
		{"1 2 3", "1,2,3"},
		{"a;b;c", "a,b,c"},
		{"x|y|z", "x,y,z"},
	} {
		res := r.Grade(GradingOrdering, tc.expected, tc.answer)
		if !res.Correct {
			t.Errorf("expected %q == %q to be correct", tc.expected, tc.answer)
		}
	}
}

func TestOrdering_Grade_WrongOrder(t *testing.T) {
	r := NewRouter()
	res := r.Grade(GradingOrdering, "1,2,3", "3,2,1")
	if res.Correct {
		t.Error("expected wrong order to be incorrect")
	}
}

func TestOrdering_Grade_WrongCount(t *testing.T) {
	r := NewRouter()
	res := r.Grade(GradingOrdering, "1,2,3", "1,2")
	if res.Correct {
		t.Error("expected wrong element count to be incorrect")
	}
}

func TestOrdering_Grade_Empty(t *testing.T) {
	r := NewRouter()
	res := r.Grade(GradingOrdering, "1,2,3", "")
	if res.Correct {
		t.Error("expected empty answer to be incorrect")
	}
}

// ---------------------------------------------------------------------------
// Symbolic grader
// ---------------------------------------------------------------------------

func TestSymbolic_Grade_ExactMatch(t *testing.T) {
	r := NewRouter()
	res := r.Grade(GradingPolynomial, "x^2 + 2x + 1", "x^2+2x+1")
	if !res.Correct {
		t.Error("expected normalised symbolic match")
	}
}

func TestSymbolic_Grade_ExponentNotation(t *testing.T) {
	r := NewRouter()
	res := r.Grade(GradingPolynomial, "x**2 + 1", "x^2 + 1")
	if !res.Correct {
		t.Error("expected ** and ^ to be equivalent")
	}
}

func TestSymbolic_Grade_Wrong(t *testing.T) {
	r := NewRouter()
	res := r.Grade(GradingPolynomial, "x^2+1", "x^2+2")
	if res.Correct {
		t.Error("expected mismatch to be incorrect")
	}
}

func TestSymbolic_Grade_Empty(t *testing.T) {
	r := NewRouter()
	res := r.Grade(GradingPolynomial, "x^2+1", "")
	if res.Correct {
		t.Error("expected empty answer to be incorrect")
	}
}

// ---------------------------------------------------------------------------
// Router dispatch
// ---------------------------------------------------------------------------

func TestRouter_UnknownType(t *testing.T) {
	r := NewRouter()
	res := r.Grade(GradingType("bogus"), "1", "1")
	if res.Correct {
		t.Error("expected unknown grading type to be incorrect")
	}
	if res.Feedback != "Unknown grading type" {
		t.Errorf("expected 'Unknown grading type' feedback, got %q", res.Feedback)
	}
}

func TestRouter_Score(t *testing.T) {
	r := NewRouter()
	res := r.Grade(GradingNumeric, "42", "42")
	if res.Score != 1 {
		t.Errorf("expected score 1 for correct answer, got %f", res.Score)
	}
	res = r.Grade(GradingNumeric, "42", "99")
	if res.Score != 0 {
		t.Errorf("expected score 0 for wrong answer, got %f", res.Score)
	}
}

func TestRouter_FeedbackCorrect(t *testing.T) {
	r := NewRouter()
	res := r.Grade(GradingNumeric, "42", "42")
	if res.Feedback != "" {
		t.Errorf("expected no feedback for correct answer, got %q", res.Feedback)
	}
}

func TestRouter_FeedbackIncorrect(t *testing.T) {
	r := NewRouter()
	res := r.Grade(GradingNumeric, "42", "99")
	if res.Feedback == "" {
		t.Error("expected feedback for incorrect answer")
	}
}
