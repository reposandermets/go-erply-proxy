package routes

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/reposandermets/go-erply-proxy/internal/erply"
	"github.com/reposandermets/go-erply-proxy/internal/handlers"
	"github.com/reposandermets/go-erply-proxy/internal/redis_utils"
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
		http.MethodGet,
		"/v1/brand",
		handlers.V1BrandGet,
	},
	Route{
		"V1BrandPost",
		http.MethodPost,
		"/v1/brand",
		handlers.V1BrandPost,
	},
}

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		inner.ServeHTTP(w, r)
		log.Printf(
			"%s %s %s %s - %s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
			r.RemoteAddr,
		)
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			// Skip authentication for the "Index" route
			next.ServeHTTP(w, r)
			return
		}

		split := strings.Split(r.Header.Get("Authorization"), ":")
		if len(split) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"message":"Unauthorized"}`))
			return
		}

		// Add custom context keys for sessionKey and clientCode
		ctx := context.WithValue(r.Context(), "ErplySessionKey", split[0])
		ctx = context.WithValue(ctx, "ErplyClientCode", split[1])

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// WithRedisContext creates a new context with the Redis client.
func WithRedisContext(handler http.Handler, redisUtil redis_utils.RedisUtil) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "redisUtil", redisUtil)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}

// WithRedisContext creates a new context with the Redis client.
func WithErplyAPIContext(handler http.Handler, erplyClient erply.ErplyAPI) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a new context with the Redis client
		ctx := context.WithValue(r.Context(), "erplyClient", erplyClient)

		// Serve the request with the new context
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}

func NewRouter(redisUtil redis_utils.RedisUtil, erplyClient erply.ErplyAPI) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		if route.Name != "Index" {
			handler = AuthMiddleware(handler)
			handler = WithRedisContext(handler, redisUtil)
			handler = WithErplyAPIContext(handler, erplyClient)
		}

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}
