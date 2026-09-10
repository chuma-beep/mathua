package auth

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"
)

const microsoftAuthority = "https://login.microsoftonline.com/common"

func microsoftClientID() string {
	if v := strings.TrimSpace(os.Getenv("MICROSOFT_CLIENT_ID")); v != "" {
		return v
	}
	return strings.TrimSpace(os.Getenv("MICROSOFT_OAUTH_CLIENT_ID"))
}

func microsoftClientSecret() string {
	if v := strings.TrimSpace(os.Getenv("MICROSOFT_CLIENT_SECRET")); v != "" {
		return v
	}
	return strings.TrimSpace(os.Getenv("MICROSOFT_OAUTH_CLIENT_SECRET"))
}

// MicrosoftLoginURL starts the Microsoft (personal + work/school) OIDC flow.
func MicrosoftLoginURL(clientID, redirectURI, state string) string {
	q := url.Values{}
	q.Set("client_id", clientID)
	q.Set("redirect_uri", redirectURI)
	q.Set("response_type", "code")
	q.Set("response_mode", "query")
	q.Set("scope", "openid profile email offline_access")
	q.Set("state", state)
	return microsoftAuthority + "/oauth2/v2.0/authorize?" + q.Encode()
}

// ExchangeMicrosoftCode trades the callback code for a normalized profile
// via the OIDC userinfo endpoint (no JWKS handling needed).
func ExchangeMicrosoftCode(ctx context.Context, code, redirectURI string) (OAuthProfile, error) {
	clientID := microsoftClientID()
	secret := microsoftClientSecret()
	if clientID == "" || secret == "" {
		return OAuthProfile{}, fmt.Errorf("MICROSOFT_CLIENT_ID/SECRET not configured")
	}
	form := url.Values{}
	form.Set("code", code)
	form.Set("client_id", clientID)
	form.Set("client_secret", secret)
	form.Set("redirect_uri", redirectURI)
	form.Set("grant_type", "authorization_code")
	var tok struct {
		AccessToken string `json:"access_token"`
	}
	if err := oauthPostForm(ctx, microsoftAuthority+"/oauth2/v2.0/token", form, &tok); err != nil {
		return OAuthProfile{}, err
	}
	if tok.AccessToken == "" {
		return OAuthProfile{}, fmt.Errorf("no access_token from microsoft")
	}
	var ui struct {
		Sub           string `json:"sub"`
		Email         string `json:"email"`
		Name          string `json:"name"`
		PreferredName string `json:"preferred_username"`
	}
	if err := oauthGetJSON(ctx, "https://graph.microsoft.com/oidc/userinfo", tok.AccessToken, &ui); err != nil {
		return OAuthProfile{}, err
	}
	if ui.Sub == "" {
		return OAuthProfile{}, fmt.Errorf("microsoft user missing sub")
	}
	name := ui.Name
	if name == "" {
		name = ui.PreferredName
	}
	return OAuthProfile{
		Provider: "microsoft", ProviderID: ui.Sub,
		Email: ui.Email, EmailVerified: ui.Email != "",
		Name: name,
	}, nil
}
