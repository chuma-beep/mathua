package mastery

type Status string

const (
	StatusUnseen     Status = "UNSEEN"
	StatusLearning   Status = "LEARNING"
	StatusPracticing Status = "PRACTICING"
	StatusMastered   Status = "MASTERED"
)

// DecayDays is how long a mastered concept stays fresh before it reads as
// DECAYING.
//
// This was written four times: as `14` in the engine, as a bare `14` twice in
// the scheduler, and as `30` in the workload estimator. The estimator therefore
// told learners their reviews came round every 30 days while the scheduler
// scheduled them at 14 — a plan built on a review cadence the app would not
// honour, which is the one number a deadline is computed from (ADR-019 prices
// reviews on this path).
//
// One value, owned here beside EffectiveStatus, which is the only place that
// interprets it.
const DecayDays = 14

// EffectiveStatus returns DECAYING if a MASTERED concept hasn't been reviewed
// beyond the decay threshold (in days). Decay is computed at read time —
// the database never stores it.
func EffectiveStatus(current Status, daysSinceLastReview float64, decayThresholdDays float64) Status {
	if current == StatusMastered && daysSinceLastReview >= decayThresholdDays {
		return "DECAYING"
	}
	return current
}
