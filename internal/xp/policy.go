// Package xp is the single source of truth for Mathua's XP award policy.
//
// Every XP number the system awards — and every number the workload
// estimator predicts — must derive from here. Estimators import this
// package; they must never duplicate its constants.
//
// Audit table (verified against the award paths in internal/engine):
//
//	Path					Base			Miss	Rush
//	practice lesson / word			effort-derived		0	-1 (<2s x2)
//	practice review				half effort, min 1	0	-1 (<2s x2)
//	quiz					effort + assessment	0 + remedial enqueue	-1
//	quiz don't-know				—			forced miss, 0, no rush penalty
//	diagnostic				—			0 XP always (no award path exists)
//
// Bases are an effort-derived approximation, not exact time accounting:
// EffortBase calibrates from the concept's expected seconds, rounded and
// clamped to 1-5 XP. Correct answers earn base x timeMult(0.5-1.25) x
// streakMult(1+0.05*streak). A typical correct answer earns 1-2 XP; a
// five-question lesson totals roughly 5-10 XP. The rush penalty (-1) is
// applied by the award call sites, not here.
package xp

import "math"

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
const RushPenalty = -1

// EffortBase calibrates a base award from expected focused effort:
// threshold seconds scaled by difficulty, rounded and clamped to 1-5 XP.
// 1 XP ≈ 1 minute is an approximation, not an invariant — rounding and
// clamping compress the bottom end (10-25s questions all price at 1 XP).
// A zero/negative threshold falls back to 10s; difficulty clamps to [0,1].
func EffortBase(thresholdSec, difficulty float64) int {
	if thresholdSec <= 0 {
		thresholdSec = 10
	}
	if difficulty < 0 {
		difficulty = 0
	}
	if difficulty > 1 {
		difficulty = 1
	}
	base := int(math.Round(thresholdSec / 60.0 * (0.5 + difficulty)))
	if base < 1 {
		base = 1
	}
	if base > 5 {
		base = 5
	}
	return base
}

// baseForTask resolves a task type to its base. Lessons and word problems
// share the effort calibration (their thresholds already differ);
// reviews cost half (known material, min 1); quizzes add an assessment
// premium of +1 (capped with the base at 5).
func baseForTask(thresholdSec float64, taskType string) int {
	switch taskType {
	case TaskReview:
		return max(1, int(math.Round(float64(EffortBase(thresholdSec, 0.5))/2)))
	case TaskQuiz:
		return min(5, EffortBase(thresholdSec, 0.8)+1)
	case TaskMultistep:
		fallthrough
	case TaskLesson:
		fallthrough
	default:
		return EffortBase(thresholdSec, 0.5)
	}
}

// BaseXP returns the base award for a task type at the default 10s
// threshold. Prefer EffortBase with the concept's real threshold.
func BaseXP(taskType string) int {
	return baseForTask(10, taskType)
}

// TimeMultiplier scales the base by speed relative to the concept's
// accommodated time threshold. Faster earns more, bounded to [0.5, 1.25].
func TimeMultiplier(elapsed, timeThreshold float64) float64 {
	ratio := elapsed / timeThreshold
	if ratio <= 0 {
		ratio = 0.01
	}
	m := 2.0 - ratio
	if m < 0.5 {
		m = 0.5
	}
	if m > 1.25 {
		m = 1.25
	}
	return m
}

// StreakMultiplier rewards consecutive correct answers, capped at 1.5x.
func StreakMultiplier(streak int) float64 {
	if streak < 0 {
		streak = 0
	}
	if streak > 10 {
		streak = 10
	}
	return 1.0 + float64(streak)*0.05
}

// Award computes the XP for one graded attempt, excluding the rush penalty
// (applied by call sites). Incorrect answers earn 0. The timeThreshold
// doubles as the effort calibration input, so awards track expected effort.
func Award(correct bool, elapsed, timeThreshold float64, streak int, taskType string) int {
	if !correct {
		return 0
	}
	return int(float64(baseForTask(timeThreshold, taskType)) * TimeMultiplier(elapsed, timeThreshold) * StreakMultiplier(streak))
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
