package handlers

import (
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "v1")
}
