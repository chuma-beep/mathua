package storage

import (
	"fmt"
	"testing"
	"time"
)

func timedQuizStore(t *testing.T, repo Repository) TimedQuizStore {
	t.Helper()
	store, ok := repo.(TimedQuizStore)
	if !ok {
		t.Fatalf("repository %T does not implement TimedQuizStore", repo)
	}
	return store
}

func TestTimedQuizBaselineCountsOnlyPositiveXPAfterIt(t *testing.T) {
	store := newTestStore(t)
	st := seedStudent(t, store)
	tq := timedQuizStore(t, store)

	before := time.Now().UTC().Add(-time.Minute)
	baseline, err := tq.EnsureTimedQuizBaseline(st.ID)
	if err != nil {
		t.Fatalf("ensure baseline: %v", err)
	}
	if baseline.BaselineXP != 0 {
		t.Fatalf("baseline xp = %d, want 0 for a newly created student", baseline.BaselineXP)
	}
	if baseline.BaselineAt.Before(before) {
		t.Fatalf("baseline time = %v, want at or after %v", baseline.BaselineAt, before)
	}
	if err := store.AwardXP(st.ID, 40, "test", "positive-40"); err != nil {
		t.Fatalf("award positive xp: %v", err)
	}
	if err := store.AwardXP(st.ID, -15, "test", "negative-15"); err != nil {
		t.Fatalf("award negative xp: %v", err)
	}
	got, err := tq.TimedQuizEligibleXP(st.ID, baseline.BaselineAt)
	if err != nil {
		t.Fatalf("eligible xp: %v", err)
	}
	if got != 40 {
		t.Fatalf("eligible xp = %d, want 40; a penalty must not reduce earned XP", got)
	}
}

func TestTimedQuizItemPersistsAdjustmentForReplay(t *testing.T) {
	store := newTestStore(t)
	st := seedStudent(t, store)
	tq := timedQuizStore(t, store)
	now := time.Now().UTC()
	run := TimedQuizRun{
		ID:         fmt.Sprintf("tq-adjust-%d", now.UnixNano()),
		StudentID:  st.ID,
		StartedAt:  now,
		DeadlineAt: now.Add(10 * time.Minute),
		Status:     "active",
		CreatedAt:  now,
	}
	items := []TimedQuizItem{
		{ID: run.ID + ":0", RunID: run.ID, ConceptID: "a", Position: 0, Question: "1+1", ExpectedAnswer: "2", ServedAt: now, Status: "served"},
		{ID: run.ID + ":1", RunID: run.ID, ConceptID: "b", Position: 1, Question: "2+2", ExpectedAnswer: "4", ServedAt: now, Status: "served"},
	}
	if err := tq.CreateTimedQuizRun(run, items); err != nil {
		t.Fatalf("create run: %v", err)
	}
	if claimed, err := tq.ClaimTimedQuizItem(items[0].ID, run.ID, "sub-1"); err != nil || !claimed {
		t.Fatalf("claim first = %v, %v", claimed, err)
	}
	if updated, err := tq.UpdateTimedQuizItemWithAdjustment(items[0].ID, run.ID, "2", "sub-1", true, 4, 3, "answered", -2, "-2 XP: poor task performance", "extremely_poor_performance", "performance"); err != nil || !updated {
		t.Fatalf("update first = %v, %v", updated, err)
	}
	if claimed, err := tq.ClaimTimedQuizItem(items[1].ID, run.ID, "sub-2"); err != nil || !claimed {
		t.Fatalf("claim second = %v, %v", claimed, err)
	}
	if updated, err := tq.UpdateTimedQuizItemWithAdjustment(items[1].ID, run.ID, "4", "sub-2", true, 4, 3, "answered", 0, "0 XP: poor task performance", "poor_performance", "performance"); err != nil || !updated {
		t.Fatalf("update second = %v, %v", updated, err)
	}

	fresh, err := tq.GetTimedQuizItems(run.ID)
	if err != nil {
		t.Fatalf("get items: %v", err)
	}
	if len(fresh) != 2 {
		t.Fatalf("items = %d, want 2", len(fresh))
	}
	if fresh[0].XPAdjustment != -2 || fresh[0].XPNote == "" || fresh[0].XPReasonCode != "extremely_poor_performance" || fresh[0].XPSource != "performance" {
		t.Fatalf("first item adjustment = %+v", fresh[0])
	}
	if fresh[1].XPAdjustment != 0 || fresh[1].XPNote == "" || fresh[1].XPReasonCode != "poor_performance" {
		t.Fatalf("zero-adjustment item = %+v", fresh[1])
	}

	// A duplicate update must not rewrite or double-charge the durable replay
	// fields; the item is already answered.
	if updated, err := tq.UpdateTimedQuizItemWithAdjustment(items[0].ID, run.ID, "2", "sub-1", true, 4, 3, "answered", -9, "changed", "changed", "changed"); err != nil || updated {
		t.Fatalf("duplicate update = %v, %v; want false, nil", updated, err)
	}
	again, err := tq.GetTimedQuizItems(run.ID)
	if err != nil {
		t.Fatalf("get items again: %v", err)
	}
	if again[0].XPAdjustment != -2 || again[0].XPNote == "" {
		t.Fatalf("duplicate update rewrote replay fields: %+v", again[0])
	}
}

func TestTimedQuizRunLifecycleIsIdempotent(t *testing.T) {
	store := newTestStore(t)
	st := seedStudent(t, store)
	tq := timedQuizStore(t, store)
	now := time.Now().UTC()
	run := TimedQuizRun{
		ID:         fmt.Sprintf("tq-%d", now.UnixNano()),
		StudentID:  st.ID,
		StartedAt:  now,
		DeadlineAt: now.Add(10 * time.Minute),
		Status:     "active",
		CreatedAt:  now,
	}
	items := []TimedQuizItem{
		{ID: run.ID + ":0", RunID: run.ID, ConceptID: "a", Position: 0, Question: "1+1", ExpectedAnswer: "2", ServedAt: now, Status: "served"},
		{ID: run.ID + ":1", RunID: run.ID, ConceptID: "b", Position: 1, Question: "2+2", ExpectedAnswer: "4", ServedAt: now, Status: "served"},
	}
	if err := tq.CreateTimedQuizRun(run, items); err != nil {
		t.Fatalf("create run: %v", err)
	}
	if active, err := tq.GetActiveTimedQuizRun(st.ID); err != nil || active == nil || active.ID != run.ID {
		t.Fatalf("active run = %+v, err = %v", active, err)
	}
	// A second active run for the same learner must be rejected by the partial
	// unique index, not silently coexist.
	if err := tq.CreateTimedQuizRun(run, items); err == nil {
		t.Fatal("second active run was created; want a uniqueness error")
	}
	if claimed, err := tq.ClaimTimedQuizItem(items[0].ID, run.ID, "sub-1"); err != nil || !claimed {
		t.Fatalf("claim item: claimed=%v err=%v", claimed, err)
	}
	if claimed, err := tq.ClaimTimedQuizItem(items[0].ID, run.ID, "sub-2"); err != nil || claimed {
		t.Fatalf("second claim: claimed=%v err=%v", claimed, err)
	}
	updated, err := tq.UpdateTimedQuizItem(items[0].ID, run.ID, "2", "sub-1", true, 4.5, 5, "answered")
	if err != nil || !updated {
		t.Fatalf("update item: updated=%v err=%v", updated, err)
	}
	if updated, err := tq.UpdateTimedQuizItem(items[1].ID, run.ID, "4", "sub-2", true, 3.0, 5, "answered"); err != nil || updated {
		t.Fatalf("update without a claim: updated=%v err=%v; want false and no error", updated, err)
	}
	if claimed, err := tq.ClaimTimedQuizItem(items[1].ID, run.ID, "sub-2"); err != nil || !claimed {
		t.Fatalf("claim second item: claimed=%v err=%v", claimed, err)
	}
	if updated, err := tq.UpdateTimedQuizItem(items[1].ID, run.ID, "4", "sub-2", true, 3.0, 5, "answered"); err != nil || !updated {
		t.Fatalf("update second item: updated=%v err=%v", updated, err)
	}
	first, err := tq.FinalizeTimedQuizRun(run.ID, st.ID, now)
	if err != nil {
		t.Fatalf("finalize: %v", err)
	}
	if first.CorrectCount != 2 || first.UnansweredCount != 0 || first.XPEarned != 10 {
		t.Fatalf("finalized run = %+v, want 2 correct, 0 unanswered, 10 XP", first)
	}
	second, err := tq.FinalizeTimedQuizRun(run.ID, st.ID, now.Add(time.Minute))
	if err != nil {
		t.Fatalf("finalize again: %v", err)
	}
	if second.CorrectCount != first.CorrectCount || second.CompletedAt == nil || !second.CompletedAt.Equal(*first.CompletedAt) {
		t.Fatalf("finalize was not idempotent: first=%+v second=%+v", first, second)
	}
}

func TestReviewEventIsUniquePerAttempt(t *testing.T) {
	store := newTestStore(t)
	st := seedStudent(t, store)
	tq := timedQuizStore(t, store)
	event := ReviewEvent{
		StudentID:      st.ID,
		ConceptID:      "a",
		AttemptID:      "sess:attempt",
		Rating:         3,
		Correct:        true,
		ElapsedSeconds: 4,
		ReviewedAt:     time.Now().UTC(),
	}
	for i := 0; i < 3; i++ {
		if err := tq.RecordReviewEvent(event); err != nil {
			t.Fatalf("record review event %d: %v", i, err)
		}
	}
	if count := seedInt(t, store, `SELECT COUNT(*) FROM review_events WHERE attempt_id = ?`, event.AttemptID); count != 1 {
		t.Fatalf("review event rows = %d, want 1", count)
	}
}

func TestTimedQuizRowsAreRemovedByResetAndDelete(t *testing.T) {
	store := newTestStore(t)
	st := seedStudent(t, store)
	tq := timedQuizStore(t, store)
	now := time.Now().UTC()
	run := TimedQuizRun{ID: fmt.Sprintf("tq-reset-%d", now.UnixNano()), StudentID: st.ID, StartedAt: now, DeadlineAt: now.Add(time.Minute), Status: "active", CreatedAt: now}
	items := []TimedQuizItem{{ID: run.ID + ":0", RunID: run.ID, ConceptID: "a", Position: 0, Question: "q", ExpectedAnswer: "a", ServedAt: now, Status: "served"}}
	if err := tq.CreateTimedQuizRun(run, items); err != nil {
		t.Fatalf("create run: %v", err)
	}
	if err := tq.RecordReviewEvent(ReviewEvent{StudentID: st.ID, ConceptID: "a", AttemptID: run.ID, Rating: 3, Correct: true, ReviewedAt: now}); err != nil {
		t.Fatalf("record review event: %v", err)
	}
	if err := store.ResetProgress(st.ID); err != nil {
		t.Fatalf("reset progress: %v", err)
	}
	if n := seedInt(t, store, `SELECT COUNT(*) FROM timed_quiz_runs WHERE student_id = ?`, st.ID); n != 0 {
		t.Fatalf("timed quiz runs after reset = %d, want 0", n)
	}
	if n := seedInt(t, store, `SELECT COUNT(*) FROM timed_quiz_items WHERE run_id = ?`, run.ID); n != 0 {
		t.Fatalf("timed quiz items after reset = %d, want 0", n)
	}
	if n := seedInt(t, store, `SELECT COUNT(*) FROM timed_quiz_baselines WHERE student_id = ?`, st.ID); n != 0 {
		t.Fatalf("timed quiz baselines after reset = %d, want 0", n)
	}
	if n := seedInt(t, store, `SELECT COUNT(*) FROM review_events WHERE student_id = ?`, st.ID); n != 0 {
		t.Fatalf("review events after reset = %d, want 0", n)
	}
	// Seed the timed-quiz rows again after the reset so the deletion path is
	// exercised too, not just the reset path.
	now = now.Add(time.Second)
	run = TimedQuizRun{ID: fmt.Sprintf("tq-delete-%d", now.UnixNano()), StudentID: st.ID, StartedAt: now, DeadlineAt: now.Add(time.Minute), Status: "active", CreatedAt: now}
	items = []TimedQuizItem{{ID: run.ID + ":0", RunID: run.ID, ConceptID: "a", Position: 0, Question: "q", ExpectedAnswer: "a", ServedAt: now, Status: "served"}}
	if err := tq.CreateTimedQuizRun(run, items); err != nil {
		t.Fatalf("recreate run: %v", err)
	}
	if err := tq.RecordReviewEvent(ReviewEvent{StudentID: st.ID, ConceptID: "a", AttemptID: run.ID, Rating: 3, Correct: true, ReviewedAt: now}); err != nil {
		t.Fatalf("recreate review event: %v", err)
	}
	if err := store.DeleteAccount(st.ID); err != nil {
		t.Fatalf("delete account: %v", err)
	}
	if n := seedInt(t, store, `SELECT COUNT(*) FROM timed_quiz_runs WHERE student_id = ?`, st.ID); n != 0 {
		t.Fatalf("timed quiz runs after delete = %d, want 0", n)
	}
	if n := seedInt(t, store, `SELECT COUNT(*) FROM timed_quiz_items WHERE run_id = ?`, run.ID); n != 0 {
		t.Fatalf("timed quiz items after delete = %d, want 0", n)
	}
	if n := seedInt(t, store, `SELECT COUNT(*) FROM timed_quiz_baselines WHERE student_id = ?`, st.ID); n != 0 {
		t.Fatalf("timed quiz baselines after delete = %d, want 0", n)
	}
	if n := seedInt(t, store, `SELECT COUNT(*) FROM review_events WHERE student_id = ?`, st.ID); n != 0 {
		t.Fatalf("review events after delete = %d, want 0", n)
	}
}
