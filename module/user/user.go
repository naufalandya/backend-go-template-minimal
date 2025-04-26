package user_module

import (
	"modular_monolith/module/user/controllers"
	"modular_monolith/server/middlewares"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(router fiber.Router) {
	router.Get("/", controllers.GetAllUsers)
	router.Get("/:id", middlewares.BearerTokenAuth, controllers.GetUserByID)
	router.Post("/", controllers.CreateUser)
}
