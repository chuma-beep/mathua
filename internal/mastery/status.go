package mastery

type Status string

const (
	StatusUnseen     Status = "UNSEEN"
	StatusLearning   Status = "LEARNING"
	StatusPracticing Status = "PRACTICING"
	StatusMastered   Status = "MASTERED"
)

// EffectiveStatus returns DECAYING if a MASTERED concept hasn't been reviewed
// beyond the decay threshold (in days). Decay is computed at read time —
// the database never stores it.
func EffectiveStatus(current Status, daysSinceLastReview float64, decayThresholdDays float64) Status {
	if current == StatusMastered && daysSinceLastReview >= decayThresholdDays {
		return "DECAYING"
	}
	return current
}
