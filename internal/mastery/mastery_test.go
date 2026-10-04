package mastery

import (
	"math"
	"testing"
)

// Next() transitions

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

// EffectiveStatus

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

// SM2Quality

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

// Status constants

func TestStatusConstants(t *testing.T) {
	if StatusUnseen != "UNSEEN" || StatusLearning != "LEARNING" || StatusPracticing != "PRACTICING" || StatusMastered != "MASTERED" {
		t.Error("status constants mismatch")
	}
}

// TestMasteryPct_NeverGoesBackwards drives the real rules — streak increments on a
// correct answer and resets to 1 on every tier advance, exactly as
// Engine.submitAnswerWithTask does — and asserts the bar never decreases.
//
// This is the assertion that would have caught the sawtooth. The old formula was
// `streak / requiredStreak`, so at the instant of each advance the bar fell from 100% to
// ~10%, three times on the way up: the learner watched their progress drop while doing
// everything right.
func TestMasteryPct_NeverGoesBackwards(t *testing.T) {
	const req = 10
	m := &Machine{}
	status := StatusUnseen
	streak := 0
	prev := MasteryPct(status, streak, req)
	advances := 0

	for answer := 0; answer < 60 && status != StatusMastered; answer++ {
		streak++ // a correct answer
		next := m.Next(status, TransitionCtx{
			Streak:            streak,
			RequiredStreak:    req,
			AvgResponseTime:   1,
			ResponseThreshold: 10,
		})
		if next != status {
			status = next
			streak = 1 // the engine's reset-on-advance
			advances++
		}

		cur := MasteryPct(status, streak, req)
		if cur < prev-1e-9 {
			t.Fatalf("bar went backwards on answer %d at status %s (streak %d): %v then %v",
				answer+1, status, streak, prev, cur)
		}
		prev = cur
	}

	if status != StatusMastered {
		t.Fatalf("never reached MASTERED in 60 correct answers, got %s", status)
	}
	if advances != 3 {
		t.Errorf("expected 3 tier advances, got %d", advances)
	}
	if prev != 1 {
		t.Errorf("final bar = %v, want 1", prev)
	}
	// And the old formula, for the record: it is this sawtooth that made the streak
	// ratio unusable as a bar length.
	if old := float64(1) / float64(req); old >= 0.2 {
		t.Errorf("sanity: the streak ratio at streak 1 is %v, which is where the bar used to land", old)
	}
}

func TestMasteryPct_TierProgressFillsEachThird(t *testing.T) {
	const req = 10
	// One third per transition, filled by progress within the tier.
	if got := MasteryPct(StatusLearning, 1, req); math.Abs(got-1.1/3) > 1e-9 {
		t.Errorf("learning streak 1 = %v, want %v", got, 1.1/3)
	}
	if got := MasteryPct(StatusLearning, 10, req); math.Abs(got-2.0/3) > 1e-9 {
		t.Errorf("learning streak met = %v, want %v", got, 2.0/3)
	}
	if got := MasteryPct(StatusPracticing, 5, req); math.Abs(got-2.5/3) > 1e-9 {
		t.Errorf("practising streak 5 = %v, want %v", got, 2.5/3)
	}
}

func TestMasteryPct_DecayingKeepsAFullBar(t *testing.T) {
	// Decay is signalled by colour only. A shorter bar would read as losing the
	// competence, which is the one thing decay does not mean.
	if got := MasteryPct("DECAYING", 1, 10); got != 1 {
		t.Errorf("decaying = %v, want 1 (colour carries decay, not length)", got)
	}
}

func TestMasteryPct_NoThresholdOrNoStreak(t *testing.T) {
	// A missing threshold must not divide by zero, and a concept with a row but no
	// correct answer must not look partly done.
	if got := MasteryPct(StatusLearning, 3, 0); got != 2.0/3 {
		t.Errorf("learning with threshold 0 = %v, want %v (tier only, no within-tier credit)", got, 2.0/3)
	}
	if got := MasteryPct(StatusLearning, 0, 10); got <= 0 {
		t.Errorf("learning with no streak = %v, want > 0 (the tier itself is progress)", got)
	}
	if got := MasteryPct(StatusUnseen, 0, 0); got != 0 {
		t.Errorf("unseen = %v, want 0", got)
	}
}
