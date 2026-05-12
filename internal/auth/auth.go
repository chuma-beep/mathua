package auth

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/chuma-beep/mathua/internal/storage"
)

var jwtSecret = generateSecret()

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
		return jwtSecret, nil
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
	return token.SignedString(jwtSecret)
}

func generateSecret() []byte {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return b
}

func GenerateSecretHex() string {
	return hex.EncodeToString(generateSecret())
}
