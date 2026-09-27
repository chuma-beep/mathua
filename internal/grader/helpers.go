package grader

const ErrEmptyAnswer = "Answer must not be empty"

func emptyResult() Result {
	return Result{Correct: false, Score: 0, Feedback: ErrEmptyAnswer}
}
