package handlers

import (
	"avito-backend-bootcamp/internal/model"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type updateFlatRequest struct {
	ID     int64  `json:"id" validate:"required,gt=0"`
	Status string `json:"status" validate:"required"`
}

func HandleUpdateFlat(validate *validator.Validate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req updateFlatRequest

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

		_, err = model.ParseFlatStatus(req.Status)
		if err != nil {
			http.Error(w, fmt.Sprintf("Validation error: %s", err), http.StatusBadRequest)
			return
		}

		fmt.Println("Handle update flat")
		w.WriteHeader(http.StatusOK)
	}
}
