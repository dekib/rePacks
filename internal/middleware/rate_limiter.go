package middleware

import (
	"github.com/dekib/rePacks/internal/errors"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"strconv"
	"time"
)

// RateLimiter creates a rate limiting middleware
func RateLimiter(rateStr string) gin.HandlerFunc {
	// Create rate limiter with in-memory store
	store := memory.NewStore()

	// Create rate limiter instance
	rate, err := limiter.NewRateFromFormatted(rateStr)
	if err != nil {
		panic(err) // This should only happen at startup if rate format is invalid
	}

	instance := limiter.New(store, rate)

	return func(c *gin.Context) {
		// Get context for the request
		context, err := instance.Get(c, c.ClientIP())
		if err != nil {
			errors.NewInternalServerError(
				errors.ErrInternalServer,
				err,
				nil,
			).Abort(c)
			return
		}

		// Set rate limit headers
		c.Header("X-RateLimit-Limit", strconv.FormatInt(context.Limit, 10))
		c.Header("X-RateLimit-Remaining", strconv.FormatInt(context.Remaining, 10))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(context.Reset, 10))

		// Check if limit exceeded
		if context.Reached {
			retryAfter := strconv.FormatInt(context.Reset-time.Now().Unix(), 10)
			c.Header("Retry-After", retryAfter)
			errors.NewRateLimitError(retryAfter).Abort(c)
			return
		}

		c.Next()
	}
}
