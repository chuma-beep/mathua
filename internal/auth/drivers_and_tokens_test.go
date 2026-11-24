package auth

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestValidateEmailAcceptsSubdomains(t *testing.T) {
	cases := []struct {
		email string
		valid bool
	}{
		{"", true},
		{"   ", true},
		{"alice@example.com", true},
		{"user.name+tag@domain.co.uk", true},
		{"not-an-email", false},
		{"@domain.com", false},
		{"user@", false},
		{"user name@domain.com", false},
	}

	for _, c := range cases {
		err := ValidateEmail(c.email)
		if c.valid && err != nil {
			t.Errorf("expected valid for %q, got error: %v", c.email, err)
		}
		if !c.valid && err == nil {
			t.Errorf("expected error for %q, got nil", c.email)
		}
	}
}

func TestGuestToken_ValidationAndExpiry(t *testing.T) {
	svc := newTestService(t)

	// Empty student ID rejected
	if _, err := IssueGuestToken(""); err == nil {
		t.Error("expected error for empty student ID, got nil")
	}
	if _, err := IssueGuestToken("   "); err == nil {
		t.Error("expected error for whitespace student ID, got nil")
	}

	// Valid guest student ID mints valid token
	tok, err := IssueGuestToken("guest_abcdef123456")
	if err != nil {
		t.Fatalf("issue guest token: %v", err)
	}
	if tok == "" {
		t.Fatal("empty token returned")
	}

	// Validate token resolves student ID
	sid, err := svc.ValidateToken(tok)
	if err != nil {
		t.Fatalf("validate token: %v", err)
	}
	if sid != "guest_abcdef123456" {
		t.Errorf("expected sid guest_abcdef123456, got %q", sid)
	}

	// Invalid tokens rejected
	if _, err := svc.ValidateToken(""); err == nil {
		t.Error("expected error for empty token string")
	}
	if _, err := svc.ValidateToken("invalid.jwt.token"); err == nil {
		t.Error("expected error for malformed token string")
	}
}

func TestCheckPasswordRejectsAnEmptyPassword(t *testing.T) {
	svc := newTestService(t)

	// Create user with password
	_, st, err := svc.Signup("PassUser", "passuser", "CorrectHorseBattery9!")
	if err != nil {
		t.Fatalf("signup: %v", err)
	}

	// Valid password check
	if err := svc.CheckPassword(st.ID, "CorrectHorseBattery9!"); err != nil {
		t.Errorf("expected correct password to succeed, got %v", err)
	}

	// Incorrect password check
	if err := svc.CheckPassword(st.ID, "WrongPassword!"); err != ErrPasswordMismatch {
		t.Errorf("expected ErrPasswordMismatch, got %v", err)
	}

	// Empty password check
	if err := svc.CheckPassword(st.ID, ""); err != ErrPasswordMismatch {
		t.Errorf("expected ErrPasswordMismatch for empty pass, got %v", err)
	}

	// Passwordless user (OAuth only)
	oauthUser, err := svc.repo.CreateUser("OAuthUser", "oauthuser", "")
	if err != nil {
		t.Fatalf("create oauth user: %v", err)
	}
	if err := svc.CheckPassword(oauthUser.ID, "AnyPassword!"); err == nil {
		t.Error("expected error checking password on passwordless account, got nil")
	}

	// Nonexistent student
	if err := svc.CheckPassword("nonexistent_id", "pass"); err == nil {
		t.Error("expected error for nonexistent student, got nil")
	}
}

func TestOAuth_CallbackURL(t *testing.T) {
	// Nil request
	if url := CallbackURL(nil, "github"); url != "" {
		t.Errorf("expected empty string for nil request, got %q", url)
	}

	// Localhost request
	reqLocal := httptest.NewRequest("GET", "http://localhost:3000/api/auth/login", nil)
	if url := CallbackURL(reqLocal, "github"); url != "http://localhost:3000/api/auth/github/callback" {
		t.Errorf("expected localhost callback URL, got %q", url)
	}

	// Forwarded proto
	reqRemote := httptest.NewRequest("GET", "http://mathua.org/api/auth/login", nil)
	reqRemote.Header.Set("X-Forwarded-Proto", "https")
	if url := CallbackURL(reqRemote, "microsoft"); url != "https://mathua.org/api/auth/microsoft/callback" {
		t.Errorf("expected https forwarded callback URL, got %q", url)
	}

	// OAUTH_REDIRECT_BASE override
	t.Setenv("OAUTH_REDIRECT_BASE", "https://auth.mathua.com/")
	if url := CallbackURL(reqLocal, "google"); url != "https://auth.mathua.com/api/auth/google/callback" {
		t.Errorf("expected OAUTH_REDIRECT_BASE callback URL, got %q", url)
	}
}

func TestOAuth_ConfiguredProviders(t *testing.T) {
	// Clean all credentials
	for _, k := range []string{
		"GOOGLE_CLIENT_ID", "GITHUB_CLIENT_ID", "GITHUB_CLIENT_SECRET",
		"FACEBOOK_APP_ID", "FACEBOOK_CLIENT_ID", "FACEBOOK_APP_SECRET",
		"FACEBOOK_CLIENT_SECRET", "MICROSOFT_CLIENT_ID", "MICROSOFT_CLIENT_SECRET",
		"APPLE_CLIENT_ID", "APPLE_TEAM_ID", "APPLE_KEY_ID", "APPLE_PRIVATE_KEY",
		"DISCORD_CLIENT_ID", "DISCORD_CLIENT_SECRET",
	} {
		t.Setenv(k, "")
	}

	if len(ConfiguredProviders()) != 0 {
		t.Errorf("expected 0 configured providers, got %v", ConfiguredProviders())
	}
	if ProviderConfigured("github") || ProviderConfigured("facebook") || ProviderConfigured("unknown") {
		t.Error("expected unconfigured providers to report false")
	}

	// An id without its secret is not configured: the config endpoint drives
	// the buttons, so a half-set provider must stay hidden rather than fail
	// after the learner has been sent to the provider.
	t.Setenv("GITHUB_CLIENT_ID", "gh-client-123")
	if ProviderConfigured("github") {
		t.Error("github with only a client id must not count as configured")
	}
	t.Setenv("GITHUB_CLIENT_SECRET", "gh-secret-123")
	if !ProviderConfigured("github") {
		t.Error("expected github configured")
	}

	// Apple is excluded from the product even when every credential is set.
	t.Setenv("APPLE_CLIENT_ID", "apple-client-456")
	t.Setenv("APPLE_TEAM_ID", "team-123")
	t.Setenv("APPLE_KEY_ID", "key-123")
	t.Setenv("APPLE_PRIVATE_KEY", "-----BEGIN PRIVATE KEY-----")
	if ProviderConfigured("apple") {
		t.Error("apple must stay unadvertised")
	}

	if ProviderConfigured("microsoft") {
		t.Error("expected microsoft unconfigured")
	}

	configured := ConfiguredProviders()
	if len(configured) != 1 || configured[0] != "github" {
		t.Errorf("expected [github], got %v", configured)
	}

	// Enable Facebook via FACEBOOK_APP_ID plus its secret.
	t.Setenv("FACEBOOK_APP_ID", "fb-app-123")
	if ProviderConfigured("facebook") {
		t.Error("facebook with only an app id must not count as configured")
	}
	t.Setenv("FACEBOOK_APP_SECRET", "fb-secret-123")
	if !ProviderConfigured("facebook") {
		t.Error("expected facebook configured via FACEBOOK_APP_ID")
	}

	// Enable Microsoft
	t.Setenv("MICROSOFT_CLIENT_ID", "ms-client-123")
	if ProviderConfigured("microsoft") {
		t.Error("microsoft with only a client id must not count as configured")
	}
	t.Setenv("MICROSOFT_CLIENT_SECRET", "ms-secret-123")
	if !ProviderConfigured("microsoft") {
		t.Error("expected microsoft configured")
	}

	// Discord needs both halves for the token exchange.
	t.Setenv("DISCORD_CLIENT_ID", "discord-client-123")
	if ProviderConfigured("discord") {
		t.Error("discord with only a client id must not count as configured")
	}
	t.Setenv("DISCORD_CLIENT_SECRET", "discord-secret-123")
	if !ProviderConfigured("discord") {
		t.Error("expected discord configured")
	}
}

func TestOAuth_DriversAndLoginURLs(t *testing.T) {
	drivers := Drivers()
	expectedDrivers := []string{"github", "facebook", "microsoft", "apple", "discord"}

	for _, name := range expectedDrivers {
		drv, ok := drivers[name]
		if !ok {
			t.Errorf("missing driver %q in Drivers()", name)
			continue
		}
		if drv.LoginURL == nil {
			t.Errorf("driver %q has nil LoginURL", name)
			continue
		}
		if drv.ClientID == nil {
			t.Errorf("driver %q has nil ClientID func", name)
		}
		if drv.Exchange == nil {
			t.Errorf("driver %q has nil Exchange func", name)
		}

		loginURL := drv.LoginURL("test-client-id", "https://example.com/callback", "state123")
		parsed, err := url.Parse(loginURL)
		if err != nil {
			t.Fatalf("driver %q LoginURL %q not valid URL: %v", name, loginURL, err)
		}

		q := parsed.Query()
		if q.Get("client_id") != "test-client-id" {
			t.Errorf("driver %q expected client_id test-client-id, got %q", name, q.Get("client_id"))
		}
		if q.Get("redirect_uri") != "https://example.com/callback" {
			t.Errorf("driver %q expected redirect_uri, got %q", name, q.Get("redirect_uri"))
		}
		if q.Get("state") != "state123" {
			t.Errorf("driver %q expected state state123, got %q", name, q.Get("state"))
		}
	}
}

func TestApple_IsRelayAddress(t *testing.T) {
	cases := []struct {
		email   string
		isRelay bool
	}{
		{"abc@privaterelay.appleid.com", true},
		{"ABC@PRIVATERELAY.APPLEID.COM", true},
		{"  user.hide@privaterelay.appleid.com  ", true},
		{"user@gmail.com", false},
		{"alice@apple.com", false},
		{"privaterelay.appleid.com", false},
		{"", false},
	}

	for _, c := range cases {
		got := isRelayAddress(c.email)
		if got != c.isRelay {
			t.Errorf("isRelayAddress(%q) = %v, want %v", c.email, got, c.isRelay)
		}
	}
}

func TestOAuth_HTTPHelpers(t *testing.T) {
	ctx := context.Background()

	// Server handling GET with bearer auth and POST with form body
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if r.URL.Path == "/error" {
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}
			if r.Header.Get("Authorization") != "Bearer secret-token" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, `{"status":"ok","value":42}`)

		case http.MethodPost:
			if r.URL.Path == "/fail" {
				http.Error(w, "server error", http.StatusInternalServerError)
				return
			}
			if err := r.ParseForm(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if r.FormValue("code") != "auth-code-123" {
				http.Error(w, "missing code", http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, `{"access_token":"tok-abc"}`)
		}
	}))
	defer ts.Close()

	// Test oauthGetJSON success
	var getResult struct {
		Status string `json:"status"`
		Value  int    `json:"value"`
	}
	if err := oauthGetJSON(ctx, ts.URL+"/get", "secret-token", &getResult); err != nil {
		t.Fatalf("oauthGetJSON success case failed: %v", err)
	}
	if getResult.Status != "ok" || getResult.Value != 42 {
		t.Errorf("unexpected getResult: %+v", getResult)
	}

	// Test oauthGetJSON non-200 error
	if err := oauthGetJSON(ctx, ts.URL+"/error", "", &getResult); err == nil {
		t.Fatal("expected error on non-200 GET, got nil")
	}

	// Test oauthPostForm success
	form := url.Values{"code": {"auth-code-123"}}
	var postResult struct {
		AccessToken string `json:"access_token"`
	}
	if err := oauthPostForm(ctx, ts.URL+"/post", form, &postResult); err != nil {
		t.Fatalf("oauthPostForm success case failed: %v", err)
	}
	if postResult.AccessToken != "tok-abc" {
		t.Errorf("unexpected postResult token: %q", postResult.AccessToken)
	}

	// Test oauthPostForm non-200 error
	if err := oauthPostForm(ctx, ts.URL+"/fail", form, &postResult); err == nil {
		t.Fatal("expected error on non-200 POST, got nil")
	}
}

func TestGoogle_RedirectURL_Variations(t *testing.T) {
	// Explicit GOOGLE_REDIRECT_URL env override
	t.Setenv("GOOGLE_REDIRECT_URL", "https://override.example.com/callback")
	if got := googleRedirectURL(nil); got != "https://override.example.com/callback" {
		t.Errorf("expected override URL, got %q", got)
	}
	t.Setenv("GOOGLE_REDIRECT_URL", "")

	// Inferred from request headers
	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/login", nil)
	req.Host = "localhost:8080"
	if got := googleRedirectURL(req); got != "http://localhost:8080/api/auth/google/callback" {
		t.Errorf("expected localhost http URL, got %q", got)
	}

	// Forwarded proto
	reqForwarded := httptest.NewRequest(http.MethodGet, "/login", nil)
	reqForwarded.Host = "mathua.org"
	reqForwarded.Header.Set("X-Forwarded-Proto", "https")
	if got := googleRedirectURL(reqForwarded); got != "https://mathua.org/api/auth/google/callback" {
		t.Errorf("expected forwarded https URL, got %q", got)
	}

	// Fallback to GOOGLE_REDIRECT_HOST
	t.Setenv("GOOGLE_REDIRECT_HOST", "custom-host.org")
	if got := googleRedirectURL(nil); got != "https://custom-host.org/api/auth/google/callback" {
		t.Errorf("expected host fallback URL, got %q", got)
	}
	t.Setenv("GOOGLE_REDIRECT_HOST", "")

	// Empty when no host available
	if got := googleRedirectURL(nil); got != "" {
		t.Errorf("expected empty string when unconfigured and request nil, got %q", got)
	}
}

func TestExchange_UnconfiguredProviders(t *testing.T) {
	ctx := context.Background()

	// Clear provider environment variables for test isolation
	t.Setenv("GITHUB_CLIENT_ID", "")
	t.Setenv("GITHUB_OAUTH_CLIENT_ID", "")
	t.Setenv("GITHUB_CLIENT_SECRET", "")
	t.Setenv("GITHUB_OAUTH_CLIENT_SECRET", "")

	t.Setenv("FACEBOOK_APP_ID", "")
	t.Setenv("FACEBOOK_CLIENT_ID", "")
	t.Setenv("FACEBOOK_APP_SECRET", "")
	t.Setenv("FACEBOOK_CLIENT_SECRET", "")

	t.Setenv("MICROSOFT_CLIENT_ID", "")
	t.Setenv("MICROSOFT_CLIENT_SECRET", "")

	t.Setenv("APPLE_CLIENT_ID", "")
	t.Setenv("APPLE_SERVICES_ID", "")

	t.Setenv("GOOGLE_CLIENT_ID", "")
	t.Setenv("GOOGLE_OAUTH_CLIENT_ID", "")
	t.Setenv("GOOGLE_CLIENT_SECRET", "")
	t.Setenv("GOOGLE_OAUTH_CLIENT_SECRET", "")

	if _, err := ExchangeGithubCode(ctx, "c", "r"); err == nil || !strings.Contains(err.Error(), "not configured") {
		t.Errorf("expected unconfigured error from ExchangeGithubCode, got %v", err)
	}

	if _, err := ExchangeFacebookCode(ctx, "c", "r"); err == nil || !strings.Contains(err.Error(), "not configured") {
		t.Errorf("expected unconfigured error from ExchangeFacebookCode, got %v", err)
	}

	if _, err := ExchangeMicrosoftCode(ctx, "c", "r"); err == nil || !strings.Contains(err.Error(), "not configured") {
		t.Errorf("expected unconfigured error from ExchangeMicrosoftCode, got %v", err)
	}

	if _, err := ExchangeAppleCode(ctx, "c", "r", ""); err == nil || !strings.Contains(err.Error(), "not configured") {
		t.Errorf("expected unconfigured error from ExchangeAppleCode, got %v", err)
	}

	if _, err := VerifyGoogleIDToken(ctx, "fake-tok"); err == nil || !strings.Contains(err.Error(), "not configured") {
		t.Errorf("expected unconfigured error from VerifyGoogleIDToken, got %v", err)
	}

	if _, err := ExchangeGoogleCode(ctx, "c", nil); err == nil || !strings.Contains(err.Error(), "not configured") {
		t.Errorf("expected unconfigured error from ExchangeGoogleCode, got %v", err)
	}
}

func TestSecretMinLength(t *testing.T) {
	if minLen := SecretMinLength(); minLen != 32 {
		t.Errorf("expected SecretMinLength 32, got %d", minLen)
	}
}

func TestReset_SMTPConfiguredAndFrontendBase(t *testing.T) {
	// SMTPConfigured checks SMTP_HOST
	t.Setenv("SMTP_HOST", "")
	if SMTPConfigured() {
		t.Error("expected SMTPConfigured false when empty")
	}
	t.Setenv("SMTP_HOST", "smtp.example.com")
	if !SMTPConfigured() {
		t.Error("expected SMTPConfigured true when populated")
	}

	// frontendBase reads FRONTEND_URL
	t.Setenv("FRONTEND_URL", "  https://frontend.example.com  ")
	if got := frontendBase(); got != "https://frontend.example.com" {
		t.Errorf("expected trimmed FRONTEND_URL, got %q", got)
	}

	t.Setenv("FRONTEND_URL", "")
	if got := frontendBase(); got != "" {
		t.Errorf("expected empty string when unconfigured, got %q", got)
	}
}
