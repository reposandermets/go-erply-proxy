package main

import (
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/reposandermets/go-erply-proxy/internal/redis_utils"
	"github.com/reposandermets/go-erply-proxy/internal/routes"
)

// "github.com/go-redis/redis/v8"
// redis_utils "github.com/reposandermets/go-erply-proxy/internal/redis_utils"
// routes "github.com/reposandermets/go-erply-proxy/internal/routes"

// type Route struct {
// 	Name        string
// 	Method      string
// 	Pattern     string
// 	HandlerFunc http.HandlerFunc
// }

// items todo
// chainerply api
// unit tests
// auth middleware
// docker compose
// readme

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	defer redisClient.Close()
	go redis_utils.PeriodicallyClearCache(redisClient)
	router := routes.NewRouter(redisClient)
	log.Fatal(http.ListenAndServe(":8081", router))
}
