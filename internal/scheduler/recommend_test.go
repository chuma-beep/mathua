package scheduler

import (
	"strings"
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/mastery"
)

// These tests use the package's existing miniDAG fixture (a -> b, a -> c, b+c -> d) so they
// agree with the scheduler tests above about graph shape. With `a` mastered, `b` and `c` are
// available and `d` is blocked; that is enough to tell "first in database order" apart from
// "most valuable".

func masteredAt(now time.Time) *ConceptSnapshot {
	return &ConceptSnapshot{Status: mastery.StatusMastered, LastAttempted: now, LastReviewed: now}
}

func startedAt(now time.Time, streak int) *ConceptSnapshot {
	return &ConceptSnapshot{Status: mastery.StatusPracticing, LastAttempted: now, Streak: streak}
}

func all(r RecommendationResponse) []Recommendation {
	return append([]Recommendation{*r.Primary}, r.Alternatives...)
}

func TestRecommend_NewLearnerIsOfferedTheDiagnosticButStillHasAgency(t *testing.T) {
	s := New(miniDAG(t))
	got := s.Recommend(RecommendInput{Now: time.Now()})
	if got.Primary == nil || got.Primary.Action.Href != "/onboard" {
		t.Fatalf("primary = %+v, want the diagnostic for a learner with no history", got.Primary)
	}
	if len(got.Alternatives) == 0 {
		t.Error("a new learner got the diagnostic and nothing else, no agency")
	}
}

func TestRecommend_ActiveLearnerIsNeverSentToTheDiagnostic(t *testing.T) {
	s := New(miniDAG(t))
	got := s.Recommend(RecommendInput{Now: time.Now(), AttemptedDays: 5})
	if got.Primary != nil && got.Primary.Action.Href == "/onboard" {
		t.Error("a learner with five days of attempts was sent to the diagnostic")
	}
}

func TestRecommend_MasteredAndBlockedConceptsAreNeverOffered(t *testing.T) {
	s := New(miniDAG(t))
	now := time.Now()
	got := s.Recommend(RecommendInput{
		Now: now, AttemptedDays: 10, ConceptsMastered: 1,
		Snapshots: map[string]*ConceptSnapshot{
			"a": masteredAt(now),
			"d": startedAt(now, 1), // started, but b and c are not done: still blocked
		},
	})
	for _, r := range all(got) {
		switch r.ConceptID {
		case "a":
			t.Errorf("mastered concept %q was offered", r.ID)
		case "d":
			t.Errorf("%q was offered while blocked on unmet prerequisites b and c", r.ID)
		}
	}
	if got.Primary.ConceptID != "b" && got.Primary.ConceptID != "c" {
		t.Errorf("primary = %q, want b or c, the unlocked, unattempted concepts", got.Primary.ConceptID)
	}
}

func TestRecommend_DecayedConceptBecomesReviewNotTeaching(t *testing.T) {
	s := New(miniDAG(t))
	now := time.Now()
	stale := now.AddDate(0, 0, -(mastery.DecayDays + 5))
	got := s.Recommend(RecommendInput{
		Now: now, AttemptedDays: 30, ConceptsMastered: 2, DueReviews: 1,
		Snapshots: map[string]*ConceptSnapshot{
			"a": {Status: mastery.StatusMastered, LastAttempted: stale, LastReviewed: stale},
		},
	})
	// ADR-037's rule, on the surface that now owns it: a decayed concept is not offered as
	// new teaching.
	for _, r := range all(got) {
		if r.ConceptID == "a" && r.Kind == KindLearn {
			t.Errorf("decayed concept offered as learn: %+v", r)
		}
	}
	if got.Primary == nil || got.Primary.Kind != KindReview {
		t.Fatalf("primary = %+v, want a review while one is due", got.Primary)
	}
	if got.Primary.Action.Href != "/review" {
		t.Errorf("review href = %q, want /review", got.Primary.Action.Href)
	}
}

func TestRecommend_WeakConceptIsOfferedAsPracticeNotNew(t *testing.T) {
	s := New(miniDAG(t))
	now := time.Now()
	got := s.Recommend(RecommendInput{
		Now: now, AttemptedDays: 20, ConceptsMastered: 1,
		Weakness:  map[string]float64{"b": 0.9},
		Snapshots: map[string]*ConceptSnapshot{"a": masteredAt(now), "b": startedAt(now, 1)},
	})
	var practice *Recommendation
	for i, r := range all(got) {
		if r.ConceptID == "b" {
			practice = &all(got)[i]
		}
	}
	if practice == nil {
		t.Fatalf("weak concept b was not offered: %+v", got)
	}
	if practice.Kind != KindPractice || practice.Reason != ReasonNeedsPractice {
		t.Errorf("b came back as %s/%s, want practice/needs_practice", practice.Kind, practice.Reason)
	}
	if practice.Action.Href != "/learn?concept=b" {
		t.Errorf("practice href = %q", practice.Action.Href)
	}
}

func TestRecommend_MasteredPrerequisiteWithWeaknessOfferedForRemediation(t *testing.T) {
	s := New(miniDAG(t))
	now := time.Now()
	// Concept "a" is mastered but elevated weakness was propagated backward from downstream struggle.
	got := s.Recommend(RecommendInput{
		Now:              now,
		AttemptedDays:    15,
		ConceptsMastered: 1,
		Weakness:         map[string]float64{"a": 0.45},
		Snapshots: map[string]*ConceptSnapshot{
			"a": {
				Status:        mastery.StatusMastered,
				Streak:        3,
				WeaknessScore: 0.45,
				LastAttempted: now,
				LastReviewed:  now,
			},
		},
	})
	var remediation *Recommendation
	for i, r := range all(got) {
		if r.ConceptID == "a" {
			remediation = &all(got)[i]
		}
	}
	if remediation == nil {
		t.Fatalf("mastered prerequisite with weakness 0.45 was not offered for remediation: %+v", got)
	}
	if remediation.Kind != KindPractice {
		t.Errorf("remediation kind = %q, want %q", remediation.Kind, KindPractice)
	}
	if remediation.Reason != ReasonPrerequisite {
		t.Errorf("remediation reason = %q, want %q", remediation.Reason, ReasonPrerequisite)
	}
	if remediation.Badge != "Remediation" {
		t.Errorf("remediation badge = %q, want Remediation", remediation.Badge)
	}
	if remediation.CTA != "Review prerequisite →" {
		t.Errorf("remediation cta = %q, want Review prerequisite →", remediation.CTA)
	}
}

func TestScheduler_BackwardRemediationOutranksDownstreamStruggle(t *testing.T) {
	s := New(miniDAG(t))
	now := time.Now()
	// In miniDAG: a -> b.
	// Learner has mastered "a" with streak 3, but downstream struggle on "b" elevated "a"'s weakness to 0.45.
	// Meanwhile "b" has weakness 0.7.
	// Prerequisite remediation for "a" must receive higher priority than repeating "b".
	candidates := buildCandidates(s.dag, map[string]*ConceptSnapshot{
		"a": {
			Status:        mastery.StatusMastered,
			Streak:        3,
			WeaknessScore: 0.45,
			LastAttempted: now,
			LastReviewed:  now,
		},
		"b": {
			Status:        mastery.StatusPracticing,
			Streak:        0,
			WeaknessScore: 0.7,
			LastAttempted: now,
		},
	}, now, "")

	var prioA, prioB float64
	foundA, foundB := false, false
	for _, c := range candidates {
		if c.Concept.ID == "a" {
			prioA = c.Priority
			foundA = true
		}
		if c.Concept.ID == "b" {
			prioB = c.Priority
			foundB = true
		}
	}
	if !foundA || !foundB {
		t.Fatalf("expected both a and b in candidates: foundA=%v foundB=%v candidates=%+v", foundA, foundB, candidates)
	}
	if prioA <= prioB {
		t.Errorf("prerequisite remediation priority %v should exceed downstream struggle priority %v", prioA, prioB)
	}
}

// A remediation card and the concept the learner is failing are two different tasks, so the
// shelf has to hold both. They used to share the practice bucket, and it has one slot: the
// mastered prerequisite took it and the failing concept disappeared from the response.
func TestRecommend_RemediationKeepsTheFailingConceptOnTheShelf(t *testing.T) {
	s := New(miniDAG(t))
	now := time.Now()
	// The weakness map is what the engine passes: zero for a mastered concept, the stored
	// score for the rest. The elevated prerequisite weakness reaches the scheduler through
	// the snapshot, which is why both are set here.
	got := s.Recommend(RecommendInput{
		Now:              now,
		AttemptedDays:    15,
		ConceptsMastered: 1,
		Weakness:         map[string]float64{"a": 0, "b": 0.8, "c": 0, "d": 0},
		Snapshots: map[string]*ConceptSnapshot{
			"a": {Status: mastery.StatusMastered, Streak: 3, WeaknessScore: 0.45, LastAttempted: now, LastReviewed: now},
			"b": {Status: mastery.StatusPracticing, WeaknessScore: 0.8, LastAttempted: now},
		},
	})
	if got.Primary == nil || got.Primary.ConceptID != "a" {
		t.Fatalf("primary = %+v, want the prerequisite remediation for a", got.Primary)
	}
	var kept *Recommendation
	for i, r := range all(got) {
		if r.ConceptID == "b" {
			kept = &all(got)[i]
		}
	}
	if kept == nil {
		t.Fatalf("the concept the learner is failing was dropped from the shelf: %+v", got)
	}
	if kept.Reason != ReasonNeedsPractice {
		t.Errorf("b reason = %q, want %q", kept.Reason, ReasonNeedsPractice)
	}
}

func TestRecommend_ExcludedConceptIsNotOfferedBack(t *testing.T) {
	s := New(miniDAG(t))
	now := time.Now()
	got := s.Recommend(RecommendInput{
		Now: now, AttemptedDays: 10, ConceptsMastered: 1,
		ExcludeConceptIds: []string{"b"},
		Snapshots:         map[string]*ConceptSnapshot{"a": masteredAt(now)},
	})
	for _, r := range all(got) {
		if r.ConceptID == "b" {
			t.Error("a just-finished concept was offered back as the next task")
		}
	}
}

func TestRecommend_SetHasNoDuplicatesAndRespectsTheCap(t *testing.T) {
	s := New(miniDAG(t))
	now := time.Now()
	snaps := map[string]*ConceptSnapshot{
		"a": masteredAt(now), "b": startedAt(now, 1), "c": startedAt(now, 1),
	}
	got := s.Recommend(RecommendInput{Now: now, AttemptedDays: 10, ConceptsMastered: 1, Snapshots: snaps})
	seen := map[string]bool{}
	list := all(got)
	for _, r := range list {
		if seen[r.ConceptID] {
			t.Errorf("duplicate recommendation for %q", r.ConceptID)
		}
		seen[r.ConceptID] = true
	}
	if len(list) > MaxRecommendations {
		t.Errorf("%d recommendations, cap is %d", len(list), MaxRecommendations)
	}
}

func TestRecommend_MasteryCheckIsATaskNotABanner(t *testing.T) {
	s := New(miniDAG(t))
	now := time.Now()
	// The review case needs a review candidate to actually exist: a concept that has decayed.
	// Passing DueReviews: 1 against an otherwise-fresh graph states a count the graph cannot
	// support, and the ranking correctly has nothing to outrank.
	stale := now.AddDate(0, 0, -(mastery.DecayDays + 3))
	decayed := map[string]*ConceptSnapshot{"a": {Status: mastery.StatusMastered,
		LastAttempted: stale, LastReviewed: stale}}
	withReview := s.Recommend(RecommendInput{Now: now, AttemptedDays: 30,
		ConceptsMastered: 5, DueReviews: 1, QuizDue: true, Snapshots: decayed})
	if withReview.Primary.Kind != KindReview {
		t.Errorf("primary = %s, want the due review to outrank the quiz", withReview.Primary.Kind)
	}
	if !hasKind(withReview.Alternatives, KindMasteryCheck) {
		t.Error("quiz due but no mastery-check task anywhere in the set")
	}

	noReview := s.Recommend(RecommendInput{Now: now, AttemptedDays: 30,
		ConceptsMastered: 5, QuizDue: true,
		Snapshots: map[string]*ConceptSnapshot{"a": masteredAt(now)}})
	if noReview.Primary.Kind != KindMasteryCheck {
		t.Fatalf("primary = %s, want the mastery check when nothing is due", noReview.Primary.Kind)
	}
	if noReview.Primary.Reason != ReasonMasteryReady {
		t.Errorf("reason = %s, want %s", noReview.Primary.Reason, ReasonMasteryReady)
	}
	if noReview.Primary.Action.Href != "/goals?quiz=1" {
		t.Errorf("href = %q", noReview.Primary.Action.Href)
	}
	// The gate is a schedule, not learner-facing economy: the interface says
	// when a check is ready, never how many XP it takes.
	if strings.Contains(noReview.Primary.Detail, "XP") {
		t.Errorf("mastery-check detail prints the gate: %q", noReview.Primary.Detail)
	}
}

func hasKind(rs []Recommendation, k Kind) bool {
	for _, r := range rs {
		if r.Kind == k {
			return true
		}
	}
	return false
}

func TestRecommend_RankingRespondsToLearnerState(t *testing.T) {
	s := New(miniDAG(t))
	now := time.Now()

	// Weakness pulls a concept up. b and c are otherwise symmetric.
	plain := s.Recommend(RecommendInput{Now: now, AttemptedDays: 10, ConceptsMastered: 1,
		Snapshots: map[string]*ConceptSnapshot{"a": masteredAt(now)}})
	weak := s.Recommend(RecommendInput{Now: now, AttemptedDays: 10, ConceptsMastered: 1,
		Weakness:  map[string]float64{"c": 0.95},
		Snapshots: map[string]*ConceptSnapshot{"a": masteredAt(now)}})

	plainB := rankOf(plain, "b")
	plainC := rankOf(plain, "c")
	weakC := rankOf(weak, "c")
	// A better rank is a lower number.
	if weakC >= plainC {
		t.Errorf("raising c's weakness to 0.95 did not improve its rank (%d -> %d); plain b=%d c=%d",
			plainC, weakC, plainB, plainC)
	}
}

func rankOf(r RecommendationResponse, id string) int {
	for i, x := range all(r) {
		if x.ConceptID == id {
			return i
		}
	}
	return 99
}

func TestRecommend_EveryRecommendationCarriesItsOwnCopy(t *testing.T) {
	// The duplication defect was two components wording one head differently ("Continue: X" in
	// one card, "Currently working towards X" in the card above it), so the wording is owned
	// here and every recommendation must carry a complete set.
	s := New(miniDAG(t))
	now := time.Now()
	got := s.Recommend(RecommendInput{Now: now, AttemptedDays: 10, ConceptsMastered: 1,
		Snapshots: map[string]*ConceptSnapshot{"a": masteredAt(now), "b": startedAt(now, 1)}})
	for _, r := range all(got) {
		if r.Badge == "" || r.Detail == "" || r.CTA == "" {
			t.Errorf("%q has incomplete copy: badge=%q detail=%q cta=%q", r.ID, r.Badge, r.Detail, r.CTA)
		}
		if r.Action.Href == "" {
			t.Errorf("%q has no href", r.ID)
		}
		if r.ConceptTitle == "" {
			t.Errorf("%q has no title", r.ID)
		}
	}
}

func TestRecommend_XPAndTimeComeFromTheOwningPackages(t *testing.T) {
	// minDAG's concepts have AvgTimeSeconds 10, which EffortBase prices at the 1-XP floor.
	// The point of the test is that the numbers arrive at all, and from the one place: the
	// client's copy of the award math had already drifted (it still halved reviews).
	s := New(miniDAG(t))
	now := time.Now()
	got := s.Recommend(RecommendInput{Now: now, AttemptedDays: 10, ConceptsMastered: 1,
		Snapshots: map[string]*ConceptSnapshot{"a": masteredAt(now)}})
	for _, r := range all(got) {
		if r.ConceptID == "" {
			continue // the diagnostic and the library fallback carry no concept
		}
		if r.XP == nil || *r.XP <= 0 {
			t.Errorf("%q carries no XP; the recommendation is what tells a learner what a task is worth", r.ID)
		}
		if r.EstimatedMinutes == nil || *r.EstimatedMinutes <= 0 {
			t.Errorf("%q carries no estimate", r.ID)
		}
	}
}

func TestRecommend_NoLearnerFacingTerminology(t *testing.T) {
	// Asserted here rather than in a copy test, so a new copy string cannot reintroduce it.
	s := New(miniDAG(t))
	now := time.Now()
	got := s.Recommend(RecommendInput{Now: now, AttemptedDays: 40, ConceptsMastered: 3,
		DueReviews: 2, QuizDue: true,
		Snapshots: map[string]*ConceptSnapshot{
			"a": {Status: mastery.StatusMastered,
				LastAttempted: now.AddDate(0, 0, -(mastery.DecayDays + 3)),
				LastReviewed:  now.AddDate(0, 0, -(mastery.DecayDays + 3))},
			"b": startedAt(now, 1), "c": startedAt(now, 1),
		}})
	bad := []string{"frontier", "Frontier", "DAG", "scheduler", "candidate", "unlocked node", "prerequisite graph"}
	for _, r := range all(got) {
		for _, f := range []string{r.Badge, r.Detail, r.CTA} {
			for _, b := range bad {
				if strings.Contains(f, b) {
					t.Errorf("learner-facing copy %q contains internal term %q", f, b)
				}
			}
		}
	}
}

// A learner with nothing eligible, nothing reachable, nothing due, nothing weak, is a real
// state, and it must still get a head: the client cannot render a null primary.
//
// It used to be answered with `/study`, and the fallback is worth keeping as a test because of
// what it says about the product. "Everything available is already learned, the reference
// library is open" told a learner with nothing due to go and read, which is the one activity
// that cannot help them. The graph is orientation: it claims nothing, it cannot 401, and it
// shows what builds on what they already hold.
func TestRecommend_NothingAvailableFallsBackToOrientation(t *testing.T) {
	s := New(miniDAG(t))
	now := time.Now()
	got := s.Recommend(RecommendInput{Now: now, AttemptedDays: 90, ConceptsMastered: 4,
		Snapshots: map[string]*ConceptSnapshot{
			"a": masteredAt(now), "b": masteredAt(now), "c": masteredAt(now), "d": masteredAt(now),
		}})
	if got.Primary == nil {
		t.Fatal("no primary when everything is learned")
	}
	if got.Primary.Action.Href != "/graph" {
		t.Errorf("primary href = %q, want /graph", got.Primary.Action.Href)
	}
	// Every destination this engine can name must be one the product still has. /study was a
	// live route when this fallback was written and is a forwarder to /learn now, so a
	// recommendation pointing at it would send a learner to a redirect on the way to nothing.
	if strings.HasPrefix(got.Primary.Action.Href, "/study") {
		t.Errorf("primary points at the closed reference library: %q", got.Primary.Action.Href)
	}
}

func TestRecommend_TargetedRemediation_KeyPrerequisitesInjectedAtFront(t *testing.T) {
	d, err := concepts.Build([]concepts.Concept{
		{
			ID: "a", Label: "Concept A", Domain: "d", Subdomain: "s-a",
			Prerequisites:    []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 10},
		},
		{
			ID: "b", Label: "Concept B", Domain: "d", Subdomain: "s-b",
			Prerequisites:    []string{"a"},
			KeyPrerequisites: []string{"a"},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 10},
		},
	})
	if err != nil {
		t.Fatalf("build DAG: %v", err)
	}
	s := New(d)
	now := time.Now()

	// Learner halts on b with key prerequisite a injected into the task shelf.
	got := s.Recommend(RecommendInput{
		Now:              now,
		AttemptedDays:    10,
		ConceptsMastered: 1,
		KeyPrerequisites: []string{"a"},
		Snapshots: map[string]*ConceptSnapshot{
			"a": masteredAt(now),
			"b": {Status: mastery.StatusPracticing, Streak: 0, WeaknessScore: 0.8, LastAttempted: now},
		},
	})

	if got.Primary == nil {
		t.Fatal("expected primary recommendation")
	}
	if got.Primary.ConceptID != "a" {
		t.Errorf("primary concept = %q, want key prerequisite a", got.Primary.ConceptID)
	}
	if got.Primary.Reason != ReasonPrerequisite {
		t.Errorf("primary reason = %q, want %q", got.Primary.Reason, ReasonPrerequisite)
	}
	if got.Primary.Badge != "Remediation" {
		t.Errorf("primary badge = %q, want Remediation", got.Primary.Badge)
	}
	if got.Primary.CTA != "Review prerequisite →" {
		t.Errorf("primary cta = %q, want Review prerequisite →", got.Primary.CTA)
	}
	if got.Primary.Action.Href != "/learn?concept=a" {
		t.Errorf("primary href = %q, want /learn?concept=a", got.Primary.Action.Href)
	}
}

func TestRecommend_TargetedRemediation_OutranksQuizDue(t *testing.T) {
	d, err := concepts.Build([]concepts.Concept{
		{
			ID: "a", Label: "Concept A", Domain: "d", Subdomain: "s-a",
			Prerequisites:    []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 10},
		},
		{
			ID: "b", Label: "Concept B", Domain: "d", Subdomain: "s-b",
			Prerequisites:    []string{"a"},
			KeyPrerequisites: []string{"a"},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 10},
		},
	})
	if err != nil {
		t.Fatalf("build DAG: %v", err)
	}
	s := New(d)
	now := time.Now()

	// Targeted remediation outranks QuizDue at the head of the task shelf.
	got := s.Recommend(RecommendInput{
		Now:              now,
		AttemptedDays:    10,
		ConceptsMastered: 1,
		QuizDue:          true,
		KeyPrerequisites: []string{"a"},
		Snapshots: map[string]*ConceptSnapshot{
			"a": masteredAt(now),
			"b": {Status: mastery.StatusPracticing, Streak: 0, WeaknessScore: 0.8, LastAttempted: now},
		},
	})

	if got.Primary == nil || got.Primary.ConceptID != "a" {
		t.Fatalf("primary = %+v, want key prerequisite a", got.Primary)
	}
	if !hasKind(got.Alternatives, KindMasteryCheck) {
		t.Errorf("expected mastery check in alternatives, got %+v", got.Alternatives)
	}
}

func TestRecommend_TargetedRemediation_MultipleKeyPrerequisites(t *testing.T) {
	d, err := concepts.Build([]concepts.Concept{
		{
			ID: "a", Label: "Concept A", Domain: "d", Subdomain: "s-a",
			Prerequisites:    []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 10},
		},
		{
			ID: "b", Label: "Concept B", Domain: "d", Subdomain: "s-b",
			Prerequisites:    []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 10},
		},
		{
			ID: "c", Label: "Concept C", Domain: "d", Subdomain: "s-c",
			Prerequisites:    []string{"a", "b"},
			KeyPrerequisites: []string{"a", "b"},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 10},
		},
	})
	if err != nil {
		t.Fatalf("build DAG: %v", err)
	}
	s := New(d)
	now := time.Now()

	got := s.Recommend(RecommendInput{
		Now:              now,
		AttemptedDays:    15,
		ConceptsMastered: 2,
		HaltedConceptID:  "c",
		Snapshots: map[string]*ConceptSnapshot{
			"a": masteredAt(now),
			"b": masteredAt(now),
			"c": {Status: mastery.StatusPracticing, Streak: 0, WeaknessScore: 0.8, LastAttempted: now},
		},
	})

	if got.Primary == nil || got.Primary.ConceptID != "a" {
		t.Fatalf("primary = %+v, want first key prerequisite a", got.Primary)
	}
	if len(got.Alternatives) == 0 || got.Alternatives[0].ConceptID != "b" {
		t.Fatalf("expected second key prerequisite b as first alternative, got %+v", got.Alternatives)
	}
}

func TestRecommend_TargetedRemediation_GracefulFallbackWhenNoKeyPrerequisites(t *testing.T) {
	d, err := concepts.Build([]concepts.Concept{
		{
			ID: "a", Label: "Concept A", Domain: "d", Subdomain: "s-a",
			Prerequisites:    []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 10},
		},
	})
	if err != nil {
		t.Fatalf("build DAG: %v", err)
	}
	s := New(d)
	now := time.Now()

	// Halted concept with no KeyPrerequisites falls back to standard remediation path.
	got := s.Recommend(RecommendInput{
		Now:              now,
		AttemptedDays:    5,
		ConceptsMastered: 0,
		HaltedConceptID:  "a",
		Snapshots: map[string]*ConceptSnapshot{
			"a": {Status: mastery.StatusPracticing, Streak: 1, WeaknessScore: 0.5, LastAttempted: now},
		},
	})

	if got.Primary == nil {
		t.Fatal("expected primary recommendation on fallback path")
	}
	if got.Primary.ConceptID != "a" {
		t.Errorf("primary = %q, want concept a", got.Primary.ConceptID)
	}
	if got.Primary.Reason != ReasonNeedsPractice {
		t.Errorf("primary reason = %q, want %q", got.Primary.Reason, ReasonNeedsPractice)
	}
}

func TestRecommend_TargetedRemediation_NormalCompletionPath(t *testing.T) {
	d := miniDAG(t)
	s := New(d)
	now := time.Now()

	// When no key prerequisites or halts exist, recommendations follow normal completion path.
	got := s.Recommend(RecommendInput{
		Now:              now,
		AttemptedDays:    10,
		ConceptsMastered: 1,
		Snapshots: map[string]*ConceptSnapshot{
			"a": masteredAt(now),
		},
	})

	if got.Primary == nil {
		t.Fatal("expected primary recommendation on normal path")
	}
	if got.Primary.Reason != ReasonNew {
		t.Errorf("primary reason = %q, want %q", got.Primary.Reason, ReasonNew)
	}
}

func TestEffectiveState_HighFSRSStability_DoesNotDecayAt15Days(t *testing.T) {
	now := time.Now()
	lastRev := now.AddDate(0, 0, -15)
	due := now.AddDate(0, 0, 30)
	stability := 60.0

	snap := &ConceptSnapshot{
		Status:         mastery.StatusMastered,
		LastReviewed:   lastRev,
		NextReviewDue:  &due,
		FSRSStability:  &stability,
		FSRSLastReview: &lastRev,
	}

	status, isReview := effectiveState(snap, now)
	if status != mastery.StatusMastered {
		t.Errorf("status = %q, want %q at 15 days with stability 60", status, mastery.StatusMastered)
	}
	if isReview {
		t.Error("expected isReview=false at 15 days when stability is 60 and NextReviewDue is in future")
	}
}

func TestRecommend_TargetedRemediation_FallsBackToDAGPrereqWhenNoKeyPrerequisites(t *testing.T) {
	d, err := concepts.Build([]concepts.Concept{
		{
			ID: "p1", Label: "Prereq 1", Domain: "d", Subdomain: "s",
			Prerequisites:    []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 10},
		},
		{
			ID: "p2", Label: "Prereq 2", Domain: "d", Subdomain: "s",
			Prerequisites:    []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 10},
		},
		{
			ID: "c", Label: "Concept C", Domain: "d", Subdomain: "s",
			Prerequisites:    []string{"p1", "p2"},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 10},
		},
	})
	if err != nil {
		t.Fatalf("build DAG: %v", err)
	}
	s := New(d)
	now := time.Now()

	// Student halted on c. Concept c has no KeyPrerequisites declared.
	// p1 is mastered (streak 5). p2 is learning (streak 1).
	// Scheduler must fall back to the weakest prerequisite (p2) for targeted remediation.
	got := s.Recommend(RecommendInput{
		Now:              now,
		AttemptedDays:    5,
		ConceptsMastered: 1,
		HaltedConceptID:  "c",
		Snapshots: map[string]*ConceptSnapshot{
			"p1": {Status: mastery.StatusMastered, Streak: 5, LastReviewed: now},
			"p2": {Status: mastery.StatusLearning, Streak: 1, LastAttempted: now},
			"c":  {Status: mastery.StatusPracticing, Streak: 0, LastAttempted: now},
		},
	})

	if got.Primary == nil {
		t.Fatal("expected primary recommendation on remedial fallback")
	}
	if got.Primary.ConceptID != "p2" {
		t.Errorf("primary concept = %q, want weakest prerequisite p2", got.Primary.ConceptID)
	}
	if got.Primary.Reason != ReasonPrerequisite {
		t.Errorf("primary reason = %q, want %q", got.Primary.Reason, ReasonPrerequisite)
	}
}

// Report 345: a diagnostic miss put an abstract algebra concept in the
// remedial queue, and it was offered as a practice task while all three of its
// prerequisites were still not started. Remediation must respect the same
// reachability rule as every other candidate.
func TestRecommend_RemedialConceptWaitsForItsPrerequisites(t *testing.T) {
	d, err := concepts.Build([]concepts.Concept{
		{
			ID: "a", Label: "Concept A", Domain: "d", Subdomain: "s-a",
			Prerequisites:    []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 10},
		},
		{
			ID: "b", Label: "Concept B", Domain: "d", Subdomain: "s-b",
			Prerequisites:    []string{"a"},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 10},
		},
	})
	if err != nil {
		t.Fatalf("build DAG: %v", err)
	}
	s := New(d)
	now := time.Now()
	in := RecommendInput{
		Now:              now,
		AttemptedDays:    10,
		ConceptsMastered: 1,
		RemedialConcepts: []string{"b"},
		Snapshots: map[string]*ConceptSnapshot{
			"b": {Status: mastery.StatusMastered, WeaknessScore: 0.9, LastAttempted: now},
		},
	}

	got := s.Recommend(in)
	if got.Primary != nil && got.Primary.ConceptID == "b" {
		t.Fatal("remedial concept b was primary while its prerequisite a was unstarted")
	}
	for _, alt := range got.Alternatives {
		if alt.ConceptID == "b" {
			t.Fatalf("remedial concept b appeared as an alternative while a was unstarted: %+v", alt)
		}
	}

	in.Snapshots["a"] = masteredAt(now)
	got = s.Recommend(in)
	found := got.Primary != nil && got.Primary.ConceptID == "b"
	for _, alt := range got.Alternatives {
		if alt.ConceptID == "b" {
			found = true
		}
	}
	if !found {
		t.Fatalf("remedial concept b was not offered after its prerequisite was mastered: %+v", got)
	}
}
