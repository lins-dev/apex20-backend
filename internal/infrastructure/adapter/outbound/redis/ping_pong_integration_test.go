package redis_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/apex20/backend/internal/infrastructure/adapter/outbound/redis"
	redis_lib "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedisPingPongIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()

	// 1. Get Redis URL from environment - MANDATORY
	redisURL := os.Getenv("REDIS_URL")
	require.NotEmpty(t, redisURL, "REDIS_URL environment variable is required for integration tests")

	// 2. Setup Clients
	opt, err := redis_lib.ParseURL(redisURL)
	require.NoError(t, err)
	client := redis_lib.NewClient(opt)
	defer client.Close()

	// Ensure connection is alive
	err = client.Ping(ctx).Err()
	require.NoError(t, err, "failed to connect to redis at %s", redisURL)

	publisher := redis.NewRedisPublisher(client)
	
	// 3. Setup Subscription (Simulating WS-Service)
	channelName := "ping-pong-test"
	pubsub := client.Subscribe(ctx, channelName)
	defer pubsub.Close()

	// Wait for subscription to be active
	_, err = pubsub.Receive(ctx)
	require.NoError(t, err)

	ch := pubsub.Channel()

	// 4. Publish Message (Backend Action)
	message := "ping-from-backend"
	err = publisher.Publish(ctx, channelName, message)
	assert.NoError(t, err)

	// 5. Receive and Verify (WS-Service Action)
	select {
	case msg := <-ch:
		assert.Equal(t, message, msg.Payload)
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for redis message")
	}
}
