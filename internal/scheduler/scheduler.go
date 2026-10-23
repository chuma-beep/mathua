package scheduler

import (
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/fsrs"
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
	// Placement marks a concept the diagnostic established as the learner's
	// frontier. It grants reachability only: a placement-seeded prerequisite
	// satisfies its successors so the learner can start where they were placed.
	// It is not mastery and is never used to decide mastery, XP, the quiz gate or
	// review (ADR-008, ADR-050).
	Placement bool
	// IsProvisional marks a conditionally completed prerequisite from diagnostic
	// placement. If a learner struggles on a dependent concept, the prerequisite
	// demotes to LEARNING and blocks progression until relearned.
	IsProvisional bool

	// FIRe: Fractional Implicit Repetition fields.
	// BankedOrdinals is the count of whole implicit repetitions banked from downstream practice.
	BankedOrdinals int
	// Repetitions is the current count of successful reviews under SM-2.
	Repetitions int
	// Interval is the current SM-2 interval in days.
	Interval int
	// EFactor is the current SM-2 easiness factor.
	EFactor float64
	// LearningSpeed scales review interval calculations.
	LearningSpeed float64
	// Encompassed marks whether this concept is encompassed by other concepts in the curriculum.
	Encompassed bool

	// FSRS-6 fields.
	FSRSStability  *float64
	FSRSDifficulty *float64
	FSRSLastReview *time.Time
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
	if len(out) < 3 {
		for _, c := range candidates {
			if added[c.Concept.ID] {
				continue
			}
			out = append(out, NextConcept{Concept: c.Concept, IsReview: c.IsReview})
			added[c.Concept.ID] = true
			if len(out) >= 3 {
				break
			}
		}
	}
	return out
}

// buildCandidatesAdvanced adds PR 1.4 edges to candidate construction:
// interleaving exclusion (same Subdomain as last 2 tasks), layering +2, non-interference -3.
func buildCandidatesAdvanced(dag *concepts.DAG, snapshots map[string]*ConceptSnapshot, now time.Time, prevIDs []string, weakness map[string]float64) []Candidate {
	prevSet := make(map[string]bool, len(prevIDs))
	for _, id := range prevIDs {
		prevSet[id] = true
	}
	windowIDs := prevIDs
	if len(windowIDs) > 2 {
		windowIDs = windowIDs[len(windowIDs)-2:]
	}
	prevSubs := make(map[string]bool)
	prevGroups := make(map[string]bool)
	for _, id := range windowIDs {
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
		// Interleaving window: skip same subdomain as recent practice (last 2 tasks),
		// but never starve the session (fall back if the filter empties the set).
		if len(prevSubs) > 0 && prevSubs[c.Subdomain] {
			continue
		}
		priority := ComputePriorityWithContext(snapshots[c.ID], now, cand.IsReview, c, weakness, prevGroups)
		out = append(out, Candidate{Concept: c, IsReview: cand.IsReview, Priority: priority})
	}
	if len(out) == 0 {
		// Interleaving would starve the session — degrade to no-window candidates.
		out = base
	}
	return out
}

// ComputePriorityWithContext computes priority including layering bonus (+2 if encompasses
// a weak prerequisite) and non-interference penalty (-3 if sharing an interference group).
func ComputePriorityWithContext(
	snap *ConceptSnapshot,
	now time.Time,
	isDecaying bool,
	c *concepts.Concept,
	weakness map[string]float64,
	recentInterferenceGroups map[string]bool,
) float64 {
	p := computePriority(snap, now, isDecaying)
	if c != nil {
		p += layeringBonus(c, weakness)
		if c.InterferenceGroup != "" && recentInterferenceGroups != nil && recentInterferenceGroups[c.InterferenceGroup] {
			p -= 3.0
		}
	}
	return p
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

// AdvanceSnapshotInterval advances the review interval and moves NextReviewDue forward
// when banked implicit credit satisfies an explicit review. It consumes one banked ordinal.
func AdvanceSnapshotInterval(snap *ConceptSnapshot, now time.Time) int {
	if snap == nil {
		return 0
	}
	interval := snap.Interval
	if interval < 1 {
		interval = 1
	}
	efactor := snap.EFactor
	if efactor < 1.3 {
		efactor = 2.5
	}
	speed := snap.LearningSpeed
	if speed <= 0 {
		speed = 1.0
	}
	reps := snap.Repetitions
	if reps <= 0 {
		if interval >= 6 {
			reps = 2
		} else if interval >= 1 {
			reps = 1
		}
	}
	nextSM2 := ComputeSM2WithSpeed(SM2{Repetitions: reps, Interval: interval, EFactor: efactor}, 4, speed)
	stepDays := nextSM2.Interval
	if stepDays < 1 {
		stepDays = 1
	}

	if snap.FSRSStability != nil && *snap.FSRSStability > 0 {
		params := fsrs.DefaultParameters()
		diff := params.InitialDifficulty(fsrs.Good)
		if snap.FSRSDifficulty != nil && *snap.FSRSDifficulty > 0 {
			diff = *snap.FSRSDifficulty
		}
		var elapsed float64
		if snap.FSRSLastReview != nil && !snap.FSRSLastReview.IsZero() {
			elapsed = now.Sub(*snap.FSRSLastReview).Hours() / 24.0
		}
		if elapsed < 0 {
			elapsed = 0
		}
		r := params.Retrievability(elapsed, *snap.FSRSStability)
		nextStab := params.NextStability(*snap.FSRSStability, diff, r, fsrs.Good)
		snap.FSRSStability = &nextStab
		fsrsInterval := params.Interval(nextStab)
		if fsrsInterval > stepDays {
			stepDays = fsrsInterval
		}
		snap.FSRSLastReview = &now
	}

	baseline := now
	if snap.NextReviewDue != nil && !snap.NextReviewDue.IsZero() {
		baseline = *snap.NextReviewDue
	}
	nextDue := baseline.AddDate(0, 0, stepDays)
	floor := now.AddDate(0, 0, 1)
	if nextDue.Before(floor) {
		nextDue = floor
	}
	snap.NextReviewDue = &nextDue
	snap.LastReviewed = now
	snap.Interval = stepDays
	if snap.BankedOrdinals > 0 {
		snap.BankedOrdinals--
	}
	return stepDays
}

func buildCandidates(dag *concepts.DAG, snapshots map[string]*ConceptSnapshot, now time.Time, skipID string) []Candidate {
	var out []Candidate
	for _, c := range dag.Order() {
		if c.ID == skipID {
			continue
		}
		snap := snapshots[c.ID]
		effStatus, isReview := effectiveState(snap, now)
		isEncompassed := (snap != nil && snap.Encompassed) || (dag != nil && len(dag.EncompassedBy(c.ID)) > 0)
		if isReview && snap != nil && snap.BankedOrdinals > 0 && isEncompassed {
			isReview = false
			if effStatus == "DECAYING" {
				effStatus = mastery.StatusMastered
			}
			AdvanceSnapshotInterval(snap, now)
		}
		if effStatus == mastery.StatusMastered && !isReview {
			if snap == nil || snap.WeaknessScore < weaknessPracticeFloor {
				continue
			}
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
	status := mastery.EffectiveStatus(snap.Status, daysSince(snap.LastReviewed), mastery.DecayDays)
	isReview := status == "DECAYING"
	if snap.Status == mastery.StatusMastered && snap.NextReviewDue != nil && !snap.NextReviewDue.After(now) {
		isReview = true
	}
	if snap.Status == mastery.StatusMastered && snap.FSRSStability != nil && *snap.FSRSStability > 0 {
		var elapsedDays float64
		if snap.FSRSLastReview != nil && !snap.FSRSLastReview.IsZero() {
			elapsedDays = now.Sub(*snap.FSRSLastReview).Hours() / 24.0
		} else if !snap.LastReviewed.IsZero() {
			elapsedDays = now.Sub(snap.LastReviewed).Hours() / 24.0
		}
		if elapsedDays < 0 {
			elapsedDays = 0
		}
		r := fsrs.DefaultParameters().Retrievability(elapsedDays, *snap.FSRSStability)
		if r < fsrs.DefaultRetention {
			isReview = true
			if status == mastery.StatusMastered {
				status = "DECAYING"
			}
		}
	}
	return status, isReview
}

func prereqsMet(dag *concepts.DAG, c *concepts.Concept, snapshots map[string]*ConceptSnapshot, now time.Time) bool {
	for _, pid := range c.Prerequisites {
		snap := snapshots[pid]
		if snap == nil {
			return false
		}
		// A demoted provisional prerequisite blocks dependent progression.
		if snap.IsProvisional && snap.Status == mastery.StatusLearning {
			return false
		}
		// Placement-seeded prerequisites count as satisfied. The diagnostic
		// established the learner's frontier; refusing to let them start the
		// concept just past it would make placement inert. This is reachability,
		// not mastery — placement never appears in countsAsMastered, the quiz
		// gate or XP.
		if snap.Placement {
			continue
		}
		status := mastery.EffectiveStatus(snap.Status, daysSince(snap.LastReviewed), mastery.DecayDays)
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
		if snap == nil || snap.BankedOrdinals == 0 || !snap.Encompassed {
			priority += 5.0
		}
		if snap != nil && snap.FSRSStability != nil && *snap.FSRSStability > 0 {
			var elapsedDays float64
			if snap.FSRSLastReview != nil && !snap.FSRSLastReview.IsZero() {
				elapsedDays = now.Sub(*snap.FSRSLastReview).Hours() / 24.0
			} else if !snap.LastReviewed.IsZero() {
				elapsedDays = now.Sub(snap.LastReviewed).Hours() / 24.0
			}
			if elapsedDays < 0 {
				elapsedDays = 0
			}
			r := fsrs.DefaultParameters().Retrievability(elapsedDays, *snap.FSRSStability)
			priority += (1.0 - r) * 5.0
		}
	}
	if snap != nil && snap.Status == mastery.StatusMastered && weakness >= weaknessPracticeFloor {
		priority += 4.0
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
