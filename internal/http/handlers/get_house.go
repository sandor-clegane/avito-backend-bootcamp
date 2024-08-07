package handlers

import (
	"avito-backend-bootcamp/internal/model"
	pkgCtx "avito-backend-bootcamp/pkg/utils/ctx"
	resp "avito-backend-bootcamp/pkg/utils/response"
	"avito-backend-bootcamp/pkg/utils/sl"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type getHouseResponse struct {
	Flats []*model.Flat
}

func HandleGetHouse(log *slog.Logger, flatService FlatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Setup logger
		const op = "handlers.HandleGetHouse"
		log := log.With(
			slog.String("op", op),
		)

		// Extract house ID from URL parameter
		houseIDStr := chi.URLParam(r, "id")
		houseID, err := strconv.ParseInt(houseIDStr, 10, 64)
		if err != nil {
			log.Error("param parsing failed", sl.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.NewError(err))
			return
		}

		// Extract audience from the context
		userTypeStr := r.Context().Value(pkgCtx.KeyUserType).(string)
		log.Info("audience extracted from request")

		// Retrieve flats associated with the house ID
		flatList, err := flatService.GetFlatListByHouseID(r.Context(), houseID, model.MustParseUserType(userTypeStr))
		if err != nil {
			log.Error("failed to get list of flats for house", sl.Err(err))
			writeInternalError(r, w, err)
			return
		}

		// Return the list of flats
		log.Info("successfully get list of flats")
		render.Status(r, http.StatusOK)
		render.JSON(w, r, getHouseResponse{
			Flats: flatList,
		})
	}
}
