package scheduler

import (
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/mastery"
)

type Candidate struct {
	Concept  *concepts.Concept
	IsReview bool
	Priority float64
}

type Scheduler struct {
	dag *concepts.DAG
}

func New(dag *concepts.DAG) *Scheduler {
	return &Scheduler{dag: dag}
}

type NextConcept struct {
	Concept  *concepts.Concept
	IsReview bool
}

// ConceptSnapshot tracks a student's progress on one concept.
// Zero value means the concept has never been attempted.
type ConceptSnapshot struct {
	Status         mastery.Status
	Streak         int
	LastAttempted  time.Time
	LastReviewed   time.Time
	NextReviewDue  *time.Time
	RequiredStreak int
	TimeThreshold  float64
	WeaknessScore  float64
}

// Next picks the next concept for a student to work on.
//
// prevConceptID is the last concept practiced ("" if none).
// sessionReviewCount / sessionNewCount are used to maintain the 70/30 balance.
func (s *Scheduler) Next(
	snapshots map[string]*ConceptSnapshot,
	prevConceptID string,
	sessionReviewCount, sessionNewCount int,
) *NextConcept {
	now := time.Now().UTC()

	candidates := buildCandidates(s.dag, snapshots, now, prevConceptID)
	if len(candidates) == 0 && prevConceptID != "" {
		candidates = buildCandidates(s.dag, snapshots, now, "")
	}
	if len(candidates) == 0 {
		return nil
	}

	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].Priority != candidates[j].Priority {
			return candidates[i].Priority > candidates[j].Priority
		}
		return candidates[i].Concept.ID < candidates[j].Concept.ID
	})

	pick := selectWithBalance(candidates, sessionReviewCount, sessionNewCount)
	return &NextConcept{Concept: pick.Concept, IsReview: pick.IsReview}
}

func buildCandidates(dag *concepts.DAG, snapshots map[string]*ConceptSnapshot, now time.Time, skipID string) []Candidate {
	var out []Candidate
	for _, c := range dag.Order() {
		if c.ID == skipID {
			continue
		}
		snap := snapshots[c.ID]
		effStatus, isReview := effectiveState(snap, now)
		if effStatus == mastery.StatusMastered {
			continue
		}
		if !prereqsMet(dag, c, snapshots, now) {
			continue
		}
		priority := computePriority(snap, now, isReview)
		out = append(out, Candidate{Concept: c, IsReview: isReview, Priority: priority})
	}
	return out
}

func effectiveState(snap *ConceptSnapshot, now time.Time) (mastery.Status, bool) {
	if snap == nil {
		return mastery.StatusUnseen, false
	}
	status := mastery.EffectiveStatus(snap.Status, daysSince(snap.LastReviewed), 14)
	isReview := status == "DECAYING"
	if snap.Status == mastery.StatusMastered && snap.NextReviewDue != nil && !snap.NextReviewDue.After(now) {
		isReview = true
	}
	return status, isReview
}

func prereqsMet(dag *concepts.DAG, c *concepts.Concept, snapshots map[string]*ConceptSnapshot, now time.Time) bool {
	for _, pid := range c.Prerequisites {
		snap := snapshots[pid]
		if snap == nil {
			return false
		}
		status := mastery.EffectiveStatus(snap.Status, daysSince(snap.LastReviewed), 14)
		if status != mastery.StatusMastered {
			return false
		}
	}
	return true
}

func computePriority(snap *ConceptSnapshot, now time.Time, isDecaying bool) float64 {
	daysSinceLastSeen := 0.0
	masteryScore := 0.0
	weakness := 0.0
	if snap != nil {
		if !snap.LastAttempted.IsZero() {
			daysSinceLastSeen = now.Sub(snap.LastAttempted).Hours() / 24
		}
		if snap.RequiredStreak > 0 {
			masteryScore = math.Min(1.0, float64(snap.Streak)/float64(snap.RequiredStreak))
		}
		weakness = snap.WeaknessScore
	}
	priority := 0.5*daysSinceLastSeen + 0.2*(1.0-masteryScore) + 0.3*weakness
	if isDecaying {
		priority += 5.0
	}
	return priority
}

func selectWithBalance(candidates []Candidate, reviewCount, newCount int) Candidate {
	targetReviewRatio := 0.3
	total := reviewCount + newCount + 1
	if total > 2 {
		currentRatio := float64(reviewCount) / float64(total)
		if currentRatio < targetReviewRatio && !candidates[0].IsReview {
			for _, cand := range candidates {
				if cand.IsReview {
					return cand
				}
			}
		}
	}
	// Shuffle among equal-priority candidates.
	topPriority := candidates[0].Priority
	same := 1
	for same < len(candidates) && candidates[same].Priority == topPriority {
		same++
	}
	return candidates[rand.Intn(same)]
}

func daysSince(t time.Time) float64 {
	if t.IsZero() {
		return 999
	}
	return time.Since(t).Hours() / 24
}
