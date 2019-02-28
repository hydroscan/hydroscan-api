package main

import (
	"time"

	"github.com/hydroscan/hydroscan-api/config"
	"github.com/hydroscan/hydroscan-api/models"
	"github.com/hydroscan/hydroscan-api/redis"
	"github.com/hydroscan/hydroscan-api/task"
	_ "github.com/joho/godotenv/autoload"
	log "github.com/sirupsen/logrus"
)

func subscribeRetryWrapper() {
	defer func() {
		if err := recover(); err != nil {
			log.Warn(err)
		}
	}()
	task.SubscribeLogs()
}

func main() {
	config.Load()

	models.Connect()
	defer models.Close()
	redis.Connect()
	task.InitEthClient()

	log.Info("subscriber running")
	for {
		subscribeRetryWrapper()
		time.Sleep(3000 * time.Millisecond)
		log.Info("subscriber retry")
	}
}
