package googleapi

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

var (
	ErrMissingHistoryID = errors.New("missing historyID")
)

type GmailService struct {
	Client *gmail.Service
}

func NewGmailService(ctx context.Context, gClient *GoogleClient) (*GmailService, error) {
	client := gClient.Config.Client(ctx, gClient.token)

	service, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("Error creating gmail service: %v", err)
	}

	return &GmailService{
		Client: service,
	}, nil
}

func (gmailService *GmailService) CreateWatch(ctx context.Context, topicName string) (uint64, int64, error) {
	req := &gmail.WatchRequest{
		LabelIds:  []string{"INBOX"},
		TopicName: topicName,
	}

	res, err := gmailService.Client.Users.Watch("me", req).Do()
	if err != nil {
		return 0, 0, fmt.Errorf("Error creando watch: %v", err)
	}

	return res.HistoryId, res.Expiration, nil
}

func (gmailService *GmailService) DeleteWatch() error {
	return gmailService.Client.Users.Stop("me").Do()
}
