package engine

import (
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/storage"
)

// The scheduler enforces a per-concept daily cap, but a cap no caller supplies is decoration.
// The engine has to count what was already served today, or a learner on a high-volume plan
// keeps meeting the same concept until that concept's supply runs out inside one sitting.
//
// The test is differential rather than absolute: the same learner is recommended before and
// after the concept reaches its cap, so a ranking change would have to be caused by the
// attempts rather than by some unrelated difference in setup.
func TestRecommendNext_DropsConceptThatReachedItsDailyCap(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("dana")
	if err := e.repo.SetDiagnosticCompleted(st.ID); err != nil {
		t.Fatalf("mark placement done: %v", err)
	}
	now := time.Now()
	past := now.Add(-24 * time.Hour)

	prog := storage.ConceptProgress{
		StudentID: st.ID, ConceptID: "a", Status: "LEARNING",
		LastAttempted: &past,
	}
	if err := e.repo.UpsertProgress(&prog); err != nil {
		t.Fatalf("upsert: %v", err)
	}

	before, err := e.RecommendNext(st.ID, nil)
	if err != nil {
		t.Fatalf("recommend before: %v", err)
	}
	if before.Primary == nil || before.Primary.ConceptID != "a" {
		t.Fatalf("setup is wrong: before any attempt today the concept is not offered")
	}

	sess, err := e.repo.CreateSession(st.ID)
	if err != nil {
		t.Fatalf("create session: %v", err)
	}
	for i := 0; i < servedPerDayCap; i++ {
		if _, err := e.repo.RecordAttempt(storage.AttemptEntry{
			SessionID: sess.ID, StudentID: st.ID, ConceptID: "a", Timestamp: now,
			Answer: "1", Expected: "1", Correct: true,
		}); err != nil {
			t.Fatalf("record attempt %d: %v", i, err)
		}
	}

	after, err := e.RecommendNext(st.ID, nil)
	if err != nil {
		t.Fatalf("recommend after: %v", err)
	}
	if after.Primary != nil && after.Primary.ConceptID == "a" {
		t.Errorf("concept a reached its daily cap of %d and must not be the primary card", servedPerDayCap)
	}
	for _, r := range after.Alternatives {
		if r.ConceptID == "a" {
			t.Errorf("concept a reached its daily cap of %d and must not appear among the alternatives", servedPerDayCap)
		}
	}
}
