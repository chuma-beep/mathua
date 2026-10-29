package engine

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/chuma-beep/mathua/internal/fire"
	"github.com/chuma-beep/mathua/internal/storage"
)

const (
	studyEffectPropagateWeakness = "propagate_weakness"
	studyEffectImplicitCredit    = "implicit_credit"
	studyEffectProvisional       = "provisional_fallback"
	studyEffectRemedial          = "remedial"
	studyEffectRemoveRemedial    = "remove_remedial"
)

type studyEffectPayload struct {
	ConceptID  string   `json:"concept_id,omitempty"`
	AttemptID  string   `json:"attempt_id,omitempty"`
	ConceptIDs []string `json:"concept_ids,omitempty"`
	Correct    bool     `json:"correct"`
	Source     string   `json:"source,omitempty"`
}

func studyEffectJSON(p studyEffectPayload) string {
	b, err := json.Marshal(p)
	if err != nil {
		log.Printf("warning: encode study effect payload: %v", err)
		return "{}"
	}
	return string(b)
}

// studyEffectsFor lists the post-commit work in the order the outbox should
// apply it. The rows are written in the same transaction as the answer, so a
// crash after commit leaves durable work instead of a log line.
func (e *Engine) studyEffectsFor(source, conceptID, attemptID string, correct, halted bool, remedial []string) []storage.StudyEffectInput {
	effects := []storage.StudyEffectInput{
		{Effect: studyEffectPropagateWeakness},
	}
	if e.fireMode.CollectsLedger() {
		effects = append(effects, storage.StudyEffectInput{
			Effect: studyEffectImplicitCredit,
			PayloadJSON: studyEffectJSON(studyEffectPayload{
				ConceptID: conceptID,
				AttemptID: attemptID,
				Correct:   correct,
				Source:    source,
			}),
		})
	}
	if correct {
		effects = append(effects, storage.StudyEffectInput{
			Effect:      studyEffectRemoveRemedial,
			PayloadJSON: studyEffectJSON(studyEffectPayload{ConceptID: conceptID}),
		})
		return effects
	}
	effects = append(effects, storage.StudyEffectInput{
		Effect:      studyEffectProvisional,
		PayloadJSON: studyEffectJSON(studyEffectPayload{ConceptID: conceptID}),
	})
	if halted && len(remedial) > 0 {
		effects = append(effects, storage.StudyEffectInput{
			Effect:      studyEffectRemedial,
			PayloadJSON: studyEffectJSON(studyEffectPayload{ConceptIDs: remedial}),
		})
	}
	return effects
}

// ProcessStudyEffects claims pending effects and applies them. It is safe to
// call after every answer and from the background worker. A failure is left
// pending with a retry time rather than failing the answer that produced it.
func (e *Engine) ProcessStudyEffects(studentID string, limit int) error {
	store, ok := e.repo.(storage.StudyEffectStore)
	if !ok {
		return nil
	}
	if limit <= 0 {
		limit = 50
	}
	now := time.Now().UTC()
	effects, err := store.ClaimStudyEffects(studentID, now, limit)
	if err != nil {
		return fmt.Errorf("claim study effects: %w", err)
	}
	for _, effect := range effects {
		if err := e.applyStudyEffect(effect); err != nil {
			log.Printf("warning: study effect %d %s for %s/%s failed: %v", effect.ID, effect.Effect, effect.StudentID, effect.SubmissionID, err)
			retryAt := now.Add(studyEffectBackoff(effect.Attempts))
			if ferr := store.FailStudyEffect(effect.ID, err.Error(), retryAt); ferr != nil {
				log.Printf("warning: fail study effect %d: %v", effect.ID, ferr)
			}
			continue
		}
		if err := store.CompleteStudyEffect(effect.ID); err != nil {
			log.Printf("warning: complete study effect %d: %v", effect.ID, err)
		}
	}
	return nil
}

// PruneStudyEffects removes completed outbox rows older than the cutoff.
// Pending and failed rows stay, so retry state remains inspectable.
func (e *Engine) PruneStudyEffects(olderThan time.Time, limit int) (int64, error) {
	store, ok := e.repo.(storage.StudyEffectStore)
	if !ok {
		return 0, nil
	}
	return store.PruneStudyEffects(olderThan, limit)
}

func studyEffectBackoff(attempts int) time.Duration {
	if attempts < 1 {
		attempts = 1
	}
	if attempts > 6 {
		attempts = 6
	}
	return time.Duration(1<<attempts) * time.Second
}

func (e *Engine) applyStudyEffect(effect storage.StudyAnswerEffect) error {
	var p studyEffectPayload
	if effect.PayloadJSON != "" {
		if err := json.Unmarshal([]byte(effect.PayloadJSON), &p); err != nil {
			return fmt.Errorf("decode %s payload: %w", effect.Effect, err)
		}
	}
	switch effect.Effect {
	case studyEffectPropagateWeakness:
		return e.propagateWeakness(effect.StudentID)
	case studyEffectImplicitCredit:
		if effect.Attempts > 0 {
			e.fireStats.Retries.Add(1)
		}
		if e.fireMode == fire.Off {
			// An effect queued before the mode changed must be drained, not
			// retried forever. Off means no ledger write and no projection.
			return nil
		}
		source := p.Source
		if source == "" {
			// Effects queued before source was recorded drained from the
			// practice path they were written for.
			source = "practice"
		}
		return e.allocateImplicitCreditErr(gradedAttempt{
			studentID: effect.StudentID,
			conceptID: p.ConceptID,
			attemptID: p.AttemptID,
			correct:   p.Correct,
			source:    source,
		})
	case studyEffectProvisional:
		_, err := e.CheckProvisionalFallback(effect.StudentID, p.ConceptID)
		return err
	case studyEffectRemedial:
		if len(p.ConceptIDs) == 0 {
			return nil
		}
		return e.repo.EnqueueRemedial(effect.StudentID, p.ConceptIDs...)
	case studyEffectRemoveRemedial:
		return e.repo.RemoveRemedial(effect.StudentID, p.ConceptID)
	default:
		return fmt.Errorf("unknown study effect %q", effect.Effect)
	}
}
