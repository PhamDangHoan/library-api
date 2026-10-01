package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const (
	rateLimitRequests = 60
	rateLimitWindow   = time.Minute
)

func RateLimit(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()

		ip := c.ClientIP()

		key := fmt.Sprintf("rate_limit:%s", ip)

		count, err := redisClient.Incr(ctx, key).Result()

		if err != nil {
			// Redis lỗi thì không chặn request.
			// API vẫn tiếp tục hoạt động.
			c.Next()
			return
		}

		if count == 1 {
			redisClient.Expire(ctx, key, rateLimitWindow)
		}

		if count > rateLimitRequests {
			c.Header("Retry-After", "60")

			c.AbortWithStatusJSON(
				http.StatusTooManyRequests,
				gin.H{
					"success": false,
					"message": "Too many requests. Please try again later.",
				},
			)

			return
		}

		c.Header(
			"X-RateLimit-Limit",
			fmt.Sprintf("%d", rateLimitRequests),
		)

		remaining := rateLimitRequests - int(count)

		if remaining < 0 {
			remaining = 0
		}

		c.Header(
			"X-RateLimit-Remaining",
			fmt.Sprintf("%d", remaining),
		)

		c.Next()
	}
}
