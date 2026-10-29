package engine

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chuma-beep/mathua/internal/fire"
	"github.com/chuma-beep/mathua/internal/mastery"
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

// SatisfyImplicitReviews applies at most one banked ordinal per eligible target
// when FIRe is active and the learner is inside the activation cohort. Shadow
// remains inert: recording an observation must never move a due date. The
// application is durable, so a crash after commit cannot apply the same ordinal
// twice and a crash before commit leaves it banked for a safe retry.
func conflictReason(reason string) bool {
	switch reason {
	case "stale_baseline", "due_moved", "epoch_closed", "cohort_mismatch", "activation_unavailable":
		return true
	default:
		return false
	}
}

func (e *Engine) SatisfyImplicitReviews(studentID string) error {
	if studentID == "" || e.fireMode != fire.Active || !e.fireConfig.ActiveFor(studentID) {
		return nil
	}
	ordinals, ok := e.repo.(storage.CreditOrdinalStore)
	if !ok {
		return nil
	}
	applications, ok := e.repo.(storage.CreditApplyStore)
	if !ok {
		return nil
	}
	details, err := ordinals.BankedOrdinalDetails(studentID)
	if err != nil {
		return fmt.Errorf("load banked ordinals: %w", err)
	}
	// The store returns every banked ordinal. Filter to targets that can
	// actually accept credit before opening a transaction, so a large backlog
	// of future targets cannot starve a due one or make every question-serving
	// call pay for transactions that can only answer not_due.
	now := nowUTC()
	seenTargets := make(map[string]bool)
	applied := 0
	for _, o := range details {
		if !creditTargetReady(o, now) || seenTargets[o.ConceptID] {
			continue
		}
		seenTargets[o.ConceptID] = true
		req := storage.CreditApplyRequest{
			StudentID:       studentID,
			TargetConceptID: o.ConceptID,
			SourceConceptID: "",
			SourceAttemptID: fmt.Sprintf("fire-review:%s:%s:%d:%d", studentID, o.ConceptID, o.Epoch, o.Ordinal),
			Epoch:           o.Epoch,
			Ordinal:         o.Ordinal,
			Cohort:          e.fireConfig.Cohort,
			Now:             now,
		}
		start := time.Now()
		res, err := applications.ApplyCredit(req)
		e.fireStats.LatencyMicros.Add(time.Since(start).Microseconds())
		e.fireStats.LatencyCount.Add(1)
		if err != nil {
			e.fireStats.Failures.Add(1)
			log.Printf("warning: apply implicit credit %s/%s: %v", studentID, o.ConceptID, err)
			continue
		}
		switch {
		case res.Replayed:
			e.fireStats.Replays.Add(1)
		case res.Applied:
			e.fireStats.Applications.Add(1)
			e.fireStats.DueShifts.Add(1)
			applied++
			log.Printf("fire: applied ordinal %d to %s/%s, due %s -> %s", o.Ordinal, studentID, o.ConceptID, res.ActualDue, res.ProposedDue)
		case conflictReason(res.Reason):
			e.fireStats.Conflicts.Add(1)
			log.Printf("fire: conflict on ordinal %d for %s/%s: %s", o.Ordinal, studentID, o.ConceptID, res.Reason)
		case res.Reason != "":
			// Waiting states such as not_due are expected while credit is
			// banked and must not feed the conflict alert.
			log.Printf("fire: deferred ordinal %d for %s/%s: %s", o.Ordinal, studentID, o.ConceptID, res.Reason)
		}
		if applied >= 32 {
			break
		}
	}
	return nil
}

// creditTargetReady is only a prefilter for SatisfyImplicitReviews. It mirrors
// the cheap, time-dependent half of creditReasonEligible. ApplyCredit checks
// the authoritative mode, activation, epoch and due baseline in its
// transaction.
func creditTargetReady(o storage.CreditOrdinalDetail, now time.Time) bool {
	if !o.HasProgress || o.PlacementSeeded != 0 || o.IsProvisional != 0 {
		return false
	}
	switch strings.ToUpper(strings.TrimSpace(o.TargetStatus)) {
	case "MASTERED", "DECAYING":
	default:
		return false
	}
	due, ok := parseCreditDue(o.TargetDue)
	return ok && !due.After(now)
}

func parseCreditDue(raw string) (time.Time, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return time.Time{}, false
	}
	for _, layout := range []string{time.RFC3339Nano, time.RFC3339, "2006-01-02 15:04:05"} {
		if t, err := time.Parse(layout, raw); err == nil {
			return t.UTC(), true
		}
	}
	return time.Time{}, false
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
	var speeds map[string]*storage.TopicSpeed
	if studentID != "" {
		if sp, err := e.repo.GetAllTopicSpeeds(studentID); err == nil {
			speeds = sp
		}
	}

	banked := map[string]int{}
	if studentID != "" {
		if cs, ok := e.repo.(storage.CreditStore); ok {
			if got, err := cs.GetBankedOrdinals(studentID); err == nil {
				banked = got
			}
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
		speed := 1.0
		if speeds != nil && speeds[cid] != nil {
			speed = speeds[cid].LearningSpeed
		}
		encompassed := false
		if c := e.dag.Concept(cid); c != nil && len(c.Encompasses) > 0 {
			encompassed = true
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
			Repetitions:    p.SM2Repetitions,
			Interval:       p.SM2Interval,
			EFactor:        p.SM2EFactor,
			LearningSpeed:  speed,
			FSRSStability:  p.FSRSStability,
			FSRSDifficulty: p.FSRSDifficulty,
			FSRSLastReview: p.FSRSLastReview,
			BankedOrdinals: banked[cid],
			Encompassed:    encompassed,
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
				Encompassed:    len(c.Encompasses) > 0,
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
	// Count only what NextReview can actually serve. Counting every stored row
	// with a stale NextReviewDue made a learning concept look like a retention
	// review, so the profile advertised work and /review answered "All caught
	// up". The scheduler owns the definition, and DueReviewCount reuses it.
	snapshots := e.conceptSnapshotsFrom(studentID, progress)
	return scheduler.New(e.dag).DueReviewCount(snapshots, time.Now().UTC()), nil
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
	if quizDue && len(e.quizCandidateWeakness(studentID, progress)) == 0 {
		// The XP gate is earned, but there is nothing a routine check could
		// assess. Recommending one would only lead to the 422 no-candidate
		// response (report 291).
		quizDue = false
	}

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
