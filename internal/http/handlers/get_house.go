package handlers

import (
	"avito-backend-bootcamp/internal/model"
	pkgCtx "avito-backend-bootcamp/pkg/utils/ctx"
	resp "avito-backend-bootcamp/pkg/utils/response"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type getHouseResponse struct {
	Flats []*model.Flat
}

func HandleGetHouse(flatService FlatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract house ID from URL parameter
		houseIDStr := chi.URLParam(r, "id")
		houseID, err := strconv.ParseInt(houseIDStr, 10, 64)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.NewError(err))
			return
		}

		// Extract audience from the context
		userTypeStr := r.Context().Value(pkgCtx.KeyUserType).(string)

		// Retrieve flats associated with the house ID
		flatList, err := flatService.GetFlatListByHouseID(r.Context(), houseID, model.MustParseUserType(userTypeStr))
		if err != nil {
			writeInternalError(r, w, err)
			return
		}

		// Return the list of flats
		render.Status(r, http.StatusOK)
		render.JSON(w, r, getHouseResponse{
			Flats: flatList,
		})
	}
}
