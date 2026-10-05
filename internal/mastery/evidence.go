package mastery

import (
	"math"
	"sort"
)

// Evidence is what the learner has actually shown on one concept, assembled from the
// attempts already recorded against it.
//
// This is the model `web/next-app/lib/progression.ts:masteryEstimate` already computes in
// the browser, moved to where the decision is actually made. That function's comment says
// "Streak is one input, not the decision", and this file is the version the server can
// act on. Keeping one model matters more than which model: the previous arrangement had
// two — a streak ratio gating progression and an evidence score rendering a label — and
// they could not both be right.
//
// Four signals, none of them new:
//
//   - Accuracy, recent-weighted, so a late miss still counts and an old one fades.
//   - Difficulty of the correct attempts, so getting it right when it was hard counts more.
//   - Variety: how many distinct problem instances the evidence spans.
//   - Time: median elapsed on correct answers against the concept's own threshold.
//
// Streak is deliberately *not* a term. Consecutive correctness measures consistency
// within a sitting; it is evidence, and it is available from `Streak` on the progress row
// for anyone who wants it, but a run of correct answers is not the same fact as competence.
type Attempt struct {
	Correct bool
	Elapsed float64
	// Difficulty the question was generated at, 0.3-1.0. Nil means unknown — see
	// attempts.difficulty, which is nullable precisely so this is expressible.
	Difficulty *float64
	// Instance distinguishes problem instances. Two attempts with the same text are the
	// same instance seen twice, which is not evidence of transfer.
	Instance string
}

// Evidence is the assembled picture. Zero fields mean "no evidence yet", which is a
// different claim from "evidence of nothing".
type Evidence struct {
	// Score is 0..1, the weighted composite.
	Score float64
	// Attempts is how many attempts the score is built from. Zero means unevidenced, and
	// callers must not read Score as meaningful when it is zero.
	Attempts int
	// Correct is the number of correct attempts among those considered.
	Correct int
	// Instances is how many distinct problem instances the correct attempts span.
	Instances int
	// AvgDifficulty is the mean difficulty of the correct attempts, 0 when unknown.
	AvgDifficulty float64
	// MedianSeconds is the median elapsed time over correct attempts.
	MedianSeconds float64
	// RecentWeightedAccuracy is the recency-weighted accuracy, before difficulty and
	// variety are folded in.
	RecentWeightedAccuracy float64
}

// EvidenceWindow is how many recent attempts the score considers. Six mirrors the browser
// model: long enough to contain a mistake and still be corrected, short enough that a
// learner is not carrying a month's of history into today's judgement.
const EvidenceWindow = 6

// EvidenceOptions carries the concept's own thresholds, so the time term means something
// relative to this concept rather than an absolute number of seconds.
type EvidenceOptions struct {
	// RequiredStreak is unused by the score itself and kept only so callers can pass the
	// concept row through unchanged. It exists in the signature because the authority
	// change needs both, and threading two arguments invites pairing them wrongly.
	_             struct{}
	TimeThreshold float64
}

// BuildEvidence assembles the score from the most recent attempts for one concept.
//
// Attempts are expected newest-last. Anything beyond the window is ignored rather than
// down-weighted, matching the browser model: a two-year-old attempt says nothing about
// whether the learner can do this now.
func BuildEvidence(attempts []Attempt, opts EvidenceOptions) Evidence {
	if len(attempts) > EvidenceWindow {
		attempts = attempts[len(attempts)-EvidenceWindow:]
	}
	ev := Evidence{Attempts: len(attempts)}
	if ev.Attempts == 0 {
		return ev
	}

	// Recency weighting: the last attempt counts twice the first.
	den := 0.0
	num := 0.0
	for i, a := range attempts {
		w := 0.5 + (float64(i)/float64(max(1, ev.Attempts-1)))*0.5
		den += w
		if a.Correct {
			num += w
		}
	}
	ev.RecentWeightedAccuracy = num / den

	// Difficulty, over the correct attempts only. Unknown difficulties are skipped rather
	// than counted as zero: a concept whose questions predate the column still has evidence,
	// just not this kind of it.
	var diffSum float64
	var diffN int
	instances := map[string]bool{}
	times := []float64{}
	for _, a := range attempts {
		if !a.Correct {
			continue
		}
		ev.Correct++
		if a.Difficulty != nil {
			diffSum += *a.Difficulty
			diffN++
		}
		if a.Instance != "" {
			instances[a.Instance] = true
		}
		times = append(times, a.Elapsed)
	}
	ev.Instances = len(instances)
	if diffN > 0 {
		ev.AvgDifficulty = diffSum / float64(diffN)
	}
	if len(times) > 0 {
		sort.Float64s(times)
		ev.MedianSeconds = times[len(times)/2]
	}

	// Difficulty weighting. `0.6 + 0.4 × difficulty` is the browser model's shape and is
	// kept so the two do not drift: at difficulty 0.3 a correct answer is worth 0.72, at
	// 1.0 it is worth 1.0.
	//
	// Unknown difficulty is credited at NeutralDifficulty, *not* dropped. Dropping the
	// term means multiplying by 1.0, which is the top of the range — so a missing value was
	// being read as "the learner got this right on the hardest question we could have
	// asked", and the comment claiming unknown attempts were "neither rewarded nor punished"
	// was describing the opposite of what the code did. It was not cosmetic: at 50%
	// accuracy the unknown-difficulty score was 0.65 and cleared the threshold, while the
	// same attempts with difficulty 0.6 scored 0.57 and did not. A learner could therefore
	// pass the mastery bar on the quiz surface and fail it on the study surface for the same
	// answers.
	difficultyTerm := difficultyWeight(NeutralDifficulty)
	if diffN > 0 {
		difficultyTerm = difficultyWeight(ev.AvgDifficulty)
	}

	// Variety: each instance beyond the first is worth up to 0.15, capped. Two instances is
	// the minimum for the bonus, so one lucky question cannot look like transfer.
	varietyBonus := math.Min(0.15, 0.05*float64(max(0, ev.Instances-1)))

	// Time, against the concept's own threshold. At or under it is full credit; beyond it
	// decays gently rather than falling off a cliff, because a slow correct answer is still
	// a correct answer.
	timeFactor := 1.0
	if opts.TimeThreshold > 0 && ev.MedianSeconds > opts.TimeThreshold {
		timeFactor = 0.8
	}

	ev.Score = clamp01((ev.RecentWeightedAccuracy*difficultyTerm + varietyBonus) * timeFactor)
	return ev
}

// MasteryThreshold is the score at which a concept counts as demonstrated.
//
// It is not a per-concept value, deliberately. `mastery_threshold.streak` was per-concept
// and is what produced a 30-answer requirement for 70% of the corpus; the score is already
// normalised against the concept's own difficulty and time threshold, so a single
// constant is the right shape. Concepts differ in what counts as a hard question, not in
// how much evidence proves competence.
//
// Measured rather than assumed, and it does bind. `TestMeasureMasteryCost` walks realistic
// interleaved accuracy profiles: 100%, 80% and 60% reach MASTERED in exactly
// MasteryEvidenceFloor answers, while 50%, 40%, 20% and 0% never reach it at any length.
// The threshold separates those two groups cleanly, so it is not decorative.
//
// An earlier version of that measurement reported the opposite — that every accuracy level
// from 40% up mastered — because the test profile was periodic in blocks of 20 (twelve
// wrong, then eight right), so a six-attempt window eventually landed entirely inside the
// clean run. A learner does not answer in blocks. The lesson is recorded because the
// failure mode recurs: a periodic probe measures the probe's period, not the model.
const MasteryThreshold = 0.6

// NeutralDifficulty is the difficulty credited to a correct answer whose attempt carries
// none — an attempt predating `attempts.difficulty`, or a curated question with no
// recorded value.
//
// It is not chosen here. It is the difficulty the engine serves when it has no information
// about the learner: `engine.difficultyFromWeakness` defaults weakness to 0.5 and computes
// `0.3 + (1 - 0.5) × 0.5 = 0.55`. `TestNeutralDifficultyMatchesEngineDefault` in the engine
// package asserts the two agree, so the scoring side cannot drift away from the serving
// side. Crediting an average question is the honest reading of "we do not know how hard
// this was".
const NeutralDifficulty = 0.55

// difficultyWeight maps a difficulty onto the score multiplier for a correct answer, in
// [0.72, 1.0]. One place, so the known and unknown paths cannot diverge.
func difficultyWeight(difficulty float64) float64 {
	return 0.6 + 0.4*clampDifficulty(difficulty)
}

func clampDifficulty(d float64) float64 {
	if d < 0.3 {
		return 0.3
	}
	if d > 1.0 {
		return 1.0
	}
	return d
}

// FirstEvidence is the bar for LEARNING: one correct answer.
//
// One answer is enough to say "this learner has begun", and it is deliberately the cheapest
// rung. The old ladder demanded a full `required_streak` — 10 consecutive correct answers for
// 459 of the 657 concepts — before it would say a learner was learning, so a first session
// was blind repetition with no signal that any of it was landing.
func FirstEvidence(ev Evidence) bool {
	return ev.Correct >= 1
}

// EnoughForPracticing is the bar for PRACTICING: consistency across a couple of attempts.
// Weaker than mastery on purpose — this rung asks "is this dependable", which does not need
// the breadth that the top rung asks for.
func EnoughForPracticing(ev Evidence) bool {
	return ev.Score >= MasteryThreshold && ev.Attempts >= 2
}

// MasteryEvidenceFloor is the minimum number of attempts the evidence must span before a
// concept can count as mastered.
//
// It is a *floor on evidence*, not a rule about what those attempts are. Six identical
// correct answers are not the requirement; six attempts whose combined score clears the
// threshold are, and the score already weights difficulty, variety of problem instance and
// time. A lone correct answer is weighted 1.0 and would clear the threshold on its own, so
// something has to stop one lucky question from being a state change — and the number that
// stops it is this, not the shape of the attempts.
//
// Six rather than three because ADR-038 measured the ladder reaching mastery in three correct
// answers on a clean run, which is below the ~9 evidence-bearing answers the design intended.
// The cost was not the goal; a reasonable amount of evidence is. Since the window is six
// attempts wide, a floor of six means the score must be clear over a full window rather than
// over the first three answers, so a learner cannot master a concept before the evidence has
// had room to include a mistake and a recovery.
const MasteryEvidenceFloor = 6

// EnoughForMastery is the bar for MASTERED: demonstrated, over a full window of evidence.
//
// **Variety deliberately scores but does not gate.** An earlier version also required
// `Instances >= 2`, and wiring this into the engine immediately produced a state that could
// never be reached: a generator whose question text does not vary between attempts holds
// Instances at 1 forever, so the concept cannot be mastered however much the learner
// demonstrates it. A gate on a derived signal can make a state unreachable, and an
// unreachable mastery is a far worse failure than a slightly generous one. Variety still
// contributes up to 0.15 through varietyBonus — worth real evidence, but never a veto.
func EnoughForMastery(ev Evidence) bool {
	return ev.Score >= MasteryThreshold && ev.Attempts >= MasteryEvidenceFloor
}

// Rung names which transition a given piece of evidence authorises. It exists so the
// decision reads as a ladder in one place and the three predicates cannot be applied out
// of order by accident.
type Rung int

const (
	// RungNone means the evidence does not move the concept.
	RungNone Rung = iota
	// RungStarted authorises UNSEEN → LEARNING.
	RungStarted
	// RungConsistent authorises LEARNING → PRACTICING.
	RungConsistent
	// RungDemonstrated authorises PRACTICING → MASTERED.
	RungDemonstrated
)

// HighestRung is the furthest the evidence reaches, which is what `Machine.Next` needs:
// the ladder advances one rung per decision, never several.
func HighestRung(ev Evidence) Rung {
	switch {
	case EnoughForMastery(ev):
		return RungDemonstrated
	case EnoughForPracticing(ev):
		return RungConsistent
	case FirstEvidence(ev):
		return RungStarted
	default:
		return RungNone
	}
}

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
