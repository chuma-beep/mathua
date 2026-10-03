package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/auth"
	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/engine"
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/planning"
	"github.com/chuma-beep/mathua/internal/storage"
)

// wipeProbe is the learner-visible surface of a reset: XP totals, concept
// progress and attempt history must all read as a fresh start, while the
// account itself survives with its name intact.
type wipeProbe struct {
	XPTotal      int
	XPToday      int
	DiagDone     bool
	Name         string
	ProgressRows int
	Attempts     int
}

func probeWipeState(t *testing.T, store storage.Repository, sid string) wipeProbe {
	t.Helper()
	st, err := store.GetStudent(sid)
	if err != nil || st == nil {
		t.Fatalf("get student %s: %v", sid, err)
	}
	prog, _ := store.GetAllProgress(sid)
	atts, _ := store.GetAttemptsForStudent(sid)
	return wipeProbe{
		XPTotal:      st.XPTotal,
		XPToday:      st.XPToday,
		DiagDone:     st.DiagnosticCompleted,
		Name:         st.Name,
		ProgressRows: len(prog),
		Attempts:     len(atts),
	}
}

// seedWipeState gives the student a full learning record, the way a learner
// who has actually used the product would have one.
func seedWipeState(t *testing.T, store storage.Repository, sid string) {
	t.Helper()
	if err := store.UpsertProgress(&storage.ConceptProgress{StudentID: sid, ConceptID: "a", Status: "PRACTICING", Streak: 3}); err != nil {
		t.Fatalf("seed progress: %v", err)
	}
	sess, err := store.CreateSession(sid)
	if err != nil {
		t.Fatalf("seed session: %v", err)
	}
	if err := store.RecordAttempt(storage.AttemptEntry{SessionID: sess.ID, StudentID: sid, ConceptID: "a", Answer: "1", Expected: "2", Correct: false, ElapsedSeconds: 4, Timestamp: time.Now().UTC()}); err != nil {
		t.Fatalf("seed attempt: %v", err)
	}
	if err := store.UpsertTopicSpeed(&storage.TopicSpeed{StudentID: sid, ConceptID: "a", EFactor: 2.5, Interval: 1, Repetitions: 1, LearningSpeed: 1.0}); err != nil {
		t.Fatalf("seed speed: %v", err)
	}
	if err := store.RecordQuizCompletion(sid, 60); err != nil {
		t.Fatalf("seed quiz: %v", err)
	}
	if err := store.AddXP(sid, 42); err != nil {
		t.Fatalf("seed xp: %v", err)
	}
}

func wipeTestStores(t *testing.T) map[string]storage.Repository {
	t.Helper()
	out := map[string]storage.Repository{}
	sq, _ := storage.NewSQLiteStore(":memory:")
	t.Cleanup(func() { sq.Close() })
	out["sqlite"] = sq
	if dsn := os.Getenv("TEST_POSTGRES_DSN"); dsn != "" {
		pg, err := storage.NewPostgresStore(dsn)
		if err != nil {
			t.Fatalf("open postgres: %v", err)
		}
		t.Cleanup(func() { pg.Close() })
		out["postgres"] = pg
	} else {
		t.Log("TEST_POSTGRES_DSN not set; Postgres leg skipped")
	}
	return out
}

func wipeTestMux(t *testing.T, store storage.Repository) *http.ServeMux {
	t.Helper()
	d, _ := concepts.Build([]concepts.Concept{
		{ID: "a", Label: "A", Domain: "d", GradingType: "numeric", Prerequisites: []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60}},
	})
	reg := generator.NewRegistry()
	reg.Register("a", &testGen{})
	s := New(engine.New(store, d, reg, nil, nil), store, auth.New(store))
	mux := http.NewServeMux()
	s.Register(mux)
	return mux
}

// POST /api/account/reset had no request-level test at all: no 405, no
// phrase-400, no {"reset":true}, and no assertion that the learning record is
// actually gone afterwards. This drives the real handler over the real auth
// middleware, for an account and for a guest, on both stores.
func TestAccountReset_WipesLearningRecord(t *testing.T) {
	for tag, store := range wipeTestStores(t) {
		t.Run(tag, func(t *testing.T) {
			runAccountResetWipe(t, store, wipeTestMux(t, store))
		})
	}
}

func runAccountResetWipe(t *testing.T, store storage.Repository, mux *http.ServeMux) {
	sfx := fmt.Sprintf("%d", time.Now().UnixNano()%100000000)
	authSvc := auth.New(store)

	post := func(token, body string) *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/account/reset", bytes.NewReader([]byte(body)))
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		mux.ServeHTTP(rec, req)
		return rec
	}

	// --- Account holder -----------------------------------------------------
	token, st, err := authSvc.Signup("Resetter", "resetter"+sfx, "Engine!n1")
	if err != nil || token == "" {
		t.Fatalf("signup: %v", err)
	}
	seedWipeState(t, store, st.ID)
	if p := probeWipeState(t, store, st.ID); p.XPTotal == 0 || p.ProgressRows == 0 || p.Attempts == 0 {
		t.Fatalf("seed did not take: %+v", p)
	}

	// Wrong phrase → 400, nothing wiped.
	if rec := post(token, `{"phrase":"nope"}`); rec.Code != 400 {
		t.Fatalf("bad phrase: expected 400, got %d: %s", rec.Code, rec.Body.String())
	}
	if p := probeWipeState(t, store, st.ID); p.XPTotal == 0 || p.ProgressRows == 0 {
		t.Errorf("bad phrase wiped state anyway: %+v", p)
	}

	// Correct phrase → 200 {"reset":true} and a genuinely fresh start.
	rec := post(token, `{"phrase":"reset my progress"}`)
	if rec.Code != 200 {
		t.Fatalf("reset: expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var res map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if res["reset"] != true {
		t.Errorf(`expected {"reset":true}, got %v`, res)
	}
	if p := probeWipeState(t, store, st.ID); p.XPTotal != 0 || p.XPToday != 0 || p.DiagDone || p.ProgressRows != 0 || p.Attempts != 0 {
		t.Errorf("learning record survives reset: %+v", p)
	}
	// The account itself must survive, credentials intact.
	if p := probeWipeState(t, store, st.ID); p.Name != "Resetter" {
		t.Errorf("account must survive reset: %+v", p)
	}
	if _, _, err := authSvc.Login("resetter"+sfx, "Engine!n1"); err != nil {
		t.Errorf("password must survive reset so the learner can sign back in: %v", err)
	}
	// Retry-safe.
	if rec := post(token, `{"phrase":"reset my progress"}`); rec.Code != 200 {
		t.Errorf("retry: expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	// Wrong method → 405.
	mrec := httptest.NewRecorder()
	mreq := httptest.NewRequest("DELETE", "/api/account/reset", bytes.NewReader([]byte(`{}`)))
	mreq.Header.Set("Authorization", "Bearer "+token)
	mux.ServeHTTP(mrec, mreq)
	if mrec.Code != 405 {
		t.Errorf("method: expected 405, got %d", mrec.Code)
	}

	// --- Guest --------------------------------------------------------------
	gr := postGuest(t, mux, fmt.Sprintf(`{"student_id":"guest_%s"}`, sfx), "")
	var g struct {
		Token     string `json:"token"`
		StudentID string `json:"student_id"`
	}
	if err := json.Unmarshal(gr.Body.Bytes(), &g); err != nil || g.Token == "" {
		t.Fatalf("guest: %v %s", err, gr.Body.String())
	}
	seedWipeState(t, store, g.StudentID)
	if p := probeWipeState(t, store, g.StudentID); p.XPTotal == 0 {
		t.Fatalf("guest seed did not take: %+v", p)
	}
	if rec := post(g.Token, `{"phrase":"reset my progress"}`); rec.Code != 200 {
		t.Fatalf("guest reset: expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	if p := probeWipeState(t, store, g.StudentID); p.XPTotal != 0 || p.ProgressRows != 0 || p.Attempts != 0 {
		t.Errorf("guest learning record survives reset: %+v", p)
	}
}

// A reset must leave the learner at zero. But the diagnostic and quiz
// sessions live in Server memory (s.diagSessions / s.quizSessions) and are
// swept only after an hour, so a tab left mid-diagnostic can answer *after* the
// wipe: handleGoalDiagnosticAnswer calls EnsureSession + RecordAttempt, which
// recreates the attempts and sessions rows ResetProgress just deleted.
// Observable symptom: the learner resets, goes back to /profile, and their
// attempts and progress are partly back.
func TestAccountReset_InFlightDiagnosticCannotResurrectProgress(t *testing.T) {
	// guestServer passes a nil planner; the goal diagnostic needs one.
	d, _ := concepts.Build([]concepts.Concept{
		{ID: "a", Label: "A", Domain: "d", GradingType: "numeric", Prerequisites: []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60}},
	})
	store, _ := storage.NewSQLiteStore(":memory:")
	t.Cleanup(func() { store.Close() })
	reg := generator.NewRegistry()
	reg.Register("a", &testGen{})
	srv := New(engine.New(store, d, reg, nil, planning.New(d)), store, auth.New(store))
	mux := http.NewServeMux()
	srv.Register(mux)
	token, sid := quizGuest(t, mux)

	post := func(path, body string) *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", path, bytes.NewReader([]byte(body)))
		req.Header.Set("Authorization", "Bearer "+token)
		mux.ServeHTTP(rec, req)
		return rec
	}

	// Start a diagnostic but do not answer it yet.
	start := post("/api/goal/diagnostic", `{"concept_ids":["a"]}`)
	if start.Code != 200 {
		t.Fatalf("diag start: %d %s", start.Code, start.Body.String())
	}
	var sres struct {
		SessionID string `json:"session_id"`
		Done      bool   `json:"done"`
	}
	if err := json.Unmarshal(start.Body.Bytes(), &sres); err != nil {
		t.Fatalf("decode start: %v", err)
	}
	if sres.Done || sres.SessionID == "" {
		t.Fatalf("expected a live diagnostic question, got %s", start.Body.String())
	}

	// Give the learner a record, then wipe it.
	seedWipeState(t, store, sid)
	if rec := post("/api/account/reset", `{"phrase":"reset my progress"}`); rec.Code != 200 {
		t.Fatalf("reset: %d %s", rec.Code, rec.Body.String())
	}
	if p := probeWipeState(t, store, sid); p.ProgressRows != 0 || p.Attempts != 0 {
		t.Fatalf("precondition: record not wiped: %+v", p)
	}

	// The abandoned tab now answers. The session must already be gone, so the
	// write is refused rather than silently re-creating the wiped rows.
	ans := post("/api/goal/diagnostic/answer",
		`{"session_id":"`+sres.SessionID+`","concept_id":"a","answer":"1","elapsed":5}`)
	if ans.Code != http.StatusNotFound {
		t.Errorf("stale diagnostic session should be evicted by reset, got %d: %s", ans.Code, ans.Body.String())
	}
	if p := probeWipeState(t, store, sid); p.ProgressRows != 0 || p.Attempts != 0 {
		t.Errorf("progress resurrected after reset: %+v", p)
	}
}

// Same hole, wider blast radius: a quiz answer writes concept_progress AND XP,
// so the learner's PositionBlock, per-domain report and daily goal all come
// back after a reset they had already wiped.
func TestAccountReset_InFlightQuizCannotResurrectProgress(t *testing.T) {
	_, mux, store := guestServer(t)
	token, sid := quizGuest(t, mux)

	post := func(path, body string) *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", path, bytes.NewReader([]byte(body)))
		req.Header.Set("Authorization", "Bearer "+token)
		mux.ServeHTTP(rec, req)
		return rec
	}

	sr := post("/api/quiz/session", `{}`)
	if sr.Code != 200 {
		t.Fatalf("quiz start: %d %s", sr.Code, sr.Body.String())
	}
	var qs struct {
		SessionID string `json:"session_id"`
	}
	if err := json.Unmarshal(sr.Body.Bytes(), &qs); err != nil || qs.SessionID == "" {
		t.Fatalf("decode quiz session: %v %s", err, sr.Body.String())
	}

	seedWipeState(t, store, sid)
	if rec := post("/api/account/reset", `{"phrase":"reset my progress"}`); rec.Code != 200 {
		t.Fatalf("reset: %d %s", rec.Code, rec.Body.String())
	}
	if p := probeWipeState(t, store, sid); p.ProgressRows != 0 || p.Attempts != 0 {
		t.Fatalf("precondition: record not wiped: %+v", p)
	}

	ans := post("/api/quiz/answer",
		`{"session_id":"`+qs.SessionID+`","concept_id":"a","answer":"1","elapsed":5}`)
	if ans.Code != http.StatusNotFound {
		t.Errorf("stale quiz session should be evicted by reset, got %d: %s", ans.Code, ans.Body.String())
	}
	if p := probeWipeState(t, store, sid); p.ProgressRows != 0 || p.Attempts != 0 || p.XPTotal != 0 {
		t.Errorf("quiz resurrected progress after reset: %+v", p)
	}
}

// Worst case: after DELETE the student row is gone, so an in-flight answer has
// no owner to write for. It must not silently recreate the wiped tables for an
// identity that no longer exists.
func TestAccountDelete_InFlightDiagnosticCannotResurrectRows(t *testing.T) {
	for tag, store := range wipeTestStores(t) {
		t.Run(tag, func(t *testing.T) { runDeleteInFlightWipe(t, store) })
	}
}

func runDeleteInFlightWipe(t *testing.T, store storage.Repository) {
	d, _ := concepts.Build([]concepts.Concept{
		{ID: "a", Label: "A", Domain: "d", GradingType: "numeric", Prerequisites: []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60}},
	})
	reg := generator.NewRegistry()
	reg.Register("a", &testGen{})
	srv := New(engine.New(store, d, reg, nil, planning.New(d)), store, auth.New(store))
	mux := http.NewServeMux()
	srv.Register(mux)
	token, sid := quizGuest(t, mux)

	call := func(method, path, body string) *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(method, path, bytes.NewReader([]byte(body)))
		req.Header.Set("Authorization", "Bearer "+token)
		mux.ServeHTTP(rec, req)
		return rec
	}

	start := call("POST", "/api/goal/diagnostic", `{"concept_ids":["a"]}`)
	var sres struct {
		SessionID string `json:"session_id"`
		Done      bool   `json:"done"`
	}
	if err := json.Unmarshal(start.Body.Bytes(), &sres); err != nil || sres.Done {
		t.Fatalf("diag start: %v %s", err, start.Body.String())
	}

	if rec := call("DELETE", "/api/account", `{"phrase":"delete my account"}`); rec.Code != 200 {
		t.Fatalf("delete: %d %s", rec.Code, rec.Body.String())
	}
	if st, _ := store.GetStudent(sid); st != nil {
		t.Fatalf("precondition: student row survives delete")
	}

	call("POST", "/api/goal/diagnostic/answer",
		`{"session_id":"`+sres.SessionID+`","concept_id":"a","answer":"1","elapsed":5}`)

	if a, _ := store.GetAttemptsForStudent(sid); len(a) != 0 {
		t.Errorf("attempts resurrected for a deleted account: %d rows", len(a))
	}
	if p, _ := store.GetAllProgress(sid); len(p) != 0 {
		t.Errorf("progress resurrected for a deleted account: %d rows", len(p))
	}
}
