package main

import (
	_ "log"
	"net/http"
	_ "os"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"sneakercollector/database"
	sc "sneakercollector/scheduler"
)

func protectedHandler(c *gin.Context) {
	action := c.Query("action")

	switch action {
	case "latest_shoes":
		brand := c.Query("brand")
		if brand == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Brand parameter is missing"})
			return
		}

		// Retrieve and return data from the sneaker_table for the selected brand
		latestShoes, err := database.GetLatestShoesForBrand(brand)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		c.JSON(http.StatusOK, latestShoes)

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
	maxRetries := 2

	r := gin.Default()

	// Set up CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"} // Add the URL of your React app
	config.AllowHeaders = []string{"Origin", "Content-Type"}

	// Use the CORS middleware
	r.Use(cors.New(config))

	r.POST("/connect", func(c *gin.Context) {
		if maxRetries <= 0 {
			c.String(http.StatusForbidden, "Maximum retries exceeded")
			return
		}

		var loginInfo struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Host     string `json:"host"`
		}

		if err := c.BindJSON(&loginInfo); err != nil {
			c.String(http.StatusBadRequest, "Invalid JSON")
			return
		}

		username := loginInfo.Username
		password := loginInfo.Password
		host := loginInfo.Host

		err := database.SetupDB(username, password, host)
		if err != nil {
			maxRetries--
			c.String(http.StatusInternalServerError, "Error connecting to the database")
			return
		}

		maxRetries = 2
		c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
	})

	r.GET("/protected", protectedHandler)

	go sc.StartScheduler()

	// Start the server in a goroutine
	r.Run(":8483")
}
