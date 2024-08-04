package middleware

import (
	"fmt"

	"net/http"
)

func AuthModerator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Auth moderator middleware")
		next.ServeHTTP(w, r)
	})
}
