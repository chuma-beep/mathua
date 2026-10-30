package grader

import "testing"

// Three answers that were reported as wrong while being right in form.
func TestReportedListAndTupleAnswers(t *testing.T) {
	r := NewRouter()
	cases := []struct {
		name        string
		gradingType GradingType
		expected    string
		answer      string
		wantCorrect bool
	}{
		{"pascal row with spaces", GradingOrdering, "1,3,3,1", "1, 3, 3, 1", true},
		{"pascal row out of order", GradingOrdering, "1,3,3,1", "1,1,3,3", false},
		{"prime factors in any order", GradingNumeric, "2,3,5", "2,5,3", true},
		{"polar without the degree symbol", GradingTuple, "(1,90°)", "(1,90)", true},
		{"polar with degrees typed out", GradingTuple, "(1,90°)", "(1, 90 degrees)", true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := r.Grade(c.gradingType, c.expected, c.answer).Correct; got != c.wantCorrect {
				t.Errorf("Grade(%s, %q, %q) = %v, want %v", c.gradingType, c.expected, c.answer, got, c.wantCorrect)
			}
		})
	}
}
