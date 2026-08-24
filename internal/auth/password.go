package auth

import (
	"errors"
	"unicode"
)

// MinPasswordLength is the signup password policy floor.
const MinPasswordLength = 8

// ValidatePassword enforces the signup password policy: at least 8
// characters, including at least one number and one special
// (non-alphanumeric) character. Login is not affected — it only ever
// checks credentials against the stored hash.
func ValidatePassword(pw string) error {
	if len(pw) < MinPasswordLength {
		return errors.New("password must be at least 8 characters")
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
