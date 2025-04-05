// ws/redis.go
package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"clipsync.com/m/db"
)

func PublishToRedis(msg Message) {
	ctx := context.Background()
	channel := fmt.Sprintf("clipboard_sync:user:%s", msg.UserID)

	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Println("Failed to marshal message for Redis:", err)
		return
	}

	err = db.RedisClient.Publish(ctx, channel, jsonData).Err()
	if err != nil {
		log.Println("Failed to publish to Redis:", err)
	}
}

func SubscribeToUserChannel(userID string, server *Server) {
	ctx := context.Background()
	channel := fmt.Sprintf("clipboard_sync:user:%s", userID)
	sub := db.RedisClient.Subscribe(ctx, channel)
	ch := sub.Channel()

	go func() {
		for {
			select {
			case msg, ok := <-ch:
				if !ok {
					log.Printf("Redis channel closed for user %s", userID)
					return
				}

				var incoming Message
				if err := json.Unmarshal([]byte(msg.Payload), &incoming); err != nil {
					log.Println("Failed to unmarshal Redis message:", err)
					continue
				}

				server.broadcast <- incoming

			case <-ctx.Done():
				log.Printf("Redis subscription context cancelled for user %s", userID)
				return
			}
		}
	}()
}
