package engine

import (
	"fmt"
	"log"
	"time"

	"github.com/chuma-beep/mathua/internal/fsrs"
	"github.com/chuma-beep/mathua/internal/mastery"
	"github.com/chuma-beep/mathua/internal/repetition"
	"github.com/chuma-beep/mathua/internal/scheduler"
	"github.com/chuma-beep/mathua/internal/storage"
)

// servedPerDayCap is how many times one concept may be served in a day before the scheduler
// drops it from the ranked list. Two is chosen rather than tuned: the spacing literature puts
// massed practice several tenths of a standard deviation below spaced practice in mathematics,
// and a concept asked at most twice a day stays inside the supply of the thinnest concepts in
// the corpus. The guard against adjacent repeats is separate and lives in the scheduler, so
// this number bounds the day rather than the sitting.
//
// Measured on 2026-10-24 by simulating 30 days for four learner profiles 25 times each, with
// and without the guard, and scoring predicted recall seven days after every review. With the
// guard a steady learner reached 0.710 against 0.644, a struggling learner 0.493 against 0.368,
// a 120-item learner 0.710 against 0.585, and a learner who misses days 0.727 against 0.648.
// Retrievability sampled at the moment of review is the wrong yardstick: it reads lower for a
// well-spaced schedule precisely because spacing lets recall decay before the next review.
const servedPerDayCap = 2

// SatisfyImplicitReviews inspects banked implicit repetition credit for concepts that appear
// due for explicit spaced review. If an encompassed concept has banked sufficient implicit credit,
// it projects the credit, advances the review interval and updates storage.
func (e *Engine) SatisfyImplicitReviews(studentID string) error {
	if studentID == "" {
		return nil
	}
	cs, ok := e.repo.(storage.CreditStore)
	if !ok {
		return nil
	}
	banked, err := cs.GetBankedOrdinals(studentID)
	if err != nil || len(banked) == 0 {
		return err
	}
	progress, err := e.repo.GetAllProgress(studentID)
	if err != nil {
		return fmt.Errorf("load progress: %w", err)
	}
	now := time.Now().UTC()
	for cid, count := range banked {
		if count <= 0 {
			continue
		}
		if e.dag == nil || len(e.dag.EncompassedBy(cid)) == 0 {
			continue
		}
		p := progress[cid]
		if p == nil {
			continue
		}
		isDue := false
		if p.NextReviewDue != nil && !p.NextReviewDue.After(now) {
			isDue = true
		}
		var elapsed float64
		if p.FSRSLastReview != nil && !p.FSRSLastReview.IsZero() {
			elapsed = now.Sub(*p.FSRSLastReview).Hours() / 24.0
		} else if p.LastReviewed != nil && !p.LastReviewed.IsZero() {
			elapsed = now.Sub(*p.LastReviewed).Hours() / 24.0
		}
		if elapsed < 0 {
			elapsed = 0
		}
		if p.FSRSStability != nil && *p.FSRSStability > 0 {
			if fsrs.DefaultParameters().Retrievability(elapsed, *p.FSRSStability) < fsrs.DefaultRetention {
				isDue = true
			}
		} else if elapsed >= mastery.DecayDays {
			isDue = true
		}
		if !isDue {
			continue
		}

		speed := 1.0
		if ts, err := e.repo.GetTopicSpeed(studentID, cid); err == nil && ts != nil {
			speed = ts.LearningSpeed
		}
		prevSM2 := scheduler.SM2{
			Interval:    p.SM2Interval,
			EFactor:     p.SM2EFactor,
			Repetitions: p.SM2Repetitions,
		}
		stepDays, err := repetition.StepDays(prevSM2, speed)
		if err != nil || stepDays < 1 {
			stepDays = 1
		}
		ords, err := cs.ProjectCredit(studentID, cid, 0, stepDays, now)
		if err != nil {
			log.Printf("warning: project credit %s/%s: %v", studentID, cid, err)
			continue
		}
		if len(ords) == 0 {
			continue
		}
		lastDueStr := ords[len(ords)-1].CounterfactualDue
		newDue, err := time.Parse(time.RFC3339, lastDueStr)
		if err != nil {
			newDue, err = time.Parse("2006-01-02 15:04:05", lastDueStr)
		}
		if err != nil {
			newDue = now.AddDate(0, 0, stepDays)
		}
		p.NextReviewDue = &newDue
		p.LastReviewed = &now
		p.SM2Interval = stepDays
		if p.Status == "DECAYING" {
			p.Status = string(mastery.StatusMastered)
		}
		if p.FSRSStability != nil && *p.FSRSStability > 0 {
			params := fsrs.DefaultParameters()
			diff := params.InitialDifficulty(fsrs.Good)
			if p.FSRSDifficulty != nil && *p.FSRSDifficulty > 0 {
				diff = *p.FSRSDifficulty
			}
			var elapsed float64
			if p.FSRSLastReview != nil && !p.FSRSLastReview.IsZero() {
				elapsed = now.Sub(*p.FSRSLastReview).Hours() / 24.0
			}
			if elapsed < 0 {
				elapsed = 0
			}
			r := params.Retrievability(elapsed, *p.FSRSStability)
			nextStab := params.NextStability(*p.FSRSStability, diff, r, fsrs.Good)
			p.FSRSStability = &nextStab
			p.FSRSLastReview = &now
		}
		if err := e.repo.UpsertProgress(p); err != nil {
			log.Printf("warning: upsert progress after credit projection %s/%s: %v", studentID, cid, err)
		}
	}
	return nil
}

// conceptSnapshots assembles the scheduler's view of a learner.
//
// This was written out three times, NextQuestion, NextReviewQuestion and the recommender,
// which is how they came to disagree. The recommender in particular could not reuse either: it
// needs the whole DAG including unseen concepts, and it needs to run outside a session.
//
// One definition, so "what should this learner do next" cannot answer differently depending on
// which caller asked.
func (e *Engine) conceptSnapshotsFrom(studentID string, progress map[string]*storage.ConceptProgress) map[string]*scheduler.ConceptSnapshot {
	var banked map[string]int
	if cs, ok := e.repo.(storage.CreditStore); ok && studentID != "" {
		if b, err := cs.GetBankedOrdinals(studentID); err == nil {
			banked = b
		}
	}
	var speeds map[string]*storage.TopicSpeed
	if studentID != "" {
		if sp, err := e.repo.GetAllTopicSpeeds(studentID); err == nil {
			speeds = sp
		}
	}

	snapshots := make(map[string]*scheduler.ConceptSnapshot, len(progress))
	for cid, p := range progress {
		reqStreak := 0
		timeThresh := 0.0
		if c := e.dag.Concept(cid); c != nil {
			reqStreak = c.MasteryThreshold.Streak
			timeThresh = c.MasteryThreshold.AvgTimeSeconds
		}
		ordinals := 0
		if banked != nil {
			ordinals = banked[cid]
		}
		speed := 1.0
		if speeds != nil && speeds[cid] != nil {
			speed = speeds[cid].LearningSpeed
		}
		isEncompassed := false
		if e.dag != nil && len(e.dag.EncompassedBy(cid)) > 0 {
			isEncompassed = true
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
			IsProvisional:  p.IsProvisional,
			BankedOrdinals: ordinals,
			Repetitions:    p.SM2Repetitions,
			Interval:       p.SM2Interval,
			EFactor:        p.SM2EFactor,
			LearningSpeed:  speed,
			Encompassed:    isEncompassed,
			FSRSStability:  p.FSRSStability,
			FSRSDifficulty: p.FSRSDifficulty,
			FSRSLastReview: p.FSRSLastReview,
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
				Encompassed:    len(e.dag.EncompassedBy(c.ID)) > 0,
			}
		}
	}
	return snapshots
}

// conceptSnapshots is the loading wrapper, for the question-serving paths.
func (e *Engine) conceptSnapshots(studentID string) (map[string]*scheduler.ConceptSnapshot, error) {
	if studentID != "" {
		_ = e.SatisfyImplicitReviews(studentID)
	}
	progress, err := e.repo.GetAllProgress(studentID)
	if err != nil {
		return nil, fmt.Errorf("load progress: %w", err)
	}
	return e.conceptSnapshotsFrom(studentID, progress), nil
}

// ReviewsDue counts concepts whose SM-2 interval has elapsed. Moved here from the server,
// which computed it inline and which the recommender needs it from.
func (e *Engine) ReviewsDue(studentID string) (int, error) {
	if studentID != "" {
		_ = e.SatisfyImplicitReviews(studentID)
	}
	progress, err := e.repo.GetAllProgress(studentID)
	if err != nil {
		return 0, fmt.Errorf("load progress: %w", err)
	}
	now := time.Now()
	count := 0
	for _, p := range progress {
		isDue := false
		if p.NextReviewDue != nil && !p.NextReviewDue.After(now) {
			isDue = true
		}
		if p.FSRSStability != nil && *p.FSRSStability > 0 {
			var elapsed float64
			if p.FSRSLastReview != nil && !p.FSRSLastReview.IsZero() {
				elapsed = now.Sub(*p.FSRSLastReview).Hours() / 24.0
			} else if p.LastReviewed != nil && !p.LastReviewed.IsZero() {
				elapsed = now.Sub(*p.LastReviewed).Hours() / 24.0
			}
			if elapsed < 0 {
				elapsed = 0
			}
			if fsrs.DefaultParameters().Retrievability(elapsed, *p.FSRSStability) < fsrs.DefaultRetention {
				isDue = true
			}
		}
		if isDue {
			count++
		}
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
	// effective MASTERED, decay applied, so a stale concept has left the count rather than
	// inflating it, and *not* "anything not unseen", which counted a concept the learner had
	// merely started. Only the new-learner gate consumes this today and AttemptedDays already
	// guarded it, so the miscount was unobservable; it is fixed because the field is named
	// after a fact it did not report, and the next caller inherits discrepancies rather than
	// reading the comment.
	mastered, attemptedDays := 0, map[string]bool{}
	for _, p := range progress {
		var elapsed float64
		if p.FSRSLastReview != nil && !p.FSRSLastReview.IsZero() {
			elapsed = daysSinceReview(p.FSRSLastReview)
		} else if p.LastReviewed != nil {
			elapsed = daysSinceReview(p.LastReviewed)
		}
		if mastery.EffectiveStatusFSRS(mastery.Status(p.Status), elapsed, p.FSRSStability, mastery.DecayDays) == mastery.StatusMastered {
			mastered++
		}
		if p.LastAttempted != nil {
			attemptedDays[p.LastAttempted.UTC().Format("2006-01-02")] = true
		}
	}

	remQueue, err := e.repo.GetRemedialQueue(studentID)
	if err != nil {
		return scheduler.RecommendationResponse{}, fmt.Errorf("read remedial queue: %w", err)
	}
	e.mu.Lock()
	if path := e.activePath[studentID]; len(path) > 0 {
		for cid := range snapshots {
			if !path[cid] {
				delete(snapshots, cid)
			}
		}
	}
	for _, as := range e.sessions {
		if as.studentID == studentID && as.halted {
			for _, cid := range as.remedialQueue {
				c := e.dag.Concept(as.conceptID)
				if c == nil {
					c = e.dag.Concept(as.lastConceptID)
				}
				if c != nil {
					isRemedialTarget := false
					for _, kp := range c.KeyPrerequisites {
						if kp == cid {
							isRemedialTarget = true
							break
						}
					}
					if !isRemedialTarget {
						for _, p := range e.dag.PrereqsOf(c.ID) {
							if p.ID == cid {
								isRemedialTarget = true
								break
							}
						}
					}
					if isRemedialTarget {
						found := false
						for _, existing := range remQueue {
							if existing == cid {
								found = true
								break
							}
						}
						if !found {
							remQueue = append(remQueue, cid)
						}
					}
				}
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

	// The most recent attempt names the concept just served. Passing it demotes that concept
	// below every other candidate, so the next card is not the one the learner finished a
	// moment ago, which on a thin concept is a repeat of the same situation.
	var lastConceptID string
	var lastAttempted time.Time
	for _, p := range progress {
		if p.LastAttempted != nil && p.LastAttempted.After(lastAttempted) {
			lastAttempted = *p.LastAttempted
			lastConceptID = p.ConceptID
		}
	}

	// Count what was already served today, so the candidate loop can drop a concept that has
	// reached its cap. Without this the cap in the scheduler is inert: it defends a limit no
	// caller supplies. A read failure leaves the map nil, which restores adjacency-only
	// behaviour rather than failing the recommendation, because a missing cap degrades variety
	// while a failed recommendation hides the next card entirely.
	servedToday := map[string]int{}
	if attempts, err := e.repo.GetAttemptsForStudent(studentID); err == nil {
		dayStart := time.Now().UTC().Truncate(24 * time.Hour)
		for _, a := range attempts {
			if a.Timestamp.UTC().Before(dayStart) {
				continue
			}
			// Placement answers are evidence about what the learner already knows, not
			// practice served from this concept, so they must not consume the practice
			// budget. Counting them made a freshly placed learner look capped out and the
			// recommendation fall back to the graph.
			if a.Source == "diagnostic" {
				continue
			}
			servedToday[a.ConceptID]++
		}
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
		KeyPrerequisites:    remQueue,
		RemedialConcepts:    remQueue,
		LastConceptID:       lastConceptID,
		ServedToday:         servedToday,
		ServedPerDay:        servedPerDayCap,
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
// It used to point at `/study` instead, with "Lessons are still open", true when it was
// written, and a dead end once the library closed: a learner who opened the app on a break was
// sent to a page that could only explain things they had already decided not to practise today.
// `/settings` is where the pause and its end both live, so it is the one destination that can
// actually answer "I am not doing reviews, now what".
func (e *Engine) RecommendPaused(studentID string) scheduler.RecommendationResponse {
	return scheduler.RecommendationResponse{
		Primary: &scheduler.Recommendation{
			ID:     "paused",
			Kind:   scheduler.KindLearn,
			Reason: scheduler.ReasonNew,
			Action: scheduler.Action{Type: scheduler.KindLearn, Href: "/settings"},
			Badge:  "Paused",
			Detail: "Reviews are on hold until you resume, due reviews are hidden while paused",
			CTA:    "Change or resume →",
		},
		Alternatives: []scheduler.Recommendation{},
		GeneratedAt:  time.Now().UTC(),
	}
}
