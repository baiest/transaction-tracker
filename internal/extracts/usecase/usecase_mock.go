package usecase

import (
	"context"

	accountsDomain "transaction-tracker/internal/accounts/domain"
	extractsDomain "transaction-tracker/internal/extracts/domain"

	"github.com/stretchr/testify/mock"
)

// MockExtractsUsecase es un mock de la interfaz ExtractsUsecase.
type MockExtractsUsecase struct {
	mock.Mock
}

func (m *MockExtractsUsecase) GetByMessageID(ctx context.Context, messageID string) (*extractsDomain.Extract, error) {
	args := m.Called(ctx, messageID)
	if extract, ok := args.Get(0).(*extractsDomain.Extract); ok {
		return extract, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockExtractsUsecase) GetExtractMessages(ctx context.Context, bankName string, account *accountsDomain.Account) ([]string, error) {
	args := m.Called(ctx, bankName, account)
	if msgs, ok := args.Get(0).([]string); ok {
		return msgs, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockExtractsUsecase) Save(ctx context.Context, extract *extractsDomain.Extract) error {
	args := m.Called(ctx, extract)
	return args.Error(0)
}

func (m *MockExtractsUsecase) Update(ctx context.Context, extract *extractsDomain.Extract) error {
	args := m.Called(ctx, extract)
	return args.Error(0)
}
