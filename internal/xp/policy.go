// Package xp is the single source of truth for Mathua's XP award policy.
//
// Every XP number the system awards — and every number the workload
// estimator predicts — must derive from here. Estimators import this
// package; they must never duplicate its constants.
//
// Audit table (verified against the award paths in internal/engine):
//
//	Path					Base	Miss	Rush
//	practice lesson (TaskLesson)		10	0	-5 (<2s x2)
//	practice review (TaskReview)		 5	0	-5 (<2s x2)
//	practice/study word (TaskMultistep)	15	0	-5 (<2s x2)
//	quiz (TaskQuiz)				20	0 + remedial enqueue	-5
//	quiz don't-know				—	forced miss, 0, no rush penalty
//	diagnostic				—	0 XP always (no award path exists)
//
// Correct answers earn base x timeMult(0.5-1.5) x streakMult(1+0.1*streak).
// The rush penalty (-5) is applied by the award call sites, not here:
// see RushPenalty and the engine's submit paths.
package xp

// Task types mirror internal/engine's TaskLesson/TaskReview/TaskMultistep/
// TaskQuiz strings. They are repeated here (not imported) because engine
// imports this package; an import cycle forbids sharing.
const (
	TaskLesson    = "lesson"
	TaskReview    = "review"
	TaskMultistep = "multistep"
	TaskQuiz      = "quiz"
)

// RushPenalty is the negative XP assigned for rushing/guessing.
const RushPenalty = -5

// BaseXP returns the base award for a task type (10/5/15/20).
func BaseXP(taskType string) int {
	switch taskType {
	case TaskReview:
		return 5
	case TaskMultistep:
		return 15
	case TaskQuiz:
		return 20
	case TaskLesson:
		fallthrough
	default:
		return 10
	}
}

// TimeMultiplier scales the base by speed relative to the concept's
// accommodated time threshold. Faster earns more, bounded to [0.5, 1.5].
func TimeMultiplier(elapsed, timeThreshold float64) float64 {
	ratio := elapsed / timeThreshold
	if ratio <= 0 {
		ratio = 0.01
	}
	m := 2.0 - ratio
	if m < 0.5 {
		m = 0.5
	}
	if m > 1.5 {
		m = 1.5
	}
	return m
}

// StreakMultiplier rewards consecutive correct answers, capped at 2x.
func StreakMultiplier(streak int) float64 {
	if streak < 0 {
		streak = 0
	}
	if streak > 10 {
		streak = 10
	}
	return 1.0 + float64(streak)*0.1
}

// Award computes the XP for one graded attempt, excluding the rush penalty
// (applied by call sites). Incorrect answers earn 0.
func Award(correct bool, elapsed, timeThreshold float64, streak int, taskType string) int {
	if !correct {
		return 0
	}
	return int(float64(BaseXP(taskType)) * TimeMultiplier(elapsed, timeThreshold) * StreakMultiplier(streak))
}

// ExpectedXP is the estimator's per-attempt expectation: award value scaled
// by probability of correctness. Never use face base value as an expectation.
func ExpectedXP(pCorrect, elapsed, timeThreshold float64, streak int, taskType string) float64 {
	if pCorrect < 0 {
		pCorrect = 0
	}
	if pCorrect > 1 {
		pCorrect = 1
	}
	return float64(Award(true, elapsed, timeThreshold, streak, taskType)) * pCorrect
}
