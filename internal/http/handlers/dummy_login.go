package handlers

import (
	"log/slog"
	"net/http"

	"avito-backend-bootcamp/internal/model"
	pkgCtx "avito-backend-bootcamp/pkg/utils/ctx"
	"avito-backend-bootcamp/pkg/utils/sl"

	"github.com/go-chi/render"
)

type dummyLoginResponse struct {
	Token string `json:"token"`
}

// HandleDummyLogin упрощенный процесс получения токена для дальнейшего прохождения авторизации
func HandleDummyLogin(log *slog.Logger, authService AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Setup logger
		const op = "handlers.HandleDummyLogin"
		log := log.With(
			slog.String("op", op),
		)

		// Extract query param
		userTypeRaw := r.URL.Query().Get(pkgCtx.KeyUserType)

		// Parse user type from param
		userType, err := model.ParseUserType(userTypeRaw)
		if err != nil {
			log.Error("invalid query param", sl.Err(err))
			writeInternalError(r, w, err)
			return
		}

		// Get token
		token, err := authService.DummyLogin(r.Context(), userType)
		if err != nil {
			log.Error("dummy login failed", sl.Err(err))
			writeInternalError(r, w, err)
			return
		}

		// Set auth token to cookie
		http.SetCookie(w, &http.Cookie{
			Name:  "Authorization",
			Value: token,
		})

		// Return created token
		log.Info("dummy login success")
		render.Status(r, http.StatusOK)
		render.JSON(w, r, dummyLoginResponse{
			Token: token,
		})
	}
}
