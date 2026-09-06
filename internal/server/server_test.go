package server

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/auth"
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
	if rec := post("/api/auth/signup", `{"name":"Ada","username":"Ada","password":"Engine!n1"}`, ""); rec.Code != 200 {
		t.Fatalf("signup: expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	rec := post("/api/auth/signup", `{"name":"Other","username":"ADA","password":"Engine!n2"}`, "")
	if rec.Code != 409 {
		t.Errorf("duplicate: expected 409, got %d: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "username is taken") {
		t.Errorf("expected friendly taken message, got %s", rec.Body.String())
	}
	// Bad username / weak password → 400, never a raw dump.
	if rec := post("/api/auth/signup", `{"name":"X","username":"ab","password":"Engine!n1"}`, ""); rec.Code != 400 {
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

	rec := put(token, `{"name":"  Ada Lovelace  "}`)
	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
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
