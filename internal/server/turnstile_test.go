package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chuma-beep/mathua/internal/auth"
	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/engine"
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/storage"
)

func turnstileFixture(t *testing.T) *Server {
	t.Helper()
	fake := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "bad form", 400)
			return
		}
		if r.Form.Get("secret") != "server-secret" || r.Form.Get("response") != "good-token" {
			_, _ = w.Write([]byte(`{"success":false}`))
			return
		}
		_, _ = w.Write([]byte(`{"success":true}`))
	}))
	t.Cleanup(fake.Close)

	d, _ := concepts.Build([]concepts.Concept{
		{ID: "a", Label: "A", Domain: "d", GradingType: "numeric", Prerequisites: []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60}},
	})
	store, _ := storage.NewSQLiteStore(":memory:")
	t.Cleanup(func() { store.Close() })
	reg := generator.NewRegistry()
	reg.Register("a", &testGen{})
	s := New(engine.New(store, d, reg, nil, nil), store, auth.New(store))
	s.turnstileSecret = "server-secret"
	s.turnstileURL = fake.URL
	s.turnstileClient = fake.Client()
	return s
}

func TestVerifyTurnstileDisabledWithoutSecret(t *testing.T) {
	s := &Server{}
	if !s.verifyTurnstile("") {
		t.Fatal("an empty secret must disable the check")
	}
}

func TestVerifyTurnstileFailsClosed(t *testing.T) {
	s := turnstileFixture(t)
	cases := []struct {
		name  string
		token string
		want  bool
	}{
		{"empty token", "", false},
		{"wrong token", "bad-token", false},
		{"valid token", "good-token", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := s.verifyTurnstile(tc.token); got != tc.want {
				t.Fatalf("verifyTurnstile(%q) = %v, want %v", tc.token, got, tc.want)
			}
		})
	}
}

func TestVerifyTurnstileRejectsTransportFailure(t *testing.T) {
	s := turnstileFixture(t)
	s.turnstileURL = "http://127.0.0.1:1/siteverify"
	if s.verifyTurnstile("good-token") {
		t.Fatal("a failed verification request must not pass")
	}
}

func TestAuthHandlersRequireTurnstile(t *testing.T) {
	s := turnstileFixture(t)
	mux := http.NewServeMux()
	s.Register(mux)
	post := func(path, body string) *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader([]byte(body)))
		mux.ServeHTTP(rec, req)
		return rec
	}

	signup := `{"name":"Ada","username":"ada","password":"Engine!n1","email":"ada@example.com"`
	if rec := post("/api/auth/signup", signup+`}`); rec.Code != 403 {
		t.Fatalf("signup without a token: got %d, want 403", rec.Code)
	}
	if rec := post("/api/auth/signup", signup+`,"turnstile_token":"bad-token"}`); rec.Code != 403 {
		t.Fatalf("signup with a bad token: got %d, want 403", rec.Code)
	}
	if rec := post("/api/auth/signup", signup+`,"turnstile_token":"good-token"}`); rec.Code != 200 {
		t.Fatalf("signup with a good token: got %d: %s", rec.Code, rec.Body.String())
	}
	if rec := post("/api/auth/login", `{"username":"ada","password":"Engine!n1"}`); rec.Code != 403 {
		t.Fatalf("login without a token: got %d, want 403", rec.Code)
	}
	if rec := post("/api/auth/login", `{"username":"ada","password":"Engine!n1","turnstile_token":"good-token"}`); rec.Code != 200 {
		t.Fatalf("login with a good token: got %d: %s", rec.Code, rec.Body.String())
	}
}
