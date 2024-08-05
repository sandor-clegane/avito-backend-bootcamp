package handlers

import (
	resp "avito-backend-bootcamp/pkg/utils/response"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type subscribeHouseRequest struct {
	Email string `json:"email" validate:"required,email"`
}

func HandleSubscribeHouse(validate *validator.Validate, subService SubscriptionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode the request body into a SubscribeHouseRequest struct
		var req subscribeHouseRequest
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

		// Extract house ID from URL parameter
		houseIDStr := chi.URLParam(r, "id")
		houseID, err := strconv.ParseInt(houseIDStr, 10, 64)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.NewError(err))
			return
		}

		// Create subscription
		err = subService.CreateSubscription(r.Context(), houseID, req.Email)
		if err != nil {
			writeInternalError(r, w, err)
			return
		}

		// Respond with success status
		render.Status(r, http.StatusOK)
	}
}
