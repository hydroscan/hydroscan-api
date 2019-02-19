package main

import (
	"github.com/hydroscan/hydroscan-api/models"
	"github.com/hydroscan/hydroscan-api/task"
	"github.com/jasonlvhit/gocron"
	_ "github.com/joho/godotenv/autoload"
	log "github.com/sirupsen/logrus"
)

type taskType func()

// Ensure cron don't crash when one task failed(panic)
func safeTask(fn taskType) taskType {
	return func() {
		defer func() {
			if err := recover(); err != nil {
				log.Warn(err)
			}
		}()
		fn()
	}
}

func main() {
	models.Connect()
	defer models.Close()
	task.InitEthClient()

	log.Info("cron running")

	// Fetch a mass of missing logs before cron task
	// task.FetchHistoricalLogs()

	// gocron.Every(1).Minutes().Do(safeTask(task.FetchHistoricalLogs))
	gocron.Every(5).Minutes().Do(safeTask(task.UpdateTokenPrices))
	gocron.Every(1).Day().Do(safeTask(task.UpdateRelayers))

	_, time := gocron.NextRun()
	log.Info(time)
	<-gocron.Start()
}
