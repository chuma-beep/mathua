package server

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// Admin triage auth: a single shared password set by the operator via the
// ADMIN_PASSWORD env var (same trust model as JWT_SECRET — whoever controls
// the deployment environment sets it; there is no in-app bootstrap and no
// super-admin role). Successful login mints a random 24h session token kept
// in server_sessions (durable across restarts/replicas).
// Legacy ADMIN_TOKEN is still accepted by adminAuthorized as a fallback.

const (
	adminSessionTTL        = 24 * time.Hour
	adminPasswordMinLength = 12
	// serverSessionAdmin is the server_sessions kind for triage logins.
	serverSessionAdmin = "admin"
)

func adminPassword() string {
	return strings.TrimSpace(os.Getenv("ADMIN_PASSWORD"))
}

// checkAdminPasswordConfig fail-fasts startup when the operator configured
// a password too weak to be a gate. Unset is fine (triage stays locked).
func checkAdminPasswordConfig() {
	if pw := adminPassword(); pw != "" && len([]rune(pw)) < adminPasswordMinLength {
		log.Fatalf("auth: ADMIN_PASSWORD must be at least %d characters long", adminPasswordMinLength)
	}
}

func newAdminSessionToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// validAdminSession reports whether token is a live admin session.
// Sessions are durable (server_sessions kind "admin"): they survive
// restarts and work across replicas. Expired rows read as missing and are
// deleted opportunistically; the New() janitor sweeps the rest.
func (s *Server) validAdminSession(token string) bool {
	if token == "" {
		return false
	}
	_, exp, found, err := s.repo.GetServerSession(serverSessionAdmin, token)
	if err != nil || !found {
		return false
	}
	if exp == "" {
		return true
	}
	if t, err := time.Parse(time.RFC3339, exp); err != nil || !t.After(time.Now().UTC()) {
		_ = s.repo.DeleteServerSession(serverSessionAdmin, token)
		return false
	}
	return true
}

// POST /api/admin/login {password} → {token, expires_at} or 401.
// 404 when no ADMIN_PASSWORD is configured (triage locked, as before).
func (s *Server) handleAdminLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	expected := adminPassword()
	if expected == "" {
		writeError(w, "not found", 404)
		return
	}
	var req struct {
		Password string `json:"password"`
	}
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, "invalid request", 400)
		return
	}
	if subtle.ConstantTimeCompare([]byte(req.Password), []byte(expected)) != 1 {
		writeError(w, "invalid password", 401)
		return
	}
	token, err := newAdminSessionToken()
	if err != nil {
		writeError(w, "failed to create session", 500)
		return
	}
	exp := time.Now().Add(adminSessionTTL)
	if err := s.repo.UpsertServerSession(serverSessionAdmin, token, "", exp.UTC().Format(time.RFC3339)); err != nil {
		writeError(w, "failed to create session", 500)
		return
	}
	writeJSON(w, map[string]interface{}{
		"token":      token,
		"expires_at": exp.UTC().Format(time.RFC3339),
	})
}

// POST /api/admin/logout — invalidates the session Bearer token.
// Idempotent: always 200, even for unknown/expired tokens.
func (s *Server) handleAdminLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	if token != "" {
		_ = s.repo.DeleteServerSession(serverSessionAdmin, token)
	}
	writeJSON(w, map[string]bool{"ok": true})
}
