package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/erply/api-go-wrapper/pkg/api"
	"github.com/erply/api-go-wrapper/pkg/api/products"
	"github.com/reposandermets/go-erply-proxy/internal/redis_utils"
)

type BrandCreateRequest struct {
	Name string `json:"name"`
}

func SaveBrandToErplyAPI(ctx context.Context, sessionKey string, clientCode string, payload BrandCreateRequest) (products.SaveBrandResult, error) {
	cli, err := api.NewClient(sessionKey, clientCode, nil)
	if err != nil {
		return products.SaveBrandResult{}, err
	}
	cli.SendParametersInRequestBody()

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	data := map[string]string{
		"name": payload.Name,
	}

	return cli.ProductManager.SaveBrand(ctx, data)
}

func V1BrandPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	ctx := r.Context()

	var brand BrandCreateRequest

	err := json.NewDecoder(r.Body).Decode(&brand)
	if err != nil {
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

	sessionKey, _ := ctx.Value("ErplySessionKey").(string)
	clientCode, _ := ctx.Value("ErplyClientCode").(string)

	res, err := SaveBrandToErplyAPI(ctx, sessionKey, clientCode, brand)
	if err != nil {
		log.Println("Error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"Internal server error"}`))
		return
	}

	responseJSON, err := json.Marshal(res)
	if err != nil {
		log.Println("Error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"Internal server error"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)

	var wg sync.WaitGroup
	wg.Add(1)

	// Clear cache for this request category: /v1/brand
	go func() {
		defer wg.Done()
		categoryKey, _ := redis_utils.GenerateUniqueKey(r)
		err := redis_utils.ClearCache(r.Context(), categoryKey)
		if err != nil {
			log.Printf("Error clearing cache: %v\n", err)
		}
	}()

	// Wait for the cache clearing goroutine to complete
	wg.Wait()
}
