package usecase

import (
	"context"
	accountsDomain "transaction-tracker/internal/accounts/domain"
	"transaction-tracker/internal/messages/domain"
)

type MessageUsecase interface {
	GetMessageByIDAndAccountID(ctx context.Context, id string, accountID string) (*domain.Message, error)
	Process(ctx context.Context, notificationID string, externalID string, account *accountsDomain.Account) (*domain.Message, error)
}
