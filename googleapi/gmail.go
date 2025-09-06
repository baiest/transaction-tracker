package googleapi

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"transaction-tracker/database/mongo/schemas"
	"transaction-tracker/googleapi/repositories"

	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type EmailExtract struct {
	email   string
	subject string
}

var (
	emailByBank = map[string]EmailExtract{
		"davivienda": {
			email:   "bancodavivienda@davivienda.com",
			subject: "Extractos",
		},
	}

	ErrMissingHistoryID = errors.New("missing historyID")

	extractsFolder = "files/extracts"
)

type GmailService struct {
	Client                  *gmail.Service
	email                   string
	gmailRepository         *repositories.GmailNotificationsRepository
	gmailMessageRepository  *repositories.GmailMessageRepository
	gmailExtractsRepository *repositories.GmailExtractsRepository
}

func NewGmailService(ctx context.Context, gClient *GoogleClient) (*GmailService, error) {
	client := gClient.Config.Client(ctx, gClient.token)

	service, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("Error creating gmail service: %v", err)
	}

	notificationRepo, err := repositories.NewGmailNotificationsRepository(ctx)
	if err != nil {
		return nil, err
	}

	messageRepo, err := repositories.NewGmailMessageRepository(ctx)
	if err != nil {
		return nil, err
	}

	extractRepo, err := repositories.NewGmailExtractsRepository(ctx)
	if err != nil {
		return nil, err
	}

	return &GmailService{
		Client:                  service,
		email:                   gClient.email,
		gmailRepository:         notificationRepo,
		gmailMessageRepository:  messageRepo,
		gmailExtractsRepository: extractRepo,
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

func (gmailService *GmailService) SaveNotification(ctx context.Context, notificationID string) (*schemas.GmailNotification, error) {
	notification := &schemas.GmailNotification{
		ID:     notificationID,
		Email:  gmailService.email,
		Status: "pending",
	}

	err := gmailService.gmailRepository.SaveNotification(ctx, notification)
	if errors.Is(err, repositories.ErrNotificationAlreadyExists) {
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

func (gmailService *GmailService) UpdateNotification(ctx context.Context, notification *schemas.GmailNotification) error {
	return gmailService.gmailRepository.UpdateNotification(ctx, notification)
}

func (gmailService *GmailService) UpdateMessage(ctx context.Context, message *schemas.Message) error {
	return gmailService.gmailMessageRepository.UpdateMessage(ctx, message)
}

func (gmailService *GmailService) GetMessageByID(ctx context.Context, messageID string, retries int) (*gmail.Message, error) {
	msg, err := gmailService.Client.Users.Messages.Get("me", messageID).Format("full").Do()
	if err != nil {
		if strings.Contains(err.Error(), "Too many concurrent requests for user") && retries > 0 {
			time.Sleep(1 * time.Second)

			return gmailService.GetMessageByID(ctx, messageID, retries-1)
		}

		return nil, err
	}

	return msg, nil
}

func (gmailService *GmailService) GetMessageAttachment(ctx context.Context, messageID string, attachmentID string) (*gmail.MessagePartBody, error) {
	return gmailService.Client.Users.Messages.Attachments.Get("me", messageID, attachmentID).Do()
}

func (gmailService *GmailService) GetExtractMessages(ctx context.Context, bankName string) (*gmail.ListMessagesResponse, error) {
	query := fmt.Sprintf("from:(%s) subject:(%s)", emailByBank[bankName].email, emailByBank[bankName].subject)

	return gmailService.Client.Users.Messages.List("me").LabelIds("INBOX").Q(query).Do()
}

func (gmailService *GmailService) DownloadAttachments(ctx context.Context, messageID string) (*schemas.GmailExtract, error) {
	var extract *schemas.GmailExtract
	var err error

	extract, err = gmailService.gmailExtractsRepository.GetExtract(ctx, messageID)
	if err != nil && !errors.Is(err, repositories.ErrExtractNotFound) {
		return nil, fmt.Errorf("error getting extract: %w", err)
	}

	if extract != nil {
		return nil, nil
	}

	msg, err := gmailService.GetMessageByID(ctx, messageID, 3)
	if err != nil {
		return nil, fmt.Errorf("error getting message: %w", err)
	}

	for _, part := range msg.Payload.Parts {
		if part.Filename != "" && part.Body != nil && part.Body.AttachmentId != "" {
			att, err := gmailService.GetMessageAttachment(ctx, messageID, part.Body.AttachmentId)
			if err != nil {
				return nil, fmt.Errorf("error getting attachment: %w", err)
			}

			data, err := base64.URLEncoding.DecodeString(att.Data)
			if err != nil {
				return nil, fmt.Errorf("error decoding attachment: %w", err)
			}

			t := time.UnixMilli(msg.InternalDate)
			year := t.Year()

			if t.Month() == time.January {
				year = year - 1
			}

			dirPath := filepath.Join(extractsFolder, gmailService.email, fmt.Sprintf("%d", year))

			err = os.MkdirAll(dirPath, os.ModePerm)
			if err != nil {
				return nil, fmt.Errorf("error creando directorio: %w", err)
			}

			filePath := filepath.Join(dirPath, part.Filename)

			err = os.WriteFile(filePath, data, 0644)
			if err != nil {
				return nil, fmt.Errorf("error writing file: %w", err)
			}

			extract = &schemas.GmailExtract{
				ID:       messageID,
				Email:    gmailService.email,
				Date:     t,
				FilePath: filePath,
			}

			err = gmailService.gmailExtractsRepository.SaveExtract(ctx, extract)
			if err != nil {
				return nil, fmt.Errorf("error saving extract: %w", err)
			}
		}
	}

	return extract, nil
}
