package grader

type Router struct {
	numeric    *numericGrader
	symbolic   *symbolicGrader
	sympy      *sympyGrader
	choice     *choiceGrader
	comparison *comparisonGrader
	ordering   *orderingGrader
	tuple      *tupleGrader
	complex    *complexGrader
}

func NewRouter() *Router {
	return &Router{
		numeric:    &numericGrader{},
		symbolic:   &symbolicGrader{},
		choice:     &choiceGrader{},
		comparison: &comparisonGrader{},
		ordering:   &orderingGrader{},
		tuple:      &tupleGrader{},
		complex:    &complexGrader{},
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
	case GradingTuple:
		return r.tuple.grade(expected, answer)
	case GradingComplex:
		return r.complex.grade(expected, answer)
	default:
		return r.tuple.grade(expected, answer)
	}
}

func (r *Router) Close() {
	if r.sympy != nil {
		r.sympy.close()
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
