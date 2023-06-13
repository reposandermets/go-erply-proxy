package handlers

import (
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func Index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./api/index.html")
}

func Swagger(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./api/swagger.json")
}
