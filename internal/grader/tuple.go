package grader

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	tupleRe = regexp.MustCompile(`[();,]+`)
	// A degree symbol is hard to type from a math keyboard, and a learner who
	// writes 90 instead of 90° is not making a mathematical error. The same
	// applies to "deg" and "degrees", which are how the value is usually typed
	// when the symbol is not available.
	tupleDegreeRe = regexp.MustCompile(`(?i)\s*(?:°|deg(?:rees?)?)\s*$`)
)

type tupleGrader struct{}

func (g *tupleGrader) grade(expected, answer string) Result {
	e := normaliseTuple(expected)
	a := normaliseTuple(answer)
	if len(e) == 0 || len(a) == 0 {
		return emptyResult()
	}
	if len(e) != len(a) {
		return Result{Correct: false, Score: 0, Feedback: fmt.Sprintf("Expected %s, got %s", strings.Join(e, ","), strings.Join(a, ","))}
	}
	for i := range e {
		if !tuplePartEqual(e[i], a[i]) {
			return Result{Correct: false, Score: 0, Feedback: fmt.Sprintf("Mismatch at position %d", i+1)}
		}
	}
	return Result{Correct: true, Score: 1}
}

// tuplePartEqual ignores a degree marker, because it is a unit rather than
// part of the value, and compares everything else exactly. The exactness
// matters: a tuple such as (4000,300,5) lists three separate quantities, so
// accepting 4000.0 for 4000 would quietly grade a different answer as correct.
func tuplePartEqual(expected, answer string) bool {
	expected = strings.TrimSpace(tupleDegreeRe.ReplaceAllString(expected, ""))
	answer = strings.TrimSpace(tupleDegreeRe.ReplaceAllString(answer, ""))
	return strings.EqualFold(expected, answer)
}

func normaliseTuple(s string) []string {
	s = strings.TrimSpace(s)
	s = tupleRe.ReplaceAllString(s, ",")
	parts := strings.Split(s, ",")
	var result []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
