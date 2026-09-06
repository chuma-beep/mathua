package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestRateLimiter_XForwardedForFirstIP(t *testing.T) {
	rl := newRateLimiter(100, 100, time.Minute)
	h := rl.middleware(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "10.0.0.1:1234"
	req.Header.Set("X-Forwarded-For", "203.0.113.7, 70.41.3.18, 150.172.238.4")
	rec := httptest.NewRecorder()
	h(rec, req)
	if rec.Code != 200 {
		t.Fatalf("expected pass-through, got %d", rec.Code)
	}
	if len(rl.visitors) != 1 {
		t.Fatalf("expected single bucket, got %d", len(rl.visitors))
	}
	for ip := range rl.visitors {
		if ip != "203.0.113.7" {
			t.Errorf("expected first-IP bucket 203.0.113.7, got %q", ip)
		}
	}
}

func TestRateLimiter_Friendly429(t *testing.T) {
	rl := newRateLimiter(100, 1, time.Hour)
	h := rl.middleware(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
	mkReq := func() *http.Request {
		req := httptest.NewRequest("GET", "/", nil)
		req.RemoteAddr = "10.9.9.9:1"
		return req
	}
	rec := httptest.NewRecorder()
	h(rec, mkReq())
	if rec.Code != 200 {
		t.Fatalf("first request should pass, got %d", rec.Code)
	}
	rec = httptest.NewRecorder()
	h(rec, mkReq())
	if rec.Code != 429 {
		t.Fatalf("expected 429, got %d", rec.Code)
	}
	body := rec.Body.String()
	if body == "" || strings.Contains(body, "rate limit exceeded") {
		t.Errorf("expected friendly 429 copy, got %q", body)
	}
}
