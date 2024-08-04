package handlers

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
	Code      int    `json:"code"`
}

func writeError(w http.ResponseWriter, err errorResponse) {
	w.WriteHeader(err.Code)
	errorResponse, _ := json.Marshal(err)
	w.Write(errorResponse)
	w.Header().Set("Retry-After", "30")
}
