package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/reposandermets/go-erply-proxy/internal/erply"
	"github.com/reposandermets/go-erply-proxy/internal/redis_utils"
	"github.com/reposandermets/go-erply-proxy/internal/routes"
)

func main() {
	dockerized := os.Getenv("DOCKERIZED")

	print("Dockerized: " + dockerized + "\n")

	redisUrl := "0.0.0.0:6379"
	if dockerized == "true" {
		redisUrl = "redis:6379"
	}
	print("redisUrl: " + redisUrl + "\n")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: "",
		DB:       0,
	})

	defer redisClient.Close()

	erplyClient := erply.NewErplyAPI()
	redisUtil := redis_utils.NewRedisUtil(redisClient)

	go redisUtil.PeriodicallyClearCache()

	router := routes.NewRouter(redisUtil, erplyClient)
	log.Fatal(http.ListenAndServe(":8081", router))
}
