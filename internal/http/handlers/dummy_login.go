package handlers

import (
	"fmt"
	"net/http"
)

func HandleDummyLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Handle dummy login")
		w.WriteHeader(http.StatusOK)
	}
}
