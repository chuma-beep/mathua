package grader

import "strings"

type symbolicGrader struct{}

func (g *symbolicGrader) grade(expected, answer string) Result {
	e := normaliseSymbolic(expected)
	a := normaliseSymbolic(answer)
	if e == "" || a == "" {
		return Result{Correct: false, Score: 0, Feedback: "Answer must not be empty"}
	}
	if e == a {
		return Result{Correct: true, Score: 1}
	}
	return Result{Correct: false, Score: 0, Feedback: "Incorrect"}
}

func normaliseSymbolic(s string) string {
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "**", "^")
	return strings.TrimSpace(s)
}
