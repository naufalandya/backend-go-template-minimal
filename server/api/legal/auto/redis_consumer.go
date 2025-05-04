package auto

import (
	"context"
	"fmt"
	"log"
	rdb "modular_monolith/server/config/redis"
)

func SubscribeToChannel(channel string) {
	ctx := context.Background()
	sub := rdb.RedisClient.Subscribe(ctx, channel)
	defer sub.Close()

	for {
		msg, err := sub.ReceiveMessage(ctx)
		if err != nil {
			log.Printf("Error receiving message from Redis: %v", err)
			return
		}
		fmt.Printf("Received message: %s from channel: %s\n", msg.Payload, channel)

	}
}
