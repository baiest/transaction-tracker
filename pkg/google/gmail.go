package google

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// EmailExtract defines the structure for holding email extraction information.
type EmailExtract struct {
	email   string
	subject string
}

const (
	maxGetMessageRetries = 5
)

var (
	emailByBank = map[string]EmailExtract{
		"davivienda": {
			email:   "bancodavivienda@davivienda.com",
			subject: "Extractos",
		},
	}

	// ErrMissingHistoryID is returned when a history ID is not found.
	ErrMissingHistoryID = errors.New("missing historyID")
	// ErrMessageNotFound is returned when a message is not found.
	ErrMessageNotFound = errors.New("message not found")

	extractsFolder = "files/%s/extracts"
)

// GmailService provides methods for interacting with the Gmail API.
type GmailService struct {
	Client *gmail.Service
}

// NewGmailClient creates a new GmailService with the provided http.Client.
func NewGmailClient(ctx context.Context, client *http.Client) (*GmailService, error) {
	service, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("error creating gmail service: %w", err)
	}

	return &GmailService{
		Client: service,
	}, nil
}

// CreateWatch sets up a watch on the user's inbox to receive push notifications.
func (g *GmailService) CreateWatch(ctx context.Context, topicName string) (uint64, int64, error) {
	req := &gmail.WatchRequest{
		LabelIds:  []string{"INBOX"},
		TopicName: topicName,
	}

	res, err := g.Client.Users.Watch("me", req).Do()
	if err != nil {
		return 0, 0, fmt.Errorf("err creating watch: %w", err)
	}

	return res.HistoryId, res.Expiration, nil
}

// DeleteWatch stops the watch on the user's inbox.
func (g *GmailService) DeleteWatch() error {
	return g.Client.Users.Stop("me").Do()
}

func (g *GmailService) getMessageWithRetries(ctx context.Context, messageID string, retries int) (*gmail.Message, error) {
	msg, err := g.Client.Users.Messages.Get("me", messageID).Format("full").Do()
	if err != nil {
		if strings.Contains(err.Error(), "Too many concurrent requests for user") && retries > 0 {
			time.Sleep(1 * time.Second)

			return g.getMessageWithRetries(ctx, messageID, retries-1)
		}

		return nil, err
	}

	return msg, nil
}

// GetMessageByID retrieves a specific message by its ID, with retries for handling concurrent requests.
func (g *GmailService) GetMessageByID(ctx context.Context, messageID string) (*gmail.Message, error) {
	return g.getMessageWithRetries(ctx, messageID, maxGetMessageRetries)
}

// GetMessageAttachment retrieves a specific attachment from a message.
func (g *GmailService) GetMessageAttachment(ctx context.Context, messageID string, attachmentID string) (*gmail.MessagePartBody, error) {
	return g.Client.Users.Messages.Attachments.Get("me", messageID, attachmentID).Do()
}

// GetExtractMessages retrieves a list of messages that match the criteria for bank extracts.
func (g *GmailService) GetExtractMessages(ctx context.Context, bankName string) (*gmail.ListMessagesResponse, error) {
	query := fmt.Sprintf("from:(%s) subject:(%s)", emailByBank[bankName].email, emailByBank[bankName].subject)

	return g.Client.Users.Messages.List("me").LabelIds("INBOX").Q(query).Do()
}

func (g *GmailService) getAttachmentData(ctx context.Context, messageID string, part *gmail.MessagePart) ([]byte, error) {
	att, err := g.GetMessageAttachment(ctx, messageID, part.Body.AttachmentId)
	if err != nil {
		return nil, fmt.Errorf("error getting attachment: %w", err)
	}

	data, err := base64.URLEncoding.DecodeString(att.Data)
	if err != nil {
		return nil, fmt.Errorf("error decoding attachment: %w", err)
	}

	return data, nil
}

func (g *GmailService) saveAttachment(accountID, fileName string, data []byte, perm fs.FileMode) (string, error) {
	// Sanitize accountID to prevent path traversal.
	safeAccountID := filepath.Clean(filepath.Base(accountID))

	exePath, _ := os.Getwd()
	currentDir := filepath.Dir(exePath)

	dirPath := filepath.Join(currentDir, fmt.Sprintf(extractsFolder, safeAccountID), fmt.Sprintf("%d", time.Now().Year()))

	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return "", fmt.Errorf("error creating directory: %w", err)
	}

	filePath := filepath.Join(dirPath, fileName)

	if err := os.WriteFile(filePath, data, perm); err != nil {
		return "", fmt.Errorf("error writing file: %w", err)
	}

	return filePath, nil
}

// DownloadAttachments downloads attachments from a specific message and saves them to the filesystem.
func (g *GmailService) DownloadAttachments(ctx context.Context, accountID string, messageID string) (time.Month, int, string, error) {
	msg, err := g.GetMessageByID(ctx, messageID)
	if err != nil {
		return 0, 0, "", fmt.Errorf("error getting message: %w", err)
	}

	for _, part := range msg.Payload.Parts {
		if part.Filename != "" && part.Body != nil && part.Body.AttachmentId != "" {
			data, err := g.getAttachmentData(ctx, messageID, part)
			if err != nil {
				return 0, 0, "", err
			}

			filePath, err := g.saveAttachment(accountID, part.Filename, data, 0600)
			if err != nil {
				return 0, 0, "", err
			}

			t := time.UnixMilli(msg.InternalDate)
			year := t.Year()

			if t.Month() == time.January {
				year = year - 1
			}

			return t.Month(), year, filePath, nil
		}
	}

	return 0, 0, "", errors.New("no attachment found")
}
