package xp

import "testing"

// Effort calibration: threshold seconds scaled by difficulty, 1-5 XP.
func TestEffortBase(t *testing.T) {
	cases := []struct {
		thresh, diff float64
		want         int
	}{
		{10, 0.5, 1},  // 10s question prices at 1 XP (bottom compression)
		{25, 0.5, 1},  // 25s still rounds down: acceptable per ADR
		{60, 0.5, 1},  // 1 minute at neutral difficulty
		{60, 0.8, 1},  // harder, still 1 (rounds 1.3 down)
		{120, 0.8, 3}, // 2 minutes hard
		{300, 1.0, 5}, // 5 minutes clamps at 5
		{0, 0.5, 1},   // degenerate threshold falls back, never 0
		{60, 9.9, 2},  // difficulty clamps to [0,1]
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

// There is exactly one multiplier, and it has two values.
//
// The old policy had a speed multiplier worth up to 1.25x for a fast answer and 0.5x for a
// slow one, plus a streak multiplier up to 1.5x — so one correct answer could pay 1.875x its
// base. Both are gone: speed now only answers "did they work normally", and the streak pays
// nothing.
func TestPerformanceMultiplierHasExactlyTwoValues(t *testing.T) {
	if m := PerformanceMultiplier(0.1, 10.0); m != 1.1 {
		t.Errorf("fast-but-inside-threshold multiplier = %v, want 1.1", m)
	}
	if m := PerformanceMultiplier(10.0, 10.0); m != 1.1 {
		t.Errorf("exactly-at-threshold multiplier = %v, want 1.1", m)
	}
	if m := PerformanceMultiplier(10.1, 10.0); m != 1.0 {
		t.Errorf("just-past-threshold multiplier = %v, want 1.0", m)
	}
	if m := PerformanceMultiplier(100.0, 10.0); m != 1.0 {
		t.Errorf("slow multiplier = %v, want 1.0 — a slow correct answer is still full credit", m)
	}
	// A zero or absent threshold must not silently grant the top tier.
	if m := PerformanceMultiplier(1.0, 0); m != 1.0 {
		t.Errorf("multiplier with no threshold = %v, want 1.0", m)
	}
}

func TestAwardMultiplierCeiling(t *testing.T) {
	// The product's table tops out at 1.1x. With the streak bonus gone there is nothing to
	// stack, so 1.1 really is the ceiling.
	for _, streak := range []int{0, 3, 10, 50, 500} {
		for _, elapsed := range []float64{0.5, 5, 12, 60} {
			base := float64(EffortBase(30, 0.5))
			got := float64(Award(true, elapsed, 30, streak, TaskLesson))
			if got > base*1.1+1e-9 {
				t.Errorf("Award(streak=%d, elapsed=%v) = %v, above base x 1.1 = %v",
					streak, elapsed, got, base*1.1)
			}
		}
	}
}

// The task type must not move the award. Reviews used to cost half and quizzes carried a
// +1 premium, which made routing around the main loop pay.
func TestTaskTypeDoesNotChangeTheBase(t *testing.T) {
	for _, task := range []string{TaskLesson, TaskReview, TaskMultistep, TaskQuiz} {
		want := Award(true, 10.0, 10.0, 0, TaskLesson)
		if got := Award(true, 10.0, 10.0, 0, task); got != want {
			t.Errorf("Award(%s) = %d, want %d — all task types earn the same effort base", task, got, want)
		}
	}
}

func TestAwardNominalTenSecondLesson(t *testing.T) {
	// base 1 x 1.1 truncates to 1.
	if got := Award(true, 10.0, 10.0, 0, TaskLesson); got != 1 {
		t.Errorf("Award(nominal 10s lesson) = %d, want 1", got)
	}
	// A longer concept earns more, which is the whole point of an effort base.
	if two := EffortBase(120, 0.5); two <= EffortBase(20, 0.5) {
		t.Errorf("a 120s concept priced at %d XP is not more than a 20s one at %d XP",
			two, EffortBase(20, 0.5))
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

// The cadence the daily goal implies, stated as a measurement rather than a hope.
//
// ADR-020 rescaled the quiz gate 150 -> 50 alongside a daily default of 30 -> 10, to keep
// "~5 lessons per quiz, ~a lesson a day". Raising the daily goal back to 30 while leaving
// the gate at 50 does not preserve that: at ~1 XP per question the gate is reached in a
// couple of days. This test reports the numbers so a future change to either constant has
// to look at them.
//
// It asserts only the *inputs* are self-consistent, not a cadence target, because the cadence
// is a product decision and not something a test should quietly re-decide.
func TestDailyGoalImpliedCadence(t *testing.T) {
	// A five-question lesson at a typical 20s threshold earns this much.
	var lesson int
	for i := 0; i < 5; i++ {
		lesson += Award(true, 14, 20, i, TaskLesson)
	}
	t.Logf("lesson (5 questions, 20s concepts)   = %d XP", lesson)

	const gate = QuizGateXP
	t.Logf("daily goal %d XP/day -> %.1f lessons/day, quiz gate %d XP reached in %.1f days",
		DefaultDailyGoal, float64(DefaultDailyGoal)/float64(lesson), gate,
		float64(gate)/float64(DefaultDailyGoal))
	t.Logf("default goal is inside the range: %d..%d", MinDailyGoal, MaxDailyGoal)

	if DefaultDailyGoal < MinDailyGoal || DefaultDailyGoal > MaxDailyGoal {
		t.Errorf("DefaultDailyGoal %d is outside MinDailyGoal..MaxDailyGoal", DefaultDailyGoal)
	}
	if MinDailyGoal >= MaxDailyGoal {
		t.Errorf("goal range is empty: %d..%d", MinDailyGoal, MaxDailyGoal)
	}
	// A goal must be able to out-earn a single lesson, or the learner clears it before the
	// unit of work it is meant to pace.
	if DefaultDailyGoal <= lesson {
		t.Errorf("DefaultDailyGoal %d is not greater than one lesson (%d XP)", DefaultDailyGoal, lesson)
	}
}
