package routes

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/reposandermets/go-erply-proxy/internal/handlers"

	redis_utils "github.com/reposandermets/go-erply-proxy/internal/redis_utils"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		handlers.Index,
	},
	Route{
		"V1BrandGet",
		strings.ToUpper("Get"),
		"/v1/brand",
		handlers.V1BrandGet,
	},
	Route{
		"V1BrandPost",
		strings.ToUpper("Post"),
		"/v1/brand",
		handlers.V1BrandPost,
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

func NewRouter(client *redis.Client) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		handler = redis_utils.WithRedisContext(handler, client)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}
