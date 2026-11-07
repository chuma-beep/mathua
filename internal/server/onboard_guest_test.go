package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGuestDiagnosticUsesBrowserGuestIDForSignupAdoption(t *testing.T) {
	s, mux, store, _, _ := setupTestServerWithPlanner(t)
	guestID := "guest_onboard_claim_test"
	body := fmt.Sprintf(`{"name":"Ada","student_id":%q,"concept_ids":["a"]}`, guestID)
	req := httptest.NewRequest(http.MethodPost, "/api/goal/diagnostic", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("diagnostic start = %d: %s", rr.Code, rr.Body.String())
	}
	var started struct {
		StudentID string `json:"student_id"`
		SessionID string `json:"session_id"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &started); err != nil {
		t.Fatalf("decode start: %v", err)
	}
	if started.StudentID != guestID || started.SessionID == "" {
		t.Fatalf("started = %+v, want guest id %q", started, guestID)
	}
	if student, err := store.GetStudent(guestID); err != nil || student == nil {
		t.Fatalf("guest row = %+v, %v", student, err)
	}
	_, st, err := s.auth.SignupWithGuest("Ada", "ada_onboard_claim", "Engine!n1", guestID)
	if err != nil {
		t.Fatalf("signup with guest: %v", err)
	}
	if st.ID != guestID {
		t.Fatalf("adopted id = %q, want %q", st.ID, guestID)
	}
}
