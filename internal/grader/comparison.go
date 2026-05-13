package grader

import "strings"

type comparisonGrader struct{}

func (g *comparisonGrader) grade(expected, answer string) Result {
	e := strings.TrimSpace(expected)
	a := strings.TrimSpace(answer)
	if e == "" || a == "" {
		return emptyResult()
	}
	eVal, err := parseComparison(e)
	if err != "" {
		return Result{Correct: false, Score: 0, Feedback: err}
	}
	aVal, aErr := parseComparison(a)
	if aErr != "" {
		return Result{Correct: false, Score: 0, Feedback: aErr}
	}
	if eVal == aVal {
		return Result{Correct: true, Score: 1}
	}
	return Result{Correct: false, Score: 0, Feedback: "Incorrect comparison"}
}

func parseComparison(s string) (string, string) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", "Answer must not be empty"
	}
	switch {
	case strings.HasPrefix(s, ">="):
		return ">=", ""
	case strings.HasPrefix(s, "<="):
		return "<=", ""
	case strings.HasPrefix(s, "!="):
		return "!=", ""
	case strings.HasPrefix(s, "=="):
		return "==", ""
	case strings.HasPrefix(s, ">"):
		return ">", ""
	case strings.HasPrefix(s, "<"):
		return "<", ""
	case strings.HasPrefix(s, "="):
		return "=", ""
	default:
		return "", "Invalid comparison operator"
	}
}
