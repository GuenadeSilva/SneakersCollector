package database

import (
	"database/sql"
	"log"

	"sneakercollector/scrapper"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func RefreshScrapedData() {
	db, err := sql.Open("postgres", "postgres://user:password@localhost/dbname?sslmode=disable")
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		return
	}
	defer db.Close()

	// Create the table if it doesn't exist
	createTableStmt := `
		CREATE TABLE IF NOT EXISTS sneaker_table (
			id SERIAL PRIMARY KEY,
			name TEXT,
			price TEXT,
			link TEXT
		)
	`
	_, err = db.Exec(createTableStmt)
	if err != nil {
		log.Printf("Failed to create table: %v", err)
		return
	}

	_, err = db.Exec("DELETE FROM sneaker_table")
	if err != nil {
		log.Printf("Error deleting existing data: %v", err)
		return
	}

	// Scrape data from Adidas
	adidasShoes := scrapper.ScrapeProducts("Adidas")
	for _, shoe := range adidasShoes {
		_, err = db.Exec("INSERT INTO sneaker_table (name, price, link) VALUES ($1, $2, $3)",
			shoe.NAME, shoe.PRICE, shoe.LINK)
		if err != nil {
			log.Printf("Error inserting Adidas data: %v", err)
			return
		}
	}

	// Scrape data from New Balance
	nbShoes := scrapper.ScrapeProducts("New Balance")
	for _, shoe := range nbShoes {
		_, err = db.Exec("INSERT INTO sneaker_table (name, price, link) VALUES ($1, $2, $3)",
			shoe.NAME, shoe.PRICE, shoe.LINK)
		if err != nil {
			log.Printf("Error inserting New Balance data: %v", err)
			return
		}
	}

	// Scrape data from Nike
	newShoes, err := scrapper.ScrapeProductsx(false)
	if err != nil {
		log.Printf("Error scraping new data: %v", err)
		return
	}

	for _, shoe := range newShoes {
		_, err = db.Exec("INSERT INTO sneaker_table (name, price, link) VALUES ($1, $2, $3)",
			shoe.NAME, shoe.PRICE, shoe.LINK)
		if err != nil {
			log.Printf("Error inserting new data: %v", err)
			return
		}
	}

	log.Println("Scraped data refreshed successfully")
}
