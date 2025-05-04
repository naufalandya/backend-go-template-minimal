package user_module

import (
	"modular_monolith/server/api/user/auto"
	"modular_monolith/server/api/user/controllers"
	rabbitmq "modular_monolith/server/config/rabbit"
	"modular_monolith/server/middlewares"

	ws "modular_monolith/server/config/websocket"

	"github.com/gofiber/fiber/v2"
)

func RegisterApp(router fiber.Router) {

	rabbitmq.InitRabbitMQ()

	go auto.SubscribeToChannel("my_channel")

	go auto.ConsumerRabbit()

	router.Get("/", middlewares.NewRateLimiterMiddleware(5, 10), controllers.GetAllUsers)
	router.Post("/", middlewares.NewRateLimiterMiddleware(3, 15), controllers.RegisterUserSimple)

	router.Put("/lol", controllers.UploadFile)
	router.Post("/hello", controllers.SayHello)
	router.Get("/ws", ws.UpgradeMiddleware())

}
