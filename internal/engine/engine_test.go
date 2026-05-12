package engine

import (
	"testing"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/storage"
)

type testGen struct {
	answer string
}

func (g *testGen) Generate(difficulty float64) generator.Problem {
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
	return New(store, d, reg)
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
	_, _ = e.NextQuestion(sess.ID, st.ID)
	result, err := e.SubmitAnswer(sess.ID, st.ID, "42", 3.0)
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
	_, _ = e.NextQuestion(sess.ID, st.ID)
	result, err := e.SubmitAnswer(sess.ID, st.ID, "99", 5.0)
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
		_, _ = e.NextQuestion(sess.ID, st.ID)
		result, err := e.SubmitAnswer(sess.ID, st.ID, "42", 1.0)
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
