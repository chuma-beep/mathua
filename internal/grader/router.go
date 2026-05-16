package grader

type Router struct {
	numeric    *numericGrader
	symbolic   *symbolicGrader
	sympy      *sympyGrader
	choice     *choiceGrader
	comparison *comparisonGrader
	ordering   *orderingGrader
}

func NewRouter() *Router {
	return &Router{
		numeric:    &numericGrader{},
		symbolic:   &symbolicGrader{},
		choice:     &choiceGrader{},
		comparison: &comparisonGrader{},
		ordering:   &orderingGrader{},
	}
}

func (r *Router) Grade(gradingType GradingType, expected, answer string) Result {
	switch gradingType {
	case GradingNumeric:
		return r.numeric.grade(expected, answer)
	case GradingPolynomial, GradingExpression:
		return r.sympyGrade(expected, answer)
	case GradingMultipleChoice:
		return r.choice.grade(expected, answer)
	case GradingComparison:
		return r.comparison.grade(expected, answer)
	case GradingOrdering:
		return r.ordering.grade(expected, answer)
	default:
		return Result{Correct: false, Score: 0, Feedback: "Unknown grading type"}
	}
}

func (r *Router) sympyGrade(expected, answer string) Result {
	if r.sympy == nil {
		var err error
		r.sympy, err = newSympyGrader()
		if err != nil {
			// Fall back to string-based symbolic grader
			return r.symbolic.grade(expected, answer)
		}
	}
	return r.sympy.grade(expected, answer)
}
