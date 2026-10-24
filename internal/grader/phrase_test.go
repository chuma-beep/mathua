package grader

import "testing"

// TestChoice_PhraseWithLostSpaces is the regression for the reported "line
// segment" false miss. MathLive spaces every letter of a word phrase, the
// client rejoins them into one word, and the choice grader compares strings
// exactly, so a correct phrase was marked wrong.
func TestChoice_PhraseWithLostSpaces(t *testing.T) {
	r := NewRouter()
	for _, answer := range []string{"linesegment", "LineSegment", "l i n e s e g m e n t"} {
		res := r.Grade(GradingMultipleChoice, "line segment", answer)
		if !res.Correct {
			t.Errorf("answer %q to expected %q was marked wrong", answer, "line segment")
		}
	}
}

// TestChoice_PhraseDoesNotAbsorbProductsOrDigits pins the guard: only an
// expected answer made of words is compared with its spaces removed.
func TestChoice_PhraseDoesNotAbsorbProductsOrDigits(t *testing.T) {
	r := NewRouter()
	cases := []struct{ expected, answer string }{
		{"line segment", "line"},         // truncated phrase
		{"line segment", "linesegments"}, // a different phrase
		{"line segment", "linesegment2"}, // digits do not belong to a phrase
		{"5 R 3", "5R3"},                 // remainder notation, not a phrase
		{"2pi r", "2pir"},                // a product, not a phrase
	}
	for _, c := range cases {
		if res := r.Grade(GradingMultipleChoice, c.expected, c.answer); res.Correct {
			t.Errorf("expected %q, answer %q was accepted", c.expected, c.answer)
		}
	}
}

// TestChoice_PhraseIgnoresOuterWhitespace covers the trailing-space form that
// broke the earlier fix attempt: the learner's phrase may arrive with padding.
func TestChoice_PhraseIgnoresOuterWhitespace(t *testing.T) {
	r := NewRouter()
	if res := r.Grade(GradingMultipleChoice, "  line segment  ", "linesegment "); !res.Correct {
		t.Error("padded phrase was marked wrong")
	}
}

// TestChoice_PhraseWithPunctuation covers phrases ending with a period or comma
// where MathLive lost whitespace on answer submission.
func TestChoice_PhraseWithPunctuation(t *testing.T) {
	r := NewRouter()
	const expected = "Division by zero is undefined."
	for _, answer := range []string{
		"division by zero is undefined.",
		"division by zero is undefined",
		"divisionbyzeroisundefined.",
		"DivisionByZeroIsUndefined",
	} {
		if res := r.Grade(GradingMultipleChoice, expected, answer); !res.Correct {
			t.Errorf("answer %q to expected %q was marked wrong", answer, expected)
		}
	}
}

// TestChoice_LabeledChoiceVariants covers student inputs such as "A)", "(A)" or "Option A".
func TestChoice_LabeledChoiceVariants(t *testing.T) {
	r := NewRouter()
	for _, answer := range []string{"A", "a", "A)", "a)", "(A)", "(a)", "Option A", "option a", "Choice A", "A."} {
		if res := r.Grade(GradingMultipleChoice, "Option A", answer); !res.Correct {
			t.Errorf("answer %q to expected %q was marked wrong", answer, "Option A")
		}
	}
}
