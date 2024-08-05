package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Manager struct {
	secretKey string
	tokenTTL  time.Duration
}

func New(secretKey string, tokenTTL time.Duration) *Manager {
	return &Manager{
		secretKey: secretKey,
		tokenTTL:  tokenTTL,
	}
}

// CreateToken creates JWT tokens with claims
func (m Manager) CreateToken(role string) (string, error) {
	// Create a new JWT token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "house-service",
		"aud": role,
		"exp": time.Now().Add(m.tokenTTL).Unix(),
		"iat": time.Now().Unix(),
	})

	// Sign the previously created token
	return token.SignedString([]byte(m.secretKey))
}

func (m Manager) ParseToken(tokenString string) (jwt.Claims, error) {
	// Parse the token with the secret key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return m.secretKey, nil
	})

	// Check for verification errors
	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Return the verified token
	return token.Claims, nil
}
