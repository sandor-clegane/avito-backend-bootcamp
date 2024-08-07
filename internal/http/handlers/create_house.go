package handlers

import (
	dbUtil "avito-backend-bootcamp/pkg/utils/db"
	resp "avito-backend-bootcamp/pkg/utils/response"
	"avito-backend-bootcamp/pkg/utils/sl"
	"log/slog"

	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type createHouseRequest struct {
	Address   string `json:"address" validate:"required"`
	Year      int64  `json:"year" validate:"required,gt=0"`
	Developer string `json:"developer"`
}

type createHouseResponse struct {
	ID        int64     `json:"id"`
	Address   string    `json:"address"`
	Year      int64     `json:"year"`
	Developer *string   `json:"developer,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func HandleCreateHouse(log *slog.Logger, validate *validator.Validate, houseService HouseService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Setup logger
		const op = "handlers.HandleCreateHouse"
		log := log.With(
			slog.String("op", op),
		)

		// Decode the request body into a CreateHouseRequest struct
		var req createHouseRequest
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

		// Create the house
		house, err := houseService.CreateHouse(r.Context(), req.Address, req.Developer, req.Year)
		if err != nil {
			log.Error("create house failed", sl.Err(err))
			writeInternalError(r, w, err)
			return
		}

		// Return the created house details
		log.Info("house created")
		render.Status(r, http.StatusOK)
		render.JSON(w, r, createHouseResponse{
			ID:        house.ID,
			Address:   house.Address,
			Year:      house.YearOfConstruction,
			Developer: dbUtil.FromNullString(house.Developer),
			CreatedAt: house.CreatedAt,
			UpdatedAt: house.UpdatedAt,
		})
	}
}
