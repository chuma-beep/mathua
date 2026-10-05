package mastery

import (
	"math"
	"testing"
)

// Next() transitions
//
// These used to be written against `streak >= required AND avg time <= threshold`, and
// they asserted the old rule faithfully: five consecutive correct answers to leave UNSEEN,
// and the same five again for each tier after it. They are rewritten here against the
// evidence ladder because the rule they pinned is the one being replaced — the ladder is
// now three different questions, and a test that still demanded a full threshold-length
// streak would be asserting the bug.

func evidenceOf(correct bool, instance string, difficulty, elapsed float64) Attempt {
	return Attempt{Correct: correct, Instance: instance, Difficulty: diff(difficulty), Elapsed: elapsed}
}

// cleanWindow builds `n` correct attempts on distinct instances at `difficulty` — what a
// learner who is reliably getting it right produces. Used where a test needs evidence that
// clears the mastery floor rather than evidence of a specific shape.
func cleanWindow(n int, difficulty float64) Evidence {
	var as []Attempt
	for i := 0; i < n; i++ {
		as = append(as, evidenceOf(true, string(rune('a'+i)), difficulty, 4))
	}
	return BuildEvidence(as, EvidenceOptions{TimeThreshold: 10})
}

// unknownDifficultyWindow is cleanWindow with no difficulty recorded at all, which is a
// different fact from a recorded 0 — see attempts.difficulty being nullable, and
// difficultyOrNil refusing to write zero.
func unknownDifficultyWindow(n int) Evidence {
	var as []Attempt
	for i := 0; i < n; i++ {
		as = append(as, Attempt{Correct: true, Elapsed: 4, Instance: string(rune('a' + i))})
	}
	return BuildEvidence(as, EvidenceOptions{TimeThreshold: 10})
}

// started is one correct answer: enough to say the learner has begun.
func started() Evidence {
	return BuildEvidence([]Attempt{evidenceOf(true, "a", 0.5, 4)}, EvidenceOptions{TimeThreshold: 10})
}

// consistent is two correct on two instances: reliable, not demonstrated.
func consistent() Evidence {
	return BuildEvidence([]Attempt{
		evidenceOf(true, "a", 0.4, 4), evidenceOf(true, "b", 0.4, 4),
	}, EvidenceOptions{TimeThreshold: 10})
}

// demonstrated is four correct across four instances at real difficulty.
func demonstrated() Evidence {
	return cleanWindow(MasteryEvidenceFloor, 0.8)
}

func TestNext_UnseenToLearning(t *testing.T) {
	m := &Machine{}
	got := m.Next(StatusUnseen, TransitionCtx{Evidence: started()})
	if got != StatusLearning {
		t.Errorf("UNSEEN → LEARNING: expected %q, got %q", StatusLearning, got)
	}
}

// The substantive change. One correct answer used to be nowhere near enough: the rule
// demanded the full `required_streak`, which is 10 for 459 of the 657 concepts.
func TestNext_OneCorrectAnswerLeavesUnseen(t *testing.T) {
	for _, required := range []int{3, 10, 15} {
		m := &Machine{}
		ctx := TransitionCtx{Evidence: started(), RequiredStreak: required}
		if got := m.Next(StatusUnseen, ctx); got != StatusLearning {
			t.Errorf("required_streak %d: one correct answer gave %q, want LEARNING", required, got)
		}
	}
}

func TestNext_UnseenStays_NothingAnswered(t *testing.T) {
	m := &Machine{}
	got := m.Next(StatusUnseen, TransitionCtx{Evidence: BuildEvidence(nil, EvidenceOptions{})})
	if got != StatusUnseen {
		t.Errorf("expected UNSEEN, got %q", got)
	}
}

// A wrong first answer must not move the concept, however fast or however many.
func TestNext_UnseenStays_WrongAnswer(t *testing.T) {
	m := &Machine{}
	ev := BuildEvidence([]Attempt{evidenceOf(false, "a", 0.5, 1)}, EvidenceOptions{TimeThreshold: 10})
	if got := m.Next(StatusUnseen, TransitionCtx{Evidence: ev}); got != StatusUnseen {
		t.Errorf("expected UNSEEN after a wrong answer, got %q", got)
	}
}

func TestNext_LearningToPracticing(t *testing.T) {
	m := &Machine{}
	if got := m.Next(StatusLearning, TransitionCtx{Evidence: consistent()}); got != StatusPracticing {
		t.Errorf("LEARNING → PRACTICING: expected PRACTICING, got %q", got)
	}
}

// Variety is rewarded, not required.
//
// Getting the same question right twice *is* reliability, so it reaches PRACTICING. It is
// weaker evidence than the same two answers on two different questions, and the score says
// so — but it must not be a veto. An earlier version gated mastery on `Instances >= 2`, and
// the engine tests immediately hit the failure that caused: a generator that returns the
// same question text every attempt drives Instances to 1 forever and the concept can never
// be mastered at all.
func TestNext_VarietyIsRewardedNotRequired(t *testing.T) {
	// The degenerate generator: a full window of correct answers whose question text never
	// varies, so every attempt is the same instance.
	var degenerate []Attempt
	for i := 0; i < MasteryEvidenceFloor; i++ {
		degenerate = append(degenerate, evidenceOf(true, "a", 0.9, 3))
	}
	same := BuildEvidence(degenerate, EvidenceOptions{TimeThreshold: 10})
	varied := cleanWindow(MasteryEvidenceFloor, 0.9)

	if same.Instances != 1 {
		t.Fatalf("the degenerate generator produced %d instances, want 1", same.Instances)
	}
	if varied.Instances != MasteryEvidenceFloor {
		t.Fatalf("varied window produced %d instances, want %d", varied.Instances, MasteryEvidenceFloor)
	}
	if varied.Score <= same.Score {
		t.Errorf("variety did not raise the score: %v vs %v", varied.Score, same.Score)
	}

	// Both reach PRACTICING; both reach MASTERED too, because variety is not a gate.
	m := &Machine{}
	if got := m.Next(StatusLearning, TransitionCtx{Evidence: same}); got != StatusPracticing {
		t.Errorf("same-instance consistency gave %q, want PRACTICING", got)
	}
	// From UNSEEN, three correct on one instance is enough to walk the whole ladder — this
	// is the regression guard for the unreachable-mastery failure.
	if got := m.Next(StatusUnseen, TransitionCtx{Evidence: same}); got != StatusLearning {
		t.Errorf("UNSEEN with three correct gave %q, want LEARNING", got)
	}
	if got := m.Next(StatusLearning, TransitionCtx{Evidence: same}); got != StatusPracticing {
		t.Errorf("LEARNING with three correct gave %q, want PRACTICING", got)
	}
	if got := m.Next(StatusPracticing, TransitionCtx{Evidence: same}); got != StatusMastered {
		t.Errorf("a non-varied generator must not make mastery unreachable; got %q", got)
	}
}

func TestNext_PracticingToMastered(t *testing.T) {
	m := &Machine{}
	if got := m.Next(StatusPracticing, TransitionCtx{Evidence: demonstrated()}); got != StatusMastered {
		t.Errorf("PRACTICING → MASTERED: expected MASTERED, got %q", got)
	}
}

// One rung per transition. Strong evidence on a LEARNING concept advances it one step, not
// straight to MASTERED — otherwise a learner with a burst of history skips the middle.
func TestNext_AdvancesOneRungOnly(t *testing.T) {
	m := &Machine{}
	if got := m.Next(StatusLearning, TransitionCtx{Evidence: demonstrated()}); got != StatusPracticing {
		t.Errorf("LEARNING with strong evidence gave %q, want PRACTICING", got)
	}
	if got := m.Next(StatusUnseen, TransitionCtx{Evidence: demonstrated()}); got != StatusLearning {
		t.Errorf("UNSEEN with strong evidence gave %q, want LEARNING", got)
	}
}

// Mastery is an exit. A later miss lowers the evidence score but must not un-learn it —
// ADR-037 made the exit unconditional, and decay is a read-time question about retention.
func TestNext_MasteredIsSticky(t *testing.T) {
	m := &Machine{}
	bad := BuildEvidence([]Attempt{
		evidenceOf(false, "z", 0.9, 1), evidenceOf(false, "y", 0.9, 1),
		evidenceOf(false, "x", 0.9, 1),
	}, EvidenceOptions{TimeThreshold: 10})
	if got := m.Next(StatusMastered, TransitionCtx{Evidence: bad}); got != StatusMastered {
		t.Errorf("MASTERED with terrible evidence gave %q, want MASTERED", got)
	}
}

// The streak fields are context, not authority. A long streak on evidence that does not
// reach the next rung must not advance the ladder.
func TestNext_StreakAloneDoesNotAdvance(t *testing.T) {
	m := &Machine{}
	ctx := TransitionCtx{
		Evidence:          BuildEvidence([]Attempt{evidenceOf(true, "a", 0.3, 2)}, EvidenceOptions{TimeThreshold: 10}),
		Streak:            99,
		RequiredStreak:    3,
		AvgResponseTime:   1,
		ResponseThreshold: 10,
	}
	if got := m.Next(StatusPracticing, ctx); got == StatusMastered {
		t.Errorf("one attempt with streak 99 reached MASTERED, giving %q", got)
	}
	if got := m.Next(StatusPracticing, TransitionCtx{
		Evidence: BuildEvidence(nil, EvidenceOptions{}),
		Streak:   99, RequiredStreak: 3, AvgResponseTime: 1, ResponseThreshold: 10,
	}); got != StatusPracticing {
		t.Errorf("a full streak with no evidence gave %q, want PRACTICING", got)
	}
}

// Streak is evidence, not authority — these two replace the old "excess streak still
// advances" and "streak == required advances", which asserted the rule being replaced.
func TestNext_ExcessStreakDoesNotSubstituteForEvidence(t *testing.T) {
	m := &Machine{}
	// Twelve in a row, comfortably over a required_streak of 5, on evidence that has only
	// ever been one question. Consistency is real; it is not the rung asked for.
	ctx := TransitionCtx{
		Evidence:          BuildEvidence([]Attempt{evidenceOf(true, "a", 0.5, 2)}, EvidenceOptions{TimeThreshold: 10}),
		Streak:            12,
		RequiredStreak:    5,
		AvgResponseTime:   2,
		ResponseThreshold: 10,
	}
	if got := m.Next(StatusPracticing, ctx); got == StatusMastered {
		t.Error("a long streak on a single question instance established mastery")
	}
}

func TestNext_ExactStreakBoundaryIsIrrelevantNow(t *testing.T) {
	m := &Machine{}
	// Exactly at both thresholds, with no evidence recorded — the old rule advanced here.
	ctx := TransitionCtx{
		Streak:            5,
		RequiredStreak:    5,
		AvgResponseTime:   10.0,
		ResponseThreshold: 10.0,
	}
	if got := m.Next(StatusLearning, ctx); got == StatusPracticing {
		t.Error("streak == required advanced the ladder with no evidence at all")
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

// TestMasteryPct_NeverGoesBackwards walks a concept from UNSEEN to MASTERED through the
// real ladder, building evidence attempt by attempt exactly as the engine does, and
// asserts the bar never decreases.
//
// This is the assertion that would have caught the sawtooth. The old formula was
// `streak / required_streak`, so at the instant of each advance the bar fell from 100% to
// ~10%, three times on the way up: the learner watched their progress drop while doing
// everything right.
func TestMasteryPct_NeverGoesBackwards(t *testing.T) {
	const req = 10
	m := &Machine{}
	status := StatusUnseen
	streak := 0
	prev := MasteryPct(status, streak, req)
	advances := 0
	answered := 0

	// Difficulty climbs as the learner does better and every answer is a new instance —
	// the shape the generator produces and the evidence model is built to read.
	difficulties := []float64{0.35, 0.45, 0.55, 0.7, 0.8, 0.9}

	for answer := 0; answer < 20 && status != StatusMastered; answer++ {
		answered++
		streak++ // a correct answer

		// A six-attempt window, rebuilt each turn exactly as the engine rebuilds it.
		window := make([]Attempt, 0, 6)
		for k := 0; k < 6; k++ {
			i := answer - k
			if i < 0 {
				i = 0
			}
			window = append(window, Attempt{
				Correct: true, Elapsed: 3,
				Difficulty: diff(difficulties[i%len(difficulties)]),
				Instance:   string(rune('a' + i)),
			})
		}
		ev := BuildEvidence(window, EvidenceOptions{TimeThreshold: 10})

		next := m.Next(status, TransitionCtx{
			Evidence: ev, Streak: streak, RequiredStreak: req,
			AvgResponseTime: 3, ResponseThreshold: 10,
		})
		if next != status {
			status = next
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
		t.Fatalf("never reached MASTERED, got %s after %d correct answers", status, answered)
	}
	if advances != 3 {
		t.Errorf("expected 3 tier advances, got %d", advances)
	}
	if prev != 1 {
		t.Errorf("final bar = %v, want 1", prev)
	}

	// The cost, measured rather than asserted. This is the number the ladder exists to
	// move: it used to be 3 × required_streak consecutive correct answers, 30 for the 459
	// concepts whose threshold is 10.
	t.Logf("%d concepts reached MASTERED in %d correct answers", len(difficulties), answered)
	if answered > 12 {
		t.Errorf("mastery took %d answers, which is back toward the old cost", answered)
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

func diff(d float64) *float64 { return &d }

func at(vals ...float64) []*float64 {
	out := make([]*float64, len(vals))
	for i, v := range vals {
		out[i] = &v
	}
	return out
}

// The evidence score must agree with the browser model it replaces, or the label in /learn
// and the decision that gates progression will drift apart again — which is the state
// ADR-037 spent a week untangling.
func TestBuildEvidence_MatchesBrowserModelShape(t *testing.T) {
	// 3 correct at difficulty 0.7, on 3 distinct instances, fast.
	ev := cleanWindow(MasteryEvidenceFloor, 0.7)

	if ev.Attempts != MasteryEvidenceFloor || ev.Correct != MasteryEvidenceFloor || ev.Instances != MasteryEvidenceFloor {
		t.Fatalf("assembled %+v", ev)
	}
	// accuracy 1.0 × (0.6 + 0.4×0.7) = 0.88, + variety 0.15 (six instances, capped),
	// × time 1.0 → 1.03, clamped to 1.
	want := 1.0
	if math.Abs(ev.Score-want) > 1e-9 {
		t.Errorf("score = %v, want %v", ev.Score, want)
	}
	if !EnoughForMastery(ev) {
		t.Errorf("three distinct correct instances at difficulty 0.7 should establish mastery, got %+v", ev)
	}
}

func TestBuildEvidence_OneCorrectAnswerCannotMaster(t *testing.T) {
	// A single correct answer is weighted 1.0 and would clear the score on its own. This
	// is the guard that stops one lucky question becoming a state change.
	ev := BuildEvidence([]Attempt{
		{Correct: true, Elapsed: 2, Difficulty: diff(1.0), Instance: "a"},
	}, EvidenceOptions{TimeThreshold: 10})
	if ev.Score < MasteryThreshold {
		t.Fatalf("a lone perfect answer should still score well, got %v", ev.Score)
	}
	if EnoughForMastery(ev) {
		t.Errorf("one attempt established mastery: %+v", ev)
	}
	// The same question three times is recorded as one instance and is *scored* lower for
	// it — but it is not blocked, because a generator that never varies its question text
	// would otherwise make the concept permanently unmasterable.
	same := BuildEvidence([]Attempt{
		{Correct: true, Elapsed: 2, Difficulty: diff(0.5), Instance: "a"},
		{Correct: true, Elapsed: 2, Difficulty: diff(0.5), Instance: "a"},
		{Correct: true, Elapsed: 2, Difficulty: diff(0.5), Instance: "a"},
	}, EvidenceOptions{TimeThreshold: 10})
	if same.Instances != 1 {
		t.Fatalf("three attempts on one instance counted as %d instances", same.Instances)
	}
	// Difficulty 0.5 rather than 1.0: at the top of the range the score saturates at the
	// clamp before variety can add anything, which is fine — there is nothing left to
	// discriminate — but it would hide the effect being asserted here.
	varied := BuildEvidence([]Attempt{
		{Correct: true, Elapsed: 2, Difficulty: diff(0.5), Instance: "a"},
		{Correct: true, Elapsed: 2, Difficulty: diff(0.5), Instance: "b"},
		{Correct: true, Elapsed: 2, Difficulty: diff(0.5), Instance: "c"},
	}, EvidenceOptions{TimeThreshold: 10})
	if varied.Score <= same.Score {
		t.Errorf("variety did not raise the score: %v vs %v", varied.Score, same.Score)
	}
}

func TestBuildEvidence_DifficultyAndTimeChangeTheScore(t *testing.T) {
	insts := func() []Attempt {
		return []Attempt{
			{Correct: true, Elapsed: 5, Difficulty: diff(0.3), Instance: "a"},
			{Correct: true, Elapsed: 5, Difficulty: diff(0.3), Instance: "b"},
			{Correct: true, Elapsed: 5, Difficulty: diff(0.3), Instance: "c"},
		}
	}
	easy := BuildEvidence(insts(), EvidenceOptions{TimeThreshold: 10})

	hard := insts()
	for i := range hard {
		hard[i].Difficulty = diff(1.0)
	}
	hardEv := BuildEvidence(hard, EvidenceOptions{TimeThreshold: 10})
	if hardEv.Score <= easy.Score {
		t.Errorf("harder questions scored %v, not above easy %v", hardEv.Score, easy.Score)
	}

	slow := insts()
	for i := range slow {
		slow[i].Elapsed = 60
	}
	slowEv := BuildEvidence(slow, EvidenceOptions{TimeThreshold: 10})
	if slowEv.Score >= easy.Score {
		t.Errorf("slow answers scored %v, not below fast %v", slowEv.Score, easy.Score)
	}
}

// Unknown difficulty must be read as neither zero nor perfect.
//
// The first version of this test asserted the opposite: it required the score to come out at
// 1.0, which is what "drop the difficulty term" produces. Dropping a *multiplier* means
// multiplying by 1.0 — the top of the range — so a missing value was being credited as
// "the hardest question we could have asked". The name said "dropped, not zeroed" and both
// halves of that were wrong: it was not dropped, and the effect was maximal.
//
// A concept whose attempts predate the column still has evidence. It just does not have
// evidence of how hard the questions were, and it is credited at NeutralDifficulty — the
// difficulty the engine serves when it knows nothing — rather than at either extreme.
func TestBuildEvidence_UnknownDifficultyIsNeitherZeroNorPerfect(t *testing.T) {
	unknown := unknownDifficultyWindow(MasteryEvidenceFloor)
	if unknown.AvgDifficulty != 0 {
		t.Errorf("AvgDifficulty = %v, want 0 when nothing recorded it", unknown.AvgDifficulty)
	}

	// Identical evidence at the neutral difficulty must score the same, because that is
	// what NeutralDifficulty claims.
	neutral := cleanWindow(MasteryEvidenceFloor, NeutralDifficulty)
	if math.Abs(unknown.Score-neutral.Score) > 1e-9 {
		t.Errorf("unknown difficulty scored %v, difficulty %v scored %v; unknown should be "+
			"credited as the neutral question", unknown.Score, NeutralDifficulty, neutral.Score)
	}

	// And it must sit strictly below the hardest question, which is what the old 1.0
	// default scored it as.
	hardest := cleanWindow(MasteryEvidenceFloor, 1.0)
	if !(unknown.Score < hardest.Score) {
		t.Errorf("unknown difficulty scored %v, not below the hardest question's %v — a "+
			"missing value is being treated as maximal evidence", unknown.Score, hardest.Score)
	}

	// The property that actually matters: it must not disqualify otherwise-good evidence.
	if !EnoughForMastery(unknown) {
		t.Errorf("unknown difficulty should not disqualify otherwise-good evidence: %+v", unknown)
	}
}

func TestBuildEvidence_OnlyTheRecentWindowCounts(t *testing.T) {
	var old []Attempt
	for i := 0; i < 20; i++ {
		old = append(old, Attempt{Correct: false, Elapsed: 5, Instance: "old"})
	}
	for i := 0; i < 4; i++ {
		old = append(old, Attempt{
			Correct: true, Elapsed: 5, Difficulty: diff(0.8),
			Instance: string(rune('a' + i)),
		})
	}
	ev := BuildEvidence(old, EvidenceOptions{TimeThreshold: 10})
	if ev.Attempts != EvidenceWindow {
		t.Errorf("Attempts = %d, want the window of %d", ev.Attempts, EvidenceWindow)
	}
	// The window reaches back two of the 20 misses, and still clears. A month of failure
	// must not prevent mastery once the recent evidence is clean — the question is whether
	// the learner can do it *now*, and a streak counter that never decays is the reason this
	// window exists.
	if !EnoughForMastery(ev) {
		t.Errorf("recent clean evidence did not establish mastery: %+v", ev)
	}

	// The converse, and the reason the window is a window: two corrects after 20 misses is
	// genuinely weak, and the model says so rather than being helped along by the old data.
	weak := BuildEvidence(append(old[:20],
		Attempt{Correct: true, Elapsed: 5, Difficulty: diff(0.8), Instance: "a"},
		Attempt{Correct: true, Elapsed: 5, Difficulty: diff(0.8), Instance: "b"},
	), EvidenceOptions{TimeThreshold: 10})
	if EnoughForMastery(weak) {
		t.Errorf("two correct answers after 20 misses should not establish mastery: %+v", weak)
	}
}

func TestBuildEvidence_EmptyIsUnevidenced(t *testing.T) {
	ev := BuildEvidence(nil, EvidenceOptions{TimeThreshold: 10})
	if ev.Attempts != 0 || ev.Score != 0 {
		t.Fatalf("empty evidence = %+v", ev)
	}
	if EnoughForMastery(ev) || EnoughForPracticing(ev) {
		t.Error("no attempts must not establish anything")
	}
}

func TestEnoughForPracticing_WeakerThanMastery(t *testing.T) {
	// Two correct attempts on two instances, all at the easiest difficulty: enough to show
	// consistency, not enough to show the concept is learned.
	ev := BuildEvidence([]Attempt{
		{Correct: true, Elapsed: 5, Difficulty: diff(0.3), Instance: "a"},
		{Correct: true, Elapsed: 5, Difficulty: diff(0.3), Instance: "b"},
	}, EvidenceOptions{TimeThreshold: 10})
	if !EnoughForPracticing(ev) {
		t.Errorf("two clean correct attempts should show consistency: %+v", ev)
	}
	if EnoughForMastery(ev) {
		t.Errorf("two attempts should not establish mastery: %+v", ev)
	}
}

// The ladder, as the state diagram draws it: three different questions, not one question
// asked three times.
//
// This is the test that would have caught the collapse. With every rung asking
// `streak >= required`, a threshold-10 concept needed 30 consecutive correct answers, and
// the cost fell entirely on the learner being stuck at UNSEEN.
func TestLadder_EachRungAsksADifferentQuestion(t *testing.T) {
	// 1. Nothing yet.
	none := BuildEvidence(nil, EvidenceOptions{TimeThreshold: 10})
	if HighestRung(none) != RungNone {
		t.Errorf("no evidence gave rung %v, want RungNone", HighestRung(none))
	}

	// 2. One correct answer: started, and nothing more. The old ladder needed a full
	// threshold-length streak to say a learner had begun learning.
	one := BuildEvidence([]Attempt{
		{Correct: true, Elapsed: 4, Difficulty: diff(0.5), Instance: "a"},
	}, EvidenceOptions{TimeThreshold: 10})
	if HighestRung(one) != RungStarted {
		t.Errorf("one correct answer gave rung %v, want RungStarted", HighestRung(one))
	}

	// 3. Two correct, two instances, easy: reliable, not yet demonstrated.
	two := BuildEvidence([]Attempt{
		{Correct: true, Elapsed: 4, Difficulty: diff(0.4), Instance: "a"},
		{Correct: true, Elapsed: 4, Difficulty: diff(0.4), Instance: "b"},
	}, EvidenceOptions{TimeThreshold: 10})
	if HighestRung(two) != RungConsistent {
		t.Errorf("two consistent correct answers gave rung %v, want RungConsistent", HighestRung(two))
	}

	// 4. A full window of correct answers across distinct instances at real difficulty:
	// demonstrated. Four is no longer enough — the floor is MasteryEvidenceFloor — which is
	// the calibration change, not a new kind of evidence.
	four := cleanWindow(MasteryEvidenceFloor, 0.8)
	if HighestRung(four) != RungDemonstrated {
		t.Errorf("a full clean window gave rung %v, want RungDemonstrated", HighestRung(four))
	}

	// Just under the floor, the top rung is out of reach however good the answers are. This
	// is the whole of the calibration: the floor is on *how much* evidence, not on what it
	// contains.
	short := cleanWindow(MasteryEvidenceFloor-1, 0.9)
	if short.Score < MasteryThreshold {
		t.Fatalf("test setup: %d attempts scored %v, which should clear the threshold",
			MasteryEvidenceFloor-1, short.Score)
	}
	if HighestRung(short) == RungDemonstrated {
		t.Errorf("%d attempts reached the top rung; the floor is %d", MasteryEvidenceFloor-1, MasteryEvidenceFloor)
	}
}

// More evidence, or better evidence, must never lower the rung.
//
// Monotone, not strictly increasing: the rungs saturate. Two correct answers on two
// instances and three correct answers on three instances are both RungConsistent, because
// both clear the threshold and the attempt floor for that rung is two. Asserting strict
// increase would be asserting a distinction the ladder does not make and should not.
func TestLadder_RungsAreMonotone(t *testing.T) {
	cases := []struct {
		name string
		ev   Evidence
	}{
		{"none", BuildEvidence(nil, EvidenceOptions{})},
		{"one", BuildEvidence([]Attempt{{Correct: true, Instance: "a"}}, EvidenceOptions{})},
		{"two easy", cleanWindow(2, 0.4)},
		{"three easy", cleanWindow(3, 0.4)},
		{"five easy", cleanWindow(MasteryEvidenceFloor-1, 0.4)},
		{"full window easy", cleanWindow(MasteryEvidenceFloor, 0.4)},
		{"full window hard", cleanWindow(MasteryEvidenceFloor, 0.9)},
	}
	for i := 1; i < len(cases); i++ {
		prev, cur := HighestRung(cases[i-1].ev), HighestRung(cases[i].ev)
		if cur < prev {
			t.Errorf("%s rung %v is below %s rung %v",
				cases[i].name, cur, cases[i-1].name, prev)
		}
	}
	// And the ends of the range are reached, so the walk above is not vacuous.
	if HighestRung(cases[0].ev) != RungNone {
		t.Error("no evidence should reach no rung")
	}
	if HighestRung(cases[len(cases)-1].ev) != RungDemonstrated {
		t.Error("a full window of hard correct answers should reach the top rung")
	}
}

// Difficulty must not be able to *lower* the rung. It weights the score upward, so this is
// the same monotonicity read through the other signal — and it is the one that would break
// if the difficulty term were ever subtracted rather than multiplied in.
func TestLadder_HarderEvidenceNeverHurts(t *testing.T) {
	easy := HighestRung(cleanWindow(MasteryEvidenceFloor, 0.3))
	hard := HighestRung(cleanWindow(MasteryEvidenceFloor, 1.0))
	if hard < easy {
		t.Errorf("harder evidence gave rung %v, below easy %v", hard, easy)
	}
	// Slow answers likewise.
	fast := HighestRung(cleanWindow(MasteryEvidenceFloor, 0.8))
	slow := BuildEvidence(func() []Attempt {
		var as []Attempt
		for i := 0; i < MasteryEvidenceFloor; i++ {
			as = append(as, evidenceOf(true, string(rune('a'+i)), 0.8, 600))
		}
		return as
	}(), EvidenceOptions{TimeThreshold: 10})
	if HighestRung(slow) > fast {
		t.Errorf("slow answers gave a higher rung (%v) than fast (%v)", HighestRung(slow), fast)
	}
}
