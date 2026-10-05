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

	// Difficulty the question was generated at, 0.3-1.0. Nil means unknown — a problem
	// that did not come from a difficulty-parameterised generator, such as a curated row
	// read straight out of the questions table.
	//
	// This field exists because difficulty used to go *into* GeneratorContext and stop
	// there. Every surface that generated a question and later graded it had to thread the
	// value back out by hand, and two of the three did not: the study and quiz paths
	// recorded `attempts.difficulty = NULL`. Mastery evidence therefore depended on which
	// screen the learner happened to use, which is not a property a learner should be able
	// to change by tapping a different button. Carrying it on the Problem makes generation
	// and grading agree by construction — the registry fills this in, so no caller can
	// generate a question and forget what it was generated at.
	Difficulty *float64
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
