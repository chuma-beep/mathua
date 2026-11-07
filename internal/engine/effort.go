package engine

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chuma-beep/mathua/internal/mastery"
	"github.com/chuma-beep/mathua/internal/storage"
	xppolicy "github.com/chuma-beep/mathua/internal/xp"
)

// effortObservation is what a graded answer contributes to accountability. It
// is deliberately smaller than a submission: no answer text, no expected
// answer, no progress row. The layer is not allowed to influence learning.
type effortObservation struct {
	StudentID     string
	TaskKind      string
	TaskID        string
	ReviewRunID   string
	ConceptID     string
	AttemptKey    string
	Answer        string
	Correct       bool
	DontKnow      bool
	Mastered      bool
	LegacyPenalty bool
	Elapsed       float64
	Difficulty    float64
	XPEarned      int
	// TaskComplete is set only when the caller has a reliable completion
	// signal. Performance waits for it; effort does not.
	TaskComplete bool
	At           time.Time
}

// effortAssessment is a decision plus the observation that produced it. It is
// built before an idempotent submission commits so the stored response can
// explain the same decision a retry will see.
type effortAssessment struct {
	decision xppolicy.Decision
	signal   storage.EffortSignal
	enabled  bool
}

// taskIDForEffort names the task a deduction is charged against. One task can
// never be charged twice because the ledger event key is derived from it.
//
// Practice and review use the session, because a session is the learner's
// working set. Study and quiz use the concept and UTC day, because those flows
// do not hand the engine a per-task id and a day is the smallest honest bound.
func taskIDForEffort(taskKind, sessionID, conceptID string, now time.Time) string {
	switch taskKind {
	case xppolicy.TaskPractice:
		if sessionID != "" {
			return sessionID
		}
	case xppolicy.TaskReview:
		if sessionID != "" {
			return sessionID + ":" + conceptID
		}
	}
	return conceptID + ":" + now.UTC().Format("2006-01-02")
}

// effortAnswerSignature is a one-way digest used only to notice the same wrong
// answer repeated. It never leaves the database and is not reversible.
func effortAnswerSignature(answer string) string {
	a := strings.ToLower(strings.TrimSpace(answer))
	if a == "" {
		return ""
	}
	sum := sha256.Sum256([]byte(a))
	return hex.EncodeToString(sum[:8])
}

func effortMastered(progress *storage.ConceptProgress) bool {
	return progress != nil && progress.Status == string(mastery.StatusMastered)
}

func effortDifficulty(v *float64) float64 {
	if v == nil {
		return 0
	}
	return *v
}

// assessEffort computes the decision without writing anything. It loads a
// bounded slice of recent observations, folds in the current answer, and caps
// the deduction at the balance the learner will actually hold once the current
// award lands.
func (e *Engine) assessEffort(o effortObservation, pendingAward int) effortAssessment {
	out := effortAssessment{decision: xppolicy.Decision{Outcome: xppolicy.OutcomeFull}}
	store, ok := e.repo.(storage.EffortStore)
	if !ok || e.effortConfig.Disabled {
		return out
	}
	cfg := xppolicy.NormalizeConfig(e.effortConfig)
	now := o.At
	if now.IsZero() {
		now = time.Now().UTC()
	}
	sig := storage.EffortSignal{
		StudentID:       o.StudentID,
		TaskKind:        o.TaskKind,
		TaskID:          o.TaskID,
		ReviewRunID:     o.ReviewRunID,
		ConceptID:       o.ConceptID,
		AttemptKey:      o.AttemptKey,
		Correct:         o.Correct,
		DontKnow:        o.DontKnow,
		Mastered:        o.Mastered,
		LegacyPenalty:   o.LegacyPenalty,
		ElapsedSeconds:  o.Elapsed,
		Difficulty:      o.Difficulty,
		XPEarned:        o.XPEarned,
		AnswerSignature: effortAnswerSignature(o.Answer),
		CreatedAt:       now,
	}
	out.signal = sig
	out.enabled = true

	// One query covers both the current task and the farming window. The
	// largest window wins so a long practice session still has its own history.
	lookback := cfg.FarmingWindow
	if lookback < 6*time.Hour {
		lookback = 6 * time.Hour
	}
	recent, err := store.EffortSignalsSince(o.StudentID, now.Add(-lookback), 400)
	if err != nil {
		log.Printf("warning: read effort signals for %s: %v", o.StudentID, err)
	}

	attempts := make([]xppolicy.AttemptSignal, 0, 8)
	recentAttempts := make([]xppolicy.AttemptSignal, 0, len(recent)+1)
	farmingFrom := now.Add(-cfg.FarmingWindow)
	seenCurrent := false
	for _, s := range recent {
		if s.TaskKind == o.TaskKind && s.TaskID == o.TaskID && s.AttemptKey == o.AttemptKey {
			seenCurrent = true
		}
		if s.TaskKind == o.TaskKind && s.TaskID == o.TaskID {
			attempts = append(attempts, effortAttemptSignal(s))
		}
		if !s.CreatedAt.Before(farmingFrom) {
			recentAttempts = append(recentAttempts, effortAttemptSignal(s))
		}
	}
	if !seenCurrent {
		attempts = append(attempts, effortAttemptSignal(sig))
		recentAttempts = append(recentAttempts, effortAttemptSignal(sig))
	}

	prior, err := store.CountXPDeductions(o.StudentID, now.Add(-cfg.EscalationWindow))
	if err != nil {
		log.Printf("warning: count xp deductions for %s: %v", o.StudentID, err)
		prior = 0
	}
	available, _, err := e.repo.GetXP(o.StudentID)
	if err != nil {
		available = 0
	}
	if pendingAward > 0 {
		available += pendingAward
	}

	in := xppolicy.Input{
		TaskKind:         o.TaskKind,
		TaskID:           o.TaskID,
		Attempts:         attempts,
		Recent:           recentAttempts,
		PriorAdjustments: prior,
		AvailableXP:      available,
		TaskComplete:     o.TaskComplete,
		Config:           cfg,
	}
	out.decision = xppolicy.Combine(xppolicy.Evaluate(in), xppolicy.EvaluatePerformance(in))
	e.accountabilityStats.recordDecision(out.decision)
	if o.TaskKind != xppolicy.TaskDiagnostic {
		if !o.TaskComplete {
			e.accountabilityStats.recordIncomplete()
		} else if len(attempts) < cfg.Performance.MinAttempts {
			e.accountabilityStats.recordShortTask()
		}
	}
	return out
}

func effortAttemptSignal(s storage.EffortSignal) xppolicy.AttemptSignal {
	return xppolicy.AttemptSignal{
		TaskID:          s.TaskID,
		Correct:         s.Correct,
		ElapsedSeconds:  s.ElapsedSeconds,
		DontKnow:        s.DontKnow,
		Mastered:        s.Mastered,
		Difficulty:      s.Difficulty,
		XPEarned:        s.XPEarned,
		AnswerSignature: s.AnswerSignature,
		LegacyPenalty:   s.LegacyPenalty,
		At:              s.CreatedAt,
	}
}

// apply records the observation and, when the decision calls for it, writes one
// ledger row for the task. It returns the amount actually applied: zero can
// mean "no adjustment", "balance already empty" or "another request won the
// same idempotency key", and none of those is an error.
func (a effortAssessment) apply(e *Engine) int {
	if !a.enabled {
		return 0
	}
	store, ok := e.repo.(storage.EffortStore)
	if !ok {
		return 0
	}
	if err := store.RecordEffortSignal(a.signal); err != nil {
		log.Printf("warning: record effort signal for %s: %v", a.signal.StudentID, err)
	}
	if a.decision.Amount == 0 {
		return 0
	}
	eventKey := "accountability:" + a.signal.TaskKind + ":" + a.signal.TaskID
	reason := a.decision.LedgerReason()
	applied, err := store.ApplyXPAdjustment(a.signal.StudentID, a.decision.Amount, reason, eventKey)
	if err != nil {
		e.accountabilityStats.recordStorageError()
		log.Printf("warning: apply xp adjustment for %s: %v", a.signal.StudentID, err)
		return 0
	}
	e.accountabilityStats.recordApplied(a.decision, applied)
	return applied
}

// noteForApplied turns a decision and its actual applied amount into the
// student-facing fields. A deduction another request already made is not
// reported twice; a balance-limited deduction reports the real number.
func noteForApplied(d xppolicy.Decision, applied int) (int, string, string) {
	if applied < 0 {
		if d.Amount < 0 && applied == d.Amount {
			return applied, d.Note, string(d.Reason)
		}
		// Capped by the balance or already applied; rebuild the sentence from
		// the amount the ledger actually holds.
		return applied, fmt.Sprintf("%d XP: %s", applied, d.Explanation), string(d.Reason)
	}
	if applied > 0 {
		return applied, d.Note, string(d.Reason)
	}
	if d.Outcome == xppolicy.OutcomeZero || d.Reason == xppolicy.ReasonLegacyPenalty {
		return 0, d.Note, string(d.Reason)
	}
	return 0, "", ""
}
