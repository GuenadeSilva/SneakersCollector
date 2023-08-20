package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"sneakercollector/database"
	sc "sneakercollector/scheduler"
)

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the query parameters
	query := r.URL.Query()
	action := query.Get("action")

	switch action {
	case "latest_run":
		// Retrieve and return the latest log entry
		logEntry, err := database.GetLatestLogEntry()
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(logEntry)

	case "sneaker_db_data":
		// Retrieve and return data from the sneaker_table
		sneakerData, err := database.GetSneakerData()
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(sneakerData)

	case "refresh_data":
		// Trigger the RefreshScrapedData function
		database.RefreshScrapedData()
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Data refresh initiated")

	default:
		http.Error(w, "Invalid action", http.StatusBadRequest)
		return
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/protected", protectedHandler).Methods("GET")
	sc.StartScheduler()
	http.Handle("/", r)
	http.ListenAndServe(":8481", nil)
}
