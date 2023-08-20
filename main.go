package main

import (
	"fmt"
	db "sneakercollector/database"
	_ "sneakercollector/scheduler"
)

func main() {
	fmt.Println("Sneaker Collector")

	// Start the scheduler
	//scheduler.RunScheduler()
	db.RefreshScrapedData()
}
