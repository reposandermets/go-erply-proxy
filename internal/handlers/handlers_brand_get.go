package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/erply/api-go-wrapper/pkg/api"
	"github.com/erply/api-go-wrapper/pkg/api/products"
	"github.com/reposandermets/go-erply-proxy/internal/redis_utils"
)

func GetBrandsFromErplyAPI(sessionKey string, clientCode string) ([]products.ProductBrand, error) {

	cli, err := api.NewClient(sessionKey, clientCode, nil)
	if err != nil {
		return nil, err
	}

	cli.SendParametersInRequestBody()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	brands, err := cli.ProductManager.GetBrands(ctx, nil)
	if err != nil {
		return nil, err
	}

	return brands, nil
}

func V1BrandGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// TODO needs middleware to check if the user is authorized
	split := strings.Split(r.Header.Get("Authorization"), ":")
	if len(split) != 2 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"Unauthorized"}`))
		return
	}

	categoryKey, urlParamsKey := redis_utils.GenerateUniqueKey(r)
	data, _ := redis_utils.GetFromCache(r.Context(), categoryKey+"|"+urlParamsKey)

	if data == "" {
		log.Println("Cache miss")
		sessionKey := split[0]
		clientCode := split[1]
		brands, err := GetBrandsFromErplyAPI(sessionKey, clientCode)
		if err != nil {
			fmt.Println("Error:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		jsonData, err := json.Marshal(brands)
		if err != nil {
			fmt.Println("Error:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Println(string(jsonData))
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)

		// Create a WaitGroup to synchronize goroutines
		var wg sync.WaitGroup
		wg.Add(1)

		// Save data to cache asynchronously
		go func() {
			defer wg.Done()
			err := redis_utils.SaveToCache(r.Context(), categoryKey, urlParamsKey, string(jsonData))
			if err != nil {
				// Handle error while saving to cache
				log.Printf("Error saving to cache: %v\n", err)
			}
		}()

		// Wait for the cache saving goroutine to complete
		wg.Wait()
		return
	}

	// Cache hit, return data
	log.Println("Cache hit")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(data))
}
