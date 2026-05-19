package mastery

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
