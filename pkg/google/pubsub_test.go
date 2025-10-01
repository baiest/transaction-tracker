package google

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewGooglePubSub(t *testing.T) {
	c := require.New(t)

	// This test will fail if you don't have a valid credentials file
	// and project ID. We are expecting an error here.
	_, err := NewGooglePubSub(context.Background(), "invalid-project-id", "invalid-credentials-file.json")
	c.Error(err)
}

func TestGetSubscription(t *testing.T) {
	c := require.New(t)

	// This test will fail if you don't have a valid credentials file
	// and project ID. We are expecting an error here.
	pubsub, err := NewGooglePubSub(context.Background(), "invalid-project-id", "invalid-credentials-file.json")
	c.Error(err)

	if pubsub != nil {
		sub, err := pubsub.GetSubscription(context.Background(), "test-sub")
		c.NoError(err)
		c.NotNil(sub)
	}
}
