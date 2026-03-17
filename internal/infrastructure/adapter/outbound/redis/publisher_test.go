package redis

import (
	"context"
	"errors"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedisPublisher_Publish(t *testing.T) {
	ctx := context.Background()
	channel := "test-channel"
	message := "test-message"

	t.Run("should publish message successfully", func(t *testing.T) {
		db, mock := redismock.NewClientMock()
		publisher := NewRedisPublisher(db)

		mock.ExpectPublish(channel, message).SetVal(1)

		err := publisher.Publish(ctx, channel, message)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when publish fails", func(t *testing.T) {
		db, mock := redismock.NewClientMock()
		publisher := NewRedisPublisher(db)

		mock.ExpectPublish(channel, message).SetErr(errors.New("redis error"))

		err := publisher.Publish(ctx, channel, message)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to publish message to redis")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
