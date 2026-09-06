package auth

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func appleClientID() string {
	if v := strings.TrimSpace(os.Getenv("APPLE_CLIENT_ID")); v != "" {
		return v
	}
	return strings.TrimSpace(os.Getenv("APPLE_SERVICES_ID"))
}

func appleTeamID() string { return strings.TrimSpace(os.Getenv("APPLE_TEAM_ID")) }
func appleKeyID() string  { return strings.TrimSpace(os.Getenv("APPLE_KEY_ID")) }

func applePrivateKey() (*ecdsa.PrivateKey, error) {
	raw := strings.ReplaceAll(strings.TrimSpace(os.Getenv("APPLE_PRIVATE_KEY")), "\\n", "\n")
	if raw == "" {
		return nil, fmt.Errorf("APPLE_PRIVATE_KEY not configured")
	}
	block, _ := pem.Decode([]byte(raw))
	if block == nil {
		return nil, fmt.Errorf("APPLE_PRIVATE_KEY is not PEM")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse apple key: %w", err)
	}
	ec, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("apple key is not ECDSA")
	}
	return ec, nil
}

// appleClientSecret builds the ES256-signed client_secret JWT (10 min life).
func appleClientSecret() (string, error) {
	clientID := appleClientID()
	teamID := appleTeamID()
	keyID := appleKeyID()
	if clientID == "" || teamID == "" || keyID == "" {
		return "", fmt.Errorf("APPLE_CLIENT_ID/TEAM_ID/KEY_ID not configured")
	}
	key, err := applePrivateKey()
	if err != nil {
		return "", err
	}
	now := time.Now()
	tok := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"iss": teamID,
		"iat": now.Unix(),
		"exp": now.Add(10 * time.Minute).Unix(),
		"aud": "https://appleid.apple.com",
		"sub": clientID,
	})
	tok.Header["kid"] = keyID
	return tok.SignedString(key)
}

// AppleLoginURL starts Sign in with Apple (the user's name arrives only on
// first authorization, via a form_post `user` field — the callback accepts POST).
func AppleLoginURL(clientID, redirectURI, state string) string {
	q := url.Values{}
	q.Set("client_id", clientID)
	q.Set("redirect_uri", redirectURI)
	q.Set("response_type", "code")
	q.Set("response_mode", "form_post")
	q.Set("scope", "name email")
	q.Set("state", state)
	return "https://appleid.apple.com/auth/authorize?" + q.Encode()
}

type appleJWK struct {
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

var (
	appleJWKSMu sync.Mutex
	appleJWKS   map[string]appleJWK
	appleJWKSAt time.Time
)

func applePublicKey(ctx context.Context, kid string) (crypto.PublicKey, error) {
	appleJWKSMu.Lock()
	defer appleJWKSMu.Unlock()
	if time.Since(appleJWKSAt) > 24*time.Hour || appleJWKS == nil {
		var doc struct {
			Keys []appleJWK `json:"keys"`
		}
		if err := oauthGetJSON(ctx, "https://appleid.apple.com/auth/keys", "", &doc); err != nil {
			return nil, err
		}
		m := make(map[string]appleJWK, len(doc.Keys))
		for _, k := range doc.Keys {
			m[k.Kid] = k
		}
		appleJWKS, appleJWKSAt = m, time.Now()
	}
	k, ok := appleJWKS[kid]
	if !ok {
		return nil, fmt.Errorf("unknown apple key %q", kid)
	}
	nBytes, err := base64.RawURLEncoding.DecodeString(k.N)
	if err != nil {
		return nil, fmt.Errorf("bad apple key N: %w", err)
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(k.E)
	if err != nil {
		return nil, fmt.Errorf("bad apple key E: %w", err)
	}
	e := 0
	for _, b := range eBytes {
		e = e<<8 | int(b)
	}
	return &rsa.PublicKey{N: new(big.Int).SetBytes(nBytes), E: e}, nil
}

// isRelayAddress reports Apple private-relay emails, which must never be
// used as matching keys (user-revocable aliases, not stable identity).
func isRelayAddress(email string) bool {
	parts := strings.Split(strings.ToLower(strings.TrimSpace(email)), "@")
	return len(parts) == 2 && parts[1] == "privaterelay.appleid.com"
}

// ExchangeAppleCode trades the callback code (+ optional first-auth `user`
// JSON) for a normalized profile. The id_token is RS256-verified against
// Apple's JWKS with iss/aud/exp checks.
func ExchangeAppleCode(ctx context.Context, code, redirectURI, userJSON string) (OAuthProfile, error) {
	clientID := appleClientID()
	if clientID == "" {
		return OAuthProfile{}, fmt.Errorf("APPLE_CLIENT_ID not configured")
	}
	secret, err := appleClientSecret()
	if err != nil {
		return OAuthProfile{}, err
	}
	form := url.Values{}
	form.Set("code", code)
	form.Set("client_id", clientID)
	form.Set("client_secret", secret)
	form.Set("redirect_uri", redirectURI)
	form.Set("grant_type", "authorization_code")
	var tok struct {
		IDToken string `json:"id_token"`
	}
	if err := oauthPostForm(ctx, "https://appleid.apple.com/auth/token", form, &tok); err != nil {
		return OAuthProfile{}, err
	}
	if tok.IDToken == "" {
		return OAuthProfile{}, fmt.Errorf("no id_token from apple")
	}
	claims := jwt.MapClaims{}
	parsed, err := jwt.ParseWithClaims(tok.IDToken, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected apple signing method")
		}
		kid, _ := t.Header["kid"].(string)
		return applePublicKey(ctx, kid)
	}, jwt.WithIssuer("https://appleid.apple.com"), jwt.WithAudience(clientID))
	if err != nil || !parsed.Valid {
		return OAuthProfile{}, fmt.Errorf("apple id_token invalid: %w", err)
	}
	sub, _ := claims["sub"].(string)
	if sub == "" {
		return OAuthProfile{}, fmt.Errorf("apple id_token missing sub")
	}
	email, _ := claims["email"].(string)
	emailVerified, _ := claims["email_verified"].(bool)
	if !emailVerified {
		if s, _ := claims["email_verified"].(string); s == "true" {
			emailVerified = true
		}
	}
	name := ""
	if userJSON != "" {
		var u struct {
			Name struct {
				First string `json:"firstName"`
				Last  string `json:"lastName"`
			} `json:"name"`
			Email string `json:"email"`
		}
		if jerr := json.Unmarshal([]byte(userJSON), &u); jerr == nil {
			name = strings.TrimSpace(strings.TrimSpace(u.Name.First) + " " + strings.TrimSpace(u.Name.Last))
			if email == "" {
				email = u.Email
			}
		}
	}
	if name == "" {
		name = email
	}
	if isRelayAddress(email) {
		emailVerified = false // never match on relay aliases
	}
	return OAuthProfile{
		Provider: "apple", ProviderID: sub,
		Email: email, EmailVerified: emailVerified, Name: name,
	}, nil
}
