package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/storage"
)

// missingSessionRepo simulates a database where the active session survived
// but its parent sessions row did not. That state is possible on Postgres,
// whose active_sessions table has no foreign key to sessions, and it is the
// state that used to reach RecordAttempt and fail with a raw 500.
type missingSessionRepo struct {
	storage.Repository
}

func (missingSessionRepo) GetSession(string) (*storage.Session, error) {
	return nil, nil
}

func TestReviewsSessionAndAnswerRoutes(t *testing.T) {
	_, mux, store, token, sid := setupTestServerWithPlanner(t)

	// POST /api/reviews/session: Method not allowed
	reqGet := httptest.NewRequest("GET", "/api/reviews/session", nil)
	reqGet.Header.Set("Authorization", "Bearer "+token)
	recGet := httptest.NewRecorder()
	mux.ServeHTTP(recGet, reqGet)
	if recGet.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405 for GET /api/reviews/session, got %d", recGet.Code)
	}

	// POST /api/reviews/session: Unauthenticated
	reqUnauth := httptest.NewRequest("POST", "/api/reviews/session", nil)
	recUnauth := httptest.NewRecorder()
	mux.ServeHTTP(recUnauth, reqUnauth)
	if recUnauth.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for unauthenticated /api/reviews/session, got %d", recUnauth.Code)
	}

	// POST /api/reviews/session: Authenticated
	reqAuth := httptest.NewRequest("POST", "/api/reviews/session", nil)
	reqAuth.Header.Set("Authorization", "Bearer "+token)
	recAuth := httptest.NewRecorder()
	mux.ServeHTTP(recAuth, reqAuth)
	if recAuth.Code != http.StatusOK {
		t.Fatalf("expected 200 for auth /api/reviews/session, got %d: %s", recAuth.Code, recAuth.Body.String())
	}
	var res startSessionRes
	if err := json.Unmarshal(recAuth.Body.Bytes(), &res); err != nil {
		t.Fatalf("unmarshal startSessionRes: %v", err)
	}
	if res.StudentID != sid || res.SessionID == "" {
		t.Errorf("unexpected session response: %+v", res)
	}

	// POST /api/reviews/answer: Method not allowed
	reqAnswerGet := httptest.NewRequest("GET", "/api/reviews/answer", nil)
	reqAnswerGet.Header.Set("Authorization", "Bearer "+token)
	recAnswerGet := httptest.NewRecorder()
	mux.ServeHTTP(recAnswerGet, reqAnswerGet)
	if recAnswerGet.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405 for GET /api/reviews/answer, got %d", recAnswerGet.Code)
	}

	// POST /api/reviews/answer: Invalid JSON
	reqBadJSON := httptest.NewRequest("POST", "/api/reviews/answer", bytes.NewReader([]byte(`invalid`)))
	reqBadJSON.Header.Set("Authorization", "Bearer "+token)
	recBadJSON := httptest.NewRecorder()
	mux.ServeHTTP(recBadJSON, reqBadJSON)
	if recBadJSON.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for bad JSON, got %d", recBadJSON.Code)
	}

	// POST /api/reviews/answer: Elapsed too fast (< MinAnswerSeconds)
	fastPayload := answerReq{
		SessionID: res.SessionID,
		AttemptID: "att-1",
		Answer:    "42",
		Elapsed:   0.1,
	}
	fastBytes, _ := json.Marshal(fastPayload)
	reqFast := httptest.NewRequest("POST", "/api/reviews/answer", bytes.NewReader(fastBytes))
	reqFast.Header.Set("Authorization", "Bearer "+token)
	recFast := httptest.NewRecorder()
	mux.ServeHTTP(recFast, reqFast)
	if recFast.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for elapsed < MinAnswerSeconds, got %d", recFast.Code)
	}

	// POST /api/reviews/answer: Session belonging to different student
	otherSt, _ := store.CreateStudent("other-student")
	otherSess, _ := store.CreateSession(otherSt.ID)
	wrongPayload := answerReq{
		SessionID: otherSess.ID,
		AttemptID: "att-2",
		Answer:    "42",
		Elapsed:   3.0,
	}
	wrongBytes, _ := json.Marshal(wrongPayload)
	reqWrong := httptest.NewRequest("POST", "/api/reviews/answer", bytes.NewReader(wrongBytes))
	reqWrong.Header.Set("Authorization", "Bearer "+token)
	recWrong := httptest.NewRecorder()
	mux.ServeHTTP(recWrong, reqWrong)
	if recWrong.Code != http.StatusForbidden {
		t.Errorf("expected 403 for other student session, got %d", recWrong.Code)
	}

	// POST /api/reviews/answer: No active question in session (409)
	noQPayload := answerReq{
		SessionID: res.SessionID,
		AttemptID: "att-3",
		Answer:    "42",
		Elapsed:   3.0,
	}
	noQBytes, _ := json.Marshal(noQPayload)
	reqNoQ := httptest.NewRequest("POST", "/api/reviews/answer", bytes.NewReader(noQBytes))
	reqNoQ.Header.Set("Authorization", "Bearer "+token)
	recNoQ := httptest.NewRecorder()
	mux.ServeHTTP(recNoQ, reqNoQ)
	if recNoQ.Code != http.StatusConflict {
		t.Errorf("expected 409 for no active question, got %d", recNoQ.Code)
	}
}

// A review answer whose parent sessions row is gone must be refused before the
// engine grades it. Submitting used to rehydrate the surviving active_sessions
// row, grade the answer and then fail the attempts.session_id foreign key with
// a 500. The contract is a 409 and a client restart, with no attempt written.
func TestReviewsAnswer_MissingSessionIs409Not500(t *testing.T) {
	srv, _, store, token, sid := setupTestServerWithPlanner(t)

	sess, err := store.CreateSession(sid)
	if err != nil {
		t.Fatalf("create session: %v", err)
	}
	if err := store.UpsertActiveSession(&storage.ActiveSession{
		SessionID:      sess.ID,
		StudentID:      sid,
		ConceptID:      "a",
		ConceptName:    "Concept A",
		ExpectedAnswer: "4",
		AttemptID:      "att-stale",
		Question:       "2 + 2 = ?",
		IsReview:       true,
	}); err != nil {
		t.Fatalf("seed active session: %v", err)
	}

	before, err := store.GetAttemptsForStudent(sid)
	if err != nil {
		t.Fatalf("attempts before: %v", err)
	}

	mux := http.NewServeMux()
	New(srv.eng, missingSessionRepo{Repository: store}, srv.auth).Register(mux)

	payload, _ := json.Marshal(answerReq{
		SessionID: sess.ID,
		AttemptID: "att-stale",
		Answer:    "4",
		Elapsed:   3.0,
	})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/reviews/answer", bytes.NewReader(payload))
	req.Header.Set("Authorization", "Bearer "+token)
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("expected 409 for missing session, got %d: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "restart the review") {
		t.Errorf("409 should tell the client to restart: %s", rec.Body.String())
	}

	after, err := store.GetAttemptsForStudent(sid)
	if err != nil {
		t.Fatalf("attempts after: %v", err)
	}
	if len(after) != len(before) {
		t.Errorf("missing-session submit wrote an attempt: before %d, after %d", len(before), len(after))
	}
}

func TestCourseDiagnosticRoute(t *testing.T) {
	_, mux, _, token, _ := setupTestServerWithPlanner(t)

	// Method not allowed
	reqGet := httptest.NewRequest("GET", "/api/courses/test-course/diagnostic", nil)
	reqGet.Header.Set("Authorization", "Bearer "+token)
	recGet := httptest.NewRecorder()
	mux.ServeHTTP(recGet, reqGet)
	if recGet.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405 for GET /api/courses/test-course/diagnostic, got %d", recGet.Code)
	}

	// Valid course
	reqValid := httptest.NewRequest("POST", "/api/courses/test-course/diagnostic", nil)
	reqValid.Header.Set("Authorization", "Bearer "+token)
	recValid := httptest.NewRecorder()
	mux.ServeHTTP(recValid, reqValid)
	if recValid.Code != http.StatusOK {
		t.Fatalf("expected 200 for valid course diagnostic, got %d: %s", recValid.Code, recValid.Body.String())
	}
	var res map[string]interface{}
	if err := json.Unmarshal(recValid.Body.Bytes(), &res); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if res["course_id"] != "test-course" {
		t.Errorf("expected course_id test-course, got %v", res["course_id"])
	}

	// Invalid course (400)
	reqInvalid := httptest.NewRequest("POST", "/api/courses/nonexistent-course/diagnostic", nil)
	reqInvalid.Header.Set("Authorization", "Bearer "+token)
	recInvalid := httptest.NewRecorder()
	mux.ServeHTTP(recInvalid, reqInvalid)
	if recInvalid.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for nonexistent course, got %d", recInvalid.Code)
	}
}

func TestGoogleAuthRoutes_Handlers(t *testing.T) {
	_, mux, _, _, _ := setupTestServerWithPlanner(t)

	// 1. GET /api/auth/google/login unconfigured
	origID := os.Getenv("GOOGLE_CLIENT_ID")
	origOauthID := os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	_ = os.Unsetenv("GOOGLE_CLIENT_ID")
	_ = os.Unsetenv("GOOGLE_OAUTH_CLIENT_ID")
	defer func() {
		if origID != "" {
			_ = os.Setenv("GOOGLE_CLIENT_ID", origID)
		}
		if origOauthID != "" {
			_ = os.Setenv("GOOGLE_OAUTH_CLIENT_ID", origOauthID)
		}
	}()

	reqUnconfigured := httptest.NewRequest("GET", "/api/auth/google/login", nil)
	recUnconfigured := httptest.NewRecorder()
	mux.ServeHTTP(recUnconfigured, reqUnconfigured)
	if recUnconfigured.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 when google auth unconfigured, got %d", recUnconfigured.Code)
	}

	// 2. GET /api/auth/google/login configured -> redirects
	_ = os.Setenv("GOOGLE_CLIENT_ID", "mock-google-client-id")
	reqConfigured := httptest.NewRequest("GET", "/api/auth/google/login?return=/dashboard", nil)
	recConfigured := httptest.NewRecorder()
	mux.ServeHTTP(recConfigured, reqConfigured)
	if recConfigured.Code != http.StatusFound {
		t.Fatalf("expected 302 redirect for configured login, got %d", recConfigured.Code)
	}
	loc := recConfigured.Header().Get("Location")
	if !strings.HasPrefix(loc, "https://accounts.google.com/o/oauth2/v2/auth") {
		t.Errorf("unexpected redirect location: %s", loc)
	}
	cookies := recConfigured.Result().Cookies()
	var hasState, hasReturn bool
	var stateVal string
	for _, c := range cookies {
		if c.Name == "oauth_state" && c.Value != "" {
			hasState = true
			stateVal = c.Value
		}
		if c.Name == "oauth_return" && c.Value == "/dashboard" {
			hasReturn = true
		}
	}
	if !hasState || !hasReturn {
		t.Errorf("expected oauth_state and oauth_return cookies set, got %v", cookies)
	}

	// 3. GET /api/auth/google/callback error param
	reqErr := httptest.NewRequest("GET", "/api/auth/google/callback?error=access_denied", nil)
	recErr := httptest.NewRecorder()
	mux.ServeHTTP(recErr, reqErr)
	if recErr.Code != http.StatusFound || !strings.Contains(recErr.Header().Get("Location"), "error=google_denied") {
		t.Errorf("expected redirect with error=google_denied, got code %d loc %s", recErr.Code, recErr.Header().Get("Location"))
	}

	// 4. GET /api/auth/google/callback invalid state
	reqBadState := httptest.NewRequest("GET", "/api/auth/google/callback?code=mock-code&state=bad-state", nil)
	reqBadState.AddCookie(&http.Cookie{Name: "oauth_state", Value: "different-state"})
	recBadState := httptest.NewRecorder()
	mux.ServeHTTP(recBadState, reqBadState)
	if recBadState.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid oauth state, got %d", recBadState.Code)
	}

	// 5. GET /api/auth/google/callback missing code
	reqMissingCode := httptest.NewRequest("GET", "/api/auth/google/callback?state="+stateVal, nil)
	reqMissingCode.AddCookie(&http.Cookie{Name: "oauth_state", Value: stateVal})
	recMissingCode := httptest.NewRecorder()
	mux.ServeHTTP(recMissingCode, reqMissingCode)
	if recMissingCode.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing code, got %d", recMissingCode.Code)
	}

	// 6. POST /api/auth/google (One-Tap)
	reqOneTapGet := httptest.NewRequest("GET", "/api/auth/google", nil)
	recOneTapGet := httptest.NewRecorder()
	mux.ServeHTTP(recOneTapGet, reqOneTapGet)
	if recOneTapGet.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405 for GET /api/auth/google, got %d", recOneTapGet.Code)
	}

	reqOneTapEmpty := httptest.NewRequest("POST", "/api/auth/google", bytes.NewReader([]byte(`{"id_token":""}`)))
	recOneTapEmpty := httptest.NewRecorder()
	mux.ServeHTTP(recOneTapEmpty, reqOneTapEmpty)
	if recOneTapEmpty.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for empty id_token, got %d", recOneTapEmpty.Code)
	}

	reqOneTapInvalid := httptest.NewRequest("POST", "/api/auth/google", bytes.NewReader([]byte(`{"id_token":"invalid.token.here"}`)))
	recOneTapInvalid := httptest.NewRecorder()
	mux.ServeHTTP(recOneTapInvalid, reqOneTapInvalid)
	if recOneTapInvalid.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for invalid google token, got %d", recOneTapInvalid.Code)
	}
}
