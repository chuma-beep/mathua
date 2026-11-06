package engine

import (
	"testing"

	xppolicy "github.com/chuma-beep/mathua/internal/xp"
)

func TestAccountabilityMetricsDistinguishCalculatedAppliedAndCapped(t *testing.T) {
	e := &Engine{}
	d := xppolicy.Decision{
		Outcome:          xppolicy.OutcomeNegative,
		Amount:           -1,
		CalculatedAmount: -2,
		Reason:           xppolicy.ReasonPoorPerformance,
		Source:           xppolicy.SourcePerformance,
	}
	e.accountabilityStats.recordDecision(d)
	e.accountabilityStats.recordApplied(d, -1)
	e.accountabilityStats.recordIncomplete()
	e.accountabilityStats.recordShortTask()
	e.accountabilityStats.recordReplayPrevented()

	snap := e.AccountabilityMetrics()
	if snap.CalculatedDeductions != 1 || snap.CalculatedDeductionXP != 2 {
		t.Fatalf("calculated = %d/%d, want 1/2", snap.CalculatedDeductions, snap.CalculatedDeductionXP)
	}
	if snap.AppliedDeductions != 1 || snap.AppliedDeductionXP != 1 || snap.AppliedXP != -1 {
		t.Fatalf("applied = %d/%d/%d, want 1/1/-1", snap.AppliedDeductions, snap.AppliedDeductionXP, snap.AppliedXP)
	}
	if snap.CappedDeductions != 1 || snap.CappedXP != 1 {
		t.Fatalf("capped = %d/%d, want 1/1", snap.CappedDeductions, snap.CappedXP)
	}
	if snap.IncompleteTasks != 1 || snap.ShortTasks != 1 {
		t.Fatalf("suppression counters = %d/%d, want 1/1", snap.IncompleteTasks, snap.ShortTasks)
	}
	if snap.ReplayOrZeroDeductions != 1 || snap.ReplaysPrevented != 1 {
		t.Fatalf("replay/zero = %d/%d, want 1/1", snap.ReplayOrZeroDeductions, snap.ReplaysPrevented)
	}
	if snap.Decisions["negative:poor_performance"] != 1 {
		t.Fatalf("decision map = %v", snap.Decisions)
	}
}
