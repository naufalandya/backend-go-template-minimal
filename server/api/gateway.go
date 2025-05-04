package api

import (
	// auth "modular_monolith/module/auth"
	legal "modular_monolith/server/api/legal"

	"github.com/gofiber/fiber/v2"
)

func ApiV1Routes(app *fiber.App) {
	api := app.Group("/api")

	v1 := api.Group("/v1")

	legalGroup := v1.Group("/legal")
	legal.RegisterApp(legalGroup)

}
