package main

import (
	//"encoding/json"
	//"fmt"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"sneakercollector/database"
	sc "sneakercollector/scheduler"
)

func protectedHandler(c *gin.Context) {
	action := c.Query("action")

	switch action {
	case "latest_run":
		// Retrieve and return the latest log entry
		logEntry, err := database.GetLatestLogEntry()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		c.JSON(http.StatusOK, logEntry)

	case "sneaker_db_data":
		// Retrieve and return data from the sneaker_table
		sneakerData, err := database.GetSneakerData()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		c.JSON(http.StatusOK, sneakerData)

	case "refresh_data":
		// Trigger the RefreshScrapedData function
		database.RefreshScrapedData()
		c.String(http.StatusOK, "Data refresh initiated")

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action"})
		return
	}
}

func getUserInput(prompt string) string {
	var input string
	fmt.Print(prompt)
	fmt.Scanln(&input)
	return input
}

func main() {
	maxRetries := 2

	for maxRetries > 0 {
		username := getUserInput("Enter database username: ")
		password := getUserInput("Enter database password: ")
		host := getUserInput("Enter database host: ")

		err := database.SetupDB(username, password, host)
		if err != nil {
			maxRetries--
			fmt.Printf("Error connecting to database: %v\n", err)
			fmt.Printf("Retries left: %d\n", maxRetries)
			if maxRetries == 0 {
				log.Fatal("Failed to establish database connection")
			}
		} else {
			// Successfully connected
			break
		}
	}

	r := gin.Default()
	r.GET("/protected", protectedHandler)
	go sc.StartScheduler() // Run the scheduler in a goroutine
	r.Run(":8482")
}
