package server

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/chuma-beep/mathua/internal/fire"
	"github.com/chuma-beep/mathua/internal/storage"
)

type fireActivationRequest struct {
	Ack       string `json:"ack"`
	Cohort    string `json:"cohort"`
	Allowlist string `json:"allowlist"`
}

func (s *Server) handleAdminFire(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	act, err := s.eng.FIReActivation()
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	cfg := s.eng.FIReConfig()
	writeJSON(w, map[string]any{
		"mode":           cfg.Mode,
		"cohort":         cfg.Cohort,
		"allowlist_size": len(cfg.Allow),
		"degraded":       cfg.Degraded,
		"description":    cfg.Describe(),
		"activation":     act,
		"metrics":        s.eng.FIReMetrics(),
		"server_now":     time.Now().UTC().Format(time.RFC3339),
	})
}

func (s *Server) handleAdminFireReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	report, err := s.eng.FIReDryRun(time.Now().UTC())
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, report)
}

func (s *Server) handleAdminFireActivate(w http.ResponseWriter, r *http.Request) {
	if !s.ownerOnly(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req fireActivationRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	cfg := fire.ParseConfig("active", req.Ack, req.Cohort, req.Allowlist)
	if cfg.Mode != fire.Active {
		writeError(w, "activation rejected: "+cfg.Degraded, http.StatusBadRequest)
		return
	}
	now := time.Now().UTC()
	act, err := s.eng.SetFIReActivation(storage.FireActivation{
		Mode:           string(fire.Active),
		Cohort:         cfg.Cohort,
		AcknowledgedBy: actorFrom(r).ID,
		AcknowledgedAt: now,
		UpdatedAt:      now,
	})
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.auditFire(r, storage.ActionFireActivation, act)
	writeJSON(w, map[string]any{
		"activation": act,
		"note":       "activation is durable; apply mode=active and the same cohort and allowlist to the process environment, then restart",
	})
}

func (s *Server) handleAdminFireDeactivate(w http.ResponseWriter, r *http.Request) {
	if !s.ownerOnly(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	act, err := s.eng.SetFIReActivation(storage.FireActivation{
		Mode:           string(fire.Shadow),
		AcknowledgedBy: actorFrom(r).ID,
		UpdatedAt:      time.Now().UTC(),
	})
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.auditFire(r, storage.ActionFireDeactivation, act)
	writeJSON(w, map[string]any{
		"activation": act,
		"note":       "durable activation disables new applications immediately; restart with mode=shadow or off to downgrade the running process. Applied rows are not reversed",
	})
}

func (s *Server) ownerOnly(w http.ResponseWriter, r *http.Request) bool {
	if actorFrom(r).Role != storage.RoleOwner {
		writeError(w, "owner permission required", http.StatusForbidden)
		return false
	}
	return true
}

func (s *Server) auditFire(r *http.Request, action string, act storage.FireActivation) {
	after, _ := json.Marshal(map[string]any{
		"mode":     act.Mode,
		"cohort":   act.Cohort,
		"revision": act.Revision,
	})
	_, _ = s.repo.RecordAdminAudit(storage.AdminAuditRecord{
		ActorID:    actorFrom(r).ID,
		Action:     action,
		EntityType: "fire_activation",
		EntityID:   strings.TrimSpace(act.Cohort),
		AfterJSON:  string(after),
		CreatedAt:  time.Now().UTC().Format(time.RFC3339),
	})
}
