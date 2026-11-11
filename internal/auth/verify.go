package auth

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/chuma-beep/mathua/internal/mailer"
)

// EmailVerifyTTL bounds an email-verification link's lifetime.
const EmailVerifyTTL = 24 * time.Hour

// LinkTokenTTL bounds a provider-connect link token's lifetime.
const LinkTokenTTL = 10 * time.Minute

func sendMail(to, subject, body string) error {
	return mailer.New().Send(mailer.Message{To: to, Subject: subject, Text: body})
}

// RequestEmailVerification mails a verification link for the student's
// current email. Requires mail + frontend base configuration.
func (a *AuthService) RequestEmailVerification(studentID string) error {
	st, err := a.repo.GetStudent(studentID)
	if err != nil || st == nil {
		return errors.New("student not found")
	}
	if strings.TrimSpace(st.Email) == "" {
		return errors.New("no email on file, add one in Settings first")
	}
	if !smtpConfigured() || frontendBase() == "" {
		return errors.New("email verification not configured")
	}
	token, err := newResetToken()
	if err != nil {
		return err
	}
	if err := a.repo.CreateEmailVerification(hashResetToken(token), st.ID, time.Now().UTC().Add(EmailVerifyTTL)); err != nil {
		return err
	}
	link := strings.TrimSuffix(frontendBase(), "/") + "/verify-email?token=" + token
	body := "Verify your Mathua email address:\r\n\r\n" + link + "\r\n\r\n" +
		"This link is valid 24 hours, single use. If this wasn't you, ignore this email.\r\n"
	if err := sendMail(st.Email, "Mathua email verification", body); err != nil {
		log.Printf("auth: verify email failed: %v", err)
		return errors.New("could not send verification email")
	}
	return nil
}

// VerifyEmail burns a verification token and marks the email verified.
// Refuses with ErrEmailInUse when another verified row already holds the
// address (e.g. an OAuth fork created before this verification landed),
// the holder logs in and connects instead of forking a duplicate.
func (a *AuthService) VerifyEmail(token string) (string, error) {
	if strings.TrimSpace(token) == "" {
		return "", errors.New("verification token is required")
	}
	studentID, ok, err := a.repo.ConsumeEmailVerification(hashResetToken(strings.TrimSpace(token)))
	if err != nil {
		return "", err
	}
	if !ok {
		return "", errors.New("verification link is invalid or expired")
	}
	st, err := a.repo.GetStudent(studentID)
	if err != nil {
		return "", err
	}
	if st != nil && st.Email != "" {
		if existing, err := a.repo.FindByEmail(st.Email); err != nil {
			return "", err
		} else if existing != nil && existing.ID != studentID && existing.EmailVerified {
			return "", ErrEmailInUse
		}
	}
	if err := a.repo.SetEmailVerified(studentID, true); err != nil {
		return "", err
	}
	return studentID, nil
}

// CreateLinkToken mints a short-lived token binding a provider-connect
// OAuth dance (which carries no ambient session) to a logged-in student.
func (a *AuthService) CreateLinkToken(studentID string) (string, error) {
	token, err := newResetToken()
	if err != nil {
		return "", err
	}
	if err := a.repo.CreateLinkToken(hashResetToken(token), studentID, time.Now().UTC().Add(LinkTokenTTL)); err != nil {
		return "", err
	}
	return token, nil
}

// ConsumeLinkToken burns a raw link token, returning the bound student.
func (a *AuthService) ConsumeLinkToken(token string) (string, bool, error) {
	if strings.TrimSpace(token) == "" {
		return "", false, nil
	}
	return a.repo.ConsumeLinkToken(hashResetToken(strings.TrimSpace(token)))
}
