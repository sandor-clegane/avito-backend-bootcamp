package handlers

import (
	h "avito-backend-bootcamp/internal/http/handlers"
	"avito-backend-bootcamp/internal/model"
	"avito-backend-bootcamp/internal/service/auth"
	resp "avito-backend-bootcamp/pkg/utils/response"
	"avito-backend-bootcamp/pkg/utils/sl"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type AuthService interface {
	Register(ctx context.Context, email, password string, role model.UserType) (uuid.UUID, error)
}

type signupRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	UserType string `json:"user_type" validate:"required"`
}

type signupResponse struct {
	UserID string `json:"user_id"`
}

func New(log *slog.Logger, validate *validator.Validate, authService AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Setup logger
		const op = "handlers.HandleSignup"
		log := log.With(
			slog.String("op", op),
		)

		// Decode the request body into a SignupRequest struct
		var req signupRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Error("invalid json", sl.Err(err))
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

		// Parse the user type
		userTypeParsed, err := model.ParseUserType(req.UserType)
		if err != nil {
			log.Error("invalud user role in request", sl.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.NewError(err))
			return
		}

		// Register the user
		userID, err := authService.Register(r.Context(), req.Email, req.Password, userTypeParsed)
		if err != nil {
			if errors.Is(err, auth.ErrEmailAlreadyUsed) {
				log.Error("user with given email already exist", sl.Err(err))
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, resp.NewError(err))
				return
			}

			log.Error("failed to register user", sl.Err(err))
			h.WriteInternalError(r, w, err)
			return
		}

		// Respond with the user ID
		log.Info("user created succesfully")
		render.Status(r, http.StatusOK)
		render.JSON(w, r, signupResponse{
			UserID: userID.String(),
		})
	}
}
