package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEdgeGuard(t *testing.T) {
	ok := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("disabled without a secret", func(t *testing.T) {
		t.Setenv("MATHUA_EDGE_SECRET", "")
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/next", nil)
		EdgeGuard(ok).ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rec.Code)
		}
	})

	t.Run("health stays open", func(t *testing.T) {
		t.Setenv("MATHUA_EDGE_SECRET", "edge-secret")
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
		EdgeGuard(ok).ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rec.Code)
		}
	})

	t.Run("API rejects a missing or wrong secret", func(t *testing.T) {
		t.Setenv("MATHUA_EDGE_SECRET", "edge-secret")
		for name, header := range map[string]string{
			"missing": "",
			"wrong":   "not-the-secret",
		} {
			t.Run(name, func(t *testing.T) {
				rec := httptest.NewRecorder()
				req := httptest.NewRequest(http.MethodGet, "/api/next", nil)
				if header != "" {
					req.Header.Set(edgeSecretHeader, header)
				}
				EdgeGuard(ok).ServeHTTP(rec, req)
				if rec.Code != http.StatusForbidden {
					t.Fatalf("status = %d, want 403", rec.Code)
				}
			})
		}
	})

	t.Run("API accepts the edge secret", func(t *testing.T) {
		t.Setenv("MATHUA_EDGE_SECRET", "edge-secret")
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/study/answer", nil)
		req.Header.Set(edgeSecretHeader, "edge-secret")
		EdgeGuard(ok).ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rec.Code)
		}
	})

	t.Run("static files are never gated", func(t *testing.T) {
		t.Setenv("MATHUA_EDGE_SECRET", "edge-secret")
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/learn", nil)
		EdgeGuard(ok).ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rec.Code)
		}
	})
}
