package auth

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"
)

func facebookAppID() string {
	if v := strings.TrimSpace(os.Getenv("FACEBOOK_APP_ID")); v != "" {
		return v
	}
	return strings.TrimSpace(os.Getenv("FACEBOOK_CLIENT_ID"))
}

func facebookAppSecret() string {
	if v := strings.TrimSpace(os.Getenv("FACEBOOK_APP_SECRET")); v != "" {
		return v
	}
	return strings.TrimSpace(os.Getenv("FACEBOOK_CLIENT_SECRET"))
}

// FacebookLoginURL starts the Meta Login flow (email needs app review).
func FacebookLoginURL(appID, redirectURI, state string) string {
	q := url.Values{}
	q.Set("client_id", appID)
	q.Set("redirect_uri", redirectURI)
	q.Set("state", state)
	q.Set("scope", "email,public_profile")
	q.Set("response_type", "code")
	return "https://www.facebook.com/v21.0/dialog/oauth?" + q.Encode()
}

// ExchangeFacebookCode trades the callback code for a normalized profile.
// Email may be absent (declined or unreviewed permission) — then unverified.
func ExchangeFacebookCode(ctx context.Context, code, redirectURI string) (OAuthProfile, error) {
	appID := facebookAppID()
	secret := facebookAppSecret()
	if appID == "" || secret == "" {
		return OAuthProfile{}, fmt.Errorf("FACEBOOK_APP_ID/SECRET not configured")
	}
	q := url.Values{}
	q.Set("client_id", appID)
	q.Set("redirect_uri", redirectURI)
	q.Set("client_secret", secret)
	q.Set("code", code)
	var tok struct {
		AccessToken string `json:"access_token"`
	}
	if err := oauthGetJSON(ctx, "https://graph.facebook.com/v21.0/oauth/access_token?"+q.Encode(), "", &tok); err != nil {
		return OAuthProfile{}, err
	}
	if tok.AccessToken == "" {
		return OAuthProfile{}, fmt.Errorf("no access_token from facebook")
	}
	meURL := "https://graph.facebook.com/v21.0/me?fields=id,name,email,picture.type(large)&access_token=" + url.QueryEscape(tok.AccessToken)
	var me struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Email   string `json:"email"`
		Picture struct {
			Data struct {
				URL string `json:"url"`
			} `json:"data"`
		} `json:"picture"`
	}
	if err := oauthGetJSON(ctx, meURL, "", &me); err != nil {
		return OAuthProfile{}, err
	}
	if me.ID == "" {
		return OAuthProfile{}, fmt.Errorf("facebook user missing id")
	}
	return OAuthProfile{
		Provider: "facebook", ProviderID: me.ID,
		Email: me.Email, EmailVerified: me.Email != "",
		Name: me.Name, AvatarURL: me.Picture.Data.URL,
	}, nil
}
