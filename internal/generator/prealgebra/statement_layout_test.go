package prealgebra

import (
	"strconv"
	"strings"
	"testing"
)

func TestRationalCountProblemRendersOneStatementPerLine(t *testing.T) {
	p := rationalCountProblem()
	if !strings.Contains(p.Question, "\n(b)") || strings.Count(p.Question, "\n(") != 7 {
		t.Fatalf("statements are not newline separated: %q", p.Question)
	}
	if _, err := strconv.Atoi(p.Answer); err != nil {
		t.Fatalf("answer %q is not an integer", p.Answer)
	}
}
