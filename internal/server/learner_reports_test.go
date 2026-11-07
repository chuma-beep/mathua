package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/auth"
	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/engine"
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/storage"
)

func reportsTestMux(t *testing.T) (*Server, *http.ServeMux, storage.Repository) {
	t.Helper()
	dag, err := concepts.Build([]concepts.Concept{
		{ID: "a", Label: "Alpha", Domain: "d", GradingType: "numeric", Prerequisites: []string{}, MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60}},
		{ID: "b", Label: "Beta", Domain: "e", GradingType: "numeric", Prerequisites: []string{}, MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60}},
	})
	if err != nil {
		t.Fatalf("build dag: %v", err)
	}
	store, err := storage.NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("new store: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	reg := generator.NewRegistry()
	reg.Register("a", &seededGen{})
	reg.Register("b", &seededGen{})
	s := New(engine.New(store, dag, reg, nil, nil), store, auth.New(store))
	mux := http.NewServeMux()
	s.Register(mux)
	return s, mux, store
}

func signupReportLearner(t *testing.T, s *Server, name, username string) (string, *storage.Student) {
	t.Helper()
	token, student, err := s.auth.Signup(name, username, "Engine!n1")
	if err != nil {
		t.Fatalf("signup %s: %v", username, err)
	}
	return token, student
}

func addReportAttempt(t *testing.T, store storage.Repository, sid, conceptID, source string, correct bool, at time.Time) {
	t.Helper()
	sess, err := store.CreateSession(sid)
	if err != nil {
		t.Fatalf("create session: %v", err)
	}
	if _, err := store.RecordAttempt(storage.AttemptEntry{
		SessionID: sess.ID, StudentID: sid, ConceptID: conceptID,
		Answer: "x", Expected: "x", Correct: correct,
		ElapsedSeconds: 5, Timestamp: at.UTC(), Question: "q", Source: source,
	}); err != nil {
		t.Fatalf("record attempt: %v", err)
	}
}

func getReport(t *testing.T, mux *http.ServeMux, token, target string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, target, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	return rr
}

func decodeReport(t *testing.T, rr *httptest.ResponseRecorder) learnerReportResponse {
	t.Helper()
	if rr.Code != http.StatusOK {
		t.Fatalf("status %d: %s", rr.Code, rr.Body.String())
	}
	var out learnerReportResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return out
}

func TestLearnerReportRequiresAuthentication(t *testing.T) {
	_, mux, _ := reportsTestMux(t)
	if rr := getReport(t, mux, "", "/api/learner/reports"); rr.Code != http.StatusUnauthorized {
		t.Fatalf("anonymous report status = %d, want 401", rr.Code)
	}
	if rr := getReport(t, mux, "", "/api/learner/reports/history"); rr.Code != http.StatusUnauthorized {
		t.Fatalf("anonymous history status = %d, want 401", rr.Code)
	}
}

func TestLearnerReportRejectsNonGet(t *testing.T) {
	s, mux, _ := reportsTestMux(t)
	token, _ := signupReportLearner(t, s, "Writer", "writer")
	req := httptest.NewRequest(http.MethodPost, "/api/learner/reports", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405", rr.Code)
	}
}

func TestLearnerReportCountsAndSections(t *testing.T) {
	s, mux, store := reportsTestMux(t)
	token, student := signupReportLearner(t, s, "Ada", "ada")
	now := time.Now().UTC()
	masteredAt := now.Add(-2 * time.Hour)
	if err := store.UpsertProgress(&storage.ConceptProgress{StudentID: student.ID, ConceptID: "a", Status: "MASTERED", Streak: 3, MasteredAt: &masteredAt}); err != nil {
		t.Fatalf("progress: %v", err)
	}
	addReportAttempt(t, store, student.ID, "a", "quiz", true, now.Add(-90*time.Minute))
	addReportAttempt(t, store, student.ID, "a", "practice", true, now.Add(-80*time.Minute))
	if err := store.AwardXP(student.ID, 42, "test", "report-42"); err != nil {
		t.Fatalf("award xp: %v", err)
	}
	if err := store.SetDiagnosticCompleted(student.ID); err != nil {
		t.Fatalf("diagnostic: %v", err)
	}
	if err := store.RecordQuizCompletion(student.ID, 42); err != nil {
		t.Fatalf("quiz completion: %v", err)
	}

	out := decodeReport(t, getReport(t, mux, token, "/api/learner/reports"))
	if out.Overview.Total != 2 {
		t.Fatalf("total = %d, want 2", out.Overview.Total)
	}
	if out.Overview.Mastered != 1 {
		t.Fatalf("mastered = %d, want 1", out.Overview.Mastered)
	}
	if out.Overview.XPLifetime != 42 {
		t.Fatalf("xp lifetime = %d, want 42", out.Overview.XPLifetime)
	}
	if len(out.Mastery) != 2 {
		t.Fatalf("mastery rows = %d, want 2", len(out.Mastery))
	}
	var stateA string
	for _, row := range out.Mastery {
		if row.ConceptID == "a" {
			stateA = row.State
		}
	}
	if stateA != "mastered" {
		t.Fatalf("state for a = %q, want mastered", stateA)
	}
	if len(out.Assessments) < 2 {
		t.Fatalf("assessments = %d, want quiz and diagnostic", len(out.Assessments))
	}
	if len(out.History) == 0 || out.HistoryTotal == 0 {
		t.Fatalf("history empty: rows=%d total=%d", len(out.History), out.HistoryTotal)
	}
	if len(out.Limitations) == 0 {
		t.Fatal("limitations should record what the current schema cannot show")
	}
}

func TestLearnerReportIsScopedToTheLearner(t *testing.T) {
	s, mux, store := reportsTestMux(t)
	tokenA, ada := signupReportLearner(t, s, "Ada", "ada")
	_, bob := signupReportLearner(t, s, "Bob", "bob")
	now := time.Now().UTC()
	masteredAt := now.Add(-time.Hour)
	if err := store.UpsertProgress(&storage.ConceptProgress{StudentID: bob.ID, ConceptID: "a", Status: "MASTERED", Streak: 3, MasteredAt: &masteredAt}); err != nil {
		t.Fatalf("bob progress: %v", err)
	}
	if err := store.AwardXP(bob.ID, 500, "test", "bob-500"); err != nil {
		t.Fatalf("bob xp: %v", err)
	}

	out := decodeReport(t, getReport(t, mux, tokenA, "/api/learner/reports"))
	if out.Overview.Mastered != 0 {
		t.Fatalf("ada mastered = %d, want 0", out.Overview.Mastered)
	}
	if out.Overview.XPLifetime != 0 {
		t.Fatalf("ada xp = %d, want 0", out.Overview.XPLifetime)
	}
	for _, row := range out.Mastery {
		if row.State == "mastered" {
			t.Fatalf("ada sees mastered concept %s", row.ConceptID)
		}
	}
	if out.HistoryTotal != 0 {
		t.Fatalf("ada history total = %d, want 0", out.HistoryTotal)
	}
	_ = ada
}

func TestLearnerReportDomainFilter(t *testing.T) {
	s, mux, store := reportsTestMux(t)
	token, student := signupReportLearner(t, s, "Ada", "ada")
	now := time.Now().UTC()
	masteredAt := now.Add(-time.Hour)
	store.UpsertProgress(&storage.ConceptProgress{StudentID: student.ID, ConceptID: "a", Status: "MASTERED", Streak: 3, MasteredAt: &masteredAt})
	store.AwardXP(student.ID, 20, "test", "x1")

	out := decodeReport(t, getReport(t, mux, token, "/api/learner/reports?domain=d"))
	if out.Overview.Total != 1 || out.Overview.Mastered != 1 {
		t.Fatalf("domain overview = total %d mastered %d, want 1/1", out.Overview.Total, out.Overview.Mastered)
	}
	if len(out.Domains) != 1 || out.Domains[0].ID != "d" {
		t.Fatalf("domains = %+v, want only d", out.Domains)
	}
	if len(out.Mastery) != 1 || out.Mastery[0].Domain != "d" {
		t.Fatalf("mastery = %+v, want only domain d", out.Mastery)
	}
	if out.Overview.XPLifetime != 20 {
		t.Fatalf("xp lifetime = %d, want 20 even with a domain filter", out.Overview.XPLifetime)
	}
}

func TestLearnerReportRejectsInvalidPeriod(t *testing.T) {
	s, mux, _ := reportsTestMux(t)
	token, _ := signupReportLearner(t, s, "Ada", "ada")
	if rr := getReport(t, mux, token, "/api/learner/reports?days=5"); rr.Code != http.StatusBadRequest {
		t.Fatalf("days=5 status = %d, want 400", rr.Code)
	}
	if rr := getReport(t, mux, token, "/api/learner/reports?days=abc"); rr.Code != http.StatusBadRequest {
		t.Fatalf("days=abc status = %d, want 400", rr.Code)
	}
}

func TestLearnerReportHistoryPaginatesWithoutOverlap(t *testing.T) {
	s, mux, store := reportsTestMux(t)
	token, student := signupReportLearner(t, s, "Ada", "ada")
	now := time.Now().UTC()
	for i := 0; i < 5; i++ {
		addReportAttempt(t, store, student.ID, "a", "practice", i%2 == 0, now.Add(-time.Duration(i)*time.Minute))
	}
	first := getReport(t, mux, token, "/api/learner/reports/history?limit=2&offset=0")
	if first.Code != http.StatusOK {
		t.Fatalf("first page status %d: %s", first.Code, first.Body.String())
	}
	var p1 learnerHistoryResponse
	if err := json.Unmarshal(first.Body.Bytes(), &p1); err != nil {
		t.Fatalf("decode first: %v", err)
	}
	second := getReport(t, mux, token, "/api/learner/reports/history?limit=2&offset=2")
	if second.Code != http.StatusOK {
		t.Fatalf("second page status %d: %s", second.Code, second.Body.String())
	}
	var p2 learnerHistoryResponse
	if err := json.Unmarshal(second.Body.Bytes(), &p2); err != nil {
		t.Fatalf("decode second: %v", err)
	}
	if p1.Total != 5 || p2.Total != 5 {
		t.Fatalf("totals = %d/%d, want 5/5", p1.Total, p2.Total)
	}
	if len(p1.Events) != 2 || len(p2.Events) != 2 {
		t.Fatalf("page sizes = %d/%d, want 2/2", len(p1.Events), len(p2.Events))
	}
	seen := map[string]bool{}
	for _, ev := range append(p1.Events, p2.Events...) {
		if seen[ev.ID] {
			t.Fatalf("event %s appeared on both pages", ev.ID)
		}
		seen[ev.ID] = true
	}
}

func TestLearnerReportHistoryIncludesNegativeXPAdjustment(t *testing.T) {
	s, mux, store := reportsTestMux(t)
	token, student := signupReportLearner(t, s, "Ada", "ada_negative")
	if err := store.AwardXP(student.ID, -5, "accountability:test", "adjustment:1"); err != nil {
		t.Fatalf("award negative: %v", err)
	}
	rr := getReport(t, mux, token, "/api/learner/reports/history?limit=10&offset=0")
	if rr.Code != http.StatusOK {
		t.Fatalf("history status %d: %s", rr.Code, rr.Body.String())
	}
	var page learnerHistoryResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &page); err != nil {
		t.Fatalf("decode history: %v", err)
	}
	for _, ev := range page.Events {
		if ev.Kind == "xp" {
			if ev.XP == nil || *ev.XP != -5 || ev.Title != "XP adjustment" || ev.Source != "accountability:test" {
				t.Fatalf("negative xp event = %+v", ev)
			}
			return
		}
	}
	t.Fatalf("negative XP adjustment missing from history: %+v", page.Events)
}

func TestLearnerReportEmptyLearner(t *testing.T) {
	s, mux, _ := reportsTestMux(t)
	token, _ := signupReportLearner(t, s, "New", "newlearner")
	out := decodeReport(t, getReport(t, mux, token, "/api/learner/reports"))
	if out.Overview.Mastered != 0 || out.Overview.XPLifetime != 0 {
		t.Fatalf("new learner overview = %+v", out.Overview)
	}
	if len(out.Mastery) != 2 {
		t.Fatalf("mastery rows = %d, want the two concepts with locked/available states", len(out.Mastery))
	}
	if len(out.Assessments) != 0 || len(out.History) != 0 || len(out.Activity) != 0 {
		t.Fatalf("new learner has data: assessments=%d history=%d activity=%d", len(out.Assessments), len(out.History), len(out.Activity))
	}
}
