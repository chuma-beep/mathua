package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"testing"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/engine"
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/planning"
	"github.com/chuma-beep/mathua/internal/storage"
)

type testGen struct{}

func (g *testGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return generator.Problem{Question: "2+2=?", Answer: "4", Explanation: "2+2=4"}
}

func testServer(t *testing.T) *Server {
	t.Helper()
	d, _ := concepts.Build([]concepts.Concept{
		{ID: "a", Label: "A", Domain: "d", GradingType: "numeric", Prerequisites: []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60}},
	})
	store, _ := storage.NewSQLiteStore(":memory:")
	t.Cleanup(func() { store.Close() })
	reg := generator.NewRegistry()
	reg.Register("a", &testGen{})
	return New(engine.New(store, d, reg, nil, nil), store, nil)
}

func TestHealth(t *testing.T) {
	s := testServer(t)
	mux := http.NewServeMux()
	s.Register(mux)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("GET", "/api/health", nil))
	if rec.Code != 200 {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestSessionStart(t *testing.T) {
	s := testServer(t)
	mux := http.NewServeMux()
	s.Register(mux)

	body, _ := json.Marshal(startSessionReq{Name: "tester"})
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("POST", "/api/session", bytes.NewReader(body)))

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var res startSessionRes
	json.Unmarshal(rec.Body.Bytes(), &res)
	if res.StudentID == "" || res.SessionID == "" {
		t.Errorf("expected non-empty IDs, got student=%q session=%q", res.StudentID, res.SessionID)
	}
	if res.Question == nil {
		t.Fatal("expected question")
	}
}

func TestAnswer(t *testing.T) {
	s := testServer(t)
	mux := http.NewServeMux()
	s.Register(mux)

	body, _ := json.Marshal(startSessionReq{Name: "tester"})
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("POST", "/api/session", bytes.NewReader(body)))
	var startRes startSessionRes
	json.Unmarshal(rec.Body.Bytes(), &startRes)

	ansBody, _ := json.Marshal(answerReq{
		SessionID: startRes.SessionID,
		AttemptID: startRes.Question.AttemptID,
		Answer:    "4",
		Elapsed:   2.0,
	})
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("POST", "/api/answer", bytes.NewReader(ansBody)))

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var res answerRes
	json.Unmarshal(rec.Body.Bytes(), &res)
	if res.Result == nil || !res.Result.Correct {
		t.Errorf("expected correct, got %v", res.Result)
	}
}

func TestScores(t *testing.T) {
	s := testServer(t)
	mux := http.NewServeMux()
	s.Register(mux)

	body, _ := json.Marshal(startSessionReq{Name: "tester"})
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("POST", "/api/session", bytes.NewReader(body)))
	var startRes startSessionRes
	json.Unmarshal(rec.Body.Bytes(), &startRes)

	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("GET", "/api/scores/"+startRes.StudentID, nil))
	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestGraph(t *testing.T) {
	s := testServer(t)
	mux := http.NewServeMux()
	s.Register(mux)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("GET", "/api/graph", nil))
	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestLeaderboard(t *testing.T) {
	s := testServer(t)
	mux := http.NewServeMux()
	s.Register(mux)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("GET", "/api/leaderboard", nil))
	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

// Expression (SymPy) grading integration test

type expressionTestGen struct{}

func (g *expressionTestGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return generator.Problem{
		Question:    "Expand (x+1)^2",
		Answer:      "x^2+2x+1",
		Explanation: "FOIL: x^2 + x + x + 1 = x^2 + 2x + 1",
	}
}

func TestAnswer_ExpressionGrading_Equivalent(t *testing.T) {
	if _, err := exec.LookPath("python3"); err != nil {
		t.Skip("python3 not available")
	}
	d, _ := concepts.Build([]concepts.Concept{
		{ID: "expr_test", Label: "Expr Test", Domain: "d",
			GradingType:      "expression",
			Prerequisites:    []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60}},
	})
	store, _ := storage.NewSQLiteStore(":memory:")
	t.Cleanup(func() { store.Close() })
	reg := generator.NewRegistry()
	reg.Register("expr_test", &expressionTestGen{})
	s := New(engine.New(store, d, reg, nil, nil), store, nil)
	mux := http.NewServeMux()
	s.Register(mux)

	body, _ := json.Marshal(startSessionReq{Name: "tester"})
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("POST", "/api/session", bytes.NewReader(body)))
	var startRes startSessionRes
	json.Unmarshal(rec.Body.Bytes(), &startRes)

	// Submit an equivalent expression
	ansBody, _ := json.Marshal(answerReq{
		SessionID: startRes.SessionID,
		AttemptID: startRes.Question.AttemptID,
		Answer:    "(x+1)^2",
		Elapsed:   5.0,
	})
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("POST", "/api/answer", bytes.NewReader(ansBody)))

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var res answerRes
	json.Unmarshal(rec.Body.Bytes(), &res)
	if res.Result == nil || !res.Result.Correct {
		t.Errorf("expected correct for equivalent expression, got %v", res.Result)
	}
}

func TestAnswer_ExpressionGrading_Nonequivalent(t *testing.T) {
	if _, err := exec.LookPath("python3"); err != nil {
		t.Skip("python3 not available")
	}
	d, _ := concepts.Build([]concepts.Concept{
		{ID: "expr_test", Label: "Expr Test", Domain: "d",
			GradingType:      "expression",
			Prerequisites:    []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60}},
	})
	store, _ := storage.NewSQLiteStore(":memory:")
	t.Cleanup(func() { store.Close() })
	reg := generator.NewRegistry()
	reg.Register("expr_test", &expressionTestGen{})
	s := New(engine.New(store, d, reg, nil, nil), store, nil)
	mux := http.NewServeMux()
	s.Register(mux)

	body, _ := json.Marshal(startSessionReq{Name: "tester"})
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("POST", "/api/session", bytes.NewReader(body)))
	var startRes startSessionRes
	json.Unmarshal(rec.Body.Bytes(), &startRes)

	// Submit a non-equivalent expression
	ansBody, _ := json.Marshal(answerReq{
		SessionID: startRes.SessionID,
		AttemptID: startRes.Question.AttemptID,
		Answer:    "x^2+3x+1",
		Elapsed:   5.0,
	})
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("POST", "/api/answer", bytes.NewReader(ansBody)))

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var res answerRes
	json.Unmarshal(rec.Body.Bytes(), &res)
	if res.Result == nil || res.Result.Correct {
		t.Errorf("expected incorrect for non-equivalent expression, got %v", res.Result)
	}
}

func TestCORS(t *testing.T) {
	s := testServer(t)
	mux := http.NewServeMux()
	s.Register(mux)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("OPTIONS", "/api/health", nil))
	if rec.Code != 204 {
		t.Errorf("expected 204 for OPTIONS, got %d", rec.Code)
	}
}

func TestGoalDiagnosticResume(t *testing.T) {
	d, _ := concepts.Build([]concepts.Concept{
		{ID: "a", Label: "A", Domain: "d", GradingType: "numeric", Prerequisites: []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60}},
	})
	store, _ := storage.NewSQLiteStore(":memory:")
	t.Cleanup(func() { store.Close() })
	reg := generator.NewRegistry()
	reg.Register("a", &testGen{})
	s := New(engine.New(store, d, reg, nil, planning.New(d)), store, nil)
	mux := http.NewServeMux()
	s.Register(mux)

	startBody, _ := json.Marshal(map[string]interface{}{
		"name": "tester", "concept_ids": []string{"a"},
	})
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("POST", "/api/goal/diagnostic", bytes.NewReader(startBody)))
	if rec.Code != 200 {
		t.Fatalf("start: expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var start map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &start)
	sid, _ := start["session_id"].(string)
	q0, _ := start["question"].(string)
	if sid == "" || q0 == "" {
		t.Fatalf("expected session_id + question, got %v", start)
	}

	// Pause, then resume: same live question, progress intact.
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("GET", "/api/goal/diagnostic/resume?session_id="+sid, nil))
	if rec.Code != 200 {
		t.Fatalf("resume: expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var resumed map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &resumed)
	if resumed["question"] != q0 {
		t.Errorf("expected resumed question %q, got %q", q0, resumed["question"])
	}
	if _, ok := resumed["progress"]; !ok {
		t.Error("expected progress in resume response")
	}

	// Answer, then resume again: advances to the next question, no loss.
	ansBody, _ := json.Marshal(map[string]interface{}{
		"session_id": sid, "concept_id": "a", "answer": "4", "elapsed": 5.0,
	})
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("POST", "/api/goal/diagnostic/answer", bytes.NewReader(ansBody)))
	if rec.Code != 200 {
		t.Fatalf("answer: expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var ans map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &ans)
	if done, _ := ans["done"].(bool); done {
		t.Skip("single-concept diagnostic completed after 1 answer; resume trivially done")
	}
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("GET", "/api/goal/diagnostic/resume?session_id="+sid, nil))
	if rec.Code != 200 {
		t.Fatalf("resume2: expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var resumed2 map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &resumed2)
	if resumed2["question"] != ans["question"] {
		t.Errorf("expected resumed question to match latest answer question")
	}

	// Unknown session → 404.
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("GET", "/api/goal/diagnostic/resume?session_id=nope", nil))
	if rec.Code != 404 {
		t.Errorf("expected 404 for unknown session, got %d", rec.Code)
	}
}
