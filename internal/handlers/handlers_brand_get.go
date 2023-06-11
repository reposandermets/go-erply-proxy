package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/erply/api-go-wrapper/pkg/api"
	redis_utils "github.com/reposandermets/go-erply-proxy/internal/redis_utils"
)

type BrandResponse struct {
	Added     int32  `json:"added,omitempty"`
	Addedby   string `json:"addedby,omitempty"`
	Changed   int32  `json:"changed,omitempty"`
	Changedby string `json:"changedby,omitempty"`
	Id        int32  `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
}

func GetBrandsFromErplyAPI() {

	sessionKey := "c2bb1db09a2ce295c67ae0c6b5ffee9c3a0327a319ba"
	clientCode := "104791"

	cli, err := api.NewClient(sessionKey, clientCode, nil)
	if err != nil {
		panic(err)
	}

	//configure the client to send the data payload in the request body instead of the query parameters.
	//Using the request body eliminates the query size limitations imposed by the maximum URL length
	cli.SendParametersInRequestBody()

	//init context to control the request flow
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	brands, err := cli.ProductManager.GetBrands(ctx, nil)
	// print the result
	if err != nil {
		panic(err)
	}
	// println(brands)

	fmt.Printf("Brand: %+v\n", brands)

}

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
