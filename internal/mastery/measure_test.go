package mastery

import (
	"fmt"
	"math"
	"sort"
	"testing"
)

// What mastery costs, and what the threshold actually does — measured rather than assumed.
//
// Two measurement lessons are baked in here, both of them learned the hard way in this
// file's own first version:
//
//  1. A probe's *period* is not the model. The first accuracy sweep used `i%20 >= 12` —
//     twelve wrong then eight right — and reported that every accuracy level from 40% up
//     reached mastery. It hadn't; a six-attempt window simply eventually landed entirely
//     inside the clean run. Learners do not answer in blocks. Every pattern below is evenly
//     distributed instead (see `correctAt`), so each window sees the nominal accuracy.
//  2. Measure the cost; do not assume it. `MasteryEvidenceFloor` was raised from 3 to 6 on
//     the strength of this file, and the honest distribution turned out to be exactly 6 —
//     which is a *finding*, not a confirmation.
//
// The distribution is reported, not asserted. The cost of mastery is a calibration decision
// and belongs to the owner. What is asserted are the properties that must hold whatever the
// cost is, because those are not calibration choices.

// correctAt reports whether attempt i is correct for a learner who gets k of every n right,
// spread as evenly as possible so that no window of EvidenceWindow is unrepresentative.
//
// `(i*k) % n < k` places exactly k correct answers in each period of n, evenly spaced. A
// learner who is right 60% of the time is wrong 40% of the time *throughout*, not in
// bursts — and the difference between those two claims is the entire measurement.
func correctAt(i, n, k int) bool {
	if k <= 0 {
		return false
	}
	if k >= n {
		return true
	}
	return (i*k)%n < k
}

type profile struct {
	name string
	// n and k describe accuracy: correct on k of every n attempts, evenly spread.
	n, k int
	// difficulty is the difficulty served on attempt i.
	difficulty func(i int) float64
	// elapsed is how long the learner took, in seconds.
	elapsed func(i int) float64
	// unknownDifficulty leaves the attempt without a difficulty, modelling a host that
	// fails to report one.
	unknownDifficulty bool
}

const midDifficulty = 0.6

func mid(int) float64 { return midDifficulty }

var profiles = []profile{
	{
		name: "consistently correct", n: 5, k: 5,
		difficulty: mid, elapsed: func(int) float64 { return 4 },
	},
	{
		name: "approximately 80% correct", n: 5, k: 4,
		difficulty: mid, elapsed: func(int) float64 { return 4 },
	},
	{
		name: "approximately 60% correct", n: 5, k: 3,
		difficulty: mid, elapsed: func(int) float64 { return 4 },
	},
	{
		name: "alternating correct/incorrect", n: 2, k: 1,
		difficulty: mid, elapsed: func(int) float64 { return 4 },
	},
	{
		name: "approximately 40% correct", n: 5, k: 2,
		difficulty: mid, elapsed: func(int) float64 { return 4 },
	},
	{
		name: "slow but accurate", n: 5, k: 5,
		difficulty: mid, elapsed: func(int) float64 { return 120 },
	},
	{
		name: "fast but inaccurate", n: 5, k: 2,
		difficulty: mid, elapsed: func(int) float64 { return 2 },
	},
	{
		name: "correct on the hardest questions", n: 5, k: 5,
		difficulty: func(int) float64 { return 1.0 }, elapsed: func(int) float64 { return 4 },
	},
	{
		name: "correct on the easiest questions", n: 5, k: 5,
		difficulty: func(int) float64 { return 0.3 }, elapsed: func(int) float64 { return 4 },
	},
	{
		name: "starts rough, then steady", n: 5, k: 4,
		// Wrong for the first three attempts, then the steady 80% profile.
		difficulty: mid, elapsed: func(int) float64 { return 4 },
	},
}

// corpusTimeThresholds are the real `avg_time_seconds` values in the corpus (with the number
// of concepts using each). They matter because the time term is scored against the
// concept's own threshold, so "slow" is relative: 120s is slow for a 15s concept and fast
// for a 150s one.
var corpusTimeThresholds = []struct {
	seconds  float64
	concepts int
}{
	{10, 32}, {15, 61}, {18, 60}, {20, 54}, {25, 37}, {30, 18}, {60, 13}, {120, 162},
}

// run drives one profile against one concept and reports how many answers it took to reach
// MASTERED, or -1. Instance text is unique per attempt so the variety term behaves the way
// it does in production, where every question is a different instance.
func run(p profile, timeThreshold float64, maxAnswers, startAt int) int {
	return runWith(p, timeThreshold, maxAnswers, startAt, func(i int) string {
		return fmt.Sprintf("q%d", i)
	})
}

func runWith(p profile, timeThreshold float64, maxAnswers, startAt int, instance func(int) string) int {
	m := &Machine{}
	status := StatusUnseen
	var window []Attempt
	for i := 0; i < maxAnswers; i++ {
		n := startAt + i
		correct := correctAt(n, p.n, p.k)
		if p.name == "starts rough, then steady" && n < 3 {
			correct = false
		}
		var d *float64
		if !p.unknownDifficulty {
			v := p.difficulty(n)
			d = &v
		}
		window = append(window, Attempt{
			Correct:    correct,
			Elapsed:    p.elapsed(n),
			Difficulty: d,
			Instance:   instance(n),
		})
		if len(window) > EvidenceWindow {
			window = window[len(window)-EvidenceWindow:]
		}
		ev := BuildEvidence(window, EvidenceOptions{TimeThreshold: timeThreshold})
		if next := m.Next(status, TransitionCtx{Evidence: ev}); next != status {
			status = next
		}
		if status == StatusMastered {
			return n + 1
		}
	}
	return -1
}

func sameInstance(i int) string { return "the-same-question" }

func pctile(sorted []int, q float64) int {
	if len(sorted) == 0 {
		return 0
	}
	return sorted[int(math.Round(q*float64(len(sorted)-1)))]
}

// TestMeasureMasteryCost is the headline number: how many answers to MASTERED, per profile,
// across every real concept time threshold. Reported, not asserted.
func TestMeasureMasteryCost(t *testing.T) {
	const maxAnswers = 60
	t.Logf("evidence floor = %d, threshold = %.2f, window = %d",
		MasteryEvidenceFloor, MasteryThreshold, EvidenceWindow)
	t.Logf("%-32s %-6s %-5s %-7s %-5s %-5s", "profile", "min", "p25", "median", "p75", "max")

	for _, p := range profiles {
		var costs []int
		for _, ct := range corpusTimeThresholds {
			c := run(p, ct.seconds, maxAnswers, 0)
			if c < 0 {
				continue
			}
			costs = append(costs, c)
		}
		if len(costs) == 0 {
			t.Logf("%-32s %s", p.name, "never masters at any concept threshold")
			continue
		}
		sort.Ints(costs)
		t.Logf("%-32s %-6d %-5d %-7d %-5d %-5d", p.name,
			costs[0], pctile(costs, 0.25), pctile(costs, 0.5), pctile(costs, 0.75), costs[len(costs)-1])
	}
}

// TestMeasureAccuracyBandThatStillMasters finds where MasteryThreshold actually bites.
//
// With the floor and the window both at 6, the floor is reached before the score has much to
// say — but the score is not inert. This sweeps evenly-distributed accuracy and reports the
// band, then asserts the two properties that make the threshold meaningful rather than
// decorative: the ladder is reachable, and it is not a rubber stamp.
func TestMeasureAccuracyBandThatStillMasters(t *testing.T) {
	const maxAnswers = 200
	// Difficulty and time held at middling values so accuracy is the only thing moving.
	base := profile{n: 5, k: 5, difficulty: mid, elapsed: func(int) float64 { return 4 }}

	var masterable, blocked []int
	for k := 0; k <= 5; k++ {
		p := base
		p.k = k
		pct := k * 20
		ok := false
		for _, ct := range corpusTimeThresholds {
			if run(p, ct.seconds, maxAnswers, 0) > 0 {
				ok = true
				break
			}
		}
		if ok {
			masterable = append(masterable, pct)
		} else {
			blocked = append(blocked, pct)
		}
	}

	t.Logf("masterable at: %v%% correct", masterable)
	t.Logf("never masters at: %v%% correct", blocked)

	if len(masterable) == 0 {
		t.Fatal("no accuracy level masters — the ladder is unreachable")
	}
	if len(blocked) == 0 {
		t.Error("every accuracy level masters — MasteryThreshold is decorative")
	}
	// Mastery is the exit from the teaching loop, so the bar has to mean something: better
	// than a coin flip, and no lower.
	if lowest := masterable[0]; lowest < 50 {
		t.Errorf("mastery is reachable at %d%% correct, which is at or below chance", lowest)
	}
	t.Logf("lowest accuracy that masters: %d%% (threshold %.2f)", masterable[0], MasteryThreshold)
}

// --- properties, not calibration ---------------------------------------------------

// No concept can become permanently unmasterable. The failure ADR-038 nearly shipped was a
// gate on distinct problem instances, which made mastery unreachable for a generator whose
// question text does not vary.
func TestNoConceptIsPermanentlyUnmasterable(t *testing.T) {
	for _, p := range profiles {
		if !reliablyCorrect(p) {
			continue // meant to fail; see TestMeasureAccuracyBandThatStillMasters
		}
		for _, ct := range corpusTimeThresholds {
			if run(p, ct.seconds, 200, 0) < 0 {
				t.Errorf("[%s] unmasterable against a %gs time threshold", p.name, ct.seconds)
			}
		}
	}
}

// A generator that returns identical question text every attempt must not make mastery
// unreachable for a learner who is reliably getting it right.
//
// This is the failure ADR-038 nearly shipped, when `EnoughForMastery` also required
// `Instances >= 2` as a hard gate: a generator that does not vary its question text holds
// Instances at 1 forever, so the concept can never be mastered however much the learner
// demonstrates it. The property under test is reachability, not equality — see
// TestVarietyRaisesTheBarWithoutBlockingMastery for the margin, which is deliberate.
func TestIdenticalInstancesStillMaster(t *testing.T) {
	for _, p := range profiles {
		if !reliablyCorrect(p) {
			continue
		}
		for _, ct := range corpusTimeThresholds {
			if runWith(p, ct.seconds, 200, 0, sameInstance) < 0 {
				t.Errorf("[%s] unmasterable with identical instances against a %gs threshold",
					p.name, ct.seconds)
			}
		}
	}
}

// reliablyCorrect reports whether a profile is accurate enough that blocked-by-variety would
// be a bug rather than a judgement.
func reliablyCorrect(p profile) bool {
	switch p.name {
	case "fast but inaccurate", "approximately 40% correct", "alternating correct/incorrect",
		"approximately 60% correct":
		return false
	}
	return true
}

// Variety is a weight, not a veto — and this pins down exactly where the line is, because
// "cannot veto" is a claim that is easy to make and easy to get wrong.
//
// Measured across the real corpus time thresholds:
//
//	 40% correct: never masters, distinct or identical
//	 60% correct: masters on distinct instances, never on identical ones
//	 80% correct: masters either way
//	100% correct: masters either way
//
// So the 0.15 variety bonus is large enough to flip the decision at 60%. That is correct and
// should not be "fixed": three right answers out of five spread across three different
// problems is real evidence of transfer, and three right answers out of five on the *same*
// problem three times is evidence about one problem. The property that matters — and the one
// ADR-038's `Instances >= 2` gate broke — is that no amount of demonstrated competence
// becomes *unreachable*. A learner at 80% masters on a generator that never varies.
func TestVarietyRaisesTheBarWithoutBlockingMastery(t *testing.T) {
	base := func(k int) profile {
		return profile{n: 5, k: k, difficulty: mid, elapsed: func(int) float64 { return 4 }}
	}
	mastersSomewhere := func(p profile, identical bool) bool {
		for _, ct := range corpusTimeThresholds {
			var got int
			if identical {
				got = runWith(p, ct.seconds, 60, 0, sameInstance)
			} else {
				got = run(p, ct.seconds, 60, 0)
			}
			if got > 0 {
				return true
			}
		}
		return false
	}

	if !mastersSomewhere(base(4), true) {
		t.Error("a learner who is right 80% of the time cannot master on a generator " +
			"that never varies — " +
			"variety has become a veto")
	}
	if !mastersSomewhere(base(5), true) {
		t.Error("a fully correct learner cannot master on a generator that never varies")
	}
	if mastersSomewhere(base(3), true) {
		t.Log("note: a learner who is right 60% of the time now masters on identical " +
			"instances too — the " +
			"variety bonus margin has moved")
	}
}

// Mastery must be reachable without an uninterrupted run of correct answers — the whole
// point of moving off a streak counter. Starts 40 answers in, then misses every fourth.
func TestMasteryReachableWithoutAnUnbrokenStreak(t *testing.T) {
	p := profile{
		name: "long history then misses", n: 5, k: 4,
		difficulty: mid, elapsed: func(int) float64 { return 4 },
	}
	for _, ct := range corpusTimeThresholds {
		if got := run(p, ct.seconds, 80, 40); got < 0 {
			t.Errorf("an established learner who now misses one in five failed to master "+
				"against a %gs threshold — something still keys on an unbroken streak", ct.seconds)
		}
	}
}

// The per-concept streak threshold must not affect cost. It no longer appears in the
// evidence model, and this is the assertion that keeps it out.
func TestStreakThresholdDoesNotAffectCost(t *testing.T) {
	d := midDifficulty
	ev := BuildEvidence([]Attempt{
		{Correct: true, Elapsed: 4, Difficulty: &d, Instance: "a"},
		{Correct: true, Elapsed: 4, Difficulty: &d, Instance: "b"},
		{Correct: true, Elapsed: 4, Difficulty: &d, Instance: "c"},
		{Correct: true, Elapsed: 4, Difficulty: &d, Instance: "d"},
		{Correct: true, Elapsed: 4, Difficulty: &d, Instance: "e"},
		{Correct: true, Elapsed: 4, Difficulty: &d, Instance: "f"},
	}, EvidenceOptions{TimeThreshold: 20})

	a := (&Machine{}).Next(StatusPracticing, TransitionCtx{Evidence: ev, RequiredStreak: 3})
	b := (&Machine{}).Next(StatusPracticing, TransitionCtx{Evidence: ev, RequiredStreak: 15})
	if a != b {
		t.Errorf("required_streak changed the outcome: %q vs %q", a, b)
	}
}

// --- issue 2: a host must not be able to change the answer ----------------------------

// Unknown difficulty must not be *better* than known difficulty.
//
// This is the regression test for the asymmetry between hosts. Unknown difficulty used to be
// scored as difficultyTerm = 1.0 — the top of the range — so a missing value read as "the
// hardest question we could have asked". At 50% accuracy that was the difference between
// scoring 0.65 and clearing the bar, and 0.57 and failing it: the same learner, the same
// answers, mastering on one surface and not on another.
func TestUnknownDifficultyIsNotMoreGenerousThanKnown(t *testing.T) {
	for k := 1; k <= 5; k++ {
		known := profile{n: 5, k: k, difficulty: mid, elapsed: func(int) float64 { return 4 }}
		unknown := known
		unknown.unknownDifficulty = true

		for _, ct := range corpusTimeThresholds {
			a := run(known, ct.seconds, 60, 0)
			b := run(unknown, ct.seconds, 60, 0)
			if b > 0 && a < 0 {
				t.Errorf("at %d%% correct against a %gs threshold: unknown difficulty reaches "+
					"MASTERED in %d answers but the same answers at a known difficulty never do",
					k*20, ct.seconds, b)
			}
		}
	}
}

// And the two must agree exactly when the unknown attempts are genuinely average, which is
// what NeutralDifficulty claims.
func TestUnknownDifficultyEqualsNeutralDifficulty(t *testing.T) {
	unknown := unknownDifficultyWindow(EvidenceWindow)
	neutral := cleanWindow(EvidenceWindow, NeutralDifficulty)
	if math.Abs(unknown.Score-neutral.Score) > 1e-9 {
		t.Errorf("unknown difficulty scored %.4f but difficulty %.2f scored %.4f; they should "+
			"be the same claim", unknown.Score, NeutralDifficulty, neutral.Score)
	}
}

// --- issue 4: mastery is attainment, retention is a separate question ------------------

// Mastery never un-learns. A concept that has been mastered stays mastered no matter how bad
// the subsequent evidence gets; recovery is a question about retention and belongs to decay
// and /review, not to this ladder.
func TestMasteryIsNeverLostToBadEvidence(t *testing.T) {
	m := &Machine{}
	status := StatusMastered
	for i := 0; i < 30; i++ {
		window := []Attempt{
			{Correct: false, Elapsed: 300, Instance: fmt.Sprintf("bad%d", i)},
			{Correct: false, Elapsed: 300, Instance: fmt.Sprintf("bad%d-b", i)},
			{Correct: false, Elapsed: 300, Instance: fmt.Sprintf("bad%d-c", i)},
			{Correct: false, Elapsed: 300, Instance: fmt.Sprintf("bad%d-d", i)},
			{Correct: false, Elapsed: 300, Instance: fmt.Sprintf("bad%d-e", i)},
			{Correct: false, Elapsed: 300, Instance: fmt.Sprintf("bad%d-f", i)},
		}
		ev := BuildEvidence(window, EvidenceOptions{TimeThreshold: 5})
		if ev.Score > 0.01 {
			t.Fatalf("test setup: a window of six misses scored %.3f", ev.Score)
		}
		if next := m.Next(status, TransitionCtx{Evidence: ev}); next != StatusMastered {
			t.Fatalf("mastery was lost to bad evidence: %v -> %v", status, next)
		}
	}
}

// The evidence model contains no spacing term, and that is deliberate: it scores attainment
// from attempts and knows nothing about elapsed time since the last one. Spaced retention is
// the deferred `tier_attained_at` work, and evidence is where it will *not* go — a learner
// who mastered a concept a year ago and has not touched it since has demonstrated exactly
// what they demonstrated then.
func TestEvidenceHasNoSpacingTerm(t *testing.T) {
	// BuildEvidence takes only attempts and the concept's thresholds — there is no clock and
	// no "when was this last reviewed" input, so no term can depend on elapsed time. That is
	// what keeps the deferred spaced work additive rather than a rewrite of this function.
	//
	// Two identical windows must therefore score identically, and a window of old-looking
	// attempts must not decay: a learner who mastered a concept a year ago and has not
	// touched it since demonstrated exactly what they demonstrated then.
	old1 := cleanWindow(EvidenceWindow, midDifficulty)
	old2 := cleanWindow(EvidenceWindow, midDifficulty)
	if old1.Score != old2.Score {
		t.Errorf("identical evidence scored differently: %.4f vs %.4f", old1.Score, old2.Score)
	}
	if fresh := cleanWindow(EvidenceWindow, midDifficulty); fresh.Score != old1.Score {
		t.Errorf("evidence scored %.4f once and %.4f another time — something is not a pure "+
			"function of the attempts", old1.Score, fresh.Score)
	}
}

// --- the band a learner reads -------------------------------------------------------

// The band is presentational, but it must not lie about what the evidence supports.
//
// The client version of this label offered "well retained" off a score with no spacing
// term at all — a retention claim drawn from evidence that knows nothing about retention.
// And with no floor it read "strong" after a *single* correct answer, because a lone answer
// carries the most recency weight there is. Both defects arrived through the label rather
// than the state, so the label needs its own tests.
func TestEvidenceBandClaimsNothingBeyondTheEvidence(t *testing.T) {
	// One correct answer scores 0.84 — above every score threshold — and must still read as
	// early, because a single attempt is not evidence.
	single := BuildEvidence([]Attempt{{
		Correct: true, Elapsed: 4, Instance: "a",
	}}, EvidenceOptions{TimeThreshold: 20})
	if single.Score < 0.8 {
		t.Fatalf("test setup: a single correct answer scored %.3f, expected it to clear the "+
			"score thresholds so the floor is what stops it", single.Score)
	}
	if got := EvidenceBand(single); got != "early" {
		t.Errorf("one correct answer reads %q, want %q — the band must be floored by evidence "+
			"the way mastery is", got, "early")
	}

	// No band may promise anything about remembering it later.
	// No clock is involved anywhere, so nothing here can distinguish an answer from
	// yesterday and one from last year. See TestEvidenceHasNoSpacingTerm.
	for _, band := range []Evidence{
		single,
		BuildEvidence(cleanAttempts(6, true), EvidenceOptions{TimeThreshold: 20}),
		BuildEvidence(cleanAttempts(6, false), EvidenceOptions{TimeThreshold: 20}),
	} {
		if got := EvidenceBand(band); got == "mastered" || got == "well retained" || got == "retained" {
			t.Errorf("band %q claims mastery or retention, which BuildEvidence cannot support", got)
		}
	}
}

// The band is a function of the evidence and nothing else, so it cannot drift from the
// ladder's own inputs — and it must move when the evidence moves.
func TestEvidenceBandTracksTheEvidence(t *testing.T) {
	strong := BuildEvidence(cleanAttempts(6, true), EvidenceOptions{TimeThreshold: 20})
	weak := BuildEvidence(cleanAttempts(6, false), EvidenceOptions{TimeThreshold: 20})
	if EvidenceBand(strong) == EvidenceBand(weak) {
		t.Errorf("six correct and six wrong both read %q", EvidenceBand(strong))
	}
	if EvidenceBand(strong) != EvidenceBand(BuildEvidence(cleanAttempts(6, true), EvidenceOptions{TimeThreshold: 20})) {
		t.Error("identical evidence produced different bands")
	}
}

func cleanAttempts(n int, correct bool) []Attempt {
	out := make([]Attempt, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, Attempt{
			Correct: correct, Elapsed: 4, Difficulty: floatPtrLocal(midDifficulty),
			Instance: fmt.Sprintf("q%d", i),
		})
	}
	return out
}

func floatPtrLocal(v float64) *float64 { return &v }
