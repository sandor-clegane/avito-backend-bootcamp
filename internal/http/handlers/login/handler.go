package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	h "avito-backend-bootcamp/internal/http/handlers"
	"avito-backend-bootcamp/internal/service/auth"
	resp "avito-backend-bootcamp/pkg/utils/response"
	"avito-backend-bootcamp/pkg/utils/sl"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type AuthService interface {
	Login(ctx context.Context, ID uuid.UUID, password string) (string, error)
}

type loginRequest struct {
	ID       string `json:"id" validate:"required,uuid"`
	Password string `json:"password" validate:"required"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func New(log *slog.Logger, validate *validator.Validate, authService AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Setup logger
		const op = "handlers.HandleLogin"
		log := log.With(
			slog.String("op", op),
		)

		// Decode the request body into a LoginRequest struct
		var req loginRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Error("invalid request json", sl.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.NewError(err))
			return
		}

		// Validate the request data
		err = validate.Struct(req)
		if err != nil {
			log.Error("input validation failed", sl.Err(err))
			errors := err.(validator.ValidationErrors)
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.NewError(fmt.Errorf("Validation error: %s", errors)))
			return
		}

		// Authenticate the user
		token, err := authService.Login(r.Context(), uuid.MustParse(req.ID), req.Password)
		if err != nil {
			if errors.Is(err, auth.ErrUserNotFound) {
				log.Error("user with given credentials not exist", sl.Err(err))
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, resp.NewError(err))
				return
			}

			if errors.Is(err, auth.ErrInvalidCredentials) {
				log.Error("user entered invalid credentials", sl.Err(err))
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, resp.NewError(err))
				return
			}

			log.Error("login failed", sl.Err(err))
			h.WriteInternalError(r, w, err)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "Authorization",
			Value: token,
		})

		// Respond with the authentication token
		log.Info("login successfully")
		render.Status(r, http.StatusOK)
		render.JSON(w, r, loginResponse{
			Token: token,
		})
	}
}
