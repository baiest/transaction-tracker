package usecase

import (
	"context"
	accountsDomain "transaction-tracker/internal/accounts/domain"
	"transaction-tracker/internal/messages/domain"

	"github.com/stretchr/testify/mock"
)

// MockMessageUsecase is a mock implementation of the MovementUsecase interface.
type MockMessageUsecase struct {
	mock.Mock
}

// GetMessageByIDAndAccountID calls the mocked GetMessageByIDAndAccountIDFunc.
func (m *MockMessageUsecase) GetMessageByIDAndAccountID(ctx context.Context, id string, accountID string) (*domain.Message, error) {
	args := m.Called(ctx, id, accountID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.Message), args.Error(1)
}

func (m *MockMessageUsecase) Process(ctx context.Context, notificationID string, externalID string, account *accountsDomain.Account) (*domain.Message, error) {
	args := m.Called(ctx, notificationID, externalID, account)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.Message), args.Error(1)
}

func (m *MockMessageUsecase) ProcessByNotification(ctx context.Context, account *accountsDomain.Account, historyID uint64) ([]*domain.Message, error) {
	args := m.Called(ctx, account, historyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*domain.Message), args.Error(1)
}

func (m *MockMessageUsecase) GetMessage(ctx context.Context, id string, accountID string) (*domain.Message, error) {
	args := m.Called(ctx, id, accountID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.Message), args.Error(1)
}

func (m *MockMessageUsecase) GetMessageIDsByNotificationID(ctx context.Context, historyID uint64, account *accountsDomain.Account) ([]string, error) {
	args := m.Called(ctx, historyID, account)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]string), args.Error(1)
}
