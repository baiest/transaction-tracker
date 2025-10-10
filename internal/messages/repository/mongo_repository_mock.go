package repository

import (
	"context"

	"transaction-tracker/internal/messages/domain"

	"github.com/stretchr/testify/mock"
)

// MockMessageRepository is a testify/mock based mock for MessageRepository.
type MockMessageRepository struct {
	mock.Mock
}

func NewMockMessageRepository() *MockMessageRepository {
	return &MockMessageRepository{}
}

func (m *MockMessageRepository) SaveMessage(ctx context.Context, message *domain.Message) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

func (m *MockMessageRepository) GetMessageByExternalID(ctx context.Context, externalID, accountID string) (*domain.Message, error) {
	args := m.Called(ctx, externalID, accountID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Message), args.Error(1)
}

func (m *MockMessageRepository) UpdateMessage(ctx context.Context, message *domain.Message) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

func (m *MockMessageRepository) GetMessageByID(ctx context.Context, id, accountID string) (*domain.Message, error) {
	args := m.Called(ctx, id, accountID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Message), args.Error(1)
}

func (m *MockMessageRepository) GetMessagesByNotificationID(ctx context.Context, notificationID string) ([]*domain.Message, error) {
	args := m.Called(ctx, notificationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Message), args.Error(1)
}
