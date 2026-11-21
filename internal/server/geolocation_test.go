package server

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/storage"
)

func TestAuthMiddlewareRecordsOneCountryPerDay(t *testing.T) {
	e := newAdminEnv(t)
	today := time.Now().UTC().Format("2006-01-02")

	send := func(code string) {
		t.Helper()
		req := httptest.NewRequest("GET", "/api/auth/identities", nil)
		req.Header.Set("Authorization", "Bearer "+e.studentTok)
		req.Header.Set("X-Mathua-Country", code)
		req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 18_0 like Mac OS X) Mobile/15E148")
		rec := httptest.NewRecorder()
		e.mux.ServeHTTP(rec, req)
		if rec.Code != 200 {
			t.Fatalf("identities with %q: %d %s", code, rec.Code, rec.Body.String())
		}
	}

	// The header arrives lowercase from no client on purpose: the Worker sends
	// uppercase, but the middleware is the boundary and normalises it.
	send("ng")
	// A second request the same day must not overwrite the first capture. A
	// learner who moves networks mid-day would otherwise turn the map into a
	// live tracking feed.
	send("GB")

	locations, err := e.store.StudentLocations(storage.LocationWindow{
		RecentFrom: today,
		PriorFrom:  today,
		DailyFrom:  today,
		Today:      today,
	})
	if err != nil {
		t.Fatalf("read locations: %v", err)
	}
	if locations.Known["NG"] != 1 {
		t.Fatalf("known NG = %d, want 1 (known %v)", locations.Known["NG"], locations.Known)
	}
	if locations.Known["GB"] != 0 {
		t.Fatalf("known GB = %d, want 0: the second same-day capture must be ignored", locations.Known["GB"])
	}
	if locations.ActiveTotal != 1 {
		t.Fatalf("active total = %d, want 1", locations.ActiveTotal)
	}
	if locations.Device["mobile"] != 1 {
		t.Fatalf("mobile learners = %d, want 1 (devices %v)", locations.Device["mobile"], locations.Device)
	}
}

func TestDeviceClassIsCoarseAndNeverTheAgent(t *testing.T) {
	cases := map[string]string{
		"": "",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 18_0 like Mac OS X) Mobile/15E148":             "mobile",
		"Mozilla/5.0 (Linux; Android 14; Pixel 8) AppleWebKit/537.36 Mobile Safari/537.36": "mobile",
		"Mozilla/5.0 (iPad; CPU OS 18_0 like Mac OS X) AppleWebKit/605.1.15":               "tablet",
		"Mozilla/5.0 (Linux; Android 14; SM-X710) AppleWebKit/537.36 Safari/537.36":        "tablet",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_5) AppleWebKit/605.1.15":                "desktop",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 Chrome/130.0 Safari/537.36":    "desktop",
	}
	for ua, want := range cases {
		if got := deviceClass(ua); got != want {
			t.Errorf("deviceClass(%q) = %q, want %q", ua, got, want)
		}
	}
}
