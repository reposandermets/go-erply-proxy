package handlers

import (
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message,omitempty"`
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "v1")
}
