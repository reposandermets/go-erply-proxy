// handlers_test.go
package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestV1BrandGet(t *testing.T) {
	// Create a request
	req, err := http.NewRequest("GET", "/v1/brand", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	handler := http.HandlerFunc(V1BrandGet)
	handler.ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("expected status code %v, got %v", http.StatusOK, rr.Code)
	}

	// Check the response body
	expectedBody := `{"id":1,"name":"Brand 1","description":{"en":"Brand 1 description"}}`
	if rr.Body.String() != expectedBody {
		t.Errorf("expected response body %v, got %v", expectedBody, rr.Body.String())
	}
}

// Add more test functions for other handlers if needed
