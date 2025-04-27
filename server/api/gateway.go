package api

import (
	// auth "modular_monolith/module/auth"
	user "modular_monolith/module/user"

	"github.com/gofiber/fiber/v2"
)

func ApiV1Routes(app *fiber.App) {
	api := app.Group("/api")

	v1 := api.Group("/v1")

	// User routes
	userGroup := v1.Group("/users")
	user.RegisterApp(userGroup)

	// Tambahin module lain di sini kalau ada ya~ (*≧ω≦)
}
