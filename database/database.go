package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

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

	// Create the sneaker_table if it doesn't exist
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

	createLogTableStmt := `
	CREATE TABLE IF NOT EXISTS log_table (
		id SERIAL PRIMARY KEY,
		message TEXT,
		timestamp TIMESTAMP
	)
`
	_, err = db.Exec(createLogTableStmt)
	if err != nil {
		log.Printf("Failed to create log table: %v", err)
		return
	}

	// Function to log each step
	logStep := func(operation, brand string, startTime time.Time, rowsInserted int, success bool) {
		endTime := time.Now()
		elapsedTime := endTime.Sub(startTime)
		status := "successful"
		if !success {
			status = "failed"
		}
		logMessage := fmt.Sprintf("Brand: %s\n", brand)
		logMessage += fmt.Sprintf("Operation: %s\n", operation)
		logMessage += fmt.Sprintf("Status: %s\n", status)
		logMessage += fmt.Sprintf("Time taken: %s\n", elapsedTime)
		logMessage += fmt.Sprintf("Rows inserted: %d\n", rowsInserted)
		logMessage += fmt.Sprintf("Timestamp: %s\n", endTime)
		logMessage += "---------------------------\n"

		// Insert the log into log_table
		_, err = db.Exec("INSERT INTO log_table (message, timestamp) VALUES ($1, $2)",
			logMessage, time.Now())
		if err != nil {
			log.Printf("Error inserting log into log_table: %v", err)
		}
	}

	// Scrape and insert data for different brands
	scrapeAndInsert := func(brand string) {
		startTime := time.Now()
		var shoes []scrapper.ShoeInfo
		if brand == "Nike" {
			shoes, _ = scrapper.ScrapeProductsx(false)
		} else {
			shoes = scrapper.ScrapeProducts(brand)
		}
		rowsInserted := 0
		for _, shoe := range shoes {
			_, err = db.Exec("INSERT INTO sneaker_table (name, price, link) VALUES ($1, $2, $3)",
				shoe.NAME, shoe.PRICE, shoe.LINK)
			if err == nil {
				rowsInserted++
			}
		}
		logStep(fmt.Sprintf("ScrapeProducts('%s')", brand), brand, startTime, rowsInserted, err == nil)
	}

	// Scrape data for different brands
	brands := []string{"Adidas", "New Balance", "Nike"}
	for _, brand := range brands {
		scrapeAndInsert(brand)
	}

	log.Println("Scraped data refresh completed")
}
