package usecase

import (
	"context"
	accountsDomain "transaction-tracker/internal/accounts/domain"
	"transaction-tracker/internal/messages/domain"
)

// MockMessageUsecase is a mock implementation of the MovementUsecase interface.
type MockMessageUsecase struct {
	GetMessageByIDAndAccountIDFunc func(ctx context.Context, id string, accountID string) (*domain.Message, error)
	ProcessFn                      func(ctx context.Context, notificationID string, externalID string, account *accountsDomain.Account) (*domain.Message, error)
}

// GetMessageByIDAndAccountID calls the mocked GetMessageByIDAndAccountIDFunc.
func (m *MockMessageUsecase) GetMessageByIDAndAccountID(ctx context.Context, id string, accountID string) (*domain.Message, error) {
	return m.GetMessageByIDAndAccountIDFunc(ctx, id, accountID)
}

func (m *MockMessageUsecase) Process(ctx context.Context, notificationID string, externalID string, account *accountsDomain.Account) (*domain.Message, error) {
	return m.ProcessFn(ctx, notificationID, externalID, account)
}
