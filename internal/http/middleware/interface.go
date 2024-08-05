package middleware

import "github.com/golang-jwt/jwt/v5"

type JWT interface {
	ParseToken(tokenString string) (jwt.Claims, error)
}
