package legal

import (
	"modular_monolith/server/api/legal/controllers"
	"modular_monolith/server/middlewares"

	"github.com/gofiber/fiber/v2"
)

func RegisterApp(router fiber.Router) {

	// rabbitmq.InitRabbitMQ()

	// go auto.SubscribeToChannel("my_channel")

	// go auto.ConsumerRabbit()

	router.Post("/",
		// middlewares.NewRateLimiterMiddleware(5, 10),
		middlewares.BearerTokenAuth,
		middlewares.MaxBodySizeMiddleware(500*1024*1024), // Allow up to 500MB
		controllers.UploadFile)
}
