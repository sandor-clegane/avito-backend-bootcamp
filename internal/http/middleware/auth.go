package middleware

import (
	"avito-backend-bootcamp/internal/model"
	ctxPkg "avito-backend-bootcamp/pkg/utils/ctx"
	"context"
	"net/http"
)

// NewAuthModerator creates a middleware that only allows Moderators
func NewAuthModerator(jwtManager JWT) func(next http.Handler) http.Handler {
	return authMiddleware(jwtManager, []model.UserType{model.Moderator})
}

// NewAuthModeratorOrClient creates a middleware that allows Moderators or Clients
func NewAuthModeratorOrClient(jwtManager JWT) func(next http.Handler) http.Handler {
	return authMiddleware(jwtManager, []model.UserType{model.Moderator, model.Client})
}

// authMiddleware is a base middleware function that handles common authorization logic
func authMiddleware(jwtManager JWT, allowedUserTypes []model.UserType) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Retrieve the token from the cookie
			cookie, err := r.Cookie("Authorization")
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Verify the token
			token, err := jwtManager.ParseToken(cookie.Value)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Get the audience (user type)
			audience, err := token.GetAudience()
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Validate the user type
			userType := model.MustParseUserType(audience[0])
			if !contains(allowedUserTypes, userType) {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Set user type in the context
			ctx := context.WithValue(r.Context(), ctxPkg.KeyUserType, userType)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// contains checks if a value is present in a slice
func contains(s []model.UserType, e model.UserType) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
