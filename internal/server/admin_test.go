package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/engine"
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/storage"
)

func adminMux(t *testing.T) (*Server, *http.ServeMux) {
	t.Helper()
	s := testServer(t)
	mux := http.NewServeMux()
	s.Register(mux)
	return s, mux
}

func loginAdmin(t *testing.T, mux *http.ServeMux, password string) (code int, token string) {
	t.Helper()
	body, _ := json.Marshal(map[string]string{"password": password})
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("POST", "/api/admin/login", bytes.NewReader(body)))
	if rec.Code != 200 {
		return rec.Code, ""
	}
	var res struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("decode login: %v", err)
	}
	return rec.Code, res.Token
}

func TestAdminLogin_Success(t *testing.T) {
	t.Setenv("ADMIN_PASSWORD", "correct-horse-secret")
	_, mux := adminMux(t)

	code, token := loginAdmin(t, mux, "correct-horse-secret")
	if code != 200 || token == "" {
		t.Fatalf("expected 200 with token, got %d token=%q", code, token)
	}
}

func TestAdminLogin_WrongPassword(t *testing.T) {
	t.Setenv("ADMIN_PASSWORD", "correct-horse-secret")
	_, mux := adminMux(t)

	code, _ := loginAdmin(t, mux, "wrong")
	if code != 401 {
		t.Errorf("expected 401, got %d", code)
	}
}

func TestAdminLogin_Unconfigured(t *testing.T) {
	t.Setenv("ADMIN_PASSWORD", "")
	_, mux := adminMux(t)

	body, _ := json.Marshal(map[string]string{"password": "anything"})
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("POST", "/api/admin/login", bytes.NewReader(body)))
	if rec.Code != 404 {
		t.Errorf("expected 404 when ADMIN_PASSWORD unset, got %d", rec.Code)
	}
}

func TestAdminSession_GrantsListAccess(t *testing.T) {
	t.Setenv("ADMIN_PASSWORD", "correct-horse-secret")
	_, mux := adminMux(t)

	_, token := loginAdmin(t, mux, "correct-horse-secret")
	req := httptest.NewRequest("GET", "/api/reports?status=open", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Errorf("expected 200 with session token, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestAdminLogout_InvalidatesSession(t *testing.T) {
	t.Setenv("ADMIN_PASSWORD", "correct-horse-secret")
	_, mux := adminMux(t)

	_, token := loginAdmin(t, mux, "correct-horse-secret")

	req := httptest.NewRequest("POST", "/api/admin/logout", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatalf("expected 200 logout, got %d", rec.Code)
	}

	req = httptest.NewRequest("GET", "/api/reports?status=open", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != 404 {
		t.Errorf("expected 404 after logout, got %d", rec.Code)
	}
}

func TestAdminSession_ExpiredRejected(t *testing.T) {
	t.Setenv("ADMIN_PASSWORD", "correct-horse-secret")
	s, mux := adminMux(t)

	_, token := loginAdmin(t, mux, "correct-horse-secret")
	// Backdate the durable row past expiry (restart-proof path, not memory).
	exp := time.Now().Add(-time.Minute).UTC().Format(time.RFC3339)
	if err := s.repo.UpsertServerSession(serverSessionAdmin, token, "", exp); err != nil {
		t.Fatalf("backdate session: %v", err)
	}

	req := httptest.NewRequest("GET", "/api/reports?status=open", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != 404 {
		t.Errorf("expected 404 for expired session, got %d", rec.Code)
	}
}

func TestAdminSession_SurvivesRestart(t *testing.T) {
	t.Setenv("ADMIN_PASSWORD", "correct-horse-secret")
	d, _ := concepts.Build([]concepts.Concept{
		{ID: "a", Label: "A", Domain: "d", GradingType: "numeric", Prerequisites: []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60}},
	})
	store, _ := storage.NewSQLiteStore(":memory:")
	t.Cleanup(func() { store.Close() })
	reg := generator.NewRegistry()
	reg.Register("a", &testGen{})
	// Two Server instances over one store = simulated restart. No
	// process-local session state may remain for this to pass.
	s1 := New(engine.New(store, d, reg, nil, nil), store, nil)
	mux1 := http.NewServeMux()
	s1.Register(mux1)
	s2 := New(engine.New(store, d, reg, nil, nil), store, nil)
	mux2 := http.NewServeMux()
	s2.Register(mux2)

	_, token := loginAdmin(t, mux1, "correct-horse-secret")
	if token == "" {
		t.Fatal("expected login token")
	}
	req := httptest.NewRequest("GET", "/api/reports?status=open", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	mux2.ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Errorf("expected 200 on fresh instance, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestAdminLegacyToken_StillAccepted(t *testing.T) {
	t.Setenv("ADMIN_TOKEN", "legacy-secret")
	_, mux := adminMux(t)

	req := httptest.NewRequest("GET", "/api/reports?status=open", nil)
	req.Header.Set("Authorization", "Bearer legacy-secret")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Errorf("expected 200 with legacy ADMIN_TOKEN, got %d", rec.Code)
	}
}
