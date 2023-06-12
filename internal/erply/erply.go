package erply

import (
	"context"
	"time"

	"github.com/erply/api-go-wrapper/pkg/api"
	"github.com/erply/api-go-wrapper/pkg/api/products"
)

type BrandCreateRequest struct {
	Name string `json:"name"`
}

type ErplyAPI interface {
	SaveBrand(ctx context.Context, sessionKey string, clientCode string, payload map[string]string) (products.SaveBrandResult, error)
	GetBrands(ctx context.Context, sessionKey string, clientCode string) ([]products.ProductBrand, error)
}

type ErplyClient struct {
}

func NewErplyAPI() ErplyAPI {
	return &ErplyClient{}
}

func (c *ErplyClient) SaveBrand(ctx context.Context, sessionKey string, clientCode string, payload map[string]string) (products.SaveBrandResult, error) {
	cli, err := api.NewClient(sessionKey, clientCode, nil)
	if err != nil {
		return products.SaveBrandResult{}, err
	}
	cli.SendParametersInRequestBody()

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	return cli.ProductManager.SaveBrand(ctx, payload)
}

func (c *ErplyClient) GetBrands(ctx context.Context, sessionKey string, clientCode string) ([]products.ProductBrand, error) {
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
