package main

import (
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/reposandermets/go-erply-proxy/internal/erply"
	"github.com/reposandermets/go-erply-proxy/internal/redis_utils"
	"github.com/reposandermets/go-erply-proxy/internal/routes"
)

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
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
