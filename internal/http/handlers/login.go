package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"avito-backend-bootcamp/internal/service/auth"
	resp "avito-backend-bootcamp/pkg/utils/response"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type loginRequest struct {
	ID       string `json:"id" validate:"required,uuid"`
	Password string `json:"password" validate:"required"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func HandleLogin(validate *validator.Validate, authService AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode the request body into a LoginRequest struct
		var req loginRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.NewError(err))
			return
		}

		// Validate the request data
		err = validate.Struct(req)
		if err != nil {
			errors := err.(validator.ValidationErrors)
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.NewError(fmt.Errorf("Validation error: %s", errors)))
			return
		}

		// Authenticate the user
		token, err := authService.Login(r.Context(), uuid.MustParse(req.ID), req.Password)
		if err != nil {
			if errors.Is(err, auth.ErrUserNotFound) {
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, resp.NewError(err))
				return
			}
			writeInternalError(r, w, err)
			return
		}

		// Respond with the authentication token
		render.Status(r, http.StatusOK)
		render.JSON(w, r, loginResponse{
			Token: token,
		})
	}
}
