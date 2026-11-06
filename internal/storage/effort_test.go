package storage

import (
	"testing"
	"time"
)

func TestApplyXPAdjustmentCapsCalculatedMinusTwoAtBalanceOne(t *testing.T) {
	store := newTestStore(t)
	es, ok := store.(EffortStore)
	if !ok {
		t.Fatal("test store does not implement EffortStore")
	}
	st := seedStudent(t, store)
	if err := store.AwardXP(st.ID, 1, "seed", "seed:one"); err != nil {
		t.Fatalf("seed: %v", err)
	}
	applied, err := es.ApplyXPAdjustment(st.ID, -2, "accountability:extremely_poor_performance", "accountability:practice:one")
	if err != nil || applied != -1 {
		t.Fatalf("applied = %d, %v; want -1, nil", applied, err)
	}
	total, _, err := store.GetXP(st.ID)
	if err != nil {
		t.Fatalf("get xp: %v", err)
	}
	if total != 0 {
		t.Fatalf("balance = %d, want 0", total)
	}
}

func TestApplyXPAdjustmentIsIdempotentAndNeverNegative(t *testing.T) {
	store := newTestStore(t)
	es, ok := store.(EffortStore)
	if !ok {
		t.Fatal("test store does not implement EffortStore")
	}
	st := seedStudent(t, store)
	if err := store.AwardXP(st.ID, 5, "study", "award:1"); err != nil {
		t.Fatalf("seed award: %v", err)
	}

	applied, err := es.ApplyXPAdjustment(st.ID, -3, "accountability:poor_performance", "accountability:practice:t1")
	if err != nil || applied != -3 {
		t.Fatalf("first adjustment = %d, %v; want -3, nil", applied, err)
	}
	applied, err = es.ApplyXPAdjustment(st.ID, -3, "accountability:poor_performance", "accountability:practice:t1")
	if err != nil || applied != 0 {
		t.Fatalf("duplicate adjustment = %d, %v; want 0, nil", applied, err)
	}
	applied, err = es.ApplyXPAdjustment(st.ID, -10, "accountability:extremely_poor_performance", "accountability:practice:t2")
	if err != nil || applied != -2 {
		t.Fatalf("capped adjustment = %d, %v; want -2, nil", applied, err)
	}
	applied, err = es.ApplyXPAdjustment(st.ID, -1, "accountability:extremely_poor_performance", "accountability:practice:t3")
	if err != nil || applied != 0 {
		t.Fatalf("zero-floor adjustment = %d, %v; want 0, nil", applied, err)
	}

	total, _, err := store.GetXP(st.ID)
	if err != nil {
		t.Fatalf("get xp: %v", err)
	}
	if total != 0 {
		t.Fatalf("xp_total = %d, want 0", total)
	}

	applied, err = es.ApplyXPAdjustment(st.ID, 2, "accountability:exceptional_performance", "accountability:practice:t4")
	if err != nil || applied != 2 {
		t.Fatalf("bonus = %d, %v; want 2, nil", applied, err)
	}
	applied, err = es.ApplyXPAdjustment(st.ID, 2, "accountability:exceptional_performance", "accountability:practice:t4")
	if err != nil || applied != 0 {
		t.Fatalf("duplicate bonus = %d, %v; want 0, nil", applied, err)
	}

	removed, err := es.SumXPDeductionsSince(st.ID, time.Time{})
	if err != nil {
		t.Fatalf("sum deductions: %v", err)
	}
	if removed != 5 {
		t.Fatalf("deduction magnitude = %d, want 5", removed)
	}
}
