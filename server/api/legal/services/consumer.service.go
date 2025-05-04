package services

import (
	"fmt"
	"log"
	rabbitmq "modular_monolith/server/config/rabbit"
)

func ConsumeMessagesRabbit(queueName string, handlerFunc func(message string) error) error {
	_, err := rabbitmq.RabbitMQChannel.QueueDeclare(
		queueName, // Queue name
		true,      // Durable
		false,     // Delete when unused
		false,     // Exclusive
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// Consume messages
	messages, err := rabbitmq.RabbitMQChannel.Consume(
		queueName,
		"",
		false, // Auto-ack = false (we will ack manually!)
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %w", err)
	}

	// Process messages
	go func() {
		for d := range messages {
			log.Printf("🐇 Received message: %s", d.Body)

			// Process the message
			if err := handlerFunc(string(d.Body)); err != nil {
				log.Printf("💔 Error processing message: %v", err)
				// Nack the message to requeue it
				d.Nack(false, true)
				continue
			}

			// Ack the message after success
			d.Ack(false)
		}
	}()

	return nil
}
