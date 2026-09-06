package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
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

// LoginOrCreateGoogle links a Google identity by email (existing username/password users keep their username) and issues a JWT.
// One-tap (id_token) and OAuth code flow both resolve to (name, email, googleID, avatar) and converge here.
func (a *AuthService) LoginOrCreateGoogle(name, email, googleID, avatarURL string) (string, *storage.Student, error) {
	if googleID == "" {
		return "", nil, jwt.ErrTokenRequiredClaimMissing
	}
	// 1. Existing Google-linked account
	if st, err := a.repo.FindByGoogleID(googleID); err != nil {
		return "", nil, err
	} else if st != nil {
		// refresh avatar
		if avatarURL != "" && st.AvatarURL != avatarURL {
			_ = a.repo.LinkGoogleID(st.ID, googleID, avatarURL)
			st.AvatarURL = avatarURL
		}
		tok, err := generateToken(st.ID)
		if err != nil {
			return "", nil, err
		}
		return tok, st, nil
	}
	// 2. Link existing email/username account (keep username unique)
	if email != "" {
		if st, err := a.repo.FindByEmail(email); err != nil {
			return "", nil, err
		} else if st != nil {
			if err := a.repo.LinkGoogleID(st.ID, googleID, avatarURL); err != nil {
				return "", nil, err
			}
			tok, err := generateToken(st.ID)
			if err != nil {
				return "", nil, err
			}
			st.GoogleID = googleID
			st.AvatarURL = avatarURL
			if st.Email == "" {
				st.Email = email
			}
			return tok, st, nil
		}
	}
	// 3. Fresh Google user
	st, err := a.repo.CreateGoogleUser(name, email, googleID, avatarURL)
	if err != nil {
		return "", nil, err
	}
	tok, err := generateToken(st.ID)
	if err != nil {
		return "", nil, err
	}
	return tok, st, nil
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
