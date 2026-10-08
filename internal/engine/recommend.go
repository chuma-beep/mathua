package engine

import (
	"fmt"
	"time"

	"github.com/chuma-beep/mathua/internal/mastery"
	"github.com/chuma-beep/mathua/internal/scheduler"
	"github.com/chuma-beep/mathua/internal/storage"
)

// conceptSnapshots assembles the scheduler's view of a learner.
//
// This was written out three times — NextQuestion, NextReviewQuestion and the recommender —
// which is how they came to disagree. The recommender in particular could not reuse either: it
// needs the whole DAG including unseen concepts, and it needs to run outside a session.
//
// One definition, so "what should this learner do next" cannot answer differently depending on
// which caller asked.
func (e *Engine) conceptSnapshotsFrom(studentID string, progress map[string]*storage.ConceptProgress) map[string]*scheduler.ConceptSnapshot {
	snapshots := make(map[string]*scheduler.ConceptSnapshot, len(progress))
	for cid, p := range progress {
		reqStreak := 0
		timeThresh := 0.0
		if c := e.dag.Concept(cid); c != nil {
			reqStreak = c.MasteryThreshold.Streak
			timeThresh = c.MasteryThreshold.AvgTimeSeconds
		}
		snapshots[cid] = &scheduler.ConceptSnapshot{
			Status:         mastery.Status(p.Status),
			Streak:         p.Streak,
			LastAttempted:  timeOrZero(p.LastAttempted),
			LastReviewed:   timeOrZero(p.LastReviewed),
			NextReviewDue:  p.NextReviewDue,
			RequiredStreak: reqStreak,
			TimeThreshold:  timeThresh,
			WeaknessScore:  p.WeaknessScore,
			Placement:      p.PlacementSeeded,
		}
	}
	// Concepts never attempted still need a snapshot: prereqsMet returns false for a missing
	// one, so an absent concept reads as blocked rather than as available-and-unstarted.
	for _, c := range e.dag.Order() {
		if _, ok := snapshots[c.ID]; !ok {
			snapshots[c.ID] = &scheduler.ConceptSnapshot{
				Status:         mastery.StatusUnseen,
				RequiredStreak: c.MasteryThreshold.Streak,
				TimeThreshold:  c.MasteryThreshold.AvgTimeSeconds,
			}
		}
	}
	return snapshots
}

// conceptSnapshots is the loading wrapper, for the question-serving paths.
func (e *Engine) conceptSnapshots(studentID string) (map[string]*scheduler.ConceptSnapshot, error) {
	progress, err := e.repo.GetAllProgress(studentID)
	if err != nil {
		return nil, fmt.Errorf("load progress: %w", err)
	}
	return e.conceptSnapshotsFrom(studentID, progress), nil
}

// ReviewsDue counts concepts whose SM-2 interval has elapsed. Moved here from the server,
// which computed it inline and which the recommender needs it from.
func (e *Engine) ReviewsDue(studentID string) (int, error) {
	progress, err := e.repo.GetAllProgress(studentID)
	if err != nil {
		return 0, fmt.Errorf("load progress: %w", err)
	}
	now := time.Now()
	count := 0
	for _, p := range progress {
		if p.NextReviewDue == nil {
			continue
		}
		if p.NextReviewDue.After(now) {
			continue
		}
		// A stored MASTERED row with a due review IS due a review — scheduler.effectiveState
		// says so explicitly, which is what put it in the review bucket in the first place.
		// The old guard skipped exactly those rows, so the count and the candidates disagreed
		// and selectTasks never took the review bucket for a learner whose only debt was a
		// decayed concept. Any NextReviewDue at or before now is a due review; that is the
		// whole condition.
		count++
	}
	return count, nil
}

// RecommendNext is the learner's head plus their alternatives. It is the answer to "what
// should I do now?", and it is answered once, here, rather than recomputed in the browser from
// five API calls and a bundled copy of the corpus.
func (e *Engine) RecommendNext(studentID string, exclude []string) (scheduler.RecommendationResponse, error) {
	if studentID == "" {
		return scheduler.New(e.dag).Recommend(scheduler.RecommendInput{
			Now: time.Now(), Weakness: map[string]float64{},
		}), nil
	}
	progress, err := e.repo.GetAllProgress(studentID)
	if err != nil {
		return scheduler.RecommendationResponse{Alternatives: []scheduler.Recommendation{}}, fmt.Errorf("load progress: %w", err)
	}
	snapshots := e.conceptSnapshotsFrom(studentID, progress)
	weakness := e.weaknessMapFromProgress(progress)
	due, err := e.ReviewsDue(studentID)
	if err != nil {
		return scheduler.RecommendationResponse{Alternatives: []scheduler.Recommendation{}}, err
	}
	quizDue, _ := e.QuizDue(studentID)

	// How far along the learner is, in the two terms the recommender needs. Mastered is
	// effective MASTERED — decay applied, so a stale concept has left the count rather than
	// inflating it, and *not* "anything not unseen", which counted a concept the learner had
	// merely started. Only the new-learner gate consumes this today and AttemptedDays already
	// guarded it, so the miscount was unobservable; it is fixed because the field is named
	// after a fact it did not report, and the next caller inherits discrepancies rather than
	// reading the comment.
	mastered, attemptedDays := 0, map[string]bool{}
	for _, p := range progress {
		if mastery.EffectiveStatus(mastery.Status(p.Status), daysSinceReview(p.LastReviewed), mastery.DecayDays) == mastery.StatusMastered {
			mastered++
		}
		if p.LastAttempted != nil {
			attemptedDays[p.LastAttempted.UTC().Format("2006-01-02")] = true
		}
	}

	e.mu.Lock()
	if path := e.activePath[studentID]; len(path) > 0 {
		for cid := range snapshots {
			if !path[cid] {
				delete(snapshots, cid)
			}
		}
	}
	e.mu.Unlock()

	// The server's record that placement happened. This is the authoritative
	// guard against offering the diagnostic back as an active task; the
	// AttemptedDays check alone is not enough, because placement deliberately
	// seeds reachability without mastery.
	diagnosticCompleted := false
	if st, err := e.repo.GetStudent(studentID); err == nil && st != nil {
		diagnosticCompleted = st.DiagnosticCompleted
	}

	return scheduler.New(e.dag).Recommend(scheduler.RecommendInput{
		Snapshots:           snapshots,
		Now:                 time.Now(),
		Weakness:            weakness,
		DueReviews:          due,
		QuizDue:             quizDue,
		ConceptsMastered:    mastered,
		AttemptedDays:       len(attemptedDays),
		ExcludeConceptIds:   exclude,
		DiagnosticCompleted: diagnosticCompleted,
	}), nil
}

func daysSinceReview(t *time.Time) float64 {
	if t == nil || t.IsZero() {
		return 999
	}
	return time.Since(*t).Hours() / 24
}

// RecommendPaused is what a paused learner is told. Reviews are hidden while paused (ADR-022
// and handleDueReviews agree on that), so a recommendation pointing at /review would
// contradict the pause the learner just set.
//
// It used to point at `/study` instead, with "Lessons are still open" — true when it was
// written, and a dead end once the library closed: a learner who opened the app on a break was
// sent to a page that could only explain things they had already decided not to practise today.
// `/settings` is where the pause and its end both live, so it is the one destination that can
// actually answer "I am not doing reviews — now what".
func (e *Engine) RecommendPaused(studentID string) scheduler.RecommendationResponse {
	return scheduler.RecommendationResponse{
		Primary: &scheduler.Recommendation{
			ID:     "paused",
			Kind:   scheduler.KindLearn,
			Reason: scheduler.ReasonNew,
			Action: scheduler.Action{Type: scheduler.KindLearn, Href: "/settings"},
			Badge:  "Paused",
			Detail: "Reviews are on hold until you resume — due reviews are hidden while paused",
			CTA:    "Change or resume →",
		},
		Alternatives: []scheduler.Recommendation{},
		GeneratedAt:  time.Now().UTC(),
	}
}
