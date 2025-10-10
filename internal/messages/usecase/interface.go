package usecase

import (
	"context"
	accountsDomain "transaction-tracker/internal/accounts/domain"
	"transaction-tracker/internal/messages/domain"
)

type MessageUsecase interface {
	GetMessage(ctx context.Context, id string, accountID string) (*domain.Message, error)
	Process(ctx context.Context, notificationID string, externalID string, account *accountsDomain.Account) (*domain.Message, error)
	ProcessByNotification(ctx context.Context, account *accountsDomain.Account, historyID uint64) ([]*domain.Message, error)
	GetMessageIDsByNotificationID(ctx context.Context, historyID uint64, account *accountsDomain.Account) ([]string, error)
}
