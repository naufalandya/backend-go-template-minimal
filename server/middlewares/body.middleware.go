package middlewares

import (
	"github.com/gofiber/fiber/v2"
)

func MaxBodySizeMiddleware(limit int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		body := c.Body()

		if len(body) > limit {
			return c.Status(fiber.StatusRequestEntityTooLarge).SendString("Request body too large")
		}

		return c.Next()
	}
}
