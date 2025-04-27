package middlewares

import (
	"fmt"
	rdb "modular_monolith/server/config/redis"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

func NewRateLimiterMiddleware(limit int, expiration int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		key := fmt.Sprintf("rate_limit:%s", c.IP())

		redisClient := rdb.RedisClient

		currentCount, err := redisClient.Get(c.Context(), key).Int()
		if err != nil && err != redis.Nil {
			return err
		}

		if currentCount >= limit {
			return c.Status(fiber.StatusTooManyRequests).SendString("Rate limit exceeded. Please try again later.")
		}

		_, err = redisClient.Incr(c.Context(), key).Result()
		if err != nil {
			return err
		}

		if currentCount == 0 {
			redisClient.Expire(c.Context(), key, time.Duration(expiration)*time.Second)
		}

		return c.Next()
	}
}
