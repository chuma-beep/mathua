package storage

import (
	"crypto/rand"
	"database/sql"
	"fmt"
)

// friendAlphabet omits characters that are easy to confuse when a learner
// reads a code aloud or types it from a screenshot: 0/O, 1/I and similar.
const friendAlphabet = "23456789ABCDEFGHJKLMNPQRSTUVWXYZ"

const (
	profileIDLength  = 16
	friendCodeLength = 8
)

func randomFriendString(n int) (string, error) {
	raw := make([]byte, n)
	if _, err := rand.Read(raw); err != nil {
		return "", fmt.Errorf("generate friend identity: %w", err)
	}
	out := make([]byte, n)
	for i, b := range raw {
		out[i] = friendAlphabet[int(b)%len(friendAlphabet)]
	}
	return string(out), nil
}

func newProfileID() (string, error) {
	body, err := randomFriendString(profileIDLength)
	if err != nil {
		return "", err
	}
	return "p_" + body, nil
}

func newFriendCode() (string, error) {
	return randomFriendString(friendCodeLength)
}

// backfillFriendIdentities gives every existing account a stable profile id and
// a rotatable friend code, then installs the partial unique indexes that make
// both lookup keys unambiguous. Empty values are legal until this runs because
// adding a unique NOT NULL column to a populated table needs a default first;
// the partial indexes keep rows written by older test helpers from colliding.
func backfillFriendIdentities(db *sql.DB, rebind func(string) string) error {
	rows, err := db.Query("SELECT id, profile_id, friend_code FROM students")
	if err != nil {
		return fmt.Errorf("backfill friend identities: %w", err)
	}
	type account struct {
		id, profileID, friendCode string
	}
	var accounts []account
	seenProfiles := map[string]bool{}
	seenCodes := map[string]bool{}
	for rows.Next() {
		var a account
		if err := rows.Scan(&a.id, &a.profileID, &a.friendCode); err != nil {
			rows.Close()
			return fmt.Errorf("backfill friend identities: %w", err)
		}
		if a.profileID != "" {
			seenProfiles[a.profileID] = true
		}
		if a.friendCode != "" {
			seenCodes[a.friendCode] = true
		}
		accounts = append(accounts, a)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return fmt.Errorf("backfill friend identities: %w", err)
	}
	rows.Close()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("backfill friend identities: %w", err)
	}
	for _, a := range accounts {
		profileID, friendCode := a.profileID, a.friendCode
		if profileID == "" {
			for {
				profileID, err = newProfileID()
				if err != nil {
					tx.Rollback()
					return err
				}
				if !seenProfiles[profileID] {
					break
				}
			}
		}
		if friendCode == "" {
			for {
				friendCode, err = newFriendCode()
				if err != nil {
					tx.Rollback()
					return err
				}
				if !seenCodes[friendCode] {
					break
				}
			}
		}
		if profileID == a.profileID && friendCode == a.friendCode {
			continue
		}
		if _, err := tx.Exec(rebind("UPDATE students SET profile_id = ?, friend_code = ? WHERE id = ?"), profileID, friendCode, a.id); err != nil {
			tx.Rollback()
			return fmt.Errorf("backfill friend identities: %w", err)
		}
		seenProfiles[profileID] = true
		seenCodes[friendCode] = true
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("backfill friend identities: %w", err)
	}

	for _, stmt := range []string{
		"CREATE UNIQUE INDEX IF NOT EXISTS uidx_students_profile_id ON students(profile_id) WHERE profile_id != ''",
		"CREATE UNIQUE INDEX IF NOT EXISTS uidx_students_friend_code ON students(friend_code) WHERE friend_code != ''",
	} {
		if _, err := db.Exec(stmt); err != nil {
			return fmt.Errorf("backfill friend identities: %w", err)
		}
	}
	return nil
}
