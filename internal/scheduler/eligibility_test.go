package scheduler

import (
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/mastery"
)

// The diagnostic is placement, not an active task. A completed diagnostic must
// never be offered again, even when placement leaves the learner with no mastered
// concept and no attempted day — which is exactly what placement does, because it
// seeds reachability rather than mastery.
func TestRecommend_CompletedDiagnosticIsNeverOfferedAgain(t *testing.T) {
	s := New(miniDAG(t))
	got := s.Recommend(RecommendInput{
		Now:                 time.Now(),
		DiagnosticCompleted: true,
	})
	if got.Primary != nil && got.Primary.Action.Href == "/onboard" {
		t.Fatalf("a completed diagnostic was offered as an active task: %+v", got.Primary)
	}
}

// A placement-seeded prerequisite satisfies its successors for reachability, so
// the learner can start at the frontier the diagnostic placed them at.
func TestRecommend_PlacementSeededPrerequisiteUnlocksSuccessor(t *testing.T) {
	s := New(miniDAG(t))
	now := time.Now()
	got := s.Recommend(RecommendInput{
		Now: now, AttemptedDays: 1, DiagnosticCompleted: true,
		Snapshots: map[string]*ConceptSnapshot{
			"b": {Status: mastery.StatusPracticing, LastAttempted: now, Placement: true},
			"c": {Status: mastery.StatusPracticing, LastAttempted: now, Placement: true},
		},
	})
	found := false
	for _, r := range all(got) {
		if r.ConceptID == "d" {
			found = true
		}
	}
	if !found {
		t.Errorf("d was not offered though both its prerequisites were placement-seeded: %+v", all(got))
	}
}

func TestEligibilityOf_States(t *testing.T) {
	dag := miniDAG(t)
	now := time.Now()
	cases := []struct {
		name     string
		snaps    map[string]*ConceptSnapshot
		concept  string
		want     string
		eligible bool
	}{
		{
			name: "unseen with met prerequisites is unlocked",
			snaps: map[string]*ConceptSnapshot{
				"a": masteredAt(now),
			},
			concept: "b", want: StateUnlocked, eligible: true,
		},
		{
			name:    "unseen with untouched prerequisites is locked",
			snaps:   map[string]*ConceptSnapshot{},
			concept: "b", want: StateLocked, eligible: false,
		},
		{
			name: "started is in progress",
			snaps: map[string]*ConceptSnapshot{
				"a": masteredAt(now), "b": startedAt(now, 1),
			},
			concept: "b", want: StateInProgress, eligible: true,
		},
		{
			name: "mastered",
			snaps: map[string]*ConceptSnapshot{
				"a": masteredAt(now), "b": masteredAt(now),
			},
			concept: "b", want: StateMastered, eligible: true,
		},
		{
			name: "untouched prerequisites that are all started are upcoming",
			snaps: map[string]*ConceptSnapshot{
				"a": masteredAt(now), "b": startedAt(now, 1), "c": startedAt(now, 1),
			},
			concept: "d", want: StateUpcoming, eligible: false,
		},
		{
			name: "a placement-seeded prerequisite is satisfied",
			snaps: map[string]*ConceptSnapshot{
				"b": {Status: mastery.StatusPracticing, LastAttempted: now, Placement: true},
				"c": {Status: mastery.StatusPracticing, LastAttempted: now, Placement: true},
			},
			concept: "d", want: StateUnlocked, eligible: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := dag.Concept(tc.concept)
			if c == nil {
				t.Fatalf("concept %q not in fixture", tc.concept)
			}
			got := EligibilityOf(tc.snaps, now, c)
			if got.State != tc.want {
				t.Errorf("state = %q, want %q", got.State, tc.want)
			}
			if got.Eligible != tc.eligible {
				t.Errorf("eligible = %v, want %v", got.Eligible, tc.eligible)
			}
		})
	}
}

func TestEligibilityOf_LockedNamesTheMissingPrerequisites(t *testing.T) {
	dag := miniDAG(t)
	got := EligibilityOf(map[string]*ConceptSnapshot{}, time.Now(), dag.Concept("d"))
	if got.Eligible {
		t.Fatal("d was eligible with no prerequisites satisfied")
	}
	if got.Reason != "missing_prerequisites" {
		t.Errorf("reason = %q, want missing_prerequisites", got.Reason)
	}
	if len(got.Missing) != 2 {
		t.Errorf("missing = %v, want both b and c", got.Missing)
	}
}
