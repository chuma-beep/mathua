package planning

import (
	"testing"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/mastery"
	"github.com/chuma-beep/mathua/internal/storage"
	"github.com/chuma-beep/mathua/internal/xp"
)

// The plan's XP number must be the XP the learner actually earns.
//
// This was wrong by 1.4-1.65x. `learningCost` prices each attempt at the bare effort base
// while `xp.Award` multiplied that base by a speed factor (0.5-1.25) and a streak factor
// (1-1.5), so the planner over-predicted remaining XP for anyone already mid-streak — the
// plan said "913 XP" and the work cost less. Both multipliers are gone, so the two now agree;
// this test is what keeps them agreeing.
func TestEstimateXPMatchesTheAwardPath(t *testing.T) {
	thresholds := []float64{6, 10, 15, 20, 30, 45, 60, 90, 120, 150}
	for _, th := range thresholds {
		for _, streak := range []int{0, 1, 3, 7, 10} {
			c := &concepts.Concept{
				ID:               "c",
				MasteryThreshold: concepts.MasteryThreshold{Streak: 10, AvgTimeSeconds: th},
			}
			// A progress row mid-streak with a realistic average response time, which is
			// the case the old multipliers made the planner over-predict.
			prog := &storage.ConceptProgress{
				ConceptID: c.ID, Streak: streak, Attempts: streak,
				AvgResponseTime: th, Status: string(mastery.StatusPracticing),
			}
			cost := learningCost(c, prog, AttemptStats{})

			// What the learner actually banks for the same number of answers — the
			// remaining ones to reach this concept's threshold — answered inside the
			// concept's own time threshold, which is the "worked normally" case.
			remaining := 10 - streak
			if remaining < 2 {
				remaining = 2
			}
			var banked int
			for i := 0; i < remaining; i++ {
				banked += xp.Award(true, th, th, streak+i, xp.TaskLesson)
			}

			if cost.XP != float64(banked) {
				t.Errorf("threshold %.0fs streak %d: estimate says %.0f XP for %d answers, "+
					"award path pays %d XP", th, streak, cost.XP, remaining, banked)
			}
		}
	}
}

// A slower correct answer is full credit, so the estimate must not assume the faster case is
// the only one that counts.
func TestSlowCorrectAnswerStillPaysTheBase(t *testing.T) {
	th := 60.0
	nominal := xp.Award(true, th, th, 0, xp.TaskLesson)
	slow := xp.Award(true, 10*th, th, 0, xp.TaskLesson)
	if slow != nominal {
		t.Errorf("a correct answer 10x over the threshold paid %d XP, nominal paid %d — slow "+
			"correct work is full credit", slow, nominal)
	}
}

// XP is effort, so the base must rise with the concept's expected time. This is the property
// that makes "1 XP ~ 1 minute" mean anything, and it is what the whole effort model is for.
func TestXPIsEffortDerivedNotFlatPerQuestion(t *testing.T) {
	quick := xp.EffortBase(10, 0.5)
	slow := xp.EffortBase(150, 0.5)
	if slow <= quick {
		t.Errorf("a 150s concept prices at %d XP and a 10s one at %d — the effort base is flat "+
			"per question, which is the thing this model exists to avoid", slow, quick)
	}
	if quick < 1 || slow > 5 {
		t.Errorf("bases out of the documented 1-5 range: %d and %d", quick, slow)
	}
}
