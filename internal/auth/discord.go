package auth

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const discordAPI = "https://discord.com/api"

func discordClientID() string {
	return strings.TrimSpace(os.Getenv("DISCORD_CLIENT_ID"))
}

func discordClientSecret() string {
	return strings.TrimSpace(os.Getenv("DISCORD_CLIENT_SECRET"))
}

// DiscordLoginURL starts the Discord OAuth2 flow. `identify` returns the
// account id, name and avatar; `email` adds the joined address, which is what
// lets an existing account merge and lets reminders reach the learner.
func DiscordLoginURL(clientID, redirectURI, state string) string {
	q := url.Values{}
	q.Set("client_id", clientID)
	q.Set("redirect_uri", redirectURI)
	q.Set("response_type", "code")
	q.Set("scope", "identify email")
	q.Set("state", state)
	return "https://discord.com/oauth2/authorize?" + q.Encode()
}

type discordUser struct {
	ID         string `json:"id"`
	Username   string `json:"username"`
	GlobalName string `json:"global_name"`
	Email      string `json:"email"`
	Verified   bool   `json:"verified"`
	Avatar     string `json:"avatar"`
}

// discordAvatar builds the CDN address for a custom avatar, or the numbered
// default the account was assigned. Discord returns only the hash, so the
// path is assembled here rather than by the browser.
func discordAvatar(u discordUser) string {
	if u.Avatar != "" {
		return fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png?size=128", u.ID, u.Avatar)
	}
	id, err := strconv.ParseUint(u.ID, 10, 64)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("https://cdn.discordapp.com/embed/avatars/%d.png", (id>>22)%6)
}

// ExchangeDiscordCode trades the callback code for a profile. Discord accepts
// the client secret in the body, so no Authorization header is needed.
func ExchangeDiscordCode(ctx context.Context, code, redirectURI string) (OAuthProfile, error) {
	clientID := discordClientID()
	secret := discordClientSecret()
	if clientID == "" || secret == "" {
		return OAuthProfile{}, fmt.Errorf("DISCORD_CLIENT_ID/SECRET not configured")
	}
	form := url.Values{}
	form.Set("client_id", clientID)
	form.Set("client_secret", secret)
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("redirect_uri", redirectURI)
	var tok struct {
		AccessToken string `json:"access_token"`
	}
	if err := oauthPostForm(ctx, discordAPI+"/oauth2/token", form, &tok); err != nil {
		return OAuthProfile{}, err
	}
	if tok.AccessToken == "" {
		return OAuthProfile{}, fmt.Errorf("no access_token from discord")
	}
	var u discordUser
	if err := oauthGetJSON(ctx, discordAPI+"/users/@me", tok.AccessToken, &u); err != nil {
		return OAuthProfile{}, err
	}
	if u.ID == "" {
		return OAuthProfile{}, fmt.Errorf("discord user missing id")
	}
	name := u.GlobalName
	if name == "" {
		name = u.Username
	}
	return OAuthProfile{
		Provider: "discord", ProviderID: u.ID,
		Email: u.Email, EmailVerified: u.Verified && u.Email != "",
		Name: name, AvatarURL: discordAvatar(u),
	}, nil
}
