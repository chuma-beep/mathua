package grader

type GradingType string

const (
	GradingNumeric        GradingType = "numeric"
	GradingPolynomial     GradingType = "polynomial"
	GradingExpression     GradingType = "expression"
	GradingSymbolic       GradingType = "symbolic"
	GradingMultipleChoice GradingType = "multiple_choice"
	GradingComparison     GradingType = "comparison"
	GradingOrdering       GradingType = "ordering"
	GradingTuple          GradingType = "tuple"
	GradingComplex        GradingType = "complex"
)

type Result struct {
	Correct  bool
	Score    float64
	Feedback string
	// Unavailable is set when the grader could not be run at all (e.g. the
	// SymPy runtime is missing or errored). Callers must NOT record this as a
	// student miss: "we couldn't evaluate" is not "you were wrong".
	Unavailable bool
	// Diagnosis names the mistake in one sentence when the submitted answer
	// makes it determinable, and is empty otherwise. It is a description of
	// what happened, produced only after grading, and never affects Correct.
	Diagnosis string
}
