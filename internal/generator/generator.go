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
	// Facts are the named values this instance was built from — the operands
	// and the derived quantities, already computed and already reduced. A
	// corpus solution schema interpolates them into authored prose; it never
	// recomputes them, so an explanation cannot contradict the answer.
	//
	// Optional: a schema with no facts behind it falls back to Explanation.
	// Keys are snake_case and stable; renaming one is a content change and the
	// contract test will say so.
	Facts map[string]string
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
