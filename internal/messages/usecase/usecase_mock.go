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
