package scrapper

import (
	"encoding/json"
	"net/http"
	"os"
)

const (
	//common query parameters
	queryid           = "filteredProductsWithContext"
	anonymousId       = "25AFE5BE9BB9BC03DE89DBE170D80669"
	language          = "pt-PT"
	country           = "PT"
	channel           = "NIKE"
	localizedRangeStr = "%7BlowestPrice%7D%E2%80%94%7BhighestPrice%7D"
	//localizedRangeStr = "%7BlowestPrice%7D%20%E2%80%94%20%7BhighestPrice%7D"

	// UUIDs
	uuids_men   = "0f64ecc7-d624-4e91-b171-b83a03dd8550,16633190-45e5-4830-a068-232ac7aea82c"
	uuids_women = "16633190-45e5-4830-a068-232ac7aea82c,193af413-39b0-4d7e-ae34-558821381d3f,7baf216c-acc6-4452-9e07-39c2ca77ba32"
)

func NikeScrapper() string {
	var nikeURL = "https://api.nike.com/cic/browse/v2?queryid=" + queryid + "&anonymousId=" + anonymousId + "&uuids=" + uuids_men + "&language=" + language + "&country=" + country + "&channel=" + channel + "&localizedRangeStr=" + localizedRangeStr
	return nikeURL
}

type NikeResponse struct {
	Data struct {
		FilteredProductsWithContext struct {
			Products []struct {
				Title string `json:"title"`
				Price struct {
					FormattedCurrentPrice string `json:"formattedCurrentPrice"`
				} `json:"price"`
				URL string `json:"url"`
			} `json:"products"`
		} `json:"filteredProductsWithContext"`
	} `json:"data"`
}

type NikeProduct struct {
	ID        string `json:"id"`
	UnitPrice struct {
		FullRetail float64 `json:"fullRetail"`
		Current    float64 `json:"current"`
	} `json:"unitPrice"`
}

func fetchDataFromNikeAPI(url string) ([]ShoeInfo, error) {
	var Shoes []ShoeInfo

	// Make HTTP GET request to Nike API
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Parse JSON response
	var nikeResponse NikeResponse
	err = json.NewDecoder(response.Body).Decode(&nikeResponse)
	if err != nil {
		return nil, err
	}

	// Extract data from the response and populate ShoeInfo
	for _, product := range nikeResponse.Data.FilteredProductsWithContext.Products {
		shoe := ShoeInfo{
			NAME:  product.Title,
			PRICE: product.Price.FormattedCurrentPrice,
			LINK:  product.URL,
		}
		Shoes = append(Shoes, shoe)
	}

	return Shoes, nil
}

func fetchDataFromLocalJSON(filePath string, Shoes *[]ShoeInfo) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var response struct {
		Data struct {
			FilteredProductsWithContext struct {
				Products []struct {
					Title string `json:"title"`
					Price struct {
						FormattedCurrentPrice string `json:"formattedCurrentPrice"`
					} `json:"price"`
					URL string `json:"url"`
				} `json:"products"`
			} `json:"filteredProductsWithContext"`
		} `json:"data"`
	}

	err = json.Unmarshal(data, &response)
	if err != nil {
		return err
	}

	for _, product := range response.Data.FilteredProductsWithContext.Products {
		shoe := ShoeInfo{
			NAME:  product.Title,
			PRICE: product.Price.FormattedCurrentPrice,
			LINK:  product.URL,
		}
		*Shoes = append(*Shoes, shoe)
	}

	return nil
}

func ScrapeProductsx(useLocalJSON bool) ([]ShoeInfo, error) {
	var Shoes []ShoeInfo

	if useLocalJSON {
		// Fetch data from local JSON file
		err := fetchDataFromLocalJSON("content.json", &Shoes)
		if err != nil {
			return nil, err
		}
	} else {
		// Fetch data from Nike API
		nikeURL := NikeScrapper()
		fetchedShoes, err := fetchDataFromNikeAPI(nikeURL)
		if err != nil {
			return nil, err
		}
		Shoes = append(Shoes, fetchedShoes...)
	}

	return Shoes, nil
}
