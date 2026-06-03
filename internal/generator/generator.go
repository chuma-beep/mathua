package generator

import "github.com/chuma-beep/mathua/internal/grader"

type GeneratorContext struct {
	Difficulty float64
	Seed       int64
}

type Problem struct {
	Question    string
	Answer      string
	Explanation string
}

type Generator interface {
	Generate(ctx GeneratorContext) Problem
}

// GradedGenerator is a Generator that can grade its own answers.
// Useful for non-standard formats like "5 R 3" or "3 sqrt(2)".
type GradedGenerator interface {
	Generator
	Grade(expected, userAnswer string) grader.Result
}
