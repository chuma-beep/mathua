package grader

import (
	"regexp"
	"strings"
)

var singleLetterRe = regexp.MustCompile(`\b[A-Za-z]\b`)

type choiceGrader struct{}

func (g *choiceGrader) grade(expected, answer string) Result {
	e := strings.TrimSpace(expected)
	a := strings.TrimSpace(answer)
	if e == "" || a == "" {
		return Result{Correct: false, Score: 0, Feedback: "Answer must not be empty"}
	}
	if strings.EqualFold(e, a) {
		return Result{Correct: true, Score: 1}
	}

	// If answer is a single letter, look for that letter in expected
	if singleLetterRe.MatchString(a) && len(strings.TrimSpace(a)) == 1 {
		letters := singleLetterRe.FindAllString(e, -1)
		for _, l := range letters {
			if strings.EqualFold(l, a) {
				return Result{Correct: true, Score: 1}
			}
		}
	}

	// If expected is a single letter, look for it in the answer's letters
	if singleLetterRe.MatchString(e) && len(strings.TrimSpace(e)) == 1 {
		letters := singleLetterRe.FindAllString(a, -1)
		for _, l := range letters {
			if strings.EqualFold(l, e) {
				return Result{Correct: true, Score: 1}
			}
		}
	}

	return Result{Correct: false, Score: 0, Feedback: "Incorrect choice"}
}
