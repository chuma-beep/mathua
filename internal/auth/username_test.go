package auth

import (
	"strings"
	"testing"
)

func TestGenerateUsername_ConformsToPolicy(t *testing.T) {
	seen := map[string]bool{}
	for i := 0; i < 500; i++ {
		u := GenerateUsername()
		if err := ValidateUsername(u); err != nil {
			t.Fatalf("generated %q violates policy: %v", u, err)
		}
		seen[u] = true
	}
	if len(seen) < 400 {
		t.Errorf("expected varied handles, got %d unique of 500", len(seen))
	}
}

func TestOAuthAssignsRandomUsername(t *testing.T) {
	svc := newTestService(t)
	_, st1, err := svc.LoginOrCreateOAuth(OAuthProfile{
		Provider: "github", ProviderID: "g1",
		Email: "g1@example.com", EmailVerified: true, Name: "G One",
	})
	if err != nil {
		t.Fatalf("oauth login: %v", err)
	}
	if strings.TrimSpace(st1.Username) == "" {
		t.Fatal("expected a random username on OAuth creation")
	}
	if err := ValidateUsername(st1.Username); err != nil {
		t.Errorf("assigned %q violates policy: %v", st1.Username, err)
	}
	_, st2, err := svc.LoginOrCreateOAuth(OAuthProfile{
		Provider: "github", ProviderID: "g2",
		Email: "g2@example.com", EmailVerified: true, Name: "G Two",
	})
	if err != nil {
		t.Fatalf("oauth login 2: %v", err)
	}
	if st2.Username == st1.Username {
		t.Errorf("expected distinct handles, both got %q", st1.Username)
	}
	// Returning login keeps the handle (no reassignment).
	_, st1b, err := svc.LoginOrCreateOAuth(OAuthProfile{Provider: "github", ProviderID: "g1"})
	if err != nil {
		t.Fatalf("oauth relogin: %v", err)
	}
	if st1b.Username != st1.Username {
		t.Errorf("expected stable handle %q, got %q", st1.Username, st1b.Username)
	}
}

func TestBackfillUsernames(t *testing.T) {
	svc := newTestService(t)
	// Pre-generator account with a blank username (direct repo creation).
	st, err := svc.repo.CreateOAuthUser("google", "legacy1", "Legacy", "legacy@example.com", true, "")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if n, err := svc.BackfillUsernames(); err != nil || n != 1 {
		t.Fatalf("backfill: n=%d err=%v", n, err)
	}
	got, _ := svc.repo.GetStudent(st.ID)
	if strings.TrimSpace(got.Username) == "" {
		t.Fatal("expected backfilled username")
	}
	if n, err := svc.BackfillUsernames(); err != nil || n != 0 {
		t.Errorf("second backfill: expected idempotent 0, got n=%d err=%v", n, err)
	}
}

func TestSignup_RequiresTrimmedBoundedName(t *testing.T) {
	svc := newTestService(t)
	for _, tc := range []struct {
		name, username, want string
	}{
		{"", "noname", "name is required"},
		{"   ", "spaces", "name is required"},
		{strings.Repeat("x", 51), "longname", "name must be 1-50 characters"},
	} {
		if _, _, err := svc.Signup(tc.name, tc.username, "Engine!n1"); err == nil || err.Error() != tc.want {
			t.Errorf("signup(%q): expected %q, got %v", tc.name, tc.want, err)
		}
	}
	if _, st, err := svc.Signup("  Ada  ", "spaced", "Engine!n1"); err != nil || st.Name != "Ada" {
		t.Errorf("expected trimmed stored name, got %q err=%v", st.Name, err)
	}
}
