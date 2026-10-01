package auth

import (
	"testing"
	"time"
)

func TestOAuth_DualVerifiedLinks(t *testing.T) {
	svc := newTestService(t)
	repo := svc.repo
	st, err := repo.CreateUser("Ada", "ada", "hash")
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.SetEmail(st.ID, "ada@example.com"); err != nil {
		t.Fatal(err)
	}
	if err := repo.SetEmailVerified(st.ID, true); err != nil {
		t.Fatal(err)
	}
	tok, linked, err := svc.LoginOrCreateOAuth(OAuthProfile{
		Provider: "github", ProviderID: "4242",
		Email: "ada@example.com", EmailVerified: true, Name: "Ada",
	})
	if err != nil || tok == "" {
		t.Fatalf("expected link+login, got %v", err)
	}
	if linked.ID != st.ID {
		t.Errorf("expected merge into %q, got %q", st.ID, linked.ID)
	}
	ids, err := repo.ListIdentities(st.ID)
	if err != nil || len(ids) != 1 || ids[0].Provider != "github" {
		t.Errorf("expected one github identity, got %+v %v", ids, err)
	}
}

func TestOAuth_UnverifiedEmailNeverMerges(t *testing.T) {
	svc := newTestService(t)
	repo := svc.repo
	victim, err := repo.CreateUser("Victim", "victim", "hash")
	if err != nil {
		t.Fatal(err)
	}
	// Attacker claims the victim's email at signup (unverified).
	if err := repo.SetEmail(victim.ID, "victim@example.com"); err != nil {
		t.Fatal(err)
	}
	// Victim signs in with a verified provider email: must NOT merge — and
	// must not fork either. The unique index would reject a duplicate row,
	// so this returns ErrEmailInUse up front: log in with the password,
	// then connect the provider in Settings.
	_, fresh, err := svc.LoginOrCreateOAuth(OAuthProfile{
		Provider: "google", ProviderID: "google-uid-1",
		Email: "victim@example.com", EmailVerified: true, Name: "Victim",
	})
	if err != ErrEmailInUse {
		t.Fatalf("expected ErrEmailInUse, got %v (id %v)", err, fresh)
	}
}

func TestOAuth_KnownIdentityLogsIn(t *testing.T) {
	svc := newTestService(t)
	tok1, st1, err := svc.LoginOrCreateOAuth(OAuthProfile{
		Provider: "github", ProviderID: "777",
		Email: "octo@example.com", EmailVerified: true, Name: "Octo",
	})
	if err != nil || tok1 == "" {
		t.Fatalf("first login: %v", err)
	}
	tok2, st2, err := svc.LoginOrCreateOAuth(OAuthProfile{
		Provider: "github", ProviderID: "777",
		Email: "changed@example.com", EmailVerified: true, Name: "Octo Renamed",
	})
	if err != nil || tok2 == "" {
		t.Fatalf("second login: %v", err)
	}
	if st1.ID != st2.ID {
		t.Errorf("expected same account, got %q vs %q", st1.ID, st2.ID)
	}
}

func TestOAuth_ProviderIDsIsolated(t *testing.T) {
	svc := newTestService(t)
	_, a, err := svc.LoginOrCreateOAuth(OAuthProfile{Provider: "github", ProviderID: "1", Name: "A"})
	if err != nil {
		t.Fatal(err)
	}
	_, b, err := svc.LoginOrCreateOAuth(OAuthProfile{Provider: "facebook", ProviderID: "1", Name: "B"})
	if err != nil {
		t.Fatal(err)
	}
	if a.ID == b.ID {
		t.Error("same numeric id on different providers must not collide")
	}
}

func TestConnectDisconnect(t *testing.T) {
	svc := newTestService(t)
	repo := svc.repo
	st, err := repo.CreateUser("Ada", "ada2", "hash")
	if err != nil {
		t.Fatal(err)
	}
	p := OAuthProfile{Provider: "github", ProviderID: "99", Email: "a@example.com", EmailVerified: true}
	if err := svc.ConnectProvider(st.ID, p); err != nil {
		t.Fatalf("connect: %v", err)
	}
	// Idempotent reconnect.
	if err := svc.ConnectProvider(st.ID, p); err != nil {
		t.Errorf("reconnect should be idempotent: %v", err)
	}
	// Conflicting account.
	other, err := repo.CreateUser("Bo", "bo", "hash")
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.ConnectProvider(other.ID, p); err == nil {
		t.Error("expected conflict connecting an owned identity")
	}
	// Disconnect allowed (password remains).
	if err := svc.DisconnectProvider(st.ID, "github"); err != nil {
		t.Errorf("disconnect: %v", err)
	}
	// Last-credential guard: password-less + single identity cannot disconnect.
	pwless, err := repo.CreateGoogleUser("G", "g@example.com", "gid-x", "")
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.DisconnectProvider(pwless.ID, "google"); err == nil {
		t.Error("expected last-credential guard to refuse")
	}
}

func TestVerifyEmailFlow(t *testing.T) {
	svc := newTestService(t)
	repo := svc.repo
	st, err := repo.CreateUser("Em", "em", "hash")
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.SetEmail(st.ID, "em@example.com"); err != nil {
		t.Fatal(err)
	}
	raw := "verify-token-abc"
	if err := repo.CreateEmailVerification(hashResetToken(raw), st.ID, time.Now().UTC().Add(EmailVerifyTTL)); err != nil {
		t.Fatal(err)
	}
	sid, err := svc.VerifyEmail(raw)
	if err != nil || sid != st.ID {
		t.Fatalf("verify: %v", err)
	}
	got, _ := repo.GetStudent(st.ID)
	if !got.EmailVerified {
		t.Error("expected email_verified set")
	}
	if _, err := svc.VerifyEmail(raw); err == nil {
		t.Error("expected replay rejected")
	}
}

func TestLinkTokenFlow(t *testing.T) {
	svc := newTestService(t)
	st, err := svc.repo.CreateUser("Li", "li", "hash")
	if err != nil {
		t.Fatal(err)
	}
	raw, err := svc.CreateLinkToken(st.ID)
	if err != nil || raw == "" {
		t.Fatalf("mint: %v", err)
	}
	sid, ok, err := svc.ConsumeLinkToken(raw)
	if err != nil || !ok || sid != st.ID {
		t.Errorf("consume: %v %v %q", err, ok, sid)
	}
	if _, ok, _ := svc.ConsumeLinkToken(raw); ok {
		t.Error("expected single use")
	}
	if _, ok, _ := svc.ConsumeLinkToken(""); ok {
		t.Error("expected empty rejected")
	}
}

func TestOAuth_MergesPasswordlessUnverifiedRow(t *testing.T) {
	svc := newTestService(t)
	repo := svc.repo
	// OAuth-created row with a provider-unverified email: no password.
	first, err := repo.CreateOAuthUser("github", "gh-1", "G", "g@example.com", false, "")
	if err != nil {
		t.Fatal(err)
	}
	// Same mailbox, provider-verified, different provider → join, not fork.
	tok, st, err := svc.LoginOrCreateOAuth(OAuthProfile{
		Provider: "google", ProviderID: "google-uid-9",
		Email: "g@example.com", EmailVerified: true, Name: "G",
	})
	if err != nil || tok == "" {
		t.Fatalf("expected merge+login, got %v", err)
	}
	if st.ID != first.ID {
		t.Errorf("expected merge into %q, got %q", first.ID, st.ID)
	}
	ids, err := repo.ListIdentities(st.ID)
	if err != nil || len(ids) != 2 {
		t.Errorf("expected two linked identities, got %+v %v", ids, err)
	}
}

func TestOAuth_CaseInsensitiveMerge(t *testing.T) {
	svc := newTestService(t)
	repo := svc.repo
	st, err := repo.CreateUser("Ada", "ada-ci", "hash")
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.SetEmail(st.ID, "ada@example.com"); err != nil {
		t.Fatal(err)
	}
	if err := repo.SetEmailVerified(st.ID, true); err != nil {
		t.Fatal(err)
	}
	tok, linked, err := svc.LoginOrCreateOAuth(OAuthProfile{
		Provider: "github", ProviderID: "4243",
		Email: "ADA@EXAMPLE.COM", EmailVerified: true, Name: "Ada",
	})
	if err != nil || tok == "" {
		t.Fatalf("expected link+login, got %v", err)
	}
	if linked.ID != st.ID {
		t.Errorf("expected merge into %q, got %q", st.ID, linked.ID)
	}
	if got, _ := repo.GetStudent(st.ID); got.Email != "ada@example.com" {
		t.Errorf("expected stored email untouched, got %q", got.Email)
	}
}
