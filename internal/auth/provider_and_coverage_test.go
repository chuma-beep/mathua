package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestAuth_LoginOrCreateGoogle_NewAndReturning(t *testing.T) {
	svc := newTestService(t)

	// New user creates account
	tok1, st1, err := svc.LoginOrCreateGoogle("Ada Lovelace", "ada.google@example.com", "google-uid-12345", "https://example.com/avatar.png")
	if err != nil {
		t.Fatalf("LoginOrCreateGoogle new student: %v", err)
	}
	if tok1 == "" {
		t.Fatal("expected non-empty token")
	}
	if st1 == nil || st1.ID == "" {
		t.Fatal("expected non-nil student with ID")
	}
	if st1.Name != "Ada Lovelace" {
		t.Errorf("expected name 'Ada Lovelace', got %q", st1.Name)
	}
	if st1.Email != "ada.google@example.com" {
		t.Errorf("expected email 'ada.google@example.com', got %q", st1.Email)
	}
	if !st1.EmailVerified {
		t.Error("expected Google email to be verified")
	}

	// Returning user returns same student
	tok2, st2, err := svc.LoginOrCreateGoogle("Ada L.", "ada.google@example.com", "google-uid-12345", "https://example.com/avatar2.png")
	if err != nil {
		t.Fatalf("LoginOrCreateGoogle returning student: %v", err)
	}
	if tok2 == "" {
		t.Fatal("expected non-empty token")
	}
	if st2.ID != st1.ID {
		t.Errorf("expected same student ID %q, got %q", st1.ID, st2.ID)
	}
}

func TestAuth_ConnectAndDisconnectProvider_Scenarios(t *testing.T) {
	svc := newTestService(t)

	// Create student 1 with password
	_, st1, err := svc.Signup("Student One", "student1", "Password!123")
	if err != nil {
		t.Fatalf("signup st1: %v", err)
	}

	// Create student 2 with password
	_, st2, err := svc.Signup("Student Two", "student2", "Password!456")
	if err != nil {
		t.Fatalf("signup st2: %v", err)
	}

	// Empty provider or provider ID returns error
	if err := svc.ConnectProvider(st1.ID, OAuthProfile{Provider: "", ProviderID: "123"}); err == nil {
		t.Error("expected error for empty provider, got nil")
	}
	if err := svc.ConnectProvider(st1.ID, OAuthProfile{Provider: "github", ProviderID: ""}); err == nil {
		t.Error("expected error for empty provider ID, got nil")
	}

	// Connect GitHub to student 1
	pGithub := OAuthProfile{
		Provider: "github", ProviderID: "gh-1001", Email: "st1@github.example.com", EmailVerified: true,
	}
	if err := svc.ConnectProvider(st1.ID, pGithub); err != nil {
		t.Fatalf("connect github to st1: %v", err)
	}

	// Idempotent connect of same provider and ID to same student
	if err := svc.ConnectProvider(st1.ID, pGithub); err != nil {
		t.Errorf("expected idempotent connect to succeed, got %v", err)
	}

	// Connect same GitHub account to student 2 must fail
	if err := svc.ConnectProvider(st2.ID, pGithub); err == nil {
		t.Error("expected error connecting already-connected provider to different student, got nil")
	}

	// Disconnect provider from non-existent student must fail
	if err := svc.DisconnectProvider("nonexistent-student-id", "github"); err == nil {
		t.Error("expected error disconnecting non-existent student, got nil")
	}

	// Student 1 has a password, so disconnecting the only OAuth identity is allowed
	if err := svc.DisconnectProvider(st1.ID, "github"); err != nil {
		t.Fatalf("disconnect github from st1: %v", err)
	}

	// Now create student 3 via OAuth only (no password)
	_, st3, err := svc.LoginOrCreateOAuth(OAuthProfile{
		Provider: "google", ProviderID: "goog-2001", Email: "st3@example.com", EmailVerified: true,
	})
	if err != nil {
		t.Fatalf("create st3: %v", err)
	}

	// Trying to disconnect the only sign-in method for passwordless student must fail
	if err := svc.DisconnectProvider(st3.ID, "google"); err == nil {
		t.Error("expected error disconnecting last sign-in method, got nil")
	}

	// Connect a second provider (GitHub) to student 3
	pGithub3 := OAuthProfile{
		Provider: "github", ProviderID: "gh-3001", Email: "st3@github.example.com", EmailVerified: true,
	}
	if err := svc.ConnectProvider(st3.ID, pGithub3); err != nil {
		t.Fatalf("connect second provider to st3: %v", err)
	}

	// Now disconnecting one of the two providers is allowed
	if err := svc.DisconnectProvider(st3.ID, "google"); err != nil {
		t.Fatalf("disconnect google with remaining github provider: %v", err)
	}
}

func TestAuth_EnsureUsername_Scenarios(t *testing.T) {
	svc := newTestService(t)

	// Non-existent student returns empty string
	if u := svc.EnsureUsername("nonexistent-student-id"); u != "" {
		t.Errorf("expected empty string for non-existent student, got %q", u)
	}

	// Student with existing username keeps their username
	_, st, err := svc.Signup("Known Name", "custom_handle", "Password!123")
	if err != nil {
		t.Fatalf("signup: %v", err)
	}
	if u := svc.EnsureUsername(st.ID); u != "custom_handle" {
		t.Errorf("expected 'custom_handle', got %q", u)
	}

	// Student without username gets assigned random handle
	rawSt, err := svc.repo.CreateStudent("No Username Student")
	if err != nil {
		t.Fatalf("create student without username: %v", err)
	}
	assigned := svc.EnsureUsername(rawSt.ID)
	if strings.TrimSpace(assigned) == "" {
		t.Fatal("expected assigned username, got empty string")
	}
	if err := ValidateUsername(assigned); err != nil {
		t.Errorf("assigned username %q violates validation: %v", assigned, err)
	}

	// Calling EnsureUsername again returns the assigned handle
	if u := svc.EnsureUsername(rawSt.ID); u != assigned {
		t.Errorf("expected stable assigned handle %q, got %q", assigned, u)
	}
}

func TestAuth_ValidateToken_EdgeCases(t *testing.T) {
	svc := newTestService(t)

	// Empty and invalid tokens
	if _, err := svc.ValidateToken(""); err == nil {
		t.Error("expected error for empty token, got nil")
	}
	if _, err := svc.ValidateToken("invalid.jwt.token"); err == nil {
		t.Error("expected error for malformed token, got nil")
	}

	// Token signed with mismatched HMAC secret
	wrongSecret := []byte("this-is-a-completely-wrong-secret-32-bytes")
	claims := &Claims{
		StudentID: "student-123",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	wrongToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := wrongToken.SignedString(wrongSecret)
	if err != nil {
		t.Fatalf("sign wrong token: %v", err)
	}
	if _, err := svc.ValidateToken(tokenStr); err == nil {
		t.Error("expected error for token with wrong signature, got nil")
	}

	// Token with empty StudentID
	emptyClaims := &Claims{
		StudentID: "",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	emptyToken := jwt.NewWithClaims(jwt.SigningMethodHS256, emptyClaims)
	emptyTokenStr, err := emptyToken.SignedString(getJWTSecret())
	if err != nil {
		t.Fatalf("sign empty claims token: %v", err)
	}
	if _, err := svc.ValidateToken(emptyTokenStr); err == nil {
		t.Error("expected error for token missing StudentID claim, got nil")
	}
}

func TestAuth_RequestEmailVerification_EdgeCases(t *testing.T) {
	svc := newTestService(t)

	// Non-existent student
	if err := svc.RequestEmailVerification("nonexistent-student-id"); err == nil || !strings.Contains(err.Error(), "student not found") {
		t.Errorf("expected 'student not found' error, got %v", err)
	}

	// Student with empty email
	st, err := svc.repo.CreateStudent("No Email")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}
	if err := svc.RequestEmailVerification(st.ID); err == nil || !strings.Contains(err.Error(), "no email on file") {
		t.Errorf("expected 'no email on file' error, got %v", err)
	}

	// Student with email but unconfigured SMTP
	if err := svc.repo.SetEmail(st.ID, "student@example.com"); err != nil {
		t.Fatalf("set email: %v", err)
	}
	t.Setenv("SMTP_HOST", "")
	t.Setenv("SENDBYTE_KEY", "")
	if err := svc.RequestEmailVerification(st.ID); err == nil || !strings.Contains(err.Error(), "not configured") {
		t.Errorf("expected 'not configured' error, got %v", err)
	}
}

func TestAuth_VerifyEmail_WhitespaceAndTrimming(t *testing.T) {
	svc := newTestService(t)

	// Whitespace-only token
	if _, err := svc.VerifyEmail("   "); err == nil || err.Error() != "verification token is required" {
		t.Fatalf("expected 'verification token is required', got: %v", err)
	}

	// Create student with email
	st, err := svc.repo.CreateStudent("Trim Student")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}
	if err := svc.repo.SetEmail(st.ID, "trim@example.com"); err != nil {
		t.Fatalf("set email: %v", err)
	}

	// Mint token
	rawToken, err := newResetToken()
	if err != nil {
		t.Fatalf("newResetToken: %v", err)
	}
	if err := svc.repo.CreateEmailVerification(hashResetToken(rawToken), st.ID, time.Now().UTC().Add(EmailVerifyTTL)); err != nil {
		t.Fatalf("CreateEmailVerification: %v", err)
	}

	// Verify with padded whitespace
	verifiedID, err := svc.VerifyEmail("   " + rawToken + "   \n")
	if err != nil {
		t.Fatalf("VerifyEmail with whitespace failed: %v", err)
	}
	if verifiedID != st.ID {
		t.Errorf("expected student ID %q, got %q", st.ID, verifiedID)
	}
}
