package mastery

import (
	"time"

	"github.com/chuma-beep/mathua/internal/fsrs"
)

type Status string

const (
	StatusUnseen     Status = "UNSEEN"
	StatusLearning   Status = "LEARNING"
	StatusPracticing Status = "PRACTICING"
	StatusMastered   Status = "MASTERED"
)

// DecayDays is how long a mastered concept stays fresh before it reads as
// DECAYING when no dynamic FSRS stability is present.
//
// This was written four times: as `14` in the engine, as a bare `14` twice in
// the scheduler, and as `30` in the workload estimator. The estimator therefore
// told learners their reviews came round every 30 days while the scheduler
// scheduled them at 14, a plan built on a review cadence the app would not
// honour, which is the one number a deadline is computed from (ADR-019 prices
// reviews on this path).
//
// One value, owned here beside EffectiveStatus, which is the only place that
// interprets it.
const DecayDays = 14

// EffectiveStatus returns DECAYING if a MASTERED concept hasn't been reviewed
// beyond the decay threshold (in days). Decay is computed at read time,
// the database never stores it.
func EffectiveStatus(current Status, daysSinceLastReview float64, decayThresholdDays float64) Status {
	return EffectiveStatusFSRS(current, daysSinceLastReview, nil, decayThresholdDays)
}

// EffectiveStatusFSRS evaluates whether a MASTERED concept has decayed.
// If stability is present and positive, it evaluates FSRS-6 retrievability:
// R(t, S) = (1 + Factor * t / S)^Decay.
// When R < DefaultRetention (0.90), the concept has decayed to "DECAYING".
// If stability is nil or non-positive, it falls back to decayThresholdDays (or DecayDays).
func EffectiveStatusFSRS(current Status, daysSinceLastReview float64, stability *float64, decayThresholdDays float64) Status {
	if current != StatusMastered {
		return current
	}
	if stability != nil && *stability > 0 {
		r := fsrs.DefaultParameters().Retrievability(daysSinceLastReview, *stability)
		if r < fsrs.DefaultRetention {
			return "DECAYING"
		}
		return StatusMastered
	}
	if decayThresholdDays <= 0 {
		decayThresholdDays = DecayDays
	}
	if daysSinceLastReview >= decayThresholdDays {
		return "DECAYING"
	}
	return current
}

// ReviewElapsedDays is the age of the most recent evidence that a mastered
// concept was reviewed, in days. It walks four timestamps in order of
// authority: an FSRS review, an SM-2 review, the moment mastery was awarded,
// then the last attempt.
//
// Mastery is the fallback that matters. A concept mastered through study has
// no review timestamp until its first review, and the two readers used to
// disagree about what that meant: the graph assumed 999 days and called it
// decayed, the scheduler assumed zero and called it fresh. That made every
// such concept appear due on the graph while the queue never offered it. Both
// now read the mastery date, so the decay clock starts when mastery is earned.
//
// A row with none of the four timestamps reports zero days: there is no
// evidence it is stale, and the versioned backfill gives legacy rows a real
// clock at boot.
func ReviewElapsedDays(fsrsLastReview, lastReviewed, masteredAt, lastAttempted *time.Time, now time.Time) float64 {
	for _, t := range []*time.Time{fsrsLastReview, lastReviewed, masteredAt, lastAttempted} {
		if t == nil || t.IsZero() {
			continue
		}
		days := now.Sub(*t).Hours() / 24.0
		if days < 0 {
			return 0
		}
		return days
	}
	return 0
}
