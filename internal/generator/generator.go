package generator

import "github.com/chuma-beep/mathua/internal/grader"

type Problem struct {
	Question    string
	Answer      string
	Explanation string
}

type Generator interface {
	Generate(difficulty float64) Problem
}

// GradedGenerator is a Generator that also knows how to grade its own answers.
// This lets generators with non-standard answer formats (e.g. "5 R 3", "3 sqrt(2)")
// provide their own grading logic instead of relying on the generic numeric grader.
type GradedGenerator interface {
	Generator
	Grade(expected, userAnswer string) grader.Result
}
