package usecase

import (
	"context"
	"transaction-tracker/api/services/accounts"
	"transaction-tracker/internal/messages/domain"
)

type MessageUsecase interface {
	GetMessageByIDAndAccountID(ctx context.Context, id string, accountID string) (*domain.Message, error)
	Process(ctx context.Context, notificationID string, externalID string, account *accounts.Account) (*domain.Message, error)
}
