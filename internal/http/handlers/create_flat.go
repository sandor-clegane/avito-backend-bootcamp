package handlers

import (
	resp "avito-backend-bootcamp/pkg/utils/response"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type createFlatRequest struct {
	HouseID int64 `json:"house_id" validate:"required,gt=0"`
	Price   int64 `json:"price" validate:"required,gt=0"`
	Rooms   int64 `json:"rooms" validate:"required,gt=0"`
}

type createFlatResponse struct {
	ID      int64  `json:"id"`
	HouseID int64  `json:"house_id"`
	Price   int64  `json:"price"`
	Rooms   int64  `json:"rooms"`
	Status  string `json:"status"`
}

func HandleCreateFlat(validate *validator.Validate, flatService FlatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode the request body into a CreateFlatRequest struct
		var req createFlatRequest
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

		// Create the flat
		flat, err := flatService.CreateFlat(r.Context(), req.HouseID, req.Price, req.Rooms)
		if err != nil {
			writeInternalError(r, w, err)
			return
		}

		// Return the created flat details
		render.Status(r, http.StatusOK)
		render.JSON(w, r, createFlatResponse{
			ID:      flat.ID,
			HouseID: flat.HouseID,
			Price:   flat.Price,
			Status:  string(flat.Status),
			Rooms:   flat.Rooms,
		})
	}
}
