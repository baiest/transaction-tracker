package gmail

import (
	"context"
	"errors"
	"strings"
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
			Email:   "banco_davivienda@davivienda.com",
			Subject: "davivienda",
		},
		{
			Email:   "juanballesteros2001@gmail.com",
			Subject: "davivienda",
		},
	}

	ErrMissingMessaageID = errors.New("missing message id")
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

func (g *GmailService) ProcessMessage(ctx context.Context, messageID string) error {
	if messageID == "" {
		return ErrMissingMessaageID
	}

	message, err := g.service.GetMessage(ctx, messageID)
	if errors.Is(err, googleapi.ErrMessageNotFound) {
		message = &schemas.Message{Status: "pending", ID: messageID, NotificationID: ""}

		return g.processMessage(ctx, message)
	}

	if err != nil {
		return err
	}

	if message.Status == "pending" {
		return nil
	}

	if message.Status == "failure" {
		return g.processMessage(ctx, message)
	}

	return nil
}

func (g *GmailService) processMessage(ctx context.Context, message *schemas.Message) error {
	err := g.service.SaveMessage(ctx, message)
	if err != nil {
		return err
	}

	msg, err := g.service.GetMessageByID(ctx, message.ID)
	if err != nil {
		return err
	}

	if !isMessageFiltered(msg) {
		message.Status = "success"
		return g.service.UpdateMessage(ctx, message)
	}

	transformer, err := transformers.NewMovementTransformer("davivienda", msg)
	if err != nil {
		message.Status = "failure"
		errUpdate := g.service.UpdateMessage(ctx, message)
		if errUpdate != nil {
			return errUpdate
		}

		return err
	}

	movement, err := transformer.Excecute()
	if err != nil {
		message.Status = "failure"
		errUpdate := g.service.UpdateMessage(ctx, message)
		if errUpdate != nil {
			return errUpdate
		}

		return err
	}

	movement.MessageID = message.ID

	err = g.movementService.CreateMovement(ctx, movement)
	if err != nil {
		message.Status = "failure"
		errUpdate := g.service.UpdateMessage(ctx, message)
		if errUpdate != nil {
			return errUpdate
		}

		return err
	}

	message.Status = "success"
	return g.service.UpdateMessage(ctx, message)
}

// Review this
func isMessageFiltered(msg *gmail.Message) bool {
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

	for _, filter := range emailFilters {
		if strings.Contains(strings.ToLower(from), filter.Email) && strings.Contains(strings.ToLower(subject), filter.Subject) {
			return true
		}
	}

	return false
}

func (g *GmailService) GetMessage(ctx context.Context, messageID string) (*schemas.Message, error) {
	if messageID == "" {
		return nil, ErrMissingMessaageID
	}

	return g.service.GetMessage(ctx, messageID)
}
