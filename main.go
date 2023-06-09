package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type BrandCreateRequestDescription struct {
	En string `json:"en,omitempty"`
}

type BrandCreateRequest struct {
	Description *BrandCreateRequestDescription `json:"description,omitempty"`

	Name string `json:"name,omitempty"`
}

type BrandResponse struct {
	Added int32 `json:"added,omitempty"`

	Addedby string `json:"addedby,omitempty"`

	Changed int32 `json:"changed,omitempty"`

	Changedby string `json:"changedby,omitempty"`

	Description *BrandCreateRequestDescription `json:"description,omitempty"`

	Id int32 `json:"id,omitempty"`

	Name string `json:"name,omitempty"`
}

type ErrorResponse struct {
	Message string `json:"message,omitempty"`
}

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func V1BrandGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func V1BrandPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},

	Route{
		"V1BrandGet",
		strings.ToUpper("Get"),
		"/v1/brand",
		V1BrandGet,
	},

	Route{
		"V1BrandPost",
		strings.ToUpper("Post"),
		"/v1/brand",
		V1BrandPost,
	},
}

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s %s %s %s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}

func main() {
	log.Printf("Server started")

	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8081", router))
}
