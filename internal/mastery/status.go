package mastery

type Status string

const (
	StatusUnseen     Status = "UNSEEN"
	StatusLearning   Status = "LEARNING"
	StatusPracticing Status = "PRACTICING"
	StatusMastered   Status = "MASTERED"
)

// EffectiveStatus returns DECAYING when a MASTERED concept has not been
// reviewed for longer than the decay threshold (in days). Decay is a
// read-time signal — the database always stores MASTERED.
func EffectiveStatus(current Status, daysSinceLastReview float64, decayThresholdDays float64) Status {
	if current == StatusMastered && daysSinceLastReview >= decayThresholdDays {
		return "DECAYING"
	}
	return current
}
