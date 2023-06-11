package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/erply/api-go-wrapper/pkg/api"
)

type BrandCreateRequest struct {
	Name string `json:"name,omitempty"`
}

func SaveBrandToErplyAPI() {

	sessionKey := "c2bb1db09a2ce295c67ae0c6b5ffee9c3a0327a319ba"
	clientCode := "104791"

	cli, err := api.NewClient(sessionKey, clientCode, nil)
	if err != nil {
		panic(err)
	}
	cli.SendParametersInRequestBody()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	payload := map[string]string{
		"name": "some brand 44",
	}

	res, _ := cli.ProductManager.SaveBrand(ctx, payload)
	fmt.Printf("Brand: %+v\n", res)

	// if err != nil {
	// 	panic(err)
	// }

	// brands, err := cli.ProductManager.GetBrands(ctx, nil)
	// // print the result
	// if err != nil {
	// 	panic(err)
	// }
	// // println(brands)

	// fmt.Printf("Brand: %+v\n", brands)

}

func V1BrandPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	//a, _ := GetRedisClientFromContext(r.Context())

	// go ClearCache(r.Context(), , "/v1/brand")

}
