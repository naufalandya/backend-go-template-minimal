package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"modular_monolith/server/api"
	"modular_monolith/server/client"
	"modular_monolith/server/config/db"
	"modular_monolith/server/config/elastic"

	"github.com/gofiber/fiber/v2"
	cors "github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func logrusLogger() fiber.Handler {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	if os.Getenv("ENV") == "development" {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}

	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()

		duration := time.Since(start)
		logger.WithFields(logrus.Fields{
			"status":  c.Response().StatusCode(),
			"method":  c.Method(),
			"path":    c.Path(),
			"latency": duration,
		}).Info("Request completed")

		return err
	}
}

func main() {

	app := fiber.New(fiber.Config{
		BodyLimit: 500 * 1024 * 1024,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // allow all origins
	}))

	err := godotenv.Load(".env")
	if err != nil {
		logrus.Fatal("Error loading .env file")
	}

	db.InitDB()

	err = elastic.InitElasticSearch()
	if err != nil {
		log.Fatalf("Error initializing Elasticsearch: %s", err)
		return
	}

	fmt.Println("Elasticsearch client is initialized and ready.")

	client.Clients, err = client.Connect()
	if err != nil {
		logrus.Fatalf("Failed to connect to gRPC: %v", err)
	}

	app.Use(logrusLogger())

	// cfg := api.LoadConfig()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// go func() {
	// 	if err := api.Start(cfg); err != nil {
	// 		logrus.Errorf("❌ gRPC serCORSver error: %v", err) // Log error without shutting down
	// 	}
	// }()

	app.Static("/swagger", "./docs")
	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL: "/swagger/v1.yaml",
	}))

	api.ApiV1Routes(app)

	go func() {
		if err := app.Listen(":8081"); err != nil {
			logrus.Errorf("❌ Fiber server error: %v", err) // Log error without shutting down
		}
	}()

	select {
	case sig := <-sigChan:
		logrus.Printf("Received signal: %s, shutting down...", sig)
		if err := app.Shutdown(); err != nil {
			logrus.Fatalf("Fiber shutdown error: %v", err)
		}
	case err := <-waitForError():
		logrus.Printf("Received error: %v, shutting down...", err)
		if err := app.Shutdown(); err != nil {
			logrus.Fatalf("Fiber shutdown error due to service failure: %v", err)
		}
	}

	<-ctx.Done()
	logrus.Println("Graceful shutdown completed")
}

func waitForError() <-chan error {
	errChan := make(chan error, 1)

	// // Example: You could listen for errors from gRPC, database, or any other service
	// go func() {
	// 	// Placeholder: error goes to channel
	// 	errChan <- fmt.Errorf("some error occurred in the system") // Simulating an error
	// }()

	return errChan
}
