package scheduler

import (
	"strings"
	"testing"
	"time"

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
		t.Error("a new learner got the diagnostic and nothing else — no agency")
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
		t.Errorf("primary = %q, want b or c — the unlocked, unattempted concepts", got.Primary.ConceptID)
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

// A learner with nothing eligible — nothing reachable, nothing due, nothing weak — is a real
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
