package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/erply/api-go-wrapper/pkg/api/products"
	"github.com/reposandermets/go-erply-proxy/internal/handlers"
)

type MockErplyAPI struct{}

func (api *MockErplyAPI) SaveBrand(ctx context.Context, sessionKey string, clientCode string, payload map[string]string) (products.SaveBrandResult, error) {

	return products.SaveBrandResult{
		BrandID: 1,
	}, nil
}

func (api *MockErplyAPI) GetBrands(ctx context.Context, sessionKey string, clientCode string, filters map[string]string) ([]products.ProductBrand, error) {
	// Mock the GetBrands method implementation here
	// Return a sample array of brands for testing
	return []products.ProductBrand{
		{ID: 1, Name: "Brand 1"},
		{ID: 2, Name: "Brand 2"},
	}, nil
}

type MockRedisUtil struct{}

func (ru *MockRedisUtil) GenerateUniqueKey(r *http.Request) (string, string) {

	return "categoryKey", "urlParamsKey"
}

func (ru *MockRedisUtil) GetFromCache(ctx context.Context, key string) (string, error) {

	return "", nil
}

func (ru *MockRedisUtil) SaveToCache(ctx context.Context, key1, key2, data string) error {

	return nil
}

func (ru *MockRedisUtil) ClearCache(ctx context.Context, categoryKey string) error {

	return nil
}

func (ru *MockRedisUtil) FlushRedis(ctx context.Context) error {

	return nil
}

func (ru *MockRedisUtil) PeriodicallyClearCache() {

}

func (ru *MockRedisUtil) ManageClearCache(wg *sync.WaitGroup, r *http.Request) {

	wg.Done()
}

func (ru *MockRedisUtil) ManageSaveToCache(wg *sync.WaitGroup, r *http.Request, categoryKey string, urlParamsKey string, jsonData []byte) {

	wg.Done()
}

func TestV1BrandGet200(t *testing.T) {
	// Create a mock request
	req, err := http.NewRequest("GET", "/v1/brands", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock response recorder
	rr := httptest.NewRecorder()

	// Create a mock context with the necessary values
	ctx := context.WithValue(req.Context(), "ErplySessionKey", "mockSessionKey")
	ctx = context.WithValue(ctx, "ErplyClientCode", "mockClientCode")
	ctx = context.WithValue(ctx, "erplyClient", &MockErplyAPI{})
	ctx = context.WithValue(ctx, "redisUtil", &MockRedisUtil{})

	// Call the handler function
	handlers.V1BrandGet(rr, req.WithContext(ctx))

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	// Check the response body
	expectedBody := `[{"brandID":1,"name":"Brand 1","added":0,"lastModified":0},{"brandID":2,"name":"Brand 2","added":0,"lastModified":0}]`
	if rr.Body.String() != expectedBody {
		t.Errorf("Handler returned unexpected body: got %v, want %v", rr.Body.String(), expectedBody)
	}
}

func TestV1BrandPost200(t *testing.T) {
	// Create a sample payload
	payload := map[string]string{
		"name": "Brand 1",
	}
	payloadJSON, _ := json.Marshal(payload)

	// Create a mock request with the payload
	req, err := http.NewRequest("POST", "/v1/brands", bytes.NewReader(payloadJSON))
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock response recorder
	rr := httptest.NewRecorder()

	// Create a mock context with the necessary values
	ctx := context.WithValue(req.Context(), "ErplySessionKey", "mockSessionKey")
	ctx = context.WithValue(ctx, "ErplyClientCode", "mockClientCode")
	ctx = context.WithValue(ctx, "erplyClient", &MockErplyAPI{})
	ctx = context.WithValue(ctx, "redisUtil", &MockRedisUtil{})

	// Call the handler function
	handlers.V1BrandPost(rr, req.WithContext(ctx))

	// Check the response status code
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	// Check the response body
	expectedBody := `{"brandID":1}`
	if rr.Body.String() != expectedBody {
		t.Errorf("Handler returned unexpected body: got %v, want %v", rr.Body.String(), expectedBody)
	}
}
