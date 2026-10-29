package server

import (
	"net/http"
	"testing"

	"github.com/chuma-beep/mathua/internal/fire"
	"github.com/chuma-beep/mathua/internal/storage"
)

func TestAdminFireActivationIsOwnerOnlyAndFailClosed(t *testing.T) {
	e := newAdminEnv(t)
	if rec := e.do(t, http.MethodGet, "/api/admin/fire", e.adminTok, ""); rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body %s", rec.Code, rec.Body.String())
	}
	if rec := e.do(t, http.MethodGet, "/api/admin/fire/report", e.adminTok, ""); rec.Code != http.StatusOK {
		t.Fatalf("report = %d, body %s", rec.Code, rec.Body.String())
	}
	if rec := e.do(t, http.MethodGet, "/api/admin/fire", e.studentTok, ""); rec.Code != http.StatusForbidden {
		t.Fatalf("learner status = %d, want 403", rec.Code)
	}
	bad := `{"ack":"wrong","cohort":"pilot","allowlist":"x"}`
	if rec := e.do(t, http.MethodPost, "/api/admin/fire/activate", e.adminTok, bad); rec.Code != http.StatusBadRequest {
		t.Fatalf("bad activation = %d, body %s", rec.Code, rec.Body.String())
	}
	good := `{"ack":"` + fire.ActivationToken + `","cohort":"pilot","allowlist":"x"}`
	if rec := e.do(t, http.MethodPost, "/api/admin/fire/activate", e.adminTok, good); rec.Code != http.StatusOK {
		t.Fatalf("activation = %d, body %s", rec.Code, rec.Body.String())
	}
	act, err := e.store.(interface {
		GetFireActivation() (*storage.FireActivation, error)
	}).GetFireActivation()
	if err != nil {
		t.Fatalf("read activation: %v", err)
	}
	if act == nil || act.Mode != string(fire.Active) || act.Cohort != "pilot" {
		t.Fatalf("activation = %+v", act)
	}
	if rec := e.do(t, http.MethodPost, "/api/admin/fire/deactivate", e.adminTok, ""); rec.Code != http.StatusOK {
		t.Fatalf("deactivate = %d, body %s", rec.Code, rec.Body.String())
	}
}
