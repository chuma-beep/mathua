package auth

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"strings"
)

// Random-username vocabulary for OAuth-created accounts (which have no
// user-chosen username). Every adjective/noun pair stays within the
// ValidateUsername rules: lowercase, 3–20 chars, [a-z0-9_.].
var usernameAdjectives = []string{
	"swift", "bright", "clever", "brave", "calm", "eager",
	"gentle", "happy", "keen", "lively", "merry", "nimble",
	"proud", "quick", "silly", "steady", "witty", "zesty",
	"cosmic", "frosty", "golden", "hidden", "ivory", "jolly",
}

var usernameNouns = []string{
	"fox", "bear", "wolf", "hawk", "otter", "raven",
	"badger", "cobalt", "dune", "ember", "fern", "glade",
	"heron", "ibis", "juniper", "koala", "lemur", "maple",
	"newt", "onyx", "panda", "quokka", "ridge", "sparrow",
}

func randIndex(n int) int {
	i, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
	if err != nil {
		return 0
	}
	return int(i.Int64())
}

// GenerateUsername returns an adjective_nounNN handle (9–19 chars) that
// always satisfies ValidateUsername. Uniqueness is NOT guaranteed — callers
// must check FindByUsername and retry (see AssignUsername).
func GenerateUsername() string {
	adj := usernameAdjectives[randIndex(len(usernameAdjectives))]
	noun := usernameNouns[randIndex(len(usernameNouns))]
	num := 10 + randIndex(90)
	name := adj + "_" + noun + fmt.Sprintf("%d", num)
	if err := ValidateUsername(name); err != nil {
		// Vocabulary bug — never ship a violating handle; fall back small.
		return "learner" + fmt.Sprintf("%d", 100+randIndex(900))
	}
	return name
}

// AssignUsername gives an OAuth-created (usernameless) student a unique
// random handle. Check-first with bounded retries; the collision window is
// tiny and usernames carry no UNIQUE constraint, so worst case is a rare
// duplicate — never a failed signup. Returns "" (and logs) when the repo
// errs, leaving the account usable exactly as before.
func (a *AuthService) AssignUsername(studentID string) string {
	for i := 0; i < 25; i++ {
		candidate := GenerateUsername()
		existing, err := a.repo.FindByUsername(candidate)
		if err != nil {
			log.Printf("auth: username availability check failed: %v", err)
			return ""
		}
		if existing != nil {
			continue
		}
		if err := a.repo.SetUsername(studentID, candidate); err != nil {
			log.Printf("auth: set username failed: %v", err)
			return ""
		}
		return candidate
	}
	log.Printf("auth: username assignment exhausted retries for %s", studentID)
	return ""
}

// EnsureUsername assigns a random handle only when the student has none
// (NULL or blank) — password accounts keep their chosen names.
func (a *AuthService) EnsureUsername(studentID string) string {
	st, err := a.repo.GetStudent(studentID)
	if err != nil || st == nil {
		return ""
	}
	if strings.TrimSpace(st.Username) != "" {
		return st.Username
	}
	return a.AssignUsername(studentID)
}

// BackfillUsernames assigns random handles to every student missing one
// (pre-random-username OAuth accounts). Idempotent — safe to run at boot.
func (a *AuthService) BackfillUsernames() (int, error) {
	ids, err := a.repo.ListStudentsMissingUsernames()
	if err != nil {
		return 0, err
	}
	assigned := 0
	for _, id := range ids {
		if u := a.EnsureUsername(id); u != "" {
			assigned++
		}
	}
	return assigned, nil
}
