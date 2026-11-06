package engine

import (
	"sync"

	xppolicy "github.com/chuma-beep/mathua/internal/xp"
)

// AccountabilityMetricsSnapshot is the process-local audit view. It is not a
// durable ledger: xp_events remains the record of truth. These counters answer
// "what did the policy calculate, what did storage apply, and where did the
// cap or a replay change the outcome" during the life of this process.
type AccountabilityMetricsSnapshot struct {
	Decisions                   map[string]int64 `json:"decisions"`
	CalculatedDeductions        int64            `json:"calculated_deductions"`
	CalculatedDeductionXP       int64            `json:"calculated_deduction_xp"`
	AppliedDeductions           int64            `json:"applied_deductions"`
	AppliedDeductionXP          int64            `json:"applied_deduction_xp"`
	Applied                     int64            `json:"applied_adjustments"`
	AppliedXP                   int64            `json:"applied_xp"`
	CappedDeductions            int64            `json:"capped_deductions"`
	CappedXP                    int64            `json:"capped_xp"`
	IncompleteTasks             int64            `json:"incomplete_task_assessments"`
	ShortTasks                  int64            `json:"short_task_assessments"`
	GenuineDifficultyExemptions int64            `json:"genuine_difficulty_exemptions"`
	ReplayOrZeroDeductions      int64            `json:"replay_or_zero_deductions"`
	ReplaysPrevented            int64            `json:"replays_prevented"`
	StorageErrors               int64            `json:"storage_errors"`
}

type accountabilityMetrics struct {
	mu                          sync.Mutex
	decisions                   map[string]int64
	calculatedDeductions        int64
	calculatedDeductionXP       int64
	appliedDeductions           int64
	appliedDeductionXP          int64
	applied                     int64
	appliedXP                   int64
	cappedDeductions            int64
	cappedXP                    int64
	incompleteTasks             int64
	shortTasks                  int64
	genuineDifficultyExemptions int64
	replayOrZeroDeductions      int64
	replaysPrevented            int64
	storageErrors               int64
}

func (m *accountabilityMetrics) recordDecision(d xppolicy.Decision) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !(d.Outcome == xppolicy.OutcomeFull && d.Reason == xppolicy.ReasonNone) {
		key := string(d.Outcome) + ":" + string(d.Reason)
		if m.decisions == nil {
			m.decisions = make(map[string]int64)
		}
		m.decisions[key]++
	}
	if d.CalculatedAmount < 0 {
		m.calculatedDeductions++
		m.calculatedDeductionXP += int64(-d.CalculatedAmount)
	}
	if d.Reason == xppolicy.ReasonGenuineDifficulty {
		m.genuineDifficultyExemptions++
	}
}

func (m *accountabilityMetrics) recordIncomplete() {
	m.mu.Lock()
	m.incompleteTasks++
	m.mu.Unlock()
}

func (m *accountabilityMetrics) recordShortTask() {
	m.mu.Lock()
	m.shortTasks++
	m.mu.Unlock()
}

func (m *accountabilityMetrics) recordApplied(d xppolicy.Decision, applied int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if applied != 0 {
		m.applied++
		m.appliedXP += int64(applied)
	}
	if d.CalculatedAmount < 0 {
		if applied < 0 {
			m.appliedDeductions++
			m.appliedDeductionXP += int64(-applied)
		} else {
			m.replayOrZeroDeductions++
		}
		if applied > d.CalculatedAmount {
			m.cappedDeductions++
			m.cappedXP += int64(applied - d.CalculatedAmount)
		}
	}
}

func (m *accountabilityMetrics) recordReplayPrevented() {
	m.mu.Lock()
	m.replayOrZeroDeductions++
	m.replaysPrevented++
	m.mu.Unlock()
}

func (m *accountabilityMetrics) recordStorageError() {
	m.mu.Lock()
	m.storageErrors++
	m.mu.Unlock()
}

// AccountabilityMetrics returns a copy that a handler can encode without
// holding the engine lock.
func (e *Engine) AccountabilityMetrics() AccountabilityMetricsSnapshot {
	e.accountabilityStats.mu.Lock()
	defer e.accountabilityStats.mu.Unlock()
	out := AccountabilityMetricsSnapshot{
		Decisions:                   make(map[string]int64, len(e.accountabilityStats.decisions)),
		CalculatedDeductions:        e.accountabilityStats.calculatedDeductions,
		CalculatedDeductionXP:       e.accountabilityStats.calculatedDeductionXP,
		AppliedDeductions:           e.accountabilityStats.appliedDeductions,
		AppliedDeductionXP:          e.accountabilityStats.appliedDeductionXP,
		Applied:                     e.accountabilityStats.applied,
		AppliedXP:                   e.accountabilityStats.appliedXP,
		CappedDeductions:            e.accountabilityStats.cappedDeductions,
		CappedXP:                    e.accountabilityStats.cappedXP,
		IncompleteTasks:             e.accountabilityStats.incompleteTasks,
		ShortTasks:                  e.accountabilityStats.shortTasks,
		GenuineDifficultyExemptions: e.accountabilityStats.genuineDifficultyExemptions,
		ReplayOrZeroDeductions:      e.accountabilityStats.replayOrZeroDeductions,
		ReplaysPrevented:            e.accountabilityStats.replaysPrevented,
		StorageErrors:               e.accountabilityStats.storageErrors,
	}
	for k, v := range e.accountabilityStats.decisions {
		out.Decisions[k] = v
	}
	return out
}

// AccountabilityConfig returns the resolved policy for the admin surface.
func (e *Engine) AccountabilityConfig() xppolicy.Config {
	return e.effortConfig
}
