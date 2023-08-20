package scrapper

import (
	"testing"
)

func TestScrapeProducts_Adidas(t *testing.T) {
	// Call the ScrapeProducts function for Adidas
	shoes := ScrapeProducts("Adidas")

	// Check if the result is not empty
	if len(shoes) == 0 {
		t.Error("ScrapeProducts should return non-empty shoe list")
	}

	// Perform more specific checks for Adidas
	// ...

}

func TestScrapeProducts_NewBalance(t *testing.T) {
	// Call the ScrapeProducts function for New Balance
	shoes := ScrapeProducts("New Balance")

	// Check if the result is not empty
	if len(shoes) == 0 {
		t.Error("ScrapeProducts should return non-empty shoe list")
	}

	// Perform more specific checks for New Balance
	// ...

}

func TestScrapeProducts_InvalidSelector(t *testing.T) {
	// Call the ScrapeProducts function with an invalid selector key
	shoes := ScrapeProducts("InvalidSelectorKey")

	// Check if the result is empty
	if len(shoes) != 0 {
		t.Error("ScrapeProducts should return an empty shoe list for invalid selector key")
	}

	// Perform more specific checks if needed
	// ...

}

// Add more tests for the extractName, extractPrice, and extractLink functions
// ...

// Add more tests for other functions and scenarios as needed
// ...

func TestWebsiteSelectorsMap(t *testing.T) {
	// Test the correctness of the websiteSelectorsMap
	// ...

	// Example test for checking if Adidas selector exists
	if _, exists := websiteSelectorsMap["Adidas"]; !exists {
		t.Error("Expected Adidas selector to exist")
	}

	// Example test for checking if New Balance selector exists
	if _, exists := websiteSelectorsMap["New Balance"]; !exists {
		t.Error("Expected New Balance selector to exist")
	}

	// Example test for checking if an invalid selector key doesn't exist
	if _, exists := websiteSelectorsMap["InvalidSelectorKey"]; exists {
		t.Error("Expected InvalidSelectorKey to not exist")
	}

	// Perform more tests and checks as needed
	// ...
}
