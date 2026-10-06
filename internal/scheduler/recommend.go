package scheduler

import (
	"fmt"
	"sort"
	"time"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/mastery"
	"github.com/chuma-beep/mathua/internal/xp"
)

// The learner-facing recommendation.
//
// This type is the answer to a question that used to be answered twice and disagreed: "what
// should this learner do next?". Until now the client recomputed the whole decision in
// `lib/nextUp.ts` from five API calls plus a bundled copy of the corpus, while
// `internal/scheduler` — which holds the eligibility rules, the decay rule and the priority
// weights — was unreachable from `/learn` and `/profile`. The two had already drifted: the
// client still halved review XP, a rule ADR-020 removed from the server, and it ranked by
// weakness-then-concept-id while the server ranked by a weighted score. So the engine answers
// once, here, and both surfaces render it.
//
// Eligibility, decay and priority reuse buildCandidates, effectiveState, prereqsMet and
// computePriority rather than restating them. A second copy of `prereqsMet` is exactly how the
// client and server came to disagree about DECAYING in the first place.

// Kind is what kind of work the task is.
type Kind string

const (
	// KindLearn is a concept the learner has not finished.
	KindLearn Kind = "learn"
	// KindPractice is more work on a concept they are struggling with.
	KindPractice Kind = "practice"
	// KindReview is a retrieval check on something already learned.
	KindReview Kind = "review"
	// KindMasteryCheck is the XP-gated quiz.
	KindMasteryCheck Kind = "mastery_check"
)

// Reason is why the task was chosen, in language a learner can act on. None of these name the
// graph, the frontier or the scheduler: "frontier concept — learn it next" told a learner
// nothing they could do, where "new · ready to learn" tells them what pressing the button
// will start.
type Reason string

const (
	ReasonNew           Reason = "new"
	ReasonReviewDue     Reason = "review_due"
	ReasonNeedsPractice Reason = "needs_practice"
	ReasonPrerequisite  Reason = "prerequisite"
	ReasonMasteryReady  Reason = "mastery_ready"
)

// Action is where the task lives and what pressing it does.
type Action struct {
	Type Kind   `json:"type"`
	Href string `json:"href"`
}

// bucket is how a recommendation is selected. It is finer than Kind because "new" and
// "resume" are both KindLearn on the wire, and the caps are per bucket.
type bucket string

const (
	bucketReview   bucket = "review"
	bucketPractice bucket = "practice"
	bucketNew      bucket = "new"
	bucketResume   bucket = "resume"
)

// Recommendation is one worthwhile task.
type Recommendation struct {
	ID           string  `json:"id"`
	ConceptID    string  `json:"conceptId"`
	ConceptTitle string  `json:"conceptTitle"`
	Kind         Kind    `json:"kind"`
	Reason       Reason  `json:"reason"`
	Priority     float64 `json:"priority"`
	// EstimatedMinutes comes from the concept's own avg_time_seconds, so a learner is not
	// shown a duration invented by the client.
	EstimatedMinutes *float64 `json:"estimatedMinutes,omitempty"`
	XP               *int     `json:"xp,omitempty"`
	Action           Action   `json:"action"`

	// Presentation, owned here so /profile and /learn cannot word the same task differently.
	// That was the duplication defect: one card said "Continue: Tens and units" and the card
	// above it said "Currently working towards Tens and units", both derived from one head.
	Badge  string `json:"badge"`
	Detail string `json:"detail"`
	CTA    string `json:"cta"`

	bucket bucket
}

// RecommendationResponse is the whole answer: one dominant task, and the alternatives the
// learner may pick instead. "The engine chooses; the learner disposes."
type RecommendationResponse struct {
	Primary      *Recommendation  `json:"primary"`
	Alternatives []Recommendation `json:"alternatives"`
	GeneratedAt  time.Time        `json:"generatedAt"`
}

// RecommendInput is everything the decision needs. It is deliberately the same information the
// client was assembling, so moving the decision server-side changes where it runs and not what
// it knows.
type RecommendInput struct {
	Snapshots map[string]*ConceptSnapshot
	Now       time.Time

	// Weakness is 0..1 per concept, from Engine.WeaknessMap.
	Weakness map[string]float64
	// DueReviews is the count of concepts whose SM-2 interval has elapsed.
	DueReviews int
	// QuizDue is whether the XP gate has been met since the last completed quiz.
	QuizDue bool
	// ConceptsMastered and AttemptedDays decide whether this learner is new enough to be
	// offered the diagnostic rather than a concept.
	ConceptsMastered int
	AttemptedDays    int
	// ExcludeConceptIds drops concepts already shown — the finished concept must not be offered
	// back as the next thing to learn.
	ExcludeConceptIds []string
}

// RecommendCaps bound the set. Five is the number of cards the existing shelf shows, so the
// learner-facing shape does not change; it is a cap and not a quota, and the set is allowed to
// be shorter. Manufacturing five cards for a learner who genuinely has two worthwhile tasks is
// how a recommendation turns back into a queue.
const (
	MaxRecommendations    = 5
	MaxConceptNew         = 2
	MaxConceptWeakness    = 1
	MaxConceptResume      = 1
	weaknessPracticeFloor = 0.35
)

// Recommend produces the learner's head and their alternatives.
//
// Diversity is not forced. Reviews come first because a retrieval check on something learned is
// worth more than a new concept, then practice on the weakest available concept, then new
// concepts, then resuming. If there is one worthwhile task, the learner gets one.
func (s *Scheduler) Recommend(in RecommendInput) RecommendationResponse {
	now := in.Now
	if now.IsZero() {
		now = time.Now()
	}
	resp := RecommendationResponse{Alternatives: []Recommendation{}, GeneratedAt: now}

	// A learner with nothing behind them cannot be recommended a concept, because no concept
	// has been shown to be learnable for them. The diagnostic is the honest answer.
	if in.ConceptsMastered == 0 && in.AttemptedDays == 0 {
		resp.Primary = &Recommendation{
			ID:     "diagnostic",
			Kind:   KindLearn,
			Reason: ReasonNew,
			Action: Action{Type: KindLearn, Href: "/onboard"},
			Badge:  "Recommended",
			Detail: "A short adaptive test finds where to start",
			CTA:    "Start test →",
		}
		// Everything available is still offered as alternatives, so the learner has agency
		// rather than a single unskippable command.
		for _, r := range s.conceptRecommendations(in, now, 4) {
			resp.Alternatives = append(resp.Alternatives, r)
		}
		return resp
	}

	picked := s.selectTasks(in, now)

	// The XP gate is a task, not a banner. It becomes the head when it is the most valuable
	// thing available, and an alternative otherwise — the point of surfacing it here is that a
	// learner is told why taking it helps.
	if in.QuizDue {
		mc := Recommendation{
			ID:       "mastery-check",
			Kind:     KindMasteryCheck,
			Reason:   ReasonMasteryReady,
			Action:   Action{Type: KindMasteryCheck, Href: "/goals?quiz=1"},
			Badge:    "Mastery check",
			Detail:   fmt.Sprintf("You have earned %d XP since your last one — take it to lock in what you have learned", xp.QuizGateXP),
			CTA:      "Take the check →",
			Priority: 100,
		}
		// A due retrieval check outranks the quiz: it is worth more to the learner, and the
		// quiz is still surfaced. Only when nothing is due does the quiz lead.
		if len(picked) > 0 && picked[0].Kind == KindReview {
			resp.Alternatives = append(resp.Alternatives, mc)
		} else {
			picked = append([]Recommendation{mc}, picked...)
		}
	}

	if len(picked) > 0 {
		resp.Primary = &picked[0]
		resp.Alternatives = append(resp.Alternatives, picked[1:]...)
	}
	if resp.Primary == nil {
		resp.Primary = &Recommendation{
			ID:     "study",
			Kind:   KindLearn,
			Reason: ReasonNew,
			Action: Action{Type: KindLearn, Href: "/study"},
			Badge:  "Library",
			Detail: "Everything available is already learned — the reference library is open",
			CTA:    "Browse lessons →",
		}
	}
	if resp.Alternatives == nil {
		resp.Alternatives = []Recommendation{}
	}
	return resp
}

// selectTasks ranks the available concepts and takes a deliberately mixed, short set.
func (s *Scheduler) selectTasks(in RecommendInput, now time.Time) []Recommendation {
	all := s.conceptRecommendations(in, now, 0) // 0 = no extra cap
	if len(all) == 0 {
		return nil
	}
	// Selection buckets are finer than the wire Kind: "new" and "resume" are both KindLearn,
	// so taking twice from KindLearn served the same concept twice. Bucketing separately is
	// what makes the caps mean what they say.
	buckets := map[bucket][]Recommendation{}
	for _, r := range all {
		buckets[r.bucket] = append(buckets[r.bucket], r)
	}

	var out []Recommendation
	seen := map[string]bool{}
	take := func(b bucket, n int) {
		for _, r := range buckets[b] {
			if n <= 0 {
				return
			}
			if seen[r.ConceptID] {
				continue
			}
			seen[r.ConceptID] = true
			out = append(out, r)
			n--
		}
	}

	// A due review outranks new work. One card, not one per concept: the review route runs
	// its own queue, so listing five concepts that all link to /review is five copies of one
	// choice.
	if in.DueReviews > 0 {
		take(bucketReview, 1)
	}
	take(bucketPractice, MaxConceptWeakness)
	take(bucketNew, MaxConceptNew)
	take(bucketResume, MaxConceptResume)

	if len(out) > MaxRecommendations {
		out = out[:MaxRecommendations]
	}
	return out
}

// conceptRecommendations scores every concept the learner is allowed to work on and renders
// each as a Recommendation, most valuable first.
func (s *Scheduler) conceptRecommendations(in RecommendInput, now time.Time, cap int) []Recommendation {
	excluded := make(map[string]bool, len(in.ExcludeConceptIds))
	for _, id := range in.ExcludeConceptIds {
		excluded[id] = true
	}

	cands := buildCandidates(s.dag, in.Snapshots, now, "")
	ranked := make([]Recommendation, 0, len(cands))
	for _, c := range cands {
		if excluded[c.Concept.ID] {
			continue
		}
		b, kind, reason := classify(c, in)
		if kind == "" {
			continue
		}
		priority := c.Priority
		if c.Concept.Subdomain != "" {
			// A small subdomain bonus so interleaved practice is preferred over two adjacent
			// nodes in one topic. This is NextSmart's interleaving rule, expressed as a
			// ranking term rather than a second selection pass.
			priority += layeringBonus(c.Concept, in.Weakness)
		}
		r := s.render(c.Concept, kind, reason, priority)
		r.bucket = b
		ranked = append(ranked, r)
	}
	sort.SliceStable(ranked, func(i, j int) bool {
		if ranked[i].Priority != ranked[j].Priority {
			return ranked[i].Priority > ranked[j].Priority
		}
		return ranked[i].ConceptID < ranked[j].ConceptID
	})
	if cap > 0 && len(ranked) > cap {
		ranked = ranked[:cap]
	}
	return ranked
}

// classify decides what kind of work a candidate is, which is the part the client used to do
// with string comparisons on a status field.
func classify(c Candidate, in RecommendInput) (bucket, Kind, Reason) {
	snap := in.Snapshots[c.Concept.ID]
	started := snap != nil
	if c.IsReview {
		return bucketReview, KindReview, ReasonReviewDue
	}
	w := in.Weakness[c.Concept.ID]
	if w >= weaknessPracticeFloor {
		return bucketPractice, KindPractice, ReasonNeedsPractice
	}
	if started {
		return bucketResume, KindLearn, ReasonPrerequisite
	}
	return bucketNew, KindLearn, ReasonNew
}

func (s *Scheduler) render(c *concepts.Concept, kind Kind, reason Reason, priority float64) Recommendation {
	r := Recommendation{
		ID:           c.ID,
		ConceptID:    c.ID,
		ConceptTitle: c.Label,
		Kind:         kind,
		Reason:       reason,
		Priority:     priority,
		Action:       Action{Type: kind, Href: "/learn?concept=" + c.ID},
	}
	if c.MasteryThreshold.AvgTimeSeconds > 0 {
		m := c.MasteryThreshold.AvgTimeSeconds / 60
		r.EstimatedMinutes = &m
	}
	// XP comes from internal/xp, the one award path. The client kept its own copy and it had
	// already drifted, still halving reviews after ADR-020 removed that from the server.
	if n := xp.EffortBase(c.MasteryThreshold.AvgTimeSeconds, xp.BaseDifficulty(c.GradingType)); n > 0 {
		v := n
		r.XP = &v
	}

	switch kind {
	case KindReview:
		r.Action.Href = "/review"
		r.Badge = "Review"
		r.Detail = "Something you have learned is due a check"
		r.CTA = "Review →"
	case KindPractice:
		r.Badge = "Practice"
		r.Detail = "Worth another go — this one has not stuck yet"
		r.CTA = "Keep practicing →"
	default:
		r.Badge = "New"
		r.Detail = "Ready to learn — everything it needs, you have done"
		r.CTA = "Start →"
	}
	return r
}

// EffectiveStatusOf reports the mastery status a learner should be shown, with decay applied,
// and whether the concept is currently due a retrieval check. Both halves are returned because
// callers that only need the status still have to write the tuple out, and the two have drifted
// apart before: DECAYING satisfies a prerequisite but must not be offered as teaching.
func EffectiveStatusOf(snap *ConceptSnapshot, now time.Time) (mastery.Status, bool) {
	return effectiveState(snap, now)
}
