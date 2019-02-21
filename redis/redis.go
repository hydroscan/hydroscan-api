package redis

import (
	"os"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

var Client *redis.Client

func Connect() {
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		panic(err)
	}
	Client = redis.NewClient(opt)
	log.Info("Reids Connected")
}
