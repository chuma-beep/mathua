package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/engine"
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/storage"
)

type testGen struct{}

func (g *testGen) Generate(difficulty float64) generator.Problem {
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
