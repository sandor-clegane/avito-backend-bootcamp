package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type signupRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	UserType string `json:"user_type" validate:"required"`
}

func HandleSignup(validate *validator.Validate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req signupRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		err = validate.Struct(req)
		if err != nil {
			errors := err.(validator.ValidationErrors)
			http.Error(w, fmt.Sprintf("Validation error: %s", errors), http.StatusBadRequest)
			return
		}

		fmt.Println("Handle signup")
		w.WriteHeader(http.StatusOK)
	}
}
