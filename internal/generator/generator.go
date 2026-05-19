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

// GradedGenerator is a Generator that can grade its own answers.
// Useful for non-standard formats like "5 R 3" or "3 sqrt(2)".
type GradedGenerator interface {
	Generator
	Grade(expected, userAnswer string) grader.Result
}
