// Package xp is the single source of truth for Mathua's XP award policy.
//
// Every XP number the system awards — and every number the workload estimator predicts —
// must derive from here. Estimators import this package; they must never duplicate its
// constants.
//
// XP means productive learning effort, priced at roughly one XP per expected minute of
// focused work. The base is the concept's own effort, so a 6-second arithmetic fact and a
// two-minute proof are not worth the same amount; everything after the base is a small
// adjustment, never the reason a question is worth earning.
//
// Audit table (verified against the award paths in internal/engine):
//
//	Path				Base			Miss	Rush
//	practice lesson / word	effort-derived		0	-1 (see RushAfter)
//	practice review		effort-derived		0	-1
//	quiz			effort-derived		0 + remedial enqueue	-1
//	quiz don't-know		—			forced miss, 0, no rush penalty
//	diagnostic		—			0 XP always (no award path exists)
//
// Three deliberate properties:
//
//   - **The task type does not change the base.** Reviews and quizzes used to cost half and
//     cost +1 respectively. Both were incentives to route around the main loop — review XP
//     being cheap made it unrewarding to keep knowledge alive, and the quiz premium made a
//     quiz worth more per question than the teaching it was meant to check. All three now
//     earn the same effort-derived base.
//   - **Speed is not rewarded.** The old TimeMultiplier was `2 - elapsed/threshold` clamped
//     to [0.5, 1.25]: a fast answer earned up to 25% more and a slow one 50% less. That pays
//     for guessing, and it punishes careful work. Elapsed time now answers one question —
//     did the learner appear to have worked normally — and nothing else.
//   - **There is no streak bonus.** StreakMultiplier reached 1.5x, above anything the product
//     intends, so a learner on a long run out-earned the same work done carefully. Combined
//     with the old time multiplier a single answer could pay 1.875x its base, which is what
//     let the workload estimator drift out of step with reality (see planning.EstimateWorkload).
//
// Correct answers earn base x PerformanceMultiplier, which is 1.1 within the concept's own
// time threshold and 1.0 outside it. A five-question lesson therefore totals ~5-6 XP.
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

// Daily XP goal policy. Owned here because the server, the estimator and the client all need
// the same numbers, and they had drifted: the shipped default was 10 while the planner's
// fallbacks said 30, so a learner with no stored goal planned against 3x the XP/day they
// were actually banking.
const (
	// DefaultDailyGoal is the XP/day a learner is assumed to target when they have not
	// chosen one. 30 XP/day is roughly 18 minutes of expected work, not 30 — see
	// EffortBase on why the one-XP-per-minute invariant is an approximation.
	DefaultDailyGoal = 30
	// MinDailyGoal and MaxDailyGoal bound the choice. The client's lib/plan.ts pins the
	// same three numbers from the other side.
	MinDailyGoal = 10
	MaxDailyGoal = 100
)

// QuizGateXP is how much learning XP earns a mastery check. It is deliberately NOT derived
// from DefaultDailyGoal: the gate is an interval between assessments, and tying it to the
// daily target would mean a learner who raised their goal started being tested more often.
// ADR-020 set it to 50 alongside a 10 XP/day default; the daily default is now 30, so the
// gate now arrives every ~1.7 days rather than ~10. That is a live product consequence, not
// an oversight — see TestDailyGoalImpliedCadence, which reports it.
const QuizGateXP = 50

// RushPenalty is the negative XP assigned for rushing/guessing.
const RushPenalty = -1

// RushAfter is how many sub-2-second wrong answers earn the penalty. It used to differ
// between hosts — the practice path counted its own rushCount, the study/quiz path required
// two misses on the concept as well — so the same behaviour was penalised on one route and
// not the other.
const RushAfter = 2

// RushSeconds is the elapsed time below which a wrong answer counts as a guess.
const RushSeconds = 2.0

// EffortBase calibrates a base award from expected focused effort:
// threshold seconds scaled by difficulty, rounded and clamped to 1-5 XP.
//
// 1 XP ~ 1 minute is an approximation, and the clamp is why: rounding plus the 1-XP floor
// compresses the bottom end hard. Measured across the corpus at the difficulty the engine
// actually serves with no learner information (0.55), **481 of 657 concepts price at 1 XP**
// and only 174 at 2 XP, so the corpus's 509 expected minutes are priced at 835 XP — 1.64 XP
// per real minute rather than 1. The worst are 6-second concepts such as arith.mult.tables
// at a 10x over-price. That is a deliberate floor: sub-minute questions still count as work,
// and a fractional XP economy would make every answer round to nothing.
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

// BaseDifficulty is the difficulty a task type's effort is calibrated at.
//
// It does not vary by task type any more: reviews and quizzes earn the same effort base as
// lessons (see the package comment). The parameter is kept because it names the task in the
// audit trail and is where a task-specific rule would belong, not because it currently
// changes the number.
func BaseDifficulty(string) float64 { return 0.5 }

// BaseXP returns the base award for a task type at the default 10s threshold.
// Prefer EffortBase with the concept's real threshold.
func BaseXP(taskType string) int {
	return EffortBase(10, BaseDifficulty(taskType))
}

// PerformanceMultiplier is the only multiplier on a correct answer.
//
// **It cannot move an integer award**, and that is worth stating rather than hiding: the
// base is clamped to 1-5 XP and Award truncates, so int(base x 1.1) == base for every base in
// range (1->1, 2->2, ... 5->5). The tier is real in the model — ExpectedXP returns it as a
// float and the estimator reads that — but on the integer award it is absorbed by rounding.
// Making it pay would mean either raising the documented 1-5 ceiling (a 5 XP question at 5.5
// rounds to 6) or dropping the clamp, and both change the effort calibration the estimator
// depends on. So the honest reading of the product's 1.1x tier is "within threshold is the
// good case", expressed here so the policy is stated once.
//
// Two tiers, from the product's own table: 1.1 for a correct answer inside the concept's
// time threshold, 1.0 for a correct answer outside it. Elapsed time is read as "did the
// learner appear to have worked normally", which is the honest use of it — the previous
// formula made a fast answer worth 25% more and a slow one 50% less, so it paid for
// guessing and penalised care.
//
// Rushing is handled separately and only when the answer is wrong (RushSeconds/RushAfter at
// the call sites), because a fast *correct* answer is a good answer.
func PerformanceMultiplier(elapsed, timeThreshold float64) float64 {
	if timeThreshold > 0 && elapsed <= timeThreshold {
		return 1.1
	}
	return 1.0
}

// Award computes the XP for one graded attempt, excluding the rush penalty (applied by call
// sites). Incorrect answers earn 0. The timeThreshold doubles as the effort calibration
// input, so awards track expected effort rather than a flat per-question number.
func Award(correct bool, elapsed, timeThreshold float64, streak int, taskType string) int {
	if !correct {
		return 0
	}
	// streak and taskType no longer move the award; they are retained so call sites and the
	// estimator keep one signature and the audit trail keeps naming the task.
	_ = streak
	base := EffortBase(timeThreshold, BaseDifficulty(taskType))
	return int(float64(base) * PerformanceMultiplier(elapsed, timeThreshold))
}

// ExpectedXP is the estimator's per-attempt expectation: award value scaled by probability of
// correctness. Never use face base value as an expectation.
func ExpectedXP(pCorrect, elapsed, timeThreshold float64, streak int, taskType string) float64 {
	if pCorrect < 0 {
		pCorrect = 0
	}
	if pCorrect > 1 {
		pCorrect = 1
	}
	return float64(Award(true, elapsed, timeThreshold, streak, taskType)) * pCorrect
}
