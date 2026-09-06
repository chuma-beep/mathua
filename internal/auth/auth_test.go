package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/storage"
)

func newTestService(t *testing.T) *AuthService {
	t.Helper()
	repo, err := storage.NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	t.Cleanup(func() { repo.Close() })
	return New(repo)
}

func TestNormalizeUsername(t *testing.T) {
	if got := NormalizeUsername("  Ada_Lovelace "); got != "ada_lovelace" {
		t.Errorf("expected trimmed lowercase, got %q", got)
	}
}

func TestSignupLogin_CaseInsensitive(t *testing.T) {
	svc := newTestService(t)
	token, st, err := svc.Signup("Ada", "Ada", "Engine!n1")
	if err != nil || token == "" {
		t.Fatalf("signup: %v", err)
	}
	if st.Username != "ada" {
		t.Errorf("expected stored lowercase username, got %q", st.Username)
	}
	for _, variant := range []string{"ada", "ADA", "  Ada  "} {
		tok, _, err := svc.Login(variant, "Engine!n1")
		if err != nil || tok == "" {
			t.Errorf("login(%q): expected success, got %v", variant, err)
		}
	}
}

func TestSignup_DuplicateUsername(t *testing.T) {
	svc := newTestService(t)
	if _, _, err := svc.Signup("Ada", "ada", "Engine!n1"); err != nil {
		t.Fatalf("first signup: %v", err)
	}
	_, _, err := svc.Signup("Other", "ADA", "Engine!n2")
	if err == nil {
		t.Fatal("expected duplicate error for case-variant username")
	}
	if err != ErrUsernameTaken && !storage.IsUniqueViolation(err) {
		t.Errorf("expected taken/unique-violation error, got: %v", err)
	}
}

func TestLogin_GoogleOnlyByEmail(t *testing.T) {
	svc := newTestService(t)
	if _, err := svc.repo.CreateGoogleUser("Gigi", "gigi@example.com", "gid-9", ""); err != nil {
		t.Fatalf("create google user: %v", err)
	}
	_, _, err := svc.Login("gigi@example.com", "whatever")
	if err != ErrGoogleOnly {
		t.Errorf("expected ErrGoogleOnly, got: %v", err)
	}
}

func TestLogin_UnknownEmailSilent(t *testing.T) {
	svc := newTestService(t)
	// Unknown identifiers fail silently (nil error, empty token) so the
	// handler can answer generic 401 without user enumeration.
	if _, _, err := svc.Login("nobody@example.com", "x"); err != nil {
		t.Errorf("expected silent miss (nil error), got: %v", err)
	}
}

func TestValidatePassword_MaxBytes(t *testing.T) {
	long := strings.Repeat("a1!", 25) // 75 bytes
	if err := ValidatePassword(long); err == nil {
		t.Error("expected rejection past 72 bytes")
	}
	ok := strings.Repeat("a1!", 24) // 72 bytes
	if err := ValidatePassword(ok); err != nil {
		t.Errorf("expected 72-byte password accepted, got %v", err)
	}
}

func TestPasswordReset_CompleteAndSingleUse(t *testing.T) {
	svc := newTestService(t)
	_, st, err := svc.Signup("Rita", "rita", "Engine!n1")
	if err != nil {
		t.Fatalf("signup: %v", err)
	}
	// Seed a token directly (no SMTP in tests).
	raw := "test-reset-token-123"
	if err := svc.repo.CreatePasswordReset(hashResetToken(raw), st.ID, time.Now().UTC().Add(ResetTokenTTL)); err != nil {
		t.Fatalf("seed token: %v", err)
	}
	tok, rst, err := svc.CompletePasswordReset(raw, "N3w!passw")
	if err != nil || tok == "" {
		t.Fatalf("complete: %v", err)
	}
	if rst.ID != st.ID {
		t.Errorf("expected student %q, got %q", st.ID, rst.ID)
	}
	// New password works.
	if _, _, err := svc.Login("rita", "N3w!passw"); err != nil {
		t.Errorf("login with new password: %v", err)
	}
	// Replay burns out.
	if _, _, err := svc.CompletePasswordReset(raw, "Other!n99"); err == nil {
		t.Error("expected replayed token rejected")
	}
}

func TestPasswordReset_ExpiredRejected(t *testing.T) {
	svc := newTestService(t)
	_, st, err := svc.Signup("Rita2", "rita2", "Engine!n1")
	if err != nil {
		t.Fatalf("signup: %v", err)
	}
	raw := "expired-token-1"
	if err := svc.repo.CreatePasswordReset(hashResetToken(raw), st.ID, time.Now().UTC().Add(-time.Minute)); err != nil {
		t.Fatalf("seed token: %v", err)
	}
	if _, _, err := svc.CompletePasswordReset(raw, "N3w!passw"); err == nil {
		t.Error("expected expired token rejected")
	}
}

func TestPasswordReset_RequestSilent(t *testing.T) {
	svc := newTestService(t)
	// Unknown identifier and missing mailer both look identical: no error.
	if _, err := svc.RequestPasswordReset("nobody-here"); err != nil {
		t.Errorf("expected silent unknown, got %v", err)
	}
	if _, err := svc.RequestPasswordReset(""); err != nil {
		t.Errorf("expected silent empty, got %v", err)
	}
}

func TestChangePassword(t *testing.T) {
	svc := newTestService(t)
	_, st, err := svc.Signup("Cid", "cid", "Engine!n1")
	if err != nil {
		t.Fatalf("signup: %v", err)
	}
	if err := svc.ChangePassword(st.ID, "Wrong!n9", "N3w!passw"); err == nil {
		t.Error("expected wrong-current rejected")
	}
	if err := svc.ChangePassword(st.ID, "Engine!n1", "weak"); err == nil {
		t.Error("expected weak-new rejected")
	}
	if err := svc.ChangePassword(st.ID, "Engine!n1", "N3w!passw"); err != nil {
		t.Fatalf("change: %v", err)
	}
	if _, _, err := svc.Login("cid", "N3w!passw"); err != nil {
		t.Errorf("login with changed password: %v", err)
	}
}
