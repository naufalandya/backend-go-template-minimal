package services

import (
	"context"
	"fmt"
	"log"
	rabbitmq "modular_monolith/server/config/rabbit"
	"modular_monolith/server/config/redis"

	"github.com/streadway/amqp"
)

func PublishMessage(message string) error {
	// Define your channel name
	channel := "my_channel"

	err := redis.RedisClient.Publish(context.Background(), channel, message).Err()
	if err != nil {
		// Log the error and return it for further handling
		log.Printf("Could not publish message: %v", err)
		return err
	}

	// Optionally, print out a success message
	fmt.Println("Message published to Redis channel:", channel)
	return nil
}

func PublishMessageRabbit(queueName, message string) error {
	// Check if RabbitMQ connection or channel is nil
	if rabbitmq.RabbitMQChannel == nil || rabbitmq.RabbitMQConn == nil {
		return fmt.Errorf("RabbitMQ connection or channel is not open. Please check the RabbitMQ connection.")
	}

	// Try publishing a message
	err := rabbitmq.RabbitMQChannel.Publish(
		"",        // Default exchange
		queueName, // Routing key (queue name)
		false,     // Mandatory
		false,     // Immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)

	// If an error occurs, it could be due to the channel being closed or a connection issue
	if err != nil {
		log.Printf("Error publishing message: %v", err)

		// Handle specific error indicating the channel is closed
		if err.Error() == "channel/connection is not open" {
			log.Println("RabbitMQ channel/connection is not open, attempting to reconnect...")
			if err := reconnectRabbitMQ(); err != nil {
				return fmt.Errorf("failed to reconnect to RabbitMQ: %v", err)
			}

			// Retry publishing after reconnect
			err = rabbitmq.RabbitMQChannel.Publish(
				"",        // Default exchange
				queueName, // Routing key (queue name)
				false,     // Mandatory
				false,     // Immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(message),
				},
			)

			if err != nil {
				return fmt.Errorf("failed to publish a message after reconnect: %v", err)
			}
		}
	}

	fmt.Println("Message published to RabbitMQ queue:", queueName)
	return nil
}

func reconnectRabbitMQ() error {
	// Close the existing connection and channel if they are open
	if rabbitmq.RabbitMQChannel != nil {
		rabbitmq.RabbitMQChannel.Close()
	}
	if rabbitmq.RabbitMQConn != nil {
		rabbitmq.RabbitMQConn.Close()
	}

	// Reinitialize RabbitMQ connection and channel
	rabbitmq.InitRabbitMQ() // No need to assign it to a variable

	return nil
}
