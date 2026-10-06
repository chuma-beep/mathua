package server

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/auth"
	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/engine"
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/lessons"
	"github.com/chuma-beep/mathua/internal/planning"
	"github.com/chuma-beep/mathua/internal/storage"
)

type testGen struct{}

func (g *testGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return generator.Problem{Question: "2+2=?", Answer: "4", Explanation: "2+2=4"}
}

// seededGen returns a different question on every call. BatchGenerateContext
// invokes Generate repeatedly with one context and dedupes by question text, so
// varying only on the seed would collapse a whole batch to a single question.
// The real generators vary per call for the same reason, and that is what the
// Learn client's exclude[] and its fresh-variant prefetch depend on.
type seededGen struct{ n int64 }

func (g *seededGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	n := atomic.AddInt64(&g.n, 1) + ctx.Seed
	return generator.Problem{
		Question: fmt.Sprintf("seeded %d + 1 = ?", n),
		Answer:   strconv.FormatInt(n+1, 10),
	}
}

// varyingGuestServer is guestServer with a generator whose questions vary.
func varyingGuestServer(t *testing.T) (*Server, *http.ServeMux, storage.Repository) {
	t.Helper()
	d, _ := concepts.Build([]concepts.Concept{
		{ID: "a", Label: "A", Domain: "d", GradingType: "numeric", Prerequisites: []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60}},
	})
	store, _ := storage.NewSQLiteStore(":memory:")
	t.Cleanup(func() { store.Close() })
	reg := generator.NewRegistry()
	reg.Register("a", &seededGen{})
	s := New(engine.New(store, d, reg, nil, nil), store, auth.New(store))
	mux := http.NewServeMux()
	s.Register(mux)
	return s, mux, store
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

func TestScores(t *testing.T) {
	s := testServer(t)
	mux := http.NewServeMux()
	s.Register(mux)

	rec := postGuest(t, mux, `{}`, "")
	if rec.Code != 200 {
		t.Fatalf("guest: %d %s", rec.Code, rec.Body.String())
	}
	var guest struct {
		StudentID string `json:"student_id"`
	}
	json.Unmarshal(rec.Body.Bytes(), &guest)

	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("GET", "/api/scores/"+guest.StudentID, nil))
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

func TestMeRoutes(t *testing.T) {
	d, _ := concepts.Build([]concepts.Concept{
		{ID: "a", Label: "A", Domain: "d", GradingType: "numeric", Prerequisites: []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60}},
	})
	store, _ := storage.NewSQLiteStore(":memory:")
	t.Cleanup(func() { store.Close() })
	reg := generator.NewRegistry()
	reg.Register("a", &testGen{})
	authSvc := auth.New(store)
	token, st, err := authSvc.Signup("Ada", "ada", "Engine!n1")
	if err != nil || token == "" || st == nil {
		t.Fatalf("signup: %v", err)
	}
	s := New(engine.New(store, d, reg, nil, nil), store, authSvc)
	mux := http.NewServeMux()
	s.Register(mux)

	for _, path := range []string{"/api/auth/me", "/api/me"} {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", path, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		mux.ServeHTTP(rec, req)
		if rec.Code != 200 {
			t.Errorf("%s: expected 200, got %d: %s", path, rec.Code, rec.Body.String())
			continue
		}
		var me map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &me)
		if me["student_id"] != st.ID {
			t.Errorf("%s: expected student_id %q, got %v", path, st.ID, me["student_id"])
		}
	}

	// Dead token (rotated secret / garbage) → 401 so the client logs out.
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/auth/me", nil)
	req.Header.Set("Authorization", "Bearer dead.token.here")
	mux.ServeHTTP(rec, req)
	if rec.Code != 401 {
		t.Errorf("expected 401 for dead token, got %d", rec.Code)
	}
}

func TestAuthLoginSignup(t *testing.T) {
	d, _ := concepts.Build([]concepts.Concept{
		{ID: "a", Label: "A", Domain: "d", GradingType: "numeric", Prerequisites: []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60}},
	})
	store, _ := storage.NewSQLiteStore(":memory:")
	t.Cleanup(func() { store.Close() })
	reg := generator.NewRegistry()
	reg.Register("a", &testGen{})
	authSvc := auth.New(store)
	s := New(engine.New(store, d, reg, nil, nil), store, authSvc)
	mux := http.NewServeMux()
	s.Register(mux)
	post := func(path, body, token string) *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", path, bytes.NewReader([]byte(body)))
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		mux.ServeHTTP(rec, req)
		return rec
	}

	// Signup normalizes case; duplicate (any case) → friendly 409.
	if rec := post("/api/auth/signup", `{"name":"Ada","username":"Ada","password":"Engine!n1","email":"ada@example.com"}`, ""); rec.Code != 200 {
		t.Fatalf("signup: expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	rec := post("/api/auth/signup", `{"name":"Other","username":"ADA","password":"Engine!n2","email":"other@example.com"}`, "")
	if rec.Code != 409 {
		t.Errorf("duplicate: expected 409, got %d: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "username is taken") {
		t.Errorf("expected friendly taken message, got %s", rec.Body.String())
	}
	// Bad username / weak password → 400, never a raw dump.
	if rec := post("/api/auth/signup", `{"name":"X","username":"ab","password":"Engine!n1","email":"x@example.com"}`, ""); rec.Code != 400 {
		t.Errorf("short username: expected 400, got %d", rec.Code)
	}
	// Login works case-insensitively.
	if rec := post("/api/auth/login", `{"username":"ADA","password":"Engine!n1"}`, ""); rec.Code != 200 {
		t.Errorf("login: expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	// Wrong password → generic 401 (no enumeration).
	if rec := post("/api/auth/login", `{"username":"ada","password":"Wrong!n9"}`, ""); rec.Code != 401 {
		t.Errorf("wrong password: expected 401, got %d", rec.Code)
	}
	// Google-only account → directed message, not generic invalid.
	if _, err := store.CreateGoogleUser("Gigi", "gigi@example.com", "gid-auth-test", ""); err != nil {
		t.Fatalf("create google user: %v", err)
	}
	rec = post("/api/auth/login", `{"username":"gigi@example.com","password":"Whatever!n1"}`, "")
	if rec.Code != 401 {
		t.Fatalf("google-only: expected 401, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Google sign-in") {
		t.Errorf("expected Google sign-in direction, got %s", rec.Body.String())
	}
	// Verified email already on a live account → friendly 409, no fork.
	rec = post("/api/auth/signup", `{"name":"Copy Cat","username":"copycat","password":"Engine!n1","email":"gigi@example.com"}`, "")
	if rec.Code != 409 {
		t.Errorf("verified email taken: expected 409, got %d: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "email already in use") {
		t.Errorf("expected email-taken message, got %s", rec.Body.String())
	}
}

// Separate server (fresh rate-limiter bucket) for the signup email rules:
// required, valid, lowercased at rest, strict one-email-one-account.
func TestAuthSignupEmail(t *testing.T) {
	d, _ := concepts.Build([]concepts.Concept{
		{ID: "a", Label: "A", Domain: "d", GradingType: "numeric", Prerequisites: []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60}},
	})
	store, _ := storage.NewSQLiteStore(":memory:")
	t.Cleanup(func() { store.Close() })
	reg := generator.NewRegistry()
	reg.Register("a", &testGen{})
	authSvc := auth.New(store)
	s := New(engine.New(store, d, reg, nil, nil), store, authSvc)
	mux := http.NewServeMux()
	s.Register(mux)
	post := func(path, body string) *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, httptest.NewRequest("POST", path, bytes.NewReader([]byte(body))))
		return rec
	}

	if rec := post("/api/auth/signup", `{"name":"No Mail","username":"nomail","password":"Engine!n1"}`); rec.Code != 400 {
		t.Errorf("missing email: expected 400, got %d: %s", rec.Code, rec.Body.String())
	} else if !strings.Contains(rec.Body.String(), "email is required") {
		t.Errorf("expected email-required message, got %s", rec.Body.String())
	}
	if rec := post("/api/auth/signup", `{"name":"Bad Mail","username":"badmail","password":"Engine!n1","email":"not-an-email"}`); rec.Code != 400 {
		t.Errorf("invalid email: expected 400, got %d: %s", rec.Code, rec.Body.String())
	}
	if rec := post("/api/auth/signup", `{"name":"Case Mail","username":"casemail","password":"Engine!n1","email":"Case@Example.COM"}`); rec.Code != 200 {
		t.Fatalf("case email signup: expected 200, got %d: %s", rec.Code, rec.Body.String())
	} else if st, err := store.FindByEmail("case@example.com"); err != nil || st == nil {
		t.Errorf("expected lowercased email stored, err=%v", err)
	}
	if rec := post("/api/auth/signup", `{"name":"Dupe One","username":"dupeone","password":"Engine!n1","email":"dupe@example.com"}`); rec.Code != 200 {
		t.Fatalf("dupe one: expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	if rec := post("/api/auth/signup", `{"name":"Dupe Two","username":"dupetwo","password":"Engine!n1","email":"dupe@example.com"}`); rec.Code != 409 {
		t.Errorf("unverified dupe: expected 409, got %d: %s", rec.Code, rec.Body.String())
	} else if !strings.Contains(rec.Body.String(), "email already in use") {
		t.Errorf("expected email-taken message, got %s", rec.Body.String())
	}
}

func TestFrontendRedirect(t *testing.T) {
	// Unset/empty base → relative redirect (single-binary unchanged).
	rel := frontendRedirect("", "/profile", "tok", "Ada", "s1")
	if !strings.HasPrefix(rel, "/profile?") {
		t.Errorf("expected relative /profile redirect, got %q", rel)
	}
	if !strings.Contains(rel, "token=tok") || !strings.Contains(rel, "id=s1") {
		t.Errorf("expected token+id query, got %q", rel)
	}
	// Split dev → absolute URL on the frontend origin.
	abs := frontendRedirect("http://localhost:3000", "/profile", "tok", "Ada", "s1")
	if !strings.HasPrefix(abs, "http://localhost:3000/profile?") {
		t.Errorf("expected absolute frontend redirect, got %q", abs)
	}
	// Trailing slash and sub-path bases join cleanly.
	sub := frontendRedirect("https://mathua.com/app/", "profile", "tok", "Ada", "s1")
	if !strings.HasPrefix(sub, "https://mathua.com/app/profile?") {
		t.Errorf("expected sub-path join, got %q", sub)
	}
	// Invalid base → safe relative fallback, never an open redirect.
	for _, bad := range []string{"javascript:alert(1)", "notaurl://", "://missing-scheme"} {
		got := frontendRedirect(bad, "/profile", "tok", "Ada", "s1")
		if !strings.HasPrefix(got, "/profile?") {
			t.Errorf("base %q: expected relative fallback, got %q", bad, got)
		}
	}
}

func TestProfileUpdate(t *testing.T) {
	d, _ := concepts.Build([]concepts.Concept{
		{ID: "a", Label: "A", Domain: "d", GradingType: "numeric", Prerequisites: []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60}},
	})
	store, _ := storage.NewSQLiteStore(":memory:")
	t.Cleanup(func() { store.Close() })
	reg := generator.NewRegistry()
	reg.Register("a", &testGen{})
	authSvc := auth.New(store)
	token, st, err := authSvc.Signup("Ada", "ada", "Engine!n1")
	if err != nil || token == "" || st == nil {
		t.Fatalf("signup: %v", err)
	}
	s := New(engine.New(store, d, reg, nil, nil), store, authSvc)
	mux := http.NewServeMux()
	s.Register(mux)
	put := func(token, body string) *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("PUT", "/api/profile", bytes.NewReader([]byte(body)))
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		mux.ServeHTTP(rec, req)
		return rec
	}

	if rec := put(token, `{"name":"  Ada Lovelace  ", "email":"Ada@Example.COM"}`); rec.Code != 200 {
		t.Fatalf("email set: expected 200, got %d: %s", rec.Code, rec.Body.String())
	} else if got, _ := store.GetStudent(st.ID); got.Email != "ada@example.com" {
		t.Errorf("expected lowercased email stored, got %q", got.Email)
	}
	token2, _, err := authSvc.Signup("Bob", "bob", "Engine!n1")
	if err != nil || token2 == "" {
		t.Fatalf("second signup: %v", err)
	}
	if rec := put(token2, `{"name":"Bob", "email":"ada@example.com"}`); rec.Code != 409 {
		t.Errorf("taken email: expected 409, got %d: %s", rec.Code, rec.Body.String())
	} else if !strings.Contains(rec.Body.String(), "email already in use") {
		t.Errorf("expected email-taken message, got %s", rec.Body.String())
	}
	// Re-saving your own address stays a no-op success.
	rec := put(token, `{"name":"Ada Lovelace", "email":"ada@example.com"}`)
	if rec.Code != 200 {
		t.Errorf("own email: expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var updated map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &updated)
	if updated["name"] != "Ada Lovelace" {
		t.Errorf("expected trimmed name, got %v", updated)
	}
	got, _ := store.GetStudent(st.ID)
	if got.Name != "Ada Lovelace" {
		t.Errorf("expected persisted name, got %q", got.Name)
	}

	for _, tc := range []struct{ name, body string }{
		{"empty", `{"name":"   "}`},
		{"missing", `{}`},
		{"too long", `{"name":"` + strings.Repeat("x", 51) + `"}`},
		{"bad json", `{"name":`},
	} {
		if rec := put(token, tc.body); rec.Code != 400 {
			t.Errorf("%s: expected 400, got %d", tc.name, rec.Code)
		}
	}
	if rec := put("", `{"name":"Bob"}`); rec.Code != 401 {
		t.Errorf("anonymous: expected 401, got %d", rec.Code)
	}
}

func avatarTestServer(t *testing.T) (*Server, *http.ServeMux, string) {
	t.Helper()
	d, _ := concepts.Build([]concepts.Concept{
		{ID: "a", Label: "A", Domain: "d", GradingType: "numeric", Prerequisites: []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60}},
	})
	store, _ := storage.NewSQLiteStore(":memory:")
	t.Cleanup(func() { store.Close() })
	reg := generator.NewRegistry()
	reg.Register("a", &testGen{})
	authSvc := auth.New(store)
	token, _, err := authSvc.Signup("Ada", "ada", "Engine!n1")
	if err != nil || token == "" {
		t.Fatalf("signup: %v", err)
	}
	s := New(engine.New(store, d, reg, nil, nil), store, authSvc)
	mux := http.NewServeMux()
	s.Register(mux)
	return s, mux, token
}

func multipartAvatar(t *testing.T, field, filename string, data []byte) (*bytes.Buffer, string) {
	t.Helper()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, err := w.CreateFormFile(field, filename)
	if err != nil {
		t.Fatal(err)
	}
	fw.Write(data)
	w.Close()
	return &buf, w.FormDataContentType()
}

// tiny valid 1x1 PNG
var tinyPNG = []byte{
	0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a,
	0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52,
	0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
	0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53,
	0xde, 0x00, 0x00, 0x00, 0x0c, 0x49, 0x44, 0x41,
	0x54, 0x08, 0xd7, 0x63, 0xf8, 0xff, 0xff, 0x3f,
	0x00, 0x05, 0xfe, 0x02, 0xfe, 0xdc, 0xcc, 0x59,
	0xe7, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4e,
	0x44, 0xae, 0x42, 0x60, 0x82,
}

func TestAvatarUploadRoundTrip(t *testing.T) {
	_, mux, token := avatarTestServer(t)
	post := func(token string, body *bytes.Buffer, ctype string) *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/avatar", body)
		req.Header.Set("Content-Type", ctype)
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		mux.ServeHTTP(rec, req)
		return rec
	}

	// Anonymous → 401.
	body, ctype := multipartAvatar(t, "avatar", "a.png", tinyPNG)
	if rec := post("", body, ctype); rec.Code != 401 {
		t.Errorf("anonymous: expected 401, got %d", rec.Code)
	}
	// Wrong field/missing file → 400.
	body, ctype = multipartAvatar(t, "notavatar", "a.png", tinyPNG)
	if rec := post(token, body, ctype); rec.Code != 400 {
		t.Errorf("missing file: expected 400, got %d", rec.Code)
	}
	// Wrong content type → 400.
	body, ctype = multipartAvatar(t, "avatar", "a.txt", []byte("hello world, this is text"))
	if rec := post(token, body, ctype); rec.Code != 400 {
		t.Errorf("text file: expected 400, got %d", rec.Code)
	}
	// Happy path.
	body, ctype = multipartAvatar(t, "avatar", "a.png", tinyPNG)
	rec := post(token, body, ctype)
	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	// Serve back with content type.
	getReq := httptest.NewRequest("GET", "/api/avatar/me", nil)
	getReq.Header.Set("Authorization", "Bearer "+token)
	getRec := httptest.NewRecorder()
	mux.ServeHTTP(getRec, getReq)
	if getRec.Code != 200 {
		t.Fatalf("expected 200 serving avatar, got %d", getRec.Code)
	}
	if getRec.Header().Get("Content-Type") != "image/png" {
		t.Errorf("expected image/png, got %q", getRec.Header().Get("Content-Type"))
	}
	if !bytes.Equal(getRec.Body.Bytes(), tinyPNG) {
		t.Error("expected identical bytes back")
	}
	// Delete → 404 on serve.
	delReq := httptest.NewRequest("DELETE", "/api/avatar", nil)
	delReq.Header.Set("Authorization", "Bearer "+token)
	delRec := httptest.NewRecorder()
	mux.ServeHTTP(delRec, delReq)
	if delRec.Code != 200 {
		t.Fatalf("expected 200 on delete, got %d", delRec.Code)
	}
	getRec2 := httptest.NewRecorder()
	mux.ServeHTTP(getRec2, getReq)
	if getRec2.Code != 404 {
		t.Errorf("expected 404 after delete, got %d", getRec2.Code)
	}
}

func TestAvatarPublicAndLeaguesSnakeCase(t *testing.T) {
	d, _ := concepts.Build([]concepts.Concept{
		{ID: "a", Label: "A", Domain: "d", GradingType: "numeric", Prerequisites: []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60}},
	})
	store, _ := storage.NewSQLiteStore(":memory:")
	t.Cleanup(func() { store.Close() })
	reg := generator.NewRegistry()
	reg.Register("a", &testGen{})
	authSvc := auth.New(store)
	token, st, err := authSvc.Signup("Ada", "ada", "Engine!n1")
	if err != nil || token == "" {
		t.Fatalf("signup: %v", err)
	}
	s := New(engine.New(store, d, reg, nil, nil), store, authSvc)
	mux := http.NewServeMux()
	s.Register(mux)
	get := func(path, token string) *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", path, nil)
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		mux.ServeHTTP(rec, req)
		return rec
	}

	// No photo yet → public 404 (no auth required).
	if rec := get("/api/avatar/"+st.ID, ""); rec.Code != 404 {
		t.Errorf("no photo: expected public 404, got %d", rec.Code)
	}
	// Upload, then the public endpoint serves identical bytes.
	body, ctype := multipartAvatar(t, "avatar", "a.png", tinyPNG)
	upRec := httptest.NewRecorder()
	upReq := httptest.NewRequest("POST", "/api/avatar", body)
	upReq.Header.Set("Content-Type", ctype)
	upReq.Header.Set("Authorization", "Bearer "+token)
	mux.ServeHTTP(upRec, upReq)
	if upRec.Code != 200 {
		t.Fatalf("upload: expected 200, got %d: %s", upRec.Code, upRec.Body.String())
	}
	pub := get("/api/avatar/"+st.ID, "")
	if pub.Code != 200 {
		t.Fatalf("public photo: expected 200, got %d", pub.Code)
	}
	if pub.Header().Get("Content-Type") != "image/png" {
		t.Errorf("expected image/png, got %q", pub.Header().Get("Content-Type"))
	}
	if !bytes.Equal(pub.Body.Bytes(), tinyPNG) {
		t.Error("expected identical bytes back")
	}
	if cc := pub.Header().Get("Cache-Control"); !strings.Contains(cc, "public") {
		t.Errorf("expected public cache header, got %q", cc)
	}
	// Exact /api/avatar/me route still requires auth (subtree must not leak).
	if rec := get("/api/avatar/me", ""); rec.Code != 401 {
		t.Errorf("own photo anonymous: expected 401, got %d", rec.Code)
	}
	// Leagues serialize snake_case (regression: untagged struct emitted
	// PascalCase keys the frontend could not read).
	rec := get("/api/leagues", token)
	if rec.Code != 200 {
		t.Fatalf("leagues: expected 200, got %d", rec.Code)
	}
	var board map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &board); err != nil {
		t.Fatalf("decode leagues: %v", err)
	}
	leagues, _ := board["leagues"].([]interface{})
	if len(leagues) == 0 {
		t.Fatal("expected at least one league tier")
	}
	members, _ := leagues[0].(map[string]interface{})["members"].([]interface{})
	if len(members) == 0 {
		t.Fatal("expected at least one league member")
	}
	m := members[0].(map[string]interface{})
	for _, k := range []string{"student_id", "name", "username", "tier", "total_mastered", "weekly_mastered", "moved"} {
		if _, ok := m[k]; !ok {
			t.Errorf("expected league member key %q, got keys %v", k, keysOf(m))
		}
	}
	if m["name"] != "Ada" || m["username"] != "ada" {
		t.Errorf("expected Ada/ada on board, got %v", m)
	}
	// Weekly board carries the username for the display fallback chain.
	wrec := get("/api/leaderboard", "")
	if wrec.Code != 200 {
		t.Fatalf("leaderboard: expected 200, got %d", wrec.Code)
	}
	var entries []map[string]interface{}
	if err := json.Unmarshal(wrec.Body.Bytes(), &entries); err != nil {
		t.Fatalf("decode leaderboard: %v", err)
	}
	if len(entries) == 0 || entries[0]["username"] != "ada" {
		t.Errorf("expected username on weekly entry, got %v", entries)
	}
}

func keysOf(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func TestAvatarOversize(t *testing.T) {
	_, mux, token := avatarTestServer(t)
	big := bytes.Repeat([]byte{0x89, 0x50}, 300*1024) // 600KB, PNG magic
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, _ := w.CreateFormFile("avatar", "big.png")
	fw.Write(big)
	w.Close()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/avatar", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	mux.ServeHTTP(rec, req)
	if rec.Code != 400 {
		t.Errorf("expected 400 for oversize, got %d", rec.Code)
	}
}

func TestPasswordResetHandlers(t *testing.T) {
	t.Setenv("SMTP_HOST", "")
	t.Setenv("FRONTEND_URL", "")
	d, _ := concepts.Build([]concepts.Concept{
		{ID: "a", Label: "A", Domain: "d", GradingType: "numeric", Prerequisites: []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60}},
	})
	store, _ := storage.NewSQLiteStore(":memory:")
	t.Cleanup(func() { store.Close() })
	reg := generator.NewRegistry()
	reg.Register("a", &testGen{})
	authSvc := auth.New(store)
	s := New(engine.New(store, d, reg, nil, nil), store, authSvc)
	mux := http.NewServeMux()
	s.Register(mux)

	// No SMTP → honest 503, not silent confusion.
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("POST", "/api/auth/reset/request", bytes.NewReader([]byte(`{"identifier":"ada"}`))))
	if rec.Code != 503 {
		t.Errorf("expected 503 without SMTP, got %d", rec.Code)
	}

	// Unknown identifier with mail configured-thin path: use SMTP set to
	// unreachable would try to send; instead assert silent-200 shape via
	// empty identifier validation.
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("POST", "/api/auth/reset/request", bytes.NewReader([]byte(`{}`))))
	if rec.Code != 400 {
		t.Errorf("expected 400 for missing identifier, got %d", rec.Code)
	}

	// Complete with garbage token → 400 invalid/expired.
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("POST", "/api/auth/reset/complete", bytes.NewReader([]byte(`{"token":"nope","password":"N3w!passw"}`))))
	if rec.Code != 400 {
		t.Errorf("expected 400 for bad token, got %d", rec.Code)
	}

	// Change password requires auth.
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("PUT", "/api/auth/password", bytes.NewReader([]byte(`{"current_password":"x","new_password":"N3w!passw"}`))))
	if rec.Code != 401 {
		t.Errorf("expected 401 anonymous, got %d", rec.Code)
	}

	// Authed change-password round trip.
	token, _, err := authSvc.Signup("Pam", "pam", "Engine!n1")
	if err != nil {
		t.Fatalf("signup: %v", err)
	}
	put := func(body, tok string) *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("PUT", "/api/auth/password", bytes.NewReader([]byte(body)))
		if tok != "" {
			req.Header.Set("Authorization", "Bearer "+tok)
		}
		mux.ServeHTTP(rec, req)
		return rec
	}
	if rec := put(`{"current_password":"Wrong!n9","new_password":"N3w!passw"}`, token); rec.Code != 400 {
		t.Errorf("wrong current: expected 400, got %d", rec.Code)
	}
	if rec := put(`{"current_password":"Engine!n1","new_password":"N3w!passw"}`, token); rec.Code != 200 {
		t.Fatalf("change: expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	loginRec := httptest.NewRecorder()
	mux.ServeHTTP(loginRec, httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader([]byte(`{"username":"pam","password":"N3w!passw"}`))))
	if loginRec.Code != 200 {
		t.Errorf("login with new password: expected 200, got %d", loginRec.Code)
	}
}

func TestIdentitiesEndpoints(t *testing.T) {
	d, _ := concepts.Build([]concepts.Concept{
		{ID: "a", Label: "A", Domain: "d", GradingType: "numeric", Prerequisites: []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60}},
	})
	store, _ := storage.NewSQLiteStore(":memory:")
	t.Cleanup(func() { store.Close() })
	reg := generator.NewRegistry()
	reg.Register("a", &testGen{})
	authSvc := auth.New(store)
	s := New(engine.New(store, d, reg, nil, nil), store, authSvc)
	mux := http.NewServeMux()
	s.Register(mux)
	token, st, err := authSvc.Signup("Ida", "ida", "Engine!n1")
	if err != nil {
		t.Fatalf("signup: %v", err)
	}
	authed := func(method, path, body string) *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()
		var rdr *bytes.Reader
		if body != "" {
			rdr = bytes.NewReader([]byte(body))
		} else {
			rdr = bytes.NewReader(nil)
		}
		req := httptest.NewRequest(method, path, rdr)
		req.Header.Set("Authorization", "Bearer "+token)
		mux.ServeHTTP(rec, req)
		return rec
	}

	// Link token mints for authed users, rejects anonymous.
	rec := authed("POST", "/api/auth/link-token", "")
	if rec.Code != 200 {
		t.Fatalf("link-token: expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var lt map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &lt)
	if lt["link_token"] == "" {
		t.Fatal("expected link_token")
	}
	anon := httptest.NewRecorder()
	mux.ServeHTTP(anon, httptest.NewRequest("POST", "/api/auth/link-token", nil))
	if anon.Code != 401 {
		t.Errorf("anonymous link-token: expected 401, got %d", anon.Code)
	}

	// Empty identities list.
	rec = authed("GET", "/api/auth/identities", "")
	if rec.Code != 200 {
		t.Fatalf("identities: expected 200, got %d", rec.Code)
	}
	// Connect a fake provider directly, then list shows it.
	if err := store.CreateIdentity("github", "42", st.ID, "ida@example.com", true); err != nil {
		t.Fatal(err)
	}
	rec = authed("GET", "/api/auth/identities", "")
	var ids []map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &ids)
	if len(ids) != 1 || ids[0]["provider"] != "github" {
		t.Errorf("expected one github identity, got %v", ids)
	}
	// Unknown provider delete → 404.
	rec = authed("DELETE", "/api/auth/identities/myspace", "")
	if rec.Code != 404 {
		t.Errorf("expected 404, got %d", rec.Code)
	}
	// Disconnect allowed (password remains).
	rec = authed("DELETE", "/api/auth/identities/github", "")
	if rec.Code != 200 {
		t.Errorf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	// Unknown OAuth provider paths → 404, unconfigured → 500.
	grec := httptest.NewRecorder()
	mux.ServeHTTP(grec, httptest.NewRequest("GET", "/api/auth/myspace/login", nil))
	if grec.Code != 404 {
		t.Errorf("expected 404 unknown provider, got %d", grec.Code)
	}
	grec = httptest.NewRecorder()
	mux.ServeHTTP(grec, httptest.NewRequest("GET", "/api/auth/github/login", nil))
	if grec.Code != 500 {
		t.Errorf("expected 500 unconfigured github, got %d", grec.Code)
	}
	// Email verify round trip via seeded token.
	if err := store.SetEmail(st.ID, "ida@example.com"); err != nil {
		t.Fatal(err)
	}
	raw := "verify-me-123"
	h := sha256.Sum256([]byte(raw))
	if err := store.CreateEmailVerification(hex.EncodeToString(h[:]), st.ID, time.Now().UTC().Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	rec = authed("POST", "/api/auth/email/verify", `{"token":"`+raw+`"}`)
	if rec.Code != 200 {
		t.Fatalf("verify: expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	got, _ := store.GetStudent(st.ID)
	if !got.EmailVerified {
		t.Error("expected email_verified set")
	}
}

// The practice endpoint is what writes the server-side answer anchor, and it
// keys that anchor by student. It could not: /api/lessons/ was registered
// without auth middleware, so the bearer token was never resolved into the
// request context, and the only other route to an anchor was a student_id the
// client does not send on that call. The anchor was therefore never written
// for anyone, and every study answer was graded against the client's own
// expected answer — the exact cheat the anchor was added to prevent.
func TestPractice_WritesTheStudyAnchor(t *testing.T) {
	_, mux, _ := guestServer(t)
	token, studentID := quizGuest(t, mux)

	rec := httptest.NewRecorder()
	get := httptest.NewRequest("GET", "/api/lessons/a/practice?count=1", nil)
	get.Header.Set("Authorization", "Bearer "+token)
	mux.ServeHTTP(rec, get)
	if rec.Code != 200 {
		t.Fatalf("practice: %d %s", rec.Code, rec.Body.String())
	}
	var set struct {
		Questions []struct{ Question, Answer, Explanation string } `json:"questions"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &set); err != nil {
		t.Fatalf("decode practice: %v", err)
	}
	if len(set.Questions) == 0 {
		t.Fatal("expected a served question")
	}
	question := set.Questions[0].Question
	answer := set.Questions[0].Answer

	// A deliberately wrong answer for that exact question. Before the
	// middleware it came back 409 (no anchor); now it grades, and the answer it
	// graded against is the served one rather than whatever the client claimed.
	body, _ := json.Marshal(map[string]interface{}{
		"concept_id": "a",
		"answer":     "definitely wrong",
		"expected":   "also wrong",
		"elapsed":    5.0,
		"question":   question,
		"student_id": studentID,
	})
	rec = httptest.NewRecorder()
	post := httptest.NewRequest("POST", "/api/study/answer", bytes.NewReader(body))
	post.Header.Set("Authorization", "Bearer "+token)
	mux.ServeHTTP(rec, post)
	if rec.Code == 409 {
		t.Fatalf("practice did not write an anchor: %s", rec.Body.String())
	}
	if rec.Code != 200 {
		t.Fatalf("study answer: %d %s", rec.Code, rec.Body.String())
	}
	var res struct {
		Correct        bool   `json:"correct"`
		ExpectedAnswer string `json:"expected_answer"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("decode answer: %v", err)
	}
	if res.Correct {
		t.Error("a wrong answer graded correct")
	}
	if res.ExpectedAnswer != answer {
		t.Errorf("expected_answer = %q, want the served %q", res.ExpectedAnswer, answer)
	}
}

// The study request carries no expected answer and needs none: the answer and
// the explanation come from the anchor written when the question was served.
// These two tests pin the guarantees that follow from that.
func TestStudyAnswer_AnchoredNotClientSupplied(t *testing.T) {
	_, mux, _ := guestServer(t)
	token, studentID := quizGuest(t, mux)

	// Serve a set, which is what writes the anchor.
	rec := httptest.NewRecorder()
	get := httptest.NewRequest("GET", "/api/lessons/a/practice?student_id="+studentID+"&count=1", nil)
	get.Header.Set("Authorization", "Bearer "+token)
	mux.ServeHTTP(rec, get)
	if rec.Code != 200 {
		t.Fatalf("practice: %d %s", rec.Code, rec.Body.String())
	}
	var set struct {
		Questions []struct{ Question, Answer string } `json:"questions"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &set); err != nil {
		t.Fatalf("decode practice: %v", err)
	}
	if len(set.Questions) == 0 {
		t.Fatal("expected a served question")
	}
	question := set.Questions[0].Question
	anchored := set.Questions[0].Answer

	// The forgery: answer something wrong while claiming it is what was
	// expected. A server that honoured the client would call this correct.
	const forged = "987654321"
	body, _ := json.Marshal(map[string]interface{}{
		"concept_id": "a",
		"answer":     forged,
		"expected":   forged,
		"elapsed":    5.0,
		"question":   question,
		"student_id": studentID,
	})
	rec = httptest.NewRecorder()
	post := httptest.NewRequest("POST", "/api/study/answer", bytes.NewReader(body))
	post.Header.Set("Authorization", "Bearer "+token)
	mux.ServeHTTP(rec, post)
	if rec.Code != 200 {
		t.Fatalf("study answer: %d %s", rec.Code, rec.Body.String())
	}
	var res struct {
		Correct        bool   `json:"correct"`
		ExpectedAnswer string `json:"expected_answer"`
		Explanation    string `json:"explanation"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("decode answer: %v", err)
	}
	if res.Correct {
		t.Errorf("a forged expected answer was honoured: %q graded correct", forged)
	}
	if res.ExpectedAnswer != anchored {
		t.Errorf("expected_answer = %q, want the anchored %q", res.ExpectedAnswer, anchored)
	}
	if res.Explanation == "" {
		t.Error("a graded answer must carry the served explanation")
	}
}

// A question the server has no record of cannot be graded, and the honest
// response is to say so rather than accept the caller's word for it. 409 (not
// 400) is what tells the client to re-serve.
func TestStudyAnswer_UnknownQuestionIs409(t *testing.T) {
	_, mux, _ := guestServer(t)
	token, studentID := quizGuest(t, mux)

	body, _ := json.Marshal(map[string]interface{}{
		"concept_id": "a",
		"answer":     "4",
		"expected":   "4", // ignored; the anchor is what matters
		"elapsed":    5.0,
		"question":   "a question the server never served",
		"student_id": studentID,
	})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/study/answer", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	mux.ServeHTTP(rec, req)
	if rec.Code != 409 {
		t.Errorf("expected 409 for an unserved question, got %d: %s", rec.Code, rec.Body.String())
	}
}

// Regression for the Learn loop dropping the questions a learner had just
// answered. The client buffers served questions and serves a second batch while
// the learner is still reading the first, so the practice endpoint is called
// twice for one concept before the first question is submitted. Replacing the
// anchor blob on the second call erased the first, and the submit came back 409
// — which the client reports as "that question had expired" and answers by
// swapping in a different question.
//
// The fix is in SetStudyAnchorBatch (merge, don't replace); this pins the HTTP
// behaviour that the learner actually sees.
func TestStudyAnswer_GradesFirstBatchAfterASecondPracticeFetch(t *testing.T) {
	_, mux, _ := varyingGuestServer(t)
	token, studentID := quizGuest(t, mux)

	serve := func(exclude []string) []struct{ Question, Answer, Explanation string } {
		t.Helper()
		u := fmt.Sprintf("/api/lessons/a/practice?count=3&seed=11&difficulty=0.4&exclude=%s",
			url.QueryEscape(strings.Join(exclude, "\n")))
		req := httptest.NewRequest("GET", u, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != 200 {
			t.Fatalf("practice: %d %s", rec.Code, rec.Body.String())
		}
		var set struct {
			Questions []struct{ Question, Answer, Explanation string } `json:"questions"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &set); err != nil {
			t.Fatalf("decode practice: %v", err)
		}
		if len(set.Questions) == 0 {
			t.Fatal("expected a served question")
		}
		return set.Questions
	}

	// Batch one is what the learner answers; batch two is the prefetch the
	// client fires while they are still reading batch one, and it excludes the
	// texts already seen -- so the two sets do not overlap, which is the point.
	first := serve(nil)
	second := serve([]string{first[0].Question, first[1%len(first)].Question})
	if len(second) == 0 {
		t.Fatal("expected the second batch to serve questions")
	}
	for _, q := range second {
		if q.Question == first[0].Question {
			t.Fatal("exclude did not filter the first batch out; the test would not reproduce the bug")
		}
	}

	body, _ := json.Marshal(map[string]interface{}{
		"concept_id": "a",
		"answer":     first[0].Answer,
		"elapsed":    5.0,
		"question":   first[0].Question,
		"student_id": studentID,
	})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/study/answer", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	mux.ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatalf("a question served before the second batch must still grade; got %d: %s", rec.Code, rec.Body.String())
	}
	var res struct {
		Correct        bool   `json:"correct"`
		ExpectedAnswer string `json:"expected_answer"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("decode answer: %v", err)
	}
	if !res.Correct {
		t.Error("expected the served answer to grade correct")
	}
	if res.ExpectedAnswer != first[0].Answer {
		t.Errorf("graded against %q, want the first batch's %q", res.ExpectedAnswer, first[0].Answer)
	}

	// A genuinely unserved question is still refused: merging must not weaken
	// the refusal that 409 exists for.
	body, _ = json.Marshal(map[string]interface{}{
		"concept_id": "a",
		"answer":     "4",
		"elapsed":    5.0,
		"question":   "a question no batch ever served",
		"student_id": studentID,
	})
	rec = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/api/study/answer", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	mux.ServeHTTP(rec, req)
	if rec.Code != 409 {
		t.Errorf("expected 409 for a question no batch served, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestStudyAnswer_UnknownConcept(t *testing.T) {
	s := testServer(t)
	mux := http.NewServeMux()
	s.Register(mux)

	body, _ := json.Marshal(map[string]interface{}{
		"concept_id": "nope.not.real",
		"answer":     "1",
		"expected":   "1",
		"elapsed":    5.0,
		"student_id": "ghost",
	})
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("POST", "/api/study/answer", bytes.NewReader(body)))
	if rec.Code != 404 {
		t.Errorf("expected 404 for unknown concept, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestStudyAnswer_TooFast(t *testing.T) {
	s := testServer(t)
	mux := http.NewServeMux()
	s.Register(mux)

	body, _ := json.Marshal(map[string]interface{}{
		"concept_id": "a",
		"answer":     "4",
		"expected":   "4",
		"elapsed":    0.1,
		"student_id": "ghost",
	})
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("POST", "/api/study/answer", bytes.NewReader(body)))
	if rec.Code != 400 {
		t.Errorf("expected 400 for rushed answer, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestProgressScores_Ownership(t *testing.T) {
	d, _ := concepts.Build([]concepts.Concept{
		{ID: "a", Label: "A", Domain: "d", GradingType: "numeric", Prerequisites: []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60}},
	})
	store, _ := storage.NewSQLiteStore(":memory:")
	t.Cleanup(func() { store.Close() })
	reg := generator.NewRegistry()
	reg.Register("a", &testGen{})
	authSvc := auth.New(store)
	victimToken, victim, err := authSvc.Signup("Victim", "victim", "Engine!n1")
	if err != nil || victimToken == "" {
		t.Fatalf("victim signup: %v", err)
	}
	attackerToken, _, err := authSvc.Signup("Attacker", "attacker", "Engine!n1")
	if err != nil || attackerToken == "" {
		t.Fatalf("attacker signup: %v", err)
	}
	// Victim has real progress.
	_ = store.UpsertProgress(&storage.ConceptProgress{StudentID: victim.ID, ConceptID: "a", Status: "LEARNING"})
	s := New(engine.New(store, d, reg, nil, nil), store, authSvc)
	mux := http.NewServeMux()
	s.Register(mux)
	get := func(path, token string) *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", path, nil)
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		mux.ServeHTTP(rec, req)
		return rec
	}

	// No token → 401 even with a valid victim ID (no enumeration).
	if rec := get("/api/progress/"+victim.ID, ""); rec.Code != 401 {
		t.Errorf("progress: expected 401 without token, got %d", rec.Code)
	}
	if rec := get("/api/scores/"+victim.ID, ""); rec.Code != 401 {
		t.Errorf("scores: expected 401 without token, got %d", rec.Code)
	}
	// Attacker asking for the victim's ID gets their OWN (empty) progress.
	rec := get("/api/progress/"+victim.ID, attackerToken)
	if rec.Code != 200 {
		t.Fatalf("progress: expected 200 for authed caller, got %d", rec.Code)
	}
	var progress map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &progress); err != nil {
		t.Fatalf("decode progress: %v", err)
	}
	if _, ok := progress["a"]; ok {
		t.Error("progress: victim concept leaked to attacker")
	}
	// Guest exception: guest_ IDs remain readable without a token.
	if rec := get("/api/progress/guest_local123", ""); rec.Code == 401 {
		t.Errorf("progress: expected guest_ IDs readable without token, got 401")
	}
}

func guestServer(t *testing.T) (*Server, *http.ServeMux, storage.Repository) {
	t.Helper()
	d, _ := concepts.Build([]concepts.Concept{
		{ID: "a", Label: "A", Domain: "d", GradingType: "numeric", Prerequisites: []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60}},
	})
	store, _ := storage.NewSQLiteStore(":memory:")
	t.Cleanup(func() { store.Close() })
	reg := generator.NewRegistry()
	reg.Register("a", &testGen{})
	s := New(engine.New(store, d, reg, nil, nil), store, auth.New(store))
	mux := http.NewServeMux()
	s.Register(mux)
	return s, mux, store
}

func postGuest(t *testing.T, mux *http.ServeMux, body, token string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/auth/guest", bytes.NewReader([]byte(body)))
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	mux.ServeHTTP(rec, req)
	return rec
}

func TestGuestToken_CreateAndReadOwn(t *testing.T) {
	_, mux, _ := guestServer(t)

	rec := postGuest(t, mux, `{}`, "")
	if rec.Code != 200 {
		t.Fatalf("create guest: expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var res map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("decode: %v", err)
	}
	token, _ := res["token"].(string)
	sid, _ := res["student_id"].(string)
	if token == "" || sid == "" {
		t.Fatalf("expected token+student_id, got %v", res)
	}
	// Guest reads own progress/scores with the token.
	for _, path := range []string{"/api/progress/" + sid, "/api/scores/" + sid} {
		r := httptest.NewRecorder()
		req := httptest.NewRequest("GET", path, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		mux.ServeHTTP(r, req)
		if r.Code != 200 {
			t.Errorf("%s: expected 200 with guest token, got %d", path, r.Code)
		}
	}
}

func TestGuestToken_ClaimAdoptsProgress(t *testing.T) {
	_, mux, store := guestServer(t)

	// Existing guest row with practice history (FKs enforced: progress rows
	// always have a students row; claim re-mints a token for it).
	ghost := "guest_claim_abc123"
	if _, err := store.ClaimGuestStudent(ghost, "Guest"); err != nil {
		t.Fatalf("seed guest row: %v", err)
	}
	_ = store.UpsertProgress(&storage.ConceptProgress{StudentID: ghost, ConceptID: "a", Status: "LEARNING"})
	rec := postGuest(t, mux, `{"student_id":"`+ghost+`"}`, "")
	if rec.Code != 200 {
		t.Fatalf("claim: expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var res map[string]interface{}
	_ = json.Unmarshal(rec.Body.Bytes(), &res)
	token, _ := res["token"].(string)
	if res["student_id"] != ghost {
		t.Fatalf("expected claimed id %q, got %v", ghost, res)
	}
	// Adopted progress is visible to the token holder.
	r := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/progress/"+ghost, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	mux.ServeHTTP(r, req)
	if r.Code != 200 {
		t.Fatalf("progress: expected 200, got %d", r.Code)
	}
	var progress map[string]interface{}
	_ = json.Unmarshal(r.Body.Bytes(), &progress)
	if _, ok := progress["a"]; !ok {
		t.Error("expected adopted concept progress after claim")
	}
}

func TestGuestToken_RefusesRegistered(t *testing.T) {
	_, mux, _ := guestServer(t)

	// Build a registered user via signup, then try to claim their ID as guest.
	srec := httptest.NewRecorder()
	sreq := httptest.NewRequest("POST", "/api/auth/signup", bytes.NewReader([]byte(`{"name":"Reg","username":"reguser","password":"Engine!n1","email":"reg@example.com"}`)))
	mux.ServeHTTP(srec, sreq)
	if srec.Code != 200 {
		t.Fatalf("signup: expected 200, got %d: %s", srec.Code, srec.Body.String())
	}
	var sres map[string]interface{}
	_ = json.Unmarshal(srec.Body.Bytes(), &sres)
	sid, _ := sres["student_id"].(string)
	if rec := postGuest(t, mux, `{"student_id":"`+sid+`"}`, ""); rec.Code != 403 {
		t.Errorf("claim registered: expected 403, got %d: %s", rec.Code, rec.Body.String())
	}
	// Unknown non-guest ID → 404, not a fresh account.
	if rec := postGuest(t, mux, `{"student_id":"01234567-89ab-cdef-0123-456789abcdef"}`, ""); rec.Code != 404 {
		t.Errorf("claim unknown uuid: expected 404, got %d: %s", rec.Code, rec.Body.String())
	}
}

// Batch 1 helpers: guest-authed JSON POST.

func quizPost(t *testing.T, mux *http.ServeMux, path, body, token string) map[string]interface{} {
	t.Helper()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", path, bytes.NewReader([]byte(body)))
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	mux.ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatalf("%s: expected 200, got %d: %s", path, rec.Code, rec.Body.String())
	}
	var res map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("%s: decode: %v", path, err)
	}
	return res
}

func quizGuest(t *testing.T, mux *http.ServeMux) (token, studentID string) {
	t.Helper()
	rec := postGuest(t, mux, `{}`, "")
	if rec.Code != 200 {
		t.Fatalf("create guest: %d %s", rec.Code, rec.Body.String())
	}
	var res map[string]interface{}
	_ = json.Unmarshal(rec.Body.Bytes(), &res)
	return res["token"].(string), res["student_id"].(string)
}

// Batch 1: quiz session carries timed closed-book contract; completing it
// records the gate baseline and offers retake.

// The quiz response never carried an explanation, so a wrong answer showed the
// grader's "Incorrect" and a right one showed nothing at all. The session
// already holds the served instance's solution; it must reach the client on
// both verdicts, exactly as handleGoalDiagnosticAnswer already did.
func TestQuizAnswer_ExplanationOnBothVerdicts(t *testing.T) {
	_, mux, _ := guestServer(t)
	token, _ := quizGuest(t, mux)

	// Wrong answer. The session holds one question, so this also ends the quiz
	// — which is fine, it is the done response that also needs the solution.
	sess := quizPost(t, mux, "/api/quiz/session", `{}`, token)
	miss := quizPost(t, mux, "/api/quiz/answer", `{"session_id":"`+sess["session_id"].(string)+`","concept_id":"a","answer":"999","elapsed":5}`, token)
	if miss["correct"] != false {
		t.Fatalf("expected the wrong answer to grade incorrect, got %v", miss)
	}
	missExpl, _ := miss["explanation"].(string)
	if missExpl == "" {
		t.Error("expected an explanation on a wrong quiz answer")
	}
	if missExpl == "Incorrect" {
		t.Errorf("explanation must be the solution, not the grader token: %q", missExpl)
	}

	// Correct answer, fresh session.
	sess2 := quizPost(t, mux, "/api/quiz/session", `{}`, token)
	hit := quizPost(t, mux, "/api/quiz/answer", `{"session_id":"`+sess2["session_id"].(string)+`","concept_id":"a","answer":"4","elapsed":5}`, token)
	if hit["correct"] != true {
		t.Fatalf("expected the exact answer to grade correct, got %v", hit)
	}
	if expl, _ := hit["explanation"].(string); expl == "" {
		t.Error("a correct quiz answer must still carry its explanation")
	}
}

func TestQuizSession_ClosedBookContractAndCompletion(t *testing.T) {
	_, mux, _ := guestServer(t)
	token, _ := quizGuest(t, mux)

	sess := quizPost(t, mux, "/api/quiz/session", `{}`, token)
	if sess["done"] == true {
		t.Fatal("expected quiz question, got done")
	}
	if sess["closed_book"] != true {
		t.Errorf("expected closed_book true, got %v", sess)
	}
	if sess["time_limit_seconds"] != float64(60) {
		t.Errorf("expected 60s limit, got %v", sess["time_limit_seconds"])
	}
	if sess["questions_total"] != float64(1) {
		t.Errorf("expected 1 question total, got %v", sess["questions_total"])
	}

	ans := quizPost(t, mux, "/api/quiz/answer", `{"session_id":"`+sess["session_id"].(string)+`","concept_id":"a","answer":"4","elapsed":5}`, token)
	if ans["done"] != true || ans["correct"] != true {
		t.Fatalf("expected done+correct, got %v", ans)
	}
	if ans["retake_available"] != true {
		t.Errorf("expected retake_available, got %v", ans)
	}
	if ans["xp"] == float64(0) {
		t.Errorf("expected TaskQuiz XP, got %v", ans["xp"])
	}

	// Gate baseline reset: scores show 0 since quiz, not due.
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/scores/x", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	mux.ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatalf("scores: %d %s", rec.Code, rec.Body.String())
	}
	var scores map[string]interface{}
	_ = json.Unmarshal(rec.Body.Bytes(), &scores)
	if scores["xp_since_quiz"] != float64(0) {
		t.Errorf("expected xp_since_quiz 0 after completion, got %v", scores["xp_since_quiz"])
	}
	if scores["quiz_due"] == true {
		t.Errorf("expected quiz_due false after completion, got %v", scores)
	}

	// Retake: a fresh session starts immediately after completion.
	sess2 := quizPost(t, mux, "/api/quiz/session", `{}`, token)
	if sess2["done"] == true || sess2["session_id"] == sess["session_id"] {
		t.Errorf("expected fresh retake session, got %v", sess2)
	}
}

// Batch 1: quiz miss returns immediate remedial; 150 XP of study earns due.

func TestQuizMiss_RemedialAndGateDue(t *testing.T) {
	s, mux, store := guestServer(t)
	token, sid := quizGuest(t, mux)

	sess := quizPost(t, mux, "/api/quiz/session", `{}`, token)
	ans := quizPost(t, mux, "/api/quiz/answer", `{"session_id":"`+sess["session_id"].(string)+`","concept_id":"a","answer":"999","elapsed":5}`, token)
	if ans["correct"] == true {
		t.Fatalf("expected incorrect, got %v", ans)
	}
	rem, _ := ans["remedial"].([]interface{})
	if len(rem) != 1 || rem[0] != "a" {
		t.Errorf("expected remedial [a], got %v", ans["remedial"])
	}
	if got := s.eng.QuizRemedial(sid); len(got) != 1 || got[0] != "a" {
		t.Errorf("expected engine queue [a], got %v", got)
	}

	// Earn 150 XP of study: gate flips due (completion earlier baselined 0).
	if err := store.AddXP(sid, 150); err != nil {
		t.Fatalf("add xp: %v", err)
	}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/scores/x", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	mux.ServeHTTP(rec, req)
	var scores map[string]interface{}
	_ = json.Unmarshal(rec.Body.Bytes(), &scores)
	if scores["quiz_due"] != true {
		t.Errorf("expected quiz_due true at 150 since, got %v", scores)
	}
	if scores["xp_since_quiz"] != float64(150) {
		t.Errorf("expected xp_since_quiz 150, got %v", scores["xp_since_quiz"])
	}
}

// Batch 2: goal diagnostic end-to-end — start, answer to completion, report
// carries placement, frontier, gaps, and completion estimates. (The legacy
// client-graded /api/diagnostic* endpoint was removed; it trusted a client
// `correct` flag and was unused by the frontend.)

func TestGoalDiagnostic_FullFlowReport(t *testing.T) {
	_, mux := twoConceptServer(t)

	code, start := postJSON(t, mux, "/api/goal/diagnostic", `{"name":"tester","concept_ids":["a","b"]}`)
	if code != 200 {
		t.Fatalf("start: %d %v", code, start)
	}
	sid, _ := start["session_id"].(string)
	cid, _ := start["concept_id"].(string)
	if sid == "" || cid == "" {
		t.Fatalf("expected session+concept, got %v", start)
	}

	var done map[string]interface{}
	for i := 0; i < 60; i++ {
		body := `{"session_id":"` + sid + `","concept_id":"` + cid + `","answer":"0","elapsed":5.0}`
		code, done = postJSON(t, mux, "/api/goal/diagnostic/answer", body)
		if code != 200 {
			t.Fatalf("answer %d: %d %v", i, code, done)
		}
		if d, _ := done["done"].(bool); d {
			break
		}
		cid, _ = done["concept_id"].(string)
		if cid == "" {
			break
		}
	}
	if d, _ := done["done"].(bool); !d {
		t.Fatalf("expected diagnostic completion, last=%v", done)
	}
	rep, _ := done["report"].(map[string]interface{})
	if rep == nil {
		t.Fatalf("expected report, got %v", done)
	}
	for _, key := range []string{"placement_course_id", "frontier_idx", "frontier_label", "gaps_by_domain", "mastery_levels", "confidence", "completion_estimates", "total_questions"} {
		if _, ok := rep[key]; !ok {
			t.Errorf("expected report key %q, got %v", key, rep)
		}
	}
}

// D: longitudinal efficacy trend endpoint (weekly buckets + retention).
func TestEfficacyTrend_Endpoint(t *testing.T) {
	_, mux, _ := guestServer(t)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("GET", "/api/efficacy/trend", nil))
	if rec.Code != 200 {
		t.Fatalf("efficacy trend: expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var res map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("decode: %v", err)
	}
	for _, key := range []string{"weeks", "total_students", "returning_students", "retention_rate", "first_pass_trend"} {
		if _, ok := res[key]; !ok {
			t.Errorf("expected key %q, got %v", key, res)
		}
	}
}

func TestAccountDelete(t *testing.T) {
	d, _ := concepts.Build([]concepts.Concept{
		{ID: "a", Label: "A", Domain: "d", GradingType: "numeric", Prerequisites: []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60}},
	})
	store, _ := storage.NewSQLiteStore(":memory:")
	t.Cleanup(func() { store.Close() })
	reg := generator.NewRegistry()
	reg.Register("a", &testGen{})
	authSvc := auth.New(store)
	s := New(engine.New(store, d, reg, nil, nil), store, authSvc)
	mux := http.NewServeMux()
	s.Register(mux)
	del := func(token, body string) *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("DELETE", "/api/account", bytes.NewReader([]byte(body)))
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		mux.ServeHTTP(rec, req)
		return rec
	}

	pwToken, pwSt, err := authSvc.Signup("Del", "deluser", "Engine!n1")
	if err != nil || pwToken == "" {
		t.Fatalf("signup: %v", err)
	}
	if err := store.SetEmail(pwSt.ID, "del@example.com"); err != nil {
		t.Fatal(err)
	}

	// Wrong phrase → 400, account intact.
	if rec := del(pwToken, `{"phrase":"nope","password":"Engine!n1"}`); rec.Code != 400 {
		t.Fatalf("phrase: expected 400, got %d: %s", rec.Code, rec.Body.String())
	}
	// Password account, missing password → 401.
	if rec := del(pwToken, `{"phrase":"delete my account"}`); rec.Code != 401 {
		t.Fatalf("missing password: expected 401, got %d: %s", rec.Code, rec.Body.String())
	}
	// Wrong password → 401.
	if rec := del(pwToken, `{"phrase":"delete my account","password":"Wrong!n9"}`); rec.Code != 401 {
		t.Fatalf("wrong password: expected 401, got %d: %s", rec.Code, rec.Body.String())
	}
	// Correct → 200, rows gone, address reusable.
	if rec := del(pwToken, `{"phrase":"delete my account","password":"Engine!n1"}`); rec.Code != 200 {
		t.Fatalf("delete: expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	if got, _ := store.GetStudent(pwSt.ID); got != nil {
		t.Error("student row survives")
	}
	if p, _ := store.GetAllProgress(pwSt.ID); len(p) != 0 {
		t.Errorf("progress survives: %d", len(p))
	}
	// Retry → 404 (row already gone).
	if rec := del(pwToken, `{"phrase":"delete my account","password":"Engine!n1"}`); rec.Code != 404 {
		t.Fatalf("retry: expected 404, got %d: %s", rec.Code, rec.Body.String())
	}
	// Wrong method → 405 (token still verifies; the row is gone).
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/account", bytes.NewReader([]byte(`{}`)))
	req.Header.Set("Authorization", "Bearer "+pwToken)
	mux.ServeHTTP(rec, req)
	if rec.Code != 405 {
		t.Fatalf("method: expected 405, got %d", rec.Code)
	}

	// OAuth-only account: phrase alone suffices.
	oauthSt, err := store.CreateGoogleUser("G", "gdel@example.com", "gid-del-acct", "")
	if err != nil {
		t.Fatal(err)
	}
	oauthTok, _, err := authSvc.LoginOrCreateOAuth(auth.OAuthProfile{Provider: "google", ProviderID: "gid-del-acct", Email: "gdel@example.com", EmailVerified: true, Name: "G"})
	if err != nil || oauthTok == "" {
		t.Fatalf("oauth login: %v", err)
	}
	if rec := del(oauthTok, `{"phrase":"delete my account"}`); rec.Code != 200 {
		t.Fatalf("oauth delete: expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	if got, _ := store.GetStudent(oauthSt.ID); got != nil {
		t.Error("oauth student row survives")
	}
}

// GET /api/lessons resolved the caller's identity only from a student_id query
// param. authedFetch attaches the guest token, but the Study page sends no such
// param for a guest, so the one progression signal Study is still allowed to
// render — the per-concept status words and badges inside LessonDetail — was
// invisible to every signed-out learner who had in fact been practising.
func TestLessons_ResolvesIdentityFromToken(t *testing.T) {
	d, _ := concepts.Build([]concepts.Concept{
		{ID: "a", Label: "A", Domain: "d", GradingType: "numeric", Prerequisites: []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60}},
	})
	store, _ := storage.NewSQLiteStore(":memory:")
	t.Cleanup(func() { store.Close() })
	reg := generator.NewRegistry()
	reg.Register("a", &testGen{})
	ll := miniLessons(t, map[string]string{"a": "Lesson A\n\nBody."})
	s := New(engine.New(store, d, reg, ll, nil), store, auth.New(store))
	mux := http.NewServeMux()
	s.Register(mux)

	token, studentID := quizGuest(t, mux)

	// Give the guest progress on concept "a" through the live answer path.
	rec := httptest.NewRecorder()
	practice := httptest.NewRequest("GET", "/api/lessons/a/practice?count=1", nil)
	practice.Header.Set("Authorization", "Bearer "+token)
	mux.ServeHTTP(rec, practice)
	if rec.Code != 200 {
		t.Fatalf("practice: %d %s", rec.Code, rec.Body.String())
	}
	var set struct {
		Questions []struct{ Question, Answer string } `json:"questions"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &set); err != nil {
		t.Fatalf("decode practice: %v", err)
	}
	if len(set.Questions) == 0 {
		t.Fatal("expected a served question")
	}
	body, _ := json.Marshal(map[string]interface{}{
		"concept_id": "a",
		"answer":     set.Questions[0].Answer,
		"elapsed":    5.0,
		"question":   set.Questions[0].Question,
		"student_id": studentID,
	})
	rec = httptest.NewRecorder()
	answer := httptest.NewRequest("POST", "/api/study/answer", bytes.NewReader(body))
	answer.Header.Set("Authorization", "Bearer "+token)
	mux.ServeHTTP(rec, answer)
	if rec.Code != 200 {
		t.Fatalf("study answer: %d %s", rec.Code, rec.Body.String())
	}

	// Token only, no student_id param — exactly what the Study page sends.
	rec = httptest.NewRecorder()
	list := httptest.NewRequest("GET", "/api/lessons", nil)
	list.Header.Set("Authorization", "Bearer "+token)
	mux.ServeHTTP(rec, list)
	if rec.Code != 200 {
		t.Fatalf("lessons: %d %s", rec.Code, rec.Body.String())
	}
	own := catalogStatus(t, rec.Body.String())
	if own == "" {
		t.Fatalf("a token-only caller got no progress for concept a: %s", rec.Body.String())
	}

	// An anonymous caller still gets the public catalog with no progress at
	// all — the catalog is not private, only the status is.
	rec = httptest.NewRecorder()
	anon := httptest.NewRequest("GET", "/api/lessons", nil)
	mux.ServeHTTP(rec, anon)
	if got := catalogStatus(t, rec.Body.String()); got != "" {
		t.Errorf("an anonymous caller saw progress %q; the catalog must stay public and untracked", got)
	}

	// A mismatched param must not redirect whose progress is returned. The
	// validated identity wins, so this is the caller's own status again.
	rec = httptest.NewRecorder()
	other := httptest.NewRequest("GET", "/api/lessons?student_id=someone-else", nil)
	other.Header.Set("Authorization", "Bearer "+token)
	mux.ServeHTTP(rec, other)
	if rec.Code != 200 {
		t.Fatalf("lessons with foreign param: %d %s", rec.Code, rec.Body.String())
	}
	if got := catalogStatus(t, rec.Body.String()); got != own {
		t.Errorf("a foreign student_id changed whose progress was returned: got %q, want the caller's own %q", got, own)
	}
}

// catalogStatus returns the reported status for concept "a" across every
// lesson in a /api/lessons response, or "" when none carries progress.
func catalogStatus(t *testing.T, body string) string {
	t.Helper()
	var catalog struct {
		Lessons map[string][]struct {
			Concepts []string                  `json:"concepts"`
			Progress map[string]map[string]any `json:"progress"`
		} `json:"lessons"`
	}
	if err := json.Unmarshal([]byte(body), &catalog); err != nil {
		t.Fatalf("decode lessons: %v — body: %s", err, body)
	}
	for _, lessons := range catalog.Lessons {
		for _, l := range lessons {
			if cp, ok := l.Progress["a"]; ok {
				if st, ok := cp["status"].(string); ok {
					return st
				}
			}
		}
	}
	return ""
}

// miniLessons builds a lesson loader over a temp dir so a test can exercise the
// catalog path. The real corpus is 565 files; a unit test needs one.
func miniLessons(t *testing.T, byConcept map[string]string) *lessons.Loader {
	t.Helper()
	dir := t.TempDir()
	var mappings []map[string]string
	i := 0
	for cid, body := range byConcept {
		name := fmt.Sprintf("lesson-%d.md", i)
		i++
		if err := os.WriteFile(filepath.Join(dir, name), []byte("# Lesson\n\n"+body), 0o600); err != nil {
			t.Fatal(err)
		}
		mappings = append(mappings, map[string]string{"concept_id": cid, "source": name})
	}
	raw, _ := json.Marshal(mappings)
	if err := os.WriteFile(filepath.Join(dir, "lessons.json"), raw, 0o600); err != nil {
		t.Fatal(err)
	}
	ll, err := lessons.Load(dir)
	if err != nil {
		t.Fatalf("mini lessons: %v", err)
	}
	return ll
}

// fixedBankGen produces a small, fixed set of question texts — the shape that ended a Learn
// session dead.
//
// `geo.basic.points_lines` really does this: four definition-recall variants, nothing else.
// Once a session had seen all four, the practice endpoint's exclusion filter emptied the
// candidate set and it answered 200 with `questions: []`, which the Learn feed read as
// success. No card, no message, no way forward.
type fixedBankGen struct{ n int64 }

func (g *fixedBankGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	bank := []struct{ q, a string }{
		{"What is: an exact location with no size?", "point"},
		{"What is: goes forever in two directions?", "line"},
		{"What is: has one endpoint and goes forever?", "ray"},
		{"What is: has two endpoints?", "line segment"},
	}
	i := int(atomic.AddInt64(&g.n, 1)+ctx.Seed) % len(bank)
	if i < 0 {
		i += len(bank)
	}
	return generator.Problem{Question: bank[i].q, Answer: bank[i].a}
}

func fixedBankServer(t *testing.T) (*Server, *http.ServeMux) {
	t.Helper()
	d, err := concepts.Build([]concepts.Concept{
		{ID: "a", Label: "A", Domain: "d", GradingType: "numeric", Prerequisites: []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60}},
	})
	if err != nil {
		t.Fatalf("build dag: %v", err)
	}
	store, _ := storage.NewSQLiteStore(":memory:")
	t.Cleanup(func() { store.Close() })
	reg := generator.NewRegistry()
	reg.Register("a", &fixedBankGen{})
	s := New(engine.New(store, d, reg, nil, nil), store, auth.New(store))
	mux := http.NewServeMux()
	s.Register(mux)
	return s, mux
}

// TestPracticeNeverReturnsEmptyOnceEveryVariantHasBeenSeen is the server half of the
// dead-end fix.
//
// Asking for practice with every one of a concept's variants in `exclude` used to return
// 200 with an empty array. Repeats are a much better answer than a stranded learner: the
// alternative is not "no practice", it is "no next step, silently".
func TestPracticeNeverReturnsEmptyOnceEveryVariantHasBeenSeen(t *testing.T) {
	_, mux := fixedBankServer(t)

	// Collect every distinct question the bank can produce.
	all := map[string]bool{}
	for i := 0; i < 12; i++ {
		req := httptest.NewRequest("GET", fmt.Sprintf("/api/lessons/a/practice?count=3&seed=%d", i), nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != 200 {
			t.Fatalf("practice %d: %d %s", i, rec.Code, rec.Body.String())
		}
		var set struct {
			Questions []struct{ Question string } `json:"questions"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &set); err != nil {
			t.Fatalf("decode: %v", err)
		}
		for _, q := range set.Questions {
			all[q.Question] = true
		}
	}
	if len(all) < 2 {
		t.Fatalf("test setup: bank produced %d distinct questions", len(all))
	}

	// Now exclude all of them at once. This is what a session does once it has seen the
	// whole bank, and it used to come back empty.
	exclude := make([]string, 0, len(all))
	for q := range all {
		exclude = append(exclude, q)
	}
	sort.Strings(exclude)
	u := "/api/lessons/a/practice?count=3&exclude=" + url.QueryEscape(strings.Join(exclude, "\n"))
	req := httptest.NewRequest("GET", u, nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Fatalf("with every variant excluded: %d %s", rec.Code, rec.Body.String())
	}
	var set struct {
		Questions []struct{ Question, Answer string } `json:"questions"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &set); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(set.Questions) == 0 {
		t.Fatal("practice returned no questions after every variant was excluded — the Learn " +
			"feed treats an empty list as success and stops appending, stranding the learner")
	}
	// And they must still be answerable: a repeat with no anchor is a 409 waiting to happen.
	for _, q := range set.Questions {
		if q.Answer == "" {
			t.Errorf("served question %q has no answer", q.Question)
		}
	}
}

// tokenFromGuest pulls the bearer token out of a /api/auth/guest response.
func tokenFromGuest(t *testing.T, rec *httptest.ResponseRecorder) string {
	t.Helper()
	var out struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil || out.Token == "" {
		t.Fatalf("no guest token in response: %s", rec.Body.String())
	}
	return out.Token
}

// GET /api/next is the single answer to "what should I do now?".
//
// This route had no HTTP-level test, which is the gap that let `/learn` and `/profile` each
// recompute the decision in the browser while internal/scheduler sat unreachable between them.
// The assertions are about the contract the client depends on: a primary plus alternatives,
// wording owned by the server, and no internal terminology in what a learner will read.
func TestNext_ReturnsAPrimaryAndAlternatives(t *testing.T) {
	_, mux, _ := guestServer(t)
	rec := postGuest(t, mux, `{}`, "")
	if rec.Code != 200 {
		t.Fatalf("guest: %d %s", rec.Code, rec.Body.String())
	}
	tok := tokenFromGuest(t, rec)

	rec = httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/next", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	mux.ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatalf("/api/next: %d %s", rec.Code, rec.Body.String())
	}

	var resp struct {
		Primary *struct {
			ID           string  `json:"id"`
			ConceptID    string  `json:"conceptId"`
			ConceptTitle string  `json:"conceptTitle"`
			Kind         string  `json:"kind"`
			Reason       string  `json:"reason"`
			Priority     float64 `json:"priority"`
			Badge        string  `json:"badge"`
			Detail       string  `json:"detail"`
			CTA          string  `json:"cta"`
			Action       struct {
				Type string `json:"type"`
				Href string `json:"href"`
			} `json:"action"`
		} `json:"primary"`
		Alternatives []struct {
			ConceptID string `json:"conceptId"`
			Badge     string `json:"badge"`
			CTA       string `json:"cta"`
			Action    struct {
				Href string `json:"href"`
			} `json:"action"`
		} `json:"alternatives"`
		GeneratedAt string `json:"generatedAt"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v (%s)", err, rec.Body.String())
	}

	if resp.Primary == nil {
		t.Fatalf("no primary recommendation: %s", rec.Body.String())
	}
	if resp.Primary.Badge == "" || resp.Primary.Detail == "" || resp.Primary.CTA == "" {
		t.Errorf("primary carries incomplete copy: %+v", resp.Primary)
	}
	if resp.Primary.Action.Href == "" {
		t.Error("primary has no href")
	}
	if resp.Alternatives == nil {
		t.Error("alternatives is null; the client iterates it without a null check")
	}
	// The one the client must be able to rely on: the head is a /learn entry.
	if resp.Primary.Kind == "learn" && !strings.HasPrefix(resp.Primary.Action.Href, "/learn?concept=") &&
		resp.Primary.Action.Href != "/onboard" && resp.Primary.Action.Href != "/study" {
		t.Errorf("learn href = %q", resp.Primary.Action.Href)
	}

	// No duplicates across head + alternatives.
	seen := map[string]bool{resp.Primary.ConceptID: true}
	for _, a := range resp.Alternatives {
		if a.ConceptID == "" {
			continue
		}
		if seen[a.ConceptID] {
			t.Errorf("duplicate concept %q in the recommendation set", a.ConceptID)
		}
		seen[a.ConceptID] = true
		if a.Badge == "" || a.CTA == "" || a.Action.Href == "" {
			t.Errorf("alternative %q has incomplete copy: %+v", a.ConceptID, a)
		}
	}

	// Learner-facing copy must not name the machinery.
	for _, f := range []string{resp.Primary.Badge, resp.Primary.Detail, resp.Primary.CTA} {
		for _, bad := range []string{"frontier", "DAG", "scheduler", "candidate"} {
			if strings.Contains(f, bad) {
				t.Errorf("learner-facing %q contains internal term %q", f, bad)
			}
		}
	}
}

// A finished concept must not come back as the next thing to learn, which is what `exclude`
// exists for and what the client's `excludeConceptIds` used to do.
func TestNext_HonoursExclude(t *testing.T) {
	_, mux, _ := guestServer(t)
	rec := postGuest(t, mux, `{}`, "")
	tok := tokenFromGuest(t, rec)

	rec = httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/next?exclude=a", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	mux.ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatalf("/api/next?exclude=a: %d", rec.Code)
	}
	if strings.Contains(rec.Body.String(), `"conceptId":"a"`) {
		t.Errorf("excluded concept was still recommended: %s", rec.Body.String())
	}
}

// Unauthenticated must not error: a guest landing on the app gets a recommendation, and an
// error would leave the home page with nothing to show.
func TestNext_UnauthenticatedStillAnswers(t *testing.T) {
	_, mux, _ := guestServer(t)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("GET", "/api/next", nil))
	if rec.Code != 200 {
		t.Fatalf("unauthenticated /api/next: %d %s", rec.Code, rec.Body.String())
	}
	var resp struct {
		Primary *struct {
			Action struct {
				Type string `json:"type"`
				Href string `json:"href"`
			} `json:"action"`
		} `json:"primary"`
		Alternatives []json.RawMessage `json:"alternatives"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Primary == nil || resp.Primary.Action.Href == "" {
		t.Errorf("no primary for a signed-out visitor: %s", rec.Body.String())
	}
}
