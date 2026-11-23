package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestGitHubStarsCachesTheFetch(t *testing.T) {
	var hits atomic.Int32
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		if got := r.Header.Get("User-Agent"); got != "mathua" {
			t.Errorf("user agent = %q, want mathua", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"stargazers_count":1234}`))
	}))
	defer upstream.Close()

	old := githubStarsEndpoint
	githubStarsEndpoint = upstream.URL
	defer func() { githubStarsEndpoint = old }()

	s := &Server{}
	for i := 0; i < 2; i++ {
		rec := httptest.NewRecorder()
		s.handleGitHubStars(rec, httptest.NewRequest(http.MethodGet, "/api/github-stars", nil))
		if rec.Code != http.StatusOK {
			t.Fatalf("request %d: status = %d, body = %s", i, rec.Code, rec.Body.String())
		}
		var body githubStarsResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("request %d: decode: %v", i, err)
		}
		if body.Count != 1234 {
			t.Fatalf("request %d: count = %d, want 1234", i, body.Count)
		}
	}
	if hits.Load() != 1 {
		t.Fatalf("upstream fetched %d times, want once because the answer is cached", hits.Load())
	}
}

func TestGitHubStarsReportsAnUpstreamFailureWhenNothingIsCached(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "rate limited", http.StatusForbidden)
	}))
	defer upstream.Close()

	old := githubStarsEndpoint
	githubStarsEndpoint = upstream.URL
	defer func() { githubStarsEndpoint = old }()

	rec := httptest.NewRecorder()
	(&Server{}).handleGitHubStars(rec, httptest.NewRequest(http.MethodGet, "/api/github-stars", nil))
	if rec.Code != http.StatusBadGateway {
		t.Fatalf("status = %d, want 502", rec.Code)
	}
}
