package handlers

import (
	h "avito-backend-bootcamp/internal/http/handlers"
	sub "avito-backend-bootcamp/internal/service/subscription"
	resp "avito-backend-bootcamp/pkg/utils/response"
	"avito-backend-bootcamp/pkg/utils/sl"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type SubscriptionService interface {
	CreateSubscription(ctx context.Context, houseID int64, email string) error
}

type subscribeHouseRequest struct {
	Email string `json:"email" validate:"required,email"`
}

func New(log *slog.Logger, validate *validator.Validate, subService SubscriptionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Setup logger
		const op = "handlers.HandleSubscribeHouse"
		log := log.With(
			slog.String("op", op),
		)

		// Decode the request body into a SubscribeHouseRequest struct
		var req subscribeHouseRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Error("invalid input json", sl.Err(err))
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

		// Extract house ID from URL parameter
		houseIDStr := chi.URLParam(r, "id")
		houseID, err := strconv.ParseInt(houseIDStr, 10, 64)
		if err != nil {
			log.Error("failed to get house id from url", sl.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.NewError(err))
			return
		}

		// Create subscription
		err = subService.CreateSubscription(r.Context(), houseID, req.Email)
		if err != nil {
			log.Error("failed to create subscription", sl.Err(err))
			if errors.Is(err, sub.ErrInvalidSubscription) {
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, resp.NewError(err))
				return
			}
			h.WriteInternalError(r, w, err)
			return
		}

		// Respond with success status
		log.Info("subscription success")
		render.Status(r, http.StatusOK)
	}
}
