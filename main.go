package main

import (
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
	redis_utils "github.com/reposandermets/go-erply-proxy/internal/redis_utils"
	routes "github.com/reposandermets/go-erply-proxy/internal/routes"
)

// type BrandCreateRequestDescription struct {
// 	En string `json:"en,omitempty"`
// }

// type BrandCreateRequest struct {
// 	Description *BrandCreateRequestDescription `json:"description,omitempty"`
// 	Name        string                         `json:"name,omitempty"`
// }

// type BrandResponse struct {
// 	Added       int32                          `json:"added,omitempty"`
// 	Addedby     string                         `json:"addedby,omitempty"`
// 	Changed     int32                          `json:"changed,omitempty"`
// 	Changedby   string                         `json:"changedby,omitempty"`
// 	Description *BrandCreateRequestDescription `json:"description,omitempty"`
// 	Id          int32                          `json:"id,omitempty"`
// 	Name        string                         `json:"name,omitempty"`
// }

// type ErrorResponse struct {
// 	Message string `json:"message,omitempty"`
// }

// type Route struct {
// 	Name        string
// 	Method      string
// 	Pattern     string
// 	HandlerFunc http.HandlerFunc
// }

func main() {
	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Update with your Redis server address
		Password: "",               // Update with your Redis server password
		DB:       0,                // Update with the desired Redis database number
	})

	// Close the Redis client when main() exits
	defer redisClient.Close()

	go redis_utils.PeriodicallyClearCache(redisClient)

	router := routes.NewRouter(redisClient)

	log.Fatal(http.ListenAndServe(":8081", router))
}
