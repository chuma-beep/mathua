package server

import (
	"net/http/httptest"
	"testing"
)

func TestCanonicalGoogleLoginURL(t *testing.T) {
	const callback = "https://mathua.com/api/auth/google/callback"
	tests := []struct {
		name      string
		host      string
		forwarded string
		query     string
		want      string
	}{
		{
			name: "canonical host does nothing",
			host: "mathua.com",
			want: "",
		},
		{
			name:      "www is moved to canonical host",
			host:      "mathua.fly.dev",
			forwarded: "www.mathua.com",
			query:     "?return=%2Fprofile",
			want:      "https://mathua.com/api/auth/google/login?canonical=1&return=%2Fprofile",
		},
		{
			name:  "direct www is moved",
			host:  "www.mathua.com",
			query: "?return=%2Flearn",
			want:  "https://mathua.com/api/auth/google/login?canonical=1&return=%2Flearn",
		},
		{
			name:  "canonical marker breaks a proxy loop",
			host:  "mathua.fly.dev",
			query: "?canonical=1&return=%2Fprofile",
			want:  "",
		},
		{
			name: "empty redirect host is ignored",
			host: "www.mathua.com",
			want: "",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			redirect := callback
			if tc.name == "empty redirect host is ignored" {
				redirect = ""
			}
			req := httptest.NewRequest("GET", "http://example.test/api/auth/google/login"+tc.query, nil)
			req.Host = tc.host
			if tc.forwarded != "" {
				req.Header.Set("X-Forwarded-Host", tc.forwarded)
			}
			if got := canonicalGoogleLoginURL(req, redirect); got != tc.want {
				t.Fatalf("canonicalGoogleLoginURL = %q, want %q", got, tc.want)
			}
		})
	}
}
