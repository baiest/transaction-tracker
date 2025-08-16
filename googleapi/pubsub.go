package googleapi

import (
	"context"
	"fmt"

	pubsub "cloud.google.com/go/pubsub"
)

// PubSubService interface to use Pub/Sub service
type PubSubService interface {
	Publish(ctx context.Context, topic string, data []byte) error
	Subscribe(ctx context.Context, subscription string, handler func(ctx context.Context, msg []byte) error) error
}

type GooglePubSub struct {
	client    *pubsub.Client
	projectID string
}

// NewGooglePubSub creates a new GooglePubSub client
func NewGooglePubSub(ctx context.Context, projectID string) (*GooglePubSub, error) {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return &GooglePubSub{client: client, projectID: projectID}, nil
}

// Publish publishes a message to a topic
func (g *GooglePubSub) Publish(ctx context.Context, topic string, data []byte) error {
	return nil
}

// Subscribe subscribes to a topic and calls the handler function for each message
func (g *GooglePubSub) Subscribe(
	ctx context.Context,
	subscriptionID string,
	handler func(ctx context.Context, msg []byte) error,
) error {
	return g.client.Subscription(subscriptionID).Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		if err := handler(ctx, msg.Data); err != nil {
			fmt.Println("Error handling message:", err)
			msg.Nack()

			return
		}

		msg.Ack()
	})
}
