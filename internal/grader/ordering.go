package grader

import (
	"regexp"
	"strings"
)

var orderingSep = regexp.MustCompile(`[\s,;|]+`)

type orderingGrader struct{}

func (g *orderingGrader) grade(expected, answer string) Result {
	eParts := splitOrdering(expected)
	aParts := splitOrdering(answer)
	if len(eParts) == 0 || len(aParts) == 0 {
		return Result{Correct: false, Score: 0, Feedback: "Answer must not be empty"}
	}
	if len(eParts) != len(aParts) {
		return Result{Correct: false, Score: 0, Feedback: "Incorrect number of elements"}
	}
	for i := range eParts {
		if !strings.EqualFold(eParts[i], aParts[i]) {
			return Result{Correct: false, Score: 0, Feedback: "Incorrect order"}
		}
	}
	return Result{Correct: true, Score: 1}
}

func splitOrdering(s string) []string {
	parts := orderingSep.Split(strings.TrimSpace(s), -1)
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
