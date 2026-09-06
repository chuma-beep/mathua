package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/chuma-beep/mathua/internal/storage"
)

// ResetTokenTTL bounds a password-reset link's lifetime.
const ResetTokenTTL = time.Hour

func hashResetToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func newResetToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// smtpConfigured reports whether outbound reset mail can be sent.
func smtpConfigured() bool {
	return strings.TrimSpace(os.Getenv("SMTP_HOST")) != ""
}

// SMTPConfigured is the exported gate for the reset-request handler.
func SMTPConfigured() bool { return smtpConfigured() }

func sendResetEmail(to, link string) error {
	host := strings.TrimSpace(os.Getenv("SMTP_HOST"))
	port := strings.TrimSpace(os.Getenv("SMTP_PORT"))
	if port == "" {
		port = "587"
	}
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	from := os.Getenv("SMTP_FROM")
	if from == "" {
		from = user
	}
	var auth smtp.Auth
	if user != "" {
		auth = smtp.PlainAuth("", user, pass, host)
	}
	subject := "Subject: Mathua password reset\r\n"
	body := "Someone requested a password reset for your Mathua account.\r\n\r\n" +
		"Reset here (valid 1 hour, single use):\r\n" + link + "\r\n\r\n" +
		"If this wasn't you, ignore this email.\r\n"
	msg := []byte("From: " + from + "\r\nTo: " + to + "\r\n" + subject + "\r\n" + body)
	return smtp.SendMail(host+":"+port, auth, from, []string{to}, msg)
}

// frontendBase returns the absolute base for emailed links.
func frontendBase() string {
	return strings.TrimSpace(os.Getenv("FRONTEND_URL"))
}

// RequestPasswordReset issues a reset email when possible and is silent
// otherwise (unknown identifier, no email on file, unreachable mailer all
// look identical) to prevent account enumeration. The boolean reports
// whether mail was actually dispatched (observability only, never to clients).
func (a *AuthService) RequestPasswordReset(identifier string) (bool, error) {
	identifier = NormalizeUsername(identifier)
	st, err := a.repo.FindByUsername(identifier)
	if err != nil {
		return false, err
	}
	if st == nil && strings.Contains(identifier, "@") {
		st, err = a.repo.FindByEmail(identifier)
		if err != nil {
			return false, err
		}
	}
	if st == nil || strings.TrimSpace(st.Email) == "" {
		return false, nil
	}
	if !smtpConfigured() {
		log.Println("auth: password reset requested but SMTP_HOST unset — skipping send")
		return false, nil
	}
	base := frontendBase()
	if base == "" {
		log.Println("auth: password reset requested but FRONTEND_URL unset — skipping send")
		return false, nil
	}
	token, err := newResetToken()
	if err != nil {
		return false, err
	}
	if err := a.repo.CreatePasswordReset(hashResetToken(token), st.ID, time.Now().UTC().Add(ResetTokenTTL)); err != nil {
		return false, err
	}
	link := strings.TrimSuffix(base, "/") + "/login?reset=" + token
	if err := sendResetEmail(st.Email, link); err != nil {
		log.Printf("auth: reset email failed: %v", err)
		return false, nil
	}
	return true, nil
}

// CompletePasswordReset burns a single-use token and sets the new password,
// returning a fresh JWT (the user is logged in on success).
func (a *AuthService) CompletePasswordReset(token, newPassword string) (string, *storage.Student, error) {
	if strings.TrimSpace(token) == "" {
		return "", nil, errors.New("reset token is required")
	}
	if err := ValidatePassword(newPassword); err != nil {
		return "", nil, err
	}
	studentID, ok, err := a.repo.ConsumePasswordReset(hashResetToken(strings.TrimSpace(token)))
	if err != nil {
		return "", nil, err
	}
	if !ok {
		return "", nil, errors.New("reset link is invalid or expired")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", nil, err
	}
	if err := a.repo.SetPasswordHash(studentID, string(hash)); err != nil {
		return "", nil, err
	}
	st, err := a.repo.GetStudent(studentID)
	if err != nil || st == nil {
		return "", nil, fmt.Errorf("load student: %w", err)
	}
	tok, err := generateToken(st.ID)
	if err != nil {
		return "", nil, err
	}
	return tok, st, nil
}

// ChangePassword verifies the current password before setting the new one
// (logged-in flow — no email involved).
func (a *AuthService) ChangePassword(studentID, currentPassword, newPassword string) error {
	st, err := a.repo.GetStudent(studentID)
	if err != nil || st == nil {
		return fmt.Errorf("load student: %w", err)
	}
	if st.PasswordHash == "" {
		return errors.New("no password set — use password reset instead")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(st.PasswordHash), []byte(currentPassword)); err != nil {
		return errors.New("current password is incorrect")
	}
	if err := ValidatePassword(newPassword); err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return a.repo.SetPasswordHash(studentID, string(hash))
}
