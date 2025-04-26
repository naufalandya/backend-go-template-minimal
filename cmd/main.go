package main

import (
	"log"
	"modular_monolith/server/api"
	db "modular_monolith/server/config/db"
	middlewares "modular_monolith/server/middlewares"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
)

func main() {
	db.InitDB()

	app := fiber.New()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	env := os.Getenv("ENV")
	if env == "development" {
		middlewares.SetupLogger(app)
		app.Use(logger.New())
	}

	// cek disini

	// http://localhost:8080/swagger/index.html#

	app.Static("/swagger", "./docs")

	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL: "/swagger/v1.yaml",
	}))

	api.ApiV1Routes(app)

	log.Fatal(app.Listen(":8080"))
}
