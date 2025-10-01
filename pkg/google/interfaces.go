package google

import (
	"context"
	"time"

	"cloud.google.com/go/pubsub"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
)

// GoogleClientAPI defines the interface for GoogleClient.
type GoogleClientAPI interface {
	SaveTokenAndInitServices(ctx context.Context, code string) (*GoogleAccount, error)
	GetAuthURL() string
	SetToken(token *oauth2.Token)
	GetUserEmail(ctx context.Context) (string, error)
	GmailService(ctx context.Context, googleAccount *GoogleAccount) (GmailAPI, error)
	RefreshToken(ctx context.Context, googleAccount *GoogleAccount) (*oauth2.Token, error)
	Config() *oauth2.Config
}

// GmailAPI defines the interface for GmailService.
type GmailAPI interface {
	CreateWatch(ctx context.Context, topicName string) (uint64, int64, error)
	DeleteWatch() error
	GetMessageByID(ctx context.Context, messageID string) (*gmail.Message, error)
	GetMessageAttachment(ctx context.Context, messageID string, attachmentID string) (*gmail.MessagePartBody, error)
	GetExtractMessages(ctx context.Context, bankName string) (*gmail.ListMessagesResponse, error)
	DownloadAttachments(ctx context.Context, accountID string, messageID string) (time.Month, int, string, error)
}

// PubSubAPI defines the interface for GooglePubSub.
type PubSubAPI interface {
	Publish(ctx context.Context, topic string, data []byte) error
	Subscribe(ctx context.Context, subscription string, handler func(ctx context.Context, msg []byte) error) error
	GetSubscription(ctx context.Context, subscriptionID string) (*pubsub.Subscription, error)
}
