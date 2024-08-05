package handlers

import (
	"net/http"

	"avito-backend-bootcamp/internal/model"
	pkgCtx "avito-backend-bootcamp/pkg/utils/ctx"

	"github.com/go-chi/render"
)

type dummyLoginResponse struct {
	Token string `json:"token"`
}

func HandleDummyLogin(authService AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userTypeRaw := r.URL.Query().Get(pkgCtx.KeyUserType)

		userType, err := model.ParseUserType(userTypeRaw)
		if err != nil {
			writeInternalError(r, w, err)
			return
		}

		token, err := authService.DummyLogin(r.Context(), userType)
		if err != nil {
			writeInternalError(r, w, err)
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, dummyLoginResponse{
			Token: token,
		})
	}
}
