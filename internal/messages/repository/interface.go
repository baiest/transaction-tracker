package repository

import (
	"context"
	"transaction-tracker/internal/messages/domain"
)

type MessageRepository interface {
	SaveMessage(ctx context.Context, message *domain.Message) error
	GetMessageByExternalID(ctx context.Context, id string, accountID string) (*domain.Message, error)
	GetMessageByID(ctx context.Context, id string, accountID string) (*domain.Message, error)
	GetMessagesByNotificationID(ctx context.Context, notificationID string) ([]*domain.Message, error)
	UpdateMessage(ctx context.Context, message *domain.Message) error
}
