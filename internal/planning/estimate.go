package planning

import (
	"math"
	"sort"
	"strings"
	"time"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/mastery"
	"github.com/chuma-beep/mathua/internal/storage"
	"github.com/chuma-beep/mathua/internal/xp"
)

// AttemptStats summarizes a learner's decisive history per concept.
// Don't-knows (admitted unknowns) are tracked separately: they cost reduced
// traversal time but must never reduce remaining work (anti-gaming floor).
type AttemptStats struct {
	Correct  int
	Failed   int
	DontKnow int
}

// EstimateOpts tunes the workload estimator. DailyXPRate is expected earned
// XP per active day (measured pace, not the aspirational goal).
type EstimateOpts struct {
	DailyXPRate     float64
	RestDaysPerWeek int
	DiagnosticMin   float64
	DecayDays       float64
	DeadlineDays    int // -1 or 0 = no deadline (plan-by-effort mode)
	FeasibleCap     float64
}

const (
	defaultDecayDays = 30
	minAccuracy      = 0.2
	maxAccuracy      = 0.95
	quizGateXP       = 150
	quizQuestions    = 5
	quizSecPerQ      = 30.0
	dontKnowTimeFrac = 0.3
)

// ConceptCost is the remaining price of one concept, both tracks.
type ConceptCost struct {
	ConceptID string  `json:"concept_id"`
	XP        float64 `json:"xp"`
	TimeMin   float64 `json:"time_min"`
	Attempts  float64 `json:"attempts"`
	Review    bool    `json:"review"`
}

// WorkloadEstimate is the full two-track estimate for a path.
type WorkloadEstimate struct {
	Total            int           `json:"total"`
	Mastered         int           `json:"mastered"`
	Remaining        []string      `json:"remaining"`
	Costs            []ConceptCost `json:"costs"`
	XPRemaining      float64       `json:"xp_remaining"`
	TimeMinRemaining float64       `json:"time_min_remaining"`
	ReviewsDue       int           `json:"reviews_due"`
	QuizzesAhead     int           `json:"quizzes_ahead"`
	AssessmentMin    float64       `json:"assessment_min"`
	DiagnosticMin    float64       `json:"diagnostic_min"`
	// Days and FinishDate answer plan-by-effort; RequiredPerDay/Feasible
	// answer plan-by-deadline (zero values when no deadline given).
	Days           float64 `json:"days"`
	FinishDate     string  `json:"finish_date"`
	RequiredPerDay float64 `json:"required_per_day"`
	Feasible       bool    `json:"feasible"`
	// Lines split the workload per the planning contract.
	Lines struct {
		Learning   float64 `json:"learning"`
		Assessment float64 `json:"assessment"`
		Diagnostic float64 `json:"diagnostic"`
	} `json:"lines"`
}

// taskBaseFor maps a concept to its XP task base: word problems are
// multistep everywhere (ADR-002), everything else is a plain lesson.
func taskBaseFor(conceptID string) float64 {
	if strings.HasSuffix(conceptID, ".word") {
		return float64(xp.BaseXP(xp.TaskMultistep))
	}
	return float64(xp.BaseXP(xp.TaskLesson))
}

// accuracy returns the decisive correct rate, excluding don't-knows, with a
// 0.5 prior for unseen concepts and a floor so estimates stay finite.
func accuracy(s AttemptStats) float64 {
	decisive := s.Correct + s.Failed
	if decisive == 0 {
		return 0.5
	}
	p := float64(s.Correct) / float64(decisive)
	if p < minAccuracy {
		p = minAccuracy
	}
	if p > maxAccuracy {
		p = maxAccuracy
	}
	return p
}

// thresholdSec returns the concept's expected seconds per attempt.
func thresholdSec(c *concepts.Concept) float64 {
	if c != nil && c.MasteryThreshold.AvgTimeSeconds > 0 {
		return c.MasteryThreshold.AvgTimeSeconds
	}
	return 10.0
}

// EstimateWorkload prices the remaining path on two tracks: XP remaining
// (the earned currency, drives dates) and focused minutes remaining.
// Pure and deterministic: same inputs always yield the same estimate.
func EstimateWorkload(path *Path, progress map[string]*storage.ConceptProgress, stats map[string]AttemptStats, opts EstimateOpts) *WorkloadEstimate {
	est := &WorkloadEstimate{}
	if path == nil {
		return est
	}
	decay := opts.DecayDays
	if decay <= 0 {
		decay = defaultDecayDays
	}
	now := time.Now().UTC()
	est.DiagnosticMin = opts.DiagnosticMin

	for _, c := range path.Concepts {
		est.Total++
		p := progress[c.ID]
		status := mastery.StatusUnseen
		if p != nil {
			status = mastery.Status(p.Status)
		}
		if status == mastery.StatusMastered && !reviewDue(p, now, decay) {
			est.Mastered++
			continue
		}
		if status == mastery.StatusMastered {
			// Mastered but decayed: one review pass, not re-learning.
			cost := reviewCost(c, p, stats[c.ID])
			cost.Review = true
			est.Costs = append(est.Costs, cost)
			est.Remaining = append(est.Remaining, c.ID)
			est.ReviewsDue++
			est.XPRemaining += cost.XP
			est.TimeMinRemaining += cost.TimeMin
			continue
		}
		cost := learningCost(c, p, stats[c.ID])
		est.Costs = append(est.Costs, cost)
		est.Remaining = append(est.Remaining, c.ID)
		est.XPRemaining += cost.XP
		est.TimeMinRemaining += cost.TimeMin
	}

	// Assessment line: one 150 XP-gate quiz per 150 XP of learning effort.
	// Quizzes grant XP, so they are time cost only — never XP debt.
	if est.XPRemaining > 0 {
		est.QuizzesAhead = int(math.Ceil(est.XPRemaining / quizGateXP))
	}
	est.AssessmentMin = float64(est.QuizzesAhead*quizQuestions) * quizSecPerQ / 60.0
	est.Lines.Learning = est.XPRemaining
	est.Lines.Assessment = 0 // assessments cost time, not XP debt
	est.Lines.Diagnostic = 0 // diagnostics cost time, not XP debt

	rate := effectiveRate(opts.DailyXPRate, opts.RestDaysPerWeek)
	if rate > 0 && est.XPRemaining > 0 {
		est.Days = est.XPRemaining / rate
		est.FinishDate = now.AddDate(0, 0, int(math.Ceil(est.Days))).Format("2006-01-02")
	}
	if opts.DeadlineDays > 0 {
		active := activeDays(opts.DeadlineDays, opts.RestDaysPerWeek)
		if active > 0 {
			est.RequiredPerDay = est.XPRemaining / float64(active)
		}
		cap := opts.FeasibleCap
		if cap <= 0 {
			cap = math.Inf(1)
		}
		est.Feasible = est.RequiredPerDay <= cap && est.RequiredPerDay > 0
	}
	return est
}

// learningCost prices an unmastered concept. Expected earned XP collapses to
// remaining-streak x base (accuracy cancels: more attempts at lower hit
// rate), while time scales with attempts. Don't-knows add reduced
// re-orientation time but never reduce the total (anti-gaming floor).
func learningCost(c *concepts.Concept, p *storage.ConceptProgress, s AttemptStats) ConceptCost {
	req := 3
	thresh := thresholdSec(c)
	if c != nil && c.MasteryThreshold.Streak > 0 {
		req = c.MasteryThreshold.Streak
	}
	streak := 0
	var elapsed float64
	if p != nil {
		streak = p.Streak
		if p.Attempts > 0 && p.AvgResponseTime > 0 {
			elapsed = p.AvgResponseTime
			if elapsed > 3*thresh {
				elapsed = 3 * thresh
			}
		} else {
			elapsed = thresh
		}
	} else {
		elapsed = thresh
	}
	remaining := req - streak
	if remaining <= 0 {
		remaining = 2 // streak met, time gate pending: a couple more tries
	}
	acc := accuracy(s)
	attempts := float64(remaining) / acc
	base := taskBaseFor(c.ID)
	cost := ConceptCost{
		ConceptID: c.ID,
		XP:        float64(remaining) * base,
		Attempts:  attempts,
	}
	cost.TimeMin = attempts*elapsed/60.0 + float64(s.DontKnow)*dontKnowTimeFrac*thresh/60.0
	// Anti-gaming floor: total time never drops below full-threshold pace,
	// so repeated don't-knows can only hold or grow the estimate.
	if floor := float64(remaining) * thresh / 60.0; cost.TimeMin < floor {
		cost.TimeMin = floor
	}
	return cost
}

// reviewCost prices one due-review pass: two attempts at review base.
func reviewCost(c *concepts.Concept, p *storage.ConceptProgress, s AttemptStats) ConceptCost {
	thresh := thresholdSec(c)
	elapsed := thresh
	if p != nil && p.Attempts > 0 && p.AvgResponseTime > 0 {
		elapsed = p.AvgResponseTime
		if elapsed > 3*thresh {
			elapsed = 3 * thresh
		}
	}
	acc := accuracy(s)
	base := float64(xp.BaseXP(xp.TaskReview))
	return ConceptCost{
		ConceptID: c.ID,
		XP:        2 * base * acc,
		TimeMin:   2 * elapsed / 60.0,
		Attempts:  2,
	}
}

// reviewDue reports whether a mastered concept needs a review pass.
func reviewDue(p *storage.ConceptProgress, now time.Time, decayDays float64) bool {
	if p == nil {
		return false
	}
	if p.NextReviewDue != nil && !p.NextReviewDue.After(now) {
		return true
	}
	if p.LastReviewed != nil && now.Sub(*p.LastReviewed).Hours()/24 >= decayDays {
		return true
	}
	return false
}

// effectiveRate applies rest days to a daily XP rate.
func effectiveRate(daily float64, restPerWeek int) float64 {
	if restPerWeek < 0 {
		restPerWeek = 0
	}
	if restPerWeek > 6 {
		restPerWeek = 6
	}
	return daily * float64(7-restPerWeek) / 7.0
}

// activeDays counts study days inside a deadline window given rest days.
// Remainder days are generous: rest is assumed to fall as late as possible.
func activeDays(deadlineDays, restPerWeek int) int {
	if restPerWeek < 0 {
		restPerWeek = 0
	}
	if restPerWeek > 6 {
		restPerWeek = 6
	}
	studyPerWeek := 7 - restPerWeek
	fullWeeks := deadlineDays / 7
	rem := deadlineDays % 7
	restInRem := rem - studyPerWeek
	if restInRem < 0 {
		restInRem = 0
	}
	active := fullWeeks*studyPerWeek + (rem - restInRem)
	if active < 1 && deadlineDays > 0 {
		active = 1
	}
	return active
}

// SampleClosure deterministically samples up to n representative concept IDs
// across the path's domains (round-robin over sorted domains and IDs).
// This is the stratified "initial knowledge estimate" probe list — never
// presented as mastery.
func SampleClosure(path *Path, n int) []string {
	if path == nil || n <= 0 {
		return nil
	}
	byDomain := map[string][]string{}
	for _, c := range path.Concepts {
		byDomain[c.Domain] = append(byDomain[c.Domain], c.ID)
	}
	var domains []string
	for d := range byDomain {
		domains = append(domains, d)
	}
	sort.Strings(domains)
	for _, ids := range byDomain {
		sort.Strings(ids)
	}
	var out []string
	for i := 0; len(out) < n; i++ {
		advanced := false
		for _, d := range domains {
			if i < len(byDomain[d]) {
				out = append(out, byDomain[d][i])
				advanced = true
				if len(out) >= n {
					break
				}
			}
		}
		if !advanced {
			break
		}
	}
	return out
}
