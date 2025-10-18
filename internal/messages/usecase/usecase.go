package usecase

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"
	accountsDomain "transaction-tracker/internal/accounts/domain"
	extractsDomain "transaction-tracker/internal/extracts/domain"
	extractsUsecase "transaction-tracker/internal/extracts/usecase"
	"transaction-tracker/internal/messages/domain"
	"transaction-tracker/internal/messages/repository"
	movementsUsecase "transaction-tracker/internal/movements/usecase"
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

	ErrMissingGmailService = errors.New("missing gmail service")
	ErrMissingMessageID    = errors.New("missing message id")
	ErrMissingExternalID   = errors.New("missing external id")
)

type messageUsecase struct {
	messageRepo    repository.MessageRepository
	mvmUsecase     movementsUsecase.MovementUsecase
	extractUsecase extractsUsecase.ExtractsUsecase
	googleClient   google.GoogleClientAPI
	log            *loggerModels.Logger
}

// NewMessageUsecase is the constructor for the use case implementation.
// It receives a repository interface as a dependency.
func NewMessageUsecase(ctx context.Context, googleClient google.GoogleClientAPI, repo repository.MessageRepository, mvmUsecase movementsUsecase.MovementUsecase, extractUsecase extractsUsecase.ExtractsUsecase) MessageUsecase {
	log := ctx.Value("logger").(*loggerModels.Logger)

	return &messageUsecase{
		messageRepo:    repo,
		mvmUsecase:     mvmUsecase,
		extractUsecase: extractUsecase,
		googleClient:   googleClient,
		log:            log,
	}
}

func (u *messageUsecase) GetMessage(ctx context.Context, id string, accountID string) (*domain.Message, error) {
	return u.messageRepo.GetMessageByID(ctx, id, accountID)
}

func (u *messageUsecase) Process(ctx context.Context, notificationID string, externalID string, account *accountsDomain.Account) (*domain.Message, error) {
	if externalID == "" {
		return nil, ErrMissingExternalID
	}

	client := u.googleClient.Client(ctx, account.GoogleAccount)

	gmailService, err := google.NewGmailClient(ctx, client)
	if err != nil {
		return nil, err
	}

	message, err := u.messageRepo.GetMessageByExternalID(ctx, externalID, account.ID)
	if err != nil && !errors.Is(err, repository.ErrMessageNotFound) {
		return nil, err
	}

	if message != nil && (message.Status == domain.Success || message.Status == domain.Pending) {
		return message, nil
	}

	if message != nil {
		err = u.filterAndUpdateMovement(ctx, gmailService, message, nil)
		if err != nil {
			message.Status = domain.Failure
		}

		return u.updateMessage(ctx, message, err)
	}

	gmailMessage, err := gmailService.GetMessageByID(ctx, externalID)
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
				to = h.Value
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

	err = u.filterAndUpdateMovement(ctx, gmailService, message, gmailMessage)
	if err != nil {
		message.Status = domain.Failure
	}

	return u.updateMessage(ctx, message, err)
}

func (u *messageUsecase) filterAndUpdateMovement(ctx context.Context, gmailService google.GmailAPI, message *domain.Message, emailMessage *gmail.Message) error {
	if message.Status != domain.Pending {
		message.Status = domain.Pending
	}

	var err error

	if emailMessage == nil {
		emailMessage, err = gmailService.GetMessageByID(ctx, message.ExternalID)
		if err != nil && strings.Contains(err.Error(), "not found") {
			u.log.Info(loggerModels.LogProperties{
				Event: "message_not_found",
				Error: err,
			})

			message.Status = domain.Success

			return nil
		}

		if err != nil {
			u.log.Error(loggerModels.LogProperties{
				Event: "get_message_from_gmail_failed",
				Error: err,
			})

			return err
		}
	}

	messageType, isSupported := isMessageFiltered(emailMessage)

	if !isSupported {
		u.log.Error(loggerModels.LogProperties{
			Event: "message_not_supported",
			Error: err,
		})

		message.Status = domain.Success

		return nil
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

		return err
	}

	mvmExtractor, err := messageextractor.NewMovementExtractor("davivienda", string(decodedBody), messageType)
	if err != nil {
		u.log.Error(loggerModels.LogProperties{
			Event: "create_extractor_failed",
			Error: err,
		})

		return err
	}

	if messageType == messageextractor.Extract {
		extract, err := u.GetExtract(ctx, gmailService, message)
		if err != nil {
			u.log.Error(loggerModels.LogProperties{
				Event: "download_attachments_failed",
				Error: err,
			})

			return err
		}

		err = u.extractUsecase.Update(ctx, extract)
		if err != nil {
			u.log.Error(loggerModels.LogProperties{
				Event: "update_extract_failed",
				Error: err,
			})

			return err
		}

		mvmExtractor.SetExtract(extract)
	}

	movements, err := mvmExtractor.Extract()
	if err != nil {
		u.log.Error(loggerModels.LogProperties{
			Event: "extract_movements_failed",
			Error: err,
		})

		return err
	}

	for _, m := range movements {
		m.AccountID = message.AccountID
		m.MessageID = message.ID

		err := u.mvmUsecase.CreateMovement(ctx, m)
		if err != nil {
			u.log.Error(loggerModels.LogProperties{
				Event: "create_movement_failed",
				Error: err,
			})

			if !errors.Is(err, movementsUsecase.ErrMustBeGreaterThanZero) {
				message.Status = domain.Failure
			}
		}
	}

	if message.Status == domain.Pending {
		message.Status = domain.Success
	}

	return nil
}

func (u *messageUsecase) GetExtract(ctx context.Context, gmailService google.GmailAPI, message *domain.Message) (*extractsDomain.Extract, error) {
	extract, err := u.extractUsecase.GetByMessageID(ctx, message.ID)
	if err != nil {
		if !errors.Is(err, extractsUsecase.ErrExtractNotFound) {
			return nil, err
		}
	}

	if extract != nil && (extract.Status == extractsDomain.ExtractStatusPending || extract.Status == extractsDomain.ExtractStatusProcessed) {
		return extract, nil
	}

	if extract == nil {
		extract = extractsDomain.NewExtract(message.AccountID, message.ID, "", "", 0, 0)

		err := u.extractUsecase.Save(ctx, extract)
		if err != nil {
			return nil, err
		}
	} else {
		err := u.mvmUsecase.DeleteMovementsByExtractID(ctx, extract.ID)
		if err != nil {
			return nil, err
		}
	}

	month, year, filePath, err := gmailService.DownloadAttachments(ctx, extract.AccountID, message.ExternalID)
	if err != nil {
		errUpdate := u.extractUsecase.Update(ctx, extract)
		if errUpdate != nil {
			return nil, errUpdate
		}

		return nil, err
	}

	extract.Month = month
	extract.Year = year
	extract.Path = filePath
	extract.Status = extractsDomain.ExtractStatusPending

	return extract, nil
}

// Review this
func isMessageFiltered(msg *gmail.Message) (messageextractor.MessageType, bool) {
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

	messageType := messageextractor.Unknown

	if strings.Contains(strings.ToLower(subject), "extractos") {
		messageType = messageextractor.Extract
	}

	if strings.ToLower(subject) == "davivienda" {
		messageType = messageextractor.Movement
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

func (u *messageUsecase) GetMessageIDsByNotificationID(ctx context.Context, historyID uint64, account *accountsDomain.Account) ([]string, error) {
	client := u.googleClient.Client(ctx, account.GoogleAccount)

	gmailService, err := google.NewGmailClient(ctx, client)
	if err != nil {
		return nil, err
	}

	messages, err := gmailService.GetMessagesByHistory(ctx, historyID)
	if err != nil {
		if errors.Is(err, google.ErrHistoryNotFound) {
			return []string{}, nil
		}

		return nil, err
	}

	messageIDs := []string{}

	for _, m := range messages {
		messageIDs = append(messageIDs, m.Id)
	}

	return messageIDs, nil
}

func (u *messageUsecase) ProcessByNotification(ctx context.Context, account *accountsDomain.Account, historyID uint64) ([]*domain.Message, error) {
	messages, err := u.GetMessageIDsByNotificationID(ctx, historyID, account)
	if err != nil {
		return nil, err
	}

	processedMessages := []*domain.Message{}

	for _, msgID := range messages {
		msg, err := u.Process(ctx, fmt.Sprintf("%d", historyID), msgID, account)
		if err != nil {
			return nil, err
		}

		processedMessages = append(processedMessages, msg)
	}

	return processedMessages, nil
}
