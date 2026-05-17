package grader

import (
	"fmt"
	"regexp"
	"strings"
)

var tupleRe = regexp.MustCompile(`[();,]+`)

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
		if e[i] != a[i] {
			return Result{Correct: false, Score: 0, Feedback: fmt.Sprintf("Mismatch at position %d", i+1)}
		}
	}
	return Result{Correct: true, Score: 1}
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
