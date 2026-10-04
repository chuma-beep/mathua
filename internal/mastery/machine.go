package mastery

import "math"

// Signals that drive a mastery state transition.
type TransitionCtx struct {
	Streak            int     // current consecutive correct streak
	RequiredStreak    int     // mastery threshold for streak
	AvgResponseTime   float64 // current average response time in seconds
	ResponseThreshold float64 // mastery threshold for response time
}

// Machine steps from one mastery status to the next based on performance.
//
// Rules:
//
//	UNSEEN    → LEARNING   when streak >= required AND time <= threshold
//	LEARNING  → PRACTICING when streak >= required AND time <= threshold
//	PRACTICING → MASTERED  when streak >= required AND time <= threshold
//	MASTERED  → MASTERED   (decay detection is the caller's job)
type Machine struct{}

// Next computes the next status.
func (m *Machine) Next(current Status, ctx TransitionCtx) Status {
	switch current {
	case StatusUnseen, StatusLearning, StatusPracticing:
		if ctx.Streak >= ctx.RequiredStreak && ctx.AvgResponseTime <= ctx.ResponseThreshold {
			return nextAdvance(current)
		}
		return current
	case StatusMastered:
		return StatusMastered
	default:
		return current
	}
}

func nextAdvance(current Status) Status {
	switch current {
	case StatusUnseen:
		return StatusLearning
	case StatusLearning:
		return StatusPracticing
	case StatusPracticing:
		return StatusMastered
	default:
		return current
	}
}

// SM2Quality turns response speed into an SM-2 score (0–5).
// timeRatio = avgResponseTime / responseThreshold.
//
//	ratio <= 0.5 → 5 (fast)
//	ratio <= 1.0 → 4 (expected speed)
//	ratio >  1.0 → 3 (slow but correct)
func SM2Quality(streakMet bool, timeRatio float64) int {
	if !streakMet {
		return 0
	}
	switch {
	case timeRatio <= 0.5:
		return 5
	case timeRatio <= 1.0:
		return 4
	default:
		return 3
	}
}

// MasteryPct is the graph's bar length, so it has to mean attainment and never go
// backwards.
//
// It used to be computed as `streak / required_streak`, which is wrong for this purpose
// because the engine resets the streak to 1 on every tier advance (Engine.submitAnswerWithTask).
// A concept therefore sat at 100% and then dropped to ~10%, three times on the way up —
// the learner watched their progress fall while doing everything right. Any streak-derived
// ratio inherits that sawtooth.
//
// Attainment is measured from the tier, not the counter: each of the three transitions is
// worth a third, and progress within the current tier fills the rest. That is monotonic
// across an advance, because the tier's baseline jump is larger than the within-tier reset
// it replaces — LEARNING at 9/10 of its streak reads 63%, and the same instant as
// PRACTICING at 1/10 reads 70%.
//
// Uses only columns that already exist (status, streak, threshold), so it needs no
// migration and no new persisted state.
func MasteryPct(status Status, streak, requiredStreak int) float64 {
	tiers := 0.0
	switch status {
	case StatusLearning:
		tiers = 1
	case StatusPracticing:
		tiers = 2
	case StatusMastered, "DECAYING":
		// Decay is not attainment: the competence was demonstrated and what is
		// outstanding is a retrieval check. The bar stays full and the state is carried
		// by colour alone, so a learner never sees decay as loss of learning.
		tiers = 3
	}
	withinTier := 1.0
	if requiredStreak > 0 {
		withinTier = math.Min(1.0, float64(streak)/float64(requiredStreak))
	} else if streak <= 0 {
		withinTier = 0
	}
	return math.Min(1.0, (tiers+withinTier)/3.0)
}
