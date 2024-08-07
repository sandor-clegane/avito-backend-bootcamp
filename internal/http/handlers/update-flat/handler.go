package handlers

import (
	h "avito-backend-bootcamp/internal/http/handlers"
	"avito-backend-bootcamp/internal/model"
	resp "avito-backend-bootcamp/pkg/utils/response"
	"avito-backend-bootcamp/pkg/utils/sl"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type FlatService interface {
	UpdateFlat(ctx context.Context, ID int64, status model.FlatStatus) (*model.Flat, error)
}

type updateFlatRequest struct {
	ID     int64  `json:"id" validate:"required,gt=0"`
	Status string `json:"status" validate:"required"`
}

type updateFlatResponse struct {
	ID      int64  `json:"id"`
	HouseID int64  `json:"house_id"`
	Price   int64  `json:"price"`
	Rooms   int64  `json:"rooms"`
	Status  string `json:"status"`
}

func New(log *slog.Logger, validate *validator.Validate, flatService FlatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Setup logger
		const op = "handlers.HandleUpdateFlat"
		log := log.With(
			slog.String("op", op),
		)

		// Decode the request body into an UpdateFlatRequest struct
		var req updateFlatRequest
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

		// Parse the flat status
		status, err := model.ParseFlatStatus(req.Status)
		if err != nil {
			log.Error("invalid flat status in request", sl.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.NewError(err))
			return
		}

		// Update the flat
		flat, err := flatService.UpdateFlat(r.Context(), req.ID, status)
		if err != nil {
			log.Error("failed to update flat", sl.Err(err))
			h.WriteInternalError(r, w, err)
			return
		}

		// Return the updated flat details
		log.Info("flat updated succesfully")
		render.Status(r, http.StatusOK)
		render.JSON(w, r, updateFlatResponse{
			ID:      flat.ID,
			HouseID: flat.HouseID,
			Price:   flat.Price,
			Status:  string(flat.Status),
			Rooms:   flat.Rooms,
		})
	}
}
