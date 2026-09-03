package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"google.golang.org/api/idtoken"
)

// GoogleProfile is the normalized identity from Google (one-tap or code flow).
type GoogleProfile struct {
	GoogleID string
	Email    string
	Name     string
	Picture  string
}

// googleClientID resolves GOOGLE_CLIENT_ID (also supports GOOGLE_OAUTH_CLIENT_ID).
func googleClientID() string {
	if v := os.Getenv("GOOGLE_CLIENT_ID"); v != "" {
		return strings.TrimSpace(v)
	}
	return strings.TrimSpace(os.Getenv("GOOGLE_OAUTH_CLIENT_ID"))
}

func googleClientSecret() string {
	if v := os.Getenv("GOOGLE_CLIENT_SECRET"); v != "" {
		return strings.TrimSpace(v)
	}
	return strings.TrimSpace(os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"))
}

func googleRedirectURL(r *http.Request) string {
	if v := os.Getenv("GOOGLE_REDIRECT_URL"); v != "" {
		return strings.TrimSpace(v)
	}
	// Default to current origin + callback path
	scheme := "https"
	if r != nil {
		if r.Header.Get("X-Forwarded-Proto") == "http" || r.TLS == nil && r.Header.Get("X-Forwarded-Proto") == "" && strings.HasPrefix(r.Host, "localhost") {
			scheme = "http"
		}
		if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
			scheme = proto
		}
		if r.URL.Scheme != "" {
			scheme = r.URL.Scheme
		}
	}
	host := ""
	if r != nil && r.Host != "" {
		host = r.Host
	}
	if host == "" {
		host = os.Getenv("GOOGLE_REDIRECT_HOST")
	}
	if host == "" {
		return ""
	}
	return fmt.Sprintf("%s://%s/api/auth/google/callback", scheme, host)
}

// VerifyGoogleIDToken validates a Google One-Tap id_token against GOOGLE_CLIENT_ID and returns the profile.
// Uses google.golang.org/api/idtoken for signature + aud/iss/exp checks.
func VerifyGoogleIDToken(ctx context.Context, idToken string) (*GoogleProfile, error) {
	aud := googleClientID()
	if aud == "" {
		return nil, fmt.Errorf("GOOGLE_CLIENT_ID not configured")
	}
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	payload, err := idtoken.Validate(ctx, idToken, aud)
	if err != nil {
		return nil, fmt.Errorf("id_token verify: %w", err)
	}
	// payload.Claims are map-like: sub, email, name, picture
	claims := payload.Claims
	gid, _ := claims["sub"].(string)
	email, _ := claims["email"].(string)
	name, _ := claims["name"].(string)
	picture, _ := claims["picture"].(string)
	if gid == "" {
		return nil, fmt.Errorf("google id_token missing sub")
	}
	if name == "" {
		name = email
	}
	return &GoogleProfile{GoogleID: gid, Email: email, Name: name, Picture: picture}, nil
}

// Exchange Google OAuth code for profile (full redirect flow).
func ExchangeGoogleCode(ctx context.Context, code string, r *http.Request) (*GoogleProfile, error) {
	clientID := googleClientID()
	secret := googleClientSecret()
	if clientID == "" || secret == "" {
		return nil, fmt.Errorf("GOOGLE_CLIENT_ID/SECRET not configured")
	}
	redirectURL := googleRedirectURL(r)
	if redirectURL == "" {
		return nil, fmt.Errorf("GOOGLE_REDIRECT_URL not configured and cannot infer from request")
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Exchange code for access_token
	form := url.Values{}
	form.Set("code", code)
	form.Set("client_id", clientID)
	form.Set("client_secret", secret)
	form.Set("redirect_uri", redirectURL)
	form.Set("grant_type", "authorization_code")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://oauth2.googleapis.com/token", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google token exchange %d", res.StatusCode)
	}
	var tokRes struct {
		AccessToken string `json:"access_token"`
		IDToken     string `json:"id_token"`
	}
	if err := json.NewDecoder(res.Body).Decode(&tokRes); err != nil {
		return nil, err
	}
	// Prefer id_token if present (contains profile), else fallback to userinfo
	if tokRes.IDToken != "" {
		payload, err := idtoken.Validate(ctx, tokRes.IDToken, clientID)
		if err == nil {
			claims := payload.Claims
			gid, _ := claims["sub"].(string)
			email, _ := claims["email"].(string)
			name, _ := claims["name"].(string)
			picture, _ := claims["picture"].(string)
			if gid != "" {
				if name == "" {
					name = email
				}
				return &GoogleProfile{GoogleID: gid, Email: email, Name: name, Picture: picture}, nil
			}
		}
	}
	if tokRes.AccessToken == "" {
		return nil, fmt.Errorf("no access_token from google")
	}
	// Fallback: userinfo
	req2, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	req2.Header.Set("Authorization", "Bearer "+tokRes.AccessToken)
	res2, err := http.DefaultClient.Do(req2)
	if err != nil {
		return nil, err
	}
	defer res2.Body.Close()
	if res2.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google userinfo %d", res2.StatusCode)
	}
	var ui struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := json.NewDecoder(res2.Body).Decode(&ui); err != nil {
		return nil, err
	}
	if ui.ID == "" {
		return nil, fmt.Errorf("google userinfo missing id")
	}
	if ui.Name == "" {
		ui.Name = ui.Email
	}
	return &GoogleProfile{GoogleID: ui.ID, Email: ui.Email, Name: ui.Name, Picture: ui.Picture}, nil
}
