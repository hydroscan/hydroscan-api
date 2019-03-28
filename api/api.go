package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hydroscan/hydroscan-api/api/v1"
)

func ApplyRoutes(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	api := r.Group("/api")
	{
		apiv1.ApplyRoutes(api)
	}
}
