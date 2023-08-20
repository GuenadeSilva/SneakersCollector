package scheduler

import (
	"fmt"

	"sneakercollector/database"

	"github.com/robfig/cron"
)

func StartScheduler() {
	c := cron.New()

	// Schedule RefreshScrapedData function to run every Sunday at midnight
	c.AddFunc("0 0 * * 0", database.RefreshScrapedData)

	c.Start()

	select {}
}

func RunScheduler() {
	fmt.Println("Scheduler started")
	StartScheduler()
}
