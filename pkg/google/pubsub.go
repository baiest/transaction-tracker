package google

import (
	"context"
	"fmt"
	"strings"

	pubsub "cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
	"google.golang.org/api/option"
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
func NewGooglePubSub(ctx context.Context, projectID string, credentialsFile string) (*GooglePubSub, error) {
	client, err := pubsub.NewClient(ctx, projectID, option.WithCredentialsFile(credentialsFile))
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

	subscription, err := g.GetSubscription(ctx, subscriptionID)
	if err != nil {
		return err
	}

	sub := g.client.Subscriber(subscription.GetName())

	return sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		if err := handler(ctx, msg.Data); err != nil {
			if !strings.Contains(err.Error(), "googleapi: Error 404: Requested entity was not found") {
				msg.Nack()

				return
			}
		}

		msg.Ack()
	})
}

// GetSubscription returns a subscription
func (g *GooglePubSub) GetSubscription(ctx context.Context, subscriptionID string) (*pubsubpb.Subscription, error) {
	return g.client.SubscriptionAdminClient.GetSubscription(ctx, &pubsubpb.GetSubscriptionRequest{
		Subscription: fmt.Sprintf("projects/%s/subscriptions/%s", g.projectID, subscriptionID),
	})
}
