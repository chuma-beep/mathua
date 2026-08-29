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

// NextSmart returns up to 3 dissimilar candidates (distinct subdomain),
// PR 1.4: layering bonus (+2 when this concept encompasses a weak prereq),
// non-interference penalty (-3 when sharing an InterferenceGroup with a
// recently practiced concept), interleaving window (exclude same Subdomain
// as the last 2 practiced). The 70/30 balance still governs the top pick.
func (s *Scheduler) NextSmart(
	snapshots map[string]*ConceptSnapshot,
	prevConceptIDs []string,
	reviewCount, newCount int,
	weakness map[string]float64,
) []NextConcept {
	now := time.Now().UTC()
	candidates := buildCandidatesAdvanced(s.dag, snapshots, now, prevConceptIDs, weakness)
	if len(candidates) == 0 {
		return nil
	}
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].Priority != candidates[j].Priority {
			return candidates[i].Priority > candidates[j].Priority
		}
		return candidates[i].Concept.ID < candidates[j].Concept.ID
	})
	pref := selectWithBalance(candidates, reviewCount, newCount)

	out := []NextConcept{{Concept: pref.Concept, IsReview: pref.IsReview}}
	added := map[string]bool{pref.Concept.ID: true}
	seenSub := map[string]bool{pref.Concept.Subdomain: true}
	for _, c := range candidates {
		if added[c.Concept.ID] || seenSub[c.Concept.Subdomain] {
			continue
		}
		out = append(out, NextConcept{Concept: c.Concept, IsReview: c.IsReview})
		added[c.Concept.ID] = true
		seenSub[c.Concept.Subdomain] = true
		if len(out) >= 3 {
			break
		}
	}
	return out
}

// buildCandidatesAdvanced adds PR 1.4 edges to candidate construction:
// interleaving exclusion (same Subdomain as last 2), layering +2, non-interference -3.
func buildCandidatesAdvanced(dag *concepts.DAG, snapshots map[string]*ConceptSnapshot, now time.Time, prevIDs []string, weakness map[string]float64) []Candidate {
	prevSet := make(map[string]bool, len(prevIDs))
	prevSubs := make(map[string]bool)
	prevGroups := make(map[string]bool)
	for _, id := range prevIDs {
		prevSet[id] = true
		if c := dag.Concept(id); c != nil {
			if c.Subdomain != "" {
				prevSubs[c.Subdomain] = true
			}
			if c.InterferenceGroup != "" {
				prevGroups[c.InterferenceGroup] = true
			}
		}
	}
	base := buildCandidates(dag, snapshots, now, "")
	var out []Candidate
	for _, cand := range base {
		c := cand.Concept
		if prevSet[c.ID] {
			continue
		}
		// Interleaving window: skip same subdomain as recent practice, but never
		// starve the session (fall back if the filter empties the set).
		if len(prevSubs) > 0 && prevSubs[c.Subdomain] {
			continue
		}
		priority := cand.Priority
		priority += layeringBonus(c, weakness)
		if prevGroups[c.InterferenceGroup] {
			priority -= 3.0
		}
		out = append(out, Candidate{Concept: c, IsReview: cand.IsReview, Priority: priority})
	}
	if len(out) == 0 {
		// Interleaving would starve the session — degrade to no-window candidates.
		out = base
	}
	return out
}

// layeringBonus: +2 when this concept encompasses a weak prerequisite
// (retroactive facilitation, improve.md:34).
func layeringBonus(c *concepts.Concept, weakness map[string]float64) float64 {
	for _, eid := range c.Encompasses {
		if weakness[eid] >= 0.6 {
			return 2.0
		}
	}
	return 0
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
		if status != mastery.StatusMastered && status != "DECAYING" {
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

// NextReview picks the next due-for-review concept.
// It only considers candidates whose IsReview flag is true.
func (s *Scheduler) NextReview(
	snapshots map[string]*ConceptSnapshot,
	prevConceptID string,
) *NextConcept {
	now := time.Now().UTC()
	candidates := buildCandidates(s.dag, snapshots, now, prevConceptID)
	if len(candidates) == 0 && prevConceptID != "" {
		candidates = buildCandidates(s.dag, snapshots, now, "")
	}
	var reviews []Candidate
	for _, c := range candidates {
		if c.IsReview {
			reviews = append(reviews, c)
		}
	}
	if len(reviews) == 0 {
		return nil
	}
	sort.Slice(reviews, func(i, j int) bool {
		if reviews[i].Priority != reviews[j].Priority {
			return reviews[i].Priority > reviews[j].Priority
		}
		return reviews[i].Concept.ID < reviews[j].Concept.ID
	})
	topPriority := reviews[0].Priority
	same := 1
	for same < len(reviews) && reviews[same].Priority == topPriority {
		same++
	}
	// Shuffle among equal-priority candidates with per-call seeded rand.
	r := rand.New(rand.NewSource(time.Now().UnixNano() ^ int64(same)))
	return &NextConcept{Concept: reviews[r.Intn(same)].Concept, IsReview: true}
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
	// Shuffle among equal-priority candidates with per-call seeded rand.
	topPriority := candidates[0].Priority
	same := 1
	for same < len(candidates) && candidates[same].Priority == topPriority {
		same++
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano() ^ int64(len(candidates)) ^ int64(reviewCount+newCount)))
	return candidates[r.Intn(same)]
}

func daysSince(t time.Time) float64 {
	if t.IsZero() {
		return 999
	}
	return time.Since(t).Hours() / 24
}
