package redis

import (
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var Client *redis.Client

func Connect() {
	opt, err := redis.ParseURL(viper.GetString("redis_url"))
	if err != nil {
		panic(err)
	}
	Client = redis.NewClient(opt)
	log.Info("Reids Connected")
}
