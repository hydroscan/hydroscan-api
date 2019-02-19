package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"time"
)

func Limit() gin.HandlerFunc {
	// Define a limit rate to 300 requests per minute.
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  120,
	}

	// Use a in-memory store with a goroutine which clears expired keys.
	store := memory.NewStore()

	return mgin.NewMiddleware(limiter.New(store, rate))
}
