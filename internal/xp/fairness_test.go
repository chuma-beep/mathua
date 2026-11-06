package xp

import (
	"fmt"
	"testing"
)

func fairnessAttempts(n, correct int, difficulty float64) []AttemptSignal {
	out := make([]AttemptSignal, 0, n)
	for i := 0; i < n; i++ {
		c := i < correct
		xp := 0
		if c {
			xp = 3
		}
		out = append(out, AttemptSignal{
			Correct:         c,
			Difficulty:      difficulty,
			ElapsedSeconds:  10,
			AnswerSignature: fmt.Sprintf("fair-%d", i),
			XPEarned:        xp,
		})
	}
	return out
}

func fairnessInput(n, correct int, difficulty float64, startingXP int) Input {
	attempts := fairnessAttempts(n, correct, difficulty)
	gross := correct * 3
	return Input{
		TaskKind:     TaskPractice,
		TaskID:       "fairness",
		Attempts:     attempts,
		AvailableXP:  startingXP + gross,
		TaskComplete: true,
		Config:       DefaultConfig(),
	}
}

// TestAccuracyMonotonicity is the core fairness property: for a fixed task
// length, difficulty and starting balance, one more correct answer can never
// produce a lower net task XP outcome. It uses the production policy and the
// three answers per correct answer used in the simulation.
func TestAccuracyMonotonicity(t *testing.T) {
	lengths := []int{4, 5, 6, 10, 20}
	difficulties := []float64{0.30, 0.40, 0.50, 0.55, 0.60, 0.80}
	starting := []int{0, 10}

	for _, n := range lengths {
		for _, difficulty := range difficulties {
			for _, start := range starting {
				previous := 0
				for correct := 0; correct <= n; correct++ {
					in := fairnessInput(n, correct, difficulty, start)
					effort := Evaluate(in)
					if effort.Amount != 0 {
						t.Fatalf("n=%d difficulty=%.2f start=%d correct=%d: effort unexpectedly triggered: %+v",
							n, difficulty, start, correct, effort)
					}
					d := Combine(effort, EvaluatePerformance(in))
					gross := correct * 3
					net := gross + d.Amount
					if correct > 0 && net < previous {
						t.Errorf("n=%d difficulty=%.2f start=%d: correct %d has net %d below correct %d net %d (%+v)",
							n, difficulty, start, correct, net, correct-1, previous, d)
					}
					previous = net
				}
			}
		}
	}
}

// TestHardTaskExemptionDoesNotProduceNegativePenalty pins the fixed edge: at
// high difficulty, zero or one correct answer may remove earned XP but must
// not add the extra negative penalty. This keeps zero and one correct from
// producing a worse net result than zero correct.
func TestHardTaskExemptionDoesNotProduceNegativePenalty(t *testing.T) {
	for _, n := range []int{6, 10, 20} {
		for _, correct := range []int{0, 1} {
			in := fairnessInput(n, correct, 0.80, 10)
			d := Combine(Evaluate(in), EvaluatePerformance(in))
			if d.Outcome == OutcomeNegative {
				t.Errorf("n=%d correct=%d produced a negative outcome: %+v", n, correct, d)
			}
			gross := correct * 3
			if net := gross + d.Amount; net < 0 {
				t.Errorf("n=%d correct=%d net = %d, want non-negative", n, correct, net)
			}
		}
	}
}

// TestNeutralLowAccuracyKeepsNetMonotonic documents the remaining
// counterintuitive display case. Zero and one correct both net -2, but the
// adjustment differs because the one correct answer first earns XP.
func TestNeutralLowAccuracyKeepsNetMonotonic(t *testing.T) {
	zero := fairnessInput(6, 0, 0.55, 10)
	one := fairnessInput(6, 1, 0.55, 10)
	zeroDecision := Combine(Evaluate(zero), EvaluatePerformance(zero))
	oneDecision := Combine(Evaluate(one), EvaluatePerformance(one))
	zeroNet := zeroDecision.Amount
	oneNet := 3 + oneDecision.Amount
	if zeroNet != oneNet {
		t.Errorf("neutral zero/one net = %d/%d, want equal", zeroNet, oneNet)
	}
	if oneDecision.Amount >= zeroDecision.Amount {
		t.Errorf("neutral one-correct adjustment = %d, zero = %d; the display difference is expected but the net must stay equal",
			oneDecision.Amount, zeroDecision.Amount)
	}
}
