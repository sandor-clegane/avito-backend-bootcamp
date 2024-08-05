package handlers

import (
	"avito-backend-bootcamp/internal/model"
	resp "avito-backend-bootcamp/pkg/utils/response"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type signupRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	UserType string `json:"user_type" validate:"required"`
}

type signupResponse struct {
	UserID string `json:"user_id"`
}

func HandleSignup(validate *validator.Validate, authService AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode the request body into a SignupRequest struct
		var req signupRequest
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

		// Parse the user type
		userTypeParsed, err := model.ParseUserType(req.UserType)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.NewError(err))
			return
		}

		// Register the user
		userID, err := authService.Register(r.Context(), req.Email, req.Password, userTypeParsed)
		if err != nil {
			writeInternalError(r, w, err)
			return
		}

		// Respond with the user ID
		render.Status(r, http.StatusOK)
		render.JSON(w, r, signupResponse{
			UserID: userID.String(),
		})
	}
}
