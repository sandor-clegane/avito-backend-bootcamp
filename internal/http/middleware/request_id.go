package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

const (
	RequestIDKey = "RequestID"
)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ID := uuid.New()
		ctx := context.WithValue(r.Context(), RequestIDKey, ID.String())
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
