package xp

import "testing"

// Audit table: base XP per task type must match the award paths.
func TestBaseXP(t *testing.T) {
	cases := map[string]int{
		TaskLesson:    10,
		TaskReview:    5,
		TaskMultistep: 15,
		TaskQuiz:      20,
		"unknown":     10,
		"":            10,
	}
	for task, want := range cases {
		if got := BaseXP(task); got != want {
			t.Errorf("BaseXP(%q) = %d, want %d", task, got, want)
		}
	}
}

// Misses earn nothing on every path; the rush penalty is separate.
func TestAwardMissIsZero(t *testing.T) {
	for _, task := range []string{TaskLesson, TaskReview, TaskMultistep, TaskQuiz} {
		if got := Award(false, 3.0, 10.0, 5, task); got != 0 {
			t.Errorf("Award(false, ..., %q) = %d, want 0", task, got)
		}
	}
	if RushPenalty != -5 {
		t.Errorf("RushPenalty = %d, want -5 (engine rush paths)", RushPenalty)
	}
}

// Multipliers stay bounded: time [0.5,1.5], streak (1,2].
func TestMultiplierBounds(t *testing.T) {
	if m := TimeMultiplier(0.1, 10.0); m != 1.5 {
		t.Errorf("fast TimeMultiplier = %v, want 1.5", m)
	}
	if m := TimeMultiplier(100.0, 10.0); m != 0.5 {
		t.Errorf("slow TimeMultiplier = %v, want 0.5", m)
	}
	if m := StreakMultiplier(0); m != 1.0 {
		t.Errorf("zero StreakMultiplier = %v, want 1.0", m)
	}
	if m := StreakMultiplier(50); m != 2.0 {
		t.Errorf("capped StreakMultiplier = %v, want 2.0", m)
	}
	// A nominal correct lesson answer earns a bounded, positive amount.
	got := Award(true, 5.0, 10.0, 0, TaskLesson)
	if got <= 0 || got > 15 {
		t.Errorf("Award(true lesson) = %d, want in (0,15]", got)
	}
}

// Estimators must use expectations, never face value.
func TestExpectedXPLessThanFace(t *testing.T) {
	face := float64(Award(true, 5.0, 10.0, 0, TaskLesson))
	exp := ExpectedXP(0.5, 5.0, 10.0, 0, TaskLesson)
	if exp >= face {
		t.Errorf("ExpectedXP(0.5) = %v, want < face %v", exp, face)
	}
	if got := ExpectedXP(1.0, 5.0, 10.0, 0, TaskLesson); got != face {
		t.Errorf("ExpectedXP(1.0) = %v, want face %v", got, face)
	}
	if got := ExpectedXP(0.0, 5.0, 10.0, 0, TaskLesson); got != 0 {
		t.Errorf("ExpectedXP(0.0) = %v, want 0", got)
	}
}
