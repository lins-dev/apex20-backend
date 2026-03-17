package port

import "context"

// EventPublisher defines the contract for broadcasting messages to other services.
// This is an Outbound Port in Hexagonal Architecture.
type EventPublisher interface {
	// Publish sends a message to a specific channel/topic.
	Publish(ctx context.Context, channel string, message interface{}) error
}
