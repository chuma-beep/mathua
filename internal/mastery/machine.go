package mastery

// TransitionCtx carries the signals used to decide the next mastery state.
type TransitionCtx struct {
	Streak            int     // current consecutive correct streak
	RequiredStreak    int     // mastery threshold for streak
	AvgResponseTime   float64 // current average response time in seconds
	ResponseThreshold float64 // mastery threshold for response time
}

// Machine advances the mastery state machine.
// It returns the new status given the current status and the latest performance context.
//
// Transition rules:
//
//	UNSEEN    → LEARNING   when streak >= required AND time <= threshold
//	LEARNING  → PRACTICING when streak >= required AND time <= threshold
//	PRACTICING → MASTERED  when streak >= required AND time <= threshold
//	MASTERED  → MASTERED   (caller is responsible for decay detection)
//
// If either condition is not met, the status stays unchanged.
type Machine struct{}

// Next returns the new mastery status.
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

// SM2Quality maps response speed to an SM-2 quality score (0–5).
// The time ratio is avgResponseTime / responseThreshold.
//
//	ratio <= 0.5 → quality 5 (perfect, fast)
//	ratio <= 1.0 → quality 4 (correct, expected speed)
//	ratio >  1.0 → quality 3 (correct but slow)
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
