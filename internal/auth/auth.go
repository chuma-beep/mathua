package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/chuma-beep/mathua/internal/storage"
)

var (
	jwtSecretOnce sync.Once
	jwtSecret     []byte
	// secretFile, when set via SetSecretFile before first use, persists a
	// generated secret across restarts (zero-config logins survive redeploys
	// on single-machine volumes). Explicit JWT_SECRET env always wins.
	secretFile string
)

// SetSecretFile configures the path used to persist an auto-generated JWT
// secret. Must be called before first token operation; later calls are ignored.
func SetSecretFile(path string) {
	if path == "" {
		return
	}
	// Only honored before the secret is loaded once.
	if jwtSecret != nil {
		return
	}
	secretFile = path
}

func loadSecret() []byte {
	if s := os.Getenv("JWT_SECRET"); s != "" {
		if b, err := hex.DecodeString(s); err == nil && len(b) == 32 {
			log.Println("auth: using JWT_SECRET from environment (hex-encoded, 32 bytes)")
			return b
		}
		if len(s) < 32 {
			log.Fatalf("auth: JWT_SECRET must be at least 32 characters long (got %d)", len(s))
		}
		log.Println("auth: using JWT_SECRET from environment")
		return []byte(s)
	}
	if p := os.Getenv("JWT_SECRET_FILE"); p != "" {
		if b, ok := loadSecretFile(p); ok {
			return b
		}
	}
	if secretFile != "" {
		if b, ok := loadSecretFile(secretFile); ok {
			return b
		}
	}
	secret := generateSecret()
	log.Println("auth: generated random JWT secret (set JWT_SECRET for persistence)")
	return secret
}

// loadSecretFile loads a persisted hex secret, generating and storing one on
// first boot. File holds 32 bytes hex-encoded with 0600 permissions.
func loadSecretFile(path string) ([]byte, bool) {
	if b, err := os.ReadFile(path); err == nil {
		raw := strings.TrimSpace(string(b))
		if decoded, derr := hex.DecodeString(raw); derr == nil && len(decoded) == 32 {
			log.Printf("auth: using persisted JWT secret from %s", path)
			return decoded, true
		}
		// Corrupt: quarantine aside and fall through to fresh generation.
		backup := path + ".bad." + time.Now().UTC().Format("20060102T150405Z")
		if rerr := os.Rename(path, backup); rerr != nil {
			log.Printf("auth: persisted JWT secret at %s is corrupt and cannot be quarantined: %v", path, rerr)
			return nil, false
		}
		log.Printf("auth: persisted JWT secret at %s was corrupt (quarantined to %s); generating fresh", path, backup)
	}
	if dir := filepath.Dir(path); dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			log.Printf("auth: cannot create secret dir %s: %v", dir, err)
			return nil, false
		}
	}
	fresh := GenerateSecretHex()
	// O_EXCL: a concurrent first boot may win; adopt the winner's secret.
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		if os.IsExist(err) {
			if b, rerr := os.ReadFile(path); rerr == nil {
				if decoded, derr := hex.DecodeString(strings.TrimSpace(string(b))); derr == nil && len(decoded) == 32 {
					log.Printf("auth: using persisted JWT secret from %s (concurrent boot)", path)
					return decoded, true
				}
			}
		}
		log.Printf("auth: cannot persist JWT secret to %s: %v", path, err)
		return nil, false
	}
	if _, werr := f.WriteString(fresh + "\n"); werr != nil {
		log.Printf("auth: cannot write JWT secret to %s: %v", path, werr)
		f.Close()
		return nil, false
	}
	if cerr := f.Close(); cerr != nil {
		log.Printf("auth: cannot close JWT secret file %s: %v", path, cerr)
		return nil, false
	}
	log.Printf("auth: generated and persisted JWT secret to %s", path)
	decoded, _ := hex.DecodeString(fresh)
	return decoded, true
}

func getJWTSecret() []byte {
	jwtSecretOnce.Do(func() {
		jwtSecret = loadSecret()
	})
	return jwtSecret
}

type Claims struct {
	StudentID string `json:"sub"`
	jwt.RegisteredClaims
}

type AuthService struct {
	repo storage.Repository
}

func New(repo storage.Repository) *AuthService {
	return &AuthService{repo: repo}
}

func SecretMinLength() int { return 32 }

// ErrGoogleOnly marks accounts created via Google Sign-In, which have no
// password — password login must point them at Google, not a generic error.
var ErrGoogleOnly = errors.New("this account uses Google sign-in")

// NormalizeUsername trims and lowercases: "Ada" and "ada" are one account.
func NormalizeUsername(u string) string { return strings.ToLower(strings.TrimSpace(u)) }

// ErrUsernameTaken marks a taken username (409, not a 400/500).
var ErrUsernameTaken = errors.New("username is taken")

func (a *AuthService) Signup(name, username, password string) (string, *storage.Student, error) {
	// Defense in depth: the handler and form already require a name, but
	// raw API/service callers bypass them — trim and enforce here so no
	// blank or oversized name ever reaches the leaderboard.
	name = strings.TrimSpace(name)
	if name == "" {
		return "", nil, errors.New("name is required")
	}
	if len([]rune(name)) > 50 {
		return "", nil, errors.New("name must be 1-50 characters")
	}
	username = NormalizeUsername(username)
	if username == "" {
		return "", nil, errors.New("username is required")
	}
	if err := ValidateUsername(username); err != nil {
		return "", nil, err
	}
	if err := ValidatePassword(password); err != nil {
		return "", nil, err
	}
	// Usernames have no UNIQUE constraint historically — check first so
	// duplicates (including case variants) get a clean 409, not a DB error.
	if existing, err := a.repo.FindByUsername(username); err != nil {
		return "", nil, err
	} else if existing != nil {
		return "", nil, ErrUsernameTaken
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil, err
	}
	st, err := a.repo.CreateUser(name, username, string(hash))
	if err != nil {
		return "", nil, err
	}
	token, err := generateToken(st.ID)
	if err != nil {
		return "", nil, err
	}
	return token, st, nil
}

func (a *AuthService) Login(username, password string) (string, *storage.Student, error) {
	username = NormalizeUsername(username)
	st, err := a.repo.FindByUsername(username)
	if err != nil {
		return "", nil, err
	}
	if st == nil && strings.Contains(username, "@") {
		// Email login attempt: route Google-only accounts to the
		// Google-sign-in message instead of a generic failure.
		st, err = a.repo.FindByEmail(username)
		if err != nil {
			return "", nil, err
		}
	}
	if st == nil {
		return "", nil, nil
	}
	if st.PasswordHash == "" {
		return "", nil, ErrGoogleOnly
	}
	if err := bcrypt.CompareHashAndPassword([]byte(st.PasswordHash), []byte(password)); err != nil {
		return "", nil, nil
	}
	token, err := generateToken(st.ID)
	if err != nil {
		return "", nil, err
	}
	return token, st, nil
}

// OAuthProfile is the normalized identity from any OAuth provider
// (Google, GitHub, Facebook, Microsoft, Apple) after its own verification.
type OAuthProfile struct {
	Provider      string // google | github | facebook | microsoft | apple
	ProviderID    string // provider-side stable user id
	Email         string
	EmailVerified bool // provider-attested (e.g. GitHub verified flag, Google)
	Name          string
	AvatarURL     string
}

// LoginOrCreateOAuth is the provider-agnostic login core:
//  1. known (provider, id) identity → log in (refresh avatar);
//  2. dual-verified email (provider-attested AND ours) → link + log in;
//  3. otherwise → fresh account. Unverified emails NEVER merge — this is
//     what stops account takeover via claimed emails.
func (a *AuthService) LoginOrCreateOAuth(p OAuthProfile) (string, *storage.Student, error) {
	if p.Provider == "" || p.ProviderID == "" {
		return "", nil, jwt.ErrTokenRequiredClaimMissing
	}
	if st, err := a.repo.FindStudentByIdentity(p.Provider, p.ProviderID); err != nil {
		return "", nil, err
	} else if st != nil {
		if p.AvatarURL != "" && st.AvatarURL != p.AvatarURL {
			if err := a.repo.SetAvatarURL(st.ID, p.AvatarURL); err != nil {
				return "", nil, err
			}
			st.AvatarURL = p.AvatarURL
		}
		tok, err := generateToken(st.ID)
		if err != nil {
			return "", nil, err
		}
		return tok, st, nil
	}
	if p.Email != "" && p.EmailVerified {
		if st, err := a.repo.FindByEmail(p.Email); err != nil {
			return "", nil, err
		} else if st != nil && st.EmailVerified {
			if err := a.repo.CreateIdentity(p.Provider, p.ProviderID, st.ID, p.Email, true); err != nil {
				return "", nil, err
			}
			if p.AvatarURL != "" && st.AvatarURL == "" {
				_ = a.repo.SetAvatarURL(st.ID, p.AvatarURL)
				st.AvatarURL = p.AvatarURL
			}
			tok, err := generateToken(st.ID)
			if err != nil {
				return "", nil, err
			}
			return tok, st, nil
		}
	}
	st, err := a.repo.CreateOAuthUser(p.Provider, p.ProviderID, p.Name, p.Email, p.EmailVerified, p.AvatarURL)
	if err != nil {
		return "", nil, err
	}
	// OAuth accounts have no user-chosen username — assign a random handle
	// so every account carries one (leaderboard fallback, future @mentions).
	// Best-effort: a failure leaves the account exactly as before.
	if u := a.EnsureUsername(st.ID); u != "" {
		st.Username = u
	}
	tok, err := generateToken(st.ID)
	if err != nil {
		return "", nil, err
	}
	return tok, st, nil
}

// LoginOrCreateGoogle keeps the Google call sites stable — both One-Tap and
// the redirect flow converge here. Google emails are provider-verified.
func (a *AuthService) LoginOrCreateGoogle(name, email, googleID, avatarURL string) (string, *storage.Student, error) {
	return a.LoginOrCreateOAuth(OAuthProfile{
		Provider: "google", ProviderID: googleID, Email: email, EmailVerified: email != "",
		Name: name, AvatarURL: avatarURL,
	})
}

// ConnectProvider links a provider identity to an already-logged-in student
// (explicit user intent — safe without email verification). Fails when the
// provider account already belongs to someone else.
func (a *AuthService) ConnectProvider(studentID string, p OAuthProfile) error {
	if p.Provider == "" || p.ProviderID == "" {
		return jwt.ErrTokenRequiredClaimMissing
	}
	if st, err := a.repo.FindStudentByIdentity(p.Provider, p.ProviderID); err != nil {
		return err
	} else if st != nil {
		if st.ID == studentID {
			return nil // idempotent
		}
		return errors.New("that account is already connected elsewhere")
	}
	return a.repo.CreateIdentity(p.Provider, p.ProviderID, studentID, p.Email, p.EmailVerified)
}

// DisconnectProvider removes a link, refusing to strand an account with no
// remaining credential (no password and a single identity).
func (a *AuthService) DisconnectProvider(studentID, provider string) error {
	st, err := a.repo.GetStudent(studentID)
	if err != nil || st == nil {
		return fmt.Errorf("load student: %w", err)
	}
	ids, err := a.repo.ListIdentities(studentID)
	if err != nil {
		return err
	}
	keep := false
	for _, id := range ids {
		if id.Provider != provider {
			keep = true
			break
		}
	}
	if !keep && st.PasswordHash == "" {
		return errors.New("cannot disconnect the last sign-in method")
	}
	return a.repo.DeleteIdentity(provider, studentID)
}

func (a *AuthService) ValidateToken(tokenStr string) (string, error) {
	claims := &Claims{}
	tok, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, jwt.ErrSignatureInvalid
		}
		return getJWTSecret(), nil
	})
	if err != nil {
		return "", err
	}
	if !tok.Valid {
		return "", jwt.ErrSignatureInvalid
	}
	if claims.StudentID == "" {
		return "", jwt.ErrTokenRequiredClaimMissing
	}
	return claims.StudentID, nil
}

func generateToken(studentID string) (string, error) {
	claims := &Claims{
		StudentID: studentID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}

// IssueGuestToken mints a bearer token for a guest student row. It carries
// no guest marker — ownership is uniform: the token proves the caller owns
// the student ID, registered or not. Package-level (not a method) so the
// guest handler works even when the AuthService is unwired (dev/test);
// the secret machinery is process-global either way.
func IssueGuestToken(studentID string) (string, error) {
	if strings.TrimSpace(studentID) == "" {
		return "", errors.New("student_id is required")
	}
	return generateToken(studentID)
}

func generateSecret() []byte {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		log.Fatalf("auth: failed to generate random secret: %v", err)
	}
	return b
}

func GenerateSecretHex() string {
	return hex.EncodeToString(generateSecret())
}

func (a *AuthService) VerifyGoogleIDToken(ctx context.Context, idToken string) (*GoogleProfile, error) {
	return VerifyGoogleIDToken(ctx, idToken)
}

func (a *AuthService) ExchangeGoogleCode(ctx context.Context, code string, r *http.Request) (*GoogleProfile, error) {
	return ExchangeGoogleCode(ctx, code, r)
}
