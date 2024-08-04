package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type createFlatRequest struct {
	HouseID int64 `json:"house_id" validate:"required,gt=0"`
	Price   int64 `json:"price" validate:"required,gt=0"`
	Rooms   int32 `json:"rooms" validate:"required,gt=0"`
}

func HandleCreateFlat(validate *validator.Validate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createFlatRequest

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

		fmt.Println("Handle create flat")
		w.WriteHeader(http.StatusOK)
	}
}
