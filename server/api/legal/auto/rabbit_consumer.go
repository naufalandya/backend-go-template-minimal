package auto

import (
	"log"
	"modular_monolith/server/api/user/services"
	rabbitmq "modular_monolith/server/config/rabbit"
	"os"
	"os/signal"
	"syscall"
)

func ConsumerRabbit() {
	err := services.ConsumeMessagesRabbit("your-queue-name", handleMessage)
	if err != nil {
		log.Fatalf("Failed to start consuming messages: %v", err)
	}

	log.Println("🐇 RabbitMQ consumer started... Waiting for messages ✨")

	waitForShutdown()

}

func handleMessage(msg string) error {
	log.Printf("🌸 Processing message: %s", msg)

	// 👉 Put your business logic here
	// Example: maybe process orders, send emails, etc

	return nil // return error if processing failed
}

func waitForShutdown() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Println("🌙 Shutting down gracefully... bye-bye~")

	defer rabbitmq.CloseRabbitMQ()

}
