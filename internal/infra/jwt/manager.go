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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"aud": role,
		"exp": time.Now().Add(m.tokenTTL).Unix(),
	})

	// Sign the previously created token
	signedToken, err := token.SignedString([]byte(m.secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (m Manager) ParseToken(tokenString string) (jwt.Claims, error) {
	// Parse the token with the secret key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.secretKey), nil
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
