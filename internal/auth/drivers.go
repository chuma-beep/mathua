package auth

import (
	"context"
	"net/http"
)

// Driver abstracts one provider's authorize + code-exchange steps so the
// server shares a single start/callback handler pair across providers.
type Driver struct {
	// ClientID reports the configured client/app id ("" when unconfigured).
	ClientID func() string
	LoginURL func(clientID, redirectURI, state string) string
	Exchange func(ctx context.Context, code, redirectURI string, r *http.Request) (OAuthProfile, error)
}

// Drivers maps provider name → driver (google keeps its bespoke handlers).
func Drivers() map[string]Driver {
	return map[string]Driver{
		"github": {
			ClientID: githubClientID,
			LoginURL: GithubLoginURL,
			Exchange: func(ctx context.Context, code, redirectURI string, _ *http.Request) (OAuthProfile, error) {
				return ExchangeGithubCode(ctx, code, redirectURI)
			},
		},
		"facebook": {
			ClientID: facebookAppID,
			LoginURL: FacebookLoginURL,
			Exchange: func(ctx context.Context, code, redirectURI string, _ *http.Request) (OAuthProfile, error) {
				return ExchangeFacebookCode(ctx, code, redirectURI)
			},
		},
		"microsoft": {
			ClientID: microsoftClientID,
			LoginURL: MicrosoftLoginURL,
			Exchange: func(ctx context.Context, code, redirectURI string, _ *http.Request) (OAuthProfile, error) {
				return ExchangeMicrosoftCode(ctx, code, redirectURI)
			},
		},
		"apple": {
			ClientID: appleClientID,
			LoginURL: AppleLoginURL,
			Exchange: func(ctx context.Context, code, redirectURI string, r *http.Request) (OAuthProfile, error) {
				userJSON := ""
				if r != nil {
					userJSON = r.FormValue("user")
				}
				return ExchangeAppleCode(ctx, code, redirectURI, userJSON)
			},
		},
	}
}
