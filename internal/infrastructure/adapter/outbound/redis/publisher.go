package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// RedisPublisher implements the EventPublisher port using Redis Pub/Sub.
type RedisPublisher struct {
	client *redis.Client
}

// NewRedisPublisher creates a new implementation of the EventPublisher.
func NewRedisPublisher(client *redis.Client) *RedisPublisher {
	return &RedisPublisher{
		client: client,
	}
}

// Publish sends a message to a Redis channel.
func (p *RedisPublisher) Publish(ctx context.Context, channel string, message interface{}) error {
	err := p.client.Publish(ctx, channel, message).Err()
	if err != nil {
		return fmt.Errorf("failed to publish message to redis: %w", err)
	}
	return nil
}
