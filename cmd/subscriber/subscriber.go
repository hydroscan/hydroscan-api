package main

import (
	"time"

	"github.com/hydroscan/hydroscan-api/models"
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
	time.Sleep(3000 * time.Millisecond)
	log.Info("subscriber retry")
}

func main() {
	models.Connect()
	defer models.Close()
	task.InitEthClient()
	log.Info("subscriber running")
	for {
		subscribeRetryWrapper()
	}
}
