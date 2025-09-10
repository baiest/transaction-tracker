package gmail

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"transaction-tracker/api/services/gmail/models"
	"transaction-tracker/api/services/gmail/transformers"
	"transaction-tracker/api/services/movements"
	"transaction-tracker/database/mongo/schemas"
	"transaction-tracker/googleapi"

	"google.golang.org/api/gmail/v1"
)

type EmailFilter struct {
	Email   string
	Subject string
}

type GmailService struct {
	movementService *movements.MovementsService
	service         *googleapi.GmailService
}

var (
	emailFilters = []EmailFilter{
		{
			Email: "banco_davivienda@davivienda.com",
		},
		{
			Email: "bancodavivienda@davivienda.com",
		},
		{
			Email: "juanballesteros2001@gmail.com",
		},
	}

	ErrMissingMessageID = errors.New("missing message id")
	ErrHistoryNotFound  = errors.New("history not found")
)

func NewGmailService(ctx context.Context) (*GmailService, error) {
	gClient, err := googleapi.NewGoogleClient(ctx)
	if err != nil {
		return nil, err
	}

	service, err := gClient.GmailService(ctx)
	if err != nil {
		return nil, err
	}

	movementsService, err := movements.NewMovementsService(ctx)
	if err != nil {
		return nil, err
	}

	return &GmailService{service: service, movementService: movementsService}, nil
}

func (g *GmailService) ProcessMessage(ctx context.Context, messageID string, notificationID string) (*schemas.Message, error) {
	if messageID == "" {
		return nil, ErrMissingMessageID
	}

	message, err := g.service.GetMessage(ctx, messageID)
	if errors.Is(err, googleapi.ErrMessageNotFound) {
		message = &schemas.Message{Status: "pending", ID: messageID, NotificationID: notificationID}

		err := g.service.SaveMessage(ctx, message)
		if err != nil {
			return nil, err
		}

		return g.filterAndCreateMovement(ctx, message)
	}

	if err != nil {
		return nil, err
	}

	if message.Status == "failure" {
		return g.filterAndCreateMovement(ctx, message)
	}

	return message, nil
}

func (g *GmailService) filterAndCreateMovement(ctx context.Context, message *schemas.Message) (*schemas.Message, error) {
	message.Status = "pending"

	errUpdate := g.service.UpdateMessage(ctx, message)
	if errUpdate != nil {
		return nil, errUpdate
	}

	msg, err := g.service.GetMessageByID(ctx, message.ID)
	if err != nil && strings.Contains(err.Error(), "not found") {
		message.Status = "success"

		errUpdate := g.service.UpdateMessage(ctx, message)
		if errUpdate != nil {
			return nil, errUpdate
		}

		return nil, nil
	}

	if err != nil {
		message.Status = "failure"

		errUpdate := g.service.UpdateMessage(ctx, message)
		if errUpdate != nil {
			return nil, errUpdate
		}

		return nil, err
	}

	messageType, isSupported := isMessageFiltered(msg)

	if !isSupported {
		message.Status = "success"

		err = g.service.UpdateMessage(ctx, message)
		if err != nil {
			return nil, err
		}

		return message, nil
	}

	transformer, err := transformers.NewMovementTransformer("davivienda", msg, messageType)
	if err != nil {
		message.Status = "failure"

		errUpdate := g.service.UpdateMessage(ctx, message)
		if errUpdate != nil {
			return nil, errUpdate
		}

		return nil, err
	}

	if messageType == models.Extract {
		extract, err := g.service.DownloadAttachments(ctx, msg.Id)
		if err != nil {
			return nil, err
		}

		transformer.SetExtract(extract)
	}

	movements, err := transformer.Excecute()
	if err != nil {
		message.Status = "failure"

		errUpdate := g.service.UpdateMessage(ctx, message)
		if errUpdate != nil {
			return nil, errUpdate
		}

		return nil, err
	}

	movementsChan := make(chan bool, len(movements))

	for _, m := range movements {
		go func() {
			defer func() {
				movementsChan <- true
			}()

			m.MessageID = message.ID
			m.Email = g.service.Email()

			err = g.movementService.CreateMovement(ctx, m)
			if err != nil {
				message.Status = "failure"

				g.service.UpdateMessage(ctx, message)
			}
		}()
	}

	for range movements {
		<-movementsChan
	}

	close(movementsChan)

	if message.Status != "pending" {
		return message, nil
	}

	message.Status = "success"

	err = g.service.UpdateMessage(ctx, message)
	if err != nil {
		return nil, err
	}

	return message, nil
}

// Review this
func isMessageFiltered(msg *gmail.Message) (models.MessageType, bool) {
	var from, subject string
	if msg.Payload != nil && msg.Payload.Headers != nil {
		for _, header := range msg.Payload.Headers {
			if header.Name == "From" {
				from = header.Value
			}

			if header.Name == "Subject" {
				subject = header.Value
			}
		}
	}

	messageType := models.Unknown

	if strings.Contains(strings.ToLower(subject), "extractos") {
		messageType = models.Extract
	}

	if strings.ToLower(subject) == "davivienda" {
		messageType = models.Movement
	}

	for _, filter := range emailFilters {
		if strings.Contains(strings.ToLower(from), filter.Email) {
			return messageType, true
		}
	}

	return messageType, false
}

func (g *GmailService) GetMessage(ctx context.Context, messageID string) (*schemas.Message, error) {
	if messageID == "" {
		return nil, ErrMissingMessageID
	}

	return g.service.GetMessage(ctx, messageID)
}

func (g *GmailService) CreateNotification(ctx context.Context, historyID uint64) (*schemas.GmailNotification, error) {
	strHistoryID := strconv.FormatInt(int64(historyID), 10)

	notification, err := g.service.SaveNotification(ctx, strHistoryID)
	if err != nil {
		return nil, err
	}

	if notification.Status != "pending" {
		return notification, nil
	}

	historyListCall := g.service.Client.Users.History.List("me").StartHistoryId(historyID)

	historyList, err := historyListCall.Do()
	if err != nil {
		return nil, err
	}

	if len(historyList.History) == 0 {
		return nil, ErrHistoryNotFound
	}

	messageIDs := map[string]bool{}

	for _, h := range historyList.History {
		for _, m := range h.Messages {
			if messageIDs[m.Id] {
				continue
			}

			messageIDs[m.Id] = true

			notification.Messages = append(notification.Messages, &schemas.Message{ID: m.Id, NotificationID: strHistoryID, Status: "pending"})
		}
	}

	return notification, nil
}

func (g *GmailService) UpdateNotification(ctx context.Context, notification *schemas.GmailNotification) error {
	return g.service.UpdateNotification(ctx, notification)
}

func (g *GmailService) GetExtractMessages(ctx context.Context, bankName string) (*gmail.ListMessagesResponse, error) {
	return g.service.GetExtractMessages(ctx, bankName)
}
