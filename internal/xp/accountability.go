package xp

import (
	"fmt"
	"math"
	"time"
)

// Effort and performance accountability separate productive work from
// exploitation of the XP economy. They are deliberately separate from Award: an
// award prices the work a question represents, while accountability asks
// whether the work was done and whether the result was good enough to reward.
//
// Three properties are non-negotiable:
//
//   - An incorrect answer never loses XP on its own. The effort rules below need
//     a repeated pattern, never one miss, one low score, one fast submission or
//     one abandoned exercise. Performance rules need a task-sized sample and a
//     result substantially below what the task's difficulty predicts.
//   - Genuine difficulty is not exploitation. Slow wrong answers, low quiz
//     scores and failed mastery checks return OutcomeFull.
//   - Mastery, FSRS and evidence never read this policy. It only writes signed
//     rows to the XP ledger, so a deduction cannot move a review date, change a
//     mastery verdict or alter the prerequisite graph.
//
// The ledger is the system of record. Every adjustment carries a stable reason
// code, a source dimension and one idempotency key per task, so a retried or
// concurrent submission cannot be charged twice. The legacy two-fast-miss rule
// in the engine is left alone; when it has already fired in a task, this layer
// stays out of the same task and lets that existing penalty stand alone.

// Outcome classifies what accountability decided for one task.
type Outcome string

const (
	// OutcomeBonus pays extra XP for sustained effort or exceptional results.
	OutcomeBonus Outcome = "bonus"
	// OutcomeFull leaves the existing award exactly as it was.
	OutcomeFull Outcome = "full"
	// OutcomePartial reduces an award the task had already earned.
	OutcomePartial Outcome = "partial"
	// OutcomeZero records that no adjustment is possible or appropriate. A
	// performance decision can carry a negative amount here when the intended
	// result is that the task's earned XP becomes zero.
	OutcomeZero Outcome = "zero"
	// OutcomeNegative deducts XP from the ledger, never below zero balance.
	OutcomeNegative Outcome = "negative"
)

// ReasonCode is the stable, machine-readable cause of an adjustment. It is
// written to the ledger as "accountability:<code>".
type ReasonCode string

const (
	ReasonNone                     ReasonCode = ""
	ReasonInsufficientEffort       ReasonCode = "insufficient_effort"
	ReasonRepeatedGuessing         ReasonCode = "repeated_guessing"
	ReasonFarming                  ReasonCode = "farming"
	ReasonGenuineDifficulty        ReasonCode = "genuine_difficulty"
	ReasonLegacyPenalty            ReasonCode = "legacy_penalty"
	ReasonDiagnosticExempt         ReasonCode = "diagnostic_exempt"
	ReasonNoBalance                ReasonCode = "no_balance"
	ReasonSustainedEffort          ReasonCode = "sustained_effort"
	ReasonAcceptablePerformance    ReasonCode = "acceptable_performance"
	ReasonPoorPerformance          ReasonCode = "poor_performance"
	ReasonExtremelyPoorPerformance ReasonCode = "extremely_poor_performance"
	ReasonExceptionalPerformance   ReasonCode = "exceptional_performance"
)

// Source identifies which accountability dimension produced a decision. It is
// returned to the learner so an explanation can say whether the result came
// from performance or from the answer pattern.
const (
	SourceNone        = ""
	SourcePerformance = "performance"
	SourceEffort      = "effort"
	SourceLegacy      = "legacy"
)

// Task kinds understood by this policy. They mirror the engine's task names and
// the two assessment kinds that are not exercised through the practice loop.
const (
	TaskPractice = "practice"
	TaskStudy    = "study"
	// TaskReview and TaskQuiz live in policy.go.
	TaskTimedQuiz    = "timed_quiz"
	TaskMasteryCheck = "mastery_check"
	TaskDiagnostic   = "diagnostic"
)

// Config holds every threshold in one place. The defaults are intentionally
// conservative: the cost of missing an exploiter is one task's XP, while the
// cost of accusing a struggling learner is losing them.
type Config struct {
	// Disabled turns both accountability dimensions off without changing any
	// other behaviour. It is the zero-value-safe spelling: a Config nobody
	// filled in is enabled with the defaults below.
	Disabled bool
	// FastAnswerSeconds is the upper bound for an answer too quick to show
	// work. It is shorter than RushSeconds because the legacy rule already
	// owns sub-RushSeconds guesses.
	FastAnswerSeconds float64
	// MinTaskAttempts is the smallest task worth judging. A short task cannot
	// establish a pattern.
	MinTaskAttempts int
	// MinFastWrong is how many too-fast wrong answers establish guessing. It is
	// at least three so one or two taps never qualify.
	MinFastWrong int
	// MinSameWrong is how many identical wrong answers establish answer churn.
	MinSameWrong int
	// SameWrongSeconds is the upper bound for a repeated wrong answer to count
	// as churn. A slow repeated misconception is evidence of difficulty, not
	// guessing, so it is left alone.
	SameWrongSeconds float64
	// FarmingWindow bounds the recent-answer window used for farming.
	FarmingWindow time.Duration
	// FarmingMinCorrect is how many fast correct answers on already mastered
	// material count as farming rather than fluency.
	FarmingMinCorrect int
	// FarmingFastSeconds is the upper bound for a correct answer that can count
	// toward farming. It is very short so ordinary fluent practice does not.
	FarmingFastSeconds float64
	// EscalationWindow is how far back prior deductions are counted for the
	// escalating ladder.
	EscalationWindow time.Duration
	// PenaltyFirst, PenaltyRepeated and PenaltyPersistent are positive XP
	// amounts deducted at escalation levels zero, one and two or more. A zero
	// value is intentional and disables that level.
	PenaltyFirst      int
	PenaltyRepeated   int
	PenaltyPersistent int
	// EnableBonus and BonusAmount configure the optional positive outcome. The
	// bonus is off by default because it changes the published XP economy.
	EnableBonus bool
	BonusAmount int
	// Performance configures task-result accountability independently from
	// effort detection.
	Performance PerformanceConfig
}

// PerformanceConfig describes when a task's result is good enough for the full
// award, only part of it, none of it or a negative correction.
//
// The ratio is observed correct answers divided by expected correct answers.
// Expected correctness comes from each item's difficulty, so a hard task can
// have a low raw score without being a poor performance. Thresholds are
// deliberately data-shaped rather than fixed percentages of raw accuracy.
type PerformanceConfig struct {
	// Disabled turns only the performance dimension off. The effort dimension
	// can remain active.
	Disabled bool
	// MinAttempts is the smallest result sample that can reduce XP. It is
	// larger than one because a single miss is never a performance verdict.
	MinAttempts int
	// NegativeMinAttempts is the smallest sample that can produce a negative
	// result, including when the task earned no XP.
	NegativeMinAttempts int
	// MinExpectedCorrect prevents a handful of very hard items from being
	// judged against a near-zero expected score.
	MinExpectedCorrect float64
	// NegativeMinExpectedCorrect is the stronger floor for a negative result.
	NegativeMinExpectedCorrect float64
	// GoodRatio is the observed/expected boundary for full XP.
	GoodRatio float64
	// AcceptableRatio is the boundary between acceptable and poor. Below it the
	// task keeps none of its earned XP.
	AcceptableRatio float64
	// NegativeRatio is the boundary for a negative result.
	NegativeRatio float64
	// AcceptableDeduction is the fraction of earned XP removed from an
	// acceptable task.
	AcceptableDeduction float64
	// NegativePenalty is the extra XP deducted beyond the task's earned XP for
	// an extremely poor result. It is what makes negative XP possible when the
	// task earned nothing.
	NegativePenalty int
	// BonusRatio and BonusMinAttempts gate the optional exceptional-result
	// bonus. The existing EnableBonus and BonusAmount switches still apply.
	BonusRatio       float64
	BonusMinAttempts int
}

// DefaultConfig returns the shipped thresholds.
func DefaultConfig() Config {
	return Config{
		Disabled:           false,
		FastAnswerSeconds:  1.5,
		MinTaskAttempts:    4,
		MinFastWrong:       3,
		MinSameWrong:       3,
		SameWrongSeconds:   5,
		FarmingWindow:      20 * time.Minute,
		FarmingMinCorrect:  40,
		FarmingFastSeconds: 1.0,
		EscalationWindow:   14 * 24 * time.Hour,
		PenaltyFirst:       1,
		PenaltyRepeated:    2,
		PenaltyPersistent:  3,
		EnableBonus:        false,
		BonusAmount:        1,
		Performance: PerformanceConfig{
			MinAttempts:                4,
			NegativeMinAttempts:        6,
			MinExpectedCorrect:         2.0,
			NegativeMinExpectedCorrect: 2.5,
			GoodRatio:                  0.8,
			AcceptableRatio:            0.5,
			NegativeRatio:              0.2,
			AcceptableDeduction:        0.25,
			NegativePenalty:            2,
			BonusRatio:                 1.25,
			BonusMinAttempts:           5,
		},
	}
}

// withDefaults fills fields whose zero value cannot mean anything useful. It
// preserves intentional zeroes for penalty and bonus amounts and for
// performance ratios or deductions, because those zeroes disable a rule.
func (c Config) withDefaults() Config {
	d := DefaultConfig()
	if c.FastAnswerSeconds <= 0 {
		c.FastAnswerSeconds = d.FastAnswerSeconds
	}
	if c.MinTaskAttempts <= 0 {
		c.MinTaskAttempts = d.MinTaskAttempts
	}
	if c.MinFastWrong <= 0 {
		c.MinFastWrong = d.MinFastWrong
	}
	if c.MinSameWrong <= 0 {
		c.MinSameWrong = d.MinSameWrong
	}
	if c.SameWrongSeconds <= 0 {
		c.SameWrongSeconds = d.SameWrongSeconds
	}
	if c.FarmingWindow <= 0 {
		c.FarmingWindow = d.FarmingWindow
	}
	if c.FarmingMinCorrect <= 0 {
		c.FarmingMinCorrect = d.FarmingMinCorrect
	}
	if c.FarmingFastSeconds <= 0 {
		c.FarmingFastSeconds = d.FarmingFastSeconds
	}
	if c.EscalationWindow <= 0 {
		c.EscalationWindow = d.EscalationWindow
	}
	if c.PenaltyFirst < 0 {
		c.PenaltyFirst = d.PenaltyFirst
	}
	if c.PenaltyRepeated < 0 {
		c.PenaltyRepeated = d.PenaltyRepeated
	}
	if c.PenaltyPersistent < 0 {
		c.PenaltyPersistent = d.PenaltyPersistent
	}
	if c.BonusAmount < 0 {
		c.BonusAmount = d.BonusAmount
	}
	c.Performance = c.Performance.withDefaults(d.Performance)
	return c
}

func (p PerformanceConfig) withDefaults(d PerformanceConfig) PerformanceConfig {
	if p.MinAttempts <= 0 {
		p.MinAttempts = d.MinAttempts
	}
	if p.NegativeMinAttempts <= 0 {
		p.NegativeMinAttempts = d.NegativeMinAttempts
	}
	if p.MinExpectedCorrect < 0 {
		p.MinExpectedCorrect = d.MinExpectedCorrect
	}
	if p.NegativeMinExpectedCorrect < 0 {
		p.NegativeMinExpectedCorrect = d.NegativeMinExpectedCorrect
	}
	if p.GoodRatio < 0 {
		p.GoodRatio = d.GoodRatio
	}
	if p.AcceptableRatio < 0 {
		p.AcceptableRatio = d.AcceptableRatio
	}
	if p.NegativeRatio < 0 {
		p.NegativeRatio = d.NegativeRatio
	}
	if p.AcceptableDeduction < 0 {
		p.AcceptableDeduction = d.AcceptableDeduction
	}
	if p.NegativePenalty < 0 {
		p.NegativePenalty = d.NegativePenalty
	}
	if p.BonusRatio < 0 {
		p.BonusRatio = d.BonusRatio
	}
	if p.BonusMinAttempts <= 0 {
		p.BonusMinAttempts = d.BonusMinAttempts
	}
	return p
}

// DisabledConfig returns a Config that leaves every task at Full.
func DisabledConfig() Config {
	c := DefaultConfig()
	c.Disabled = true
	return c
}

// NormalizeConfig fills unset thresholds with the shipped defaults. The engine
// calls it once at construction so call sites read the same numbers the
// evaluator will use.
func NormalizeConfig(c Config) Config {
	return c.withDefaults()
}

// AttemptSignal is the accountability view of one graded answer. It carries no
// concept label, no answer text and no mastery verdict, because the layer must
// not be able to influence learning evidence.
type AttemptSignal struct {
	TaskID          string
	Correct         bool
	ElapsedSeconds  float64
	DontKnow        bool
	Mastered        bool
	Difficulty      float64
	AnswerSignature string
	LegacyPenalty   bool
	// XPEarned is the gross positive award for this answer. It is zero for a
	// wrong answer and for an answer whose award was already zero.
	XPEarned int
	At       time.Time
}

// Input is one accountability decision. Attempts is the current task; Recent is
// a bounded window across the learner's tasks, used only for farming.
type Input struct {
	TaskKind string
	TaskID   string
	Attempts []AttemptSignal
	Recent   []AttemptSignal
	// PriorAdjustments is the number of confirmed deductions in the escalation
	// window before this task.
	PriorAdjustments int
	// AvailableXP is the learner's current balance. A deduction is capped by it
	// so the balance can never go below zero.
	AvailableXP int
	// TaskComplete is true only when the engine has a reliable signal that the
	// task is finished. Performance never judges an unfinished task; effort
	// detection remains immediate.
	TaskComplete bool
	Config       Config
}

// Decision is the result for one task. Amount is signed: negative deducts,
// positive bonuses, zero changes nothing. Note is student-facing and never
// accusatory.
type Decision struct {
	Outcome Outcome
	// Amount is the adjustment after the policy's available-balance cap.
	Amount int
	// CalculatedAmount is the uncapped amount before that cap. It exists so
	// observability can distinguish "calculated -2" from "applied -1".
	CalculatedAmount int
	Reason           ReasonCode
	Explanation      string
	Note             string
	Level            int
	Source           string
}

// LedgerReason is the reason written to xp_events. It is stable and safe to
// query, so history can report adjustments without parsing prose.
func (d Decision) LedgerReason() string {
	if d.Reason == ReasonNone {
		return ""
	}
	return "accountability:" + string(d.Reason)
}

// taskStats is the summary the effort rules read. None of these numbers is an
// accuracy threshold used on its own; each rule also requires a pattern of
// repeated or implausibly fast answers.
type taskStats struct {
	attempts     int
	correct      int
	fastWrong    int
	sameWrong    int
	fastCorrect  int
	farmingFast  int
	legacyFired  bool
	maxDiff      float64
	totalElapsed float64
	xpEarned     int
}

func summarize(attempts []AttemptSignal, cfg Config) taskStats {
	var s taskStats
	counts := map[string]int{}
	for _, a := range attempts {
		s.attempts++
		s.totalElapsed += a.ElapsedSeconds
		if a.Difficulty > s.maxDiff {
			s.maxDiff = a.Difficulty
		}
		if a.XPEarned > 0 {
			s.xpEarned += a.XPEarned
		}
		if a.LegacyPenalty {
			s.legacyFired = true
		}
		if a.Correct {
			s.correct++
			if a.ElapsedSeconds >= 0 && a.ElapsedSeconds < cfg.FarmingFastSeconds {
				s.fastCorrect++
				if a.Mastered {
					s.farmingFast++
				}
			}
			continue
		}
		if a.DontKnow {
			// An admitted unknown is honesty, not rushing. It never counts
			// toward a fast-wrong or churn rule.
			continue
		}
		if a.ElapsedSeconds >= 0 && a.ElapsedSeconds < cfg.FastAnswerSeconds {
			s.fastWrong++
		}
		if a.AnswerSignature != "" && a.ElapsedSeconds >= 0 && a.ElapsedSeconds < cfg.SameWrongSeconds {
			counts[a.AnswerSignature]++
			if counts[a.AnswerSignature] > s.sameWrong {
				s.sameWrong = counts[a.AnswerSignature]
			}
		}
	}
	return s
}

// Evaluate decides the effort dimension of one task. The order matters: an
// exempt context returns immediately, then the strongest exploitation signals,
// then the guards for genuine difficulty, and only at the end the optional
// bonus.
func Evaluate(in Input) Decision {
	cfg := in.Config.withDefaults()
	if in.TaskKind == TaskDiagnostic {
		return Decision{
			Outcome:     OutcomeFull,
			Reason:      ReasonDiagnosticExempt,
			Explanation: "Diagnostic answers choose a starting point and never change XP.",
			Source:      SourceEffort,
		}
	}
	if cfg.Disabled || len(in.Attempts) == 0 {
		return Decision{Outcome: OutcomeFull}
	}

	task := summarize(in.Attempts, cfg)
	recent := summarize(in.Recent, cfg)

	// The legacy rule already charged this task for rapid wrong answers. Do not
	// stack a second deduction on the same task; the existing penalty stands.
	if task.legacyFired {
		return Decision{
			Outcome:     OutcomeFull,
			Reason:      ReasonLegacyPenalty,
			Explanation: "Two very fast wrong answers triggered the existing rush rule.",
			Note:        fmt.Sprintf("%d XP: Two very fast wrong answers triggered the existing rush rule.", RushPenalty),
			Source:      SourceLegacy,
		}
	}

	// Farming is independent of accuracy: a long run of very fast correct
	// answers on material already mastered is not practice, whether or not any
	// answer was wrong.
	if recent.farmingFast >= cfg.FarmingMinCorrect {
		return penalize(in, cfg, ReasonFarming, task)
	}

	// A pattern of implausibly fast wrong answers. The legacy rule owns the
	// sub-RushSeconds case when it fires; this catches the band between
	// FastAnswerSeconds and RushSeconds, and hosts whose legacy counter is cold.
	if task.attempts >= cfg.MinTaskAttempts && task.fastWrong >= cfg.MinFastWrong {
		return penalize(in, cfg, ReasonInsufficientEffort, task)
	}

	// The same wrong answer repeated unanswered is answer churn, not thinking.
	// It is not a single wrong answer and not a low score.
	if task.sameWrong >= cfg.MinSameWrong && task.correct == 0 {
		return penalize(in, cfg, ReasonRepeatedGuessing, task)
	}

	// Genuine difficulty: a full task of considered wrong answers, a failed
	// mastery check or a low quiz score lands here with no note and no change.
	if task.correct == 0 && task.maxDiff >= 0.7 {
		return Decision{Outcome: OutcomeFull, Reason: ReasonGenuineDifficulty, Source: SourceEffort}
	}

	if cfg.EnableBonus && task.attempts >= 5 && task.correct == task.attempts && task.fastWrong == 0 && task.maxDiff >= 0.7 {
		return Decision{
			Outcome:     OutcomeBonus,
			Amount:      cfg.BonusAmount,
			Reason:      ReasonSustainedEffort,
			Explanation: "A full task was completed correctly at a difficulty that takes sustained effort.",
			Note:        fmt.Sprintf("+%d XP: Sustained effort on a difficult task.", cfg.BonusAmount),
			Source:      SourceEffort,
		}
	}

	return Decision{Outcome: OutcomeFull}
}

// EvaluatePerformance decides the result dimension of one task. It compares
// observed correctness with the correctness the item difficulties predict, so
// hard material is judged against a lower expectation rather than against the
// same raw score as easy material.
func EvaluatePerformance(in Input) Decision {
	cfg := in.Config.withDefaults()
	if in.TaskKind == TaskDiagnostic {
		return Decision{
			Outcome:     OutcomeFull,
			Reason:      ReasonDiagnosticExempt,
			Explanation: "Diagnostic answers choose a starting point and never change XP.",
			Source:      SourceEffort,
		}
	}
	if cfg.Disabled || cfg.Performance.Disabled || len(in.Attempts) == 0 {
		return Decision{Outcome: OutcomeFull}
	}

	attempts := make([]AttemptSignal, 0, len(in.Attempts))
	for _, a := range in.Attempts {
		if !a.DontKnow {
			attempts = append(attempts, a)
		}
	}
	if len(attempts) == 0 {
		return Decision{Outcome: OutcomeFull}
	}
	// A task key is idempotent, so a performance verdict written before the
	// learner finishes cannot be revised. Wait for a completion signal.
	if !in.TaskComplete {
		return Decision{Outcome: OutcomeFull}
	}

	task := summarize(attempts, cfg)
	p := cfg.Performance
	if task.attempts < p.MinAttempts {
		return Decision{Outcome: OutcomeFull}
	}
	if task.legacyFired {
		return Decision{Outcome: OutcomeFull, Reason: ReasonLegacyPenalty, Source: SourceLegacy}
	}

	if task.correct == 0 && task.maxDiff >= 0.7 {
		return Decision{Outcome: OutcomeFull, Reason: ReasonGenuineDifficulty, Source: SourcePerformance}
	}

	expected := expectedCorrect(attempts)
	if expected <= 0 {
		return Decision{Outcome: OutcomeFull}
	}
	ratio := float64(task.correct) / expected

	// A high-difficulty task is not made more culpable by a single correct
	// answer. Without this branch, one correct answer could fall into the
	// negative band while zero correct answers were fully exempt, so a learner
	// could be punished for improving. Remove earned XP at most, never add the
	// negative penalty on a task the difficulty itself explains.
	if task.maxDiff >= 0.7 && ratio < p.NegativeRatio {
		if task.xpEarned <= 0 {
			return Decision{Outcome: OutcomeFull, Reason: ReasonGenuineDifficulty, Source: SourcePerformance}
		}
		return performanceAdjustment(OutcomeZero, -task.xpEarned, ReasonGenuineDifficulty,
			"The task was difficult. XP earned on it was removed, but no negative penalty was applied.", in.AvailableXP)
	}

	if task.attempts >= p.NegativeMinAttempts && expected >= p.NegativeMinExpectedCorrect && ratio < p.NegativeRatio {
		amount := -(task.xpEarned + p.NegativePenalty)
		if amount >= 0 {
			return Decision{
				Outcome:     OutcomeZero,
				Reason:      ReasonExtremelyPoorPerformance,
				Explanation: "Task performance was extremely poor, but there was no XP to remove.",
				Note:        "0 XP: Task performance was extremely poor.",
				Source:      SourcePerformance,
			}
		}
		return performanceAdjustment(OutcomeNegative, amount, ReasonExtremelyPoorPerformance,
			"Task performance was extremely poor across this task.", in.AvailableXP)
	}

	if expected >= p.MinExpectedCorrect && ratio < p.AcceptableRatio {
		if task.xpEarned <= 0 {
			return Decision{
				Outcome:     OutcomeZero,
				Reason:      ReasonPoorPerformance,
				Explanation: "Task performance was poor, but the task earned no XP to remove.",
				Note:        "0 XP: Task performance was poor.",
				Source:      SourcePerformance,
			}
		}
		return performanceAdjustment(OutcomeZero, -task.xpEarned, ReasonPoorPerformance,
			"Task performance was poor, so the XP earned on this task was removed.", in.AvailableXP)
	}

	if expected >= p.MinExpectedCorrect && ratio < p.GoodRatio {
		if task.xpEarned <= 0 {
			return Decision{Outcome: OutcomeFull}
		}
		deduction := int(math.Round(float64(task.xpEarned) * p.AcceptableDeduction))
		if deduction < 1 && p.AcceptableDeduction > 0 {
			deduction = 1
		}
		if deduction > task.xpEarned {
			deduction = task.xpEarned
		}
		if deduction <= 0 {
			return Decision{Outcome: OutcomeFull}
		}
		return performanceAdjustment(OutcomePartial, -deduction, ReasonAcceptablePerformance,
			"Task performance was acceptable, so part of the XP earned was removed.", in.AvailableXP)
	}

	if cfg.EnableBonus && cfg.BonusAmount > 0 && task.attempts >= p.BonusMinAttempts &&
		task.correct == task.attempts && (task.maxDiff >= 0.7 || ratio >= p.BonusRatio) {
		return Decision{
			Outcome:     OutcomeBonus,
			Amount:      cfg.BonusAmount,
			Reason:      ReasonExceptionalPerformance,
			Explanation: "The whole task was correct at a difficulty that predicts a much lower score.",
			Note:        fmt.Sprintf("+%d XP: Exceptional task performance.", cfg.BonusAmount),
			Source:      SourcePerformance,
		}
	}

	return Decision{Outcome: OutcomeFull}
}

func expectedCorrect(attempts []AttemptSignal) float64 {
	var expected float64
	for _, a := range attempts {
		if a.DontKnow {
			continue
		}
		expected += difficultyCorrectPrior(a.Difficulty)
	}
	return expected
}

// difficultyCorrectPrior maps a served difficulty to an expected correct rate.
//
// Difficulty in Mathua is a generator targeting parameter in [0.3, 1.0], not a
// calibrated probability. It is derived from learner weakness and used by the
// mastery scorer as an ordinal weight (`0.6 + 0.4×difficulty`), so `1-difficulty`
// was not a defensible probability. The bands below are observed correct rates
// from production attempts measured on 2026-11-06, excluding diagnostics:
//
//	0.30       6,630 attempts   83.0% correct
//	0.40–0.45  58,384 attempts  89.9% / 84.7% correct
//	0.55       43,612 attempts  94.4% correct
//	0.80       7,622 attempts   87.2% correct
//
// Diagnostics are excluded because they are exempt from this policy and their
// 0.3/0.6/0.8 bands have very different accuracy. No non-diagnostic 0.6 band
// was observed, so it uses the conservative lower edge of the surrounding
// non-diagnostic rates. The scale is not monotonic across bands, which is
// itself evidence that one number cannot be read as a probability across
// question families. Unknown or zero difficulty is treated as a neutral 50%,
// not as a maximally easy question.
func difficultyCorrectPrior(difficulty float64) float64 {
	if difficulty <= 0 {
		return 0.50
	}
	if difficulty > 1 {
		difficulty = 1
	}
	switch {
	case difficulty < 0.35:
		return 0.83
	case difficulty < 0.50:
		return 0.89
	case difficulty < 0.575:
		return 0.944
	case difficulty < 0.70:
		return 0.85
	default:
		return 0.872
	}
}

func performanceAdjustment(outcome Outcome, amount int, reason ReasonCode, explanation string, available int) Decision {
	calculated := amount
	if available < 0 {
		available = 0
	}
	if amount < -available {
		amount = -available
	}
	if amount == 0 {
		return Decision{
			Outcome:          OutcomeZero,
			Amount:           0,
			CalculatedAmount: calculated,
			Reason:           ReasonNoBalance,
			Explanation:      explanation,
			Note:             "No XP deducted. The XP balance is already zero.",
			Source:           SourcePerformance,
		}
	}
	return Decision{
		Outcome:          outcome,
		Amount:           amount,
		CalculatedAmount: calculated,
		Reason:           reason,
		Explanation:      explanation,
		Note:             fmt.Sprintf("%d XP: %s", amount, explanation),
		Source:           SourcePerformance,
	}
}

// Combine folds the two dimensions into one ledger decision. A task never
// receives two negative adjustments for the same work: effort-specific
// findings keep their existing reason and amount, while a performance verdict
// applies only when the effort dimension is neutral. A legacy rush penalty and
// a genuine-difficulty exemption suppress the performance verdict entirely.
func Combine(effort, performance Decision) Decision {
	if effort.Reason == ReasonLegacyPenalty || effort.Reason == ReasonGenuineDifficulty {
		return effort
	}
	if effort.Amount < 0 {
		return effort
	}
	if performance.Amount < 0 {
		return performance
	}
	if effort.Amount > 0 && performance.Amount > 0 {
		if performance.Amount >= effort.Amount {
			return performance
		}
		return effort
	}
	if effort.Amount > 0 {
		return effort
	}
	if performance.Amount > 0 {
		return performance
	}
	if effort.Outcome == OutcomeZero {
		return effort
	}
	if performance.Outcome == OutcomeZero {
		return performance
	}
	return effort
}

func penalize(in Input, cfg Config, reason ReasonCode, task taskStats) Decision {
	level := in.PriorAdjustments
	if level > 2 {
		level = 2
	}
	penalty := cfg.PenaltyFirst
	switch level {
	case 1:
		penalty = cfg.PenaltyRepeated
	case 2:
		penalty = cfg.PenaltyPersistent
	}
	message := "Insufficient effort detected across this task."
	if level == 1 {
		message = "Repeated low-effort attempts."
	} else if level >= 2 {
		message = "Low effort continued after earlier feedback."
	}

	outcome := OutcomeNegative
	if task.correct > 0 {
		// Some work in the task was productive, so this is a reduction of what
		// the task earned rather than an extra charge on top of nothing.
		outcome = OutcomePartial
	}
	available := in.AvailableXP
	if available < 0 {
		available = 0
	}
	calculated := -penalty
	amount := calculated
	if amount < -available {
		amount = -available
	}
	if amount == 0 {
		return Decision{
			Outcome:          OutcomeZero,
			Amount:           0,
			CalculatedAmount: calculated,
			Reason:           ReasonNoBalance,
			Explanation:      "A low-effort pattern was detected, but the XP balance is already zero, so no deduction was made.",
			Note:             "No XP deducted. Review the explanation and try again.",
			Level:            level,
			Source:           SourceEffort,
		}
	}
	return Decision{
		Outcome:          outcome,
		Amount:           amount,
		CalculatedAmount: calculated,
		Reason:           reason,
		Explanation:      message,
		Note:             fmt.Sprintf("%d XP: %s", amount, message),
		Level:            level,
		Source:           SourceEffort,
	}
}
