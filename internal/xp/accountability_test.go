package xp

import "testing"

func attempt(correct bool, difficulty, elapsed float64, xp int, signature string) AttemptSignal {
	return AttemptSignal{
		Correct:         correct,
		Difficulty:      difficulty,
		ElapsedSeconds:  elapsed,
		XPEarned:        xp,
		AnswerSignature: signature,
	}
}

func performanceInput(attempts []AttemptSignal) Input {
	return Input{
		TaskKind:     TaskPractice,
		TaskID:       "task-1",
		Attempts:     attempts,
		AvailableXP:  100,
		TaskComplete: true,
		Config:       DefaultConfig(),
	}
}

func TestSingleWrongAnswerDoesNotDeduct(t *testing.T) {
	in := performanceInput([]AttemptSignal{attempt(false, 0.5, 10, 0, "wrong")})
	if d := EvaluatePerformance(in); d.Amount != 0 {
		t.Fatalf("single wrong answer deducted %d XP: %+v", d.Amount, d)
	}
}

func TestGoodPerformanceEarnsFullXP(t *testing.T) {
	in := performanceInput([]AttemptSignal{
		attempt(true, 0.3, 10, 1, "a"),
		attempt(true, 0.3, 10, 1, "b"),
		attempt(true, 0.3, 10, 1, "c"),
		attempt(false, 0.3, 10, 0, "d"),
	})
	d := EvaluatePerformance(in)
	if d.Amount != 0 || d.Outcome != OutcomeFull {
		t.Fatalf("good performance = %+v, want full XP", d)
	}
}

func TestAcceptablePerformanceEarnsPartialXP(t *testing.T) {
	in := performanceInput([]AttemptSignal{
		attempt(true, 0.6, 10, 3, "a"),
		attempt(true, 0.6, 10, 0, "b"),
		attempt(false, 0.6, 10, 0, "c"),
		attempt(false, 0.6, 10, 0, "d"),
	})
	d := EvaluatePerformance(in)
	if d.Outcome != OutcomePartial || d.Amount != -1 {
		t.Fatalf("acceptable performance = %+v, want partial -1", d)
	}
	if d.Source != SourcePerformance {
		t.Fatalf("source = %q, want performance", d.Source)
	}
}

func TestPoorPerformanceEarnsZeroXP(t *testing.T) {
	in := performanceInput([]AttemptSignal{
		attempt(true, 0.6, 10, 3, "a"),
		attempt(true, 0.6, 10, 0, "b"),
		attempt(false, 0.6, 10, 0, "c"),
		attempt(false, 0.6, 10, 0, "d"),
		attempt(false, 0.6, 10, 0, "e"),
		attempt(false, 0.6, 10, 0, "f"),
	})
	d := EvaluatePerformance(in)
	if d.Outcome != OutcomeZero || d.Amount != -3 {
		t.Fatalf("poor performance = %+v, want zero task XP", d)
	}
}

func TestExtremelyPoorPerformanceCanGoNegative(t *testing.T) {
	var attempts []AttemptSignal
	for i := 0; i < 6; i++ {
		attempts = append(attempts, attempt(false, 0.3, 10, 0, string(rune('a'+i))))
	}
	d := EvaluatePerformance(performanceInput(attempts))
	if d.Outcome != OutcomeNegative || d.Amount != -2 {
		t.Fatalf("extremely poor performance = %+v, want -2", d)
	}
}

func TestNegativePerformanceIsCappedByBalance(t *testing.T) {
	var attempts []AttemptSignal
	for i := 0; i < 6; i++ {
		attempts = append(attempts, attempt(false, 0.3, 10, 0, string(rune('a'+i))))
	}
	in := performanceInput(attempts)
	in.AvailableXP = 1
	d := EvaluatePerformance(in)
	if d.Amount != -1 {
		t.Fatalf("capped negative = %+v, want -1", d)
	}
}

func TestHardTaskIsNotPenalizedForLegitimateDifficulty(t *testing.T) {
	var attempts []AttemptSignal
	for i := 0; i < 4; i++ {
		attempts = append(attempts, attempt(false, 0.8, 30, 0, string(rune('a'+i))))
	}
	in := performanceInput(attempts)
	if d := EvaluatePerformance(in); d.Amount != 0 {
		t.Fatalf("hard wrong answers deducted %d XP: %+v", d.Amount, d)
	}
}

func TestThreeWrongAnswersAreNotEnoughForPerformance(t *testing.T) {
	var attempts []AttemptSignal
	for i := 0; i < 3; i++ {
		attempts = append(attempts, attempt(false, 0.5, 10, 0, string(rune('a'+i))))
	}
	in := performanceInput(attempts)
	if d := EvaluatePerformance(in); d.Amount != 0 {
		t.Fatalf("three wrong answers deducted %d XP: %+v", d.Amount, d)
	}
}

func TestEffortAndPerformanceDoNotDoubleCharge(t *testing.T) {
	var attempts []AttemptSignal
	for i := 0; i < 4; i++ {
		attempts = append(attempts, attempt(false, 0.5, 1, 0, string(rune('a'+i))))
	}
	in := performanceInput(attempts)
	effort := Evaluate(in)
	performance := EvaluatePerformance(in)
	combined := Combine(effort, performance)
	if combined.Amount != -1 {
		t.Fatalf("combined = %+v, want one -1 adjustment", combined)
	}
	if combined.Source != SourceEffort {
		t.Fatalf("combined source = %q, want effort", combined.Source)
	}
}

func TestRepeatedGuessingStillDetected(t *testing.T) {
	var attempts []AttemptSignal
	for i := 0; i < 4; i++ {
		attempts = append(attempts, attempt(false, 0.5, 3, 0, "same"))
	}
	d := Evaluate(performanceInput(attempts))
	if d.Reason != ReasonRepeatedGuessing || d.Amount >= 0 {
		t.Fatalf("repeated guessing = %+v, want a negative repeated_guessing decision", d)
	}
}

func TestSlowRepeatedWrongAnswersAreNotGuessing(t *testing.T) {
	var attempts []AttemptSignal
	for i := 0; i < 4; i++ {
		attempts = append(attempts, attempt(false, 0.5, 10, 0, "same"))
	}
	d := Evaluate(performanceInput(attempts))
	if d.Reason == ReasonRepeatedGuessing {
		t.Fatalf("slow repeated answers were treated as guessing: %+v", d)
	}
}

func TestLegacyRushPenaltySuppressesOtherAdjustments(t *testing.T) {
	in := performanceInput([]AttemptSignal{{
		Correct:        false,
		Difficulty:     0.5,
		ElapsedSeconds: 1,
		LegacyPenalty:  true,
	}})
	effort := Evaluate(in)
	performance := EvaluatePerformance(in)
	combined := Combine(effort, performance)
	if combined.Reason != ReasonLegacyPenalty || combined.Source != SourceLegacy || combined.Amount != 0 {
		t.Fatalf("legacy precedence = %+v", combined)
	}
}

func TestExceptionalPerformanceBonusWhenEnabled(t *testing.T) {
	cfg := DefaultConfig()
	cfg.EnableBonus = true
	in := performanceInput(nil)
	in.Config = cfg
	for i := 0; i < 5; i++ {
		in.Attempts = append(in.Attempts, attempt(true, 0.8, 10, 1, string(rune('a'+i))))
	}
	d := EvaluatePerformance(in)
	if d.Outcome != OutcomeBonus || d.Amount != 1 {
		t.Fatalf("exceptional performance = %+v, want +1", d)
	}
}

func TestDiagnosticIsExempt(t *testing.T) {
	in := performanceInput([]AttemptSignal{attempt(false, 0.5, 1, 0, "a")})
	in.TaskKind = TaskDiagnostic
	if d := EvaluatePerformance(in); d.Reason != ReasonDiagnosticExempt {
		t.Fatalf("diagnostic = %+v, want diagnostic_exempt", d)
	}
}

func TestDisabledPolicyPreservesExistingXP(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Disabled = true
	in := performanceInput([]AttemptSignal{attempt(false, 0.5, 1, 0, "a")})
	in.Config = cfg
	if d := EvaluatePerformance(in); d.Amount != 0 || d.Outcome != OutcomeFull {
		t.Fatalf("disabled performance = %+v, want full", d)
	}
	if d := Evaluate(in); d.Amount != 0 || d.Outcome != OutcomeFull {
		t.Fatalf("disabled effort = %+v, want full", d)
	}
}

func TestNormalizePreservesZeroPenaltyAndFillsSameWrongSeconds(t *testing.T) {
	cfg := NormalizeConfig(Config{PenaltyFirst: 0})
	if cfg.PenaltyFirst != 0 {
		t.Fatalf("PenaltyFirst = %d, want intentional zero", cfg.PenaltyFirst)
	}
	if cfg.SameWrongSeconds <= 0 {
		t.Fatalf("SameWrongSeconds = %v, want a default", cfg.SameWrongSeconds)
	}
}

func TestUnfinishedTaskCannotReceivePerformanceDeduction(t *testing.T) {
	var attempts []AttemptSignal
	for i := 0; i < 6; i++ {
		attempts = append(attempts, attempt(false, 0.3, 10, 0, string(rune('a'+i))))
	}
	in := performanceInput(attempts)
	in.TaskComplete = false
	if d := EvaluatePerformance(in); d.Amount != 0 || d.Outcome != OutcomeFull {
		t.Fatalf("unfinished task = %+v, want full", d)
	}
}

func TestDifficultyPriorUsesObservedBands(t *testing.T) {
	cases := []struct {
		difficulty float64
		want       float64
	}{
		{0, 0.50},
		{0.3, 0.83},
		{0.4, 0.89},
		{0.55, 0.944},
		{0.6, 0.85},
		{0.8, 0.872},
	}
	for _, tc := range cases {
		if got := difficultyCorrectPrior(tc.difficulty); got != tc.want {
			t.Errorf("difficultyCorrectPrior(%v) = %v, want %v", tc.difficulty, got, tc.want)
		}
	}
}

func TestOneCorrectOnHardTaskDoesNotCauseNegativePenalty(t *testing.T) {
	attempts := []AttemptSignal{attempt(true, 0.8, 20, 3, "correct")}
	for i := 0; i < 5; i++ {
		attempts = append(attempts, attempt(false, 0.8, 20, 0, string(rune('a'+i))))
	}
	d := EvaluatePerformance(performanceInput(attempts))
	if d.Outcome != OutcomeZero || d.Amount != -3 {
		t.Fatalf("hard task with one correct = %+v, want zeroed earned XP without the extra penalty", d)
	}
	if d.Reason != ReasonGenuineDifficulty {
		t.Fatalf("reason = %q, want genuine_difficulty", d.Reason)
	}
}

func TestHarderBandIsNotMorePenalizedWhenAllWrong(t *testing.T) {
	var medium, hard []AttemptSignal
	for i := 0; i < 6; i++ {
		medium = append(medium, attempt(false, 0.6, 20, 0, string(rune('a'+i))))
		hard = append(hard, attempt(false, 0.8, 20, 0, string(rune('b'+i))))
	}
	md := EvaluatePerformance(performanceInput(medium))
	hd := EvaluatePerformance(performanceInput(hard))
	if md.Amount >= 0 {
		t.Fatalf("medium all-wrong should be a deduction, got %+v", md)
	}
	if hd.Amount != 0 {
		t.Fatalf("hard all-wrong should be exempt as genuine difficulty, got %+v", hd)
	}
}

func TestExtremelyPoorPerformanceNeedsNoEffortTrigger(t *testing.T) {
	var attempts []AttemptSignal
	for i := 0; i < 6; i++ {
		attempts = append(attempts, attempt(false, 0.3, 10, 0, string(rune('a'+i))))
	}
	in := performanceInput(attempts)
	effort := Evaluate(in)
	performance := EvaluatePerformance(in)
	combined := Combine(effort, performance)
	if effort.Amount != 0 {
		t.Fatalf("effort unexpectedly triggered: %+v", effort)
	}
	if performance.Amount != -2 {
		t.Fatalf("performance = %+v, want -2", performance)
	}
	if combined.Amount != -2 || combined.Source != SourcePerformance {
		t.Fatalf("combined = %+v, want a -2 performance adjustment", combined)
	}
}

func TestFiveWrongAnswersDoNotProduceNegativeXP(t *testing.T) {
	var attempts []AttemptSignal
	for i := 0; i < 5; i++ {
		attempts = append(attempts, attempt(false, 0.3, 10, 0, string(rune('a'+i))))
	}
	d := EvaluatePerformance(performanceInput(attempts))
	if d.Amount != 0 || d.Outcome == OutcomeNegative {
		t.Fatalf("five wrong answers = %+v, want no negative adjustment", d)
	}
}

func TestExtremelyPoorPerformanceRemovesEarnedXPAndPenalty(t *testing.T) {
	attempts := []AttemptSignal{attempt(true, 0.3, 10, 3, "correct")}
	for i := 0; i < 6; i++ {
		attempts = append(attempts, attempt(false, 0.3, 10, 0, string(rune('a'+i))))
	}
	d := EvaluatePerformance(performanceInput(attempts))
	if d.Amount != -5 || d.Outcome != OutcomeNegative {
		t.Fatalf("earned-XP task = %+v, want -5", d)
	}
}

func TestSixCorrectAnswersRemainNonNegative(t *testing.T) {
	var attempts []AttemptSignal
	for i := 0; i < 6; i++ {
		attempts = append(attempts, attempt(true, 0.6, 10, 3, string(rune('a'+i))))
	}
	d := EvaluatePerformance(performanceInput(attempts))
	if d.Amount < 0 {
		t.Fatalf("six correct answers = %+v, want a non-negative outcome", d)
	}
}
