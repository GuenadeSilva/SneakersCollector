package main

import (
	"fmt"
	sc "sneakercollector/scrapper"
)

func main() {

	// Scrape from Adidas
	shoes1 := sc.ScrapeProducts("Adidas")

	// Scrape from New Balance
	shoes2 := sc.ScrapeProducts("New Balance")

	// Print the scraped data for Adidas
	fmt.Println("Adidas:")
	printShoeInfo(shoes1)

	// Print the scraped data for New balance
	fmt.Println("New Balance:")
	printShoeInfo(shoes2)

	useLocalJSON := true
	Shoes, err := sc.ScrapeProductsx(useLocalJSON)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print or process the scraped shoe data
	for _, shoe := range Shoes {
		fmt.Printf("Name: %s\nPrice: %s\nLink: %s\n\n", shoe.NAME, shoe.PRICE, shoe.LINK)
		fmt.Println("---------------------------")
	}

}

func printShoeInfo(shoes []sc.ShoeInfo) {
	for _, shoe := range shoes {
		fmt.Printf("Shoe Name: %s\n", shoe.NAME)
		fmt.Printf("Price: %s\n", shoe.PRICE)
		fmt.Printf("Link: %s\n", shoe.LINK)
		fmt.Println("---------------------------")
	}
}
