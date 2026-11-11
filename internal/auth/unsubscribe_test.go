package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestNotificationUnsubscribeTokenRoundTrip(t *testing.T) {
	raw, err := NotificationUnsubscribeToken("stu_123")
	if err != nil {
		t.Fatalf("mint: %v", err)
	}
	got, err := ParseNotificationUnsubscribeToken(raw)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if got != "stu_123" {
		t.Fatalf("subject = %q, want stu_123", got)
	}
}

func TestNotificationUnsubscribeTokenRejectsBadInput(t *testing.T) {
	if _, err := NotificationUnsubscribeToken("  "); err == nil {
		t.Error("minted a token without a subject")
	}
	if _, err := ParseNotificationUnsubscribeToken(""); err == nil {
		t.Error("parsed an empty token")
	}
	if _, err := ParseNotificationUnsubscribeToken("not-a-token"); err == nil {
		t.Error("parsed a malformed token")
	}

	raw, err := NotificationUnsubscribeToken("stu_123")
	if err != nil {
		t.Fatalf("mint: %v", err)
	}
	tampered := raw[:len(raw)-1] + "x"
	if _, err := ParseNotificationUnsubscribeToken(tampered); err == nil {
		t.Error("accepted a tampered token")
	}
}

func TestNotificationUnsubscribeTokenRejectsWrongAudienceAndPurpose(t *testing.T) {
	now := time.Now().UTC()
	forged := jwt.NewWithClaims(jwt.SigningMethodHS256, notificationUnsubscribeClaims{
		Purpose: "session",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "stu_123",
			Issuer:    notificationUnsubscribeIssuer,
			Audience:  jwt.ClaimStrings{NotificationUnsubscribeAudience},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
		},
	})
	raw, err := forged.SignedString(getJWTSecret())
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	if _, err := ParseNotificationUnsubscribeToken(raw); err == nil || !strings.Contains(err.Error(), "invalid") {
		t.Fatalf("wrong purpose = %v, want invalid", err)
	}

	wrongAudience := jwt.NewWithClaims(jwt.SigningMethodHS256, notificationUnsubscribeClaims{
		Purpose: notificationUnsubscribePurpose,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "stu_123",
			Issuer:    notificationUnsubscribeIssuer,
			Audience:  jwt.ClaimStrings{"something-else"},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
		},
	})
	raw, err = wrongAudience.SignedString(getJWTSecret())
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	if _, err := ParseNotificationUnsubscribeToken(raw); err == nil {
		t.Error("accepted a token for another audience")
	}
}

func TestNotificationUnsubscribeTokenRejectsExpired(t *testing.T) {
	now := time.Now().UTC().Add(-48 * time.Hour)
	claims := notificationUnsubscribeClaims{
		Purpose: notificationUnsubscribePurpose,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "stu_123",
			Issuer:    notificationUnsubscribeIssuer,
			Audience:  jwt.ClaimStrings{NotificationUnsubscribeAudience},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
		},
	}
	raw, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(getJWTSecret())
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	if _, err := ParseNotificationUnsubscribeToken(raw); err == nil {
		t.Error("accepted an expired token")
	}
}
