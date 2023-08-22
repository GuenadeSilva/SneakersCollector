package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"sneakercollector/scrapper"

	_ "github.com/lib/pq" // PostgreSQL driver
)

var db *sql.DB // Global database connection

func SetupDB(username, password, host string) error {
	connStr := fmt.Sprintf("postgres://%s:%s@%s/sneaker_db?sslmode=disable", username, password, host)
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	// Set connection pool settings if needed
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	// Test the database connection
	err = db.Ping()
	if err != nil {
		return err
	}

	return nil
}

func RefreshScrapedData() {
	// Drop the sneaker_table if it exists
	dropTableStmt := `
		DROP TABLE IF EXISTS sneaker_table
	`
	_, err := db.Exec(dropTableStmt)
	if err != nil {
		log.Printf("Failed to drop table: %v", err)
		return
	}

	// Create the sneaker_table
	createTableStmt := `
		CREATE TABLE IF NOT EXISTS sneaker_table (
			id SERIAL PRIMARY KEY,
			name TEXT,
			price TEXT,
			link TEXT,
			brand TEXT
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
		timestamp TIMESTAMP,
		status TEXT,
		elapsed_time INTERVAL
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
		_, err := db.Exec("INSERT INTO log_table (message, timestamp, status, elapsed_time) VALUES ($1, $2, $3, $4)",
			logMessage, time.Now(), status, elapsedTime)
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
			_, err := db.Exec("INSERT INTO sneaker_table (name, price, link, brand) VALUES ($1, $2, $3, $4)",
				shoe.NAME, shoe.PRICE, shoe.LINK, brand)
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

// GetSneakerData retrieves data from the sneaker_table
func GetSneakerData() ([]scrapper.ShoeInfo, error) {
	rows, err := db.Query("SELECT name, price, link FROM sneaker_table")
	if err != nil {
		log.Printf("Error querying database: %v", err)
		return nil, err
	}
	defer rows.Close()

	var sneakers []scrapper.ShoeInfo
	for rows.Next() {
		var shoe scrapper.ShoeInfo
		err := rows.Scan(&shoe.NAME, &shoe.PRICE, &shoe.LINK)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}
		sneakers = append(sneakers, shoe)
	}

	return sneakers, nil
}

// GetLatestShoesForBrand retrieves the latest scraped shoes for a specific brand
func GetLatestShoesForBrand(brand string) ([]scrapper.ShoeInfo, error) {
	rows, err := db.Query("SELECT name, price, link FROM sneaker_table WHERE brand=$1 ORDER BY id DESC", brand)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		return nil, err
	}
	defer rows.Close()

	var shoes []scrapper.ShoeInfo
	for rows.Next() {
		var shoe scrapper.ShoeInfo
		err := rows.Scan(&shoe.NAME, &shoe.PRICE, &shoe.LINK)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}
		shoes = append(shoes, shoe)
	}

	return shoes, nil
}
