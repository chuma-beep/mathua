package grader

type GradingType string

const (
	GradingNumeric        GradingType = "numeric"
	GradingPolynomial     GradingType = "polynomial"
	GradingExpression     GradingType = "expression"
	GradingMultipleChoice GradingType = "multiple_choice"
	GradingComparison     GradingType = "comparison"
	GradingOrdering       GradingType = "ordering"
)

type Result struct {
	Correct  bool
	Score    float64
	Feedback string
}
