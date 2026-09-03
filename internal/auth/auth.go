package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/chuma-beep/mathua/internal/storage"
)

var (
	jwtSecretOnce sync.Once
	jwtSecret     []byte
)

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
	secret := generateSecret()
	log.Println("auth: generated random JWT secret (set JWT_SECRET for persistence)")
	return secret
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

func (a *AuthService) Signup(name, username, password string) (string, *storage.Student, error) {
	if err := ValidatePassword(password); err != nil {
		return "", nil, err
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
	st, err := a.repo.FindByUsername(username)
	if err != nil || st == nil {
		return "", nil, err
	}
	if st.PasswordHash == "" {
		return "", nil, nil
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
