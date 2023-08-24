package main

import (
	//"encoding/json"
	"fmt"
	_ "log"
	"net/http"
	_ "os"
	"os/exec"
	"runtime"
	"time"

	"html/template"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"sneakercollector/database"
	sc "sneakercollector/scheduler"
)

const formTemplate = `
<!DOCTYPE html>
<html>
<head>
	<title>Database Connection</title>
</head>
<body>
	<h1>Database Connection Details</h1>
	<form action="/connect" method="post">
		<label for="username">Username:</label>
		<input type="text" id="username" name="username"><br><br>

		<label for="password">Password:</label>
		<input type="password" id="password" name="password"><br><br>

		<label for="host">Host:</label>
		<input type="text" id="host" name="host"><br><br>

		<input type="submit" value="Connect">
	</form>
</body>
</html>
`

const optionsTemplate = `
<!DOCTYPE html>
<html>
<head>
	<title>Options</title>
</head>
<body>
	<h1>Select an Option</h1>
	<ul>
		<li><a href="/protected?action=sneaker_db_data">View all data as a table</a></li>
		<li><a href="/protected?action=refresh_data">Refresh data</a></li>
		<li>View latest shoes by brand:
			<ul>
				<li><a href="/protected?action=latest_shoes&brand=Nike">Nike</a></li>
				<li><a href="/protected?action=latest_shoes&brand=Adidas">Adidas</a></li>
				<li><a href="/protected?action=latest_shoes&brand=New%20Balance">New Balance</a></li>
			</ul>
		</li>
	</ul>
</body>
</html>

`

const dataTableTemplate = `
<!DOCTYPE html>
<html>
<head>
	<title>All Data</title>
	<style>
		table {
			border-collapse: collapse;
			width: 100%;
		}
		th, td {
			border: 1px solid black;
			padding: 8px;
			text-align: left;
		}
	</style>
</head>
<body>
	<h1>All Data</h1>
	<table>
		<thead>
			<tr>
				<th>ID</th>
				<th>Name</th>
				<th>Brand</th>
				<!-- Add more columns as needed -->
			</tr>
		</thead>
		<tbody>
			{{ range . }}
				<tr>
					<td>{{ .ID }}</td>
					<td>{{ .NAME }}</td>
					<td>{{ .BRAND }}</td>
					<!-- Add more columns as needed -->
				</tr>
			{{ end }}
		</tbody>
	</table>
</body>
</html>
`

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

	r.GET("/connect-form", func(c *gin.Context) {
		tmpl, err := template.New("form").Parse(formTemplate)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error rendering form")
			return
		}
		tmpl.Execute(c.Writer, nil)
	})

	r.POST("/connect", func(c *gin.Context) {
		if maxRetries <= 0 {
			c.String(http.StatusForbidden, "Maximum retries exceeded")
			return
		}

		username := c.PostForm("username")
		password := c.PostForm("password")
		host := c.PostForm("host")

		err := database.SetupDB(username, password, host)
		if err != nil {
			maxRetries--
			c.String(http.StatusInternalServerError, "Error connecting to the database")
			return
		}

		maxRetries = 2
		c.Redirect(http.StatusSeeOther, "/options")
	})

	r.GET("/options", func(c *gin.Context) {
		tmpl, err := template.New("options").Parse(optionsTemplate)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error rendering options page")
			return
		}
		tmpl.Execute(c.Writer, nil)
	})

	r.GET("/view-all-data", func(c *gin.Context) {
		// Retrieve and return all data from the sneaker_table
		allData, err := database.GetSneakerData()
		if err != nil {
			c.String(http.StatusInternalServerError, "Error retrieving data from the database")
			return
		}

		// Render the data table template with the retrieved data
		tmpl, err := template.New("datatable").Parse(dataTableTemplate)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error rendering data table")
			return
		}
		tmpl.Execute(c.Writer, allData)
	})

	r.GET("/protected", protectedHandler)

	go sc.StartScheduler()

	// Start the server in a goroutine
	r.Run(":8483")

	// Wait for a short time to ensure the server has started
	// before attempting to open the browser
	time.Sleep(time.Second * 2)

	// Open the connect-form URL in the default browser
	openBrowser("http://localhost:8483/connect-form")
}

func openBrowser(url string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		fmt.Println("Opening browser on Linux:", url)
		cmd = exec.Command("xdg-open", url)
	case "darwin":
		fmt.Println("Opening browser on macOS:", url)
		cmd = exec.Command("open", url)
	case "windows":
		fmt.Println("Opening browser on Windows:", url)
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		fmt.Println("Unsupported operating system")
		return
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("Failed to open browser:", err)
	}
}
