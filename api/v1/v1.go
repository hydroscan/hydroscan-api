package apiv1

import (
	"github.com/gin-gonic/gin"
)

func ApplyRoutes(r *gin.RouterGroup) {
	v1 := r.Group("/v1")
	{
		v1.GET("/tokens", GetTokens)
		v1.GET("/tokens/:address", GetToken)
		v1.GET("/tokens/:address/chart", GetTokenChart)

		v1.GET("/relayers", GetRelayers)
		v1.GET("/relayers/:slug", GetRelayer)

		v1.GET("/trades", GetTrades)
		v1.GET("/trades/:uuid", GetTrade)
		v1.GET("/trades_chart", GetTradesChart)
		v1.GET("/trades_indicators", GetTradesIndicators)
	}
}
