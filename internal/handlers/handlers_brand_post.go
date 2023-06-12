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

type BrandCreateRequest struct {
	Name string `json:"name"`
}

func SaveBrandToErplyAPI(sessionKey string, clientCode string, payload BrandCreateRequest) (result products.SaveBrandResult, err error) {

	cli, err := api.NewClient(sessionKey, clientCode, nil)
	if err != nil {
		return result, err
	}
	cli.SendParametersInRequestBody()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	data := map[string]string{
		"name": payload.Name,
	}

	return cli.ProductManager.SaveBrand(ctx, data)
}

func V1BrandPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var brand BrandCreateRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&brand); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message":"Invalid request body"}`))
		return
	}

	defer r.Body.Close()

	if brand.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message":"Brand name is required"}`))
		return
	}

	// TODO needs middleware to check if the user is authorized
	split := strings.Split(r.Header.Get("Authorization"), ":")
	if len(split) != 2 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"Unauthorized"}`))
		return
	}

	sessionKey := split[0]
	clientCode := split[1]

	res, err := SaveBrandToErplyAPI(sessionKey, clientCode, brand)
	if err != nil {
		fmt.Println("Error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"Internal server error"}`))
		return
	}

	responseJSON, _ := json.Marshal(res)

	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)

	var wg sync.WaitGroup
	wg.Add(1)

	// Clear cache for this request category: /v1/brand
	go func(req *http.Request) {
		defer wg.Done()
		categoryKey, _ := redis_utils.GenerateUniqueKey(req)
		err := redis_utils.ClearCache(r.Context(), categoryKey)
		if err != nil {
			// Handle error while saving to cache
			log.Printf("Error saving to cache: %v\n", err)
		}
	}(r)

	// Wait for the cache saving goroutine to complete
	wg.Wait()
}
