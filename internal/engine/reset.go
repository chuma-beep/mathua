package engine

import (
	"fmt"
	"strings"
	"time"
)

// ResetPhrase is the exact typed confirmation the reset endpoint requires.
// Client-side gating is UX only; the server enforces this value.
const ResetPhrase = "reset my progress"

var errResetPhraseMismatch = fmt.Errorf("confirmation phrase does not match")

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
	if strings.TrimSpace(phrase) != ResetPhrase {
		return errResetPhraseMismatch
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
		}
	}
	for k := range e.studyExpected {
		if strings.HasPrefix(k, prefix) {
			delete(e.studyExpected, k)
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
