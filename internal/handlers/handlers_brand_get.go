package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/erply/api-go-wrapper/pkg/api"
	"github.com/erply/api-go-wrapper/pkg/api/products"
	"github.com/reposandermets/go-erply-proxy/internal/redis_utils"
)

func GetBrandsFromErplyAPI(ctx context.Context, sessionKey string, clientCode string) ([]products.ProductBrand, error) {
	cli, err := api.NewClient(sessionKey, clientCode, nil)
	if err != nil {
		return nil, err
	}

	cli.SendParametersInRequestBody()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	result, err := cli.ProductManager.GetBrands(ctx, nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func V1BrandGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	ctx := r.Context()

	sessionKey, _ := ctx.Value("ErplySessionKey").(string)
	clientCode, _ := ctx.Value("ErplyClientCode").(string)

	categoryKey, urlParamsKey := redis_utils.GenerateUniqueKey(r)
	data, _ := redis_utils.GetFromCache(ctx, categoryKey+"|"+urlParamsKey)

	if data == "" {
		log.Println("Cache miss")

		brands, err := GetBrandsFromErplyAPI(ctx, sessionKey, clientCode)
		if err != nil {
			log.Printf("Error retrieving brands: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message":"Failed to retrieve brands"}`))
			return
		}

		jsonData, err := json.Marshal(brands)
		if err != nil {
			log.Printf("Error marshaling JSON: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message":"Failed to marshal JSON"}`))
			return
		}

		log.Println(string(jsonData))
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)

		var wg sync.WaitGroup
		wg.Add(1)

		// Save data to cache asynchronously
		go func() {
			defer wg.Done()
			err := redis_utils.SaveToCache(ctx, categoryKey, urlParamsKey, string(jsonData))
			if err != nil {
				log.Printf("Error saving to cache: %v\n", err)
				// Handle the error accordingly, e.g., return an error response or log it for later analysis
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
