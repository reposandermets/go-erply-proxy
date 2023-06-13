package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/reposandermets/go-erply-proxy/internal/erply"
	"github.com/reposandermets/go-erply-proxy/internal/redis_utils"
)

// V1BrandGet handles the GET request for retrieving a list of brands.
func V1BrandGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	ctx := r.Context()

	redisUtil := ctx.Value("redisUtil").(redis_utils.RedisUtil)
	erplyClient := ctx.Value("erplyClient").(erply.ErplyAPI)

	sessionKey, _ := ctx.Value("ErplySessionKey").(string)
	clientCode, _ := ctx.Value("ErplyClientCode").(string)

	categoryKey, urlParamsKey := redisUtil.GenerateUniqueKey(r)
	data, _ := redisUtil.GetFromCache(ctx, categoryKey+"|"+urlParamsKey)

	if data == "" {
		log.Println("Cache miss", categoryKey, urlParamsKey)

		brands, err := erplyClient.GetBrands(ctx, sessionKey, clientCode)
		if err != nil {
			log.Printf("Error retrieving brands: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message":` + err.Error() + `}`))
			return
		}

		jsonData, err := json.Marshal(brands)
		if err != nil {
			log.Printf("Error marshaling JSON: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message":"Failed to marshal JSON"}`))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)

		var wg sync.WaitGroup
		wg.Add(1)
		go redisUtil.ManageSaveToCache(&wg, r, categoryKey, urlParamsKey, jsonData)
		wg.Wait()

		return
	}

	log.Println("Cache hit", categoryKey, urlParamsKey)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(data))
}
