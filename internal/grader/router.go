package grader

type Router struct {
	numeric    *numericGrader
	symbolic   *symbolicGrader
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
		return gradeSymPy(expected, answer)
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

func (r *Router) Close() {}
