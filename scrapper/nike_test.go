package scrapper

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestFetchDataFromNikeAPI(t *testing.T) {
	// Mock Nike API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Serve mock JSON response
		mockResponse := NikeResponse{
			Data: struct {
				FilteredProductsWithContext struct {
					Products []struct {
						Title string `json:"title"`
						Price struct {
							FormattedCurrentPrice string `json:"formattedCurrentPrice"`
						} `json:"price"`
						URL string `json:"url"`
					} `json:"products"`
				} `json:"filteredProductsWithContext"`
			}{},
		}
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// Fetch data using the mock server
	shoes, err := fetchDataFromNikeAPI(server.URL)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Check if the result is not empty
	if len(shoes) == 0 {
		t.Error("fetchDataFromNikeAPI should return non-empty shoe list")
	}
}

func TestFetchDataFromLocalJSON(t *testing.T) {
	// Create a temporary mock JSON file
	mockJSONFile := "mock.json"
	defer os.Remove(mockJSONFile)

	// Write mock JSON data to the file
	mockData := []byte(`{
		"data": {
			"filteredProductsWithContext": {
				"products": [
					{
						"title": "Mock Shoe",
						"price": {
							"formattedCurrentPrice": "$100"
						},
						"url": "https://example.com/mock"
					}
				]
			}
		}
	}`)
	err := os.WriteFile(mockJSONFile, mockData, 0644)
	if err != nil {
		t.Fatalf("Failed to create mock JSON file: %v", err)
	}

	// Fetch data from the mock JSON file
	var shoes []ShoeInfo
	err = fetchDataFromLocalJSON(mockJSONFile, &shoes)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Check if the result is not empty
	if len(shoes) == 0 {
		t.Error("fetchDataFromLocalJSON should return non-empty shoe list")
	}
}

func TestScrapeProductsx(t *testing.T) {
	// Test fetching data from Nike API
	shoes, err := ScrapeProductsx(false)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if len(shoes) == 0 {
		t.Error("ScrapeProductsx should return non-empty shoe list")
	}

	// Test fetching data from local JSON
	shoes, err = ScrapeProductsx(true)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if len(shoes) == 0 {
		t.Error("ScrapeProductsx should return non-empty shoe list")
	}
}
