package repository

import (
	"context"
	"transaction-tracker/internal/accounts/domain"
	"transaction-tracker/pkg/google"

	"github.com/stretchr/testify/mock"
)

// MockAccountsRepository is a mock of AccountsRepository
type MockAccountsRepository struct {
	mock.Mock
}

func (m *MockAccountsRepository) CreateAccount(ctx context.Context, account *domain.Account) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}

func (m *MockAccountsRepository) GetAccount(ctx context.Context, accountID string) (*domain.Account, error) {
	args := m.Called(ctx, accountID)
	return args.Get(0).(*domain.Account), args.Error(1)
}

func (m *MockAccountsRepository) GetAccountByEmail(ctx context.Context, email string) (*domain.Account, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*domain.Account), args.Error(1)
}

func (m *MockAccountsRepository) SaveGoogleAccount(ctx context.Context, accountID string, googleAccount *google.GoogleAccount) error {
	args := m.Called(ctx, accountID, googleAccount)
	return args.Error(0)
}

func (m *MockAccountsRepository) UpdateAccount(ctx context.Context, account *domain.Account) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}

func (m *MockAccountsRepository) DeleteWatcher(ctx context.Context, accountID string) error {
	args := m.Called(ctx, accountID)
	return args.Error(0)
}

func (m *MockAccountsRepository) CreateWatcher(ctx context.Context, accountID string) error {
	args := m.Called(ctx, accountID)
	return args.Error(0)
}
