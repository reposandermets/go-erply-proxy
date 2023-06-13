package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/reposandermets/go-erply-proxy/internal/erply"
	"github.com/reposandermets/go-erply-proxy/internal/redis_utils"
)

type BrandCreateRequest struct {
	Name string `json:"name"`
}

// V1BrandPost handles the POST request for creating a brand.
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

	erplyClient := ctx.Value("erplyClient").(erply.ErplyAPI)

	payload := map[string]string{
		"name": brand.Name,
	}

	res, err := erplyClient.SaveBrand(ctx, sessionKey, clientCode, payload)
	if err != nil {
		log.Println("Error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":` + err.Error() + `}`))
		return
	}

	responseJSON, err := json.Marshal(res)
	if err != nil {
		log.Println("Error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"Internal server error"}`))
		return
	}

	var wg sync.WaitGroup
	redisUtil := ctx.Value("redisUtil").(redis_utils.RedisUtil)
	wg.Add(1)

	go redisUtil.ManageClearCache(&wg, r)

	w.WriteHeader(http.StatusCreated)
	w.Write(responseJSON)

	wg.Wait()
}
