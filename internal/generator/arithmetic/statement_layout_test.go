package arithmetic

import (
	"strconv"
	"strings"
	"testing"
)

func TestCountArithVerdictsRendersOneStatementPerLine(t *testing.T) {
	pool := []arithVerdict{
		{label: "one", yes: true}, {label: "two", yes: true}, {label: "three", yes: false}, {label: "four", yes: false},
		{label: "five", yes: true}, {label: "six", yes: true}, {label: "seven", yes: false}, {label: "eight", yes: false},
	}
	p := countArithVerdicts("How many statements are true?", pool)
	if !strings.Contains(p.Question, "\n(b) ") || strings.Count(p.Question, "\n(") != 7 {
		t.Fatalf("statements are not newline separated: %q", p.Question)
	}
	if _, err := strconv.Atoi(p.Answer); err != nil {
		t.Fatalf("answer %q is not an integer", p.Answer)
	}
}
