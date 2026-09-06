package auth

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"
)

func githubClientID() string {
	if v := strings.TrimSpace(os.Getenv("GITHUB_CLIENT_ID")); v != "" {
		return v
	}
	return strings.TrimSpace(os.Getenv("GITHUB_OAUTH_CLIENT_ID"))
}

func githubClientSecret() string {
	if v := strings.TrimSpace(os.Getenv("GITHUB_CLIENT_SECRET")); v != "" {
		return v
	}
	return strings.TrimSpace(os.Getenv("GITHUB_OAUTH_CLIENT_SECRET"))
}

// GithubLoginURL starts the GitHub OAuth flow (user:email gives verified flags).
func GithubLoginURL(clientID, redirectURI, state string) string {
	q := url.Values{}
	q.Set("client_id", clientID)
	q.Set("redirect_uri", redirectURI)
	q.Set("scope", "read:user user:email")
	q.Set("state", state)
	return "https://github.com/login/oauth/authorize?" + q.Encode()
}

type githubUser struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

type githubEmail struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

// ExchangeGithubCode trades the callback code for a normalized profile.
// Email comes from the verified-flagged emails API, preferring primary.
func ExchangeGithubCode(ctx context.Context, code, redirectURI string) (OAuthProfile, error) {
	clientID := githubClientID()
	secret := githubClientSecret()
	if clientID == "" || secret == "" {
		return OAuthProfile{}, fmt.Errorf("GITHUB_CLIENT_ID/SECRET not configured")
	}
	form := url.Values{}
	form.Set("code", code)
	form.Set("client_id", clientID)
	form.Set("client_secret", secret)
	form.Set("redirect_uri", redirectURI)
	var tok struct {
		AccessToken string `json:"access_token"`
		Error       string `json:"error"`
	}
	if err := oauthPostForm(ctx, "https://github.com/login/oauth/access_token", form, &tok); err != nil {
		return OAuthProfile{}, err
	}
	if tok.AccessToken == "" {
		return OAuthProfile{}, fmt.Errorf("no access_token from github")
	}
	var u githubUser
	if err := oauthGetJSON(ctx, "https://api.github.com/user", tok.AccessToken, &u); err != nil {
		return OAuthProfile{}, err
	}
	if u.ID == 0 {
		return OAuthProfile{}, fmt.Errorf("github user missing id")
	}
	var emails []githubEmail
	if err := oauthGetJSON(ctx, "https://api.github.com/user/emails", tok.AccessToken, &emails); err != nil {
		return OAuthProfile{}, err
	}
	email, verified := "", false
	for _, e := range emails {
		if e.Primary {
			email, verified = e.Email, e.Verified
			break
		}
	}
	if email == "" {
		for _, e := range emails {
			if e.Verified {
				email, verified = e.Email, true
				break
			}
		}
	}
	if email == "" && u.Email != "" {
		email = u.Email // public email is GitHub-confirmed
		verified = true
	}
	name := u.Name
	if name == "" {
		name = u.Login
	}
	return OAuthProfile{
		Provider: "github", ProviderID: fmt.Sprintf("%d", u.ID),
		Email: email, EmailVerified: verified, Name: name, AvatarURL: u.AvatarURL,
	}, nil
}
