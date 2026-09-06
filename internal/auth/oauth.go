package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// callbackBase infers the public base URL for OAuth callbacks:
// OAUTH_REDIRECT_BASE wins, else the request host (http on localhost).
func callbackBase(r *http.Request) string {
	if v := strings.TrimSpace(os.Getenv("OAUTH_REDIRECT_BASE")); v != "" {
		return strings.TrimRight(v, "/")
	}
	scheme := "https"
	host := ""
	if r != nil {
		if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
			scheme = proto
		} else if strings.HasPrefix(r.Host, "localhost") {
			scheme = "http"
		}
		host = r.Host
	}
	if host == "" {
		return ""
	}
	return scheme + "://" + host
}

// CallbackURL is the provider callback address sent to the OAuth server.
func CallbackURL(r *http.Request, provider string) string {
	if base := callbackBase(r); base != "" {
		return base + "/api/auth/" + provider + "/callback"
	}
	return ""
}

func oauthGetJSON(ctx context.Context, url, bearer string, out interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}
	req.Header.Set("Accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("oauth api %d", res.StatusCode)
	}
	return json.NewDecoder(res.Body).Decode(out)
}

func oauthPostForm(ctx context.Context, url string, form url.Values, out interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("oauth token %d: %.120s", res.StatusCode, string(body))
	}
	return json.Unmarshal(body, out)
}

// ProviderConfigured reports which OAuth providers have credentials set.
// Used by /api/config so the frontend renders exactly those buttons.
func ProviderConfigured(provider string) bool {
	switch provider {
	case "google":
		return googleClientID() != ""
	case "github":
		return strings.TrimSpace(os.Getenv("GITHUB_CLIENT_ID")) != ""
	case "facebook":
		return strings.TrimSpace(os.Getenv("FACEBOOK_APP_ID")) != "" ||
			strings.TrimSpace(os.Getenv("FACEBOOK_CLIENT_ID")) != ""
	case "microsoft":
		return strings.TrimSpace(os.Getenv("MICROSOFT_CLIENT_ID")) != ""
	case "apple":
		return strings.TrimSpace(os.Getenv("APPLE_CLIENT_ID")) != ""
	}
	return false
}

// ConfiguredProviders lists providers in display order.
func ConfiguredProviders() []string {
	var out []string
	for _, p := range []string{"google", "github", "facebook", "microsoft", "apple"} {
		if ProviderConfigured(p) {
			out = append(out, p)
		}
	}
	return out
}
