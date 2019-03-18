package main

import (
	"github.com/hydroscan/hydroscan-api/config"
	"github.com/hydroscan/hydroscan-api/models"
	"github.com/hydroscan/hydroscan-api/redis"
	"github.com/hydroscan/hydroscan-api/task"
	"github.com/jasonlvhit/gocron"
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
	config.Load()

	models.Connect()
	defer models.Close()
	redis.Connect()
	task.InitEthClient()

	log.Info("cron running")
	gocron.Every(5).Minutes().Do(safeTask(task.UpdateTokenPrices))
	gocron.Every(1).Day().Do(safeTask(task.UpdateRelayers))

	gocron.Every(5).Minutes().Do(safeTask(task.UpdateTokenTrades24h))

	gocron.Every(5).Minutes().Do(safeTask(task.UpdateTokenVolume24h))
	gocron.Every(30).Minutes().Do(safeTask(task.UpdateTokenVolume7d))
	// gocron.Every(60).Minutes().Do(safeTask(task.UpdateTokenVolume1m))
	gocron.Every(60).Minutes().Do(safeTask(task.UpdateTokenVolumeAll))

	gocron.Every(5).Minutes().Do(safeTask(task.UpdateIndicators))
	// Update Trade PriceUSD and VolumeUSD which is nil
	gocron.Every(60).Minutes().Do(safeTask(task.UpdateHistoryTradePrice))
	gocron.Every(60).Minutes().Do(safeTask(task.UpdateOnlyVolumeUSD))

	// first time start run all tasks right now
	gocron.RunAll()

	_, time := gocron.NextRun()
	log.Info("next corn: ", time)
	// Starts the schedule as per normal
	<-gocron.Start()
}
