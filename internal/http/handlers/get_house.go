package handlers

import (
	"fmt"
	"net/http"
)

func HandleGetHouse() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Handle get house")
		w.WriteHeader(http.StatusOK)
	}
}
