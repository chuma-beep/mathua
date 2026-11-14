package server

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// turnstileSiteverifyURL is a variable so tests can point at a local server.
var turnstileSiteverifyURL = "https://challenges.cloudflare.com/turnstile/v0/siteverify"

// verifyTurnstile checks a Cloudflare Turnstile token when TURNSTILE_SECRET is
// configured. An empty secret disables the check, which keeps local
// development and self-hosted deployments without a widget working. A
// configured secret with an empty or invalid token fails closed.
func (s *Server) verifyTurnstile(token string) bool {
	secret := strings.TrimSpace(s.turnstileSecret)
	if secret == "" {
		return true
	}
	token = strings.TrimSpace(token)
	if token == "" {
		return false
	}

	form := url.Values{}
	form.Set("secret", secret)
	form.Set("response", token)
	req, err := http.NewRequest(http.MethodPost, s.turnstileURL, strings.NewReader(form.Encode()))
	if err != nil {
		return false
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := s.turnstileClient
	if client == nil {
		client = &http.Client{Timeout: 5 * time.Second}
	}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return false
	}

	var out struct {
		Success bool `json:"success"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&out); err != nil {
		return false
	}
	return out.Success
}
