package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/auth"
	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/engine"
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/storage"
)

// gatedServer builds a two-node DAG (a -> b) with a guest-capable server, so a
// locked topic can be exercised over HTTP.
func gatedServer(t *testing.T) (*http.ServeMux, storage.Repository) {
	t.Helper()
	d, _ := concepts.Build([]concepts.Concept{
		{ID: "a", Label: "A", Domain: "d", GradingType: "numeric", Prerequisites: []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60}},
		{ID: "b", Label: "B", Domain: "d", GradingType: "numeric", Prerequisites: []string{"a"},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60}},
	})
	store, _ := storage.NewSQLiteStore(":memory:")
	t.Cleanup(func() { store.Close() })
	reg := generator.NewRegistry()
	reg.Register("a", &testGen{})
	reg.Register("b", &testGen{})
	s := New(engine.New(store, d, reg, nil, nil), store, auth.New(store))
	mux := http.NewServeMux()
	s.Register(mux)
	return mux, store
}

func guestTokenAndID(t *testing.T, mux *http.ServeMux) (string, string) {
	t.Helper()
	rec := postGuest(t, mux, `{}`, "")
	if rec.Code != 200 {
		t.Fatalf("guest: %d %s", rec.Code, rec.Body.String())
	}
	var res struct {
		Token     string `json:"token"`
		StudentID string `json:"student_id"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("decode guest: %v", err)
	}
	return res.Token, res.StudentID
}

func authedGet(t *testing.T, mux *http.ServeMux, path, token string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", path, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	mux.ServeHTTP(rec, req)
	return rec
}

// A locked topic cannot be started through /learn, even when the request reaches
// the practice endpoint directly. This is invariant 3: the frontend's choice does
// not stand in for demonstrated knowledge.
func TestPractice_LockedTopicIsRefused(t *testing.T) {
	mux, _ := gatedServer(t)
	token, _ := guestTokenAndID(t, mux)

	rec := authedGet(t, mux, "/api/lessons/b/practice", token)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("locked practice = %d, want 403: %s", rec.Code, rec.Body.String())
	}
	var body struct {
		Error struct {
			Code    string `json:"code"`
			Details struct {
				PrerequisiteConceptIDs []string `json:"prerequisiteConceptIds"`
			} `json:"details"`
		} `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode error: %v (%s)", err, rec.Body.String())
	}
	if body.Error.Code != "TOPIC_LOCKED" {
		t.Errorf("code = %q, want TOPIC_LOCKED", body.Error.Code)
	}
	foundA := false
	for _, id := range body.Error.Details.PrerequisiteConceptIDs {
		if id == "a" {
			foundA = true
		}
	}
	if !foundA {
		t.Errorf("prerequisiteConceptIds = %v, want to include a", body.Error.Details.PrerequisiteConceptIDs)
	}
}

// And the answer path refuses too, so no route records evidence for a topic the
// learner was not ready for.
func TestStudyAnswer_LockedTopicIsRefused(t *testing.T) {
	mux, _ := gatedServer(t)
	token, _ := guestTokenAndID(t, mux)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/study/answer", bytes.NewReader([]byte(
		`{"concept_id":"b","answer":"4","elapsed":3}`)))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("locked study answer = %d, want 403: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "TOPIC_LOCKED") {
		t.Errorf("body = %s, want TOPIC_LOCKED", rec.Body.String())
	}
}

// Placement makes the successor startable without writing mastery.
func TestPractice_PlacementUnlocksSuccessor(t *testing.T) {
	mux, store := gatedServer(t)
	token, sid := guestTokenAndID(t, mux)

	if err := store.UpsertProgress(&storage.ConceptProgress{
		StudentID: sid, ConceptID: "a", Status: "PRACTICING",
		SM2EFactor: 2.5, PlacementSeeded: true,
	}); err != nil {
		t.Fatalf("seed placement: %v", err)
	}

	rec := authedGet(t, mux, "/api/lessons/b/readiness", token)
	if rec.Code != 200 {
		t.Fatalf("readiness = %d: %s", rec.Code, rec.Body.String())
	}
	var body struct {
		Eligible bool   `json:"eligible"`
		State    string `json:"state"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if !body.Eligible || body.State != "unlocked" {
		t.Errorf("readiness = %+v, want eligible unlocked via placement", body)
	}

	if rec := authedGet(t, mux, "/api/lessons/b/practice", token); rec.Code != 200 {
		t.Fatalf("placement-unlocked practice = %d, want 200: %s", rec.Code, rec.Body.String())
	}
}

// A completed diagnostic is never offered back as the active task, even when the
// seeded rows leave mastery at zero.
func TestNext_CompletedDiagnosticIsNotAnActiveTask(t *testing.T) {
	mux, store := gatedServer(t)
	token, sid := guestTokenAndID(t, mux)

	if err := store.SetDiagnosticCompleted(sid); err != nil {
		t.Fatalf("mark completed: %v", err)
	}
	rec := authedGet(t, mux, "/api/next", token)
	if rec.Code != 200 {
		t.Fatalf("/api/next = %d: %s", rec.Code, rec.Body.String())
	}
	var body struct {
		Primary *struct {
			ID     string `json:"id"`
			Action struct {
				Href string `json:"href"`
			} `json:"action"`
		} `json:"primary"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Primary == nil {
		t.Fatal("no primary recommendation for a placed learner")
	}
	if body.Primary.ID == "diagnostic" || body.Primary.Action.Href == "/onboard" {
		t.Errorf("completed diagnostic offered as active task: %+v", body.Primary)
	}
	// And the learner gets an actual concept to learn — the placement fix is only
	// useful if real tasks replace the diagnostic.
	if !strings.Contains(body.Primary.Action.Href, "concept=a") {
		t.Errorf("primary = %+v, want a real concept task", body.Primary)
	}
}

// The curriculum browser reads state from the server; these endpoints are that
// source, and they use the same eligibility rule as /learn.
func TestCurriculumDomains_StatesComeFromTheServer(t *testing.T) {
	mux, _ := gatedServer(t)
	token, _ := guestTokenAndID(t, mux)

	rec := authedGet(t, mux, "/api/curriculum/domains", token)
	if rec.Code != 200 {
		t.Fatalf("/api/curriculum/domains = %d: %s", rec.Code, rec.Body.String())
	}
	var body struct {
		Domains []struct {
			ID              string `json:"id"`
			ConceptCount    int    `json:"conceptCount"`
			UnlockedCount   int    `json:"unlockedCount"`
			InProgressCount int    `json:"inProgressCount"`
		} `json:"domains"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(body.Domains) != 1 || body.Domains[0].ID != "d" {
		t.Fatalf("domains = %+v, want one domain d", body.Domains)
	}
	d := body.Domains[0]
	if d.ConceptCount != 2 || d.UnlockedCount != 1 || d.InProgressCount != 0 {
		t.Errorf("counts = %+v, want 2 concepts and 1 unlocked", d)
	}
}

func TestCurriculumDomain_TopicsCarryState(t *testing.T) {
	mux, _ := gatedServer(t)
	token, _ := guestTokenAndID(t, mux)

	rec := authedGet(t, mux, "/api/curriculum/domains/d", token)
	if rec.Code != 200 {
		t.Fatalf("/api/curriculum/domains/d = %d: %s", rec.Code, rec.Body.String())
	}
	var body struct {
		Domain struct {
			ID string `json:"id"`
		} `json:"domain"`
		Topics []struct {
			ConceptID string `json:"conceptId"`
			State     string `json:"state"`
		} `json:"topics"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	states := map[string]string{}
	for _, tp := range body.Topics {
		states[tp.ConceptID] = tp.State
	}
	if states["a"] != "unlocked" || states["b"] != "locked" {
		t.Errorf("states = %v, want a unlocked and b locked", states)
	}
}
