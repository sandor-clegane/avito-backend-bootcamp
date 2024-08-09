package handlers

import (
	"net/http"

	mwr "avito-backend-bootcamp/internal/http/middleware"

	"github.com/go-chi/render"
)

type ErrorResponse struct {
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
	Code      int    `json:"code"`
}

func WriteInternalError(r *http.Request, w http.ResponseWriter, err error) {
	render.Status(r, http.StatusInternalServerError)
	w.Header().Set("Retry-After", "30")
	render.JSON(w, r, ErrorResponse{
		Message:   err.Error(),
		RequestID: r.Context().Value(mwr.RequestIDKey).(string),
		Code:      http.StatusInternalServerError,
	})
}
