package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateReport_Guest(t *testing.T) {
	s := testServer(t)
	mux := http.NewServeMux()
	s.Register(mux)

	body, _ := json.Marshal(map[string]string{
		"concept_id": "a",
		"kind":       "question",
		"question":   "2+2=?",
		"expected":   "4",
		"reason":     "wrong_answer",
		"detail":     "looks off",
		"reporter_id": "guest_123",
	})
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("POST", "/api/reports", bytes.NewReader(body)))
	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var res map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if res["status"] != "open" {
		t.Errorf("expected open status, got %v", res)
	}
}

func TestCreateReport_InvalidReason(t *testing.T) {
	s := testServer(t)
	mux := http.NewServeMux()
	s.Register(mux)

	body, _ := json.Marshal(map[string]string{
		"concept_id": "a",
		"kind":       "question",
		"reason":     "bogus",
	})
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("POST", "/api/reports", bytes.NewReader(body)))
	if rec.Code != 400 {
		t.Errorf("expected 400, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestListReports_RequiresAdmin(t *testing.T) {
	s := testServer(t)
	mux := http.NewServeMux()
	s.Register(mux)

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("GET", "/api/reports", nil))
	if rec.Code != 404 {
		t.Errorf("expected 404 without ADMIN_TOKEN, got %d", rec.Code)
	}
}
