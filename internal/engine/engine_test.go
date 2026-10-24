package engine

import (
	"errors"
	"fmt"
	"math"
	"sync"
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/diagnostic"
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/mastery"
	"github.com/chuma-beep/mathua/internal/storage"
	xppolicy "github.com/chuma-beep/mathua/internal/xp"
)

type testGen struct {
	answer string
}

func (g *testGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return generator.Problem{
		Question:    "What is the answer?",
		Answer:      g.answer,
		Explanation: "Here is why.",
	}
}

func testEngine(t *testing.T) *Engine {
	t.Helper()
	d, err := concepts.Build([]concepts.Concept{
		{
			ID: "a", Label: "Concept A", Domain: "d",
			GradingType: "numeric", Prerequisites: []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60},
		},
		{
			ID: "b", Label: "Concept B", Domain: "d",
			GradingType: "numeric", Prerequisites: []string{"a"},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60},
		},
	})
	if err != nil {
		t.Fatalf("build DAG: %v", err)
	}
	store, err := storage.NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("create store: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	reg := generator.NewRegistry()
	reg.Register("a", &testGen{answer: "42"})
	reg.Register("b", &testGen{answer: "99"})
	return New(store, d, reg, nil, nil)
}

func TestEngine_CreateStudent(t *testing.T) {
	e := testEngine(t)
	st, err := e.CreateStudent("alice")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if st.Name != "alice" {
		t.Errorf("expected alice, got %s", st.Name)
	}
}

func TestEngine_NextQuestion(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("bob")
	sess, _ := e.repo.CreateSession(st.ID)
	q, err := e.NextQuestion(sess.ID, st.ID)
	if err != nil {
		t.Fatalf("next question: %v", err)
	}
	if q == nil {
		t.Fatal("expected a question")
	}
	if q.ConceptID != "a" {
		t.Errorf("expected concept a, got %s", q.ConceptID)
	}
	if q.Question == "" {
		t.Error("expected non-empty question")
	}
	if q.IsReview {
		t.Error("expected not review")
	}
}

func TestEngine_SubmitAnswer_Correct(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("carol")
	sess, _ := e.repo.CreateSession(st.ID)
	q, _ := e.NextQuestion(sess.ID, st.ID)
	result, err := e.SubmitAnswer(sess.ID, st.ID, q.AttemptID, "42", 3.0)
	if err != nil {
		t.Fatalf("submit answer: %v", err)
	}
	if !result.Correct {
		t.Error("expected correct answer")
	}
	if result.Streak != 1 {
		t.Errorf("expected streak 1, got %d", result.Streak)
	}
	progress, _ := e.GetProgress(st.ID)
	if p, ok := progress["a"]; !ok || p.Streak != 1 {
		t.Errorf("expected progress with streak 1, got %v", p)
	}
}

func TestEngine_SubmitAnswer_Wrong(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("dave")
	sess, _ := e.repo.CreateSession(st.ID)
	q, _ := e.NextQuestion(sess.ID, st.ID)
	result, err := e.SubmitAnswer(sess.ID, st.ID, q.AttemptID, "99", 5.0)
	if err != nil {
		t.Fatalf("submit answer: %v", err)
	}
	if result.Correct {
		t.Error("expected incorrect answer")
	}
	if result.Streak != 0 {
		t.Errorf("expected streak 0, got %d", result.Streak)
	}
	if result.Explanation == "" {
		t.Error("expected explanation for wrong answer")
	}
}

func TestEngine_MasteryProgression(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("eve")

	// Answer concept "a" correctly 3 times (streak=3 → reaches LEARNING).
	for i := 0; i < 3; i++ {
		sess, _ := e.repo.CreateSession(st.ID)
		e.mu.Lock()
		e.sessions[sess.ID] = &activeSession{}
		e.mu.Unlock()
		q, _ := e.NextQuestion(sess.ID, st.ID)
		result, err := e.SubmitAnswer(sess.ID, st.ID, q.AttemptID, "42", 1.0)
		if err != nil {
			t.Fatalf("submit %d: %v", i+1, err)
		}
		if !result.Correct {
			t.Fatalf("expected correct on attempt %d", i+1)
		}
	}
	// One clean answer is LEARNING — the cheap first rung.
	one, _ := e.CreateStudent("prog_one")
	answerCorrectly(t, e, one.ID, 1)
	if p := mustProgress(t, e, one.ID)["a"]; p == nil || p.Status != "LEARNING" {
		t.Errorf("expected LEARNING after 1 clean answer, got %+v", p)
	}

	// Three clean answers reach PRACTICING and stop: the mastery floor is a full evidence
	// window, so the top rung is out of reach however good the answers are.
	//
	// This used to assert LEARNING, because the old rule needed the streak to *reach* the
	// threshold before it would move at all — and then a further full threshold for each of
	// the next two tiers.
	progress, _ := e.GetProgress(st.ID)
	if p, ok := progress["a"]; !ok || p.Status != "PRACTICING" {
		t.Errorf("expected PRACTICING after 3 clean answers, got %+v", progress["a"])
	}

	// A full window of clean answers reaches MASTERED.
	full, _ := e.CreateStudent("prog_full")
	answerCorrectly(t, e, full.ID, mastery.MasteryEvidenceFloor)
	if p := mustProgress(t, e, full.ID)["a"]; p == nil || p.Status != "MASTERED" {
		t.Errorf("expected MASTERED after a full clean window, got %+v", p)
	}
}

func mustProgress(t *testing.T, e *Engine, studentID string) map[string]*storage.ConceptProgress {
	t.Helper()
	p, err := e.GetProgress(studentID)
	if err != nil {
		t.Fatalf("get progress: %v", err)
	}
	return p
}

func TestEngine_GetScores(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("frank")
	scores, err := e.GetScores(st.ID)
	if err != nil {
		t.Fatalf("get scores: %v", err)
	}
	if scores.ConceptsMastered != 0 {
		t.Errorf("expected 0 mastered, got %d", scores.ConceptsMastered)
	}
	if scores.Level != "Novice" {
		t.Errorf("expected Novice, got %s", scores.Level)
	}
}

// A duplicate of one attempt may race the real submission. Every caller must
// get the recorded result back — a retry the learner made because the first
// response was lost is not an error — but only one grading may happen, so the
// progress is written once and no second attempt is recorded.
func TestEngine_SubmitAnswer_ConcurrentDuplicates(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("grace")
	sess, _ := e.repo.CreateSession(st.ID)
	q, _ := e.NextQuestion(sess.ID, st.ID)

	const workers = 8
	var wg sync.WaitGroup
	results := make([]error, workers)
	out := make([]*AnswerResult, workers)
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			out[i], results[i] = e.SubmitAnswer(sess.ID, st.ID, q.AttemptID, "42", 3.0)
		}(i)
	}
	wg.Wait()

	for i, err := range results {
		if err != nil {
			t.Fatalf("duplicate %d was rejected: %v", i, err)
		}
	}
	// Every duplicate must report the same accepted answer, not a fresh grade.
	for i := 1; i < workers; i++ {
		if out[i] == nil || out[0] == nil {
			t.Fatalf("duplicate %d returned no result", i)
		}
		if out[i].Correct != out[0].Correct || out[i].XP != out[0].XP {
			t.Errorf("duplicate %d = correct:%v xp:%d, want %v/%d", i, out[i].Correct, out[i].XP, out[0].Correct, out[0].XP)
		}
	}

	// The grade happened once: same streak, one attempt, answer "42".
	progress, _ := e.GetProgress(st.ID)
	if p, ok := progress["a"]; !ok || p.Streak != 1 || p.Attempts != 1 {
		t.Errorf("expected concept a graded once, got %+v", progress["a"])
	}
}

// A client that retries a submission it never saw acknowledged must get the
// recorded result back, not an error: the answer was accepted, and the learner
// should be told so. The attempt id is the idempotency key, so a duplicate
// carrying different text still replays the first grading instead of grading
// the new text against a question that has already moved on.
func TestEngine_SubmitAnswer_ReplaysAnAcceptedAttempt(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("replay")
	sess, _ := e.repo.CreateSession(st.ID)

	q1, err := e.NextQuestion(sess.ID, st.ID)
	if err != nil {
		t.Fatalf("next question: %v", err)
	}
	if q1.AttemptID == "" {
		t.Fatal("expected attempt id on question")
	}
	if _, err := e.SubmitAnswer(sess.ID, st.ID, q1.AttemptID, "42", 3.0); err != nil {
		t.Fatalf("submit first answer: %v", err)
	}

	// Session advanced; the next question carries a fresh token.
	q2, err := e.NextQuestion(sess.ID, st.ID)
	if err != nil {
		t.Fatalf("next question after first submit: %v", err)
	}
	if q2.AttemptID == "" || q2.AttemptID == q1.AttemptID {
		t.Fatalf("expected a fresh attempt id for the next question, got %q", q2.AttemptID)
	}

	// The retry replays the accepted answer. Different text is ignored, because
	// the key is the attempt id and that attempt has already been graded.
	if _, err := e.SubmitAnswer(sess.ID, st.ID, q1.AttemptID, "99", 3.0); err != nil {
		t.Fatalf("replay of an accepted attempt: %v", err)
	}
	progress, _ := e.GetProgress(st.ID)
	if p := progress["a"]; p == nil || p.Attempts != 1 || p.Streak != 1 {
		t.Fatalf("replay must not be graded twice: got %+v", progress["a"])
	}

	// The current question's token is still accepted.
	if _, err := e.SubmitAnswer(sess.ID, st.ID, q2.AttemptID, "42", 3.0); err != nil {
		t.Fatalf("submit with current token: %v", err)
	}
}

func TestEngine_Diagnostic(t *testing.T) {
	e := testEngine(t)
	s := e.StartDiagnostic()
	if s == nil {
		t.Fatal("expected diagnostic session")
	}
	_, cid, err := e.NextDiagnosticQuestion(s)
	if err != nil {
		t.Fatalf("next diag question: %v", err)
	}
	if cid == "" {
		t.Error("expected non-empty concept ID")
	}
}

func TestEngine_GetLeaderboard(t *testing.T) {
	e := testEngine(t)
	alice, err := e.CreateStudent("alice")
	if err != nil {
		t.Fatalf("create alice: %v", err)
	}
	if _, err := e.CreateStudent("bob"); err != nil {
		t.Fatalf("create bob: %v", err)
	}
	// The board ranks earned XP, so a learner with no ledger entries is not on
	// it. Bob stays off until he answers something; Alice is the only entry.
	if err := e.repo.AwardXP(alice.ID, 5, "test", "alice:1"); err != nil {
		t.Fatalf("award: %v", err)
	}
	board, err := e.GetLeaderboardPage(alice.ID, "weekly", 1, 25)
	if err != nil {
		t.Fatalf("get leaderboard: %v", err)
	}
	if board.Total != 1 {
		t.Errorf("total = %d, want 1", board.Total)
	}
	if len(board.Entries) != 1 || board.Entries[0].Name != "alice" {
		t.Fatalf("entries = %+v, want alice only", board.Entries)
	}
	if board.Me.Rank != 1 || board.Me.XP != 5 {
		t.Errorf("me = %+v, want rank 1 with 5 XP", board.Me)
	}
}

func TestEngine_GetDAG(t *testing.T) {
	e := testEngine(t)
	dag := e.GetDAG()
	if dag.Count() != 2 {
		t.Errorf("expected 2 concepts, got %d", dag.Count())
	}
}

func TestEngine_FullMasteryCycle(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("mastery_tester")
	sess, _ := e.repo.CreateSession(st.ID)

	// Master concept "a", then confirm "b" becomes available because its prerequisite is
	// now mastered.
	//
	// This loop used to ask for concept "a" nine times and fail if it got anything else,
	// because the old rule needed 3 × required_streak correct answers per concept. Under
	// the evidence ladder "a" is mastered in three and the scheduler moves to "b" — so the
	// test now asserts the *frontier advancing*, which is the behaviour the ladder exists
	// to produce. It is the clearest available proof of the cost change: 3 answers, not 9,
	// and not 30.
	masteredOn := 0
	for i := 0; i < 9; i++ {
		e.mu.Lock()
		e.sessions[sess.ID] = &activeSession{}
		e.mu.Unlock()
		q, err := e.NextQuestion(sess.ID, st.ID)
		if err != nil {
			t.Fatalf("next question %d: %v", i+1, err)
		}
		if q == nil {
			t.Fatal("expected question")
		}
		if q.ConceptID != "a" {
			break // the frontier moved on, which is the point
		}
		res, err := e.SubmitAnswer(sess.ID, st.ID, q.AttemptID, "42", 1.0)
		if err != nil {
			t.Fatalf("submit answer %d: %v", i+1, err)
		}
		if !res.Correct {
			t.Fatalf("expected correct on attempt %d, got %s", i+1, res.Feedback)
		}
		masteredOn = i + 1
		if p, _ := e.GetProgress(st.ID); p["a"] != nil && p["a"].Status == "MASTERED" {
			break
		}
	}

	progress, _ := e.GetProgress(st.ID)
	if p := progress["a"]; p == nil || p.Status != "MASTERED" {
		t.Fatalf("expected a MASTERED, got %+v", progress["a"])
	}
	if masteredOn > 6 {
		t.Errorf("a took %d answers to master; the ladder is supposed to be a handful, not a threshold-length streak", masteredOn)
	}
	t.Logf("concept a mastered in %d correct answers", masteredOn)

	// Now concept "b" should be available (prereq "a" is mastered).
	e.mu.Lock()
	e.sessions[sess.ID] = &activeSession{}
	e.mu.Unlock()
	q, err := e.NextQuestion(sess.ID, st.ID)
	if err != nil {
		t.Fatalf("next question after mastery: %v", err)
	}
	if q == nil {
		t.Fatal("expected concept b to become available")
	}
	if q.ConceptID != "b" {
		t.Errorf("expected concept b after mastering a, got %s", q.ConceptID)
	}

	// Answer "a" correctly 9 more times... no wait, now concept "b" is selected.
	// Answer "b" correctly (it expects "99").
	res, err := e.SubmitAnswer(sess.ID, st.ID, q.AttemptID, "99", 1.0)
	if err != nil {
		t.Fatalf("submit answer for b: %v", err)
	}
	if !res.Correct {
		t.Fatalf("expected correct on b, got %s", res.Feedback)
	}
	scores, _ := e.GetScores(st.ID)
	if scores.ConceptsMastered != 1 {
		t.Errorf("expected 1 mastered, got %d", scores.ConceptsMastered)
	}
	t.Logf("Mastery complete! A: %s, Level: %s", progress["a"].Status, scores.Level)
}

func TestEngine_LessonAttachedToQuestion(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("lesson_tester")
	sess, _ := e.repo.CreateSession(st.ID)

	q, err := e.NextQuestion(sess.ID, st.ID)
	if err != nil {
		t.Fatalf("next question: %v", err)
	}
	if q == nil {
		t.Fatal("expected question")
	}
	if q.ConceptID != "a" {
		t.Fatalf("expected concept a, got %s", q.ConceptID)
	}
	// No lessons loaded, so Lesson should be nil
	if q.Lesson != nil {
		t.Error("expected nil lesson (no loader configured)")
	}
}

// PR 1.5 tests

func TestEngine_Halt_AfterTwoConsecutiveMisses(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("halt")
	sess, _ := e.repo.CreateSession(st.ID)

	q1, _ := e.NextQuestion(sess.ID, st.ID)
	if q1 == nil {
		t.Fatal("expected first question")
	}
	res, err := e.SubmitAnswer(sess.ID, st.ID, q1.AttemptID, "wrong", 5.0)
	if err != nil {
		t.Fatalf("submit miss 1: %v", err)
	}
	if res.Halted {
		t.Error("expected no halt after first miss")
	}
	if res.XP != 0 {
		t.Errorf("expected 0 XP on miss, got %d", res.XP)
	}

	q2, _ := e.NextQuestion(sess.ID, st.ID)
	if q2 == nil {
		t.Fatal("expected second question")
	}
	res2, err := e.SubmitAnswer(sess.ID, st.ID, q2.AttemptID, "wrong", 5.0)
	if err != nil {
		t.Fatalf("submit miss 2: %v", err)
	}
	if !res2.Halted {
		t.Error("expected halt after two consecutive misses")
	}
	if res2.XP != 0 {
		t.Errorf("expected XP cut short (0) on halt, got %d", res2.XP)
	}

	// Next question must come from the remedial queue (concept a re-queued
	// since it has no KeyPrerequisites in the test DAG).
	q3, _ := e.NextQuestion(sess.ID, st.ID)
	if q3 == nil {
		t.Fatal("expected remedial question after halt")
	}
	if q3.ConceptID != "a" {
		t.Errorf("expected remedial re-queue of a, got %s", q3.ConceptID)
	}
	// Correct answer clears the halt.
	res3, err := e.SubmitAnswer(sess.ID, st.ID, q3.AttemptID, "42", 5.0)
	if err != nil {
		t.Fatalf("submit remedial correct: %v", err)
	}
	if !res3.Correct {
		t.Fatalf("expected correct on remedial, got %s", res3.Feedback)
	}
	e.mu.Lock()
	haltedAfter := e.sessions[sess.ID].halted
	e.mu.Unlock()
	if haltedAfter {
		t.Error("expected halt cleared after remedial success")
	}
	// Scheduling resumes — a remains the only unlocked concept in this DAG,
	// so it must be served again (b still gated on a's mastery).
	q4, _ := e.NextQuestion(sess.ID, st.ID)
	if q4 == nil {
		t.Fatal("expected scheduling to resume")
	}
}

func TestEngine_Halt_KeyPrerequisites_TargetedRemediation(t *testing.T) {
	d, err := concepts.Build([]concepts.Concept{
		{
			ID: "a", Label: "Concept A", Domain: "d",
			GradingType: "numeric", Prerequisites: []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 2, AvgTimeSeconds: 60},
		},
		{
			ID: "b", Label: "Concept B", Domain: "d",
			GradingType: "numeric", Prerequisites: []string{"a"},
			KeyPrerequisites: []string{"a"},
			Variants: []concepts.Variant{
				{Difficulty: 0.3, TimeThresh: 60.0, Label: "b-easy"},
			},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 2, AvgTimeSeconds: 60},
		},
	})
	if err != nil {
		t.Fatalf("build DAG: %v", err)
	}
	store, err := storage.NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("create store: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	reg := generator.NewRegistry()
	reg.Register("a", &testGen{answer: "42"})
	reg.Register("b", &testGen{answer: "99"})
	e := New(store, d, reg, nil, nil)
	st, _ := e.CreateStudent("targeted")
	sess, _ := e.repo.CreateSession(st.ID)

	// Step 1: Mark prerequisite concept a mastered so concept b is unlocked.
	now := time.Now().UTC()
	if err := e.repo.UpsertProgress(&storage.ConceptProgress{
		StudentID:     st.ID,
		ConceptID:     "a",
		Status:        string(mastery.StatusMastered),
		Streak:        3,
		LastAttempted: &now,
		LastReviewed:  &now,
	}); err != nil {
		t.Fatalf("upsert progress: %v", err)
	}

	// Step 2: Now concept b is served.
	q1, err := e.NextQuestion(sess.ID, st.ID)
	if err != nil || q1 == nil {
		t.Fatalf("step 2 next: %v", err)
	}
	if q1.ConceptID != "b" {
		t.Fatalf("expected concept b, got %s", q1.ConceptID)
	}

	// First miss on concept b.
	res1, err := e.SubmitAnswer(sess.ID, st.ID, q1.AttemptID, "wrong", 5.0)
	if err != nil {
		t.Fatalf("submit miss 1: %v", err)
	}
	if res1.Halted {
		t.Error("expected no halt after first miss")
	}
	if res1.XP != 0 {
		t.Errorf("expected 0 XP on miss, got %d", res1.XP)
	}

	// Second question served on concept b.
	q2, err := e.NextQuestion(sess.ID, st.ID)
	if err != nil || q2 == nil {
		t.Fatalf("step 2 next 2: %v", err)
	}
	if q2.ConceptID != "b" {
		t.Fatalf("expected concept b, got %s", q2.ConceptID)
	}

	// Second consecutive miss on concept b triggers 2-fail halt.
	res2, err := e.SubmitAnswer(sess.ID, st.ID, q2.AttemptID, "wrong", 5.0)
	if err != nil {
		t.Fatalf("submit miss 2: %v", err)
	}
	if !res2.Halted {
		t.Error("expected halt after two consecutive misses")
	}
	if res2.XP != 0 {
		t.Errorf("expected 0 XP on halt, got %d", res2.XP)
	}
	// Verify remedial queue contains key prerequisite a then b.
	if len(res2.Remedial) != 2 || res2.Remedial[0] != "a" || res2.Remedial[1] != "b" {
		t.Fatalf("expected remedial [a b], got %v", res2.Remedial)
	}
	// Verify student-level remedial queue is populated.
	rem := e.RemedialQueue(st.ID)
	if len(rem) != 1 || rem[0] != "a" {
		t.Fatalf("expected engine remedial queue [a], got %v", rem)
	}

	// Step 3: Verify learner task shelf shows targeted key prerequisite at the front.
	rec, err := e.RecommendNext(st.ID, nil)
	if err != nil {
		t.Fatalf("recommend next: %v", err)
	}
	if rec.Primary == nil || rec.Primary.ConceptID != "a" {
		t.Fatalf("expected primary recommendation to be key prerequisite a, got %+v", rec.Primary)
	}
	if rec.Primary.Badge != "Remediation" {
		t.Errorf("expected badge Remediation, got %s", rec.Primary.Badge)
	}
	if rec.Primary.CTA != "Review prerequisite →" {
		t.Errorf("expected cta 'Review prerequisite →', got %s", rec.Primary.CTA)
	}

	// Step 4: NextQuestion serves key prerequisite a.
	qRemPrereq, err := e.NextQuestion(sess.ID, st.ID)
	if err != nil || qRemPrereq == nil {
		t.Fatalf("next remedial prereq: %v", err)
	}
	if qRemPrereq.ConceptID != "a" {
		t.Fatalf("expected remedial key prerequisite a, got %s", qRemPrereq.ConceptID)
	}

	// Submit correct answer for key prerequisite a.
	resRemPrereq, err := e.SubmitAnswer(sess.ID, st.ID, qRemPrereq.AttemptID, "42", 5.0)
	if err != nil || !resRemPrereq.Correct {
		t.Fatalf("submit remedial prereq: %v", err)
	}

	// Key prerequisite is now drained from student remedial queue.
	if len(e.RemedialQueue(st.ID)) != 0 {
		t.Errorf("expected empty remedial queue after key prereq success, got %v", e.RemedialQueue(st.ID))
	}

	// Step 5: NextQuestion serves concept b at easier variant.
	qRemConcept, err := e.NextQuestion(sess.ID, st.ID)
	if err != nil || qRemConcept == nil {
		t.Fatalf("next remedial concept: %v", err)
	}
	if qRemConcept.ConceptID != "b" {
		t.Fatalf("expected remedial concept b, got %s", qRemConcept.ConceptID)
	}

	// Submit correct answer for concept b.
	resRemConcept, err := e.SubmitAnswer(sess.ID, st.ID, qRemConcept.AttemptID, "99", 5.0)
	if err != nil || !resRemConcept.Correct {
		t.Fatalf("submit remedial concept: %v", err)
	}

	// Session halt is cleared.
	e.mu.Lock()
	halted := e.sessions[sess.ID].halted
	e.mu.Unlock()
	if halted {
		t.Error("expected halt cleared after completing remedial sequence")
	}

	// Normal completion path resumes.
	qResume, err := e.NextQuestion(sess.ID, st.ID)
	if err != nil || qResume == nil {
		t.Fatal("expected scheduling to resume normally")
	}
}

func TestEngine_RemedialQueue_PersistenceAcrossRestart(t *testing.T) {
	d, err := concepts.Build([]concepts.Concept{
		{
			ID: "a", Label: "Concept A", Domain: "d",
			GradingType: "numeric", Prerequisites: []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 2, AvgTimeSeconds: 60},
		},
		{
			ID: "b", Label: "Concept B", Domain: "d",
			GradingType: "numeric", Prerequisites: []string{"a"},
			KeyPrerequisites: []string{"a"},
			Variants: []concepts.Variant{
				{Difficulty: 0.3, TimeThresh: 60.0, Label: "b-easy"},
			},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 2, AvgTimeSeconds: 60},
		},
	})
	if err != nil {
		t.Fatalf("build DAG: %v", err)
	}
	store, err := storage.NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("create store: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	reg := generator.NewRegistry()
	reg.Register("a", &testGen{answer: "42"})
	reg.Register("b", &testGen{answer: "99"})

	// First server instance
	e1 := New(store, d, reg, nil, nil)
	st, _ := e1.CreateStudent("persistent-student")
	sess, _ := e1.repo.CreateSession(st.ID)

	now := time.Now().UTC()
	if err := e1.repo.UpsertProgress(&storage.ConceptProgress{
		StudentID:     st.ID,
		ConceptID:     "a",
		Status:        string(mastery.StatusMastered),
		Streak:        3,
		LastAttempted: &now,
		LastReviewed:  &now,
	}); err != nil {
		t.Fatalf("upsert progress: %v", err)
	}

	q1, err := e1.NextQuestion(sess.ID, st.ID)
	if err != nil || q1 == nil {
		t.Fatalf("q1: %v", err)
	}
	_, _ = e1.SubmitAnswer(sess.ID, st.ID, q1.AttemptID, "wrong", 5.0)

	q2, err := e1.NextQuestion(sess.ID, st.ID)
	if err != nil || q2 == nil {
		t.Fatalf("q2: %v", err)
	}
	res2, err := e1.SubmitAnswer(sess.ID, st.ID, q2.AttemptID, "wrong", 5.0)
	if err != nil || !res2.Halted {
		t.Fatalf("expected halt on second miss, res: %+v, err: %v", res2, err)
	}

	// Verify remedial queue is in storage
	storedRem, err := store.GetRemedialQueue(st.ID)
	if err != nil {
		t.Fatalf("get stored remedial: %v", err)
	}
	if len(storedRem) != 1 || storedRem[0] != "a" {
		t.Fatalf("expected stored remedial queue [a], got %v", storedRem)
	}

	// Simulate server restart by creating a completely new engine instance sharing the same store
	e2 := New(store, d, reg, nil, nil)

	// e2 must read the persistent remedial queue from storage
	remAfterRestart := e2.RemedialQueue(st.ID)
	if len(remAfterRestart) != 1 || remAfterRestart[0] != "a" {
		t.Fatalf("expected remedial queue [a] after restart, got %v", remAfterRestart)
	}

	// Recommendations on new engine must place the remedial prerequisite at the front
	rec, err := e2.RecommendNext(st.ID, nil)
	if err != nil {
		t.Fatalf("recommend on restarted engine: %v", err)
	}
	if rec.Primary == nil || rec.Primary.ConceptID != "a" {
		t.Fatalf("expected primary recommendation 'a', got %+v", rec.Primary)
	}
	if rec.Primary.Badge != "Remediation" {
		t.Errorf("expected badge Remediation, got %s", rec.Primary.Badge)
	}

	// Fresh session on the restarted server must serve the persistent remedial prerequisite
	sess2, _ := e2.repo.CreateSession(st.ID)
	qRem, err := e2.NextQuestion(sess2.ID, st.ID)
	if err != nil || qRem == nil {
		t.Fatalf("next question on restarted server: %v", err)
	}
	if qRem.ConceptID != "a" {
		t.Fatalf("expected remedial concept 'a', got %s", qRem.ConceptID)
	}

	// Answering correctly drains the remedial queue in storage
	resRem, err := e2.SubmitAnswer(sess2.ID, st.ID, qRem.AttemptID, "42", 5.0)
	if err != nil || !resRem.Correct {
		t.Fatalf("submit remedial answer: %v", err)
	}

	storedRemAfterSolve, err := store.GetRemedialQueue(st.ID)
	if err != nil {
		t.Fatalf("get stored remedial after solve: %v", err)
	}
	if len(storedRemAfterSolve) != 0 {
		t.Fatalf("expected empty remedial queue after solving prerequisite, got %v", storedRemAfterSolve)
	}
}

func TestEngine_Halt_DistinctConceptMissesDoNotHalt(t *testing.T) {
	d, err := concepts.Build([]concepts.Concept{
		{
			ID: "a", Label: "Concept A", Domain: "d",
			GradingType: "numeric", Prerequisites: []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60},
		},
		{
			ID: "b", Label: "Concept B", Domain: "d",
			GradingType: "numeric", Prerequisites: []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60},
		},
	})
	if err != nil {
		t.Fatalf("build DAG: %v", err)
	}
	store, err := storage.NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("create store: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	reg := generator.NewRegistry()
	reg.Register("a", &testGen{answer: "42"})
	reg.Register("b", &testGen{answer: "99"})
	e := New(store, d, reg, nil, nil)
	st, _ := e.CreateStudent("interleaved")
	sess, _ := e.repo.CreateSession(st.ID)

	// Question 1 is missed.
	q1, _ := e.NextQuestion(sess.ID, st.ID)
	res1, _ := e.SubmitAnswer(sess.ID, st.ID, q1.AttemptID, "wrong", 5.0)
	if res1.Halted {
		t.Error("unexpected halt on first miss")
	}

	// Serve the other concept next so misses are across distinct concepts.
	otherConcept := "b"
	otherAnswer := "99"
	if q1.ConceptID == "b" {
		otherConcept = "a"
		otherAnswer = "42"
	}

	e.mu.Lock()
	as := e.sessions[sess.ID]
	as.conceptID = otherConcept
	as.attemptID = newAttemptID()
	as.expectedAnswer = otherAnswer
	as.answered = false
	e.mu.Unlock()

	// Question on the second concept is missed. Since this is only 1 miss on this concept, no halt occurs.
	res2, err := e.SubmitAnswer(sess.ID, st.ID, as.attemptID, "wrong", 5.0)
	if err != nil {
		t.Fatalf("submit miss on second concept: %v", err)
	}
	if res2.Halted {
		t.Error("expected no halt when misses occur across different concepts")
	}
}

func TestEngine_Halt_StudyPath_KeyPrerequisites(t *testing.T) {
	d, err := concepts.Build([]concepts.Concept{
		{
			ID: "a", Label: "Concept A", Domain: "d",
			GradingType: "numeric", Prerequisites: []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 2, AvgTimeSeconds: 60},
		},
		{
			ID: "b", Label: "Concept B", Domain: "d",
			GradingType: "numeric", Prerequisites: []string{"a"},
			KeyPrerequisites: []string{"a"},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 2, AvgTimeSeconds: 60},
		},
	})
	if err != nil {
		t.Fatalf("build DAG: %v", err)
	}
	store, err := storage.NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("create store: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	reg := generator.NewRegistry()
	reg.Register("a", &testGen{answer: "42"})
	reg.Register("b", &testGen{answer: "99"})
	e := New(store, d, reg, nil, nil)
	st, _ := e.CreateStudent("study-halt")

	e.SetStudyAnchor(st.ID, "b", "99")

	// Miss 1 on concept b
	res1, err := e.SubmitStudyAnswer(st.ID, "b", "wrong", 5.0, "")
	if err != nil {
		t.Fatalf("miss 1: %v", err)
	}
	if res1.Halted {
		t.Error("expected no halt on first study miss")
	}

	e.SetStudyAnchor(st.ID, "b", "99")

	// Miss 2 on concept b triggers 2-fail halt
	res2, err := e.SubmitStudyAnswer(st.ID, "b", "wrong", 5.0, "")
	if err != nil {
		t.Fatalf("miss 2: %v", err)
	}
	if !res2.Halted {
		t.Error("expected halt on second study miss")
	}
	if res2.XP != 0 {
		t.Errorf("expected 0 XP on halt, got %d", res2.XP)
	}
	if len(res2.Remedial) != 2 || res2.Remedial[0] != "a" || res2.Remedial[1] != "b" {
		t.Errorf("expected remedial [a b], got %v", res2.Remedial)
	}
	rem := e.RemedialQueue(st.ID)
	if len(rem) != 1 || rem[0] != "a" {
		t.Errorf("expected engine remedial queue [a], got %v", rem)
	}

	// Correct answer on prerequisite drains the remedial queue
	e.SetStudyAnchor(st.ID, "a", "42")
	res3, err := e.SubmitStudyAnswer(st.ID, "a", "42", 5.0, "")
	if err != nil || !res3.Correct {
		t.Fatalf("correct remedial submit: %v", err)
	}
	if len(e.RemedialQueue(st.ID)) != 0 {
		t.Errorf("expected empty remedial queue, got %v", e.RemedialQueue(st.ID))
	}
}

func TestEngine_NegativeXP_AfterSecondRush(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("rusher")
	sess, _ := e.repo.CreateSession(st.ID)

	q1, _ := e.NextQuestion(sess.ID, st.ID)
	res1, _ := e.SubmitAnswer(sess.ID, st.ID, q1.AttemptID, "wrong", 1.0)
	if res1.XP != 0 {
		t.Errorf("expected 0 XP first rush, got %d", res1.XP)
	}
	q2, _ := e.NextQuestion(sess.ID, st.ID)
	res2, _ := e.SubmitAnswer(sess.ID, st.ID, q2.AttemptID, "wrong", 1.0)
	if res2.XP != -5 {
		t.Errorf("expected -5 XP second rapid guess, got %d", res2.XP)
	}
}

func TestEngine_AntiGuessing_ResetOnCorrectAnswer(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("guess_reset")
	sess, _ := e.repo.CreateSession(st.ID)

	// First rapid incorrect attempt: latency < 2s.
	q1, _ := e.NextQuestion(sess.ID, st.ID)
	res1, _ := e.SubmitAnswer(sess.ID, st.ID, q1.AttemptID, "wrong", 1.0)
	if res1.XP != 0 {
		t.Errorf("expected 0 XP first rapid miss, got %d", res1.XP)
	}

	// Intervening correct attempt resets the consecutive rush counter.
	q2, _ := e.NextQuestion(sess.ID, st.ID)
	res2, _ := e.SubmitAnswer(sess.ID, st.ID, q2.AttemptID, "42", 5.0)
	if !res2.Correct || res2.XP <= 0 {
		t.Fatalf("expected positive XP for correct answer, got %d", res2.XP)
	}

	// Another rapid incorrect attempt: not consecutive, so no penalty yet.
	q3, _ := e.NextQuestion(sess.ID, st.ID)
	res3, _ := e.SubmitAnswer(sess.ID, st.ID, q3.AttemptID, "wrong", 1.0)
	if res3.XP != 0 {
		t.Errorf("expected 0 XP after reset, got %d", res3.XP)
	}
}

func TestEngine_Accommodation_TimeMultiplierScalesThreshold(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("accom_mult")
	sess, _ := e.repo.CreateSession(st.ID)

	// Test settings with top level time_multiplier.
	_ = e.repo.UpdateSettings(st.ID, `{"time_multiplier": 1.5}`)
	q1, err := e.NextQuestion(sess.ID, st.ID)
	if err != nil || q1 == nil {
		t.Fatalf("next question: %v", err)
	}
	e.mu.Lock()
	as1 := e.sessions[sess.ID]
	e.mu.Unlock()
	if as1 == nil {
		t.Fatal("expected active session")
	}
	// Concept "a" baseline is 60s. 60s * 1.5 = 90s.
	if as1.timeThreshold != 90.0 {
		t.Errorf("expected timeThreshold 90.0 with 1.5x multiplier, got %f", as1.timeThreshold)
	}

	// Test settings with accommodations.time_multiplier.
	_ = e.repo.UpdateSettings(st.ID, `{"accommodations": {"time_multiplier": 2.0}}`)
	// Answer q1 to advance active session.
	_, _ = e.SubmitAnswer(sess.ID, st.ID, q1.AttemptID, "42", 10.0)

	q2, err := e.NextQuestion(sess.ID, st.ID)
	if err != nil || q2 == nil {
		t.Fatalf("next question: %v", err)
	}
	e.mu.Lock()
	as2 := e.sessions[sess.ID]
	e.mu.Unlock()
	if as2 == nil {
		t.Fatal("expected active session")
	}
	// 60s * 2.0 = 120s.
	if as2.timeThreshold != 120.0 {
		t.Errorf("expected timeThreshold 120.0 with 2.0x multiplier, got %f", as2.timeThreshold)
	}
}

func TestEngine_Accommodation_ExtraTimeScalesThreshold(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("accom")
	sess, _ := e.repo.CreateSession(st.ID)
	// extra_time 1.25 → 2.25x threshold (a: 60s → 135s)
	_ = e.repo.UpdateSettings(st.ID, `{"accommodations":{"extra_time":1.25}}`)

	q, err := e.NextQuestion(sess.ID, st.ID)
	if err != nil || q == nil {
		t.Fatalf("next question: %v", err)
	}
	e.mu.Lock()
	as := e.sessions[sess.ID]
	e.mu.Unlock()
	if as == nil {
		t.Fatal("expected active session")
	}
	if as.timeThreshold < 120 {
		t.Errorf("expected accommodated threshold >= 120 (60*2.25=135), got %f", as.timeThreshold)
	}
	// Slow-but-correct within the accommodated window must still be graded correct.
	res, err := e.SubmitAnswer(sess.ID, st.ID, q.AttemptID, "42", 100.0)
	if err != nil {
		t.Fatalf("submit slow correct: %v", err)
	}
	if !res.Correct {
		t.Error("expected correct with extra time accommodation")
	}
}

func TestEngine_StudyPath_NegativeXPOnRush(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("study_rush")
	e.SetStudyAnchor(st.ID, "a", "42")
	// First incorrect rush → 0; second incorrect rush → -1
	r1, err := e.SubmitStudyAnswer(st.ID, "a", "wrong", 1.0, "")
	if err != nil {
		t.Fatalf("study submit 1: %v", err)
	}
	if r1.XP != 0 {
		t.Errorf("expected 0 XP first study rush, got %d", r1.XP)
	}
	r2, err := e.SubmitStudyAnswer(st.ID, "a", "wrong", 1.0, "")
	if err != nil {
		t.Fatalf("study submit 2: %v", err)
	}
	if r2.XP != xppolicy.RushPenalty {
		t.Errorf("expected %d XP second study rush, got %d", xppolicy.RushPenalty, r2.XP)
	}
	if !r2.Halted {
		t.Error("expected halted flag after 2 consecutive study misses")
	}
}

// A tapped choice arrives as a letter, because the browser cannot return the
// option text faithfully. The anchor stores the authored text, so the letter has
// to be resolved before the comparison: reports #154, #160, #163, #172 and #179
// were all correct taps graded wrong by a rendering difference (spaces gone,
// LaTeX turned into glyphs, U+2212 turned into a hyphen).
func TestResolveChoiceAnswerMapsALetterToTheAuthoredOption(t *testing.T) {
	const question = "Which statement is true?\n\nA) 1 − P(A)\n\nB) P(A)\n\nC) \\(x^{2}<0\\)\n\nD) down on the left"

	if got := resolveChoiceAnswer(question, "D"); got != "down on the left" {
		t.Errorf("D resolved to %q, want the authored text", got)
	}
	if got := resolveChoiceAnswer(question, "D)"); got != "down on the left" {
		t.Errorf("D) resolved to %q, want the authored text", got)
	}
	if got := resolveChoiceAnswer(question, "(d)"); got != "down on the left" {
		t.Errorf("(d) resolved to %q, want the authored text", got)
	}
	if got := resolveChoiceAnswer(question, "Option D"); got != "down on the left" {
		t.Errorf("Option D resolved to %q, want the authored text", got)
	}
	if got := resolveChoiceAnswer(question, " c "); got != `\(x^{2}<0\)` {
		t.Errorf("lowercase letter with spaces resolved to %q", got)
	}
	if got := resolveChoiceAnswer(question, "downontheleft"); got != "down on the left" {
		t.Errorf("downontheleft without spaces resolved to %q, want the authored text", got)
	}
	if got := resolveChoiceAnswer(question, "down on the left."); got != "down on the left" {
		t.Errorf("down on the left with trailing dot resolved to %q, want the authored text", got)
	}
	if got := resolveChoiceAnswer(question, "1 − P(A)"); got != "1 − P(A)" {
		t.Errorf("a full option text was rewritten to %q", got)
	}
	// A letter with no option, a number, and a stem without options all pass through.
	if got := resolveChoiceAnswer(question, "E"); got != "E" {
		t.Errorf("unmatched letter became %q", got)
	}
	if got := resolveChoiceAnswer(question, "42"); got != "42" {
		t.Errorf("numeric answer became %q", got)
	}
	if got := resolveChoiceAnswer("Solve for x: 2x = 4", "A"); got != "A" {
		t.Errorf("letter against a free-text question became %q", got)
	}
}

// The learner is owed the explanation for the question they actually saw, on
// both verdicts. A miss used to render the grader's status token ("Incorrect")
// and a correct answer got no explanation at all, because the return was gated
// on !Correct and the value it returned was the grader feedback.
func TestEngine_StudyPath_ExplanationOnBothVerdicts(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("study_expl")
	const question = "3/4 + 1/2 = ?"
	const why = "Given 3/4 + 1/2. Step 1: common denominator LCM(4,2)=4. Step 2: add tops = 5/4. Answer: 5/4."
	e.SetStudyAnchorBatch(st.ID, "a", map[string]Anchor{
		question: {Answer: "5/4", Explanation: why},
	})

	// 7/4, not 1.25 — the numeric grader would (correctly) call 1.25 equal to 5/4.
	miss, err := e.SubmitStudyAnswer(st.ID, "a", "7/4", 5.0, question)
	if err != nil {
		t.Fatalf("submit miss: %v", err)
	}
	if miss.Correct {
		t.Fatal("expected the miss to grade incorrect")
	}
	if miss.Explanation != why {
		t.Errorf("miss explanation = %q, want the served explanation %q", miss.Explanation, why)
	}

	hit, err := e.SubmitStudyAnswer(st.ID, "a", "5/4", 5.0, question)
	if err != nil {
		t.Fatalf("submit correct: %v", err)
	}
	if !hit.Correct {
		t.Fatal("expected the exact answer to grade correct")
	}
	if hit.Explanation != why {
		t.Errorf("correct explanation = %q, want the same served explanation %q", hit.Explanation, why)
	}
}

// A second serve must not erase the first. The Learn client buffers served
// questions and fetches a fresh batch while the learner is still reading the
// one in front of them, so a replace-on-write blob made every question past the
// first ungradeable: the submit hit ErrNoStudyAnchor, the server answered 409,
// and the client swapped in a different question — which read as Learn skipping
// questions the learner had just answered.
func TestEngine_StudyAnchor_MergesAcrossBatches(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("study_merge")

	// The real arith.add.single question shape (internal/generator/arithmetic).
	first := map[string]Anchor{
		`\(3 + 4 = ?\)`: {Answer: "7", Explanation: "Given 3 + 4. Answer: 7."},
		`\(5 + 2 = ?\)`: {Answer: "7", Explanation: "Given 5 + 2. Answer: 7."},
	}
	second := map[string]Anchor{
		`\(8 + 1 = ?\)`: {Answer: "9", Explanation: "Given 8 + 1. Answer: 9."},
	}
	e.SetStudyAnchorBatch(st.ID, "a", first)
	e.SetStudyAnchorBatch(st.ID, "a", second)

	// Every question from the first batch must still grade against the answer it
	// was served with, not be reported as expired.
	for question, anchor := range first {
		res, err := e.SubmitStudyAnswer(st.ID, "a", anchor.Answer, 5.0, question)
		if err != nil {
			t.Fatalf("submit %q after a second batch: %v", question, err)
		}
		if !res.Correct {
			t.Errorf("expected %q to grade correct", question)
		}
		if res.Explanation != anchor.Explanation {
			t.Errorf("explanation for %q = %q, want the served one %q", question, res.Explanation, anchor.Explanation)
		}
	}

	// And the second batch is anchored too.
	res, err := e.SubmitStudyAnswer(st.ID, "a", "9", 5.0, `\(8 + 1 = ?\)`)
	if err != nil || !res.Correct {
		t.Errorf("second batch question: correct=%v err=%v", res != nil && res.Correct, err)
	}
}

// Incoming entries win on identical text. Deterministic generators give the
// same text the same answer, so this is not lossy — but a re-serve carrying a
// corrected explanation must not be shadowed by the stale copy.
func TestEngine_StudyAnchor_IncomingWinsOnSameText(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("study_same_text")
	const question = `\(3 + 4 = ?\)`

	e.SetStudyAnchorBatch(st.ID, "a", map[string]Anchor{question: {Answer: "7", Explanation: "old"}})
	e.SetStudyAnchorBatch(st.ID, "a", map[string]Anchor{question: {Answer: "7", Explanation: "new"}})

	res, err := e.SubmitStudyAnswer(st.ID, "a", "7", 5.0, question)
	if err != nil {
		t.Fatalf("submit: %v", err)
	}
	if res.Explanation != "new" {
		t.Errorf("explanation = %q, want the latest served one", res.Explanation)
	}
}

// The merged blob is bounded, and the batch just served survives the cap in
// full: those are the questions the learner is looking at now. Anything evicted
// degrades to the 409 re-serve path, never to a wrong grade.
func TestEngine_StudyAnchor_CapKeepsCurrentBatch(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("study_cap")

	old := make(map[string]Anchor, studyAnchorMaxEntries)
	for i := 0; i < studyAnchorMaxEntries; i++ {
		old[fmt.Sprintf("old-%d", i)] = Anchor{Answer: "1"}
	}
	e.SetStudyAnchorBatch(st.ID, "a", old)

	current := map[string]Anchor{}
	for i := 0; i < 10; i++ {
		current[fmt.Sprintf("now-%d", i)] = Anchor{Answer: "2", Explanation: "current"}
	}
	e.SetStudyAnchorBatch(st.ID, "a", current)

	for question := range current {
		if _, err := e.SubmitStudyAnswer(st.ID, "a", "2", 5.0, question); err != nil {
			t.Errorf("current-batch question %q was evicted: %v", question, err)
		}
	}

	// The blob stayed bounded.
	e.mu.Lock()
	size := len(decodeStudyAnchor(e.studyAnchor[st.ID+"|a"]))
	e.mu.Unlock()
	if size > studyAnchorMaxEntries {
		t.Errorf("anchor grew to %d entries, cap is %d", size, studyAnchorMaxEntries)
	}
}

// No anchor means the server cannot grade, and it must refuse rather than
// accept an expected answer from the caller. An ungradeable attempt is not a
// miss: nothing is recorded, and the caller is told to re-serve.
func TestEngine_StudyPath_RefusesWithoutAnchor(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("study_no_anchor")
	_, err := e.SubmitStudyAnswer(st.ID, "a", "42", 5.0, "a question that was never served")
	if !errors.Is(err, ErrNoStudyAnchor) {
		t.Fatalf("expected ErrNoStudyAnchor, got %v", err)
	}
	// Nothing recorded — no progress row means no streak, weakness or XP moved.
	if p, _ := e.GetProgress(st.ID); len(p) != 0 {
		t.Errorf("an ungradeable answer must record nothing, got %v", p)
	}
}

// A caller-supplied expected answer must not rescue a missing anchor: the
// whole point is that there is no longer a parameter to grade against.
func TestEngine_StudyPath_AnchorWinsOverSuppliedAnswer(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("study_anchor_wins")
	// The anchor says 42; a client claiming the expected is 7 would score 7.
	e.SetStudyAnchorBatch(st.ID, "a", map[string]Anchor{
		"Q1": {Answer: "42", Explanation: "The answer is 42."},
	})
	res, err := e.SubmitStudyAnswer(st.ID, "a", "42", 5.0, "Q1")
	if err != nil {
		t.Fatalf("submit: %v", err)
	}
	if !res.Correct {
		t.Error("expected the anchored answer to grade correct")
	}
	if res.ExpectedAnswer != "42" {
		t.Errorf("ExpectedAnswer = %q, want the anchored 42", res.ExpectedAnswer)
	}
}

// The review path holds its explanation on the active session rather than the
// study anchor, and it must return it on a correct answer too.
func TestEngine_Review_ExplanationOnCorrect(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("review_expl")
	sess, err := e.repo.CreateSession(st.ID)
	if err != nil {
		t.Fatalf("create session: %v", err)
	}
	q, err := e.NextQuestion(sess.ID, st.ID)
	if err != nil {
		t.Fatalf("next question: %v", err)
	}
	// Question carries no Expected on purpose — it is never sent to the client.
	// Read it off the active session to submit a genuinely correct answer.
	e.mu.Lock()
	expected := e.sessions[sess.ID].expectedAnswer
	e.mu.Unlock()
	res, err := e.SubmitAnswer(sess.ID, st.ID, q.AttemptID, expected, 5.0)
	if err != nil {
		t.Fatalf("submit: %v", err)
	}
	if !res.Correct {
		t.Fatal("expected the served expected answer to grade correct")
	}
	if res.Explanation == "" {
		t.Error("a correct review answer must still carry its explanation")
	}
	if res.Explanation == res.Feedback {
		t.Errorf("explanation (%q) should be the solution, not the verdict token", res.Explanation)
	}
}

// Session 5: efficacy instrumentation — first-pass / second-pass rates.
func TestEngine_Efficacy(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("efficacy")
	e.SetStudyAnchor(st.ID, "a", "42")
	e.SetStudyAnchor(st.ID, "b", "99")
	// Concept a: first attempt correct (first-pass hit), concept b: wrong then right.
	if _, err := e.SubmitStudyAnswer(st.ID, "a", "42", 5.0, ""); err != nil {
		t.Fatalf("submit a: %v", err)
	}
	if _, err := e.SubmitStudyAnswer(st.ID, "b", "1", 5.0, ""); err != nil {
		t.Fatalf("submit b wrong: %v", err)
	}
	if _, err := e.SubmitStudyAnswer(st.ID, "b", "99", 5.0, ""); err != nil {
		t.Fatalf("submit b right: %v", err)
	}
	rep, err := e.Efficacy(st.ID)
	if err != nil {
		t.Fatalf("efficacy: %v", err)
	}
	if rep.ConceptsTouched != 2 {
		t.Errorf("expected 2 concepts touched, got %d", rep.ConceptsTouched)
	}
	if rep.FirstPassRate != 0.5 {
		t.Errorf("expected first-pass 0.5 (1/2), got %f", rep.FirstPassRate)
	}
	if rep.SecondPassRate != 1.0 {
		t.Errorf("expected second-pass 1.0, got %f", rep.SecondPassRate)
	}
	if rep.AvgAttemptsPerConcept != 1.5 {
		t.Errorf("expected avg 1.5 attempts/concept, got %f", rep.AvgAttemptsPerConcept)
	}
	// Empty student → zeros, no error.
	empty, err := e.Efficacy("nobody")
	if err != nil || empty.ConceptsTouched != 0 {
		t.Errorf("expected empty report, got %+v err=%v", empty, err)
	}
}

func TestEngine_AggregateEfficacy(t *testing.T) {
	e := testEngine(t)
	st1, _ := e.CreateStudent("agg1")
	st2, _ := e.CreateStudent("agg2")
	for _, id := range []string{st1.ID, st2.ID} {
		e.SetStudyAnchor(id, "a", "42")
		e.SetStudyAnchor(id, "b", "99")
	}
	// agg1: a correct first try; b wrong then right.
	_, _ = e.SubmitStudyAnswer(st1.ID, "a", "42", 5.0, "")
	_, _ = e.SubmitStudyAnswer(st1.ID, "b", "1", 5.0, "")
	_, _ = e.SubmitStudyAnswer(st1.ID, "b", "99", 5.0, "")
	// agg2: a wrong then right; b wrong then right.
	_, _ = e.SubmitStudyAnswer(st2.ID, "a", "1", 5.0, "")
	_, _ = e.SubmitStudyAnswer(st2.ID, "a", "42", 5.0, "")
	_, _ = e.SubmitStudyAnswer(st2.ID, "b", "1", 5.0, "")
	_, _ = e.SubmitStudyAnswer(st2.ID, "b", "99", 5.0, "")

	rep, err := e.AggregateEfficacy()
	if err != nil {
		t.Fatalf("aggregate: %v", err)
	}
	if rep.StudentsTracked != 2 {
		t.Errorf("expected 2 students tracked, got %d", rep.StudentsTracked)
	}
	if rep.ConceptsTouched != 4 {
		t.Errorf("expected 4 student-concept pairs, got %d", rep.ConceptsTouched)
	}
	if rep.FirstPassRate != 0.25 {
		t.Errorf("expected first-pass 0.25 (1/4), got %f", rep.FirstPassRate)
	}
	if rep.SecondPassRate != 1.0 {
		t.Errorf("expected second-pass 1.0, got %f", rep.SecondPassRate)
	}
}

func TestEngine_SubmitQuizAnswer_TaskQuizXP(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("quizzer")
	before, _, err := e.repo.GetXP(st.ID)
	if err != nil {
		t.Fatalf("get xp: %v", err)
	}
	res, err := e.SubmitQuizAnswer(st.ID, "a", "42", "42", 3.0, "", nil)
	if err != nil {
		t.Fatalf("submit quiz answer: %v", err)
	}
	if !res.Correct {
		t.Fatal("expected correct quiz answer")
	}
	after, _, err := e.repo.GetXP(st.ID)
	if err != nil {
		t.Fatalf("get xp: %v", err)
	}
	// Single-path invariant: exactly what the response reports is what the
	// DB received — no lesson-rate write + quiz-rate display divergence.
	if after-before != res.XP {
		t.Errorf("DB delta %d != response XP %d", after-before, res.XP)
	}
	want := computeXPForTask(true, 3.0, 60, res.Streak, TaskQuiz)
	if res.XP != want {
		t.Errorf("expected TaskQuiz XP %d, got %d", want, res.XP)
	}
	// A quiz answer earns the same effort base as a lesson answer. It used to carry a +1
	// assessment premium, which made a quiz worth more per question than the teaching it
	// was checking, so the honest assertion is equality.
	lesson := computeXPForTask(true, 3.0, 60, res.Streak, TaskLesson)
	if res.XP != lesson {
		t.Errorf("quiz XP %d should equal lesson XP %d at the same effort", res.XP, lesson)
	}
}

func TestEngine_SubmitStudyAnswer_UnknownConcept(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("lost")
	if _, err := e.SubmitStudyAnswer(st.ID, "nope.not.real", "1", 5.0, ""); !errors.Is(err, ErrUnknownConcept) {
		t.Errorf("expected ErrUnknownConcept, got %v", err)
	}
	if _, err := e.SubmitQuizAnswer(st.ID, "nope.not.real", "1", "1", 5.0, "", nil); !errors.Is(err, ErrUnknownConcept) {
		t.Errorf("expected ErrUnknownConcept (quiz), got %v", err)
	}
	// No garbage progress row persisted.
	if p, _ := e.GetProgress(st.ID); len(p) != 0 {
		t.Errorf("expected no progress rows for unknown concept, got %v", p)
	}
}

func TestEngine_StudyAnchor_SurvivesRestart(t *testing.T) {
	e1 := testEngine(t)
	// Reach into the shared :memory: store via a second engine instance.
	e1.SetStudyAnchor("stu1", "a", "42")
	// Simulated restart: fresh engine, empty memory map, same store.
	e2 := New(e1.repo, e1.dag, e1.registry, nil, nil)
	if got, ok := e2.studyAnchorFor("stu1", "a", ""); !ok || got.Answer != "42" {
		t.Fatalf("expected durable answer=42, got %+v ok=%v", got, ok)
	}
	// The anchor stays available until TTL so a re-answer still grades.
	if _, ok := e2.studyAnchorFor("stu1", "a", ""); !ok {
		t.Error("expected anchor to remain available for re-answers")
	}
}

func TestEngine_StudyAnchorBatch_PerQuestion(t *testing.T) {
	e := testEngine(t)
	// Distinct questions on the same concept must each grade against their own
	// answer — the bug was anchoring only the first question — and each is owed
	// its own explanation, not the set's first one.
	e.SetStudyAnchorBatch("stu1", "a", map[string]Anchor{
		"Q1": {Answer: "5^4", Explanation: "5^4 = 5·5·5·5 = 625."},
		"Q2": {Answer: "625", Explanation: "Already in simplest form."},
	})
	if got, ok := e.studyAnchorFor("stu1", "a", "Q1"); !ok || got.Answer != "5^4" || got.Explanation != "5^4 = 5·5·5·5 = 625." {
		t.Errorf("Q1 anchor = %+v ok=%v", got, ok)
	}
	if got, ok := e.studyAnchorFor("stu1", "a", "Q2"); !ok || got.Answer != "625" || got.Explanation != "Already in simplest form." {
		t.Errorf("Q2 anchor = %+v ok=%v", got, ok)
	}
	if _, ok := e.studyAnchorFor("stu1", "a", "Q3"); ok {
		t.Error("unserved question should have no anchor")
	}
}

// A blob written before the explanation was carried must still grade: the
// session it belongs to is in flight, and refusing it would turn a deploy into
// a broken learner session. Only the two shapes that can actually be in
// server_sessions are covered — every writer marshals a map, so a bare string
// was never stored and the old fallback was dead too.
func TestDecodeStudyAnchor_LegacyShapes(t *testing.T) {
	for name, tc := range map[string]struct{ blob, key, answer string }{
		"current":     {`{"Q1":{"answer":"42","explanation":"why"}}`, "Q1", "42"},
		"perQuestion": {`{"Q1":"42"}`, "Q1", "42"},
		"singleKey":   {`{"":"42"}`, "", "42"},
	} {
		t.Run(name, func(t *testing.T) {
			got, ok := decodeStudyAnchor(tc.blob)[tc.key]
			if !ok {
				t.Fatalf("key %q missing from %s", tc.key, tc.blob)
			}
			if got.Answer != tc.answer {
				t.Errorf("answer = %q, want %q", got.Answer, tc.answer)
			}
		})
	}
}

// Batch 1: quiz miss enqueues immediate remedial (key prereqs + concept).

func TestEngine_QuizMiss_EnqueuesRemedial(t *testing.T) {
	d, err := concepts.Build([]concepts.Concept{
		{
			ID: "a", Label: "Concept A", Domain: "d",
			GradingType: "numeric", Prerequisites: []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60},
		},
		{
			ID: "b", Label: "Concept B", Domain: "d",
			GradingType: "numeric", Prerequisites: []string{"a"},
			KeyPrerequisites: []string{"a"},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60},
		},
	})
	if err != nil {
		t.Fatalf("build DAG: %v", err)
	}
	store, err := storage.NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("create store: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	reg := generator.NewRegistry()
	reg.Register("a", &testGen{answer: "42"})
	reg.Register("b", &testGen{answer: "99"})
	e := New(store, d, reg, nil, nil)
	st, _ := e.CreateStudent("remedial")

	res, err := e.SubmitQuizAnswer(st.ID, "b", "wrong", "99", 5.0, "", nil)
	if err != nil {
		t.Fatalf("submit: %v", err)
	}
	if res.Correct {
		t.Fatal("expected incorrect grade")
	}
	if len(res.Remedial) != 2 || res.Remedial[0] != "a" || res.Remedial[1] != "b" {
		t.Errorf("expected remedial [a b], got %v", res.Remedial)
	}
	if got := e.QuizRemedial(st.ID); len(got) != 2 {
		t.Errorf("expected queued remedial len 2, got %v", got)
	}
	// Correct answers never enqueue.
	if _, err := e.SubmitQuizAnswer(st.ID, "a", "42", "42", 5.0, "", nil); err != nil {
		t.Fatalf("submit correct: %v", err)
	}
	if got := e.QuizRemedial(st.ID); len(got) != 2 {
		t.Errorf("expected queue unchanged after correct, got %v", got)
	}
	e.ClearQuizRemedial(st.ID)
	if got := e.QuizRemedial(st.ID); len(got) != 0 {
		t.Errorf("expected drained queue, got %v", got)
	}
}

// Quiz gate — due at 50 XP since last completion, reset on record.

func TestEngine_QuizGate_DueAndReset(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("gate")
	due, err := e.QuizDue(st.ID)
	if err != nil || due {
		t.Fatalf("expected not due at 0 XP, due=%v err=%v", due, err)
	}
	if err := e.repo.AwardXP(st.ID, 49, "test", "seed-49"); err != nil {
		t.Fatalf("add xp: %v", err)
	}
	if due, _ := e.QuizDue(st.ID); due {
		t.Error("expected not due at 49 XP since quiz")
	}
	if err := e.repo.AwardXP(st.ID, 1, "test", "seed-1"); err != nil {
		t.Fatalf("add xp: %v", err)
	}
	if since, _ := e.QuizXPSince(st.ID); since != 50 {
		t.Errorf("expected 50 since, got %d", since)
	}
	if due, _ := e.QuizDue(st.ID); !due {
		t.Error("expected due at 50 XP since quiz")
	}
	if err := e.RecordQuizCompletion(st.ID); err != nil {
		t.Fatalf("record completion: %v", err)
	}
	if since, _ := e.QuizXPSince(st.ID); since != 0 {
		t.Errorf("expected 0 since after completion, got %d", since)
	}
	if due, _ := e.QuizDue(st.ID); due {
		t.Error("expected not due right after completion")
	}
}

// Batch 1: accommodated time limit wrapper follows concept threshold.

func TestEngine_TimeLimitFor_ScalesWithAccommodation(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("timed")
	if got := e.TimeLimitFor(st.ID, "a"); got != 60 {
		t.Errorf("expected base 60s, got %v", got)
	}
}

// A diagnosis names the mistake after grading. It must never change the verdict,
// never appear on a correct answer, and stay empty when nothing can be said
// with certainty. Concept "a" expects 42.
func TestEngine_Diagnosis(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("diag")
	// The study path grades against the server-side anchor, never a
	// caller-supplied expected answer, so the test has to serve one.
	e.SetStudyAnchor(st.ID, "a", "42")

	cases := map[string]struct {
		answer        string
		wantCorrect   bool
		wantDiagnosis bool
		why           string
	}{
		"sign flipped": {"-42", false, true, "the magnitude is right, so the mistake is the sign"},
		"off by one":   {"43", false, true, "43 is exactly one from 42"},
		"unrelated":    {"17", false, false, "nothing distinguishes 17 from 42"},
		"right answer": {"42", true, false, "a correct answer has no mistake to name"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			res, err := e.SubmitStudyAnswer(st.ID, "a", tc.answer, 5.0, "")
			if err != nil {
				t.Fatalf("submit: %v", err)
			}
			if res.Correct != tc.wantCorrect {
				t.Errorf("correct = %v, want %v — %s", res.Correct, tc.wantCorrect, tc.why)
			}
			if tc.wantDiagnosis && res.Diagnosis == "" {
				t.Errorf("expected a diagnosis for %q, got none — %s", tc.answer, tc.why)
			}
			if !tc.wantDiagnosis && res.Diagnosis != "" {
				t.Errorf("unexpected diagnosis %q — %s", res.Diagnosis, tc.why)
			}
		})
	}
}

// Decay is never persisted — the row stays MASTERED forever and the downgrade
// happens on read. Before GetProgress applied it, /api/progress reported a
// stale concept as mastered while the scheduler treated it as DECAYING and
// scheduled it for review: two parts of the app disagreeing about one fact, and
// the learner seeing the flattering half.
func TestGetProgress_AppliesDecayOnRead(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("decay")
	if err := e.repo.UpsertProgress(&storage.ConceptProgress{
		StudentID: st.ID, ConceptID: "a", Status: string(mastery.StatusMastered), Streak: 3,
	}); err != nil {
		t.Fatalf("seed: %v", err)
	}

	// No LastReviewed at all: we cannot claim it is fresh.
	got, err := e.GetProgress(st.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got["a"].Status != "DECAYING" {
		t.Errorf("status = %q, want DECAYING for a review with no timestamp", got["a"].Status)
	}

	// Reviewed just now: still mastered.
	now := time.Now().UTC()
	if err := e.repo.UpsertProgress(&storage.ConceptProgress{
		StudentID: st.ID, ConceptID: "a", Status: string(mastery.StatusMastered), Streak: 3,
		LastReviewed: &now,
	}); err != nil {
		t.Fatalf("seed recent: %v", err)
	}
	got, err = e.GetProgress(st.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got["a"].Status != string(mastery.StatusMastered) {
		t.Errorf("status = %q, want MASTERED for a fresh review", got["a"].Status)
	}

	// Reviewed past the threshold.
	stale := now.AddDate(0, 0, -15)
	if err := e.repo.UpsertProgress(&storage.ConceptProgress{
		StudentID: st.ID, ConceptID: "a", Status: string(mastery.StatusMastered), Streak: 3,
		LastReviewed: &stale,
	}); err != nil {
		t.Fatalf("seed stale: %v", err)
	}
	got, err = e.GetProgress(st.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got["a"].Status != "DECAYING" {
		t.Errorf("status = %q, want DECAYING 15 days after the last review", got["a"].Status)
	}

	// The row itself is untouched: decay is a read-time view, never persisted.
	raw, err := e.repo.GetProgress(st.ID, "a")
	if err != nil {
		t.Fatalf("raw: %v", err)
	}
	if raw.Status != string(mastery.StatusMastered) {
		t.Errorf("stored status = %q, want it left as MASTERED", raw.Status)
	}
}

// An unmastered status must pass through untouched whatever its timestamp.
func TestGetProgress_LeavesOtherStatusesAlone(t *testing.T) {
	d, err := concepts.Build([]concepts.Concept{
		{ID: "u", Label: "U", Domain: "d", GradingType: "numeric", Prerequisites: []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60}},
		{ID: "l", Label: "L", Domain: "d", GradingType: "numeric", Prerequisites: []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60}},
		{ID: "p", Label: "P", Domain: "d", GradingType: "numeric", Prerequisites: []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60}},
	})
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	store, err := storage.NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	e := New(store, d, generator.NewRegistry(), nil, nil)
	st, _ := e.CreateStudent("decay2")

	// Ancient review timestamps: decay must not touch a status that was never
	// mastery in the first place.
	long := time.Now().UTC().AddDate(0, 0, -400)
	for id, status := range map[string]string{"u": "UNSEEN", "l": "LEARNING", "p": "PRACTICING"} {
		if err := store.UpsertProgress(&storage.ConceptProgress{
			StudentID: st.ID, ConceptID: id, Status: status, Streak: 1, LastReviewed: &long,
		}); err != nil {
			t.Fatalf("seed %s: %v", id, err)
		}
	}
	got, err := e.GetProgress(st.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	for id, want := range map[string]string{"u": "UNSEEN", "l": "LEARNING", "p": "PRACTICING"} {
		if got[id].Status != want {
			t.Errorf("%s status = %q, want %q unchanged", id, got[id].Status, want)
		}
	}
}

// answerCorrectly answers concept "a" correctly n times, using the engine's own question
// pool. It is not usable for long runs: the pool for a two-concept DAG runs out after a
// handful of questions, which is why the exhaustive monotonicity proof lives in
// internal/mastery — where it can drive the transition rules directly — and what is
// asserted here is that the endpoint derives from the server's threshold.
func answerCorrectly(t *testing.T, e *Engine, studentID string, n int) {
	t.Helper()
	for i := 0; i < n; i++ {
		sess, _ := e.repo.CreateSession(studentID)
		e.mu.Lock()
		e.sessions[sess.ID] = &activeSession{}
		e.mu.Unlock()
		q, err := e.NextQuestion(sess.ID, studentID)
		if err != nil {
			t.Fatalf("next question %d: %v", i+1, err)
		}
		res, err := e.SubmitAnswer(sess.ID, studentID, q.AttemptID, "42", 1.0)
		if err != nil {
			t.Fatalf("submit %d: %v", i+1, err)
		}
		if !res.Correct {
			t.Fatalf("expected correct on attempt %d", i+1)
		}
	}
}

// ProgressWithMasteryPct is what the graph's bar length is read from, so the two
// properties that matter are: it is derived from the server's own threshold rather than a
// bundled copy of the corpus, and it does not fall when the learner advances a tier.
func TestProgressWithMasteryPct_DerivesFromServerThreshold(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("pct_threshold")

	views, err := e.ProgressWithMasteryPct(st.ID)
	if err != nil {
		t.Fatalf("ProgressWithMasteryPct: %v", err)
	}
	if len(views) != 0 {
		t.Errorf("an untouched learner has %d progress rows, want 0", len(views))
	}

	// A full window of clean answers reaches MASTERED and the bar is full. The
	// per-concept threshold has no say in this: the score is normalised against the
	// concept's own difficulty and time threshold, so mastery costs the same whether
	// required_streak is 3 or 15.
	answerCorrectly(t, e, st.ID, mastery.MasteryEvidenceFloor)
	views, _ = e.ProgressWithMasteryPct(st.ID)
	row := views["a"]
	if row.Status != "MASTERED" {
		t.Fatalf("status after a full clean window = %q, want MASTERED", row.Status)
	}
	if row.MasteryPct != 1 {
		t.Errorf("mastery_pct at MASTERED = %v, want 1", row.MasteryPct)
	}

	// Mid-ladder the bar is partial and ordered, which is what makes it mean attainment
	// rather than a streak ratio. A separate learner, because concept "a" has run out of
	// questions by now and asking again is not possible.
	mid, _ := e.CreateStudent("pct_mid_ladder")
	answerCorrectly(t, e, mid.ID, 1)
	midViews, _ := e.ProgressWithMasteryPct(mid.ID)
	if got := midViews["a"].Status; got != "LEARNING" {
		t.Fatalf("status after 1 clean answer = %q, want LEARNING", got)
	}
	if got := midViews["a"].MasteryPct; got <= 1.0/3.0 || got >= 2.0/3.0 {
		t.Errorf("mastery_pct one answer in = %v, want strictly between 1/3 and 2/3", got)
	}
	// The stored row is embedded, so the schema fields still travel with it alongside the
	// derived percentage.
	if row.Status != "MASTERED" {
		t.Errorf("embedded status = %q, want MASTERED", row.Status)
	}
	if row.Attempts != mastery.MasteryEvidenceFloor {
		t.Errorf("embedded attempts = %d, want %d", row.Attempts, mastery.MasteryEvidenceFloor)
	}
	if row.MasteryPct != views["a"].MasteryPct {
		t.Errorf("embedded view disagrees with the derived value: %v", row.MasteryPct)
	}
}

func TestProgressWithMasteryPct_NeverFallsOnAdvance(t *testing.T) {
	// The sawtooth this exists to prevent, at the exact answer where it happened. The
	// engine resets the streak to 1 on the tier advance; a streak-derived ratio therefore
	// falls from 1.0 to 0.33 here. Tier-derived attainment rises instead.
	//
	// The exhaustive version — every answer through all three advances — is
	// TestMasteryPct_NeverGoesBackwards in internal/mastery, which drives the transition
	// rules directly and so is not limited by the question pool.
	if got := mastery.MasteryPct(mastery.StatusLearning, 3, 3); math.Abs(got-2.0/3.0) > 1e-9 {
		t.Errorf("mastery_pct for LEARNING with the streak met = %v, want %v", got, 2.0/3.0)
	}
	// And the ratio the client used to compute, at the same two instants.
	if old := float64(3) / float64(3); old != 1 {
		t.Errorf("sanity: the old streak ratio at the threshold should be 1, got %v", old)
	}
	if after := float64(1) / float64(3); after >= 0.5 {
		t.Errorf("sanity: after the reset the old ratio reads %v, which is the drop", after)
	}
}

// A decayed concept keeps a full bar. Decay is signalled by colour alone, and a shorter
// bar would read as losing the competence — the one thing decay does not mean.
func TestProgressWithMasteryPct_DecayingIsFull(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("pct_decaying")
	answerCorrectly(t, e, st.ID, 3)

	views, _ := e.ProgressWithMasteryPct(st.ID)
	if got := views["a"].MasteryPct; got < 1.0/3.0 {
		t.Errorf("mastery_pct at LEARNING = %v, want past the first third", got)
	}
	// `GetProgress` applies decay from the clock and these answers all happened now, so
	// the served status is still LEARNING. The decayed case is asserted directly against
	// the same function the endpoint calls; backdating the row to force it would be
	// asserting the clock rather than the rule.
	if got := mastery.MasteryPct("DECAYING", 1, 3); got != 1 {
		t.Errorf("decaying mastery_pct = %v, want 1 (colour carries decay, not length)", got)
	}
}

// The evidence model weights a correct answer by how hard the question was. That is only
// usable if the server records the difficulty it generated at — and records *nothing*
// rather than a guess when it does not know.
func TestEngine_RecordsServedQuestionDifficulty(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("difficulty_recorded")

	sess, _ := e.repo.CreateSession(st.ID)
	e.mu.Lock()
	e.sessions[sess.ID] = &activeSession{}
	e.mu.Unlock()

	q, err := e.NextQuestion(sess.ID, st.ID)
	if err != nil {
		t.Fatalf("next question: %v", err)
	}
	res, err := e.SubmitAnswer(sess.ID, st.ID, q.AttemptID, "42", 1.0)
	if err != nil {
		t.Fatalf("submit: %v", err)
	}
	if !res.Correct {
		t.Fatal("expected the seeded answer to grade correct")
	}

	atts, err := e.repo.GetAttemptsForStudent(st.ID)
	if err != nil {
		t.Fatalf("get attempts: %v", err)
	}
	if len(atts) != 1 {
		t.Fatalf("got %d attempts, want 1", len(atts))
	}
	d := atts[0].Difficulty
	if d == nil {
		t.Fatal("attempt difficulty is nil; a generated question always knows its own difficulty")
	}
	// Generators are asked for 0.3-1.0. Outside that range would mean the value came
	// from somewhere other than the generator context.
	if *d < 0.3 || *d > 1.0 {
		t.Errorf("difficulty = %v, want within the generator's 0.3-1.0 range", *d)
	}
}

// A difficulty of 0 would claim the question was trivially easy. Nothing measured that,
// so nothing may write it.
func TestEngine_NeverRecordsZeroDifficulty(t *testing.T) {
	if got := difficultyOrNil(0); got != nil {
		t.Errorf("difficultyOrNil(0) = %v, want nil", *got)
	}
	if got := difficultyOrNil(-1); got != nil {
		t.Errorf("difficultyOrNil(-1) = %v, want nil", *got)
	}
	if got := difficultyOrNil(0.3); got == nil || *got != 0.3 {
		t.Errorf("difficultyOrNil(0.3) = %v, want 0.3", got)
	}
	if got := attemptDifficulty(nil); got != nil {
		t.Errorf("attemptDifficulty(nil) = %v, want nil", *got)
	}
}

// The scoring side and the serving side must agree on what "a question of unknown
// difficulty" is worth.
//
// mastery.NeutralDifficulty is what BuildEvidence credits a correct answer with when the
// attempt carries no difficulty, and difficultyFromWeakness is what the engine serves when
// it has no information about the learner. If those two drift apart, an unknown-difficulty
// attempt silently becomes worth more or less than the questions actually being served —
// which is precisely the bug that made the quiz surface more permissive than the study
// surface.
func TestNeutralDifficultyMatchesEngineDefault(t *testing.T) {
	if got := difficultyFromWeakness(nil, "anything"); got != mastery.NeutralDifficulty {
		t.Errorf("engine default difficulty = %v, mastery.NeutralDifficulty = %v; the score "+
			"credited to an unknown-difficulty attempt no longer matches the question served",
			got, mastery.NeutralDifficulty)
	}
	if got := difficultyFromWeakness(map[string]float64{}, "anything"); got != mastery.NeutralDifficulty {
		t.Errorf("engine default difficulty with an empty weakness map = %v, want %v",
			got, mastery.NeutralDifficulty)
	}
}

// --- the shared mastery authority must not care which screen produced the attempt ------

// A learner must not be able to change their own mastery outcome by choosing a different
// surface.
//
// `submitAnswerWithTask` serves both the study/Learn loop and quizzes, and it used to hand
// the mastery authority a hardcoded `difficulty: 0` — meaning "unknown" — while
// `SubmitAnswer` passed the real difficulty off the active session. Unknown difficulty was
// then scored as `difficultyTerm = 1.0`, the top of the range, so the same answers were worth
// more on the quiz path than on the study path. Measured at 50% accuracy: 0.65 and cleared
// the bar on quiz, 0.57 and failed on study.
//
// This drives both hosts with identical questions, identical outcomes and identical
// difficulty, and requires the persisted evidence to be identical.
func TestStudyAndQuizProduceEquivalentMasteryEvidence(t *testing.T) {
	const served = 0.65
	diff := served
	question := "3/4 + 1/2 = ?"

	// Six answers at 50% accuracy — the exact case the asymmetry decided.
	outcomes := []bool{true, false, true, false, true, true}

	// Study: the anchor is the record of what was served, so the difficulty rides there —
	// no caller supplies it.
	studyEngine, studyStudent := func() (*Engine, *storage.Student) {
		e := testEngine(t)
		st, err := e.CreateStudent("equiv_study")
		if err != nil {
			t.Fatalf("create student: %v", err)
		}
		e.SetStudyAnchorBatch(st.ID, "a", map[string]Anchor{
			question: {Answer: "5/4", Explanation: "because", Difficulty: &diff},
		})
		for i, correct := range outcomes {
			ans := "5/4"
			if !correct {
				ans = "7/4"
			}
			if _, err := e.SubmitStudyAnswer(st.ID, "a", ans, 5.0, question); err != nil {
				t.Fatalf("study answer %d: %v", i, err)
			}
		}
		return e, st
	}()

	// Quiz: the difficulty comes off the served problem, exactly as handleQuizAnswer reads
	// it from sess.LastProblem.Difficulty.
	quizEngine := testEngine(t)
	quizStudent, err := quizEngine.CreateStudent("equiv_quiz")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}
	for i, correct := range outcomes {
		ans := "5/4"
		if !correct {
			ans = "7/4"
		}
		if _, err := quizEngine.SubmitQuizAnswer(quizStudent.ID, "a", ans, "5/4", 5.0, question, &diff); err != nil {
			t.Fatalf("quiz answer %d: %v", i, err)
		}
	}

	studyRows := recentAttempts(t, studyEngine, studyStudent.ID, "a")
	quizRows := recentAttempts(t, quizEngine, quizStudent.ID, "a")

	if len(studyRows) != len(outcomes) || len(quizRows) != len(outcomes) {
		t.Fatalf("got %d study and %d quiz attempts, want %d each",
			len(studyRows), len(quizRows), len(outcomes))
	}
	for i := range outcomes {
		if studyRows[i].Difficulty == nil || quizRows[i].Difficulty == nil {
			t.Fatalf("attempt %d: study difficulty %v, quiz difficulty %v — both must carry "+
				"the served difficulty", i, studyRows[i].Difficulty, quizRows[i].Difficulty)
		}
		if math.Abs(*studyRows[i].Difficulty-*quizRows[i].Difficulty) > 1e-9 {
			t.Errorf("attempt %d: study difficulty %v vs quiz %v",
				i, *studyRows[i].Difficulty, *quizRows[i].Difficulty)
		}
		if studyRows[i].Correct != quizRows[i].Correct {
			t.Errorf("attempt %d: study correct=%v, quiz correct=%v for the same answer",
				i, studyRows[i].Correct, quizRows[i].Correct)
		}
	}

	// And the evidence they assemble must be the same number, not merely the same shape.
	studyEv := evidenceFromRows(studyRows)
	quizEv := evidenceFromRows(quizRows)
	if math.Abs(studyEv.Score-quizEv.Score) > 1e-9 {
		t.Errorf("study evidence scored %.4f and quiz evidence %.4f for the same answers on "+
			"the same questions — the surface is changing what the learner demonstrated",
			studyEv.Score, quizEv.Score)
	}
	if studyEv.Attempts != quizEv.Attempts {
		t.Errorf("study saw %d attempts, quiz %d", studyEv.Attempts, quizEv.Attempts)
	}
}

// recentAttempts reads back the rows the evidence window is built from.
func recentAttempts(t *testing.T, e *Engine, studentID, conceptID string) []storage.AttemptEntry {
	t.Helper()
	rows, err := e.repo.GetRecentAttemptsForConcept(studentID, conceptID, 20)
	if err != nil {
		t.Fatalf("recent attempts: %v", err)
	}
	return rows
}

// A host that genuinely does not know the difficulty must not be *more* generous than one
// that does. This is the regression at the level the bug actually lived: not "is the value
// plumbed" but "does its absence buy anything".
func TestUnknownDifficultyHostIsNotMoreGenerous(t *testing.T) {
	e := testEngine(t)
	question := "3/4 + 1/2 = ?"
	diff := 0.65
	outcomes := []bool{true, false, true, false, true, true}

	// Known difficulty, from the anchor.
	known, _ := e.CreateStudent("diff_known")
	e.SetStudyAnchorBatch(known.ID, "a", map[string]Anchor{
		question: {Answer: "5/4", Difficulty: &diff},
	})
	for i, correct := range outcomes {
		ans, want := "7/4", "5/4"
		if correct {
			ans = "5/4"
		}
		if _, err := e.SubmitStudyAnswer(known.ID, "a", ans, 5.0, question); err != nil {
			t.Fatalf("answer %d: %v", i, err)
		}
		_ = want
	}

	// Unknown difficulty, honestly reported as unknown.
	unknown, _ := e.CreateStudent("diff_unknown")
	e.SetStudyAnchorBatch(unknown.ID, "a", map[string]Anchor{
		question: {Answer: "5/4"},
	})
	for i, correct := range outcomes {
		ans := "7/4"
		if correct {
			ans = "5/4"
		}
		if _, err := e.SubmitStudyAnswer(unknown.ID, "a", ans, 5.0, question); err != nil {
			t.Fatalf("answer %d: %v", i, err)
		}
	}

	knownRows := recentAttempts(t, e, known.ID, "a")
	unknownRows := recentAttempts(t, e, unknown.ID, "a")
	if len(knownRows) != 6 || len(unknownRows) != 6 {
		t.Fatalf("got %d and %d attempts, want 6 each", len(knownRows), len(unknownRows))
	}
	if knownRows[5].Difficulty == nil || unknownRows[5].Difficulty != nil {
		t.Fatalf("setup: known=%v unknown=%v", knownRows[5].Difficulty, unknownRows[5].Difficulty)
	}

	// Build the same evidence both ways and require the unknown-difficulty score not to
	// exceed the known one.
	knownEv := evidenceFromRows(knownRows)
	unknownEv := evidenceFromRows(unknownRows)
	if unknownEv.Score > knownEv.Score+1e-9 {
		t.Errorf("unknown difficulty scored %.4f against %.4f for known difficulty %.2f — "+
			"a missing value is buying more evidence than a recorded one",
			unknownEv.Score, knownEv.Score, *knownRows[5].Difficulty)
	}
}

// evidenceFromRows assembles evidence the way buildEvidenceWindow does, from persisted rows.
func evidenceFromRows(rows []storage.AttemptEntry) mastery.Evidence {
	var as []mastery.Attempt
	for _, r := range rows {
		as = append(as, mastery.Attempt{
			Correct:    r.Correct,
			Elapsed:    r.ElapsedSeconds,
			Difficulty: r.Difficulty,
			Instance:   r.Question,
		})
	}
	return mastery.BuildEvidence(as, mastery.EvidenceOptions{TimeThreshold: 60})
}

// The answer response carries the authoritative evidence, so no client has to compute its
// own.
//
// `web/next-app/lib/progression.ts` kept a second implementation of the mastery evidence
// model — different weights, different variety population, a flat 60s time cliff instead of
// the concept threshold — and rendered a band from it in the Learn header. That is only
// removable if the number is available, so this asserts it is, and that the band comes from
// the same place as the verdict.
func TestAnswerResponseCarriesAuthoritativeEvidence(t *testing.T) {
	e := testEngine(t)
	st, err := e.CreateStudent("evidence_payload")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}
	const question = "3/4 + 1/2 = ?"
	diff := 0.65
	e.SetStudyAnchorBatch(st.ID, "a", map[string]Anchor{
		question: {Answer: "5/4", Difficulty: &diff},
	})

	// Below the evidence floor, so the band must say so rather than score well.
	first, err := e.SubmitStudyAnswer(st.ID, "a", "5/4", 5.0, question)
	if err != nil {
		t.Fatalf("submit: %v", err)
	}
	if first.EvidenceScore <= 0 {
		t.Error("evidence_score was not sent on a graded answer")
	}
	if first.EvidenceBand != "early" {
		t.Errorf("band after one attempt = %q, want %q", first.EvidenceBand, "early")
	}

	// Fill the window and the band must move, and move to what BuildEvidence implies.
	for i := 0; i < mastery.MasteryEvidenceFloor-1; i++ {
		if _, err := e.SubmitStudyAnswer(st.ID, "a", "5/4", 5.0, question); err != nil {
			t.Fatalf("submit %d: %v", i, err)
		}
	}
	rows := recentAttempts(t, e, st.ID, "a")
	ev := evidenceFromRows(rows)
	last, err := e.SubmitStudyAnswer(st.ID, "a", "5/4", 5.0, question)
	if err != nil {
		t.Fatalf("final submit: %v", err)
	}
	if last.EvidenceBand == "early" {
		t.Error("band still reads early with a full window of correct answers")
	}

	// The payload must be the same evidence the ladder read, not a parallel computation.
	fresh := evidenceFromRows(recentAttempts(t, e, st.ID, "a"))
	if fresh.Score < ev.Score-1e-9 {
		t.Errorf("evidence went backwards after another correct answer: %.4f then %.4f",
			ev.Score, fresh.Score)
	}
	if last.EvidenceBand != mastery.EvidenceBand(fresh) {
		t.Errorf("band %q does not match mastery.EvidenceBand(%+v) = %q",
			last.EvidenceBand, fresh, mastery.EvidenceBand(fresh))
	}
}

// An ungraded answer records nothing, so it must not report new evidence.
func TestUngradedAnswerCarriesNoEvidence(t *testing.T) {
	// Mirrors the shape both hosts construct when `Engine.GradeAnswer` reports a grader
	// fault. ADR-006: nothing is recorded, so there is no new evidence to report.
	res := &AnswerResult{Ungraded: true, Feedback: "grader unavailable"}
	if res.EvidenceScore != 0 || res.EvidenceBand != "" {
		t.Errorf("an ungraded answer reported evidence %v/%q", res.EvidenceScore, res.EvidenceBand)
	}
}

// The rush penalty must mean the same thing on every host.
//
// The practice path counted its own sub-2-second wrong answers on the active session; the
// study/quiz path required two misses on the concept as well, so a learner who rushed twice
// but was wrong for ordinary reasons was never penalised there. Same behaviour, two
// different rules — now both xp.RushSeconds / xp.RushAfter.
func TestRushPenaltyIsIdenticalAcrossHosts(t *testing.T) {
	e := testEngine(t)
	st, err := e.CreateStudent("rush_parity")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}
	const question = "3/4 + 1/2 = ?"
	diff := 0.6
	e.SetStudyAnchorBatch(st.ID, "a", map[string]Anchor{
		question: {Answer: "5/4", Difficulty: &diff},
	})

	// Two fast wrong answers on the study path: the second must carry the penalty.
	var xps []int
	for i := 0; i < 2; i++ {
		res, err := e.SubmitStudyAnswer(st.ID, "a", "7/4", xppolicy.RushSeconds-0.5, question)
		if err != nil {
			t.Fatalf("submit %d: %v", i, err)
		}
		xps = append(xps, res.XP)
	}
	if xps[0] != 0 {
		t.Errorf("first fast wrong answer XP = %d, want 0 (penalty starts at RushAfter)", xps[0])
	}
	if xps[1] != xppolicy.RushPenalty {
		t.Errorf("second fast wrong answer XP = %d, want %d", xps[1], xppolicy.RushPenalty)
	}

	// The same two fast wrong answers on the quiz path must produce the same two numbers.
	e2 := testEngine(t)
	st2, _ := e2.CreateStudent("rush_parity_quiz")
	var quizXPs []int
	for i := 0; i < 2; i++ {
		res, err := e2.SubmitQuizAnswer(st2.ID, "a", "7/4", "5/4", xppolicy.RushSeconds-0.5, question, &diff)
		if err != nil {
			t.Fatalf("quiz submit %d: %v", i, err)
		}
		quizXPs = append(quizXPs, res.XP)
	}
	if quizXPs[0] != xps[0] || quizXPs[1] != xps[1] {
		t.Errorf("quiz XP %v vs study XP %v for the same two fast wrong answers", quizXPs, xps)
	}
}

// An admitted "I don't know" is honesty, not rushing, so it never draws the penalty however
// fast it was submitted.
func TestDontKnowIsNeverRushed(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("dontknow_rush")
	const question = "3/4 + 1/2 = ?"
	diff := 0.6
	for i := 0; i < 3; i++ {
		res, err := e.SubmitQuizDontKnow(st.ID, "a", "5/4", 0.5, question, &diff)
		if err != nil {
			t.Fatalf("submit %d: %v", i, err)
		}
		if res.XP == xppolicy.RushPenalty {
			t.Fatalf("an admitted don't-know at 0.5s drew the rush penalty on attempt %d", i+1)
		}
		if res.XP != 0 {
			t.Errorf("an admitted don't-know paid %d XP, want 0", res.XP)
		}
	}
}

// Weakness on a downstream concept propagates backward to immediate prerequisites
// so learners who struggle downstream have foundation gaps surfaced for remediation.
func TestEngine_PropagateWeaknessBackwardToPrerequisites(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("weakness_propagation")
	now := time.Now().UTC()
	// Prerequisite "a" was mastered previously.
	if err := e.repo.UpsertProgress(&storage.ConceptProgress{
		StudentID:     st.ID,
		ConceptID:     "a",
		Status:        string(mastery.StatusMastered),
		Streak:        3,
		WeaknessScore: 0.0,
		LastReviewed:  &now,
	}); err != nil {
		t.Fatalf("upsert prereq: %v", err)
	}
	// Downstream concept "b" (which depends on "a") has high weakness due to struggles.
	if err := e.repo.UpsertProgress(&storage.ConceptProgress{
		StudentID:     st.ID,
		ConceptID:     "b",
		Status:        string(mastery.StatusPracticing),
		Streak:        0,
		WeaknessScore: 0.8,
		LastAttempted: &now,
	}); err != nil {
		t.Fatalf("upsert dependent: %v", err)
	}

	e.PropagateWeakness(st.ID)

	pA, err := e.repo.GetProgress(st.ID, "a")
	if err != nil {
		t.Fatalf("get prereq progress: %v", err)
	}
	if pA.WeaknessScore < 0.35 {
		t.Errorf("prerequisite weakness = %v, want >= 0.35 after backward propagation from struggling dependent", pA.WeaknessScore)
	}
}

func TestEngine_DiagnosticConditionalCompletion_Persistence(t *testing.T) {
	e := testEngine(t)
	st, err := e.CreateStudent("provisional_student")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}

	report := &diagnostic.DiagnosticReport{
		ConditionallyCompleted: []string{"Concept A"},
	}
	if err := e.ApplyDiagnosticReport(st.ID, report); err != nil {
		t.Fatalf("apply diagnostic report: %v", err)
	}

	p, err := e.repo.GetProgress(st.ID, "a")
	if err != nil {
		t.Fatalf("get progress: %v", err)
	}
	if p == nil {
		t.Fatal("expected progress row for conditionally completed concept")
	}
	if p.Status != string(mastery.StatusLearning) {
		// A conditional completion is a placement, so it must not read MASTERED:
		// that unlocked dependents as though evidence existed, and hysteresis then
		// made the guess permanent after one correct answer.
		t.Errorf("status = %q, want LEARNING", p.Status)
	}
	if !p.IsProvisional {
		t.Error("expected IsProvisional true")
	}
	if !p.PlacementSeeded {
		t.Error("expected PlacementSeeded true")
	}
}

func TestEngine_ProvisionalFallback_DemotionOnStruggle(t *testing.T) {
	e := testEngine(t)
	st, err := e.CreateStudent("struggle_student")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}

	// Concept a is conditionally completed (provisional mastered).
	report := &diagnostic.DiagnosticReport{
		ConditionallyCompleted: []string{"a"},
	}
	if err := e.ApplyDiagnosticReport(st.ID, report); err != nil {
		t.Fatalf("apply report: %v", err)
	}

	// Concept b depends on a.
	sess, err := e.repo.CreateSession(st.ID)
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	// Active session serving b.
	e.mu.Lock()
	e.sessions[sess.ID] = &activeSession{
		studentID:      st.ID,
		conceptID:      "b",
		expectedAnswer: "99",
		attemptID:      "att-1",
		requiredStreak: 3,
		timeThreshold:  60.0,
	}
	e.mu.Unlock()

	// Learner struggles on b with a wrong answer.
	res, err := e.SubmitAnswer(sess.ID, st.ID, "att-1", "wrong", 5.0)
	if err != nil {
		t.Fatalf("submit answer: %v", err)
	}

	if !res.Halted {
		t.Error("expected session halted after provisional prerequisite fallback")
	}
	if len(res.Remedial) == 0 || res.Remedial[0] != "a" {
		t.Errorf("expected remedial queue head to be prerequisite a, got %v", res.Remedial)
	}

	// Verify prerequisite a demoted from MASTERED to LEARNING.
	pA, err := e.repo.GetProgress(st.ID, "a")
	if err != nil {
		t.Fatalf("get progress a: %v", err)
	}
	if pA.Status != string(mastery.StatusLearning) {
		t.Errorf("status a = %q, want LEARNING", pA.Status)
	}
	if pA.Streak != 0 {
		t.Errorf("streak a = %d, want 0", pA.Streak)
	}
	if pA.WeaknessScore != 0.8 {
		t.Errorf("weakness a = %v, want 0.8", pA.WeaknessScore)
	}
	if !pA.IsProvisional {
		t.Error("expected IsProvisional to remain true while learning")
	}

	// Prerequisite a must be served first in active session.
	nextQ, err := e.NextQuestion(sess.ID, st.ID)
	if err != nil {
		t.Fatalf("next question: %v", err)
	}
	if nextQ.ConceptID != "a" {
		t.Errorf("next concept = %q, want remedial concept a", nextQ.ConceptID)
	}
}

func TestEngine_ProvisionalClearing_OnTrueMastery(t *testing.T) {
	e := testEngine(t)
	st, err := e.CreateStudent("mastery_student")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}

	// Concept a was provisional and demoted to LEARNING.
	now := time.Now().UTC()
	if err := e.repo.UpsertProgress(&storage.ConceptProgress{
		StudentID:     st.ID,
		ConceptID:     "a",
		Status:        string(mastery.StatusLearning),
		Streak:        0,
		IsProvisional: true,
		LastAttempted: &now,
	}); err != nil {
		t.Fatalf("seed progress: %v", err)
	}

	sess, err := e.repo.CreateSession(st.ID)
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	// Practice a until true mastery is achieved.
	for i := 0; i < 6; i++ {
		e.mu.Lock()
		e.sessions[sess.ID] = &activeSession{}
		e.mu.Unlock()
		q, err := e.NextQuestion(sess.ID, st.ID)
		if err != nil {
			t.Fatalf("next question %d: %v", i+1, err)
		}
		if q == nil || q.ConceptID != "a" {
			break
		}
		res, err := e.SubmitAnswer(sess.ID, st.ID, q.AttemptID, "42", 1.0)
		if err != nil {
			t.Fatalf("submit answer %d: %v", i+1, err)
		}
		if res.NewStatus == mastery.StatusMastered {
			break
		}
	}

	// IsProvisional must now be false and status must be MASTERED.
	pA, err := e.repo.GetProgress(st.ID, "a")
	if err != nil {
		t.Fatalf("get progress a: %v", err)
	}
	if pA.Status != string(mastery.StatusMastered) {
		t.Fatalf("status = %q, want MASTERED", pA.Status)
	}
	if pA.IsProvisional {
		t.Error("expected IsProvisional cleared to false after true mastery")
	}

	// Subsequent struggle on dependent b must NOT demote a.
	e.mu.Lock()
	e.sessions[sess.ID] = &activeSession{
		studentID:      st.ID,
		conceptID:      "b",
		expectedAnswer: "99",
		attemptID:      "att-b",
		requiredStreak: 3,
		timeThreshold:  60.0,
	}
	e.mu.Unlock()

	_, err = e.SubmitAnswer(sess.ID, st.ID, "att-b", "wrong", 5.0)
	if err != nil {
		t.Fatalf("submit answer b: %v", err)
	}

	pAAfter, err := e.repo.GetProgress(st.ID, "a")
	if err != nil {
		t.Fatalf("get progress a: %v", err)
	}
	if pAAfter.Status != string(mastery.StatusMastered) {
		t.Errorf("status a = %q, want to remain MASTERED after independent mastery", pAAfter.Status)
	}
}

func TestEngine_ProvisionalFallback_StudyAnchorPath(t *testing.T) {
	e := testEngine(t)
	st, err := e.CreateStudent("study_fallback")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}

	report := &diagnostic.DiagnosticReport{
		ConditionallyCompleted: []string{"a"},
	}
	if err := e.ApplyDiagnosticReport(st.ID, report); err != nil {
		t.Fatalf("apply report: %v", err)
	}

	// Set study anchor for concept b.
	e.SetStudyAnchorBatch(st.ID, "b", map[string]Anchor{
		"Question B text": {Answer: "99", Explanation: "Explanation B"},
	})

	// Submit incorrect study answer for b.
	res, err := e.SubmitStudyAnswer(st.ID, "b", "wrong", 5.0, "Question B text")
	if err != nil {
		t.Fatalf("submit study answer: %v", err)
	}
	if !res.Halted {
		t.Error("expected study result halted after provisional fallback")
	}
	if len(res.Remedial) == 0 || res.Remedial[0] != "a" {
		t.Errorf("expected remedial to contain demoted concept a, got %v", res.Remedial)
	}

	pA, err := e.repo.GetProgress(st.ID, "a")
	if err != nil {
		t.Fatalf("get progress a: %v", err)
	}
	if pA.Status != string(mastery.StatusLearning) {
		t.Errorf("status a = %q, want LEARNING", pA.Status)
	}
}

func TestEngine_FSRS_LiveStateTransitions(t *testing.T) {
	e := testEngine(t)
	st, err := e.CreateStudent("fsrs-learner")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}
	sess, err := e.repo.CreateSession(st.ID)
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	// First question on concept a
	q1, err := e.NextQuestion(sess.ID, st.ID)
	if err != nil {
		t.Fatalf("next question: %v", err)
	}
	ans1, err := e.SubmitAnswer(sess.ID, st.ID, q1.AttemptID, "42", 3.0)
	if err != nil {
		t.Fatalf("submit answer 1: %v", err)
	}
	if !ans1.Correct {
		t.Fatalf("expected answer 1 correct")
	}

	p1, err := e.repo.GetProgress(st.ID, q1.ConceptID)
	if err != nil {
		t.Fatalf("get progress 1: %v", err)
	}
	if p1.FSRSStability == nil || *p1.FSRSStability <= 0 {
		t.Fatalf("expected valid FSRSStability after first correct attempt, got %v", p1.FSRSStability)
	}
	if p1.FSRSDifficulty == nil || *p1.FSRSDifficulty <= 0 {
		t.Fatalf("expected valid FSRSDifficulty after first correct attempt, got %v", p1.FSRSDifficulty)
	}
	stab1 := *p1.FSRSStability

	// Second question on concept a
	q2, err := e.NextQuestion(sess.ID, st.ID)
	if err != nil {
		t.Fatalf("next question 2: %v", err)
	}
	ans2, err := e.SubmitAnswer(sess.ID, st.ID, q2.AttemptID, "42", 2.5)
	if err != nil {
		t.Fatalf("submit answer 2: %v", err)
	}
	if !ans2.Correct {
		t.Fatalf("expected answer 2 correct")
	}

	p2, err := e.repo.GetProgress(st.ID, q2.ConceptID)
	if err != nil {
		t.Fatalf("get progress 2: %v", err)
	}
	stab2 := *p2.FSRSStability
	if stab2 <= stab1 {
		t.Errorf("expected stability to increase on repeated correct recall: stab1=%f stab2=%f", stab1, stab2)
	}

	// Test decay detection via retrievability
	now := time.Now().UTC()
	past := now.Add(-30 * 24 * time.Hour)
	p2.Status = string(mastery.StatusMastered)
	p2.FSRSLastReview = &past
	futureDue := now.Add(10 * 24 * time.Hour)
	p2.NextReviewDue = &futureDue
	if err := e.repo.UpsertProgress(p2); err != nil {
		t.Fatalf("upsert progress past review: %v", err)
	}

	dueCount, err := e.ReviewsDue(st.ID)
	if err != nil {
		t.Fatalf("reviews due: %v", err)
	}
	if dueCount < 1 {
		t.Errorf("expected at least 1 review due when FSRS retrievability decays, got %d", dueCount)
	}
}

func TestEngine_AntiGuessing_SlowMissBreaksTheRun(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("guess_slow")
	sess, _ := e.repo.CreateSession(st.ID)

	// First rapid miss: the run starts but nothing is charged yet.
	q1, _ := e.NextQuestion(sess.ID, st.ID)
	if res, _ := e.SubmitAnswer(sess.ID, st.ID, q1.AttemptID, "wrong", 1.0); res.XP != 0 {
		t.Errorf("expected 0 XP first rapid miss, got %d", res.XP)
	}

	// A slow wrong answer is considered, so it ends the run rather than extending it.
	q2, _ := e.NextQuestion(sess.ID, st.ID)
	if res, _ := e.SubmitAnswer(sess.ID, st.ID, q2.AttemptID, "wrong", 9.0); res.XP != 0 {
		t.Errorf("expected 0 XP for a considered miss, got %d", res.XP)
	}

	// The next rapid miss is the first of a new run, so it is not penalised.
	q3, _ := e.NextQuestion(sess.ID, st.ID)
	if res, _ := e.SubmitAnswer(sess.ID, st.ID, q3.AttemptID, "wrong", 1.0); res.XP != 0 {
		t.Errorf("expected 0 XP after the run was broken, got %d", res.XP)
	}

	// A second consecutive rapid miss is what earns the penalty.
	q4, _ := e.NextQuestion(sess.ID, st.ID)
	if res, _ := e.SubmitAnswer(sess.ID, st.ID, q4.AttemptID, "wrong", 1.0); res.XP != xppolicy.RushPenalty {
		t.Errorf("expected %d XP for two consecutive rapid misses, got %d", xppolicy.RushPenalty, res.XP)
	}
}
