package usecase

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"
	"transaction-tracker/api/services/accounts"
	"transaction-tracker/api/services/gmail/models"
	"transaction-tracker/internal/messages/domain"
	"transaction-tracker/internal/messages/repository"
	movementRepository "transaction-tracker/internal/movements/repository"
	loggerModels "transaction-tracker/logger/models"
	"transaction-tracker/pkg/google"
	messageextractor "transaction-tracker/pkg/message-extractor"

	"google.golang.org/api/gmail/v1"
)

type EmailFilter struct {
	Email   string
	Subject string
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

	ErrMissingMessageID  = errors.New("missing message id")
	ErrHistoryNotFound   = errors.New("history not found")
	ErrMissingExternalID = errors.New("missing external id")
)

type messageUsecase struct {
	messageRepo  repository.MessageRepository
	movementRepo movementRepository.MovementRepository
	gmailService *google.GmailService
	log          *loggerModels.Logger
}

// NewMessageUsecase is the constructor for the use case implementation.
// It receives a repository interface as a dependency.
func NewMessageUsecase(ctx context.Context, log *loggerModels.Logger, repo repository.MessageRepository, movementsRepo movementRepository.MovementRepository) MessageUsecase {
	return &messageUsecase{
		messageRepo:  repo,
		movementRepo: movementsRepo,
		log:          log,
	}
}

func (u *messageUsecase) GetMessageByIDAndAccountID(ctx context.Context, id string, accountID string) (*domain.Message, error) {
	return u.messageRepo.GetMessageByID(ctx, id, accountID)
}

func (u *messageUsecase) Process(ctx context.Context, notificationID string, externalID string, account *accounts.Account) (*domain.Message, error) {
	if externalID == "" {
		return nil, ErrMissingExternalID
	}

	googleClient, err := google.NewGoogleClient(ctx, account)
	if err != nil {
		return nil, err
	}

	_, err = googleClient.RefreshToken(ctx)
	if err != nil {
		return nil, err
	}

	gmailService, err := google.NewGmailService(ctx, googleClient)
	if err != nil {
		return nil, err
	}

	u.gmailService = gmailService

	message, err := u.messageRepo.GetMessageByExternalID(ctx, externalID, account.ID)
	if err != nil && !errors.Is(err, repository.ErrMessageNotFound) {
		return nil, err
	}

	if message != nil && (message.Status == domain.Success || message.Status == domain.Pending) {
		return message, nil
	}

	if message != nil {
		return u.filterAndUpdateMovement(ctx, message, nil)
	}

	gmailMessage, err := u.gmailService.GetMessageByID(ctx, externalID)
	if err != nil {
		return nil, err
	}

	from := ""
	to := ""
	date := time.Now()

	if gmailMessage != nil && gmailMessage.Payload != nil {
		for _, h := range gmailMessage.Payload.Headers {
			if h.Name == "From" {
				from = h.Value
			}

			if h.Name == "To" {
				from = h.Value
			}

			if h.Name == "Date" {
				parsed, err := mail.ParseDate(h.Value)
				if err != nil {
					continue
				}

				date = parsed
			}
		}
	}

	message = domain.NewMessage(account.ID, from, to, gmailMessage.Id, notificationID, notificationID, date)

	u.messageRepo.SaveMessage(ctx, message)

	return u.filterAndUpdateMovement(ctx, message, gmailMessage)
}

func (u *messageUsecase) filterAndUpdateMovement(ctx context.Context, message *domain.Message, emailMessage *gmail.Message) (*domain.Message, error) {
	if message.Status != domain.Pending {
		message.Status = domain.Pending

		errUpdate := u.messageRepo.UpdateMessage(ctx, message)
		if errUpdate != nil {
			u.log.Error(loggerModels.LogProperties{
				Event: "update_message_failed",
				Error: errUpdate,
			})

			return nil, errUpdate
		}
	}

	var err error

	if emailMessage == nil {
		emailMessage, err = u.gmailService.GetMessageByID(ctx, message.ExternalID)
		if err != nil && strings.Contains(err.Error(), "not found") {
			u.log.Info(loggerModels.LogProperties{
				Event: "message_not_found",
				Error: err,
			})

			message.Status = domain.Success

			return u.updateMessage(ctx, message, err)
		}

		if err != nil {
			u.log.Error(loggerModels.LogProperties{
				Event: "get_message_from_gmail_failed",
				Error: err,
			})

			message.Status = domain.Failure

			return u.updateMessage(ctx, message, err)
		}
	}

	messageType, isSupported := isMessageFiltered(emailMessage)

	if !isSupported {
		u.log.Error(loggerModels.LogProperties{
			Event: "message_not_supported",
			Error: err,
		})

		message.Status = domain.Success

		return u.updateMessage(ctx, message, err)
	}

	body := ""
	if len(emailMessage.Payload.Parts) > 0 {
		body = emailMessage.Payload.Parts[0].Body.Data
	} else {
		body = emailMessage.Payload.Body.Data
	}

	decodedBody, err := base64.URLEncoding.DecodeString(body)
	if err != nil {
		u.log.Error(loggerModels.LogProperties{
			Event: "decode_message_body_failed",
			Error: err,
		})

		message.Status = domain.Failure

		return u.updateMessage(ctx, message, err)
	}

	mvmExtractor, err := messageextractor.NewMovementExtractor("davivienda", string(decodedBody), messageType)
	if err != nil {
		u.log.Error(loggerModels.LogProperties{
			Event: "create_extractor_failed",
			Error: err,
		})

		message.Status = domain.Failure

		return u.updateMessage(ctx, message, err)
	}

	movements, err := mvmExtractor.Extract()
	if err != nil {
		u.log.Error(loggerModels.LogProperties{
			Event: "extract_movements_failed",
			Error: err,
		})

		message.Status = domain.Failure

		fmt.Println("fallo extrayendo")

		return u.updateMessage(ctx, message, err)
	}

	for _, m := range movements {
		m.AccountID = message.AccountID
		m.MessageID = message.ID

		err := u.movementRepo.CreateMovement(ctx, m)
		if err != nil {
			u.log.Error(loggerModels.LogProperties{
				Event: "create_movement_failed",
				Error: err,
			})

			message.Status = domain.Failure
		}
	}

	if message.Status == domain.Pending {
		message.Status = domain.Success
	}

	return u.updateMessage(ctx, message, err)
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

func (u *messageUsecase) updateMessage(ctx context.Context, message *domain.Message, err error) (*domain.Message, error) {
	errUpdate := u.messageRepo.UpdateMessage(ctx, message)
	if errUpdate != nil {
		u.log.Error(loggerModels.LogProperties{
			Event: "update_message_failed",
			Error: errUpdate,
		})

		return nil, errUpdate
	}

	return message, err
}
