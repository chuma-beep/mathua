package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	// NotificationUnsubscribeAudience scopes an unsubscribe token so it cannot
	// be replayed as a session token even if both are signed with the same key.
	NotificationUnsubscribeAudience = "notification-unsubscribe"
	notificationUnsubscribeIssuer   = "mathua"
	notificationUnsubscribePurpose  = "notification_unsubscribe"
	notificationUnsubscribeTTL      = 30 * 24 * time.Hour
)

type notificationUnsubscribeClaims struct {
	Purpose string `json:"purpose"`
	jwt.RegisteredClaims
}

// NotificationUnsubscribeToken mints a signed, stateless link target for the
// unsubscribe footer of a notification email. It carries no password and no
// session authority; the only action it authorises is turning email delivery
// off for its subject.
func NotificationUnsubscribeToken(studentID string) (string, error) {
	studentID = strings.TrimSpace(studentID)
	if studentID == "" {
		return "", errors.New("unsubscribe token: student id is required")
	}
	now := time.Now().UTC()
	claims := notificationUnsubscribeClaims{
		Purpose: notificationUnsubscribePurpose,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   studentID,
			Issuer:    notificationUnsubscribeIssuer,
			Audience:  jwt.ClaimStrings{NotificationUnsubscribeAudience},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(notificationUnsubscribeTTL)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(getJWTSecret())
}

// ParseNotificationUnsubscribeToken verifies the signature, issuer, audience,
// purpose and expiry, and returns the learner id the link may change.
func ParseNotificationUnsubscribeToken(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", errors.New("unsubscribe token is required")
	}
	claims := &notificationUnsubscribeClaims{}
	tok, err := jwt.ParseWithClaims(raw, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok || t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("unexpected signing method")
		}
		return getJWTSecret(), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(notificationUnsubscribeIssuer),
		jwt.WithAudience(NotificationUnsubscribeAudience))
	if err != nil || !tok.Valid {
		return "", errors.New("unsubscribe link is invalid or expired")
	}
	if claims.Purpose != notificationUnsubscribePurpose || strings.TrimSpace(claims.Subject) == "" {
		return "", errors.New("unsubscribe link is invalid")
	}
	return claims.Subject, nil
}
