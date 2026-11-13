package server

import (
	"crypto/subtle"
	"net/http"
	"os"
	"strings"
)

// edgeSecretHeader is the header the Cloudflare Worker adds to every proxied
// API request. It is not a credential the browser holds; the Worker keeps the
// value in its environment and the Go origin compares it before trusting the
// client address the Worker forwards.
const edgeSecretHeader = "X-Mathua-Edge-Secret"

// EdgeGuard rejects API requests that did not arrive through the edge Worker
// when MATHUA_EDGE_SECRET is set.
//
// The Worker preserves the learner's real address in CF-Connecting-IP, which
// the server trusts through AI_TRUSTED_IP_HEADER. Without this guard a direct
// request to mathua.fly.dev could set that header itself and choose the
// rate-limit bucket it lands in. Health checks stay open: the platform probes
// the origin directly and has no secret to send.
//
// An empty secret disables the guard, which keeps local development and the
// test suite working without edge configuration.
func EdgeGuard(next http.Handler) http.Handler {
	secret := strings.TrimSpace(os.Getenv("MATHUA_EDGE_SECRET"))
	if secret == "" {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/api/") || r.URL.Path == "/api/health" {
			next.ServeHTTP(w, r)
			return
		}
		got := r.Header.Get(edgeSecretHeader)
		if subtle.ConstantTimeCompare([]byte(got), []byte(secret)) != 1 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`{"error":"forbidden"}`))
			return
		}
		next.ServeHTTP(w, r)
	})
}
