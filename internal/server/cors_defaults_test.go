package server

import (
	"sync"
	"testing"
)

// resetOrigins clears the cached parse so each case sees its own environment.
func resetOrigins() {
	allowedOriginsOnce = sync.Once{}
	allowedOriginsExact = nil
	allowedOriginPatterns = nil
	allowedOriginsAllowAll = false
}

// TestCORS_MissingConfigurationIsNotAllowAll is the D8 regression. An unset
// variable used to permit every origin, so a forgotten environment variable
// silently opened the API; local development must still work.
func TestCORS_MissingConfigurationIsNotAllowAll(t *testing.T) {
	t.Setenv("CORS_ALLOWED_ORIGINS", "")
	t.Setenv("FLY_APP_NAME", "mathua")
	resetOrigins()

	if isAllowedOrigin("https://evil.example") {
		t.Error("an unset CORS_ALLOWED_ORIGINS allowed a foreign origin")
	}
	if !isAllowedOrigin("http://localhost:3000") {
		t.Error("local development origin was refused")
	}
}

func TestCORS_UsesTheConfiguredList(t *testing.T) {
	t.Setenv("CORS_ALLOWED_ORIGINS", "https://mathua.com, https://*.example.com")
	resetOrigins()

	for _, ok := range []string{"https://mathua.com", "https://branch.example.com"} {
		if !isAllowedOrigin(ok) {
			t.Errorf("%s was refused", ok)
		}
	}
	if isAllowedOrigin("https://mathua.example") {
		t.Error("an origin outside the list was allowed")
	}
}

func TestCORS_UnusableListFallsBackToLocalOnly(t *testing.T) {
	t.Setenv("CORS_ALLOWED_ORIGINS", " , ")
	resetOrigins()

	if isAllowedOrigin("https://evil.example") {
		t.Error("an empty list allowed a foreign origin")
	}
	if !isAllowedOrigin("http://127.0.0.1:3000") {
		t.Error("local development origin was refused")
	}
}
