package auth

import (
	"crypto/rand"
	"encoding/hex"
	"log"
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

func (a *AuthService) ValidateToken(tokenStr string) (string, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return getJWTSecret(), nil
	})
	if err != nil {
		return "", err
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
