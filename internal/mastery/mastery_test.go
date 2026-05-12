package mastery

import "testing"

// ---------------------------------------------------------------------------
// Next() transitions
// ---------------------------------------------------------------------------

func TestNext_UnseenToLearning(t *testing.T) {
	m := &Machine{}
	ctx := TransitionCtx{
		Streak:            5,
		RequiredStreak:    5,
		AvgResponseTime:   3.0,
		ResponseThreshold: 10.0,
	}
	got := m.Next(StatusUnseen, ctx)
	if got != StatusLearning {
		t.Errorf("UNSEEN → LEARNING: expected %q, got %q", StatusLearning, got)
	}
}

func TestNext_UnseenStays_StreakNotMet(t *testing.T) {
	m := &Machine{}
	ctx := TransitionCtx{
		Streak:            3,
		RequiredStreak:    5,
		AvgResponseTime:   3.0,
		ResponseThreshold: 10.0,
	}
	got := m.Next(StatusUnseen, ctx)
	if got != StatusUnseen {
		t.Errorf("expected UNSEEN, got %q", got)
	}
}

func TestNext_UnseenStays_TimeTooSlow(t *testing.T) {
	m := &Machine{}
	ctx := TransitionCtx{
		Streak:            5,
		RequiredStreak:    5,
		AvgResponseTime:   12.0,
		ResponseThreshold: 10.0,
	}
	got := m.Next(StatusUnseen, ctx)
	if got != StatusUnseen {
		t.Errorf("expected UNSEEN, got %q", got)
	}
}

func TestNext_UnseenStays_NeitherMet(t *testing.T) {
	m := &Machine{}
	ctx := TransitionCtx{
		Streak:            2,
		RequiredStreak:    5,
		AvgResponseTime:   15.0,
		ResponseThreshold: 10.0,
	}
	got := m.Next(StatusUnseen, ctx)
	if got != StatusUnseen {
		t.Errorf("expected UNSEEN, got %q", got)
	}
}

func TestNext_LearningToPracticing(t *testing.T) {
	m := &Machine{}
	ctx := TransitionCtx{
		Streak:            5,
		RequiredStreak:    5,
		AvgResponseTime:   4.0,
		ResponseThreshold: 8.0,
	}
	got := m.Next(StatusLearning, ctx)
	if got != StatusPracticing {
		t.Errorf("LEARNING → PRACTICING: expected %q, got %q", StatusPracticing, got)
	}
}

func TestNext_LearningStays_StreakNotMet(t *testing.T) {
	m := &Machine{}
	ctx := TransitionCtx{
		Streak:            4,
		RequiredStreak:    5,
		AvgResponseTime:   4.0,
		ResponseThreshold: 8.0,
	}
	got := m.Next(StatusLearning, ctx)
	if got != StatusLearning {
		t.Errorf("expected LEARNING, got %q", got)
	}
}

func TestNext_PracticingToMastered(t *testing.T) {
	m := &Machine{}
	ctx := TransitionCtx{
		Streak:            7,
		RequiredStreak:    7,
		AvgResponseTime:   5.0,
		ResponseThreshold: 6.0,
	}
	got := m.Next(StatusPracticing, ctx)
	if got != StatusMastered {
		t.Errorf("PRACTICING → MASTERED: expected %q, got %q", StatusMastered, got)
	}
}

func TestNext_PracticingStays_StreakBroken(t *testing.T) {
	m := &Machine{}
	ctx := TransitionCtx{
		Streak:            6,
		RequiredStreak:    7,
		AvgResponseTime:   5.0,
		ResponseThreshold: 6.0,
	}
	got := m.Next(StatusPracticing, ctx)
	if got != StatusPracticing {
		t.Errorf("expected PRACTICING, got %q", got)
	}
}

func TestNext_PracticingStays_TimeTooSlow(t *testing.T) {
	m := &Machine{}
	ctx := TransitionCtx{
		Streak:            7,
		RequiredStreak:    7,
		AvgResponseTime:   9.0,
		ResponseThreshold: 6.0,
	}
	got := m.Next(StatusPracticing, ctx)
	if got != StatusPracticing {
		t.Errorf("expected PRACTICING, got %q", got)
	}
}

func TestNext_MasteredStays(t *testing.T) {
	m := &Machine{}
	ctx := TransitionCtx{
		Streak:            10,
		RequiredStreak:    7,
		AvgResponseTime:   2.0,
		ResponseThreshold: 6.0,
	}
	got := m.Next(StatusMastered, ctx)
	if got != StatusMastered {
		t.Errorf("expected MASTERED, got %q", got)
	}
}

// Streak exceeds required — still advances (excess is irrelevant).
func TestNext_ExcessStreak(t *testing.T) {
	m := &Machine{}
	ctx := TransitionCtx{
		Streak:            12,
		RequiredStreak:    5,
		AvgResponseTime:   3.0,
		ResponseThreshold: 10.0,
	}
	got := m.Next(StatusUnseen, ctx)
	if got != StatusLearning {
		t.Errorf("expected LEARNING (excess streak ok), got %q", got)
	}
}

// Exact boundary: streak == required, time == threshold.
func TestNext_ExactBoundary(t *testing.T) {
	m := &Machine{}
	ctx := TransitionCtx{
		Streak:            5,
		RequiredStreak:    5,
		AvgResponseTime:   10.0,
		ResponseThreshold: 10.0,
	}
	got := m.Next(StatusLearning, ctx)
	if got != StatusPracticing {
		t.Errorf("expected PRACTICING (exact boundary ok), got %q", got)
	}
}

// ---------------------------------------------------------------------------
// EffectiveStatus
// ---------------------------------------------------------------------------

func TestEffectiveStatus_Decaying(t *testing.T) {
	if got := EffectiveStatus(StatusMastered, 15, 14); got != "DECAYING" {
		t.Errorf("expected DECAYING, got %q", got)
	}
}

func TestEffectiveStatus_ExactBoundaryDecay(t *testing.T) {
	if got := EffectiveStatus(StatusMastered, 14, 14); got != "DECAYING" {
		t.Errorf("expected DECAYING at exact boundary, got %q", got)
	}
}

func TestEffectiveStatus_NotDecaying(t *testing.T) {
	if got := EffectiveStatus(StatusMastered, 13, 14); got != StatusMastered {
		t.Errorf("expected MASTERED, got %q", got)
	}
}

func TestEffectiveStatus_NonMasteredUnaffected(t *testing.T) {
	for _, s := range []Status{StatusUnseen, StatusLearning, StatusPracticing} {
		if got := EffectiveStatus(s, 30, 14); got != s {
			t.Errorf("expected %q to be unchanged, got %q", s, got)
		}
	}
}

// ---------------------------------------------------------------------------
// SM2Quality
// ---------------------------------------------------------------------------

func TestSM2Quality_Perfect(t *testing.T) {
	if q := SM2Quality(true, 0.3); q != 5 {
		t.Errorf("expected quality 5, got %d", q)
	}
}

func TestSM2Quality_FastBoundary(t *testing.T) {
	if q := SM2Quality(true, 0.5); q != 5 {
		t.Errorf("expected quality 5 at boundary, got %d", q)
	}
}

func TestSM2Quality_Normal(t *testing.T) {
	if q := SM2Quality(true, 0.7); q != 4 {
		t.Errorf("expected quality 4, got %d", q)
	}
}

func TestSM2Quality_Slow(t *testing.T) {
	if q := SM2Quality(true, 1.5); q != 3 {
		t.Errorf("expected quality 3, got %d", q)
	}
}

func TestSM2Quality_Failed(t *testing.T) {
	if q := SM2Quality(false, 0.3); q != 0 {
		t.Errorf("expected quality 0 when streak not met, got %d", q)
	}
}

func TestSM2Quality_ExactOne(t *testing.T) {
	if q := SM2Quality(true, 1.0); q != 4 {
		t.Errorf("expected quality 4 at exact 1.0, got %d", q)
	}
}

// ---------------------------------------------------------------------------
// Status constants
// ---------------------------------------------------------------------------

func TestStatusConstants(t *testing.T) {
	if StatusUnseen != "UNSEEN" || StatusLearning != "LEARNING" || StatusPracticing != "PRACTICING" || StatusMastered != "MASTERED" {
		t.Error("status constants mismatch")
	}
}
