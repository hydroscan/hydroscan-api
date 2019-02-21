package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hydroscan/hydroscan-api/api"
	"github.com/hydroscan/hydroscan-api/middleware"
	"github.com/hydroscan/hydroscan-api/models"
	"github.com/hydroscan/hydroscan-api/redis"
	"github.com/hydroscan/hydroscan-api/task"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	models.Connect()
	defer models.Close()
	redis.Connect()
	task.InitEthClient()

	r := gin.Default()
	r.ForwardedByClientIP = true
	r.Use(middleware.Limit())
	r.Use(middleware.CORS())
	api.ApplyRoutes(r)
	r.Run()
}
