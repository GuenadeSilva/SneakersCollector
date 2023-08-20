package main

import (
	//"encoding/json"
	//"fmt"
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

func main() {
	r := gin.Default()
	r.GET("/protected", protectedHandler)
	go sc.StartScheduler() // Run the scheduler in a goroutine
	r.Run(":8481")
}
