package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	redis_utils "github.com/reposandermets/go-erply-proxy/internal/redis_utils"
)

func V1BrandGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	split := strings.Split(r.Header.Get("Authorization"), ":")
	if len(split) != 2 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"Unauthorized"}`))
		return
	}

	key1, key2 := redis_utils.GenerateUniqueKey(r)
	data, _ := redis_utils.GetFromCache(r.Context(), key1+"|"+key2)

	if data == "" {
		log.Println("Cache miss")
		// Cache miss, get data from the database
		// clientCode := split[0]
		// sessionKey := split[1]
		// data = GetBrandFromErplyAPI(r)

		data = `{"id":1,"name":"Brand 1","description":{"en":"Brand 1 description"}}`
		// Save data to cache
		err := redis_utils.SaveToCache(r.Context(), key1, key2, data)
		if err != nil {
			// Error occurred while saving to Redis
			log.Printf("Error saving to cache: %v\n", err)
		}
	} else {
		// Cache hit, return data
		log.Println("Cache hit")
	}
	fmt.Printf("reqLookupKey: %s\n", data)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(data))
}

func V1BrandPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	//a, _ := GetRedisClientFromContext(r.Context())

	// go ClearCache(r.Context(), , "/v1/brand")

}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}
