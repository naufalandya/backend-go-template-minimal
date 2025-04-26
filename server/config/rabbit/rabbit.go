package rabbitmq

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

var RabbitMQChannel *amqp.Channel
var RabbitMQConn *amqp.Connection

// InitRabbitMQ initializes RabbitMQ connection and channel.
func InitRabbitMQ() {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Read environment variables
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		rabbitMQURL = "amqp://guest:guest@localhost:5672/"
	}

	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	RabbitMQConn = conn
	fmt.Println("Connected to RabbitMQ~ ₍˶ˆ꒳ˆ˶₎✨")

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	RabbitMQChannel = ch
	fmt.Println("RabbitMQ channel created successfully~ ✨")

	// Close the channel and connection when the application stops
	// This will be done using defer in the main function (or the place where you close the app)
}

func CloseRabbitMQ() {
	if RabbitMQChannel != nil {
		RabbitMQChannel.Close()
	}
	if RabbitMQConn != nil {
		RabbitMQConn.Close()
	}
}
