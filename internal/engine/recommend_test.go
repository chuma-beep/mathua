package engine

import (
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/scheduler"
	"github.com/chuma-beep/mathua/internal/storage"
)

// ReviewsDue must count exactly the concepts NextReview can serve. It used to
// count any row with a stale NextReviewDue, including LEARNING and PRACTICING
// rows the scheduler never puts in the review bucket. That made the profile
// advertise a review while /review answered "All caught up" (report 274).
func TestReviewsDue_CountsOnlyServableReviews(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("dana")
	past := time.Now().Add(-time.Hour)

	rows := []struct {
		name string
		prog storage.ConceptProgress
	}{
		{"decayed mastered, review due", storage.ConceptProgress{
			StudentID: st.ID, ConceptID: "a", Status: "MASTERED",
			NextReviewDue: &past, LastReviewed: &past, LastAttempted: &past,
		}},
		{"mastered, review not yet due", storage.ConceptProgress{
			StudentID: st.ID, ConceptID: "b", Status: "MASTERED",
			NextReviewDue: ptr(time.Now().Add(48 * time.Hour)),
		}},
		{"mastered, no review scheduled", storage.ConceptProgress{
			StudentID: st.ID, ConceptID: "c", Status: "MASTERED",
		}},
		{"learning, review due", storage.ConceptProgress{
			StudentID: st.ID, ConceptID: "d", Status: "LEARNING",
			NextReviewDue: &past, LastAttempted: &past,
		}},
	}
	for _, r := range rows {
		if err := e.repo.UpsertProgress(&r.prog); err != nil {
			t.Fatalf("%s: upsert: %v", r.name, err)
		}
	}

	n, err := e.ReviewsDue(st.ID)
	if err != nil {
		t.Fatalf("reviews due: %v", err)
	}
	// Only a. b is scheduled but not due, c was never scheduled, and d is
	// learning, which NextReview never serves as a retention review.
	if n != 1 {
		t.Errorf("ReviewsDue = %d, want 1 (only the decayed mastered concept)", n)
	}

	sess, err := e.repo.CreateSession(st.ID)
	if err != nil {
		t.Fatalf("create session: %v", err)
	}
	q, err := e.NextReviewQuestion(sess.ID, st.ID)
	if err != nil {
		t.Fatalf("next review question: %v", err)
	}
	if q == nil || q.ConceptID != "a" {
		t.Fatalf("NextReviewQuestion = %+v, want concept a", q)
	}
}

func TestReviewsDue_LearningConceptIsNotAReview(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("lee")
	past := time.Now().Add(-time.Hour)
	if err := e.repo.UpsertProgress(&storage.ConceptProgress{
		StudentID: st.ID, ConceptID: "a", Status: "LEARNING",
		NextReviewDue: &past, LastAttempted: &past, LastReviewed: &past,
	}); err != nil {
		t.Fatalf("upsert: %v", err)
	}

	n, err := e.ReviewsDue(st.ID)
	if err != nil {
		t.Fatalf("reviews due: %v", err)
	}
	if n != 0 {
		t.Errorf("ReviewsDue = %d, want 0 for a learning concept", n)
	}
	sess, err := e.repo.CreateSession(st.ID)
	if err != nil {
		t.Fatalf("create session: %v", err)
	}
	q, err := e.NextReviewQuestion(sess.ID, st.ID)
	if err != nil {
		t.Fatalf("next review question: %v", err)
	}
	if q != nil {
		t.Errorf("NextReviewQuestion served %s, want none", q.ConceptID)
	}
}

func TestReviewsDue_MatchesGraphDecayWhenPrerequisiteIsUnmet(t *testing.T) {
	// Reported from Discord: the graph counted 172 decayed concepts while the
	// review queue offered 18. A decayed concept whose prerequisite is no
	// longer met was counted by GetProgress and skipped by the scheduler.
	e := testEngine(t)
	st, _ := e.CreateStudent("riley")
	old := time.Now().Add(-30 * 24 * time.Hour)
	if err := e.repo.UpsertProgress(&storage.ConceptProgress{
		StudentID: st.ID, ConceptID: "b", Status: "MASTERED",
		LastReviewed: &old, LastAttempted: &old,
	}); err != nil {
		t.Fatalf("upsert: %v", err)
	}

	progress, err := e.GetProgress(st.ID)
	if err != nil {
		t.Fatalf("get progress: %v", err)
	}
	if got := progress["b"].Status; got != "DECAYING" {
		t.Fatalf("graph status for b = %s, want DECAYING", got)
	}

	n, err := e.ReviewsDue(st.ID)
	if err != nil {
		t.Fatalf("reviews due: %v", err)
	}
	if n != 1 {
		t.Fatalf("ReviewsDue = %d, want 1 to match the graph's decaying count", n)
	}

	sess, err := e.repo.CreateSession(st.ID)
	if err != nil {
		t.Fatalf("create session: %v", err)
	}
	q, err := e.NextReviewQuestion(sess.ID, st.ID)
	if err != nil {
		t.Fatalf("next review question: %v", err)
	}
	if q == nil || q.ConceptID != "b" {
		t.Fatalf("NextReviewQuestion = %+v, want the decayed concept b", q)
	}
}

func TestReviewsDue_MatchesGraphDecayWhenMasteryPredatesReviewClock(t *testing.T) {
	// The same Discord report, second cause: a concept mastered through study
	// has a mastery date but no review date. The graph used to call it decayed
	// (999 days) while the scheduler called it fresh (0 days). Both now start
	// the decay clock at mastery, so the graph and the queue agree.
	e := testEngine(t)
	st, _ := e.CreateStudent("sam")
	old := time.Now().Add(-30 * 24 * time.Hour)
	if err := e.repo.UpsertProgress(&storage.ConceptProgress{
		StudentID: st.ID, ConceptID: "b", Status: "MASTERED", MasteredAt: &old,
	}); err != nil {
		t.Fatalf("upsert: %v", err)
	}

	progress, err := e.GetProgress(st.ID)
	if err != nil {
		t.Fatalf("get progress: %v", err)
	}
	if got := progress["b"].Status; got != "DECAYING" {
		t.Fatalf("graph status for b = %s, want DECAYING", got)
	}

	n, err := e.ReviewsDue(st.ID)
	if err != nil {
		t.Fatalf("reviews due: %v", err)
	}
	if n != 1 {
		t.Fatalf("ReviewsDue = %d, want 1 to match the graph's decaying count", n)
	}

	sess, err := e.repo.CreateSession(st.ID)
	if err != nil {
		t.Fatalf("create session: %v", err)
	}
	q, err := e.NextReviewQuestion(sess.ID, st.ID)
	if err != nil {
		t.Fatalf("next review question: %v", err)
	}
	if q == nil || q.ConceptID != "b" {
		t.Fatalf("NextReviewQuestion = %+v, want the decayed concept b", q)
	}
}

func TestRecommendNext_DoesNotOfferAQuizWithNoCandidates(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("nora")
	now := time.Now()
	if err := e.repo.UpsertProgress(&storage.ConceptProgress{
		StudentID: st.ID, ConceptID: "a", Status: "MASTERED",
		Attempts: 3, LastReviewed: &now, LastAttempted: &now,
	}); err != nil {
		t.Fatal(err)
	}
	if err := e.repo.AwardXP(st.ID, QuizGateXP, "test", "seed-quiz-ready"); err != nil {
		t.Fatal(err)
	}

	resp, err := e.RecommendNext(st.ID, nil)
	if err != nil {
		t.Fatalf("recommend: %v", err)
	}
	items := append([]scheduler.Recommendation{}, resp.Alternatives...)
	if resp.Primary != nil {
		items = append(items, *resp.Primary)
	}
	for _, item := range items {
		if item.Action.Href == "/goals?quiz=1" {
			t.Fatalf("recommended a quiz with no eligible concept: %+v", item)
		}
	}
}

func ptr(t time.Time) *time.Time { return &t }

// ConceptsMastered counted every non-UNSEEN row, so a concept the learner had merely *started*
// was counted as mastered. Fixed to count effective MASTERED.
//
// This test does NOT discriminate the fix, and the comment is here so nobody later reads it as a
// guard: the only consumer is the new-learner gate `ConceptsMastered == 0 && AttemptedDays == 0`,
// and AttemptedDays already excluded anyone with history, so the miscount never changed a
// response. It is worth stating plainly rather than shipping the fix as though a test proved it.
func TestRecommendNext_StartedConceptsAreNotSentToTheDiagnostic(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("erin")
	now := time.Now()

	// Enough evidence to master `a`, left reviewed just now so decay does not apply.
	if err := e.repo.UpsertProgress(&storage.ConceptProgress{
		StudentID: st.ID, ConceptID: "a", Status: "MASTERED",
		AvgResponseTime: 1, LastReviewed: &now, LastAttempted: &now,
	}); err != nil {
		t.Fatal(err)
	}
	// Started, never mastered.
	if err := e.repo.UpsertProgress(&storage.ConceptProgress{
		StudentID: st.ID, ConceptID: "b", Status: "PRACTICING",
		AvgResponseTime: 1, LastReviewed: &now, LastAttempted: &now,
	}); err != nil {
		t.Fatal(err)
	}

	resp, err := e.RecommendNext(st.ID, nil)
	if err != nil {
		t.Fatalf("recommend: %v", err)
	}
	if resp.Primary == nil {
		t.Fatal("no primary recommendation")
	}
	if resp.Primary.Action.Href == "/onboard" {
		t.Error("a learner with attempts and one mastered concept was sent to the diagnostic")
	}
}
