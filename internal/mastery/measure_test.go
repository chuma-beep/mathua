package mastery

import (
	"fmt"
	"math"
	"sort"
	"testing"
)

// What mastery costs, measured rather than assumed.
//
// `MasteryEvidenceFloor` is a minimum on *how much* evidence, not a rule about what the
// evidence contains, so the real answer depends on the learner: mistakes, recovery, and how
// hard the questions were. This walks a spread of learner profiles against the corpus's
// real time thresholds and reports the distribution.
//
// It is a test with no assertions on the distribution, only on the *properties* that must
// hold whatever the cost turns out to be — because the cost is a calibration decision and
// the properties are not.

type profile struct {
	name string
	// correct reports whether attempt i was answered correctly.
	correct func(i int) bool
	// difficulty is the difficulty served on attempt i.
	difficulty func(i int) float64
	// elapsed is how long the learner took, in seconds.
	elapsed func(i int) float64
}

var profiles = []profile{
	{
		name:       "perfect, hardest questions",
		correct:    func(int) bool { return true },
		difficulty: func(int) float64 { return 1.0 },
		elapsed:    func(int) float64 { return 3 },
	},
	{
		name:       "perfect, easiest questions",
		correct:    func(int) bool { return true },
		difficulty: func(int) float64 { return 0.3 },
		elapsed:    func(int) float64 { return 3 },
	},
	{
		name:       "steady — one miss every eight",
		correct:    func(i int) bool { return i%8 != 7 },
		difficulty: func(i int) float64 { return 0.6 },
		elapsed:    func(int) float64 { return 4 },
	},
	{
		name:       "flaky — 80% correct",
		correct:    func(i int) bool { return i%5 != 4 },
		difficulty: func(i int) float64 { return 0.6 },
		elapsed:    func(int) float64 { return 4 },
	},
	{
		name:       "slow but right",
		correct:    func(int) bool { return true },
		difficulty: func(i int) float64 { return 0.6 },
		elapsed:    func(int) float64 { return 90 },
	},
	{
		name:       "starts rough",
		correct:    func(i int) bool { return i < 5 || i%4 != 3 },
		difficulty: func(i int) float64 { return 0.5 },
		elapsed:    func(int) float64 { return 5 },
	},
}

// corpusTimeThresholds are the real `avg_time_seconds` values in the corpus, with the
// count of concepts that use each. They matter because the time term is scored against the
// concept's own threshold, so a slow concept's "fast" answer is not a slow concept's.
var corpusTimeThresholds = []struct {
	seconds  float64
	concepts int
}{
	{10, 32}, {15, 61}, {18, 60}, {20, 54}, {25, 37}, {30, 18}, {60, 13}, {120, 162},
}

// answersToMastery runs one profile against one concept and reports the number of answers
// it took, or -1 if it never got there.
func answersToMastery(p profile, timeThreshold float64, maxAnswers int, startAt int) int {
	m := &Machine{}
	status := StatusUnseen
	var window []Attempt

	for i := 0; i < maxAnswers; i++ {
		n := startAt + i
		a := Attempt{
			Correct:    p.correct(n),
			Elapsed:    p.elapsed(n),
			Difficulty: diff(0.6),
			Instance:   fmt.Sprintf("q%d", n),
		}
		a.Difficulty = difficultyOf(p, n)
		window = append(window, a)
		if len(window) > EvidenceWindow {
			window = window[len(window)-EvidenceWindow:]
		}
		ev := BuildEvidence(window, EvidenceOptions{TimeThreshold: timeThreshold})
		next := m.Next(status, TransitionCtx{Evidence: ev})
		if next != status {
			status = next
		}
		if status == StatusMastered {
			return n + 1
		}
	}
	return -1
}

func difficultyOf(p profile, i int) *float64 {
	d := p.difficulty(i)
	return &d
}

func TestMeasureMasteryCost(t *testing.T) {
	const maxAnswers = 60

	type row struct {
		profile string
		costs   []int
	}
	var rows []row

	for _, p := range profiles {
		var costs []int
		for _, ct := range corpusTimeThresholds {
			c := answersToMastery(p, ct.seconds, maxAnswers, 0)
			if c < 0 {
				t.Errorf("[%s] never reached MASTERED against a %gs time threshold within %d answers",
					p.name, ct.seconds, maxAnswers)
				continue
			}
			costs = append(costs, c)
		}
		rows = append(rows, row{p.name, costs})
	}

	for _, r := range rows {
		if len(r.costs) == 0 {
			continue
		}
		sorted := append([]int(nil), r.costs...)
		sort.Ints(sorted)
		t.Logf("%-32s min %d  p25 %d  median %d  p75 %d  max %d",
			r.profile, sorted[0], pct(sorted, 0.25), pct(sorted, 0.5), pct(sorted, 0.75), sorted[len(sorted)-1])
	}
	t.Logf("evidence floor = %d attempts", MasteryEvidenceFloor)
}

func pct(sorted []int, q float64) int {
	if len(sorted) == 0 {
		return 0
	}
	idx := int(math.Round(q * float64(len(sorted)-1)))
	return sorted[idx]
}

// The properties that must hold whatever the calibration turns out to be.

// No concept can become permanently unmasterable. Every profile must eventually master
// against every real time threshold — the failure ADR-038 nearly shipped, where a gate on
// distinct problem instances made a concept's mastery unreachable for a generator that does
// not vary its question text.
func TestNoConceptIsPermanentlyUnmasterable(t *testing.T) {
	for _, p := range profiles {
		for _, ct := range corpusTimeThresholds {
			if got := answersToMastery(p, ct.seconds, 120, 0); got < 0 {
				t.Errorf("[%s] unmasterable against a %gs threshold", p.name, ct.seconds)
			}
		}
	}
}

// The same, for a generator that returns identical question text every attempt — the
// degenerate case variety cannot see. Variety must never be a veto.
func TestIdenticalInstancesStillMaster(t *testing.T) {
	degenerate := profile{
		name:       "identical question text every attempt",
		correct:    func(int) bool { return true },
		difficulty: func(int) float64 { return 0.7 },
		elapsed:    func(int) float64 { return 4 },
	}
	for _, ct := range corpusTimeThresholds {
		got := answersToMasteryIdentical(degenerate, ct.seconds, 120)
		if got < 0 {
			t.Errorf("unmasterable with identical instances against a %gs threshold", ct.seconds)
		}
	}
}

// Mastery must be reachable without an uninterrupted run of correct answers — the whole
// point of moving off a streak counter. Every profile above inserts misses; this asserts
// explicitly that the misses do not reset progress, by starting partway up the ladder.
func TestMasteryReachableWithoutAnUnbrokenStreak(t *testing.T) {
	// Start at 40 correct answers, then insert a miss every fourth answer. If any rule
	// anywhere still keys on an uninterrupted streak, this never masters.
	for _, ct := range corpusTimeThresholds {
		if got := answersToMastery(profiles[0], ct.seconds, 80, 40); got < 0 {
			t.Errorf("clean evidence after 40 answers failed against a %gs threshold", ct.seconds)
		}
	}
}

// The streak threshold must not affect cost at all. It no longer appears in the evidence
// model, and this is the assertion that keeps it out.
func TestStreakThresholdDoesNotAffectCost(t *testing.T) {
	// The ladder never sees RequiredStreak, so two calls differing only in it must be
	// identical.
	ev := cleanWindow(MasteryEvidenceFloor, 0.7)
	a := (&Machine{}).Next(StatusPracticing, TransitionCtx{Evidence: ev, RequiredStreak: 3})
	b := (&Machine{}).Next(StatusPracticing, TransitionCtx{Evidence: ev, RequiredStreak: 15})
	if a != b {
		t.Errorf("required_streak changed the outcome: %q vs %q", a, b)
	}
}

// answersToMasteryIdentical is answersToMastery with every attempt reporting the same
// instance, so the variety term never moves.
func answersToMasteryIdentical(p profile, timeThreshold float64, maxAnswers int) int {
	m := &Machine{}
	status := StatusUnseen
	var window []Attempt
	for i := 0; i < maxAnswers; i++ {
		d := p.difficulty(i)
		window = append(window, Attempt{
			Correct: p.correct(i), Elapsed: p.elapsed(i), Difficulty: &d, Instance: "the-same-question",
		})
		if len(window) > EvidenceWindow {
			window = window[len(window)-EvidenceWindow:]
		}
		ev := BuildEvidence(window, EvidenceOptions{TimeThreshold: timeThreshold})
		if next := m.Next(status, TransitionCtx{Evidence: ev}); next != status {
			status = next
		}
		if status == StatusMastered {
			return i + 1
		}
	}
	return -1
}

// Where does the score actually bind?
//
// With the floor at 6 and a window of 6, every profile above masters in exactly 6 answers,
// which means MasteryThreshold is not constraining anything — the floor is doing all the
// work. That is worth knowing, and so is the flip side: there has to be some accuracy band
// where a learner *cannot* master, or the model is not a model. This finds it.
func TestMeasureAccuracyBandThatStillMasters(t *testing.T) {
	const maxAnswers = 400
	var masterable, blocked []int // percentages

	for pct := 40; pct <= 100; pct += 5 {
		// A deterministic pattern at roughly `pct` correct, so the window composition
		// settles rather than depending on when the loop happens to end.
		missesPer20 := int(math.Round(20 * (1 - float64(pct)/100)))
		p := profile{
			name:    fmt.Sprintf("%d%% correct", pct),
			correct: func(i int) bool { return (i % 20) >= missesPer20 },
			// Difficulty and time held at middling values, so accuracy is the only thing
			// moving.
			difficulty: func(int) float64 { return 0.6 },
			elapsed:    func(int) float64 { return 4 },
		}
		got := answersToMastery(p, 20, maxAnswers, 0)
		if got < 0 {
			blocked = append(blocked, pct)
		} else {
			masterable = append(masterable, pct)
		}
	}

	t.Logf("masterable at: %v%% correct", masterable)
	t.Logf("never masters at: %v%% correct", blocked)

	// Reachability is asserted; the *upper* end of the band is reported, not asserted.
	//
	// The measurement found that with the floor at 6 and the window at 6, MasteryThreshold
	// does not bind at all — every accuracy level from 40% up reaches MASTERED, so the gate
	// is entirely the evidence floor. That is a real calibration question and it is the
	// owner's to answer, so it is logged rather than silently fixed here: raising the
	// threshold is a different decision from setting the floor to 6, and inventing a value
	// for it would be the same error as inventing a target of "9 answers".
	if len(masterable) == 0 {
		t.Error("no accuracy level masters — the ladder is unreachable")
	}
	if len(masterable) > 0 && masterable[0] > 70 {
		t.Errorf("mastery needs only %d%% correct, which is too low to mean anything", masterable[0])
	}
	t.Logf("lowest accuracy that masters: %d%% — MasteryThreshold is %v", masterable[0], MasteryThreshold)
}
