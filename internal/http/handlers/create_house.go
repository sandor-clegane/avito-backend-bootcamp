package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type createHouseRequest struct {
	Address   string `json:"address" validate:"required"`
	Year      int64  `json:"year" validate:"required,gt=0"`
	Developer int32  `json:"developer"`
}

func HandleCreateHouse(validate *validator.Validate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createHouseRequest

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

		fmt.Println("Handle create house")
		w.WriteHeader(http.StatusOK)
	}
}
