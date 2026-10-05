package engine

import (
	"fmt"
	"strings"
	"time"
)

// ResetPhrase is the typed confirmation the reset endpoint requires.
// Client-side gating is UX only; the server enforces this value.
const ResetPhrase = "reset my progress"

// PhraseMatches reports whether a typed confirmation is acceptable for want.
//
// Case-insensitive, and deliberately so. The confirmation is typed into a `type="text"`
// field, which Safari capitalises as you type and Android keyboards routinely
// autocorrect — "reset my progress" arrives as "Reset my progress". An exact comparison
// therefore rejected people who had typed precisely what the label asked for, and the
// error they got back said they had mistyped it.
//
// The phrase confirms intent; it is not a secret. Casing and surrounding whitespace say
// nothing about how deliberate the action was, and the two phrases stay separate values
// so a reset confirmation can never authorize a deletion.
func PhraseMatches(typed, want string) bool {
	return strings.EqualFold(strings.TrimSpace(typed), want)
}

// DeletePhrase is the typed confirmation the account-deletion
// endpoint requires. Separate from ResetPhrase so a reset confirmation can
// never authorize a deletion.
const DeletePhrase = "delete my account"

// ErrResetPhraseMismatch and ErrDeletePhraseMismatch are the only expected
// failures on these paths. They are exported so the HTTP layer can tell "you
// typed the wrong words" apart from "the database refused": collapsing both
// into one 400 told a learner with a real infrastructure fault that their
// confirmation phrase was wrong, which is unsendable advice.
var ErrResetPhraseMismatch = fmt.Errorf("confirmation phrase does not match")

var ErrDeletePhraseMismatch = fmt.Errorf("confirmation phrase does not match")

// ResetAccountProgress wipes one student's learning record and returns them
// to a fresh start. Order matters and every step is retry-safe:
//
//  1. Repository reset (single transaction; itself idempotent).
//  2. In-memory engine state for the student (anchors, misses, remedials,
//     active path) — otherwise stale grading state survives the wipe.
//  3. Evidence baseline delete (KV last: if the process dies between 1
//     and 3, a retry simply repeats idempotent deletes).
//  4. Prefs reset marker (keeps destination/deadline/pace, stamps ResetAt
//     so Plan shows the neutral fresh-start state until a new baseline).
func (e *Engine) ResetAccountProgress(studentID, phrase string) error {
	if !PhraseMatches(phrase, ResetPhrase) {
		return ErrResetPhraseMismatch
	}
	if err := e.repo.ResetProgress(studentID); err != nil {
		return err
	}
	e.mu.Lock()
	delete(e.studySessions, studentID)
	delete(e.quizRemedial, studentID)
	delete(e.activePath, studentID)
	prefix := studentID + "|"
	for k := range e.studyMisses {
		if strings.HasPrefix(k, prefix) {
			delete(e.studyMisses, k)
			delete(e.studyRush, k)
		}
	}
	for k := range e.studyAnchor {
		if strings.HasPrefix(k, prefix) {
			delete(e.studyAnchor, k)
		}
	}
	e.mu.Unlock()
	// KV last and retry-safe: deletes tolerate missing rows, the prefs
	// upsert is idempotent, and ResetProgress itself is a no-op rerun.
	if err := e.repo.DeleteServerSession(serverSessionStudyPlan, studentID); err != nil {
		return err
	}
	prefs, _ := e.GetPlanPrefs(studentID)
	if prefs == nil {
		prefs = &StudyPlanPrefs{}
	}
	prefs.ResetAt = time.Now().UTC().Format(time.RFC3339)
	return e.savePlanPrefs(studentID, prefs)
}

// DeleteAccount removes one student and every owned row, then drops their
// in-memory engine state (same maps as the reset path). Order matters and
// every step is retry-safe: the repository delete is one idempotent
// transaction, and map deletes tolerate absence. Password verification (for
// password accounts) happens in the server handler before this runs.
func (e *Engine) DeleteAccount(studentID, phrase string) error {
	if !PhraseMatches(phrase, DeletePhrase) {
		return ErrDeletePhraseMismatch
	}
	if err := e.repo.DeleteAccount(studentID); err != nil {
		return err
	}
	e.mu.Lock()
	delete(e.studySessions, studentID)
	delete(e.quizRemedial, studentID)
	delete(e.activePath, studentID)
	prefix := studentID + "|"
	for k := range e.studyMisses {
		if strings.HasPrefix(k, prefix) {
			delete(e.studyMisses, k)
			delete(e.studyRush, k)
		}
	}
	for k := range e.studyAnchor {
		if strings.HasPrefix(k, prefix) {
			delete(e.studyAnchor, k)
		}
	}
	e.mu.Unlock()
	return nil
}
