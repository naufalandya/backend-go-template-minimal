package user_module

import (
	"modular_monolith/module/user/auto"
	"modular_monolith/module/user/controllers"
	rabbitmq "modular_monolith/server/config/rabbit"

	"github.com/gofiber/fiber/v2"
)

func RegisterApp(router fiber.Router) {

	rabbitmq.InitRabbitMQ()

	go auto.SubscribeToChannel("my_channel")

	go auto.ConsumerRabbit()

	router.Get("/", controllers.GetAllUsers)
	// router.Get("/:id", middlewares.BearerTokenAuth, controllers.GetUserByID)
	router.Post("/", controllers.RegisterUserSimple)
	router.Put("/lol", controllers.UploadFile)
}
