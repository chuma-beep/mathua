package auth

import (
	"errors"
	"net/mail"
	"strings"
	"unicode"
)

// MinPasswordLength is the signup password policy floor.
const MinPasswordLength = 8

// MaxPasswordBytes caps passwords at bcrypt's 72-byte limit — longer inputs
// fail loudly here instead of erroring (or silently truncating) in bcrypt.
const MaxPasswordBytes = 72

// ValidateEmail accepts empty (optional recovery email) or a parseable address.
func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return nil
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return errors.New("email address is invalid")
	}
	return nil
}

// Username rules: 3-20 chars of lowercase letters, digits, underscore, dot.
// Normalization (NormalizeUsername) runs first, so this sees the final form.
func ValidateUsername(u string) error {
	if len([]rune(u)) < 3 || len([]rune(u)) > 20 {
		return errors.New("username must be 3-20 characters")
	}
	for _, r := range u {
		if r >= 'a' && r <= 'z' {
			continue
		}
		if unicode.IsDigit(r) {
			continue
		}
		if r == '_' || r == '.' {
			continue
		}
		return errors.New("username may only contain letters, numbers, underscore, and dot")
	}
	return nil
}

// ValidatePassword enforces the signup password policy: at least 8
// characters, including at least one number and one special
// (non-alphanumeric) character. Login is not affected — it only ever
// checks credentials against the stored hash.
func ValidatePassword(pw string) error {
	if len(pw) < MinPasswordLength {
		return errors.New("password must be at least 8 characters")
	}
	if len([]byte(pw)) > MaxPasswordBytes {
		return errors.New("password must be at most 72 characters")
	}
	hasDigit := false
	hasSpecial := false
	for _, r := range pw {
		switch {
		case unicode.IsDigit(r):
			hasDigit = true
		case !unicode.IsLetter(r):
			hasSpecial = true
		}
	}
	if !hasDigit {
		return errors.New("password must include at least one number")
	}
	if !hasSpecial {
		return errors.New("password must include at least one special character")
	}
	return nil
}
