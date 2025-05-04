package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"modular_monolith/server/api"
	"modular_monolith/server/client"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// Custom Fiber Logger Middleware using logrus
func logrusLogger() fiber.Handler {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Set log level based on environment
	if os.Getenv("ENV") == "development" {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}

	return func(c *fiber.Ctx) error {
		// Log the request
		start := time.Now()
		err := c.Next()

		// Log response details
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
	// // Initialize database and Redis connections
	// if err := db.InitDB(); err != nil {
	// 	logrus.Fatalf("Failed to initialize database: %v", err)
	// }
	// if err := redis.InitRedis(); err != nil {
	// 	logrus.Fatalf("Failed to initialize Redis: %v", err)
	// }

	// Create Fiber app
	app := fiber.New(fiber.Config{
		BodyLimit: 500 * 1024 * 1024,
	})
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		logrus.Fatal("Error loading .env file")
	}

	// Connect to gRPC server
	err = client.Connect()
	if err != nil {
		logrus.Fatalf("Failed to connect to gRPC: %v", err)
	}

	// Use custom logrus logger middleware
	app.Use(logrusLogger())

	// Load API configuration
	// cfg := api.LoadConfig()

	// Create a channel to listen for termination signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create a context with a timeout for the graceful shutdown process
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Start the gRPC server in a goroutine
	// go func() {
	// 	if err := api.Start(cfg); err != nil {
	// 		logrus.Errorf("❌ gRPC server error: %v", err) // Log error without shutting down
	// 	}
	// }()

	// Serve Swagger API documentation
	app.Static("/swagger", "./docs")
	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL: "/swagger/v1.yaml",
	}))

	// Setup routes for the API
	api.ApiV1Routes(app)

	// Start the Fiber server in a goroutine
	go func() {
		if err := app.Listen(":8081"); err != nil {
			logrus.Errorf("❌ Fiber server error: %v", err) // Log error without shutting down
		}
	}()

	// Wait for termination signal or error
	select {
	case sig := <-sigChan:
		logrus.Printf("Received signal: %s, shutting down...", sig)
		// Initiate graceful shutdown for Fiber
		if err := app.Shutdown(); err != nil {
			logrus.Fatalf("Fiber shutdown error: %v", err)
		}
	case err := <-waitForError():
		logrus.Printf("Received error: %v, shutting down...", err)
		// Handle error and initiate shutdown
		if err := app.Shutdown(); err != nil {
			logrus.Fatalf("Fiber shutdown error due to service failure: %v", err)
		}
	}

	<-ctx.Done()
	logrus.Println("Graceful shutdown completed")
}

// waitForError waits for an error from gRPC or other sources
func waitForError() <-chan error {
	errChan := make(chan error, 1)

	// // Example: You could listen for errors from gRPC, database, or any other service
	// go func() {
	// 	// Placeholder: error goes to channel
	// 	errChan <- fmt.Errorf("some error occurred in the system") // Simulating an error
	// }()

	return errChan
}
