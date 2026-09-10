package engine

import (
	"errors"
	"sync"
	"testing"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/storage"
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
	progress, _ := e.GetProgress(st.ID)
	// 3 correct fast answers: UNSEEN → LEARNING (streak hits 3)
	if p, ok := progress["a"]; !ok || p.Status != "LEARNING" {
		t.Errorf("expected LEARNING after 3-streak, got %+v", progress["a"])
	}
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

// A duplicate/stale submission racing the real one must be rejected, never
// graded against a question the client was not shown (cross-question grading).
func TestEngine_SubmitAnswer_ConcurrentDuplicates(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("grace")
	sess, _ := e.repo.CreateSession(st.ID)
	q, _ := e.NextQuestion(sess.ID, st.ID)

	const workers = 8
	var wg sync.WaitGroup
	results := make([]error, workers)
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, results[i] = e.SubmitAnswer(sess.ID, st.ID, q.AttemptID, "42", 3.0)
		}(i)
	}
	wg.Wait()

	var ok, rejected int
	for _, err := range results {
		switch {
		case err == nil:
			ok++
		case errors.Is(err, ErrNoActiveQuestion):
			rejected++
		default:
			t.Fatalf("unexpected error: %v", err)
		}
	}
	if ok != 1 {
		t.Errorf("expected exactly 1 successful submission, got %d (rejected %d)", ok, rejected)
	}

	// The winning grade must match the question that was displayed (concept "a", answer "42").
	progress, _ := e.GetProgress(st.ID)
	if p, ok := progress["a"]; !ok || p.Streak != 1 {
		t.Errorf("expected concept a to be graded once with streak 1, got %+v", progress["a"])
	}
}

// A submission carrying the token of a previous question must be rejected, never
// graded against the session's current question (stale/cross-question grading).
func TestEngine_SubmitAnswer_StaleAttemptIDRejected(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("stale")
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

	// A stale submission for Q1 must be rejected — never graded against Q2.
	if _, err := e.SubmitAnswer(sess.ID, st.ID, q1.AttemptID, "99", 3.0); !errors.Is(err, ErrNoActiveQuestion) {
		t.Fatalf("expected stale submission rejected, got %v", err)
	}
	progress, _ := e.GetProgress(st.ID)
	if p := progress["a"]; p == nil || p.Attempts != 1 || p.Streak != 1 {
		t.Fatalf("stale submission must not be graded: got %+v", progress["a"])
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
	e.CreateStudent("alice")
	e.CreateStudent("bob")
	entries, err := e.GetLeaderboard()
	if err != nil {
		t.Fatalf("get leaderboard: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
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

	// Master concept "a" (requiredStreak=3 per level, 9 total correct answers).
	// Then verify concept "b" becomes available (prerequisite "a" now mastered).
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
			t.Fatalf("expected concept a on attempt %d, got %s", i+1, q.ConceptID)
		}
		res, err := e.SubmitAnswer(sess.ID, st.ID, q.AttemptID, "42", 1.0)
		if err != nil {
			t.Fatalf("submit answer %d: %v", i+1, err)
		}
		if !res.Correct {
			t.Fatalf("expected correct on attempt %d, got %s", i+1, res.Feedback)
		}
	}
	progress, _ := e.GetProgress(st.ID)
	if p := progress["a"]; p == nil || p.Status != "MASTERED" {
		t.Fatalf("expected MASTERED after 9 correct, got %+v", progress["a"])
	}

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
		t.Errorf("expected -5 XP second rush, got %d", res2.XP)
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
	// First incorrect rush → 0; second incorrect rush → -5
	r1, err := e.SubmitStudyAnswer(st.ID, "a", "wrong", "42", 1.0)
	if err != nil {
		t.Fatalf("study submit 1: %v", err)
	}
	if r1.XP != 0 {
		t.Errorf("expected 0 XP first study rush, got %d", r1.XP)
	}
	r2, err := e.SubmitStudyAnswer(st.ID, "a", "wrong", "42", 1.0)
	if err != nil {
		t.Fatalf("study submit 2: %v", err)
	}
	if r2.XP != -5 {
		t.Errorf("expected -5 XP second study rush, got %d", r2.XP)
	}
	if !r2.Halted {
		t.Error("expected halted flag after 2 consecutive study misses")
	}
}

// Session 5: efficacy instrumentation — first-pass / second-pass rates.
func TestEngine_Efficacy(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("efficacy")
	// Concept a: first attempt correct (first-pass hit), concept b: wrong then right.
	if _, err := e.SubmitStudyAnswer(st.ID, "a", "42", "42", 5.0); err != nil {
		t.Fatalf("submit a: %v", err)
	}
	if _, err := e.SubmitStudyAnswer(st.ID, "b", "1", "99", 5.0); err != nil {
		t.Fatalf("submit b wrong: %v", err)
	}
	if _, err := e.SubmitStudyAnswer(st.ID, "b", "99", "99", 5.0); err != nil {
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
	// agg1: a correct first try; b wrong then right.
	_, _ = e.SubmitStudyAnswer(st1.ID, "a", "42", "42", 5.0)
	_, _ = e.SubmitStudyAnswer(st1.ID, "b", "1", "99", 5.0)
	_, _ = e.SubmitStudyAnswer(st1.ID, "b", "99", "99", 5.0)
	// agg2: a wrong then right; b wrong then right.
	_, _ = e.SubmitStudyAnswer(st2.ID, "a", "1", "42", 5.0)
	_, _ = e.SubmitStudyAnswer(st2.ID, "a", "42", "42", 5.0)
	_, _ = e.SubmitStudyAnswer(st2.ID, "b", "1", "99", 5.0)
	_, _ = e.SubmitStudyAnswer(st2.ID, "b", "99", "99", 5.0)

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
	res, err := e.SubmitQuizAnswer(st.ID, "a", "42", "42", 3.0)
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
	lesson := computeXPForTask(true, 3.0, 60, res.Streak, TaskLesson)
	if res.XP <= lesson {
		t.Errorf("expected quiz XP %d to exceed lesson XP %d", res.XP, lesson)
	}
}

func TestEngine_SubmitStudyAnswer_UnknownConcept(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("lost")
	if _, err := e.SubmitStudyAnswer(st.ID, "nope.not.real", "1", "1", 5.0); !errors.Is(err, ErrUnknownConcept) {
		t.Errorf("expected ErrUnknownConcept, got %v", err)
	}
	if _, err := e.SubmitQuizAnswer(st.ID, "nope.not.real", "1", "1", 5.0); !errors.Is(err, ErrUnknownConcept) {
		t.Errorf("expected ErrUnknownConcept (quiz), got %v", err)
	}
	// No garbage progress row persisted.
	if p, _ := e.GetProgress(st.ID); len(p) != 0 {
		t.Errorf("expected no progress rows for unknown concept, got %v", p)
	}
}

func TestEngine_StudyExpected_SurvivesRestart(t *testing.T) {
	e1 := testEngine(t)
	// Reach into the shared :memory: store via a second engine instance.
	e1.SetStudyExpected("stu1", "a", "42")
	// Simulated restart: fresh engine, empty memory map, same store.
	e2 := New(e1.repo, e1.dag, e1.registry, nil, nil)
	if v, ok := e2.popStudyExpected("stu1", "a"); !ok || v != "42" {
		t.Fatalf("expected durable expected=42, got %q ok=%v", v, ok)
	}
	// Single-use: second pop misses everywhere.
	if _, ok := e2.popStudyExpected("stu1", "a"); ok {
		t.Error("expected consumed anchor to miss on second pop")
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

	res, err := e.SubmitQuizAnswer(st.ID, "b", "wrong", "99", 5.0)
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
	if _, err := e.SubmitQuizAnswer(st.ID, "a", "42", "42", 5.0); err != nil {
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

// Batch 1: 150 XP gate — due at 150 since last completion, reset on record.

func TestEngine_QuizGate_DueAndReset(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("gate")
	due, err := e.QuizDue(st.ID)
	if err != nil || due {
		t.Fatalf("expected not due at 0 XP, due=%v err=%v", due, err)
	}
	if err := e.repo.AddXP(st.ID, 149); err != nil {
		t.Fatalf("add xp: %v", err)
	}
	if due, _ := e.QuizDue(st.ID); due {
		t.Error("expected not due at 149 XP since quiz")
	}
	if err := e.repo.AddXP(st.ID, 1); err != nil {
		t.Fatalf("add xp: %v", err)
	}
	if since, _ := e.QuizXPSince(st.ID); since != 150 {
		t.Errorf("expected 150 since, got %d", since)
	}
	if due, _ := e.QuizDue(st.ID); !due {
		t.Error("expected due at 150 XP since quiz")
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
