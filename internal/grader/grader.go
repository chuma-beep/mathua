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
}
