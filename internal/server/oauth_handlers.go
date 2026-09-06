package server

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/chuma-beep/mathua/internal/auth"
)

func oauthProviderFromLoginPath(p string) (string, bool) {
	rest, ok := strings.CutPrefix(p, "/api/auth/")
	if !ok {
		return "", false
	}
	name, ok := strings.CutSuffix(rest, "/login")
	if !ok {
		return "", false
	}
	if _, known := auth.Drivers()[name]; !known {
		return "", false
	}
	return name, true
}

func oauthProviderFromCallbackPath(p string) (string, bool) {
	rest, ok := strings.CutPrefix(p, "/api/auth/")
	if !ok {
		return "", false
	}
	name, ok := strings.CutSuffix(rest, "/callback")
	if !ok {
		return "", false
	}
	if _, known := auth.Drivers()[name]; !known {
		return "", false
	}
	return name, true
}

func oauthCookie(w http.ResponseWriter, r *http.Request, name, value string, clear bool) {
	c := &http.Cookie{Name: name, Value: value, Path: "/", HttpOnly: true, SameSite: http.SameSiteLaxMode, MaxAge: 600}
	c.Secure = r.Header.Get("X-Forwarded-Proto") == "https" || r.TLS != nil
	if clear {
		c.Value = ""
		c.MaxAge = -1
	}
	http.SetCookie(w, c)
}

// GET /api/auth/{provider}/login?intent=login|link&link_token=..&return=..
func (s *Server) handleOAuthLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	if s.auth == nil {
		writeError(w, "authentication is disabled", 400)
		return
	}
	provider, ok := oauthProviderFromLoginPath(r.URL.Path)
	if !ok {
		writeError(w, "unknown provider", 404)
		return
	}
	driver := auth.Drivers()[provider]
	clientID := driver.ClientID()
	if clientID == "" {
		writeError(w, provider+" auth not configured", 500)
		return
	}
	redirectURI := auth.CallbackURL(r, provider)
	if redirectURI == "" {
		writeError(w, "cannot determine callback URL", 500)
		return
	}
	state := newUUID()
	oauthCookie(w, r, "oauth_state", state, false)
	intent := r.URL.Query().Get("intent")
	if intent != "link" {
		intent = "login"
	}
	oauthCookie(w, r, "oauth_intent", intent, false)
	ret := r.URL.Query().Get("return")
	if ret == "" {
		ret = "/profile"
	}
	oauthCookieValue(w, r, "oauth_return", ret)
	if intent == "link" {
		oauthCookieValue(w, r, "oauth_link", r.URL.Query().Get("link_token"))
	}
	http.Redirect(w, r, driver.LoginURL(clientID, redirectURI, state), http.StatusFound)
}

func oauthCookieValue(w http.ResponseWriter, r *http.Request, name, value string) {
	oauthCookie(w, r, name, value, false)
}

// handleOAuthCallback serves GET (github/facebook/microsoft) and GET+POST
// (apple form_post) callbacks.
func (s *Server) handleOAuthCallback(w http.ResponseWriter, r *http.Request) {
	if s.auth == nil {
		writeError(w, "authentication is disabled", 400)
		return
	}
	provider, ok := oauthProviderFromCallbackPath(r.URL.Path)
	if !ok {
		writeError(w, "unknown provider", 404)
		return
	}
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			writeError(w, "invalid form", 400)
			return
		}
	} else if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	ret := "/profile"
	if c, err := r.Cookie("oauth_return"); err == nil && c.Value != "" {
		ret = c.Value
	}
	fail := func(code string) {
		u := ret
		sep := "?"
		if strings.Contains(u, "?") {
			sep = "&"
		}
		http.Redirect(w, r, u+sep+"error="+code, http.StatusFound)
	}
	if r.FormValue("error") != "" {
		fail(provider + "_denied")
		return
	}
	state := r.FormValue("state")
	cookie, _ := r.Cookie("oauth_state")
	if state == "" || cookie == nil || cookie.Value != state {
		writeError(w, "invalid oauth state", 400)
		return
	}
	code := r.FormValue("code")
	if strings.TrimSpace(code) == "" {
		writeError(w, "missing code", 400)
		return
	}
	driver := auth.Drivers()[provider]
	redirectURI := auth.CallbackURL(r, provider)
	profile, err := driver.Exchange(r.Context(), code, redirectURI, r)
	if err != nil {
		log.Printf("oauth %s exchange: %v", provider, err)
		fail(provider + "_failed")
		return
	}
	intent := ""
	if c, err := r.Cookie("oauth_intent"); err == nil {
		intent = c.Value
	}
	oauthCookie(w, r, "oauth_state", "", true)
	oauthCookie(w, r, "oauth_intent", "", true)
	oauthCookie(w, r, "oauth_return", "", true)
	oauthCookie(w, r, "oauth_link", "", true)
	if intent == "link" {
		linkTok := ""
		if c, err := r.Cookie("oauth_link"); err == nil {
			linkTok = c.Value
		}
		studentID, ok, err := s.auth.ConsumeLinkToken(linkTok)
		if err != nil || !ok {
			fail("link_expired")
			return
		}
		if err := s.auth.ConnectProvider(studentID, profile); err != nil {
			log.Printf("oauth %s connect: %v", provider, err)
			fail("link_taken")
			return
		}
		u := ret
		sep := "?"
		if strings.Contains(u, "?") {
			sep = "&"
		}
		http.Redirect(w, r, u+sep+"linked="+provider, http.StatusFound)
		return
	}
	token, st, err := s.auth.LoginOrCreateOAuth(profile)
	if err != nil {
		log.Printf("oauth %s login: %v", provider, err)
		fail(provider + "_failed")
		return
	}
	u := frontendRedirect(os.Getenv("FRONTEND_URL"), ret, token, st.Name, st.ID)
	http.Redirect(w, r, u, http.StatusFound)
}

// POST /api/auth/link-token (authed) — binds a connect dance to this student.
func (s *Server) handleLinkToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	if studentID == "" {
		writeError(w, "not authenticated", 401)
		return
	}
	tok, err := s.auth.CreateLinkToken(studentID)
	if err != nil {
		writeError(w, "failed to create link token", 500)
		return
	}
	writeJSON(w, map[string]interface{}{"link_token": tok})
}

// GET /api/auth/identities (authed) — connected providers for Settings UI.
func (s *Server) handleIdentities(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	if studentID == "" {
		writeError(w, "not authenticated", 401)
		return
	}
	ids, err := s.repo.ListIdentities(studentID)
	if err != nil {
		writeError(w, "failed to list identities", 500)
		return
	}
	out := make([]map[string]interface{}, 0, len(ids))
	for _, id := range ids {
		out = append(out, map[string]interface{}{
			"provider": id.Provider,
			"email":    id.Email,
		})
	}
	writeJSON(w, out)
}

// DELETE /api/auth/identities/{provider} (authed).
func (s *Server) handleIdentityDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	if studentID == "" {
		writeError(w, "not authenticated", 401)
		return
	}
	provider := strings.TrimPrefix(r.URL.Path, "/api/auth/identities/")
	if provider == "" || strings.Contains(provider, "/") {
		writeError(w, "unknown provider", 404)
		return
	}
	if _, known := auth.Drivers()[provider]; !known && provider != "google" {
		writeError(w, "unknown provider", 404)
		return
	}
	if err := s.auth.DisconnectProvider(studentID, provider); err != nil {
		writeError(w, err.Error(), 400)
		return
	}
	writeJSON(w, map[string]interface{}{"ok": true})
}

// POST /api/auth/email/request (authed) — send verification link.
func (s *Server) handleEmailRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	if studentID == "" {
		writeError(w, "not authenticated", 401)
		return
	}
	if err := s.auth.RequestEmailVerification(studentID); err != nil {
		writeError(w, err.Error(), 400)
		return
	}
	writeJSON(w, map[string]interface{}{"ok": true})
}

// POST /api/auth/email/verify {token} — burn token, mark verified.
func (s *Server) handleEmailVerify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	var req struct {
		Token string `json:"token"`
	}
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, "invalid request", 400)
		return
	}
	if _, err := s.auth.VerifyEmail(req.Token); err != nil {
		writeError(w, err.Error(), 400)
		return
	}
	writeJSON(w, map[string]interface{}{"ok": true})
}

// POST /api/auth/google/connect {id_token} (authed) — link the Google
// account behind a One-Tap credential to the current student. Redirect
// flows connect via intent=link + link_token instead.
func (s *Server) handleGoogleConnect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	if studentID == "" {
		writeError(w, "not authenticated", 401)
		return
	}
	var req struct {
		IDToken string `json:"id_token"`
	}
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, "invalid request", 400)
		return
	}
	profile, err := s.auth.VerifyGoogleIDToken(r.Context(), req.IDToken)
	if err != nil {
		writeError(w, "invalid google token: "+err.Error(), 401)
		return
	}
	if err := s.auth.ConnectProvider(studentID, auth.OAuthProfile{
		Provider: "google", ProviderID: profile.GoogleID,
		Email: profile.Email, EmailVerified: profile.Email != "",
		Name: profile.Name, AvatarURL: profile.Picture,
	}); err != nil {
		writeError(w, err.Error(), 400)
		return
	}
	writeJSON(w, map[string]interface{}{"ok": true})
}
