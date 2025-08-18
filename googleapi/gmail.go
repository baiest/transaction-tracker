package googleapi

import (
	"context"
	"errors"
	"fmt"
	"transaction-tracker/database/mongo/schemas"
	"transaction-tracker/googleapi/repository"

	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

var (
	ErrMissingHistoryID = errors.New("missing historyID")
)

type GmailService struct {
	Client                 *gmail.Service
	email                  string
	gmailRepository        *repository.GmailNotificationsRepository
	gmailMessageRepository *repository.GmailMessageRepository
}

func NewGmailService(ctx context.Context, gClient *GoogleClient) (*GmailService, error) {
	client := gClient.Config.Client(ctx, gClient.token)

	service, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("Error creating gmail service: %v", err)
	}

	notificationRepo, err := repository.NewGmailNotificationsRepository(ctx)
	if err != nil {
		return nil, err
	}

	messageRepo, err := repository.NewGmailMessageRepository(ctx)
	if err != nil {
		return nil, err
	}

	return &GmailService{
		Client:                 service,
		email:                  gClient.email,
		gmailRepository:        notificationRepo,
		gmailMessageRepository: messageRepo,
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

func (gmailService *GmailService) SaveNotification(ctx context.Context, notificationID string) (*schemas.GamilNotification, error) {
	notification := &schemas.GamilNotification{
		ID:     notificationID,
		Email:  gmailService.email,
		Status: "pending",
	}

	err := gmailService.gmailRepository.SaveNotification(ctx, notification)
	if errors.Is(err, repository.ErrNotificationAlreadyExists) {
		notification, err := gmailService.gmailRepository.GetNotificationByID(ctx, notification.ID)
		if err != nil {
			return nil, err
		}

		messages, err := gmailService.gmailMessageRepository.GetMessagesByNotificationID(ctx, notificationID)
		if err != nil {
			return nil, err
		}

		notification.Messages = messages

		return notification, err
	}

	if err != nil {
		return nil, err
	}

	return notification, nil
}

func (gmailService *GmailService) SaveMessage(ctx context.Context, message *schemas.Message) error {
	return gmailService.gmailMessageRepository.SaveMessage(ctx, message)
}

func (gmailService *GmailService) GetMessagesByNotificationID(ctx context.Context, notificationID string) ([]*schemas.Message, error) {
	return gmailService.gmailMessageRepository.GetMessagesByNotificationID(ctx, notificationID)
}

func (gmailService *GmailService) UpdateNotification(ctx context.Context, notification *schemas.GamilNotification) error {
	return gmailService.gmailRepository.UpdateNotification(ctx, notification)
}

func (gmailService *GmailService) UpdateMessage(ctx context.Context, message *schemas.Message) error {
	return gmailService.gmailMessageRepository.UpdateMessage(ctx, message)
}
