package xp

import "testing"

// Effort calibration: threshold seconds scaled by difficulty, 1-5 XP.
func TestEffortBase(t *testing.T) {
	cases := []struct {
		thresh, diff float64
		want         int
	}{
		{10, 0.5, 1},   // 10s question prices at 1 XP (bottom compression)
		{25, 0.5, 1},   // 25s still rounds down: acceptable per ADR
		{60, 0.5, 1},   // 1 minute at neutral difficulty
		{60, 0.8, 1},   // harder, still 1 (rounds 1.3 down)
		{120, 0.8, 3},  // 2 minutes hard
		{300, 1.0, 5},  // 5 minutes clamps at 5
		{0, 0.5, 1},    // degenerate threshold falls back, never 0
		{60, 9.9, 2},   // difficulty clamps to [0,1]
	}
	for _, c := range cases {
		if got := EffortBase(c.thresh, c.diff); got != c.want {
			t.Errorf("EffortBase(%v,%v) = %d, want %d", c.thresh, c.diff, got, c.want)
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
	if RushPenalty != -1 {
		t.Errorf("RushPenalty = %d, want -1 (engine rush paths)", RushPenalty)
	}
}

// Multipliers stay bounded: time [0.5,1.25], streak (1,1.5].
// A nominal 10s lesson answer earns exactly its 1 XP base.
func TestMultiplierBounds(t *testing.T) {
	if m := TimeMultiplier(0.1, 10.0); m != 1.25 {
		t.Errorf("fast TimeMultiplier = %v, want 1.25", m)
	}
	if m := TimeMultiplier(100.0, 10.0); m != 0.5 {
		t.Errorf("slow TimeMultiplier = %v, want 0.5", m)
	}
	if m := StreakMultiplier(0); m != 1.0 {
		t.Errorf("zero StreakMultiplier = %v, want 1.0", m)
	}
	if m := StreakMultiplier(50); m != 1.5 {
		t.Errorf("capped StreakMultiplier = %v, want 1.5", m)
	}
	if got := Award(true, 10.0, 10.0, 0, TaskLesson); got != 1 {
		t.Errorf("Award(nominal 10s lesson) = %d, want 1", got)
	}
	// Quiz carries an assessment premium over the same threshold.
	if q, l := Award(true, 10.0, 10.0, 0, TaskQuiz), Award(true, 10.0, 10.0, 0, TaskLesson); q <= l {
		t.Errorf("quiz %d must exceed lesson %d at equal threshold", q, l)
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
